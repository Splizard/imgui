package imgui

//-----------------------------------------------------------------------------
// [SECTION] Flags & Enumerations
//-----------------------------------------------------------------------------

// Flags for ImGui::Begin()
const (
	ImGuiWindowFlags_None                      ImGuiWindowFlags = 0
	ImGuiWindowFlags_NoTitleBar                ImGuiWindowFlags = 1 << 0  // Disable title-bar
	ImGuiWindowFlags_NoResize                  ImGuiWindowFlags = 1 << 1  // Disable user resizing with the lower-right grip
	ImGuiWindowFlags_NoMove                    ImGuiWindowFlags = 1 << 2  // Disable user moving the window
	ImGuiWindowFlags_NoScrollbar               ImGuiWindowFlags = 1 << 3  // Disable scrollbars (window can still scroll with mouse or programmatically)
	ImGuiWindowFlags_NoScrollWithMouse         ImGuiWindowFlags = 1 << 4  // Disable user vertically scrolling with mouse wheel. On child window, mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	ImGuiWindowFlags_NoCollapse                ImGuiWindowFlags = 1 << 5  // Disable user collapsing window by double-clicking on it
	ImGuiWindowFlags_AlwaysAutoResize          ImGuiWindowFlags = 1 << 6  // Resize every window to its content every frame
	ImGuiWindowFlags_NoBackground              ImGuiWindowFlags = 1 << 7  // Disable drawing background color (WindowBg, etc.) and outside border. Similar as using SetNextWindowBgAlpha(0.0f).
	ImGuiWindowFlags_NoSavedSettings           ImGuiWindowFlags = 1 << 8  // Never load/save settings in .ini file
	ImGuiWindowFlags_NoMouseInputs             ImGuiWindowFlags = 1 << 9  // Disable catching mouse, hovering test with pass through.
	ImGuiWindowFlags_MenuBar                   ImGuiWindowFlags = 1 << 10 // Has a menu-bar
	ImGuiWindowFlags_HorizontalScrollbar       ImGuiWindowFlags = 1 << 11 // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width,0.0f)); prior to calling Begin() to specify width. Read code in imgui_demo in the "Horizontal Scrolling" section.
	ImGuiWindowFlags_NoFocusOnAppearing        ImGuiWindowFlags = 1 << 12 // Disable taking focus when transitioning from hidden to visible state
	ImGuiWindowFlags_NoBringToFrontOnFocus     ImGuiWindowFlags = 1 << 13 // Disable bringing window to front when taking focus (e.g. clicking on it or programmatically giving it focus)
	ImGuiWindowFlags_AlwaysVerticalScrollbar   ImGuiWindowFlags = 1 << 14 // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	ImGuiWindowFlags_AlwaysHorizontalScrollbar ImGuiWindowFlags = 1 << 15 // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	ImGuiWindowFlags_AlwaysUseWindowPadding    ImGuiWindowFlags = 1 << 16 // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows, because more convenient)
	ImGuiWindowFlags_NoNavInputs               ImGuiWindowFlags = 1 << 18 // No gamepad/keyboard navigation within the window
	ImGuiWindowFlags_NoNavFocus                ImGuiWindowFlags = 1 << 19 // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	ImGuiWindowFlags_UnsavedDocument           ImGuiWindowFlags = 1 << 20 // Display a dot next to the title. When used in a tab/docking context, tab is selected when clicking the X + closure is not assumed (will wait for user to stop submitting the tab). Otherwise closure is assumed when pressing the X, so if you keep submitting the tab may reappear at end of tab bar.
	ImGuiWindowFlags_NoNav                     ImGuiWindowFlags = ImGuiWindowFlags_NoNavInputs | ImGuiWindowFlags_NoNavFocus
	ImGuiWindowFlags_NoDecoration              ImGuiWindowFlags = ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoScrollbar | ImGuiWindowFlags_NoCollapse
	ImGuiWindowFlags_NoInputs                  ImGuiWindowFlags = ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs | ImGuiWindowFlags_NoNavFocus

	// [Internal]
	ImGuiWindowFlags_NavFlattened ImGuiWindowFlags = 1 << 23 // [BETA] Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	ImGuiWindowFlags_ChildWindow  ImGuiWindowFlags = 1 << 24 // Don't use! For internal use by BeginChild()
	ImGuiWindowFlags_Tooltip      ImGuiWindowFlags = 1 << 25 // Don't use! For internal use by BeginTooltip()
	ImGuiWindowFlags_Popup        ImGuiWindowFlags = 1 << 26 // Don't use! For internal use by BeginPopup()
	ImGuiWindowFlags_Modal        ImGuiWindowFlags = 1 << 27 // Don't use! For internal use by BeginPopupModal()
	ImGuiWindowFlags_ChildMenu    ImGuiWindowFlags = 1 << 28 // Don't use! For internal use by BeginMenu()

	// [Obsolete]
	//ImGuiWindowFlags_ResizeFromAnySide    = 1 << 17// --> Set io.ConfigWindowsResizeFromEdges=true and make sure mouse cursors are supported by backend (io.BackendFlags & ImGuiBackendFlags_HasMouseCursors)
)

// Flags for ImGui::InputText()
const (
	ImGuiInputTextFlags_None                ImGuiInputTextFlags = 0
	ImGuiInputTextFlags_CharsDecimal        ImGuiInputTextFlags = 1 << 0  // Allow 0123456789.+-*/
	ImGuiInputTextFlags_CharsHexadecimal    ImGuiInputTextFlags = 1 << 1  // Allow 0123456789ABCDEFabcdef
	ImGuiInputTextFlags_CharsUppercase      ImGuiInputTextFlags = 1 << 2  // Turn a..z into A..Z
	ImGuiInputTextFlags_CharsNoBlank        ImGuiInputTextFlags = 1 << 3  // Filter out spaces, tabs
	ImGuiInputTextFlags_AutoSelectAll       ImGuiInputTextFlags = 1 << 4  // Select entire text when first taking mouse focus
	ImGuiInputTextFlags_EnterReturnsTrue    ImGuiInputTextFlags = 1 << 5  // Return 'true' when Enter is pressed (as opposed to every time the value was modified). Consider looking at the IsItemDeactivatedAfterEdit() function.
	ImGuiInputTextFlags_CallbackCompletion  ImGuiInputTextFlags = 1 << 6  // Callback on pressing TAB (for completion handling)
	ImGuiInputTextFlags_CallbackHistory     ImGuiInputTextFlags = 1 << 7  // Callback on pressing Up/Down arrows (for history handling)
	ImGuiInputTextFlags_CallbackAlways      ImGuiInputTextFlags = 1 << 8  // Callback on each iteration. User code may query cursor position, modify text buffer.
	ImGuiInputTextFlags_CallbackCharFilter  ImGuiInputTextFlags = 1 << 9  // Callback on character inputs to replace or discard them. Modify 'EventChar' to replace or discard, or return 1 in callback to discard.
	ImGuiInputTextFlags_AllowTabInput       ImGuiInputTextFlags = 1 << 10 // Pressing TAB input a '\t' character into the text field
	ImGuiInputTextFlags_CtrlEnterForNewLine ImGuiInputTextFlags = 1 << 11 // In multi-line mode, unfocus with Enter, add new line with Ctrl+Enter (default is opposite: unfocus with Ctrl+Enter, add line with Enter).
	ImGuiInputTextFlags_NoHorizontalScroll  ImGuiInputTextFlags = 1 << 12 // Disable following the cursor horizontally
	ImGuiInputTextFlags_AlwaysOverwrite     ImGuiInputTextFlags = 1 << 13 // Overwrite mode
	ImGuiInputTextFlags_ReadOnly            ImGuiInputTextFlags = 1 << 14 // Read-only mode
	ImGuiInputTextFlags_Password            ImGuiInputTextFlags = 1 << 15 // Password mode, display all characters as '*'
	ImGuiInputTextFlags_NoUndoRedo          ImGuiInputTextFlags = 1 << 16 // Disable undo/redo. Note that input text owns the text data while active, if you want to provide your own undo/redo stack you need e.g. to call ClearActiveID().
	ImGuiInputTextFlags_CharsScientific     ImGuiInputTextFlags = 1 << 17 // Allow 0123456789.+-*/eE (Scientific notation input)
	ImGuiInputTextFlags_CallbackResize      ImGuiInputTextFlags = 1 << 18 // Callback on buffer capacity changes request (beyond 'buf_size' parameter value), allowing the string to grow. Notify when the string wants to be resized (for string types which hold a cache of their Size). You will be provided a new BufSize in the callback and NEED to honor it. (see misc/cpp/imgui_stdlib.h for an example of using this)
	ImGuiInputTextFlags_CallbackEdit        ImGuiInputTextFlags = 1 << 19 // Callback on any edit (note that InputText() already returns true on edit, the callback is useful mainly to manipulate the underlying buffer while focus is active)
)

// Flags for ImGui::TreeNodeEx(), ImGui::CollapsingHeader*()
const (
	ImGuiTreeNodeFlags_None                 ImGuiTreeNodeFlags = 0
	ImGuiTreeNodeFlags_Selected             ImGuiTreeNodeFlags = 1 << 0  // Draw as selected
	ImGuiTreeNodeFlags_Framed               ImGuiTreeNodeFlags = 1 << 1  // Draw frame with background (e.g. for CollapsingHeader)
	ImGuiTreeNodeFlags_AllowItemOverlap     ImGuiTreeNodeFlags = 1 << 2  // Hit testing to allow subsequent widgets to overlap this one
	ImGuiTreeNodeFlags_NoTreePushOnOpen     ImGuiTreeNodeFlags = 1 << 3  // Don't do a TreePush() when open (e.g. for CollapsingHeader) = no extra indent nor pushing on ID stack
	ImGuiTreeNodeFlags_NoAutoOpenOnLog      ImGuiTreeNodeFlags = 1 << 4  // Don't automatically and temporarily open node when Logging is active (by default logging will automatically open tree nodes)
	ImGuiTreeNodeFlags_DefaultOpen          ImGuiTreeNodeFlags = 1 << 5  // Default node to be open
	ImGuiTreeNodeFlags_OpenOnDoubleClick    ImGuiTreeNodeFlags = 1 << 6  // Need double-click to open node
	ImGuiTreeNodeFlags_OpenOnArrow          ImGuiTreeNodeFlags = 1 << 7  // Only open when clicking on the arrow part. If ImGuiTreeNodeFlags_OpenOnDoubleClick is also set, single-click arrow or double-click all box to open.
	ImGuiTreeNodeFlags_Leaf                 ImGuiTreeNodeFlags = 1 << 8  // No collapsing, no arrow (use as a convenience for leaf nodes).
	ImGuiTreeNodeFlags_Bullet               ImGuiTreeNodeFlags = 1 << 9  // Display a bullet instead of arrow
	ImGuiTreeNodeFlags_FramePadding         ImGuiTreeNodeFlags = 1 << 10 // Use FramePadding (even for an unframed text node) to vertically align text baseline to regular widget height. Equivalent to calling AlignTextToFramePadding().
	ImGuiTreeNodeFlags_SpanAvailWidth       ImGuiTreeNodeFlags = 1 << 11 // Extend hit box to the right-most edge, even if not framed. This is not the default in order to allow adding other items on the same line. In the future we may refactor the hit system to be front-to-back, allowing natural overlaps and then this can become the default.
	ImGuiTreeNodeFlags_SpanFullWidth        ImGuiTreeNodeFlags = 1 << 12 // Extend hit box to the left-most and right-most edges (bypass the indented area).
	ImGuiTreeNodeFlags_NavLeftJumpsBackHere ImGuiTreeNodeFlags = 1 << 13 // (WIP) Nav: left direction may move to this TreeNode() from any of its child (items submitted between TreeNode and TreePop)
	//ImGuiTreeNodeFlags_NoScrollOnOpen     = 1 << 14// FIXME: TODO: Disable automatic scroll on TreePop() if node got just open and contents is not visible
	ImGuiTreeNodeFlags_CollapsingHeader ImGuiTreeNodeFlags = ImGuiTreeNodeFlags_Framed | ImGuiTreeNodeFlags_NoTreePushOnOpen | ImGuiTreeNodeFlags_NoAutoOpenOnLog
)

// Flags for OpenPopup*(), BeginPopupContext*(), IsPopupOpen() functions.
// - To be backward compatible with older API which took an 'int mouse_button = 1' argument, we need to treat
//   small flags values as a mouse button index, so we encode the mouse button in the first few bits of the flags.
//   It is therefore guaranteed to be legal to pass a mouse button index in ImGuiPopupFlags.
// - For the same reason, we exceptionally default the ImGuiPopupFlags argument of BeginPopupContextXXX functions to 1 instead of 0.
//   IMPORTANT: because the default parameter is 1 (==ImGuiPopupFlags_MouseButtonRight), if you rely on the default parameter
//   and want to another another flag, you need to pass in the ImGuiPopupFlags_MouseButtonRight flag.
// - Multiple buttons currently cannot be combined/or-ed in those functions (we could allow it later).
const (
	ImGuiPopupFlags_None                    ImGuiPopupFlags = 0
	ImGuiPopupFlags_MouseButtonLeft         ImGuiPopupFlags = 0 // For BeginPopupContext*(): open on Left Mouse release. Guaranteed to always be == 0 (same as ImGuiMouseButton_Left)
	ImGuiPopupFlags_MouseButtonRight        ImGuiPopupFlags = 1 // For BeginPopupContext*(): open on Right Mouse release. Guaranteed to always be == 1 (same as ImGuiMouseButton_Right)
	ImGuiPopupFlags_MouseButtonMiddle       ImGuiPopupFlags = 2 // For BeginPopupContext*(): open on Middle Mouse release. Guaranteed to always be == 2 (same as ImGuiMouseButton_Middle)
	ImGuiPopupFlags_MouseButtonMask_        ImGuiPopupFlags = 0x1F
	ImGuiPopupFlags_MouseButtonDefault_     ImGuiPopupFlags = 1
	ImGuiPopupFlags_NoOpenOverExistingPopup ImGuiPopupFlags = 1 << 5 // For OpenPopup*(), BeginPopupContext*(): don't open if there's already a popup at the same level of the popup stack
	ImGuiPopupFlags_NoOpenOverItems         ImGuiPopupFlags = 1 << 6 // For BeginPopupContextWindow(): don't return true when hovering items, only when hovering empty space
	ImGuiPopupFlags_AnyPopupId              ImGuiPopupFlags = 1 << 7 // For IsPopupOpen(): ignore the ImGuiID parameter and test for any popup.
	ImGuiPopupFlags_AnyPopupLevel           ImGuiPopupFlags = 1 << 8 // For IsPopupOpen(): search/test at any level of the popup stack (default test in the current level)
	ImGuiPopupFlags_AnyPopup                ImGuiPopupFlags = ImGuiPopupFlags_AnyPopupId | ImGuiPopupFlags_AnyPopupLevel
)

// Flags for ImGui::Selectable()
const (
	ImGuiSelectableFlags_None             ImGuiSelectableFlags = 0
	ImGuiSelectableFlags_DontClosePopups  ImGuiSelectableFlags = 1 << 0 // Clicking this don't close parent popup window
	ImGuiSelectableFlags_SpanAllColumns   ImGuiSelectableFlags = 1 << 1 // Selectable frame can span all columns (text will still fit in current column)
	ImGuiSelectableFlags_AllowDoubleClick ImGuiSelectableFlags = 1 << 2 // Generate press events on double clicks too
	ImGuiSelectableFlags_Disabled         ImGuiSelectableFlags = 1 << 3 // Cannot be selected, display grayed out text
	ImGuiSelectableFlags_AllowItemOverlap ImGuiSelectableFlags = 1 << 4 // (WIP) Hit testing to allow subsequent widgets to overlap this one
)

// Flags for ImGui::BeginCombo()
const (
	ImGuiComboFlags_None           ImGuiComboFlags = 0
	ImGuiComboFlags_PopupAlignLeft ImGuiComboFlags = 1 << 0 // Align the popup toward the left by default
	ImGuiComboFlags_HeightSmall    ImGuiComboFlags = 1 << 1 // Max ~4 items visible. Tip: If you want your combo popup to be a specific size you can use SetNextWindowSizeConstraints() prior to calling BeginCombo()
	ImGuiComboFlags_HeightRegular  ImGuiComboFlags = 1 << 2 // Max ~8 items visible (default)
	ImGuiComboFlags_HeightLarge    ImGuiComboFlags = 1 << 3 // Max ~20 items visible
	ImGuiComboFlags_HeightLargest  ImGuiComboFlags = 1 << 4 // As many fitting items as possible
	ImGuiComboFlags_NoArrowButton  ImGuiComboFlags = 1 << 5 // Display on the preview box without the square arrow button
	ImGuiComboFlags_NoPreview      ImGuiComboFlags = 1 << 6 // Display only a square arrow button
	ImGuiComboFlags_HeightMask_    ImGuiComboFlags = ImGuiComboFlags_HeightSmall | ImGuiComboFlags_HeightRegular | ImGuiComboFlags_HeightLarge | ImGuiComboFlags_HeightLargest
)

// Flags for ImGui::BeginTabBar()
const (
	ImGuiTabBarFlags_None                         ImGuiTabBarFlags = 0
	ImGuiTabBarFlags_Reorderable                  ImGuiTabBarFlags = 1 << 0 // Allow manually dragging tabs to re-order them + New tabs are appended at the end of list
	ImGuiTabBarFlags_AutoSelectNewTabs            ImGuiTabBarFlags = 1 << 1 // Automatically select new tabs when they appear
	ImGuiTabBarFlags_TabListPopupButton           ImGuiTabBarFlags = 1 << 2 // Disable buttons to open the tab list popup
	ImGuiTabBarFlags_NoCloseWithMiddleMouseButton ImGuiTabBarFlags = 1 << 3 // Disable behavior of closing tabs (that are submitted with p_open != NULL) with middle mouse button. You can still repro this behavior on user's side with if (IsItemHovered() && IsMouseClicked(2)) *p_open = false.
	ImGuiTabBarFlags_NoTabListScrollingButtons    ImGuiTabBarFlags = 1 << 4 // Disable scrolling buttons (apply when fitting policy is ImGuiTabBarFlags_FittingPolicyScroll)
	ImGuiTabBarFlags_NoTooltip                    ImGuiTabBarFlags = 1 << 5 // Disable tooltips when hovering a tab
	ImGuiTabBarFlags_FittingPolicyResizeDown      ImGuiTabBarFlags = 1 << 6 // Resize tabs when they don't fit
	ImGuiTabBarFlags_FittingPolicyScroll          ImGuiTabBarFlags = 1 << 7 // Add scroll buttons when tabs don't fit
	ImGuiTabBarFlags_FittingPolicyMask_           ImGuiTabBarFlags = ImGuiTabBarFlags_FittingPolicyResizeDown | ImGuiTabBarFlags_FittingPolicyScroll
	ImGuiTabBarFlags_FittingPolicyDefault_        ImGuiTabBarFlags = ImGuiTabBarFlags_FittingPolicyResizeDown
)

// Flags for ImGui::BeginTabItem()
const (
	ImGuiTabItemFlags_None                         ImGuiTabItemFlags = 0
	ImGuiTabItemFlags_UnsavedDocument              ImGuiTabItemFlags = 1 << 0 // Display a dot next to the title + tab is selected when clicking the X + closure is not assumed (will wait for user to stop submitting the tab). Otherwise closure is assumed when pressing the X, so if you keep submitting the tab may reappear at end of tab bar.
	ImGuiTabItemFlags_SetSelected                  ImGuiTabItemFlags = 1 << 1 // Trigger flag to programmatically make the tab selected when calling BeginTabItem()
	ImGuiTabItemFlags_NoCloseWithMiddleMouseButton ImGuiTabItemFlags = 1 << 2 // Disable behavior of closing tabs (that are submitted with p_open != NULL) with middle mouse button. You can still repro this behavior on user's side with if (IsItemHovered() && IsMouseClicked(2)) *p_open = false.
	ImGuiTabItemFlags_NoPushId                     ImGuiTabItemFlags = 1 << 3 // Don't call PushID(tab->ID)/PopID() on BeginTabItem()/EndTabItem()
	ImGuiTabItemFlags_NoTooltip                    ImGuiTabItemFlags = 1 << 4 // Disable tooltip for the given tab
	ImGuiTabItemFlags_NoReorder                    ImGuiTabItemFlags = 1 << 5 // Disable reordering this tab or having another tab cross over this tab
	ImGuiTabItemFlags_Leading                      ImGuiTabItemFlags = 1 << 6 // Enforce the tab position to the left of the tab bar (after the tab list popup button)
	ImGuiTabItemFlags_Trailing                     ImGuiTabItemFlags = 1 << 7 // Enforce the tab position to the right of the tab bar (before the scrolling buttons)
)

// Flags for ImGui::BeginTable()
// [BETA API] API may evolve slightly! If you use this, please update to the next version when it comes out!
// - Important! Sizing policies have complex and subtle side effects, more so than you would expect.
//   Read comments/demos carefully + experiment with live demos to get acquainted with them.
// - The DEFAULT sizing policies are:
//    - Default to ImGuiTableFlags_SizingFixedFit    if ScrollX is on, or if host window has ImGuiWindowFlags_AlwaysAutoResize.
//    - Default to ImGuiTableFlags_SizingStretchSame if ScrollX is off.
// - When ScrollX is off:
//    - Table defaults to ImGuiTableFlags_SizingStretchSame -> all Columns defaults to ImGuiTableColumnFlags_WidthStretch with same weight.
//    - Columns sizing policy allowed: Stretch (default), Fixed/Auto.
//    - Fixed Columns will generally obtain their requested width (unless the table cannot fit them all).
//    - Stretch Columns will share the remaining width.
//    - Mixed Fixed/Stretch columns is possible but has various side-effects on resizing behaviors.
//      The typical use of mixing sizing policies is: any number of LEADING Fixed columns, followed by one or two TRAILING Stretch columns.
//      (this is because the visible order of columns have subtle but necessary effects on how they react to manual resizing).
// - When ScrollX is on:
//    - Table defaults to ImGuiTableFlags_SizingFixedFit -> all Columns defaults to ImGuiTableColumnFlags_WidthFixed
//    - Columns sizing policy allowed: Fixed/Auto mostly.
//    - Fixed Columns can be enlarged as needed. Table will show an horizontal scrollbar if needed.
//    - When using auto-resizing (non-resizable) fixed columns, querying the content width to use item right-alignment e.g. SetNextItemWidth(-FLT_MIN) doesn't make sense, would create a feedback loop.
//    - Using Stretch columns OFTEN DOES NOT MAKE SENSE if ScrollX is on, UNLESS you have specified a value for 'inner_width' in BeginTable().
//      If you specify a value for 'inner_width' then effectively the scrolling space is known and Stretch or mixed Fixed/Stretch columns become meaningful again.
// - Read on documentation at the top of imgui_tables.cpp for details.
const (

	// Features
	ImGuiTableFlags_None              ImGuiTableFlags = 0
	ImGuiTableFlags_Resizable         ImGuiTableFlags = 1 << 0 // Enable resizing columns.
	ImGuiTableFlags_Reorderable       ImGuiTableFlags = 1 << 1 // Enable reordering columns in header row (need calling TableSetupColumn() + TableHeadersRow() to display headers)
	ImGuiTableFlags_Hideable          ImGuiTableFlags = 1 << 2 // Enable hiding/disabling columns in context menu.
	ImGuiTableFlags_Sortable          ImGuiTableFlags = 1 << 3 // Enable sorting. Call TableGetSortSpecs() to obtain sort specs. Also see ImGuiTableFlags_SortMulti and ImGuiTableFlags_SortTristate.
	ImGuiTableFlags_NoSavedSettings   ImGuiTableFlags = 1 << 4 // Disable persisting columns order, width and sort settings in the .ini file.
	ImGuiTableFlags_ContextMenuInBody ImGuiTableFlags = 1 << 5 // Right-click on columns body/contents will display table context menu. By default it is available in TableHeadersRow().
	// Decorations
	ImGuiTableFlags_RowBg                      ImGuiTableFlags = 1 << 6                                                        // Set each RowBg color with ImGuiCol_TableRowBg or ImGuiCol_TableRowBgAlt (equivalent of calling TableSetBgColor with ImGuiTableBgFlags_RowBg0 on each row manually)
	ImGuiTableFlags_BordersInnerH              ImGuiTableFlags = 1 << 7                                                        // Draw horizontal borders between rows.
	ImGuiTableFlags_BordersOuterH              ImGuiTableFlags = 1 << 8                                                        // Draw horizontal borders at the top and bottom.
	ImGuiTableFlags_BordersInnerV              ImGuiTableFlags = 1 << 9                                                        // Draw vertical borders between columns.
	ImGuiTableFlags_BordersOuterV              ImGuiTableFlags = 1 << 10                                                       // Draw vertical borders on the left and right sides.
	ImGuiTableFlags_BordersH                   ImGuiTableFlags = ImGuiTableFlags_BordersInnerH | ImGuiTableFlags_BordersOuterH // Draw horizontal borders.
	ImGuiTableFlags_BordersV                   ImGuiTableFlags = ImGuiTableFlags_BordersInnerV | ImGuiTableFlags_BordersOuterV // Draw vertical borders.
	ImGuiTableFlags_BordersInner               ImGuiTableFlags = ImGuiTableFlags_BordersInnerV | ImGuiTableFlags_BordersInnerH // Draw inner borders.
	ImGuiTableFlags_BordersOuter               ImGuiTableFlags = ImGuiTableFlags_BordersOuterV | ImGuiTableFlags_BordersOuterH // Draw outer borders.
	ImGuiTableFlags_Borders                    ImGuiTableFlags = ImGuiTableFlags_BordersInner | ImGuiTableFlags_BordersOuter   // Draw all borders.
	ImGuiTableFlags_NoBordersInBody            ImGuiTableFlags = 1 << 11                                                       // [ALPHA] Disable vertical borders in columns Body (borders will always appears in Headers). -> May move to style
	ImGuiTableFlags_NoBordersInBodyUntilResize ImGuiTableFlags = 1 << 12                                                       // [ALPHA] Disable vertical borders in columns Body until hovered for resize (borders will always appears in Headers). -> May move to style
	// Sizing Policy (read above for defaults)
	ImGuiTableFlags_SizingFixedFit    ImGuiTableFlags = 1 << 13 // Columns default to _WidthFixed or _WidthAuto (if resizable or not resizable), matching contents width.
	ImGuiTableFlags_SizingFixedSame   ImGuiTableFlags = 2 << 13 // Columns default to _WidthFixed or _WidthAuto (if resizable or not resizable), matching the maximum contents width of all columns. Implicitly enable ImGuiTableFlags_NoKeepColumnsVisible.
	ImGuiTableFlags_SizingStretchProp ImGuiTableFlags = 3 << 13 // Columns default to _WidthStretch with default weights proportional to each columns contents widths.
	ImGuiTableFlags_SizingStretchSame ImGuiTableFlags = 4 << 13 // Columns default to _WidthStretch with default weights all equal, unless overridden by TableSetupColumn().
	// Sizing Extra Options
	ImGuiTableFlags_NoHostExtendX        ImGuiTableFlags = 1 << 16 // Make outer width auto-fit to columns, overriding outer_size.x value. Only available when ScrollX/ScrollY are disabled and Stretch columns are not used.
	ImGuiTableFlags_NoHostExtendY        ImGuiTableFlags = 1 << 17 // Make outer height stop exactly at outer_size.y (prevent auto-extending table past the limit). Only available when ScrollX/ScrollY are disabled. Data below the limit will be clipped and not visible.
	ImGuiTableFlags_NoKeepColumnsVisible ImGuiTableFlags = 1 << 18 // Disable keeping column always minimally visible when ScrollX is off and table gets too small. Not recommended if columns are resizable.
	ImGuiTableFlags_PreciseWidths        ImGuiTableFlags = 1 << 19 // Disable distributing remainder width to stretched columns (width allocation on a 100-wide table with 3 columns: Without this flag: 33,33,34. With this flag: 33,33,33). With larger number of columns, resizing will appear to be less smooth.
	// Clipping
	ImGuiTableFlags_NoClip ImGuiTableFlags = 1 << 20 // Disable clipping rectangle for every individual columns (reduce draw command count, items will be able to overflow into other columns). Generally incompatible with TableSetupScrollFreeze().
	// Padding
	ImGuiTableFlags_PadOuterX   ImGuiTableFlags = 1 << 21 // Default if BordersOuterV is on. Enable outer-most padding. Generally desirable if you have headers.
	ImGuiTableFlags_NoPadOuterX ImGuiTableFlags = 1 << 22 // Default if BordersOuterV is off. Disable outer-most padding.
	ImGuiTableFlags_NoPadInnerX ImGuiTableFlags = 1 << 23 // Disable inner padding between columns (double inner padding if BordersOuterV is on, single inner padding if BordersOuterV is off).
	// Scrolling
	ImGuiTableFlags_ScrollX ImGuiTableFlags = 1 << 24 // Enable horizontal scrolling. Require 'outer_size' parameter of BeginTable() to specify the container size. Changes default sizing policy. Because this create a child window, ScrollY is currently generally recommended when using ScrollX.
	ImGuiTableFlags_ScrollY ImGuiTableFlags = 1 << 25 // Enable vertical scrolling. Require 'outer_size' parameter of BeginTable() to specify the container size.
	// Sorting
	ImGuiTableFlags_SortMulti    ImGuiTableFlags = 1 << 26 // Hold shift when clicking headers to sort on multiple column. TableGetSortSpecs() may return specs where (SpecsCount > 1).
	ImGuiTableFlags_SortTristate ImGuiTableFlags = 1 << 27 // Allow no sorting, disable default sorting. TableGetSortSpecs() may return specs where (SpecsCount == 0).

	// [Internal] Combinations and masks
	ImGuiTableFlags_SizingMask_ ImGuiTableFlags = ImGuiTableFlags_SizingFixedFit | ImGuiTableFlags_SizingFixedSame | ImGuiTableFlags_SizingStretchProp | ImGuiTableFlags_SizingStretchSame
)

// Flags for ImGui::TableSetupColumn()
const (

	// Input configuration flags
	ImGuiTableColumnFlags_None                 ImGuiTableColumnFlags = 0
	ImGuiTableColumnFlags_Disabled             ImGuiTableColumnFlags = 1 << 0  // Overriding/master disable flag: hide column, won't show in context menu (unlike calling TableSetColumnEnabled() which manipulates the user accessible state)
	ImGuiTableColumnFlags_DefaultHide          ImGuiTableColumnFlags = 1 << 1  // Default as a hidden/disabled column.
	ImGuiTableColumnFlags_DefaultSort          ImGuiTableColumnFlags = 1 << 2  // Default as a sorting column.
	ImGuiTableColumnFlags_WidthStretch         ImGuiTableColumnFlags = 1 << 3  // Column will stretch. Preferable with horizontal scrolling disabled (default if table sizing policy is _SizingStretchSame or _SizingStretchProp).
	ImGuiTableColumnFlags_WidthFixed           ImGuiTableColumnFlags = 1 << 4  // Column will not stretch. Preferable with horizontal scrolling enabled (default if table sizing policy is _SizingFixedFit and table is resizable).
	ImGuiTableColumnFlags_NoResize             ImGuiTableColumnFlags = 1 << 5  // Disable manual resizing.
	ImGuiTableColumnFlags_NoReorder            ImGuiTableColumnFlags = 1 << 6  // Disable manual reordering this column, this will also prevent other columns from crossing over this column.
	ImGuiTableColumnFlags_NoHide               ImGuiTableColumnFlags = 1 << 7  // Disable ability to hide/disable this column.
	ImGuiTableColumnFlags_NoClip               ImGuiTableColumnFlags = 1 << 8  // Disable clipping for this column (all NoClip columns will render in a same draw command).
	ImGuiTableColumnFlags_NoSort               ImGuiTableColumnFlags = 1 << 9  // Disable ability to sort on this field (even if ImGuiTableFlags_Sortable is set on the table).
	ImGuiTableColumnFlags_NoSortAscending      ImGuiTableColumnFlags = 1 << 10 // Disable ability to sort in the ascending direction.
	ImGuiTableColumnFlags_NoSortDescending     ImGuiTableColumnFlags = 1 << 11 // Disable ability to sort in the descending direction.
	ImGuiTableColumnFlags_NoHeaderLabel        ImGuiTableColumnFlags = 1 << 12 // TableHeadersRow() will not submit label for this column. Convenient for some small columns. Name will still appear in context menu.
	ImGuiTableColumnFlags_NoHeaderWidth        ImGuiTableColumnFlags = 1 << 13 // Disable header text width contribution to automatic column width.
	ImGuiTableColumnFlags_PreferSortAscending  ImGuiTableColumnFlags = 1 << 14 // Make the initial sort direction Ascending when first sorting on this column (default).
	ImGuiTableColumnFlags_PreferSortDescending ImGuiTableColumnFlags = 1 << 15 // Make the initial sort direction Descending when first sorting on this column.
	ImGuiTableColumnFlags_IndentEnable         ImGuiTableColumnFlags = 1 << 16 // Use current Indent value when entering cell (default for column 0).
	ImGuiTableColumnFlags_IndentDisable        ImGuiTableColumnFlags = 1 << 17 // Ignore current Indent value when entering cell (default for columns > 0). Indentation changes _within_ the cell will still be honored.

	// Output status flags, read-only via TableGetColumnFlags()
	ImGuiTableColumnFlags_IsEnabled ImGuiTableColumnFlags = 1 << 24 // Status: is enabled == not hidden by user/api (referred to as "Hide" in _DefaultHide and _NoHide) flags.
	ImGuiTableColumnFlags_IsVisible ImGuiTableColumnFlags = 1 << 25 // Status: is visible == is enabled AND not clipped by scrolling.
	ImGuiTableColumnFlags_IsSorted  ImGuiTableColumnFlags = 1 << 26 // Status: is currently part of the sort specs
	ImGuiTableColumnFlags_IsHovered ImGuiTableColumnFlags = 1 << 27 // Status: is hovered by mouse

	// [Internal] Combinations and masks
	ImGuiTableColumnFlags_WidthMask_      ImGuiTableColumnFlags = ImGuiTableColumnFlags_WidthStretch | ImGuiTableColumnFlags_WidthFixed
	ImGuiTableColumnFlags_IndentMask_     ImGuiTableColumnFlags = ImGuiTableColumnFlags_IndentEnable | ImGuiTableColumnFlags_IndentDisable
	ImGuiTableColumnFlags_StatusMask_     ImGuiTableColumnFlags = ImGuiTableColumnFlags_IsEnabled | ImGuiTableColumnFlags_IsVisible | ImGuiTableColumnFlags_IsSorted | ImGuiTableColumnFlags_IsHovered
	ImGuiTableColumnFlags_NoDirectResize_ ImGuiTableColumnFlags = 1 << 30 // [Internal] Disable user resizing this column directly (it may however we resized indirectly from its left edge)
)

// Flags for ImGui::TableNextRow()
const (
	ImGuiTableRowFlags_None    ImGuiTableRowFlags = 0
	ImGuiTableRowFlags_Headers ImGuiTableRowFlags = 1 << 0 // Identify header row (set default background color + width of its contents accounted different for auto column width)
)

// const (
// Background colors are rendering in 3 layers:
//  - Layer 0: draw with RowBg0 color if set, otherwise draw with ColumnBg0 if set.
//  - Layer 1: draw with RowBg1 color if set, otherwise draw with ColumnBg1 if set.
//  - Layer 2: draw with CellBg color if set.
// The purpose of the two row/columns layers is to let you decide if a background color changes should override or blend with the existing color.
// When using ImGuiTableFlags_RowBg on the table, each row has the RowBg0 color automatically set for odd/even rows.
// If you set the color of RowBg0 target, your color will override the existing RowBg0 color.
// If you set the color of RowBg1 or ColumnBg1 target, your color will blend over the RowBg0 color.
const (
	ImGuiTableBgTarget_None   ImGuiTableBgTarget = 0
	ImGuiTableBgTarget_RowBg0 ImGuiTableBgTarget = 1 // Set row background color 0 (generally used for background, automatically set when ImGuiTableFlags_RowBg is used)
	ImGuiTableBgTarget_RowBg1 ImGuiTableBgTarget = 2 // Set row background color 1 (generally used for selection marking)
	ImGuiTableBgTarget_CellBg ImGuiTableBgTarget = 3 // Set cell background color (top-most color)
)

// Flags for ImGui::IsWindowFocused()
const (
	ImGuiFocusedFlags_None                ImGuiFocusedFlags = 0
	ImGuiFocusedFlags_ChildWindows        ImGuiFocusedFlags = 1 << 0 // IsWindowFocused(): Return true if any children of the window is focused
	ImGuiFocusedFlags_RootWindow          ImGuiFocusedFlags = 1 << 1 // IsWindowFocused(): Test from root window (top most parent of the current hierarchy)
	ImGuiFocusedFlags_AnyWindow           ImGuiFocusedFlags = 1 << 2 // IsWindowFocused(): Return true if any window is focused. Important: If you are trying to tell how to dispatch your low-level inputs, do NOT use this. Use 'io.WantCaptureMouse' instead! Please read the FAQ!
	ImGuiFocusedFlags_RootAndChildWindows ImGuiFocusedFlags = ImGuiFocusedFlags_RootWindow | ImGuiFocusedFlags_ChildWindows
)

// Flags for ImGui::IsItemHovered(), ImGui::IsWindowHovered()
// Note: if you are trying to check whether your mouse should be dispatched to Dear ImGui or to your app, you should use 'io.WantCaptureMouse' instead! Please read the FAQ!
// Note: windows with the ImGuiWindowFlags_NoInputs flag are ignored by IsWindowHovered() calls.
const (
	ImGuiHoveredFlags_None                    ImGuiHoveredFlags = 0      // Return true if directly over the item/window, not obstructed by another window, not obstructed by an active popup or modal blocking inputs under them.
	ImGuiHoveredFlags_ChildWindows            ImGuiHoveredFlags = 1 << 0 // IsWindowHovered() only: Return true if any children of the window is hovered
	ImGuiHoveredFlags_RootWindow              ImGuiHoveredFlags = 1 << 1 // IsWindowHovered() only: Test from root window (top most parent of the current hierarchy)
	ImGuiHoveredFlags_AnyWindow               ImGuiHoveredFlags = 1 << 2 // IsWindowHovered() only: Return true if any window is hovered
	ImGuiHoveredFlags_AllowWhenBlockedByPopup ImGuiHoveredFlags = 1 << 3 // Return true even if a popup window is normally blocking access to this item/window
	//ImGuiHoveredFlags_AllowWhenBlockedByModal     = 1 << 4 // Return true even if a modal popup window is normally blocking access to this item/window. FIXME-TODO: Unavailable yet.
	ImGuiHoveredFlags_AllowWhenBlockedByActiveItem ImGuiHoveredFlags = 1 << 5 // Return true even if an active item is blocking access to this item/window. Useful for Drag and Drop patterns.
	ImGuiHoveredFlags_AllowWhenOverlapped          ImGuiHoveredFlags = 1 << 6 // Return true even if the position is obstructed or overlapped by another window
	ImGuiHoveredFlags_AllowWhenDisabled            ImGuiHoveredFlags = 1 << 7 // Return true even if the item is disabled
	ImGuiHoveredFlags_RectOnly                     ImGuiHoveredFlags = ImGuiHoveredFlags_AllowWhenBlockedByPopup | ImGuiHoveredFlags_AllowWhenBlockedByActiveItem | ImGuiHoveredFlags_AllowWhenOverlapped
	ImGuiHoveredFlags_RootAndChildWindows          ImGuiHoveredFlags = ImGuiHoveredFlags_RootWindow | ImGuiHoveredFlags_ChildWindows
)

// Flags for ImGui::BeginDragDropSource(), ImGui::AcceptDragDropPayload()
const (
	ImGuiDragDropFlags_None ImGuiDragDropFlags = 0
	// BeginDragDropSource() flags
	ImGuiDragDropFlags_SourceNoPreviewTooltip   ImGuiDragDropFlags = 1 << 0 // By default, a successful call to BeginDragDropSource opens a tooltip so you can display a preview or description of the source contents. This flag disable this behavior.
	ImGuiDragDropFlags_SourceNoDisableHover     ImGuiDragDropFlags = 1 << 1 // By default, when dragging we clear data so that IsItemHovered() will return false, to avoid subsequent user code submitting tooltips. This flag disable this behavior so you can still call IsItemHovered() on the source item.
	ImGuiDragDropFlags_SourceNoHoldToOpenOthers ImGuiDragDropFlags = 1 << 2 // Disable the behavior that allows to open tree nodes and collapsing header by holding over them while dragging a source item.
	ImGuiDragDropFlags_SourceAllowNullID        ImGuiDragDropFlags = 1 << 3 // Allow items such as Text(), Image() that have no unique identifier to be used as drag source, by manufacturing a temporary identifier based on their window-relative position. This is extremely unusual within the dear imgui ecosystem and so we made it explicit.
	ImGuiDragDropFlags_SourceExtern             ImGuiDragDropFlags = 1 << 4 // External source (from outside of dear imgui), won't attempt to read current item/window info. Will always return true. Only one Extern source can be active simultaneously.
	ImGuiDragDropFlags_SourceAutoExpirePayload  ImGuiDragDropFlags = 1 << 5 // Automatically expire the payload if the source cease to be submitted (otherwise payloads are persisting while being dragged)
	// AcceptDragDropPayload() flags
	ImGuiDragDropFlags_AcceptBeforeDelivery    ImGuiDragDropFlags = 1 << 10                                                                              // AcceptDragDropPayload() will returns true even before the mouse button is released. You can then call IsDelivery() to test if the payload needs to be delivered.
	ImGuiDragDropFlags_AcceptNoDrawDefaultRect ImGuiDragDropFlags = 1 << 11                                                                              // Do not draw the default highlight rectangle when hovering over target.
	ImGuiDragDropFlags_AcceptNoPreviewTooltip  ImGuiDragDropFlags = 1 << 12                                                                              // Request hiding the BeginDragDropSource tooltip from the BeginDragDropTarget site.
	ImGuiDragDropFlags_AcceptPeekOnly          ImGuiDragDropFlags = ImGuiDragDropFlags_AcceptBeforeDelivery | ImGuiDragDropFlags_AcceptNoDrawDefaultRect // For peeking ahead and inspecting the payload before delivery.
)

// Standard Drag and Drop payload types. You can define you own payload types using short strings. Types starting with '_' are defined by Dear ImGui.
const IMGUI_PAYLOAD_TYPE_COLOR_3F = "_COL3F" // float[3]: Standard type for colors, without alpha. User code may use this type.
const IMGUI_PAYLOAD_TYPE_COLOR_4F = "_COL4F" // float[4]: Standard type for colors. User code may use this type.

// A primary data type
const (
	ImGuiDataType_S8     ImGuiDataType = iota // signed char / char (with sensible compilers)
	ImGuiDataType_U8                          // unsigned char
	ImGuiDataType_S16                         // short
	ImGuiDataType_U16                         // unsigned short
	ImGuiDataType_S32                         // int
	ImGuiDataType_U32                         // unsigned int
	ImGuiDataType_S64                         // long long / __int64
	ImGuiDataType_U64                         // unsigned long long / unsigned __int64
	ImGuiDataType_Float                       // float
	ImGuiDataType_Double                      // double
	ImGuiDataType_COUNT
)

// A cardinal direction
const (
	ImGuiDir_None  ImGuiDir = -1
	ImGuiDir_Left  ImGuiDir = 0
	ImGuiDir_Right ImGuiDir = 1
	ImGuiDir_Up    ImGuiDir = 2
	ImGuiDir_Down  ImGuiDir = 3
	ImGuiDir_COUNT
)

// A sorting direction
const (
	ImGuiSortDirection_None       ImGuiSortDirection = 0
	ImGuiSortDirection_Ascending  ImGuiSortDirection = 1 // Ascending = 0->9, A->Z etc.
	ImGuiSortDirection_Descending ImGuiSortDirection = 2 // Descending = 9->0, Z->A etc.
)

// User fill ImGuiIO.KeyMap[] array with indices into the ImGuiIO.KeysDown[512] array
const (
	ImGuiKey_Tab ImGuiKey = iota
	ImGuiKey_LeftArrow
	ImGuiKey_RightArrow
	ImGuiKey_UpArrow
	ImGuiKey_DownArrow
	ImGuiKey_PageUp
	ImGuiKey_PageDown
	ImGuiKey_Home
	ImGuiKey_End
	ImGuiKey_Insert
	ImGuiKey_Delete
	ImGuiKey_Backspace
	ImGuiKey_Space
	ImGuiKey_Enter
	ImGuiKey_Escape
	ImGuiKey_KeyPadEnter
	ImGuiKey_A // for text edit CTRL+A: select all
	ImGuiKey_C // for text edit CTRL+C: copy
	ImGuiKey_V // for text edit CTRL+V: paste
	ImGuiKey_X // for text edit CTRL+X: cut
	ImGuiKey_Y // for text edit CTRL+Y: redo
	ImGuiKey_Z // for text edit CTRL+Z: undo
	ImGuiKey_COUNT
)

// To test io.KeyMods (which is a combination of individual fields io.KeyCtrl, io.KeyShift, io.KeyAlt set by user/backend)
const (
	ImGuiKeyModFlags_None  ImGuiKeyModFlags = 0
	ImGuiKeyModFlags_Ctrl  ImGuiKeyModFlags = 1 << 0
	ImGuiKeyModFlags_Shift ImGuiKeyModFlags = 1 << 1
	ImGuiKeyModFlags_Alt   ImGuiKeyModFlags = 1 << 2
	ImGuiKeyModFlags_Super ImGuiKeyModFlags = 1 << 3
)

// Gamepad/Keyboard navigation
// Keyboard: Set io.ConfigFlags |= ImGuiConfigFlags_NavEnableKeyboard to enable. NewFrame() will automatically fill io.NavInputs[] based on your io.KeysDown[] + io.KeyMap[] arrays.
// Gamepad:  Set io.ConfigFlags |= ImGuiConfigFlags_NavEnableGamepad to enable. Backend: set ImGuiBackendFlags_HasGamepad and fill the io.NavInputs[] fields before calling NewFrame(). Note that io.NavInputs[] is cleared by EndFrame().
// Read instructions in imgui.cpp for more details. Download PNG/PSD at http://dearimgui.org/controls_sheets.
const (

	// Gamepad Mapping
	ImGuiNavInput_Activate    ImGuiNavInput = iota // activate / open / toggle / tweak value       // e.g. Cross  (PS4), A (Xbox), A (Switch), Space (Keyboard)
	ImGuiNavInput_Cancel                           // cancel / close / exit                        // e.g. Circle (PS4), B (Xbox), B (Switch), Escape (Keyboard)
	ImGuiNavInput_Input                            // text input / on-screen keyboard              // e.g. Triang.(PS4), Y (Xbox), X (Switch), Return (Keyboard)
	ImGuiNavInput_Menu                             // tap: toggle menu / hold: focus, move, resize // e.g. Square (PS4), X (Xbox), Y (Switch), Alt (Keyboard)
	ImGuiNavInput_DpadLeft                         // move / tweak / resize window (w/ PadMenu)    // e.g. D-pad Left/Right/Up/Down (Gamepads), Arrow keys (Keyboard)
	ImGuiNavInput_DpadRight                        //
	ImGuiNavInput_DpadUp                           //
	ImGuiNavInput_DpadDown                         //
	ImGuiNavInput_LStickLeft                       // scroll / move window (w/ PadMenu)            // e.g. Left Analog Stick Left/Right/Up/Down
	ImGuiNavInput_LStickRight                      //
	ImGuiNavInput_LStickUp                         //
	ImGuiNavInput_LStickDown                       //
	ImGuiNavInput_FocusPrev                        // next window (w/ PadMenu)                     // e.g. L1 or L2 (PS4), LB or LT (Xbox), L or ZL (Switch)
	ImGuiNavInput_FocusNext                        // prev window (w/ PadMenu)                     // e.g. R1 or R2 (PS4), RB or RT (Xbox), R or ZL (Switch)
	ImGuiNavInput_TweakSlow                        // slower tweaks                                // e.g. L1 or L2 (PS4), LB or LT (Xbox), L or ZL (Switch)
	ImGuiNavInput_TweakFast                        // faster tweaks                                // e.g. R1 or R2 (PS4), RB or RT (Xbox), R or ZL (Switch)

	// [Internal] Don't use directly! This is used internally to differentiate keyboard from gamepad inputs for behaviors that require to differentiate them.
	// Keyboard behavior that have no corresponding gamepad mapping (e.g. CTRL+TAB) will be directly reading from io.KeysDown[] instead of io.NavInputs[].
	ImGuiNavInput_KeyLeft_  // move left                                    // = Arrow keys
	ImGuiNavInput_KeyRight_ // move right
	ImGuiNavInput_KeyUp_    // move up
	ImGuiNavInput_KeyDown_  // move down
	ImGuiNavInput_COUNT
	ImGuiNavInput_InternalStart_ = ImGuiNavInput_KeyLeft_
)

// Configuration flags stored in io.ConfigFlags. Set by user/application.
const (
	ImGuiConfigFlags_None                 ImGuiConfigFlags = 0
	ImGuiConfigFlags_NavEnableKeyboard    ImGuiConfigFlags = 1 << 0 // Master keyboard navigation enable flag. NewFrame() will automatically fill io.NavInputs[] based on io.KeysDown[].
	ImGuiConfigFlags_NavEnableGamepad     ImGuiConfigFlags = 1 << 1 // Master gamepad navigation enable flag. This is mostly to instruct your imgui backend to fill io.NavInputs[]. Backend also needs to set ImGuiBackendFlags_HasGamepad.
	ImGuiConfigFlags_NavEnableSetMousePos ImGuiConfigFlags = 1 << 2 // Instruct navigation to move the mouse cursor. May be useful on TV/console systems where moving a virtual mouse is awkward. Will update io.MousePos and set io.WantSetMousePos=true. If enabled you MUST honor io.WantSetMousePos requests in your backend, otherwise ImGui will react as if the mouse is jumping around back and forth.
	ImGuiConfigFlags_NavNoCaptureKeyboard ImGuiConfigFlags = 1 << 3 // Instruct navigation to not set the io.WantCaptureKeyboard flag when io.NavActive is set.
	ImGuiConfigFlags_NoMouse              ImGuiConfigFlags = 1 << 4 // Instruct imgui to clear mouse position/buttons in NewFrame(). This allows ignoring the mouse information set by the backend.
	ImGuiConfigFlags_NoMouseCursorChange  ImGuiConfigFlags = 1 << 5 // Instruct backend to not alter mouse cursor shape and visibility. Use if the backend cursor changes are interfering with yours and you don't want to use SetMouseCursor() to change mouse cursor. You may want to honor requests from imgui by reading GetMouseCursor() yourself instead.

	// User storage (to allow your backend/engine to communicate to code that may be shared between multiple projects. Those flags are not used by core Dear ImGui)
	ImGuiConfigFlags_IsSRGB        = 1 << 20 // Application is SRGB-aware.
	ImGuiConfigFlags_IsTouchScreen = 1 << 21 // Application is using a touch screen instead of a mouse.
)

// Backend capabilities flags stored in io.BackendFlags. Set by imgui_impl_xxx or custom backend.
const (
	ImGuiBackendFlags_None                 ImGuiBackendFlags = 0
	ImGuiBackendFlags_HasGamepad           ImGuiBackendFlags = 1 << 0 // Backend Platform supports gamepad and currently has one connected.
	ImGuiBackendFlags_HasMouseCursors      ImGuiBackendFlags = 1 << 1 // Backend Platform supports honoring GetMouseCursor() value to change the OS cursor shape.
	ImGuiBackendFlags_HasSetMousePos       ImGuiBackendFlags = 1 << 2 // Backend Platform supports io.WantSetMousePos requests to reposition the OS mouse position (only used if ImGuiConfigFlags_NavEnableSetMousePos is set).
	ImGuiBackendFlags_RendererHasVtxOffset ImGuiBackendFlags = 1 << 3 // Backend Renderer supports ImDrawCmd::VtxOffset. This enables output of large meshes (64K+ vertices) while still using 16-bit indices.
)

// Enumeration for PushStyleColor() / PopStyleColor()
const (
	ImGuiCol_Text ImGuiCol = iota
	ImGuiCol_TextDisabled
	ImGuiCol_WindowBg // Background of normal windows
	ImGuiCol_ChildBg  // Background of child windows
	ImGuiCol_PopupBg  // Background of popups, menus, tooltips windows
	ImGuiCol_Border
	ImGuiCol_BorderShadow
	ImGuiCol_FrameBg // Background of checkbox, radio button, plot, slider, text input
	ImGuiCol_FrameBgHovered
	ImGuiCol_FrameBgActive
	ImGuiCol_TitleBg
	ImGuiCol_TitleBgActive
	ImGuiCol_TitleBgCollapsed
	ImGuiCol_MenuBarBg
	ImGuiCol_ScrollbarBg
	ImGuiCol_ScrollbarGrab
	ImGuiCol_ScrollbarGrabHovered
	ImGuiCol_ScrollbarGrabActive
	ImGuiCol_CheckMark
	ImGuiCol_SliderGrab
	ImGuiCol_SliderGrabActive
	ImGuiCol_Button
	ImGuiCol_ButtonHovered
	ImGuiCol_ButtonActive
	ImGuiCol_Header // Header* colors are used for CollapsingHeader, TreeNode, Selectable, MenuItem
	ImGuiCol_HeaderHovered
	ImGuiCol_HeaderActive
	ImGuiCol_Separator
	ImGuiCol_SeparatorHovered
	ImGuiCol_SeparatorActive
	ImGuiCol_ResizeGrip
	ImGuiCol_ResizeGripHovered
	ImGuiCol_ResizeGripActive
	ImGuiCol_Tab
	ImGuiCol_TabHovered
	ImGuiCol_TabActive
	ImGuiCol_TabUnfocused
	ImGuiCol_TabUnfocusedActive
	ImGuiCol_PlotLines
	ImGuiCol_PlotLinesHovered
	ImGuiCol_PlotHistogram
	ImGuiCol_PlotHistogramHovered
	ImGuiCol_TableHeaderBg     // Table header background
	ImGuiCol_TableBorderStrong // Table outer and header borders (prefer using Alpha=1.0 here)
	ImGuiCol_TableBorderLight  // Table inner borders (prefer using Alpha=1.0 here)
	ImGuiCol_TableRowBg        // Table row background (even rows)
	ImGuiCol_TableRowBgAlt     // Table row background (odd rows)
	ImGuiCol_TextSelectedBg
	ImGuiCol_DragDropTarget
	ImGuiCol_NavHighlight          // Gamepad/keyboard: current highlighted item
	ImGuiCol_NavWindowingHighlight // Highlight window when using CTRL+TAB
	ImGuiCol_NavWindowingDimBg     // Darken/colorize entire screen behind the CTRL+TAB window list, when active
	ImGuiCol_ModalWindowDimBg      // Darken/colorize entire screen behind a modal window, when one is active
	ImGuiCol_COUNT
)

// Enumeration for PushStyleVar() / PopStyleVar() to temporarily modify the ImGuiStyle structure.
// - The const (
//   During initialization or between frames, feel free to just poke into ImGuiStyle directly.
// - Tip: Use your programming IDE navigation facilities on the names in the _second column_ below to find the actual members and their description.
//   In Visual Studio IDE: CTRL+comma ("Edit.NavigateTo") can follow symbols in comments, whereas CTRL+F12 ("Edit.GoToImplementation") cannot.
//   With Visual Assist installed: ALT+G ("VAssistX.GoToImplementation") can also follow symbols in comments.
// - When changing this enum, you need to update the associated internal table GStyleVarInfo[] accordingly. This is where we link const (
const (

	// const (
	ImGuiStyleVar_Alpha               ImGuiStyleVar = iota // float     Alpha
	ImGuiStyleVar_DisabledAlpha                            // float     DisabledAlpha
	ImGuiStyleVar_WindowPadding                            // ImVec2    WindowPadding
	ImGuiStyleVar_WindowRounding                           // float     WindowRounding
	ImGuiStyleVar_WindowBorderSize                         // float     WindowBorderSize
	ImGuiStyleVar_WindowMinSize                            // ImVec2    WindowMinSize
	ImGuiStyleVar_WindowTitleAlign                         // ImVec2    WindowTitleAlign
	ImGuiStyleVar_ChildRounding                            // float     ChildRounding
	ImGuiStyleVar_ChildBorderSize                          // float     ChildBorderSize
	ImGuiStyleVar_PopupRounding                            // float     PopupRounding
	ImGuiStyleVar_PopupBorderSize                          // float     PopupBorderSize
	ImGuiStyleVar_FramePadding                             // ImVec2    FramePadding
	ImGuiStyleVar_FrameRounding                            // float     FrameRounding
	ImGuiStyleVar_FrameBorderSize                          // float     FrameBorderSize
	ImGuiStyleVar_ItemSpacing                              // ImVec2    ItemSpacing
	ImGuiStyleVar_ItemInnerSpacing                         // ImVec2    ItemInnerSpacing
	ImGuiStyleVar_IndentSpacing                            // float     IndentSpacing
	ImGuiStyleVar_CellPadding                              // ImVec2    CellPadding
	ImGuiStyleVar_ScrollbarSize                            // float     ScrollbarSize
	ImGuiStyleVar_ScrollbarRounding                        // float     ScrollbarRounding
	ImGuiStyleVar_GrabMinSize                              // float     GrabMinSize
	ImGuiStyleVar_GrabRounding                             // float     GrabRounding
	ImGuiStyleVar_TabRounding                              // float     TabRounding
	ImGuiStyleVar_ButtonTextAlign                          // ImVec2    ButtonTextAlign
	ImGuiStyleVar_SelectableTextAlign                      // ImVec2    SelectableTextAlign
	ImGuiStyleVar_COUNT
)

// Flags for InvisibleButton() [extended in imgui_internal.h]
const (
	ImGuiButtonFlags_None              ImGuiButtonFlags = 0
	ImGuiButtonFlags_MouseButtonLeft   ImGuiButtonFlags = 1 << 0 // React on left mouse button (default)
	ImGuiButtonFlags_MouseButtonRight  ImGuiButtonFlags = 1 << 1 // React on right mouse button
	ImGuiButtonFlags_MouseButtonMiddle ImGuiButtonFlags = 1 << 2 // React on center mouse button

	// [Internal]
	ImGuiButtonFlags_MouseButtonMask_    ImGuiButtonFlags = ImGuiButtonFlags_MouseButtonLeft | ImGuiButtonFlags_MouseButtonRight | ImGuiButtonFlags_MouseButtonMiddle
	ImGuiButtonFlags_MouseButtonDefault_ ImGuiButtonFlags = ImGuiButtonFlags_MouseButtonLeft
)

// Flags for ColorEdit3() / ColorEdit4() / ColorPicker3() / ColorPicker4() / ColorButton()
const (
	ImGuiColorEditFlags_None           ImGuiColorEditFlags = 0
	ImGuiColorEditFlags_NoAlpha        ImGuiColorEditFlags = 1 << 1  //              // ColorEdit, ColorPicker, ColorButton: ignore Alpha component (will only read 3 components from the input pointer).
	ImGuiColorEditFlags_NoPicker       ImGuiColorEditFlags = 1 << 2  //              // ColorEdit: disable picker when clicking on color square.
	ImGuiColorEditFlags_NoOptions      ImGuiColorEditFlags = 1 << 3  //              // ColorEdit: disable toggling options menu when right-clicking on inputs/small preview.
	ImGuiColorEditFlags_NoSmallPreview ImGuiColorEditFlags = 1 << 4  //              // ColorEdit, ColorPicker: disable color square preview next to the inputs. (e.g. to show only the inputs)
	ImGuiColorEditFlags_NoInputs       ImGuiColorEditFlags = 1 << 5  //              // ColorEdit, ColorPicker: disable inputs sliders/text widgets (e.g. to show only the small preview color square).
	ImGuiColorEditFlags_NoTooltip      ImGuiColorEditFlags = 1 << 6  //              // ColorEdit, ColorPicker, ColorButton: disable tooltip when hovering the preview.
	ImGuiColorEditFlags_NoLabel        ImGuiColorEditFlags = 1 << 7  //              // ColorEdit, ColorPicker: disable display of inline text label (the label is still forwarded to the tooltip and picker).
	ImGuiColorEditFlags_NoSidePreview  ImGuiColorEditFlags = 1 << 8  //              // ColorPicker: disable bigger color preview on right side of the picker, use small color square preview instead.
	ImGuiColorEditFlags_NoDragDrop     ImGuiColorEditFlags = 1 << 9  //              // ColorEdit: disable drag and drop target. ColorButton: disable drag and drop source.
	ImGuiColorEditFlags_NoBorder       ImGuiColorEditFlags = 1 << 10 //              // ColorButton: disable border (which is enforced by default)

	// User Options (right-click on widget to change some of them).
	ImGuiColorEditFlags_AlphaBar         ImGuiColorEditFlags = 1 << 16 //              // ColorEdit, ColorPicker: show vertical alpha bar/gradient in picker.
	ImGuiColorEditFlags_AlphaPreview     ImGuiColorEditFlags = 1 << 17 //              // ColorEdit, ColorPicker, ColorButton: display preview as a transparent color over a checkerboard, instead of opaque.
	ImGuiColorEditFlags_AlphaPreviewHalf ImGuiColorEditFlags = 1 << 18 //              // ColorEdit, ColorPicker, ColorButton: display half opaque / half checkerboard, instead of opaque.
	ImGuiColorEditFlags_HDR              ImGuiColorEditFlags = 1 << 19 //              // (WIP) ColorEdit: Currently only disable 0.0f..1.0f limits in RGBA edition (note: you probably want to use ImGuiColorEditFlags_Float flag as well).
	ImGuiColorEditFlags_DisplayRGB       ImGuiColorEditFlags = 1 << 20 // [Display]    // ColorEdit: override _display_ type among RGB/HSV/Hex. ColorPicker: select any combination using one or more of RGB/HSV/Hex.
	ImGuiColorEditFlags_DisplayHSV       ImGuiColorEditFlags = 1 << 21 // [Display]    // "
	ImGuiColorEditFlags_DisplayHex       ImGuiColorEditFlags = 1 << 22 // [Display]    // "
	ImGuiColorEditFlags_Uint8            ImGuiColorEditFlags = 1 << 23 // [DataType]   // ColorEdit, ColorPicker, ColorButton: _display_ values formatted as 0..255.
	ImGuiColorEditFlags_Float            ImGuiColorEditFlags = 1 << 24 // [DataType]   // ColorEdit, ColorPicker, ColorButton: _display_ values formatted as 0.0f..1.0f floats instead of 0..255 integers. No round-trip of value via integers.
	ImGuiColorEditFlags_PickerHueBar     ImGuiColorEditFlags = 1 << 25 // [Picker]     // ColorPicker: bar for Hue, rectangle for Sat/Value.
	ImGuiColorEditFlags_PickerHueWheel   ImGuiColorEditFlags = 1 << 26 // [Picker]     // ColorPicker: wheel for Hue, triangle for Sat/Value.
	ImGuiColorEditFlags_InputRGB         ImGuiColorEditFlags = 1 << 27 // [Input]      // ColorEdit, ColorPicker: input and output data in RGB format.
	ImGuiColorEditFlags_InputHSV         ImGuiColorEditFlags = 1 << 28 // [Input]      // ColorEdit, ColorPicker: input and output data in HSV format.

	// Defaults Options. You can set application defaults using SetColorEditOptions(). The intent is that you probably don't want to
	// override them in most of your calls. Let the user choose via the option menu and/or call SetColorEditOptions() once during startup.
	ImGuiColorEditFlags_DefaultOptions_ ImGuiColorEditFlags = ImGuiColorEditFlags_Uint8 | ImGuiColorEditFlags_DisplayRGB | ImGuiColorEditFlags_InputRGB | ImGuiColorEditFlags_PickerHueBar

	// [Internal] Masks
	ImGuiColorEditFlags_DisplayMask_  ImGuiColorEditFlags = ImGuiColorEditFlags_DisplayRGB | ImGuiColorEditFlags_DisplayHSV | ImGuiColorEditFlags_DisplayHex
	ImGuiColorEditFlags_DataTypeMask_ ImGuiColorEditFlags = ImGuiColorEditFlags_Uint8 | ImGuiColorEditFlags_Float
	ImGuiColorEditFlags_PickerMask_   ImGuiColorEditFlags = ImGuiColorEditFlags_PickerHueWheel | ImGuiColorEditFlags_PickerHueBar
	ImGuiColorEditFlags_InputMask_    ImGuiColorEditFlags = ImGuiColorEditFlags_InputRGB | ImGuiColorEditFlags_InputHSV
)

// Flags for DragFloat(), DragInt(), SliderFloat(), SliderInt() etc.
// We use the same sets of flags for DragXXX() and SliderXXX() functions as the features are the same and it makes it easier to swap them.
const (
	ImGuiSliderFlags_None            ImGuiSliderFlags = 0
	ImGuiSliderFlags_AlwaysClamp     ImGuiSliderFlags = 1 << 4     // Clamp value to min/max bounds when input manually with CTRL+Click. By default CTRL+Click allows going out of bounds.
	ImGuiSliderFlags_Logarithmic     ImGuiSliderFlags = 1 << 5     // Make the widget logarithmic (linear otherwise). Consider using ImGuiSliderFlags_NoRoundToFormat with this if using a format-string with small amount of digits.
	ImGuiSliderFlags_NoRoundToFormat ImGuiSliderFlags = 1 << 6     // Disable rounding underlying value to match precision of the display format string (e.g. %.3f values are rounded to those 3 digits)
	ImGuiSliderFlags_NoInput         ImGuiSliderFlags = 1 << 7     // Disable CTRL+Click or Enter key allowing to input text directly into the widget
	ImGuiSliderFlags_InvalidMask_    ImGuiSliderFlags = 0x7000000F // [Internal] We treat using those bits as being potentially a 'float power' argument from the previous API that has got miscast to this enum, and will trigger an assert if needed.
)

// Identify a mouse button.
// Those values are guaranteed to be stable and we frequently use 0/1 directly. Named enums provided for convenience.
const (
	ImGuiMouseButton_Left   ImGuiMouseButton = 0
	ImGuiMouseButton_Right  ImGuiMouseButton = 1
	ImGuiMouseButton_Middle ImGuiMouseButton = 2
	ImGuiMouseButton_COUNT  ImGuiMouseButton = 5
)

// Enumeration for GetMouseCursor()
// User code may request backend to display given cursor by calling SetMouseCursor(), which is why we have some cursors that are marked unused here
const (
	ImGuiMouseCursor_None       ImGuiMouseCursor = -1
	ImGuiMouseCursor_Arrow      ImGuiMouseCursor = iota
	ImGuiMouseCursor_TextInput                   // When hovering over InputText, etc.
	ImGuiMouseCursor_ResizeAll                   // (Unused by Dear ImGui functions)
	ImGuiMouseCursor_ResizeNS                    // When hovering over an horizontal border
	ImGuiMouseCursor_ResizeEW                    // When hovering over a vertical border or a column
	ImGuiMouseCursor_ResizeNESW                  // When hovering over the bottom-left corner of a window
	ImGuiMouseCursor_ResizeNWSE                  // When hovering over the bottom-right corner of a window
	ImGuiMouseCursor_Hand                        // (Unused by Dear ImGui functions. Use for e.g. hyperlinks)
	ImGuiMouseCursor_NotAllowed                  // When hovering something with disallowed interaction. Usually a crossed circle.
	ImGuiMouseCursor_COUNT
)

// Enumeration for ImGui::SetWindow***(), SetNextWindow***(), SetNextItem***() functions
// Represent a condition.
// Important: Treat as a regular enum! Do NOT combine multiple values using binary operators! All the functions above treat 0 as a shortcut to ImGuiCond_Always.
const (
	ImGuiCond_None         ImGuiCond = 0      // No condition (always set the variable), same as _Always
	ImGuiCond_Always       ImGuiCond = 1 << 0 // No condition (always set the variable)
	ImGuiCond_Once         ImGuiCond = 1 << 1 // Set the variable once per runtime session (only the first call will succeed)
	ImGuiCond_FirstUseEver ImGuiCond = 1 << 2 // Set the variable if the object/window has no persistently saved data (no entry in .ini file)
	ImGuiCond_Appearing    ImGuiCond = 1 << 3 // Set the variable if the object/window is appearing after being hidden/inactive (or the first time)
)

// Flags for ImDrawList functions
// (Legacy: bit 0 must always correspond to ImDrawFlags_Closed to be backward compatible with old API using a bool. Bits 1..3 must be unused)
const (
	ImDrawFlags_None                    ImDrawFlags = 0
	ImDrawFlags_Closed                  ImDrawFlags = 1 << 0 // PathStroke(), AddPolyline(): specify that shape should be closed (Important: this is always == 1 for legacy reason)
	ImDrawFlags_RoundCornersTopLeft     ImDrawFlags = 1 << 4 // AddRect(), AddRectFilled(), PathRect(): enable rounding top-left corner only (when rounding > 0.0f, we default to all corners). Was 0x01.
	ImDrawFlags_RoundCornersTopRight    ImDrawFlags = 1 << 5 // AddRect(), AddRectFilled(), PathRect(): enable rounding top-right corner only (when rounding > 0.0f, we default to all corners). Was 0x02.
	ImDrawFlags_RoundCornersBottomLeft  ImDrawFlags = 1 << 6 // AddRect(), AddRectFilled(), PathRect(): enable rounding bottom-left corner only (when rounding > 0.0f, we default to all corners). Was 0x04.
	ImDrawFlags_RoundCornersBottomRight ImDrawFlags = 1 << 7 // AddRect(), AddRectFilled(), PathRect(): enable rounding bottom-right corner only (when rounding > 0.0f, we default to all corners). Wax 0x08.
	ImDrawFlags_RoundCornersNone        ImDrawFlags = 1 << 8 // AddRect(), AddRectFilled(), PathRect(): disable rounding on all corners (when rounding > 0.0f). This is NOT zero, NOT an implicit flag!
	ImDrawFlags_RoundCornersTop         ImDrawFlags = ImDrawFlags_RoundCornersTopLeft | ImDrawFlags_RoundCornersTopRight
	ImDrawFlags_RoundCornersBottom      ImDrawFlags = ImDrawFlags_RoundCornersBottomLeft | ImDrawFlags_RoundCornersBottomRight
	ImDrawFlags_RoundCornersLeft        ImDrawFlags = ImDrawFlags_RoundCornersBottomLeft | ImDrawFlags_RoundCornersTopLeft
	ImDrawFlags_RoundCornersRight       ImDrawFlags = ImDrawFlags_RoundCornersBottomRight | ImDrawFlags_RoundCornersTopRight
	ImDrawFlags_RoundCornersAll         ImDrawFlags = ImDrawFlags_RoundCornersTopLeft | ImDrawFlags_RoundCornersTopRight | ImDrawFlags_RoundCornersBottomLeft | ImDrawFlags_RoundCornersBottomRight
	ImDrawFlags_RoundCornersDefault_    ImDrawFlags = ImDrawFlags_RoundCornersAll // Default to ALL corners if none of the _RoundCornersXX flags are specified.
	ImDrawFlags_RoundCornersMask_       ImDrawFlags = ImDrawFlags_RoundCornersAll | ImDrawFlags_RoundCornersNone
)

// Flags for ImDrawList instance. Those are set automatically by ImGui:: functions from ImGuiIO settings, and generally not manipulated directly.
// It is however possible to temporarily alter flags between calls to ImDrawList:: functions.
const (
	ImDrawListFlags_None                   ImDrawListFlags = 0
	ImDrawListFlags_AntiAliasedLines       ImDrawListFlags = 1 << 0 // Enable anti-aliased lines/borders (*2 the number of triangles for 1.0f wide line or lines thin enough to be drawn using textures, otherwise *3 the number of triangles)
	ImDrawListFlags_AntiAliasedLinesUseTex ImDrawListFlags = 1 << 1 // Enable anti-aliased lines/borders using textures when possible. Require backend to render with bilinear filtering.
	ImDrawListFlags_AntiAliasedFill        ImDrawListFlags = 1 << 2 // Enable anti-aliased edge around filled shapes (rounded rectangles, circles).
	ImDrawListFlags_AllowVtxOffset         ImDrawListFlags = 1 << 3 // Can emit 'VtxOffset > 0' to allow large meshes. Set when 'ImGuiBackendFlags_RendererHasVtxOffset' is enabled.
)

// Flags for ImFontAtlas build
type ImFontAtlasFlags_ int

const (
	ImFontAtlasFlags_None               ImFontAtlasFlags_ = 0
	ImFontAtlasFlags_NoPowerOfTwoHeight ImFontAtlasFlags_ = 1 << 0 // Don't round the height to next power of two
	ImFontAtlasFlags_NoMouseCursors     ImFontAtlasFlags_ = 1 << 1 // Don't build software mouse cursors into the atlas (save a little texture memory)
	ImFontAtlasFlags_NoBakedLines       ImFontAtlasFlags_ = 1 << 2 // Don't build thick line textures into the atlas (save a little texture memory). The AntiAliasedLinesUseTex features uses them, otherwise they will be rendered using polygons (more expensive for CPU/GPU).
)

// Flags stored in ImGuiViewport::Flags, giving indications to the platform backends.
const (
	ImGuiViewportFlags_None              ImGuiViewportFlags = 0
	ImGuiViewportFlags_IsPlatformWindow  ImGuiViewportFlags = 1 << 0 // Represent a Platform Window
	ImGuiViewportFlags_IsPlatformMonitor ImGuiViewportFlags = 1 << 1 // Represent a Platform Monitor (unused yet)
	ImGuiViewportFlags_OwnedByApp        ImGuiViewportFlags = 1 << 2 // Platform Window: is created/managed by the application (rather than a dear imgui backend)
)
