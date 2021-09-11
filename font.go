package imgui

import (
	"math"
)

// 'max_width' stops rendering after a certain width (could be turned into a 2d size). FLT_MAX to disable.
// 'wrap_width' enable automatic word-wrapping across multiple lines to fit into given width. 0.0f to disable.
func (this *ImFont) CalcWordWrapPositionA(scale float, text string, wrap_width float) int {
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
	var inside_word bool = true

	var i int
	for i = 0; i < int(len(text)); {
		var c uint = uint(text[i])

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

		var char_width float = this.FallbackAdvanceX
		if (int)(c) < int(len(this.IndexAdvanceX)) {
			char_width = this.IndexAdvanceX[c]
		}

		if ImCharIsBlankW(uint(c)) {
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
			inside_word = (c != '.' && c != ',' && c != ';' && c != '!' && c != '?' && c != '"')
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

func (this *ImFont) CalcTextSizeA(size, max_width, wrap_width float, text string, remaining *string) ImVec2 {
	var line_height = size
	var scale = size / this.FontSize

	var text_size ImVec2 = ImVec2{}
	var line_width float = 0.0

	var word_wrap_enabled = (wrap_width > 0.0)
	var word_wrap_eol int = -1

	var i int

	for i = 0; i < int(len(text)); {
		if word_wrap_enabled {
			// Calculate how far we can render. Requires two passes on the string data but keeps the code simple and not intrusive for what's essentially an uncommon feature.
			if word_wrap_eol == -1 {
				word_wrap_eol = i + this.CalcWordWrapPositionA(scale, text[i:], wrap_width-line_width)
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
		var c uint = uint(text[i])
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

		var char_width float = this.FallbackAdvanceX
		if (int)(c) < int(len(this.IndexAdvanceX)) {
			char_width = this.IndexAdvanceX[c]
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

	var atlas *ImFontAtlas = g.Font.ContainerAtlas
	g.DrawListSharedData.TexUvWhitePixel = atlas.TexUvWhitePixel
	g.DrawListSharedData.TexUvLines = atlas.TexUvLines[:]
	g.DrawListSharedData.Font = g.Font
	g.DrawListSharedData.FontSize = g.FontSize
}

func GetDefaultFont() *ImFont {
	var g *ImGuiContext = GImGui
	if g.IO.FontDefault != nil {
		return g.IO.FontDefault
	}
	return g.IO.Fonts.Fonts[0]
}

func (this *ImFont) ClearOutputData() {
	this.FontSize = 0
	this.FallbackAdvanceX = 0
	this.Glyphs = this.Glyphs[:0]
	this.IndexAdvanceX = this.IndexAdvanceX[:0]
	this.IndexLookup = this.IndexLookup[:0]
	this.FallbackGlyph = nil
	this.ContainerAtlas = nil
	this.DirtyLookupTables = true
	this.Ascent = 0
	this.Descent = 0
	this.MetricsTotalSurface = 0
}

func (this *ImFont) FindGlyph(c ImWchar) *ImFontGlyph {
	if (size_t)(c) >= (size_t)(len(this.IndexLookup)) {
		return this.FallbackGlyph
	}
	var i ImWchar = this.IndexLookup[c]
	if i == math.MaxUint16 {
		return this.FallbackGlyph
	}
	return &this.Glyphs[i]
}

func (this *ImFont) SetGlyphVisible(c ImWchar, visible bool) {
	if glyph := this.FindGlyph((ImWchar)(c)); glyph != nil {
		if visible {
			glyph.Visible = 1
		} else {
			glyph.Visible = 0
		}
	}
}

// x0/y0/x1/y1 are offset from the character upper-left layout position, in pixels. Therefore x0/y0 are often fairly close to zero.
// Not to be mistaken with texture coordinates, which are held by u0/v0/u1/v1 in normalized format (0.0..1.0 on each texture axis).
// 'cfg' is not necessarily == 'this.ConfigData' because multiple source fonts+configs can be used to build one target font.
func (this *ImFont) AddGlyph(cfg *ImFontConfig, codepoint ImWchar, x0, y0, x1, y1, u0, v0, u1, v1, advance_x float) {
	if cfg != nil {
		// Clamp & recenter if needed
		var advance_x_original float = advance_x
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

	this.Glyphs = append(this.Glyphs, ImFontGlyph{})
	var glyph *ImFontGlyph = &this.Glyphs[len(this.Glyphs)-1]
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
	var pad float = float(this.ContainerAtlas.TexGlyphPadding) + 0.99
	this.DirtyLookupTables = true
	this.MetricsTotalSurface += (int)((glyph.U1-glyph.U0)*float(this.ContainerAtlas.TexWidth)+pad) * (int)((glyph.V1-glyph.V0)*float(this.ContainerAtlas.TexHeight)+pad)
}

func FindFirstExistingGlyph(font *ImFont, candidate_chars []ImWchar, candidate_chars_count int) ImWchar {
	for n := int(0); n < candidate_chars_count; n++ {
		if font.FindGlyphNoFallback(candidate_chars[n]) != nil {
			return candidate_chars[n]
		}
	}
	return MaxImWchar
}

func (this *ImFont) FindGlyphNoFallback(c ImWchar) *ImFontGlyph {
	if size_t(c) >= (size_t)(len(this.IndexLookup)) {
		return nil
	}
	var i ImWchar = this.IndexLookup[c]
	if i == MaxImWchar {
		return nil
	}
	return &this.Glyphs[i]
}

func (this *ImFont) GrowIndex(new_size int) {
	IM_ASSERT(len(this.IndexAdvanceX) == len(this.IndexLookup))
	if new_size <= int(len(this.IndexLookup)) {
		return
	}
	for int(len(this.IndexAdvanceX)) < new_size {
		this.IndexAdvanceX = append(this.IndexAdvanceX, -1)
	}
	for int(len(this.IndexLookup)) < new_size {
		this.IndexLookup = append(this.IndexLookup, MaxImWchar)
	}
}

func (this *ImFont) BuildLookupTable() {
	var max_codepoint int = 0
	for i := range this.Glyphs {
		max_codepoint = ImMaxInt(max_codepoint, (int)(this.Glyphs[i].Codepoint))
	}

	// Build lookup table
	IM_ASSERT(len(this.Glyphs) < 0xFFFF) // -1 is reserved
	this.IndexAdvanceX = this.IndexAdvanceX[:0]
	this.IndexLookup = this.IndexLookup[:0]
	this.DirtyLookupTables = false
	this.Used4kPagesMap = [2]byte{}
	this.GrowIndex(max_codepoint + 1)
	for i := range this.Glyphs {
		var codepoint int = (int)(this.Glyphs[i].Codepoint)
		this.IndexAdvanceX[codepoint] = this.Glyphs[i].AdvanceX
		this.IndexLookup[codepoint] = (ImWchar)(i)

		// Mark 4K page as used
		var page_n int = codepoint / 4096
		this.Used4kPagesMap[page_n>>3] |= 1 << (page_n & 7)
	}

	// Create a glyph to handle TAB
	// FIXME: Needs proper TAB handling but it needs to be contextualized (or we could arbitrary say that each string starts at "column 0" ?)
	if this.FindGlyph((ImWchar)(' ')) != nil {
		if this.Glyphs[len(this.Glyphs)-1].Codepoint != '\t' { // So we can call this function multiple times (FIXME: Flaky)
			this.Glyphs = append(this.Glyphs, ImFontGlyph{})
		}
		var tab_glyph *ImFontGlyph = &this.Glyphs[len(this.Glyphs)-1]
		*tab_glyph = *this.FindGlyph((ImWchar)(' '))
		tab_glyph.Codepoint = '\t'
		tab_glyph.AdvanceX *= IM_TABSIZE
		this.IndexAdvanceX[(int)(tab_glyph.Codepoint)] = (float)(tab_glyph.AdvanceX)
		this.IndexLookup[(int)(tab_glyph.Codepoint)] = (ImWchar)(len(this.Glyphs) - 1)
	}

	// Mark special glyphs as not visible (note that AddGlyph already mark as non-visible glyphs with zero-size polygons)
	this.SetGlyphVisible((ImWchar)(' '), false)
	this.SetGlyphVisible((ImWchar)('\t'), false)

	// Ellipsis character is required for rendering elided text. We prefer using U+2026 (horizontal ellipsis).
	// However some old fonts may contain ellipsis at U+0085. Here we auto-detect most suitable ellipsis character.
	// FIXME: Note that 0x2026 is rarely included in our font ranges. Because of this we are more likely to use three individual dots.
	var ellipsis_chars = []ImWchar{(ImWchar)(0x2026), (ImWchar)(0x0085)}
	var dots_chars = []ImWchar{(ImWchar)('.'), (ImWchar)(0xFF0E)}
	if this.EllipsisChar == MaxImWchar {
		this.EllipsisChar = FindFirstExistingGlyph(this, ellipsis_chars, int(len(ellipsis_chars)))
	}
	if this.DotChar == MaxImWchar {
		this.DotChar = FindFirstExistingGlyph(this, dots_chars, int(len(dots_chars)))
	}

	// Setup fallback character
	var fallback_chars = []ImWchar{(ImWchar)(IM_UNICODE_CODEPOINT_INVALID), (ImWchar)('?'), (ImWchar)(' ')}
	this.FallbackGlyph = this.FindGlyphNoFallback(this.FallbackChar)
	if this.FallbackGlyph == nil {
		this.FallbackChar = FindFirstExistingGlyph(this, fallback_chars, int(len(fallback_chars)))
		this.FallbackGlyph = this.FindGlyphNoFallback(this.FallbackChar)
		if this.FallbackGlyph == nil {
			this.FallbackGlyph = &this.Glyphs[len(this.Glyphs)-1]
			this.FallbackChar = (ImWchar)(this.FallbackGlyph.Codepoint)
		}
	}

	this.FallbackAdvanceX = this.FallbackGlyph.AdvanceX
	for i := int(0); i < max_codepoint+1; i++ {
		if this.IndexAdvanceX[i] < 0.0 {
			this.IndexAdvanceX[i] = this.FallbackAdvanceX
		}
	}
}
