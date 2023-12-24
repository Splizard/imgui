package imgui

func UpdateTabFocus() {
	g := guiContext

	// Pressing TAB activate widget focus
	g.TabFocusPressed = g.NavWindow != nil && g.NavWindow.Active && (g.NavWindow.Flags&ImGuiWindowFlags_NoNavInputs == 0) && !g.IO.KeyCtrl && IsKeyPressedMap(ImGuiKey_Tab, true)
	if g.ActiveId == 0 && g.TabFocusPressed {
		// - This path is only taken when no widget are active/tabbed-into yet.
		//   Subsequent tabbing will be processed by FocusableItemRegister()
		// - Note that SetKeyboardFocusHere() sets the Next fields mid-frame. To be consistent we also
		//   manipulate the Next fields here even though they will be turned into Curr fields below.
		g.TabFocusRequestNextWindow = g.NavWindow
		g.TabFocusRequestNextCounterRegular = INT_MAX

		var shift int
		if g.IO.KeyShift {
			shift = -1
		}

		if g.NavId != 0 && g.NavIdTabCounter != INT_MAX {
			g.TabFocusRequestNextCounterTabStop = g.NavIdTabCounter + shift
		} else {
			g.TabFocusRequestNextCounterTabStop = shift
		}
	}

	// Turn queued focus request into current one
	g.TabFocusRequestCurrWindow = nil
	g.TabFocusRequestCurrCounterRegular = INT_MAX
	g.TabFocusRequestCurrCounterTabStop = INT_MAX
	if g.TabFocusRequestNextWindow != nil {
		var window = g.TabFocusRequestNextWindow
		g.TabFocusRequestCurrWindow = window
		if g.TabFocusRequestNextCounterRegular != INT_MAX && window.DC.FocusCounterRegular != -1 {
			g.TabFocusRequestCurrCounterRegular = ImModPositive(g.TabFocusRequestNextCounterRegular, window.DC.FocusCounterRegular+1)
		}
		if g.TabFocusRequestNextCounterTabStop != INT_MAX && window.DC.FocusCounterTabStop != -1 {
			g.TabFocusRequestCurrCounterTabStop = ImModPositive(g.TabFocusRequestNextCounterTabStop, window.DC.FocusCounterTabStop+1)
		}
		g.TabFocusRequestNextWindow = nil
		g.TabFocusRequestNextCounterRegular = INT_MAX
		g.TabFocusRequestNextCounterTabStop = INT_MAX
	}

	g.NavIdTabCounter = INT_MAX
}
