package imgui

// Fonts, drawing

func getForegroundDrawList(window *ImGuiWindow) *ImDrawList { return GetForegroundDrawList(nil) } // This seemingly unnecessary wrapper simplifies compatibility between the 'master' and 'docking' branches.
func getBackgroundDrawList(t *ImGuiViewport) *ImDrawList    { panic("not implemented") }          // get background draw list for the given viewport. this draw list will be the first rendering one. Useful to quickly draw shapes/text behind dear imgui contents.

func getForegroundDrawListViewport(viewport *ImGuiViewport, drawlist_no int, drawlist_name string) *ImDrawList {
	// Create the draw list on demand, because they are not frequently used for all viewports
	var g = *GImGui
	IM_ASSERT(drawlist_no < int(len(viewport.DrawLists)))
	var draw_list = viewport.DrawLists[drawlist_no]
	if draw_list == nil {
		l := NewImDrawList(&g.DrawListSharedData)
		draw_list = &l
		draw_list._OwnerName = drawlist_name
		viewport.DrawLists[drawlist_no] = draw_list
	}

	// Our ImDrawList system requires that there is always a command
	if viewport.DrawListsLastFrame[drawlist_no] != g.FrameCount {
		draw_list._ResetForNewFrame()
		draw_list.PushTextureID(g.IO.Fonts.TexID)
		draw_list.PushClipRect(viewport.Pos, viewport.Pos.Add(viewport.Size), false)
		viewport.DrawListsLastFrame[drawlist_no] = g.FrameCount
	}
	return draw_list
}

// get foreground draw list for the given viewport. this draw list will be the last rendered one. Useful to quickly draw shapes/text over dear imgui contents.
func GetForegroundDrawListViewport(viewport *ImGuiViewport) *ImDrawList {
	return getForegroundDrawListViewport(viewport, 1, "###Foreground")
}

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

// Helper for ColorPicker4()
// NB: This is rather brittle and will show artifact when rounding this enabled if rounded corners overlap multiple cells. Caller currently responsible for avoiding that.
// Spent a non reasonable amount of time trying to getting this right for ColorButton with rounding+anti-aliasing+ImGuiColorEditFlags_HalfAlphaPreview flag + various grid sizes and offsets, and eventually gave up... probably more reasonable to disable rounding altogether.
// FIXME: uses ImGui::GetColorU32
func RenderColorRectWithAlphaCheckerboard(draw_list *ImDrawList, p_min ImVec2, p_max ImVec2, col ImU32, grid_step float, grid_off ImVec2, rounding float, flags ImDrawFlags) {
	if (flags & ImDrawFlags_RoundCornersMask_) == 0 {
		flags = ImDrawFlags_RoundCornersDefault_
	}
	if ((col & IM_COL32_A_MASK) >> IM_COL32_A_SHIFT) < 0xFF {
		var col_bg1 ImU32 = GetColorU32FromInt(ImAlphaBlendColors(IM_COL32(204, 204, 204, 255), col))
		var col_bg2 ImU32 = GetColorU32FromInt(ImAlphaBlendColors(IM_COL32(128, 128, 128, 255), col))
		draw_list.AddRectFilled(p_min, p_max, col_bg1, rounding, flags)

		var yi int = 0
		for y := p_min.y + grid_off.y; y < p_max.y; y, yi = y+grid_step, yi+1 {
			var y1, y2 = ImClamp(y, p_min.y, p_max.y), ImMin(y+grid_step, p_max.y)
			if y2 <= y1 {
				continue
			}
			for x := p_min.x + grid_off.x + float(yi&1)*grid_step; x < p_max.x; x += grid_step * 2.0 {
				var x1, x2 = ImClamp(x, p_min.x, p_max.x), ImMin(x+grid_step, p_max.x)
				if x2 <= x1 {
					continue
				}
				var cell_flags ImDrawFlags = ImDrawFlags_RoundCornersNone
				if y1 <= p_min.y {
					if x1 <= p_min.x {
						cell_flags |= ImDrawFlags_RoundCornersTopLeft
					}
					if x2 >= p_max.x {
						cell_flags |= ImDrawFlags_RoundCornersTopRight
					}
				}
				if y2 >= p_max.y {
					if x1 <= p_min.x {
						cell_flags |= ImDrawFlags_RoundCornersBottomLeft
					}
					if x2 >= p_max.x {
						cell_flags |= ImDrawFlags_RoundCornersBottomRight
					}
				}

				// Combine flags

				if flags == ImDrawFlags_RoundCornersNone || cell_flags == ImDrawFlags_RoundCornersNone {
					cell_flags = ImDrawFlags_RoundCornersNone
				} else {
					cell_flags = (cell_flags & flags)
				}

				draw_list.AddRectFilled(ImVec2{x1, y1}, ImVec2{x2, y2}, col_bg2, rounding, cell_flags)
			}
		}
	} else {
		draw_list.AddRectFilled(p_min, p_max, col, rounding, flags)
	}
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

func RenderMouseCursor(draw_list *ImDrawList, pos ImVec2, scale float, mouse_cursor ImGuiMouseCursor, col_fill ImU32, col_border ImU32, col_shadow ImU32) {
	if mouse_cursor == ImGuiMouseCursor_None {
		return
	}
	IM_ASSERT(mouse_cursor > ImGuiMouseCursor_None && mouse_cursor < ImGuiMouseCursor_COUNT)

	var font_atlas *ImFontAtlas = draw_list._Data.Font.ContainerAtlas
	var offset, size ImVec2
	var uv1, uv2 [2]ImVec2
	if font_atlas.GetMouseCursorTexData(mouse_cursor, &offset, &size, &uv1, &uv2) {
		pos = pos.Sub(offset)
		var tex_id ImTextureID = font_atlas.TexID
		draw_list.PushTextureID(tex_id)
		draw_list.AddImage(tex_id, pos.Add(ImVec2{1, 0}.Scale(scale)), pos.Add((ImVec2{1, 0}.Add(size)).Scale(scale)), &uv2[0], &uv2[1], col_shadow)
		draw_list.AddImage(tex_id, pos.Add(ImVec2{2, 0}.Scale(scale)), pos.Add((ImVec2{2, 0}.Add(size)).Scale(scale)), &uv2[0], &uv2[1], col_shadow)
		draw_list.AddImage(tex_id, pos, pos.Add(size.Scale(scale)), &uv2[0], &uv2[1], col_border)
		draw_list.AddImage(tex_id, pos, pos.Add(size.Scale(scale)), &uv1[0], &uv1[1], col_fill)
		draw_list.PopTextureID()
	}
}

// Render an arrow. 'pos' is position of the arrow tip. half_sz.x is length from base to tip. half_sz.y is length on each side.
func RenderArrowPointingAt(draw_list *ImDrawList, pos ImVec2, half_sz ImVec2, direction ImGuiDir, col ImU32) {
	switch direction {
	case ImGuiDir_Left:
		draw_list.AddTriangleFilled(&ImVec2{pos.x + half_sz.x, pos.y - half_sz.y}, &ImVec2{pos.x + half_sz.x, pos.y + half_sz.y}, pos, col)
		return
	case ImGuiDir_Right:
		draw_list.AddTriangleFilled(&ImVec2{pos.x - half_sz.x, pos.y + half_sz.y}, &ImVec2{pos.x - half_sz.x, pos.y - half_sz.y}, pos, col)
		return
	case ImGuiDir_Up:
		draw_list.AddTriangleFilled(&ImVec2{pos.x + half_sz.x, pos.y + half_sz.y}, &ImVec2{pos.x - half_sz.x, pos.y + half_sz.y}, pos, col)
		return
	case ImGuiDir_Down:
		draw_list.AddTriangleFilled(&ImVec2{pos.x - half_sz.x, pos.y - half_sz.y}, &ImVec2{pos.x + half_sz.x, pos.y - half_sz.y}, pos, col)
		return
	case ImGuiDir_None:
	case ImGuiDir_COUNT:
		break // Fix warnings
	}
}

// FIXME: Cleanup and move code to ImDrawList.
func RenderRectFilledRangeH(draw_list *ImDrawList, rect *ImRect, col ImU32, x_start_norm float, x_end_norm float, rounding float) {
	if x_end_norm == x_start_norm {
		return
	}
	if x_start_norm > x_end_norm {
		ImSwap(x_start_norm, x_end_norm)
	}

	var p0 ImVec2 = ImVec2{ImLerp(rect.Min.x, rect.Max.x, x_start_norm), rect.Min.y}
	var p1 ImVec2 = ImVec2{ImLerp(rect.Min.x, rect.Max.x, x_end_norm), rect.Max.y}
	if rounding == 0.0 {
		draw_list.AddRectFilled(p0, p1, col, 0.0, 0)
		return
	}

	rounding = ImClamp(ImMin((rect.Max.x-rect.Min.x)*0.5, (rect.Max.y-rect.Min.y)*0.5)-1.0, 0.0, rounding)
	var inv_rounding = 1.0 / rounding
	var arc0_b float = ImAcos01(1.0 - (p0.x-rect.Min.x)*inv_rounding)
	var arc0_e float = ImAcos01(1.0 - (p1.x-rect.Min.x)*inv_rounding)
	var half_pi float = IM_PI * 0.5 // We will == compare to this because we know this is the exact value ImAcos01 can return.
	var x0 = ImMax(p0.x, rect.Min.x+rounding)
	if arc0_b == arc0_e {
		draw_list.PathLineTo(ImVec2{x0, p1.y})
		draw_list.PathLineTo(ImVec2{x0, p0.y})
	} else if arc0_b == 0.0 && arc0_e == half_pi {
		draw_list.PathArcToFast(ImVec2{x0, p1.y - rounding}, rounding, 3, 6) // BL
		draw_list.PathArcToFast(ImVec2{x0, p0.y + rounding}, rounding, 6, 9) // TR
	} else {
		draw_list.PathArcTo(ImVec2{x0, p1.y - rounding}, rounding, IM_PI-arc0_e, IM_PI-arc0_b, 3) // BL
		draw_list.PathArcTo(ImVec2{x0, p0.y + rounding}, rounding, IM_PI+arc0_b, IM_PI+arc0_e, 3) // TR
	}
	if p1.x > rect.Min.x+rounding {
		var arc1_b float = ImAcos01(1.0 - (rect.Max.x-p1.x)*inv_rounding)
		var arc1_e float = ImAcos01(1.0 - (rect.Max.x-p0.x)*inv_rounding)
		var x1 float = ImMin(p1.x, rect.Max.x-rounding)
		if arc1_b == arc1_e {
			draw_list.PathLineTo(ImVec2{x1, p0.y})
			draw_list.PathLineTo(ImVec2{x1, p1.y})
		} else if arc1_b == 0.0 && arc1_e == half_pi {
			draw_list.PathArcToFast(ImVec2{x1, p0.y + rounding}, rounding, 9, 12) // TR
			draw_list.PathArcToFast(ImVec2{x1, p1.y - rounding}, rounding, 0, 3)  // BR
		} else {
			draw_list.PathArcTo(ImVec2{x1, p0.y + rounding}, rounding, -arc1_e, -arc1_b, 3) // TR
			draw_list.PathArcTo(ImVec2{x1, p1.y - rounding}, rounding, +arc1_b, +arc1_e, 3) // BR
		}
	}
	draw_list.PathFillConvex(col)
}

func RenderRectFilledWithHole(draw_list *ImDrawList, outer ImRect, inner ImRect, col ImU32, rounding float) {
	var fill_L = (inner.Min.x > outer.Min.x)
	var fill_R = (inner.Max.x < outer.Max.x)
	var fill_U = (inner.Min.y > outer.Min.y)
	var fill_D = (inner.Max.y < outer.Max.y)
	if fill_L {
		var flags ImDrawFlags
		if !fill_U {
			flags |= ImDrawFlags_RoundCornersTopLeft
		}
		if fill_D {
			flags |= ImDrawFlags_RoundCornersBottomLeft
		}
		draw_list.AddRectFilled(ImVec2{outer.Min.x, inner.Min.y}, ImVec2{inner.Min.x, inner.Max.y}, col, rounding, flags)
	}
	if fill_R {
		var flags ImDrawFlags
		if !fill_U {
			flags |= ImDrawFlags_RoundCornersTopRight
		}
		if fill_D {
			flags |= ImDrawFlags_RoundCornersBottomRight
		}
		draw_list.AddRectFilled(ImVec2{inner.Max.x, inner.Min.y}, ImVec2{outer.Max.x, inner.Max.y}, col, rounding, flags)
	}
	if fill_U {
		var flags ImDrawFlags
		if !fill_L {
			flags |= ImDrawFlags_RoundCornersTopLeft
		}
		if fill_R {
			flags |= ImDrawFlags_RoundCornersTopRight
		}
		draw_list.AddRectFilled(ImVec2{inner.Min.x, outer.Min.y}, ImVec2{inner.Max.x, inner.Min.y}, col, rounding, flags)
	}
	if fill_D {
		var flags ImDrawFlags
		if !fill_L {
			flags |= ImDrawFlags_RoundCornersBottomLeft
		}
		if fill_R {
			flags |= ImDrawFlags_RoundCornersBottomRight
		}
		draw_list.AddRectFilled(ImVec2{inner.Min.x, inner.Max.y}, ImVec2{inner.Max.x, outer.Max.y}, col, rounding, flags)
	}
	if fill_L && fill_U {
		draw_list.AddRectFilled(ImVec2{outer.Min.x, outer.Min.y}, ImVec2{inner.Min.x, inner.Min.y}, col, rounding, ImDrawFlags_RoundCornersTopLeft)
	}
	if fill_R && fill_U {
		draw_list.AddRectFilled(ImVec2{inner.Max.x, outer.Min.y}, ImVec2{outer.Max.x, inner.Min.y}, col, rounding, ImDrawFlags_RoundCornersTopRight)
	}
	if fill_L && fill_D {
		draw_list.AddRectFilled(ImVec2{outer.Min.x, inner.Max.y}, ImVec2{inner.Min.x, outer.Max.y}, col, rounding, ImDrawFlags_RoundCornersBottomLeft)
	}
	if fill_R && fill_D {
		draw_list.AddRectFilled(ImVec2{inner.Max.x, inner.Max.y}, ImVec2{outer.Max.x, outer.Max.y}, col, rounding, ImDrawFlags_RoundCornersBottomRight)
	}
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

func ImFontAtlasBuildMultiplyCalcLookupTable(out_table []byte, in_brighten_factor float32) {
	for i := uint(0); i < 256; i++ {
		var value = uint(float(i) * in_brighten_factor)
		if value > 255 {
			value = 255
		}
		out_table[i] = byte(value & 0xFF)
	}
}

func ImFontAtlasBuildMultiplyRectAlpha8(table []byte, pixels []byte, x, y, w, h, stride int) {
	var data = pixels[x+y*stride:]
	for j := h; j > 0; j, data = j-1, data[stride:] {
		for i := int(0); i < w; i++ {
			data[i] = table[data[i]]
		}
	}
}
