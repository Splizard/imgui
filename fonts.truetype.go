package imgui

import (
	"github.com/splizard/imgui/stb/stbrp"
	"github.com/splizard/imgui/stb/stbtt"
)

type fontBuildSrcData struct {
	FontInfo      stbtt.FontInfo
	PackRange     stbtt.PackRange    // Hold the list of codepoints to pack (essentially points to Codepoints.Data)
	Rects         []stbrp.Rect       // Rectangle to pack. We first fill in their size and the packer will give us their position.
	PackedChars   []stbtt.PackedChar // Output glyphs
	SrcRanges     string             // Ranges as requested by user (user is allowed to request too much, e.g. 0x0020..0xFFFF)
	DstIndex      int                // Index into atlas->Fonts[] and dst_tmp_array[]
	GlyphsHighest int                // Highest requested codepoint
	GlyphsCount   int                // Glyph count (excluding missing glyphs and glyphs already set by an earlier source font)
	GlyphsSet     BitVector          // Glyph bit map (random access, 1-bit per codepoint. This will be a maximum of 8KB)
	GlyphsList    []int              // Glyph codepoints list (flattened version of GlyphsMap)
}

// Temporary data for one destination ImFont* (multiple source fonts can be merged into one destination ImFont)
type fontBuildDstData struct {
	SrcCount      int // Number of source fonts targeting this destination font.
	GlyphsHighest int
	GlyphsCount   int
	GlyphsSet     BitVector // This is used to resolve collision when multiple sources are merged into a same destination font.
}

func fontAtlasBuildWithStbTruetype(atlas *FontAtlas) bool {
	if len(atlas.configData) == 0 {
		panic("no config data")
	}

	FontAtlasBuildInit(atlas)

	// Clear atlas
	atlas.TexID = 0
	atlas.texWidth = 0
	atlas.texHeight = 0
	atlas.texUVScale = Vec2{}
	atlas.texUVWhitePixel = Vec2{}

	// Temporary storage for building
	var src_tmp_array = make([]fontBuildSrcData, len(atlas.configData))
	//var dst_tmp_array = make([]fontBuildDstData, len(atlas.fonts))

	// 1. Initialize font loading structure, check font data validity
	for i := range atlas.configData {
		var src_tmp = &src_tmp_array[i]
		var cfg = &atlas.configData[i]

		// Find index from cfg.DstFont (we allow the user to set cfg.DstFont. Also it makes casual debugging nicer than when storing indices)
		src_tmp.DstIndex = -1
		for j := range atlas.fonts {
			if atlas.fonts[j] == cfg.dstFont {
				src_tmp.DstIndex = j
				break
			}
		}
		if src_tmp.DstIndex == -1 {
			panic("dstFont not found")
		}
		font_offset := stbtt.GetFontOffsetForIndex(cfg.FontData, cfg.FontNo)
		if font_offset == -1 {
			panic("FontData is incorrect, or FontNo cannot be found.")
		}
		// Measure highest codepoints
		/*dst_temp := dst_tmp_array[src_tmp.DstIndex]
		src_tmp.SrcRanges = cfg.SrcRanges
		if cfg.GlyphRanges != "" {
			src_tmp.SrcRanges = atlas.getGlyphRangesDefault()
		}*/

		panic("THIS IS WHERE I GOT TO not implemented")
	}

	panic("THIS IS WHERE I GOT TO not implemented")

}
