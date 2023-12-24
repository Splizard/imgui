package imgui

import "fmt"

// ImFontAtlasCustomRect See ImFontAtlas::AddCustomRectXXX functions.
type ImFontAtlasCustomRect struct {
	Width, Height uint16  // Input    // Desired rectangle dimension
	X, Y          uint16  // Output   // Packed position in Atlas
	GlyphID       uint    // Input    // For custom font glyphs only (ID < 0x110000)
	GlyphAdvanceX float   // Input    // For custom font glyphs only: glyph xadvance
	GlyphOffset   ImVec2  // Input    // For custom font glyphs only: glyph display offset
	Font          *ImFont // Input    // For custom font glyphs only: target font
}

func NewImFontAtlasCustomRect() ImFontAtlasCustomRect {
	return ImFontAtlasCustomRect{
		Width:         0,
		Height:        0,
		X:             0xFFFF,
		Y:             0xFFFF,
		GlyphID:       0,
		GlyphAdvanceX: 0,
		GlyphOffset:   ImVec2{},
		Font:          nil,
	}
}

func (r ImFontAtlasCustomRect) IsPacked() bool {
	return r.X != 0xFFFF
}

// ImFontAtlas Load and rasterize multiple TTF/OTF fonts into a same texture. The font atlas will build a single texture holding:
//   - One or more fonts.
//   - Custom graphics data needed to render the shapes needed by Dear ImGui.
//   - Mouse cursor shapes for software cursor rendering (unless setting 'Flags |= ImFontAtlasFlags_NoMouseCursors' in the font atlas).
//
// It is the user-code responsibility to setup/build the atlas, then upload the pixel data into a texture accessible by your graphics api.
//   - Optionally, call any of the AddFont*** functions. If you don't call any, the default font embedded in the code will be loaded for you.
//   - Call GetTexDataAsAlpha8() or GetTexDataAsRGBA32() to build and retrieve pixels data.
//   - Upload the pixels data into a texture within your graphics system (see imgui_impl_xxxx.cpp examples)
//   - Call SetTexID(my_tex_id); and pass the pointer/identifier to your texture in a format natural to your graphics API.
//     This value will be passed back to you during rendering to identify the texture. Read FAQ entry about ImTextureID for more details.
//
// Common pitfalls:
//   - If you pass a 'glyph_ranges' array to AddFont*** functions, you need to make sure that your array persist up until the
//     atlas is build (when calling GetTexData*** or Build()). We only copy the pointer, not the data.
//   - Important: By default, AddFontFromMemoryTTF() takes ownership of the data. Even though we are not writing to it, we will free the pointer on destruction.
//     You can set font_cfg->FontDataOwnedByAtlas=false to keep ownership of your data and it won't be freed,
//   - Even though many functions are suffixed with "TTF", OTF data is supported just as well.
//   - This is an old API and it is currently awkward for those and and various other reasons! We will address them in the future!
type ImFontAtlas struct {
	//-------------------------------------------
	// Members
	//-------------------------------------------

	Flags           ImFontAtlasFlags // Build flags (see ImFontAtlasFlags_)
	TexID           ImTextureID      // User data to refer to the texture once it has been uploaded to user's graphic systems. It is passed back to you during rendering via the ImDrawCmd structure.
	TexDesiredWidth int              // Texture width desired by user before Build(). Must be a power-of-two. If have many glyphs your graphics API have texture size restrictions you may want to increase texture width to decrease height.
	TexGlyphPadding int              // Padding between glyphs within texture in pixels. Defaults to 1. If your rendering method doesn't rely on bilinear filtering you may set this to 0.
	Locked          bool             // Marked as Locked by ImGui::NewFrame() so attempt to modify the atlas will assert.

	// [Internal]
	// NB: Access texture data via GetTexData*() calls! Which will setup a default font for you.
	TexReady           bool                                        // Set when texture was built matching current font input
	TexPixelsUseColors bool                                        // Tell whether our texture data is known to use colors (rather than just alpha channel), in order to help backend select a format.
	TexPixelsAlpha8    []byte                                      // 1 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight
	TexPixelsRGBA32    []uint                                      // 4 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight * 4
	TexWidth           int                                         // Texture width calculated during Build().
	TexHeight          int                                         // Texture height calculated during Build().
	TexUvScale         ImVec2                                      // = (1.0f/TexWidth, 1.0f/TexHeight)
	TexUvWhitePixel    ImVec2                                      // Texture coordinates to a white pixel
	Fonts              []*ImFont                                   // Hold all the fonts returned by AddFont*. Fonts[0] is the default font upon calling ImGui::NewFrame(), use ImGui::PushFont()/PopFont() to change the current font.
	CustomRects        []ImFontAtlasCustomRect                     // Rectangles for packing custom texture data into the atlas.
	ConfigData         []ImFontConfig                              // Configuration data
	TexUvLines         [IM_DRAWLIST_TEX_LINES_WIDTH_MAX + 1]ImVec4 // UVs for baked anti-aliased lines

	// [Internal] Font builder
	FontBuilderIO    *ImFontBuilderIO // Opaque interface to a font builder (default to stb_truetype, can be changed to use FreeType by defining IMGUI_ENABLE_FREETYPE).
	FontBuilderFlags uint             // Shared flags (for all fonts) for custom font builder. THIS IS BUILD IMPLEMENTATION DEPENDENT. Per-font override is also available in ImFontConfig.

	// [Internal] Packing data
	PackIdMouseCursors int // Custom texture rectangle ID for white pixel and mouse cursors
	PackIdLines        int // Custom texture rectangle ID for baked anti-aliased lines
}

func NewImFontAtlas() ImFontAtlas {
	return ImFontAtlas{
		TexGlyphPadding:    1,
		PackIdMouseCursors: -1,
		PackIdLines:        -1,
	}
}

func (atlas *ImFontAtlas) AddFont(font_cfg *ImFontConfig) *ImFont {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render")
	IM_ASSERT(font_cfg.FontData != nil && font_cfg.FontDataSize > 0)
	IM_ASSERT(font_cfg.SizePixels > 0.0)

	// Create new font
	if !font_cfg.MergeMode {
		f := NewImFont()
		atlas.Fonts = append(atlas.Fonts, &f)
	} else {
		IM_ASSERT_USER_ERROR(len(atlas.Fonts) != 0, "Cannot use MergeMode for the first font") // When using MergeMode make sure that a font has already been added before. You can use ImGui::GetIO().Fonts.AddFontDefault() to add the default imgui font.
	}

	atlas.ConfigData = append(atlas.ConfigData, *font_cfg)
	var new_font_cfg *ImFontConfig = &atlas.ConfigData[len(atlas.ConfigData)-1]
	if new_font_cfg.DstFont == nil {
		new_font_cfg.DstFont = atlas.Fonts[len(atlas.ConfigData)-1]
	}
	if !new_font_cfg.FontDataOwnedByAtlas {
		new_font_cfg.FontData = make([]byte, new_font_cfg.FontDataSize)
		new_font_cfg.FontDataOwnedByAtlas = true
		copy(new_font_cfg.FontData, font_cfg.FontData[:(size_t)(new_font_cfg.FontDataSize)])
	}

	if new_font_cfg.DstFont.EllipsisChar == (ImWchar)(-1) {
		new_font_cfg.DstFont.EllipsisChar = font_cfg.EllipsisChar
	}

	// Invalidate texture
	atlas.TexReady = false
	atlas.ClearTexData()
	return new_font_cfg.DstFont
}

func (atlas *ImFontAtlas) AddFontFromFileTTF(filename string, size_pixels float, font_cfg_template *ImFontConfig, glyph_ranges []ImWchar) *ImFont {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	var data_size size_t = 0
	var data = ImFileLoadToMemory(filename, "rb", &data_size, 0)
	if data == nil {
		IM_ASSERT_USER_ERROR(false, "Could not load font file!")
		return nil
	}

	var font_cfg ImFontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg = NewImFontConfig()
	}

	if font_cfg.Name == "" {
		// Store a copy of filename into into the font name for convenience
		font_cfg.Name = fmt.Sprintf("%s, %.0fpx", filename, size_pixels)
	}
	return atlas.AddFontFromMemoryTTF(data, (int)(data_size), size_pixels, &font_cfg, glyph_ranges)
}

func (atlas *ImFontAtlas) AddFontFromMemoryTTF(ttf_data []byte, ttf_size int, size_pixels float, font_cfg_template *ImFontConfig, glyph_ranges []ImWchar) *ImFont {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	var font_cfg ImFontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg = NewImFontConfig()
	}
	IM_ASSERT(font_cfg.FontData == nil)
	font_cfg.FontData = ttf_data
	font_cfg.FontDataSize = ttf_size
	if size_pixels > 0.0 {
		font_cfg.SizePixels = size_pixels
	}
	if glyph_ranges != nil {
		font_cfg.GlyphRanges = glyph_ranges
	}
	return atlas.AddFont(&font_cfg)
}

// ClearInputData Note: Transfer ownership of 'ttf_data' to ImFontAtlas! Will be deleted after destruction of the atlas. Set font_cfg->FontDataOwnedByAtlas=false to keep ownership of your data and it won't be freed.
// 'compressed_font_data_base85' still owned by caller. Compress with binary_to_compressed_c.cpp with -base85 parameter.
// Clear input data (all ImFontConfig structures including sizes, TTF data, glyph ranges, etc.) = all the data used to build the texture and fonts.
func (atlas *ImFontAtlas) ClearInputData() {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	for i := range atlas.ConfigData {
		if atlas.ConfigData[i].FontData != nil && atlas.ConfigData[i].FontDataOwnedByAtlas {
			atlas.ConfigData[i].FontData = nil
		}
	}

	atlas.ConfigData = nil
	atlas.CustomRects = nil
	atlas.PackIdMouseCursors = -1
	atlas.PackIdLines = -1
	// Important: we leave TexReady untouched
}

func (atlas *ImFontAtlas) ClearTexData() {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	atlas.TexPixelsAlpha8 = nil
	atlas.TexPixelsRGBA32 = nil
	atlas.TexPixelsUseColors = false
}

// ClearFonts Clear output texture data (CPU side). Saves RAM once the texture has been copied to graphics memory.
// Clear output font data (glyphs storage, UV coordinates).
func (atlas *ImFontAtlas) ClearFonts() {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")
	atlas.Fonts = nil
	atlas.TexReady = false
}

// Clear all input and output.
func (atlas *ImFontAtlas) Clear() {
	atlas.ClearInputData()
	atlas.ClearTexData()
	atlas.ClearFonts()
}

// GetTexDataAsRGBA32 4 bytes-per-pixel
func (atlas *ImFontAtlas) GetTexDataAsRGBA32(out_pixels *[]uint32, out_width, out_height, out_bytes_per_pixel *int) {
	// Convert to RGBA32 format on demand
	// Although it is likely to be the most commonly used format, our font rendering is 1 channel / 8 bpp
	if atlas.TexPixelsRGBA32 == nil {
		var pixels []byte = nil
		atlas.GetTexDataAsAlpha8(&pixels, nil, nil, nil)
		if pixels != nil {
			atlas.TexPixelsRGBA32 = make([]uint32, atlas.TexWidth*atlas.TexHeight*4)
			var src = pixels
			var dst = atlas.TexPixelsRGBA32
			for n := atlas.TexWidth * atlas.TexHeight; n > 0; n-- {
				dst[0] = IM_COL32(255, 255, 255, src[0])
				dst = dst[1:]
				src = src[1:]
			}
		}
	}

	*out_pixels = atlas.TexPixelsRGBA32
	if out_width != nil {
		*out_width = atlas.TexWidth
	}
	if out_height != nil {
		*out_height = atlas.TexHeight
	}
	if out_bytes_per_pixel != nil {
		*out_bytes_per_pixel = 4
	}
}

func (atlas *ImFontAtlas) IsBuilt() bool           { return len(atlas.Fonts) > 0 && atlas.TexReady } // Bit ambiguous: used to detect when user didn't built texture but effectively we should check TexID != 0 except that would be backend dependent...
func (atlas *ImFontAtlas) SetTexID(id ImTextureID) { atlas.TexID = id }

//-------------------------------------------
// Glyph Ranges
//-------------------------------------------

// Helpers to retrieve list of common Unicode ranges (2 value per range, values are inclusive, zero-terminated list)
// NB: Make sure that your string are UTF-8 and NOT in your local code page. In C++11, you can create UTF-8 string literal using the u8"Hello world" syntax. See FAQ for details.
// NB: Consider using ImFontGlyphRangesBuilder to build glyph ranges from textual data.

func UnpackAccumulativeOffsetsIntoRanges(base_codepoint int, accumulative_offsets []int16, accumulative_offsets_count int, out_ranges []ImWchar) {
	for n := int(0); n < accumulative_offsets_count; n, out_ranges = n+1, out_ranges[2:] {
		out_ranges[0] = (ImWchar)(base_codepoint + ImWchar(accumulative_offsets[n]))
		out_ranges[1] = (ImWchar)(base_codepoint + ImWchar(accumulative_offsets[n]))
		base_codepoint += int(accumulative_offsets[n])
	}
	out_ranges[0] = 0
}

//-------------------------------------------
// [BETA] Custom Rectangles/Glyphs API
//-------------------------------------------

// AddCustomRectRegular You can request arbitrary rectangles to be packed into the atlas, for your own purposes.
//   - After calling Build(), you can query the rectangle position and render your pixels.
//   - If you render colored output, set 'atlas->TexPixelsUseColors = true' as this may help some backends decide of prefered texture format.
//   - You can also request your rectangles to be mapped as font glyph (given a font + Unicode point),
//     so you can render e.guiContext. custom colorful icons and use them as regular glyphs.
//   - Read docs/FONTS.md for more details about using colorful icons.
//   - Note: this API may be redesigned later in order to support multi-monitor varying DPI settings.
func (atlas *ImFontAtlas) AddCustomRectRegular(width, height int) int {
	IM_ASSERT(width > 0 && width <= 0xFFFF)
	IM_ASSERT(height > 0 && height <= 0xFFFF)
	var r ImFontAtlasCustomRect
	r.Width = (uint16)(width)
	r.Height = (uint16)(height)
	atlas.CustomRects = append(atlas.CustomRects, r)
	return int(len(atlas.CustomRects)) - 1 // Return index
}

func (atlas *ImFontAtlas) AddCustomRectFontGlyph(font *ImFont, id ImWchar, width, height int, advance_x float, offset *ImVec2) int {
	IM_ASSERT(font != nil)
	IM_ASSERT(width > 0 && width <= 0xFFFF)
	IM_ASSERT(height > 0 && height <= 0xFFFF)
	var r ImFontAtlasCustomRect
	r.Width = (uint16)(width)
	r.Height = (uint16)(height)
	r.GlyphID = uint(id)
	r.GlyphAdvanceX = advance_x
	r.GlyphOffset = *offset
	r.Font = font
	atlas.CustomRects = append(atlas.CustomRects, r)
	return int(len(atlas.CustomRects) - 1) // Return index
}

func (atlas *ImFontAtlas) GetCustomRectByIndex(index int) *ImFontAtlasCustomRect {
	IM_ASSERT(index >= 0)
	return &atlas.CustomRects[index]
}

// CalcCustomRectUV [Internal]
func (atlas *ImFontAtlas) CalcCustomRectUV(rect *ImFontAtlasCustomRect, out_uv_min, out_uv_max *ImVec2) {
	IM_ASSERT(atlas.TexWidth > 0 && atlas.TexHeight > 0) // Font atlas needs to be built before we can calculate UV coordinates
	IM_ASSERT(rect.IsPacked())                           // Make sure the rectangle has been packed
	*out_uv_min = ImVec2{(float)(rect.X) * atlas.TexUvScale.x, (float)(rect.Y) * atlas.TexUvScale.y}
	*out_uv_max = ImVec2{(float)(rect.X+rect.Width) * atlas.TexUvScale.x, (float)(rect.Y+rect.Height) * atlas.TexUvScale.y}
}

func (atlas *ImFontAtlas) GetMouseCursorTexData(cursor_type ImGuiMouseCursor, out_offset, out_size *ImVec2, out_uv_border *[2]ImVec2, out_uv_fill *[2]ImVec2) bool {
	if cursor_type <= ImGuiMouseCursor_None || cursor_type >= ImGuiMouseCursor_COUNT {
		return false
	}
	if atlas.Flags&ImFontAtlasFlags_NoMouseCursors != 0 {
		return false
	}

	IM_ASSERT(atlas.PackIdMouseCursors != -1)
	var r *ImFontAtlasCustomRect = atlas.GetCustomRectByIndex(atlas.PackIdMouseCursors)
	var pos ImVec2 = FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA[cursor_type][0].Add(ImVec2{(float)(r.X), (float)(r.Y)})
	var size ImVec2 = FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA[cursor_type][1]
	*out_size = size
	*out_offset = FONT_ATLAS_DEFAULT_TEX_CURSOR_DATA[cursor_type][2]
	out_uv_border[0] = (pos).Mul(atlas.TexUvScale)
	out_uv_border[1] = (pos.Add(size)).Mul(atlas.TexUvScale)
	pos.x += float(FONT_ATLAS_DEFAULT_TEX_DATA_W + 1)
	out_uv_fill[0] = (pos).Mul(atlas.TexUvScale)
	out_uv_fill[1] = (pos.Add(size)).Mul(atlas.TexUvScale)
	return true
}
