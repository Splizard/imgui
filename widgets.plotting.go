package imgui

//-------------------------------------------------------------------------
// [SECTION] Widgets: PlotLines, PlotHistogram
//-------------------------------------------------------------------------
// - PlotEx() [Internal]
// - PlotLines()
// - PlotHistogram()
//-------------------------------------------------------------------------
// Plot/Graph widgets are not very good.
// Consider writing your own, or using a third-party one, see:
// - ImPlot https://github.com/epezent/implot
// - others https://github.com/ocornut/imgui/wiki/Useful-Extensions
//-------------------------------------------------------------------------

type ImGuiPlotArrayGetterData struct {
	Values []float
	Stride int
}

func Plot_ArrayGetter(data any, idx int) float {
	var plot_data = (data).(*ImGuiPlotArrayGetterData)
	return plot_data.Values[idx*plot_data.Stride]
}

// Widgets: Data Plotting
// - Consider using ImPlot (https://github.com/epezent/implot) which is much better!
func PlotLines(label string, values []float, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/, stride int /*= sizeof(float)*/) {
	var data = ImGuiPlotArrayGetterData{values, stride}
	PlotEx(ImGuiPlotType_Lines, label, Plot_ArrayGetter, &data, values_count, values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func PlotLinesFunc(label string, values_getter func(data any, idx int) float, data any, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/) {
	PlotEx(ImGuiPlotType_Lines, label, values_getter, data, values_count, values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func PlotHistogram(label string, values []float, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/, stride int /* = sizeof(float)*/) {
	var data = ImGuiPlotArrayGetterData{values, stride}
	PlotEx(ImGuiPlotType_Histogram, label, Plot_ArrayGetter, &data, values_count, values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func PlotHistogramFunc(label string, values_getter func(data any, idx int) float, data any, values_count int, values_offset int /*= 0*/, overlay_text string /*= L*/, scale_min float /*= X*/, scale_max float /*= X*/, graph_size ImVec2 /*= 0*/) {
	PlotEx(ImGuiPlotType_Histogram, label, values_getter, data, values_count, values_offset, overlay_text, scale_min, scale_max, graph_size)
}

func PlotEx(plot_type ImGuiPlotType, label string, values_getter func(data any, idx int) float, data any, values_count int, values_offset int, overlay_text string, scale_min float, scale_max float, frame_size ImVec2) int {
	g := GImGui
	var window = GetCurrentWindow()
	if window.SkipItems {
		return -1
	}

	var style = g.Style
	var id = window.GetIDs(label)

	var label_size = CalcTextSize(label, true, -1)
	if frame_size.x == 0.0 {
		frame_size.x = CalcItemWidth()
	}
	if frame_size.y == 0.0 {
		frame_size.y = label_size.y + (style.FramePadding.y * 2)
	}

	var padding float = 0
	if label_size.x > 0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var frame_bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(frame_size)}
	var inner_bb = ImRect{frame_bb.Min.Add(style.FramePadding), frame_bb.Max.Sub(style.FramePadding)}
	var total_bb = ImRect{frame_bb.Min, frame_bb.Max.Add(ImVec2{padding, 0})}
	ItemSizeRect(&total_bb, style.FramePadding.y)
	if !ItemAdd(&total_bb, 0, &frame_bb, 0) {
		return -1
	}
	var hovered = ItemHoverable(&frame_bb, id)

	// Determine scale from values if not specified
	if scale_min == FLT_MAX || scale_max == FLT_MAX {
		var v_min float = FLT_MAX
		var v_max float = -FLT_MAX
		for i := int(0); i < values_count; i++ {
			var v = values_getter(data, i)
			if v != v { // Ignore NaN values
				continue
			}
			v_min = ImMin(v_min, v)
			v_max = ImMax(v_max, v)
		}
		if scale_min == FLT_MAX {
			scale_min = v_min
		}
		if scale_max == FLT_MAX {
			scale_max = v_max
		}
	}

	RenderFrame(frame_bb.Min, frame_bb.Max, GetColorU32FromID(ImGuiCol_FrameBg, 1), true, style.FrameRounding)

	var values_count_min int = 1
	if plot_type == ImGuiPlotType_Lines {
		values_count_min = 2
	}
	var idx_hovered int = -1
	if values_count >= values_count_min {
		var b int = 0
		if plot_type == ImGuiPlotType_Lines {
			b = -1
		}
		var res_w = ImMinInt((int)(frame_size.x), values_count) + b
		var item_count = values_count + b

		// Tooltip on hover
		if hovered && inner_bb.ContainsVec(g.IO.MousePos) {
			var t = ImClamp((g.IO.MousePos.x-inner_bb.Min.x)/(inner_bb.Max.x-inner_bb.Min.x), 0.0, 0.9999)
			var v_idx = (int)(t * float(item_count))
			IM_ASSERT(v_idx >= 0 && v_idx < values_count)

			var v0 = values_getter(data, (v_idx+values_offset)%values_count)
			var v1 = values_getter(data, (v_idx+1+values_offset)%values_count)
			if plot_type == ImGuiPlotType_Lines {
				SetTooltip("%d: %8.4g\n%d: %8.4g", v_idx, v0, v_idx+1, v1)
			} else if plot_type == ImGuiPlotType_Histogram {
				SetTooltip("%d: %8.4g", v_idx, v0)
			}
			idx_hovered = v_idx
		}

		var t_step = 1.0 / (float)(res_w)
		var inv_scale = (1.0 / (scale_max - scale_min))
		if scale_min == scale_max {
			inv_scale = 0
		}

		var v0 = values_getter(data, (0+values_offset)%values_count)
		var t0 float = 0.0
		var tp0 = ImVec2{t0, 1.0 - ImSaturate((v0-scale_min)*inv_scale)} // Point in the normalized space of our target rectangle
		var histogram_zero_line_t float                                  // Where does the zero line stands
		if scale_min*scale_max < 0.0 {
			histogram_zero_line_t = (1 + scale_min*inv_scale)
		} else {
			if scale_min < 0 {
				histogram_zero_line_t = 0
			} else {
				histogram_zero_line_t = 1
			}
		}

		c := ImGuiCol_PlotHistogram
		if plot_type == ImGuiPlotType_Lines {
			c = ImGuiCol_PlotLines
		}

		var col_base = GetColorU32FromID(c, 1)

		c = ImGuiCol_PlotHistogramHovered
		if plot_type == ImGuiPlotType_Lines {
			c = ImGuiCol_PlotLinesHovered
		}

		var col_hovered = GetColorU32FromID(c, 1)

		for n := int(0); n < res_w; n++ {
			var t1 = t0 + t_step
			var v1_idx = (int)(t0*float(item_count) + 0.5)
			IM_ASSERT(v1_idx >= 0 && v1_idx < values_count)
			var v1 = values_getter(data, (v1_idx+values_offset+1)%values_count)
			var tp1 = ImVec2{t1, 1.0 - ImSaturate((v1-scale_min)*inv_scale)}

			// NB: Draw calls are merged together by the DrawList system. Still, we should render our batch are lower level to save a bit of CPU.
			var pos0 = ImLerpVec2WithVec2(&inner_bb.Min, &inner_bb.Max, tp0)

			var t = ImVec2{tp1.x, histogram_zero_line_t}
			if plot_type == ImGuiPlotType_Lines {
				t = tp1
			}

			var pos1 = ImLerpVec2WithVec2(&inner_bb.Min, &inner_bb.Max, t)
			if plot_type == ImGuiPlotType_Lines {
				c := col_base
				if idx_hovered == v1_idx {
					c = col_hovered
				}
				window.DrawList.AddLine(&pos0, &pos1, c, 1)
			} else if plot_type == ImGuiPlotType_Histogram {
				if pos1.x >= pos0.x+2.0 {
					pos1.x -= 1.0
				}
				c := col_base
				if idx_hovered == v1_idx {
					c = col_hovered
				}
				window.DrawList.AddRectFilled(pos0, pos1, c, 0, 0)
			}

			t0 = t1
			tp0 = tp1
		}
	}

	// Text overlay
	if overlay_text != "" {
		RenderTextClipped(&ImVec2{frame_bb.Min.x, frame_bb.Min.y + style.FramePadding.y}, &frame_bb.Max, overlay_text, nil, &ImVec2{0.5, 0.0}, nil)
	}
	if label_size.x > 0.0 {
		RenderText(ImVec2{frame_bb.Max.x + style.ItemInnerSpacing.x, inner_bb.Min.y}, label, true)
	}

	// Return hovered index or -1 if none are hovered.
	// This is currently not exposed in the public API because we need a larger redesign of the whole thing, but in the short-term we are making it available in PlotEx().
	return idx_hovered
}
