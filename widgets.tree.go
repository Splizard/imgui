package imgui

func TreeNode(label string) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	return TreeNodeBehavior(window.GetIDs(label), 0, label)
}

func TreePushOverrideID(id ImGuiID) {
	var g = GImGui
	var window = g.CurrentWindow
	Indent(0)
	window.DC.TreeDepth++
	window.IDStack = append(window.IDStack, id)
}

// ~ Unindent()+PopId()
func TreePop() {
	var g = GImGui
	var window = g.CurrentWindow
	Unindent(0)

	window.DC.TreeDepth--
	var tree_depth_mask ImU32 = (1 << window.DC.TreeDepth)

	// Handle Left arrow to move to parent tree node (when ImGuiTreeNodeFlags_NavLeftJumpsBackHere is enabled)
	if g.NavMoveDir == ImGuiDir_Left && g.NavWindow == window && NavMoveRequestButNoResultYet() {
		if g.NavIdIsAlive && (window.DC.TreeJumpToParentOnPopMask&tree_depth_mask) != 0 {
			SetNavID(window.IDStack[len(window.IDStack)-1], g.NavLayer, 0, &ImRect{})
			NavMoveRequestCancel()
		}
	}
	window.DC.TreeJumpToParentOnPopMask &= tree_depth_mask - 1

	IM_ASSERT(len(window.IDStack) > 1) // There should always be 1 element in the IDStack (pushed during window creation). If this triggers you called TreePop/PopID too much.
	PopID()
}

// Consume previous SetNextItemOpen() data, if any. May return true when logging
func TreeNodeBehaviorIsOpen(id ImGuiID, flags ImGuiTreeNodeFlags) bool {
	if flags&ImGuiTreeNodeFlags_Leaf != 0 {
		return true
	}

	// We only write to the tree storage if the user clicks (or explicitly use the SetNextItemOpen function)
	var g = GImGui
	var window = g.CurrentWindow
	var storage = &window.DC.StateStorage

	var is_open bool
	if g.NextItemData.Flags&ImGuiNextItemDataFlags_HasOpen != 0 {
		if g.NextItemData.OpenCond&ImGuiCond_Always != 0 {
			is_open = g.NextItemData.OpenVal
			storage.SetInt(id, bool2int(is_open))
		} else {
			// We treat ImGuiCond_Once and ImGuiCond_FirstUseEver the same because tree node state are not saved persistently.
			var stored_value int = storage.GetInt(id, -1)
			if stored_value == -1 {
				is_open = g.NextItemData.OpenVal
				storage.SetInt(id, bool2int(is_open))
			} else {
				is_open = stored_value != 0
			}
		}
	} else {
		if (flags & ImGuiTreeNodeFlags_DefaultOpen) != 0 {
			is_open = storage.GetInt(id, 1) != 0
		} else {
			is_open = storage.GetInt(id, 0) != 0
		}
	}

	// When logging is enabled, we automatically expand tree nodes (but *NOT* collapsing headers.. seems like sensible behavior).
	// NB- If we are above max depth we still allow manually opened nodes to be logged.
	if g.LogEnabled && 0 == (flags&ImGuiTreeNodeFlags_NoAutoOpenOnLog) && (window.DC.TreeDepth-g.LogDepthRef) < g.LogDepthToExpand {
		is_open = true
	}

	return is_open
}

func TreeNodeBehavior(id ImGuiID, flags ImGuiTreeNodeFlags, label string) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style *ImGuiStyle = &g.Style
	var display_frame bool = (flags & ImGuiTreeNodeFlags_Framed) != 0
	var padding ImVec2
	if display_frame || (flags&ImGuiTreeNodeFlags_FramePadding) != 0 {
		padding = style.FramePadding
	} else {
		padding = ImVec2{style.FramePadding.x, ImMin(window.DC.CurrLineTextBaseOffset, style.FramePadding.y)}
	}

	var label_size ImVec2 = CalcTextSize(label, false, 0)

	// We vertically grow up to current line height up the typical widget height.
	var frame_height float = ImMax(ImMin(window.DC.CurrLineSize.y, g.FontSize+style.FramePadding.y*2), label_size.y+padding.y*2)
	var frame_bb ImRect
	if (flags & ImGuiTreeNodeFlags_SpanFullWidth) != 0 {
		frame_bb.Min.x = window.WorkRect.Min.x
	} else {
		frame_bb.Min.x = window.DC.CursorPos.x
	}
	frame_bb.Min.y = window.DC.CursorPos.y
	frame_bb.Max.x = window.WorkRect.Max.x
	frame_bb.Max.y = window.DC.CursorPos.y + frame_height
	if display_frame {
		// Framed header expand a little outside the default padding, to the edge of InnerClipRect
		// (FIXME: May remove this at some point and make InnerClipRect align with WindowPadding.x instead of WindowPadding.x*0.5f)
		frame_bb.Min.x -= IM_FLOOR(window.WindowPadding.x*0.5 - 1.0)
		frame_bb.Max.x += IM_FLOOR(window.WindowPadding.x * 0.5)
	}

	// Collapser arrow width + Spacing
	var text_offset_x float = g.FontSize
	if display_frame {
		text_offset_x += padding.x * 3
	} else {
		text_offset_x += padding.x * 2
	}

	var text_offset_y float = ImMax(padding.y, window.DC.CurrLineTextBaseOffset) // Latch before ItemSize changes it
	var text_width float = g.FontSize
	if label_size.x > 0.0 { // Include collapser
		text_width += padding.x * 2
	}
	var text_pos = ImVec2{window.DC.CursorPos.x + text_offset_x, window.DC.CursorPos.y + text_offset_y}
	ItemSizeVec(&ImVec2{text_width, frame_height}, padding.y)

	// For regular tree nodes, we arbitrary allow to click past 2 worth of ItemSpacing
	var interact_bb ImRect = frame_bb
	if !display_frame && (flags&(ImGuiTreeNodeFlags_SpanAvailWidth|ImGuiTreeNodeFlags_SpanFullWidth)) == 0 {
		interact_bb.Max.x = frame_bb.Min.x + text_width + label_size.x + style.ItemSpacing.x*2.0
	}

	// Store a flag for the current depth to tell if we will allow closing this node when navigating one of its child.
	// For this purpose we essentially compare if g.NavIdIsAlive went from 0 to 1 between TreeNode() and TreePop().
	// This is currently only support 32 level deep and we are fine with (1 << Depth) overflowing into a zero.
	var is_leaf bool = (flags & ImGuiTreeNodeFlags_Leaf) != 0
	var is_open bool = TreeNodeBehaviorIsOpen(id, flags)
	if is_open && !g.NavIdIsAlive && (flags&ImGuiTreeNodeFlags_NavLeftJumpsBackHere) != 0 && 0 == (flags&ImGuiTreeNodeFlags_NoTreePushOnOpen) {
		window.DC.TreeJumpToParentOnPopMask |= (1 << window.DC.TreeDepth)
	}

	var item_add bool = ItemAdd(&interact_bb, id, nil, 0)
	g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HasDisplayRect
	g.LastItemData.DisplayRect = frame_bb

	if !item_add {
		if is_open && 0 == (flags&ImGuiTreeNodeFlags_NoTreePushOnOpen) {
			TreePushOverrideID(id)
		}
		return is_open
	}

	var button_flags ImGuiButtonFlags = ImGuiButtonFlags_None
	if flags&ImGuiTreeNodeFlags_AllowItemOverlap != 0 {
		button_flags |= ImGuiButtonFlags_AllowItemOverlap
	}
	if !is_leaf {
		button_flags |= ImGuiButtonFlags_PressedOnDragDropHold
	}

	// We allow clicking on the arrow section with keyboard modifiers held, in order to easily
	// allow browsing a tree while preserving selection with code implementing multi-selection patterns.
	// When clicking on the rest of the tree node we always disallow keyboard modifiers.
	var arrow_hit_x1 float = (text_pos.x - text_offset_x) - style.TouchExtraPadding.x
	var arrow_hit_x2 float = (text_pos.x - text_offset_x) + (g.FontSize + padding.x*2.0) + style.TouchExtraPadding.x
	var is_mouse_x_over_arrow bool = (g.IO.MousePos.x >= arrow_hit_x1 && g.IO.MousePos.x < arrow_hit_x2)
	if window != g.HoveredWindow || !is_mouse_x_over_arrow {
		button_flags |= ImGuiButtonFlags_NoKeyModifiers
	}

	// Open behaviors can be altered with the _OpenOnArrow and _OnOnDoubleClick flags.
	// Some alteration have subtle effects (e.g. toggle on MouseUp vs MouseDown events) due to requirements for multi-selection and drag and drop support.
	// - Single-click on label = Toggle on MouseUp (default, when _OpenOnArrow=0)
	// - Single-click on arrow = Toggle on MouseDown (when _OpenOnArrow=0)
	// - Single-click on arrow = Toggle on MouseDown (when _OpenOnArrow=1)
	// - Double-click on label = Toggle on MouseDoubleClick (when _OpenOnDoubleClick=1)
	// - Double-click on arrow = Toggle on MouseDoubleClick (when _OpenOnDoubleClick=1 and _OpenOnArrow=0)
	// It is rather standard that arrow click react on Down rather than Up.
	// We set ImGuiButtonFlags_PressedOnClickRelease on OpenOnDoubleClick because we want the item to be active on the initial MouseDown in order for drag and drop to work.
	if is_mouse_x_over_arrow {
		button_flags |= ImGuiButtonFlags_PressedOnClick
	} else if flags&ImGuiTreeNodeFlags_OpenOnDoubleClick != 0 {
		button_flags |= ImGuiButtonFlags_PressedOnClickRelease | ImGuiButtonFlags_PressedOnDoubleClick
	} else {
		button_flags |= ImGuiButtonFlags_PressedOnClickRelease
	}

	var selected bool = (flags & ImGuiTreeNodeFlags_Selected) != 0
	var was_selected bool = selected

	var hovered, held bool
	var pressed bool = ButtonBehavior(&interact_bb, id, &hovered, &held, button_flags)
	var toggled bool = false
	if !is_leaf {
		if pressed && g.DragDropHoldJustPressedId != id {
			if (flags&(ImGuiTreeNodeFlags_OpenOnArrow|ImGuiTreeNodeFlags_OpenOnDoubleClick)) == 0 || (g.NavActivateId == id) {
				toggled = true
			}
			if flags&ImGuiTreeNodeFlags_OpenOnArrow != 0 {
				toggled = toggled || (is_mouse_x_over_arrow && !g.NavDisableMouseHover) // Lightweight equivalent of IsMouseHoveringRect() since ButtonBehavior() already did the job
			}
			if (flags&ImGuiTreeNodeFlags_OpenOnDoubleClick) != 0 && g.IO.MouseDoubleClicked[0] {
				toggled = true
			}
		} else if pressed && g.DragDropHoldJustPressedId == id {
			IM_ASSERT(button_flags&ImGuiButtonFlags_PressedOnDragDropHold != 0)
			if !is_open { // When using Drag and Drop "hold to open" we keep the node highlighted after opening, but never close it again.
				toggled = true
			}
		}

		if g.NavId == id && g.NavMoveDir == ImGuiDir_Left && is_open {
			toggled = true
			NavMoveRequestCancel()
		}
		if g.NavId == id && g.NavMoveDir == ImGuiDir_Right && !is_open { // If there's something upcoming on the line we may want to give it the priority?
			toggled = true
			NavMoveRequestCancel()
		}

		if toggled {
			is_open = !is_open
			window.DC.StateStorage.SetInt(id, bool2int(is_open))
			g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_ToggledOpen
		}
	}
	if flags&ImGuiTreeNodeFlags_AllowItemOverlap != 0 {
		SetItemAllowOverlap()
	}

	// In this branch, TreeNodeBehavior() cannot toggle the selection so this will never trigger.
	if selected != was_selected { //-V547
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_ToggledSelection
	}

	// Render
	var text_col ImU32 = GetColorU32FromID(ImGuiCol_Text, 1)
	var nav_highlight_flags ImGuiNavHighlightFlags = ImGuiNavHighlightFlags_TypeThin
	if display_frame {
		// Framed type
		var bg_col ImU32
		if held && hovered {
			bg_col = GetColorU32FromID(ImGuiCol_HeaderActive, 1)
		} else if hovered {
			bg_col = GetColorU32FromID(ImGuiCol_HeaderHovered, 1)
		} else {
			bg_col = GetColorU32FromID(ImGuiCol_Header, 1)
		}
		RenderFrame(frame_bb.Min, frame_bb.Max, bg_col, true, style.FrameRounding)
		RenderNavHighlight(&frame_bb, id, nav_highlight_flags)
		if flags&ImGuiTreeNodeFlags_Bullet != 0 {
			RenderBullet(window.DrawList, ImVec2{text_pos.x - text_offset_x*0.60, text_pos.y + g.FontSize*0.5}, text_col)
		} else if !is_leaf {
			var arrow ImGuiDir
			if is_open {
				arrow = ImGuiDir_Down
			} else {
				arrow = ImGuiDir_Right
			}
			RenderArrow(window.DrawList, ImVec2{text_pos.x - text_offset_x + padding.x, text_pos.y}, text_col, arrow, 1.0)
		} else { // Leaf without bullet, left-adjusted text
			text_pos.x -= text_offset_x
		}
		if flags&ImGuiTreeNodeFlags_ClipLabelForTrailingButton != 0 {
			frame_bb.Max.x -= g.FontSize + style.FramePadding.x
		}

		if g.LogEnabled {
			LogSetNextTextDecoration("###", "###")
		}
		RenderTextClipped(&text_pos, &frame_bb.Max, label, &label_size, nil, nil)
	} else {
		// Unframed typed for tree nodes
		if hovered || selected {
			var bg_col ImU32
			if held && hovered {
				bg_col = GetColorU32FromID(ImGuiCol_HeaderActive, 1)
			} else if hovered {
				bg_col = GetColorU32FromID(ImGuiCol_HeaderHovered, 1)
			} else {
				bg_col = GetColorU32FromID(ImGuiCol_Header, 1)
			}
			RenderFrame(frame_bb.Min, frame_bb.Max, bg_col, false, 0)
		}
		RenderNavHighlight(&frame_bb, id, nav_highlight_flags)
		if flags&ImGuiTreeNodeFlags_Bullet != 0 {
			RenderBullet(window.DrawList, ImVec2{text_pos.x - text_offset_x*0.5, text_pos.y + g.FontSize*0.5}, text_col)
		} else if !is_leaf {
			var arrow ImGuiDir
			if is_open {
				arrow = ImGuiDir_Down
			} else {
				arrow = ImGuiDir_Right
			}
			RenderArrow(window.DrawList, ImVec2{text_pos.x - text_offset_x + padding.x, text_pos.y + g.FontSize*0.15}, text_col, arrow, 0.70)
		}
		if g.LogEnabled {
			LogSetNextTextDecoration(">", "")
		}
		RenderText(text_pos, label, false)
	}

	if is_open && 0 == (flags&ImGuiTreeNodeFlags_NoTreePushOnOpen) {
		TreePushOverrideID(id)
	}
	return is_open
}
