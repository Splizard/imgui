package imgui

import (
	"fmt"
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

	// Layout & Scrolling
	show_scrolling_decoration bool
	show_scrolling_track      bool
	scrolling_track_item      int
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
	ShowDemoWindowWidgets()
	ShowDemoWindowLayout()
	ShowDemoWindowPopups()
	//ShowDemoWindowTables()
	ShowDemoWindowMisc()

	// End of ShowDemoWindow()
	PopItemWidth()
	End()
}

// ImVec4FromHSV converts HSV values to an ImVec4 color
func ImVec4FromHSV(h, s, v float) ImVec4 {
	var r, g, b float
	ColorConvertHSVtoRGB(h, s, v, &r, &g, &b)
	return ImVec4{r, g, b, 1.0}
}

// State for ShowDemoWindowWidgets
var widgetsState struct {
	// Basic
	clicked       int
	check         bool
	e             int
	counter       int
	str0          []byte
	str1          []byte
	i0            int
	f0            float
	f1            float
	vec4a         [4]float
	i1, i2        int
	f1drag, f2drag float
	i1slider      int
	f1slider      float
	f2slider      float
	angle         float
	elem          int
	col1          [3]float
	col2          [4]float
	item_current_combo int
	item_current_list  int

	// Trees
	base_flags                        ImGuiTreeNodeFlags
	align_label_with_current_x_position bool
	test_drag_and_drop                bool
	selection_mask                    int

	// Collapsing Headers
	closable_group bool

	// Text
	wrap_width float

	// Images
	pressed_count int

	// Combo (advanced)
	combo_flags          ImGuiComboFlags
	combo_item_current_idx int

	// List boxes
	listbox_item_current_idx int

	// Selectables
	selectable_basic     [5]bool
	selectable_single    int
	selectable_multiple  [5]bool
	selectable_render    [3]bool
	selectable_columns   [10]bool
	selectable_grid      [4][4]bool

	// Tabs
	tab_bar_flags    ImGuiTabBarFlags
	tabs_opened      [4]bool

	// Plots
	plots_animate     bool
	plots_values      [90]float
	plots_offset      int
	plots_refresh     float64
	plots_phase       float
	plots_progress    float
	plots_progress_dir float

	// Color/Picker Widgets
	color_picker_color ImVec4
	color_alpha_preview bool
	color_alpha_half_preview bool
	color_drag_and_drop bool
	color_options_menu bool
	color_hdr bool
	color_no_border bool
	color_alpha bool
	color_alpha_bar bool
	color_side_preview bool
	color_ref_color bool
	color_ref_color_v ImVec4
	color_display_mode int
	color_picker_mode int
	color_hsv ImVec4

	// Drag/Slider Flags
	slider_flags ImGuiSliderFlags
	drag_f float
	drag_i int
	slider_f float
	slider_i int

	// Range Widgets
	range_begin float
	range_end float
	range_begin_i int
	range_end_i int

	// Disable all
	disable_all bool
}

func init() {
	// Initialize default values
	widgetsState.check = true
	widgetsState.str0 = []byte("Hello, world!")
	widgetsState.str1 = []byte{}
	widgetsState.i0 = 123
	widgetsState.f0 = 0.001
	widgetsState.f1 = 1.0e10
	widgetsState.vec4a = [4]float{0.10, 0.20, 0.30, 0.44}
	widgetsState.i1 = 50
	widgetsState.i2 = 42
	widgetsState.f1drag = 1.00
	widgetsState.f2drag = 0.0067
	widgetsState.f1slider = 0.123
	widgetsState.col1 = [3]float{1.0, 0.0, 0.2}
	widgetsState.col2 = [4]float{0.4, 0.7, 0.0, 0.5}
	widgetsState.item_current_list = 1
	widgetsState.base_flags = ImGuiTreeNodeFlags_OpenOnArrow | ImGuiTreeNodeFlags_OpenOnDoubleClick | ImGuiTreeNodeFlags_SpanAvailWidth
	widgetsState.selection_mask = 1 << 2
	widgetsState.closable_group = true
	widgetsState.wrap_width = 200.0

	// Selectables
	widgetsState.selectable_basic[1] = true
	widgetsState.selectable_single = -1

	// Tabs
	widgetsState.tab_bar_flags = ImGuiTabBarFlags_Reorderable
	widgetsState.tabs_opened = [4]bool{true, true, true, true}

	// Plots
	widgetsState.plots_animate = true
	widgetsState.plots_progress_dir = 1.0

	// Color/Picker Widgets
	widgetsState.color_picker_color = ImVec4{114.0 / 255.0, 144.0 / 255.0, 154.0 / 255.0, 200.0 / 255.0}
	widgetsState.color_alpha_preview = true
	widgetsState.color_drag_and_drop = true
	widgetsState.color_options_menu = true
	widgetsState.color_alpha = true
	widgetsState.color_alpha_bar = true
	widgetsState.color_side_preview = true
	widgetsState.color_ref_color_v = ImVec4{1.0, 0.0, 1.0, 0.5}
	widgetsState.color_hsv = ImVec4{0.23, 1.0, 1.0, 1.0}

	// Drag/Slider Flags
	widgetsState.drag_f = 0.5
	widgetsState.drag_i = 50
	widgetsState.slider_f = 0.5
	widgetsState.slider_i = 50

	// Range Widgets
	widgetsState.range_begin = 10
	widgetsState.range_end = 90
	widgetsState.range_begin_i = 100
	widgetsState.range_end_i = 1000
}

func ShowDemoWindowWidgets() {
	if !CollapsingHeader("Widgets", 0) {
		return
	}

	if widgetsState.disable_all {
		BeginDisabled(true)
	}

	if TreeNode("Basic") {
		if Button("Button") {
			widgetsState.clicked++
		}
		if widgetsState.clicked&1 != 0 {
			SameLine(0, 0)
			Text("Thanks for clicking me!")
		}

		Checkbox("checkbox", &widgetsState.check)

		RadioButtonInt("radio a", &widgetsState.e, 0)
		SameLine(0, 0)
		RadioButtonInt("radio b", &widgetsState.e, 1)
		SameLine(0, 0)
		RadioButtonInt("radio c", &widgetsState.e, 2)

		// Color buttons, demonstrate using PushID() to add unique identifier in the ID stack, and changing style.
		for i := int(0); i < 7; i++ {
			if i > 0 {
				SameLine(0, 0)
			}
			PushID(i)
			col := ImVec4FromHSV(float(i)/7.0, 0.6, 0.6)
			PushStyleColorVec(ImGuiCol_Button, &col)
			col2 := ImVec4FromHSV(float(i)/7.0, 0.7, 0.7)
			PushStyleColorVec(ImGuiCol_ButtonHovered, &col2)
			col3 := ImVec4FromHSV(float(i)/7.0, 0.8, 0.8)
			PushStyleColorVec(ImGuiCol_ButtonActive, &col3)
			Button("Click")
			PopStyleColor(3)
			PopID()
		}

		// Use AlignTextToFramePadding() to align text baseline to the baseline of framed widgets elements
		AlignTextToFramePadding()
		Text("Hold to repeat:")
		SameLine(0, 0)

		// Arrow buttons with Repeater
		spacing := GetStyle().ItemInnerSpacing.x
		PushButtonRepeat(true)
		if ArrowButton("##left", ImGuiDir_Left) {
			widgetsState.counter--
		}
		SameLine(0.0, spacing)
		if ArrowButton("##right", ImGuiDir_Right) {
			widgetsState.counter++
		}
		PopButtonRepeat()
		SameLine(0, 0)
		Text("%d", widgetsState.counter)

		Text("Hover over me")
		if IsItemHovered(0) {
			SetTooltip("I am a tooltip")
		}

		SameLine(0, 0)
		Text("- or me")
		if IsItemHovered(0) {
			BeginTooltip()
			Text("I am a fancy tooltip")
			arr := []float{0.6, 0.1, 1.0, 0.5, 0.92, 0.1, 0.2}
			PlotLines("Curve", arr, int(len(arr)), 0, "", FLT_MAX, FLT_MAX, ImVec2{}, 4)
			EndTooltip()
		}

		Separator()

		LabelText("label", "Value")

		{
			// Using the _simplified_ one-liner Combo() api here
			items := []string{"AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIIIIII", "JJJJ", "KKKKKKK"}
			Combo("combo", &widgetsState.item_current_combo, items, int(len(items)), -1)
			SameLine(0, 0)
			HelpMarker("Using the simplified one-liner Combo API here.\nRefer to the \"Combo\" section below for an explanation of how to use the more flexible and general BeginCombo/EndCombo API.")
		}

		{
			InputText("input text", &widgetsState.str0, 0, nil, nil)
			SameLine(0, 0)
			HelpMarker(
				"USER:\n" +
					"Hold SHIFT or use mouse to select text.\n" +
					"CTRL+Left/Right to word jump.\n" +
					"CTRL+A or double-click to select all.\n" +
					"CTRL+X,CTRL+C,CTRL+V clipboard.\n" +
					"CTRL+Z,CTRL+Y undo/redo.\n" +
					"ESCAPE to revert.\n\n" +
					"PROGRAMMER:\n" +
					"You can use the ImGuiInputTextFlags_CallbackResize facility if you need to wire InputText() " +
					"to a dynamic string type. See misc/cpp/imgui_stdlib.h for an example (this is not demonstrated " +
					"in imgui_demo.cpp).")

			InputTextWithHint("input text (w/ hint)", "enter text here", &widgetsState.str1, 0, nil, nil)

			InputInt("input int", &widgetsState.i0, 1, 100, 0)
			SameLine(0, 0)
			HelpMarker(
				"You can apply arithmetic operators +,*,/ on numerical values.\n" +
					"  e.g. [ 100 ], input '*2', result becomes [ 200 ]\n" +
					"Use +- to subtract.")

			InputFloat("input float", &widgetsState.f0, 0.01, 1.0, "%.3f", 0)

			InputFloat("input scientific", &widgetsState.f1, 0.0, 0.0, "%e", 0)
			SameLine(0, 0)
			HelpMarker("You can input value using the scientific notation,\n  e.g. \"1e+8\" becomes \"100000000\".")

			var vec3 = [3]float{widgetsState.vec4a[0], widgetsState.vec4a[1], widgetsState.vec4a[2]}
			InputFloat3("input float3", &vec3, "%.3f", 0)
			widgetsState.vec4a[0], widgetsState.vec4a[1], widgetsState.vec4a[2] = vec3[0], vec3[1], vec3[2]
		}

		{
			DragInt("drag int", &widgetsState.i1, 1, 0, 0, "%d", 0)
			SameLine(0, 0)
			HelpMarker(
				"Click and drag to edit value.\n" +
					"Hold SHIFT/ALT for faster/slower edit.\n" +
					"Double-click or CTRL+click to input value.")

			DragInt("drag int 0..100", &widgetsState.i2, 1, 0, 100, "%d%%", ImGuiSliderFlags_AlwaysClamp)

			DragFloat("drag float", &widgetsState.f1drag, 0.005, 0.0, 0.0, "%.3f", 0)
			DragFloat("drag small float", &widgetsState.f2drag, 0.0001, 0.0, 0.0, "%.06f ns", 0)
		}

		{
			SliderInt("slider int", &widgetsState.i1slider, -1, 3, "%d", 0)
			SameLine(0, 0)
			HelpMarker("CTRL+click to input value.")

			SliderFloat("slider float", &widgetsState.f1slider, 0.0, 1.0, "ratio = %.3f", 0)
			SliderFloat("slider float (log)", &widgetsState.f2slider, -10.0, 10.0, "%.4f", ImGuiSliderFlags_Logarithmic)

			SliderAngle("slider angle", &widgetsState.angle, -360.0, 360.0, "%.0f deg", 0)

			// Using the format string to display a name instead of an integer.
			elems_names := []string{"Fire", "Earth", "Air", "Water"}
			elem_name := "Unknown"
			if widgetsState.elem >= 0 && widgetsState.elem < int(len(elems_names)) {
				elem_name = elems_names[widgetsState.elem]
			}
			SliderInt("slider enum", &widgetsState.elem, 0, int(len(elems_names)-1), elem_name, 0)
			SameLine(0, 0)
			HelpMarker("Using the format string parameter to display a name instead of the underlying integer.")
		}

		{
			ColorEdit3("color 1", &widgetsState.col1, 0)
			SameLine(0, 0)
			HelpMarker(
				"Click on the color square to open a color picker.\n" +
					"Click and hold to use drag and drop.\n" +
					"Right-click on the color square to show options.\n" +
					"CTRL+click on individual component to input value.\n")

			ColorEdit4("color 2", &widgetsState.col2, 0)
		}

		{
			// Using the _simplified_ one-liner ListBox() api here
			items := []string{"Apple", "Banana", "Cherry", "Kiwi", "Mango", "Orange", "Pineapple", "Strawberry", "Watermelon"}
			ListBox("listbox", &widgetsState.item_current_list, items, int(len(items)), 4)
			SameLine(0, 0)
			HelpMarker("Using the simplified one-liner ListBox API here.\nRefer to the \"List boxes\" section below for an explanation of how to use the more flexible and general BeginListBox/EndListBox API.")
		}

		TreePop()
	}

	if TreeNode("Trees") {
		if TreeNode("Basic trees") {
			for i := int(0); i < 5; i++ {
				// Use SetNextItemOpen() so set the default state of a node to be open.
				if i == 0 {
					SetNextItemOpen(true, ImGuiCond_Once)
				}

				if TreeNodeInterface(i, "Child %d", i) {
					Text("blah blah")
					SameLine(0, 0)
					if SmallButton("button") {
					}
					TreePop()
				}
			}
			TreePop()
		}

		if TreeNode("Advanced, with Selectable nodes") {
			HelpMarker(
				"This is a more typical looking tree with selectable nodes.\n" +
					"Click to select, CTRL+Click to toggle, click on arrows or double-click to open.")
			CheckboxFlagsInt("ImGuiTreeNodeFlags_OpenOnArrow", (*int32)(&widgetsState.base_flags), int32(ImGuiTreeNodeFlags_OpenOnArrow))
			CheckboxFlagsInt("ImGuiTreeNodeFlags_OpenOnDoubleClick", (*int32)(&widgetsState.base_flags), int32(ImGuiTreeNodeFlags_OpenOnDoubleClick))
			CheckboxFlagsInt("ImGuiTreeNodeFlags_SpanAvailWidth", (*int32)(&widgetsState.base_flags), int32(ImGuiTreeNodeFlags_SpanAvailWidth))
			SameLine(0, 0)
			HelpMarker("Extend hit area to all available width instead of allowing more items to be laid out after the node.")
			CheckboxFlagsInt("ImGuiTreeNodeFlags_SpanFullWidth", (*int32)(&widgetsState.base_flags), int32(ImGuiTreeNodeFlags_SpanFullWidth))
			Checkbox("Align label with current X position", &widgetsState.align_label_with_current_x_position)
			Checkbox("Test tree node as drag source", &widgetsState.test_drag_and_drop)
			Text("Hello!")
			if widgetsState.align_label_with_current_x_position {
				Unindent(GetTreeNodeToLabelSpacing())
			}

			node_clicked := int(-1)
			for i := int(0); i < 6; i++ {
				node_flags := widgetsState.base_flags
				is_selected := (widgetsState.selection_mask & (1 << i)) != 0
				if is_selected {
					node_flags |= ImGuiTreeNodeFlags_Selected
				}
				if i < 3 {
					// Items 0..2 are Tree Node
					node_open := TreeNodeInterfaceEx(i, node_flags, "Selectable Node %d", i)
					if IsItemClicked(0) {
						node_clicked = int(i)
					}
					if widgetsState.test_drag_and_drop && BeginDragDropSource(0) {
						SetDragDropPayload("_TREENODE", nil, 0, 0)
						Text("This is a drag and drop source")
						EndDragDropSource()
					}
					if node_open {
						BulletText("Blah blah\nBlah Blah")
						TreePop()
					}
				} else {
					// Items 3..5 are Tree Leaves
					node_flags |= ImGuiTreeNodeFlags_Leaf | ImGuiTreeNodeFlags_NoTreePushOnOpen
					TreeNodeInterfaceEx(i, node_flags, "Selectable Leaf %d", i)
					if IsItemClicked(0) {
						node_clicked = int(i)
					}
					if widgetsState.test_drag_and_drop && BeginDragDropSource(0) {
						SetDragDropPayload("_TREENODE", nil, 0, 0)
						Text("This is a drag and drop source")
						EndDragDropSource()
					}
				}
			}
			if node_clicked != -1 {
				// Update selection state
				if GetIO().KeyCtrl {
					widgetsState.selection_mask ^= 1 << node_clicked // CTRL+click to toggle
				} else {
					widgetsState.selection_mask = 1 << node_clicked // Click to single-select
				}
			}
			if widgetsState.align_label_with_current_x_position {
				Indent(GetTreeNodeToLabelSpacing())
			}
			TreePop()
		}
		TreePop()
	}

	if TreeNode("Collapsing Headers") {
		Checkbox("Show 2nd header", &widgetsState.closable_group)
		if CollapsingHeader("Header", 0) {
			Text("IsItemHovered: %d", bool2int(IsItemHovered(0)))
			for i := int(0); i < 5; i++ {
				Text("Some content %d", i)
			}
		}
		if CollapsingHeaderVisible("Header with a close button", &widgetsState.closable_group, 0) {
			Text("IsItemHovered: %d", bool2int(IsItemHovered(0)))
			for i := int(0); i < 5; i++ {
				Text("More content %d", i)
			}
		}
		TreePop()
	}

	if TreeNode("Bullets") {
		BulletText("Bullet point 1")
		BulletText("Bullet point 2\nOn multiple lines")
		if TreeNode("Tree node") {
			BulletText("Another bullet point")
			TreePop()
		}
		Bullet()
		Text("Bullet point 3 (two calls)")
		Bullet()
		SmallButton("Button")
		TreePop()
	}

	if TreeNode("Text") {
		if TreeNode("Colorful Text") {
			// Using shortcut. You can use PushStyleColor()/PopStyleColor() for more flexibility.
			pink := ImVec4{1.0, 0.0, 1.0, 1.0}
			TextColored(&pink, "Pink")
			yellow := ImVec4{1.0, 1.0, 0.0, 1.0}
			TextColored(&yellow, "Yellow")
			TextDisabled("Disabled")
			SameLine(0, 0)
			HelpMarker("The TextDisabled color is stored in ImGuiStyle.")
			TreePop()
		}

		if TreeNode("Word Wrapping") {
			// Using shortcut. You can use PushTextWrapPos()/PopTextWrapPos() for more flexibility.
			TextWrapped(
				"This text should automatically wrap on the edge of the window. The current implementation " +
					"for text wrapping follows simple rules suitable for English and possibly other languages.")
			Spacing()

			SliderFloat("Wrap width", &widgetsState.wrap_width, -20, 600, "%.0f", 0)

			draw_list := GetWindowDrawList()
			for n := int(0); n < 2; n++ {
				Text("Test paragraph %d:", n)
				pos := GetCursorScreenPos()
				marker_min := ImVec2{pos.x + widgetsState.wrap_width, pos.y}
				marker_max := ImVec2{pos.x + widgetsState.wrap_width + 10, pos.y + GetTextLineHeight()}
				PushTextWrapPos(GetCursorPos().x + widgetsState.wrap_width)
				if n == 0 {
					Text("The lazy dog is a good dog. This paragraph should fit within %.0f pixels. Testing a 1 character word. The quick brown fox jumps over the lazy dog.", widgetsState.wrap_width)
				} else {
					Text("aaaaaaaa bbbbbbbb, c cccccccc,dddddddd. d eeeeeeee   ffffffff. gggggggg!hhhhhhhh")
				}

				// Draw actual text bounding box, following by marker of our expected limit
				draw_list.AddRect(GetItemRectMin(), GetItemRectMax(), IM_COL32(255, 255, 0, 255), 0, 0, 1.0)
				draw_list.AddRectFilled(marker_min, marker_max, IM_COL32(255, 0, 255, 255), 0, 0)
				PopTextWrapPos()
			}

			TreePop()
		}

		if TreeNode("UTF-8 Text") {
			TextWrapped(
				"CJK text will only appears if the font was loaded with the appropriate CJK character ranges. " +
					"Call io.Fonts->AddFontFromFileTTF() manually to load extra character ranges. " +
					"Read docs/FONTS.md for details.")
			Text("Hiragana: かきくけこ (kakikukeko)")
			Text("Kanjis: 日本語 (nihongo)")
			TreePop()
		}
		TreePop()
	}

	if TreeNode("Images") {
		io := GetIO()
		TextWrapped(
			"Below we are displaying the font texture (which is the only texture we have access to in this demo). " +
				"Use the 'ImTextureID' type as storage to pass pointers or identifier to your own texture data. " +
				"Hover the texture for a zoomed view!")

		// Display the font texture
		my_tex_id := io.Fonts.TexID
		my_tex_w := float(io.Fonts.TexWidth)
		my_tex_h := float(io.Fonts.TexHeight)
		{
			Text("%.0fx%.0f", my_tex_w, my_tex_h)
			pos := GetCursorScreenPos()
			uv_min := ImVec2{0.0, 0.0}
			uv_max := ImVec2{1.0, 1.0}
			tint_col := ImVec4{1.0, 1.0, 1.0, 1.0}
			border_col := ImVec4{1.0, 1.0, 1.0, 0.5}
			Image(my_tex_id, ImVec2{my_tex_w, my_tex_h}, uv_min, uv_max, tint_col, border_col)
			if IsItemHovered(0) {
				BeginTooltip()
				region_sz := float(32.0)
				region_x := io.MousePos.x - pos.x - region_sz*0.5
				region_y := io.MousePos.y - pos.y - region_sz*0.5
				zoom := float(4.0)
				if region_x < 0.0 {
					region_x = 0.0
				} else if region_x > my_tex_w-region_sz {
					region_x = my_tex_w - region_sz
				}
				if region_y < 0.0 {
					region_y = 0.0
				} else if region_y > my_tex_h-region_sz {
					region_y = my_tex_h - region_sz
				}
				Text("Min: (%.2f, %.2f)", region_x, region_y)
				Text("Max: (%.2f, %.2f)", region_x+region_sz, region_y+region_sz)
				uv0 := ImVec2{region_x / my_tex_w, region_y / my_tex_h}
				uv1 := ImVec2{(region_x + region_sz) / my_tex_w, (region_y + region_sz) / my_tex_h}
				Image(my_tex_id, ImVec2{region_sz * zoom, region_sz * zoom}, uv0, uv1, tint_col, border_col)
				EndTooltip()
			}
		}
		TextWrapped("And now some textured buttons..")
		for i := int(0); i < 8; i++ {
			PushID(i)
			frame_padding := int(-1 + i)
			size := ImVec2{32.0, 32.0}
			uv0 := ImVec2{0.0, 0.0}
			uv1 := ImVec2{32.0 / my_tex_w, 32.0 / my_tex_h}
			bg_col := ImVec4{0.0, 0.0, 0.0, 1.0}
			tint_col := ImVec4{1.0, 1.0, 1.0, 1.0}
			if ImageButton(my_tex_id, size, uv0, uv1, frame_padding, bg_col, tint_col) {
				widgetsState.pressed_count++
			}
			PopID()
			SameLine(0, 0)
		}
		NewLine()
		Text("Pressed %d times.", widgetsState.pressed_count)
		TreePop()
	}

	if TreeNode("Combo") {
		// Expose flags as checkbox for the demo
		combo_flags_int := int(widgetsState.combo_flags)
		CheckboxFlagsInt("ImGuiComboFlags_PopupAlignLeft", &combo_flags_int, int(ImGuiComboFlags_PopupAlignLeft))
		widgetsState.combo_flags = ImGuiComboFlags(combo_flags_int)
		SameLine(0, 0)
		HelpMarker("Only makes a difference if the popup is larger than the combo")
		combo_flags_int = int(widgetsState.combo_flags)
		if CheckboxFlagsInt("ImGuiComboFlags_NoArrowButton", &combo_flags_int, int(ImGuiComboFlags_NoArrowButton)) {
			combo_flags_int &= ^int(ImGuiComboFlags_NoPreview)
		}
		widgetsState.combo_flags = ImGuiComboFlags(combo_flags_int)
		combo_flags_int = int(widgetsState.combo_flags)
		if CheckboxFlagsInt("ImGuiComboFlags_NoPreview", &combo_flags_int, int(ImGuiComboFlags_NoPreview)) {
			combo_flags_int &= ^int(ImGuiComboFlags_NoArrowButton)
		}
		widgetsState.combo_flags = ImGuiComboFlags(combo_flags_int)

		// Using the generic BeginCombo() API
		items := []string{"AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO"}
		combo_preview_value := items[widgetsState.combo_item_current_idx]
		if BeginCombo("combo 1", combo_preview_value, widgetsState.combo_flags) {
			for n := int(0); n < int(len(items)); n++ {
				is_selected := widgetsState.combo_item_current_idx == n
				if Selectable(items[n], is_selected, 0, ImVec2{}) {
					widgetsState.combo_item_current_idx = n
				}
				if is_selected {
					SetItemDefaultFocus()
				}
			}
			EndCombo()
		}
		TreePop()
	}

	if TreeNode("List boxes") {
		items := []string{"AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO"}
		if BeginListBox("listbox 1", ImVec2{}) {
			for n := int(0); n < int(len(items)); n++ {
				is_selected := widgetsState.listbox_item_current_idx == n
				if Selectable(items[n], is_selected, 0, ImVec2{}) {
					widgetsState.listbox_item_current_idx = n
				}
				if is_selected {
					SetItemDefaultFocus()
				}
			}
			EndListBox()
		}

		// Custom size: use all width, 5 items tall
		Text("Full-width:")
		if BeginListBox("##listbox 2", ImVec2{-FLT_MIN, 5 * GetTextLineHeightWithSpacing()}) {
			for n := int(0); n < int(len(items)); n++ {
				is_selected := widgetsState.listbox_item_current_idx == n
				if Selectable(items[n], is_selected, 0, ImVec2{}) {
					widgetsState.listbox_item_current_idx = n
				}
				if is_selected {
					SetItemDefaultFocus()
				}
			}
			EndListBox()
		}
		TreePop()
	}

	if TreeNode("Selectables") {
		if TreeNode("Basic") {
			SelectablePointer("1. I am selectable", &widgetsState.selectable_basic[0], 0, ImVec2{})
			SelectablePointer("2. I am selectable", &widgetsState.selectable_basic[1], 0, ImVec2{})
			Text("(I am not selectable)")
			SelectablePointer("4. I am selectable", &widgetsState.selectable_basic[3], 0, ImVec2{})
			if Selectable("5. I am double clickable", widgetsState.selectable_basic[4], ImGuiSelectableFlags_AllowDoubleClick, ImVec2{}) {
				if IsMouseDoubleClicked(0) {
					widgetsState.selectable_basic[4] = !widgetsState.selectable_basic[4]
				}
			}
			TreePop()
		}
		if TreeNode("Selection State: Single Selection") {
			for n := int(0); n < 5; n++ {
				label := fmt.Sprintf("Object %d", n)
				if Selectable(label, widgetsState.selectable_single == n, 0, ImVec2{}) {
					widgetsState.selectable_single = n
				}
			}
			TreePop()
		}
		if TreeNode("Selection State: Multiple Selection") {
			HelpMarker("Hold CTRL and click to select multiple items.")
			for n := int(0); n < 5; n++ {
				label := fmt.Sprintf("Object %d", n)
				if Selectable(label, widgetsState.selectable_multiple[n], 0, ImVec2{}) {
					if !GetIO().KeyCtrl {
						// Clear selection when CTRL is not held
						for i := range widgetsState.selectable_multiple {
							widgetsState.selectable_multiple[i] = false
						}
					}
					widgetsState.selectable_multiple[n] = !widgetsState.selectable_multiple[n]
				}
			}
			TreePop()
		}
		if TreeNode("Rendering more text into the same line") {
			SelectablePointer("main.c", &widgetsState.selectable_render[0], 0, ImVec2{})
			SameLine(300, 0)
			Text(" 2,345 bytes")
			SelectablePointer("Hello.cpp", &widgetsState.selectable_render[1], 0, ImVec2{})
			SameLine(300, 0)
			Text("12,345 bytes")
			SelectablePointer("Hello.h", &widgetsState.selectable_render[2], 0, ImVec2{})
			SameLine(300, 0)
			Text(" 2,345 bytes")
			TreePop()
		}
		if TreeNode("In columns") {
			if BeginTable("split1", 3, ImGuiTableFlags_Resizable|ImGuiTableFlags_NoSavedSettings|ImGuiTableFlags_Borders, ImVec2{}, 0) {
				for i := int(0); i < 10; i++ {
					label := fmt.Sprintf("Item %d", i)
					TableNextColumn()
					SelectablePointer(label, &widgetsState.selectable_columns[i], 0, ImVec2{})
				}
				EndTable()
			}
			Separator()
			if BeginTable("split2", 3, ImGuiTableFlags_Resizable|ImGuiTableFlags_NoSavedSettings|ImGuiTableFlags_Borders, ImVec2{}, 0) {
				for i := int(0); i < 10; i++ {
					label := fmt.Sprintf("Item %d", i)
					TableNextRow(0, 0)
					TableNextColumn()
					SelectablePointer(label, &widgetsState.selectable_columns[i], ImGuiSelectableFlags_SpanAllColumns, ImVec2{})
					TableNextColumn()
					Text("Some other contents")
					TableNextColumn()
					Text("123456")
				}
				EndTable()
			}
			TreePop()
		}
		if TreeNode("Grid") {
			// Add in a bit of silly fun...
			time := float(GetTime())
			winning_state := true
			for ri := 0; ri < 4; ri++ {
				for ci := 0; ci < 4; ci++ {
					if widgetsState.selectable_grid[ri][ci] {
						winning_state = false
					}
				}
			}
			if winning_state {
				PushStyleVec(ImGuiStyleVar_SelectableTextAlign, ImVec2{0.5 + 0.5*ImCos(time*2.0), 0.5 + 0.5*ImSin(time*3.0)})
			}

			for ri := 0; ri < 4; ri++ {
				for ci := 0; ci < 4; ci++ {
					if ci > 0 {
						SameLine(0, 0)
					}
					PushID(int(ri*4 + ci))
					if Selectable("Sailor", widgetsState.selectable_grid[ri][ci], 0, ImVec2{50, 50}) {
						// Toggle clicked cell + toggle neighbors
						widgetsState.selectable_grid[ri][ci] = !widgetsState.selectable_grid[ri][ci]
						if ci > 0 {
							widgetsState.selectable_grid[ri][ci-1] = !widgetsState.selectable_grid[ri][ci-1]
						}
						if ci < 3 {
							widgetsState.selectable_grid[ri][ci+1] = !widgetsState.selectable_grid[ri][ci+1]
						}
						if ri > 0 {
							widgetsState.selectable_grid[ri-1][ci] = !widgetsState.selectable_grid[ri-1][ci]
						}
						if ri < 3 {
							widgetsState.selectable_grid[ri+1][ci] = !widgetsState.selectable_grid[ri+1][ci]
						}
					}
					PopID()
				}
			}

			if winning_state {
				PopStyleVar(1)
			}
			TreePop()
		}
		TreePop()
	}

	// Tabs
	if TreeNode("Tabs") {
		if TreeNode("Basic") {
			tab_bar_flags := ImGuiTabBarFlags_None
			if BeginTabBar("MyTabBar", tab_bar_flags) {
				if BeginTabItem("Avocado", nil, 0) {
					Text("This is the Avocado tab!\nblah blah blah blah blah")
					EndTabItem()
				}
				if BeginTabItem("Broccoli", nil, 0) {
					Text("This is the Broccoli tab!\nblah blah blah blah blah")
					EndTabItem()
				}
				if BeginTabItem("Cucumber", nil, 0) {
					Text("This is the Cucumber tab!\nblah blah blah blah blah")
					EndTabItem()
				}
				EndTabBar()
			}
			Separator()
			TreePop()
		}

		if TreeNode("Advanced & Close Button") {
			// Expose a couple of the available flags
			tab_bar_flags_int := int(widgetsState.tab_bar_flags)
			CheckboxFlagsInt("ImGuiTabBarFlags_Reorderable", &tab_bar_flags_int, int(ImGuiTabBarFlags_Reorderable))
			CheckboxFlagsInt("ImGuiTabBarFlags_AutoSelectNewTabs", &tab_bar_flags_int, int(ImGuiTabBarFlags_AutoSelectNewTabs))
			CheckboxFlagsInt("ImGuiTabBarFlags_TabListPopupButton", &tab_bar_flags_int, int(ImGuiTabBarFlags_TabListPopupButton))
			CheckboxFlagsInt("ImGuiTabBarFlags_NoCloseWithMiddleMouseButton", &tab_bar_flags_int, int(ImGuiTabBarFlags_NoCloseWithMiddleMouseButton))
			widgetsState.tab_bar_flags = ImGuiTabBarFlags(tab_bar_flags_int)

			// Tab Bar
			names := []string{"Artichoke", "Beetroot", "Celery", "Daikon"}
			for n := int(0); n < int(len(names)); n++ {
				if n > 0 {
					SameLine(0, 0)
				}
				Checkbox(names[n], &widgetsState.tabs_opened[n])
			}

			// Passing a bool* to BeginTabItem() is similar to passing one to Begin()
			if BeginTabBar("MyTabBar2", widgetsState.tab_bar_flags) {
				for n := int(0); n < int(len(names)); n++ {
					if widgetsState.tabs_opened[n] && BeginTabItem(names[n], &widgetsState.tabs_opened[n], ImGuiTabItemFlags_None) {
						Text("This is the %s tab!", names[n])
						if n&1 != 0 {
							Text("I am an odd tab.")
						}
						EndTabItem()
					}
				}
				EndTabBar()
			}
			Separator()
			TreePop()
		}
		TreePop()
	}

	// Plots Widgets
	if TreeNode("Plots Widgets") {
		Checkbox("Animate", &widgetsState.plots_animate)

		// Plot as lines and plot as histogram
		arr := []float{0.6, 0.1, 1.0, 0.5, 0.92, 0.1, 0.2}
		PlotLines("Frame Times", arr, int(len(arr)), 0, "", FLT_MAX, FLT_MAX, ImVec2{}, 4)
		PlotHistogram("Histogram", arr, int(len(arr)), 0, "", 0.0, 1.0, ImVec2{0, 80.0}, 4)

		// Fill an array of contiguous float values to plot
		if !widgetsState.plots_animate || widgetsState.plots_refresh == 0.0 {
			widgetsState.plots_refresh = GetTime()
		}
		for widgetsState.plots_refresh < GetTime() {
			widgetsState.plots_values[widgetsState.plots_offset] = ImCos(widgetsState.plots_phase)
			widgetsState.plots_offset = (widgetsState.plots_offset + 1) % int(len(widgetsState.plots_values))
			widgetsState.plots_phase += 0.10 * float(widgetsState.plots_offset)
			widgetsState.plots_refresh += 1.0 / 60.0
		}

		// Plots can display overlay texts
		{
			average := float(0.0)
			for n := int(0); n < int(len(widgetsState.plots_values)); n++ {
				average += widgetsState.plots_values[n]
			}
			average /= float(len(widgetsState.plots_values))
			overlay := fmt.Sprintf("avg %f", average)
			PlotLines("Lines", widgetsState.plots_values[:], int(len(widgetsState.plots_values)), widgetsState.plots_offset, overlay, -1.0, 1.0, ImVec2{0, 80.0}, 4)
		}

		Separator()

		// Animate a simple progress bar
		if widgetsState.plots_animate {
			widgetsState.plots_progress += widgetsState.plots_progress_dir * 0.4 * GetIO().DeltaTime
			if widgetsState.plots_progress >= 1.1 {
				widgetsState.plots_progress = 1.1
				widgetsState.plots_progress_dir *= -1.0
			}
			if widgetsState.plots_progress <= -0.1 {
				widgetsState.plots_progress = -0.1
				widgetsState.plots_progress_dir *= -1.0
			}
		}

		// Typically we would use ImVec2(-1.0f,0.0f) or ImVec2(-FLT_MIN,0.0f) to use all available width
		ProgressBar(widgetsState.plots_progress, ImVec2{0.0, 0.0}, "")
		SameLine(0, GetStyle().ItemInnerSpacing.x)
		Text("Progress Bar")

		progress_saturated := widgetsState.plots_progress
		if progress_saturated < 0.0 {
			progress_saturated = 0.0
		}
		if progress_saturated > 1.0 {
			progress_saturated = 1.0
		}
		buf := fmt.Sprintf("%d/%d", int(progress_saturated*1753), 1753)
		ProgressBar(widgetsState.plots_progress, ImVec2{0.0, 0.0}, buf)
		TreePop()
	}

	// Color/Picker Widgets
	if TreeNode("Color/Picker Widgets") {
		Checkbox("With Alpha Preview", &widgetsState.color_alpha_preview)
		Checkbox("With Half Alpha Preview", &widgetsState.color_alpha_half_preview)
		Checkbox("With Drag and Drop", &widgetsState.color_drag_and_drop)
		Checkbox("With Options Menu", &widgetsState.color_options_menu)
		SameLine(0, 0)
		HelpMarker("Right-click on the individual color widget to show options.")
		Checkbox("With HDR", &widgetsState.color_hdr)
		SameLine(0, 0)
		HelpMarker("Currently all this does is to lift the 0..1 limits on dragging widgets.")

		misc_flags := ImGuiColorEditFlags(0)
		if widgetsState.color_hdr {
			misc_flags |= ImGuiColorEditFlags_HDR
		}
		if !widgetsState.color_drag_and_drop {
			misc_flags |= ImGuiColorEditFlags_NoDragDrop
		}
		if widgetsState.color_alpha_half_preview {
			misc_flags |= ImGuiColorEditFlags_AlphaPreviewHalf
		} else if widgetsState.color_alpha_preview {
			misc_flags |= ImGuiColorEditFlags_AlphaPreview
		}
		if !widgetsState.color_options_menu {
			misc_flags |= ImGuiColorEditFlags_NoOptions
		}

		Text("Color widget:")
		SameLine(0, 0)
		HelpMarker("Click on the color square to open a color picker.\nCTRL+click on individual component to input value.")
		col3 := [3]float{widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z}
		ColorEdit3("MyColor##1", &col3, misc_flags)
		widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z = col3[0], col3[1], col3[2]

		Text("Color widget HSV with Alpha:")
		col4 := [4]float{widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w}
		ColorEdit4("MyColor##2", &col4, ImGuiColorEditFlags_DisplayHSV|misc_flags)
		widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w = col4[0], col4[1], col4[2], col4[3]

		Text("Color widget with Float Display:")
		col4 = [4]float{widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w}
		ColorEdit4("MyColor##2f", &col4, ImGuiColorEditFlags_Float|misc_flags)
		widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w = col4[0], col4[1], col4[2], col4[3]

		Text("Color button with Picker:")
		SameLine(0, 0)
		HelpMarker("With the ImGuiColorEditFlags_NoInputs flag you can hide all the slider/text inputs.\nWith the ImGuiColorEditFlags_NoLabel flag you can pass a non-empty label which will only be used for the tooltip and picker popup.")
		col4 = [4]float{widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w}
		ColorEdit4("MyColor##3", &col4, ImGuiColorEditFlags_NoInputs|ImGuiColorEditFlags_NoLabel|misc_flags)
		widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w = col4[0], col4[1], col4[2], col4[3]

		Text("Color button only:")
		Checkbox("ImGuiColorEditFlags_NoBorder", &widgetsState.color_no_border)
		no_border_flag := ImGuiColorEditFlags(0)
		if widgetsState.color_no_border {
			no_border_flag = ImGuiColorEditFlags_NoBorder
		}
		ColorButton("MyColor##3c", widgetsState.color_picker_color, misc_flags|no_border_flag, ImVec2{80, 80})

		Text("Color picker:")
		Checkbox("With Alpha", &widgetsState.color_alpha)
		Checkbox("With Alpha Bar", &widgetsState.color_alpha_bar)
		Checkbox("With Side Preview", &widgetsState.color_side_preview)
		if widgetsState.color_side_preview {
			SameLine(0, 0)
			Checkbox("With Ref Color", &widgetsState.color_ref_color)
			if widgetsState.color_ref_color {
				SameLine(0, 0)
				ref_col := [4]float{widgetsState.color_ref_color_v.x, widgetsState.color_ref_color_v.y, widgetsState.color_ref_color_v.z, widgetsState.color_ref_color_v.w}
				ColorEdit4("##RefColor", &ref_col, ImGuiColorEditFlags_NoInputs|misc_flags)
				widgetsState.color_ref_color_v.x, widgetsState.color_ref_color_v.y, widgetsState.color_ref_color_v.z, widgetsState.color_ref_color_v.w = ref_col[0], ref_col[1], ref_col[2], ref_col[3]
			}
		}
		Combo("Display Mode", &widgetsState.color_display_mode, []string{"Auto/Current", "None", "RGB Only", "HSV Only", "Hex Only"}, 5, -1)
		SameLine(0, 0)
		HelpMarker("ColorEdit defaults to displaying RGB inputs if you don't specify a display mode, but the user can change it with a right-click.\n\nColorPicker defaults to displaying RGB+HSV+Hex if you don't specify a display mode.\n\nYou can change the defaults using SetColorEditOptions().")
		Combo("Picker Mode", &widgetsState.color_picker_mode, []string{"Auto/Current", "Hue bar + SV rect", "Hue wheel + SV triangle"}, 3, -1)
		SameLine(0, 0)
		HelpMarker("User can right-click the picker to change mode.")

		picker_flags := misc_flags
		if !widgetsState.color_alpha {
			picker_flags |= ImGuiColorEditFlags_NoAlpha
		}
		if widgetsState.color_alpha_bar {
			picker_flags |= ImGuiColorEditFlags_AlphaBar
		}
		if !widgetsState.color_side_preview {
			picker_flags |= ImGuiColorEditFlags_NoSidePreview
		}
		if widgetsState.color_picker_mode == 1 {
			picker_flags |= ImGuiColorEditFlags_PickerHueBar
		}
		if widgetsState.color_picker_mode == 2 {
			picker_flags |= ImGuiColorEditFlags_PickerHueWheel
		}
		if widgetsState.color_display_mode == 1 {
			picker_flags |= ImGuiColorEditFlags_NoInputs
		}
		if widgetsState.color_display_mode == 2 {
			picker_flags |= ImGuiColorEditFlags_DisplayRGB
		}
		if widgetsState.color_display_mode == 3 {
			picker_flags |= ImGuiColorEditFlags_DisplayHSV
		}
		if widgetsState.color_display_mode == 4 {
			picker_flags |= ImGuiColorEditFlags_DisplayHex
		}

		col4 = [4]float{widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w}
		var ref_col []float
		if widgetsState.color_ref_color {
			ref_col = []float{widgetsState.color_ref_color_v.x, widgetsState.color_ref_color_v.y, widgetsState.color_ref_color_v.z, widgetsState.color_ref_color_v.w}
		}
		ColorPicker4("MyColor##4", &col4, picker_flags, ref_col)
		widgetsState.color_picker_color.x, widgetsState.color_picker_color.y, widgetsState.color_picker_color.z, widgetsState.color_picker_color.w = col4[0], col4[1], col4[2], col4[3]

		Text("Set defaults in code:")
		SameLine(0, 0)
		HelpMarker("SetColorEditOptions() is designed to allow you to set boot-time default.\nWe don't have Push/Pop functions because you can force options on a per-widget basis if needed, and the user can change non-forced ones with the options menu.\nWe don't have a getter to avoid encouraging you to persistently save values that aren't forward-compatible.")
		if Button("Default: Uint8 + HSV + Hue Bar") {
			SetColorEditOptions(ImGuiColorEditFlags_Uint8 | ImGuiColorEditFlags_DisplayHSV | ImGuiColorEditFlags_PickerHueBar)
		}
		if Button("Default: Float + HDR + Hue Wheel") {
			SetColorEditOptions(ImGuiColorEditFlags_Float | ImGuiColorEditFlags_HDR | ImGuiColorEditFlags_PickerHueWheel)
		}

		// HSV encoded support
		Spacing()
		Text("HSV encoded colors")
		SameLine(0, 0)
		HelpMarker("By default, colors are given to ColorEdit and ColorPicker in RGB, but ImGuiColorEditFlags_InputHSV allows you to store colors as HSV and pass them to ColorEdit and ColorPicker as HSV. This comes with the added benefit that you can manipulate hue values with the picker even when saturation or value are zero.")
		Text("Color widget with InputHSV:")
		hsv_col := [4]float{widgetsState.color_hsv.x, widgetsState.color_hsv.y, widgetsState.color_hsv.z, widgetsState.color_hsv.w}
		ColorEdit4("HSV shown as RGB##1", &hsv_col, ImGuiColorEditFlags_DisplayRGB|ImGuiColorEditFlags_InputHSV|ImGuiColorEditFlags_Float)
		ColorEdit4("HSV shown as HSV##1", &hsv_col, ImGuiColorEditFlags_DisplayHSV|ImGuiColorEditFlags_InputHSV|ImGuiColorEditFlags_Float)
		widgetsState.color_hsv.x, widgetsState.color_hsv.y, widgetsState.color_hsv.z, widgetsState.color_hsv.w = hsv_col[0], hsv_col[1], hsv_col[2], hsv_col[3]

		TreePop()
	}

	// Drag/Slider Flags
	if TreeNode("Drag/Slider Flags") {
		// Demonstrate using advanced flags for DragXXX and SliderXXX functions
		slider_flags_int := int(widgetsState.slider_flags)
		CheckboxFlagsInt("ImGuiSliderFlags_AlwaysClamp", &slider_flags_int, int(ImGuiSliderFlags_AlwaysClamp))
		SameLine(0, 0)
		HelpMarker("Always clamp value to min/max bounds (if any) when input manually with CTRL+Click.")
		CheckboxFlagsInt("ImGuiSliderFlags_Logarithmic", &slider_flags_int, int(ImGuiSliderFlags_Logarithmic))
		SameLine(0, 0)
		HelpMarker("Enable logarithmic editing (more precision for small values).")
		CheckboxFlagsInt("ImGuiSliderFlags_NoRoundToFormat", &slider_flags_int, int(ImGuiSliderFlags_NoRoundToFormat))
		SameLine(0, 0)
		HelpMarker("Disable rounding underlying value to match precision of the format string (e.g. %.3f values are rounded to those 3 digits).")
		CheckboxFlagsInt("ImGuiSliderFlags_NoInput", &slider_flags_int, int(ImGuiSliderFlags_NoInput))
		SameLine(0, 0)
		HelpMarker("Disable CTRL+Click or Enter key allowing to input text directly into the widget.")
		widgetsState.slider_flags = ImGuiSliderFlags(slider_flags_int)

		// Drags
		Text("Underlying float value: %f", widgetsState.drag_f)
		DragFloat("DragFloat (0 -> 1)", &widgetsState.drag_f, 0.005, 0.0, 1.0, "%.3f", widgetsState.slider_flags)
		DragFloat("DragFloat (0 -> +inf)", &widgetsState.drag_f, 0.005, 0.0, FLT_MAX, "%.3f", widgetsState.slider_flags)
		DragFloat("DragFloat (-inf -> 1)", &widgetsState.drag_f, 0.005, -FLT_MAX, 1.0, "%.3f", widgetsState.slider_flags)
		DragFloat("DragFloat (-inf -> +inf)", &widgetsState.drag_f, 0.005, -FLT_MAX, FLT_MAX, "%.3f", widgetsState.slider_flags)
		DragInt("DragInt (0 -> 100)", &widgetsState.drag_i, 0.5, 0, 100, "%d", widgetsState.slider_flags)

		// Sliders
		Text("Underlying float value: %f", widgetsState.slider_f)
		SliderFloat("SliderFloat (0 -> 1)", &widgetsState.slider_f, 0.0, 1.0, "%.3f", widgetsState.slider_flags)
		SliderInt("SliderInt (0 -> 100)", &widgetsState.slider_i, 0, 100, "%d", widgetsState.slider_flags)

		TreePop()
	}

	// Range Widgets
	if TreeNode("Range Widgets") {
		DragFloatRange2("range float", &widgetsState.range_begin, &widgetsState.range_end, 0.25, 0.0, 100.0, "Min: %.1f %%", "Max: %.1f %%", ImGuiSliderFlags_AlwaysClamp)
		DragIntRange2("range int", &widgetsState.range_begin_i, &widgetsState.range_end_i, 5, 0, 1000, "Min: %d units", "Max: %d units", 0)
		DragIntRange2("range int (no bounds)", &widgetsState.range_begin_i, &widgetsState.range_end_i, 5, 0, 0, "Min: %d units", "Max: %d units", 0)
		TreePop()
	}

	if widgetsState.disable_all {
		EndDisabled()
	}
}

// State for ShowDemoWindowLayout
var layoutState struct {
	// Child windows
	disable_mouse_wheel bool
	disable_menu        bool
	offset_x            int

	// Widgets Width
	width_f              float
	show_indented_items  bool

	// Basic Horizontal Layout
	c1, c2, c3, c4 bool
	f0, f1, f2     float
	item           int

	// Groups
	group_values [4]float
}

func init() {
	// Initialize layout state
	layoutState.show_indented_items = true
	layoutState.f0 = 1.0
	layoutState.f1 = 2.0
	layoutState.f2 = 3.0
	layoutState.item = -1
}

func ShowDemoWindowLayout() {
	if !CollapsingHeader("Layout & Scrolling", 0) {
		return
	}

	if TreeNode("Child windows") {
		HelpMarker("Use child windows to begin into a self-contained independent scrolling/clipping regions within a host window.")
		Checkbox("Disable Mouse Wheel", &layoutState.disable_mouse_wheel)
		Checkbox("Disable Menu", &layoutState.disable_menu)

		// Child 1: no border, enable horizontal scrollbar
		{
			window_flags := ImGuiWindowFlags_HorizontalScrollbar
			if layoutState.disable_mouse_wheel {
				window_flags |= ImGuiWindowFlags_NoScrollWithMouse
			}
			BeginChild("ChildL", ImVec2{GetContentRegionAvail().x * 0.5, 260}, false, window_flags)
			for i := int(0); i < 100; i++ {
				Text("%04d: scrollable region", i)
			}
			EndChild()
		}

		SameLine(0, 0)

		// Child 2: rounded border
		{
			window_flags := ImGuiWindowFlags_None
			if layoutState.disable_mouse_wheel {
				window_flags |= ImGuiWindowFlags_NoScrollWithMouse
			}
			if !layoutState.disable_menu {
				window_flags |= ImGuiWindowFlags_MenuBar
			}
			PushStyleFloat(ImGuiStyleVar_ChildRounding, 5.0)
			BeginChild("ChildR", ImVec2{0, 260}, true, window_flags)
			if !layoutState.disable_menu && BeginMenuBar() {
				if BeginMenu("Menu", true) {
					// Simplified menu
					MenuItem("New", "", nil, true)
					MenuItem("Open", "Ctrl+O", nil, true)
					EndMenu()
				}
				EndMenuBar()
			}
			if BeginTable("split", 2, ImGuiTableFlags_Resizable|ImGuiTableFlags_NoSavedSettings, ImVec2{}, 0) {
				for i := int(0); i < 100; i++ {
					buf := fmt.Sprintf("%03d", i)
					TableNextColumn()
					Button(buf)
				}
				EndTable()
			}
			EndChild()
			PopStyleVar(1)
		}

		Separator()

		// Demonstrate a few extra things
		{
			SetNextItemWidth(GetFontSize() * 8)
			DragInt("Offset X", &layoutState.offset_x, 1.0, -1000, 1000, "%d", 0)

			SetCursorPosX(GetCursorPosX() + float(layoutState.offset_x))
			PushStyleColorInt(ImGuiCol_ChildBg, IM_COL32(255, 0, 0, 100))
			BeginChild("Red", ImVec2{200, 100}, true, ImGuiWindowFlags_None)
			for n := int(0); n < 50; n++ {
				Text("Some test %d", n)
			}
			EndChild()
			child_is_hovered := IsItemHovered(0)
			child_rect_min := GetItemRectMin()
			child_rect_max := GetItemRectMax()
			PopStyleColor(1)
			Text("Hovered: %d", bool2int(child_is_hovered))
			Text("Rect of child window is: (%.0f,%.0f) (%.0f,%.0f)", child_rect_min.x, child_rect_min.y, child_rect_max.x, child_rect_max.y)
		}

		TreePop()
	}

	if TreeNode("Widgets Width") {
		Checkbox("Show indented items", &layoutState.show_indented_items)

		// Use SetNextItemWidth() to set the width of a single upcoming item.
		Text("SetNextItemWidth/PushItemWidth(100)")
		SameLine(0, 0)
		HelpMarker("Fixed width.")
		PushItemWidth(100)
		DragFloat("float##1b", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
		if layoutState.show_indented_items {
			Indent(0)
			DragFloat("float (indented)##1b", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
			Unindent(0)
		}
		PopItemWidth()

		Text("SetNextItemWidth/PushItemWidth(-100)")
		SameLine(0, 0)
		HelpMarker("Align to right edge minus 100")
		PushItemWidth(-100)
		DragFloat("float##2a", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
		if layoutState.show_indented_items {
			Indent(0)
			DragFloat("float (indented)##2b", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
			Unindent(0)
		}
		PopItemWidth()

		Text("SetNextItemWidth/PushItemWidth(GetContentRegionAvail().x * 0.5f)")
		SameLine(0, 0)
		HelpMarker("Half of available width.\n(~ right-cursor_pos)\n(works within a column set)")
		PushItemWidth(GetContentRegionAvail().x * 0.5)
		DragFloat("float##3a", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
		if layoutState.show_indented_items {
			Indent(0)
			DragFloat("float (indented)##3b", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
			Unindent(0)
		}
		PopItemWidth()

		Text("SetNextItemWidth/PushItemWidth(-FLT_MIN)")
		SameLine(0, 0)
		HelpMarker("Align to right edge")
		PushItemWidth(-FLT_MIN)
		DragFloat("##float5a", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
		if layoutState.show_indented_items {
			Indent(0)
			DragFloat("float (indented)##5b", &layoutState.width_f, 1, 0, 0, "%.3f", 0)
			Unindent(0)
		}
		PopItemWidth()

		TreePop()
	}

	if TreeNode("Basic Horizontal Layout") {
		TextWrapped("(Use ImGui::SameLine() to keep adding items to the right of the preceding item)")

		// Text
		Text("Two items: Hello")
		SameLine(0, 0)
		yellow := ImVec4{1, 1, 0, 1}
		TextColored(&yellow, "Sailor")

		// Adjust spacing
		Text("More spacing: Hello")
		SameLine(0, 20)
		TextColored(&yellow, "Sailor")

		// Button
		AlignTextToFramePadding()
		Text("Normal buttons")
		SameLine(0, 0)
		Button("Banana")
		SameLine(0, 0)
		Button("Apple")
		SameLine(0, 0)
		Button("Corniflower")

		// Button
		Text("Small buttons")
		SameLine(0, 0)
		SmallButton("Like this one")
		SameLine(0, 0)
		Text("can fit within a text block.")

		// Aligned to arbitrary position
		Text("Aligned")
		SameLine(150, 0)
		Text("x=150")
		SameLine(300, 0)
		Text("x=300")
		Text("Aligned")
		SameLine(150, 0)
		SmallButton("x=150")
		SameLine(300, 0)
		SmallButton("x=300")

		// Checkbox
		Checkbox("My", &layoutState.c1)
		SameLine(0, 0)
		Checkbox("Tailor", &layoutState.c2)
		SameLine(0, 0)
		Checkbox("Is", &layoutState.c3)
		SameLine(0, 0)
		Checkbox("Rich", &layoutState.c4)

		// Various
		PushItemWidth(80)
		items := []string{"AAAA", "BBBB", "CCCC", "DDDD"}
		Combo("Combo", &layoutState.item, items, int(len(items)), -1)
		SameLine(0, 0)
		SliderFloat("X", &layoutState.f0, 0.0, 5.0, "%.3f", 0)
		SameLine(0, 0)
		SliderFloat("Y", &layoutState.f1, 0.0, 5.0, "%.3f", 0)
		SameLine(0, 0)
		SliderFloat("Z", &layoutState.f2, 0.0, 5.0, "%.3f", 0)
		PopItemWidth()

		TreePop()
	}

	if TreeNode("Groups") {
		HelpMarker("BeginGroup() basically locks the horizontal position for new line. EndGroup() bundles the whole group so that you can use \"item\" functions such as IsItemHovered()/IsItemActive() or SameLine() etc. on the whole group.")
		BeginGroup()
		{
			BeginGroup()
			Button("AAA")
			SameLine(0, 0)
			Button("BBB")
			SameLine(0, 0)
			BeginGroup()
			Button("CCC")
			Button("DDD")
			EndGroup()
			SameLine(0, 0)
			Button("EEE")
			EndGroup()
			if IsItemHovered(0) {
				SetTooltip("First group hovered")
			}
		}
		// Capture the group size and create widgets based on that size
		size := GetItemRectSize()
		values := layoutState.group_values[:]

		PlotHistogram("##values", values, int(len(values)), 0, "", 0.0, 1.0, size, 4)

		Button("ACTION")
		SameLine(0, 0)
		Button("REACTION")
		SameLine(0, 0)
		Button("DEMO")
		SameLine(0, 0)
		Button("DEMO")
		SameLine(0, 0)
		Button("DEMO")
		EndGroup()
		SameLine(0, 0)

		Button("HIERARCHYHHH")
		SameLine(0, 0)
		BeginGroup()
		Button("Track")
		SameLine(0, 0)
		Button("Move")
		SameLine(0, 0)
		Button("Rotate")
		EndGroup()

		TreePop()
	}

	if TreeNode("Text Baseline Alignment") {
		{
			BulletText("Text baseline:")
			SameLine(0, 0)
			HelpMarker("This is testing the vertical alignment that gets applied on text to keep it aligned with widgets. Lines only composed of text or \"small\" widgets use less vertical space than lines with framed widgets.")
			Indent(0)

			Text("KO Alarm")
				SameLine(0, 0)
				Button("Some framed item")
				SameLine(0, 0)
				Text("Some text")

				// If your line starts with text, call AlignTextToFramePadding() to align text to upcoming widgets.
				AlignTextToFramePadding()
				Text("OK Alarm")
				SameLine(0, 0)
				Button("Some framed item")
				SameLine(0, 0)
				Text("Some text")

			Unindent(0)
		}

		{
			BulletText("Multi-line text:")
			Indent(0)
			Text("One\nTwo\nThree")
			SameLine(0, 0)
			Text("Hello\nWorld")
			SameLine(0, 0)
			Text("Banana")

			Text("Banana")
			SameLine(0, 0)
			Text("Hello\nWorld")
			SameLine(0, 0)
			Text("One\nTwo\nThree")

			Button("HOP##1")
			SameLine(0, 0)
			Text("Banana")
			SameLine(0, 0)
			Text("Hello\nWorld")
			SameLine(0, 0)
			Text("Banana")

			Button("HOP##2")
			SameLine(0, 0)
			Text("Hello\nWorld")
			SameLine(0, 0)
			Text("Banana")
			Unindent(0)
		}

		TreePop()
	}

	if TreeNode("Scrolling") {
		// Vertical scroll functions
		HelpMarker("Use SetScrollHereY() or SetScrollFromPosY() to scroll to a given vertical position.")

		Checkbox("Decoration", &demoState.show_scrolling_decoration)

		if demoState.show_scrolling_decoration {
			SameLine(0, 0)
			Checkbox("Track", &demoState.show_scrolling_track)
			SameLine(0, 0)
			PushItemWidth(100)
			Combo("##combo", &demoState.scrolling_track_item, []string{"track", "track centered"}, 2, -1)
			PopItemWidth()
		}

		if demoState.show_scrolling_track && demoState.show_scrolling_decoration {
			for i := int(0); i < 5; i++ {
				if i > 0 {
					SameLine(0, 0)
				}
				BeginGroup()
				scroll_to_off := Button(fmt.Sprintf("Scroll Offset##%d", i))
				scroll_to_pos := Button(fmt.Sprintf("Scroll Pos##%d", i))
				if scroll_to_off || scroll_to_pos {
					demoState.scrolling_track_item = i
				}
				EndGroup()
			}
		}

		TreePop()
	}
}

// State for ShowDemoWindowPopups
var popupsState struct {
	selected_fish int
	toggles       [5]bool
	context_value float
	name          [32]byte
}

func init() {
	popupsState.selected_fish = -1
	popupsState.toggles[0] = true
	copy(popupsState.name[:], "Label1")
}

func ShowDemoWindowPopups() {
	if !CollapsingHeader("Popups & Modal windows", 0) {
		return
	}

	if TreeNode("Popups") {
		TextWrapped("When a popup is active, it inhibits interacting with windows that are behind the popup. Clicking outside the popup closes it.")

		names := []string{"Bream", "Haddock", "Mackerel", "Pollock", "Tilefish"}

		// Simple selection popup
		if Button("Select..") {
			OpenPopup("my_select_popup", 0)
		}
		SameLine(0, 0)
		if popupsState.selected_fish == -1 {
			TextUnformatted("<None>")
		} else {
			TextUnformatted(names[popupsState.selected_fish])
		}
		if BeginPopup("my_select_popup", 0) {
			Text("Aquarium")
			Separator()
			for i := 0; i < len(names); i++ {
				if Selectable(names[i], false, 0, ImVec2{}) {
					popupsState.selected_fish = int(i)
				}
			}
			EndPopup()
		}

		// Showing a menu with toggles
		if Button("Toggle..") {
			OpenPopup("my_toggle_popup", 0)
		}
		if BeginPopup("my_toggle_popup", 0) {
			for i := 0; i < len(names); i++ {
				MenuItemSelected(names[i], "", &popupsState.toggles[i], true)
			}
			if BeginMenu("Sub-menu", true) {
				MenuItem("Click me", "", nil, true)
				EndMenu()
			}

			Separator()
			Text("Tooltip here")
			if IsItemHovered(0) {
				SetTooltip("I am a tooltip over a popup")
			}

			if Button("Stacked Popup") {
				OpenPopup("another popup", 0)
			}
			if BeginPopup("another popup", 0) {
				for i := 0; i < len(names); i++ {
					MenuItemSelected(names[i], "", &popupsState.toggles[i], true)
				}
				if BeginMenu("Sub-menu", true) {
					MenuItem("Click me", "", nil, true)
					if Button("Stacked Popup") {
						OpenPopup("another popup", 0)
					}
					if BeginPopup("another popup", 0) {
						Text("I am the last one here.")
						EndPopup()
					}
					EndMenu()
				}
				EndPopup()
			}
			EndPopup()
		}

		TreePop()
	}

	if TreeNode("Context menus") {
		HelpMarker("\"Context\" functions are simple helpers to associate a Popup to a given Item or Window identifier.")

		// Example 1
			{
				context_names := []string{"Label1", "Label2", "Label3", "Label4", "Label5"}
				for n := 0; n < 5; n++ {
					Selectable(context_names[n], false, 0, ImVec2{})
					if BeginPopupContextItem("", ImGuiPopupFlags_MouseButtonRight) {
						Text("This a popup for \"%s\"!", context_names[n])
						if Button("Close") {
							CloseCurrentPopup()
						}
						EndPopup()
					}
					if IsItemHovered(0) {
						SetTooltip("Right-click to open popup")
					}
				}
			}

		// Example 2
		{
			HelpMarker("Text() elements don't have stable identifiers so we need to provide one.")
			Text("Value = %.3f <-- (1) right-click this value", popupsState.context_value)
			if BeginPopupContextItem("my popup", ImGuiPopupFlags_MouseButtonRight) {
				if Selectable("Set to zero", false, 0, ImVec2{}) {
					popupsState.context_value = 0.0
				}
				if Selectable("Set to PI", false, 0, ImVec2{}) {
					popupsState.context_value = 3.1415
				}
				SetNextItemWidth(-FLT_MIN)
				DragFloat("##Value", &popupsState.context_value, 0.1, 0.0, 0.0, "%.3f", 0)
				EndPopup()
			}

			Text("(2) Or right-click this text")
			OpenPopupOnItemClick("my popup", ImGuiPopupFlags_MouseButtonRight)

			if Button("(3) Or click this button") {
				OpenPopup("my popup", 0)
			}
		}

		// Example 3
		{
			HelpMarker("Showcase using a popup ID linked to item ID, with the item having a changing label + stable ID using the ### operator.")
			name := string(popupsState.name[:])
			for i, c := range popupsState.name {
				if c == 0 {
					name = string(popupsState.name[:i])
					break
				}
			}
			buf := fmt.Sprintf("Button: %s###Button", name)
			Button(buf)
			if BeginPopupContextItem("", ImGuiPopupFlags_MouseButtonRight) {
				Text("Edit name:")
				nameSlice := popupsState.name[:]
				InputText("##edit", &nameSlice, 0, nil, nil)
				if Button("Close") {
					CloseCurrentPopup()
				}
				EndPopup()
			}
			SameLine(0, 0)
			Text("(<-- right-click here)")
		}

		TreePop()
	}

	if TreeNode("Modals") {
		TextWrapped("Modal windows are like popups but the user cannot close them by clicking outside.")

		if Button("Delete..") {
			OpenPopup("Delete?", 0)
		}

		// Always center this window when appearing
		center := GetMainViewport().GetCenter()
		SetNextWindowPos(&center, ImGuiCond_Appearing, ImVec2{0.5, 0.5})

		if BeginPopupModal("Delete?", nil, ImGuiWindowFlags_AlwaysAutoResize) {
			Text("All those beautiful files will be deleted.\nThis operation cannot be undone!\n\n")
			Separator()

			PushStyleVec(ImGuiStyleVar_FramePadding, ImVec2{0, 0})
			Checkbox("Don't ask me next time", &demoState.no_close)
			PopStyleVar(1)

			if Button("OK") {
				CloseCurrentPopup()
			}
			SetItemDefaultFocus()
			SameLine(0, 0)
			if Button("Cancel") {
				CloseCurrentPopup()
			}
			EndPopup()
		}

		if Button("Stacked modals..") {
			OpenPopup("Stacked 1", 0)
		}
		if BeginPopupModal("Stacked 1", nil, 0) {
			Text("Hello from Stacked The First\nUsing style.Colors[ImGuiCol_ModalWindowDimBg] behind it.")

			if Button("Another modal..") {
				OpenPopup("Stacked 2", 0)
			}
			SetItemDefaultFocus()
			if BeginPopupModal("Stacked 2", nil, 0) {
				Text("Hello from Stacked The Second!")
				if Button("Close") {
					CloseCurrentPopup()
				}
				EndPopup()
			}

			if Button("Close") {
				CloseCurrentPopup()
			}
			EndPopup()
		}

		TreePop()
	}

	if TreeNode("Menus inside a regular window") {
		TextWrapped("Below we are testing adding menu items to a regular window. It's rather unusual but should work!")
		Separator()

		MenuItem("Menu item", "CTRL+M", nil, true)
		if BeginMenu("Menu inside a regular window", true) {
			MenuItem("Submenu item", "", nil, true)
			EndMenu()
		}
		Separator()
		TreePop()
	}
}

// State for ShowDemoWindowMisc
var miscState struct {
	filter ImGuiTextFilter
}

func ShowDemoWindowMisc() {
	if !CollapsingHeader("Filtering", 0) {
		return
	}

	// Helper to display a little (?) mark which shows a tooltip when hovered.
	HelpMarker("Filter usage:\n" +
		"  \"\"         display all lines\n" +
		"  \"xxx\"      display lines containing \"xxx\"\n" +
		"  \"xxx,yyy\"  display lines containing \"xxx\" or \"yyy\"\n" +
		"  \"-xxx\"     hide lines containing \"xxx\"")

	miscState.filter.Draw("Filter (inc,-exc)", 0)
	lines := []string{"aaa1.c", "bbb1.c", "ccc1.c", "aaa2.cpp", "bbb2.cpp", "ccc2.cpp", "abc.h", "hello, world"}
	for i := 0; i < len(lines); i++ {
		if miscState.filter.PassFilter(lines[i]) {
			BulletText("%s", lines[i])
		}
	}
}
