package imgui

import (
	"math"
)

// HelpMarker Helper to display a little (?) mark which shows a tooltip when hovered.
// In your own code you may want to display an actual icon if you are using a merged icon fonts (see docs/FONTS.md)
func HelpMarker(desc string) {
	TextDisabled(" (?)")

	// TODO: this is never true for some reason
	if IsItemHovered(0) {
		BeginTooltip()
		PushTextWrapPos(GetFontSize() * 35.0)
		TextUnformatted(desc)
		PopTextWrapPos()
		EndTooltip()
	}
}

// ShowUserGuide Helper to display basic user controls.
func ShowUserGuide() {
	var io = GetIO()
	BulletText("Double-click on title bar to collapse window.")
	BulletText(
		"Click and drag on lower corner to resize window\n" +
			"(double-click to auto fit window to its contents).")
	BulletText("CTRL+Click on a slider or drag box to input value as text.")
	BulletText("TAB/SHIFT+TAB to cycle through keyboard editable fields.")
	if io.FontAllowUserScaling {
		BulletText("CTRL+Mouse Wheel to zoom window contents.")
	}
	BulletText("While inputing text:\n")
	Indent(0)
	BulletText("CTRL+Left/Right to word jump.")
	BulletText("CTRL+A or double-click to select all.")
	BulletText("CTRL+X/C/V to use clipboard cut/copy/paste.")
	BulletText("CTRL+Z,CTRL+Y to undo/redo.")
	BulletText("ESCAPE to revert.")
	BulletText("You can apply arithmetic operators +,*,/ on numerical values.\nUse +- to subtract.")
	Unindent(0)
	BulletText("With keyboard navigation enabled:")
	Indent(0)
	BulletText("Arrow keys to navigate.")
	BulletText("Space to activate a widget.")
	BulletText("Return to input text into a widget.")
	BulletText("Escape to deactivate a widget, close popup, exit child window.")
	BulletText("Alt to jump to the menu layer of a window.")
	BulletText("CTRL+Tab to select a window.")
	Unindent(0)
}

var demoState struct {
	// Examples Apps (accessible from the "Examples" menu)
	show_app_main_menu_bar, show_app_documents,
	show_app_console, show_app_log, show_app_layout,
	show_app_property_editor, show_app_long_text,
	show_app_auto_resize, show_app_constrained_resize,
	show_app_simple_overlay, show_app_fullscreen,
	show_app_window_titles, show_app_custom_rendering bool

	// Dear ImGui Apps (accessible from the "Tools" menu)
	show_app_metrics, show_app_style_editor, show_app_about bool

	// Demonstrate the various window flags. Typically you would just use the default!
	no_titlebar, no_scrollbar, no_menu, no_move, no_resize,
	no_collapse, no_close, no_nav, no_background,
	no_bring_to_front, unsaved_document bool
}

// ShowDemoWindow Demonstrate most Dear ImGui features (this is big function!)
// You may execute this function to experiment with the UI and understand what it does.
// You may then search for keywords in the code when you are interested by a specific feature.
// create Demo window. demonstrate most ImGui features. call this to learn about the library! try to make it always available in your application!
func ShowDemoWindow(p_open *bool) {
	// Exceptionally add an extra assert here for people confused about initial Dear ImGui setup
	// Most ImGui functions would normally just crash if the context is missing.
	IM_ASSERT_USER_ERROR(GetCurrentContext() != nil, "Missing dear imgui context. Refer to examples app!")

	if demoState.show_app_main_menu_bar {
		//ShowExampleAppMainMenuBar();
	}
	if demoState.show_app_documents {
		//ShowExampleAppDocuments(&state.show_app_documents);
	}
	if demoState.show_app_console {
		//ShowExampleAppConsole(&show_app_console);
	}
	if demoState.show_app_log {
		//ShowExampleAppLog(&show_app_log);
	}
	if demoState.show_app_layout {
		//ShowExampleAppLayout(&show_app_layout);
	}
	if demoState.show_app_property_editor {
		//ShowExampleAppPropertyEditor(&show_app_property_editor);
	}
	if demoState.show_app_long_text {
		//ShowExampleAppLongText(&show_app_long_text);
	}
	if demoState.show_app_auto_resize {
		//ShowExampleAppAutoResize(&show_app_auto_resize);
	}
	if demoState.show_app_constrained_resize {
		//ShowExampleAppConstrainedResize(&show_app_constrained_resize);
	}
	if demoState.show_app_simple_overlay {
		//ShowExampleAppSimpleOverlay(&show_app_simple_overlay);
	}
	if demoState.show_app_fullscreen {
		//ShowExampleAppFullscreen(&show_app_fullscreen);
	}
	if demoState.show_app_window_titles {
		//ShowExampleAppWindowTitles(&show_app_window_titles);
	}
	if demoState.show_app_custom_rendering {
		//ShowExampleAppCustomRendering(&show_app_custom_rendering);
	}

	if demoState.show_app_metrics {
		ShowMetricsWindow(&demoState.show_app_metrics)
	}
	if demoState.show_app_about {
		ShowAboutWindow(&demoState.show_app_about)
	}
	if demoState.show_app_style_editor {
		Begin("Dear ImGui Style Editor", &demoState.show_app_style_editor, 0)
		ShowStyleEditor(nil)
		End()
	}

	var window_flags ImGuiWindowFlags = 0
	if demoState.no_titlebar {
		window_flags |= ImGuiWindowFlags_NoTitleBar
	}
	if demoState.no_scrollbar {
		window_flags |= ImGuiWindowFlags_NoScrollbar
	}
	if !demoState.no_menu {
		window_flags |= ImGuiWindowFlags_MenuBar
	}
	if demoState.no_move {
		window_flags |= ImGuiWindowFlags_NoMove
	}
	if demoState.no_resize {
		window_flags |= ImGuiWindowFlags_NoResize
	}
	if demoState.no_collapse {
		window_flags |= ImGuiWindowFlags_NoCollapse
	}
	if demoState.no_nav {
		window_flags |= ImGuiWindowFlags_NoNav
	}
	if demoState.no_background {
		window_flags |= ImGuiWindowFlags_NoBackground
	}
	if demoState.no_bring_to_front {
		window_flags |= ImGuiWindowFlags_NoBringToFrontOnFocus
	}
	if demoState.unsaved_document {
		window_flags |= ImGuiWindowFlags_UnsavedDocument
	}
	if demoState.no_close {
		p_open = nil // Don't pass our bool* to Begin
	}

	// We specify a default position/size in case there's no data in the .ini file.
	// We only do it to make the demo applications a little more welcoming, but typically this isn't required.
	var main_viewport = GetMainViewport()
	SetNextWindowPos(NewImVec2(main_viewport.WorkPos.X()+650, main_viewport.WorkPos.Y()+20), ImGuiCond_FirstUseEver, ImVec2{})
	SetNextWindowSize(NewImVec2(550, 680), ImGuiCond_FirstUseEver)

	// Main body of the Demo window starts here.
	if !Begin("Dear ImGui Demo", p_open, window_flags) {
		// Early out if the window is collapsed, as an optimization.
		End()
		return
	}

	// Most "big" widgets share a common width settings by default. See 'Demo->Layout->Widgets Width' for details.

	// e.g. Use 2/3 of the space for widgets and 1/3 for labels (right align)
	//PushItemWidth(-GetWindowWidth() * 0.35f);

	// e.g. Leave a fixed amount of width for labels (by passing a negative value), the rest goes to widgets.
	PushItemWidth(GetFontSize() * -12)

	// Menu Bar
	if BeginMenuBar() {
		if BeginMenu("Menu", true) {
			//ShowExampleMenuFile()
			EndMenu()
		}
		if BeginMenu("Examples", true) {
			MenuItem("Main menu bar", "", &demoState.show_app_main_menu_bar, true)
			MenuItem("Console", "", &demoState.show_app_console, true)
			MenuItem("Log", "", &demoState.show_app_log, true)
			MenuItem("Simple layout", "", &demoState.show_app_layout, true)
			MenuItem("Property editor", "", &demoState.show_app_property_editor, true)
			MenuItem("Long text display", "", &demoState.show_app_long_text, true)
			MenuItem("Auto-resizing window", "", &demoState.show_app_auto_resize, true)
			MenuItem("Constrained-resizing window", "", &demoState.show_app_constrained_resize, true)
			MenuItem("Simple overlay", "", &demoState.show_app_simple_overlay, true)
			MenuItem("Fullscreen window", "", &demoState.show_app_fullscreen, true)
			MenuItem("Manipulating window titles", "", &demoState.show_app_window_titles, true)
			MenuItem("Custom rendering", "", &demoState.show_app_custom_rendering, true)
			MenuItem("Documents", "", &demoState.show_app_documents, true)
			EndMenu()
		}
		if BeginMenu("Tools", true) {
			MenuItem("Metrics/Debugger", "", &demoState.show_app_metrics, true)
			MenuItem("Style Editor", "", &demoState.show_app_style_editor, true)
			MenuItem("About Dear ImGui", "", &demoState.show_app_about, true)
			EndMenu()
		}
		EndMenuBar()
	}

	Text("dear imgui says hello. (%s)", IMGUI_VERSION)
	Spacing()

	if CollapsingHeader("Help", 0) {
		Text("ABOUT THIS DEMO:")
		BulletText("Sections below are demonstrating many aspects of the library.")
		BulletText("The \"Examples\" menu above leads to more demo contents.")
		BulletText("The \"Tools\" menu above gives access to: About Box, Style Editor,\n" +
			"and Metrics/Debugger (general purpose Dear ImGui debugging tool).")
		Separator()

		Text("PROGRAMMER GUIDE:")
		BulletText("See the ShowDemoWindow() code in imgui_demo.cpp. <- you are here!")
		BulletText("See comments in cpp.")
		BulletText("See example applications in the examples/ folder.")
		BulletText("Read the FAQ at http://www.dearorg/faq/")
		BulletText("Set 'io.ConfigFlags |= NavEnableKeyboard' for keyboard controls.")
		BulletText("Set 'io.ConfigFlags |= NavEnableGamepad' for gamepad controls.")
		Separator()

		Text("USER GUIDE:")
		ShowUserGuide()
	}

	if CollapsingHeader("Configuration", 0) {
		var io = GetIO()

		if TreeNode("Configuration##2") {
			CheckboxFlagsInt("io.ConfigFlags: NavEnableKeyboard", (*int32)(&io.ConfigFlags), int32(ImGuiConfigFlags_NavEnableKeyboard))
			SameLine(0, 0)
			HelpMarker("Enable keyboard controls.")
			CheckboxFlagsInt("io.ConfigFlags: NavEnableGamepad", (*int32)(&io.ConfigFlags), int32(ImGuiConfigFlags_NavEnableGamepad))
			SameLine(0, 0)
			HelpMarker("Enable gamepad controls. Require backend to set io.BackendFlags |= ImGuiBackendFlags_HasGamepad.\n\nRead instructions in cpp for details.")
			CheckboxFlagsInt("io.ConfigFlags: NavEnableSetMousePos", (*int32)(&io.ConfigFlags), int32(ImGuiConfigFlags_NavEnableSetMousePos))
			SameLine(0, 0)
			HelpMarker("Instruct navigation to move the mouse cursor. See comment for ImGuiConfigFlags_NavEnableSetMousePos.")
			CheckboxFlagsInt("io.ConfigFlags: NoMouse", (*int32)(&io.ConfigFlags), int32(ImGuiConfigFlags_NoMouse))
			if io.ConfigFlags&ImGuiConfigFlags_NoMouse != 0 {
				// The "NoMouse" option can get us stuck with a disabled mouse! Let's provide an alternative way to fix it:
				if math.Mod((float64)(GetTime()), 0.40) < 0.20 {
					SameLine(0, 0)
					Text("<<PRESS SPACE TO DISABLE>>")
				}
				if IsKeyPressed(GetKeyIndex(ImGuiKey_Space), true) {
					io.ConfigFlags &= ^ImGuiConfigFlags_NoMouse
				}
			}
			CheckboxFlagsInt("io.ConfigFlags: NoMouseCursorChange", (*int32)(&io.ConfigFlags), int32(ImGuiConfigFlags_NoMouseCursorChange))
			SameLine(0, 0)
			HelpMarker("Instruct backend to not alter mouse cursor shape and visibility.")
			Checkbox("io.ConfigInputTextCursorBlink", &io.ConfigInputTextCursorBlink)
			SameLine(0, 0)
			HelpMarker("Enable blinking cursor (optional as some users consider it to be distracting)")
			Checkbox("io.ConfigDragClickToInputText", &io.ConfigDragClickToInputText)
			SameLine(0, 0)
			HelpMarker("Enable turning DragXXX widgets into text input with a simple mouse click-release (without moving).")
			Checkbox("io.ConfigWindowsResizeFromEdges", &io.ConfigWindowsResizeFromEdges)
			SameLine(0, 0)
			HelpMarker("Enable resizing of windows from their edges and from the lower-left corner.\nThis requires (io.BackendFlags & ImGuiBackendFlags_HasMouseCursors) because it needs mouse cursor feedback.")
			Checkbox("io.ConfigWindowsMoveFromTitleBarOnly", &io.ConfigWindowsMoveFromTitleBarOnly)
			Checkbox("io.MouseDrawCursor", &io.MouseDrawCursor)
			SameLine(0, 0)
			HelpMarker("Instruct Dear ImGui to render a mouse cursor itself. Note that a mouse cursor rendered via your application GPU rendering path will feel more laggy than hardware cursor, but will be more in sync with your other visuals.\n\nSome desktop applications may use both kinds of cursors (e.g. enable software cursor only when resizing/dragging something).")
			Text("Also see Style->Rendering for rendering options.")
			TreePop()
			Separator()
		}

		if TreeNode("Backend Flags") {
			HelpMarker(
				"Those flags are set by the backends (imgui_impl_xxx files) to specify their capabilities.\n" +
					"Here we expose them as read-only fields to avoid breaking interactions with your backend.")

			// Make a local copy to avoid modifying actual backend flags.
			var backend_flags = io.BackendFlags
			CheckboxFlagsInt("io.BackendFlags: HasGamepad", (*int32)(&backend_flags), int32(ImGuiBackendFlags_HasGamepad))
			CheckboxFlagsInt("io.BackendFlags: HasMouseCursors", (*int32)(&backend_flags), int32(ImGuiBackendFlags_HasMouseCursors))
			CheckboxFlagsInt("io.BackendFlags: HasSetMousePos", (*int32)(&backend_flags), int32(ImGuiBackendFlags_HasSetMousePos))
			CheckboxFlagsInt("io.BackendFlags: RendererHasVtxOffset", (*int32)(&backend_flags), int32(ImGuiBackendFlags_RendererHasVtxOffset))
			TreePop()
			Separator()
		}

		if TreeNode("Style") {
			HelpMarker("The same contents can be accessed in 'Tools->Style Editor' or by calling the ShowStyleEditor() function.")
			ShowStyleEditor(nil)
			TreePop()
			Separator()
		}

		if TreeNode("Capture/Logging") {
			HelpMarker(
				"The logging API redirects all text output so you can easily capture the content of " +
					"a window or a block. Tree nodes can be automatically expanded.\n" +
					"Try opening any of the contents below in this window and then click one of the \"Log To\" button.")
			LogButtons()

			HelpMarker("You can also call LogText() to output directly to the log without a visual output.")
			if Button("Copy \"Hello, world!\" to clipboard") {
				LogToClipboard(-1)
				LogText("Hello, world!")
				LogFinish()
			}
			TreePop()
		}
	}

	if CollapsingHeader("Window options", 0) {
		if BeginTable("split", 3, 0, ImVec2{}, 0) {
			TableNextColumn()
			Checkbox("No titlebar", &demoState.no_titlebar)
			TableNextColumn()
			Checkbox("No scrollbar", &demoState.no_scrollbar)
			TableNextColumn()
			Checkbox("No menu", &demoState.no_menu)
			TableNextColumn()
			Checkbox("No move", &demoState.no_move)
			TableNextColumn()
			Checkbox("No resize", &demoState.no_resize)
			TableNextColumn()
			Checkbox("No collapse", &demoState.no_collapse)
			TableNextColumn()
			Checkbox("No close", &demoState.no_close)
			TableNextColumn()
			Checkbox("No nav", &demoState.no_nav)
			TableNextColumn()
			Checkbox("No background", &demoState.no_background)
			TableNextColumn()
			Checkbox("No bring to front", &demoState.no_bring_to_front)
			TableNextColumn()
			Checkbox("Unsaved document", &demoState.unsaved_document)
			EndTable()
		}
	}

	// All demo contents
	//ShowDemoWindowWidgets()
	//ShowDemoWindowLayout()
	//ShowDemoWindowPopups()
	//ShowDemoWindowTables()
	//ShowDemoWindowMisc()

	// End of ShowDemoWindow()
	PopItemWidth()
	End()
}
