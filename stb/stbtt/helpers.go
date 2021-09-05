package stbtt

import (
	"log"
	"math"
)

type ttBYTE = uint8
type ttCHAR = int8

func ttUSHORT(p []byte) uint16 {
	return uint16(p[0])*256 + uint16(p[1])
}

func ttSHORT(p []byte) int16 {
	return int16(p[0])*256 + int16(p[1])
}

func ttULONG(p []byte) uint32 {
	return (uint32(p[0]) << 24) + (uint32(p[1]) << 16) + (uint32(p[2]) << 8) + uint32(p[3])
}

func ttLONG(p []byte) int32 {
	return int32(ttULONG(p))
}

func ttFixed(p []byte) uint32 {
	return (uint32(p[0]) << 24) + (uint32(p[1]) << 16) + (uint32(p[2]) << 8) + uint32(p[3])
}

func tag4(p []byte, c0, c1, c2, c3 byte) bool {
	return p[0] == c0 && p[1] == c1 && p[2] == c2 && p[3] == c3
}

func tag(p []byte, str string) bool {
	return tag4(p, str[0], str[1], str[2], str[3])
}

func isFont(font []byte) bool {
	switch {
	case tag4(font, '1', 0, 0, 0): // TrueType 1
		return true
	case tag(font, "typ1"): // TrueType with type 1 font -- we don't support this!
		return true
	case tag(font, "OTTO"): // OpenType with CFF
		return true
	case tag4(font, 0, 1, 0, 0): // OpenType 1.0
		return true
	case tag(font, "true"): // Apple specification for TrueType fonts
		return true
	default:
		return false
	}
}

func findTable(data []byte, fontstart uint32, nametag string) uint32 {
	var num_tables = ttUSHORT(data[fontstart+4:])
	var tabledir = fontstart + 12
	for i := uint32(0); i < uint32(num_tables); i++ {
		var loc = tabledir + 16*i
		if tag(data[loc+0:], nametag) {
			return ttULONG(data[loc+8:])
		}
	}
	return 0
}

func getFontOffsetForIndex(font_collection []byte, index int32) int32 {
	// if it's just a font, there's only one valid index
	if isFont(font_collection) {
		if index == 0 {
			return 0
		}
		return -1
	}

	// check if it's a TTC
	if tag(font_collection, "ttcf") {
		// version 1?
		if ttULONG(font_collection[4:]) == 0x00010000 || ttULONG(font_collection[4:]) == 0x00020000 {
			var n = ttLONG(font_collection[8:])
			if index >= n {
				return -1
			}
			return int32(ttULONG(font_collection[12+index*4:]))
		}
	}
	return -1
}

func getNumberOfFonts(font_collection []byte) int32 {
	// if it's just a font, there's only one valid font
	if isFont(font_collection) {
		return 1
	}

	// check if it's a TTC
	if tag(font_collection, "ttcf") {
		// version 1?
		if ttULONG(font_collection[4:]) == 0x00010000 || ttULONG(font_collection[4:]) == 0x00020000 {
			return ttLONG(font_collection[8:])
		}
	}

	return 0
}

func getSubrs(cff buf, fontdict buf) buf {
	var subrsoff [1]uint32
	var private_loc [2]uint32
	dictGetInts(&fontdict, 18, private_loc[:])
	if private_loc[1] == 0 || private_loc[0] == 0 {
		return buf{}
	}
	pdict := bufRange(&cff, int32(private_loc[1]), int32(private_loc[0]))
	dictGetInts(&pdict, 19, subrsoff[:])
	if subrsoff[0] == 0 {
		return buf{}
	}
	bufSeek(&cff, int32(private_loc[1]+subrsoff[0]))
	return cffGetIndex(&cff)
}

func getSVG(info *FontInfo) int32 {
	var t uint32
	if info.svg < 0 {
		t = findTable(info.data, uint32(info.fontstart), "SVG ")
		if t == 0 {
			info.svg = 0
		} else {
			info.svg = int32(t + ttULONG(info.data[t+2:]))
		}
	}
	return info.svg
}

func initFont(info *FontInfo, data []byte, fontstart int32) int32 {
	var cmap, t uint32
	var i, numTables uint32

	info.data = data
	info.fontstart = fontstart
	info.cff = buf{}

	ufontstart := uint32(fontstart)

	cmap = findTable(data, ufontstart, "cmap")             // required
	info.loca = int32(findTable(data, ufontstart, "loca")) // required
	info.head = int32(findTable(data, ufontstart, "head")) // required
	info.glyf = int32(findTable(data, ufontstart, "glyf")) // required
	info.hhea = int32(findTable(data, ufontstart, "hhea")) // required
	info.hmtx = int32(findTable(data, ufontstart, "hmtx")) // required
	info.kern = int32(findTable(data, ufontstart, "kern")) // not required
	info.gpos = int32(findTable(data, ufontstart, "GPOS")) // not required

	if cmap == 0 || info.head == 0 || info.hhea == 0 || info.hmtx == 0 {
		return 0
	}
	if info.glyf != 0 {
		if info.loca == 0 {
			return 0
		}
	} else {
		// initialization for CFF / Type2 fonts (OTF)
		var b, topdict, topdictidx buf
		var cstype = [1]uint32{2}
		var charstrings, fdarrayoff, fdselectoff [1]uint32

		cff := findTable(data, ufontstart, "CFF ")
		if cff == 0 {
			return 0
		}

		info.fontdicts = buf{}
		info.fdselect = buf{}

		// @TODO this should use size from table (not 512MB)
		info.cff = buf{data: data[cff:]}
		b = info.cff

		// read the header
		bufSkip(&b, 2)
		bufSeek(&b, int32(bufGet8(&b))) // hdrsize

		// @TODO the name INDEX could list multiple fonts,
		// but we just use the first one.
		cffGetIndex(&b) // name INDEX
		topdictidx = cffGetIndex(&b)
		topdict = cffIndexGet(topdictidx, 0)
		cffGetIndex(&b)
		info.gsubrs = cffGetIndex(&b)

		dictGetInts(&topdict, 17, charstrings[:])
		dictGetInts(&topdict, 0x100|6, cstype[:])
		dictGetInts(&topdict, 0x100|36, fdarrayoff[:])
		dictGetInts(&topdict, 0x100|37, fdselectoff[:])
		info.subrs = getSubrs(b, topdict)

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
			bufSeek(&b, int32(fdarrayoff[0]))
			info.fontdicts = cffGetIndex(&b)
			info.fdselect = bufRange(&b,
				int32(fdselectoff[0]),
				int32(len(b.data))-int32(fdselectoff[0]),
			)
		}

		bufSeek(&b, int32(charstrings[0]))
		info.charstrings = cffGetIndex(&b)
	}

	t = findTable(data, ufontstart, "maxp")
	if t != 0 {
		info.numGlyphs = int32(ttUSHORT(data[t+4:]))
	} else {
		info.numGlyphs = 0xffff
	}

	info.svg = -1

	// find a cmap encoding table we understand *now* to avoid searching
	// later. (todo: could make this installable)
	// the same regardless of glyph.
	numTables = uint32(ttUSHORT(data[cmap+2:]))
	info.index_map = 0
	for i = 0; i < numTables; i++ {
		var encoding_record = cmap + 4 + 8*i
		// find an encoding we understand:
		switch ttUSHORT(data[encoding_record:]) {
		case PlatformIDMicrosoft:
			switch ttUSHORT(data[encoding_record+2:]) {
			case MSEIDUnicodeBMP:
			case MSEIDUnicodeFull:
				// MS/Unicode
				info.index_map = int32(cmap + ttULONG(data[encoding_record+4:]))
			}
		case PlatformIDUnicode:
			// Mac/iOS has these
			// all the encodingIDs are unicode, so we don't bother to check it
			info.index_map = int32(cmap + ttULONG(data[encoding_record+4:]))
		}
	}
	if info.index_map == 0 {
		return 0
	}

	info.indexToLocFormat = int32(ttUSHORT(data[info.head+50:]))
	return 1
}

func setVertex(v *Vertex, typ uint8, x, y, cx, cy int32) {
	v.typ = typ
	v.x = int16(x)
	v.y = int16(y)
	v.cx = int16(cx)
	v.cy = int16(cy)
}

func getGlyfOffset(info *FontInfo, glyph_index int32) int32 {
	var g1, g2 int32

	if glyph_index >= info.numGlyphs {
		return -1
	}
	if info.indexToLocFormat >= 2 {
		return -1
	}

	if info.indexToLocFormat == 0 {
		g1 = info.glyf + int32(ttUSHORT(info.data[info.loca+glyph_index*2:])*2)
		g2 = info.glyf + int32(ttUSHORT(info.data[info.loca+glyph_index*2+2:])*2)
	} else {
		g1 = info.glyf + int32(ttULONG(info.data[info.loca+glyph_index*4:]))
		g2 = info.glyf + int32(ttULONG(info.data[info.loca+glyph_index*4+4:]))
	}

	if g1 == g2 {
		return -1
	}
	return g1
}

func closeShape(vertices []Vertex, num_vertices int32, was_off, start_off bool,
	sx, sy, scx, scy, cx, cy int32) int32 {
	if start_off {
		if was_off {
			setVertex(&vertices[num_vertices], vcurve, (cx+scx)>>1, (cy+scy)>>1, cx, cy)
			num_vertices++
		}
		setVertex(&vertices[num_vertices], vcurve, sx, sy, scx, scy)
		num_vertices++
	} else {
		if was_off {
			setVertex(&vertices[num_vertices], vcurve, sx, sy, cx, cy)
			num_vertices++
		} else {
			setVertex(&vertices[num_vertices], vline, sx, sy, 0, 0)
			num_vertices++
		}
	}
	return num_vertices
}

func getGlyphShapeTT(info *FontInfo, glyph_index int32) []Vertex {
	var numberOfContours int16
	var endPtsOfContours []byte
	var data = info.data
	var vertices []Vertex
	var numVertices int32
	var g = getGlyfOffset(info, glyph_index)
	if g < 0 {
		return nil
	}

	numberOfContours = ttSHORT(data[g:])
	if numberOfContours > 0 {
		var flags, flagcount uint8
		var ins, i, j, m, n, next_move, off int32
		var start_off, was_off bool
		var x, y, cx, cy, sx, sy, scx, scy int32
		var points []byte
		endPtsOfContours = data[g+10:]
		ins = int32(ttUSHORT(data[g+10+int32(numberOfContours)*2:]))
		points = data[g+10+int32(numberOfContours)*2+2+ins:]

		n = 1 + int32(ttUSHORT(endPtsOfContours[numberOfContours*2-2:]))
		vertices = make([]Vertex, 0, n+int32(numberOfContours)*2)
		next_move = 0
		flagcount = 0

		// in first pass, we load uninterpreted data into the allocated array
		// above, shifted to the end of the array so we won't overwrite it when
		// we create our final data starting from the front

		off = m - n

		// first load flags
		for i = 0; i < n; i++ {
			if flagcount == 0 {
				points = points[1:]
				flags = points[0]
				if flags&8 != 0 {
					points = points[1:]
					flagcount = points[0]
				}
			} else {
				flagcount--
			}
			vertices[off+i].typ = flags
		}

		// now load x coordinates
		x = 0
		for i = 0; i < n; i++ {
			flags = vertices[off+i].typ
			if flags&2 != 0 {
				points = points[1:]
				delta := int32(ttSHORT(points))
				if flags&16 == 0 {
					delta = -delta // ???
				}
				x += delta
			} else {
				if flags&16 == 0 {
					x += int32(ttSHORT(points))
					points = points[2:]
				}
			}
			vertices[off+i].x = int16(x)
		}

		// now load y coordinates
		y = 0
		for i = 0; i < n; i++ {
			flags = vertices[off+i].typ
			if flags&4 != 0 {
				points = points[1:]
				delta := int32(ttSHORT(points))
				if flags&32 == 0 {
					delta = -delta // ???
				}
				y += delta
			} else {
				if flags&32 == 0 {
					y += int32(ttSHORT(points))
					points = points[2:]
				}
			}
			vertices[off+i].y = int16(y)
		}

		// now convert them to our format
		for i = 0; i < n; i++ {
			flags = vertices[off+i].typ
			x = int32(vertices[off+i].x)
			y = int32(vertices[off+i].y)

			if next_move == i {
				if i != 0 {
					numVertices = closeShape(vertices, numVertices, was_off, start_off, sx, sy, scx, scy, cx, cy)
				}

				// now start the new one
				start_off = !(flags&1 == 0)
				if start_off {
					// if we start off with an off-curve point, then when we need to find a point on the curve
					// where we can start, and we need to save some state for when we wraparound.
					scx = x
					scy = y
					if vertices[off+i+1].typ&1 == 0 {
						// next point is also a curve point, so interpolate an on-point curve
						sx = (x + int32(vertices[off+i+1].x)) >> 1
						sy = (y + int32(vertices[off+i+1].y)) >> 1
					} else {
						// otherwise just use the next point as our start point
						sx = int32(vertices[off+i+1].x)
						sy = int32(vertices[off+i+1].y)
						i++ // we're using point i+1 as the starting point, so skip it
					}
				} else {
					sx = x
					sy = y
				}
				numVertices++
				setVertex(&vertices[numVertices], vmove, sx, sy, 0, 0)
				was_off = false
				next_move = 1 + int32(ttUSHORT(endPtsOfContours[j*2:]))
				j++
			} else {
				if flags&1 == 0 { // if it's a curve
					if was_off { // two off-curve control points in a row means interpolate an on-curve midpoint
						numVertices++
						setVertex(&vertices[numVertices], vcurve, (cx+x)>>1, (cy+y)>>1, cx, cy)
					}
					cx = x
					cy = y
					was_off = true
				} else {
					if was_off {
						numVertices++
						setVertex(&vertices[numVertices], vcurve, x, y, cx, cy)
					} else {
						numVertices++
						setVertex(&vertices[numVertices], vline, x, y, 0, 0)
					}
					was_off = false
				}
			}
		}
		numVertices = closeShape(vertices, numVertices, was_off, start_off, sx, sy, scx, scy, cx, cy)
	} else if numberOfContours < 0 {
		// Compound shapes.
		var more = true
		var comp = data[g+10:]
		numVertices = 0
		vertices = nil
		for more {
			var flags, gidx uint16
			var comp_verts, tmp []Vertex
			var mtx = [...]float32{1, 0, 0, 1, 0, 0}
			var m, n float32

			flags = uint16(ttSHORT(comp))
			comp = comp[2:]
			glyph_index = int32(ttSHORT(comp))
			comp = comp[2:]

			if flags&2 != 0 { // XY values
				if flags&1 != 0 {
					mtx[4] = float32(ttSHORT(comp))
					comp = comp[2:]
					mtx[5] = float32(ttSHORT(comp))
					comp = comp[2:]
				} else {
					mtx[4] = float32(ttCHAR(comp[0]))
					comp = comp[1:]
					mtx[5] = float32(ttCHAR(comp[0]))
					comp = comp[1:]
				}
			} else {
				panic("TODO handle matching point")
			}

			switch {
			case flags&(1<<3) != 0: // WE_HAVE_A_SCALE
				scale := float32(ttSHORT(comp)) / 16384
				mtx[0] = scale
				mtx[3] = scale
				comp = comp[2:]
			case flags&(1<<6) != 0: // WE_HAVE_AN_X_AND_YSCALE
				mtx[0] = float32(ttSHORT(comp)) / 16384
				comp = comp[2:]
				mtx[3] = float32(ttSHORT(comp)) / 16384
				comp = comp[2:]
			case flags&(1<<7) != 0: // WE_HAVE_A_TWO_BY_TWO
				mtx[0] = float32(ttSHORT(comp)) / 16384
				comp = comp[2:]
				mtx[1] = float32(ttSHORT(comp)) / 16384
				comp = comp[2:]
				mtx[2] = float32(ttSHORT(comp)) / 16384
				comp = comp[2:]
				mtx[3] = float32(ttSHORT(comp)) / 16384
				comp = comp[2:]
			}

			// Find transformation scales.
			m = (float32)(math.Sqrt(float64(mtx[0]*mtx[0] + mtx[1]*mtx[1])))
			n = (float32)(math.Sqrt(float64(mtx[2]*mtx[2] + mtx[3]*mtx[3])))

			// Get indexed glyph.
			comp_verts = GetGlyphShape(info, int32(gidx))
			if len(comp_verts) > 0 {
				// Transform vertices
				for i := range comp_verts {
					v := &comp_verts[i]
					x, y := float32(v.x), float32(v.y)
					v.x = int16(m * (mtx[0]*x + mtx[2]*y + mtx[4]))
					v.y = int16(n * (mtx[1]*x + mtx[3]*y + mtx[5]))
					x, y = float32(v.cx), float32(v.cy)
					v.cx = int16(m * (mtx[0]*x + mtx[2]*y + mtx[4]))
					v.cy = int16(n * (mtx[1]*x + mtx[3]*y + mtx[5]))
				}
				// Append vertices.
				tmp = make([]Vertex, len(comp_verts))
				copy(tmp, comp_verts)
				vertices = tmp
				numVertices += int32(len(comp_verts))
			}
			// More components ?
			more = flags&(1<<5) != 0
		}
	} else {
		// numberOfCounters == 0, do nothing
	}

	return vertices
}

func cidGetGlyphSubrs(info *FontInfo, glyph_index int32) buf {
	fdselect := info.fdselect
	var nranges, start, end, v, fmt, fdselector int32
	fdselector = -1

	var i int32

	bufSeek(&fdselect, 0)
	fmt = int32(bufGet8(&fdselect))
	if fmt == 0 {
		// untested
		bufSkip(&fdselect, glyph_index)
		fdselector = int32(bufGet8(&fdselect))
	} else if fmt == 3 {
		nranges = int32(bufGet16(&fdselect))
		start = int32(bufGet16(&fdselect))
		for i = 0; i < nranges; i++ {
			v = int32(bufGet8(&fdselect))
			end = int32(bufGet16(&fdselect))
			if glyph_index >= start && glyph_index < end {
				fdselector = v
				break
			}
			start = end
		}
	}
	if fdselector == -1 {
		return buf{}
	}
	return getSubrs(info.cff, cffIndexGet(info.fontdicts, fdselector))
}

func getGlyphShapeT2(info *FontInfo, glyph_index int32) []Vertex {
	// runs the charstring twice, once to count and once to output (to avoid realloc)
	count_ctx := csctx{bounds: true}
	output_ctx := csctx{bounds: false}

	var err error
	if err = runCharString(info, glyph_index, &count_ctx); err == nil {
		vertices := make([]Vertex, count_ctx.num_vertices)
		output_ctx.vertices = vertices
		if err = runCharString(info, glyph_index, &output_ctx); err == nil {
			if !(output_ctx.num_vertices == count_ctx.num_vertices) {
				panic("!(output_ctx.num_vertices == count_ctx.num_vertices)")
			}
			return vertices
		}
	}

	//TODO PORTING need a more flexible way to return this error up.
	if err != nil {
		log.Println(err)
	}

	return nil
}

func getGlyphInfoT2(info *FontInfo, glyph_index int32) (x0, y0, x1, y1 int32, err error) {
	var c = csctx{bounds: true}
	err = runCharString(info, glyph_index, &c)
	if err != nil {
		return
	}
	return c.min_x, c.min_y, c.max_x, c.max_y, nil
}

func getGlyphKernInfoAdvance(info *FontInfo, glyph1, glyph2 int32) int32 {
	var data = info.data[info.kern:]
	var needle, straw uint32
	var l, r, m int32

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
	r = int32(ttUSHORT(data[10:]) - 1)
	needle = uint32(glyph1)<<16 | uint32(glyph2)
	for l <= r {
		m = (l + r) >> 1
		straw = ttULONG(data[18+(m*6):])
		if needle < straw {
			r = m - 1
		} else if needle > straw {
			l = m + 1
		} else {
			return int32(ttSHORT(data[22+(m*6):]))
		}
	}
	return 0
}

func getCoverageIndex(coverageTable []byte, glyph int32) int32 {
	var converageFormat = ttUSHORT(coverageTable)
	switch converageFormat {
	case 1:
		var glyphCount = ttUSHORT(coverageTable[2:])
		// Binary search.
		var l, r, m int32 = 0, int32(glyphCount - 1), 0
		var straw, needle int32 = 0, glyph
		for l <= r {
			var glyphArray = coverageTable[4:]
			var glyphID int32
			m = (l + r) >> 1
			glyphID = int32(ttUSHORT(glyphArray[2*m:]))
			straw = glyphID
			if needle < straw {
				r = m - 1
			} else if needle > straw {
				l = m + 1
			} else {
				return m
			}
		}
	case 2:
		var rangeCount = ttUSHORT(coverageTable[2:])
		var rangeArray = coverageTable[4:]
		// Binary search.
		var l, r, m int32 = 0, int32(rangeCount - 1), 0
		var strawStart, strawEnd, needle int32 = 0, 0, glyph
		for l <= r {
			m = (l + r) >> 1
			var rangeRecord = rangeArray[6*m:]
			strawStart = int32(ttUSHORT(rangeRecord))
			strawEnd = int32(ttUSHORT(rangeRecord[2:]))
			if needle < strawStart {
				r = m - 1
			} else if needle > strawEnd {
				l = m + 1
			} else {
				var startCoverageIndex = ttUSHORT(rangeRecord[4:])
				return int32(startCoverageIndex) + glyph - strawStart
			}
		}
	default:
		return -1 // unsupported
	}
	return -1
}

func getGlyphClass(classDefTable []byte, glyph int32) int32 {
	var classDefFormat = ttUSHORT(classDefTable)
	switch classDefFormat {
	case 1:
		var startGlyphID = ttUSHORT(classDefTable[2:])
		var glyphCount = ttUSHORT(classDefTable[4:])
		var classDef1ValueArray = classDefTable[6:]

		if glyph >= int32(startGlyphID) && glyph < int32(startGlyphID+glyphCount) {
			return int32(ttUSHORT(classDef1ValueArray[2*(glyph-int32(startGlyphID)):]))
		}
	case 2:
		var classRangeCount = ttUSHORT(classDefTable[2:])
		var classRangeRecords = classDefTable[4:]

		// Binary search.
		var l, r, m int32 = 0, int32(classRangeCount - 1), 0
		var strawStart, strawEnd, needle int32 = 0, 0, glyph
		for l <= r {
			m = (l + r) >> 1
			var classRangeRecord = classRangeRecords[6*m:]
			strawStart = int32(ttUSHORT(classRangeRecord))
			strawEnd = int32(ttUSHORT(classRangeRecord[2:]))
			if needle < strawStart {
				r = m - 1
			} else if needle > strawEnd {
				l = m + 1
			} else {
				return int32(ttUSHORT(classRangeRecord[4:]))
			}
		}
	default:
		return -1 // unsupported
	}
	return -1
}

func getGlyphGPOSInfoAdvance(info *FontInfo, glyph1, glyph2 int32) int32 {
	if info.gpos == 0 {
		return 0
	}
	data := info.data[info.gpos:]

	if ttUSHORT(data[0:]) != 1 {
		return 0 // Major version 1
	}
	if ttUSHORT(data[2:]) != 0 {
		return 0 // Minor version 0
	}

	lookupListOffset := ttUSHORT(data[8:])
	lookupList := data[lookupListOffset:]
	lookupCount := ttUSHORT(lookupList)

	for i := uint16(0); i < lookupCount; i++ {
		lookupOffset := ttUSHORT(lookupList[2+2*i:])
		lookupTable := lookupList[lookupOffset:]

		lookupType := ttUSHORT(lookupTable)
		subTableCount := ttUSHORT(lookupTable[4:])
		subTableOffsets := lookupTable[6:]

		if lookupType != 2 { // Pair Adjustment Positioning Subtable
			continue
		}

		for sti := uint16(0); sti < subTableCount; sti++ {
			subtableOffset := ttUSHORT(subTableOffsets[2*sti:])
			table := lookupTable[subtableOffset:]

			posFormat := ttUSHORT(table)
			coverageOffset := ttUSHORT(table[2:])
			coverageIndex := getCoverageIndex(table[coverageOffset:], glyph1)
			if coverageIndex == -1 {
				continue
			}

			switch posFormat {
			case 1:
				valueFormat1 := ttUSHORT(table[4:])
				valueFormat2 := ttUSHORT(table[6:])
				if valueFormat1 == 4 && valueFormat2 == 0 { // Support more formats?
					var valueRecordPairSizeInBytes uint16 = 2
					pairSetCount := ttUSHORT(table[8:])
					pairPosOffset := ttUSHORT(table[10+2*coverageIndex:])
					pairValueTable := table[pairPosOffset:]
					pairValueCount := ttUSHORT(pairValueTable)
					pairValueArray := pairValueTable[2:]

					if coverageIndex >= int32(pairSetCount) {
						return 0
					}

					needle := glyph2
					r := pairValueCount - 1
					var l uint16 = 0

					// Binary search.
					for l <= r {
						m := (l + r) >> 1
						pairValue := pairValueArray[(2+valueRecordPairSizeInBytes)*m:]
						secondGlyph := ttUSHORT(pairValue)
						straw := secondGlyph
						if needle < int32(straw) {
							r = m - 1
						} else if needle > int32(straw) {
							l = m + 1
						} else {
							xAdvance := ttSHORT(pairValue[2:])
							return int32(xAdvance)
						}
					}
				} else {
					return 0
				}
			case 2:
				valueFormat1 := ttUSHORT(table[4:])
				valueFormat2 := ttUSHORT(table[6:])
				if valueFormat1 == 4 && valueFormat2 == 0 { // Support more formats?
					classDef1Offset := ttUSHORT(table[8:])
					classDef2Offset := ttUSHORT(table[10:])
					glyph1class := getGlyphClass(table[classDef1Offset:], glyph1)
					glyph2class := getGlyphClass(table[classDef2Offset:], glyph2)

					class1Count := ttUSHORT(table[12:])
					class2Count := ttUSHORT(table[14:])

					if glyph1class < 0 || glyph1class >= int32(class1Count) {
						return 0 // malformed
					}
					if glyph2class < 0 || glyph2class >= int32(class2Count) {
						return 0 // malformed
					}

					class1Records := table[16:]
					class2Records := class1Records[2*(glyph1class*int32(class2Count)):]
					xAdvance := ttSHORT(class2Records[2*glyph2class:])
					return int32(xAdvance)
				} else {
					return 0
				}
			default:
				return 0 // Unsupported position format
			}
		}
	}

	return 0
}
