package imgui

const WINDOWS_MOUSE_WHEEL_SCROLL_LOCK_TIMER float = 2.00 // Lock scrolled window (so it doesn't pick child windows that are scrolling through) for a certain time, unless mouse moved.

// We typically use ImVec2(-FLT_MAX,-FLT_MAX) to denote an invalid mouse position.
func IsMousePosValid(mouse_pos *ImVec2) bool {
	// The assert is only to silence a false-positive in XCode Static Analysis.
	// Because GImGui is not dereferenced in every code path, the static analyzer assume that it may be nil (which it doesn't for other functions).
	IM_ASSERT(GImGui != nil)
	var MOUSE_INVALID float = -256000.0
	var p ImVec2
	if mouse_pos != nil {
		p = *mouse_pos
	} else {
		p = GImGui.IO.MousePos
	}
	return p.x >= MOUSE_INVALID && p.y >= MOUSE_INVALID
}

func StartLockWheelingWindow(window *ImGuiWindow) {
	var g = GImGui
	if g.WheelingWindow == window {
		return
	}
	g.WheelingWindow = window
	g.WheelingWindowRefMousePos = g.IO.MousePos
	g.WheelingWindowTimer = WINDOWS_MOUSE_WHEEL_SCROLL_LOCK_TIMER
}

// Test if mouse cursor is hovering given rectangle
// NB- Rectangle is clipped by our current clip setting
// NB- Expand the rectangle to be generous on imprecise inputs systems (g.Style.TouchExtraPadding)
// is mouse hovering given bounding rect (in screen space). clipped by current clipping settings, but disregarding of other consideration of focus/window ordering/popup-block.
func IsMouseHoveringRect(r_min, r_max ImVec2, clip bool /*= true*/) bool {
	var g = GImGui

	// Clip
	var rect_clipped = ImRect{r_min, r_max}
	if clip {
		rect_clipped.ClipWith(g.CurrentWindow.ClipRect)
	}

	// Expand for touch input
	var rect_for_touch = ImRect{rect_clipped.Min.Sub(g.Style.TouchExtraPadding), rect_clipped.Max.Add(g.Style.TouchExtraPadding)}
	if !rect_for_touch.ContainsVec(g.IO.MousePos) {
		return false
	}
	return true
}

func IsWindowContentHoverable(window *ImGuiWindow, flags ImGuiHoveredFlags) bool {
	// An active popup disable hovering on other windows (apart from its own children)
	// FIXME-OPT: This could be cached/stored within the window.
	var g = GImGui
	if g.NavWindow != nil {
		if focused_root_window := g.NavWindow.RootWindow; focused_root_window != nil {
			if focused_root_window.WasActive && focused_root_window != window.RootWindow {
				// For the purpose of those flags we differentiate "standard popup" from "modal popup"
				// NB: The order of those two tests is important because Modal windows are also Popups.
				if focused_root_window.Flags&ImGuiWindowFlags_Modal != 0 {
					return false
				}
				if (focused_root_window.Flags&ImGuiWindowFlags_Popup != 0) && 0 == (flags&ImGuiHoveredFlags_AllowWhenBlockedByPopup) {
					return false
				}
			}
		}
	}
	return true
}

func UpdateMouseMovingWindowEndFrame() {
	var g = GImGui
	if g.ActiveId != 0 || g.HoveredId != 0 {
		return
	}

	// Unless we just made a window/popup appear
	if g.NavWindow != nil && g.NavWindow.Appearing {
		return
	}

	// Click on empty space to focus window and start moving
	// (after we're done with all our widgets)
	if g.IO.MouseClicked[0] {
		// Handle the edge case of a popup being closed while clicking in its empty space.
		// If we try to focus it, FocusWindow() > ClosePopupsOverWindow() will accidentally close any parent popups because they are not linked together any more.
		var root_window *ImGuiWindow
		if g.HoveredWindow != nil {
			root_window = g.HoveredWindow.RootWindow
		}
		var is_closed_popup = root_window != nil && (root_window.Flags&ImGuiWindowFlags_Popup != 0) && !IsPopupOpenID(root_window.PopupId, ImGuiPopupFlags_AnyPopupLevel)

		if root_window != nil && !is_closed_popup {
			StartMouseMovingWindow(g.HoveredWindow) //-V595

			// Cancel moving if clicked outside of title bar
			if g.IO.ConfigWindowsMoveFromTitleBarOnly && 0 == (root_window.Flags&ImGuiWindowFlags_NoTitleBar) {
				rect := root_window.TitleBarRect()
				if !rect.ContainsVec(g.IO.MouseClickedPos[0]) {
					g.MovingWindow = nil
				}
			}

			// Cancel moving if clicked over an item which was disabled or inhibited by popups (note that we know HoveredId == 0 already)
			if g.HoveredIdDisabled {
				g.MovingWindow = nil
			}
		} else if root_window == nil && g.NavWindow != nil && GetTopMostPopupModal() == nil {
			// Clicking on void disable focus
			FocusWindow(nil)
		}
	}

	// With right mouse button we close popups without changing focus based on where the mouse is aimed
	// Instead, focus will be restored to the window under the bottom-most closed popup.
	// (The left mouse button path calls FocusWindow on the hovered window, which will lead NewFrame.ClosePopupsOverWindow to trigger)
	if g.IO.MouseClicked[1] {
		// Find the top-most window between HoveredWindow and the top-most Modal Window.
		// This is where we can trim the popup stack.
		var modal *ImGuiWindow = GetTopMostPopupModal()
		var hovered_window_above_modal bool = g.HoveredWindow != nil && IsWindowAbove(g.HoveredWindow, modal)

		win := modal
		if hovered_window_above_modal {
			win = g.HoveredWindow
		}

		ClosePopupsOverWindow(win, true)
	}
}

func UpdateMouseWheel() {
	var g = GImGui

	// Reset the locked window if we move the mouse or after the timer elapses
	if g.WheelingWindow != nil {
		g.WheelingWindowTimer -= g.IO.DeltaTime
		if IsMousePosValid(nil) && ImLengthSqrVec2(g.IO.MousePos.Sub(g.WheelingWindowRefMousePos)) > g.IO.MouseDragThreshold*g.IO.MouseDragThreshold {
			g.WheelingWindowTimer = 0.0
		}
		if g.WheelingWindowTimer <= 0.0 {
			g.WheelingWindow = nil
			g.WheelingWindowTimer = 0.0
		}
	}

	if g.IO.MouseWheel == 0.0 && g.IO.MouseWheelH == 0.0 {
		return
	}

	if (g.ActiveId != 0 && g.ActiveIdUsingMouseWheel) || (g.HoveredIdPreviousFrame != 0 && g.HoveredIdPreviousFrameUsingMouseWheel) {
		return
	}

	var window *ImGuiWindow
	if g.WheelingWindow != nil {
		window = g.WheelingWindow
	} else {
		window = g.HoveredWindow
	}
	if window == nil || window.Collapsed {
		return
	}

	// Zoom / Scale window
	// FIXME-OBSOLETE: This is an old feature, it still works but pretty much nobody is using it and may be best redesigned.
	if g.IO.MouseWheel != 0.0 && g.IO.KeyCtrl && g.IO.FontAllowUserScaling {
		StartLockWheelingWindow(window)
		var new_font_scale float = ImClamp(window.FontWindowScale+g.IO.MouseWheel*0.10, 0.50, 2.50)
		var scale float = new_font_scale / window.FontWindowScale
		window.FontWindowScale = new_font_scale
		if window == window.RootWindow {
			var offset ImVec2 = window.Size.Scale((1.0 - scale)).Mul(g.IO.MousePos.Sub(window.Pos)).Div(window.Size)
			p := window.Pos.Add(offset)
			setWindowPos(window, &p, 0)
			scaled := window.Size.Scale(scale)
			scaledFull := window.SizeFull.Scale(scale)
			window.Size = *ImFloorVec(&scaled)
			window.SizeFull = *ImFloorVec(&scaledFull)
		}
		return
	}

	// Mouse wheel scrolling
	// If a child window has the ImGuiWindowFlags_NoScrollWithMouse flag, we give a chance to scroll its parent
	if g.IO.KeyCtrl {
		return
	}

	// As a standard behavior holding SHIFT while using Vertical Mouse Wheel triggers Horizontal scroll instead
	// (we avoid doing it on OSX as it the OS input layer handles this already)
	var swap_axis bool = g.IO.KeyShift && !g.IO.ConfigMacOSXBehaviors

	var wheel_y float
	if !swap_axis {
		wheel_y = g.IO.MouseWheel
	}
	var wheel_x float
	if swap_axis {
		wheel_x = g.IO.MouseWheel
	} else {
		wheel_x = g.IO.MouseWheelH
	}

	// Vertical Mouse Wheel scrolling
	if wheel_y != 0.0 {
		StartLockWheelingWindow(window)
		for (window.Flags&ImGuiWindowFlags_ChildWindow != 0) && ((window.ScrollMax.y == 0.0) || ((window.Flags&ImGuiWindowFlags_NoScrollWithMouse != 0) && (window.Flags&ImGuiWindowFlags_NoMouseInputs == 0))) {
			window = window.ParentWindow
		}
		if (window.Flags&ImGuiWindowFlags_NoScrollWithMouse == 0) && (window.Flags&ImGuiWindowFlags_NoMouseInputs == 0) {
			var max_step float = window.InnerRect.GetHeight() * 0.67
			var scroll_step float = ImFloor(ImMin(5*window.CalcFontSize(), max_step))
			setScrollY(window, window.Scroll.y-wheel_y*scroll_step)
		}
	}

	// Horizontal Mouse Wheel scrolling, or Vertical Mouse Wheel w/ Shift held
	if wheel_x != 0.0 {
		StartLockWheelingWindow(window)
		for (window.Flags&ImGuiWindowFlags_ChildWindow != 0) && ((window.ScrollMax.x == 0.0) || ((window.Flags&ImGuiWindowFlags_NoScrollWithMouse != 0) && (window.Flags&ImGuiWindowFlags_NoMouseInputs == 0))) {
			window = window.ParentWindow
		}
		if (window.Flags&ImGuiWindowFlags_NoScrollWithMouse == 0) && (window.Flags&ImGuiWindowFlags_NoMouseInputs == 0) {
			var max_step float = window.InnerRect.GetWidth() * 0.67
			var scroll_step float = ImFloor(ImMin(2*window.CalcFontSize(), max_step))
			setScrollX(window, window.Scroll.x-wheel_x*scroll_step)
		}
	}
}

func UpdateMouseInputs() {
	var g = GImGui

	// Round mouse position to avoid spreading non-rounded position (e.g. UpdateManualResize doesn't support them well)
	if IsMousePosValid(&g.IO.MousePos) {
		g.IO.MousePos = *ImFloorVec(&g.IO.MousePos)
		g.MouseLastValidPos = *ImFloorVec(&g.IO.MousePos)
	}

	// If mouse just appeared or disappeared (usually denoted by -FLT_MAX components) we cancel out movement in MouseDelta
	if IsMousePosValid(&g.IO.MousePos) && IsMousePosValid(&g.IO.MousePosPrev) {
		g.IO.MouseDelta = g.IO.MousePos.Sub(g.IO.MousePosPrev)
	} else {
		g.IO.MouseDelta = ImVec2{}
	}

	// If mouse moved we re-enable mouse hovering in case it was disabled by gamepad/keyboard. In theory should use a >0.0f threshold but would need to reset in everywhere we set this to true.
	if g.IO.MouseDelta.x != 0.0 || g.IO.MouseDelta.y != 0.0 {
		g.NavDisableMouseHover = false
	}

	g.IO.MousePosPrev = g.IO.MousePos
	for i := 0; i < len(g.IO.MouseDown); i++ {
		g.IO.MouseClicked[i] = g.IO.MouseDown[i] && g.IO.MouseDownDuration[i] < 0.0
		g.IO.MouseReleased[i] = !g.IO.MouseDown[i] && g.IO.MouseDownDuration[i] >= 0.0
		g.IO.MouseDownDurationPrev[i] = g.IO.MouseDownDuration[i]
		if g.IO.MouseDown[i] {
			if g.IO.MouseDownDuration[i] < 0.0 {
				g.IO.MouseDownDuration[i] = 0.0
			} else {
				g.IO.MouseDownDuration[i] += g.IO.DeltaTime
			}
		} else {
			g.IO.MouseDownDuration[i] = -1.0
		}
		g.IO.MouseDoubleClicked[i] = false
		if g.IO.MouseClicked[i] {
			if (float)(g.Time-g.IO.MouseClickedTime[i]) < g.IO.MouseDoubleClickTime {
				var delta_from_click_pos ImVec2
				if IsMousePosValid(&g.IO.MousePos) {
					delta_from_click_pos = (g.IO.MousePos.Sub(g.IO.MouseClickedPos[i]))
				}
				if ImLengthSqrVec2(delta_from_click_pos) < g.IO.MouseDoubleClickMaxDist*g.IO.MouseDoubleClickMaxDist {
					g.IO.MouseDoubleClicked[i] = true
				}
				g.IO.MouseClickedTime[i] = -double(g.IO.MouseDoubleClickTime) * 2.0 // Mark as "old enough" so the third click isn't turned into a double-click
			} else {
				g.IO.MouseClickedTime[i] = g.Time
			}
			g.IO.MouseClickedPos[i] = g.IO.MousePos
			g.IO.MouseDownWasDoubleClick[i] = g.IO.MouseDoubleClicked[i]
			g.IO.MouseDragMaxDistanceAbs[i] = ImVec2{}
			g.IO.MouseDragMaxDistanceSqr[i] = 0.0
		} else if g.IO.MouseDown[i] {
			// Maintain the maximum distance we reaching from the initial click position, which is used with dragging threshold
			var delta_from_click_pos ImVec2
			if IsMousePosValid(&g.IO.MousePos) {
				delta_from_click_pos = (g.IO.MousePos.Sub(g.IO.MouseClickedPos[i]))
			}
			g.IO.MouseDragMaxDistanceSqr[i] = ImMax(g.IO.MouseDragMaxDistanceSqr[i], ImLengthSqrVec2(delta_from_click_pos))

			var dx, dy float
			if delta_from_click_pos.x < 0 {
				dx = -delta_from_click_pos.x
			}
			if delta_from_click_pos.y < 0 {
				dy = -delta_from_click_pos.y
			}

			g.IO.MouseDragMaxDistanceAbs[i].x = ImMax(g.IO.MouseDragMaxDistanceAbs[i].x, dx)
			g.IO.MouseDragMaxDistanceAbs[i].y = ImMax(g.IO.MouseDragMaxDistanceAbs[i].y, dy)
		}
		if !g.IO.MouseDown[i] && !g.IO.MouseReleased[i] {
			g.IO.MouseDownWasDoubleClick[i] = false
		}
		if g.IO.MouseClicked[i] { // Clicking any mouse button reactivate mouse hovering which may have been deactivated by gamepad/keyboard navigation
			g.NavDisableMouseHover = false
		}
	}
}