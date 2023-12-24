package imgui

// Item/Widgets Utilities and Query Functions
// - Most of the functions are referring to the previous Item that has been submitted.
// - See Demo Window under "Widgets->Querying Status" for an interactive visualization of most of those functions.

// IsItemHovered is the last item hovered? (and usable, aka not blocked by a popup, etc.). See ImGuiHoveredFlags for more options.
// This is roughly matching the behavior of internal-facing ItemHoverable()
// - we allow hovering to be true when ActiveId==window.MoveID, so that clicking on non-interactive items such as a Text() item still returns true with IsItemHovered()
// - this should work even for non-interactive items that have no ID, so we cannot use LastItemId
func IsItemHovered(flags ImGuiHoveredFlags) bool {
	window := g.CurrentWindow
	if g.NavDisableMouseHover && !g.NavDisableHighlight {
		if (g.LastItemData.InFlags&ImGuiItemFlags_Disabled != 0) && (flags&ImGuiHoveredFlags_AllowWhenDisabled == 0) {
			return false
		}
		return IsItemFocused()
	}

	// Test for bounding box overlap, as updated as ItemAdd()
	var status_flags = g.LastItemData.StatusFlags
	if status_flags&ImGuiItemStatusFlags_HoveredRect == 0 {
		return false
	}
	IM_ASSERT((flags & (ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows)) == 0) // Flags not supported by this function

	// Test if we are hovering the right window (our window could be behind another window)
	// [2021/03/02] Reworked / reverted the revert, finally. Note we want e.g. BeginGroup/ItemAdd/EndGroup to work as well. (#3851)
	// [2017/10/16] Reverted commit 344d48be3 and testing RootWindow instead. I believe it is correct to NOT test for RootWindow but this leaves us unable
	// to use IsItemHovered() after EndChild() itself. Until a solution is found I believe reverting to the test from 2017/09/27 is safe since this was
	// the test that has been running for a long while.
	if g.HoveredWindow != window && (status_flags&ImGuiItemStatusFlags_HoveredWindow) == 0 {
		if (flags & ImGuiHoveredFlags_AllowWhenOverlapped) == 0 {
			return false
		}
	}

	// Test if another item is active (e.g. being dragged)
	if (flags & ImGuiHoveredFlags_AllowWhenBlockedByActiveItem) == 0 {
		if g.ActiveId != 0 && g.ActiveId != g.LastItemData.ID && !g.ActiveIdAllowOverlap && g.ActiveId != window.MoveId {
			return false
		}
	}

	// Test if interactions on this window are blocked by an active popup or modal.
	// The ImGuiHoveredFlags_AllowWhenBlockedByPopup flag will be tested here.
	if !IsWindowContentHoverable(window, flags) {
		return false
	}

	// Test if the item is disabled
	if (g.LastItemData.InFlags&ImGuiItemFlags_Disabled != 0) && (flags&ImGuiHoveredFlags_AllowWhenDisabled == 0) {
		return false
	}

	// Special handling for calling after Begin() which represent the title bar or tab.
	// When the window is collapsed (SkipItems==true) that last item will never be overwritten so we need to detect the case.
	return g.LastItemData.ID == window.MoveId && window.WriteAccessed
}

// IsItemFocused is the last item focused for keyboard/gamepad navigation?
// == GetItemID() == GetFocusID()
func IsItemFocused() bool {
	return !(g.NavId != g.LastItemData.ID || g.NavId == 0)
}

// IsItemClicked is the last item hovered and mouse clicked on? (**)  == IsMouseClicked(mouse_button) && IsItemHovered()Important. (**) this it NOT equivalent to the behavior of e.g. Button(). Read comments in function definition.
// Important: this can be useful but it is NOT equivalent to the behavior of e.g.Button()!
// Most widgets have specific reactions based on mouse-up/down state, mouse position etc.
func IsItemClicked(mouse_button ImGuiMouseButton) bool {
	return IsMouseClicked(mouse_button, false) && IsItemHovered(ImGuiHoveredFlags_None)
}

// IsItemVisible is the last item visible? (items may be out of sight because of clipping/scrolling)
func IsItemVisible() bool {
	return g.CurrentWindow.ClipRect.Overlaps(g.LastItemData.Rect)
}

// IsItemEdited did the last item modify its underlying value this frame? or was pressed? This is generally the same as the "bool" return value of many widgets.
func IsItemEdited() bool {
	return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_Edited) != 0
}

// IsItemActivated was the last item just made active (item was previously inactive).
func IsItemActivated() bool {
	if g.ActiveId != 0 {
		return g.ActiveId == g.LastItemData.ID && g.ActiveIdPreviousFrame != g.LastItemData.ID
	}
	return false
}

// IsItemDeactivated was the last item just made inactive (item was previously active). Useful for Undo/Redo patterns with widgets that requires continuous editing.
func IsItemDeactivated() bool {
	if g.LastItemData.StatusFlags&ImGuiItemStatusFlags_HasDeactivated != 0 {
		return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_Deactivated) != 0
	}
	return g.ActiveIdPreviousFrame == g.LastItemData.ID && g.ActiveIdPreviousFrame != 0 && g.ActiveId != g.LastItemData.ID
}

// IsItemDeactivatedAfterEdit was the last item just made inactive and made a value change when it was active? (e.g. Slider/Drag moved). Useful for Undo/Redo patterns with widgets that requires continuous editing. Note that you may get false positives (some widgets such as Combo()/ListBox()/Selectable() will return true even when clicking an already selected item).
func IsItemDeactivatedAfterEdit() bool {
	return IsItemDeactivated() && (g.ActiveIdPreviousFrameHasBeenEditedBefore || (g.ActiveId == 0 && g.ActiveIdHasBeenEditedBefore))
}

// IsItemToggledOpen was the last item open state toggled? set by TreeNode().
func IsItemToggledOpen() bool {
	return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_ToggledOpen) != 0
}

// IsAnyItemHovered is any item hovered?
func IsAnyItemHovered() bool {
	return g.HoveredId != 0 || g.HoveredIdPreviousFrame != 0
}

// IsAnyItemActive is any item active?
func IsAnyItemActive() bool {
	return g.ActiveId != 0
}

// IsAnyItemFocused is any item focused?
func IsAnyItemFocused() bool {
	return g.NavId != 0 && !g.NavDisableHighlight
}

// GetItemRectMin get upper-left bounding rectangle of the last item (screen space)
func GetItemRectMin() ImVec2 {
	return g.LastItemData.Rect.Min
}

// GetItemRectMax get lower-right bounding rectangle of the last item (screen space)
func GetItemRectMax() ImVec2 {
	return g.LastItemData.Rect.Max
}

// GetItemRectSize get size of last item
func GetItemRectSize() ImVec2 {
	return g.LastItemData.Rect.GetSize()
}

// SetItemAllowOverlap allow last item to be overlapped by a subsequent item. sometimes useful with invisible buttons, selectables, etc. to catch unused area.
// Allow last item to be overlapped by a subsequent item. Both may be activated during the same frame before the later one takes priority.
// FIXME: Although this is exposed, its interaction and ideal idiom with using ImGuiButtonFlags_AllowItemOverlap flag are extremely confusing, need rework.
func SetItemAllowOverlap() {
	var id = g.LastItemData.ID
	if g.HoveredId == id {
		g.HoveredIdAllowOverlap = true
	}
	if g.ActiveId == id {
		g.ActiveIdAllowOverlap = true
	}
}
