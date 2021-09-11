package imgui

// Context creation and access
// - Each context create its own ImFontAtlas by default. You may instance one yourself and pass it to CreateContext() to share a font atlas between contexts.
// - DLL users: heaps and globals are not shared across DLL boundaries! You will need to call SetCurrentContext() + SetAllocatorFunctions()
//   for each static/DLL boundary you are calling from. Read "Context and Memory Allocators" section of imgui.cpp for details.
func CreateContext(shared_font_atlas *ImFontAtlas) *ImGuiContext {
	var ctx ImGuiContext = NewImGuiContext(shared_font_atlas)
	if GImGui == nil {
		SetCurrentContext(&ctx)
	}
	Initialize(&ctx)
	return &ctx
}

func DestroyContext(ctx *ImGuiContext) { panic("not implemented") } // NULL = destroy current context
func GetCurrentContext() *ImGuiContext { panic("not implemented") }

func SetCurrentContext(ctx *ImGuiContext) {
	GImGui = ctx
}

// access the IO structure (mouse/keyboard/gamepad inputs, time, various configuration options/flags)
func GetIO() *ImGuiIO {
	IM_ASSERT_USER_ERROR(GImGui != nil, "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?")
	return &GImGui.IO
}

func GetStyle() *ImGuiStyle { panic("not implemented") } // access the Style structure (colors, sizes). Always use PushStyleCol(), PushStyleVar() to modify style mid-frame!

// Pass this to your backend rendering function! Valid after Render() and until the next call to NewFrame()
func GetDrawData() *ImDrawData {
	var g = GImGui
	var viewport = g.Viewports[0]
	if viewport.DrawDataP.Valid {
		return &viewport.DrawDataP
	} else {
		return nil
	}
} // valid after Render() and until the next call to NewFrame(). this is what you have to render.

// Demo, Debug, Information
func ShowDemoWindow(p_open *bool)         { panic("not implemented") } // create Demo window. demonstrate most ImGui features. call this to learn about the library! try to make it always available in your application!
func ShowMetricsWindow(p_open *bool)      { panic("not implemented") } // create Metrics/Debugger window. display Dear ImGui internals: windows, draw commands, various internal state, etc.
func ShowAboutWindow(p_open *bool)        { panic("not implemented") } // create About window. display Dear ImGui version, credits and build/system information.
func ShowStyleEditor(ref *ImGuiStyle)     { panic("not implemented") } // add style editor block (not a window). you can pass in a reference ImGuiStyle structure to compare to, revert to and save to (else it uses the default style)
func ShowStyleSelector(label string) bool { panic("not implemented") } // add style selector block (not a window), essentially a combo listing the default styles.
func ShowFontSelector(label string)       { panic("not implemented") } // add font selector block (not a window), essentially a combo listing the loaded fonts.
func ShowUserGuide()                      { panic("not implemented") } // add basic help/info block (not a window): how to manipulate ImGui as a end-user (mouse/keyboard controls).
func GetVersion() string                  { panic("not implemented") } // get the compiled version string e.g. "1.80 WIP" (essentially the value for IMGUI_VERSION from the compiled version of imgui.cpp)

// Styles
func StyleColorsLight(dst *ImGuiStyle)   { panic("not implemented") } // best used with borders and a custom, thicker font
func StyleColorsClassic(dst *ImGuiStyle) { panic("not implemented") } // classic imgui style

// Child Windows
// - Use child windows to begin into a self-contained independent scrolling/clipping regions within a host window. Child windows can embed their own child.
// - For each independent axis of 'size': ==0.0: use remaining host window size / >0.0: fixed size / <0.0: use remaining window size minus abs(size) / Each axis can use a different mode, e.g. ImVec2(0,400).
// - BeginChild() returns false to indicate the window is collapsed or fully clipped, so you may early out and omit submitting anything to the window.
//   Always call a matching EndChild() for each BeginChild() call, regardless of its return value.
//   [Important: due to legacy reason, this is inconsistent with most other functions such as BeginMenu/EndMenu,
//    BeginPopup/EndPopup, etc. where the EndXXX call should only be called if the corresponding BeginXXX function
//    returned true. Begin and BeginChild are the only odd ones out. Will be fixed in a future update.]
func BeginChild(str_id string, size ImVec2, border bool, flags ImGuiWindowFlags) bool {
	panic("not implemented")
}
func BeginChildID(id ImGuiID, size ImVec2, border bool, flags ImGuiWindowFlags) bool {
	panic("not implemented")
}
func EndChild() { panic("not implemented") }

// Windows Utilities
// - 'current window' = the window we are appending into while inside a Begin()/End() block. 'next window' = next window we will Begin() into.
func IsWindowAppearing() bool                      { panic("not implemented") }
func IsWindowCollapsed() bool                      { panic("not implemented") }
func IsWindowFocused(flags ImGuiFocusedFlags) bool { panic("not implemented") } // is current window focused? or its root/child, depending on flags. see flags for options.
func IsWindowHovered(flags ImGuiHoveredFlags) bool { panic("not implemented") } // is current window hovered (and typically: not blocked by a popup/modal)? see flags for options. NB: If you are trying to check whether your mouse should be dispatched to imgui or to your app, you should use the 'io.WantCaptureMouse' boolean for that! Please read the FAQ!
func GetWindowDrawList() *ImDrawList               { panic("not implemented") } // get draw list associated to the current window, to append your own drawing primitives
func GetWindowPos() ImVec2                         { panic("not implemented") } // get current window position in screen space (useful if you want to do your own drawing via the DrawList API)
func GetWindowSize() ImVec2                        { panic("not implemented") } // get current window size
func GetWindowWidth() float                        { panic("not implemented") } // get current window width (shortcut for GetWindowSize().x)
func GetWindowHeight() float                       { panic("not implemented") } // get current window height (shortcut for GetWindowSize().y)

// Window manipulation
// - Prefer using SetNextXXX functions (before Begin) rather that SetXXX functions (after Begin).
func SetNextWindowPos(pos *ImVec2, cond ImGuiCond, pivot ImVec2) { panic("not implemented") } // set next window position. call before Begin(). use pivot=(0.5,0.5) to center on given point, etc.

func SetNextWindowSizeConstraints(size_min ImVec2, size_max ImVec2, custom_callback ImGuiSizeCallback, custom_callback_data interface{}) {
	panic("not implemented")
}                                                                         // set next window size limits. use -1,-1 on either X/Y axis to preserve the current size. Sizes will be rounded down. Use callback to apply non-trivial programmatic constraints.
func SetNextWindowContentSize(size ImVec2)                                { panic("not implemented") } // set next window content size (~ scrollable client area, which enforce the range of scrollbars). Not including window decorations (title bar, menu bar, etc.) nor WindowPadding. set an axis to 0.0 to leave it automatic. call before Begin()
func SetNextWindowCollapsed(collapsed bool, cond ImGuiCond)               { panic("not implemented") } // set next window collapsed state. call before Begin()
func SetNextWindowFocus()                                                 { panic("not implemented") } // set next window to be focused / top-most. call before Begin()
func SetNextWindowBgAlpha(alpha float)                                    { panic("not implemented") } // set next window background color alpha. helper to easily override the Alpha component of ImGuiCol_WindowBg/ChildBg/PopupBg. you may also use ImGuiWindowFlags_NoBackground.
func SetWindowPos(pos ImVec2, cond ImGuiCond)                             { panic("not implemented") } // (not recommended) set current window position - call within Begin()/End(). prefer using SetNextWindowPos(), as this may incur tearing and side-effects.
func SetWindowSize(size ImVec2, cond ImGuiCond)                           { panic("not implemented") } // (not recommended) set current window size - call within Begin()/End(). set to ImVec2(0, 0) to force an auto-fit. prefer using SetNextWindowSize(), as this may incur tearing and minor side-effects.
func SetWindowCollapsed(collapsed bool, cond ImGuiCond)                   { panic("not implemented") } // (not recommended) set current window collapsed state. prefer using SetNextWindowCollapsed().
func SetWindowFocus()                                                     { panic("not implemented") } // (not recommended) set current window to be focused / top-most. prefer using SetNextWindowFocus().
func SetWindowFontScale(scale float)                                      { panic("not implemented") } // [OBSOLETE] set font scale. Adjust IO.FontGlobalScale if you want to scale all windows. This is an old API! For correct scaling, prefer to reload font + rebuild ImFontAtlas + call style.ScaleAllSizes().
func SetNamedWindowPos(name string, pos ImVec2, cond ImGuiCond)           { panic("not implemented") } // set named window position.
func SetNamedWindowSize(name string, size ImVec2, cond ImGuiCond)         { panic("not implemented") } // set named window size. set axis to 0.0 to force an auto-fit on this axis.
func SetNamedWindowCollapsed(name string, collapsed bool, cond ImGuiCond) { panic("not implemented") } // set named window collapsed state
func SetNamedWindowFocus(name string)                                     { panic("not implemented") } // set named window to be focused / top-most. use NULL to remove focus.

// Content region
// - Retrieve available space from a given point. GetContentRegionAvail() is frequently useful.
// - Those functions are bound to be redesigned (they are confusing, incomplete and the Min/Max return values are in local window coordinates which increases confusion)
func GetContentRegionAvail() ImVec2     { panic("not implemented") } // == GetContentRegionMax() - GetCursorPos()
func GetContentRegionMax() ImVec2       { panic("not implemented") } // current content boundaries (typically window boundaries including scrolling, or current column boundaries), in windows coordinates
func GetWindowContentRegionMin() ImVec2 { panic("not implemented") } // content boundaries min for the full window (roughly (0,0)-Scroll), in window coordinates
func GetWindowContentRegionMax() ImVec2 { panic("not implemented") } // content boundaries max for the full window (roughly (0,0)+Size-Scroll) where Size can be override with SetNextWindowContentSize(), in window coordinates

// Windows Scrolling
func GetScrollX() float                                         { panic("not implemented") } // get scrolling amount [0 .. GetScrollMaxX()]
func GetScrollY() float                                         { panic("not implemented") } // get scrolling amount [0 .. GetScrollMaxY()]
func SetScrollX(scroll_x float)                                 { panic("not implemented") } // set scrolling amount [0 .. GetScrollMaxX()]
func SetScrollY(scroll_y float)                                 { panic("not implemented") } // set scrolling amount [0 .. GetScrollMaxY()]
func GetScrollMaxX() float                                      { panic("not implemented") } // get maximum scrolling amount ~~ ContentSize.x - WindowSize.x - DecorationsSize.x
func GetScrollMaxY() float                                      { panic("not implemented") } // get maximum scrolling amount ~~ ContentSize.y - WindowSize.y - DecorationsSize.y
func SetScrollHereX(center_x_ratio float /*= 0.5*/)             { panic("not implemented") } // adjust scrolling amount to make current cursor position visible. center_x_ratio=0.0: left, 0.5: center, 1.0: right. When using to make a "default/current item" visible, consider using SetItemDefaultFocus() instead.
func SetScrollHereY(center_y_ratio float /*= 0.5*/)             { panic("not implemented") } // adjust scrolling amount to make current cursor position visible. center_y_ratio=0.0: top, 0.5: center, 1.0: bottom. When using to make a "default/current item" visible, consider using SetItemDefaultFocus() instead.
func SetScrollFromPosX(local_x, center_x_ratio float /*= 0.5*/) { panic("not implemented") } // adjust scrolling amount to make given position visible. Generally GetCursorStartPos() + offset to compute a valid position.
func SetScrollFromPosY(local_y, center_y_ratio float /*= 0.5*/) { panic("not implemented") } // adjust scrolling amount to make given position visible. Generally GetCursorStartPos() + offset to compute a valid position.

// Parameters stacks (shared)
func PushFont(font ImFont)                             { panic("not implemented") } // use NULL as a shortcut to push default font
func PopFont()                                         { panic("not implemented") }
func PushStyleColorInt(idx ImGuiCol, col ImU32)        { panic("not implemented") } // modify a style color. always use this if you modify the style after NewFrame().
func PushStyleColorVec(idx ImGuiCol, col ImVec4)       { panic("not implemented") }
func PopStyleColor(count int /*= 1*/)                  { panic("not implemented") }
func PushStyleFloat(idx ImGuiStyleVar, val float)      { panic("not implemented") } // modify a style variable float. always use this if you modify the style after NewFrame().
func PushStyleVec(idx ImGuiStyleVar, val ImVec2)       { panic("not implemented") } // modify a style variable ImVec2. always use this if you modify the style after NewFrame().
func PopStyleVar(count int /*= 1*/)                    { panic("not implemented") }
func PushAllowKeyboardFocus(allow_keyboard_focus bool) { panic("not implemented") } // == tab stop enable. Allow focusing using TAB/Shift-TAB, enabled by default but you can disable it for certain widgets
func PopAllowKeyboardFocus()                           { panic("not implemented") }
func PushButtonRepeat(repeat bool)                     { panic("not implemented") } // in 'repeat' mode, Button*() functions return repeated true in a typematic manner (using io.KeyRepeatDelay/io.KeyRepeatRate setting). Note that you can call IsItemActive() after any Button() to tell if the button is held in the current frame.
func PopButtonRepeat()                                 { panic("not implemented") }

// Parameters stacks (current window)
func PushItemWidth(item_width float)         { panic("not implemented") } // push width of items for common large "item+label" widgets. >0.0: width in pixels, <0.0 align xx pixels to the right of window (so -FLT_MIN always align width to the right side).
func PopItemWidth()                          { panic("not implemented") }
func SetNextItemWidth(item_width float)      { panic("not implemented") } // set width of the _next_ common large "item+label" widget. >0.0: width in pixels, <0.0 align xx pixels to the right of window (so -FLT_MIN always align width to the right side)
func CalcItemWidth() float                   { panic("not implemented") } // width of item given pushed settings and current cursor position. NOT necessarily the width of last item unlike most 'Item' functions.
func PushTextWrapPos(wrap_local_pos_x float) { panic("not implemented") } // push word-wrapping position for Text*() commands. < 0.0: no wrapping; 0.0: wrap to end of window (or column)  {panic("not implemented")} > 0.0: wrap at 'wrap_pos_x' position in window local space
func PopTextWrapPos()                        { panic("not implemented") }

// Style read access
// - Use the style editor (ShowStyleEditor() function) to interactively see what the colors are)
func GetFont() *ImFont               { panic("not implemented") } // get current font
func GetFontSize() float             { panic("not implemented") } // get current font size (= height in pixels) of current font with current scale applied
func GetFontTexUvWhitePixel() ImVec2 { panic("not implemented") } // get UV coordinate for a while pixel, useful to draw custom shapes via the ImDrawList API

func GetColorU32FromID(idx ImGuiCol, alpha_mul float /*= 1.0*/) ImU32 {
	var style = GImGui.Style
	var c ImVec4 = style.Colors[idx]
	c.w *= style.Alpha * alpha_mul
	return ColorConvertFloat4ToU32(c)
}

// retrieve given style color with style alpha applied and optional extra alpha multiplier, packed as a 32-bit value suitable for ImDrawList
func GetColorU32FromVec(col ImVec4) ImU32 { panic("not implemented") } // retrieve given color with style alpha applied, packed as a 32-bit value suitable for ImDrawList

func GetColorU32FromInt(col ImU32) ImU32 { panic("not implemented") } // retrieve given color with style alpha applied, packed as a 32-bit value suitable for ImDrawList

func GetStyleColorVec4(idx ImGuiCol) *ImVec4 { panic("not implemented") } // retrieve style color as stored in ImGuiStyle structure. use to feed back into PushStyleColor(), otherwise use GetColorU32() to get style color with style alpha baked in.

// Cursor / Layout
// - By "cursor" we mean the current output position.
// - The typical widget behavior is to output themselves at the current cursor position, then move the cursor one line down.
// - You can call SameLine() between widgets to undo the last carriage return and output at the right of the preceding widget.
// - Attention! We currently have inconsistencies between window-local and absolute positions we will aim to fix with future API:
//    Window-local coordinates:   SameLine(), GetCursorPos(), SetCursorPos(), GetCursorStartPos(), GetContentRegionMax(), GetWindowContentRegion*(), PushTextWrapPos()
//    Absolute coordinate:        GetCursorScreenPos(), SetCursorScreenPos(), all ImDrawList:: functions.
func Separator()                          { panic("not implemented") } // separator, generally horizontal. inside a menu bar or in horizontal layout mode, this becomes a vertical separator.
func NewLine()                            { panic("not implemented") } // undo a SameLine() or force a new line when in an horizontal-layout context.
func Spacing()                            { panic("not implemented") } // add vertical spacing.
func Dummy(size ImVec2)                   { panic("not implemented") } // add a dummy item of given size. unlike InvisibleButton(), Dummy() won't take the mouse click or be navigable into.
func Indent(indent_w float)               { panic("not implemented") } // move content position toward the right, by indent_w, or style.IndentSpacing if indent_w <= 0
func Unindent(indent_w float)             { panic("not implemented") } // move content position back to the left, by indent_w, or style.IndentSpacing if indent_w <= 0
func BeginGroup()                         { panic("not implemented") } // lock horizontal starting position
func EndGroup()                           { panic("not implemented") } // unlock horizontal starting position + capture the whole group bounding box into one "item" (so you can use IsItemHovered() or layout primitives such as SameLine() on whole group, etc.)
func GetCursorPos() ImVec2                { panic("not implemented") } // cursor position in window coordinates (relative to window position)
func GetCursorPosX() float                { panic("not implemented") } //   (some functions are using window-relative coordinates, such as: GetCursorPos, GetCursorStartPos, GetContentRegionMax, GetWindowContentRegion* etc.
func GetCursorPosY() float                { panic("not implemented") } //    other functions such as GetCursorScreenPos or everything in ImDrawList::
func SetCursorPos(local_pos *ImVec2)      { panic("not implemented") } //    are using the main, absolute coordinate system.
func SetCursorPosX(local_x float)         { panic("not implemented") } //    GetWindowPos() + GetCursorPos() == GetCursorScreenPos() etc.)
func SetCursorPosY(local_y float)         { panic("not implemented") } //
func GetCursorStartPos() ImVec2           { panic("not implemented") } // initial cursor position in window coordinates
func GetCursorScreenPos() ImVec2          { panic("not implemented") } // cursor position in absolute coordinates (useful to work with ImDrawList API). generally top-left == GetMainViewport()->Pos == (0,0) in single viewport mode, and bottom-right == GetMainViewport()->Pos+Size == io.DisplaySize in single-viewport mode.
func SetCursorScreenPos(pos ImVec2)       { panic("not implemented") } // cursor position in absolute coordinates
func AlignTextToFramePadding()            { panic("not implemented") } // vertically align upcoming text baseline to FramePadding.y so that it will align properly to regularly framed items (call if you have text on a line before a framed item)
func GetTextLineHeight() float            { panic("not implemented") } // ~ FontSize
func GetTextLineHeightWithSpacing() float { panic("not implemented") } // ~ FontSize + style.ItemSpacing.y (distance in pixels between 2 consecutive lines of text)
func GetFrameHeight() float               { panic("not implemented") } // ~ FontSize + style.FramePadding.y * 2
func GetFrameHeightWithSpacing() float    { panic("not implemented") } // ~ FontSize + style.FramePadding.y * 2 + style.ItemSpacing.y (distance in pixels between 2 consecutive lines of framed widgets)

// ID stack/scopes
// Read the FAQ (docs/FAQ.md or http://dearimgui.org/faq) for more details about how ID are handled in dear imgui.
// - Those questions are answered and impacted by understanding of the ID stack system:
//   - "Q: Why is my widget not reacting when I click on it?"
//   - "Q: How can I have widgets with an empty label?"
//   - "Q: How can I have multiple widgets with the same label?"
// - Short version: ID are hashes of the entire ID stack. If you are creating widgets in a loop you most likely
//   want to push a unique identifier (e.g. object pointer, loop index) to uniquely differentiate them.
// - You can also use the "Label##foobar" syntax within widget label to distinguish them from each others.
// - In this header file we use the "label"/"name" terminology to denote a string that will be displayed + used as an ID,
//   whereas "str_id" denote a string that is only used as an ID and not normally displayed.
func PushString(str_id string) {
	var g = GImGui
	var window = g.CurrentWindow
	var id = window.GetIDNoKeepAlive(str_id)
	window.IDStack = append(window.IDStack, id)
}                                                    // push string into the ID stack (will hash string).
func PushIDs(str_id_begin string, str_id_end string) { panic("not implemented") } // push string into the ID stack (will hash string).
func PushInterface(ptr_id interface{})               { panic("not implemented") } // push pointer into the ID stack (will hash pointer).
func PushID(int_id int)                              { panic("not implemented") } // push integer into the ID stack (will hash integer).

func PopID() {
	var window = GImGui.CurrentWindow
	IM_ASSERT(len(window.IDStack) > 1) // Too many PopID(), or could be popping in a wrong/different window?
	window.IDStack = window.IDStack[:len(window.IDStack)-1]
} // pop from the ID stack.

func GetIDFromString(str_id string) ImGuiID {
	return GImGui.CurrentWindow.GetIDs(str_id, "")

} // calculate unique ID (hash of whole ID stack + given parameter). e.g. if you want to query into ImGuiStorage yourself

func GetIDs(str_id_begin string, str_id_end string) ImGuiID {
	return GImGui.CurrentWindow.GetIDs(str_id_begin, str_id_end)
}

func GetIDFromInterface(ptr_id interface{}) ImGuiID { panic("not implemented") }

// Widgets: Text
func TextUnformatted(text string, text_end string)            { panic("not implemented") } // raw text without formatting. Roughly equivalent to Text("%s", text) but: A) doesn't require null terminated string if 'text_end' is specified, B) it's faster, no memory copy is done, no buffer size limits, recommended for long chunks of text.
func TextColored(col ImVec4, fmt string, args ...interface{}) { panic("not implemented") } // shortcut for PushStyleColor(ImGuiCol_Text, col); Text(fmt, ...); PopStyleColor()  {panic("not implemented")}
func TextDisabled(fmt string, args ...interface{})            { panic("not implemented") } // shortcut for PushStyleColor(ImGuiCol_Text, style.Colors[ImGuiCol_TextDisabled]); Text(fmt, ...); PopStyleColor()  {panic("not implemented")}
func TextWrapped(fmt string, args ...interface{})             { panic("not implemented") } // shortcut for PushTextWrapPos(0.0); Text(fmt, ...); PopTextWrapPos()  {panic("not implemented")}. Note that this won't work on an auto-resizing window if there's no other widgets to extend the window width, yoy may need to set a size using SetNextWindowSize().
func LabelText(label string, fmt string, args ...interface{}) { panic("not implemented") } // display text+label aligned the same way as value+label widgets
func BulletText(fmt string, args ...interface{})              { panic("not implemented") } // shortcut for Bullet()+Text()

// Widgets: Main
// - Most widgets return true when the value has been changed or when pressed/selected
// - You may also use one of the many IsItemXXX functions (e.g. IsItemActive, IsItemHovered, etc.) to query widget state.
func Button(label string, size ImVec2) bool { panic("not implemented") } // button
func SmallButton(label string) bool         { panic("not implemented") } // button with FramePadding=(0,0) to easily embed within text
func InvisibleButton(str_id string, size ImVec2, flsgs ImGuiButtonFlags) bool {
	panic("not implemented")
}                                                  // flexible button behavior without the visuals, frequently useful to build custom behaviors using the public api (along with IsItemActive, IsItemHovered, etc.)
func ArrowButton(str_id string, dir ImGuiDir) bool { panic("not implemented") } // square button with an arrow shape
func Image(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, tint_col ImVec4, border_col ImVec4) {
	panic("not implemented")
}
func ImageButton(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, frame_padding int /*/*= /*/, bg_col ImVec4, tint_col ImVec4) bool {
	panic("not implemented")
}                                                                         // <0 frame_padding uses default frame padding settings. 0 for no padding
func Checkbox(label string, v *bool) bool                                 { panic("not implemented") }
func CheckboxFlagsInt(label string, flags *int, flags_value int) bool     { panic("not implemented") }
func CheckboxFlagsUint(label string, flags *uint, uflags_value uint) bool { panic("not implemented") }
func RadioButtonBool(label string, active bool)                           { panic("not implemented") } // use with e.g. if (RadioButton("one", my_value==1)) { my_value = 1 bool {panic("not implemented")} }
func RadioButtonInt(label string, v *int, v_button int) bool              { panic("not implemented") } // shortcut to handle the above pattern when value is an integer
func ProgressBar(fraction float, size_arg ImVec2 /*= ImVec2(-FLT_MIN, 0)*/, overlay string) {
	panic("not implemented")
}
func Bullet() { panic("not implemented") } // draw a small circle + keep the cursor on the same line. advance cursor x position by GetTreeNodeToLabelSpacing(), same distance that TreeNode() uses

// Widgets: Combo Box
// - The BeginCombo()/EndCombo() api allows you to manage your contents and selection state however you want it, by creating e.g. Selectable() items.
// - The old Combo() api are helpers over BeginCombo()/EndCombo() which are kept available for convenience purpose. This is analogous to how ListBox are created.
func BeginCombo(label string, preview_value string, flsgs ImGuiComboFlags) bool {
	panic("not implemented")
}
func EndCombo() { panic("not implemented") } // only call EndCombo() if BeginCombo() returns true!
func ComboSlice(label string, current_item *int, items []string, items_count int, popup_max_height_in_items int /*= -1*/) bool {
	panic("not implemented")
}
func ComboString(label string, current_item *int, items_separated_by_zeros string, popup_max_height_in_items int /*= -1*/) bool {
	panic("not implemented")
} // Separate items with \0 within a string, end item-list with \0\0. e.g. "One\0Two\0Three\0"
func ComboFunc(label string, current_item *int, items_getter func(data, idx int, out_text *string) bool, data interface{}, items_count, popup_max_height_in_items int /*= -1*/) bool {
	panic("not implemented")
}

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
func TreeNode(label string) bool                                    { panic("not implemented") }
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
func TreePop()                                                     { panic("not implemented") } // ~ Unindent()+PopId()
func GetTreeNodeToLabelSpacing() float                             { panic("not implemented") } // horizontal distance preceding label when using TreeNode*() or Bullet() == (g.FontSize + style.FramePadding.x*2) for a regular unframed TreeNode
func CollapsingHeader(label string, flsgs ImGuiTreeNodeFlags) bool { panic("not implemented") } // if returning 'true' the header is open. doesn't indent nor push on ID stack. user doesn't have to call TreePop().
func CollapsingHeaderVisible(label string, p_visible *bool, flsgs ImGuiTreeNodeFlags) bool {
	panic("not implemented")
}                                                  // when 'p_visible != NULL': if '*p_visible==true' display an additional small close button on upper right of the header which will set the to bool false when clicked, if '*p_visible==false' don't display the header.
func SetNextItemOpen(is_open bool, cond ImGuiCond) { panic("not implemented") } // set next TreeNode/CollapsingHeader open state.

// Widgets: Selectables
// - A selectable highlights when hovered, and can display another color when selected.
// - Neighbors selectable extend their highlight bounds in order to leave no gap between them. This is so a series of selected Selectable appear contiguous.
func Selectable(label string, selected bool, flsgs ImGuiSelectableFlags, size ImVec2) bool {
	panic("not implemented")
} // "selected bool" carry the selection state (read-only). Selectable() is clicked is returns true so you can modify your selection state. size.x==0.0: use remaining width, size.x>0.0: specify width. size.y==0.0: use label height, size.y>0.0: specify height
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

// Tooltips
// - Tooltip are windows following the mouse. They do not take focus away.
func BeginTooltip()                              { panic("not implemented") } // begin/append a tooltip window. to create full-featured tooltip (with any kind of items).
func EndTooltip()                                { panic("not implemented") }
func SetTooltip(fmt string, args ...interface{}) { panic("not implemented") } // set a text-only tooltip, typically use with ImGui::IsItemHovered(). override any previous call to SetTooltip().
func SetTooltipV(fmt string, args []interface{}) { panic("not implemented") }

// Popups, Modals
//  - They block normal mouse hovering detection (and therefore most mouse interactions) behind them.
//  - If not modal: they can be closed by clicking anywhere outside them, or by pressing ESCAPE.
//  - Their visibility state (~bool) is held internally instead of being held by the programmer as we are used to with regular Begin*() calls.
//  - The 3 properties above are related: we need to retain popup visibility state in the library because popups may be closed as any time.
//  - You can bypass the hovering restriction by using ImGuiHoveredFlags_AllowWhenBlockedByPopup when calling IsItemHovered() or IsWindowHovered().
//  - IMPORTANT: Popup identifiers are relative to the current ID stack, so OpenPopup and BeginPopup generally needs to be at the same level of the stack.
//    This is sometimes leading to confusing mistakes. May rework this in the future.

// Popups: begin/end functions
//  - BeginPopup(): query popup state, if open start appending into the window. Call EndPopup() afterwards. ImGuiWindowFlags are forwarded to the window.
//  - BeginPopupModal(): block every interactions behind the window, cannot be closed by user, add a dimming background, has a title bar.
func BeginPopup(str_id string, flsgs ImGuiWindowFlags) bool { panic("not implemented") } // return true if the popup is open, and you can start outputting to it.
func BeginPopupModal(name string, p_open *bool, flsgs ImGuiWindowFlags) bool {
	panic("not implemented")
}               // return true if the modal is open, and you can start outputting to it.
func EndPopup() { panic("not implemented") } // only call EndPopup() if BeginPopupXXX() returns true!

// Popups: open/close functions
//  - OpenPopup(): set popup state to open. are ImGuiPopupFlags available for opening options.
//  - If not modal: they can be closed by clicking anywhere outside them, or by pressing ESCAPE.
//  - CloseCurrentPopup(): use inside the BeginPopup()/EndPopup() scope to close manually.
//  - CloseCurrentPopup() is called by default by Selectable()/MenuItem() when activated (FIXME: need some options).
//  - Use ImGuiPopupFlags_NoOpenOverExistingPopup to a opening a popup if there's already one at the same level. This is equivalent to e.g. testing for !IsAnyPopupOpen() prior to OpenPopup().
//  - Use IsWindowAppearing() after BeginPopup() to tell if a window just opened.
func OpenPopup(str_id string, popup_flags ImGuiPopupFlags) { panic("not implemented") } // call to mark popup as open (don't call every frame!).
func OpenPopupID(id ImGuiID, popup_flags ImGuiPopupFlags)  { panic("not implemented") } // id overload to facilitate calling from nested stacks
func OpenPopupOnItemClick(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) {
	panic("not implemented")
}                        // helper to open popup when clicked on last item. Default to ImGuiPopupFlags_MouseButtonRight == 1. (note: actually triggers on the mouse _released_ event to be consistent with popup behaviors)
func CloseCurrentPopup() { panic("not implemented") } // manually close the popup we have begin-ed into.

// Popups: open+begin combined functions helpers
//  - Helpers to do OpenPopup+BeginPopup where the Open action is triggered by e.g. hovering an item and right-clicking.
//  - They are convenient to easily create context menus, hence the name.
//  - IMPORTANT: Notice that BeginPopupContextXXX takes just ImGuiPopupFlags like OpenPopup() and unlike BeginPopup(). For full consistency, we may add ImGuiWindowFlags to the BeginPopupContextXXX functions in the future.
//  - IMPORTANT: we exceptionally default their flags to 1 (== ImGuiPopupFlags_MouseButtonRight) for backward compatibility with older API taking 'mouse_button int/*= r*/,so if you add other flags remember to re-add the ImGuiPopupFlags_MouseButtonRight.
func BeginPopupContextItem(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) bool {
	panic("not implemented")
} // open+begin popup when clicked on last item. Use str_id==NULL to associate the popup to previous item. If you want to use that on a non-interactive item such as Text() you need to pass in an explicit ID here. read comments in .cpp!
func BeginPopupContextWindow(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) bool {
	panic("not implemented")
} // open+begin popup when clicked on current window.
func BeginPopupContext(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) bool {
	panic("not implemented")
} // open+begin popup when clicked in  (where there are no windows).

// Popups: query functions
//  - IsPopupOpen(): return true if the popup is open at the current BeginPopup() level of the popup stack.
//  - IsPopupOpen() with ImGuiPopupFlags_AnyPopupId: return true if any popup is open at the current BeginPopup() level of the popup stack.
//  - IsPopupOpen() with ImGuiPopupFlags_AnyPopupId + ImGuiPopupFlags_AnyPopupLevel: return true if any popup is open.
func IsPopupOpen(str_id string, flsgs ImGuiPopupFlags) bool { panic("not implemented") } // return true if the popup is open.

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

// Logging/Capture
// - All text output from the interface can be captured into tty/file/clipboard. By default, tree nodes are automatically opened during logging.
func LogToTTY(auto_open_depth int /*= -1*/)                  { panic("not implemented") } // start logging to tty (stdout)
func LogToFile(auto_open_depth int /*= 1*/, filename string) { panic("not implemented") } // start logging to file
func LogToClipboard(auto_open_depth int /*= -1*/)            { panic("not implemented") } // start logging to OS clipboard

func LogFinish() {
	var g = GImGui
	if !g.LogEnabled {
		return
	}

	LogText(IM_NEWLINE)
	switch g.LogType {
	case ImGuiLogType_TTY:
		//g.LogFile
		break
	case ImGuiLogType_File:
		ImFileClose(g.LogFile)
		break
	case ImGuiLogType_Buffer:
		break
	case ImGuiLogType_Clipboard:
		if len(g.LogBuffer) > 0 {
			SetClipboardText(string(g.LogBuffer))
		}
		break
	case ImGuiLogType_None:
		IM_ASSERT(false)
		break
	}

	g.LogEnabled = false
	g.LogType = ImGuiLogType_None
	g.LogFile = nil
	g.LogBuffer = g.LogBuffer[:0]
} // stop logging (close file, etc.)

func LogButtons()                             { panic("not implemented") } // helper to display buttons for logging to tty/file/clipboard
func LogText(fmt string, args ...interface{}) { panic("not implemented") } // pass text data straight to log (without being displayed)
func LogTextV(fmt string, args []interface{}) { panic("not implemented") }

// Drag and Drop
// - On source items, call BeginDragDropSource(), if it returns true also call SetDragDropPayload() + EndDragDropSource().
// - On target candidates, call BeginDragDropTarget(), if it returns true also call AcceptDragDropPayload() + EndDragDropTarget().
// - If you stop calling BeginDragDropSource() the payload is preserved however it won't have a preview tooltip (we currently display a fallback "..." tooltip, see #1725)
// - An item can be both drag source and drop target.
func BeginDragDropSource(flags ImGuiDragDropFlags) bool { panic("not implemented") } // call after submitting an item which may be dragged. when this return true, you can call SetDragDropPayload() + EndDragDropSource()
func SetDragDropPayload(ptype string, data interface{}, sz uintptr, cond ImGuiCond) bool {
	panic("not implemented")
}                               // type is a user defined string of maximum 32 characters. Strings starting with '_' are reserved for dear imgui internal types. Data is copied and held by imgui.
func EndDragDropSource()        { panic("not implemented") } // only call EndDragDropSource() if BeginDragDropSource() returns true!
func BeginDragDropTarget() bool { panic("not implemented") } // call after submitting an item that may receive a payload. If this returns true, you can call AcceptDragDropPayload() + EndDragDropTarget()
func AcceptDragDropPayload(ptype string, flsgs ImGuiDragDropFlags) *ImGuiPayload {
	panic("not implemented")
}                                       // accept contents of a given type. If ImGuiDragDropFlags_AcceptBeforeDelivery is set you can peek into the payload before the mouse button is released.
func EndDragDropTarget()                { panic("not implemented") } // only call EndDragDropTarget() if BeginDragDropTarget() returns true!
func GetDragDropPayload() *ImGuiPayload { panic("not implemented") } // peek directly into the current payload from anywhere. may return NULL. use ImGuiPayload::IsDataType() to test for the payload type.

// Disabling [BETA API]
// - Disable all user interactions and dim items visuals (applying style.DisabledAlpha over current colors)
// - Those can be nested but it cannot be used to enable an already disabled section (a single BeginDisabled(true) in the stack is enough to keep everything disabled)
// - BeginDisabled(false) essentially does nothing useful but is provided to facilitate use of boolean expressions. If you can a calling BeginDisabled(False)/EndDisabled() best to a it.
func BeginDisabled(disabled bool /*= true*/) { panic("not implemented") }
func EndDisabled()                           { panic("not implemented") }

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

// Focus, Activation
// - Prefer using "SetItemDefaultFocus()" over "if (IsWindowAppearing()) SetScrollHereY()" when applicable to signify "this is the default item"
func SetItemDefaultFocus()            { panic("not implemented") } // make last item the default focused item of a window.
func SetKeyboardFocusHere(offset int) { panic("not implemented") } // focus keyboard on the next widget. Use positive 'offset' to access sub components of a multiple component widget. Use -1 to access previous widget.

// Item/Widgets Utilities and Query Functions
// - Most of the functions are referring to the previous Item that has been submitted.
// - See Demo Window under "Widgets->Querying Status" for an interactive visualization of most of those functions.
func IsItemHovered(flags ImGuiHoveredFlags) bool       { panic("not implemented") } // is the last item hovered? (and usable, aka not blocked by a popup, etc.). See ImGuiHoveredFlags for more options.
func IsItemFocused() bool                              { panic("not implemented") } // is the last item focused for keyboard/gamepad navigation?
func IsItemClicked(mouse_button ImGuiMouseButton) bool { panic("not implemented") } // is the last item hovered and mouse clicked on? (**)  == IsMouseClicked(mouse_button) && IsItemHovered()Important. (**) this it NOT equivalent to the behavior of e.g. Button(). Read comments in function definition.
func IsItemVisible() bool                              { panic("not implemented") } // is the last item visible? (items may be out of sight because of clipping/scrolling)
func IsItemEdited() bool                               { panic("not implemented") } // did the last item modify its underlying value this frame? or was pressed? This is generally the same as the "bool" return value of many widgets.
func IsItemActivated() bool                            { panic("not implemented") } // was the last item just made active (item was previously inactive).
func IsItemDeactivated() bool                          { panic("not implemented") } // was the last item just made inactive (item was previously active). Useful for Undo/Redo patterns with widgets that requires continuous editing.
func IsItemDeactivatedAfterEdit() bool                 { panic("not implemented") } // was the last item just made inactive and made a value change when it was active? (e.g. Slider/Drag moved). Useful for Undo/Redo patterns with widgets that requires continuous editing. Note that you may get false positives (some widgets such as Combo()/ListBox()/Selectable() will return true even when clicking an already selected item).
func IsItemToggledOpen() bool                          { panic("not implemented") } // was the last item open state toggled? set by TreeNode().
func IsAnyItemHovered() bool                           { panic("not implemented") } // is any item hovered?
func IsAnyItemActive() bool                            { panic("not implemented") } // is any item active?
func IsAnyItemFocused() bool                           { panic("not implemented") } // is any item focused?
func GetItemRectMin() ImVec2                           { panic("not implemented") } // get upper-left bounding rectangle of the last item (screen space)
func GetItemRectMax() ImVec2                           { panic("not implemented") } // get lower-right bounding rectangle of the last item (screen space)
func GetItemRectSize() ImVec2                          { panic("not implemented") } // get size of last item
func SetItemAllowOverlap()                             { panic("not implemented") } // allow last item to be overlapped by a subsequent item. sometimes useful with invisible buttons, selectables, etc. to catch unused area.

// Viewports
// - Currently represents the Platform Window created by the application which is hosting our Dear ImGui windows.
// - In 'docking' branch with multi-viewport enabled, we extend this concept to have multiple active viewports.
// - In the future we will extend this concept further to also represent Platform Monitor and support a "no main platform window" operation mode.
func GetMainViewport() *ImGuiViewport {
	var g = GImGui
	return g.Viewports[0]
} // return primary/default viewport. This can never be NULL.

// Miscellaneous Utilities
func IsRectVisible(size ImVec2) bool                     { panic("not implemented") } // test if rectangle (of given size, starting from cursor position) is visible / not clipped.
func IsRectVisibleMinMax(rect_min, rect_max ImVec2) bool { panic("not implemented") } // test if rectangle (in screen space) is visible / not clipped. to perform coarse clipping on user's side.
func GetTime() double                                    { panic("not implemented") } // get global imgui time. incremented by io.DeltaTime every frame.
func GetFrameCount() int                                 { panic("not implemented") } // get global imgui frame count. incremented by 1 every frame.
func GetBackgroundDrawList() *ImDrawList                 { panic("not implemented") } // this draw list will be the first rendering one. Useful to quickly draw shapes/text behind dear imgui contents.
func GetForegroundDrawList() *ImDrawList                 { panic("not implemented") } // this draw list will be the last rendered one. Useful to quickly draw shapes/text over dear imgui contents.
func GetDrawListSharedData() *ImDrawListSharedData       { panic("not implemented") } // you may use this when creating your own ImDrawList instances.
func GetStyleColorName(idx ImGuiCol) string              { panic("not implemented") } // get a string corresponding to the enum value (for display, saving, etc.).
func SetStateStorage(storage *ImGuiStorage)              { panic("not implemented") } // replace current window storage with our own (if you want to manipulate it yourself, typically clear subsection of it)
func GetStateStorage() *ImGuiStorage                     { panic("not implemented") }
func CalcListClipping(items_count int, items_height float, out_items_display_start *int, out_items_display_end *int) {
	panic("not implemented")
}                                                                          // calculate coarse clipping for large list of evenly sized items. Prefer using the ImGuiListClipper higher-level helper if you can.
func BeginChildFrame(id ImGuiID, size ImVec2, flsgs ImGuiWindowFlags) bool { panic("not implemented") } // helper to create a child window / scrolling region that looks like a normal widget frame
func EndChildFrame()                                                       { panic("not implemented") } // always call EndChildFrame() regardless of BeginChildFrame() return values (which indicates a collapsed/clipped window)

// Color Utilities
func ColorConvertU32ToFloat4(in ImU32) ImVec4 { panic("not implemented") }

func ColorConvertRGBtoHSV(r float, g float, b float, out_h, out_s, out_v *float) {
	panic("not implemented")
}
func ColorConvertHSVtoRGB(h float, s float, v float, out_r, out_g, out_b *float) {
	panic("not implemented")
}

// Inputs Utilities: Keyboard
// - For 'user_key_index int' you can use your own indices/enums according to how your backend/engine stored them in io.KeysDown[].
// - We don't know the meaning of those value. You can use GetKeyIndex() to map a ImGuiKey_ value into the user index.
func GetKeyIndex(imgui_key ImGuiKey) int { panic("not implemented") } // map ImGuiKey_* values into user's key index. == io.KeyMap[key]
func IsKeyDown(user_key_index int) bool  { panic("not implemented") } // is key being held. == io.KeysDown[user_key_index].

func IsKeyReleased(user_key_index int) bool                                 { panic("not implemented") } // was key released (went from Down to !Down)?
func GetKeyPressedAmount(key_index int, repeat_delay float, rate float) int { panic("not implemented") } // uses provided repeat rate/delay. return a count, most often 0 or 1 but might be >1 if RepeatRate is small enough that DeltaTime > RepeatRate
func CaptureKeyboardFromApp(want_capture_keyboard_value bool /*= true*/)    { panic("not implemented") } // attention: misleading name! manually override io.WantCaptureKeyboard flag next frame (said flag is entirely left for your application to handle). e.g. force capture keyboard when your widget is being hovered. This is equivalent to setting "io.WantCaptureKeyboard = want_capture_keyboard_value"  {panic("not implemented")} after the next NewFrame() call.

// Inputs Utilities: Mouse
// - To refer to a mouse button, you may use named enums in your code e.g. ImGuiMouseButton_Left, ImGuiMouseButton_Right.
// - You can also use regular integer: it is forever guaranteed that 0=Left, 1=Right, 2=Middle.
// - Dragging operations are only reported after mouse has moved a certain distance away from the initial clicking position (see 'lock_threshold' and 'io.MouseDraggingThreshold')
func IsMouseDown(button ImGuiMouseButton) bool                 { panic("not implemented") } // is mouse button held?
func IsMouseClicked(button ImGuiMouseButton, repeat bool) bool { panic("not implemented") } // did mouse button clicked? (went from !Down to Down)
func IsMouseReleased(button ImGuiMouseButton) bool             { panic("not implemented") } // did mouse button released? (went from Down to !Down)
func IsMouseDoubleClicked(button ImGuiMouseButton) bool        { panic("not implemented") } // did mouse button double-clicked? (note that a double-click will also report IsMouseClicked() == true)
func IsAnyMouseDown() bool                                     { panic("not implemented") } // is any mouse button held?
func GetMousePos() ImVec2                                      { panic("not implemented") } // shortcut to ImGui::GetIO().MousePos provided by user, to be consistent with other calls
func GetMousePosOnOpeningCurrentPopup() ImVec2                 { panic("not implemented") } // retrieve mouse position at the time of opening popup we have BeginPopup() into (helper to a user backing that value themselves)
func IsMouseDragging(button ImGuiMouseButton, lock_threshold float /*= -1.0*/) bool {
	panic("not implemented")
} // is mouse dragging? (if lock_threshold < -1.0, uses io.MouseDraggingThreshold)
func GetMouseDragDelta(button ImGuiMouseButton /*= 0*/, lock_threshold float /*= -1.0*/) ImVec2 {
	panic("not implemented")
}                                                                  // return the delta from the initial clicking position while the mouse button is pressed or was just released. This is locked and return 0.0 until the mouse moves past a distance threshold at least once (if lock_threshold < -1.0, uses io.MouseDraggingThreshold)
func ResetMouseDragDelta(button ImGuiMouseButton)                  { panic("not implemented") } //
func GetMouseCursor() ImGuiMouseCursor                             { panic("not implemented") } // get desired cursor type, reset in ImGui::NewFrame(), this is updated during the frame. valid before Render(). If you use software rendering by setting io.MouseDrawCursor ImGui will render those for you
func SetMouseCursor(cursor_type ImGuiMouseCursor)                  { panic("not implemented") } // set desired cursor type
func CaptureMouseFromApp(want_capture_mouse_value bool /*= true*/) { panic("not implemented") } // attention: misleading name! manually override io.WantCaptureMouse flag next frame (said flag is entirely left for your application to handle). This is equivalent to setting "io.WantCaptureMouse = want_capture_mouse_value  {panic("not implemented")}" after the next NewFrame() call.

// Clipboard Utilities
// - Also see the LogToClipboard() function to capture GUI into clipboard, or easily output text data to the clipboard.
func GetClipboardText() string     { panic("not implemented") }
func SetClipboardText(text string) { panic("not implemented") }

// Settings/.Ini Utilities
// - The disk functions are automatically called if io.IniFilename != NULL (default is "imgui.ini").
// - Set io.IniFilename to NULL to load/save manually. Read io.WantSaveIniSettings description about handling .ini saving manually.
// - Important: default value "imgui.ini" is relative to current working dir! Most apps will want to lock this to an absolute path (e.g. same path as executables).
func LoadIniSettingsFromDisk(ini_filename string)                 { panic("not implemented") } // call after CreateContext() and before the first call to NewFrame(). NewFrame() automatically calls LoadIniSettingsFromDisk(io.IniFilename).
func LoadIniSettingsFromMemory(ini_data string, ini_size uintptr) { panic("not implemented") } // call after CreateContext() and before the first call to NewFrame() to provide .ini data from your own data source.
func SaveIniSettingsToDisk(ini_filename string)                   { panic("not implemented") } // this is automatically called (if io.IniFilename is not empty) a few seconds after any modification that should be reflected in the .ini file (and also by DestroyContext).
func SaveIniSettingsToMemory(out_ini_size *uintptr) string        { panic("not implemented") } // return a zero-terminated string with the .ini data which you can save by your own mean. call when io.WantSaveIniSettings is set, then save data by your own mean and clear io.WantSaveIniSettings.

// Debug Utilities
// - This is used by the IMGUI_CHECKVERSION() macro.
func DebugCheckVersionAndDataLayout(version_str string, sz_io uintptr, sz_style uintptr, sz_vec2 uintptr, sz_vec4 uintptr, sz_drawvert uintptr, sz_drawidx uintptr) bool {
	panic("not implemented")
} // This is called by IMGUI_CHECKVERSION() macro.
