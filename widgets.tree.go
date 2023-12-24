package imgui

import "fmt"

// Widgets: Trees
// - TreeNode functions return true when the node is open, in which case you need to also call TreePop() when you are finished displaying the tree node contents.

// helper variation to easily decorelate the id from the displayed string. Read the FAQ about why and how to use ID. to align arbitrary text at the same level as a TreeNode() you can use Bullet().
func TreeNodeF(str_id string, format string, args ...any) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	return TreeNodeBehavior(window.GetIDs(str_id), 0, fmt.Sprintf(format, args...))
}

func TreeNodeInterface(ptr_id any, format string, args ...any) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	return TreeNodeBehavior(window.GetIDInterface(ptr_id), 0, fmt.Sprintf(format, args...))
}

func TreeNodeEx(str_id string, flags ImGuiTreeNodeFlags, format string, args ...any) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	return TreeNodeBehavior(window.GetIDs(str_id), flags, fmt.Sprintf(format, args...))
}

func TreeNodeInterfaceEx(ptr_id any, flags ImGuiTreeNodeFlags, format string, args ...any) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	return TreeNodeBehavior(window.GetIDInterface(ptr_id), flags, fmt.Sprintf(format, args...))
}

// ~ Indent()+PushId(). Already called by TreeNode() when returning true, but you can call TreePush/TreePop yourself if desired.
func TreePush(str_id string) {
	window := GetCurrentWindow()
	Indent(0)
	window.DC.TreeDepth++
	if str_id != "" {
		PushString(str_id)
	} else {
		PushString("#TreePush")
	}
}

func TreePushInterface(ptr_id any) {
	window := GetCurrentWindow()
	Indent(0)
	window.DC.TreeDepth++
	if ptr_id != nil {
		PushInterface(ptr_id)
	} else {
		PushString("#TreePush")
	}
} // "

// horizontal distance preceding label when using TreeNode*() or Bullet() == (guiContext.FontSize + style.FramePadding.x*2) for a regular unframed TreeNode
func GetTreeNodeToLabelSpacing() float {
	return guiContext.FontSize + (guiContext.Style.FramePadding.x * 2.0)
}

// CollapsingHeader returns true when opened but do not indent nor push into the ID stack (because of the ImGuiTreeNodeFlags_NoTreePushOnOpen flag).
// This is basically the same as calling TreeNodeEx(label, ImGuiTreeNodeFlags_CollapsingHeader). You can remove the _NoTreePushOnOpen flag if you want behavior closer to normal TreeNode().
// if returning 'true' the header is open. doesn't indent nor push on ID stack. user doesn't have to call TreePop().
func CollapsingHeader(label string, flags ImGuiTreeNodeFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	return TreeNodeBehavior(window.GetIDs(label), flags|ImGuiTreeNodeFlags_CollapsingHeader, label)
}

// when 'p_visible != NULL': if '*p_visible==true' display an additional small close button on upper right of the header which will set the to bool false when clicked, if '*p_visible==false' don't display the header.
// p_visible == nil                        : regular collapsing header
// p_visible != nil && *p_visible == true  : show a small close button on the corner of the header, clicking the button will set *p_visible = false
// p_visible != nil && *p_visible == false : do not show the header at all
// Do not mistake this with the Open state of the header itself, which you can adjust with SetNextItemOpen() or ImGuiTreeNodeFlags_DefaultOpen.
func CollapsingHeaderVisible(label string, p_visible *bool, flags ImGuiTreeNodeFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	if p_visible != nil && !*p_visible {
		return false
	}

	var id = window.GetIDs(label)
	flags |= ImGuiTreeNodeFlags_CollapsingHeader
	if p_visible != nil {
		flags |= ImGuiTreeNodeFlags_AllowItemOverlap | ImGuiTreeNodeFlags_ClipLabelForTrailingButton
	}
	var is_open = TreeNodeBehavior(id, flags, label)
	if p_visible != nil {
		// Create a small overlapping close button
		// FIXME: We can evolve this into user accessible helpers to add extra buttons on title bars, headers, etc.
		// FIXME: CloseButton can overlap into text, need find a way to clip the text somehow.
		var last_item_backup = guiContext.LastItemData
		var button_size = guiContext.FontSize
		var button_x = max(guiContext.LastItemData.Rect.Min.x, guiContext.LastItemData.Rect.Max.x-guiContext.Style.FramePadding.x*2.0-button_size)
		var button_y = guiContext.LastItemData.Rect.Min.y
		var close_button_id = GetIDWithSeed("#CLOSE", id)
		if CloseButton(close_button_id, &ImVec2{button_x, button_y}) {
			*p_visible = false
		}
		guiContext.LastItemData = last_item_backup
	}

	return is_open
}

// set next TreeNode/CollapsingHeader open state.
func SetNextItemOpen(is_open bool, cond ImGuiCond) {
	if guiContext.CurrentWindow.SkipItems {
		return
	}
	guiContext.NextItemData.Flags |= ImGuiNextItemDataFlags_HasOpen
	guiContext.NextItemData.OpenVal = is_open
	if cond != 0 {
		guiContext.NextItemData.OpenCond = cond
	} else {
		guiContext.NextItemData.OpenCond = ImGuiCond_Always
	}
}

func TreeNode(label string) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	return TreeNodeBehavior(window.GetIDs(label), 0, label)
}

func TreePushOverrideID(id ImGuiID) {
	window := guiContext.CurrentWindow
	Indent(0)
	window.DC.TreeDepth++
	window.IDStack = append(window.IDStack, id)
}

// ~ Unindent()+PopId()
func TreePop() {
	window := guiContext.CurrentWindow
	Unindent(0)

	window.DC.TreeDepth--
	var tree_depth_mask ImU32 = (1 << window.DC.TreeDepth)

	// Handle Left arrow to move to parent tree node (when ImGuiTreeNodeFlags_NavLeftJumpsBackHere is enabled)
	if guiContext.NavMoveDir == ImGuiDir_Left && guiContext.NavWindow == window && NavMoveRequestButNoResultYet() {
		if guiContext.NavIdIsAlive && (window.DC.TreeJumpToParentOnPopMask&tree_depth_mask) != 0 {
			SetNavID(window.IDStack[len(window.IDStack)-1], guiContext.NavLayer, 0, &ImRect{})
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
	window := guiContext.CurrentWindow
	var storage = &window.DC.StateStorage

	var is_open bool
	if guiContext.NextItemData.Flags&ImGuiNextItemDataFlags_HasOpen != 0 {
		if guiContext.NextItemData.OpenCond&ImGuiCond_Always != 0 {
			is_open = guiContext.NextItemData.OpenVal
			storage.SetInt(id, bool2int(is_open))
		} else {
			// We treat ImGuiCond_Once and ImGuiCond_FirstUseEver the same because tree node state are not saved persistently.
			var stored_value = storage.GetInt(id, -1)
			if stored_value == -1 {
				is_open = guiContext.NextItemData.OpenVal
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
	if guiContext.LogEnabled && (flags&ImGuiTreeNodeFlags_NoAutoOpenOnLog) == 0 && (window.DC.TreeDepth-guiContext.LogDepthRef) < guiContext.LogDepthToExpand {
		is_open = true
	}

	return is_open
}

func TreeNodeBehavior(id ImGuiID, flags ImGuiTreeNodeFlags, label string) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var style = &guiContext.Style
	var display_frame = (flags & ImGuiTreeNodeFlags_Framed) != 0
	var padding ImVec2
	if display_frame || (flags&ImGuiTreeNodeFlags_FramePadding) != 0 {
		padding = style.FramePadding
	} else {
		padding = ImVec2{style.FramePadding.x, min(window.DC.CurrLineTextBaseOffset, style.FramePadding.y)}
	}

	label = FindRenderedTextEnd(label)
	var label_size = CalcTextSize(label, false, 0)

	// We vertically grow up to current line height up the typical widget height.
	var frame_height = max(min(window.DC.CurrLineSize.y, guiContext.FontSize+style.FramePadding.y*2), label_size.y+padding.y*2)
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
	var text_offset_x = guiContext.FontSize
	if display_frame {
		text_offset_x += padding.x * 3
	} else {
		text_offset_x += padding.x * 2
	}

	var text_offset_y = max(padding.y, window.DC.CurrLineTextBaseOffset) // Latch before ItemSize changes it
	var text_width = guiContext.FontSize
	if label_size.x > 0.0 { // Include collapser
		text_width += padding.x * 2
	}
	var text_pos = ImVec2{window.DC.CursorPos.x + text_offset_x, window.DC.CursorPos.y + text_offset_y}
	ItemSizeVec(&ImVec2{text_width, frame_height}, padding.y)

	// For regular tree nodes, we arbitrary allow to click past 2 worth of ItemSpacing
	var interact_bb = frame_bb
	if !display_frame && (flags&(ImGuiTreeNodeFlags_SpanAvailWidth|ImGuiTreeNodeFlags_SpanFullWidth)) == 0 {
		interact_bb.Max.x = frame_bb.Min.x + text_width + label_size.x + style.ItemSpacing.x*2.0
	}

	// Store a flag for the current depth to tell if we will allow closing this node when navigating one of its child.
	// For this purpose we essentially compare if guiContext.NavIdIsAlive went from 0 to 1 between TreeNode() and TreePop().
	// This is currently only support 32 level deep and we are fine with (1 << Depth) overflowing into a zero.
	var is_leaf = (flags & ImGuiTreeNodeFlags_Leaf) != 0
	var is_open = TreeNodeBehaviorIsOpen(id, flags)
	if is_open && !guiContext.NavIdIsAlive && (flags&ImGuiTreeNodeFlags_NavLeftJumpsBackHere) != 0 && flags&ImGuiTreeNodeFlags_NoTreePushOnOpen == 0 {
		window.DC.TreeJumpToParentOnPopMask |= (1 << window.DC.TreeDepth)
	}

	var item_add = ItemAdd(&interact_bb, id, nil, 0)
	guiContext.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HasDisplayRect
	guiContext.LastItemData.DisplayRect = frame_bb

	if !item_add {
		if is_open && flags&ImGuiTreeNodeFlags_NoTreePushOnOpen == 0 {
			TreePushOverrideID(id)
		}
		return is_open
	}

	var button_flags = ImGuiButtonFlags_None
	if flags&ImGuiTreeNodeFlags_AllowItemOverlap != 0 {
		button_flags |= ImGuiButtonFlags_AllowItemOverlap
	}
	if !is_leaf {
		button_flags |= ImGuiButtonFlags_PressedOnDragDropHold
	}

	// We allow clicking on the arrow section with keyboard modifiers held, in order to easily
	// allow browsing a tree while preserving selection with code implementing multi-selection patterns.
	// When clicking on the rest of the tree node we always disallow keyboard modifiers.
	var arrow_hit_x1 = (text_pos.x - text_offset_x) - style.TouchExtraPadding.x
	var arrow_hit_x2 = (text_pos.x - text_offset_x) + (guiContext.FontSize + padding.x*2.0) + style.TouchExtraPadding.x
	var is_mouse_x_over_arrow = (guiContext.IO.MousePos.x >= arrow_hit_x1 && guiContext.IO.MousePos.x < arrow_hit_x2)
	if window != guiContext.HoveredWindow || !is_mouse_x_over_arrow {
		button_flags |= ImGuiButtonFlags_NoKeyModifiers
	}

	// Open behaviors can be altered with the _OpenOnArrow and _OnOnDoubleClick flags.
	// Some alteration have subtle effects (e.guiContext. toggle on MouseUp vs MouseDown events) due to requirements for multi-selection and drag and drop support.
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

	var selected = (flags & ImGuiTreeNodeFlags_Selected) != 0
	var was_selected = selected

	var hovered, held bool
	var pressed = ButtonBehavior(&interact_bb, id, &hovered, &held, button_flags)
	var toggled = false
	if !is_leaf {
		if pressed && guiContext.DragDropHoldJustPressedId != id {
			if (flags&(ImGuiTreeNodeFlags_OpenOnArrow|ImGuiTreeNodeFlags_OpenOnDoubleClick)) == 0 || (guiContext.NavActivateId == id) {
				toggled = true
			}
			if flags&ImGuiTreeNodeFlags_OpenOnArrow != 0 {
				toggled = toggled || (is_mouse_x_over_arrow && !guiContext.NavDisableMouseHover) // Lightweight equivalent of IsMouseHoveringRect() since ButtonBehavior() already did the job
			}
			if (flags&ImGuiTreeNodeFlags_OpenOnDoubleClick) != 0 && guiContext.IO.MouseDoubleClicked[0] {
				toggled = true
			}
		} else if pressed && guiContext.DragDropHoldJustPressedId == id {
			IM_ASSERT(button_flags&ImGuiButtonFlags_PressedOnDragDropHold != 0)
			if !is_open { // When using Drag and Drop "hold to open" we keep the node highlighted after opening, but never close it again.
				toggled = true
			}
		}

		if guiContext.NavId == id && guiContext.NavMoveDir == ImGuiDir_Left && is_open {
			toggled = true
			NavMoveRequestCancel()
		}
		if guiContext.NavId == id && guiContext.NavMoveDir == ImGuiDir_Right && !is_open { // If there's something upcoming on the line we may want to give it the priority?
			toggled = true
			NavMoveRequestCancel()
		}

		if toggled {
			is_open = !is_open
			window.DC.StateStorage.SetInt(id, bool2int(is_open))
			guiContext.LastItemData.StatusFlags |= ImGuiItemStatusFlags_ToggledOpen
		}
	}
	if flags&ImGuiTreeNodeFlags_AllowItemOverlap != 0 {
		SetItemAllowOverlap()
	}

	// In this branch, TreeNodeBehavior() cannot toggle the selection so this will never trigger.
	if selected != was_selected { //-V547
		guiContext.LastItemData.StatusFlags |= ImGuiItemStatusFlags_ToggledSelection
	}

	// Render
	var text_col = GetColorU32FromID(ImGuiCol_Text, 1)
	var nav_highlight_flags = ImGuiNavHighlightFlags_TypeThin
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
			RenderBullet(window.DrawList, ImVec2{text_pos.x - text_offset_x*0.60, text_pos.y + guiContext.FontSize*0.5}, text_col)
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
			frame_bb.Max.x -= guiContext.FontSize + style.FramePadding.x
		}

		if guiContext.LogEnabled {
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
			RenderBullet(window.DrawList, ImVec2{text_pos.x - text_offset_x*0.5, text_pos.y + guiContext.FontSize*0.5}, text_col)
		} else if !is_leaf {
			var arrow ImGuiDir
			if is_open {
				arrow = ImGuiDir_Down
			} else {
				arrow = ImGuiDir_Right
			}
			RenderArrow(window.DrawList, ImVec2{text_pos.x - text_offset_x + padding.x, text_pos.y + guiContext.FontSize*0.15}, text_col, arrow, 0.70)
		}
		if guiContext.LogEnabled {
			LogSetNextTextDecoration(">", "")
		}
		RenderText(text_pos, label, false)
	}

	if is_open && flags&ImGuiTreeNodeFlags_NoTreePushOnOpen == 0 {
		TreePushOverrideID(id)
	}
	return is_open
}
