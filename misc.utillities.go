package imgui

// Miscellaneous Utilities

func GetViewportDrawList(viewport *ImGuiViewportP, drawlist_no size_t, drawlist_name string) *ImDrawList {
	// Create the draw list on demand, because they are not frequently used for all viewports
	g := GImGui
	IM_ASSERT(drawlist_no < size_t(len(viewport.DrawLists)))
	var draw_list = viewport.DrawLists[drawlist_no]
	if draw_list == nil {
		l := NewImDrawList(&g.DrawListSharedData)
		draw_list = &l
		draw_list._OwnerName = drawlist_name
		viewport.DrawLists[drawlist_no] = draw_list
	}

	// Our ImDrawList system requires that there is always a command
	if viewport.DrawListsLastFrame[drawlist_no] != g.FrameCount {
		draw_list._ResetForNewFrame()
		draw_list.PushTextureID(g.IO.Fonts.TexID)
		draw_list.PushClipRect(viewport.Pos, viewport.Pos.Add(viewport.Size), false)
		viewport.DrawListsLastFrame[drawlist_no] = g.FrameCount
	}
	return draw_list
}

// test if rectangle (of given size, starting from cursor position) is visible / not clipped.
func IsRectVisible(size ImVec2) bool {
	var window = GImGui.CurrentWindow
	return window.ClipRect.Overlaps(ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)})
}

// test if rectangle (in screen space) is visible / not clipped. to perform coarse clipping on user's side.
func IsRectVisibleMinMax(rect_min, rect_max ImVec2) bool {
	var window = GImGui.CurrentWindow
	return window.ClipRect.Overlaps(ImRect{rect_min, rect_max})
}

func GetTime() double    { return GImGui.Time }       // get global imgui time. incremented by io.DeltaTime every frame.
func GetFrameCount() int { return GImGui.FrameCount } // get global imgui frame count. incremented by 1 every frame.

// this draw list will be the first rendering one. Useful to quickly draw shapes/text behind dear imgui contents.
func GetBackgroundDrawList(viewport *ImGuiViewport) *ImDrawList {
	g := GImGui
	if viewport == nil {
		viewport = g.Viewports[0]
	}
	return GetViewportDrawList(viewport, 0, "##Background")
}

// this draw list will be the last rendered one. Useful to quickly draw shapes/text over dear imgui contents.
func GetForegroundDrawList(viewport *ImGuiViewport) *ImDrawList {
	g := GImGui
	if viewport == nil {
		viewport = g.Viewports[0]
	}
	return GetViewportDrawList(viewport, 1, "##Foreground")
}

// you may use this when creating your own ImDrawList instances.
func GetDrawListSharedData() *ImDrawListSharedData {
	return &GImGui.DrawListSharedData
}

// replace current window storage with our own (if you want to manipulate it yourself, typically clear subsection of it)
func SetStateStorage(storage *ImGuiStorage) {
	var window = GImGui.CurrentWindow
	if storage != nil {
		window.DC.StateStorage = *storage
	} else {
		window.DC.StateStorage = window.StateStorage
	}
}

func GetStateStorage() ImGuiStorage {
	var window = GImGui.CurrentWindow
	return window.DC.StateStorage
}
