package imgui

import (
	"reflect"
	"unsafe"

	"github.com/Splizard/imgui/golang"
)

// Enumeration for PushStyleVar() / PopStyleVar() to temporarily modify the ImGuiStyle structure.
//   - The const (
//     During initialization or between frames, feel free to just poke into ImGuiStyle directly.
//   - Tip: Use your programming IDE navigation facilities on the names in the _second column_ below to find the actual members and their description.
//     In Visual Studio IDE: CTRL+comma ("Edit.NavigateTo") can follow symbols in comments, whereas CTRL+F12 ("Edit.GoToImplementation") cannot.
//     With Visual Assist installed: ALT+G ("VAssistX.GoToImplementation") can also follow symbols in comments.
//   - When changing this enum, you need to update the associated internal table GStyleVarInfo[] accordingly. This is where we link const (
const (
	ImGuiStyleVar_Alpha               ImGuiStyleVar = iota // float     Alpha
	ImGuiStyleVar_DisabledAlpha                            // float     DisabledAlpha
	ImGuiStyleVar_WindowPadding                            // ImVec2    WindowPadding
	ImGuiStyleVar_WindowRounding                           // float     WindowRounding
	ImGuiStyleVar_WindowBorderSize                         // float     WindowBorderSize
	ImGuiStyleVar_WindowMinSize                            // ImVec2    WindowMinSize
	ImGuiStyleVar_WindowTitleAlign                         // ImVec2    WindowTitleAlign
	ImGuiStyleVar_ChildRounding                            // float     ChildRounding
	ImGuiStyleVar_ChildBorderSize                          // float     ChildBorderSize
	ImGuiStyleVar_PopupRounding                            // float     PopupRounding
	ImGuiStyleVar_PopupBorderSize                          // float     PopupBorderSize
	ImGuiStyleVar_FramePadding                             // ImVec2    FramePadding
	ImGuiStyleVar_FrameRounding                            // float     FrameRounding
	ImGuiStyleVar_FrameBorderSize                          // float     FrameBorderSize
	ImGuiStyleVar_ItemSpacing                              // ImVec2    ItemSpacing
	ImGuiStyleVar_ItemInnerSpacing                         // ImVec2    ItemInnerSpacing
	ImGuiStyleVar_IndentSpacing                            // float     IndentSpacing
	ImGuiStyleVar_CellPadding                              // ImVec2    CellPadding
	ImGuiStyleVar_ScrollbarSize                            // float     ScrollbarSize
	ImGuiStyleVar_ScrollbarRounding                        // float     ScrollbarRounding
	ImGuiStyleVar_GrabMinSize                              // float     GrabMinSize
	ImGuiStyleVar_GrabRounding                             // float     GrabRounding
	ImGuiStyleVar_TabRounding                              // float     TabRounding
	ImGuiStyleVar_ButtonTextAlign                          // ImVec2    ButtonTextAlign
	ImGuiStyleVar_SelectableTextAlign                      // ImVec2    SelectableTextAlign
	ImGuiStyleVar_COUNT
)

type ImGuiStyleMod struct {
	VarIdx      ImGuiStyleVar
	BackupValue [2]int
}

func NewImGuiStyleModInt(idx ImGuiStyleVar, v int) ImGuiStyleMod {
	return ImGuiStyleMod{VarIdx: idx, BackupValue: [2]int{v, 0}}
}

func NewImGuiStyleModFloat(idx ImGuiStyleVar, v float32) ImGuiStyleMod {
	return ImGuiStyleMod{VarIdx: idx, BackupValue: [2]int{*(*int)(unsafe.Pointer(&v)), 0}}
}

func NewImGuiStyleModVec(idx ImGuiStyleVar, v ImVec2) ImGuiStyleMod {
	return ImGuiStyleMod{VarIdx: idx, BackupValue: [2]int{*(*int)(unsafe.Pointer(&v.x)), *(*int)(unsafe.Pointer(&v.y))}}
}

func (mod ImGuiStyleMod) Int() int {
	return mod.BackupValue[0]
}

func (mod ImGuiStyleMod) Float() float {
	return *(*float)(unsafe.Pointer(&mod.BackupValue[0]))
}

func (mod ImGuiStyleMod) Vec2() ImVec2 {
	return ImVec2{*(*float)(unsafe.Pointer(&mod.BackupValue[0])), *(*float)(unsafe.Pointer(&mod.BackupValue[1]))}
}

// -----------------------------------------------------------------------------
// [SECTION] ImGuiStyle
// -----------------------------------------------------------------------------
// You may modify the ImGui::GetStyle() main instance during initialization and before NewFrame().
// During the frame, use ImGui::PushStyleVar(ImGuiStyleVar_XXXX)/PopStyleVar() to alter the main style values,
// and ImGui::PushStyleColor(ImGuiCol_XXX)/PopStyleColor() for colors.
// -----------------------------------------------------------------------------
type ImGuiStyle struct {
	Alpha               float  // Global alpha applies to everything in Dear ImGui.
	DisabledAlpha       float  // Additional alpha multiplier applied by BeginDisabled(). Multiply over current value of Alpha.
	WindowPadding       ImVec2 // Padding within a window.
	WindowRounding      float  // Radius of window corners rounding. Set to 0.0f to have rectangular windows. Large values tend to lead to variety of artifacts and are not recommended.
	WindowBorderSize    float  // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	WindowMinSize       ImVec2 // Minimum window size. This is a global setting. If you want to constraint individual windows, use SetNextWindowSizeConstraints().
	WindowTitleAlign    ImVec2 // Alignment for title bar text. Defaults to (0.0,0.5) for left-aligned,vertically centered.
	ChildRounding       float  // Radius of child window corners rounding. Set to 0.0f to have rectangular windows.
	ChildBorderSize     float  // Thickness of border around child windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	PopupRounding       float  // Radius of popup window corners rounding. (Note that tooltip windows use WindowRounding)
	PopupBorderSize     float  // Thickness of border around popup/tooltip windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	FramePadding        ImVec2 // Padding within a framed rectangle (used by most widgets).
	FrameRounding       float  // Radius of frame corners rounding. Set to 0.0f to have rectangular frame (used by most widgets).
	FrameBorderSize     float  // Thickness of border around frames. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	ItemSpacing         ImVec2 // Horizontal and vertical spacing between widgets/lines.
	ItemInnerSpacing    ImVec2 // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label).
	IndentSpacing       float  // Horizontal indentation when e.g. entering a tree node. Generally :: (FontSize + FramePadding.x*2).
	CellPadding         ImVec2 // Padding within a table cell
	ScrollbarSize       float  // Width of the vertical scrollbar, Height of the horizontal scrollbar.
	ScrollbarRounding   float  // Radius of grab corners for scrollbar.
	GrabMinSize         float  // Minimum width/height of a grab box for slider/scrollbar.
	GrabRounding        float  // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	TabRounding         float  // Radius of upper corners of a tab. Set to 0.0f to have rectangular tabs.
	ButtonTextAlign     ImVec2 // Alignment of button text when button is larger than text. Defaults to (0.5, 0.5) (centered).
	SelectableTextAlign ImVec2 // Alignment of selectable text. Defaults to (0.0, 0.0) (top-left aligned). It's generally important to keep this left-aligned if you want to lay multiple items on a same line.

	ColumnsMinSpacing          float    // Minimum horizontal spacing between two columns. Preferably > (FramePadding.x + 1).
	LogSliderDeadzone          float    // The size in pixels of the dead-zone around zero on logarithmic sliders that cross zero.
	TabBorderSize              float    // Thickness of border around tabs.
	TabMinWidthForCloseButton  float    // Minimum width for close button to appears on an unselected tab when hovered. Set to 0.0f to always show when hovering, set to FLT_MAX to never show close button unless selected.
	ColorButtonPosition        ImGuiDir // Side of the color button in the ColorEdit4 widget (left/right). Defaults to ImGuiDir_Right.
	DisplayWindowPadding       ImVec2   // Window position are clamped to be visible within the display area or monitors by at least this amount. Only applies to regular windows.
	DisplaySafeAreaPadding     ImVec2   // If you cannot see the edges of your screen (e.g. on a TV) increase the safe area padding. Apply to popups/tooltips as well regular windows. NB: Prefer configuring your TV sets correctly!
	MouseCursorScale           float    // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	AntiAliasedLines           bool     // Enable anti-aliased lines/borders. Disable if you are really tight on CPU/GPU. Latched at the beginning of the frame (copied to ImDrawList).
	AntiAliasedLinesUseTex     bool     // Enable anti-aliased lines/borders using textures where possible. Require backend to render with bilinear filtering. Latched at the beginning of the frame (copied to ImDrawList).
	AntiAliasedFill            bool     // Enable anti-aliased edges around filled shapes (rounded rectangles, circles, etc.). Disable if you are really tight on CPU/GPU. Latched at the beginning of the frame (copied to ImDrawList).
	CurveTessellationTol       float    // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	CircleTessellationMaxError float    // Maximum error (in pixels) allowed when using AddCircle()/AddCircleFilled() or drawing rounded corner rectangles with no explicit segment count specified. Decrease for higher quality but more geometry.

	TouchExtraPadding        ImVec2   // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	WindowMenuButtonPosition ImGuiDir // Side of the collapsing/docking button in the title bar (None/Left/Right). Defaults to ImGuiDir_Left.
	Colors                   [ImGuiCol_COUNT]ImVec4
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

// access the Style structure (colors, sizes). Always use PushStyleCol(), PushStyleVar() to modify style mid-frame!
func GetStyle() *ImGuiStyle {
	IM_ASSERT_USER_ERROR(g != nil, "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?")
	return &g.Style
}

// To scale your entire UI (e.g. if you want your app to use High DPI or generally be DPI aware) you may use this helper function. Scaling the fonts is done separately and is up to you.
// Important: This operation is lossy because we round all sizes to integer. If you need to change your scale multiples, call this over a freshly initialized ImGuiStyle structure rather than scaling multiple times.
func (style *ImGuiStyle) ScaleAllSizes(scale_factor float) {
	winpad := style.WindowPadding.Scale(scale_factor)
	style.WindowPadding = *ImFloorVec(&winpad)
	style.WindowRounding = ImFloor(style.WindowRounding * scale_factor)
	minsize := style.WindowMinSize.Scale(scale_factor)
	style.WindowMinSize = *ImFloorVec(&minsize)
	style.ChildRounding = ImFloor(style.ChildRounding * scale_factor)
	style.PopupRounding = ImFloor(style.PopupRounding * scale_factor)
	framepad := style.FramePadding.Scale(scale_factor)
	style.FramePadding = *ImFloorVec(&framepad)
	style.FrameRounding = ImFloor(style.FrameRounding * scale_factor)
	itemspacing := style.ItemSpacing.Scale(scale_factor)
	style.ItemSpacing = *ImFloorVec(&itemspacing)
	innerspacing := style.ItemInnerSpacing.Scale(scale_factor)
	style.ItemInnerSpacing = *ImFloorVec(&innerspacing)
	cellpad := style.CellPadding.Scale(scale_factor)
	style.CellPadding = *ImFloorVec(&cellpad)
	touchxtra := style.TouchExtraPadding.Scale(scale_factor)
	style.TouchExtraPadding = *ImFloorVec(&touchxtra)
	style.IndentSpacing = ImFloor(style.IndentSpacing * scale_factor)
	style.ColumnsMinSpacing = ImFloor(style.ColumnsMinSpacing * scale_factor)
	style.ScrollbarSize = ImFloor(style.ScrollbarSize * scale_factor)
	style.ScrollbarRounding = ImFloor(style.ScrollbarRounding * scale_factor)
	style.GrabMinSize = ImFloor(style.GrabMinSize * scale_factor)
	style.GrabRounding = ImFloor(style.GrabRounding * scale_factor)
	style.LogSliderDeadzone = ImFloor(style.LogSliderDeadzone * scale_factor)
	style.TabRounding = ImFloor(style.TabRounding * scale_factor)
	if style.TabMinWidthForCloseButton != FLT_MAX {
		style.TabMinWidthForCloseButton = ImFloor(style.TabMinWidthForCloseButton * scale_factor)
	}
	diswinpad := style.DisplayWindowPadding.Scale(scale_factor)
	style.DisplayWindowPadding = *ImFloorVec(&diswinpad)
	safepad := style.DisplaySafeAreaPadding.Scale(scale_factor)
	style.DisplaySafeAreaPadding = *ImFloorVec(&safepad)
	style.MouseCursorScale = ImFloor(style.MouseCursorScale * scale_factor)
}

// modify a style variable float. always use this if you modify the style after NewFrame().
func PushStyleFloat(idx ImGuiStyleVar, val float) {
	var pvar = reflect.ValueOf(&g.Style).Elem().Field(golang.Int(idx)).Addr().Interface().(*float)
	g.StyleVarStack = append(g.StyleVarStack, NewImGuiStyleModFloat(idx, *pvar))
	*pvar = val
}

// modify a style variable ImVec2. always use this if you modify the style after NewFrame().
func PushStyleVec(idx ImGuiStyleVar, val ImVec2) {
	var pvar = reflect.ValueOf(&g.Style).Elem().Field(golang.Int(idx)).Addr().Interface().(*ImVec2)
	g.StyleVarStack = append(g.StyleVarStack, NewImGuiStyleModVec(idx, *pvar))
	*pvar = val
}

func PopStyleVar(count int /*= 1*/) {
	for count > 0 {
		// We avoid a generic memcpy(data, &backup.Backup.., GDataTypeSize[info.Type] * info.Count), the overhead in Debug is not worth it.
		var backup = &g.StyleVarStack[len(g.StyleVarStack)-1]

		field := reflect.ValueOf(&g.Style).Elem().Field(golang.Int(backup.VarIdx))
		switch field.Type() {
		case reflect.TypeOf(ImVec2{}):
			field.Set(reflect.ValueOf(backup.Vec2()))
		case reflect.TypeOf(float(0)):
			field.Set(reflect.ValueOf(backup.Float()))
		case reflect.TypeOf(int(0)):
			field.Set(reflect.ValueOf(backup.Int()))
		}
		g.StyleVarStack = g.StyleVarStack[:len(g.StyleVarStack)-1]
		count--
	}
}

// retrieve given style color with style alpha applied and optional extra alpha multiplier, packed as a 32-bit value suitable for ImDrawList
func GetColorU32FromID(idx ImGuiCol, alpha_mul float /*= 1.0*/) ImU32 {
	var style = g.Style
	var c = style.Colors[idx]
	c.w *= style.Alpha * alpha_mul
	return ColorConvertFloat4ToU32(c)
}

// retrieve given color with style alpha applied, packed as a 32-bit value suitable for ImDrawList
func GetColorU32FromVec(col ImVec4) ImU32 {
	var style = g.Style
	var c = col
	c.w *= style.Alpha
	return ColorConvertFloat4ToU32(c)
}

// retrieve given color with style alpha applied, packed as a 32-bit value suitable for ImDrawList
func GetColorU32FromInt(col ImU32) ImU32 {
	var style = g.Style
	if style.Alpha >= 1.0 {
		return col
	}
	var a = (col & IM_COL32_A_MASK) >> IM_COL32_A_SHIFT
	a = (ImU32)(float(a) * style.Alpha) // We don't need to clamp 0..255 because Style.Alpha is in 0..1 range.
	return (col &^ IM_COL32_A_MASK) | (a << IM_COL32_A_SHIFT)
}

// retrieve style color as stored in ImGuiStyle structure. use to feed back into PushStyleColor(), otherwise use GetColorU32() to get style color with style alpha baked in.
func GetStyleColorVec4(idx ImGuiCol) *ImVec4 {
	var style = g.Style
	return &style.Colors[idx]
}

// FIXME: This may incur a round-trip (if the end user got their data from a float4) but eventually we aim to store the in-flight colors as ImU32
func PushStyleColorInt(idx ImGuiCol, col ImU32) {
	var backup ImGuiColorMod
	backup.Col = idx
	backup.BackupValue = g.Style.Colors[idx]
	g.ColorStack = append(g.ColorStack, backup)
	g.Style.Colors[idx] = ColorConvertU32ToFloat4(col)
}

func PushStyleColorVec(idx ImGuiCol, col *ImVec4) {
	var backup ImGuiColorMod
	backup.Col = idx
	backup.BackupValue = g.Style.Colors[idx]
	g.ColorStack = append(g.ColorStack, backup)
	g.Style.Colors[idx] = *col
}

func PopStyleColor(count int /*= 1*/) {
	for count > 0 {
		var backup = &g.ColorStack[len(g.ColorStack)-1]
		g.Style.Colors[backup.Col] = backup.BackupValue
		g.ColorStack = g.ColorStack[:len(g.ColorStack)-1]
		count--
	}
}

// get a string corresponding to the enum value (for display, saving, etc.).
func GetStyleColorName(idx ImGuiCol) string {
	switch idx {
	case ImGuiCol_Text:
		return "Text"
	case ImGuiCol_TextDisabled:
		return "TextDisabled"
	case ImGuiCol_WindowBg:
		return "WindowBg"
	case ImGuiCol_ChildBg:
		return "ChildBg"
	case ImGuiCol_PopupBg:
		return "PopupBg"
	case ImGuiCol_Border:
		return "Border"
	case ImGuiCol_BorderShadow:
		return "BorderShadow"
	case ImGuiCol_FrameBg:
		return "FrameBg"
	case ImGuiCol_FrameBgHovered:
		return "FrameBgHovered"
	case ImGuiCol_FrameBgActive:
		return "FrameBgActive"
	case ImGuiCol_TitleBg:
		return "TitleBg"
	case ImGuiCol_TitleBgActive:
		return "TitleBgActive"
	case ImGuiCol_TitleBgCollapsed:
		return "TitleBgCollapsed"
	case ImGuiCol_MenuBarBg:
		return "MenuBarBg"
	case ImGuiCol_ScrollbarBg:
		return "ScrollbarBg"
	case ImGuiCol_ScrollbarGrab:
		return "ScrollbarGrab"
	case ImGuiCol_ScrollbarGrabHovered:
		return "ScrollbarGrabHovered"
	case ImGuiCol_ScrollbarGrabActive:
		return "ScrollbarGrabActive"
	case ImGuiCol_CheckMark:
		return "CheckMark"
	case ImGuiCol_SliderGrab:
		return "SliderGrab"
	case ImGuiCol_SliderGrabActive:
		return "SliderGrabActive"
	case ImGuiCol_Button:
		return "Button"
	case ImGuiCol_ButtonHovered:
		return "ButtonHovered"
	case ImGuiCol_ButtonActive:
		return "ButtonActive"
	case ImGuiCol_Header:
		return "Header"
	case ImGuiCol_HeaderHovered:
		return "HeaderHovered"
	case ImGuiCol_HeaderActive:
		return "HeaderActive"
	case ImGuiCol_Separator:
		return "Separator"
	case ImGuiCol_SeparatorHovered:
		return "SeparatorHovered"
	case ImGuiCol_SeparatorActive:
		return "SeparatorActive"
	case ImGuiCol_ResizeGrip:
		return "ResizeGrip"
	case ImGuiCol_ResizeGripHovered:
		return "ResizeGripHovered"
	case ImGuiCol_ResizeGripActive:
		return "ResizeGripActive"
	case ImGuiCol_Tab:
		return "Tab"
	case ImGuiCol_TabHovered:
		return "TabHovered"
	case ImGuiCol_TabActive:
		return "TabActive"
	case ImGuiCol_TabUnfocused:
		return "TabUnfocused"
	case ImGuiCol_TabUnfocusedActive:
		return "TabUnfocusedActive"
	case ImGuiCol_PlotLines:
		return "PlotLines"
	case ImGuiCol_PlotLinesHovered:
		return "PlotLinesHovered"
	case ImGuiCol_PlotHistogram:
		return "PlotHistogram"
	case ImGuiCol_PlotHistogramHovered:
		return "PlotHistogramHovered"
	case ImGuiCol_TableHeaderBg:
		return "TableHeaderBg"
	case ImGuiCol_TableBorderStrong:
		return "TableBorderStrong"
	case ImGuiCol_TableBorderLight:
		return "TableBorderLight"
	case ImGuiCol_TableRowBg:
		return "TableRowBg"
	case ImGuiCol_TableRowBgAlt:
		return "TableRowBgAlt"
	case ImGuiCol_TextSelectedBg:
		return "TextSelectedBg"
	case ImGuiCol_DragDropTarget:
		return "DragDropTarget"
	case ImGuiCol_NavHighlight:
		return "NavHighlight"
	case ImGuiCol_NavWindowingHighlight:
		return "NavWindowingHighlight"
	case ImGuiCol_NavWindowingDimBg:
		return "NavWindowingDimBg"
	case ImGuiCol_ModalWindowDimBg:
		return "ModalWindowDimBg"
	}
	IM_ASSERT(false)
	return "Unknown"
}
