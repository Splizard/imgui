package imgui

// BeginDragDropTargetCustom Drag and Drop
func BeginDragDropTargetCustom(bb *ImRect, id ImGuiID) bool {
	if !g.DragDropActive {
		return false
	}

	window := g.CurrentWindow
	var hovered_window = g.HoveredWindowUnderMovingWindow
	if hovered_window == nil || window.RootWindow != hovered_window.RootWindow {
		return false
	}
	IM_ASSERT(id != 0)
	if !IsMouseHoveringRect(bb.Min, bb.Max, true) || (id == g.DragDropPayload.SourceId) {
		return false
	}
	if window.SkipItems {
		return false
	}

	IM_ASSERT(!g.DragDropWithinTarget)
	g.DragDropTargetRect = *bb
	g.DragDropTargetId = id
	g.DragDropWithinTarget = true
	return true
}

func ClearDragDrop() {
	g.DragDropActive = false
	g.DragDropPayload = ImGuiPayload{}
	g.DragDropAcceptFlags = ImGuiDragDropFlags_None
	g.DragDropAcceptIdCurr = 0
	g.DragDropAcceptIdPrev = 0
	g.DragDropAcceptIdCurrRectSurface = FLT_MAX
	g.DragDropAcceptFrameCount = -1

	g.DragDropPayloadBufHeap = g.DragDropPayloadBufHeap[:0]
	g.DragDropPayloadBufLocal = [len(g.DragDropPayloadBufLocal)]byte{}
}

func IsDragDropPayloadBeingAccepted() bool {
	return g.DragDropActive && g.DragDropAcceptIdPrev != 0
}

// Drag and Drop
// - On source items, call BeginDragDropSource(), if it returns true also call SetDragDropPayload() + EndDragDropSource().
// - On target candidates, call BeginDragDropTarget(), if it returns true also call AcceptDragDropPayload() + EndDragDropTarget().
// - If you stop calling BeginDragDropSource() the payload is preserved however it won't have a preview tooltip (we currently display a fallback "..." tooltip, see #1725)
// - An item can be both drag source and drop target.

// BeginDragDropSource When this returns true you need to: a) call SetDragDropPayload() exactly once, b) you may render the payload visual/description, c) call EndDragDropSource()
// If the item has an identifier:
// - This assume/require the item to be activated (typically via ButtonBehavior).
// - Therefore if you want to use this with a mouse button other than left mouse button, it is up to the item itself to activate with another button.
// - We then pull and use the mouse button that was used to activate the item and use it to carry on the drag.
// If the item has no identifier:
// - Currently always assume left mouse button.
// call after submitting an item which may be dragged. when this return true, you can call SetDragDropPayload() + EndDragDropSource()
func BeginDragDropSource(flags ImGuiDragDropFlags) bool {
	window := g.CurrentWindow

	// FIXME-DRAGDROP: While in the common-most "drag from non-zero active id" case we can tell the mouse button,
	// in both SourceExtern and id==0 cases we may requires something else (explicit flags or some heuristic).
	var mouse_button = ImGuiMouseButton_Left

	var source_drag_active = false
	var source_id ImGuiID = 0
	var source_parent_id ImGuiID = 0
	if flags&ImGuiDragDropFlags_SourceExtern == 0 {
		source_id = g.LastItemData.ID
		if source_id != 0 {
			// Common path: items with ID
			if g.ActiveId != source_id {
				return false
			}
			if g.ActiveIdMouseButton != -1 {
				mouse_button = g.ActiveIdMouseButton
			}
			if !g.IO.MouseDown[mouse_button] {
				return false
			}
			g.ActiveIdAllowOverlap = false
		} else {
			// Uncommon path: items without ID
			if !g.IO.MouseDown[mouse_button] {
				return false
			}

			// If you want to use BeginDragDropSource() on an item with no unique identifier for interaction, such as Text() or Image(), you need to:
			// A) Read the explanation below, B) Use the ImGuiDragDropFlags_SourceAllownilID flag, C) Swallow your programmer pride.
			if flags&ImGuiDragDropFlags_SourceAllowNullID == 0 {
				IM_ASSERT(false)
				return false
			}

			// Early out
			if (g.LastItemData.StatusFlags&ImGuiItemStatusFlags_HoveredRect) == 0 && (g.ActiveId == 0 || g.ActiveIdWindow != window) {
				return false
			}

			// Magic fallback (=somehow reprehensible) to handle items with no assigned ID, e.g. Text(), Image()
			// We build a throwaway ID based on current ID stack + relative AABB of items in window.
			// THE IDENTIFIER WON'T SURVIVE ANY REPOSITIONING OF THE WIDGET, so if your widget moves your dragging operation will be canceled.
			// We don't need to maintain/call ClearActiveID() as releasing the button will early out this function and trigger !ActiveIdIsAlive.
			// Rely on keeping other window.LastItemXXX fields intact.
			source_id = window.GetIDFromRectangle(g.LastItemData.Rect)
			g.LastItemData.ID = source_id
			var is_hovered = ItemHoverable(&g.LastItemData.Rect, source_id)
			if is_hovered && g.IO.MouseClicked[mouse_button] {
				SetActiveID(source_id, window)
				FocusWindow(window)
			}
			if g.ActiveId == source_id { // Allow the underlying widget to display/return hovered during the mouse release frame, else we would get a flicker.
				g.ActiveIdAllowOverlap = is_hovered
			}
		}
		if g.ActiveId != source_id {
			return false
		}
		source_parent_id = window.IDStack[len(window.IDStack)-1]
		source_drag_active = IsMouseDragging(mouse_button, -1)

		// Disable navigation and key inputs while dragging + cancel existing request if any
		SetActiveIdUsingNavAndKeys()
	} else {
		window = nil
		source_id = ImHashStr("#SourceExtern", uintptr(len("#SourceExtern")), 0)
		source_drag_active = true
	}

	if source_drag_active {
		if !g.DragDropActive {
			IM_ASSERT(source_id != 0)
			ClearDragDrop()
			var payload = &g.DragDropPayload
			payload.SourceId = source_id
			payload.SourceParentId = source_parent_id
			g.DragDropActive = true
			g.DragDropSourceFlags = flags
			g.DragDropMouseButton = mouse_button
			if payload.SourceId == g.ActiveId {
				g.ActiveIdNoClearOnFocusLoss = true
			}
		}
		g.DragDropSourceFrameCount = g.FrameCount
		g.DragDropWithinSource = true

		if flags&ImGuiDragDropFlags_SourceNoPreviewTooltip == 0 {
			// Target can request the Source to not display its tooltip (we use a dedicated flag to make this request explicit)
			// We unfortunately can't just modify the source flags and skip the call to BeginTooltip, as caller may be emitting contents.
			BeginTooltip()
			if g.DragDropAcceptIdPrev != 0 && (g.DragDropAcceptFlags&ImGuiDragDropFlags_AcceptNoPreviewTooltip != 0) {
				var tooltip_window = g.CurrentWindow
				tooltip_window.Hidden = true
				tooltip_window.SkipItems = true
				tooltip_window.HiddenFramesCanSkipItems = 1
			}
		}

		if flags&ImGuiDragDropFlags_SourceNoDisableHover == 0 && flags&ImGuiDragDropFlags_SourceExtern == 0 {
			g.LastItemData.StatusFlags &= ^ImGuiItemStatusFlags_HoveredRect
		}

		return true
	}
	return false
}

// SetDragDropPayload Use 'cond' to choose to submit payload on drag start or every frame
// type is a user defined string of maximum 32 characters. Strings starting with '_' are reserved for dear imgui internal types. Data is copied and held by imgui.
func SetDragDropPayload(ptype string, data any, data_size uintptr, cond ImGuiCond) bool {
	var payload = &g.DragDropPayload
	if cond == 0 {
		cond = ImGuiCond_Always
	}

	IM_ASSERT(ptype != "")
	IM_ASSERT_USER_ERROR(len(ptype) < len(payload.DataType), "Payload type can be at most 32 characters long")
	IM_ASSERT((data != nil && data_size > 0) || (data == nil && data_size == 0))
	IM_ASSERT(cond == ImGuiCond_Always || cond == ImGuiCond_Once)
	IM_ASSERT(payload.SourceId != 0) // Not called between BeginDragDropSource() and EndDragDropSource()

	if cond == ImGuiCond_Always || payload.DataFrameCount == -1 {
		// Copy payload
		copy(payload.DataType[:], ptype[len(payload.DataType):])
		g.DragDropPayloadBufHeap = g.DragDropPayloadBufHeap[:0]
		payload.Data = data
		payload.DataSize = (int)(data_size)
	}
	payload.DataFrameCount = g.FrameCount

	return (g.DragDropAcceptFrameCount == g.FrameCount) || (g.DragDropAcceptFrameCount == g.FrameCount-1)
}

// EndDragDropSource only call EndDragDropSource() if BeginDragDropSource() returns true!
func EndDragDropSource() {
	IM_ASSERT(g.DragDropActive)
	IM_ASSERT_USER_ERROR(g.DragDropWithinSource, "Not after a BeginDragDropSource()?")

	if g.DragDropSourceFlags&ImGuiDragDropFlags_SourceNoPreviewTooltip != 0 {
		EndTooltip()
	}

	// Discard the drag if have not called SetDragDropPayload()
	if g.DragDropPayload.DataFrameCount == -1 {
		ClearDragDrop()
	}
	g.DragDropWithinSource = false
}

// BeginDragDropTarget call after submitting an item that may receive a payload. If this returns true, you can call AcceptDragDropPayload() + EndDragDropTarget()\
// We don't use BeginDragDropTargetCustom() and duplicate its code because:
// 1) we use LastItemRectHoveredRect which handles items that pushes a temporarily clip rectangle in their code. Calling BeginDragDropTargetCustom(LastItemRect) would not handle them.
// 2) and it's faster. as this code may be very frequently called, we want to early out as fast as we can.
// Also note how the HoveredWindow test is positioned differently in both functions (in both functions we optimize for the cheapest early out case)
func BeginDragDropTarget() bool {
	if !g.DragDropActive {
		return false
	}

	window := g.CurrentWindow
	if g.LastItemData.StatusFlags&ImGuiItemStatusFlags_HoveredRect == 0 {
		return false
	}
	var hovered_window = g.HoveredWindowUnderMovingWindow
	if hovered_window == nil || window.RootWindow != hovered_window.RootWindow {
		return false
	}

	var display_rect = g.LastItemData.Rect
	if g.LastItemData.StatusFlags&ImGuiItemStatusFlags_HasDisplayRect != 0 {
		display_rect = g.LastItemData.DisplayRect
	}
	var id = g.LastItemData.ID
	if id == 0 {
		id = window.GetIDFromRectangle(display_rect)
	}
	if g.DragDropPayload.SourceId == id {
		return false
	}

	IM_ASSERT(!g.DragDropWithinTarget)
	g.DragDropTargetRect = display_rect
	g.DragDropTargetId = id
	g.DragDropWithinTarget = true
	return true
}

// AcceptDragDropPayload accept contents of a given type. If ImGuiDragDropFlags_AcceptBeforeDelivery is set you can peek into the payload before the mouse button is released.
func AcceptDragDropPayload(ptype string, flags ImGuiDragDropFlags) *ImGuiPayload {
	window := g.CurrentWindow
	var payload = g.DragDropPayload
	IM_ASSERT(g.DragDropActive)             // Not called between BeginDragDropTarget() and EndDragDropTarget() ?
	IM_ASSERT(payload.DataFrameCount != -1) // Forgot to call EndDragDropTarget() ?
	if ptype != "" && !payload.IsDataType(ptype) {
		return nil
	}

	// Accept smallest drag target bounding box, this allows us to nest drag targets conveniently without ordering constraints.
	// NB: We currently accept nil id as target. However, overlapping targets requires a unique ID to function!
	var was_accepted_previously = (g.DragDropAcceptIdPrev == g.DragDropTargetId)
	var r = g.DragDropTargetRect
	var r_surface = r.GetWidth() * r.GetHeight()
	if r_surface <= g.DragDropAcceptIdCurrRectSurface {
		g.DragDropAcceptFlags = flags
		g.DragDropAcceptIdCurr = g.DragDropTargetId
		g.DragDropAcceptIdCurrRectSurface = r_surface
	}

	// Render default drop visuals
	// FIXME-DRAGDROP: Settle on a proper default visuals for drop target.
	payload.Preview = was_accepted_previously
	flags |= (g.DragDropSourceFlags & ImGuiDragDropFlags_AcceptNoDrawDefaultRect) // Source can also inhibit the preview (useful for external sources that lives for 1 frame)
	if (flags&ImGuiDragDropFlags_AcceptNoDrawDefaultRect == 0) && payload.Preview {
		window.DrawList.AddRect(r.Min.Sub(ImVec2{3.5, 3.5}), r.Max.Add(ImVec2{3.5, 3.5}), GetColorU32FromID(ImGuiCol_DragDropTarget, 1), 0.0, 0, 2.0)
	}

	g.DragDropAcceptFrameCount = g.FrameCount
	payload.Delivery = was_accepted_previously && !IsMouseDown(g.DragDropMouseButton) // For extern drag sources affecting os window focus, it's easier to just test !IsMouseDown() instead of IsMouseReleased()
	if !payload.Delivery && (flags&ImGuiDragDropFlags_AcceptBeforeDelivery == 0) {
		return nil
	}

	return &payload
}

// EndDragDropTarget We don't really use/need this now, but added it for the sake of consistency and because we might need it later.
// only call EndDragDropTarget() if BeginDragDropTarget() returns true!
func EndDragDropTarget() {
	IM_ASSERT(g.DragDropActive)
	IM_ASSERT(g.DragDropWithinTarget)
	g.DragDropWithinTarget = false
}

// GetDragDropPayload peek directly into the current payload from anywhere. may return NULL. use ImGuiPayload::IsDataType() to test for the payload type.
func GetDragDropPayload() *ImGuiPayload {
	if g.DragDropActive {
		return &g.DragDropPayload
	}
	return nil
}
