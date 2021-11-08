package imgui

// Basic Helpers for widget code

// Remotely activate a button, checkbox, tree node etc. given its unique ID. activation is queued and processed on the next frame when the item is encountered again.
func ActivateItem(id ImGuiID) {
	var g = GImGui
	g.NavNextActivateId = id
}

// Called by ItemAdd()
// Process TAB/Shift+TAB. Be mindful that this function may _clear_ the ActiveID when tabbing out.
// [WIP] This will eventually be refactored and moved into NavProcessItem()
func ItemInputable(window *ImGuiWindow, id ImGuiID) {
	var g = GImGui
	IM_ASSERT(id != 0 && id == g.LastItemData.ID)

	// Increment counters
	// FIXME: ImGuiItemFlags_Disabled should disable more.
	var is_tab_stop bool = (g.LastItemData.InFlags & (ImGuiItemFlags_NoTabStop | ImGuiItemFlags_Disabled)) == 0
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

	var g = GImGui
	var window = g.CurrentWindow
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

//Was the last item selection toggled? (after Selectable(), TreeNode() etc. We only returns toggle _event_ in order to handle clipping correctly)
func IsItemToggledSelection() bool {
	var g = GImGui
	return (g.LastItemData.StatusFlags & ImGuiItemStatusFlags_ToggledSelection) != 0
}

func ShrinkWidths(s *ImGuiShrinkWidthItem, count int, width_excess float) {
	panic("not implemented")
}

// Inputs
// FIXME: Eventually we should aim to move e.g. IsActiveIdUsingKey() into IsKeyXXX functions.
func SetItemUsingMouseWheel() {
	var g = GImGui
	var id ImGuiID = g.LastItemData.ID
	if g.HoveredId == id {
		g.HoveredIdUsingMouseWheel = true
	}
	if g.ActiveId == id {
		g.ActiveIdUsingMouseWheel = true
	}
}
