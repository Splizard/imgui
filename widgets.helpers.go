package imgui

import (
	"sort"

	"github.com/Splizard/imgui/golang"
)

// Remotely activate a button, checkbox, tree node etc. given its unique ID. activation is queued and processed on the next frame when the item is encountered again.
func ActivateItem(id ImGuiID) {
	g := GImGui
	g.NavNextActivateId = id
}

// Called by ItemAdd()
// Process TAB/Shift+TAB. Be mindful that this function may _clear_ the ActiveID when tabbing out.
// [WIP] This will eventually be refactored and moved into NavProcessItem()
func ItemInputable(window *ImGuiWindow, id ImGuiID) {
	g := GImGui
	IM_ASSERT(id != 0 && id == g.LastItemData.ID)

	// Increment counters
	// FIXME: ImGuiItemFlags_Disabled should disable more.
	var is_tab_stop = (g.LastItemData.InFlags & (ImGuiItemFlags_NoTabStop | ImGuiItemFlags_Disabled)) == 0
	window.DC.FocusCounterRegular++
	if is_tab_stop {
		window.DC.FocusCounterTabStop++
		if g.NavId == id {
			g.NavIdTabCounter = window.DC.FocusCounterTabStop
		}
	}

	// Process TAB/Shift-TAB to tab *OUT* of the currently focused item.
	// (Note that we can always TAB out of a widget that doesn't allow tabbing in)
	if g.ActiveId == id && g.TabFocusPressed && !IsActiveIdUsingKey(ImGuiKey_Tab) && g.TabFocusRequestNextWindow == nil {
		g.TabFocusRequestNextWindow = window

		var add int
		if g.IO.KeyShift {
			if is_tab_stop {
				add = -1
			}
		} else {
			add = +1
		}

		g.TabFocusRequestNextCounterTabStop = window.DC.FocusCounterTabStop + add // Modulo on index will be applied at the end of frame once we've got the total counter of items.
	}

	// Handle focus requests
	if g.TabFocusRequestCurrWindow == window {
		if window.DC.FocusCounterRegular == g.TabFocusRequestCurrCounterRegular {
			g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_FocusedByCode
			return
		}
		if is_tab_stop && window.DC.FocusCounterTabStop == g.TabFocusRequestCurrCounterTabStop {
			g.NavJustTabbedId = id
			g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_FocusedByTabbing
			return
		}

		// If another item is about to be focused, we clear our own active id
		if g.ActiveId == id {
			ClearActiveID()
		}
	}
}

func CalcWrapWidthForPos(pos *ImVec2, wrap_pos_x float) float {
	if wrap_pos_x < 0.0 {
		return 0.0
	}

	g := GImGui
	window := g.CurrentWindow
	if wrap_pos_x == 0.0 {
		// We could decide to setup a default wrapping max point for auto-resizing windows,
		// or have auto-wrap (with unspecified wrapping pos) behave as a ContentSize extending function?
		//if (window.Hidden && (window.Flags & ImGuiWindowFlags_AlwaysAutoResize))
		//    wrap_pos_x = ImMax(window.WorkRect.Min.x + g.FontSize * 10.0f, window.WorkRect.Max.x);
		//else
		wrap_pos_x = window.WorkRect.Max.x
	} else if wrap_pos_x > 0.0 {
		wrap_pos_x += window.Pos.x - window.Scroll.x // wrap_pos_x is provided is window local space
	}

	return ImMax(wrap_pos_x-pos.x, 1.0)
}

// Was the last item selection toggled? (after Selectable(), TreeNode() etc. We only returns toggle _event_ in order to handle clipping correctly)
func IsItemToggledSelection() bool {
	g := GImGui
	return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_ToggledSelection) != 0
}

// Shrink excess width from a set of item, by removing width from the larger items first.
// Set items Width to -1.0f to disable shrinking this item.
func ShrinkWidths(items []ImGuiShrinkWidthItem, count int, width_excess float) {
	if count == 1 {
		if items[0].Width >= 0.0 {
			items[0].Width = ImMax(items[0].Width-width_excess, 1.0)
		}
		return
	}
	sort.Slice(items, func(i, j golang.Int) bool {
		var a = items[i]
		var b = items[j]
		if d := (int)(b.Width - a.Width); d != 0 {
			return true
		}
		return (b.Index - a.Index) != 0
	})
	var count_same_width int = 1
	for width_excess > 0.0 && count_same_width < count {
		for count_same_width < count && items[0].Width <= items[count_same_width].Width {
			count_same_width++
		}
		var max_width_to_remove_per_item = (items[0].Width - 1.0)
		if count_same_width < count && items[count_same_width].Width >= 0.0 {
			max_width_to_remove_per_item = (items[0].Width - items[count_same_width].Width)
		}
		if max_width_to_remove_per_item <= 0.0 {
			break
		}
		var width_to_remove_per_item = ImMin(width_excess/float(count_same_width), max_width_to_remove_per_item)
		for item_n := int(0); item_n < count_same_width; item_n++ {
			items[item_n].Width -= width_to_remove_per_item
		}
		width_excess -= width_to_remove_per_item * float(count_same_width)
	}

	// Round width and redistribute remainder left-to-right (could make it an option of the function?)
	// Ensure that e.g. the right-most tab of a shrunk tab-bar always reaches exactly at the same distance from the right-most edge of the tab bar separator.
	width_excess = 0.0
	for n := int(0); n < count; n++ {
		var width_rounded = ImFloor(items[n].Width)
		width_excess += items[n].Width - width_rounded
		items[n].Width = width_rounded
	}
	if width_excess > 0.0 {
		for n := int(0); n < count; n++ {
			if items[n].Index < (int)(width_excess+0.01) {
				items[n].Width += 1.0
			}
		}
	}
}

// Inputs
// FIXME: Eventually we should aim to move e.g. IsActiveIdUsingKey() into IsKeyXXX functions.
func SetItemUsingMouseWheel() {
	g := GImGui
	var id = g.LastItemData.ID
	if g.HoveredId == id {
		g.HoveredIdUsingMouseWheel = true
	}
	if g.ActiveId == id {
		g.ActiveIdUsingMouseWheel = true
	}
}
