package imgui

// Widgets: Menus
// - Use BeginMenuBar() on a window ImGuiWindowFlags_MenuBar to append to its menu bar.
// - Use BeginMainMenuBar() to create a menu bar at the top of the screen and append to it.
// - Use BeginMenu() to create a menu. You can call BeginMenu() multiple time with the same identifier to append more items to it.
// - Not that MenuItem() keyboardshortcuts are displayed as a convenience but _not processed_ by Dear ImGui at the moment.

// append to menu-bar of current window (requires ImGuiWindowFlags_MenuBar flag set on parent window).
// FIXME: Provided a rectangle perhaps e.g. a BeginMenuBarEx() could be used anywhere..
// Currently the main responsibility of this function being to setup clip-rect + horizontal layout + menu navigation layer.
// Ideally we also want this to be responsible for claiming space out of the main window scrolling rectangle, in which case ImGuiWindowFlags_MenuBar will become unnecessary.
// Then later the same system could be used for multiple menu-bars, scrollbars, side-bars.
func BeginMenuBar() bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	if window.Flags&ImGuiWindowFlags_MenuBar != 0 {
		return false
	}

	IM_ASSERT(!window.DC.MenuBarAppending)
	BeginGroup() // Backup position on layer 0 // FIXME: Misleading to use a group for that backup/restore
	PushString("##menubar")

	// We don't clip with current window clipping rectangle as it is already set to the area below. However we clip with window full rect.
	// We remove 1 worth of rounding to Max.x to that text in long menus and small windows don't tend to display over the lower-right rounded area, which looks particularly glitchy.
	var bar_rect = window.MenuBarRect()
	var clip_rect = ImRect{
		ImVec2{
			IM_ROUND(bar_rect.Min.x + window.WindowBorderSize),
			IM_ROUND(bar_rect.Min.y + window.WindowBorderSize),
		},
		ImVec2{
			IM_ROUND(ImMax(bar_rect.Min.x, bar_rect.Max.x-ImMax(window.WindowRounding, window.WindowBorderSize))),
			IM_ROUND(bar_rect.Max.y),
		},
	}
	clip_rect.ClipWith(window.OuterRectClipped)
	PushClipRect(clip_rect.Min, clip_rect.Max, false)

	// We overwrite CursorMaxPos because BeginGroup sets it to CursorPos (essentially the .EmitItem hack in EndMenuBar() would need something analogous here, maybe a BeginGroupEx() with flags).
	window.DC.CursorPos = ImVec2{bar_rect.Min.x + window.DC.MenuBarOffset.x, bar_rect.Min.y + window.DC.MenuBarOffset.y}
	window.DC.CursorMaxPos = window.DC.CursorPos
	window.DC.LayoutType = ImGuiLayoutType_Horizontal
	window.DC.NavLayerCurrent = ImGuiNavLayer_Menu
	window.DC.MenuBarAppending = true
	AlignTextToFramePadding()
	return true
}

// only call EndMenuBar() if BeginMenuBar() returns true!
func EndMenuBar() {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return
	}
	g := GImGui

	// Nav: When a move request within one of our child menu failed, capture the request to navigate among our siblings.
	if NavMoveRequestButNoResultYet() && (g.NavMoveDir == ImGuiDir_Left || g.NavMoveDir == ImGuiDir_Right) && (g.NavWindow.Flags&ImGuiWindowFlags_ChildMenu) != 0 {
		// Try to find out if the request is for one of our child menu
		var nav_earliest_child = g.NavWindow
		for nav_earliest_child.ParentWindow != nil && (nav_earliest_child.ParentWindow.Flags&ImGuiWindowFlags_ChildMenu) != 0 {
			nav_earliest_child = nav_earliest_child.ParentWindow
		}
		if nav_earliest_child.ParentWindow == window && nav_earliest_child.DC.ParentLayoutType == ImGuiLayoutType_Horizontal && (g.NavMoveFlags&ImGuiNavMoveFlags_Forwarded) == 0 {
			// To do so we claim focus back, restore NavId and then process the movement request for yet another frame.
			// This involve a one-frame delay which isn't very problematic in this situation. We could remove it by scoring in advance for multiple window (probably not worth bothering)
			var layer = ImGuiNavLayer_Menu
			IM_ASSERT(window.DC.NavLayersActiveMaskNext&(1<<layer) != 0) // Sanity check
			FocusWindow(window)
			SetNavID(window.NavLastIds[layer], layer, 0, &window.NavRectRel[layer])
			g.NavDisableHighlight = true // Hide highlight for the current frame so we don't see the intermediary selection.
			g.NavDisableMouseHover = true
			g.NavMousePosDirty = true
			NavMoveRequestForward(g.NavMoveDir, g.NavMoveClipDir, g.NavMoveFlags) // Repeat
		}
	}

	IM_ASSERT(window.Flags&ImGuiWindowFlags_MenuBar != 0)
	IM_ASSERT(window.DC.MenuBarAppending)
	PopClipRect()
	PopID()
	window.DC.MenuBarOffset.x = window.DC.CursorPos.x - window.Pos.x // Save horizontal position so next append can reuse it. This is kinda equivalent to a per-layer CursorPos.
	g.GroupStack[len(g.GroupStack)-1].EmitItem = false
	EndGroup() // Restore position on layer 0
	window.DC.LayoutType = ImGuiLayoutType_Vertical
	window.DC.NavLayerCurrent = ImGuiNavLayer_Main
	window.DC.MenuBarAppending = false
}

// create and append to a full screen menu-bar.
func BeginMainMenuBar() bool {
	g := GImGui
	var viewport = GetMainViewport()

	// For the main menu bar, which cannot be moved, we honor g.Style.DisplaySafeAreaPadding to ensure text can be visible on a TV set.
	// FIXME: This could be generalized as an opt-in way to clamp window.DC.CursorStartPos to avoid SafeArea?
	// FIXME: Consider removing support for safe area down the line... it's messy. Nowadays consoles have support for TV calibration in OS settings.
	g.NextWindowData.MenuBarOffsetMinVal = ImVec2{g.Style.DisplaySafeAreaPadding.x, ImMax(g.Style.DisplaySafeAreaPadding.y-g.Style.FramePadding.y, 0.0)}
	var window_flags = ImGuiWindowFlags_NoScrollbar | ImGuiWindowFlags_NoSavedSettings | ImGuiWindowFlags_MenuBar
	var height = GetFrameHeight()
	var is_open = BeginViewportSideBar("##MainMenuBar", viewport, ImGuiDir_Up, height, window_flags)
	g.NextWindowData.MenuBarOffsetMinVal = ImVec2{}

	if is_open {
		BeginMenuBar()
	} else {
		End()
	}

	return is_open
}

// only call EndMainMenuBar() if BeginMainMenuBar() returns true!
func EndMainMenuBar() {
	EndMenuBar()

	// When the user has left the menu layer (typically: closed menus through activation of an item), we restore focus to the previous window
	// FIXME: With this strategy we won't be able to restore a nil focus.
	g := GImGui
	if g.CurrentWindow == g.NavWindow && g.NavLayer == ImGuiNavLayer_Main && !g.NavAnyRequest {
		FocusTopMostWindowUnderOne(g.NavWindow, nil)
	}

	End()
}

// create a sub-menu entry. only call EndMenu() if this returns true!
func BeginMenu(label string, enabled bool /*= true*/) bool {
	return BeginMenuEx(label, "", enabled)
}

// only call EndMenu() if BeginMenu() returns true!
func EndMenu() {
	// Nav: When a left move request _within our child menu_ failed, close ourselves (the _parent_ menu).
	// A menu doesn't close itself because EndMenuBar() wants the catch the last Left<>Right inputs.
	// However, it means that with the current code, a BeginMenu() from outside another menu or a menu-bar won't be closable with the Left direction.
	g := GImGui
	var window = g.CurrentWindow
	if g.NavWindow != nil && g.NavWindow.ParentWindow == window && g.NavMoveDir == ImGuiDir_Left && NavMoveRequestButNoResultYet() && window.DC.LayoutType == ImGuiLayoutType_Vertical {
		ClosePopupToLevel(int(len(g.BeginPopupStack)), true)
		NavMoveRequestCancel()
	}

	EndPopup()
}

// return true when activated.
func MenuItem(label string, shortcut string /*= L*/, selected *bool /*= e*/, enabled bool /*= true*/) bool {
	return MenuItemEx(label, "", shortcut, selected, enabled)
}

// return true when activated + toggle (*p_selected) if p_selected != NULL
func MenuItemSelected(label string, shortcut string, p_selected *bool, enabled bool /*= true*/) bool {
	var b = false
	if p_selected != nil {
		b = *p_selected
	}
	if MenuItemEx(label, "", shortcut, &b, enabled) {
		if p_selected != nil {
			*p_selected = !*p_selected
		}
		return true
	}
	return false
}

// Important: calling order matters!
// FIXME: Somehow overlapping with docking tech.
// FIXME: The "rect-cut" aspect of this could be formalized into a lower-level helper (rect-cut: https://halt.software/dead-simple-layouts)
func BeginViewportSideBar(name string, viewport_p *ImGuiViewport, dir ImGuiDir, axis_size float, window_flags ImGuiWindowFlags) bool {
	IM_ASSERT(dir != ImGuiDir_None)

	var bar_window = FindWindowByName(name)
	if bar_window == nil || bar_window.BeginCount == 0 {
		// Calculate and set window size/position
		var viewport = viewport_p
		if viewport_p == nil {
			viewport = GetMainViewport()
		}

		var avail_rect = viewport.GetBuildWorkRect()
		var axis = ImGuiAxis_X
		if dir == ImGuiDir_Up || dir == ImGuiDir_Down {
			axis = ImGuiAxis_Y
		}
		var pos = avail_rect.Min
		if dir == ImGuiDir_Right || dir == ImGuiDir_Down {
			switch axis {
			case ImGuiAxis_X:
				pos.x = avail_rect.Max.x - axis_size
			case ImGuiAxis_Y:
				pos.y = avail_rect.Max.y - axis_size
			}
		}
		var size = avail_rect.GetSize()
		switch axis {
		case ImGuiAxis_X:
			size.x = axis_size
		case ImGuiAxis_Y:
			size.y = axis_size
		}
		SetNextWindowPos(&pos, 0, ImVec2{})
		SetNextWindowSize(&size, 0)

		// Report our size into work area (for next frame) using actual window size
		if dir == ImGuiDir_Up || dir == ImGuiDir_Left {
			switch axis {
			case ImGuiAxis_X:
				viewport.BuildWorkOffsetMin.x += axis_size
			case ImGuiAxis_Y:
				viewport.BuildWorkOffsetMin.y += axis_size
			}
		} else if dir == ImGuiDir_Down || dir == ImGuiDir_Right {
			switch axis {
			case ImGuiAxis_X:
				viewport.BuildWorkOffsetMin.x -= axis_size
			case ImGuiAxis_Y:
				viewport.BuildWorkOffsetMin.y -= axis_size
			}
		}
	}

	window_flags |= ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoMove
	PushStyleFloat(ImGuiStyleVar_WindowRounding, 0.0)
	PushStyleVec(ImGuiStyleVar_WindowMinSize, ImVec2{}) // Lift normal size constraint
	var is_open = Begin(name, nil, window_flags)
	PopStyleVar(2)

	return is_open
}

// Menus
func BeginMenuEx(label string, icon string, enabled bool /*= true*/) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	var style = g.Style
	var id = window.GetIDs(label)
	var menu_is_open = IsPopupOpenID(id, ImGuiPopupFlags_None)

	// Sub-menus are ChildWindow so that mouse can be hovering across them (otherwise top-most popup menu would steal focus and not allow hovering on parent menu)
	var flags = ImGuiWindowFlags_ChildMenu | ImGuiWindowFlags_AlwaysAutoResize | ImGuiWindowFlags_NoMove | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoSavedSettings | ImGuiWindowFlags_NoNavFocus
	if (window.Flags & (ImGuiWindowFlags_Popup | ImGuiWindowFlags_ChildMenu)) != 0 {
		flags |= ImGuiWindowFlags_ChildWindow
	}

	// If a menu with same the ID was already submitted, we will append to it, matching the behavior of Begin().
	// We are relying on a O(N) search - so O(N log N) over the frame - which seems like the most efficient for the expected small amount of BeginMenu() calls per frame.
	// If somehow this is ever becoming a problem we can switch to use e.g. ImGuiStorage mapping key to last frame used.
	var contains bool
	for _, menu := range g.MenusIdSubmittedThisFrame {
		if menu == id {
			contains = true
			break
		}
	}
	if contains {
		if menu_is_open {
			menu_is_open = BeginPopupEx(id, flags) // menu_is_open can be 'false' when the popup is completely clipped (e.g. zero size display)
		} else {
			g.NextWindowData.ClearFlags() // we behave like Begin() and need to consume those values
		}
		return menu_is_open
	}

	// Tag menu as used. Next time BeginMenu() with same ID is called it will append to existing menu
	g.MenusIdSubmittedThisFrame = append(g.MenusIdSubmittedThisFrame, id)

	var label_size = CalcTextSize(label, true, -1)
	var pressed bool
	var menuset_is_open = (window.Flags&ImGuiWindowFlags_Popup) == 0 && (len(g.OpenPopupStack) > len(g.BeginPopupStack) && g.OpenPopupStack[len(g.BeginPopupStack)].OpenParentId == window.IDStack[len(window.IDStack)-1])
	var backed_nav_window = g.NavWindow
	if menuset_is_open {
		g.NavWindow = window // Odd hack to allow hovering across menus of a same menu-set (otherwise we wouldn't be able to hover parent)
	}

	// The reference position stored in popup_pos will be used by Begin() to find a suitable position for the child menu,
	// However the final position is going to be different! It is chosen by FindBestWindowPosForPopup().
	// e.g. Menus tend to overlap each other horizontally to amplify relative Z-ordering.
	var popup_pos ImVec2
	var pos = window.DC.CursorPos
	PushString(label)
	if !enabled {
		BeginDisabled(true)
	}
	var offsets = &window.DC.MenuColumns
	if window.DC.LayoutType == ImGuiLayoutType_Horizontal {
		// Menu inside an horizontal menu bar
		// Selectable extend their highlight by half ItemSpacing in each direction.
		// For ChildMenu, the popup position will be overwritten by the call to FindBestWindowPosForPopup() in Begin()
		popup_pos = ImVec2{pos.x - 1.0 - IM_FLOOR(style.ItemSpacing.x*0.5), pos.y - style.FramePadding.y + window.MenuBarHeight()}
		window.DC.CursorPos.x += IM_FLOOR(style.ItemSpacing.x * 0.5)
		PushStyleVec(ImGuiStyleVar_ItemSpacing, ImVec2{style.ItemSpacing.x * 2.0, style.ItemSpacing.y})
		var w = label_size.x
		var text_pos = ImVec2{window.DC.CursorPos.x + float(offsets.OffsetLabel), window.DC.CursorPos.y + window.DC.CurrLineTextBaseOffset}
		pressed = Selectable("", menu_is_open, ImGuiSelectableFlags_NoHoldingActiveID|ImGuiSelectableFlags_SelectOnClick|ImGuiSelectableFlags_DontClosePopups, ImVec2{w, 0.0})
		RenderText(text_pos, label, true)
		PopStyleVar(1)
		window.DC.CursorPos.x += IM_FLOOR(style.ItemSpacing.x * (-1.0 + 0.5)) // -1 spacing to compensate the spacing added when Selectable() did a SameLine(). It would also work to call SameLine() ourselves after the PopStyleVar().
	} else {
		// Menu inside a menu
		// (In a typical menu window where all items are BeginMenu() or MenuItem() calls, extra_w will always be 0.0f.
		//  Only when they are other items sticking out we're going to add spacing, yet only register minimum width into the layout system.
		popup_pos = ImVec2{pos.x, pos.y - style.WindowPadding.y}
		var icon_w float
		if icon != "" {
			icon_w = CalcTextSize(icon, true, -1).x
		}
		var checkmark_w = IM_FLOOR(g.FontSize * 1.20)
		var min_w = window.DC.MenuColumns.DeclColumns(icon_w, label_size.x, 0.0, checkmark_w) // Feedback to next frame
		var extra_w = ImMax(0.0, GetContentRegionAvail().x-min_w)
		var text_pos = ImVec2{window.DC.CursorPos.x + float(offsets.OffsetLabel), window.DC.CursorPos.y + window.DC.CurrLineTextBaseOffset}
		pressed = Selectable("", menu_is_open, ImGuiSelectableFlags_NoHoldingActiveID|ImGuiSelectableFlags_SelectOnClick|ImGuiSelectableFlags_DontClosePopups|ImGuiSelectableFlags_SpanAvailWidth, ImVec2{min_w, 0.0})
		RenderText(text_pos, label, true)
		if icon_w > 0.0 {
			RenderText(pos.Add(ImVec2{float(offsets.OffsetIcon), 0.0}), icon, true)
		}
		RenderArrow(window.DrawList, pos.Add(ImVec2{float(offsets.OffsetMark) + extra_w + g.FontSize*0.30, 0.0}), GetColorU32FromID(ImGuiCol_Text, 1), ImGuiDir_Right, 1)
	}
	if !enabled {
		EndDisabled()
	}

	var hovered = (g.HoveredId == id) && enabled
	if menuset_is_open {
		g.NavWindow = backed_nav_window
	}

	var want_open = false
	var want_close = false
	if window.DC.LayoutType == ImGuiLayoutType_Vertical { // (window.Flags & (ImGuiWindowFlags_Popup|ImGuiWindowFlags_ChildMenu))

		// Close menu when not hovering it anymore unless we are moving roughly in the direction of the menu
		// Implement http://bjk5.com/post/44698559168/breaking-down-amazons-mega-dropdown to avoid using timers, so menus feels more reactive.
		var moving_toward_other_child_menu = false

		var child_menu_window *ImGuiWindow = nil
		if len(g.BeginPopupStack) < len(g.OpenPopupStack) && g.OpenPopupStack[len(g.BeginPopupStack)].SourceWindow == window {
			child_menu_window = g.OpenPopupStack[len(g.BeginPopupStack)].Window
		}
		if g.HoveredWindow == window && child_menu_window != nil && window.Flags&ImGuiWindowFlags_MenuBar == 0 {
			// FIXME-DPI: Values should be derived from a master "scale" factor.
			var next_window_rect = child_menu_window.Rect()
			var ta = g.IO.MousePos.Sub(g.IO.MouseDelta)
			var tb ImVec2
			var tc ImVec2
			if window.Pos.x < child_menu_window.Pos.x {
				tb, tc = next_window_rect.GetTL(), next_window_rect.GetBL()
			} else {
				tb, tc = next_window_rect.GetTR(), next_window_rect.GetBR()
			}
			var extra = ImClamp(ImFabs(ta.x-tb.x)*0.30, 5.0, 30.0) // add a bit of extra slack.
			if window.Pos.x < child_menu_window.Pos.x {
				ta.x += -0.5
			} else { // to avoid numerical issues
				ta.x += +0.5
			}
			tb.y = ta.y + ImMax((tb.y-extra)-ta.y, -100.0) // triangle is maximum 200 high to limit the slope and the bias toward large sub-menus // FIXME: Multiply by fb_scale?
			tc.y = ta.y + ImMin((tc.y+extra)-ta.y, +100.0)
			moving_toward_other_child_menu = ImTriangleContainsPoint(&ta, &tb, &tc, &g.IO.MousePos)
			//GetForegroundDrawList().AddTriangleFilled(ta, tb, tc, moving_within_opened_triangle ? IM_COL32(0,128,0,128) : IM_COL32(128,0,0,128)); // [DEBUG]
		}

		// FIXME: Hovering a disabled BeginMenu or MenuItem won't close us
		if menu_is_open && !hovered && g.HoveredWindow == window && g.HoveredIdPreviousFrame != 0 && g.HoveredIdPreviousFrame != id && !moving_toward_other_child_menu {
			want_close = true
		}

		if !menu_is_open && hovered && pressed { // Click to open
			want_open = true
		} else if !menu_is_open && hovered && !moving_toward_other_child_menu { // Hover to open
			want_open = true
		}

		if g.NavActivateId == id {
			want_close = menu_is_open
			want_open = !menu_is_open
		}
		if g.NavId == id && g.NavMoveDir == ImGuiDir_Right { // Nav-Right to open
			want_open = true
			NavMoveRequestCancel()
		}
	} else {
		// Menu bar
		if menu_is_open && pressed && menuset_is_open { // Click an open menu again to close it
			want_close = true
			want_open = false
			menu_is_open = false
		} else if pressed || (hovered && menuset_is_open && !menu_is_open) { // First click to open, then hover to open others
			want_open = true
		} else if g.NavId == id && g.NavMoveDir == ImGuiDir_Down { // Nav-Down to open
			want_open = true
			NavMoveRequestCancel()
		}
	}

	if !enabled { // explicitly close if an open menu becomes disabled, facilitate users code a lot in pattern such as 'if (BeginMenu("options", has_object)) { ..use object.. }'
		want_close = true
	}
	if want_close && IsPopupOpenID(id, ImGuiPopupFlags_None) {
		ClosePopupToLevel(int(len(g.BeginPopupStack)), true)
	}
	PopID()

	if !menu_is_open && want_open && len(g.OpenPopupStack) > len(g.BeginPopupStack) {
		// Don't recycle same menu level in the same frame, first close the other menu and yield for a frame.
		OpenPopup(label, 0)
		return false
	}

	menu_is_open = menu_is_open || want_open
	if want_open {
		OpenPopup(label, 0)
	}

	if menu_is_open {
		SetNextWindowPos(&popup_pos, ImGuiCond_Always, ImVec2{}) // Note: this is super misleading! The value will serve as reference for FindBestWindowPosForPopup(), not actual pos.
		menu_is_open = BeginPopupEx(id, flags)                   // menu_is_open can be 'false' when the popup is completely clipped (e.g. zero size display)
	} else {
		g.NextWindowData.ClearFlags() // We behave like Begin() and need to consume those values
	}

	return menu_is_open
}

func MenuItemEx(label string, icon string, shortcut string, selected *bool, enabled bool /*= true*/) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	var style = g.Style
	var pos = window.DC.CursorPos
	var label_size = CalcTextSize(label, true, -1)

	// We've been using the equivalent of ImGuiSelectableFlags_SetNavIdOnHover on all Selectable() since early Nav system days (commit 43ee5d73),
	// but I am unsure whether this should be kept at all. For now moved it to be an opt-in feature used by menus only.
	var pressed bool
	PushString(label)
	if !enabled {
		BeginDisabled(true)
	}
	var flags = ImGuiSelectableFlags_SelectOnRelease | ImGuiSelectableFlags_SetNavIdOnHover
	var offsets = &window.DC.MenuColumns
	if window.DC.LayoutType == ImGuiLayoutType_Horizontal {
		// Mimic the exact layout spacing of BeginMenu() to allow MenuItem() inside a menu bar, which is a little misleading but may be useful
		// Note that in this situation: we don't render the shortcut, we render a highlight instead of the selected tick mark.
		var w = label_size.x
		window.DC.CursorPos.x += IM_FLOOR(style.ItemSpacing.x * 0.5)
		PushStyleVec(ImGuiStyleVar_ItemSpacing, ImVec2{style.ItemSpacing.x * 2.0, style.ItemSpacing.y})
		pressed = SelectablePointer("", selected, flags, ImVec2{w, 0.0})
		PopStyleVar(1)
		RenderText(pos.Add(ImVec2{float(offsets.OffsetLabel), 0.0}), label, true)
		window.DC.CursorPos.x += IM_FLOOR(style.ItemSpacing.x * (-1.0 + 0.5)) // -1 spacing to compensate the spacing added when Selectable() did a SameLine(). It would also work to call SameLine() ourselves after the PopStyleVar().
	} else {
		// Menu item inside a vertical menu
		// (In a typical menu window where all items are BeginMenu() or MenuItem() calls, extra_w will always be 0.0f.
		//  Only when they are other items sticking out we're going to add spacing, yet only register minimum width into the layout system.
		var icon_w float
		if icon != "" {
			icon_w = CalcTextSize(icon, true, -1).x
		}
		var shortcut_w float
		if shortcut != "" {
			shortcut_w = CalcTextSize(shortcut, true, -1).x
		}
		var checkmark_w = IM_FLOOR(g.FontSize * 1.20)
		var min_w = window.DC.MenuColumns.DeclColumns(icon_w, label_size.x, shortcut_w, checkmark_w) // Feedback for next frame
		var stretch_w = ImMax(0.0, GetContentRegionAvail().x-min_w)
		pressed = Selectable("", false, flags|ImGuiSelectableFlags_SpanAvailWidth, ImVec2{min_w, 0.0})
		RenderText(pos.Add(ImVec2{float(offsets.OffsetLabel), 0.0}), label, true)
		if icon_w > 0.0 {
			RenderText(pos.Add(ImVec2{float(offsets.OffsetIcon), 0.0}), icon, true)
		}
		if shortcut_w > 0.0 {
			PushStyleColorVec(ImGuiCol_Text, &style.Colors[ImGuiCol_TextDisabled])
			RenderText(pos.Add(ImVec2{float(offsets.OffsetShortcut) + stretch_w, 0.0}), shortcut, false)
			PopStyleColor(1)
		}
		if *selected {
			RenderCheckMark(window.DrawList, pos.Add(ImVec2{float(offsets.OffsetMark) + stretch_w + g.FontSize*0.40, g.FontSize * 0.134 * 0.5}), GetColorU32FromID(ImGuiCol_Text, 1), g.FontSize*0.866)
		}
	}
	if !enabled {
		EndDisabled()
	}
	PopID()

	return pressed
}
