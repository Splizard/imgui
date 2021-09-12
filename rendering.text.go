package imgui

import "fmt"

// read one character. return input UTF-8 bytes count
func ImTextCharFromUtf8(out_char *uint, text string) int {
	for i, c := range text {
		*out_char = uint(c)
		return int(i) + 1
	}
	*out_char = 0
	return 0
}

// Render helpers
// AVOID USING OUTSIDE OF IMGUI.CPP! NOT FOR PUBLIC CONSUMPTION. THOSE FUNCTIONS ARE A MESS. THEIR SIGNATURE AND BEHAVIOR WILL CHANGE, THEY NEED TO BE REFACTORED INTO SOMETHING DECENT.
// NB: All position are in absolute pixels coordinates (we are never using window coordinates internally)
func RenderText(pos ImVec2, t string, hide_text_after_hash bool /*= true*/) {
	panic("not implemented")
}

func RenderTextWrapped(pos ImVec2, text string, wrap_width float) {
	var g = GImGui
	var window = g.CurrentWindow

	if len(text) > 0 {
		window.DrawList.AddTextV(g.Font, g.FontSize, pos, GetColorU32FromID(ImGuiCol_Text, 1), text, wrap_width, nil)
		if g.LogEnabled {
			LogRenderedText(&pos, text)
		}
	}
}

// formatted text
func Text(format string, args ...interface{}) {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return
	}
	TextEx(fmt.Sprintf(format, args...), ImGuiTextFlags_NoWidthForLargeClippedText)
}

func TextEx(text string, flags ImGuiTextFlags) {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var g = GImGui

	var text_pos = ImVec2{window.DC.CursorPos.x, window.DC.CursorPos.y + window.DC.CurrLineTextBaseOffset}
	var wrap_pos_x float = window.DC.TextWrapPos
	var wrap_enabled bool = (wrap_pos_x >= 0.0)
	if len(text) > 2000 && !wrap_enabled {
		// Long text!
		// Perform manual coarse clipping to optimize for long multi-line text
		// - From this point we will only compute the width of lines that are visible. Optimization only available when word-wrapping is disabled.
		// - We also don't vertically center the text within the line full height, which is unlikely to matter because we are likely the biggest and only item on the line.
		// - We use memchr(), pay attention that well optimized versions of those str/mem functions are much faster than a casually written loop.
		var line int = 0
		var line_height float = GetTextLineHeight()
		var text_size ImVec2

		// Lines to skip (can't skip when logging text)
		var pos ImVec2 = text_pos
		if !g.LogEnabled {
			var lines_skippable int = (int)((window.ClipRect.Min.y - text_pos.y) / line_height)
			if lines_skippable > 0 {
				var lines_skipped int = 0
				for line < int(len(text)) && lines_skipped < lines_skippable {
					var line_end int
					for line_end = line; line_end < int(len(text)) && text[line_end] != '\n'; line_end++ {
					}
					if (flags & ImGuiTextFlags_NoWidthForLargeClippedText) == 0 {
						text_size.x = ImMax(text_size.x, CalcTextSize(text[line:line_end], false, 0).x)
					}
					line = line_end + 1
					lines_skipped++
				}
				pos.y += float(lines_skipped) * line_height
			}
		}

		// Lines to render
		if line < int(len(text)) {
			var line_rect = ImRect{pos, pos.Add(ImVec2{FLT_MAX, line_height})}
			for line < int(len(text)) {
				if IsClippedEx(&line_rect, 0, false) {
					break
				}

				var line_end int
				for line_end = line; line_end < int(len(text)) && text[line_end] != '\n'; line_end++ {
				}
				text_size.x = ImMax(text_size.x, CalcTextSize(text[line:], false, 0).x)
				RenderText(pos, text[line:line_end], false)
				line = line_end + 1
				line_rect.Min.y += line_height
				line_rect.Max.y += line_height
				pos.y += line_height
			}

			// Count remaining lines
			var lines_skipped int = 0
			for line < int(len(text)) {
				var line_end int
				for line_end = line; line_end < int(len(text)) && text[line_end] != '\n'; line_end++ {
				}
				if (flags & ImGuiTextFlags_NoWidthForLargeClippedText) == 0 {
					text_size.x = ImMax(text_size.x, CalcTextSize(text[line:line_end], false, 0).x)
				}
				line = line_end + 1
				lines_skipped++
			}
			pos.y += float(lines_skipped) * line_height
		}
		text_size.y = (pos.Sub(text_pos)).y

		var bb = ImRect{text_pos, text_pos.Add(text_size)}
		ItemSizeVec(&text_size, 0.0)
		ItemAdd(&bb, 0, nil, 0)
	} else {
		var wrap_width float
		if wrap_enabled {
			wrap_width = CalcWrapWidthForPos(&window.DC.CursorPos, wrap_pos_x)
		}
		var text_size ImVec2 = CalcTextSize(text, false, wrap_width)

		var bb = ImRect{text_pos, text_pos.Add(text_size)}
		ItemSizeVec(&text_size, 0.0)
		if !ItemAdd(&bb, 0, nil, 0) {
			return
		}

		// Render (we don't hide text after ## in this end-user function)
		RenderTextWrapped(bb.Min, text, wrap_width)
	}
}

func (this *ImDrawList) AddText(pos ImVec2, col ImU32, text string) {
	this.AddTextV(nil, 0.0, pos, col, text, 0, nil)
}
func (this *ImDrawList) AddTextV(font *ImFont, font_size float, pos ImVec2, col ImU32, text string, wrap_width float, cpu_fine_clip_rect *ImVec4) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}
	if len(text) == 0 {
		return
	}

	// Pull default font/size from the shared ImDrawListSharedData instance
	if font == nil {
		font = this._Data.Font
	}
	if font_size == 0.0 {
		font_size = this._Data.FontSize
	}

	IM_ASSERT(font.ContainerAtlas.TexID == this._CmdHeader.TextureId) // Use high-level ImGui::PushFont() or low-level ImDrawList::PushTextureId() to change font.

	var clip_rect ImVec4 = this._CmdHeader.ClipRect
	if cpu_fine_clip_rect != nil {
		clip_rect.x = ImMax(clip_rect.x, cpu_fine_clip_rect.x)
		clip_rect.y = ImMax(clip_rect.y, cpu_fine_clip_rect.y)
		clip_rect.z = ImMin(clip_rect.z, cpu_fine_clip_rect.z)
		clip_rect.w = ImMin(clip_rect.w, cpu_fine_clip_rect.w)
	}
	font.RenderText(this, font_size, pos, col, &clip_rect, text, wrap_width, cpu_fine_clip_rect != nil)
}

// Find the optional ## from which we stop displaying text.
func FindRenderedTextEnd(t string) string {
	var i int = 0
	for i < int(len(t)) {
		if t[i] == '#' {
			if i+1 < int(len(t)) && t[i+1] == '#' {
				return t[:i]
			}
		}
		i++
	}
	return t
}

// Text Utilities
func CalcTextSize(text string, hide_text_after_double_hash bool /*= e*/, wrap_width float /*= -1.0*/) ImVec2 {
	var g = GImGui

	var text_display_end string
	if hide_text_after_double_hash {
		text = FindRenderedTextEnd(text) // Hide anything after a '##' string
	}

	var font = g.Font
	var font_size float = g.FontSize
	if text == text_display_end {
		return ImVec2{0.0, font_size}
	}
	var text_size ImVec2 = font.CalcTextSizeA(font_size, FLT_MAX, wrap_width, text, nil)

	// Round
	// FIXME: This has been here since Dec 2015 (7b0bf230) but down the line we want this out.
	// FIXME: Investigate using ceilf or e.g.
	// - https://git.musl-libc.org/cgit/musl/tree/src/math/ceilf.c
	// - https://embarkstudios.github.io/rust-gpu/api/src/libm/math/ceilf.rs.html
	text_size.x = IM_FLOOR(text_size.x + 0.99999)

	return text_size
}

// Default clip_rect uses (pos_min,pos_max)
// Handle clipping on CPU immediately (vs typically let the GPU clip the triangles that are overlapping the clipping rectangle edges)
func RenderTextClippedEx(draw_list *ImDrawList, pos_min *ImVec2, pos_max *ImVec2, text string, text_size_if_known *ImVec2, align *ImVec2, clip_rect *ImRect) {
	// Perform CPU side clipping for single clipped element to avoid using scissor state
	var pos ImVec2 = *pos_min
	var text_size ImVec2
	if text_size_if_known != nil {
		text_size = *text_size_if_known
	} else {
		text_size = CalcTextSize(text, false, 0.0)
	}

	var clip_min *ImVec2 = pos_min
	var clip_max *ImVec2 = pos_max
	if clip_rect != nil {
		clip_min = &clip_rect.Min
		clip_max = &clip_rect.Max
	}
	var need_clipping bool = (pos.x+text_size.x >= clip_max.x) || (pos.y+text_size.y >= clip_max.y)
	if clip_rect != nil { // If we had no explicit clipping rectangle then pos==clip_min
		need_clipping = need_clipping || (pos.x < clip_min.x) || (pos.y < clip_min.y)
	}

	// Align whole block. We should defer that to the better rendering function when we'll have support for individual line alignment.
	if align.x > 0.0 {
		pos.x = ImMax(pos.x, pos.x+(pos_max.x-pos.x-text_size.x)*align.x)
	}
	if align.y > 0.0 {
		pos.y = ImMax(pos.y, pos.y+(pos_max.y-pos.y-text_size.y)*align.y)
	}

	// Render
	if need_clipping {
		var fine_clip_rect = ImVec4{clip_min.x, clip_min.y, clip_max.x, clip_max.y}
		draw_list.AddTextV(nil, 0.0, pos, GetColorU32FromID(ImGuiCol_Text, 1), text, 0.0, &fine_clip_rect)
	} else {
		draw_list.AddTextV(nil, 0.0, pos, GetColorU32FromID(ImGuiCol_Text, 1), text, 0.0, nil)
	}
}

func RenderTextClipped(pos_min *ImVec2, pos_max *ImVec2, text string, text_size_if_known *ImVec2, align *ImVec2, clip_rect *ImRect) {
	// Hide anything after a '##' string
	text = FindRenderedTextEnd(text)
	if len(text) == 0 {
		return
	}

	var g = GImGui
	var window *ImGuiWindow = g.CurrentWindow
	RenderTextClippedEx(window.DrawList, pos_min, pos_max, text, text_size_if_known, align, clip_rect)
	if g.LogEnabled {
		LogRenderedText(pos_min, text)
	}
}

func (this *ImFont) RenderText(draw_list *ImDrawList, size float, pos ImVec2, col ImU32, clip_rect *ImVec4, text string, wrap_width float, cpu_fine_clip bool) {

	// Align to be pixel perfect
	pos.x = IM_FLOOR(pos.x)
	pos.y = IM_FLOOR(pos.y)
	var x float = pos.x
	var y float = pos.y
	if y > clip_rect.w {
		return
	}

	var scale float = size / this.FontSize
	var line_height float = this.FontSize * scale
	var word_wrap_enabled bool = (wrap_width > 0.0)
	var word_wrap_eol int = -1

	// Fast-forward to first visible line
	var i int = 0
	if y+line_height < clip_rect.y && !word_wrap_enabled {
		for y+line_height < clip_rect.y && i < int(len(text)) {
			for ; i < int(len(text)) && text[i] != '\n'; i++ {
			}
			i++
			y += line_height
		}
	}

	// For large text, scan for the last visible line in order to avoid over-reserving in the call to PrimReserve()
	// Note that very large horizontal line will still be affected by the issue (e.g. a one megabyte string buffer without a newline will likely crash atm)
	if int(len(text))-i > 10000 && !word_wrap_enabled {
		var i_end = i
		var y_end float = y
		for y_end < clip_rect.w && i_end < int(len(text)) {
			for ; i_end < int(len(text)) && text[i_end] != '\n'; i_end++ {
			}
			i_end++
			y_end += line_height
		}
	}
	if i == int(len(text)) {
		return
	}

	// Reserve vertices for remaining worse case (over-reserving is useful and easily amortized)
	var vtx_count_max int = (int)(int(len(text))-i) * 4
	var idx_count_max int = (int)(int(len(text))-i) * 6
	var idx_expected_size int = int(len(draw_list.IdxBuffer)) + idx_count_max

	//TODO/FIXME these can be negative?
	if idx_count_max < 0 || vtx_count_max < 0 {
		return
	}

	draw_list.PrimReserve(idx_count_max, vtx_count_max)

	var vtx_write []ImDrawVert = draw_list._VtxWritePtr
	var idx_write []ImDrawIdx = draw_list._IdxWritePtr
	var vtx_current_idx uint = draw_list._VtxCurrentIdx

	var col_untinted ImU32 = ImU32(uint64(col) | ^uint64(IM_COL32_A_MASK))

	for i < int(len(text)) {

		if word_wrap_enabled {
			// Calculate how far we can render. Requires two passes on the string data but keeps the code simple and not intrusive for what's essentially an uncommon feature.
			if word_wrap_eol == -1 {
				word_wrap_eol = i + this.CalcWordWrapPositionA(scale, text[i:], wrap_width-(x-pos.x))
				if word_wrap_eol == i { // Wrap_width is too small to fit anything. Force displaying 1 character to minimize the height discontinuity.
					word_wrap_eol++ // +1 may not be a character start point in UTF-8 but it's ok because we use s >= word_wrap_eol below
				}
			}

			if i >= word_wrap_eol {
				x = pos.x
				y += line_height
				word_wrap_eol = -1

				// Wrapping skips upcoming blanks
				for i < int(len(text)) {
					var c = text[i]
					if ImCharIsBlankA(c) {
						i++
					} else if c == '\n' {
						i++
						break
					} else {
						break
					}
				}
				continue
			}
		}

		// Decode and advance source
		var c uint = (uint)(text[i])
		if c < 0x80 {
			i += 1
		} else {
			i += ImTextCharFromUtf8(&c, text[i:])
			if c == 0 { // Malformed UTF-8?
				break
			}
		}

		if c < 32 {
			if c == '\n' {
				x = pos.x
				y += line_height
				if y > clip_rect.w {
					break // break out of main loop
				}
				continue
			}
			if c == '\r' {
				continue
			}
		}

		var glyph *ImFontGlyph = this.FindGlyph((ImWchar)(c))
		if glyph == nil {
			continue
		}

		var char_width float = glyph.AdvanceX * scale
		if glyph.Visible != 0 {
			// We don't do a second finer clipping test on the Y axis as we've already skipped anything before clip_rect.y and exit once we pass clip_rect.w
			var x1 float = x + glyph.X0*scale
			var x2 float = x + glyph.X1*scale
			var y1 float = y + glyph.Y0*scale
			var y2 float = y + glyph.Y1*scale
			if x1 <= clip_rect.z && x2 >= clip_rect.x {
				// Render a character
				var u1 float = glyph.U0
				var v1 float = glyph.V0
				var u2 float = glyph.U1
				var v2 float = glyph.V1

				// CPU side clipping used to fit text in their frame when the frame is too small. Only does clipping for axis aligned quads.
				if cpu_fine_clip {
					if x1 < clip_rect.x {
						u1 = u1 + (1.0-(x2-clip_rect.x)/(x2-x1))*(u2-u1)
						x1 = clip_rect.x
					}
					if y1 < clip_rect.y {
						v1 = v1 + (1.0-(y2-clip_rect.y)/(y2-y1))*(v2-v1)
						y1 = clip_rect.y
					}
					if x2 > clip_rect.z {
						u2 = u1 + ((clip_rect.z-x1)/(x2-x1))*(u2-u1)
						x2 = clip_rect.z
					}
					if y2 > clip_rect.w {
						v2 = v1 + ((clip_rect.w-y1)/(y2-y1))*(v2-v1)
						y2 = clip_rect.w
					}
					if y1 >= y2 {
						x += char_width
						continue
					}
				}

				// Support for untinted glyphs
				var glyph_col ImU32 = col
				if glyph.Colored != 0 {
					glyph_col = col_untinted
				}
				// We are NOT calling PrimRectUV() here because non-inlined causes too much overhead in a debug builds. Inlined here:
				{
					idx_write[0] = (ImDrawIdx)(vtx_current_idx)
					idx_write[1] = (ImDrawIdx)(vtx_current_idx + 1)
					idx_write[2] = (ImDrawIdx)(vtx_current_idx + 2)
					idx_write[3] = (ImDrawIdx)(vtx_current_idx)
					idx_write[4] = (ImDrawIdx)(vtx_current_idx + 2)
					idx_write[5] = (ImDrawIdx)(vtx_current_idx + 3)
					vtx_write[0].pos.x = x1
					vtx_write[0].pos.y = y1
					vtx_write[0].col = glyph_col
					vtx_write[0].uv.x = u1
					vtx_write[0].uv.y = v1
					vtx_write[1].pos.x = x2
					vtx_write[1].pos.y = y1
					vtx_write[1].col = glyph_col
					vtx_write[1].uv.x = u2
					vtx_write[1].uv.y = v1
					vtx_write[2].pos.x = x2
					vtx_write[2].pos.y = y2
					vtx_write[2].col = glyph_col
					vtx_write[2].uv.x = u2
					vtx_write[2].uv.y = v2
					vtx_write[3].pos.x = x1
					vtx_write[3].pos.y = y2
					vtx_write[3].col = glyph_col
					vtx_write[3].uv.x = u1
					vtx_write[3].uv.y = v2
					vtx_write = vtx_write[4:]
					vtx_current_idx += 4
					idx_write = idx_write[6:]
				}
			}
		}
		x += char_width
	}

	// Give back unused vertices (clipped ones, blanks) ~ this is essentially a PrimUnreserve() action.
	draw_list.VtxBuffer = draw_list.VtxBuffer[:len(draw_list.VtxBuffer)-len(vtx_write)]
	draw_list.IdxBuffer = draw_list.IdxBuffer[:len(draw_list.IdxBuffer)-len(idx_write)]
	draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1].ElemCount -= uint(idx_expected_size - int(len(draw_list.IdxBuffer)))
	draw_list._VtxWritePtr = vtx_write
	draw_list._IdxWritePtr = idx_write
	draw_list._VtxCurrentIdx = vtx_current_idx
}
