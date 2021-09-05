package imgui

type ImGuiInputTextCallback func(data *InputTextCallbackData) int
type ImGuiSizeCallback func(data *SizeCallbackData)

// Enums/Flags (declared as int for compatibility with old C++, to allow using as flags without overhead, and to not pollute the top of this file)
// - Tip: Use your programming IDE navigation facilities on the names in the _central column_ below to find the actual flags/enum lists!
//   In Visual Studio IDE: CTRL+comma ("Edit.NavigateTo") can follow symbols in comments, whereas CTRL+F12 ("Edit.GoToImplementation") cannot.
//   With Visual Assist installed: ALT+G ("VAssistX.GoToImplementation") can also follow symbols in comments.
type Col int              // -> enum ImGuiCol_             // Enum: A color identifier for styling
type Cond int             // -> enum ImGuiCond_            // Enum: A condition for many Set*() functions
type DataType int         // -> enum ImGuiDataType_        // Enum: A primary data type
type Dir int              // -> enum ImGuiDir_             // Enum: A cardinal direction
type Key int              // -> enum ImGuiKey_             // Enum: A key identifier (ImGui-side enum)
type NavInput int         // -> enum ImGuiNavInput_        // Enum: An input identifier for navigation
type MouseButton int      // -> enum ImGuiMouseButton_     // Enum: A mouse button identifier (0=left, 1=right, 2=middle)
type MouseCursor int      // -> enum ImGuiMouseCursor_     // Enum: A mouse cursor identifier
type SortDirection int    // -> enum ImGuiSortDirection_   // Enum: A sorting direction (ascending or descending)
type StyleVar int         // -> enum ImGuiStyleVar_        // Enum: A variable identifier for styling
type TableBgTarget int    // -> enum ImGuiTableBgTarget_   // Enum: A color target for TableSetBgColor()
type DrawFlags int        // -> enum ImDrawFlags_          // Flags: for ImDrawList functions
type DrawListFlags int    // -> enum ImDrawListFlags_      // Flags: for ImDrawList instance
type FontAtlasFlags int   // -> enum ImFontAtlasFlags_     // Flags: for ImFontAtlas build
type BackendFlags int     // -> enum ImGuiBackendFlags_    // Flags: for io.BackendFlags
type ButtonFlags int      // -> enum ImGuiButtonFlags_     // Flags: for InvisibleButton()
type ColorEditFlags int   // -> enum ImGuiColorEditFlags_  // Flags: for ColorEdit4(), ColorPicker4() etc.
type ConfigFlags int      // -> enum ImGuiConfigFlags_     // Flags: for io.ConfigFlags
type ComboFlags int       // -> enum ImGuiComboFlags_      // Flags: for BeginCombo()
type DragDropFlags int    // -> enum ImGuiDragDropFlags_   // Flags: for BeginDragDropSource(), AcceptDragDropPayload()
type FocusedFlags int     // -> enum ImGuiFocusedFlags_    // Flags: for IsWindowFocused()
type HoveredFlags int     // -> enum ImGuiHoveredFlags_    // Flags: for IsItemHovered(), IsWindowHovered() etc.
type InputTextFlags int   // -> enum ImGuiInputTextFlags_  // Flags: for InputText(), InputTextMultiline()
type KeyModFlags int      // -> enum ImGuiKeyModFlags_     // Flags: for io.KeyMods (Ctrl/Shift/Alt/Super)
type PopupFlags int       // -> enum ImGuiPopupFlags_      // Flags: for OpenPopup*(), BeginPopupContext*(), IsPopupOpen()
type SelectableFlags int  // -> enum ImGuiSelectableFlags_ // Flags: for Selectable()
type SliderFlags int      // -> enum ImGuiSliderFlags_     // Flags: for DragFloat(), DragInt(), SliderFloat(), SliderInt() etc.
type TabBarFlags int      // -> enum ImGuiTabBarFlags_     // Flags: for BeginTabBar()
type TabItemFlags int     // -> enum ImGuiTabItemFlags_    // Flags: for BeginTabItem()
type TableFlags int       // -> enum ImGuiTableFlags_      // Flags: For BeginTable()
type TableColumnFlags int // -> enum ImGuiTableColumnFlags_// Flags: For TableSetupColumn()
type TableRowFlags int    // -> enum ImGuiTableRowFlags_   // Flags: For TableNextRow()
type TreeNodeFlags int    // -> enum ImGuiTreeNodeFlags_   // Flags: for TreeNode(), TreeNodeEx(), CollapsingHeader()
type ViewportFlags int    // -> enum ImGuiViewportFlags_   // Flags: for ImGuiViewport
type WindowFlags int      // -> enum ImGuiWindowFlags_     // Flags: for Begin(), BeginChild()

type DrawChannel struct{}        // Temporary storage to output draw commands out of order, used by ImDrawListSplitter and ImDrawList::ChannelsSplit()
type DrawListSharedData struct{} // Data shared among multiple draw lists (typically owned by parent ImGui context, but you may create one yourself)
type DrawListSplitter struct{}   // Helper to split a draw list into different layers which can be drawn into out of order, then flattened back.
type FontBuilderIO struct{}      // Opaque interface to a font builder (stb_truetype or FreeType).

type FontGlyph struct{}              // A single font glyph (code point + coordinates within in ImFontAtlas + offset)
type FontGlyphRangesBuilder struct{} // Helper to build glyph ranges from text/string data
type InputTextCallbackData struct{}  // Shared state of InputText() when using custom ImGuiInputTextCallback (rare/advanced use)
type ListClipper struct{}            // Helper to manually clip large list of items
type OnceUponAFrame struct{}         // Helper for running a block of code not more than once a frame, used by IMGUI_ONCE_UPON_A_FRAME macro
type Payload struct{}                // User data payload for drag and drop operations
type SizeCallbackData struct{}       // Callback data when using SetNextWindowSizeConstraints() (rare/advanced use)
type Storage struct{}                // Helper for key->value storage
type TableSortSpecs struct{}         // Sorting specifications for a table (often handling sort specs for a single column, occasionally more)
type TableColumnSortSpecs struct{}   // Sorting specification for one column of a table
type TextBuffer struct{}             // Helper to hold and append into a text buffer (~string builder)
type TextFilter struct{}             // Helper to parse and apply text filters (e.g. "aaaaa[,bbbbb][,ccccc]")
type Viewport struct{}               // A Platform Window (always only one in 'master' branch), in the future may represent Platform Monitor
