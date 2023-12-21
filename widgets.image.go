package imgui

func Image(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, tint_col ImVec4, border_col ImVec4) {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)}
	if border_col.w > 0.0 {
		bb.Max = bb.Max.Add(ImVec2{2, 2})
	}
	ItemSizeRect(&bb, 0)
	if !ItemAdd(&bb, 0, nil, 0) {
		return
	}

	if border_col.w > 0.0 {
		window.DrawList.AddRect(bb.Min, bb.Max, GetColorU32FromVec(border_col), 0.0, 0, 1)
		window.DrawList.AddImage(user_texture_id, bb.Min.Add(ImVec2{1, 1}), bb.Max.Sub(ImVec2{1, 1}), &uv0, &uv1, GetColorU32FromVec(tint_col))
	} else {
		window.DrawList.AddImage(user_texture_id, bb.Min, bb.Max, &uv0, &uv1, GetColorU32FromVec(tint_col))
	}
}

// ImageButton() is flawed as 'id' is always derived from 'texture_id' (see #2464 #1390)
// We provide this internal helper to write your own variant while we figure out how to redesign the public ImageButton() API.
func ImageButtonEx(id ImGuiID, texture_id ImTextureID, size *ImVec2, uv0 *ImVec2, uv1 *ImVec2, padding *ImVec2, bg_col *ImVec4, tint_col *ImVec4) bool {
	g := GImGui
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(*size).Add(padding.Scale(2))}
	ItemSizeRect(&bb, 0)
	if !ItemAdd(&bb, id, nil, 0) {
		return false
	}

	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, 0)

	// Render
	var c = ImGuiCol_Button
	if held && hovered {
		c = ImGuiCol_ButtonActive
	} else if hovered {
		c = ImGuiCol_ButtonHovered
	}
	var col = GetColorU32FromID(c, 1)
	RenderNavHighlight(&bb, id, 0)
	RenderFrame(bb.Min, bb.Max, col, true, ImClamp((float)(min(padding.x, padding.y)), 0.0, g.Style.FrameRounding))
	if bg_col.w > 0.0 {
		window.DrawList.AddRectFilled(bb.Min.Add(*padding), bb.Max.Sub(*padding), GetColorU32FromVec(*bg_col), 0, 0)
	}
	window.DrawList.AddImage(texture_id, bb.Min.Add(*padding), bb.Max.Sub(*padding), uv0, uv1, GetColorU32FromVec(*tint_col))

	return pressed
}

// frame_padding < 0: uses FramePadding from style (default)
// frame_padding = 0: no framing
// frame_padding > 0: set framing size
func ImageButton(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, frame_padding int /*/*= /*/, bg_col ImVec4, tint_col ImVec4) bool {
	g := GImGui
	window := g.CurrentWindow
	if window.SkipItems {
		return false
	}

	// Default to using texture ID as ID. User can still push string/integer prefixes.
	PushID(int(user_texture_id))
	var id = window.GetIDs("#image")
	PopID()

	var padding = g.Style.FramePadding
	if frame_padding >= 0 {
		padding = ImVec2{(float)(frame_padding), (float)(frame_padding)}
	}
	return ImageButtonEx(id, user_texture_id, &size, &uv0, &uv1, &padding, &bg_col, &tint_col)
}

// Image primitives
// - Read FAQ to understand what ImTextureID is.
// - "p_min" and "p_max" represent the upper-left and lower-right corners of the rectangle.
// - "uv_min" and "uv_max" represent the normalized texture coordinates to use for those corners. Using (0,0)->(1,1) texture coordinates will generally display the entire texture.
