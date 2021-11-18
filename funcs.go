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
func InvisibleButton(str_id string, size ImVec2, flsgs ImGuiButtonFlags) bool {
	panic("not implemented")
}                                                  // flexible button behavior without the visuals, frequently useful to build custom behaviors using the public api (along with IsItemActive, IsItemHovered, etc.)
func ArrowButton(str_id string, dir ImGuiDir) bool { panic("not implemented") } // square button with an arrow shape

func RadioButtonBool(label string, active bool)              { panic("not implemented") } // use with e.g. if (RadioButton("one", my_value==1)) { my_value = 1 bool {panic("not implemented")} }
func RadioButtonInt(label string, v *int, v_button int) bool { panic("not implemented") } // shortcut to handle the above pattern when value is an integer
func ProgressBar(fraction float, size_arg ImVec2 /*= ImVec2(-FLT_MIN, 0)*/, overlay string) {
	panic("not implemented")
}
func Bullet() { panic("not implemented") } // draw a small circle + keep the cursor on the same line. advance cursor x position by GetTreeNodeToLabelSpacing(), same distance that TreeNode() uses

// Widgets: Regular Sliders
// - CTRL+Click on any slider to turn them into an input box. Manually input values aren't clamped and can go off-bounds.
// - Adjust format string to decorate the value with a prefix, a suffix, or adapt the editing and display precision e.g. "%.3f" -> 1.234; "%5.2 secs" -> 01.23 secs; "Biscuit: %.0f" -> Biscuit: 1; etc.
// - Format string may also be set to NULL or use the default format ("%f" or "%d").
// - Legacy: Pre-1.78 there are SliderXXX() function signatures that takes a final `power float=1.0' argument instead of the `ImGuiSliderFlags flags=0' argument.
//   If you get a warning converting a to float ImGuiSliderFlags, read https://github.com/ocornut/imgui/issues/3361
func SliderFloat(label string, v *float, v_min float, v_max float, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
} // adjust format to decorate the value with a prefix or a suffix for in-slider labels or unit display.
func SliderFloat2(label string, v [2]float, v_min float, v_max float, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderFloat3(label string, v [3]float, v_min float, v_max float, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderFloat4(label string, v [4]float, v_min float, v_max float, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderAngle(label string, v_rad *float, v_degrees_min float /*= 0*/, v_degrees_max float /*= 0*/, format string /* = "%.0f deg"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderInt(label string, v *int, v_min int, v_max int, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderInt2(label string, v [2]int, v_min int, v_max int, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderInt3(label string, v [3]int, v_min int, v_max int, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderInt4(label string, v [4]int, v_min int, v_max int, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderScalar(label string, data_type ImGuiDataType, p_data interface{}, p_min interface{}, p_max interface{}, format string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderScalarN(label string, data_type ImGuiDataType, p_data interface{}, components int, p_min interface{}, p_max interface{}, format string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func VSliderFloat(label string, size ImVec2, v *float, v_min float, v_max float, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func VSliderInt(label string, size ImVec2, v *int, v_min int, v_max int, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func VSliderScalar(label string, size ImVec2, data_type ImGuiDataType, p_data interface{}, p_min interface{}, p_max interface{}, format string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}

// Widgets: Input with Keyboard
// - If you want to use InputText() with std::string or any custom dynamic string type, see misc/cpp/imgui_stdlib.h and comments in imgui_demo.cpp.
// - Most of the ImGuiInputTextFlags flags are only useful for InputText() and not for InputFloatX, InputIntX, Inputetc double.
func InputText(label string, char []byte, buf_size uintptr, flsgs ImGuiInputTextFlags, callback ImGuiInputTextCallback /*= L*/, user_data interface{}) bool {
	panic("not implemented")
}
func InputTextMultiline(label string, buf []byte, buf_size uintptr, size ImVec2, flsgs ImGuiInputTextFlags, callback ImGuiInputTextCallback /*= L*/, user_data interface{}) bool {
	panic("not implemented")
}
func InputTextWithHint(label string, hstring int, char []byte, buf_size uintptr, flsgs ImGuiInputTextFlags, callback ImGuiInputTextCallback /*= L*/, user_data interface{}) bool {
	panic("not implemented")
}
func InputFloat(label string, v *float, step float /*= 0*/, step_fast float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputFloat2(label string, v [2]float, format string /*= "%.3f"*/, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputFloat3(label string, v [3]float, format string /*= "%.3f"*/, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputFloat4(label string, v [4]float, format string /*= "%.3f"*/, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputInt(label string, v *int, step int /*= 1*/, step_fast int /*= 100*/, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputInt2(label string, v [2]int, flsgs ImGuiInputTextFlags) bool { panic("not implemented") }
func InputInt3(label string, v [3]int, flsgs ImGuiInputTextFlags) bool { panic("not implemented") }
func InputInt4(label string, v [4]int, flsgs ImGuiInputTextFlags) bool { panic("not implemented") }
func InputDouble(label string, v *double, step double /*= 0*/, step_fast double /*= 0*/, format string /*= "%.6f"*/, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputScalar(label string, data_type ImGuiDataType, p_data interface{}, p_step interface{} /*= L*/, p_step_fast interface{} /*= L*/, format string, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func InputScalarN(label string, data_type ImGuiDataType, p_data interface{}, components int, p_step interface{} /*= L*/, p_step_fast interface{} /*= L*/, format string, flsgs ImGuiInputTextFlags) bool {
	panic("not implemented")
}

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

// Widgets: List Boxes
// - This is essentially a thin wrapper to using BeginChild/EndChild with some stylistic changes.
// - The BeginListBox()/EndListBox() api allows you to manage your contents and selection state however you want it, by creating e.g. Selectable() or any items.
// - The simplified/old ListBox() api are helpers over BeginListBox()/EndListBox() which are kept available for convenience purpose. This is analoguous to how Combos are created.
// - Choose frame width:   size.x > 0.0: custom  /  size.x < 0.0 or -FLT_MIN: right-align   /  size.x = 0.0 (default): use current ItemWidth
// - Choose frame height:  size.y > 0.0: custom  /  size.y < 0.0 or -FLT_MIN: bottom-align  /  size.y = 0.0 (default): arbitrary default height which can fit ~7 items
func BeginListBox(label string, size ImVec2) bool { panic("not implemented") } // open a framed scrolling region
func EndListBox()                                 { panic("not implemented") } // only call EndListBox() if BeginListBox() returned true!
func ListBox(label string, current_item *int, items []string, items_count int, height_in_items int /*= -1*/) bool {
	panic("not implemented")
}
func ListBoxFunc(label string, current_item *int, items_getter func(data interface{}, idx int, out_text *string) bool, data interface{}, items_count int, height_in_items int /*= -1*/) bool {
	panic("not implemented")
}

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

// Legacy Columns API (prefer using Tables!)
// - You can also use SameLine(pos_x) to mimic simplified columns.
func Columns(count int /*= 1*/, id string /*= L*/, border bool /*= true*/) { panic("not implemented") }
func NextColumn()                                                          { panic("not implemented") } // next column, defaults to current row or next row if the current row is finished
func GetColumnIndex() int                                                  { panic("not implemented") } // get current column index
func GetColumnWidth(column_index int /*= -1*/) float                       { panic("not implemented") } // get column width (in pixels). pass -1 to use current column
func SetColumnWidth(column_index int, width float)                         { panic("not implemented") } // set column width (in pixels). pass -1 to use current column
func GetColumnOffset(column_index int /*= -1*/) float                      { panic("not implemented") } // get position of column line (in pixels, from the left side of the contents region). pass -1 to use current column, otherwise 0..GetColumnsCount() inclusive. column 0 is typically 0.0
func SetColumnOffset(column_index int, offset_x float)                     { panic("not implemented") } // set position of column line (in pixels, from the left side of the contents region). pass -1 to use current column
func GetColumnsCount() int                                                 { panic("not implemented") }

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
