package imgui

type BitVector struct{}            // Store 1-bit per value
type Rect struct{}                 // An axis-aligned rectangle (2 points)
type DrawDataBuilder struct{}      // Helper to build a ImDrawData instance
type ColorMod struct{}             // Stacked color modifier, backup of modified data so we can restore it
type ContextHook struct{}          // Hook for extensions like ImGuiTestEngine
type DataTypeInfo struct{}         // Type information associated to a ImGuiDataType enum
type GroupData struct{}            // Stacked storage data for BeginGroup()/EndGroup()
type InputTextState struct{}       // Internal state of the currently focused/edited text input box
type LastItemData struct{}         // Status storage for last submitted items
type MenuColumns struct{}          // Simple column measurement, currently used for MenuItem() only
type NavItemData struct{}          // Result of a gamepad/keyboard directional navigation move query result
type MetricsConfig struct{}        // Storage for ShowMetricsWindow() and DebugNodeXXX() functions
type NextWindowData struct{}       // Storage for SetNextWindow** functions
type NextItemData struct{}         // Storage for SetNextItem** functions
type OldColumnData struct{}        // Storage data for a single column for legacy Columns() api
type OldColumns struct{}           // Storage data for a columns set for legacy Columns() api
type PopupData struct{}            // Storage for current popup stack
type SettingsHandler struct{}      // Storage for one type registered in the .ini file
type StackSizes struct{}           // Storage of stack sizes for debugging/asserting
type StyleMod struct{}             // Stacked style modifier, backup of modified data so we can restore it
type TabBar struct{}               // Storage for a tab bar
type TabItem struct{}              // Storage for a tab item (within a tab bar)
type Table struct{}                // Storage for a table
type TableColumn struct{}          // Storage for one column of a table
type TableTempData struct{}        // Temporary storage for one table (one per table in the stack), shared between tables.
type TableSettings struct{}        // Storage for a table .ini settings
type TableColumnsSettings struct{} // Storage for a column .ini settings
type Window struct{}               // Storage for one window
type WindowTempData struct{}       // Temporary storage for one window (that's the data which in theory we could ditch at the end of the frame, in practice we currently keep it for each window)
type WindowSettings struct{}       // Storage for a window .ini settings (we keep one of those even if the actual window wasn't instanced during this session)

// Use your programming IDE "Go to definition" facility on the names of the center columns to find the actual flags/enum lists.
type LayoutType int          // -> enum ImGuiLayoutType_         // Enum: Horizontal or vertical
type ItemFlags int           // -> enum ImGuiItemFlags_          // Flags: for PushItemFlag()
type ItemStatusFlags int     // -> enum ImGuiItemStatusFlags_    // Flags: for DC.LastItemStatusFlags
type OldColumnFlags int      // -> enum ImGuiOldColumnFlags_     // Flags: for BeginColumns()
type NavHighlightFlags int   // -> enum ImGuiNavHighlightFlags_  // Flags: for RenderNavHighlight()
type NavDirSourceFlags int   // -> enum ImGuiNavDirSourceFlags_  // Flags: for GetNavInputAmount2d()
type NavMoveFlags int        // -> enum ImGuiNavMoveFlags_       // Flags: for navigation requests
type NextItemDataFlags int   // -> enum ImGuiNextItemDataFlags_  // Flags: for SetNextItemXXX() functions
type NextWindowDataFlags int // -> enum ImGuiNextWindowDataFlags_// Flags: for SetNextWindowXXX() functions
type SeparatorFlags int      // -> enum ImGuiSeparatorFlags_     // Flags: for SeparatorEx()
type TextFlags int           // -> enum ImGuiTextFlags_          // Flags: for TextEx()
type TooltipFlags int        // -> enum ImGuiTooltipFlags_       // Flags: for BeginTooltipEx()
