package imgui

// Debug Tools

// Experimental recovery from incorrect usage of BeginXXX/EndXXX/PushXXX/PopXXX calls.
// Must be called during or before EndFrame().
// This is generally flawed as we are not necessarily End/Popping things in the right order.
// FIXME: Can't recover from inside BeginTabItem/EndTabItem yet.
// FIXME: Can't recover from interleaved BeginTabBar/Begin
func ErrorCheckEndFrameRecover(log_callback ImGuiErrorLogCallback, user_data interface{}) {
	var g = GImGui
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
	var g = GImGui
	var window = g.CurrentWindow

	var scale ImVec2 = bb.GetSize().Div(viewport.Size)
	var off ImVec2 = bb.Min.Sub(viewport.Pos.Mul(scale))
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

// Avoid naming collision with imgui_demo.cpp's HelpMarker() for unity builds.
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
	var g = GImGui
	var window = g.CurrentWindow

	// We don't display full monitor bounds (we could, but it often looks awkward), instead we display just enough to cover all of our viewports.
	const SCALE = 1.0 / 8.0
	var bb_full = ImRect{ImVec2{FLT_MAX, FLT_MAX}, ImVec2{-FLT_MAX, -FLT_MAX}}
	for n := range g.Viewports {
		bb_full.AddRect(g.Viewports[n].GetMainRect())
	}
	var p ImVec2 = window.DC.CursorPos
	var off ImVec2 = p.Sub(bb_full.Min.Scale(SCALE))
	for n := range g.Viewports {
		var viewport *ImGuiViewportP = g.Viewports[n]
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

// [DEBUG] Item picker tool - start with DebugStartItemPicker() - useful to visually select an item and break into its call-stack.
func UpdateDebugToolItemPicker() {
	var g = GImGui
	g.DebugItemPickerBreakId = 0
	if g.DebugItemPickerActive {
		var hovered_id ImGuiID = g.HoveredIdPreviousFrame
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

// create Metrics/Debugger window. display Dear ImGui internals: windows, draw commands, various internal state, etc.
func ShowMetricsWindow(p_open *bool) {
	if !Begin("Dear ImGui Metrics/Debugger", p_open, 0) {
		End()
		return
	}

	var g = GImGui
	var io = g.IO
	var cfg = &g.DebugMetricsConfig

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

	/*GetTableRect := func(table *ImGuiTable, rect_type, n int) ImRect {
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
		case TRT_ColumnsRect:
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
			return ImRect(c.WorkMinX, table.InnerClipRect.Min.y+table.LastFirstRowHeight, c.ContentMaxXUnfrozen, table.InnerClipRect.Max.y)
		}
		IM_ASSERT(false)
		return ImRect{}
	}*/

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
			var min ImVec2 = window.InnerRect.Min.Sub(window.Scroll).Add(window.WindowPadding)
			return ImRect{min, min.Add(window.ContentSize)}
		case WRT_ContentIdeal:
			var min ImVec2 = window.InnerRect.Min.Sub(window.Scroll).Add(window.WindowPadding)
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
		cfg.ShowWindowsRects = cfg.ShowWindowsRects || ComboSlice("##show_windows_rect_type", &cfg.ShowWindowsRectsType, wrt_rects_names, WRT_Count, WRT_Count)
		if cfg.ShowWindowsRects && g.NavWindow != nil {
			BulletText("'%s':", g.NavWindow.Name)
			Indent(0)
			for rect_n := int(0); rect_n < WRT_Count; rect_n++ {
				var r ImRect = GetWindowRect(g.NavWindow, rect_n)
				Text("(%6.1f,%6.1f) (%6.1f,%6.1f) Size (%6.1f,%6.1f) %s", r.Min.x, r.Min.y, r.Max.x, r.Max.y, r.GetWidth(), r.GetHeight(), wrt_rects_names[rect_n])
			}
			Unindent(0)
		}

		TreePop()
	}

	End()
}
