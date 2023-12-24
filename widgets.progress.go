package imgui

import "fmt"

// size_arg (for each axis) < 0.0f: align to end, 0.0f: auto, > 0.0f: specified size
func ProgressBar(fraction float, size_arg ImVec2 /*= ImVec2(-FLT_MIN, 0)*/, overlay string) {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	style := guiContext.Style

	var pos = window.DC.CursorPos
	var size = CalcItemSize(size_arg, CalcItemWidth(), guiContext.FontSize+style.FramePadding.y*2.0)
	var bb = ImRect{pos, pos.Add(size)}
	ItemSizeVec(&size, style.FramePadding.y)
	if !ItemAdd(&bb, 0, nil, 0) {
		return
	}

	// Render
	fraction = ImSaturate(fraction)
	RenderFrame(bb.Min, bb.Max, GetColorU32FromID(ImGuiCol_FrameBg, 1), true, style.FrameRounding)
	bb.ExpandVec(ImVec2{-style.FrameBorderSize, -style.FrameBorderSize})
	var fill_br = ImVec2{ImLerp(bb.Min.x, bb.Max.x, fraction), bb.Max.y}
	RenderRectFilledRangeH(window.DrawList, &bb, GetColorU32FromID(ImGuiCol_PlotHistogram, 1), 0.0, fraction, style.FrameRounding)

	// Default displaying the fraction as percentage string, but user can override it
	if overlay == "" {
		overlay = fmt.Sprintf("%.0f%%", fraction*100+0.01)
	}

	var overlay_size = CalcTextSize(overlay, true, -1)
	if overlay_size.x > 0.0 {
		RenderTextClipped(&ImVec2{ImClamp(fill_br.x+style.ItemSpacing.x, bb.Min.x, bb.Max.x-overlay_size.x-style.ItemInnerSpacing.x), bb.Min.y}, &bb.Max, overlay, &overlay_size, &ImVec2{0.0, 0.5}, &bb)
	}
}
