

// FIXME-TABLE: This is a mess, need to redesign how we render borders (as some are also done in TableEndRow)
void ImGui::TableDrawBorders(ImGuiTable* table)
{
    ImGuiWindow* inner_window = table.InnerWindow;
    if (!table.OuterWindow.ClipRect.Overlaps(table.OuterRect))
        return;

    ImDrawList* inner_drawlist = inner_window.DrawList;
    table.DrawSplitter.SetCurrentChannel(inner_drawlist, TABLE_DRAW_CHANNEL_BG0);
    inner_drawlist.PushClipRect(table.Bg0ClipRectForDrawCmd.Min, table.Bg0ClipRectForDrawCmd.Max, false);

    // Draw inner border and resizing feedback
    const float border_size = TABLE_BORDER_SIZE;
    const float draw_y1 = table.InnerRect.Min.y;
    const float draw_y2_body = table.InnerRect.Max.y;
    const float draw_y2_head = table.IsUsingHeaders ? ImMin(table.InnerRect.Max.y, (table.FreezeRowsCount >= 1 ? table.InnerRect.Min.y : table.WorkRect.Min.y) + table.LastFirstRowHeight) : draw_y1;
    if (table.Flags & ImGuiTableFlags_BordersInnerV)
    {
        for (int order_n = 0; order_n < table.ColumnsCount; order_n++)
        {
            if (!(table.EnabledMaskByDisplayOrder & ((ImU64)1 << order_n)))
                continue;

            const int column_n = table.DisplayOrderToIndex[order_n];
            ImGuiTableColumn* column = &table.Columns[column_n];
            const bool is_hovered = (table.HoveredColumnBorder == column_n);
            const bool is_resized = (table.ResizedColumn == column_n) && (table.InstanceInteracted == table.InstanceCurrent);
            const bool is_resizable = (column.Flags & (ImGuiTableColumnFlags_NoResize | ImGuiTableColumnFlags_NoDirectResize_)) == 0;
            const bool is_frozen_separator = (table.FreezeColumnsCount == order_n + 1);
            if (column.MaxX > table.InnerClipRect.Max.x && !is_resized)
                continue;

            // Decide whether right-most column is visible
            if (column.NextEnabledColumn == -1 && !is_resizable)
                if ((table.Flags & ImGuiTableFlags_SizingMask_) != ImGuiTableFlags_SizingFixedSame || (table.Flags & ImGuiTableFlags_NoHostExtendX))
                    continue;
            if (column.MaxX <= column.ClipRect.Min.x) // FIXME-TABLE FIXME-STYLE: Assume BorderSize==1, this is problematic if we want to increase the border size..
                continue;

            // Draw in outer window so right-most column won't be clipped
            // Always draw full height border when being resized/hovered, or on the delimitation of frozen column scrolling.
            ImU32 col;
            float draw_y2;
            if (is_hovered || is_resized || is_frozen_separator)
            {
                draw_y2 = draw_y2_body;
                col = is_resized ? GetColorU32(ImGuiCol_SeparatorActive) : is_hovered ? GetColorU32(ImGuiCol_SeparatorHovered) : table.BorderColorStrong;
            }
            else
            {
                draw_y2 = (table.Flags & (ImGuiTableFlags_NoBordersInBody | ImGuiTableFlags_NoBordersInBodyUntilResize)) ? draw_y2_head : draw_y2_body;
                col = (table.Flags & (ImGuiTableFlags_NoBordersInBody | ImGuiTableFlags_NoBordersInBodyUntilResize)) ? table.BorderColorStrong : table.BorderColorLight;
            }

            if (draw_y2 > draw_y1)
                inner_drawlist.AddLine(ImVec2(column.MaxX, draw_y1), ImVec2(column.MaxX, draw_y2), col, border_size);
        }
    }

    // Draw outer border
    // FIXME: could use AddRect or explicit VLine/HLine helper?
    if (table.Flags & ImGuiTableFlags_BordersOuter)
    {
        // Display outer border offset by 1 which is a simple way to display it without adding an extra draw call
        // (Without the offset, in outer_window it would be rendered behind cells, because child windows are above their
        // parent. In inner_window, it won't reach out over scrollbars. Another weird solution would be to display part
        // of it in inner window, and the part that's over scrollbars in the outer window..)
        // Either solution currently won't allow us to use a larger border size: the border would clipped.
        const ImRect outer_border = table.OuterRect;
        const ImU32 outer_col = table.BorderColorStrong;
        if ((table.Flags & ImGuiTableFlags_BordersOuter) == ImGuiTableFlags_BordersOuter)
        {
            inner_drawlist.AddRect(outer_border.Min, outer_border.Max, outer_col, 0.0f, 0, border_size);
        }
        else if (table.Flags & ImGuiTableFlags_BordersOuterV)
        {
            inner_drawlist.AddLine(outer_border.Min, ImVec2(outer_border.Min.x, outer_border.Max.y), outer_col, border_size);
            inner_drawlist.AddLine(ImVec2(outer_border.Max.x, outer_border.Min.y), outer_border.Max, outer_col, border_size);
        }
        else if (table.Flags & ImGuiTableFlags_BordersOuterH)
        {
            inner_drawlist.AddLine(outer_border.Min, ImVec2(outer_border.Max.x, outer_border.Min.y), outer_col, border_size);
            inner_drawlist.AddLine(ImVec2(outer_border.Min.x, outer_border.Max.y), outer_border.Max, outer_col, border_size);
        }
    }
    if ((table.Flags & ImGuiTableFlags_BordersInnerH) && table.RowPosY2 < table.OuterRect.Max.y)
    {
        // Draw bottom-most row border
        const float border_y = table.RowPosY2;
        if (border_y >= table.BgClipRect.Min.y && border_y < table.BgClipRect.Max.y)
            inner_drawlist.AddLine(ImVec2(table.BorderX1, border_y), ImVec2(table.BorderX2, border_y), table.BorderColorLight, border_size);
    }

    inner_drawlist.PopClipRect();
}

//-------------------------------------------------------------------------
// [SECTION] Tables: Sorting
//-------------------------------------------------------------------------
// - TableGetSortSpecs()
// - TableFixColumnSortDirection() [Internal]
// - TableGetColumnNextSortDirection() [Internal]
// - TableSetColumnSortDirection() [Internal]
// - TableSortSpecsSanitize() [Internal]
// - TableSortSpecsBuild() [Internal]
//-------------------------------------------------------------------------

// Return NULL if no sort specs (most often when ImGuiTableFlags_Sortable is not set)
// You can sort your data again when 'SpecsChanged == true'. It will be true with sorting specs have changed since
// last call, or the first time.
// Lifetime: don't hold on this pointer over multiple frames or past any subsequent call to BeginTable()!
ImGuiTableSortSpecs* ImGui::TableGetSortSpecs()
{
    ImGuiContext& g = *GImGui;
    ImGuiTable* table = g.CurrentTable;
    IM_ASSERT(table != NULL);

    if (!(table.Flags & ImGuiTableFlags_Sortable))
        return NULL;

    // Require layout (in case TableHeadersRow() hasn't been called) as it may alter IsSortSpecsDirty in some paths.
    if (!table.IsLayoutLocked)
        TableUpdateLayout(table);

    TableSortSpecsBuild(table);

    return &table.SortSpecs;
}

static inline ImGuiSortDirection TableGetColumnAvailSortDirection(ImGuiTableColumn* column, int n)
{
    IM_ASSERT(n < column.SortDirectionsAvailCount);
    return (column.SortDirectionsAvailList >> (n << 1)) & 0x03;
}

// Fix sort direction if currently set on a value which is unavailable (e.g. activating NoSortAscending/NoSortDescending)
void ImGui::TableFixColumnSortDirection(ImGuiTable* table, ImGuiTableColumn* column)
{
    if (column.SortOrder == -1 || (column.SortDirectionsAvailMask & (1 << column.SortDirection)) != 0)
        return;
    column.SortDirection = (ImU8)TableGetColumnAvailSortDirection(column, 0);
    table.IsSortSpecsDirty = true;
}

// Calculate next sort direction that would be set after clicking the column
// - If the PreferSortDescending flag is set, we will default to a Descending direction on the first click.
// - Note that the PreferSortAscending flag is never checked, it is essentially the default and therefore a no-op.
IM_STATIC_ASSERT(ImGuiSortDirection_None == 0 && ImGuiSortDirection_Ascending == 1 && ImGuiSortDirection_Descending == 2);
ImGuiSortDirection ImGui::TableGetColumnNextSortDirection(ImGuiTableColumn* column)
{
    IM_ASSERT(column.SortDirectionsAvailCount > 0);
    if (column.SortOrder == -1)
        return TableGetColumnAvailSortDirection(column, 0);
    for (int n = 0; n < 3; n++)
        if (column.SortDirection == TableGetColumnAvailSortDirection(column, n))
            return TableGetColumnAvailSortDirection(column, (n + 1) % column.SortDirectionsAvailCount);
    IM_ASSERT(0);
    return ImGuiSortDirection_None;
}

// Note that the NoSortAscending/NoSortDescending flags are processed in TableSortSpecsSanitize(), and they may change/revert
// the value of SortDirection. We could technically also do it here but it would be unnecessary and duplicate code.
void ImGui::TableSetColumnSortDirection(int column_n, ImGuiSortDirection sort_direction, bool append_to_sort_specs)
{
    ImGuiContext& g = *GImGui;
    ImGuiTable* table = g.CurrentTable;

    if (!(table.Flags & ImGuiTableFlags_SortMulti))
        append_to_sort_specs = false;
    if (!(table.Flags & ImGuiTableFlags_SortTristate))
        IM_ASSERT(sort_direction != ImGuiSortDirection_None);

    ImGuiTableColumnIdx sort_order_max = 0;
    if (append_to_sort_specs)
        for (int other_column_n = 0; other_column_n < table.ColumnsCount; other_column_n++)
            sort_order_max = ImMax(sort_order_max, table.Columns[other_column_n].SortOrder);

    ImGuiTableColumn* column = &table.Columns[column_n];
    column.SortDirection = (ImU8)sort_direction;
    if (column.SortDirection == ImGuiSortDirection_None)
        column.SortOrder = -1;
    else if (column.SortOrder == -1 || !append_to_sort_specs)
        column.SortOrder = append_to_sort_specs ? sort_order_max + 1 : 0;

    for (int other_column_n = 0; other_column_n < table.ColumnsCount; other_column_n++)
    {
        ImGuiTableColumn* other_column = &table.Columns[other_column_n];
        if (other_column != column && !append_to_sort_specs)
            other_column.SortOrder = -1;
        TableFixColumnSortDirection(table, other_column);
    }
    table.IsSettingsDirty = true;
    table.IsSortSpecsDirty = true;
}

void ImGui::TableSortSpecsSanitize(ImGuiTable* table)
{
    IM_ASSERT(table.Flags & ImGuiTableFlags_Sortable);

    // Clear SortOrder from hidden column and verify that there's no gap or duplicate.
    int sort_order_count = 0;
    ImU64 sort_order_mask = 0x00;
    for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
    {
        ImGuiTableColumn* column = &table.Columns[column_n];
        if (column.SortOrder != -1 && !column.IsEnabled)
            column.SortOrder = -1;
        if (column.SortOrder == -1)
            continue;
        sort_order_count++;
        sort_order_mask |= ((ImU64)1 << column.SortOrder);
        IM_ASSERT(sort_order_count < (int)sizeof(sort_order_mask) * 8);
    }

    const bool need_fix_linearize = ((ImU64)1 << sort_order_count) != (sort_order_mask + 1);
    const bool need_fix_single_sort_order = (sort_order_count > 1) && !(table.Flags & ImGuiTableFlags_SortMulti);
    if (need_fix_linearize || need_fix_single_sort_order)
    {
        ImU64 fixed_mask = 0x00;
        for (int sort_n = 0; sort_n < sort_order_count; sort_n++)
        {
            // Fix: Rewrite sort order fields if needed so they have no gap or duplicate.
            // (e.g. SortOrder 0 disappeared, SortOrder 1..2 exists -. rewrite then as SortOrder 0..1)
            int column_with_smallest_sort_order = -1;
            for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
                if ((fixed_mask & ((ImU64)1 << (ImU64)column_n)) == 0 && table.Columns[column_n].SortOrder != -1)
                    if (column_with_smallest_sort_order == -1 || table.Columns[column_n].SortOrder < table.Columns[column_with_smallest_sort_order].SortOrder)
                        column_with_smallest_sort_order = column_n;
            IM_ASSERT(column_with_smallest_sort_order != -1);
            fixed_mask |= ((ImU64)1 << column_with_smallest_sort_order);
            table.Columns[column_with_smallest_sort_order].SortOrder = (ImGuiTableColumnIdx)sort_n;

            // Fix: Make sure only one column has a SortOrder if ImGuiTableFlags_MultiSortable is not set.
            if (need_fix_single_sort_order)
            {
                sort_order_count = 1;
                for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
                    if (column_n != column_with_smallest_sort_order)
                        table.Columns[column_n].SortOrder = -1;
                break;
            }
        }
    }

    // Fallback default sort order (if no column had the ImGuiTableColumnFlags_DefaultSort flag)
    if (sort_order_count == 0 && !(table.Flags & ImGuiTableFlags_SortTristate))
        for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
        {
            ImGuiTableColumn* column = &table.Columns[column_n];
            if (column.IsEnabled && !(column.Flags & ImGuiTableColumnFlags_NoSort))
            {
                sort_order_count = 1;
                column.SortOrder = 0;
                column.SortDirection = (ImU8)TableGetColumnAvailSortDirection(column, 0);
                break;
            }
        }

    table.SortSpecsCount = (ImGuiTableColumnIdx)sort_order_count;
}

void ImGui::TableSortSpecsBuild(ImGuiTable* table)
{
    bool dirty = table.IsSortSpecsDirty;
    if (dirty)
    {
        TableSortSpecsSanitize(table);
        table.SortSpecsMulti.resize(table.SortSpecsCount <= 1 ? 0 : table.SortSpecsCount);
        table.SortSpecs.SpecsDirty = true; // Mark as dirty for user
        table.IsSortSpecsDirty = false; // Mark as not dirty for us
    }

    // Write output
    ImGuiTableColumnSortSpecs* sort_specs = (table.SortSpecsCount == 0) ? NULL : (table.SortSpecsCount == 1) ? &table.SortSpecsSingle : table.SortSpecsMulti.Data;
    if (dirty && sort_specs != NULL)
        for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
        {
            ImGuiTableColumn* column = &table.Columns[column_n];
            if (column.SortOrder == -1)
                continue;
            IM_ASSERT(column.SortOrder < table.SortSpecsCount);
            ImGuiTableColumnSortSpecs* sort_spec = &sort_specs[column.SortOrder];
            sort_spec.ColumnUserID = column.UserID;
            sort_spec.ColumnIndex = (ImGuiTableColumnIdx)column_n;
            sort_spec.SortOrder = (ImGuiTableColumnIdx)column.SortOrder;
            sort_spec.SortDirection = column.SortDirection;
        }

    table.SortSpecs.Specs = sort_specs;
    table.SortSpecs.SpecsCount = table.SortSpecsCount;
}

//-------------------------------------------------------------------------
// [SECTION] Tables: Headers
//-------------------------------------------------------------------------
// - TableGetHeaderRowHeight() [Internal]
// - TableHeadersRow()
// - TableHeader()
//-------------------------------------------------------------------------

float ImGui::TableGetHeaderRowHeight()
{
    // Caring for a minor edge case:
    // Calculate row height, for the unlikely case that some labels may be taller than others.
    // If we didn't do that, uneven header height would highlight but smaller one before the tallest wouldn't catch input for all height.
    // In your custom header row you may omit this all together and just call TableNextRow() without a height...
    float row_height = GetTextLineHeight();
    int columns_count = TableGetColumnCount();
    for (int column_n = 0; column_n < columns_count; column_n++)
    {
        ImGuiTableColumnFlags flags = TableGetColumnFlags(column_n);
        if ((flags & ImGuiTableColumnFlags_IsEnabled) && !(flags & ImGuiTableColumnFlags_NoHeaderLabel))
            row_height = ImMax(row_height, CalcTextSize(TableGetColumnName(column_n)).y);
    }
    row_height += GetStyle().CellPadding.y * 2.0f;
    return row_height;
}

// [Public] This is a helper to output TableHeader() calls based on the column names declared in TableSetupColumn().
// The intent is that advanced users willing to create customized headers would not need to use this helper
// and can create their own! For example: TableHeader() may be preceeded by Checkbox() or other custom widgets.
// See 'Demo.Tables.Custom headers' for a demonstration of implementing a custom version of this.
// This code is constructed to not make much use of internal functions, as it is intended to be a template to copy.
// FIXME-TABLE: TableOpenContextMenu() and TableGetHeaderRowHeight() are not public.
void ImGui::TableHeadersRow()
{
    ImGuiContext& g = *GImGui;
    ImGuiTable* table = g.CurrentTable;
    IM_ASSERT(table != NULL && "Need to call TableHeadersRow() after BeginTable()!");

    // Layout if not already done (this is automatically done by TableNextRow, we do it here solely to facilitate stepping in debugger as it is frequent to step in TableUpdateLayout)
    if (!table.IsLayoutLocked)
        TableUpdateLayout(table);

    // Open row
    const float row_y1 = GetCursorScreenPos().y;
    const float row_height = TableGetHeaderRowHeight();
    TableNextRow(ImGuiTableRowFlags_Headers, row_height);
    if (table.HostSkipItems) // Merely an optimization, you may skip in your own code.
        return;

    const int columns_count = TableGetColumnCount();
    for (int column_n = 0; column_n < columns_count; column_n++)
    {
        if (!TableSetColumnIndex(column_n))
            continue;

        // Push an id to allow unnamed labels (generally accidental, but let's behave nicely with them)
        // - in your own code you may omit the PushID/PopID all-together, provided you know they won't collide
        // - table.InstanceCurrent is only >0 when we use multiple BeginTable/EndTable calls with same identifier.
        const char* name = (TableGetColumnFlags(column_n) & ImGuiTableColumnFlags_NoHeaderLabel) ? "" : TableGetColumnName(column_n);
        PushID(table.InstanceCurrent * table.ColumnsCount + column_n);
        TableHeader(name);
        PopID();
    }

    // Allow opening popup from the right-most section after the last column.
    ImVec2 mouse_pos = ImGui::GetMousePos();
    if (IsMouseReleased(1) && TableGetHoveredColumn() == columns_count)
        if (mouse_pos.y >= row_y1 && mouse_pos.y < row_y1 + row_height)
            TableOpenContextMenu(-1); // Will open a non-column-specific popup.
}

// Emit a column header (text + optional sort order)
// We cpu-clip text here so that all columns headers can be merged into a same draw call.
// Note that because of how we cpu-clip and display sorting indicators, you _cannot_ use SameLine() after a TableHeader()
void ImGui::TableHeader(const char* label)
{
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = g.CurrentWindow;
    if (window.SkipItems)
        return;

    ImGuiTable* table = g.CurrentTable;
    IM_ASSERT(table != NULL && "Need to call TableHeader() after BeginTable()!");
    IM_ASSERT(table.CurrentColumn != -1);
    const int column_n = table.CurrentColumn;
    ImGuiTableColumn* column = &table.Columns[column_n];

    // Label
    if (label == NULL)
        label = "";
    const char* label_end = FindRenderedTextEnd(label);
    ImVec2 label_size = CalcTextSize(label, label_end, true);
    ImVec2 label_pos = window.DC.CursorPos;

    // If we already got a row height, there's use that.
    // FIXME-TABLE: Padding problem if the correct outer-padding CellBgRect strays off our ClipRect?
    ImRect cell_r = TableGetCellBgRect(table, column_n);
    float label_height = ImMax(label_size.y, table.RowMinHeight - table.CellPaddingY * 2.0f);

    // Calculate ideal size for sort order arrow
    float w_arrow = 0.0f;
    float w_sort_text = 0.0f;
    char sort_order_suf[4] = "";
    const float ARROW_SCALE = 0.65f;
    if ((table.Flags & ImGuiTableFlags_Sortable) && !(column.Flags & ImGuiTableColumnFlags_NoSort))
    {
        w_arrow = ImFloor(g.FontSize * ARROW_SCALE + g.Style.FramePadding.x);
        if (column.SortOrder > 0)
        {
            ImFormatString(sort_order_suf, IM_ARRAYSIZE(sort_order_suf), "%d", column.SortOrder + 1);
            w_sort_text = g.Style.ItemInnerSpacing.x + CalcTextSize(sort_order_suf).x;
        }
    }

    // We feed our unclipped width to the column without writing on CursorMaxPos, so that column is still considering for merging.
    float max_pos_x = label_pos.x + label_size.x + w_sort_text + w_arrow;
    column.ContentMaxXHeadersUsed = ImMax(column.ContentMaxXHeadersUsed, column.WorkMaxX);
    column.ContentMaxXHeadersIdeal = ImMax(column.ContentMaxXHeadersIdeal, max_pos_x);

    // Keep header highlighted when context menu is open.
    const bool selected = (table.IsContextPopupOpen && table.ContextPopupColumn == column_n && table.InstanceInteracted == table.InstanceCurrent);
    ImGuiID id = window.GetID(label);
    ImRect bb(cell_r.Min.x, cell_r.Min.y, cell_r.Max.x, ImMax(cell_r.Max.y, cell_r.Min.y + label_height + g.Style.CellPadding.y * 2.0f));
    ItemSize(ImVec2(0.0f, label_height)); // Don't declare unclipped width, it'll be fed ContentMaxPosHeadersIdeal
    if (!ItemAdd(bb, id))
        return;

    //GetForegroundDrawList().AddRect(cell_r.Min, cell_r.Max, IM_COL32(255, 0, 0, 255)); // [DEBUG]
    //GetForegroundDrawList().AddRect(bb.Min, bb.Max, IM_COL32(255, 0, 0, 255)); // [DEBUG]

    // Using AllowItemOverlap mode because we cover the whole cell, and we want user to be able to submit subsequent items.
    bool hovered, held;
    bool pressed = ButtonBehavior(bb, id, &hovered, &held, ImGuiButtonFlags_AllowItemOverlap);
    if (g.ActiveId != id)
        SetItemAllowOverlap();
    if (held || hovered || selected)
    {
        const ImU32 col = GetColorU32(held ? ImGuiCol_HeaderActive : hovered ? ImGuiCol_HeaderHovered : ImGuiCol_Header);
        //RenderFrame(bb.Min, bb.Max, col, false, 0.0f);
        TableSetBgColor(ImGuiTableBgTarget_CellBg, col, table.CurrentColumn);
    }
    else
    {
        // Submit single cell bg color in the case we didn't submit a full header row
        if ((table.RowFlags & ImGuiTableRowFlags_Headers) == 0)
            TableSetBgColor(ImGuiTableBgTarget_CellBg, GetColorU32(ImGuiCol_TableHeaderBg), table.CurrentColumn);
    }
    RenderNavHighlight(bb, id, ImGuiNavHighlightFlags_TypeThin | ImGuiNavHighlightFlags_NoRounding);
    if (held)
        table.HeldHeaderColumn = (ImGuiTableColumnIdx)column_n;
    window.DC.CursorPos.y -= g.Style.ItemSpacing.y * 0.5f;

    // Drag and drop to re-order columns.
    // FIXME-TABLE: Scroll request while reordering a column and it lands out of the scrolling zone.
    if (held && (table.Flags & ImGuiTableFlags_Reorderable) && IsMouseDragging(0) && !g.DragDropActive)
    {
        // While moving a column it will jump on the other side of the mouse, so we also test for MouseDelta.x
        table.ReorderColumn = (ImGuiTableColumnIdx)column_n;
        table.InstanceInteracted = table.InstanceCurrent;

        // We don't reorder: through the frozen<>unfrozen line, or through a column that is marked with ImGuiTableColumnFlags_NoReorder.
        if (g.IO.MouseDelta.x < 0.0f && g.IO.MousePos.x < cell_r.Min.x)
            if (ImGuiTableColumn* prev_column = (column.PrevEnabledColumn != -1) ? &table.Columns[column.PrevEnabledColumn] : NULL)
                if (!((column.Flags | prev_column.Flags) & ImGuiTableColumnFlags_NoReorder))
                    if ((column.IndexWithinEnabledSet < table.FreezeColumnsRequest) == (prev_column.IndexWithinEnabledSet < table.FreezeColumnsRequest))
                        table.ReorderColumnDir = -1;
        if (g.IO.MouseDelta.x > 0.0f && g.IO.MousePos.x > cell_r.Max.x)
            if (ImGuiTableColumn* next_column = (column.NextEnabledColumn != -1) ? &table.Columns[column.NextEnabledColumn] : NULL)
                if (!((column.Flags | next_column.Flags) & ImGuiTableColumnFlags_NoReorder))
                    if ((column.IndexWithinEnabledSet < table.FreezeColumnsRequest) == (next_column.IndexWithinEnabledSet < table.FreezeColumnsRequest))
                        table.ReorderColumnDir = +1;
    }

    // Sort order arrow
    const float ellipsis_max = cell_r.Max.x - w_arrow - w_sort_text;
    if ((table.Flags & ImGuiTableFlags_Sortable) && !(column.Flags & ImGuiTableColumnFlags_NoSort))
    {
        if (column.SortOrder != -1)
        {
            float x = ImMax(cell_r.Min.x, cell_r.Max.x - w_arrow - w_sort_text);
            float y = label_pos.y;
            if (column.SortOrder > 0)
            {
                PushStyleColor(ImGuiCol_Text, GetColorU32(ImGuiCol_Text, 0.70f));
                RenderText(ImVec2(x + g.Style.ItemInnerSpacing.x, y), sort_order_suf);
                PopStyleColor();
                x += w_sort_text;
            }
            RenderArrow(window.DrawList, ImVec2(x, y), GetColorU32(ImGuiCol_Text), column.SortDirection == ImGuiSortDirection_Ascending ? ImGuiDir_Up : ImGuiDir_Down, ARROW_SCALE);
        }

        // Handle clicking on column header to adjust Sort Order
        if (pressed && table.ReorderColumn != column_n)
        {
            ImGuiSortDirection sort_direction = TableGetColumnNextSortDirection(column);
            TableSetColumnSortDirection(column_n, sort_direction, g.IO.KeyShift);
        }
    }

    // Render clipped label. Clipping here ensure that in the majority of situations, all our header cells will
    // be merged into a single draw call.
    //window.DrawList.AddCircleFilled(ImVec2(ellipsis_max, label_pos.y), 40, IM_COL32_WHITE);
    RenderTextEllipsis(window.DrawList, label_pos, ImVec2(ellipsis_max, label_pos.y + label_height + g.Style.FramePadding.y), ellipsis_max, ellipsis_max, label, label_end, &label_size);

    const bool text_clipped = label_size.x > (ellipsis_max - label_pos.x);
    if (text_clipped && hovered && g.HoveredIdNotActiveTimer > g.TooltipSlowDelay)
        SetTooltip("%.*s", (int)(label_end - label), label);

    // We don't use BeginPopupContextItem() because we want the popup to stay up even after the column is hidden
    if (IsMouseReleased(1) && IsItemHovered())
        TableOpenContextMenu(column_n);
}

//-------------------------------------------------------------------------
// [SECTION] Tables: Context Menu
//-------------------------------------------------------------------------
// - TableOpenContextMenu() [Internal]
// - TableDrawContextMenu() [Internal]
//-------------------------------------------------------------------------

// Use -1 to open menu not specific to a given column.
void ImGui::TableOpenContextMenu(int column_n)
{
    ImGuiContext& g = *GImGui;
    ImGuiTable* table = g.CurrentTable;
    if (column_n == -1 && table.CurrentColumn != -1)   // When called within a column automatically use this one (for consistency)
        column_n = table.CurrentColumn;
    if (column_n == table.ColumnsCount)                // To facilitate using with TableGetHoveredColumn()
        column_n = -1;
    IM_ASSERT(column_n >= -1 && column_n < table.ColumnsCount);
    if (table.Flags & (ImGuiTableFlags_Resizable | ImGuiTableFlags_Reorderable | ImGuiTableFlags_Hideable))
    {
        table.IsContextPopupOpen = true;
        table.ContextPopupColumn = (ImGuiTableColumnIdx)column_n;
        table.InstanceInteracted = table.InstanceCurrent;
        const ImGuiID context_menu_id = ImHashStr("##ContextMenu", 0, table.ID);
        OpenPopupEx(context_menu_id, ImGuiPopupFlags_None);
    }
}

// Output context menu into current window (generally a popup)
// FIXME-TABLE: Ideally this should be writable by the user. Full programmatic access to that data?
void ImGui::TableDrawContextMenu(ImGuiTable* table)
{
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = g.CurrentWindow;
    if (window.SkipItems)
        return;

    bool want_separator = false;
    const int column_n = (table.ContextPopupColumn >= 0 && table.ContextPopupColumn < table.ColumnsCount) ? table.ContextPopupColumn : -1;
    ImGuiTableColumn* column = (column_n != -1) ? &table.Columns[column_n] : NULL;

    // Sizing
    if (table.Flags & ImGuiTableFlags_Resizable)
    {
        if (column != NULL)
        {
            const bool can_resize = !(column.Flags & ImGuiTableColumnFlags_NoResize) && column.IsEnabled;
            if (MenuItem("Size column to fit###SizeOne", NULL, false, can_resize))
                TableSetColumnWidthAutoSingle(table, column_n);
        }

        const char* size_all_desc;
        if (table.ColumnsEnabledFixedCount == table.ColumnsEnabledCount && (table.Flags & ImGuiTableFlags_SizingMask_) != ImGuiTableFlags_SizingFixedSame)
            size_all_desc = "Size all columns to fit###SizeAll";        // All fixed
        else
            size_all_desc = "Size all columns to default###SizeAll";    // All stretch or mixed
        if (MenuItem(size_all_desc, NULL))
            TableSetColumnWidthAutoAll(table);
        want_separator = true;
    }

    // Ordering
    if (table.Flags & ImGuiTableFlags_Reorderable)
    {
        if (MenuItem("Reset order", NULL, false, !table.IsDefaultDisplayOrder))
            table.IsResetDisplayOrderRequest = true;
        want_separator = true;
    }

    // Reset all (should work but seems unnecessary/noisy to expose?)
    //if (MenuItem("Reset all"))
    //    table.IsResetAllRequest = true;

    // Sorting
    // (modify TableOpenContextMenu() to add _Sortable flag if enabling this)
#if 0
    if ((table.Flags & ImGuiTableFlags_Sortable) && column != NULL && (column.Flags & ImGuiTableColumnFlags_NoSort) == 0)
    {
        if (want_separator)
            Separator();
        want_separator = true;

        bool append_to_sort_specs = g.IO.KeyShift;
        if (MenuItem("Sort in Ascending Order", NULL, column.SortOrder != -1 && column.SortDirection == ImGuiSortDirection_Ascending, (column.Flags & ImGuiTableColumnFlags_NoSortAscending) == 0))
            TableSetColumnSortDirection(table, column_n, ImGuiSortDirection_Ascending, append_to_sort_specs);
        if (MenuItem("Sort in Descending Order", NULL, column.SortOrder != -1 && column.SortDirection == ImGuiSortDirection_Descending, (column.Flags & ImGuiTableColumnFlags_NoSortDescending) == 0))
            TableSetColumnSortDirection(table, column_n, ImGuiSortDirection_Descending, append_to_sort_specs);
    }
#endif

    // Hiding / Visibility
    if (table.Flags & ImGuiTableFlags_Hideable)
    {
        if (want_separator)
            Separator();
        want_separator = true;

        PushItemFlag(ImGuiItemFlags_SelectableDontClosePopup, true);
        for (int other_column_n = 0; other_column_n < table.ColumnsCount; other_column_n++)
        {
            ImGuiTableColumn* other_column = &table.Columns[other_column_n];
            if (other_column.Flags & ImGuiTableColumnFlags_Disabled)
                continue;

            const char* name = TableGetColumnName(table, other_column_n);
            if (name == NULL || name[0] == 0)
                name = "<Unknown>";

            // Make sure we can't hide the last active column
            bool menu_item_active = (other_column.Flags & ImGuiTableColumnFlags_NoHide) ? false : true;
            if (other_column.IsUserEnabled && table.ColumnsEnabledCount <= 1)
                menu_item_active = false;
            if (MenuItem(name, NULL, other_column.IsUserEnabled, menu_item_active))
                other_column.IsUserEnabledNextFrame = !other_column.IsUserEnabled;
        }
        PopItemFlag();
    }
}

//-------------------------------------------------------------------------
// [SECTION] Tables: Settings (.ini data)
//-------------------------------------------------------------------------
// FIXME: The binding/finding/creating flow are too confusing.
//-------------------------------------------------------------------------
// - TableSettingsInit() [Internal]
// - TableSettingsCalcChunkSize() [Internal]
// - TableSettingsCreate() [Internal]
// - TableSettingsFindByID() [Internal]
// - TableGetBoundSettings() [Internal]
// - TableResetSettings()
// - TableSaveSettings() [Internal]
// - TableLoadSettings() [Internal]
// - TableSettingsHandler_ClearAll() [Internal]
// - TableSettingsHandler_ApplyAll() [Internal]
// - TableSettingsHandler_ReadOpen() [Internal]
// - TableSettingsHandler_ReadLine() [Internal]
// - TableSettingsHandler_WriteAll() [Internal]
// - TableSettingsInstallHandler() [Internal]
//-------------------------------------------------------------------------
// [Init] 1: TableSettingsHandler_ReadXXXX()   Load and parse .ini file into TableSettings.
// [Main] 2: TableLoadSettings()               When table is created, bind Table to TableSettings, serialize TableSettings data into Table.
// [Main] 3: TableSaveSettings()               When table properties are modified, serialize Table data into bound or new TableSettings, mark .ini as dirty.
// [Main] 4: TableSettingsHandler_WriteAll()   When .ini file is dirty (which can come from other source), save TableSettings into .ini file.
//-------------------------------------------------------------------------

// Clear and initialize empty settings instance
static void TableSettingsInit(ImGuiTableSettings* settings, ImGuiID id, int columns_count, int columns_count_max)
{
    IM_PLACEMENT_NEW(settings) ImGuiTableSettings();
    ImGuiTableColumnSettings* settings_column = settings.GetColumnSettings();
    for (int n = 0; n < columns_count_max; n++, settings_column++)
        IM_PLACEMENT_NEW(settings_column) ImGuiTableColumnSettings();
    settings.ID = id;
    settings.ColumnsCount = (ImGuiTableColumnIdx)columns_count;
    settings.ColumnsCountMax = (ImGuiTableColumnIdx)columns_count_max;
    settings.WantApply = true;
}

static size_t TableSettingsCalcChunkSize(int columns_count)
{
    return sizeof(ImGuiTableSettings) + (size_t)columns_count * sizeof(ImGuiTableColumnSettings);
}

ImGuiTableSettings* ImGui::TableSettingsCreate(ImGuiID id, int columns_count)
{
    ImGuiContext& g = *GImGui;
    ImGuiTableSettings* settings = g.SettingsTables.alloc_chunk(TableSettingsCalcChunkSize(columns_count));
    TableSettingsInit(settings, id, columns_count, columns_count);
    return settings;
}

// Find existing settings
ImGuiTableSettings* ImGui::TableSettingsFindByID(ImGuiID id)
{
    // FIXME-OPT: Might want to store a lookup map for this?
    ImGuiContext& g = *GImGui;
    for (ImGuiTableSettings* settings = g.SettingsTables.begin(); settings != NULL; settings = g.SettingsTables.next_chunk(settings))
        if (settings.ID == id)
            return settings;
    return NULL;
}

// Get settings for a given table, NULL if none
ImGuiTableSettings* ImGui::TableGetBoundSettings(ImGuiTable* table)
{
    if (table.SettingsOffset != -1)
    {
        ImGuiContext& g = *GImGui;
        ImGuiTableSettings* settings = g.SettingsTables.ptr_from_offset(table.SettingsOffset);
        IM_ASSERT(settings.ID == table.ID);
        if (settings.ColumnsCountMax >= table.ColumnsCount)
            return settings; // OK
        settings.ID = 0; // Invalidate storage, we won't fit because of a count change
    }
    return NULL;
}

// Restore initial state of table (with or without saved settings)
void ImGui::TableResetSettings(ImGuiTable* table)
{
    table.IsInitializing = table.IsSettingsDirty = true;
    table.IsResetAllRequest = false;
    table.IsSettingsRequestLoad = false;                   // Don't reload from ini
    table.SettingsLoadedFlags = ImGuiTableFlags_None;      // Mark as nothing loaded so our initialized data becomes authoritative
}

void ImGui::TableSaveSettings(ImGuiTable* table)
{
    table.IsSettingsDirty = false;
    if (table.Flags & ImGuiTableFlags_NoSavedSettings)
        return;

    // Bind or create settings data
    ImGuiContext& g = *GImGui;
    ImGuiTableSettings* settings = TableGetBoundSettings(table);
    if (settings == NULL)
    {
        settings = TableSettingsCreate(table.ID, table.ColumnsCount);
        table.SettingsOffset = g.SettingsTables.offset_from_ptr(settings);
    }
    settings.ColumnsCount = (ImGuiTableColumnIdx)table.ColumnsCount;

    // Serialize ImGuiTable/ImGuiTableColumn into ImGuiTableSettings/ImGuiTableColumnSettings
    IM_ASSERT(settings.ID == table.ID);
    IM_ASSERT(settings.ColumnsCount == table.ColumnsCount && settings.ColumnsCountMax >= settings.ColumnsCount);
    ImGuiTableColumn* column = table.Columns.Data;
    ImGuiTableColumnSettings* column_settings = settings.GetColumnSettings();

    bool save_ref_scale = false;
    settings.SaveFlags = ImGuiTableFlags_None;
    for (int n = 0; n < table.ColumnsCount; n++, column++, column_settings++)
    {
        const float width_or_weight = (column.Flags & ImGuiTableColumnFlags_WidthStretch) ? column.StretchWeight : column.WidthRequest;
        column_settings.WidthOrWeight = width_or_weight;
        column_settings.Index = (ImGuiTableColumnIdx)n;
        column_settings.DisplayOrder = column.DisplayOrder;
        column_settings.SortOrder = column.SortOrder;
        column_settings.SortDirection = column.SortDirection;
        column_settings.IsEnabled = column.IsUserEnabled;
        column_settings.IsStretch = (column.Flags & ImGuiTableColumnFlags_WidthStretch) ? 1 : 0;
        if ((column.Flags & ImGuiTableColumnFlags_WidthStretch) == 0)
            save_ref_scale = true;

        // We skip saving some data in the .ini file when they are unnecessary to restore our state.
        // Note that fixed width where initial width was derived from auto-fit will always be saved as InitStretchWeightOrWidth will be 0.0f.
        // FIXME-TABLE: We don't have logic to easily compare SortOrder to DefaultSortOrder yet so it's always saved when present.
        if (width_or_weight != column.InitStretchWeightOrWidth)
            settings.SaveFlags |= ImGuiTableFlags_Resizable;
        if (column.DisplayOrder != n)
            settings.SaveFlags |= ImGuiTableFlags_Reorderable;
        if (column.SortOrder != -1)
            settings.SaveFlags |= ImGuiTableFlags_Sortable;
        if (column.IsUserEnabled != ((column.Flags & ImGuiTableColumnFlags_DefaultHide) == 0))
            settings.SaveFlags |= ImGuiTableFlags_Hideable;
    }
    settings.SaveFlags &= table.Flags;
    settings.RefScale = save_ref_scale ? table.RefScale : 0.0f;

    MarkIniSettingsDirty();
}

void ImGui::TableLoadSettings(ImGuiTable* table)
{
    ImGuiContext& g = *GImGui;
    table.IsSettingsRequestLoad = false;
    if (table.Flags & ImGuiTableFlags_NoSavedSettings)
        return;

    // Bind settings
    ImGuiTableSettings* settings;
    if (table.SettingsOffset == -1)
    {
        settings = TableSettingsFindByID(table.ID);
        if (settings == NULL)
            return;
        if (settings.ColumnsCount != table.ColumnsCount) // Allow settings if columns count changed. We could otherwise decide to return...
            table.IsSettingsDirty = true;
        table.SettingsOffset = g.SettingsTables.offset_from_ptr(settings);
    }
    else
    {
        settings = TableGetBoundSettings(table);
    }

    table.SettingsLoadedFlags = settings.SaveFlags;
    table.RefScale = settings.RefScale;

    // Serialize ImGuiTableSettings/ImGuiTableColumnSettings into ImGuiTable/ImGuiTableColumn
    ImGuiTableColumnSettings* column_settings = settings.GetColumnSettings();
    ImU64 display_order_mask = 0;
    for (int data_n = 0; data_n < settings.ColumnsCount; data_n++, column_settings++)
    {
        int column_n = column_settings.Index;
        if (column_n < 0 || column_n >= table.ColumnsCount)
            continue;

        ImGuiTableColumn* column = &table.Columns[column_n];
        if (settings.SaveFlags & ImGuiTableFlags_Resizable)
        {
            if (column_settings.IsStretch)
                column.StretchWeight = column_settings.WidthOrWeight;
            else
                column.WidthRequest = column_settings.WidthOrWeight;
            column.AutoFitQueue = 0x00;
        }
        if (settings.SaveFlags & ImGuiTableFlags_Reorderable)
            column.DisplayOrder = column_settings.DisplayOrder;
        else
            column.DisplayOrder = (ImGuiTableColumnIdx)column_n;
        display_order_mask |= (ImU64)1 << column.DisplayOrder;
        column.IsUserEnabled = column.IsUserEnabledNextFrame = column_settings.IsEnabled;
        column.SortOrder = column_settings.SortOrder;
        column.SortDirection = column_settings.SortDirection;
    }

    // Validate and fix invalid display order data
    const ImU64 expected_display_order_mask = (settings.ColumnsCount == 64) ? ~0 : ((ImU64)1 << settings.ColumnsCount) - 1;
    if (display_order_mask != expected_display_order_mask)
        for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
            table.Columns[column_n].DisplayOrder = (ImGuiTableColumnIdx)column_n;

    // Rebuild index
    for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
        table.DisplayOrderToIndex[table.Columns[column_n].DisplayOrder] = (ImGuiTableColumnIdx)column_n;
}

static void TableSettingsHandler_ClearAll(ImGuiContext* ctx, ImGuiSettingsHandler*)
{
    ImGuiContext& g = *ctx;
    for (int i = 0; i != g.Tables.GetMapSize(); i++)
        if (ImGuiTable* table = g.Tables.TryGetMapData(i))
            table.SettingsOffset = -1;
    g.SettingsTables.clear();
}

// Apply to existing windows (if any)
static void TableSettingsHandler_ApplyAll(ImGuiContext* ctx, ImGuiSettingsHandler*)
{
    ImGuiContext& g = *ctx;
    for (int i = 0; i != g.Tables.GetMapSize(); i++)
        if (ImGuiTable* table = g.Tables.TryGetMapData(i))
        {
            table.IsSettingsRequestLoad = true;
            table.SettingsOffset = -1;
        }
}

static void* TableSettingsHandler_ReadOpen(ImGuiContext*, ImGuiSettingsHandler*, const char* name)
{
    ImGuiID id = 0;
    int columns_count = 0;
    if (sscanf(name, "0x%08X,%d", &id, &columns_count) < 2)
        return NULL;

    if (ImGuiTableSettings* settings = ImGui::TableSettingsFindByID(id))
    {
        if (settings.ColumnsCountMax >= columns_count)
        {
            TableSettingsInit(settings, id, columns_count, settings.ColumnsCountMax); // Recycle
            return settings;
        }
        settings.ID = 0; // Invalidate storage, we won't fit because of a count change
    }
    return ImGui::TableSettingsCreate(id, columns_count);
}

static void TableSettingsHandler_ReadLine(ImGuiContext*, ImGuiSettingsHandler*, void* entry, const char* line)
{
    // "Column 0  UserID=0x42AD2D21 Width=100 Visible=1 Order=0 Sort=0v"
    ImGuiTableSettings* settings = (ImGuiTableSettings*)entry;
    float f = 0.0f;
    int column_n = 0, r = 0, n = 0;

    if (sscanf(line, "RefScale=%f", &f) == 1) { settings.RefScale = f; return; }

    if (sscanf(line, "Column %d%n", &column_n, &r) == 1)
    {
        if (column_n < 0 || column_n >= settings.ColumnsCount)
            return;
        line = ImStrSkipBlank(line + r);
        char c = 0;
        ImGuiTableColumnSettings* column = settings.GetColumnSettings() + column_n;
        column.Index = (ImGuiTableColumnIdx)column_n;
        if (sscanf(line, "UserID=0x%08X%n", (ImU32*)&n, &r)==1) { line = ImStrSkipBlank(line + r); column.UserID = (ImGuiID)n; }
        if (sscanf(line, "Width=%d%n", &n, &r) == 1)            { line = ImStrSkipBlank(line + r); column.WidthOrWeight = (float)n; column.IsStretch = 0; settings.SaveFlags |= ImGuiTableFlags_Resizable; }
        if (sscanf(line, "Weight=%f%n", &f, &r) == 1)           { line = ImStrSkipBlank(line + r); column.WidthOrWeight = f; column.IsStretch = 1; settings.SaveFlags |= ImGuiTableFlags_Resizable; }
        if (sscanf(line, "Visible=%d%n", &n, &r) == 1)          { line = ImStrSkipBlank(line + r); column.IsEnabled = (ImU8)n; settings.SaveFlags |= ImGuiTableFlags_Hideable; }
        if (sscanf(line, "Order=%d%n", &n, &r) == 1)            { line = ImStrSkipBlank(line + r); column.DisplayOrder = (ImGuiTableColumnIdx)n; settings.SaveFlags |= ImGuiTableFlags_Reorderable; }
        if (sscanf(line, "Sort=%d%c%n", &n, &c, &r) == 2)       { line = ImStrSkipBlank(line + r); column.SortOrder = (ImGuiTableColumnIdx)n; column.SortDirection = (c == '^') ? ImGuiSortDirection_Descending : ImGuiSortDirection_Ascending; settings.SaveFlags |= ImGuiTableFlags_Sortable; }
    }
}

static void TableSettingsHandler_WriteAll(ImGuiContext* ctx, ImGuiSettingsHandler* handler, ImGuiTextBuffer* buf)
{
    ImGuiContext& g = *ctx;
    for (ImGuiTableSettings* settings = g.SettingsTables.begin(); settings != NULL; settings = g.SettingsTables.next_chunk(settings))
    {
        if (settings.ID == 0) // Skip ditched settings
            continue;

        // TableSaveSettings() may clear some of those flags when we establish that the data can be stripped
        // (e.g. Order was unchanged)
        const bool save_size    = (settings.SaveFlags & ImGuiTableFlags_Resizable) != 0;
        const bool save_visible = (settings.SaveFlags & ImGuiTableFlags_Hideable) != 0;
        const bool save_order   = (settings.SaveFlags & ImGuiTableFlags_Reorderable) != 0;
        const bool save_sort    = (settings.SaveFlags & ImGuiTableFlags_Sortable) != 0;
        if (!save_size && !save_visible && !save_order && !save_sort)
            continue;

        buf.reserve(buf.size() + 30 + settings.ColumnsCount * 50); // ballpark reserve
        buf.appendf("[%s][0x%08X,%d]\n", handler.TypeName, settings.ID, settings.ColumnsCount);
        if (settings.RefScale != 0.0f)
            buf.appendf("RefScale=%g\n", settings.RefScale);
        ImGuiTableColumnSettings* column = settings.GetColumnSettings();
        for (int column_n = 0; column_n < settings.ColumnsCount; column_n++, column++)
        {
            // "Column 0  UserID=0x42AD2D21 Width=100 Visible=1 Order=0 Sort=0v"
            bool save_column = column.UserID != 0 || save_size || save_visible || save_order || (save_sort && column.SortOrder != -1);
            if (!save_column)
                continue;
            buf.appendf("Column %-2d", column_n);
            if (column.UserID != 0)                    buf.appendf(" UserID=%08X", column.UserID);
            if (save_size && column.IsStretch)         buf.appendf(" Weight=%.4f", column.WidthOrWeight);
            if (save_size && !column.IsStretch)        buf.appendf(" Width=%d", (int)column.WidthOrWeight);
            if (save_visible)                           buf.appendf(" Visible=%d", column.IsEnabled);
            if (save_order)                             buf.appendf(" Order=%d", column.DisplayOrder);
            if (save_sort && column.SortOrder != -1)   buf.appendf(" Sort=%d%c", column.SortOrder, (column.SortDirection == ImGuiSortDirection_Ascending) ? 'v' : '^');
            buf.append("\n");
        }
        buf.append("\n");
    }
}

void ImGui::TableSettingsInstallHandler(ImGuiContext* context)
{
    ImGuiContext& g = *context;
    ImGuiSettingsHandler ini_handler;
    ini_handler.TypeName = "Table";
    ini_handler.TypeHash = ImHashStr("Table");
    ini_handler.ClearAllFn = TableSettingsHandler_ClearAll;
    ini_handler.ReadOpenFn = TableSettingsHandler_ReadOpen;
    ini_handler.ReadLineFn = TableSettingsHandler_ReadLine;
    ini_handler.ApplyAllFn = TableSettingsHandler_ApplyAll;
    ini_handler.WriteAllFn = TableSettingsHandler_WriteAll;
    g.SettingsHandlers.push_back(ini_handler);
}

//-------------------------------------------------------------------------
// [SECTION] Tables: Garbage Collection
//-------------------------------------------------------------------------
// - TableRemove() [Internal]
// - TableGcCompactTransientBuffers() [Internal]
// - TableGcCompactSettings() [Internal]
//-------------------------------------------------------------------------

// Remove Table (currently only used by TestEngine)
void ImGui::TableRemove(ImGuiTable* table)
{
    //IMGUI_DEBUG_LOG("TableRemove() id=0x%08X\n", table.ID);
    ImGuiContext& g = *GImGui;
    int table_idx = g.Tables.GetIndex(table);
    //memset(table.RawData.Data, 0, table.RawData.size_in_bytes());
    //memset(table, 0, sizeof(ImGuiTable));
    g.Tables.Remove(table.ID, table);
    g.TablesLastTimeActive[table_idx] = -1.0f;
}

// Free up/compact internal Table buffers for when it gets unused
void ImGui::TableGcCompactTransientBuffers(ImGuiTable* table)
{
    //IMGUI_DEBUG_LOG("TableGcCompactTransientBuffers() id=0x%08X\n", table.ID);
    ImGuiContext& g = *GImGui;
    IM_ASSERT(table.MemoryCompacted == false);
    table.SortSpecs.Specs = NULL;
    table.SortSpecsMulti.clear();
    table.IsSortSpecsDirty = true; // FIXME: shouldn't have to leak into user performing a sort
    table.ColumnsNames.clear();
    table.MemoryCompacted = true;
    for (int n = 0; n < table.ColumnsCount; n++)
        table.Columns[n].NameOffset = -1;
    g.TablesLastTimeActive[g.Tables.GetIndex(table)] = -1.0f;
}

void ImGui::TableGcCompactTransientBuffers(ImGuiTableTempData* temp_data)
{
    temp_data.DrawSplitter.ClearFreeMemory();
    temp_data.LastTimeActive = -1.0f;
}

// Compact and remove unused settings data (currently only used by TestEngine)
void ImGui::TableGcCompactSettings()
{
    ImGuiContext& g = *GImGui;
    int required_memory = 0;
    for (ImGuiTableSettings* settings = g.SettingsTables.begin(); settings != NULL; settings = g.SettingsTables.next_chunk(settings))
        if (settings.ID != 0)
            required_memory += (int)TableSettingsCalcChunkSize(settings.ColumnsCount);
    if (required_memory == g.SettingsTables.Buf.Size)
        return;
    ImChunkStream<ImGuiTableSettings> new_chunk_stream;
    new_chunk_stream.Buf.reserve(required_memory);
    for (ImGuiTableSettings* settings = g.SettingsTables.begin(); settings != NULL; settings = g.SettingsTables.next_chunk(settings))
        if (settings.ID != 0)
            memcpy(new_chunk_stream.alloc_chunk(TableSettingsCalcChunkSize(settings.ColumnsCount)), settings, TableSettingsCalcChunkSize(settings.ColumnsCount));
    g.SettingsTables.swap(new_chunk_stream);
}


//-------------------------------------------------------------------------
// [SECTION] Tables: Debugging
//-------------------------------------------------------------------------
// - DebugNodeTable() [Internal]
//-------------------------------------------------------------------------

#ifndef IMGUI_DISABLE_METRICS_WINDOW

static const char* DebugNodeTableGetSizingPolicyDesc(ImGuiTableFlags sizing_policy)
{
    sizing_policy &= ImGuiTableFlags_SizingMask_;
    if (sizing_policy == ImGuiTableFlags_SizingFixedFit)    { return "FixedFit"; }
    if (sizing_policy == ImGuiTableFlags_SizingFixedSame)   { return "FixedSame"; }
    if (sizing_policy == ImGuiTableFlags_SizingStretchProp) { return "StretchProp"; }
    if (sizing_policy == ImGuiTableFlags_SizingStretchSame) { return "StretchSame"; }
    return "N/A";
}

void ImGui::DebugNodeTable(ImGuiTable* table)
{
    char buf[512];
    char* p = buf;
    const char* buf_end = buf + IM_ARRAYSIZE(buf);
    const bool is_active = (table.LastFrameActive >= ImGui::GetFrameCount() - 2); // Note that fully clipped early out scrolling tables will appear as inactive here.
    ImFormatString(p, buf_end - p, "Table 0x%08X (%d columns, in '%s')%s", table.ID, table.ColumnsCount, table.OuterWindow.Name, is_active ? "" : " *Inactive*");
    if (!is_active) { PushStyleColor(ImGuiCol_Text, GetStyleColorVec4(ImGuiCol_TextDisabled)); }
    bool open = TreeNode(table, "%s", buf);
    if (!is_active) { PopStyleColor(); }
    if (IsItemHovered())
        GetForegroundDrawList().AddRect(table.OuterRect.Min, table.OuterRect.Max, IM_COL32(255, 255, 0, 255));
    if (IsItemVisible() && table.HoveredColumnBody != -1)
        GetForegroundDrawList().AddRect(GetItemRectMin(), GetItemRectMax(), IM_COL32(255, 255, 0, 255));
    if (!open)
        return;
    bool clear_settings = SmallButton("Clear settings");
    BulletText("OuterRect: Pos: (%.1f,%.1f) Size: (%.1f,%.1f) Sizing: '%s'", table.OuterRect.Min.x, table.OuterRect.Min.y, table.OuterRect.GetWidth(), table.OuterRect.GetHeight(), DebugNodeTableGetSizingPolicyDesc(table.Flags));
    BulletText("ColumnsGivenWidth: %.1f, ColumnsAutoFitWidth: %.1f, InnerWidth: %.1f%s", table.ColumnsGivenWidth, table.ColumnsAutoFitWidth, table.InnerWidth, table.InnerWidth == 0.0f ? " (auto)" : "");
    BulletText("CellPaddingX: %.1f, CellSpacingX: %.1f/%.1f, OuterPaddingX: %.1f", table.CellPaddingX, table.CellSpacingX1, table.CellSpacingX2, table.OuterPaddingX);
    BulletText("HoveredColumnBody: %d, HoveredColumnBorder: %d", table.HoveredColumnBody, table.HoveredColumnBorder);
    BulletText("ResizedColumn: %d, ReorderColumn: %d, HeldHeaderColumn: %d", table.ResizedColumn, table.ReorderColumn, table.HeldHeaderColumn);
    //BulletText("BgDrawChannels: %d/%d", 0, table.BgDrawChannelUnfrozen);
    float sum_weights = 0.0f;
    for (int n = 0; n < table.ColumnsCount; n++)
        if (table.Columns[n].Flags & ImGuiTableColumnFlags_WidthStretch)
            sum_weights += table.Columns[n].StretchWeight;
    for (int n = 0; n < table.ColumnsCount; n++)
    {
        ImGuiTableColumn* column = &table.Columns[n];
        const char* name = TableGetColumnName(table, n);
        ImFormatString(buf, IM_ARRAYSIZE(buf),
            "Column %d order %d '%s': offset %+.2f to %+.2f%s\n"
            "Enabled: %d, VisibleX/Y: %d/%d, RequestOutput: %d, SkipItems: %d, DrawChannels: %d,%d\n"
            "WidthGiven: %.1f, Request/Auto: %.1f/%.1f, StretchWeight: %.3f (%.1f%%)\n"
            "MinX: %.1f, MaxX: %.1f (%+.1f), ClipRect: %.1f to %.1f (+%.1f)\n"
            "ContentWidth: %.1f,%.1f, HeadersUsed/Ideal %.1f/%.1f\n"
            "Sort: %d%s, UserID: 0x%08X, Flags: 0x%04X: %s%s%s..",
            n, column.DisplayOrder, name, column.MinX - table.WorkRect.Min.x, column.MaxX - table.WorkRect.Min.x, (n < table.FreezeColumnsRequest) ? " (Frozen)" : "",
            column.IsEnabled, column.IsVisibleX, column.IsVisibleY, column.IsRequestOutput, column.IsSkipItems, column.DrawChannelFrozen, column.DrawChannelUnfrozen,
            column.WidthGiven, column.WidthRequest, column.WidthAuto, column.StretchWeight, column.StretchWeight > 0.0f ? (column.StretchWeight / sum_weights) * 100.0f : 0.0f,
            column.MinX, column.MaxX, column.MaxX - column.MinX, column.ClipRect.Min.x, column.ClipRect.Max.x, column.ClipRect.Max.x - column.ClipRect.Min.x,
            column.ContentMaxXFrozen - column.WorkMinX, column.ContentMaxXUnfrozen - column.WorkMinX, column.ContentMaxXHeadersUsed - column.WorkMinX, column.ContentMaxXHeadersIdeal - column.WorkMinX,
            column.SortOrder, (column.SortDirection == ImGuiSortDirection_Ascending) ? " (Asc)" : (column.SortDirection == ImGuiSortDirection_Descending) ? " (Des)" : "", column.UserID, column.Flags,
            (column.Flags & ImGuiTableColumnFlags_WidthStretch) ? "WidthStretch " : "",
            (column.Flags & ImGuiTableColumnFlags_WidthFixed) ? "WidthFixed " : "",
            (column.Flags & ImGuiTableColumnFlags_NoResize) ? "NoResize " : "");
        Bullet();
        Selectable(buf);
        if (IsItemHovered())
        {
            ImRect r(column.MinX, table.OuterRect.Min.y, column.MaxX, table.OuterRect.Max.y);
            GetForegroundDrawList().AddRect(r.Min, r.Max, IM_COL32(255, 255, 0, 255));
        }
    }
    if (ImGuiTableSettings* settings = TableGetBoundSettings(table))
        DebugNodeTableSettings(settings);
    if (clear_settings)
        table.IsResetAllRequest = true;
    TreePop();
}

void ImGui::DebugNodeTableSettings(ImGuiTableSettings* settings)
{
    if (!TreeNode((void*)(intptr_t)settings.ID, "Settings 0x%08X (%d columns)", settings.ID, settings.ColumnsCount))
        return;
    BulletText("SaveFlags: 0x%08X", settings.SaveFlags);
    BulletText("ColumnsCount: %d (max %d)", settings.ColumnsCount, settings.ColumnsCountMax);
    for (int n = 0; n < settings.ColumnsCount; n++)
    {
        ImGuiTableColumnSettings* column_settings = &settings.GetColumnSettings()[n];
        ImGuiSortDirection sort_dir = (column_settings.SortOrder != -1) ? (ImGuiSortDirection)column_settings.SortDirection : ImGuiSortDirection_None;
        BulletText("Column %d Order %d SortOrder %d %s Vis %d %s %7.3f UserID 0x%08X",
            n, column_settings.DisplayOrder, column_settings.SortOrder,
            (sort_dir == ImGuiSortDirection_Ascending) ? "Asc" : (sort_dir == ImGuiSortDirection_Descending) ? "Des" : "---",
            column_settings.IsEnabled, column_settings.IsStretch ? "Weight" : "Width ", column_settings.WidthOrWeight, column_settings.UserID);
    }
    TreePop();
}

#else // #ifndef IMGUI_DISABLE_METRICS_WINDOW

void ImGui::DebugNodeTable(ImGuiTable*) {}
void ImGui::DebugNodeTableSettings(ImGuiTableSettings*) {}

#endif


//-------------------------------------------------------------------------
// [SECTION] Columns, BeginColumns, EndColumns, etc.
// (This is a legacy API, prefer using BeginTable/EndTable!)
//-------------------------------------------------------------------------
// FIXME: sizing is lossy when columns width is very small (default width may turn negative etc.)
//-------------------------------------------------------------------------
// - SetWindowClipRectBeforeSetChannel() [Internal]
// - GetColumnIndex()
// - GetColumnsCount()
// - GetColumnOffset()
// - GetColumnWidth()
// - SetColumnOffset()
// - SetColumnWidth()
// - PushColumnClipRect() [Internal]
// - PushColumnsBackground() [Internal]
// - PopColumnsBackground() [Internal]
// - FindOrCreateColumns() [Internal]
// - GetColumnsID() [Internal]
// - BeginColumns()
// - NextColumn()
// - EndColumns()
// - Columns()
//-------------------------------------------------------------------------

// [Internal] Small optimization to avoid calls to PopClipRect/SetCurrentChannel/PushClipRect in sequences,
// they would meddle many times with the underlying ImDrawCmd.
// Instead, we do a preemptive overwrite of clipping rectangle _without_ altering the command-buffer and let
// the subsequent single call to SetCurrentChannel() does it things once.
void ImGui::SetWindowClipRectBeforeSetChannel(ImGuiWindow* window, const ImRect& clip_rect)
{
    ImVec4 clip_rect_vec4 = clip_rect.ToVec4();
    window.ClipRect = clip_rect;
    window.DrawList._CmdHeader.ClipRect = clip_rect_vec4;
    window.DrawList._ClipRectStack.Data[window.DrawList._ClipRectStack.Size - 1] = clip_rect_vec4;
}

int ImGui::GetColumnIndex()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CurrentColumns ? window.DC.CurrentColumns.Current : 0;
}

int ImGui::GetColumnsCount()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CurrentColumns ? window.DC.CurrentColumns.Count : 1;
}

float ImGui::GetColumnOffsetFromNorm(const ImGuiOldColumns* columns, float offset_norm)
{
    return offset_norm * (columns.OffMaxX - columns.OffMinX);
}

float ImGui::GetColumnNormFromOffset(const ImGuiOldColumns* columns, float offset)
{
    return offset / (columns.OffMaxX - columns.OffMinX);
}

static const float COLUMNS_HIT_RECT_HALF_WIDTH = 4.0f;

static float GetDraggedColumnOffset(ImGuiOldColumns* columns, int column_index)
{
    // Active (dragged) column always follow mouse. The reason we need this is that dragging a column to the right edge of an auto-resizing
    // window creates a feedback loop because we store normalized positions. So while dragging we enforce absolute positioning.
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = g.CurrentWindow;
    IM_ASSERT(column_index > 0); // We are not supposed to drag column 0.
    IM_ASSERT(g.ActiveId == columns.ID + ImGuiID(column_index));

    float x = g.IO.MousePos.x - g.ActiveIdClickOffset.x + COLUMNS_HIT_RECT_HALF_WIDTH - window.Pos.x;
    x = ImMax(x, ImGui::GetColumnOffset(column_index - 1) + g.Style.ColumnsMinSpacing);
    if ((columns.Flags & ImGuiOldColumnFlags_NoPreserveWidths))
        x = ImMin(x, ImGui::GetColumnOffset(column_index + 1) - g.Style.ColumnsMinSpacing);

    return x;
}

float ImGui::GetColumnOffset(int column_index)
{
    ImGuiWindow* window = GetCurrentWindowRead();
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    if (columns == NULL)
        return 0.0f;

    if (column_index < 0)
        column_index = columns.Current;
    IM_ASSERT(column_index < columns.Columns.Size);

    const float t = columns.Columns[column_index].OffsetNorm;
    const float x_offset = ImLerp(columns.OffMinX, columns.OffMaxX, t);
    return x_offset;
}

static float GetColumnWidthEx(ImGuiOldColumns* columns, int column_index, bool before_resize = false)
{
    if (column_index < 0)
        column_index = columns.Current;

    float offset_norm;
    if (before_resize)
        offset_norm = columns.Columns[column_index + 1].OffsetNormBeforeResize - columns.Columns[column_index].OffsetNormBeforeResize;
    else
        offset_norm = columns.Columns[column_index + 1].OffsetNorm - columns.Columns[column_index].OffsetNorm;
    return ImGui::GetColumnOffsetFromNorm(columns, offset_norm);
}

float ImGui::GetColumnWidth(int column_index)
{
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = g.CurrentWindow;
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    if (columns == NULL)
        return GetContentRegionAvail().x;

    if (column_index < 0)
        column_index = columns.Current;
    return GetColumnOffsetFromNorm(columns, columns.Columns[column_index + 1].OffsetNorm - columns.Columns[column_index].OffsetNorm);
}

void ImGui::SetColumnOffset(int column_index, float offset)
{
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = g.CurrentWindow;
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    IM_ASSERT(columns != NULL);

    if (column_index < 0)
        column_index = columns.Current;
    IM_ASSERT(column_index < columns.Columns.Size);

    const bool preserve_width = !(columns.Flags & ImGuiOldColumnFlags_NoPreserveWidths) && (column_index < columns.Count - 1);
    const float width = preserve_width ? GetColumnWidthEx(columns, column_index, columns.IsBeingResized) : 0.0f;

    if (!(columns.Flags & ImGuiOldColumnFlags_NoForceWithinWindow))
        offset = ImMin(offset, columns.OffMaxX - g.Style.ColumnsMinSpacing * (columns.Count - column_index));
    columns.Columns[column_index].OffsetNorm = GetColumnNormFromOffset(columns, offset - columns.OffMinX);

    if (preserve_width)
        SetColumnOffset(column_index + 1, offset + ImMax(g.Style.ColumnsMinSpacing, width));
}

void ImGui::SetColumnWidth(int column_index, float width)
{
    ImGuiWindow* window = GetCurrentWindowRead();
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    IM_ASSERT(columns != NULL);

    if (column_index < 0)
        column_index = columns.Current;
    SetColumnOffset(column_index + 1, GetColumnOffset(column_index) + width);
}

void ImGui::PushColumnClipRect(int column_index)
{
    ImGuiWindow* window = GetCurrentWindowRead();
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    if (column_index < 0)
        column_index = columns.Current;

    ImGuiOldColumnData* column = &columns.Columns[column_index];
    PushClipRect(column.ClipRect.Min, column.ClipRect.Max, false);
}

// Get into the columns background draw command (which is generally the same draw command as before we called BeginColumns)
void ImGui::PushColumnsBackground()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    if (columns.Count == 1)
        return;

    // Optimization: avoid SetCurrentChannel() + PushClipRect()
    columns.HostBackupClipRect = window.ClipRect;
    SetWindowClipRectBeforeSetChannel(window, columns.HostInitialClipRect);
    columns.Splitter.SetCurrentChannel(window.DrawList, 0);
}

void ImGui::PopColumnsBackground()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    if (columns.Count == 1)
        return;

    // Optimization: avoid PopClipRect() + SetCurrentChannel()
    SetWindowClipRectBeforeSetChannel(window, columns.HostBackupClipRect);
    columns.Splitter.SetCurrentChannel(window.DrawList, columns.Current + 1);
}

ImGuiOldColumns* ImGui::FindOrCreateColumns(ImGuiWindow* window, ImGuiID id)
{
    // We have few columns per window so for now we don't need bother much with turning this into a faster lookup.
    for (int n = 0; n < window.ColumnsStorage.Size; n++)
        if (window.ColumnsStorage[n].ID == id)
            return &window.ColumnsStorage[n];

    window.ColumnsStorage.push_back(ImGuiOldColumns());
    ImGuiOldColumns* columns = &window.ColumnsStorage.back();
    columns.ID = id;
    return columns;
}

ImGuiID ImGui::GetColumnsID(const char* str_id, int columns_count)
{
    ImGuiWindow* window = GetCurrentWindow();

    // Differentiate column ID with an arbitrary prefix for cases where users name their columns set the same as another widget.
    // In addition, when an identifier isn't explicitly provided we include the number of columns in the hash to make it uniquer.
    PushID(0x11223347 + (str_id ? 0 : columns_count));
    ImGuiID id = window.GetID(str_id ? str_id : "columns");
    PopID();

    return id;
}

void ImGui::BeginColumns(const char* str_id, int columns_count, ImGuiOldColumnFlags flags)
{
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = GetCurrentWindow();

    IM_ASSERT(columns_count >= 1);
    IM_ASSERT(window.DC.CurrentColumns == NULL);   // Nested columns are currently not supported

    // Acquire storage for the columns set
    ImGuiID id = GetColumnsID(str_id, columns_count);
    ImGuiOldColumns* columns = FindOrCreateColumns(window, id);
    IM_ASSERT(columns.ID == id);
    columns.Current = 0;
    columns.Count = columns_count;
    columns.Flags = flags;
    window.DC.CurrentColumns = columns;

    columns.HostCursorPosY = window.DC.CursorPos.y;
    columns.HostCursorMaxPosX = window.DC.CursorMaxPos.x;
    columns.HostInitialClipRect = window.ClipRect;
    columns.HostBackupParentWorkRect = window.ParentWorkRect;
    window.ParentWorkRect = window.WorkRect;

    // Set state for first column
    // We aim so that the right-most column will have the same clipping width as other after being clipped by parent ClipRect
    const float column_padding = g.Style.ItemSpacing.x;
    const float half_clip_extend_x = ImFloor(ImMax(window.WindowPadding.x * 0.5f, window.WindowBorderSize));
    const float max_1 = window.WorkRect.Max.x + column_padding - ImMax(column_padding - window.WindowPadding.x, 0.0f);
    const float max_2 = window.WorkRect.Max.x + half_clip_extend_x;
    columns.OffMinX = window.DC.Indent.x - column_padding + ImMax(column_padding - window.WindowPadding.x, 0.0f);
    columns.OffMaxX = ImMax(ImMin(max_1, max_2) - window.Pos.x, columns.OffMinX + 1.0f);
    columns.LineMinY = columns.LineMaxY = window.DC.CursorPos.y;

    // Clear data if columns count changed
    if (columns.Columns.Size != 0 && columns.Columns.Size != columns_count + 1)
        columns.Columns.resize(0);

    // Initialize default widths
    columns.IsFirstFrame = (columns.Columns.Size == 0);
    if (columns.Columns.Size == 0)
    {
        columns.Columns.reserve(columns_count + 1);
        for (int n = 0; n < columns_count + 1; n++)
        {
            ImGuiOldColumnData column;
            column.OffsetNorm = n / (float)columns_count;
            columns.Columns.push_back(column);
        }
    }

    for (int n = 0; n < columns_count; n++)
    {
        // Compute clipping rectangle
        ImGuiOldColumnData* column = &columns.Columns[n];
        float clip_x1 = IM_ROUND(window.Pos.x + GetColumnOffset(n));
        float clip_x2 = IM_ROUND(window.Pos.x + GetColumnOffset(n + 1) - 1.0f);
        column.ClipRect = ImRect(clip_x1, -FLT_MAX, clip_x2, +FLT_MAX);
        column.ClipRect.ClipWithFull(window.ClipRect);
    }

    if (columns.Count > 1)
    {
        columns.Splitter.Split(window.DrawList, 1 + columns.Count);
        columns.Splitter.SetCurrentChannel(window.DrawList, 1);
        PushColumnClipRect(0);
    }

    // We don't generally store Indent.x inside ColumnsOffset because it may be manipulated by the user.
    float offset_0 = GetColumnOffset(columns.Current);
    float offset_1 = GetColumnOffset(columns.Current + 1);
    float width = offset_1 - offset_0;
    PushItemWidth(width * 0.65f);
    window.DC.ColumnsOffset.x = ImMax(column_padding - window.WindowPadding.x, 0.0f);
    window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x);
    window.WorkRect.Max.x = window.Pos.x + offset_1 - column_padding;
}

void ImGui::NextColumn()
{
    ImGuiWindow* window = GetCurrentWindow();
    if (window.SkipItems || window.DC.CurrentColumns == NULL)
        return;

    ImGuiContext& g = *GImGui;
    ImGuiOldColumns* columns = window.DC.CurrentColumns;

    if (columns.Count == 1)
    {
        window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x);
        IM_ASSERT(columns.Current == 0);
        return;
    }

    // Next column
    if (++columns.Current == columns.Count)
        columns.Current = 0;

    PopItemWidth();

    // Optimization: avoid PopClipRect() + SetCurrentChannel() + PushClipRect()
    // (which would needlessly attempt to update commands in the wrong channel, then pop or overwrite them),
    ImGuiOldColumnData* column = &columns.Columns[columns.Current];
    SetWindowClipRectBeforeSetChannel(window, column.ClipRect);
    columns.Splitter.SetCurrentChannel(window.DrawList, columns.Current + 1);

    const float column_padding = g.Style.ItemSpacing.x;
    columns.LineMaxY = ImMax(columns.LineMaxY, window.DC.CursorPos.y);
    if (columns.Current > 0)
    {
        // Columns 1+ ignore IndentX (by canceling it out)
        // FIXME-COLUMNS: Unnecessary, could be locked?
        window.DC.ColumnsOffset.x = GetColumnOffset(columns.Current) - window.DC.Indent.x + column_padding;
    }
    else
    {
        // New row/line: column 0 honor IndentX.
        window.DC.ColumnsOffset.x = ImMax(column_padding - window.WindowPadding.x, 0.0f);
        columns.LineMinY = columns.LineMaxY;
    }
    window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x);
    window.DC.CursorPos.y = columns.LineMinY;
    window.DC.CurrLineSize = ImVec2(0.0f, 0.0f);
    window.DC.CurrLineTextBaseOffset = 0.0f;

    // FIXME-COLUMNS: Share code with BeginColumns() - move code on columns setup.
    float offset_0 = GetColumnOffset(columns.Current);
    float offset_1 = GetColumnOffset(columns.Current + 1);
    float width = offset_1 - offset_0;
    PushItemWidth(width * 0.65f);
    window.WorkRect.Max.x = window.Pos.x + offset_1 - column_padding;
}

void ImGui::EndColumns()
{
    ImGuiContext& g = *GImGui;
    ImGuiWindow* window = GetCurrentWindow();
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    IM_ASSERT(columns != NULL);

    PopItemWidth();
    if (columns.Count > 1)
    {
        PopClipRect();
        columns.Splitter.Merge(window.DrawList);
    }

    const ImGuiOldColumnFlags flags = columns.Flags;
    columns.LineMaxY = ImMax(columns.LineMaxY, window.DC.CursorPos.y);
    window.DC.CursorPos.y = columns.LineMaxY;
    if (!(flags & ImGuiOldColumnFlags_GrowParentContentsSize))
        window.DC.CursorMaxPos.x = columns.HostCursorMaxPosX;  // Restore cursor max pos, as columns don't grow parent

    // Draw columns borders and handle resize
    // The IsBeingResized flag ensure we preserve pre-resize columns width so back-and-forth are not lossy
    bool is_being_resized = false;
    if (!(flags & ImGuiOldColumnFlags_NoBorder) && !window.SkipItems)
    {
        // We clip Y boundaries CPU side because very long triangles are mishandled by some GPU drivers.
        const float y1 = ImMax(columns.HostCursorPosY, window.ClipRect.Min.y);
        const float y2 = ImMin(window.DC.CursorPos.y, window.ClipRect.Max.y);
        int dragging_column = -1;
        for (int n = 1; n < columns.Count; n++)
        {
            ImGuiOldColumnData* column = &columns.Columns[n];
            float x = window.Pos.x + GetColumnOffset(n);
            const ImGuiID column_id = columns.ID + ImGuiID(n);
            const float column_hit_hw = COLUMNS_HIT_RECT_HALF_WIDTH;
            const ImRect column_hit_rect(ImVec2(x - column_hit_hw, y1), ImVec2(x + column_hit_hw, y2));
            KeepAliveID(column_id);
            if (IsClippedEx(column_hit_rect, column_id, false))
                continue;

            bool hovered = false, held = false;
            if (!(flags & ImGuiOldColumnFlags_NoResize))
            {
                ButtonBehavior(column_hit_rect, column_id, &hovered, &held);
                if (hovered || held)
                    g.MouseCursor = ImGuiMouseCursor_ResizeEW;
                if (held && !(column.Flags & ImGuiOldColumnFlags_NoResize))
                    dragging_column = n;
            }

            // Draw column
            const ImU32 col = GetColorU32(held ? ImGuiCol_SeparatorActive : hovered ? ImGuiCol_SeparatorHovered : ImGuiCol_Separator);
            const float xi = IM_FLOOR(x);
            window.DrawList.AddLine(ImVec2(xi, y1 + 1.0f), ImVec2(xi, y2), col);
        }

        // Apply dragging after drawing the column lines, so our rendered lines are in sync with how items were displayed during the frame.
        if (dragging_column != -1)
        {
            if (!columns.IsBeingResized)
                for (int n = 0; n < columns.Count + 1; n++)
                    columns.Columns[n].OffsetNormBeforeResize = columns.Columns[n].OffsetNorm;
            columns.IsBeingResized = is_being_resized = true;
            float x = GetDraggedColumnOffset(columns, dragging_column);
            SetColumnOffset(dragging_column, x);
        }
    }
    columns.IsBeingResized = is_being_resized;

    window.WorkRect = window.ParentWorkRect;
    window.ParentWorkRect = columns.HostBackupParentWorkRect;
    window.DC.CurrentColumns = NULL;
    window.DC.ColumnsOffset.x = 0.0f;
    window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x);
}

void ImGui::Columns(int columns_count, const char* id, bool border)
{
    ImGuiWindow* window = GetCurrentWindow();
    IM_ASSERT(columns_count >= 1);

    ImGuiOldColumnFlags flags = (border ? 0 : ImGuiOldColumnFlags_NoBorder);
    //flags |= ImGuiOldColumnFlags_NoPreserveWidths; // NB: Legacy behavior
    ImGuiOldColumns* columns = window.DC.CurrentColumns;
    if (columns != NULL && columns.Count == columns_count && columns.Flags == flags)
        return;

    if (columns != NULL)
        EndColumns();

    if (columns_count != 1)
        BeginColumns(id, columns_count, flags);
}

//-------------------------------------------------------------------------

#endif // #ifndef IMGUI_DISABLE
