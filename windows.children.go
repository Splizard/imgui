package imgui

import "fmt"

// Popups, Modals, Tooltips
func BeginChildEx(name string, id ImGuiID, size_arg *ImVec2, border bool, flags ImGuiWindowFlags) bool {
	var g = GImGui
	var parent_window *ImGuiWindow = g.CurrentWindow

	flags |= ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoSavedSettings | ImGuiWindowFlags_ChildWindow
	flags |= (parent_window.Flags & ImGuiWindowFlags_NoMove) // Inherit the NoMove flag

	// Size
	var content_avail ImVec2 = GetContentRegionAvail()
	var size ImVec2 = *ImFloorVec(size_arg)

	var auto_fit_axises int
	if size.x == 0.0 {
		auto_fit_axises = (1 << ImGuiAxis_X)
	}
	if size.y == 0.0 {
		auto_fit_axises |= (1 << ImGuiAxis_Y)
	}

	if size.x <= 0.0 {
		size.x = ImMax(content_avail.x+size.x, 4.0) // Arbitrary minimum child size (0.0f causing too much issues)
	}
	if size.y <= 0.0 {
		size.y = ImMax(content_avail.y+size.y, 4.0)
	}
	SetNextWindowSize(&size, 0)

	// Build up name. If you need to append to a same child from multiple location in the ID stack, use BeginChild(ImGuiID id) with a stable value.
	if name != "" {
		g.TempBuffer = fmt.Sprintf("%s/%s_%08X", parent_window.Name, name, id)
	} else {
		g.TempBuffer = fmt.Sprintf("%s/%08X", parent_window.Name, id)
	}

	var backup_border_size float = g.Style.ChildBorderSize
	if !border {
		g.Style.ChildBorderSize = 0.0
	}
	var ret bool = Begin(string(g.TempBuffer[:]), nil, flags)
	g.Style.ChildBorderSize = backup_border_size

	var child_window *ImGuiWindow = g.CurrentWindow
	child_window.ChildId = id
	child_window.AutoFitChildAxises = (ImS8)(auto_fit_axises)

	// Set the cursor to handle case where the user called SetNextWindowPos()+BeginChild() manually.
	// While this is not really documented/defined, it seems that the expected thing to do.
	if child_window.BeginCount == 1 {
		parent_window.DC.CursorPos = child_window.Pos
	}

	// Process navigation-in immediately so NavInit can run on first frame
	if g.NavActivateId == id && (flags&ImGuiWindowFlags_NavFlattened == 0) && (child_window.DC.NavLayersActiveMask != 0 || child_window.DC.NavHasScroll) {
		FocusWindow(child_window)
		NavInitWindow(child_window, false)
		SetActiveID(id+1, child_window) // Steal ActiveId with another arbitrary id so that key-press won't activate child item
		g.ActiveIdSource = ImGuiInputSource_Nav
	}
	return ret
}

// Child Windows
// - Use child windows to begin into a self-contained independent scrolling/clipping regions within a host window. Child windows can embed their own child.
// - For each independent axis of 'size': ==0.0: use remaining host window size / >0.0: fixed size / <0.0: use remaining window size minus abs(size) / Each axis can use a different mode, e.g. ImVec2(0,400).
// - BeginChild() returns false to indicate the window is collapsed or fully clipped, so you may early out and omit submitting anything to the window.
//   Always call a matching EndChild() for each BeginChild() call, regardless of its return value.
//   [Important: due to legacy reason, this is inconsistent with most other functions such as BeginMenu/EndMenu,
//    BeginPopup/EndPopup, etc. where the EndXXX call should only be called if the corresponding BeginXXX function
//    returned true. Begin and BeginChild are the only odd ones out. Will be fixed in a future update.]
func BeginChild(str_id string, size ImVec2, border bool, flags ImGuiWindowFlags) bool {
	var window = GetCurrentWindow()
	return BeginChildEx(str_id, window.GetIDs(str_id), &size, border, flags)
}

func BeginChildID(id ImGuiID, size ImVec2, border bool, flags ImGuiWindowFlags) bool {
	IM_ASSERT(id != 0)
	return BeginChildEx("", id, &size, border, flags)
}

func EndChild() {
	var g = GImGui
	var window = g.CurrentWindow

	IM_ASSERT(g.WithinEndChild == false)
	IM_ASSERT(window.Flags&ImGuiWindowFlags_ChildWindow != 0) // Mismatched BeginChild()/EndChild() calls

	g.WithinEndChild = true
	if window.BeginCount > 1 {
		End()
	} else {
		var sz ImVec2 = window.Size
		if window.AutoFitChildAxises&(1<<ImGuiAxis_X) != 0 { // Arbitrary minimum zero-ish child size of 4.0f causes less trouble than a 0.0f
			sz.x = ImMax(4.0, sz.x)
		}
		if window.AutoFitChildAxises&(1<<ImGuiAxis_Y) != 0 {
			sz.y = ImMax(4.0, sz.y)
		}
		End()

		var parent_window *ImGuiWindow = g.CurrentWindow
		var bb = ImRect{parent_window.DC.CursorPos, parent_window.DC.CursorPos.Add(sz)}
		ItemSizeVec(&sz, 0)
		if (window.DC.NavLayersActiveMask != 0 || window.DC.NavHasScroll) && (window.Flags&ImGuiWindowFlags_NavFlattened == 0) {
			ItemAdd(&bb, window.ChildId, nil, 0)
			RenderNavHighlight(&bb, window.ChildId, 0)

			// When browsing a window that has no activable items (scroll only) we keep a highlight on the child
			if window.DC.NavLayersActiveMask == 0 && window == g.NavWindow {
				RenderNavHighlight(&ImRect{bb.Min.Sub(ImVec2{2, 2}), bb.Max.Add(ImVec2{2, 2})}, g.NavId, ImGuiNavHighlightFlags_TypeThin)
			}
		} else {
			// Not navigable into
			ItemAdd(&bb, 0, nil, 0)
		}
		if g.HoveredWindow == window {
			g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HoveredWindow
		}
	}
	g.WithinEndChild = false
	g.LogLinePosY = -FLT_MAX // To enforce a carriage return
}

// calculate coarse clipping for large list of evenly sized items. Prefer using the ImGuiListClipper higher-level helper if you can.
// helper to create a child window / scrolling region that looks like a normal widget frame
func BeginChildFrame(id ImGuiID, size ImVec2, flags ImGuiWindowFlags) bool {
	var g = GImGui
	var style = g.Style
	PushStyleColorVec(ImGuiCol_ChildBg, &style.Colors[ImGuiCol_FrameBg])
	PushStyleFloat(ImGuiStyleVar_ChildRounding, style.FrameRounding)
	PushStyleFloat(ImGuiStyleVar_ChildBorderSize, style.FrameBorderSize)
	PushStyleVec(ImGuiStyleVar_WindowPadding, style.FramePadding)
	var ret bool = BeginChildID(id, size, true, ImGuiWindowFlags_NoMove|ImGuiWindowFlags_AlwaysUseWindowPadding|flags)
	PopStyleVar(3)
	PopStyleColor(1)
	return ret
}

// always call EndChildFrame() regardless of BeginChildFrame() return values (which indicates a collapsed/clipped window)
func EndChildFrame() {
	EndChild()
}
