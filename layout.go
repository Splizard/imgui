package imgui

// Cursor / Layout
// - By "cursor" we mean the current output position.
// - The typical widget behavior is to output themselves at the current cursor position, then move the cursor one line down.
// - You can call SameLine() between widgets to undo the last carriage return and output at the right of the preceding widget.
// - Attention! We currently have inconsistencies between window-local and absolute positions we will aim to fix with future API:
//    Window-local coordinates:   SameLine(), GetCursorPos(), SetCursorPos(), GetCursorStartPos(), GetContentRegionMax(), GetWindowContentRegion*(), PushTextWrapPos()
//    Absolute coordinate:        GetCursorScreenPos(), SetCursorScreenPos(), all ImDrawList:: functions.
func NewLine() { panic("not implemented") } // undo a SameLine() or force a new line when in an horizontal-layout context.
func Spacing() { panic("not implemented") } // add vertical spacing.

// add a dummy item of given size. unlike InvisibleButton(), Dummy() won't take the mouse click or be navigable into.
func Dummy(size ImVec2) {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	ItemSizeVec(&size, 0)
	ItemAdd(&bb, 0, nil, 0)
}

// Lock horizontal starting position + capture group bounding box into one "item" (so you can use IsItemHovered() or layout primitives such as SameLine() on whole group, etc.)
// Groups are currently a mishmash of functionalities which should perhaps be clarified and separated.
func BeginGroup() {
	var g = GImGui
	var window = g.CurrentWindow

	g.GroupStack = append(g.GroupStack, ImGuiGroupData{})

	var group_data = &g.GroupStack[len(g.GroupStack)-1]
	group_data.WindowID = window.ID
	group_data.BackupCursorPos = window.DC.CursorPos
	group_data.BackupCursorMaxPos = window.DC.CursorMaxPos
	group_data.BackupIndent = window.DC.Indent
	group_data.BackupGroupOffset = window.DC.GroupOffset
	group_data.BackupCurrLineSize = window.DC.CurrLineSize
	group_data.BackupCurrLineTextBaseOffset = window.DC.CurrLineTextBaseOffset
	group_data.BackupActiveIdIsAlive = g.ActiveIdIsAlive
	group_data.BackupHoveredIdIsAlive = g.HoveredId != 0
	group_data.BackupActiveIdPreviousFrameIsAlive = g.ActiveIdPreviousFrameIsAlive
	group_data.EmitItem = true

	window.DC.GroupOffset.x = window.DC.CursorPos.x - window.Pos.x - window.DC.ColumnsOffset.x
	window.DC.Indent = window.DC.GroupOffset
	window.DC.CursorMaxPos = window.DC.CursorPos
	window.DC.CurrLineSize = ImVec2{0.0, 0.0}
	if g.LogEnabled {
		g.LogLinePosY = -FLT_MAX // To enforce a carriage return
	}
}

// unlock horizontal starting position + capture the whole group bounding box into one "item" (so you can use IsItemHovered() or layout primitives such as SameLine() on whole group, etc.)
func EndGroup() {
	var g = GImGui
	var window *ImGuiWindow = g.CurrentWindow
	IM_ASSERT(len(g.GroupStack) > 0) // Mismatched BeginGroup()/EndGroup() calls

	var group_data *ImGuiGroupData = &g.GroupStack[len(g.GroupStack)-1]
	IM_ASSERT(group_data.WindowID == window.ID) // EndGroup() in wrong window?

	var group_bb = ImRect{group_data.BackupCursorPos, ImMaxVec2(&window.DC.CursorMaxPos, &group_data.BackupCursorPos)}

	window.DC.CursorPos = group_data.BackupCursorPos
	window.DC.CursorMaxPos = ImMaxVec2(&group_data.BackupCursorMaxPos, &window.DC.CursorMaxPos)
	window.DC.Indent = group_data.BackupIndent
	window.DC.GroupOffset = group_data.BackupGroupOffset
	window.DC.CurrLineSize = group_data.BackupCurrLineSize
	window.DC.CurrLineTextBaseOffset = group_data.BackupCurrLineTextBaseOffset
	if g.LogEnabled {
		g.LogLinePosY = -FLT_MAX // To enforce a carriage return
	}

	if !group_data.EmitItem {
		g.GroupStack = g.GroupStack[:len(g.GroupStack)-1]
		return
	}

	window.DC.CurrLineTextBaseOffset = ImMax(window.DC.PrevLineTextBaseOffset, group_data.BackupCurrLineTextBaseOffset) // FIXME: Incorrect, we should grab the base offset from the *first line* of the group but it is hard to obtain now.

	size := group_bb.GetSize()
	ItemSizeVec(&size, 0)
	ItemAdd(&group_bb, 0, nil, 0)

	// If the current ActiveId was declared within the boundary of our group, we copy it to LastItemId so IsItemActive(), IsItemDeactivated() etc. will be functional on the entire group.
	// It would be be neater if we replaced window.DC.LastItemId by e.g. 'bool LastItemIsActive', but would put a little more burden on individual widgets.
	// Also if you grep for LastItemId you'll notice it is only used in that context.
	// (The two tests not the same because ActiveIdIsAlive is an ID itself, in order to be able to handle ActiveId being overwritten during the frame.)
	var group_contains_curr_active_id bool = (group_data.BackupActiveIdIsAlive != g.ActiveId) && (g.ActiveIdIsAlive == g.ActiveId) && g.ActiveId != 0
	var group_contains_prev_active_id bool = (group_data.BackupActiveIdPreviousFrameIsAlive == false) && (g.ActiveIdPreviousFrameIsAlive == true)
	if group_contains_curr_active_id {
		g.LastItemData.ID = g.ActiveId
	} else if group_contains_prev_active_id {
		g.LastItemData.ID = g.ActiveIdPreviousFrame
	}
	g.LastItemData.Rect = group_bb

	// Forward Hovered flag
	var group_contains_curr_hovered_id bool = (group_data.BackupHoveredIdIsAlive == false) && g.HoveredId != 0
	if group_contains_curr_hovered_id {
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HoveredWindow
	}

	// Forward Edited flag
	if group_contains_curr_active_id && g.ActiveIdHasBeenEditedThisFrame {
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_Edited
	}

	// Forward Deactivated flag
	g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HasDeactivated
	if group_contains_prev_active_id && g.ActiveId != g.ActiveIdPreviousFrame {
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_Deactivated
	}

	g.GroupStack = g.GroupStack[:len(g.GroupStack)-1]
	//window.DrawList.AddRect(group_bb.Min, group_bb.Max, IM_COL32(255,0,255,255));   // [Debug]
}

// cursor position in window coordinates (relative to window position)
func GetCursorPos() ImVec2 {

	// User generally sees positions in window coordinates. Internally we store CursorPos in absolute screen coordinates because it is more convenient.
	// Conversion happens as we pass the value to user, but it makes our naming convention confusing because GetCursorPos() == (DC.CursorPos - window.Pos). May want to rename 'DC.CursorPos'.

	var window = GetCurrentWindowRead()
	return window.DC.CursorPos.Sub(window.Pos).Add(window.Scroll)
}

//   (some functions are using window-relative coordinates, such as: GetCursorPos, GetCursorStartPos, GetContentRegionMax, GetWindowContentRegion* etc.
//    other functions such as GetCursorScreenPos or everything in ImDrawList::
//    are using the main, absolute coordinate system.
//    GetWindowPos() + GetCursorPos() == GetCursorScreenPos() etc.)

func GetCursorPosX() float {
	var window = GetCurrentWindowRead()
	return window.DC.CursorPos.x - window.Pos.x + window.Scroll.x
}

func GetCursorPosY() float {
	var window = GetCurrentWindowRead()
	return window.DC.CursorPos.y - window.Pos.y + window.Scroll.y
}

func SetCursorPos(local_pos *ImVec2) {
	var window = GetCurrentWindow()
	window.DC.CursorPos = window.Pos.Sub(window.Scroll).Add(*local_pos)
	window.DC.CursorMaxPos = ImMaxVec2(&window.DC.CursorMaxPos, &window.DC.CursorPos)
}

func SetCursorPosX(local_x float) {
	var window = GetCurrentWindow()
	window.DC.CursorPos.x = window.Pos.x - window.Scroll.x + local_x
	window.DC.CursorMaxPos.x = ImMax(window.DC.CursorMaxPos.x, window.DC.CursorPos.x)
}

func SetCursorPosY(local_y float) {
	var window = GetCurrentWindow()
	window.DC.CursorPos.y = window.Pos.y - window.Scroll.y + local_y
	window.DC.CursorMaxPos.y = ImMax(window.DC.CursorMaxPos.y, window.DC.CursorPos.y)
}

// initial cursor position in window coordinates
func GetCursorStartPos() ImVec2 {
	var window = GetCurrentWindowRead()
	return window.DC.CursorStartPos.Sub(window.Pos)
}

// cursor position in absolute coordinates (useful to work with ImDrawList API). generally top-left == GetMainViewport()->Pos == (0,0) in single viewport mode, and bottom-right == GetMainViewport()->Pos+Size == io.DisplaySize in single-viewport mode.
func GetCursorScreenPos() ImVec2 {
	var window = GetCurrentWindowRead()
	return window.DC.CursorPos
}

// cursor position in absolute coordinates
func SetCursorScreenPos(pos ImVec2) {
	var window = GetCurrentWindow()
	window.DC.CursorPos = pos
	window.DC.CursorMaxPos = ImMaxVec2(&window.DC.CursorMaxPos, &window.DC.CursorPos)
}

func AlignTextToFramePadding() { panic("not implemented") } // vertically align upcoming text baseline to FramePadding.y so that it will align properly to regularly framed items (call if you have text on a line before a framed item)

// ~ FontSize
func GetTextLineHeight() float {
	var g = GImGui
	return g.FontSize
}

// ~ FontSize + style.ItemSpacing.y (distance in pixels between 2 consecutive lines of text)
func GetTextLineHeightWithSpacing() float {
	var g = GImGui
	return g.FontSize + g.Style.ItemSpacing.y
}

// ~ FontSize + style.FramePadding.y * 2
func GetFrameHeight() float {
	var g = GImGui
	return g.FontSize + g.Style.FramePadding.y*2.0
}

// ~ FontSize + style.FramePadding.y * 2 + style.ItemSpacing.y (distance in pixels between 2 consecutive lines of framed widgets)
func GetFrameHeightWithSpacing() float {
	var g = GImGui
	return g.FontSize + g.Style.FramePadding.y*2.0 + g.Style.ItemSpacing.y
}

// Parameters stacks (current window)

// push width of items for common large "item+label" widgets. >0.0: width in pixels, <0.0 align xx pixels to the right of window (so -FLT_MIN always align width to the right side).
func PushItemWidth(item_width float) {
	// FIXME: Remove the == 0.0f behavior?

	var g = GImGui
	var window = g.CurrentWindow
	window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, window.DC.ItemWidth) // Backup current width
	if item_width == 0 {
		window.DC.ItemWidth = window.ItemWidthDefault
	} else {
		window.DC.ItemWidth = item_width
	}
	g.NextItemData.Flags &= ^ImGuiNextItemDataFlags_HasWidth
}

func PushMultiItemsWidths(components int, width_full float) {
	var g = GImGui
	var window = g.CurrentWindow
	var style = g.Style
	var w_item_one float = ImMax(1.0, IM_FLOOR((width_full-(style.ItemInnerSpacing.x)*float(components-1))/(float)(components)))
	var w_item_last float = ImMax(1.0, IM_FLOOR(width_full-(w_item_one+style.ItemInnerSpacing.x)*float(components-1)))
	window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, window.DC.ItemWidth) // Backup current width
	window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, w_item_last)
	for i := int(0); i < components-2; i++ {
		window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, w_item_one)
	}
	if components == 1 {
		window.DC.ItemWidth = w_item_last
	} else {
		window.DC.ItemWidth = w_item_one
	}
	g.NextItemData.Flags &= ^ImGuiNextItemDataFlags_HasWidth
}

func PopItemWidth() {
	var window = GetCurrentWindow()
	window.DC.ItemWidth = window.DC.ItemWidthStack[len(window.DC.ItemWidthStack)-1]
	window.DC.ItemWidthStack = window.DC.ItemWidthStack[:len(window.DC.ItemWidthStack)-1]
}

// set width of the _next_ common large "item+label" widget. >0.0: width in pixels, <0.0 align xx pixels to the right of window (so -FLT_MIN always align width to the right side)
// Affect large frame+labels widgets only.
func SetNextItemWidth(item_width float) {
	var g = GImGui
	g.NextItemData.Flags |= ImGuiNextItemDataFlags_HasWidth
	g.NextItemData.Width = item_width
}

// Calculate default item width given value passed to PushItemWidth() or SetNextItemWidth().
// The SetNextItemWidth() data is generally cleared/consumed by ItemAdd() or NextItemData.ClearFlags()
// width of item given pushed settings and current cursor position. NOT necessarily the width of last item unlike most 'Item' functions.
func CalcItemWidth() float {
	var g = GImGui
	var window = g.CurrentWindow
	var w float
	if g.NextItemData.Flags&ImGuiNextItemDataFlags_HasWidth != 0 {
		w = g.NextItemData.Width
	} else {
		w = window.DC.ItemWidth
	}
	if w < 0.0 {
		var region_max_x float = GetContentRegionMaxAbs().x
		w = ImMax(1.0, region_max_x-window.DC.CursorPos.x+w)
	}
	w = IM_FLOOR(w)
	return w
}

// FIXME: All the Contents Region function are messy or misleading. WE WILL AIM TO OBSOLETE ALL OF THEM WITH A NEW "WORK RECT" API. Thanks for your patience!

// Content region
// - Retrieve available space from a given point. GetContentRegionAvail() is frequently useful.
// - Those functions are bound to be redesigned (they are confusing, incomplete and the Min/Max return values are in local window coordinates which increases confusion)

// == GetContentRegionMax() - GetCursorPos()
func GetContentRegionAvail() ImVec2 {
	var window = GImGui.CurrentWindow
	return GetContentRegionMaxAbs().Sub(window.DC.CursorPos)
}

// current content boundaries (typically window boundaries including scrolling, or current column boundaries), in windows coordinates
// FIXME: This is in window space (not screen space!).
func GetContentRegionMax() ImVec2 {
	var g = GImGui
	var window = g.CurrentWindow
	var mx = window.ContentRegionRect.Max.Sub(window.Pos)
	if window.DC.CurrentColumns != nil || g.CurrentTable != nil {
		mx.x = window.WorkRect.Max.x - window.Pos.x
	}
	return mx
}

// [Internal] Absolute coordinate. Saner. This is not exposed until we finishing refactoring work rect features.
func GetContentRegionMaxAbs() ImVec2 {
	var g = GImGui
	var window = g.CurrentWindow
	var mx = window.ContentRegionRect.Max
	if window.DC.CurrentColumns != nil || g.CurrentTable != nil {
		mx.x = window.WorkRect.Max.x
	}
	return mx
}

// content boundaries min for the full window (roughly (0,0)-Scroll), in window coordinates
func GetWindowContentRegionMin() ImVec2 {
	var window = GImGui.CurrentWindow
	return window.ContentRegionRect.Min.Sub(window.Pos)
}

// content boundaries max for the full window (roughly (0,0)+Size-Scroll) where Size can be override with SetNextWindowContentSize(), in window coordinates
func GetWindowContentRegionMax() ImVec2 {
	var window = GImGui.CurrentWindow
	return window.ContentRegionRect.Max.Sub(window.Pos)
}
