package imgui

type SizeCallback func()

const (
	WindowFlagsNone                      = 0
	WindowFlagsNoTitleBar                = 1 << 0  // Disable title-bar
	WindowFlagsNoResize                  = 1 << 1  // Disable user resizing with the lower-right grip
	WindowFlagsNoMove                    = 1 << 2  // Disable user moving the window
	WindowFlagsNoScrollbar               = 1 << 3  // Disable scrollbars (window can still scroll with mouse or programmatically)
	WindowFlagsNoScrollWithMouse         = 1 << 4  // Disable user vertically scrolling with mouse wheel. On child window, mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	WindowFlagsNoCollapse                = 1 << 5  // Disable user collapsing window by double-clicking on it
	WindowFlagsAlwaysAutoResize          = 1 << 6  // Resize every window to its content every frame
	WindowFlagsNoBackground              = 1 << 7  // Disable drawing background color (WindowBg, etc.) and outside border. Similar as using SetNextWindowBgAlpha(0.0f).
	WindowFlagsNoSavedSettings           = 1 << 8  // Never load/save settings in .ini file
	WindowFlagsNoMouseInputs             = 1 << 9  // Disable catching mouse, hovering test with pass through.
	WindowFlagsMenuBar                   = 1 << 10 // Has a menu-bar
	WindowFlagsHorizontalScrollbar       = 1 << 11 // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width,0.0f)); prior to calling Begin() to specify width. Read code in imguidemo in the "Horizontal Scrolling" section.
	WindowFlagsNoFocusOnAppearing        = 1 << 12 // Disable taking focus when transitioning from hidden to visible state
	WindowFlagsNoBringToFrontOnFocus     = 1 << 13 // Disable bringing window to front when taking focus (e.g. clicking on it or programmatically giving it focus)
	WindowFlagsAlwaysVerticalScrollbar   = 1 << 14 // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	WindowFlagsAlwaysHorizontalScrollbar = 1 << 15 // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	WindowFlagsAlwaysUseWindowPadding    = 1 << 16 // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows, because more convenient)
	WindowFlagsNoNavInputs               = 1 << 18 // No gamepad/keyboard navigation within the window
	WindowFlagsNoNavFocus                = 1 << 19 // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	WindowFlagsUnsavedDocument           = 1 << 20 // Display a dot next to the title. When used in a tab/docking context, tab is selected when clicking the X + closure is not assumed (will wait for user to stop submitting the tab). Otherwise closure is assumed when pressing the X, so if you keep submitting the tab may reappear at end of tab bar.
	WindowFlagsNoNav                     = WindowFlagsNoNavInputs | WindowFlagsNoNavFocus
	WindowFlagsNoDecoration              = WindowFlagsNoTitleBar | WindowFlagsNoResize | WindowFlagsNoScrollbar | WindowFlagsNoCollapse
	WindowFlagsNoInputs                  = WindowFlagsNoMouseInputs | WindowFlagsNoNavInputs | WindowFlagsNoNavFocus

	// [Internal]
	WindowFlagsNavFlattened = 1 << 23 // [BETA] Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	WindowFlagsChildWindow  = 1 << 24 // Don't use! For internal use by BeginChild()
	WindowFlagsTooltip      = 1 << 25 // Don't use! For internal use by BeginTooltip()
	WindowFlagsPopup        = 1 << 26 // Don't use! For internal use by BeginPopup()
	WindowFlagsModal        = 1 << 27 // Don't use! For internal use by BeginPopupModal()
	WindowFlagsChildMenu    = 1 << 28 // Don't use! For internal use by BeginMenu()
)

func Begin(name string, open bool, flags WindowFlags) bool {
	panic("not implemented")
	return false
}

func End() { panic("not implemented") }

func BeginChild(id string, size Vec2, border bool, flags WindowFlags) bool {
	panic("not implemented")
	return false
}
func EndChild() {}

func IsWindowAppearing() bool {
	panic("not implemented")
	return false
}

func IsWindowCollapsed() bool {
	panic("not implemented")
	return false
}

func IsWindowFocused(flags FocusedFlags) bool {
	panic("not implemented")
	return false
}

func IsWindowHovered(flags HoveredFlags) bool {
	panic("not implemented")
	return false
}

func GetWindowDrawList() DrawList {
	panic("not implemented")
	return DrawList{}
}

func GetWindowPos() Vec2 {
	panic("not implemented")
	return Vec2{}
}

func GetWindowSize() Vec2 {
	panic("not implemented")
	return Vec2{}
}

func GetWindowWidth() float32 {
	panic("not implemented")
	return 0
}

func GetWindowHeight() float32 {
	panic("not implemented")
	return 0
}

func SetNextWindowPos(pos Vec2, cond Cond) { panic("not implemented") }

func SetNextWindowSize(size Vec2, cond Cond) { panic("not implemented") }

func SetNextWindowSizeConstraints(sizeMin, sizeMax Vec2, callback SizeCallback) {
	panic("not implemented")
}

func SetNextWindowContentSize(size Vec2) { panic("not implemented") }

func SetNextWindowCollapsed(collapsed bool, cond Cond) { panic("not implemented") }

func SetNextWindowFocus() { panic("not implemented") }

func SetNextWindowBgAlpha(alpha float32) { panic("not implemented") }

func SetWindowPos(pos Vec2, cond Cond) { panic("not implemented") }

func SetWindowSize(size Vec2, cond Cond) { panic("not implemented") }

func SetWindowCollapsed(collapsed bool, cond Cond) { panic("not implemented") }

func SetWindowFocus() { panic("not implemented") }

func SetWindowFontScale(scale float32) { panic("not implemented") }

func GetContentRegionAvail() Vec2 {
	panic("not implemented")
	return Vec2{}
}

func GetContentRegionMax() Vec2 {
	panic("not implemented")
	return Vec2{}
}

func GetWindowContentRegionMin() Vec2 {
	panic("not implemented")
	return Vec2{}
}

func GetWindowContentRegionMax() Vec2 {
	panic("not implemented")
	return Vec2{}
}
