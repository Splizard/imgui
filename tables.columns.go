package imgui

const COLUMNS_HIT_RECT_HALF_WIDTH float = 4

// PushColumnsBackground Get into the columns background draw command (which is generally the same draw command as before we called BeginColumns)
func PushColumnsBackground() {
	window := GetCurrentWindowRead()
	var columns = window.DC.CurrentColumns
	if columns.Count == 1 {
		return
	}

	// Optimization: avoid SetCurrentChannel() + PushClipRect()
	columns.HostBackupClipRect = window.ClipRect
	SetWindowClipRectBeforeSetChannel(window, &columns.HostInitialClipRect)
	columns.Splitter.SetCurrentChannel(window.DrawList, 0)
}

// Internal Columns API (this is not exposed because we will encourage transitioning to the Tables API)

// SetWindowClipRectBeforeSetChannel [Internal] Small optimization to avoid calls to PopClipRect/SetCurrentChannel/PushClipRect in sequences,
// they would meddle many times with the underlying ImDrawCmd.
// Instead, we do a preemptive overwrite of clipping rectangle _without_ altering the command-buffer and let
// the subsequent single call to SetCurrentChannel() does it things once.
func SetWindowClipRectBeforeSetChannel(window *ImGuiWindow, clip_rect *ImRect) {
	var clip_rect_vec4 = clip_rect.ToVec4()
	window.ClipRect = *clip_rect
	window.DrawList._CmdHeader.ClipRect = clip_rect_vec4
	window.DrawList._ClipRectStack[len(window.DrawList._ClipRectStack)-1] = clip_rect_vec4
}

// BeginColumns setup number of columns. use an identifier to distinguish multiple column sets. close with EndColumns().
func BeginColumns(str_id string, columns_count int, flags ImGuiOldColumnFlags) {
	g := GImGui
	window := GetCurrentWindow()

	IM_ASSERT(columns_count >= 1)
	IM_ASSERT(window.DC.CurrentColumns == nil) // Nested columns are currently not supported

	// Acquire storage for the columns set
	var id = GetColumnsID(str_id, columns_count)
	var columns = FindOrCreateColumns(window, id)
	IM_ASSERT(columns.ID == id)
	columns.Current = 0
	columns.Count = columns_count
	columns.Flags = flags
	window.DC.CurrentColumns = columns

	columns.HostCursorPosY = window.DC.CursorPos.y
	columns.HostCursorMaxPosX = window.DC.CursorMaxPos.x
	columns.HostInitialClipRect = window.ClipRect
	columns.HostBackupParentWorkRect = window.ParentWorkRect
	window.ParentWorkRect = window.WorkRect

	// Set state for first column
	// We aim so that the right-most column will have the same clipping width as other after being clipped by parent ClipRect
	var column_padding = g.Style.ItemSpacing.x
	var half_clip_extend_x = ImFloor(max(window.WindowPadding.x*0.5, window.WindowBorderSize))
	var max_1 = window.WorkRect.Max.x + column_padding - max(column_padding-window.WindowPadding.x, 0.0)
	var max_2 = window.WorkRect.Max.x + half_clip_extend_x
	columns.OffMinX = window.DC.Indent.x - column_padding + max(column_padding-window.WindowPadding.x, 0.0)
	columns.OffMaxX = max(min(max_1, max_2)-window.Pos.x, columns.OffMinX+1.0)
	columns.LineMinY = window.DC.CursorPos.y
	columns.LineMaxY = window.DC.CursorPos.y

	// Clear data if columns count changed
	if len(columns.Columns) != 0 && int(len(columns.Columns)) != columns_count+1 {
		columns.Columns = columns.Columns[:0]
	}

	// Initialize default widths
	columns.IsFirstFrame = (len(columns.Columns) == 0)
	if len(columns.Columns) == 0 {
		for n := int(0); n < columns_count+1; n++ {
			var column ImGuiOldColumnData
			column.OffsetNorm = float(n) / (float)(columns_count)
			columns.Columns = append(columns.Columns, column)
		}
	}

	for n := int(0); n < int(columns_count); n++ {
		// Compute clipping rectangle
		var column = &columns.Columns[n]
		var clip_x1 = IM_ROUND(window.Pos.x + GetColumnOffset(n))
		var clip_x2 = IM_ROUND(window.Pos.x + GetColumnOffset(n+1) - 1.0)
		column.ClipRect = ImRect{ImVec2{clip_x1, -FLT_MAX}, ImVec2{clip_x2, +FLT_MAX}}
		column.ClipRect.ClipWithFull(window.ClipRect)
	}

	if columns.Count > 1 {
		columns.Splitter.Split(window.DrawList, 1+columns.Count)
		columns.Splitter.SetCurrentChannel(window.DrawList, 1)
		PushColumnClipRect(0)
	}

	// We don't generally store Indent.x inside ColumnsOffset because it may be manipulated by the user.
	var offset_0 = GetColumnOffset(columns.Current)
	var offset_1 = GetColumnOffset(columns.Current + 1)
	var width = offset_1 - offset_0
	PushItemWidth(width * 0.65)
	window.DC.ColumnsOffset.x = max(column_padding-window.WindowPadding.x, 0.0)
	window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x)
	window.WorkRect.Max.x = window.Pos.x + offset_1 - column_padding
}

// EndColumns close columns
func EndColumns() {
	g := GImGui
	window := GetCurrentWindow()
	var columns = window.DC.CurrentColumns
	IM_ASSERT(columns != nil)

	PopItemWidth()
	if columns.Count > 1 {
		PopClipRect()
		columns.Splitter.Merge(window.DrawList)
	}

	var flags = columns.Flags
	columns.LineMaxY = max(columns.LineMaxY, window.DC.CursorPos.y)
	window.DC.CursorPos.y = columns.LineMaxY
	if (flags & ImGuiOldColumnFlags_GrowParentContentsSize) == 0 {
		window.DC.CursorMaxPos.x = columns.HostCursorMaxPosX // Restore cursor max pos, as columns don't grow parent
	}

	// Draw columns borders and handle resize
	// The IsBeingResized flag ensure we preserve pre-resize columns width so back-and-forth are not lossy
	var is_being_resized = false
	if (flags&ImGuiOldColumnFlags_NoBorder) == 0 && !window.SkipItems {
		// We clip Y boundaries CPU side because very long triangles are mishandled by some GPU drivers.
		var y1 = max(columns.HostCursorPosY, window.ClipRect.Min.y)
		var y2 = min(window.DC.CursorPos.y, window.ClipRect.Max.y)
		var dragging_column int = -1
		for n := int(1); n < columns.Count; n++ {
			var column = &columns.Columns[n]
			var x = window.Pos.x + GetColumnOffset(n)
			var column_id = columns.ID + ImGuiID(n)
			var column_hit_hw = COLUMNS_HIT_RECT_HALF_WIDTH
			var column_hit_rect = ImRect{ImVec2{x - column_hit_hw, y1}, ImVec2{x + column_hit_hw, y2}}
			KeepAliveID(column_id)
			if IsClippedEx(&column_hit_rect, column_id, false) {
				continue
			}

			var hovered, held = false, false
			if (flags & ImGuiOldColumnFlags_NoResize) == 0 {
				ButtonBehavior(&column_hit_rect, column_id, &hovered, &held, 0)
				if hovered || held {
					g.MouseCursor = ImGuiMouseCursor_ResizeEW
				}
				if held && (column.Flags&ImGuiOldColumnFlags_NoResize) == 0 {
					dragging_column = n
				}
			}

			// Draw column
			var col ImU32
			if held {
				col = GetColorU32FromID(ImGuiCol_SeparatorActive, 1)
			} else if hovered {
				col = GetColorU32FromID(ImGuiCol_SeparatorHovered, 1)
			} else {
				col = GetColorU32FromID(ImGuiCol_Separator, 1)
			}
			var xi = IM_FLOOR(x)
			window.DrawList.AddLine(&ImVec2{xi, y1 + 1.0}, &ImVec2{xi, y2}, col, 1)
		}

		// Apply dragging after drawing the column lines, so our rendered lines are in sync with how items were displayed during the frame.
		if dragging_column != -1 {
			if !columns.IsBeingResized {
				for n := int(0); n < columns.Count+1; n++ {
					columns.Columns[n].OffsetNormBeforeResize = columns.Columns[n].OffsetNorm
				}
			}
			columns.IsBeingResized = true
			is_being_resized = true
			var x = GetDraggedColumnOffset(columns, dragging_column)
			SetColumnOffset(dragging_column, x)
		}
	}
	columns.IsBeingResized = is_being_resized

	window.WorkRect = window.ParentWorkRect
	window.ParentWorkRect = columns.HostBackupParentWorkRect
	window.DC.CurrentColumns = nil
	window.DC.ColumnsOffset.x = 0.0
	window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x)
}

func PushColumnClipRect(column_index int) {
	window := GetCurrentWindowRead()
	var columns = window.DC.CurrentColumns
	if column_index < 0 {
		column_index = columns.Current
	}

	var column = &columns.Columns[column_index]
	PushClipRect(column.ClipRect.Min, column.ClipRect.Max, false)
}

func PopColumnsBackground() {
	window := GetCurrentWindowRead()
	var columns = window.DC.CurrentColumns
	if columns.Count == 1 {
		return
	}

	// Optimization: avoid PopClipRect() + SetCurrentChannel()
	SetWindowClipRectBeforeSetChannel(window, &columns.HostBackupClipRect)
	columns.Splitter.SetCurrentChannel(window.DrawList, columns.Current+1)
}

func GetColumnsID(str_id string, columns_count int) ImGuiID {
	window := GetCurrentWindow()

	var id ImGuiID
	// Differentiate column ID with an arbitrary prefix for cases where users name their columns set the same as another widget.
	// In addition, when an identifier isn't explicitly provided we include the number of columns in the hash to make it uniquer.
	if str_id != "" {
		PushID(0x11223347)
		id = window.GetIDs(str_id)
	} else {
		PushID(0x11223347 + columns_count)
		id = window.GetIDs("columns")
	}
	PopID()

	return id
}

func FindOrCreateColumns(window *ImGuiWindow, id ImGuiID) *ImGuiOldColumns {
	// We have few columns per window so for now we don't need bother much with turning this into a faster lookup.
	for n := range window.ColumnsStorage {
		if window.ColumnsStorage[n].ID == id {
			return &window.ColumnsStorage[n]
		}
	}

	window.ColumnsStorage = append(window.ColumnsStorage, ImGuiOldColumns{})
	var columns = &window.ColumnsStorage[len(window.ColumnsStorage)-1]
	columns.ID = id
	return columns
}

func GetColumnOffsetFromNorm(columns *ImGuiOldColumns, offset_norm float) float {
	return offset_norm * (columns.OffMaxX - columns.OffMinX)
}

func GetColumnNormFromOffset(columns *ImGuiOldColumns, offset float) float {
	return offset / (columns.OffMaxX - columns.OffMinX)
}

// Legacy Columns API (prefer using Tables!)
// - You can also use SameLine(pos_x) to mimic simplified columns.

func Columns(columns_count int /*= 1*/, id string /*= L*/, border bool /*= true*/) {
	window := GetCurrentWindow()
	IM_ASSERT(columns_count >= 1)

	var flags = ImGuiOldColumnFlags_NoBorder
	if border {
		flags = 0
	}
	//flags |= ImGuiOldColumnFlags_NoPreserveWidths; // NB: Legacy behavior
	var columns = window.DC.CurrentColumns
	if columns != nil && columns.Count == columns_count && columns.Flags == flags {
		return
	}

	if columns != nil {
		EndColumns()
	}

	if columns_count != 1 {
		BeginColumns(id, columns_count, flags)
	}
}

// NextColumn next column, defaults to current row or next row if the current row is finished
func NextColumn() {
	window := GetCurrentWindow()
	if window.SkipItems || window.DC.CurrentColumns == nil {
		return
	}

	g := GImGui
	var columns = window.DC.CurrentColumns

	if columns.Count == 1 {
		window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x)
		IM_ASSERT(columns.Current == 0)
		return
	}

	// Next column
	columns.Current++
	if columns.Current == columns.Count {
		columns.Current = 0
	}

	PopItemWidth()

	// Optimization: avoid PopClipRect() + SetCurrentChannel() + PushClipRect()
	// (which would needlessly attempt to update commands in the wrong channel, then pop or overwrite them),
	var column = &columns.Columns[columns.Current]
	SetWindowClipRectBeforeSetChannel(window, &column.ClipRect)
	columns.Splitter.SetCurrentChannel(window.DrawList, columns.Current+1)

	var column_padding = g.Style.ItemSpacing.x
	columns.LineMaxY = max(columns.LineMaxY, window.DC.CursorPos.y)
	if columns.Current > 0 {
		// Columns 1+ ignore IndentX (by canceling it out)
		// FIXME-COLUMNS: Unnecessary, could be locked?
		window.DC.ColumnsOffset.x = GetColumnOffset(columns.Current) - window.DC.Indent.x + column_padding
	} else {
		// New row/line: column 0 honor IndentX.
		window.DC.ColumnsOffset.x = max(column_padding-window.WindowPadding.x, 0.0)
		columns.LineMinY = columns.LineMaxY
	}
	window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x)
	window.DC.CursorPos.y = columns.LineMinY
	window.DC.CurrLineSize = ImVec2{0.0, 0.0}
	window.DC.CurrLineTextBaseOffset = 0.0

	// FIXME-COLUMNS: Share code with BeginColumns() - move code on columns setup.
	var offset_0 = GetColumnOffset(columns.Current)
	var offset_1 = GetColumnOffset(columns.Current + 1)
	var width = offset_1 - offset_0
	PushItemWidth(width * 0.65)
	window.WorkRect.Max.x = window.Pos.x + offset_1 - column_padding
}

// GetColumnIndex get current column index
func GetColumnIndex() int {
	window := GetCurrentWindowRead()
	if window.DC.CurrentColumns != nil {
		return window.DC.CurrentColumns.Current
	}
	return 0
}

// GetColumnWidth get column width (in pixels). pass -1 to use current column
func GetColumnWidth(column_index int /*= -1*/) float {
	g := GImGui
	window := g.CurrentWindow
	var columns = window.DC.CurrentColumns
	if columns == nil {
		return GetContentRegionAvail().x
	}

	if column_index < 0 {
		column_index = columns.Current
	}
	return GetColumnOffsetFromNorm(columns, columns.Columns[column_index+1].OffsetNorm-columns.Columns[column_index].OffsetNorm)
}

// SetColumnWidth set column width (in pixels). pass -1 to use current column
func SetColumnWidth(column_index int, width float) {
	window := GetCurrentWindowRead()
	var columns = window.DC.CurrentColumns
	IM_ASSERT(columns != nil)

	if column_index < 0 {
		column_index = columns.Current
	}
	SetColumnOffset(column_index+1, GetColumnOffset(column_index)+width)
}

// GetColumnOffset get position of column line (in pixels, from the left side of the contents region). pass -1 to use current column, otherwise 0..GetColumnsCount() inclusive. column 0 is typically 0.0
func GetColumnOffset(column_index int /*= -1*/) float {
	window := GetCurrentWindowRead()
	var columns = window.DC.CurrentColumns
	if columns == nil {
		return 0.0
	}

	if column_index < 0 {
		column_index = columns.Current
	}
	IM_ASSERT(column_index < int(len(columns.Columns)))

	var t = columns.Columns[column_index].OffsetNorm
	var x_offset = ImLerp(columns.OffMinX, columns.OffMaxX, t)
	return x_offset
}

// SetColumnOffset set position of column line (in pixels, from the left side of the contents region). pass -1 to use current column
func SetColumnOffset(column_index int, offset float) {
	g := GImGui
	window := g.CurrentWindow
	var columns = window.DC.CurrentColumns
	IM_ASSERT(columns != nil)

	if column_index < 0 {
		column_index = columns.Current
	}
	IM_ASSERT(column_index < int(len(columns.Columns)))

	var preserve_width = (columns.Flags&ImGuiOldColumnFlags_NoPreserveWidths) == 0 && (column_index < columns.Count-1)
	var width float = 0.0
	if preserve_width {
		width = GetColumnWidthEx(columns, column_index, columns.IsBeingResized)
	}

	if (columns.Flags & ImGuiOldColumnFlags_NoForceWithinWindow) == 0 {
		offset = min(offset, columns.OffMaxX-g.Style.ColumnsMinSpacing*float(columns.Count-column_index))
	}
	columns.Columns[column_index].OffsetNorm = GetColumnNormFromOffset(columns, offset-columns.OffMinX)

	if preserve_width {
		SetColumnOffset(column_index+1, offset+max(g.Style.ColumnsMinSpacing, width))
	}
}

func GetColumnsCount() int {
	window := GetCurrentWindowRead()
	if window.DC.CurrentColumns != nil {
		return window.DC.CurrentColumns.Count
	}
	return 1
}

func GetDraggedColumnOffset(columns *ImGuiOldColumns, column_index int) float {
	// Active (dragged) column always follow mouse. The reason we need this is that dragging a column to the right edge of an auto-resizing
	// window creates a feedback loop because we store normalized positions. So while dragging we enforce absolute positioning.
	g := GImGui
	window := g.CurrentWindow
	IM_ASSERT(column_index > 0) // We are not supposed to drag column 0.
	IM_ASSERT(g.ActiveId == columns.ID+ImGuiID(column_index))

	var x = g.IO.MousePos.x - g.ActiveIdClickOffset.x + COLUMNS_HIT_RECT_HALF_WIDTH - window.Pos.x
	x = max(x, GetColumnOffset(column_index-1)+g.Style.ColumnsMinSpacing)
	if (columns.Flags & ImGuiOldColumnFlags_NoPreserveWidths) != 0 {
		x = min(x, GetColumnOffset(column_index+1)-g.Style.ColumnsMinSpacing)
	}

	return x
}

func GetColumnWidthEx(columns *ImGuiOldColumns, column_index int, before_resize bool) float {
	if column_index < 0 {
		column_index = columns.Current
	}

	var offset_norm float
	if before_resize {
		offset_norm = columns.Columns[column_index+1].OffsetNormBeforeResize - columns.Columns[column_index].OffsetNormBeforeResize
	} else {
		offset_norm = columns.Columns[column_index+1].OffsetNorm - columns.Columns[column_index].OffsetNorm
	}
	return GetColumnOffsetFromNorm(columns, offset_norm)
}
