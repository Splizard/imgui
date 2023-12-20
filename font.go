package imgui

// Style read access
// - Use the style editor (ShowStyleEditor() function) to interactively see what the colors are)

var FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA = [ImGuiMouseCursor_COUNT][3]ImVec2{
	// Pos ........ Size ......... Offset ......
	{ImVec2{0, 3}, ImVec2{12, 19}, ImVec2{0, 0}},    // ImGuiMouseCursor_Arrow
	{ImVec2{13, 0}, ImVec2{7, 16}, ImVec2{1, 8}},    // ImGuiMouseCursor_TextInput
	{ImVec2{31, 0}, ImVec2{23, 23}, ImVec2{11, 11}}, // ImGuiMouseCursor_ResizeAll
	{ImVec2{21, 0}, ImVec2{9, 23}, ImVec2{4, 11}},   // ImGuiMouseCursor_ResizeNS
	{ImVec2{55, 18}, ImVec2{23, 9}, ImVec2{11, 4}},  // ImGuiMouseCursor_ResizeEW
	{ImVec2{73, 0}, ImVec2{17, 17}, ImVec2{8, 8}},   // ImGuiMouseCursor_ResizeNESW
	{ImVec2{55, 0}, ImVec2{17, 17}, ImVec2{8, 8}},   // ImGuiMouseCursor_ResizeNWSE
	{ImVec2{91, 0}, ImVec2{17, 22}, ImVec2{5, 0}},   // ImGuiMouseCursor_Hand
}

// GetFont get current font
func GetFont() *ImFont { return GImGui.Font }

// GetFontSize get current font size (= height in pixels) of current font with current scale applied
func GetFontSize() float { return GImGui.FontSize }

func GetFontTexUvWhitePixel() ImVec2 { return GImGui.DrawListSharedData.TexUvWhitePixel } // get UV coordinate for a while pixel, useful to draw custom shapes via the ImDrawList API

// PushFont Parameters stacks (shared)
// use NULL as a shortcut to push default font
func PushFont(font *ImFont) {
	var g = GImGui
	if font == nil {
		font = GetDefaultFont()
	}
	SetCurrentFont(font)
	g.FontStack = append(g.FontStack, font)
	g.CurrentWindow.DrawList.PushTextureID(font.ContainerAtlas.TexID)
}
func PopFont() {
	var g = GImGui
	g.CurrentWindow.DrawList.PopTextureID()
	g.FontStack = g.FontStack[:len(g.FontStack)-1]
	if len(g.FontStack) == 0 {
		SetCurrentFont(GetDefaultFont())
	} else {
		SetCurrentFont(g.FontStack[len(g.FontStack)-1])
	}
}

// CalcWordWrapPositionA 'max_width' stops rendering after a certain width (could be turned into a 2d size). FLT_MAX to disable.
// 'wrap_width' enable automatic word-wrapping across multiple lines to fit into given width. 0.0f to disable.
func (f *ImFont) CalcWordWrapPositionA(scale float, text string, wrap_width float) int {
	// Simple word-wrapping for English, not full-featured. Please submit failing cases!
	// FIXME: Much possible improvements (don't cut things like "word !", "word!!!" but cut within "word,,,,", more sensible support for punctuations, support for Unicode punctuations, etc.)

	// For references, possible wrap point marked with ^
	//  "aaa bbb, ccc,ddd. eee   fff. ggg!"
	//      ^    ^    ^   ^   ^__    ^    ^

	// List of hardcoded separators: .,;!?'"

	// Skip extra blanks after a line returns (that includes not counting them in width computation)
	// e.g. "Hello    world" -. "Hello" "World"

	// Cut words that cannot possibly fit within one line.
	// e.g.: "The tropical fish" with ~5 characters worth of width -. "The tr" "opical" "fish"

	var line_width float = 0.0
	var word_width float = 0.0
	var blank_width float = 0.0
	wrap_width /= scale // We work with unscaled widths to avoid scaling every characters

	var word_end int = 0
	var prev_word_end int = -1
	var inside_word = true

	var i int
	for i = 0; i < int(len(text)); {
		var c = rune(text[i])

		var next_i int
		if c < 0x80 {
			next_i = i + 1
		} else {
			next_i = i + ImTextCharFromUtf8(&c, text[i:])
		}

		if c == 0 {
			break
		}

		if c < 32 {
			if c == '\n' {
				line_width = 0
				word_width = 0
				blank_width = 0.0
				inside_word = true
				i = next_i
				continue
			}
			if c == '\r' {
				i = next_i
				continue
			}
		}

		var char_width = f.FallbackAdvanceX
		if c < int(len(f.IndexAdvanceX)) {
			char_width = f.IndexAdvanceX[c]
		}

		if ImCharIsBlankW(c) {
			if inside_word {
				line_width += blank_width
				blank_width = 0.0
				word_end = i
			}
			blank_width += char_width
			inside_word = false
		} else {
			word_width += char_width
			if inside_word {
				word_end = next_i
			} else {
				prev_word_end = word_end
				line_width += word_width + blank_width
				word_width = 0
				blank_width = 0.0
			}

			// Allow wrapping after punctuation.
			inside_word = c != '.' && c != ',' && c != ';' && c != '!' && c != '?' && c != '"'
		}

		// We ignore blank width at the end of the line (they can be skipped)
		if line_width+word_width > wrap_width {
			// Words that cannot possibly fit within an entire line will be cut anywhere.
			if word_width < wrap_width {
				if prev_word_end != 0 {
					i = prev_word_end
				} else {
					i = word_end
				}
			}
			break
		}

		i = next_i
	}

	return i
}

func (f *ImFont) CalcTextSizeA(size, max_width, wrap_width float, text string, remaining *string) ImVec2 {
	var line_height = size
	var scale = size / f.FontSize

	var text_size = ImVec2{}
	var line_width float = 0.0

	var word_wrap_enabled = wrap_width > 0.0
	var word_wrap_eol int = -1

	var i int

	for i = 0; i < int(len(text)); {
		if word_wrap_enabled {
			// Calculate how far we can render. Requires two passes on the string data but keeps the code simple and not intrusive for what's essentially an uncommon feature.
			if word_wrap_eol == -1 {
				word_wrap_eol = i + f.CalcWordWrapPositionA(scale, text[i:], wrap_width-line_width)
				if word_wrap_eol == i { // Wrap_width is too small to fit anything. Force displaying 1 character to minimize the height discontinuity.
					word_wrap_eol++ // +1 may not be a character start point in UTF-8 but it's ok because we use s >= word_wrap_eol below
				}
			}

			if i >= word_wrap_eol {
				if text_size.x < line_width {
					text_size.x = line_width
				}
				text_size.y += line_height
				line_width = 0.0
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
		var prev_i = i
		var c = rune(text[i])
		if c < 0x80 {
			i += 1
		} else {
			i += ImTextCharFromUtf8(&c, text)
			if c == 0 { // Malformed UTF-8?
				break
			}
		}

		if c < 32 {
			if c == '\n' {
				text_size.x = ImMax(text_size.x, line_width)
				text_size.y += line_height
				line_width = 0.0
				continue
			}
			if c == '\r' {
				continue
			}
		}

		var char_width = f.FallbackAdvanceX
		if c < int(len(f.IndexAdvanceX)) {
			char_width = f.IndexAdvanceX[c]
		}
		char_width *= scale

		if line_width+char_width >= max_width {
			i = prev_i
			break
		}

		line_width += char_width
	}

	if text_size.x < line_width {
		text_size.x = line_width
	}

	if line_width > 0 || text_size.y == 0.0 {
		text_size.y += line_height
	}

	if remaining != nil {
		*remaining = text[i:]
	}

	return text_size
}

func SetCurrentFont(font *ImFont) {
	var g = GImGui
	IM_ASSERT(font != nil && font.IsLoaded()) // Font Atlas not created. Did you call io.Fonts.GetTexDataAsRGBA32 / GetTexDataAsAlpha8 ?
	IM_ASSERT(font.Scale > 0.0)

	g.Font = font
	g.FontBaseSize = ImMax(1.0, g.IO.FontGlobalScale*g.Font.FontSize*g.Font.Scale)

	if g.CurrentWindow != nil {
		g.FontSize = g.CurrentWindow.CalcFontSize()
	} else {
		g.FontSize = 0
	}

	var atlas = g.Font.ContainerAtlas
	g.DrawListSharedData.TexUvWhitePixel = atlas.TexUvWhitePixel
	g.DrawListSharedData.TexUvLines = atlas.TexUvLines[:]
	g.DrawListSharedData.Font = g.Font
	g.DrawListSharedData.FontSize = g.FontSize
}

func GetDefaultFont() *ImFont {
	var g = GImGui
	if g.IO.FontDefault != nil {
		return g.IO.FontDefault
	}
	return g.IO.Fonts.Fonts[0]
}

func (f *ImFont) ClearOutputData() {
	f.FontSize = 0
	f.FallbackAdvanceX = 0
	f.Glyphs = f.Glyphs[:0]
	f.IndexAdvanceX = f.IndexAdvanceX[:0]
	f.IndexLookup = f.IndexLookup[:0]
	f.FallbackGlyph = nil
	f.ContainerAtlas = nil
	f.DirtyLookupTables = true
	f.Ascent = 0
	f.Descent = 0
	f.MetricsTotalSurface = 0
}

func (f *ImFont) FindGlyph(c ImWchar) *ImFontGlyph {
	if (size_t)(c) >= (size_t)(len(f.IndexLookup)) {
		return f.FallbackGlyph
	}
	var i = f.IndexLookup[c]
	if i == (ImWchar)(-1) {
		return f.FallbackGlyph
	}
	return &f.Glyphs[i]
}

func (f *ImFont) SetGlyphVisible(c ImWchar, visible bool) {
	if glyph := f.FindGlyph(c); glyph != nil {
		if visible {
			glyph.Visible = 1
		} else {
			glyph.Visible = 0
		}
	}
}

// AddGlyph x0/y0/x1/y1 are offset from the character upper-left layout position, in pixels. Therefore x0/y0 are often fairly close to zero.
// Not to be mistaken with texture coordinates, which are held by u0/v0/u1/v1 in normalized format (0.0..1.0 on each texture axis).
// 'cfg' is not necessarily == 'this.ConfigData' because multiple source fonts+configs can be used to build one target font.
func (f *ImFont) AddGlyph(cfg *ImFontConfig, codepoint ImWchar, x0, y0, x1, y1, u0, v0, u1, v1, advance_x float) {
	if cfg != nil {
		// Clamp & recenter if needed
		var advance_x_original = advance_x
		advance_x = ImClamp(advance_x, cfg.GlyphMinAdvanceX, cfg.GlyphMaxAdvanceX)
		if advance_x != advance_x_original {
			var char_off_x float
			if cfg.PixelSnapH {
				char_off_x = ImFloor((advance_x - advance_x_original) * 0.5)
			} else {
				char_off_x = (advance_x - advance_x_original) * 0.5
			}
			x0 += char_off_x
			x1 += char_off_x
		}

		// Snap to pixel
		if cfg.PixelSnapH {
			advance_x = IM_ROUND(advance_x)
		}

		// Bake spacing
		advance_x += cfg.GlyphExtraSpacing.x
	}

	f.Glyphs = append(f.Glyphs, ImFontGlyph{})
	var glyph = &f.Glyphs[len(f.Glyphs)-1]
	glyph.Codepoint = (uint)(codepoint)
	glyph.Visible = uint(bool2int((x0 != x1) && (y0 != y1)))
	glyph.Colored = uint(bool2int(false))
	glyph.X0 = x0
	glyph.Y0 = y0
	glyph.X1 = x1
	glyph.Y1 = y1
	glyph.U0 = u0
	glyph.V0 = v0
	glyph.U1 = u1
	glyph.V1 = v1
	glyph.AdvanceX = advance_x

	// Compute rough surface usage metrics (+1 to account for average padding, +0.99 to round)
	// We use (U1-U0)*TexWidth instead of X1-X0 to account for oversampling.
	var pad = float(f.ContainerAtlas.TexGlyphPadding) + 0.99
	f.DirtyLookupTables = true
	f.MetricsTotalSurface += (int)((glyph.U1-glyph.U0)*float(f.ContainerAtlas.TexWidth)+pad) * (int)((glyph.V1-glyph.V0)*float(f.ContainerAtlas.TexHeight)+pad)
}

func FindFirstExistingGlyph(font *ImFont, candidate_chars []ImWchar, candidate_chars_count int) ImWchar {
	for n := int(0); n < candidate_chars_count; n++ {
		if font.FindGlyphNoFallback(candidate_chars[n]) != nil {
			return candidate_chars[n]
		}
	}
	return (ImWchar)(-1)
}

func (f *ImFont) FindGlyphNoFallback(c ImWchar) *ImFontGlyph {
	if size_t(c) >= (size_t)(len(f.IndexLookup)) {
		return nil
	}
	var i = f.IndexLookup[c]
	if i == (ImWchar)(-1) {
		return nil
	}
	return &f.Glyphs[i]
}

func (f *ImFont) GrowIndex(new_size int) {
	IM_ASSERT(len(f.IndexAdvanceX) == len(f.IndexLookup))
	if new_size <= int(len(f.IndexLookup)) {
		return
	}
	for int(len(f.IndexAdvanceX)) < new_size {
		f.IndexAdvanceX = append(f.IndexAdvanceX, -1)
	}
	for int(len(f.IndexLookup)) < new_size {
		f.IndexLookup = append(f.IndexLookup, (ImWchar)(-1))
	}
}

func (f *ImFont) BuildLookupTable() {
	var max_codepoint int = 0
	for i := range f.Glyphs {
		max_codepoint = ImMaxInt(max_codepoint, (int)(f.Glyphs[i].Codepoint))
	}

	// Build lookup table
	IM_ASSERT(len(f.Glyphs) < 0xFFFF) // -1 is reserved
	f.IndexAdvanceX = f.IndexAdvanceX[:0]
	f.IndexLookup = f.IndexLookup[:0]
	f.DirtyLookupTables = false
	f.Used4kPagesMap = [2]byte{}
	f.GrowIndex(max_codepoint + 1)
	for i := range f.Glyphs {
		var codepoint = (int)(f.Glyphs[i].Codepoint)
		f.IndexAdvanceX[codepoint] = f.Glyphs[i].AdvanceX
		f.IndexLookup[codepoint] = (ImWchar)(i)

		// Mark 4K page as used
		var page_n = codepoint / 4096
		f.Used4kPagesMap[page_n>>3] |= 1 << (page_n & 7)
	}

	// Create a glyph to handle TAB
	// FIXME: Needs proper TAB handling but it needs to be contextualized (or we could arbitrary say that each string starts at "column 0" ?)
	if f.FindGlyph(' ') != nil {
		if f.Glyphs[len(f.Glyphs)-1].Codepoint != '\t' { // So we can call f function multiple times (FIXME: Flaky)
			f.Glyphs = append(f.Glyphs, ImFontGlyph{})
		}
		var tab_glyph = &f.Glyphs[len(f.Glyphs)-1]
		*tab_glyph = *f.FindGlyph(' ')
		tab_glyph.Codepoint = '\t'
		tab_glyph.AdvanceX *= IM_TABSIZE
		f.IndexAdvanceX[(int)(tab_glyph.Codepoint)] = tab_glyph.AdvanceX
		f.IndexLookup[(int)(tab_glyph.Codepoint)] = (ImWchar)(len(f.Glyphs) - 1)
	}

	// Mark special glyphs as not visible (note that AddGlyph already mark as non-visible glyphs with zero-size polygons)
	f.SetGlyphVisible(' ', false)
	f.SetGlyphVisible('\t', false)

	// Ellipsis character is required for rendering elided text. We prefer using U+2026 (horizontal ellipsis).
	// However some old fonts may contain ellipsis at U+0085. Here we auto-detect most suitable ellipsis character.
	// FIXME: Note that 0x2026 is rarely included in our font ranges. Because of f we are more likely to use three individual dots.
	var ellipsis_chars = []ImWchar{(ImWchar)(0x2026), (ImWchar)(0x0085)}
	var dots_chars = []ImWchar{'.', (ImWchar)(0xFF0E)}
	if f.EllipsisChar == (ImWchar)(-1) {
		f.EllipsisChar = FindFirstExistingGlyph(f, ellipsis_chars, int(len(ellipsis_chars)))
	}
	if f.DotChar == (ImWchar)(-1) {
		f.DotChar = FindFirstExistingGlyph(f, dots_chars, int(len(dots_chars)))
	}

	// Setup fallback character
	var fallback_chars = []ImWchar{IM_UNICODE_CODEPOINT_INVALID, '?', ' '}
	f.FallbackGlyph = f.FindGlyphNoFallback(f.FallbackChar)
	if f.FallbackGlyph == nil {
		f.FallbackChar = FindFirstExistingGlyph(f, fallback_chars, int(len(fallback_chars)))
		f.FallbackGlyph = f.FindGlyphNoFallback(f.FallbackChar)
		if f.FallbackGlyph == nil {
			f.FallbackGlyph = &f.Glyphs[len(f.Glyphs)-1]
			f.FallbackChar = (ImWchar)(f.FallbackGlyph.Codepoint)
		}
	}

	f.FallbackAdvanceX = f.FallbackGlyph.AdvanceX
	for i := int(0); i < max_codepoint+1; i++ {
		if f.IndexAdvanceX[i] < 0.0 {
			f.IndexAdvanceX[i] = f.FallbackAdvanceX
		}
	}
}
