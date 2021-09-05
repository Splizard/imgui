//
// Finding The Right Font
//
// You should really just solve this offline, keep your own tables
// of what font is what, and don't try to get it out of the .ttf file.
// That's because getting it out of the .ttf file is really hard, because
// the names in the file can appear in many possible encodings, in many
// possible languages, and e.g. if you need a case-insensitive comparison,
// the details of that depend on the encoding & language in a complex way
// (actually underspecified in truetype, but also gigantic).
//
// But you can use the provided functions in two possible ways:
//     stbtt.FindMatchingFont() will use *case-sensitive* comparisons on
//             unicode-encoded names to try to find the font you want;
//             you can run this before calling stbtt_InitFont()
//
//     stbtt.GetFontNameString() lets you get any of the various strings
//             from the file yourself and do your own comparisons on them.
//             You have to have called stbtt_InitFont() first.
//
package stbtt

import (
	"math"

	"github.com/splizard/imgui/stb/stbrp"
)

const (
	vmove = iota
	vline
	vcurve
	vcubic
)

const (
	maxOversample     = 8
	rasterizerVersion = 2
)

const (
	MacstyleDontCare   = 0
	MacstyleBold       = 1
	MacstyleItalic     = 2
	MacstyleUnderscore = 4
	MacstyleNone       = 8
)

// some of the values for the IDs are below; for more see the truetype spec:
//     http://developer.apple.com/textfonts/TTRefMan/RM06/Chap6name.html
//     http://www.microsoft.com/typography/otspec/name.htm
const (
	MSLangEnglish  = 0x0409
	MSLangChinese  = 0x0804
	MSLangDutch    = 0x0413
	MSLangFrench   = 0x040c
	MSLangGerman   = 0x0407
	MSLangHebrew   = 0x040d
	MSLangItalian  = 0x0410
	MSLangJapanese = 0x0411
	MSLangKorean   = 0x0412
	MSLangRussian  = 0x0419
	MSLangSpanish  = 0x0409
	MSLangSwedish  = 0x041D

	PlatformIDUnicode = 0
	PlatformIDMac     = 1
	PlatformIDISO     = 2

	PlatformIDMicrosoft     = 3
	UnicodeEIDUnicode10     = 0
	UnicodeEIDUnicode11     = 1
	UnicodeEIDISO10646      = 2
	UnicodeEIDUnicode20BMP  = 3
	UnicodeEIDUnicode20Full = 4

	MSEIDSymbol      = 0
	MSEIDUnicodeBMP  = 1
	MSEIDShiftJIS    = 2
	MSEIDUnicodeFull = 10

	MacEIDRoman              = 0
	MacEIDJapanese           = 1
	MacEIDChineseTraditional = 2
	MacEIDKorean             = 3
	MacEIDArabic             = 4
	MacEIDHebrew             = 5
	MacEIDGreek              = 6
	MacEIDRussian            = 7
)

type Vertex struct {
	x, y, cx, cy, cx1, cy1 int16
	typ, padding           byte
}

type point struct {
	x, y float32
}

type bitmap struct {
	w, h, stride int32
	pixels       []byte
}

type BakedChar struct {
	X0, Y0, X1, Y1       uint16
	XOff, YOff, XAdvance float32
}

type AlignedQuad struct {
	X0, Y0, S0, T0 float32 // top-left
	X1, Y1, S1, T1 float32 // bottom-right
}

type PackedChar struct {
	X0, Y0, X1, Y1       uint16 // coordinates of bbox in bitmap
	XOff, YOff, XAdvance float32
	XOff2, YOff2         float32
}

type PackRange struct {
	FontSize                     float32
	FirstUnicodeCodepointInRange int32   // if non-zero, then the chars are continuous, and this is the first codepoint
	ArrayOfUnicodeCodepoints     []int32 // if non-zero, then this is an array of unicode codepoints

	ChardataForRange []PackedChar // output

	h_oversample, v_oversample byte
}

// this is an opaque structure that you shouldn't mess with which holds
// all the context needed from PackBegin to PackEnd.
type PackContext struct {
	user_allocator_context     interface{}
	pack_info                  interface{}
	width                      int32
	height                     int32
	stride_in_bytes            int32
	padding                    int32
	skip_missing               int32
	h_oversample, v_oversample uint32
	pixels                     []byte
	nodes                      interface{}
}

// The following structure is defined publicly so you can declare one on
// the stack or as a global or etc, but you should treat it as opaque.
type FontInfo struct {
	userdata  interface{}
	data      []byte
	fontstart int32 // offset of start of font

	numGlyphs int32 // number of glyphs, needed for range checking

	loca, head, glyf, hhea, hmtx, kern, gpos, svg int32 // table locations as offset from start of .ttf
	index_map                                     int32 // a cmap mapping for our chosen character encoding
	indexToLocFormat                              int32 // format needed to map from glyph index to glyph

	cff         buf // cff font data
	charstrings buf // the charstring index
	gsubrs      buf // global charstring subroutines index
	subrs       buf // private charstring subroutines index
	fontdicts   buf // array of font dicts
	fdselect    buf // map from glyph to fontdict
}

type KerningEntry struct {
	glyph1, glyph2 int32 // use stbtt_FindGlyphIndex
	advance        int32
}

// if result is positive, the first unused row of the bitmap
// if result is negative, returns the negative of the number of characters that fit
// if result is 0, no characters fit and no rows were used
// This uses a very crappy packing.
func BakeFontBitmap(data []byte, offset int32, pixel_height float32, pw, ph int32, chars []BakedChar) (pixels []byte, result int32) {
	panic("not implemented")
}

// Call GetBakedQuad with char_index = 'character - first_char', and it
// creates the quad you need to draw and advances the current position.
//
// The coordinate system used assumes y increases downwards.
//
// Characters will extend both above and below the current position;
// see discussion of "BASELINE" above.
//
// It's inefficient; you might want to c&p it and optimize it.
func GetBakedQuad(chardata []BakedChar, pw, ph, index int32, xpos, ypos *float32, opengl_fillrule int32) AlignedQuad {
	panic("not implemented")
}

// Query the font vertical metrics without having to create a font first.
func GetScaledFontVMetrics(fontdata []byte, index int32, size float32) (ascent, descent, lineGap float32) {
	panic("not implemented")
}

// Initializes a packing context stored in the passed-in stbtt_pack_context.
// Future calls using this context will pack characters into the bitmap passed
// in here: a 1-channel bitmap that is width * height. stride_in_bytes is
// the distance from one row to the next (or 0 to mean they are packed tightly
// together). "padding" is the amount of padding to leave between each
// character (normally you want '1' for bitmaps you'll use as textures with
// bilinear filtering).
//
// Returns 0 on failure, 1 on success.
func PackBegin(spc *PackContext, pixels []byte, width, height, stride_in_bytes, padding int32, alloc_context interface{}) int32 {
	panic("not implemented")
}

// Cleans up the packing context and frees all memory.
func PackEnd(spc *PackContext) {}

// Creates character bitmaps from the font_index'th font found in fontdata (use
// font_index=0 if you don't know what that is). It creates num_chars_in_range
// bitmaps for characters with unicode values starting at first_unicode_char_in_range
// and increasing. Data for how to render them is stored in chardata_for_range;
// pass these to stbtt_GetPackedQuad to get back renderable quads.
//
// font_size is the full height of the character from ascender to descender,
// as computed by stbtt_ScaleForPixelHeight. To use a point size as computed
// by stbtt_ScaleForMappingEmToPixels, wrap the point size in STBTT_POINT_SIZE()
// and pass that result as 'font_size':
//       ...,                  20 , ... // font max minus min y is 20 pixels tall
//       ..., STBTT_POINT_SIZE(20), ... // 'M' is 20 pixels tall
func PackFontRange(spc *PackContext, fontdata []byte, fontindex int32, fontsize float32,
	first_unicode_char_in_range, num_chars_in_range int32, chardata_for_range []PackedChar) int32 {

	panic("not implemented")
}

// Creates character bitmaps from multiple ranges of characters stored in
// ranges. This will usually create a better-packed bitmap than multiple
// calls to stbtt_PackFontRange. Note that you can call this multiple
// times within a single PackBegin/PackEnd.
func PackFontRanges(spc *PackContext, fontdata []byte, font_index int32, ranges []PackRange) int32 {
	panic("not implemented")
}

// Oversampling a font increases the quality by allowing higher-quality subpixel
// positioning, and is especially valuable at smaller text sizes.
//
// This function sets the amount of oversampling for all following calls to
// stbtt_PackFontRange(s) or stbtt_PackFontRangesGatherRects for a given
// pack context. The default (no oversampling) is achieved by h_oversample=1
// and v_oversample=1. The total number of pixels required is
// h_oversample*v_oversample larger than the default; for example, 2x2
// oversampling requires 4x the storage of 1x1. For best results, render
// oversampled textures with bilinear filtering. Look at the readme in
// stb/tests/oversample for information about oversampled fonts
//
// To use with PackFontRangesGather etc., you must set it before calls
// call to PackFontRangesGatherRects.
func PackSetOversampling(spc *PackContext, h_oversample, v_oversample uint32) {
	panic("not implemented")
}

// If skip != 0, this tells stb_truetype to skip any codepoints for which
// there is no corresponding glyph. If skip=0, which is the default, then
// codepoints without a glyph recived the font's "missing character" glyph,
// typically an empty box by convention.
func PackSetSkipMissingCodepoints(spc *PackContext, skip int32) {
	panic("not implemented")
}

func GetPackedQuad(chardata []PackedChar, pw, ph int32, char_index int32, xpos, ypos *float32, align_to_integer int32) AlignedQuad {
	panic("not implemented")
}

// Calling these functions in sequence is roughly equivalent to calling
// stbtt_PackFontRanges(). If you more control over the packing of multiple
// fonts, or if you want to pack custom data into a font texture, take a look
// at the source to of stbtt_PackFontRanges() and create a custom version
// using these functions, e.g. call GatherRects multiple times,
// building up a single array of rects, then call PackRects once,
// then call RenderIntoRects repeatedly. This may result in a
// better packing than calling PackFontRanges multiple times
// (or it may not).

func PackFontRangesGatherRects(spc *PackContext, info *FontInfo, ranges []PackRange, rects []stbrp.Rect) int32 {
	panic("not implemented")
}
func PackFontRangesPackRects(spc *PackContext, rects []stbrp.Rect) {
	panic("not implemented")
}

func PackFontRangesRenderIntoRects(spc *PackContext, info *FontInfo, ranges []PackRange, rects []stbrp.Rect) {
	panic("not implemented")
}

// This function will determine the number of fonts in a font file.  TrueType
// collection (.ttc) files may contain multiple fonts, while TrueType font
// (.ttf) files only contain one font. The number of fonts can be used for
// indexing with the previous function where the index is between zero and one
// less than the total fonts. If an error occurs, -1 is returned.
func GetNumberOfFonts(data []byte) (int32, error) {
	panic("not implemented")
}

// Each .ttf/.ttc file may have more than one font. Each font has a sequential
// index number starting from 0. Call this function to get the font offset for
// a given index; it returns -1 if the index is out of range. A regular .ttf
// file will only define one font and it always be at offset 0, so it will
// return '0' for index 0, and -1 for all other indices.
func GetFontOffsetForIndex(data []byte, index int) int32 {
	panic("not implemented")
}

// Given an offset into the file that defines a font, this function builds
// the necessary cached info for the rest of the system. You must allocate
// the stbtt_fontinfo yourself, and stbtt_InitFont will fill it out. You don't
// need to do anything special to free it, because the contents are pure
// value data with no additional data structures. Returns 0 on failure.
func InitFont(info *FontInfo, data []byte, offset int) int32 {
	panic("not implemented")
}

// If you're going to perform multiple operations on the same character
// and you want a speed-up, call this function with the character you're
// going to process, then use glyph-based functions instead of the
// codepoint-based functions.
// Returns 0 if the character codepoint is not defined in the font.
func FindGlyphIndex(info *FontInfo, unicode_codepoint rune) int32 {
	var data = info.data
	var index_map = info.index_map

	var format = ttUSHORT(data[index_map:])
	switch format {
	case 0: // apple byte encoding
		var bytes = ttUSHORT(data[index_map+2:])
		if uint16(unicode_codepoint) < bytes-6 {
			return int32(ttBYTE(data[index_map+6+unicode_codepoint]))
		}
		return 0
	case 6: // apple 16-bit encoding
		var first = ttUSHORT(data[index_map+6:])
		var count = ttUSHORT(data[index_map+8:])
		if uint16(unicode_codepoint) >= first && uint16(unicode_codepoint) < first+count {
			return int32(ttUSHORT(data[index_map+10+2*(unicode_codepoint-int32(first)):]))
		}
		return 0
	case 2:
		panic("not implemented (TODO: high-byte mapping for japanese/chinese/korean)")
	case 4: // standard mapping for windows fonts: binary search collection of ranges
		var segcount = ttUSHORT(data[index_map+6:]) >> 1
		var searchRange = ttUSHORT(data[index_map+8:]) >> 1
		var entrySelector = ttUSHORT(data[index_map+10:])
		var rangeShift = ttUSHORT(data[index_map+12:]) >> 1

		// do a binary search of the segments
		var endCount = index_map + 14
		var search = endCount

		if unicode_codepoint > 0xffff {
			return 0
		}

		// they lie from endCount .. endCount + segCount
		// but searchRange is the nearest power of two, so...
		if unicode_codepoint >= int32(ttUSHORT(data[search+int32(rangeShift)*2:])) {
			search += int32(rangeShift) * 2
		}

		// now decrement to bias correctly to find smallest
		search -= 2
		for entrySelector != 0 {
			var end uint16
			searchRange >>= 1
			end = ttUSHORT(data[search+int32(searchRange)*2:])
			if unicode_codepoint > int32(end) {
				search += int32(searchRange * 2)
			}
			entrySelector--
		}
		search += 2

		{
			var offset, start, last uint16
			var item uint16 = uint16((search - endCount) >> 1)

			start = ttUSHORT(data[index_map+14+int32(segcount)*2+2+2*int32(item):])
			last = ttUSHORT(data[int32(endCount)+2+2*int32(item):])
			if uint16(unicode_codepoint) < start || unicode_codepoint > int32(last) {
				return 0
			}

			offset = ttUSHORT(data[index_map+14+int32(segcount)*6+2+2*int32(item):])
			if offset == 0 {
				return int32(unicode_codepoint + int32(ttSHORT(data[index_map+14+int32(segcount)*4+2+2*int32(item):])))
			}

			return int32(ttUSHORT(data[int32(offset)+(unicode_codepoint-int32(start))*2+index_map+14+int32(segcount)*6+2+2*int32(item):]))
		}
	case 12, 13: // standard mapping for macintosh fonts, java: high-byte mapping zero
		var ngroups = ttULONG(data[index_map+12:])
		var low, high int32
		low = 0
		high = int32(ngroups)
		// Binary search the right group.
		for low < high {
			var mid = (low + high) >> 1
			var start_char = ttULONG(data[index_map+16+mid*12:])
			var end_char = ttULONG(data[index_map+16+mid*12+4:])
			if uint32(unicode_codepoint) < start_char {
				high = mid
			} else if uint32(unicode_codepoint) > end_char {
				low = mid + 1
			} else {
				var start_glyph = ttULONG(data[index_map+16+mid*12+8:])
				if format == 12 {
					return int32(start_glyph + uint32(unicode_codepoint) - start_char)
				} else { // format == 13
					return int32(start_glyph)
				}
			}
		}
		return 0 //not found
	}
	panic("unreachable")
}

// computes a scale factor to produce a font whose "height" is 'pixels' tall.
// Height is measured as the distance from the highest ascender to the lowest
// descender; in other words, it's equivalent to calling stbtt_GetFontVMetrics
// and computing:
//       scale = pixels / (ascent - descent)
// so if you prefer to measure height by the ascent only, use a similar calculation.
func ScaleForPixelHeight(info *FontInfo, pixels float32) float32 {
	var fheight = ttSHORT(info.data[info.hhea+4:]) - ttSHORT(info.data[info.hhea+6:])
	return float32(pixels) / float32(fheight)
}

// computes a scale factor to produce a font whose EM size is mapped to
// 'pixels' tall. This is probably what traditional APIs compute, but
// I'm not positive.
func ScaleForMappingEmToPixels(info *FontInfo, pixels float32) float32 {
	var unitsPerEm = ttUSHORT(info.data[info.head+18:])
	return pixels / float32(unitsPerEm)
}

// ascent is the coordinate above the baseline the font extends; descent
// is the coordinate below the baseline the font extends (i.e. it is typically negative)
// lineGap is the spacing between one row's descent and the next row's ascent...
// so you should advance the vertical position by "*ascent - *descent + *lineGap"
//   these are expressed in unscaled coordinates, so you must multiply by
//   the scale factor for a given size
func GetFontVMetrics(info *FontInfo) (accent, descent, lineGap int32) {
	return int32(ttSHORT(info.data[info.hhea+4:])),
		int32(ttSHORT(info.data[info.hhea+6:])),
		int32(ttSHORT(info.data[info.hhea+8:]))
}

// analogous to GetFontVMetrics, but returns the "typographic" values from the OS/2
// table (specific to MS/Windows TTF files).
//
// Returns 1 on success (table present), 0 on failure.
func GetFontVMetricsOS2(info *FontInfo) (ok bool, accent, descent, lineGap int32) {
	var tab = findTable(info.data, uint32(info.fontstart), "OS/2")
	if tab == 0 {
		return false, 0, 0, 0
	}
	return true,
		int32(ttSHORT(info.data[tab+68:])),
		int32(ttSHORT(info.data[tab+70:])),
		int32(ttSHORT(info.data[tab+72:]))

}

// the bounding box around all possible characters
func GetFontBoundingBox(info *FontInfo) (x0, y0, x1, y1 int32) {
	return int32(ttSHORT(info.data[info.head+36:])),
		int32(ttSHORT(info.data[info.head+38:])),
		int32(ttSHORT(info.data[info.head+40:])),
		int32(ttSHORT(info.data[info.head+42:]))
}

// leftSideBearing is the offset from the current horizontal position to the left edge of the character
// advanceWidth is the offset from the current horizontal position to the next horizontal position
//   these are expressed in unscaled coordinates
func GetCodepointHMetrics(info *FontInfo, codepoint rune) (advanceWidth, leftSideBearing int32) {
	return GetGlyphHMetrics(info, FindGlyphIndex(info, codepoint))
}

// an additional amount to add to the 'advance' value between ch1 and ch2
func GetCodepointKernAdvance(info *FontInfo, ch1, ch2 int32) int32 {
	if info.kern == 0 && info.gpos == 0 {
		return 0 // if no kerning table, don't waste time looking up both codepoint->glyphs
	}
	return GetGlyphKernAdvance(info, FindGlyphIndex(info, ch1), FindGlyphIndex(info, ch2))
}

// Gets the bounding box of the visible part of the glyph, in unscaled coordinates
func GetCodepointBox(info *FontInfo, codepoint rune) (ok bool, x0, y0, x1, y1 int32) {
	return GetGlyphBox(info, FindGlyphIndex(info, codepoint))
}

func GetGlyphHMetrics(info *FontInfo, glyphIndex int32) (advanceWidth, leftSideBearing int32) {
	var numOfLongHorMetrics = int32(ttUSHORT(info.data[info.hhea+34:]))
	if glyphIndex < numOfLongHorMetrics {
		return int32(ttSHORT(info.data[info.hmtx+4*glyphIndex:])),
			int32(ttSHORT(info.data[info.hmtx+4*glyphIndex+2:]))
	} else {
		return int32(ttSHORT(info.data[info.hmtx+4*(numOfLongHorMetrics-1):])),
			int32(ttSHORT(info.data[info.hmtx+4*numOfLongHorMetrics+2*(glyphIndex-numOfLongHorMetrics):]))
	}
}

func GetGlyphKernAdvance(info *FontInfo, glyph1, glyph2 int32) int32 {
	var xAdvance int32 = 0
	if info.gpos != 0 {
		xAdvance += getGlyphGPOSInfoAdvance(info, glyph1, glyph2)
	} else if info.kern != 0 {
		xAdvance += getGlyphKernInfoAdvance(info, glyph1, glyph2)
	}
	return xAdvance
}

func GetGlyphBox(info *FontInfo, glyphIndex int32) (ok bool, x0, y0, x1, y1 int32) {
	if len(info.cff.data) != 0 {
		x0, y0, x1, y1, err := getGlyphInfoT2(info, glyphIndex)
		return err == nil, x0, y0, x1, y1
	} else {
		var g = getGlyfOffset(info, glyphIndex)
		if g < 0 {
			return false, 0, 0, 0, 0
		}

		x0 = int32(ttSHORT(info.data[g+2:]))
		y0 = int32(ttSHORT(info.data[g+4:]))
		x1 = int32(ttSHORT(info.data[g+6:]))
		y1 = int32(ttSHORT(info.data[g+8:]))
	}
	return true, x0, y0, x1, y1
}

func GetKerningTableLength(info *FontInfo) int32 {
	var data = info.data[info.kern:]

	// we only look at the first table. it must be 'horizontal' and format 0.
	if info.kern == 0 {
		return 0
	}
	if ttUSHORT(data[2:]) < 1 { // number of tables, need at least 1
		return 0
	}
	if ttUSHORT(data[8:]) != 1 { // horizontal flag must be set in format
		return 0
	}

	return int32(ttUSHORT(data[10:]))
}

// Retrieves a complete list of all of the kerning pairs provided by the font
// stbtt_GetKerningTable never writes more than table_length entries and returns how many entries it did write.
// The table will be sorted by (a.glyph1 == b.glyph1)?(a.glyph2 < b.glyph2):(a.glyph1 < b.glyph1)
func GetKerningTable(info *FontInfo) (table []KerningEntry) {
	var data = info.data[info.kern:]

	if info.kern == 0 {
		return nil
	}
	if ttUSHORT(data[2:]) < 1 { // number of tables, need at least 1
		return nil
	}
	if ttUSHORT(data[8:]) != 1 { // horizontal flag must be set in format
		return nil
	}

	var length = int32(ttUSHORT(data[10:]))
	table = make([]KerningEntry, length)

	var i int32
	for i = 0; i < length; i++ {
		table[i].glyph1 = int32(ttUSHORT(data[18+(i*6):]))
		table[i].glyph2 = int32(ttUSHORT(data[20+(i*6):]))
		table[i].advance = int32(ttSHORT(data[22+(i*6):]))
	}

	return table
}

// returns true if nothing is drawn for this glyph
func IsGlyphEmpty(info *FontInfo, glyph_index int32) bool {
	var numberOfContours int16
	if len(info.cff.data) != 0 {
		_, _, _, _, err := getGlyphInfoT2(info, glyph_index)
		return err != nil
	} else {
		g := getGlyfOffset(info, glyph_index)
		if g < 0 {
			return true
		}
		numberOfContours = int16(ttSHORT(info.data[g:]))
	}
	return numberOfContours == 0
}

func GetCodepointShape(info *FontInfo, codepoint rune) (vertices []Vertex) {
	return GetGlyphShape(info, FindGlyphIndex(info, codepoint))
}

// returns # of vertices and fills *vertices with the pointer to them
//   these are expressed in "unscaled" coordinates
//
// The shape is a series of contours. Each one starts with
// a STBTT_moveto, then consists of a series of mixed
// STBTT_lineto and STBTT_curveto segments. A lineto
// draws a line from previous endpoint to its x,y; a curveto
// draws a quadratic bezier from previous endpoint to
// its x,y, using cx,cy as the bezier control point.
func GetGlyphShape(info *FontInfo, glyph_index int32) (vertices []Vertex) {
	if len(info.cff.data) == 0 {
		return getGlyphShapeTT(info, glyph_index)
	}
	return getGlyphShapeT2(info, glyph_index)
}

func FindSVGDoc(info *FontInfo, gl int32) []byte {
	var data = info.data
	var svg_doc_list = data[info.svg:]

	var numEntries = int(ttUSHORT(svg_doc_list))
	var svg_docs = svg_doc_list[2:]

	for i := 0; i < numEntries; i++ {
		var svg_doc = svg_docs[12*i:]
		if (gl >= int32(ttUSHORT(svg_doc))) && (gl <= int32(ttUSHORT(svg_doc[2:]))) {
			return svg_doc
		}
	}
	return nil
}

func GetCodepointSVG(info *FontInfo, unicode_codepoint rune) (svg []byte) {
	return GetGlyphSVG(info, FindGlyphIndex(info, unicode_codepoint))
}

// fills svg with the character's SVG data.
// returns data size or 0 if SVG not found.
func GetGlyphSVG(info *FontInfo, gl int32) (svg []byte) {
	var data = info.data
	if len(info.data) == 0 {
		return nil
	}

	svg_doc := FindSVGDoc(info, gl)
	if svg_doc == nil {
		return nil
	}

	return data[info.svg+int32(ttULONG(svg_doc[4:])) : info.svg+int32(ttULONG(svg_doc[8:]))]
}

// allocates a large-enough single-channel 8bpp bitmap and renders the
// specified character/glyph at the specified scale into it, with
// antialiasing. 0 is no coverage (transparent), 255 is fully covered (opaque).
// *width & *height are filled out with the width & height of the bitmap,
// which is stored left-to-right, top-to-bottom.
//
// xoff/yoff are the offset it pixel space from the glyph origin to the top-left of the bitmap
func GetCodepointBitmap(info *FontInfo, scale_x, scale_y float32, codepoint rune) (bitmap []byte, width, height, xoff, yoff int32) {
	panic("not implemented")
}

// the same as GetCodepointBitmap, but you can specify a subpixel
// shift for the character
func GetCodepointBitmapSubpixel(info *FontInfo, scale_x, scale_y, shift_x, shift_y float32, codepoint rune) (bitmap []byte, width, height, xoff, yoff int32) {
	panic("not implemented")
}

// the same as stbtt_GetCodepointBitmap, but you pass in storage for the bitmap
// in the form of 'output', with row spacing of 'out_stride' bytes. the bitmap
// is clipped to out_w/out_h bytes. Call stbtt_GetCodepointBitmapBox to get the
// width and height and positioning info for it first.
func MakeCodepointBitmap(info *FontInfo, output []byte, width, height, stride int32, scale_x, scale_y float32, codepoint rune) {
	panic("not implemented")
}

// same as stbtt_MakeCodepointBitmap, but you can specify a subpixel
// shift for the character
func MakeCodepointBitmapSubpixel(info *FontInfo, output []byte, width, height, stride int32, scale_x, scale_y, shift_x, shift_y float32, codepoint rune) {
	panic("not implemented")
}

// same as stbtt_MakeCodepointBitmapSubpixel, but prefiltering
// is performed (see stbtt_PackSetOversampling)
func MakeCodepointBitmapSubpixelPrefilter(info *FontInfo, output []byte, width, height, stride int32, scale_x, scale_y, shift_x, shift_y float32, oversample_x, oversample_y int, sub_x, sub_y []float32, codepoint rune) {
	panic("not implemented")
}

// get the bbox of the bitmap centered around the glyph origin; so the
// bitmap width is ix1-ix0, height is iy1-iy0, and location to place
// the bitmap top left is (leftSideBearing*scale,iy0).
// (Note that the bitmap uses y-increases-down, but the shape uses
// y-increases-up, so CodepointBitmapBox and CodepointBox are inverted.)
func GetCodepointBitmapBox(info *FontInfo, codepoint rune, scale_x, scale_y float32) (x0, y0, x1, y1 int32) {
	return GetCodepointBitmapBoxSubpixel(info, codepoint, scale_x, scale_y, 0.0, 0.0)
}

// same as stbtt_GetCodepointBitmapBox, but you can specify a subpixel
// shift for the character
func GetCodepointBitmapBoxSubpixel(info *FontInfo, codepoint rune, scale_x, scale_y, shift_x, shift_y float32) (x0, y0, x1, y1 int32) {
	return GetGlyphBitmapBoxSubpixel(info, FindGlyphIndex(info, codepoint), scale_x, scale_y, shift_x, shift_y)
}

func GetGlyphBitmap(info *FontInfo, scale_x, scale_y float32, glyph int) (bitmap []byte, width, height, xoff, yoff int32) {
	panic("not implemented")
}

func GetGlyphBitmapSubpixel(info *FontInfo, scale_x, scale_y, shift_x, shift_y float32, glyph int) (bitmap []byte, width, height, xoff, yoff int32) {
	panic("not implemented")
}

func MakeGlyphBitmap(info *FontInfo, output []byte, width, height, stride int32, scale_x, scale_y float32, glyph int) {
	panic("not implemented")
}

func MakeGlyphBitmapSubpixel(info *FontInfo, output []byte, width, height, stride int32, scale_x, scale_y, shift_x, shift_y float32, glyph int) {
	panic("not implemented")
}

func MakeGlyphBitmapSubpixelPrefilter(info *FontInfo, output []byte, width, height, stride int32, scale_x, scale_y, shift_x, shift_y float32, oversample_x, oversample_y int, sub_x, sub_y []float32, glyph int) {
	panic("not implemented")
}

func GetGlyphBitmapBox(info *FontInfo, glyph int32, scale_x, scale_y float32) (x0, y0, x1, y1 int32) {
	return GetGlyphBitmapBoxSubpixel(info, glyph, scale_x, scale_y, 0.0, 0.0)
}

func GetGlyphBitmapBoxSubpixel(info *FontInfo, glyph int32, scale_x, scale_y, shift_x, shift_y float32) (x0, y0, x1, y1 int32) {
	ok, x0, y0, x1, y1 := GetGlyphBox(info, glyph)
	if !ok {
		return 0, 0, 0, 0
	}
	return int32(math.Floor(float64(x0)*float64(scale_x) + float64(shift_x))),
		int32(math.Floor(-float64(y1)*float64(scale_y) + float64(shift_y))),
		int32(math.Ceil(float64(x1)*float64(scale_x) + float64(shift_x))),
		int32(math.Ceil(-float64(y0)*float64(scale_y) + float64(shift_y)))
}

// rasterize a shape with quadratic beziers into a bitmap
func Rasterize(result *bitmap, // 1-channel bitmap to draw into
	flatness_in_pixels float32, // allowable error of curve in pixels
	vertices []Vertex, // array of vertices defining shape
	scale_x, scale_y, // scale applied to input vertices
	shift_x, shift_y float32, // translation applied to input vertices
	x_off, y_off int32, // another translation applied to input
	invert bool, // if non-zero, vertically flip shape
) {
	var scale float32
	if scale_x > scale_y {
		scale = scale_y
	} else {
		scale = scale_x
	}
	points, lengths := flattenCurves(vertices, flatness_in_pixels/scale)
	if lengths != nil {
		rasterize(result, points, lengths, scale_x, scale_y, shift_x, shift_y, x_off, y_off, invert)
	}
}

func GetGlyphSDF(info *FontInfo, scale float32, glyph int, padding int, onedge_value byte, pixel_dist_scale float32) (result []byte, width, height, xoff, yoff int) {
	panic("not implemented")
}

// These functions compute a discretized SDF field for a single character, suitable for storing
// in a single-channel texture, sampling with bilinear filtering, and testing against
// larger than some threshold to produce scalable fonts.
//        info              --  the font
//        scale             --  controls the size of the resulting SDF bitmap, same as it would be creating a regular bitmap
//        glyph/codepoint   --  the character to generate the SDF for
//        padding           --  extra "pixels" around the character which are filled with the distance to the character (not 0),
//                                 which allows effects like bit outlines
//        onedge_value      --  value 0-255 to test the SDF against to reconstruct the character (i.e. the isocontour of the character)
//        pixel_dist_scale  --  what value the SDF should increase by when moving one SDF "pixel" away from the edge (on the 0..255 scale)
//                                 if positive, > onedge_value is inside; if negative, < onedge_value is inside
//        width,height      --  output height & width of the SDF bitmap (including padding)
//        xoff,yoff         --  output origin of the character
//        return value      --  a 2D array of bytes 0..255, width*height in size
//
// pixel_dist_scale & onedge_value are a scale & bias that allows you to make
// optimal use of the limited 0..255 for your application, trading off precision
// and special effects. SDF values outside the range 0..255 are clamped to 0..255.
//
// Example:
//      scale = stbtt_ScaleForPixelHeight(22)
//      padding = 5
//      onedge_value = 180
//      pixel_dist_scale = 180/5.0 = 36.0
//
//      This will create an SDF bitmap in which the character is about 22 pixels
//      high but the whole bitmap is about 22+5+5=32 pixels high. To produce a filled
//      shape, sample the SDF at each pixel and fill the pixel if the SDF value
//      is greater than or equal to 180/255. (You'll actually want to antialias,
//      which is beyond the scope of this example.) Additionally, you can compute
//      offset outlines (e.g. to stroke the character border inside & outside,
//      or only outside). For example, to fill outside the character up to 3 SDF
//      pixels, you would compare against (180-36.0*3)/255 = 72/255. The above
//      choice of variables maps a range from 5 pixels outside the shape to
//      2 pixels inside the shape to 0..255; this is intended primarily for apply
//      outside effects only (the interior range is needed to allow proper
//      antialiasing of the font at *smaller* sizes)
//
// The function computes the SDF analytically at each SDF pixel, not by e.g.
// building a higher-res bitmap and approximating it. In theory the quality
// should be as high as possible for an SDF of this size & representation, but
// unclear if this is true in practice (perhaps building a higher-res bitmap
// and computing from that can allow drop-out prevention).
//
// The algorithm has not been optimized at all, so expect it to be slow
// if computing lots of characters or very large sizes.
func GetCodepointSDF(info *FontInfo, scale float32, codepoint rune, padding int, onedge_value byte, pixel_dist_scale float32) (result []byte, width, height, xoff, yoff int) {
	panic("not implemented")
}

// returns the offset (not index) of the font that matches, or -1 if none
//   if you use STBTT_MACSTYLE_DONTCARE, use a font name like "Arial Bold".
//   if you use any other flag, use a font name like "Arial"; this checks
//     the 'macStyle' header field; i don't know if fonts set this consistently
func FindMatchingFont(fontdata []byte, name string, flags int) int {
	panic("not implemented")
}

// returns 1/0 whether the first string interpreted as utf8 is identical to
// the second string interpreted as big-endian utf16... useful for strings from next func
func CompareUTF8toUTF16Bigendian(s1 string, s2 string) bool {
	panic("not implemented")
}

// returns the string (which may be big-endian double byte, e.g. for unicode)
// and puts the length in bytes in *length.
//
// some of the values for the IDs are below; for more see the truetype spec:
//     http://developer.apple.com/textfonts/TTRefMan/RM06/Chap6name.html
//     http://www.microsoft.com/typography/otspec/name.htm
func GetFontNameString(font *FontInfo, platformID, encodingID, languageID, nameID int) string {
	panic("not implemented")
}
