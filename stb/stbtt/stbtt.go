package stbtt

import (
	"math"
	"unsafe"

	"github.com/Splizard/imgui/golang"
	"github.com/Splizard/imgui/stb/stbrp"
)

//PORTING STATUS = DONE - Quentin Quaadgras

func iszero(x uint) bool {
	return x == 0
}

type double = float64
type int = int32
type uint = uint32
type float = float32
type size_t = uintptr
type char = byte

func isfalse(x int) bool {
	return x == 0
}

func istrue(x int) bool {
	return x != 0
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

type (
	stbtt_uint8  = uint8
	stbtt_int8   = int8
	stbtt_uint16 = uint16
	stbtt_int16  = int16
	stbtt_uint32 = uint32
	stbtt_int32  = int32
)

func STBTT_ifloor(x float32) int {
	return int(math.Floor(float64(x)))
}

func STBTT_iceil(x float32) int {
	return int(math.Ceil(float64(x)))
}

func STBTT_sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func STBTT_pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func STBTT_fmod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

func STBTT_cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func STBTT_acos(x float32) float32 {
	return float32(math.Acos(float64(x)))
}

func STBTT_fabs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func STBTT_assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}

func STBTT_strlen(str string) int {
	return int(len(str))
}

// private structure
type stbtt__buf struct {
	data   []byte
	cursor int
	size   int
}

type stbtt_bakedchar struct {
	x0, y0, x1, y1       uint16 // coordinates of bbox in bitmap
	xoff, yoff, xadvance float
}

// if return is positive, the first unused row of the bitmap
// if return is negative, returns the negative of the number of characters that fit
// if return is 0, no characters fit and no rows were used
// This uses a very crappy packing.
func stbtt_BakeFontBitmap(data []byte, offset int, // font location (use offset=0 for plain .ttf)
	pixel_height float, // height of font in pixels
	pixels []byte, pw, ph, // bitmap to be filled in
	first_char, num_chars int, // characters to bake
	chardata []stbtt_bakedchar) int { // you allocate this, it's num_chars long

	return stbtt_BakeFontBitmap_internal(data, offset, pixel_height, pixels, pw, ph, first_char, num_chars, chardata)
}

type AlignedQuad struct {
	X0, Y0, S0, T0 float // top-left
	X1, Y1, S1, T1 float // bottom-right
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
func stbtt_GetBakedQuad(chardata []stbtt_bakedchar, pw, ph, // same data as above
	char_index int, // character to display
	xpos, ypos *float, // pointers to current position in screen pixel space
	q *AlignedQuad, // output: quad to draw
	opengl_fillrule int) { // true if opengl fill rule; false if DX9 or earlier
	var d3d_bias float
	if opengl_fillrule == 0 {
		d3d_bias = -0.5
	}
	var ipw = 1.0 / float(pw)
	var iph = 1.0 / float(ph)
	var b = &chardata[char_index]
	var round_x = STBTT_ifloor((*xpos + b.xoff) + 0.5)
	var round_y = STBTT_ifloor((*ypos + b.yoff) + 0.5)

	q.X0 = float(round_x) + d3d_bias
	q.Y0 = float(round_y) + d3d_bias
	q.X1 = float(round_x) + float(b.x1) - float(b.x0) + d3d_bias
	q.Y1 = float(round_y) + float(b.y1) - float(b.y0) + d3d_bias

	q.S0 = float(b.x0) * ipw
	q.T0 = float(b.y0) * iph
	q.S1 = float(b.x1) * ipw
	q.T1 = float(b.y1) * iph

	*xpos += b.xadvance
}

// Query the font vertical metrics without having to create a font first.
func stbtt_GetScaledFontVMetrics(fontdata []byte, index int, size float, ascent, descent, lineGap *float) {
	var i_ascent, i_descent, i_lineGap int
	var scale float
	var info FontInfo
	InitFont(&info, fontdata, GetFontOffsetForIndex(fontdata, index))
	if size > 0 {
		scale = ScaleForPixelHeight(&info, size)
	} else {
		scale = ScaleForMappingEmToPixels(&info, -size)
	}
	GetFontVMetrics(&info, &i_ascent, &i_descent, &i_lineGap)
	*ascent = (float)(i_ascent) * scale
	*descent = (float)(i_descent) * scale
	*lineGap = (float)(i_lineGap) * scale
}

type PackedChar struct {
	x0, y0, x1, y1       uint16 // coordinates of bbox in bitmap
	xoff, yoff, AdvanceX float
	xoff2, yoff2         float
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
func PackBegin(spc *PackContext, pixels []byte, width, height, stride_in_bytes, padding int, alloc_context any) int {
	var context = new(stbrp.Context)
	var num_nodes = width - padding
	var nodes = make([]stbrp.Node, num_nodes)

	if context == nil || nodes == nil {
		//if (context != nil) {STBTT_free(context, alloc_context);
		//if (nodes   != nil) STBTT_free(nodes  , alloc_context);
		return 0
	}

	spc.user_allocator_context = alloc_context
	spc.width = width
	spc.Height = height
	spc.Pixels = pixels
	spc.PackInfo = context
	spc.nodes = nodes
	spc.padding = padding
	if stride_in_bytes != 0 {
		spc.stride_in_bytes = stride_in_bytes
	} else {
		spc.stride_in_bytes = width
	}
	spc.h_oversample = 1
	spc.v_oversample = 1
	spc.skip_missing = 0

	stbrp.InitTarget(context, width-padding, height-padding, nodes, num_nodes)

	return 1
}

func STBTT_POINT_SIZE(x int) int {
	return -x
}

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
//
//	...,                  20 , ... // font max minus min y is 20 pixels tall
//	..., STBTT_POINT_SIZE(20), ... // 'M' is 20 pixels tall
func stbtt_PackFontRange(spc *PackContext, fontdata []byte, font_index int, font_size float,
	first_unicode_char_in_range, num_chars_in_range int, chardata_for_range []PackedChar) int {

	var prange [1]PackRange
	prange[0].FirstUnicodeCodepointInRange = first_unicode_char_in_range
	prange[0].ArrayOfUnicodeCodepoints = nil
	prange[0].NumChars = num_chars_in_range
	prange[0].ChardataForRange = chardata_for_range
	prange[0].FontSize = font_size
	return stbtt_PackFontRanges(spc, fontdata, font_index, prange[:], 1)
}

type PackRange struct {
	FontSize                     float
	FirstUnicodeCodepointInRange int   // if non-zero, then the chars are continuous, and this is the first codepoint
	ArrayOfUnicodeCodepoints     []int // if non-zero, then this is an array of unicode codepoints
	NumChars                     int
	ChardataForRange             []PackedChar // output

	Oversample struct {
		H, V byte
	}
}

// Creates character bitmaps from multiple ranges of characters stored in
// ranges. This will usually create a better-packed bitmap than multiple
// calls to stbtt_PackFontRange. Note that you can call this multiple
// times within a single PackBegin/PackEnd.
func stbtt_PackFontRanges(spc *PackContext, fontdata []byte, font_index int, ranges []PackRange, num_ranges int) int {
	var info FontInfo
	var i, j, n, return_value int // [DEAR IMGUI] removed = 1
	//stbrp.Context *context = (stbrp.Context *) spc.pack_info;
	var rects []stbrp.Rect

	// flag all characters as NOT packed
	for i = 0; i < num_ranges; i++ {
		for j = 0; j < ranges[i].NumChars; j++ {
			ranges[i].ChardataForRange[j].x0 = 0
			ranges[i].ChardataForRange[j].y0 = 0
			ranges[i].ChardataForRange[j].x1 = 0
			ranges[i].ChardataForRange[j].y1 = 0
		}
	}

	n = 0
	for i = 0; i < num_ranges; i++ {
		n += ranges[i].NumChars
	}

	rects = make([]stbrp.Rect, n)
	if rects == nil {
		return 0
	}

	info.userdata = spc.user_allocator_context
	InitFont(&info, fontdata, GetFontOffsetForIndex(fontdata, font_index))

	n = stbtt_PackFontRangesGatherRects(spc, &info, ranges, num_ranges, rects)

	stbtt_PackFontRangesPackRects(spc, rects, n)

	return_value = PackFontRangesRenderIntoRects(spc, &info, ranges, num_ranges, rects)

	return return_value
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
func PackSetOversampling(spc *PackContext, h_oversample, v_oversample uint) {
	STBTT_assert(h_oversample <= STBTT_MAX_OVERSAMPLE)
	STBTT_assert(v_oversample <= STBTT_MAX_OVERSAMPLE)
	if h_oversample <= STBTT_MAX_OVERSAMPLE {
		spc.h_oversample = h_oversample
	}
	if v_oversample <= STBTT_MAX_OVERSAMPLE {
		spc.v_oversample = v_oversample
	}
}

// If skip != 0, this tells stb_truetype to skip any codepoints for which
// there is no corresponding glyph. If skip=0, which is the default, then
// codepoints without a glyph recived the font's "missing character" glyph,
// typically an empty box by convention.
func stbtt_PackSetSkipMissingCodepoints(spc *PackContext, skip int) {
	spc.skip_missing = skip
}

func GetPackedQuad(chardata []PackedChar, pw, ph, // same data as above
	char_index int, // character to display
	xpos, ypos *float, // pointers to current position in screen pixel space
	q *AlignedQuad, // output: quad to draw
	align_to_integer int) {

	var ipw = 1.0 / float(pw)
	var iph = 1.0 / float(ph)
	var b = &chardata[char_index]

	if align_to_integer != 0 {
		var x = (float)(STBTT_ifloor((*xpos + b.xoff) + 0.5))
		var y = (float)(STBTT_ifloor((*ypos + b.yoff) + 0.5))
		q.X0 = x
		q.Y0 = y
		q.X1 = x + b.xoff2 - b.xoff
		q.Y1 = y + b.yoff2 - b.yoff
	} else {
		q.X0 = *xpos + b.xoff
		q.Y0 = *ypos + b.yoff
		q.X1 = *xpos + b.xoff2
		q.Y1 = *ypos + b.yoff2
	}

	q.S0 = float(b.x0) * ipw
	q.T0 = float(b.y0) * iph
	q.S1 = float(b.x1) * ipw
	q.T1 = float(b.y1) * iph

	*xpos += b.AdvanceX
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

func stbtt_PackFontRangesGatherRects(spc *PackContext, info *FontInfo, ranges []PackRange, num_ranges int, rects []stbrp.Rect) int {
	var i, j, k int

	k = 0
	for i = 0; i < num_ranges; i++ {
		var fh = ranges[i].FontSize
		var scale float
		if fh > 0 {
			scale = ScaleForPixelHeight(info, fh)
		} else {
			scale = ScaleForMappingEmToPixels(info, -fh)
		}
		ranges[i].Oversample.H = (byte)(spc.h_oversample)
		ranges[i].Oversample.V = (byte)(spc.v_oversample)
		for j = 0; j < ranges[i].NumChars; j++ {
			var x0, y0, x1, y1 int
			var codepoint int
			if ranges[i].ArrayOfUnicodeCodepoints == nil {
				codepoint = ranges[i].FirstUnicodeCodepointInRange + j
			} else {
				codepoint = ranges[i].ArrayOfUnicodeCodepoints[j]
			}
			var glyph = FindGlyphIndex(info, codepoint)
			if glyph == 0 && spc.skip_missing != 0 {
				rects[k].W = 0
				rects[k].H = 0
			} else {
				GetGlyphBitmapBoxSubpixel(info, glyph,
					scale*float(spc.h_oversample),
					scale*float(spc.v_oversample),
					0, 0,
					&x0, &y0, &x1, &y1)
				rects[k].W = (stbrp.Coord)(x1 - x0 + spc.padding + int(spc.h_oversample) - 1)
				rects[k].H = (stbrp.Coord)(y1 - y0 + spc.padding + int(spc.v_oversample) - 1)
			}
			k++
		}
	}

	return k
}

func stbtt_PackFontRangesPackRects(spc *PackContext, rects []stbrp.Rect, num_rects int) {
	stbrp.PackRects(spc.PackInfo.(*stbrp.Context), rects, num_rects)
}

func PackFontRangesRenderIntoRects(spc *PackContext, info *FontInfo, ranges []PackRange, num_ranges int, rects []stbrp.Rect) int {
	var i, j, k, return_value int = 0, 0, 0, 1

	// save current values
	var old_h_over = int(spc.h_oversample)
	var old_v_over = int(spc.v_oversample)

	k = 0
	for i = 0; i < num_ranges; i++ {
		var fh = ranges[i].FontSize
		var scale float
		if fh > 0 {
			scale = ScaleForPixelHeight(info, fh)
		} else {
			scale = ScaleForMappingEmToPixels(info, -fh)
		}
		var recip_h, recip_v, sub_x, sub_y float
		spc.h_oversample = uint(ranges[i].Oversample.H)
		spc.v_oversample = uint(ranges[i].Oversample.V)
		recip_h = 1.0 / float(spc.h_oversample)
		recip_v = 1.0 / float(spc.v_oversample)
		sub_x = stbtt__oversample_shift(int(spc.h_oversample))
		sub_y = stbtt__oversample_shift(int(spc.v_oversample))
		for j = 0; j < ranges[i].NumChars; j++ {
			var r = &rects[k]
			if r.WasPacked != 0 && r.W != 0 && r.H != 0 {
				var bc = &ranges[i].ChardataForRange[j]
				var advance, lsb, x0, y0, x1, y1 int
				var codepoint int
				if ranges[i].ArrayOfUnicodeCodepoints == nil {
					codepoint = ranges[i].FirstUnicodeCodepointInRange + j
				} else {
					codepoint = ranges[i].ArrayOfUnicodeCodepoints[j]
				}
				var glyph = FindGlyphIndex(info, codepoint)
				var pad = (stbrp.Coord)(spc.padding)

				// pad on left and top
				r.X += pad
				r.Y += pad
				r.W -= pad
				r.H -= pad
				stbtt_GetGlyphHMetrics(info, glyph, &advance, &lsb)
				stbtt_GetGlyphBitmapBox(info, glyph,
					scale*float(spc.h_oversample),
					scale*float(spc.v_oversample),
					&x0, &y0, &x1, &y1)

				stbtt_MakeGlyphBitmapSubpixel(info,
					spc.Pixels[int(r.X)+int(r.Y)*spc.stride_in_bytes:],
					int(r.W)-int(spc.h_oversample)+1,
					int(r.H)-int(spc.v_oversample)+1,
					spc.stride_in_bytes,
					scale*float(spc.h_oversample),
					scale*float(spc.v_oversample),
					0, 0,
					glyph)

				if spc.h_oversample > 1 {
					stbtt__h_prefilter(spc.Pixels[int(r.X)+int(r.Y)*spc.stride_in_bytes:],
						int(r.W), int(r.H), spc.stride_in_bytes,
						spc.h_oversample)
				}

				if spc.v_oversample > 1 {
					stbtt__v_prefilter(spc.Pixels[int(r.X)+int(r.Y)*spc.stride_in_bytes:],
						int(r.W), int(r.H), spc.stride_in_bytes,
						spc.v_oversample)
				}

				bc.x0 = uint16((stbtt_int16)(r.X))
				bc.y0 = uint16((stbtt_int16)(r.Y))
				bc.x1 = uint16((stbtt_int16)(r.X + r.W))
				bc.y1 = uint16((stbtt_int16)(r.Y + r.H))
				bc.AdvanceX = scale * float(advance)
				bc.xoff = (float)(float(x0)*recip_h + sub_x)
				bc.yoff = (float)(float(y0)*recip_v + sub_y)
				bc.xoff2 = float(x0+int(r.W))*recip_h + sub_x
				bc.yoff2 = float(y0+int(r.H))*recip_v + sub_y
			} else {
				return_value = 0 // if any fail, report failure
			}

			k++
		}
	}

	// restore original values
	spc.h_oversample = uint(old_h_over)
	spc.v_oversample = uint(old_v_over)

	return return_value
}

// this is an opaque structure that you shouldn't mess with which holds
// all the context needed from PackBegin to PackEnd.
type PackContext struct {
	user_allocator_context     any
	PackInfo                   any
	width                      int
	Height                     int
	stride_in_bytes            int
	padding                    int
	skip_missing               int
	h_oversample, v_oversample uint
	Pixels                     []byte
	nodes                      any
}

// This function will determine the number of fonts in a font file.  TrueType
// collection (.ttc) files may contain multiple fonts, while TrueType font
// (.ttf) files only contain one font. The number of fonts can be used for
// indexing with the previous function where the index is between zero and one
// less than the total fonts. If an error occurs, -1 is returned.
func stbtt_GetNumberOfFonts(data []byte) int {
	return stbtt_GetNumberOfFonts_internal(data)
}

// Each .ttf/.ttc file may have more than one font. Each font has a sequential
// index number starting from 0. Call this function to get the font offset for
// a given index; it returns -1 if the index is out of range. A regular .ttf
// file will only define one font and it always be at offset 0, so it will
// return '0' for index 0, and -1 for all other indices.
func GetFontOffsetForIndex(data []byte, index int) int {
	return stbtt_GetFontOffsetForIndex_internal(data, index)
}

// The following structure is defined publicly so you can declare one on
// the stack or as a global or etc, but you should treat it as opaque.
type FontInfo struct {
	userdata  any
	data      []byte // pointer to .ttf file
	fontstart int    // offset of start of font

	numGlyphs int // number of glyphs, needed for range checking

	loca, head, glyf, hhea, hmtx, kern, gpos int // table locations as offset from start of .ttf
	index_map                                int // a cmap mapping for our chosen character encoding
	indexToLocFormat                         int // format needed to map from glyph index to glyph

	cff         stbtt__buf // cff font data
	charstrings stbtt__buf // the charstring index
	gsubrs      stbtt__buf // global charstring subroutines index
	subrs       stbtt__buf // private charstring subroutines index
	fontdicts   stbtt__buf // array of font dicts
	fdselect    stbtt__buf // map from glyph to fontdict
}

// Given an offset into the file that defines a font, this function builds
// the necessary cached info for the rest of the system. You must allocate
// the stbtt_fontinfo yourself, and InitFont will fill it out. You don't
// need to do anything special to free it, because the contents are pure
// value data with no additional data structures. Returns 0 on failure.
func InitFont(info *FontInfo, data []byte, offset int) int {
	return stbtt_InitFont_internal(info, data, offset)
}

// If you're going to perform multiple operations on the same character
// and you want a speed-up, call this function with the character you're
// going to process, then use glyph-based functions instead of the
// codepoint-based functions.
// Returns 0 if the character codepoint is not defined in the font.
func FindGlyphIndex(info *FontInfo, unicode_codepoint int) int {
	var data = info.data
	var index_map = stbtt_uint32(info.index_map)

	var format = ttUSHORT(data[index_map+0:])
	if format == 0 { // apple byte encoding
		var bytes = stbtt_int32(ttUSHORT(data[index_map+2:]))
		if unicode_codepoint < bytes-6 {
			return int(ttBYTE(data[index_map+6+uint(unicode_codepoint):]))
		}
		return 0
	} else if format == 6 {
		var first = stbtt_uint32(ttUSHORT(data[index_map+6:]))
		var count = stbtt_uint32(ttUSHORT(data[index_map+8:]))
		if stbtt_uint32(unicode_codepoint) >= first && (stbtt_uint32)(unicode_codepoint) < first+count {
			return int(ttUSHORT(data[index_map+10+(uint(unicode_codepoint)-first)*2:]))
		}
		return 0
	} else if format == 2 {
		STBTT_assert(false) // @TODO: high-byte mapping for japanese/chinese/korean
		return 0
	} else if format == 4 { // standard mapping for windows fonts: binary search collection of ranges
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
		if unicode_codepoint >= int(ttUSHORT(data[search+uint(rangeShift)*2:])) {
			search += uint(rangeShift) * 2
		}

		// now decrement to bias correctly to find smallest
		search -= 2
		for entrySelector != 0 {
			var end stbtt_uint16
			searchRange >>= 1
			end = ttUSHORT(data[search+uint(searchRange)*2:])
			if unicode_codepoint > int(end) {
				search += uint(searchRange) * 2
			}
			entrySelector--
		}
		search += 2

		{
			var offset, start stbtt_uint16
			var item = (stbtt_uint16)((search - endCount) >> 1)

			STBTT_assert(unicode_codepoint <= int(ttUSHORT(data[endCount+2*uint(item):])))
			start = ttUSHORT(data[index_map+14+uint(segcount)*2+2+2*uint(item):])
			if unicode_codepoint < int(start) {
				return 0
			}

			offset = ttUSHORT(data[index_map+14+uint(segcount)*6+2+2*uint(item):])
			if offset == 0 {
				return int((stbtt_uint16)(unicode_codepoint + int(ttSHORT(data[index_map+14+uint(segcount)*4+2+2*uint(item):]))))
			}

			return int(ttUSHORT(data[uint(offset)+(uint(unicode_codepoint)-uint(start))*2+index_map+14+uint(segcount)*6+2+2*uint(item):]))
		}
	} else if format == 12 || format == 13 {
		var ngroups = ttULONG(data[index_map+12:])
		var low, high stbtt_int32
		low, high = 0, (stbtt_int32)(ngroups)
		// Binary search the right group.
		for low < high {
			var mid = low + ((high - low) >> 1) // rounds down, so low <= mid < high
			var start_char = ttULONG(data[index_map+16+uint(mid)*12:])
			var end_char = ttULONG(data[index_map+16+uint(mid)*12+4:])
			if (stbtt_uint32)(unicode_codepoint) < start_char {
				high = mid
			} else if (stbtt_uint32)(unicode_codepoint) > end_char {
				low = mid + 1
			} else {
				var start_glyph = ttULONG(data[index_map+16+uint(mid)*12+8:])
				if format == 12 {
					return int(start_glyph + uint(unicode_codepoint) - uint(start_char))
				} else { // format == 13
					return int(start_glyph)
				}
			}
		}
		return 0 // not found
	}
	// @TODO
	STBTT_assert(false)
	return 0
}

// computes a scale factor to produce a font whose "height" is 'pixels' tall.
// Height is measured as the distance from the highest ascender to the lowest
// descender; in other words, it's equivalent to calling stbtt_GetFontVMetrics
// and computing:
//
//	scale = pixels / (ascent - descent)
//
// so if you prefer to measure height by the ascent only, use a similar calculation.
func ScaleForPixelHeight(info *FontInfo, pixels float) float {
	var fheight = int(ttSHORT(info.data[info.hhea+4:])) - int(ttSHORT(info.data[info.hhea+6:]))
	return pixels / (float)(fheight)
}

// computes a scale factor to produce a font whose EM size is mapped to
// 'pixels' tall. This is probably what traditional APIs compute, but
// I'm not positive.
func ScaleForMappingEmToPixels(info *FontInfo, pixels float) float {
	var unitsPerEm = int(ttUSHORT(info.data[info.head+18:]))
	return pixels / float(unitsPerEm)
}

// ascent is the coordinate above the baseline the font extends; descent
// is the coordinate below the baseline the font extends (i.e. it is typically negative)
// lineGap is the spacing between one row's descent and the next row's ascent...
// so you should advance the vertical position by "*ascent - *descent + *lineGap"
//
//	these are expressed in unscaled coordinates, so you must multiply by
//	the scale factor for a given size
func GetFontVMetrics(info *FontInfo, ascent, descent, lineGap *int) {
	if ascent != nil {
		*ascent = int(ttSHORT(info.data[info.hhea+4:]))
	}
	if descent != nil {
		*descent = int(ttSHORT(info.data[info.hhea+6:]))
	}
	if lineGap != nil {
		*lineGap = int(ttSHORT(info.data[info.hhea+8:]))
	}
}

// analogous to GetFontVMetrics, but returns the "typographic" values from the OS/2
// table (specific to MS/Windows TTF files).
//
// Returns 1 on success (table present), 0 on failure.
func stbtt_GetFontVMetricsOS2(info *FontInfo, typoAscent, typoDescent, typoLineGap *int) int {
	var tab = int(stbtt__find_table(info.data, uint(info.fontstart), "OS/2"))
	if tab == 0 {
		return 0
	}
	if typoAscent != nil {
		*typoAscent = int(ttSHORT(info.data[tab+68:]))
	}
	if typoDescent != nil {
		*typoDescent = int(ttSHORT(info.data[tab+70:]))
	}
	if typoLineGap != nil {
		*typoLineGap = int(ttSHORT(info.data[tab+72:]))
	}
	return 1
}

// the bounding box around all possible characters
func stbtt_GetFontBoundingBox(info *FontInfo, x0, y0, x1, y1 *int) {
	*x0 = int(ttSHORT(info.data[info.head+36:]))
	*y0 = int(ttSHORT(info.data[info.head+38:]))
	*x1 = int(ttSHORT(info.data[info.head+40:]))
	*y1 = int(ttSHORT(info.data[info.head+42:]))
}

// leftSideBearing is the offset from the current horizontal position to the left edge of the character
// advanceWidth is the offset from the current horizontal position to the next horizontal position
//
//	these are expressed in unscaled coordinates
func stbtt_GetCodepointHMetrics(info *FontInfo, codepoint int, advanceWidth, leftSideBearing *int) {
	stbtt_GetGlyphHMetrics(info, FindGlyphIndex(info, codepoint), advanceWidth, leftSideBearing)
}

// an additional amount to add to the 'advance' value between ch1 and ch2
func stbtt_GetCodepointKernAdvance(info *FontInfo, ch1, ch2 int) int {
	if info.kern == 0 && info.gpos == 0 { // if no kerning table, don't waste time looking up both codepoint.glyphs
		return 0
	}
	return stbtt_GetGlyphKernAdvance(info, FindGlyphIndex(info, ch1), FindGlyphIndex(info, ch2))
}

// Gets the bounding box of the visible part of the glyph, in unscaled coordinates
func stbtt_GetCodepointBox(info *FontInfo, codepoint int, x0, y0, x1, y1 *int) int {
	return stbtt_GetGlyphBox(info, FindGlyphIndex(info, codepoint), x0, y0, x1, y1)
}

// as above, but takes one or more glyph indices for greater efficiency

func stbtt_GetGlyphHMetrics(info *FontInfo, glyph_index int, advanceWidth, leftSideBearing *int) {
	var numOfLongHorMetrics = ttUSHORT(info.data[info.hhea+34:])
	if glyph_index < int(numOfLongHorMetrics) {
		if advanceWidth != nil {
			*advanceWidth = int(ttSHORT(info.data[info.hmtx+4*glyph_index:]))
		}
		if leftSideBearing != nil {
			*leftSideBearing = int(ttSHORT(info.data[info.hmtx+4*glyph_index+2:]))
		}
	} else {
		if advanceWidth != nil {
			*advanceWidth = int(ttSHORT(info.data[info.hmtx+4*(int(numOfLongHorMetrics)-1):]))
		}
		if leftSideBearing != nil {
			*leftSideBearing = int(ttSHORT(info.data[info.hmtx+4*int(numOfLongHorMetrics)+2*(glyph_index-int(numOfLongHorMetrics)):]))
		}
	}
}

func stbtt_GetGlyphKernAdvance(info *FontInfo, glyph1, glyph2 int) int {
	var xAdvance int = 0

	if info.gpos != 0 {
		xAdvance += stbtt__GetGlyphGPOSInfoAdvance(info, glyph1, glyph2)
	}

	if info.kern != 0 {
		xAdvance += stbtt__GetGlyphKernInfoAdvance(info, glyph1, glyph2)
	}

	return xAdvance
}

func stbtt_GetGlyphBox(info *FontInfo, glyph_index int, x0, y0, x1, y1 *int) int {
	if info.cff.size != 0 {
		stbtt__GetGlyphInfoT2(info, glyph_index, x0, y0, x1, y1)
	} else {
		var g = stbtt__GetGlyfOffset(info, glyph_index)
		if g < 0 {
			return 0
		}

		if x0 != nil {
			*x0 = int(ttSHORT(info.data[g+2:]))
		}
		if y0 != nil {
			*y0 = int(ttSHORT(info.data[g+4:]))
		}
		if x1 != nil {
			*x1 = int(ttSHORT(info.data[g+6:]))
		}
		if y1 != nil {
			*y1 = int(ttSHORT(info.data[g+8:]))
		}
	}
	return 1
}

const (
	STBTT_vmove = iota + 1
	STBTT_vline
	STBTT_vcurve
	STBTT_vcubic
)

type stbtt_vertex_type = int16
type stbtt_vertex struct {
	x, y, cx, cy, cx1, cy1 stbtt_vertex_type
	vtype, padding         byte
}

// returns non-zero if nothing is drawn for this glyph
func stbtt_IsGlyphEmpty(info *FontInfo, glyph_index int) int {
	var numberOfContours stbtt_int16
	var g int
	if info.cff.size != 0 {
		return bool2int(stbtt__GetGlyphInfoT2(info, glyph_index, nil, nil, nil, nil) == 0)
	}
	g = stbtt__GetGlyfOffset(info, glyph_index)
	if g < 0 {
		return 1
	}
	numberOfContours = ttSHORT(info.data[g:])
	return bool2int(numberOfContours == 0)
}

// returns # of vertices and fills *vertices with the pointer to them
//
//	these are expressed in "unscaled" coordinates
//
// The shape is a series of contours. Each one starts with
// a STBTT_moveto, then consists of a series of mixed
// STBTT_lineto and STBTT_curveto segments. A lineto
// draws a line from previous endpoint to its x,y; a curveto
// draws a quadratic bezier from previous endpoint to
// its x,y, using cx,cy as the bezier control point.
func stbtt_GetCodepointShape(info *FontInfo, unicode_codepoint int, vertices *[]stbtt_vertex) int {
	return stbtt_GetGlyphShape(info, FindGlyphIndex(info, unicode_codepoint), vertices)
}

// frees the data allocated above
func stbtt_FreeShape(info *FontInfo, vertices []stbtt_vertex) {}

// frees the bitmap allocated below
func stbtt_FreeBitmap(bitmap []byte, userdata any) {}

// allocates a large-enough single-channel 8bpp bitmap and renders the
// specified character/glyph at the specified scale into it, with
// antialiasing. 0 is no coverage (transparent), 255 is fully covered (opaque).
// *width & *height are filled out with the width & height of the bitmap,
// which is stored left-to-right, top-to-bottom.
//
// xoff/yoff are the offset it pixel space from the glyph origin to the top-left of the bitmap
func stbtt_GetCodepointBitmap(info *FontInfo, scale_x, scale_y float, codepoint int, width, height, xoff, yoff *int) []byte {
	return stbtt_GetCodepointBitmapSubpixel(info, scale_x, scale_y, 0.0, 0.0, codepoint, width, height, xoff, yoff)
}

// the same as stbtt_GetCodepoitnBitmap, but you can specify a subpixel
// shift for the character
func stbtt_GetCodepointBitmapSubpixel(info *FontInfo, scale_x, scale_y, shift_x, shift_y float, codepoint int, width, height, xoff, yoff *int) []byte {
	return stbtt_GetGlyphBitmapSubpixel(info, scale_x, scale_y, shift_x, shift_y, FindGlyphIndex(info, codepoint), width, height, xoff, yoff)
}

// the same as stbtt_GetCodepointBitmap, but you pass in storage for the bitmap
// in the form of 'output', with row spacing of 'out_stride' bytes. the bitmap
// is clipped to out_w/out_h bytes. Call stbtt_GetCodepointBitmapBox to get the
// width and height and positioning info for it first.
func stbtt_MakeCodepointBitmap(info *FontInfo, output []byte, out_w, out_h, out_stride int, scale_x, scale_y float, codepoint int) {
	stbtt_MakeCodepointBitmapSubpixel(info, output, out_w, out_h, out_stride, scale_x, scale_y, 0.0, 0.0, codepoint)
}

// same as stbtt_MakeCodepointBitmap, but you can specify a subpixel
// shift for the character
func stbtt_MakeCodepointBitmapSubpixel(info *FontInfo, output []byte, out_w, out_h, out_stride int, scale_x, scale_y, shift_x, shift_y float, codepoint int) {
	stbtt_MakeGlyphBitmapSubpixel(info, output, out_w, out_h, out_stride, scale_x, scale_y, shift_x, shift_y, FindGlyphIndex(info, codepoint))
}

// same as stbtt_MakeCodepointBitmapSubpixel, but prefiltering
// is performed (see stbtt_PackSetOversampling)
func stbtt_MakeCodepointBitmapSubpixelPrefilter(info *FontInfo, output []byte, out_w, out_h, out_stride int, scale_x, scale_y, shift_x, shift_y float, oversample_x, oversample_y int, sub_x, sub_y *float, codepoint int) {
	stbtt_MakeGlyphBitmapSubpixelPrefilter(info, output, out_w, out_h, out_stride, scale_x, scale_y, shift_x, shift_y, oversample_x, oversample_y, sub_x, sub_y, FindGlyphIndex(info, codepoint))
}

// get the bbox of the bitmap centered around the glyph origin; so the
// bitmap width is ix1-ix0, height is iy1-iy0, and location to place
// the bitmap top left is (leftSideBearing*scale,iy0).
// (Note that the bitmap uses y-increases-down, but the shape uses
// y-increases-up, so CodepointBitmapBox and CodepointBox are inverted.)
func stbtt_GetCodepointBitmapBox(font *FontInfo, codepoint int, scale_x, scale_y float, ix0, iy0, ix1, iy1 *int) {
	stbtt_GetCodepointBitmapBoxSubpixel(font, codepoint, scale_x, scale_y, 0.0, 0.0, ix0, iy0, ix1, iy1)
}

// same as stbtt_GetCodepointBitmapBox, but you can specify a subpixel
// shift for the character
func stbtt_GetCodepointBitmapBoxSubpixel(font *FontInfo, codepoint int, scale_x, scale_y, shift_x, shift_y float, ix0, iy0, ix1, iy1 *int) {
	GetGlyphBitmapBoxSubpixel(font, FindGlyphIndex(font, codepoint), scale_x, scale_y, shift_x, shift_y, ix0, iy0, ix1, iy1)
}

// the following functions are equivalent to the above functions, but operate
// on glyph indices instead of Unicode codepoints (for efficiency)
func stbtt_GetGlyphBitmap(info *FontInfo, scale_x, scale_y float, glyph int, width, height, xoff, yoff *int) []byte {
	return stbtt_GetGlyphBitmapSubpixel(info, scale_x, scale_y, 0.0, 0.0, glyph, width, height, xoff, yoff)
}

func stbtt_GetGlyphBitmapSubpixel(info *FontInfo, scale_x, scale_y, shift_x, shift_y float, glyph int, width, height, xoff, yoff *int) []byte {
	var ix0, iy0, ix1, iy1 int
	var gbm stbtt__bitmap
	var vertices []stbtt_vertex
	var num_verts = stbtt_GetGlyphShape(info, glyph, &vertices)

	if scale_x == 0 {
		scale_x = scale_y
	}
	if scale_y == 0 {
		if scale_x == 0 {
			return nil
		}
		scale_y = scale_x
	}

	GetGlyphBitmapBoxSubpixel(info, glyph, scale_x, scale_y, shift_x, shift_y, &ix0, &iy0, &ix1, &iy1)

	// now we get the size
	gbm.w = (ix1 - ix0)
	gbm.h = (iy1 - iy0)
	gbm.pixels = nil // in case we error

	if width != nil {
		*width = gbm.w
	}
	if height != nil {
		*height = gbm.h
	}
	if xoff != nil {
		*xoff = ix0
	}
	if yoff != nil {
		*yoff = iy0
	}

	if gbm.w != 0 && gbm.h != 0 {
		gbm.pixels = make([]byte, gbm.w*gbm.h)
		if gbm.pixels != nil {
			gbm.stride = gbm.w

			stbtt_Rasterize(&gbm, 0.35, vertices, num_verts, scale_x, scale_y, shift_x, shift_y, ix0, iy0, 1, info.userdata)
		}
	}
	return gbm.pixels
}

func stbtt_MakeGlyphBitmap(info *FontInfo, output []byte, out_w, out_h, out_stride int, scale_x, scale_y float, glyph int) {
	stbtt_MakeGlyphBitmapSubpixel(info, output, out_w, out_h, out_stride, scale_x, scale_y, 0.0, 0.0, glyph)
}

func stbtt_MakeGlyphBitmapSubpixel(info *FontInfo, output []byte, out_w, out_h, out_stride int, scale_x, scale_y, shift_x, shift_y float, glyph int) {
	var ix0, iy0 int
	var vertices []stbtt_vertex
	var num_verts = stbtt_GetGlyphShape(info, glyph, &vertices)
	var gbm stbtt__bitmap

	GetGlyphBitmapBoxSubpixel(info, glyph, scale_x, scale_y, shift_x, shift_y, &ix0, &iy0, nil, nil)
	gbm.pixels = output
	gbm.w = out_w
	gbm.h = out_h
	gbm.stride = out_stride

	if gbm.w != 0 && gbm.h != 0 {

		stbtt_Rasterize(&gbm, 0.35, vertices, num_verts, scale_x, scale_y, shift_x, shift_y, ix0, iy0, 1, info.userdata)
	}
}

func stbtt_MakeGlyphBitmapSubpixelPrefilter(info *FontInfo, output []byte, out_w, out_h, out_stride int, scale_x, scale_y, shift_x, shift_y float, oversample_x, oversample_y int, sub_x, sub_y *float, glyph int) {
	stbtt_MakeGlyphBitmapSubpixel(info,
		output,
		out_w-(oversample_x-1),
		out_h-(oversample_y-1),
		out_stride,
		scale_x,
		scale_y,
		shift_x,
		shift_y,
		glyph)

	if oversample_x > 1 {
		stbtt__h_prefilter(output, out_w, out_h, out_stride, uint(oversample_x))
	}

	if oversample_y > 1 {
		stbtt__v_prefilter(output, out_w, out_h, out_stride, uint(oversample_y))
	}

	*sub_x = stbtt__oversample_shift(oversample_x)
	*sub_y = stbtt__oversample_shift(oversample_y)
}

func stbtt_GetGlyphBitmapBox(font *FontInfo, glyph int, scale_x, scale_y float, ix0, iy0, ix1, iy1 *int) {
	GetGlyphBitmapBoxSubpixel(font, glyph, scale_x, scale_y, 0.0, 0.0, ix0, iy0, ix1, iy1)
}

func GetGlyphBitmapBoxSubpixel(font *FontInfo, glyph int, scale_x, scale_y, shift_x, shift_y float, ix0, iy0, ix1, iy1 *int) {
	var x0, y0, x1, y1 int // =0 suppresses compiler warning
	if stbtt_GetGlyphBox(font, glyph, &x0, &y0, &x1, &y1) == 0 {
		// e.g. space character
		if ix0 != nil {
			*ix0 = 0
		}
		if iy0 != nil {
			*iy0 = 0
		}
		if ix1 != nil {
			*ix1 = 0
		}
		if iy1 != nil {
			*iy1 = 0
		}
	} else {
		// move to integral bboxes (treating pixels as little squares, what pixels get touched)?
		if ix0 != nil {
			*ix0 = STBTT_ifloor(float32(x0)*scale_x + shift_x)
		}
		if iy0 != nil {
			*iy0 = STBTT_ifloor(-float32(y1)*scale_y + shift_y)
		}
		if ix1 != nil {
			*ix1 = STBTT_iceil(float32(x1)*scale_x + shift_x)
		}
		if iy1 != nil {
			*iy1 = STBTT_iceil(-float32(y0)*scale_y + shift_y)
		}
	}
}

type stbtt__bitmap struct {
	w, h, stride int
	pixels       []byte
}

// rasterize a shape with quadratic beziers into a bitmap
func stbtt_Rasterize(result *stbtt__bitmap, // 1-channel bitmap to draw into
	flatness_in_pixels float, // allowable error of curve in pixels
	vertices []stbtt_vertex, // array of vertices defining shape
	num_verts int, // number of vertices in above array
	scale_x, scale_y, // scale applied to input vertices
	shift_x, shift_y float, // translation applied to input vertices
	x_off, y_off int, // another translation applied to input
	invert int, // if non-zero, vertically flip shape
	userdata any) { // context for to STBTT_MALLOC

	var scale float
	if scale_x > scale_y {
		scale = scale_y
	} else {
		scale = scale_x
	}
	var winding_count int = 0
	var winding_lengths []int = nil
	var windings = stbtt_FlattenCurves(vertices, num_verts, flatness_in_pixels/scale, &winding_lengths, &winding_count, userdata)
	if windings != nil {
		//fmt.Println("rasterize2")
		stbtt__rasterize(result, windings, winding_lengths, winding_count, scale_x, scale_y, shift_x, shift_y, x_off, y_off, invert, userdata)
	}
}

// frees the SDF bitmap allocated below
func stbtt_FreeSDF(bitmap []byte, userdata any) {
	panic("not implemented")
}

func stbtt_GetGlyphSDF(info *FontInfo, scale float, glyph, padding int, onedge_value byte, pixel_dist_scale float, width, height, xoff, yoff *int) []byte {
	var scale_x, scale_y = scale, scale
	var ix0, iy0, ix1, iy1 int
	var w, h int
	var data []byte

	// if one scale is 0, use same scale for both
	if scale_x == 0 {
		scale_x = scale_y
	}

	if scale_y == 0 {
		if scale_x == 0 {
			return nil // if both scales are 0, return nil
		}
		scale_y = scale_x
	}

	GetGlyphBitmapBoxSubpixel(info, glyph, scale, scale, 0.0, 0.0, &ix0, &iy0, &ix1, &iy1)

	// if empty, return nil
	if ix0 == ix1 || iy0 == iy1 {
		return nil
	}

	ix0 -= padding
	iy0 -= padding
	ix1 += padding
	iy1 += padding

	w = (ix1 - ix0)
	h = (iy1 - iy0)

	if width != nil {
		*width = w
	}
	if height != nil {
		*height = h
	}
	if xoff != nil {
		*xoff = ix0
	}
	if yoff != nil {
		*yoff = iy0
	}

	// invert for y-downwards bitmaps
	scale_y = -scale_y

	{
		var x, y, i, j int
		var precompute []float
		var verts []stbtt_vertex
		var num_verts = stbtt_GetGlyphShape(info, glyph, &verts)
		data = make([]byte, w*h)
		precompute = make([]float, num_verts)

		for j = num_verts - 1; i < num_verts; j, i = i, i+1 {
			if verts[i].vtype == STBTT_vline {
				var x0 = float(verts[i].x) * scale_x
				var y0 = float(verts[i].y) * scale_y
				var x1 = float(verts[j].x) * scale_x
				var y1 = float(verts[j].y) * scale_y
				var dist = (float)(STBTT_sqrt((x1-x0)*(x1-x0) + (y1-y0)*(y1-y0)))
				if dist == 0 {
					precompute[i] = 0
				} else {
					precompute[i] = 1.0 / dist
				}
			} else if verts[i].vtype == STBTT_vcurve {
				var x2 = float(verts[j].x) * scale_x
				var y2 = float(verts[j].y) * scale_y
				var x1 = float(verts[i].cx) * scale_x
				var y1 = float(verts[i].cy) * scale_y
				var x0 = float(verts[i].x) * scale_x
				var y0 = float(verts[i].y) * scale_y
				var bx = x0 - 2*x1 + x2
				var by = y0 - 2*y1 + y2
				var len2 = bx*bx + by*by
				if len2 != 0.0 {
					precompute[i] = 1.0 / (bx*bx + by*by)
				} else {
					precompute[i] = 0.0
				}
			} else {
				precompute[i] = 0.0
			}
		}

		for y = iy0; y < iy1; y++ {
			for x = ix0; x < ix1; x++ {
				var val float
				var min_dist float = 999999.0
				var sx = (float)(x) + 0.5
				var sy = (float)(y) + 0.5
				var x_gspace = (sx / scale_x)
				var y_gspace = (sy / scale_y)

				var winding = stbtt__compute_crossings_x(x_gspace, y_gspace, num_verts, verts) // @OPTIMIZE: this could just be a rasterization, but needs to be line vs. non-tesselated curves so a new path

				for i = 0; i < num_verts; i++ {
					var x0 = float(verts[i].x) * scale_x
					var y0 = float(verts[i].y) * scale_y

					// check against every point here rather than inside line/curve primitives -- @TODO: wrong if multiple 'moves' in a row produce a garbage point, and given culling, probably more efficient to do within line/curve
					var dist2 = (x0-sx)*(x0-sx) + (y0-sy)*(y0-sy)
					if dist2 < min_dist*min_dist {
						min_dist = (float)(STBTT_sqrt(dist2))
					}

					if verts[i].vtype == STBTT_vline {
						var x1 = float(verts[i-1].x) * scale_x
						var y1 = float(verts[i-1].y) * scale_y

						// coarse culling against bbox
						//if (sx > STBTT_min(x0,x1)-min_dist && sx < STBTT_max(x0,x1)+min_dist &&
						//    sy > STBTT_min(y0,y1)-min_dist && sy < STBTT_max(y0,y1)+min_dist)
						var dist = (float)(STBTT_fabs((x1-x0)*(y0-sy)-(y1-y0)*(x0-sx)) * precompute[i])
						STBTT_assert(i != 0)
						if dist < min_dist {
							// check position along line
							// x' = x0 + t*(x1-x0), y' = y0 + t*(y1-y0)
							// minimize (x'-sx)*(x'-sx)+(y'-sy)*(y'-sy)
							var dx = x1 - x0
							var dy = y1 - y0
							var px = x0 - sx
							var py = y0 - sy
							// minimize (px+t*dx)^2 + (py+t*dy)^2 = px*px + 2*px*dx*t + t^2*dx*dx + py*py + 2*py*dy*t + t^2*dy*dy
							// derivative: 2*px*dx + 2*py*dy + (2*dx*dx+2*dy*dy)*t, set to 0 and solve
							var t = -(px*dx + py*dy) / (dx*dx + dy*dy)
							if t >= 0.0 && t <= 1.0 {
								min_dist = dist
							}
						}
					} else if verts[i].vtype == STBTT_vcurve {
						var x2 = float(verts[i-1].x) * scale_x
						var y2 = float(verts[i-1].y) * scale_y
						var x1 = float(verts[i].cx) * scale_x
						var y1 = float(verts[i].cy) * scale_y
						var box_x0 = STBTT_minf(STBTT_minf(x0, x1), x2)
						var box_y0 = STBTT_minf(STBTT_minf(y0, y1), y2)
						var box_x1 = STBTT_maxf(STBTT_maxf(x0, x1), x2)
						var box_y1 = STBTT_maxf(STBTT_maxf(y0, y1), y2)
						// coarse culling against bbox to avoid computing cubic unnecessarily
						if sx > box_x0-min_dist && sx < box_x1+min_dist && sy > box_y0-min_dist && sy < box_y1+min_dist {
							var num int = 0
							var ax = x1 - x0
							var ay = y1 - y0
							var bx = x0 - 2*x1 + x2
							var by = y0 - 2*y1 + y2
							var mx = x0 - sx
							var my = y0 - sy
							var res [3]float
							var px, py, t, it float
							var a_inv = precompute[i]
							if a_inv == 0.0 { // if a_inv is 0, it's 2nd degree so use quadratic formula
								var a = 3 * (ax*bx + ay*by)
								var b = 2*(ax*ax+ay*ay) + (mx*bx + my*by)
								var c = mx*ax + my*ay
								if a == 0.0 { // if a is 0, it's linear
									if b != 0.0 {
										res[num] = -c / b
										num++
									}
								} else {
									var discriminant = b*b - 4*a*c
									if discriminant < 0 {
										num = 0
									} else {
										var root = (float)(STBTT_sqrt(discriminant))
										res[0] = (-b - root) / (2 * a)
										res[1] = (-b + root) / (2 * a)
										num = 2 // don't bother distinguishing 1-solution case, as code below will still work
									}
								}
							} else {
								var b = 3 * (ax*bx + ay*by) * a_inv // could precompute this as it doesn't depend on sample point
								var c = (2*(ax*ax+ay*ay) + (mx*bx + my*by)) * a_inv
								var d = (mx*ax + my*ay) * a_inv
								num = stbtt__solve_cubic(b, c, d, res[:])
							}
							if num >= 1 && res[0] >= 0.0 && res[0] <= 1.0 {
								t = res[0]
								it = 1.0 - t
								px = it*it*x0 + 2*t*it*x1 + t*t*x2
								py = it*it*y0 + 2*t*it*y1 + t*t*y2
								dist2 = (px-sx)*(px-sx) + (py-sy)*(py-sy)
								if dist2 < min_dist*min_dist {
									min_dist = (float)(STBTT_sqrt(dist2))
								}
							}
							if num >= 2 && res[1] >= 0.0 && res[1] <= 1.0 {
								t = res[1]
								it = 1.0 - t
								px = it*it*x0 + 2*t*it*x1 + t*t*x2
								py = it*it*y0 + 2*t*it*y1 + t*t*y2
								dist2 = (px-sx)*(px-sx) + (py-sy)*(py-sy)
								if dist2 < min_dist*min_dist {
									min_dist = (float)(STBTT_sqrt(dist2))
								}
							}
							if num >= 3 && res[2] >= 0.0 && res[2] <= 1.0 {
								t = res[2]
								it = 1.0 - t
								px = it*it*x0 + 2*t*it*x1 + t*t*x2
								py = it*it*y0 + 2*t*it*y1 + t*t*y2
								dist2 = (px-sx)*(px-sx) + (py-sy)*(py-sy)
								if dist2 < min_dist*min_dist {
									min_dist = (float)(STBTT_sqrt(dist2))
								}
							}
						}
					}
				}
				if winding == 0 {
					min_dist = -min_dist // if outside the shape, value is negative
				}
				val = float(onedge_value) + pixel_dist_scale*min_dist
				if val < 0 {
					val = 0
				} else if val > 255 {
					val = 255
				}
				data[(y-iy0)*w+(x-ix0)] = (byte)(val)
			}
		}
	}
	return data
}

// These functions compute a discretized SDF field for a single character, suitable for storing
// in a single-channel texture, sampling with bilinear filtering, and testing against
// larger than some threshold to produce scalable fonts.
//
//	info              --  the font
//	scale             --  controls the size of the resulting SDF bitmap, same as it would be creating a regular bitmap
//	glyph/codepoint   --  the character to generate the SDF for
//	padding           --  extra "pixels" around the character which are filled with the distance to the character (not 0),
//	                         which allows effects like bit outlines
//	onedge_value      --  value 0-255 to test the SDF against to reconstruct the character (i.e. the isocontour of the character)
//	pixel_dist_scale  --  what value the SDF should increase by when moving one SDF "pixel" away from the edge (on the 0..255 scale)
//	                         if positive, > onedge_value is inside; if negative, < onedge_value is inside
//	width,height      --  output height & width of the SDF bitmap (including padding)
//	xoff,yoff         --  output origin of the character
//	return value      --  a 2D array of bytes 0..255, width*height in size
//
// pixel_dist_scale & onedge_value are a scale & bias that allows you to make
// optimal use of the limited 0..255 for your application, trading off precision
// and special effects. SDF values outside the range 0..255 are clamped to 0..255.
//
// Example:
//
//	scale = stbtt_ScaleForPixelHeight(22)
//	padding = 5
//	onedge_value = 180
//	pixel_dist_scale = 180/5.0 = 36.0
//
//	This will create an SDF bitmap in which the character is about 22 pixels
//	high but the whole bitmap is about 22+5+5=32 pixels high. To produce a filled
//	shape, sample the SDF at each pixel and fill the pixel if the SDF value
//	is greater than or equal to 180/255. (You'll actually want to antialias,
//	which is beyond the scope of this example.) Additionally, you can compute
//	offset outlines (e.g. to stroke the character border inside & outside,
//	or only outside). For example, to fill outside the character up to 3 SDF
//	pixels, you would compare against (180-36.0*3)/255 = 72/255. The above
//	choice of variables maps a range from 5 pixels outside the shape to
//	2 pixels inside the shape to 0..255; this is intended primarily for apply
//	outside effects only (the interior range is needed to allow proper
//	antialiasing of the font at *smaller* sizes)
//
// The function computes the SDF analytically at each SDF pixel, not by e.g.
// building a higher-res bitmap and approximating it. In theory the quality
// should be as high as possible for an SDF of this size & representation, but
// unclear if this is true in practice (perhaps building a higher-res bitmap
// and computing from that can allow drop-out prevention).
//
// The algorithm has not been optimized at all, so expect it to be slow
// if computing lots of characters or very large sizes.
func stbtt_GetCodepointSDF(info *FontInfo, scale float, codepoint, padding int, onedge_value byte, pixel_dist_scale float, width, height, xoff, yoff *int) []byte {
	return stbtt_GetGlyphSDF(info, scale, FindGlyphIndex(info, codepoint), padding, onedge_value, pixel_dist_scale, width, height, xoff, yoff)
}

//////////////////////////////////////////////////////////////////////////////
//
// Finding the right font...
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
//     stbtt_FindMatchingFont() will use *case-sensitive* comparisons on
//             unicode-encoded names to try to find the font you want;
//             you can run this before calling stbtt_InitFont()
//
//     stbtt_GetFontNameString() lets you get any of the various strings
//             from the file yourself and do your own comparisons on them.
//             You have to have called stbtt_InitFont() first.

// returns the offset (not index) of the font that matches, or -1 if none
//
//	if you use STBTT_MACSTYLE_DONTCARE, use a font name like "Arial Bold".
//	if you use any other flag, use a font name like "Arial"; this checks
//	  the 'macStyle' header field; i don't know if fonts set this consistently
func stbtt_FindMatchingFont(fontdata []byte, name string, flags int) int {
	return stbtt_FindMatchingFont_internal(fontdata, []byte(name), flags)
}

const (
	STBTT_MACSTYLE_DONTCARE   = 0
	STBTT_MACSTYLE_BOLD       = 1
	STBTT_MACSTYLE_ITALIC     = 2
	STBTT_MACSTYLE_UNDERSCORE = 4
	STBTT_MACSTYLE_NONE       = 8 // <= not same as 0, this makes us check the bitfield is 0
)

// returns 1/0 whether the first string interpreted as utf8 is identical to
// the second string interpreted as big-endian utf16... useful for strings from next func
func stbtt_CompareUTF8toUTF16_bigendian(s1 string, len1 int, s2 string, len2 int) int {
	return stbtt_CompareUTF8toUTF16_bigendian_internal([]byte(s1), len1, []byte(s2), len2)
}

// returns the string (which may be big-endian double byte, e.g. for unicode)
// and puts the length in bytes in *length.
//
// some of the values for the IDs are below; for more see the truetype spec:
//
//	http://developer.apple.com/textfonts/TTRefMan/RM06/Chap6name.html
//	http://www.microsoft.com/typography/otspec/name.htm
func stbtt_GetFontNameString(font FontInfo, length *int, platformID, encodingID, languageID, nameID int) string {
	var i, count, stringOffset stbtt_int32
	var fc = font.data
	var offset = uint(font.fontstart)
	var nm = stbtt__find_table(fc, offset, "name")
	if nm == 0 {
		return ""
	}

	count = int(ttUSHORT(fc[nm+2:]))
	stringOffset = int(nm + uint(ttUSHORT(fc[nm+4:])))
	for i = 0; i < count; i++ {
		var loc = nm + 6 + uint(12*i)
		if platformID == int(ttUSHORT(fc[loc+0:])) &&
			encodingID == int(ttUSHORT(fc[loc+2:])) &&
			languageID == int(ttUSHORT(fc[loc+4:])) &&
			nameID == int(ttUSHORT(fc[loc+6:])) {

			*length = int(ttUSHORT(fc[loc+8:]))
			pos := stringOffset + int(ttUSHORT(fc[loc+10:]))
			return string(fc[pos : pos+*length])
		}
	}
	return ""
}

const ( // platformID
	STBTT_PLATFORM_ID_UNICODE   = 0
	STBTT_PLATFORM_ID_MAC       = 1
	STBTT_PLATFORM_ID_ISO       = 2
	STBTT_PLATFORM_ID_MICROSOFT = 3
)

const ( // encodingID for STBTT_PLATFORM_ID_UNICODE
	STBTT_UNICODE_EID_UNICODE_1_0      = 0
	STBTT_UNICODE_EID_UNICODE_1_1      = 1
	STBTT_UNICODE_EID_ISO_10646        = 2
	STBTT_UNICODE_EID_UNICODE_2_0_BMP  = 3
	STBTT_UNICODE_EID_UNICODE_2_0_FULL = 4
)

const ( // encodingID for STBTT_PLATFORM_ID_MICROSOFT
	STBTT_MS_EID_SYMBOL       = 0
	STBTT_MS_EID_UNICODE_BMP  = 1
	STBTT_MS_EID_SHIFTJIS     = 2
	STBTT_MS_EID_UNICODE_FULL = 10
)

const ( // encodingID for STBTT_PLATFORM_ID_MAC; same as Script Manager codes
	STBTT_MAC_EID_ROMAN        = 0
	STBTT_MAC_EID_ARABIC       = 4
	STBTT_MAC_EID_JAPANESE     = 1
	STBTT_MAC_EID_HEBREW       = 5
	STBTT_MAC_EID_CHINESE_TRAD = 2
	STBTT_MAC_EID_GREEK        = 6
	STBTT_MAC_EID_KOREAN       = 3
	STBTT_MAC_EID_RUSSIAN      = 7
)

const ( // languageID for STBTT_PLATFORM_ID_MICROSOFT; same as LCID...
	// problematic because there are e.g. 16 english LCIDs and 16 arabic LCIDs
	STBTT_MS_LANG_ENGLISH  = 0x0409
	STBTT_MS_LANG_ITALIAN  = 0x0410
	STBTT_MS_LANG_CHINESE  = 0x0804
	STBTT_MS_LANG_JAPANESE = 0x0411
	STBTT_MS_LANG_DUTCH    = 0x0413
	STBTT_MS_LANG_KOREAN   = 0x0412
	STBTT_MS_LANG_FRENCH   = 0x040c
	STBTT_MS_LANG_RUSSIAN  = 0x0419
	STBTT_MS_LANG_GERMAN   = 0x0407
	STBTT_MS_LANG_SPANISH  = 0x0409
	STBTT_MS_LANG_HEBREW   = 0x040d
	STBTT_MS_LANG_SWEDISH  = 0x041D
)

const ( // languageID for STBTT_PLATFORM_ID_MAC
	STBTT_MAC_LANG_ENGLISH            = 0
	STBTT_MAC_LANG_JAPANESE           = 11
	STBTT_MAC_LANG_ARABIC             = 12
	STBTT_MAC_LANG_KOREAN             = 23
	STBTT_MAC_LANG_DUTCH              = 4
	STBTT_MAC_LANG_RUSSIAN            = 32
	STBTT_MAC_LANG_FRENCH             = 1
	STBTT_MAC_LANG_SPANISH            = 6
	STBTT_MAC_LANG_GERMAN             = 2
	STBTT_MAC_LANG_SWEDISH            = 5
	STBTT_MAC_LANG_HEBREW             = 10
	STBTT_MAC_LANG_CHINESE_SIMPLIFIED = 33
	STBTT_MAC_LANG_ITALIAN            = 3
	STBTT_MAC_LANG_CHINESE_TRAD       = 19
)

const STBTT_MAX_OVERSAMPLE = 8
const STBTT_RASTERIZER_VERSION = 2

func stbtt__buf_get8(b *stbtt__buf) stbtt_uint8 {
	if b.cursor >= b.size {
		return 0
	}
	res := b.data[b.cursor]
	b.cursor++
	return res
}

func stbtt__buf_peek8(b *stbtt__buf) stbtt_uint8 {
	if b.cursor >= b.size {
		return 0
	}
	return b.data[b.cursor]
}

func stbtt__buf_seek(b *stbtt__buf, o int) {
	STBTT_assert(!(o > b.size || o < 0))
	if o > b.size || o < 0 {
		b.cursor = b.size
	} else {
		b.cursor = o
	}
}

func stbtt__buf_skip(b *stbtt__buf, o int) {
	stbtt__buf_seek(b, b.cursor+o)
}

func stbtt__buf_get(b *stbtt__buf, n int) stbtt_uint32 {
	var v stbtt_uint32 = 0
	var i int
	STBTT_assert(n >= 1 && n <= 4)
	for i = 0; i < n; i++ {
		v = (v << 8) | uint32(stbtt__buf_get8(b))
	}
	return v
}

func stbtt__new_buf(p []byte, size golang.Int) stbtt__buf {
	var r stbtt__buf
	STBTT_assert(size < 0x40000000)
	r.data = p
	r.size = int(size)
	r.cursor = 0
	return r
}

func stbtt__buf_get16(b *stbtt__buf) stbtt_uint16 {
	return stbtt_uint16(stbtt__buf_get(b, 2))
}

func stbtt__buf_get32(b *stbtt__buf) uint {
	return uint(stbtt__buf_get(b, 4))
}

func stbtt__buf_range(b *stbtt__buf, o, s int) stbtt__buf {
	var r = stbtt__new_buf(nil, 0)
	if o < 0 || s < 0 || o > b.size || s > b.size-o {
		return r
	}
	r.data = b.data[o:]
	r.size = s
	return r
}

func stbtt__cff_get_index(b *stbtt__buf) stbtt__buf {
	var count, start, offsize int
	start = b.cursor
	count = int(stbtt__buf_get16(b))
	if istrue(count) {
		offsize = int(stbtt__buf_get8(b))
		STBTT_assert(offsize >= 1 && offsize <= 4)
		stbtt__buf_skip(b, offsize*count)
		stbtt__buf_skip(b, int(stbtt__buf_get(b, offsize))-1)
	}
	return stbtt__buf_range(b, start, b.cursor-start)
}

func stbtt__cff_int(b *stbtt__buf) stbtt_uint32 {
	var b0 = int(stbtt__buf_get8(b))
	switch {
	case (b0 >= 32 && b0 <= 246):
		return stbtt_uint32(b0 - 139)
	case (b0 >= 247 && b0 <= 250):
		return stbtt_uint32((b0-247)*256 + int(stbtt__buf_get8(b)) + 108)
	case (b0 >= 251 && b0 <= 254):
		return stbtt_uint32(-(b0-251)*256 - int(stbtt__buf_get8(b)) - 108)
	case (b0 == 28):
		return stbtt_uint32(stbtt__buf_get16(b))
	case (b0 == 29):
		return stbtt_uint32(stbtt__buf_get32(b))
	}
	STBTT_assert(false)
	return 0
}

func stbtt__cff_skip_operand(b *stbtt__buf) {
	var v, b0 int = 0, int(stbtt__buf_peek8(b))
	STBTT_assert(b0 >= 28)
	if b0 == 30 {
		stbtt__buf_skip(b, 1)
		for b.cursor < b.size {
			v = int(stbtt__buf_get8(b))
			if (v&0xF) == 0xF || (v>>4) == 0xF {
				break
			}
		}
	} else {
		stbtt__cff_int(b)
	}
}

func stbtt__dict_get(b *stbtt__buf, key int) stbtt__buf {
	stbtt__buf_seek(b, 0)
	for b.cursor < b.size {
		var start, end, op int = b.cursor, 0, 0
		for stbtt__buf_peek8(b) >= 28 {
			stbtt__cff_skip_operand(b)
		}
		end = b.cursor
		op = int(stbtt__buf_get8(b))
		if op == 12 {
			//FIXME/TODO is this correct?
			op = int(stbtt__buf_get8(b)) | 0x100
		}
		if op == key {
			return stbtt__buf_range(b, start, end-start)
		}
	}
	return stbtt__buf_range(b, 0, 0)
}

func stbtt__dict_get_ints(b *stbtt__buf, key, outcount int, out []stbtt_uint32) {
	var i int
	var operands = stbtt__dict_get(b, key)
	for i = 0; i < outcount && operands.cursor < operands.size; i++ {
		out[i] = stbtt__cff_int(&operands)
	}
}

func stbtt__cff_index_count(b *stbtt__buf) int {
	stbtt__buf_seek(b, 0)
	return int(stbtt__buf_get16(b))
}

func stbtt__cff_index_get(b stbtt__buf, i int) stbtt__buf {
	var count, offsize, start, end int
	stbtt__buf_seek(&b, 0)
	count = int(stbtt__buf_get16(&b))
	offsize = int(stbtt__buf_get8(&b))
	STBTT_assert(i >= 0 && i < count)
	STBTT_assert(offsize >= 1 && offsize <= 4)
	stbtt__buf_skip(&b, i*offsize)
	start = int(stbtt__buf_get(&b, offsize))
	end = int(stbtt__buf_get(&b, offsize))
	return stbtt__buf_range(&b, 2+(count+1)*offsize+start, end-start)
}

func ttBYTE(p []byte) stbtt_uint8 {
	return p[0]
}

func ttCHAR(p []byte) stbtt_int8 {
	//FIXME/TODO is this correct?
	return *(*int8)(unsafe.Pointer(&p[0]))
}

func ttFIXED(p []byte) stbtt_int32 {
	return ttLONG(p)
}

func ttUSHORT(p []byte) stbtt_uint16 { return stbtt_uint16(p[0])*256 + stbtt_uint16(p[1]) }
func ttSHORT(p []byte) stbtt_int16   { return stbtt_int16(p[0])*256 + stbtt_int16(p[1]) }
func ttULONG(p []byte) stbtt_uint32 {
	return (stbtt_uint32(p[0]) << 24) + (stbtt_uint32(p[1]) << 16) + (stbtt_uint32(p[2]) << 8) + stbtt_uint32(p[3])
}
func ttLONG(p []byte) stbtt_int32 {
	return (stbtt_int32(p[0]) << 24) + (stbtt_int32(p[1]) << 16) + (stbtt_int32(p[2]) << 8) + stbtt_int32(p[3])
}

func stbtt_tag4(p []byte, c0, c1, c2, c3 byte) bool {
	return p[0] == c0 && p[1] == c1 && p[2] == c2 && p[3] == c3
}

func stbtt_tag(p []byte, str string) bool {
	return stbtt_tag4(p, str[0], str[1], str[2], str[3])
}

func stbtt__isfont(font []byte) int {
	// check the version number
	switch {
	case (stbtt_tag4(font, '1', 0, 0, 0)), // TrueType 1
		(stbtt_tag(font, "typ1")),      // TrueType with type 1 font -- we don't support this!
		(stbtt_tag(font, "OTTO")),      // OpenType with CFF
		(stbtt_tag4(font, 0, 1, 0, 0)), // OpenType 1.0
		(stbtt_tag(font, "true")):      // Apple specification for TrueType fonts
		return 1
	default:
		return 0
	}
}

// @OPTIMIZE: binary search
func stbtt__find_table(data []byte, fontstart stbtt_uint32, tag string) stbtt_uint32 {
	var num_tables = stbtt_int32(ttUSHORT(data[fontstart+4:]))
	var tabledir = fontstart + 12
	var i stbtt_int32
	for i = 0; i < num_tables; i++ {
		var loc = tabledir + stbtt_uint32(16*i)
		if stbtt_tag(data[loc+0:], tag) {
			return ttULONG(data[loc+8:])
		}
	}
	return 0
}

func stbtt_GetFontOffsetForIndex_internal(font_collection []byte, index int) int {
	// if it's just a font, there's only one valid index
	if istrue(stbtt__isfont(font_collection)) {
		if index == 0 {
			return 0
		}
		return -1
	}

	// check if it's a TTC
	if stbtt_tag(font_collection, "ttcf") {
		// version 1?
		if ttULONG(font_collection[4:]) == 0x00010000 || ttULONG(font_collection[4:]) == 0x00020000 {
			var n = ttLONG(font_collection[8:])
			if index >= n {
				return -1
			}
			return int(ttULONG(font_collection[12+index*4:]))
		}
	}
	return -1
}

func stbtt_GetNumberOfFonts_internal(font_collection []byte) int {
	// if it's just a font, there's only one valid font
	if stbtt__isfont(font_collection) != 0 {
		return 1
	}

	// check if it's a TTC
	if stbtt_tag(font_collection, "ttcf") {
		// version 1?
		if ttULONG(font_collection[4:]) == 0x00010000 || ttULONG(font_collection[4:]) == 0x00020000 {
			return ttLONG(font_collection[8:])
		}
	}
	return 0
}

func stbtt__get_subrs(cff stbtt__buf, fontdict stbtt__buf) stbtt__buf {
	var subrsoff [1]stbtt_uint32
	var private_loc = [2]stbtt_uint32{0, 0}
	var pdict stbtt__buf
	stbtt__dict_get_ints(&fontdict, 18, 2, private_loc[:])
	if isfalse(int(private_loc[1])) || isfalse(int(private_loc[0])) {
		return stbtt__new_buf(nil, 0)
	}
	pdict = stbtt__buf_range(&cff, int(private_loc[1]), int(private_loc[0]))
	stbtt__dict_get_ints(&pdict, 19, 1, subrsoff[:])
	if isfalse(int(subrsoff[0])) {
		return stbtt__new_buf(nil, 0)
	}
	stbtt__buf_seek(&cff, int(private_loc[1]+subrsoff[0]))
	return stbtt__cff_get_index(&cff)
}

func stbtt_InitFont_internal(info *FontInfo, data []byte, fontstart int) int {
	var cmap, t stbtt_uint32
	var i, numTables stbtt_int32

	info.data = data
	info.fontstart = fontstart
	info.cff = stbtt__new_buf(nil, 0)

	cmap = stbtt__find_table(data, uint(fontstart), "cmap")           // required
	info.loca = int(stbtt__find_table(data, uint(fontstart), "loca")) // required
	info.head = int(stbtt__find_table(data, uint(fontstart), "head")) // required
	info.glyf = int(stbtt__find_table(data, uint(fontstart), "glyf")) // required
	info.hhea = int(stbtt__find_table(data, uint(fontstart), "hhea")) // required
	info.hmtx = int(stbtt__find_table(data, uint(fontstart), "hmtx")) // required
	info.kern = int(stbtt__find_table(data, uint(fontstart), "kern")) // not required
	info.gpos = int(stbtt__find_table(data, uint(fontstart), "GPOS")) // not required

	if cmap == 0 || info.head == 0 || info.hhea == 0 || info.hmtx == 0 {
		return 0
	}

	if istrue(info.glyf) {
		// required for truetype
		if info.loca == 0 {
			return 0
		}
	} else {
		// initialization for CFF / Type2 fonts (OTF)
		var b, topdict, topdictidx stbtt__buf
		var cstype, charstrings, fdarrayoff, fdselectoff = [1]stbtt_uint32{2},
			[1]stbtt_uint32{0},
			[1]stbtt_uint32{0},
			[1]stbtt_uint32{0}

		var cff stbtt_uint32

		cff = stbtt__find_table(data, uint(fontstart), "CFF ")
		if cff == 0 {
			return 0
		}

		info.fontdicts = stbtt__new_buf(nil, 0)
		info.fdselect = stbtt__new_buf(nil, 0)

		// @TODO this should use size from table (not 512MB)
		info.cff = stbtt__new_buf(data[cff:], 512*1024*1024)
		b = info.cff

		// read the header
		stbtt__buf_skip(&b, 2)
		stbtt__buf_seek(&b, int(stbtt__buf_get8(&b))) // hdrsize

		// @TODO the name INDEX could list multiple fonts,
		// but we just use the first one.
		stbtt__cff_get_index(&b) // name INDEX
		topdictidx = stbtt__cff_get_index(&b)
		topdict = stbtt__cff_index_get(topdictidx, 0)
		stbtt__cff_get_index(&b) // string INDEX
		info.gsubrs = stbtt__cff_get_index(&b)

		stbtt__dict_get_ints(&topdict, 17, 1, charstrings[:])
		stbtt__dict_get_ints(&topdict, 0x100|6, 1, cstype[:])
		stbtt__dict_get_ints(&topdict, 0x100|36, 1, fdarrayoff[:])
		stbtt__dict_get_ints(&topdict, 0x100|37, 1, fdselectoff[:])
		info.subrs = stbtt__get_subrs(b, topdict)

		// we only support Type 2 charstrings
		if cstype[0] != 2 {
			return 0
		}
		if charstrings[0] == 0 {
			return 0
		}

		if fdarrayoff[0] != 0 {
			// looks like a CID font
			if fdselectoff[0] == 0 {
				return 0
			}
			stbtt__buf_seek(&b, int(fdarrayoff[0]))
			info.fontdicts = stbtt__cff_get_index(&b)
			info.fdselect = stbtt__buf_range(&b, int(fdselectoff[0]), b.size-int(fdselectoff[0]))
		}

		stbtt__buf_seek(&b, int(charstrings[0]))
		info.charstrings = stbtt__cff_get_index(&b)
	}

	t = stbtt__find_table(data, uint(fontstart), "maxp")
	if t != 0 {
		info.numGlyphs = int(ttUSHORT(data[t+4:]))
	} else {
		info.numGlyphs = 0xffff
	}

	// find a cmap encoding table we understand *now* to avoid searching
	// later. (todo: could make this installable)
	// the same regardless of glyph.
	numTables = int(ttUSHORT(data[cmap+2:]))
	info.index_map = 0
	for i = 0; i < numTables; i++ {
		var encoding_record = cmap + 4 + uint(8*i)
		// find an encoding we understand:
		switch ttUSHORT(data[encoding_record:]) {
		case STBTT_PLATFORM_ID_MICROSOFT:
			switch ttUSHORT(data[encoding_record+2:]) {
			case STBTT_MS_EID_UNICODE_BMP:
			case STBTT_MS_EID_UNICODE_FULL:
				// MS/Unicode
				info.index_map = int(cmap + ttULONG(data[encoding_record+4:]))
				break
			}
			break
		case STBTT_PLATFORM_ID_UNICODE:
			// Mac/iOS has these
			// all the encodingIDs are unicode, so we don't bother to check it
			info.index_map = int(cmap + ttULONG(data[encoding_record+4:]))
			break
		}
	}
	if info.index_map == 0 {
		return 0
	}

	info.indexToLocFormat = int(ttUSHORT(data[info.head+50:]))
	return 1
}

func stbtt_setvertex(v *stbtt_vertex, vtype stbtt_uint8, x, y, cx, cy stbtt_int32) {
	v.vtype = vtype
	v.x = (stbtt_int16)(x)
	v.y = (stbtt_int16)(y)
	v.cx = (stbtt_int16)(cx)
	v.cy = (stbtt_int16)(cy)
}

func stbtt__GetGlyfOffset(info *FontInfo, glyph_index int) int {
	var g1, g2 int

	STBTT_assert(info.cff.size == 0)

	if glyph_index >= info.numGlyphs {
		return -1 // glyph index out of range
	}
	if info.indexToLocFormat >= 2 {
		return -1 // unknown index.glyph map format
	}

	if info.indexToLocFormat == 0 {
		g1 = info.glyf + int(ttUSHORT(info.data[info.loca+glyph_index*2:])*2)
		g2 = info.glyf + int(ttUSHORT(info.data[info.loca+glyph_index*2+2:])*2)
	} else {
		g1 = info.glyf + int(ttULONG(info.data[info.loca+glyph_index*4:]))
		g2 = info.glyf + int(ttULONG(info.data[info.loca+glyph_index*4+4:]))
	}

	if g1 == g2 {
		return -1 // if length is 0, return -1
	}
	return g1
}

func stbtt__GetGlyphInfoT2(info *FontInfo, glyph_index int, x0, y0, x1, y1 *int) int {
	panic("unimplemented")
}

func stbtt__close_shape(vertices []stbtt_vertex, num_vertices, was_off, start_off int,
	sx, sy, scx, scy, cx, cy stbtt_int32) int {
	if start_off != 0 {
		if was_off != 0 {
			stbtt_setvertex(&vertices[num_vertices], STBTT_vcurve, (cx+scx)>>1, (cy+scy)>>1, cx, cy)
			num_vertices++
		}
		stbtt_setvertex(&vertices[num_vertices], STBTT_vcurve, sx, sy, scx, scy)
		num_vertices++
	} else {
		if was_off != 0 {
			stbtt_setvertex(&vertices[num_vertices], STBTT_vcurve, sx, sy, cx, cy)
			num_vertices++
		} else {
			stbtt_setvertex(&vertices[num_vertices], STBTT_vline, sx, sy, 0, 0)
			num_vertices++
		}
	}
	return num_vertices
}

func stbtt__GetGlyphShapeTT(info *FontInfo, glyph_index int, pvertices *[]stbtt_vertex) int {
	var numberOfContours stbtt_int16
	var endPtsOfContours []stbtt_uint8
	var data = info.data
	var vertices []stbtt_vertex

	var num_vertices int = 0
	var g = stbtt__GetGlyfOffset(info, glyph_index)

	*pvertices = nil

	if g < 0 {
		return 0
	}

	numberOfContours = ttSHORT(data[g:])

	if numberOfContours > 0 {
		var flags, flagcount stbtt_uint8 = 0, 0
		var ins, i, j, m, n, next_move, was_off, off, start_off stbtt_int32
		var x, y, cx, cy, sx, sy, scx, scy stbtt_int32
		var points []stbtt_uint8
		endPtsOfContours = (data[g+10:])
		ins = int(ttUSHORT(data[g+10+int(numberOfContours)*2:]))
		points = data[g+10+int(numberOfContours)*2+2+ins:]

		n = int(1 + ttUSHORT(endPtsOfContours[numberOfContours*2-2:]))

		m = n + 2*int(numberOfContours) // a loose bound on how many vertices we might need
		vertices = make([]stbtt_vertex, m)
		if vertices == nil {
			return 0
		}

		next_move = 0
		flagcount = 0

		// in first pass, we load uninterpreted data into the allocated array
		// above, shifted to the end of the array so we won't overwrite it when
		// we create our final data starting from the front

		off = m - n // starting offset for uninterpreted data, regardless of how m ends up being calculated

		// first load flags

		for i = 0; i < n; i++ {
			if flagcount == 0 {
				flags = points[0]
				points = points[1:]
				if flags&8 != 0 {
					flagcount = points[0]
					points = points[1:]
				}
			} else {
				flagcount--
			}
			vertices[off+i].vtype = flags
		}

		// now load x coordinates
		x = 0
		for i = 0; i < n; i++ {
			flags = vertices[off+i].vtype
			if flags&2 != 0 {
				var dx = stbtt_int16(points[0])
				points = points[1:]
				if (flags & 16) != 0 {
					x += int(dx)
				} else {
					x -= int(dx)
				}
			} else {
				if flags&16 == 0 {
					x = x + int(stbtt_int16(points[0])*256+stbtt_int16(points[1]))
					points = points[2:]
				}
			}
			vertices[off+i].x = (stbtt_int16)(x)
		}

		// now load y coordinates
		y = 0
		for i = 0; i < n; i++ {
			flags = vertices[off+i].vtype
			if flags&4 != 0 {
				var dy = stbtt_int16(points[0])
				points = points[1:]

				if (flags & 32) != 0 {
					y += int(dy)
				} else {
					y -= int(dy)
				}
			} else {
				if flags&32 == 0 {
					y = y + int(stbtt_int16(points[0])*256+stbtt_int16(points[1]))
					points = points[2:]
				}
			}
			vertices[off+i].y = (stbtt_int16)(y)
		}

		num_vertices = 0
		sx = 0
		sy = 0
		cx = 0
		cy = 0
		scx = 0
		scy = 0
		// now convert them to our format
		for i = 0; i < n; i++ {
			flags = vertices[off+i].vtype
			x = (int)(vertices[off+i].x)
			y = (int)(vertices[off+i].y)

			if next_move == i {
				if i != 0 {
					num_vertices = stbtt__close_shape(vertices, num_vertices, was_off, start_off, sx, sy, scx, scy, cx, cy)
				}

				// now start the new one
				start_off = bool2int((flags & 1) == 0)
				if start_off != 0 {
					// if we start off with an off-curve point, then when we need to find a point on the curve
					// where we can start, and we need to save some state for when we wraparound.
					scx = x
					scy = y
					if (vertices[off+i+1].vtype & 1) == 0 {
						// next point is also a curve point, so interpolate an on-point curve
						sx = (x + (stbtt_int32)(vertices[off+i+1].x)) >> 1
						sy = (y + (stbtt_int32)(vertices[off+i+1].y)) >> 1
					} else {
						// otherwise just use the next point as our start point
						sx = (stbtt_int32)(vertices[off+i+1].x)
						sy = (stbtt_int32)(vertices[off+i+1].y)
						i++ // we're using point i+1 as the starting point, so skip it
					}
				} else {
					sx = x
					sy = y
				}
				stbtt_setvertex(&vertices[num_vertices], STBTT_vmove, sx, sy, 0, 0)
				num_vertices++
				was_off = 0
				next_move = 1 + int(ttUSHORT(endPtsOfContours[j*2:]))
				j++
			} else {
				if flags&1 == 0 { // if it's a curve
					if was_off != 0 { // two off-curve control points in a row means interpolate an on-curve midpoint
						stbtt_setvertex(&vertices[num_vertices], STBTT_vcurve, (cx+x)>>1, (cy+y)>>1, cx, cy)
						num_vertices++
					}
					cx = x
					cy = y
					was_off = 1
				} else {
					if was_off != 0 {
						stbtt_setvertex(&vertices[num_vertices], STBTT_vcurve, x, y, cx, cy)
						num_vertices++
					} else {
						stbtt_setvertex(&vertices[num_vertices], STBTT_vline, x, y, 0, 0)
						num_vertices++
					}
					was_off = 0
				}
			}
		}
		num_vertices = stbtt__close_shape(vertices, num_vertices, was_off, start_off, sx, sy, scx, scy, cx, cy)
	} else if numberOfContours == -1 {
		// Compound shapes.
		var more int = 1
		var comp = data[g+10:]
		num_vertices = 0
		vertices = nil
		for more != 0 {
			var flags, gidx stbtt_uint16
			var comp_num_verts, i int
			var comp_verts, tmp []stbtt_vertex
			var mtx = [6]float{1, 0, 0, 1, 0, 0}
			var m, n float

			flags = uint16(ttSHORT(comp))
			comp = comp[2:]
			gidx = uint16(ttSHORT(comp))
			comp = comp[2:]

			if flags&2 != 0 { // XY values
				if flags&1 != 0 { // shorts
					mtx[4] = float(ttSHORT(comp))
					comp = comp[2:]
					mtx[5] = float(ttSHORT(comp))
					comp = comp[2:]
				} else {
					mtx[4] = float(ttCHAR(comp))
					comp = comp[1:]
					mtx[5] = float(ttCHAR(comp))
					comp = comp[1:]
				}
			} else {
				// @TODO handle matching point
				STBTT_assert(false)
			}
			if flags&(1<<3) != 0 { // WE_HAVE_A_SCALE
				mtx[0] = float(ttSHORT(comp)) / 16384.0
				mtx[3] = mtx[0]
				comp = comp[2:]
				mtx[1] = 0
				mtx[2] = 0
			} else if flags&(1<<6) != 0 { // WE_HAVE_AN_X_AND_YSCALE
				mtx[0] = float(ttSHORT(comp)) / 16384.0
				comp = comp[2:]
				mtx[1] = 0
				mtx[2] = 0
				mtx[3] = float(ttSHORT(comp)) / 16384.0
				comp = comp[2:]
			} else if flags&(1<<7) != 0 { // WE_HAVE_A_TWO_BY_TWO
				mtx[0] = float(ttSHORT(comp)) / 16384.0
				comp = comp[2:]
				mtx[1] = float(ttSHORT(comp)) / 16384.0
				comp = comp[2:]
				mtx[2] = float(ttSHORT(comp)) / 16384.0
				comp = comp[2:]
				mtx[3] = float(ttSHORT(comp)) / 16384.0
				comp = comp[2:]
			}

			// Find transformation scales.
			m = STBTT_sqrt(mtx[0]*mtx[0] + mtx[1]*mtx[1])
			n = STBTT_sqrt(mtx[2]*mtx[2] + mtx[3]*mtx[3])

			// Get indexed glyph.
			comp_num_verts = stbtt_GetGlyphShape(info, int(gidx), &comp_verts)
			if comp_num_verts > 0 {
				// Transform vertices.
				for i = 0; i < comp_num_verts; i++ {
					var v = &comp_verts[i]
					var x, y stbtt_vertex_type
					x = v.x
					y = v.y
					v.x = (stbtt_vertex_type)(m * (mtx[0]*float32(x) + mtx[2]*float32(y) + mtx[4]))
					v.y = (stbtt_vertex_type)(n * (mtx[1]*float32(x) + mtx[3]*float32(y) + mtx[5]))
					x = v.cx
					y = v.cy
					v.cx = (stbtt_vertex_type)(m * (mtx[0]*float32(x) + mtx[2]*float32(y) + mtx[4]))
					v.cy = (stbtt_vertex_type)(n * (mtx[1]*float32(x) + mtx[3]*float32(y) + mtx[5]))
				}
				// Append vertices.
				tmp = make([]stbtt_vertex, num_vertices+comp_num_verts)
				if tmp == nil {
					if vertices != nil {
						//STBTT_free(vertices, info.userdata)
					}
					if comp_verts != nil {
						//STBTT_free(comp_verts, info.userdata)
					}
					return 0
				}
				if num_vertices > 0 {
					copy(tmp, vertices[:num_vertices]) //-V595
				}
				copy(tmp[num_vertices:], comp_verts[:comp_num_verts])
				if vertices != nil {
					//STBTT_free(vertices, info.userdata)
				}
				vertices = tmp
				//STBTT_free(comp_verts, info.userdata)
				num_vertices += comp_num_verts
			}
			// More components ?
			more = int(flags & (1 << 5))
		}
	} else if numberOfContours < 0 {
		// @TODO other compound variations?
		STBTT_assert(false)
	} else {
		// numberOfCounters == 0, do nothing
	}

	*pvertices = vertices[:num_vertices]

	return num_vertices
}

type stbtt__csctx struct {
	bounds                     int
	started                    int
	first_x, first_y           float
	x, y                       float
	min_x, max_x, min_y, max_y stbtt_int32

	pvertices    []stbtt_vertex
	num_vertices int
}

func stbtt__track_vertex(c *stbtt__csctx, x, y stbtt_int32) {
	if x > c.max_x || c.started == 0 {
		c.max_x = x
	}
	if y > c.max_y || c.started == 0 {
		c.max_y = y
	}
	if x < c.min_x || c.started == 0 {
		c.min_x = x
	}
	if y < c.min_y || c.started == 0 {
		c.min_y = y
	}
	c.started = 1
}

func stbtt__csctx_v(c *stbtt__csctx, vtype stbtt_uint8, x, y, cx, cy, cx1, cy1 stbtt_int32) {
	if c.bounds != 0 {
		stbtt__track_vertex(c, x, y)
		if vtype == STBTT_vcubic {
			stbtt__track_vertex(c, cx, cy)
			stbtt__track_vertex(c, cx1, cy1)
		}
	} else {
		stbtt_setvertex(&c.pvertices[c.num_vertices], vtype, x, y, cx, cy)
		c.pvertices[c.num_vertices].cx1 = (stbtt_int16)(cx1)
		c.pvertices[c.num_vertices].cy1 = (stbtt_int16)(cy1)
	}
	c.num_vertices++
}

func stbtt__csctx_close_shape(ctx *stbtt__csctx) {
	if ctx.first_x != ctx.x || ctx.first_y != ctx.y {
		stbtt__csctx_v(ctx, STBTT_vline, (int)(ctx.first_x), (int)(ctx.first_y), 0, 0, 0, 0)
	}
}

func stbtt__csctx_rmove_to(ctx *stbtt__csctx, dx, dy float) {
	stbtt__csctx_close_shape(ctx)
	//TODO/FIXME check that this is correct.
	ctx.first_x, ctx.x = ctx.x+dx, ctx.x+dx
	ctx.first_y, ctx.y = ctx.y+dy, ctx.y+dy
	stbtt__csctx_v(ctx, STBTT_vmove, (int)(ctx.x), (int)(ctx.y), 0, 0, 0, 0)
}

func stbtt__csctx_rline_to(ctx *stbtt__csctx, dx, dy float) {
	ctx.x += dx
	ctx.y += dy
	stbtt__csctx_v(ctx, STBTT_vline, (int)(ctx.x), (int)(ctx.y), 0, 0, 0, 0)
}

func stbtt__csctx_rccurve_to(ctx *stbtt__csctx, dx1, dy1, dx2, dy2, dx3, dy3 float) {
	var cx1 = ctx.x + dx1
	var cy1 = ctx.y + dy1
	var cx2 = cx1 + dx2
	var cy2 = cy1 + dy2
	ctx.x = cx2 + dx3
	ctx.y = cy2 + dy3
	stbtt__csctx_v(ctx, STBTT_vcubic, (int)(ctx.x), (int)(ctx.y), (int)(cx1), (int)(cy1), (int)(cx2), (int)(cy2))
}

func stbtt__get_subr(idx stbtt__buf, n int) stbtt__buf {
	var count = stbtt__cff_index_count(&idx)
	var bias int = 107
	if count >= 33900 {
		bias = 32768
	} else if count >= 1240 {
		bias = 1131
	}
	n += bias
	if n < 0 || n >= count {
		return stbtt__new_buf(nil, 0)
	}
	return stbtt__cff_index_get(idx, n)
}

func stbtt__cid_get_glyph_subrs(info *FontInfo, glyph_index int) stbtt__buf {
	var fdselect = info.fdselect
	var nranges, start, end, v, fmt, fdselector, i int = 0, 0, 0, 0, 0, -1, 0

	stbtt__buf_seek(&fdselect, 0)
	fmt = int(stbtt__buf_get8(&fdselect))
	if fmt == 0 {
		// untested
		stbtt__buf_skip(&fdselect, glyph_index)
		fdselector = int(stbtt__buf_get8(&fdselect))
	} else if fmt == 3 {
		nranges = int(stbtt__buf_get16(&fdselect))
		start = int(stbtt__buf_get16(&fdselect))
		for i = 0; i < nranges; i++ {
			v = int(stbtt__buf_get8(&fdselect))
			end = int(stbtt__buf_get16(&fdselect))
			if glyph_index >= start && glyph_index < end {
				fdselector = v
				break
			}
			start = end
		}
	}
	if fdselector == -1 {
		stbtt__new_buf(nil, 0)
	}
	return stbtt__get_subrs(info.cff, stbtt__cff_index_get(info.fontdicts, fdselector))
}

func stbtt__run_charstring(info *FontInfo, glyph_index int, c *stbtt__csctx) int {
	var in_header, maskbits, subr_stack_height, sp, v, i, b0 int = 1, 0, 0, 0, 0, 0, 0
	var has_subrs, clear_stack int
	var s [48]float
	var subr_stack [10]stbtt__buf
	var subrs = info.subrs
	var b stbtt__buf
	var f float

	STBTT__CSERR := func(s string) int {
		return 0
	}

	// this currently ignores the initial width value, which isn't needed if we have hmtx
	b = stbtt__cff_index_get(info.charstrings, glyph_index)
	for b.cursor < b.size {
		i = 0
		clear_stack = 1
		b0 = int(stbtt__buf_get8(&b))
		switch b0 {
		// @TODO implement hinting
		case 0x13: // hintmask
			fallthrough
		case 0x14: // cntrmask
			if in_header != 0 {
				maskbits += (sp / 2) // implicit "vstem"
			}
			in_header = 0
			stbtt__buf_skip(&b, (maskbits+7)/8)
			break

		case 0x01: // hstem
			fallthrough
		case 0x03: // vstem
			fallthrough
		case 0x12: // hstemhm
			fallthrough
		case 0x17: // vstemhm
			maskbits += (sp / 2)
			break

		case 0x15: // rmoveto
			in_header = 0
			if sp < 2 {
				return STBTT__CSERR("rmoveto stack")
			}
			stbtt__csctx_rmove_to(c, s[sp-2], s[sp-1])
			break
		case 0x04: // vmoveto
			in_header = 0
			if sp < 1 {
				return STBTT__CSERR("vmoveto stack")
			}
			stbtt__csctx_rmove_to(c, 0, s[sp-1])
			break
		case 0x16: // hmoveto
			in_header = 0
			if sp < 1 {
				return STBTT__CSERR("hmoveto stack")
			}
			stbtt__csctx_rmove_to(c, s[sp-1], 0)
			break

		case 0x05: // rlineto
			if sp < 2 {
				return STBTT__CSERR("rlineto stack")
			}
			for ; i+1 < sp; i += 2 {
				stbtt__csctx_rline_to(c, s[i], s[i+1])
			}
			break

		// hlineto/vlineto and vhcurveto/hvcurveto alternate horizontal and vertical
		// starting from a different place.

		case 0x07: // vlineto
			if sp < 1 {
				return STBTT__CSERR("vlineto stack")
			}
			if i >= sp {
				break
			}
			stbtt__csctx_rline_to(c, 0, s[i])
			i++
			fallthrough
		case 0x06: // hlineto
			if sp < 1 {
				return STBTT__CSERR("hlineto stack")
			}
			for {
				if i >= sp {
					break
				}
				stbtt__csctx_rline_to(c, s[i], 0)
				i++
				if i >= sp {
					break
				}
				stbtt__csctx_rline_to(c, 0, s[i])
				i++
			}
			break

		case 0x1F: // hvcurveto
			if sp < 4 {
				return STBTT__CSERR("hvcurveto stack")
			}
			if i+3 >= sp {
				break
			}
			var last float32
			if sp-i == 5 {
				last = 0
			} else {
				last = s[i+3]
			}
			stbtt__csctx_rccurve_to(c, s[i], 0, s[i+1], s[i+2], last, s[i+3])
			i += 4
			fallthrough
		case 0x1E: // vhcurveto
			if sp < 4 {
				return STBTT__CSERR("vhcurveto stack")
			}
			for {
				if i+3 >= sp {
					break
				}
				var last float32
				if sp-i == 5 {
					last = s[i+4]
				}
				stbtt__csctx_rccurve_to(c, 0, s[i], s[i+1], s[i+2], s[i+3], last)
				i += 4
				if i+3 >= sp {
					break
				}
				if sp-i == 5 {
					last = 0
				} else {
					last = s[i+3]
				}
				stbtt__csctx_rccurve_to(c, s[i], 0, s[i+1], s[i+2], last, s[i+3])
				i += 4
			}
			break

		case 0x08: // rrcurveto
			if sp < 6 {
				return STBTT__CSERR("rcurveline stack")
			}
			for ; i+5 < sp; i += 6 {
				stbtt__csctx_rccurve_to(c, s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5])
			}
			break

		case 0x18: // rcurveline
			if sp < 8 {
				return STBTT__CSERR("rcurveline stack")
			}
			for ; i+5 < sp-2; i += 6 {
				stbtt__csctx_rccurve_to(c, s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5])
			}
			if i+1 >= sp {
				return STBTT__CSERR("rcurveline stack")
			}
			stbtt__csctx_rline_to(c, s[i], s[i+1])
			break

		case 0x19: // rlinecurve
			if sp < 8 {
				return STBTT__CSERR("rlinecurve stack")
			}
			for ; i+1 < sp-6; i += 2 {
				stbtt__csctx_rline_to(c, s[i], s[i+1])
			}
			if i+5 >= sp {
				return STBTT__CSERR("rlinecurve stack")
			}
			stbtt__csctx_rccurve_to(c, s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5])
			break

		case 0x1A: // vvcurveto
			fallthrough
		case 0x1B: // hhcurveto
			if sp < 4 {
				return STBTT__CSERR("(vv|hh)curveto stack")
			}
			f = 0.0
			if sp&1 != 0 {
				f = s[i]
				i++
			}
			for ; i+3 < sp; i += 4 {
				if b0 == 0x1B {
					stbtt__csctx_rccurve_to(c, s[i], f, s[i+1], s[i+2], s[i+3], 0.0)
				} else {
					stbtt__csctx_rccurve_to(c, f, s[i], s[i+1], s[i+2], 0.0, s[i+3])
				}
				f = 0.0
			}
			break

		case 0x0A: // callsubr
			if has_subrs == 0 {
				if info.fdselect.size != 0 {
					subrs = stbtt__cid_get_glyph_subrs(info, glyph_index)
				}
				has_subrs = 1
			}
			// fallthrough
		case 0x1D: // callgsubr
			if sp < 1 {
				return STBTT__CSERR("call(g|)subr stack")
			}
			sp--
			v = (int)(s[sp])
			if subr_stack_height >= 10 {
				return STBTT__CSERR("recursion limit")
			}

			subr_stack[subr_stack_height] = b
			subr_stack_height++

			var subr stbtt__buf
			if b0 == 0x0A {
				subr = subrs
			} else {
				subr = info.gsubrs
			}

			b = stbtt__get_subr(subr, v)
			if b.size == 0 {
				return STBTT__CSERR("subr not found")
			}
			b.cursor = 0
			clear_stack = 0
			break

		case 0x0B: // return
			if subr_stack_height <= 0 {
				return STBTT__CSERR("return outside subr")
			}
			subr_stack_height--
			b = subr_stack[subr_stack_height]
			clear_stack = 0
			break

		case 0x0E: // endchar
			stbtt__csctx_close_shape(c)
			return 1

		case 0x0C:
			{ // two-byte escape
				var dx1, dx2, dx3, dx4, dx5, dx6, dy1, dy2, dy3, dy4, dy5, dy6 float
				var dx, dy float
				var b1 = int(stbtt__buf_get8(&b))
				switch b1 {
				// @TODO These "flex" implementations ignore the flex-depth and resolution,
				// and always draw beziers.
				case 0x22: // hflex
					if sp < 7 {
						return STBTT__CSERR("hflex stack")
					}
					dx1 = s[0]
					dx2 = s[1]
					dy2 = s[2]
					dx3 = s[3]
					dx4 = s[4]
					dx5 = s[5]
					dx6 = s[6]
					stbtt__csctx_rccurve_to(c, dx1, 0, dx2, dy2, dx3, 0)
					stbtt__csctx_rccurve_to(c, dx4, 0, dx5, -dy2, dx6, 0)
					break

				case 0x23: // flex
					if sp < 13 {
						return STBTT__CSERR("flex stack")
					}
					dx1 = s[0]
					dy1 = s[1]
					dx2 = s[2]
					dy2 = s[3]
					dx3 = s[4]
					dy3 = s[5]
					dx4 = s[6]
					dy4 = s[7]
					dx5 = s[8]
					dy5 = s[9]
					dx6 = s[10]
					dy6 = s[11]
					//fd is s[12]
					stbtt__csctx_rccurve_to(c, dx1, dy1, dx2, dy2, dx3, dy3)
					stbtt__csctx_rccurve_to(c, dx4, dy4, dx5, dy5, dx6, dy6)
					break

				case 0x24: // hflex1
					if sp < 9 {
						return STBTT__CSERR("hflex1 stack")
					}
					dx1 = s[0]
					dy1 = s[1]
					dx2 = s[2]
					dy2 = s[3]
					dx3 = s[4]
					dx4 = s[5]
					dx5 = s[6]
					dy5 = s[7]
					dx6 = s[8]
					stbtt__csctx_rccurve_to(c, dx1, dy1, dx2, dy2, dx3, 0)
					stbtt__csctx_rccurve_to(c, dx4, 0, dx5, dy5, dx6, -(dy1 + dy2 + dy5))
					break

				case 0x25: // flex1
					if sp < 11 {
						return STBTT__CSERR("flex1 stack")
					}
					dx1 = s[0]
					dy1 = s[1]
					dx2 = s[2]
					dy2 = s[3]
					dx3 = s[4]
					dy3 = s[5]
					dx4 = s[6]
					dy4 = s[7]
					dx5 = s[8]
					dy5 = s[9]
					dx6 = s[10]
					dy6 = s[10]
					dx = dx1 + dx2 + dx3 + dx4 + dx5
					dy = dy1 + dy2 + dy3 + dy4 + dy5
					if STBTT_fabs(dx) > STBTT_fabs(dy) {
						dy6 = -dy
					} else {
						dx6 = -dx
					}
					stbtt__csctx_rccurve_to(c, dx1, dy1, dx2, dy2, dx3, dy3)
					stbtt__csctx_rccurve_to(c, dx4, dy4, dx5, dy5, dx6, dy6)
					break

				default:
					return STBTT__CSERR("unimplemented")
				}
			}
			break

		default:
			if b0 != 255 && b0 != 28 && (b0 < 32 || b0 > 254) { //-V560
				return STBTT__CSERR("reserved operator")
			}

			// push immediate
			if b0 == 255 {
				f = (float)((stbtt_int32)(stbtt__buf_get32(&b))) / 0x10000
			} else {
				stbtt__buf_skip(&b, -1)
				f = (float)((stbtt_int16)(stbtt__cff_int(&b)))
			}
			if sp >= 48 {
				return STBTT__CSERR("push stack overflow")
			}
			s[sp] = f
			sp++
			clear_stack = 0
			break
		}
		if clear_stack != 0 {
			sp = 0
		}
	}
	return STBTT__CSERR("no endchar")
}

func stbtt__GetGlyphShapeT2(info *FontInfo, glyph_index int, pvertices *[]stbtt_vertex) int {
	// runs the charstring twice, once to count and once to output (to avoid realloc)
	var count_ctx = stbtt__csctx{bounds: 1}
	var output_ctx = stbtt__csctx{bounds: 0}
	if stbtt__run_charstring(info, glyph_index, &count_ctx) != 0 {
		*pvertices = make([]stbtt_vertex, count_ctx.num_vertices)
		output_ctx.pvertices = *pvertices
		if stbtt__run_charstring(info, glyph_index, &output_ctx) != 0 {
			STBTT_assert(output_ctx.num_vertices == count_ctx.num_vertices)
			return output_ctx.num_vertices
		}
	}
	*pvertices = nil
	return 0
}

func stbtt_GetGlyphShape(info *FontInfo, glyph_index int, pvertices *[]stbtt_vertex) int {
	if info.cff.size == 0 {
		return stbtt__GetGlyphShapeTT(info, glyph_index, pvertices)
	} else {
		return stbtt__GetGlyphShapeT2(info, glyph_index, pvertices)
	}
}

func stbtt__GetGlyphKernInfoAdvance(info *FontInfo, glyph1, glyph2 int) int {
	var data = info.data[info.kern:]
	var needle, straw stbtt_uint32
	var l, r, m int

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

	l = 0
	r = int(ttUSHORT(data[10:]) - 1)
	needle = stbtt_uint32(glyph1)<<16 | stbtt_uint32(glyph2)
	for l <= r {
		m = (l + r) >> 1
		straw = ttULONG(data[18+(m*6):]) // note: unaligned read
		if needle < straw {
			r = m - 1
		} else if needle > straw {
			l = m + 1
		} else {
			return int(ttSHORT(data[22+(m*6):]))
		}
	}
	return 0
}

func stbtt__GetCoverageIndex(coverageTable []byte, glyph int) stbtt_int32 {
	var coverageFormat = ttUSHORT(coverageTable)
	switch coverageFormat {
	case 1:
		{
			var glyphCount = ttUSHORT(coverageTable[2:])

			// Binary search.
			var l, r, m stbtt_int32 = 0, stbtt_int32(glyphCount) - 1, 0
			var straw, needle int = 0, glyph
			for l <= r {
				var glyphArray = coverageTable[4:]
				var glyphID stbtt_uint16
				m = (l + r) >> 1
				glyphID = ttUSHORT(glyphArray[2*m:])
				straw = int(glyphID)
				if needle < straw {
					r = m - 1
				} else if needle > straw {
					l = m + 1
				} else {
					return m
				}
			}
		}
		break

	case 2:
		{
			var rangeCount = ttUSHORT(coverageTable[2:])
			var rangeArray = coverageTable[4:]

			// Binary search.
			var l, r, m stbtt_int32 = 0, stbtt_int32(rangeCount) - 1, 0
			var strawStart, strawEnd, needle int = 0, 0, glyph
			for l <= r {
				var rangeRecord []byte
				m = (l + r) >> 1
				rangeRecord = rangeArray[6*m:]
				strawStart = stbtt_int32(ttUSHORT(rangeRecord))
				strawEnd = stbtt_int32(ttUSHORT(rangeRecord[2:]))
				if needle < strawStart {
					r = m - 1
				} else if needle > strawEnd {
					l = m + 1
				} else {
					var startCoverageIndex = ttUSHORT(rangeRecord[4:])
					return int(startCoverageIndex) + glyph - strawStart
				}
			}
		}
		break

	default:
		{
			// There are no other cases.
			STBTT_assert(false)
		}
		break
	}

	return -1
}

func stbtt__GetGlyphClass(classDefTable []byte, glyph int) stbtt_int32 {
	var classDefFormat = ttUSHORT(classDefTable)
	switch classDefFormat {
	case 1:
		{
			var startGlyphID = ttUSHORT(classDefTable[2:])
			var glyphCount = ttUSHORT(classDefTable[4:])
			var classDef1ValueArray = classDefTable[6:]

			if glyph >= int(startGlyphID) && glyph < int(startGlyphID+glyphCount) {
				return (stbtt_int32)(ttUSHORT(classDef1ValueArray[2*(glyph-int(startGlyphID)):]))
			}

			// [DEAR IMGUI] Commented to fix static analyzer warning
			//classDefTable = classDef1ValueArray + 2 * glyphCount;
		}
		break

	case 2:
		{
			var classRangeCount = ttUSHORT(classDefTable[2:])
			var classRangeRecords = classDefTable[4:]

			// Binary search.
			var l, r, m stbtt_int32 = 0, int(classRangeCount) - 1, 0
			var strawStart, strawEnd, needle int = 0, 0, glyph
			for l <= r {
				var classRangeRecord []byte
				m = (l + r) >> 1
				classRangeRecord = classRangeRecords[6*m:]
				strawStart = int(ttUSHORT(classRangeRecord))
				strawEnd = int(ttUSHORT(classRangeRecord[2:]))
				if needle < strawStart {
					r = m - 1
				} else if needle > strawEnd {
					l = m + 1
				} else {
					return (stbtt_int32)(ttUSHORT(classRangeRecord[4:]))
				}
			}

			// [DEAR IMGUI] Commented to fix static analyzer warning
			//classDefTable = classRangeRecords + 6 * classRangeCount;
		}
		break

	default:
		{
			// There are no other cases.
			STBTT_assert(false)
		}
		break
	}

	return -1
}

func STBTT_GPOS_TODO_assert(x bool) {
	if !x {
		panic("not implemented")
	}
}

func stbtt__GetGlyphGPOSInfoAdvance(info *FontInfo, glyph1, glyph2 int) stbtt_int32 {
	var lookupListOffset stbtt_uint16
	var lookupList []byte
	var lookupCount stbtt_uint16
	var data []byte
	var i stbtt_int32

	if info.gpos == 0 {
		return 0
	}

	data = info.data[info.gpos:]

	if ttUSHORT(data[0:]) != 1 {
		return 0 // Major version 1
	}
	if ttUSHORT(data[2:]) != 0 {
		return 0 // Minor version 0
	}

	lookupListOffset = ttUSHORT(data[8:])
	lookupList = data[lookupListOffset:]
	lookupCount = ttUSHORT(lookupList)

	for i = 0; i < int(lookupCount); i++ {
		var lookupOffset = ttUSHORT(lookupList[+2+2*i:])
		var lookupTable = lookupList[lookupOffset:]

		var lookupType = ttUSHORT(lookupTable)
		var subTableCount = ttUSHORT(lookupTable[4:])
		var subTableOffsets = lookupTable[6:]
		switch lookupType {
		case 2:
			{ // Pair Adjustment Positioning Subtable
				var sti stbtt_int32
				for sti = 0; sti < int(subTableCount); sti++ {
					var subtableOffset = ttUSHORT(subTableOffsets[2*sti:])
					var table = lookupTable[subtableOffset:]
					var posFormat = ttUSHORT(table)
					var coverageOffset = ttUSHORT(table[2:])
					var coverageIndex = stbtt__GetCoverageIndex(table[coverageOffset:], glyph1)
					if coverageIndex == -1 {
						continue
					}

					switch posFormat {
					case 1:
						{
							var l, r, m stbtt_int32
							var straw, needle int
							var valueFormat1 = ttUSHORT(table[4:])
							var valueFormat2 = ttUSHORT(table[6:])
							var valueRecordPairSizeInBytes stbtt_int32 = 2
							var pairSetCount = ttUSHORT(table[8:])
							var pairPosOffset = ttUSHORT(table[10+2*coverageIndex:])
							var pairValueTable = table[pairPosOffset:]
							var pairValueCount = ttUSHORT(pairValueTable)
							var pairValueArray = pairValueTable[2:]
							// TODO: Support more formats.
							STBTT_GPOS_TODO_assert(valueFormat1 == 4)
							if valueFormat1 != 4 {
								return 0
							}
							STBTT_GPOS_TODO_assert(valueFormat2 == 0)
							if valueFormat2 != 0 {
								return 0
							}

							STBTT_assert(coverageIndex < int(pairSetCount))

							needle = glyph2
							r = int(pairValueCount) - 1
							l = 0

							// Binary search.
							for l <= r {
								var secondGlyph stbtt_uint16
								var pairValue []byte
								m = (l + r) >> 1
								pairValue = pairValueArray[(2+valueRecordPairSizeInBytes)*m:]
								secondGlyph = ttUSHORT(pairValue)
								straw = int(secondGlyph)
								if needle < straw {
									r = m - 1
								} else if needle > straw {
									l = m + 1
								} else {
									var xAdvance = ttSHORT(pairValue[2:])
									return int(xAdvance)
								}
							}
						}
						break

					case 2:
						{
							var valueFormat1 = ttUSHORT(table[4:])
							var valueFormat2 = ttUSHORT(table[6:])

							var classDef1Offset = ttUSHORT(table[8:])
							var classDef2Offset = ttUSHORT(table[10:])
							var glyph1class = stbtt__GetGlyphClass(table[classDef1Offset:], glyph1)
							var glyph2class = stbtt__GetGlyphClass(table[classDef2Offset:], glyph2)

							var class1Count = ttUSHORT(table[12:])
							var class2Count = ttUSHORT(table[14:])
							STBTT_assert(glyph1class < int(class1Count))
							STBTT_assert(glyph2class < int(class2Count))

							// TODO: Support more formats.
							STBTT_GPOS_TODO_assert(valueFormat1 == 4)
							if valueFormat1 != 4 {
								return 0
							}
							STBTT_GPOS_TODO_assert(valueFormat2 == 0)
							if valueFormat2 != 0 {
								return 0
							}

							if glyph1class >= 0 && glyph1class < int(class1Count) && glyph2class >= 0 && glyph2class < int(class2Count) {
								var class1Records = table[16:]
								var class2Records = class1Records[2*(glyph1class*int(class2Count)):]
								var xAdvance = ttSHORT(class2Records[2*glyph2class:])
								return int(xAdvance)
							}
						}
						break

					default:
						{
							// There are no other cases.
							STBTT_assert(false)
							break
						} // [DEAR IMGUI] removed ;
					}
				}
				break
			} // [DEAR IMGUI] removed ;

		default:
			// TODO: Implement other stuff.
			break
		}
	}

	return 0
}

type stbtt__edge struct {
	x0, y0, x1, y1 float
	invert         int
}

type stbtt__active_edge struct {
	next         *stbtt__active_edge
	fx, fdx, fdy float
	direction    float
	sy           float
	ey           float
}

// the edge passed in here does not cross the vertical line at x or the vertical line at x+1
// (i.e. it has already been clipped to those)
func stbtt__handle_clipped_edge(scanline []float, x int, e *stbtt__active_edge, x0, y0, x1, y1 float) {
	if y0 == y1 {
		return
	}
	STBTT_assert(y0 < y1)
	STBTT_assert(e.sy <= e.ey)
	if y0 > e.ey {
		return
	}
	if y1 < e.sy {
		return
	}
	if y0 < e.sy {
		x0 += (x1 - x0) * (e.sy - y0) / (y1 - y0)
		y0 = e.sy
	}
	if y1 > e.ey {
		x1 += (x1 - x0) * (e.ey - y1) / (y1 - y0)
		y1 = e.ey
	}

	if x0 == float(x) {
		STBTT_assert(x1 <= float(x)+1)
	} else if x0 == float(x)+1 {
		STBTT_assert(x1 >= float(x))
	} else if x0 <= float(x) {
		STBTT_assert(x1 <= float(x))
	} else if x0 >= float(x)+1 {
		STBTT_assert(x1 >= float(x)+1)
	} else {
		STBTT_assert(x1 >= float(x) && x1 <= float(x)+1)
	}

	if x0 <= float(x) && x1 <= float(x) {
		scanline[x] += e.direction * (y1 - y0)
	} else if x0 >= float(x)+1 && x1 >= float(x)+1 {

	} else {
		STBTT_assert(x0 >= float(x) && x0 <= float(x)+1 && x1 >= float(x) && x1 <= float(x)+1)
		scanline[x] += e.direction * (y1 - y0) * (1 - ((x0-float(x))+(x1-float(x)))/2) // coverage = 1 - average x position
	}
}

func stbtt__new_active(e *stbtt__edge, off_x int, start_point float, userdata any) *stbtt__active_edge {
	var z = new(stbtt__active_edge)
	var dxdy = (e.x1 - e.x0) / (e.y1 - e.y0)
	STBTT_assert(z != nil)
	//STBTT_assert(e.y0 <= start_point);
	if z == nil {
		return z
	}
	z.fdx = dxdy
	if dxdy != 0.0 {
		z.fdy = 1.0 / dxdy
	}
	z.fx = e.x0 + dxdy*(start_point-e.y0)
	z.fx -= float(off_x)
	if e.invert != 0 {
		z.direction = 1
	} else {
		z.direction = -1
	}
	z.sy = e.y0
	z.ey = e.y1
	z.next = nil
	return z
}

func stbtt__fill_active_edges_new(scanline []float, scanline_fill []float, scanline_fill_idx, len int, e *stbtt__active_edge, y_top float) {
	var y_bottom = y_top + 1

	for e != nil {
		// brute force every pixel

		// compute intersection points with top & bottom
		STBTT_assert(e.ey >= y_top)

		if e.fdx == 0 {
			var x0 = e.fx
			if x0 < float(len) {
				if x0 >= 0 {
					stbtt__handle_clipped_edge(scanline, (int)(x0), e, x0, y_top, x0, y_bottom)
					stbtt__handle_clipped_edge(scanline_fill[scanline_fill_idx-1:], (int)(x0)+1, e, x0, y_top, x0, y_bottom)
				} else {
					stbtt__handle_clipped_edge(scanline_fill[scanline_fill_idx-1:], 0, e, x0, y_top, x0, y_bottom)
				}
			}
		} else {
			var x0 = e.fx
			var dx = e.fdx
			var xb = x0 + dx
			var x_top, x_bottom float
			var sy0, sy1 float
			var dy = e.fdy
			STBTT_assert(e.sy <= y_bottom && e.ey >= y_top)

			// compute endpoints of line segment clipped to this scanline (if the
			// line segment starts on this scanline. x0 is the intersection of the
			// line with y_top, but that may be off the line segment.
			if e.sy > y_top {
				x_top = x0 + dx*(e.sy-y_top)
				sy0 = e.sy
			} else {
				x_top = x0
				sy0 = y_top
			}
			if e.ey < y_bottom {
				x_bottom = x0 + dx*(e.ey-y_top)
				sy1 = e.ey
			} else {
				x_bottom = xb
				sy1 = y_bottom
			}

			if x_top >= 0 && x_bottom >= 0 && x_top < float(len) && x_bottom < float(len) {
				// from here on, we don't have to range check x values

				if (int)(x_top) == (int)(x_bottom) {
					var height float
					// simple case, only spans one pixel
					var x = (int)(x_top)
					height = sy1 - sy0
					STBTT_assert(x >= 0 && x < len)
					scanline[x] += e.direction * (1 - ((x_top-float(x))+(x_bottom-float(x)))/2) * height
					scanline_fill[scanline_fill_idx+x] += e.direction * height // everything right of this pixel is filled
				} else {
					var x, x1, x2 int
					var y_crossing, step, sign, area float
					// covers 2+ pixels
					if x_top > x_bottom {
						// flip scanline vertically; signed area is the same
						var t float
						sy0 = y_bottom - (sy0 - y_top)
						sy1 = y_bottom - (sy1 - y_top)
						t = sy0
						sy0 = sy1
						sy1 = t
						t = x_bottom
						x_bottom = x_top
						x_top = t
						dx = -dx
						dy = -dy
						t = x0
						x0 = xb
						xb = t
					}

					x1 = (int)(x_top)
					x2 = (int)(x_bottom)
					// compute intersection with y axis at x1+1
					y_crossing = (float(x1)+1-float(x0))*dy + y_top

					sign = e.direction
					// area of the rectangle covered from y0..y_crossing
					area = sign * (y_crossing - sy0)
					// area of the triangle (x_top,y0), (x+1,y0), (x+1,y_crossing)
					scanline[x1] += area * (1 - ((float(x_top)-float(x1))+(float(x1)+1-float(x1)))/2)

					step = sign * dy
					for x = x1 + 1; x < x2; x++ {
						scanline[x] += area + step/2
						area += step
					}
					y_crossing += dy * float(x2-(x1+1))

					STBTT_assert(STBTT_fabs(area) <= 1.01)

					scanline[x2] += area + sign*(1-(float(x2-x2)+(x_bottom-float(x2)))/2)*(sy1-y_crossing)

					scanline_fill[scanline_fill_idx+x2] += sign * (sy1 - sy0)
				}
			} else {
				// if edge goes outside of box we're drawing, we require
				// clipping logic. since this does not match the intended use
				// of this library, we use a different, very slow brute
				// force implementation
				var x int
				for x = 0; x < len; x++ {
					// cases:
					//
					// there can be up to two intersections with the pixel. any intersection
					// with left or right edges can be handled by splitting into two (or three)
					// regions. intersections with top & bottom do not necessitate case-wise logic.
					//
					// the old way of doing this found the intersections with the left & right edges,
					// then used some simple logic to produce up to three segments in sorted order
					// from top-to-bottom. however, this had a problem: if an x edge was epsilon
					// across the x border, then the corresponding y position might not be distinct
					// from the other y segment, and it might ignored as an empty segment. to avoid
					// that, we need to explicitly produce segments based on x positions.

					// rename variables to clearly-defined pairs
					var y0 = y_top
					var x1 = (float)(x)
					var x2 = (float)(x + 1)
					var x3 = xb
					var y3 = y_bottom

					// x = e.x + e.dx * (y-y_top)
					// (y-y_top) = (x - e.x) / e.dx
					// y = (x - e.x) / e.dx + y_top
					var y1 = (float32(x)-x0)/dx + y_top
					var y2 = (float32(x)+1-x0)/dx + y_top

					if x0 < x1 && x3 > x2 { // three segments descending down-right
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x1, y1)
						stbtt__handle_clipped_edge(scanline, x, e, x1, y1, x2, y2)
						stbtt__handle_clipped_edge(scanline, x, e, x2, y2, x3, y3)
					} else if x3 < x1 && x0 > x2 { // three segments descending down-left
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x2, y2)
						stbtt__handle_clipped_edge(scanline, x, e, x2, y2, x1, y1)
						stbtt__handle_clipped_edge(scanline, x, e, x1, y1, x3, y3)
					} else if x0 < x1 && x3 > x1 { // two segments across x, down-right
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x1, y1)
						stbtt__handle_clipped_edge(scanline, x, e, x1, y1, x3, y3)
					} else if x3 < x1 && x0 > x1 { // two segments across x, down-left
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x1, y1)
						stbtt__handle_clipped_edge(scanline, x, e, x1, y1, x3, y3)
					} else if x0 < x2 && x3 > x2 { // two segments across x+1, down-right
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x2, y2)
						stbtt__handle_clipped_edge(scanline, x, e, x2, y2, x3, y3)
					} else if x3 < x2 && x0 > x2 { // two segments across x+1, down-left
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x2, y2)
						stbtt__handle_clipped_edge(scanline, x, e, x2, y2, x3, y3)
					} else { // one segment
						stbtt__handle_clipped_edge(scanline, x, e, x0, y0, x3, y3)
					}
				}
			}
		}
		e = e.next
	}
}

// directly AA rasterize edges w/o supersampling
func stbtt__rasterize_sorted_edges(result *stbtt__bitmap, e []stbtt__edge, n, vsubsample, off_x, off_y int, userdata any) {
	var active *stbtt__active_edge = nil
	var y, j, i int
	var scanline_data [129]float
	var scanline, scanline2 []float

	if result.w > 64 {
		scanline = make([]float, result.w*2+1)
	} else {
		scanline = scanline_data[:]
	}

	scanline2 = scanline[result.w:]

	y = off_y
	e[n].y0 = (float)(off_y+result.h) + 1

	for j < result.h {
		// find center of pixel for this scanline
		var scan_y_top = float(y) + 0.0
		var scan_y_bottom = y + 1.0
		var step = &active

		for i := 0; int(i) < result.w; i++ {
			scanline[i] = 0
		}
		for i := 0; int(i) < result.w+1; i++ {
			scanline2[i] = 0
		}

		// update all active edges;
		// remove all active edges that terminate before the top of this scanline
		for *step != nil {
			var z = *step
			if z.ey <= scan_y_top {
				*step = z.next // delete from list
				STBTT_assert(z.direction != 0)
				z.direction = 0
			} else {
				step = &((*step).next) // advance through list
			}
		}

		// insert all edges that start before the bottom of this scanline
		for e[0].y0 <= float(scan_y_bottom) {
			if e[0].y0 != e[0].y1 {
				var z = stbtt__new_active(&e[0], off_x, scan_y_top, userdata)
				if z != nil {
					if j == 0 && off_y != 0 {
						if z.ey < scan_y_top {
							// this can happen due to subpixel positioning and some kind of fp rounding error i think
							z.ey = scan_y_top
						}
					}
					STBTT_assert(z.ey >= scan_y_top) // if we get really unlucky a tiny bit of an edge can be out of bounds
					// insert at front
					z.next = active
					active = z
				}
			}
			e = e[1:]
		}

		// now process all active edges
		if active != nil {
			stbtt__fill_active_edges_new(scanline, scanline2, 1, result.w, active, scan_y_top)
		}

		{
			var sum float = 0
			for i = 0; i < result.w; i++ {
				var k float
				var m int
				sum += scanline2[i]
				k = scanline[i] + sum
				k = (float)(STBTT_fabs(k)*255 + 0.5)
				m = (int)(k)
				if m > 255 {
					m = 255
				}
				result.pixels[j*result.stride+i] = (byte)(m)
			}
		}
		// advance all the edges
		step = &active
		for *step != nil {
			var z = *step
			z.fx += z.fdx          // advance to position for current scanline
			step = &((*step).next) // advance through list
		}

		y++
		j++
	}
}

func STBTT__COMPARE(a, b *stbtt__edge) int {
	return bool2int(a.y0 < b.y0)
}

func stbtt__sort_edges_ins_sort(p []stbtt__edge, n int) {
	var i, j int
	for i = 1; i < n; i++ {
		var t = p[i]
		var a = &t
		j = i
		for j > 0 {
			var b = &p[j-1]
			var c = STBTT__COMPARE(a, b)
			if c == 0 {
				break
			}
			p[j] = p[j-1]
			j--
		}
		if i != j {
			p[j] = t
		}
	}
}

func stbtt__sort_edges_quicksort(p []stbtt__edge, n int) {
	/* threshold for transitioning to insertion sort */
	for n > 12 {
		var t stbtt__edge
		var c01, c12, c, m, i, j int

		/* compute median of three */
		m = n >> 1
		c01 = STBTT__COMPARE(&p[0], &p[m])
		c12 = STBTT__COMPARE(&p[m], &p[n-1])
		/* if 0 >= mid >= end, or 0 < mid < end, then use mid */
		if c01 != c12 {
			/* otherwise, we'll need to swap something else to middle */
			var z int
			c = STBTT__COMPARE(&p[0], &p[n-1])
			/* 0>mid && mid<n:  0>n => n; 0<n => 0 */
			/* 0<mid && mid>n:  0>n => 0; 0<n => n */
			if c == c12 {
				z = 0
			} else {
				z = n - 1
			}
			t = p[z]
			p[z] = p[m]
			p[m] = t
		}
		/* now p[m] is the median-of-three */
		/* swap it to the beginning so it won't move around */
		t = p[0]
		p[0] = p[m]
		p[m] = t

		/* partition loop */
		i = 1
		j = n - 1
		for {
			/* handling of equality is crucial here */
			/* for sentinels & efficiency with duplicates */
			for ; ; i++ {
				if STBTT__COMPARE(&p[i], &p[0]) == 0 {
					break
				}
			}
			for ; ; j-- {
				if STBTT__COMPARE(&p[0], &p[j]) == 0 {
					break
				}
			}
			/* make sure we haven't crossed */
			if i >= j {
				break
			}
			t = p[i]
			p[i] = p[j]
			p[j] = t

			i++
			j--
		}
		/* recurse on smaller side, iterate on larger */
		if j < (n - i) {
			stbtt__sort_edges_quicksort(p, j)
			p = p[i:]
			n = n - i
		} else {
			stbtt__sort_edges_quicksort(p[i:], n-i)
			n = j
		}
	}
}

func stbtt__sort_edges(p []stbtt__edge, n int) {
	stbtt__sort_edges_quicksort(p, n)
	stbtt__sort_edges_ins_sort(p, n)
}

type stbtt__point struct {
	x, y float
}

func stbtt__rasterize(result *stbtt__bitmap, pts []stbtt__point, wcount []int, windings int, scale_x, scale_y, shift_x, shift_y float, off_x, off_y, invert int, userdata any) {
	var y_scale_inv float
	if invert != 0 {
		y_scale_inv = -scale_y
	} else {
		y_scale_inv = scale_y
	}
	var e []stbtt__edge
	var n, i, j, k, m int

	var vsubsample float = 1
	// vsubsample should divide 255 evenly; otherwise we won't reach full opacity

	// now we have to blow out the windings into explicit edge lists
	n = 0
	for i = 0; i < windings; i++ {
		n += wcount[i]
	}

	e = make([]stbtt__edge, (n + 1)) // add an extra one as a sentinel
	if e == nil {
		return
	}
	n = 0

	m = 0
	for i = 0; i < windings; i++ {
		var p = pts[m:]
		m += wcount[i]
		j = wcount[i] - 1
		for k = 0; k < wcount[i]; j, k = k, k+1 {
			var a, b = k, j
			// skip the edge if horizontal
			if p[j].y == p[k].y {
				continue
			}
			// add edge from j to k to the list
			e[n].invert = 0

			var condition bool
			if invert != 0 {
				condition = p[j].y > p[k].y
			} else {
				condition = p[j].y < p[k].y
			}

			if condition {
				e[n].invert = 1
				a, b = j, k
			}
			e[n].x0 = p[a].x*scale_x + shift_x
			e[n].y0 = (p[a].y*y_scale_inv + shift_y) * vsubsample
			e[n].x1 = p[b].x*scale_x + shift_x
			e[n].y1 = (p[b].y*y_scale_inv + shift_y) * vsubsample
			n++
		}
	}

	// now sort the edges by their highest point (should snap to integer, and then by x)
	//STBTT_sort(e, n, sizeof(e[0]), stbtt__edge_compare);
	stbtt__sort_edges(e, n)

	// now, traverse the scanlines and find the intersections on each scanline, use xor winding rule
	stbtt__rasterize_sorted_edges(result, e, n, int(vsubsample), off_x, off_y, userdata)
}

func stbtt__add_point(points []stbtt__point, n int, x, y float) {
	if points == nil {
		return // during first pass, it's unallocated
	}
	points[n].x = x
	points[n].y = y
}

// tessellate until threshold p is happy... @TODO warped to compensate for non-linear stretching
func stbtt__tesselate_curve(points []stbtt__point, num_points *int, x0, y0, x1, y1, x2, y2, objspace_flatness_squared float, n int) int {
	// midpoint
	var mx = (x0 + 2*x1 + x2) / 4
	var my = (y0 + 2*y1 + y2) / 4
	// versus directly drawn line
	var dx = (x0+x2)/2 - mx
	var dy = (y0+y2)/2 - my
	if n > 16 { // 65536 segments on one curve better be enough!
		return 1
	}
	if dx*dx+dy*dy > objspace_flatness_squared { // half-pixel error allowed... need to be smaller if AA
		stbtt__tesselate_curve(points, num_points, x0, y0, (x0+x1)/2.0, (y0+y1)/2.0, mx, my, objspace_flatness_squared, n+1)
		stbtt__tesselate_curve(points, num_points, mx, my, (x1+x2)/2.0, (y1+y2)/2.0, x2, y2, objspace_flatness_squared, n+1)
	} else {
		stbtt__add_point(points, *num_points, x2, y2)
		*num_points = *num_points + 1
	}
	return 1
}

func stbtt__tesselate_cubic(points []stbtt__point, num_points *int, x0, y0, x1, y1, x2, y2, x3, y3, objspace_flatness_squared float, n int) {
	// @TODO this "flatness" calculation is just made-up nonsense that seems to work well enough
	var dx0 = x1 - x0
	var dy0 = y1 - y0
	var dx1 = x2 - x1
	var dy1 = y2 - y1
	var dx2 = x3 - x2
	var dy2 = y3 - y2
	var dx = x3 - x0
	var dy = y3 - y0
	var longlen = (float)(STBTT_sqrt(dx0*dx0+dy0*dy0) + STBTT_sqrt(dx1*dx1+dy1*dy1) + STBTT_sqrt(dx2*dx2+dy2*dy2))
	var shortlen = (float)(STBTT_sqrt(dx*dx + dy*dy))
	var flatness_squared = longlen*longlen - shortlen*shortlen

	if n > 16 { // 65536 segments on one curve better be enough!
		return
	}

	if flatness_squared > objspace_flatness_squared {
		var x01 = (x0 + x1) / 2
		var y01 = (y0 + y1) / 2
		var x12 = (x1 + x2) / 2
		var y12 = (y1 + y2) / 2
		var x23 = (x2 + x3) / 2
		var y23 = (y2 + y3) / 2

		var xa = (x01 + x12) / 2
		var ya = (y01 + y12) / 2
		var xb = (x12 + x23) / 2
		var yb = (y12 + y23) / 2

		var mx = (xa + xb) / 2
		var my = (ya + yb) / 2

		stbtt__tesselate_cubic(points, num_points, x0, y0, x01, y01, xa, ya, mx, my, objspace_flatness_squared, n+1)
		stbtt__tesselate_cubic(points, num_points, mx, my, xb, yb, x23, y23, x3, y3, objspace_flatness_squared, n+1)
	} else {
		stbtt__add_point(points, *num_points, x3, y3)
		*num_points = *num_points + 1
	}
}

// returns number of contours
func stbtt_FlattenCurves(vertices []stbtt_vertex, num_verts int, objspace_flatness float, contour_lengths *[]int, num_contours *int, userdata any) []stbtt__point {
	var points []stbtt__point
	var num_points int = 0

	var objspace_flatness_squared = objspace_flatness * objspace_flatness
	var i, n, start, pass int

	// count how many "moves" there are to get the contour count
	for i = 0; i < num_verts; i++ {
		if vertices[i].vtype == STBTT_vmove {
			n++
		}
	}

	*num_contours = n
	if n == 0 {
		return nil
	}

	*contour_lengths = make([]int, n)

	if *contour_lengths == nil {
		*num_contours = 0
		return nil
	}

	// make two passes through the points so we don't need to realloc
	for pass = 0; pass < 2; pass++ {
		var x, y float
		if pass == 1 {
			points = make([]stbtt__point, num_points)
			if points == nil {
				goto error
			}
		}
		num_points = 0
		n = -1
		for i = 0; i < num_verts; i++ {
			switch vertices[i].vtype {
			case STBTT_vmove:
				// start the next contour
				if n >= 0 {
					(*contour_lengths)[n] = num_points - start
				}
				n++
				start = num_points

				x = float(vertices[i].x)
				y = float(vertices[i].y)
				stbtt__add_point(points, num_points, x, y)
				break
			case STBTT_vline:
				x = float(vertices[i].x)
				y = float(vertices[i].y)
				stbtt__add_point(points, num_points, x, y)
				num_points++
				break
			case STBTT_vcurve:
				stbtt__tesselate_curve(points, &num_points, x, y,
					float(vertices[i].cx), float(vertices[i].cy),
					float(vertices[i].x), float(vertices[i].y),
					objspace_flatness_squared, 0)
				x = float(vertices[i].x)
				y = float(vertices[i].y)
				break
			case STBTT_vcubic:
				stbtt__tesselate_cubic(points, &num_points, x, y,
					float(vertices[i].cx), float(vertices[i].cy),
					float(vertices[i].cx1), float(vertices[i].cy1),
					float(vertices[i].x), float(vertices[i].y),
					objspace_flatness_squared, 0)
				x = float(vertices[i].x)
				y = float(vertices[i].y)
				break
			default:
				panic("unknown vtype")
			}
		}
		(*contour_lengths)[n] = num_points - start
	}

	return points
error:
	*contour_lengths = nil
	*num_contours = 0
	return nil
}

func stbtt_BakeFontBitmap_internal(data []byte, offset int, // font location (use offset=0 for plain .ttf)
	pixel_height float, // height of font in pixels
	pixels []byte, pw, ph, // bitmap to be filled in
	first_char, num_chars int, // characters to bake
	chardata []stbtt_bakedchar) int {
	var scale float
	var x, y, bottom_y, i int
	var f FontInfo
	f.userdata = nil
	if InitFont(&f, data, offset) == 0 {
		return -1
	}
	x = 1
	y = 1
	bottom_y = 1

	scale = ScaleForPixelHeight(&f, pixel_height)

	for i = 0; i < num_chars; i++ {
		var advance, lsb, x0, y0, x1, y1, gw, gh int
		var g = FindGlyphIndex(&f, first_char+i)
		stbtt_GetGlyphHMetrics(&f, g, &advance, &lsb)
		stbtt_GetGlyphBitmapBox(&f, g, scale, scale, &x0, &y0, &x1, &y1)
		gw = x1 - x0
		gh = y1 - y0
		if x+gw+1 >= pw {
			y = bottom_y
			x = 1 // advance to next row
		}
		if y+gh+1 >= ph { // check if it fits vertically AFTER potentially moving to next row
			return -i
		}
		STBTT_assert(x+gw < pw)
		STBTT_assert(y+gh < ph)
		stbtt_MakeGlyphBitmap(&f, pixels[x+y*pw:], gw, gh, pw, scale, scale, g)
		chardata[i].x0 = uint16((stbtt_int16)(x))
		chardata[i].y0 = uint16((stbtt_int16)(y))
		chardata[i].x1 = uint16((stbtt_int16)(x + gw))
		chardata[i].y1 = uint16((stbtt_int16)(y + gh))
		chardata[i].xadvance = scale * float(advance)
		chardata[i].xoff = (float)(x0)
		chardata[i].yoff = (float)(y0)
		x = x + gw + 1
		if y+gh+1 > bottom_y {
			bottom_y = y + gh + 1
		}
	}
	return bottom_y
}

const STBTT__OVER_MASK = STBTT_MAX_OVERSAMPLE - 1

func stbtt__h_prefilter(pixels []byte, w, h, stride_in_bytes int, kernel_width uint) {
	var buffer [STBTT_MAX_OVERSAMPLE]byte
	var safe_w = w - int(kernel_width)
	var j int
	for j = 0; j < h; j++ {
		var i int
		var total uint

		total = 0

		// make kernel_width a constant in common cases so compiler can optimize out the divide
		switch kernel_width {
		case 2:
			for i = 0; i <= safe_w; i++ {
				total += uint(pixels[i]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i]
				pixels[i] = (byte)(total / 2)
			}
			break
		case 3:
			for i = 0; i <= safe_w; i++ {
				total += uint(pixels[i]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i]
				pixels[i] = (byte)(total / 3)
			}
			break
		case 4:
			for i = 0; i <= safe_w; i++ {
				total += uint(pixels[i]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i]
				pixels[i] = (byte)(total / 4)
			}
			break
		case 5:
			for i = 0; i <= safe_w; i++ {
				total += uint(pixels[i]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i]
				pixels[i] = (byte)(total / 5)
			}
			break
		default:
			for i = 0; i <= safe_w; i++ {
				total += uint(pixels[i]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i]
				pixels[i] = (byte)(total / kernel_width)
			}
			break
		}

		for ; i < w; i++ {
			STBTT_assert(pixels[i] == 0)
			total -= uint(buffer[i&STBTT__OVER_MASK])
			pixels[i] = (byte)(total / kernel_width)
		}

		pixels = pixels[stride_in_bytes:]
	}
}

func stbtt__v_prefilter(pixels []byte, w, h, stride_in_bytes int, kernel_width uint) {
	var buffer [STBTT_MAX_OVERSAMPLE]byte
	var safe_h = h - int(kernel_width)
	var j int

	for j = 0; j < w; j++ {
		var i int
		var total uint

		total = 0

		// make kernel_width a constant in common cases so compiler can optimize out the divide
		switch kernel_width {
		case 2:
			for i = 0; i <= safe_h; i++ {
				total += uint(pixels[i*stride_in_bytes]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i*stride_in_bytes]
				pixels[i*stride_in_bytes] = (byte)(total / 2)
			}
			break
		case 3:
			for i = 0; i <= safe_h; i++ {
				total += uint(pixels[i*stride_in_bytes]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i*stride_in_bytes]
				pixels[i*stride_in_bytes] = (byte)(total / 3)
			}
			break
		case 4:
			for i = 0; i <= safe_h; i++ {
				total += uint(pixels[i*stride_in_bytes]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i*stride_in_bytes]
				pixels[i*stride_in_bytes] = (byte)(total / 4)
			}
			break
		case 5:
			for i = 0; i <= safe_h; i++ {
				total += uint(pixels[i*stride_in_bytes]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i*stride_in_bytes]
				pixels[i*stride_in_bytes] = (byte)(total / 5)
			}
			break
		default:
			for i = 0; i <= safe_h; i++ {
				total += uint(pixels[i*stride_in_bytes]) - uint(buffer[i&STBTT__OVER_MASK])
				buffer[(uint(i)+kernel_width)&STBTT__OVER_MASK] = pixels[i*stride_in_bytes]
				pixels[i*stride_in_bytes] = (byte)(total / kernel_width)
			}
			break
		}

		for ; i < h; i++ {
			STBTT_assert(pixels[i*stride_in_bytes] == 0)
			total -= uint(buffer[i&STBTT__OVER_MASK])
			pixels[i*stride_in_bytes] = (byte)(total / kernel_width)
		}

		pixels = pixels[1:]
	}
}

func stbtt__oversample_shift(oversample int) float {
	if oversample == 0 {
		return 0.0
	}

	// The prefilter is a box filter of width "oversample",
	// which shifts phase by (oversample - 1)/2 pixels in
	// oversampled space. We want to shift in the opposite
	// direction to counter this.
	return (float)(-(oversample - 1)) / (2.0 * (float)(oversample))
}

func STBTT_min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func STBTT_max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func STBTT_minf(a, b float) float {
	if a < b {
		return a
	}
	return b
}

func STBTT_maxf(a, b float) float {
	if a > b {
		return a
	}
	return b
}

func stbtt__ray_intersect_bezier(orig, ray, q0, q1, q2 [2]float, hits [2][2]float) int {
	var q0perp = q0[1]*ray[0] - q0[0]*ray[1]
	var q1perp = q1[1]*ray[0] - q1[0]*ray[1]
	var q2perp = q2[1]*ray[0] - q2[0]*ray[1]
	var roperp = orig[1]*ray[0] - orig[0]*ray[1]

	var a = q0perp - 2*q1perp + q2perp
	var b = q1perp - q0perp
	var c = q0perp - roperp

	var s0 float = 0.0
	var s1 float = 0.0
	var num_s int = 0

	if a != 0.0 {
		var discr = b*b - a*c
		if discr > 0.0 {
			var rcpna = -1 / a
			var d = (float)(STBTT_sqrt(discr))
			s0 = (b + d) * rcpna
			s1 = (b - d) * rcpna
			if s0 >= 0.0 && s0 <= 1.0 {
				num_s = 1
			}
			if d > 0.0 && s1 >= 0.0 && s1 <= 1.0 {
				if num_s == 0 {
					s0 = s1
				}
				num_s++
			}
		}
	} else {
		// 2*b*s + c = 0
		// s = -c / (2*b)
		s0 = c / (-2 * b)
		if s0 >= 0.0 && s0 <= 1.0 {
			num_s = 1
		}
	}

	if num_s == 0 {
		return 0
	} else {
		var rcp_len2 = 1 / (ray[0]*ray[0] + ray[1]*ray[1])
		var rayn_x = ray[0] * rcp_len2
		var rayn_y = ray[1] * rcp_len2

		var q0d = q0[0]*rayn_x + q0[1]*rayn_y
		var q1d = q1[0]*rayn_x + q1[1]*rayn_y
		var q2d = q2[0]*rayn_x + q2[1]*rayn_y
		var rod = orig[0]*rayn_x + orig[1]*rayn_y

		var q10d = q1d - q0d
		var q20d = q2d - q0d
		var q0rd = q0d - rod

		hits[0][0] = q0rd + s0*(2.0-2.0*s0)*q10d + s0*s0*q20d
		hits[0][1] = a*s0 + b

		if num_s > 1 {
			hits[1][0] = q0rd + s1*(2.0-2.0*s1)*q10d + s1*s1*q20d
			hits[1][1] = a*s1 + b
			return 2
		} else {
			return 1
		}
	}
}

func equal(a, b []float) int {
	return bool2int(a[0] == b[0] && a[1] == b[1])
}

func stbtt__compute_crossings_x(x, y float, nverts int, verts []stbtt_vertex) int {
	var i int
	var orig, ray = [2]float{}, [2]float{1, 0}

	var y_frac float
	var winding int = 0

	orig[0] = x
	//orig[1] = y; // [DEAR IMGUI] commented double assignment

	// make sure y never passes through a vertex of the shape
	y_frac = (float)(STBTT_fmod(y, 1.0))
	if y_frac < 0.01 {
		y += 0.01
	} else if y_frac > 0.99 {
		y -= 0.01
	}
	orig[1] = y

	// test a ray from (-infinity,y) to (x,y)
	for i = 0; i < nverts; i++ {
		if verts[i].vtype == STBTT_vline {
			var x0 = (int)(verts[i-1].x)
			var y0 = (int)(verts[i-1].y)
			var x1 = (int)(verts[i].x)
			var y1 = (int)(verts[i].y)
			if y > float(STBTT_min(y0, y1)) && y < float(STBTT_max(y0, y1)) && x > float(STBTT_min(x0, x1)) {
				var x_inter = (y-float(y0))/(float(y1)-float(y0))*(float(x1)-float(x0)) + float(x0)
				if x_inter < x {
					if y0 < y1 {
						winding += 1
					} else {
						winding -= 1
					}
				}
			}
		}
		if verts[i].vtype == STBTT_vcurve {
			var x0 = (int)(verts[i-1].x)
			var y0 = (int)(verts[i-1].y)
			var x1 = (int)(verts[i].cx)
			var y1 = (int)(verts[i].cy)
			var x2 = (int)(verts[i].x)
			var y2 = (int)(verts[i].y)
			var ax = STBTT_min(x0, STBTT_min(x1, x2))
			var ay = STBTT_min(y0, STBTT_min(y1, y2))
			var by = STBTT_max(y0, STBTT_max(y1, y2))
			if y > float(ay) && y < float(by) && x > float(ax) {
				var q0, q1, q2 [2]float
				var hits [2][2]float
				q0[0] = (float)(x0)
				q0[1] = (float)(y0)
				q1[0] = (float)(x1)
				q1[1] = (float)(y1)
				q2[0] = (float)(x2)
				q2[1] = (float)(y2)
				if equal(q0[:], q1[:]) != 0 || equal(q1[:], q2[:]) != 0 {
					x0 = (int)(verts[i-1].x)
					y0 = (int)(verts[i-1].y)
					x1 = (int)(verts[i].x)
					y1 = (int)(verts[i].y)
					if y > float(STBTT_min(y0, y1)) && y < float(STBTT_max(y0, y1)) && x > float(STBTT_min(x0, x1)) {
						var x_inter = (y-float(y0))/(float(y1)-float(y0))*(float(x1)-float(x0)) + float(x0)
						if x_inter < x {
							if y0 < y1 {
								winding += 1
							} else {
								winding -= 1
							}
						}
					}
				} else {
					var num_hits = stbtt__ray_intersect_bezier(orig, ray, q0, q1, q2, hits)
					if num_hits >= 1 {
						if hits[0][0] < 0 {
							if hits[0][1] < 0 {
								winding += 1
							} else {
								winding -= 1
							}
						}
					}
					if num_hits >= 2 {
						if hits[1][0] < 0 {
							if hits[1][1] < 0 {
								winding -= 1
							} else {
								winding += 1
							}
						}
					}
				}
			}
		}
	}
	return winding
}

func stbtt__cuberoot(x float) float {
	if x < 0 {
		return -(float)(STBTT_pow(-x, 1.0/3.0))
	} else {
		return (float)(STBTT_pow(x, 1.0/3.0))
	}
}

// x^3 + c*x^2 + b*x + a = 0
func stbtt__solve_cubic(a, b, c float, r []float) int {
	var s = -a / 3
	var p = b - a*a/3
	var q = a*(2*a*a-9*b)/27 + c
	var p3 = p * p * p
	var d = q*q + 4*p3/27
	if d >= 0 {
		var z = (float)(STBTT_sqrt(d))
		var u = (-q + z) / 2
		var v = (-q - z) / 2
		u = stbtt__cuberoot(u)
		v = stbtt__cuberoot(v)
		r[0] = s + u + v
		return 1
	} else {
		var u = (float)(STBTT_sqrt(-p / 3))
		var v = (float)(STBTT_acos(-STBTT_sqrt(-27/p3)*q/2) / 3) // p3 must be negative, since d is negative
		var m = (float)(STBTT_cos(v))
		var n = (float)(STBTT_cos(v-3.141592/2) * 1.732050808)
		r[0] = s + u*2*m
		r[1] = s - u*(m+n)
		r[2] = s - u*(m-n)

		//STBTT_assert( STBTT_fabs(((r[0]+a)*r[0]+b)*r[0]+c) < 0.05f);  // these asserts may not be safe at all scales, though they're in bezier t parameter units so maybe?
		//STBTT_assert( STBTT_fabs(((r[1]+a)*r[1]+b)*r[1]+c) < 0.05f);
		//STBTT_assert( STBTT_fabs(((r[2]+a)*r[2]+b)*r[2]+c) < 0.05f);
		return 3
	}
}

// check if a utf8 string contains a prefix which is the utf16 string; if so return length of matching utf8 string
func stbtt__CompareUTF8toUTF16_bigendian_prefix(s1 []stbtt_uint8, len1 stbtt_int32, s2 []stbtt_uint8, len2 stbtt_int32) stbtt_int32 {
	var i stbtt_int32 = 0

	// convert utf16 to utf8 and compare the results while converting
	for len2 != 0 {
		var ch = stbtt_uint16(s2[0])*256 + stbtt_uint16(s2[1])
		if ch < 0x80 {
			if i >= len1 {
				return -1
			}

			if stbtt_uint16(s1[i]) != ch {
				return -1
			}
			i++
		} else if ch < 0x800 {
			if i+1 >= len1 {
				return -1
			}

			if stbtt_uint16(s1[i]) != 0xc0+(ch>>6) {
				return -1
			}
			i++
			if stbtt_uint16(s1[i]) != 0x80+(ch&0x3f) {
				return -1
			}
			i++
		} else if ch >= 0xd800 && ch < 0xdc00 {
			var c stbtt_uint32
			var ch2 = stbtt_uint16(s2[2])*256 + stbtt_uint16(s2[3])
			if i+3 >= len1 {
				return -1
			}
			c = (stbtt_uint32(ch-0xd800) << 10) + stbtt_uint32(ch2-0xdc00) + stbtt_uint32(0x10000)

			if uint(s1[i]) != 0xf0+(c>>18) {
				return -1
			}
			i++
			if uint(s1[i]) != 0x80+((c>>12)&0x3f) {
				return -1
			}
			i++
			if uint(s1[i]) != 0x80+((c>>6)&0x3f) {
				return -1
			}
			i++
			if uint(s1[i]) != 0x80+((c)&0x3f) {
				return -1
			}
			i++
			s2 = s2[2:] // plus another 2 below
			len2 -= 2
		} else if ch >= 0xdc00 && ch < 0xe000 {
			return -1
		} else {
			if i+2 >= len1 {
				return -1
			}

			if stbtt_uint16(s1[i]) != 0xe0+(ch>>12) {
				return -1
			}
			i++
			if stbtt_uint16(s1[i]) != 0x80+((ch>>6)&0x3f) {
				return -1
			}
			i++
			if stbtt_uint16(s1[i]) != 0x80+((ch)&0x3f) {
				return -1
			}
			i++
		}
		s2 = s2[2:]
		len2 -= 2
	}
	return i
}

func stbtt_CompareUTF8toUTF16_bigendian_internal(s1 []byte, len1 int, s2 []byte, len2 int) int {
	return bool2int(len1 == stbtt__CompareUTF8toUTF16_bigendian_prefix(s1, len1, s2, len2))
}

func stbtt__matchpair(fc []byte, nm stbtt_uint32, name []byte, nlen stbtt_int32, target_id stbtt_int32, next_id stbtt_int32) int {
	var i stbtt_int32
	var count = int(ttUSHORT(fc[nm+2:]))
	var stringOffset = int(nm + uint(ttUSHORT(fc[nm+4:])))

	for i = 0; i < count; i++ {
		var loc = nm + 6 + uint(12*i)
		var id = int(ttUSHORT(fc[loc+6:]))
		if id == target_id {
			// find the encoding
			var platform = int(ttUSHORT(fc[loc+0:]))
			var encoding = int(ttUSHORT(fc[loc+2:]))
			var language = int(ttUSHORT(fc[loc+4:]))

			// is this a Unicode encoding?
			if platform == 0 || (platform == 3 && encoding == 1) || (platform == 3 && encoding == 10) {
				var slen = int(ttUSHORT(fc[loc+8:]))
				var off = int(ttUSHORT(fc[loc+10:]))

				// check if there's a prefix match
				var matchlen = stbtt__CompareUTF8toUTF16_bigendian_prefix(name, nlen, fc[stringOffset+off:], slen)
				if matchlen >= 0 {
					// check for target_id+1 immediately following, with same encoding & language
					if i+1 < count && int(ttUSHORT(fc[loc+12+6:])) == next_id &&
						int(ttUSHORT(fc[loc+12:])) == platform && int(ttUSHORT(fc[loc+12+2:])) == encoding &&
						int(ttUSHORT(fc[loc+12+4:])) == language {

						slen = int(ttUSHORT(fc[loc+12+8:]))
						off = int(ttUSHORT(fc[loc+12+10:]))
						if slen == 0 {
							if matchlen == nlen {
								return 1
							}
						} else if matchlen < nlen && name[matchlen] == ' ' {
							matchlen++
							if stbtt_CompareUTF8toUTF16_bigendian_internal(name[matchlen:], nlen-matchlen, fc[stringOffset+off:], slen) != 0 {
								return 1
							}
						}
					} else {
						// if nothing immediately following
						if matchlen == nlen {
							return 1
						}
					}
				}
			}

			// @TODO handle other encodings
		}
	}
	return 0
}

func stbtt__matches(fc []byte, offset stbtt_uint32, name []byte, flags stbtt_int32) int {
	var nlen = int(len(name))
	var nm, hd stbtt_uint32
	if stbtt__isfont(fc[offset:]) == 0 {
		return 0
	}

	// check italics/bold/underline flags in macStyle...
	if flags != 0 {
		hd = stbtt__find_table(fc, offset, "head")
		if bool2int(ttUSHORT(fc[hd+44:])&7 != 0) != bool2int(flags&7 != 0) {
			return 0
		}
	}

	nm = stbtt__find_table(fc, offset, "name")
	if nm == 0 {
		return 0
	}

	if flags != 0 {
		// if we checked the macStyle flags, then just check the family and ignore the subfamily
		if stbtt__matchpair(fc, nm, name, nlen, 16, -1) != 0 {
			return 1
		}
		if stbtt__matchpair(fc, nm, name, nlen, 1, -1) != 0 {
			return 1
		}
		if stbtt__matchpair(fc, nm, name, nlen, 3, -1) != 0 {
			return 1
		}
	} else {
		if stbtt__matchpair(fc, nm, name, nlen, 16, 17) != 0 {
			return 1
		}
		if stbtt__matchpair(fc, nm, name, nlen, 1, 2) != 0 {
			return 1
		}
		if stbtt__matchpair(fc, nm, name, nlen, 3, -1) != 0 {
			return 1
		}
	}

	return 0
}

func stbtt_FindMatchingFont_internal(font_collection []byte, name_utf8 []byte, flags stbtt_int32) int {
	var i stbtt_int32
	for i = 0; ; i++ {
		var off = GetFontOffsetForIndex(font_collection, i)
		if off < 0 {
			return off
		}
		if stbtt__matches(font_collection, uint(off), name_utf8, flags) != 0 {
			return off
		}
	}
}

// FULL VERSION HISTORY
//
//   1.19 (2021-10-14) Ported to Go
//   1.19 (2018-02-11) OpenType GPOS kerning (horizontal only), STBTT_fmod
//   1.18 (2018-01-29) add missing function
//   1.17 (2017-07-23) make more arguments const; doc fix
//   1.16 (2017-07-12) SDF support
//   1.15 (2017-03-03) make more arguments const
//   1.14 (2017-01-16) num-fonts-in-TTC function
//   1.13 (2017-01-02) support OpenType fonts, certain Apple fonts
//   1.12 (2016-10-25) suppress warnings about casting away const with -Wcast-qual
//   1.11 (2016-04-02) fix unused-variable warning
//   1.10 (2016-04-02) allow user-defined fabs() replacement
//                     fix memory leak if fontsize=0.0
//                     fix warning from duplicate typedef
//   1.09 (2016-01-16) warning fix; avoid crash on outofmem; use alloc userdata for PackFontRanges
//   1.08 (2015-09-13) document stbtt_Rasterize(); fixes for vertical & horizontal edges
//   1.07 (2015-08-01) allow PackFontRanges to accept arrays of sparse codepoints;
//                     allow PackFontRanges to pack and render in separate phases;
//                     fix stbtt_GetFontOFfsetForIndex (never worked for non-0 input?);
//                     fixed an assert() bug in the new rasterizer
//                     replace assert() with STBTT_assert() in new rasterizer
//   1.06 (2015-07-14) performance improvements (~35% faster on x86 and x64 on test machine)
//                     also more precise AA rasterizer, except if shapes overlap
//                     remove need for STBTT_sort
//   1.05 (2015-04-15) fix misplaced definitions for STBTT_STATIC
//   1.04 (2015-04-15) typo in example
//   1.03 (2015-04-12) STBTT_STATIC, fix memory leak in new packing, various fixes
//   1.02 (2014-12-10) fix various warnings & compile issues w/ stb_rect_pack, C++
//   1.01 (2014-12-08) fix subpixel position when oversampling to exactly match
//                        non-oversampled; STBTT_POINT_SIZE for packed case only
//   1.00 (2014-12-06) add new PackBegin etc. API, w/ support for oversampling
//   0.99 (2014-09-18) fix multiple bugs with subpixel rendering (ryg)
//   0.9  (2014-08-07) support certain mac/iOS fonts without an MS platformID
//   0.8b (2014-07-07) fix a warning
//   0.8  (2014-05-25) fix a few more warnings
//   0.7  (2013-09-25) bugfix: subpixel glyph bug fixed in 0.5 had come back
//   0.6c (2012-07-24) improve documentation
//   0.6b (2012-07-20) fix a few more warnings
//   0.6  (2012-07-17) fix warnings; added stbtt_ScaleForMappingEmToPixels,
//                        stbtt_GetFontBoundingBox, stbtt_IsGlyphEmpty
//   0.5  (2011-12-09) bugfixes:
//                        subpixel glyph renderer computed wrong bounding box
//                        first vertex of shape can be off-curve (FreeSans)
//   0.4b (2011-12-03) fixed an error in the font baking example
//   0.4  (2011-12-01) kerning, subpixel rendering (tor)
//                    bugfixes for:
//                        codepoint-to-glyph conversion using table fmt=12
//                        codepoint-to-glyph conversion using table fmt=4
//                        stbtt_GetBakedQuad with non-square texture (Zer)
//                    updated Hello World! sample to use kerning and subpixel
//                    fixed some warnings
//   0.3  (2009-06-24) cmap fmt=12, compound shapes (MM)
//                    userdata, malloc-from-userdata, non-zero fill (stb)
//   0.2  (2009-03-11) Fix unsigned/signed char warnings
//   0.1  (2009-03-09) First public release
//

/*
------------------------------------------------------------------------------
This software is available under 2 licenses -- choose whichever you prefer.
------------------------------------------------------------------------------
ALTERNATIVE A - MIT License
Copyright (c) 2017 Sean Barrett
Copyright (c) 2021 Quentin Quaadgras
Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
------------------------------------------------------------------------------
ALTERNATIVE B - Public Domain (www.unlicense.org)
This is free and unencumbered software released into the public domain.
Anyone is free to copy, modify, publish, use, compile, sell, or distribute this
software, either in source code form or as a compiled binary, for any purpose,
commercial or non-commercial, and by any means.
In jurisdictions that recognize copyright laws, the author or authors of this
software dedicate any and all copyright interest in the software to the public
domain. We make this dedication for the benefit of the public at large and to
the detriment of our heirs and successors. We intend this dedication to be an
overt act of relinquishment in perpetuity of all present and future rights to
this software under copyright law.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
------------------------------------------------------------------------------
*/
