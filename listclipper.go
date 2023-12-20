package imgui

// Helper: Manually clip large list of items.
// If you are submitting lots of evenly spaced items and you have a random access to the list, you can perform coarse
// clipping based on visibility to save yourself from processing those items at all.
// The clipper calculates the range of visible items and advance the cursor to compensate for the non-visible items we have skipped.
// (Dear ImGui already clip items based on their bounds but it needs to measure text size to do so, whereas manual coarse clipping before submission makes this cost and your own data fetching/submission cost almost null)
// Usage:
//
//	clipper ImGuiListClipper
//
// clipper.Begin(1000) //         // We have 1000 elements, evenly spaced.
//
//	while (clipper.Step())
//	    for (int i clipper.DisplayStart = i clipper.DisplayEnd < i++)
//	        ImGui::Text("line number i) %d",
//
// Generally what happens is:
// - Clipper lets you process the first element (DisplayStart  DisplayEnd = 1) regardless of it being visible or not.
// - User code submit one element.
// - Clipper can measure the height of the first element
// - Clipper calculate the actual range of elements to display based on the current clipping rectangle, position the cursor before the first visible element.
// - User code submit visible elements.
type ImGuiListClipper struct {
	DisplayStart int
	DisplayEnd   int

	// [Internal]
	ItemsCount  int
	StepNo      int
	ItemsFrozen int
	ItemsHeight float
	StartPosY   float
}

func NewImGuiListClipper() ImGuiListClipper {
	return ImGuiListClipper{
		ItemsCount: -1,
	}
}

// Use case A: Begin() called from constructor with items_height<0, then called again from Step() in StepNo 1
// Use case B: Begin() called from constructor with items_height>0
// FIXME-LEGACY: Ideally we should remove the Begin/End functions but they are part of the legacy API we still support. This is why some of the code in Step() calling Begin() and reassign some fields, spaghetti style.
// items_count: Use INT_MAX if you don't know how many items you have (in which case the cursor won't be advanced in the final step)
// items_height: Use -1.0f to be calculated automatically on first step. Otherwise pass in the distance between your items, typically GetTextLineHeightWithSpacing() or GetFrameHeightWithSpacing().
func (this ImGuiListClipper) Begin(items_count int, items_height float /*= -1.0f*/) {
	var g = GImGui
	var window = g.CurrentWindow

	if table := g.CurrentTable; table != nil {
		if table.IsInsideRow {
			TableEndRow(table)
		}
	}

	this.StartPosY = window.DC.CursorPos.y
	this.ItemsHeight = items_height
	this.ItemsCount = items_count
	this.ItemsFrozen = 0
	this.StepNo = 0
	this.DisplayStart = -1
	this.DisplayEnd = 0
}

// Automatically called on the last call of Step() that returns false.
func (this ImGuiListClipper) End() {
	if this.ItemsCount < 0 { // Already ended
		return
	}

	// In theory here we should assert that ImGui::GetCursorPosY() == StartPosY + DisplayEnd * ItemsHeight, but it feels saner to just seek at the end and not assert/crash the user.
	if this.ItemsCount < INT_MAX && this.DisplayStart >= 0 {
		SetCursorPosYAndSetupForPrevLine(this.StartPosY+float(this.ItemsCount-this.ItemsFrozen)*this.ItemsHeight, this.ItemsHeight)
	}
	this.ItemsCount = -1
	this.StepNo = 3
}

func (this ImGuiListClipper) Step() bool {
	var g = GImGui
	var window = g.CurrentWindow

	var table = g.CurrentTable
	if table != nil && table.IsInsideRow {
		TableEndRow(table)
	}

	// No items
	if this.ItemsCount == 0 || GetSkipItemForListClipping() {
		this.End()
		return false
	}

	// Step 0: Let you process the first element (regardless of it being visible or not, so we can measure the element height)
	if this.StepNo == 0 {
		// While we are in frozen row state, keep displaying items one by one, unclipped
		// FIXME: Could be stored as a table-agnostic state.
		if table != nil && !table.IsUnfrozenRows {
			this.DisplayStart = this.ItemsFrozen
			this.DisplayEnd = this.ItemsFrozen + 1
			this.ItemsFrozen++
			return true
		}

		this.StartPosY = window.DC.CursorPos.y
		if this.ItemsHeight <= 0.0 {
			// Submit the first item so we can measure its height (generally it is 0..1)
			this.DisplayStart = this.ItemsFrozen
			this.DisplayEnd = this.ItemsFrozen + 1
			this.StepNo = 1
			return true
		}

		// Already has item height (given by user in Begin): skip to calculating step
		this.DisplayStart = this.DisplayEnd
		this.StepNo = 2
	}

	// Step 1: the clipper infer height from first element
	if this.StepNo == 1 {
		IM_ASSERT(this.ItemsHeight <= 0.0)
		if table != nil {
			var pos_y1 = table.RowPosY1 // Using this instead of StartPosY to handle clipper straddling the frozen row
			var pos_y2 = table.RowPosY2 // Using this instead of CursorPos.y to take account of tallest cell.
			this.ItemsHeight = pos_y2 - pos_y1
			window.DC.CursorPos.y = pos_y2
		} else {
			this.ItemsHeight = window.DC.CursorPos.y - this.StartPosY
		}
		IM_ASSERT_USER_ERROR(this.ItemsHeight > 0.0, "Unable to calculate item height! First item hasn't moved the cursor vertically!")
		this.StepNo = 2
	}

	// Reached end of list
	if this.DisplayEnd >= this.ItemsCount {
		this.End()
		return false
	}

	// Step 2: calculate the actual range of elements to display, and position the cursor before the first element
	if this.StepNo == 2 {
		IM_ASSERT(this.ItemsHeight > 0.0)

		var already_submitted = this.DisplayEnd
		CalcListClipping(this.ItemsCount-already_submitted, this.ItemsHeight, &this.DisplayStart, &this.DisplayEnd)
		this.DisplayStart += already_submitted
		this.DisplayEnd += already_submitted

		// Seek cursor
		if this.DisplayStart > already_submitted {
			SetCursorPosYAndSetupForPrevLine(this.StartPosY+float(this.DisplayStart-this.ItemsFrozen)*this.ItemsHeight, this.ItemsHeight)
		}

		this.StepNo = 3
		return true
	}

	// Step 3: the clipper validate that we have reached the expected Y position (corresponding to element DisplayEnd),
	// Advance the cursor to the end of the list and then returns 'false' to end the loop.
	if this.StepNo == 3 {
		// Seek cursor
		if this.ItemsCount < INT_MAX {
			SetCursorPosYAndSetupForPrevLine(this.StartPosY+float(this.ItemsCount-this.ItemsFrozen)*this.ItemsHeight, this.ItemsHeight) // advance cursor
		}
		this.ItemsCount = -1
		return false
	}

	IM_ASSERT(false)
	return false
}

// FIXME-TABLE: This prevents us from using ImGuiListClipper _inside_ a table cell.
// The problem we have is that without a Begin/End scheme for rows using the clipper is ambiguous.
func GetSkipItemForListClipping() bool {
	var g = GImGui
	if g.CurrentTable != nil {
		return g.CurrentTable.HostSkipItems
	}
	return g.CurrentWindow.SkipItems
}

// Helper to calculate coarse clipping of large list of evenly sized items.
// NB: Prefer using the ImGuiListClipper higher-level helper if you can! Read comments and instructions there on how those use this sort of pattern.
// NB: 'items_count' is only used to clamp the result, if you don't know your count you can use INT_MAX
func CalcListClipping(items_count int, items_height float, out_items_display_start *int, out_items_display_end *int) {
	var g = GImGui
	var window = g.CurrentWindow
	if g.LogEnabled {
		// If logging is active, do not perform any clipping
		*out_items_display_start = 0
		*out_items_display_end = items_count
		return
	}
	if GetSkipItemForListClipping() {
		*out_items_display_start = 0
		*out_items_display_end = 0
		return
	}

	// We create the union of the ClipRect and the scoring rect which at worst should be 1 page away from ClipRect
	var unclipped_rect = window.ClipRect
	if g.NavMoveScoringItems {
		unclipped_rect.AddRect(g.NavScoringRect)
	}
	if g.NavJustMovedToId != 0 && window.NavLastIds[0] == g.NavJustMovedToId {
		// Could store and use NavJustMovedToRectRe
		unclipped_rect.AddRect(ImRect{window.Pos.Add(window.NavRectRel[0].Min), window.Pos.Add(window.NavRectRel[0].Max)})
	}

	var pos = window.DC.CursorPos
	var start = (int)((unclipped_rect.Min.y - pos.y) / items_height)
	var end = (int)((unclipped_rect.Max.y - pos.y) / items_height)

	// When performing a navigation request, ensure we have one item extra in the direction we are moving to
	if g.NavMoveScoringItems && g.NavMoveClipDir == ImGuiDir_Up {
		start--
	}
	if g.NavMoveScoringItems && g.NavMoveClipDir == ImGuiDir_Down {
		end++
	}

	start = ImClampInt(start, 0, items_count)
	end = ImClampInt(end+1, start, items_count)
	*out_items_display_start = start
	*out_items_display_end = end
}

func SetCursorPosYAndSetupForPrevLine(pos_y, line_height float) {
	// Set cursor position and a few other things so that SetScrollHereY() and Columns() can work when seeking cursor.
	// FIXME: It is problematic that we have to do that here, because custom/equivalent end-user code would stumble on the same issue.
	// The clipper should probably have a 4th step to display the last item in a regular manner.
	var g = GImGui
	var window = g.CurrentWindow
	var off_y = pos_y - window.DC.CursorPos.y
	window.DC.CursorPos.y = pos_y
	window.DC.CursorMaxPos.y = ImMax(window.DC.CursorMaxPos.y, pos_y)
	window.DC.CursorPosPrevLine.y = window.DC.CursorPos.y - line_height // Setting those fields so that SetScrollHereY() can properly function after the end of our clipper usage.
	window.DC.PrevLineSize.y = (line_height - g.Style.ItemSpacing.y)    // If we end up needing more accurate data (to e.g. use SameLine) we may as well make the clipper have a fourth step to let user process and display the last item in their list.
	if columns := window.DC.CurrentColumns; columns != nil {
		columns.LineMinY = window.DC.CursorPos.y // Setting this so that cell Y position are set properly
	}
	if table := g.CurrentTable; table != nil {
		if table.IsInsideRow {
			TableEndRow(table)
		}
		table.RowPosY2 = window.DC.CursorPos.y
		var row_increase = (int)((off_y / line_height) + 0.5)
		//table.CurrentRow += row_increase; // Can't do without fixing TableEndRow()
		table.RowBgColorCounter += row_increase
	}
}
