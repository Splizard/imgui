package imgui

// "bool* p_selected" poto int the selection state (read-write), as a convenient helper.
func SelectablePointer(label string, p_selected *bool, flags ImGuiSelectableFlags, size_arg ImVec2) bool {
	if Selectable(label, *p_selected, flags, size_arg) {
		*p_selected = !*p_selected
		return true
	}
	return false
}

// Widgets: Selectables
// - A selectable highlights when hovered, and can display another color when selected.
// - Neighbors selectable extend their highlight bounds in order to leave no gap between them. This is so a series of selected Selectable appear contiguous.
// "selected bool" carry the selection state (read-only). Selectable() is clicked is returns true so you can modify your selection state. size.x==0.0: use remaining width, size.x>0.0: specify width. size.y==0.0: use label height, size.y>0.0: specify height
// Tip: pass a non-visible label (e.g. "##hello") then you can use the space to draw other text or image.
// But you need to make sure the ID is unique, e.g. enclose calls in PushID/PopID or use ##unique_id.
// With this scheme, ImGuiSelectableFlags_SpanAllColumns and ImGuiSelectableFlags_AllowItemOverlap are also frequently used flags.
// FIXME: Selectable() with (size.x == 0.0f) and (SelectableTextAlign.x > 0.0f) followed by SameLine() is currently not supported.
func Selectable(label string, selected bool, flags ImGuiSelectableFlags, size_arg ImVec2) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := g.Style

	// Submit label or explicit size to ItemSize(), whereas ItemAdd() will submit a larger/spanning rectangle.
	var id = window.GetIDs(label)
	var label_size = CalcTextSize(label, true, -1)
	var size = ImVec2{label_size.x, label_size.y}
	if size_arg.x != 0.0 {
		size.x = size_arg.x
	}
	if size_arg.y != 0.0 {
		size.y = size_arg.y
	}
	var pos = window.DC.CursorPos
	pos.y += window.DC.CurrLineTextBaseOffset
	ItemSizeVec(&size, 0.0)

	// Fill horizontal space
	// We don't support (size < 0.0f) in Selectable() because the ItemSpacing extension would make explicitly right-aligned sizes not visibly match other widgets.
	var span_all_columns = (flags & ImGuiSelectableFlags_SpanAllColumns) != 0
	var min_x = pos.x
	var max_x = window.WorkRect.Max.x
	if span_all_columns {
		min_x = window.ParentWorkRect.Min.x
		max_x = window.ParentWorkRect.Max.x
	}
	if size_arg.x == 0.0 || (flags&ImGuiSelectableFlags_SpanAvailWidth != 0) {
		size.x = max(label_size.x, max_x-min_x)
	}

	// Text stays at the submission position, but bounding box may be extended on both sides
	var text_min = pos
	var text_max = ImVec2{min_x + size.x, pos.y + size.y}

	// Selectables are meant to be tightly packed together with no click-gap, so we extend their box to cover spacing between selectable.
	var bb = ImRect{ImVec2{min_x, pos.y}, ImVec2{text_max.x, text_max.y}}
	if (flags & ImGuiSelectableFlags_NoPadWithHalfSpacing) == 0 {
		var spacing_x float
		if !span_all_columns {
			spacing_x = style.ItemSpacing.x
		}
		var spacing_y = style.ItemSpacing.y
		var spacing_L = IM_FLOOR(spacing_x * 0.50)
		var spacing_U = IM_FLOOR(spacing_y * 0.50)
		bb.Min.x -= spacing_L
		bb.Min.y -= spacing_U
		bb.Max.x += (spacing_x - spacing_L)
		bb.Max.y += (spacing_y - spacing_U)
	}
	//if (g.IO.KeyCtrl) { GetForegroundDrawList().AddRect(bb.Min, bb.Max, IM_COL32(0, 255, 0, 255)); }

	// Modify ClipRect for the ItemAdd(), faster than doing a PushColumnsBackground/PushTableBackground for every Selectable..
	var backup_clip_rect_min_x = window.ClipRect.Min.x
	var backup_clip_rect_max_x = window.ClipRect.Max.x
	if span_all_columns {
		window.ClipRect.Min.x = window.ParentWorkRect.Min.x
		window.ClipRect.Max.x = window.ParentWorkRect.Max.x
	}

	var disabled_item = (flags & ImGuiSelectableFlags_Disabled) != 0

	var disabled_flags = ImGuiItemFlags_None
	if disabled_item {
		disabled_flags |= ImGuiItemFlags_Disabled
	}

	var item_add = ItemAdd(&bb, id, nil, disabled_flags)
	if span_all_columns {
		window.ClipRect.Min.x = backup_clip_rect_min_x
		window.ClipRect.Max.x = backup_clip_rect_max_x
	}

	if !item_add {
		return false
	}

	var disabled_global = (g.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0
	if disabled_item && !disabled_global { // Only testing this as an optimization
		BeginDisabled(true)
	}

	// FIXME: We can standardize the behavior of those two, we could also keep the fast path of override ClipRect + full push on render only,
	// which would be advantageous since most selectable are not selected.
	if span_all_columns && window.DC.CurrentColumns != nil {
		PushColumnsBackground()
	} else if span_all_columns && g.CurrentTable != nil {
		TablePushBackgroundChannel()
	}

	// We use NoHoldingActiveID on menus so user can click and _hold_ on a menu then drag to browse child entries
	var button_flags ImGuiButtonFlags = 0
	if flags&ImGuiSelectableFlags_NoHoldingActiveID != 0 {
		button_flags |= ImGuiButtonFlags_NoHoldingActiveId
	}
	if flags&ImGuiSelectableFlags_SelectOnClick != 0 {
		button_flags |= ImGuiButtonFlags_PressedOnClick
	}
	if flags&ImGuiSelectableFlags_SelectOnRelease != 0 {
		button_flags |= ImGuiButtonFlags_PressedOnRelease
	}
	if flags&ImGuiSelectableFlags_AllowDoubleClick != 0 {
		button_flags |= ImGuiButtonFlags_PressedOnClickRelease | ImGuiButtonFlags_PressedOnDoubleClick
	}
	if flags&ImGuiSelectableFlags_AllowItemOverlap != 0 {
		button_flags |= ImGuiButtonFlags_AllowItemOverlap
	}

	var was_selected = selected
	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, button_flags)

	// Auto-select when moved into
	// - This will be more fully fleshed in the range-select branch
	// - This is not exposed as it won't nicely work with some user side handling of shift/control
	// - We cannot do 'if (g.NavJustMovedToId != id) { selected = false; pressed = was_selected; }' for two reasons
	//   - (1) it would require focus scope to be set, need exposing PushFocusScope() or equivalent (e.g. BeginSelection() calling PushFocusScope())
	//   - (2) usage will fail with clipped items
	//   The multi-select API aim to fix those issues, e.g. may be replaced with a BeginSelection() API.
	if (flags&ImGuiSelectableFlags_SelectOnNav != 0) && g.NavJustMovedToId != 0 && g.NavJustMovedToFocusScopeId == window.DC.NavFocusScopeIdCurrent {
		if g.NavJustMovedToId == id {
			selected = true
			pressed = true
		}
	}

	// Update NavId when clicking or when Hovering (this doesn't happen on most widgets), so navigation can be resumed with gamepad/keyboard
	if pressed || (hovered && (flags&ImGuiSelectableFlags_SetNavIdOnHover != 0)) {
		if !g.NavDisableMouseHover && g.NavWindow == window && g.NavLayer == window.DC.NavLayerCurrent {
			SetNavID(id, window.DC.NavLayerCurrent, window.DC.NavFocusScopeIdCurrent, &ImRect{bb.Min.Sub(window.Pos), bb.Max.Sub(window.Pos)}) // (bb == NavRect)
			g.NavDisableHighlight = true
		}
	}
	if pressed {
		MarkItemEdited(id)
	}

	if flags&ImGuiSelectableFlags_AllowItemOverlap != 0 {
		SetItemAllowOverlap()
	}

	// In this branch, Selectable() cannot toggle the selection so this will never trigger.
	if selected != was_selected { //-V547
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_ToggledSelection
	}

	// Render
	if held && (flags&ImGuiSelectableFlags_DrawHoveredWhenHeld != 0) {
		hovered = true
	}
	if hovered || selected {
		var c = ImGuiCol_Header
		switch {
		case (held && hovered):
			c = ImGuiCol_HeaderActive
		case hovered:
			c = ImGuiCol_HeaderHovered
		}
		var col = GetColorU32FromID(c, 1)
		RenderFrame(bb.Min, bb.Max, col, false, 0.0)
	}
	RenderNavHighlight(&bb, id, ImGuiNavHighlightFlags_TypeThin|ImGuiNavHighlightFlags_NoRounding)

	if span_all_columns && window.DC.CurrentColumns != nil {
		PopColumnsBackground()
	} else if span_all_columns && g.CurrentTable != nil {
		TablePopBackgroundChannel()
	}

	RenderTextClipped(&text_min, &text_max, label, &label_size, &style.SelectableTextAlign, &bb)

	// Automatically close popups
	if pressed && (window.Flags&ImGuiWindowFlags_Popup != 0) && (flags&ImGuiSelectableFlags_DontClosePopups == 0) && (g.LastItemData.InFlags&ImGuiItemFlags_SelectableDontClosePopup == 0) {
		CloseCurrentPopup()
	}

	if disabled_item && !disabled_global {
		EndDisabled()
	}

	return pressed //-V1020
}

// Focus, Activation
// - Prefer using "SetItemDefaultFocus()" over "if (IsWindowAppearing()) SetScrollHereY()" when applicable to signify "this is the default item"

// make last item the default focused item of a window.
func SetItemDefaultFocus() {
	window := g.CurrentWindow
	if !window.Appearing {
		return
	}
	if g.NavWindow == window.RootWindowForNav && (g.NavInitRequest || g.NavInitResultId != 0) && g.NavLayer == window.DC.NavLayerCurrent {
		g.NavInitRequest = false
		g.NavInitResultId = g.LastItemData.ID
		g.NavInitResultRectRel = ImRect{g.LastItemData.Rect.Min.Sub(window.Pos), g.LastItemData.Rect.Max.Sub(window.Pos)}
		NavUpdateAnyRequestFlag()
		if !IsItemVisible() {
			SetScrollHereY(0.5)
		}
	}
}

// focus keyboard on the next widget. Use positive 'offset' to access sub components of a multiple component widget. Use -1 to access previous widget.
func SetKeyboardFocusHere(offset int) {
	IM_ASSERT(offset >= -1) // -1 is allowed but not below
	window := g.CurrentWindow
	g.TabFocusRequestNextWindow = window
	g.TabFocusRequestNextCounterRegular = window.DC.FocusCounterRegular + 1 + offset
	g.TabFocusRequestNextCounterTabStop = INT_MAX
}

// Focus Scope (WIP)
// This is generally used to identify a selection set (multiple of which may be in the same window), as selection
// patterns generally need to react (e.g. clear selection) when landing on an item of the set.
func PushFocusScope(id ImGuiID) {
	window := g.CurrentWindow
	g.FocusScopeStack = append(g.FocusScopeStack, window.DC.NavFocusScopeIdCurrent)
	window.DC.NavFocusScopeIdCurrent = id
}

func PopFocusScope() {
	window := g.CurrentWindow
	IM_ASSERT(len(g.FocusScopeStack) > 0) // Too many PopFocusScope() ?
	window.DC.NavFocusScopeIdCurrent = g.FocusScopeStack[len(g.FocusScopeStack)-1]
	g.FocusScopeStack = g.FocusScopeStack[:len(g.FocusScopeStack)-1]
}

func GetFocusedFocusScope() ImGuiID { g := g; return g.NavFocusScopeId } // Focus scope which is actually active
func GetFocusScope() ImGuiID {
	return g.CurrentWindow.DC.NavFocusScopeIdCurrent
} // Focus scope we are outputting into, set by PushFocusScope()

// Notifies Dear ImGui when hosting platform windows lose or gain input focus
func (io *ImGuiIO) AddFocusEvent(focused bool) {
	if focused {
		return
	}

	// Clear buttons state when focus is lost
	// (this is useful so e.g. releasing Alt after focus loss on Alt-Tab doesn't trigger the Alt menu toggle)
	io.KeysDown = [len(io.KeysDown)]bool{}
	for n := range io.KeysDownDuration {
		io.KeysDownDuration[n] = -1
		io.KeysDownDurationPrev[n] = -1
	}
	io.KeyCtrl = false
	io.KeyShift = false
	io.KeyAlt = false
	io.KeySuper = false
	io.KeyMods = ImGuiKeyModFlags_None
	io.KeyModsPrev = ImGuiKeyModFlags_None
	for n := range io.NavInputsDownDuration {
		io.NavInputsDownDuration[n] = -1
		io.NavInputsDownDurationPrev[n] = -1
	}
}

// Windows: Display Order and Focus Order
func FocusWindow(window *ImGuiWindow) {
	g := g

	if g.NavWindow != window {
		g.NavWindow = window
		if window != nil && g.NavDisableMouseHover {
			g.NavMousePosDirty = true
		}
		if window != nil {
			g.NavId = window.NavLastIds[0]
		} else {
			g.NavId = 0
		}
		g.NavFocusScopeId = 0
		g.NavIdIsAlive = false
		g.NavLayer = ImGuiNavLayer_Main
		g.NavInitRequest = false
		g.NavMoveSubmitted = false
		g.NavMoveScoringItems = false
		NavUpdateAnyRequestFlag()
		//IMGUI_DEBUG_LOG("FocusWindow(\"%s\")\n", window ? window.Name : nil);
	}

	// Close popups if any
	ClosePopupsOverWindow(window, false)

	// Move the root window to the top of the pile
	IM_ASSERT(window == nil || window.RootWindow != nil)

	var focus_front_window *ImGuiWindow
	var display_front_window *ImGuiWindow
	if window != nil {
		focus_front_window = window
		display_front_window = window
	}

	// Steal active widgets. Some of the cases it triggers includes:
	// - Focus a window while an InputText in another window is active, if focus happens before the old InputText can run.
	// - When using Nav to activate menu items (due to timing of activating on press.new window appears.losing ActiveId)
	if g.ActiveId != 0 && g.ActiveIdWindow != nil && g.ActiveIdWindow.RootWindow != focus_front_window {
		if !g.ActiveIdNoClearOnFocusLoss {
			ClearActiveID()
		}
	}

	// Passing nil allow to disable keyboard focus
	if window == nil {
		return
	}

	// Bring to front
	BringWindowToFocusFront(focus_front_window)
	if ((window.Flags | display_front_window.Flags) & ImGuiWindowFlags_NoBringToFrontOnFocus) == 0 {
		BringWindowToDisplayFront(display_front_window)
	}
}

func BringWindowToFocusFront(window *ImGuiWindow) {
	IM_ASSERT(window == window.RootWindow)

	var cur_order = window.FocusOrder
	IM_ASSERT(g.WindowsFocusOrder[cur_order] == window)
	if g.WindowsFocusOrder[len(g.WindowsFocusOrder)-1] == window {
		return
	}

	var new_order = int(len(g.WindowsFocusOrder)) - 1
	for n := int(cur_order); n < new_order; n++ {
		g.WindowsFocusOrder[n] = g.WindowsFocusOrder[n+1]
		g.WindowsFocusOrder[n].FocusOrder--
		IM_ASSERT(int(g.WindowsFocusOrder[n].FocusOrder) == n)
	}
	g.WindowsFocusOrder[new_order] = window
	window.FocusOrder = (short)(new_order)
}

func BringWindowToDisplayFront(window *ImGuiWindow) {
	var current_front_window = g.Windows[len(g.Windows)-1]
	if current_front_window == window || current_front_window.RootWindow == window { // Cheap early out (could be better)
		return
	}
	for i := len(g.Windows) - 2; i >= 0; i-- { // We can ignore the top-most window
		if g.Windows[i] == window {
			amount := len(g.Windows) - i - 1
			copy(g.Windows[i:], g.Windows[i+1:i+1+amount])
			g.Windows[len(g.Windows)-1] = window
			break
		}
	}
}

func FocusTopMostWindowUnderOne(under_this_window *ImGuiWindow, ignore_window *ImGuiWindow) {
	g := g

	var start_idx int
	if under_this_window != nil {
		start_idx = FindWindowFocusIndex(under_this_window)
	} else {
		start_idx = int(len(g.WindowsFocusOrder) - 1)
	}

	for i := start_idx; i >= 0; i-- {
		// We may later decide to test for different NoXXXInputs based on the active navigation input (mouse vs nav) but that may feel more confusing to the user.
		var window = g.WindowsFocusOrder[i]
		IM_ASSERT(window == window.RootWindow)
		if window != ignore_window && window.WasActive {
			if (window.Flags & (ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs)) != (ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs) {
				var focus_window = NavRestoreLastChildNavWindow(window)
				FocusWindow(focus_window)
				return
			}
		}
	}
	FocusWindow(nil)
}
