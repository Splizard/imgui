package imgui

func Image(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, tint_col ImVec4, border_col ImVec4) {
	var window = GetCurrentWindow()
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

func ImageButton(user_texture_id ImTextureID, size ImVec2, uv0 ImVec2, uv1 ImVec2, frame_padding int /*/*= /*/, bg_col ImVec4, tint_col ImVec4) bool {
	panic("not implemented")
} // <0 frame_padding uses default frame padding settings. 0 for no padding

// Image primitives
// - Read FAQ to understand what ImTextureID is.
// - "p_min" and "p_max" represent the upper-left and lower-right corners of the rectangle.
// - "uv_min" and "uv_max" represent the normalized texture coordinates to use for those corners. Using (0,0)->(1,1) texture coordinates will generally display the entire texture.
