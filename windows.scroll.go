package imgui

// Helper to snap on edges when aiming at an item very close to the edge,
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
	var scroll ImVec2 = window.Scroll
	if window.ScrollTarget.x < FLT_MAX {
		var decoration_total_width float = window.ScrollbarSizes.x
		var center_x_ratio float = window.ScrollTargetCenterRatio.x
		var scroll_target_x float = window.ScrollTarget.x
		if window.ScrollTargetEdgeSnapDist.x > 0.0 {
			var snap_x_min float = 0.0
			var snap_x_max float = window.ScrollMax.x + window.SizeFull.x - decoration_total_width
			scroll_target_x = CalcScrollEdgeSnap(scroll_target_x, snap_x_min, snap_x_max, window.ScrollTargetEdgeSnapDist.x, center_x_ratio)
		}
		scroll.x = scroll_target_x - center_x_ratio*(window.SizeFull.x-decoration_total_width)
	}
	if window.ScrollTarget.y < FLT_MAX {
		var decoration_total_height float = window.TitleBarHeight() + window.MenuBarHeight() + window.ScrollbarSizes.y
		var center_y_ratio float = window.ScrollTargetCenterRatio.y
		var scroll_target_y float = window.ScrollTarget.y
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
