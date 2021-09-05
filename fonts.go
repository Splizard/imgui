package imgui

const (
	DrawlistTexLinesWidthMax = 63

	FontAtlasDefaultTexDataWidth  = 108
	FontAtlasDefaultTexDataHeight = 27

	DrawListTexLinesWidthMax = 63
)

const (
	FontAtlasFlagsNone               FontAtlasFlags = 0
	FontAtlasFlagsNoPowerOfTwoHeight FontAtlasFlags = 1 << 0 // Don't round the height to next power of two
	FontAtlasFlagsNoMouseCursors     FontAtlasFlags = 1 << 1 // Don't build software mouse cursors into the atlas (save a little texture memory)
	FontAtlasFlagsNoBakedLines       FontAtlasFlags = 1 << 2 // Don't build thick line textures into the atlas (save a little texture memory). The AntiAliasedLinesUseTex features uses them, otherwise they will be rendered using polygons (more expensive for CPU/GPU).
)

type FontBuilder func(atlas *FontAtlas) bool

type fontAtlasCustomRect struct{}

//FontConfig contains Configuration data when adding a font or merging fonts
type FontConfig struct {
	FontData           []byte  //          // TTF/OTF data
	FontDataSize       int     //          // TTF/OTF data size
	FontNo             int     // true     // Index of font within TTF/OTF file
	SizePixels         float32 //          // Size in pixels for rasterizer (more or less maps to the resulting font height).
	OversampleH        int     // 3        // Rasterize at higher quality for sub-pixel positioning. Note the difference between 2 and 3 is minimal so you can reduce this to 2 to save memory. Read https://github.com/nothings/stb/blob/master/tests/oversample/README.md for details.
	OversampleV        int     // 1        // Rasterize at higher quality for sub-pixel positioning. This is not really useful as we don't use sub-pixel positions on the Y axis.
	PixelSnapH         bool    // false    // Align every glyph to pixel boundary. Useful e.g. if you are merging a non-pixel aligned font with the default font. If enabled, you can set OversampleH/V to 1.
	GlyphExtraSpacing  Vec2    // 0, 0     // Extra spacing (in pixels) between glyphs. Only X axis is supported for now.
	GlyphOffset        Vec2    // 0, 0     // Offset all glyphs from this font input.
	GlyphRanges        string  // NULL     // Pointer to a user-provided list of Unicode range (2 value per range, values are inclusive, zero-terminated list). THE ARRAY DATA NEEDS TO PERSIST AS LONG AS THE FONT IS ALIVE.
	GlyphMinAdvanceX   float32 // 0        // Minimum AdvanceX for glyphs, set Min to align font icons, set both Min/Max to enforce mono-space font
	GlyphMaxAdvanceX   float32 // FLT_MAX  // Maximum AdvanceX for glyphs
	MergeMode          bool    // false    // Merge into previous ImFont, so you can combine multiple inputs font into one ImFont (e.g. ASCII font + icons + Japanese glyphs). You may want to use GlyphOffset.y when merge font of different heights.
	FontBuilderFlags   uint    // 0        // Settings for custom font builder. THIS IS BUILDER IMPLEMENTATION DEPENDENT. Leave as zero if unsure.
	RasterizerMultiply float32 // 1.0f     // Brighten (>1.0f) or darken (<1.0f) font output. Brightening small fonts may be a good workaround to make them more readable.
	EllipsisChar       rune    // -1       // Explicitly specify unicode codepoint of ellipsis character. When fonts are being merged first specified ellipsis will be used.

	// [Internal]
	name    string // Name (strictly to ease debugging)
	dstFont *Font
}

type FontAtlas struct {
	Flags           FontAtlasFlags
	TexID           TextureID
	TexDesiredWidth int
	TexGlyphPadding int
	Locked          bool

	// [Internal]
	// NB: Access texture data via GetTexData*() calls! Which will setup a default font for you.
	texReady           bool
	texPixelsUseColors bool
	texPixelsAlpha8    []byte
	texPixelsRGBA32    []int32
	texWidth           int
	texHeight          int
	texUVScale         Vec2
	texUVWhitePixel    Vec2
	fonts              []*Font
	customRects        []fontAtlasCustomRect
	configData         []FontConfig
	texUVLines         [DrawlistTexLinesWidthMax + 1]Vec4

	// [Internal] Font builder
	fontBuilderIO    FontBuilder
	fontBuilderFlags int

	// [Internal] Packing data
	packIDMouseCursors int
	packIDLines        int
}

func newFontAtlas() *FontAtlas {
	return &FontAtlas{
		TexGlyphPadding:    1,
		packIDMouseCursors: -1,
		packIDLines:        -1,
	}
}

func (f *FontAtlas) AddFontDefault() {
	panic("not implemented")
}

func (f *FontAtlas) Build() bool {
	if f.Locked {
		panic("Cannot modify a locked FontAtlas between NewFrame() and EndFrame/Render()!")
	}

	if len(f.configData) == 0 {
		f.AddFontDefault()
	}

	if f.fontBuilderIO == nil {
		f.fontBuilderIO = fontAtlasBuildWithStbTruetype
	}

	return f.fontBuilderIO(f)
}

func (f *FontAtlas) AddCustomRectRegular(width, height int) int {
	panic("not implemented")
}

func (f *FontAtlas) GetTexDataAsAlpha8() (pixels []byte, w int, h int) {
	// Build atlas on demand
	if f.texPixelsAlpha8 == nil {
		f.Build()
	}
	return f.texPixelsAlpha8, f.texWidth, f.texHeight
}

func (f *FontAtlas) SetTexID(id TextureID) {
	panic("not implemented")
}

type Font struct {
}

func FontAtlasBuildInit(atlas *FontAtlas) {
	// Register texture region for mouse cursors or standard white pixels
	if atlas.packIDMouseCursors < 0 {
		if atlas.Flags&FontAtlasFlagsNoMouseCursors == 0 {
			atlas.packIDMouseCursors = atlas.AddCustomRectRegular(FontAtlasDefaultTexDataWidth*2+1, FontAtlasDefaultTexDataHeight)
		} else {
			atlas.packIDMouseCursors = atlas.AddCustomRectRegular(2, 2)
		}
	}

	// Register texture region for thick lines
	// The +2 here is to give space for the end caps, whilst height +1 is to accommodate the fact we have a zero-width row
	if atlas.packIDLines < 0 {
		if atlas.Flags&FontAtlasFlagsNoBakedLines == 0 {
			atlas.packIDLines = atlas.AddCustomRectRegular(DrawListTexLinesWidthMax+2, DrawListTexLinesWidthMax+1)
		}
	}
}
