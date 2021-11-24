package imgui

// Widgets: List Boxes
// - This is essentially a thin wrapper to using BeginChild/EndChild with some stylistic changes.
// - The BeginListBox()/EndListBox() api allows you to manage your contents and selection state however you want it, by creating e.g. Selectable() or any items.
// - The simplified/old ListBox() api are helpers over BeginListBox()/EndListBox() which are kept available for convenience purpose. This is analoguous to how Combos are created.
// - Choose frame width:   size.x > 0.0: custom  /  size.x < 0.0 or -FLT_MIN: right-align   /  size.x = 0.0 (default): use current ItemWidth
// - Choose frame height:  size.y > 0.0: custom  /  size.y < 0.0 or -FLT_MIN: bottom-align  /  size.y = 0.0 (default): arbitrary default height which can fit ~7 items
func BeginListBox(label string, size ImVec2) bool { panic("not implemented") } // open a framed scrolling region
func EndListBox()                                 { panic("not implemented") } // only call EndListBox() if BeginListBox() returned true!

func ListBox(label string, current_item *int, items []string, items_count int, height_in_items int /*= -1*/) bool {
	var value_changed = ListBoxFunc(label, current_item, func(data interface{}, idx int, out_text *string) bool {
		var items = data.([]string)
		if out_text != nil {
			*out_text = items[idx]
		}
		return true
	}, items, items_count, height_in_items)
	return value_changed
}

func ListBoxFunc(label string, current_item *int, items_getter func(data interface{}, idx int, out_text *string) bool, data interface{}, items_count int, height_in_items int /*= -1*/) bool {
	panic("not implemented")
}
