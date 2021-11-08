package imgui

// Widgets: Selectables
// - A selectable highlights when hovered, and can display another color when selected.
// - Neighbors selectable extend their highlight bounds in order to leave no gap between them. This is so a series of selected Selectable appear contiguous.
// "selected bool" carry the selection state (read-only). Selectable() is clicked is returns true so you can modify your selection state. size.x==0.0: use remaining width, size.x>0.0: specify width. size.y==0.0: use label height, size.y>0.0: specify height
func Selectable(label string, selected bool, flsgs ImGuiSelectableFlags, size ImVec2) bool {
	panic("not implemented")
}

// Focus, Activation
// - Prefer using "SetItemDefaultFocus()" over "if (IsWindowAppearing()) SetScrollHereY()" when applicable to signify "this is the default item"

// make last item the default focused item of a window.
func SetItemDefaultFocus() {
	var g = GImGui
	var window = g.CurrentWindow
	if !window.Appearing {
		return
	}
	if g.NavWindow == window.RootWindowForNav && (g.NavInitRequest || g.NavInitResultId != 0) && g.NavLayer == window.DC.NavLayerCurrent {
		g.NavInitRequest = false
		g.NavInitResultId = g.LastItemData.ID
		g.NavInitResultRectRel = ImRect{g.LastItemData.Rect.Min.Sub(window.Pos), g.LastItemData.Rect.Max.Sub(window.Pos)}
		NavUpdateAnyRequestFlag()
		if !IsItemVisible() {
			SetScrollHereY(0.5)
		}
	}
}

// focus keyboard on the next widget. Use positive 'offset' to access sub components of a multiple component widget. Use -1 to access previous widget.
func SetKeyboardFocusHere(offset int) {
	IM_ASSERT(offset >= -1) // -1 is allowed but not below
	var g = GImGui
	var window *ImGuiWindow = g.CurrentWindow
	g.TabFocusRequestNextWindow = window
	g.TabFocusRequestNextCounterRegular = window.DC.FocusCounterRegular + 1 + offset
	g.TabFocusRequestNextCounterTabStop = INT_MAX
}

// Focus Scope (WIP)
// This is generally used to identify a selection set (multiple of which may be in the same window), as selection
// patterns generally need to react (e.g. clear selection) when landing on an item of the set.
func PushFocusScope(id ImGuiID) {
	var g = GImGui
	var window = g.CurrentWindow
	g.FocusScopeStack = append(g.FocusScopeStack, window.DC.NavFocusScopeIdCurrent)
	window.DC.NavFocusScopeIdCurrent = id
}

func PopFocusScope() {
	var g = GImGui
	var window = g.CurrentWindow
	IM_ASSERT(len(g.FocusScopeStack) > 0) // Too many PopFocusScope() ?
	window.DC.NavFocusScopeIdCurrent = g.FocusScopeStack[len(g.FocusScopeStack)-1]
	g.FocusScopeStack = g.FocusScopeStack[:len(g.FocusScopeStack)-1]
}

func GetFocusedFocusScope() ImGuiID { var g *ImGuiContext = GImGui; return g.NavFocusScopeId } // Focus scope which is actually active
func GetFocusScope() ImGuiID {
	var g *ImGuiContext = GImGui
	return g.CurrentWindow.DC.NavFocusScopeIdCurrent
} // Focus scope we are outputting into, set by PushFocusScope()

// Notifies Dear ImGui when hosting platform windows lose or gain input focus
func (io *ImGuiIO) AddFocusEvent(focused bool) {
	if focused {
		return
	}

	// Clear buttons state when focus is lost
	// (this is useful so e.g. releasing Alt after focus loss on Alt-Tab doesn't trigger the Alt menu toggle)
	io.KeysDown = [len(io.KeysDown)]bool{}
	for n := range io.KeysDownDuration {
		io.KeysDownDuration[n] = -1
		io.KeysDownDurationPrev[n] = -1
	}
	io.KeyCtrl = false
	io.KeyShift = false
	io.KeyAlt = false
	io.KeySuper = false
	io.KeyMods = ImGuiKeyModFlags_None
	io.KeyModsPrev = ImGuiKeyModFlags_None
	for n := range io.NavInputsDownDuration {
		io.NavInputsDownDuration[n] = -1
		io.NavInputsDownDurationPrev[n] = -1
	}
}

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
			amount := len(g.Windows) - i - 1
			copy(g.Windows[i:], g.Windows[i+1:i+1+amount])
			g.Windows[len(g.Windows)-1] = window
			break
		}
	}
}

func FocusTopMostWindowUnderOne(under_this_window *ImGuiWindow, ignore_window *ImGuiWindow) {
	var g = GImGui

	var start_idx int
	if under_this_window != nil {
		start_idx = FindWindowFocusIndex(under_this_window)
	} else {
		start_idx = int(len(g.WindowsFocusOrder) - 1)
	}

	for i := start_idx; i >= 0; i-- {
		// We may later decide to test for different NoXXXInputs based on the active navigation input (mouse vs nav) but that may feel more confusing to the user.
		var window *ImGuiWindow = g.WindowsFocusOrder[i]
		IM_ASSERT(window == window.RootWindow)
		if window != ignore_window && window.WasActive {
			if (window.Flags & (ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs)) != (ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs) {
				var focus_window *ImGuiWindow = NavRestoreLastChildNavWindow(window)
				FocusWindow(focus_window)
				return
			}
		}
	}
	FocusWindow(nil)
}
