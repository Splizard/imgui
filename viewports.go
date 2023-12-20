package imgui

// Update viewports and monitor infos
func UpdateViewportsNewFrame() {
	g := GImGui
	IM_ASSERT(len(g.Viewports) == 1)

	// Update main viewport with current platform position.
	// FIXME-VIEWPORT: Size is driven by backend/user code for backward-compatibility but we should aim to make this more consistent.
	var main_viewport = g.Viewports[0]
	main_viewport.Flags = ImGuiViewportFlags_IsPlatformWindow | ImGuiViewportFlags_OwnedByApp
	main_viewport.Pos = ImVec2{}
	main_viewport.Size = g.IO.DisplaySize

	for n := range g.Viewports {
		var viewport = g.Viewports[n]

		// Lock down space taken by menu bars and status bars, reset the offset for fucntions like BeginMainMenuBar() to alter them again.
		viewport.WorkOffsetMin = viewport.BuildWorkOffsetMin
		viewport.WorkOffsetMax = viewport.BuildWorkOffsetMax
		viewport.BuildWorkOffsetMin = ImVec2{}
		viewport.BuildWorkOffsetMax = ImVec2{}
		viewport.UpdateWorkRect()
	}
}

func SetupViewportDrawData(viewport *ImGuiViewportP, draw_lists *[]*ImDrawList) {
	var io = GetIO()
	var draw_data = &viewport.DrawDataP
	draw_data.Valid = true
	if len(*draw_lists) > 0 {
		draw_data.CmdLists = *draw_lists
	}
	draw_data.CmdListsCount = int(len(*draw_lists))
	draw_data.TotalVtxCount = 0
	draw_data.TotalIdxCount = 0
	draw_data.DisplayPos = viewport.Pos
	draw_data.DisplaySize = viewport.Size
	draw_data.FramebufferScale = io.DisplayFramebufferScale

	for _, v := range *draw_lists {

		draw_data.TotalVtxCount += int(len(v.VtxBuffer))
		draw_data.TotalIdxCount += int(len(v.IdxBuffer))
	}
}
