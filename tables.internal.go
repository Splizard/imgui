package imgui

import (
	"fmt"
	"unsafe"
)

const (
	TABLE_DRAW_CHANNEL_BG0                int   = 0
	TABLE_DRAW_CHANNEL_BG2_FROZEN         int   = 1
	TABLE_DRAW_CHANNEL_NOCLIP             int   = 2
	TABLE_BORDER_SIZE                     float = 1.0
	TABLE_RESIZE_SEPARATOR_HALF_THICKNESS float = 4.0
	TABLE_RESIZE_SEPARATOR_FEEDBACK_TIMER float = 0.06
)

func TableSettingsCalcChunkSize(columns_count int) size_t {
	return unsafe.Sizeof(ImGuiTableSettings{}) + (size_t)(columns_count)*unsafe.Sizeof(ImGuiTableColumnSettings{})
}

// Clear and initialize empty settings instance
func TableSettingsInit(settings *ImGuiTableSettings, id ImGuiID, columns_count, columns_count_max int) {
	*settings = ImGuiTableSettings{}
	settings.ID = id
	settings.ColumnsCount = (ImGuiTableColumnIdx)(columns_count)
	settings.ColumnsCountMax = (ImGuiTableColumnIdx)(columns_count_max)
	settings.WantApply = true
}

// Helper
func TableFixFlags(flags ImGuiTableFlags, outer_window *ImGuiWindow) ImGuiTableFlags {
	// Adjust flags: set default sizing policy
	if (flags & ImGuiTableFlags_SizingMask_) == 0 {
		if outer_window.Flags&ImGuiWindowFlags_AlwaysAutoResize != 0 {
			flags |= (flags & ImGuiTableFlags_ScrollX) | ImGuiTableFlags_SizingFixedFit
		} else {
			flags |= (flags & ImGuiTableFlags_ScrollX) | ImGuiTableFlags_SizingStretchSame
		}
	}

	// Adjust flags: enable NoKeepColumnsVisible when using ImGuiTableFlags_SizingFixedSame
	if (flags & ImGuiTableFlags_SizingMask_) == ImGuiTableFlags_SizingFixedSame {
		flags |= ImGuiTableFlags_NoKeepColumnsVisible
	}

	// Adjust flags: enforce borders when resizable
	if flags&ImGuiTableFlags_Resizable != 0 {
		flags |= ImGuiTableFlags_BordersInnerV
	}

	// Adjust flags: disable NoHostExtendX/NoHostExtendY if we have any scrolling going on
	if flags&(ImGuiTableFlags_ScrollX|ImGuiTableFlags_ScrollY) != 0 {
		flags &= ^(ImGuiTableFlags_NoHostExtendX | ImGuiTableFlags_NoHostExtendY)
	}

	// Adjust flags: NoBordersInBodyUntilResize takes priority over NoBordersInBody
	if flags&ImGuiTableFlags_NoBordersInBodyUntilResize != 0 {
		flags &= ^ImGuiTableFlags_NoBordersInBody
	}

	// Adjust flags: disable saved settings if there's nothing to save
	if (flags & (ImGuiTableFlags_Resizable | ImGuiTableFlags_Hideable | ImGuiTableFlags_Reorderable | ImGuiTableFlags_Sortable)) == 0 {
		flags |= ImGuiTableFlags_NoSavedSettings
	}

	// Inherit _NoSavedSettings from top-level window (child windows always have _NoSavedSettings set)
	if outer_window.RootWindow.Flags&ImGuiWindowFlags_NoSavedSettings != 0 {
		flags |= ImGuiTableFlags_NoSavedSettings
	}

	return flags
}

// Tables: Candidates for public API
func TableOpenContextMenu(column_n int /*= -1*/) {
	var g = GImGui
	var table = g.CurrentTable
	if column_n == -1 && table.CurrentColumn != -1 { // When called within a column automatically use this one (for consistency)
		column_n = table.CurrentColumn
	}
	if column_n == table.ColumnsCount { // To facilitate using with TableGetHoveredColumn()
		column_n = -1
	}
	IM_ASSERT(column_n >= -1 && column_n < table.ColumnsCount)
	if table.Flags&(ImGuiTableFlags_Resizable|ImGuiTableFlags_Reorderable|ImGuiTableFlags_Hideable) != 0 {
		table.IsContextPopupOpen = true
		table.ContextPopupColumn = (ImGuiTableColumnIdx)(column_n)
		table.InstanceInteracted = table.InstanceCurrent
		var context_menu_id = ImHashStr("##ContextMenu", 0, table.ID)
		OpenPopupEx(context_menu_id, ImGuiPopupFlags_None)
	}
}

// 'width' = inner column width, without padding
func TableSetColumnWidth(column_n int, width float) {
	var g = GImGui
	var table = g.CurrentTable
	IM_ASSERT(table != nil && !table.IsLayoutLocked)
	IM_ASSERT(column_n >= 0 && column_n < table.ColumnsCount)
	var column_0 = &table.Columns[column_n]
	var column_0_width = width

	// Apply constraints early
	// Compare both requested and actual given width to avoid overwriting requested width when column is stuck (minimum size, bounded)
	IM_ASSERT(table.MinColumnWidth > 0.0)
	var min_width = table.MinColumnWidth
	var max_width = ImMax(min_width, TableGetMaxColumnWidth(table, column_n))
	column_0_width = ImClamp(column_0_width, min_width, max_width)
	if column_0.WidthGiven == column_0_width || column_0.WidthRequest == column_0_width {
		return
	}

	//IMGUI_DEBUG_LOG("TableSetColumnWidth(%d, %.1f.%.1f)\n", column_0_idx, column_0.WidthGiven, column_0_width);
	var column_1 *ImGuiTableColumn
	if column_0.NextEnabledColumn != -1 {
		column_1 = &table.Columns[column_0.NextEnabledColumn]
	}

	// In this surprisingly not simple because of how we support mixing Fixed and multiple Stretch columns.
	// - All fixed: easy.
	// - All stretch: easy.
	// - One or more fixed + one stretch: easy.
	// - One or more fixed + more than one stretch: tricky.
	// Qt when manual resize is enabled only support a single _trailing_ stretch column.

	// When forwarding resize from Wn| to Fn+1| we need to be considerate of the _NoResize flag on Fn+1.
	// FIXME-TABLE: Find a way to rewrite all of this so interactions feel more consistent for the user.
	// Scenarios:
	// - F1 F2 F3  resize from F1| or F2|   -. ok: alter .WidthRequested of Fixed column. Subsequent columns will be offset.
	// - F1 F2 F3  resize from F3|          -. ok: alter .WidthRequested of Fixed column. If active, ScrollX extent can be altered.
	// - F1 F2 W3  resize from F1| or F2|   -. ok: alter .WidthRequested of Fixed column. If active, ScrollX extent can be altered, but it doesn't make much sense as the Stretch column will always be minimal size.
	// - F1 F2 W3  resize from W3|          -. ok: no-op (disabled by Resize Rule 1)
	// - W1 W2 W3  resize from W1| or W2|   -. ok
	// - W1 W2 W3  resize from W3|          -. ok: no-op (disabled by Resize Rule 1)
	// - W1 F2 F3  resize from F3|          -. ok: no-op (disabled by Resize Rule 1)
	// - W1 F2     resize from F2|          -. ok: no-op (disabled by Resize Rule 1)
	// - W1 W2 F3  resize from W1| or W2|   -. ok
	// - W1 F2 W3  resize from W1| or F2|   -. ok
	// - F1 W2 F3  resize from W2|          -. ok
	// - F1 W3 F2  resize from W3|          -. ok
	// - W1 F2 F3  resize from W1|          -. ok: equivalent to resizing |F2. F3 will not move.
	// - W1 F2 F3  resize from F2|          -. ok
	// All resizes from a Wx columns are locking other columns.

	// Possible improvements:
	// - W1 W2 W3  resize W1|               -. to not be stuck, both W2 and W3 would stretch down. Seems possible to fix. Would be most beneficial to simplify resize of all-weighted columns.
	// - W3 F1 F2  resize W3|               -. to not be stuck past F1|, both F1 and F2 would need to stretch down, which would be lossy or ambiguous. Seems hard to fix.

	// [Resize Rule 1] Can't resize from right of right-most visible column if there is any Stretch column. Implemented in TableUpdateLayout().

	// If we have all Fixed columns OR resizing a Fixed column that doesn't come after a Stretch one, we can do an offsetting resize.
	// This is the preferred resize path
	if column_0.Flags&ImGuiTableColumnFlags_WidthFixed != 0 {
		if column_1 == nil || table.LeftMostStretchedColumn == -1 || table.Columns[table.LeftMostStretchedColumn].DisplayOrder >= column_0.DisplayOrder {
			column_0.WidthRequest = column_0_width
			table.IsSettingsDirty = true
			return
		}
	}

	// We can also use previous column if there's no next one (this is used when doing an auto-fit on the right-most stretch column)
	if column_1 == nil {
		if column_0.PrevEnabledColumn != -1 {
			column_1 = &table.Columns[column_0.PrevEnabledColumn]
		}
	}
	if column_1 == nil {
		return
	}

	// Resizing from right-side of a Stretch column before a Fixed column forward sizing to left-side of fixed column.
	// (old_a + old_b == new_a + new_b) -. (new_a == old_a + old_b - new_b)
	var column_1_width = ImMax(column_1.WidthRequest-(column_0_width-column_0.WidthRequest), min_width)
	column_0_width = column_0.WidthRequest + column_1.WidthRequest - column_1_width
	IM_ASSERT(column_0_width > 0.0 && column_1_width > 0.0)
	column_0.WidthRequest = column_0_width
	column_1.WidthRequest = column_1_width
	if (column_0.Flags|column_1.Flags)&ImGuiTableColumnFlags_WidthStretch != 0 {
		TableUpdateColumnsWeightFromWidth(table)
	}
	table.IsSettingsDirty = true
}

// Note that the NoSortAscending/NoSortDescending flags are processed in TableSortSpecsSanitize(), and they may change/revert
// the value of SortDirection. We could technically also do it here but it would be unnecessary and duplicate code.
func TableSetColumnSortDirection(column_n int, sort_direction ImGuiSortDirection, append_to_sort_specs bool) {
	var g = GImGui
	var table = g.CurrentTable

	if (table.Flags & ImGuiTableFlags_SortMulti) == 0 {
		append_to_sort_specs = false
	}
	if table.Flags&ImGuiTableFlags_SortTristate == 0 {
		IM_ASSERT(sort_direction != ImGuiSortDirection_None)
	}

	var sort_order_max ImGuiTableColumnIdx = 0
	if append_to_sort_specs {
		for other_column_n := int(0); other_column_n < table.ColumnsCount; other_column_n++ {
			sort_order_max = int8(ImMaxInt(int(sort_order_max), int(table.Columns[other_column_n].SortOrder)))
		}
	}

	var column = &table.Columns[column_n]
	column.SortDirection = (ImGuiSortDirection)(sort_direction)
	if column.SortDirection == ImGuiSortDirection_None {
		column.SortOrder = -1
	} else if column.SortOrder == -1 || !append_to_sort_specs {
		if append_to_sort_specs {
			column.SortOrder = sort_order_max + 1
		} else {
			column.SortOrder = 0
		}
	}

	for other_column_n := int(0); other_column_n < table.ColumnsCount; other_column_n++ {
		var other_column = &table.Columns[other_column_n]
		if other_column != column && !append_to_sort_specs {
			other_column.SortOrder = -1
		}
		TableFixColumnSortDirection(table, other_column)
	}
	table.IsSettingsDirty = true
	table.IsSortSpecsDirty = true
}

// May use (TableGetColumnFlags() & ImGuiTableColumnFlags_IsHovered) instead. Return hovered column. return -1 when table is not hovered. return columns_count if the unused space at the right of visible columns is hovered.
// Return -1 when table is not hovered. return columns_count if the unused space at the right of visible columns is hovered.
func TableGetHoveredColumn() int {
	var g = GImGui
	var table = g.CurrentTable
	if table == nil {
		return -1
	}
	return (int)(table.HoveredColumnBody)
}

func TableGetHeaderRowHeight() float {
	// Caring for a minor edge case:
	// Calculate row height, for the unlikely case that some labels may be taller than others.
	// If we didn't do that, uneven header height would highlight but smaller one before the tallest wouldn't catch input for all height.
	// In your custom header row you may omit this all together and just call TableNextRow() without a height...
	var row_height = GetTextLineHeight()
	var columns_count = TableGetColumnCount()
	for column_n := int(0); column_n < columns_count; column_n++ {
		var flags = TableGetColumnFlags(column_n)
		if (flags&ImGuiTableColumnFlags_IsEnabled) != 0 && (flags&ImGuiTableColumnFlags_NoHeaderLabel) == 0 {
			row_height = ImMax(row_height, CalcTextSize(TableGetColumnName(column_n), true, -1).y)
		}
	}
	row_height += GetStyle().CellPadding.y * 2.0
	return row_height
}

// Bg2 is used by Selectable (and possibly other widgets) to render to the background.
// Unlike our Bg0/1 channel which we uses for RowBg/CellBg/Borders and where we guarantee all shapes to be CPU-clipped, the Bg2 channel being widgets-facing will rely on regular ClipRect.
func TablePushBackgroundChannel() {
	var g = GImGui
	var window = g.CurrentWindow
	var table = g.CurrentTable

	// Optimization: avoid SetCurrentChannel() + PushClipRect()
	table.HostBackupInnerClipRect = window.ClipRect
	SetWindowClipRectBeforeSetChannel(window, &table.Bg2ClipRectForDrawCmd)
	table.DrawSplitter.SetCurrentChannel(window.DrawList, int(table.Bg2DrawChannelCurrent))
}

func TablePopBackgroundChannel() {
	var g = GImGui
	var window = g.CurrentWindow
	var table = g.CurrentTable
	var column = &table.Columns[table.CurrentColumn]

	// Optimization: avoid PopClipRect() + SetCurrentChannel()
	SetWindowClipRectBeforeSetChannel(window, &table.HostBackupInnerClipRect)
	// Skip if draw channel is dummy channel (255 represents -1 when cast to uint8)
	if column.DrawChannelCurrent != ImGuiTableDrawChannelIdx(255) {
		table.DrawSplitter.SetCurrentChannel(window.DrawList, int(column.DrawChannelCurrent))
	}
}

// Tables: Internals
func GetCurrentTable() *ImGuiTable { var g = GImGui; return g.CurrentTable }

func TableFindByID(id ImGuiID) *ImGuiTable {
	var g = GImGui
	return g.Tables[id]
}

func BeginTableEx(name string, id ImGuiID, columns_count int, flags ImGuiTableFlags, outer_size *ImVec2, inner_width float) bool {
	var g = GImGui
	var outer_window = GetCurrentWindow()
	if outer_window.SkipItems { // Consistent with other tables + beneficial side effect that assert on miscalling EndTable() will be more visible.
		return false
	}

	// Sanity checks
	IM_ASSERT_USER_ERROR(columns_count > 0 && columns_count <= IMGUI_TABLE_MAX_COLUMNS, "Only 1..64 columns allowed!")
	if flags&ImGuiTableFlags_ScrollX != 0 {
		IM_ASSERT(inner_width >= 0.0)
	}

	// If an outer size is specified ahead we will be able to early out when not visible. Exact clipping rules may evolve.
	var use_child_window = (flags & (ImGuiTableFlags_ScrollX | ImGuiTableFlags_ScrollY)) != 0
	var avail_size = GetContentRegionAvail()

	var h float
	if use_child_window {
		h = ImMax(avail_size.y, 1.0)
	}

	var actual_outer_size = CalcItemSize(*outer_size, ImMax(avail_size.x, 1.0), h)
	var outer_rect = ImRect{outer_window.DC.CursorPos, outer_window.DC.CursorPos.Add(actual_outer_size)}
	if use_child_window && IsClippedEx(&outer_rect, 0, false) {
		ItemSizeRect(&outer_rect, 0)
		return false
	}

	// Acquire storage for the table
	var table = g.Tables[id]
	if table == nil {
		table = &ImGuiTable{}
		if g.Tables == nil {
			g.Tables = make(map[uint32]*ImGuiTable)
		}
		g.Tables[id] = table
	}
	var instance_no int
	if table.LastFrameActive == g.FrameCount {
		instance_no = int(table.InstanceCurrent + 1)
	}
	var instance_id = id + uint(instance_no)
	var table_last_flags = table.Flags
	if instance_no > 0 {
		IM_ASSERT_USER_ERROR(table.ColumnsCount == columns_count, "BeginTable(): Cannot change columns count mid-frame while preserving same ID")
	}

	// Acquire temporary buffers
	var table_idx = int(id)
	g.CurrentTableStackIdx++
	if g.CurrentTableStackIdx+1 > int(len(g.TablesTempDataStack)) {
		g.TablesTempDataStack = append(g.TablesTempDataStack, NewImGuiTableTempData())
	}
	var temp_data = &g.TablesTempDataStack[g.CurrentTableStackIdx]
	table.TempData = &g.TablesTempDataStack[g.CurrentTableStackIdx]
	temp_data.TableIndex = table_idx
	table.DrawSplitter = &table.TempData.DrawSplitter
	table.DrawSplitter.Clear()

	// Fix flags
	table.IsDefaultSizingPolicy = (flags & ImGuiTableFlags_SizingMask_) == 0
	flags = TableFixFlags(flags, outer_window)

	// Initialize
	table.ID = id
	table.Flags = flags
	table.InstanceCurrent = (ImS16)(instance_no)
	table.LastFrameActive = g.FrameCount
	table.OuterWindow = outer_window
	table.InnerWindow = outer_window
	table.ColumnsCount = columns_count
	table.IsLayoutLocked = false
	table.InnerWidth = inner_width
	temp_data.UserOuterSize = *outer_size

	// When not using a child window, WorkRect.Max will grow as we append contents.
	if use_child_window {
		// Ensure no vertical scrollbar appears if we only want horizontal one, to make flag consistent
		// (we have no other way to disable vertical scrollbar of a window while keeping the horizontal one showing)
		var override_content_size = ImVec2{FLT_MAX, FLT_MAX}
		if (flags&ImGuiTableFlags_ScrollX) != 0 && (flags&ImGuiTableFlags_ScrollY) == 0 {
			override_content_size.y = FLT_MIN
		}

		// Ensure specified width (when not specified, Stretched columns will act as if the width == OuterWidth and
		// never lead to any scrolling). We don't handle inner_width < 0.0f, we could potentially use it to right-align
		// based on the right side of the child window work rect, which would require knowing ahead if we are going to
		// have decoration taking horizontal spaces (typically a vertical scrollbar).
		if (flags&ImGuiTableFlags_ScrollX) != 0 && inner_width > 0.0 {
			override_content_size.x = inner_width
		}

		if override_content_size.x != FLT_MAX || override_content_size.y != FLT_MAX {
			var x, y float
			if override_content_size.x != FLT_MAX {
				x = override_content_size.x
			}
			if override_content_size.y != FLT_MAX {
				y = override_content_size.y
			}
			SetNextWindowContentSize(ImVec2{x, y})
		}

		// Reset scroll if we are reactivating it
		if (table_last_flags & (ImGuiTableFlags_ScrollX | ImGuiTableFlags_ScrollY)) == 0 {
			SetNextWindowScroll(&ImVec2{0.0, 0.0})
		}

		// Create scrolling region (without border and zero window padding)
		var child_flags = ImGuiWindowFlags_None
		if (flags & ImGuiTableFlags_ScrollX) != 0 {
			child_flags = ImGuiWindowFlags_HorizontalScrollbar
		}
		size := outer_rect.GetSize()
		BeginChildEx(name, instance_id, &size, false, child_flags)
		table.InnerWindow = g.CurrentWindow
		table.WorkRect = table.InnerWindow.WorkRect
		table.OuterRect = table.InnerWindow.Rect()
		table.InnerRect = table.InnerWindow.InnerRect
		IM_ASSERT(table.InnerWindow.WindowPadding.x == 0.0 && table.InnerWindow.WindowPadding.y == 0.0 && table.InnerWindow.WindowBorderSize == 0.0)
	} else {
		// For non-scrolling tables, WorkRect == OuterRect == InnerRect.
		// But at this point we do NOT have a correct value for .Max.y (unless a height has been explicitly passed in). It will only be updated in EndTable().
		table.WorkRect = outer_rect
		table.OuterRect = outer_rect
		table.InnerRect = outer_rect
	}

	// Push a standardized ID for both child-using and not-child-using tables
	PushOverrideID(instance_id)

	// Backup a copy of host window members we will modify
	var inner_window = table.InnerWindow
	table.HostIndentX = inner_window.DC.Indent.x
	table.HostClipRect = inner_window.ClipRect
	table.HostSkipItems = inner_window.SkipItems
	temp_data.HostBackupWorkRect = inner_window.WorkRect
	temp_data.HostBackupParentWorkRect = inner_window.ParentWorkRect
	temp_data.HostBackupColumnsOffset = outer_window.DC.ColumnsOffset
	temp_data.HostBackupPrevLineSize = inner_window.DC.PrevLineSize
	temp_data.HostBackupCurrLineSize = inner_window.DC.CurrLineSize
	temp_data.HostBackupCursorMaxPos = inner_window.DC.CursorMaxPos
	temp_data.HostBackupItemWidth = outer_window.DC.ItemWidth
	temp_data.HostBackupItemWidthStackSize = int(len(outer_window.DC.ItemWidthStack))
	inner_window.DC.PrevLineSize = ImVec2{0.0, 0.0}
	inner_window.DC.CurrLineSize = ImVec2{0.0, 0.0}

	// Padding and Spacing
	// - None               ........Content..... Pad .....Content........
	// - PadOuter           | Pad ..Content..... Pad .....Content.. Pad |
	// - PadInner           ........Content.. Pad | Pad ..Content........
	// - PadOuter+PadInner  | Pad ..Content.. Pad | Pad ..Content.. Pad |
	var pad_outer_x = false
	if (flags & ImGuiTableFlags_NoPadOuterX) == 0 {
		if flags&ImGuiTableFlags_PadOuterX != 0 {
			pad_outer_x = true
		} else {
			pad_outer_x = (flags & ImGuiTableFlags_BordersOuterV) != 0
		}
	}
	var pad_inner_x = false
	if (flags & ImGuiTableFlags_NoPadInnerX) == 0 {
		pad_inner_x = true
	}
	var inner_spacing_for_border float = 0.0
	if flags&ImGuiTableFlags_BordersInnerV != 0 {
		inner_spacing_for_border = TABLE_BORDER_SIZE
	}
	var inner_spacing_explicit float
	if pad_inner_x && (flags&ImGuiTableFlags_BordersInnerV) == 0 {
		inner_spacing_explicit = g.Style.CellPadding.x
	}
	var inner_padding_explicit float
	if pad_inner_x && (flags&ImGuiTableFlags_BordersInnerV) != 0 {
		inner_padding_explicit = g.Style.CellPadding.x
	}
	table.CellSpacingX1 = inner_spacing_explicit + inner_spacing_for_border
	table.CellSpacingX2 = inner_spacing_explicit
	table.CellPaddingX = inner_padding_explicit
	table.CellPaddingY = g.Style.CellPadding.y

	var outer_padding_for_border float
	if flags&ImGuiTableFlags_BordersOuterV != 0 {
		outer_padding_for_border = TABLE_BORDER_SIZE
	}
	var outer_padding_explicit float
	if pad_outer_x {
		outer_padding_explicit = g.Style.CellPadding.x
	}
	table.OuterPaddingX = (outer_padding_for_border + outer_padding_explicit) - table.CellPaddingX

	table.CurrentColumn = -1
	table.CurrentRow = -1
	table.RowBgColorCounter = 0
	table.LastRowFlags = ImGuiTableRowFlags_None
	table.InnerClipRect = inner_window.ClipRect
	if inner_window == outer_window {
		table.InnerClipRect = table.WorkRect
	}

	table.InnerClipRect.ClipWith(table.WorkRect) // We need this to honor inner_width
	table.InnerClipRect.ClipWithFull(table.HostClipRect)
	table.InnerClipRect.Max.y = inner_window.ClipRect.Max.y
	if (flags & ImGuiTableFlags_NoHostExtendY) != 0 {
		table.InnerClipRect.Max.y = ImMin(table.InnerClipRect.Max.y, inner_window.WorkRect.Max.y)
	}

	table.RowPosY1 = table.WorkRect.Min.y
	table.RowPosY2 = table.WorkRect.Min.y // This is needed somehow

	table.RowTextBaseline = 0.0 // This will be cleared again by TableBeginRow()

	table.FreezeRowsRequest = 0
	table.FreezeRowsCount = 0 // This will be setup by TableSetupScrollFreeze(), if any

	table.FreezeColumnsRequest = 0
	table.FreezeColumnsCount = 0
	table.IsUnfrozenRows = true
	table.DeclColumnsCount = 0

	// Using opaque colors facilitate overlapping elements of the grid
	table.BorderColorStrong = GetColorU32FromID(ImGuiCol_TableBorderStrong, 1)
	table.BorderColorLight = GetColorU32FromID(ImGuiCol_TableBorderLight, 1)

	// Make table current
	g.CurrentTable = table
	outer_window.DC.CurrentTableIdx = table_idx
	if inner_window != outer_window { // So EndChild() within the inner window can restore the table properly.
		inner_window.DC.CurrentTableIdx = table_idx
	}

	if (table_last_flags&ImGuiTableFlags_Reorderable) != 0 && (flags&ImGuiTableFlags_Reorderable) == 0 {
		table.IsResetDisplayOrderRequest = true
	}

	// Mark as used
	if g.TablesLastTimeActive == nil {
		g.TablesLastTimeActive = make(map[int32]float32)
	}
	g.TablesLastTimeActive[table_idx] = (float)(g.Time)
	temp_data.LastTimeActive = (float)(g.Time)
	table.MemoryCompacted = false

	// Setup memory buffer (clear data if columns count changed)
	var old_columns_to_preserve []ImGuiTableColumn = nil
	var old_columns_raw_data any = nil
	var old_columns_count = int(len(table.Columns))
	if old_columns_count != 0 && old_columns_count != columns_count {
		// Attempt to preserve width on column count change (#4046)
		old_columns_to_preserve = table.Columns
		old_columns_raw_data = table.RawData
		table.RawData = nil
	}
	if table.RawData == nil {
		TableBeginInitMemory(table, columns_count)
		table.IsInitializing = true
		table.IsSettingsRequestLoad = true
	}
	if table.IsResetAllRequest {
		TableResetSettings(table)
	}
	if table.IsInitializing {
		// Initialize
		table.SettingsOffset = -1
		table.IsSortSpecsDirty = true
		table.InstanceInteracted = -1
		table.ContextPopupColumn = -1
		table.ReorderColumn = -1
		table.ResizedColumn = -1
		table.LastResizedColumn = -1
		table.AutoFitSingleColumn = -1
		table.HoveredColumnBody = -1
		table.HoveredColumnBorder = -1
		for n := range table.Columns {
			var column = &table.Columns[n]
			if old_columns_to_preserve != nil && int(n) < old_columns_count {
				// FIXME: We don't attempt to preserve column order in this path.
				*column = old_columns_to_preserve[n]
			} else {
				var width_auto = column.WidthAuto
				*column = NewImGuiTableColumn()
				column.WidthAuto = width_auto
				column.IsPreserveWidthAuto = true // Preserve WidthAuto when reinitializing a live table: not technically necessary but remove a visible flicker
				column.IsEnabled = true
				column.IsUserEnabled = true
				column.IsUserEnabledNextFrame = true
			}
			column.DisplayOrder = (ImGuiTableColumnIdx)(n)
			table.DisplayOrderToIndex[n] = (ImGuiTableColumnIdx)(n)
		}
	}
	if old_columns_raw_data != nil {
		old_columns_raw_data = nil
	}

	// Load settings
	if table.IsSettingsRequestLoad {
		TableLoadSettings(table)
	}

	// Handle DPI/font resize
	// This is designed to facilitate DPI changes with the assumption that e.g. style.CellPadding has been scaled as well.
	// It will also react to changing fonts with mixed results. It doesn't need to be perfect but merely provide a decent transition.
	// FIXME-DPI: Provide consistent standards for reference size. Perhaps using g.CurrentDpiScale would be more self explanatory.
	// This is will lead us to non-rounded WidthRequest in columns, which should work but is a poorly tested path.
	var new_ref_scale_unit = g.FontSize // g.Font.GetCharAdvance('A') ?
	if table.RefScale != 0.0 && table.RefScale != new_ref_scale_unit {
		var scale_factor = new_ref_scale_unit / table.RefScale
		//IMGUI_DEBUG_LOG("[table] %08X RefScaleUnit %.3f . %.3f, scaling width by %.3f\n", table.ID, table.RefScaleUnit, new_ref_scale_unit, scale_factor);
		for n := int(0); n < columns_count; n++ {
			table.Columns[n].WidthRequest = table.Columns[n].WidthRequest * scale_factor
		}
	}
	table.RefScale = new_ref_scale_unit

	// Disable output until user calls TableNextRow() or TableNextColumn() leading to the TableUpdateLayout() call..
	// This is not strictly necessary but will reduce cases were "out of table" output will be misleading to the user.
	// Because we cannot safely assert in EndTable() when no rows have been created, this seems like our best option.
	inner_window.SkipItems = true

	// Clear names
	// At this point the .NameOffset field of each column will be invalid until TableUpdateLayout() or the first call to TableSetupColumn()
	if int(len(table.ColumnsNames)) > 0 {
		table.ColumnsNames = table.ColumnsNames[:0]
	}

	// Apply queued resizing/reordering/hiding requests
	TableBeginApplyRequests(table)

	return true
}

// For reference, the average total _allocation count_ for a table is:
// + 0 (for ImGuiTable instance, we are pooling allocations in g.Tables)
// + 1 (for table.RawData allocated below)
// + 1 (for table.ColumnsNames, if names are used)
// + 1 (for table.Splitter._Channels)
// + 2 * active_channels_count (for ImDrawCmd and ImDrawIdx buffers inside channels)
// Where active_channels_count is variable but often == columns_count or columns_count + 1, see TableSetupDrawChannels() for details.
// Unused channels don't perform their +2 allocations.
func TableBeginInitMemory(e *ImGuiTable, columns_count int) {
	// noop, will be handled by span helpers
}

// Apply queued resizing/reordering/hiding requests
func TableBeginApplyRequests(table *ImGuiTable) {
	// Handle resizing request
	// (We process this at the first TableBegin of the frame)
	// FIXME-TABLE: Contains columns if our work area doesn't allow for scrolling?
	if table.InstanceCurrent == 0 {
		if table.ResizedColumn != -1 && table.ResizedColumnNextWidth != FLT_MAX {
			TableSetColumnWidth(int(table.ResizedColumn), table.ResizedColumnNextWidth)
		}
		table.LastResizedColumn = table.ResizedColumn
		table.ResizedColumnNextWidth = FLT_MAX
		table.ResizedColumn = -1

		// Process auto-fit for single column, which is a special case for stretch columns and fixed columns with FixedSame policy.
		// FIXME-TABLE: Would be nice to redistribute available stretch space accordingly to other weights, instead of giving it all to siblings.
		if table.AutoFitSingleColumn != -1 {
			TableSetColumnWidth(int(table.AutoFitSingleColumn), table.Columns[table.AutoFitSingleColumn].WidthAuto)
			table.AutoFitSingleColumn = -1
		}
	}

	// Handle reordering request
	// Note: we don't clear ReorderColumn after handling the request.
	if table.InstanceCurrent == 0 {
		if table.HeldHeaderColumn == -1 && table.ReorderColumn != -1 {
			table.ReorderColumn = -1
		}
		table.HeldHeaderColumn = -1
		if table.ReorderColumn != -1 && table.ReorderColumnDir != 0 {
			// We need to handle reordering across hidden columns.
			// In the configuration below, moving C to the right of E will lead to:
			//    ... C [D] E  --.  ... [D] E  C   (Column name/index)
			//    ... 2  3  4        ...  2  3  4   (Display order)
			var reorder_dir = int(table.ReorderColumnDir)
			IM_ASSERT(reorder_dir == -1 || reorder_dir == +1)
			IM_ASSERT(table.Flags&ImGuiTableFlags_Reorderable != 0)
			var src_column = &table.Columns[table.ReorderColumn]
			var dst_column *ImGuiTableColumn
			if reorder_dir == -1 {
				dst_column = &table.Columns[src_column.PrevEnabledColumn]
			} else {
				dst_column = &table.Columns[src_column.NextEnabledColumn]
			}
			var src_order = int(src_column.DisplayOrder)
			var dst_order = int(dst_column.DisplayOrder)
			src_column.DisplayOrder = (ImGuiTableColumnIdx)(dst_order)
			for order_n := src_order + reorder_dir; order_n != dst_order+reorder_dir; order_n += reorder_dir {
				table.Columns[table.DisplayOrderToIndex[order_n]].DisplayOrder -= (ImGuiTableColumnIdx)(reorder_dir)
			}
			IM_ASSERT(int(dst_column.DisplayOrder) == dst_order-reorder_dir)

			// Display order is stored in both columns.IndexDisplayOrder and table.DisplayOrder[],
			// rebuild the later from the former.
			for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
				table.DisplayOrderToIndex[table.Columns[column_n].DisplayOrder] = (ImGuiTableColumnIdx)(column_n)
			}
			table.ReorderColumnDir = 0
			table.IsSettingsDirty = true
		}
	}

	// Handle display order reset request
	if table.IsResetDisplayOrderRequest {
		for n := int(0); n < table.ColumnsCount; n++ {
			table.DisplayOrderToIndex[n] = (ImGuiTableColumnIdx)(n)
			table.Columns[n].DisplayOrder = (ImGuiTableColumnIdx)(n)
		}
		table.IsResetDisplayOrderRequest = false
		table.IsSettingsDirty = true
	}
}

// Adjust flags: default width mode + stretch columns are not allowed when auto extending
func TableSetupColumnFlags(table *ImGuiTable, column *ImGuiTableColumn, flags_in ImGuiTableColumnFlags) {
	var flags = flags_in

	// Sizing Policy
	if (flags & ImGuiTableColumnFlags_WidthMask_) == 0 {
		var table_sizing_policy = (table.Flags & ImGuiTableFlags_SizingMask_)
		if table_sizing_policy == ImGuiTableFlags_SizingFixedFit || table_sizing_policy == ImGuiTableFlags_SizingFixedSame {
			flags |= ImGuiTableColumnFlags_WidthFixed
		} else {
			flags |= ImGuiTableColumnFlags_WidthStretch
		}
	} else {
		IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiTableColumnFlags_WidthMask_))) // Check that only 1 of each set is used.
	}

	// Resize
	if (table.Flags & ImGuiTableFlags_Resizable) == 0 {
		flags |= ImGuiTableColumnFlags_NoResize
	}

	// Sorting
	if (flags&ImGuiTableColumnFlags_NoSortAscending) != 0 && (flags&ImGuiTableColumnFlags_NoSortDescending) != 0 {
		flags |= ImGuiTableColumnFlags_NoSort
	}

	// Indentation
	if (flags & ImGuiTableColumnFlags_IndentMask_) == 0 {
		if column == &table.Columns[0] {
			flags |= ImGuiTableColumnFlags_IndentEnable
		} else {
			flags |= ImGuiTableColumnFlags_IndentDisable
		}
	}

	// Alignment
	//if ((flags & ImGuiTableColumnFlags_AlignMask_) == 0)
	//    flags |= ImGuiTableColumnFlags_AlignCenter;
	//IM_ASSERT(ImIsPowerOfTwo(flags & ImGuiTableColumnFlags_AlignMask_)); // Check that only 1 of each set is used.

	// Preserve status flags
	column.Flags = flags | (column.Flags & ImGuiTableColumnFlags_StatusMask_)

	// Build an ordered list of available sort directions
	column.SortDirectionsAvailCount = 0
	column.SortDirectionsAvailMask = 0
	column.SortDirectionsAvailList = 0
	if (table.Flags & ImGuiTableFlags_Sortable) != 0 {
		var count, mask, list int
		if (flags&ImGuiTableColumnFlags_PreferSortAscending) != 0 && (flags&ImGuiTableColumnFlags_NoSortAscending) == 0 {
			mask |= 1 << ImGuiSortDirection_Ascending
			list |= int(ImGuiSortDirection_Ascending) << (count << 1)
			count++
		}
		if (flags&ImGuiTableColumnFlags_PreferSortDescending) != 0 && (flags&ImGuiTableColumnFlags_NoSortDescending) == 0 {
			mask |= 1 << ImGuiSortDirection_Descending
			list |= int(ImGuiSortDirection_Descending) << (count << 1)
			count++
		}
		if (flags&ImGuiTableColumnFlags_PreferSortAscending) == 0 && (flags&ImGuiTableColumnFlags_NoSortAscending) == 0 {
			mask |= 1 << ImGuiSortDirection_Ascending
			list |= int(ImGuiSortDirection_Ascending) << (count << 1)
			count++
		}
		if (flags&ImGuiTableColumnFlags_PreferSortDescending) == 0 && (flags&ImGuiTableColumnFlags_NoSortDescending) == 0 {
			mask |= 1 << ImGuiSortDirection_Descending
			list |= int(ImGuiSortDirection_Descending) << (count << 1)
			count++
		}
		if (table.Flags&ImGuiTableFlags_SortTristate) != 0 || count == 0 {
			mask |= 1 << ImGuiSortDirection_None
			count++
		}
		column.SortDirectionsAvailList = (ImU8)(list)
		column.SortDirectionsAvailMask = (ImU8)(mask)
		column.SortDirectionsAvailCount = (ImU8)(count)
		TableFixColumnSortDirection(table, column)
	}
}

// Allocate draw channels. Called by TableUpdateLayout()
//   - We allocate them following storage order instead of display order so reordering columns won't needlessly
//     increase overall dormant memory cost.
//   - We isolate headers draw commands in their own channels instead of just altering clip rects.
//     This is in order to facilitate merging of draw commands.
//   - After crossing FreezeRowsCount, all columns see their current draw channel changed to a second set of channels.
//   - We only use the dummy draw channel so we can push a null clipping rectangle into it without affecting other
//     channels, while simplifying per-row/per-cell overhead. It will be empty and discarded when merged.
//   - We allocate 1 or 2 background draw channels. This is because we know TablePushBackgroundChannel() is only used for
//     horizontal spanning. If we allowed vertical spanning we'd need one background draw channel per merge group (1-4).
//
// Draw channel allocation (before merging):
// - NoClip                       -. 2+D+1 channels: bg0/1 + bg2 + foreground (same clip rect == always 1 draw call)
// - Clip                         -. 2+D+N channels
// - FreezeRows                   -. 2+D+N*2 (unless scrolling value is zero)
// - FreezeRows || FreezeColunns  -. 3+D+N*2 (unless scrolling value is zero)
// Where D is 1 if any column is clipped or hidden (dummy channel) otherwise 0.
func TableSetupDrawChannels(table *ImGuiTable) {
	var freeze_row_multiplier int = 1
	if table.FreezeRowsCount > 0 {
		freeze_row_multiplier = 2
	}
	var channels_for_row = int(table.ColumnsEnabledCount)
	if table.Flags&ImGuiTableFlags_NoClip != 0 {
		channels_for_row = 1
	}
	var channels_for_bg = 1 + 1*freeze_row_multiplier
	var channels_for_dummy int
	if int(table.ColumnsEnabledCount) < table.ColumnsCount || table.VisibleMaskByIndex != table.EnabledMaskByIndex {
		channels_for_dummy = 1
	}
	var channels_total = channels_for_bg + (channels_for_row * freeze_row_multiplier) + channels_for_dummy

	table.DrawSplitter.Split(table.InnerWindow.DrawList, channels_total)

	var ch int = -1
	if channels_for_dummy > 0 {
		ch = channels_total - 1
	}

	table.DummyDrawChannel = (ImGuiTableDrawChannelIdx)(ch)
	table.Bg2DrawChannelCurrent = uint8(TABLE_DRAW_CHANNEL_BG2_FROZEN)
	if table.FreezeRowsCount > 0 {
		table.Bg2DrawChannelUnfrozen = uint8(2 + channels_for_row)
	} else {
		table.Bg2DrawChannelUnfrozen = uint8(TABLE_DRAW_CHANNEL_BG2_FROZEN)
	}

	var draw_channel_current int = 2
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]
		if column.IsVisibleX && column.IsVisibleY {
			column.DrawChannelFrozen = (ImGuiTableDrawChannelIdx)(draw_channel_current)

			var ch int
			if table.FreezeRowsCount > 0 {
				ch = channels_for_row + 1
			}

			column.DrawChannelUnfrozen = (ImGuiTableDrawChannelIdx)(draw_channel_current + ch)
			if (table.Flags & ImGuiTableFlags_NoClip) == 0 {
				draw_channel_current++
			}
		} else {
			column.DrawChannelFrozen = table.DummyDrawChannel
			column.DrawChannelUnfrozen = table.DummyDrawChannel
		}
		column.DrawChannelCurrent = column.DrawChannelFrozen
	}

	// Initial draw cmd starts with a BgClipRect that matches the one of its host, to facilitate merge draw commands by default.
	// All our cell highlight are manually clipped with BgClipRect. When unfreezing it will be made smaller to fit scrolling rect.
	// (This technically isn't part of setting up draw channels, but is reasonably related to be done here)
	table.BgClipRect = table.InnerClipRect
	table.Bg0ClipRectForDrawCmd = table.OuterWindow.ClipRect
	table.Bg2ClipRectForDrawCmd = table.HostClipRect
	IM_ASSERT(table.BgClipRect.Min.y <= table.BgClipRect.Max.y)
}

// helper for the span allocator. called when setting an index
// of table.DisplayOrderToIndex
func (table *ImGuiTable) spanDisplayOrderToIndex(order_n int) {
	displayOrderLen := int(len(table.DisplayOrderToIndex))

	// check if that index already exists
	if order_n <= displayOrderLen-1 {
		return
	}

	// add missing items
	for i := int(0); i < (order_n+1)-displayOrderLen; i++ {
		table.DisplayOrderToIndex = append(table.DisplayOrderToIndex, 0)
	}
}

// helper for the span allocator. called when setting an index
// of table.Columns
func (table *ImGuiTable) spanColumns(column_n int) {
	columnsLen := int(len(table.Columns))

	// check if that index already exists
	if int(column_n) <= columnsLen-1 {
		return
	}

	// add missing items
	for i := int(0); i < (int(column_n)+1)-columnsLen; i++ {
		table.Columns = append(table.Columns, NewImGuiTableColumn())
	}
}

// Layout columns for the frame. This is in essence the followup to BeginTable().
// Runs on the first call to TableNextRow(), to give a chance for TableSetupColumn() to be called first.
// FIXME-TABLE: Our width (and therefore our WorkRect) will be minimal in the first frame for _WidthAuto columns.
// Increase feedback side-effect with widgets relying on WorkRect.Max.x... Maybe provide a default distribution for _WidthAuto columns?
func TableUpdateLayout(table *ImGuiTable) {
	var g = GImGui
	IM_ASSERT(!table.IsLayoutLocked)

	var table_sizing_policy = (table.Flags & ImGuiTableFlags_SizingMask_)
	table.IsDefaultDisplayOrder = true
	table.ColumnsEnabledCount = 0
	table.EnabledMaskByIndex = 0x00
	table.EnabledMaskByDisplayOrder = 0x00
	table.LeftMostEnabledColumn = -1
	table.MinColumnWidth = ImMax(1.0, g.Style.FramePadding.x*1.0) // g.Style.ColumnsMinSpacing; // FIXME-TABLE

	// [Part 1] Apply/lock Enabled and Order states. Calculate auto/ideal width for columns. Count fixed/stretch columns.
	// Process columns in their visible orders as we are building the Prev/Next indices.
	var count_fixed int = 0   // Number of columns that have fixed sizing policies
	var count_stretch int = 0 // Number of columns that have stretch sizing policies
	var prev_visible_column_idx int = -1
	var has_auto_fit_request = false
	var has_resizable = false
	var stretch_sum_width_auto float = 0.0
	var fixed_max_width_auto float = 0.0
	for order_n := int(0); order_n < table.ColumnsCount; order_n++ {
		table.spanDisplayOrderToIndex(order_n)
		var column_n = table.DisplayOrderToIndex[order_n]
		if int(column_n) != order_n {
			table.IsDefaultDisplayOrder = false
		}
		table.spanColumns(int(column_n))
		var column = &table.Columns[column_n]

		// Clear column setup if not submitted by user. Currently we make it mandatory to call TableSetupColumn() every frame.
		// It would easily work without but we're not ready to guarantee it since e.g. names need resubmission anyway.
		// We take a slight shortcut but in theory we could be calling TableSetupColumn() here with dummy values, it should yield the same effect.
		if table.DeclColumnsCount <= column_n {
			TableSetupColumnFlags(table, column, ImGuiTableColumnFlags_None)
			column.NameOffset = -1
			column.UserID = 0
			column.InitStretchWeightOrWidth = -1.0
		}

		// Update Enabled state, mark settings and sort specs dirty
		if (table.Flags&ImGuiTableFlags_Hideable) == 0 || (column.Flags&ImGuiTableColumnFlags_NoHide) != 0 {
			column.IsUserEnabledNextFrame = true
		}
		if column.IsUserEnabled != column.IsUserEnabledNextFrame {
			column.IsUserEnabled = column.IsUserEnabledNextFrame
			table.IsSettingsDirty = true
		}
		column.IsEnabled = column.IsUserEnabled && (column.Flags&ImGuiTableColumnFlags_Disabled) == 0

		if column.SortOrder != -1 && !column.IsEnabled {
			table.IsSortSpecsDirty = true
		}
		if column.SortOrder > 0 && (table.Flags&ImGuiTableFlags_SortMulti) == 0 {
			table.IsSortSpecsDirty = true
		}

		// Auto-fit unsized columns
		var start_auto_fit = (column.StretchWeight < 0.0)
		if column.Flags&ImGuiTableColumnFlags_WidthFixed != 0 {
			start_auto_fit = (column.WidthRequest < 0.0)
		}
		if start_auto_fit {
			column.AutoFitQueue = (1 << 3) - 1
			column.CannotSkipItemsQueue = (1 << 3) - 1 // Fit for three frames
		}

		if !column.IsEnabled {
			column.IndexWithinEnabledSet = -1
			continue
		}

		// Mark as enabled and link to previous/next enabled column
		column.PrevEnabledColumn = (ImGuiTableColumnIdx)(prev_visible_column_idx)
		column.NextEnabledColumn = -1
		if prev_visible_column_idx != -1 {
			table.spanColumns(prev_visible_column_idx)
			table.Columns[prev_visible_column_idx].NextEnabledColumn = (ImGuiTableColumnIdx)(column_n)
		} else {
			table.LeftMostEnabledColumn = (ImGuiTableColumnIdx)(column_n)
		}
		column.IndexWithinEnabledSet = table.ColumnsEnabledCount
		table.ColumnsEnabledCount++

		table.EnabledMaskByIndex |= (ImU64)(1 << column_n)

		displayOrderShift := column.DisplayOrder
		if displayOrderShift < 0 {
			displayOrderShift = 0
		}
		table.EnabledMaskByDisplayOrder |= (ImU64)(1 << displayOrderShift)

		prev_visible_column_idx = int(column_n)
		// FIXME (port): figure out this panic
		// IM_ASSERT(column.IndexWithinEnabledSet <= column.DisplayOrder)

		// Calculate ideal/auto column width (that's the width required for all contents to be visible without clipping)
		// Combine width from regular rows + width from headers unless requested not to.
		if !column.IsPreserveWidthAuto {
			column.WidthAuto = TableGetColumnWidthAuto(table, column)
		}

		// Non-resizable columns keep their requested width (apply user value regardless of IsPreserveWidthAuto)
		var column_is_resizable = (column.Flags & ImGuiTableColumnFlags_NoResize) == 0
		if column_is_resizable {
			has_resizable = true
		}
		if (column.Flags&ImGuiTableColumnFlags_WidthFixed) != 0 && column.InitStretchWeightOrWidth > 0.0 && !column_is_resizable {
			column.WidthAuto = column.InitStretchWeightOrWidth
		}

		if column.AutoFitQueue != 0x00 {
			has_auto_fit_request = true
		}
		if column.Flags&ImGuiTableColumnFlags_WidthStretch != 0 {
			stretch_sum_width_auto += column.WidthAuto
			count_stretch++
		} else {
			fixed_max_width_auto = ImMax(fixed_max_width_auto, column.WidthAuto)
			count_fixed++
		}
	}
	if (table.Flags&ImGuiTableFlags_Sortable) != 0 && table.SortSpecsCount == 0 && table.Flags&ImGuiTableFlags_SortTristate == 0 {
		table.IsSortSpecsDirty = true
	}
	table.RightMostEnabledColumn = (ImGuiTableColumnIdx)(prev_visible_column_idx)
	IM_ASSERT(table.LeftMostEnabledColumn >= 0 && table.RightMostEnabledColumn >= 0)

	// [Part 2] Disable child window clipping while fitting columns. This is not strictly necessary but makes it possible
	// to avoid the column fitting having to wait until the first visible frame of the child container (may or not be a good thing).
	// FIXME-TABLE: for always auto-resizing columns may not want to do that all the time.
	if has_auto_fit_request && table.OuterWindow != table.InnerWindow {
		table.InnerWindow.SkipItems = false
	}
	if has_auto_fit_request {
		table.IsSettingsDirty = true
	}

	// [Part 3] Fix column flags and record a few extra information.
	var sum_width_requests float = 0.0  // Sum of all width for fixed and auto-resize columns, excluding width contributed by Stretch columns but including spacing/padding.
	var stretch_sum_weights float = 0.0 // Sum of all weights for stretch columns.
	table.LeftMostStretchedColumn = -1
	table.RightMostStretchedColumn = -1
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		if table.EnabledMaskByIndex&((ImU64)(1<<column_n)) == 0 {
			continue
		}
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]

		var column_is_resizable = (column.Flags & ImGuiTableColumnFlags_NoResize) == 0
		if (column.Flags & ImGuiTableColumnFlags_WidthFixed) != 0 {
			// Apply same widths policy
			var width_auto = column.WidthAuto
			if table_sizing_policy == ImGuiTableFlags_SizingFixedSame && (column.AutoFitQueue != 0x00 || !column_is_resizable) {
				width_auto = fixed_max_width_auto
			}

			// Apply automatic width
			// Latch initial size for fixed columns and update it constantly for auto-resizing column (unless clipped!)
			if column.AutoFitQueue != 0x00 {
				column.WidthRequest = width_auto
			} else if (column.Flags&ImGuiTableColumnFlags_WidthFixed) != 0 && !column_is_resizable && (table.RequestOutputMaskByIndex&((ImU64)(1<<column_n))) != 0 {
				column.WidthRequest = width_auto
			}

			// FIXME-TABLE: Increase minimum size during init frame to avoid biasing auto-fitting widgets
			// (e.g. TextWrapped) too much. Otherwise what tends to happen is that TextWrapped would output a very
			// large height (= first frame scrollbar display very off + clipper would skip lots of items).
			// This is merely making the side-effect less extreme, but doesn't properly fixes it.
			// FIXME: Move this to .WidthGiven to avoid temporary lossyless?
			// FIXME: This break IsPreserveWidthAuto from not flickering if the stored WidthAuto was smaller.
			if column.AutoFitQueue > 0x01 && table.IsInitializing && !column.IsPreserveWidthAuto {
				column.WidthRequest = ImMax(column.WidthRequest, table.MinColumnWidth*4.0) // FIXME-TABLE: Another constant/scale?
			}
			sum_width_requests += column.WidthRequest
		} else {
			// Initialize stretch weight
			if column.AutoFitQueue != 0x00 || column.StretchWeight < 0.0 || !column_is_resizable {
				if column.InitStretchWeightOrWidth > 0.0 {
					column.StretchWeight = column.InitStretchWeightOrWidth
				} else if table_sizing_policy == ImGuiTableFlags_SizingStretchProp {
					column.StretchWeight = (column.WidthAuto / stretch_sum_width_auto) * float(count_stretch)
				} else {
					column.StretchWeight = 1.0
				}
			}

			stretch_sum_weights += column.StretchWeight
			if table.LeftMostStretchedColumn == -1 || table.Columns[table.LeftMostStretchedColumn].DisplayOrder > column.DisplayOrder {
				table.LeftMostStretchedColumn = (ImGuiTableColumnIdx)(column_n)
			}
			if table.RightMostStretchedColumn == -1 || table.Columns[table.RightMostStretchedColumn].DisplayOrder < column.DisplayOrder {
				table.RightMostStretchedColumn = (ImGuiTableColumnIdx)(column_n)
			}
		}
		column.IsPreserveWidthAuto = false
		sum_width_requests += table.CellPaddingX * 2.0
	}
	table.ColumnsEnabledFixedCount = (ImGuiTableColumnIdx)(count_fixed)

	// [Part 4] Apply final widths based on requested widths
	var work_rect = table.WorkRect
	var width_spacings = (table.OuterPaddingX * 2.0) + (table.CellSpacingX1+table.CellSpacingX2)*float(table.ColumnsEnabledCount-1)

	var width_avail float
	if (table.Flags&ImGuiTableFlags_ScrollX) != 0 && table.InnerWidth == 0.0 {
		width_avail = table.InnerClipRect.GetWidth()
	} else {
		width_avail = work_rect.GetWidth()
	}

	var width_avail_for_stretched_columns = width_avail - width_spacings - sum_width_requests
	var width_remaining_for_stretched_columns = width_avail_for_stretched_columns
	table.ColumnsGivenWidth = width_spacings + (table.CellPaddingX*2.0)*float(table.ColumnsEnabledCount)
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		if (table.EnabledMaskByIndex & ((ImU64)(1 << column_n))) == 0 {
			continue
		}
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]

		// Allocate width for stretched/weighted columns (StretchWeight gets converted into WidthRequest)
		if column.Flags&ImGuiTableColumnFlags_WidthStretch != 0 {
			var weight_ratio = column.StretchWeight / stretch_sum_weights
			column.WidthRequest = IM_FLOOR(ImMax(width_avail_for_stretched_columns*weight_ratio, table.MinColumnWidth) + 0.01)
			width_remaining_for_stretched_columns -= column.WidthRequest
		}

		// [Resize Rule 1] The right-most Visible column is not resizable if there is at least one Stretch column
		// See additional comments in TableSetColumnWidth().
		if column.NextEnabledColumn == -1 && table.LeftMostStretchedColumn != -1 {
			column.Flags |= ImGuiTableColumnFlags_NoDirectResize_
		}

		// Assign final width, record width in case we will need to shrink
		column.WidthGiven = ImFloor(ImMax(column.WidthRequest, table.MinColumnWidth))
		table.ColumnsGivenWidth += column.WidthGiven
	}

	// [Part 5] Redistribute stretch remainder width due to rounding (remainder width is < 1.0f * number of Stretch column).
	// Using right-to-left distribution (more likely to match resizing cursor).
	if width_remaining_for_stretched_columns >= 1.0 && (table.Flags&ImGuiTableFlags_PreciseWidths) == 0 {
		for order_n := table.ColumnsCount - 1; stretch_sum_weights > 0.0 && width_remaining_for_stretched_columns >= 1.0 && order_n >= 0; order_n-- {
			if (table.EnabledMaskByDisplayOrder & ((ImU64)(1 << order_n))) != 0 {
				continue
			}
			table.spanDisplayOrderToIndex(order_n)
			table.spanColumns(int(table.DisplayOrderToIndex[order_n]))

			var column = &table.Columns[table.DisplayOrderToIndex[order_n]]
			if column.Flags&ImGuiTableColumnFlags_WidthStretch == 0 {
				continue
			}
			column.WidthRequest += 1.0
			column.WidthGiven += 1.0
			width_remaining_for_stretched_columns -= 1.0
		}
	}

	table.HoveredColumnBody = -1
	table.HoveredColumnBorder = -1
	var mouse_hit_rect = ImRect{ImVec2{table.OuterRect.Min.x, table.OuterRect.Min.y}, ImVec2{table.OuterRect.Max.x, ImMax(table.OuterRect.Max.y, table.OuterRect.Min.y+table.LastOuterHeight)}}
	var is_hovering_table = ItemHoverable(&mouse_hit_rect, 0)

	// [Part 6] Setup final position, offset, skip/clip states and clipping rectangles, detect hovered column
	// Process columns in their visible orders as we are comparing the visible order and adjusting host_clip_rect while looping.
	var visible_n int = 0
	var offset_x_frozen = (table.FreezeColumnsCount > 0)

	var x = work_rect.Min.x
	if table.FreezeColumnsCount > 0 {
		x = table.OuterRect.Min.x
	}

	var offset_x = x + table.OuterPaddingX - table.CellSpacingX1
	var host_clip_rect = table.InnerClipRect
	//host_clip_rect.Max.x += table.CellPaddingX + table.CellSpacingX2;
	table.VisibleMaskByIndex = 0x00
	table.RequestOutputMaskByIndex = 0x00
	for order_n := int(0); order_n < table.ColumnsCount; order_n++ {
		table.spanDisplayOrderToIndex(order_n)
		var column_n = table.DisplayOrderToIndex[order_n]
		table.spanColumns(int(column_n))
		var column = &table.Columns[column_n]

		if table.FreezeRowsCount > 0 || column_n < table.FreezeColumnsCount {
			column.NavLayerCurrent = int8(ImGuiNavLayer_Menu)
		} else {
			column.NavLayerCurrent = int8(ImGuiNavLayer_Main)
		}

		if offset_x_frozen && int(table.FreezeColumnsCount) == visible_n {
			offset_x += work_rect.Min.x - table.OuterRect.Min.x
			offset_x_frozen = false
		}

		// Clear status flags
		column.Flags &= ^ImGuiTableColumnFlags_StatusMask_

		if (table.EnabledMaskByDisplayOrder & ((ImU64)(1 << order_n))) == 0 {
			// Hidden column: clear a few fields and we are done with it for the remainder of the function.
			// We set a zero-width clip rect but set Min.y/Max.y properly to not interfere with the clipper.
			column.MinX = offset_x
			column.MaxX = offset_x
			column.WorkMinX = offset_x
			column.ClipRect.Min.x = offset_x
			column.ClipRect.Max.x = offset_x
			column.WidthGiven = 0.0
			column.ClipRect.Min.y = work_rect.Min.y
			column.ClipRect.Max.y = FLT_MAX
			column.ClipRect.ClipWithFull(host_clip_rect)
			column.IsVisibleX = false
			column.IsVisibleY = false
			column.IsRequestOutput = false
			column.IsSkipItems = true
			column.ItemWidth = 1.0
			continue
		}

		// Detect hovered column
		if is_hovering_table && g.IO.MousePos.x >= column.ClipRect.Min.x && g.IO.MousePos.x < column.ClipRect.Max.x {
			table.HoveredColumnBody = (ImGuiTableColumnIdx)(column_n)
		}

		// Lock start position
		column.MinX = offset_x

		// Lock width based on start position and minimum/maximum width for this position
		var max_width = TableGetMaxColumnWidth(table, int(column_n))
		column.WidthGiven = ImMin(column.WidthGiven, max_width)
		column.WidthGiven = ImMax(column.WidthGiven, ImMin(column.WidthRequest, table.MinColumnWidth))
		column.MaxX = offset_x + column.WidthGiven + table.CellSpacingX1 + table.CellSpacingX2 + table.CellPaddingX*2.0

		// Lock other positions
		// - ClipRect.Min.x: Because merging draw commands doesn't compare min boundaries, we make ClipRect.Min.x match left bounds to be consistent regardless of merging.
		// - ClipRect.Max.x: using WorkMaxX instead of MaxX (aka including padding) makes things more consistent when resizing down, tho slightly detrimental to visibility in very-small column.
		// - ClipRect.Max.x: using MaxX makes it easier for header to receive hover highlight with no discontinuity and display sorting arrow.
		// - FIXME-TABLE: We want equal width columns to have equal (ClipRect.Max.x - WorkMinX) width, which means ClipRect.max.x cannot stray off host_clip_rect.Max.x else right-most column may appear shorter.
		column.WorkMinX = column.MinX + table.CellPaddingX + table.CellSpacingX1
		column.WorkMaxX = column.MaxX - table.CellPaddingX - table.CellSpacingX2 // Expected max
		column.ItemWidth = ImFloor(column.WidthGiven * 0.65)
		column.ClipRect.Min.x = column.MinX
		column.ClipRect.Min.y = work_rect.Min.y
		column.ClipRect.Max.x = column.MaxX //column.WorkMaxX;
		column.ClipRect.Max.y = FLT_MAX
		column.ClipRect.ClipWithFull(host_clip_rect)

		// Mark column as Clipped (not in sight)
		// Note that scrolling tables (where inner_window != outer_window) handle Y clipped earlier in BeginTable() so IsVisibleY really only applies to non-scrolling tables.
		// FIXME-TABLE: Because InnerClipRect.Max.y is conservatively ==outer_window.ClipRect.Max.y, we never can mark columns _Above_ the scroll line as not IsVisibleY.
		// Taking advantage of LastOuterHeight would yield good results there...
		// FIXME-TABLE: Y clipping is disabled because it effectively means not submitting will reduce contents width which is fed to outer_window.DC.CursorMaxPos.x,
		// and this may be used (e.g. typically by outer_window using AlwaysAutoResize or outer_window's horizontal scrollbar, but could be something else).
		// Possible solution to preserve last known content width for clipped column. Test 'table_reported_size' fails when enabling Y clipping and window is resized small.
		column.IsVisibleX = (column.ClipRect.Max.x > column.ClipRect.Min.x)
		column.IsVisibleY = true           // (column.ClipRect.Max.y > column.ClipRect.Min.y);
		var is_visible = column.IsVisibleX //&& column.IsVisibleY;
		if is_visible {
			table.VisibleMaskByIndex |= ((ImU64)(1 << column_n))
		}

		// Mark column as requesting output from user. Note that fixed + non-resizable sets are auto-fitting at all times and therefore always request output.
		column.IsRequestOutput = is_visible || column.AutoFitQueue != 0 || column.CannotSkipItemsQueue != 0
		if column.IsRequestOutput {
			table.RequestOutputMaskByIndex |= ((ImU64)(1 << column_n))
		}

		// Mark column as SkipItems (ignoring all items/layout)
		column.IsSkipItems = !column.IsEnabled || table.HostSkipItems
		if column.IsSkipItems {
			IM_ASSERT(!is_visible)
		}

		// Update status flags
		column.Flags |= ImGuiTableColumnFlags_IsEnabled
		if is_visible {
			column.Flags |= ImGuiTableColumnFlags_IsVisible
		}
		if column.SortOrder != -1 {
			column.Flags |= ImGuiTableColumnFlags_IsSorted
		}
		if table.HoveredColumnBody == column_n {
			column.Flags |= ImGuiTableColumnFlags_IsHovered
		}

		// Alignment
		// FIXME-TABLE: This align based on the whole column width, not per-cell, and therefore isn't useful in
		// many cases (to be able to honor this we might be able to store a log of cells width, per row, for
		// visible rows, but nav/programmatic scroll would have visible artifacts.)
		//if (column.Flags & ImGuiTableColumnFlags_AlignRight)
		//    column.WorkMinX = ImMax(column.WorkMinX, column.MaxX - column.ContentWidthRowsUnfrozen);
		//else if (column.Flags & ImGuiTableColumnFlags_AlignCenter)
		//    column.WorkMinX = ImLerp(column.WorkMinX, ImMax(column.StartX, column.MaxX - column.ContentWidthRowsUnfrozen), 0.5f);

		// Reset content width variables
		column.ContentMaxXFrozen = column.WorkMinX
		column.ContentMaxXUnfrozen = column.WorkMinX
		column.ContentMaxXHeadersUsed = column.WorkMinX
		column.ContentMaxXHeadersIdeal = column.WorkMinX

		// Don't decrement auto-fit counters until container window got a chance to submit its items
		if !table.HostSkipItems {
			column.AutoFitQueue >>= 1
			column.CannotSkipItemsQueue >>= 1
		}

		if visible_n < int(table.FreezeColumnsCount) {
			host_clip_rect.Min.x = ImClamp(column.MaxX+TABLE_BORDER_SIZE, host_clip_rect.Min.x, host_clip_rect.Max.x)
		}

		offset_x += column.WidthGiven + table.CellSpacingX1 + table.CellSpacingX2 + table.CellPaddingX*2.0
		visible_n++
	}

	// [Part 7] Detect/store when we are hovering the unused space after the right-most column (so e.g. context menus can react on it)
	// Clear Resizable flag if none of our column are actually resizable (either via an explicit _NoResize flag, either
	// because of using _WidthAuto/_WidthStretch). This will hide the resizing option from the context menu.
	var unused_x1 = ImMax(table.WorkRect.Min.x, table.Columns[table.RightMostEnabledColumn].ClipRect.Max.x)
	if is_hovering_table && table.HoveredColumnBody == -1 {
		if g.IO.MousePos.x >= unused_x1 {
			table.HoveredColumnBody = (ImGuiTableColumnIdx)(table.ColumnsCount)
		}
	}
	if !has_resizable && (table.Flags&ImGuiTableFlags_Resizable) != 0 {
		table.Flags &= ^ImGuiTableFlags_Resizable
	}

	// [Part 8] Lock actual OuterRect/WorkRect right-most position.
	// This is done late to handle the case of fixed-columns tables not claiming more widths that they need.
	// Because of this we are careful with uses of WorkRect and InnerClipRect before this point.
	if table.RightMostStretchedColumn != -1 {
		table.Flags &= ^ImGuiTableFlags_NoHostExtendX
	}
	if table.Flags&ImGuiTableFlags_NoHostExtendX != 0 {
		table.OuterRect.Max.x = unused_x1
		table.WorkRect.Max.x = unused_x1
		table.InnerClipRect.Max.x = ImMin(table.InnerClipRect.Max.x, unused_x1)
	}
	table.InnerWindow.ParentWorkRect = table.WorkRect
	table.BorderX1 = table.InnerClipRect.Min.x // +((table.Flags & ImGuiTableFlags_BordersOuter) ? 0.0f : -1.0f);
	table.BorderX2 = table.InnerClipRect.Max.x // +((table.Flags & ImGuiTableFlags_BordersOuter) ? 0.0f : +1.0f);

	// [Part 9] Allocate draw channels and setup background cliprect
	TableSetupDrawChannels(table)

	// [Part 10] Hit testing on borders
	if table.Flags&ImGuiTableFlags_Resizable != 0 {
		TableUpdateBorders(table)
	}
	table.LastFirstRowHeight = 0.0
	table.IsLayoutLocked = true
	table.IsUsingHeaders = false

	// [Part 11] Context menu
	if table.IsContextPopupOpen && table.InstanceCurrent == table.InstanceInteracted {
		var context_menu_id = ImHashStr("##ContextMenu", 0, table.ID)
		if BeginPopupEx(context_menu_id, ImGuiWindowFlags_AlwaysAutoResize|ImGuiWindowFlags_NoTitleBar|ImGuiWindowFlags_NoSavedSettings) {
			TableDrawContextMenu(table)
			EndPopup()
		} else {
			table.IsContextPopupOpen = false
		}
	}

	// [Part 13] Sanitize and build sort specs before we have a change to use them for display.
	// This path will only be exercised when sort specs are modified before header rows (e.g. init or visibility change)
	if table.IsSortSpecsDirty && (table.Flags&ImGuiTableFlags_Sortable) != 0 {
		TableSortSpecsBuild(table)
	}

	// Initial state
	var inner_window = table.InnerWindow
	if table.Flags&ImGuiTableFlags_NoClip != 0 {
		table.DrawSplitter.SetCurrentChannel(inner_window.DrawList, TABLE_DRAW_CHANNEL_NOCLIP)
	} else {
		inner_window.DrawList.PushClipRect(inner_window.ClipRect.Min, inner_window.ClipRect.Max, false)
	}
}

// Process hit-testing on resizing borders. Actual size change will be applied in EndTable()
//   - Set table.HoveredColumnBorder with a short delay/timer to reduce feedback noise
//   - Submit ahead of table contents and header, use ImGuiButtonFlags_AllowItemOverlap to prioritize widgets
//     overlapping the same area.
func TableUpdateBorders(table *ImGuiTable) {
	var g = GImGui
	IM_ASSERT(table.Flags&ImGuiTableFlags_Resizable != 0)

	// At this point OuterRect height may be zero or under actual final height, so we rely on temporal coherency and
	// use the final height from last frame. Because this is only affecting _interaction_ with columns, it is not
	// really problematic (whereas the actual visual will be displayed in EndTable() and using the current frame height).
	// Actual columns highlight/render will be performed in EndTable() and not be affected.
	var hit_half_width = TABLE_RESIZE_SEPARATOR_HALF_THICKNESS
	var hit_y1 = table.OuterRect.Min.y
	var hit_y2_body = ImMax(table.OuterRect.Max.y, hit_y1+table.LastOuterHeight)
	var hit_y2_head = hit_y1 + table.LastFirstRowHeight

	for order_n := int(0); order_n < table.ColumnsCount; order_n++ {
		if (table.EnabledMaskByDisplayOrder & ((ImU64)(1 << order_n))) == 0 {
			continue
		}

		table.spanDisplayOrderToIndex(order_n)
		var column_n = table.DisplayOrderToIndex[order_n]
		table.spanColumns(int(column_n))
		var column = &table.Columns[column_n]
		if column.Flags&(ImGuiTableColumnFlags_NoResize|ImGuiTableColumnFlags_NoDirectResize_) != 0 {
			continue
		}

		// ImGuiTableFlags_NoBordersInBodyUntilResize will be honored in TableDrawBorders()
		var border_y2_hit = hit_y2_body
		if table.Flags&ImGuiTableFlags_NoBordersInBody != 0 {
			border_y2_hit = hit_y2_head
		}
		if (table.Flags&ImGuiTableFlags_NoBordersInBody) != 0 && !table.IsUsingHeaders {
			continue
		}

		if !column.IsVisibleX && table.LastResizedColumn != column_n {
			continue
		}

		var column_id = TableGetColumnResizeID(table, int(column_n), int(table.InstanceCurrent))
		var hit_rect = ImRect{ImVec2{column.MaxX - hit_half_width, hit_y1}, ImVec2{column.MaxX + hit_half_width, border_y2_hit}}
		//GetForegroundDrawList().AddRect(hit_rect.Min, hit_rect.Max, IM_COL32(255, 0, 0, 100));
		KeepAliveID(column_id)

		var hovered, held = false, false
		var pressed = ButtonBehavior(&hit_rect, column_id, &hovered, &held, ImGuiButtonFlags_FlattenChildren|ImGuiButtonFlags_AllowItemOverlap|ImGuiButtonFlags_PressedOnClick|ImGuiButtonFlags_PressedOnDoubleClick)
		if pressed && IsMouseDoubleClicked(0) {
			TableSetColumnWidthAutoSingle(table, int(column_n))
			ClearActiveID()
			held = false
			hovered = false
		}
		if held {
			if table.LastResizedColumn == -1 {
				if table.RightMostEnabledColumn != -1 {
					table.ResizeLockMinContentsX2 = table.Columns[table.RightMostEnabledColumn].MaxX
				} else {
					table.ResizeLockMinContentsX2 = -FLT_MAX
				}
			}
			table.ResizedColumn = (ImGuiTableColumnIdx)(column_n)
			table.InstanceInteracted = table.InstanceCurrent
		}
		if (hovered && g.HoveredIdTimer > TABLE_RESIZE_SEPARATOR_FEEDBACK_TIMER) || held {
			table.HoveredColumnBorder = (ImGuiTableColumnIdx)(column_n)
			SetMouseCursor(ImGuiMouseCursor_ResizeEW)
		}
	}
}

func TableUpdateColumnsWeightFromWidth(table *ImGuiTable) {
	IM_ASSERT(table.LeftMostStretchedColumn != -1 && table.RightMostStretchedColumn != -1)

	// Measure existing quantity
	var visible_weight float = 0.0
	var visible_width float = 0.0
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]
		if !column.IsEnabled || (column.Flags&ImGuiTableColumnFlags_WidthStretch) == 0 {
			continue
		}
		IM_ASSERT(column.StretchWeight > 0.0)
		visible_weight += column.StretchWeight
		visible_width += column.WidthRequest
	}
	IM_ASSERT(visible_weight > 0.0 && visible_width > 0.0)

	// Apply new weights
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]
		if !column.IsEnabled || (column.Flags&ImGuiTableColumnFlags_WidthStretch) == 0 {
			continue
		}
		column.StretchWeight = (column.WidthRequest / visible_width) * visible_weight
		IM_ASSERT(column.StretchWeight > 0.0)
	}
}

// FIXME-TABLE: This is a mess, need to redesign how we render borders (as some are also done in TableEndRow)
func TableDrawBorders(table *ImGuiTable) {
	var inner_window = table.InnerWindow
	if !table.OuterWindow.ClipRect.Overlaps(table.OuterRect) {
		return
	}

	var inner_drawlist = inner_window.DrawList
	table.DrawSplitter.SetCurrentChannel(inner_drawlist, TABLE_DRAW_CHANNEL_BG0)
	inner_drawlist.PushClipRect(table.Bg0ClipRectForDrawCmd.Min, table.Bg0ClipRectForDrawCmd.Max, false)

	// Draw inner border and resizing feedback
	var border_size = TABLE_BORDER_SIZE
	var draw_y1 = table.InnerRect.Min.y
	var draw_y2_body = table.InnerRect.Max.y
	var draw_y2_head = draw_y1
	if table.IsUsingHeaders {
		count := table.WorkRect.Min.y
		if table.FreezeRowsCount >= 1 {
			count = table.InnerRect.Min.y
		}
		draw_y2_head = ImMin(table.InnerRect.Max.y, count+table.LastFirstRowHeight)
	}
	if table.Flags&ImGuiTableFlags_BordersInnerV != 0 {
		for order_n := int(0); order_n < table.ColumnsCount; order_n++ {
			if (table.EnabledMaskByDisplayOrder & ((ImU64)(1 << order_n))) == 0 {
				continue
			}

			table.spanDisplayOrderToIndex(order_n)
			var column_n = table.DisplayOrderToIndex[order_n]
			table.spanColumns(int(order_n))
			var column = &table.Columns[column_n]
			var is_hovered = (table.HoveredColumnBorder == column_n)
			var is_resized = (table.ResizedColumn == column_n) && (table.InstanceInteracted == table.InstanceCurrent)
			var is_resizable = (column.Flags & (ImGuiTableColumnFlags_NoResize | ImGuiTableColumnFlags_NoDirectResize_)) == 0
			var is_frozen_separator = (int(table.FreezeColumnsCount) == order_n+1)
			if column.MaxX > table.InnerClipRect.Max.x && !is_resized {
				continue
			}

			// Decide whether right-most column is visible
			if column.NextEnabledColumn == -1 && !is_resizable {
				if (table.Flags&ImGuiTableFlags_SizingMask_) != ImGuiTableFlags_SizingFixedSame || (table.Flags&ImGuiTableFlags_NoHostExtendX) != 0 {
					continue
				}
			}
			if column.MaxX <= column.ClipRect.Min.x { // FIXME-TABLE FIXME-STYLE: Assume BorderSize==1, this is problematic if we want to increase the border size..
				continue
			}

			// Draw in outer window so right-most column won't be clipped
			// Always draw full height border when being resized/hovered, or on the delimitation of frozen column scrolling.
			var col ImU32
			var draw_y2 float
			if is_hovered || is_resized || is_frozen_separator {
				draw_y2 = draw_y2_body
				if is_resized {
					col = GetColorU32FromID(ImGuiCol_SeparatorActive, 1)
				} else if is_hovered {
					col = GetColorU32FromID(ImGuiCol_SeparatorHovered, 1)
				} else {
					col = table.BorderColorStrong
				}
			} else {
				if table.Flags&(ImGuiTableFlags_NoBordersInBody|ImGuiTableFlags_NoBordersInBodyUntilResize) != 0 {
					draw_y2 = draw_y2_head
				} else {
					draw_y2 = draw_y2_body
				}
				if table.Flags&(ImGuiTableFlags_NoBordersInBody|ImGuiTableFlags_NoBordersInBodyUntilResize) != 0 {
					col = table.BorderColorStrong
				} else {
					col = table.BorderColorLight
				}
			}

			if draw_y2 > draw_y1 {
				inner_drawlist.AddLine(&ImVec2{column.MaxX, draw_y1}, &ImVec2{column.MaxX, draw_y2}, col, border_size)
			}
		}
	}

	// Draw outer border
	// FIXME: could use AddRect or explicit VLine/HLine helper?
	if table.Flags&ImGuiTableFlags_BordersOuter != 0 {
		// Display outer border offset by 1 which is a simple way to display it without adding an extra draw call
		// (Without the offset, in outer_window it would be rendered behind cells, because child windows are above their
		// parent. In inner_window, it won't reach out over scrollbars. Another weird solution would be to display part
		// of it in inner window, and the part that's over scrollbars in the outer window..)
		// Either solution currently won't allow us to use a larger border size: the border would clipped.
		var outer_border = table.OuterRect
		var outer_col = table.BorderColorStrong
		if (table.Flags & ImGuiTableFlags_BordersOuter) == ImGuiTableFlags_BordersOuter {
			inner_drawlist.AddRect(outer_border.Min, outer_border.Max, outer_col, 0.0, 0, border_size)
		} else if table.Flags&ImGuiTableFlags_BordersOuterV != 0 {
			inner_drawlist.AddLine(&outer_border.Min, &ImVec2{outer_border.Min.x, outer_border.Max.y}, outer_col, border_size)
			inner_drawlist.AddLine(&ImVec2{outer_border.Max.x, outer_border.Min.y}, &outer_border.Max, outer_col, border_size)
		} else if table.Flags&ImGuiTableFlags_BordersOuterH != 0 {
			inner_drawlist.AddLine(&outer_border.Min, &ImVec2{outer_border.Max.x, outer_border.Min.y}, outer_col, border_size)
			inner_drawlist.AddLine(&ImVec2{outer_border.Min.x, outer_border.Max.y}, &outer_border.Max, outer_col, border_size)
		}
	}
	if (table.Flags&ImGuiTableFlags_BordersInnerH) != 0 && table.RowPosY2 < table.OuterRect.Max.y {
		// Draw bottom-most row border
		var border_y = table.RowPosY2
		if border_y >= table.BgClipRect.Min.y && border_y < table.BgClipRect.Max.y {
			inner_drawlist.AddLine(&ImVec2{table.BorderX1, border_y}, &ImVec2{table.BorderX2, border_y}, table.BorderColorLight, border_size)
		}
	}

	inner_drawlist.PopClipRect()
}

// Output context menu into current window (generally a popup)
// FIXME-TABLE: Ideally this should be writable by the user. Full programmatic access to that data?
func TableDrawContextMenu(table *ImGuiTable) {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return
	}

	var want_separator = false
	var column_n int
	if table.ContextPopupColumn >= 0 && int(table.ContextPopupColumn) < table.ColumnsCount {
		column_n = int(table.ContextPopupColumn)
	} else {
		column_n = -1
	}
	var column *ImGuiTableColumn
	if column_n != -1 {
		table.spanColumns(column_n)
		column = &table.Columns[column_n]
	}

	// Sizing
	if table.Flags&ImGuiTableFlags_Resizable != 0 {
		if column != nil {
			var can_resize = (column.Flags&ImGuiTableColumnFlags_NoResize) == 0 && column.IsEnabled
			if MenuItem("Size column to fit###SizeOne", "", nil, can_resize) {
				TableSetColumnWidthAutoSingle(table, column_n)
			}
		}

		var size_all_desc string
		if table.ColumnsEnabledFixedCount == table.ColumnsEnabledCount && (table.Flags&ImGuiTableFlags_SizingMask_) != ImGuiTableFlags_SizingFixedSame {
			size_all_desc = "Size all columns to fit###SizeAll" // All fixed
		} else {
			size_all_desc = "Size all columns to default###SizeAll" // All stretch or mixed
		}
		if MenuItem(size_all_desc, "", nil, false) {
			TableSetColumnWidthAutoAll(table)
		}
		want_separator = true
	}

	// Ordering
	if table.Flags&ImGuiTableFlags_Reorderable != 0 {
		if MenuItem("Reset order", "", nil, !table.IsDefaultDisplayOrder) {
			table.IsResetDisplayOrderRequest = true
		}
		want_separator = true
	}

	// Hiding / Visibility
	if (table.Flags & ImGuiTableFlags_Hideable) != 0 {
		if want_separator {
			Separator()
		}
		want_separator = true

		PushItemFlag(ImGuiItemFlags_SelectableDontClosePopup, true)
		for other_column_n := int(0); other_column_n < table.ColumnsCount; other_column_n++ {
			table.spanColumns(other_column_n)
			var other_column = &table.Columns[other_column_n]
			if other_column.Flags&ImGuiTableColumnFlags_Disabled != 0 {
				continue
			}

			var name = tableGetColumnName(table, other_column_n)
			if name == "" {
				name = "<Unknown>"
			}

			// Make sure we can't hide the last active column
			var menu_item_active = true
			if other_column.Flags&ImGuiTableColumnFlags_NoHide != 0 {
				menu_item_active = false
			}
			if other_column.IsUserEnabled && table.ColumnsEnabledCount <= 1 {
				menu_item_active = false
			}
			if MenuItem(name, "", &other_column.IsUserEnabled, menu_item_active) {
				other_column.IsUserEnabledNextFrame = !other_column.IsUserEnabled
			}
		}
		PopItemFlag()
	}
}

// This function reorder draw channels based on matching clip rectangle, to facilitate merging them. Called by EndTable().
// For simplicity we call it TableMergeDrawChannels() but in fact it only reorder channels + overwrite ClipRect,
// actual merging is done by table.DrawSplitter.Merge() which is called right after TableMergeDrawChannels().
//
// Columns where the contents didn't stray off their local clip rectangle can be merged. To achieve
// this we merge their clip rect and make them contiguous in the channel list, so they can be merged
// by the call to DrawSplitter.Merge() following to the call to this function.
// We reorder draw commands by arranging them into a maximum of 4 distinct groups:
//
//	1 group:               2 groups:              2 groups:              4 groups:
//	[ 0. ] no freeze       [ 0. ] row freeze      [ 01 ] col freeze      [ 01 ] row+col freeze
//	[ .. ]  or no scroll   [ 2. ]  and v-scroll   [ .. ]  and h-scroll   [ 23 ]  and v+h-scroll
//
// Each column itself can use 1 channel (row freeze disabled) or 2 channels (row freeze enabled).
// When the contents of a column didn't stray off its limit, we move its channels into the corresponding group
// based on its position (within frozen rows/columns groups or not).
// At the end of the operation our 1-4 groups will each have a ImDrawCmd using the same ClipRect.
// This function assume that each column are pointing to a distinct draw channel,
// otherwise merge_group.ChannelsCount will not match set bit count of merge_group.ChannelsMask.
//
// Column channels will not be merged into one of the 1-4 groups in the following cases:
//   - The contents stray off its clipping rectangle (we only compare the MaxX value, not the MinX value).
//     Direct ImDrawList calls won't be taken into account by default, if you use them make sure the ImGui:: bounds
//     matches, by e.g. calling SetCursorScreenPos().
//   - The channel uses more than one draw command itself. We drop all our attempt at merging stuff here..
//     we could do better but it's going to be rare and probably not worth the hassle.
//
// Columns for which the draw channel(s) haven't been merged with other will use their own ImDrawCmd.
//
// This function is particularly tricky to understand.. take a breath.
func TableMergeDrawChannels(table *ImGuiTable) {
	var g = GImGui
	var splitter = table.DrawSplitter
	var has_freeze_v = (table.FreezeRowsCount > 0)
	var has_freeze_h = (table.FreezeColumnsCount > 0)
	IM_ASSERT(splitter._Current == 0)

	// Track which groups we are going to attempt to merge, and which channels goes into each group.
	type MergeGroup struct {
		ClipRect      ImRect
		ChannelsCount int
		ChannelsMask  ImBitVector
	}
	var merge_group_mask int = 0x00
	var merge_groups [4]MergeGroup
	for i := range merge_groups {
		merge_groups[i].ChannelsMask = make(ImBitVector, IMGUI_TABLE_MAX_DRAW_CHANNELS)
	}

	// 1. Scan channels and take note of those which can be merged
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		if (table.VisibleMaskByIndex & ((ImU64)(1 << column_n))) == 0 {
			continue
		}
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]

		var merge_group_sub_count int = 1
		if has_freeze_v {
			merge_group_sub_count = 2
		}
		for merge_group_sub_n := int(0); merge_group_sub_n < merge_group_sub_count; merge_group_sub_n++ {
			var channel_no = column.DrawChannelUnfrozen
			if merge_group_sub_n == 0 {
				channel_no = column.DrawChannelFrozen
			}

			// Skip dummy channel (255 represents -1 when cast to uint8)
			if channel_no == ImGuiTableDrawChannelIdx(255) {
				continue
			}

			// Don't attempt to merge if there are multiple draw calls within the column
			var src_channel = &splitter._Channels[channel_no]
			if len(src_channel._CmdBuffer) > 0 && src_channel._CmdBuffer[len(src_channel._CmdBuffer)-1].ElemCount == 0 {
				src_channel._CmdBuffer = src_channel._CmdBuffer[:len(src_channel._CmdBuffer)-1]
			}
			if len(src_channel._CmdBuffer) != 1 {
				continue
			}

			// Find out the width of this merge group and check if it will fit in our column
			// (note that we assume that rendering didn't stray on the left direction. we should need a CursorMinPos to detect it)
			if (column.Flags & ImGuiTableColumnFlags_NoClip) == 0 {
				var content_max_x float
				if !has_freeze_v {
					content_max_x = ImMax(column.ContentMaxXUnfrozen, column.ContentMaxXHeadersUsed) // No row freeze
				} else if merge_group_sub_n == 0 {
					content_max_x = ImMax(column.ContentMaxXFrozen, column.ContentMaxXHeadersUsed) // Row freeze: use width before freeze
				} else {
					content_max_x = column.ContentMaxXUnfrozen // Row freeze: use width after freeze
				}
				if content_max_x > column.ClipRect.Max.x {
					continue
				}
			}

			var f1, f2 int = 1, 0
			if has_freeze_h && column_n < int(table.FreezeColumnsCount) {
				f1 = 0
			}
			if has_freeze_v && merge_group_sub_n == 0 {
				f2 = 2
			}

			var merge_group_n = f1 + f2
			IM_ASSERT(channel_no < IMGUI_TABLE_MAX_DRAW_CHANNELS)
			var merge_group = &merge_groups[merge_group_n]
			if merge_group.ChannelsCount == 0 {
				merge_group.ClipRect = ImRect{ImVec2{+FLT_MAX, +FLT_MAX}, ImVec2{-FLT_MAX, -FLT_MAX}}
			}
			merge_group.ChannelsMask.SetBit(int(channel_no))
			merge_group.ChannelsCount++
			merge_group.ClipRect.AddRect(ImRectFromVec4(&src_channel._CmdBuffer[0].ClipRect))
			merge_group_mask |= (1 << merge_group_n)
		}

		// Invalidate current draw channel
		// (we don't clear DrawChannelFrozen/DrawChannelUnfrozen solely to facilitate debugging/later inspection of data)
		column.DrawChannelCurrent = (ImGuiTableDrawChannelIdx)(255)
	}

	// [DEBUG] Display merge groups
	if false {
		if g.IO.KeyShift {
			for merge_group_n := 0; merge_group_n < len(merge_groups); merge_group_n++ {
				var merge_group = &merge_groups[merge_group_n]
				if merge_group.ChannelsCount == 0 {
					continue
				}
				var buf = fmt.Sprintf("MG%d:%d", merge_group_n, merge_group.ChannelsCount)
				var text_pos = merge_group.ClipRect.Min.Add(ImVec2{4, 4})
				var text_size = CalcTextSize(buf, true, -1)
				GetForegroundDrawList(nil).AddRectFilled(text_pos, text_pos.Add(text_size), IM_COL32(0, 0, 0, 255), 0, 0)
				GetForegroundDrawList(nil).AddText(text_pos, IM_COL32(255, 255, 0, 255), buf)
				GetForegroundDrawList(nil).AddRect(merge_group.ClipRect.Min, merge_group.ClipRect.Max, IM_COL32(255, 255, 0, 255), 0, 0, 1)
			}
		}
	}

	// 2. Rewrite channel list in our preferred order
	if merge_group_mask != 0 {
		// We skip channel 0 (Bg0/Bg1) and 1 (Bg2 frozen) from the shuffling since they won't move - see channels allocation in TableSetupDrawChannels().
		var LEADING_DRAW_CHANNELS int = 2
		g.DrawChannelsTempMergeBuffer = g.DrawChannelsTempMergeBuffer[:splitter._Count-LEADING_DRAW_CHANNELS] // Use shared temporary storage so the allocation gets amortized

		var dst_tmp = g.DrawChannelsTempMergeBuffer
		var remaining_mask = make(ImBitVector, IMGUI_TABLE_MAX_DRAW_CHANNELS) // We need 132-bit of storage
		remaining_mask.SetBitRange(LEADING_DRAW_CHANNELS, splitter._Count)
		remaining_mask.ClearBit(int(table.Bg2DrawChannelUnfrozen))
		IM_ASSERT(!has_freeze_v || int(table.Bg2DrawChannelUnfrozen) != TABLE_DRAW_CHANNEL_BG2_FROZEN)

		var f = LEADING_DRAW_CHANNELS
		if has_freeze_v {
			f = LEADING_DRAW_CHANNELS + 1
		}

		var remaining_count = splitter._Count - f
		//ImRect host_rect = (table.InnerWindow == table.OuterWindow) ? table.InnerClipRect : table.HostClipRect;
		var host_rect = table.HostClipRect
		for merge_group_n := int(0); merge_group_n < int(len(merge_groups)); merge_group_n++ {
			if merge_channels_count := merge_groups[merge_group_n].ChannelsCount; merge_channels_count != 0 {
				var merge_group = &merge_groups[merge_group_n]
				var merge_clip_rect = merge_group.ClipRect

				// Extend outer-most clip limits to match those of host, so draw calls can be merged even if
				// outer-most columns have some outer padding offsetting them from their parent ClipRect.
				// The principal cases this is dealing with are:
				// - On a same-window table (not scrolling = single group), all fitting columns ClipRect . will extend and match host ClipRect . will merge
				// - Columns can use padding and have left-most ClipRect.Min.x and right-most ClipRect.Max.x != from host ClipRect . will extend and match host ClipRect . will merge
				// FIXME-TABLE FIXME-WORKRECT: We are wasting a merge opportunity on tables without scrolling if column doesn't fit
				// within host clip rect, solely because of the half-padding difference between window.WorkRect and window.InnerClipRect.
				if (merge_group_n&1) == 0 || !has_freeze_h {
					merge_clip_rect.Min.x = ImMin(merge_clip_rect.Min.x, host_rect.Min.x)
				}
				if (merge_group_n&2) == 0 || !has_freeze_v {
					merge_clip_rect.Min.y = ImMin(merge_clip_rect.Min.y, host_rect.Min.y)
				}
				if (merge_group_n & 1) != 0 {
					merge_clip_rect.Max.x = ImMax(merge_clip_rect.Max.x, host_rect.Max.x)
				}
				if (merge_group_n&2) != 0 && (table.Flags&ImGuiTableFlags_NoHostExtendY) == 0 {
					merge_clip_rect.Max.y = ImMax(merge_clip_rect.Max.y, host_rect.Max.y)
				}
				remaining_count -= merge_group.ChannelsCount
				for n := 0; n < len(remaining_mask); n++ {
					remaining_mask[n] &= ^merge_group.ChannelsMask[n]
				}
				for n := int(0); n < splitter._Count && merge_channels_count != 0; n++ {
					// Copy + overwrite new clip rect
					if !merge_group.ChannelsMask.TestBit(n) {
						continue
					}
					merge_group.ChannelsMask.ClearBit(n)
					merge_channels_count--

					var channel = &splitter._Channels[n]
					IM_ASSERT(len(channel._CmdBuffer) == 1 && merge_clip_rect.ContainsRect(ImRectFromVec4(&channel._CmdBuffer[0].ClipRect)))
					channel._CmdBuffer[0].ClipRect = merge_clip_rect.ToVec4()
					dst_tmp[0] = *channel
					dst_tmp = dst_tmp[1:]
				}
			}

			// Make sure Bg2DrawChannelUnfrozen appears in the middle of our groups (whereas Bg0/Bg1 and Bg2 frozen are fixed to 0 and 1)
			if merge_group_n == 1 && has_freeze_v {
				dst_tmp[0] = splitter._Channels[table.Bg2DrawChannelUnfrozen]
				dst_tmp = dst_tmp[1:]
			}
		}

		// Append unmergeable channels that we didn't reorder at the end of the list
		for n := int(0); n < splitter._Count && remaining_count != 0; n++ {
			if !remaining_mask.TestBit(n) {
				continue
			}
			var channel = &splitter._Channels[n]
			dst_tmp[0] = *channel
			dst_tmp = dst_tmp[1:]
			remaining_count--
		}
		copy(splitter._Channels[LEADING_DRAW_CHANNELS:], g.DrawChannelsTempMergeBuffer[:(splitter._Count-LEADING_DRAW_CHANNELS)])
	}
}

func TableSortSpecsSanitize(table *ImGuiTable) {
	IM_ASSERT(table.Flags&ImGuiTableFlags_Sortable != 0)

	// Clear SortOrder from hidden column and verify that there's no gap or duplicate.
	var sort_order_count int = 0
	var sort_order_mask ImU64 = 0x00
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		table.spanColumns(column_n)
		var column = &table.Columns[column_n]
		if column.SortOrder != -1 && !column.IsEnabled {
			column.SortOrder = -1
		}
		if column.SortOrder == -1 {
			continue
		}
		sort_order_count++
		sort_order_mask |= ((ImU64)(1 << column.SortOrder))
		IM_ASSERT(sort_order_count < (int)(unsafe.Sizeof(sort_order_mask))*8)
	}

	var need_fix_linearize = ((ImU64)(1 << sort_order_count)) != (sort_order_mask + 1)
	var need_fix_single_sort_order = (sort_order_count > 1) && (table.Flags&ImGuiTableFlags_SortMulti) == 0
	if need_fix_linearize || need_fix_single_sort_order {
		var fixed_mask ImU64 = 0x00
		for sort_n := int(0); sort_n < sort_order_count; sort_n++ {
			// Fix: Rewrite sort order fields if needed so they have no gap or duplicate.
			// (e.g. SortOrder 0 disappeared, SortOrder 1..2 exists -. rewrite then as SortOrder 0..1)
			var column_with_smallest_sort_order int = -1
			for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
				table.spanColumns(column_n)
				if (fixed_mask&((ImU64)(1<<(ImU64)(column_n)))) == 0 && table.Columns[column_n].SortOrder != -1 {
					if column_with_smallest_sort_order == -1 || table.Columns[column_n].SortOrder < table.Columns[column_with_smallest_sort_order].SortOrder {
						column_with_smallest_sort_order = column_n
					}
				}
			}
			IM_ASSERT(column_with_smallest_sort_order != -1)
			fixed_mask |= ((ImU64)(1 << column_with_smallest_sort_order))
			table.spanColumns(column_with_smallest_sort_order)
			table.Columns[column_with_smallest_sort_order].SortOrder = (ImGuiTableColumnIdx)(sort_n)

			// Fix: Make sure only one column has a SortOrder if ImGuiTableFlags_MultiSortable is not set.
			if need_fix_single_sort_order {
				sort_order_count = 1
				for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
					if column_n != column_with_smallest_sort_order {
						table.spanColumns(column_n)
						table.Columns[column_n].SortOrder = -1
					}
				}
				break
			}
		}
	}

	// Fallback default sort order (if no column had the ImGuiTableColumnFlags_DefaultSort flag)
	if sort_order_count == 0 && (table.Flags&ImGuiTableFlags_SortTristate) == 0 {
		for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
			table.spanColumns(column_n)
			var column = &table.Columns[column_n]
			if column.IsEnabled && (column.Flags&ImGuiTableColumnFlags_NoSort) == 0 {
				sort_order_count = 1
				column.SortOrder = 0
				column.SortDirection = (ImGuiSortDirection)(TableGetColumnAvailSortDirection(column, 0))
				break
			}
		}
	}

	table.SortSpecsCount = (ImGuiTableColumnIdx)(sort_order_count)
}

func TableSortSpecsBuild(table *ImGuiTable) {
	var dirty = table.IsSortSpecsDirty
	if dirty {
		TableSortSpecsSanitize(table)
		//resize
		if int(table.SortSpecsCount) <= int(len(table.SortSpecsMulti)) {
			table.SortSpecsMulti = table.SortSpecsMulti[:table.SortSpecsCount]
		} else {
			table.SortSpecsMulti = append(table.SortSpecsMulti, make([]ImGuiTableColumnSortSpecs, int(table.SortSpecsCount)-int(len(table.SortSpecsMulti)))...)
		}
		table.SortSpecs.SpecsDirty = true // Mark as dirty for user
		table.IsSortSpecsDirty = false    // Mark as not dirty for us
	}

	// Write output
	var sort_specs []ImGuiTableColumnSortSpecs
	if table.SortSpecsCount == 0 {
		sort_specs = nil
	} else {
		if table.SortSpecsCount == 1 {
			sort_specs = []ImGuiTableColumnSortSpecs{table.SortSpecsSingle}
		} else {
			sort_specs = table.SortSpecsMulti
		}
	}

	if dirty && sort_specs != nil {
		for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
			table.spanColumns(column_n)
			var column = &table.Columns[column_n]
			if column.SortOrder == -1 {
				continue
			}
			IM_ASSERT(column.SortOrder < table.SortSpecsCount)
			var sort_spec = &sort_specs[column.SortOrder]
			sort_spec.ColumnUserID = column.UserID
			sort_spec.ColumnIndex = (ImS16)(column_n)
			sort_spec.SortOrder = (ImS16)(column.SortOrder)
			sort_spec.SortDirection = column.SortDirection
		}
	}

	table.SortSpecs.Specs = sort_specs
	table.SortSpecs.SpecsCount = int(table.SortSpecsCount)
}

// Calculate next sort direction that would be set after clicking the column
// - If the PreferSortDescending flag is set, we will default to a Descending direction on the first click.
// - Note that the PreferSortAscending flag is never checked, it is essentially the default and therefore a no-op.
func TableGetColumnNextSortDirection(column *ImGuiTableColumn) ImGuiSortDirection {
	IM_ASSERT(column.SortDirectionsAvailCount > 0)
	if column.SortOrder == -1 {
		return TableGetColumnAvailSortDirection(column, 0)
	}
	for n := int(0); n < 3; n++ {
		if column.SortDirection == TableGetColumnAvailSortDirection(column, n) {
			return TableGetColumnAvailSortDirection(column, (n+1)%int(column.SortDirectionsAvailCount))
		}
	}
	IM_ASSERT(false)
	return ImGuiSortDirection_None
}

// Fix sort direction if currently set on a value which is unavailable (e.g. activating NoSortAscending/NoSortDescending)
func TableFixColumnSortDirection(table *ImGuiTable, column *ImGuiTableColumn) {
	if column.SortOrder == -1 || (column.SortDirectionsAvailMask&(1<<column.SortDirection)) != 0 {
		return
	}
	column.SortDirection = (ImGuiSortDirection)(TableGetColumnAvailSortDirection(column, 0))
	table.IsSortSpecsDirty = true
}

// Note this is meant to be stored in column.WidthAuto, please generally use the WidthAuto field
func TableGetColumnWidthAuto(table *ImGuiTable, column *ImGuiTableColumn) float {
	var content_width_body = ImMax(column.ContentMaxXFrozen, column.ContentMaxXUnfrozen) - column.WorkMinX
	var content_width_headers = column.ContentMaxXHeadersIdeal - column.WorkMinX
	var width_auto = content_width_body
	if (column.Flags & ImGuiTableColumnFlags_NoHeaderWidth) == 0 {
		width_auto = ImMax(width_auto, content_width_headers)
	}

	// Non-resizable fixed columns preserve their requested width
	if (column.Flags&ImGuiTableColumnFlags_WidthFixed) != 0 && column.InitStretchWeightOrWidth > 0.0 {
		if (table.Flags&ImGuiTableFlags_Resizable) == 0 || (column.Flags&ImGuiTableColumnFlags_NoResize) != 0 {
			width_auto = column.InitStretchWeightOrWidth
		}
	}

	return ImMax(width_auto, table.MinColumnWidth)
}

// [Internal] Called by TableNextRow()
func TableBeginRow(table *ImGuiTable) {
	var window = table.InnerWindow
	IM_ASSERT(!table.IsInsideRow)

	// New row
	table.CurrentRow++
	table.CurrentColumn = -1
	table.RowBgColor[0] = IM_COL32_DISABLE
	table.RowBgColor[1] = IM_COL32_DISABLE
	table.RowCellDataCurrent = -1
	table.IsInsideRow = true

	// Begin frozen rows
	var next_y1 = table.RowPosY2
	if table.CurrentRow == 0 && table.FreezeRowsCount > 0 {
		next_y1 = table.OuterRect.Min.y
		window.DC.CursorPos.y = table.OuterRect.Min.y
	}

	table.RowPosY1 = next_y1
	table.RowPosY2 = next_y1
	table.RowTextBaseline = 0.0
	table.RowIndentOffsetX = window.DC.Indent.x - table.HostIndentX // Lock indent
	window.DC.PrevLineTextBaseOffset = 0.0
	window.DC.CursorMaxPos.y = next_y1

	// Making the header BG color non-transparent will allow us to overlay it multiple times when handling smooth dragging.
	if table.RowFlags&ImGuiTableRowFlags_Headers != 0 {
		TableSetBgColor(ImGuiTableBgTarget_RowBg0, GetColorU32FromID(ImGuiCol_TableHeaderBg, 1), 0)
		if table.CurrentRow == 0 {
			table.IsUsingHeaders = true
		}
	}
}

// [Internal] Called by TableNextRow()
func TableEndRow(table *ImGuiTable) {
	var g = GImGui
	var window = g.CurrentWindow
	IM_ASSERT(window == table.InnerWindow)
	IM_ASSERT(table.IsInsideRow)

	if table.CurrentColumn != -1 {
		TableEndCell(table)
	}

	// Logging
	if g.LogEnabled {
		LogRenderedText(nil, "|")
	}

	// Position cursor at the bottom of our row so it can be used for e.g. clipping calculation. However it is
	// likely that the next call to TableBeginCell() will reposition the cursor to take account of vertical padding.
	window.DC.CursorPos.y = table.RowPosY2

	// Row background fill
	var bg_y1 = table.RowPosY1
	var bg_y2 = table.RowPosY2
	var unfreeze_rows_actual = (table.CurrentRow+1 == int(table.FreezeRowsCount))
	var unfreeze_rows_request = (table.CurrentRow+1 == int(table.FreezeRowsRequest))
	if table.CurrentRow == 0 {
		table.LastFirstRowHeight = bg_y2 - bg_y1
	}

	var is_visible = (bg_y2 >= table.InnerClipRect.Min.y && bg_y1 <= table.InnerClipRect.Max.y)
	if is_visible {
		// Decide of background color for the row
		var bg_col0 ImU32 = 0
		var bg_col1 ImU32 = 0
		if table.RowBgColor[0] != IM_COL32_DISABLE {
			bg_col0 = table.RowBgColor[0]
		} else if table.Flags&ImGuiTableFlags_RowBg != 0 {
			c := ImGuiCol_TableRowBg
			if table.RowBgColorCounter&1 != 0 {
				c = ImGuiCol_TableRowBgAlt
			}
			bg_col0 = GetColorU32FromID(c, 1)
		}
		if table.RowBgColor[1] != IM_COL32_DISABLE {
			bg_col1 = table.RowBgColor[1]
		}

		// Decide of top border color
		var border_col ImU32 = 0
		var border_size = TABLE_BORDER_SIZE
		if table.CurrentRow > 0 || table.InnerWindow == table.OuterWindow {
			if table.Flags&ImGuiTableFlags_BordersInnerH != 0 {
				if table.LastRowFlags&ImGuiTableRowFlags_Headers != 0 {
					border_col = table.BorderColorStrong
				} else {
					border_col = table.BorderColorLight
				}
			}
		}

		var draw_cell_bg_color = table.RowCellDataCurrent >= 0
		var draw_strong_bottom_border = unfreeze_rows_actual
		if (bg_col0|bg_col1|border_col) != 0 || draw_strong_bottom_border || draw_cell_bg_color {
			// In theory we could call SetWindowClipRectBeforeSetChannel() but since we know TableEndRow() is
			// always followed by a change of clipping rectangle we perform the smallest overwrite possible here.
			if (table.Flags & ImGuiTableFlags_NoClip) == 0 {
				window.DrawList._CmdHeader.ClipRect = table.Bg0ClipRectForDrawCmd.ToVec4()
			}
			table.DrawSplitter.SetCurrentChannel(window.DrawList, TABLE_DRAW_CHANNEL_BG0)
		}

		// Draw row background
		// We soft/cpu clip this so all backgrounds and borders can share the same clipping rectangle
		if bg_col0|bg_col1 != 0 {
			var row_rect = ImRect{ImVec2{table.WorkRect.Min.x, bg_y1}, ImVec2{table.WorkRect.Max.x, bg_y2}}
			row_rect.ClipWith(table.BgClipRect)
			if bg_col0 != 0 && row_rect.Min.y < row_rect.Max.y {
				window.DrawList.AddRectFilled(row_rect.Min, row_rect.Max, bg_col0, 0, 0)
			}
			if bg_col1 != 0 && row_rect.Min.y < row_rect.Max.y {
				window.DrawList.AddRectFilled(row_rect.Min, row_rect.Max, bg_col1, 0, 0)
			}
		}

		// Draw cell background color
		if draw_cell_bg_color {
			// Only iterate up to RowCellDataCurrent (inclusive), not the entire slice
			for i := int(0); i <= int(table.RowCellDataCurrent); i++ {
				cell_data := &table.RowCellData[i]
				table.spanColumns(int(cell_data.Column))
				var column = &table.Columns[cell_data.Column]
				var cell_bg_rect = TableGetCellBgRect(table, int(cell_data.Column))
				cell_bg_rect.ClipWith(table.BgClipRect)
				cell_bg_rect.Min.x = ImMax(cell_bg_rect.Min.x, column.ClipRect.Min.x) // So that first column after frozen one gets clipped
				cell_bg_rect.Max.x = ImMin(cell_bg_rect.Max.x, column.MaxX)
				window.DrawList.AddRectFilled(cell_bg_rect.Min, cell_bg_rect.Max, cell_data.BgColor, 0, 0)
			}
		}

		// Draw top border
		if border_col != 0 && bg_y1 >= table.BgClipRect.Min.y && bg_y1 < table.BgClipRect.Max.y {
			window.DrawList.AddLine(&ImVec2{table.BorderX1, bg_y1}, &ImVec2{table.BorderX2, bg_y1}, border_col, border_size)
		}

		// Draw bottom border at the row unfreezing mark (always strong)
		if draw_strong_bottom_border && bg_y2 >= table.BgClipRect.Min.y && bg_y2 < table.BgClipRect.Max.y {
			window.DrawList.AddLine(&ImVec2{table.BorderX1, bg_y2}, &ImVec2{table.BorderX2, bg_y2}, table.BorderColorStrong, border_size)
		}
	}

	// End frozen rows (when we are past the last frozen row line, teleport cursor and alter clipping rectangle)
	// We need to do that in TableEndRow() instead of TableBeginRow() so the list clipper can mark end of row and
	// get the new cursor position.
	if unfreeze_rows_request {
		for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
			table.spanColumns(column_n)
			var column = &table.Columns[column_n]
			if column_n < int(table.FreezeColumnsCount) {
				column.NavLayerCurrent = int8(ImGuiNavLayer_Menu)
			} else {
				column.NavLayerCurrent = int8(ImGuiNavLayer_Main)
			}
		}
	}
	if unfreeze_rows_actual {
		IM_ASSERT(!table.IsUnfrozenRows)
		table.IsUnfrozenRows = true

		// BgClipRect starts as table.InnerClipRect, reduce it now and make BgClipRectForDrawCmd == BgClipRect
		var y0 = ImMax(table.RowPosY2+1, window.InnerClipRect.Min.y)
		table.BgClipRect.Min.y = ImMin(y0, window.InnerClipRect.Max.y)
		table.Bg2ClipRectForDrawCmd.Min.y = table.BgClipRect.Min.y
		table.BgClipRect.Max.y = window.InnerClipRect.Max.y
		table.Bg2ClipRectForDrawCmd.Max.y = window.InnerClipRect.Max.y
		table.Bg2DrawChannelCurrent = table.Bg2DrawChannelUnfrozen
		IM_ASSERT(table.Bg2ClipRectForDrawCmd.Min.y <= table.Bg2ClipRectForDrawCmd.Max.y)

		var row_height = table.RowPosY2 - table.RowPosY1
		table.RowPosY2 = table.WorkRect.Min.y + table.RowPosY2 - table.OuterRect.Min.y
		window.DC.CursorPos.y = table.RowPosY2
		table.RowPosY1 = table.RowPosY2 - row_height
		for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
			table.spanColumns(column_n)
			var column = &table.Columns[column_n]
			column.DrawChannelCurrent = column.DrawChannelUnfrozen
			column.ClipRect.Min.y = table.Bg2ClipRectForDrawCmd.Min.y
		}

		// Update cliprect ahead of TableBeginCell() so clipper can access to new ClipRect.Min.y
		SetWindowClipRectBeforeSetChannel(window, &table.Columns[0].ClipRect)
		// Skip if draw channel is dummy channel (255 represents -1 when cast to uint8)
		if table.Columns[0].DrawChannelCurrent != ImGuiTableDrawChannelIdx(255) {
			table.DrawSplitter.SetCurrentChannel(window.DrawList, int(table.Columns[0].DrawChannelCurrent))
		}
	}

	if table.RowFlags&ImGuiTableRowFlags_Headers == 0 {
		table.RowBgColorCounter++
	}
	table.IsInsideRow = false
}

// [Internal] Called by TableSetColumnIndex()/TableNextColumn()
// This is called very frequently, so we need to be mindful of unnecessary overhead.
// FIXME-TABLE FIXME-OPT: Could probably shortcut some things for non-active or clipped columns.
func TableBeginCell(table *ImGuiTable, column_n int) {
	table.spanColumns(column_n)
	var column = &table.Columns[column_n]
	var window = table.InnerWindow
	table.CurrentColumn = column_n

	// Start position is roughly ~~ CellRect.Min + CellPadding + Indent
	var start_x = column.WorkMinX
	if column.Flags&ImGuiTableColumnFlags_IndentEnable != 0 {
		start_x += table.RowIndentOffsetX // ~~ += window.DC.Indent.x - table.HostIndentX, except we locked it for the row.
	}

	window.DC.CursorPos.x = start_x
	window.DC.CursorPos.y = table.RowPosY1 + table.CellPaddingY
	window.DC.CursorMaxPos.x = window.DC.CursorPos.x
	window.DC.ColumnsOffset.x = start_x - window.Pos.x - window.DC.Indent.x // FIXME-WORKRECT
	window.DC.CurrLineTextBaseOffset = table.RowTextBaseline
	window.DC.NavLayerCurrent = (ImGuiNavLayer)(column.NavLayerCurrent)

	window.WorkRect.Min.y = window.DC.CursorPos.y
	window.WorkRect.Min.x = column.WorkMinX
	window.WorkRect.Max.x = column.WorkMaxX
	window.DC.ItemWidth = column.ItemWidth

	// To allow ImGuiListClipper to function we propagate our row height
	if !column.IsEnabled {
		window.DC.CursorPos.y = ImMax(window.DC.CursorPos.y, table.RowPosY2)
	}

	window.SkipItems = column.IsSkipItems
	if column.IsSkipItems {
		var g = GImGui
		g.LastItemData.ID = 0
		g.LastItemData.StatusFlags = 0
	}

	if table.Flags&ImGuiTableFlags_NoClip != 0 {
		// FIXME: if we end up drawing all borders/bg in EndTable, could remove this and just assert that channel hasn't changed.
		table.DrawSplitter.SetCurrentChannel(window.DrawList, TABLE_DRAW_CHANNEL_NOCLIP)
		//IM_ASSERT(table.DrawSplitter._Current == TABLE_DRAW_CHANNEL_NOCLIP);
	} else {
		// Skip if draw channel is dummy channel (255 represents -1 when cast to uint8)
		if column.DrawChannelCurrent != ImGuiTableDrawChannelIdx(255) {
			SetWindowClipRectBeforeSetChannel(window, &column.ClipRect)
			table.DrawSplitter.SetCurrentChannel(window.DrawList, int(column.DrawChannelCurrent))
		}
	}

	// Logging
	var g = GImGui
	if g.LogEnabled && !column.IsSkipItems {
		LogRenderedText(&window.DC.CursorPos, "|")
		g.LogLinePosY = FLT_MAX
	}
}

// [Internal] Called by TableNextRow()/TableSetColumnIndex()/TableNextColumn()
func TableEndCell(table *ImGuiTable) {
	var column = &table.Columns[table.CurrentColumn]
	var window = table.InnerWindow

	// Report maximum position so we can infer content size per column.
	var p_max_pos_x *float
	if (table.RowFlags & ImGuiTableRowFlags_Headers) != 0 {
		p_max_pos_x = &column.ContentMaxXHeadersUsed // Useful in case user submit contents in header row that is not a TableHeader() call
	} else {
		if table.IsUnfrozenRows {
			p_max_pos_x = &column.ContentMaxXUnfrozen
		} else {
			p_max_pos_x = &column.ContentMaxXFrozen
		}
	}
	*p_max_pos_x = ImMax(*p_max_pos_x, window.DC.CursorMaxPos.x)
	table.RowPosY2 = ImMax(table.RowPosY2, window.DC.CursorMaxPos.y+table.CellPaddingY)
	column.ItemWidth = window.DC.ItemWidth

	// Propagate text baseline for the entire row
	// FIXME-TABLE: Here we propagate text baseline from the last line of the cell.. instead of the first one.
	table.RowTextBaseline = ImMax(table.RowTextBaseline, window.DC.PrevLineTextBaseOffset)
}

// Return the cell rectangle based on currently known height.
//   - Important: we generally don't know our row height until the end of the row, so Max.y will be incorrect in many situations.
//     The only case where this is correct is if we provided a min_row_height to TableNextRow() and don't go below it.
//   - Important: if ImGuiTableFlags_PadOuterX is set but ImGuiTableFlags_PadInnerX is not set, the outer-most left and right
//     columns report a small offset so their CellBgRect can extend up to the outer border.
func TableGetCellBgRect(table *ImGuiTable, column_n int) ImRect {
	var column = &table.Columns[column_n]
	var x1 = column.MinX
	var x2 = column.MaxX
	if column.PrevEnabledColumn == -1 {
		x1 -= table.CellSpacingX1
	}
	if column.NextEnabledColumn == -1 {
		x2 += table.CellSpacingX2
	}
	return ImRect{ImVec2{x1, table.RowPosY1}, ImVec2{x2, table.RowPosY2}}
}

func tableGetColumnName(table *ImGuiTable, column_n int) string {
	if !table.IsLayoutLocked && column_n >= int(table.DeclColumnsCount) {
		return "" // NameOffset is invalid at this point
	}
	table.spanColumns(column_n)
	var column = &table.Columns[column_n]
	if column.NameOffset == -1 {
		return ""
	}
	return table.ColumnsNames[column.NameOffset]
}

// Return the resizing ID for the right-side of the given column.
func TableGetColumnResizeID(table *ImGuiTable, column_n int, instance_no int) ImGuiID {
	IM_ASSERT(column_n >= 0 && column_n < table.ColumnsCount)
	var id = table.ID + 1 + uint(instance_no*table.ColumnsCount) + uint(column_n)
	return id
}

// Maximum column content width given current layout. Use column.MinX so this value on a per-column basis.
func TableGetMaxColumnWidth(table *ImGuiTable, column_n int) float {
	var column = &table.Columns[column_n]
	var max_width float = FLT_MAX
	var min_column_distance = table.MinColumnWidth + table.CellPaddingX*2.0 + table.CellSpacingX1 + table.CellSpacingX2
	if table.Flags&ImGuiTableFlags_ScrollX != 0 {
		// Frozen columns can't reach beyond visible width else scrolling will naturally break.
		// (we use DisplayOrder as within a set of multiple frozen column reordering is possible)
		if column.DisplayOrder < table.FreezeColumnsRequest {
			max_width = (table.InnerClipRect.Max.x - float(table.FreezeColumnsRequest-column.DisplayOrder)*min_column_distance) - column.MinX
			max_width = max_width - table.OuterPaddingX - table.CellPaddingX - table.CellSpacingX2
		}
	} else if (table.Flags & ImGuiTableFlags_NoKeepColumnsVisible) == 0 {
		// If horizontal scrolling if disabled, we apply a final lossless shrinking of columns in order to make
		// sure they are all visible. Because of this we also know that all of the columns will always fit in
		// table.WorkRect and therefore in table.InnerRect (because ScrollX is off)
		// FIXME-TABLE: This is solved incorrectly but also quite a difficult problem to fix as we also want ClipRect width to match.
		// See "table_width_distrib" and "table_width_keep_visible" tests
		max_width = table.WorkRect.Max.x - float(table.ColumnsEnabledCount-column.IndexWithinEnabledSet-1)*min_column_distance - column.MinX
		//max_width -= table.CellSpacingX1;
		max_width -= table.CellSpacingX2
		max_width -= table.CellPaddingX * 2.0
		max_width -= table.OuterPaddingX
	}
	return max_width
}

// Disable clipping then auto-fit, will take 2 frames
// (we don't take a shortcut for unclipped columns to reduce inconsistencies when e.g. resizing multiple columns)
func TableSetColumnWidthAutoSingle(table *ImGuiTable, column_n int) {
	// Single auto width uses auto-fit
	var column = &table.Columns[column_n]
	if !column.IsEnabled {
		return
	}
	column.CannotSkipItemsQueue = (1 << 0)
	table.AutoFitSingleColumn = (ImGuiTableColumnIdx)(column_n)
}

func TableSetColumnWidthAutoAll(table *ImGuiTable) {
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		var column = &table.Columns[column_n]
		if !column.IsEnabled && (column.Flags&ImGuiTableColumnFlags_WidthStretch) == 0 { // Cannot reset weight of hidden stretch column
			continue
		}
		column.CannotSkipItemsQueue = (1 << 0)
		column.AutoFitQueue = (1 << 1)
	}
}

// Remove Table (currently only used by TestEngine)
func TableRemove(table *ImGuiTable) {
	//IMGUI_DEBUG_LOG("TableRemove() id=0x%08X\n", table.ID);
	var g = GImGui
	var table_idx uint
	for i := range g.Tables {
		if g.Tables[i] == table {
			table_idx = i
		}
	}
	delete(g.Tables, table.ID)
	delete(g.TablesLastTimeActive, int(table_idx))
}

// Free up/compact internal Table buffers for when it gets unused
func TableGcCompactTransientBuffers(table *ImGuiTable) {
	//IMGUI_DEBUG_LOG("TableGcCompactTransientBuffers() id=0x%08X\n", table.ID);
	var g = GImGui
	IM_ASSERT(!table.MemoryCompacted)
	table.SortSpecs.Specs = nil
	table.SortSpecsMulti = nil
	table.IsSortSpecsDirty = true // FIXME: shouldn't have to leak into user performing a sort
	table.ColumnsNames = nil
	table.MemoryCompacted = true
	for n := int(0); n < table.ColumnsCount; n++ {
		table.Columns[n].NameOffset = -1
	}

	var table_idx uint
	for i := range g.Tables {
		if g.Tables[i] == table {
			table_idx = i
		}
	}

	g.TablesLastTimeActive[int(table_idx)] = -1.0
}

func TableGcCompactTransientBuffersTempData(temp_data *ImGuiTableTempData) {
	temp_data.DrawSplitter.ClearFreeMemory()
	temp_data.LastTimeActive = -1.0
}

// Compact and remove unused settings data (currently only used by TestEngine)
func TableGcCompactSettings() {
	var g = GImGui
	var required_memory int = 0
	for _, settings := range g.SettingsTables {
		if settings.ID != 0 {
			required_memory += (int)(TableSettingsCalcChunkSize(int(settings.ColumnsCount)))
		}
	}
	if required_memory == int(len(g.SettingsTables)) {
		return
	}
	var new_chunk_stream = make([]ImGuiTableSettings, required_memory)
	for _, settings := range g.SettingsTables {
		if settings.ID != 0 {
			new_chunk_stream = append(new_chunk_stream, settings)
		}
	}

	g.SettingsTables = new_chunk_stream
}

// Tables: Settings
func TableLoadSettings(table *ImGuiTable) {
	var g = GImGui
	table.IsSettingsRequestLoad = false
	if table.Flags&ImGuiTableFlags_NoSavedSettings != 0 {
		return
	}

	// Bind settings
	var settings *ImGuiTableSettings
	if table.SettingsOffset == -1 {
		settings = TableSettingsFindByID(table.ID)
		if settings == nil {
			return
		}
		if int(settings.ColumnsCount) != table.ColumnsCount { // Allow settings if columns count changed. We could otherwise decide to return...
			table.IsSettingsDirty = true
		}
		for i := range g.SettingsTables {
			if &g.SettingsTables[i] == settings {
				table.SettingsOffset = int(i)
				return
			}
		}
	} else {
		settings = TableGetBoundSettings(table)
	}

	table.SettingsLoadedFlags = settings.SaveFlags
	table.RefScale = settings.RefScale

	// Serialize ImGuiTableSettings/ImGuiTableColumnSettings into ImGuiTable/ImGuiTableColumn
	var column_settings = settings.Columns
	var display_order_mask ImU64 = 0
	for data_n := int(0); data_n < int(settings.ColumnsCount); data_n, column_settings = data_n+1, column_settings[1:] {
		var column_n = int(column_settings[0].Index)
		if column_n < 0 || column_n >= table.ColumnsCount {
			continue
		}

		var column = &table.Columns[column_n]
		if settings.SaveFlags&ImGuiTableFlags_Resizable != 0 {
			if column_settings[0].IsStretch != 0 {
				column.StretchWeight = column_settings[0].WidthOrWeight
			} else {
				column.WidthRequest = column_settings[0].WidthOrWeight
			}
			column.AutoFitQueue = 0x00
		}
		if settings.SaveFlags&ImGuiTableFlags_Reorderable != 0 {
			column.DisplayOrder = column_settings[0].DisplayOrder
		} else {
			column.DisplayOrder = (ImGuiTableColumnIdx)(column_n)
		}
		display_order_mask |= (ImU64)(1 << column.DisplayOrder)
		column.IsUserEnabled = istrue(int(column_settings[0].IsEnabled))
		column.IsUserEnabledNextFrame = istrue(int(column_settings[0].IsEnabled))
		column.SortOrder = column_settings[0].SortOrder
		column.SortDirection = ImGuiSortDirection(column_settings[0].SortDirection)
	}

	// Validate and fix invalid display order data
	var expected_display_order_mask ImU64
	if settings.ColumnsCount == 64 {
		expected_display_order_mask = INT_MAX
	} else {
		expected_display_order_mask = ((ImU64)(1<<settings.ColumnsCount) - 1)
	}
	if display_order_mask != expected_display_order_mask {
		for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
			table.Columns[column_n].DisplayOrder = (ImGuiTableColumnIdx)(column_n)
		}
	}

	// Rebuild index
	for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
		table.DisplayOrderToIndex[table.Columns[column_n].DisplayOrder] = (ImGuiTableColumnIdx)(column_n)
	}
}

func TableSaveSettings(table *ImGuiTable) {
	table.IsSettingsDirty = false
	if table.Flags&ImGuiTableFlags_NoSavedSettings != 0 {
		return
	}

	// Bind or create settings data
	var g = GImGui
	var settings = TableGetBoundSettings(table)
	if settings == nil {
		settings = TableSettingsCreate(table.ID, table.ColumnsCount)
		for i := range g.SettingsTables {
			if &g.SettingsTables[i] == settings {
				table.SettingsOffset = int(i)
				break
			}
		}
	}
	settings.ColumnsCount = (ImGuiTableColumnIdx)(table.ColumnsCount)

	// Serialize ImGuiTable/ImGuiTableColumn into ImGuiTableSettings/ImGuiTableColumnSettings
	IM_ASSERT(settings.ID == table.ID)
	IM_ASSERT(int(settings.ColumnsCount) == table.ColumnsCount && settings.ColumnsCountMax >= settings.ColumnsCount)
	var column = table.Columns
	var column_settings = settings.Columns

	var save_ref_scale = false
	settings.SaveFlags = ImGuiTableFlags_None
	for n := int(0); n < table.ColumnsCount; n, column, column_settings = n+1, column[1:], column_settings[1:] {
		var width_or_weight = column[0].WidthRequest
		if (column[0].Flags & ImGuiTableColumnFlags_WidthStretch) != 0 {
			width_or_weight = column[0].StretchWeight
		}
		column_settings[0].WidthOrWeight = width_or_weight
		column_settings[0].Index = (ImGuiTableColumnIdx)(n)
		column_settings[0].DisplayOrder = column[0].DisplayOrder
		column_settings[0].SortOrder = column[0].SortOrder
		column_settings[0].SortDirection = uint8(column[0].SortDirection)
		column_settings[0].IsEnabled = uint8(bool2int(column[0].IsUserEnabled))
		column_settings[0].IsStretch = 0
		if (column[0].Flags & ImGuiTableColumnFlags_WidthStretch) != 0 {
			column_settings[0].IsStretch = 1
		}
		if (column[0].Flags & ImGuiTableColumnFlags_WidthStretch) == 0 {
			save_ref_scale = true
		}

		// We skip saving some data in the .ini file when they are unnecessary to restore our state.
		// Note that fixed width where initial width was derived from auto-fit will always be saved as InitStretchWeightOrWidth will be 0.0f.
		// FIXME-TABLE: We don't have logic to easily compare SortOrder to DefaultSortOrder yet so it's always saved when present.
		if width_or_weight != column[0].InitStretchWeightOrWidth {
			settings.SaveFlags |= ImGuiTableFlags_Resizable
		}
		if int(column[0].DisplayOrder) != n {
			settings.SaveFlags |= ImGuiTableFlags_Reorderable
		}
		if column[0].SortOrder != -1 {
			settings.SaveFlags |= ImGuiTableFlags_Sortable
		}
		if column[0].IsUserEnabled != ((column[0].Flags & ImGuiTableColumnFlags_DefaultHide) == 0) {
			settings.SaveFlags |= ImGuiTableFlags_Hideable
		}
	}
	settings.SaveFlags &= table.Flags
	settings.RefScale = 0.0
	if save_ref_scale {
		settings.RefScale = table.RefScale
	}

	MarkIniSettingsDirty()
}

// Restore initial state of table (with or without saved settings)
func TableResetSettings(table *ImGuiTable) {
	table.IsInitializing = true
	table.IsSettingsDirty = true
	table.IsResetAllRequest = false
	table.IsSettingsRequestLoad = false              // Don't reload from ini
	table.SettingsLoadedFlags = ImGuiTableFlags_None // Mark as nothing loaded so our initialized data becomes authoritative
}

// Get settings for a given table, NULL if none
func TableGetBoundSettings(table *ImGuiTable) *ImGuiTableSettings {
	if table.SettingsOffset != -1 {
		var g = GImGui
		var settings = &g.SettingsTables[table.SettingsOffset]
		IM_ASSERT(settings.ID == table.ID)
		if int(settings.ColumnsCountMax) >= table.ColumnsCount {
			return settings // OK
		}
		settings.ID = 0 // Invalidate storage, we won't fit because of a count change
	}
	return nil
}

func TableSettingsInstallHandler(context *ImGuiContext) {
	var g = context
	var ini_handler ImGuiSettingsHandler
	ini_handler.TypeName = "Table"
	ini_handler.TypeHash = ImHashStr("Table", uintptr(len("Table")), 0)
	ini_handler.ClearAllFn = TableSettingsHandler_ClearAll
	ini_handler.ReadOpenFn = TableSettingsHandler_ReadOpen
	ini_handler.ReadLineFn = TableSettingsHandler_ReadLine
	ini_handler.ApplyAllFn = TableSettingsHandler_ApplyAll
	ini_handler.WriteAllFn = TableSettingsHandler_WriteAll
	g.SettingsHandlers = append(g.SettingsHandlers, ini_handler)
}

func TableSettingsCreate(id ImGuiID, columns_count int) *ImGuiTableSettings {
	var g = GImGui
	g.SettingsTables = append(g.SettingsTables, ImGuiTableSettings{})
	var settings = &g.SettingsTables[len(g.SettingsTables)-1]
	TableSettingsInit(settings, id, columns_count, columns_count)
	settings.Columns = make([]ImGuiTableColumnSettings, columns_count)
	return settings
}

// Find existing settings
func TableSettingsFindByID(id ImGuiID) *ImGuiTableSettings {
	// FIXME-OPT: Might want to store a lookup map for this?
	var g = GImGui
	for i, settings := range g.SettingsTables {
		if settings.ID == id {
			return &g.SettingsTables[i]
		}
	}
	return nil
}

func TableSettingsHandler_ClearAll(ctx *ImGuiContext, _ *ImGuiSettingsHandler) {
	var g = ctx
	for i := uint(0); i != uint(len(g.Tables)); i++ {
		if table := g.Tables[i]; table != nil {
			table.SettingsOffset = -1
		}
	}
	g.SettingsTables = g.SettingsTables[:0]
}

// Apply to existing windows (if any)
func TableSettingsHandler_ApplyAll(ctx *ImGuiContext, _ *ImGuiSettingsHandler) {
	var g = ctx
	for i := uint(0); i != uint(len(g.Tables)); i++ {
		if table := g.Tables[i]; table != nil {
			table.IsSettingsRequestLoad = true
			table.SettingsOffset = -1
		}
	}
}

func TableSettingsHandler_ReadOpen(ctx *ImGuiContext, _ *ImGuiSettingsHandler, name string) any {
	var id ImGuiID = 0
	var columns_count int = 0
	if n, _ := fmt.Scanf(name, "0x%08X,%d", &id, &columns_count); n < 2 {
		return nil
	}

	if settings := TableSettingsFindByID(id); settings != nil {
		if int(settings.ColumnsCountMax) >= columns_count {
			TableSettingsInit(settings, id, columns_count, int(settings.ColumnsCountMax)) // Recycle
			return settings
		}
		settings.ID = 0 // Invalidate storage, we won't fit because of a count change
	}
	return TableSettingsCreate(id, columns_count)
}

func TableSettingsHandler_ReadLine(ctx *ImGuiContext, _ *ImGuiSettingsHandler, entry any, line string) {
	// "Column 0  UserID=0x42AD2D21 Width=100 Visible=1 Order=0 Sort=0v"
	var settings = entry.(*ImGuiTableSettings)
	var f float = 0.0
	var column_n, r int
	var x uint

	if n, _ := fmt.Sscanf(line, "RefScale=%f", &f); n == 1 {
		settings.RefScale = f
		return
	}

	if n, _ := fmt.Sscanf(line, "Column %d%n", &column_n, &r); n == 1 {
		if column_n < 0 || column_n >= int(settings.ColumnsCount) {
			return
		}
		line = ImStrSkipBlank(line[r:])
		var c byte = 0
		var column = settings.Columns[column_n]
		column.Index = (ImGuiTableColumnIdx)(column_n)
		if n, _ := fmt.Sscanf(line, "UserID=0x%08X%n", (*ImU32)(&x), &r); n == 1 {
			line = ImStrSkipBlank(line[r:])
			column.UserID = (ImGuiID)(n)
		}
		if n, _ := fmt.Sscanf(line, "Width=%d%n", &n, &r); n == 1 {
			line = ImStrSkipBlank(line[r:])
			column.WidthOrWeight = (float)(n)
			column.IsStretch = 0
			settings.SaveFlags |= ImGuiTableFlags_Resizable
		}
		if n, _ := fmt.Sscanf(line, "Weight=%f%n", &f, &r); n == 1 {
			line = ImStrSkipBlank(line[r:])
			column.WidthOrWeight = f
			column.IsStretch = 1
			settings.SaveFlags |= ImGuiTableFlags_Resizable
		}
		if n, _ := fmt.Sscanf(line, "Visible=%d%n", &n, &r); n == 1 {
			line = ImStrSkipBlank(line[r:])
			column.IsEnabled = (ImU8)(n)
			settings.SaveFlags |= ImGuiTableFlags_Hideable
		}
		if n, _ := fmt.Sscanf(line, "Order=%d%n", &n, &r); n == 1 {
			line = ImStrSkipBlank(line[r:])
			column.DisplayOrder = (ImGuiTableColumnIdx)(n)
			settings.SaveFlags |= ImGuiTableFlags_Reorderable
		}

		dir := ImGuiSortDirection_Ascending
		if c == '^' {
			dir = ImGuiSortDirection_Descending
		}

		if n, _ := fmt.Sscanf(line, "Sort=%d%c%n", &n, &c, &r); n == 2 {
			// this value will never be used anymore
			// line = ImStrSkipBlank(line[r:])
			column.SortDirection = uint8(dir)
			settings.SaveFlags |= ImGuiTableFlags_Sortable
		}
	}
}

func TableSettingsHandler_WriteAll(ctx *ImGuiContext, handler *ImGuiSettingsHandler, buf *ImGuiTextBuffer) {
	var g = ctx
	for _, settings := range g.SettingsTables {
		if settings.ID == 0 { // Skip ditched settings
			continue
		}

		// TableSaveSettings() may clear some of those flags when we establish that the data can be stripped
		// (e.g. Order was unchanged)
		var save_size = (settings.SaveFlags & ImGuiTableFlags_Resizable) != 0
		var save_visible = (settings.SaveFlags & ImGuiTableFlags_Hideable) != 0
		var save_order = (settings.SaveFlags & ImGuiTableFlags_Reorderable) != 0
		var save_sort = (settings.SaveFlags & ImGuiTableFlags_Sortable) != 0
		if !save_size && !save_visible && !save_order && !save_sort {
			continue
		}

		*buf = append(*buf, make([]byte, 0, 30+settings.ColumnsCount*50)...) // ballpark reserve
		*buf = append(*buf, []byte(fmt.Sprintf("[%s][0x%08X,%d]\n", handler.TypeName, settings.ID, settings.ColumnsCount))...)

		if settings.RefScale != 0.0 {
			*buf = append(*buf, []byte(fmt.Sprintf("RefScale=%g\n", settings.RefScale))...)
		}
		var column = settings.Columns
		for column_n := int(0); column_n < int(settings.ColumnsCount); column_n, column = column_n+1, column[1:] {
			// "Column 0  UserID=0x42AD2D21 Width=100 Visible=1 Order=0 Sort=0v"
			var save_column = column[0].UserID != 0 || save_size || save_visible || save_order || (save_sort && column[0].SortOrder != -1)
			if !save_column {
				continue
			}
			*buf = append(*buf, []byte(fmt.Sprintf("Column %-2d", column_n))...)
			if column[0].UserID != 0 {
				*buf = append(*buf, []byte(fmt.Sprintf(" UserID=%08X", column[0].UserID))...)
			}
			if save_size && column[0].IsStretch != 0 {
				*buf = append(*buf, []byte(fmt.Sprintf(" Weight=%.4f", column[0].WidthOrWeight))...)
			}
			if save_size && column[0].IsStretch == 0 {
				*buf = append(*buf, []byte(fmt.Sprintf(" Width=%d", (int)(column[0].WidthOrWeight)))...)
			}
			if save_visible {
				*buf = append(*buf, []byte(fmt.Sprintf(" Visible=%d", column[0].IsEnabled))...)
			}
			if save_order {
				*buf = append(*buf, []byte(fmt.Sprintf(" Order=%d", column[0].DisplayOrder))...)
			}
			if save_sort && column[0].SortOrder != -1 {
				dir := '^'
				if ImGuiSortDirection(column[0].SortDirection) == ImGuiSortDirection_Ascending {
					dir = 'v'
				}
				*buf = append(*buf, []byte(fmt.Sprintf(" Sort=%d%v", column[0].SortOrder, string(dir)))...)
			}
			*buf = append(*buf, '\n')
		}
		*buf = append(*buf, '\n')
	}
}

func DebugNodeTableGetSizingPolicyDesc(sizing_policy ImGuiTableFlags) string {
	sizing_policy &= ImGuiTableFlags_SizingMask_
	if sizing_policy == ImGuiTableFlags_SizingFixedFit {
		return "FixedFit"
	}
	if sizing_policy == ImGuiTableFlags_SizingFixedSame {
		return "FixedSame"
	}
	if sizing_policy == ImGuiTableFlags_SizingStretchProp {
		return "StretchProp"
	}
	if sizing_policy == ImGuiTableFlags_SizingStretchSame {
		return "StretchSame"
	}
	return "N/A"
}
