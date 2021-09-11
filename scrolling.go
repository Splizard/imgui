package imgui

// Return scrollbar rectangle, must only be called for corresponding axis if window.ScrollbarX/Y is set.
func GetWindowScrollbarRect(window *ImGuiWindow, axis ImGuiAxis) ImRect {
	var outer_rect ImRect = window.Rect()
	var inner_rect ImRect = window.InnerRect
	var border_size float = window.WindowBorderSize
	var scrollbar_size float

	// (ScrollbarSizes.x = width of Y scrollbar; ScrollbarSizes.y = height of X scrollbar)
	//TODO/FIXME is this correct?
	switch axis {
	case ImGuiAxis_X:
		scrollbar_size = window.ScrollbarSizes.x
	case ImGuiAxis_Y:
		scrollbar_size = window.ScrollbarSizes.y
	}

	IM_ASSERT(scrollbar_size > 0.0)
	if axis == ImGuiAxis_X {
		return ImRect{ImVec2{inner_rect.Min.x, ImMax(outer_rect.Min.y, outer_rect.Max.y-border_size-scrollbar_size)}, ImVec2{inner_rect.Max.x, outer_rect.Max.y}}
	} else {
		return ImRect{ImVec2{ImMax(outer_rect.Min.x, outer_rect.Max.x-border_size-scrollbar_size), inner_rect.Min.y}, ImVec2{outer_rect.Max.x, inner_rect.Max.y}}
	}
}

func GetWindowScrollbarID(window *ImGuiWindow, axis ImGuiAxis) ImGuiID {
	switch axis {
	case ImGuiAxis_X:
		return window.GetIDNoKeepAlive("#SCROLLX")
	case ImGuiAxis_Y:
		return window.GetIDNoKeepAlive("#SCROLLY")
	}
	return 0
}

func Scrollbar(axis ImGuiAxis) {
	var g = GImGui
	var window = g.CurrentWindow

	var id ImGuiID = GetWindowScrollbarID(window, axis)
	KeepAliveID(id)

	// Calculate scrollbar bounding box
	var bb ImRect = GetWindowScrollbarRect(window, axis)
	var rounding_corners ImDrawFlags = ImDrawFlags_RoundCornersNone
	if axis == ImGuiAxis_X {
		rounding_corners |= ImDrawFlags_RoundCornersBottomLeft
		if !window.ScrollbarY {
			rounding_corners |= ImDrawFlags_RoundCornersBottomRight
		}
	} else {
		if (window.Flags&ImGuiWindowFlags_NoTitleBar != 0) && 0 == (window.Flags&ImGuiWindowFlags_MenuBar) {
			rounding_corners |= ImDrawFlags_RoundCornersTopRight
		}
		if !window.ScrollbarX {
			rounding_corners |= ImDrawFlags_RoundCornersBottomRight
		}
	}
	var size_avail float
	var size_contents float
	var amount float
	switch axis {
	case ImGuiAxis_X:
		size_avail = window.InnerRect.Max.x - window.InnerRect.Min.x
		amount = window.Scroll.x
	case ImGuiAxis_Y:
		size_contents = window.ContentSize.y + window.WindowPadding.y*2.0
		amount = window.Scroll.y
	}

	ScrollbarEx(&bb, id, axis, &amount, size_avail, size_contents, rounding_corners)
}

// Vertical/Horizontal scrollbar
// The entire piece of code below is rather confusing because:
// - We handle absolute seeking (when first clicking outside the grab) and relative manipulation (afterward or when clicking inside the grab)
// - We store values as normalized ratio and in a form that allows the window content to change while we are holding on a scrollbar
// - We handle both horizontal and vertical scrollbars, which makes the terminology not ideal.
// Still, the code should probably be made simpler..
func ScrollbarEx(bb_frame *ImRect, id ImGuiID, axis ImGuiAxis, p_scroll_v *float, size_avail_v float, size_contents_v float, flags ImDrawFlags) bool {
	var g = GImGui
	var window = g.CurrentWindow
	if window.SkipItems {
		return false
	}

	var bb_frame_width float = bb_frame.GetWidth()
	var bb_frame_height float = bb_frame.GetHeight()
	if bb_frame_width <= 0.0 || bb_frame_height <= 0.0 {
		return false
	}

	// When we are too small, start hiding and disabling the grab (this reduce visual noise on very small window and facilitate using the window resize grab)
	var alpha float = 1.0
	if (axis == ImGuiAxis_Y) && bb_frame_height < g.FontSize+g.Style.FramePadding.y*2.0 {
		alpha = ImSaturate((bb_frame_height - g.FontSize) / (g.Style.FramePadding.y * 2.0))
	}
	if alpha <= 0.0 {
		return false
	}

	var style *ImGuiStyle = &g.Style
	var allow_interaction bool = (alpha >= 1.0)

	var bb ImRect = *bb_frame
	bb.ExpandVec(ImVec2{-ImClamp(IM_FLOOR((bb_frame_width-2.0)*0.5), 0.0, 3.0), -ImClamp(IM_FLOOR((bb_frame_height-2.0)*0.5), 0.0, 3.0)})

	// V denote the main, longer axis of the scrollbar (= height for a vertical scrollbar)
	var scrollbar_size_v float

	if axis == ImGuiAxis_X {
		scrollbar_size_v = bb.GetWidth()
	} else {
		scrollbar_size_v = bb.GetHeight()
	}

	// Calculate the height of our grabbable box. It generally represent the amount visible (vs the total scrollable amount)
	// But we maintain a minimum size in pixel to allow for the user to still aim inside.
	IM_ASSERT(ImMax(size_contents_v, size_avail_v) > 0.0) // Adding this assert to check if the ImMax(XXX,1.0f) is still needed. PLEASE CONTACT ME if this triggers.
	var win_size_v float = ImMax(ImMax(size_contents_v, size_avail_v), 1.0)
	var grab_h_pixels float = ImClamp(scrollbar_size_v*(size_avail_v/win_size_v), style.GrabMinSize, scrollbar_size_v)
	var grab_h_norm float = grab_h_pixels / scrollbar_size_v

	// Handle input right away. None of the code of Begin() is relying on scrolling position before calling Scrollbar().
	var held bool = false
	var hovered bool = false
	ButtonBehavior(&bb, id, &hovered, &held, ImGuiButtonFlags_NoNavFocus)

	var scroll_max float = ImMax(1.0, size_contents_v-size_avail_v)
	var scroll_ratio float = ImSaturate(*p_scroll_v / scroll_max)
	var grab_v_norm float = scroll_ratio * (scrollbar_size_v - grab_h_pixels) / scrollbar_size_v // Grab position in normalized space
	if held && allow_interaction && grab_h_norm < 1.0 {
		var scrollbar_pos_v float
		var mouse_pos_v float
		switch axis {
		case ImGuiAxis_X:
			scrollbar_pos_v = bb.Min.x
			mouse_pos_v = g.IO.MousePos.x
		case ImGuiAxis_Y:
			scrollbar_pos_v = bb.Min.y
			mouse_pos_v = g.IO.MousePos.y
		}

		// Click position in scrollbar normalized space (0.0f.1.0f)
		var clicked_v_norm float = ImSaturate((mouse_pos_v - scrollbar_pos_v) / scrollbar_size_v)
		SetHoveredID(id)

		var seek_absolute bool = false
		if g.ActiveIdIsJustActivated {
			// On initial click calculate the distance between mouse and the center of the grab
			seek_absolute = (clicked_v_norm < grab_v_norm || clicked_v_norm > grab_v_norm+grab_h_norm)
			if seek_absolute {
				g.ScrollbarClickDeltaToGrabCenter = 0.0
			} else {
				g.ScrollbarClickDeltaToGrabCenter = clicked_v_norm - grab_v_norm - grab_h_norm*0.5
			}
		}

		// Apply scroll (p_scroll_v will generally point on one member of window.Scroll)
		// It is ok to modify Scroll here because we are being called in Begin() after the calculation of ContentSize and before setting up our starting position
		var scroll_v_norm float = ImSaturate((clicked_v_norm - g.ScrollbarClickDeltaToGrabCenter - grab_h_norm*0.5) / (1.0 - grab_h_norm))
		*p_scroll_v = IM_ROUND(scroll_v_norm * scroll_max) //(win_size_contents_v - win_size_v));

		// Update values for rendering
		scroll_ratio = ImSaturate(*p_scroll_v / scroll_max)
		grab_v_norm = scroll_ratio * (scrollbar_size_v - grab_h_pixels) / scrollbar_size_v

		// Update distance to grab now that we have seeked and saturated
		if seek_absolute {
			g.ScrollbarClickDeltaToGrabCenter = clicked_v_norm - grab_v_norm - grab_h_norm*0.5
		}
	}

	// Render
	var bg_col ImU32 = GetColorU32FromID(ImGuiCol_ScrollbarBg, 1)
	var grab_col ImU32 = GetColorU32FromID(ImGuiCol_ScrollbarGrab, alpha)

	if held {
		grab_col = GetColorU32FromID(ImGuiCol_ScrollbarGrabActive, 1)
	} else if hovered {
		grab_col = GetColorU32FromID(ImGuiCol_ScrollbarGrabHovered, 1)
	}

	window.DrawList.AddRectFilled(bb_frame.Min, bb_frame.Max, bg_col, window.WindowRounding, flags)
	var grab_rect ImRect
	if axis == ImGuiAxis_X {
		grab_rect = ImRect{ImVec2{ImLerp(bb.Min.x, bb.Max.x, grab_v_norm), bb.Min.y}, ImVec2{ImLerp(bb.Min.x, bb.Max.x, grab_v_norm) + grab_h_pixels, bb.Max.y}}
	} else {
		grab_rect = ImRect{ImVec2{bb.Min.x, ImLerp(bb.Min.y, bb.Max.y, grab_v_norm)}, ImVec2{bb.Max.x, ImLerp(bb.Min.y, bb.Max.y, grab_v_norm) + grab_h_pixels}}
	}

	window.DrawList.AddRectFilled(grab_rect.Min, grab_rect.Max, grab_col, style.ScrollbarRounding, 0)

	return held
}
