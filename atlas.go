package imgui

import (
	_ "embed"
	"fmt"
	"unsafe"

	"github.com/splizard/imgui/stb/stbrp"
	"github.com/splizard/imgui/stb/stbtt"
)

//-----------------------------------------------------------------------------
// [SECTION] Default font data (ProggyClean.ttf)
//-----------------------------------------------------------------------------
// ProggyClean.ttf
// Copyright (c) 2004, 2005 Tristan Grimmer
// MIT license (see License.txt in http://www.upperbounds.net/download/ProggyClean.ttf.zip)
// Download and more information at http://upperbounds.net
//-----------------------------------------------------------------------------
// File: 'ProggyClean.ttf' (41208 bytes)
// Exported using misc/fonts/binary_to_compressed_c.cpp (with compression + base85 string encoding).
// The purpose of encoding as base85 instead of "0x00,0x01,..." style is only save on _source code_ size.
//-----------------------------------------------------------------------------
//go:embed proggy.ttf
var proggy_clean_ttf_decompressed_data_base85 []byte

// A work of art lies ahead! (. = white layer, X = black layer, others are blank)
// The 2x2 white texels on the top left are the ones we'll use everywhere in Dear ImGui to render filled shapes.
const FONT_ATLAS_DEFAULT_TEX_DATA_W int = 108 // Actual texture will be 2 times that + 1 spacing.
const FONT_ATLAS_DEFAULT_TEX_DATA_H int = 27
const FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS = "" +
	"..-         -XXXXXXX-    X    -           X           -XXXXXXX          -          XXXXXXX-     XX          " +
	"..-         -X.....X-   X.X   -          X.X          -X.....X          -          X.....X-    X..X         " +
	"---         -XXX.XXX-  X...X  -         X...X         -X....X           -           X....X-    X..X         " +
	"X           -  X.X  - X.....X -        X.....X        -X...X            -            X...X-    X..X         " +
	"XX          -  X.X  -X.......X-       X.......X       -X..X.X           -           X.X..X-    X..X         " +
	"X.X         -  X.X  -XXXX.XXXX-       XXXX.XXXX       -X.X X.X          -          X.X X.X-    X..XXX       " +
	"X..X        -  X.X  -   X.X   -          X.X          -XX   X.X         -         X.X   XX-    X..X..XXX    " +
	"X...X       -  X.X  -   X.X   -    XX    X.X    XX    -      X.X        -        X.X      -    X..X..X..XX  " +
	"X....X      -  X.X  -   X.X   -   X.X    X.X    X.X   -       X.X       -       X.X       -    X..X..X..X.X " +
	"X.....X     -  X.X  -   X.X   -  X..X    X.X    X..X  -        X.X      -      X.X        -XXX X..X..X..X..X" +
	"X......X    -  X.X  -   X.X   - X...XXXXXX.XXXXXX...X -         X.X   XX-XX   X.X         -X..XX........X..X" +
	"X.......X   -  X.X  -   X.X   -X.....................X-          X.X X.X-X.X X.X          -X...X...........X" +
	"X........X  -  X.X  -   X.X   - X...XXXXXX.XXXXXX...X -           X.X..X-X..X.X           - X..............X" +
	"X.........X -XXX.XXX-   X.X   -  X..X    X.X    X..X  -            X...X-X...X            -  X.............X" +
	"X..........X-X.....X-   X.X   -   X.X    X.X    X.X   -           X....X-X....X           -  X.............X" +
	"X......XXXXX-XXXXXXX-   X.X   -    XX    X.X    XX    -          X.....X-X.....X          -   X............X" +
	"X...X..X    ---------   X.X   -          X.X          -          XXXXXXX-XXXXXXX          -   X...........X " +
	"X..X X..X   -       -XXXX.XXXX-       XXXX.XXXX       -------------------------------------    X..........X " +
	"X.X  X..X   -       -X.......X-       X.......X       -    XX           XX    -           -    X..........X " +
	"XX    X..X  -       - X.....X -        X.....X        -   X.X           X.X   -           -     X........X  " +
	"      X..X          -  X...X  -         X...X         -  X..X           X..X  -           -     X........X  " +
	"       XX           -   X.X   -          X.X          - X...XXXXXXXXXXXXX...X -           -     XXXXXXXXXX  " +
	"------------        -    X    -           X           -X.....................X-           ------------------" +
	"                    ----------------------------------- X...XXXXXXXXXXXXX...X -                             " +
	"                                                      -  X..X           X..X  -                             " +
	"                                                      -   X.X           X.X   -                             " +
	"                                                      -    XX           XX    -                             "

// Temporary data for one source font           u tiple source fonts can be merged into one destination ImFont)
// (C doesn't allow instancing ImVectr>wth function-local types so we declare the type here.)
type ImFontBuildSrcData struct {
	FontInfo      stbtt.FontInfo
	PackRange     [1]stbtt.PackRange // Hold the list of codepoints to pack (essentially points to Codepoints.Data)
	Rects         []stbrp.Rect       // Rectangle to pack. We first fill in their size and the packer will give us their position.
	PackedChars   []stbtt.PackedChar // Output glyphs
	SrcRanges     []ImWchar          // Ranges as requested by user (user is allowed to request too much, e.g. 0x0020..0xFFFF)
	DstIndex      int                // Index into atlas.Fonts[] and dst_tmp_array[]
	GlyphsHighest int                // Highest requested codepoint
	GlyphsCount   int                // Glyph count (excluding missing glyphs and glyphs already set by an earlier source font)
	GlyphsSet     ImBitVector        // Glyph bit map (random access, 1-bit per codepoint. This will be a maximum of 8KB)
	GlyphsList    []int              // Glyph codepoints list (flattened version of GlyphsMap)
}

// Temporary data for one destination ImFont* (multiple source fonts can be merged into one destination ImFont)
type ImFontBuildDstData struct {
	SrcCount      int // Number of source fonts targeting this destination font.
	GlyphsHighest int
	GlyphsCount   int
	GlyphsSet     ImBitVector // This is used to resolve collision when multiple sources are merged into a same destination font.
}

func UnpackBitVectorToFlatIndexList(in *ImBitVector, out *[]int) {
	IM_ASSERT(unsafe.Sizeof((*in)[0]) == unsafe.Sizeof(int(0)))
	for i, it := range *in {
		if entries_32 := it; entries_32 != 0 {
			for bit_n := 0; bit_n < 32; bit_n++ {
				if entries_32&((ImU32)(1<<bit_n)) != 0 {
					*out = append(*out, ((int)(((i) << 5) + bit_n)))
				}
			}
		}
	}
}

func IM_ASSERT(x bool) {
	if !x {
		panic("imgui: IM_ASSERT failed")
	}
}

// 1 byte per-pixel
func (atlas *ImFontAtlas) GetTexDataAsAlpha8(out_pixels *[]byte, out_width, out_height, out_bytes_per_pixel *int) {
	// Build atlas on demand
	if atlas.TexPixelsAlpha8 == nil {
		atlas.Build()
	}

	*out_pixels = atlas.TexPixelsAlpha8
	if out_width != nil {
		*out_width = atlas.TexWidth
	}
	if out_height != nil {
		*out_height = atlas.TexHeight
	}
	if out_bytes_per_pixel != nil {
		*out_bytes_per_pixel = 1
	}
}

// Load embedded ProggyClean.ttf at size 13, disable oversampling
func (atlas *ImFontAtlas) AddFontDefault(font_cfg_template *ImFontConfig) *ImFont {
	var font_cfg ImFontConfig
	if font_cfg_template != nil {
		font_cfg = *font_cfg_template
	} else {
		font_cfg = NewImFontConfig()
	}

	if font_cfg_template == nil {
		font_cfg.OversampleH = 1
		font_cfg.OversampleV = 1
		font_cfg.PixelSnapH = true
	}
	if font_cfg.SizePixels <= 0.0 {
		font_cfg.SizePixels = 13.0 * 1.0
	}
	if font_cfg.Name == "" {
		font_cfg.Name = fmt.Sprintf("ProggyClean.ttf, %dpx", (int)(font_cfg.SizePixels))
	}
	font_cfg.EllipsisChar = (ImWchar)(0x0085)
	font_cfg.GlyphOffset.y = 1.0 * IM_FLOOR(font_cfg.SizePixels/13.0) // Add +1 offset per 13 units

	var glyph_ranges []ImWchar
	if font_cfg.GlyphRanges != nil {
		glyph_ranges = font_cfg.GlyphRanges
	} else {
		glyph_ranges = atlas.GetGlyphRangesDefault()
	}

	data := proggy_clean_ttf_decompressed_data_base85

	return atlas.AddFontFromMemoryTTF(proggy_clean_ttf_decompressed_data_base85, int32(len(data)), font_cfg.SizePixels, &font_cfg, glyph_ranges)
}

// Build atlas, retrieve pixel data.
// User is in charge of copying the pixels into graphics memory (e.g. create a texture with your engine). Then store your texture handle with SetTexID().
// The pitch is always = Width * BytesPerPixels (1 or 4)
// Building in RGBA32 format is provided for convenience and compatibility, but note that unless you manually manipulate or copy color data into
// the texture (e.g. when using the AddCustomRect*** api), then the RGB pixels emitted will always be white (~75% of memory/bandwidth waste.
// Build pixels data. This is called automatically for you by the GetTexData*** functions.
func (atlas *ImFontAtlas) Build() bool {
	IM_ASSERT_USER_ERROR(!atlas.Locked, "Cannot modify a locked ImFontAtlas between NewFrame() and EndFrame/Render()!")

	// Default font is none are specified
	if len(atlas.ConfigData) == 0 {
		atlas.AddFontDefault(nil)
	}

	// Select builder
	// - Note that we do not reassign to atlas->FontBuilderIO, since it is likely to point to static data which
	//   may mess with some hot-reloading schemes. If you need to assign to this (for dynamic selection) AND are
	//   using a hot-reloading scheme that messes up static data, store your own instance of ImFontBuilderIO somewhere
	//   and point to it instead of pointing directly to return value of the GetBuilderXXX functions.
	var builder_io = atlas.FontBuilderIO
	if builder_io == nil {
		builder_io = ImFontAtlasGetBuilderForStbTruetype()
	}

	// Build
	return builder_io.FontBuilder_Build(atlas)
}

// Retrieve list of range (2 int per range, values are inclusive)
func (atlas *ImFontAtlas) GetGlyphRangesDefault() []ImWchar {
	return []ImWchar{
		0x0020, 0x00FF, // Basic Latin + Latin Supplement
		0,
	}
}

func Decode85Byte(c char) uint {
	if c >= '\\' {
		return uint(c) - 36
	}
	return uint(c) - 35
}

func Decode85(src string, dst []byte) {
	for len(src) >= 0 {
		var tmp uint = Decode85Byte(src[0]) + 85*(Decode85Byte(src[1])+85*(Decode85Byte(src[2])+85*(Decode85Byte(src[3])+85*Decode85Byte(src[4]))))
		dst[0] = (char(tmp>>0) & 0xFF)
		dst[1] = (char(tmp>>8) & 0xFF)
		dst[2] = (char(tmp>>16) & 0xFF)
		dst[3] = (char(tmp>>24) & 0xFF) // We can't assume little-endianness.
		src = src[5:]
		dst = dst[4:]
	}
}

func ImFontAtlasGetBuilderForStbTruetype() *ImFontBuilderIO {
	var io ImFontBuilderIO
	io.FontBuilder_Build = ImFontAtlasBuildWithStbTruetype
	return &io
}

// Note: this is called / shared by both the stb_truetype and the FreeType builder
func ImFontAtlasBuildInit(atlas *ImFontAtlas) {
	// Register texture region for mouse cursors or standard white pixels
	if atlas.PackIdMouseCursors < 0 {
		if atlas.Flags&ImFontAtlasFlags_NoMouseCursors == 0 {
			atlas.PackIdMouseCursors = atlas.AddCustomRectRegular(FONT_ATLAS_DEFAULT_TEX_DATA_W*2+1, FONT_ATLAS_DEFAULT_TEX_DATA_H)
		} else {
			atlas.PackIdMouseCursors = atlas.AddCustomRectRegular(2, 2)
		}
	}

	// Register texture region for thick lines
	// The +2 here is to give space for the end caps, whilst height +1 is to accommodate the fact we have a zero-width row
	if atlas.PackIdLines < 0 {
		if atlas.Flags&ImFontAtlasFlags_NoBakedLines == 0 {
			atlas.PackIdLines = atlas.AddCustomRectRegular(IM_DRAWLIST_TEX_LINES_WIDTH_MAX+2, IM_DRAWLIST_TEX_LINES_WIDTH_MAX+1)
		}
	}
}

func ImFontAtlasBuildPackCustomRects(atlas *ImFontAtlas, stbrp_context_opaque interface{}) {
	var pack_context *stbrp.Context = stbrp_context_opaque.(*stbrp.Context)
	IM_ASSERT(pack_context != nil)

	var user_rects []ImFontAtlasCustomRect = atlas.CustomRects
	IM_ASSERT(len(user_rects) >= 1) // We expect at least the default custom rects to be registered, else something went wrong.

	var pack_rects []stbrp.Rect = make([]stbrp.Rect, len(user_rects))

	for i := range user_rects {
		pack_rects[i].W = stbrp.Coord(user_rects[i].Width)
		pack_rects[i].H = stbrp.Coord(user_rects[i].Height)
	}
	stbrp.PackRects(pack_context, pack_rects, int(len(pack_rects)))
	for i := range pack_rects {
		if pack_rects[i].WasPacked != 0 {
			user_rects[i].X = uint16(pack_rects[i].X)
			user_rects[i].Y = uint16(pack_rects[i].Y)
			IM_ASSERT(uint16(pack_rects[i].W) == user_rects[i].Width && uint16(pack_rects[i].H) == user_rects[i].Height)
			atlas.TexHeight = ImMaxInt(atlas.TexHeight, int(pack_rects[i].Y)+int(pack_rects[i].H))
		}
	}
}

func ImFontAtlasBuildSetupFont(atlas *ImFontAtlas, font *ImFont, font_config *ImFontConfig, ascent float32, descent float32) {
	if !font_config.MergeMode {
		font.ClearOutputData()
		font.FontSize = font_config.SizePixels
		font.ConfigData = []ImFontConfig{*font_config}
		font.ConfigDataCount = 0
		font.ContainerAtlas = atlas
		font.Ascent = ascent
		font.Descent = descent
	}
	font.ConfigDataCount++
}

func ImFontAtlasBuildRender8bppRectFromString(atlas *ImFontAtlas, x, y, w, h int, in_str string, in_marker_char byte, in_marker_pixel_value byte) {
	IM_ASSERT(x >= 0 && x+w <= atlas.TexWidth)
	IM_ASSERT(y >= 0 && y+h <= atlas.TexHeight)
	var out_pixel []byte = atlas.TexPixelsAlpha8[x+(y*atlas.TexWidth):]
	for off_y := int(0); off_y < h; off_y, out_pixel, in_str = off_y+1, out_pixel[atlas.TexWidth:], in_str[w:] {
		for off_x := int(0); off_x < w; off_x++ {
			if in_str[off_x] == in_marker_char {
				out_pixel[off_x] = in_marker_pixel_value
			} else {
				out_pixel[off_x] = 0x00
			}
		}
	}
}

func ImFontAtlasBuildRender32bppRectFromString(atlas *ImFontAtlas, x, y, w, h int, in_str string, in_marker_char byte, in_marker_pixel_value uint32) {
	IM_ASSERT(x >= 0 && x+w <= atlas.TexWidth)
	IM_ASSERT(y >= 0 && y+h <= atlas.TexHeight)
	var out_pixel []uint = atlas.TexPixelsRGBA32[x+(y*atlas.TexWidth):]
	for off_y := int(0); off_y < h; off_y, out_pixel, in_str = off_y+1, out_pixel[atlas.TexWidth:], in_str[w:] {
		for off_x := int(0); off_x < w; off_x++ {
			if in_str[off_x] == in_marker_char {
				out_pixel[off_x] = in_marker_pixel_value
			} else {
				out_pixel[off_x] = IM_COL32_BLACK_TRANS
			}
		}
	}
}

func ImFontAtlasBuildRenderDefaultTexData(atlas *ImFontAtlas) {
	var r *ImFontAtlasCustomRect = atlas.GetCustomRectByIndex(atlas.PackIdMouseCursors)
	IM_ASSERT(r.IsPacked())

	var w int = atlas.TexWidth
	if atlas.Flags&ImFontAtlasFlags_NoMouseCursors == 0 {
		// Render/copy pixels
		IM_ASSERT(int(r.Width) == FONT_ATLAS_DEFAULT_TEX_DATA_W*2+1 && int(r.Height) == FONT_ATLAS_DEFAULT_TEX_DATA_H)
		var x_for_white int = int(r.X)
		var x_for_black int = int(r.X) + FONT_ATLAS_DEFAULT_TEX_DATA_W + 1
		if atlas.TexPixelsAlpha8 != nil {
			ImFontAtlasBuildRender8bppRectFromString(atlas, x_for_white, int(r.Y), FONT_ATLAS_DEFAULT_TEX_DATA_W, FONT_ATLAS_DEFAULT_TEX_DATA_H, FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS, '.', 0xFF)
			ImFontAtlasBuildRender8bppRectFromString(atlas, x_for_black, int(r.Y), FONT_ATLAS_DEFAULT_TEX_DATA_W, FONT_ATLAS_DEFAULT_TEX_DATA_H, FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS, 'X', 0xFF)
		} else {
			ImFontAtlasBuildRender32bppRectFromString(atlas, x_for_white, int(r.Y), FONT_ATLAS_DEFAULT_TEX_DATA_W, FONT_ATLAS_DEFAULT_TEX_DATA_H, FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS, '.', IM_COL32_WHITE)
			ImFontAtlasBuildRender32bppRectFromString(atlas, x_for_black, int(r.Y), FONT_ATLAS_DEFAULT_TEX_DATA_W, FONT_ATLAS_DEFAULT_TEX_DATA_H, FONT_ATLAS_DEFAULT_TEX_DATA_PIXELS, 'X', IM_COL32_WHITE)
		}
	} else {
		// Render 4 white pixels
		IM_ASSERT(r.Width == 2 && r.Height == 2)
		var offset int = (int)(r.X) + (int)(r.Y)*w
		if atlas.TexPixelsAlpha8 != nil {
			atlas.TexPixelsAlpha8[offset] = 0xFF
			atlas.TexPixelsAlpha8[offset+1] = 0xFF
			atlas.TexPixelsAlpha8[offset+w] = 0xFF
			atlas.TexPixelsAlpha8[offset+w+1] = 0xFF
		} else {
			atlas.TexPixelsRGBA32[offset] = IM_COL32_WHITE
			atlas.TexPixelsRGBA32[offset+1] = IM_COL32_WHITE
			atlas.TexPixelsRGBA32[offset+w] = IM_COL32_WHITE
			atlas.TexPixelsRGBA32[offset+w+1] = IM_COL32_WHITE
		}
	}
	atlas.TexUvWhitePixel = ImVec2{(float(r.X) + 0.5) * atlas.TexUvScale.x, (float(r.Y) + 0.5) * atlas.TexUvScale.y}
}

func ImFontAtlasBuildRenderLinesTexData(atlas *ImFontAtlas) {
	if atlas.Flags&ImFontAtlasFlags_NoBakedLines != 0 {
		return
	}

	// This generates a triangular shape in the texture, with the various line widths stacked on top of each other to allow interpolation between them
	var r *ImFontAtlasCustomRect = atlas.GetCustomRectByIndex(atlas.PackIdLines)
	IM_ASSERT(r.IsPacked())
	for n := 0; n < IM_DRAWLIST_TEX_LINES_WIDTH_MAX+1; n++ { // +1 because of the zero-width row
		// Each line consists of at least two empty pixels at the ends, with a line of solid pixels in the middle
		var y uint = uint(n)
		var line_width uint = uint(n)
		var pad_left uint = (uint(r.Width) - uint(line_width)) / 2
		var pad_right uint = uint(r.Width) - (uint(pad_left) + uint(line_width))

		// Write each slice
		IM_ASSERT(pad_left+line_width+pad_right == uint(r.Width) && y < uint(r.Height)) // Make sure we're inside the texture bounds before we start writing pixels
		if atlas.TexPixelsAlpha8 != nil {
			var write_ptr []byte = atlas.TexPixelsAlpha8[uint(r.X)+((uint(r.Y)+uint(y))*uint(atlas.TexWidth)):]
			for i := uint(0); i < pad_left; i++ {
				write_ptr[i] = 0x00
			}

			for i := uint(0); i < line_width; i++ {
				write_ptr[pad_left+i] = 0xFF
			}

			for i := uint(0); i < pad_right; i++ {
				write_ptr[pad_left+line_width+i] = 0x00
			}
		} else {
			var write_ptr []uint = atlas.TexPixelsRGBA32[uint(r.X)+((uint(r.Y)+uint(y))*uint(atlas.TexWidth)):]
			for i := uint(0); i < pad_left; i++ {
				write_ptr[i] = IM_COL32_BLACK_TRANS
			}

			for i := uint(0); i < line_width; i++ {
				write_ptr[pad_left+i] = IM_COL32_WHITE
			}

			for i := uint(0); i < pad_right; i++ {
				write_ptr[pad_left+line_width+i] = IM_COL32_BLACK_TRANS
			}
		}

		// Calculate UVs for this line
		var uv0 ImVec2 = ImVec2{(float)(uint(r.X) + pad_left - 1), (float)(uint(r.Y) + y)}.Mul(atlas.TexUvScale)
		var uv1 ImVec2 = ImVec2{(float)(uint(r.X) + pad_left + line_width + 1), (float)(uint(r.Y) + y + 1)}.Mul(atlas.TexUvScale)
		var half_v float = (uv0.y + uv1.y) * 0.5 // Calculate a constant V in the middle of the row to avoid sampling artifacts
		atlas.TexUvLines[n] = ImVec4{uv0.x, half_v, uv1.x, half_v}
	}
}

// This is called/shared by both the stb_truetype and the FreeType builder.
func ImFontAtlasBuildFinish(atlas *ImFontAtlas) {
	// Render into our custom data blocks
	IM_ASSERT(atlas.TexPixelsAlpha8 != nil || atlas.TexPixelsRGBA32 != nil)
	ImFontAtlasBuildRenderDefaultTexData(atlas)
	ImFontAtlasBuildRenderLinesTexData(atlas)

	// Register custom rectangle glyphs
	for i := range atlas.CustomRects {
		var r *ImFontAtlasCustomRect = &atlas.CustomRects[i]
		if r.Font == nil || r.GlyphID == 0 {
			continue
		}

		// Will ignore ImFontConfig settings: GlyphMinAdvanceX, GlyphMinAdvanceY, GlyphExtraSpacing, PixelSnapH
		IM_ASSERT(r.Font.ContainerAtlas == atlas)
		var uv0, uv1 ImVec2
		atlas.CalcCustomRectUV(r, &uv0, &uv1)
		r.Font.AddGlyph(nil, (ImWchar)(r.GlyphID), r.GlyphOffset.x, r.GlyphOffset.y, r.GlyphOffset.x+float(r.Width), r.GlyphOffset.y+float(r.Height), uv0.x, uv0.y, uv1.x, uv1.y, r.GlyphAdvanceX)
	}

	// Build all fonts lookup tables
	for i := range atlas.Fonts {
		if atlas.Fonts[i].DirtyLookupTables {
			atlas.Fonts[i].BuildLookupTable()
		}
	}

	atlas.TexReady = true
}

func ImFontAtlasBuildWithStbTruetype(atlas *ImFontAtlas) bool {
	IM_ASSERT(len(atlas.ConfigData) > 0)

	ImFontAtlasBuildInit(atlas)

	// Clear atlas
	atlas.TexID = (ImTextureID)(0)
	atlas.TexWidth = 0
	atlas.TexHeight = 0
	atlas.TexUvScale = ImVec2{}
	atlas.TexUvWhitePixel = ImVec2{}
	atlas.ClearTexData()

	// Temporary storage for building
	var src_tmp_array = make([]ImFontBuildSrcData, len(atlas.ConfigData))
	var dst_tmp_array = make([]ImFontBuildDstData, len(atlas.Fonts))

	// 1. Initialize font loading structure, check font data validity
	for src_i := range atlas.ConfigData {
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		var cfg *ImFontConfig = &atlas.ConfigData[src_i]
		IM_ASSERT(cfg.DstFont != nil && (!cfg.DstFont.IsLoaded() || cfg.DstFont.ContainerAtlas == atlas))

		// Find index from cfg.DstFont (we allow the user to set cfg.DstFont. Also it makes casual debugging nicer than when storing indices)
		src_tmp.DstIndex = -1
		for output_i := 0; output_i < len(atlas.Fonts) && src_tmp.DstIndex == -1; output_i++ {
			if cfg.DstFont == atlas.Fonts[output_i] {
				src_tmp.DstIndex = int(output_i)
			}
		}
		if src_tmp.DstIndex == -1 {
			IM_ASSERT(src_tmp.DstIndex != -1) // cfg.DstFont not pointing within atlas.Fonts[] array?
			return false
		}
		// Initialize helper structure for font loading and verify that the TTF/OTF data is correct
		var font_offset int = stbtt.GetFontOffsetForIndex(cfg.FontData, cfg.FontNo)
		IM_ASSERT_USER_ERROR(font_offset >= 0, "FontData is incorrect, or FontNo cannot be found.")
		if stbtt.InitFont(&src_tmp.FontInfo, cfg.FontData, font_offset) == 0 {
			return false
		}

		// Measure highest codepoints
		var dst_tmp *ImFontBuildDstData = &dst_tmp_array[src_tmp.DstIndex]
		if cfg.GlyphRanges != nil {
			src_tmp.SrcRanges = cfg.GlyphRanges
		} else {
			src_tmp.SrcRanges = atlas.GetGlyphRangesDefault()
		}

		for src_range := src_tmp.SrcRanges; src_range[0] != 0 && src_range[1] != 0; src_range = src_range[2:] {
			src_tmp.GlyphsHighest = ImMaxInt(src_tmp.GlyphsHighest, (int)(src_range[1]))
		}
		dst_tmp.SrcCount++
		dst_tmp.GlyphsHighest = ImMaxInt(dst_tmp.GlyphsHighest, src_tmp.GlyphsHighest)
	}

	// 2. For every requested codepoint, check for their presence in the font data, and handle redundancy or overlaps between source fonts to avoid unused glyphs.
	var total_glyphs_count int = 0
	for src_i := range src_tmp_array {
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		var dst_tmp *ImFontBuildDstData = &dst_tmp_array[src_tmp.DstIndex]
		src_tmp.GlyphsSet.Create(src_tmp.GlyphsHighest + 1)
		if len(dst_tmp.GlyphsSet) == 0 {
			dst_tmp.GlyphsSet.Create(dst_tmp.GlyphsHighest + 1)
		}

		for src_range := src_tmp.SrcRanges; src_range[0] != 0 && src_range[1] != 0; src_range = src_range[2:] {
			for codepoint := src_range[0]; codepoint <= src_range[1]; codepoint++ {
				if dst_tmp.GlyphsSet.TestBit(int(codepoint)) { // Don't overwrite existing glyphs. We could make this an option for MergeMode (e.g. MergeOverwrite==true)
					continue
				}
				if stbtt.FindGlyphIndex(&src_tmp.FontInfo, int(codepoint)) == 0 { // It is actually in the font?
					continue
				}

				// Add to avail set/counters
				src_tmp.GlyphsCount++
				dst_tmp.GlyphsCount++
				src_tmp.GlyphsSet.SetBit(int(codepoint))
				dst_tmp.GlyphsSet.SetBit(int(codepoint))
				total_glyphs_count++
			}
		}
	}

	// 3. Unpack our bit map into a flat list (we now have all the Unicode points that we know are requested _and_ available _and_ not overlapping another)
	for src_i := range src_tmp_array {
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		src_tmp.GlyphsList = make([]int, 0, src_tmp.GlyphsCount)
		UnpackBitVectorToFlatIndexList(&src_tmp.GlyphsSet, &src_tmp.GlyphsList)
		src_tmp.GlyphsSet.Clear()
		IM_ASSERT(int(len(src_tmp.GlyphsList)) == src_tmp.GlyphsCount)

	}
	for dst_i := range dst_tmp_array {
		dst_tmp_array[dst_i].GlyphsSet.Clear()
	}
	dst_tmp_array = dst_tmp_array[:0]

	// Allocate packing character data and flag packed characters buffer as non-packed (x0=y0=x1=y1=0)
	// (We technically don't need to zero-clear buf_rects, but let's do it for the sake of sanity)
	var buf_rects []stbrp.Rect = make([]stbrp.Rect, total_glyphs_count)
	var buf_packedchars []stbtt.PackedChar = make([]stbtt.PackedChar, total_glyphs_count)

	// 4. Gather glyphs sizes so we can pack them in our virtual canvas.
	var total_surface int = 0
	var buf_rects_out_n int = 0
	var buf_packedchars_out_n int = 0
	for src_i := range src_tmp_array {
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		if src_tmp.GlyphsCount == 0 {
			continue
		}

		src_tmp.Rects = buf_rects[buf_rects_out_n:]
		src_tmp.PackedChars = buf_packedchars[buf_packedchars_out_n:]
		buf_rects_out_n += src_tmp.GlyphsCount
		buf_packedchars_out_n += src_tmp.GlyphsCount

		// Convert our ranges in the format stb_truetype wants
		var cfg *ImFontConfig = &atlas.ConfigData[src_i]
		src_tmp.PackRange[0].FontSize = cfg.SizePixels
		src_tmp.PackRange[0].FirstUnicodeCodepointInRange = 0
		src_tmp.PackRange[0].ArrayOfUnicodeCodepoints = src_tmp.GlyphsList
		src_tmp.PackRange[0].NumChars = int(len(src_tmp.GlyphsList))
		src_tmp.PackRange[0].ChardataForRange = src_tmp.PackedChars
		src_tmp.PackRange[0].Oversample.H = (byte)(cfg.OversampleH)
		src_tmp.PackRange[0].Oversample.V = (byte)(cfg.OversampleV)

		// Gather the sizes of all rectangles we will need to pack (this loop is based on stbtt_PackFontRangesGatherRects)
		var scale float
		if cfg.SizePixels > 0 {
			scale = stbtt.ScaleForPixelHeight(&src_tmp.FontInfo, cfg.SizePixels)
		} else {
			scale = stbtt.ScaleForMappingEmToPixels(&src_tmp.FontInfo, -cfg.SizePixels)
		}
		var padding int = atlas.TexGlyphPadding
		for glyph_i := range src_tmp.GlyphsList {
			var x0, y0, x1, y1 int
			var glyph_index_in_font int = stbtt.FindGlyphIndex(&src_tmp.FontInfo, src_tmp.GlyphsList[glyph_i])
			IM_ASSERT(glyph_index_in_font != 0)
			stbtt.GetGlyphBitmapBoxSubpixel(&src_tmp.FontInfo, glyph_index_in_font, scale*float(cfg.OversampleH), scale*float(cfg.OversampleV), 0, 0, &x0, &y0, &x1, &y1)
			src_tmp.Rects[glyph_i].W = (stbrp.Coord)(x1 - x0 + padding + cfg.OversampleH - 1)
			src_tmp.Rects[glyph_i].H = (stbrp.Coord)(y1 - y0 + padding + cfg.OversampleV - 1)
			total_surface += int(src_tmp.Rects[glyph_i].W) * int(src_tmp.Rects[glyph_i].H)
		}
	}

	// We need a width for the skyline algorithm, any width!
	// The exact width doesn't really matter much, but some API/GPU have texture size limitations and increasing width can decrease height.
	// User can override TexDesiredWidth and TexGlyphPadding if they wish, otherwise we use a simple heuristic to select the width based on expected surface.
	var surface_sqrt int = (int)(ImSqrt((float)(total_surface)) + 1)
	atlas.TexHeight = 0
	if atlas.TexDesiredWidth > 0 {
		atlas.TexWidth = atlas.TexDesiredWidth
	} else {
		if surface_sqrt >= int(ImFloor(4096*0.7)) {
			atlas.TexWidth = 4096
		} else {
			if surface_sqrt >= int(ImFloor(2048*0.7)) {
				atlas.TexWidth = 2048
			} else {
				if surface_sqrt >= int(ImFloor(1024*0.7)) {
					atlas.TexWidth = 1024
				} else {
					atlas.TexWidth = 512
				}
			}
		}
	}

	// 5. Start packing
	// Pack our extra data rectangles first, so it will be on the upper-left corner of our texture (UV will have small values).
	const TEX_HEIGHT_MAX int = 1024 * 32
	var spc stbtt.PackContext
	stbtt.PackBegin(&spc, nil, atlas.TexWidth, TEX_HEIGHT_MAX, 0, atlas.TexGlyphPadding, nil)
	ImFontAtlasBuildPackCustomRects(atlas, spc.PackInfo)

	// 6. Pack each source font. No rendering yet, we are working with rectangles in an infinitely tall texture at this point.
	for src_i := range src_tmp_array {
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		if src_tmp.GlyphsCount == 0 {
			continue
		}

		stbrp.PackRects(spc.PackInfo.(*stbrp.Context), src_tmp.Rects, src_tmp.GlyphsCount)

		// Extend texture height and mark missing glyphs as non-packed so we won't render them.
		// FIXME: We are not handling packing failure here (would happen if we got off TEX_HEIGHT_MAX or if a single if larger than TexWidth?)
		for glyph_i := int(0); glyph_i < src_tmp.GlyphsCount; glyph_i++ {
			if src_tmp.Rects[glyph_i].WasPacked != 0 {
				atlas.TexHeight = ImMaxInt(atlas.TexHeight, int(src_tmp.Rects[glyph_i].Y)+int(src_tmp.Rects[glyph_i].H))
			}
		}
	}

	// 7. Allocate texture
	if atlas.Flags&ImFontAtlasFlags_NoPowerOfTwoHeight != 0 {
		atlas.TexHeight = (atlas.TexHeight + 1)
	} else {
		atlas.TexHeight = ImUpperPowerOfTwo(atlas.TexHeight)
	}
	atlas.TexUvScale = ImVec2{1.0 / float(atlas.TexWidth), 1.0 / float(atlas.TexHeight)}
	atlas.TexPixelsAlpha8 = make([]byte, atlas.TexWidth*atlas.TexHeight)
	spc.Pixels = atlas.TexPixelsAlpha8
	spc.Height = atlas.TexHeight

	// 8. Render/rasterize font characters into the texture
	for src_i := range src_tmp_array {
		var cfg *ImFontConfig = &atlas.ConfigData[src_i]
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		if src_tmp.GlyphsCount == 0 {
			continue
		}

		stbtt.PackFontRangesRenderIntoRects(&spc, &src_tmp.FontInfo, src_tmp.PackRange[:], 1, src_tmp.Rects)

		// Apply multiply operator
		if cfg.RasterizerMultiply != 1.0 {
			var multiply_table [256]byte
			ImFontAtlasBuildMultiplyCalcLookupTable(multiply_table[:], cfg.RasterizerMultiply)
			var r []stbrp.Rect = src_tmp.Rects
			for glyph_i := int(0); glyph_i < src_tmp.GlyphsCount; glyph_i, r = glyph_i+1, r[1:] {
				if r[0].WasPacked != 0 {
					ImFontAtlasBuildMultiplyRectAlpha8(multiply_table[:], atlas.TexPixelsAlpha8, int(r[0].X), int(r[0].Y), int(r[0].W), int(r[0].H), atlas.TexWidth*1)
				}
			}
		}
		src_tmp.Rects = nil
	}

	// End packing
	buf_rects = buf_rects[:0]

	// 9. Setup ImFont and glyphs for runtime
	for src_i := range src_tmp_array {
		var src_tmp *ImFontBuildSrcData = &src_tmp_array[src_i]
		if src_tmp.GlyphsCount == 0 {
			continue
		}

		// When merging fonts with MergeMode=true:
		// - We can have multiple input fonts writing into a same destination font.
		// - dst_font.ConfigData is != from cfg which is our source configuration.
		var cfg *ImFontConfig = &atlas.ConfigData[src_i]
		var dst_font *ImFont = cfg.DstFont

		var font_scale float = stbtt.ScaleForPixelHeight(&src_tmp.FontInfo, cfg.SizePixels)
		var unscaled_ascent, unscaled_descent, unscaled_line_gap int
		stbtt.GetFontVMetrics(&src_tmp.FontInfo, &unscaled_ascent, &unscaled_descent, &unscaled_line_gap)

		var dir float
		if unscaled_ascent > 0.0 {
			dir = +1
		} else {
			dir = -1
		}

		var ascent float = ImFloor(float(unscaled_ascent)*font_scale + dir)
		var descent float = ImFloor(float(unscaled_descent)*font_scale + dir)
		ImFontAtlasBuildSetupFont(atlas, dst_font, cfg, ascent, descent)
		var font_off_x float = cfg.GlyphOffset.x
		var font_off_y float = cfg.GlyphOffset.y + IM_ROUND(dst_font.Ascent)

		for glyph_i := int(0); glyph_i < src_tmp.GlyphsCount; glyph_i++ {
			// Register glyph
			var codepoint int = src_tmp.GlyphsList[glyph_i]
			var pc *stbtt.PackedChar = &src_tmp.PackedChars[glyph_i]
			var q stbtt.AlignedQuad
			var unused_x float = 0.0
			var unused_y float = 0.0
			stbtt.GetPackedQuad(src_tmp.PackedChars, atlas.TexWidth, atlas.TexHeight, glyph_i, &unused_x, &unused_y, &q, 0)
			dst_font.AddGlyph(cfg, (ImWchar)(codepoint), q.X0+font_off_x, q.Y0+font_off_y, q.X1+font_off_x, q.Y1+font_off_y, q.S0, q.T0, q.S1, q.T1, pc.AdvanceX)
		}
	}

	ImFontAtlasBuildFinish(atlas)
	return true
}
