package imgui

const WINDOWS_MOUSE_WHEEL_SCROLL_LOCK_TIMER float = 2.00 // Lock scrolled window (so it doesn't pick child windows that are scrolling through) for a certain time, unless mouse moved.

// Inputs Utilities: Mouse
// - To refer to a mouse button, you may use named enums in your code e.g. ImGuiMouseButton_Left, ImGuiMouseButton_Right.
// - You can also use regular integer: it is forever guaranteed that 0=Left, 1=Right, 2=Middle.
// - Dragging operations are only reported after mouse has moved a certain distance away from the initial clicking position (see 'lock_threshold' and 'io.MouseDraggingThreshold')

// IsMouseDown is mouse button held?
func IsMouseDown(button ImGuiMouseButton) bool {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	return g.IO.MouseDown[button]
}

// IsMouseReleased did mouse button released? (went from Down to !Down)
func IsMouseReleased(button ImGuiMouseButton) bool {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	return g.IO.MouseReleased[button]
}

// IsMouseDoubleClicked did mouse button double-clicked? (note that a double-click will also report IsMouseClicked() == true)
func IsMouseDoubleClicked(button ImGuiMouseButton) bool {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	return g.IO.MouseDoubleClicked[button]
}

// IsAnyMouseDown is any mouse button held?
func IsAnyMouseDown() bool {
	var g = GImGui
	for n := range g.IO.MouseDown {
		if g.IO.MouseDown[n] {
			return true
		}
	}
	return false
}

// GetMousePos shortcut to ImGui::GetIO().MousePos provided by user, to be consistent with other calls
func GetMousePos() ImVec2 {
	var g = GImGui
	return g.IO.MousePos
}

// GetMousePosOnOpeningCurrentPopup retrieve mouse position at the time of opening popup we have BeginPopup() into (helper to a user backing that value themselves)
func GetMousePosOnOpeningCurrentPopup() ImVec2 {
	var g = GImGui
	if len(g.BeginPopupStack) > 0 {
		return g.OpenPopupStack[len(g.BeginPopupStack)-1].OpenMousePos
	}
	return g.IO.MousePos
}

// GetMouseDragDelta Return the delta from the initial clicking position while the mouse button is clicked or was just released.
// This is locked and return 0.0f until the mouse moves past a distance threshold at least once.
// NB: This is only valid if IsMousePosValid(). backends in theory should always keep mouse position valid when dragging even outside the client window.
// return the delta from the initial clicking position while the mouse button is pressed or was just released. This is locked and return 0.0 until the mouse moves past a distance threshold at least once (if lock_threshold < -1.0, uses io.MouseDraggingThreshold)
func GetMouseDragDelta(button ImGuiMouseButton /*= 0*/, lock_threshold float /*= -1.0*/) ImVec2 {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	if lock_threshold < 0.0 {
		lock_threshold = g.IO.MouseDragThreshold
	}
	if g.IO.MouseDown[button] || g.IO.MouseReleased[button] {
		if g.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold*lock_threshold {
			if IsMousePosValid(&g.IO.MousePos) && IsMousePosValid(&g.IO.MouseClickedPos[button]) {
				return g.IO.MousePos.Sub(g.IO.MouseClickedPos[button])
			}
		}
	}
	return ImVec2{0.0, 0.0}
}

func ResetMouseDragDelta(button ImGuiMouseButton) {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	// NB: We don't need to reset g.IO.MouseDragMaxDistanceSqr
	g.IO.MouseClickedPos[button] = g.IO.MousePos
}

// GetMouseCursor get desired cursor type, reset in ImGui::NewFrame(), this is updated during the frame. valid before Render(). If you use software rendering by setting io.MouseDrawCursor ImGui will render those for you
func GetMouseCursor() ImGuiMouseCursor {
	return GImGui.MouseCursor
}

// CaptureMouseFromApp attention: misleading name! manually override io.WantCaptureMouse flag next frame (said flag is entirely left for your application to handle). This is equivalent to setting "io.WantCaptureMouse = want_capture_mouse_value  {panic("not implemented")}" after the next NewFrame() call.
func CaptureMouseFromApp(want_capture_mouse_value bool /*= true*/) {
	if want_capture_mouse_value {
		GImGui.WantCaptureMouseNextFrame = 1
	} else {
		GImGui.WantCaptureMouseNextFrame = 0
	}
}

// IsMousePosValid We typically use ImVec2(-FLT_MAX,-FLT_MAX) to denote an invalid mouse position.
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

// IsMouseClicked did mouse button clicked? (went from !Down to Down)
func IsMouseClicked(button ImGuiMouseButton, repeat bool) bool {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	var t float = g.IO.MouseDownDuration[button]
	if t == 0.0 {
		return true
	}

	if repeat && t > g.IO.KeyRepeatDelay {
		// FIXME: 2019/05/03: Our old repeat code was wrong here and led to doubling the repeat rate, which made it an ok rate for repeat on mouse hold.
		var amount int = CalcTypematicRepeatAmount(t-g.IO.DeltaTime, t, g.IO.KeyRepeatDelay, g.IO.KeyRepeatRate*0.50)
		if amount > 0 {
			return true
		}
	}
	return false
}

// SetMouseCursor set desired cursor type
func SetMouseCursor(cursor_type ImGuiMouseCursor) {
	GImGui.MouseCursor = cursor_type
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

// IsMouseDragging is mouse dragging? (if lock_threshold < -1.0, uses io.MouseDraggingThreshold)
func IsMouseDragging(button ImGuiMouseButton, lock_threshold float /*= -1.0*/) bool {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	if !g.IO.MouseDown[button] {
		return false
	}
	return IsMouseDragPastThreshold(button, lock_threshold)
}

// IsMouseDragPastThreshold Return if a mouse click/drag went past the given threshold. Valid to call during the MouseReleased frame.
// [Internal] This doesn't test if the button is pressed
func IsMouseDragPastThreshold(button ImGuiMouseButton, lock_threshold float /*= -1.0f*/) bool {
	var g = GImGui
	IM_ASSERT(button >= 0 && button < ImGuiMouseButton(len(g.IO.MouseDown)))
	if lock_threshold < 0.0 {
		lock_threshold = g.IO.MouseDragThreshold
	}
	return g.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold*lock_threshold
}

func StartMouseMovingWindow(window *ImGuiWindow) {
	// Set ActiveId even if the _NoMove flag is set. Without it, dragging away from a window with _NoMove would activate hover on other windows.
	// We _also_ call this when clicking in a window empty space when io.ConfigWindowsMoveFromTitleBarOnly is set, but clear g.MovingWindow afterward.
	// This is because we want ActiveId to be set even when the window is not permitted to move.
	var g = GImGui
	FocusWindow(window)
	SetActiveID(window.MoveId, window)
	g.NavDisableHighlight = true
	g.ActiveIdClickOffset = g.IO.MouseClickedPos[0].Sub(window.RootWindow.Pos)
	g.ActiveIdNoClearOnFocusLoss = true
	SetActiveIdUsingNavAndKeys()

	var can_move_window bool = true
	if (window.Flags&ImGuiWindowFlags_NoMove != 0) || (window.RootWindow.Flags&ImGuiWindowFlags_NoMove != 0) {
		can_move_window = false
	}
	if can_move_window {
		g.MovingWindow = window
	}
}

// IsMouseHoveringRect Test if mouse cursor is hovering given rectangle
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
	return rect_for_touch.ContainsVec(g.IO.MousePos)
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
				if (focused_root_window.Flags&ImGuiWindowFlags_Popup != 0) && flags&ImGuiHoveredFlags_AllowWhenBlockedByPopup == 0 {
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
			if g.IO.ConfigWindowsMoveFromTitleBarOnly && root_window.Flags&ImGuiWindowFlags_NoTitleBar == 0 {
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
