package imgui

// Windows Utilities
// - 'current window' = the window we are appending into while inside a Begin()/End() block. 'next window' = next window we will Begin() into.

func IsWindowAppearing() bool {
	window := GetCurrentWindowRead()
	return window.Appearing
}

func IsWindowCollapsed() bool {
	window := GetCurrentWindowRead()
	return window.Collapsed
}

// Can we focus this window with CTRL+TAB (or PadMenu + PadFocusPrev/PadFocusNext)
// Note that NoNavFocus makes the window not reachable with CTRL+TAB but it can still be focused with mouse or programmatically.
// If you want a window to never be focused, you may use the e.guiContext. NoInputs flag.
func IsWindowNavFocusable(window *ImGuiWindow) bool {
	return window.WasActive && window == window.RootWindow && (window.Flags&ImGuiWindowFlags_NoNavFocus == 0)
}

// is current window focused? or its root/child, depending on flags. see flags for options.
func IsWindowFocused(flags ImGuiFocusedFlags) bool {
	g := guiContext

	if flags&ImGuiFocusedFlags_AnyWindow != 0 {
		return g.NavWindow != nil
	}

	IM_ASSERT(g.CurrentWindow != nil) // Not inside a Begin()/End()
	switch flags & (ImGuiFocusedFlags_RootWindow | ImGuiFocusedFlags_ChildWindows) {
	case ImGuiFocusedFlags_RootWindow | ImGuiFocusedFlags_ChildWindows:
		return g.NavWindow != nil && g.NavWindow.RootWindow == g.CurrentWindow.RootWindow
	case ImGuiFocusedFlags_RootWindow:
		return g.NavWindow == g.CurrentWindow.RootWindow
	case ImGuiFocusedFlags_ChildWindows:
		return g.NavWindow != nil && IsWindowChildOf(g.NavWindow, g.CurrentWindow)
	default:
		return g.NavWindow == g.CurrentWindow
	}
}

// is current window hovered (and typically: not blocked by a popup/modal)? see flags for options. NB: If you are trying to check whether your mouse should be dispatched to imgui or to your app, you should use the 'io.WantCaptureMouse' boolean for that! Please read the FAQ!
func IsWindowHovered(flags ImGuiHoveredFlags) bool {
	IM_ASSERT((flags & ImGuiHoveredFlags_AllowWhenOverlapped) == 0) // Flags not supported by this function
	if guiContext.HoveredWindow == nil {
		return false
	}

	if (flags & ImGuiHoveredFlags_AnyWindow) == 0 {
		window := guiContext.CurrentWindow
		switch flags & (ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows) {
		case ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows:
			if guiContext.HoveredWindow.RootWindow != window.RootWindow {
				return false
			}
		case ImGuiHoveredFlags_RootWindow:
			if guiContext.HoveredWindow != window.RootWindow {
				return false
			}
		case ImGuiHoveredFlags_ChildWindows:
			if !IsWindowChildOf(guiContext.HoveredWindow, window) {
				return false
			}
		default:
			if guiContext.HoveredWindow != window {
				return false
			}
		}
	}

	if !IsWindowContentHoverable(guiContext.HoveredWindow, flags) {
		return false
	}
	if flags&ImGuiHoveredFlags_AllowWhenBlockedByActiveItem == 0 {
		if guiContext.ActiveId != 0 && !guiContext.ActiveIdAllowOverlap && guiContext.ActiveId != guiContext.HoveredWindow.MoveId {
			return false
		}
	}
	return true
}

// get draw list associated to the current window, to append your own drawing primitives
func GetWindowDrawList() *ImDrawList {
	window := GetCurrentWindow()
	return window.DrawList
}

// get current window position in screen space (useful if you want to do your own drawing via the DrawList API)
func GetWindowPos() ImVec2 {
	window := guiContext.CurrentWindow
	return window.Pos
}

// get current window size
func GetWindowSize() ImVec2 {
	window := GetCurrentWindowRead()
	return window.Size
}

// get current window width (shortcut for GetWindowSize().x)
func GetWindowWidth() float {
	window := guiContext.CurrentWindow
	return window.Size.x
}

// get current window height (shortcut for GetWindowSize().y)
func GetWindowHeight() float {
	window := guiContext.CurrentWindow
	return window.Size.y
}
