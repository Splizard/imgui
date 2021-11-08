
#ifndef IMGUI_DISABLE_DEMO_WINDOWS
namespace ImGui { void ShowFontAtlas(ImFontAtlas* atlas); }
#endif

void ImGui::ShowMetricsWindow(bool* p_open)
{
    
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
