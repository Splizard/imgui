package imgui

import (
	"unicode"
	"unsafe"
)

const IMGUI_VERSION = "1.85 WIP"
const IMGUI_VERSION_NUM = 18411

func IMGUI_CHECKVERSION() {
	panic("imgui: IMGUI_CHECKVERSION not supported")
}

const IMGUI_HAS_TABLE = true

// ImGuiCol Enums/Flags (declared as for int compatibility with old C++, to allow using as flags without overhead, and to not pollute the top of this file)
//   - Tip: Use your programming IDE navigation facilities on the names in the _central column_ below to find the actual flags/enum lists!
//     In Visual Studio IDE: CTRL+comma ("Edit.NavigateTo") can follow symbols in comments, whereas CTRL+F12 ("Edit.GoToImplementation") cannot.
//     With Visual Assist installed: ALT+G ("VAssistX.GoToImplementation") can also follow symbols in comments.
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

// ImGuiID Scalar data types
type ImGuiID = uint // A unique ID used by widgets (typically the result of hashing a stack of string)
type ImS8 = int8    // 8-bit signed integer
type ImU8 = uint8   // 8-bit uinteger
type ImS16 = int16  // 16-bit signed integer
type ImU16 = uint16 // 16-bit uinteger
type ImS32 = int    // 32-bit signed integer == int
type ImU32 = uint   // 32-bit uinteger (often used to store packed colors)
type ImS64 = int64  // 64-bit signed integer (pre and post C++11 with Visual Studio)
type ImU64 = uint64 // 64-bit uinteger (pre and post C++11 with Visual Studio)

// ImWchar16 Character types
// (we generally use UTF-8 encoded string in the API. This is storage specifically for a decoded character used for keyboard input and display)
type ImWchar16 = uint16 // A single decoded U16 character/code point. We encode them as multi bytes UTF-8 when used in strings.
type ImWchar32 = int    // A single decoded U32 character/code point. We encode them as multi bytes UTF-8 when used in strings.
type ImWchar = rune     // ImWchar [configurable type: override in imconfig.h with '#define IMGUI_USE_WCHAR32' to support Unicode planes 1-16]

// ImGuiInputTextCallback Callback and functions types
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
	UserData                any                 // = NULL           // Store your own data for retrieval by callbacks.

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
	BackendPlatformName     string // = NULL
	BackendRendererName     string // = NULL
	BackendPlatformUserData any    // = NULL           // User data for platform backend
	BackendRendererUserData any    // = NULL           // User data for renderer backend
	BackendLanguageUserData any    // = NULL           // User data for non C++ programming language backend

	// Optional: Access OS clipboard
	// (default to use native Win32 clipboard on Windows, otherwise uses a private clipboard. Override to access OS clipboard on other architectures)
	GetClipboardTextFn func(user_data any) string
	SetClipboardTextFn func(user_data any, text string)
	ClipboardUserData  any

	// Optional: Notify OS Input Method Editor of the screen position of your cursor for text input position (e.g. when using Japanese/Chinese IME on Windows)
	// (default to use native imm32 api on Windows)
	ImeSetInputScreenPosFn func(x, y int)
	ImeWindowHandle        any // = NULL           // (Windows) Set this to your HWND to get automatic IME cursor positioning.

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
	io.GetClipboardTextFn = GetClipboardTextFn_DefaultImpl // Platform dependent default implementations
	io.SetClipboardTextFn = SetClipboardTextFn_DefaultImpl
	io.ImeSetInputScreenPosFn = func(x, y int) {}

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

// ImGuiInputTextCallbackData Shared state of InputText(), passed as an argument to your callback when a ImGuiInputTextFlags_Callback* flag is used.
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
	UserData  any                 // What user passed to InputText()      // Read-only

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

func NewImGuiInputTextCallbackData() *ImGuiInputTextCallbackData {
	return new(ImGuiInputTextCallbackData)
}

// DeleteChars Public API to manipulate UTF-8 text
// We expose UTF-8 to the user (unlike the STB_TEXTEDIT_* functions which are manipulating wchar)
// FIXME: The existence of this rarely exercised code path is a bit of a nuisance.
func (this *ImGuiInputTextCallbackData) DeleteChars(pos, bytes_count int) {
	IM_ASSERT(pos+bytes_count <= this.BufTextLen)
	var dst = this.Buf[pos:]
	var src = this.Buf[pos+bytes_count:]

	copy(dst, src)

	if this.CursorPos >= pos+bytes_count {
		this.CursorPos -= bytes_count
	} else if this.CursorPos >= pos {
		this.CursorPos = pos
	}
	this.SelectionStart = this.CursorPos
	this.SelectionEnd = this.CursorPos
	this.BufDirty = true
	this.BufTextLen -= bytes_count
}

func (this *ImGuiInputTextCallbackData) InsertChars(pos int, new_text string) {
	var is_resizable = (this.Flags & ImGuiInputTextFlags_CallbackResize) != 0
	var new_text_len = int(len(new_text))
	if new_text_len+this.BufTextLen >= this.BufSize {
		if !is_resizable {
			return
		}

		// Contrary to STB_TEXTEDIT_INSERTCHARS() this is working in the UTF8 buffer, hence the mildly similar code (until we remove the U16 buffer altogether!)
		var g = GImGui
		var edit_state = &g.InputTextState
		IM_ASSERT(edit_state.ID != 0 && g.ActiveId == edit_state.ID)
		//IM_ASSERT(this.Buf == edit_state.TextA)
		var new_buf_size = this.BufTextLen + ImClampInt(new_text_len*4, 32, ImMaxInt(256, new_text_len)) + 1
		edit_state.TextA = append(edit_state.TextA, make([]byte, new_buf_size-int(len(edit_state.TextA)))...)

		this.Buf = edit_state.TextA
		this.BufSize = new_buf_size
		edit_state.BufCapacityA = new_buf_size
	}

	if this.BufTextLen != pos {
		copy(this.Buf[pos+new_text_len:], this.Buf[pos:])
	}
	copy(this.Buf[pos:], new_text)

	if this.CursorPos >= pos {
		this.CursorPos += new_text_len
	}
	this.SelectionStart = this.CursorPos
	this.SelectionEnd = this.CursorPos
	this.BufDirty = true
	this.BufTextLen += new_text_len
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

// ImGuiSizeCallbackData Resizing callback data to apply custom constraint. As enabled by SetNextWindowSizeConstraints(). Callback is called during the next Begin().
// NB: For basic min/max size constraint on each axis you don't need to use the callback! The SetNextWindowSizeConstraints() parameters are enough.
type ImGuiSizeCallbackData struct {
	UserData    any    // Read-only.   What user passed to SetNextWindowSizeConstraints()
	Pos         ImVec2 // Read-only.   Window position, for reference.
	CurrentSize ImVec2 // Read-only.   Current window size.
	DesiredSize ImVec2 // Read-write.  Desired size, based on user's mouse position. Write to this field to restrain resizing.
}

// ImGuiPayload Data payload for Drag and Drop operations: AcceptDragDropPayload(), GetDragDropPayload()
type ImGuiPayload struct {
	// Members
	Data     any // Data (copied and owned by dear imgui)
	DataSize int // Data size

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

// ImGuiTableColumnSortSpecs Sorting specification for one column of a table (sizeof == 12 bytes)
type ImGuiTableColumnSortSpecs struct {
	ColumnUserID  ImGuiID            // User id of the column (if specified by a TableSetupColumn() call)
	ColumnIndex   ImS16              // Index of the column
	SortOrder     ImS16              // Index within parent ImGuiTableSortSpecs (always stored in order starting from 0, tables sorted on a single criteria will always have a 0 here)
	SortDirection ImGuiSortDirection // ImGuiSortDirection_Ascending or ImGuiSortDirection_Descending (you can use this or SortSign, whichever is more convenient for your sort function)
}

func NewImGuiTableColumnSortSpecs() ImGuiTableColumnSortSpecs {
	return ImGuiTableColumnSortSpecs{}
}

// ImGuiTableSortSpecs Sorting specifications for a table (often handling sort specs for a single column, occasionally more)
// Obtained by calling TableGetSortSpecs().
// When 'SpecsDirty == true' you can sort your data. It will be true with sorting specs have changed since last call, or the first time.
// Make sure to set 'SpecsDirty = false' after sorting, else you may wastefully sort your data every frame!
type ImGuiTableSortSpecs struct {
	Specs      []ImGuiTableColumnSortSpecs // Pointer to sort spec array.
	SpecsCount int                         // Sort spec count. Most often 1. May be > 1 when ImGuiTableFlags_SortMulti is enabled. May be == 0 when ImGuiTableFlags_SortTristate is enabled.
	SpecsDirty bool                        // Set to true when specs have changed since last time! Use this to sort again, then clear the flag.
}

const IM_UNICODE_CODEPOINT_INVALID = unicode.ReplacementChar

// IM_UNICODE_CODEPOINT_MAX /*
const IM_UNICODE_CODEPOINT_MAX = 0xFFFF

// ImGuiOnceUponAFrame Helper: Execute a block of code at maximum once a frame. Convenient if you want to quickly create an UI within deep-nested code that runs multiple times every frame.
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
	var current_frame = GetFrameCount()
	if this.RefFrame == current_frame {
		return false
	}
	this.RefFrame = current_frame
	return true
}

const IM_COL32_R_SHIFT = 0
const IM_COL32_G_SHIFT = 8
const IM_COL32_B_SHIFT = 16
const IM_COL32_A_SHIFT = 24
const IM_COL32_A_MASK = 0xFF000000

func IM_COL32(R, G, B, A byte) ImU32 {
	return ((ImU32)(A) << IM_COL32_A_SHIFT) | ((ImU32)(B) << IM_COL32_B_SHIFT) | ((ImU32)(G) << IM_COL32_G_SHIFT) | ((ImU32)(R) << IM_COL32_R_SHIFT)
}

const IM_COL32_WHITE = 0xFFFFFFF
const IM_COL32_BLACK = 0xFF00000
const IM_COL32_BLACK_TRANS = 0x0000000

// IM_DRAWLIST_TEX_LINES_WIDTH_MAX The maximum line width to bake anti-aliased textures for. Build atlas with ImFontAtlasFlags_NoBakedLines to disable baking.
const IM_DRAWLIST_TEX_LINES_WIDTH_MAX = 63

// ImDrawCallback ImDrawCallback: Draw callbacks for advanced uses [configurable type: override in imconfig.h]
// NB: You most likely do NOT need to use draw callbacks just to create your own widget or customized UI rendering,
// you can poke into the draw list for that! Draw callback may be useful for example to:
//
//	A) Change your GPU render state,
//	B) render a complex 3D scene inside a UI element without an intermediate texture/render target, etc.
//
// The expected behavior from your rendering function is 'if (cmd.UserCallback != NULL) { cmd) cmd.UserCallback(parent_list, } else { RenderTriangles() }'
// If you want to override the signature of ImDrawCallback, you can simply use e.g. '#define ImDrawCallback MyDrawCallback' (in imconfig.h) + update rendering backend accordingly.
type ImDrawCallback func(parent_list *ImDrawList, cmd *ImDrawCmd)

// ImDrawCmd Typically, 1 command = 1 GPU draw call (unless command is a callback)
//   - VtxOffset/IdxOffset: When 'io.BackendFlags & ImGuiBackendFlags_RendererHasVtxOffset' is enabled,
//     those fields allow us to render meshes larger than 64K vertices while keeping 16-bit indices.
//     Pre-1.71 backends will typically ignore the VtxOffset/IdxOffset fields.
//   - The ClipRect/TextureId/VtxOffset fields must be contiguous as we memcmp() them together (this is asserted for).
type ImDrawCmd struct {
	ClipRect         ImVec4         // 4*4  // Clipping rectangle (x1, y1, x2, y2). Subtract ImDrawData->DisplayPos to get clipping rectangle in "viewport" coordinates
	TextureId        ImTextureID    // 4-8  // User-provided texture ID. Set by user in ImfontAtlas::SetTexID() for fonts or passed to Image*() functions. Ignore if never using images or multiple fonts atlas.
	VtxOffset        uint           // 4    // Start offset in vertex buffer. ImGuiBackendFlags_RendererHasVtxOffset: always 0, otherwise may be >0 to support meshes larger than 64K vertices with 16-bit indices.
	IdxOffset        uint           // 4    // Start offset in index buffer. Always equal to sum of ElemCount drawn so far.
	ElemCount        uint           // 4    // Number of indices (multiple of 3) to be rendered as triangles. Vertices are stored in the callee ImDrawList's vtx_buffer[] array, indices in idx_buffer[].
	UserCallback     ImDrawCallback // 4-8  // If != NULL, call the function instead of rendering the vertices. clip_rect and texture_id will be set normally.
	UserCallbackData any            // 4-8  // The draw callback code can access this.
}

func (this *ImDrawCmd) HeaderEquals(other *ImDrawCmd) bool {
	return this.ClipRect == other.ClipRect && this.TextureId == other.TextureId && this.VtxOffset == other.VtxOffset
}

func (this *ImDrawCmd) HeaderEqualsHeader(other *ImDrawCmdHeader) bool {
	return this.ClipRect == other.ClipRect && this.TextureId == other.TextureId && this.VtxOffset == other.VtxOffset
}

func (this *ImDrawCmd) HeaderCopyFromHeader(other ImDrawCmdHeader) {
	this.ClipRect = other.ClipRect
	this.TextureId = other.TextureId
	this.VtxOffset = other.VtxOffset
}

func (this *ImDrawCmd) GetTexID() ImTextureID {
	return this.TextureId
}

type ImDrawVert struct {
	Pos ImVec2
	Uv  ImVec2
	Col ImU32
}

func ImDrawVertSizeAndOffset() (size, o1, o2, o3 uintptr) {
	return unsafe.Sizeof(ImDrawVert{}),
		unsafe.Offsetof(ImDrawVert{}.Pos),
		unsafe.Offsetof(ImDrawVert{}.Uv),
		unsafe.Offsetof(ImDrawVert{}.Col)
}

// ImDrawCmdHeader [Internal] For use by ImDrawList
type ImDrawCmdHeader struct {
	ClipRect  ImVec4
	TextureId ImTextureID
	VtxOffset uint
}

// ImDrawChannel [Internal] For use by ImDrawListSplitter
type ImDrawChannel struct {
	_CmdBuffer []ImDrawCmd
	_IdxBuffer []ImDrawIdx
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
		EllipsisChar:         (ImWchar)(-1),
	}
}

// ImFont Font runtime data and rendering
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
	ConfigData          []ImFontConfig                                  // 4-8   // in  //            // Pointer within ContainerAtlas->ConfigData
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

// NewImFont Methods
func NewImFont() ImFont {
	return ImFont{
		FallbackChar: (ImWchar)(-1),
		EllipsisChar: (ImWchar)(-1),
		DotChar:      (ImWchar)(-1),
		Scale:        1,
	}
}

func (f *ImFont) GetCharAdvance(c ImWchar) float {
	if (int)(c) < int(len(f.IndexAdvanceX)) {
		return f.IndexAdvanceX[(int)(c)]
	}
	return f.FallbackAdvanceX
}

func (f *ImFont) IsLoaded() bool { return f.ContainerAtlas != nil }
func (f *ImFont) GetDebugName() string {
	if f.ConfigData != nil {
		return string(f.ConfigData[0].Name[:])
	}
	return "<unknown>"
}

func (f *ImFont) RenderChar(draw_list *ImDrawList, size float, pos ImVec2, col ImU32, c ImWchar) {
	var glyph = f.FindGlyph(c)
	if glyph == nil || glyph.Visible == 0 {
		return
	}
	if glyph.Colored != 0 {
		col |= IM_COL32_A_MASK
	}
	var scale float = 1.0
	if size >= 0.0 {
		scale = size / f.FontSize
	}
	pos.x = IM_FLOOR(pos.x)
	pos.y = IM_FLOOR(pos.y)
	draw_list.PrimReserve(6, 4)
	draw_list.PrimRectUV(&ImVec2{pos.x + glyph.X0*scale, pos.y + glyph.Y0*scale}, &ImVec2{pos.x + glyph.X1*scale, pos.y + glyph.Y1*scale}, &ImVec2{glyph.U0, glyph.V0}, &ImVec2{glyph.U1, glyph.V1}, col)

}

// [Internal] Don't use!

// AddRemapChar Makes 'dst' character/glyph points to 'src' character/glyph. Currently needs to be called AFTER fonts have been built.
func (f *ImFont) AddRemapChar(dst, src ImWchar, overwrite_dst bool /*= true*/) {
	IM_ASSERT(len(f.IndexLookup) > 0) // Currently f can only be called AFTER the font has been built, aka after calling ImFontAtlas::GetTexDataAs*() function.
	var index_size = ImWchar(len(f.IndexLookup))

	if dst < index_size && f.IndexLookup[dst] == (ImWchar)(-1) && !overwrite_dst { // 'dst' already exists
		return
	}
	if src >= index_size && dst >= index_size { // both 'dst' and 'src' don't exist . no-op
		return
	}

	f.GrowIndex(dst + 1)
	if src < index_size {
		f.IndexLookup[dst] = f.IndexLookup[src]
		f.IndexAdvanceX[dst] = f.IndexAdvanceX[src]
	} else {
		f.IndexLookup[dst] = (ImWchar)(-1)
		f.IndexAdvanceX[dst] = 1
	}
}

// IsGlyphRangeUnused API is designed this way to avoid exposing the 4K page size
// e.g. use with IsGlyphRangeUnused(0, 255)
func (f *ImFont) IsGlyphRangeUnused(c_begin, c_last uint) bool {
	var page_begin = c_begin / 4096
	var page_last = c_last / 4096
	for page_n := page_begin; page_n <= page_last; page_n++ {
		if uintptr(page_n>>3) < unsafe.Sizeof(f.Used4kPagesMap) {
			if f.Used4kPagesMap[page_n>>3]&(1<<(page_n&7)) != 0 {
				return false
			}
		}
	}
	return true
}

// ImGuiViewport - Currently represents the Platform Window created by the application which is hosting our Dear ImGui windows.
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

// GetCenter Helpers
func (p *ImGuiViewport) GetCenter() ImVec2 {
	return ImVec2{p.Pos.x + p.Size.x*0.5, p.Pos.y + p.Size.y*0.5}
}
func (p *ImGuiViewport) GetWorkCenter() ImVec2 {
	return ImVec2{p.WorkPos.x + p.WorkSize.x*0.5, p.WorkPos.y + p.WorkSize.y*0.5}
}
