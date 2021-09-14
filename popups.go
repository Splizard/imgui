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

// Note that this is used for popups, which can overlap the non work-area of individual viewports.
func GetPopupAllowedExtentRect(*ImGuiWindow) ImRect {
	var g = GImGui
	var r_screen ImRect = GetMainViewport().GetMainRect()
	var padding ImVec2 = g.Style.DisplaySafeAreaPadding

	var x, y float
	if r_screen.GetWidth() > padding.x*2 {
		x = -padding.x
	}
	if r_screen.GetHeight() > padding.y*2 {
		y = -padding.y
	}

	r_screen.ExpandVec(ImVec2{x, y})
	return r_screen
}

func FindBestWindowPosForPopup(window *ImGuiWindow) ImVec2 {
	var g = GImGui

	var r_outer ImRect = GetPopupAllowedExtentRect(window)
	if window.Flags&ImGuiWindowFlags_ChildMenu != 0 {
		// Child menus typically request _any_ position within the parent menu item, and then we move the new menu outside the parent bounds.
		// This is how we end up with child menus appearing (most-commonly) on the right of the parent menu.
		IM_ASSERT(g.CurrentWindow == window)
		var parent_window *ImGuiWindow = g.CurrentWindowStack[len(g.CurrentWindowStack)-2].Window
		var horizontal_overlap float = g.Style.ItemInnerSpacing.x // We want some overlap to convey the relative depth of each menu (currently the amount of overlap is hard-coded to style.ItemSpacing.x).
		var r_avoid ImRect
		if parent_window.DC.MenuBarAppending {
			r_avoid = ImRect{ImVec2{-FLT_MAX, parent_window.ClipRect.Min.y}, ImVec2{FLT_MAX, parent_window.ClipRect.Max.y}} // Avoid parent menu-bar. If we wanted multi-line menu-bar, we may instead want to have the calling window setup e.g. a NextWindowData.PosConstraintAvoidRect field
		} else {
			r_avoid = ImRect{ImVec2{parent_window.Pos.x + horizontal_overlap, -FLT_MAX}, ImVec2{parent_window.Pos.x + parent_window.Size.x - horizontal_overlap - parent_window.ScrollbarSizes.x, FLT_MAX}}
		}
		return FindBestWindowPosForPopupEx(&window.Pos, &window.Size, &window.AutoPosLastDirection, &r_outer, &r_avoid, ImGuiPopupPositionPolicy_Default)
	}
	if window.Flags&ImGuiWindowFlags_Popup != 0 {
		var r_avoid ImRect = ImRect{ImVec2{window.Pos.x - 1, window.Pos.y - 1}, ImVec2{window.Pos.x + 1, window.Pos.y + 1}}
		return FindBestWindowPosForPopupEx(&window.Pos, &window.Size, &window.AutoPosLastDirection, &r_outer, &r_avoid, ImGuiPopupPositionPolicy_Default)
	}
	if window.Flags&ImGuiWindowFlags_Tooltip != 0 {
		// Position tooltip (always follows mouse)
		var sc float = g.Style.MouseCursorScale
		var ref_pos ImVec2 = NavCalcPreferredRefPos()
		var r_avoid ImRect
		if !g.NavDisableHighlight && g.NavDisableMouseHover && 0 == (g.IO.ConfigFlags&ImGuiConfigFlags_NavEnableSetMousePos) {
			r_avoid = ImRect{ImVec2{ref_pos.x - 16, ref_pos.y - 8}, ImVec2{ref_pos.x + 16, ref_pos.y + 8}}
		} else {
			r_avoid = ImRect{ImVec2{ref_pos.x - 16, ref_pos.y - 8}, ImVec2{ref_pos.x + 24*sc, ref_pos.y + 24*sc}} // FIXME: Hard-coded based on mouse cursor shape expectation. Exact dimension not very important.
		}
		return FindBestWindowPosForPopupEx(&ref_pos, &window.Size, &window.AutoPosLastDirection, &r_outer, &r_avoid, ImGuiPopupPositionPolicy_Tooltip)
	}
	IM_ASSERT(false)
	return window.Pos
}

// r_avoid = the rectangle to avoid (e.g. for tooltip it is a rectangle around the mouse cursor which we want to avoid. for popups it's a small point around the cursor.)
// r_outer = the visible area rectangle, minus safe area padding. If our popup size won't fit because of safe area padding we ignore it.
// (r_outer is usually equivalent to the viewport rectangle minus padding, but when multi-viewports are enabled and monitor
//  information are available, it may represent the entire platform monitor from the frame of reference of the current viewport.
//  this allows us to have tooltips/popups displayed out of the parent viewport.)
func FindBestWindowPosForPopupEx(ref_pos *ImVec2, size *ImVec2, last_dir *ImGuiDir, r_outer *ImRect, r_avoid *ImRect, policy ImGuiPopupPositionPolicy) ImVec2 {
	var base_pos_clamped ImVec2 = ImClampVec2(ref_pos, &r_outer.Min, r_outer.Max.Sub(*size))
	//GetForegroundDrawList().AddRect(r_avoid.Min, r_avoid.Max, IM_COL32(255,0,0,255));
	//GetForegroundDrawList().AddRect(r_outer.Min, r_outer.Max, IM_COL32(0,255,0,255));

	// Combo Box policy (we want a connecting edge)
	if policy == ImGuiPopupPositionPolicy_ComboBox {
		var dir_prefered_order = [ImGuiDir_COUNT]ImGuiDir{ImGuiDir_Down, ImGuiDir_Right, ImGuiDir_Left, ImGuiDir_Up}

		var start int
		if *last_dir != ImGuiDir_None {
			start = -1
		}

		for n := start; n < int(ImGuiDir_COUNT); n++ {
			var dir ImGuiDir
			if n == -1 {
				dir = *last_dir
			} else {
				dir = dir_prefered_order[n]
			}
			if n != -1 && dir == *last_dir { // Already tried this direction?
				continue
			}
			var pos ImVec2
			switch dir {
			case ImGuiDir_Down:
				pos = ImVec2{r_avoid.Min.x, r_avoid.Max.y} // Below, Toward Right (default)
			case ImGuiDir_Right:
				pos = ImVec2{r_avoid.Min.x, r_avoid.Min.y - size.y} // Above, Toward Right
			case ImGuiDir_Left:
				pos = ImVec2{r_avoid.Max.x - size.x, r_avoid.Max.y} // Below, Toward Left
			case ImGuiDir_Up:
				pos = ImVec2{r_avoid.Max.x - size.x, r_avoid.Min.y - size.y} // Above, Toward Left
			}
			if !r_outer.ContainsRect(ImRect{pos, pos.Add(*size)}) {
				continue
			}
			*last_dir = dir
			return pos
		}
	}

	// Tooltip and Default popup policy
	// (Always first try the direction we used on the last frame, if any)
	if policy == ImGuiPopupPositionPolicy_Tooltip || policy == ImGuiPopupPositionPolicy_Default {
		var dir_prefered_order = [ImGuiDir_COUNT]ImGuiDir{ImGuiDir_Right, ImGuiDir_Down, ImGuiDir_Up, ImGuiDir_Left}

		var start int
		if *last_dir != ImGuiDir_None {
			start = -1
		}

		for n := start; n < int(ImGuiDir_COUNT); n++ {
			var dir ImGuiDir
			if n == -1 {
				dir = *last_dir
			} else {
				dir = dir_prefered_order[n]
			}
			if n != -1 && dir == *last_dir { // Already tried this direction?
				continue
			}

			var wl, wr float
			if dir == ImGuiDir_Left {
				wl = r_avoid.Min.x
			} else {
				wl = r_outer.Max.x
			}
			if dir == ImGuiDir_Right {
				wr = r_avoid.Max.x
			} else {
				wr = r_outer.Min.x
			}

			var avail_w float = wl - wr

			var hl, hr float
			if dir == ImGuiDir_Up {
				hl = r_avoid.Min.y
			} else {
				hl = r_outer.Max.y
			}
			if dir == ImGuiDir_Down {
				hr = r_avoid.Max.y
			} else {
				hr = r_outer.Min.y
			}

			var avail_h float = hl - hr

			// If there not enough room on one axis, there's no point in positioning on a side on this axis (e.g. when not enough width, use a top/bottom position to maximize available width)
			if avail_w < size.x && (dir == ImGuiDir_Left || dir == ImGuiDir_Right) {
				continue
			}
			if avail_h < size.y && (dir == ImGuiDir_Up || dir == ImGuiDir_Down) {
				continue
			}

			var pos ImVec2
			if dir == ImGuiDir_Left {
				pos.x = r_avoid.Min.x - size.x
			} else if dir == ImGuiDir_Right {
				pos.x = r_avoid.Max.x
			} else {
				pos.x = base_pos_clamped.x
			}
			if dir == ImGuiDir_Up {
				pos.y = r_avoid.Min.y - size.y
			} else if dir == ImGuiDir_Down {
				pos.y = r_avoid.Max.y
			} else {
				pos.y = base_pos_clamped.y
			}

			// Clamp top-left corner of popup
			pos.x = ImMax(pos.x, r_outer.Min.x)
			pos.y = ImMax(pos.y, r_outer.Min.y)

			*last_dir = dir
			return pos
		}
	}

	// Fallback when not enough room:
	*last_dir = ImGuiDir_None

	// For tooltip we prefer avoiding the cursor at all cost even if it means that part of the tooltip won't be visible.
	if policy == ImGuiPopupPositionPolicy_Tooltip {
		return ref_pos.Add(ImVec2{2, 2})
	}

	// Otherwise try to keep within display
	var pos ImVec2 = *ref_pos
	pos.x = ImMax(ImMin(pos.x+size.x, r_outer.Max.x)-size.x, r_outer.Min.x)
	pos.y = ImMax(ImMin(pos.y+size.y, r_outer.Max.y)-size.y, r_outer.Min.y)
	return pos
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
