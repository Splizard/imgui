package imgui

// separator, generally horizontal. inside a menu bar or in horizontal layout mode, this becomes a vertical separator.
func Separator() {
	g := GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return
	}

	// Those flags should eventually be overridable by the user
	var flags ImGuiSeparatorFlags
	if window.DC.LayoutType == ImGuiLayoutType_Horizontal {
		flags = ImGuiSeparatorFlags_Vertical
	} else {
		flags = ImGuiSeparatorFlags_Horizontal
	}

	flags |= ImGuiSeparatorFlags_SpanAllColumns
	SeparatorEx(flags)
}

// Horizontal/vertical separating line
func SeparatorEx(flags ImGuiSeparatorFlags) {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	g := GImGui
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & (ImGuiSeparatorFlags_Horizontal | ImGuiSeparatorFlags_Vertical)))) // Check that only 1 option is selected

	var thickness_draw float = 1.0
	var thickness_layout float = 0.0
	if flags&ImGuiSeparatorFlags_Vertical != 0 {
		// Vertical separator, for menu bars (use current line height). Not exposed because it is misleading and it doesn't have an effect on regular layout.
		var y1 = window.DC.CursorPos.y
		var y2 = window.DC.CursorPos.y + window.DC.CurrLineSize.y
		var bb = ImRect{ImVec2{window.DC.CursorPos.x, y1}, ImVec2{window.DC.CursorPos.x + thickness_draw, y2}}
		ItemSizeVec(&ImVec2{thickness_layout, 0.0}, 0)
		if !ItemAdd(&bb, 0, nil, 0) {
			return
		}

		// Draw
		window.DrawList.AddLine(&ImVec2{bb.Min.x, bb.Min.y}, &ImVec2{bb.Min.x, bb.Max.y}, GetColorU32FromID(ImGuiCol_Separator, 1), 1)
		if g.LogEnabled {
			LogText(" |")
		}
	} else if flags&ImGuiSeparatorFlags_Horizontal != 0 {
		// Horizontal Separator
		var x1 = window.Pos.x
		var x2 = window.Pos.x + window.Size.x

		// FIXME-WORKRECT: old hack (#205) until we decide of consistent behavior with WorkRect/Indent and Separator
		if len(g.GroupStack) > 0 && g.GroupStack[len(g.GroupStack)-1].WindowID == window.ID {
			x1 += window.DC.Indent.x
		}

		var columns *ImGuiOldColumns
		if flags&ImGuiSeparatorFlags_SpanAllColumns != 0 {
			columns = window.DC.CurrentColumns
		}
		if columns != nil {
			PushColumnsBackground()
		}

		// We don't provide our width to the layout so that it doesn't get feed back into AutoFit
		var bb = ImRect{ImVec2{x1, window.DC.CursorPos.y}, ImVec2{x2, window.DC.CursorPos.y + thickness_draw}}
		ItemSizeVec(&ImVec2{0.0, thickness_layout}, 0)
		var item_visible = ItemAdd(&bb, 0, nil, 0)
		if item_visible {
			// Draw
			window.DrawList.AddLine(&bb.Min, &ImVec2{bb.Max.x, bb.Min.y}, GetColorU32FromID(ImGuiCol_Separator, 1), 1)
			if g.LogEnabled {
				LogRenderedText(&bb.Min, "--------------------------------\n")
			}
		}
		if columns != nil {
			PopColumnsBackground()
			columns.LineMinY = window.DC.CursorPos.y
		}
	}
}
