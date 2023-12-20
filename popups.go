package imgui

import "fmt"

// Popups, Modals
//  - They block normal mouse hovering detection (and therefore most mouse interactions) behind them.
//  - If not modal: they can be closed by clicking anywhere outside them, or by pressing ESCAPE.
//  - Their visibility state (~bool) is held internally instead of being held by the programmer as we are used to with regular Begin*() calls.
//  - The 3 properties above are related: we need to retain popup visibility state in the library because popups may be closed as any time.
//  - You can bypass the hovering restriction by using ImGuiHoveredFlags_AllowWhenBlockedByPopup when calling IsItemHovered() or IsWindowHovered().
//  - IMPORTANT: Popup identifiers are relative to the current ID stack, so OpenPopup and BeginPopup generally needs to be at the same level of the stack.
//    This is sometimes leading to confusing mistakes. May rework this in the future.

// BeginPopup Popups: begin/end functions
//   - BeginPopup(): query popup state, if open start appending into the window. Call EndPopup() afterwards. ImGuiWindowFlags are forwarded to the window.
//   - BeginPopupModal(): block every interactions behind the window, cannot be closed by user, add a dimming background, has a title bar.
//
// return true if the popup is open, and you can start outputting to it.
func BeginPopup(str_id string, flags ImGuiWindowFlags) bool {
	var g = GImGui
	if len(g.OpenPopupStack) <= len(g.BeginPopupStack) {
		g.NextWindowData.ClearFlags() // We behave like Begin() and need to consume those values
		return false
	}
	flags |= ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoSavedSettings
	return BeginPopupEx(g.CurrentWindow.GetIDs(str_id), flags)
}

// BeginPopupModal If 'p_open' is specified for a modal popup window, the popup will have a regular close button which will close the popup.
// Note that popup visibility status is owned by Dear ImGui (and manipulated with e.g. OpenPopup) so the actual value of *p_open is meaningless here.
// return true if the modal is open, and you can start outputting to it.
func BeginPopupModal(name string, p_open *bool, flags ImGuiWindowFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	var id ImGuiID = window.GetIDs(name)
	if !isPopupOpen(id, ImGuiPopupFlags_None) {
		g.NextWindowData.ClearFlags() // We behave like Begin() and need to consume those values
		return false
	}

	// Center modal windows by default for increased visibility
	// (this won't really last as settings will kick in, and is mostly for backward compatibility. user may do the same themselves)
	// FIXME: Should test for (PosCond & window.SetWindowPosAllowFlags) with the upcoming window.
	if (g.NextWindowData.Flags & ImGuiNextWindowDataFlags_HasPos) == 0 {
		var viewport *ImGuiViewport = GetMainViewport()
		center := viewport.GetCenter()
		SetNextWindowPos(&center, ImGuiCond_FirstUseEver, ImVec2{0.5, 0.5})
	}

	flags |= ImGuiWindowFlags_Popup | ImGuiWindowFlags_Modal | ImGuiWindowFlags_NoCollapse
	var is_open = Begin(name, p_open, flags)
	if !is_open || (p_open != nil && !*p_open) { // NB: is_open can be 'false' when the popup is completely clipped (e.g. zero size display)
		EndPopup()
		if is_open {
			ClosePopupToLevel(int(len(g.BeginPopupStack)), true)
		}
		return false
	}
	return is_open
}

// EndPopup only call EndPopup() if BeginPopupXXX() returns true!
func EndPopup() {
	var g = GImGui
	var window = g.CurrentWindow
	IM_ASSERT(window.Flags&ImGuiWindowFlags_Popup != 0) // Mismatched BeginPopup()/EndPopup() calls
	IM_ASSERT(len(g.BeginPopupStack) > 0)

	// Make all menus and popups wrap around for now, may need to expose that policy (e.g. focus scope could include wrap/loop policy flags used by new move requests)
	if g.NavWindow == window {
		NavMoveRequestTryWrapping(window, ImGuiNavMoveFlags_LoopY)
	}

	// Child-popups don't need to be laid out
	IM_ASSERT(!g.WithinEndChild)
	if window.Flags&ImGuiWindowFlags_ChildWindow != 0 {
		g.WithinEndChild = true
	}
	End()
	g.WithinEndChild = false
}

// OpenPopup Popups: open/close functions
//   - OpenPopup(): set popup state to open. are ImGuiPopupFlags available for opening options.
//   - If not modal: they can be closed by clicking anywhere outside them, or by pressing ESCAPE.
//   - CloseCurrentPopup(): use inside the BeginPopup()/EndPopup() scope to close manually.
//   - CloseCurrentPopup() is called by default by Selectable()/MenuItem() when activated (FIXME: need some options).
//   - Use ImGuiPopupFlags_NoOpenOverExistingPopup to a opening a popup if there's already one at the same level. This is equivalent to e.g. testing for !IsAnyPopupOpen() prior to OpenPopup().
//   - Use IsWindowAppearing() after BeginPopup() to tell if a window just opened.
//
// call to mark popup as open (don't call every frame!).
func OpenPopup(str_id string, popup_flags ImGuiPopupFlags) {
	var g = GImGui
	OpenPopupEx(g.CurrentWindow.GetIDs(str_id), popup_flags)
}

// OpenPopupID id overload to facilitate calling from nested stacks
func OpenPopupID(id ImGuiID, popup_flags ImGuiPopupFlags) {
	OpenPopupEx(id, popup_flags)
}

// OpenPopupOnItemClick Helper to open a popup if mouse button is released over the item
// - This is essentially the same as BeginPopupContextItem() but without the trailing BeginPopup()
// helper to open popup when clicked on last item. Default to ImGuiPopupFlags_MouseButtonRight == 1. (note: actually triggers on the mouse _released_ event to be consistent with popup behaviors)
func OpenPopupOnItemClick(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) {
	var g = GImGui
	var window = g.CurrentWindow
	var mouse_button = ImGuiMouseButton(popup_flags & ImGuiPopupFlags_MouseButtonMask_)
	if IsMouseReleased(mouse_button) && IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup) {
		var id ImGuiID = g.LastItemData.ID // If user hasn't passed an ID, we can use the LastItemID. Using LastItemID as a Popup ID won't conflict!
		if str_id != "" {
			id = window.GetIDs(str_id)
		}
		IM_ASSERT(id != 0) // You cannot pass a nil str_id if the last item has no identifier (e.g. a Text() item)
		OpenPopupEx(id, popup_flags)
	}
}

// CloseCurrentPopup manually close the popup we have begin-ed into.
func CloseCurrentPopup() {
	var g = GImGui
	var popup_idx = int(len(g.BeginPopupStack)) - 1
	if popup_idx < 0 || popup_idx >= int(len(g.OpenPopupStack)) || g.BeginPopupStack[popup_idx].PopupId != g.OpenPopupStack[popup_idx].PopupId {
		return
	}

	// Closing a menu closes its top-most parent popup (unless a modal)
	for popup_idx > 0 {
		var popup_window = g.OpenPopupStack[popup_idx].Window
		var parent_popup_window = g.OpenPopupStack[popup_idx-1].Window
		var close_parent = false
		if popup_window != nil && (popup_window.Flags&ImGuiWindowFlags_ChildMenu != 0) {
			if parent_popup_window == nil || (parent_popup_window.Flags&ImGuiWindowFlags_Modal == 0) {
				close_parent = true
			}
		}
		if !close_parent {
			break
		}
		popup_idx--
	}
	//IMGUI_DEBUG_LOG_POPUP("CloseCurrentPopup %d . %d\n", g.BeginPopupStack.Size-1, popup_idx)
	ClosePopupToLevel(popup_idx, true)

	// A common pattern is to close a popup when selecting a menu item/selectable that will open another window.
	// To improve this usage pattern, we avoid nav highlight for a single frame in the parent window.
	// Similarly, we could avoid mouse hover highlight in this window but it is less visually problematic.
	if window := g.NavWindow; window != nil {
		window.DC.NavHideHighlightOneFrame = true
	}
}

// Popups: open+begin combined functions helpers
//  - Helpers to do OpenPopup+BeginPopup where the Open action is triggered by e.g. hovering an item and right-clicking.
//  - They are convenient to easily create context menus, hence the name.
//  - IMPORTANT: Notice that BeginPopupContextXXX takes just ImGuiPopupFlags like OpenPopup() and unlike BeginPopup(). For full consistency, we may add ImGuiWindowFlags to the BeginPopupContextXXX functions in the future.
//  - IMPORTANT: we exceptionally default their flags to 1 (== ImGuiPopupFlags_MouseButtonRight) for backward compatibility with older API taking 'mouse_button int/*= r*/,so if you add other flags remember to re-add the ImGuiPopupFlags_MouseButtonRight.

// BeginPopupContextItem This is a helper to handle the simplest case of associating one named popup to one given widget.
// - To create a popup associated to the last item, you generally want to pass a nil value to str_id.
// - To create a popup with a specific identifier, pass it in str_id.
//   - This is useful when using using BeginPopupContextItem() on an item which doesn't have an identifier, e.g. a Text() call.
//   - This is useful when multiple code locations may want to manipulate/open the same popup, given an explicit id.
//   - You may want to handle the whole on user side if you have specific needs (e.g. tweaking IsItemHovered() parameters).
//     This is essentially the same as:
//     id = str_id ? GetID(str_id) : GetItemID();
//     OpenPopupOnItemClick(str_id);
//     return BeginPopup(id);
//     Which is essentially the same as:
//     id = str_id ? GetID(str_id) : GetItemID();
//     if (IsItemHovered() && IsMouseReleased(ImGuiMouseButton_Right))
//     OpenPopup(id);
//     return BeginPopup(id);
//     The main difference being that this is tweaked to avoid computing the ID twice.
//
// open+begin popup when clicked on last item. Use str_id==NULL to associate the popup to previous item. If you want to use that on a non-interactive item such as Text() you need to pass in an explicit ID here. read comments in .cpp!
func BeginPopupContextItem(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}
	var id ImGuiID = g.LastItemData.ID // If user hasn't passed an ID, we can use the LastItemID. Using LastItemID as a Popup ID won't conflict!
	if str_id != "" {
		id = window.GetIDs(str_id)
	}
	IM_ASSERT(id != 0) // You cannot pass a nil str_id if the last item has no identifier (e.g. a Text() item)
	var mouse_button = ImGuiMouseButton(popup_flags & ImGuiPopupFlags_MouseButtonMask_)
	if IsMouseReleased(mouse_button) && IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup) {
		OpenPopupEx(id, popup_flags)
	}
	return BeginPopupEx(id, ImGuiWindowFlags_AlwaysAutoResize|ImGuiWindowFlags_NoTitleBar|ImGuiWindowFlags_NoSavedSettings)
}

func BeginPopupContextVoid(str_id string, popup_flags ImGuiPopupFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if str_id == "" {
		str_id = "void_context"
	}
	var id ImGuiID = window.GetIDs(str_id)
	var mouse_button = ImGuiMouseButton(popup_flags & ImGuiPopupFlags_MouseButtonMask_)
	if IsMouseReleased(mouse_button) && !IsWindowHovered(ImGuiHoveredFlags_AnyWindow) {
		if GetTopMostPopupModal() == nil {
			OpenPopupEx(id, popup_flags)
		}
	}
	return BeginPopupEx(id, ImGuiWindowFlags_AlwaysAutoResize|ImGuiWindowFlags_NoTitleBar|ImGuiWindowFlags_NoSavedSettings)
}

// BeginPopupContextWindow open+begin popup when clicked on current window.
func BeginPopupContextWindow(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) bool {
	panic("not implemented")
}

// BeginPopupContext open+begin popup when clicked in  (where there are no windows).
func BeginPopupContext(str_id string /*= L*/, popup_flags ImGuiPopupFlags /*= 1*/) bool {
	panic("not implemented")
}

// IsPopupOpen Popups: query functions
//   - IsPopupOpen(): return true if the popup is open at the current BeginPopup() level of the popup stack.
//   - IsPopupOpen() with ImGuiPopupFlags_AnyPopupId: return true if any popup is open at the current BeginPopup() level of the popup stack.
//   - IsPopupOpen() with ImGuiPopupFlags_AnyPopupId + ImGuiPopupFlags_AnyPopupLevel: return true if any popup is open.
//
// return true if the popup is open.
func IsPopupOpen(str_id string, flags ImGuiPopupFlags) bool {
	var g = GImGui
	var id ImGuiID
	if flags&ImGuiPopupFlags_AnyPopupId == 0 {
		id = g.CurrentWindow.GetIDs(str_id)
	}
	if (flags&ImGuiPopupFlags_AnyPopupLevel != 0) && id != 0 {
		IM_ASSERT_USER_ERROR(false, "Cannot use IsPopupOpen() with a string id and ImGuiPopupFlags_AnyPopupLevel.") // But non-string version is legal and used internally
	}
	return IsPopupOpenID(id, flags)
}

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

// GetPopupAllowedExtentRect Note that this is used for popups, which can overlap the non work-area of individual viewports.
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
		if !g.NavDisableHighlight && g.NavDisableMouseHover && g.IO.ConfigFlags&ImGuiConfigFlags_NavEnableSetMousePos == 0 {
			r_avoid = ImRect{ImVec2{ref_pos.x - 16, ref_pos.y - 8}, ImVec2{ref_pos.x + 16, ref_pos.y + 8}}
		} else {
			r_avoid = ImRect{ImVec2{ref_pos.x - 16, ref_pos.y - 8}, ImVec2{ref_pos.x + 24*sc, ref_pos.y + 24*sc}} // FIXME: Hard-coded based on mouse cursor shape expectation. Exact dimension not very important.
		}
		return FindBestWindowPosForPopupEx(&ref_pos, &window.Size, &window.AutoPosLastDirection, &r_outer, &r_avoid, ImGuiPopupPositionPolicy_Tooltip)
	}
	IM_ASSERT(false)
	return window.Pos
}

// FindBestWindowPosForPopupEx r_avoid = the rectangle to avoid (e.g. for tooltip it is a rectangle around the mouse cursor which we want to avoid. for popups it's a small point around the cursor.)
// r_outer = the visible area rectangle, minus safe area padding. If our popup size won't fit because of safe area padding we ignore it.
// (r_outer is usually equivalent to the viewport rectangle minus padding, but when multi-viewports are enabled and monitor
//
//	information are available, it may represent the entire platform monitor from the frame of reference of the current viewport.
//	this allows us to have tooltips/popups displayed out of the parent viewport.)
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

// ClosePopupsOverWindow When popups are stacked, clicking on a lower level popups puts focus back to it and close popups above it.
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

// IsPopupOpenID Supported flags: ImGuiPopupFlags_AnyPopupId, ImGuiPopupFlags_AnyPopupLevel
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

// OpenPopupEx Mark popup as open (toggle toward open state).
// Popups are closed when user click outside, or activate a pressable item, or CloseCurrentPopup() is called within a BeginPopup()/EndPopup() block.
// Popup identifiers are relative to the current ID-stack (so OpenPopup and BeginPopup needs to be at the same level).
// One open popup per level of the popup hierarchy (NB: when assigning we reset the Window member of ImGuiPopupRef to nil)
func OpenPopupEx(id ImGuiID, popup_flags ImGuiPopupFlags) {
	var g = GImGui
	var parent_window = g.CurrentWindow
	var current_stack_size = int(len(g.BeginPopupStack))

	if popup_flags&ImGuiPopupFlags_NoOpenOverExistingPopup != 0 {
		if IsPopupOpen("", ImGuiPopupFlags_AnyPopupId) {
			return
		}
	}

	var popup_ref ImGuiPopupData // Tagged as new ref as Window will be set back to nil if we write this into OpenPopupStack.
	popup_ref.PopupId = id
	popup_ref.Window = nil
	popup_ref.SourceWindow = g.NavWindow
	popup_ref.OpenFrameCount = g.FrameCount
	popup_ref.OpenParentId = parent_window.IDStack[len(parent_window.IDStack)-1]
	popup_ref.OpenPopupPos = NavCalcPreferredRefPos()
	if IsMousePosValid(&g.IO.MousePos) {
		popup_ref.OpenMousePos = g.IO.MousePos
	} else {
		popup_ref.OpenMousePos = popup_ref.OpenPopupPos
	}

	//IMGUI_DEBUG_LOG_POPUP("OpenPopupEx(0x%08X)\n", id)
	if int(len(g.OpenPopupStack)) < current_stack_size+1 {
		g.OpenPopupStack = append(g.OpenPopupStack, popup_ref)
	} else {
		// Gently handle the user mistakenly calling OpenPopup() every frame. It is a programming mistake! However, if we were to run the regular code path, the ui
		// would become completely unusable because the popup will always be in hidden-while-calculating-size state _while_ claiming focus. Which would be a very confusing
		// situation for the programmer. Instead, we silently allow the popup to proceed, it will keep reappearing and the programming error will be more obvious to understand.
		if g.OpenPopupStack[current_stack_size].PopupId == id && g.OpenPopupStack[current_stack_size].OpenFrameCount == g.FrameCount-1 {
			g.OpenPopupStack[current_stack_size].OpenFrameCount = popup_ref.OpenFrameCount
		} else {
			// Close child popups if any, then flag popup for open/reopen
			ClosePopupToLevel(current_stack_size, false)
			g.OpenPopupStack = append(g.OpenPopupStack, popup_ref)
		}

		// When reopening a popup we first refocus its parent, otherwise if its parent is itself a popup it would get closed by ClosePopupsOverWindow().
		// This is equivalent to what ClosePopupToLevel() does.
		//if (g.OpenPopupStack[current_stack_size].PopupId == id)
		//    FocusWindow(parent_window);
	}
}

func ClosePopupToLevel(remaining int, restore_focus_to_window_under_popup bool) {
	var g = GImGui
	//IMGUI_DEBUG_LOG_POPUP("ClosePopupToLevel(%d), restore_focus_to_window_under_popup=%d\n", remaining, restore_focus_to_window_under_popup)
	IM_ASSERT(remaining >= 0 && remaining < int(len(g.OpenPopupStack)))

	// Trim open popup stack
	var focus_window = g.OpenPopupStack[remaining].SourceWindow
	var popup_window = g.OpenPopupStack[remaining].Window
	g.OpenPopupStack = g.OpenPopupStack[:remaining]

	if restore_focus_to_window_under_popup {
		if focus_window != nil && !focus_window.WasActive && popup_window != nil {
			// Fallback
			FocusTopMostWindowUnderOne(popup_window, nil)
		} else {
			if g.NavLayer == ImGuiNavLayer_Main && focus_window != nil {
				focus_window = NavRestoreLastChildNavWindow(focus_window)
			}
			FocusWindow(focus_window)
		}
	}
}

func isPopupOpen(id ImGuiID, popup_flags ImGuiPopupFlags) bool {
	var g = GImGui
	if (popup_flags & ImGuiPopupFlags_AnyPopupId) != 0 {
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

// BeginPopupEx Attention! BeginPopup() adds default flags which BeginPopupEx()!
func BeginPopupEx(id ImGuiID, flags ImGuiWindowFlags) bool {
	var g = GImGui
	if !isPopupOpen(id, ImGuiPopupFlags_None) {
		g.NextWindowData.ClearFlags() // We behave like Begin() and need to consume those values
		return false
	}

	var name string
	if flags&ImGuiWindowFlags_ChildMenu != 0 {
		name = fmt.Sprintf("##Menu_%02d", int(len(g.BeginPopupStack))) // Recycle windows based on depth
	} else {
		name = fmt.Sprintf("##Popup_%08x", id) // Not recycling, so we can close/open during the same frame
	}

	flags |= ImGuiWindowFlags_Popup
	var is_open bool = Begin(string(name[:]), nil, flags)
	if !is_open { // NB: Begin can return false when the popup is completely clipped (e.g. zero size display)
		EndPopup()
	}

	return is_open
}
