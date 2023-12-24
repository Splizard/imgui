package imgui

import "fmt"

// BeginTable Tables
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
//   - In most situations you can use TableNextRow() + TableSetColumnIndex(N) to start appending into a column.
//   - If you are using tables as a sort of grid, where every columns is holding the same type of contents,
//     you may prefer using TableNextColumn() instead of TableNextRow() + TableSetColumnIndex().
//     TableNextColumn() will automatically wrap-around into the next row if needed.
//   - IMPORTANT: Comparatively to the old Columns() API, we need to call TableNextColumn() for the first column!
//   - Summary of possible call flow:
//     --------------------------------------------------------------------------------------------------------
//     TableNextRow() -> TableSetColumnIndex(0) -> Text("Hello 0") -> TableSetColumnIndex(1) -> Text("Hello 1")  // OK
//     TableNextRow() -> TableNextColumn()      -> Text("Hello 0") -> TableNextColumn()      -> Text("Hello 1")  // OK
//     TableNextColumn()      -> Text("Hello 0") -> TableNextColumn()      -> Text("Hello 1")  // OK: TableNextColumn() automatically gets to next row!
//     TableNextRow()                           -> Text("Hello 0")                                               // Not OK! Missing TableSetColumnIndex() or TableNextColumn()! Text will not appear!
//     --------------------------------------------------------------------------------------------------------
//
// - 5. Call EndTable()
// Read about "TABLE SIZING" at the top of this file.
func BeginTable(str_id string, columns_count int, flags ImGuiTableFlags, outer_size ImVec2, inner_width float) bool {
	var id = GetIDs(str_id)
	return BeginTableEx(str_id, id, columns_count, flags, &outer_size, inner_width)
}

// EndTable only call EndTable() if BeginTable() returns true!
func EndTable() {
	var table = guiContext.CurrentTable
	IM_ASSERT_USER_ERROR(table != nil, "Only call EndTable() if BeginTable() returns true!")

	// This assert would be very useful to catch a common error... unfortunately it would probably trigger in some
	// cases, and for consistency user may sometimes output empty tables (and still benefit from e.guiContext. outer border)
	//IM_ASSERT(table.IsLayoutLocked && "Table unused: never called TableNextRow(), is that the intent?");

	// If the user never got to call TableNextRow() or TableNextColumn(), we call layout ourselves to ensure all our
	// code paths are consistent (instead of just hoping that TableBegin/TableEnd will work), get borders drawn, etc.
	if !table.IsLayoutLocked {
		TableUpdateLayout(table)
	}

	var flags = table.Flags
	var inner_window = table.InnerWindow
	var outer_window = table.OuterWindow
	var temp_data = table.TempData
	IM_ASSERT(inner_window == guiContext.CurrentWindow)
	IM_ASSERT(outer_window == inner_window || outer_window == inner_window.ParentWindow)

	if table.IsInsideRow {
		TableEndRow(table)
	}

	// Context menu in columns body
	if (flags & ImGuiTableFlags_ContextMenuInBody) != 0 {
		if table.HoveredColumnBody != -1 && !IsAnyItemHovered() && IsMouseReleased(ImGuiMouseButton_Right) {
			TableOpenContextMenu((int)(table.HoveredColumnBody))
		}
	}

	// Finalize table height
	inner_window.DC.PrevLineSize = temp_data.HostBackupPrevLineSize
	inner_window.DC.CurrLineSize = temp_data.HostBackupCurrLineSize
	inner_window.DC.CursorMaxPos = temp_data.HostBackupCursorMaxPos
	var inner_content_max_y = table.RowPosY2
	IM_ASSERT(table.RowPosY2 == inner_window.DC.CursorPos.y)
	if inner_window != outer_window {
		inner_window.DC.CursorMaxPos.y = inner_content_max_y
	} else if (flags & ImGuiTableFlags_NoHostExtendY) == 0 {
		table.OuterRect.Max.y = max(table.OuterRect.Max.y, inner_content_max_y) // Patch OuterRect/InnerRect height
		table.InnerRect.Max.y = table.OuterRect.Max.y
	}
	table.WorkRect.Max.y = max(table.WorkRect.Max.y, table.OuterRect.Max.y)
	table.LastOuterHeight = table.OuterRect.GetHeight()

	// Setup inner scrolling range
	// FIXME: This ideally should be done earlier, in BeginTable() SetNextWindowContentSize call, just like writing to inner_window.DC.CursorMaxPos.y,
	// but since the later is likely to be impossible to do we'd rather update both axises together.
	if table.Flags&ImGuiTableFlags_ScrollX != 0 {
		var outer_padding_for_border float
		if table.Flags&ImGuiTableFlags_BordersOuterV != 0 {
			outer_padding_for_border = TABLE_BORDER_SIZE
		}
		var max_pos_x = table.InnerWindow.DC.CursorMaxPos.x
		if table.RightMostEnabledColumn != -1 {
			max_pos_x = max(max_pos_x, table.Columns[table.RightMostEnabledColumn].WorkMaxX+table.CellPaddingX+table.OuterPaddingX-outer_padding_for_border)
		}
		if table.ResizedColumn != -1 {
			max_pos_x = max(max_pos_x, table.ResizeLockMinContentsX2)
		}
		table.InnerWindow.DC.CursorMaxPos.x = max_pos_x
	}

	// Pop clipping rect
	if (flags & ImGuiTableFlags_NoClip) == 0 {
		inner_window.DrawList.PopClipRect()
	}
	inner_window.ClipRect = ImRectFromVec4(&inner_window.DrawList._ClipRectStack[len(inner_window.DrawList._ClipRectStack)-1])

	// Draw borders
	if (flags & ImGuiTableFlags_Borders) != 0 {
		TableDrawBorders(table)
	}

	if false {
		// Strip out dummy channel draw calls
		// We have no way to prevent user submitting direct ImDrawList calls into a hidden column (but ImGui:: calls will be clipped out)
		// Pros: remove draw calls which will have no effect. since they'll have zero-size cliprect they may be early out anyway.
		// Cons: making it harder for users watching metrics/debugger to spot the wasted vertices.
		if ImGuiTableColumnIdx(table.DummyDrawChannel) != (ImGuiTableColumnIdx)(-1) {
			var dummy_channel = &table.DrawSplitter._Channels[table.DummyDrawChannel]
			dummy_channel._CmdBuffer = dummy_channel._CmdBuffer[:0]
			dummy_channel._IdxBuffer = dummy_channel._IdxBuffer[:0]
		}
	}

	// Flatten channels and merge draw calls
	var splitter = table.DrawSplitter
	splitter.SetCurrentChannel(inner_window.DrawList, 0)
	if (table.Flags & ImGuiTableFlags_NoClip) == 0 {
		TableMergeDrawChannels(table)
	}
	splitter.Merge(inner_window.DrawList)

	// Update ColumnsAutoFitWidth to get us ahead for host using our size to auto-resize without waiting for next BeginTable()
	var width_spacings = (table.OuterPaddingX * 2.0) + (table.CellSpacingX1+table.CellSpacingX2)*float(table.ColumnsEnabledCount-1)
	table.ColumnsAutoFitWidth = width_spacings + (table.CellPaddingX*2.0)*float(table.ColumnsEnabledCount)
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		if table.EnabledMaskByIndex&((ImU64)(1<<column_n)) != 0 {
			var column = &table.Columns[column_n]
			if (column.Flags&ImGuiTableColumnFlags_WidthFixed) != 0 && column.Flags&ImGuiTableColumnFlags_NoResize == 0 {
				table.ColumnsAutoFitWidth += column.WidthRequest
			} else {
				table.ColumnsAutoFitWidth += TableGetColumnWidthAuto(table, column)
			}
		}
	}

	// Update scroll
	if (table.Flags&ImGuiTableFlags_ScrollX) == 0 && inner_window != outer_window {
		inner_window.Scroll.x = 0.0
	} else if table.LastResizedColumn != -1 && table.ResizedColumn == -1 && inner_window.ScrollbarX && table.InstanceInteracted == table.InstanceCurrent {
		// When releasing a column being resized, scroll to keep the resulting column in sight
		var neighbor_width_to_keep_visible = table.MinColumnWidth + table.CellPaddingX*2.0
		var column = &table.Columns[table.LastResizedColumn]
		if column.MaxX < table.InnerClipRect.Min.x {
			setScrollFromPosX(inner_window, column.MaxX-inner_window.Pos.x-neighbor_width_to_keep_visible, 1.0)
		} else if column.MaxX > table.InnerClipRect.Max.x {
			setScrollFromPosX(inner_window, column.MaxX-inner_window.Pos.x+neighbor_width_to_keep_visible, 1.0)
		}
	}

	// Apply resizing/dragging at the end of the frame
	if table.ResizedColumn != -1 && table.InstanceCurrent == table.InstanceInteracted {
		var column = &table.Columns[table.ResizedColumn]
		var new_x2 = (guiContext.IO.MousePos.x - guiContext.ActiveIdClickOffset.x + TABLE_RESIZE_SEPARATOR_HALF_THICKNESS)
		var new_width = ImFloor(new_x2 - column.MinX - table.CellSpacingX1 - table.CellPaddingX*2.0)
		table.ResizedColumnNextWidth = new_width
	}

	// Pop from id stack
	IM_ASSERT_USER_ERROR(inner_window.IDStack[len(inner_window.IDStack)-1] == table.ID+uint(table.InstanceCurrent), "Mismatching PushID/PopID!")
	IM_ASSERT_USER_ERROR(int(len(outer_window.DC.ItemWidthStack)) >= temp_data.HostBackupItemWidthStackSize, "Too many PopItemWidth!")
	PopID()

	// Restore window data that we modified
	var backup_outer_max_pos = outer_window.DC.CursorMaxPos
	inner_window.WorkRect = temp_data.HostBackupWorkRect
	inner_window.ParentWorkRect = temp_data.HostBackupParentWorkRect
	inner_window.SkipItems = table.HostSkipItems
	outer_window.DC.CursorPos = table.OuterRect.Min
	outer_window.DC.ItemWidth = temp_data.HostBackupItemWidth
	//outer_window.DC.ItemWidthStack.Size = temp_data.HostBackupItemWidthStackSize //FIXME?
	outer_window.DC.ColumnsOffset = temp_data.HostBackupColumnsOffset

	// Layout in outer window
	// (FIXME: To allow auto-fit and allow desirable effect of SameLine() we dissociate 'used' vs 'ideal' size by overriding
	// CursorPosPrevLine and CursorMaxPos manually. That should be a more general layout feature, see same problem e.guiContext. #3414)
	if inner_window != outer_window {
		EndChild()
	} else {
		size := table.OuterRect.GetSize()
		ItemSizeVec(&size, 0)
		ItemAdd(&table.OuterRect, 0, nil, 0)
	}

	// Override declared contents width/height to enable auto-resize while not needlessly adding a scrollbar
	if table.Flags&ImGuiTableFlags_NoHostExtendX != 0 {
		// FIXME-TABLE: Could we remove this section?
		// ColumnsAutoFitWidth may be one frame ahead here since for Fixed+NoResize is calculated from latest contents
		IM_ASSERT((table.Flags & ImGuiTableFlags_ScrollX) == 0)
		outer_window.DC.CursorMaxPos.x = max(backup_outer_max_pos.x, table.OuterRect.Min.x+table.ColumnsAutoFitWidth)
	} else if temp_data.UserOuterSize.x <= 0.0 {
		var decoration_size float = 0.0
		if table.Flags&ImGuiTableFlags_ScrollX != 0 {
			decoration_size = inner_window.ScrollbarSizes.x
		}
		outer_window.DC.IdealMaxPos.x = max(outer_window.DC.IdealMaxPos.x, table.OuterRect.Min.x+table.ColumnsAutoFitWidth+decoration_size-temp_data.UserOuterSize.x)
		outer_window.DC.CursorMaxPos.x = max(backup_outer_max_pos.x, min(table.OuterRect.Max.x, table.OuterRect.Min.x+table.ColumnsAutoFitWidth))
	} else {
		outer_window.DC.CursorMaxPos.x = max(backup_outer_max_pos.x, table.OuterRect.Max.x)
	}
	if temp_data.UserOuterSize.y <= 0.0 {
		var decoration_size float = 0.0
		if (table.Flags & ImGuiTableFlags_ScrollY) != 0 {
			decoration_size = inner_window.ScrollbarSizes.y
		}
		outer_window.DC.IdealMaxPos.y = max(outer_window.DC.IdealMaxPos.y, inner_content_max_y+decoration_size-temp_data.UserOuterSize.y)
		outer_window.DC.CursorMaxPos.y = max(backup_outer_max_pos.y, min(table.OuterRect.Max.y, inner_content_max_y))
	} else {
		// OuterRect.Max.y may already have been pushed downward from the initial value (unless ImGuiTableFlags_NoHostExtendY is set)
		outer_window.DC.CursorMaxPos.y = max(backup_outer_max_pos.y, table.OuterRect.Max.y)
	}

	// Save settings
	if table.IsSettingsDirty {
		TableSaveSettings(table)
	}
	table.IsInitializing = false

	// Clear or restore current table, if any
	IM_ASSERT(guiContext.CurrentWindow == outer_window && guiContext.CurrentTable == table)
	IM_ASSERT(guiContext.CurrentTableStackIdx >= 0)
	guiContext.CurrentTableStackIdx--
	if guiContext.CurrentTableStackIdx >= 0 {
		temp_data = &guiContext.TablesTempDataStack[guiContext.CurrentTableStackIdx]
	} else {
		temp_data = nil
	}
	if temp_data != nil {
		guiContext.CurrentTable = guiContext.Tables[uint(temp_data.TableIndex)]
	} else {
		guiContext.CurrentTable = nil
	}
	if guiContext.CurrentTable != nil {
		guiContext.CurrentTable.TempData = temp_data
		guiContext.CurrentTable.DrawSplitter = &temp_data.DrawSplitter
	}
	outer_window.DC.CurrentTableIdx = -1
	if guiContext.CurrentTable != nil {
		for i, table := range guiContext.Tables {
			if table == guiContext.CurrentTable {
				outer_window.DC.CurrentTableIdx = int(i)
				break
			}
		}
	}
}

// TableNextRow [Public] Starts into the first cell of a new row
func TableNextRow(row_flags ImGuiTableRowFlags /*= 0*/, row_min_height float) {
	var table = guiContext.CurrentTable

	if !table.IsLayoutLocked {
		TableUpdateLayout(table)
	}
	if table.IsInsideRow {
		TableEndRow(table)
	}

	table.LastRowFlags = table.RowFlags
	table.RowFlags = row_flags
	table.RowMinHeight = row_min_height
	TableBeginRow(table)

	// We honor min_row_height requested by user, but cannot guarantee per-row maximum height,
	// because that would essentially require a unique clipping rectangle per-cell.
	table.RowPosY2 += table.CellPaddingY * 2.0
	table.RowPosY2 = max(table.RowPosY2, table.RowPosY1+row_min_height)

	// Disable output until user calls TableNextColumn()
	table.InnerWindow.SkipItems = true
}

// TableNextColumn [Public] Append into the next column, wrap and create a new row when already on last column
// append into the first cell of a new row.
// append into the next column (or first column of next row if currently in last column). Return true when column is visible.
func TableNextColumn() bool {
	var table = guiContext.CurrentTable
	if table == nil {
		return false
	}

	if table.IsInsideRow && table.CurrentColumn+1 < table.ColumnsCount {
		if table.CurrentColumn != -1 {
			TableEndCell(table)
		}
		TableBeginCell(table, table.CurrentColumn+1)
	} else {
		TableNextRow(0, 0)
		TableBeginCell(table, 0)
	}

	// Return whether the column is visible. User may choose to skip submitting items based on this return value,
	// however they shouldn't skip submitting for columns that may have the tallest contribution to row height.
	var column_n = table.CurrentColumn
	return (table.RequestOutputMaskByIndex & ((ImU64)(1 << column_n))) != 0
}

// TableSetColumnIndex [Public] Append into a specific column
// append into the specified column. Return true when column is visible.
func TableSetColumnIndex(column_n int) bool {
	var table = guiContext.CurrentTable
	if table == nil {
		return false
	}

	if table.CurrentColumn != column_n {
		if table.CurrentColumn != -1 {
			TableEndCell(table)
		}
		IM_ASSERT(column_n >= 0 && table.ColumnsCount != 0)
		TableBeginCell(table, column_n)
	}

	// Return whether the column is visible. User may choose to skip submitting items based on this return value,
	// however they shouldn't skip submitting for columns that may have the tallest contribution to row height.
	return (table.RequestOutputMaskByIndex & ((ImU64)(1 << column_n))) != 0
}

// TableSetupColumn Tables: Headers & Columns declaration
//   - Use TableSetupColumn() to specify label, resizing policy, default width/weight, id, various other flags etc.
//   - Use TableHeadersRow() to create a header row and automatically submit a TableHeader() for each column.
//     Headers are required to perform: reordering, sorting, and opening the context menu.
//     The context menu can also be made available in columns body using ImGuiTableFlags_ContextMenuInBody.
//   - You may manually submit headers using TableNextRow() + TableHeader() calls, but this is only useful in
//     some advanced use cases (e.guiContext. adding custom widgets in header row).
//   - Use TableSetupScrollFreeze() to lock columns/rows so they stay visible when scrolled.
//
// See "COLUMN SIZING POLICIES" comments at the top of this file
// If (init_width_or_weight <= 0.0f) it is ignored
func TableSetupColumn(label string, flags ImGuiTableColumnFlags, init_width_or_weight float /*= 0*/, user_id ImGuiID) {
	var table = guiContext.CurrentTable
	IM_ASSERT_USER_ERROR(table != nil, "Need to call TableSetupColumn() after BeginTable()!")
	IM_ASSERT_USER_ERROR(!table.IsLayoutLocked, "Need to call call TableSetupColumn() before first row!")
	IM_ASSERT_USER_ERROR((flags&ImGuiTableColumnFlags_StatusMask_) == 0, "Illegal to pass StatusMask values to TableSetupColumn()")
	if int(table.DeclColumnsCount) >= table.ColumnsCount {
		IM_ASSERT_USER_ERROR(int(table.DeclColumnsCount) < table.ColumnsCount, "Called TableSetupColumn() too many times!")
		return
	}

	var column = &table.Columns[table.DeclColumnsCount]
	table.DeclColumnsCount++

	// Assert when passing a width or weight if policy is entirely left to default, to avoid storing width into weight and vice-versa.
	// Give a grace to users of ImGuiTableFlags_ScrollX.
	if table.IsDefaultSizingPolicy && (flags&ImGuiTableColumnFlags_WidthMask_) == 0 && (flags&ImGuiTableColumnFlags(ImGuiTableFlags_ScrollX)) == 0 {
		IM_ASSERT_USER_ERROR(init_width_or_weight <= 0.0, "Can only specify width/weight if sizing policy is set explicitly in either Table or Column.")
	}

	// When passing a width automatically enforce WidthFixed policy
	// (whereas TableSetupColumnFlags would default to WidthAuto if table is not Resizable)
	if (flags&ImGuiTableColumnFlags_WidthMask_) == 0 && init_width_or_weight > 0.0 {
		if (table.Flags&ImGuiTableFlags_SizingMask_) == ImGuiTableFlags_SizingFixedFit || (table.Flags&ImGuiTableFlags_SizingMask_) == ImGuiTableFlags_SizingFixedSame {
			flags |= ImGuiTableColumnFlags_WidthFixed
		}
	}

	TableSetupColumnFlags(table, column, flags)
	column.UserID = user_id
	flags = column.Flags

	// Initialize defaults
	column.InitStretchWeightOrWidth = init_width_or_weight
	if table.IsInitializing {
		// Init width or weight
		if column.WidthRequest < 0.0 && column.StretchWeight < 0.0 {
			if (flags&ImGuiTableColumnFlags_WidthFixed) != 0 && init_width_or_weight > 0.0 {
				column.WidthRequest = init_width_or_weight
			}
			if flags&ImGuiTableColumnFlags_WidthStretch != 0 {
				if init_width_or_weight > 0.0 {
					column.StretchWeight = init_width_or_weight
				} else {
					column.StretchWeight = -1.0
				}
			}

			// Disable auto-fit if an explicit width/weight has been specified
			if init_width_or_weight > 0.0 {
				column.AutoFitQueue = 0x00
			}
		}

		// Init default visibility/sort state
		if (flags&ImGuiTableColumnFlags_DefaultHide) != 0 && (table.SettingsLoadedFlags&ImGuiTableFlags_Hideable) == 0 {
			column.IsUserEnabled = false
			column.IsUserEnabledNextFrame = false
		}
		if flags&ImGuiTableColumnFlags_DefaultSort != 0 && (table.SettingsLoadedFlags&ImGuiTableFlags_Sortable) == 0 {
			column.SortOrder = 0 // Multiple columns using _DefaultSort will be reassigned unique SortOrder values when building the sort specs.
			if column.Flags&ImGuiTableColumnFlags_PreferSortDescending != 0 {
				column.SortDirection = ImGuiSortDirection_Descending
			} else {
				column.SortDirection = ImGuiSortDirection_Ascending
			}
		}
	}

	// Store name (append with zero-terminator in contiguous buffer)
	column.NameOffset = -1
	if label != "" {
		column.NameOffset = (ImS16)(len(table.ColumnsNames))
		table.ColumnsNames = append(table.ColumnsNames, label)
	}
}

// TableSetupScrollFreeze [Public]
// lock columns/rows so they stay visible when scrolled.
func TableSetupScrollFreeze(columns int, rows int) {
	var table = guiContext.CurrentTable
	IM_ASSERT_USER_ERROR(table != nil, "Need to call TableSetupColumn() after BeginTable()!")
	IM_ASSERT_USER_ERROR(!table.IsLayoutLocked, "Need to call TableSetupColumn() before first row!")
	IM_ASSERT(columns >= 0 && columns < IMGUI_TABLE_MAX_COLUMNS)
	IM_ASSERT(rows >= 0 && rows < 128) // Arbitrary limit

	if (table.Flags & ImGuiTableFlags_ScrollX) != 0 {
		table.FreezeColumnsRequest = (ImGuiTableColumnIdx)(min(columns, table.ColumnsCount))
	} else {
		table.FreezeColumnsRequest = 0
	}

	if table.InnerWindow.Scroll.x != 0.0 {
		table.FreezeColumnsCount = table.FreezeColumnsRequest
	} else {
		table.FreezeColumnsCount = 0
	}
	if table.Flags&ImGuiTableFlags_ScrollY != 0 {
		table.FreezeRowsRequest = (ImGuiTableColumnIdx)(rows)
	} else {
		table.FreezeRowsRequest = 0
	}
	if table.InnerWindow.Scroll.y != 0.0 {
		table.FreezeRowsCount = table.FreezeRowsRequest
	} else {
		table.FreezeRowsCount = 0
	}

	table.IsUnfrozenRows = (table.FreezeRowsCount == 0) // Make sure this is set before TableUpdateLayout() so ImGuiListClipper can benefit from it.b

	// Ensure frozen columns are ordered in their section. We still allow multiple frozen columns to be reordered.
	for column_n := int8(0); column_n < table.FreezeColumnsRequest; column_n++ {
		var order_n = int(table.DisplayOrderToIndex[column_n])
		if order_n != int(column_n) && order_n >= int(table.FreezeColumnsRequest) {
			//swap
			table.Columns[table.DisplayOrderToIndex[order_n]].DisplayOrder, table.Columns[table.DisplayOrderToIndex[column_n]].DisplayOrder = table.Columns[table.DisplayOrderToIndex[column_n]].DisplayOrder, table.Columns[table.DisplayOrderToIndex[order_n]].DisplayOrder
			table.DisplayOrderToIndex[order_n], table.DisplayOrderToIndex[column_n] = table.DisplayOrderToIndex[column_n], table.DisplayOrderToIndex[order_n]
		}
	}
}

// TableHeadersRow [Public] This is a helper to output TableHeader() calls based on the column names declared in TableSetupColumn().
// The intent is that advanced users willing to create customized headers would not need to use this helper
// and can create their own! For example: TableHeader() may be preceeded by Checkbox() or other custom widgets.
// See 'Demo.Tables.Custom headers' for a demonstration of implementing a custom version of this.
// This code is constructed to not make much use of internal functions, as it is intended to be a template to copy.
// FIXME-TABLE: TableOpenContextMenu() and TableGetHeaderRowHeight() are not public.
// submit all headers cells based on data provided to TableSetupColumn() + submit context menu
func TableHeadersRow() {
	var table = guiContext.CurrentTable
	IM_ASSERT_USER_ERROR(table != nil, "Need to call TableHeadersRow() after BeginTable()!")

	// Layout if not already done (this is automatically done by TableNextRow, we do it here solely to facilitate stepping in debugger as it is frequent to step in TableUpdateLayout)
	if !table.IsLayoutLocked {
		TableUpdateLayout(table)
	}

	// Open row
	var row_y1 = GetCursorScreenPos().y
	var row_height = TableGetHeaderRowHeight()
	TableNextRow(ImGuiTableRowFlags_Headers, row_height)
	if table.HostSkipItems { // Merely an optimization, you may skip in your own code.
		return
	}

	var columns_count = TableGetColumnCount()
	for column_n := int(0); column_n < columns_count; column_n++ {
		if !TableSetColumnIndex(column_n) {
			continue
		}

		// Push an id to allow unnamed labels (generally accidental, but let's behave nicely with them)
		// - in your own code you may omit the PushID/PopID all-together, provided you know they won't collide
		// - table.InstanceCurrent is only >0 when we use multiple BeginTable/EndTable calls with same identifier.
		var name string
		if (TableGetColumnFlags(column_n) & ImGuiTableColumnFlags_NoHeaderLabel) != 0 {
			name = TableGetColumnName(column_n)
		}
		PushID(int(table.InstanceCurrent)*table.ColumnsCount + column_n)
		TableHeader(name)
		PopID()
	}

	// Allow opening popup from the right-most section after the last column.
	var mouse_pos = GetMousePos()
	if IsMouseReleased(1) && TableGetHoveredColumn() == columns_count {
		if mouse_pos.y >= row_y1 && mouse_pos.y < row_y1+row_height {
			TableOpenContextMenu(-1) // Will open a non-column-specific popup.
		}
	}
}

// TableHeader Emit a column header (text + optional sort order)
// We cpu-clip text here so that all columns headers can be merged into a same draw call.
// Note that because of how we cpu-clip and display sorting indicators, you _cannot_ use SameLine() after a TableHeader()
// submit one header cell manually (rarely used)
func TableHeader(label string) {
	window := guiContext.CurrentWindow
	if window.SkipItems {
		return
	}

	var table = guiContext.CurrentTable
	IM_ASSERT_USER_ERROR(table != nil, "Need to call TableHeader() after BeginTable()!")
	IM_ASSERT(table.CurrentColumn != -1)
	var column_n = table.CurrentColumn
	var column = &table.Columns[column_n]

	// Label
	var label_size = CalcTextSize(label, true, -1)
	var label_pos = window.DC.CursorPos

	// If we already got a row height, there's use that.
	// FIXME-TABLE: Padding problem if the correct outer-padding CellBgRect strays off our ClipRect?
	var cell_r = TableGetCellBgRect(table, column_n)
	var label_height = max(label_size.y, table.RowMinHeight-table.CellPaddingY*2.0)

	// Calculate ideal size for sort order arrow
	var w_arrow float = 0.0
	var w_sort_text float = 0.0
	var sort_order_suf string
	var ARROW_SCALE float = 0.65
	if (table.Flags&ImGuiTableFlags_Sortable) != 0 && (column.Flags&ImGuiTableColumnFlags_NoSort) == 0 {
		w_arrow = ImFloor(guiContext.FontSize*ARROW_SCALE + guiContext.Style.FramePadding.x)
		if column.SortOrder > 0 {
			sort_order_suf = fmt.Sprintf("%d", column.SortOrder+1)
			w_sort_text = guiContext.Style.ItemInnerSpacing.x + CalcTextSize(sort_order_suf, true, -1).x
		}
	}

	// We feed our unclipped width to the column without writing on CursorMaxPos, so that column is still considering for merging.
	var max_pos_x = label_pos.x + label_size.x + w_sort_text + w_arrow
	column.ContentMaxXHeadersUsed = max(column.ContentMaxXHeadersUsed, column.WorkMaxX)
	column.ContentMaxXHeadersIdeal = max(column.ContentMaxXHeadersIdeal, max_pos_x)

	// Keep header highlighted when context menu is open.
	var selected = (table.IsContextPopupOpen && int(table.ContextPopupColumn) == column_n && table.InstanceInteracted == table.InstanceCurrent)
	var id = window.GetIDs(label)
	var bb = ImRect{ImVec2{cell_r.Min.x, cell_r.Min.y}, ImVec2{cell_r.Max.x, max(cell_r.Max.y, cell_r.Min.y+label_height+guiContext.Style.CellPadding.y*2.0)}}
	ItemSizeVec(&ImVec2{0.0, label_height}, 0) // Don't declare unclipped width, it'll be fed ContentMaxPosHeadersIdeal
	if !ItemAdd(&bb, id, nil, 0) {
		return
	}

	//GetForegroundDrawList().AddRect(cell_r.Min, cell_r.Max, IM_COL32(255, 0, 0, 255)); // [DEBUG]
	//GetForegroundDrawList().AddRect(bb.Min, bb.Max, IM_COL32(255, 0, 0, 255)); // [DEBUG]

	// Using AllowItemOverlap mode because we cover the whole cell, and we want user to be able to submit subsequent items.
	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, ImGuiButtonFlags_AllowItemOverlap)
	if guiContext.ActiveId != id {
		SetItemAllowOverlap()
	}
	if held || hovered || selected {
		var col ImU32
		if held {
			col = GetColorU32FromID(ImGuiCol_HeaderActive, 1)
		} else if hovered {
			col = GetColorU32FromID(ImGuiCol_HeaderHovered, 1)
		} else {
			col = GetColorU32FromID(ImGuiCol_Header, 1)
		}
		//RenderFrame(bb.Min, bb.Max, col, false, 0.0f);
		TableSetBgColor(ImGuiTableBgTarget_CellBg, col, table.CurrentColumn)
	} else {
		// Submit single cell bg color in the case we didn't submit a full header row
		if (table.RowFlags & ImGuiTableRowFlags_Headers) == 0 {
			TableSetBgColor(ImGuiTableBgTarget_CellBg, GetColorU32FromID(ImGuiCol_TableHeaderBg, 1), table.CurrentColumn)
		}
	}
	RenderNavHighlight(&bb, id, ImGuiNavHighlightFlags_TypeThin|ImGuiNavHighlightFlags_NoRounding)
	if held {
		table.HeldHeaderColumn = (ImGuiTableColumnIdx)(column_n)
	}
	window.DC.CursorPos.y -= guiContext.Style.ItemSpacing.y * 0.5

	// Drag and drop to re-order columns.
	// FIXME-TABLE: Scroll request while reordering a column and it lands out of the scrolling zone.
	if held && (table.Flags&ImGuiTableFlags_Reorderable) != 0 && IsMouseDragging(0, -1) && !guiContext.DragDropActive {
		// While moving a column it will jump on the other side of the mouse, so we also test for MouseDelta.x
		table.ReorderColumn = (ImGuiTableColumnIdx)(column_n)
		table.InstanceInteracted = table.InstanceCurrent

		// We don't reorder: through the frozen<>unfrozen line, or through a column that is marked with ImGuiTableColumnFlags_NoReorder.
		if guiContext.IO.MouseDelta.x < 0.0 && guiContext.IO.MousePos.x < cell_r.Min.x {

			var prev_column *ImGuiTableColumn
			if column.PrevEnabledColumn != -1 {
				prev_column = &table.Columns[column.PrevEnabledColumn]
			}

			if prev_column != nil {
				if ((column.Flags | prev_column.Flags) & ImGuiTableColumnFlags_NoReorder) == 0 {
					if (column.IndexWithinEnabledSet < table.FreezeColumnsRequest) == (prev_column.IndexWithinEnabledSet < table.FreezeColumnsRequest) {
						table.ReorderColumnDir = -1
					}
				}
			}
		}
		if guiContext.IO.MouseDelta.x > 0.0 && guiContext.IO.MousePos.x > cell_r.Max.x {

			var next_column *ImGuiTableColumn
			if column.NextEnabledColumn != -1 {
				next_column = &table.Columns[column.NextEnabledColumn]
			}

			if next_column != nil {
				if ((column.Flags | next_column.Flags) & ImGuiTableColumnFlags_NoReorder) == 0 {
					if (column.IndexWithinEnabledSet < table.FreezeColumnsRequest) == (next_column.IndexWithinEnabledSet < table.FreezeColumnsRequest) {
						table.ReorderColumnDir = +1
					}
				}
			}
		}
	}

	// Sort order arrow
	var ellipsis_max = cell_r.Max.x - w_arrow - w_sort_text
	if (table.Flags&ImGuiTableFlags_Sortable) != 0 && (column.Flags&ImGuiTableColumnFlags_NoSort) == 0 {
		if column.SortOrder != -1 {
			var x = max(cell_r.Min.x, cell_r.Max.x-w_arrow-w_sort_text)
			var y = label_pos.y
			if column.SortOrder > 0 {
				PushStyleColorInt(ImGuiCol_Text, GetColorU32FromID(ImGuiCol_Text, 0.70))
				RenderText(ImVec2{x + guiContext.Style.ItemInnerSpacing.x, y}, sort_order_suf, true)
				PopStyleColor(1)
				x += w_sort_text
			}

			var dir = ImGuiDir_Down
			if column.SortDirection == ImGuiSortDirection_Ascending {
				dir = ImGuiDir_Up
			}

			RenderArrow(window.DrawList, ImVec2{x, y}, GetColorU32FromID(ImGuiCol_Text, 1), dir, ARROW_SCALE)
		}

		// Handle clicking on column header to adjust Sort Order
		if pressed && int(table.ReorderColumn) != column_n {
			var sort_direction = TableGetColumnNextSortDirection(column)
			TableSetColumnSortDirection(column_n, sort_direction, guiContext.IO.KeyShift)
		}
	}

	// Render clipped label. Clipping here ensure that in the majority of situations, all our header cells will
	// be merged into a single draw call.
	//window.DrawList.AddCircleFilled(ImVec2(ellipsis_max, label_pos.y), 40, IM_COL32_WHITE);
	RenderTextEllipsis(window.DrawList, &label_pos, &ImVec2{ellipsis_max, label_pos.y + label_height + guiContext.Style.FramePadding.y}, ellipsis_max, ellipsis_max, label, &label_size)

	var text_clipped = label_size.x > (ellipsis_max - label_pos.x)
	if text_clipped && hovered && guiContext.HoveredIdNotActiveTimer > guiContext.TooltipSlowDelay {
		SetTooltip("%.*s", (int)(len(label)), label)
	}

	// We don't use BeginPopupContextItem() because we want the popup to stay up even after the column is hidden
	if IsMouseReleased(1) && IsItemHovered(0) {
		TableOpenContextMenu(column_n)
	}
}

// TableGetSortSpecs Tables: Sorting
//   - Call TableGetSortSpecs() to retrieve latest sort specs for the table. NULL when not sorting.
//   - When 'SpecsDirty == true' you should sort your data. It will be true when sorting specs have changed
//     since last call, or the first time. Make sure to set 'SpecsDirty/*= guiContext*/,else you may
//     wastefully sort your data every frame!
//   - Lifetime: don't hold on this pointer over multiple frames or past any subsequent call to BeginTable().
//
// get latest sort specs for the table (NULL if not sorting).
func TableGetSortSpecs() *ImGuiTableSortSpecs {
	var table = guiContext.CurrentTable
	IM_ASSERT(table != nil)

	if (table.Flags & ImGuiTableFlags_Sortable) == 0 {
		return nil
	}

	// Require layout (in case TableHeadersRow() hasn't been called) as it may alter IsSortSpecsDirty in some paths.
	if !table.IsLayoutLocked {
		TableUpdateLayout(table)
	}

	TableSortSpecsBuild(table)

	return &table.SortSpecs
}

func TableGetColumnAvailSortDirection(column *ImGuiTableColumn, n int) ImGuiSortDirection {
	IM_ASSERT(n < int(column.SortDirectionsAvailCount))
	return ImGuiSortDirection((column.SortDirectionsAvailList >> (n << 1)) & 0x03)
}

// TableGetColumnCount Tables: Miscellaneous functions
// - Functions args 'column_n int' treat the default value of -1 as the same as passing the current column index.
// return number of columns (value passed to BeginTable)
func TableGetColumnCount() int {
	var table = guiContext.CurrentTable
	if table != nil {
		return int(table.ColumnsCount)
	}
	return 0
}

// TableGetColumnIndex return current column index.
func TableGetColumnIndex() int {
	var table = guiContext.CurrentTable
	if table == nil {
		return 0
	}
	return table.CurrentColumn
}

// TableGetRowIndex [Public] Note: for row coloring we use .RowBgColorCounter which is the same value without counting header rows
// return current row index.
func TableGetRowIndex() int {
	var table = guiContext.CurrentTable
	if table == nil {
		return 0
	}
	return table.CurrentRow
}

// TableGetColumnName return "" if column didn't have a name declared by TableSetupColumn(). Pass -1 to use current column.
func TableGetColumnName(column_n int /*= -1*/) string {
	var table = guiContext.CurrentTable
	if table == nil {
		return ""
	}
	if column_n < 0 {
		column_n = table.CurrentColumn
	}
	return tableGetColumnName(table, column_n)
}

// TableGetColumnFlags We allow querying for an extra column in order to poll the IsHovered state of the right-most section
// return column flags so you can query their Enabled/Visible/Sorted/Hovered status flags. Pass -1 to use current column.
func TableGetColumnFlags(column_n int /*= -1*/) ImGuiTableColumnFlags {
	var table = guiContext.CurrentTable
	if table == nil {
		return ImGuiTableColumnFlags_None
	}
	if column_n < 0 {
		column_n = table.CurrentColumn
	}
	if column_n == table.ColumnsCount {
		if int(table.HoveredColumnBody) == column_n {
			return ImGuiTableColumnFlags_IsHovered
		}
		return ImGuiTableColumnFlags_None
	}
	return table.Columns[column_n].Flags
}

// TableSetColumnEnabled Change user accessible enabled/disabled state of a column (often perceived as "showing/hiding" from users point of view)
// Note that end-user can use the context menu to change this themselves (right-click in headers, or right-click in columns body with ImGuiTableFlags_ContextMenuInBody)
// - Require table to have the ImGuiTableFlags_Hideable flag because we are manipulating user accessible state.
// - Request will be applied during next layout, which happens on the first call to TableNextRow() after BeginTable().
// - For the getter you can test (TableGetColumnFlags() & ImGuiTableColumnFlags_IsEnabled) != 0.
// - Alternative: the ImGuiTableColumnFlags_Disabled is an overriding/master disable flag which will also hide the column from context menu.
// change user accessible enabled/disabled state of a column. Set to false to hide the column. User can use the context menu to change this themselves (right-click in headers, or right-click in columns body with ImGuiTableFlags_ContextMenuInBody)
func TableSetColumnEnabled(column_n int, enabled bool) {
	var table = guiContext.CurrentTable
	IM_ASSERT(table != nil)
	if table == nil {
		return
	}
	IM_ASSERT(table.Flags&ImGuiTableFlags_Hideable != 0) // See comments above
	if column_n < 0 {
		column_n = table.CurrentColumn
	}
	IM_ASSERT(column_n >= 0 && column_n < table.ColumnsCount)
	var column = &table.Columns[column_n]
	column.IsUserEnabledNextFrame = enabled
}

// TableSetBgColor change the color of a cell, row, or column. See ImGuiTableBgTarget_ flags for details.
func TableSetBgColor(target ImGuiTableBgTarget, color ImU32, column_n int /*= -1*/) {
	var table = guiContext.CurrentTable
	IM_ASSERT(target != ImGuiTableBgTarget_None)

	if color == IM_COL32_DISABLE {
		color = 0
	}

	// We cannot draw neither the cell or row background immediately as we don't know the row height at this point in time.
	switch target {
	case ImGuiTableBgTarget_CellBg:
		if table.RowPosY1 > table.InnerClipRect.Max.y { // Discard
			return
		}
		if column_n == -1 {
			column_n = table.CurrentColumn
		}
		if (table.VisibleMaskByIndex & ((ImU64)(1 << column_n))) == 0 {
			return
		}
		if table.RowCellDataCurrent < 0 || int(table.RowCellData[table.RowCellDataCurrent].Column) != column_n {
			table.RowCellDataCurrent++
		}
		var cell_data = &table.RowCellData[table.RowCellDataCurrent]
		cell_data.BgColor = color
		cell_data.Column = (ImGuiTableColumnIdx)(column_n)
	case ImGuiTableBgTarget_RowBg0:
		fallthrough
	case ImGuiTableBgTarget_RowBg1:
		if table.RowPosY1 > table.InnerClipRect.Max.y { // Discard
			return
		}
		IM_ASSERT(column_n == -1)
		var bg_idx int
		if target == ImGuiTableBgTarget_RowBg1 {
			bg_idx = 1
		}
		table.RowBgColor[bg_idx] = color
	default:
		IM_ASSERT(false)
	}
}
