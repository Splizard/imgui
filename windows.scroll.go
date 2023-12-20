package imgui

// Scrolling
// Windows Scrolling

// GetScrollX get scrolling amount [0 .. GetScrollMaxX()]
func GetScrollX() float {
	window := GImGui.CurrentWindow
	return window.Scroll.x
}

// GetScrollY get scrolling amount [0 .. GetScrollMaxY()]
func GetScrollY() float {
	window := GImGui.CurrentWindow
	return window.Scroll.y
}

// SetScrollX set scrolling amount [0 .. GetScrollMaxX()]
func SetScrollX(scroll_x float) {
	g := GImGui
	setScrollX(g.CurrentWindow, scroll_x)
}

// SetScrollY set scrolling amount [0 .. GetScrollMaxY()]
func SetScrollY(scroll_y float) {
	g := GImGui
	setScrollY(g.CurrentWindow, scroll_y)
}

// GetScrollMaxX get maximum scrolling amount ~~ ContentSize.x - WindowSize.x - DecorationsSize.x
func GetScrollMaxX() float {
	window := GImGui.CurrentWindow
	return window.ScrollMax.x
}

// GetScrollMaxY get maximum scrolling amount ~~ ContentSize.y - WindowSize.y - DecorationsSize.y
func GetScrollMaxY() float {
	window := GImGui.CurrentWindow
	return window.ScrollMax.y
}

// SetScrollHereX adjust scrolling amount to make current cursor position visible. center_x_ratio=0.0: left, 0.5: center, 1.0: right. When using to make a "default/current item" visible, consider using SetItemDefaultFocus() instead.
// center_x_ratio: 0.0f left of last item, 0.5f horizontal center of last item, 1.0f right of last item.
func SetScrollHereX(center_x_ratio float /*= 0.5*/) {
	g := GImGui
	window := g.CurrentWindow
	var spacing_x = ImMax(window.WindowPadding.x, g.Style.ItemSpacing.x)
	var target_pos_x = ImLerp(g.LastItemData.Rect.Min.x-spacing_x, g.LastItemData.Rect.Max.x+spacing_x, center_x_ratio)
	setScrollFromPosX(window, target_pos_x-window.Pos.x, center_x_ratio) // Convert from absolute to local pos

	// Tweak: snap on edges when aiming at an item very close to the edge
	window.ScrollTargetEdgeSnapDist.x = ImMax(0.0, window.WindowPadding.x-spacing_x)
}

// SetScrollHereY adjust scrolling amount to make current cursor position visible. center_y_ratio=0.0: top, 0.5: center, 1.0: bottom. When using to make a "default/current item" visible, consider using SetItemDefaultFocus() instead.
// center_y_ratio: 0.0f top of last item, 0.5f vertical center of last item, 1.0f bottom of last item.
func SetScrollHereY(center_y_ratio float /*= 0.5*/) {
	g := GImGui
	window := g.CurrentWindow
	var spacing_y = ImMax(window.WindowPadding.y, g.Style.ItemSpacing.y)
	var target_pos_y = ImLerp(window.DC.CursorPosPrevLine.y-spacing_y, window.DC.CursorPosPrevLine.y+window.DC.PrevLineSize.y+spacing_y, center_y_ratio)
	setScrollFromPosY(window, target_pos_y-window.Pos.y, center_y_ratio) // Convert from absolute to local pos

	// Tweak: snap on edges when aiming at an item very close to the edge
	window.ScrollTargetEdgeSnapDist.y = ImMax(0.0, window.WindowPadding.y-spacing_y)
}

// SetScrollFromPosX adjust scrolling amount to make given position visible. Generally GetCursorStartPos() + offset to compute a valid position
func SetScrollFromPosX(local_x, center_x_ratio float /*= 0.5*/) {
	g := GImGui
	setScrollFromPosX(g.CurrentWindow, local_x, center_x_ratio)
}

// SetScrollFromPosY adjust scrolling amount to make given position visible. Generally GetCursorStartPos() + offset to compute a valid position.
func SetScrollFromPosY(local_y, center_y_ratio float /*= 0.5*/) {
	g := GImGui
	setScrollFromPosY(g.CurrentWindow, local_y, center_y_ratio)
}

// SetNextWindowScroll Use -1.0f on one axis to leave as-is
func SetNextWindowScroll(scroll *ImVec2) {
	g := GImGui
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasScroll
	g.NextWindowData.ScrollVal = *scroll
}

func setScrollX(window *ImGuiWindow, scroll_x float) {
	window.ScrollTarget.x = scroll_x
	window.ScrollTargetCenterRatio.x = 0.0
	window.ScrollTargetEdgeSnapDist.x = 0.0
}

func setScrollY(window *ImGuiWindow, scroll_y float) {
	window.ScrollTarget.y = scroll_y
	window.ScrollTargetCenterRatio.y = 0.0
	window.ScrollTargetEdgeSnapDist.y = 0.0
}

// Note that a local position will vary depending on initial scroll value,
// This is a little bit confusing so bear with us:
//   - local_pos = (absolution_pos - window.Pos)
//   - So local_x/local_y are 0.0f for a position at the upper-left corner of a window,
//     and generally local_x/local_y are >(padding+decoration) && <(size-padding-decoration) when in the visible area.
//   - They mostly exists because of legacy API.
//
// Following the rules above, when trying to work with scrolling code, consider that:
//   - SetScrollFromPosY(0.0f) == SetScrollY(0.0f + scroll.y) == has no effect!
//   - SetScrollFromPosY(-scroll.y) == SetScrollY(-scroll.y + scroll.y) == SetScrollY(0.0f) == reset scroll. Of course writing SetScrollY(0.0f) directly then makes more sense
//
// We store a target position so centering and clamping can occur on the next frame when we are guaranteed to have a known window size
func setScrollFromPosX(window *ImGuiWindow, local_x float, center_x_ratio float) {
	IM_ASSERT(center_x_ratio >= 0.0 && center_x_ratio <= 1.0)
	window.ScrollTarget.x = IM_FLOOR(local_x + window.Scroll.x) // Convert local position to scroll offset
	window.ScrollTargetCenterRatio.x = center_x_ratio
	window.ScrollTargetEdgeSnapDist.x = 0.0
}

func setScrollFromPosY(window *ImGuiWindow, local_y float, center_y_ratio float) {
	IM_ASSERT(center_y_ratio >= 0.0 && center_y_ratio <= 1.0)
	var decoration_up_height = window.TitleBarHeight() + window.MenuBarHeight() // FIXME: Would be nice to have a more standardized access to our scrollable/client rect;
	local_y -= decoration_up_height
	window.ScrollTarget.y = IM_FLOOR(local_y + window.Scroll.y) // Convert local position to scroll offset
	window.ScrollTargetCenterRatio.y = center_y_ratio
	window.ScrollTargetEdgeSnapDist.y = 0.0
}

func ScrollToBringRectIntoView(window *ImGuiWindow, item_rect *ImRect) ImVec2 {
	g := GImGui
	var window_rect = ImRect{window.InnerRect.Min.Sub(ImVec2{1, 1}), window.InnerRect.Max.Add(ImVec2{1, 1})}
	//GetForegroundDrawList(window).AddRect(window_rect.Min, window_rect.Max, IM_COL32_WHITE); // [DEBUG]

	var delta_scroll ImVec2
	if !window_rect.ContainsRect(*item_rect) {
		if window.ScrollbarX && item_rect.Min.x < window_rect.Min.x {
			setScrollFromPosX(window, item_rect.Min.x-window.Pos.x-g.Style.ItemSpacing.x, 0.0)
		} else if window.ScrollbarX && item_rect.Max.x >= window_rect.Max.x {
			setScrollFromPosX(window, item_rect.Max.x-window.Pos.x+g.Style.ItemSpacing.x, 1.0)
		}
		if item_rect.Min.y < window_rect.Min.y {
			setScrollFromPosY(window, item_rect.Min.y-window.Pos.y-g.Style.ItemSpacing.y, 0.0)
		} else if item_rect.Max.y >= window_rect.Max.y {
			setScrollFromPosY(window, item_rect.Max.y-window.Pos.y+g.Style.ItemSpacing.y, 1.0)
		}

		var next_scroll = CalcNextScrollFromScrollTargetAndClamp(window)
		delta_scroll = next_scroll.Sub(window.Scroll)
	}

	// Also scroll parent window to keep us into view if necessary
	if window.Flags&ImGuiWindowFlags_ChildWindow != 0 {
		delta_scroll = delta_scroll.Add(ScrollToBringRectIntoView(window.ParentWindow,
			&ImRect{item_rect.Min.Sub(delta_scroll), item_rect.Max.Sub(delta_scroll)}))
	}

	return delta_scroll
}

// CalcScrollEdgeSnap Helper to snap on edges when aiming at an item very close to the edge,
// So the difference between WindowPadding and ItemSpacing will be in the visible area after scrolling.
// When we refactor the scrolling API this may be configurable with a flag?
// Note that the effect for this won't be visible on X axis with default Style settings as WindowPadding.x == ItemSpacing.x by default.
func CalcScrollEdgeSnap(target, snap_min, snap_max, snap_threshold, center_ratio float) float {
	if target <= snap_min+snap_threshold {
		return ImLerp(snap_min, target, center_ratio)
	}
	if target >= snap_max-snap_threshold {
		return ImLerp(target, snap_max, center_ratio)
	}
	return target
}

func CalcNextScrollFromScrollTargetAndClamp(window *ImGuiWindow) ImVec2 {
	var scroll = window.Scroll
	if window.ScrollTarget.x < FLT_MAX {
		var decoration_total_width = window.ScrollbarSizes.x
		var center_x_ratio = window.ScrollTargetCenterRatio.x
		var scroll_target_x = window.ScrollTarget.x
		if window.ScrollTargetEdgeSnapDist.x > 0.0 {
			var snap_x_min float = 0.0
			var snap_x_max = window.ScrollMax.x + window.SizeFull.x - decoration_total_width
			scroll_target_x = CalcScrollEdgeSnap(scroll_target_x, snap_x_min, snap_x_max, window.ScrollTargetEdgeSnapDist.x, center_x_ratio)
		}
		scroll.x = scroll_target_x - center_x_ratio*(window.SizeFull.x-decoration_total_width)
	}
	if window.ScrollTarget.y < FLT_MAX {
		var decoration_total_height = window.TitleBarHeight() + window.MenuBarHeight() + window.ScrollbarSizes.y
		var center_y_ratio = window.ScrollTargetCenterRatio.y
		var scroll_target_y = window.ScrollTarget.y
		if window.ScrollTargetEdgeSnapDist.y > 0.0 {
			var snap_y_min float = 0.0
			var snap_y_max = window.ScrollMax.y + window.SizeFull.y - decoration_total_height
			scroll_target_y = CalcScrollEdgeSnap(scroll_target_y, snap_y_min, snap_y_max, window.ScrollTargetEdgeSnapDist.y, center_y_ratio)
		}
		scroll.y = scroll_target_y - center_y_ratio*(window.SizeFull.y-decoration_total_height)
	}
	scroll.x = IM_FLOOR(ImMax(scroll.x, 0.0))
	scroll.y = IM_FLOOR(ImMax(scroll.y, 0.0))
	if !window.Collapsed && !window.SkipItems {
		scroll.x = ImMin(scroll.x, window.ScrollMax.x)
		scroll.y = ImMin(scroll.y, window.ScrollMax.y)
	}
	return scroll
}
