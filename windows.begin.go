package imgui

// Push a new Dear ImGui window to add widgets to.
// - A default window called "Debug" is automatically stacked at the beginning of every frame so you can use widgets without explicitly calling a Begin/End pair.
// - Begin/End can be called multiple times during the frame with the same window name to append content.
// - The window name is used as a unique identifier to preserve window information across frames (and save rudimentary information to the .ini file).
//   You can use the "##" or "###" markers to use the same label with different id, or same id with different label. See documentation at the top of this file.
// - Return false when window is collapsed, so you can early out in your code. You always need to call ImGui::End() even if false is returned.
// - Passing 'bool* p_open' displays a Close button on the upper-right corner of the window, the pointed value will be set to false when the button is pressed.
func Begin(name string, p_open *bool, flags ImGuiWindowFlags) bool {
	var g = GImGui
	var style = g.Style
	IM_ASSERT(name != "")                        // Window name required
	IM_ASSERT(g.WithinFrameScope)                // Forgot to call ImGui::NewFrame()
	IM_ASSERT(g.FrameCountEnded != g.FrameCount) // Called ImGui::Render() or ImGui::EndFrame() and haven't called ImGui::NewFrame() again yet

	// Find or create
	var window = FindWindowByName(name)
	var window_just_created = (window == nil)
	if window_just_created {
		window = CreateNewWindow(name, flags)
	}

	// Automatically disable manual moving/resizing when NoInputs is set
	if (flags & ImGuiWindowFlags_NoInputs) == ImGuiWindowFlags_NoInputs {
		flags |= ImGuiWindowFlags_NoMove | ImGuiWindowFlags_NoResize
	}

	if flags&ImGuiWindowFlags_NavFlattened != 0 {
		IM_ASSERT(flags&ImGuiWindowFlags_ChildWindow != 0)
	}

	var current_frame int = g.FrameCount
	var first_begin_of_the_frame bool = (window.LastFrameActive != current_frame)
	window.IsFallbackWindow = (len(g.CurrentWindowStack) == 0 && g.WithinFrameScopeWithImplicitWindow)

	// Update the Appearing flag
	var window_just_activated_by_user bool = (window.LastFrameActive < current_frame-1) // Not using !WasActive because the implicit "Debug" window would always toggle off.on
	if flags&ImGuiWindowFlags_Popup != 0 {
		var popup_ref *ImGuiPopupData = &g.OpenPopupStack[len(g.BeginPopupStack)]
		window_just_activated_by_user = window_just_activated_by_user || (window.PopupId != popup_ref.PopupId) // We recycle popups so treat window as activated if popup id changed
		window_just_activated_by_user = window_just_activated_by_user || (window != popup_ref.Window)
	}
	window.Appearing = window_just_activated_by_user
	if window.Appearing {
		SetWindowConditionAllowFlags(window, ImGuiCond_Appearing, true)
	}

	// Update Flags, LastFrameActive, BeginOrderXXX fields
	if first_begin_of_the_frame {
		window.Flags = (ImGuiWindowFlags)(flags)
		window.LastFrameActive = current_frame
		window.LastTimeActive = (float)(g.Time)
		window.BeginOrderWithinParent = 0
		window.BeginOrderWithinContext = (short)(g.WindowsActiveCount)
		g.WindowsActiveCount++
	} else {
		flags = window.Flags
	}

	// Parent window is latched only on the first call to Begin() of the frame, so further append-calls can be done from a different window stack
	var parent_window_in_stack *ImGuiWindow
	if len(g.CurrentWindowStack) > 0 {
		parent_window_in_stack = g.CurrentWindowStack[len(g.CurrentWindowStack)-1].Window
	}
	var parent_window *ImGuiWindow
	if first_begin_of_the_frame {
		if flags&(ImGuiWindowFlags_ChildWindow|ImGuiWindowFlags_Popup) != 0 {
			parent_window = parent_window_in_stack
		} else {
			parent_window = nil
		}
	} else {
		parent_window = window.ParentWindow
	}
	IM_ASSERT(parent_window != nil || (flags&ImGuiWindowFlags_ChildWindow) == 0)

	// We allow window memory to be compacted so recreate the base stack when needed.
	if len(window.IDStack) == 0 {
		window.IDStack = append(window.IDStack, window.ID)
	}

	// Add to stack
	// We intentionally set g.CurrentWindow to nil to prevent usage until when the viewport is set, then will call SetCurrentWindow()
	var window_stack_data ImGuiWindowStackData
	window_stack_data.Window = window
	window_stack_data.ParentLastItemDataBackup = g.LastItemData
	g.CurrentWindowStack = append(g.CurrentWindowStack, window_stack_data)
	g.CurrentWindow = window
	window.DC.StackSizesOnBegin.SetToCurrentState()
	g.CurrentWindow = nil

	if flags&ImGuiWindowFlags_Popup != 0 {
		var popup_ref *ImGuiPopupData = &g.OpenPopupStack[len(g.BeginPopupStack)]
		popup_ref.Window = window
		g.BeginPopupStack = append(g.BeginPopupStack, *popup_ref)
		window.PopupId = popup_ref.PopupId
	}

	// Update .RootWindow and others pointers (before any possible call to FocusWindow)
	if first_begin_of_the_frame {
		UpdateWindowParentAndRootLinks(window, flags, parent_window)
	}

	// Process SetNextWindow***() calls
	// (FIXME: Consider splitting the HasXXX flags into X/Y components
	var window_pos_set_by_api bool = false
	var window_size_x_set_by_api bool = false
	var window_size_y_set_by_api bool = false
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasPos != 0 {
		window_pos_set_by_api = (window.SetWindowPosAllowFlags & g.NextWindowData.PosCond) != 0
		if window_pos_set_by_api && ImLengthSqrVec2(g.NextWindowData.PosPivotVal) > 0.00001 {
			// May be processed on the next frame if this is our first frame and we are measuring size
			// FIXME: Look into removing the branch so everything can go through this same code path for consistency.
			window.SetWindowPosVal = g.NextWindowData.PosVal
			window.SetWindowPosPivot = g.NextWindowData.PosPivotVal
			window.SetWindowPosAllowFlags &= ^(ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing)
		} else {
			setWindowPos(window, &g.NextWindowData.PosVal, g.NextWindowData.PosCond)
		}
	}
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasSize != 0 {
		window_size_x_set_by_api = (window.SetWindowSizeAllowFlags&g.NextWindowData.SizeCond) != 0 && (g.NextWindowData.SizeVal.x > 0.0)
		window_size_y_set_by_api = (window.SetWindowSizeAllowFlags&g.NextWindowData.SizeCond) != 0 && (g.NextWindowData.SizeVal.y > 0.0)
		setWindowSize(window, &g.NextWindowData.SizeVal, g.NextWindowData.SizeCond)
	}
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasScroll != 0 {
		if g.NextWindowData.ScrollVal.x >= 0.0 {
			window.ScrollTarget.x = g.NextWindowData.ScrollVal.x
			window.ScrollTargetCenterRatio.x = 0.0
		}
		if g.NextWindowData.ScrollVal.y >= 0.0 {
			window.ScrollTarget.y = g.NextWindowData.ScrollVal.y
			window.ScrollTargetCenterRatio.y = 0.0
		}
	}
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasContentSize != 0 {
		window.ContentSizeExplicit = g.NextWindowData.ContentSizeVal
	} else if first_begin_of_the_frame {
		window.ContentSizeExplicit = ImVec2{}
	}
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasCollapsed != 0 {
		setWindowCollapsed(window, g.NextWindowData.CollapsedVal, g.NextWindowData.CollapsedCond)
	}
	if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasFocus != 0 {
		FocusWindow(window)
	}
	if window.Appearing {
		SetWindowConditionAllowFlags(window, ImGuiCond_Appearing, false)
	}

	// When reusing window again multiple times a frame, just append content (don't need to setup again)
	if first_begin_of_the_frame {
		// Initialize
		var window_is_child_tooltip bool = (flags&ImGuiWindowFlags_ChildWindow != 0) && (flags&ImGuiWindowFlags_Tooltip != 0) // FIXME-WIP: Undocumented behavior of Child+Tooltip for pinned tooltip (#1345)
		window.Active = true
		window.HasCloseButton = (p_open != nil)
		window.ClipRect = ImRect{ImVec2{-FLT_MAX, -FLT_MAX}, ImVec2{+FLT_MAX, +FLT_MAX}}
		window.IDStack = window.IDStack[:1]
		window.DrawList._ResetForNewFrame()
		window.DC.CurrentTableIdx = -1

		// Restore buffer capacity when woken from a compacted state, to avoid
		if window.MemoryCompacted {
			GcAwakeTransientWindowBuffers(window)
		}

		// Update stored window name when it changes (which can _only_ happen with the "###" operator, so the ID would stay unchanged).
		// The title bar always display the 'name' parameter, so we only update the string storage if it needs to be visible to the end-user elsewhere.
		var window_title_visible_elsewhere bool = false
		if g.NavWindowingListWindow != nil && (window.Flags&ImGuiWindowFlags_NoNavFocus) == 0 { // Window titles visible when using CTRL+TAB
			window_title_visible_elsewhere = true
		}
		if window_title_visible_elsewhere && !window_just_created && name == window.Name {
			window.Name = name
		}

		// UPDATE CONTENTS SIZE, UPDATE HIDDEN STATUS

		// Update contents size from last frame for auto-fitting (or use explicit size)
		var window_just_appearing_after_hidden_for_resize bool = (window.HiddenFramesCannotSkipItems > 0)
		CalcWindowContentSizes(window, &window.ContentSize, &window.ContentSizeIdeal)
		if window.HiddenFramesCanSkipItems > 0 {
			window.HiddenFramesCanSkipItems--
		}
		if window.HiddenFramesCannotSkipItems > 0 {
			window.HiddenFramesCannotSkipItems--
		}
		if window.HiddenFramesForRenderOnly > 0 {
			window.HiddenFramesForRenderOnly--
		}

		// Hide new windows for one frame until they calculate their size
		if window_just_created && (!window_size_x_set_by_api || !window_size_y_set_by_api) {
			window.HiddenFramesCannotSkipItems = 1
		}

		// Hide popup/tooltip window when re-opening while we measure size (because we recycle the windows)
		// We reset Size/ContentSize for reappearing popups/tooltips early in this function, so further code won't be tempted to use the old size.
		if window_just_activated_by_user && (flags&(ImGuiWindowFlags_Popup|ImGuiWindowFlags_Tooltip)) != 0 {
			window.HiddenFramesCannotSkipItems = 1
			if flags&ImGuiWindowFlags_AlwaysAutoResize != 0 {
				if !window_size_x_set_by_api {
					window.Size.x = 0
					window.SizeFull.x = 0
				}
				if !window_size_y_set_by_api {
					window.Size.y = 0
					window.SizeFull.y = 0
				}
				window.ContentSize = ImVec2{}
				window.ContentSizeIdeal = ImVec2{}
			}
		}

		// SELECT VIEWPORT
		// FIXME-VIEWPORT: In the docking/viewport branch, this is the point where we select the current viewport (which may affect the style)
		SetCurrentWindow(window)

		// LOCK BORDER SIZE AND PADDING FOR THE FRAME (so that altering them doesn't cause inconsistencies)

		if flags&ImGuiWindowFlags_ChildWindow != 0 {
			window.WindowBorderSize = style.ChildBorderSize
		} else {
			if (flags&(ImGuiWindowFlags_Popup|ImGuiWindowFlags_Tooltip) != 0) && (flags&ImGuiWindowFlags_Modal == 0) {
				window.WindowBorderSize = style.PopupBorderSize
			}
			window.WindowBorderSize = style.WindowBorderSize
		}

		window.WindowPadding = style.WindowPadding
		if (flags&ImGuiWindowFlags_ChildWindow != 0) && (flags&(ImGuiWindowFlags_AlwaysUseWindowPadding|ImGuiWindowFlags_Popup) == 0) && window.WindowBorderSize == 0.0 {
			var y float
			if flags&ImGuiWindowFlags_MenuBar != 0 {
				y = style.WindowPadding.y
			}
			window.WindowPadding = ImVec2{0.0, y}
		}

		// Lock menu offset so size calculation can use it as menu-bar windows need a minimum size.
		window.DC.MenuBarOffset.x = ImMax(ImMax(window.WindowPadding.x, style.ItemSpacing.x), g.NextWindowData.MenuBarOffsetMinVal.x)
		window.DC.MenuBarOffset.y = g.NextWindowData.MenuBarOffsetMinVal.y

		// Collapse window by double-clicking on title bar
		// At this point we don't have a clipping rectangle setup yet, so we can use the title bar area for hit detection and drawing
		if (0 == flags&ImGuiWindowFlags_NoTitleBar) && (0 == flags&ImGuiWindowFlags_NoCollapse) {
			// We don't use a regular button+id to test for double-click on title bar (mostly due to legacy reason, could be fixed), so verify that we don't have items over the title bar.
			var title_bar_rect ImRect = window.TitleBarRect()
			if g.HoveredWindow == window && g.HoveredId == 0 && g.HoveredIdPreviousFrame == 0 && IsMouseHoveringRect(title_bar_rect.Min, title_bar_rect.Max, true) && g.IO.MouseDoubleClicked[0] {
				window.WantCollapseToggle = true
			}
			if window.WantCollapseToggle {
				window.Collapsed = !window.Collapsed
				MarkIniSettingsDirtyWindow(window)
			}
		} else {
			window.Collapsed = false
		}
		window.WantCollapseToggle = false

		// SIZE

		// Calculate auto-fit size, handle automatic resize
		var size_auto_fit ImVec2 = CalcWindowAutoFitSize(window, &window.ContentSizeIdeal)
		var use_current_size_for_scrollbar_x bool = window_just_created
		var use_current_size_for_scrollbar_y bool = window_just_created
		if (flags&ImGuiWindowFlags_AlwaysAutoResize != 0) && !window.Collapsed {
			// Using SetNextWindowSize() overrides ImGuiWindowFlags_AlwaysAutoResize, so it can be used on tooltips/popups, etc.
			if !window_size_x_set_by_api {
				window.SizeFull.x = size_auto_fit.x
				use_current_size_for_scrollbar_x = true
			}
			if !window_size_y_set_by_api {
				window.SizeFull.y = size_auto_fit.y
				use_current_size_for_scrollbar_y = true
			}
		} else if window.AutoFitFramesX > 0 || window.AutoFitFramesY > 0 {
			// Auto-fit may only grow window during the first few frames
			// We still process initial auto-fit on collapsed windows to get a window width, but otherwise don't honor ImGuiWindowFlags_AlwaysAutoResize when collapsed.
			if !window_size_x_set_by_api && window.AutoFitFramesX > 0 {
				if window.AutoFitOnlyGrows {
					window.SizeFull.x = ImMax(window.SizeFull.x, size_auto_fit.x)
				} else {
					window.SizeFull.x = size_auto_fit.x
				}
				use_current_size_for_scrollbar_x = true
			}
			if !window_size_y_set_by_api && window.AutoFitFramesY > 0 {
				if window.AutoFitOnlyGrows {
					window.SizeFull.y = ImMax(window.SizeFull.y, size_auto_fit.y)
				} else {
					window.SizeFull.y = size_auto_fit.y
				}
				use_current_size_for_scrollbar_y = true
			}
			if !window.Collapsed {
				MarkIniSettingsDirtyWindow(window)
			}
		}

		// Apply minimum/maximum window size constraints and final size
		window.SizeFull = CalcWindowSizeAfterConstraint(window, &window.SizeFull)
		if window.Collapsed && 0 == (flags&ImGuiWindowFlags_ChildWindow) {
			titlebar := window.TitleBarRect()
			window.SizeFull = titlebar.GetSize()
		}

		// Decoration size
		var decoration_up_height float = window.TitleBarHeight() + window.MenuBarHeight()

		// POSITION

		// Popup latch its initial position, will position itself when it appears next frame
		if window_just_activated_by_user {
			window.AutoPosLastDirection = ImGuiDir_None
			if (flags&ImGuiWindowFlags_Popup) != 0 && 0 == (flags&ImGuiWindowFlags_Modal) && !window_pos_set_by_api { // FIXME: BeginPopup() could use SetNextWindowPos()
				window.Pos = g.BeginPopupStack[len(g.BeginPopupStack)-1].OpenPopupPos
			}
		}

		// Position child window
		if flags&ImGuiWindowFlags_ChildWindow != 0 {
			IM_ASSERT(parent_window != nil && parent_window.Active)
			window.BeginOrderWithinParent = (short)(len(parent_window.DC.ChildWindows))
			parent_window.DC.ChildWindows = append(parent_window.DC.ChildWindows, window)
			if 0 == (flags&ImGuiWindowFlags_Popup) && !window_pos_set_by_api && !window_is_child_tooltip {
				window.Pos = parent_window.DC.CursorPos
			}
		}

		var window_pos_with_pivot bool = (window.SetWindowPosVal.x != FLT_MAX && window.HiddenFramesCannotSkipItems == 0)
		if window_pos_with_pivot {
			p := window.SetWindowPosVal.Sub(window.Size.Mul(window.SetWindowPosPivot))
			setWindowPos(window, &p, 0) // Position given a pivot (e.g. for centering)
		} else if (flags & ImGuiWindowFlags_ChildMenu) != 0 {
			window.Pos = FindBestWindowPosForPopup(window)
		} else if (flags&ImGuiWindowFlags_Popup) != 0 && !window_pos_set_by_api && window_just_appearing_after_hidden_for_resize {
			window.Pos = FindBestWindowPosForPopup(window)
		} else if (flags&ImGuiWindowFlags_Tooltip) != 0 && !window_pos_set_by_api && !window_is_child_tooltip {
			window.Pos = FindBestWindowPosForPopup(window)
		}

		// Calculate the range of allowed position for that window (to be movable and visible past safe area padding)
		// When clamping to stay visible, we will enforce that window.Pos stays inside of visibility_rect.
		var viewport = GetMainViewport()
		var viewport_rect = ImRect(viewport.GetMainRect())
		var viewport_work_rect = ImRect(viewport.GetWorkRect())
		var visibility_padding = ImMaxVec2(&style.DisplayWindowPadding, &style.DisplaySafeAreaPadding)
		var visibility_rect = ImRect{viewport_work_rect.Min.Add(visibility_padding), viewport_work_rect.Max.Sub(visibility_padding)}

		// Clamp position/size so window stays visible within its viewport or monitor
		// Ignore zero-sized display explicitly to avoid losing positions if a window manager reports zero-sized window when initializing or minimizing.
		if !window_pos_set_by_api && 0 == (flags&ImGuiWindowFlags_ChildWindow) && window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0 {
			if viewport_rect.GetWidth() > 0.0 && viewport_rect.GetHeight() > 0.0 {
				ClampWindowRect(window, &visibility_rect)
			}
		}
		window.Pos = *ImFloorVec(&window.Pos)

		// Lock window rounding for the frame (so that altering them doesn't cause inconsistencies)
		// Large values tend to lead to variety of artifacts and are not recommended.
		if flags&ImGuiWindowFlags_ChildWindow != 0 {
			window.WindowRounding = style.ChildRounding
		} else {
			if (flags&ImGuiWindowFlags_Popup != 0) && 0 == (flags&ImGuiWindowFlags_Modal) {
				window.WindowRounding = style.PopupRounding
			} else {
				window.WindowRounding = style.WindowRounding
			}
		}

		// For windows with title bar or menu bar, we clamp to FrameHeight(FontSize + FramePadding.y * 2.0f) to completely hide artifacts.
		//if ((window.Flags & ImGuiWindowFlags_MenuBar) || !(window.Flags & ImGuiWindowFlags_NoTitleBar))
		//    window.WindowRounding = ImMin(window.WindowRounding, g.FontSize + style.FramePadding.y * 2.0f);

		// Apply window focus (new and reactivated windows are moved to front)
		var want_focus bool = false
		if window_just_activated_by_user && 0 == (flags&ImGuiWindowFlags_NoFocusOnAppearing) {
			if flags&ImGuiWindowFlags_Popup != 0 {
				want_focus = true
			} else if (flags & (ImGuiWindowFlags_ChildWindow | ImGuiWindowFlags_Tooltip)) == 0 {
				want_focus = true
			}
		}

		// Handle manual resize: Resize Grips, Borders, Gamepad
		var border_held int = -1
		var resize_grip_col [4]ImU32

		// Allow resize from lower-left if we have the mouse cursor feedback for it.
		var resize_grip_count int = 1
		if g.IO.ConfigWindowsResizeFromEdges {
			resize_grip_count = 2
		}
		var resize_grip_draw_size float = IM_FLOOR(ImMax(g.FontSize*1.10, window.WindowRounding+1.0+g.FontSize*0.2))
		if !window.Collapsed {
			if UpdateWindowManualResize(window, &size_auto_fit, &border_held, resize_grip_count, &resize_grip_col, &visibility_rect) {
				use_current_size_for_scrollbar_x = true
				use_current_size_for_scrollbar_y = true
			}
		}
		window.ResizeBorderHeld = int8((byte)(border_held))

		// SCROLLBAR VISIBILITY

		// Update scrollbar visibility (based on the Size that was effective during last frame or the auto-resized Size).
		if !window.Collapsed {
			// When reading the current size we need to read it after size constraints have been applied.
			// When we use InnerRect here we are intentionally reading last frame size, same for ScrollbarSizes values before we set them again.
			var avail_size_from_current_frame ImVec2 = ImVec2{window.SizeFull.x, window.SizeFull.y - decoration_up_height}
			var avail_size_from_last_frame ImVec2 = window.InnerRect.GetSize().Add(window.ScrollbarSizes)
			var needed_size_from_last_frame ImVec2
			if !window_just_created {
				needed_size_from_last_frame = window.ContentSize.Add(window.WindowPadding.Scale(2.0))
			}
			var size_x_for_scrollbars float = avail_size_from_last_frame.x
			if use_current_size_for_scrollbar_x {
				size_x_for_scrollbars = avail_size_from_current_frame.x
			}
			var size_y_for_scrollbars float = avail_size_from_last_frame.y
			if use_current_size_for_scrollbar_y {
				size_x_for_scrollbars = avail_size_from_current_frame.y
			}
			//bool scrollbar_y_from_last_frame = window.ScrollbarY; // FIXME: May want to use that in the ScrollbarX expression? How many pros vs cons?
			window.ScrollbarY = (flags&ImGuiWindowFlags_AlwaysVerticalScrollbar != 0) || ((needed_size_from_last_frame.y > size_y_for_scrollbars) && 0 == (flags&ImGuiWindowFlags_NoScrollbar))

			var scroll_bar_size_y float
			if window.ScrollbarY {
				scroll_bar_size_y = style.ScrollbarSize
			}

			window.ScrollbarX = (flags&ImGuiWindowFlags_AlwaysHorizontalScrollbar != 0) || ((needed_size_from_last_frame.x > size_x_for_scrollbars-(scroll_bar_size_y)) && 0 == (flags&ImGuiWindowFlags_NoScrollbar) && (flags&ImGuiWindowFlags_HorizontalScrollbar != 0))
			if window.ScrollbarX && !window.ScrollbarY {
				window.ScrollbarY = (needed_size_from_last_frame.y > size_y_for_scrollbars) && 0 == (flags&ImGuiWindowFlags_NoScrollbar)
			}
			var scroll_bar_size_x float
			if window.ScrollbarX {
				scroll_bar_size_x = style.ScrollbarSize
			}
			window.ScrollbarSizes = ImVec2{scroll_bar_size_y, scroll_bar_size_x}
		}

		// UPDATE RECTANGLES (1- THOSE NOT AFFECTED BY SCROLLING)
		// Update various regions. Variables they depends on should be set above in this function.
		// We set this up after processing the resize grip so that our rectangles doesn't lag by a frame.

		// Outer rectangle
		// Not affected by window border size. Used by:
		// - FindHoveredWindow() (w/ extra padding when border resize is enabled)
		// - Begin() initial clipping rect for drawing window background and borders.
		// - Begin() clipping whole child
		var host_rect ImRect = viewport_rect
		if (flags&ImGuiWindowFlags_ChildWindow != 0) && 0 == (flags&ImGuiWindowFlags_Popup) && !window_is_child_tooltip {
			host_rect = parent_window.ClipRect
		}
		var outer_rect ImRect = window.Rect()
		var title_bar_rect ImRect = window.TitleBarRect()
		window.OuterRectClipped = outer_rect
		window.OuterRectClipped.ClipWith(host_rect)

		// Inner rectangle
		// Not affected by window border size. Used by:
		// - InnerClipRect
		// - ScrollToBringRectIntoView()
		// - NavUpdatePageUpPageDown()
		// - Scrollbar()
		window.InnerRect.Min.x = window.Pos.x
		window.InnerRect.Min.y = window.Pos.y + decoration_up_height
		window.InnerRect.Max.x = window.Pos.x + window.Size.x - window.ScrollbarSizes.x
		window.InnerRect.Max.y = window.Pos.y + window.Size.y - window.ScrollbarSizes.y

		// Inner clipping rectangle.
		// Will extend a little bit outside the normal work region.
		// This is to allow e.g. Selectable or CollapsingHeader or some separators to cover that space.
		// Force round operator last to ensure that e.g. (int)(max.x-min.x) in user's render code produce correct result.
		// Note that if our window is collapsed we will end up with an inverted (~nil) clipping rectangle which is the correct behavior.
		// Affected by window/frame border size. Used by:
		// - Begin() initial clip rect
		var top_border_size float
		if (flags&ImGuiWindowFlags_MenuBar != 0) || 0 == (flags&ImGuiWindowFlags_NoTitleBar) {
			top_border_size = style.FrameBorderSize
		}
		window.InnerClipRect.Min.x = ImFloor(0.5 + window.InnerRect.Min.x + ImMax(ImFloor(window.WindowPadding.x*0.5), window.WindowBorderSize))
		window.InnerClipRect.Min.y = ImFloor(0.5 + window.InnerRect.Min.y + top_border_size)
		window.InnerClipRect.Max.x = ImFloor(0.5 + window.InnerRect.Max.x - ImMax(ImFloor(window.WindowPadding.x*0.5), window.WindowBorderSize))
		window.InnerClipRect.Max.y = ImFloor(0.5 + window.InnerRect.Max.y - window.WindowBorderSize)
		window.InnerClipRect.ClipWithFull(host_rect)

		// Default item width. Make it proportional to window size if window manually resizes
		if window.Size.x > 0.0 && 0 == (flags&ImGuiWindowFlags_Tooltip) && 0 == (flags&ImGuiWindowFlags_AlwaysAutoResize) {
			window.ItemWidthDefault = ImFloor(window.Size.x * 0.65)
		} else {
			window.ItemWidthDefault = ImFloor(g.FontSize * 16.0)
		}

		// SCROLLING

		// Lock down maximum scrolling
		// The value of ScrollMax are ahead from ScrollbarX/ScrollbarY which is intentionally using InnerRect from previous rect in order to accommodate
		// for right/bottom aligned items without creating a scrollbar.
		window.ScrollMax.x = ImMax(0.0, window.ContentSize.x+window.WindowPadding.x*2.0-window.InnerRect.GetWidth())
		window.ScrollMax.y = ImMax(0.0, window.ContentSize.y+window.WindowPadding.y*2.0-window.InnerRect.GetHeight())

		// Apply scrolling
		window.Scroll = CalcNextScrollFromScrollTargetAndClamp(window)
		window.ScrollTarget = ImVec2{FLT_MAX, FLT_MAX}

		// DRAWING

		// Setup draw list and outer clipping rectangle
		IM_ASSERT(len(window.DrawList.CmdBuffer) == 1 && window.DrawList.CmdBuffer[0].ElemCount == 0)
		window.DrawList.PushTextureID(g.Font.ContainerAtlas.TexID)
		PushClipRect(host_rect.Min, host_rect.Max, false)

		// Draw modal window background (darkens what is behind them, all viewports)
		var dim_bg_for_modal bool = (flags&ImGuiWindowFlags_Modal != 0) && window == GetTopMostPopupModal() && window.HiddenFramesCannotSkipItems <= 0
		var dim_bg_for_window_list = g.NavWindowingTargetAnim != nil && (window == g.NavWindowingTargetAnim.RootWindow)
		if dim_bg_for_modal || dim_bg_for_window_list {

			var c ImGuiCol
			if dim_bg_for_modal {
				c = ImGuiCol_ModalWindowDimBg
			} else {
				c = ImGuiCol_NavWindowingDimBg
			}

			var dim_bg_col ImU32 = GetColorU32FromID(c, g.DimBgRatio)
			window.DrawList.AddRectFilled(viewport_rect.Min, viewport_rect.Max, dim_bg_col, 0, 0)
		}

		// Draw navigation selection/windowing rectangle background
		if dim_bg_for_window_list && window == g.NavWindowingTargetAnim {
			var bb ImRect = window.Rect()
			bb.Expand(g.FontSize)
			if !bb.ContainsRect(viewport_rect) { // Avoid drawing if the window covers all the viewport anyway
				window.DrawList.AddRectFilled(bb.Min, bb.Max, GetColorU32FromID(ImGuiCol_NavWindowingHighlight, g.NavWindowingHighlightAlpha*0.25), g.Style.WindowRounding, 0)
			}
		}

		// Child windows can render their decoration (bg color, border, scrollbars, etc.) within their parent to save a draw call (since 1.71)
		// When using overlapping child windows, this will break the assumption that child z-order is mapped to submission order.
		// FIXME: User code may rely on explicit sorting of overlapping child window and would need to disable this somehow. Please get in contact if you are affected (github #4493)
		{
			var render_decorations_in_parent bool = false
			if (flags&ImGuiWindowFlags_ChildWindow != 0) && 0 == (flags&ImGuiWindowFlags_Popup) && !window_is_child_tooltip {
				// - We test overlap with the previous child window only (testing all would end up being O(log N) not a good investment here)
				// - We disable this when the parent window has zero vertices, which is a common pattern leading to laying out multiple overlapping childs
				var previous_child *ImGuiWindow
				var previous_child_overlapping bool
				if len(parent_window.DC.ChildWindows) >= 2 {
					previous_child = parent_window.DC.ChildWindows[len(parent_window.DC.ChildWindows)-2]

				}
				if previous_child != nil {
					pr := previous_child.Rect()
					previous_child_overlapping = pr.Overlaps(window.Rect())
				}
				var parent_is_empty bool = len(parent_window.DrawList.VtxBuffer) > 0
				if window.DrawList.CmdBuffer[len(window.DrawList.CmdBuffer)-1].ElemCount == 0 && parent_is_empty && !previous_child_overlapping {
					render_decorations_in_parent = true
				}
			}
			if render_decorations_in_parent {
				window.DrawList = parent_window.DrawList
			}

			// Handle title bar, scrollbar, resize grips and resize borders
			var window_to_highlight *ImGuiWindow = g.NavWindow
			if g.NavWindowingTarget != nil {
				window_to_highlight = g.NavWindowingTarget
			}
			var title_bar_is_highlight bool = want_focus || (window_to_highlight != nil && window.RootWindowForTitleBarHighlight == window_to_highlight.RootWindowForTitleBarHighlight)
			RenderWindowDecorations(window, &title_bar_rect, title_bar_is_highlight, resize_grip_count, resize_grip_col, resize_grip_draw_size)

			if render_decorations_in_parent {
				window.DrawList = &window.DrawListInst
			}
		}

		// Draw navigation selection/windowing rectangle border
		if g.NavWindowingTargetAnim == window {
			var rounding float = ImMax(window.WindowRounding, g.Style.WindowRounding)
			var bb ImRect = window.Rect()
			bb.Expand(g.FontSize)
			if bb.ContainsRect(viewport_rect) { // If a window fits the entire viewport, adjust its highlight inward
				bb.Expand(-g.FontSize - 1.0)
				rounding = window.WindowRounding
			}
			window.DrawList.AddRect(bb.Min, bb.Max, GetColorU32FromID(ImGuiCol_NavWindowingHighlight, g.NavWindowingHighlightAlpha), rounding, 0, 3.0)
		}

		// UPDATE RECTANGLES (2- THOSE AFFECTED BY SCROLLING)

		// Work rectangle.
		// Affected by window padding and border size. Used by:
		// - Columns() for right-most edge
		// - TreeNode(), CollapsingHeader() for right-most edge
		// - BeginTabBar() for right-most edge
		var allow_scrollbar_x bool = 0 == (flags&ImGuiWindowFlags_NoScrollbar) && (flags&ImGuiWindowFlags_HorizontalScrollbar != 0)
		var allow_scrollbar_y bool = 0 == (flags & ImGuiWindowFlags_NoScrollbar)

		var work_rect_size_x float
		if window.ContentSizeExplicit.x != 0.0 {
			work_rect_size_x = window.ContentSizeExplicit.x
		} else {
			var a float
			if allow_scrollbar_x {
				a = window.ContentSize.x
			}
			work_rect_size_x = ImMax(a, window.Size.x-window.WindowPadding.x*2.0-window.ScrollbarSizes.x)
		}

		var work_rect_size_y float
		if window.ContentSizeExplicit.y != 0.0 {
			work_rect_size_y = window.ContentSizeExplicit.y
		} else {
			var a float
			if allow_scrollbar_y {
				a = window.ContentSize.y
			}
			work_rect_size_y = ImMax(a, window.Size.y-window.WindowPadding.y*2.0-window.ScrollbarSizes.y)
		}

		window.WorkRect.Min.x = ImFloor(window.InnerRect.Min.x - window.Scroll.x + ImMax(window.WindowPadding.x, window.WindowBorderSize))
		window.WorkRect.Min.y = ImFloor(window.InnerRect.Min.y - window.Scroll.y + ImMax(window.WindowPadding.y, window.WindowBorderSize))
		window.WorkRect.Max.x = window.WorkRect.Min.x + work_rect_size_x
		window.WorkRect.Max.y = window.WorkRect.Min.y + work_rect_size_y
		window.ParentWorkRect = window.WorkRect

		// [LEGACY] Content Region
		// FIXME-OBSOLETE: window.ContentRegionRect.Max is currently very misleading / partly faulty, but some BeginChild() patterns relies on it.
		// Used by:
		// - Mouse wheel scrolling + many other things
		window.ContentRegionRect.Min.x = window.Pos.x - window.Scroll.x + window.WindowPadding.x
		window.ContentRegionRect.Min.y = window.Pos.y - window.Scroll.y + window.WindowPadding.y + decoration_up_height

		var explicit_x float
		if window.ContentSizeExplicit.x != 0.0 {
			explicit_x = window.ContentSizeExplicit.x
		} else {
			explicit_x = window.Size.x - window.WindowPadding.x*2.0 - window.ScrollbarSizes.x
		}
		var explicit_y float
		if window.ContentSizeExplicit.y != 0.0 {
			explicit_y = window.ContentSizeExplicit.y
		} else {
			explicit_y = window.Size.y - window.WindowPadding.y*2.0 - window.ScrollbarSizes.y
		}

		window.ContentRegionRect.Max.x = window.ContentRegionRect.Min.x + explicit_x
		window.ContentRegionRect.Max.y = window.ContentRegionRect.Min.y + explicit_y

		// Setup drawing context
		// (NB: That term "drawing context / DC" lost its meaning a long time ago. Initially was meant to hold transient data only. Nowadays difference between window. and window.DC. is dubious.)
		window.DC.Indent.x = 0.0 + window.WindowPadding.x - window.Scroll.x
		window.DC.GroupOffset.x = 0.0
		window.DC.ColumnsOffset.x = 0.0
		window.DC.CursorStartPos = window.Pos.Add(ImVec2{window.DC.Indent.x + window.DC.ColumnsOffset.x, decoration_up_height + window.WindowPadding.y - window.Scroll.y})
		window.DC.CursorPos = window.DC.CursorStartPos
		window.DC.CursorPosPrevLine = window.DC.CursorPos
		window.DC.CursorMaxPos = window.DC.CursorStartPos
		window.DC.IdealMaxPos = window.DC.CursorStartPos
		window.DC.CurrLineSize = ImVec2{}
		window.DC.PrevLineSize = ImVec2{}
		window.DC.CurrLineTextBaseOffset = 0
		window.DC.PrevLineTextBaseOffset = 0

		window.DC.NavLayerCurrent = ImGuiNavLayer_Main
		window.DC.NavLayersActiveMask = window.DC.NavLayersActiveMaskNext
		window.DC.NavLayersActiveMaskNext = 0x00
		window.DC.NavHideHighlightOneFrame = false
		window.DC.NavHasScroll = (window.ScrollMax.y > 0.0)

		window.DC.MenuBarAppending = false
		window.DC.MenuColumns.Update(style.ItemSpacing.x, window_just_activated_by_user)
		window.DC.TreeDepth = 0
		window.DC.TreeJumpToParentOnPopMask = 0x00
		window.DC.ChildWindows = window.DC.ChildWindows[:0]
		window.DC.StateStorage = window.StateStorage
		window.DC.CurrentColumns = nil
		window.DC.LayoutType = ImGuiLayoutType_Vertical
		if parent_window != nil {
			window.DC.ParentLayoutType = parent_window.DC.LayoutType
		} else {
			window.DC.ParentLayoutType = ImGuiLayoutType_Vertical
		}
		window.DC.FocusCounterRegular = -1
		window.DC.FocusCounterTabStop = -1

		window.DC.ItemWidth = window.ItemWidthDefault
		window.DC.TextWrapPos = -1.0 // disabled
		window.DC.ItemWidthStack = window.DC.ItemWidthStack[:0]
		window.DC.TextWrapPosStack = window.DC.TextWrapPosStack[:0]

		if window.AutoFitFramesX > 0 {
			window.AutoFitFramesX--
		}
		if window.AutoFitFramesY > 0 {
			window.AutoFitFramesY--
		}

		// Apply focus (we need to call FocusWindow() AFTER setting DC.CursorStartPos so our initial navigation reference rectangle can start around there)
		if want_focus {
			FocusWindow(window)
			NavInitWindow(window, false) // <-- this is in the way for us to be able to defer and sort reappearing FocusWindow() calls
		}

		// Title bar
		if 0 == (flags & ImGuiWindowFlags_NoTitleBar) {
			RenderWindowTitleBarContents(window, &ImRect{ImVec2{title_bar_rect.Min.x + window.WindowBorderSize, title_bar_rect.Min.y}, ImVec2{title_bar_rect.Max.x - window.WindowBorderSize, title_bar_rect.Max.y}}, name, p_open)
		}

		// Clear hit test shape every frame
		window.HitTestHoleSize.x = 0
		window.HitTestHoleSize.y = 0

		// Pressing CTRL+C while holding on a window copy its content to the clipboard
		// This works but 1. doesn't handle multiple Begin/End pairs, 2. recursing into another Begin/End pair - so we need to work that out and add better logging scope.
		// Maybe we can support CTRL+C on every element?
		/*
		   //if (g.NavWindow == window && g.ActiveId == 0)
		   if (g.ActiveId == window.MoveId)
		       if (g.IO.KeyCtrl && IsKeyPressedMap(ImGuiKey_C))
		           LogToClipboard();
		*/

		// We fill last item data based on Title Bar/Tab, in order for IsItemHovered() and IsItemActive() to be usable after Begin().
		// This is useful to allow creating context menus on title bar only, etc.
		g.LastItemData.ID = window.MoveId
		g.LastItemData.InFlags = g.CurrentItemFlags
		if IsMouseHoveringRect(title_bar_rect.Min, title_bar_rect.Max, false) {
			g.LastItemData.StatusFlags = ImGuiItemStatusFlags_HoveredRect
		} else {
			g.LastItemData.StatusFlags = 0
		}
		g.LastItemData.Rect = title_bar_rect
	} else {
		// Append
		SetCurrentWindow(window)
	}

	// Pull/inherit current state
	if flags&ImGuiWindowFlags_ChildWindow != 0 {
		window.DC.NavFocusScopeIdCurrent = parent_window.DC.NavFocusScopeIdCurrent
	} else {
		window.GetIDs("#FOCUSSCOPE", "")
	}

	PushClipRect(window.InnerClipRect.Min, window.InnerClipRect.Max, true)

	// Clear 'accessed' flag last thing (After PushClipRect which will set the flag. We want the flag to stay false when the default "Debug" window is unused)
	window.WriteAccessed = false
	window.BeginCount++
	g.NextWindowData.ClearFlags()

	// Update visibility
	if first_begin_of_the_frame {
		if flags&ImGuiWindowFlags_ChildWindow != 0 {
			// Child window can be out of sight and have "negative" clip windows.
			// Mark them as collapsed so commands are skipped earlier (we can't manually collapse them because they have no title bar).
			IM_ASSERT((flags & ImGuiWindowFlags_NoTitleBar) != 0)
			if 0 == (flags&ImGuiWindowFlags_AlwaysAutoResize) && window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0 { // FIXME: Doesn't make sense for ChildWindow??
				if !g.LogEnabled {
					if window.OuterRectClipped.Min.x >= window.OuterRectClipped.Max.x || window.OuterRectClipped.Min.y >= window.OuterRectClipped.Max.y {
						window.HiddenFramesCanSkipItems = 1
					}
				}
			}

			// Hide along with parent or if parent is collapsed
			if parent_window != nil && (parent_window.Collapsed || parent_window.HiddenFramesCanSkipItems > 0) {
				window.HiddenFramesCanSkipItems = 1
			}
			if parent_window != nil && (parent_window.Collapsed || parent_window.HiddenFramesCannotSkipItems > 0) {
				window.HiddenFramesCannotSkipItems = 1
			}
		}

		// Don't render if style alpha is 0.0 at the time of Begin(). This is arbitrary and inconsistent but has been there for a long while (may remove at some point)
		if style.Alpha <= 0.0 {
			window.HiddenFramesCanSkipItems = 1
		}

		// Update the Hidden flag
		window.Hidden = (window.HiddenFramesCanSkipItems > 0) || (window.HiddenFramesCannotSkipItems > 0) || (window.HiddenFramesForRenderOnly > 0)

		// Disable inputs for requested number of frames
		if window.DisableInputsFrames > 0 {
			window.DisableInputsFrames--
			window.Flags |= ImGuiWindowFlags_NoInputs
		}

		// Update the SkipItems flag, used to early out of all items functions (no layout required)
		var skip_items bool = false
		if window.Collapsed || !window.Active || window.Hidden {
			if window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0 && window.HiddenFramesCannotSkipItems <= 0 {
				skip_items = true
			}
		}
		window.SkipItems = skip_items
	}

	return !window.SkipItems
}
