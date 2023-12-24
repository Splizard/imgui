package imgui

const WINDOWS_RESIZE_FROM_EDGES_FEEDBACK_TIMER float = 0.04 // Reduce visual noise by only highlighting the border after a certain time.

// Data for resizing from corner
type ImGuiResizeGripDef struct {
	CornerPosN             ImVec2
	InnerDir               ImVec2
	AngleMin12, AngleMax12 int
}

var resize_grip_def = [4]ImGuiResizeGripDef{
	{ImVec2{1, 1}, ImVec2{-1, -1}, 0, 3},  // Lower-right
	{ImVec2{0, 1}, ImVec2{+1, -1}, 3, 6},  // Lower-left
	{ImVec2{0, 0}, ImVec2{+1, +1}, 6, 9},  // Upper-left (Unused)
	{ImVec2{1, 0}, ImVec2{-1, +1}, 9, 12}, // Upper-right (Unused)
}

// Data for resizing from borders
type ImGuiResizeBorderDef struct {
	InnerDir             ImVec2
	SegmentN1, SegmentN2 ImVec2
	OuterAngle           float
}

var resize_border_def = [4]ImGuiResizeBorderDef{
	{ImVec2{+1, 0}, ImVec2{0, 1}, ImVec2{0, 0}, IM_PI * 1.00}, // Left
	{ImVec2{-1, 0}, ImVec2{1, 0}, ImVec2{1, 1}, IM_PI * 0.00}, // Right
	{ImVec2{0, +1}, ImVec2{0, 0}, ImVec2{1, 0}, IM_PI * 1.50}, // Up
	{ImVec2{0, -1}, ImVec2{1, 1}, ImVec2{0, 1}, IM_PI * 0.50}, // Down
}

func CalcResizePosSizeFromAnyCorner(window *ImGuiWindow, corner_target, corner_norm *ImVec2, out_pos, out_size *ImVec2) {
	a := window.Pos.Add(window.Size)
	var pos_min = ImLerpVec2WithVec2(corner_target, &window.Pos, *corner_norm) // Expected window upper-left
	var pos_max = ImLerpVec2WithVec2(&a, corner_target, *corner_norm)          // Expected window lower-right
	var size_expected = pos_max.Sub(pos_min)
	var size_constrained = CalcWindowSizeAfterConstraint(window, &size_expected)
	*out_pos = pos_min
	if corner_norm.x == 0.0 {
		out_pos.x -= (size_constrained.x - size_expected.x)
	}
	if corner_norm.y == 0.0 {
		out_pos.y -= (size_constrained.y - size_expected.y)
	}
	*out_size = size_constrained
}

func GetResizeBorderRect(window *ImGuiWindow, border_n int, perp_padding, thickness float) ImRect {
	var rect = window.Rect()
	if thickness == 0.0 {
		rect.Max = rect.Max.Sub(ImVec2{1, 1})
	}
	if border_n == int(ImGuiDir_Left) {
		return ImRect{ImVec2{rect.Min.x - thickness, rect.Min.y + perp_padding}, ImVec2{rect.Min.x + thickness, rect.Max.y - perp_padding}}
	}
	if border_n == int(ImGuiDir_Right) {
		return ImRect{ImVec2{rect.Max.x - thickness, rect.Min.y + perp_padding}, ImVec2{rect.Max.x + thickness, rect.Max.y - perp_padding}}
	}
	if border_n == int(ImGuiDir_Up) {
		return ImRect{ImVec2{rect.Min.x + perp_padding, rect.Min.y - thickness}, ImVec2{rect.Max.x - perp_padding, rect.Min.y + thickness}}
	}
	if border_n == int(ImGuiDir_Down) {
		return ImRect{ImVec2{rect.Min.x + perp_padding, rect.Max.y - thickness}, ImVec2{rect.Max.x - perp_padding, rect.Max.y + thickness}}
	}
	IM_ASSERT(false)
	return ImRect{}
}

// Handle resize for: Resize Grips, Borders, Gamepad
// Return true when using auto-fit (double click on resize grip)
func UpdateWindowManualResize(window *ImGuiWindow, size_auto_fit *ImVec2, border_held *int, resize_grip_count int, resize_grip_col *[4]ImU32, visibility_rect *ImRect) bool {
	var flags = window.Flags

	if (flags&ImGuiWindowFlags_NoResize != 0) || (flags&ImGuiWindowFlags_AlwaysAutoResize != 0) || window.AutoFitFramesX > 0 || window.AutoFitFramesY > 0 {
		return false
	}
	if !window.WasActive { // Early out to avoid running this code for e.guiContext. an hidden implicit/fallback Debug window.
		return false
	}

	var ret_auto_fit = false
	var resize_border_count int
	if guiContext.IO.ConfigWindowsResizeFromEdges {
		resize_border_count = 4
	}
	var grip_draw_size = IM_FLOOR(max(guiContext.FontSize*1.35, window.WindowRounding+1.0+guiContext.FontSize*0.2))
	var grip_hover_inner_size = IM_FLOOR(grip_draw_size * 0.75)
	var grip_hover_outer_size float
	if guiContext.IO.ConfigWindowsResizeFromEdges {
		grip_hover_outer_size = WINDOWS_HOVER_PADDING
	}

	var pos_target = ImVec2{FLT_MAX, FLT_MAX}
	var size_target = ImVec2{FLT_MAX, FLT_MAX}

	// Resize grips and borders are on layer 1
	window.DC.NavLayerCurrent = ImGuiNavLayer_Menu

	// Manual resize grips
	PushString("#RESIZE")
	for resize_grip_n := int(0); resize_grip_n < resize_grip_count; resize_grip_n++ {
		var def = &resize_grip_def[resize_grip_n]

		size := window.Pos.Add(window.Size)
		var corner = ImLerpVec2WithVec2(&window.Pos, &size, def.CornerPosN)

		// Using the FlattenChilds button flag we make the resize button accessible even if we are hovering over a child window
		var hovered, held bool
		var resize_rect = ImRect{corner.Sub(def.InnerDir.Scale(grip_hover_outer_size)), corner.Add(def.InnerDir.Scale(grip_hover_inner_size))}
		if resize_rect.Min.x > resize_rect.Max.x {
			resize_rect.Min.x, resize_rect.Max.x = resize_rect.Max.x, resize_rect.Min.x
		}
		if resize_rect.Min.y > resize_rect.Max.y {
			resize_rect.Min.y, resize_rect.Max.y = resize_rect.Max.y, resize_rect.Min.y
		}
		var resize_grip_id = window.GetIDInt(resize_grip_n) // == GetWindowResizeCornerID()
		ButtonBehavior(&resize_rect, resize_grip_id, &hovered, &held, ImGuiButtonFlags_FlattenChildren|ImGuiButtonFlags_NoNavFocus)
		//GetForegroundDrawList(window).AddRect(resize_rect.Min, resize_rect.Max, IM_COL32(255, 255, 0, 255));
		if hovered || held {
			if resize_grip_n&1 != 0 {
				guiContext.MouseCursor = ImGuiMouseCursor_ResizeNESW
			} else {
				guiContext.MouseCursor = ImGuiMouseCursor_ResizeNWSE
			}
		}

		if held && guiContext.IO.MouseDoubleClicked[0] && resize_grip_n == 0 {
			// Manual auto-fit when double-clicking
			size_target = CalcWindowSizeAfterConstraint(window, size_auto_fit)
			ret_auto_fit = true
			ClearActiveID()
		} else if held {
			// Resize from any of the four corners
			// We don't use an incremental MouseDelta but rather compute an absolute target size based on mouse position
			var clamp_min = ImVec2{-FLT_MAX, -FLT_MAX}
			if def.CornerPosN.x == 1.0 {
				clamp_min.x = visibility_rect.Min.x
			}
			if def.CornerPosN.y == 1.0 {
				clamp_min.y = visibility_rect.Min.y
			}
			var clamp_max = ImVec2{+FLT_MAX, +FLT_MAX}
			if def.CornerPosN.x == 0 {
				clamp_max.x = visibility_rect.Max.x
			}
			if def.CornerPosN.y == 0 {
				clamp_max.y = visibility_rect.Max.y
			}

			ls := def.InnerDir.Scale(grip_hover_outer_size)
			rs := def.InnerDir.Scale(-grip_hover_inner_size)

			var corner_target = guiContext.IO.MousePos.Sub(guiContext.ActiveIdClickOffset).Add(
				ImLerpVec2WithVec2(&ls, &rs, def.CornerPosN)) // Corner of the window corresponding to our corner grip
			corner_target = ImClampVec2(&corner_target, &clamp_min, clamp_max)
			CalcResizePosSizeFromAnyCorner(window, &corner_target, &def.CornerPosN, &pos_target, &size_target)
		}

		// Only lower-left grip is visible before hovering/activating
		if resize_grip_n == 0 || held || hovered {
			var c = ImGuiCol_ResizeGrip
			if held {
				c = ImGuiCol_ResizeGripActive
			} else {
				if hovered {
					c = ImGuiCol_ResizeGripHovered
				}
			}
			resize_grip_col[resize_grip_n] = GetColorU32FromID(c, 1)
		}
	}
	for border_n := ImGuiDir(0); border_n < ImGuiDir(resize_border_count); border_n++ {
		var def = &resize_border_def[border_n]
		var axis = ImGuiAxis_Y
		if border_n == ImGuiDir_Left || border_n == ImGuiDir_Right {
			axis = ImGuiAxis_X
		}

		var hovered, held bool
		var border_rect = GetResizeBorderRect(window, int(border_n), grip_hover_inner_size, WINDOWS_HOVER_PADDING)
		var border_id = window.GetIDInt(int(border_n) + 4) // == GetWindowResizeBorderID()
		ButtonBehavior(&border_rect, border_id, &hovered, &held, ImGuiButtonFlags_FlattenChildren)
		//GetForegroundDrawLists(window).AddRect(border_rect.Min, border_rect.Max, IM_COL32(255, 255, 0, 255));
		if (hovered && guiContext.HoveredIdTimer > WINDOWS_RESIZE_FROM_EDGES_FEEDBACK_TIMER) || held {
			if axis == ImGuiAxis_X {
				guiContext.MouseCursor = ImGuiMouseCursor_ResizeEW
			} else {
				guiContext.MouseCursor = ImGuiMouseCursor_ResizeNS
			}
			if held {
				*border_held = int(border_n)
			}
		}
		if held {
			var clamp_min = ImVec2{-FLT_MAX, -FLT_MAX}
			if border_n == ImGuiDir_Right {
				clamp_min.x = visibility_rect.Min.x
			}
			if border_n == ImGuiDir_Down {
				clamp_min.y = visibility_rect.Min.y
			}
			var clamp_max = ImVec2{+FLT_MAX, +FLT_MAX}
			if border_n == ImGuiDir_Left {
				clamp_max.x = visibility_rect.Max.x
			}
			if border_n == ImGuiDir_Up {
				clamp_max.y = visibility_rect.Max.y
			}

			var border_target = window.Pos
			switch axis {
			case ImGuiAxis_X:
				border_target.x = guiContext.IO.MousePos.x - guiContext.ActiveIdClickOffset.x + WINDOWS_HOVER_PADDING
			case ImGuiAxis_Y:
				border_target.y = guiContext.IO.MousePos.y - guiContext.ActiveIdClickOffset.y + WINDOWS_HOVER_PADDING
			}

			border_target = ImClampVec2(&border_target, &clamp_min, clamp_max)

			min := ImMinVec2(&def.SegmentN1, &def.SegmentN2)
			CalcResizePosSizeFromAnyCorner(window, &border_target, &min, &pos_target, &size_target)
		}
	}
	PopID()

	// Restore nav layer
	window.DC.NavLayerCurrent = ImGuiNavLayer_Main

	// Navigation resize (keyboard/gamepad)
	if guiContext.NavWindowingTarget != nil && guiContext.NavWindowingTarget.RootWindow == window {
		var nav_resize_delta ImVec2
		if guiContext.NavInputSource == ImGuiInputSource_Keyboard && guiContext.IO.KeyShift {
			nav_resize_delta = GetNavInputAmount2d(ImGuiNavDirSourceFlags_Keyboard, ImGuiInputReadMode_Down, 0, 0)
		}
		if guiContext.NavInputSource == ImGuiInputSource_Gamepad {
			nav_resize_delta = GetNavInputAmount2d(ImGuiNavDirSourceFlags_PadDPad, ImGuiInputReadMode_Down, 0, 0)
		}
		if nav_resize_delta.x != 0.0 || nav_resize_delta.y != 0.0 {
			const NAV_RESIZE_SPEED float = 600
			nav_resize_delta = nav_resize_delta.Scale(ImFloor(NAV_RESIZE_SPEED * guiContext.IO.DeltaTime * min(guiContext.IO.DisplayFramebufferScale.x, guiContext.IO.DisplayFramebufferScale.y)))
			yd := visibility_rect.Min.Sub(window.Pos).Sub(window.Size)
			nav_resize_delta = ImMaxVec2(&nav_resize_delta, &yd)
			guiContext.NavWindowingToggleLayer = false
			guiContext.NavDisableMouseHover = true
			resize_grip_col[0] = GetColorU32FromInt(uint(ImGuiCol_ResizeGripActive))
			// FIXME-NAV: Should store and accumulate into a separate size buffer to handle sizing constraints properly, right now a constraint will make us stuck.

			desired := window.SizeFull.Add(nav_resize_delta)
			size_target = CalcWindowSizeAfterConstraint(window, &desired)
		}
	}

	// Apply back modified position/size to window
	if size_target.x != FLT_MAX {
		window.SizeFull = size_target
		MarkIniSettingsDirtyWindow(window)
	}
	if pos_target.x != FLT_MAX {
		window.Pos = *ImFloorVec(&pos_target)
		MarkIniSettingsDirtyWindow(window)
	}

	window.Size = window.SizeFull
	return ret_auto_fit
}
