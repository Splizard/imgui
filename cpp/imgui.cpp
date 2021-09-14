

//-----------------------------------------------------------------------------
// [SECTION] MISC HELPERS/UTILITIES (Color functions)
// Note: The Convert functions are early design which are not consistent with other API.
//-----------------------------------------------------------------------------

IMGUI_API ImU32 ImAlphaBlendColors(ImU32 col_a, ImU32 col_b)
{
    float t = ((col_b >> IM_COL32_A_SHIFT) & 0xFF) / 255.f;
    int r = ImLerp((int)(col_a >> IM_COL32_R_SHIFT) & 0xFF, (int)(col_b >> IM_COL32_R_SHIFT) & 0xFF, t);
    int g = ImLerp((int)(col_a >> IM_COL32_G_SHIFT) & 0xFF, (int)(col_b >> IM_COL32_G_SHIFT) & 0xFF, t);
    int b = ImLerp((int)(col_a >> IM_COL32_B_SHIFT) & 0xFF, (int)(col_b >> IM_COL32_B_SHIFT) & 0xFF, t);
    return IM_COL32(r, g, b, 0xFF);
}

ImVec4 ImGui::ColorConvertU32ToFloat4(ImU32 in)
{
    float s = 1.0f / 255.0f;
    return ImVec4(
        ((in >> IM_COL32_R_SHIFT) & 0xFF) * s,
        ((in >> IM_COL32_G_SHIFT) & 0xFF) * s,
        ((in >> IM_COL32_B_SHIFT) & 0xFF) * s,
        ((in >> IM_COL32_A_SHIFT) & 0xFF) * s);
}

ImU32 ImGui::ColorConvertFloat4ToU32(const ImVec4& in)
{
    ImU32 out;
    out  = ((ImU32)IM_F32_TO_INT8_SAT(in.x)) << IM_COL32_R_SHIFT;
    out |= ((ImU32)IM_F32_TO_INT8_SAT(in.y)) << IM_COL32_G_SHIFT;
    out |= ((ImU32)IM_F32_TO_INT8_SAT(in.z)) << IM_COL32_B_SHIFT;
    out |= ((ImU32)IM_F32_TO_INT8_SAT(in.w)) << IM_COL32_A_SHIFT;
    return out;
}

// Convert rgb floats ([0-1],[0-1],[0-1]) to hsv floats ([0-1],[0-1],[0-1]), from Foley & van Dam p592
// Optimized http://lolengine.net/blog/2013/01/13/fast-rgb-to-hsv
void ImGui::ColorConvertRGBtoHSV(float r, float g, float b, float& out_h, float& out_s, float& out_v)
{
    float K = 0.f;
    if (g < b)
    {
        ImSwap(g, b);
        K = -1.f;
    }
    if (r < g)
    {
        ImSwap(r, g);
        K = -2.f / 6.f - K;
    }

    const float chroma = r - (g < b ? g : b);
    out_h = ImFabs(K + (g - b) / (6.f * chroma + 1e-20f));
    out_s = chroma / (r + 1e-20f);
    out_v = r;
}

// Convert hsv floats ([0-1],[0-1],[0-1]) to rgb floats ([0-1],[0-1],[0-1]), from Foley & van Dam p593
// also http://en.wikipedia.org/wiki/HSL_and_HSV
void ImGui::ColorConvertHSVtoRGB(float h, float s, float v, float& out_r, float& out_g, float& out_b)
{
    if (s == 0.0f)
    {
        // gray
        out_r = out_g = out_b = v;
        return;
    }

    h = ImFmod(h, 1.0f) / (60.0f / 360.0f);
    int   i = (int)h;
    float f = h - (float)i;
    float p = v * (1.0f - s);
    float q = v * (1.0f - s * f);
    float t = v * (1.0f - s * (1.0f - f));

    switch (i)
    {
    case 0: out_r = v; out_g = t; out_b = p; break;
    case 1: out_r = q; out_g = v; out_b = p; break;
    case 2: out_r = p; out_g = v; out_b = t; break;
    case 3: out_r = p; out_g = q; out_b = v; break;
    case 4: out_r = t; out_g = p; out_b = v; break;
    case 5: default: out_r = v; out_g = p; out_b = q; break;
    }
}


// For quicker full rebuild of a storage (instead of an incremental one), you may add all your contents and then sort once.
void ImGuiStorage::BuildSortByKey()
{
    struct StaticFunc
    {
        static int IMGUI_CDECL PairCompareByID(const void* lhs, const void* rhs)
        {
            // We can't just do a subtraction because qsort uses signed integers and subtracting our ID doesn't play well with that.
            if (((const ImGuiStoragePair*)lhs).key > ((const ImGuiStoragePair*)rhs).key) return +1;
            if (((const ImGuiStoragePair*)lhs).key < ((const ImGuiStoragePair*)rhs).key) return -1;
            return 0;
        }
    };
    if (Data.Size > 1)
        ImQsort(Data.Data, (size_t)Data.Size, sizeof(ImGuiStoragePair), StaticFunc::PairCompareByID);
}


bool ImGuiStorage::GetBool(ImGuiID key, bool default_val) const
{
    return GetInt(key, default_val ? 1 : 0) != 0;
}

float ImGuiStorage::GetFloat(ImGuiID key, float default_val) const
{
    ImGuiStoragePair* it = LowerBound(const_cast<ImVector<ImGuiStoragePair>&>(Data), key);
    if (it == Data.end() || it.key != key)
        return default_val;
    return it.val_f;
}

void* ImGuiStorage::GetVoidPtr(ImGuiID key) const
{
    ImGuiStoragePair* it = LowerBound(const_cast<ImVector<ImGuiStoragePair>&>(Data), key);
    if (it == Data.end() || it.key != key)
        return nil;
    return it.val_p;
}

// References are only valid until a new value is added to the storage. Calling a Set***() function or a Get***Ref() function invalidates the pointer.
int* ImGuiStorage::GetIntRef(ImGuiID key, int default_val)
{
    ImGuiStoragePair* it = LowerBound(Data, key);
    if (it == Data.end() || it.key != key)
        it = Data.insert(it, ImGuiStoragePair(key, default_val));
    return &it.val_i;
}

bool* ImGuiStorage::GetBoolRef(ImGuiID key, bool default_val)
{
    return (bool*)GetIntRef(key, default_val ? 1 : 0);
}

float* ImGuiStorage::GetFloatRef(ImGuiID key, float default_val)
{
    ImGuiStoragePair* it = LowerBound(Data, key);
    if (it == Data.end() || it.key != key)
        it = Data.insert(it, ImGuiStoragePair(key, default_val));
    return &it.val_f;
}

void** ImGuiStorage::GetVoidPtrRef(ImGuiID key, void* default_val)
{
    ImGuiStoragePair* it = LowerBound(Data, key);
    if (it == Data.end() || it.key != key)
        it = Data.insert(it, ImGuiStoragePair(key, default_val));
    return &it.val_p;
}


void ImGuiStorage::SetBool(ImGuiID key, bool val)
{
    SetInt(key, val ? 1 : 0);
}

void ImGuiStorage::SetFloat(ImGuiID key, float val)
{
    ImGuiStoragePair* it = LowerBound(Data, key);
    if (it == Data.end() || it.key != key)
    {
        Data.insert(it, ImGuiStoragePair(key, val));
        return;
    }
    it.val_f = val;
}

void ImGuiStorage::SetVoidPtr(ImGuiID key, void* val)
{
    ImGuiStoragePair* it = LowerBound(Data, key);
    if (it == Data.end() || it.key != key)
    {
        Data.insert(it, ImGuiStoragePair(key, val));
        return;
    }
    it.val_p = val;
}

void ImGuiStorage::SetAllInt(int v)
{
    for (int i = 0; i < Data.Size; i++)
        Data[i].val_i = v;
}

//-----------------------------------------------------------------------------
// [SECTION] ImGuiTextFilter
//-----------------------------------------------------------------------------

// Helper: Parse and apply text filters. In format "aaaaa[,bbbb][,ccccc]"
ImGuiTextFilter::ImGuiTextFilter(const char* default_filter)
{
    if (default_filter)
    {
        ImStrncpy(InputBuf, default_filter, IM_ARRAYSIZE(InputBuf));
        Build();
    }
    else
    {
        InputBuf[0] = 0;
        CountGrep = 0;
    }
}

bool ImGuiTextFilter::Draw(const char* label, float width)
{
    if (width != 0.0f)
        ImGui::SetNextItemWidth(width);
    bool value_changed = ImGui::InputText(label, InputBuf, IM_ARRAYSIZE(InputBuf));
    if (value_changed)
        Build();
    return value_changed;
}

void ImGuiTextFilter::ImGuiTextRange::split(char separator, ImVector<ImGuiTextRange>* out) const
{
    out.resize(0);
    const char* wb = b;
    const char* we = wb;
    while (we < e)
    {
        if (*we == separator)
        {
            out.push_back(ImGuiTextRange(wb, we));
            wb = we + 1;
        }
        we++;
    }
    if (wb != we)
        out.push_back(ImGuiTextRange(wb, we));
}

void ImGuiTextFilter::Build()
{
    Filters.resize(0);
    ImGuiTextRange input_range(InputBuf, InputBuf + strlen(InputBuf));
    input_range.split(',', &Filters);

    CountGrep = 0;
    for (int i = 0; i != Filters.Size; i++)
    {
        ImGuiTextRange& f = Filters[i];
        while (f.b < f.e && ImCharIsBlankA(f.b[0]))
            f.b++;
        while (f.e > f.b && ImCharIsBlankA(f.e[-1]))
            f.e--;
        if (f.empty())
            continue;
        if (Filters[i].b[0] != '-')
            CountGrep += 1;
    }
}

bool ImGuiTextFilter::PassFilter(const char* text, const char* text_end) const
{
    if (Filters.empty())
        return true;

    if (text == nil)
        text = "";

    for (int i = 0; i != Filters.Size; i++)
    {
        const ImGuiTextRange& f = Filters[i];
        if (f.empty())
            continue;
        if (f.b[0] == '-')
        {
            // Subtract
            if (ImStristr(text, text_end, f.b + 1, f.e) != nil)
                return false;
        }
        else
        {
            // Grep
            if (ImStristr(text, text_end, f.b, f.e) != nil)
                return true;
        }
    }

    // Implicit * grep
    if (CountGrep == 0)
        return true;

    return false;
}

//-----------------------------------------------------------------------------
// [SECTION] ImGuiTextBuffer
//-----------------------------------------------------------------------------

// On some platform vsnprintf() takes va_list by reference and modifies it.
// va_copy is the 'correct' way to copy a va_list but Visual Studio prior to 2013 doesn't have it.
#ifndef va_copy
#if defined(__GNUC__) || defined(__clang__)
#define va_copy(dest, src) __builtin_va_copy(dest, src)
#else
#define va_copy(dest, src) (dest = src)
#endif
#endif

char ImGuiTextBuffer::EmptyString[1] = { 0 };

void ImGuiTextBuffer::append(const char* str, const char* str_end)
{
    int len = str_end ? (int)(str_end - str) : (int)strlen(str);

    // Add zero-terminator the first time
    const int write_off = (Buf.Size != 0) ? Buf.Size : 1;
    const int needed_sz = write_off + len;
    if (write_off + len >= Buf.Capacity)
    {
        int new_capacity = Buf.Capacity * 2;
        Buf.reserve(needed_sz > new_capacity ? needed_sz : new_capacity);
    }

    Buf.resize(needed_sz);
    memcpy(&Buf[write_off - 1], str, (size_t)len);
    Buf[write_off - 1 + len] = 0;
}

void ImGuiTextBuffer::appendf(const char* fmt, ...)
{
    va_list args;
    va_start(args, fmt);
    appendfv(fmt, args);
    va_end(args);
}

// Helper: Text buffer for logging/accumulating text
void ImGuiTextBuffer::appendfv(const char* fmt, va_list args)
{
    va_list args_copy;
    va_copy(args_copy, args);

    int len = ImFormatStringV(nil, 0, fmt, args);         // FIXME-OPT: could do a first pass write attempt, likely successful on first pass.
    if (len <= 0)
    {
        va_end(args_copy);
        return;
    }

    // Add zero-terminator the first time
    const int write_off = (Buf.Size != 0) ? Buf.Size : 1;
    const int needed_sz = write_off + len;
    if (write_off + len >= Buf.Capacity)
    {
        int new_capacity = Buf.Capacity * 2;
        Buf.reserve(needed_sz > new_capacity ? needed_sz : new_capacity);
    }

    Buf.resize(needed_sz);
    ImFormatStringV(&Buf[write_off - 1], (size_t)len + 1, fmt, args_copy);
    va_end(args_copy);
}

//-----------------------------------------------------------------------------
// [SECTION] ImGuiListClipper
// This is currently not as flexible/powerful as it should be and really confusing/spaghetti, mostly because we changed
// the API mid-way through development and support two ways to using the clipper, needs some rework (see TODO)
//-----------------------------------------------------------------------------

// FIXME-TABLE: This prevents us from using ImGuiListClipper _inside_ a table cell.
// The problem we have is that without a Begin/End scheme for rows using the clipper is ambiguous.
static bool GetSkipItemForListClipping()
{
   var g = GImGui
    return (g.CurrentTable ? g.CurrentTable.HostSkipItems : g.CurrentWindow.SkipItems);
}

// Helper to calculate coarse clipping of large list of evenly sized items.
// NB: Prefer using the ImGuiListClipper higher-level helper if you can! Read comments and instructions there on how those use this sort of pattern.
// NB: 'items_count' is only used to clamp the result, if you don't know your count you can use INT_MAX
void ImGui::CalcListClipping(int items_count, float items_height, int* out_items_display_start, int* out_items_display_end)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (g.LogEnabled)
    {
        // If logging is active, do not perform any clipping
        *out_items_display_start = 0;
        *out_items_display_end = items_count;
        return;
    }
    if (GetSkipItemForListClipping())
    {
        *out_items_display_start = *out_items_display_end = 0;
        return;
    }

    // We create the union of the ClipRect and the scoring rect which at worst should be 1 page away from ClipRect
    ImRect unclipped_rect = window.ClipRect;
    if (g.NavMoveScoringItems)
        unclipped_rect.Add(g.NavScoringRect);
    if (g.NavJustMovedToId && window.NavLastIds[0] == g.NavJustMovedToId)
        unclipped_rect.Add(ImRect(window.Pos + window.NavRectRel[0].Min, window.Pos + window.NavRectRel[0].Max)); // Could store and use NavJustMovedToRectRel

    const ImVec2 pos = window.DC.CursorPos;
    int start = (int)((unclipped_rect.Min.y - pos.y) / items_height);
    int end = (int)((unclipped_rect.Max.y - pos.y) / items_height);

    // When performing a navigation request, ensure we have one item extra in the direction we are moving to
    if (g.NavMoveScoringItems && g.NavMoveClipDir == ImGuiDir_Up)
        start--;
    if (g.NavMoveScoringItems && g.NavMoveClipDir == ImGuiDir_Down)
        end++;

    start = ImClamp(start, 0, items_count);
    end = ImClamp(end + 1, start, items_count);
    *out_items_display_start = start;
    *out_items_display_end = end;
}

static void SetCursorPosYAndSetupForPrevLine(float pos_y, float line_height)
{
    // Set cursor position and a few other things so that SetScrollHereY() and Columns() can work when seeking cursor.
    // FIXME: It is problematic that we have to do that here, because custom/equivalent end-user code would stumble on the same issue.
    // The clipper should probably have a 4th step to display the last item in a regular manner.
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    float off_y = pos_y - window.DC.CursorPos.y;
    window.DC.CursorPos.y = pos_y;
    window.DC.CursorMaxPos.y = ImMax(window.DC.CursorMaxPos.y, pos_y);
    window.DC.CursorPosPrevLine.y = window.DC.CursorPos.y - line_height;  // Setting those fields so that SetScrollHereY() can properly function after the end of our clipper usage.
    window.DC.PrevLineSize.y = (line_height - g.Style.ItemSpacing.y);      // If we end up needing more accurate data (to e.g. use SameLine) we may as well make the clipper have a fourth step to let user process and display the last item in their list.
    if (ImGuiOldColumns* columns = window.DC.CurrentColumns)
        columns.LineMinY = window.DC.CursorPos.y;                         // Setting this so that cell Y position are set properly
    if (ImGuiTable* table = g.CurrentTable)
    {
        if (table.IsInsideRow)
            ImGui::TableEndRow(table);
        table.RowPosY2 = window.DC.CursorPos.y;
        const int row_increase = (int)((off_y / line_height) + 0.5f);
        //table.CurrentRow += row_increase; // Can't do without fixing TableEndRow()
        table.RowBgColorCounter += row_increase;
    }
}

ImGuiListClipper::ImGuiListClipper()
{
    memset(this, 0, sizeof(*this));
    ItemsCount = -1;
}

ImGuiListClipper::~ImGuiListClipper()
{
    IM_ASSERT(ItemsCount == -1 && "Forgot to call End(), or to Step() until false?");
}

// Use case A: Begin() called from constructor with items_height<0, then called again from Step() in StepNo 1
// Use case B: Begin() called from constructor with items_height>0
// FIXME-LEGACY: Ideally we should remove the Begin/End functions but they are part of the legacy API we still support. This is why some of the code in Step() calling Begin() and reassign some fields, spaghetti style.
void ImGuiListClipper::Begin(int items_count, float items_height)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    if (ImGuiTable* table = g.CurrentTable)
        if (table.IsInsideRow)
            ImGui::TableEndRow(table);

    StartPosY = window.DC.CursorPos.y;
    ItemsHeight = items_height;
    ItemsCount = items_count;
    ItemsFrozen = 0;
    StepNo = 0;
    DisplayStart = -1;
    DisplayEnd = 0;
}

void ImGuiListClipper::End()
{
    if (ItemsCount < 0) // Already ended
        return;

    // In theory here we should assert that ImGui::GetCursorPosY() == StartPosY + DisplayEnd * ItemsHeight, but it feels saner to just seek at the end and not assert/crash the user.
    if (ItemsCount < INT_MAX && DisplayStart >= 0)
        SetCursorPosYAndSetupForPrevLine(StartPosY + (ItemsCount - ItemsFrozen) * ItemsHeight, ItemsHeight);
    ItemsCount = -1;
    StepNo = 3;
}

bool ImGuiListClipper::Step()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    ImGuiTable* table = g.CurrentTable;
    if (table && table.IsInsideRow)
        ImGui::TableEndRow(table);

    // No items
    if (ItemsCount == 0 || GetSkipItemForListClipping())
    {
        End();
        return false;
    }

    // Step 0: Let you process the first element (regardless of it being visible or not, so we can measure the element height)
    if (StepNo == 0)
    {
        // While we are in frozen row state, keep displaying items one by one, unclipped
        // FIXME: Could be stored as a table-agnostic state.
        if (table != nil && !table.IsUnfrozenRows)
        {
            DisplayStart = ItemsFrozen;
            DisplayEnd = ItemsFrozen + 1;
            ItemsFrozen++;
            return true;
        }

        StartPosY = window.DC.CursorPos.y;
        if (ItemsHeight <= 0.0f)
        {
            // Submit the first item so we can measure its height (generally it is 0..1)
            DisplayStart = ItemsFrozen;
            DisplayEnd = ItemsFrozen + 1;
            StepNo = 1;
            return true;
        }

        // Already has item height (given by user in Begin): skip to calculating step
        DisplayStart = DisplayEnd;
        StepNo = 2;
    }

    // Step 1: the clipper infer height from first element
    if (StepNo == 1)
    {
        IM_ASSERT(ItemsHeight <= 0.0f);
        if (table)
        {
            const float pos_y1 = table.RowPosY1;   // Using this instead of StartPosY to handle clipper straddling the frozen row
            const float pos_y2 = table.RowPosY2;   // Using this instead of CursorPos.y to take account of tallest cell.
            ItemsHeight = pos_y2 - pos_y1;
            window.DC.CursorPos.y = pos_y2;
        }
        else
        {
            ItemsHeight = window.DC.CursorPos.y - StartPosY;
        }
        IM_ASSERT(ItemsHeight > 0.0f && "Unable to calculate item height! First item hasn't moved the cursor vertically!");
        StepNo = 2;
    }

    // Reached end of list
    if (DisplayEnd >= ItemsCount)
    {
        End();
        return false;
    }

    // Step 2: calculate the actual range of elements to display, and position the cursor before the first element
    if (StepNo == 2)
    {
        IM_ASSERT(ItemsHeight > 0.0f);

        int already_submitted = DisplayEnd;
        ImGui::CalcListClipping(ItemsCount - already_submitted, ItemsHeight, &DisplayStart, &DisplayEnd);
        DisplayStart += already_submitted;
        DisplayEnd += already_submitted;

        // Seek cursor
        if (DisplayStart > already_submitted)
            SetCursorPosYAndSetupForPrevLine(StartPosY + (DisplayStart - ItemsFrozen) * ItemsHeight, ItemsHeight);

        StepNo = 3;
        return true;
    }

    // Step 3: the clipper validate that we have reached the expected Y position (corresponding to element DisplayEnd),
    // Advance the cursor to the end of the list and then returns 'false' to end the loop.
    if (StepNo == 3)
    {
        // Seek cursor
        if (ItemsCount < INT_MAX)
            SetCursorPosYAndSetupForPrevLine(StartPosY + (ItemsCount - ItemsFrozen) * ItemsHeight, ItemsHeight); // advance cursor
        ItemsCount = -1;
        return false;
    }

    IM_ASSERT(0);
    return false;
}

//-----------------------------------------------------------------------------
// [SECTION] STYLING
//-----------------------------------------------------------------------------

ImGuiStyle& ImGui::GetStyle()
{
    IM_ASSERT(GImGui != nil && "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?");
    return GImGui.Style;
}

ImU32 ImGui::GetColorU32(const ImVec4& col)
{
    ImGuiStyle& style = GImGui.Style;
    ImVec4 c = col;
    c.w *= style.Alpha;
    return ColorConvertFloat4ToU32(c);
}

ImU32 ImGui::GetColorU32(ImU32 col)
{
    ImGuiStyle& style = GImGui.Style;
    if (style.Alpha >= 1.0f)
        return col;
    ImU32 a = (col & IM_COL32_A_MASK) >> IM_COL32_A_SHIFT;
    a = (ImU32)(a * style.Alpha); // We don't need to clamp 0..255 because Style.Alpha is in 0..1 range.
    return (col & ~IM_COL32_A_MASK) | (a << IM_COL32_A_SHIFT);
}

struct ImGuiStyleVarInfo
{
    ImGuiDataType   Type;
    ImU32           Count;
    ImU32           Offset;
    void*           GetVarPtr(ImGuiStyle* style) const { return (void*)((unsigned char*)style + Offset); }
};

static const ImGuiStyleVarInfo GStyleVarInfo[] =
{
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, Alpha) },               // ImGuiStyleVar_Alpha
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, DisabledAlpha) },       // ImGuiStyleVar_DisabledAlpha
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, WindowPadding) },       // ImGuiStyleVar_WindowPadding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, WindowRounding) },      // ImGuiStyleVar_WindowRounding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, WindowBorderSize) },    // ImGuiStyleVar_WindowBorderSize
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, WindowMinSize) },       // ImGuiStyleVar_WindowMinSize
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, WindowTitleAlign) },    // ImGuiStyleVar_WindowTitleAlign
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, ChildRounding) },       // ImGuiStyleVar_ChildRounding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, ChildBorderSize) },     // ImGuiStyleVar_ChildBorderSize
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, PopupRounding) },       // ImGuiStyleVar_PopupRounding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, PopupBorderSize) },     // ImGuiStyleVar_PopupBorderSize
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, FramePadding) },        // ImGuiStyleVar_FramePadding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, FrameRounding) },       // ImGuiStyleVar_FrameRounding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, FrameBorderSize) },     // ImGuiStyleVar_FrameBorderSize
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, ItemSpacing) },         // ImGuiStyleVar_ItemSpacing
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, ItemInnerSpacing) },    // ImGuiStyleVar_ItemInnerSpacing
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, IndentSpacing) },       // ImGuiStyleVar_IndentSpacing
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, CellPadding) },         // ImGuiStyleVar_CellPadding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, ScrollbarSize) },       // ImGuiStyleVar_ScrollbarSize
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, ScrollbarRounding) },   // ImGuiStyleVar_ScrollbarRounding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, GrabMinSize) },         // ImGuiStyleVar_GrabMinSize
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, GrabRounding) },        // ImGuiStyleVar_GrabRounding
    { ImGuiDataType_Float, 1, (ImU32)IM_OFFSETOF(ImGuiStyle, TabRounding) },         // ImGuiStyleVar_TabRounding
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, ButtonTextAlign) },     // ImGuiStyleVar_ButtonTextAlign
    { ImGuiDataType_Float, 2, (ImU32)IM_OFFSETOF(ImGuiStyle, SelectableTextAlign) }, // ImGuiStyleVar_SelectableTextAlign
};

static const ImGuiStyleVarInfo* GetStyleVarInfo(ImGuiStyleVar idx)
{
    IM_ASSERT(idx >= 0 && idx < ImGuiStyleVar_COUNT);
    IM_ASSERT(IM_ARRAYSIZE(GStyleVarInfo) == ImGuiStyleVar_COUNT);
    return &GStyleVarInfo[idx];
}

void ImGui::PushStyleVar(ImGuiStyleVar idx, float val)
{
    const ImGuiStyleVarInfo* var_info = GetStyleVarInfo(idx);
    if (var_info.Type == ImGuiDataType_Float && var_info.Count == 1)
    {
       var g = GImGui
        float* pvar = (float*)var_info.GetVarPtr(&g.Style);
        g.StyleVarStack.push_back(ImGuiStyleMod(idx, *pvar));
        *pvar = val;
        return;
    }
    IM_ASSERT(0 && "Called PushStyleVar() float variant but variable is not a float!");
}

void ImGui::PushStyleVar(ImGuiStyleVar idx, const ImVec2& val)
{
    const ImGuiStyleVarInfo* var_info = GetStyleVarInfo(idx);
    if (var_info.Type == ImGuiDataType_Float && var_info.Count == 2)
    {
       var g = GImGui
        ImVec2* pvar = (ImVec2*)var_info.GetVarPtr(&g.Style);
        g.StyleVarStack.push_back(ImGuiStyleMod(idx, *pvar));
        *pvar = val;
        return;
    }
    IM_ASSERT(0 && "Called PushStyleVar() ImVec2 variant but variable is not a ImVec2!");
}

void ImGui::PopStyleVar(int count)
{
   var g = GImGui
    while (count > 0)
    {
        // We avoid a generic memcpy(data, &backup.Backup.., GDataTypeSize[info.Type] * info.Count), the overhead in Debug is not worth it.
        ImGuiStyleMod& backup = g.StyleVarStack.back();
        const ImGuiStyleVarInfo* info = GetStyleVarInfo(backup.VarIdx);
        void* data = info.GetVarPtr(&g.Style);
        if (info.Type == ImGuiDataType_Float && info.Count == 1)      { ((float*)data)[0] = backup.BackupFloat[0]; }
        else if (info.Type == ImGuiDataType_Float && info.Count == 2) { ((float*)data)[0] = backup.BackupFloat[0]; ((float*)data)[1] = backup.BackupFloat[1]; }
        g.StyleVarStack.pop_back();
        count--;
    }
}

const char* ImGui::GetStyleColorName(ImGuiCol idx)
{
    // Create switch-case from enum with regexp: ImGuiCol_{.*}, -. case ImGuiCol_\1: return "\1";
    switch (idx)
    {
    case ImGuiCol_Text: return "Text";
    case ImGuiCol_TextDisabled: return "TextDisabled";
    case ImGuiCol_WindowBg: return "WindowBg";
    case ImGuiCol_ChildBg: return "ChildBg";
    case ImGuiCol_PopupBg: return "PopupBg";
    case ImGuiCol_Border: return "Border";
    case ImGuiCol_BorderShadow: return "BorderShadow";
    case ImGuiCol_FrameBg: return "FrameBg";
    case ImGuiCol_FrameBgHovered: return "FrameBgHovered";
    case ImGuiCol_FrameBgActive: return "FrameBgActive";
    case ImGuiCol_TitleBg: return "TitleBg";
    case ImGuiCol_TitleBgActive: return "TitleBgActive";
    case ImGuiCol_TitleBgCollapsed: return "TitleBgCollapsed";
    case ImGuiCol_MenuBarBg: return "MenuBarBg";
    case ImGuiCol_ScrollbarBg: return "ScrollbarBg";
    case ImGuiCol_ScrollbarGrab: return "ScrollbarGrab";
    case ImGuiCol_ScrollbarGrabHovered: return "ScrollbarGrabHovered";
    case ImGuiCol_ScrollbarGrabActive: return "ScrollbarGrabActive";
    case ImGuiCol_CheckMark: return "CheckMark";
    case ImGuiCol_SliderGrab: return "SliderGrab";
    case ImGuiCol_SliderGrabActive: return "SliderGrabActive";
    case ImGuiCol_Button: return "Button";
    case ImGuiCol_ButtonHovered: return "ButtonHovered";
    case ImGuiCol_ButtonActive: return "ButtonActive";
    case ImGuiCol_Header: return "Header";
    case ImGuiCol_HeaderHovered: return "HeaderHovered";
    case ImGuiCol_HeaderActive: return "HeaderActive";
    case ImGuiCol_Separator: return "Separator";
    case ImGuiCol_SeparatorHovered: return "SeparatorHovered";
    case ImGuiCol_SeparatorActive: return "SeparatorActive";
    case ImGuiCol_ResizeGrip: return "ResizeGrip";
    case ImGuiCol_ResizeGripHovered: return "ResizeGripHovered";
    case ImGuiCol_ResizeGripActive: return "ResizeGripActive";
    case ImGuiCol_Tab: return "Tab";
    case ImGuiCol_TabHovered: return "TabHovered";
    case ImGuiCol_TabActive: return "TabActive";
    case ImGuiCol_TabUnfocused: return "TabUnfocused";
    case ImGuiCol_TabUnfocusedActive: return "TabUnfocusedActive";
    case ImGuiCol_PlotLines: return "PlotLines";
    case ImGuiCol_PlotLinesHovered: return "PlotLinesHovered";
    case ImGuiCol_PlotHistogram: return "PlotHistogram";
    case ImGuiCol_PlotHistogramHovered: return "PlotHistogramHovered";
    case ImGuiCol_TableHeaderBg: return "TableHeaderBg";
    case ImGuiCol_TableBorderStrong: return "TableBorderStrong";
    case ImGuiCol_TableBorderLight: return "TableBorderLight";
    case ImGuiCol_TableRowBg: return "TableRowBg";
    case ImGuiCol_TableRowBgAlt: return "TableRowBgAlt";
    case ImGuiCol_TextSelectedBg: return "TextSelectedBg";
    case ImGuiCol_DragDropTarget: return "DragDropTarget";
    case ImGuiCol_NavHighlight: return "NavHighlight";
    case ImGuiCol_NavWindowingHighlight: return "NavWindowingHighlight";
    case ImGuiCol_NavWindowingDimBg: return "NavWindowingDimBg";
    case ImGuiCol_ModalWindowDimBg: return "ModalWindowDimBg";
    }
    IM_ASSERT(0);
    return "Unknown";
}


//-----------------------------------------------------------------------------
// [SECTION] RENDER HELPERS
// Some of those (internal) functions are currently quite a legacy mess - their signature and behavior will change,
// we need a nicer separation between low-level functions and high-level functions relying on the ImGui context.
// Also see imgui_draw.cpp for some more which have been reworked to not rely on ImGui:: context.
//-----------------------------------------------------------------------------

const char* ImGui::FindRenderedTextEnd(const char* text, const char* text_end)
{
    const char* text_display_end = text;
    if (!text_end)
        text_end = (const char*)-1;

    while (text_display_end < text_end && *text_display_end != '\0' && (text_display_end[0] != '#' || text_display_end[1] != '#'))
        text_display_end++;
    return text_display_end;
}

// Internal ImGui functions to render text
// RenderText***() functions calls ImDrawList::AddText() calls ImBitmapFont::RenderText()
void ImGui::RenderText(ImVec2 pos, const char* text, const char* text_end, bool hide_text_after_hash)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    // Hide anything after a '##' string
    const char* text_display_end;
    if (hide_text_after_hash)
    {
        text_display_end = FindRenderedTextEnd(text, text_end);
    }
    else
    {
        if (!text_end)
            text_end = text + strlen(text); // FIXME-OPT
        text_display_end = text_end;
    }

    if (text != text_display_end)
    {
        window.DrawList.AddText(g.Font, g.FontSize, pos, GetColorU32(ImGuiCol_Text), text, text_display_end);
        if (g.LogEnabled)
            LogRenderedText(&pos, text, text_display_end);
    }
}

void ImGui::RenderTextWrapped(ImVec2 pos, const char* text, const char* text_end, float wrap_width)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    if (!text_end)
        text_end = text + strlen(text); // FIXME-OPT

    if (text != text_end)
    {
        window.DrawList.AddText(g.Font, g.FontSize, pos, GetColorU32(ImGuiCol_Text), text, text_end, wrap_width);
        if (g.LogEnabled)
            LogRenderedText(&pos, text, text_end);
    }
}


// Another overly complex function until we reorganize everything into a nice all-in-one helper.
// This is made more complex because we have dissociated the layout rectangle (pos_min..pos_max) which define _where_ the ellipsis is, from actual clipping of text and limit of the ellipsis display.
// This is because in the context of tabs we selectively hide part of the text when the Close Button appears, but we don't want the ellipsis to move.
void ImGui::RenderTextEllipsis(ImDrawList* draw_list, const ImVec2& pos_min, const ImVec2& pos_max, float clip_max_x, float ellipsis_max_x, const char* text, const char* text_end_full, const ImVec2* text_size_if_known)
{
   var g = GImGui
    if (text_end_full == nil)
        text_end_full = FindRenderedTextEnd(text);
    const ImVec2 text_size = text_size_if_known ? *text_size_if_known : CalcTextSize(text, text_end_full, false, 0.0f);

    //draw_list.AddLine(ImVec2(pos_max.x, pos_min.y - 4), ImVec2(pos_max.x, pos_max.y + 4), IM_COL32(0, 0, 255, 255));
    //draw_list.AddLine(ImVec2(ellipsis_max_x, pos_min.y-2), ImVec2(ellipsis_max_x, pos_max.y+2), IM_COL32(0, 255, 0, 255));
    //draw_list.AddLine(ImVec2(clip_max_x, pos_min.y), ImVec2(clip_max_x, pos_max.y), IM_COL32(255, 0, 0, 255));
    // FIXME: We could technically remove (last_glyph.AdvanceX - last_glyph.X1) from text_size.x here and save a few pixels.
    if (text_size.x > pos_max.x - pos_min.x)
    {
        // Hello wo...
        // |       |   |
        // min   max   ellipsis_max
        //          <. this is generally some padding value

        const ImFont* font = draw_list._Data.Font;
        const float font_size = draw_list._Data.FontSize;
        const char* text_end_ellipsis = nil;

        ImWchar ellipsis_char = font.EllipsisChar;
        int ellipsis_char_count = 1;
        if (ellipsis_char == (ImWchar)-1)
        {
            ellipsis_char = font.DotChar;
            ellipsis_char_count = 3;
        }
        const ImFontGlyph* glyph = font.FindGlyph(ellipsis_char);

        float ellipsis_glyph_width = glyph.X1;                 // Width of the glyph with no padding on either side
        float ellipsis_total_width = ellipsis_glyph_width;      // Full width of entire ellipsis

        if (ellipsis_char_count > 1)
        {
            // Full ellipsis size without free spacing after it.
            const float spacing_between_dots = 1.0f * (draw_list._Data.FontSize / font.FontSize);
            ellipsis_glyph_width = glyph.X1 - glyph.X0 + spacing_between_dots;
            ellipsis_total_width = ellipsis_glyph_width * (float)ellipsis_char_count - spacing_between_dots;
        }

        // We can now claim the space between pos_max.x and ellipsis_max.x
        const float text_avail_width = ImMax((ImMax(pos_max.x, ellipsis_max_x) - ellipsis_total_width) - pos_min.x, 1.0f);
        float text_size_clipped_x = font.CalcTextSizeA(font_size, text_avail_width, 0.0f, text, text_end_full, &text_end_ellipsis).x;
        if (text == text_end_ellipsis && text_end_ellipsis < text_end_full)
        {
            // Always display at least 1 character if there's no room for character + ellipsis
            text_end_ellipsis = text + ImTextCountUtf8BytesFromChar(text, text_end_full);
            text_size_clipped_x = font.CalcTextSizeA(font_size, FLT_MAX, 0.0f, text, text_end_ellipsis).x;
        }
        while (text_end_ellipsis > text && ImCharIsBlankA(text_end_ellipsis[-1]))
        {
            // Trim trailing space before ellipsis (FIXME: Supporting non-ascii blanks would be nice, for this we need a function to backtrack in UTF-8 text)
            text_end_ellipsis--;
            text_size_clipped_x -= font.CalcTextSizeA(font_size, FLT_MAX, 0.0f, text_end_ellipsis, text_end_ellipsis + 1).x; // Ascii blanks are always 1 byte
        }

        // Render text, render ellipsis
        RenderTextClippedEx(draw_list, pos_min, ImVec2(clip_max_x, pos_max.y), text, text_end_ellipsis, &text_size, ImVec2(0.0f, 0.0f));
        float ellipsis_x = pos_min.x + text_size_clipped_x;
        if (ellipsis_x + ellipsis_total_width <= ellipsis_max_x)
            for (int i = 0; i < ellipsis_char_count; i++)
            {
                font.RenderChar(draw_list, font_size, ImVec2(ellipsis_x, pos_min.y), GetColorU32(ImGuiCol_Text), ellipsis_char);
                ellipsis_x += ellipsis_glyph_width;
            }
    }
    else
    {
        RenderTextClippedEx(draw_list, pos_min, ImVec2(clip_max_x, pos_max.y), text, text_end_full, &text_size, ImVec2(0.0f, 0.0f));
    }

    if (g.LogEnabled)
        LogRenderedText(&pos_min, text, text_end_full);
}

void ImGui::RenderFrameBorder(ImVec2 p_min, ImVec2 p_max, float rounding)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    const float border_size = g.Style.FrameBorderSize;
    if (border_size > 0.0f)
    {
        window.DrawList.AddRect(p_min + ImVec2(1, 1), p_max + ImVec2(1, 1), GetColorU32(ImGuiCol_BorderShadow), rounding, 0, border_size);
        window.DrawList.AddRect(p_min, p_max, GetColorU32(ImGuiCol_Border), rounding, 0, border_size);
    }
}


//-----------------------------------------------------------------------------
// [SECTION] MAIN CODE (most of the code! lots of stuff, needs tidying up!)
//-----------------------------------------------------------------------------

ImGuiID ImGuiWindow::GetID(const void* ptr)
{
    ImGuiID seed = IDStack.back();
    ImGuiID id = ImHashData(&ptr, sizeof(void*), seed);
    ImGui::KeepAliveID(id);
#ifdef IMGUI_ENABLE_TEST_ENGINE
   var g = GImGui
    IMGUI_TEST_ENGINE_ID_INFO(id, ImGuiDataType_Pointer, ptr);
#endif
    return id;
}



ImGuiID ImGuiWindow::GetIDNoKeepAlive(int n)
{
    ImGuiID seed = IDStack.back();
    ImGuiID id = ImHashData(&n, sizeof(n), seed);
#ifdef IMGUI_ENABLE_TEST_ENGINE
   var g = GImGui
    IMGUI_TEST_ENGINE_ID_INFO(id, ImGuiDataType_S32, (intptr_t)n);
#endif
    return id;
}

// This is only used in rare/specific situations to manufacture an ID out of nowhere.
ImGuiID ImGuiWindow::GetIDFromRectangle(const ImRect& r_abs)
{
    ImGuiID seed = IDStack.back();
    const int r_rel[4] = { (int)(r_abs.Min.x - Pos.x), (int)(r_abs.Min.y - Pos.y), (int)(r_abs.Max.x - Pos.x), (int)(r_abs.Max.y - Pos.y) };
    ImGuiID id = ImHashData(&r_rel, sizeof(r_rel), seed);
    ImGui::KeepAliveID(id);
    return id;
}



void ImGui::GcCompactTransientMiscBuffers()
{
   var g = GImGui
    g.ItemFlagsStack.clear();
    g.GroupStack.clear();
    TableGcCompactSettings();
}

// Free up/compact internal window buffers, we can use this when a window becomes unused.
// Not freed:
// - ImGuiWindow, ImGuiWindowSettings, Name, StateStorage, ColumnsStorage (may hold useful data)
// This should have no noticeable visual effect. When the window reappear however, expect new allocation/buffer growth/copy cost.
void ImGui::GcCompactTransientWindowBuffers(ImGuiWindow* window)
{
    window.MemoryCompacted = true;
    window.MemoryDrawListIdxCapacity = window.DrawList.IdxBuffer.Capacity;
    window.MemoryDrawListVtxCapacity = window.DrawList.VtxBuffer.Capacity;
    window.IDStack.clear();
    window.DrawList._ClearFreeMemory();
    window.DC.ChildWindows.clear();
    window.DC.ItemWidthStack.clear();
    window.DC.TextWrapPosStack.clear();
}

void ImGui::GcAwakeTransientWindowBuffers(ImGuiWindow* window)
{
    // We stored capacity of the ImDrawList buffer to reduce growth-caused allocation/copy when awakening.
    // The other buffers tends to amortize much faster.
    window.MemoryCompacted = false;
    window.DrawList.IdxBuffer.reserve(window.MemoryDrawListIdxCapacity);
    window.DrawList.VtxBuffer.reserve(window.MemoryDrawListVtxCapacity);
    window.MemoryDrawListIdxCapacity = window.MemoryDrawListVtxCapacity = 0;
}

void ImGui::MarkItemEdited(ImGuiID id)
{
    // This marking is solely to be able to provide info for IsItemDeactivatedAfterEdit().
    // ActiveId might have been released by the time we call this (as in the typical press/release button behavior) but still need need to fill the data.
   var g = GImGui
    IM_ASSERT(g.ActiveId == id || g.ActiveId == 0 || g.DragDropActive);
    IM_UNUSED(id); // Avoid unused variable warnings when asserts are compiled out.
    //IM_ASSERT(g.CurrentWindow.DC.LastItemId == id);
    g.ActiveIdHasBeenEditedThisFrame = true;
    g.ActiveIdHasBeenEditedBefore = true;
    g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_Edited;
}

// This is roughly matching the behavior of internal-facing ItemHoverable()
// - we allow hovering to be true when ActiveId==window.MoveID, so that clicking on non-interactive items such as a Text() item still returns true with IsItemHovered()
// - this should work even for non-interactive items that have no ID, so we cannot use LastItemId
bool ImGui::IsItemHovered(ImGuiHoveredFlags flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (g.NavDisableMouseHover && !g.NavDisableHighlight)
    {
        if ((g.LastItemData.InFlags & ImGuiItemFlags_Disabled) && !(flags & ImGuiHoveredFlags_AllowWhenDisabled))
            return false;
        return IsItemFocused();
    }

    // Test for bounding box overlap, as updated as ItemAdd()
    ImGuiItemStatusFlags status_flags = g.LastItemData.StatusFlags;
    if (!(status_flags & ImGuiItemStatusFlags_HoveredRect))
        return false;
    IM_ASSERT((flags & (ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows)) == 0);   // Flags not supported by this function

    // Test if we are hovering the right window (our window could be behind another window)
    // [2021/03/02] Reworked / reverted the revert, finally. Note we want e.g. BeginGroup/ItemAdd/EndGroup to work as well. (#3851)
    // [2017/10/16] Reverted commit 344d48be3 and testing RootWindow instead. I believe it is correct to NOT test for RootWindow but this leaves us unable
    // to use IsItemHovered() after EndChild() itself. Until a solution is found I believe reverting to the test from 2017/09/27 is safe since this was
    // the test that has been running for a long while.
    if (g.HoveredWindow != window && (status_flags & ImGuiItemStatusFlags_HoveredWindow) == 0)
        if ((flags & ImGuiHoveredFlags_AllowWhenOverlapped) == 0)
            return false;

    // Test if another item is active (e.g. being dragged)
    if ((flags & ImGuiHoveredFlags_AllowWhenBlockedByActiveItem) == 0)
        if (g.ActiveId != 0 && g.ActiveId != g.LastItemData.ID && !g.ActiveIdAllowOverlap && g.ActiveId != window.MoveId)
            return false;

    // Test if interactions on this window are blocked by an active popup or modal.
    // The ImGuiHoveredFlags_AllowWhenBlockedByPopup flag will be tested here.
    if (!IsWindowContentHoverable(window, flags))
        return false;

    // Test if the item is disabled
    if ((g.LastItemData.InFlags & ImGuiItemFlags_Disabled) && !(flags & ImGuiHoveredFlags_AllowWhenDisabled))
        return false;

    // Special handling for calling after Begin() which represent the title bar or tab.
    // When the window is collapsed (SkipItems==true) that last item will never be overwritten so we need to detect the case.
    if (g.LastItemData.ID == window.MoveId && window.WriteAccessed)
        return false;
    return true;
}

// Called by ItemAdd()
// Process TAB/Shift+TAB. Be mindful that this function may _clear_ the ActiveID when tabbing out.
// [WIP] This will eventually be refactored and moved into NavProcessItem()
void ImGui::ItemInputable(ImGuiWindow* window, ImGuiID id)
{
   var g = GImGui
    IM_ASSERT(id != 0 && id == g.LastItemData.ID);

    // Increment counters
    // FIXME: ImGuiItemFlags_Disabled should disable more.
    const bool is_tab_stop = (g.LastItemData.InFlags & (ImGuiItemFlags_NoTabStop | ImGuiItemFlags_Disabled)) == 0;
    window.DC.FocusCounterRegular++;
    if (is_tab_stop)
    {
        window.DC.FocusCounterTabStop++;
        if (g.NavId == id)
            g.NavIdTabCounter = window.DC.FocusCounterTabStop;
    }

    // Process TAB/Shift-TAB to tab *OUT* of the currently focused item.
    // (Note that we can always TAB out of a widget that doesn't allow tabbing in)
    if (g.ActiveId == id && g.TabFocusPressed && !IsActiveIdUsingKey(ImGuiKey_Tab) && g.TabFocusRequestNextWindow == nil)
    {
        g.TabFocusRequestNextWindow = window;
        g.TabFocusRequestNextCounterTabStop = window.DC.FocusCounterTabStop + (g.IO.KeyShift ? (is_tab_stop ? -1 : 0) : +1); // Modulo on index will be applied at the end of frame once we've got the total counter of items.
    }

    // Handle focus requests
    if (g.TabFocusRequestCurrWindow == window)
    {
        if (window.DC.FocusCounterRegular == g.TabFocusRequestCurrCounterRegular)
        {
            g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_FocusedByCode;
            return;
        }
        if (is_tab_stop && window.DC.FocusCounterTabStop == g.TabFocusRequestCurrCounterTabStop)
        {
            g.NavJustTabbedId = id;
            g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_FocusedByTabbing;
            return;
        }

        // If another item is about to be focused, we clear our own active id
        if (g.ActiveId == id)
            ClearActiveID();
    }
}

float ImGui::CalcWrapWidthForPos(const ImVec2& pos, float wrap_pos_x)
{
    if (wrap_pos_x < 0.0f)
        return 0.0f;

   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (wrap_pos_x == 0.0f)
    {
        // We could decide to setup a default wrapping max point for auto-resizing windows,
        // or have auto-wrap (with unspecified wrapping pos) behave as a ContentSize extending function?
        //if (window.Hidden && (window.Flags & ImGuiWindowFlags_AlwaysAutoResize))
        //    wrap_pos_x = ImMax(window.WorkRect.Min.x + g.FontSize * 10.0f, window.WorkRect.Max.x);
        //else
        wrap_pos_x = window.WorkRect.Max.x;
    }
    else if (wrap_pos_x > 0.0f)
    {
        wrap_pos_x += window.Pos.x - window.Scroll.x; // wrap_pos_x is provided is window local space
    }

    return ImMax(wrap_pos_x - pos.x, 1.0f);
}

// IM_ALLOC() == ImGui::MemAlloc()
void* ImGui::MemAlloc(size_t size)
{
    if (ImGuiContext* ctx = GImGui)
        ctx.IO.MetricsActiveAllocations++;
    return (*GImAllocatorAllocFunc)(size, GImAllocatorUserData);
}

// IM_FREE() == ImGui::MemFree()
void ImGui::MemFree(void* ptr)
{
    if (ptr)
        if (ImGuiContext* ctx = GImGui)
            ctx.IO.MetricsActiveAllocations--;
    return (*GImAllocatorFreeFunc)(ptr, GImAllocatorUserData);
}

const char* ImGui::GetClipboardText()
{
   var g = GImGui
    return g.IO.GetClipboardTextFn ? g.IO.GetClipboardTextFn(g.IO.ClipboardUserData) : "";
}

void ImGui::SetClipboardText(const char* text)
{
   var g = GImGui
    if (g.IO.SetClipboardTextFn)
        g.IO.SetClipboardTextFn(g.IO.ClipboardUserData, text);
}

// Internal state access - if you want to share Dear ImGui state between modules (e.g. DLL) or allocate it yourself
// Note that we still point to some static data and members (such as GFontAtlas), so the state instance you end up using will point to the static data within its module
ImGuiContext* ImGui::GetCurrentContext()
{
    return GImGui;
}


void ImGui::SetAllocatorFunctions(ImGuiMemAllocFunc alloc_func, ImGuiMemFreeFunc free_func, void* user_data)
{
    GImAllocatorAllocFunc = alloc_func;
    GImAllocatorFreeFunc = free_func;
    GImAllocatorUserData = user_data;
}

// This is provided to facilitate copying allocators from one static/DLL boundary to another (e.g. retrieve default allocator of your executable address space)
void ImGui::GetAllocatorFunctions(ImGuiMemAllocFunc* p_alloc_func, ImGuiMemFreeFunc* p_free_func, void** p_user_data)
{
    *p_alloc_func = GImAllocatorAllocFunc;
    *p_free_func = GImAllocatorFreeFunc;
    *p_user_data = GImAllocatorUserData;
}


void ImGui::DestroyContext(ImGuiContext* ctx)
{
    if (ctx == nil)
        ctx = GImGui;
    Shutdown(ctx);
    if (GImGui == ctx)
        SetCurrentContext(nil);
    IM_DELETE(ctx);
}

// No specific ordering/dependency support, will see as needed
ImGuiID ImGui::AddContextHook(ImGuiContext* ctx, const ImGuiContextHook* hook)
{
    ImGuiContext& g = *ctx;
    IM_ASSERT(hook.Callback != nil && hook.HookId == 0 && hook.Type != ImGuiContextHookType_PendingRemoval_);
    g.Hooks.push_back(*hook);
    g.Hooks.back().HookId = ++g.HookIdNext;
    return g.HookIdNext;
}

// Deferred removal, avoiding issue with changing vector while iterating it
void ImGui::RemoveContextHook(ImGuiContext* ctx, ImGuiID hook_id)
{
    ImGuiContext& g = *ctx;
    IM_ASSERT(hook_id != 0);
    for (int n = 0; n < g.Hooks.Size; n++)
        if (g.Hooks[n].HookId == hook_id)
            g.Hooks[n].Type = ImGuiContextHookType_PendingRemoval_;
}



ImGuiIO& ImGui::GetIO()
{
    IM_ASSERT(GImGui != nil && "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?");
    return GImGui.IO;
}

double ImGui::GetTime()
{
    return GImGui.Time;
}

int ImGui::GetFrameCount()
{
    return GImGui.FrameCount;
}

static ImDrawList* GetViewportDrawList(ImGuiViewportP* viewport, size_t drawlist_no, const char* drawlist_name)
{
    // Create the draw list on demand, because they are not frequently used for all viewports
   var g = GImGui
    IM_ASSERT(drawlist_no < IM_ARRAYSIZE(viewport.DrawLists));
    ImDrawList* draw_list = viewport.DrawLists[drawlist_no];
    if (draw_list == nil)
    {
        draw_list = IM_NEW(ImDrawList)(&g.DrawListSharedData);
        draw_list._OwnerName = drawlist_name;
        viewport.DrawLists[drawlist_no] = draw_list;
    }

    // Our ImDrawList system requires that there is always a command
    if (viewport.DrawListsLastFrame[drawlist_no] != g.FrameCount)
    {
        draw_list._ResetForNewFrame();
        draw_list.PushTextureID(g.IO.Fonts.TexID);
        draw_list.PushClipRect(viewport.Pos, viewport.Pos + viewport.Size, false);
        viewport.DrawListsLastFrame[drawlist_no] = g.FrameCount;
    }
    return draw_list;
}

ImDrawList* ImGui::GetBackgroundDrawList(ImGuiViewport* viewport)
{
    return GetViewportDrawList((ImGuiViewportP*)viewport, 0, "##Background");
}

ImDrawList* ImGui::GetBackgroundDrawList()
{
   var g = GImGui
    return GetBackgroundDrawList(g.Viewports[0]);
}

ImDrawList* ImGui::GetForegroundDrawList(ImGuiViewport* viewport)
{
    return GetViewportDrawList((ImGuiViewportP*)viewport, 1, "##Foreground");
}

ImDrawList* ImGui::GetForegroundDrawList()
{
   var g = GImGui
    return GetForegroundDrawList(g.Viewports[0]);
}

ImDrawListSharedData* ImGui::GetDrawListSharedData()
{
    return &GImGui.DrawListSharedData;
}

void ImGui::Initialize(ImGuiContext* context)
{
    ImGuiContext& g = *context;
    IM_ASSERT(!g.Initialized && !g.SettingsLoaded);

    // Add .ini handle for ImGuiWindow type
    {
        ImGuiSettingsHandler ini_handler;
        ini_handler.TypeName = "Window";
        ini_handler.TypeHash = ImHashStr("Window");
        ini_handler.ClearAllFn = WindowSettingsHandler_ClearAll;
        ini_handler.ReadOpenFn = WindowSettingsHandler_ReadOpen;
        ini_handler.ReadLineFn = WindowSettingsHandler_ReadLine;
        ini_handler.ApplyAllFn = WindowSettingsHandler_ApplyAll;
        ini_handler.WriteAllFn = WindowSettingsHandler_WriteAll;
        g.SettingsHandlers.push_back(ini_handler);
    }

    // Add .ini handle for ImGuiTable type
    TableSettingsInstallHandler(context);

    // Create default viewport
    ImGuiViewportP* viewport = IM_NEW(ImGuiViewportP)();
    g.Viewports.push_back(viewport);

#ifdef IMGUI_HAS_DOCK
#endif

    g.Initialized = true;
}

// This function is merely here to free heap allocations.
void ImGui::Shutdown(ImGuiContext* context)
{
    // The fonts atlas can be used prior to calling NewFrame(), so we clear it even if g.Initialized is FALSE (which would happen if we never called NewFrame)
    ImGuiContext& g = *context;
    if (g.IO.Fonts && g.FontAtlasOwnedByContext)
    {
        g.IO.Fonts.Locked = false;
        IM_DELETE(g.IO.Fonts);
    }
    g.IO.Fonts = nil;

    // Cleanup of other data are conditional on actually having initialized Dear ImGui.
    if (!g.Initialized)
        return;

    // Save settings (unless we haven't attempted to load them: CreateContext/DestroyContext without a call to NewFrame shouldn't save an empty file)
    if (g.SettingsLoaded && g.IO.IniFilename != nil)
    {
        ImGuiContext* backup_context = GImGui;
        SetCurrentContext(&g);
        SaveIniSettingsToDisk(g.IO.IniFilename);
        SetCurrentContext(backup_context);
    }

    CallContextHooks(&g, ImGuiContextHookType_Shutdown);

    // Clear everything else
    g.Windows.clear_delete();
    g.WindowsFocusOrder.clear();
    g.WindowsTempSortBuffer.clear();
    g.CurrentWindow = nil;
    g.CurrentWindowStack.clear();
    g.WindowsById.Clear();
    g.NavWindow = nil;
    g.HoveredWindow = g.HoveredWindowUnderMovingWindow = nil;
    g.ActiveIdWindow = g.ActiveIdPreviousFrameWindow = nil;
    g.MovingWindow = nil;
    g.ColorStack.clear();
    g.StyleVarStack.clear();
    g.FontStack.clear();
    g.OpenPopupStack.clear();
    g.BeginPopupStack.clear();

    g.Viewports.clear_delete();

    g.TabBars.Clear();
    g.CurrentTabBarStack.clear();
    g.ShrinkWidthBuffer.clear();

    g.Tables.Clear();
    g.TablesTempDataStack.clear_destruct();
    g.DrawChannelsTempMergeBuffer.clear();

    g.ClipboardHandlerData.clear();
    g.MenusIdSubmittedThisFrame.clear();
    g.InputTextState.ClearFreeMemory();

    g.SettingsWindows.clear();
    g.SettingsHandlers.clear();

    if (g.LogFile)
    {
#ifndef IMGUI_DISABLE_TTY_FUNCTIONS
        if (g.LogFile != stdout)
#endif
            ImFileClose(g.LogFile);
        g.LogFile = nil;
    }
    g.LogBuffer.clear();

    g.Initialized = false;
}

int ImGui::GetKeyIndex(ImGuiKey imgui_key)
{
    IM_ASSERT(imgui_key >= 0 && imgui_key < ImGuiKey_COUNT);
   var g = GImGui
    return g.IO.KeyMap[imgui_key];
}

// Note that dear imgui doesn't know the semantic of each entry of io.KeysDown[]!
// Use your own indices/enums according to how your backend/engine stored them into io.KeysDown[]!
bool ImGui::IsKeyDown(int user_key_index)
{
    if (user_key_index < 0)
        return false;
   var g = GImGui
    IM_ASSERT(user_key_index >= 0 && user_key_index < IM_ARRAYSIZE(g.IO.KeysDown));
    return g.IO.KeysDown[user_key_index];
}


int ImGui::GetKeyPressedAmount(int key_index, float repeat_delay, float repeat_rate)
{
   var g = GImGui
    if (key_index < 0)
        return 0;
    IM_ASSERT(key_index >= 0 && key_index < IM_ARRAYSIZE(g.IO.KeysDown));
    const float t = g.IO.KeysDownDuration[key_index];
    return CalcTypematicRepeatAmount(t - g.IO.DeltaTime, t, repeat_delay, repeat_rate);
}


bool ImGui::IsKeyReleased(int user_key_index)
{
   var g = GImGui
    if (user_key_index < 0) return false;
    IM_ASSERT(user_key_index >= 0 && user_key_index < IM_ARRAYSIZE(g.IO.KeysDown));
    return g.IO.KeysDownDurationPrev[user_key_index] >= 0.0f && !g.IO.KeysDown[user_key_index];
}

bool ImGui::IsMouseDown(ImGuiMouseButton button)
{
   var g = GImGui
    IM_ASSERT(button >= 0 && button < IM_ARRAYSIZE(g.IO.MouseDown));
    return g.IO.MouseDown[button];
}

bool ImGui::IsMouseReleased(ImGuiMouseButton button)
{
   var g = GImGui
    IM_ASSERT(button >= 0 && button < IM_ARRAYSIZE(g.IO.MouseDown));
    return g.IO.MouseReleased[button];
}

bool ImGui::IsMouseDoubleClicked(ImGuiMouseButton button)
{
   var g = GImGui
    IM_ASSERT(button >= 0 && button < IM_ARRAYSIZE(g.IO.MouseDown));
    return g.IO.MouseDoubleClicked[button];
}

// Return if a mouse click/drag went past the given threshold. Valid to call during the MouseReleased frame.
// [Internal] This doesn't test if the button is pressed
bool ImGui::IsMouseDragPastThreshold(ImGuiMouseButton button, float lock_threshold)
{
   var g = GImGui
    IM_ASSERT(button >= 0 && button < IM_ARRAYSIZE(g.IO.MouseDown));
    if (lock_threshold < 0.0f)
        lock_threshold = g.IO.MouseDragThreshold;
    return g.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold * lock_threshold;
}

ImVec2 ImGui::GetMousePos()
{
   var g = GImGui
    return g.IO.MousePos;
}

// NB: prefer to call right after BeginPopup(). At the time Selectable/MenuItem is activated, the popup is already closed!
ImVec2 ImGui::GetMousePosOnOpeningCurrentPopup()
{
   var g = GImGui
    if (g.BeginPopupStack.Size > 0)
        return g.OpenPopupStack[g.BeginPopupStack.Size - 1].OpenMousePos;
    return g.IO.MousePos;
}


bool ImGui::IsAnyMouseDown()
{
   var g = GImGui
    for (int n = 0; n < IM_ARRAYSIZE(g.IO.MouseDown); n++)
        if (g.IO.MouseDown[n])
            return true;
    return false;
}

// Return the delta from the initial clicking position while the mouse button is clicked or was just released.
// This is locked and return 0.0f until the mouse moves past a distance threshold at least once.
// NB: This is only valid if IsMousePosValid(). backends in theory should always keep mouse position valid when dragging even outside the client window.
ImVec2 ImGui::GetMouseDragDelta(ImGuiMouseButton button, float lock_threshold)
{
   var g = GImGui
    IM_ASSERT(button >= 0 && button < IM_ARRAYSIZE(g.IO.MouseDown));
    if (lock_threshold < 0.0f)
        lock_threshold = g.IO.MouseDragThreshold;
    if (g.IO.MouseDown[button] || g.IO.MouseReleased[button])
        if (g.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold * lock_threshold)
            if (IsMousePosValid(&g.IO.MousePos) && IsMousePosValid(&g.IO.MouseClickedPos[button]))
                return g.IO.MousePos - g.IO.MouseClickedPos[button];
    return ImVec2(0.0f, 0.0f);
}

void ImGui::ResetMouseDragDelta(ImGuiMouseButton button)
{
   var g = GImGui
    IM_ASSERT(button >= 0 && button < IM_ARRAYSIZE(g.IO.MouseDown));
    // NB: We don't need to reset g.IO.MouseDragMaxDistanceSqr
    g.IO.MouseClickedPos[button] = g.IO.MousePos;
}

ImGuiMouseCursor ImGui::GetMouseCursor()
{
    return GImGui.MouseCursor;
}

void ImGui::SetMouseCursor(ImGuiMouseCursor cursor_type)
{
    GImGui.MouseCursor = cursor_type;
}

void ImGui::CaptureKeyboardFromApp(bool capture)
{
    GImGui.WantCaptureKeyboardNextFrame = capture ? 1 : 0;
}

void ImGui::CaptureMouseFromApp(bool capture)
{
    GImGui.WantCaptureMouseNextFrame = capture ? 1 : 0;
}

bool ImGui::IsItemActivated()
{
   var g = GImGui
    if (g.ActiveId)
        if (g.ActiveId == g.LastItemData.ID && g.ActiveIdPreviousFrame != g.LastItemData.ID)
            return true;
    return false;
}

bool ImGui::IsItemDeactivated()
{
   var g = GImGui
    if (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_HasDeactivated)
        return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_Deactivated) != 0;
    return (g.ActiveIdPreviousFrame == g.LastItemData.ID && g.ActiveIdPreviousFrame != 0 && g.ActiveId != g.LastItemData.ID);
}

bool ImGui::IsItemDeactivatedAfterEdit()
{
   var g = GImGui
    return IsItemDeactivated() && (g.ActiveIdPreviousFrameHasBeenEditedBefore || (g.ActiveId == 0 && g.ActiveIdHasBeenEditedBefore));
}

// == GetItemID() == GetFocusID()
bool ImGui::IsItemFocused()
{
   var g = GImGui
    if (g.NavId != g.LastItemData.ID || g.NavId == 0)
        return false;
    return true;
}

// Important: this can be useful but it is NOT equivalent to the behavior of e.g.Button()!
// Most widgets have specific reactions based on mouse-up/down state, mouse position etc.
bool ImGui::IsItemClicked(ImGuiMouseButton mouse_button)
{
    return IsMouseClicked(mouse_button) && IsItemHovered(ImGuiHoveredFlags_None);
}

bool ImGui::IsItemToggledOpen()
{
   var g = GImGui
    return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_ToggledOpen) ? true : false;
}

bool ImGui::IsItemToggledSelection()
{
   var g = GImGui
    return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_ToggledSelection) ? true : false;
}

bool ImGui::IsAnyItemHovered()
{
   var g = GImGui
    return g.HoveredId != 0 || g.HoveredIdPreviousFrame != 0;
}

bool ImGui::IsAnyItemActive()
{
   var g = GImGui
    return g.ActiveId != 0;
}

bool ImGui::IsAnyItemFocused()
{
   var g = GImGui
    return g.NavId != 0 && !g.NavDisableHighlight;
}

bool ImGui::IsItemVisible()
{
   var g = GImGui
    return g.CurrentWindow.ClipRect.Overlaps(g.LastItemData.Rect);
}

bool ImGui::IsItemEdited()
{
   var g = GImGui
    return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_Edited) != 0;
}

// Allow last item to be overlapped by a subsequent item. Both may be activated during the same frame before the later one takes priority.
// FIXME: Although this is exposed, its interaction and ideal idiom with using ImGuiButtonFlags_AllowItemOverlap flag are extremely confusing, need rework.
void ImGui::SetItemAllowOverlap()
{
   var g = GImGui
    ImGuiID id = g.LastItemData.ID;
    if (g.HoveredId == id)
        g.HoveredIdAllowOverlap = true;
    if (g.ActiveId == id)
        g.ActiveIdAllowOverlap = true;
}

void ImGui::SetItemUsingMouseWheel()
{
   var g = GImGui
    ImGuiID id = g.LastItemData.ID;
    if (g.HoveredId == id)
        g.HoveredIdUsingMouseWheel = true;
    if (g.ActiveId == id)
        g.ActiveIdUsingMouseWheel = true;
}

ImVec2 ImGui::GetItemRectMin()
{
   var g = GImGui
    return g.LastItemData.Rect.Min;
}

ImVec2 ImGui::GetItemRectMax()
{
   var g = GImGui
    return g.LastItemData.Rect.Max;
}

ImVec2 ImGui::GetItemRectSize()
{
   var g = GImGui
    return g.LastItemData.Rect.GetSize();
}

bool ImGui::BeginChildEx(const char* name, ImGuiID id, const ImVec2& size_arg, bool border, ImGuiWindowFlags flags)
{
   var g = GImGui
    ImGuiWindow* parent_window = g.CurrentWindow;

    flags |= ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoSavedSettings | ImGuiWindowFlags_ChildWindow;
    flags |= (parent_window.Flags & ImGuiWindowFlags_NoMove);  // Inherit the NoMove flag

    // Size
    const ImVec2 content_avail = GetContentRegionAvail();
    ImVec2 size = ImFloor(size_arg);
    const int auto_fit_axises = ((size.x == 0.0f) ? (1 << ImGuiAxis_X) : 0x00) | ((size.y == 0.0f) ? (1 << ImGuiAxis_Y) : 0x00);
    if (size.x <= 0.0f)
        size.x = ImMax(content_avail.x + size.x, 4.0f); // Arbitrary minimum child size (0.0f causing too much issues)
    if (size.y <= 0.0f)
        size.y = ImMax(content_avail.y + size.y, 4.0f);
    SetNextWindowSize(size);

    // Build up name. If you need to append to a same child from multiple location in the ID stack, use BeginChild(ImGuiID id) with a stable value.
    if (name)
        ImFormatString(g.TempBuffer, IM_ARRAYSIZE(g.TempBuffer), "%s/%s_%08X", parent_window.Name, name, id);
    else
        ImFormatString(g.TempBuffer, IM_ARRAYSIZE(g.TempBuffer), "%s/%08X", parent_window.Name, id);

    const float backup_border_size = g.Style.ChildBorderSize;
    if (!border)
        g.Style.ChildBorderSize = 0.0f;
    bool ret = Begin(g.TempBuffer, nil, flags);
    g.Style.ChildBorderSize = backup_border_size;

    ImGuiWindow* child_window = g.CurrentWindow;
    child_window.ChildId = id;
    child_window.AutoFitChildAxises = (ImS8)auto_fit_axises;

    // Set the cursor to handle case where the user called SetNextWindowPos()+BeginChild() manually.
    // While this is not really documented/defined, it seems that the expected thing to do.
    if (child_window.BeginCount == 1)
        parent_window.DC.CursorPos = child_window.Pos;

    // Process navigation-in immediately so NavInit can run on first frame
    if (g.NavActivateId == id && !(flags & ImGuiWindowFlags_NavFlattened) && (child_window.DC.NavLayersActiveMask != 0 || child_window.DC.NavHasScroll))
    {
        FocusWindow(child_window);
        NavInitWindow(child_window, false);
        SetActiveID(id + 1, child_window); // Steal ActiveId with another arbitrary id so that key-press won't activate child item
        g.ActiveIdSource = ImGuiInputSource_Nav;
    }
    return ret;
}

bool ImGui::BeginChild(const char* str_id, const ImVec2& size_arg, bool border, ImGuiWindowFlags extra_flags)
{
    ImGuiWindow* window = GetCurrentWindow();
    return BeginChildEx(str_id, window.GetID(str_id), size_arg, border, extra_flags);
}

bool ImGui::BeginChild(ImGuiID id, const ImVec2& size_arg, bool border, ImGuiWindowFlags extra_flags)
{
    IM_ASSERT(id != 0);
    return BeginChildEx(nil, id, size_arg, border, extra_flags);
}

void ImGui::EndChild()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    IM_ASSERT(g.WithinEndChild == false);
    IM_ASSERT(window.Flags & ImGuiWindowFlags_ChildWindow);   // Mismatched BeginChild()/EndChild() calls

    g.WithinEndChild = true;
    if (window.BeginCount > 1)
    {
        End();
    }
    else
    {
        ImVec2 sz = window.Size;
        if (window.AutoFitChildAxises & (1 << ImGuiAxis_X)) // Arbitrary minimum zero-ish child size of 4.0f causes less trouble than a 0.0f
            sz.x = ImMax(4.0f, sz.x);
        if (window.AutoFitChildAxises & (1 << ImGuiAxis_Y))
            sz.y = ImMax(4.0f, sz.y);
        End();

        ImGuiWindow* parent_window = g.CurrentWindow;
        ImRect bb(parent_window.DC.CursorPos, parent_window.DC.CursorPos + sz);
        ItemSize(sz);
        if ((window.DC.NavLayersActiveMask != 0 || window.DC.NavHasScroll) && !(window.Flags & ImGuiWindowFlags_NavFlattened))
        {
            ItemAdd(bb, window.ChildId);
            RenderNavHighlight(bb, window.ChildId);

            // When browsing a window that has no activable items (scroll only) we keep a highlight on the child
            if (window.DC.NavLayersActiveMask == 0 && window == g.NavWindow)
                RenderNavHighlight(ImRect(bb.Min - ImVec2(2, 2), bb.Max + ImVec2(2, 2)), g.NavId, ImGuiNavHighlightFlags_TypeThin);
        }
        else
        {
            // Not navigable into
            ItemAdd(bb, 0);
        }
        if (g.HoveredWindow == window)
            g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HoveredWindow;
    }
    g.WithinEndChild = false;
    g.LogLinePosY = -FLT_MAX; // To enforce a carriage return
}

// Helper to create a child window / scrolling region that looks like a normal widget frame.
bool ImGui::BeginChildFrame(ImGuiID id, const ImVec2& size, ImGuiWindowFlags extra_flags)
{
   var g = GImGui
    const var style = g.Style
    PushStyleColor(ImGuiCol_ChildBg, style.Colors[ImGuiCol_FrameBg]);
    PushStyleVar(ImGuiStyleVar_ChildRounding, style.FrameRounding);
    PushStyleVar(ImGuiStyleVar_ChildBorderSize, style.FrameBorderSize);
    PushStyleVar(ImGuiStyleVar_WindowPadding, style.FramePadding);
    bool ret = BeginChild(id, size, true, ImGuiWindowFlags_NoMove | ImGuiWindowFlags_AlwaysUseWindowPadding | extra_flags);
    PopStyleVar(3);
    PopStyleColor();
    return ret;
}

void ImGui::EndChildFrame()
{
    EndChild();
}


ImVec2 ImGui::CalcWindowNextAutoFitSize(ImGuiWindow* window)
{
    ImVec2 size_contents_current;
    ImVec2 size_contents_ideal;
    CalcWindowContentSizes(window, &size_contents_current, &size_contents_ideal);
    ImVec2 size_auto_fit = CalcWindowAutoFitSize(window, size_contents_ideal);
    ImVec2 size_final = CalcWindowSizeAfterConstraint(window, size_auto_fit);
    return size_final;
}


// 0..3: corners (Lower-right, Lower-left, Unused, Unused)
ImGuiID ImGui::GetWindowResizeCornerID(ImGuiWindow* window, int n)
{
    IM_ASSERT(n >= 0 && n < 4);
    ImGuiID id = window.ID;
    id = ImHashStr("#RESIZE", 0, id);
    id = ImHashData(&n, sizeof(int), id);
    return id;
}

// Borders (Left, Right, Up, Down)
ImGuiID ImGui::GetWindowResizeBorderID(ImGuiWindow* window, ImGuiDir dir)
{
    IM_ASSERT(dir >= 0 && dir < 4);
    int n = (int)dir + 4;
    ImGuiID id = window.ID;
    id = ImHashStr("#RESIZE", 0, id);
    id = ImHashData(&n, sizeof(int), id);
    return id;
}

void ImGui::UpdateWindowParentAndRootLinks(ImGuiWindow* window, ImGuiWindowFlags flags, ImGuiWindow* parent_window)
{
    window.ParentWindow = parent_window;
    window.RootWindow = window.RootWindowForTitleBarHighlight = window.RootWindowForNav = window;
    if (parent_window && (flags & ImGuiWindowFlags_ChildWindow) && !(flags & ImGuiWindowFlags_Tooltip))
        window.RootWindow = parent_window.RootWindow;
    if (parent_window && !(flags & ImGuiWindowFlags_Modal) && (flags & (ImGuiWindowFlags_ChildWindow | ImGuiWindowFlags_Popup)))
        window.RootWindowForTitleBarHighlight = parent_window.RootWindowForTitleBarHighlight;
    while (window.RootWindowForNav.Flags & ImGuiWindowFlags_NavFlattened)
    {
        IM_ASSERT(window.RootWindowForNav.ParentWindow != nil);
        window.RootWindowForNav = window.RootWindowForNav.ParentWindow;
    }
}

void ImGui::BringWindowToDisplayBack(ImGuiWindow* window)
{
   var g = GImGui
    if (g.Windows[0] == window)
        return;
    for (int i = 0; i < g.Windows.Size; i++)
        if (g.Windows[i] == window)
        {
            memmove(&g.Windows[1], &g.Windows[0], (size_t)i * sizeof(ImGuiWindow*));
            g.Windows[0] = window;
            break;
        }
}

void ImGui::FocusTopMostWindowUnderOne(ImGuiWindow* under_this_window, ImGuiWindow* ignore_window)
{
   var g = GImGui

    const int start_idx = ((under_this_window != nil) ? FindWindowFocusIndex(under_this_window) : g.WindowsFocusOrder.Size) - 1;
    for (int i = start_idx; i >= 0; i--)
    {
        // We may later decide to test for different NoXXXInputs based on the active navigation input (mouse vs nav) but that may feel more confusing to the user.
        ImGuiWindow* window = g.WindowsFocusOrder[i];
        IM_ASSERT(window == window.RootWindow);
        if (window != ignore_window && window.WasActive)
            if ((window.Flags & (ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs)) != (ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs))
            {
                ImGuiWindow* focus_window = NavRestoreLastChildNavWindow(window);
                FocusWindow(focus_window);
                return;
            }
    }
    FocusWindow(nil);
}


void ImGui::PushFont(ImFont* font)
{
   var g = GImGui
    if (!font)
        font = GetDefaultFont();
    SetCurrentFont(font);
    g.FontStack.push_back(font);
    g.CurrentWindow.DrawList.PushTextureID(font.ContainerAtlas.TexID);
}

void  ImGui::PopFont()
{
   var g = GImGui
    g.CurrentWindow.DrawList.PopTextureID();
    g.FontStack.pop_back();
    SetCurrentFont(g.FontStack.empty() ? GetDefaultFont() : g.FontStack.back());
}

void ImGui::PushItemFlag(ImGuiItemFlags option, bool enabled)
{
   var g = GImGui
    ImGuiItemFlags item_flags = g.CurrentItemFlags;
    IM_ASSERT(item_flags == g.ItemFlagsStack.back());
    if (enabled)
        item_flags |= option;
    else
        item_flags &= ~option;
    g.CurrentItemFlags = item_flags;
    g.ItemFlagsStack.push_back(item_flags);
}

void ImGui::PopItemFlag()
{
   var g = GImGui
    IM_ASSERT(g.ItemFlagsStack.Size > 1); // Too many calls to PopItemFlag() - we always leave a 0 at the bottom of the stack.
    g.ItemFlagsStack.pop_back();
    g.CurrentItemFlags = g.ItemFlagsStack.back();
}

// BeginDisabled()/EndDisabled()
// - Those can be nested but it cannot be used to enable an already disabled section (a single BeginDisabled(true) in the stack is enough to keep everything disabled)
// - Visually this is currently altering alpha, but it is expected that in a future styling system this would work differently.
// - Feedback welcome at https://github.com/ocornut/imgui/issues/211
// - BeginDisabled(false) essentially does nothing useful but is provided to facilitate use of boolean expressions. If you can avoid calling BeginDisabled(False)/EndDisabled() best to avoid it.
// - Optimized shortcuts instead of PushStyleVar() + PushItemFlag()
void ImGui::BeginDisabled(bool disabled)
{
   var g = GImGui
    bool was_disabled = (g.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0;
    if (!was_disabled && disabled)
    {
        g.DisabledAlphaBackup = g.Style.Alpha;
        g.Style.Alpha *= g.Style.DisabledAlpha; // PushStyleVar(ImGuiStyleVar_Alpha, g.Style.Alpha * g.Style.DisabledAlpha);
    }
    if (was_disabled || disabled)
        g.CurrentItemFlags |= ImGuiItemFlags_Disabled;
    g.ItemFlagsStack.push_back(g.CurrentItemFlags);
}

void ImGui::EndDisabled()
{
   var g = GImGui
    bool was_disabled = (g.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0;
    //PopItemFlag();
    g.ItemFlagsStack.pop_back();
    g.CurrentItemFlags = g.ItemFlagsStack.back();
    if (was_disabled && (g.CurrentItemFlags & ImGuiItemFlags_Disabled) == 0)
        g.Style.Alpha = g.DisabledAlphaBackup; //PopStyleVar();
}

// FIXME: Look into renaming this once we have settled the new Focus/Activation/TabStop system.
void ImGui::PushAllowKeyboardFocus(bool allow_keyboard_focus)
{
    PushItemFlag(ImGuiItemFlags_NoTabStop, !allow_keyboard_focus);
}

void ImGui::PopAllowKeyboardFocus()
{
    PopItemFlag();
}

void ImGui::PushButtonRepeat(bool repeat)
{
    PushItemFlag(ImGuiItemFlags_ButtonRepeat, repeat);
}

void ImGui::PopButtonRepeat()
{
    PopItemFlag();
}

void ImGui::PushTextWrapPos(float wrap_pos_x)
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.TextWrapPosStack.push_back(window.DC.TextWrapPos);
    window.DC.TextWrapPos = wrap_pos_x;
}

void ImGui::PopTextWrapPos()
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.TextWrapPos = window.DC.TextWrapPosStack.back();
    window.DC.TextWrapPosStack.pop_back();
}

bool ImGui::IsWindowChildOf(ImGuiWindow* window, ImGuiWindow* potential_parent)
{
    if (window.RootWindow == potential_parent)
        return true;
    while (window != nil)
    {
        if (window == potential_parent)
            return true;
        window = window.ParentWindow;
    }
    return false;
}

bool ImGui::IsWindowAbove(ImGuiWindow* potential_above, ImGuiWindow* potential_below)
{
   var g = GImGui
    for (int i = g.Windows.Size - 1; i >= 0; i--)
    {
        ImGuiWindow* candidate_window = g.Windows[i];
        if (candidate_window == potential_above)
            return true;
        if (candidate_window == potential_below)
            return false;
    }
    return false;
}

bool ImGui::IsWindowHovered(ImGuiHoveredFlags flags)
{
    IM_ASSERT((flags & ImGuiHoveredFlags_AllowWhenOverlapped) == 0);   // Flags not supported by this function
   var g = GImGui
    if (g.HoveredWindow == nil)
        return false;

    if ((flags & ImGuiHoveredFlags_AnyWindow) == 0)
    {
        ImGuiWindow* window = g.CurrentWindow;
        switch (flags & (ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows))
        {
        case ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows:
            if (g.HoveredWindow.RootWindow != window.RootWindow)
                return false;
            break;
        case ImGuiHoveredFlags_RootWindow:
            if (g.HoveredWindow != window.RootWindow)
                return false;
            break;
        case ImGuiHoveredFlags_ChildWindows:
            if (!IsWindowChildOf(g.HoveredWindow, window))
                return false;
            break;
        default:
            if (g.HoveredWindow != window)
                return false;
            break;
        }
    }

    if (!IsWindowContentHoverable(g.HoveredWindow, flags))
        return false;
    if (!(flags & ImGuiHoveredFlags_AllowWhenBlockedByActiveItem))
        if (g.ActiveId != 0 && !g.ActiveIdAllowOverlap && g.ActiveId != g.HoveredWindow.MoveId)
            return false;
    return true;
}

bool ImGui::IsWindowFocused(ImGuiFocusedFlags flags)
{
   var g = GImGui

    if (flags & ImGuiFocusedFlags_AnyWindow)
        return g.NavWindow != nil;

    IM_ASSERT(g.CurrentWindow);     // Not inside a Begin()/End()
    switch (flags & (ImGuiFocusedFlags_RootWindow | ImGuiFocusedFlags_ChildWindows))
    {
    case ImGuiFocusedFlags_RootWindow | ImGuiFocusedFlags_ChildWindows:
        return g.NavWindow && g.NavWindow.RootWindow == g.CurrentWindow.RootWindow;
    case ImGuiFocusedFlags_RootWindow:
        return g.NavWindow == g.CurrentWindow.RootWindow;
    case ImGuiFocusedFlags_ChildWindows:
        return g.NavWindow && IsWindowChildOf(g.NavWindow, g.CurrentWindow);
    default:
        return g.NavWindow == g.CurrentWindow;
    }
}

// Can we focus this window with CTRL+TAB (or PadMenu + PadFocusPrev/PadFocusNext)
// Note that NoNavFocus makes the window not reachable with CTRL+TAB but it can still be focused with mouse or programmatically.
// If you want a window to never be focused, you may use the e.g. NoInputs flag.
bool ImGui::IsWindowNavFocusable(ImGuiWindow* window)
{
    return window.WasActive && window == window.RootWindow && !(window.Flags & ImGuiWindowFlags_NoNavFocus);
}

float ImGui::GetWindowWidth()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.Size.x;
}

float ImGui::GetWindowHeight()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.Size.y;
}

ImVec2 ImGui::GetWindowPos()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    return window.Pos;
}


void ImGui::SetWindowPos(const ImVec2& pos, ImGuiCond cond)
{
    ImGuiWindow* window = GetCurrentWindowRead();
    SetWindowPos(window, pos, cond);
}

void ImGui::SetWindowPos(const char* name, const ImVec2& pos, ImGuiCond cond)
{
    if (ImGuiWindow* window = FindWindowByName(name))
        SetWindowPos(window, pos, cond);
}

ImVec2 ImGui::GetWindowSize()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.Size;
}

void ImGui::SetWindowSize(const ImVec2& size, ImGuiCond cond)
{
    SetWindowSize(GImGui.CurrentWindow, size, cond);
}

void ImGui::SetWindowSize(const char* name, const ImVec2& size, ImGuiCond cond)
{
    if (ImGuiWindow* window = FindWindowByName(name))
        SetWindowSize(window, size, cond);
}

void ImGui::SetWindowCollapsed(ImGuiWindow* window, bool collapsed, ImGuiCond cond)
{
    // Test condition (NB: bit 0 is always true) and clear flags for next time
    if (cond && (window.SetWindowCollapsedAllowFlags & cond) == 0)
        return;
    window.SetWindowCollapsedAllowFlags &= ~(ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing);

    // Set
    window.Collapsed = collapsed;
}

void ImGui::SetWindowHitTestHole(ImGuiWindow* window, const ImVec2& pos, const ImVec2& size)
{
    IM_ASSERT(window.HitTestHoleSize.x == 0);     // We don't support multiple holes/hit test filters
    window.HitTestHoleSize = ImVec2ih(size);
    window.HitTestHoleOffset = ImVec2ih(pos - window.Pos);
}

void ImGui::SetWindowCollapsed(bool collapsed, ImGuiCond cond)
{
    SetWindowCollapsed(GImGui.CurrentWindow, collapsed, cond);
}

bool ImGui::IsWindowCollapsed()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.Collapsed;
}

bool ImGui::IsWindowAppearing()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.Appearing;
}

void ImGui::SetWindowCollapsed(const char* name, bool collapsed, ImGuiCond cond)
{
    if (ImGuiWindow* window = FindWindowByName(name))
        SetWindowCollapsed(window, collapsed, cond);
}

void ImGui::SetWindowFocus()
{
    FocusWindow(GImGui.CurrentWindow);
}

void ImGui::SetWindowFocus(const char* name)
{
    if (name)
    {
        if (ImGuiWindow* window = FindWindowByName(name))
            FocusWindow(window);
    }
    else
    {
        FocusWindow(nil);
    }
}

void ImGui::SetNextWindowPos(const ImVec2& pos, ImGuiCond cond, const ImVec2& pivot)
{
   var g = GImGui
    IM_ASSERT(cond == 0 || ImIsPowerOfTwo(cond)); // Make sure the user doesn't attempt to combine multiple condition flags.
    g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasPos;
    g.NextWindowData.PosVal = pos;
    g.NextWindowData.PosPivotVal = pivot;
    g.NextWindowData.PosCond = cond ? cond : ImGuiCond_Always;
}


void ImGui::SetNextWindowSizeConstraints(const ImVec2& size_min, const ImVec2& size_max, ImGuiSizeCallback custom_callback, void* custom_callback_user_data)
{
   var g = GImGui
    g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasSizeConstraint;
    g.NextWindowData.SizeConstraintRect = ImRect(size_min, size_max);
    g.NextWindowData.SizeCallback = custom_callback;
    g.NextWindowData.SizeCallbackUserData = custom_callback_user_data;
}

// Content size = inner scrollable rectangle, padded with WindowPadding.
// SetNextWindowContentSize(ImVec2(100,100) + ImGuiWindowFlags_AlwaysAutoResize will always allow submitting a 100x100 item.
void ImGui::SetNextWindowContentSize(const ImVec2& size)
{
   var g = GImGui
    g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasContentSize;
    g.NextWindowData.ContentSizeVal = ImFloor(size);
}

void ImGui::SetNextWindowScroll(const ImVec2& scroll)
{
   var g = GImGui
    g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasScroll;
    g.NextWindowData.ScrollVal = scroll;
}

void ImGui::SetNextWindowCollapsed(bool collapsed, ImGuiCond cond)
{
   var g = GImGui
    IM_ASSERT(cond == 0 || ImIsPowerOfTwo(cond)); // Make sure the user doesn't attempt to combine multiple condition flags.
    g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasCollapsed;
    g.NextWindowData.CollapsedVal = collapsed;
    g.NextWindowData.CollapsedCond = cond ? cond : ImGuiCond_Always;
}

void ImGui::SetNextWindowFocus()
{
   var g = GImGui
    g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasFocus;
}


ImDrawList* ImGui::GetWindowDrawList()
{
    ImGuiWindow* window = GetCurrentWindow();
    return window.DrawList;
}

ImFont* ImGui::GetFont()
{
    return GImGui.Font;
}

float ImGui::GetFontSize()
{
    return GImGui.FontSize;
}

ImVec2 ImGui::GetFontTexUvWhitePixel()
{
    return GImGui.DrawListSharedData.TexUvWhitePixel;
}

void ImGui::SetWindowFontScale(float scale)
{
    IM_ASSERT(scale > 0.0f);
   var g = GImGui
    ImGuiWindow* window = GetCurrentWindow();
    window.FontWindowScale = scale;
    g.FontSize = g.DrawListSharedData.FontSize = window.CalcFontSize();
}

void ImGui::ActivateItem(ImGuiID id)
{
   var g = GImGui
    g.NavNextActivateId = id;
}

void ImGui::PushFocusScope(ImGuiID id)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    g.FocusScopeStack.push_back(window.DC.NavFocusScopeIdCurrent);
    window.DC.NavFocusScopeIdCurrent = id;
}

void ImGui::PopFocusScope()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    IM_ASSERT(g.FocusScopeStack.Size > 0); // Too many PopFocusScope() ?
    window.DC.NavFocusScopeIdCurrent = g.FocusScopeStack.back();
    g.FocusScopeStack.pop_back();
}

void ImGui::SetKeyboardFocusHere(int offset)
{
    IM_ASSERT(offset >= -1);    // -1 is allowed but not below
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    g.TabFocusRequestNextWindow = window;
    g.TabFocusRequestNextCounterRegular = window.DC.FocusCounterRegular + 1 + offset;
    g.TabFocusRequestNextCounterTabStop = INT_MAX;
}

void ImGui::SetItemDefaultFocus()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (!window.Appearing)
        return;
    if (g.NavWindow == window.RootWindowForNav && (g.NavInitRequest || g.NavInitResultId != 0) && g.NavLayer == window.DC.NavLayerCurrent)
    {
        g.NavInitRequest = false;
        g.NavInitResultId = g.LastItemData.ID;
        g.NavInitResultRectRel = ImRect(g.LastItemData.Rect.Min - window.Pos, g.LastItemData.Rect.Max - window.Pos);
        NavUpdateAnyRequestFlag();
        if (!IsItemVisible())
            SetScrollHereY();
    }
}

void ImGui::SetStateStorage(ImGuiStorage* tree)
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    window.DC.StateStorage = tree ? tree : &window.StateStorage;
}

ImGuiStorage* ImGui::GetStateStorage()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.DC.StateStorage;
}

void ImGui::PushID(const char* str_id_begin, const char* str_id_end)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    ImGuiID id = window.GetIDNoKeepAlive(str_id_begin, str_id_end);
    window.IDStack.push_back(id);
}

void ImGui::PushID(const void* ptr_id)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    ImGuiID id = window.GetIDNoKeepAlive(ptr_id);
    window.IDStack.push_back(id);
}

void ImGui::PushID(int int_id)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    ImGuiID id = window.GetIDNoKeepAlive(int_id);
    window.IDStack.push_back(id);
}

// Push a given id value ignoring the ID stack as a seed.
void ImGui::PushOverrideID(ImGuiID id)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    window.IDStack.push_back(id);
}

// Helper to avoid a common series of PushOverrideID . GetID() . PopID() call
// (note that when using this pattern, TestEngine's "Stack Tool" will tend to not display the intermediate stack level.
//  for that to work we would need to do PushOverrideID() . ItemAdd() . PopID() which would alter widget code a little more)
ImGuiID ImGui::GetIDWithSeed(const char* str, const char* str_end, ImGuiID seed)
{
    ImGuiID id = ImHashStr(str, str_end ? (str_end - str) : 0, seed);
    ImGui::KeepAliveID(id);
#ifdef IMGUI_ENABLE_TEST_ENGINE
   var g = GImGui
    IMGUI_TEST_ENGINE_ID_INFO2(id, ImGuiDataType_String, str, str_end);
#endif
    return id;
}


ImGuiID ImGui::GetID(const void* ptr_id)
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.GetID(ptr_id);
}

bool ImGui::IsRectVisible(const ImVec2& size)
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.ClipRect.Overlaps(ImRect(window.DC.CursorPos, window.DC.CursorPos + size));
}

bool ImGui::IsRectVisible(const ImVec2& rect_min, const ImVec2& rect_max)
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.ClipRect.Overlaps(ImRect(rect_min, rect_max));
}


//-----------------------------------------------------------------------------
// [SECTION] ERROR CHECKING
//-----------------------------------------------------------------------------

// Helper function to verify ABI compatibility between caller code and compiled version of Dear ImGui.
// Verify that the type sizes are matching between the calling file's compilation unit and imgui.cpp's compilation unit
// If the user has inconsistent compilation settings, imgui configuration #define, packing pragma, etc. your user code
// may see different structures than what imgui.cpp sees, which is problematic.
// We usually require settings to be in imconfig.h to make sure that they are accessible to all compilation units involved with Dear ImGui.
bool ImGui::DebugCheckVersionAndDataLayout(const char* version, size_t sz_io, size_t sz_style, size_t sz_vec2, size_t sz_vec4, size_t sz_vert, size_t sz_idx)
{
    bool error = false;
    if (strcmp(version, IMGUI_VERSION) != 0) { error = true; IM_ASSERT(strcmp(version, IMGUI_VERSION) == 0 && "Mismatched version string!"); }
    if (sz_io != sizeof(ImGuiIO)) { error = true; IM_ASSERT(sz_io == sizeof(ImGuiIO) && "Mismatched struct layout!"); }
    if (sz_style != sizeof(ImGuiStyle)) { error = true; IM_ASSERT(sz_style == sizeof(ImGuiStyle) && "Mismatched struct layout!"); }
    if (sz_vec2 != sizeof(ImVec2)) { error = true; IM_ASSERT(sz_vec2 == sizeof(ImVec2) && "Mismatched struct layout!"); }
    if (sz_vec4 != sizeof(ImVec4)) { error = true; IM_ASSERT(sz_vec4 == sizeof(ImVec4) && "Mismatched struct layout!"); }
    if (sz_vert != sizeof(ImDrawVert)) { error = true; IM_ASSERT(sz_vert == sizeof(ImDrawVert) && "Mismatched struct layout!"); }
    if (sz_idx != sizeof(ImDrawIdx)) { error = true; IM_ASSERT(sz_idx == sizeof(ImDrawIdx) && "Mismatched struct layout!"); }
    return !error;
}

// Experimental recovery from incorrect usage of BeginXXX/EndXXX/PushXXX/PopXXX calls.
// Must be called during or before EndFrame().
// This is generally flawed as we are not necessarily End/Popping things in the right order.
// FIXME: Can't recover from inside BeginTabItem/EndTabItem yet.
// FIXME: Can't recover from interleaved BeginTabBar/Begin
void    ImGui::ErrorCheckEndFrameRecover(ImGuiErrorLogCallback log_callback, void* user_data)
{
    // PVS-Studio V1044 is "Loop break conditions do not depend on the number of iterations"
   var g = GImGui
    while (g.CurrentWindowStack.Size > 0)
    {
        while (g.CurrentTable && (g.CurrentTable.OuterWindow == g.CurrentWindow || g.CurrentTable.InnerWindow == g.CurrentWindow))
        {
            if (log_callback) log_callback(user_data, "Recovered from missing EndTable() in '%s'", g.CurrentTable.OuterWindow.Name);
            EndTable();
        }
        ImGuiWindow* window = g.CurrentWindow;
        IM_ASSERT(window != nil);
        while (g.CurrentTabBar != nil) //-V1044
        {
            if (log_callback) log_callback(user_data, "Recovered from missing EndTabBar() in '%s'", window.Name);
            EndTabBar();
        }
        while (window.DC.TreeDepth > 0)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing TreePop() in '%s'", window.Name);
            TreePop();
        }
        while (g.GroupStack.Size > window.DC.StackSizesOnBegin.SizeOfGroupStack)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing EndGroup() in '%s'", window.Name);
            EndGroup();
        }
        while (window.IDStack.Size > 1)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing PopID() in '%s'", window.Name);
            PopID();
        }
        while (g.ColorStack.Size > window.DC.StackSizesOnBegin.SizeOfColorStack)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing PopStyleColor() in '%s' for ImGuiCol_%s", window.Name, GetStyleColorName(g.ColorStack.back().Col));
            PopStyleColor();
        }
        while (g.StyleVarStack.Size > window.DC.StackSizesOnBegin.SizeOfStyleVarStack)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing PopStyleVar() in '%s'", window.Name);
            PopStyleVar();
        }
        while (g.FocusScopeStack.Size > window.DC.StackSizesOnBegin.SizeOfFocusScopeStack)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing PopFocusScope() in '%s'", window.Name);
            PopFocusScope();
        }
        if (g.CurrentWindowStack.Size == 1)
        {
            IM_ASSERT(g.CurrentWindow.IsFallbackWindow);
            break;
        }
        IM_ASSERT(window == g.CurrentWindow);
        if (window.Flags & ImGuiWindowFlags_ChildWindow)
        {
            if (log_callback) log_callback(user_data, "Recovered from missing EndChild() for '%s'", window.Name);
            EndChild();
        }
        else
        {
            if (log_callback) log_callback(user_data, "Recovered from missing End() for '%s'", window.Name);
            End();
        }
    }
}


//-----------------------------------------------------------------------------
// [SECTION] LAYOUT
//-----------------------------------------------------------------------------
// - ItemSize()
// - ItemAdd()
// - SameLine()
// - GetCursorScreenPos()
// - SetCursorScreenPos()
// - GetCursorPos(), GetCursorPosX(), GetCursorPosY()
// - SetCursorPos(), SetCursorPosX(), SetCursorPosY()
// - GetCursorStartPos()
// - Indent()
// - Unindent()
// - SetNextItemWidth()
// - PushItemWidth()
// - PushMultiItemsWidths()
// - PopItemWidth()
// - CalcItemWidth()
// - CalcItemSize()
// - GetTextLineHeight()
// - GetTextLineHeightWithSpacing()
// - GetFrameHeight()
// - GetFrameHeightWithSpacing()
// - GetContentRegionMax()
// - GetContentRegionMaxAbs() [Internal]
// - GetContentRegionAvail(),
// - GetWindowContentRegionMin(), GetWindowContentRegionMax()
// - BeginGroup()
// - EndGroup()
// Also see in imgui_widgets: tab bars, columns.
//-----------------------------------------------------------------------------



ImVec2 ImGui::GetCursorScreenPos()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CursorPos;
}

void ImGui::SetCursorScreenPos(const ImVec2& pos)
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.CursorPos = pos;
    window.DC.CursorMaxPos = ImMax(window.DC.CursorMaxPos, window.DC.CursorPos);
}

// User generally sees positions in window coordinates. Internally we store CursorPos in absolute screen coordinates because it is more convenient.
// Conversion happens as we pass the value to user, but it makes our naming convention confusing because GetCursorPos() == (DC.CursorPos - window.Pos). May want to rename 'DC.CursorPos'.
ImVec2 ImGui::GetCursorPos()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CursorPos - window.Pos + window.Scroll;
}

float ImGui::GetCursorPosX()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CursorPos.x - window.Pos.x + window.Scroll.x;
}

float ImGui::GetCursorPosY()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CursorPos.y - window.Pos.y + window.Scroll.y;
}

void ImGui::SetCursorPos(const ImVec2& local_pos)
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.CursorPos = window.Pos - window.Scroll + local_pos;
    window.DC.CursorMaxPos = ImMax(window.DC.CursorMaxPos, window.DC.CursorPos);
}

void ImGui::SetCursorPosX(float x)
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.CursorPos.x = window.Pos.x - window.Scroll.x + x;
    window.DC.CursorMaxPos.x = ImMax(window.DC.CursorMaxPos.x, window.DC.CursorPos.x);
}

void ImGui::SetCursorPosY(float y)
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.CursorPos.y = window.Pos.y - window.Scroll.y + y;
    window.DC.CursorMaxPos.y = ImMax(window.DC.CursorMaxPos.y, window.DC.CursorPos.y);
}

ImVec2 ImGui::GetCursorStartPos()
{
    ImGuiWindow* window = GetCurrentWindowRead();
    return window.DC.CursorStartPos - window.Pos;
}

// Affect large frame+labels widgets only.
void ImGui::SetNextItemWidth(float item_width)
{
   var g = GImGui
    g.NextItemData.Flags |= ImGuiNextItemDataFlags_HasWidth;
    g.NextItemData.Width = item_width;
}

// FIXME: Remove the == 0.0f behavior?
void ImGui::PushItemWidth(float item_width)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    window.DC.ItemWidthStack.push_back(window.DC.ItemWidth); // Backup current width
    window.DC.ItemWidth = (item_width == 0.0f ? window.ItemWidthDefault : item_width);
    g.NextItemData.Flags &= ~ImGuiNextItemDataFlags_HasWidth;
}

void ImGui::PushMultiItemsWidths(int components, float w_full)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    const var style = g.Style
    const float w_item_one  = ImMax(1.0f, IM_FLOOR((w_full - (style.ItemInnerSpacing.x) * (components - 1)) / (float)components));
    const float w_item_last = ImMax(1.0f, IM_FLOOR(w_full - (w_item_one + style.ItemInnerSpacing.x) * (components - 1)));
    window.DC.ItemWidthStack.push_back(window.DC.ItemWidth); // Backup current width
    window.DC.ItemWidthStack.push_back(w_item_last);
    for (int i = 0; i < components - 2; i++)
        window.DC.ItemWidthStack.push_back(w_item_one);
    window.DC.ItemWidth = (components == 1) ? w_item_last : w_item_one;
    g.NextItemData.Flags &= ~ImGuiNextItemDataFlags_HasWidth;
}

void ImGui::PopItemWidth()
{
    ImGuiWindow* window = GetCurrentWindow();
    window.DC.ItemWidth = window.DC.ItemWidthStack.back();
    window.DC.ItemWidthStack.pop_back();
}

// Calculate default item width given value passed to PushItemWidth() or SetNextItemWidth().
// The SetNextItemWidth() data is generally cleared/consumed by ItemAdd() or NextItemData.ClearFlags()
float ImGui::CalcItemWidth()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    float w;
    if (g.NextItemData.Flags & ImGuiNextItemDataFlags_HasWidth)
        w = g.NextItemData.Width;
    else
        w = window.DC.ItemWidth;
    if (w < 0.0f)
    {
        float region_max_x = GetContentRegionMaxAbs().x;
        w = ImMax(1.0f, region_max_x - window.DC.CursorPos.x + w);
    }
    w = IM_FLOOR(w);
    return w;
}


float ImGui::GetTextLineHeight()
{
   var g = GImGui
    return g.FontSize;
}

float ImGui::GetTextLineHeightWithSpacing()
{
   var g = GImGui
    return g.FontSize + g.Style.ItemSpacing.y;
}

float ImGui::GetFrameHeight()
{
   var g = GImGui
    return g.FontSize + g.Style.FramePadding.y * 2.0f;
}

float ImGui::GetFrameHeightWithSpacing()
{
   var g = GImGui
    return g.FontSize + g.Style.FramePadding.y * 2.0f + g.Style.ItemSpacing.y;
}

// FIXME: All the Contents Region function are messy or misleading. WE WILL AIM TO OBSOLETE ALL OF THEM WITH A NEW "WORK RECT" API. Thanks for your patience!

// FIXME: This is in window space (not screen space!).
ImVec2 ImGui::GetContentRegionMax()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    ImVec2 mx = window.ContentRegionRect.Max - window.Pos;
    if (window.DC.CurrentColumns || g.CurrentTable)
        mx.x = window.WorkRect.Max.x - window.Pos.x;
    return mx;
}

// [Internal] Absolute coordinate. Saner. This is not exposed until we finishing refactoring work rect features.
ImVec2 ImGui::GetContentRegionMaxAbs()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    ImVec2 mx = window.ContentRegionRect.Max;
    if (window.DC.CurrentColumns || g.CurrentTable)
        mx.x = window.WorkRect.Max.x;
    return mx;
}

ImVec2 ImGui::GetContentRegionAvail()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return GetContentRegionMaxAbs() - window.DC.CursorPos;
}

// In window space (not screen space!)
ImVec2 ImGui::GetWindowContentRegionMin()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.ContentRegionRect.Min - window.Pos;
}

ImVec2 ImGui::GetWindowContentRegionMax()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.ContentRegionRect.Max - window.Pos;
}

// Lock horizontal starting position + capture group bounding box into one "item" (so you can use IsItemHovered() or layout primitives such as SameLine() on whole group, etc.)
// Groups are currently a mishmash of functionalities which should perhaps be clarified and separated.
void ImGui::BeginGroup()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    g.GroupStack.resize(g.GroupStack.Size + 1);
    ImGuiGroupData& group_data = g.GroupStack.back();
    group_data.WindowID = window.ID;
    group_data.BackupCursorPos = window.DC.CursorPos;
    group_data.BackupCursorMaxPos = window.DC.CursorMaxPos;
    group_data.BackupIndent = window.DC.Indent;
    group_data.BackupGroupOffset = window.DC.GroupOffset;
    group_data.BackupCurrLineSize = window.DC.CurrLineSize;
    group_data.BackupCurrLineTextBaseOffset = window.DC.CurrLineTextBaseOffset;
    group_data.BackupActiveIdIsAlive = g.ActiveIdIsAlive;
    group_data.BackupHoveredIdIsAlive = g.HoveredId != 0;
    group_data.BackupActiveIdPreviousFrameIsAlive = g.ActiveIdPreviousFrameIsAlive;
    group_data.EmitItem = true;

    window.DC.GroupOffset.x = window.DC.CursorPos.x - window.Pos.x - window.DC.ColumnsOffset.x;
    window.DC.Indent = window.DC.GroupOffset;
    window.DC.CursorMaxPos = window.DC.CursorPos;
    window.DC.CurrLineSize = ImVec2(0.0f, 0.0f);
    if (g.LogEnabled)
        g.LogLinePosY = -FLT_MAX; // To enforce a carriage return
}

void ImGui::EndGroup()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    IM_ASSERT(g.GroupStack.Size > 0); // Mismatched BeginGroup()/EndGroup() calls

    ImGuiGroupData& group_data = g.GroupStack.back();
    IM_ASSERT(group_data.WindowID == window.ID); // EndGroup() in wrong window?

    ImRect group_bb(group_data.BackupCursorPos, ImMax(window.DC.CursorMaxPos, group_data.BackupCursorPos));

    window.DC.CursorPos = group_data.BackupCursorPos;
    window.DC.CursorMaxPos = ImMax(group_data.BackupCursorMaxPos, window.DC.CursorMaxPos);
    window.DC.Indent = group_data.BackupIndent;
    window.DC.GroupOffset = group_data.BackupGroupOffset;
    window.DC.CurrLineSize = group_data.BackupCurrLineSize;
    window.DC.CurrLineTextBaseOffset = group_data.BackupCurrLineTextBaseOffset;
    if (g.LogEnabled)
        g.LogLinePosY = -FLT_MAX; // To enforce a carriage return

    if (!group_data.EmitItem)
    {
        g.GroupStack.pop_back();
        return;
    }

    window.DC.CurrLineTextBaseOffset = ImMax(window.DC.PrevLineTextBaseOffset, group_data.BackupCurrLineTextBaseOffset);      // FIXME: Incorrect, we should grab the base offset from the *first line* of the group but it is hard to obtain now.
    ItemSize(group_bb.GetSize());
    ItemAdd(group_bb, 0);

    // If the current ActiveId was declared within the boundary of our group, we copy it to LastItemId so IsItemActive(), IsItemDeactivated() etc. will be functional on the entire group.
    // It would be be neater if we replaced window.DC.LastItemId by e.g. 'bool LastItemIsActive', but would put a little more burden on individual widgets.
    // Also if you grep for LastItemId you'll notice it is only used in that context.
    // (The two tests not the same because ActiveIdIsAlive is an ID itself, in order to be able to handle ActiveId being overwritten during the frame.)
    const bool group_contains_curr_active_id = (group_data.BackupActiveIdIsAlive != g.ActiveId) && (g.ActiveIdIsAlive == g.ActiveId) && g.ActiveId;
    const bool group_contains_prev_active_id = (group_data.BackupActiveIdPreviousFrameIsAlive == false) && (g.ActiveIdPreviousFrameIsAlive == true);
    if (group_contains_curr_active_id)
        g.LastItemData.ID = g.ActiveId;
    else if (group_contains_prev_active_id)
        g.LastItemData.ID = g.ActiveIdPreviousFrame;
    g.LastItemData.Rect = group_bb;

    // Forward Hovered flag
    const bool group_contains_curr_hovered_id = (group_data.BackupHoveredIdIsAlive == false) && g.HoveredId != 0;
    if (group_contains_curr_hovered_id)
        g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HoveredWindow;

    // Forward Edited flag
    if (group_contains_curr_active_id && g.ActiveIdHasBeenEditedThisFrame)
        g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_Edited;

    // Forward Deactivated flag
    g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HasDeactivated;
    if (group_contains_prev_active_id && g.ActiveId != g.ActiveIdPreviousFrame)
        g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_Deactivated;

    g.GroupStack.pop_back();
    //window.DrawList.AddRect(group_bb.Min, group_bb.Max, IM_COL32(255,0,255,255));   // [Debug]
}


//-----------------------------------------------------------------------------
// [SECTION] SCROLLING
//-----------------------------------------------------------------------------



// Scroll to keep newly navigated item fully into view
ImVec2 ImGui::ScrollToBringRectIntoView(ImGuiWindow* window, const ImRect& item_rect)
{
   var g = GImGui
    ImRect window_rect(window.InnerRect.Min - ImVec2(1, 1), window.InnerRect.Max + ImVec2(1, 1));
    //GetForegroundDrawList(window).AddRect(window_rect.Min, window_rect.Max, IM_COL32_WHITE); // [DEBUG]

    ImVec2 delta_scroll;
    if (!window_rect.Contains(item_rect))
    {
        if (window.ScrollbarX && item_rect.Min.x < window_rect.Min.x)
            SetScrollFromPosX(window, item_rect.Min.x - window.Pos.x - g.Style.ItemSpacing.x, 0.0f);
        else if (window.ScrollbarX && item_rect.Max.x >= window_rect.Max.x)
            SetScrollFromPosX(window, item_rect.Max.x - window.Pos.x + g.Style.ItemSpacing.x, 1.0f);
        if (item_rect.Min.y < window_rect.Min.y)
            SetScrollFromPosY(window, item_rect.Min.y - window.Pos.y - g.Style.ItemSpacing.y, 0.0f);
        else if (item_rect.Max.y >= window_rect.Max.y)
            SetScrollFromPosY(window, item_rect.Max.y - window.Pos.y + g.Style.ItemSpacing.y, 1.0f);

        ImVec2 next_scroll = CalcNextScrollFromScrollTargetAndClamp(window);
        delta_scroll = next_scroll - window.Scroll;
    }

    // Also scroll parent window to keep us into view if necessary
    if (window.Flags & ImGuiWindowFlags_ChildWindow)
        delta_scroll += ScrollToBringRectIntoView(window.ParentWindow, ImRect(item_rect.Min - delta_scroll, item_rect.Max - delta_scroll));

    return delta_scroll;
}

float ImGui::GetScrollX()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.Scroll.x;
}

float ImGui::GetScrollY()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.Scroll.y;
}

float ImGui::GetScrollMaxX()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.ScrollMax.x;
}

float ImGui::GetScrollMaxY()
{
    ImGuiWindow* window = GImGui.CurrentWindow;
    return window.ScrollMax.y;
}

void ImGui::SetScrollX(ImGuiWindow* window, float scroll_x)
{
    window.ScrollTarget.x = scroll_x;
    window.ScrollTargetCenterRatio.x = 0.0f;
    window.ScrollTargetEdgeSnapDist.x = 0.0f;
}


void ImGui::SetScrollX(float scroll_x)
{
   var g = GImGui
    SetScrollX(g.CurrentWindow, scroll_x);
}

void ImGui::SetScrollY(float scroll_y)
{
   var g = GImGui
    SetScrollY(g.CurrentWindow, scroll_y);
}

// Note that a local position will vary depending on initial scroll value,
// This is a little bit confusing so bear with us:
//  - local_pos = (absolution_pos - window.Pos)
//  - So local_x/local_y are 0.0f for a position at the upper-left corner of a window,
//    and generally local_x/local_y are >(padding+decoration) && <(size-padding-decoration) when in the visible area.
//  - They mostly exists because of legacy API.
// Following the rules above, when trying to work with scrolling code, consider that:
//  - SetScrollFromPosY(0.0f) == SetScrollY(0.0f + scroll.y) == has no effect!
//  - SetScrollFromPosY(-scroll.y) == SetScrollY(-scroll.y + scroll.y) == SetScrollY(0.0f) == reset scroll. Of course writing SetScrollY(0.0f) directly then makes more sense
// We store a target position so centering and clamping can occur on the next frame when we are guaranteed to have a known window size
void ImGui::SetScrollFromPosX(ImGuiWindow* window, float local_x, float center_x_ratio)
{
    IM_ASSERT(center_x_ratio >= 0.0f && center_x_ratio <= 1.0f);
    window.ScrollTarget.x = IM_FLOOR(local_x + window.Scroll.x); // Convert local position to scroll offset
    window.ScrollTargetCenterRatio.x = center_x_ratio;
    window.ScrollTargetEdgeSnapDist.x = 0.0f;
}

void ImGui::SetScrollFromPosY(ImGuiWindow* window, float local_y, float center_y_ratio)
{
    IM_ASSERT(center_y_ratio >= 0.0f && center_y_ratio <= 1.0f);
    const float decoration_up_height = window.TitleBarHeight() + window.MenuBarHeight(); // FIXME: Would be nice to have a more standardized access to our scrollable/client rect;
    local_y -= decoration_up_height;
    window.ScrollTarget.y = IM_FLOOR(local_y + window.Scroll.y); // Convert local position to scroll offset
    window.ScrollTargetCenterRatio.y = center_y_ratio;
    window.ScrollTargetEdgeSnapDist.y = 0.0f;
}

void ImGui::SetScrollFromPosX(float local_x, float center_x_ratio)
{
   var g = GImGui
    SetScrollFromPosX(g.CurrentWindow, local_x, center_x_ratio);
}

void ImGui::SetScrollFromPosY(float local_y, float center_y_ratio)
{
   var g = GImGui
    SetScrollFromPosY(g.CurrentWindow, local_y, center_y_ratio);
}

// center_x_ratio: 0.0f left of last item, 0.5f horizontal center of last item, 1.0f right of last item.
void ImGui::SetScrollHereX(float center_x_ratio)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    float spacing_x = ImMax(window.WindowPadding.x, g.Style.ItemSpacing.x);
    float target_pos_x = ImLerp(g.LastItemData.Rect.Min.x - spacing_x, g.LastItemData.Rect.Max.x + spacing_x, center_x_ratio);
    SetScrollFromPosX(window, target_pos_x - window.Pos.x, center_x_ratio); // Convert from absolute to local pos

    // Tweak: snap on edges when aiming at an item very close to the edge
    window.ScrollTargetEdgeSnapDist.x = ImMax(0.0f, window.WindowPadding.x - spacing_x);
}

// center_y_ratio: 0.0f top of last item, 0.5f vertical center of last item, 1.0f bottom of last item.
void ImGui::SetScrollHereY(float center_y_ratio)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    float spacing_y = ImMax(window.WindowPadding.y, g.Style.ItemSpacing.y);
    float target_pos_y = ImLerp(window.DC.CursorPosPrevLine.y - spacing_y, window.DC.CursorPosPrevLine.y + window.DC.PrevLineSize.y + spacing_y, center_y_ratio);
    SetScrollFromPosY(window, target_pos_y - window.Pos.y, center_y_ratio); // Convert from absolute to local pos

    // Tweak: snap on edges when aiming at an item very close to the edge
    window.ScrollTargetEdgeSnapDist.y = ImMax(0.0f, window.WindowPadding.y - spacing_y);
}

//-----------------------------------------------------------------------------
// [SECTION] POPUPS
//-----------------------------------------------------------------------------

bool ImGui::IsPopupOpen(const char* str_id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    ImGuiID id = (popup_flags & ImGuiPopupFlags_AnyPopupId) ? 0 : g.CurrentWindow.GetID(str_id);
    if ((popup_flags & ImGuiPopupFlags_AnyPopupLevel) && id != 0)
        IM_ASSERT(0 && "Cannot use IsPopupOpen() with a string id and ImGuiPopupFlags_AnyPopupLevel."); // But non-string version is legal and used internally
    return IsPopupOpen(id, popup_flags);
}



void ImGui::OpenPopup(const char* str_id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    OpenPopupEx(g.CurrentWindow.GetID(str_id), popup_flags);
}

void ImGui::OpenPopup(ImGuiID id, ImGuiPopupFlags popup_flags)
{
    OpenPopupEx(id, popup_flags);
}

// Mark popup as open (toggle toward open state).
// Popups are closed when user click outside, or activate a pressable item, or CloseCurrentPopup() is called within a BeginPopup()/EndPopup() block.
// Popup identifiers are relative to the current ID-stack (so OpenPopup and BeginPopup needs to be at the same level).
// One open popup per level of the popup hierarchy (NB: when assigning we reset the Window member of ImGuiPopupRef to nil)
void ImGui::OpenPopupEx(ImGuiID id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    ImGuiWindow* parent_window = g.CurrentWindow;
    const int current_stack_size = g.BeginPopupStack.Size;

    if (popup_flags & ImGuiPopupFlags_NoOpenOverExistingPopup)
        if (IsPopupOpen(0u, ImGuiPopupFlags_AnyPopupId))
            return;

    ImGuiPopupData popup_ref; // Tagged as new ref as Window will be set back to nil if we write this into OpenPopupStack.
    popup_ref.PopupId = id;
    popup_ref.Window = nil;
    popup_ref.SourceWindow = g.NavWindow;
    popup_ref.OpenFrameCount = g.FrameCount;
    popup_ref.OpenParentId = parent_window.IDStack.back();
    popup_ref.OpenPopupPos = NavCalcPreferredRefPos();
    popup_ref.OpenMousePos = IsMousePosValid(&g.IO.MousePos) ? g.IO.MousePos : popup_ref.OpenPopupPos;

    IMGUI_DEBUG_LOG_POPUP("OpenPopupEx(0x%08X)\n", id);
    if (g.OpenPopupStack.Size < current_stack_size + 1)
    {
        g.OpenPopupStack.push_back(popup_ref);
    }
    else
    {
        // Gently handle the user mistakenly calling OpenPopup() every frame. It is a programming mistake! However, if we were to run the regular code path, the ui
        // would become completely unusable because the popup will always be in hidden-while-calculating-size state _while_ claiming focus. Which would be a very confusing
        // situation for the programmer. Instead, we silently allow the popup to proceed, it will keep reappearing and the programming error will be more obvious to understand.
        if (g.OpenPopupStack[current_stack_size].PopupId == id && g.OpenPopupStack[current_stack_size].OpenFrameCount == g.FrameCount - 1)
        {
            g.OpenPopupStack[current_stack_size].OpenFrameCount = popup_ref.OpenFrameCount;
        }
        else
        {
            // Close child popups if any, then flag popup for open/reopen
            ClosePopupToLevel(current_stack_size, false);
            g.OpenPopupStack.push_back(popup_ref);
        }

        // When reopening a popup we first refocus its parent, otherwise if its parent is itself a popup it would get closed by ClosePopupsOverWindow().
        // This is equivalent to what ClosePopupToLevel() does.
        //if (g.OpenPopupStack[current_stack_size].PopupId == id)
        //    FocusWindow(parent_window);
    }
}

void ImGui::ClosePopupToLevel(int remaining, bool restore_focus_to_window_under_popup)
{
   var g = GImGui
    IMGUI_DEBUG_LOG_POPUP("ClosePopupToLevel(%d), restore_focus_to_window_under_popup=%d\n", remaining, restore_focus_to_window_under_popup);
    IM_ASSERT(remaining >= 0 && remaining < g.OpenPopupStack.Size);

    // Trim open popup stack
    ImGuiWindow* focus_window = g.OpenPopupStack[remaining].SourceWindow;
    ImGuiWindow* popup_window = g.OpenPopupStack[remaining].Window;
    g.OpenPopupStack.resize(remaining);

    if (restore_focus_to_window_under_popup)
    {
        if (focus_window && !focus_window.WasActive && popup_window)
        {
            // Fallback
            FocusTopMostWindowUnderOne(popup_window, nil);
        }
        else
        {
            if (g.NavLayer == ImGuiNavLayer_Main && focus_window)
                focus_window = NavRestoreLastChildNavWindow(focus_window);
            FocusWindow(focus_window);
        }
    }
}

// Close the popup we have begin-ed into.
void ImGui::CloseCurrentPopup()
{
   var g = GImGui
    int popup_idx = g.BeginPopupStack.Size - 1;
    if (popup_idx < 0 || popup_idx >= g.OpenPopupStack.Size || g.BeginPopupStack[popup_idx].PopupId != g.OpenPopupStack[popup_idx].PopupId)
        return;

    // Closing a menu closes its top-most parent popup (unless a modal)
    while (popup_idx > 0)
    {
        ImGuiWindow* popup_window = g.OpenPopupStack[popup_idx].Window;
        ImGuiWindow* parent_popup_window = g.OpenPopupStack[popup_idx - 1].Window;
        bool close_parent = false;
        if (popup_window && (popup_window.Flags & ImGuiWindowFlags_ChildMenu))
            if (parent_popup_window == nil || !(parent_popup_window.Flags & ImGuiWindowFlags_Modal))
                close_parent = true;
        if (!close_parent)
            break;
        popup_idx--;
    }
    IMGUI_DEBUG_LOG_POPUP("CloseCurrentPopup %d . %d\n", g.BeginPopupStack.Size - 1, popup_idx);
    ClosePopupToLevel(popup_idx, true);

    // A common pattern is to close a popup when selecting a menu item/selectable that will open another window.
    // To improve this usage pattern, we avoid nav highlight for a single frame in the parent window.
    // Similarly, we could avoid mouse hover highlight in this window but it is less visually problematic.
    if (ImGuiWindow* window = g.NavWindow)
        window.DC.NavHideHighlightOneFrame = true;
}

// Attention! BeginPopup() adds default flags which BeginPopupEx()!
bool ImGui::BeginPopupEx(ImGuiID id, ImGuiWindowFlags flags)
{
   var g = GImGui
    if (!IsPopupOpen(id, ImGuiPopupFlags_None))
    {
        g.NextWindowData.ClearFlags(); // We behave like Begin() and need to consume those values
        return false;
    }

    char name[20];
    if (flags & ImGuiWindowFlags_ChildMenu)
        ImFormatString(name, IM_ARRAYSIZE(name), "##Menu_%02d", g.BeginPopupStack.Size); // Recycle windows based on depth
    else
        ImFormatString(name, IM_ARRAYSIZE(name), "##Popup_%08x", id); // Not recycling, so we can close/open during the same frame

    flags |= ImGuiWindowFlags_Popup;
    bool is_open = Begin(name, nil, flags);
    if (!is_open) // NB: Begin can return false when the popup is completely clipped (e.g. zero size display)
        EndPopup();

    return is_open;
}

bool ImGui::BeginPopup(const char* str_id, ImGuiWindowFlags flags)
{
   var g = GImGui
    if (g.OpenPopupStack.Size <= g.BeginPopupStack.Size) // Early out for performance
    {
        g.NextWindowData.ClearFlags(); // We behave like Begin() and need to consume those values
        return false;
    }
    flags |= ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoSavedSettings;
    return BeginPopupEx(g.CurrentWindow.GetID(str_id), flags);
}

// If 'p_open' is specified for a modal popup window, the popup will have a regular close button which will close the popup.
// Note that popup visibility status is owned by Dear ImGui (and manipulated with e.g. OpenPopup) so the actual value of *p_open is meaningless here.
bool ImGui::BeginPopupModal(const char* name, bool* p_open, ImGuiWindowFlags flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    const ImGuiID id = window.GetID(name);
    if (!IsPopupOpen(id, ImGuiPopupFlags_None))
    {
        g.NextWindowData.ClearFlags(); // We behave like Begin() and need to consume those values
        return false;
    }

    // Center modal windows by default for increased visibility
    // (this won't really last as settings will kick in, and is mostly for backward compatibility. user may do the same themselves)
    // FIXME: Should test for (PosCond & window.SetWindowPosAllowFlags) with the upcoming window.
    if ((g.NextWindowData.Flags & ImGuiNextWindowDataFlags_HasPos) == 0)
    {
        const ImGuiViewport* viewport = GetMainViewport();
        SetNextWindowPos(viewport.GetCenter(), ImGuiCond_FirstUseEver, ImVec2(0.5f, 0.5f));
    }

    flags |= ImGuiWindowFlags_Popup | ImGuiWindowFlags_Modal | ImGuiWindowFlags_NoCollapse;
    const bool is_open = Begin(name, p_open, flags);
    if (!is_open || (p_open && !*p_open)) // NB: is_open can be 'false' when the popup is completely clipped (e.g. zero size display)
    {
        EndPopup();
        if (is_open)
            ClosePopupToLevel(g.BeginPopupStack.Size, true);
        return false;
    }
    return is_open;
}

void ImGui::EndPopup()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    IM_ASSERT(window.Flags & ImGuiWindowFlags_Popup);  // Mismatched BeginPopup()/EndPopup() calls
    IM_ASSERT(g.BeginPopupStack.Size > 0);

    // Make all menus and popups wrap around for now, may need to expose that policy (e.g. focus scope could include wrap/loop policy flags used by new move requests)
    if (g.NavWindow == window)
        NavMoveRequestTryWrapping(window, ImGuiNavMoveFlags_LoopY);

    // Child-popups don't need to be laid out
    IM_ASSERT(g.WithinEndChild == false);
    if (window.Flags & ImGuiWindowFlags_ChildWindow)
        g.WithinEndChild = true;
    End();
    g.WithinEndChild = false;
}

// Helper to open a popup if mouse button is released over the item
// - This is essentially the same as BeginPopupContextItem() but without the trailing BeginPopup()
void ImGui::OpenPopupOnItemClick(const char* str_id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    int mouse_button = (popup_flags & ImGuiPopupFlags_MouseButtonMask_);
    if (IsMouseReleased(mouse_button) && IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup))
    {
        ImGuiID id = str_id ? window.GetID(str_id) : g.LastItemData.ID;    // If user hasn't passed an ID, we can use the LastItemID. Using LastItemID as a Popup ID won't conflict!
        IM_ASSERT(id != 0);                                             // You cannot pass a nil str_id if the last item has no identifier (e.g. a Text() item)
        OpenPopupEx(id, popup_flags);
    }
}

// This is a helper to handle the simplest case of associating one named popup to one given widget.
// - To create a popup associated to the last item, you generally want to pass a nil value to str_id.
// - To create a popup with a specific identifier, pass it in str_id.
//    - This is useful when using using BeginPopupContextItem() on an item which doesn't have an identifier, e.g. a Text() call.
//    - This is useful when multiple code locations may want to manipulate/open the same popup, given an explicit id.
// - You may want to handle the whole on user side if you have specific needs (e.g. tweaking IsItemHovered() parameters).
//   This is essentially the same as:
//       id = str_id ? GetID(str_id) : GetItemID();
//       OpenPopupOnItemClick(str_id);
//       return BeginPopup(id);
//   Which is essentially the same as:
//       id = str_id ? GetID(str_id) : GetItemID();
//       if (IsItemHovered() && IsMouseReleased(ImGuiMouseButton_Right))
//           OpenPopup(id);
//       return BeginPopup(id);
//   The main difference being that this is tweaked to avoid computing the ID twice.
bool ImGui::BeginPopupContextItem(const char* str_id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (window.SkipItems)
        return false;
    ImGuiID id = str_id ? window.GetID(str_id) : g.LastItemData.ID;    // If user hasn't passed an ID, we can use the LastItemID. Using LastItemID as a Popup ID won't conflict!
    IM_ASSERT(id != 0);                                             // You cannot pass a nil str_id if the last item has no identifier (e.g. a Text() item)
    int mouse_button = (popup_flags & ImGuiPopupFlags_MouseButtonMask_);
    if (IsMouseReleased(mouse_button) && IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup))
        OpenPopupEx(id, popup_flags);
    return BeginPopupEx(id, ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoSavedSettings);
}

bool ImGui::BeginPopupContextWindow(const char* str_id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (!str_id)
        str_id = "window_context";
    ImGuiID id = window.GetID(str_id);
    int mouse_button = (popup_flags & ImGuiPopupFlags_MouseButtonMask_);
    if (IsMouseReleased(mouse_button) && IsWindowHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup))
        if (!(popup_flags & ImGuiPopupFlags_NoOpenOverItems) || !IsAnyItemHovered())
            OpenPopupEx(id, popup_flags);
    return BeginPopupEx(id, ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoSavedSettings);
}

bool ImGui::BeginPopupContextVoid(const char* str_id, ImGuiPopupFlags popup_flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    if (!str_id)
        str_id = "void_context";
    ImGuiID id = window.GetID(str_id);
    int mouse_button = (popup_flags & ImGuiPopupFlags_MouseButtonMask_);
    if (IsMouseReleased(mouse_button) && !IsWindowHovered(ImGuiHoveredFlags_AnyWindow))
        if (GetTopMostPopupModal() == nil)
            OpenPopupEx(id, popup_flags);
    return BeginPopupEx(id, ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoSavedSettings);
}

ImGuiDir ImGetDirQuadrantFromDelta(float dx, float dy)
{
    if (ImFabs(dx) > ImFabs(dy))
        return (dx > 0.0f) ? ImGuiDir_Right : ImGuiDir_Left;
    return (dy > 0.0f) ? ImGuiDir_Down : ImGuiDir_Up;
}


// Forward will reuse the move request again on the next frame (generally with modifications done to it)
void ImGui::NavMoveRequestForward(ImGuiDir move_dir, ImGuiDir clip_dir, ImGuiNavMoveFlags move_flags)
{
   var g = GImGui
    IM_ASSERT(g.NavMoveForwardToNextFrame == false);
    NavMoveRequestCancel();
    g.NavMoveForwardToNextFrame = true;
    g.NavMoveDir = move_dir;
    g.NavMoveClipDir = clip_dir;
    g.NavMoveFlags = move_flags | ImGuiNavMoveFlags_Forwarded;
}

// Navigation wrap-around logic is delayed to the end of the frame because this operation is only valid after entire
// popup is assembled and in case of appended popups it is not clear which EndPopup() call is final.
void ImGui::NavMoveRequestTryWrapping(ImGuiWindow* window, ImGuiNavMoveFlags wrap_flags)
{
   var g = GImGui
    IM_ASSERT(wrap_flags != 0); // Call with _WrapX, _WrapY, _LoopX, _LoopY
    // In theory we should test for NavMoveRequestButNoResultYet() but there's no point doing it, NavEndFrame() will do the same test
    if (g.NavWindow == window && g.NavMoveScoringItems && g.NavLayer == ImGuiNavLayer_Main)
        g.NavMoveFlags |= wrap_flags;
}

//-----------------------------------------------------------------------------
// [SECTION] DRAG AND DROP
//-----------------------------------------------------------------------------

void ImGui::ClearDragDrop()
{
   var g = GImGui
    g.DragDropActive = false;
    g.DragDropPayload.Clear();
    g.DragDropAcceptFlags = ImGuiDragDropFlags_None;
    g.DragDropAcceptIdCurr = g.DragDropAcceptIdPrev = 0;
    g.DragDropAcceptIdCurrRectSurface = FLT_MAX;
    g.DragDropAcceptFrameCount = -1;

    g.DragDropPayloadBufHeap.clear();
    memset(&g.DragDropPayloadBufLocal, 0, sizeof(g.DragDropPayloadBufLocal));
}

// When this returns true you need to: a) call SetDragDropPayload() exactly once, b) you may render the payload visual/description, c) call EndDragDropSource()
// If the item has an identifier:
// - This assume/require the item to be activated (typically via ButtonBehavior).
// - Therefore if you want to use this with a mouse button other than left mouse button, it is up to the item itself to activate with another button.
// - We then pull and use the mouse button that was used to activate the item and use it to carry on the drag.
// If the item has no identifier:
// - Currently always assume left mouse button.
bool ImGui::BeginDragDropSource(ImGuiDragDropFlags flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    // FIXME-DRAGDROP: While in the common-most "drag from non-zero active id" case we can tell the mouse button,
    // in both SourceExtern and id==0 cases we may requires something else (explicit flags or some heuristic).
    ImGuiMouseButton mouse_button = ImGuiMouseButton_Left;

    bool source_drag_active = false;
    ImGuiID source_id = 0;
    ImGuiID source_parent_id = 0;
    if (!(flags & ImGuiDragDropFlags_SourceExtern))
    {
        source_id = g.LastItemData.ID;
        if (source_id != 0)
        {
            // Common path: items with ID
            if (g.ActiveId != source_id)
                return false;
            if (g.ActiveIdMouseButton != -1)
                mouse_button = g.ActiveIdMouseButton;
            if (g.IO.MouseDown[mouse_button] == false)
                return false;
            g.ActiveIdAllowOverlap = false;
        }
        else
        {
            // Uncommon path: items without ID
            if (g.IO.MouseDown[mouse_button] == false)
                return false;

            // If you want to use BeginDragDropSource() on an item with no unique identifier for interaction, such as Text() or Image(), you need to:
            // A) Read the explanation below, B) Use the ImGuiDragDropFlags_SourceAllownilID flag, C) Swallow your programmer pride.
            if (!(flags & ImGuiDragDropFlags_SourceAllownilID))
            {
                IM_ASSERT(0);
                return false;
            }

            // Early out
            if ((g.LastItemData.StatusFlags & ImGuiItemStatusFlags_HoveredRect) == 0 && (g.ActiveId == 0 || g.ActiveIdWindow != window))
                return false;

            // Magic fallback (=somehow reprehensible) to handle items with no assigned ID, e.g. Text(), Image()
            // We build a throwaway ID based on current ID stack + relative AABB of items in window.
            // THE IDENTIFIER WON'T SURVIVE ANY REPOSITIONING OF THE WIDGET, so if your widget moves your dragging operation will be canceled.
            // We don't need to maintain/call ClearActiveID() as releasing the button will early out this function and trigger !ActiveIdIsAlive.
            // Rely on keeping other window.LastItemXXX fields intact.
            source_id = g.LastItemData.ID = window.GetIDFromRectangle(g.LastItemData.Rect);
            bool is_hovered = ItemHoverable(g.LastItemData.Rect, source_id);
            if (is_hovered && g.IO.MouseClicked[mouse_button])
            {
                SetActiveID(source_id, window);
                FocusWindow(window);
            }
            if (g.ActiveId == source_id) // Allow the underlying widget to display/return hovered during the mouse release frame, else we would get a flicker.
                g.ActiveIdAllowOverlap = is_hovered;
        }
        if (g.ActiveId != source_id)
            return false;
        source_parent_id = window.IDStack.back();
        source_drag_active = IsMouseDragging(mouse_button);

        // Disable navigation and key inputs while dragging + cancel existing request if any
        SetActiveIdUsingNavAndKeys();
    }
    else
    {
        window = nil;
        source_id = ImHashStr("#SourceExtern");
        source_drag_active = true;
    }

    if (source_drag_active)
    {
        if (!g.DragDropActive)
        {
            IM_ASSERT(source_id != 0);
            ClearDragDrop();
            ImGuiPayload& payload = g.DragDropPayload;
            payload.SourceId = source_id;
            payload.SourceParentId = source_parent_id;
            g.DragDropActive = true;
            g.DragDropSourceFlags = flags;
            g.DragDropMouseButton = mouse_button;
            if (payload.SourceId == g.ActiveId)
                g.ActiveIdNoClearOnFocusLoss = true;
        }
        g.DragDropSourceFrameCount = g.FrameCount;
        g.DragDropWithinSource = true;

        if (!(flags & ImGuiDragDropFlags_SourceNoPreviewTooltip))
        {
            // Target can request the Source to not display its tooltip (we use a dedicated flag to make this request explicit)
            // We unfortunately can't just modify the source flags and skip the call to BeginTooltip, as caller may be emitting contents.
            BeginTooltip();
            if (g.DragDropAcceptIdPrev && (g.DragDropAcceptFlags & ImGuiDragDropFlags_AcceptNoPreviewTooltip))
            {
                ImGuiWindow* tooltip_window = g.CurrentWindow;
                tooltip_window.Hidden = tooltip_window.SkipItems = true;
                tooltip_window.HiddenFramesCanSkipItems = 1;
            }
        }

        if (!(flags & ImGuiDragDropFlags_SourceNoDisableHover) && !(flags & ImGuiDragDropFlags_SourceExtern))
            g.LastItemData.StatusFlags &= ~ImGuiItemStatusFlags_HoveredRect;

        return true;
    }
    return false;
}

void ImGui::EndDragDropSource()
{
   var g = GImGui
    IM_ASSERT(g.DragDropActive);
    IM_ASSERT(g.DragDropWithinSource && "Not after a BeginDragDropSource()?");

    if (!(g.DragDropSourceFlags & ImGuiDragDropFlags_SourceNoPreviewTooltip))
        EndTooltip();

    // Discard the drag if have not called SetDragDropPayload()
    if (g.DragDropPayload.DataFrameCount == -1)
        ClearDragDrop();
    g.DragDropWithinSource = false;
}

// Use 'cond' to choose to submit payload on drag start or every frame
bool ImGui::SetDragDropPayload(const char* type, const void* data, size_t data_size, ImGuiCond cond)
{
   var g = GImGui
    ImGuiPayload& payload = g.DragDropPayload;
    if (cond == 0)
        cond = ImGuiCond_Always;

    IM_ASSERT(type != nil);
    IM_ASSERT(strlen(type) < IM_ARRAYSIZE(payload.DataType) && "Payload type can be at most 32 characters long");
    IM_ASSERT((data != nil && data_size > 0) || (data == nil && data_size == 0));
    IM_ASSERT(cond == ImGuiCond_Always || cond == ImGuiCond_Once);
    IM_ASSERT(payload.SourceId != 0);                               // Not called between BeginDragDropSource() and EndDragDropSource()

    if (cond == ImGuiCond_Always || payload.DataFrameCount == -1)
    {
        // Copy payload
        ImStrncpy(payload.DataType, type, IM_ARRAYSIZE(payload.DataType));
        g.DragDropPayloadBufHeap.resize(0);
        if (data_size > sizeof(g.DragDropPayloadBufLocal))
        {
            // Store in heap
            g.DragDropPayloadBufHeap.resize((int)data_size);
            payload.Data = g.DragDropPayloadBufHeap.Data;
            memcpy(payload.Data, data, data_size);
        }
        else if (data_size > 0)
        {
            // Store locally
            memset(&g.DragDropPayloadBufLocal, 0, sizeof(g.DragDropPayloadBufLocal));
            payload.Data = g.DragDropPayloadBufLocal;
            memcpy(payload.Data, data, data_size);
        }
        else
        {
            payload.Data = nil;
        }
        payload.DataSize = (int)data_size;
    }
    payload.DataFrameCount = g.FrameCount;

    return (g.DragDropAcceptFrameCount == g.FrameCount) || (g.DragDropAcceptFrameCount == g.FrameCount - 1);
}

bool ImGui::BeginDragDropTargetCustom(const ImRect& bb, ImGuiID id)
{
   var g = GImGui
    if (!g.DragDropActive)
        return false;

    ImGuiWindow* window = g.CurrentWindow;
    ImGuiWindow* hovered_window = g.HoveredWindowUnderMovingWindow;
    if (hovered_window == nil || window.RootWindow != hovered_window.RootWindow)
        return false;
    IM_ASSERT(id != 0);
    if (!IsMouseHoveringRect(bb.Min, bb.Max) || (id == g.DragDropPayload.SourceId))
        return false;
    if (window.SkipItems)
        return false;

    IM_ASSERT(g.DragDropWithinTarget == false);
    g.DragDropTargetRect = bb;
    g.DragDropTargetId = id;
    g.DragDropWithinTarget = true;
    return true;
}

// We don't use BeginDragDropTargetCustom() and duplicate its code because:
// 1) we use LastItemRectHoveredRect which handles items that pushes a temporarily clip rectangle in their code. Calling BeginDragDropTargetCustom(LastItemRect) would not handle them.
// 2) and it's faster. as this code may be very frequently called, we want to early out as fast as we can.
// Also note how the HoveredWindow test is positioned differently in both functions (in both functions we optimize for the cheapest early out case)
bool ImGui::BeginDragDropTarget()
{
   var g = GImGui
    if (!g.DragDropActive)
        return false;

    ImGuiWindow* window = g.CurrentWindow;
    if (!(g.LastItemData.StatusFlags & ImGuiItemStatusFlags_HoveredRect))
        return false;
    ImGuiWindow* hovered_window = g.HoveredWindowUnderMovingWindow;
    if (hovered_window == nil || window.RootWindow != hovered_window.RootWindow)
        return false;

    const ImRect& display_rect = (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_HasDisplayRect) ? g.LastItemData.DisplayRect : g.LastItemData.Rect;
    ImGuiID id = g.LastItemData.ID;
    if (id == 0)
        id = window.GetIDFromRectangle(display_rect);
    if (g.DragDropPayload.SourceId == id)
        return false;

    IM_ASSERT(g.DragDropWithinTarget == false);
    g.DragDropTargetRect = display_rect;
    g.DragDropTargetId = id;
    g.DragDropWithinTarget = true;
    return true;
}

bool ImGui::IsDragDropPayloadBeingAccepted()
{
   var g = GImGui
    return g.DragDropActive && g.DragDropAcceptIdPrev != 0;
}

const ImGuiPayload* ImGui::AcceptDragDropPayload(const char* type, ImGuiDragDropFlags flags)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    ImGuiPayload& payload = g.DragDropPayload;
    IM_ASSERT(g.DragDropActive);                        // Not called between BeginDragDropTarget() and EndDragDropTarget() ?
    IM_ASSERT(payload.DataFrameCount != -1);            // Forgot to call EndDragDropTarget() ?
    if (type != nil && !payload.IsDataType(type))
        return nil;

    // Accept smallest drag target bounding box, this allows us to nest drag targets conveniently without ordering constraints.
    // NB: We currently accept nil id as target. However, overlapping targets requires a unique ID to function!
    const bool was_accepted_previously = (g.DragDropAcceptIdPrev == g.DragDropTargetId);
    ImRect r = g.DragDropTargetRect;
    float r_surface = r.GetWidth() * r.GetHeight();
    if (r_surface <= g.DragDropAcceptIdCurrRectSurface)
    {
        g.DragDropAcceptFlags = flags;
        g.DragDropAcceptIdCurr = g.DragDropTargetId;
        g.DragDropAcceptIdCurrRectSurface = r_surface;
    }

    // Render default drop visuals
    // FIXME-DRAGDROP: Settle on a proper default visuals for drop target.
    payload.Preview = was_accepted_previously;
    flags |= (g.DragDropSourceFlags & ImGuiDragDropFlags_AcceptNoDrawDefaultRect); // Source can also inhibit the preview (useful for external sources that lives for 1 frame)
    if (!(flags & ImGuiDragDropFlags_AcceptNoDrawDefaultRect) && payload.Preview)
        window.DrawList.AddRect(r.Min - ImVec2(3.5f,3.5f), r.Max + ImVec2(3.5f, 3.5f), GetColorU32(ImGuiCol_DragDropTarget), 0.0f, 0, 2.0f);

    g.DragDropAcceptFrameCount = g.FrameCount;
    payload.Delivery = was_accepted_previously && !IsMouseDown(g.DragDropMouseButton); // For extern drag sources affecting os window focus, it's easier to just test !IsMouseDown() instead of IsMouseReleased()
    if (!payload.Delivery && !(flags & ImGuiDragDropFlags_AcceptBeforeDelivery))
        return nil;

    return &payload;
}

const ImGuiPayload* ImGui::GetDragDropPayload()
{
   var g = GImGui
    return g.DragDropActive ? &g.DragDropPayload : nil;
}

// We don't really use/need this now, but added it for the sake of consistency and because we might need it later.
void ImGui::EndDragDropTarget()
{
   var g = GImGui
    IM_ASSERT(g.DragDropActive);
    IM_ASSERT(g.DragDropWithinTarget);
    g.DragDropWithinTarget = false;
}

//-----------------------------------------------------------------------------
// [SECTION] LOGGING/CAPTURING
//-----------------------------------------------------------------------------
// All text output from the interface can be captured into tty/file/clipboard.
// By default, tree nodes are automatically opened during logging.
//-----------------------------------------------------------------------------

// Pass text data straight to log (without being displayed)
static inline void LogTextV(ImGuiContext& g, const char* fmt, va_list args)
{
    if (g.LogFile)
    {
        g.LogBuffer.Buf.resize(0);
        g.LogBuffer.appendfv(fmt, args);
        ImFileWrite(g.LogBuffer.c_str(), sizeof(char), (ImU64)g.LogBuffer.size(), g.LogFile);
    }
    else
    {
        g.LogBuffer.appendfv(fmt, args);
    }
}

void ImGui::LogText(const char* fmt, ...)
{
   var g = GImGui
    if (!g.LogEnabled)
        return;

    va_list args;
    va_start(args, fmt);
    LogTextV(g, fmt, args);
    va_end(args);
}

void ImGui::LogTextV(const char* fmt, va_list args)
{
   var g = GImGui
    if (!g.LogEnabled)
        return;

    LogTextV(g, fmt, args);
}

// Internal version that takes a position to decide on newline placement and pad items according to their depth.
// We split text into individual lines to add current tree level padding
// FIXME: This code is a little complicated perhaps, considering simplifying the whole system.
void ImGui::LogRenderedText(const ImVec2* ref_pos, const char* text, const char* text_end)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    const char* prefix = g.LogNextPrefix;
    const char* suffix = g.LogNextSuffix;
    g.LogNextPrefix = g.LogNextSuffix = nil;

    if (!text_end)
        text_end = FindRenderedTextEnd(text, text_end);

    const bool log_new_line = ref_pos && (ref_pos.y > g.LogLinePosY + g.Style.FramePadding.y + 1);
    if (ref_pos)
        g.LogLinePosY = ref_pos.y;
    if (log_new_line)
    {
        LogText(IM_NEWLINE);
        g.LogLineFirstItem = true;
    }

    if (prefix)
        LogRenderedText(ref_pos, prefix, prefix + strlen(prefix)); // Calculate end ourself to ensure "##" are included here.

    // Re-adjust padding if we have popped out of our starting depth
    if (g.LogDepthRef > window.DC.TreeDepth)
        g.LogDepthRef = window.DC.TreeDepth;
    const int tree_depth = (window.DC.TreeDepth - g.LogDepthRef);

    const char* text_remaining = text;
    for (;;)
    {
        // Split the string. Each new line (after a '\n') is followed by indentation corresponding to the current depth of our log entry.
        // We don't add a trailing \n yet to allow a subsequent item on the same line to be captured.
        const char* line_start = text_remaining;
        const char* line_end = ImStreolRange(line_start, text_end);
        const bool is_last_line = (line_end == text_end);
        if (line_start != line_end || !is_last_line)
        {
            const int line_length = (int)(line_end - line_start);
            const int indentation = g.LogLineFirstItem ? tree_depth * 4 : 1;
            LogText("%*s%.*s", indentation, "", line_length, line_start);
            g.LogLineFirstItem = false;
            if (*line_end == '\n')
            {
                LogText(IM_NEWLINE);
                g.LogLineFirstItem = true;
            }
        }
        if (is_last_line)
            break;
        text_remaining = line_end + 1;
    }

    if (suffix)
        LogRenderedText(ref_pos, suffix, suffix + strlen(suffix));
}

// Start logging/capturing text output
void ImGui::LogBegin(ImGuiLogType type, int auto_open_depth)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;
    IM_ASSERT(g.LogEnabled == false);
    IM_ASSERT(g.LogFile == nil);
    IM_ASSERT(g.LogBuffer.empty());
    g.LogEnabled = true;
    g.LogType = type;
    g.LogNextPrefix = g.LogNextSuffix = nil;
    g.LogDepthRef = window.DC.TreeDepth;
    g.LogDepthToExpand = ((auto_open_depth >= 0) ? auto_open_depth : g.LogDepthToExpandDefault);
    g.LogLinePosY = FLT_MAX;
    g.LogLineFirstItem = true;
}

// Important: doesn't copy underlying data, use carefully (prefix/suffix must be in scope at the time of the next LogRenderedText)
void ImGui::LogSetNextTextDecoration(const char* prefix, const char* suffix)
{
   var g = GImGui
    g.LogNextPrefix = prefix;
    g.LogNextSuffix = suffix;
}

void ImGui::LogToTTY(int auto_open_depth)
{
   var g = GImGui
    if (g.LogEnabled)
        return;
    IM_UNUSED(auto_open_depth);
#ifndef IMGUI_DISABLE_TTY_FUNCTIONS
    LogBegin(ImGuiLogType_TTY, auto_open_depth);
    g.LogFile = stdout;
#endif
}

// Start logging/capturing text output to given file
void ImGui::LogToFile(int auto_open_depth, const char* filename)
{
   var g = GImGui
    if (g.LogEnabled)
        return;

    // FIXME: We could probably open the file in text mode "at", however note that clipboard/buffer logging will still
    // be subject to outputting OS-incompatible carriage return if within strings the user doesn't use IM_NEWLINE.
    // By opening the file in binary mode "ab" we have consistent output everywhere.
    if (!filename)
        filename = g.IO.LogFilename;
    if (!filename || !filename[0])
        return;
    ImFileHandle f = ImFileOpen(filename, "ab");
    if (!f)
    {
        IM_ASSERT(0);
        return;
    }

    LogBegin(ImGuiLogType_File, auto_open_depth);
    g.LogFile = f;
}

// Start logging/capturing text output to clipboard
void ImGui::LogToClipboard(int auto_open_depth)
{
   var g = GImGui
    if (g.LogEnabled)
        return;
    LogBegin(ImGuiLogType_Clipboard, auto_open_depth);
}

void ImGui::LogToBuffer(int auto_open_depth)
{
   var g = GImGui
    if (g.LogEnabled)
        return;
    LogBegin(ImGuiLogType_Buffer, auto_open_depth);
}

void ImGui::LogFinish()
{
   var g = GImGui
    if (!g.LogEnabled)
        return;

    LogText(IM_NEWLINE);
    switch (g.LogType)
    {
    case ImGuiLogType_TTY:
#ifndef IMGUI_DISABLE_TTY_FUNCTIONS
        fflush(g.LogFile);
#endif
        break;
    case ImGuiLogType_File:
        ImFileClose(g.LogFile);
        break;
    case ImGuiLogType_Buffer:
        break;
    case ImGuiLogType_Clipboard:
        if (!g.LogBuffer.empty())
            SetClipboardText(g.LogBuffer.begin());
        break;
    case ImGuiLogType_None:
        IM_ASSERT(0);
        break;
    }

    g.LogEnabled = false;
    g.LogType = ImGuiLogType_None;
    g.LogFile = nil;
    g.LogBuffer.clear();
}

// Helper to display logging buttons
// FIXME-OBSOLETE: We should probably obsolete this and let the user have their own helper (this is one of the oldest function alive!)
void ImGui::LogButtons()
{
   var g = GImGui

    PushID("LogButtons");
#ifndef IMGUI_DISABLE_TTY_FUNCTIONS
    const bool log_to_tty = Button("Log To TTY"); SameLine();
#else
    const bool log_to_tty = false;
#endif
    const bool log_to_file = Button("Log To File"); SameLine();
    const bool log_to_clipboard = Button("Log To Clipboard"); SameLine();
    PushAllowKeyboardFocus(false);
    SetNextItemWidth(80.0f);
    SliderInt("Default Depth", &g.LogDepthToExpandDefault, 0, 9, nil);
    PopAllowKeyboardFocus();
    PopID();

    // Start logging at the end of the function so that the buttons don't appear in the log
    if (log_to_tty)
        LogToTTY();
    if (log_to_file)
        LogToFile();
    if (log_to_clipboard)
        LogToClipboard();
}


//-----------------------------------------------------------------------------
// [SECTION] SETTINGS
//-----------------------------------------------------------------------------
// - UpdateSettings() [Internal]
// - MarkIniSettingsDirty() [Internal]
// - CreateNewWindowSettings() [Internal]
// - FindWindowSettings() [Internal]
// - FindOrCreateWindowSettings() [Internal]
// - FindSettingsHandler() [Internal]
// - ClearIniSettings() [Internal]
// - LoadIniSettingsFromDisk()
// - LoadIniSettingsFromMemory()
// - SaveIniSettingsToDisk()
// - SaveIniSettingsToMemory()
// - WindowSettingsHandler_***() [Internal]
//-----------------------------------------------------------------------------

void ImGui::ClearIniSettings()
{
   var g = GImGui
    g.SettingsIniData.clear();
    for (int handler_n = 0; handler_n < g.SettingsHandlers.Size; handler_n++)
        if (g.SettingsHandlers[handler_n].ClearAllFn)
            g.SettingsHandlers[handler_n].ClearAllFn(&g, &g.SettingsHandlers[handler_n]);
}


static void WindowSettingsHandler_WriteAll(ImGuiContext* ctx, ImGuiSettingsHandler* handler, ImGuiTextBuffer* buf)
{
    // Gather data from windows that were active during this session
    // (if a window wasn't opened in this session we preserve its settings)
    ImGuiContext& g = *ctx;
    for (int i = 0; i != g.Windows.Size; i++)
    {
        ImGuiWindow* window = g.Windows[i];
        if (window.Flags & ImGuiWindowFlags_NoSavedSettings)
            continue;

        ImGuiWindowSettings* settings = (window.SettingsOffset != -1) ? g.SettingsWindows.ptr_from_offset(window.SettingsOffset) : ImGui::FindWindowSettings(window.ID);
        if (!settings)
        {
            settings = ImGui::CreateNewWindowSettings(window.Name);
            window.SettingsOffset = g.SettingsWindows.offset_from_ptr(settings);
        }
        IM_ASSERT(settings.ID == window.ID);
        settings.Pos = ImVec2ih(window.Pos);
        settings.Size = ImVec2ih(window.SizeFull);

        settings.Collapsed = window.Collapsed;
    }

    // Write to text buffer
    buf.reserve(buf.size() + g.SettingsWindows.size() * 6); // ballpark reserve
    for (ImGuiWindowSettings* settings = g.SettingsWindows.begin(); settings != nil; settings = g.SettingsWindows.next_chunk(settings))
    {
        const char* settings_name = settings.GetName();
        buf.appendf("[%s][%s]\n", handler.TypeName, settings_name);
        buf.appendf("Pos=%d,%d\n", settings.Pos.x, settings.Pos.y);
        buf.appendf("Size=%d,%d\n", settings.Size.x, settings.Size.y);
        buf.appendf("Collapsed=%d\n", settings.Collapsed);
        buf.append("\n");
    }
}

//-----------------------------------------------------------------------------
// [SECTION] DOCKING
//-----------------------------------------------------------------------------

// (this section is filled in the 'docking' branch)


//-----------------------------------------------------------------------------
// [SECTION] PLATFORM DEPENDENT HELPERS
//-----------------------------------------------------------------------------

#if defined(_WIN32) && !defined(IMGUI_DISABLE_WIN32_FUNCTIONS) && !defined(IMGUI_DISABLE_WIN32_DEFAULT_CLIPBOARD_FUNCTIONS)

#ifdef _MSC_VER
#pragma comment(lib, "user32")
#pragma comment(lib, "kernel32")
#endif

// Win32 clipboard implementation
// We use g.ClipboardHandlerData for temporary storage to ensure it is freed on Shutdown()
static const char* GetClipboardTextFn_DefaultImpl(void*)
{
   var g = GImGui
    g.ClipboardHandlerData.clear();
    if (!::OpenClipboard(nil))
        return nil;
    HANDLE wbuf_handle = ::GetClipboardData(CF_UNICODETEXT);
    if (wbuf_handle == nil)
    {
        ::CloseClipboard();
        return nil;
    }
    if (const WCHAR* wbuf_global = (const WCHAR*)::GlobalLock(wbuf_handle))
    {
        int buf_len = ::WideCharToMultiByte(CP_UTF8, 0, wbuf_global, -1, nil, 0, nil, nil);
        g.ClipboardHandlerData.resize(buf_len);
        ::WideCharToMultiByte(CP_UTF8, 0, wbuf_global, -1, g.ClipboardHandlerData.Data, buf_len, nil, nil);
    }
    ::GlobalUnlock(wbuf_handle);
    ::CloseClipboard();
    return g.ClipboardHandlerData.Data;
}

static void SetClipboardTextFn_DefaultImpl(void*, const char* text)
{
    if (!::OpenClipboard(nil))
        return;
    const int wbuf_length = ::MultiByteToWideChar(CP_UTF8, 0, text, -1, nil, 0);
    HGLOBAL wbuf_handle = ::GlobalAlloc(GMEM_MOVEABLE, (SIZE_T)wbuf_length * sizeof(WCHAR));
    if (wbuf_handle == nil)
    {
        ::CloseClipboard();
        return;
    }
    WCHAR* wbuf_global = (WCHAR*)::GlobalLock(wbuf_handle);
    ::MultiByteToWideChar(CP_UTF8, 0, text, -1, wbuf_global, wbuf_length);
    ::GlobalUnlock(wbuf_handle);
    ::EmptyClipboard();
    if (::SetClipboardData(CF_UNICODETEXT, wbuf_handle) == nil)
        ::GlobalFree(wbuf_handle);
    ::CloseClipboard();
}

#elif defined(__APPLE__) && TARGET_OS_OSX && defined(IMGUI_ENABLE_OSX_DEFAULT_CLIPBOARD_FUNCTIONS)

#include <Carbon/Carbon.h>  // Use old API to avoid need for separate .mm file
static PasteboardRef main_clipboard = 0;

// OSX clipboard implementation
// If you enable this you will need to add '-framework ApplicationServices' to your linker command-line!
static void SetClipboardTextFn_DefaultImpl(void*, const char* text)
{
    if (!main_clipboard)
        PasteboardCreate(kPasteboardClipboard, &main_clipboard);
    PasteboardClear(main_clipboard);
    CFDataRef cf_data = CFDataCreate(kCFAllocatorDefault, (const UInt8*)text, strlen(text));
    if (cf_data)
    {
        PasteboardPutItemFlavor(main_clipboard, (PasteboardItemID)1, CFSTR("public.utf8-plain-text"), cf_data, 0);
        CFRelease(cf_data);
    }
}

static const char* GetClipboardTextFn_DefaultImpl(void*)
{
    if (!main_clipboard)
        PasteboardCreate(kPasteboardClipboard, &main_clipboard);
    PasteboardSynchronize(main_clipboard);

    ItemCount item_count = 0;
    PasteboardGetItemCount(main_clipboard, &item_count);
    for (ItemCount i = 0; i < item_count; i++)
    {
        PasteboardItemID item_id = 0;
        PasteboardGetItemIdentifier(main_clipboard, i + 1, &item_id);
        CFArrayRef flavor_type_array = 0;
        PasteboardCopyItemFlavors(main_clipboard, item_id, &flavor_type_array);
        for (CFIndex j = 0, nj = CFArrayGetCount(flavor_type_array); j < nj; j++)
        {
            CFDataRef cf_data;
            if (PasteboardCopyItemFlavorData(main_clipboard, item_id, CFSTR("public.utf8-plain-text"), &cf_data) == noErr)
            {
               var g = GImGui
                g.ClipboardHandlerData.clear();
                int length = (int)CFDataGetLength(cf_data);
                g.ClipboardHandlerData.resize(length + 1);
                CFDataGetBytes(cf_data, CFRangeMake(0, length), (UInt8*)g.ClipboardHandlerData.Data);
                g.ClipboardHandlerData[length] = 0;
                CFRelease(cf_data);
                return g.ClipboardHandlerData.Data;
            }
        }
    }
    return nil;
}

#else

// Local Dear ImGui-only clipboard implementation, if user hasn't defined better clipboard handlers.
static const char* GetClipboardTextFn_DefaultImpl(void*)
{
   var g = GImGui
    return g.ClipboardHandlerData.empty() ? nil : g.ClipboardHandlerData.begin();
}

static void SetClipboardTextFn_DefaultImpl(void*, const char* text)
{
   var g = GImGui
    g.ClipboardHandlerData.clear();
    const char* text_end = text + strlen(text);
    g.ClipboardHandlerData.resize((int)(text_end - text) + 1);
    memcpy(&g.ClipboardHandlerData[0], text, (size_t)(text_end - text));
    g.ClipboardHandlerData[(int)(text_end - text)] = 0;
}

#endif

// Win32 API IME support (for Asian languages, etc.)
#if defined(_WIN32) && !defined(IMGUI_DISABLE_WIN32_FUNCTIONS) && !defined(IMGUI_DISABLE_WIN32_DEFAULT_IME_FUNCTIONS)

#include <imm.h>
#ifdef _MSC_VER
#pragma comment(lib, "imm32")
#endif

static void ImeSetInputScreenPosFn_DefaultImpl(int x, int y)
{
    // Notify OS Input Method Editor of text input position
    ImGuiIO& io = ImGui::GetIO();
    if (HWND hwnd = (HWND)io.ImeWindowHandle)
        if (HIMC himc = ::ImmGetContext(hwnd))
        {
            COMPOSITIONFORM cf;
            cf.ptCurrentPos.x = x;
            cf.ptCurrentPos.y = y;
            cf.dwStyle = CFS_FORCE_POSITION;
            ::ImmSetCompositionWindow(himc, &cf);
            ::ImmReleaseContext(hwnd, himc);
        }
}

#else

static void ImeSetInputScreenPosFn_DefaultImpl(int, int) {}

#endif

//-----------------------------------------------------------------------------
// [SECTION] METRICS/DEBUGGER WINDOW
//-----------------------------------------------------------------------------
// - RenderViewportThumbnail() [Internal]
// - RenderViewportsThumbnails() [Internal]
// - MetricsHelpMarker() [Internal]
// - ShowMetricsWindow()
// - DebugNodeColumns() [Internal]
// - DebugNodeDrawList() [Internal]
// - DebugNodeDrawCmdShowMeshAndBoundingBox() [Internal]
// - DebugNodeStorage() [Internal]
// - DebugNodeTabBar() [Internal]
// - DebugNodeViewport() [Internal]
// - DebugNodeWindow() [Internal]
// - DebugNodeWindowSettings() [Internal]
// - DebugNodeWindowsList() [Internal]
//-----------------------------------------------------------------------------

#ifndef IMGUI_DISABLE_METRICS_WINDOW

void ImGui::DebugRenderViewportThumbnail(ImDrawList* draw_list, ImGuiViewportP* viewport, const ImRect& bb)
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    ImVec2 scale = bb.GetSize() / viewport.Size;
    ImVec2 off = bb.Min - viewport.Pos * scale;
    float alpha_mul = 1.0f;
    window.DrawList.AddRectFilled(bb.Min, bb.Max, GetColorU32(ImGuiCol_Border, alpha_mul * 0.40f));
    for (int i = 0; i != g.Windows.Size; i++)
    {
        ImGuiWindow* thumb_window = g.Windows[i];
        if (!thumb_window.WasActive || (thumb_window.Flags & ImGuiWindowFlags_ChildWindow))
            continue;

        ImRect thumb_r = thumb_window.Rect();
        ImRect title_r = thumb_window.TitleBarRect();
        thumb_r = ImRect(ImFloor(off + thumb_r.Min * scale), ImFloor(off +  thumb_r.Max * scale));
        title_r = ImRect(ImFloor(off + title_r.Min * scale), ImFloor(off +  ImVec2(title_r.Max.x, title_r.Min.y) * scale) + ImVec2(0,5)); // Exaggerate title bar height
        thumb_r.ClipWithFull(bb);
        title_r.ClipWithFull(bb);
        const bool window_is_focused = (g.NavWindow && thumb_window.RootWindowForTitleBarHighlight == g.NavWindow.RootWindowForTitleBarHighlight);
        window.DrawList.AddRectFilled(thumb_r.Min, thumb_r.Max, GetColorU32(ImGuiCol_WindowBg, alpha_mul));
        window.DrawList.AddRectFilled(title_r.Min, title_r.Max, GetColorU32(window_is_focused ? ImGuiCol_TitleBgActive : ImGuiCol_TitleBg, alpha_mul));
        window.DrawList.AddRect(thumb_r.Min, thumb_r.Max, GetColorU32(ImGuiCol_Border, alpha_mul));
        window.DrawList.AddText(g.Font, g.FontSize * 1.0f, title_r.Min, GetColorU32(ImGuiCol_Text, alpha_mul), thumb_window.Name, FindRenderedTextEnd(thumb_window.Name));
    }
    draw_list.AddRect(bb.Min, bb.Max, GetColorU32(ImGuiCol_Border, alpha_mul));
}

static void RenderViewportsThumbnails()
{
   var g = GImGui
    ImGuiWindow* window = g.CurrentWindow;

    // We don't display full monitor bounds (we could, but it often looks awkward), instead we display just enough to cover all of our viewports.
    float SCALE = 1.0f / 8.0f;
    ImRect bb_full(FLT_MAX, FLT_MAX, -FLT_MAX, -FLT_MAX);
    for (int n = 0; n < g.Viewports.Size; n++)
        bb_full.Add(g.Viewports[n].GetMainRect());
    ImVec2 p = window.DC.CursorPos;
    ImVec2 off = p - bb_full.Min * SCALE;
    for (int n = 0; n < g.Viewports.Size; n++)
    {
        ImGuiViewportP* viewport = g.Viewports[n];
        ImRect viewport_draw_bb(off + (viewport.Pos) * SCALE, off + (viewport.Pos + viewport.Size) * SCALE);
        ImGui::DebugRenderViewportThumbnail(window.DrawList, viewport, viewport_draw_bb);
    }
    ImGui::Dummy(bb_full.GetSize() * SCALE);
}

// Avoid naming collision with imgui_demo.cpp's HelpMarker() for unity builds.
static void MetricsHelpMarker(const char* desc)
{
    ImGui::TextDisabled("(?)");
    if (ImGui::IsItemHovered())
    {
        ImGui::BeginTooltip();
        ImGui::PushTextWrapPos(ImGui::GetFontSize() * 35.0f);
        ImGui::TextUnformatted(desc);
        ImGui::PopTextWrapPos();
        ImGui::EndTooltip();
    }
}

#ifndef IMGUI_DISABLE_DEMO_WINDOWS
namespace ImGui { void ShowFontAtlas(ImFontAtlas* atlas); }
#endif

void ImGui::ShowMetricsWindow(bool* p_open)
{
    // Tools
    if (TreeNode("Tools"))
    {
        // The Item Picker tool is super useful to visually select an item and break into the call-stack of where it was submitted.
        if (Button("Item Picker.."))
            DebugStartItemPicker();
        SameLine();
        MetricsHelpMarker("Will call the IM_DEBUG_BREAK() macro to break in debugger.\nWarning: If you don't have a debugger attached, this will probably crash.");

        Checkbox("Show windows begin order", &cfg.ShowWindowsBeginOrder);
        Checkbox("Show windows rectangles", &cfg.ShowWindowsRects);
        SameLine();
        SetNextItemWidth(GetFontSize() * 12);
        cfg.ShowWindowsRects |= Combo("##show_windows_rect_type", &cfg.ShowWindowsRectsType, wrt_rects_names, WRT_Count, WRT_Count);
        if (cfg.ShowWindowsRects && g.NavWindow != nil)
        {
            BulletText("'%s':", g.NavWindow.Name);
            Indent();
            for (int rect_n = 0; rect_n < WRT_Count; rect_n++)
            {
                ImRect r = Funcs::GetWindowRect(g.NavWindow, rect_n);
                Text("(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), wrt_rects_names[rect_n]);
            }
            Unindent();
        }
        Checkbox("Show ImDrawCmd mesh when hovering", &cfg.ShowDrawCmdMesh);
        Checkbox("Show ImDrawCmd bounding boxes when hovering", &cfg.ShowDrawCmdBoundingBoxes);

        Checkbox("Show tables rectangles", &cfg.ShowTablesRects);
        SameLine();
        SetNextItemWidth(GetFontSize() * 12);
        cfg.ShowTablesRects |= Combo("##show_table_rects_type", &cfg.ShowTablesRectsType, trt_rects_names, TRT_Count, TRT_Count);
        if (cfg.ShowTablesRects && g.NavWindow != nil)
        {
            for (int table_n = 0; table_n < g.Tables.GetMapSize(); table_n++)
            {
                ImGuiTable* table = g.Tables.TryGetMapData(table_n);
                if (table == nil || table.LastFrameActive < g.FrameCount - 1 || (table.OuterWindow != g.NavWindow && table.InnerWindow != g.NavWindow))
                    continue;

                BulletText("Table 0x%08X (%d columns, in '%s')", table.ID, table.ColumnsCount, table.OuterWindow.Name);
                if (IsItemHovered())
                    GetForegroundDrawList().AddRect(table.OuterRect.Min - ImVec2(1, 1), table.OuterRect.Max + ImVec2(1, 1), IM_COL32(255, 255, 0, 255), 0.0f, 0, 2.0f);
                Indent();
                char buf[128];
                for (int rect_n = 0; rect_n < TRT_Count; rect_n++)
                {
                    if (rect_n >= TRT_ColumnsRect)
                    {
                        if (rect_n != TRT_ColumnsRect && rect_n != TRT_ColumnsClipRect)
                            continue;
                        for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
                        {
                            ImRect r = Funcs::GetTableRect(table, rect_n, column_n);
                            ImFormatString(buf, IM_ARRAYSIZE(buf), "(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) Col %d %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), column_n, trt_rects_names[rect_n]);
                            Selectable(buf);
                            if (IsItemHovered())
                                GetForegroundDrawList().AddRect(r.Min - ImVec2(1, 1), r.Max + ImVec2(1, 1), IM_COL32(255, 255, 0, 255), 0.0f, 0, 2.0f);
                        }
                    }
                    else
                    {
                        ImRect r = Funcs::GetTableRect(table, rect_n, -1);
                        ImFormatString(buf, IM_ARRAYSIZE(buf), "(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), trt_rects_names[rect_n]);
                        Selectable(buf);
                        if (IsItemHovered())
                            GetForegroundDrawList().AddRect(r.Min - ImVec2(1, 1), r.Max + ImVec2(1, 1), IM_COL32(255, 255, 0, 255), 0.0f, 0, 2.0f);
                    }
                }
                Unindent();
            }
        }

        TreePop();
    }

    // Windows
    DebugNodeWindowsList(&g.Windows, "Windows");
    //DebugNodeWindowsList(&g.WindowsFocusOrder, "WindowsFocusOrder");

    // DrawLists
    int drawlist_count = 0;
    for (int viewport_i = 0; viewport_i < g.Viewports.Size; viewport_i++)
        drawlist_count += g.Viewports[viewport_i].DrawDataBuilder.GetDrawListCount();
    if (TreeNode("DrawLists", "DrawLists (%d)", drawlist_count))
    {
        for (int viewport_i = 0; viewport_i < g.Viewports.Size; viewport_i++)
        {
            ImGuiViewportP* viewport = g.Viewports[viewport_i];
            for (int layer_i = 0; layer_i < IM_ARRAYSIZE(viewport.DrawDataBuilder.Layers); layer_i++)
                for (int draw_list_i = 0; draw_list_i < viewport.DrawDataBuilder.Layers[layer_i].Size; draw_list_i++)
                    DebugNodeDrawList(nil, viewport.DrawDataBuilder.Layers[layer_i][draw_list_i], "DrawList");
        }
        TreePop();
    }

    // Viewports
    if (TreeNode("Viewports", "Viewports (%d)", g.Viewports.Size))
    {
        Indent(GetTreeNodeToLabelSpacing());
        RenderViewportsThumbnails();
        Unindent(GetTreeNodeToLabelSpacing());
        for (int i = 0; i < g.Viewports.Size; i++)
            DebugNodeViewport(g.Viewports[i]);
        TreePop();
    }

    // Details for Popups
    if (TreeNode("Popups", "Popups (%d)", g.OpenPopupStack.Size))
    {
        for (int i = 0; i < g.OpenPopupStack.Size; i++)
        {
            ImGuiWindow* window = g.OpenPopupStack[i].Window;
            BulletText("PopupID: %08x, Window: '%s'%s%s", g.OpenPopupStack[i].PopupId, window ? window.Name : "nil", window && (window.Flags & ImGuiWindowFlags_ChildWindow) ? " ChildWindow" : "", window && (window.Flags & ImGuiWindowFlags_ChildMenu) ? " ChildMenu" : "");
        }
        TreePop();
    }

    // Details for TabBars
    if (TreeNode("TabBars", "Tab Bars (%d)", g.TabBars.GetAliveCount()))
    {
        for (int n = 0; n < g.TabBars.GetMapSize(); n++)
            if (ImGuiTabBar* tab_bar = g.TabBars.TryGetMapData(n))
            {
                PushID(tab_bar);
                DebugNodeTabBar(tab_bar, "TabBar");
                PopID();
            }
        TreePop();
    }

    // Details for Tables
    if (TreeNode("Tables", "Tables (%d)", g.Tables.GetAliveCount()))
    {
        for (int n = 0; n < g.Tables.GetMapSize(); n++)
            if (ImGuiTable* table = g.Tables.TryGetMapData(n))
                DebugNodeTable(table);
        TreePop();
    }

    // Details for Fonts
#ifndef IMGUI_DISABLE_DEMO_WINDOWS
    ImFontAtlas* atlas = g.IO.Fonts;
    if (TreeNode("Fonts", "Fonts (%d)", atlas.Fonts.Size))
    {
        ShowFontAtlas(atlas);
        TreePop();
    }
#endif

    // Details for Docking
#ifdef IMGUI_HAS_DOCK
    if (TreeNode("Docking"))
    {
        TreePop();
    }
#endif // #ifdef IMGUI_HAS_DOCK

    // Settings
    if (TreeNode("Settings"))
    {
        if (SmallButton("Clear"))
            ClearIniSettings();
        SameLine();
        if (SmallButton("Save to memory"))
            SaveIniSettingsToMemory();
        SameLine();
        if (SmallButton("Save to disk"))
            SaveIniSettingsToDisk(g.IO.IniFilename);
        SameLine();
        if (g.IO.IniFilename)
            Text("\"%s\"", g.IO.IniFilename);
        else
            TextUnformatted("<nil>");
        Text("SettingsDirtyTimer %.2f", g.SettingsDirtyTimer);
        if (TreeNode("SettingsHandlers", "Settings handlers: (%d)", g.SettingsHandlers.Size))
        {
            for (int n = 0; n < g.SettingsHandlers.Size; n++)
                BulletText("%s", g.SettingsHandlers[n].TypeName);
            TreePop();
        }
        if (TreeNode("SettingsWindows", "Settings packed data: Windows: %d bytes", g.SettingsWindows.size()))
        {
            for (ImGuiWindowSettings* settings = g.SettingsWindows.begin(); settings != nil; settings = g.SettingsWindows.next_chunk(settings))
                DebugNodeWindowSettings(settings);
            TreePop();
        }

        if (TreeNode("SettingsTables", "Settings packed data: Tables: %d bytes", g.SettingsTables.size()))
        {
            for (ImGuiTableSettings* settings = g.SettingsTables.begin(); settings != nil; settings = g.SettingsTables.next_chunk(settings))
                DebugNodeTableSettings(settings);
            TreePop();
        }

#ifdef IMGUI_HAS_DOCK
#endif // #ifdef IMGUI_HAS_DOCK

        if (TreeNode("SettingsIniData", "Settings unpacked data (.ini): %d bytes", g.SettingsIniData.size()))
        {
            InputTextMultiline("##Ini", (char*)(void*)g.SettingsIniData.c_str(), g.SettingsIniData.Buf.Size, ImVec2(-FLT_MIN, GetTextLineHeight() * 20), ImGuiInputTextFlags_ReadOnly);
            TreePop();
        }
        TreePop();
    }

    // Misc Details
    if (TreeNode("Internal state"))
    {
        const char* input_source_names[] = { "None", "Mouse", "Keyboard", "Gamepad", "Nav", "Clipboard" }; IM_ASSERT(IM_ARRAYSIZE(input_source_names) == ImGuiInputSource_COUNT);

        Text("WINDOWING");
        Indent();
        Text("HoveredWindow: '%s'", g.HoveredWindow ? g.HoveredWindow.Name : "nil");
        Text("HoveredWindow.Root: '%s'", g.HoveredWindow ? g.HoveredWindow.RootWindow.Name : "nil");
        Text("HoveredWindowUnderMovingWindow: '%s'", g.HoveredWindowUnderMovingWindow ? g.HoveredWindowUnderMovingWindow.Name : "nil");
        Text("MovingWindow: '%s'", g.MovingWindow ? g.MovingWindow.Name : "nil");
        Unindent();

        Text("ITEMS");
        Indent();
        Text("ActiveId: 0x%08X/0x%08X (%.2f sec), AllowOverlap: %d, Source: %s", g.ActiveId, g.ActiveIdPreviousFrame, g.ActiveIdTimer, g.ActiveIdAllowOverlap, input_source_names[g.ActiveIdSource]);
        Text("ActiveIdWindow: '%s'", g.ActiveIdWindow ? g.ActiveIdWindow.Name : "nil");
        Text("ActiveIdUsing: Wheel: %d, NavDirMask: %X, NavInputMask: %X, KeyInputMask: %llX", g.ActiveIdUsingMouseWheel, g.ActiveIdUsingNavDirMask, g.ActiveIdUsingNavInputMask, g.ActiveIdUsingKeyInputMask);
        Text("HoveredId: 0x%08X (%.2f sec), AllowOverlap: %d", g.HoveredIdPreviousFrame, g.HoveredIdTimer, g.HoveredIdAllowOverlap); // Not displaying g.HoveredId as it is update mid-frame
        Text("DragDrop: %d, SourceId = 0x%08X, Payload \"%s\" (%d bytes)", g.DragDropActive, g.DragDropPayload.SourceId, g.DragDropPayload.DataType, g.DragDropPayload.DataSize);
        Unindent();

        Text("NAV,FOCUS");
        Indent();
        Text("NavWindow: '%s'", g.NavWindow ? g.NavWindow.Name : "nil");
        Text("NavId: 0x%08X, NavLayer: %d", g.NavId, g.NavLayer);
        Text("NavInputSource: %s", input_source_names[g.NavInputSource]);
        Text("NavActive: %d, NavVisible: %d", g.IO.NavActive, g.IO.NavVisible);
        Text("NavActivateId: 0x%08X, NavInputId: 0x%08X", g.NavActivateId, g.NavInputId);
        Text("NavDisableHighlight: %d, NavDisableMouseHover: %d", g.NavDisableHighlight, g.NavDisableMouseHover);
        Text("NavFocusScopeId = 0x%08X", g.NavFocusScopeId);
        Text("NavWindowingTarget: '%s'", g.NavWindowingTarget ? g.NavWindowingTarget.Name : "nil");
        Unindent();

        TreePop();
    }

    // Overlay: Display windows Rectangles and Begin Order
    if (cfg.ShowWindowsRects || cfg.ShowWindowsBeginOrder)
    {
        for (int n = 0; n < g.Windows.Size; n++)
        {
            ImGuiWindow* window = g.Windows[n];
            if (!window.WasActive)
                continue;
            ImDrawList* draw_list = GetForegroundDrawList(window);
            if (cfg.ShowWindowsRects)
            {
                ImRect r = Funcs::GetWindowRect(window, cfg.ShowWindowsRectsType);
                draw_list.AddRect(r.Min, r.Max, IM_COL32(255, 0, 128, 255));
            }
            if (cfg.ShowWindowsBeginOrder && !(window.Flags & ImGuiWindowFlags_ChildWindow))
            {
                char buf[32];
                ImFormatString(buf, IM_ARRAYSIZE(buf), "%d", window.BeginOrderWithinContext);
                float font_size = GetFontSize();
                draw_list.AddRectFilled(window.Pos, window.Pos + ImVec2(font_size, font_size), IM_COL32(200, 100, 100, 255));
                draw_list.AddText(window.Pos, IM_COL32(255, 255, 255, 255), buf);
            }
        }
    }

    // Overlay: Display Tables Rectangles
    if (cfg.ShowTablesRects)
    {
        for (int table_n = 0; table_n < g.Tables.GetMapSize(); table_n++)
        {
            ImGuiTable* table = g.Tables.TryGetMapData(table_n);
            if (table == nil || table.LastFrameActive < g.FrameCount - 1)
                continue;
            ImDrawList* draw_list = GetForegroundDrawList(table.OuterWindow);
            if (cfg.ShowTablesRectsType >= TRT_ColumnsRect)
            {
                for (int column_n = 0; column_n < table.ColumnsCount; column_n++)
                {
                    ImRect r = Funcs::GetTableRect(table, cfg.ShowTablesRectsType, column_n);
                    ImU32 col = (table.HoveredColumnBody == column_n) ? IM_COL32(255, 255, 128, 255) : IM_COL32(255, 0, 128, 255);
                    float thickness = (table.HoveredColumnBody == column_n) ? 3.0f : 1.0f;
                    draw_list.AddRect(r.Min, r.Max, col, 0.0f, 0, thickness);
                }
            }
            else
            {
                ImRect r = Funcs::GetTableRect(table, cfg.ShowTablesRectsType, -1);
                draw_list.AddRect(r.Min, r.Max, IM_COL32(255, 0, 128, 255));
            }
        }
    }

#ifdef IMGUI_HAS_DOCK
    // Overlay: Display Docking info
    if (show_docking_nodes && g.IO.KeyCtrl)
    {
    }
#endif // #ifdef IMGUI_HAS_DOCK

    End();
}

// [DEBUG] List fonts in a font atlas and display its texture
void ImGui::ShowFontAtlas(ImFontAtlas* atlas)
{
    for (int i = 0; i < atlas.Fonts.Size; i++)
    {
        ImFont* font = atlas.Fonts[i];
        PushID(font);
        DebugNodeFont(font);
        PopID();
    }
    if (TreeNode("Atlas texture", "Atlas texture (%dx%d pixels)", atlas.TexWidth, atlas.TexHeight))
    {
        ImVec4 tint_col = ImVec4(1.0f, 1.0f, 1.0f, 1.0f);
        ImVec4 border_col = ImVec4(1.0f, 1.0f, 1.0f, 0.5f);
        Image(atlas.TexID, ImVec2((float)atlas.TexWidth, (float)atlas.TexHeight), ImVec2(0.0f, 0.0f), ImVec2(1.0f, 1.0f), tint_col, border_col);
        TreePop();
    }
}

// [DEBUG] Display contents of Columns
void ImGui::DebugNodeColumns(ImGuiOldColumns* columns)
{
    if (!TreeNode((void*)(uintptr_t)columns.ID, "Columns Id: 0x%08X, Count: %d, Flags: 0x%04X", columns.ID, columns.Count, columns.Flags))
        return;
    BulletText("Width: %.1f (MinX: %.1f, MaxX: %.1f)", columns.OffMaxX - columns.OffMinX, columns.OffMinX, columns.OffMaxX);
    for (int column_n = 0; column_n < columns.Columns.Size; column_n++)
        BulletText("Column %02d: OffsetNorm %.3f (= %.1f px)", column_n, columns.Columns[column_n].OffsetNorm, GetColumnOffsetFromNorm(columns, columns.Columns[column_n].OffsetNorm));
    TreePop();
}

// [DEBUG] Display contents of ImDrawList
void ImGui::DebugNodeDrawList(ImGuiWindow* window, const ImDrawList* draw_list, const char* label)
{
   var g = GImGui
    ImGuiMetricsConfig* cfg = &g.DebugMetricsConfig;
    int cmd_count = draw_list.CmdBuffer.Size;
    if (cmd_count > 0 && draw_list.CmdBuffer.back().ElemCount == 0 && draw_list.CmdBuffer.back().UserCallback == nil)
        cmd_count--;
    bool node_open = TreeNode(draw_list, "%s: '%s' %d vtx, %d indices, %d cmds", label, draw_list._OwnerName ? draw_list._OwnerName : "", draw_list.VtxBuffer.Size, draw_list.IdxBuffer.Size, cmd_count);
    if (draw_list == GetWindowDrawList())
    {
        SameLine();
        TextColored(ImVec4(1.0f, 0.4f, 0.4f, 1.0f), "CURRENTLY APPENDING"); // Can't display stats for active draw list! (we don't have the data double-buffered)
        if (node_open)
            TreePop();
        return;
    }

    ImDrawList* fg_draw_list = GetForegroundDrawList(window); // Render additional visuals into the top-most draw list
    if (window && IsItemHovered())
        fg_draw_list.AddRect(window.Pos, window.Pos + window.Size, IM_COL32(255, 255, 0, 255));
    if (!node_open)
        return;

    if (window && !window.WasActive)
        TextDisabled("Warning: owning Window is inactive. This DrawList is not being rendered!");

    for (const ImDrawCmd* pcmd = draw_list.CmdBuffer.Data; pcmd < draw_list.CmdBuffer.Data + cmd_count; pcmd++)
    {
        if (pcmd.UserCallback)
        {
            BulletText("Callback %p, user_data %p", pcmd.UserCallback, pcmd.UserCallbackData);
            continue;
        }

        char buf[300];
        ImFormatString(buf, IM_ARRAYSIZE(buf), "DrawCmd:%5d tris, Tex 0x%p, ClipRect (%4.0f,%4.0f)-(%4.0f,%4.0f)",
            pcmd.ElemCount / 3, (void*)(intptr_t)pcmd.TextureId,
            pcmd.ClipRect.x, pcmd.ClipRect.y, pcmd.ClipRect.z, pcmd.ClipRect.w);
        bool pcmd_node_open = TreeNode((void*)(pcmd - draw_list.CmdBuffer.begin()), "%s", buf);
        if (IsItemHovered() && (cfg.ShowDrawCmdMesh || cfg.ShowDrawCmdBoundingBoxes) && fg_draw_list)
            DebugNodeDrawCmdShowMeshAndBoundingBox(fg_draw_list, draw_list, pcmd, cfg.ShowDrawCmdMesh, cfg.ShowDrawCmdBoundingBoxes);
        if (!pcmd_node_open)
            continue;

        // Calculate approximate coverage area (touched pixel count)
        // This will be in pixels squared as long there's no post-scaling happening to the renderer output.
        const ImDrawIdx* idx_buffer = (draw_list.IdxBuffer.Size > 0) ? draw_list.IdxBuffer.Data : nil;
        const ImDrawVert* vtx_buffer = draw_list.VtxBuffer.Data + pcmd.VtxOffset;
        float total_area = 0.0f;
        for (unsigned int idx_n = pcmd.IdxOffset; idx_n < pcmd.IdxOffset + pcmd.ElemCount; )
        {
            ImVec2 triangle[3];
            for (int n = 0; n < 3; n++, idx_n++)
                triangle[n] = vtx_buffer[idx_buffer ? idx_buffer[idx_n] : idx_n].pos;
            total_area += ImTriangleArea(triangle[0], triangle[1], triangle[2]);
        }

        // Display vertex information summary. Hover to get all triangles drawn in wire-frame
        ImFormatString(buf, IM_ARRAYSIZE(buf), "Mesh: ElemCount: %d, VtxOffset: +%d, IdxOffset: +%d, Area: ~%0.f px", pcmd.ElemCount, pcmd.VtxOffset, pcmd.IdxOffset, total_area);
        Selectable(buf);
        if (IsItemHovered() && fg_draw_list)
            DebugNodeDrawCmdShowMeshAndBoundingBox(fg_draw_list, draw_list, pcmd, true, false);

        // Display individual triangles/vertices. Hover on to get the corresponding triangle highlighted.
        ImGuiListClipper clipper;
        clipper.Begin(pcmd.ElemCount / 3); // Manually coarse clip our print out of individual vertices to save CPU, only items that may be visible.
        while (clipper.Step())
            for (int prim = clipper.DisplayStart, idx_i = pcmd.IdxOffset + clipper.DisplayStart * 3; prim < clipper.DisplayEnd; prim++)
            {
                char* buf_p = buf, * buf_end = buf + IM_ARRAYSIZE(buf);
                ImVec2 triangle[3];
                for (int n = 0; n < 3; n++, idx_i++)
                {
                    const ImDrawVert& v = vtx_buffer[idx_buffer ? idx_buffer[idx_i] : idx_i];
                    triangle[n] = v.pos;
                    buf_p += ImFormatString(buf_p, buf_end - buf_p, "%s %04d: pos (%8.2f,%8.2f), uv (%.6f,%.6f), col %08X\n",
                        (n == 0) ? "Vert:" : "     ", idx_i, v.pos.x, v.pos.y, v.uv.x, v.uv.y, v.col);
                }

                Selectable(buf, false);
                if (fg_draw_list && IsItemHovered())
                {
                    ImDrawListFlags backup_flags = fg_draw_list.Flags;
                    fg_draw_list.Flags &= ~ImDrawListFlags_AntiAliasedLines; // Disable AA on triangle outlines is more readable for very large and thin triangles.
                    fg_draw_list.AddPolyline(triangle, 3, IM_COL32(255, 255, 0, 255), ImDrawFlags_Closed, 1.0f);
                    fg_draw_list.Flags = backup_flags;
                }
            }
        TreePop();
    }
    TreePop();
}

// [DEBUG] Display mesh/aabb of a ImDrawCmd
void ImGui::DebugNodeDrawCmdShowMeshAndBoundingBox(ImDrawList* out_draw_list, const ImDrawList* draw_list, const ImDrawCmd* draw_cmd, bool show_mesh, bool show_aabb)
{
    IM_ASSERT(show_mesh || show_aabb);

    // Draw wire-frame version of all triangles
    ImRect clip_rect = draw_cmd.ClipRect;
    ImRect vtxs_rect(FLT_MAX, FLT_MAX, -FLT_MAX, -FLT_MAX);
    ImDrawListFlags backup_flags = out_draw_list.Flags;
    out_draw_list.Flags &= ~ImDrawListFlags_AntiAliasedLines; // Disable AA on triangle outlines is more readable for very large and thin triangles.
    for (unsigned int idx_n = draw_cmd.IdxOffset, idx_end = draw_cmd.IdxOffset + draw_cmd.ElemCount; idx_n < idx_end; )
    {
        ImDrawIdx* idx_buffer = (draw_list.IdxBuffer.Size > 0) ? draw_list.IdxBuffer.Data : nil; // We don't hold on those pointers past iterations as .AddPolyline() may invalidate them if out_draw_list==draw_list
        ImDrawVert* vtx_buffer = draw_list.VtxBuffer.Data + draw_cmd.VtxOffset;

        ImVec2 triangle[3];
        for (int n = 0; n < 3; n++, idx_n++)
            vtxs_rect.Add((triangle[n] = vtx_buffer[idx_buffer ? idx_buffer[idx_n] : idx_n].pos));
        if (show_mesh)
            out_draw_list.AddPolyline(triangle, 3, IM_COL32(255, 255, 0, 255), ImDrawFlags_Closed, 1.0f); // In yellow: mesh triangles
    }
    // Draw bounding boxes
    if (show_aabb)
    {
        out_draw_list.AddRect(ImFloor(clip_rect.Min), ImFloor(clip_rect.Max), IM_COL32(255, 0, 255, 255)); // In pink: clipping rectangle submitted to GPU
        out_draw_list.AddRect(ImFloor(vtxs_rect.Min), ImFloor(vtxs_rect.Max), IM_COL32(0, 255, 255, 255)); // In cyan: bounding box of triangles
    }
    out_draw_list.Flags = backup_flags;
}

// [DEBUG] Display details for a single font, called by ShowStyleEditor().
void ImGui::DebugNodeFont(ImFont* font)
{
    bool opened = TreeNode(font, "Font: \"%s\"\n%.2f px, %d glyphs, %d file(s)",
        font.ConfigData ? font.ConfigData[0].Name : "", font.FontSize, font.Glyphs.Size, font.ConfigDataCount);
    SameLine();
    if (SmallButton("Set as default"))
        GetIO().FontDefault = font;
    if (!opened)
        return;

    // Display preview text
    PushFont(font);
    Text("The quick brown fox jumps over the lazy dog");
    PopFont();

    // Display details
    SetNextItemWidth(GetFontSize() * 8);
    DragFloat("Font scale", &font.Scale, 0.005f, 0.3f, 2.0f, "%.1f");
    SameLine(); MetricsHelpMarker(
        "Note than the default embedded font is NOT meant to be scaled.\n\n"
        "Font are currently rendered into bitmaps at a given size at the time of building the atlas. "
        "You may oversample them to get some flexibility with scaling. "
        "You can also render at multiple sizes and select which one to use at runtime.\n\n"
        "(Glimmer of hope: the atlas system will be rewritten in the future to make scaling more flexible.)");
    Text("Ascent: %f, Descent: %f, Height: %f", font.Ascent, font.Descent, font.Ascent - font.Descent);
    char c_str[5];
    Text("Fallback character: '%s' (U+%04X)", ImTextCharToUtf8(c_str, font.FallbackChar), font.FallbackChar);
    Text("Ellipsis character: '%s' (U+%04X)", ImTextCharToUtf8(c_str, font.EllipsisChar), font.EllipsisChar);
    const int surface_sqrt = (int)ImSqrt((float)font.MetricsTotalSurface);
    Text("Texture Area: about %d px ~%dx%d px", font.MetricsTotalSurface, surface_sqrt, surface_sqrt);
    for (int config_i = 0; config_i < font.ConfigDataCount; config_i++)
        if (font.ConfigData)
            if (const ImFontConfig* cfg = &font.ConfigData[config_i])
                BulletText("Input %d: \'%s\', Oversample: (%d,%d), PixelSnapH: %d, Offset: (%.1f,%.1f)",
                    config_i, cfg.Name, cfg.OversampleH, cfg.OversampleV, cfg.PixelSnapH, cfg.GlyphOffset.x, cfg.GlyphOffset.y);

    // Display all glyphs of the fonts in separate pages of 256 characters
    if (TreeNode("Glyphs", "Glyphs (%d)", font.Glyphs.Size))
    {
        ImDrawList* draw_list = GetWindowDrawList();
        const ImU32 glyph_col = GetColorU32(ImGuiCol_Text);
        const float cell_size = font.FontSize * 1;
        const float cell_spacing = GetStyle().ItemSpacing.y;
        for (unsigned int base = 0; base <= IM_UNICODE_CODEPOINT_MAX; base += 256)
        {
            // Skip ahead if a large bunch of glyphs are not present in the font (test in chunks of 4k)
            // This is only a small optimization to reduce the number of iterations when IM_UNICODE_MAX_CODEPOINT
            // is large // (if ImWchar==ImWchar32 we will do at least about 272 queries here)
            if (!(base & 4095) && font.IsGlyphRangeUnused(base, base + 4095))
            {
                base += 4096 - 256;
                continue;
            }

            int count = 0;
            for (unsigned int n = 0; n < 256; n++)
                if (font.FindGlyphNoFallback((ImWchar)(base + n)))
                    count++;
            if (count <= 0)
                continue;
            if (!TreeNode((void*)(intptr_t)base, "U+%04X..U+%04X (%d %s)", base, base + 255, count, count > 1 ? "glyphs" : "glyph"))
                continue;

            // Draw a 16x16 grid of glyphs
            ImVec2 base_pos = GetCursorScreenPos();
            for (unsigned int n = 0; n < 256; n++)
            {
                // We use ImFont::RenderChar as a shortcut because we don't have UTF-8 conversion functions
                // available here and thus cannot easily generate a zero-terminated UTF-8 encoded string.
                ImVec2 cell_p1(base_pos.x + (n % 16) * (cell_size + cell_spacing), base_pos.y + (n / 16) * (cell_size + cell_spacing));
                ImVec2 cell_p2(cell_p1.x + cell_size, cell_p1.y + cell_size);
                const ImFontGlyph* glyph = font.FindGlyphNoFallback((ImWchar)(base + n));
                draw_list.AddRect(cell_p1, cell_p2, glyph ? IM_COL32(255, 255, 255, 100) : IM_COL32(255, 255, 255, 50));
                if (glyph)
                    font.RenderChar(draw_list, cell_size, cell_p1, glyph_col, (ImWchar)(base + n));
                if (glyph && IsMouseHoveringRect(cell_p1, cell_p2))
                {
                    BeginTooltip();
                    Text("Codepoint: U+%04X", base + n);
                    Separator();
                    Text("Visible: %d", glyph.Visible);
                    Text("AdvanceX: %.1f", glyph.AdvanceX);
                    Text("Pos: (%.2f,%.2f).(%.2f,%.2f)", glyph.X0, glyph.Y0, glyph.X1, glyph.Y1);
                    Text("UV: (%.3f,%.3f).(%.3f,%.3f)", glyph.U0, glyph.V0, glyph.U1, glyph.V1);
                    EndTooltip();
                }
            }
            Dummy(ImVec2((cell_size + cell_spacing) * 16, (cell_size + cell_spacing) * 16));
            TreePop();
        }
        TreePop();
    }
    TreePop();
}

// [DEBUG] Display contents of ImGuiStorage
void ImGui::DebugNodeStorage(ImGuiStorage* storage, const char* label)
{
    if (!TreeNode(label, "%s: %d entries, %d bytes", label, storage.Data.Size, storage.Data.size_in_bytes()))
        return;
    for (int n = 0; n < storage.Data.Size; n++)
    {
        const ImGuiStorage::ImGuiStoragePair& p = storage.Data[n];
        BulletText("Key 0x%08X Value { i: %d }", p.key, p.val_i); // Important: we currently don't store a type, real value may not be integer.
    }
    TreePop();
}

// [DEBUG] Display contents of ImGuiTabBar
void ImGui::DebugNodeTabBar(ImGuiTabBar* tab_bar, const char* label)
{
    // Standalone tab bars (not associated to docking/windows functionality) currently hold no discernible strings.
    char buf[256];
    char* p = buf;
    const char* buf_end = buf + IM_ARRAYSIZE(buf);
    const bool is_active = (tab_bar.PrevFrameVisible >= GetFrameCount() - 2);
    p += ImFormatString(p, buf_end - p, "%s 0x%08X (%d tabs)%s", label, tab_bar.ID, tab_bar.Tabs.Size, is_active ? "" : " *Inactive*");
    p += ImFormatString(p, buf_end - p, "  { ");
    for (int tab_n = 0; tab_n < ImMin(tab_bar.Tabs.Size, 3); tab_n++)
    {
        ImGuiTabItem* tab = &tab_bar.Tabs[tab_n];
        p += ImFormatString(p, buf_end - p, "%s'%s'",
            tab_n > 0 ? ", " : "", (tab.NameOffset != -1) ? tab_bar.GetTabName(tab) : "???");
    }
    p += ImFormatString(p, buf_end - p, (tab_bar.Tabs.Size > 3) ? " ... }" : " } ");
    if (!is_active) { PushStyleColor(ImGuiCol_Text, GetStyleColorVec4(ImGuiCol_TextDisabled)); }
    bool open = TreeNode(label, "%s", buf);
    if (!is_active) { PopStyleColor(); }
    if (is_active && IsItemHovered())
    {
        ImDrawList* draw_list = GetForegroundDrawList();
        draw_list.AddRect(tab_bar.BarRect.Min, tab_bar.BarRect.Max, IM_COL32(255, 255, 0, 255));
        draw_list.AddLine(ImVec2(tab_bar.ScrollingRectMinX, tab_bar.BarRect.Min.y), ImVec2(tab_bar.ScrollingRectMinX, tab_bar.BarRect.Max.y), IM_COL32(0, 255, 0, 255));
        draw_list.AddLine(ImVec2(tab_bar.ScrollingRectMaxX, tab_bar.BarRect.Min.y), ImVec2(tab_bar.ScrollingRectMaxX, tab_bar.BarRect.Max.y), IM_COL32(0, 255, 0, 255));
    }
    if (open)
    {
        for (int tab_n = 0; tab_n < tab_bar.Tabs.Size; tab_n++)
        {
            const ImGuiTabItem* tab = &tab_bar.Tabs[tab_n];
            PushID(tab);
            if (SmallButton("<")) { TabBarQueueReorder(tab_bar, tab, -1); } SameLine(0, 2);
            if (SmallButton(">")) { TabBarQueueReorder(tab_bar, tab, +1); } SameLine();
            Text("%02d%c Tab 0x%08X '%s' Offset: %.1f, Width: %.1f/%.1f",
                tab_n, (tab.ID == tab_bar.SelectedTabId) ? '*' : ' ', tab.ID, (tab.NameOffset != -1) ? tab_bar.GetTabName(tab) : "???", tab.Offset, tab.Width, tab.ContentWidth);
            PopID();
        }
        TreePop();
    }
}

void ImGui::DebugNodeViewport(ImGuiViewportP* viewport)
{
    SetNextItemOpen(true, ImGuiCond_Once);
    if (TreeNode("viewport0", "Viewport #%d", 0))
    {
        ImGuiWindowFlags flags = viewport.Flags;
        BulletText("Main Pos: (%.0f,%.0f), Size: (%.0f,%.0f)\nWorkArea Offset Left: %.0f Top: %.0f, Right: %.0f, Bottom: %.0f",
            viewport.Pos.x, viewport.Pos.y, viewport.Size.x, viewport.Size.y,
            viewport.WorkOffsetMin.x, viewport.WorkOffsetMin.y, viewport.WorkOffsetMax.x, viewport.WorkOffsetMax.y);
        BulletText("Flags: 0x%04X =%s%s%s", viewport.Flags,
            (flags & ImGuiViewportFlags_IsPlatformWindow)  ? " IsPlatformWindow"  : "",
            (flags & ImGuiViewportFlags_IsPlatformMonitor) ? " IsPlatformMonitor" : "",
            (flags & ImGuiViewportFlags_OwnedByApp)        ? " OwnedByApp"        : "");
        for (int layer_i = 0; layer_i < IM_ARRAYSIZE(viewport.DrawDataBuilder.Layers); layer_i++)
            for (int draw_list_i = 0; draw_list_i < viewport.DrawDataBuilder.Layers[layer_i].Size; draw_list_i++)
                DebugNodeDrawList(nil, viewport.DrawDataBuilder.Layers[layer_i][draw_list_i], "DrawList");
        TreePop();
    }
}

void ImGui::DebugNodeWindow(ImGuiWindow* window, const char* label)
{
    if (window == nil)
    {
        BulletText("%s: nil", label);
        return;
    }

   var g = GImGui
    const bool is_active = window.WasActive;
    ImGuiTreeNodeFlags tree_node_flags = (window == g.NavWindow) ? ImGuiTreeNodeFlags_Selected : ImGuiTreeNodeFlags_None;
    if (!is_active) { PushStyleColor(ImGuiCol_Text, GetStyleColorVec4(ImGuiCol_TextDisabled)); }
    const bool open = TreeNodeEx(label, tree_node_flags, "%s '%s'%s", label, window.Name, is_active ? "" : " *Inactive*");
    if (!is_active) { PopStyleColor(); }
    if (IsItemHovered() && is_active)
        GetForegroundDrawList(window).AddRect(window.Pos, window.Pos + window.Size, IM_COL32(255, 255, 0, 255));
    if (!open)
        return;

    if (window.MemoryCompacted)
        TextDisabled("Note: some memory buffers have been compacted/freed.");

    var flags = window.Flags
    DebugNodeDrawList(window, window.DrawList, "DrawList");
    BulletText("Pos: (%.1f,%.1f), Size: (%.1f,%.1f), ContentSize (%.1f,%.1f) Ideal (%.1f,%.1f)", window.Pos.x, window.Pos.y, window.Size.x, window.Size.y, window.ContentSize.x, window.ContentSize.y, window.ContentSizeIdeal.x, window.ContentSizeIdeal.y);
    BulletText("Flags: 0x%08X (%s%s%s%s%s%s%s%s%s..)", flags,
        (flags & ImGuiWindowFlags_ChildWindow)  ? "Child " : "",      (flags & ImGuiWindowFlags_Tooltip)     ? "Tooltip "   : "",  (flags & ImGuiWindowFlags_Popup) ? "Popup " : "",
        (flags & ImGuiWindowFlags_Modal)        ? "Modal " : "",      (flags & ImGuiWindowFlags_ChildMenu)   ? "ChildMenu " : "",  (flags & ImGuiWindowFlags_NoSavedSettings) ? "NoSavedSettings " : "",
        (flags & ImGuiWindowFlags_NoMouseInputs)? "NoMouseInputs":"", (flags & ImGuiWindowFlags_NoNavInputs) ? "NoNavInputs" : "", (flags & ImGuiWindowFlags_AlwaysAutoResize) ? "AlwaysAutoResize" : "");
    BulletText("Scroll: (%.2f/%.2f,%.2f/%.2f) Scrollbar:%s%s", window.Scroll.x, window.ScrollMax.x, window.Scroll.y, window.ScrollMax.y, window.ScrollbarX ? "X" : "", window.ScrollbarY ? "Y" : "");
    BulletText("Active: %d/%d, WriteAccessed: %d, BeginOrderWithinContext: %d", window.Active, window.WasActive, window.WriteAccessed, (window.Active || window.WasActive) ? window.BeginOrderWithinContext : -1);
    BulletText("Appearing: %d, Hidden: %d (CanSkip %d Cannot %d), SkipItems: %d", window.Appearing, window.Hidden, window.HiddenFramesCanSkipItems, window.HiddenFramesCannotSkipItems, window.SkipItems);
    for (int layer = 0; layer < ImGuiNavLayer_COUNT; layer++)
    {
        ImRect r = window.NavRectRel[layer];
        if (r.Min.x >= r.Max.y && r.Min.y >= r.Max.y)
        {
            BulletText("NavLastIds[%d]: 0x%08X", layer, window.NavLastIds[layer]);
            continue;
        }
        BulletText("NavLastIds[%d]: 0x%08X at +(%.1f,%.1f)(%.1f,%.1f)", layer, window.NavLastIds[layer], r.Min.x, r.Min.y, r.Max.x, r.Max.y);
        if (IsItemHovered())
            GetForegroundDrawList(window).AddRect(r.Min + window.Pos, r.Max + window.Pos, IM_COL32(255, 255, 0, 255));
    }
    BulletText("NavLayersActiveMask: %X, NavLastChildNavWindow: %s", window.DC.NavLayersActiveMask, window.NavLastChildNavWindow ? window.NavLastChildNavWindow.Name : "nil");
    if (window.RootWindow != window)       { DebugNodeWindow(window.RootWindow, "RootWindow"); }
    if (window.ParentWindow != nil)       { DebugNodeWindow(window.ParentWindow, "ParentWindow"); }
    if (window.DC.ChildWindows.Size > 0)   { DebugNodeWindowsList(&window.DC.ChildWindows, "ChildWindows"); }
    if (window.ColumnsStorage.Size > 0 && TreeNode("Columns", "Columns sets (%d)", window.ColumnsStorage.Size))
    {
        for (int n = 0; n < window.ColumnsStorage.Size; n++)
            DebugNodeColumns(&window.ColumnsStorage[n]);
        TreePop();
    }
    DebugNodeStorage(&window.StateStorage, "Storage");
    TreePop();
}

void ImGui::DebugNodeWindowSettings(ImGuiWindowSettings* settings)
{
    Text("0x%08X \"%s\" Pos (%d,%d) Size (%d,%d) Collapsed=%d",
        settings.ID, settings.GetName(), settings.Pos.x, settings.Pos.y, settings.Size.x, settings.Size.y, settings.Collapsed);
}

void ImGui::DebugNodeWindowsList(ImVector<ImGuiWindow*>* windows, const char* label)
{
    if (!TreeNode(label, "%s (%d)", label, windows.Size))
        return;
    Text("(In front-to-back order:)");
    for (int i = windows.Size - 1; i >= 0; i--) // Iterate front to back
    {
        PushID((*windows)[i]);
        DebugNodeWindow((*windows)[i], "Window");
        PopID();
    }
    TreePop();
}

#else

void ImGui::ShowMetricsWindow(bool*) {}
void ImGui::ShowFontAtlas(ImFontAtlas*) {}
void ImGui::DebugNodeColumns(ImGuiOldColumns*) {}
void ImGui::DebugNodeDrawList(ImGuiWindow*, const ImDrawList*, const char*) {}
void ImGui::DebugNodeDrawCmdShowMeshAndBoundingBox(ImDrawList*, const ImDrawList*, const ImDrawCmd*, bool, bool) {}
void ImGui::DebugNodeFont(ImFont*) {}
void ImGui::DebugNodeStorage(ImGuiStorage*, const char*) {}
void ImGui::DebugNodeTabBar(ImGuiTabBar*, const char*) {}
void ImGui::DebugNodeWindow(ImGuiWindow*, const char*) {}
void ImGui::DebugNodeWindowSettings(ImGuiWindowSettings*) {}
void ImGui::DebugNodeWindowsList(ImVector<ImGuiWindow*>*, const char*) {}
void ImGui::DebugNodeViewport(ImGuiViewportP*) {}

#endif

//-----------------------------------------------------------------------------

// Include imgui_user.inl at the end of imgui.cpp to access private data/functions that aren't exposed.
// Prefer just including imgui_internal.h from your code rather than using this define. If a declaration is missing from imgui_internal.h add it or request it on the github.
#ifdef IMGUI_INCLUDE_IMGUI_USER_INL
#include "imgui_user.inl"
#endif

//-----------------------------------------------------------------------------

#endif // #ifndef IMGUI_DISABLE
