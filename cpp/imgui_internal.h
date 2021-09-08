

//-----------------------------------------------------------------------------
// [SECTION] Generic context hooks
//-----------------------------------------------------------------------------

typedef void (*ImGuiContextHookCallback)(ImGuiContext* ctx, ImGuiContextHook* hook);
enum ImGuiContextHookType { ImGuiContextHookType_NewFramePre, ImGuiContextHookType_NewFramePost, ImGuiContextHookType_EndFramePre, ImGuiContextHookType_EndFramePost, ImGuiContextHookType_RenderPre, ImGuiContextHookType_RenderPost, ImGuiContextHookType_Shutdown, ImGuiContextHookType_PendingRemoval_ };

struct ImGuiContextHook
{
    ImGuiID                     HookId;     // A unique ID assigned by AddContextHook()
    ImGuiContextHookType        Type;
    ImGuiID                     Owner;
    ImGuiContextHookCallback    Callback;
    void*                       UserData;

    ImGuiContextHook()          { memset(this, 0, sizeof(*this)); }
};

//-----------------------------------------------------------------------------
// [SECTION] ImGuiContext (main imgui context)
//-----------------------------------------------------------------------------

struct ImGuiContext
{
    bool                    Initialized;
    bool                    FontAtlasOwnedByContext;            // IO.Fonts-> is owned by the ImGuiContext and will be destructed along with it.
    ImGuiIO                 IO;
    ImGuiStyle              Style;
    ImFont*                 Font;                               // (Shortcut) == FontStack.empty() ? IO.Font : FontStack.back()
    float                   FontSize;                           // (Shortcut) == FontBaseSize * g.CurrentWindow->FontWindowScale == window->FontSize(). Text height for current window.
    float                   FontBaseSize;                       // (Shortcut) == IO.FontGlobalScale * Font->Scale * Font->FontSize. Base text height.
    ImDrawListSharedData    DrawListSharedData;
    double                  Time;
    int                     FrameCount;
    int                     FrameCountEnded;
    int                     FrameCountRendered;
    bool                    WithinFrameScope;                   // Set by NewFrame(), cleared by EndFrame()
    bool                    WithinFrameScopeWithImplicitWindow; // Set by NewFrame(), cleared by EndFrame() when the implicit debug window has been pushed
    bool                    WithinEndChild;                     // Set within EndChild()
    bool                    GcCompactAll;                       // Request full GC
    bool                    TestEngineHookItems;                // Will call test engine hooks: ImGuiTestEngineHook_ItemAdd(), ImGuiTestEngineHook_ItemInfo(), ImGuiTestEngineHook_Log()
    ImGuiID                 TestEngineHookIdInfo;               // Will call test engine hooks: ImGuiTestEngineHook_IdInfo() from GetID()
    void*                   TestEngine;                         // Test engine user data

    // Windows state
    ImVector<ImGuiWindow*>  Windows;                            // Windows, sorted in display order, back to front
    ImVector<ImGuiWindow*>  WindowsFocusOrder;                  // Root windows, sorted in focus order, back to front.
    ImVector<ImGuiWindow*>  WindowsTempSortBuffer;              // Temporary buffer used in EndFrame() to reorder windows so parents are kept before their child
    ImVector<ImGuiWindowStackData> CurrentWindowStack;
    ImGuiStorage            WindowsById;                        // Map window's ImGuiID to ImGuiWindow*
    int                     WindowsActiveCount;                 // Number of unique windows submitted by frame
    ImVec2                  WindowsHoverPadding;                // Padding around resizable windows for which hovering on counts as hovering the window == ImMax(style.TouchExtraPadding, WINDOWS_HOVER_PADDING)
    ImGuiWindow*            CurrentWindow;                      // Window being drawn into
    ImGuiWindow*            HoveredWindow;                      // Window the mouse is hovering. Will typically catch mouse inputs.
    ImGuiWindow*            HoveredWindowUnderMovingWindow;     // Hovered window ignoring MovingWindow. Only set if MovingWindow is set.
    ImGuiWindow*            MovingWindow;                       // Track the window we clicked on (in order to preserve focus). The actual window that is moved is generally MovingWindow->RootWindow.
    ImGuiWindow*            WheelingWindow;                     // Track the window we started mouse-wheeling on. Until a timer elapse or mouse has moved, generally keep scrolling the same window even if during the course of scrolling the mouse ends up hovering a child window.
    ImVec2                  WheelingWindowRefMousePos;
    float                   WheelingWindowTimer;

    // Item/widgets state and tracking information
    ImGuiID                 HoveredId;                          // Hovered widget, filled during the frame
    ImGuiID                 HoveredIdPreviousFrame;
    bool                    HoveredIdAllowOverlap;
    bool                    HoveredIdUsingMouseWheel;           // Hovered widget will use mouse wheel. Blocks scrolling the underlying window.
    bool                    HoveredIdPreviousFrameUsingMouseWheel;
    bool                    HoveredIdDisabled;                  // At least one widget passed the rect test, but has been discarded by disabled flag or popup inhibit. May be true even if HoveredId == 0.
    float                   HoveredIdTimer;                     // Measure contiguous hovering time
    float                   HoveredIdNotActiveTimer;            // Measure contiguous hovering time where the item has not been active
    ImGuiID                 ActiveId;                           // Active widget
    ImGuiID                 ActiveIdIsAlive;                    // Active widget has been seen this frame (we can't use a bool as the ActiveId may change within the frame)
    float                   ActiveIdTimer;
    bool                    ActiveIdIsJustActivated;            // Set at the time of activation for one frame
    bool                    ActiveIdAllowOverlap;               // Active widget allows another widget to steal active id (generally for overlapping widgets, but not always)
    bool                    ActiveIdNoClearOnFocusLoss;         // Disable losing active id if the active id window gets unfocused.
    bool                    ActiveIdHasBeenPressedBefore;       // Track whether the active id led to a press (this is to allow changing between PressOnClick and PressOnRelease without pressing twice). Used by range_select branch.
    bool                    ActiveIdHasBeenEditedBefore;        // Was the value associated to the widget Edited over the course of the Active state.
    bool                    ActiveIdHasBeenEditedThisFrame;
    bool                    ActiveIdUsingMouseWheel;            // Active widget will want to read mouse wheel. Blocks scrolling the underlying window.
    ImU32                   ActiveIdUsingNavDirMask;            // Active widget will want to read those nav move requests (e.g. can activate a button and move away from it)
    ImU32                   ActiveIdUsingNavInputMask;          // Active widget will want to read those nav inputs.
    ImU64                   ActiveIdUsingKeyInputMask;          // Active widget will want to read those key inputs. When we grow the ImGuiKey enum we'll need to either to order the enum to make useful keys come first, either redesign this into e.g. a small array.
    ImVec2                  ActiveIdClickOffset;                // Clicked offset from upper-left corner, if applicable (currently only set by ButtonBehavior)
    ImGuiWindow*            ActiveIdWindow;
    ImGuiInputSource        ActiveIdSource;                     // Activating with mouse or nav (gamepad/keyboard)
    int                     ActiveIdMouseButton;
    ImGuiID                 ActiveIdPreviousFrame;
    bool                    ActiveIdPreviousFrameIsAlive;
    bool                    ActiveIdPreviousFrameHasBeenEditedBefore;
    ImGuiWindow*            ActiveIdPreviousFrameWindow;
    ImGuiID                 LastActiveId;                       // Store the last non-zero ActiveId, useful for animation.
    float                   LastActiveIdTimer;                  // Store the last non-zero ActiveId timer since the beginning of activation, useful for animation.

    // Next window/item data
    ImGuiItemFlags          CurrentItemFlags;                      // == g.ItemFlagsStack.back()
    ImGuiNextItemData       NextItemData;                       // Storage for SetNextItem** functions
    ImGuiLastItemData       LastItemData;                       // Storage for last submitted item (setup by ItemAdd)
    ImGuiNextWindowData     NextWindowData;                     // Storage for SetNextWindow** functions

    // Shared stacks
    ImVector<ImGuiColorMod> ColorStack;                         // Stack for PushStyleColor()/PopStyleColor() - inherited by Begin()
    ImVector<ImGuiStyleMod> StyleVarStack;                      // Stack for PushStyleVar()/PopStyleVar() - inherited by Begin()
    ImVector<ImFont*>       FontStack;                          // Stack for PushFont()/PopFont() - inherited by Begin()
    ImVector<ImGuiID>       FocusScopeStack;                    // Stack for PushFocusScope()/PopFocusScope() - not inherited by Begin(), unless child window
    ImVector<ImGuiItemFlags>ItemFlagsStack;                     // Stack for PushItemFlag()/PopItemFlag() - inherited by Begin()
    ImVector<ImGuiGroupData>GroupStack;                         // Stack for BeginGroup()/EndGroup() - not inherited by Begin()
    ImVector<ImGuiPopupData>OpenPopupStack;                     // Which popups are open (persistent)
    ImVector<ImGuiPopupData>BeginPopupStack;                    // Which level of BeginPopup() we are in (reset every frame)

    // Viewports
    ImVector<ImGuiViewportP*> Viewports;                        // Active viewports (Size==1 in 'master' branch). Each viewports hold their copy of ImDrawData.

    // Gamepad/keyboard Navigation
    ImGuiWindow*            NavWindow;                          // Focused window for navigation. Could be called 'FocusWindow'
    ImGuiID                 NavId;                              // Focused item for navigation
    ImGuiID                 NavFocusScopeId;                    // Identify a selection scope (selection code often wants to "clear other items" when landing on an item of the selection set)
    ImGuiID                 NavActivateId;                      // ~~ (g.ActiveId == 0) && IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0, also set when calling ActivateItem()
    ImGuiID                 NavActivateDownId;                  // ~~ IsNavInputDown(ImGuiNavInput_Activate) ? NavId : 0
    ImGuiID                 NavActivatePressedId;               // ~~ IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0
    ImGuiID                 NavInputId;                         // ~~ IsNavInputPressed(ImGuiNavInput_Input) ? NavId : 0
    ImGuiID                 NavJustTabbedId;                    // Just tabbed to this id.
    ImGuiID                 NavJustMovedToId;                   // Just navigated to this id (result of a successfully MoveRequest).
    ImGuiID                 NavJustMovedToFocusScopeId;         // Just navigated to this focus scope id (result of a successfully MoveRequest).
    ImGuiKeyModFlags        NavJustMovedToKeyMods;
    ImGuiID                 NavNextActivateId;                  // Set by ActivateItem(), queued until next frame.
    ImGuiInputSource        NavInputSource;                     // Keyboard or Gamepad mode? THIS WILL ONLY BE None or NavGamepad or NavKeyboard.
    ImGuiNavLayer           NavLayer;                           // Layer we are navigating on. For now the system is hard-coded for 0=main contents and 1=menu/title bar, may expose layers later.
    int                     NavIdTabCounter;                    // == NavWindow->DC.FocusIdxTabCounter at time of NavId processing
    bool                    NavIdIsAlive;                       // Nav widget has been seen this frame ~~ NavRectRel is valid
    bool                    NavMousePosDirty;                   // When set we will update mouse position if (io.ConfigFlags & ImGuiConfigFlags_NavEnableSetMousePos) if set (NB: this not enabled by default)
    bool                    NavDisableHighlight;                // When user starts using mouse, we hide gamepad/keyboard highlight (NB: but they are still available, which is why NavDisableHighlight isn't always != NavDisableMouseHover)
    bool                    NavDisableMouseHover;               // When user starts using gamepad/keyboard, we hide mouse hovering highlight until mouse is touched again.

    // Navigation: Init & Move Requests
    bool                    NavAnyRequest;                      // ~~ NavMoveRequest || NavInitRequest this is to perform early out in ItemAdd()
    bool                    NavInitRequest;                     // Init request for appearing window to select first item
    bool                    NavInitRequestFromMove;
    ImGuiID                 NavInitResultId;                    // Init request result (first item of the window, or one for which SetItemDefaultFocus() was called)
    ImRect                  NavInitResultRectRel;               // Init request result rectangle (relative to parent window)
    bool                    NavMoveSubmitted;                   // Move request submitted, will process result on next NewFrame()
    bool                    NavMoveScoringItems;                // Move request submitted, still scoring incoming items
    bool                    NavMoveForwardToNextFrame;
    ImGuiNavMoveFlags       NavMoveFlags;
    ImGuiKeyModFlags        NavMoveKeyMods;
    ImGuiDir                NavMoveDir;                         // Direction of the move request (left/right/up/down)
    ImGuiDir                NavMoveDirForDebug;
    ImGuiDir                NavMoveClipDir;                     // FIXME-NAV: Describe the purpose of this better. Might want to rename?
    ImRect                  NavScoringRect;                     // Rectangle used for scoring, in screen space. Based of window->NavRectRel[], modified for directional navigation scoring.
    int                     NavScoringDebugCount;               // Metrics for debugging
    ImGuiNavItemData        NavMoveResultLocal;                 // Best move request candidate within NavWindow
    ImGuiNavItemData        NavMoveResultLocalVisible;          // Best move request candidate within NavWindow that are mostly visible (when using ImGuiNavMoveFlags_AlsoScoreVisibleSet flag)
    ImGuiNavItemData        NavMoveResultOther;                 // Best move request candidate within NavWindow's flattened hierarchy (when using ImGuiWindowFlags_NavFlattened flag)

    // Navigation: Windowing (CTRL+TAB for list, or Menu button + keys or directional pads to move/resize)
    ImGuiWindow*            NavWindowingTarget;                 // Target window when doing CTRL+Tab (or Pad Menu + FocusPrev/Next), this window is temporarily displayed top-most!
    ImGuiWindow*            NavWindowingTargetAnim;             // Record of last valid NavWindowingTarget until DimBgRatio and NavWindowingHighlightAlpha becomes 0.0f, so the fade-out can stay on it.
    ImGuiWindow*            NavWindowingListWindow;             // Internal window actually listing the CTRL+Tab contents
    float                   NavWindowingTimer;
    float                   NavWindowingHighlightAlpha;
    bool                    NavWindowingToggleLayer;

    // Legacy Focus/Tabbing system (older than Nav, active even if Nav is disabled, misnamed. FIXME-NAV: This needs a redesign!)
    ImGuiWindow*            TabFocusRequestCurrWindow;          //
    ImGuiWindow*            TabFocusRequestNextWindow;          //
    int                     TabFocusRequestCurrCounterRegular;  // Any item being requested for focus, stored as an index (we on layout to be stable between the frame pressing TAB and the next frame, semi-ouch)
    int                     TabFocusRequestCurrCounterTabStop;  // Tab item being requested for focus, stored as an index
    int                     TabFocusRequestNextCounterRegular;  // Stored for next frame
    int                     TabFocusRequestNextCounterTabStop;  // "
    bool                    TabFocusPressed;                    // Set in NewFrame() when user pressed Tab

    // Render
    float                   DimBgRatio;                         // 0.0..1.0 animation when fading in a dimming background (for modal window and CTRL+TAB list)
    ImGuiMouseCursor        MouseCursor;

    // Drag and Drop
    bool                    DragDropActive;
    bool                    DragDropWithinSource;               // Set when within a BeginDragDropXXX/EndDragDropXXX block for a drag source.
    bool                    DragDropWithinTarget;               // Set when within a BeginDragDropXXX/EndDragDropXXX block for a drag target.
    ImGuiDragDropFlags      DragDropSourceFlags;
    int                     DragDropSourceFrameCount;
    int                     DragDropMouseButton;
    ImGuiPayload            DragDropPayload;
    ImRect                  DragDropTargetRect;                 // Store rectangle of current target candidate (we favor small targets when overlapping)
    ImGuiID                 DragDropTargetId;
    ImGuiDragDropFlags      DragDropAcceptFlags;
    float                   DragDropAcceptIdCurrRectSurface;    // Target item surface (we resolve overlapping targets by prioritizing the smaller surface)
    ImGuiID                 DragDropAcceptIdCurr;               // Target item id (set at the time of accepting the payload)
    ImGuiID                 DragDropAcceptIdPrev;               // Target item id from previous frame (we need to store this to allow for overlapping drag and drop targets)
    int                     DragDropAcceptFrameCount;           // Last time a target expressed a desire to accept the source
    ImGuiID                 DragDropHoldJustPressedId;          // Set when holding a payload just made ButtonBehavior() return a press.
    ImVector<unsigned char> DragDropPayloadBufHeap;             // We don't expose the ImVector<> directly, ImGuiPayload only holds pointer+size
    unsigned char           DragDropPayloadBufLocal[16];        // Local buffer for small payloads

    // Table
    ImGuiTable*                     CurrentTable;
    int                             CurrentTableStackIdx;
    ImPool<ImGuiTable>              Tables;
    ImVector<ImGuiTableTempData>    TablesTempDataStack;
    ImVector<float>                 TablesLastTimeActive;       // Last used timestamp of each tables (SOA, for efficient GC)
    ImVector<ImDrawChannel>         DrawChannelsTempMergeBuffer;

    // Tab bars
    ImGuiTabBar*                    CurrentTabBar;
    ImPool<ImGuiTabBar>             TabBars;
    ImVector<ImGuiPtrOrIndex>       CurrentTabBarStack;
    ImVector<ImGuiShrinkWidthItem>  ShrinkWidthBuffer;

    // Widget state
    ImVec2                  MouseLastValidPos;
    ImGuiInputTextState     InputTextState;
    ImFont                  InputTextPasswordFont;
    ImGuiID                 TempInputId;                        // Temporary text input when CTRL+clicking on a slider, etc.
    ImGuiColorEditFlags     ColorEditOptions;                   // Store user options for color edit widgets
    float                   ColorEditLastHue;                   // Backup of last Hue associated to LastColor[3], so we can restore Hue in lossy RGB<>HSV round trips
    float                   ColorEditLastSat;                   // Backup of last Saturation associated to LastColor[3], so we can restore Saturation in lossy RGB<>HSV round trips
    float                   ColorEditLastColor[3];
    ImVec4                  ColorPickerRef;                     // Initial/reference color at the time of opening the color picker.
    ImGuiComboPreviewData   ComboPreviewData;
    float                   SliderCurrentAccum;                 // Accumulated slider delta when using navigation controls.
    bool                    SliderCurrentAccumDirty;            // Has the accumulated slider delta changed since last time we tried to apply it?
    bool                    DragCurrentAccumDirty;
    float                   DragCurrentAccum;                   // Accumulator for dragging modification. Always high-precision, not rounded by end-user precision settings
    float                   DragSpeedDefaultRatio;              // If speed == 0.0f, uses (max-min) * DragSpeedDefaultRatio
    float                   DisabledAlphaBackup;                // Backup for style.Alpha for BeginDisabled()
    float                   ScrollbarClickDeltaToGrabCenter;    // Distance between mouse and center of grab box, normalized in parent space. Use storage?
    int                     TooltipOverrideCount;
    float                   TooltipSlowDelay;                   // Time before slow tooltips appears (FIXME: This is temporary until we merge in tooltip timer+priority work)
    ImVector<char>          ClipboardHandlerData;               // If no custom clipboard handler is defined
    ImVector<ImGuiID>       MenusIdSubmittedThisFrame;          // A list of menu IDs that were rendered at least once

    // Platform support
    ImVec2                  PlatformImePos;                     // Cursor position request & last passed to the OS Input Method Editor
    ImVec2                  PlatformImeLastPos;
    char                    PlatformLocaleDecimalPoint;         // '.' or *localeconv()->decimal_point

    // Settings
    bool                    SettingsLoaded;
    float                   SettingsDirtyTimer;                 // Save .ini Settings to memory when time reaches zero
    ImGuiTextBuffer         SettingsIniData;                    // In memory .ini settings
    ImVector<ImGuiSettingsHandler>      SettingsHandlers;       // List of .ini settings handlers
    ImChunkStream<ImGuiWindowSettings>  SettingsWindows;        // ImGuiWindow .ini settings entries
    ImChunkStream<ImGuiTableSettings>   SettingsTables;         // ImGuiTable .ini settings entries
    ImVector<ImGuiContextHook>          Hooks;                  // Hooks for extensions (e.g. test engine)
    ImGuiID                             HookIdNext;             // Next available HookId

    // Capture/Logging
    bool                    LogEnabled;                         // Currently capturing
    ImGuiLogType            LogType;                            // Capture target
    ImFileHandle            LogFile;                            // If != NULL log to stdout/ file
    ImGuiTextBuffer         LogBuffer;                          // Accumulation buffer when log to clipboard. This is pointer so our GImGui static constructor doesn't call heap allocators.
    const char*             LogNextPrefix;
    const char*             LogNextSuffix;
    float                   LogLinePosY;
    bool                    LogLineFirstItem;
    int                     LogDepthRef;
    int                     LogDepthToExpand;
    int                     LogDepthToExpandDefault;            // Default/stored value for LogDepthMaxExpand if not specified in the LogXXX function call.

    // Debug Tools
    bool                    DebugItemPickerActive;              // Item picker is active (started with DebugStartItemPicker())
    ImGuiID                 DebugItemPickerBreakId;             // Will call IM_DEBUG_BREAK() when encountering this id
    ImGuiMetricsConfig      DebugMetricsConfig;

    // Misc
    float                   FramerateSecPerFrame[120];          // Calculate estimate of framerate for user over the last 2 seconds.
    int                     FramerateSecPerFrameIdx;
    int                     FramerateSecPerFrameCount;
    float                   FramerateSecPerFrameAccum;
    int                     WantCaptureMouseNextFrame;          // Explicit capture via CaptureKeyboardFromApp()/CaptureMouseFromApp() sets those flags
    int                     WantCaptureKeyboardNextFrame;
    int                     WantTextInputNextFrame;
    char                    TempBuffer[1024 * 3 + 1];           // Temporary text buffer

    ImGuiContext(ImFontAtlas* shared_font_atlas)
    {
        Initialized = false;
        FontAtlasOwnedByContext = shared_font_atlas ? false : true;
        Font = NULL;
        FontSize = FontBaseSize = 0.0f;
        IO.Fonts = shared_font_atlas ? shared_font_atlas : IM_NEW(ImFontAtlas)();
        Time = 0.0f;
        FrameCount = 0;
        FrameCountEnded = FrameCountRendered = -1;
        WithinFrameScope = WithinFrameScopeWithImplicitWindow = WithinEndChild = false;
        GcCompactAll = false;
        TestEngineHookItems = false;
        TestEngineHookIdInfo = 0;
        TestEngine = NULL;

        WindowsActiveCount = 0;
        CurrentWindow = NULL;
        HoveredWindow = NULL;
        HoveredWindowUnderMovingWindow = NULL;
        MovingWindow = NULL;
        WheelingWindow = NULL;
        WheelingWindowTimer = 0.0f;

        HoveredId = HoveredIdPreviousFrame = 0;
        HoveredIdAllowOverlap = false;
        HoveredIdUsingMouseWheel = HoveredIdPreviousFrameUsingMouseWheel = false;
        HoveredIdDisabled = false;
        HoveredIdTimer = HoveredIdNotActiveTimer = 0.0f;
        ActiveId = 0;
        ActiveIdIsAlive = 0;
        ActiveIdTimer = 0.0f;
        ActiveIdIsJustActivated = false;
        ActiveIdAllowOverlap = false;
        ActiveIdNoClearOnFocusLoss = false;
        ActiveIdHasBeenPressedBefore = false;
        ActiveIdHasBeenEditedBefore = false;
        ActiveIdHasBeenEditedThisFrame = false;
        ActiveIdUsingMouseWheel = false;
        ActiveIdUsingNavDirMask = 0x00;
        ActiveIdUsingNavInputMask = 0x00;
        ActiveIdUsingKeyInputMask = 0x00;
        ActiveIdClickOffset = ImVec2(-1, -1);
        ActiveIdWindow = NULL;
        ActiveIdSource = ImGuiInputSource_None;
        ActiveIdMouseButton = -1;
        ActiveIdPreviousFrame = 0;
        ActiveIdPreviousFrameIsAlive = false;
        ActiveIdPreviousFrameHasBeenEditedBefore = false;
        ActiveIdPreviousFrameWindow = NULL;
        LastActiveId = 0;
        LastActiveIdTimer = 0.0f;

        CurrentItemFlags = ImGuiItemFlags_None;

        NavWindow = NULL;
        NavId = NavFocusScopeId = NavActivateId = NavActivateDownId = NavActivatePressedId = NavInputId = 0;
        NavJustTabbedId = NavJustMovedToId = NavJustMovedToFocusScopeId = NavNextActivateId = 0;
        NavJustMovedToKeyMods = ImGuiKeyModFlags_None;
        NavInputSource = ImGuiInputSource_None;
        NavLayer = ImGuiNavLayer_Main;
        NavIdTabCounter = INT_MAX;
        NavIdIsAlive = false;
        NavMousePosDirty = false;
        NavDisableHighlight = true;
        NavDisableMouseHover = false;
        NavAnyRequest = false;
        NavInitRequest = false;
        NavInitRequestFromMove = false;
        NavInitResultId = 0;
        NavMoveSubmitted = false;
        NavMoveScoringItems = false;
        NavMoveForwardToNextFrame = false;
        NavMoveFlags = ImGuiNavMoveFlags_None;
        NavMoveKeyMods = ImGuiKeyModFlags_None;
        NavMoveDir = NavMoveDirForDebug = NavMoveClipDir = ImGuiDir_None;
        NavScoringDebugCount = 0;

        NavWindowingTarget = NavWindowingTargetAnim = NavWindowingListWindow = NULL;
        NavWindowingTimer = NavWindowingHighlightAlpha = 0.0f;
        NavWindowingToggleLayer = false;

        TabFocusRequestCurrWindow = TabFocusRequestNextWindow = NULL;
        TabFocusRequestCurrCounterRegular = TabFocusRequestCurrCounterTabStop = INT_MAX;
        TabFocusRequestNextCounterRegular = TabFocusRequestNextCounterTabStop = INT_MAX;
        TabFocusPressed = false;

        DimBgRatio = 0.0f;
        MouseCursor = ImGuiMouseCursor_Arrow;

        DragDropActive = DragDropWithinSource = DragDropWithinTarget = false;
        DragDropSourceFlags = ImGuiDragDropFlags_None;
        DragDropSourceFrameCount = -1;
        DragDropMouseButton = -1;
        DragDropTargetId = 0;
        DragDropAcceptFlags = ImGuiDragDropFlags_None;
        DragDropAcceptIdCurrRectSurface = 0.0f;
        DragDropAcceptIdPrev = DragDropAcceptIdCurr = 0;
        DragDropAcceptFrameCount = -1;
        DragDropHoldJustPressedId = 0;
        memset(DragDropPayloadBufLocal, 0, sizeof(DragDropPayloadBufLocal));

        CurrentTable = NULL;
        CurrentTableStackIdx = -1;
        CurrentTabBar = NULL;

        TempInputId = 0;
        ColorEditOptions = ImGuiColorEditFlags_DefaultOptions_;
        ColorEditLastHue = ColorEditLastSat = 0.0f;
        ColorEditLastColor[0] = ColorEditLastColor[1] = ColorEditLastColor[2] = FLT_MAX;
        SliderCurrentAccum = 0.0f;
        SliderCurrentAccumDirty = false;
        DragCurrentAccumDirty = false;
        DragCurrentAccum = 0.0f;
        DragSpeedDefaultRatio = 1.0f / 100.0f;
        DisabledAlphaBackup = 0.0f;
        ScrollbarClickDeltaToGrabCenter = 0.0f;
        TooltipOverrideCount = 0;
        TooltipSlowDelay = 0.50f;

        PlatformImePos = PlatformImeLastPos = ImVec2(FLT_MAX, FLT_MAX);
        PlatformLocaleDecimalPoint = '.';

        SettingsLoaded = false;
        SettingsDirtyTimer = 0.0f;
        HookIdNext = 0;

        LogEnabled = false;
        LogType = ImGuiLogType_None;
        LogNextPrefix = LogNextSuffix = NULL;
        LogFile = NULL;
        LogLinePosY = FLT_MAX;
        LogLineFirstItem = false;
        LogDepthRef = 0;
        LogDepthToExpand = LogDepthToExpandDefault = 2;

        DebugItemPickerActive = false;
        DebugItemPickerBreakId = 0;

        memset(FramerateSecPerFrame, 0, sizeof(FramerateSecPerFrame));
        FramerateSecPerFrameIdx = FramerateSecPerFrameCount = 0;
        FramerateSecPerFrameAccum = 0.0f;
        WantCaptureMouseNextFrame = WantCaptureKeyboardNextFrame = WantTextInputNextFrame = -1;
        memset(TempBuffer, 0, sizeof(TempBuffer));
    }
};

//-----------------------------------------------------------------------------
// [SECTION] ImGuiWindowTempData, ImGuiWindow
//-----------------------------------------------------------------------------

// Transient per-window data, reset at the beginning of the frame. This used to be called ImGuiDrawContext, hence the DC variable name in ImGuiWindow.
// (That's theory, in practice the delimitation between ImGuiWindow and ImGuiWindowTempData is quite tenuous and could be reconsidered..)
// (This doesn't need a constructor because we zero-clear it as part of ImGuiWindow and all frame-temporary data are setup on Begin)
struct IMGUI_API ImGuiWindowTempData
{
    // Layout
    ImVec2                  CursorPos;              // Current emitting position, in absolute coordinates.
    ImVec2                  CursorPosPrevLine;
    ImVec2                  CursorStartPos;         // Initial position after Begin(), generally ~ window position + WindowPadding.
    ImVec2                  CursorMaxPos;           // Used to implicitly calculate ContentSize at the beginning of next frame, for scrolling range and auto-resize. Always growing during the frame.
    ImVec2                  IdealMaxPos;            // Used to implicitly calculate ContentSizeIdeal at the beginning of next frame, for auto-resize only. Always growing during the frame.
    ImVec2                  CurrLineSize;
    ImVec2                  PrevLineSize;
    float                   CurrLineTextBaseOffset; // Baseline offset (0.0f by default on a new line, generally == style.FramePadding.y when a framed item has been added).
    float                   PrevLineTextBaseOffset;
    ImVec1                  Indent;                 // Indentation / start position from left of window (increased by TreePush/TreePop, etc.)
    ImVec1                  ColumnsOffset;          // Offset to the current column (if ColumnsCurrent > 0). FIXME: This and the above should be a stack to allow use cases like Tree->Column->Tree. Need revamp columns API.
    ImVec1                  GroupOffset;

    // Keyboard/Gamepad navigation
    ImGuiNavLayer           NavLayerCurrent;        // Current layer, 0..31 (we currently only use 0..1)
    short                   NavLayersActiveMask;    // Which layers have been written to (result from previous frame)
    short                   NavLayersActiveMaskNext;// Which layers have been written to (accumulator for current frame)
    ImGuiID                 NavFocusScopeIdCurrent; // Current focus scope ID while appending
    bool                    NavHideHighlightOneFrame;
    bool                    NavHasScroll;           // Set when scrolling can be used (ScrollMax > 0.0f)

    // Miscellaneous
    bool                    MenuBarAppending;       // FIXME: Remove this
    ImVec2                  MenuBarOffset;          // MenuBarOffset.x is sort of equivalent of a per-layer CursorPos.x, saved/restored as we switch to the menu bar. The only situation when MenuBarOffset.y is > 0 if when (SafeAreaPadding.y > FramePadding.y), often used on TVs.
    ImGuiMenuColumns        MenuColumns;            // Simplified columns storage for menu items measurement
    int                     TreeDepth;              // Current tree depth.
    ImU32                   TreeJumpToParentOnPopMask; // Store a copy of !g.NavIdIsAlive for TreeDepth 0..31.. Could be turned into a ImU64 if necessary.
    ImVector<ImGuiWindow*>  ChildWindows;
    ImGuiStorage*           StateStorage;           // Current persistent per-window storage (store e.g. tree node open/close state)
    ImGuiOldColumns*        CurrentColumns;         // Current columns set
    int                     CurrentTableIdx;        // Current table index (into g.Tables)
    ImGuiLayoutType         LayoutType;
    ImGuiLayoutType         ParentLayoutType;       // Layout type of parent window at the time of Begin()
    int                     FocusCounterRegular;    // (Legacy Focus/Tabbing system) Sequential counter, start at -1 and increase when ImGuiItemFlags_Inputable (FIXME-NAV: Needs redesign)
    int                     FocusCounterTabStop;    // (Legacy Focus/Tabbing system) Same, but only count widgets which you can Tab through.

    // Local parameters stacks
    // We store the current settings outside of the vectors to increase memory locality (reduce cache misses). The vectors are rarely modified. Also it allows us to not heap allocate for short-lived windows which are not using those settings.
    float                   ItemWidth;              // Current item width (>0.0: width in pixels, <0.0: align xx pixels to the right of window).
    float                   TextWrapPos;            // Current text wrap pos.
    ImVector<float>         ItemWidthStack;         // Store item widths to restore (attention: .back() is not == ItemWidth)
    ImVector<float>         TextWrapPosStack;       // Store text wrap pos to restore (attention: .back() is not == TextWrapPos)
    ImGuiStackSizes         StackSizesOnBegin;      // Store size of various stacks for asserting
};

// Storage for one window
struct IMGUI_API ImGuiWindow
{
    char*                   Name;                               // Window name, owned by the window.
    ImGuiID                 ID;                                 // == ImHashStr(Name)
    ImGuiWindowFlags        Flags;                              // See enum ImGuiWindowFlags_
    ImVec2                  Pos;                                // Position (always rounded-up to nearest pixel)
    ImVec2                  Size;                               // Current size (==SizeFull or collapsed title bar size)
    ImVec2                  SizeFull;                           // Size when non collapsed
    ImVec2                  ContentSize;                        // Size of contents/scrollable client area (calculated from the extents reach of the cursor) from previous frame. Does not include window decoration or window padding.
    ImVec2                  ContentSizeIdeal;
    ImVec2                  ContentSizeExplicit;                // Size of contents/scrollable client area explicitly request by the user via SetNextWindowContentSize().
    ImVec2                  WindowPadding;                      // Window padding at the time of Begin().
    float                   WindowRounding;                     // Window rounding at the time of Begin(). May be clamped lower to avoid rendering artifacts with title bar, menu bar etc.
    float                   WindowBorderSize;                   // Window border size at the time of Begin().
    int                     NameBufLen;                         // Size of buffer storing Name. May be larger than strlen(Name)!
    ImGuiID                 MoveId;                             // == window->GetID("#MOVE")
    ImGuiID                 ChildId;                            // ID of corresponding item in parent window (for navigation to return from child window to parent window)
    ImVec2                  Scroll;
    ImVec2                  ScrollMax;
    ImVec2                  ScrollTarget;                       // target scroll position. stored as cursor position with scrolling canceled out, so the highest point is always 0.0f. (FLT_MAX for no change)
    ImVec2                  ScrollTargetCenterRatio;            // 0.0f = scroll so that target position is at top, 0.5f = scroll so that target position is centered
    ImVec2                  ScrollTargetEdgeSnapDist;           // 0.0f = no snapping, >0.0f snapping threshold
    ImVec2                  ScrollbarSizes;                     // Size taken by each scrollbars on their smaller axis. Pay attention! ScrollbarSizes.x == width of the vertical scrollbar, ScrollbarSizes.y = height of the horizontal scrollbar.
    bool                    ScrollbarX, ScrollbarY;             // Are scrollbars visible?
    bool                    Active;                             // Set to true on Begin(), unless Collapsed
    bool                    WasActive;
    bool                    WriteAccessed;                      // Set to true when any widget access the current window
    bool                    Collapsed;                          // Set when collapsing window to become only title-bar
    bool                    WantCollapseToggle;
    bool                    SkipItems;                          // Set when items can safely be all clipped (e.g. window not visible or collapsed)
    bool                    Appearing;                          // Set during the frame where the window is appearing (or re-appearing)
    bool                    Hidden;                             // Do not display (== HiddenFrames*** > 0)
    bool                    IsFallbackWindow;                   // Set on the "Debug##Default" window.
    bool                    HasCloseButton;                     // Set when the window has a close button (p_open != NULL)
    signed char             ResizeBorderHeld;                   // Current border being held for resize (-1: none, otherwise 0-3)
    short                   BeginCount;                         // Number of Begin() during the current frame (generally 0 or 1, 1+ if appending via multiple Begin/End pairs)
    short                   BeginOrderWithinParent;             // Begin() order within immediate parent window, if we are a child window. Otherwise 0.
    short                   BeginOrderWithinContext;            // Begin() order within entire imgui context. This is mostly used for debugging submission order related issues.
    short                   FocusOrder;                         // Order within WindowsFocusOrder[], altered when windows are focused.
    ImGuiID                 PopupId;                            // ID in the popup stack when this window is used as a popup/menu (because we use generic Name/ID for recycling)
    ImS8                    AutoFitFramesX, AutoFitFramesY;
    ImS8                    AutoFitChildAxises;
    bool                    AutoFitOnlyGrows;
    ImGuiDir                AutoPosLastDirection;
    ImS8                    HiddenFramesCanSkipItems;           // Hide the window for N frames
    ImS8                    HiddenFramesCannotSkipItems;        // Hide the window for N frames while allowing items to be submitted so we can measure their size
    ImS8                    HiddenFramesForRenderOnly;          // Hide the window until frame N at Render() time only
    ImS8                    DisableInputsFrames;                // Disable window interactions for N frames
    ImGuiCond               SetWindowPosAllowFlags : 8;         // store acceptable condition flags for SetNextWindowPos() use.
    ImGuiCond               SetWindowSizeAllowFlags : 8;        // store acceptable condition flags for SetNextWindowSize() use.
    ImGuiCond               SetWindowCollapsedAllowFlags : 8;   // store acceptable condition flags for SetNextWindowCollapsed() use.
    ImVec2                  SetWindowPosVal;                    // store window position when using a non-zero Pivot (position set needs to be processed when we know the window size)
    ImVec2                  SetWindowPosPivot;                  // store window pivot for positioning. ImVec2(0, 0) when positioning from top-left corner; ImVec2(0.5f, 0.5f) for centering; ImVec2(1, 1) for bottom right.

    ImVector<ImGuiID>       IDStack;                            // ID stack. ID are hashes seeded with the value at the top of the stack. (In theory this should be in the TempData structure)
    ImGuiWindowTempData     DC;                                 // Temporary per-window data, reset at the beginning of the frame. This used to be called ImGuiDrawContext, hence the "DC" variable name.

    // The best way to understand what those rectangles are is to use the 'Metrics->Tools->Show Windows Rectangles' viewer.
    // The main 'OuterRect', omitted as a field, is window->Rect().
    ImRect                  OuterRectClipped;                   // == Window->Rect() just after setup in Begin(). == window->Rect() for root window.
    ImRect                  InnerRect;                          // Inner rectangle (omit title bar, menu bar, scroll bar)
    ImRect                  InnerClipRect;                      // == InnerRect shrunk by WindowPadding*0.5f on each side, clipped within viewport or parent clip rect.
    ImRect                  WorkRect;                           // Initially covers the whole scrolling region. Reduced by containers e.g columns/tables when active. Shrunk by WindowPadding*1.0f on each side. This is meant to replace ContentRegionRect over time (from 1.71+ onward).
    ImRect                  ParentWorkRect;                     // Backup of WorkRect before entering a container such as columns/tables. Used by e.g. SpanAllColumns functions to easily access. Stacked containers are responsible for maintaining this. // FIXME-WORKRECT: Could be a stack?
    ImRect                  ClipRect;                           // Current clipping/scissoring rectangle, evolve as we are using PushClipRect(), etc. == DrawList->clip_rect_stack.back().
    ImRect                  ContentRegionRect;                  // FIXME: This is currently confusing/misleading. It is essentially WorkRect but not handling of scrolling. We currently rely on it as right/bottom aligned sizing operation need some size to rely on.
    ImVec2ih                HitTestHoleSize;                    // Define an optional rectangular hole where mouse will pass-through the window.
    ImVec2ih                HitTestHoleOffset;

    int                     LastFrameActive;                    // Last frame number the window was Active.
    float                   LastTimeActive;                     // Last timestamp the window was Active (using float as we don't need high precision there)
    float                   ItemWidthDefault;
    ImGuiStorage            StateStorage;
    ImVector<ImGuiOldColumns> ColumnsStorage;
    float                   FontWindowScale;                    // User scale multiplier per-window, via SetWindowFontScale()
    int                     SettingsOffset;                     // Offset into SettingsWindows[] (offsets are always valid as we only grow the array from the back)

    ImDrawList*             DrawList;                           // == &DrawListInst (for backward compatibility reason with code using imgui_internal.h we keep this a pointer)
    ImDrawList              DrawListInst;
    ImGuiWindow*            ParentWindow;                       // If we are a child _or_ popup window, this is pointing to our parent. Otherwise NULL.
    ImGuiWindow*            RootWindow;                         // Point to ourself or first ancestor that is not a child window == Top-level window.
    ImGuiWindow*            RootWindowForTitleBarHighlight;     // Point to ourself or first ancestor which will display TitleBgActive color when this window is active.
    ImGuiWindow*            RootWindowForNav;                   // Point to ourself or first ancestor which doesn't have the NavFlattened flag.

    ImGuiWindow*            NavLastChildNavWindow;              // When going to the menu bar, we remember the child window we came from. (This could probably be made implicit if we kept g.Windows sorted by last focused including child window.)
    ImGuiID                 NavLastIds[ImGuiNavLayer_COUNT];    // Last known NavId for this window, per layer (0/1)
    ImRect                  NavRectRel[ImGuiNavLayer_COUNT];    // Reference rectangle, in window relative space

    int                     MemoryDrawListIdxCapacity;          // Backup of last idx/vtx count, so when waking up the window we can preallocate and avoid iterative alloc/copy
    int                     MemoryDrawListVtxCapacity;
    bool                    MemoryCompacted;                    // Set when window extraneous data have been garbage collected

public:
    ImGuiWindow(ImGuiContext* context, const char* name);
    ~ImGuiWindow();

    ImGuiID     GetID(const char* str, const char* str_end = NULL);
    ImGuiID     GetID(const void* ptr);
    ImGuiID     GetID(int n);
    ImGuiID     GetIDNoKeepAlive(const char* str, const char* str_end = NULL);
    ImGuiID     GetIDNoKeepAlive(const void* ptr);
    ImGuiID     GetIDNoKeepAlive(int n);
    ImGuiID     GetIDFromRectangle(const ImRect& r_abs);

    // We don't use g.FontSize because the window may be != g.CurrentWidow.
    ImRect      Rect() const            { return ImRect(Pos.x, Pos.y, Pos.x + Size.x, Pos.y + Size.y); }
    float       CalcFontSize() const    { ImGuiContext& g = *GImGui; float scale = g.FontBaseSize * FontWindowScale; if (ParentWindow) scale *= ParentWindow->FontWindowScale; return scale; }
    float       TitleBarHeight() const  { ImGuiContext& g = *GImGui; return (Flags & ImGuiWindowFlags_NoTitleBar) ? 0.0f : CalcFontSize() + g.Style.FramePadding.y * 2.0f; }
    ImRect      TitleBarRect() const    { return ImRect(Pos, ImVec2(Pos.x + SizeFull.x, Pos.y + TitleBarHeight())); }
    float       MenuBarHeight() const   { ImGuiContext& g = *GImGui; return (Flags & ImGuiWindowFlags_MenuBar) ? DC.MenuBarOffset.y + CalcFontSize() + g.Style.FramePadding.y * 2.0f : 0.0f; }
    ImRect      MenuBarRect() const     { float y1 = Pos.y + TitleBarHeight(); return ImRect(Pos.x, y1, Pos.x + SizeFull.x, y1 + MenuBarHeight()); }
};

//-----------------------------------------------------------------------------
// [SECTION] Tab bar, Tab item support
//-----------------------------------------------------------------------------

// Extend ImGuiTabBarFlags_
enum ImGuiTabBarFlagsPrivate_
{
    ImGuiTabBarFlags_DockNode                   = 1 << 20,  // Part of a dock node [we don't use this in the master branch but it facilitate branch syncing to keep this around]
    ImGuiTabBarFlags_IsFocused                  = 1 << 21,
    ImGuiTabBarFlags_SaveSettings               = 1 << 22   // FIXME: Settings are handled by the docking system, this only request the tab bar to mark settings dirty when reordering tabs
};

// Extend ImGuiTabItemFlags_
enum ImGuiTabItemFlagsPrivate_
{
    ImGuiTabItemFlags_SectionMask_              = ImGuiTabItemFlags_Leading | ImGuiTabItemFlags_Trailing,
    ImGuiTabItemFlags_NoCloseButton             = 1 << 20,  // Track whether p_open was set or not (we'll need this info on the next frame to recompute ContentWidth during layout)
    ImGuiTabItemFlags_Button                    = 1 << 21   // Used by TabItemButton, change the tab item behavior to mimic a button
};

// Storage for one active tab item (sizeof() 40 bytes)
struct ImGuiTabItem
{
    ImGuiID             ID;
    ImGuiTabItemFlags   Flags;
    int                 LastFrameVisible;
    int                 LastFrameSelected;      // This allows us to infer an ordered list of the last activated tabs with little maintenance
    float               Offset;                 // Position relative to beginning of tab
    float               Width;                  // Width currently displayed
    float               ContentWidth;           // Width of label, stored during BeginTabItem() call
    ImS32               NameOffset;             // When Window==NULL, offset to name within parent ImGuiTabBar::TabsNames
    ImS16               BeginOrder;             // BeginTabItem() order, used to re-order tabs after toggling ImGuiTabBarFlags_Reorderable
    ImS16               IndexDuringLayout;      // Index only used during TabBarLayout()
    bool                WantClose;              // Marked as closed by SetTabItemClosed()

    ImGuiTabItem()      { memset(this, 0, sizeof(*this)); LastFrameVisible = LastFrameSelected = -1; NameOffset = -1; BeginOrder = IndexDuringLayout = -1; }
};

// Storage for a tab bar (sizeof() 152 bytes)
struct ImGuiTabBar
{
    ImVector<ImGuiTabItem> Tabs;
    ImGuiTabBarFlags    Flags;
    ImGuiID             ID;                     // Zero for tab-bars used by docking
    ImGuiID             SelectedTabId;          // Selected tab/window
    ImGuiID             NextSelectedTabId;      // Next selected tab/window. Will also trigger a scrolling animation
    ImGuiID             VisibleTabId;           // Can occasionally be != SelectedTabId (e.g. when previewing contents for CTRL+TAB preview)
    int                 CurrFrameVisible;
    int                 PrevFrameVisible;
    ImRect              BarRect;
    float               CurrTabsContentsHeight;
    float               PrevTabsContentsHeight; // Record the height of contents submitted below the tab bar
    float               WidthAllTabs;           // Actual width of all tabs (locked during layout)
    float               WidthAllTabsIdeal;      // Ideal width if all tabs were visible and not clipped
    float               ScrollingAnim;
    float               ScrollingTarget;
    float               ScrollingTargetDistToVisibility;
    float               ScrollingSpeed;
    float               ScrollingRectMinX;
    float               ScrollingRectMaxX;
    ImGuiID             ReorderRequestTabId;
    ImS16               ReorderRequestOffset;
    ImS8                BeginCount;
    bool                WantLayout;
    bool                VisibleTabWasSubmitted;
    bool                TabsAddedNew;           // Set to true when a new tab item or button has been added to the tab bar during last frame
    ImS16               TabsActiveCount;        // Number of tabs submitted this frame.
    ImS16               LastTabItemIdx;         // Index of last BeginTabItem() tab for use by EndTabItem()
    float               ItemSpacingY;
    ImVec2              FramePadding;           // style.FramePadding locked at the time of BeginTabBar()
    ImVec2              BackupCursorPos;
    ImGuiTextBuffer     TabsNames;              // For non-docking tab bar we re-append names in a contiguous buffer.

    ImGuiTabBar();
    int                 GetTabOrder(const ImGuiTabItem* tab) const  { return Tabs.index_from_ptr(tab); }
    const char*         GetTabName(const ImGuiTabItem* tab) const
    {
        IM_ASSERT(tab->NameOffset != -1 && tab->NameOffset < TabsNames.Buf.Size);
        return TabsNames.Buf.Data + tab->NameOffset;
    }
};

//-----------------------------------------------------------------------------
// [SECTION] Table support
//-----------------------------------------------------------------------------

#define IM_COL32_DISABLE                IM_COL32(0,0,0,1)   // Special sentinel code which cannot be used as a regular color.
#define IMGUI_TABLE_MAX_COLUMNS         64                  // sizeof(ImU64) * 8. This is solely because we frequently encode columns set in a ImU64.
#define IMGUI_TABLE_MAX_DRAW_CHANNELS   (4 + 64 * 2)        // See TableSetupDrawChannels()

// Our current column maximum is 64 but we may raise that in the future.
typedef ImS8 ImGuiTableColumnIdx;
typedef ImU8 ImGuiTableDrawChannelIdx;

// [Internal] sizeof() ~ 104
// We use the terminology "Enabled" to refer to a column that is not Hidden by user/api.
// We use the terminology "Clipped" to refer to a column that is out of sight because of scrolling/clipping.
// This is in contrast with some user-facing api such as IsItemVisible() / IsRectVisible() which use "Visible" to mean "not clipped".
struct ImGuiTableColumn
{
    ImGuiTableColumnFlags   Flags;                          // Flags after some patching (not directly same as provided by user). See ImGuiTableColumnFlags_
    float                   WidthGiven;                     // Final/actual width visible == (MaxX - MinX), locked in TableUpdateLayout(). May be > WidthRequest to honor minimum width, may be < WidthRequest to honor shrinking columns down in tight space.
    float                   MinX;                           // Absolute positions
    float                   MaxX;
    float                   WidthRequest;                   // Master width absolute value when !(Flags & _WidthStretch). When Stretch this is derived every frame from StretchWeight in TableUpdateLayout()
    float                   WidthAuto;                      // Automatic width
    float                   StretchWeight;                  // Master width weight when (Flags & _WidthStretch). Often around ~1.0f initially.
    float                   InitStretchWeightOrWidth;       // Value passed to TableSetupColumn(). For Width it is a content width (_without padding_).
    ImRect                  ClipRect;                       // Clipping rectangle for the column
    ImGuiID                 UserID;                         // Optional, value passed to TableSetupColumn()
    float                   WorkMinX;                       // Contents region min ~(MinX + CellPaddingX + CellSpacingX1) == cursor start position when entering column
    float                   WorkMaxX;                       // Contents region max ~(MaxX - CellPaddingX - CellSpacingX2)
    float                   ItemWidth;                      // Current item width for the column, preserved across rows
    float                   ContentMaxXFrozen;              // Contents maximum position for frozen rows (apart from headers), from which we can infer content width.
    float                   ContentMaxXUnfrozen;
    float                   ContentMaxXHeadersUsed;         // Contents maximum position for headers rows (regardless of freezing). TableHeader() automatically softclip itself + report ideal desired size, to avoid creating extraneous draw calls
    float                   ContentMaxXHeadersIdeal;
    ImS16                   NameOffset;                     // Offset into parent ColumnsNames[]
    ImGuiTableColumnIdx     DisplayOrder;                   // Index within Table's IndexToDisplayOrder[] (column may be reordered by users)
    ImGuiTableColumnIdx     IndexWithinEnabledSet;          // Index within enabled/visible set (<= IndexToDisplayOrder)
    ImGuiTableColumnIdx     PrevEnabledColumn;              // Index of prev enabled/visible column within Columns[], -1 if first enabled/visible column
    ImGuiTableColumnIdx     NextEnabledColumn;              // Index of next enabled/visible column within Columns[], -1 if last enabled/visible column
    ImGuiTableColumnIdx     SortOrder;                      // Index of this column within sort specs, -1 if not sorting on this column, 0 for single-sort, may be >0 on multi-sort
    ImGuiTableDrawChannelIdx DrawChannelCurrent;            // Index within DrawSplitter.Channels[]
    ImGuiTableDrawChannelIdx DrawChannelFrozen;             // Draw channels for frozen rows (often headers)
    ImGuiTableDrawChannelIdx DrawChannelUnfrozen;           // Draw channels for unfrozen rows
    bool                    IsEnabled;                      // IsUserEnabled && (Flags & ImGuiTableColumnFlags_Disabled) == 0
    bool                    IsUserEnabled;                  // Is the column not marked Hidden by the user? (unrelated to being off view, e.g. clipped by scrolling).
    bool                    IsUserEnabledNextFrame;
    bool                    IsVisibleX;                     // Is actually in view (e.g. overlapping the host window clipping rectangle, not scrolled).
    bool                    IsVisibleY;
    bool                    IsRequestOutput;                // Return value for TableSetColumnIndex() / TableNextColumn(): whether we request user to output contents or not.
    bool                    IsSkipItems;                    // Do we want item submissions to this column to be completely ignored (no layout will happen).
    bool                    IsPreserveWidthAuto;
    ImS8                    NavLayerCurrent;                // ImGuiNavLayer in 1 byte
    ImU8                    AutoFitQueue;                   // Queue of 8 values for the next 8 frames to request auto-fit
    ImU8                    CannotSkipItemsQueue;           // Queue of 8 values for the next 8 frames to disable Clipped/SkipItem
    ImU8                    SortDirection : 2;              // ImGuiSortDirection_Ascending or ImGuiSortDirection_Descending
    ImU8                    SortDirectionsAvailCount : 2;   // Number of available sort directions (0 to 3)
    ImU8                    SortDirectionsAvailMask : 4;    // Mask of available sort directions (1-bit each)
    ImU8                    SortDirectionsAvailList;        // Ordered of available sort directions (2-bits each)

    ImGuiTableColumn()
    {
        memset(this, 0, sizeof(*this));
        StretchWeight = WidthRequest = -1.0f;
        NameOffset = -1;
        DisplayOrder = IndexWithinEnabledSet = -1;
        PrevEnabledColumn = NextEnabledColumn = -1;
        SortOrder = -1;
        SortDirection = ImGuiSortDirection_None;
        DrawChannelCurrent = DrawChannelFrozen = DrawChannelUnfrozen = (ImU8)-1;
    }
};

// Transient cell data stored per row.
// sizeof() ~ 6
struct ImGuiTableCellData
{
    ImU32                       BgColor;    // Actual color
    ImGuiTableColumnIdx         Column;     // Column number
};

// FIXME-TABLE: more transient data could be stored in a per-stacked table structure: DrawSplitter, SortSpecs, incoming RowData
struct ImGuiTable
{
    ImGuiID                     ID;
    ImGuiTableFlags             Flags;
    void*                       RawData;                    // Single allocation to hold Columns[], DisplayOrderToIndex[] and RowCellData[]
    ImGuiTableTempData*         TempData;                   // Transient data while table is active. Point within g.CurrentTableStack[]
    ImSpan<ImGuiTableColumn>    Columns;                    // Point within RawData[]
    ImSpan<ImGuiTableColumnIdx> DisplayOrderToIndex;        // Point within RawData[]. Store display order of columns (when not reordered, the values are 0...Count-1)
    ImSpan<ImGuiTableCellData>  RowCellData;                // Point within RawData[]. Store cells background requests for current row.
    ImU64                       EnabledMaskByDisplayOrder;  // Column DisplayOrder -> IsEnabled map
    ImU64                       EnabledMaskByIndex;         // Column Index -> IsEnabled map (== not hidden by user/api) in a format adequate for iterating column without touching cold data
    ImU64                       VisibleMaskByIndex;         // Column Index -> IsVisibleX|IsVisibleY map (== not hidden by user/api && not hidden by scrolling/cliprect)
    ImU64                       RequestOutputMaskByIndex;   // Column Index -> IsVisible || AutoFit (== expect user to submit items)
    ImGuiTableFlags             SettingsLoadedFlags;        // Which data were loaded from the .ini file (e.g. when order is not altered we won't save order)
    int                         SettingsOffset;             // Offset in g.SettingsTables
    int                         LastFrameActive;
    int                         ColumnsCount;               // Number of columns declared in BeginTable()
    int                         CurrentRow;
    int                         CurrentColumn;
    ImS16                       InstanceCurrent;            // Count of BeginTable() calls with same ID in the same frame (generally 0). This is a little bit similar to BeginCount for a window, but multiple table with same ID look are multiple tables, they are just synched.
    ImS16                       InstanceInteracted;         // Mark which instance (generally 0) of the same ID is being interacted with
    float                       RowPosY1;
    float                       RowPosY2;
    float                       RowMinHeight;               // Height submitted to TableNextRow()
    float                       RowTextBaseline;
    float                       RowIndentOffsetX;
    ImGuiTableRowFlags          RowFlags : 16;              // Current row flags, see ImGuiTableRowFlags_
    ImGuiTableRowFlags          LastRowFlags : 16;
    int                         RowBgColorCounter;          // Counter for alternating background colors (can be fast-forwarded by e.g clipper), not same as CurrentRow because header rows typically don't increase this.
    ImU32                       RowBgColor[2];              // Background color override for current row.
    ImU32                       BorderColorStrong;
    ImU32                       BorderColorLight;
    float                       BorderX1;
    float                       BorderX2;
    float                       HostIndentX;
    float                       MinColumnWidth;
    float                       OuterPaddingX;
    float                       CellPaddingX;               // Padding from each borders
    float                       CellPaddingY;
    float                       CellSpacingX1;              // Spacing between non-bordered cells
    float                       CellSpacingX2;
    float                       LastOuterHeight;            // Outer height from last frame
    float                       LastFirstRowHeight;         // Height of first row from last frame
    float                       InnerWidth;                 // User value passed to BeginTable(), see comments at the top of BeginTable() for details.
    float                       ColumnsGivenWidth;          // Sum of current column width
    float                       ColumnsAutoFitWidth;        // Sum of ideal column width in order nothing to be clipped, used for auto-fitting and content width submission in outer window
    float                       ResizedColumnNextWidth;
    float                       ResizeLockMinContentsX2;    // Lock minimum contents width while resizing down in order to not create feedback loops. But we allow growing the table.
    float                       RefScale;                   // Reference scale to be able to rescale columns on font/dpi changes.
    ImRect                      OuterRect;                  // Note: for non-scrolling table, OuterRect.Max.y is often FLT_MAX until EndTable(), unless a height has been specified in BeginTable().
    ImRect                      InnerRect;                  // InnerRect but without decoration. As with OuterRect, for non-scrolling tables, InnerRect.Max.y is
    ImRect                      WorkRect;
    ImRect                      InnerClipRect;
    ImRect                      BgClipRect;                 // We use this to cpu-clip cell background color fill
    ImRect                      Bg0ClipRectForDrawCmd;      // Actual ImDrawCmd clip rect for BG0/1 channel. This tends to be == OuterWindow->ClipRect at BeginTable() because output in BG0/BG1 is cpu-clipped
    ImRect                      Bg2ClipRectForDrawCmd;      // Actual ImDrawCmd clip rect for BG2 channel. This tends to be a correct, tight-fit, because output to BG2 are done by widgets relying on regular ClipRect.
    ImRect                      HostClipRect;               // This is used to check if we can eventually merge our columns draw calls into the current draw call of the current window.
    ImRect                      HostBackupInnerClipRect;    // Backup of InnerWindow->ClipRect during PushTableBackground()/PopTableBackground()
    ImGuiWindow*                OuterWindow;                // Parent window for the table
    ImGuiWindow*                InnerWindow;                // Window holding the table data (== OuterWindow or a child window)
    ImGuiTextBuffer             ColumnsNames;               // Contiguous buffer holding columns names
    ImDrawListSplitter*         DrawSplitter;               // Shortcut to TempData->DrawSplitter while in table. Isolate draw commands per columns to avoid switching clip rect constantly
    ImGuiTableColumnSortSpecs   SortSpecsSingle;
    ImVector<ImGuiTableColumnSortSpecs> SortSpecsMulti;     // FIXME-OPT: Using a small-vector pattern would be good.
    ImGuiTableSortSpecs         SortSpecs;                  // Public facing sorts specs, this is what we return in TableGetSortSpecs()
    ImGuiTableColumnIdx         SortSpecsCount;
    ImGuiTableColumnIdx         ColumnsEnabledCount;        // Number of enabled columns (<= ColumnsCount)
    ImGuiTableColumnIdx         ColumnsEnabledFixedCount;   // Number of enabled columns (<= ColumnsCount)
    ImGuiTableColumnIdx         DeclColumnsCount;           // Count calls to TableSetupColumn()
    ImGuiTableColumnIdx         HoveredColumnBody;          // Index of column whose visible region is being hovered. Important: == ColumnsCount when hovering empty region after the right-most column!
    ImGuiTableColumnIdx         HoveredColumnBorder;        // Index of column whose right-border is being hovered (for resizing).
    ImGuiTableColumnIdx         AutoFitSingleColumn;        // Index of single column requesting auto-fit.
    ImGuiTableColumnIdx         ResizedColumn;              // Index of column being resized. Reset when InstanceCurrent==0.
    ImGuiTableColumnIdx         LastResizedColumn;          // Index of column being resized from previous frame.
    ImGuiTableColumnIdx         HeldHeaderColumn;           // Index of column header being held.
    ImGuiTableColumnIdx         ReorderColumn;              // Index of column being reordered. (not cleared)
    ImGuiTableColumnIdx         ReorderColumnDir;           // -1 or +1
    ImGuiTableColumnIdx         LeftMostEnabledColumn;      // Index of left-most non-hidden column.
    ImGuiTableColumnIdx         RightMostEnabledColumn;     // Index of right-most non-hidden column.
    ImGuiTableColumnIdx         LeftMostStretchedColumn;    // Index of left-most stretched column.
    ImGuiTableColumnIdx         RightMostStretchedColumn;   // Index of right-most stretched column.
    ImGuiTableColumnIdx         ContextPopupColumn;         // Column right-clicked on, of -1 if opening context menu from a neutral/empty spot
    ImGuiTableColumnIdx         FreezeRowsRequest;          // Requested frozen rows count
    ImGuiTableColumnIdx         FreezeRowsCount;            // Actual frozen row count (== FreezeRowsRequest, or == 0 when no scrolling offset)
    ImGuiTableColumnIdx         FreezeColumnsRequest;       // Requested frozen columns count
    ImGuiTableColumnIdx         FreezeColumnsCount;         // Actual frozen columns count (== FreezeColumnsRequest, or == 0 when no scrolling offset)
    ImGuiTableColumnIdx         RowCellDataCurrent;         // Index of current RowCellData[] entry in current row
    ImGuiTableDrawChannelIdx    DummyDrawChannel;           // Redirect non-visible columns here.
    ImGuiTableDrawChannelIdx    Bg2DrawChannelCurrent;      // For Selectable() and other widgets drawing across columns after the freezing line. Index within DrawSplitter.Channels[]
    ImGuiTableDrawChannelIdx    Bg2DrawChannelUnfrozen;
    bool                        IsLayoutLocked;             // Set by TableUpdateLayout() which is called when beginning the first row.
    bool                        IsInsideRow;                // Set when inside TableBeginRow()/TableEndRow().
    bool                        IsInitializing;
    bool                        IsSortSpecsDirty;
    bool                        IsUsingHeaders;             // Set when the first row had the ImGuiTableRowFlags_Headers flag.
    bool                        IsContextPopupOpen;         // Set when default context menu is open (also see: ContextPopupColumn, InstanceInteracted).
    bool                        IsSettingsRequestLoad;
    bool                        IsSettingsDirty;            // Set when table settings have changed and needs to be reported into ImGuiTableSetttings data.
    bool                        IsDefaultDisplayOrder;      // Set when display order is unchanged from default (DisplayOrder contains 0...Count-1)
    bool                        IsResetAllRequest;
    bool                        IsResetDisplayOrderRequest;
    bool                        IsUnfrozenRows;             // Set when we got past the frozen row.
    bool                        IsDefaultSizingPolicy;      // Set if user didn't explicitly set a sizing policy in BeginTable()
    bool                        MemoryCompacted;
    bool                        HostSkipItems;              // Backup of InnerWindow->SkipItem at the end of BeginTable(), because we will overwrite InnerWindow->SkipItem on a per-column basis

    IMGUI_API ImGuiTable()      { memset(this, 0, sizeof(*this)); LastFrameActive = -1; }
    IMGUI_API ~ImGuiTable()     { IM_FREE(RawData); }
};

// Transient data that are only needed between BeginTable() and EndTable(), those buffers are shared (1 per level of stacked table).
// - Accessing those requires chasing an extra pointer so for very frequently used data we leave them in the main table structure.
// - We also leave out of this structure data that tend to be particularly useful for debugging/metrics.
struct ImGuiTableTempData
{
    int                         TableIndex;                 // Index in g.Tables.Buf[] pool
    float                       LastTimeActive;             // Last timestamp this structure was used

    ImVec2                      UserOuterSize;              // outer_size.x passed to BeginTable()
    ImDrawListSplitter          DrawSplitter;

    ImRect                      HostBackupWorkRect;         // Backup of InnerWindow->WorkRect at the end of BeginTable()
    ImRect                      HostBackupParentWorkRect;   // Backup of InnerWindow->ParentWorkRect at the end of BeginTable()
    ImVec2                      HostBackupPrevLineSize;     // Backup of InnerWindow->DC.PrevLineSize at the end of BeginTable()
    ImVec2                      HostBackupCurrLineSize;     // Backup of InnerWindow->DC.CurrLineSize at the end of BeginTable()
    ImVec2                      HostBackupCursorMaxPos;     // Backup of InnerWindow->DC.CursorMaxPos at the end of BeginTable()
    ImVec1                      HostBackupColumnsOffset;    // Backup of OuterWindow->DC.ColumnsOffset at the end of BeginTable()
    float                       HostBackupItemWidth;        // Backup of OuterWindow->DC.ItemWidth at the end of BeginTable()
    int                         HostBackupItemWidthStackSize;//Backup of OuterWindow->DC.ItemWidthStack.Size at the end of BeginTable()

    IMGUI_API ImGuiTableTempData() { memset(this, 0, sizeof(*this)); LastTimeActive = -1.0f; }
};

// sizeof() ~ 12
struct ImGuiTableColumnSettings
{
    float                   WidthOrWeight;
    ImGuiID                 UserID;
    ImGuiTableColumnIdx     Index;
    ImGuiTableColumnIdx     DisplayOrder;
    ImGuiTableColumnIdx     SortOrder;
    ImU8                    SortDirection : 2;
    ImU8                    IsEnabled : 1; // "Visible" in ini file
    ImU8                    IsStretch : 1;

    ImGuiTableColumnSettings()
    {
        WidthOrWeight = 0.0f;
        UserID = 0;
        Index = -1;
        DisplayOrder = SortOrder = -1;
        SortDirection = ImGuiSortDirection_None;
        IsEnabled = 1;
        IsStretch = 0;
    }
};

// This is designed to be stored in a single ImChunkStream (1 header followed by N ImGuiTableColumnSettings, etc.)
struct ImGuiTableSettings
{
    ImGuiID                     ID;                     // Set to 0 to invalidate/delete the setting
    ImGuiTableFlags             SaveFlags;              // Indicate data we want to save using the Resizable/Reorderable/Sortable/Hideable flags (could be using its own flags..)
    float                       RefScale;               // Reference scale to be able to rescale columns on font/dpi changes.
    ImGuiTableColumnIdx         ColumnsCount;
    ImGuiTableColumnIdx         ColumnsCountMax;        // Maximum number of columns this settings instance can store, we can recycle a settings instance with lower number of columns but not higher
    bool                        WantApply;              // Set when loaded from .ini data (to enable merging/loading .ini data into an already running context)

    ImGuiTableSettings()        { memset(this, 0, sizeof(*this)); }
    ImGuiTableColumnSettings*   GetColumnSettings()     { return (ImGuiTableColumnSettings*)(this + 1); }
};

//-----------------------------------------------------------------------------
// [SECTION] ImGui internal API
// No guarantee of forward compatibility here!
//-----------------------------------------------------------------------------

namespace ImGui
{
    // Windows
    // We should always have a CurrentWindow in the stack (there is an implicit "Debug" window)
    // If this ever crash because g.CurrentWindow is NULL it means that either
    // - ImGui::NewFrame() has never been called, which is illegal.
    // - You are calling ImGui functions after ImGui::EndFrame()/ImGui::Render() and before the next ImGui::NewFrame(), which is also illegal.
    inline    ImGuiWindow*  GetCurrentWindowRead()      { ImGuiContext& g = *GImGui; return g.CurrentWindow; }
    inline    ImGuiWindow*  GetCurrentWindow()          { ImGuiContext& g = *GImGui; g.CurrentWindow->WriteAccessed = true; return g.CurrentWindow; }
    IMGUI_API ImGuiWindow*  FindWindowByID(ImGuiID id);
    IMGUI_API ImGuiWindow*  FindWindowByName(const char* name);
    IMGUI_API void          UpdateWindowParentAndRootLinks(ImGuiWindow* window, ImGuiWindowFlags flags, ImGuiWindow* parent_window);
    IMGUI_API ImVec2        CalcWindowNextAutoFitSize(ImGuiWindow* window);
    IMGUI_API bool          IsWindowChildOf(ImGuiWindow* window, ImGuiWindow* potential_parent);
    IMGUI_API bool          IsWindowAbove(ImGuiWindow* potential_above, ImGuiWindow* potential_below);
    IMGUI_API bool          IsWindowNavFocusable(ImGuiWindow* window);
    IMGUI_API void          SetWindowPos(ImGuiWindow* window, const ImVec2& pos, ImGuiCond cond = 0);
    IMGUI_API void          SetWindowSize(ImGuiWindow* window, const ImVec2& size, ImGuiCond cond = 0);
    IMGUI_API void          SetWindowCollapsed(ImGuiWindow* window, bool collapsed, ImGuiCond cond = 0);
    IMGUI_API void          SetWindowHitTestHole(ImGuiWindow* window, const ImVec2& pos, const ImVec2& size);

    // Windows: Display Order and Focus Order
    IMGUI_API void          FocusWindow(ImGuiWindow* window);
    IMGUI_API void          FocusTopMostWindowUnderOne(ImGuiWindow* under_this_window, ImGuiWindow* ignore_window);
    IMGUI_API void          BringWindowToFocusFront(ImGuiWindow* window);
    IMGUI_API void          BringWindowToDisplayFront(ImGuiWindow* window);
    IMGUI_API void          BringWindowToDisplayBack(ImGuiWindow* window);

    // Fonts, drawing
    IMGUI_API void          SetCurrentFont(ImFont* font);
    inline ImFont*          GetDefaultFont() { ImGuiContext& g = *GImGui; return g.IO.FontDefault ? g.IO.FontDefault : g.IO.Fonts->Fonts[0]; }
    inline ImDrawList*      GetForegroundDrawList(ImGuiWindow* window) { IM_UNUSED(window); return GetForegroundDrawList(); } // This seemingly unnecessary wrapper simplifies compatibility between the 'master' and 'docking' branches.
    IMGUI_API ImDrawList*   GetBackgroundDrawList(ImGuiViewport* viewport);                     // get background draw list for the given viewport. this draw list will be the first rendering one. Useful to quickly draw shapes/text behind dear imgui contents.
    IMGUI_API ImDrawList*   GetForegroundDrawList(ImGuiViewport* viewport);                     // get foreground draw list for the given viewport. this draw list will be the last rendered one. Useful to quickly draw shapes/text over dear imgui contents.

    // Init
    IMGUI_API void          Initialize(ImGuiContext* context);
    IMGUI_API void          Shutdown(ImGuiContext* context);    // Since 1.60 this is a _private_ function. You can call DestroyContext() to destroy the context created by CreateContext().

    // NewFrame
    IMGUI_API void          UpdateHoveredWindowAndCaptureFlags();
    IMGUI_API void          StartMouseMovingWindow(ImGuiWindow* window);
    IMGUI_API void          UpdateMouseMovingWindowNewFrame();
    IMGUI_API void          UpdateMouseMovingWindowEndFrame();

    // Generic context hooks
    IMGUI_API ImGuiID       AddContextHook(ImGuiContext* context, const ImGuiContextHook* hook);
    IMGUI_API void          RemoveContextHook(ImGuiContext* context, ImGuiID hook_to_remove);
    IMGUI_API void          CallContextHooks(ImGuiContext* context, ImGuiContextHookType type);

    // Settings
    IMGUI_API void                  MarkIniSettingsDirty();
    IMGUI_API void                  MarkIniSettingsDirty(ImGuiWindow* window);
    IMGUI_API void                  ClearIniSettings();
    IMGUI_API ImGuiWindowSettings*  CreateNewWindowSettings(const char* name);
    IMGUI_API ImGuiWindowSettings*  FindWindowSettings(ImGuiID id);
    IMGUI_API ImGuiWindowSettings*  FindOrCreateWindowSettings(const char* name);
    IMGUI_API ImGuiSettingsHandler* FindSettingsHandler(const char* type_name);

    // Scrolling
    IMGUI_API void          SetNextWindowScroll(const ImVec2& scroll); // Use -1.0f on one axis to leave as-is
    IMGUI_API void          SetScrollX(ImGuiWindow* window, float scroll_x);
    IMGUI_API void          SetScrollY(ImGuiWindow* window, float scroll_y);
    IMGUI_API void          SetScrollFromPosX(ImGuiWindow* window, float local_x, float center_x_ratio);
    IMGUI_API void          SetScrollFromPosY(ImGuiWindow* window, float local_y, float center_y_ratio);
    IMGUI_API ImVec2        ScrollToBringRectIntoView(ImGuiWindow* window, const ImRect& item_rect);

    // Basic Accessors
    inline ImGuiID          GetItemID()     { ImGuiContext& g = *GImGui; return g.LastItemData.ID; }   // Get ID of last item (~~ often same ImGui::GetID(label) beforehand)
    inline ImGuiItemStatusFlags GetItemStatusFlags(){ ImGuiContext& g = *GImGui; return g.LastItemData.StatusFlags; }
    inline ImGuiItemFlags   GetItemFlags()  { ImGuiContext& g = *GImGui; return g.LastItemData.InFlags; }
    inline ImGuiID          GetActiveID()   { ImGuiContext& g = *GImGui; return g.ActiveId; }
    inline ImGuiID          GetFocusID()    { ImGuiContext& g = *GImGui; return g.NavId; }
    IMGUI_API void          SetActiveID(ImGuiID id, ImGuiWindow* window);
    IMGUI_API void          SetFocusID(ImGuiID id, ImGuiWindow* window);
    IMGUI_API void          ClearActiveID();
    IMGUI_API ImGuiID       GetHoveredID();
    IMGUI_API void          SetHoveredID(ImGuiID id);
    IMGUI_API void          KeepAliveID(ImGuiID id);
    IMGUI_API void          MarkItemEdited(ImGuiID id);     // Mark data associated to given item as "edited", used by IsItemDeactivatedAfterEdit() function.
    IMGUI_API void          PushOverrideID(ImGuiID id);     // Push given value as-is at the top of the ID stack (whereas PushID combines old and new hashes)
    IMGUI_API ImGuiID       GetIDWithSeed(const char* str_id_begin, const char* str_id_end, ImGuiID seed);

    // Basic Helpers for widget code
    IMGUI_API void          ItemSize(const ImVec2& size, float text_baseline_y = -1.0f);
    IMGUI_API void          ItemSize(const ImRect& bb, float text_baseline_y = -1.0f);
    IMGUI_API bool          ItemAdd(const ImRect& bb, ImGuiID id, const ImRect* nav_bb = NULL, ImGuiItemFlags extra_flags = 0);
    IMGUI_API bool          ItemHoverable(const ImRect& bb, ImGuiID id);
    IMGUI_API void          ItemInputable(ImGuiWindow* window, ImGuiID id);
    IMGUI_API bool          IsClippedEx(const ImRect& bb, ImGuiID id, bool clip_even_when_logged);
    IMGUI_API ImVec2        CalcItemSize(ImVec2 size, float default_w, float default_h);
    IMGUI_API float         CalcWrapWidthForPos(const ImVec2& pos, float wrap_pos_x);
    IMGUI_API void          PushMultiItemsWidths(int components, float width_full);
    IMGUI_API bool          IsItemToggledSelection();                                   // Was the last item selection toggled? (after Selectable(), TreeNode() etc. We only returns toggle _event_ in order to handle clipping correctly)
    IMGUI_API ImVec2        GetContentRegionMaxAbs();
    IMGUI_API void          ShrinkWidths(ImGuiShrinkWidthItem* items, int count, float width_excess);

    // Parameter stacks
    IMGUI_API void          PushItemFlag(ImGuiItemFlags option, bool enabled);
    IMGUI_API void          PopItemFlag();

#ifndef IMGUI_DISABLE_OBSOLETE_FUNCTIONS
    // Currently refactoring focus/nav/tabbing system
    // If you have old/custom copy-and-pasted widgets that used FocusableItemRegister():
    //  (Old) IMGUI_VERSION_NUM  < 18209: using 'ItemAdd(....)'                              and 'bool focused = FocusableItemRegister(...)'
    //  (Old) IMGUI_VERSION_NUM >= 18209: using 'ItemAdd(..., ImGuiItemAddFlags_Focusable)'  and 'bool focused = (GetItemStatusFlags() & ImGuiItemStatusFlags_Focused) != 0'
    //  (New) IMGUI_VERSION_NUM >= 18411: using 'ItemAdd(..., ImGuiItemAddFlags_Inputable)'  and 'bool focused = (GetItemStatusFlags() & ImGuiItemStatusFlags_Focused) != 0'
    // Widget code are simplified as there's no need to call FocusableItemUnregister() while managing the transition from regular widget to TempInputText()
    inline bool FocusableItemRegister(ImGuiWindow* window, ImGuiID id)  { IM_ASSERT(0); IM_UNUSED(window); IM_UNUSED(id); return false; } // -> pass ImGuiItemAddFlags_Inputable flag to ItemAdd()
    inline void FocusableItemUnregister(ImGuiWindow* window)            { IM_ASSERT(0); IM_UNUSED(window); }                              // -> unnecessary: TempInputText() uses ImGuiInputTextFlags_MergedItem
#endif

    // Logging/Capture
    IMGUI_API void          LogBegin(ImGuiLogType type, int auto_open_depth);           // -> BeginCapture() when we design v2 api, for now stay under the radar by using the old name.
    IMGUI_API void          LogToBuffer(int auto_open_depth = -1);                      // Start logging/capturing to internal buffer
    IMGUI_API void          LogRenderedText(const ImVec2* ref_pos, const char* text, const char* text_end = NULL);
    IMGUI_API void          LogSetNextTextDecoration(const char* prefix, const char* suffix);

    // Popups, Modals, Tooltips
    IMGUI_API bool          BeginChildEx(const char* name, ImGuiID id, const ImVec2& size_arg, bool border, ImGuiWindowFlags flags);
    IMGUI_API void          OpenPopupEx(ImGuiID id, ImGuiPopupFlags popup_flags = ImGuiPopupFlags_None);
    IMGUI_API void          ClosePopupToLevel(int remaining, bool restore_focus_to_window_under_popup);
    IMGUI_API void          ClosePopupsOverWindow(ImGuiWindow* ref_window, bool restore_focus_to_window_under_popup);
    IMGUI_API bool          IsPopupOpen(ImGuiID id, ImGuiPopupFlags popup_flags);
    IMGUI_API bool          BeginPopupEx(ImGuiID id, ImGuiWindowFlags extra_flags);
    IMGUI_API void          BeginTooltipEx(ImGuiWindowFlags extra_flags, ImGuiTooltipFlags tooltip_flags);
    IMGUI_API ImRect        GetPopupAllowedExtentRect(ImGuiWindow* window);
    IMGUI_API ImGuiWindow*  GetTopMostPopupModal();
    IMGUI_API ImVec2        FindBestWindowPosForPopup(ImGuiWindow* window);
    IMGUI_API ImVec2        FindBestWindowPosForPopupEx(const ImVec2& ref_pos, const ImVec2& size, ImGuiDir* last_dir, const ImRect& r_outer, const ImRect& r_avoid, ImGuiPopupPositionPolicy policy);
    IMGUI_API bool          BeginViewportSideBar(const char* name, ImGuiViewport* viewport, ImGuiDir dir, float size, ImGuiWindowFlags window_flags);

    // Menus
    IMGUI_API bool          BeginMenuEx(const char* label, const char* icon, bool enabled = true);
    IMGUI_API bool          MenuItemEx(const char* label, const char* icon, const char* shortcut = NULL, bool selected = false, bool enabled = true);

    // Combos
    IMGUI_API bool          BeginComboPopup(ImGuiID popup_id, const ImRect& bb, ImGuiComboFlags flags);
    IMGUI_API bool          BeginComboPreview();
    IMGUI_API void          EndComboPreview();

    // Gamepad/Keyboard Navigation
    IMGUI_API void          NavInitWindow(ImGuiWindow* window, bool force_reinit);
    IMGUI_API bool          NavMoveRequestButNoResultYet();
    IMGUI_API void          NavMoveRequestSubmit(ImGuiDir move_dir, ImGuiDir clip_dir, ImGuiNavMoveFlags move_flags);
    IMGUI_API void          NavMoveRequestForward(ImGuiDir move_dir, ImGuiDir clip_dir, ImGuiNavMoveFlags move_flags);
    IMGUI_API void          NavMoveRequestCancel();
    IMGUI_API void          NavMoveRequestApplyResult();
    IMGUI_API void          NavMoveRequestTryWrapping(ImGuiWindow* window, ImGuiNavMoveFlags move_flags);
    IMGUI_API float         GetNavInputAmount(ImGuiNavInput n, ImGuiInputReadMode mode);
    IMGUI_API ImVec2        GetNavInputAmount2d(ImGuiNavDirSourceFlags dir_sources, ImGuiInputReadMode mode, float slow_factor = 0.0f, float fast_factor = 0.0f);
    IMGUI_API int           CalcTypematicRepeatAmount(float t0, float t1, float repeat_delay, float repeat_rate);
    IMGUI_API void          ActivateItem(ImGuiID id);   // Remotely activate a button, checkbox, tree node etc. given its unique ID. activation is queued and processed on the next frame when the item is encountered again.
    IMGUI_API void          SetNavID(ImGuiID id, ImGuiNavLayer nav_layer, ImGuiID focus_scope_id, const ImRect& rect_rel);

    // Focus Scope (WIP)
    // This is generally used to identify a selection set (multiple of which may be in the same window), as selection
    // patterns generally need to react (e.g. clear selection) when landing on an item of the set.
    IMGUI_API void          PushFocusScope(ImGuiID id);
    IMGUI_API void          PopFocusScope();
    inline ImGuiID          GetFocusedFocusScope()          { ImGuiContext& g = *GImGui; return g.NavFocusScopeId; }                            // Focus scope which is actually active
    inline ImGuiID          GetFocusScope()                 { ImGuiContext& g = *GImGui; return g.CurrentWindow->DC.NavFocusScopeIdCurrent; }   // Focus scope we are outputting into, set by PushFocusScope()

    // Inputs
    // FIXME: Eventually we should aim to move e.g. IsActiveIdUsingKey() into IsKeyXXX functions.
    IMGUI_API void          SetItemUsingMouseWheel();
    IMGUI_API void          SetActiveIdUsingNavAndKeys();
    inline bool             IsActiveIdUsingNavDir(ImGuiDir dir)                         { ImGuiContext& g = *GImGui; return (g.ActiveIdUsingNavDirMask & (1 << dir)) != 0; }
    inline bool             IsActiveIdUsingNavInput(ImGuiNavInput input)                { ImGuiContext& g = *GImGui; return (g.ActiveIdUsingNavInputMask & (1 << input)) != 0; }
    inline bool             IsActiveIdUsingKey(ImGuiKey key)                            { ImGuiContext& g = *GImGui; IM_ASSERT(key < 64); return (g.ActiveIdUsingKeyInputMask & ((ImU64)1 << key)) != 0; }
    IMGUI_API bool          IsMouseDragPastThreshold(ImGuiMouseButton button, float lock_threshold = -1.0f);
    inline bool             IsKeyPressedMap(ImGuiKey key, bool repeat = true)           { ImGuiContext& g = *GImGui; const int key_index = g.IO.KeyMap[key]; return (key_index >= 0) ? IsKeyPressed(key_index, repeat) : false; }
    inline bool             IsNavInputDown(ImGuiNavInput n)                             { ImGuiContext& g = *GImGui; return g.IO.NavInputs[n] > 0.0f; }
    inline bool             IsNavInputTest(ImGuiNavInput n, ImGuiInputReadMode rm)      { return (GetNavInputAmount(n, rm) > 0.0f); }
    IMGUI_API ImGuiKeyModFlags GetMergedKeyModFlags();

    // Drag and Drop
    IMGUI_API bool          BeginDragDropTargetCustom(const ImRect& bb, ImGuiID id);
    IMGUI_API void          ClearDragDrop();
    IMGUI_API bool          IsDragDropPayloadBeingAccepted();

    // Internal Columns API (this is not exposed because we will encourage transitioning to the Tables API)
    IMGUI_API void          SetWindowClipRectBeforeSetChannel(ImGuiWindow* window, const ImRect& clip_rect);
    IMGUI_API void          BeginColumns(const char* str_id, int count, ImGuiOldColumnFlags flags = 0); // setup number of columns. use an identifier to distinguish multiple column sets. close with EndColumns().
    IMGUI_API void          EndColumns();                                                               // close columns
    IMGUI_API void          PushColumnClipRect(int column_index);
    IMGUI_API void          PushColumnsBackground();
    IMGUI_API void          PopColumnsBackground();
    IMGUI_API ImGuiID       GetColumnsID(const char* str_id, int count);
    IMGUI_API ImGuiOldColumns* FindOrCreateColumns(ImGuiWindow* window, ImGuiID id);
    IMGUI_API float         GetColumnOffsetFromNorm(const ImGuiOldColumns* columns, float offset_norm);
    IMGUI_API float         GetColumnNormFromOffset(const ImGuiOldColumns* columns, float offset);

    // Tables: Candidates for public API
    IMGUI_API void          TableOpenContextMenu(int column_n = -1);
    IMGUI_API void          TableSetColumnWidth(int column_n, float width);
    IMGUI_API void          TableSetColumnSortDirection(int column_n, ImGuiSortDirection sort_direction, bool append_to_sort_specs);
    IMGUI_API int           TableGetHoveredColumn(); // May use (TableGetColumnFlags() & ImGuiTableColumnFlags_IsHovered) instead. Return hovered column. return -1 when table is not hovered. return columns_count if the unused space at the right of visible columns is hovered.
    IMGUI_API float         TableGetHeaderRowHeight();
    IMGUI_API void          TablePushBackgroundChannel();
    IMGUI_API void          TablePopBackgroundChannel();

    // Tables: Internals
    inline    ImGuiTable*   GetCurrentTable() { ImGuiContext& g = *GImGui; return g.CurrentTable; }
    IMGUI_API ImGuiTable*   TableFindByID(ImGuiID id);
    IMGUI_API bool          BeginTableEx(const char* name, ImGuiID id, int columns_count, ImGuiTableFlags flags = 0, const ImVec2& outer_size = ImVec2(0, 0), float inner_width = 0.0f);
    IMGUI_API void          TableBeginInitMemory(ImGuiTable* table, int columns_count);
    IMGUI_API void          TableBeginApplyRequests(ImGuiTable* table);
    IMGUI_API void          TableSetupDrawChannels(ImGuiTable* table);
    IMGUI_API void          TableUpdateLayout(ImGuiTable* table);
    IMGUI_API void          TableUpdateBorders(ImGuiTable* table);
    IMGUI_API void          TableUpdateColumnsWeightFromWidth(ImGuiTable* table);
    IMGUI_API void          TableDrawBorders(ImGuiTable* table);
    IMGUI_API void          TableDrawContextMenu(ImGuiTable* table);
    IMGUI_API void          TableMergeDrawChannels(ImGuiTable* table);
    IMGUI_API void          TableSortSpecsSanitize(ImGuiTable* table);
    IMGUI_API void          TableSortSpecsBuild(ImGuiTable* table);
    IMGUI_API ImGuiSortDirection TableGetColumnNextSortDirection(ImGuiTableColumn* column);
    IMGUI_API void          TableFixColumnSortDirection(ImGuiTable* table, ImGuiTableColumn* column);
    IMGUI_API float         TableGetColumnWidthAuto(ImGuiTable* table, ImGuiTableColumn* column);
    IMGUI_API void          TableBeginRow(ImGuiTable* table);
    IMGUI_API void          TableEndRow(ImGuiTable* table);
    IMGUI_API void          TableBeginCell(ImGuiTable* table, int column_n);
    IMGUI_API void          TableEndCell(ImGuiTable* table);
    IMGUI_API ImRect        TableGetCellBgRect(const ImGuiTable* table, int column_n);
    IMGUI_API const char*   TableGetColumnName(const ImGuiTable* table, int column_n);
    IMGUI_API ImGuiID       TableGetColumnResizeID(const ImGuiTable* table, int column_n, int instance_no = 0);
    IMGUI_API float         TableGetMaxColumnWidth(const ImGuiTable* table, int column_n);
    IMGUI_API void          TableSetColumnWidthAutoSingle(ImGuiTable* table, int column_n);
    IMGUI_API void          TableSetColumnWidthAutoAll(ImGuiTable* table);
    IMGUI_API void          TableRemove(ImGuiTable* table);
    IMGUI_API void          TableGcCompactTransientBuffers(ImGuiTable* table);
    IMGUI_API void          TableGcCompactTransientBuffers(ImGuiTableTempData* table);
    IMGUI_API void          TableGcCompactSettings();

    // Tables: Settings
    IMGUI_API void                  TableLoadSettings(ImGuiTable* table);
    IMGUI_API void                  TableSaveSettings(ImGuiTable* table);
    IMGUI_API void                  TableResetSettings(ImGuiTable* table);
    IMGUI_API ImGuiTableSettings*   TableGetBoundSettings(ImGuiTable* table);
    IMGUI_API void                  TableSettingsInstallHandler(ImGuiContext* context);
    IMGUI_API ImGuiTableSettings*   TableSettingsCreate(ImGuiID id, int columns_count);
    IMGUI_API ImGuiTableSettings*   TableSettingsFindByID(ImGuiID id);

    // Tab Bars
    IMGUI_API bool          BeginTabBarEx(ImGuiTabBar* tab_bar, const ImRect& bb, ImGuiTabBarFlags flags);
    IMGUI_API ImGuiTabItem* TabBarFindTabByID(ImGuiTabBar* tab_bar, ImGuiID tab_id);
    IMGUI_API void          TabBarRemoveTab(ImGuiTabBar* tab_bar, ImGuiID tab_id);
    IMGUI_API void          TabBarCloseTab(ImGuiTabBar* tab_bar, ImGuiTabItem* tab);
    IMGUI_API void          TabBarQueueReorder(ImGuiTabBar* tab_bar, const ImGuiTabItem* tab, int offset);
    IMGUI_API void          TabBarQueueReorderFromMousePos(ImGuiTabBar* tab_bar, const ImGuiTabItem* tab, ImVec2 mouse_pos);
    IMGUI_API bool          TabBarProcessReorder(ImGuiTabBar* tab_bar);
    IMGUI_API bool          TabItemEx(ImGuiTabBar* tab_bar, const char* label, bool* p_open, ImGuiTabItemFlags flags);
    IMGUI_API ImVec2        TabItemCalcSize(const char* label, bool has_close_button);
    IMGUI_API void          TabItemBackground(ImDrawList* draw_list, const ImRect& bb, ImGuiTabItemFlags flags, ImU32 col);
    IMGUI_API void          TabItemLabelAndCloseButton(ImDrawList* draw_list, const ImRect& bb, ImGuiTabItemFlags flags, ImVec2 frame_padding, const char* label, ImGuiID tab_id, ImGuiID close_button_id, bool is_contents_visible, bool* out_just_closed, bool* out_text_clipped);

    // Render helpers
    // AVOID USING OUTSIDE OF IMGUI.CPP! NOT FOR PUBLIC CONSUMPTION. THOSE FUNCTIONS ARE A MESS. THEIR SIGNATURE AND BEHAVIOR WILL CHANGE, THEY NEED TO BE REFACTORED INTO SOMETHING DECENT.
    // NB: All position are in absolute pixels coordinates (we are never using window coordinates internally)
    IMGUI_API void          RenderText(ImVec2 pos, const char* text, const char* text_end = NULL, bool hide_text_after_hash = true);
    IMGUI_API void          RenderTextWrapped(ImVec2 pos, const char* text, const char* text_end, float wrap_width);
    IMGUI_API void          RenderTextClipped(const ImVec2& pos_min, const ImVec2& pos_max, const char* text, const char* text_end, const ImVec2* text_size_if_known, const ImVec2& align = ImVec2(0, 0), const ImRect* clip_rect = NULL);
    IMGUI_API void          RenderTextClippedEx(ImDrawList* draw_list, const ImVec2& pos_min, const ImVec2& pos_max, const char* text, const char* text_end, const ImVec2* text_size_if_known, const ImVec2& align = ImVec2(0, 0), const ImRect* clip_rect = NULL);
    IMGUI_API void          RenderTextEllipsis(ImDrawList* draw_list, const ImVec2& pos_min, const ImVec2& pos_max, float clip_max_x, float ellipsis_max_x, const char* text, const char* text_end, const ImVec2* text_size_if_known);
    IMGUI_API void          RenderFrame(ImVec2 p_min, ImVec2 p_max, ImU32 fill_col, bool border = true, float rounding = 0.0f);
    IMGUI_API void          RenderFrameBorder(ImVec2 p_min, ImVec2 p_max, float rounding = 0.0f);
    IMGUI_API void          RenderColorRectWithAlphaCheckerboard(ImDrawList* draw_list, ImVec2 p_min, ImVec2 p_max, ImU32 fill_col, float grid_step, ImVec2 grid_off, float rounding = 0.0f, ImDrawFlags flags = 0);
    IMGUI_API void          RenderNavHighlight(const ImRect& bb, ImGuiID id, ImGuiNavHighlightFlags flags = ImGuiNavHighlightFlags_TypeDefault); // Navigation highlight
    IMGUI_API const char*   FindRenderedTextEnd(const char* text, const char* text_end = NULL); // Find the optional ## from which we stop displaying text.

    // Render helpers (those functions don't access any ImGui state!)
    IMGUI_API void          RenderArrow(ImDrawList* draw_list, ImVec2 pos, ImU32 col, ImGuiDir dir, float scale = 1.0f);
    IMGUI_API void          RenderBullet(ImDrawList* draw_list, ImVec2 pos, ImU32 col);
    IMGUI_API void          RenderCheckMark(ImDrawList* draw_list, ImVec2 pos, ImU32 col, float sz);
    IMGUI_API void          RenderMouseCursor(ImDrawList* draw_list, ImVec2 pos, float scale, ImGuiMouseCursor mouse_cursor, ImU32 col_fill, ImU32 col_border, ImU32 col_shadow);
    IMGUI_API void          RenderArrowPointingAt(ImDrawList* draw_list, ImVec2 pos, ImVec2 half_sz, ImGuiDir direction, ImU32 col);
    IMGUI_API void          RenderRectFilledRangeH(ImDrawList* draw_list, const ImRect& rect, ImU32 col, float x_start_norm, float x_end_norm, float rounding);
    IMGUI_API void          RenderRectFilledWithHole(ImDrawList* draw_list, ImRect outer, ImRect inner, ImU32 col, float rounding);

#ifndef IMGUI_DISABLE_OBSOLETE_FUNCTIONS
    // [1.71: 2019/06/07: Updating prototypes of some of the internal functions. Leaving those for reference for a short while]
    inline void RenderArrow(ImVec2 pos, ImGuiDir dir, float scale=1.0f) { ImGuiWindow* window = GetCurrentWindow(); RenderArrow(window->DrawList, pos, GetColorU32(ImGuiCol_Text), dir, scale); }
    inline void RenderBullet(ImVec2 pos)                                { ImGuiWindow* window = GetCurrentWindow(); RenderBullet(window->DrawList, pos, GetColorU32(ImGuiCol_Text)); }
#endif

    // Widgets
    IMGUI_API void          TextEx(const char* text, const char* text_end = NULL, ImGuiTextFlags flags = 0);
    IMGUI_API bool          ButtonEx(const char* label, const ImVec2& size_arg = ImVec2(0, 0), ImGuiButtonFlags flags = 0);
    IMGUI_API bool          CloseButton(ImGuiID id, const ImVec2& pos);
    IMGUI_API bool          CollapseButton(ImGuiID id, const ImVec2& pos);
    IMGUI_API bool          ArrowButtonEx(const char* str_id, ImGuiDir dir, ImVec2 size_arg, ImGuiButtonFlags flags = 0);
    IMGUI_API void          Scrollbar(ImGuiAxis axis);
    IMGUI_API bool          ScrollbarEx(const ImRect& bb, ImGuiID id, ImGuiAxis axis, float* p_scroll_v, float avail_v, float contents_v, ImDrawFlags flags);
    IMGUI_API bool          ImageButtonEx(ImGuiID id, ImTextureID texture_id, const ImVec2& size, const ImVec2& uv0, const ImVec2& uv1, const ImVec2& padding, const ImVec4& bg_col, const ImVec4& tint_col);
    IMGUI_API ImRect        GetWindowScrollbarRect(ImGuiWindow* window, ImGuiAxis axis);
    IMGUI_API ImGuiID       GetWindowScrollbarID(ImGuiWindow* window, ImGuiAxis axis);
    IMGUI_API ImGuiID       GetWindowResizeCornerID(ImGuiWindow* window, int n); // 0..3: corners
    IMGUI_API ImGuiID       GetWindowResizeBorderID(ImGuiWindow* window, ImGuiDir dir);
    IMGUI_API void          SeparatorEx(ImGuiSeparatorFlags flags);
    IMGUI_API bool          CheckboxFlags(const char* label, ImS64* flags, ImS64 flags_value);
    IMGUI_API bool          CheckboxFlags(const char* label, ImU64* flags, ImU64 flags_value);

    // Widgets low-level behaviors
    IMGUI_API bool          ButtonBehavior(const ImRect& bb, ImGuiID id, bool* out_hovered, bool* out_held, ImGuiButtonFlags flags = 0);
    IMGUI_API bool          DragBehavior(ImGuiID id, ImGuiDataType data_type, void* p_v, float v_speed, const void* p_min, const void* p_max, const char* format, ImGuiSliderFlags flags);
    IMGUI_API bool          SliderBehavior(const ImRect& bb, ImGuiID id, ImGuiDataType data_type, void* p_v, const void* p_min, const void* p_max, const char* format, ImGuiSliderFlags flags, ImRect* out_grab_bb);
    IMGUI_API bool          SplitterBehavior(const ImRect& bb, ImGuiID id, ImGuiAxis axis, float* size1, float* size2, float min_size1, float min_size2, float hover_extend = 0.0f, float hover_visibility_delay = 0.0f);
    IMGUI_API bool          TreeNodeBehavior(ImGuiID id, ImGuiTreeNodeFlags flags, const char* label, const char* label_end = NULL);
    IMGUI_API bool          TreeNodeBehaviorIsOpen(ImGuiID id, ImGuiTreeNodeFlags flags = 0);                     // Consume previous SetNextItemOpen() data, if any. May return true when logging
    IMGUI_API void          TreePushOverrideID(ImGuiID id);

    // Template functions are instantiated in imgui_widgets.cpp for a finite number of types.
    // To use them externally (for custom widget) you may need an "extern template" statement in your code in order to link to existing instances and silence Clang warnings (see #2036).
    // e.g. " extern template IMGUI_API float RoundScalarWithFormatT<float, float>(const char* format, ImGuiDataType data_type, float v); "
    template<typename T, typename SIGNED_T, typename FLOAT_T>   IMGUI_API float ScaleRatioFromValueT(ImGuiDataType data_type, T v, T v_min, T v_max, bool is_logarithmic, float logarithmic_zero_epsilon, float zero_deadzone_size);
    template<typename T, typename SIGNED_T, typename FLOAT_T>   IMGUI_API T     ScaleValueFromRatioT(ImGuiDataType data_type, float t, T v_min, T v_max, bool is_logarithmic, float logarithmic_zero_epsilon, float zero_deadzone_size);
    template<typename T, typename SIGNED_T, typename FLOAT_T>   IMGUI_API bool  DragBehaviorT(ImGuiDataType data_type, T* v, float v_speed, T v_min, T v_max, const char* format, ImGuiSliderFlags flags);
    template<typename T, typename SIGNED_T, typename FLOAT_T>   IMGUI_API bool  SliderBehaviorT(const ImRect& bb, ImGuiID id, ImGuiDataType data_type, T* v, T v_min, T v_max, const char* format, ImGuiSliderFlags flags, ImRect* out_grab_bb);
    template<typename T, typename SIGNED_T>                     IMGUI_API T     RoundScalarWithFormatT(const char* format, ImGuiDataType data_type, T v);
    template<typename T>                                        IMGUI_API bool  CheckboxFlagsT(const char* label, T* flags, T flags_value);

    // Data type helpers
    IMGUI_API const ImGuiDataTypeInfo*  DataTypeGetInfo(ImGuiDataType data_type);
    IMGUI_API int           DataTypeFormatString(char* buf, int buf_size, ImGuiDataType data_type, const void* p_data, const char* format);
    IMGUI_API void          DataTypeApplyOp(ImGuiDataType data_type, int op, void* output, const void* arg_1, const void* arg_2);
    IMGUI_API bool          DataTypeApplyOpFromText(const char* buf, const char* initial_value_buf, ImGuiDataType data_type, void* p_data, const char* format);
    IMGUI_API int           DataTypeCompare(ImGuiDataType data_type, const void* arg_1, const void* arg_2);
    IMGUI_API bool          DataTypeClamp(ImGuiDataType data_type, void* p_data, const void* p_min, const void* p_max);

    // InputText
    IMGUI_API bool          InputTextEx(const char* label, const char* hint, char* buf, int buf_size, const ImVec2& size_arg, ImGuiInputTextFlags flags, ImGuiInputTextCallback callback = NULL, void* user_data = NULL);
    IMGUI_API bool          TempInputText(const ImRect& bb, ImGuiID id, const char* label, char* buf, int buf_size, ImGuiInputTextFlags flags);
    IMGUI_API bool          TempInputScalar(const ImRect& bb, ImGuiID id, const char* label, ImGuiDataType data_type, void* p_data, const char* format, const void* p_clamp_min = NULL, const void* p_clamp_max = NULL);
    inline bool             TempInputIsActive(ImGuiID id)       { ImGuiContext& g = *GImGui; return (g.ActiveId == id && g.TempInputId == id); }
    inline ImGuiInputTextState* GetInputTextState(ImGuiID id)   { ImGuiContext& g = *GImGui; return (g.InputTextState.ID == id) ? &g.InputTextState : NULL; } // Get input text state if active

    // Color
    IMGUI_API void          ColorTooltip(const char* text, const float* col, ImGuiColorEditFlags flags);
    IMGUI_API void          ColorEditOptionsPopup(const float* col, ImGuiColorEditFlags flags);
    IMGUI_API void          ColorPickerOptionsPopup(const float* ref_col, ImGuiColorEditFlags flags);

    // Plot
    IMGUI_API int           PlotEx(ImGuiPlotType plot_type, const char* label, float (*values_getter)(void* data, int idx), void* data, int values_count, int values_offset, const char* overlay_text, float scale_min, float scale_max, ImVec2 frame_size);

    // Shade functions (write over already created vertices)
    IMGUI_API void          ShadeVertsLinearColorGradientKeepAlpha(ImDrawList* draw_list, int vert_start_idx, int vert_end_idx, ImVec2 gradient_p0, ImVec2 gradient_p1, ImU32 col0, ImU32 col1);
    IMGUI_API void          ShadeVertsLinearUV(ImDrawList* draw_list, int vert_start_idx, int vert_end_idx, const ImVec2& a, const ImVec2& b, const ImVec2& uv_a, const ImVec2& uv_b, bool clamp);

    // Garbage collection
    IMGUI_API void          GcCompactTransientMiscBuffers();
    IMGUI_API void          GcCompactTransientWindowBuffers(ImGuiWindow* window);
    IMGUI_API void          GcAwakeTransientWindowBuffers(ImGuiWindow* window);

    // Debug Tools
    IMGUI_API void          ErrorCheckEndFrameRecover(ImGuiErrorLogCallback log_callback, void* user_data = NULL);
    inline void             DebugDrawItemRect(ImU32 col = IM_COL32(255,0,0,255))    { ImGuiContext& g = *GImGui; ImGuiWindow* window = g.CurrentWindow; GetForegroundDrawList(window)->AddRect(g.LastItemData.Rect.Min, g.LastItemData.Rect.Max, col); }
    inline void             DebugStartItemPicker()                                  { ImGuiContext& g = *GImGui; g.DebugItemPickerActive = true; }

    IMGUI_API void          ShowFontAtlas(ImFontAtlas* atlas);
    IMGUI_API void          DebugNodeColumns(ImGuiOldColumns* columns);
    IMGUI_API void          DebugNodeDrawList(ImGuiWindow* window, const ImDrawList* draw_list, const char* label);
    IMGUI_API void          DebugNodeDrawCmdShowMeshAndBoundingBox(ImDrawList* out_draw_list, const ImDrawList* draw_list, const ImDrawCmd* draw_cmd, bool show_mesh, bool show_aabb);
    IMGUI_API void          DebugNodeFont(ImFont* font);
    IMGUI_API void          DebugNodeStorage(ImGuiStorage* storage, const char* label);
    IMGUI_API void          DebugNodeTabBar(ImGuiTabBar* tab_bar, const char* label);
    IMGUI_API void          DebugNodeTable(ImGuiTable* table);
    IMGUI_API void          DebugNodeTableSettings(ImGuiTableSettings* settings);
    IMGUI_API void          DebugNodeWindow(ImGuiWindow* window, const char* label);
    IMGUI_API void          DebugNodeWindowSettings(ImGuiWindowSettings* settings);
    IMGUI_API void          DebugNodeWindowsList(ImVector<ImGuiWindow*>* windows, const char* label);
    IMGUI_API void          DebugNodeViewport(ImGuiViewportP* viewport);
    IMGUI_API void          DebugRenderViewportThumbnail(ImDrawList* draw_list, ImGuiViewportP* viewport, const ImRect& bb);

} // namespace ImGui


//-----------------------------------------------------------------------------
// [SECTION] ImFontAtlas internal API
//-----------------------------------------------------------------------------

// This structure is likely to evolve as we add support for incremental atlas updates
struct ImFontBuilderIO
{
    bool    (*FontBuilder_Build)(ImFontAtlas* atlas);
};

// Helper for font builder
IMGUI_API const ImFontBuilderIO* ImFontAtlasGetBuilderForStbTruetype();
IMGUI_API void      ImFontAtlasBuildInit(ImFontAtlas* atlas);
IMGUI_API void      ImFontAtlasBuildSetupFont(ImFontAtlas* atlas, ImFont* font, ImFontConfig* font_config, float ascent, float descent);
IMGUI_API void      ImFontAtlasBuildPackCustomRects(ImFontAtlas* atlas, void* stbrp_context_opaque);
IMGUI_API void      ImFontAtlasBuildFinish(ImFontAtlas* atlas);
IMGUI_API void      ImFontAtlasBuildRender8bppRectFromString(ImFontAtlas* atlas, int x, int y, int w, int h, const char* in_str, char in_marker_char, unsigned char in_marker_pixel_value);
IMGUI_API void      ImFontAtlasBuildRender32bppRectFromString(ImFontAtlas* atlas, int x, int y, int w, int h, const char* in_str, char in_marker_char, unsigned int in_marker_pixel_value);
IMGUI_API void      ImFontAtlasBuildMultiplyCalcLookupTable(unsigned char out_table[256], float in_multiply_factor);
IMGUI_API void      ImFontAtlasBuildMultiplyRectAlpha8(const unsigned char table[256], unsigned char* pixels, int x, int y, int w, int h, int stride);

//-----------------------------------------------------------------------------
// [SECTION] Test Engine specific hooks (imgui_test_engine)
//-----------------------------------------------------------------------------

#ifdef IMGUI_ENABLE_TEST_ENGINE
extern void         ImGuiTestEngineHook_ItemAdd(ImGuiContext* ctx, const ImRect& bb, ImGuiID id);
extern void         ImGuiTestEngineHook_ItemInfo(ImGuiContext* ctx, ImGuiID id, const char* label, ImGuiItemStatusFlags flags);
extern void         ImGuiTestEngineHook_IdInfo(ImGuiContext* ctx, ImGuiDataType data_type, ImGuiID id, const void* data_id);
extern void         ImGuiTestEngineHook_IdInfo(ImGuiContext* ctx, ImGuiDataType data_type, ImGuiID id, const void* data_id, const void* data_id_end);
extern void         ImGuiTestEngineHook_Log(ImGuiContext* ctx, const char* fmt, ...);
#define IMGUI_TEST_ENGINE_ITEM_ADD(_BB,_ID)                 if (g.TestEngineHookItems) ImGuiTestEngineHook_ItemAdd(&g, _BB, _ID)               // Register item bounding box
#define IMGUI_TEST_ENGINE_ITEM_INFO(_ID,_LABEL,_FLAGS)      if (g.TestEngineHookItems) ImGuiTestEngineHook_ItemInfo(&g, _ID, _LABEL, _FLAGS)   // Register item label and status flags (optional)
#define IMGUI_TEST_ENGINE_LOG(_FMT,...)                     if (g.TestEngineHookItems) ImGuiTestEngineHook_Log(&g, _FMT, __VA_ARGS__)          // Custom log entry from user land into test log
#define IMGUI_TEST_ENGINE_ID_INFO(_ID,_TYPE,_DATA)          if (g.TestEngineHookIdInfo == _ID) ImGuiTestEngineHook_IdInfo(&g, _TYPE, _ID, (const void*)(_DATA));
#define IMGUI_TEST_ENGINE_ID_INFO2(_ID,_TYPE,_DATA,_DATA2)  if (g.TestEngineHookIdInfo == _ID) ImGuiTestEngineHook_IdInfo(&g, _TYPE, _ID, (const void*)(_DATA), (const void*)(_DATA2));
#else
#define IMGUI_TEST_ENGINE_ITEM_INFO(_ID,_LABEL,_FLAGS)      ((void)0)
#endif

//-----------------------------------------------------------------------------

#if defined(__clang__)
#pragma clang diagnostic pop
#elif defined(__GNUC__)
#pragma GCC diagnostic pop
#endif

#ifdef _MSC_VER
#pragma warning (pop)
#endif

#endif // #ifndef IMGUI_DISABLE
