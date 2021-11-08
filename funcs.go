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

// Styles
func StyleColorsLight(dst *ImGuiStyle)   { panic("not implemented") } // best used with borders and a custom, thicker font
func StyleColorsClassic(dst *ImGuiStyle) { panic("not implemented") } // classic imgui style

// Widgets: Main
// - Most widgets return true when the value has been changed or when pressed/selected
// - You may also use one of the many IsItemXXX functions (e.g. IsItemActive, IsItemHovered, etc.) to query widget state.
func SmallButton(label string) bool { panic("not implemented") } // button with FramePadding=(0,0) to easily embed within text
func InvisibleButton(str_id string, size ImVec2, flsgs ImGuiButtonFlags) bool {
	panic("not implemented")
}                                                  // flexible button behavior without the visuals, frequently useful to build custom behaviors using the public api (along with IsItemActive, IsItemHovered, etc.)
func ArrowButton(str_id string, dir ImGuiDir) bool { panic("not implemented") } // square button with an arrow shape
func Image(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, tint_col ImVec4, border_col ImVec4) {
	panic("not implemented")
}
func ImageButton(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, frame_padding int /*/*= /*/, bg_col ImVec4, tint_col ImVec4) bool {
	panic("not implemented")
} // <0 frame_padding uses default frame padding settings. 0 for no padding

func RadioButtonBool(label string, active bool)              { panic("not implemented") } // use with e.g. if (RadioButton("one", my_value==1)) { my_value = 1 bool {panic("not implemented")} }
func RadioButtonInt(label string, v *int, v_button int) bool { panic("not implemented") } // shortcut to handle the above pattern when value is an integer
func ProgressBar(fraction float, size_arg ImVec2 /*= ImVec2(-FLT_MIN, 0)*/, overlay string) {
	panic("not implemented")
}
func Bullet() { panic("not implemented") } // draw a small circle + keep the cursor on the same line. advance cursor x position by GetTreeNodeToLabelSpacing(), same distance that TreeNode() uses

// Widgets: Drag Sliders
// - CTRL+Click on any drag box to turn them into an input box. Manually input values aren't clamped and can go off-bounds.
// - For all the Float2/Float3/Float4/Int2/Int3/Int4 versions of every functions, note that a 'v float[X]' function argument is the same as 'float* v', the array syntax is just a way to document the number of elements that are expected to be accessible. You can pass address of your first element out of a contiguous set, e.g. &myvector.x
// - Adjust format string to decorate the value with a prefix, a suffix, or adapt the editing and display precision e.g. "%.3f" -> 1.234; "%5.2 secs" -> 01.23 secs; "Biscuit: %.0f" -> Biscuit: 1; etc.
// - Format string may also be set to NULL or use the default format ("%f" or "%d").
// - Speed are per-pixel of mouse movement (v_speed=0.2: mouse needs to move by 5 pixels to increase value by 1). For gamepad/keyboard navigation, minimum speed is Max(v_speed, minimum_step_at_given_precision).
// - Use v_min < v_max to clamp edits to given limits. Note that CTRL+Click manual input can override those limits.
// - Use v_max/*= m*/,same with v_min = -FLT_MAX / INT_MIN to a clamping to a minimum.
// - We use the same sets of flags for DragXXX() and SliderXXX() functions as the features are the same and it makes it easier to swap them.
// - Legacy: Pre-1.78 there are DragXXX() function signatures that takes a final `power float=1.0' argument instead of the `ImGuiSliderFlags flags=0' argument.
//   If you get a warning converting a to float ImGuiSliderFlags, read https://github.com/ocornut/imgui/issues/3361
func DragFloat2(label string, v [2]float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragFloat3(label string, v [3]float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragFloat4(label string, v [4]float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragFloatRange2(label string, v_current_min *float, v_current_max *float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "*/, format_max string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragInt(label string, v *int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
} // If v_min >= v_max we have no bound
func DragInt2(label string, v [2]int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragInt3(label string, v [3]int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragInt4(label string, v [4]int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragIntRange2(label string, v_current_min *int, v_current_max *int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "*/, format_max string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragScalar(label string, data_type ImGuiDataType, p_data interface{}, v_speed float /*= 0*/, p_min interface{} /*= L*/, p_max interface{} /*= L*/, format string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragScalarN(label string, data_type ImGuiDataType, p_data interface{}, components int, v_speed float /*= 0*/, p_min interface{} /*= L*/, p_max interface{} /*= L*/, format string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}

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

// Widgets: Trees
// - TreeNode functions return true when the node is open, in which case you need to also call TreePop() when you are finished displaying the tree node contents.
func TreeNodeF(str_id string, fmt string, args ...interface{}) bool { panic("not implemented") } // helper variation to easily decorelate the id from the displayed string. Read the FAQ about why and how to use ID. to align arbitrary text at the same level as a TreeNode() you can use Bullet().
func TreeNodeInterface(ptr_id interface{}, fmt string, args ...interface{}) bool {
	panic("not implemented")
}                                                                  // "
func TreeNodeV(str_id string, fmt string, args []interface{}) bool { panic("not implemented") }
func TreeNodeInterfaceV(ptr_id interface{}, fmt string, args []interface{}) bool {
	panic("not implemented")
}
func TreeNodeEx(label string, flsgs ImGuiTreeNodeFlags) bool { panic("not implemented") }
func TreeNodeExF(str_id string, flags ImGuiTreeNodeFlags, fmt string, args ...interface{}) bool {
	panic("not implemented")
}
func TreeNodeInterfaceEx(ptr_id interface{}, flags ImGuiTreeNodeFlags, fmt string, args ...interface{}) bool {
	panic("not implemented")
}
func TreeNodeExV(str_id string, flags ImGuiTreeNodeFlags, fmt string, args []interface{}) bool {
	panic("not implemented")
}
func TreeNodeInterfaceExV(ptr_id interface{}, flags ImGuiTreeNodeFlags, fmt string, args []interface{}) bool {
	panic("not implemented")
}
func TreePush(str_id string)                                       { panic("not implemented") } // ~ Indent()+PushId(). Already called by TreeNode() when returning true, but you can call TreePush/TreePop yourself if desired.
func TreePushInterface(ptr_id interface{})                         { panic("not implemented") } // "
func GetTreeNodeToLabelSpacing() float                             { panic("not implemented") } // horizontal distance preceding label when using TreeNode*() or Bullet() == (g.FontSize + style.FramePadding.x*2) for a regular unframed TreeNode
func CollapsingHeader(label string, flsgs ImGuiTreeNodeFlags) bool { panic("not implemented") } // if returning 'true' the header is open. doesn't indent nor push on ID stack. user doesn't have to call TreePop().
func CollapsingHeaderVisible(label string, p_visible *bool, flsgs ImGuiTreeNodeFlags) bool {
	panic("not implemented")
}                                                  // when 'p_visible != NULL': if '*p_visible==true' display an additional small close button on upper right of the header which will set the to bool false when clicked, if '*p_visible==false' don't display the header.
func SetNextItemOpen(is_open bool, cond ImGuiCond) { panic("not implemented") } // set next TreeNode/CollapsingHeader open state.

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

// Tables
// [BETA API] API may evolve slightly! If you use this, please update to the next version when it comes out!
// - Full-featured replacement for old Columns API.
// - See Demo->Tables for demo code.
// - See top of imgui_tables.cpp for general commentary.
// - See ImGuiTableFlags_ and ImGuiTableColumnFlags_ enums for a description of available flags.
// The typical call flow is:
// - 1. Call BeginTable().
// - 2. Optionally call TableSetupColumn() to submit column name/flags/defaults.
// - 3. Optionally call TableSetupScrollFreeze() to request scroll freezing of columns/rows.
// - 4. Optionally call TableHeadersRow() to submit a header row. Names are pulled from TableSetupColumn() data.
// - 5. Populate contents:
//    - In most situations you can use TableNextRow() + TableSetColumnIndex(N) to start appending into a column.
//    - If you are using tables as a sort of grid, where every columns is holding the same type of contents,
//      you may prefer using TableNextColumn() instead of TableNextRow() + TableSetColumnIndex().
//      TableNextColumn() will automatically wrap-around into the next row if needed.
//    - IMPORTANT: Comparatively to the old Columns() API, we need to call TableNextColumn() for the first column!
//    - Summary of possible call flow:
//        --------------------------------------------------------------------------------------------------------
//        TableNextRow() -> TableSetColumnIndex(0) -> Text("Hello 0") -> TableSetColumnIndex(1) -> Text("Hello 1")  // OK
//        TableNextRow() -> TableNextColumn()      -> Text("Hello 0") -> TableNextColumn()      -> Text("Hello 1")  // OK
//                          TableNextColumn()      -> Text("Hello 0") -> TableNextColumn()      -> Text("Hello 1")  // OK: TableNextColumn() automatically gets to next row!
//        TableNextRow()                           -> Text("Hello 0")                                               // Not OK! Missing TableSetColumnIndex() or TableNextColumn()! Text will not appear!
//        --------------------------------------------------------------------------------------------------------
// - 5. Call EndTable()
func BeginTable(str_id string, column int, flsgs ImGuiTableFlags, outer_size ImVec2, inner_width float) bool {
	panic("not implemented")
}
func EndTable() { panic("not implemented") } // only call EndTable() if BeginTable() returns true!
func TableNextRow(row_flags ImGuiTableRowFlags /*= 0*/, min_row_height float) {
	panic("not implemented")
}                                           // append into the first cell of a new row.
func TableNextColumn() bool                 { panic("not implemented") } // append into the next column (or first column of next row if currently in last column). Return true when column is visible.
func TableSetColumnIndex(column_n int) bool { panic("not implemented") } // append into the specified column. Return true when column is visible.

// Tables: Headers & Columns declaration
// - Use TableSetupColumn() to specify label, resizing policy, default width/weight, id, various other flags etc.
// - Use TableHeadersRow() to create a header row and automatically submit a TableHeader() for each column.
//   Headers are required to perform: reordering, sorting, and opening the context menu.
//   The context menu can also be made available in columns body using ImGuiTableFlags_ContextMenuInBody.
// - You may manually submit headers using TableNextRow() + TableHeader() calls, but this is only useful in
//   some advanced use cases (e.g. adding custom widgets in header row).
// - Use TableSetupScrollFreeze() to lock columns/rows so they stay visible when scrolled.
func TableSetupColumn(label string, flsgs ImGuiTableColumnFlags, init_width_or_weight float /*= 0*/, user_id ImGuiID) {
	panic("not implemented")
}
func TableSetupScrollFreeze(cols int, rows int) { panic("not implemented") } // lock columns/rows so they stay visible when scrolled.
func TableHeadersRow()                          { panic("not implemented") } // submit all headers cells based on data provided to TableSetupColumn() + submit context menu
func TableHeader(label string)                  { panic("not implemented") } // submit one header cell manually (rarely used)

// Tables: Sorting
// - Call TableGetSortSpecs() to retrieve latest sort specs for the table. NULL when not sorting.
// - When 'SpecsDirty == true' you should sort your data. It will be true when sorting specs have changed
//   since last call, or the first time. Make sure to set 'SpecsDirty/*= g*/,else you may
//   wastefully sort your data every frame!
// - Lifetime: don't hold on this pointer over multiple frames or past any subsequent call to BeginTable().
func TableGetSortSpecs() *ImGuiTableSortSpecs { panic("not implemented") } // get latest sort specs for the table (NULL if not sorting).

// Tables: Miscellaneous functions
// - Functions args 'column_n int' treat the default value of -1 as the same as passing the current column index.
func TableGetColumnCount() int                                        { panic("not implemented") } // return number of columns (value passed to BeginTable)
func TableGetColumnIndex() int                                        { panic("not implemented") } // return current column index.
func TableGetRowIndex() int                                           { panic("not implemented") } // return current row index.
func TableGetColumnName(column_n int /*= -1*/) string                 { panic("not implemented") } // return "" if column didn't have a name declared by TableSetupColumn(). Pass -1 to use current column.
func TableGetColumnFlags(column_n int /*= -1*/) ImGuiTableColumnFlags { panic("not implemented") } // return column flags so you can query their Enabled/Visible/Sorted/Hovered status flags. Pass -1 to use current column.
func TableSetColumnEnabled(column_n int, v bool)                      { panic("not implemented") } // change user accessible enabled/disabled state of a column. Set to false to hide the column. User can use the context menu to change this themselves (right-click in headers, or right-click in columns body with ImGuiTableFlags_ContextMenuInBody)
func TableSetBgColor(target ImGuiTableBgTarget, color ImU32, column_n int /*= -1*/) {
	panic("not implemented")
} // change the color of a cell, row, or column. See ImGuiTableBgTarget_ flags for details.

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
