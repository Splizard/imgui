package imgui

import "fmt"

// Widgets: Combo Box
// - The BeginCombo()/EndCombo() api allows you to manage your contents and selection state however you want it, by creating e.g. Selectable() items.
// - The old Combo() api are helpers over BeginCombo()/EndCombo() which are kept available for convenience purpose. This is analogous to how ListBox are created.

func BeginCombo(label string, preview_value string, flags ImGuiComboFlags) bool {
	var g = GImGui
	var window = GetCurrentWindow()

	var backup_next_window_data_flags ImGuiNextWindowDataFlags = g.NextWindowData.Flags
	g.NextWindowData.ClearFlags() // We behave like Begin() and need to consume those values
	if window.SkipItems {
		return false
	}

	var style = g.Style
	var id ImGuiID = window.GetIDs(label)
	IM_ASSERT((flags & (ImGuiComboFlags_NoArrowButton | ImGuiComboFlags_NoPreview)) != (ImGuiComboFlags_NoArrowButton | ImGuiComboFlags_NoPreview)) // Can't use both flags together

	var arrow_size float = 0.0
	if (flags & ImGuiComboFlags_NoArrowButton) == 0 {
		arrow_size = GetFrameHeight()
	}
	var label_size ImVec2 = CalcTextSize(label, true, -1)
	var w = arrow_size
	if (flags & ImGuiComboFlags_NoPreview) == 0 {
		w = CalcItemWidth()
	}
	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(ImVec2{w, label_size.y + style.FramePadding.y*2.0})}

	var add float
	if label_size.x > 0.0 {
		add = style.ItemInnerSpacing.x + +label_size.x
	}

	var total_bb = ImRect{bb.Min, bb.Max.Add(ImVec2{add, 0.0})}
	ItemSizeRect(&total_bb, style.FramePadding.y)
	if !ItemAdd(&total_bb, id, &bb, 0) {
		return false
	}

	// Open on click
	var hovered, held bool
	var pressed bool = ButtonBehavior(&bb, id, &hovered, &held, 0)
	var popup_id ImGuiID = ImHashStr("##ComboPopup", 0, id)
	var popup_open bool = isPopupOpen(popup_id, ImGuiPopupFlags_None)
	if (pressed || g.NavActivateId == id) && !popup_open {
		OpenPopupEx(popup_id, ImGuiPopupFlags_None)
		popup_open = true
	}

	var c = ImGuiCol_FrameBg
	if hovered {
		c = ImGuiCol_FrameBgHovered
	}

	var rounding = ImDrawFlags_RoundCornersLeft
	if flags&ImGuiComboFlags_NoArrowButton != 0 {
		rounding = ImDrawFlags_RoundCornersAll
	}

	// Render shape
	var frame_col ImU32 = GetColorU32FromID(c, 1)
	var value_x2 float = ImMax(bb.Min.x, bb.Max.x-arrow_size)
	RenderNavHighlight(&bb, id, 0)
	if flags&ImGuiComboFlags_NoPreview == 0 {
		window.DrawList.AddRectFilled(bb.Min, ImVec2{value_x2, bb.Max.y}, frame_col, style.FrameRounding, rounding)
	}
	if flags&ImGuiComboFlags_NoArrowButton == 0 {
		var c = ImGuiCol_Button
		if popup_open || hovered {
			c = ImGuiCol_ButtonHovered
		}

		var bg_col ImU32 = GetColorU32FromID(c, 1)
		var text_col ImU32 = GetColorU32FromID(ImGuiCol_Text, 1)

		var rounding = ImDrawFlags_RoundCornersRight
		if w <= arrow_size {
			rounding = ImDrawFlags_RoundCornersAll
		}

		window.DrawList.AddRectFilled(ImVec2{value_x2, bb.Min.y}, bb.Max, bg_col, style.FrameRounding, rounding)
		if value_x2+arrow_size-style.FramePadding.x <= bb.Max.x {
			RenderArrow(window.DrawList, ImVec2{value_x2 + style.FramePadding.y, bb.Min.y + style.FramePadding.y}, text_col, ImGuiDir_Down, 1.0)
		}
	}
	RenderFrameBorder(bb.Min, bb.Max, style.FrameRounding)

	// Custom preview
	if flags&ImGuiComboFlags_CustomPreview != 0 {
		g.ComboPreviewData.PreviewRect = ImRect{ImVec2{bb.Min.x, bb.Min.y}, ImVec2{value_x2, bb.Max.y}}
		IM_ASSERT(preview_value == "" || preview_value[0] == 0)
		preview_value = ""
	}

	// Render preview and label
	if preview_value != "" && (flags&ImGuiComboFlags_NoPreview == 0) {
		if g.LogEnabled {
			LogSetNextTextDecoration("{", "}")
		}
		min := bb.Min.Add(style.FramePadding)
		RenderTextClipped(&min, &ImVec2{value_x2, bb.Max.y}, preview_value, nil, nil, nil)
	}
	if label_size.x > 0 {
		RenderText(ImVec2{bb.Max.x + style.ItemInnerSpacing.x, bb.Min.y + style.FramePadding.y}, label, true)
	}

	if !popup_open {
		return false
	}

	g.NextWindowData.Flags = backup_next_window_data_flags
	return BeginComboPopup(popup_id, &bb, flags)
}

func BeginComboPreview() bool { panic("not implemented") }
func EndComboPreview()        { panic("not implemented") }

func BeginComboPopup(popup_id ImGuiID, bb *ImRect, flags ImGuiComboFlags) bool {
	var g = GImGui
	if !isPopupOpen(popup_id, ImGuiPopupFlags_None) {
		g.NextWindowData.ClearFlags()
		return false
	}

	// Set popup size
	var w float = bb.GetWidth()
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasSizeConstraint != 0 {
		g.NextWindowData.SizeConstraintRect.Min.x = ImMax(g.NextWindowData.SizeConstraintRect.Min.x, w)
	} else {
		if (flags & ImGuiComboFlags_HeightMask_) == 0 {
			flags |= ImGuiComboFlags_HeightRegular
		}
		IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiComboFlags_HeightMask_))) // Only one
		var popup_max_height_in_items int = -1
		if flags&ImGuiComboFlags_HeightRegular != 0 {
			popup_max_height_in_items = 8
		} else if flags&ImGuiComboFlags_HeightSmall != 0 {
			popup_max_height_in_items = 4
		} else if flags&ImGuiComboFlags_HeightLarge != 0 {
			popup_max_height_in_items = 20
		}
		SetNextWindowSizeConstraints(ImVec2{w, 0.0}, ImVec2{FLT_MAX, CalcMaxPopupHeightFromItemCount(popup_max_height_in_items)}, nil, nil)
	}

	// This is essentially a specialized version of BeginPopupEx()
	var name = fmt.Sprintf("##Combo_%d", len(g.BeginPopupStack)) // Recycle windows based on depth

	// Set position given a custom constraint (peak into expected window size so we can position it)
	// FIXME: This might be easier to express with an hypothetical SetNextWindowPosConstraints() function?
	// FIXME: This might be moved to Begin() or at least around the same spot where Tooltips and other Popups are calling FindBestWindowPosForPopupEx()?
	if popup_window := FindWindowByName(string(name[:])); popup_window != nil {
		if popup_window.WasActive {
			// Always override 'AutoPosLastDirection' to not leave a chance for a past value to affect us.
			var size_expected ImVec2 = CalcWindowNextAutoFitSize(popup_window)

			dir := ImGuiDir_Down
			if (flags & ImGuiComboFlags_PopupAlignLeft) != 0 {
				dir = ImGuiDir_Left
			}

			popup_window.AutoPosLastDirection = dir // Left = "Below, Toward Left", Down = "Below, Toward Right (default)"
			var r_outer ImRect = GetPopupAllowedExtentRect(popup_window)
			var bl = bb.GetBL()
			var pos ImVec2 = FindBestWindowPosForPopupEx(&bl, &size_expected, &popup_window.AutoPosLastDirection, &r_outer, bb, ImGuiPopupPositionPolicy_ComboBox)
			SetNextWindowPos(&pos, 0, ImVec2{})
		}
	}

	// We don't use BeginPopupEx() solely because we have a custom name string, which we could make an argument to BeginPopupEx()
	var window_flags ImGuiWindowFlags = ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_Popup | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoSavedSettings | ImGuiWindowFlags_NoMove
	PushStyleVec(ImGuiStyleVar_WindowPadding, ImVec2{g.Style.FramePadding.x, g.Style.WindowPadding.y}) // Horizontally align ourselves with the framed text
	var ret bool = Begin(string(name[:]), nil, window_flags)
	PopStyleVar(1)
	if !ret {
		EndPopup()
		IM_ASSERT(false) // This should never happen as we tested for IsPopupOpen() above
		return false
	}
	return true
}

// only call EndCombo() if BeginCombo() returns true!
func EndCombo() {
	EndPopup()
}

func ComboString(label string, current_item *int, items_separated_by_zeros string, popup_max_height_in_items int /*= -1*/) bool {
	panic("not implemented")
} // Separate items with \0 within a string, end item-list with \0\0. e.g. "One\0Two\0Three\0"

func CalcMaxPopupHeightFromItemCount(items_count int) float32 {
	var g = GImGui
	if items_count <= 0 {
		return FLT_MAX
	}
	return (g.FontSize+g.Style.ItemSpacing.y)*float(items_count) - g.Style.ItemSpacing.y + (g.Style.WindowPadding.y * 2)
}

// Combo box helper allowing to pass an array of strings.
func ComboSlice(label string, current_item *int, items []string, items_count int, popup_max_height_in_items int /*= -1*/) bool {
	var value_changed = ComboFunc(label, current_item, func(slice interface{}, idx int, val *string) bool {
		*val = slice.([]string)[idx]
		return true
	}, items, items_count, popup_max_height_in_items)
	return value_changed
}

// Old API, prefer using BeginCombo() nowadays if you can.
func ComboFunc(label string, current_item *int, items_getter func(data interface{}, idx int, out_text *string) bool, data interface{}, items_count, popup_max_height_in_items int /*= -1*/) bool {
	var g = GImGui

	// Call the getter to obtain the preview string which is a parameter to BeginCombo()
	var preview_value string
	if *current_item >= 0 && *current_item < items_count {
		items_getter(data, *current_item, &preview_value)
	}

	// The old Combo() API exposed "popup_max_height_in_items". The new more general BeginCombo() API doesn't have/need it, but we emulate it here.
	if popup_max_height_in_items != -1 && (g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasSizeConstraint == 0) {
		SetNextWindowSizeConstraints(ImVec2{}, ImVec2{FLT_MAX, CalcMaxPopupHeightFromItemCount(popup_max_height_in_items)}, nil, nil)
	}

	if !BeginCombo(label, preview_value, ImGuiComboFlags_None) {
		return false
	}

	// Display items
	// FIXME-OPT: Use clipper (but we need to disable it on the appearing frame to make sure our call to SetItemDefaultFocus() is processed)
	var value_changed bool = false
	for i := int(0); i < items_count; i++ {
		PushID(i)
		var item_selected bool = (i == *current_item)
		var item_text string
		if !items_getter(data, i, &item_text) {
			item_text = "*Unknown item*"
		}
		if Selectable(item_text, item_selected, 0, ImVec2{}) {
			value_changed = true
			*current_item = i
		}
		if item_selected {
			SetItemDefaultFocus()
		}
		PopID()
	}

	EndCombo()

	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}
