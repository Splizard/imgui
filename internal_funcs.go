package imgui

// Windows
// We should always have a CurrentWindow in the stack (there is an implicit "Debug" window)
// If this ever crash because g.CurrentWindow is NULL it means that either
// - ImGui::NewFrame() has never been called, which is illegal.
// - You are calling ImGui functions after ImGui::EndFrame()/ImGui::Render() and before the next ImGui::NewFrame(), which is also illegal.
func GetCurrentWindowRead() *ImGuiWindow {
	var g *ImGuiContext = GImGui
	return g.CurrentWindow
}
func GetCurrentWindow() *ImGuiWindow {
	var g *ImGuiContext = GImGui
	g.CurrentWindow.WriteAccessed = true
	return g.CurrentWindow
}
func FindWindowByID(id ImGuiID) *ImGuiWindow { panic("not implemented") }
func FindWindowByName(e string) *ImGuiWindow { panic("not implemented") }
func UpdateWindowParentAndRootLinks(window *ImGuiWindow, flags ImGuiWindowFlags, parent_window *ImGuiWindow) {
	panic("not implemented")
}
func CalcWindowNextAutoFitSize(w *ImGuiWindow) ImVec2     { panic("not implemented") }
func IsWindowChildOf(w *ImGuiWindow, t *ImGuiWindow) bool { panic("not implemented") }
func IsWindowAbove(e *ImGuiWindow, w *ImGuiWindow) bool   { panic("not implemented") }
func IsWindowNavFocusable(w *ImGuiWindow) bool            { panic("not implemented") }

func setWindowPos(w *ImGuiWindow, pos *ImVec2, cond ImGuiCond)   { panic("not implemented") }
func setWindowSize(w *ImGuiWindow, size *ImVec2, cond ImGuiCond) { panic("not implemented") }
func setWindowCollapsed(w *ImGuiWindow, collapsed bool, cond ImGuiCond) {
	panic("not implemented")
}
func SetWindowHitTestHole(w *ImGuiWindow, pos *ImVec2, size *ImVec2) { panic("not implemented") }

// Windows: Display Order and Focus Order
func FocusWindow(w *ImGuiWindow) { panic("not implemented") }
func FocusTopMostWindowUnderOne(under_this_window *ImGuiWindow, ignore_window *ImGuiWindow) {
	panic("not implemented")
}
func BringWindowToFocusFront(w *ImGuiWindow)   { panic("not implemented") }
func BringWindowToDisplayFront(w *ImGuiWindow) { panic("not implemented") }
func BringWindowToDisplayBack(w *ImGuiWindow)  { panic("not implemented") }

// Fonts, drawing
func SetCurrentFont(t *ImFont) { panic("not implemented") }
func GetDefaultFont() *ImFont {
	var g *ImGuiContext = GImGui
	if g.IO.FontDefault != nil {
		return g.IO.FontDefault
	}
	return g.IO.Fonts.Fonts[0]
}
func getForegroundDrawList(window *ImGuiWindow) *ImDrawList      { return GetForegroundDrawList() } // This seemingly unnecessary wrapper simplifies compatibility between the 'master' and 'docking' branches.
func getBackgroundDrawList(t *ImGuiViewport) *ImDrawList         { panic("not implemented") }       // get background draw list for the given viewport. this draw list will be the first rendering one. Useful to quickly draw shapes/text behind dear imgui contents.
func GetForegroundDrawListViewport(t *ImGuiViewport) *ImDrawList { panic("not implemented") }       // get foreground draw list for the given viewport. this draw list will be the last rendered one. Useful to quickly draw shapes/text over dear imgui contents.

// Init
func Initialize(context *ImGuiContext) {
	var g *ImGuiContext = context
	IM_ASSERT(!g.Initialized && !g.SettingsLoaded)

	// Add .ini handle for ImGuiWindow type
	{
		var ini_handler ImGuiSettingsHandler
		ini_handler.TypeName = "Window"
		ini_handler.TypeHash = ImHashStr("Window", size_t(len("Window")), 0)
		ini_handler.ClearAllFn = WindowSettingsHandler_ClearAll
		ini_handler.ReadOpenFn = WindowSettingsHandler_ReadOpen
		ini_handler.ReadLineFn = WindowSettingsHandler_ReadLine
		ini_handler.ApplyAllFn = WindowSettingsHandler_ApplyAll
		ini_handler.WriteAllFn = WindowSettingsHandler_WriteAll
		g.SettingsHandlers = append(g.SettingsHandlers, ini_handler)
	}

	// Add .ini handle for ImGuiTable type
	//TableSettingsInstallHandler(context)

	// Create default viewport
	var viewport ImGuiViewportP = NewImGuiViewportP()
	g.Viewports = append(g.Viewports, &viewport)

	g.Initialized = true
}

func Shutdown(t *ImGuiContext) { panic("not implemented") } // Since 1.60 this is a _private_ function. You can call DestroyContext() to destroy the context created by CreateContext().

// NewFrame
func UpdateHoveredWindowAndCaptureFlags()   { panic("not implemented") }
func StartMouseMovingWindow(w *ImGuiWindow) { panic("not implemented") }
func UpdateMouseMovingWindowNewFrame()      { panic("not implemented") }
func UpdateMouseMovingWindowEndFrame()      { panic("not implemented") }

// Generic context hooks
func AddContextHook(context *ImGuiContext, hook *ImGuiContextHook) ImGuiID { panic("not implemented") }
func RemoveContextHook(context *ImGuiContext, hook_to_remove ImGuiID)      { panic("not implemented") }
func CallContextHooks(context *ImGuiContext, htype ImGuiContextHookType) {
	panic("not implemented")
}

// Settings
func MarkIniSettingsDirty()                                    { panic("not implemented") }
func MarkIniSettingsDirtyWindow(w *ImGuiWindow)                { panic("not implemented") }
func ClearIniSettings()                                        { panic("not implemented") }
func CreateNewWindowSettings(e string) *ImGuiWindowSettings    { panic("not implemented") }
func FindWindowSettings(id ImGuiID) *ImGuiWindowSettings       { panic("not implemented") }
func FindOrCreateWindowSettings(e string) *ImGuiWindowSettings { panic("not implemented") }
func FindSettingsHandler(e string) *ImGuiSettingsHandler       { panic("not implemented") }

// Scrolling
func SetNextWindowScroll(scroll *ImVec2)        { panic("not implemented") } // Use -1.0f on one axis to leave as-is
func setScrollX(w *ImGuiWindow, scroll_x float) { panic("not implemented") }
func setScrollY(w *ImGuiWindow, scroll_y float) { panic("not implemented") }
func setScrollFromPosX(w *ImGuiWindow, local_x float, center_x_ratio float) {
	panic("not implemented")
}
func setScrollFromPosY(w *ImGuiWindow, local_y float, center_y_ratio float) {
	panic("not implemented")
}
func ScrollToBringRectIntoView(w *ImGuiWindow, item_rect *ImRect) ImVec2 { panic("not implemented") }

// Basic Accessors
func GetItemID() ImGuiID { var g *ImGuiContext = GImGui; return g.LastItemData.ID } // Get ID of last item (~~ often same ImGui::GetID(label) beforehand)
func GetItemStatusFlags() ImGuiItemStatusFlags {
	var g *ImGuiContext = GImGui
	return g.LastItemData.StatusFlags
}
func GetItemFlags() ImGuiItemFlags                           { var g *ImGuiContext = GImGui; return g.LastItemData.InFlags }
func GetActiveID() ImGuiID                                   { var g *ImGuiContext = GImGui; return g.ActiveId }
func GetFocusID() ImGuiID                                    { var g *ImGuiContext = GImGui; return g.NavId }
func SetActiveID(id ImGuiID, w *ImGuiWindow)                 { panic("not implemented") }
func SetFocusID(id ImGuiID, w *ImGuiWindow)                  { panic("not implemented") }
func ClearActiveID()                                         { panic("not implemented") }
func GetHoveredID() ImGuiID                                  { panic("not implemented") }
func SetHoveredID(id ImGuiID)                                { panic("not implemented") }
func KeepAliveID(id ImGuiID)                                 { panic("not implemented") }
func MarkItemEdited(id ImGuiID)                              { panic("not implemented") } // Mark data associated to given item as "edited", used by IsItemDeactivatedAfterEdit() function.
func PushOverrideID(id ImGuiID)                              { panic("not implemented") } // Push given value as-is at the top of the ID stack (whereas PushID combines old and new hashes)
func GetIDWithSeed(n string, d string, seed ImGuiID) ImGuiID { panic("not implemented") }

// Basic Helpers for widget code
func ItemSizeVec(size *ImVec2, text_baseline_y float /*= -1.0f*/) { panic("not implemented") }
func ItemSizeRect(bb *ImRect, text_baseline_y float /*= -1.0f*/)  { panic("not implemented") }
func ItemAdd(bb *ImRect, id ImGuiID, b *ImRect, extra_flags ImGuiItemFlags) bool {
	panic("not implemented")
}
func ItemHoverable(bb *ImRect, id ImGuiID) bool                           { panic("not implemented") }
func ItemInputable(w *ImGuiWindow, id ImGuiID)                            { panic("not implemented") }
func IsClippedEx(bb *ImRect, id ImGuiID, clip_even_when_logged bool) bool { panic("not implemented") }
func CalcItemSize(size ImVec2, default_w float, default_h float) ImVec2   { panic("not implemented") }
func CalcWrapWidthForPos(pos *ImVec2, wrap_pos_x float) float             { panic("not implemented") }
func PushMultiItemsWidths(components int, width_full float)               { panic("not implemented") }
func IsItemToggledSelection() bool                                        { panic("not implemented") } // Was the last item selection toggled? (after Selectable(), TreeNode() etc. We only returns toggle _event_ in order to handle clipping correctly)
func GetContentRegionMaxAbs() ImVec2                                      { panic("not implemented") }
func ShrinkWidths(s *ImGuiShrinkWidthItem, count int, width_excess float) {
	panic("not implemented")
}

// Parameter stacks
func PushItemFlag(option ImGuiItemFlags, enabled bool) { panic("not implemented") }
func PopItemFlag()                                     { panic("not implemented") }

// Logging/Capture
func LogBegin(ltype ImGuiLogType, auto_open_depth int)      { panic("not implemented") } // . BeginCapture() when we design v2 api, for now stay under the radar by using the old name.
func LogToBuffer(auto_open_depth int /*= -1*/)              { panic("not implemented") } // Start logging/capturing to internal buffer
func LogRenderedText(s *ImVec2, t string, d string)         { panic("not implemented") }
func LogSetNextTextDecoration(prefix string, suffix string) { panic("not implemented") }

// Popups, Modals, Tooltips
func BeginChildEx(e string, id ImGuiID, size_arg *ImVec2, border bool, flags ImGuiWindowFlags) bool {
	panic("not implemented")
}
func OpenPopupEx(id ImGuiID, popup_flags ImGuiPopupFlags) { panic("not implemented") }
func ClosePopupToLevel(remaining int, restore_focus_to_window_under_popup bool) {
	panic("not implemented")
}
func ClosePopupsOverWindow(w *ImGuiWindow, restore_focus_to_window_under_popup bool) {
	panic("not implemented")
}
func isPopupOpen(id ImGuiID, popup_flags ImGuiPopupFlags) bool   { panic("not implemented") }
func BeginPopupEx(id ImGuiID, extra_flags ImGuiWindowFlags) bool { panic("not implemented") }
func BeginTooltipEx(extra_flags ImGuiWindowFlags, tooltip_flags ImGuiTooltipFlags) {
	panic("not implemented")
}
func GetPopupAllowedExtentRect(w *ImGuiWindow) ImRect { panic("not implemented") }
func GetTopMostPopupModal() *ImGuiWindow              { panic("not implemented") }
func FindBestWindowPosForPopup(w *ImGuiWindow) ImVec2 { panic("not implemented") }
func FindBestWindowPosForPopupEx(ref_pos *ImVec2, size *ImVec2, r *ImGuiDir, r_outer *ImRect, r_avoid *ImRect, policy ImGuiPopupPositionPolicy) ImVec2 {
	panic("not implemented")
}
func BeginViewportSideBar(e string, t *ImGuiViewport, dir ImGuiDir, size float, window_flags ImGuiWindowFlags) bool {
	panic("not implemented")
}

// Menus
func BeginMenuEx(l string, n string, enabled bool /*= true*/) bool { panic("not implemented") }
func MenuItemEx(l string, n string, t string, selected bool, enabled bool /*= true*/) bool {
	panic("not implemented")
}

// Combos
func BeginComboPopup(popup_id ImGuiID, bb *ImRect, flags ImGuiComboFlags) bool {
	panic("not implemented")
}
func BeginComboPreview() bool { panic("not implemented") }
func EndComboPreview()        { panic("not implemented") }

// Gamepad/Keyboard Navigation
func NavInitWindow(w *ImGuiWindow, force_reinit bool) { panic("not implemented") }
func NavMoveRequestButNoResultYet() bool              { panic("not implemented") }
func NavMoveRequestSubmit(move_dir ImGuiDir, clip_dir ImGuiDir, move_flags ImGuiNavMoveFlags) {
	panic("not implemented")
}
func NavMoveRequestForward(move_dir ImGuiDir, clip_dir ImGuiDir, move_flags ImGuiNavMoveFlags) {
	panic("not implemented")
}
func NavMoveRequestCancel()      { panic("not implemented") }
func NavMoveRequestApplyResult() { panic("not implemented") }
func NavMoveRequestTryWrapping(w *ImGuiWindow, move_flags ImGuiNavMoveFlags) {
	panic("not implemented")
}
func GetNavInputAmount(n ImGuiNavInput, mode ImGuiInputReadMode) float { panic("not implemented") }
func GetNavInputAmount2d(dir_sources ImGuiNavDirSourceFlags, mode ImGuiInputReadMode, slow_factor float, fast_factor float) ImVec2 {
	panic("not implemented")
}
func CalcTypematicRepeatAmount(t0 float, t1 float, repeat_delay float, repeat_rate float) int {
	panic("not implemented")
}
func ActivateItem(id ImGuiID) { panic("not implemented") } // Remotely activate a button, checkbox, tree node etc. given its unique ID. activation is queued and processed on the next frame when the item is encountered again.
func SetNavID(id ImGuiID, nav_layer ImGuiNavLayer, focus_scope_id ImGuiID, rect_rel *ImRect) {
	panic("not implemented")
}

// Focus Scope (WIP)
// This is generally used to identify a selection set (multiple of which may be in the same window), as selection
// patterns generally need to react (e.g. clear selection) when landing on an item of the set.
func PushFocusScope(id ImGuiID)     { panic("not implemented") }
func PopFocusScope()                { panic("not implemented") }
func GetFocusedFocusScope() ImGuiID { var g *ImGuiContext = GImGui; return g.NavFocusScopeId } // Focus scope which is actually active
func GetFocusScope() ImGuiID {
	var g *ImGuiContext = GImGui
	return g.CurrentWindow.DC.NavFocusScopeIdCurrent
} // Focus scope we are outputting into, set by PushFocusScope()

// Inputs
// FIXME: Eventually we should aim to move e.g. IsActiveIdUsingKey() into IsKeyXXX functions.
func SetItemUsingMouseWheel()     { panic("not implemented") }
func SetActiveIdUsingNavAndKeys() { panic("not implemented") }
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
func IsMouseDragPastThreshold(button ImGuiMouseButton, lock_threshold float /*= -1.0f*/) bool {
	panic("not implemented")
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
func GetMergedKeyModFlags() ImGuiKeyModFlags { panic("not implemented") }

// Drag and Drop
func BeginDragDropTargetCustom(bb *ImRect, id ImGuiID) bool { panic("not implemented") }
func ClearDragDrop()                                        { panic("not implemented") }
func IsDragDropPayloadBeingAccepted() bool                  { panic("not implemented") }

// Internal Columns API (this is not exposed because we will encourage transitioning to the Tables API)
func SetWindowClipRectBeforeSetChannel(w *ImGuiWindow, clip_rect *ImRect) {
	panic("not implemented")
}
func BeginColumns(d string, count int, flags ImGuiOldColumnFlags)         { panic("not implemented") } // setup number of columns. use an identifier to distinguish multiple column sets. close with EndColumns().
func EndColumns()                                                         { panic("not implemented") } // close columns
func PushColumnClipRect(column_index int)                                 { panic("not implemented") }
func PushColumnsBackground()                                              { panic("not implemented") }
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

// Render helpers
// AVOID USING OUTSIDE OF IMGUI.CPP! NOT FOR PUBLIC CONSUMPTION. THOSE FUNCTIONS ARE A MESS. THEIR SIGNATURE AND BEHAVIOR WILL CHANGE, THEY NEED TO BE REFACTORED INTO SOMETHING DECENT.
// NB: All position are in absolute pixels coordinates (we are never using window coordinates internally)
func RenderText(pos ImVec2, t string, d string, hide_text_after_hash bool /*= true*/) {
	panic("not implemented")
}
func RenderTextWrapped(pos ImVec2, t string, d string, wrap_width float) {
	panic("not implemented")
}
func RenderTextClipped(pos_min *ImVec2, pos_max *ImVec2, text string, text_end string, n *ImVec2, align *ImVec2, t *ImRect) {
	panic("not implemented")
}
func RenderTextClippedEx(t *ImDrawList, pos_min *ImVec2, pos_max *ImVec2, ttext string, text_end string, n *ImVec2, align *ImVec2, clip_rect *ImRect) {
	panic("not implemented")
}
func RenderTextEllipsis(t *ImDrawList, pos_min *ImVec2, pos_max *ImVec2, clip_max_x float, ellipsis_max_x float, text string, text_end string, n *ImVec2) {
	panic("not implemented")
}
func RenderFrame(p_min ImVec2, p_max ImVec2, fill_col ImU32, border bool /*= true*/, rounding float) {
	panic("not implemented")
}
func RenderFrameBorder(p_min ImVec2, p_max ImVec2, rounding float) { panic("not implemented") }
func RenderColorRectWithAlphaCheckerboard(t *ImDrawList, p_min ImVec2, p_max ImVec2, fill_col ImU32, grid_step float, grid_off ImVec2, rounding float, flags ImDrawFlags) {
	panic("not implemented")
}
func RenderNavHighlight(bb *ImRect, id ImGuiID, flags ImGuiNavHighlightFlags) {
	panic("not implemented")
}                                                   // Navigation highlight
func FindRenderedTextEnd(t string, d string) string { panic("not implemented") } // Find the optional ## from which we stop displaying text.

// Render helpers (those functions don't access any ImGui state!)
func RenderArrow(t *ImDrawList, pos ImVec2, col ImU32, dir ImGuiDir, scale float /*= 1.0f*/) {
	panic("not implemented")
}
func RenderBullet(t *ImDrawList, pos ImVec2, col ImU32)              { panic("not implemented") }
func RenderCheckMark(t *ImDrawList, pos ImVec2, col ImU32, sz float) { panic("not implemented") }
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
func TextEx(t string, d string, flags ImGuiTextFlags)                  { panic("not implemented") }
func ButtonEx(l string, size_arg *ImVec2, flags ImGuiButtonFlags) bool { panic("not implemented") }
func CloseButton(id ImGuiID, pos *ImVec2) bool                         { panic("not implemented") }
func CollapseButton(id ImGuiID, pos *ImVec2) bool                      { panic("not implemented") }
func ArrowButtonEx(d string, dir ImGuiDir, size_arg ImVec2, flags ImGuiButtonFlags) bool {
	panic("not implemented")
}
func Scrollbar(axis ImGuiAxis) { panic("not implemented") }
func ScrollbarEx(bb *ImRect, id ImGuiID, axis ImGuiAxis, v *float, avail_v float, contents_v float, flags ImDrawFlags) bool {
	panic("not implemented")
}
func ImageButtonEx(id ImGuiID, texture_id ImTextureID, size *ImVec2, uv0 *ImVec2, uv1 *ImVec2, padding *ImVec2, bg_col *ImVec4, tint_col *ImVec4) bool {
	panic("not implemented")
}
func GetWindowScrollbarRect(w *ImGuiWindow, axis ImGuiAxis) ImRect { panic("not implemented") }
func GetWindowScrollbarID(w *ImGuiWindow, axis ImGuiAxis) ImGuiID  { panic("not implemented") }
func GetWindowResizeCornerID(w *ImGuiWindow, n int) ImGuiID        { panic("not implemented") } // 0..3: corners
func GetWindowResizeBorderID(w *ImGuiWindow, dir ImGuiDir) ImGuiID { panic("not implemented") }
func SeparatorEx(flags ImGuiSeparatorFlags)                        { panic("not implemented") }
func CheckboxFlagsU(l string, s *ImS64, flags_value ImS64) bool    { panic("not implemented") }
func CheckboxFlagsS(l string, s *ImU64, flags_value ImU64) bool    { panic("not implemented") }

// Widgets low-level behaviors
func ButtonBehavior(bb *ImRect, id ImGuiID, out_hovered *bool, out_held *bool, flags ImGuiButtonFlags) bool {
	panic("not implemented")
}
func DragBehavior(id ImGuiID, data_type ImGuiDataType, v interface{}, v_speed float, n interface{}, x interface{}, t string, flags ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderBehavior(bb *ImRect, id ImGuiID, data_type ImGuiDataType, v interface{}, n interface{}, x interface{}, t string, flags ImGuiSliderFlags, b *ImRect) bool {
	panic("not implemented")
}
func SplitterBehavior(bb *ImRect, id ImGuiID, axis ImGuiAxis, size1 *float, size2 *float, min_size1 float, min_size2 float, hover_extend float, hover_visibility_delay float) bool {
	panic("not implemented")
}
func TreeNodeBehavior(id ImGuiID, flags ImGuiTreeNodeFlags, l string, d string) bool {
	panic("not implemented")
}
func TreeNodeBehaviorIsOpen(id ImGuiID, flags ImGuiTreeNodeFlags) bool { panic("not implemented") } // Consume previous SetNextItemOpen() data, if any. May return true when logging
func TreePushOverrideID(id ImGuiID)                                    { panic("not implemented") }

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

// Garbage collection
func GcCompactTransientMiscBuffers()                 { panic("not implemented") }
func GcCompactTransientWindowBuffers(w *ImGuiWindow) { panic("not implemented") }
func GcAwakeTransientWindowBuffers(w *ImGuiWindow)   { panic("not implemented") }

// Debug Tools
func ErrorCheckEndFrameRecover(log_callback ImGuiErrorLogCallback, a interface{}) {
	panic("not implemented")
}
func DebugDrawItemRect(col ImU32 /*= IM_COL32(255,0,0,255)*/) {
	var g *ImGuiContext = GImGui
	var window *ImGuiWindow = g.CurrentWindow
	getForegroundDrawList(window).AddRect(g.LastItemData.Rect.Min, g.LastItemData.Rect.Max, col, 0, 0, 1)
}
func DebugStartItemPicker() { var g *ImGuiContext = GImGui; g.DebugItemPickerActive = true }

func ShowFontAtlas(s *ImFontAtlas)                              { panic("not implemented") }
func DebugNodeColumns(s *ImGuiOldColumns)                       { panic("not implemented") }
func DebugNodeDrawList(w *ImGuiWindow, t *ImDrawList, l string) { panic("not implemented") }
func DebugNodeDrawCmdShowMeshAndBoundingBox(out_draw_list *ImDrawList, draw_list *ImDrawList, d *ImDrawCmd, show_mesh bool, show_aabb bool) {
	panic("not implemented")
}
func DebugNodeFont(t *ImFont)                               { panic("not implemented") }
func DebugNodeStorage(e *ImGuiStorage, l string)            { panic("not implemented") }
func DebugNodeTabBar(r *ImGuiTabBar, l string)              { panic("not implemented") }
func DebugNodeTable(e *ImGuiTable)                          { panic("not implemented") }
func DebugNodeTableSettings(s *ImGuiTableSettings)          { panic("not implemented") }
func DebugNodeWindow(w *ImGuiWindow, l string)              { panic("not implemented") }
func DebugNodeWindowSettings(s *ImGuiWindowSettings)        { panic("not implemented") }
func DebugNodeWindowsList(windows []*ImGuiWindow, l string) { panic("not implemented") }
func DebugNodeViewport(t *ImGuiViewportP)                   { panic("not implemented") }
func DebugRenderViewportThumbnail(draw_list *ImDrawList, viewport *ImGuiViewportP, bb *ImRect) {
	panic("not implemented")
}

func ImFontAtlasBuildSetupFont(atlas *ImFontAtlas, font *ImFont, font_config *ImFontConfig, ascent float32, descent float32) {
	panic("not implemented")
}
func ImFontAtlasBuildPackCustomRects(atlas *ImFontAtlas, stbrp_context_opaque interface{}) {
	panic("not implemented")
}
func ImFontAtlasBuildFinish(atlas *ImFontAtlas) { panic("not implemented") }
func ImFontAtlasBuildRender8bppRectFromString(atlas *ImFontAtlas, x, y, w, h int, in_str string, in_marker_char byte, in_marker_pixel_value byte) {
	panic("not implemented")
}
func ImFontAtlasBuildRender32bppRectFromString(atlas *ImFontAtlas, x, y, w, h int, in_str string, in_marker_char byte, in_marker_pixel_value uint32) {
	panic("not implemented")
}
func ImFontAtlasBuildMultiplyCalcLookupTable(out_table []byte, in_multiply_factor float32) {
	panic("not implemented")
}
func ImFontAtlasBuildMultiplyRectAlpha8(table []byte, pixels []byte, x, y, w, h, stride int) {
	panic("not implemented")
}
