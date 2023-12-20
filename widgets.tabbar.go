package imgui

import (
	"sort"

	"github.com/Splizard/imgui/golang"
)

type ImGuiTabItem struct {
	ID                ImGuiID
	Flags             ImGuiTabItemFlags
	LastFrameVisible  int
	LastFrameSelected int   // This allows us to infer an ordered list of the last activated tabs with little maintenance
	Offset            float // Position relative to beginning of tab
	Width             float // Width currently displayed
	ContentWidth      float // Width of label, stored during BeginTabItem() call
	NameOffset        int   // When Window==NULL, offset to name within parent ImGuiTabBar::TabsNames
	BeginOrder        int   // BeginTabItem() order, used to re-order tabs after toggling ImGuiTabBarFlags_Reorderable
	IndexDuringLayout int   // Index only used during TabBarLayout()
	WantClose         bool  // Marked as closed by SetTabItemClosed()
}

func NewImGuiTabItem() ImGuiTabItem {
	return ImGuiTabItem{
		LastFrameVisible:  -1,
		LastFrameSelected: -1,
		NameOffset:        -1,
		BeginOrder:        -1,
		IndexDuringLayout: -1,
	}
}

type ImGuiTabBar struct {
	Tabs                            []ImGuiTabItem
	Flags                           ImGuiTabBarFlags
	ID                              ImGuiID  // Zero for tab-bars used by docking
	SelectedTabId                   ImGuiID  // Selected tab/window
	NextSelectedTabId               ImGuiID  // Next selected tab/window. Will also trigger a scrolling animation
	VisibleTabId                    ImGuiID  // Can occasionally be != SelectedTabId (e.g. when previewing contents for CTRL+TAB preview)
	CurrFrameVisible                int      //
	PrevFrameVisible                int      //
	BarRect                         ImRect   //
	CurrTabsContentsHeight          float    //
	PrevTabsContentsHeight          float    // Record the height of contents submitted below the tab bar
	WidthAllTabs                    float    // Actual width of all tabs (locked during layout)
	WidthAllTabsIdeal               float    // Ideal width if all tabs were visible and not clipped
	ScrollingAnim                   float    //
	ScrollingTarget                 float    //
	ScrollingTargetDistToVisibility float    //
	ScrollingSpeed                  float    //
	ScrollingRectMinX               float    //
	ScrollingRectMaxX               float    //
	ReorderRequestTabId             ImGuiID  //
	ReorderRequestOffset            int      //
	BeginCount                      int      //
	WantLayout                      bool     //
	VisibleTabWasSubmitted          bool     // Set to true when a new tab item or button has been added to the tab bar during last frame
	TabsAddedNew                    bool     // Set to true when a new tab item or button has been added to the tab bar during last frame
	TabsActiveCount                 int      // Number of tabs submitted this frame.
	LastTabItemIdx                  int      // Index of last BeginTabItem() tab for use by EndTabItem()
	ItemSpacingY                    float    //
	FramePadding                    ImVec2   // style.FramePadding locked at the time of BeginTabBar()
	BackupCursorPos                 ImVec2   //
	TabsNames                       []string // For non-docking tab bar we re-append names in a contiguous buffer.
}

func NewImGuiTabBar() ImGuiTabBar {
	var tabbar ImGuiTabBar
	tabbar.CurrFrameVisible = -1
	tabbar.PrevFrameVisible = -1
	tabbar.LastTabItemIdx = -1
	return tabbar
}

func (this ImGuiTabBar) GetTabOrder(tab *ImGuiTabItem) int {
	for i := range this.Tabs {
		if tab == &this.Tabs[i] {
			return int(i)
		}
	}
	return -1
}

func (this ImGuiTabBar) GetTabName(tab *ImGuiTabItem) string {
	IM_ASSERT(tab.NameOffset != -1 && tab.NameOffset < int(len(this.TabsNames)))
	return string(this.TabsNames[tab.NameOffset])
}

type ImGuiTabBarSection struct {
	TabCount int   // Number of tabs in this section.
	Width    float // Sum of width of tabs in this section (after shrinking down)
	Spacing  float // Horizontal spacing at the end of the section.
}

func TabItemGetSectionIdx(tab *ImGuiTabItem) int {
	if tab.Flags&ImGuiTabItemFlags_Leading != 0 {
		return 0
	} else {
		if (tab.Flags & ImGuiTabItemFlags_Trailing) != 0 {
			return 2
		}
	}
	return 1
}

func GetTabBarFromTabBarRef(ref ImGuiPtrOrIndex) *ImGuiTabBar {
	var g = GImGui
	if ref.Ptr != nil {
		return ref.Ptr.(*ImGuiTabBar)
	} else {
		return g.TabBars[uint(ref.Index)]
	}
}

func GetTabBarRefFromTabBar(tab_bar *ImGuiTabBar) ImGuiPtrOrIndex {
	var g = GImGui
	for idx, tab_bar_ref := range g.TabBars {
		if tab_bar == tab_bar_ref {
			return ImGuiPtrOrIndex{Index: int(idx)}
		}
	}
	return ImGuiPtrOrIndex{Ptr: tab_bar}
}

// Tab Bars, Tabs

// create and append into a TabBar
func BeginTabBar(str_id string, flags ImGuiTabBarFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}

	var id = window.GetIDs(str_id)
	var tab_bar, ok = g.TabBars[id]
	if !ok {
		tab_bar = new(ImGuiTabBar)
		tab_bar.ID = id
		g.TabBars[id] = tab_bar
	}
	var tab_bar_bb = ImRect{ImVec2{window.DC.CursorPos.x, window.DC.CursorPos.y}, ImVec2{window.WorkRect.Max.x, window.DC.CursorPos.y + g.FontSize + g.Style.FramePadding.y*2}}
	tab_bar.ID = id
	return BeginTabBarEx(tab_bar, &tab_bar_bb, flags|ImGuiTabBarFlags_IsFocused)
}

// only call EndTabBar() if BeginTabBar() returns true!
func EndTabBar() {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return
	}

	var tab_bar = g.CurrentTabBar
	if tab_bar == nil {
		IM_ASSERT_USER_ERROR(tab_bar != nil, "Mismatched BeginTabBar()/EndTabBar()!")
		return
	}

	// Fallback in case no TabItem have been submitted
	if tab_bar.WantLayout {
		TabBarLayout(tab_bar)
	}

	// Restore the last visible height if no tab is visible, this reduce vertical flicker/movement when a tabs gets removed without calling SetTabItemClosed().
	var tab_bar_appearing = (tab_bar.PrevFrameVisible+1 < g.FrameCount)
	if tab_bar.VisibleTabWasSubmitted || tab_bar.VisibleTabId == 0 || tab_bar_appearing {
		tab_bar.CurrTabsContentsHeight = ImMax(window.DC.CursorPos.y-tab_bar.BarRect.Max.y, tab_bar.CurrTabsContentsHeight)
		window.DC.CursorPos.y = tab_bar.BarRect.Max.y + tab_bar.CurrTabsContentsHeight
	} else {
		window.DC.CursorPos.y = tab_bar.BarRect.Max.y + tab_bar.PrevTabsContentsHeight
	}
	if tab_bar.BeginCount > 1 {
		window.DC.CursorPos = tab_bar.BackupCursorPos
	}

	if (tab_bar.Flags & ImGuiTabBarFlags_DockNode) == 0 {
		PopID()
	}

	g.CurrentTabBarStack = g.CurrentTabBarStack[:len(g.CurrentTabBarStack)-1]
	g.CurrentTabBar = nil
	if len(g.CurrentTabBarStack) > 0 {
		g.CurrentTabBar = GetTabBarFromTabBarRef(g.CurrentTabBarStack[len(g.CurrentTabBarStack)-1])
	}
}

// create a Tab. Returns true if the Tab is selected.
func BeginTabItem(label string, p_open *bool, flags ImGuiTabItemFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}

	var tab_bar = g.CurrentTabBar
	if tab_bar == nil {
		IM_ASSERT_USER_ERROR(tab_bar != nil, "Needs to be called between BeginTabBar() and EndTabBar()!")
		return false
	}
	IM_ASSERT((flags & ImGuiTabItemFlags_Button) == 0) // BeginTabItem() Can't be used with button flags, use TabItemButton() instead!

	var ret = TabItemEx(tab_bar, label, p_open, flags)
	if ret && (flags&ImGuiTabItemFlags_NoPushId) == 0 {
		var tab = &tab_bar.Tabs[tab_bar.LastTabItemIdx]
		PushOverrideID(tab.ID) // We already hashed 'label' so push into the ID stack directly instead of doing another hash through PushID(label)
	}
	return ret
}

// only call EndTabItem() if BeginTabItem() returns true!
func EndTabItem() {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return
	}

	var tab_bar = g.CurrentTabBar
	if tab_bar == nil {
		IM_ASSERT_USER_ERROR(tab_bar != nil, "Needs to be called between BeginTabBar() and EndTabBar()!")
		return
	}
	IM_ASSERT(tab_bar.LastTabItemIdx >= 0)
	var tab = &tab_bar.Tabs[tab_bar.LastTabItemIdx]
	if (tab.Flags & ImGuiTabItemFlags_NoPushId) == 0 {
		PopID()
	}
}

// create a Tab behaving like a button. return true when clicked. cannot be selected in the tab bar.
func TabItemButton(label string, flags ImGuiTabItemFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}

	var tab_bar = g.CurrentTabBar
	if tab_bar == nil {
		IM_ASSERT_USER_ERROR(tab_bar != nil, "Needs to be called between BeginTabBar() and EndTabBar()!")
		return false
	}
	return TabItemEx(tab_bar, label, nil, flags|ImGuiTabItemFlags_Button|ImGuiTabItemFlags_NoReorder)
}

// notify TabBar or Docking system of a closed tab/window ahead (useful to reduce visual flicker on reorderable tab bars).
// For tab-bar: call after BeginTabBar() and before Tab submissions. Otherwise call with a window name.
// [Public] This is call is 100% optional but it allows to remove some one-frame glitches when a tab has been unexpectedly removed.
// To use it to need to call the function SetTabItemClosed() between BeginTabBar() and EndTabBar().
// Tabs closed by the close button will automatically be flagged to avoid this issue.
func SetTabItemClosed(tab_or_docked_window_label string) {
	var g = GImGui
	var is_within_manual_tab_bar = g.CurrentTabBar != nil && (g.CurrentTabBar.Flags&ImGuiTabBarFlags_DockNode) == 0
	if is_within_manual_tab_bar {
		var tab_bar = g.CurrentTabBar
		var tab_id = TabBarCalcTabID(tab_bar, tab_or_docked_window_label)
		if tab := TabBarFindTabByID(tab_bar, tab_id); tab != nil {
			tab.WantClose = true // Will be processed by next call to TabBarLayout()
		}
	}
}

// Tab Bars
func BeginTabBarEx(tab_bar *ImGuiTabBar, tab_bar_bb *ImRect, flags ImGuiTabBarFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}

	if (flags & ImGuiTabBarFlags_DockNode) == 0 {
		PushOverrideID(tab_bar.ID)
	}

	// Add to stack
	g.CurrentTabBarStack = append(g.CurrentTabBarStack, GetTabBarRefFromTabBar(tab_bar))
	g.CurrentTabBar = tab_bar

	// Append with multiple BeginTabBar()/EndTabBar() pairs.
	tab_bar.BackupCursorPos = window.DC.CursorPos
	if tab_bar.CurrFrameVisible == g.FrameCount {
		window.DC.CursorPos = ImVec2{tab_bar.BarRect.Min.x, tab_bar.BarRect.Max.y + tab_bar.ItemSpacingY}
		tab_bar.BeginCount++
		return true
	}

	// Ensure correct ordering when toggling ImGuiTabBarFlags_Reorderable flag, or when a new tab was added while being not reorderable
	if (flags&ImGuiTabBarFlags_Reorderable) != (tab_bar.Flags&ImGuiTabBarFlags_Reorderable) || (tab_bar.TabsAddedNew && flags&ImGuiTabBarFlags_Reorderable == 0) {
		if len(tab_bar.Tabs) > 1 {
			sort.Slice(tab_bar.Tabs, func(i, j golang.Int) bool {
				a := tab_bar.Tabs[i]
				b := tab_bar.Tabs[j]
				return a.BeginOrder < b.BeginOrder
			})
		}
	}
	tab_bar.TabsAddedNew = false

	// Flags
	if (flags & ImGuiTabBarFlags_FittingPolicyMask_) == 0 {
		flags |= ImGuiTabBarFlags_FittingPolicyDefault_
	}

	tab_bar.Flags = flags
	tab_bar.BarRect = *tab_bar_bb
	tab_bar.WantLayout = true // Layout will be done on the first call to ItemTab()
	tab_bar.PrevFrameVisible = tab_bar.CurrFrameVisible
	tab_bar.CurrFrameVisible = g.FrameCount
	tab_bar.PrevTabsContentsHeight = tab_bar.CurrTabsContentsHeight
	tab_bar.CurrTabsContentsHeight = 0.0
	tab_bar.ItemSpacingY = g.Style.ItemSpacing.y
	tab_bar.FramePadding = g.Style.FramePadding
	tab_bar.TabsActiveCount = 0
	tab_bar.BeginCount = 1

	// Set cursor pos in a way which only be used in the off-chance the user erroneously submits item before BeginTabItem(): items will overlap
	window.DC.CursorPos = ImVec2{tab_bar.BarRect.Min.x, tab_bar.BarRect.Max.y + tab_bar.ItemSpacingY}

	// Draw separator
	c := ImGuiCol_TabUnfocusedActive
	if (flags & ImGuiTabBarFlags_IsFocused) != 0 {
		c = ImGuiCol_TabActive
	}

	var col = GetColorU32FromID(c, 1)
	var y = tab_bar.BarRect.Max.y - 1.0
	{
		var separator_min_x = tab_bar.BarRect.Min.x - IM_FLOOR(window.WindowPadding.x*0.5)
		var separator_max_x = tab_bar.BarRect.Max.x + IM_FLOOR(window.WindowPadding.x*0.5)
		window.DrawList.AddLine(&ImVec2{separator_min_x, y}, &ImVec2{separator_max_x, y}, col, 1.0)
	}
	return true
}

func TabBarFindTabByID(tab_bar *ImGuiTabBar, tab_id ImGuiID) *ImGuiTabItem {
	if tab_id != 0 {
		for n := range tab_bar.Tabs {
			if tab_bar.Tabs[n].ID == tab_id {
				return &tab_bar.Tabs[n]
			}
		}
	}
	return nil
}

// The *TabId fields be already set by the docking system _before_ the actual TabItem was created, so we clear them regardless.
func TabBarRemoveTab(tab_bar *ImGuiTabBar, tab_id ImGuiID) {
	if tab := TabBarFindTabByID(tab_bar, tab_id); tab != nil {
		//tab_bar.Tabs.erase(tab)
		for i, t := range tab_bar.Tabs {
			if t.ID == tab_id {
				tab_bar.Tabs = append(tab_bar.Tabs[:i], tab_bar.Tabs[i+1:]...)
				break
			}
		}
	}
	if tab_bar.VisibleTabId == tab_id {
		tab_bar.VisibleTabId = 0
	}
	if tab_bar.SelectedTabId == tab_id {
		tab_bar.SelectedTabId = 0
	}
	if tab_bar.NextSelectedTabId == tab_id {
		tab_bar.NextSelectedTabId = 0
	}
}

// Called on manual closure attempt
func TabBarCloseTab(tab_bar *ImGuiTabBar, tab *ImGuiTabItem) {
	IM_ASSERT((tab.Flags & ImGuiTabItemFlags_Button) == 0)
	if (tab.Flags & ImGuiTabItemFlags_UnsavedDocument) == 0 {
		// This will remove a frame of lag for selecting another tab on closure.
		// However we don't run it in the case where the 'Unsaved' flag is set, so user gets a chance to fully undo the closure
		tab.WantClose = true
		if tab_bar.VisibleTabId == tab.ID {
			tab.LastFrameVisible = -1
			tab_bar.SelectedTabId = 0
			tab_bar.NextSelectedTabId = 0
		}
	} else {
		// Actually select before expecting closure attempt (on an UnsavedDocument tab user is expect to e.g. show a popup)
		if tab_bar.VisibleTabId != tab.ID {
			tab_bar.NextSelectedTabId = tab.ID
		}
	}
}

func TabBarQueueReorder(tab_bar *ImGuiTabBar, tab *ImGuiTabItem, offset int) {
	IM_ASSERT(offset != 0)
	IM_ASSERT(tab_bar.ReorderRequestTabId == 0)
	tab_bar.ReorderRequestTabId = tab.ID
	tab_bar.ReorderRequestOffset = (int)(offset)
}

func TabBarQueueReorderFromMousePos(tab_bar *ImGuiTabBar, src_tab *ImGuiTabItem, mouse_pos ImVec2) {
	var g = GImGui
	IM_ASSERT(tab_bar.ReorderRequestTabId == 0)
	if (tab_bar.Flags & ImGuiTabBarFlags_Reorderable) == 0 {
		return
	}

	var is_central_section = (src_tab.Flags & ImGuiTabItemFlags_SectionMask_) == 0
	var bar_offset = tab_bar.BarRect.Min.x
	if is_central_section {
		bar_offset -= tab_bar.ScrollingTarget
	}

	// Count number of contiguous tabs we are crossing over
	var dir int = 1
	if (bar_offset + src_tab.Offset) > mouse_pos.x {
		dir = -1
	}
	var src_idx int
	for n := range tab_bar.Tabs {
		if src_tab == &tab_bar.Tabs[n] {
			src_idx = int(n)
			break
		}
	}
	var dst_idx = src_idx
	for i := src_idx; i >= 0 && i < int(len(tab_bar.Tabs)); i += dir {
		// Reordered tabs must share the same section
		var dst_tab = &tab_bar.Tabs[i]
		if dst_tab.Flags&ImGuiTabItemFlags_NoReorder != 0 {
			break
		}
		if (dst_tab.Flags & ImGuiTabItemFlags_SectionMask_) != (src_tab.Flags & ImGuiTabItemFlags_SectionMask_) {
			break
		}
		dst_idx = i

		// Include spacing after tab, so when mouse cursor is between tabs we would not continue checking further tabs that are not hovered.
		var x1 = bar_offset + dst_tab.Offset - g.Style.ItemInnerSpacing.x
		var x2 = bar_offset + dst_tab.Offset + dst_tab.Width + g.Style.ItemInnerSpacing.x
		//GetForegroundDrawList().AddRect(ImVec2(x1, tab_bar.BarRect.Min.y), ImVec2(x2, tab_bar.BarRect.Max.y), IM_COL32(255, 0, 0, 255));
		if (dir < 0 && mouse_pos.x > x1) || (dir > 0 && mouse_pos.x < x2) {
			break
		}
	}

	if dst_idx != src_idx {
		TabBarQueueReorder(tab_bar, src_tab, dst_idx-src_idx)
	}
}

func TabBarProcessReorder(tab_bar *ImGuiTabBar) bool {
	var tab1 = TabBarFindTabByID(tab_bar, tab_bar.ReorderRequestTabId)
	if tab1 == nil || (tab1.Flags&ImGuiTabItemFlags_NoReorder) != 0 {
		return false
	}
	var tab1_idx int
	for n := range tab_bar.Tabs {
		if &tab_bar.Tabs[n] == tab1 {
			tab1_idx = int(n)
			break
		}
	}
	var tab1_slice = tab_bar.Tabs[tab1_idx:]

	//IM_ASSERT(tab_bar.Flags & ImGuiTabBarFlags_Reorderable); // <- this may happen when using debug tools
	var tab2_order = tab_bar.GetTabOrder(tab1) + tab_bar.ReorderRequestOffset
	if tab2_order < 0 || tab2_order >= int(len(tab_bar.Tabs)) {
		return false
	}

	// Reordered tabs must share the same section
	// (Note: TabBarQueueReorderFromMousePos() also has a similar test but since we allow direct calls to TabBarQueueReorder() we do it here too)
	var tab2 = tab_bar.Tabs[tab2_order:]
	if (tab2[0].Flags & ImGuiTabItemFlags_NoReorder) != 0 {
		return false
	}
	if (tab1.Flags & ImGuiTabItemFlags_SectionMask_) != (tab2[0].Flags & ImGuiTabItemFlags_SectionMask_) {
		return false
	}

	var item_tmp = *tab1
	var src_tab = tab2
	var dst_tab = tab2[1:]
	if tab_bar.ReorderRequestOffset > 0 {
		src_tab, dst_tab = tab1_slice[1:], tab1_slice
	}
	var move_count = -tab_bar.ReorderRequestOffset
	if tab_bar.ReorderRequestOffset > 0 {
		move_count = tab_bar.ReorderRequestOffset
	}
	copy(dst_tab, src_tab[:move_count])
	tab2[0] = item_tmp

	if tab_bar.Flags&ImGuiTabBarFlags_SaveSettings != 0 {
		MarkIniSettingsDirty()
	}
	return true
}

func TabItemEx(tab_bar *ImGuiTabBar, label string, p_open *bool, flags ImGuiTabItemFlags) bool {
	// Layout whole tab bar if not already done
	if tab_bar.WantLayout {
		TabBarLayout(tab_bar)
	}

	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}

	var style = g.Style
	var id = TabBarCalcTabID(tab_bar, label)

	// If the user called us with *p_open == false, we early out and don't render.
	// We make a call to ItemAdd() so that attempts to use a contextual popup menu with an implicit ID won't use an older ID.
	if p_open != nil && !*p_open {
		PushItemFlag(ImGuiItemFlags_NoNav|ImGuiItemFlags_NoNavDefaultFocus, true)
		ItemAdd(&ImRect{}, id, nil, 0)
		PopItemFlag()
		return false
	}

	IM_ASSERT(p_open == nil || (flags&ImGuiTabItemFlags_Button) == 0)
	IM_ASSERT((flags & (ImGuiTabItemFlags_Leading | ImGuiTabItemFlags_Trailing)) != (ImGuiTabItemFlags_Leading | ImGuiTabItemFlags_Trailing)) // Can't use both Leading and Trailing

	// Store into ImGuiTabItemFlags_NoCloseButton, also honor ImGuiTabItemFlags_NoCloseButton passed by user (although not documented)
	if flags&ImGuiTabItemFlags_NoCloseButton != 0 {
		p_open = nil
	} else if p_open == nil {
		flags |= ImGuiTabItemFlags_NoCloseButton
	}

	// Calculate tab contents size
	var size = TabItemCalcSize(label, p_open != nil)

	// Acquire tab data
	var tab = TabBarFindTabByID(tab_bar, id)
	var tab_is_new = false
	if tab == nil {
		tab_bar.Tabs = append(tab_bar.Tabs, NewImGuiTabItem())
		tab = &tab_bar.Tabs[len(tab_bar.Tabs)-1]
		tab.ID = id
		tab.Width = size.x
		tab_bar.TabsAddedNew = true
		tab_is_new = true
	}
	for n := range tab_bar.Tabs {
		if &tab_bar.Tabs[n] == tab {
			tab_bar.LastTabItemIdx = int(n)
			break
		}
	}
	tab.ContentWidth = size.x
	tab.BeginOrder = tab_bar.TabsActiveCount
	tab_bar.TabsActiveCount++

	var tab_bar_appearing = (tab_bar.PrevFrameVisible+1 < g.FrameCount)
	var tab_bar_focused = (tab_bar.Flags & ImGuiTabBarFlags_IsFocused) != 0
	var tab_appearing = (tab.LastFrameVisible+1 < g.FrameCount)
	var is_tab_button = (flags & ImGuiTabItemFlags_Button) != 0
	tab.LastFrameVisible = g.FrameCount
	tab.Flags = flags

	// Append name with zero-terminator
	tab.NameOffset = (ImS32)(len(tab_bar.TabsNames))
	tab_bar.TabsNames = append(tab_bar.TabsNames, label)

	// Update selected tab
	if tab_appearing && (tab_bar.Flags&ImGuiTabBarFlags_AutoSelectNewTabs) != 0 && tab_bar.NextSelectedTabId == 0 {
		if !tab_bar_appearing || tab_bar.SelectedTabId == 0 {
			if !is_tab_button {
				tab_bar.NextSelectedTabId = id // New tabs gets activated
			}
		}
	}
	if (flags&ImGuiTabItemFlags_SetSelected) != 0 && (tab_bar.SelectedTabId != id) { // SetSelected can only be passed on explicit tab bar
		if !is_tab_button {
			tab_bar.NextSelectedTabId = id
		}
	}

	// Lock visibility
	// (Note: tab_contents_visible != tab_selected... because CTRL+TAB operations may preview some tabs without selecting them!)
	var tab_contents_visible = (tab_bar.VisibleTabId == id)
	if tab_contents_visible {
		tab_bar.VisibleTabWasSubmitted = true
	}

	// On the very first frame of a tab bar we let first tab contents be visible to minimize appearing glitches
	if !tab_contents_visible && tab_bar.SelectedTabId == 0 && tab_bar_appearing {
		if len(tab_bar.Tabs) == 1 && (tab_bar.Flags&ImGuiTabBarFlags_AutoSelectNewTabs) == 0 {
			tab_contents_visible = true
		}
	}

	// Note that tab_is_new is not necessarily the same as tab_appearing! When a tab bar stops being submitted
	// and then gets submitted again, the tabs will have 'tab_appearing=true' but 'tab_is_new=false'.
	if tab_appearing && (!tab_bar_appearing || tab_is_new) {
		PushItemFlag(ImGuiItemFlags_NoNav|ImGuiItemFlags_NoNavDefaultFocus, true)
		ItemAdd(&ImRect{}, id, nil, 0)
		PopItemFlag()
		if is_tab_button {
			return false
		}
		return tab_contents_visible
	}

	if tab_bar.SelectedTabId == id {
		tab.LastFrameSelected = g.FrameCount
	}

	// Backup current layout position
	var backup_main_cursor_pos = window.DC.CursorPos

	// Layout
	var is_central_section = (tab.Flags & ImGuiTabItemFlags_SectionMask_) == 0
	size.x = tab.Width
	if is_central_section {
		window.DC.CursorPos = tab_bar.BarRect.Min.Add(ImVec2{IM_FLOOR(tab.Offset - tab_bar.ScrollingAnim), 0.0})
	} else {
		window.DC.CursorPos = tab_bar.BarRect.Min.Add(ImVec2{tab.Offset, 0.0})
	}
	var pos = window.DC.CursorPos
	var bb = ImRect{pos, pos.Add(size)}

	// We don't have CPU clipping primitives to clip the CloseButton (until it becomes a texture), so need to add an extra draw call (temporary in the case of vertical animation)
	var want_clip_rect = is_central_section && (bb.Min.x < tab_bar.ScrollingRectMinX || bb.Max.x > tab_bar.ScrollingRectMaxX)
	if want_clip_rect {
		PushClipRect(ImVec2{ImMax(bb.Min.x, tab_bar.ScrollingRectMinX), bb.Min.y - 1}, ImVec2{tab_bar.ScrollingRectMaxX, bb.Max.y}, true)
	}

	var backup_cursor_max_pos = window.DC.CursorMaxPos
	sz := bb.GetSize()
	ItemSizeVec(&sz, style.FramePadding.y)
	window.DC.CursorMaxPos = backup_cursor_max_pos

	if !ItemAdd(&bb, id, nil, 0) {
		if want_clip_rect {
			PopClipRect()
		}
		window.DC.CursorPos = backup_main_cursor_pos
		return tab_contents_visible
	}

	var f = ImGuiButtonFlags_PressedOnClick
	if is_tab_button {
		f = ImGuiButtonFlags_PressedOnClickRelease
	}

	// Click to Select a tab
	var button_flags = f | ImGuiButtonFlags_AllowItemOverlap
	if g.DragDropActive {
		button_flags |= ImGuiButtonFlags_PressedOnDragDropHold
	}
	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, button_flags)
	if pressed && !is_tab_button {
		tab_bar.NextSelectedTabId = id
	}

	// Allow the close button to overlap unless we are dragging (in which case we don't want any overlapping tabs to be hovered)
	if g.ActiveId != id {
		SetItemAllowOverlap()
	}

	// Drag and drop: re-order tabs
	if held && !tab_appearing && IsMouseDragging(0, -1) {
		if !g.DragDropActive && (tab_bar.Flags&ImGuiTabBarFlags_Reorderable) != 0 {
			// While moving a tab it will jump on the other side of the mouse, so we also test for MouseDelta.x
			if g.IO.MouseDelta.x < 0.0 && g.IO.MousePos.x < bb.Min.x {
				TabBarQueueReorderFromMousePos(tab_bar, tab, g.IO.MousePos)
			} else if g.IO.MouseDelta.x > 0.0 && g.IO.MousePos.x > bb.Max.x {
				TabBarQueueReorderFromMousePos(tab_bar, tab, g.IO.MousePos)
			}
		}
	}

	var c = ImGuiCol_TabUnfocused
	if held || hovered {
		c = ImGuiCol_TabHovered
	} else {
		if tab_contents_visible {
			if tab_bar_focused {
				c = ImGuiCol_TabActive
			} else {
				c = ImGuiCol_TabUnfocusedActive
			}
		} else {
			if tab_bar_focused {
				c = ImGuiCol_Tab
			} else {
				c = ImGuiCol_TabUnfocused
			}
		}
	}

	// Render tab shape
	var display_draw_list = window.DrawList
	var tab_col = GetColorU32FromID(c, 1)
	TabItemBackground(display_draw_list, &bb, flags, tab_col)
	RenderNavHighlight(&bb, id, 0)

	// Select with right mouse button. This is so the common idiom for context menu automatically highlight the current widget.
	var hovered_unblocked = IsItemHovered(ImGuiHoveredFlags_AllowWhenBlockedByPopup)
	if hovered_unblocked && (IsMouseClicked(1, false) || IsMouseReleased(1)) {
		if !is_tab_button {
			tab_bar.NextSelectedTabId = id
		}
	}

	if (tab_bar.Flags & ImGuiTabBarFlags_NoCloseWithMiddleMouseButton) != 0 {
		flags |= ImGuiTabItemFlags_NoCloseWithMiddleMouseButton
	}

	// Render tab label, process close button
	var close_button_id uint = 0
	if p_open != nil {
		close_button_id = GetIDWithSeed("#CLOSE", id)
	}
	var just_closed bool
	var text_clipped bool
	TabItemLabelAndCloseButton(display_draw_list, &bb, flags, tab_bar.FramePadding, label, id, close_button_id, tab_contents_visible, &just_closed, &text_clipped)
	if just_closed && p_open != nil {
		*p_open = false
		TabBarCloseTab(tab_bar, tab)
	}

	// Restore main window position so user can draw there
	if want_clip_rect {
		PopClipRect()
	}
	window.DC.CursorPos = backup_main_cursor_pos

	// Tooltip
	// (Won't work over the close button because ItemOverlap systems messes up with HoveredIdTimer. seems ok)
	// (We test IsItemHovered() to discard e.g. when another item is active or drag and drop over the tab bar, which g.HoveredId ignores)
	// FIXME: This is a mess.
	// FIXME: We may want disabled tab to still display the tooltip?
	if text_clipped && g.HoveredId == id && !held && g.HoveredIdNotActiveTimer > g.TooltipSlowDelay && IsItemHovered(0) {
		if (tab_bar.Flags&ImGuiTabBarFlags_NoTooltip) == 0 && (tab.Flags&ImGuiTabItemFlags_NoTooltip) == 0 {
			SetTooltip("%v", label)
		}
	}

	IM_ASSERT(!is_tab_button || !(tab_bar.SelectedTabId == tab.ID && is_tab_button)) // TabItemButton should not be selected
	if is_tab_button {
		return pressed
	}
	return tab_contents_visible
}

func TabItemCalcSize(label string, has_close_button bool) ImVec2 {
	var g = GImGui
	var label_size = CalcTextSize(label, true, -1)
	var size = ImVec2{label_size.x + g.Style.FramePadding.x, label_size.y + g.Style.FramePadding.y*2.0}
	if has_close_button {
		size.x += g.Style.FramePadding.x + (g.Style.ItemInnerSpacing.x + g.FontSize) // We use Y intentionally to fit the close button circle.
	} else {
		size.x += g.Style.FramePadding.x + 1.0
	}
	return ImVec2{ImMin(size.x, TabBarCalcMaxTabWidth()), size.y}
}

func TabItemBackground(draw_list *ImDrawList, bb *ImRect, flags ImGuiTabItemFlags, col ImU32) {
	// While rendering tabs, we trim 1 pixel off the top of our bounding box so they can fit within a regular frame height while looking "detached" from it.
	var g = GImGui
	var width = bb.GetWidth()
	IM_ASSERT(width > 0.0)

	var r = g.Style.TabRounding
	if (flags & ImGuiTabItemFlags_Button) != 0 {
		r = g.Style.FrameRounding
	}

	var rounding = ImMax(0.0, ImMin(r, width*0.5-1.0))
	var y1 = bb.Min.y + 1.0
	var y2 = bb.Max.y - 1.0
	draw_list.PathLineTo(ImVec2{bb.Min.x, y2})
	draw_list.PathArcToFast(ImVec2{bb.Min.x + rounding, y1 + rounding}, rounding, 6, 9)
	draw_list.PathArcToFast(ImVec2{bb.Max.x - rounding, y1 + rounding}, rounding, 9, 12)
	draw_list.PathLineTo(ImVec2{bb.Max.x, y2})
	draw_list.PathFillConvex(col)
	if g.Style.TabBorderSize > 0.0 {
		draw_list.PathLineTo(ImVec2{bb.Min.x + 0.5, y2})
		draw_list.PathArcToFast(ImVec2{bb.Min.x + rounding + 0.5, y1 + rounding + 0.5}, rounding, 6, 9)
		draw_list.PathArcToFast(ImVec2{bb.Max.x - rounding - 0.5, y1 + rounding + 0.5}, rounding, 9, 12)
		draw_list.PathLineTo(ImVec2{bb.Max.x - 0.5, y2})
		draw_list.PathStroke(GetColorU32FromID(ImGuiCol_Border, 1), 0, g.Style.TabBorderSize)
	}
}

// Render text label (with custom clipping) + Unsaved Document marker + Close Button logic
// We tend to lock style.FramePadding for a given tab-bar, hence the 'frame_padding' parameter.
func TabItemLabelAndCloseButton(draw_list *ImDrawList, bb *ImRect, flags ImGuiTabItemFlags, frame_padding ImVec2, label string, tab_id ImGuiID, close_button_id ImGuiID, is_contents_visible bool, out_just_closed *bool, out_text_clipped *bool) {
	var g = GImGui
	var label_size = CalcTextSize(label, true, -1)

	if out_just_closed != nil {
		*out_just_closed = false
	}
	if out_text_clipped != nil {
		*out_text_clipped = false
	}

	if bb.GetWidth() <= 1.0 {
		return
	}

	// Render text label (with clipping + alpha gradient) + unsaved marker
	var text_pixel_clip_bb = ImRect{ImVec2{bb.Min.x + frame_padding.x, bb.Min.y + frame_padding.y}, ImVec2{bb.Max.x - frame_padding.x, bb.Max.y}}
	var text_ellipsis_clip_bb = text_pixel_clip_bb

	// Return clipped state ignoring the close button
	if out_text_clipped != nil {
		*out_text_clipped = (text_ellipsis_clip_bb.Min.x + label_size.x) > text_pixel_clip_bb.Max.x
		//draw_list.AddCircle(text_ellipsis_clip_bb.Min, 3.0f, *out_text_clipped ? IM_COL32(255, 0, 0, 255) : IM_COL32(0, 255, 0, 255));
	}

	var button_sz = g.FontSize
	var button_pos = ImVec2{ImMax(bb.Min.x, bb.Max.x-frame_padding.x*2.0-button_sz), bb.Min.y}

	// Close Button & Unsaved Marker
	// We are relying on a subtle and confusing distinction between 'hovered' and 'g.HoveredId' which happens because we are using ImGuiButtonFlags_AllowOverlapMode + SetItemAllowOverlap()
	//  'hovered' will be true when hovering the Tab but NOT when hovering the close button
	//  'g.HoveredId==id' will be true when hovering the Tab including when hovering the close button
	//  'g.ActiveId==close_button_id' will be true when we are holding on the close button, in which case both hovered booleans are false
	var close_button_pressed = false
	var close_button_visible = false
	if close_button_id != 0 {
		if is_contents_visible || bb.GetWidth() >= ImMax(button_sz, g.Style.TabMinWidthForCloseButton) {
			if g.HoveredId == tab_id || g.HoveredId == close_button_id || g.ActiveId == tab_id || g.ActiveId == close_button_id {
				close_button_visible = true
			}
		}
	}
	var unsaved_marker_visible = (flags&ImGuiTabItemFlags_UnsavedDocument) != 0 && (button_pos.x+button_sz <= bb.Max.x)

	if close_button_visible {
		var last_item_backup = g.LastItemData
		PushStyleVec(ImGuiStyleVar_FramePadding, frame_padding)
		if CloseButton(close_button_id, &button_pos) {
			close_button_pressed = true
		}
		PopStyleVar(1)
		g.LastItemData = last_item_backup

		// Close with middle mouse button
		if flags&ImGuiTabItemFlags_NoCloseWithMiddleMouseButton == 0 && IsMouseClicked(2, false) {
			close_button_pressed = true
		}
	} else if unsaved_marker_visible {
		var bullet_bb = ImRect{button_pos, button_pos.Add(ImVec2{button_sz, button_sz}).Add(g.Style.FramePadding.Scale(2.0))}
		RenderBullet(draw_list, bullet_bb.GetCenter(), GetColorU32FromID(ImGuiCol_Text, 1))
	}

	// This is all rather complicated
	// (the main idea is that because the close button only appears on hover, we don't want it to alter the ellipsis position)
	// FIXME: if FramePadding is noticeably large, ellipsis_max_x will be wrong here (e.g. #3497), maybe for consistency that parameter of RenderTextEllipsis() shouldn't exist..
	var ellipsis_max_x = bb.Max.x - 1.0
	if close_button_visible {
		ellipsis_max_x = text_pixel_clip_bb.Max.x
	}
	if close_button_visible || unsaved_marker_visible {
		if close_button_visible {
			text_pixel_clip_bb.Max.x -= (button_sz)
		} else {
			text_pixel_clip_bb.Max.x -= (button_sz * 0.80)
		}
		if unsaved_marker_visible {
			text_ellipsis_clip_bb.Max.x -= (button_sz * 0.80)
		}
		ellipsis_max_x = text_pixel_clip_bb.Max.x
	}
	RenderTextEllipsis(draw_list, &text_ellipsis_clip_bb.Min, &text_ellipsis_clip_bb.Max, text_pixel_clip_bb.Max.x, ellipsis_max_x, label, &label_size)

	if out_just_closed != nil {
		*out_just_closed = close_button_pressed
	}
}

// Note: we may scroll to tab that are not selected! e.g. using keyboard arrow keys
func TabBarScrollToTab(tab_bar *ImGuiTabBar, tab_id ImGuiID, sections [3]ImGuiTabBarSection) {
	var tab = TabBarFindTabByID(tab_bar, tab_id)
	if tab == nil {
		return
	}
	if (tab.Flags & ImGuiTabItemFlags_SectionMask_) != 0 {
		return
	}

	var g = GImGui
	var margin = g.FontSize * 1.0 // When to scroll to make Tab N+1 visible always make a bit of N visible to suggest more scrolling area (since we don't have a scrollbar)
	var order = tab_bar.GetTabOrder(tab)

	// Scrolling happens only in the central section (leading/trailing sections are not scrolling)
	// FIXME: This is all confusing.
	var scrollable_width = tab_bar.BarRect.GetWidth() - sections[0].Width - sections[2].Width - sections[1].Spacing

	// We make all tabs positions all relative Sections[0].Width to make code simpler
	var tab_x1 = tab.Offset - sections[0].Width
	if order > sections[0].TabCount-1 {
		tab_x1 += -margin
	}
	var tab_x2 = tab.Offset - sections[0].Width + tab.Width + 1
	if order+1 < int(len(tab_bar.Tabs))-sections[2].TabCount {
		tab_x2 += margin - 1
	}
	tab_bar.ScrollingTargetDistToVisibility = 0.0
	if tab_bar.ScrollingTarget > tab_x1 || (tab_x2-tab_x1 >= scrollable_width) {
		// Scroll to the left
		tab_bar.ScrollingTargetDistToVisibility = ImMax(tab_bar.ScrollingAnim-tab_x2, 0.0)
		tab_bar.ScrollingTarget = tab_x1
	} else if tab_bar.ScrollingTarget < tab_x2-scrollable_width {
		// Scroll to the right
		tab_bar.ScrollingTargetDistToVisibility = ImMax((tab_x1-scrollable_width)-tab_bar.ScrollingAnim, 0.0)
		tab_bar.ScrollingTarget = tab_x2 - scrollable_width
	}
}

func TabBarTabListPopupButton(tab_bar *ImGuiTabBar) *ImGuiTabItem {
	var g = GImGui
	var window = g.CurrentWindow

	// We use g.Style.FramePadding.y to match the square ArrowButton size
	var tab_list_popup_button_width = g.FontSize + g.Style.FramePadding.y
	var backup_cursor_pos = window.DC.CursorPos
	window.DC.CursorPos = ImVec2{tab_bar.BarRect.Min.x - g.Style.FramePadding.y, tab_bar.BarRect.Min.y}
	tab_bar.BarRect.Min.x += tab_list_popup_button_width

	var arrow_col = g.Style.Colors[ImGuiCol_Text]
	arrow_col.w *= 0.5
	PushStyleColorVec(ImGuiCol_Text, &arrow_col)
	PushStyleColorVec(ImGuiCol_Button, &ImVec4{})
	var open = BeginCombo("##v", "", ImGuiComboFlags_NoPreview|ImGuiComboFlags_HeightLargest)
	PopStyleColor(2)

	var tab_to_select *ImGuiTabItem = nil
	if open {
		for tab_n := int(0); tab_n < int(len(tab_bar.Tabs)); tab_n++ {
			var tab = &tab_bar.Tabs[tab_n]
			if (tab.Flags & ImGuiTabItemFlags_Button) != 0 {
				continue
			}

			var tab_name = tab_bar.GetTabName(tab)
			if Selectable(tab_name, tab_bar.SelectedTabId == tab.ID, 0, ImVec2{}) {
				tab_to_select = tab
			}
		}
		EndCombo()
	}

	window.DC.CursorPos = backup_cursor_pos
	return tab_to_select
}

func TabBarScrollClamp(tab_bar *ImGuiTabBar, scrolling float) float {
	scrolling = ImMin(scrolling, tab_bar.WidthAllTabs-tab_bar.BarRect.GetWidth())
	return ImMax(scrolling, 0.0)
}

// Dockables uses Name/ID in the global namespace. Non-dockable items use the ID stack.
func TabBarCalcTabID(tab_bar *ImGuiTabBar, label string) ImU32 {
	if (tab_bar.Flags & ImGuiTabBarFlags_DockNode) != 0 {
		var id = ImHashStr(label, uintptr(len(label)), 0)
		KeepAliveID(id)
		return id
	} else {
		var window = GImGui.CurrentWindow
		return window.GetIDs(label)
	}
}

func TabBarCalcMaxTabWidth() float {
	var g = GImGui
	return g.FontSize * 20.0
}

func TabBarScrollingButtons(tab_bar *ImGuiTabBar) *ImGuiTabItem {
	var g = GImGui
	var window = g.CurrentWindow

	var arrow_button_size = ImVec2{g.FontSize - 2.0, g.FontSize + g.Style.FramePadding.y*2.0}
	var scrolling_buttons_width = arrow_button_size.x * 2.0

	var backup_cursor_pos = window.DC.CursorPos
	//window.DrawList.AddRect(ImVec2(tab_bar.BarRect.Max.x - scrolling_buttons_width, tab_bar.BarRect.Min.y), ImVec2(tab_bar.BarRect.Max.x, tab_bar.BarRect.Max.y), IM_COL32(255,0,0,255));

	var select_dir int = 0
	var arrow_col = g.Style.Colors[ImGuiCol_Text]
	arrow_col.w *= 0.5

	PushStyleColorVec(ImGuiCol_Text, &arrow_col)
	PushStyleColorVec(ImGuiCol_Button, &ImVec4{})
	var backup_repeat_delay = g.IO.KeyRepeatDelay
	var backup_repeat_rate = g.IO.KeyRepeatRate
	g.IO.KeyRepeatDelay = 0.250
	g.IO.KeyRepeatRate = 0.200
	var x = ImMax(tab_bar.BarRect.Min.x, tab_bar.BarRect.Max.x-scrolling_buttons_width)
	window.DC.CursorPos = ImVec2{x, tab_bar.BarRect.Min.y}
	if ArrowButtonEx("##<", ImGuiDir_Left, arrow_button_size, ImGuiButtonFlags_PressedOnClick|ImGuiButtonFlags_Repeat) {
		select_dir = -1
	}
	window.DC.CursorPos = ImVec2{x + arrow_button_size.x, tab_bar.BarRect.Min.y}
	if ArrowButtonEx("##>", ImGuiDir_Right, arrow_button_size, ImGuiButtonFlags_PressedOnClick|ImGuiButtonFlags_Repeat) {
		select_dir = +1
	}
	PopStyleColor(2)
	g.IO.KeyRepeatRate = backup_repeat_rate
	g.IO.KeyRepeatDelay = backup_repeat_delay

	var tab_to_scroll_to *ImGuiTabItem = nil
	if select_dir != 0 {
		if tab_item := TabBarFindTabByID(tab_bar, tab_bar.SelectedTabId); tab_item != nil {
			var selected_order = tab_bar.GetTabOrder(tab_item)
			var target_order = selected_order + select_dir

			// Skip tab item buttons until another tab item is found or end is reached
			for tab_to_scroll_to == nil {
				// If we are at the end of the list, still scroll to make our tab visible
				if target_order >= 0 && target_order < int(len(tab_bar.Tabs)) {
					tab_to_scroll_to = &tab_bar.Tabs[target_order]
				} else {
					tab_to_scroll_to = &tab_bar.Tabs[selected_order]
				}

				// Cross through buttons
				// (even if first/last item is a button, return it so we can update the scroll)
				if tab_to_scroll_to.Flags&ImGuiTabItemFlags_Button != 0 {
					target_order += select_dir
					selected_order += select_dir
					if !(target_order < 0 || target_order >= int(len(tab_bar.Tabs))) {
						tab_to_scroll_to = nil
					}
				}
			}
		}
	}
	window.DC.CursorPos = backup_cursor_pos
	tab_bar.BarRect.Max.x -= scrolling_buttons_width + 1.0

	return tab_to_scroll_to
}

// This is called only once a frame before by the first call to ItemTab()
// The reason we're not calling it in BeginTabBar() is to leave a chance to the user to call the SetTabItemClosed() functions.
func TabBarLayout(tab_bar *ImGuiTabBar) {
	var g = GImGui
	tab_bar.WantLayout = false

	// Garbage collect by compacting list
	// Detect if we need to sort out tab list (e.g. in rare case where a tab changed section)
	var tab_dst_n int = 0
	var need_sort_by_section = false
	var sections [3]ImGuiTabBarSection // Layout sections: Leading, Central, Trailing
	for tab_src_n := int(0); tab_src_n < int(len(tab_bar.Tabs)); tab_src_n++ {
		var tab = &tab_bar.Tabs[tab_src_n]
		if tab.LastFrameVisible < tab_bar.PrevFrameVisible || tab.WantClose {
			// Remove tab
			if tab_bar.VisibleTabId == tab.ID {
				tab_bar.VisibleTabId = 0
			}
			if tab_bar.SelectedTabId == tab.ID {
				tab_bar.SelectedTabId = 0
			}
			if tab_bar.NextSelectedTabId == tab.ID {
				tab_bar.NextSelectedTabId = 0
			}
			continue
		}
		if tab_dst_n != tab_src_n {
			tab_bar.Tabs[tab_dst_n] = tab_bar.Tabs[tab_src_n]
		}

		tab = &tab_bar.Tabs[tab_dst_n]
		tab.IndexDuringLayout = (int)(tab_dst_n)

		// We will need sorting if tabs have changed section (e.g. moved from one of Leading/Central/Trailing to another)
		var curr_tab_section_n = TabItemGetSectionIdx(tab)
		if tab_dst_n > 0 {
			var prev_tab = &tab_bar.Tabs[tab_dst_n-1]
			var prev_tab_section_n = TabItemGetSectionIdx(prev_tab)
			if curr_tab_section_n == 0 && prev_tab_section_n != 0 {
				need_sort_by_section = true
			}
			if prev_tab_section_n == 2 && curr_tab_section_n != 2 {
				need_sort_by_section = true
			}
		}

		sections[curr_tab_section_n].TabCount++
		tab_dst_n++
	}
	if int(len(tab_bar.Tabs)) != tab_dst_n {
		//tab_bar.Tabs.resize(tab_dst_n);
		if tab_dst_n > int(len(tab_bar.Tabs)) {
			tab_bar.Tabs = append(tab_bar.Tabs, make([]ImGuiTabItem, tab_dst_n-int(len(tab_bar.Tabs)))...)
		} else {
			tab_bar.Tabs = tab_bar.Tabs[:tab_dst_n]
		}
	}

	if need_sort_by_section {
		sort.Slice(tab_bar.Tabs, func(i, j golang.Int) bool {
			var a = &tab_bar.Tabs[i]
			var b = &tab_bar.Tabs[j]
			var a_section = TabItemGetSectionIdx(a)
			var b_section = TabItemGetSectionIdx(b)
			if a_section != b_section {
				return a_section < b_section
			}
			return a.IndexDuringLayout < b.IndexDuringLayout
		})
	}

	// Calculate spacing between sections
	sections[0].Spacing = 0
	if sections[0].TabCount > 0 && (sections[1].TabCount+sections[2].TabCount) > 0 {
		sections[0].Spacing = g.Style.ItemInnerSpacing.x
	}
	sections[1].Spacing = 0
	if sections[1].TabCount > 0 && sections[2].TabCount > 0 {
		sections[1].Spacing = g.Style.ItemInnerSpacing.x
	}

	// Setup next selected tab
	var scroll_to_tab_id ImGuiID = 0
	if tab_bar.NextSelectedTabId != 0 {
		tab_bar.SelectedTabId = tab_bar.NextSelectedTabId
		tab_bar.NextSelectedTabId = 0
		scroll_to_tab_id = tab_bar.SelectedTabId
	}

	// Process order change request (we could probably process it when requested but it's just saner to do it in a single spot).
	if tab_bar.ReorderRequestTabId != 0 {
		if TabBarProcessReorder(tab_bar) {
			if tab_bar.ReorderRequestTabId == tab_bar.SelectedTabId {
				scroll_to_tab_id = tab_bar.ReorderRequestTabId
			}
		}
		tab_bar.ReorderRequestTabId = 0
	}

	// Tab List Popup (will alter tab_bar.BarRect and therefore the available width!)
	var tab_list_popup_button = (tab_bar.Flags & ImGuiTabBarFlags_TabListPopupButton) != 0
	if tab_list_popup_button {
		if tab_to_select := TabBarTabListPopupButton(tab_bar); tab_to_select != nil { // NB: Will alter BarRect.Min.x!
			scroll_to_tab_id = tab_to_select.ID
			tab_bar.SelectedTabId = tab_to_select.ID
		}
	}

	// Leading/Trailing tabs will be shrink only if central one aren't visible anymore, so layout the shrink data as: leading, trailing, central
	// (whereas our tabs are stored as: leading, central, trailing)
	var shrink_buffer_indexes = [3]int{0, sections[0].TabCount + sections[2].TabCount, sections[0].TabCount}

	//resize to tab_bar.Tabs.Size
	if len(tab_bar.Tabs) > len(g.ShrinkWidthBuffer) {
		g.ShrinkWidthBuffer = append(g.ShrinkWidthBuffer, make([]ImGuiShrinkWidthItem, len(tab_bar.Tabs)-len(g.ShrinkWidthBuffer))...)
	} else {
		g.ShrinkWidthBuffer = g.ShrinkWidthBuffer[:len(tab_bar.Tabs)]
	}

	// Compute ideal tabs widths + store them into shrink buffer
	var most_recently_selected_tab *ImGuiTabItem = nil
	var curr_section_n int = -1
	var found_selected_tab_id = false
	for tab_n := int(0); tab_n < int(len(tab_bar.Tabs)); tab_n++ {
		var tab = &tab_bar.Tabs[tab_n]
		IM_ASSERT(tab.LastFrameVisible >= tab_bar.PrevFrameVisible)

		if (most_recently_selected_tab == nil || most_recently_selected_tab.LastFrameSelected < tab.LastFrameSelected) && tab.Flags&ImGuiTabItemFlags_Button == 0 {
			most_recently_selected_tab = tab
		}
		if tab.ID == tab_bar.SelectedTabId {
			found_selected_tab_id = true
		}
		if scroll_to_tab_id == 0 && g.NavJustMovedToId == tab.ID {
			scroll_to_tab_id = tab.ID
		}

		// Refresh tab width immediately, otherwise changes of style e.g. style.FramePadding.x would noticeably lag in the tab bar.
		// Additionally, when using TabBarAddTab() to manipulate tab bar order we occasionally insert new tabs that don't have a width yet,
		// and we cannot wait for the next BeginTabItem() call. We cannot compute this width within TabBarAddTab() because font size depends on the active window.
		var tab_name = tab_bar.GetTabName(tab)
		var has_close_button = true
		if tab.Flags&ImGuiTabItemFlags_NoCloseButton != 0 {
			has_close_button = false
		}
		tab.ContentWidth = TabItemCalcSize(tab_name, has_close_button).x

		var section_n = TabItemGetSectionIdx(tab)
		var section = &sections[section_n]
		section.Width += tab.ContentWidth
		if section_n == curr_section_n {
			section.Width += g.Style.ItemInnerSpacing.x
		}
		curr_section_n = section_n

		// Store data so we can build an array sorted by width if we need to shrink tabs down
		var shrink_buffer_index = shrink_buffer_indexes[section_n]
		shrink_buffer_indexes[section_n]++
		g.ShrinkWidthBuffer[shrink_buffer_index].Index = tab_n
		g.ShrinkWidthBuffer[shrink_buffer_index].Width = tab.ContentWidth

		IM_ASSERT(tab.ContentWidth > 0.0)
		tab.Width = tab.ContentWidth
	}

	// Compute total ideal width (used for e.g. auto-resizing a window)
	tab_bar.WidthAllTabsIdeal = 0.0
	for section_n := 0; section_n < 3; section_n++ {
		tab_bar.WidthAllTabsIdeal += sections[section_n].Width + sections[section_n].Spacing
	}

	// Horizontal scrolling buttons
	// (note that TabBarScrollButtons() will alter BarRect.Max.x)
	if (tab_bar.WidthAllTabsIdeal > tab_bar.BarRect.GetWidth() && len(tab_bar.Tabs) > 1) && tab_bar.Flags&ImGuiTabBarFlags_NoTabListScrollingButtons == 0 && (tab_bar.Flags&ImGuiTabBarFlags_FittingPolicyScroll) != 0 {
		if scroll_and_select_tab := TabBarScrollingButtons(tab_bar); scroll_and_select_tab != nil {
			scroll_to_tab_id = scroll_and_select_tab.ID
			if (scroll_and_select_tab.Flags & ImGuiTabItemFlags_Button) == 0 {
				tab_bar.SelectedTabId = scroll_to_tab_id
			}
		}
	}

	// Shrink widths if full tabs don't fit in their allocated space
	var section_0_w = sections[0].Width + sections[0].Spacing
	var section_1_w = sections[1].Width + sections[1].Spacing
	var section_2_w = sections[2].Width + sections[2].Spacing
	var central_section_is_visible = (section_0_w + section_2_w) < tab_bar.BarRect.GetWidth()
	var width_excess float
	if central_section_is_visible {
		width_excess = ImMax(section_1_w-(tab_bar.BarRect.GetWidth()-section_0_w-section_2_w), 0.0) // Excess used to shrink central section
	} else {
		width_excess = (section_0_w + section_2_w) - tab_bar.BarRect.GetWidth() // Excess used to shrink leading/trailing section
	}

	// With ImGuiTabBarFlags_FittingPolicyScroll policy, we will only shrink leading/trailing if the central section is not visible anymore
	if width_excess > 0.0 && ((tab_bar.Flags&ImGuiTabBarFlags_FittingPolicyResizeDown) != 0 || !central_section_is_visible) {
		var shrink_data_count int
		var shrink_data_offset int
		if central_section_is_visible {
			shrink_data_count = sections[1].TabCount
			shrink_data_offset = sections[0].TabCount + sections[2].TabCount
		} else {
			shrink_data_count = sections[0].TabCount + sections[2].TabCount
		}
		ShrinkWidths(g.ShrinkWidthBuffer[shrink_data_offset:], shrink_data_count, width_excess)

		// Apply shrunk values into tabs and sections
		for tab_n := shrink_data_offset; tab_n < shrink_data_offset+shrink_data_count; tab_n++ {
			var tab = &tab_bar.Tabs[g.ShrinkWidthBuffer[tab_n].Index]
			var shrinked_width = IM_FLOOR(g.ShrinkWidthBuffer[tab_n].Width)
			if shrinked_width < 0.0 {
				continue
			}

			var section_n = TabItemGetSectionIdx(tab)
			sections[section_n].Width -= (tab.Width - shrinked_width)
			tab.Width = shrinked_width
		}
	}

	// Layout all active tabs
	var section_tab_index int = 0
	var tab_offset float
	tab_bar.WidthAllTabs = 0.0
	for section_n := int(0); section_n < 3; section_n++ {
		var section = &sections[section_n]
		if section_n == 2 {
			tab_offset = ImMin(ImMax(0.0, tab_bar.BarRect.GetWidth()-section.Width), tab_offset)
		}

		for tab_n := int(0); tab_n < section.TabCount; tab_n++ {
			var tab = &tab_bar.Tabs[section_tab_index+tab_n]
			tab.Offset = tab_offset
			tab_offset += tab.Width
			if tab_n < section.TabCount-1 {
				tab_offset += g.Style.ItemInnerSpacing.x
			}
		}
		tab_bar.WidthAllTabs += ImMax(section.Width+section.Spacing, 0.0)
		tab_offset += section.Spacing
		section_tab_index += section.TabCount
	}

	// If we have lost the selected tab, select the next most recently active one
	if !found_selected_tab_id {
		tab_bar.SelectedTabId = 0
	}
	if tab_bar.SelectedTabId == 0 && tab_bar.NextSelectedTabId == 0 && most_recently_selected_tab != nil {
		scroll_to_tab_id = most_recently_selected_tab.ID
		tab_bar.SelectedTabId = most_recently_selected_tab.ID
	}

	// Lock in visible tab
	tab_bar.VisibleTabId = tab_bar.SelectedTabId
	tab_bar.VisibleTabWasSubmitted = false

	// Update scrolling
	if scroll_to_tab_id != 0 {
		TabBarScrollToTab(tab_bar, scroll_to_tab_id, sections)
	}
	tab_bar.ScrollingAnim = TabBarScrollClamp(tab_bar, tab_bar.ScrollingAnim)
	tab_bar.ScrollingTarget = TabBarScrollClamp(tab_bar, tab_bar.ScrollingTarget)
	if tab_bar.ScrollingAnim != tab_bar.ScrollingTarget {
		// Scrolling speed adjust itself so we can always reach our target in 1/3 seconds.
		// Teleport if we are aiming far off the visible line
		tab_bar.ScrollingSpeed = ImMax(tab_bar.ScrollingSpeed, 70.0*g.FontSize)
		tab_bar.ScrollingSpeed = ImMax(tab_bar.ScrollingSpeed, ImFabs(tab_bar.ScrollingTarget-tab_bar.ScrollingAnim)/0.3)
		var teleport = (tab_bar.PrevFrameVisible+1 < g.FrameCount) || (tab_bar.ScrollingTargetDistToVisibility > 10.0*g.FontSize)
		if teleport {
			tab_bar.ScrollingAnim = tab_bar.ScrollingTarget
		} else {
			tab_bar.ScrollingAnim = ImLinearSweep(tab_bar.ScrollingAnim, tab_bar.ScrollingTarget, g.IO.DeltaTime*tab_bar.ScrollingSpeed)
		}
	} else {
		tab_bar.ScrollingSpeed = 0.0
	}
	tab_bar.ScrollingRectMinX = tab_bar.BarRect.Min.x + sections[0].Width + sections[0].Spacing
	tab_bar.ScrollingRectMaxX = tab_bar.BarRect.Max.x - sections[2].Width - sections[1].Spacing

	// Clear name buffers
	if (tab_bar.Flags & ImGuiTabBarFlags_DockNode) == 0 {
		tab_bar.TabsNames = tab_bar.TabsNames[:0]
	}

	// Actual layout in host window (we don't do it in BeginTabBar() so as not to waste an extra frame)
	var window = g.CurrentWindow
	window.DC.CursorPos = tab_bar.BarRect.Min
	ItemSizeVec(&ImVec2{tab_bar.WidthAllTabs, tab_bar.BarRect.GetHeight()}, tab_bar.FramePadding.y)
	window.DC.IdealMaxPos.x = ImMax(window.DC.IdealMaxPos.x, tab_bar.BarRect.Min.x+tab_bar.WidthAllTabsIdeal)
}
