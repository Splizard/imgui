package imgui

// Windows: Display Order and Focus Order
func FocusWindow(window *ImGuiWindow) {
	var g = GImGui

	if g.NavWindow != window {
		g.NavWindow = window
		if window != nil && g.NavDisableMouseHover {
			g.NavMousePosDirty = true
		}
		if window != nil {
			g.NavId = window.NavLastIds[0]
		} else {
			g.NavId = 0
		}
		g.NavFocusScopeId = 0
		g.NavIdIsAlive = false
		g.NavLayer = ImGuiNavLayer_Main
		g.NavInitRequest = false
		g.NavMoveSubmitted = false
		g.NavMoveScoringItems = false
		NavUpdateAnyRequestFlag()
		//IMGUI_DEBUG_LOG("FocusWindow(\"%s\")\n", window ? window.Name : nil);
	}

	// Close popups if any
	ClosePopupsOverWindow(window, false)

	// Move the root window to the top of the pile
	IM_ASSERT(window == nil || window.RootWindow != nil)

	var focus_front_window *ImGuiWindow
	var display_front_window *ImGuiWindow
	if window != nil {
		focus_front_window = window
		display_front_window = window
	}

	// Steal active widgets. Some of the cases it triggers includes:
	// - Focus a window while an InputText in another window is active, if focus happens before the old InputText can run.
	// - When using Nav to activate menu items (due to timing of activating on press.new window appears.losing ActiveId)
	if g.ActiveId != 0 && g.ActiveIdWindow != nil && g.ActiveIdWindow.RootWindow != focus_front_window {
		if !g.ActiveIdNoClearOnFocusLoss {
			ClearActiveID()
		}
	}

	// Passing nil allow to disable keyboard focus
	if window == nil {
		return
	}

	// Bring to front
	BringWindowToFocusFront(focus_front_window)
	if ((window.Flags | display_front_window.Flags) & ImGuiWindowFlags_NoBringToFrontOnFocus) == 0 {
		BringWindowToDisplayFront(display_front_window)
	}
}

func BringWindowToFocusFront(window *ImGuiWindow) {
	var g = GImGui
	IM_ASSERT(window == window.RootWindow)

	var cur_order = window.FocusOrder
	IM_ASSERT(g.WindowsFocusOrder[cur_order] == window)
	if g.WindowsFocusOrder[len(g.WindowsFocusOrder)-1] == window {
		return
	}

	var new_order int = int(len(g.WindowsFocusOrder)) - 1
	for n := int(cur_order); n < new_order; n++ {
		g.WindowsFocusOrder[n] = g.WindowsFocusOrder[n+1]
		g.WindowsFocusOrder[n].FocusOrder--
		IM_ASSERT(int(g.WindowsFocusOrder[n].FocusOrder) == n)
	}
	g.WindowsFocusOrder[new_order] = window
	window.FocusOrder = (short)(new_order)
}

func BringWindowToDisplayFront(window *ImGuiWindow) {
	var g = GImGui
	var current_front_window *ImGuiWindow = g.Windows[len(g.Windows)-1]
	if current_front_window == window || current_front_window.RootWindow == window { // Cheap early out (could be better)
		return
	}
	for i := len(g.Windows) - 2; i >= 0; i-- { // We can ignore the top-most window
		if g.Windows[i] == window {
			copy(g.Windows[i:], g.Windows[i+1:i+1+(len(g.Windows)-i-len(g.Windows)-i-11)])
			g.Windows[len(g.Windows)-1] = window
			break
		}
	}
}
