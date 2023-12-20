package imgui

// access the IO structure (mouse/keyboard/gamepad inputs, time, various configuration options/flags)
func GetIO() *ImGuiIO {
	IM_ASSERT_USER_ERROR(GImGui != nil, "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?")
	return &GImGui.IO
}

// Demo, Debug, Information
func ShowAboutWindow(p_open *bool)        { panic("not implemented") } // create About window. display Dear ImGui version, credits and build/system information.
func ShowStyleEditor(ref *ImGuiStyle)     { panic("not implemented") } // add style editor block (not a window). you can pass in a reference ImGuiStyle structure to compare to, revert to and save to (else it uses the default style)
func ShowStyleSelector(label string) bool { panic("not implemented") } // add style selector block (not a window), essentially a combo listing the default styles.
func ShowFontSelector(label string)       { panic("not implemented") } // add font selector block (not a window), essentially a combo listing the loaded fonts.

// get the compiled version string e.g. "1.80 WIP" (essentially the value for IMGUI_VERSION from the compiled version of imgui.cpp)
func GetVersion() string {
	return IMGUI_VERSION
}

// Clipping
// - Mouse hovering is affected by ImGui::PushClipRect() calls, unlike direct calls to ImDrawList::PushClipRect() which are render only.
// Push a clipping rectangle for both ImGui logic (hit-testing etc.) and low-level ImDrawList rendering.
//   - When using this function it is sane to ensure that float are perfectly rounded to integer values,
//     so that e.g. (int)(max.x-min.x) in user's render produce correct result.
//   - If the code here changes, may need to update code of functions like NextColumn() and PushColumnClipRect():
//     some frequently called functions which to modify both channels and clipping simultaneously tend to use the
//     more specialized SetWindowClipRectBeforeSetChannel() to avoid extraneous updates of underlying ImDrawCmds.
func PushClipRect(cr_min ImVec2, cr_max ImVec2, intersect_with_current_clip_rect bool) {
	var window = GetCurrentWindow()
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
