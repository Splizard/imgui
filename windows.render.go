package imgui

// Draw background and borders
// Draw and handle scrollbars
func RenderWindowDecorations(window *ImGuiWindow, title_bar_rect *ImRect, title_bar_is_highlight bool, resize_grip_count int, resize_grip_col [4]ImU32, resize_grip_draw_size float) {
	var g = GImGui
	var style = g.Style
	var flags = window.Flags

	// Ensure that ScrollBar doesn't read last frame's SkipItems
	IM_ASSERT(window.BeginCount == 0)
	window.SkipItems = false

	// Draw window + handle manual resize
	// As we highlight the title bar when want_focus is set, multiple reappearing windows will have have their title bar highlighted on their reappearing frame.
	var window_rounding float = window.WindowRounding
	var window_border_size float = window.WindowBorderSize
	if window.Collapsed {
		// Title bar only
		var backup_border_size float = style.FrameBorderSize
		g.Style.FrameBorderSize = window.WindowBorderSize
		var title_bar_col ImU32
		if title_bar_is_highlight && !g.NavDisableHighlight {
			title_bar_col = GetColorU32FromID(ImGuiCol_TitleBgActive, 1)
		} else {
			title_bar_col = GetColorU32FromID(ImGuiCol_TitleBgCollapsed, 1)
		}
		RenderFrame(title_bar_rect.Min, title_bar_rect.Max, title_bar_col, true, window_rounding)
		g.Style.FrameBorderSize = backup_border_size
	} else {
		// Window background
		if 0 == (flags & ImGuiWindowFlags_NoBackground) {
			var bg_col ImU32 = GetColorU32FromID(GetWindowBgColorIdxFromFlags(flags), 1)

			var override_alpha bool = false
			var alpha float = 1.0
			if g.NextWindowData.Flags&ImGuiNextWindowDataFlags_HasBgAlpha != 0 {
				alpha = g.NextWindowData.BgAlphaVal
				override_alpha = true
			}
			if override_alpha {
				bg_col = uint32(uint64(uint64(bg_col)&^uint64(IM_COL32_A_MASK)) | uint64(IM_F32_TO_INT8_SAT(alpha)<<IM_COL32_A_SHIFT))
			}

			f := ImDrawFlags_RoundCornersBottom
			if flags&ImGuiWindowFlags_NoTitleBar != 0 {
				f = 0
			}

			window.DrawList.AddRectFilled(window.Pos.Add(ImVec2{0, window.TitleBarHeight()}), window.Pos.Add(window.Size), bg_col, window_rounding, f)
		}

		// Title bar
		if 0 == (flags & ImGuiWindowFlags_NoTitleBar) {
			var title_bar_col ImU32 = GetColorU32FromID(ImGuiCol_TitleBg, 1)
			if title_bar_is_highlight {
				title_bar_col = GetColorU32FromID(ImGuiCol_TitleBgActive, 1)
			}
			window.DrawList.AddRectFilled(title_bar_rect.Min, title_bar_rect.Max, title_bar_col, window_rounding, ImDrawFlags_RoundCornersTop)
		}

		// Menu bar
		if flags&ImGuiWindowFlags_MenuBar != 0 {
			var menu_bar_rect ImRect = window.MenuBarRect()
			menu_bar_rect.ClipWith(window.Rect()) // Soft clipping, in particular child window don't have minimum size covering the menu bar so this is useful for them.

			round := window_rounding
			if flags&ImGuiWindowFlags_NoTitleBar != 0 {
				round = 0
			}

			window.DrawList.AddRectFilled(menu_bar_rect.Min.Add(ImVec2{window_border_size, 0}), menu_bar_rect.Max.Sub(ImVec2{window_border_size, 0}), GetColorU32FromID(ImGuiCol_MenuBarBg, 1), round, ImDrawFlags_RoundCornersTop)
			if style.FrameBorderSize > 0.0 && menu_bar_rect.Max.y < window.Pos.y+window.Size.y {
				from, to := menu_bar_rect.GetBL(), menu_bar_rect.GetBR()
				window.DrawList.AddLine(&from, &to, GetColorU32FromID(ImGuiCol_Border, 1), style.FrameBorderSize)
			}
		}

		// Scrollbars
		if window.ScrollbarX {
			Scrollbar(ImGuiAxis_X)
		}
		if window.ScrollbarY {
			Scrollbar(ImGuiAxis_Y)
		}

		// Render resize grips (after their input handling so we don't have a frame of latency)
		if 0 == (flags & ImGuiWindowFlags_NoResize) {
			for resize_grip_n := int(0); resize_grip_n < resize_grip_count; resize_grip_n++ {
				var grip *ImGuiResizeGripDef = &resize_grip_def[resize_grip_n]

				lt := window.Pos.Add(window.Size)
				var corner ImVec2 = ImLerpVec2WithVec2(&window.Pos, &lt, grip.CornerPosN)

				var to1, to2 ImVec2 = ImVec2{resize_grip_draw_size, window_border_size}, ImVec2{window_border_size, resize_grip_draw_size}
				if resize_grip_n&1 != 0 {
					to1 = ImVec2{window_border_size, resize_grip_draw_size}
					to2 = ImVec2{resize_grip_draw_size, window_border_size}
				}

				window.DrawList.PathLineTo(corner.Add(grip.InnerDir.Mul(to1)))
				window.DrawList.PathLineTo(corner.Add(grip.InnerDir.Mul(to2)))
				window.DrawList.PathArcToFast(ImVec2{corner.x + grip.InnerDir.x*(window_rounding+window_border_size), corner.y + grip.InnerDir.y*(window_rounding+window_border_size)}, window_rounding, grip.AngleMin12, grip.AngleMax12)
				window.DrawList.PathFillConvex(resize_grip_col[resize_grip_n])
			}
		}

		// Borders
		RenderWindowOuterBorders(window)
	}
}

func RenderWindowOuterBorders(window *ImGuiWindow) {
	var g = GImGui
	var rounding float = window.WindowRounding
	var border_size float = window.WindowBorderSize
	if border_size > 0.0 && 0 == (window.Flags&ImGuiWindowFlags_NoBackground) {
		window.DrawList.AddRect(window.Pos, window.Pos.Add(window.Size), GetColorU32FromID(ImGuiCol_Border, 1), rounding, 0, border_size)
	}

	var border_held int = int(window.ResizeBorderHeld)
	if border_held != -1 {
		var def ImGuiResizeBorderDef = resize_border_def[border_held]
		var border_r ImRect = GetResizeBorderRect(window, border_held, rounding, 0.0)
		window.DrawList.PathArcTo(
			ImLerpVec2WithVec2(&border_r.Min, &border_r.Max, def.SegmentN1).
				Add(ImVec2{0.5, 0.5}).Add(def.InnerDir.Scale(rounding)), rounding, def.OuterAngle-IM_PI*0.25, def.OuterAngle, 0)
		window.DrawList.PathArcTo(
			ImLerpVec2WithVec2(&border_r.Min, &border_r.Max, def.SegmentN2).
				Add(ImVec2{0.5, 0.5}).Add(def.InnerDir.Scale(rounding)), rounding, def.OuterAngle, def.OuterAngle+IM_PI*0.25, 0)
		window.DrawList.PathStroke(GetColorU32FromID(ImGuiCol_SeparatorActive, 1), 0, ImMax(2.0, border_size)) // Thicker than usual
	}
	if g.Style.FrameBorderSize > 0 && 0 == (window.Flags&ImGuiWindowFlags_NoTitleBar) {
		var y float = window.Pos.y + window.TitleBarHeight() - 1
		window.DrawList.AddLine(&ImVec2{window.Pos.x + border_size, y}, &ImVec2{window.Pos.x + window.Size.x - border_size, y}, GetColorU32FromID(ImGuiCol_Border, 1), g.Style.FrameBorderSize)
	}
}

// Render title text, collapse button, close button
func RenderWindowTitleBarContents(window *ImGuiWindow, title_bar_rect *ImRect, name string, p_open *bool) {
	var g = GImGui
	var style = g.Style
	var flags = window.Flags

	var has_close_button bool = (p_open != nil)
	var has_collapse_button bool = 0 == (flags&ImGuiWindowFlags_NoCollapse) && (style.WindowMenuButtonPosition != ImGuiDir_None)

	// Close & Collapse button are on the Menu NavLayer and don't default focus (unless there's nothing else on that layer)
	var item_flags_backup ImGuiItemFlags = g.CurrentItemFlags
	g.CurrentItemFlags |= ImGuiItemFlags_NoNavDefaultFocus
	window.DC.NavLayerCurrent = ImGuiNavLayer_Menu

	// Layout buttons
	// FIXME: Would be nice to generalize the subtleties expressed here into reusable code.
	var pad_l float = style.FramePadding.x
	var pad_r float = style.FramePadding.x
	var button_sz float = g.FontSize
	var close_button_pos ImVec2
	var collapse_button_pos ImVec2
	if has_close_button {
		pad_r += button_sz
		close_button_pos = ImVec2{title_bar_rect.Max.x - pad_r - style.FramePadding.x, title_bar_rect.Min.y}
	}
	if has_collapse_button && style.WindowMenuButtonPosition == ImGuiDir_Right {
		pad_r += button_sz
		collapse_button_pos = ImVec2{title_bar_rect.Max.x - pad_r - style.FramePadding.x, title_bar_rect.Min.y}
	}
	if has_collapse_button && style.WindowMenuButtonPosition == ImGuiDir_Left {
		collapse_button_pos = ImVec2{title_bar_rect.Min.x + pad_l - style.FramePadding.x, title_bar_rect.Min.y}
		pad_l += button_sz
	}

	// Collapse button (submitting first so it gets priority when choosing a navigation init fallback)
	if has_collapse_button {
		if CollapseButton(window.GetIDs("#COLLAPSE", ""), &collapse_button_pos) {
			window.WantCollapseToggle = true // Defer actual collapsing to next frame as we are too far in the Begin() function
		}
	}

	// Close button
	if has_close_button {
		if CloseButton(window.GetIDs("#CLOSE", ""), &close_button_pos) {
			*p_open = false
		}
	}

	window.DC.NavLayerCurrent = ImGuiNavLayer_Main
	g.CurrentItemFlags = item_flags_backup

	// Title bar text (with: horizontal alignment, avoiding collapse/close button, optional "unsaved document" marker)
	// FIXME: Refactor text alignment facilities along with RenderText helpers, this is WAY too much messy code..
	var marker_size_x float
	if flags&ImGuiWindowFlags_UnsavedDocument != 0 {
		marker_size_x = button_sz * 0.80
	}
	var text_size ImVec2 = CalcTextSize(name, true, -1).Add(ImVec2{marker_size_x, 0.0})

	// As a nice touch we try to ensure that centered title text doesn't get affected by visibility of Close/Collapse button,
	// while uncentered title text will still reach edges correctly.
	if pad_l > style.FramePadding.x {
		pad_l += g.Style.ItemInnerSpacing.x
	}
	if pad_r > style.FramePadding.x {
		pad_r += g.Style.ItemInnerSpacing.x
	}
	if style.WindowTitleAlign.x > 0.0 && style.WindowTitleAlign.x < 1.0 {
		var centerness float = ImSaturate(1.0 - ImFabs(style.WindowTitleAlign.x-0.5)*2.0) // 0.0f on either edges, 1.0f on center
		var pad_extend float = ImMin(ImMax(pad_l, pad_r), title_bar_rect.GetWidth()-pad_l-pad_r-text_size.x)
		pad_l = ImMax(pad_l, pad_extend*centerness)
		pad_r = ImMax(pad_r, pad_extend*centerness)
	}

	var layout_r = ImRect{ImVec2{title_bar_rect.Min.x + pad_l, title_bar_rect.Min.y}, ImVec2{title_bar_rect.Max.x - pad_r, title_bar_rect.Max.y}}
	var clip_r = ImRect{ImVec2{layout_r.Min.x, layout_r.Min.y}, ImVec2{ImMin(layout_r.Max.x+g.Style.ItemInnerSpacing.x, title_bar_rect.Max.x), layout_r.Max.y}}
	if flags&ImGuiWindowFlags_UnsavedDocument != 0 {
		var marker_pos ImVec2
		marker_pos.x = ImClamp(layout_r.Min.x+(layout_r.GetWidth()-text_size.x)*style.WindowTitleAlign.x+text_size.x, layout_r.Min.x, layout_r.Max.x)
		marker_pos.y = (layout_r.Min.y + layout_r.Max.y) * 0.5
		if marker_pos.x > layout_r.Min.x {
			RenderBullet(window.DrawList, marker_pos, GetColorU32FromID(ImGuiCol_Text, 1))
			clip_r.Max.x = ImMin(clip_r.Max.x, marker_pos.x-float((int)(marker_size_x*0.5)))
		}
	}
	//if (g.IO.KeyShift) window.DrawList.AddRect(layout_r.Min, layout_r.Max, IM_COL32(255, 128, 0, 255)); // [DEBUG]
	//if (g.IO.KeyCtrl) window.DrawList.AddRect(clip_r.Min, clip_r.Max, IM_COL32(255, 128, 0, 255)); // [DEBUG]
	RenderTextClipped(&layout_r.Min, &layout_r.Max, name, &text_size, &style.WindowTitleAlign, &clip_r)
}
