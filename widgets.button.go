package imgui

// in 'repeat' mode, Button*() functions return repeated true in a typematic manner (using io.KeyRepeatDelay/io.KeyRepeatRate setting). Note that you can call IsItemActive() after any Button() to tell if the button is held in the current frame.
func PushButtonRepeat(repeat bool) {
	PushItemFlag(ImGuiItemFlags_ButtonRepeat, repeat)
}

func PopButtonRepeat() { PopItemFlag() }

// button
func Button(label string) bool {
	return ButtonEx(label, &ImVec2{}, ImGuiButtonFlags_None)
}

// Small buttons fits within text without additional vertical spacing.
// button with FramePadding=(0,0) to easily embed within text
func SmallButton(label string) bool {
	var backup_padding_y = guiContext.Style.FramePadding.y
	guiContext.Style.FramePadding.y = 0.0
	var pressed = ButtonEx(label, &ImVec2{}, ImGuiButtonFlags_AlignTextBaseLine)
	guiContext.Style.FramePadding.y = backup_padding_y
	return pressed
}

func ButtonEx(label string, size_arg *ImVec2, flags ImGuiButtonFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := guiContext.Style
	var id = window.GetIDs(label)
	var label_size = CalcTextSize(label, true, 0)

	var pos = window.DC.CursorPos
	if (flags&ImGuiButtonFlags_AlignTextBaseLine) != 0 && style.FramePadding.y < window.DC.CurrLineTextBaseOffset { // Try to vertically align buttons that are smaller/have no padding so that text baseline matches (bit hacky, since it shouldn't be a flag)
		pos.y += window.DC.CurrLineTextBaseOffset - style.FramePadding.y
	}
	var size = CalcItemSize(*size_arg, label_size.x+style.FramePadding.x*2.0, label_size.y+style.FramePadding.y*2.0)

	var bb = ImRect{pos, pos.Add(size)}
	ItemSizeVec(&size, style.FramePadding.y)
	if !ItemAdd(&bb, id, nil, 0) {
		return false
	}

	if guiContext.LastItemData.InFlags&ImGuiItemFlags_ButtonRepeat != 0 {
		flags |= ImGuiButtonFlags_Repeat
	}

	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, flags)

	// Render
	var col ImU32
	if held && hovered {
		col = GetColorU32FromID(ImGuiCol_ButtonActive, 1)
	} else if hovered {
		col = GetColorU32FromID(ImGuiCol_ButtonHovered, 1)
	} else {
		col = GetColorU32FromID(ImGuiCol_Button, 1)
	}
	RenderNavHighlight(&bb, id, 0)
	RenderFrame(bb.Min, bb.Max, col, true, style.FrameRounding)

	if guiContext.LogEnabled {
		LogSetNextTextDecoration("[", "]")
	}
	min, max := bb.Min.Add(style.FramePadding), bb.Max.Sub(style.FramePadding)
	RenderTextClipped(&min, &max, label, &label_size, &style.ButtonTextAlign, &bb)

	return pressed
}

// The ButtonBehavior() function is key to many interactions and used by many/most widgets.
// Because we handle so many cases (keyboard/gamepad navigation, drag and drop) and many specific behavior (via ImGuiButtonFlags_),
// this code is a little complex.
// By far the most common path is interacting with the Mouse using the default ImGuiButtonFlags_PressedOnClickRelease button behavior.
// See the series of events below and the corresponding state reported by dear imgui:
// ------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnClickRelease:             return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//
//	Frame N+0 (mouse is outside bb)        -             -                -               -                  -                    -
//	Frame N+1 (mouse moves inside bb)      -             true             -               -                  -                    -
//	Frame N+2 (mouse button is down)       -             true             true            true               -                    true
//	Frame N+3 (mouse button is down)       -             true             true            -                  -                    -
//	Frame N+4 (mouse moves outside bb)     -             -                true            -                  -                    -
//	Frame N+5 (mouse moves inside bb)      -             true             true            -                  -                    -
//	Frame N+6 (mouse button is released)   true          true             -               -                  true                 -
//	Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//	Frame N+8 (mouse moves outside bb)     -             -                -               -                  -                    -
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnClick:                    return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//
//	Frame N+2 (mouse button is down)       true          true             true            true               -                    true
//	Frame N+3 (mouse button is down)       -             true             true            -                  -                    -
//	Frame N+6 (mouse button is released)   -             true             -               -                  true                 -
//	Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnRelease:                  return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//
//	Frame N+2 (mouse button is down)       -             true             -               -                  -                    true
//	Frame N+3 (mouse button is down)       -             true             -               -                  -                    -
//	Frame N+6 (mouse button is released)   true          true             -               -                  -                    -
//	Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnDoubleClick:              return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//
//	Frame N+0 (mouse button is down)       -             true             -               -                  -                    true
//	Frame N+1 (mouse button is down)       -             true             -               -                  -                    -
//	Frame N+2 (mouse button is released)   -             true             -               -                  -                    -
//	Frame N+3 (mouse button is released)   -             true             -               -                  -                    -
//	Frame N+4 (mouse button is down)       true          true             true            true               -                    true
//	Frame N+5 (mouse button is down)       -             true             true            -                  -                    -
//	Frame N+6 (mouse button is released)   -             true             -               -                  true                 -
//	Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
// Note that some combinations are supported,
// - PressedOnDragDropHold can generally be associated with any flag.
// - PressedOnDoubleClick can be associated by PressedOnClickRelease/PressedOnRelease, in which case the second release event won't be reported.
// ------------------------------------------------------------------------------------------------------------------------------------------------
// The behavior of the return-value changes when ImGuiButtonFlags_Repeat is set:
//
//	Repeat+                  Repeat+           Repeat+             Repeat+
//	PressedOnClickRelease    PressedOnClick    PressedOnRelease    PressedOnDoubleClick
//
// -------------------------------------------------------------------------------------------------------------------------------------------------
//
//	Frame N+0 (mouse button is down)       -                        true              -                   true
//	...                                    -                        -                 -                   -
//	Frame N + RepeatDelay                  true                     true              -                   true
//	...                                    -                        -                 -                   -
//	Frame N + RepeatDelay + RepeatRate*N   true                     true              -                   true
//
// -------------------------------------------------------------------------------------------------------------------------------------------------
func ButtonBehavior(bb *ImRect, id ImGuiID, out_hovered *bool, out_held *bool, flags ImGuiButtonFlags) bool {
	window := GetCurrentWindow()

	// Default only reacts to left mouse button
	if (flags & ImGuiButtonFlags_MouseButtonMask_) == 0 {
		flags |= ImGuiButtonFlags_MouseButtonDefault_
	}

	// Default behavior requires click + release inside bounding box
	if (flags & ImGuiButtonFlags_PressedOnMask_) == 0 {
		flags |= ImGuiButtonFlags_PressedOnDefault_
	}

	var backup_hovered_window = guiContext.HoveredWindow
	var flatten_hovered_children = (flags&ImGuiButtonFlags_FlattenChildren != 0) && guiContext.HoveredWindow != nil && guiContext.HoveredWindow.RootWindow == window
	if flatten_hovered_children {
		guiContext.HoveredWindow = window
	}

	var pressed = false
	var hovered = ItemHoverable(bb, id)

	// Drag source doesn't report as hovered
	if hovered && guiContext.DragDropActive && guiContext.DragDropPayload.SourceId == id && guiContext.DragDropSourceFlags&ImGuiDragDropFlags_SourceNoDisableHover == 0 {
		hovered = false
	}

	// Special mode for Drag and Drop where holding button pressed for a long time while dragging another item triggers the button
	if guiContext.DragDropActive && (flags&ImGuiButtonFlags_PressedOnDragDropHold != 0) && guiContext.DragDropSourceFlags&ImGuiDragDropFlags_SourceNoHoldToOpenOthers == 0 {
		if IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByActiveItem) {
			hovered = true
			SetHoveredID(id)
			if guiContext.HoveredIdTimer-guiContext.IO.DeltaTime <= DRAGDROP_HOLD_TO_OPEN_TIMER && guiContext.HoveredIdTimer >= DRAGDROP_HOLD_TO_OPEN_TIMER {
				pressed = true
				guiContext.DragDropHoldJustPressedId = id
				FocusWindow(window)
			}
		}
	}

	if flatten_hovered_children {
		guiContext.HoveredWindow = backup_hovered_window
	}

	// AllowOverlap mode (rarely used) requires previous frame HoveredId to be null or to match. This allows using patterns where a later submitted widget overlaps a previous one.
	if hovered && (flags&ImGuiButtonFlags_AllowItemOverlap != 0) && (guiContext.HoveredIdPreviousFrame != id && guiContext.HoveredIdPreviousFrame != 0) {
		hovered = false
	}

	// Mouse handling
	if hovered {
		if flags&ImGuiButtonFlags_NoKeyModifiers == 0 || (!guiContext.IO.KeyCtrl && !guiContext.IO.KeyShift && !guiContext.IO.KeyAlt) {
			// Poll buttons
			var mouse_button_clicked ImGuiMouseButton = -1
			var mouse_button_released ImGuiMouseButton = -1
			if (flags&ImGuiButtonFlags_MouseButtonLeft != 0) && guiContext.IO.MouseClicked[0] {
				mouse_button_clicked = 0
			} else if (flags&ImGuiButtonFlags_MouseButtonRight != 0) && guiContext.IO.MouseClicked[1] {
				mouse_button_clicked = 1
			} else if (flags&ImGuiButtonFlags_MouseButtonMiddle != 0) && guiContext.IO.MouseClicked[2] {
				mouse_button_clicked = 2
			}
			if (flags&ImGuiButtonFlags_MouseButtonLeft != 0) && guiContext.IO.MouseReleased[0] {
				mouse_button_released = 0
			} else if (flags&ImGuiButtonFlags_MouseButtonRight != 0) && guiContext.IO.MouseReleased[1] {
				mouse_button_released = 1
			} else if (flags&ImGuiButtonFlags_MouseButtonMiddle != 0) && guiContext.IO.MouseReleased[2] {
				mouse_button_released = 2
			}

			if mouse_button_clicked != -1 && guiContext.ActiveId != id {
				if flags&(ImGuiButtonFlags_PressedOnClickRelease|ImGuiButtonFlags_PressedOnClickReleaseAnywhere) != 0 {
					SetActiveID(id, window)
					guiContext.ActiveIdMouseButton = mouse_button_clicked
					if flags&ImGuiButtonFlags_NoNavFocus == 0 {
						SetFocusID(id, window)
					}
					FocusWindow(window)
				}
				if (flags&ImGuiButtonFlags_PressedOnClick != 0) || ((flags&ImGuiButtonFlags_PressedOnDoubleClick != 0) && guiContext.IO.MouseDoubleClicked[mouse_button_clicked]) {
					pressed = true
					if flags&ImGuiButtonFlags_NoHoldingActiveId != 0 {
						ClearActiveID()
					} else {
						SetActiveID(id, window) // Hold on ID
					}
					guiContext.ActiveIdMouseButton = mouse_button_clicked
					FocusWindow(window)
				}
			}
			if (flags&ImGuiButtonFlags_PressedOnRelease != 0) && mouse_button_released != -1 {
				// Repeat mode trumps on release behavior
				var has_repeated_at_least_once = (flags&ImGuiButtonFlags_Repeat != 0) && guiContext.IO.MouseDownDurationPrev[mouse_button_released] >= guiContext.IO.KeyRepeatDelay
				if !has_repeated_at_least_once {
					pressed = true
				}
				ClearActiveID()
			}

			// 'Repeat' mode acts when held regardless of _PressedOn flags (see table above).
			// Relies on repeat logic of IsMouseClicked() but we may as well do it ourselves if we end up exposing finer RepeatDelay/RepeatRate settings.
			if guiContext.ActiveId == id && (flags&ImGuiButtonFlags_Repeat != 0) {
				if guiContext.IO.MouseDownDuration[guiContext.ActiveIdMouseButton] > 0.0 && IsMouseClicked(guiContext.ActiveIdMouseButton, true) {
					pressed = true
				}
			}
		}

		if pressed {
			guiContext.NavDisableHighlight = true
		}
	}

	// Gamepad/Keyboard navigation
	// We report navigated item as hovered but we don't set guiContext.HoveredId to not interfere with mouse.
	if guiContext.NavId == id && !guiContext.NavDisableHighlight && guiContext.NavDisableMouseHover && (guiContext.ActiveId == 0 || guiContext.ActiveId == id || guiContext.ActiveId == window.MoveId) {
		if flags&ImGuiButtonFlags_NoHoveredOnFocus == 0 {
			hovered = true
		}
	}
	if guiContext.NavActivateDownId == id {
		var nav_activated_by_code = (guiContext.NavActivateId == id)
		var nav_activated_by_inputs bool

		if (flags & ImGuiButtonFlags_Repeat) != 0 {
			nav_activated_by_inputs = IsNavInputTest(ImGuiNavInput_Activate, ImGuiInputReadMode_Repeat)
		} else {
			nav_activated_by_inputs = IsNavInputTest(ImGuiNavInput_Activate, ImGuiInputReadMode_Pressed)
		}

		if nav_activated_by_code || nav_activated_by_inputs {
			pressed = true
		}
		if nav_activated_by_code || nav_activated_by_inputs || guiContext.ActiveId == id {
			// Set active id so it can be queried by user via IsItemActive(), equivalent of holding the mouse button.
			guiContext.NavActivateId = id // This is so SetActiveId assign a Nav source
			SetActiveID(id, window)
			if (nav_activated_by_code || nav_activated_by_inputs) && flags&ImGuiButtonFlags_NoNavFocus == 0 {
				SetFocusID(id, window)
			}
		}
	}

	// Process while held
	var held = false
	if guiContext.ActiveId == id {
		if guiContext.ActiveIdSource == ImGuiInputSource_Mouse {
			if guiContext.ActiveIdIsJustActivated {
				guiContext.ActiveIdClickOffset = guiContext.IO.MousePos.Sub(bb.Min)
			}

			var mouse_button = guiContext.ActiveIdMouseButton
			IM_ASSERT(mouse_button >= 0 && mouse_button < ImGuiMouseButton_COUNT)
			if guiContext.IO.MouseDown[mouse_button] {
				held = true
			} else {
				var release_in = hovered && (flags&ImGuiButtonFlags_PressedOnClickRelease) != 0
				var release_anywhere = (flags & ImGuiButtonFlags_PressedOnClickReleaseAnywhere) != 0
				if (release_in || release_anywhere) && !guiContext.DragDropActive {
					// Report as pressed when releasing the mouse (this is the most common path)
					var is_double_click_release = (flags&ImGuiButtonFlags_PressedOnDoubleClick != 0) && guiContext.IO.MouseDownWasDoubleClick[mouse_button]
					var is_repeating_already = (flags&ImGuiButtonFlags_Repeat != 0) && guiContext.IO.MouseDownDurationPrev[mouse_button] >= guiContext.IO.KeyRepeatDelay // Repeat mode trumps <on release>
					if !is_double_click_release && !is_repeating_already {
						pressed = true
					}
				}
				ClearActiveID()
			}
			if flags&ImGuiButtonFlags_NoNavFocus == 0 {
				guiContext.NavDisableHighlight = true
			}
		} else if guiContext.ActiveIdSource == ImGuiInputSource_Nav {
			// When activated using Nav, we hold on the ActiveID until activation button is released
			if guiContext.NavActivateDownId != id {
				ClearActiveID()
			}
		}
		if pressed {
			guiContext.ActiveIdHasBeenPressedBefore = true
		}
	}

	if out_hovered != nil {
		*out_hovered = hovered
	}
	if out_held != nil {
		*out_held = held
	}

	return pressed
}

func CollapseButton(id ImGuiID, pos *ImVec2) bool {
	window := guiContext.CurrentWindow

	var bb = ImRect{*pos, pos.Add(ImVec2{guiContext.FontSize, guiContext.FontSize}).Add(guiContext.Style.FramePadding.Scale(2.0))}
	ItemAdd(&bb, id, nil, 0)
	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, ImGuiButtonFlags_None)

	// Render
	var bg_col = GetColorU32FromID(ImGuiCol_Button, 1)
	if held && hovered {
		bg_col = GetColorU32FromID(ImGuiCol_ButtonActive, 1)
	} else if hovered {
		bg_col = GetColorU32FromID(ImGuiCol_ButtonHovered, 1)
	}

	var text_col = GetColorU32FromID(ImGuiCol_Text, 1)
	var center = bb.GetCenter()
	if hovered || held {
		window.DrawList.AddCircleFilled(center /*+ ImVec2(0.0f, -0.5f)*/, guiContext.FontSize*0.5+1.0, bg_col, 12)
	}

	var arrow = ImGuiDir_Down
	if window.Collapsed {
		arrow = ImGuiDir_Right
	}

	RenderArrow(window.DrawList, bb.Min.Add(guiContext.Style.FramePadding), text_col, arrow, 1.0)

	// Switch to moving the window after mouse is moved beyond the initial drag threshold
	if IsItemActive() && IsMouseDragging(0, -1) {
		StartMouseMovingWindow(window)
	}

	return pressed
}

// flexible button behavior without the visuals, frequently useful to build custom behaviors using the public api (along with IsItemActive, IsItemHovered, etc.)
// Tip: use ImGui::PushID()/PopID() to push indices or pointers in the ID stack.
// Then you can keep 'str_id' empty or the same for all your buttons (instead of creating a string based on a non-string id)
func InvisibleButton(str_id string, size_arg ImVec2, flags ImGuiButtonFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	// Cannot use zero-size for InvisibleButton(). Unlike Button() there is not way to fallback using the label size.
	IM_ASSERT(size_arg.x != 0.0 && size_arg.y != 0.0)

	var id = window.GetIDs(str_id)
	var size = CalcItemSize(size_arg, 0.0, 0.0)
	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	ItemSizeVec(&size, 0)
	if !ItemAdd(&bb, id, nil, 0) {
		return false
	}

	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, flags)

	return pressed
}

// square button with an arrow shape
func ArrowButton(str_id string, dir ImGuiDir) bool {
	var sz = GetFrameHeight()
	return ArrowButtonEx(str_id, dir, ImVec2{sz, sz}, ImGuiButtonFlags_None)
}

// use with e.guiContext. if (RadioButton("one", my_value==1)) { my_value = 1 bool {panic("not implemented")} }
func RadioButtonBool(label string, active bool) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := guiContext.Style
	var id = window.GetIDs(label)
	var label_size = CalcTextSize(label, true, -1)

	var square_sz = GetFrameHeight()
	var pos = window.DC.CursorPos
	var check_bb = ImRect{pos, pos.Add(ImVec2{square_sz, square_sz})}

	var padding float
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var total_bb = ImRect{pos, pos.Add(ImVec2{square_sz + padding, label_size.y + style.FramePadding.y*2.0})}
	ItemSizeRect(&total_bb, style.FramePadding.y)
	if !ItemAdd(&total_bb, id, nil, 0) {
		return false
	}

	var center = check_bb.GetCenter()
	center.x = IM_ROUND(center.x)
	center.y = IM_ROUND(center.y)
	var radius = (square_sz - 1.0) * 0.5

	var hovered, held bool
	var pressed = ButtonBehavior(&total_bb, id, &hovered, &held, 0)
	if pressed {
		MarkItemEdited(id)
	}

	RenderNavHighlight(&total_bb, id, 0)

	var c = ImGuiCol_FrameBg
	if held && hovered {
		c = ImGuiCol_FrameBgHovered
	} else if hovered {
		c = ImGuiCol_FrameBgHovered
	}

	window.DrawList.AddCircleFilled(center, radius, GetColorU32FromID(c, 1), 16)
	if active {
		var pad = max(1.0, IM_FLOOR(square_sz/6.0))
		window.DrawList.AddCircleFilled(center, radius-pad, GetColorU32FromID(ImGuiCol_CheckMark, 1), 16)
	}

	if style.FrameBorderSize > 0.0 {
		window.DrawList.AddCircle(center.Add(ImVec2{1, 1}), radius, GetColorU32FromID(ImGuiCol_BorderShadow, 1), 16, style.FrameBorderSize)
		window.DrawList.AddCircle(center, radius, GetColorU32FromID(ImGuiCol_Border, 1), 16, style.FrameBorderSize)
	}

	var label_pos = ImVec2{check_bb.Max.x + style.ItemInnerSpacing.x, check_bb.Min.y + style.FramePadding.y}
	if guiContext.LogEnabled {
		s := "( )"
		if active {
			s = "(X)"
		}
		LogRenderedText(&label_pos, s)
	}
	if label_size.x > 0.0 {
		RenderText(label_pos, label, true)
	}

	return pressed
}

// shortcut to handle the above pattern when value is an integer
func RadioButtonInt(label string, v *int, v_button int) bool {
	// FIXME: This would work nicely if it was a public template, e.guiContext. 'template<T> RadioButton(const char* label, T* v, T v_button)', but I'm not sure how we would expose it..
	var pressed = RadioButtonBool(label, *v == v_button)
	if pressed {
		*v = v_button
	}
	return pressed
}

// Button to close a window
func CloseButton(id ImGuiID, pos *ImVec2) bool {
	window := guiContext.CurrentWindow

	// Tweak 1: Shrink hit-testing area if button covers an abnormally large proportion of the visible region. That's in order to facilitate moving the window away. (#3825)
	// This may better be applied as a general hit-rect reduction mechanism for all widgets to ensure the area to move window is always accessible?
	var bb = ImRect{*pos, pos.Add(ImVec2{guiContext.FontSize, guiContext.FontSize}).Add(guiContext.Style.FramePadding.Scale(2.0))}
	var bb_interact = bb
	var area_to_visible_ratio = window.OuterRectClipped.GetArea() / bb.GetArea()
	if area_to_visible_ratio < 1.5 {
		expansion := bb_interact.GetSize().Scale(-0.25)
		bb_interact.ExpandVec(*ImFloorVec(&expansion))
	}

	// Tweak 2: We intentionally allow interaction when clipped so that a mechanical Alt,Right,Activate sequence can always close a window.
	// (this isn't the regular behavior of buttons, but it doesn't affect the user much because navigation tends to keep items visible).
	var is_clipped = !ItemAdd(&bb_interact, id, nil, 0)

	var hovered, held bool
	var pressed = ButtonBehavior(&bb_interact, id, &hovered, &held, 0)
	if is_clipped {
		return pressed
	}

	// Render
	// FIXME: Clarify this mess
	var c = ImGuiCol_ButtonHovered
	if held {
		c = ImGuiCol_ButtonActive
	}

	var col = GetColorU32FromID(c, 1)
	var center = bb.GetCenter()
	if hovered {
		window.DrawList.AddCircleFilled(center, max(2.0, guiContext.FontSize*0.5+1.0), col, 12)
	}

	var cross_extent = guiContext.FontSize*0.5*0.7071 - 1.0
	var cross_col = GetColorU32FromID(ImGuiCol_Text, 1)
	center = center.Sub(ImVec2{0.5, 0.5})
	a, b := center.Add(ImVec2{+cross_extent, +cross_extent}), center.Add(ImVec2{-cross_extent, -cross_extent})
	window.DrawList.AddLine(&a, &b, cross_col, 1.0)
	a, b = center.Add(ImVec2{+cross_extent, -cross_extent}), center.Add(ImVec2{-cross_extent, +cross_extent})
	window.DrawList.AddLine(&a, &b, cross_col, 1.0)

	return pressed
}

func ArrowButtonEx(str_id string, dir ImGuiDir, size ImVec2, flags ImGuiButtonFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var id = window.GetIDs(str_id)
	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	var default_size = GetFrameHeight()

	var baseline float = -1
	if size.y >= default_size {
		baseline = guiContext.Style.FramePadding.y
	}

	ItemSizeVec(&size, baseline)
	if !ItemAdd(&bb, id, nil, 0) {
		return false
	}

	if (guiContext.LastItemData.InFlags & ImGuiItemFlags_ButtonRepeat) != 0 {
		flags |= ImGuiButtonFlags_Repeat
	}

	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, flags)

	var c = ImGuiCol_Button
	if held && hovered {
		c = ImGuiCol_ButtonActive
	} else if hovered {
		c = ImGuiCol_ButtonHovered
	}

	// Render
	var bg_col = GetColorU32FromID(c, 1.0)
	var text_col = GetColorU32FromID(ImGuiCol_Text, 1)
	RenderNavHighlight(&bb, id, 0)
	RenderFrame(bb.Min, bb.Max, bg_col, true, guiContext.Style.FrameRounding)
	RenderArrow(window.DrawList, bb.Min.Add(ImVec2{max(0.0, (size.x-guiContext.FontSize)*0.5), max(0.0, (size.y-guiContext.FontSize)*0.5)}), text_col, dir, 1)

	return pressed
}
