package imgui

//-----------------------------------------------------------------------------
// [SECTION] ImGuiStyle
//-----------------------------------------------------------------------------
// You may modify the ImGui::GetStyle() main instance during initialization and before NewFrame().
// During the frame, use ImGui::PushStyleVar(ImGuiStyleVar_XXXX)/PopStyleVar() to alter the main style values,
// and ImGui::PushStyleColor(ImGuiCol_XXX)/PopStyleColor() for colors.
//-----------------------------------------------------------------------------
type ImGuiStyle struct {
	Alpha                      float    // Global alpha applies to everything in Dear ImGui.
	DisabledAlpha              float    // Additional alpha multiplier applied by BeginDisabled(). Multiply over current value of Alpha.
	WindowPadding              ImVec2   // Padding within a window.
	WindowRounding             float    // Radius of window corners rounding. Set to 0.0f to have rectangular windows. Large values tend to lead to variety of artifacts and are not recommended.
	WindowBorderSize           float    // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	WindowMinSize              ImVec2   // Minimum window size. This is a global setting. If you want to constraint individual windows, use SetNextWindowSizeConstraints().
	WindowTitleAlign           ImVec2   // Alignment for title bar text. Defaults to (0.0,0.5) for left-aligned,vertically centered.
	WindowMenuButtonPosition   ImGuiDir // Side of the collapsing/docking button in the title bar (None/Left/Right). Defaults to ImGuiDir_Left.
	ChildRounding              float    // Radius of child window corners rounding. Set to 0.0f to have rectangular windows.
	ChildBorderSize            float    // Thickness of border around child windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	PopupRounding              float    // Radius of popup window corners rounding. (Note that tooltip windows use WindowRounding)
	PopupBorderSize            float    // Thickness of border around popup/tooltip windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	FramePadding               ImVec2   // Padding within a framed rectangle (used by most widgets).
	FrameRounding              float    // Radius of frame corners rounding. Set to 0.0f to have rectangular frame (used by most widgets).
	FrameBorderSize            float    // Thickness of border around frames. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	ItemSpacing                ImVec2   // Horizontal and vertical spacing between widgets/lines.
	ItemInnerSpacing           ImVec2   // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label).
	CellPadding                ImVec2   // Padding within a table cell
	TouchExtraPadding          ImVec2   // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	IndentSpacing              float    // Horizontal indentation when e.g. entering a tree node. Generally :: (FontSize + FramePadding.x*2).
	ColumnsMinSpacing          float    // Minimum horizontal spacing between two columns. Preferably > (FramePadding.x + 1).
	ScrollbarSize              float    // Width of the vertical scrollbar, Height of the horizontal scrollbar.
	ScrollbarRounding          float    // Radius of grab corners for scrollbar.
	GrabMinSize                float    // Minimum width/height of a grab box for slider/scrollbar.
	GrabRounding               float    // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	LogSliderDeadzone          float    // The size in pixels of the dead-zone around zero on logarithmic sliders that cross zero.
	TabRounding                float    // Radius of upper corners of a tab. Set to 0.0f to have rectangular tabs.
	TabBorderSize              float    // Thickness of border around tabs.
	TabMinWidthForCloseButton  float    // Minimum width for close button to appears on an unselected tab when hovered. Set to 0.0f to always show when hovering, set to FLT_MAX to never show close button unless selected.
	ColorButtonPosition        ImGuiDir // Side of the color button in the ColorEdit4 widget (left/right). Defaults to ImGuiDir_Right.
	ButtonTextAlign            ImVec2   // Alignment of button text when button is larger than text. Defaults to (0.5, 0.5) (centered).
	SelectableTextAlign        ImVec2   // Alignment of selectable text. Defaults to (0.0, 0.0) (top-left aligned). It's generally important to keep this left-aligned if you want to lay multiple items on a same line.
	DisplayWindowPadding       ImVec2   // Window position are clamped to be visible within the display area or monitors by at least this amount. Only applies to regular windows.
	DisplaySafeAreaPadding     ImVec2   // If you cannot see the edges of your screen (e.g. on a TV) increase the safe area padding. Apply to popups/tooltips as well regular windows. NB: Prefer configuring your TV sets correctly!
	MouseCursorScale           float    // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	AntiAliasedLines           bool     // Enable anti-aliased lines/borders. Disable if you are really tight on CPU/GPU. Latched at the beginning of the frame (copied to ImDrawList).
	AntiAliasedLinesUseTex     bool     // Enable anti-aliased lines/borders using textures where possible. Require backend to render with bilinear filtering. Latched at the beginning of the frame (copied to ImDrawList).
	AntiAliasedFill            bool     // Enable anti-aliased edges around filled shapes (rounded rectangles, circles, etc.). Disable if you are really tight on CPU/GPU. Latched at the beginning of the frame (copied to ImDrawList).
	CurveTessellationTol       float    // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	CircleTessellationMaxError float    // Maximum error (in pixels) allowed when using AddCircle()/AddCircleFilled() or drawing rounded corner rectangles with no explicit segment count specified. Decrease for higher quality but more geometry.
	Colors                     [ImGuiCol_COUNT]ImVec4
}

func NewImGuiStyle() ImGuiStyle {
	var style = ImGuiStyle{
		Alpha:                      1.0,              // Global alpha applies to everything in Dear ImGui.
		DisabledAlpha:              0.60,             // Additional alpha multiplier applied by BeginDisabled(). Multiply over current value of Alpha.
		WindowPadding:              ImVec2{8, 8},     // Padding within a window
		WindowRounding:             0.0,              // Radius of window corners rounding. Set to 0.0f to have rectangular windows. Large values tend to lead to variety of artifacts and are not recommended.
		WindowBorderSize:           1.0,              // Thickness of border around windows. Generally set to 0.0f or 1.0f. Other values not well tested.
		WindowMinSize:              ImVec2{32, 32},   // Minimum window size
		WindowTitleAlign:           ImVec2{0.0, 0.5}, // Alignment for title bar text
		WindowMenuButtonPosition:   ImGuiDir_Left,    // Position of the collapsing/docking button in the title bar (left/right). Defaults to ImGuiDir_Left.
		ChildRounding:              0.0,              // Radius of child window corners rounding. Set to 0.0f to have rectangular child windows
		ChildBorderSize:            1.0,              // Thickness of border around child windows. Generally set to 0.0f or 1.0f. Other values not well tested.
		PopupRounding:              0.0,              // Radius of popup window corners rounding. Set to 0.0f to have rectangular child windows
		PopupBorderSize:            1.0,              // Thickness of border around popup or tooltip windows. Generally set to 0.0f or 1.0f. Other values not well tested.
		FramePadding:               ImVec2{4, 3},     // Padding within a framed rectangle (used by most widgets)
		FrameRounding:              0.0,              // Radius of frame corners rounding. Set to 0.0f to have rectangular frames (used by most widgets).
		FrameBorderSize:            0.0,              // Thickness of border around frames. Generally set to 0.0f or 1.0f. Other values not well tested.
		ItemSpacing:                ImVec2{8, 4},     // Horizontal and vertical spacing between widgets/lines
		ItemInnerSpacing:           ImVec2{4, 4},     // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label)
		CellPadding:                ImVec2{4, 2},     // Padding within a table cell
		TouchExtraPadding:          ImVec2{0, 0},     // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
		IndentSpacing:              21.0,             // Horizontal spacing when e.g. entering a tree node. Generally :: (FontSize + FramePadding.x*2).
		ColumnsMinSpacing:          6.0,              // Minimum horizontal spacing between two columns. Preferably > (FramePadding.x + 1).
		ScrollbarSize:              14.0,             // Width of the vertical scrollbar, Height of the horizontal scrollbar
		ScrollbarRounding:          9.0,              // Radius of grab corners rounding for scrollbar
		GrabMinSize:                10.0,             // Minimum width/height of a grab box for slider/scrollbar
		GrabRounding:               0.0,              // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
		LogSliderDeadzone:          4.0,              // The size in pixels of the dead-zone around zero on logarithmic sliders that cross zero.
		TabRounding:                4.0,              // Radius of upper corners of a tab. Set to 0.0f to have rectangular tabs.
		TabBorderSize:              0.0,              // Thickness of border around tabs.
		TabMinWidthForCloseButton:  0.0,              // Minimum width for close button to appears on an unselected tab when hovered. Set to 0.0f to always show when hovering, set to FLT_MAX to never show close button unless selected.
		ColorButtonPosition:        ImGuiDir_Right,   // Side of the color button in the ColorEdit4 widget (left/right). Defaults to ImGuiDir_Right.
		ButtonTextAlign:            ImVec2{0.5, 0.5}, // Alignment of button text when button is larger than text.
		SelectableTextAlign:        ImVec2{0.0, 0.0}, // Alignment of selectable text. Defaults to (0.0, 0.0) (top-left aligned). It's generally important to keep this left-aligned if you want to lay multiple items on a same line.
		DisplayWindowPadding:       ImVec2{19, 19},   // Window position are clamped to be visible within the display area or monitors by at least this amount. Only applies to regular windows.
		DisplaySafeAreaPadding:     ImVec2{3, 3},     // If you cannot see the edge of your screen (e.g. on a TV) increase the safe area padding. Covers popups/tooltips as well regular windows.
		MouseCursorScale:           1.0,              // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
		AntiAliasedLines:           true,             // Enable anti-aliased lines/borders. Disable if you are really tight on CPU/GPU.
		AntiAliasedLinesUseTex:     true,             // Enable anti-aliased lines/borders using textures where possible. Require backend to render with bilinear filtering.
		AntiAliasedFill:            true,             // Enable anti-aliased filled shapes (rounded rectangles, circles, etc.).
		CurveTessellationTol:       1.25,             // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
		CircleTessellationMaxError: 0.30,             // Maximum error (in pixels) allowed when using AddCircle()/AddCircleFilled() or drawing rounded corner rectangles with no explicit segment count specified. Decrease for higher quality but more geometry.
	}
	StyleColorsDark(&style)
	return style
}

func (ImGuiStyle) ScaleAllSizes(scale_factor float) { panic("not implemented") }

func StyleColorsDark(dst *ImGuiStyle) {
	var style *ImGuiStyle = dst
	if dst == nil {
		style = GetStyle()
	}
	var colors = &style.Colors

	colors[ImGuiCol_Text] = ImVec4{1.00, 1.00, 1.00, 1.00}
	colors[ImGuiCol_TextDisabled] = ImVec4{0.50, 0.50, 0.50, 1.00}
	colors[ImGuiCol_WindowBg] = ImVec4{0.06, 0.06, 0.06, 0.94}
	colors[ImGuiCol_ChildBg] = ImVec4{0.00, 0.00, 0.00, 0.00}
	colors[ImGuiCol_PopupBg] = ImVec4{0.08, 0.08, 0.08, 0.94}
	colors[ImGuiCol_Border] = ImVec4{0.43, 0.43, 0.50, 0.50}
	colors[ImGuiCol_BorderShadow] = ImVec4{0.00, 0.00, 0.00, 0.00}
	colors[ImGuiCol_FrameBg] = ImVec4{0.16, 0.29, 0.48, 0.54}
	colors[ImGuiCol_FrameBgHovered] = ImVec4{0.26, 0.59, 0.98, 0.40}
	colors[ImGuiCol_FrameBgActive] = ImVec4{0.26, 0.59, 0.98, 0.67}
	colors[ImGuiCol_TitleBg] = ImVec4{0.04, 0.04, 0.04, 1.00}
	colors[ImGuiCol_TitleBgActive] = ImVec4{0.16, 0.29, 0.48, 1.00}
	colors[ImGuiCol_TitleBgCollapsed] = ImVec4{0.00, 0.00, 0.00, 0.51}
	colors[ImGuiCol_MenuBarBg] = ImVec4{0.14, 0.14, 0.14, 1.00}
	colors[ImGuiCol_ScrollbarBg] = ImVec4{0.02, 0.02, 0.02, 0.53}
	colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.31, 0.31, 0.31, 1.00}
	colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.41, 0.41, 0.41, 1.00}
	colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.51, 0.51, 0.51, 1.00}
	colors[ImGuiCol_CheckMark] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_SliderGrab] = ImVec4{0.24, 0.52, 0.88, 1.00}
	colors[ImGuiCol_SliderGrabActive] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_Button] = ImVec4{0.26, 0.59, 0.98, 0.40}
	colors[ImGuiCol_ButtonHovered] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_ButtonActive] = ImVec4{0.06, 0.53, 0.98, 1.00}
	colors[ImGuiCol_Header] = ImVec4{0.26, 0.59, 0.98, 0.31}
	colors[ImGuiCol_HeaderHovered] = ImVec4{0.26, 0.59, 0.98, 0.80}
	colors[ImGuiCol_HeaderActive] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_Separator] = colors[ImGuiCol_Border]
	colors[ImGuiCol_SeparatorHovered] = ImVec4{0.10, 0.40, 0.75, 0.78}
	colors[ImGuiCol_SeparatorActive] = ImVec4{0.10, 0.40, 0.75, 1.00}
	colors[ImGuiCol_ResizeGrip] = ImVec4{0.26, 0.59, 0.98, 0.20}
	colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.26, 0.59, 0.98, 0.67}
	colors[ImGuiCol_ResizeGripActive] = ImVec4{0.26, 0.59, 0.98, 0.95}
	colors[ImGuiCol_Tab] = ImLerpVec4(&colors[ImGuiCol_Header], &colors[ImGuiCol_TitleBgActive], 0.80)
	colors[ImGuiCol_TabHovered] = colors[ImGuiCol_HeaderHovered]
	colors[ImGuiCol_TabActive] = ImLerpVec4(&colors[ImGuiCol_HeaderActive], &colors[ImGuiCol_TitleBgActive], 0.60)
	colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&colors[ImGuiCol_Tab], &colors[ImGuiCol_TitleBg], 0.80)
	colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&colors[ImGuiCol_TabActive], &colors[ImGuiCol_TitleBg], 0.40)
	colors[ImGuiCol_PlotLines] = ImVec4{0.61, 0.61, 0.61, 1.00}
	colors[ImGuiCol_PlotLinesHovered] = ImVec4{1.00, 0.43, 0.35, 1.00}
	colors[ImGuiCol_PlotHistogram] = ImVec4{0.90, 0.70, 0.00, 1.00}
	colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.00, 0.60, 0.00, 1.00}
	colors[ImGuiCol_TableHeaderBg] = ImVec4{0.19, 0.19, 0.20, 1.00}
	colors[ImGuiCol_TableBorderStrong] = ImVec4{0.31, 0.31, 0.35, 1.00} // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableBorderLight] = ImVec4{0.23, 0.23, 0.25, 1.00}  // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableRowBg] = ImVec4{0.00, 0.00, 0.00, 0.00}
	colors[ImGuiCol_TableRowBgAlt] = ImVec4{1.00, 1.00, 1.00, 0.06}
	colors[ImGuiCol_TextSelectedBg] = ImVec4{0.26, 0.59, 0.98, 0.35}
	colors[ImGuiCol_DragDropTarget] = ImVec4{1.00, 1.00, 0.00, 0.90}
	colors[ImGuiCol_NavHighlight] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_NavWindowingHighlight] = ImVec4{1.00, 1.00, 1.00, 0.70}
	colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.80, 0.80, 0.80, 0.20}
	colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.80, 0.80, 0.80, 0.35}
}
