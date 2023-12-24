package imgui

// NAV_WINDOWING_LIST_APPEAR_DELAY When using CTRL+TAB (or Gamepad Square+L/R) we delay the visual a little in order to reduce visual noise doing a fast switch.
const NAV_WINDOWING_LIST_APPEAR_DELAY float = 0.15 // Time before the window list starts to appear

const NAV_WINDOWING_HIGHLIGHT_DELAY float = 0.20 // Time before the highlight and screen dimming starts fading in

func FindWindowNavFocusable(i_start, i_stop, dir int) *ImGuiWindow { // FIXME-OPT O(N)
	for i := i_start; i >= 0 && i < int(len(guiContext.WindowsFocusOrder)) && i != i_stop; i += dir {
		if IsWindowNavFocusable(guiContext.WindowsFocusOrder[i]) {
			return guiContext.WindowsFocusOrder[i]
		}
	}
	return nil
}

// NavMoveRequestForward Forward will reuse the move request again on the next frame (generally with modifications done to it)
func NavMoveRequestForward(move_dir ImGuiDir, clip_dir ImGuiDir, move_flags ImGuiNavMoveFlags) {
	IM_ASSERT(!guiContext.NavMoveForwardToNextFrame)
	NavMoveRequestCancel()
	guiContext.NavMoveForwardToNextFrame = true
	guiContext.NavMoveDir = move_dir
	guiContext.NavMoveClipDir = clip_dir
	guiContext.NavMoveFlags = move_flags | ImGuiNavMoveFlags_Forwarded
}

// NavMoveRequestTryWrapping Navigation wrap-around logic is delayed to the end of the frame because this operation is only valid after entire
// popup is assembled and in case of appended popups it is not clear which EndPopup() call is final.
func NavMoveRequestTryWrapping(window *ImGuiWindow, move_flags ImGuiNavMoveFlags) {
	IM_ASSERT(move_flags != 0) // Call with _WrapX, _WrapY, _LoopX, _LoopY
	// In theory we should test for NavMoveRequestButNoResultYet() but there's no point doing it, NavEndFrame() will do the same test
	if guiContext.NavWindow == window && guiContext.NavMoveScoringItems && guiContext.NavLayer == ImGuiNavLayer_Main {
		guiContext.NavMoveFlags |= move_flags
	}
}

// SetNavID FIXME-NAV: The existence of SetNavID vs SetFocusID properly needs to be clarified/reworked.
// In our terminology those should be interchangeable. Those two functions are merely a legacy artifact, so at minimum naming should be clarified.
func SetNavID(id ImGuiID, nav_layer ImGuiNavLayer, focus_scope_id ImGuiID, rect_rel *ImRect) {
	IM_ASSERT(guiContext.NavWindow != nil)
	IM_ASSERT(nav_layer == ImGuiNavLayer_Main || nav_layer == ImGuiNavLayer_Menu)
	guiContext.NavId = id
	guiContext.NavLayer = nav_layer
	guiContext.NavFocusScopeId = focus_scope_id
	guiContext.NavWindow.NavLastIds[nav_layer] = id
	guiContext.NavWindow.NavRectRel[nav_layer] = *rect_rel
	//guiContext.NavDisableHighlight = false;
	//guiContext.NavDisableMouseHover = guiContext.NavMousePosDirty = true;
}

func NavUpdateAnyRequestFlag() {
	guiContext.NavAnyRequest = guiContext.NavMoveScoringItems || guiContext.NavInitRequest
	if guiContext.NavAnyRequest {
		IM_ASSERT(guiContext.NavWindow != nil)
	}
}

// NavRestoreLastChildNavWindow Restore the last focused child.
// Call when we are expected to land on the Main Layer (0) after FocusWindow()
func NavRestoreLastChildNavWindow(window *ImGuiWindow) *ImGuiWindow {
	if window.NavLastChildNavWindow != nil && window.NavLastChildNavWindow.WasActive {
		return window.NavLastChildNavWindow
	}
	return window
}

func GetNavInputAmount(n ImGuiNavInput, mode ImGuiInputReadMode) float {
	if mode == ImGuiInputReadMode_Down {
		return guiContext.IO.NavInputs[n] // Instant, read analog input (0.0f..1.0f, as provided by user)
	}

	var t = guiContext.IO.NavInputsDownDuration[n]
	if t < 0.0 && mode == ImGuiInputReadMode_Released { // Return 1.0f when just released, no repeat, ignore analog input.
		if guiContext.IO.NavInputsDownDurationPrev[n] >= 0.0 {
			return 1
		}
		return 0
	}

	if t < 0.0 {
		return 0.0
	}

	if mode == ImGuiInputReadMode_Pressed { // Return 1.0f when just pressed, no repeat, ignore analog input.
		if t == 0.0 {
			return 1
		}
		return 0.0
	}

	if mode == ImGuiInputReadMode_Repeat {
		return (float)(CalcTypematicRepeatAmount(t-guiContext.IO.DeltaTime, t, guiContext.IO.KeyRepeatDelay*0.72, guiContext.IO.KeyRepeatRate*0.80))
	}
	if mode == ImGuiInputReadMode_RepeatSlow {
		return (float)(CalcTypematicRepeatAmount(t-guiContext.IO.DeltaTime, t, guiContext.IO.KeyRepeatDelay*1.25, guiContext.IO.KeyRepeatRate*2.00))
	}
	if mode == ImGuiInputReadMode_RepeatFast {
		return (float)(CalcTypematicRepeatAmount(t-guiContext.IO.DeltaTime, t, guiContext.IO.KeyRepeatDelay*0.72, guiContext.IO.KeyRepeatRate*0.30))
	}
	return 0.0
}

// NavMoveRequestButNoResultYet Gamepad/Keyboard Navigation
func NavMoveRequestButNoResultYet() bool {
	return guiContext.NavMoveScoringItems && guiContext.NavMoveResultLocal.ID == 0 && guiContext.NavMoveResultOther.ID == 0
}

func NavApplyItemToResult(result *ImGuiNavItemData) {
	window := guiContext.CurrentWindow
	result.Window = window
	result.ID = guiContext.LastItemData.ID
	result.FocusScopeId = window.DC.NavFocusScopeIdCurrent
	result.RectRel = ImRect{guiContext.LastItemData.NavRect.Min.Sub(window.Pos), guiContext.LastItemData.NavRect.Max.Sub(window.Pos)}
}

// NavProcessItem We get there when either NavId == id, or when guiContext.NavAnyRequest is set (which is updated by NavUpdateAnyRequestFlag above)
// This is called after LastItemData is set.
func NavProcessItem() {
	window := guiContext.CurrentWindow
	var id = guiContext.LastItemData.ID
	var nav_bb = guiContext.LastItemData.NavRect
	var item_flags = guiContext.LastItemData.InFlags

	// Process Init Request
	if guiContext.NavInitRequest && guiContext.NavLayer == window.DC.NavLayerCurrent {
		// Even if 'ImGuiItemFlags_NoNavDefaultFocus' is on (typically collapse/close button) we record the first ResultId so they can be used as a fallback
		var candidate_for_nav_default_focus = (item_flags & (ImGuiItemFlags_NoNavDefaultFocus | ImGuiItemFlags_Disabled)) == 0
		if candidate_for_nav_default_focus || guiContext.NavInitResultId == 0 {
			guiContext.NavInitResultId = id
			guiContext.NavInitResultRectRel = ImRect{nav_bb.Min.Sub(window.Pos), nav_bb.Max.Sub(window.Pos)}
		}
		if candidate_for_nav_default_focus {
			guiContext.NavInitRequest = false // Found a match, clear request
			NavUpdateAnyRequestFlag()
		}
	}

	// Process Move Request (scoring for navigation)
	// FIXME-NAV: Consider policy for double scoring (scoring from NavScoringRect + scoring from a rect wrapped according to current wrapping policy)
	if guiContext.NavMoveScoringItems {
		if (guiContext.NavId != id || (guiContext.NavMoveFlags&ImGuiNavMoveFlags_AllowCurrentNavId != 0)) && item_flags&(ImGuiItemFlags_Disabled|ImGuiItemFlags_NoNav) == 0 {
			var result *ImGuiNavItemData
			if window == guiContext.NavWindow {
				result = &guiContext.NavMoveResultLocal
			} else {
				result = &guiContext.NavMoveResultOther
			}
			if NavScoreItem(result) {
				NavApplyItemToResult(result)
			}

			// Features like PageUp/PageDown need to maintain a separate score for the visible set of items.
			const VISIBLE_RATIO float = 0.70
			if (guiContext.NavMoveFlags&ImGuiNavMoveFlags_AlsoScoreVisibleSet) != 0 && window.ClipRect.Overlaps(nav_bb) {
				if ImClamp(nav_bb.Max.y, window.ClipRect.Min.y, window.ClipRect.Max.y)-ImClamp(nav_bb.Min.y, window.ClipRect.Min.y, window.ClipRect.Max.y) >= (nav_bb.Max.y-nav_bb.Min.y)*VISIBLE_RATIO {
					if NavScoreItem(&guiContext.NavMoveResultLocalVisible) {
						NavApplyItemToResult(&guiContext.NavMoveResultLocalVisible)
					}
				}
			}
		}
	}

	// Update window-relative bounding box of navigated item
	if guiContext.NavId == id {
		guiContext.NavWindow = window // Always refresh guiContext.NavWindow, because some operations such as FocusItem() don't have a window.
		guiContext.NavLayer = window.DC.NavLayerCurrent
		guiContext.NavFocusScopeId = window.DC.NavFocusScopeIdCurrent
		guiContext.NavIdIsAlive = true
		window.NavRectRel[window.DC.NavLayerCurrent] = ImRect{nav_bb.Min.Sub(window.Pos), nav_bb.Max.Sub(window.Pos)} // Store item bounding box (relative to window position)
	}
}

// NavUpdatePageUpPageDown Handle PageUp/PageDown/Home/End keys
// Called from NavUpdateCreateMoveRequest() which will use our output to create a move request
// FIXME-NAV: This doesn't work properly with NavFlattened siblings as we use NavWindow rectangle for reference
// FIXME-NAV: how to get Home/End to aim at the beginning/end of a 2D grid?
func NavUpdatePageUpPageDown() float {
	io := guiContext.IO

	var window = guiContext.NavWindow
	if (window.Flags&ImGuiWindowFlags_NoNavInputs != 0) || guiContext.NavWindowingTarget != nil || guiContext.NavLayer != ImGuiNavLayer_Main {
		return 0.0
	}

	var page_up_held = IsKeyDown(io.KeyMap[ImGuiKey_PageUp]) && !IsActiveIdUsingKey(ImGuiKey_PageUp)
	var page_down_held = IsKeyDown(io.KeyMap[ImGuiKey_PageDown]) && !IsActiveIdUsingKey(ImGuiKey_PageDown)
	var home_pressed = IsKeyPressed(io.KeyMap[ImGuiKey_Home], true) && !IsActiveIdUsingKey(ImGuiKey_Home)
	var end_pressed = IsKeyPressed(io.KeyMap[ImGuiKey_End], true) && !IsActiveIdUsingKey(ImGuiKey_End)
	if page_up_held == page_down_held && home_pressed == end_pressed { // Proceed if either (not both) are pressed, otherwise early out
		return 0.0
	}

	if window.DC.NavLayersActiveMask == 0x00 && window.DC.NavHasScroll {
		// Fallback manual-scroll when window has no navigable item
		if IsKeyPressed(io.KeyMap[ImGuiKey_PageUp], true) {
			setScrollY(window, window.Scroll.y-window.InnerRect.GetHeight())
		} else if IsKeyPressed(io.KeyMap[ImGuiKey_PageDown], true) {
			setScrollY(window, window.Scroll.y+window.InnerRect.GetHeight())
		} else if home_pressed {
			setScrollY(window, 0.0)
		} else if end_pressed {
			setScrollY(window, window.ScrollMax.y)
		}
	} else {
		var nav_rect_rel = &window.NavRectRel[guiContext.NavLayer]
		var page_offset_y = max(0.0, window.InnerRect.GetHeight()-window.CalcFontSize()*1.0+nav_rect_rel.GetHeight())
		var nav_scoring_rect_offset_y float = 0.0
		if IsKeyPressed(io.KeyMap[ImGuiKey_PageUp], true) {
			nav_scoring_rect_offset_y = -page_offset_y
			guiContext.NavMoveDir = ImGuiDir_Down // Because our scoring rect is offset up, we request the down direction (so we can always land on the last item)
			guiContext.NavMoveClipDir = ImGuiDir_Up
			guiContext.NavMoveFlags = ImGuiNavMoveFlags_AllowCurrentNavId | ImGuiNavMoveFlags_AlsoScoreVisibleSet
		} else if IsKeyPressed(io.KeyMap[ImGuiKey_PageDown], true) {
			nav_scoring_rect_offset_y = +page_offset_y
			guiContext.NavMoveDir = ImGuiDir_Up // Because our scoring rect is offset down, we request the up direction (so we can always land on the last item)
			guiContext.NavMoveClipDir = ImGuiDir_Down
			guiContext.NavMoveFlags = ImGuiNavMoveFlags_AllowCurrentNavId | ImGuiNavMoveFlags_AlsoScoreVisibleSet
		} else if home_pressed {
			// FIXME-NAV: handling of Home/End is assuming that the top/bottom most item will be visible with Scroll.y == 0/ScrollMax.y
			// Scrolling will be handled via the ImGuiNavMoveFlags_ScrollToEdge flag, we don't scroll immediately to avoid scrolling happening before nav result.
			// Preserve current horizontal position if we have any.
			nav_rect_rel.Min.y = -window.Scroll.y
			nav_rect_rel.Max.y = -window.Scroll.y
			if nav_rect_rel.IsInverted() {
				nav_rect_rel.Min.x = 0
				nav_rect_rel.Max.x = 0.0
			}
			guiContext.NavMoveDir = ImGuiDir_Down
			guiContext.NavMoveFlags = ImGuiNavMoveFlags_AllowCurrentNavId | ImGuiNavMoveFlags_ScrollToEdge
			// FIXME-NAV: MoveClipDir left to _None, intentional?
		} else if end_pressed {
			nav_rect_rel.Min.y = window.ScrollMax.y + window.SizeFull.y - window.Scroll.y
			nav_rect_rel.Max.y = window.ScrollMax.y + window.SizeFull.y - window.Scroll.y
			if nav_rect_rel.IsInverted() {
				nav_rect_rel.Min.x = 0
				nav_rect_rel.Max.x = 0.0
			}
			guiContext.NavMoveDir = ImGuiDir_Up
			guiContext.NavMoveFlags = ImGuiNavMoveFlags_AllowCurrentNavId | ImGuiNavMoveFlags_ScrollToEdge
			// FIXME-NAV: MoveClipDir left to _None, intentional?
		}
		return nav_scoring_rect_offset_y
	}
	return 0.0
}

func NavUpdateCreateMoveRequest() {
	io := guiContext.IO
	var window = guiContext.NavWindow

	if guiContext.NavMoveForwardToNextFrame && window != nil {
		// Forwarding previous request (which has been modified, e.guiContext. wrap around menus rewrite the requests with a starting rectangle at the other side of the window)
		// (preserve most state, which were already set by the NavMoveRequestForward() function)
		IM_ASSERT(guiContext.NavMoveDir != ImGuiDir_None && guiContext.NavMoveClipDir != ImGuiDir_None)
		IM_ASSERT(guiContext.NavMoveFlags&ImGuiNavMoveFlags_Forwarded != 0)
		//IMGUI_DEBUG_LOG_NAV("[nav] NavMoveRequestForward %d\n", guiContext.NavMoveDir)
	} else {
		// Initiate directional inputs request
		guiContext.NavMoveDir = ImGuiDir_None
		guiContext.NavMoveFlags = ImGuiNavMoveFlags_None
		if window != nil && guiContext.NavWindowingTarget == nil && window.Flags&ImGuiWindowFlags_NoNavInputs == 0 {
			var read_mode = ImGuiInputReadMode_Repeat
			if !IsActiveIdUsingNavDir(ImGuiDir_Left) && (IsNavInputTest(ImGuiNavInput_DpadLeft, read_mode) || IsNavInputTest(ImGuiNavInput_KeyLeft_, read_mode)) {
				guiContext.NavMoveDir = ImGuiDir_Left
			}
			if !IsActiveIdUsingNavDir(ImGuiDir_Right) && (IsNavInputTest(ImGuiNavInput_DpadRight, read_mode) || IsNavInputTest(ImGuiNavInput_KeyRight_, read_mode)) {
				guiContext.NavMoveDir = ImGuiDir_Right
			}
			if !IsActiveIdUsingNavDir(ImGuiDir_Up) && (IsNavInputTest(ImGuiNavInput_DpadUp, read_mode) || IsNavInputTest(ImGuiNavInput_KeyUp_, read_mode)) {
				guiContext.NavMoveDir = ImGuiDir_Up
			}
			if !IsActiveIdUsingNavDir(ImGuiDir_Down) && (IsNavInputTest(ImGuiNavInput_DpadDown, read_mode) || IsNavInputTest(ImGuiNavInput_KeyDown_, read_mode)) {
				guiContext.NavMoveDir = ImGuiDir_Down
			}
		}
		guiContext.NavMoveClipDir = guiContext.NavMoveDir
	}

	// Update PageUp/PageDown/Home/End scroll
	// FIXME-NAV: Consider enabling those keys even without the master ImGuiConfigFlags_NavEnableKeyboard flag?
	var nav_keyboard_active = (io.ConfigFlags & ImGuiConfigFlags_NavEnableKeyboard) != 0
	var scoring_rect_offset_y float = 0.0
	if window != nil && guiContext.NavMoveDir == ImGuiDir_None && nav_keyboard_active {
		scoring_rect_offset_y = NavUpdatePageUpPageDown()
	}

	// Submit
	guiContext.NavMoveForwardToNextFrame = false
	if guiContext.NavMoveDir != ImGuiDir_None {
		NavMoveRequestSubmit(guiContext.NavMoveDir, guiContext.NavMoveClipDir, guiContext.NavMoveFlags)
	}

	// Moving with no reference triggers a init request (will be used as a fallback if the direction fails to find a match)
	if guiContext.NavMoveSubmitted && guiContext.NavId == 0 {
		//IMGUI_DEBUG_LOG_NAV("[nav] NavInitRequest: from move, window \"%s\", layer=%d\n", guiContext.NavWindow.Name, guiContext.NavLayer)
		guiContext.NavInitRequest = true
		guiContext.NavInitRequestFromMove = true
		guiContext.NavInitResultId = 0
		guiContext.NavDisableHighlight = false
	}

	// When using gamepad, we project the reference nav bounding box into window visible area.
	// This is to allow resuming navigation inside the visible area after doing a large amount of scrolling, since with gamepad every movements are relative
	// (can't focus a visible object like we can with the mouse).
	if guiContext.NavMoveSubmitted && guiContext.NavInputSource == ImGuiInputSource_Gamepad && guiContext.NavLayer == ImGuiNavLayer_Main && window != nil {
		var window_rect_rel = ImRect{window.InnerRect.Min.Sub(window.Pos).Sub(ImVec2{1, 1}), window.InnerRect.Max.Sub(window.Pos).Add(ImVec2{1, 1})}
		if !window_rect_rel.ContainsRect(window.NavRectRel[guiContext.NavLayer]) {
			//IMGUI_DEBUG_LOG_NAV("[nav] NavMoveRequest: clamp NavRectRel\n")
			var pad = window.CalcFontSize() * 0.5
			window_rect_rel.ExpandVec(ImVec2{-min(window_rect_rel.GetWidth(), pad), -min(window_rect_rel.GetHeight(), pad)}) // Terrible approximation for the intent of starting navigation from first fully visible item
			window.NavRectRel[guiContext.NavLayer].ClipWithFull(window_rect_rel)
			guiContext.NavId = 0
			guiContext.NavFocusScopeId = 0
		}
	}

	// For scoring we use a single segment on the left side our current item bounding box (not touching the edge to avoid box overlap with zero-spaced items)
	var scoring_rect ImRect
	if window != nil {
		var nav_rect_rel ImRect
		if !window.NavRectRel[guiContext.NavLayer].IsInverted() {
			nav_rect_rel = window.NavRectRel[guiContext.NavLayer]
		}
		scoring_rect = ImRect{window.Pos.Add(nav_rect_rel.Min), window.Pos.Add(nav_rect_rel.Max)}
		scoring_rect.TranslateY(scoring_rect_offset_y)
		scoring_rect.Min.x = min(scoring_rect.Min.x+1.0, scoring_rect.Max.x)
		scoring_rect.Max.x = scoring_rect.Min.x
		IM_ASSERT(!scoring_rect.IsInverted()) // Ensure if we have a finite, non-inverted bounding box here will allows us to remove extraneous ImFabs() calls in NavScoreItem().
		//GetForegroundDrawList().AddRect(scoring_rect.Min, scoring_rect.Max, IM_COL32(255,200,0,255)); // [DEBUG]
	}
	guiContext.NavScoringRect = scoring_rect
}

// NavUpdateCancelRequest Process NavCancel input (to close a popup, get back to parent, clear focus)
// FIXME: In order to support e.guiContext. Escape to clear a selection we'll need:
// - either to store the equivalent of ActiveIdUsingKeyInputMask for a FocusScope and test for it.
// - either to move most/all of those tests to the epilogue/end functions of the scope they are dealing with (e.guiContext. exit child window in EndChild()) or in EndFrame(), to allow an earlier intercept
func NavUpdateCancelRequest() {
	if !IsNavInputTest(ImGuiNavInput_Cancel, ImGuiInputReadMode_Pressed) {
		return
	}

	//IMGUI_DEBUG_LOG_NAV("[nav] ImGuiNavInput_Cancel\n")
	if guiContext.ActiveId != 0 {
		if !IsActiveIdUsingNavInput(ImGuiNavInput_Cancel) {
			ClearActiveID()
		}
	} else if guiContext.NavLayer != ImGuiNavLayer_Main {
		// Leave the "menu" layer
		NavRestoreLayer(ImGuiNavLayer_Main)
	} else if guiContext.NavWindow != nil && guiContext.NavWindow != guiContext.NavWindow.RootWindow && guiContext.NavWindow.Flags&ImGuiWindowFlags_Popup == 0 && guiContext.NavWindow.ParentWindow != nil {
		// Exit child window
		var child_window = guiContext.NavWindow
		var parent_window = guiContext.NavWindow.ParentWindow
		IM_ASSERT(child_window.ChildId != 0)
		var child_rect = child_window.Rect()
		FocusWindow(parent_window)
		SetNavID(child_window.ChildId, ImGuiNavLayer_Main, 0, &ImRect{child_rect.Min.Sub(parent_window.Pos), child_rect.Max.Sub(parent_window.Pos)})
	} else if len(guiContext.OpenPopupStack) > 0 {
		// Close open popup/menu
		if guiContext.OpenPopupStack[len(guiContext.OpenPopupStack)-1].Window.Flags&ImGuiWindowFlags_Modal == 0 {
			ClosePopupToLevel(int(len(guiContext.OpenPopupStack)-1), true)
		}
	} else {
		// Clear NavLastId for popups but keep it for regular child window so we can leave one and come back where we were
		if guiContext.NavWindow != nil && ((guiContext.NavWindow.Flags&ImGuiWindowFlags_Popup != 0) || guiContext.NavWindow.Flags&ImGuiWindowFlags_ChildWindow == 0) {
			guiContext.NavWindow.NavLastIds[0] = 0
		}
		guiContext.NavId = 0
		guiContext.NavFocusScopeId = 0
	}
}

func NavRestoreLayer(layer ImGuiNavLayer) {
	if layer == ImGuiNavLayer_Main {
		guiContext.NavWindow = NavRestoreLastChildNavWindow(guiContext.NavWindow)
	}
	var window = guiContext.NavWindow
	if window.NavLastIds[layer] != 0 {
		SetNavID(window.NavLastIds[layer], layer, 0, &window.NavRectRel[layer])
	} else {
		guiContext.NavLayer = layer
		NavInitWindow(window, true)
	}
	guiContext.NavDisableHighlight = false
	guiContext.NavDisableMouseHover = true
	guiContext.NavMousePosDirty = true
}

func FindWindowFocusIndex(window *ImGuiWindow) int {
	var order = int(window.FocusOrder)
	IM_ASSERT(guiContext.WindowsFocusOrder[order] == window)
	return order
}

func NavUpdateWindowingHighlightWindow(focus_change_dir int) {
	IM_ASSERT(guiContext.NavWindowingTarget != nil)
	if guiContext.NavWindowingTarget.Flags&ImGuiWindowFlags_Modal != 0 {
		return
	}

	var i_current = FindWindowFocusIndex(guiContext.NavWindowingTarget)
	var window_target = FindWindowNavFocusable(i_current+focus_change_dir, -INT_MAX, focus_change_dir)
	if window_target == nil {
		var start int
		if focus_change_dir < 0 {
			start = int(len(guiContext.WindowsFocusOrder) - 1)
		}
		window_target = FindWindowNavFocusable(start, i_current, focus_change_dir)
	}
	if window_target != nil { // Don't reset windowing target if there's a single window in the list
		guiContext.NavWindowingTarget = window_target
		guiContext.NavWindowingTargetAnim = window_target
	}
	guiContext.NavWindowingToggleLayer = false
}

func NavUpdateInitResult() {
	// In very rare cases guiContext.NavWindow may be nil (e.guiContext. clearing focus after requesting an init request, which does happen when releasing Alt while clicking on void)
	if guiContext.NavWindow == nil {
		return
	}

	// Apply result from previous navigation init request (will typically select the first item, unless SetItemDefaultFocus() has been called)
	// FIXME-NAV: On _NavFlattened windows, guiContext.NavWindow will only be updated during subsequent frame. Not a problem currently.
	//IMGUI_DEBUG_LOG_NAV("[nav] NavInitRequest: result NavID 0x%08X in Layer %d Window \"%s\"\n", guiContext.NavInitResultId, guiContext.NavLayer, guiContext.NavWindow.Name)
	SetNavID(guiContext.NavInitResultId, guiContext.NavLayer, 0, &guiContext.NavInitResultRectRel)
	guiContext.NavIdIsAlive = true // Mark as alive from previous frame as we got a result
	if guiContext.NavInitRequestFromMove {
		guiContext.NavDisableHighlight = false
		guiContext.NavDisableMouseHover = true
		guiContext.NavMousePosDirty = true
	}
}

func NavCalcPreferredRefPos() ImVec2 {
	if guiContext.NavDisableHighlight || !guiContext.NavDisableMouseHover || guiContext.NavWindow == nil {
		// Mouse (we need a fallback in case the mouse becomes invalid after being used)
		if IsMousePosValid(&guiContext.IO.MousePos) {
			return guiContext.IO.MousePos
		}
		return guiContext.MouseLastValidPos
	} else {
		// When navigation is active and mouse is disabled, decide on an arbitrary position around the bottom left of the currently navigated item.
		var rect_rel = &guiContext.NavWindow.NavRectRel[guiContext.NavLayer]
		var pos = guiContext.NavWindow.Pos.Add(ImVec2{rect_rel.Min.x + min(guiContext.Style.FramePadding.x*4, rect_rel.GetWidth()), rect_rel.Max.y - min(guiContext.Style.FramePadding.y, rect_rel.GetHeight())})
		var viewport = GetMainViewport()

		clamped := ImClampVec2(&pos, &viewport.Pos, viewport.Pos.Add(viewport.Size))
		return *ImFloorVec(&clamped) // ImFloor() is important because non-integer mouse position application in backend might be lossy and result in undesirable non-zero delta.
	}
}

// NavSaveLastChildNavWindowIntoParent FIXME: This could be replaced by updating a frame number in each window when (window == NavWindow) and (NavLayer == 0).
// This way we could find the last focused window among our children. It would be much less confusing this way?
func NavSaveLastChildNavWindowIntoParent(nav_window *ImGuiWindow) {
	var parent = nav_window
	for parent != nil && parent.RootWindow != parent && (parent.Flags&(ImGuiWindowFlags_Popup|ImGuiWindowFlags_ChildMenu)) == 0 {
		parent = parent.ParentWindow
	}
	if parent != nil && parent != nav_window {
		parent.NavLastChildNavWindow = nav_window
	}
}

// NavUpdateWindowing Windowing management mode
// Keyboard: CTRL+Tab (change focus/move/resize), Alt (toggle menu layer)
// Gamepad:  Hold Menu/Square (change focus/move/resize), Tap Menu/Square (toggle menu layer)
func NavUpdateWindowing() {
	io := guiContext.IO

	var apply_focus_window *ImGuiWindow = nil
	var apply_toggle_layer = false

	var modal_window = GetTopMostPopupModal()
	var allow_windowing = (modal_window == nil)
	if !allow_windowing {
		guiContext.NavWindowingTarget = nil
	}

	// Fade out
	if guiContext.NavWindowingTargetAnim != nil && guiContext.NavWindowingTarget == nil {
		guiContext.NavWindowingHighlightAlpha = max(guiContext.NavWindowingHighlightAlpha-io.DeltaTime*10.0, 0.0)
		if guiContext.DimBgRatio <= 0.0 && guiContext.NavWindowingHighlightAlpha <= 0.0 {
			guiContext.NavWindowingTargetAnim = nil
		}
	}

	// Start CTRL-TAB or Square+L/R window selection
	var start_windowing_with_gamepad = allow_windowing && guiContext.NavWindowingTarget == nil && IsNavInputTest(ImGuiNavInput_Menu, ImGuiInputReadMode_Pressed)
	var start_windowing_with_keyboard = allow_windowing && guiContext.NavWindowingTarget == nil && io.KeyCtrl && IsKeyPressedMap(ImGuiKey_Tab, true) && (io.ConfigFlags&ImGuiConfigFlags_NavEnableKeyboard) != 0
	if start_windowing_with_gamepad || start_windowing_with_keyboard {
		var window *ImGuiWindow
		if guiContext.NavWindow != nil {
			window = guiContext.NavWindow
		} else {
			window = FindWindowNavFocusable(int(len(guiContext.WindowsFocusOrder)-1), -INT_MAX, -1)
		}
		if window != nil {
			guiContext.NavWindowingTarget = window.RootWindow
			guiContext.NavWindowingTargetAnim = window.RootWindow
			guiContext.NavWindowingTimer = 0.0
			guiContext.NavWindowingHighlightAlpha = 0.0
			if start_windowing_with_gamepad { // Gamepad starts toggling layer
				guiContext.NavWindowingToggleLayer = true
			} else {
				guiContext.NavWindowingToggleLayer = false
			}
			if start_windowing_with_keyboard {
				guiContext.NavInputSource = ImGuiInputSource_Keyboard
			} else {
				guiContext.NavInputSource = ImGuiInputSource_Gamepad
			}
		}
	}

	// Gamepad update
	guiContext.NavWindowingTimer += io.DeltaTime
	if guiContext.NavWindowingTarget != nil && guiContext.NavInputSource == ImGuiInputSource_Gamepad {
		// Highlight only appears after a brief time holding the button, so that a fast tap on PadMenu (to toggle NavLayer) doesn't add visual noise
		guiContext.NavWindowingHighlightAlpha = max(guiContext.NavWindowingHighlightAlpha, ImSaturate((guiContext.NavWindowingTimer-NAV_WINDOWING_HIGHLIGHT_DELAY)/0.05))

		// Select window to focus
		var focus_change_dir = bool2int(IsNavInputTest(ImGuiNavInput_FocusPrev, ImGuiInputReadMode_RepeatSlow)) - bool2int(IsNavInputTest(ImGuiNavInput_FocusNext, ImGuiInputReadMode_RepeatSlow))
		if focus_change_dir != 0 {
			NavUpdateWindowingHighlightWindow(focus_change_dir)
			guiContext.NavWindowingHighlightAlpha = 1.0
		}

		// Single press toggles NavLayer, long press with L/R apply actual focus on release (until then the window was merely rendered top-most)
		if !IsNavInputDown(ImGuiNavInput_Menu) {
			guiContext.NavWindowingToggleLayer = guiContext.NavWindowingToggleLayer && (guiContext.NavWindowingHighlightAlpha < 1.0) // Once button was held long enough we don't consider it a tap-to-toggle-layer press anymore.
			if guiContext.NavWindowingToggleLayer && guiContext.NavWindow != nil {
				apply_toggle_layer = true
			} else if !guiContext.NavWindowingToggleLayer {
				apply_focus_window = guiContext.NavWindowingTarget
			}
			guiContext.NavWindowingTarget = nil
		}
	}

	// Keyboard: Focus
	if guiContext.NavWindowingTarget != nil && guiContext.NavInputSource == ImGuiInputSource_Keyboard {
		// Visuals only appears after a brief time after pressing TAB the first time, so that a fast CTRL+TAB doesn't add visual noise
		guiContext.NavWindowingHighlightAlpha = max(guiContext.NavWindowingHighlightAlpha, ImSaturate((guiContext.NavWindowingTimer-NAV_WINDOWING_HIGHLIGHT_DELAY)/0.05)) // 1.0f
		if IsKeyPressedMap(ImGuiKey_Tab, true) {
			if io.KeyShift {
				NavUpdateWindowingHighlightWindow(+1)
			} else {
				NavUpdateWindowingHighlightWindow(-1)
			}

		}
		if !io.KeyCtrl {
			apply_focus_window = guiContext.NavWindowingTarget
		}
	}

	// Keyboard: Press and Release ALT to toggle menu layer
	// - Testing that only Alt is tested prevents Alt+Shift or AltGR from toggling menu layer.
	// - AltGR is normally Alt+Ctrl but we can't reliably detect it (not all backends/systems/layout emit it as Alt+Ctrl). But even on keyboards without AltGR we don't want Alt+Ctrl to open menu anyway.
	if io.KeyMods == ImGuiKeyModFlags_Alt && (io.KeyModsPrev&ImGuiKeyModFlags_Alt) == 0 {
		guiContext.NavWindowingToggleLayer = true
		guiContext.NavInputSource = ImGuiInputSource_Keyboard
	}
	if guiContext.NavWindowingToggleLayer && guiContext.NavInputSource == ImGuiInputSource_Keyboard {
		// We cancel toggling nav layer when any text has been typed (generally while holding Alt). (See #370)
		// We cancel toggling nav layer when other modifiers are pressed. (See #4439)
		if len(io.InputQueueCharacters) > 0 || io.KeyCtrl || io.KeyShift || io.KeySuper {
			guiContext.NavWindowingToggleLayer = false
		}

		// Apply layer toggle on release
		// Important: we don't assume that Alt was previously held in order to handle loss of focus when backend calls io.AddFocusEvent(false)
		// Important: as before version <18314 we lacked an explicit IO event for focus gain/loss, we also compare mouse validity to detect old backends clearing mouse pos on focus loss.
		if io.KeyMods&ImGuiKeyModFlags_Alt == 0 && (io.KeyModsPrev&ImGuiKeyModFlags_Alt != 0) && guiContext.NavWindowingToggleLayer {
			if guiContext.ActiveId == 0 || guiContext.ActiveIdAllowOverlap {
				if IsMousePosValid(&io.MousePos) == IsMousePosValid(&io.MousePosPrev) {
					apply_toggle_layer = true
				}
			}
		}
		if !io.KeyAlt {
			guiContext.NavWindowingToggleLayer = false
		}
	}

	// Move window
	if guiContext.NavWindowingTarget != nil && guiContext.NavWindowingTarget.Flags&ImGuiWindowFlags_NoMove == 0 {
		var move_delta ImVec2
		if guiContext.NavInputSource == ImGuiInputSource_Keyboard && !io.KeyShift {
			move_delta = GetNavInputAmount2d(ImGuiNavDirSourceFlags_Keyboard, ImGuiInputReadMode_Down, 0, 0)
		}
		if guiContext.NavInputSource == ImGuiInputSource_Gamepad {
			move_delta = GetNavInputAmount2d(ImGuiNavDirSourceFlags_PadLStick, ImGuiInputReadMode_Down, 0, 0)
		}
		if move_delta.x != 0.0 || move_delta.y != 0.0 {
			const NAV_MOVE_SPEED float = 800.0
			var move_speed = ImFloor(NAV_MOVE_SPEED * io.DeltaTime * min(io.DisplayFramebufferScale.x, io.DisplayFramebufferScale.y)) // FIXME: Doesn't handle variable framerate very well
			var moving_window = guiContext.NavWindowingTarget.RootWindow
			p := moving_window.Pos.Add(move_delta.Scale(move_speed))
			setWindowPos(moving_window, &p, ImGuiCond_Always)
			MarkIniSettingsDirtyWindow(moving_window)
			guiContext.NavDisableMouseHover = true
		}
	}

	// Apply final focus
	if apply_focus_window != nil && (guiContext.NavWindow == nil || apply_focus_window != guiContext.NavWindow.RootWindow) {
		ClearActiveID()
		guiContext.NavDisableHighlight = false
		guiContext.NavDisableMouseHover = true
		apply_focus_window = NavRestoreLastChildNavWindow(apply_focus_window)
		ClosePopupsOverWindow(apply_focus_window, false)
		FocusWindow(apply_focus_window)
		if apply_focus_window.NavLastIds[0] == 0 {
			NavInitWindow(apply_focus_window, false)
		}

		// If the window has ONLY a menu layer (no main layer), select it directly
		// Use NavLayersActiveMaskNext since windows didn't have a chance to be Begin()-ed on this frame,
		// so CTRL+Tab where the keys are only held for 1 frame will be able to use correct layers mask since
		// the target window as already been previewed once.
		// FIXME-NAV: This should be done in NavInit.. or in FocusWindow... However in both of those cases,
		// we won't have a guarantee that windows has been visible before and therefore NavLayersActiveMask*
		// won't be valid.
		if apply_focus_window.DC.NavLayersActiveMaskNext == (1 << ImGuiNavLayer_Menu) {
			guiContext.NavLayer = ImGuiNavLayer_Menu
		}
	}
	if apply_focus_window != nil {
		guiContext.NavWindowingTarget = nil
	}

	// Apply menu/layer toggle
	if apply_toggle_layer && guiContext.NavWindow != nil {
		ClearActiveID()

		// Move to parent menu if necessary
		var new_nav_window = guiContext.NavWindow
		for new_nav_window.ParentWindow != nil &&
			(new_nav_window.DC.NavLayersActiveMask&(1<<ImGuiNavLayer_Menu)) == 0 &&
			(new_nav_window.Flags&ImGuiWindowFlags_ChildWindow) != 0 &&
			(new_nav_window.Flags&(ImGuiWindowFlags_Popup|ImGuiWindowFlags_ChildMenu)) == 0 {
			new_nav_window = new_nav_window.ParentWindow
		}
		if new_nav_window != guiContext.NavWindow {
			var old_nav_window = guiContext.NavWindow
			FocusWindow(new_nav_window)
			new_nav_window.NavLastChildNavWindow = old_nav_window
		}

		// Toggle layer
		var new_nav_layer ImGuiNavLayer
		if guiContext.NavWindow.DC.NavLayersActiveMask&(1<<ImGuiNavLayer_Menu) != 0 {
			new_nav_layer = (ImGuiNavLayer)((int)(guiContext.NavLayer ^ 1))
		} else {
			new_nav_layer = ImGuiNavLayer_Main
		}
		if new_nav_layer != guiContext.NavLayer {
			// Reinitialize navigation when entering menu bar with the Alt key (FIXME: could be a properly of the layer?)
			if new_nav_layer == ImGuiNavLayer_Menu {
				guiContext.NavWindow.NavLastIds[new_nav_layer] = 0
			}
			NavRestoreLayer(new_nav_layer)
		}
	}
}

func NavUpdate() {
	io := guiContext.IO

	io.WantSetMousePos = false

	// Set input source as Gamepad when buttons are pressed (as some features differs when used with Gamepad vs Keyboard)
	// (do it before we map Keyboard input!)
	var nav_keyboard_active = (io.ConfigFlags & ImGuiConfigFlags_NavEnableKeyboard) != 0
	var nav_gamepad_active = (io.ConfigFlags&ImGuiConfigFlags_NavEnableGamepad) != 0 && (io.BackendFlags&ImGuiBackendFlags_HasGamepad) != 0
	if nav_gamepad_active && guiContext.NavInputSource != ImGuiInputSource_Gamepad {
		if io.NavInputs[ImGuiNavInput_Activate] > 0.0 || io.NavInputs[ImGuiNavInput_Input] > 0.0 || io.NavInputs[ImGuiNavInput_Cancel] > 0.0 || io.NavInputs[ImGuiNavInput_Menu] > 0.0 ||
			io.NavInputs[ImGuiNavInput_DpadLeft] > 0.0 || io.NavInputs[ImGuiNavInput_DpadRight] > 0.0 || io.NavInputs[ImGuiNavInput_DpadUp] > 0.0 || io.NavInputs[ImGuiNavInput_DpadDown] > 0.0 {
			guiContext.NavInputSource = ImGuiInputSource_Gamepad
		}
	}

	// Update Keyboard.Nav inputs mapping
	if nav_keyboard_active {
		var NAV_MAP_KEY = func(key ImGuiKey, input ImGuiNavInput) {
			if IsKeyDown(io.KeyMap[key]) {
				io.NavInputs[input] = 1.0
				guiContext.NavInputSource = ImGuiInputSource_Keyboard
			}
		}
		NAV_MAP_KEY(ImGuiKey_Space, ImGuiNavInput_Activate)
		NAV_MAP_KEY(ImGuiKey_Enter, ImGuiNavInput_Input)
		NAV_MAP_KEY(ImGuiKey_Escape, ImGuiNavInput_Cancel)
		NAV_MAP_KEY(ImGuiKey_LeftArrow, ImGuiNavInput_KeyLeft_)
		NAV_MAP_KEY(ImGuiKey_RightArrow, ImGuiNavInput_KeyRight_)
		NAV_MAP_KEY(ImGuiKey_UpArrow, ImGuiNavInput_KeyUp_)
		NAV_MAP_KEY(ImGuiKey_DownArrow, ImGuiNavInput_KeyDown_)
		if io.KeyCtrl {
			io.NavInputs[ImGuiNavInput_TweakSlow] = 1.0
		}
		if io.KeyShift {
			io.NavInputs[ImGuiNavInput_TweakFast] = 1.0
		}
	}
	copy(io.NavInputsDownDurationPrev[:], io.NavInputsDownDuration[:])
	for i := range io.NavInputs {

		if io.NavInputs[i] > 0.0 {
			if io.NavInputsDownDuration[i] < 0.0 {
				io.NavInputsDownDuration[i] = 0
			} else {
				io.NavInputsDownDuration[i] += guiContext.IO.DeltaTime
			}
		} else {
			io.NavInputsDownDuration[i] = -1.0
		}
	}

	// Process navigation init request (select first/default focus)
	if guiContext.NavInitResultId != 0 {
		NavUpdateInitResult()
	}
	guiContext.NavInitRequest = false
	guiContext.NavInitRequestFromMove = false
	guiContext.NavInitResultId = 0
	guiContext.NavJustMovedToId = 0

	// Process navigation move request
	if guiContext.NavMoveSubmitted {
		NavMoveRequestApplyResult()
	}
	guiContext.NavMoveSubmitted = false
	guiContext.NavMoveScoringItems = false

	// Apply application mouse position movement, after we had a chance to process move request result.
	if guiContext.NavMousePosDirty && guiContext.NavIdIsAlive {
		// Set mouse position given our knowledge of the navigated item position from last frame
		if (io.ConfigFlags&ImGuiConfigFlags_NavEnableSetMousePos != 0) && (io.BackendFlags&ImGuiBackendFlags_HasSetMousePos != 0) {
			if !guiContext.NavDisableHighlight && guiContext.NavDisableMouseHover && guiContext.NavWindow != nil {
				p := NavCalcPreferredRefPos()
				io.MousePos = p
				io.MousePosPrev = p
				io.WantSetMousePos = true
				//IMGUI_DEBUG_LOG("SetMousePos: (%.1f,%.1f)\n", io.MousePos.x, io.MousePos.y);
			}
		}
		guiContext.NavMousePosDirty = false
	}
	guiContext.NavIdIsAlive = false
	guiContext.NavJustTabbedId = 0
	IM_ASSERT(guiContext.NavLayer == 0 || guiContext.NavLayer == 1)

	// Store our return window (for returning from Menu Layer to Main Layer) and clear it as soon as we step back in our own Layer 0
	if guiContext.NavWindow != nil {
		NavSaveLastChildNavWindowIntoParent(guiContext.NavWindow)
	}
	if guiContext.NavWindow != nil && guiContext.NavWindow.NavLastChildNavWindow != nil && guiContext.NavLayer == ImGuiNavLayer_Main {
		guiContext.NavWindow.NavLastChildNavWindow = nil
	}

	// Update CTRL+TAB and Windowing features (hold Square to move/resize/etc.)
	NavUpdateWindowing()

	// Set output flags for user application
	io.NavActive = (nav_keyboard_active || nav_gamepad_active) && guiContext.NavWindow != nil && guiContext.NavWindow.Flags&ImGuiWindowFlags_NoNavInputs == 0
	io.NavVisible = (io.NavActive && guiContext.NavId != 0 && !guiContext.NavDisableHighlight) || (guiContext.NavWindowingTarget != nil)

	// Process NavCancel input (to close a popup, get back to parent, clear focus)
	NavUpdateCancelRequest()

	// Process manual activation request
	guiContext.NavActivateId = 0
	guiContext.NavActivateDownId = 0
	guiContext.NavActivatePressedId = 0
	guiContext.NavInputId = 0
	if guiContext.NavId != 0 && !guiContext.NavDisableHighlight && guiContext.NavWindowingTarget == nil && guiContext.NavWindow != nil && guiContext.NavWindow.Flags&ImGuiWindowFlags_NoNavInputs == 0 {
		var activate_down = IsNavInputDown(ImGuiNavInput_Activate)
		var activate_pressed = activate_down && IsNavInputTest(ImGuiNavInput_Activate, ImGuiInputReadMode_Pressed)
		if guiContext.ActiveId == 0 && activate_pressed {
			guiContext.NavActivateId = guiContext.NavId
		}
		if (guiContext.ActiveId == 0 || guiContext.ActiveId == guiContext.NavId) && activate_down {
			guiContext.NavActivateDownId = guiContext.NavId
		}
		if (guiContext.ActiveId == 0 || guiContext.ActiveId == guiContext.NavId) && activate_pressed {
			guiContext.NavActivatePressedId = guiContext.NavId
		}
		if (guiContext.ActiveId == 0 || guiContext.ActiveId == guiContext.NavId) && IsNavInputTest(ImGuiNavInput_Input, ImGuiInputReadMode_Pressed) {
			guiContext.NavInputId = guiContext.NavId
		}
	}
	if guiContext.NavWindow != nil && (guiContext.NavWindow.Flags&ImGuiWindowFlags_NoNavInputs != 0) {
		guiContext.NavDisableHighlight = true
	}
	if guiContext.NavActivateId != 0 {
		IM_ASSERT(guiContext.NavActivateDownId == guiContext.NavActivateId)
	}

	// Process programmatic activation request
	if guiContext.NavNextActivateId != 0 {
		guiContext.NavActivateId = guiContext.NavNextActivateId
		guiContext.NavActivateDownId = guiContext.NavNextActivateId
		guiContext.NavActivatePressedId = guiContext.NavNextActivateId
		guiContext.NavInputId = guiContext.NavNextActivateId
	}
	guiContext.NavNextActivateId = 0

	// Process move requests
	NavUpdateCreateMoveRequest()
	NavUpdateAnyRequestFlag()

	// Scrolling
	if guiContext.NavWindow != nil && (guiContext.NavWindow.Flags&ImGuiWindowFlags_NoNavInputs == 0) && guiContext.NavWindowingTarget == nil {
		// *Fallback* manual-scroll with Nav directional keys when window has no navigable item
		var window = guiContext.NavWindow
		var scroll_speed = IM_ROUND(window.CalcFontSize() * 100 * io.DeltaTime) // We need round the scrolling speed because sub-pixel scroll isn't reliably supported.
		var move_dir = guiContext.NavMoveDir
		if window.DC.NavLayersActiveMask == 0x00 && window.DC.NavHasScroll && move_dir != ImGuiDir_None {
			if move_dir == ImGuiDir_Left || move_dir == ImGuiDir_Right {
				var dir float
				if move_dir == ImGuiDir_Left {
					dir = -1.0
				} else {
					dir = 1.0
				}
				setScrollX(window, ImFloor(window.Scroll.x+dir*scroll_speed))
			}
			if move_dir == ImGuiDir_Up || move_dir == ImGuiDir_Down {
				var dir float
				if move_dir == ImGuiDir_Up {
					dir = -1.0
				} else {
					dir = 1.0
				}
				setScrollY(window, ImFloor(window.Scroll.y+dir*scroll_speed))
			}
		}

		// *Normal* Manual scroll with NavScrollXXX keys
		// Next movement request will clamp the NavId reference rectangle to the visible area, so navigation will resume within those bounds.
		var scroll_dir = GetNavInputAmount2d(ImGuiNavDirSourceFlags_PadLStick, ImGuiInputReadMode_Down, 1.0/10.0, 10.0)
		if scroll_dir.x != 0.0 && window.ScrollbarX {
			setScrollX(window, ImFloor(window.Scroll.x+scroll_dir.x*scroll_speed))
		}
		if scroll_dir.y != 0.0 {
			setScrollY(window, ImFloor(window.Scroll.y+scroll_dir.y*scroll_speed))
		}
	}

	// Always prioritize mouse highlight if navigation is disabled
	if !nav_keyboard_active && !nav_gamepad_active {
		guiContext.NavDisableHighlight = true
		guiContext.NavDisableMouseHover = false
		guiContext.NavMousePosDirty = false
	}

	// [DEBUG]
	guiContext.NavScoringDebugCount = 0
}

// NavInitWindow This needs to be called before we submit any widget (aka in or before Begin)
func NavInitWindow(window *ImGuiWindow, force_reinit bool) {
	IM_ASSERT(window == guiContext.NavWindow)

	if window.Flags&ImGuiWindowFlags_NoNavInputs != 0 {
		guiContext.NavId = 0
		guiContext.NavFocusScopeId = 0
		return
	}

	var init_for_nav = false
	if window == window.RootWindow || (window.Flags&ImGuiWindowFlags_Popup != 0) || (window.NavLastIds[0] == 0) || force_reinit {
		init_for_nav = true
	}

	//IMGUI_DEBUG_LOG_NAV("[nav] NavInitRequest: from NavInitWindow(), init_for_nav=%d, window=\"%s\", layer=%d\n", init_for_nav, window.Name, guiContext.NavLayer)
	if init_for_nav {
		SetNavID(0, guiContext.NavLayer, 0, &ImRect{})
		guiContext.NavInitRequest = true
		guiContext.NavInitRequestFromMove = false
		guiContext.NavInitResultId = 0
		guiContext.NavInitResultRectRel = ImRect{}
		NavUpdateAnyRequestFlag()
	} else {
		guiContext.NavId = window.NavLastIds[0]
		guiContext.NavFocusScopeId = 0
	}
}

func NavScoreItemDistInterval(a0, a1, b0, b1 float) float {
	if a1 < b0 {
		return a1 - b0
	}
	if b1 < a0 {
		return a0 - b1
	}
	return 0.0
}

// NavScoreItem Scoring function for gamepad/keyboard directional navigation. Based on https://gist.github.com/rygorous/6981057
func NavScoreItem(result *ImGuiNavItemData) bool {
	window := guiContext.CurrentWindow
	if guiContext.NavLayer != window.DC.NavLayerCurrent {
		return false
	}

	// FIXME: Those are not good variables names
	var cand = guiContext.LastItemData.NavRect // Current item nav rectangle
	var curr = guiContext.NavScoringRect       // Current modified source rect (NB: we've applied Max.x = Min.x in NavUpdate() to inhibit the effect of having varied item width)
	guiContext.NavScoringDebugCount++

	// When entering through a NavFlattened border, we consider child window items as fully clipped for scoring
	if window.ParentWindow == guiContext.NavWindow {
		IM_ASSERT((window.Flags|guiContext.NavWindow.Flags)&ImGuiWindowFlags_NavFlattened != 0)
		if !window.ClipRect.Overlaps(cand) {
			return false
		}
		cand.ClipWithFull(window.ClipRect) // This allows the scored item to not overlap other candidates in the parent window
	}

	// We perform scoring on items bounding box clipped by the current clipping rectangle on the other axis (clipping on our movement axis would give us equal scores for all clipped items)
	// For example, this ensure that items in one column are not reached when moving vertically from items in another column.
	NavClampRectToVisibleAreaForMoveDir(guiContext.NavMoveClipDir, &cand, &window.ClipRect)

	// Compute distance between boxes
	// FIXME-NAV: Introducing biases for vertical navigation, needs to be removed.
	var dbx = NavScoreItemDistInterval(cand.Min.x, cand.Max.x, curr.Min.x, curr.Max.x)
	var dby = NavScoreItemDistInterval(ImLerp(cand.Min.y, cand.Max.y, 0.2), ImLerp(cand.Min.y, cand.Max.y, 0.8), ImLerp(curr.Min.y, curr.Max.y, 0.2), ImLerp(curr.Min.y, curr.Max.y, 0.8)) // Scale down on Y to keep using box-distance for vertically touching items
	if dby != 0.0 && dbx != 0.0 {
		var dir float
		if dbx > 0.0 {
			dir = 1.0
		} else {
			dir = -1.0
		}
		dbx = (dbx / 1000.0) + dir
	}
	var dist_box = ImFabs(dbx) + ImFabs(dby)

	// Compute distance between centers (this is off by a factor of 2, but we only compare center distances with each other so it doesn't matter)
	var dcx = (cand.Min.x + cand.Max.x) - (curr.Min.x + curr.Max.x)
	var dcy = (cand.Min.y + cand.Max.y) - (curr.Min.y + curr.Max.y)
	var dist_center = ImFabs(dcx) + ImFabs(dcy) // L1 metric (need this for our connectedness guarantee)

	// Determine which quadrant of 'curr' our candidate item 'cand' lies in based on distance
	var quadrant ImGuiDir
	var dax, day, dist_axial float
	if dbx != 0.0 || dby != 0.0 {
		// For non-overlapping boxes, use distance between boxes
		dax = dbx
		day = dby
		dist_axial = dist_box
		quadrant = ImGetDirQuadrantFromDelta(dbx, dby)
	} else if dcx != 0.0 || dcy != 0.0 {
		// For overlapping boxes with different centers, use distance between centers
		dax = dcx
		day = dcy
		dist_axial = dist_center
		quadrant = ImGetDirQuadrantFromDelta(dcx, dcy)
	} else {
		// Degenerate case: two overlapping buttons with same center, break ties arbitrarily (note that LastItemId here is really the _previous_ item order, but it doesn't matter)
		if guiContext.LastItemData.ID < guiContext.NavId {
			quadrant = ImGuiDir_Left
		} else {
			quadrant = ImGuiDir_Right
		}
	}

	// Is it in the quadrant we're interesting in moving to?
	var new_best = false
	var move_dir = guiContext.NavMoveDir
	if quadrant == move_dir {
		// Does it beat the current best candidate?
		if dist_box < result.DistBox {
			result.DistBox = dist_box
			result.DistCenter = dist_center
			return true
		}
		if dist_box == result.DistBox {
			// Try using distance between center points to break ties
			if dist_center < result.DistCenter {
				result.DistCenter = dist_center
				new_best = true
			} else if dist_center == result.DistCenter {
				var check float
				if move_dir == ImGuiDir_Up || move_dir == ImGuiDir_Down {
					check = dby
				} else {
					check = dbx
				}
				// Still tied! we need to be extra-careful to make sure everything gets linked properly. We consistently break ties by symbolically moving "later" items
				// (with higher index) to the right/downwards by an infinitesimal amount since we the current "best" button already (so it must have a lower index),
				// this is fairly easy. This rule ensures that all buttons with dx==dy==0 will end up being linked in order of appearance along the x axis.
				if check < 0.0 { // moving bj to the right/down decreases distance
					new_best = true
				}
			}
		}
	}

	// Axial check: if 'curr' has no link at all in some direction and 'cand' lies roughly in that direction, add a tentative link. This will only be kept if no "real" matches
	// are found, so it only augments the graph produced by the above method using extra links. (important, since it doesn't guarantee strong connectedness)
	// This is just to avoid buttons having no links in a particular direction when there's a suitable neighbor. you get good graphs without this too.
	// 2017/09/29: FIXME: This now currently only enabled inside menu bars, ideally we'd disable it everywhere. Menus in particular need to catch failure. For general navigation it feels awkward.
	// Disabling it may lead to disconnected graphs when nodes are very spaced out on different axis. Perhaps consider offering this as an option?
	if result.DistBox == FLT_MAX && dist_axial < result.DistAxial { // Check axial match
		if guiContext.NavLayer == ImGuiNavLayer_Menu && guiContext.NavWindow.Flags&ImGuiWindowFlags_ChildMenu == 0 {
			if (move_dir == ImGuiDir_Left && dax < 0.0) || (move_dir == ImGuiDir_Right && dax > 0.0) || (move_dir == ImGuiDir_Up && day < 0.0) || (move_dir == ImGuiDir_Down && day > 0.0) {
				result.DistAxial = dist_axial
				new_best = true
			}
		}
	}

	return new_best
}

func NavClampRectToVisibleAreaForMoveDir(move_dir ImGuiDir, r *ImRect, clip_rect *ImRect) {
	if move_dir == ImGuiDir_Left || move_dir == ImGuiDir_Right {
		r.Min.y = ImClamp(r.Min.y, clip_rect.Min.y, clip_rect.Max.y)
		r.Max.y = ImClamp(r.Max.y, clip_rect.Min.y, clip_rect.Max.y)
	} else { // FIXME: PageUp/PageDown are leaving move_dir == None
		r.Min.x = ImClamp(r.Min.x, clip_rect.Min.x, clip_rect.Max.x)
		r.Max.x = ImClamp(r.Max.x, clip_rect.Min.x, clip_rect.Max.x)
	}
}

func NavEndFrame() {
	g := guiContext

	// Show CTRL+TAB list window
	if g.NavWindowingTarget != nil {
		NavUpdateWindowingOverlay()
	}

	// Perform wrap-around in menus
	// FIXME-NAV: Wrap (not Loop) support could be handled by the scoring function and then WrapX would function without an extra frame.
	var window = g.NavWindow
	var move_flags = g.NavMoveFlags
	var wanted_flags = ImGuiNavMoveFlags_WrapX | ImGuiNavMoveFlags_LoopX | ImGuiNavMoveFlags_WrapY | ImGuiNavMoveFlags_LoopY
	if window != nil && NavMoveRequestButNoResultYet() && (g.NavMoveFlags&wanted_flags != 0) && (g.NavMoveFlags&ImGuiNavMoveFlags_Forwarded) == 0 {
		var do_forward = false
		var bb_rel = window.NavRectRel[g.NavLayer]
		var clip_dir = g.NavMoveDir
		if g.NavMoveDir == ImGuiDir_Left && (move_flags&(ImGuiNavMoveFlags_WrapX|ImGuiNavMoveFlags_LoopX) != 0) {
			m := max(window.SizeFull.x, window.ContentSize.x+window.WindowPadding.x*2.0) - window.Scroll.x
			bb_rel.Min.x = m
			bb_rel.Max.x = m

			if move_flags&ImGuiNavMoveFlags_WrapX != 0 {
				bb_rel.TranslateY(-bb_rel.GetHeight())
				clip_dir = ImGuiDir_Up
			}
			do_forward = true
		}
		if g.NavMoveDir == ImGuiDir_Right && (move_flags&(ImGuiNavMoveFlags_WrapX|ImGuiNavMoveFlags_LoopX) != 0) {
			bb_rel.Min.x = -window.Scroll.x
			bb_rel.Max.x = -window.Scroll.x
			if move_flags&ImGuiNavMoveFlags_WrapX != 0 {
				bb_rel.TranslateY(+bb_rel.GetHeight())
				clip_dir = ImGuiDir_Down
			}
			do_forward = true
		}
		if g.NavMoveDir == ImGuiDir_Up && (move_flags&(ImGuiNavMoveFlags_WrapY|ImGuiNavMoveFlags_LoopY) != 0) {
			m := max(window.SizeFull.y, window.ContentSize.y+window.WindowPadding.y*2.0) - window.Scroll.y
			bb_rel.Min.y = m
			bb_rel.Max.y = m
			if move_flags&ImGuiNavMoveFlags_WrapY != 0 {
				bb_rel.TranslateX(-bb_rel.GetWidth())
				clip_dir = ImGuiDir_Left
			}
			do_forward = true
		}
		if g.NavMoveDir == ImGuiDir_Down && (move_flags&(ImGuiNavMoveFlags_WrapY|ImGuiNavMoveFlags_LoopY) != 0) {
			bb_rel.Min.y = -window.Scroll.y
			bb_rel.Max.y = -window.Scroll.y
			if move_flags&ImGuiNavMoveFlags_WrapY != 0 {
				bb_rel.TranslateX(+bb_rel.GetWidth())
				clip_dir = ImGuiDir_Right
			}
			do_forward = true
		}
		if do_forward {
			window.NavRectRel[g.NavLayer] = bb_rel
			NavMoveRequestForward(g.NavMoveDir, clip_dir, move_flags)
		}
	}
}

// NavUpdateWindowingOverlay Overlay displayed when using CTRL+TAB. Called by EndFrame().
func NavUpdateWindowingOverlay() {
	IM_ASSERT(guiContext.NavWindowingTarget != nil)

	if guiContext.NavWindowingTimer < NAV_WINDOWING_LIST_APPEAR_DELAY {
		return
	}

	if guiContext.NavWindowingListWindow == nil {
		guiContext.NavWindowingListWindow = FindWindowByName("###NavWindowingList")
	}
	var viewport = GetMainViewport()
	SetNextWindowSizeConstraints(ImVec2{viewport.Size.x * 0.20, viewport.Size.y * 0.20}, ImVec2{FLT_MAX, FLT_MAX}, nil, nil)
	center := viewport.GetCenter()
	SetNextWindowPos(&center, ImGuiCond_Always, ImVec2{0.5, 0.5})
	PushStyleVec(ImGuiStyleVar_WindowPadding, guiContext.Style.WindowPadding.Scale(2.0))
	Begin("###NavWindowingList", nil, ImGuiWindowFlags_NoTitleBar|ImGuiWindowFlags_NoFocusOnAppearing|ImGuiWindowFlags_NoResize|ImGuiWindowFlags_NoMove|ImGuiWindowFlags_NoInputs|ImGuiWindowFlags_AlwaysAutoResize|ImGuiWindowFlags_NoSavedSettings)
	for n := range guiContext.WindowsFocusOrder {
		var window = guiContext.WindowsFocusOrder[n]
		IM_ASSERT(window != nil) // Fix static analyzers
		if !IsWindowNavFocusable(window) {
			continue
		}
		var label = window.Name
		if label == FindRenderedTextEnd(label) {
			label = GetFallbackWindowNameForWindowingList(window)
		}
		Selectable(label, guiContext.NavWindowingTarget == window, 0, ImVec2{})
	}
	End()
	PopStyleVar(1)
}

// NavMoveRequestSubmit FIXME: ScoringRect is not set
func NavMoveRequestSubmit(move_dir ImGuiDir, clip_dir ImGuiDir, move_flags ImGuiNavMoveFlags) {
	IM_ASSERT(guiContext.NavWindow != nil)
	guiContext.NavMoveSubmitted = true
	guiContext.NavMoveScoringItems = true
	guiContext.NavMoveDir = move_dir
	guiContext.NavMoveDirForDebug = move_dir
	guiContext.NavMoveClipDir = clip_dir
	guiContext.NavMoveFlags = move_flags
	guiContext.NavMoveForwardToNextFrame = false
	guiContext.NavMoveKeyMods = guiContext.IO.KeyMods
	guiContext.NavMoveResultLocal.Clear()
	guiContext.NavMoveResultLocalVisible.Clear()
	guiContext.NavMoveResultOther.Clear()
}

func GetNavInputAmount2d(dir_sources ImGuiNavDirSourceFlags, mode ImGuiInputReadMode, slow_factor float, fast_factor float) ImVec2 {
	var delta ImVec2
	if dir_sources&ImGuiNavDirSourceFlags_Keyboard != 0 {
		delta = delta.Add(ImVec2{GetNavInputAmount(ImGuiNavInput_KeyRight_, mode) - GetNavInputAmount(ImGuiNavInput_KeyLeft_, mode), GetNavInputAmount(ImGuiNavInput_KeyDown_, mode) - GetNavInputAmount(ImGuiNavInput_KeyUp_, mode)})
	}
	if dir_sources&ImGuiNavDirSourceFlags_PadDPad != 0 {
		delta = delta.Add(ImVec2{GetNavInputAmount(ImGuiNavInput_DpadRight, mode) - GetNavInputAmount(ImGuiNavInput_DpadLeft, mode), GetNavInputAmount(ImGuiNavInput_DpadDown, mode) - GetNavInputAmount(ImGuiNavInput_DpadUp, mode)})
	}
	if dir_sources&ImGuiNavDirSourceFlags_PadLStick != 0 {
		delta = delta.Add(ImVec2{GetNavInputAmount(ImGuiNavInput_LStickRight, mode) - GetNavInputAmount(ImGuiNavInput_LStickLeft, mode), GetNavInputAmount(ImGuiNavInput_LStickDown, mode) - GetNavInputAmount(ImGuiNavInput_LStickUp, mode)})
	}
	if slow_factor != 0.0 && IsNavInputDown(ImGuiNavInput_TweakSlow) {
		delta = delta.Scale(slow_factor)
	}
	if fast_factor != 0.0 && IsNavInputDown(ImGuiNavInput_TweakFast) {
		delta = delta.Scale(fast_factor)
	}
	return delta
}

// NavMoveRequestApplyResult Apply result from previous frame navigation directional move request. Always called from NavUpdate()
func NavMoveRequestApplyResult() {
	g := guiContext

	// No result
	// In a situation when there is no results but NavId != 0, re-enable the Navigation highlight (because guiContext.NavId is not considered as a possible result)
	if g.NavMoveResultLocal.ID == 0 && g.NavMoveResultOther.ID == 0 {
		if g.NavId != 0 {
			g.NavDisableHighlight = false
			g.NavDisableMouseHover = true
		}
		return
	}

	// Select which result to use
	var result *ImGuiNavItemData
	if g.NavMoveResultLocal.ID != 0 {
		result = &g.NavMoveResultLocal
	} else {
		result = &g.NavMoveResultOther
	}

	// PageUp/PageDown behavior first jumps to the bottom/top mostly visible item, _otherwise_ use the result from the previous/next page.
	if g.NavMoveFlags&ImGuiNavMoveFlags_AlsoScoreVisibleSet != 0 {
		if g.NavMoveResultLocalVisible.ID != 0 && g.NavMoveResultLocalVisible.ID != g.NavId {
			result = &g.NavMoveResultLocalVisible
		}
	}

	// Maybe entering a flattened child from the outside? In this case solve the tie using the regular scoring rules.
	if result != &g.NavMoveResultOther && g.NavMoveResultOther.ID != 0 && g.NavMoveResultOther.Window.ParentWindow == g.NavWindow {
		if (g.NavMoveResultOther.DistBox < result.DistBox) || (g.NavMoveResultOther.DistBox == result.DistBox && g.NavMoveResultOther.DistCenter < result.DistCenter) {
			result = &g.NavMoveResultOther
		}
	}
	IM_ASSERT(g.NavWindow != nil && result.Window != nil)

	// Scroll to keep newly navigated item fully into view.
	if g.NavLayer == ImGuiNavLayer_Main {
		var delta_scroll ImVec2
		if g.NavMoveFlags&ImGuiNavMoveFlags_ScrollToEdge != 0 {
			var scroll_target float
			if g.NavMoveDir == ImGuiDir_Up {
				scroll_target = result.Window.ScrollMax.y
			}
			delta_scroll.y = result.Window.Scroll.y - scroll_target
			setScrollY(result.Window, scroll_target)
		} else {
			var rect_abs = ImRect{result.RectRel.Min.Add(result.Window.Pos), result.RectRel.Max.Add(result.Window.Pos)}
			delta_scroll = ScrollToBringRectIntoView(result.Window, &rect_abs)
		}

		// Offset our result position so mouse position can be applied immediately after in NavUpdate()
		result.RectRel.TranslateX(-delta_scroll.x)
		result.RectRel.TranslateY(-delta_scroll.y)
	}

	ClearActiveID()
	g.NavWindow = result.Window
	if g.NavId != result.ID {
		// Don't set NavJustMovedToId if just landed on the same spot (which may happen with ImGuiNavMoveFlags_AllowCurrentNavId)
		g.NavJustMovedToId = result.ID
		g.NavJustMovedToFocusScopeId = result.FocusScopeId
		g.NavJustMovedToKeyMods = g.NavMoveKeyMods
	}

	// Focus
	//IMGUI_DEBUG_LOG_NAV("[nav] NavMoveRequest: result NavID 0x%08X in Layer %d Window \"%s\"\n", result.ID, guiContext.NavLayer, guiContext.NavWindow.Name)
	SetNavID(result.ID, g.NavLayer, result.FocusScopeId, &result.RectRel)

	// Enable nav highlight
	g.NavDisableHighlight = false
	g.NavDisableMouseHover = true
	g.NavMousePosDirty = true
}

func NavMoveRequestCancel() {
	guiContext.NavMoveSubmitted = false
	guiContext.NavMoveScoringItems = false
	NavUpdateAnyRequestFlag()
}
