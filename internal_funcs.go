package imgui

// Fonts, drawing

func getForegroundDrawList(window *ImGuiWindow) *ImDrawList      { return GetForegroundDrawList(nil) } // This seemingly unnecessary wrapper simplifies compatibility between the 'master' and 'docking' branches.
func getBackgroundDrawList(t *ImGuiViewport) *ImDrawList         { panic("not implemented") }          // get background draw list for the given viewport. this draw list will be the first rendering one. Useful to quickly draw shapes/text behind dear imgui contents.
func GetForegroundDrawListViewport(t *ImGuiViewport) *ImDrawList { panic("not implemented") }          // get foreground draw list for the given viewport. this draw list will be the last rendered one. Useful to quickly draw shapes/text over dear imgui contents.

// Call context hooks (used by e.g. test engine)
// We assume a small number of hooks so all stored in same array
func CallContextHooks(ctx *ImGuiContext, hook_type ImGuiContextHookType) {
	var g = ctx
	for n := range g.Hooks {
		if g.Hooks[n].Type == hook_type {
			g.Hooks[n].Callback(g, &g.Hooks[n])
		}
	}
}

// Basic Accessors
func GetItemID() ImGuiID { var g *ImGuiContext = GImGui; return g.LastItemData.ID } // Get ID of last item (~~ often same ImGui::GetID(label) beforehand)
func GetItemStatusFlags() ImGuiItemStatusFlags {
	var g *ImGuiContext = GImGui
	return g.LastItemData.StatusFlags
}
func GetItemFlags() ImGuiItemFlags { var g *ImGuiContext = GImGui; return g.LastItemData.InFlags }
func GetActiveID() ImGuiID         { var g *ImGuiContext = GImGui; return g.ActiveId }
func GetFocusID() ImGuiID          { var g *ImGuiContext = GImGui; return g.NavId }

func SetActiveID(id ImGuiID, window *ImGuiWindow) {
	var g = GImGui
	g.ActiveIdIsJustActivated = (g.ActiveId != id)
	if g.ActiveIdIsJustActivated {
		g.ActiveIdTimer = 0.0
		g.ActiveIdHasBeenPressedBefore = false
		g.ActiveIdHasBeenEditedBefore = false
		g.ActiveIdMouseButton = -1
		if id != 0 {
			g.LastActiveId = id
			g.LastActiveIdTimer = 0.0
		}
	}
	g.ActiveId = id
	g.ActiveIdAllowOverlap = false
	g.ActiveIdNoClearOnFocusLoss = false
	g.ActiveIdWindow = window
	g.ActiveIdHasBeenEditedThisFrame = false
	if id != 0 {
		g.ActiveIdIsAlive = id
		if g.NavActivateId == id || g.NavInputId == id || g.NavJustTabbedId == id || g.NavJustMovedToId == id {
			g.ActiveIdSource = ImGuiInputSource_Nav
		} else {
			g.ActiveIdSource = ImGuiInputSource_Mouse
		}
	}

	// Clear declaration of inputs claimed by the widget
	// (Please note that this is WIP and not all keys/inputs are thoroughly declared by all widgets yet)
	g.ActiveIdUsingMouseWheel = false
	g.ActiveIdUsingNavDirMask = 0x00
	g.ActiveIdUsingNavInputMask = 0x00
	g.ActiveIdUsingKeyInputMask = 0x00
}

func SetFocusID(id ImGuiID, window *ImGuiWindow) {
	var g = GImGui
	IM_ASSERT(id != 0)

	// Assume that SetFocusID() is called in the context where its window.DC.NavLayerCurrent and window.DC.NavFocusScopeIdCurrent are valid.
	// Note that window may be != g.CurrentWindow (e.g. SetFocusID call in InputTextEx for multi-line text)
	var nav_layer ImGuiNavLayer = window.DC.NavLayerCurrent
	if g.NavWindow != window {
		g.NavInitRequest = false
	}
	g.NavWindow = window
	g.NavId = id
	g.NavLayer = nav_layer
	g.NavFocusScopeId = window.DC.NavFocusScopeIdCurrent
	window.NavLastIds[nav_layer] = id
	if g.LastItemData.ID == id {
		window.NavRectRel[nav_layer] = ImRect{g.LastItemData.NavRect.Min.Sub(window.Pos), g.LastItemData.NavRect.Max.Sub(window.Pos)}
	}

	if g.ActiveIdSource == ImGuiInputSource_Nav {
		g.NavDisableMouseHover = true
	} else {
		g.NavDisableHighlight = true
	}
}

func ClearActiveID() {
	SetActiveID(0, nil) // g.ActiveId = 0
}

func GetHoveredID() ImGuiID {
	var g = GImGui
	if g.HoveredId != 0 {
		return g.HoveredId
	}
	return g.HoveredIdPreviousFrame
}

func SetHoveredID(id ImGuiID) {
	var g = GImGui
	g.HoveredId = id
	g.HoveredIdAllowOverlap = false
	g.HoveredIdUsingMouseWheel = false
	if id != 0 && g.HoveredIdPreviousFrame != id {
		g.HoveredIdTimer = 0
		g.HoveredIdNotActiveTimer = 0.0
	}
}

func KeepAliveID(id ImGuiID) {
	var g = GImGui
	if g.ActiveId == id {
		g.ActiveIdIsAlive = id
	}
	if g.ActiveIdPreviousFrame == id {
		g.ActiveIdPreviousFrameIsAlive = true
	}
}

// Mark data associated to given item as "edited", used by IsItemDeactivatedAfterEdit() function.
func MarkItemEdited(id ImGuiID) {
	// This marking is solely to be able to provide info for IsItemDeactivatedAfterEdit().
	// ActiveId might have been released by the time we call this (as in the typical press/release button behavior) but still need need to fill the data.
	var g = GImGui
	IM_ASSERT(g.ActiveId == id || g.ActiveId == 0 || g.DragDropActive)

	//IM_ASSERT(g.CurrentWindow.DC.LastItemId == id);
	g.ActiveIdHasBeenEditedThisFrame = true
	g.ActiveIdHasBeenEditedBefore = true
	g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_Edited
}

// Parameter stacks
func PushItemFlag(option ImGuiItemFlags, enabled bool) {
	var g = GImGui
	var item_flags ImGuiItemFlags = g.CurrentItemFlags
	IM_ASSERT(item_flags == g.ItemFlagsStack[len(g.ItemFlagsStack)-1])
	if enabled {
		item_flags |= option
	} else {
		item_flags &= ^option
	}
	g.CurrentItemFlags = item_flags
	g.ItemFlagsStack = append(g.ItemFlagsStack, item_flags)
}

func PopItemFlag() {
	var g = GImGui
	IM_ASSERT(len(g.ItemFlagsStack) > 1) // Too many calls to PopItemFlag() - we always leave a 0 at the bottom of the stack.
	g.ItemFlagsStack = g.ItemFlagsStack[:len(g.ItemFlagsStack)-1]
	g.CurrentItemFlags = g.ItemFlagsStack[len(g.ItemFlagsStack)-1]
}

func BeginViewportSideBar(e string, t *ImGuiViewport, dir ImGuiDir, size float, window_flags ImGuiWindowFlags) bool {
	panic("not implemented")
}

// Menus
func BeginMenuEx(l string, n string, enabled bool /*= true*/) bool { panic("not implemented") }
func MenuItemEx(l string, n string, t string, selected bool, enabled bool /*= true*/) bool {
	panic("not implemented")
}

// t0 = previous time (e.g.: g.Time - g.IO.DeltaTime)
// t1 = current time (e.g.: g.Time)
// An event is triggered at:
//  t = 0.0f     t = repeat_delay,    t = repeat_delay + repeat_rate*N
func CalcTypematicRepeatAmount(t0, t1, repeat_delay, repeat_rate float) int {
	if t1 == 0.0 {
		return 1
	}
	if t0 >= t1 {
		return 0
	}
	if repeat_rate <= 0.0 {
		return bool2int((t0 < repeat_delay) && (t1 >= repeat_delay))
	}
	var count_t0 int = -1
	if t0 >= repeat_delay {
		count_t0 = (int)((t0 - repeat_delay) / repeat_rate)
	}
	var count_t1 int = -1
	if t1 >= repeat_delay {
		count_t1 = (int)((t1 - repeat_delay) / repeat_rate)
	}
	var count int = count_t1 - count_t0
	return count
}

func SetActiveIdUsingNavAndKeys() {
	var g = GImGui
	IM_ASSERT(g.ActiveId != 0)
	g.ActiveIdUsingNavDirMask = ^(ImU32)(0)
	g.ActiveIdUsingNavInputMask = ^(ImU32)(0)
	g.ActiveIdUsingKeyInputMask = ^(ImU64)(0)
	NavMoveRequestCancel()
}

func IsActiveIdUsingNavDir(dir ImGuiDir) bool {
	var g *ImGuiContext = GImGui
	return (g.ActiveIdUsingNavDirMask & (1 << dir)) != 0
}
func IsActiveIdUsingNavInput(input ImGuiNavInput) bool {
	var g *ImGuiContext = GImGui
	return (g.ActiveIdUsingNavInputMask & (1 << input)) != 0
}
func IsActiveIdUsingKey(key ImGuiKey) bool {
	var g *ImGuiContext = GImGui
	IM_ASSERT(key < 64)
	return (g.ActiveIdUsingKeyInputMask & ((ImU64)(1) << key)) != 0
}

func IsKeyPressedMap(key ImGuiKey, repeat bool /*= true*/) bool {
	var g *ImGuiContext = GImGui
	var key_index int = g.IO.KeyMap[key]
	if key_index >= 0 {
		return IsKeyPressed(key_index, repeat)
	}
	return false
}
func IsNavInputDown(n ImGuiNavInput) bool {
	var g *ImGuiContext = GImGui
	return g.IO.NavInputs[n] > 0.0
}
func IsNavInputTest(n ImGuiNavInput, rm ImGuiInputReadMode) bool {
	return (GetNavInputAmount(n, rm) > 0.0)
}

// Internal Columns API (this is not exposed because we will encourage transitioning to the Tables API)
func SetWindowClipRectBeforeSetChannel(w *ImGuiWindow, clip_rect *ImRect) {
	panic("not implemented")
}
func BeginColumns(d string, count int, flags ImGuiOldColumnFlags)         { panic("not implemented") } // setup number of columns. use an identifier to distinguish multiple column sets. close with EndColumns().
func EndColumns()                                                         { panic("not implemented") } // close columns
func PushColumnClipRect(column_index int)                                 { panic("not implemented") }
func PopColumnsBackground()                                               { panic("not implemented") }
func GetColumnsID(d string, count int) ImGuiID                            { panic("not implemented") }
func FindOrCreateColumns(w *ImGuiWindow, id ImGuiID) *ImGuiOldColumns     { panic("not implemented") }
func GetColumnOffsetFromNorm(s *ImGuiOldColumns, offset_norm float) float { panic("not implemented") }
func GetColumnNormFromOffset(s *ImGuiOldColumns, offset float) float      { panic("not implemented") }

// Tables: Candidates for public API
func TableOpenContextMenu(column_n int /*= -1*/)    { panic("not implemented") }
func TableSetColumnWidth(column_n int, width float) { panic("not implemented") }
func TableSetColumnSortDirection(column_n int, sort_direction ImGuiSortDirection, append_to_sort_specs bool) {
	panic("not implemented")
}
func TableGetHoveredColumn() int     { panic("not implemented") } // May use (TableGetColumnFlags() & ImGuiTableColumnFlags_IsHovered) instead. Return hovered column. return -1 when table is not hovered. return columns_count if the unused space at the right of visible columns is hovered.
func TableGetHeaderRowHeight() float { panic("not implemented") }
func TablePushBackgroundChannel()    { panic("not implemented") }
func TablePopBackgroundChannel()     { panic("not implemented") }

// Tables: Internals
func GetCurrentTable() *ImGuiTable         { var g *ImGuiContext = GImGui; return g.CurrentTable }
func TableFindByID(id ImGuiID) *ImGuiTable { panic("not implemented") }
func BeginTableEx(e string, id ImGuiID, columns_count int, flags ImGuiTableFlags, outer_size *ImVec2, inner_width float) bool {
	panic("not implemented")
}
func TableBeginInitMemory(e *ImGuiTable, columns_count int) { panic("not implemented") }
func TableBeginApplyRequests(e *ImGuiTable)                 { panic("not implemented") }
func TableSetupDrawChannels(e *ImGuiTable)                  { panic("not implemented") }
func TableUpdateLayout(e *ImGuiTable)                       { panic("not implemented") }
func TableUpdateBorders(e *ImGuiTable)                      { panic("not implemented") }
func TableUpdateColumnsWeightFromWidth(e *ImGuiTable)       { panic("not implemented") }
func TableDrawBorders(e *ImGuiTable)                        { panic("not implemented") }
func TableDrawContextMenu(e *ImGuiTable)                    { panic("not implemented") }
func TableMergeDrawChannels(e *ImGuiTable)                  { panic("not implemented") }
func TableSortSpecsSanitize(e *ImGuiTable)                  { panic("not implemented") }
func TableSortSpecsBuild(e *ImGuiTable)                     { panic("not implemented") }
func TableGetColumnNextSortDirection(n *ImGuiTableColumn) ImGuiSortDirection {
	panic("not implemented")
}
func TableFixColumnSortDirection(e *ImGuiTable, n *ImGuiTableColumn)   { panic("not implemented") }
func TableGetColumnWidthAuto(e *ImGuiTable, n *ImGuiTableColumn) float { panic("not implemented") }
func TableBeginRow(e *ImGuiTable)                                      { panic("not implemented") }
func TableEndRow(e *ImGuiTable)                                        { panic("not implemented") }
func TableBeginCell(e *ImGuiTable, column_n int)                       { panic("not implemented") }
func TableEndCell(e *ImGuiTable)                                       { panic("not implemented") }
func TableGetCellBgRect(e *ImGuiTable, column_n int) ImRect            { panic("not implemented") }
func tableGetColumnName(e *ImGuiTable, column_n int) string            { panic("not implemented") }
func TableGetColumnResizeID(e *ImGuiTable, column_n int, instance_no int) ImGuiID {
	panic("not implemented")
}
func TableGetMaxColumnWidth(e *ImGuiTable, column_n int) float     { panic("not implemented") }
func TableSetColumnWidthAutoSingle(e *ImGuiTable, column_n int)    { panic("not implemented") }
func TableSetColumnWidthAutoAll(e *ImGuiTable)                     { panic("not implemented") }
func TableRemove(e *ImGuiTable)                                    { panic("not implemented") }
func TableGcCompactTransientBuffers(e *ImGuiTable)                 { panic("not implemented") }
func TableGcCompactTransientBuffersTempData(e *ImGuiTableTempData) { panic("not implemented") }
func TableGcCompactSettings()                                      { panic("not implemented") }

// Tables: Settings
func TableLoadSettings(e *ImGuiTable)                                       { panic("not implemented") }
func TableSaveSettings(e *ImGuiTable)                                       { panic("not implemented") }
func TableResetSettings(e *ImGuiTable)                                      { panic("not implemented") }
func TableGetBoundSettings(e *ImGuiTable) *ImGuiTableSettings               { panic("not implemented") }
func TableSettingsInstallHandler(t *ImGuiContext)                           { panic("not implemented") }
func TableSettingsCreate(id ImGuiID, columns_count int) *ImGuiTableSettings { panic("not implemented") }
func TableSettingsFindByID(id ImGuiID) *ImGuiTableSettings                  { panic("not implemented") }

// Tab Bars
func BeginTabBarEx(r *ImGuiTabBar, bb *ImRect, flags ImGuiTabBarFlags) bool { panic("not implemented") }
func TabBarFindTabByID(r *ImGuiTabBar, tab_id ImGuiID) *ImGuiTabItem        { panic("not implemented") }
func TabBarRemoveTab(r *ImGuiTabBar, tab_id ImGuiID)                        { panic("not implemented") }
func TabBarCloseTab(r *ImGuiTabBar, b *ImGuiTabItem)                        { panic("not implemented") }
func TabBarQueueReorder(r *ImGuiTabBar, b *ImGuiTabItem, offset int)        { panic("not implemented") }
func TabBarQueueReorderFromMousePos(r *ImGuiTabBar, b *ImGuiTabItem, mouse_pos ImVec2) {
	panic("not implemented")
}
func TabBarProcessReorder(r *ImGuiTabBar) bool { panic("not implemented") }
func TabItemEx(r *ImGuiTabBar, l string, n *bool, flags ImGuiTabItemFlags) bool {
	panic("not implemented")
}
func TabItemCalcSize(l string, has_close_button bool) ImVec2 { panic("not implemented") }
func TabItemBackground(t *ImDrawList, bb *ImRect, flags ImGuiTabItemFlags, col ImU32) {
	panic("not implemented")
}
func TabItemLabelAndCloseButton(t *ImDrawList, bb *ImRect, flags ImGuiTabItemFlags, frame_padding ImVec2, l string, tab_id ImGuiID, close_button_id ImGuiID, is_contents_visible bool, out_just_closed *bool, out_text_clipped *bool) {
	panic("not implemented")
}

func RenderColorRectWithAlphaCheckerboard(t *ImDrawList, p_min ImVec2, p_max ImVec2, fill_col ImU32, grid_step float, grid_off ImVec2, rounding float, flags ImDrawFlags) {
	panic("not implemented")
}

// Navigation highlight
func RenderNavHighlight(bb *ImRect, id ImGuiID, flags ImGuiNavHighlightFlags) {
	var g = GImGui
	if id != g.NavId {
		return
	}
	if g.NavDisableHighlight && 0 == (flags&ImGuiNavHighlightFlags_AlwaysDraw) {
		return
	}
	var window = g.CurrentWindow
	if window.DC.NavHideHighlightOneFrame {
		return
	}

	var rounding float
	if 0 == (flags & ImGuiNavHighlightFlags_NoRounding) {
		rounding = g.Style.FrameRounding
	}

	var display_rect ImRect = *bb
	display_rect.ClipWith(window.ClipRect)
	if flags&ImGuiNavHighlightFlags_TypeDefault != 0 {
		var THICKNESS float = 2.0
		var DISTANCE float = 3.0 + THICKNESS*0.5
		display_rect.ExpandVec(ImVec2{DISTANCE, DISTANCE})
		var fully_visible bool = window.ClipRect.ContainsRect(display_rect)
		if !fully_visible {
			window.DrawList.PushClipRect(display_rect.Min, display_rect.Max, false)
		}
		window.DrawList.AddRect(display_rect.Min.Add(ImVec2{THICKNESS * 0.5, THICKNESS * 0.5}), display_rect.Max.Sub(ImVec2{THICKNESS * 0.5, THICKNESS * 0.5}), GetColorU32FromID(ImGuiCol_NavHighlight, 1), rounding, 0, THICKNESS)
		if !fully_visible {
			window.DrawList.PopClipRect()
		}
	}
	if flags&ImGuiNavHighlightFlags_TypeThin != 0 {
		window.DrawList.AddRect(display_rect.Min, display_rect.Max, GetColorU32FromID(ImGuiCol_NavHighlight, 1), rounding, 0, 1.0)
	}
}

func RenderBullet(t *ImDrawList, pos ImVec2, col ImU32) { panic("not implemented") }

func RenderMouseCursor(t *ImDrawList, pos ImVec2, scale float, mouse_cursor ImGuiMouseCursor, col_fill ImU32, col_border ImU32, col_shadow ImU32) {
	panic("not implemented")
}
func RenderArrowPointingAt(t *ImDrawList, pos ImVec2, half_sz ImVec2, direction ImGuiDir, col ImU32) {
	panic("not implemented")
}
func RenderRectFilledRangeH(t *ImDrawList, rect *ImRect, col ImU32, x_start_norm float, x_end_norm float, rounding float) {
	panic("not implemented")
}
func RenderRectFilledWithHole(t *ImDrawList, outer ImRect, inner ImRect, col ImU32, rounding float) {
	panic("not implemented")
}

// Widgets

func CloseButton(id ImGuiID, pos *ImVec2) bool { panic("not implemented") }
func ArrowButtonEx(d string, dir ImGuiDir, size_arg ImVec2, flags ImGuiButtonFlags) bool {
	panic("not implemented")
}

func ImageButtonEx(id ImGuiID, texture_id ImTextureID, size *ImVec2, uv0 *ImVec2, uv1 *ImVec2, padding *ImVec2, bg_col *ImVec4, tint_col *ImVec4) bool {
	panic("not implemented")
}

func CheckboxFlagsU(l string, s *ImS64, flags_value ImS64) bool { panic("not implemented") }
func CheckboxFlagsS(l string, s *ImU64, flags_value ImU64) bool { panic("not implemented") }

// Widgets low-level behaviors

func DragBehavior(id ImGuiID, data_type ImGuiDataType, v interface{}, v_speed float, n interface{}, x interface{}, t string, flags ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderBehavior(bb *ImRect, id ImGuiID, data_type ImGuiDataType, v interface{}, n interface{}, x interface{}, t string, flags ImGuiSliderFlags, b *ImRect) bool {
	panic("not implemented")
}
func SplitterBehavior(bb *ImRect, id ImGuiID, axis ImGuiAxis, size1 *float, size2 *float, min_size1 float, min_size2 float, hover_extend float, hover_visibility_delay float) bool {
	panic("not implemented")
}

// Template functions are instantiated in imgui_widgets.cpp for a finite number of types.
// To use them externally (for custom widget) you may need an "extern template" statement in your code in order to link to existing instances and silence Clang warnings (see #2036).
// e.g. " extern template func  RoundScalarWithFormatT<float, float>(t string, data_type ImGuiDataType, v float) float {panic("not implemented")} "
//template<typename T, typename SIGNED_T, typename FLOAT_T>   func  ScaleRatioFromValueT(data_type ImGuiDataType, T v, T v_min, T v_max, is_logarithmic bool, logarithmic_zero_epsilon float, zero_deadzone_size float) float {panic("not implemented")}
//template<typename T, typename SIGNED_T, typename FLOAT_T>   func  ScaleValueFromRatioT(data_type ImGuiDataType, t float, T v_min, T v_max, is_logarithmic bool, logarithmic_zero_epsilon float, zero_deadzone_size float) T {panic("not implemented")}
//template<typename T, typename SIGNED_T, typename FLOAT_T>   func  DragBehaviorT(data_type ImGuiDataType, v *T, v_speed float, T v_min, T v_max, t string, flags ImGuiSliderFlags) bool {panic("not implemented")}
//template<typename T, typename SIGNED_T, typename FLOAT_T>   func  SliderBehaviorT(bb *ImRect, id ImGuiID, data_type ImGuiDataType, v *T, T v_min, T v_max, t string, flags ImGuiSliderFlags, b *ImRect) bool {panic("not implemented")}
//template<typename T, typename SIGNED_T>                     func  RoundScalarWithFormatT(t string, data_type ImGuiDataType, T v) T {panic("not implemented")}
//template<typename T>                                        func  CheckboxFlagsT(l string, s *T, T flags_value) bool {panic("not implemented")}

// Data type helpers
func DataTypeGetInfo(data_type ImGuiDataType) *ImGuiDataTypeInfo { panic("not implemented") }
func DataTypeFormatString(f *char, buf_size int, data_type ImGuiDataType, a interface{}, t string) int {
	panic("not implemented")
}
func DataTypeApplyOp(data_type ImGuiDataType, op int, t interface{}, arg_1 interface{}, arg_2 interface{}) {
	panic("not implemented")
}
func DataTypeApplyOpFromText(buf string, initial_value_buf string, data_type ImGuiDataType, a interface{}, t string) bool {
	panic("not implemented")
}
func DataTypeCompare(data_type ImGuiDataType, arg_1 interface{}, arg_2 interface{}) int {
	panic("not implemented")
}
func DataTypeClamp(data_type ImGuiDataType, a interface{}, n interface{}, x interface{}) bool {
	panic("not implemented")
}

// InputText
func InputTextEx(l string, t string, f *char, buf_size int, size_arg *ImVec2, flags ImGuiInputTextFlags, callback ImGuiInputTextCallback, a interface{}) bool {
	panic("not implemented")
}
func TempInputText(bb *ImRect, id ImGuiID, l string, f *char, buf_size int, flags ImGuiInputTextFlags) bool {
	panic("not implemented")
}
func TempInputScalar(bb *ImRect, id ImGuiID, l string, data_type ImGuiDataType, a interface{}, t string, n interface{}, x interface{}) bool {
	panic("not implemented")
}
func TempInputIsActive(id ImGuiID) bool {
	var g *ImGuiContext = GImGui
	return (g.ActiveId == id && g.TempInputId == id)
}
func GetInputTextState(id ImGuiID) *ImGuiInputTextState {
	var g *ImGuiContext = GImGui
	if g.InputTextState.ID == id {
		return &g.InputTextState
	}
	return nil
} // Get input text state if active

// Color
func ColorTooltip(t string, l *float, flags ImGuiColorEditFlags)  { panic("not implemented") }
func ColorEditOptionsPopup(l *float, flags ImGuiColorEditFlags)   { panic("not implemented") }
func ColorPickerOptionsPopup(l *float, flags ImGuiColorEditFlags) { panic("not implemented") }

// Plot
func PlotEx(plot_type ImGuiPlotType, l string, values_getter func(data interface{}, idx int) float, a interface{}, values_count int, values_offset int, t string, scale_min float, scale_max float, frame_size ImVec2) int {
	panic("not implemented")
}

// Shade functions (write over already created vertices)
func ShadeVertsLinearColorGradientKeepAlpha(t *ImDrawList, vert_start_idx int, vert_end_idx int, gradient_p0 ImVec2, gradient_p1 ImVec2, col0 ImU32, col1 ImU32) {
	panic("not implemented")
}
func ShadeVertsLinearUV(t *ImDrawList, vert_start_idx int, vert_end_idx int, a *ImVec2, b *ImVec2, uv_a *ImVec2, uv_b *ImVec2, clamp bool) {
	panic("not implemented")
}

func ImFontAtlasBuildMultiplyCalcLookupTable(out_table []byte, in_multiply_factor float32) {
	panic("not implemented")
}
func ImFontAtlasBuildMultiplyRectAlpha8(table []byte, pixels []byte, x, y, w, h, stride int) {
	panic("not implemented")
}
