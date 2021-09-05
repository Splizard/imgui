package demo

import "github.com/splizard/imgui"

// Helper to display a little (?) mark which shows a tooltip when hovered.
// In your own code you may want to display an actual icon if you are using a merged icon fonts (see docs/FONTS.md)
func HelpMarker(desc string) {
	imgui.TextDisabled("(?)")
	if imgui.IsItemHovered() {
		imgui.BeginTooltip()
		imgui.PushTextWrapPos(imgui.GetFontSize() * 35.0)
		imgui.TextUnformatted(desc)
		imgui.PopTextWrapPos()
		imgui.EndTooltip()
	}
}

// Helper to display basic user controls.
func ShowUserGuide() {
	var io = imgui.GetIO()
	imgui.BulletText("Double-click on title bar to collapse window.")
	imgui.BulletText(
		"Click and drag on lower corner to resize window\n" +
			"(double-click to auto fit window to its contents).")
	imgui.BulletText("CTRL+Click on a slider or drag box to input value as text.")
	imgui.BulletText("TAB/SHIFT+TAB to cycle through keyboard editable fields.")
	if io.FontAllowUserScaling {
		imgui.BulletText("CTRL+Mouse Wheel to zoom window contents.")
	}
	imgui.BulletText("While inputing text:\n")
	imgui.Indent()
	imgui.BulletText("CTRL+Left/Right to word jump.")
	imgui.BulletText("CTRL+A or double-click to select all.")
	imgui.BulletText("CTRL+X/C/V to use clipboard cut/copy/paste.")
	imgui.BulletText("CTRL+Z,CTRL+Y to undo/redo.")
	imgui.BulletText("ESCAPE to revert.")
	imgui.BulletText("You can apply arithmetic operators +,*,/ on numerical values.\nUse +- to subtract.")
	imgui.Unindent()
	imgui.BulletText("With keyboard navigation enabled:")
	imgui.Indent()
	imgui.BulletText("Arrow keys to navigate.")
	imgui.BulletText("Space to activate a widget.")
	imgui.BulletText("Return to input text into a widget.")
	imgui.BulletText("Escape to deactivate a widget, close popup, exit child window.")
	imgui.BulletText("Alt to jump to the menu layer of a window.")
	imgui.BulletText("CTRL+Tab to select a window.")
	imgui.Unindent()
}

var (
	// Examples Apps (accessible from the "Examples" menu)
	show_app_main_menu_bar      = false
	show_app_documents          = false
	show_app_console            = false
	show_app_log                = false
	show_app_layout             = false
	show_app_property_editor    = false
	show_app_long_text          = false
	show_app_auto_resize        = false
	show_app_constrained_resize = false
	show_app_simple_overlay     = false
	show_app_fullscreen         = false
	show_app_window_titles      = false
	show_app_custom_rendering   = false

	// Dear ImGui Apps (accessible from the "Tools" menu)
	show_app_metrics      = false
	show_app_style_editor = false
	show_app_about        = false

	// Demonstrate the various window flags. Typically you would just use the default!
	no_titlebar       = false
	no_scrollbar      = false
	no_menu           = false
	no_move           = false
	no_resize         = false
	no_collapse       = false
	no_close          = false
	no_nav            = false
	no_background     = false
	no_bring_to_front = false
	unsaved_document  = false
)

// Demonstrate most Dear ImGui features (this is big function!)
// You may execute this function to experiment with the UI and understand what it does.
// You may then search for keywords in the code when you are interested by a specific feature.
func ShowDemoWindow(p_open *bool) {
	// Exceptionally add an extra assert here for people confused about initial Dear ImGui setup
	// Most ImGui functions would normally just crash if the context is missing.
	IM_ASSERT(imgui.GetCurrentContext() != NULL && "Missing dear imgui context. Refer to examples app!")

	if show_app_main_menu_bar {
		ShowExampleAppMainMenuBar()
	}
	if show_app_documents {
		ShowExampleAppDocuments(&show_app_documents)
	}
	if show_app_console {
		ShowExampleAppConsole(&show_app_console)
	}

	if show_app_log {
		ShowExampleAppLog(&show_app_log)
	}
	if show_app_layout {
		ShowExampleAppLayout(&show_app_layout)
	}
	if show_app_property_editor {
		ShowExampleAppPropertyEditor(&show_app_property_editor)
	}
	if show_app_long_text {
		ShowExampleAppLongText(&show_app_long_text)
	}
	if show_app_auto_resize {
		ShowExampleAppAutoResize(&show_app_auto_resize)
	}
	if show_app_constrained_resize {
		ShowExampleAppConstrainedResize(&show_app_constrained_resize)
	}
	if show_app_simple_overlay {
		ShowExampleAppSimpleOverlay(&show_app_simple_overlay)
	}
	if show_app_fullscreen {
		ShowExampleAppFullscreen(&show_app_fullscreen)
	}
	if show_app_window_titles {
		ShowExampleAppWindowTitles(&show_app_window_titles)
	}
	if show_app_custom_rendering {
		ShowExampleAppCustomRendering(&show_app_custom_rendering)
	}

	if show_app_metrics {
		imgui.ShowMetricsWindow(&show_app_metrics)
	}
	if show_app_about {
		imgui.ShowAboutWindow(&show_app_about)
	}
	if show_app_style_editor {
		imgui.Begin("Dear ImGui Style Editor", &show_app_style_editor)
		imgui.ShowStyleEditor()
		imgui.End()
	}

	var window_flags ImGuiWindowFlags
	if no_titlebar {
		window_flags |= ImGuiWindowFlags_NoTitleBar
	}
	if no_scrollbar {
		window_flags |= ImGuiWindowFlags_NoScrollbar
	}
	if !no_menu {
		window_flags |= ImGuiWindowFlags_MenuBar
	}
	if no_move {
		window_flags |= ImGuiWindowFlags_NoMove
	}
	if no_resize {
		window_flags |= ImGuiWindowFlags_NoResize
	}
	if no_collapse {
		window_flags |= ImGuiWindowFlags_NoCollapse
	}
	if no_nav {
		window_flags |= ImGuiWindowFlags_NoNav
	}
	if no_background {
		window_flags |= ImGuiWindowFlags_NoBackground
	}
	if no_bring_to_front {
		window_flags |= ImGuiWindowFlags_NoBringToFrontOnFocus
	}
	if unsaved_document {
		window_flags |= ImGuiWindowFlags_UnsavedDocument
	}
	if no_close {
		p_open = nil // Don't pass our bool* to Begin
	}

	// We specify a default position/size in case there's no data in the .ini file.
	// We only do it to make the demo applications a little more welcoming, but typically this isn't required.
	var main_viewport = imgui.GetMainViewport()
	imgui.SetNextWindowPos(ImVec2(main_viewport.WorkPos.X+650, main_viewport.WorkPos.Y+20), ImGuiCond_FirstUseEver)
	imgui.SetNextWindowSize(ImVec2(550, 680), ImGuiCond_FirstUseEver)

	// Main body of the Demo window starts here.
	if !imgui.Begin("Dear ImGui Demo", p_open, window_flags) {
		// Early out if the window is collapsed, as an optimization.
		imgui.End()
		return
	}

	// Most "big" widgets share a common width settings by default. See 'Demo->Layout->Widgets Width' for details.

	// e.g. Use 2/3 of the space for widgets and 1/3 for labels (right align)
	//imgui.PushItemWidth(-imgui.GetWindowWidth() * 0.35f);

	// e.g. Leave a fixed amount of width for labels (by passing a negative value), the rest goes to widgets.
	imgui.PushItemWidth(imgui.GetFontSize() * -12)

	// Menu Bar
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("Menu") {
			ShowExampleMenuFile()
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Examples") {
			imgui.MenuItem("Main menu bar", NULL, &show_app_main_menu_bar)
			imgui.MenuItem("Console", NULL, &show_app_console)
			imgui.MenuItem("Log", NULL, &show_app_log)
			imgui.MenuItem("Simple layout", NULL, &show_app_layout)
			imgui.MenuItem("Property editor", NULL, &show_app_property_editor)
			imgui.MenuItem("Long text display", NULL, &show_app_long_text)
			imgui.MenuItem("Auto-resizing window", NULL, &show_app_auto_resize)
			imgui.MenuItem("Constrained-resizing window", NULL, &show_app_constrained_resize)
			imgui.MenuItem("Simple overlay", NULL, &show_app_simple_overlay)
			imgui.MenuItem("Fullscreen window", NULL, &show_app_fullscreen)
			imgui.MenuItem("Manipulating window titles", NULL, &show_app_window_titles)
			imgui.MenuItem("Custom rendering", NULL, &show_app_custom_rendering)
			imgui.MenuItem("Documents", NULL, &show_app_documents)
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Tools") {
			imgui.MenuItem("Metrics/Debugger", NULL, &show_app_metrics)
			imgui.MenuItem("Style Editor", NULL, &show_app_style_editor)
			imgui.MenuItem("About Dear ImGui", NULL, &show_app_about)
			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}

	imgui.Text("dear imgui says hello. (%s)", IMGUI_VERSION)
	imgui.Spacing()

	if imgui.CollapsingHeader("Help") {
		imgui.Text("ABOUT THIS DEMO:")
		imgui.BulletText("Sections below are demonstrating many aspects of the library.")
		imgui.BulletText("The \"Examples\" menu above leads to more demo contents.")
		imgui.BulletText("The \"Tools\" menu above gives access to: About Box, Style Editor,\n" +
			"and Metrics/Debugger (general purpose Dear ImGui debugging tool).")
		imgui.Separator()

		imgui.Text("PROGRAMMER GUIDE:")
		imgui.BulletText("See the ShowDemoWindow() code in imgui_demo.cpp. <- you are here!")
		imgui.BulletText("See comments in imgui.cpp.")
		imgui.BulletText("See example applications in the examples/ folder.")
		imgui.BulletText("Read the FAQ at http://www.dearimgui.org/faq/")
		imgui.BulletText("Set 'io.ConfigFlags |= NavEnableKeyboard' for keyboard controls.")
		imgui.BulletText("Set 'io.ConfigFlags |= NavEnableGamepad' for gamepad controls.")
		imgui.Separator()

		imgui.Text("USER GUIDE:")
		imgui.ShowUserGuide()
	}

	if imgui.CollapsingHeader("Configuration") {
		ImGuiIO & io = imgui.GetIO()

		if imgui.TreeNode("Configuration##2") {
			imgui.CheckboxFlags("io.ConfigFlags: NavEnableKeyboard", &io.ConfigFlags, ImGuiConfigFlags_NavEnableKeyboard)
			imgui.SameLine()
			HelpMarker("Enable keyboard controls.")
			imgui.CheckboxFlags("io.ConfigFlags: NavEnableGamepad", &io.ConfigFlags, ImGuiConfigFlags_NavEnableGamepad)
			imgui.SameLine()
			HelpMarker("Enable gamepad controls. Require backend to set io.BackendFlags |= ImGuiBackendFlags_HasGamepad.\n\nRead instructions in imgui.cpp for details.")
			imgui.CheckboxFlags("io.ConfigFlags: NavEnableSetMousePos", &io.ConfigFlags, ImGuiConfigFlags_NavEnableSetMousePos)
			imgui.SameLine()
			HelpMarker("Instruct navigation to move the mouse cursor. See comment for ImGuiConfigFlags_NavEnableSetMousePos.")
			imgui.CheckboxFlags("io.ConfigFlags: NoMouse", &io.ConfigFlags, ImGuiConfigFlags_NoMouse)
			if io.ConfigFlags & ImGuiConfigFlags_NoMouse {
				// The "NoMouse" option can get us stuck with a disabled mouse! Let's provide an alternative way to fix it:
				if fmodf(float32(imgui.GetTime()), 0.40) < 0.20 {
					imgui.SameLine()
					imgui.Text("<<PRESS SPACE TO DISABLE>>")
				}
				if imgui.IsKeyPressed(imgui.GetKeyIndex(ImGuiKey_Space)) {
					io.ConfigFlags &= ^ImGuiConfigFlags_NoMouse
				}
			}
			imgui.CheckboxFlags("io.ConfigFlags: NoMouseCursorChange", &io.ConfigFlags, ImGuiConfigFlags_NoMouseCursorChange)
			imgui.SameLine()
			HelpMarker("Instruct backend to not alter mouse cursor shape and visibility.")
			imgui.Checkbox("io.ConfigInputTextCursorBlink", &io.ConfigInputTextCursorBlink)
			imgui.SameLine()
			HelpMarker("Enable blinking cursor (optional as some users consider it to be distracting)")
			imgui.Checkbox("io.ConfigDragClickToInputText", &io.ConfigDragClickToInputText)
			imgui.SameLine()
			HelpMarker("Enable turning DragXXX widgets into text input with a simple mouse click-release (without moving).")
			imgui.Checkbox("io.ConfigWindowsResizeFromEdges", &io.ConfigWindowsResizeFromEdges)
			imgui.SameLine()
			HelpMarker("Enable resizing of windows from their edges and from the lower-left corner.\nThis requires (io.BackendFlags & ImGuiBackendFlags_HasMouseCursors) because it needs mouse cursor feedback.")
			imgui.Checkbox("io.ConfigWindowsMoveFromTitleBarOnly", &io.ConfigWindowsMoveFromTitleBarOnly)
			imgui.Checkbox("io.MouseDrawCursor", &io.MouseDrawCursor)
			imgui.SameLine()
			HelpMarker("Instruct Dear ImGui to render a mouse cursor itself. Note that a mouse cursor rendered via your application GPU rendering path will feel more laggy than hardware cursor, but will be more in sync with your other visuals.\n\nSome desktop applications may use both kinds of cursors (e.g. enable software cursor only when resizing/dragging something).")
			imgui.Text("Also see Style->Rendering for rendering options.")
			imgui.TreePop()
			imgui.Separator()
		}

		if imgui.TreeNode("Backend Flags") {
			HelpMarker(
				"Those flags are set by the backends (imgui_impl_xxx files) to specify their capabilities.\n" +
					"Here we expose them as read-only fields to avoid breaking interactions with your backend.")

			// Make a local copy to avoid modifying actual backend flags.
			var backend_flags = io.BackendFlags
			imgui.CheckboxFlags("io.BackendFlags: HasGamepad", &backend_flags, ImGuiBackendFlags_HasGamepad)
			imgui.CheckboxFlags("io.BackendFlags: HasMouseCursors", &backend_flags, ImGuiBackendFlags_HasMouseCursors)
			imgui.CheckboxFlags("io.BackendFlags: HasSetMousePos", &backend_flags, ImGuiBackendFlags_HasSetMousePos)
			imgui.CheckboxFlags("io.BackendFlags: RendererHasVtxOffset", &backend_flags, ImGuiBackendFlags_RendererHasVtxOffset)
			imgui.TreePop()
			imgui.Separator()
		}

		if imgui.TreeNode("Style") {
			HelpMarker("The same contents can be accessed in 'Tools->Style Editor' or by calling the ShowStyleEditor() function.")
			imgui.ShowStyleEditor()
			imgui.TreePop()
			imgui.Separator()
		}

		if imgui.TreeNode("Capture/Logging") {
			HelpMarker(
				"The logging API redirects all text output so you can easily capture the content of " +
					"a window or a block. Tree nodes can be automatically expanded.\n" +
					"Try opening any of the contents below in this window and then click one of the \"Log To\" button.")
			imgui.LogButtons()

			HelpMarker("You can also call imgui.LogText() to output directly to the log without a visual output.")
			if imgui.Button("Copy \"Hello, world!\" to clipboard") {
				imgui.LogToClipboard()
				imgui.LogText("Hello, world!")
				imgui.LogFinish()
			}
			imgui.TreePop()
		}
	}

	if imgui.CollapsingHeader("Window options") {
		if imgui.BeginTable("split", 3) {
			imgui.TableNextColumn()
			imgui.Checkbox("No titlebar", &no_titlebar)
			imgui.TableNextColumn()
			imgui.Checkbox("No scrollbar", &no_scrollbar)
			imgui.TableNextColumn()
			imgui.Checkbox("No menu", &no_menu)
			imgui.TableNextColumn()
			imgui.Checkbox("No move", &no_move)
			imgui.TableNextColumn()
			imgui.Checkbox("No resize", &no_resize)
			imgui.TableNextColumn()
			imgui.Checkbox("No collapse", &no_collapse)
			imgui.TableNextColumn()
			imgui.Checkbox("No close", &no_close)
			imgui.TableNextColumn()
			imgui.Checkbox("No nav", &no_nav)
			imgui.TableNextColumn()
			imgui.Checkbox("No background", &no_background)
			imgui.TableNextColumn()
			imgui.Checkbox("No bring to front", &no_bring_to_front)
			imgui.TableNextColumn()
			imgui.Checkbox("Unsaved document", &unsaved_document)
			imgui.EndTable()
		}
	}

	// All demo contents
	ShowDemoWindowWidgets()
	ShowDemoWindowLayout()
	ShowDemoWindowPopups()
	ShowDemoWindowTables()
	ShowDemoWindowMisc()

	// End of ShowDemoWindow()
	imgui.PopItemWidth()
	imgui.End()
}
