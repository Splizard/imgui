package imgui

func RenderFrameBorder(p_min ImVec2, p_max ImVec2, rounding float) {
	g := GImGui
	window := g.CurrentWindow
	var border_size = g.Style.FrameBorderSize
	if border_size > 0.0 {
		window.DrawList.AddRect(p_min.Add(ImVec2{1, 1}), p_max.Add(ImVec2{1, 1}), GetColorU32FromID(ImGuiCol_BorderShadow, 1), rounding, 0, border_size)
		window.DrawList.AddRect(p_min, p_max, GetColorU32FromID(ImGuiCol_Border, 1), rounding, 0, border_size)
	}
}

func RenderBullet(draw_list *ImDrawList, pos ImVec2, col ImU32) {
	draw_list.AddCircleFilled(pos, draw_list._Data.FontSize*0.20, col, 8)
}

// RenderArrow Render an arrow aimed to be aligned with text (p_min is a position in the same space text would be positioned). To e.g. denote expanded/collapsed state
func RenderArrow(draw_list *ImDrawList, pos ImVec2, col ImU32, dir ImGuiDir, scale float /*= 1.0f*/) {
	var h = draw_list._Data.FontSize * 1.00
	var r = h * 0.40 * scale
	var center = pos.Add(ImVec2{h * 0.50, h * 0.50 * scale})

	var a, b, c ImVec2
	switch dir {
	case ImGuiDir_Up:
		fallthrough
	case ImGuiDir_Down:
		if dir == ImGuiDir_Up {
			r = -r
		}
		a = ImVec2{+0.000, +0.750}.Scale(r)
		b = ImVec2{-0.866, -0.750}.Scale(r)
		c = ImVec2{+0.866, -0.750}.Scale(r)
	case ImGuiDir_Left:
		fallthrough
	case ImGuiDir_Right:
		if dir == ImGuiDir_Left {
			r = -r
		}
		a = ImVec2{+0.750, +0.000}.Scale(r)
		b = ImVec2{-0.750, +0.866}.Scale(r)
		c = ImVec2{-0.750, -0.866}.Scale(r)
	case ImGuiDir_None:
		fallthrough
	case ImGuiDir_COUNT:
		IM_ASSERT(false)
	}

	p1, p2 := center.Add(a), center.Add(b)
	draw_list.AddTriangleFilled(&p1, &p2, center.Add(c), col)
}

// Render ends the Dear ImGui frame, finalize the draw data. You can then get call GetDrawData()
// Prepare the data for rendering so you can call GetDrawData()
// (As with anything within the ImGui:: namspace this doesn't touch your GPU or graphics API at all:
// it is the role of the ImGui_ImplXXXX_RenderDrawData() function provided by the renderer backend).
func Render() {
	g := GImGui
	IM_ASSERT(g.Initialized)

	if g.FrameCountEnded != g.FrameCount {
		EndFrame()
	}
	g.FrameCountRendered = g.FrameCount
	g.IO.MetricsRenderWindows = 0

	CallContextHooks(g, ImGuiContextHookType_RenderPre)

	// Add background ImDrawList (for each active viewport)
	for n := range g.Viewports {
		var viewport = g.Viewports[n]
		viewport.DrawDataBuilder.Clear()
		if viewport.DrawLists[0] != nil {
			AddDrawListToDrawData(&viewport.DrawDataBuilder[0], getBackgroundDrawList(viewport))
		}
	}

	// Add ImDrawList to render
	var windows_to_render_top_most [2]*ImGuiWindow
	if g.NavWindowingTarget != nil && g.NavWindowingTarget.Flags&ImGuiWindowFlags_NoBringToFrontOnFocus == 0 {
		windows_to_render_top_most[0] = g.NavWindowingTarget.RootWindow
	}
	if g.NavWindowingTarget != nil {
		windows_to_render_top_most[1] = g.NavWindowingListWindow
	}

	for n := range g.Windows {
		var window = g.Windows[n]
		if IsWindowActiveAndVisible(window) && (window.Flags&ImGuiWindowFlags_ChildWindow) == 0 && window != windows_to_render_top_most[0] && window != windows_to_render_top_most[1] {
			AddRootWindowToDrawData(window)
		}
	}
	for n := range windows_to_render_top_most {
		if windows_to_render_top_most[n] != nil && IsWindowActiveAndVisible(windows_to_render_top_most[n]) { // NavWindowingTarget is always temporarily displayed as the top-most window
			AddRootWindowToDrawData(windows_to_render_top_most[n])
		}
	}

	// Setup ImDrawData structures for end-user
	g.IO.MetricsRenderVertices = 0
	g.IO.MetricsRenderIndices = 0
	for n := range g.Viewports {
		var viewport = g.Viewports[n]
		viewport.DrawDataBuilder.FlattenIntoSingleLayer()

		// Draw software mouse cursor if requested by io.MouseDrawCursor flag
		if g.IO.MouseDrawCursor {
			RenderMouseCursor(GetForegroundDrawListViewport(viewport), g.IO.MousePos, g.Style.MouseCursorScale, g.MouseCursor, IM_COL32_WHITE, IM_COL32_BLACK, IM_COL32(0, 0, 0, 48))
		}

		// Add foreground ImDrawList (for each active viewport)
		if viewport.DrawLists[1] != nil {
			AddDrawListToDrawData(&viewport.DrawDataBuilder[0], GetForegroundDrawListViewport(viewport))
		}

		SetupViewportDrawData(viewport, &viewport.DrawDataBuilder[0])
		var draw_data = &viewport.DrawDataP
		g.IO.MetricsRenderVertices += draw_data.TotalVtxCount
		g.IO.MetricsRenderIndices += draw_data.TotalIdxCount
	}

	CallContextHooks(g, ImGuiContextHookType_RenderPost)
}
