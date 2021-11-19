

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
