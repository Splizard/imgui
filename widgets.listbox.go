package imgui

// Widgets: List Boxes
// - This is essentially a thin wrapper to using BeginChild/EndChild with some stylistic changes.
// - The BeginListBox()/EndListBox() api allows you to manage your contents and selection state however you want it, by creating e.g. Selectable() or any items.
// - The simplified/old ListBox() api are helpers over BeginListBox()/EndListBox() which are kept available for convenience purpose. This is analoguous to how Combos are created.
// - Choose frame width:   size.x > 0.0: custom  /  size.x < 0.0 or -FLT_MIN: right-align   /  size.x = 0.0 (default): use current ItemWidth
// - Choose frame height:  size.y > 0.0: custom  /  size.y < 0.0 or -FLT_MIN: bottom-align  /  size.y = 0.0 (default): arbitrary default height which can fit ~7 items

// open a framed scrolling region
// Tip: To have a list filling the entire window width, use size.x = -FLT_MIN and pass an non-visible label e.g. "##empty"
// Tip: If your vertical size is calculated from an item count (e.g. 10 * item_height) consider adding a fractional part to facilitate seeing scrolling boundaries (e.g. 10.25 * item_height).
func BeginListBox(label string, size_arg ImVec2) bool {
	g := GImGui
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var style = g.Style
	var id = GetIDs(label)
	var label_size = CalcTextSize(label, true, -1)

	// Size default to hold ~7.25 items.
	// Fractional number of items helps seeing that we can scroll down/up without looking at scrollbar.
	var size_not_floored = CalcItemSize(size_arg, CalcItemWidth(), GetTextLineHeightWithSpacing()*7.25+style.FramePadding.y*2.0)
	var size = ImFloorVec(&size_not_floored)
	var frame_size = ImVec2{size.x, ImMax(size.y, label_size.y)}
	var frame_bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(frame_size)}

	var padding float = 0
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var bb = ImRect{frame_bb.Min, frame_bb.Max.Add(ImVec2{padding, 0.0})}
	g.NextItemData.ClearFlags()

	if !IsRectVisibleMinMax(bb.Min, bb.Max) {
		v := bb.GetSize()
		ItemSizeVec(&v, style.FramePadding.y)
		ItemAdd(&bb, 0, &frame_bb, 0)
		return false
	}

	// FIXME-OPT: We could omit the BeginGroup() if label_size.x but would need to omit the EndGroup() as well.
	BeginGroup()
	if label_size.x > 0.0 {
		var label_pos = ImVec2{frame_bb.Max.x + style.ItemInnerSpacing.x, frame_bb.Min.y + style.FramePadding.y}
		RenderText(label_pos, label, true)
		b := label_pos.Add(label_size)
		window.DC.CursorMaxPos = ImMaxVec2(&window.DC.CursorMaxPos, &b)
	}

	BeginChildFrame(id, frame_bb.GetSize(), 0)
	return true
}

// only call EndListBox() if BeginListBox() returned true!
func EndListBox() {
	g := GImGui
	var window = g.CurrentWindow
	IM_ASSERT_USER_ERROR((window.Flags&ImGuiWindowFlags_ChildWindow) != 0, "Mismatched BeginListBox/EndListBox calls. Did you test the return value of BeginListBox?")

	EndChildFrame()
	EndGroup() // This is only required to be able to do IsItemXXX query on the whole ListBox including label
}

func ListBox(label string, current_item *int, items []string, items_count int, height_in_items int /*= -1*/) bool {
	var value_changed = ListBoxFunc(label, current_item, func(data any, idx int, out_text *string) bool {
		var items = data.([]string)
		if out_text != nil {
			*out_text = items[idx]
		}
		return true
	}, items, items_count, height_in_items)
	return value_changed
}

// This is merely a helper around BeginListBox(), EndListBox().
// Considering using those directly to submit custom data or store selection differently.
func ListBoxFunc(label string, current_item *int, items_getter func(data any, idx int, out_text *string) bool, data any, items_count int, height_in_items int /*= -1*/) bool {
	g := GImGui

	// Calculate size from "height_in_items"
	if height_in_items < 0 {
		height_in_items = ImMinInt(items_count, 7)
	}
	var height_in_items_f = float(height_in_items) + 0.25
	var size = ImVec2{0.0, ImFloor(GetTextLineHeightWithSpacing()*height_in_items_f + g.Style.FramePadding.y*2.0)}

	if !BeginListBox(label, size) {
		return false
	}

	// Assume all items have even height (= 1 line of text). If you need items of different height,
	// you can create a custom version of ListBox() in your code without using the clipper.
	var value_changed = false
	var clipper ImGuiListClipper
	clipper.Begin(items_count, GetTextLineHeightWithSpacing()) // We know exactly our line height here so we pass it as a minor optimization, but generally you don't need to.
	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			var item_text string
			if !items_getter(data, i, &item_text) {
				item_text = "*Unknown item*"
			}

			PushID(i)
			var item_selected = (i == *current_item)
			if (Selectable(item_text, item_selected, 0, ImVec2{})) {
				*current_item = i
				value_changed = true
			}
			if item_selected {
				SetItemDefaultFocus()
			}
			PopID()
		}
	}
	EndListBox()

	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}
