package imgui

func ErrorCheckNewFrameSanityChecks() {
	g := guiContext

	// Check user IM_ASSERT macro
	// (IF YOU GET A WARNING OR COMPILE ERROR HERE: it means your assert macro is incorrectly defined!
	//  If your macro uses multiple statements, it NEEDS to be surrounded by a 'do { ... } while (0)' block.
	//  This is a common C/C++ idiom to allow multiple statements macros to be used in control flow blocks.)
	// #define IM_ASSERT(EXPR)   if (SomeCode(EXPR)) SomeMoreCode();                    // Wrong!
	// #define IM_ASSERT(EXPR)   do { if (SomeCode(EXPR)) SomeMoreCode(); } while (0)   // Correct!
	if true {
		IM_ASSERT(true)
	} else {
		IM_ASSERT(false)
	}

	// Check user data
	// (We pass an error message in the assert expression to make it visible to programmers who are not using a debugger, as most assert handlers display their argument)
	IM_ASSERT(g.Initialized)
	IM_ASSERT_USER_ERROR((g.IO.DeltaTime > 0.0 || g.FrameCount == 0), "Need a positive DeltaTime!")
	IM_ASSERT_USER_ERROR((g.FrameCount == 0 || g.FrameCountEnded == g.FrameCount), "Forgot to call Render() or EndFrame() at the end of the previous frame?")
	IM_ASSERT_USER_ERROR(g.IO.DisplaySize.x >= 0.0 && g.IO.DisplaySize.y >= 0.0, "Invalid DisplaySize value!")
	IM_ASSERT_USER_ERROR(g.IO.Fonts.IsBuilt(), "Font Atlas not built! Make sure you called ImGui_ImplXXXX_NewFrame() function for renderer backend, which should call io.Fonts.GetTexDataAsRGBA32() / GetTexDataAsAlpha8()")
	IM_ASSERT_USER_ERROR(g.Style.CurveTessellationTol > 0, "Invalid style setting!")
	IM_ASSERT_USER_ERROR(g.Style.CircleTessellationMaxError > 0.0, "Invalid style setting!")
	IM_ASSERT_USER_ERROR(g.Style.Alpha >= 0.0 && g.Style.Alpha <= 1.0, "Invalid style setting!") // Allows us to avoid a few clamps in color computations
	IM_ASSERT_USER_ERROR(g.Style.WindowMinSize.x >= 1.0 && g.Style.WindowMinSize.y >= 1.0, "Invalid style setting.")
	IM_ASSERT(g.Style.WindowMenuButtonPosition == ImGuiDir_None || g.Style.WindowMenuButtonPosition == ImGuiDir_Left || g.Style.WindowMenuButtonPosition == ImGuiDir_Right)
	for n := ImGuiKey(0); n < ImGuiKey_COUNT; n++ {
		IM_ASSERT_USER_ERROR(g.IO.KeyMap[n] >= -1 && g.IO.KeyMap[n] < int(len(g.IO.KeysDown)), "io.KeyMap[] contains an out of bound value (need to be 0..512, or -1 for unmapped key)")
	}

	// Check: required key mapping (we intentionally do NOT check all keys to not pressure user into setting up everything, but Space is required and was only added in 1.60 WIP)
	if g.IO.ConfigFlags&ImGuiConfigFlags_NavEnableKeyboard != 0 {
		IM_ASSERT_USER_ERROR(g.IO.KeyMap[ImGuiKey_Space] != -1, "ImGuiKey_Space is not mapped, required for keyboard navigation.")
	}

	// Check: the io.ConfigWindowsResizeFromEdges option requires backend to honor mouse cursor changes and set the ImGuiBackendFlags_HasMouseCursors flag accordingly.
	if g.IO.ConfigWindowsResizeFromEdges && g.IO.BackendFlags&ImGuiBackendFlags_HasMouseCursors == 0 {
		g.IO.ConfigWindowsResizeFromEdges = false
	}
}

// start a new Dear ImGui frame, you can submit any command from this pountil int Render()/EndFrame().
func NewFrame() {
	IM_ASSERT_USER_ERROR(guiContext != nil, "No current context. Did you call ImGui::CreateContext() and ImGui::SetCurrentContext() ?")
	g := guiContext

	// Remove pending delete hooks before frame start.
	// This deferred removal avoid issues of removal while iterating the hook vector
	for n := len(g.Hooks) - 1; n >= 0; n-- {
		if g.Hooks[n].Type == ImGuiContextHookType_PendingRemoval_ {
			//erase
			g.Hooks[n] = g.Hooks[len(g.Hooks)-1]
			g.Hooks = g.Hooks[:len(g.Hooks)-1]
		}
	}

	CallContextHooks(g, ImGuiContextHookType_NewFramePre)

	// Check and assert for various common IO and Configuration mistakes
	ErrorCheckNewFrameSanityChecks()

	// Load settings on first frame, save settings when modified (after a delay)
	UpdateSettings()

	g.Time += double(g.IO.DeltaTime)
	g.WithinFrameScope = true
	g.FrameCount += 1
	g.TooltipOverrideCount = 0
	g.WindowsActiveCount = 0
	g.MenusIdSubmittedThisFrame = g.MenusIdSubmittedThisFrame[:0]

	// Calculate frame-rate for the user, as a purely luxurious feature
	g.FramerateSecPerFrameAccum += g.IO.DeltaTime - g.FramerateSecPerFrame[g.FramerateSecPerFrameIdx]
	g.FramerateSecPerFrame[g.FramerateSecPerFrameIdx] = g.IO.DeltaTime
	g.FramerateSecPerFrameIdx = (g.FramerateSecPerFrameIdx + 1) % int(len(g.FramerateSecPerFrame))
	g.FramerateSecPerFrameCount = min(g.FramerateSecPerFrameCount+1, int(len(g.FramerateSecPerFrame)))
	if g.FramerateSecPerFrameAccum > 0.0 {
		g.IO.Framerate = (1.0 / (g.FramerateSecPerFrameAccum / (float)(g.FramerateSecPerFrameCount)))
	} else {
		g.IO.Framerate = FLT_MAX
	}

	UpdateViewportsNewFrame()

	// Setup current font and draw list shared data
	g.IO.Fonts.Locked = true
	SetCurrentFont(GetDefaultFont())
	IM_ASSERT(g.Font.IsLoaded())
	var virtual_space = ImRect{ImVec2{FLT_MAX, FLT_MAX}, ImVec2{-FLT_MAX, -FLT_MAX}}
	for n := range g.Viewports {
		virtual_space.AddRect(g.Viewports[n].GetMainRect())
	}
	g.DrawListSharedData.ClipRectFullscreen = virtual_space.ToVec4()
	g.DrawListSharedData.CurveTessellationTol = g.Style.CurveTessellationTol
	g.DrawListSharedData.SetCircleTessellationMaxError(g.Style.CircleTessellationMaxError)
	g.DrawListSharedData.InitialFlags = ImDrawListFlags_None
	if g.Style.AntiAliasedLines {
		g.DrawListSharedData.InitialFlags |= ImDrawListFlags_AntiAliasedLines
	}
	if g.Style.AntiAliasedLinesUseTex && g.Font.ContainerAtlas.Flags&ImFontAtlasFlags_NoBakedLines == 0 {
		g.DrawListSharedData.InitialFlags |= ImDrawListFlags_AntiAliasedLinesUseTex
	}
	if g.Style.AntiAliasedFill {
		g.DrawListSharedData.InitialFlags |= ImDrawListFlags_AntiAliasedFill
	}
	if g.IO.BackendFlags&ImGuiBackendFlags_RendererHasVtxOffset != 0 {
		g.DrawListSharedData.InitialFlags |= ImDrawListFlags_AllowVtxOffset
	}

	// Mark rendering data as invalid to prevent user who may have a handle on it to use it.
	for n := range g.Viewports {
		var viewport = g.Viewports[n]
		viewport.DrawDataP.Clear()
	}

	// Drag and drop keep the source ID alive so even if the source disappear our state is consistent
	if g.DragDropActive && g.DragDropPayload.SourceId == g.ActiveId {
		KeepAliveID(g.DragDropPayload.SourceId)
	}

	// Update HoveredId data
	if g.HoveredIdPreviousFrame == 0 {
		g.HoveredIdTimer = 0.0
	}
	if g.HoveredIdPreviousFrame == 0 || (g.HoveredId != 0 && g.ActiveId == g.HoveredId) {
		g.HoveredIdNotActiveTimer = 0.0
	}
	if g.HoveredId != 0 {
		g.HoveredIdTimer += g.IO.DeltaTime
	}
	if g.HoveredId != 0 && g.ActiveId != g.HoveredId {
		g.HoveredIdNotActiveTimer += g.IO.DeltaTime
	}
	g.HoveredIdPreviousFrame = g.HoveredId
	g.HoveredIdPreviousFrameUsingMouseWheel = g.HoveredIdUsingMouseWheel
	g.HoveredId = 0
	g.HoveredIdAllowOverlap = false
	g.HoveredIdUsingMouseWheel = false
	g.HoveredIdDisabled = false

	// Update ActiveId data (clear reference to active widget if the widget isn't alive anymore)
	if g.ActiveIdIsAlive != g.ActiveId && g.ActiveIdPreviousFrame == g.ActiveId && g.ActiveId != 0 {
		ClearActiveID()
	}
	if g.ActiveId != 0 {
		g.ActiveIdTimer += g.IO.DeltaTime
	}
	g.LastActiveIdTimer += g.IO.DeltaTime
	g.ActiveIdPreviousFrame = g.ActiveId
	g.ActiveIdPreviousFrameWindow = g.ActiveIdWindow
	g.ActiveIdPreviousFrameHasBeenEditedBefore = g.ActiveIdHasBeenEditedBefore
	g.ActiveIdIsAlive = 0
	g.ActiveIdHasBeenEditedThisFrame = false
	g.ActiveIdPreviousFrameIsAlive = false
	g.ActiveIdIsJustActivated = false
	if g.TempInputId != 0 && g.ActiveId != g.TempInputId {
		g.TempInputId = 0
	}
	if g.ActiveId == 0 {
		g.ActiveIdUsingNavDirMask = 0x00
		g.ActiveIdUsingNavInputMask = 0x00
		g.ActiveIdUsingKeyInputMask = 0x00
	}

	// Drag and drop
	g.DragDropAcceptIdPrev = g.DragDropAcceptIdCurr
	g.DragDropAcceptIdCurr = 0
	g.DragDropAcceptIdCurrRectSurface = FLT_MAX
	g.DragDropWithinSource = false
	g.DragDropWithinTarget = false
	g.DragDropHoldJustPressedId = 0

	// Update keyboard input state
	// Synchronize io.KeyMods with individual modifiers io.KeyXXX bools
	g.IO.KeyMods = GetMergedKeyModFlags()
	copy(g.IO.KeysDownDurationPrev[:], g.IO.KeysDownDuration[:])
	for i := 0; i < len(g.IO.KeysDown); i++ {
		if g.IO.KeysDown[i] {
			if g.IO.KeysDownDuration[i] < 0.0 {
				g.IO.KeysDownDuration[i] = 0.0
			} else {
				g.IO.KeysDownDuration[i] += g.IO.DeltaTime
			}
		} else {
			g.IO.KeysDownDuration[i] = -1.0
		}
	}

	// Update gamepad/keyboard navigation
	NavUpdate()

	// Update mouse input state
	UpdateMouseInputs()

	// Find hovered window
	// (needs to be before UpdateMouseMovingWindowNewFrame so we fill guiContext.HoveredWindowUnderMovingWindow on the mouse release frame)
	UpdateHoveredWindowAndCaptureFlags()

	// Handle user moving window with mouse (at the beginning of the frame to avoid input lag or sheering)
	UpdateMouseMovingWindowNewFrame()

	// Background darkening/whitening
	if GetTopMostPopupModal() != nil || (g.NavWindowingTarget != nil && g.NavWindowingHighlightAlpha > 0.0) {
		g.DimBgRatio = min(g.DimBgRatio+g.IO.DeltaTime*6.0, 1.0)
	} else {
		g.DimBgRatio = max(g.DimBgRatio-g.IO.DeltaTime*10.0, 0.0)
	}

	g.MouseCursor = ImGuiMouseCursor_Arrow
	g.WantCaptureMouseNextFrame = -1
	g.WantCaptureKeyboardNextFrame = -1
	g.WantTextInputNextFrame = -1
	g.PlatformImePos = ImVec2{1.0, 1.0} // OS Input Method Editor showing on top-left of our window by default

	// Mouse wheel scrolling, scale
	UpdateMouseWheel()

	// Update legacy TAB focus
	UpdateTabFocus()

	// Mark all windows as not visible and compact unused memory.
	IM_ASSERT(len(g.WindowsFocusOrder) <= len(g.Windows))
	var memory_compact_start_time float
	if g.GcCompactAll || g.IO.ConfigMemoryCompactTimer < 0.0 {
		memory_compact_start_time = FLT_MAX
	} else {
		memory_compact_start_time = (float)(g.Time - double(g.IO.ConfigMemoryCompactTimer))
	}
	for i := range g.Windows {
		var window = g.Windows[i]
		window.WasActive = window.Active
		window.BeginCount = 0
		window.Active = false
		window.WriteAccessed = false

		// Garbage collect transient buffers of recently unused windows
		if !window.WasActive && !window.MemoryCompacted && window.LastTimeActive < memory_compact_start_time {
			GcCompactTransientWindowBuffers(window)
		}
	}

	// Garbage collect transient buffers of recently unused tables
	/*for i := range guiContext.TablesLastTimeActive {
		if guiContext.TablesLastTimeActive[i] >= 0.0 && guiContext.TablesLastTimeActive[i] < memory_compact_start_time {
			TableGcCompactTransientBuffers(guiContext.Tables[i])
		}
	}
	for i := range guiContext.TablesTempDataStack {
		if guiContext.TablesTempDataStack[i].LastTimeActive >= 0.0 && guiContext.TablesTempDataStack[i].LastTimeActive < memory_compact_start_time {
			TableGcCompactTransientBuffers(&guiContext.TablesTempDataStack[i])
		}
	}*/
	if g.GcCompactAll {
		GcCompactTransientMiscBuffers()
	}
	g.GcCompactAll = false

	// Closing the focused window restore focus to the first active root window in descending z-order
	if g.NavWindow != nil && !g.NavWindow.WasActive {
		FocusTopMostWindowUnderOne(nil, nil)
	}

	// No window should be open at the beginning of the frame.
	// But in order to allow the user to call NewFrame() multiple times without calling Render(), we are doing an explicit clear.
	g.CurrentWindowStack = g.CurrentWindowStack[:0]
	g.BeginPopupStack = g.BeginPopupStack[:0]
	g.ItemFlagsStack = g.ItemFlagsStack[:0]
	g.ItemFlagsStack = append(g.ItemFlagsStack, ImGuiItemFlags_None)
	g.GroupStack = g.GroupStack[:0]

	// [DEBUG] Item picker tool - start with DebugStartItemPicker() - useful to visually select an item and break into its call-stack.
	UpdateDebugToolItemPicker()

	// Create implicit/fallback window - which we will only render it if the user has added something to it.
	// We don't use "Debug" to avoid colliding with user trying to create a "Debug" window with custom flags.
	// This fallback is particularly important as it avoid ImGui:: calls from crashing.
	g.WithinFrameScopeWithImplicitWindow = true
	SetNextWindowSize(&ImVec2{400, 400}, ImGuiCond_FirstUseEver)
	Begin("Debug##Default", nil, 0)
	IM_ASSERT(g.CurrentWindow.IsFallbackWindow)

	CallContextHooks(g, ImGuiContextHookType_NewFramePost)
}

func ErrorCheckEndFrameSanityChecks() {
	g := guiContext

	// Verify that io.KeyXXX fields haven't been tampered with. Key mods should not be modified between NewFrame() and EndFrame()
	// One possible reason leading to this assert is that your backends update inputs _AFTER_ NewFrame().
	// It is known that when some modal native windows called mid-frame takes focus away, some backends such as GLFW will
	// send key release events mid-frame. This would normally trigger this assertion and lead to sheared inputs.
	// We silently accommodate for this case by ignoring/ the case where all io.KeyXXX modifiers were released (aka key_mod_flags == 0),
	// while still correctly asserting on mid-frame key press events.
	var key_mod_flags = GetMergedKeyModFlags()
	IM_ASSERT_USER_ERROR((key_mod_flags == 0 || g.IO.KeyMods == key_mod_flags), "Mismatching io.KeyCtrl/io.KeyShift/io.KeyAlt/io.KeySuper vs io.KeyMods")

	// Recover from errors
	ErrorCheckEndFrameRecover(nil, nil)

	// Report when there is a mismatch of Begin/BeginChild vs End/EndChild calls. Important: Remember that the Begin/BeginChild API requires you
	// to always call End/EndChild even if Begin/BeginChild returns false! (this is unfortunately inconsistent with most other Begin* API).
	if len(g.CurrentWindowStack) != 1 {
		if len(g.CurrentWindowStack) > 1 {
			IM_ASSERT_USER_ERROR(len(g.CurrentWindowStack) == 1, "Mismatched Begin/BeginChild vs End/EndChild calls: did you forget to call End/EndChild?")
			for len(g.CurrentWindowStack) > 1 {
				End()
			}
		} else {
			IM_ASSERT_USER_ERROR(len(g.CurrentWindowStack) == 1, "Mismatched Begin/BeginChild vs End/EndChild calls: did you call End/EndChild too much?")
		}
	}

	IM_ASSERT_USER_ERROR(len(g.GroupStack) == 0, "Missing EndGroup call!")
}

// ends the Dear ImGui frame. automatically called by Render(). If you don't need to render data (skipping rendering) you may call EndFrame() without Render()... but you'll have wasted CPU already! If you don't need to render, better to not create any windows and not call NewFrame() at all!
func EndFrame() {
	IM_ASSERT(guiContext.Initialized)

	// Don't process EndFrame() multiple times.
	if guiContext.FrameCountEnded == guiContext.FrameCount {
		return
	}
	IM_ASSERT_USER_ERROR(guiContext.WithinFrameScope, "Forgot to call ImGui::NewFrame()?")

	CallContextHooks(guiContext, ImGuiContextHookType_EndFramePre)

	ErrorCheckEndFrameSanityChecks()

	// Notify OS when our Input Method Editor cursor has moved (e.guiContext. CJK inputs using Microsoft IME)
	if guiContext.IO.ImeSetInputScreenPosFn != nil && (guiContext.PlatformImeLastPos.x == FLT_MAX || ImLengthSqrVec2(guiContext.PlatformImeLastPos.Sub(guiContext.PlatformImePos)) > 0.0001) {
		guiContext.IO.ImeSetInputScreenPosFn((int)(guiContext.PlatformImePos.x), (int)(guiContext.PlatformImePos.y))
		guiContext.PlatformImeLastPos = guiContext.PlatformImePos
	}

	// Hide implicit/fallback "Debug" window if it hasn't been used
	guiContext.WithinFrameScopeWithImplicitWindow = false
	if guiContext.CurrentWindow != nil && !guiContext.CurrentWindow.WriteAccessed {
		guiContext.CurrentWindow.Active = false
	}
	End()

	// Update navigation: CTRL+Tab, wrap-around requests
	NavEndFrame()

	// Drag and Drop: Elapse payload (if delivered, or if source stops being submitted)
	if guiContext.DragDropActive {
		var is_delivered = guiContext.DragDropPayload.Delivery
		var is_elapsed = (guiContext.DragDropPayload.DataFrameCount+1 < guiContext.FrameCount) && ((guiContext.DragDropSourceFlags&ImGuiDragDropFlags_SourceAutoExpirePayload != 0) || !IsMouseDown(guiContext.DragDropMouseButton))
		if is_delivered || is_elapsed {
			ClearDragDrop()
		}
	}

	// Drag and Drop: Fallback for source tooltip. This is not ideal but better than nothing.
	if guiContext.DragDropActive && guiContext.DragDropSourceFrameCount < guiContext.FrameCount && guiContext.DragDropSourceFlags&ImGuiDragDropFlags_SourceNoPreviewTooltip == 0 {
		guiContext.DragDropWithinSource = true
		SetTooltip("...")
		guiContext.DragDropWithinSource = false
	}

	// End frame
	guiContext.WithinFrameScope = false
	guiContext.FrameCountEnded = guiContext.FrameCount

	// Initiate moving window + handle left-click and right-click focus
	UpdateMouseMovingWindowEndFrame()

	// Sort the window list so that all child windows are after their parent
	// We cannot do that on FocusWindow() because children may not exist yet
	guiContext.WindowsTempSortBuffer = guiContext.WindowsTempSortBuffer[:0]
	guiContext.WindowsTempSortBuffer = make([]*ImGuiWindow, 0, len(guiContext.Windows))
	for i := range guiContext.Windows {
		var window = guiContext.Windows[i]
		if window.Active && (window.Flags&ImGuiWindowFlags_ChildWindow != 0) { // if a child is active its parent will add it
			continue
		}
		AddWindowToSortBuffer(&guiContext.WindowsTempSortBuffer, window)
	}

	// This usually assert if there is a mismatch between the ImGuiWindowFlags_ChildWindow / ParentWindow values and DC.ChildWindows[] in parents, aka we've done something wrong.
	IM_ASSERT(len(guiContext.Windows) == len(guiContext.WindowsTempSortBuffer))
	guiContext.Windows, guiContext.WindowsTempSortBuffer = guiContext.WindowsTempSortBuffer, guiContext.Windows
	guiContext.IO.MetricsActiveWindows = guiContext.WindowsActiveCount

	// Unlock font atlas
	guiContext.IO.Fonts.Locked = false

	// Clear Input data for next frame
	guiContext.IO.MouseWheel = 0
	guiContext.IO.MouseWheelH = 0.0
	guiContext.IO.InputQueueCharacters = guiContext.IO.InputQueueCharacters[:0]
	guiContext.IO.KeyModsPrev = guiContext.IO.KeyMods // doing it here is better than in NewFrame() as we'll tolerate backend writing to KeyMods. If we want to firmly disallow it we should detect it.
	guiContext.IO.NavInputs = [20]float{}

	CallContextHooks(guiContext, ImGuiContextHookType_EndFramePost)
}
