package imgui

// Window manipulation
// - Prefer using SetNextXXX functions (before Begin) rather that SetXXX functions (after Begin).

// set next window position. call before Begin(). use pivot=(0.5,0.5) to center on given point, etc.
func SetNextWindowPos(pos *ImVec2, cond ImGuiCond, pivot ImVec2) {
	var g = GImGui
	IM_ASSERT(cond == 0 || ImIsPowerOfTwoInt(int(cond))) // Make sure the user doesn't attempt to combine multiple condition flags.
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasPos
	g.NextWindowData.PosVal = *pos
	g.NextWindowData.PosPivotVal = pivot
	if cond != 0 {
		g.NextWindowData.PosCond = cond
	} else {
		g.NextWindowData.PosCond = ImGuiCond_Always
	}
}

// set next window size limits. use -1,-1 on either X/Y axis to preserve the current size. Sizes will be rounded down. Use callback to apply non-trivial programmatic constraints.
func SetNextWindowSizeConstraints(size_min ImVec2, size_max ImVec2, custom_callback ImGuiSizeCallback, custom_callback_data any) {
	var g = GImGui
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasSizeConstraint
	g.NextWindowData.SizeConstraintRect = ImRect{size_min, size_max}
	g.NextWindowData.SizeCallback = custom_callback
	g.NextWindowData.SizeCallbackUserData = custom_callback_data
}

// Content size = inner scrollable rectangle, padded with WindowPadding.
// SetNextWindowContentSize(ImVec2(100,100) + ImGuiWindowFlags_AlwaysAutoResize will always allow submitting a 100x100 item.
// set next window content size (~ scrollable client area, which enforce the range of scrollbars). Not including window decorations (title bar, menu bar, etc.) nor WindowPadding. set an axis to 0.0 to leave it automatic. call before Begin()
func SetNextWindowContentSize(size ImVec2) {
	var g = GImGui
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasContentSize
	g.NextWindowData.ContentSizeVal = *ImFloorVec(&size)
}

// set next window collapsed state. call before Begin()
func SetNextWindowCollapsed(collapsed bool, cond ImGuiCond) {
	var g = GImGui
	IM_ASSERT(cond == 0 || ImIsPowerOfTwoInt(int(cond))) // Make sure the user doesn't attempt to combine multiple condition flags.
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasCollapsed
	g.NextWindowData.CollapsedVal = collapsed
	if cond != 0 {
		g.NextWindowData.CollapsedCond = cond
	} else {
		g.NextWindowData.CollapsedCond = ImGuiCond_Always
	}
}

// set next window to be focused / top-most. call before Begin()
func SetNextWindowFocus() {
	var g = GImGui
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasFocus
}

// (not recommended) set current window position - call within Begin()/End(). prefer using SetNextWindowPos(), as this may incur tearing and side-effects.
func SetWindowPos(pos ImVec2, cond ImGuiCond) {
	var window = GetCurrentWindowRead()
	setWindowPos(window, &pos, cond)
}

// (not recommended) set current window size - call within Begin()/End(). set to ImVec2(0, 0) to force an auto-fit. prefer using SetNextWindowSize(), as this may incur tearing and minor side-effects.
func SetWindowSize(size ImVec2, cond ImGuiCond) {
	setWindowSize(GImGui.CurrentWindow, &size, cond)
}

func setWindowCollapsed(window *ImGuiWindow, collapsed bool, cond ImGuiCond) {
	// Test condition (NB: bit 0 is always true) and clear flags for next time
	if cond != 0 && (window.SetWindowCollapsedAllowFlags&cond) == 0 {
		return
	}
	window.SetWindowCollapsedAllowFlags &= ^(ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing)

	// Set
	window.Collapsed = collapsed
}

// (not recommended) set current window collapsed state. prefer using SetNextWindowCollapsed().
func SetWindowCollapsed(collapsed bool, cond ImGuiCond) {
	setWindowCollapsed(GImGui.CurrentWindow, collapsed, cond)
}

// (not recommended) set current window to be focused / top-most. prefer using SetNextWindowFocus().
func SetWindowFocus() {
	FocusWindow(GImGui.CurrentWindow)
}

// [OBSOLETE] set font scale. Adjust IO.FontGlobalScale if you want to scale all windows. This is an old API! For correct scaling, prefer to reload font + rebuild ImFontAtlas + call style.ScaleAllSizes().
func SetWindowFontScale(scale float) {
	IM_ASSERT(scale > 0.0)
	var g = GImGui
	var window = GetCurrentWindow()
	window.FontWindowScale = scale
	calculated := window.CalcFontSize()
	g.FontSize = calculated
	g.DrawListSharedData.FontSize = calculated
}

// set named window position.
func SetNamedWindowPos(name string, pos ImVec2, cond ImGuiCond) {
	if window := FindWindowByName(name); window != nil {
		setWindowPos(window, &pos, cond)
	}
}

// set named window size. set axis to 0.0 to force an auto-fit on this axis.
func SetNamedWindowSize(name string, size ImVec2, cond ImGuiCond) {
	if window := FindWindowByName(name); window != nil {
		setWindowSize(window, &size, cond)
	}
}

// set named window collapsed state
func SetNamedWindowCollapsed(name string, collapsed bool, cond ImGuiCond) {
	if window := FindWindowByName(name); window != nil {
		setWindowCollapsed(window, collapsed, cond)
	}
}

// set named window to be focused / top-most. use NULL to remove focus.
func SetNamedWindowFocus(name string) {
	if name != "" {
		if window := FindWindowByName(name); window != nil {
			FocusWindow(window)
		}
	} else {
		FocusWindow(nil)
	}
}
