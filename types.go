package imgui

import (
	"math"
	"unicode"
	"unsafe"
)

const IMGUI_VERSION = "1.85 WIP"
const IMGUI_VERSION_NUM = 18411

func IMGUI_CHECKVERSION() {
	panic("imgui: IMGUI_CHECKVERSION not supported")
}

const IMGUI_HAS_TABLE = true

// Enums/Flags (declared as for int compatibility with old C++, to allow using as flags without overhead, and to not pollute the top of this file)
// - Tip: Use your programming IDE navigation facilities on the names in the _central column_ below to find the actual flags/enum lists!
//   In Visual Studio IDE: CTRL+comma ("Edit.NavigateTo") can follow symbols in comments, whereas CTRL+F12 ("Edit.GoToImplementation") cannot.
//   With Visual Assist installed: ALT+G ("VAssistX.GoToImplementation") can also follow symbols in comments.
type ImGuiCol int              // -> enum ImGuiCol_             // Enum: A color identifier for styling
type ImGuiCond int             // -> enum ImGuiCond_            // Enum: A condition for many Set*() functions
type ImGuiDataType int         // -> enum ImGuiDataType_        // Enum: A primary data type
type ImGuiDir int              // -> enum ImGuiDir_             // Enum: A cardinal direction
type ImGuiKey int              // -> enum ImGuiKey_             // Enum: A key identifier (ImGui-side enum)
type ImGuiNavInput int         // -> enum ImGuiNavInput_        // Enum: An input identifier for navigation
type ImGuiMouseButton int      // -> enum ImGuiMouseButton_     // Enum: A mouse button identifier (0=left, 1=right, 2=middle)
type ImGuiMouseCursor int      // -> enum ImGuiMouseCursor_     // Enum: A mouse cursor identifier
type ImGuiSortDirection int    // -> enum ImGuiSortDirection_   // Enum: A sorting direction (ascending or descending)
type ImGuiStyleVar int         // -> enum ImGuiStyleVar_        // Enum: A variable identifier for styling
type ImGuiTableBgTarget int    // -> enum ImGuiTableBgTarget_   // Enum: A color target for TableSetBgColor()
type ImDrawFlags int           // -> enum ImDrawFlags_          // Flags: for ImDrawList functions
type ImDrawListFlags int       // -> enum ImDrawListFlags_      // Flags: for ImDrawList instance
type ImFontAtlasFlags int      // -> enum ImFontAtlasFlags_     // Flags: for ImFontAtlas build
type ImGuiBackendFlags int     // -> enum ImGuiBackendFlags_    // Flags: for io.BackendFlags
type ImGuiButtonFlags int      // -> enum ImGuiButtonFlags_     // Flags: for InvisibleButton()
type ImGuiColorEditFlags int   // -> enum ImGuiColorEditFlags_  // Flags: for ColorEdit4(), ColorPicker4() etc.
type ImGuiConfigFlags int      // -> enum ImGuiConfigFlags_     // Flags: for io.ConfigFlags
type ImGuiComboFlags int       // -> enum ImGuiComboFlags_      // Flags: for BeginCombo()
type ImGuiDragDropFlags int    // -> enum ImGuiDragDropFlags_   // Flags: for BeginDragDropSource(), AcceptDragDropPayload()
type ImGuiFocusedFlags int     // -> enum ImGuiFocusedFlags_    // Flags: for IsWindowFocused()
type ImGuiHoveredFlags int     // -> enum ImGuiHoveredFlags_    // Flags: for IsItemHovered(), IsWindowHovered() etc.
type ImGuiInputTextFlags int   // -> enum ImGuiInputTextFlags_  // Flags: for InputText(), InputTextMultiline()
type ImGuiKeyModFlags int      // -> enum ImGuiKeyModFlags_     // Flags: for io.KeyMods (Ctrl/Shift/Alt/Super)
type ImGuiPopupFlags int       // -> enum ImGuiPopupFlags_      // Flags: for OpenPopup*(), BeginPopupContext*(), IsPopupOpen()
type ImGuiSelectableFlags int  // -> enum ImGuiSelectableFlags_ // Flags: for Selectable()
type ImGuiSliderFlags int      // -> enum ImGuiSliderFlags_     // Flags: for DragFloat(), DragInt(), SliderFloat(), SliderInt() etc.
type ImGuiTabBarFlags int      // -> enum ImGuiTabBarFlags_     // Flags: for BeginTabBar()
type ImGuiTabItemFlags int     // -> enum ImGuiTabItemFlags_    // Flags: for BeginTabItem()
type ImGuiTableFlags int       // -> enum ImGuiTableFlags_      // Flags: For BeginTable()
type ImGuiTableColumnFlags int // -> enum ImGuiTableColumnFlags_// Flags: For TableSetupColumn()
type ImGuiTableRowFlags int    // -> enum ImGuiTableRowFlags_   // Flags: For TableNextRow()
type ImGuiTreeNodeFlags int    // -> enum ImGuiTreeNodeFlags_   // Flags: for TreeNode(), TreeNodeEx(), CollapsingHeader()
type ImGuiViewportFlags int    // -> enum ImGuiViewportFlags_   // Flags: for ImGuiViewport
type ImGuiWindowFlags int      // -> enum ImGuiWindowFlags_     // Flags: for Begin(), BeginChild()

type ImTextureID uintptr // Default: store a pointer or an integer fitting in a pointer (most renderer backends are ok with that)

type ImDrawIdx uint16 // Default: 16-bit (for maximum compatibility with renderer backends)

// Scalar data types
type ImGuiID = uint // A unique ID used by widgets (typically the result of hashing a stack of string)
type ImS8 = int8    // 8-bit signed integer
type ImU8 = uint8   // 8-bit uinteger
type ImS16 = int16  // 16-bit signed integer
type ImU16 = uint16 // 16-bit uinteger
type ImS32 = int    // 32-bit signed integer == int
type ImU32 = uint   // 32-bit uinteger (often used to store packed colors)
type ImS64 = int64  // 64-bit signed integer (pre and post C++11 with Visual Studio)
type ImU64 = uint64 // 64-bit uinteger (pre and post C++11 with Visual Studio)

// Character types
// (we generally use UTF-8 encoded string in the API. This is storage specifically for a decoded character used for keyboard input and display)
type ImWchar16 = uint16 // A single decoded U16 character/code point. We encode them as multi bytes UTF-8 when used in strings.
type ImWchar32 = int    // A single decoded U32 character/code point. We encode them as multi bytes UTF-8 when used in strings.
type ImWchar = rune     // ImWchar [configurable type: override in imconfig.h with '#define IMGUI_USE_WCHAR32' to support Unicode planes 1-16]
const MaxImWchar ImWchar = math.MaxInt32

// Callback and functions types
type ImGuiInputTextCallback func(data *ImGuiInputTextCallbackData) int // Callback function for ImGui::InputText()
type ImGuiSizeCallback func(data *ImGuiSizeCallbackData)               // Callback function for ImGui::SetNextWindowSizeConstraints()

//-----------------------------------------------------------------------------
// [SECTION] ImGuiIO
//-----------------------------------------------------------------------------
// Communicate most settings and inputs/outputs to Dear ImGui using this structure.
// Access via ImGui::GetIO(). Read 'Programmer guide' section in .cpp file for general usage.
//-----------------------------------------------------------------------------

type ImGuiIO struct {
	//------------------------------------------------------------------
	// Configuration (fill once)                // Default value
	//------------------------------------------------------------------

	ConfigFlags             ImGuiConfigFlags    // = 0              // See ImGuiConfigFlags_ enum. Set by user/application. Gamepad/keyboard navigation options, etc.
	BackendFlags            ImGuiBackendFlags   // = 0              // See ImGuiBackendFlags_ enum. Set by backend (imgui_impl_xxx files or custom backend) to communicate features supported by the backend.
	DisplaySize             ImVec2              // <unset>          // Main display size, in pixels (generally == GetMainViewport()->Size)
	DeltaTime               float               // /*= 1.0f*//60.0f     // Time elapsed since last frame, in seconds.
	IniSavingRate           float               // = 5.0f           // Minimum time between saving positions/sizes to .ini file, in seconds.
	IniFilename             string              // = "imgui.ini"    // Path to .ini file (important: default "imgui.ini" is relative to current working dir!). Set NULL to disable automatic .ini loading/saving or if you want to manually call LoadIniSettingsXXX() / SaveIniSettingsXXX() functions.
	LogFilename             string              // = "imgui_log.txt"// Path to .log file (default parameter to ImGui::LogToFile when no file is specified).
	MouseDoubleClickTime    float               // = 0.30f          // Time for a double-click, in seconds.
	MouseDoubleClickMaxDist float               // = 6.0f           // Distance threshold to stay in to validate a double-click, in pixels.
	MouseDragThreshold      float               // = 6.0f           // Distance threshold before considering we are dragging.
	KeyMap                  [ImGuiKey_COUNT]int // <unset>          // Map of indices into the KeysDown[512] entries array which represent your "native" keyboard state.
	KeyRepeatDelay          float               // = 0.250f         // When holding a key/button, time before it starts repeating, in seconds (for buttons in Repeat mode, etc.).
	KeyRepeatRate           float               // = 0.050f         // When holding a key/button, rate at which it repeats, in seconds.
	UserData                interface{}         // = NULL           // Store your own data for retrieval by callbacks.

	Fonts                   *ImFontAtlas // <auto>           // Font atlas: load, rasterize and pack one or more fonts into a single texture.
	FontGlobalScale         float        // /*= 1.0f*/           // Global scale all fonts
	FontAllowUserScaling    bool         // = false          // Allow user scaling text of individual window with CTRL+Wheel.
	FontDefault             *ImFont      // = NULL           // Font to use on NewFrame(). Use NULL to uses Fonts->Fonts[0].
	DisplayFramebufferScale ImVec2       // = (1, 1)         // For retina display or other situations where window coordinates are different from framebuffer coordinates. This generally ends up in ImDrawData::FramebufferScale.

	// Miscellaneous options
	MouseDrawCursor                   bool  // = false          // Request ImGui to draw a mouse cursor for you (if you are on a platform without a mouse cursor). Cannot be easily renamed to 'io.ConfigXXX' because this is frequently used by backend implementations.
	ConfigMacOSXBehaviors             bool  // = defined(__APPLE__) // OS X style: Text editing cursor movement using Alt instead of Ctrl, Shortcuts using Cmd/Super instead of Ctrl, Line/Text Start and End using Cmd+Arrows instead of Home/End, Double click selects by word instead of selecting whole text, Multi-selection in lists uses Cmd/Super instead of Ctrl.
	ConfigInputTextCursorBlink        bool  // = true           // Enable blinking cursor (optional as some users consider it to be distracting).
	ConfigDragClickToInputText        bool  // = false          // [BETA] Enable turning DragXXX widgets into text input with a simple mouse click-release (without moving). Not desirable on devices without a keyboard.
	ConfigWindowsResizeFromEdges      bool  // = true           // Enable resizing of windows from their edges and from the lower-left corner. This requires (io.BackendFlags & ImGuiBackendFlags_HasMouseCursors) because it needs mouse cursor feedback. (This used to be a per-window ImGuiWindowFlags_ResizeFromAnySide flag)
	ConfigWindowsMoveFromTitleBarOnly bool  // = false       // Enable allowing to move windows only when clicking on their title bar. Does not apply to windows without a title bar.
	ConfigMemoryCompactTimer          float // = 60.0f          // Timer (in seconds) to free transient windows/tables memory buffers when unused. Set to -1.0f to disable.

	//------------------------------------------------------------------
	// Platform Functions
	// (the imgui_impl_xxxx backend files are setting those up for you)
	//------------------------------------------------------------------

	// Optional: Platform/Renderer backend name (informational only! will be displayed in About Window) + User data for backend/wrappers to store their own stuff.
	BackendPlatformName     string      // = NULL
	BackendRendererName     string      // = NULL
	BackendPlatformUserData interface{} // = NULL           // User data for platform backend
	BackendRendererUserData interface{} // = NULL           // User data for renderer backend
	BackendLanguageUserData interface{} // = NULL           // User data for non C++ programming language backend

	// Optional: Access OS clipboard
	// (default to use native Win32 clipboard on Windows, otherwise uses a private clipboard. Override to access OS clipboard on other architectures)
	GetClipboardTextFn func(user_data interface{}) string
	SetClipboardTextFn func(user_data interface{}, text string) string
	ClipboardUserData  interface{}

	// Optional: Notify OS Input Method Editor of the screen position of your cursor for text input position (e.g. when using Japanese/Chinese IME on Windows)
	// (default to use native imm32 api on Windows)
	ImeSetInputScreenPosFn func(x, y int) int
	ImeWindowHandle        interface{} // = NULL           // (Windows) Set this to your HWND to get automatic IME cursor positioning.

	//------------------------------------------------------------------
	// Input - Fill before calling NewFrame()
	//------------------------------------------------------------------

	MousePos    ImVec2                     // Mouse position, in pixels. Set to ImVec2(-FLT_MAX, -FLT_MAX) if mouse is unavailable (on another screen, etc.)
	MouseDown   [5]bool                    // Mouse buttons: 0=left, 1=right, 2=middle + extras (ImGuiMouseButton_COUNT == 5). Dear ImGui mostly uses left and right buttons. Others buttons allows us to track if the mouse is being used by your application + available to user as a convenience via IsMouse** API.
	MouseWheel  float                      // Mouse wheel Vertical: 1 unit scrolls about 5 lines text.
	MouseWheelH float                      // Mouse wheel Horizontal. Most users don't have a mouse with an horizontal wheel, may not be filled by all backends.
	KeyCtrl     bool                       // Keyboard modifier pressed: Control
	KeyShift    bool                       // Keyboard modifier pressed: Shift
	KeyAlt      bool                       // Keyboard modifier pressed: Alt
	KeySuper    bool                       // Keyboard modifier pressed: Cmd/Super/Windows
	KeysDown    [512]bool                  // Keyboard keys that are pressed (ideally left in the "native" order your engine has access to keyboard keys, so you can use your own defines/enums for keys).
	NavInputs   [ImGuiNavInput_COUNT]float // Gamepad inputs. Cleared back to zero by EndFrame(). Keyboard keys will be auto-mapped and be written here by NewFrame().
	// Notifies Dear ImGui when hosting platform windows lose or gain input focus

	//------------------------------------------------------------------
	// Output - Updated by NewFrame() or EndFrame()/Render()
	// (when reading from the io.WantCaptureMouse, io.WantCaptureKeyboard flags to dispatch your inputs, it is
	//  generally easier and more correct to use their state BEFORE calling NewFrame(). See FAQ for details!)
	//------------------------------------------------------------------

	WantCaptureMouse         bool   // Set when Dear ImGui will use mouse inputs, in this case do not dispatch them to your main game/application (either way, always pass on mouse inputs to imgui). (e.g. unclicked mouse is hovering over an imgui window, widget is active, mouse was clicked over an imgui window, etc.).
	WantCaptureKeyboard      bool   // Set when Dear ImGui will use keyboard inputs, in this case do not dispatch them to your main game/application (either way, always pass keyboard inputs to imgui). (e.g. InputText active, or an imgui window is focused and navigation is enabled, etc.).
	WantTextInput            bool   // Mobile/console: when set, you may display an on-screen keyboard. This is set by Dear ImGui when it wants textual keyboard input to happen (e.g. when a InputText widget is active).
	WantSetMousePos          bool   // MousePos has been altered, backend should reposition mouse on next frame. Rarely used! Set only when ImGuiConfigFlags_NavEnableSetMousePos flag is enabled.
	WantSaveIniSettings      bool   // When manual .ini load/save is active (io.IniFilename == NULL), this will be set to notify your application that you can call SaveIniSettingsToMemory() and save yourself. Important: clear io.WantSaveIniSettings yourself after saving!
	NavActive                bool   // Keyboard/Gamepad navigation is currently allowed (will handle ImGuiKey_NavXXX events) = a window is focused and it doesn't use the ImGuiWindowFlags_NoNavInputs flag.
	NavVisible               bool   // Keyboard/Gamepad navigation is visible and allowed (will handle ImGuiKey_NavXXX events).
	Framerate                float  // Rough estimate of application framerate, in frame per second. Solely for convenience. Rolling average estimation based on io.DeltaTime over 120 frames.
	MetricsRenderVertices    int    // Vertices output during last call to Render()
	MetricsRenderIndices     int    // Indices output during last call to Render() = number of triangles * 3
	MetricsRenderWindows     int    // Number of visible windows
	MetricsActiveWindows     int    // Number of active windows
	MetricsActiveAllocations int    // Number of active allocations, updated by MemAlloc/MemFree based on current context. May be off if you have multiple imgui contexts.
	MouseDelta               ImVec2 // Mouse delta. Note that this is zero if either current or previous position are invalid (-FLT_MAX,-FLT_MAX), so a disappearing/reappearing mouse won't have a huge delta.

	//------------------------------------------------------------------
	// [Internal] Dear ImGui will maintain those fields. Forward compatibility not guaranteed!
	//------------------------------------------------------------------

	WantCaptureMouseUnlessPopupClose bool             // Alternative to WantCaptureMouse: (WantCaptureMouse == true && WantCaptureMouseUnlessPopupClose == false) when a click over void is expected to close a popup.
	KeyMods                          ImGuiKeyModFlags // Key mods flags (same as io.KeyCtrl/KeyShift/KeyAlt/KeySuper but merged into flags), updated by NewFrame()
	KeyModsPrev                      ImGuiKeyModFlags // Previous key mods
	MousePosPrev                     ImVec2           // Previous mouse position (note that MouseDelta is not necessary == MousePos-MousePosPrev, in case either position is invalid)
	MouseClickedPos                  [5]ImVec2        // Position at time of clicking
	MouseClickedTime                 [5]double        // Time of last click (used to figure out double-click)
	MouseClicked                     [5]bool          // Mouse button went from !Down to Down
	MouseDoubleClicked               [5]bool          // Has mouse button been double-clicked?
	MouseReleased                    [5]bool          // Mouse button went from Down to !Down
	MouseDownOwned                   [5]bool          // Track if button was clicked inside a dear imgui window or over void blocked by a popup. We don't request mouse capture from the application if click started outside ImGui bounds.
	MouseDownOwnedUnlessPopupClose   [5]bool          //Track if button was clicked inside a dear imgui window.
	MouseDownWasDoubleClick          [5]bool          // Track if button down was a double-click
	MouseDownDuration                [5]float         // Duration the mouse button has been down (0.0f == just clicked)
	MouseDownDurationPrev            [5]float         // Previous time the mouse button has been down
	MouseDragMaxDistanceAbs          [5]ImVec2        // Maximum distance, absolute, on each axis, of how much mouse has traveled from the clicking point
	MouseDragMaxDistanceSqr          [5]float         // Squared maximum distance of how much mouse has traveled from the clicking point
	KeysDownDuration                 [512]float       // Duration the keyboard key has been down (0.0f == just pressed)
	KeysDownDurationPrev             [512]float       // Previous duration the key has been down
	NavInputsDownDuration            [ImGuiNavInput_COUNT]float
	NavInputsDownDurationPrev        [ImGuiNavInput_COUNT]float
	PenPressure                      float     // Touch/Pen pressure (0.0f to 1.0f, should be >0.0f only when MouseDown[0] == true). Helper storage currently unused by Dear ImGui.
	InputQueueSurrogate              ImWchar16 // For AddInputCharacterUTF16
	InputQueueCharacters             []ImWchar // Queue of _characters_ input (obtained by platform backend). Fill using AddInputCharacter() helper.
}

func NewImGuiIO() ImGuiIO {
	var io ImGuiIO
	IM_ASSERT(int(len(io.MouseDown)) == int(ImGuiMouseButton_COUNT) && int(len(io.MouseClicked)) == int(ImGuiMouseButton_COUNT)) // Our pre-C++11 IM_STATIC_ASSERT() macros triggers warning on modern compilers so we don't use it here.

	// Settings
	io.ConfigFlags = ImGuiConfigFlags_None
	io.BackendFlags = ImGuiBackendFlags_None
	io.DisplaySize = ImVec2{-1.0, -1.0}
	io.DeltaTime = 1.0 / 60.0
	io.IniSavingRate = 5.0
	io.IniFilename = "imgui.ini" // Important: "imgui.ini" is relative to current working dir, most apps will want to lock this to an absolute path (e.g. same path as executables).
	io.LogFilename = "imgui_log.txt"
	io.MouseDoubleClickTime = 0.30
	io.MouseDoubleClickMaxDist = 6.0
	for i := ImGuiKey(0); i < ImGuiKey_COUNT; i++ {
		io.KeyMap[i] = -1
	}
	io.KeyRepeatDelay = 0.275
	io.KeyRepeatRate = 0.050

	io.FontGlobalScale = 1.0
	io.FontAllowUserScaling = false
	io.DisplayFramebufferScale = ImVec2{1.0, 1.0}

	// Miscellaneous options
	io.MouseDrawCursor = false
	/*#ifdef __APPLE__
	      ConfigMacOSXBehaviors = true;  // Set Mac OS X style defaults based on __APPLE__ compile time flag
	  #else
	      ConfigMacOSXBehaviors = false;
	  #endif*/
	io.ConfigInputTextCursorBlink = true
	io.ConfigWindowsResizeFromEdges = true
	io.ConfigWindowsMoveFromTitleBarOnly = false
	io.ConfigMemoryCompactTimer = 60.0

	// Platform Functions
	//io.GetClipboardTextFn = GetClipboardTextFn_DefaultImpl // Platform dependent default implementations
	//io.SetClipboardTextFn = SetClipboardTextFn_DefaultImpl
	//io.ImeSetInputScreenPosFn = ImeSetInputScreenPosFn_DefaultImpl

	// Input (NB: we already have memset zero the entire structure!)
	io.MousePos = ImVec2{-FLT_MAX, -FLT_MAX}
	io.MousePosPrev = ImVec2{-FLT_MAX, -FLT_MAX}
	io.MouseDragThreshold = 6.0
	for i := 0; i < len(io.MouseDownDuration); i++ {
		io.MouseDownDuration[i] = -1.0
		io.MouseDownDurationPrev[i] = -1.0
	}
	for i := 0; i < len(io.KeysDownDuration); i++ {
		io.KeysDownDuration[i] = -1
		io.KeysDownDurationPrev[i] = -1.0
	}
	for i := 0; i < len(io.NavInputsDownDuration); i++ {
		io.NavInputsDownDuration[i] = -1.0
	}
	return io
}

// Shared state of InputText(), passed as an argument to your callback when a ImGuiInputTextFlags_Callback* flag is used.
// The callback function should return 0 by default.
// Callbacks (follow a flag name and see comments in ImGuiInputTextFlags_ declarations for more details)
// - ImGuiInputTextFlags_CallbackEdit:        Callback on buffer edit (note that InputText() already returns true on edit, the callback is useful mainly to manipulate the underlying buffer while focus is active)
// - ImGuiInputTextFlags_CallbackAlways:      Callback on each iteration
// - ImGuiInputTextFlags_CallbackCompletion:  Callback on pressing TAB
// - ImGuiInputTextFlags_CallbackHistory:     Callback on pressing Up/Down arrows
// - ImGuiInputTextFlags_CallbackCharFilter:  Callback on character inputs to replace or discard them. Modify 'EventChar' to replace or discard, or return 1 in callback to discard.
// - ImGuiInputTextFlags_CallbackResize:      Callback on buffer capacity changes request (beyond 'buf_size' parameter value), allowing the string to grow.
type ImGuiInputTextCallbackData struct {
	EventFlag ImGuiInputTextFlags // One ImGuiInputTextFlags_Callback*    // Read-only
	Flags     ImGuiInputTextFlags // What user passed to InputText()      // Read-only
	UserData  interface{}         // What user passed to InputText()      // Read-only

	// Arguments for the different callback events
	// - To modify the text buffer in a callback, prefer using the InsertChars() / DeleteChars() function. InsertChars() will take care of calling the resize callback if necessary.
	// - If you know your edits are not going to resize the underlying buffer allocation, you may modify the contents of 'Buf[]' directly. You need to update 'BufTextLen' accordingly (0 <= BufTextLen < BufSize) and set 'BufDirty'' to true so InputText can update its internal state.
	EventChar      ImWchar  // Character input                      // Read-write   // [CharFilter] Replace character with another one, or set to zero to drop. return 1 is equivalent to EventChar=0 setting
	EventKey       ImGuiKey // Key pressed (Up/Down/TAB)            // Read-only    // [Completion,History]
	Buf            []byte   // Text buffer                          // Read-write   // [Resize] Can replace pointer / [Completion,History,Always] Only write to pointed data, don't replace the actual pointer!
	BufTextLen     int      // Text length (in bytes)               // Read-write   // [Resize,Completion,History,Always] Exclude zero-terminator storage. In C land: == strlen(some_text), in C++ land: string.length()
	BufSize        int      // Buffer size (in bytes) = capacity+1  // Read-only    // [Resize,Completion,History,Always] Include zero-terminator storage. In C land == ARRAYSIZE(my_char_array), in C++ land: string.capacity()+1
	BufDirty       bool     // Set if you modify Buf/BufTextLen!    // Write        // [Completion,History,Always]
	CursorPos      int      //                                      // Read-write   // [Completion,History,Always]
	SelectionStart int      //                                      // Read-write   // [Completion,History,Always] == to SelectionEnd when no selection)
	SelectionEnd   int      //                                      // Read-write   // [Completion,History,Always]
}

func NewImGuiInputTextCallbackData() *ImGuiInputTextCallbackData     { panic("not implemented") }
func (*ImGuiInputTextCallbackData) DeleteChars(pos, bytes_count int) { panic("not implemented") }
func (*ImGuiInputTextCallbackData) InsertChars(pos int, text, text_end string) {
	panic("not implemented")
}
func (this *ImGuiInputTextCallbackData) SelectAll() {
	this.SelectionStart = 0
	this.SelectionEnd = this.BufTextLen
}
func (this *ImGuiInputTextCallbackData) ClearSelection() {
	this.SelectionStart = this.BufTextLen
	this.SelectionEnd = this.BufTextLen
}
func (this *ImGuiInputTextCallbackData) HasSelection() bool {
	return this.SelectionStart != this.SelectionEnd
}

// Resizing callback data to apply custom constraint. As enabled by SetNextWindowSizeConstraints(). Callback is called during the next Begin().
// NB: For basic min/max size constraint on each axis you don't need to use the callback! The SetNextWindowSizeConstraints() parameters are enough.
type ImGuiSizeCallbackData struct {
	UserData    interface{} // Read-only.   What user passed to SetNextWindowSizeConstraints()
	Pos         ImVec2      // Read-only.   Window position, for reference.
	CurrentSize ImVec2      // Read-only.   Current window size.
	DesiredSize ImVec2      // Read-write.  Desired size, based on user's mouse position. Write to this field to restrain resizing.
}

// Data payload for Drag and Drop operations: AcceptDragDropPayload(), GetDragDropPayload()
type ImGuiPayload struct {
	// Members
	Data     interface{} // Data (copied and owned by dear imgui)
	DataSize int         // Data size

	// [Internal]
	SourceId       ImGuiID      // Source item id
	SourceParentId ImGuiID      // Source parent id (if available)
	DataFrameCount int          // Data timestamp
	DataType       [32 + 1]byte // Data type tag (short user-supplied string, 32 characters max)
	Preview        bool         // Set when AcceptDragDropPayload() was called and mouse has been hovering the target item (nb: handle overlapping drag targets)
	Delivery       bool         // Set when AcceptDragDropPayload() was called and mouse button is released over the target item.
}

func NewImGuiPayload() ImGuiPayload {
	return ImGuiPayload{
		DataFrameCount: -1,
	}
}

func (this ImGuiPayload) IsDataType(dtype string) bool {
	return this.DataFrameCount != -1 && dtype == string(this.DataType[:])
}
func (this ImGuiPayload) IsPreview() bool  { return this.Preview }
func (this ImGuiPayload) IsDelivery() bool { return this.Delivery }

// Sorting specification for one column of a table (sizeof == 12 bytes)
type ImGuiTableColumnSortSpecs struct {
	ColumnUserID  ImGuiID            // User id of the column (if specified by a TableSetupColumn() call)
	ColumnIndex   ImS16              // Index of the column
	SortOrder     ImS16              // Index within parent ImGuiTableSortSpecs (always stored in order starting from 0, tables sorted on a single criteria will always have a 0 here)
	SortDirection ImGuiSortDirection // ImGuiSortDirection_Ascending or ImGuiSortDirection_Descending (you can use this or SortSign, whichever is more convenient for your sort function)
}

func NewImGuiTableColumnSortSpecs() ImGuiTableColumnSortSpecs {
	return ImGuiTableColumnSortSpecs{}
}

// Sorting specifications for a table (often handling sort specs for a single column, occasionally more)
// Obtained by calling TableGetSortSpecs().
// When 'SpecsDirty == true' you can sort your data. It will be true with sorting specs have changed since last call, or the first time.
// Make sure to set 'SpecsDirty = false' after sorting, else you may wastefully sort your data every frame!
type ImGuiTableSortSpecs struct {
	Specs      *ImGuiTableColumnSortSpecs // Pointer to sort spec array.
	SpecsCount int                        // Sort spec count. Most often 1. May be > 1 when ImGuiTableFlags_SortMulti is enabled. May be == 0 when ImGuiTableFlags_SortTristate is enabled.
	SpecsDirty bool                       // Set to true when specs have changed since last time! Use this to sort again, then clear the flag.
}

const IM_UNICODE_CODEPOINT_INVALID = unicode.ReplacementChar

/*
#ifdef IMGUI_USE_WCHAR32
#define IM_UNICODE_CODEPOINT_MAX     0x10FFFF   // Maximum Unicode code point supported by this build.
#else
*/
const IM_UNICODE_CODEPOINT_MAX = 0xFFFF

// Helper: Execute a block of code at maximum once a frame. Convenient if you want to quickly create an UI within deep-nested code that runs multiple times every frame.
// Usage: static oaf ImGuiOnceUponAFrame if (oaf) ImGui::Text("This will be called only once frame") per
type ImGuiOnceUponAFrame struct {
	RefFrame int
}

func NewImGuiOnceUponAFrame() ImGuiOnceUponAFrame {
	return ImGuiOnceUponAFrame{
		RefFrame: -1,
	}
}

func (this ImGuiOnceUponAFrame) Bool() bool {
	var current_frame int = GetFrameCount()
	if this.RefFrame == current_frame {
		return false
	}
	this.RefFrame = current_frame
	return true
}

// Helper: Manually clip large list of items.
// If you are submitting lots of evenly spaced items and you have a random access to the list, you can perform coarse
// clipping based on visibility to save yourself from processing those items at all.
// The clipper calculates the range of visible items and advance the cursor to compensate for the non-visible items we have skipped.
// (Dear ImGui already clip items based on their bounds but it needs to measure text size to do so, whereas manual coarse clipping before submission makes this cost and your own data fetching/submission cost almost null)
// Usage:
//   clipper ImGuiListClipper
//clipper.Begin(1000) //         // We have 1000 elements, evenly spaced.
//   while (clipper.Step())
//       for (int i clipper.DisplayStart = i clipper.DisplayEnd < i++)
//           ImGui::Text("line number i) %d",
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

func NewImGuiListClipper() ImGuiListClipper { panic("not implemented") }

// items_count: Use INT_MAX if you don't know how many items you have (in which case the cursor won't be advanced in the final step)
// items_height: Use -1.0f to be calculated automatically on first step. Otherwise pass in the distance between your items, typically GetTextLineHeightWithSpacing() or GetFrameHeightWithSpacing().
func (this ImGuiListClipper) Begin(items_count int, items_height float /*= -1.0f*/) {
	panic("not implemented")
}                                        // Automatically called by constructor if you passed 'items_count' or by Step() in Step 1.
func (this ImGuiListClipper) End()       { panic("not implemented") } // Automatically called on the last call of Step() that returns false.
func (this ImGuiListClipper) Step() bool { panic("not implemented") }

const IM_COL32_R_SHIFT = 0
const IM_COL32_G_SHIFT = 8
const IM_COL32_B_SHIFT = 16
const IM_COL32_A_SHIFT = 24
const IM_COL32_A_MASK = 0xFF000000

func IM_COL32(R, G, B, A byte) ImU32 {
	return (((ImU32)(A) << IM_COL32_A_SHIFT) | ((ImU32)(B) << IM_COL32_B_SHIFT) | ((ImU32)(G) << IM_COL32_G_SHIFT) | ((ImU32)(R) << IM_COL32_R_SHIFT))
}

const IM_COL32_WHITE = 0xFFFFFFF
const IM_COL32_BLACK = 0xFF00000
const IM_COL32_BLACK_TRANS = 0x0000000

// The maximum line width to bake anti-aliased textures for. Build atlas with ImFontAtlasFlags_NoBakedLines to disable baking.
const IM_DRAWLIST_TEX_LINES_WIDTH_MAX = 63

// ImDrawCallback: Draw callbacks for advanced uses [configurable type: override in imconfig.h]
// NB: You most likely do NOT need to use draw callbacks just to create your own widget or customized UI rendering,
// you can poke into the draw list for that! Draw callback may be useful for example to:
//  A) Change your GPU render state,
//  B) render a complex 3D scene inside a UI element without an intermediate texture/render target, etc.
// The expected behavior from your rendering function is 'if (cmd.UserCallback != NULL) { cmd) cmd.UserCallback(parent_list, } else { RenderTriangles() }'
// If you want to override the signature of ImDrawCallback, you can simply use e.g. '#define ImDrawCallback MyDrawCallback' (in imconfig.h) + update rendering backend accordingly.
type ImDrawCallback func(parent_list *ImDrawList, cmd *ImDrawCmd)

// Typically, 1 command = 1 GPU draw call (unless command is a callback)
// - VtxOffset/IdxOffset: When 'io.BackendFlags & ImGuiBackendFlags_RendererHasVtxOffset' is enabled,
//   those fields allow us to render meshes larger than 64K vertices while keeping 16-bit indices.
//   Pre-1.71 backends will typically ignore the VtxOffset/IdxOffset fields.
// - The ClipRect/TextureId/VtxOffset fields must be contiguous as we memcmp() them together (this is asserted for).
type ImDrawCmd struct {
	ClipRect         ImVec4         // 4*4  // Clipping rectangle (x1, y1, x2, y2). Subtract ImDrawData->DisplayPos to get clipping rectangle in "viewport" coordinates
	TextureId        ImTextureID    // 4-8  // User-provided texture ID. Set by user in ImfontAtlas::SetTexID() for fonts or passed to Image*() functions. Ignore if never using images or multiple fonts atlas.
	VtxOffset        uint           // 4    // Start offset in vertex buffer. ImGuiBackendFlags_RendererHasVtxOffset: always 0, otherwise may be >0 to support meshes larger than 64K vertices with 16-bit indices.
	IdxOffset        uint           // 4    // Start offset in index buffer. Always equal to sum of ElemCount drawn so far.
	ElemCount        uint           // 4    // Number of indices (multiple of 3) to be rendered as triangles. Vertices are stored in the callee ImDrawList's vtx_buffer[] array, indices in idx_buffer[].
	UserCallback     ImDrawCallback // 4-8  // If != NULL, call the function instead of rendering the vertices. clip_rect and texture_id will be set normally.
	UserCallbackData interface{}    // 4-8  // The draw callback code can access this.
}

func (this *ImDrawCmd) GetTexID() ImTextureID {
	return this.TextureId
}

type ImDrawVert struct {
	pos ImVec2
	uv  ImVec2
	col ImU32
}

func ImDrawVertSizeAndOffset() (size, o1, o2, o3 uintptr) {
	return unsafe.Sizeof(ImDrawVert{}),
		unsafe.Offsetof(ImDrawVert{}.pos),
		unsafe.Offsetof(ImDrawVert{}.uv),
		unsafe.Offsetof(ImDrawVert{}.col)
}

// [Internal] For use by ImDrawList
type ImDrawCmdHeader struct {
	ClipRect  ImVec4
	TextureId ImTextureID
	VtxOffset uint
}

// [Internal] For use by ImDrawListSplitter
type ImDrawChannel struct {
	_CmdBuffer []ImDrawCmd
	_IdxBuffer []ImDrawIdx
}

// Split/Merge functions are used to split the draw list into different layers which can be drawn into out of order.
// This is used by the Columns/Tables API, so items of each column can be batched together in a same draw call.
type ImDrawListSplitter struct {
	_Current  int             // Current channel number (0)
	_Count    int             // Number of active channels (1+)
	_Channels []ImDrawChannel // Draw channels (not resized down so _Count might be < Channels.Size)
}

func (this *ImDrawListSplitter) Clear() {
	this._Current = 0
	this._Count = 1 // Do not clear Channels[] so our allocations are reused next frame
}
func (this *ImDrawListSplitter) ClearFreeMemory()                       { panic("not implemented") }
func (this *ImDrawListSplitter) Split(draw_list *ImDrawList, count int) { panic("not implemented") }
func (this *ImDrawListSplitter) Merge(draw_list *ImDrawList)            { panic("not implemented") }
func (this *ImDrawListSplitter) SetCurrentChannel(draw_list *ImDrawList, channel_idx int) {
	panic("not implemented")
}

// Draw command list
// This is the low-level list of polygons that ImGui:: functions are filling. At the end of the frame,
// all command lists are passed to your ImGuiIO::RenderDrawListFn function for rendering.
// Each dear imgui window contains its own ImDrawList. You can use ImGui::GetWindowDrawList() to
// access the current window draw list and draw custom primitives.
// You can interleave normal ImGui:: calls and adding primitives to the current draw list.
// In single viewport mode, top-left is == GetMainViewport()->Pos (generally 0,0), bottom-right is == GetMainViewport()->Pos+Size (generally io.DisplaySize).
// You are totally free to apply whatever transformation matrix to want to the data (depending on the use of the transformation you may want to apply it to ClipRect as well!)
// Important: Primitives are always added to the list and not culled (culling is done at higher-level by ImGui:: functions), if you use this API a lot consider coarse culling your drawn objects.
type ImDrawList struct {
	// This is what you have to render
	CmdBuffer []ImDrawCmd     // Draw commands. Typically 1 command = 1 GPU draw call, unless the command is a callback.
	IdxBuffer []ImDrawIdx     // Index buffer. Each command consume ImDrawCmd::ElemCount of those
	VtxBuffer []ImDrawVert    // Vertex buffer.
	Flags     ImDrawListFlags // Flags, you may poke into these to adjust anti-aliasing settings per-primitive.

	// [Internal, used while building lists]
	_VtxCurrentIdx  uint                  // [Internal] generally == VtxBuffer.Size unless we are past 64K vertices, in which case this gets reset to 0.
	_Data           *ImDrawListSharedData // Pointer to shared draw data (you can use ImGui::GetDrawListSharedData() to get the one from current ImGui context)
	_OwnerName      string                // Pointer to owner window's name for debugging
	_VtxWritePtr    int                   // [Internal] point within VtxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_IdxWritePtr    int                   // [Internal] point within IdxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_ClipRectStack  []ImVec4              // [Internal]
	_TextureIdStack []ImTextureID         // [Internal]
	_Path           []ImVec2              // [Internal] current path building
	_CmdHeader      ImDrawCmdHeader       // [Internal] template of active commands. Fields should match those of CmdBuffer.back().
	_Splitter       ImDrawListSplitter    // [Internal] for channels api (note: prefer using your own persistent instance of ImDrawListSplitter!)
	_FringeScale    float                 // [Internal] anti-alias fringe is scaled by this value, this helps to keep things sharp while zooming at vertex buffer content
}

func NewImDrawList(shared_data *ImDrawListSharedData) ImDrawList {
	return ImDrawList{
		_Data: shared_data,
	}
}

func (this *ImDrawList) PushClipRect(cr_min, cr_max ImVec2, intersect_with_current_clip_rect bool) {
	var cr = ImVec4{cr_min.x, cr_min.y, cr_max.x, cr_max.y}
	if intersect_with_current_clip_rect {
		var current ImVec4 = this._CmdHeader.ClipRect
		if cr.x < current.x {
			cr.x = current.x
		}
		if cr.y < current.y {
			cr.y = current.y
		}
		if cr.z > current.z {
			cr.z = current.z
		}
		if cr.w > current.w {
			cr.w = current.w
		}
	}
	cr.z = ImMax(cr.x, cr.z)
	cr.w = ImMax(cr.y, cr.w)

	this._ClipRectStack = append(this._ClipRectStack, cr)
	this._CmdHeader.ClipRect = cr
	this._OnChangedClipRect()
}

// Render-level scissoring. This is passed down to your render function but not used for CPU-side coarse clipping. Prefer using higher-level ImGui::PushClipRect() to affect logic (hit-testing and widget culling)
func (this *ImDrawList) PushClipRectFullScreen() { panic("not implemented") }

func (this *ImDrawList) PopClipRect() {
	this._ClipRectStack = this._ClipRectStack[len(this._ClipRectStack)-1:]
	if len(this._ClipRectStack) == 0 {
		this._CmdHeader.ClipRect = this._Data.ClipRectFullscreen
	} else {
		this._CmdHeader.ClipRect = this._ClipRectStack[len(this._ClipRectStack)-1]
	}
	this._OnChangedClipRect()
}

func (this *ImDrawList) PopTextureID() { panic("not implemented") }
func (this *ImDrawList) GetClipRectMin() ImVec2 {
	var cr *ImVec4 = &this._ClipRectStack[len(this._ClipRectStack)-1]
	return ImVec2{cr.x, cr.y}
}
func (this *ImDrawList) GetClipRectMax() ImVec2 {
	var cr *ImVec4 = &this._ClipRectStack[len(this._ClipRectStack)-1]
	return ImVec2{cr.x, cr.y}
}

// Primitives
// - For rectangular primitives, "p_min" and "p_max" represent the upper-left and lower-right corners.
// - For circle primitives, use "num_segments == 0" to automatically calculate tessellation (preferred).
//   In older versions (until Dear ImGui 1.77) the AddCircle functions defaulted to num_segments == 12.
//   In future versions we will use textures to provide cheaper and higher-quality circles.
//   Use AddNgon() and AddNgonFilled() functions if you need to guaranteed a specific number of sides.
func (this *ImDrawList) AddLine(p1 *ImVec2, p2 *ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}
	this.PathLineTo(p1.Add(ImVec2{0.5, 0.5}))
	this.PathLineTo(p2.Add(ImVec2{0.5, 0.5}))
	this.PathStroke(col, 0, thickness)
}

func (this *ImDrawList) AddRectFilledMultiColor(cp_min ImVec2, p_max ImVec2, col_upr_left, col_upr_right, col_bot_right, col_bot_left ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) AddQuad(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	panic("not implemented")
}
func (this *ImDrawList) AddQuadFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) AddTriangle(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	panic("not implemented")
}

func (this *ImDrawList) AddCircle(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	panic("not implemented")
}

func (this *ImDrawList) AddNgon(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	panic("not implemented")
}
func (this *ImDrawList) AddNgonFilled(center ImVec2, radius float, col ImU32, num_segments int) {
	panic("not implemented")
}

func (this *ImDrawList) AddBezierCubic(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32, thickness float, num_segments int) {
	panic("not implemented")
} // Cubic Bezier (4 control points)
func (this *ImDrawList) AddBezierQuadratic(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32, thickness float, num_segments int) {
	panic("not implemented")
} // Quadratic Bezier (3 control points)

// Image primitives
// - Read FAQ to understand what ImTextureID is.
// - "p_min" and "p_max" represent the upper-left and lower-right corners of the rectangle.
// - "uv_min" and "uv_max" represent the normalized texture coordinates to use for those corners. Using (0,0)->(1,1) texture coordinates will generally display the entire texture.
func (this *ImDrawList) AddImage(user_texture_id ImTextureID, cp_min ImVec2, p_max ImVec2, uv_min *ImVec2, uv_max *ImVec2, col ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) AddImageQuad(user_texture_id ImTextureID, p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, uv1 *ImVec2, uv2 *ImVec2 /*= ImVec2(1, 0)*/, uv3 ImVec2 /*ImVec2(1, 1)*/, uv4 ImVec2 /*ImVec2(0, 1)*/, col ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) AddImageRounded(user_texture_id ImTextureID, cp_min ImVec2, p_max ImVec2, uv_min, uv_max *ImVec2, col ImU32, rounding float, flags ImDrawFlags) {
	panic("not implemented")
}

// Stateful path API, add points then finish with PathFillConvex() or PathStroke()
func (this *ImDrawList) PathClear() {
	this._Path = this._Path[:0]
}
func (this *ImDrawList) PathLineTo(pos ImVec2) {
	this._Path = append(this._Path, pos)
}
func (this *ImDrawList) PathLineToMergeDuplicate(pos ImVec2) {
	if len(this._Path) == 0 || this._Path[len(this._Path)-1] == pos {
		this._Path = append(this._Path, pos)
	}
}

// Note: Anti-aliased filling requires points to be in clockwise order.
func (this *ImDrawList) PathFillConvex(col ImU32) {
	this.AddConvexPolyFilled(this._Path, int(len(this._Path)), col)
	this._Path = this._Path[:0]
}

func (this *ImDrawList) PathStroke(col ImU32, flags ImDrawFlags, thickness float /*= 1.0f*/) {
	this.AddPolyline(this._Path, int(len(this._Path)), col, flags, thickness)
	this._Path = this._Path[:0]
}

func (this *ImDrawList) PathBezierCubicCurveTo(p2 *ImVec2, p3 ImVec2, p4 ImVec2, num_segments int) {
	panic("not implemented")
} // Cubic Bezier (4 control points)
func (this *ImDrawList) PathBezierQuadraticCurveTo(p2 *ImVec2, p3 ImVec2, num_segments int) {
	panic("not implemented")
} // Quadratic Bezier (3 control points)

func FixRectCornerFlags(flags ImDrawFlags) ImDrawFlags {
	// If this triggers, please update your code replacing hardcoded values with new ImDrawFlags_RoundCorners* values.
	// Note that ImDrawFlags_Closed (== 0x01) is an invalid flag for AddRect(), AddRectFilled(), PathRect() etc...
	IM_ASSERT_USER_ERROR((flags&0x0F) == 0, "Misuse of legacy hardcoded ImDrawCornerFlags values!")

	if (flags & ImDrawFlags_RoundCornersMask_) == 0 {
		flags |= ImDrawFlags_RoundCornersAll
	}

	return flags
}

func (this *ImDrawList) PathRect(a, b *ImVec2, rounding float, flags ImDrawFlags) {
	flags = FixRectCornerFlags(flags)

	var xamount, yamount float = 1, 1
	if ((flags & ImDrawFlags_RoundCornersTop) == ImDrawFlags_RoundCornersTop) || ((flags & ImDrawFlags_RoundCornersBottom) == ImDrawFlags_RoundCornersBottom) {
		xamount = 0.5
	}
	if ((flags & ImDrawFlags_RoundCornersLeft) == ImDrawFlags_RoundCornersLeft) || ((flags & ImDrawFlags_RoundCornersRight) == ImDrawFlags_RoundCornersRight) {
		yamount = 0.5
	}

	rounding = ImMin(rounding, ImFabs(b.x-a.x)*(xamount)-1.0)
	rounding = ImMin(rounding, ImFabs(b.y-a.y)*(yamount)-1.0)

	if rounding <= 0.0 || (flags&ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone {
		this.PathLineTo(*a)
		this.PathLineTo(ImVec2{b.x, a.y})
		this.PathLineTo(*b)
		this.PathLineTo(ImVec2{a.x, b.y})
	} else {
		var rounding_tl, rounding_tr, rounding_br, rounding_bl float
		if (flags & ImDrawFlags_RoundCornersTopLeft) != 0 {
			rounding_tl = rounding
		}
		if (flags & ImDrawFlags_RoundCornersTopRight) != 0 {
			rounding_tr = rounding
		}
		if (flags & ImDrawFlags_RoundCornersBottomRight) != 0 {
			rounding_br = rounding
		}
		if (flags & ImDrawFlags_RoundCornersBottomLeft) != 0 {
			rounding_bl = rounding
		}
		this.PathArcToFast(ImVec2{a.x + rounding_tl, a.y + rounding_tl}, rounding_tl, 6, 9)
		this.PathArcToFast(ImVec2{b.x - rounding_tr, a.y + rounding_tr}, rounding_tr, 9, 12)
		this.PathArcToFast(ImVec2{b.x - rounding_br, b.y - rounding_br}, rounding_br, 0, 3)
		this.PathArcToFast(ImVec2{a.x + rounding_bl, b.y - rounding_bl}, rounding_bl, 3, 6)
	}
}

// Advanced
func AddCallback(callback ImDrawCallback, callback_data interface{}) { panic("not implemented") } // Your rendering function must check for 'UserCallback' in ImDrawCmd and call the function instead of rendering triangles.

// This is useful if you need to forcefully create a new draw call (to allow for dependent rendering / blending). Otherwise primitives are merged into the same draw-call as much as possible
func (this *ImDrawList) AddDrawCmd() {
	var draw_cmd ImDrawCmd
	draw_cmd.ClipRect = this._CmdHeader.ClipRect // Same as calling ImDrawCmd_HeaderCopy()
	draw_cmd.TextureId = this._CmdHeader.TextureId

	draw_cmd.VtxOffset = this._CmdHeader.VtxOffset
	draw_cmd.IdxOffset = uint(len(this.IdxBuffer))

	IM_ASSERT(draw_cmd.ClipRect.x <= draw_cmd.ClipRect.z && draw_cmd.ClipRect.y <= draw_cmd.ClipRect.w)
	this.CmdBuffer = append(this.CmdBuffer, draw_cmd)
}

func (this *ImDrawList) CloneOutput() *ImDrawList { panic("not implemented") } // Create a clone of the CmdBuffer/IdxBuffer/VtxBuffer.

// Advanced: Channels
// - Use to split render into layers. By switching channels to can render out-of-order (e.g. submit FG primitives before BG primitives)
// - Use to minimize draw calls (e.g. if going back-and-forth between multiple clipping rectangles, prefer to append into separate channels then merge at the end)
// - FIXME-OBSOLETE: This API shouldn't have been in ImDrawList in the first place!
//   Prefer using your own persistent instance of ImDrawListSplitter as you can stack them.
//   Using the ImDrawList::ChannelsXXXX you cannot stack a split over another.
func (this *ImDrawList) ChannelsSplit(count int)  { this._Splitter.Split(this, count) }
func (this *ImDrawList) ChannelsMerge()           { this._Splitter.Merge(this) }
func (this *ImDrawList) ChannelsSetCurrent(n int) { this._Splitter.SetCurrentChannel(this, n) }

func (this *ImDrawList) PrimUnreserve(idx_count, vtx_count int) { panic("not implemented") }

// Axis aligned rectangle (composed of two triangles)
func (this *ImDrawList) PrimRectUV(a, b, uv_a, uv_b *ImVec2, col ImU32) { panic("not implemented") }
func (this *ImDrawList) PrimQuadUV(a, b, c, d *ImVec2, uv_a, uv_b, yv_c, uv_d *ImVec2, col ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) PrimWriteVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	this.VtxBuffer[this._VtxWritePtr].pos = pos
	this.VtxBuffer[this._VtxWritePtr].uv = *uv
	this.VtxBuffer[this._VtxWritePtr].col = col
	this._VtxWritePtr++

	this._VtxCurrentIdx++
}
func (this *ImDrawList) PrimWriteIdx(idx ImDrawIdx) {
	this.IdxBuffer[this._IdxWritePtr] = idx
	this._IdxWritePtr++
}

func (this *ImDrawList) PrimVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx))
	this.PrimWriteVtx(pos, uv, col)
} // Write vertex with unique index

func (this *ImDrawList) _ClearFreeMemory() { panic("not implemented") }

// Pop trailing draw command (used before merging or presenting to user)
// Note that this leaves the ImDrawList in a state unfit for further commands, as most code assume that CmdBuffer.Size > 0 && CmdBuffer.back().UserCallback == nil
func (this *ImDrawList) _PopUnusedDrawCmd() {
	if len(this.CmdBuffer) == 0 {
		return
	}
	var curr_cmd *ImDrawCmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	if curr_cmd.ElemCount == 0 && curr_cmd.UserCallback == nil {
		this.CmdBuffer = this.CmdBuffer[:len(this.CmdBuffer)-1]
	}
}

func (this *ImDrawList) _TryMergeDrawCmds() { panic("not implemented") }

func (this *ImDrawList) _OnChangedVtxOffset() { panic("not implemented") }

func (this *ImDrawList) _CalcCircleAutoSegmentCount(radius float) int {
	// Automatic segment count
	var radius_idx int = (int)(radius + 0.999999) // ceil to never reduce accuracy
	if radius_idx < int(len(this._Data.CircleSegmentCounts)) {
		return int(this._Data.CircleSegmentCounts[radius_idx]) // Use cached value
	} else {
		return int(IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(radius, this._Data.CircleSegmentMaxError))
	}
}

func (this *ImDrawList) _PathArcToN(center ImVec2, radius, a_min, a_max float, num_segments int) {
	if radius <= 0.0 {
		this._Path = append(this._Path, center)
		return
	}

	// Note that we are adding a point at both a_min and a_max.
	// If you are trying to draw a full closed circle you don't want the overlapping points!
	this._Path = reserveVec2Slice(this._Path, int(len(this._Path))+(num_segments+1))
	for i := int(0); i <= num_segments; i++ {
		var a float = a_min + ((float)(i)/(float)(num_segments))*(a_max-a_min)
		this._Path = append(this._Path, ImVec2{center.x + ImCos(a)*radius, center.y + ImSin(a)*radius})
	}
}

//-----------------------------------------------------------------------------
// [SECTION] Font API (ImFontConfig, ImFontGlyph, ImFontAtlasFlags, ImFontAtlas, ImFontGlyphRangesBuilder, ImFont)
//-----------------------------------------------------------------------------

type ImFontConfig struct {
	FontData             []byte    //          // TTF/OTF data
	FontDataSize         int       //          // TTF/OTF data size
	FontDataOwnedByAtlas bool      // true     // TTF/OTF data ownership taken by the container ImFontAtlas (will delete memory itself).
	FontNo               int       // 0        // Index of font within TTF/OTF file
	SizePixels           float     //          // Size in pixels for rasterizer (more or less maps to the resulting font height).
	OversampleH          int       // 3        // Rasterize at higher quality for sub-pixel positioning. Note the difference between 2 and 3 is minimal so you can reduce this to 2 to save memory. Read https://github.com/nothings/stb/blob/master/tests/oversample/README.md for details.
	OversampleV          int       // 1        // Rasterize at higher quality for sub-pixel positioning. This is not really useful as we don't use sub-pixel positions on the Y axis.
	PixelSnapH           bool      // false    // Align every glyph to pixel boundary. Useful e.g. if you are merging a non-pixel aligned font with the default font. If enabled, you can set OversampleH/V to 1.
	GlyphExtraSpacing    ImVec2    // 0, 0     // Extra spacing (in pixels) between glyphs. Only X axis is supported for now.
	GlyphOffset          ImVec2    // 0, 0     // Offset all glyphs from this font input.
	GlyphRanges          []ImWchar // NULL     // Pointer to a user-provided list of Unicode range (2 value per range, values are inclusive, zero-terminated list). THE ARRAY DATA NEEDS TO PERSIST AS LONG AS THE FONT IS ALIVE.
	GlyphMinAdvanceX     float     // 0        // Minimum AdvanceX for glyphs, set Min to align font icons, set both Min/Max to enforce mono-space font
	GlyphMaxAdvanceX     float     // FLT_MAX  // Maximum AdvanceX for glyphs
	MergeMode            bool      // false    // Merge into previous ImFont, so you can combine multiple inputs font into one ImFont (e.g. ASCII font + icons + Japanese glyphs). You may want to use GlyphOffset.y when merge font of different heights.
	FontBuilderFlags     uint      // 0        // Settings for custom font builder. THIS IS BUILDER IMPLEMENTATION DEPENDENT. Leave as zero if unsure.
	RasterizerMultiply   float     // 1.0f     // Brighten (>1.0f) or darken (<1.0f) font output. Brightening small fonts may be a good workaround to make them more readable.
	EllipsisChar         ImWchar   // -1       // Explicitly specify unicode codepoint of ellipsis character. When fonts are being merged first specified ellipsis will be used.

	// [Internal]
	Name    string
	DstFont *ImFont
}

func NewImFontConfig() ImFontConfig {
	return ImFontConfig{
		FontDataOwnedByAtlas: true,
		OversampleH:          3, // FIXME: 2 may be a better default?
		OversampleV:          1,
		GlyphMaxAdvanceX:     FLT_MAX,
		RasterizerMultiply:   1.0,
		EllipsisChar:         255,
	}
}

// Hold rendering data for one glyph.
// (Note: some language parsers may fail to convert the 31+1 bitfield members, in this case maybe drop store a single u32 or we can rework this)
type ImFontGlyph struct {
	Colored        uint  // Flag to indicate glyph is colored and should generally ignore tinting (make it usable with no shift on little-endian as this is used in loops)
	Visible        uint  // Flag to indicate glyph has no visible pixels (e.g. space). Allow early out when rendering.
	Codepoint      uint  // 0x0000..0x10FFFF
	AdvanceX       float // Distance to next character (= data from font + ImFontConfig::GlyphExtraSpacing.x baked in)
	X0, Y0, X1, Y1 float // Glyph corners
	U0, V0, U1, V1 float // Texture coordinates
}

func NewImFontGlyph() ImFontGlyph {
	return ImFontGlyph{}
}

// Helper to build glyph ranges from text/string data. Feed your application strings/characters to it then call BuildRanges().
// This is essentially a tightly packed of vector of 64k booleans = 8KB storage.
type ImFontGlyphRangesBuilder []ImU32

func NewImFontGlyphRangesBuilder() ImFontGlyphRangesBuilder {
	return make(ImFontGlyphRangesBuilder, (IM_UNICODE_CODEPOINT_MAX+1)/8)
}

func (this ImFontGlyphRangesBuilder) GetBit(n uintptr) bool {
	var off int = (int)(n >> 5)
	var mask ImU32 = uint(1) << (n & 31)
	return (this[off] & mask) != 0
}

func (this ImFontGlyphRangesBuilder) SetBit(n uintptr) {
	var off int = (int)(n >> 5)
	var mask ImU32 = uint(1) << (n & 31)
	this[off] |= mask
}

func (this ImFontGlyphRangesBuilder) AddChar(c ImWchar) {
	this.SetBit(uintptr(c))
}

func AddText(text, text_end string)     { panic("not implemented") } // Add string (each character of the UTF-8 string are added)
func AddRanges(ranges []ImWchar)        { panic("not implemented") } // Add ranges, e.g. builder.AddRanges(ImFontAtlas::GetGlyphRangesDefault()) to force add all of ASCII/Latin+Ext
func BuildRanges(out_ranges *[]ImWchar) { panic("not implemented") } // Output new ranges

// See ImFontAtlas::AddCustomRectXXX functions.
type ImFontAtlasCustomRect struct {
	Width, Height uint16  // Input    // Desired rectangle dimension
	X, Y          uint16  // Output   // Packed position in Atlas
	GlyphID       uint    // Input    // For custom font glyphs only (ID < 0x110000)
	GlyphAdvanceX float   // Input    // For custom font glyphs only: glyph xadvance
	GlyphOffset   ImVec2  // Input    // For custom font glyphs only: glyph display offset
	Font          *ImFont // Input    // For custom font glyphs only: target font
}

func NewImFontAtlasCustomRect() ImFontAtlasCustomRect {
	return ImFontAtlasCustomRect{
		Width:         0,
		Height:        0,
		X:             0xFFFF,
		Y:             0xFFFF,
		GlyphID:       0,
		GlyphAdvanceX: 0,
		GlyphOffset:   ImVec2{},
		Font:          nil,
	}
}

func (this ImFontAtlasCustomRect) IsPacked() bool {
	return this.X != 0xFFFF
}

// Load and rasterize multiple TTF/OTF fonts into a same texture. The font atlas will build a single texture holding:
//  - One or more fonts.
//  - Custom graphics data needed to render the shapes needed by Dear ImGui.
//  - Mouse cursor shapes for software cursor rendering (unless setting 'Flags |= ImFontAtlasFlags_NoMouseCursors' in the font atlas).
// It is the user-code responsibility to setup/build the atlas, then upload the pixel data into a texture accessible by your graphics api.
//  - Optionally, call any of the AddFont*** functions. If you don't call any, the default font embedded in the code will be loaded for you.
//  - Call GetTexDataAsAlpha8() or GetTexDataAsRGBA32() to build and retrieve pixels data.
//  - Upload the pixels data into a texture within your graphics system (see imgui_impl_xxxx.cpp examples)
//  - Call SetTexID(my_tex_id); and pass the pointer/identifier to your texture in a format natural to your graphics API.
//    This value will be passed back to you during rendering to identify the texture. Read FAQ entry about ImTextureID for more details.
// Common pitfalls:
// - If you pass a 'glyph_ranges' array to AddFont*** functions, you need to make sure that your array persist up until the
//   atlas is build (when calling GetTexData*** or Build()). We only copy the pointer, not the data.
// - Important: By default, AddFontFromMemoryTTF() takes ownership of the data. Even though we are not writing to it, we will free the pointer on destruction.
//   You can set font_cfg->FontDataOwnedByAtlas=false to keep ownership of your data and it won't be freed,
// - Even though many functions are suffixed with "TTF", OTF data is supported just as well.
// - This is an old API and it is currently awkward for those and and various other reasons! We will address them in the future!
type ImFontAtlas struct {
	//-------------------------------------------
	// Members
	//-------------------------------------------

	Flags           ImFontAtlasFlags // Build flags (see ImFontAtlasFlags_)
	TexID           ImTextureID      // User data to refer to the texture once it has been uploaded to user's graphic systems. It is passed back to you during rendering via the ImDrawCmd structure.
	TexDesiredWidth int              // Texture width desired by user before Build(). Must be a power-of-two. If have many glyphs your graphics API have texture size restrictions you may want to increase texture width to decrease height.
	TexGlyphPadding int              // Padding between glyphs within texture in pixels. Defaults to 1. If your rendering method doesn't rely on bilinear filtering you may set this to 0.
	Locked          bool             // Marked as Locked by ImGui::NewFrame() so attempt to modify the atlas will assert.

	// [Internal]
	// NB: Access texture data via GetTexData*() calls! Which will setup a default font for you.
	TexReady           bool                                        // Set when texture was built matching current font input
	TexPixelsUseColors bool                                        // Tell whether our texture data is known to use colors (rather than just alpha channel), in order to help backend select a format.
	TexPixelsAlpha8    []byte                                      // 1 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight
	TexPixelsRGBA32    []uint                                      // 4 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight * 4
	TexWidth           int                                         // Texture width calculated during Build().
	TexHeight          int                                         // Texture height calculated during Build().
	TexUvScale         ImVec2                                      // = (1.0f/TexWidth, 1.0f/TexHeight)
	TexUvWhitePixel    ImVec2                                      // Texture coordinates to a white pixel
	Fonts              []*ImFont                                   // Hold all the fonts returned by AddFont*. Fonts[0] is the default font upon calling ImGui::NewFrame(), use ImGui::PushFont()/PopFont() to change the current font.
	CustomRects        []ImFontAtlasCustomRect                     // Rectangles for packing custom texture data into the atlas.
	ConfigData         []ImFontConfig                              // Configuration data
	TexUvLines         [IM_DRAWLIST_TEX_LINES_WIDTH_MAX + 1]ImVec4 // UVs for baked anti-aliased lines

	// [Internal] Font builder
	FontBuilderIO    *ImFontBuilderIO // Opaque interface to a font builder (default to stb_truetype, can be changed to use FreeType by defining IMGUI_ENABLE_FREETYPE).
	FontBuilderFlags uint             // Shared flags (for all fonts) for custom font builder. THIS IS BUILD IMPLEMENTATION DEPENDENT. Per-font override is also available in ImFontConfig.

	// [Internal] Packing data
	PackIdMouseCursors int // Custom texture rectangle ID for white pixel and mouse cursors
	PackIdLines        int // Custom texture rectangle ID for baked anti-aliased lines
}

func NewImFontAtlas() ImFontAtlas {
	return ImFontAtlas{
		TexGlyphPadding:    1,
		PackIdMouseCursors: -1,
		PackIdLines:        -1,
	}
}

func (atlas *ImFontAtlas) AddFont(font_cfg *ImFontConfig) *ImFont {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render")
	IM_ASSERT(font_cfg.FontData != nil && font_cfg.FontDataSize > 0)
	IM_ASSERT(font_cfg.SizePixels > 0.0)

	// Create new font
	if !font_cfg.MergeMode {
		f := NewImFont()
		atlas.Fonts = append(atlas.Fonts, &f)
	} else {
		IM_ASSERT_USER_ERROR(len(atlas.Fonts) != 0, "Cannot use MergeMode for the first font") // When using MergeMode make sure that a font has already been added before. You can use ImGui::GetIO().Fonts.AddFontDefault() to add the default imgui font.
	}

	atlas.ConfigData = append(atlas.ConfigData, *font_cfg)
	var new_font_cfg *ImFontConfig = &atlas.ConfigData[len(atlas.ConfigData)-1]
	if new_font_cfg.DstFont == nil {
		new_font_cfg.DstFont = atlas.Fonts[len(atlas.ConfigData)-1]
	}
	if !new_font_cfg.FontDataOwnedByAtlas {
		new_font_cfg.FontData = make([]byte, new_font_cfg.FontDataSize)
		new_font_cfg.FontDataOwnedByAtlas = true
		copy(new_font_cfg.FontData, font_cfg.FontData[:(size_t)(new_font_cfg.FontDataSize)])
	}

	if new_font_cfg.DstFont.EllipsisChar == MaxImWchar {
		new_font_cfg.DstFont.EllipsisChar = font_cfg.EllipsisChar
	}

	// Invalidate texture
	atlas.TexReady = false
	atlas.ClearTexData()
	return new_font_cfg.DstFont
}

func (atlas *ImFontAtlas) AddFontFromFileTTF(filename string, size_pixels float, font_cfg *ImFontConfig, glyph_ranges []ImWchar) *ImFont {
	panic("not implemented")
}

func (atlas *ImFontAtlas) AddFontFromMemoryTTF(ttf_data []byte, ttf_size int, size_pixels float, font_cfg_template *ImFontConfig, glyph_ranges []ImWchar) *ImFont {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	var font_cfg ImFontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg = NewImFontConfig()
	}
	IM_ASSERT(font_cfg.FontData == nil)
	font_cfg.FontData = ttf_data
	font_cfg.FontDataSize = ttf_size
	if size_pixels > 0.0 {
		font_cfg.SizePixels = size_pixels
	}
	if glyph_ranges != nil {
		font_cfg.GlyphRanges = glyph_ranges
	}
	return atlas.AddFont(&font_cfg)
}

// Note: Transfer ownership of 'ttf_data' to ImFontAtlas! Will be deleted after destruction of the atlas. Set font_cfg->FontDataOwnedByAtlas=false to keep ownership of your data and it won't be freed.
// 'compressed_font_data_base85' still owned by caller. Compress with binary_to_compressed_c.cpp with -base85 parameter.
func (atlas *ImFontAtlas) ClearInputData() { panic("not implemented") } // Clear input data (all ImFontConfig structures including sizes, TTF data, glyph ranges, etc.) = all the data used to build the texture and fonts.

func (atlas *ImFontAtlas) ClearTexData() {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	atlas.TexPixelsAlpha8 = nil
	atlas.TexPixelsRGBA32 = nil
	atlas.TexPixelsUseColors = false
}

// Clear output texture data (CPU side). Saves RAM once the texture has been copied to graphics memory.
func (atlas *ImFontAtlas) ClearFonts() { panic("not implemented") } // Clear output font data (glyphs storage, UV coordinates).
func (atlas *ImFontAtlas) Clear()      { panic("not implemented") } // Clear all input and output.

func (atlas *ImFontAtlas) GetTexDataAsRGBA32(out_pixels *[]byte, out_width, out_height, iout_bytes_per_pixel *int) {
	panic("not implemented")
}                                                  // 4 bytes-per-pixel
func (atlas *ImFontAtlas) IsBuilt() bool           { return len(atlas.Fonts) > 0 && atlas.TexReady } // Bit ambiguous: used to detect when user didn't built texture but effectively we should check TexID != 0 except that would be backend dependent...
func (atlas *ImFontAtlas) SetTexID(id ImTextureID) { atlas.TexID = id }

//-------------------------------------------
// Glyph Ranges
//-------------------------------------------

// Helpers to retrieve list of common Unicode ranges (2 value per range, values are inclusive, zero-terminated list)
// NB: Make sure that your string are UTF-8 and NOT in your local code page. In C++11, you can create UTF-8 string literal using the u8"Hello world" syntax. See FAQ for details.
// NB: Consider using ImFontGlyphRangesBuilder to build glyph ranges from textual data.
func (atlas *ImFontAtlas) GetGlyphRangesKorean() []ImWchar                  { panic("not implemented") } // Default + Korean characters
func (atlas *ImFontAtlas) GetGlyphRangesJapanese() []ImWchar                { panic("not implemented") } // Default + Hiragana, Katakana, Half-Width, Selection of 2999 Ideographs
func (atlas *ImFontAtlas) GetGlyphRangesChineseFull() []ImWchar             { panic("not implemented") } // Default + Half-Width + Japanese Hiragana/Katakana + full set of about 21000 CJK Unified Ideographs
func (atlas *ImFontAtlas) GetGlyphRangesChineseSimplifiedCommon() []ImWchar { panic("not implemented") } // Default + Half-Width + Japanese Hiragana/Katakana + set of 2500 CJK Unified Ideographs for common simplified Chinese
func (atlas *ImFontAtlas) GetGlyphRangesCyrillic() []ImWchar                { panic("not implemented") } // Default + about 400 Cyrillic characters
func (atlas *ImFontAtlas) GetGlyphRangesThai() []ImWchar                    { panic("not implemented") } // Default + Thai characters
func (atlas *ImFontAtlas) GetGlyphRangesVietnamese() []ImWchar              { panic("not implemented") } // Default + Vietnamese characters

//-------------------------------------------
// [BETA] Custom Rectangles/Glyphs API
//-------------------------------------------

// You can request arbitrary rectangles to be packed into the atlas, for your own purposes.
// - After calling Build(), you can query the rectangle position and render your pixels.
// - If you render colored output, set 'atlas->TexPixelsUseColors = true' as this may help some backends decide of prefered texture format.
// - You can also request your rectangles to be mapped as font glyph (given a font + Unicode point),
//   so you can render e.g. custom colorful icons and use them as regular glyphs.
// - Read docs/FONTS.md for more details about using colorful icons.
// - Note: this API may be redesigned later in order to support multi-monitor varying DPI settings.
func (atlas *ImFontAtlas) AddCustomRectRegular(width, height int) int {
	IM_ASSERT(width > 0 && width <= 0xFFFF)
	IM_ASSERT(height > 0 && height <= 0xFFFF)
	var r ImFontAtlasCustomRect
	r.Width = (uint16)(width)
	r.Height = (uint16)(height)
	atlas.CustomRects = append(atlas.CustomRects, r)
	return int(len(atlas.CustomRects)) - 1 // Return index
}

func (atlas *ImFontAtlas) AddCustomRectFontGlyph(font *ImFont, id ImWchar, width, height int, advance_x float, offset *ImVec2) int {
	panic("not implemented")
}
func (atlas *ImFontAtlas) GetCustomRectByIndex(index int) *ImFontAtlasCustomRect {
	IM_ASSERT(index >= 0)
	return &atlas.CustomRects[index]
}

// [Internal]
func (atlas *ImFontAtlas) CalcCustomRectUV(rect *ImFontAtlasCustomRect, out_uv_min, out_uv_max *ImVec2) {
	panic("not implemented")
}
func (atlas *ImFontAtlas) GetMouseCursorTexData(cursor ImGuiMouseCursor, out_offset, out_size *ImVec2, out_uv_border [2]ImVec2, out_uv_fill [2]ImVec2) bool {
	panic("not implemented")
}

// Font runtime data and rendering
// ImFontAtlas automatically loads a default embedded font for you when you call GetTexDataAsAlpha8() or GetTexDataAsRGBA32().
type ImFont struct {
	// Members: Hot ~20/24 bytes (for CalcTextSize)
	IndexAdvanceX    []float // 12-16 // out //            // Sparse. Glyphs->AdvanceX in a directly indexable way (cache-friendly for CalcTextSize functions which only this this info, and are often bottleneck in large UI).
	FallbackAdvanceX float   // 4     // out // = FallbackGlyph->AdvanceX
	FontSize         float   // 4     // in  //            // Height of characters/line, set during loading (don't change after loading)

	// Members: Hot ~28/40 bytes (for CalcTextSize + render loop)
	IndexLookup   []ImWchar     // 12-16 // out //            // Sparse. Index glyphs by Unicode code-point.
	Glyphs        []ImFontGlyph // 12-16 // out //            // All glyphs.
	FallbackGlyph *ImFontGlyph  // 4-8   // out // = FindGlyph(FontFallbackChar)

	// Members: Cold ~32/40 bytes
	ContainerAtlas      *ImFontAtlas                                    // 4-8   // out //            // What we has been loaded into
	ConfigData          *ImFontConfig                                   // 4-8   // in  //            // Pointer within ContainerAtlas->ConfigData
	ConfigDataCount     short                                           // 2     // in  // ~ 1        // Number of ImFontConfig involved in creating this font. Bigger than 1 when merging multiple font sources into one ImFont.
	FallbackChar        ImWchar                                         // 2     // out // = FFFD/'?' // Character used if a glyph isn't found.
	EllipsisChar        ImWchar                                         // 2     // out // = '...'    // Character used for ellipsis rendering.
	DotChar             ImWchar                                         // 2     // out // = '.'      // Character used for ellipsis rendering (if a single '...' character isn't found)
	DirtyLookupTables   bool                                            // 1     // out //
	Scale               float                                           // 4     // in  // = 1.f      // Base font scale, multiplied by the per-window font scale which you can adjust with SetWindowFontScale()
	Ascent, Descent     float                                           // 4+4   // out //            // Ascent: distance from top to bottom of e.g. 'A' [0..FontSize]
	MetricsTotalSurface int                                             // 4     // out //            // Total surface in pixels to get an idea of the font rasterization/texture cost (not exact, we approximate the cost of padding between glyphs)
	Used4kPagesMap      [(IM_UNICODE_CODEPOINT_MAX + 1) / 4096 / 8]ImU8 // 2 bytes if ImWchar=ImWchar16, 34 bytes if ImWchar==ImWchar32. Store 1-bit for each block of 4K codepoints that has one active glyph. This is mainly used to facilitate iterations across all used codepoints.
}

// Methods
func NewImFont() ImFont {
	return ImFont{
		FallbackChar: MaxImWchar,
		EllipsisChar: MaxImWchar,
		DotChar:      MaxImWchar,
		Scale:        1,
	}
}

func (this *ImFont) GetCharAdvance(c ImWchar) float {
	if (int)(c) < int(len(this.IndexAdvanceX)) {
		return this.IndexAdvanceX[(int)(c)]
	}
	return this.FallbackAdvanceX
}

func (this *ImFont) IsLoaded() bool { return this.ContainerAtlas != nil }
func (this *ImFont) GetDebugName() string {
	if this.ConfigData != nil {
		return string(this.ConfigData.Name[:])
	}
	return "<unknown>"
}

func (this *ImFont) RenderChar(draw_list *ImDrawList, size float, pos ImVec2, col ImU32, c ImWchar) {
	panic("not implemented")
}

// [Internal] Don't use!

func (this *ImFont) AddRemapChar(dst, src ImWchar, overwite_dst bool /*= true*/) {
	panic("not implemented")
} // Makes 'dst' character/glyph points to 'src' character/glyph. Currently needs to be called AFTER fonts have been built.

func (this *ImFont) IsGlyphRangeUnused(c_begin, c_last uint) bool { panic("not implemented") }

// - Currently represents the Platform Window created by the application which is hosting our Dear ImGui windows.
// - In 'docking' branch with multi-viewport enabled, we extend this concept to have multiple active viewports.
// - In the future we will extend this concept further to also represent Platform Monitor and support a "no main platform window" operation mode.
// - About Main Area vs Work Area:
//   - Main Area = entire viewport.
//   - Work Area = entire viewport minus sections used by main menu bars (for platform windows), or by task bar (for platform monitor).
//   - Windows are generally trying to stay within the Work Area of their host viewport.
type ImGuiViewport struct {
	Flags    ImGuiViewportFlags // See ImGuiViewportFlags_
	Pos      ImVec2             // Main Area: Position of the viewport (Dear ImGui coordinates are the same as OS desktop/native coordinates)
	Size     ImVec2             // Main Area: Size of the viewport.
	WorkPos  ImVec2             // Work Area: Position of the viewport minus task bars, menus bars, status bars (>= Pos)
	WorkSize ImVec2             // Work Area: Size of the viewport minus task bars, menu bars, status bars (<= Size)

	DrawListsLastFrame [2]int         // Last frame number the background (0) and foreground (1) draw lists were used
	DrawLists          [2]*ImDrawList // Convenience background (0) and foreground (1) draw lists. We use them to draw software mouser cursor when io.MouseDrawCursor is set and to draw most debug overlays.
	DrawDataP          ImDrawData
	DrawDataBuilder    ImDrawDataBuilder

	WorkOffsetMin      ImVec2 // Work Area: Offset from Pos to top-left corner of Work Area. Generally (0,0) or (0,+main_menu_bar_height). Work Area is Full Area but without menu-bars/status-bars (so WorkArea always fit inside Pos/Size!)
	WorkOffsetMax      ImVec2 // Work Area: Offset from Pos+Size to bottom-right corner of Work Area. Generally (0,0) or (0,-status_bar_height).
	BuildWorkOffsetMin ImVec2 // Work Area: Offset being built during current frame. Generally >= 0.0f.
	BuildWorkOffsetMax ImVec2 // Work Area: Offset being built during current frame. Generally <= 0.0f.
}

// Helpers
func (this *ImGuiViewport) GetCenter() ImVec2 {
	return ImVec2{this.Pos.x + this.Size.x*0.5, this.Pos.y + this.Size.y*0.5}
}
func (this *ImGuiViewport) GetWorkCenter() ImVec2 {
	return ImVec2{this.WorkPos.x + this.WorkSize.x*0.5, this.WorkPos.y + this.WorkSize.y*0.5}
}
