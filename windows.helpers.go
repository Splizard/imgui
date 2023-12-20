package imgui

import "unsafe"

// Windows
// We should always have a CurrentWindow in the stack (there is an implicit "Debug" window)
// If this ever crash because g.CurrentWindow is NULL it means that either
// - ImGui::NewFrame() has never been called, which is illegal.
// - You are calling ImGui functions after ImGui::EndFrame()/ImGui::Render() and before the next ImGui::NewFrame(), which is also illegal.
func GetCurrentWindowRead() *ImGuiWindow {
	var g = GImGui
	return g.CurrentWindow
}
func GetCurrentWindow() *ImGuiWindow {
	var g = GImGui
	g.CurrentWindow.WriteAccessed = true
	return g.CurrentWindow
}

func UpdateWindowParentAndRootLinks(window *ImGuiWindow, flags ImGuiWindowFlags, parent_window *ImGuiWindow) {
	window.ParentWindow = parent_window
	window.RootWindow = window
	window.RootWindowForTitleBarHighlight = window
	window.RootWindowForNav = window
	if parent_window != nil && (flags&ImGuiWindowFlags_ChildWindow != 0) && flags&ImGuiWindowFlags_Tooltip == 0 {
		window.RootWindow = parent_window.RootWindow
	}
	if parent_window != nil && flags&ImGuiWindowFlags_Modal == 0 && (flags&(ImGuiWindowFlags_ChildWindow|ImGuiWindowFlags_Popup) != 0) {
		window.RootWindowForTitleBarHighlight = parent_window.RootWindowForTitleBarHighlight
	}
	for window.RootWindowForNav.Flags&ImGuiWindowFlags_NavFlattened != 0 {
		IM_ASSERT(window.RootWindowForNav.ParentWindow != nil)
		window.RootWindowForNav = window.RootWindowForNav.ParentWindow
	}
}

func CalcWindowNextAutoFitSize(window *ImGuiWindow) ImVec2 {
	var size_contents_current ImVec2
	var size_contents_ideal ImVec2
	CalcWindowContentSizes(window, &size_contents_current, &size_contents_ideal)
	var size_auto_fit = CalcWindowAutoFitSize(window, &size_contents_ideal)
	var size_final = CalcWindowSizeAfterConstraint(window, &size_auto_fit)
	return size_final
}

func IsWindowChildOf(window *ImGuiWindow, potential_parent *ImGuiWindow) bool {
	if window.RootWindow == potential_parent {
		return true
	}
	for window != nil {
		if window == potential_parent {
			return true
		}
		window = window.ParentWindow
	}
	return false
}

func IsWindowAbove(potential_above *ImGuiWindow, potential_below *ImGuiWindow) bool {
	var g = GImGui
	for i := len(g.Windows) - 1; i >= 0; i-- {
		var candidate_window = g.Windows[i]
		if candidate_window == potential_above {
			return true
		}
		if candidate_window == potential_below {
			return false
		}
	}
	return false
}

func SetWindowHitTestHole(window *ImGuiWindow, pos *ImVec2, size *ImVec2) {
	IM_ASSERT(window.HitTestHoleSize.x == 0) // We don't support multiple holes/hit test filters
	window.HitTestHoleSize = ImVec2ih{x: int16(size.x), y: int16(size.y)}
	diff := pos.Sub(window.Pos)
	window.HitTestHoleOffset = ImVec2ih{int16(diff.x), int16(diff.y)}
}

func BringWindowToDisplayBack(window *ImGuiWindow) {
	var g = GImGui
	if g.Windows[0] == window {
		return
	}
	for i := range g.Windows {
		if g.Windows[i] == window {
			*g.Windows[1] = *g.Windows[0]
			g.Windows[0] = window
			break
		}
	}
}

// 0..3: corners (Lower-right, Lower-left, Unused, Unused)
func GetWindowResizeCornerID(window *ImGuiWindow, n int) ImGuiID {
	IM_ASSERT(n >= 0 && n < 4)
	var id = window.ID
	id = ImHashStr("#RESIZE", 0, id)
	id = ImHashData(unsafe.Pointer(&n), unsafe.Sizeof(n), id)
	return id
}

// Borders (Left, Right, Up, Down)
func GetWindowResizeBorderID(window *ImGuiWindow, dir ImGuiDir) ImGuiID {
	IM_ASSERT(dir >= 0 && dir < 4)
	var n = (int)(dir) + 4
	var id = window.ID
	id = ImHashStr("#RESIZE", 0, id)
	id = ImHashData(unsafe.Pointer(&n), unsafe.Sizeof(n), id)
	return id
}
