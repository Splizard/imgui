package imgui

import (
	"math"
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
type ImWchar16 = uint16  // A single decoded U16 character/code point. We encode them as multi bytes UTF-8 when used in strings.
type ImWchar32 = int     // A single decoded U32 character/code point. We encode them as multi bytes UTF-8 when used in strings.
type ImWchar = ImWchar16 // ImWchar [configurable type: override in imconfig.h with '#define IMGUI_USE_WCHAR32' to support Unicode planes 1-16]

// Callback and functions types
type ImGuiInputTextCallback func(data *ImGuiInputTextCallbackData) int // Callback function for ImGui::InputText()
type ImGuiSizeCallback func(data *ImGuiSizeCallbackData)               // Callback function for ImGui::SetNextWindowSizeConstraints()

//-----------------------------------------------------------------------------
// [SECTION] ImGuiStyle
//-----------------------------------------------------------------------------
// You may modify the ImGui::GetStyle() main instance during initialization and before NewFrame().
// During the frame, use ImGui::PushStyleVar(ImGuiStyleVar_XXXX)/PopStyleVar() to alter the main style values,
// and ImGui::PushStyleColor(ImGuiCol_XXX)/PopStyleColor() for colors.
//-----------------------------------------------------------------------------
type ImGuiStyle struct {
	Alpha                      float    // Global alpha applies to everything in Dear ImGui.
	DisabledAlpha              float    // Additional alpha multiplier applied by BeginDisabled(). Multiply over current value of Alpha.
	WindowPadding              ImVec2   // Padding within a window.
	WindowRounding             float    // Radius of window corners rounding. Set to 0.0f to have rectangular windows. Large values tend to lead to variety of artifacts and are not recommended.
	WindowBorderSize           float    // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	WindowMinSize              ImVec2   // Minimum window size. This is a global setting. If you want to constraint individual windows, use SetNextWindowSizeConstraints().
	WindowTitleAlign           ImVec2   // Alignment for title bar text. Defaults to (0.0f,0.5f) for left-aligned,vertically centered.
	WindowMenuButtonPosition   ImGuiDir // Side of the collapsing/docking button in the title bar (None/Left/Right). Defaults to ImGuiDir_Left.
	ChildRounding              float    // Radius of child window corners rounding. Set to 0.0f to have rectangular windows.
	ChildBorderSize            float    // Thickness of border around child windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	PopupRounding              float    // Radius of popup window corners rounding. (Note that tooltip windows use WindowRounding)
	PopupBorderSize            float    // Thickness of border around popup/tooltip windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	FramePadding               ImVec2   // Padding within a framed rectangle (used by most widgets).
	FrameRounding              float    // Radius of frame corners rounding. Set to 0.0f to have rectangular frame (used by most widgets).
	FrameBorderSize            float    // Thickness of border around frames. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	ItemSpacing                ImVec2   // Horizontal and vertical spacing between widgets/lines.
	ItemInnerSpacing           ImVec2   // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label).
	CellPadding                ImVec2   // Padding within a table cell
	TouchExtraPadding          ImVec2   // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	IndentSpacing              float    // Horizontal indentation when e.g. entering a tree node. Generally == (FontSize + FramePadding.x*2).
	ColumnsMinSpacing          float    // Minimum horizontal spacing between two columns. Preferably > (FramePadding.x + 1).
	ScrollbarSize              float    // Width of the vertical scrollbar, Height of the horizontal scrollbar.
	ScrollbarRounding          float    // Radius of grab corners for scrollbar.
	GrabMinSize                float    // Minimum width/height of a grab box for slider/scrollbar.
	GrabRounding               float    // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	LogSliderDeadzone          float    // The size in pixels of the dead-zone around zero on logarithmic sliders that cross zero.
	TabRounding                float    // Radius of upper corners of a tab. Set to 0.0f to have rectangular tabs.
	TabBorderSize              float    // Thickness of border around tabs.
	TabMinWidthForCloseButton  float    // Minimum width for close button to appears on an unselected tab when hovered. Set to 0.0f to always show when hovering, set to FLT_MAX to never show close button unless selected.
	ColorButtonPosition        ImGuiDir // Side of the color button in the ColorEdit4 widget (left/right). Defaults to ImGuiDir_Right.
	ButtonTextAlign            ImVec2   // Alignment of button text when button is larger than text. Defaults to (0.5f, 0.5f) (centered).
	SelectableTextAlign        ImVec2   // Alignment of selectable text. Defaults to (0.0f, 0.0f) (top-left aligned). It's generally important to keep this left-aligned if you want to lay multiple items on a same line.
	DisplayWindowPadding       ImVec2   // Window position are clamped to be visible within the display area or monitors by at least this amount. Only applies to regular windows.
	DisplaySafeAreaPadding     ImVec2   // If you cannot see the edges of your screen (e.g. on a TV) increase the safe area padding. Apply to popups/tooltips as well regular windows. NB: Prefer configuring your TV sets correctly!
	MouseCursorScale           float    // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	AntiAliasedLines           bool     // Enable anti-aliased lines/borders. Disable if you are really tight on CPU/GPU. Latched at the beginning of the frame (copied to ImDrawList).
	AntiAliasedLinesUseTex     bool     // Enable anti-aliased lines/borders using textures where possible. Require backend to render with bilinear filtering. Latched at the beginning of the frame (copied to ImDrawList).
	AntiAliasedFill            bool     // Enable anti-aliased edges around filled shapes (rounded rectangles, circles, etc.). Disable if you are really tight on CPU/GPU. Latched at the beginning of the frame (copied to ImDrawList).
	CurveTessellationTol       float    // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	CircleTessellationMaxError float    // Maximum error (in pixels) allowed when using AddCircle()/AddCircleFilled() or drawing rounded corner rectangles with no explicit segment count specified. Decrease for higher quality but more geometry.
	Colors                     [ImGuiCol_COUNT]ImVec4
}

func NewImGuiStyle() *ImGuiStyle {
	panic("not implemented")
}

func (ImGuiStyle) ScaleAllSizes(scale_factor float) { panic("not implemented") }

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

func NewImGuiIO() *ImGuiIO { panic("not implemented") }

func (*ImGuiIO) AddInputCharacter(c uint)           { panic("not implemented") } // Queue new character input
func (*ImGuiIO) AddInputCharacterUTF16(c ImWchar16) { panic("not implemented") } // Queue new character input from an UTF-16 character, it can be a surrogate
func (*ImGuiIO) AddInputCharactersUTF8(str string)  { panic("not implemented") } // Queue new characters input from an UTF-8 string
func (*ImGuiIO) ClearInputCharacters()              { panic("not implemented") } // Clear the text input buffer manually
func (*ImGuiIO) AddFocusEvent(focused bool)         { panic("not implemented") } // Notifies Dear ImGui when hosting platform windows lose or gain input focus

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

const IM_UNICODE_CODEPOINT_INVALID = 0xFFFD

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

type ImGuiTextRange struct {
	b string
	e string
}

func (this ImGuiTextRange) split(separator byte, out []ImGuiTextRange) { panic("not implemented") }

// Helper: Parse and apply text filters. In format "aaaaa[,bbbb][,ccccc]"
type ImGuiTextFilter struct {
	InputBuf  [256]byte
	Filters   []ImGuiTextRange
	CountGrep int
}

func NewImGuiTextFilter(default_filter string) ImGuiTextFilter { panic("not implemented") }

func (this *ImGuiTextFilter) Draw(label string /*= "Filter (inc,-exc)"*/, width float) bool {
	panic("not implemented")
}                                                                   // Helper calling InputText+Build
func (this *ImGuiTextFilter) PassFilter(text, text_end string) bool { panic("not implemented") }
func (this *ImGuiTextFilter) Build()                                { panic("not implemented") }
func (this *ImGuiTextFilter) Clear() {
	this.InputBuf[0] = 0
	this.Build()
}
func (this *ImGuiTextFilter) IsActive() bool { return len(this.Filters) > 0 }

// Helper: Growable text buffer for logging/accumulating text
// (this could be called 'ImGuiTextBuilder' / 'ImGuiStringBuilder')
type ImGuiTextBuffer []byte

// [Internal]
type ImGuiStoragePair struct {
	key ImGuiID
	val int
}

func (this *ImGuiStoragePair) SetInt(key ImGuiID, val int) {
	this.key = key
	this.val = val
}

func (this *ImGuiStoragePair) SetFloat(key ImGuiID, val float) {
	this.key = key
	*(*float)(unsafe.Pointer(&this.val)) = val
}

// Helper: Key->Value storage
// Typically you don't have to worry about this since a storage is held within each Window.
// We use it to e.g. store collapse state for a tree (Int 0/1)
// This is optimized for efficient lookup (dichotomy into a contiguous buffer) and rare insertion (typically tied to user interactions aka max once a frame)
// You can use it as custom user storage for temporary values. Declare your own storage if, for example:
// - You want to manipulate the open/close state of a particular sub-tree in your interface (tree node uses Int 0/1 to store their state).
// - You want to store custom debug data easily without adding or editing structures in your code (probably not efficient, but convenient)
// Types are NOT stored, so it is up to you to make sure your Key don't collide with different types.
type ImGuiStorage struct {
	Data     []ImGuiStoragePair
	Pointers map[ImGuiID]interface{}
}

func (this *ImGuiStorage) Clear() {
	this.Data = this.Data[:0]
	this.Pointers = make(map[ImGuiID]interface{})
}

func (this *ImGuiStorage) GetInt(key ImGuiID, default_val int) int {
	panic("not implemented")
}

func (this *ImGuiStorage) SetInt(key ImGuiID, val int) {
	panic("not implemented")
}

func (this *ImGuiStorage) GetBool(key ImGuiID, default_val bool) bool {
	panic("not implemented")
}

func (this *ImGuiStorage) SetBool(key ImGuiID, val bool) {
	panic("not implemented")
}

func (this *ImGuiStorage) GetFloat(key ImGuiID, default_val float) float {
	panic("not implemented")
}

func (this *ImGuiStorage) SetFloat(key ImGuiID, val float) {
	panic("not implemented")
}

func (this *ImGuiStorage) GetInterface(key ImGuiID) interface{} {
	panic("not implemented")
}

func (this *ImGuiStorage) SetInterface(key ImGuiID, val interface{}) {
	panic("not implemented")
}

// - Get***Ref() functions finds pair, insert on demand if missing, return pointer. Useful if you intend to do Get+Set.
// - References are only valid until a new value is added to the storage. Calling a Set***() function or a Get***Ref() function invalidates the pointer.
// - A typical use case where this is convenient for quick hacking (e.g. add storage during a live Edit&Continue session if you can't modify existing struct)
//      float* pvar ImGui::GetFloatRef(key) = ImGui::SliderFloat("var", pvar, 100.0f) 0, some_var *pvar +=

func (this *ImGuiStorage) GetIntRef(key ImGuiID, default_val int) *int {
	panic("not implemented")
}

func (this *ImGuiStorage) GetBoolRef(key ImGuiID, default_val bool) bool {
	panic("not implemented")
}

func (this *ImGuiStorage) GetFloatRef(key ImGuiID, default_val float) float {
	panic("not implemented")
}

func (this *ImGuiStorage) GetInterfaceRef(key ImGuiID, default_val interface{}) *interface{} {
	panic("not implemented")
}

// Use on your own storage if you know only integer are being stored (open/close all tree nodes)
func (this *ImGuiStorage) SetAllInt(val int) { panic("not implemented") }

// For quicker full rebuild of a storage (instead of an incremental one), you may add all your contents and then sort once.
func (this *ImGuiStorage) BuildSortByKey() {}

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

const IM_COL32_R_SHIFT = 16
const IM_COL32_G_SHIFT = 8
const IM_COL32_B_SHIFT = 0
const IM_COL32_A_SHIFT = 24
const IM_COL32_A_MASK = 0xFF000000

func IM_COL32(R, G, B, A byte) ImU32 {
	return (((ImU32)(A) << IM_COL32_A_SHIFT) | ((ImU32)(B) << IM_COL32_B_SHIFT) | ((ImU32)(G) << IM_COL32_G_SHIFT) | ((ImU32)(R) << IM_COL32_R_SHIFT))
}

const IM_COL32_WHITE = 0xFFFFFFFF
const IM_COL32_BLACK = 0xFF000000
const IM_COL32_BLACK_TRANS = 0x00000000

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
	_VtxWritePtr    []ImDrawVert          // [Internal] point within VtxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_IdxWritePtr    ImDrawIdx             // [Internal] point within IdxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_ClipRectStack  []ImVec4              // [Internal]
	_TextureIdStack []ImTextureID         // [Internal]
	_Path           []ImVec2              // [Internal] current path building
	_CmdHeader      ImDrawCmdHeader       // [Internal] template of active commands. Fields should match those of CmdBuffer.back().
	_Splitter       ImDrawListSplitter    // [Internal] for channels api (note: prefer using your own persistent instance of ImDrawListSplitter!)
	_FringeScale    float                 // [Internal] anti-alias fringe is scaled by this value, this helps to keep things sharp while zooming at vertex buffer content

}

func (this *ImDrawList) PushClipRect(clip_rect_min, clip_rect_max ImVec2, intersect_with_current_clip_rect bool) {
	panic("not implemented")
}                                                             // Render-level scissoring. This is passed down to your render function but not used for CPU-side coarse clipping. Prefer using higher-level ImGui::PushClipRect() to affect logic (hit-testing and widget culling)
func (this *ImDrawList) PushClipRectFullScreen()              { panic("not implemented") }
func (this *ImDrawList) PopClipRect()                         { panic("not implemented") }
func (this *ImDrawList) PushTextureID(texture_id ImTextureID) { panic("not implemented") }
func (this *ImDrawList) PopTextureID()                        { panic("not implemented") }
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
	panic("not implemented")
}
func (this *ImDrawList) AddRect(cp_min ImVec2, p_max ImVec2, col ImU32, rounding float, flags ImDrawFlags, thickness float /*= 1.0f*/) {
	panic("not implemented")
} // a: upper-left, b: lower-right (== upper-left + size)
func (this *ImDrawList) AddRectFilled(cp_min ImVec2, p_max ImVec2, col ImU32, rounding float, flags ImDrawFlags) {
	panic("not implemented")
} // a: upper-left, b: lower-right (== upper-left + size)
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
func (this *ImDrawList) AddTriangleFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) AddCircle(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	panic("not implemented")
}
func (this *ImDrawList) AddCircleFilled(center ImVec2, radius float, col ImU32, num_segments int) {
	panic("not implemented")
}
func (this *ImDrawList) AddNgon(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	panic("not implemented")
}
func (this *ImDrawList) AddNgonFilled(center ImVec2, radius float, col ImU32, num_segments int) {
	panic("not implemented")
}
func (this *ImDrawList) AddText(pos ImVec2, col ImU32, text_begin string, text_end string) {
	panic("not implemented")
}
func (this *ImDrawList) AddTextV(font *ImFont, font_size float, pos ImVec2, col ImU32, text_begin, text_end string, wrap_width float, cpu_fine_clip_rect *ImVec4) {
	panic("not implemented")
}
func (this *ImDrawList) AddPolyline(points []ImVec2, num_points int, col ImU32, flags ImDrawFlags, thickness float) {
	panic("not implemented")
}
func (this *ImDrawList) AddConvexPolyFilled(points []ImVec2, num_points int, col ImU32) {
	panic("not implemented")
} // Note: Anti-aliased filling requires points to be in clockwise order.
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

func (this *ImDrawList) PathArcTo(center ImVec2, radius, a_min, a_max float, num_segments int) {
	panic("not implemented")
}
func (this *ImDrawList) PathArcToFast(center ImVec2, radius float, a_min_of_12, a_max_of_12 int) {
	panic("not implemented")
} // Use precomputed angles for a 12 steps circle
func (this *ImDrawList) PathBezierCubicCurveTo(p2 *ImVec2, p3 ImVec2, p4 ImVec2, num_segments int) {
	panic("not implemented")
} // Cubic Bezier (4 control points)
func (this *ImDrawList) PathBezierQuadraticCurveTo(p2 *ImVec2, p3 ImVec2, num_segments int) {
	panic("not implemented")
} // Quadratic Bezier (3 control points)
func (this *ImDrawList) PathRect(rect_min, rect_max *ImVec2, rounding float, flags ImDrawFlags) {
	panic("not implemented")
}

// Advanced
func AddCallback(callback ImDrawCallback, callback_data interface{}) { panic("not implemented") } // Your rendering function must check for 'UserCallback' in ImDrawCmd and call the function instead of rendering triangles.
func (this *ImDrawList) AddDrawCmd()                                 { panic("not implemented") } // This is useful if you need to forcefully create a new draw call (to allow for dependent rendering / blending). Otherwise primitives are merged into the same draw-call as much as possible
func (this *ImDrawList) CloneOutput() *ImDrawList                    { panic("not implemented") } // Create a clone of the CmdBuffer/IdxBuffer/VtxBuffer.

// Advanced: Channels
// - Use to split render into layers. By switching channels to can render out-of-order (e.g. submit FG primitives before BG primitives)
// - Use to minimize draw calls (e.g. if going back-and-forth between multiple clipping rectangles, prefer to append into separate channels then merge at the end)
// - FIXME-OBSOLETE: This API shouldn't have been in ImDrawList in the first place!
//   Prefer using your own persistent instance of ImDrawListSplitter as you can stack them.
//   Using the ImDrawList::ChannelsXXXX you cannot stack a split over another.
func (this *ImDrawList) ChannelsSplit(count int)  { this._Splitter.Split(this, count) }
func (this *ImDrawList) ChannelsMerge()           { this._Splitter.Merge(this) }
func (this *ImDrawList) ChannelsSetCurrent(n int) { this._Splitter.SetCurrentChannel(this, n) }

// Advanced: Primitives allocations
// - We render triangles (three vertices)
// - All primitives needs to be reserved via PrimReserve() beforehand.
func (this *ImDrawList) PrimReserve(idx_count, vtx_count int)           { panic("not implemented") }
func (this *ImDrawList) PrimUnreserve(idx_count, vtx_count int)         { panic("not implemented") }
func (this *ImDrawList) PrimRect(a, b *ImVec2, col ImU32)               { panic("not implemented") } // Axis aligned rectangle (composed of two triangles)
func (this *ImDrawList) PrimRectUV(a, b, uv_a, uv_b *ImVec2, col ImU32) { panic("not implemented") }
func (this *ImDrawList) PrimQuadUV(a, b, c, d *ImVec2, uv_a, uv_b, yv_c, uv_d *ImVec2, col ImU32) {
	panic("not implemented")
}
func (this *ImDrawList) PrimWriteVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	this._VtxWritePtr[0].pos = pos
	this._VtxWritePtr[0].uv = *uv
	this._VtxWritePtr[0].col = col
	this._VtxWritePtr = this._VtxWritePtr[1:]
	this._VtxCurrentIdx++
}
func (this *ImDrawList) PrimWriteIdx(idx ImDrawIdx) {
	this._IdxWritePtr = idx
	this._IdxWritePtr++
}

func (this *ImDrawList) PrimVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx))
	this.PrimWriteVtx(pos, uv, col)
} // Write vertex with unique index

// [Internal helpers]
func (this *ImDrawList) _ResetForNewFrame()                       { panic("not implemented") }
func (this *ImDrawList) _ClearFreeMemory()                        { panic("not implemented") }
func (this *ImDrawList) _PopUnusedDrawCmd()                       { panic("not implemented") }
func (this *ImDrawList) _TryMergeDrawCmds()                       { panic("not implemented") }
func (this *ImDrawList) _OnChangedClipRect()                      { panic("not implemented") }
func (this *ImDrawList) _OnChangedTextureID()                     { panic("not implemented") }
func (this *ImDrawList) _OnChangedVtxOffset()                     { panic("not implemented") }
func (this *ImDrawList) _CalcCircleAutoSegmentCount(radius float) { panic("not implemented") }
func (this *ImDrawList) _PathArcToFastEx(center ImVec2, radius, a_min_sample, a_max_sample float, a_step int) {
	panic("not implemented")
}
func (this *ImDrawList) _PathArcToN(center ImVec2, radius, a_min, a_max float, num_segments int) {
	panic("not implemented")
}

// All draw data to render a Dear ImGui frame
// (NB: the style and the naming convention here is a little inconsistent, we currently preserve them for backward compatibility purpose,
// as this is one of the oldest structure exposed by the library! Basically, ImDrawList == CmdList)
type ImDrawData struct {
	Valid            bool          // Only valid after Render() is called and before the next NewFrame() is called.
	CmdListsCount    int           // Number of ImDrawList* to render
	TotalIdxCount    int           // For convenience, sum of all ImDrawList's IdxBuffer.Size
	TotalVtxCount    int           // For convenience, sum of all ImDrawList's VtxBuffer.Size
	CmdLists         []*ImDrawList // Array of ImDrawList* to render. The ImDrawList are owned by ImGuiContext and only pointed to from here.
	DisplayPos       ImVec2        // Top-left position of the viewport to render (== top-left of the orthogonal projection matrix to use) (== GetMainViewport()->Pos for the main viewport, == (0.0) in most single-viewport applications)
	DisplaySize      ImVec2        // Size of the viewport to render (== GetMainViewport()->Size for the main viewport, == io.DisplaySize in most single-viewport applications)
	FramebufferScale ImVec2        // Amount of pixels for each unit of DisplaySize. Based on io.DisplayFramebufferScale. Generally (1,1) on normal display, (2,2) on OSX with Retina display.
}

// Functions
func (ImDrawData) DeIndexAllBuffers()              { panic("not implemented") } // Helper to convert all buffers from indexed to non-indexed, in case you cannot render indexed. Note: this is slow and most likely a waste of resources. Always prefer indexed rendering!
func (ImDrawData) ScaleClipRects(fb_scale *ImVec2) { panic("not implemented") } // Helper to scale the ClipRect field of each ImDrawCmd. Use if your final output buffer is at a different scale than Dear ImGui expects, or if there is a difference between your window resolution and framebuffer resolution.

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
	TexPixelsRGBA32    []int                                       // 4 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight * 4
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

	if new_font_cfg.DstFont.EllipsisChar == math.MaxInt16 {
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
	FallbackGlyph []ImFontGlyph // 4-8   // out // = FindGlyph(FontFallbackChar)

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
		FallbackChar: math.MaxInt16,
		EllipsisChar: math.MaxInt16,
		DotChar:      math.MaxInt16,
		Scale:        1,
	}
}

func (this *ImFont) FindGlyph(c ImWchar) *ImFontGlyph           { panic("not implemented") }
func (this *ImFont) FindGlyphNoFallback(c ImWchar) *ImFontGlyph { panic("not implemented") }

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

// 'max_width' stops rendering after a certain width (could be turned into a 2d size). FLT_MAX to disable.
// 'wrap_width' enable automatic word-wrapping across multiple lines to fit into given width. 0.0f to disable.
func (this *ImFont) CalcTextSizeA(size, max_width, wrap_width float, text_begin, text_end string, remaining *string) ImVec2 {
	panic("not implemented")
} // utf8
func (this *ImFont) CalcWordWrapPositionA(scale float, text, text_end string, wrap_width float) string {
	panic("not implemented")
}
func (this *ImFont) RenderChar(draw_list *ImDrawList, size float, pos ImVec2, col ImU32, c ImWchar) {
	panic("not implemented")
}
func (this *ImFont) RenderText(draw_list *ImDrawList, size float, pos ImVec2, col ImU32, clip_rect *ImVec4, text_begin, text_end string, wrap_width float, cpu_fine_clip bool) {
	panic("not implemented")
}

// [Internal] Don't use!
func (this *ImFont) BuildLookupTable()      { panic("not implemented") }
func (this *ImFont) ClearOutputData()       { panic("not implemented") }
func (this *ImFont) GrowIndex(new_size int) { panic("not implemented") }
func (this *ImFont) AddGlyph(src_cfg *ImFontConfig, c ImWchar, x0, y0, x1, y1, u0, v0, u1, v1, advance_x float) {
	panic("not implemented")
}
func (this *ImFont) AddRemapChar(dst, src ImWchar, overwite_dst bool /*= true*/) {
	panic("not implemented")
}                                                                 // Makes 'dst' character/glyph points to 'src' character/glyph. Currently needs to be called AFTER fonts have been built.
func (this *ImFont) SetGlyphVisible(c ImWchar, visible bool)      { panic("not implemented") }
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

}

// Helpers
func (this *ImGuiViewport) GetCenter() ImVec2 {
	return ImVec2{this.Pos.x + this.Size.x*0.5, this.Pos.y + this.Size.y*0.5}
}
func (this *ImGuiViewport) GetWorkCenter() ImVec2 {
	return ImVec2{this.WorkPos.x + this.WorkSize.x*0.5, this.WorkPos.y + this.WorkSize.y*0.5}
}
