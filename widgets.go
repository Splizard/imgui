package imgui

const DRAGDROP_HOLD_TO_OPEN_TIMER float = 0.70 // Time for drag-hold to activate items accepting the ImGuiButtonFlags_PressedOnDragDropHold button behavior.
const DRAG_MOUSE_THRESHOLD_FACTOR float = 0.50 // Multiplier for the default value of io.MouseDragThreshold to make DragFloat/DragInt react faster to mouse drags.

func IsClippedEx(bb *ImRect, id ImGuiID, clip_even_when_logged bool) bool {
	window := g.CurrentWindow
	if !bb.Overlaps(window.ClipRect) {
		if id == 0 || (id != g.ActiveId && id != g.NavId) {
			if clip_even_when_logged || !g.LogEnabled {
				return true
			}
		}
	}
	return false
}

// is the last item active? (e.g. button being held, text field being edited. This will continuously return true while holding mouse button on an item. Items that don't interact will always return false)
func IsItemActive() bool {
	if g.ActiveId != 0 {
		return g.ActiveId == g.LastItemData.ID
	}
	return false
}

// Internal facing ItemHoverable() used when submitting widgets. Differs slightly from IsItemHovered().
func ItemHoverable(bb *ImRect, id ImGuiID) bool {
	if g.HoveredId != 0 && g.HoveredId != id && !g.HoveredIdAllowOverlap {
		return false
	}

	window := g.CurrentWindow
	if g.HoveredWindow != window {
		return false
	}
	if g.ActiveId != 0 && g.ActiveId != id && !g.ActiveIdAllowOverlap {
		return false
	}
	if !IsMouseHoveringRect(bb.Min, bb.Max, true) {
		return false
	}
	if g.NavDisableMouseHover {
		return false
	}
	if !IsWindowContentHoverable(window, ImGuiHoveredFlags_None) {
		g.HoveredIdDisabled = true
		return false
	}

	// We exceptionally allow this function to be called with id==0 to allow using it for easy high-level
	// hover test in widgets code. We could also decide to split this function is two.
	if id != 0 {
		SetHoveredID(id)
	}

	// When disabled we'll return false but still set HoveredId
	var item_flags = g.CurrentItemFlags
	if g.LastItemData.ID == id {
		item_flags = g.LastItemData.InFlags
	}
	if item_flags&ImGuiItemFlags_Disabled != 0 {
		// Release active id if turning disabled
		if g.ActiveId == id {
			ClearActiveID()
		}
		g.HoveredIdDisabled = true
		return false
	}

	/*if id != 0 {
		// [DEBUG] Item Picker tool!
		// We perform the check here because SetHoveredID() is not frequently called (1~ time a frame), making
		// the cost of this tool near-zero. We can get slightly better call-stack and support picking non-hovered
		// items if we perform the test in ItemAdd(), but that would incur a small runtime cost.
		// #define IMGUI_DEBUG_TOOL_ITEM_PICKER_EX in imconfig.h if you want this check to also be performed in ItemAdd().
		if g.DebugItemPickerActive && g.HoveredIdPreviousFrame == id {
			GetForegroundDrawList().AddRect(bb.Min, bb.Max, IM_COL32(255, 255, 0, 255))
		}
		if g.DebugItemPickerBreakId == id {
			IM_DEBUG_BREAK()
		}
	}*/

	return true
}

// [Internal] Calculate full item size given user provided 'size' parameter and default width/height. Default width is often == CalcItemWidth().
// Those two functions CalcItemWidth vs CalcItemSize are awkwardly named because they are not fully symmetrical.
// Note that only CalcItemWidth() is publicly exposed.
// The 4.0f here may be changed to match CalcItemWidth() and/or BeginChild() (right now we have a mismatch which is harmless but undesirable)
func CalcItemSize(size ImVec2, default_w float, default_h float) ImVec2 {
	window := g.CurrentWindow

	var region_max ImVec2
	if size.x < 0.0 || size.y < 0.0 {
		region_max = GetContentRegionMaxAbs()
	}

	if size.x == 0.0 {
		size.x = default_w
	} else if size.x < 0.0 {
		size.x = max(4.0, region_max.x-window.DC.CursorPos.x+size.x)
	}

	if size.y == 0.0 {
		size.y = default_h
	} else if size.y < 0.0 {
		size.y = max(4.0, region_max.y-window.DC.CursorPos.y+size.y)
	}

	return size
}

// Declare item bounding box for clipping and interaction.
// Note that the size can be different than the one provided to ItemSize(). Typically, widgets that spread over available surface
// declare their minimum size requirement to ItemSize() and provide a larger region to ItemAdd() which is used drawing/interaction.
func ItemAdd(bb *ImRect, id ImGuiID, nav_bb_arg *ImRect, extra_flags ImGuiItemFlags) bool {
	window := g.CurrentWindow

	// Set item data
	// (DisplayRect is left untouched, made valid when ImGuiItemStatusFlags_HasDisplayRect is set)
	g.LastItemData.ID = id
	g.LastItemData.Rect = *bb
	if nav_bb_arg != nil {
		g.LastItemData.NavRect = *nav_bb_arg
	} else {
		g.LastItemData.NavRect = *bb
	}
	g.LastItemData.InFlags = g.CurrentItemFlags | extra_flags
	g.LastItemData.StatusFlags = ImGuiItemStatusFlags_None

	// Directional navigation processing
	if id != 0 {
		// Runs prior to clipping early-out
		//  (a) So that NavInitRequest can be honored, for newly opened windows to select a default widget
		//  (b) So that we can scroll up/down past clipped items. This adds a small O(N) cost to regular navigation requests
		//      unfortunately, but it is still limited to one window. It may not scale very well for windows with ten of
		//      thousands of item, but at least NavMoveRequest is only set on user interaction, aka maximum once a frame.
		//      We could early out with "if (is_clipped && !g.NavInitRequest) return false;" but when we wouldn't be able
		//      to reach unclipped widgets. This would work if user had explicit scrolling control (e.g. mapped on a stick).
		// We intentionally don't check if g.NavWindow != nil because g.NavAnyRequest should only be set when it is non nil.
		// If we crash on a nil g.NavWindow we need to fix the bug elsewhere.
		window.DC.NavLayersActiveMaskNext |= (1 << window.DC.NavLayerCurrent)
		if g.NavId == id || g.NavAnyRequest {
			if g.NavWindow.RootWindowForNav == window.RootWindowForNav {
				if window == g.NavWindow || ((window.Flags|g.NavWindow.Flags)&ImGuiWindowFlags_NavFlattened != 0) {
					NavProcessItem()
				}
			}
		}
	}
	// Clipping test
	var is_clipped = IsClippedEx(bb, id, false)
	if is_clipped {
		return false
	}
	//if (g.IO.KeyAlt) window.DrawList.AddRect(bb.Min, bb.Max, IM_COL32(255,255,0,120)); // [DEBUG]

	// [WIP] Tab stop handling (previously was using internal FocusableItemRegister() api)
	// FIXME-NAV: We would now want to move this before the clipping test, but this would require being able to scroll and currently this would mean an extra frame. (#4079, #343)
	if extra_flags&ImGuiItemFlags_Inputable != 0 {
		ItemInputable(window, id)
	}

	// We need to calculate this now to take account of the current clipping rectangle (as items like Selectable may change them)
	if IsMouseHoveringRect(bb.Min, bb.Max, true) {
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HoveredRect
	}
	return true
}

// Advance cursor given item size for layout.
// Register minimum needed size so it can extend the bounding box used for auto-fit calculation.
// See comments in ItemAdd() about how/why the size provided to ItemSize() vs ItemAdd() may often different.
func ItemSizeVec(size *ImVec2, text_baseline_y float) {
	window := g.CurrentWindow
	if window.SkipItems {
		return
	}

	// We increase the height in this function to accommodate for baseline offset.
	// In theory we should be offsetting the starting position (window.DC.CursorPos), that will be the topic of a larger refactor,
	// but since ItemSize() is not yet an API that moves the cursor (to handle e.g. wrapping) enlarging the height has the same effect.
	var offset_to_match_baseline_y float
	if text_baseline_y >= 0 {
		offset_to_match_baseline_y = max(0.0, window.DC.CurrLineTextBaseOffset-text_baseline_y)
	}
	var line_height = max(window.DC.CurrLineSize.y, size.y+offset_to_match_baseline_y)

	// Always align ourselves on pixel boundaries
	//if (g.IO.KeyAlt) window.DrawList.AddRect(window.DC.CursorPos, window.DC.CursorPos + ImVec2(size.x, line_height), IM_COL32(255,0,0,200)); // [DEBUG]
	window.DC.CursorPosPrevLine.x = window.DC.CursorPos.x + size.x
	window.DC.CursorPosPrevLine.y = window.DC.CursorPos.y
	window.DC.CursorPos.x = IM_FLOOR(window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x) // Next line
	window.DC.CursorPos.y = IM_FLOOR(window.DC.CursorPos.y + line_height + g.Style.ItemSpacing.y)   // Next line
	window.DC.CursorMaxPos.x = max(window.DC.CursorMaxPos.x, window.DC.CursorPosPrevLine.x)
	window.DC.CursorMaxPos.y = max(window.DC.CursorMaxPos.y, window.DC.CursorPos.y-g.Style.ItemSpacing.y)
	//if (g.IO.KeyAlt) window.DrawList.AddCircle(window.DC.CursorMaxPos, 3.0f, IM_COL32(255,0,0,255), 4); // [DEBUG]

	window.DC.PrevLineSize.y = line_height
	window.DC.CurrLineSize.y = 0.0
	window.DC.PrevLineTextBaseOffset = max(window.DC.CurrLineTextBaseOffset, text_baseline_y)
	window.DC.CurrLineTextBaseOffset = 0.0

	// Horizontal layout mode
	if window.DC.LayoutType == ImGuiLayoutType_Horizontal {
		SameLine(0, -1)
	}
}

// Gets back to previous line and continue with horizontal layout
//
//	offset_from_start_x == 0 : follow right after previous item
//	offset_from_start_x != 0 : align to specified x position (relative to window/group left)
//	spacing_w < 0            : use default spacing if pos_x == 0, no spacing if pos_x != 0
//	spacing_w >= 0           : enforce spacing amount
func SameLine(offset_from_start_x, spacing_w float) {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	if offset_from_start_x != 0.0 {
		if spacing_w < 0.0 {
			spacing_w = 0.0
		}
		window.DC.CursorPos.x = window.Pos.x - window.Scroll.x + offset_from_start_x + spacing_w + window.DC.GroupOffset.x + window.DC.ColumnsOffset.x
		window.DC.CursorPos.y = window.DC.CursorPosPrevLine.y
	} else {
		if spacing_w < 0.0 {
			spacing_w = g.Style.ItemSpacing.x
		}
		window.DC.CursorPos.x = window.DC.CursorPosPrevLine.x + spacing_w
		window.DC.CursorPos.y = window.DC.CursorPosPrevLine.y
	}
	window.DC.CurrLineSize = window.DC.PrevLineSize
	window.DC.CurrLineTextBaseOffset = window.DC.PrevLineTextBaseOffset
}

func ItemSizeRect(bb *ImRect, text_baseline_y float) {
	size := bb.GetSize()
	ItemSizeVec(&size, text_baseline_y)
}
