package imgui

// access the IO structure (mouse/keyboard/gamepad inputs, time, various configuration options/flags)
func GetIO() *ImGuiIO {
	IM_ASSERT_USER_ERROR(GImGui != nil, "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?")
	return &GImGui.IO
}

// Demo, Debug, Information
func ShowDemoWindow(p_open *bool)         { panic("not implemented") } // create Demo window. demonstrate most ImGui features. call this to learn about the library! try to make it always available in your application!
func ShowAboutWindow(p_open *bool)        { panic("not implemented") } // create About window. display Dear ImGui version, credits and build/system information.
func ShowStyleEditor(ref *ImGuiStyle)     { panic("not implemented") } // add style editor block (not a window). you can pass in a reference ImGuiStyle structure to compare to, revert to and save to (else it uses the default style)
func ShowStyleSelector(label string) bool { panic("not implemented") } // add style selector block (not a window), essentially a combo listing the default styles.
func ShowFontSelector(label string)       { panic("not implemented") } // add font selector block (not a window), essentially a combo listing the loaded fonts.
func ShowUserGuide()                      { panic("not implemented") } // add basic help/info block (not a window): how to manipulate ImGui as a end-user (mouse/keyboard controls).

// get the compiled version string e.g. "1.80 WIP" (essentially the value for IMGUI_VERSION from the compiled version of imgui.cpp)
func GetVersion() string {
	return IMGUI_VERSION
}

// Widgets: Main
// - Most widgets return true when the value has been changed or when pressed/selected
// - You may also use one of the many IsItemXXX functions (e.g. IsItemActive, IsItemHovered, etc.) to query widget state.

// Widgets: Input with Keyboard
// - If you want to use InputText() with std::string or any custom dynamic string type, see misc/cpp/imgui_stdlib.h and comments in imgui_demo.cpp.
// - Most of the ImGuiInputTextFlags flags are only useful for InputText() and not for InputFloatX, InputIntX, Inputetc double.

// Widgets: Color Editor/Picker (tip: the ColorEdit* functions have a little color square that can be left-clicked to open a picker, and right-clicked to open an option menu.)
// - Note that in C++ a 'v float[X]' function argument is the _same_ as 'float* v', the array syntax is just a way to document the number of elements that are expected to be accessible.
// - You can pass the address of a first element float out of a contiguous structure, e.g. &myvector.x
func ColorEdit3(label string, col [3]float, flsgs ImGuiColorEditFlags) bool { panic("not implemented") }
func ColorEdit4(label string, col [4]float, flsgs ImGuiColorEditFlags) bool { panic("not implemented") }
func ColorPicker3(label string, col [3]float, flsgs ImGuiColorEditFlags) bool {
	panic("not implemented")
}
func ColorPicker4(label string, col [4]float, flsgs ImGuiColorEditFlags, ref_col []float) bool {
	panic("not implemented")
}
func ColorButton(desc_id string, col ImVec4, flsgs ImGuiColorEditFlags, size ImVec2 /*= 0*/) bool {
	panic("not implemented")
}                                                   // display a color square/button, hover for details, return true when pressed.
func SetColorEditOptions(flags ImGuiColorEditFlags) { panic("not implemented") } // initialize current options (generally on application startup) if you want to select a default format, picker type, etc. User will be able to change many settings, unless you pass the _NoOptions flag to your calls.

func SelectablePointer(label string, p_selected *bool, flsgs ImGuiSelectableFlags, size ImVec2) bool {
	panic("not implemented")
} // "bool* p_selected" poto int the selection state (read-write), as a convenient helper.

// Widgets: Data Plotting
// - Consider using ImPlot (https://github.com/epezent/implot) which is much better!
func PlotLines(label string, values []float, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/, stride int /*= sizeof(float)*/) {
	panic("not implemented")
}
func PlotLinesFunc(label string, values_getter func(data interface{}, idx int) float, data interface{}, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/) {
	panic("not implemented")
}
func PlotHistogram(label string, values []float, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/, stride int /* = sizeof(float)*/) {
	panic("not implemented")
}
func PlotHistogramFunc(label string, values_getter func(data interface{}, idx int) float, data interface{}, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/) {
	panic("not implemented")
}

// Widgets: Value() Helpers.
// - Those are merely shortcut to calling Text() with a format string. Output single value in "name: value" format (tip: freely declare more in your code to handle your types. you can add functions to the ImGui namespace)
func ValueBool(prefix string, b bool)                        { panic("not implemented") }
func ValueInt(prefix string, v int)                          { panic("not implemented") }
func ValueUint(prefix string, uv uint)                       { panic("not implemented") }
func ValueFloat(prefix string, v float, float_format string) { panic("not implemented") }

// Widgets: Menus
// - Use BeginMenuBar() on a window ImGuiWindowFlags_MenuBar to append to its menu bar.
// - Use BeginMainMenuBar() to create a menu bar at the top of the screen and append to it.
// - Use BeginMenu() to create a menu. You can call BeginMenu() multiple time with the same identifier to append more items to it.
// - Not that MenuItem() keyboardshortcuts are displayed as a convenience but _not processed_ by Dear ImGui at the moment.
func BeginMenuBar() bool                                   { panic("not implemented") } // append to menu-bar of current window (requires ImGuiWindowFlags_MenuBar flag set on parent window).
func EndMenuBar()                                          { panic("not implemented") } // only call EndMenuBar() if BeginMenuBar() returns true!
func BeginMainMenuBar() bool                               { panic("not implemented") } // create and append to a full screen menu-bar.
func EndMainMenuBar()                                      { panic("not implemented") } // only call EndMainMenuBar() if BeginMainMenuBar() returns true!
func BeginMenu(label string, enabled bool /*= true*/) bool { panic("not implemented") } // create a sub-menu entry. only call EndMenu() if this returns true!
func EndMenu()                                             { panic("not implemented") } // only call EndMenu() if BeginMenu() returns true!
func MenuItem(label string, shortcut string /*= L*/, selected bool /*= e*/, enabled bool /*= true*/) bool {
	panic("not implemented")
} // return true when activated.
func MenuItemSelected(label string, shortcut string, p_selected *bool, enabled bool /*= true*/) bool {
	panic("not implemented")
} // return true when activated + toggle (*p_selected) if p_selected != NULL

// Tab Bars, Tabs
func BeginTabBar(str_id string, flsgs ImGuiTabBarFlags) bool                { panic("not implemented") } // create and append into a TabBar
func EndTabBar()                                                            { panic("not implemented") } // only call EndTabBar() if BeginTabBar() returns true!
func BeginTabItem(label string, p_open *bool, flsgs ImGuiTabItemFlags) bool { panic("not implemented") } // create a Tab. Returns true if the Tab is selected.
func EndTabItem()                                                           { panic("not implemented") } // only call EndTabItem() if BeginTabItem() returns true!
func TabItemButton(label string, flsgs ImGuiTabItemFlags) bool              { panic("not implemented") } // create a Tab behaving like a button. return true when clicked. cannot be selected in the tab bar.
func SetTabItemClosed(tab_or_docked_window_label string)                    { panic("not implemented") } // notify TabBar or Docking system of a closed tab/window ahead (useful to reduce visual flicker on reorderable tab bars). For tab-bar: call after BeginTabBar() and before Tab submissions. Otherwise call with a window name.

// Clipping
// - Mouse hovering is affected by ImGui::PushClipRect() calls, unlike direct calls to ImDrawList::PushClipRect() which are render only.
// Push a clipping rectangle for both ImGui logic (hit-testing etc.) and low-level ImDrawList rendering.
// - When using this function it is sane to ensure that float are perfectly rounded to integer values,
//   so that e.g. (int)(max.x-min.x) in user's render produce correct result.
// - If the code here changes, may need to update code of functions like NextColumn() and PushColumnClipRect():
//   some frequently called functions which to modify both channels and clipping simultaneously tend to use the
//   more specialized SetWindowClipRectBeforeSetChannel() to avoid extraneous updates of underlying ImDrawCmds.
func PushClipRect(cr_min ImVec2, cr_max ImVec2, intersect_with_current_clip_rect bool) {
	var window *ImGuiWindow = GetCurrentWindow()
	window.DrawList.PushClipRect(cr_min, cr_max, intersect_with_current_clip_rect)
	window.ClipRect = ImRectFromVec4(&window.DrawList._ClipRectStack[len(window.DrawList._ClipRectStack)-1])
}

func PopClipRect() {
	var window = GetCurrentWindow()
	window.DrawList.PopClipRect()
	window.ClipRect = ImRectFromVec4(&window.DrawList._ClipRectStack[len(window.DrawList._ClipRectStack)-1])
}

// Viewports
// - Currently represents the Platform Window created by the application which is hosting our Dear ImGui windows.
// - In 'docking' branch with multi-viewport enabled, we extend this concept to have multiple active viewports.
// - In the future we will extend this concept further to also represent Platform Monitor and support a "no main platform window" operation mode.
func GetMainViewport() *ImGuiViewport {
	var g = GImGui
	return g.Viewports[0]
} // return primary/default viewport. This can never be NULL.
