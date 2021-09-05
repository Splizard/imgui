package imgui

const (
	KeyTab = iota
	KeyLeftArrow
	KeyRightArrow
	KeyUpArrow
	KeyDownArrow
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyInsert
	KeyDelete
	KeyBackspace
	KeySpace
	KeyEnter
	KeyEscape
	KeyKeyPadEnter
	KeyA // for text edit CTRL+A: select all
	KeyC // for text edit CTRL+C: copy
	KeyV // for text edit CTRL+V: paste
	KeyX // for text edit CTRL+X: cut
	KeyY // for text edit CTRL+Y: redo
	KeyZ // for text edit CTRL+Z: undo
	keyCOUNT
)

const (
	// Gamepad Mapping
	KeyNavInputActivate    = iota // activate / open / toggle / tweak value       // e.g. Cross  (PS4), A (Xbox), A (Switch), Space (Keyboard)
	KeyNavInputCancel             // cancel / close / exit                        // e.g. Circle (PS4), B (Xbox), B (Switch), Escape (Keyboard)
	KeyNavInputInput              // text input / on-screen keyboard              // e.g. Triang.(PS4), Y (Xbox), X (Switch), Return (Keyboard)
	KeyNavInputMenu               // tap: toggle menu / hold: focus, move, resize // e.g. Square (PS4), X (Xbox), Y (Switch), Alt (Keyboard)
	KeyNavInputDpadLeft           // move / tweak / resize window (w/ PadMenu)    // e.g. D-pad Left/Right/Up/Down (Gamepads), Arrow keys (Keyboard)
	KeyNavInputDpadRight          //
	KeyNavInputDpadUp             //
	KeyNavInputDpadDown           //
	KeyNavInputLStickLeft         // scroll / move window (w/ PadMenu)            // e.g. Left Analog Stick Left/Right/Up/Down
	KeyNavInputLStickRight        //
	KeyNavInputLStickUp           //
	KeyNavInputLStickDown         //
	KeyNavInputFocusPrev          // next window (w/ PadMenu)                     // e.g. L1 or L2 (PS4), LB or LT (Xbox), L or ZL (Switch)
	KeyNavInputFocusNext          // prev window (w/ PadMenu)                     // e.g. R1 or R2 (PS4), RB or RT (Xbox), R or ZL (Switch)
	KeyNavInputTweakSlow          // slower tweaks                                // e.g. L1 or L2 (PS4), LB or LT (Xbox), L or ZL (Switch)
	KeyNavInputTweakFast          // faster tweaks                                // e.g. R1 or R2 (PS4), RB or RT (Xbox), R or ZL (Switch)

	// [Internal] Don't use directly! This is used internally to differentiate keyboard from gamepad inputs for behaviors that require to differentiate them.
	// Keyboard behavior that have no corresponding gamepad mapping (e.g. CTRL+TAB) will be directly reading from io.KeysDown[] instead of io.NavInputs[].
	keyNavInputKeyLeft  // move left                                    // = Arrow keys
	keyNavInputKeyRight // move right
	keyNavInputKeyUp    // move up
	keyNavInputKeyDown  // move down
	keyNavInputCOUNT
	keyNavInputInternalStart = keyNavInputKeyLeft
)

type IO struct {
	//------------------------------------------------------------------
	// Configuration (fill once)                // Default value
	//------------------------------------------------------------------

	ConfigFlags             ConfigFlags   // = 0              // See ImGuiConfigFlags_ enum. Set by user/application. Gamepad/keyboard navigation options, etc.
	BackendFlags            BackendFlags  // = 0              // See ImGuiBackendFlags_ enum. Set by backend (imgui_impl_xxx files or custom backend) to communicate features supported by the backend.
	DisplaySize             Vec2          // <unset>          // Main display size, in pixels (generally == GetMainViewport()->Size)
	DeltaTime               float32       // = 1.0f/60.0f     // Time elapsed since last frame, in seconds.
	IniSavingRate           float32       // = 5.0f           // Minimum time between saving positions/sizes to .ini file, in seconds.
	IniFilename             string        // = "imgui.ini"    // Path to .ini file (important: default "imgui.ini" is relative to current working dir!). Set NULL to disable automatic .ini loading/saving or if you want to manually call LoadIniSettingsXXX() / SaveIniSettingsXXX() functions.
	LogFilename             string        // = "imgui_log.txt"// Path to .log file (default parameter to ImGui::LogToFile when no file is specified).
	MouseDoubleClickTime    float32       // = 0.30f          // Time for a double-click, in seconds.
	MouseDoubleClickMaxDist float32       // = 6.0f           // Distance threshold to stay in to validate a double-click, in pixels.
	MouseDragThreshold      float32       // = 6.0f           // Distance threshold before considering we are dragging.
	KeyMap                  [keyCOUNT]int // <unset>          // Map of indices into the KeysDown[512] entries array which represent your "native" keyboard state.
	KeyRepeatDelay          float32       // = 0.250f         // When holding a key/button, time before it starts repeating, in seconds (for buttons in Repeat mode, etc.).
	KeyRepeatRate           float32       // = 0.050f         // When holding a key/button, rate at which it repeats, in seconds.
	UserData                interface{}   // = NULL           // Store your own data for retrieval by callbacks.

	Fonts                   *FontAtlas // <auto>           // Font atlas: load, rasterize and pack one or more fonts into a single texture.
	FontGlobalScale         float32    // = 1.0f           // Global scale all fonts
	FontAllowUserScaling    bool       // = false          // Allow user scaling text of individual window with CTRL+Wheel.
	FontDefault             *Font      // = NULL           // Font to use on NewFrame(). Use NULL to uses Fonts->Fonts[0].
	DisplayFramebufferScale Vec2       // = (1, 1)         // For retina display or other situations where window coordinates are different from framebuffer coordinates. This generally ends up in ImDrawData::FramebufferScale.

	// Miscellaneous options
	MouseDrawCursor                   bool    // = false          // Request ImGui to draw a mouse cursor for you (if you are on a platform without a mouse cursor). Cannot be easily renamed to 'io.ConfigXXX' because this is frequently used by backend implementations.
	ConfigMacOSXBehaviors             bool    // = defined(__APPLE__) // OS X style: Text editing cursor movement using Alt instead of Ctrl, Shortcuts using Cmd/Super instead of Ctrl, Line/Text Start and End using Cmd+Arrows instead of Home/End, Double click selects by word instead of selecting whole text, Multi-selection in lists uses Cmd/Super instead of Ctrl.
	ConfigInputTextCursorBlink        bool    // = true           // Enable blinking cursor (optional as some users consider it to be distracting).
	ConfigDragClickToInputText        bool    // = false          // [BETA] Enable turning DragXXX widgets into text input with a simple mouse click-release (without moving). Not desirable on devices without a keyboard.
	ConfigWindowsResizeFromEdges      bool    // = true           // Enable resizing of windows from their edges and from the lower-left corner. This requires (io.BackendFlags & ImGuiBackendFlags_HasMouseCursors) because it needs mouse cursor feedback. (This used to be a per-window ImGuiWindowFlags_ResizeFromAnySide flag)
	ConfigWindowsMoveFromTitleBarOnly bool    // = false       // Enable allowing to move windows only when clicking on their title bar. Does not apply to windows without a title bar.
	ConfigMemoryCompactTimer          float32 // = 60.0f          // Timer (in seconds) to free transient windows/tables memory buffers when unused. Set to -1.0f to disable.

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
	GetClipboardTextFn func(interface{}) string
	SetClipboardTextFn func(interface{}, string)
	ClipboardUserData  interface{}

	// Optional: Notify OS Input Method Editor of the screen position of your cursor for text input position (e.g. when using Japanese/Chinese IME on Windows)
	// (default to use native imm32 api on Windows)
	ImeSetInputScreenPosFn func(x, y int)
	ImeWindowHandle        interface{} // = NULL           // (Windows) Set this to your HWND to get automatic IME cursor positioning.

	//------------------------------------------------------------------
	// Input - Fill before calling NewFrame()
	//------------------------------------------------------------------

	MousePos    Vec2                      // Mouse position, in pixels. Set to ImVec2(-FLT_MAX, -FLT_MAX) if mouse is unavailable (on another screen, etc.)
	MouseDown   [5]bool                   // Mouse buttons: 0=left, 1=right, 2=middle + extras (ImGuiMouseButton_COUNT == 5). Dear ImGui mostly uses left and right buttons. Others buttons allows us to track if the mouse is being used by your application + available to user as a convenience via IsMouse** API.
	MouseWheel  float32                   // Mouse wheel Vertical: 1 unit scrolls about 5 lines text.
	MouseWheelH float32                   // Mouse wheel Horizontal. Most users don't have a mouse with an horizontal wheel, may not be filled by all backends.
	KeyCtrl     bool                      // Keyboard modifier pressed: Control
	KeyShift    bool                      // Keyboard modifier pressed: Shift
	KeyAlt      bool                      // Keyboard modifier pressed: Alt
	KeySuper    bool                      // Keyboard modifier pressed: Cmd/Super/Windows
	KeysDown    [512]bool                 // Keyboard keys that are pressed (ideally left in the "native" order your engine has access to keyboard keys, so you can use your own defines/enums for keys).
	NavInputs   [keyNavInputCOUNT]float32 // Gamepad inputs. Cleared back to zero by EndFrame(). Keyboard keys will be auto-mapped and be written here by NewFrame().

	//------------------------------------------------------------------
	// Output - Updated by NewFrame() or EndFrame()/Render()
	// (when reading from the io.WantCaptureMouse, io.WantCaptureKeyboard flags to dispatch your inputs, it is
	//  generally easier and more correct to use their state BEFORE calling NewFrame(). See FAQ for details!)
	//------------------------------------------------------------------

	WantCaptureMouse         bool    // Set when Dear ImGui will use mouse inputs, in this case do not dispatch them to your main game/application (either way, always pass on mouse inputs to imgui). (e.g. unclicked mouse is hovering over an imgui window, widget is active, mouse was clicked over an imgui window, etc.).
	WantCaptureKeyboard      bool    // Set when Dear ImGui will use keyboard inputs, in this case do not dispatch them to your main game/application (either way, always pass keyboard inputs to imgui). (e.g. InputText active, or an imgui window is focused and navigation is enabled, etc.).
	WantTextInput            bool    // Mobile/console: when set, you may display an on-screen keyboard. This is set by Dear ImGui when it wants textual keyboard input to happen (e.g. when a InputText widget is active).
	WantSetMousePos          bool    // MousePos has been altered, backend should reposition mouse on next frame. Rarely used! Set only when ImGuiConfigFlags_NavEnableSetMousePos flag is enabled.
	WantSaveIniSettings      bool    // When manual .ini load/save is active (io.IniFilename == NULL), this will be set to notify your application that you can call SaveIniSettingsToMemory() and save yourself. Important: clear io.WantSaveIniSettings yourself after saving!
	NavActive                bool    // Keyboard/Gamepad navigation is currently allowed (will handle KeyNavXXX events) = a window is focused and it doesn't use the ImGuiWindowFlags_NoNavInputs flag.
	NavVisible               bool    // Keyboard/Gamepad navigation is visible and allowed (will handle KeyNavXXX events).
	Framerate                float32 // Rough estimate of application framerate, in frame per second. Solely for convenience. Rolling average estimation based on io.DeltaTime over 120 frames.
	MetricsRenderVertices    int     // Vertices output during last call to Render()
	MetricsRenderIndices     int     // Indices output during last call to Render() = number of triangles * 3
	MetricsRenderWindows     int     // Number of visible windows
	MetricsActiveWindows     int     // Number of active windows
	MetricsActiveAllocations int     // Number of active allocations, updated by MemAlloc/MemFree based on current context. May be off if you have multiple imgui contexts.
	MouseDelta               Vec2    // Mouse delta. Note that this is zero if either current or previous position are invalid (-FLT_MAX,-FLT_MAX), so a disappearing/reappearing mouse won't have a huge delta.

	//------------------------------------------------------------------
	// [Internal] Dear ImGui will maintain those fields. Forward compatibility not guaranteed!
	//------------------------------------------------------------------

	wantCaptureMouseUnlessPopupClose bool         // Alternative to WantCaptureMouse: (WantCaptureMouse == true && WantCaptureMouseUnlessPopupClose == false) when a click over void is expected to close a popup.
	keyMods                          KeyModFlags  // Key mods flags (same as io.KeyCtrl/KeyShift/KeyAlt/KeySuper but merged into flags), updated by NewFrame()
	keyModsPrev                      KeyModFlags  // Previous key mods
	mousePosPrev                     Vec2         // Previous mouse position (note that MouseDelta is not necessary == MousePos-MousePosPrev, in case either position is invalid)
	mouseClickedPos                  [5]Vec2      // Position at time of clicking
	mouseClickedTime                 [5]float64   // Time of last click (used to figure out double-click)
	mouseClicked                     [5]bool      // Mouse button went from !Down to Down
	mouseDoubleClicked               [5]bool      // Has mouse button been double-clicked?
	mouseReleased                    [5]bool      // Mouse button went from Down to !Down
	mouseDownOwned                   [5]bool      // Track if button was clicked inside a dear imgui window or over void blocked by a popup. We don't request mouse capture from the application if click started outside ImGui bounds.
	mouseDownOwnedUnlessPopupClose   [5]bool      //Track if button was clicked inside a dear imgui window.
	mouseDownWasDoubleClick          [5]bool      // Track if button down was a double-click
	mouseDownDuration                [5]float32   // Duration the mouse button has been down (0.0f == just clicked)
	mouseDownDurationPrev            [5]float32   // Previous time the mouse button has been down
	mouseDragMaxDistanceAbs          [5]Vec2      // Maximum distance, absolute, on each axis, of how much mouse has traveled from the clicking point
	mouseDragMaxDistanceSqr          [5]float32   // Squared maximum distance of how much mouse has traveled from the clicking point
	keysDownDuration                 [512]float32 // Duration the keyboard key has been down (0.0f == just pressed)
	keysDownDurationPrev             [512]float32 // Previous duration the key has been down
	navInputsDownDuration            [keyNavInputCOUNT]float32
	navInputsDownDurationPrev        [keyNavInputCOUNT]float32
	penPressure                      float32 // Touch/Pen pressure (0.0f to 1.0f, should be >0.0f only when MouseDown[0] == true). Helper storage currently unused by Dear ImGui.
	inputQueueCharacters             string  // Queue of _characters_ input (obtained by platform backend). Fill using AddInputCharacter() helper.
}

func NewIO() {
	panic("not implemented")
}

func (io *IO) AddInputCharacter(r rune)    { panic("not implemented") }
func (io *IO) AddInputCharacters(s string) { panic("not implemented") }
func (io *IO) ClearInputCharacters()       { panic("not implemented") }
func (io *IO) AddFocusEvent(focused bool)  { panic("not implemented") }
