package imgui

import (
	"sort"

	"github.com/Splizard/imgui/golang"
)

const WINDOWS_HOVER_PADDING float = 4.0 // Extend outside window for hovering/resizing (maxxed with TouchPadding) and inside windows for borders. Affect FindHoveredWindow().

func IsWindowActiveAndVisible(window *ImGuiWindow) bool {
	return (window.Active) && (!window.Hidden)
}

// set next window background color alpha. helper to easily override the Alpha component of ImGuiCol_WindowBg/ChildBg/PopupBg. you may also use ImGuiWindowFlags_NoBackground.
func SetNextWindowBgAlpha(alpha float) {
	g := GImGui
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasBgAlpha
	g.NextWindowData.BgAlphaVal = alpha
}

// Layer is locked for the root window, however child windows may use a different viewport (e.g. extruding menu)
func AddRootWindowToDrawData(window *ImGuiWindow) {
	var layer int
	if window.Flags&ImGuiWindowFlags_Tooltip != 0 {
		layer = 1
	}
	AddWindowToDrawData(window, layer)
}

// Window has already passed the IsWindowNavFocusable()
func GetFallbackWindowNameForWindowingList(window *ImGuiWindow) string {
	if window.Flags&ImGuiWindowFlags_Popup != 0 {
		return "(Popup)"
	}
	if (window.Flags&ImGuiWindowFlags_MenuBar != 0) && window.Name == "##MainMenuBar" {
		return "(Main menu bar)"
	}
	return "(Untitled)"
}

func AddWindowToDrawData(window *ImGuiWindow, layer int) {
	g := GImGui
	var viewport = g.Viewports[0]
	g.IO.MetricsRenderWindows++
	AddDrawListToDrawData(&viewport.DrawDataBuilder[layer], window.DrawList)
	for i := range window.DC.ChildWindows {
		var child = window.DC.ChildWindows[i]
		if IsWindowActiveAndVisible(child) { // Clipped children may have been marked not active
			AddWindowToDrawData(child, layer)
		}
	}
}

// FIXME: Add a more explicit sort order in the window structure.
func ChildWindowComparer(a, b *ImGuiWindow) int {
	if d := (a.Flags & ImGuiWindowFlags_Popup) - (b.Flags & ImGuiWindowFlags_Popup); d != 0 {
		return int(d)
	}
	if d := (a.Flags & ImGuiWindowFlags_Tooltip) - (b.Flags & ImGuiWindowFlags_Tooltip); d != 0 {
		return int(d)
	}
	return int(a.BeginOrderWithinParent - b.BeginOrderWithinParent)
}

func AddWindowToSortBuffer(out_sorted_windows *[]*ImGuiWindow, window *ImGuiWindow) {
	*out_sorted_windows = append(*out_sorted_windows, window)
	if window.Active {
		var count = int(len(window.DC.ChildWindows))
		if count > 1 {
			sort.Slice(window.DC.ChildWindows, func(i, j golang.Int) bool {
				return ChildWindowComparer(window.DC.ChildWindows[i], window.DC.ChildWindows[j]) < 0
			})
		}
		for i := int(0); i < count; i++ {
			var child = window.DC.ChildWindows[i]
			if child.Active {
				AddWindowToSortBuffer(out_sorted_windows, child)
			}
		}
	}
}

func FindWindowByID(id ImGuiID) *ImGuiWindow {
	g := GImGui
	ptr := g.WindowsById.GetInterface(id)
	if ptr == nil {
		return nil
	}
	return ptr.(*ImGuiWindow)
}

func FindWindowByName(e string) *ImGuiWindow {
	var id = ImHashStr(e, 0, 0)
	return FindWindowByID(id)
}

func SetCurrentWindow(window *ImGuiWindow) {
	g := GImGui
	g.CurrentWindow = window
	if window != nil && window.DC.CurrentTableIdx != -1 {
		g.CurrentTable = g.Tables[uint(window.DC.CurrentTableIdx)]
	}
	if window != nil {
		size := window.CalcFontSize()
		g.FontSize = size
		g.DrawListSharedData.FontSize = size
	}

}

func setWindowPos(window *ImGuiWindow, pos *ImVec2, cond ImGuiCond) {
	// Test condition (NB: bit 0 is always true) and clear flags for next time
	if cond != 0 && (window.SetWindowPosAllowFlags&cond) == 0 {
		return
	}

	IM_ASSERT(cond == 0 || ImIsPowerOfTwoInt(int(cond))) // Make sure the user doesn't attempt to combine multiple condition flags.
	window.SetWindowPosAllowFlags &= ^(ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing)
	window.SetWindowPosVal = ImVec2{FLT_MAX, FLT_MAX}

	// Set
	var old_pos = window.Pos
	window.Pos = *ImFloorVec(pos)
	var offset = window.Pos.Sub(old_pos)
	window.DC.CursorPos = window.DC.CursorPos.Add(offset)       // As we happen to move the window while it is being appended to (which is a bad idea - will smear) let's at least offset the cursor
	window.DC.CursorMaxPos = window.DC.CursorMaxPos.Add(offset) // And more importantly we need to offset CursorMaxPos/CursorStartPos this so ContentSize calculation doesn't get affected.
	window.DC.IdealMaxPos = window.DC.IdealMaxPos.Add(offset)
	window.DC.CursorStartPos = window.DC.CursorStartPos.Add(offset)
}

func setWindowSize(window *ImGuiWindow, size *ImVec2, cond ImGuiCond) {
	// Test condition (NB: bit 0 is always true) and clear flags for next time
	if cond != 0 && (window.SetWindowSizeAllowFlags&cond) == 0 {
		return
	}

	IM_ASSERT(cond == 0 || ImIsPowerOfTwoInt(int(cond))) // Make sure the user doesn't attempt to combine multiple condition flags.
	window.SetWindowSizeAllowFlags &= ^(ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing)

	// Set
	if size.x > 0.0 {
		window.AutoFitFramesX = 0
		window.SizeFull.x = IM_FLOOR(size.x)
	} else {
		window.AutoFitFramesX = 2
		window.AutoFitOnlyGrows = false
	}
	if size.y > 0.0 {
		window.AutoFitFramesY = 0
		window.SizeFull.y = IM_FLOOR(size.y)
	} else {
		window.AutoFitFramesY = 2
		window.AutoFitOnlyGrows = false
	}
}

// set next window size. set axis to 0.0 to force an auto-fit on this axis. call before Begin()
func SetNextWindowSize(size *ImVec2, cond ImGuiCond) {
	g := GImGui
	IM_ASSERT(cond == 0 || ImIsPowerOfTwoInt(int(cond))) // Make sure the user doesn't attempt to combine multiple condition flags.
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasSize
	g.NextWindowData.SizeVal = *size
	if cond != 0 {
		g.NextWindowData.SizeCond = cond
	} else {
		g.NextWindowData.SizeCond = ImGuiCond_Always
	}
}

func SetWindowConditionAllowFlags(window *ImGuiWindow, flags ImGuiCond, enabled bool) {
	if enabled {
		window.SetWindowPosAllowFlags |= flags
		window.SetWindowSizeAllowFlags |= flags
		window.SetWindowCollapsedAllowFlags |= flags
	} else {
		window.SetWindowPosAllowFlags &= ^flags
		window.SetWindowSizeAllowFlags &= ^flags
		window.SetWindowCollapsedAllowFlags &= ^flags
	}
}

func CreateNewWindow(name string, flags ImGuiWindowFlags) *ImGuiWindow {
	g := GImGui

	// Create window the first time
	var window = NewImGuiWindow(g, name)
	window.Flags = flags
	g.WindowsById.SetInterface(window.ID, window)

	// Default/arbitrary window position. Use SetNextWindowPos() with the appropriate condition flag to change the initial position of a window.
	var main_viewport = GetMainViewport()
	window.Pos = main_viewport.Pos.Add(ImVec2{60, 60})

	// User can disable loading and saving of settings. Tooltip and child windows also don't store settings.
	if flags&ImGuiWindowFlags_NoSavedSettings == 0 {
		if settings := FindWindowSettings(window.ID); settings != nil {
			// Retrieve settings from .ini file
			for i := range g.SettingsWindows {
				if &g.SettingsWindows[i] == settings {
					window.SettingsOffset = int(i)
					break
				}
			}

			SetWindowConditionAllowFlags(window, ImGuiCond_FirstUseEver, false)
			ApplyWindowSettings(window, settings)
		}
	}
	window.DC.CursorStartPos = window.Pos // So first call to CalcContentSize() doesn't return crazy values
	window.DC.CursorMaxPos = window.Pos   // So first call to CalcContentSize() doesn't return crazy values

	if (flags & ImGuiWindowFlags_AlwaysAutoResize) != 0 {
		window.AutoFitFramesX = 2
		window.AutoFitFramesY = 2
		window.AutoFitOnlyGrows = false
	} else {
		if window.Size.x <= 0.0 {
			window.AutoFitFramesX = 2
		}
		if window.Size.y <= 0.0 {
			window.AutoFitFramesY = 2
		}
		window.AutoFitOnlyGrows = (window.AutoFitFramesX > 0) || (window.AutoFitFramesY > 0)
	}

	if flags&ImGuiWindowFlags_ChildWindow == 0 {
		g.WindowsFocusOrder = append(g.WindowsFocusOrder, window)
		window.FocusOrder = (short)(len(g.WindowsFocusOrder) - 1)
	}

	if flags&ImGuiWindowFlags_NoBringToFrontOnFocus != 0 {
		g.Windows = append([]*ImGuiWindow{window}, g.Windows...)
	} else {
		g.Windows = append(g.Windows, window)
	}
	return window
}

func GetWindowBgColorIdxFromFlags(flags ImGuiWindowFlags) ImGuiCol {
	if flags&(ImGuiWindowFlags_Tooltip|ImGuiWindowFlags_Popup) != 0 {
		return ImGuiCol_PopupBg
	}
	if flags&ImGuiWindowFlags_ChildWindow != 0 {
		return ImGuiCol_ChildBg
	}
	return ImGuiCol_WindowBg
}

func CalcWindowAutoFitSize(window *ImGuiWindow, size_contents *ImVec2) ImVec2 {
	g := GImGui
	style := g.Style
	var decoration_up_height = window.TitleBarHeight() + window.MenuBarHeight()
	var size_pad = window.WindowPadding.Scale(2)
	var size_desired = size_contents.Add(size_pad).Add(ImVec2{0.0, decoration_up_height})
	if window.Flags&ImGuiWindowFlags_Tooltip != 0 {
		// Tooltip always resize
		return size_desired
	} else {
		// Maximum window size is determined by the viewport size or monitor size
		var is_popup = (window.Flags & ImGuiWindowFlags_Popup) != 0
		var is_menu = (window.Flags & ImGuiWindowFlags_ChildMenu) != 0
		var size_min = style.WindowMinSize
		if is_popup || is_menu { // Popups and menus bypass style.WindowMinSize by default, but we give then a non-zero minimum size to facilitate understanding problematic cases (e.g. empty popups)
			size_min = ImMinVec2(&size_min, &ImVec2{4.0, 4.0})
		}

		// FIXME-VIEWPORT-WORKAREA: May want to use GetWorkSize() instead of Size depending on the type of windows?
		var avail_size = GetMainViewport().Size
		var s = avail_size.Sub(style.DisplaySafeAreaPadding.Scale(2.0))
		var size_auto_fit = ImClampVec2(&size_desired, &size_min, ImMaxVec2(&size_min, &s))

		// When the window cannot fit all contents (either because of constraints, either because screen is too small),
		// we are growing the size on the other axis to compensate for expected scrollbar. FIXME: Might turn bigger than ViewportSize-WindowPadding.
		var size_auto_fit_after_constraint = CalcWindowSizeAfterConstraint(window, &size_auto_fit)
		var will_have_scrollbar_x = (size_auto_fit_after_constraint.x-size_pad.x-0.0 < size_contents.x && (window.Flags&ImGuiWindowFlags_NoScrollbar == 0) && (window.Flags&ImGuiWindowFlags_HorizontalScrollbar != 0)) || (window.Flags&ImGuiWindowFlags_AlwaysHorizontalScrollbar != 0)
		var will_have_scrollbar_y = (size_auto_fit_after_constraint.y-size_pad.y-decoration_up_height < size_contents.y && (window.Flags&ImGuiWindowFlags_NoScrollbar == 0) || (window.Flags&ImGuiWindowFlags_AlwaysVerticalScrollbar != 0))
		if will_have_scrollbar_x {
			size_auto_fit.y += style.ScrollbarSize
		}
		if will_have_scrollbar_y {
			size_auto_fit.x += style.ScrollbarSize
		}
		return size_auto_fit
	}
}

func ClampWindowRect(window *ImGuiWindow, visibility_rect *ImRect) {
	g := GImGui
	var size_for_clamping = window.Size
	if g.IO.ConfigWindowsMoveFromTitleBarOnly && window.Flags&ImGuiWindowFlags_NoTitleBar == 0 {
		size_for_clamping.y = window.TitleBarHeight()
	}
	sub := visibility_rect.Min.Sub(size_for_clamping)
	window.Pos = ImClampVec2(&window.Pos, &sub, visibility_rect.Max)
}

func CalcWindowSizeAfterConstraint(window *ImGuiWindow, size_desired *ImVec2) ImVec2 {
	g := GImGui
	var new_size = *size_desired
	if (g.NextWindowData.Flags & ImGuiNextWindowDataFlags_HasSizeConstraint) != 0 {
		// Using -1,-1 on either X/Y axis to preserve the current size.
		var cr = g.NextWindowData.SizeConstraintRect
		if cr.Min.x >= 0 && cr.Max.x >= 0 {
			new_size.x = ImClamp(new_size.x, cr.Min.x, cr.Max.x)
		} else {
			new_size.x = window.SizeFull.x
		}
		if cr.Min.y >= 0 && cr.Max.y >= 0 {
			new_size.y = ImClamp(new_size.y, cr.Min.y, cr.Max.y)
		} else {
			new_size.y = window.SizeFull.y
		}
		if g.NextWindowData.SizeCallback != nil {
			var data ImGuiSizeCallbackData
			data.UserData = g.NextWindowData.SizeCallbackUserData
			data.Pos = window.Pos
			data.CurrentSize = window.SizeFull
			data.DesiredSize = new_size
			g.NextWindowData.SizeCallback(&data)
			new_size = data.DesiredSize
		}
		new_size.x = IM_FLOOR(new_size.x)
		new_size.y = IM_FLOOR(new_size.y)
	}

	// Minimum size
	if window.Flags&(ImGuiWindowFlags_ChildWindow|ImGuiWindowFlags_AlwaysAutoResize) == 0 {
		var window_for_height = window
		var decoration_up_height = window_for_height.TitleBarHeight() + window_for_height.MenuBarHeight()
		new_size = ImMaxVec2(&new_size, &g.Style.WindowMinSize)
		new_size.y = max(new_size.y, decoration_up_height+max(0.0, g.Style.WindowRounding-1.0)) // Reduce artifacts with very small windows
	}
	return new_size
}

func CalcWindowContentSizes(window *ImGuiWindow, content_size_current, content_size_ideal *ImVec2) {
	var preserve_old_content_sizes = false
	if window.Collapsed && window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0 {
		preserve_old_content_sizes = true
	} else if window.Hidden && window.HiddenFramesCannotSkipItems == 0 && window.HiddenFramesCanSkipItems > 0 {
		preserve_old_content_sizes = true
	}
	if preserve_old_content_sizes {
		*content_size_current = window.ContentSize
		*content_size_ideal = window.ContentSizeIdeal
		return
	}
	if window.ContentSizeExplicit.x != 0.0 {
		content_size_current.x = window.ContentSizeExplicit.x
		content_size_current.y = window.ContentSizeExplicit.y
		content_size_ideal.x = window.ContentSizeExplicit.x
		content_size_ideal.y = window.ContentSizeExplicit.y
	} else {
		content_size_current.x = IM_FLOOR(window.DC.CursorMaxPos.x - window.DC.CursorStartPos.x)
		content_size_current.y = IM_FLOOR(window.DC.CursorMaxPos.y - window.DC.CursorStartPos.y)
		content_size_ideal.x = IM_FLOOR(max(window.DC.CursorMaxPos.x, window.DC.IdealMaxPos.x) - window.DC.CursorStartPos.x)
		content_size_ideal.y = IM_FLOOR(max(window.DC.CursorMaxPos.y, window.DC.IdealMaxPos.y) - window.DC.CursorStartPos.y)
	}
}

// Windows
// - Begin() = push window to the stack and start appending to it. End() = pop window from the stack.
// - Passing 'bool* p_open != NULL' shows a window-closing widget in the upper-right corner of the window,
//   which clicking will set the boolean to false when clicked.
// - You may append multiple times to the same window during the same frame by calling Begin()/End() pairs multiple times.
//   Some information such as 'flags' or 'p_open' will only be considered by the first call to Begin().
// - Begin() return false to indicate the window is collapsed or fully clipped, so you may early out and omit submitting
//   anything to the window. Always call a matching End() for each Begin() call, regardless of its return value!
//   [Important: due to legacy reason, this is inconsistent with most other functions such as BeginMenu/EndMenu,
//    BeginPopup/EndPopup, etc. where the EndXXX call should only be called if the corresponding BeginXXX function
//    returned true. Begin and BeginChild are the only odd ones out. Will be fixed in a future update.]
// - Note that the bottom of window stack always contains a window called "Debug".

func End() {
	g := GImGui
	window := g.CurrentWindow

	window.StateStorage = window.DC.StateStorage

	// Error checking: verify that user hasn't called End() too many times!
	if len(g.CurrentWindowStack) <= 1 && g.WithinFrameScopeWithImplicitWindow {
		IM_ASSERT_USER_ERROR(len(g.CurrentWindowStack) > 1, "Calling End() too many times!")
		return
	}
	IM_ASSERT(len(g.CurrentWindowStack) > 0)

	// Error checking: verify that user doesn't directly call End() on a child window.
	if window.Flags&ImGuiWindowFlags_ChildWindow != 0 {
		IM_ASSERT_USER_ERROR(g.WithinEndChild, "Must call EndChild() and not End()!")
	}

	// Close anything that is open
	if window.DC.CurrentColumns != nil {
		EndColumns()
	}
	PopClipRect() // Inner window clip rectangle

	// Stop logging
	if window.Flags&ImGuiWindowFlags_ChildWindow == 0 { // FIXME: add more options for scope of logging
		LogFinish()
	}

	// Pop from window stack
	g.LastItemData = g.CurrentWindowStack[len(g.CurrentWindowStack)-1].ParentLastItemDataBackup
	g.CurrentWindowStack = g.CurrentWindowStack[:len(g.CurrentWindowStack)-1]
	if window.Flags&ImGuiWindowFlags_Popup != 0 {
		g.BeginPopupStack = g.BeginPopupStack[:len(g.BeginPopupStack)-1]
	}
	window.DC.StackSizesOnBegin.CompareWithCurrentState()

	var current *ImGuiWindow
	if len(g.CurrentWindowStack) > 0 {
		current = g.CurrentWindowStack[len(g.CurrentWindowStack)-1].Window
	}
	SetCurrentWindow(current)
}

// Find window given position, search front-to-back
// FIXME: Note that we have an inconsequential lag here: OuterRectClipped is updated in Begin(), so windows moved programmatically
// with SetWindowPos() and not SetNextWindowPos() will have that rectangle lagging by a frame at the time FindHoveredWindow() is
// called, aka before the next Begin(). Moving window isn't affected.
func FindHoveredWindow() {
	g := GImGui

	var hovered_window *ImGuiWindow = nil
	var hovered_window_ignoring_moving_window *ImGuiWindow = nil
	if g.MovingWindow != nil && (g.MovingWindow.Flags&ImGuiWindowFlags_NoMouseInputs == 0) {
		hovered_window = g.MovingWindow
	}

	var padding_regular = g.Style.TouchExtraPadding
	var padding_for_resize ImVec2
	if g.IO.ConfigWindowsResizeFromEdges {
		padding_for_resize = g.WindowsHoverPadding
	} else {
		padding_for_resize = padding_regular
	}

	for i := len(g.Windows) - 1; i >= 0; i-- {
		var window = g.Windows[i]

		if !window.Active || window.Hidden {
			continue
		}
		if window.Flags&ImGuiWindowFlags_NoMouseInputs != 0 {
			continue
		}

		// Using the clipped AABB, a child window will typically be clipped by its parent (not always)
		var bb = ImRect(window.OuterRectClipped)
		if window.Flags&(ImGuiWindowFlags_ChildWindow|ImGuiWindowFlags_NoResize|ImGuiWindowFlags_AlwaysAutoResize) != 0 {
			bb.ExpandVec(padding_regular)
		} else {
			bb.ExpandVec(padding_for_resize)
		}
		if !bb.ContainsVec(g.IO.MousePos) {
			continue
		}

		// Support for one rectangular hole in any given window
		// FIXME: Consider generalizing hit-testing override (with more generic data, callback, etc.) (#1512)
		if window.HitTestHoleSize.x != 0 {
			var hole_pos = ImVec2{window.Pos.x + (float)(window.HitTestHoleOffset.x), window.Pos.y + (float)(window.HitTestHoleOffset.y)}
			var hole_size = ImVec2{(float)(window.HitTestHoleSize.x), (float)(window.HitTestHoleSize.y)}
			if (&ImRect{hole_pos, hole_pos.Add(hole_size)}).ContainsVec(g.IO.MousePos) {
				continue
			}
		}

		if hovered_window == nil {
			hovered_window = window
		}
		if hovered_window_ignoring_moving_window == nil && (g.MovingWindow == nil || window.RootWindow != g.MovingWindow.RootWindow) {
			hovered_window_ignoring_moving_window = window
		}
		if hovered_window != nil && hovered_window_ignoring_moving_window != nil {
			break
		}
	}

	g.HoveredWindow = hovered_window
	g.HoveredWindowUnderMovingWindow = hovered_window_ignoring_moving_window
}

// The reason this is exposed in imgui_internal.h is: on touch-based system that don't have hovering, we want to dispatch inputs to the right target (imgui vs imgui+app)
func UpdateHoveredWindowAndCaptureFlags() {
	g := GImGui
	io := g.IO
	g.WindowsHoverPadding = ImMaxVec2(&g.Style.TouchExtraPadding, &ImVec2{WINDOWS_HOVER_PADDING, WINDOWS_HOVER_PADDING})

	// Find the window hovered by mouse:
	// - Child windows can extend beyond the limit of their parent so we need to derive HoveredRootWindow from HoveredWindow.
	// - When moving a window we can skip the search, which also conveniently bypasses the fact that window.WindowRectClipped is lagging as this point of the frame.
	// - We also support the moved window toggling the NoInputs flag after moving has started in order to be able to detect windows below it, which is useful for e.g. docking mechanisms.
	var clear_hovered_windows = false
	FindHoveredWindow()

	// Modal windows prevents mouse from hovering behind them.
	var modal_window = GetTopMostPopupModal()
	if modal_window != nil && g.HoveredWindow != nil && !IsWindowChildOf(g.HoveredWindow.RootWindow, modal_window) {
		clear_hovered_windows = true
	}

	// Disabled mouse?
	if io.ConfigFlags&ImGuiConfigFlags_NoMouse != 0 {
		clear_hovered_windows = true
	}

	// We track click ownership. When clicked outside of a window the click is owned by the application and
	// won't report hovering nor request capture even while dragging over our windows afterward.
	var has_open_popup = (len(g.OpenPopupStack) > 0)
	var has_open_modal = (modal_window != nil)
	var mouse_earliest_down int = -1
	var mouse_any_down = false
	for i := range io.MouseDown {
		if io.MouseClicked[i] {
			io.MouseDownOwned[i] = (g.HoveredWindow != nil) || has_open_popup
			io.MouseDownOwnedUnlessPopupClose[i] = (g.HoveredWindow != nil) || has_open_modal
		}
		mouse_any_down = mouse_any_down || io.MouseDown[i]
		if io.MouseDown[i] {
			if mouse_earliest_down == -1 || io.MouseClickedTime[i] < io.MouseClickedTime[mouse_earliest_down] {
				mouse_earliest_down = int(i)
			}
		}
	}
	var mouse_avail = (mouse_earliest_down == -1) || io.MouseDownOwned[mouse_earliest_down]
	var mouse_avail_unless_popup_close = (mouse_earliest_down == -1) || io.MouseDownOwnedUnlessPopupClose[mouse_earliest_down]

	// If mouse was first clicked outside of ImGui bounds we also cancel out hovering.
	// FIXME: For patterns of drag and drop across OS windows, we may need to rework/remove this test (first committed 311c0ca9 on 2015/02)
	var mouse_dragging_extern_payload = g.DragDropActive && (g.DragDropSourceFlags&ImGuiDragDropFlags_SourceExtern) != 0
	if !mouse_avail && !mouse_dragging_extern_payload {
		clear_hovered_windows = true
	}

	if clear_hovered_windows {
		g.HoveredWindow = nil
		g.HoveredWindowUnderMovingWindow = nil
	}

	// Update io.WantCaptureMouse for the user application (true = dispatch mouse info to Dear ImGui only, false = dispatch mouse to Dear ImGui + underlying app)
	// Update io.WantCaptureMouseAllowPopupClose (experimental) to give a chance for app to react to popup closure with a drag
	if g.WantCaptureMouseNextFrame != -1 {
		io.WantCaptureMouse = (g.WantCaptureMouseNextFrame != 0)
		io.WantCaptureMouseUnlessPopupClose = (g.WantCaptureMouseNextFrame != 0)
	} else {
		io.WantCaptureMouse = (mouse_avail && (g.HoveredWindow != nil || mouse_any_down)) || has_open_popup
		io.WantCaptureMouseUnlessPopupClose = (mouse_avail_unless_popup_close && (g.HoveredWindow != nil || mouse_any_down)) || has_open_modal
	}

	// Update io.WantCaptureKeyboard for the user application (true = dispatch keyboard info to Dear ImGui only, false = dispatch keyboard info to Dear ImGui + underlying app)
	if g.WantCaptureKeyboardNextFrame != -1 {
		io.WantCaptureKeyboard = (g.WantCaptureKeyboardNextFrame != 0)
	} else {
		io.WantCaptureKeyboard = (g.ActiveId != 0) || (modal_window != nil)
	}
	if io.NavActive && (io.ConfigFlags&ImGuiConfigFlags_NavEnableKeyboard != 0) && (io.ConfigFlags&ImGuiConfigFlags_NavNoCaptureKeyboard == 0) {
		io.WantCaptureKeyboard = true
	}

	// Update io.WantTextInput flag, this is to allow systems without a keyboard (e.g. mobile, hand-held) to show a software keyboard if possible
	if g.WantTextInputNextFrame != -1 {
		io.WantTextInput = (g.WantTextInputNextFrame != 0)
	} else {
		io.WantTextInput = false
	}
}

// Handle mouse moving window
// Note: moving window with the navigation keys (Square + d-pad / CTRL+TAB + Arrows) are processed in NavUpdateWindowing()
// FIXME: We don't have strong guarantee that g.MovingWindow stay synched with g.ActiveId == g.MovingWindow.MoveId.
// This is currently enforced by the fact that BeginDragDropSource() is setting all g.ActiveIdUsingXXXX flags to inhibit navigation inputs,
// but if we should more thoroughly test cases where g.ActiveId or g.MovingWindow gets changed and not the other.
func UpdateMouseMovingWindowNewFrame() {
	g := GImGui
	if g.MovingWindow != nil {
		// We actually want to move the root window. g.MovingWindow == window we clicked on (could be a child window).
		// We track it to preserve Focus and so that generally ActiveIdWindow == MovingWindow and ActiveId == MovingWindow.MoveId for consistency.
		KeepAliveID(g.ActiveId)
		IM_ASSERT(g.MovingWindow != nil && g.MovingWindow.RootWindow != nil)
		var moving_window = g.MovingWindow.RootWindow
		if g.IO.MouseDown[0] && IsMousePosValid(&g.IO.MousePos) {
			var pos = g.IO.MousePos.Sub(g.ActiveIdClickOffset)
			if moving_window.Pos.x != pos.x || moving_window.Pos.y != pos.y {
				MarkIniSettingsDirtyWindow(moving_window)
				setWindowPos(moving_window, &pos, ImGuiCond_Always)
			}
			FocusWindow(g.MovingWindow)
		} else {
			g.MovingWindow = nil
			ClearActiveID()
		}
	} else {
		// When clicking/dragging from a window that has the _NoMove flag, we still set the ActiveId in order to prevent hovering others.
		if g.ActiveIdWindow != nil && g.ActiveIdWindow.MoveId == g.ActiveId {
			KeepAliveID(g.ActiveId)
			if !g.IO.MouseDown[0] {
				ClearActiveID()
			}
		}
	}
}
