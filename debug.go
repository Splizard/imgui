package imgui

import "fmt"

// Debug Tools

// ErrorCheckEndFrameRecover Experimental recovery from incorrect usage of BeginXXX/EndXXX/PushXXX/PopXXX calls.
// Must be called during or before EndFrame().
// This is generally flawed as we are not necessarily End/Popping things in the right order.
// FIXME: Can't recover from inside BeginTabItem/EndTabItem yet.
// FIXME: Can't recover from interleaved BeginTabBar/Begin
func ErrorCheckEndFrameRecover(log_callback ImGuiErrorLogCallback, user_data any) {
	g := GImGui
	for len(g.CurrentWindowStack) > 0 {
		for g.CurrentTable != nil && (g.CurrentTable.OuterWindow == g.CurrentWindow || g.CurrentTable.InnerWindow == g.CurrentWindow) {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing EndTable() in '%s'", g.CurrentTable.OuterWindow.Name)
			}
			EndTable()
		}
		var window = g.CurrentWindow
		IM_ASSERT(window != nil)
		for g.CurrentTabBar != nil {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing EndTabBar() in '%s'", window.Name)
			}
			EndTabBar()
		}
		for window.DC.TreeDepth > 0 {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing TreePop() in '%s'", window.Name)
			}
			TreePop()
		}
		for int(len(g.GroupStack)) > int(window.DC.StackSizesOnBegin.SizeOfGroupStack) {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing EndGroup() in '%s'", window.Name)
			}
			EndGroup()
		}
		for len(window.IDStack) > 1 {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing PopID() in '%s'", window.Name)
			}
			PopID()
		}
		for int(len(g.ColorStack)) > int(window.DC.StackSizesOnBegin.SizeOfColorStack) {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing PopStyleColor() in '%s' for ImGuiCol_%s", window.Name, GetStyleColorName(g.ColorStack[len(g.ColorStack)-1].Col))
			}
			PopStyleColor(1)
		}
		for int(len(g.StyleVarStack)) > int(window.DC.StackSizesOnBegin.SizeOfStyleVarStack) {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing PopStyleVar() in '%s'", window.Name)
			}
			PopStyleVar(1)
		}
		for int(len(g.FocusScopeStack)) > int(window.DC.StackSizesOnBegin.SizeOfFocusScopeStack) {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing PopFocusScope() in '%s'", window.Name)
			}
			PopFocusScope()
		}
		if len(g.CurrentWindowStack) == 1 {
			IM_ASSERT(g.CurrentWindow.IsFallbackWindow)
			break
		}
		IM_ASSERT(window == g.CurrentWindow)
		if window.Flags&ImGuiWindowFlags_ChildWindow != 0 {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing EndChild() for '%s'", window.Name)
			}
			EndChild()
		} else {
			if log_callback != nil {
				log_callback(user_data, "Recovered from missing End() for '%s'", window.Name)
			}
			End()
		}
	}
}

func DebugDrawItemRect(col ImU32 /*= IM_COL32(255,0,0,255)*/) {
	g := GImGui
	var window = g.CurrentWindow
	getForegroundDrawList(window).AddRect(g.LastItemData.Rect.Min, g.LastItemData.Rect.Max, col, 0, 0, 1)
}
func DebugStartItemPicker() { g := GImGui; g.DebugItemPickerActive = true }

// ShowFontAtlas [DEBUG] List fonts in a font atlas and display its texture
func ShowFontAtlas(atlas *ImFontAtlas) {
	for i := range atlas.Fonts {
		var font = atlas.Fonts[i]
		PushInterface(font)
		DebugNodeFont(font)
		PopID()
	}
	if TreeNodeF("Atlas texture", "Atlas texture (%dx%d pixels)", atlas.TexWidth, atlas.TexHeight) {
		var tint_col = ImVec4{1.0, 1.0, 1.0, 1.0}
		var border_col = ImVec4{1.0, 1.0, 1.0, 0.5}
		Image(atlas.TexID, ImVec2{(float)(atlas.TexWidth), (float)(atlas.TexHeight)}, ImVec2{}, ImVec2{1, 1}, tint_col, border_col)
		TreePop()
	}
}

func DebugNodeColumns(columns *ImGuiOldColumns) {
	if !TreeNodeInterface(columns.ID, "Columns Id: 0x%08X, Count: %d, Flags: 0x%04X", columns.ID, columns.Count, columns.Flags) {
		return
	}
	BulletText("Width: %.1f (MinX: %.1f, MaxX: %.1f)", columns.OffMaxX-columns.OffMinX, columns.OffMinX, columns.OffMaxX)
	for column_n := range columns.Columns {
		BulletText("Column %02d: OffsetNorm %.3f (= %.1f px)", column_n, columns.Columns[column_n].OffsetNorm, GetColumnOffsetFromNorm(columns, columns.Columns[column_n].OffsetNorm))
	}
	TreePop()
}

func DebugNodeDrawList(window *ImGuiWindow, draw_list *ImDrawList, label string) {
	g := GImGui
	cfg := &g.DebugMetricsConfig
	var cmd_count = int(len(draw_list.CmdBuffer))
	if cmd_count > 0 && draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1].ElemCount == 0 && draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1].UserCallback == nil {
		cmd_count--
	}
	var node_open = TreeNodeInterface(draw_list, "%s: '%s' %d vtx, %d indices, %d cmds", label, draw_list._OwnerName, len(draw_list.VtxBuffer), len(draw_list.IdxBuffer), cmd_count)
	if draw_list == GetWindowDrawList() {
		SameLine(0, 0)
		TextColored(&ImVec4{1.0, 0.4, 0.4, 1.0}, "CURRENTLY APPENDING") // Can't display stats for active draw list! (we don't have the data double-buffered)
		if node_open {
			TreePop()
		}
		return
	}

	var fg_draw_list = getForegroundDrawList(window) // Render additional visuals into the top-most draw list
	if window != nil && IsItemHovered(0) {
		fg_draw_list.AddRect(window.Pos, window.Pos.Add(window.Size), IM_COL32(255, 255, 0, 255), 0, 0, 1)
	}
	if !node_open {
		return
	}

	if window != nil && !window.WasActive {
		TextDisabled("Warning: owning Window is inactive. This DrawList is not being rendered!")
	}

	for i := range draw_list.CmdBuffer {
		pcmd := &draw_list.CmdBuffer[i]

		if pcmd.UserCallback != nil {
			BulletText("Callback %p, user_data %p", pcmd.UserCallback, pcmd.UserCallbackData)
			continue
		}

		var buf = fmt.Sprintf("DrawCmd:%5d tris, Tex 0x%v, ClipRect (%4.0f,%4.0f)-(%4.0f,%4.0f)",
			pcmd.ElemCount/3, pcmd.TextureId,
			pcmd.ClipRect.x, pcmd.ClipRect.y, pcmd.ClipRect.z, pcmd.ClipRect.w,
		)
		var pcmd_node_open = TreeNodeInterface(pcmd, "%s", buf)
		if IsItemHovered(0) && (cfg.ShowDrawCmdMesh || cfg.ShowDrawCmdBoundingBoxes) && fg_draw_list != nil {
			DebugNodeDrawCmdShowMeshAndBoundingBox(fg_draw_list, draw_list, pcmd, cfg.ShowDrawCmdMesh, cfg.ShowDrawCmdBoundingBoxes)
		}
		if !pcmd_node_open {
			continue
		}

		// Calculate approximate coverage area (touched pixel count)
		// This will be in pixels squared as long there's no post-scaling happening to the renderer output.
		var idx_buffer []ImDrawIdx = nil
		if len(draw_list.IdxBuffer) > 0 {
			idx_buffer = draw_list.IdxBuffer
		}
		var vtx_buffer = draw_list.VtxBuffer[pcmd.VtxOffset:]
		var total_area float = 0.0
		for idx_n := pcmd.IdxOffset; idx_n < pcmd.IdxOffset+pcmd.ElemCount; idx_n++ {
			var triangle [3]ImVec2
			for n := range triangle {
				if idx_buffer != nil {
					triangle[n] = vtx_buffer[idx_buffer[idx_n]].Pos
				} else {
					triangle[n] = vtx_buffer[idx_n].Pos
				}
			}
			total_area += ImTriangleArea(&triangle[0], &triangle[1], &triangle[2])
		}

		// Display vertex information summary. Hover to get all triangles drawn in wire-frame
		Selectable(fmt.Sprintf("Mesh: ElemCount: %d, VtxOffset: +%d, IdxOffset: +%d, Area: ~%0.f px", pcmd.ElemCount, pcmd.VtxOffset, pcmd.IdxOffset, total_area), false, 0, ImVec2{0, 0})
		if IsItemHovered(0) && fg_draw_list != nil {
			DebugNodeDrawCmdShowMeshAndBoundingBox(fg_draw_list, draw_list, pcmd, true, false)
		}

		// Display individual triangles/vertices. Hover on to get the corresponding triangle highlighted.
		var clipper ImGuiListClipper
		clipper.Begin(int(pcmd.ElemCount/3), -1) // Manually coarse clip our print out of individual vertices to save CPU, only items that may be visible.
		for clipper.Step() {
			for prim, idx_i := clipper.DisplayStart, pcmd.IdxOffset+uint(clipper.DisplayStart*3); prim < clipper.DisplayEnd; prim++ {
				var info string
				var triangle [3]ImVec2
				for n := range triangle {
					var v *ImDrawVert
					if idx_buffer != nil {
						v = &vtx_buffer[idx_buffer[idx_i]]
					} else {
						v = &vtx_buffer[idx_i]
					}
					var prefix = "     "
					if n == 0 {
						prefix = "Vert:"
					}
					triangle[n] = v.Pos
					info += fmt.Sprintf("%s %04d: pos (%8.2f,%8.2f), uv (%.6f,%.6f), col %08X\n",
						prefix, idx_i, v.Pos.x, v.Pos.y, v.Uv.x, v.Uv.y, v.Col)
				}

				Selectable(buf, false, 0, ImVec2{})
				if fg_draw_list != nil && IsItemHovered(0) {
					var backup_flags = fg_draw_list.Flags
					fg_draw_list.Flags &= ^ImDrawListFlags_AntiAliasedLines // Disable AA on triangle outlines is more readable for very large and thin triangles.
					fg_draw_list.AddPolyline(triangle[:], 3, IM_COL32(255, 255, 0, 255), ImDrawFlags_Closed, 1.0)
					fg_draw_list.Flags = backup_flags
				}
			}
		}
		TreePop()
	}
	TreePop()
}

func DebugNodeDrawCmdShowMeshAndBoundingBox(out_draw_list *ImDrawList, draw_list *ImDrawList, draw_cmd *ImDrawCmd, show_mesh bool, show_aabb bool) {
	IM_ASSERT(show_mesh || show_aabb)

	// Draw wire-frame version of all triangles
	var clip_rect = ImRectFromVec4(&draw_cmd.ClipRect)
	var vtxs_rect = ImRect{ImVec2{FLT_MAX, FLT_MAX}, ImVec2{-FLT_MAX, -FLT_MAX}}
	var backup_flags = out_draw_list.Flags
	out_draw_list.Flags &= ^ImDrawListFlags_AntiAliasedLines // Disable AA on triangle outlines is more readable for very large and thin triangles.
	for idx_n, idx_end := draw_cmd.IdxOffset, draw_cmd.IdxOffset+draw_cmd.ElemCount; idx_n < idx_end; {
		var idx_buffer = draw_list.IdxBuffer // We don't hold on those pointers past iterations as .AddPolyline() may invalidate them if out_draw_list==draw_list
		var vtx_buffer = draw_list.VtxBuffer

		var triangle [3]ImVec2
		for n := range triangle {
			if idx_buffer != nil {
				triangle[n] = vtx_buffer[idx_buffer[idx_n]].Pos
			} else {
				triangle[n] = vtx_buffer[idx_n].Pos
			}
			vtxs_rect.AddVec(triangle[n])
		}
		if show_mesh {
			out_draw_list.AddPolyline(triangle[:], 3, IM_COL32(255, 255, 0, 255), ImDrawFlags_Closed, 1.0) // In yellow: mesh triangles
		}
	}
	// Draw bounding boxes
	if show_aabb {
		out_draw_list.AddRect(*ImFloorVec(&clip_rect.Min), *ImFloorVec(&clip_rect.Max), IM_COL32(255, 0, 255, 255), 0, 0, 1) // In pink: clipping rectangle submitted to GPU
		out_draw_list.AddRect(*ImFloorVec(&vtxs_rect.Min), *ImFloorVec(&vtxs_rect.Max), IM_COL32(0, 255, 255, 255), 0, 0, 1) // In cyan: bounding box of triangles
	}
	out_draw_list.Flags = backup_flags
}

func DebugNodeFont(font *ImFont) {
	var name string
	if font.ConfigData != nil {
		name = font.ConfigData[0].Name
	}

	var opened = TreeNodeInterface(font, "Font: \"%s\"\n%.2f px, %d glyphs, %d file(s)",
		name, font.FontSize, len(font.Glyphs), font.ConfigDataCount)
	SameLine(0, 0)
	if SmallButton("Set as default") {
		GetIO().FontDefault = font
	}
	if !opened {
		return
	}

	// Display preview text
	PushFont(font)
	Text("The quick brown fox jumps over the lazy dog")
	PopFont()

	// Display details
	SetNextItemWidth(GetFontSize() * 8)
	DragFloat("Font scale", &font.Scale, 0.005, 0.3, 2.0, "%.1f", 0)
	SameLine(0, 0)
	MetricsHelpMarker(
		"Note than the default embedded font is NOT meant to be scaled.\n\n" +
			"Font are currently rendered into bitmaps at a given size at the time of building the atlas. " +
			"You may oversample them to get some flexibility with scaling. " +
			"You can also render at multiple sizes and select which one to use at runtime.\n\n" +
			"(Glimmer of hope: the atlas system will be rewritten in the future to make scaling more flexible.)")
	Text("Ascent: %f, Descent: %f, Height: %f", font.Ascent, font.Descent, font.Ascent-font.Descent)
	Text("Fallback character: '%v' (U+%04X)", string(font.FallbackChar), font.FallbackChar)
	Text("Ellipsis character: '%v' (U+%04X)", string(font.EllipsisChar), font.EllipsisChar)
	var surface_sqrt = (int)(ImSqrt((float)(font.MetricsTotalSurface)))
	Text("Texture Area: about %d px ~%dx%d px", font.MetricsTotalSurface, surface_sqrt, surface_sqrt)
	for config_i := int16(0); config_i < font.ConfigDataCount; config_i++ {
		if font.ConfigData != nil {
			if cfg := &font.ConfigData[config_i]; cfg != nil {
				BulletText(`Input %d: \'%s\', Oversample: (%d,%d), PixelSnapH: %v, Offset: (%.1f,%.1f)`,
					config_i, cfg.Name, cfg.OversampleH, cfg.OversampleV, cfg.PixelSnapH, cfg.GlyphOffset.x, cfg.GlyphOffset.y)
			}
		}
	}

	// Display all glyphs of the fonts in separate pages of 256 characters
	if TreeNodeF("Glyphs", "Glyphs (%d)", len(font.Glyphs)) {
		var draw_list = GetWindowDrawList()
		var glyph_col = GetColorU32FromID(ImGuiCol_Text, 1)
		var cell_size = font.FontSize * 1
		var cell_spacing = GetStyle().ItemSpacing.y
		for base := uint(0); base <= IM_UNICODE_CODEPOINT_MAX; base += 256 {
			// Skip ahead if a large bunch of glyphs are not present in the font (test in chunks of 4k)
			// This is only a small optimization to reduce the number of iterations when IM_UNICODE_MAX_CODEPOINT
			// is large // (if ImWchar==ImWchar32 we will do at least about 272 queries here)
			if (base&4095) == 0 && font.IsGlyphRangeUnused(base, base+4095) {
				base += 4096 - 256
				continue
			}

			var count int = 0
			for n := uint(0); n < 256; n++ {
				if font.FindGlyphNoFallback((ImWchar)(base+n)) != nil {
					count++
				}
			}
			if count <= 0 {
				continue
			}

			var plural = "glyphs"
			if count == 1 {
				plural = "glyph"
			}
			if !TreeNodeInterface(base, "U+%04X..U+%04X (%d %s)", base, base+255, count, plural) {
				continue
			}

			// Draw a 16x16 grid of glyphs
			var base_pos = GetCursorScreenPos()
			for n := 0; n < 256; n++ {
				// We use ImFont::RenderChar as a shortcut because we don't have UTF-8 conversion functions
				// available here and thus cannot easily generate a zero-terminated UTF-8 encoded string.
				var cell_p1 = ImVec2{base_pos.x + float(n%16)*(cell_size+cell_spacing), base_pos.y + float(n/16)*(cell_size+cell_spacing)}
				var cell_p2 = ImVec2{cell_p1.x + cell_size, cell_p1.y + cell_size}
				var glyph = font.FindGlyphNoFallback((ImWchar)(base + uint(n)))

				var c = IM_COL32(255, 255, 255, 50)
				if glyph != nil {
					c = IM_COL32(255, 255, 255, 100)
				}

				draw_list.AddRect(cell_p1, cell_p2, c, 0, 0, 1)
				if glyph != nil {
					font.RenderChar(draw_list, cell_size, cell_p1, glyph_col, (ImWchar)(base+uint(n)))
				}
				if glyph != nil && IsMouseHoveringRect(cell_p1, cell_p2, true) {
					BeginTooltip()
					Text("Codepoint: U+%04X", base+uint(n))
					Separator()
					Text("Visible: %d", glyph.Visible)
					Text("AdvanceX: %.1f", glyph.AdvanceX)
					Text("Pos: (%.2f,%.2f).(%.2f,%.2f)", glyph.X0, glyph.Y0, glyph.X1, glyph.Y1)
					Text("UV: (%.3f,%.3f).(%.3f,%.3f)", glyph.U0, glyph.V0, glyph.U1, glyph.V1)
					EndTooltip()
				}
			}
			Dummy(ImVec2{(cell_size + cell_spacing) * 16, (cell_size + cell_spacing) * 16})
			TreePop()
		}
		TreePop()
	}
	TreePop()
}

func DebugNodeStorage(storage *ImGuiStorage, label string) {
	if !TreeNodeF(label, "%s: %d entries, %d bytes", label, len(storage.Data), len(storage.Data)) {
		return
	}
	for key, val := range storage.Data {
		BulletText("Key 0x%08X Value { i: %d }", key, val) // Important: we currently don't store a type, real value may not be integer.
	}
	for key, val := range storage.Pointers {
		BulletText("Key 0x%08X Value { i: %d }", key, val) // Important: we currently don't store a type, real value may not be integer.
	}
	TreePop()
}

func DebugNodeTabBar(tab_bar *ImGuiTabBar, label string) {
	// Standalone tab bars (not associated to docking/windows functionality) currently hold no discernible strings.
	var p string
	var is_active = (tab_bar.PrevFrameVisible >= GetFrameCount()-2)

	var inactive string
	if !is_active {
		inactive = " *Inactive*"
	}

	p += fmt.Sprintf("%s 0x%08X (%d tabs)%s  { ", label, tab_bar.ID, len(tab_bar.Tabs), inactive)
	for tab_n := int(0); tab_n < ImMinInt(int(len(tab_bar.Tabs)), 3); tab_n++ {
		var tab = &tab_bar.Tabs[tab_n]
		if tab_n > 0 {
			p += ", "
		}
		if tab.NameOffset != -1 {
			p += fmt.Sprint(tab_bar.GetTabName(tab))
		} else {
			p += "???"
		}
	}
	if len(tab_bar.Tabs) > 3 {
		p += " ... }"
	} else {
		p += " } "
	}
	if !is_active {
		PushStyleColorVec(ImGuiCol_Text, GetStyleColorVec4(ImGuiCol_TextDisabled))
	}
	var open = TreeNodeF(label, "%s", p)
	if !is_active {
		PopStyleColor(1)
	}
	if is_active && IsItemHovered(0) {
		var draw_list = GetForegroundDrawList(nil)
		draw_list.AddRect(tab_bar.BarRect.Min, tab_bar.BarRect.Max, IM_COL32(255, 255, 0, 255), 0, 0, 1)
		draw_list.AddLine(&ImVec2{tab_bar.ScrollingRectMinX, tab_bar.BarRect.Min.y}, &ImVec2{tab_bar.ScrollingRectMinX, tab_bar.BarRect.Max.y}, IM_COL32(0, 255, 0, 255), 1)
		draw_list.AddLine(&ImVec2{tab_bar.ScrollingRectMaxX, tab_bar.BarRect.Min.y}, &ImVec2{tab_bar.ScrollingRectMaxX, tab_bar.BarRect.Max.y}, IM_COL32(0, 255, 0, 255), 1)
	}
	if open {
		for tab_n := range tab_bar.Tabs {
			var tab = &tab_bar.Tabs[tab_n]
			PushInterface(tab)
			if SmallButton("<") {
				TabBarQueueReorder(tab_bar, tab, -1)
			}
			SameLine(0, 2)
			if SmallButton(">") {
				TabBarQueueReorder(tab_bar, tab, +1)
			}
			SameLine(0, 0)

			var a, b = " ", "???"
			if tab.ID == tab_bar.SelectedTabId {
				a = "*"
			}
			if tab.NameOffset != -1 {
				b = tab_bar.GetTabName(tab)
			}

			Text("%02d%v Tab 0x%08X '%s' Offset: %.1f, Width: %.1f/%.1f",
				tab_n, a, tab.ID, b, tab.Offset, tab.Width, tab.ContentWidth)
			PopID()
		}
		TreePop()
	}
}

func DebugNodeTable(table *ImGuiTable) {
	buf := ""
	var is_active = (table.LastFrameActive >= GetFrameCount()-2) // Note that fully clipped early out scrolling tables will appear as inactive here.

	var active string
	if !is_active {
		active = " *Inactive*"
	}
	buf = fmt.Sprintf("Table 0x%08X (%d columns, in '%s')%s", table.ID, table.ColumnsCount, table.OuterWindow.Name, active)
	if !is_active {
		PushStyleColorVec(ImGuiCol_Text, GetStyleColorVec4(ImGuiCol_TextDisabled))
	}
	var open = TreeNodeInterface(table, "%s", buf)
	if !is_active {
		PopStyleColor(1)
	}
	if IsItemHovered(0) {
		GetForegroundDrawList(nil).AddRect(table.OuterRect.Min, table.OuterRect.Max, IM_COL32(255, 255, 0, 255), 0, 0, 1)
	}
	if IsItemVisible() && table.HoveredColumnBody != -1 {
		GetForegroundDrawList(nil).AddRect(GetItemRectMin(), GetItemRectMax(), IM_COL32(255, 255, 0, 255), 0, 0, 1)
	}
	if !open {
		return
	}
	var clear_settings = SmallButton("Clear settings")
	BulletText("OuterRect: Pos: (%.1f,%.1f) Size: (%.1f,%.1f) Sizing: '%s'", table.OuterRect.Min.x, table.OuterRect.Min.y, table.OuterRect.GetWidth(), table.OuterRect.GetHeight(), DebugNodeTableGetSizingPolicyDesc(table.Flags))
	var auto string
	if table.InnerWidth == 0.0 {
		auto = " (auto)"
	}
	BulletText("ColumnsGivenWidth: %.1f, ColumnsAutoFitWidth: %.1f, InnerWidth: %.1f%s", table.ColumnsGivenWidth, table.ColumnsAutoFitWidth, table.InnerWidth, auto)
	BulletText("CellPaddingX: %.1f, CellSpacingX: %.1f/%.1f, OuterPaddingX: %.1f", table.CellPaddingX, table.CellSpacingX1, table.CellSpacingX2, table.OuterPaddingX)
	BulletText("HoveredColumnBody: %d, HoveredColumnBorder: %d", table.HoveredColumnBody, table.HoveredColumnBorder)
	BulletText("ResizedColumn: %d, ReorderColumn: %d, HeldHeaderColumn: %d", table.ResizedColumn, table.ReorderColumn, table.HeldHeaderColumn)
	//BulletText("BgDrawChannels: %d/%d", 0, table.BgDrawChannelUnfrozen);
	var sum_weights float = 0.0
	for n := int(0); n < table.ColumnsCount; n++ {
		if table.Columns[n].Flags&ImGuiTableColumnFlags_WidthStretch != 0 {
			sum_weights += table.Columns[n].StretchWeight
		}
	}
	for n := int(0); n < table.ColumnsCount; n++ {
		var column = &table.Columns[n]
		var name = tableGetColumnName(table, n)

		var frozen string
		if n < int(table.FreezeColumnsRequest) {
			frozen = " (Frozen)"
		}
		var stretch float
		if column.StretchWeight > 0.0 {
			stretch = (column.StretchWeight / sum_weights) * 100
		}

		var sorting string
		if column.SortDirection == ImGuiSortDirection_Ascending {
			sorting = " (Asc)"
		} else {
			if column.SortDirection == ImGuiSortDirection_Descending {
				sorting = " (Desc)"
			}
		}

		var flags string
		if column.Flags&ImGuiTableColumnFlags_WidthStretch != 0 {
			flags += "WidthStretch "
		}
		if column.Flags&ImGuiTableColumnFlags_WidthFixed != 0 {
			flags += "WidthFixed "
		}
		if column.Flags&ImGuiTableColumnFlags_NoResize != 0 {
			flags += "NoResize "
		}

		buf = fmt.Sprintf(
			"Column %v order %v '%s': offset %+.2f to %+.2f%s\n"+
				"Enabled: %v, VisibleX/Y: %v/%v, RequestOutput: %v, SkipItems: %v, DrawChannels: %v,%v\n"+
				"WidthGiven: %.1f, Request/Auto: %.1f/%.1f, StretchWeight: %.3f (%.1f%%)\n"+
				"MinX: %.1f, MaxX: %.1f (%+.1f), ClipRect: %.1f to %.1f (+%.1f)\n"+
				"ContentWidth: %.1f,%.1f, HeadersUsed/Ideal %.1f/%.1f\n"+
				"Sort: %d%s, UserID: 0x%08X, Flags: 0x%04X: %s..",
			n, column.DisplayOrder, name, column.MinX-table.WorkRect.Min.x, column.MaxX-table.WorkRect.Min.x, frozen,
			column.IsEnabled, column.IsVisibleX, column.IsVisibleY, column.IsRequestOutput, column.IsSkipItems, column.DrawChannelFrozen, column.DrawChannelUnfrozen,
			column.WidthGiven, column.WidthRequest, column.WidthAuto, column.StretchWeight, stretch,
			column.MinX, column.MaxX, column.MaxX-column.MinX, column.ClipRect.Min.x, column.ClipRect.Max.x, column.ClipRect.Max.x-column.ClipRect.Min.x,
			column.ContentMaxXFrozen-column.WorkMinX, column.ContentMaxXUnfrozen-column.WorkMinX, column.ContentMaxXHeadersUsed-column.WorkMinX, column.ContentMaxXHeadersIdeal-column.WorkMinX,
			column.SortOrder, sorting, column.UserID, column.Flags,
			flags,
		)
		Bullet()
		Selectable(buf, false, 0, ImVec2{})
		if IsItemHovered(0) {
			var r = ImRect{ImVec2{column.MinX, table.OuterRect.Min.y}, ImVec2{column.MaxX, table.OuterRect.Max.y}}
			GetForegroundDrawList(nil).AddRect(r.Min, r.Max, IM_COL32(255, 255, 0, 255), 0, 0, 1)
		}
	}
	if settings := TableGetBoundSettings(table); settings != nil {
		DebugNodeTableSettings(settings)
	}
	if clear_settings {
		table.IsResetAllRequest = true
	}
	TreePop()
}

func DebugNodeTableSettings(settings *ImGuiTableSettings) {
	if !TreeNodeInterface(&settings.ID, "Settings 0x%08X (%d columns)", settings.ID, settings.ColumnsCount) {
		return
	}
	BulletText("SaveFlags: 0x%08X", settings.SaveFlags)
	BulletText("ColumnsCount: %d (max %d)", settings.ColumnsCount, settings.ColumnsCountMax)
	for n := int(0); n < int(settings.ColumnsCount); n++ {
		var column_settings = &settings.Columns[n]
		var sort_dir = ImGuiSortDirection_None
		if column_settings.SortOrder != -1 {
			sort_dir = (ImGuiSortDirection)(column_settings.SortDirection)
		}

		var sorting string
		if sort_dir == ImGuiSortDirection_Ascending {
			sorting = "Asc"
		} else if sort_dir == ImGuiSortDirection_Descending {
			sorting = "Des"
		} else {
			sorting = "---"
		}

		var stretch = "Width "
		if column_settings.IsStretch != 0 {
			stretch = "Weight"
		}

		BulletText("Column %d Order %d SortOrder %d %s Vis %d %s %7.3f UserID 0x%08X",
			n, column_settings.DisplayOrder, column_settings.SortOrder, sorting,
			column_settings.IsEnabled, stretch, column_settings.WidthOrWeight, column_settings.UserID)
	}
	TreePop()
}

func DebugNodeWindow(window *ImGuiWindow, label string) {
	if window == nil {
		BulletText("%s: nil", label)
		return
	}
	g := GImGui

	var selected_flags = ImGuiTreeNodeFlags_None
	if window == g.NavWindow {
		selected_flags = ImGuiTreeNodeFlags_Selected
	}

	var is_active = window.WasActive
	var tree_node_flags = selected_flags
	if !is_active {
		PushStyleColorVec(ImGuiCol_Text, GetStyleColorVec4(ImGuiCol_TextDisabled))
	}

	var active_string = " *Inactive*"
	if is_active {
		active_string = ""
	}

	var open = TreeNodeEx(label, tree_node_flags, "%s '%s'%s", label, window.Name, active_string)
	if !is_active {
		PopStyleColor(1)
	}
	if IsItemHovered(0) && is_active {
		getForegroundDrawList(window).AddRect(window.Pos, window.Pos.Add(window.Size), IM_COL32(255, 255, 0, 255), 0, 0, 1)
	}
	if !open {
		return
	}

	if window.MemoryCompacted {
		TextDisabled("Note: some memory buffers have been compacted/freed.")
	}

	var flags = window.Flags

	var kind string
	if (flags & ImGuiWindowFlags_ChildWindow) != 0 {
		kind += " Child "
	}
	if (flags & ImGuiWindowFlags_Tooltip) != 0 {
		kind = " Tooltip "
	}
	if (flags & ImGuiWindowFlags_Popup) != 0 {
		kind += " Popup "
	}
	if (flags & ImGuiWindowFlags_Modal) != 0 {
		kind += " Modal "
	}
	if (flags & ImGuiWindowFlags_ChildMenu) != 0 {
		kind += " ChildMenu "
	}
	if (flags & ImGuiWindowFlags_NoSavedSettings) != 0 {
		kind += " NoSavedSettings "
	}
	if (flags & ImGuiWindowFlags_NoMouseInputs) != 0 {
		kind += " NoMouseInputs "
	}
	if (flags & ImGuiWindowFlags_NoNavInputs) != 0 {
		kind += " NoNavInputs "
	}
	if (flags & ImGuiWindowFlags_AlwaysAutoResize) != 0 {
		kind += " AlwaysAutoResize "
	}

	DebugNodeDrawList(window, window.DrawList, "DrawList")
	BulletText("Pos: (%.1f,%.1f), Size: (%.1f,%.1f), ContentSize (%.1f,%.1f) Ideal (%.1f,%.1f)", window.Pos.x, window.Pos.y, window.Size.x, window.Size.y, window.ContentSize.x, window.ContentSize.y, window.ContentSizeIdeal.x, window.ContentSizeIdeal.y)
	BulletText("Flags: 0x%08X (%s..)", flags, kind)

	var scrollX, scrollY string
	if window.ScrollbarX {
		scrollX = "X"
	}
	if window.ScrollbarY {
		scrollY = "Y"
	}

	var BeginOrderWithinContext = window.BeginOrderWithinContext
	if window.Active || window.WasActive {
		BeginOrderWithinContext = -1
	}

	BulletText("Scroll: (%.2f/%.2f,%.2f/%.2f) Scrollbar:%s%s", window.Scroll.x, window.ScrollMax.x, window.Scroll.y, window.ScrollMax.y, scrollX, scrollY)
	BulletText("Active: %v/%v, WriteAccessed: %v, BeginOrderWithinContext: %d", window.Active, window.WasActive, window.WriteAccessed, BeginOrderWithinContext)
	BulletText("Appearing: %v, Hidden: %v (CanSkip %v Cannot %v), SkipItems: %v", window.Appearing, window.Hidden, window.HiddenFramesCanSkipItems, window.HiddenFramesCannotSkipItems, window.SkipItems)
	for layer := 0; layer < ImGuiNavLayer_COUNT; layer++ {
		var r = window.NavRectRel[layer]
		if r.Min.x >= r.Max.y && r.Min.y >= r.Max.y {
			BulletText("NavLastIds[%d]: 0x%08X", layer, window.NavLastIds[layer])
			continue
		}
		BulletText("NavLastIds[%d]: 0x%08X at +(%.1f,%.1f)(%.1f,%.1f)", layer, window.NavLastIds[layer], r.Min.x, r.Min.y, r.Max.x, r.Max.y)
		if IsItemHovered(0) {
			getForegroundDrawList(window).AddRect(r.Min.Add(window.Pos), r.Max.Add(window.Pos), IM_COL32(255, 255, 0, 255), 0, 0, 1)
		}
	}
	var lastChildNavWindowName = "nil"
	if window.NavLastChildNavWindow != nil {
		lastChildNavWindowName = window.NavLastChildNavWindow.Name
	}

	BulletText("NavLayersActiveMask: %X, NavLastChildNavWindow: %s", window.DC.NavLayersActiveMask, lastChildNavWindowName)
	if window.RootWindow != window {
		DebugNodeWindow(window.RootWindow, "RootWindow")
	}
	if window.ParentWindow != nil {
		DebugNodeWindow(window.ParentWindow, "ParentWindow")
	}
	if len(window.DC.ChildWindows) > 0 {
		DebugNodeWindowsList(window.DC.ChildWindows, "ChildWindows")
	}
	if len(window.ColumnsStorage) > 0 && TreeNodeEx("Columns", 0, "Columns sets (%d)", len(window.ColumnsStorage)) {
		for n := range window.ColumnsStorage {
			DebugNodeColumns(&window.ColumnsStorage[n])
		}
		TreePop()
	}
	DebugNodeStorage(&window.StateStorage, "Storage")
	TreePop()
}

func DebugNodeWindowSettings(settings *ImGuiWindowSettings) {
	Text("0x%08X \"%s\" Pos (%d,%d) Size (%d,%d) Collapsed=%v",
		settings.ID, settings.GetName(), settings.Pos.x, settings.Pos.y, settings.Size.x, settings.Size.y, settings.Collapsed)
}

func DebugNodeWindowsList(windows []*ImGuiWindow, label string) {
	if !TreeNodeEx(label, 0, "%s (%d)", label, len(windows)) {
		return
	}
	Text("(In front-to-back order:)")
	for i := len(windows) - 1; i >= 0; i-- { // Iterate front to back
		PushID(int32(windows[i].ID))
		DebugNodeWindow(windows[i], "Window")
		PopID()
	}
	TreePop()
}

func DebugNodeViewport(viewport *ImGuiViewportP) {
	SetNextItemOpen(true, ImGuiCond_Once)
	if TreeNodeF("viewport0", "Viewport #%d", 0) {
		var flags = viewport.Flags
		BulletText("Main Pos: (%.0f,%.0f), Size: (%.0f,%.0f)\nWorkArea Offset Left: %.0f Top: %.0f, Right: %.0f, Bottom: %.0f",
			viewport.Pos.x, viewport.Pos.y, viewport.Size.x, viewport.Size.y,
			viewport.WorkOffsetMin.x, viewport.WorkOffsetMin.y, viewport.WorkOffsetMax.x, viewport.WorkOffsetMax.y)

		var flagString string
		if (flags & ImGuiViewportFlags_IsPlatformWindow) != 0 {
			flagString += " IsPlatformWindow "
		}
		if (flags & ImGuiViewportFlags_IsPlatformMonitor) != 0 {
			flagString += " IsPlatformMonitor "
		}
		if (flags & ImGuiViewportFlags_OwnedByApp) != 0 {
			flagString += " OwnedByApp "
		}

		BulletText("Flags: 0x%04X =%s", viewport.Flags, flagString)
		for layer_i := range viewport.DrawDataBuilder {
			for draw_list_i := range viewport.DrawDataBuilder[layer_i] {
				DebugNodeDrawList(nil, viewport.DrawDataBuilder[layer_i][draw_list_i], "DrawList")
			}
		}

		TreePop()
	}
}

func DebugRenderViewportThumbnail(draw_list *ImDrawList, viewport *ImGuiViewportP, bb *ImRect) {
	g := GImGui
	var window = g.CurrentWindow

	var scale = bb.GetSize().Div(viewport.Size)
	var off = bb.Min.Sub(viewport.Pos.Mul(scale))
	var alpha_mul float = 1.0
	window.DrawList.AddRectFilled(bb.Min, bb.Max, GetColorU32FromID(ImGuiCol_Border, alpha_mul*0.40), 0, 0)
	for i := range g.Windows {
		var thumb_window = g.Windows[i]
		if !thumb_window.WasActive || (thumb_window.Flags&ImGuiWindowFlags_ChildWindow != 0) {
			continue
		}

		var thumb_r = thumb_window.Rect()
		var title_r = thumb_window.TitleBarRect()

		a, b := off.Add(thumb_r.Min.Mul(scale)), off.Add(thumb_r.Max.Mul(scale))
		c, d := off.Add(title_r.Min.Mul(scale)), off.Add(ImVec2{title_r.Max.x, title_r.Min.y}.Mul(scale)).Add(ImVec2{0, 5}) // Exaggerate title bar height

		thumb_r = ImRect{*ImFloorVec(&a), *ImFloorVec(&b)}
		title_r = ImRect{*ImFloorVec(&c), *ImFloorVec(&d)}
		thumb_r.ClipWithFull(*bb)
		title_r.ClipWithFull(*bb)
		var window_is_focused = (g.NavWindow != nil && thumb_window.RootWindowForTitleBarHighlight == g.NavWindow.RootWindowForTitleBarHighlight)

		focused_color := ImGuiCol_TitleBg
		if window_is_focused {
			focused_color = ImGuiCol_TitleBgActive
		}

		window.DrawList.AddRectFilled(thumb_r.Min, thumb_r.Max, GetColorU32FromID(ImGuiCol_WindowBg, alpha_mul), 0, 0)
		window.DrawList.AddRectFilled(title_r.Min, title_r.Max, GetColorU32FromID(focused_color, alpha_mul), 0, 0)
		window.DrawList.AddRect(thumb_r.Min, thumb_r.Max, GetColorU32FromID(ImGuiCol_Border, alpha_mul), 0, 0, 1)
		window.DrawList.AddTextV(g.Font, g.FontSize*1.0, title_r.Min, GetColorU32FromID(ImGuiCol_Text, alpha_mul), thumb_window.Name, 0, nil)
	}
	draw_list.AddRect(bb.Min, bb.Max, GetColorU32FromID(ImGuiCol_Border, alpha_mul), 0, 0, 1)
}

// MetricsHelpMarker Avoid naming collision with imgui_demo.cpp's HelpMarker() for unity builds.
func MetricsHelpMarker(desc string) {
	TextDisabled("(?)")
	if IsItemHovered(0) {
		BeginTooltip()
		PushTextWrapPos(GetFontSize() * 35.0)
		TextUnformatted(desc)
		PopTextWrapPos()
		EndTooltip()
	}
}

func RenderViewportsThumbnails() {
	g := GImGui
	var window = g.CurrentWindow

	// We don't display full monitor bounds (we could, but it often looks awkward), instead we display just enough to cover all of our viewports.
	const SCALE = 1.0 / 8.0
	var bb_full = ImRect{ImVec2{FLT_MAX, FLT_MAX}, ImVec2{-FLT_MAX, -FLT_MAX}}
	for n := range g.Viewports {
		bb_full.AddRect(g.Viewports[n].GetMainRect())
	}
	var p = window.DC.CursorPos
	var off = p.Sub(bb_full.Min.Scale(SCALE))
	for n := range g.Viewports {
		var viewport = g.Viewports[n]
		var viewport_draw_bb = ImRect{off.Add((viewport.Pos).Scale(SCALE)), off.Add((viewport.Pos.Add(viewport.Size)).Scale(SCALE))}
		DebugRenderViewportThumbnail(window.DrawList, viewport, &viewport_draw_bb)
	}
	Dummy(bb_full.GetSize().Scale(SCALE))
}

// Debugging enums
const (
	WRT_OuterRect = iota
	WRT_InnerRect
	WRT_OuterRectClipped
	WRT_InnerClipRect
	WRT_WorkRect
	WRT_Content
	WRT_ContentIdeal
	WRT_ContentRegionRect
	WRT_Count
)

var wrt_rects_names = []string{
	"OuterRect", "OuterRectClipped", "InnerRect", "InnerClipRect", "WorkRect", "Content", "ContentIdeal", "ContentRegionRect",
}

const (
	TRT_OuterRect = iota
	TRT_InnerRect
	TRT_WorkRect
	TRT_HostClipRect
	TRT_InnerClipRect
	TRT_BackgroundClipRect
	TRT_ColumnsRect
	TRT_ColumnsWorkRect
	TRT_ColumnsClipRect
	TRT_ColumnsContentHeadersUsed
	TRT_ColumnsContentHeadersIdeal
	TRT_ColumnsContentFrozen
	TRT_ColumnsContentUnfrozen
	TRT_Count
)

var trt_rects_names = []string{
	"OuterRect", "InnerRect", "WorkRect", "HostClipRect", "InnerClipRect", "BackgroundClipRect", "ColumnsRect", "ColumnsWorkRect", "ColumnsClipRect", "ColumnsContentHeadersUsed", "ColumnsContentHeadersIdeal", "ColumnsContentFrozen", "ColumnsContentUnfrozen",
}

// UpdateDebugToolItemPicker [DEBUG] Item picker tool - start with DebugStartItemPicker() - useful to visually select an item and break into its call-stack.
func UpdateDebugToolItemPicker() {
	g := GImGui
	g.DebugItemPickerBreakId = 0
	if g.DebugItemPickerActive {
		var hovered_id = g.HoveredIdPreviousFrame
		SetMouseCursor(ImGuiMouseCursor_Hand)
		if IsKeyPressedMap(ImGuiKey_Escape, true) {
			g.DebugItemPickerActive = false
		}
		if IsMouseClicked(0, true) && hovered_id != 0 {
			g.DebugItemPickerBreakId = hovered_id
			g.DebugItemPickerActive = false
		}
		SetNextWindowBgAlpha(0.60)
		BeginTooltip()
		Text("HoveredId: 0x%08X", hovered_id)
		Text("Press ESC to abort picking.")
		var c *ImVec4
		if hovered_id != 0 {
			c = GetStyleColorVec4(ImGuiCol_Text)
		} else {
			c = GetStyleColorVec4(ImGuiCol_TextDisabled)
		}
		TextColored(c, "Click to break in debugger!")
		EndTooltip()
	}
}

// ShowMetricsWindow create Metrics/Debugger window. display Dear ImGui internals: windows, draw commands, various internal state, etc.
func ShowMetricsWindow(p_open *bool) {
	if !Begin("Dear ImGui Metrics/Debugger", p_open, 0) {
		End()
		return
	}

	g := GImGui
	io := g.IO
	cfg := &g.DebugMetricsConfig

	Text("Dear ImGui %s", GetVersion())
	Text("Application average %.3f ms/frame (%.1f FPS)", 1000.0/io.Framerate, io.Framerate)
	Text("%d vertices, %d indices (%d triangles)", io.MetricsRenderVertices, io.MetricsRenderIndices, io.MetricsRenderIndices/3)
	Text("%d active windows (%d visible)", io.MetricsActiveWindows, io.MetricsRenderWindows)
	Text("%d active allocations", io.MetricsActiveAllocations)

	Separator()

	if cfg.ShowWindowsRectsType < 0 {
		cfg.ShowWindowsRectsType = WRT_WorkRect
	}
	if cfg.ShowTablesRectsType < 0 {
		cfg.ShowTablesRectsType = TRT_WorkRect
	}

	GetTableRect := func(table *ImGuiTable, rect_type, n int) ImRect {
		switch rect_type {
		case TRT_OuterRect:
			return table.OuterRect
		case TRT_InnerRect:
			return table.InnerRect
		case TRT_WorkRect:
			return table.WorkRect
		case TRT_HostClipRect:
			return table.HostClipRect
		case TRT_InnerClipRect:
			return table.InnerClipRect
		case TRT_BackgroundClipRect:
			return table.BgClipRect
			/*case TRT_ColumnsRect:
				var c = &table.Columns[n]
				return ImRect(c.MinX, table.InnerClipRect.Min.y, c.MaxX, table.InnerClipRect.Min.y+table.LastOuterHeight)
			case TRT_ColumnsWorkRect:
				var c = &table.Columns[n]
				return ImRect(c.WorkMinX, table.WorkRect.Min.y, c.WorkMaxX, table.WorkRect.Max.y)
			case TRT_ColumnsClipRect:
				var c = &table.Columns[n]
				return c.ClipRect
			case TRT_ColumnsContentHeadersUsed:
				var c = &table.Columns[n] // Note: y1/y2 not always accurate
				return ImRect(c.WorkMinX, table.InnerClipRect.Min.y, c.ContentMaxXHeadersUsed, table.InnerClipRect.Min.y+table.LastFirstRowHeight)
			case TRT_ColumnsContentHeadersIdeal:
				var c = &table.Columns[n]
				return ImRect(c.WorkMinX, table.InnerClipRect.Min.y, c.ContentMaxXHeadersIdeal, table.InnerClipRect.Min.y+table.LastFirstRowHeight)
			case TRT_ColumnsContentFrozen:
				var c = &table.Columns[n]
				return ImRect(c.WorkMinX, table.InnerClipRect.Min.y, c.ContentMaxXFrozen, table.InnerClipRect.Min.y+table.LastFirstRowHeight)
			case TRT_ColumnsContentUnfrozen:
				var c = &table.Columns[n]
				return ImRect(c.WorkMinX, table.InnerClipRect.Min.y+table.LastFirstRowHeight, c.ContentMaxXUnfrozen, table.InnerClipRect.Max.y)*/
		}
		IM_ASSERT(false)
		return ImRect{}
	}

	GetWindowRect := func(window *ImGuiWindow, rect_type int) ImRect {
		switch rect_type {
		case WRT_OuterRect:
			return window.Rect()
		case WRT_OuterRectClipped:
			return window.OuterRectClipped
		case WRT_InnerRect:
			return window.InnerRect
		case WRT_InnerClipRect:
			return window.InnerClipRect
		case WRT_WorkRect:
			return window.WorkRect
		case WRT_Content:
			var min = window.InnerRect.Min.Sub(window.Scroll).Add(window.WindowPadding)
			return ImRect{min, min.Add(window.ContentSize)}
		case WRT_ContentIdeal:
			var min = window.InnerRect.Min.Sub(window.Scroll).Add(window.WindowPadding)
			return ImRect{min, min.Add(window.ContentSizeIdeal)}
		case WRT_ContentRegionRect:
			return window.ContentRegionRect
		}
		IM_ASSERT(false)
		return ImRect{}
	}

	//Tools
	if TreeNode("Tools") {
		if Button("Item Picker") {
			DebugStartItemPicker()
		}
		SameLine(0, 0)
		MetricsHelpMarker("Will call the IM_DEBUG_BREAK() macro to break in debugger.\nWarning: If you don't have a debugger attached, this will probably crash.")

		Checkbox("Show windows begin order", &cfg.ShowWindowsBeginOrder)
		Checkbox("Show windows rectangles", &cfg.ShowWindowsRects)
		SameLine(0, 0)

		SetNextItemWidth(GetFontSize() * 12)
		cfg.ShowWindowsRects = Combo("##show_windows_rect_type", &cfg.ShowWindowsRectsType, wrt_rects_names, WRT_Count, WRT_Count) || cfg.ShowWindowsRects
		if cfg.ShowWindowsRects && g.NavWindow != nil {
			BulletText("'%s':", g.NavWindow.Name)
			Indent(0)
			for rect_n := int(0); rect_n < WRT_Count; rect_n++ {
				var r = GetWindowRect(g.NavWindow, rect_n)
				Text("(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), wrt_rects_names[rect_n])
			}
			Unindent(0)
		}
		Checkbox("Show ImDrawCmd mesh when hovering", &cfg.ShowDrawCmdMesh)
		Checkbox("Show ImDrawCmd bounding boxes when hovering", &cfg.ShowDrawCmdBoundingBoxes)

		Checkbox("Show tables rectangles", &cfg.ShowTablesRects)
		SameLine(0, 0)
		SetNextItemWidth(GetFontSize() * 12)
		cfg.ShowTablesRects = Combo("##show_table_rects_type", &cfg.ShowTablesRectsType, trt_rects_names, TRT_Count, TRT_Count) || cfg.ShowTablesRects

		if cfg.ShowTablesRects && g.NavWindow != nil {
			for _, table := range g.Tables {
				if table == nil || table.LastFrameActive < g.FrameCount-1 || (table.OuterWindow != g.NavWindow && table.InnerWindow != g.NavWindow) {
					continue
				}

				BulletText("Table 0x%08X (%d columns, in '%s')", table.ID, table.ColumnsCount, table.OuterWindow.Name)
				if IsItemHovered(0) {
					GetForegroundDrawList(nil).AddRect(table.OuterRect.Min.Sub(ImVec2{1, 1}), table.OuterRect.Max.Add(ImVec2{1, 1}), IM_COL32(255, 255, 0, 255), 0.0, 0, 2.0)
				}
				Indent(0)
				buf := ""
				for rect_n := int(0); rect_n < TRT_Count; rect_n++ {
					if rect_n >= TRT_ColumnsRect {
						if rect_n != TRT_ColumnsRect && rect_n != TRT_ColumnsClipRect {
							continue
						}
						for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
							var r = GetTableRect(table, rect_n, column_n)
							buf = fmt.Sprintf("(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) Col %d %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), column_n, trt_rects_names[rect_n])
							Selectable(buf, false, 0, ImVec2{})
							if IsItemHovered(0) {
								GetForegroundDrawList(nil).AddRect(r.Min.Sub(ImVec2{1, 1}), r.Max.Add(ImVec2{1, 1}), IM_COL32(255, 255, 0, 255), 0.0, 0, 2.0)
							}
						}
					} else {
						var r = GetTableRect(table, rect_n, -1)
						buf = fmt.Sprintf("(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), trt_rects_names[rect_n])
						Selectable(buf, false, 0, ImVec2{})
						if IsItemHovered(0) {
							GetForegroundDrawList(nil).AddRect(r.Min.Sub(ImVec2{1, 1}), r.Max.Add(ImVec2{1, 1}), IM_COL32(255, 255, 0, 255), 0.0, 0, 2.0)
						}
					}
				}
				Unindent(0)
			}
		}
		TreePop()
	}

	DebugNodeWindowsList(g.Windows, "Windows")

	// DrawLists
	var drawlist_count int = 0
	for viewport_i := int(0); viewport_i < int(len(g.Viewports)); viewport_i++ {
		drawlist_count += g.Viewports[viewport_i].DrawDataBuilder.GetDrawListCount()
	}
	if TreeNodeF("DrawLists", "DrawLists (%d)", drawlist_count) {
		for viewport_i := 0; viewport_i < len(g.Viewports); viewport_i++ {
			var viewport = g.Viewports[viewport_i]
			for layer_i := 0; layer_i < len(viewport.DrawDataBuilder); layer_i++ {
				for draw_list_i := 0; draw_list_i < len(viewport.DrawDataBuilder[layer_i]); draw_list_i++ {
					DebugNodeDrawList(nil, viewport.DrawDataBuilder[layer_i][draw_list_i], "DrawList")
				}
			}
		}
		TreePop()
	}

	// Viewports
	if TreeNodeF("Viewports", "Viewports (%d)", len(g.Viewports)) {
		Indent(GetTreeNodeToLabelSpacing())
		RenderViewportsThumbnails()
		//FIXME causes infinite loop.
		/*Unindent(GetTreeNodeToLabelSpacing())
		for i := 0; i < len(g.Viewports); i++ {
			DebugNodeViewport(g.Viewports[i])
		}*/
		TreePop()
	}

	// Details for Popups
	if TreeNodeF("Popups", "Popups (%d)", len(g.OpenPopupStack)) {
		for i := range g.OpenPopupStack {
			var window = g.OpenPopupStack[i].Window
			var winName = "nil"
			if window.Name != "" {
				winName = window.Name
			}
			var kind string
			if (window.Flags & ImGuiWindowFlags_ChildWindow) != 0 {
				kind += " ChildWindow"
			}
			if (window.Flags & ImGuiWindowFlags_ChildMenu) != 0 {
				kind += " ChildMenu"
			}
			BulletText("PopupID: %08x, Window: '%s'%s", g.OpenPopupStack[i].PopupId, winName, kind)
		}
		TreePop()
	}

	// Details for TabBars
	if TreeNodeF("TabBars", "Tab Bars (%d)", len(g.TabBars)) {
		for _, tab_bar := range g.TabBars {
			PushInterface(tab_bar)
			DebugNodeTabBar(tab_bar, "TabBar")
			PopID()
		}
		TreePop()
	}

	// Details for Tables
	if TreeNodeF("Tables", "Tables (%d)", len(g.Tables)) {
		for _, table := range g.Tables {
			PushInterface(table)
			DebugNodeTable(table)
			PopID()
		}
		TreePop()
	}

	// Details for Fonts
	var atlas = g.IO.Fonts
	if TreeNodeF("Fonts", "Fonts (%d)", len(atlas.Fonts)) {
		ShowFontAtlas(atlas)
		TreePop()
	}

	// Settings
	if TreeNode("Settings") {
		if SmallButton("Clear") {
			ClearIniSettings()
		}
		SameLine(0, 0)
		if SmallButton("Save to memory") {
			SaveIniSettingsToMemory(nil)
		}
		SameLine(0, 0)
		if SmallButton("Save to disk") {
			SaveIniSettingsToDisk(g.IO.IniFilename)
		}
		SameLine(0, 0)
		if g.IO.IniFilename != "" {
			Text("\"%s\"", g.IO.IniFilename)
		} else {
			TextUnformatted("<nil>")
		}
		Text("SettingsDirtyTimer %.2f", g.SettingsDirtyTimer)
		if TreeNodeF("SettingsHandlers", "Settings handlers: (%d)", len(g.SettingsHandlers)) {
			for n := range g.SettingsHandlers {
				BulletText("%s", g.SettingsHandlers[n].TypeName)
			}
			TreePop()
		}
		if TreeNodeF("SettingsWindows", "Settings packed data: Windows: %d bytes", len(g.SettingsWindows)) {
			for i := range g.SettingsWindows {
				DebugNodeWindowSettings(&g.SettingsWindows[i])
			}
			TreePop()
		}

		if TreeNodeF("SettingsTables", "Settings packed data: Tables: %d bytes", len(g.SettingsTables)) {
			for i := range g.SettingsTables {
				DebugNodeTableSettings(&g.SettingsTables[i])
			}
			TreePop()
		}
		TreePop()
	}

	// Misc Details
	if TreeNode("Internal state") {
		var input_source_names = []string{"None", "Mouse", "Keyboard", "Gamepad", "Nav", "Clipboard"}
		IM_ASSERT(int(len(input_source_names)) == int(ImGuiInputSource_COUNT))

		Text("WINDOWING")
		Indent(0)

		var name, rootName, underName, movingName = "nil", "nil", "nil", "nil"
		if g.HoveredWindow != nil {
			name, rootName, underName, movingName = g.HoveredWindow.Name, g.HoveredWindow.RootWindow.Name,
				g.HoveredWindow.RootWindow.Name, g.MovingWindow.Name
		}

		Text("HoveredWindow: '%s'", name)
		Text("HoveredWindow.Root: '%s'", rootName)
		Text("HoveredWindowUnderMovingWindow: '%s'", underName)
		Text("MovingWindow: '%s'", movingName)
		Unindent(0)

		var activeName = "nil"
		if g.ActiveIdWindow != nil {
			activeName = g.ActiveIdWindow.Name
		}

		Text("ITEMS")
		Indent(0)
		Text("ActiveId: 0x%08X/0x%08X (%.2f sec), AllowOverlap: %d, Source: %s", g.ActiveId, g.ActiveIdPreviousFrame, g.ActiveIdTimer, g.ActiveIdAllowOverlap, input_source_names[g.ActiveIdSource])
		Text("ActiveIdWindow: '%s'", activeName)
		Text("ActiveIdUsing: Wheel: %d, NavDirMask: %X, NavInputMask: %X, KeyInputMask: %llX", g.ActiveIdUsingMouseWheel, g.ActiveIdUsingNavDirMask, g.ActiveIdUsingNavInputMask, g.ActiveIdUsingKeyInputMask)
		Text("HoveredId: 0x%08X (%.2f sec), AllowOverlap: %d", g.HoveredIdPreviousFrame, g.HoveredIdTimer, g.HoveredIdAllowOverlap) // Not displaying g.HoveredId as it is update mid-frame
		Text("DragDrop: %d, SourceId = 0x%08X, Payload \"%s\" (%d bytes)", g.DragDropActive, g.DragDropPayload.SourceId, g.DragDropPayload.DataType, g.DragDropPayload.DataSize)
		Unindent(0)

		var navWindowName, navTargetName = "nil", "nil"
		if g.NavWindow != nil {
			navWindowName, navTargetName = g.NavWindow.Name, g.NavWindowingTarget.Name
		}

		Text("NAV,FOCUS")
		Indent(0)
		Text("NavWindow: '%s'", navWindowName)
		Text("NavId: 0x%08X, NavLayer: %d", g.NavId, g.NavLayer)
		Text("NavInputSource: %s", input_source_names[g.NavInputSource])
		Text("NavActive: %d, NavVisible: %d", g.IO.NavActive, g.IO.NavVisible)
		Text("NavActivateId: 0x%08X, NavInputId: 0x%08X", g.NavActivateId, g.NavInputId)
		Text("NavDisableHighlight: %d, NavDisableMouseHover: %d", g.NavDisableHighlight, g.NavDisableMouseHover)
		Text("NavFocusScopeId = 0x%08X", g.NavFocusScopeId)
		Text("NavWindowingTarget: '%s'", navTargetName)
		Unindent(0)

		TreePop()
	}

	// Overlay: Display windows Rectangles and Begin Order
	if cfg.ShowWindowsRects || cfg.ShowWindowsBeginOrder {
		for n := range g.Windows {
			var window = g.Windows[n]
			if !window.WasActive {
				continue
			}
			var draw_list = getForegroundDrawList(window)
			if cfg.ShowWindowsRects {
				var r = GetWindowRect(window, cfg.ShowWindowsRectsType)
				draw_list.AddRect(r.Min, r.Max, IM_COL32(255, 0, 128, 255), 0, 0, 1)
			}
			if cfg.ShowWindowsBeginOrder && (window.Flags&ImGuiWindowFlags_ChildWindow == 0) {
				var font_size = GetFontSize()
				draw_list.AddRectFilled(window.Pos, window.Pos.Add(ImVec2{font_size, font_size}), IM_COL32(200, 100, 100, 255), 0, 0)
				draw_list.AddText(window.Pos, IM_COL32(255, 255, 255, 255), fmt.Sprint(window.BeginOrderWithinContext))
			}
		}
	}

	// Overlay: Display Tables Rectangles
	if cfg.ShowTablesRects {
		for _, table := range g.Tables {
			if table == nil || table.LastFrameActive < g.FrameCount-1 {
				continue
			}
			var draw_list = getForegroundDrawList(table.OuterWindow)
			if cfg.ShowTablesRectsType >= TRT_ColumnsRect {
				for column_n := int(0); column_n < table.ColumnsCount; column_n++ {
					var r = GetTableRect(table, cfg.ShowTablesRectsType, column_n)
					var col = IM_COL32(255, 0, 128, 255)
					if int(table.HoveredColumnBody) == column_n {
						col = IM_COL32(255, 255, 128, 255)
					}
					var thickness float = 1.0
					if int(table.HoveredColumnBody) == column_n {
						thickness = 3.0
					}
					draw_list.AddRect(r.Min, r.Max, col, 0.0, 0, thickness)
				}
			} else {
				var r = GetTableRect(table, cfg.ShowTablesRectsType, -1)
				draw_list.AddRect(r.Min, r.Max, IM_COL32(255, 0, 128, 255), 0.0, 0, 1)
			}
		}
	}

	End()
}
