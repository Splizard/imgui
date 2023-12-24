package imgui

import "bytes"

type ImGuiContext struct {
	Initialized                        bool
	IO                                 ImGuiIO
	Style                              ImGuiStyle
	Font                               *ImFont
	FontSize                           float
	FontBaseSize                       float
	DrawListSharedData                 ImDrawListSharedData
	Time                               double
	FrameCount                         int
	FrameCountEnded                    int
	FrameCountRendered                 int
	WithinFrameScope                   bool // Set by NewFrame(), cleared by EndFrame()
	WithinFrameScopeWithImplicitWindow bool // Set by NewFrame(), cleared by EndFrame() when the implicit debug window has been pushed
	WithinEndChild                     bool // Set within EndChild()
	GcCompactAll                       bool

	// Windows state
	Windows                        []*ImGuiWindow // Windows, sorted in display order, back to front
	WindowsFocusOrder              []*ImGuiWindow // Root windows, sorted in focus order, back to front.
	WindowsTempSortBuffer          []*ImGuiWindow // Temporary buffer used in EndFrame() to reorder windows so parents are kept before their child
	CurrentWindowStack             []ImGuiWindowStackData
	WindowsById                    ImGuiStorage // Map window's ImGuiID to *ImGuiWindow
	WindowsActiveCount             int          // Number of unique windows submitted by frame
	WindowsHoverPadding            ImVec2       // Padding around resizable windows for which hovering on counts as hovering the window == max(style.TouchExtraPadding, WINDOWS_HOVER_PADDING)
	CurrentWindow                  *ImGuiWindow // Window being drawn into
	HoveredWindow                  *ImGuiWindow // Window the mouse is hovering. Will typically catch mouse inputs.
	HoveredWindowUnderMovingWindow *ImGuiWindow // Hovered window ignoring MovingWindow. Only set if MovingWindow is set.
	MovingWindow                   *ImGuiWindow // Track the window we clicked on (in order to preserve focus). The actual window that is moved is generally MovingWindow->RootWindow.
	WheelingWindow                 *ImGuiWindow // Track the window we started mouse-wheeling on. Until a timer elapse or mouse has moved, generally keep scrolling the same window even if during the course of scrolling the mouse ends up hovering a child window.
	WheelingWindowRefMousePos      ImVec2
	WheelingWindowTimer            float

	// Item/widgets state and tracking information
	HoveredId                                ImGuiID // Hovered widget, filled during the frame
	HoveredIdPreviousFrame                   ImGuiID
	HoveredIdAllowOverlap                    bool
	HoveredIdUsingMouseWheel                 bool // Hovered widget will use mouse wheel. Blocks scrolling the underlying window.
	HoveredIdPreviousFrameUsingMouseWheel    bool
	HoveredIdDisabled                        bool    // At least one widget passed the rect test, but has been discarded by disabled flag or popup inhibit. May be true even if HoveredId == 0.
	HoveredIdTimer                           float   // Measure contiguous hovering time
	HoveredIdNotActiveTimer                  float   // Measure contiguous hovering time where the item has not been active
	ActiveId                                 ImGuiID // Active widget
	ActiveIdIsAlive                          ImGuiID // Active widget has been seen this frame (we can't use a bool as the ActiveId may change within the frame)
	ActiveIdTimer                            float
	ActiveIdIsJustActivated                  bool // Set at the time of activation for one frame
	ActiveIdAllowOverlap                     bool // Active widget allows another widget to steal active id (generally for overlapping widgets, but not always)
	ActiveIdNoClearOnFocusLoss               bool // Disable losing active id if the active id window gets unfocused.
	ActiveIdHasBeenPressedBefore             bool // Track whether the active id led to a press (this is to allow changing between PressOnClick and PressOnRelease without pressing twice). Used by range_select branch.
	ActiveIdHasBeenEditedBefore              bool // Was the value associated to the widget Edited over the course of the Active state.
	ActiveIdHasBeenEditedThisFrame           bool
	ActiveIdUsingMouseWheel                  bool   // Active widget will want to read mouse wheel. Blocks scrolling the underlying window.
	ActiveIdUsingNavDirMask                  ImU32  // Active widget will want to read those nav move requests (e.g. can activate a button and move away from it)
	ActiveIdUsingNavInputMask                ImU32  // Active widget will want to read those nav inputs.
	ActiveIdUsingKeyInputMask                ImU64  // Active widget will want to read those key inputs. When we grow the ImGuiKey enum we'll need to either to order the enum to make useful keys come first, either redesign this into e.g. a small array.
	ActiveIdClickOffset                      ImVec2 // Clicked offset from upper-left corner, if applicable (currently only set by ButtonBehavior)
	ActiveIdWindow                           *ImGuiWindow
	ActiveIdSource                           ImGuiInputSource // Activating with mouse or nav (gamepad/keyboard)
	ActiveIdMouseButton                      ImGuiMouseButton
	ActiveIdPreviousFrame                    ImGuiID
	ActiveIdPreviousFrameIsAlive             bool
	ActiveIdPreviousFrameHasBeenEditedBefore bool
	ActiveIdPreviousFrameWindow              *ImGuiWindow
	LastActiveId                             ImGuiID // Store the last non-zero ActiveId, useful for animation.
	LastActiveIdTimer                        float   // Store the last non-zero ActiveId timer since the beginning of activation, useful for animation.

	// Next window/item data
	CurrentItemFlags ImGuiItemFlags      // == g.ItemFlagsStack.back()
	NextItemData     ImGuiNextItemData   // Storage for SetNextItem** functions
	LastItemData     ImGuiLastItemData   // Storage for last submitted item (setup by ItemAdd)
	NextWindowData   ImGuiNextWindowData // Storage for SetNextWindow** functions

	// Shared stacks
	ColorStack      []ImGuiColorMod  // Stack for PushStyleColor()/PopStyleColor() - inherited by Begin()
	StyleVarStack   []ImGuiStyleMod  // Stack for PushStyleVar()/PopStyleVar() - inherited by Begin()
	FontStack       []*ImFont        // Stack for PushFont()/PopFont() - inherited by Begin()
	FocusScopeStack []ImGuiID        // Stack for PushFocusScope()/PopFocusScope() - not inherited by Begin(), unless child window
	ItemFlagsStack  []ImGuiItemFlags // Stack for PushItemFlag()/PopItemFlag() - inherited by Begin()
	GroupStack      []ImGuiGroupData // Stack for BeginGroup()/EndGroup() - not inherited by Begin()
	OpenPopupStack  []ImGuiPopupData // Which popups are open (persistent)
	BeginPopupStack []ImGuiPopupData // Which level of BeginPopup() we are in (reset every frame)

	// Viewports
	Viewports []*ImGuiViewportP // Active viewports (Size==1 in 'master' branch). Each viewports hold their copy of ImDrawData.

	// Gamepad/keyboard Navigation
	NavWindow                  *ImGuiWindow // Focused window for navigation. Could be called 'FocusWindow'
	NavId                      ImGuiID      // Focused item for navigation
	NavFocusScopeId            ImGuiID      // Identify a selection scope (selection code often wants to "clear other items" when landing on an item of the selection set)
	NavActivateId              ImGuiID      // ~~ (g.ActiveId == 0) && IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0, also set when calling ActivateItem()
	NavActivateDownId          ImGuiID      // ~~ IsNavInputDown(ImGuiNavInput_Activate) ? NavId : 0
	NavActivatePressedId       ImGuiID      // ~~ IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0
	NavInputId                 ImGuiID      // ~~ IsNavInputPressed(ImGuiNavInput_Input) ? NavId : 0
	NavJustTabbedId            ImGuiID      // Just tabbed to this id.
	NavJustMovedToId           ImGuiID      // Just navigated to this id (result of a successfully MoveRequest).
	NavJustMovedToFocusScopeId ImGuiID      // Just navigated to this focus scope id (result of a successfully MoveRequest).
	NavJustMovedToKeyMods      ImGuiKeyModFlags
	NavNextActivateId          ImGuiID          // Set by ActivateItem(), queued until next frame.
	NavInputSource             ImGuiInputSource // Keyboard or Gamepad mode? THIS WILL ONLY BE None or NavGamepad or NavKeyboard.
	NavLayer                   ImGuiNavLayer    // Layer we are navigating on. For now the system is hard-coded for 0=main contents and 1=menu/title bar, may expose layers later.
	NavIdTabCounter            int              // == NavWindow->DC.FocusIdxTabCounter at time of NavId processing
	NavIdIsAlive               bool             // Nav widget has been seen this frame ~~ NavRectRel is valid
	NavMousePosDirty           bool             // When set we will update mouse position if (io.ConfigFlags & ImGuiConfigFlags_NavEnableSetMousePos) if set (NB: this not enabled by default)
	NavDisableHighlight        bool             // When user starts using mouse, we hide gamepad/keyboard highlight (NB: but they are still available, which is why NavDisableHighlight isn't always != NavDisableMouseHover)
	NavDisableMouseHover       bool             // When user starts using gamepad/keyboard, we hide mouse hovering highlight until mouse is touched again.

	// Navigation: Init & Move Requests
	NavAnyRequest             bool // ~~ NavMoveRequest || NavInitRequest this is to perform early out in ItemAdd()
	NavInitRequest            bool // Init request for appearing window to select first item
	NavInitRequestFromMove    bool
	NavInitResultId           ImGuiID // Init request result (first item of the window, or one for which SetItemDefaultFocus() was called)
	NavInitResultRectRel      ImRect  // Init request result rectangle (relative to parent window)
	NavMoveSubmitted          bool    // Move request submitted, will process result on next NewFrame()
	NavMoveScoringItems       bool    // Move request submitted, still scoring incoming items
	NavMoveForwardToNextFrame bool
	NavMoveFlags              ImGuiNavMoveFlags
	NavMoveKeyMods            ImGuiKeyModFlags
	NavMoveDir                ImGuiDir // Direction of the move request (left/right/up/down)
	NavMoveDirForDebug        ImGuiDir
	NavMoveClipDir            ImGuiDir         // FIXME-NAV: Describe the purpose of this better. Might want to rename?
	NavScoringRect            ImRect           // Rectangle used for scoring, in screen space. Based of window->NavRectRel[], modified for directional navigation scoring.
	NavScoringDebugCount      int              // Metrics for debugging
	NavMoveResultLocal        ImGuiNavItemData // Best move request candidate within NavWindow
	NavMoveResultLocalVisible ImGuiNavItemData // Best move request candidate within NavWindow that are mostly visible (when using ImGuiNavMoveFlags_AlsoScoreVisibleSet flag)
	NavMoveResultOther        ImGuiNavItemData // Best move request candidate within NavWindow's flattened hierarchy (when using ImGuiWindowFlags_NavFlattened flag)

	// Navigation: Windowing (CTRL+TAB for list, or Menu button + keys or directional pads to move/resize)
	NavWindowingTarget         *ImGuiWindow // Target window when doing CTRL+Tab (or Pad Menu + FocusPrev/Next), this window is temporarily displayed top-most!
	NavWindowingTargetAnim     *ImGuiWindow // Record of last valid NavWindowingTarget until DimBgRatio and NavWindowingHighlightAlpha becomes 0.0f, so the fade-out can stay on it.
	NavWindowingListWindow     *ImGuiWindow // Internal window actually listing the CTRL+Tab contents
	NavWindowingTimer          float
	NavWindowingHighlightAlpha float
	NavWindowingToggleLayer    bool

	// Legacy Focus/Tabbing system (older than Nav, active even if Nav is disabled, misnamed. FIXME-NAV: This needs a redesign!)
	TabFocusRequestCurrWindow         *ImGuiWindow //
	TabFocusRequestNextWindow         *ImGuiWindow //
	TabFocusRequestCurrCounterRegular int          // Any item being requested for focus, stored as an index (we on layout to be stable between the frame pressing TAB and the next frame, semi-ouch)
	TabFocusRequestCurrCounterTabStop int          // Tab item being requested for focus, stored as an index
	TabFocusRequestNextCounterRegular int          // Stored for next frame
	TabFocusRequestNextCounterTabStop int          // "
	TabFocusPressed                   bool         // Set in NewFrame() when user pressed Tab

	// Render
	DimBgRatio  float // 0.0..1.0 animation when fading in a dimming background (for modal window and CTRL+TAB list)
	MouseCursor ImGuiMouseCursor

	// Drag and Drop
	DragDropActive                  bool
	DragDropWithinSource            bool // Set when within a BeginDragDropXXX/EndDragDropXXX block for a drag source.
	DragDropWithinTarget            bool // Set when within a BeginDragDropXXX/EndDragDropXXX block for a drag target.
	DragDropSourceFlags             ImGuiDragDropFlags
	DragDropSourceFrameCount        int
	DragDropMouseButton             ImGuiMouseButton
	DragDropPayload                 ImGuiPayload
	DragDropTargetRect              ImRect // Store rectangle of current target candidate (we favor small targets when overlapping)
	DragDropTargetId                ImGuiID
	DragDropAcceptFlags             ImGuiDragDropFlags
	DragDropAcceptIdCurrRectSurface float    // Target item surface (we resolve overlapping targets by prioritizing the smaller surface)
	DragDropAcceptIdCurr            ImGuiID  // Target item id (set at the time of accepting the payload)
	DragDropAcceptIdPrev            ImGuiID  // Target item id from previous frame (we need to store this to allow for overlapping drag and drop targets)
	DragDropAcceptFrameCount        int      // Last time a target expressed a desire to accept the source
	DragDropHoldJustPressedId       ImGuiID  // Set when holding a payload just made ButtonBehavior() return a press.
	DragDropPayloadBufHeap          []byte   // We don't expose the ImVector<> directly, ImGuiPayload only holds pointer+size
	DragDropPayloadBufLocal         [16]byte // Local buffer for small payloads

	// Table
	CurrentTable                *ImGuiTable
	CurrentTableStackIdx        int
	Tables                      map[ImGuiID]*ImGuiTable
	TablesTempDataStack         []ImGuiTableTempData
	TablesLastTimeActive        map[int]float // Last used timestamp of each tables (SOA, for efficient GC)
	DrawChannelsTempMergeBuffer []ImDrawChannel

	// Tab bars
	CurrentTabBar      *ImGuiTabBar
	TabBars            map[ImGuiID]*ImGuiTabBar
	CurrentTabBarStack []ImGuiPtrOrIndex
	ShrinkWidthBuffer  []ImGuiShrinkWidthItem

	// Widget state
	MouseLastValidPos               ImVec2
	InputTextState                  ImGuiInputTextState
	InputTextPasswordFont           ImFont
	TempInputId                     ImGuiID             // Temporary text input when CTRL+clicking on a slider, etc.
	ColorEditOptions                ImGuiColorEditFlags // Store user options for color edit widgets
	ColorEditLastHue                float               // Backup of last Hue associated to LastColor[3], so we can restore Hue in lossy RGB<>HSV round trips
	ColorEditLastSat                float               // Backup of last Saturation associated to LastColor[3], so we can restore Saturation in lossy RGB<>HSV round trips
	ColorEditLastColor              [3]float
	ColorPickerRef                  ImVec4 // Initial/reference color at the time of opening the color picker.
	ComboPreviewData                ImGuiComboPreviewData
	SliderCurrentAccum              float // Accumulated slider delta when using navigation controls.
	SliderCurrentAccumDirty         bool  // Has the accumulated slider delta changed since last time we tried to apply it?
	DragCurrentAccumDirty           bool
	DragCurrentAccum                float // Accumulator for dragging modification. Always high-precision, not rounded by end-user precision settings
	DragSpeedDefaultRatio           float // If speed == 0.0f, uses (max-min) * DragSpeedDefaultRatio
	DisabledAlphaBackup             float // Backup for style.Alpha for BeginDisabled()
	ScrollbarClickDeltaToGrabCenter float // Distance between mouse and center of grab box, normalized in parent space. Use storage?
	TooltipOverrideCount            int
	TooltipSlowDelay                float     // Time before slow tooltips appears (FIXME: This is temporary until we merge in tooltip timer+priority work)
	ClipboardHandlerData            []char    // If no custom clipboard handler is defined
	MenusIdSubmittedThisFrame       []ImGuiID // A list of menu IDs that were rendered at least once

	// Platform support
	PlatformImePos             ImVec2 // Cursor position request & last passed to the OS Input Method Editor
	PlatformImeLastPos         ImVec2
	PlatformLocaleDecimalPoint char // '.' or *localeconv()->decimal_point

	// Settings
	SettingsLoaded     bool
	SettingsDirtyTimer float                  // Save .ini Settings to memory when time reaches zero
	SettingsIniData    ImGuiTextBuffer        // In memory .ini settings
	SettingsHandlers   []ImGuiSettingsHandler // List of .ini settings handlers
	SettingsWindows    []ImGuiWindowSettings  // ImGuiWindow .ini settings entries
	SettingsTables     []ImGuiTableSettings   // ImGuiTable .ini settings entries
	Hooks              []ImGuiContextHook     // Hooks for extensions (e.g. test engine)
	HookIdNext         ImGuiID                // Next available HookId

	// Capture/Logging
	LogEnabled              bool         // Currently capturing
	LogType                 ImGuiLogType // Capture target
	LogFile                 ImFileHandle // If != NULL log to stdout/ file
	LogBuffer               bytes.Buffer // Accumulation buffer when log to clipboard. This is pointer so our g static constructor doesn't call heap allocators.
	LogNextPrefix           string
	LogNextSuffix           string
	LogLinePosY             float
	LogLineFirstItem        bool
	LogDepthRef             int
	LogDepthToExpand        int
	LogDepthToExpandDefault int // Default/stored value for LogDepthMaxExpand if not specified in the LogXXX function call.

	// Debug Tools
	DebugItemPickerActive  bool    // Item picker is active (started with DebugStartItemPicker())
	DebugItemPickerBreakId ImGuiID // Will call IM_DEBUG_BREAK() when encountering this id
	DebugMetricsConfig     ImGuiMetricsConfig

	// Misc
	FramerateSecPerFrame         [120]float // Calculate estimate of framerate for user over the last 2 seconds.
	FramerateSecPerFrameIdx      int
	FramerateSecPerFrameCount    int
	FramerateSecPerFrameAccum    float
	WantCaptureMouseNextFrame    int // Explicit capture via CaptureKeyboardFromApp()/CaptureMouseFromApp() sets those flags
	WantCaptureKeyboardNextFrame int
	WantTextInputNextFrame       int
	TempBuffer                   string // Temporary text buffer

	FontAtlasOwnedByContext bool
}

func NewImGuiContext(atlas *ImFontAtlas) ImGuiContext {
	if atlas == nil {
		ptr := NewImFontAtlas()
		atlas = &ptr
	}
	var io = NewImGuiIO()
	io.Fonts = atlas
	return ImGuiContext{
		IO:                                io,
		DrawListSharedData:                NewImDrawListSharedData(),
		Style:                             NewImGuiStyle(),
		FrameCountEnded:                   -1,
		FrameCountRendered:                -1,
		ActiveIdClickOffset:               ImVec2{-1, -1},
		ActiveIdMouseButton:               -1,
		NavIdTabCounter:                   INT_MAX,
		TabFocusRequestCurrCounterRegular: INT_MAX,
		TabFocusRequestNextCounterRegular: INT_MAX,
		TabFocusRequestCurrCounterTabStop: INT_MAX,
		TabFocusRequestNextCounterTabStop: INT_MAX,
		MouseCursor:                       ImGuiMouseCursor_Arrow,
		DragDropSourceFrameCount:          -1,
		DragDropMouseButton:               -1,
		DragDropAcceptFrameCount:          -1,
		CurrentTableStackIdx:              -1,
		ColorEditLastColor:                [3]float{FLT_MAX, FLT_MAX, FLT_MAX},
		DragSpeedDefaultRatio:             1 / 100.0,
		TooltipSlowDelay:                  0.5,
		PlatformImePos:                    ImVec2{FLT_MAX, FLT_MAX},
		PlatformImeLastPos:                ImVec2{FLT_MAX, FLT_MAX},
		PlatformLocaleDecimalPoint:        '.',
		LogLinePosY:                       FLT_MAX,
		LogDepthToExpand:                  2,
		LogDepthToExpandDefault:           2,
		WantCaptureMouseNextFrame:         -1,
		WantCaptureKeyboardNextFrame:      -1,
		WantTextInputNextFrame:            -1,
	}
}

// CreateContext Context creation and access
//   - Each context create its own ImFontAtlas by default. You may instance one yourself and pass it to CreateContext() to share a font atlas between contexts.
//   - DLL users: heaps and globals are not shared across DLL boundaries! You will need to call SetCurrentContext() + SetAllocatorFunctions()
//     for each static/DLL boundary you are calling from. Read "Context and Memory Allocators" section of imgui.cpp for details.
func CreateContext(shared_font_atlas *ImFontAtlas) *ImGuiContext {
	var ctx = NewImGuiContext(shared_font_atlas)
	if g == nil {
		SetCurrentContext(&ctx)
	}
	Initialize(&ctx)
	return &ctx
}

// DestroyContext NULL = destroy current context
func DestroyContext(ctx *ImGuiContext) {
	if ctx == nil {
		ctx = g
	}
	Shutdown(ctx)
	if g == ctx {
		SetCurrentContext(nil)
	}
}

// GetCurrentContext Internal state access - if you want to share Dear ImGui state between modules (e.g. DLL) or allocate it yourself
// Note that we still point to some static data and members (such as GFontAtlas), so the state instance you end up using will point to the static data within its module
func GetCurrentContext() *ImGuiContext { return g }

func SetCurrentContext(ctx *ImGuiContext) { g = ctx }

// AddContextHook Generic context hooks
// No specific ordering/dependency support, will see as needed
func AddContextHook(context *ImGuiContext, hook *ImGuiContextHook) ImGuiID {
	var g = context
	IM_ASSERT(hook.Callback != nil && hook.HookId == 0 && hook.Type != ImGuiContextHookType_PendingRemoval_)
	g.Hooks = append(g.Hooks, *hook)
	g.HookIdNext++
	g.Hooks[len(g.Hooks)-1].HookId = g.HookIdNext
	return g.HookIdNext
}

// RemoveContextHook Deferred removal, avoiding issue with changing vector while iterating it
func RemoveContextHook(context *ImGuiContext, hook_to_remove ImGuiID) {
	var g = context
	IM_ASSERT(hook_to_remove != 0)
	for n := range g.Hooks {
		if g.Hooks[n].HookId == hook_to_remove {
			g.Hooks[n].Type = ImGuiContextHookType_PendingRemoval_
		}
	}
}

// Initialize Init
func Initialize(context *ImGuiContext) {
	var g = context
	IM_ASSERT(!g.Initialized && !g.SettingsLoaded)

	// Add .ini handle for ImGuiWindow type
	{
		var ini_handler ImGuiSettingsHandler
		ini_handler.TypeName = "Window"
		ini_handler.TypeHash = ImHashStr("Window", 0, 0)
		ini_handler.ClearAllFn = WindowSettingsHandler_ClearAll
		ini_handler.ReadOpenFn = WindowSettingsHandler_ReadOpen
		ini_handler.ReadLineFn = WindowSettingsHandler_ReadLine
		ini_handler.ApplyAllFn = WindowSettingsHandler_ApplyAll
		ini_handler.WriteAllFn = WindowSettingsHandler_WriteAll
		g.SettingsHandlers = append(g.SettingsHandlers, ini_handler)
	}

	// Add .ini handle for ImGuiTable type
	//TableSettingsInstallHandler(context)

	// Create default viewport
	var viewport = NewImGuiViewportP()
	g.Viewports = append(g.Viewports, &viewport)

	g.Initialized = true
}

// Shutdown Since 1.60 this is a _private_ function. You can call DestroyContext() to destroy the context created by CreateContext().
func Shutdown(context *ImGuiContext) {
	// The fonts atlas can be used prior to calling NewFrame(), so we clear it even if g.Initialized is FALSE (which would happen if we never called NewFrame)
	var g = context
	if g.IO.Fonts != nil && g.FontAtlasOwnedByContext {
		g.IO.Fonts.Locked = false
		g.IO.Fonts = nil
	}
	g.IO.Fonts = nil

	// Cleanup of other data are conditional on actually having initialized Dear ImGui.
	if !g.Initialized {
		return
	}

	// Save settings (unless we haven't attempted to load them: CreateContext/DestroyContext without a call to NewFrame shouldn't save an empty file)
	if g.SettingsLoaded && g.IO.IniFilename != "" {
		var backup_context = g
		SetCurrentContext(g)
		SaveIniSettingsToDisk(g.IO.IniFilename)
		SetCurrentContext(backup_context)
	}

	CallContextHooks(g, ImGuiContextHookType_Shutdown)

	// Clear everything else
	g.Windows = nil
	g.WindowsFocusOrder = nil
	g.WindowsTempSortBuffer = nil
	g.CurrentWindow = nil
	g.CurrentWindowStack = nil
	g.WindowsById.Clear()
	g.NavWindow = nil
	g.HoveredWindow = nil
	g.HoveredWindowUnderMovingWindow = nil
	g.ActiveIdWindow = nil
	g.ActiveIdPreviousFrameWindow = nil
	g.MovingWindow = nil
	g.ColorStack = nil
	g.StyleVarStack = nil
	g.FontStack = nil
	g.OpenPopupStack = nil
	g.BeginPopupStack = nil

	g.Viewports = nil

	g.TabBars = nil
	g.CurrentTabBarStack = nil
	g.ShrinkWidthBuffer = nil

	g.Tables = nil
	g.TablesTempDataStack = nil
	g.DrawChannelsTempMergeBuffer = nil

	g.ClipboardHandlerData = nil
	g.MenusIdSubmittedThisFrame = nil
	g.InputTextState = ImGuiInputTextState{}

	g.SettingsWindows = nil
	g.SettingsHandlers = nil

	if g.LogFile != nil {
		ImFileClose(g.LogFile)
		g.LogFile = nil
	}
	g.LogBuffer.Reset()

	g.Initialized = false
}
