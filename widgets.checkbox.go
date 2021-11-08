package imgui

func Checkbox(label string, v *bool) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style
	var id ImGuiID = window.GetIDs(label)
	var label_size ImVec2 = CalcTextSize(label, true, -1)

	var square_sz float = GetFrameHeight()
	var pos ImVec2 = window.DC.CursorPos

	var x float
	if label_size.x > 0 {
		x = style.ItemInnerSpacing.x + label_size.x
	}

	var total_bb = ImRect{pos, pos.Add(ImVec2{square_sz + x, label_size.y + style.FramePadding.y*2.0})}
	ItemSizeRect(&total_bb, style.FramePadding.y)
	if !ItemAdd(&total_bb, id, nil, 0) {
		return false
	}

	var hovered, held bool
	var pressed bool = ButtonBehavior(&total_bb, id, &hovered, &held, 0)
	if pressed {
		*v = !(*v)
		MarkItemEdited(id)
	}

	c := ImGuiCol_FrameBg
	if held && hovered {
		c = ImGuiCol_FrameBgActive
	} else if hovered {
		c = ImGuiCol_FrameBgHovered
	}

	var check_bb = ImRect{pos, pos.Add(ImVec2{square_sz, square_sz})}
	RenderNavHighlight(&total_bb, id, 0)
	RenderFrame(check_bb.Min, check_bb.Max, GetColorU32FromID(c, 1), true, style.FrameRounding)
	var check_col ImU32 = GetColorU32FromID(ImGuiCol_CheckMark, 1)
	var mixed_value bool = (g.LastItemData.InFlags & ImGuiItemFlags_MixedValue) != 0
	if mixed_value {
		// Undocumented tristate/mixed/indeterminate checkbox (#2644)
		// This may seem awkwardly designed because the aim is to make ImGuiItemFlags_MixedValue supported by all widgets (not just checkbox)
		var pad = ImVec2{ImMax(1.0, IM_FLOOR(square_sz/3.6)), ImMax(1.0, IM_FLOOR(square_sz/3.6))}
		window.DrawList.AddRectFilled(check_bb.Min.Add(pad), check_bb.Max.Sub(pad), check_col, style.FrameRounding, 0)
	} else if *v {
		var pad float = ImMax(1.0, IM_FLOOR(square_sz/6.0))
		RenderCheckMark(window.DrawList, check_bb.Min.Add(ImVec2{pad, pad}), check_col, square_sz-pad*2.0)
	}

	var label_pos ImVec2 = ImVec2{check_bb.Max.x + style.ItemInnerSpacing.x, check_bb.Min.y + style.FramePadding.y}
	if g.LogEnabled {
		s := "[ ]"
		if mixed_value {
			s = "[~]"
		} else if *v {
			s = "[x]"
		}
		LogRenderedText(&label_pos, s)
	}
	if label_size.x > 0.0 {
		RenderText(label_pos, label, true)
	}

	return pressed
}

func CheckboxFlagsInt(label string, flags *int, flags_value int) bool     { panic("not implemented") }
func CheckboxFlagsUint(label string, flags *uint, uflags_value uint) bool { panic("not implemented") }

func RenderCheckMark(draw_list *ImDrawList, pos ImVec2, col ImU32, sz float) {
	var thickness float = ImMax(sz/5.0, 1.0)
	sz -= thickness * 0.5
	pos = pos.Add(ImVec2{thickness * 0.25, thickness * 0.25})

	var third float = sz / 3.0
	var bx float = pos.x + third
	var by float = pos.y + sz - third*0.5
	draw_list.PathLineTo(ImVec2{bx - third, by - third})
	draw_list.PathLineTo(ImVec2{bx, by})
	draw_list.PathLineTo(ImVec2{bx + third*2.0, by - third*2.0})
	draw_list.PathStroke(col, 0, thickness)
}
