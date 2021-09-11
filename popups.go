package imgui

func GetTopMostPopupModal() *ImGuiWindow {
	var g = GImGui
	for n := len(g.OpenPopupStack) - 1; n >= 0; n-- {
		if popup := g.OpenPopupStack[n].Window; popup != nil {
			if popup.Flags&ImGuiWindowFlags_Modal != 0 {
				return popup
			}
		}
	}
	return nil
}

// When popups are stacked, clicking on a lower level popups puts focus back to it and close popups above it.
// This function closes any popups that are over 'ref_window'.
func ClosePopupsOverWindow(ref_window *ImGuiWindow, restore_focus_to_window_under_popup bool) {
	var g = GImGui
	if len(g.OpenPopupStack) == 0 {
		return
	}

	// Don't close our own child popup windows.
	var popup_count_to_keep int = 0
	if ref_window != nil {
		// Find the highest popup which is a descendant of the reference window (generally reference window = NavWindow)
		for ; popup_count_to_keep < int(len(g.OpenPopupStack)); popup_count_to_keep++ {
			var popup *ImGuiPopupData = &g.OpenPopupStack[popup_count_to_keep]
			if popup.Window == nil {
				continue
			}
			IM_ASSERT((popup.Window.Flags & ImGuiWindowFlags_Popup) != 0)
			if popup.Window.Flags&ImGuiWindowFlags_ChildWindow != 0 {
				continue
			}

			// Trim the stack unless the popup is a direct parent of the reference window (the reference window is often the NavWindow)
			// - With this stack of window, clicking/focusing Popup1 will close Popup2 and Popup3:
			//     Window . Popup1 . Popup2 . Popup3
			// - Each popups may contain child windows, which is why we compare .RootWindow!
			//     Window . Popup1 . Popup1_Child . Popup2 . Popup2_Child
			var ref_window_is_descendent_of_popup bool = false
			for n := popup_count_to_keep; n < int(len(g.OpenPopupStack)); n++ {
				if popup_window := g.OpenPopupStack[n].Window; popup_window != nil {
					if popup_window.RootWindow == ref_window.RootWindow {
						ref_window_is_descendent_of_popup = true
						break
					}
				}
			}
			if !ref_window_is_descendent_of_popup {
				break
			}
		}
	}
	if popup_count_to_keep < int(len(g.OpenPopupStack)) { // This test is not required but it allows to set a convenient breakpoint on the statement below
		//IMGUI_DEBUG_LOG_POPUP("ClosePopupsOverWindow(\"%s\") . ClosePopupToLevel(%d)\n", ref_window.Name, popup_count_to_keep)
		ClosePopupToLevel(popup_count_to_keep, restore_focus_to_window_under_popup)
	}
}

// Supported flags: ImGuiPopupFlags_AnyPopupId, ImGuiPopupFlags_AnyPopupLevel
func IsPopupOpenID(id ImGuiID, popup_flags ImGuiPopupFlags) bool {
	var g = GImGui
	if popup_flags&ImGuiPopupFlags_AnyPopupId != 0 {
		// Return true if any popup is open at the current BeginPopup() level of the popup stack
		// This may be used to e.g. test for another popups already opened to handle popups priorities at the same level.
		IM_ASSERT(id == 0)
		if popup_flags&ImGuiPopupFlags_AnyPopupLevel != 0 {
			return len(g.OpenPopupStack) > 0
		} else {
			return len(g.OpenPopupStack) > len(g.BeginPopupStack)
		}
	} else {
		if popup_flags&ImGuiPopupFlags_AnyPopupLevel != 0 {
			// Return true if the popup is open anywhere in the popup stack
			for n := range g.OpenPopupStack {
				if g.OpenPopupStack[n].PopupId == id {
					return true
				}
			}

			return false
		} else {
			// Return true if the popup is open at the current BeginPopup() level of the popup stack (this is the most-common query)
			return len(g.OpenPopupStack) > len(g.BeginPopupStack) && g.OpenPopupStack[len(g.BeginPopupStack)].PopupId == id
		}
	}
}
