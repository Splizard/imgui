package imgui

// button
func Button(label string) bool {
	return ButtonEx(label, &ImVec2{}, ImGuiButtonFlags_None)
}

func ButtonEx(label string, size_arg *ImVec2, flags ImGuiButtonFlags) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style
	var id = window.GetIDs(label)
	var label_size = CalcTextSize(label, true, 0)

	var pos ImVec2 = window.DC.CursorPos
	if (flags&ImGuiButtonFlags_AlignTextBaseLine) != 0 && style.FramePadding.y < window.DC.CurrLineTextBaseOffset { // Try to vertically align buttons that are smaller/have no padding so that text baseline matches (bit hacky, since it shouldn't be a flag)
		pos.y += window.DC.CurrLineTextBaseOffset - style.FramePadding.y
	}
	var size ImVec2 = CalcItemSize(*size_arg, label_size.x+style.FramePadding.x*2.0, label_size.y+style.FramePadding.y*2.0)

	var bb = ImRect{pos, pos.Add(size)}
	ItemSizeVec(&size, style.FramePadding.y)
	if !ItemAdd(&bb, id, nil, 0) {
		return false
	}

	if g.LastItemData.InFlags&ImGuiItemFlags_ButtonRepeat != 0 {
		flags |= ImGuiButtonFlags_Repeat
	}

	var hovered, held bool
	var pressed bool = ButtonBehavior(&bb, id, &hovered, &held, flags)

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

	if g.LogEnabled {
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
//------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnClickRelease:             return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//   Frame N+0 (mouse is outside bb)        -             -                -               -                  -                    -
//   Frame N+1 (mouse moves inside bb)      -             true             -               -                  -                    -
//   Frame N+2 (mouse button is down)       -             true             true            true               -                    true
//   Frame N+3 (mouse button is down)       -             true             true            -                  -                    -
//   Frame N+4 (mouse moves outside bb)     -             -                true            -                  -                    -
//   Frame N+5 (mouse moves inside bb)      -             true             true            -                  -                    -
//   Frame N+6 (mouse button is released)   true          true             -               -                  true                 -
//   Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//   Frame N+8 (mouse moves outside bb)     -             -                -               -                  -                    -
//------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnClick:                    return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//   Frame N+2 (mouse button is down)       true          true             true            true               -                    true
//   Frame N+3 (mouse button is down)       -             true             true            -                  -                    -
//   Frame N+6 (mouse button is released)   -             true             -               -                  true                 -
//   Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnRelease:                  return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//   Frame N+2 (mouse button is down)       -             true             -               -                  -                    true
//   Frame N+3 (mouse button is down)       -             true             -               -                  -                    -
//   Frame N+6 (mouse button is released)   true          true             -               -                  -                    -
//   Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//------------------------------------------------------------------------------------------------------------------------------------------------
// with PressedOnDoubleClick:              return-value  IsItemHovered()  IsItemActive()  IsItemActivated()  IsItemDeactivated()  IsItemClicked()
//   Frame N+0 (mouse button is down)       -             true             -               -                  -                    true
//   Frame N+1 (mouse button is down)       -             true             -               -                  -                    -
//   Frame N+2 (mouse button is released)   -             true             -               -                  -                    -
//   Frame N+3 (mouse button is released)   -             true             -               -                  -                    -
//   Frame N+4 (mouse button is down)       true          true             true            true               -                    true
//   Frame N+5 (mouse button is down)       -             true             true            -                  -                    -
//   Frame N+6 (mouse button is released)   -             true             -               -                  true                 -
//   Frame N+7 (mouse button is released)   -             true             -               -                  -                    -
//------------------------------------------------------------------------------------------------------------------------------------------------
// Note that some combinations are supported,
// - PressedOnDragDropHold can generally be associated with any flag.
// - PressedOnDoubleClick can be associated by PressedOnClickRelease/PressedOnRelease, in which case the second release event won't be reported.
//------------------------------------------------------------------------------------------------------------------------------------------------
// The behavior of the return-value changes when ImGuiButtonFlags_Repeat is set:
//                                         Repeat+                  Repeat+           Repeat+             Repeat+
//                                         PressedOnClickRelease    PressedOnClick    PressedOnRelease    PressedOnDoubleClick
//-------------------------------------------------------------------------------------------------------------------------------------------------
//   Frame N+0 (mouse button is down)       -                        true              -                   true
//   ...                                    -                        -                 -                   -
//   Frame N + RepeatDelay                  true                     true              -                   true
//   ...                                    -                        -                 -                   -
//   Frame N + RepeatDelay + RepeatRate*N   true                     true              -                   true
//-------------------------------------------------------------------------------------------------------------------------------------------------
func ButtonBehavior(bb *ImRect, id ImGuiID, out_hovered *bool, out_held *bool, flags ImGuiButtonFlags) bool {
	var g = GImGui
	var window = GetCurrentWindow()

	// Default only reacts to left mouse button
	if (flags & ImGuiButtonFlags_MouseButtonMask_) == 0 {
		flags |= ImGuiButtonFlags_MouseButtonDefault_
	}

	// Default behavior requires click + release inside bounding box
	if (flags & ImGuiButtonFlags_PressedOnMask_) == 0 {
		flags |= ImGuiButtonFlags_PressedOnDefault_
	}

	var backup_hovered_window *ImGuiWindow = g.HoveredWindow
	var flatten_hovered_children bool = (flags&ImGuiButtonFlags_FlattenChildren != 0) && g.HoveredWindow != nil && g.HoveredWindow.RootWindow == window
	if flatten_hovered_children {
		g.HoveredWindow = window
	}

	var pressed bool = false
	var hovered bool = ItemHoverable(bb, id)

	// Drag source doesn't report as hovered
	if hovered && g.DragDropActive && g.DragDropPayload.SourceId == id && 0 == (g.DragDropSourceFlags&ImGuiDragDropFlags_SourceNoDisableHover) {
		hovered = false
	}

	// Special mode for Drag and Drop where holding button pressed for a long time while dragging another item triggers the button
	if g.DragDropActive && (flags&ImGuiButtonFlags_PressedOnDragDropHold != 0) && 0 == (g.DragDropSourceFlags&ImGuiDragDropFlags_SourceNoHoldToOpenOthers) {
		if IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByActiveItem) {
			hovered = true
			SetHoveredID(id)
			if g.HoveredIdTimer-g.IO.DeltaTime <= DRAGDROP_HOLD_TO_OPEN_TIMER && g.HoveredIdTimer >= DRAGDROP_HOLD_TO_OPEN_TIMER {
				pressed = true
				g.DragDropHoldJustPressedId = id
				FocusWindow(window)
			}
		}
	}

	if flatten_hovered_children {
		g.HoveredWindow = backup_hovered_window
	}

	// AllowOverlap mode (rarely used) requires previous frame HoveredId to be null or to match. This allows using patterns where a later submitted widget overlaps a previous one.
	if hovered && (flags&ImGuiButtonFlags_AllowItemOverlap != 0) && (g.HoveredIdPreviousFrame != id && g.HoveredIdPreviousFrame != 0) {
		hovered = false
	}

	// Mouse handling
	if hovered {
		if 0 == (flags&ImGuiButtonFlags_NoKeyModifiers) || (!g.IO.KeyCtrl && !g.IO.KeyShift && !g.IO.KeyAlt) {
			// Poll buttons
			var mouse_button_clicked ImGuiMouseButton = -1
			var mouse_button_released ImGuiMouseButton = -1
			if (flags&ImGuiButtonFlags_MouseButtonLeft != 0) && g.IO.MouseClicked[0] {
				mouse_button_clicked = 0
			} else if (flags&ImGuiButtonFlags_MouseButtonRight != 0) && g.IO.MouseClicked[1] {
				mouse_button_clicked = 1
			} else if (flags&ImGuiButtonFlags_MouseButtonMiddle != 0) && g.IO.MouseClicked[2] {
				mouse_button_clicked = 2
			}
			if (flags&ImGuiButtonFlags_MouseButtonLeft != 0) && g.IO.MouseReleased[0] {
				mouse_button_released = 0
			} else if (flags&ImGuiButtonFlags_MouseButtonRight != 0) && g.IO.MouseReleased[1] {
				mouse_button_released = 1
			} else if (flags&ImGuiButtonFlags_MouseButtonMiddle != 0) && g.IO.MouseReleased[2] {
				mouse_button_released = 2
			}

			if mouse_button_clicked != -1 && g.ActiveId != id {
				if flags&(ImGuiButtonFlags_PressedOnClickRelease|ImGuiButtonFlags_PressedOnClickReleaseAnywhere) != 0 {
					SetActiveID(id, window)
					g.ActiveIdMouseButton = mouse_button_clicked
					if 0 == (flags & ImGuiButtonFlags_NoNavFocus) {
						SetFocusID(id, window)
					}
					FocusWindow(window)
				}
				if (flags&ImGuiButtonFlags_PressedOnClick != 0) || ((flags&ImGuiButtonFlags_PressedOnDoubleClick != 0) && g.IO.MouseDoubleClicked[mouse_button_clicked]) {
					pressed = true
					if flags&ImGuiButtonFlags_NoHoldingActiveId != 0 {
						ClearActiveID()
					} else {
						SetActiveID(id, window) // Hold on ID
					}
					g.ActiveIdMouseButton = mouse_button_clicked
					FocusWindow(window)
				}
			}
			if (flags&ImGuiButtonFlags_PressedOnRelease != 0) && mouse_button_released != -1 {
				// Repeat mode trumps on release behavior
				var has_repeated_at_least_once bool = (flags&ImGuiButtonFlags_Repeat != 0) && g.IO.MouseDownDurationPrev[mouse_button_released] >= g.IO.KeyRepeatDelay
				if !has_repeated_at_least_once {
					pressed = true
				}
				ClearActiveID()
			}

			// 'Repeat' mode acts when held regardless of _PressedOn flags (see table above).
			// Relies on repeat logic of IsMouseClicked() but we may as well do it ourselves if we end up exposing finer RepeatDelay/RepeatRate settings.
			if g.ActiveId == id && (flags&ImGuiButtonFlags_Repeat != 0) {
				if g.IO.MouseDownDuration[g.ActiveIdMouseButton] > 0.0 && IsMouseClicked(g.ActiveIdMouseButton, true) {
					pressed = true
				}
			}
		}

		if pressed {
			g.NavDisableHighlight = true
		}
	}

	// Gamepad/Keyboard navigation
	// We report navigated item as hovered but we don't set g.HoveredId to not interfere with mouse.
	if g.NavId == id && !g.NavDisableHighlight && g.NavDisableMouseHover && (g.ActiveId == 0 || g.ActiveId == id || g.ActiveId == window.MoveId) {
		if 0 == (flags & ImGuiButtonFlags_NoHoveredOnFocus) {
			hovered = true
		}
	}
	if g.NavActivateDownId == id {
		var nav_activated_by_code bool = (g.NavActivateId == id)
		var nav_activated_by_inputs bool

		if (flags & ImGuiButtonFlags_Repeat) != 0 {
			nav_activated_by_inputs = IsNavInputTest(ImGuiNavInput_Activate, ImGuiInputReadMode_Repeat)
		} else {
			nav_activated_by_inputs = IsNavInputTest(ImGuiNavInput_Activate, ImGuiInputReadMode_Pressed)
		}

		if nav_activated_by_code || nav_activated_by_inputs {
			pressed = true
		}
		if nav_activated_by_code || nav_activated_by_inputs || g.ActiveId == id {
			// Set active id so it can be queried by user via IsItemActive(), equivalent of holding the mouse button.
			g.NavActivateId = id // This is so SetActiveId assign a Nav source
			SetActiveID(id, window)
			if (nav_activated_by_code || nav_activated_by_inputs) && 0 == (flags&ImGuiButtonFlags_NoNavFocus) {
				SetFocusID(id, window)
			}
		}
	}

	// Process while held
	var held bool = false
	if g.ActiveId == id {
		if g.ActiveIdSource == ImGuiInputSource_Mouse {
			if g.ActiveIdIsJustActivated {
				g.ActiveIdClickOffset = g.IO.MousePos.Sub(bb.Min)
			}

			var mouse_button = g.ActiveIdMouseButton
			IM_ASSERT(mouse_button >= 0 && mouse_button < ImGuiMouseButton_COUNT)
			if g.IO.MouseDown[mouse_button] {
				held = true
			} else {
				var release_in bool = hovered && (flags&ImGuiButtonFlags_PressedOnClickRelease) != 0
				var release_anywhere bool = (flags & ImGuiButtonFlags_PressedOnClickReleaseAnywhere) != 0
				if (release_in || release_anywhere) && !g.DragDropActive {
					// Report as pressed when releasing the mouse (this is the most common path)
					var is_double_click_release bool = (flags&ImGuiButtonFlags_PressedOnDoubleClick != 0) && g.IO.MouseDownWasDoubleClick[mouse_button]
					var is_repeating_already bool = (flags&ImGuiButtonFlags_Repeat != 0) && g.IO.MouseDownDurationPrev[mouse_button] >= g.IO.KeyRepeatDelay // Repeat mode trumps <on release>
					if !is_double_click_release && !is_repeating_already {
						pressed = true
					}
				}
				ClearActiveID()
			}
			if 0 == (flags & ImGuiButtonFlags_NoNavFocus) {
				g.NavDisableHighlight = true
			}
		} else if g.ActiveIdSource == ImGuiInputSource_Nav {
			// When activated using Nav, we hold on the ActiveID until activation button is released
			if g.NavActivateDownId != id {
				ClearActiveID()
			}
		}
		if pressed {
			g.ActiveIdHasBeenPressedBefore = true
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
	var g = GImGui
	var window = g.CurrentWindow

	var bb = ImRect{*pos, pos.Add(ImVec2{g.FontSize, g.FontSize}).Add(g.Style.FramePadding.Scale(2.0))}
	ItemAdd(&bb, id, nil, 0)
	var hovered, held bool
	var pressed bool = ButtonBehavior(&bb, id, &hovered, &held, ImGuiButtonFlags_None)

	// Render
	var bg_col ImU32 = GetColorU32FromID(ImGuiCol_Button, 1)
	if held && hovered {
		bg_col = GetColorU32FromID(ImGuiCol_ButtonActive, 1)
	} else if hovered {
		bg_col = GetColorU32FromID(ImGuiCol_ButtonHovered, 1)
	}

	var text_col ImU32 = GetColorU32FromID(ImGuiCol_Text, 1)
	var center ImVec2 = bb.GetCenter()
	if hovered || held {
		window.DrawList.AddCircleFilled(center /*+ ImVec2(0.0f, -0.5f)*/, g.FontSize*0.5+1.0, bg_col, 12)
	}

	var arrow = ImGuiDir_Down
	if window.Collapsed {
		arrow = ImGuiDir_Right
	}

	RenderArrow(window.DrawList, bb.Min.Add(g.Style.FramePadding), text_col, arrow, 1.0)

	// Switch to moving the window after mouse is moved beyond the initial drag threshold
	if IsItemActive() && IsMouseDragging(0, -1) {
		StartMouseMovingWindow(window)
	}

	return pressed
}
