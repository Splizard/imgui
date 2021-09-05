package stbtt

import (
	"errors"
	"math"
)

func runCharString(info *FontInfo, glyph_index int32, c *csctx) error {
	var in_header = true
	var maskbits, subr_stack_height, sp int32
	var has_subrs, clear_stack bool
	var s [48]float32
	var subr_stack [10]buf
	var b buf
	var subrs = info.subrs
	var f float32

	// this currently ignores the initial width value, which isn't needed if we have hmtx
	b = cffIndexGet(info.charstrings, glyph_index)
	for b.cursor < int32(len(b.data)) {
		i := 0
		clear_stack = true
		b0 := bufGet8(&b)
		switch b0 {
		// @TODO implement hinting
		case 0x13: // hintmask
			fallthrough
		case 0x14: // cntrmask
			if in_header {
				maskbits += sp / 2
			}
			in_header = false
			bufSkip(&b, (maskbits+7)/8)
		case 0x01: // hstem
			fallthrough
		case 0x03: // vstem
			fallthrough
		case 0x12: // hstemhm
			fallthrough
		case 0x17: // vstemhm
			maskbits += sp / 2

		case 0x15: // rmoveto
			in_header = false
			if sp < 2 {
				return errors.New("rmoveto stack")
			}
			csctxRmoveTo(c, s[sp-2], s[sp-1])
		case 0x04: // vmoveto
			in_header = false
			if sp < 1 {
				return errors.New("vmoveto stack")
			}
			csctxRmoveTo(c, 0, s[sp-1])
		case 0x16: // hmoveto
			in_header = false
			if sp < 1 {
				return errors.New("hmoveto stack")
			}
			csctxRmoveTo(c, s[sp-1], 0)
		case 0x05: // rlineto
			if sp < 2 {
				return errors.New("rlineto stack")
			}
			for ; int32(i+1) < sp; i++ {
				csctxRlineTo(c, s[i], s[i+1])
			}

		// hlineto/vlineto and vhcurveto/hvcurveto alternate horizontal and vertical
		// starting from a different place.

		case 0x07: // vlineto
			if sp < 1 {
				return errors.New("vlineto stack")
			}
			if int32(i) >= sp {
				break
			}
			csctxRlineTo(c, 0, s[i])
			i++
			for {
				if int32(i) >= sp {
					break
				}
				csctxRlineTo(c, 0, s[i])
				i++
				if int32(i) >= sp {
					break
				}
				csctxRlineTo(c, s[i], 0)
				i++
			}
		case 0x06: // hlineto
			if sp < 1 {
				return errors.New("hlineto stack")
			}
			for {
				if int32(i) >= sp {
					break
				}
				csctxRlineTo(c, s[i], 0)
				i++
				if int32(i) >= sp {
					break
				}
				csctxRlineTo(c, 0, s[i])
				i++
			}

		case 0x1F: //hvcurveto
			if sp < 4 {
				return errors.New("hvcurveto stack")
			}
			for {
				var last float32
				if sp-int32(i) == 5 {
					last = s[i+4]
				}

				if int32(i+3) >= sp {
					break
				}
				csctxRccurveTo(c, s[i], 0, s[i+1], s[i+2], last, s[i+3])
				i += 4
				if int32(i+3) >= sp {
					break
				}
				csctxRccurveTo(c, 0, s[i], s[i+1], s[i+2], s[i+3], last)
				i += 4
			}

		case 0x1E: // vhcurveto
			if sp < 4 {
				return errors.New("vhcurveto stack")
			}
			for {
				var last float32
				if sp-int32(i) == 5 {
					last = s[i+4]
				}
				if int32(i+3) >= sp {
					break
				}
				csctxRccurveTo(c, 0, s[i], s[i+1], s[i+2], s[i+3], last)
				i += 4
				if int32(i+3) >= sp {
					break
				}
				csctxRccurveTo(c, s[i], 0, s[i+1], s[i+2], last, s[i+3])
				i += 4
			}

		case 0x08: // rrcurveto
			if sp < 6 {
				return errors.New("rrcurveto stack")
			}
			for ; int32(i+5) < sp; i += 6 {
				csctxRccurveTo(c, s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5])
			}

		case 0x18: //rcurveline
			if sp < 8 {
				return errors.New("rcurveline stack")
			}
			for ; int32(i+5) < sp; i += 6 {
				csctxRccurveTo(c, s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5])
			}
			if i+1 >= int(sp) {
				return errors.New("rcurveline stack")
			}
			csctxRlineTo(c, s[i], s[i+1])
		case 0x19: //rlinecurve

			if sp < 8 {
				return errors.New("rlinecurve stack")
			}
			for ; int32(i+1) < sp-6; i += 2 {
				csctxRlineTo(c, s[i], s[i+1])
			}
			if i+5 >= int(sp) {
				return errors.New("rlinecurve stack")
			}
			csctxRccurveTo(c, s[i], s[i+1], s[i+2], s[i+3], s[i+4], s[i+5])
		case 0x1A: //vvcurveto
			if sp < 4 {
				return errors.New("vvcurveto stack")
			}
			fallthrough
		case 0x1B: //hhcurveto
			if sp < 4 {
				return errors.New("hhcurveto stack")
			}
			var f float32
			if sp&1 != 0 {
				f = s[i]
				i++
			}
			for ; int32(i+3) < sp; i += 4 {
				if b0 == 0x1B {
					csctxRccurveTo(c, s[i], f, s[i+1], s[i+2], s[i+3], 0)
				} else {
					csctxRccurveTo(c, s[i], 0, s[i+1], f, s[i+2], s[i+3])
				}
				f = 0
			}
		case 0x0A: // callsubr
			if !has_subrs {
				if len(info.fdselect.data) > 0 {
					subrs = cidGetGlyphSubrs(info, glyph_index)
				}
				has_subrs = true
			}
			if sp < 1 {
				return errors.New("callsubr stack")
			}
			fallthrough
		case 0x1D: // callgsubr
			if sp < 1 {
				return errors.New("callgsubr stack")
			}
			sp--
			v := int32(s[sp])
			if subr_stack_height >= 10 {
				return errors.New("subr recursion limit reached")
			}
			subr_stack_height++
			subr_stack[subr_stack_height] = b
			arg := info.gsubrs
			if b0 == 0x0A {
				arg = subrs
			}
			b = getSubr(arg, v)
			if len(b.data) == 0 {
				return errors.New("subr not found")
			}
			b.cursor = 0
			clear_stack = false
		case 0x0B: // return
			if subr_stack_height <= 0 {
				return errors.New("return outside subroutine")
			}
			subr_stack_height--
			b = subr_stack[subr_stack_height]
			clear_stack = false
		case 0x0E: // endchar
			csctxCloseShape(c)
			return nil
		case 0x0C: // two-byte escape
			var dx1, dx2, dx3, dx4, dx5, dx6, dy1, dy2, dy3, dy4, dy5, dy6 float32
			var dx, dy float32
			var b1 int32 = int32(bufGet8(&b))
			switch b1 {
			// @TODO These "flex" implementations ignore the flex-depth and resolution,
			// and always draw beziers.
			case 0x22: // hflex
				if sp < 7 {
					return errors.New("hflex stack")
				}
				dx1 = s[0]
				dx2 = s[1]
				dy2 = s[2]
				dx3 = s[3]
				dx4 = s[4]
				dx5 = s[5]
				dx6 = s[6]
				csctxRccurveTo(c, dx1, 0, dx2, dy2, dx3, 0)
				csctxRccurveTo(c, dx4, 0, dx5, -dy2, dx6, 0)
			case 0x23: // flex
				if sp < 13 {
					return errors.New("flex stack")
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
				csctxRccurveTo(c, dx1, dy1, dx2, dy2, dx3, dy3)
				csctxRccurveTo(c, dx4, dy4, dx5, dy5, dx6, dy6)
			case 0x24: // hflex1
				if sp < 9 {
					return errors.New("hflex1 stack")
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
				csctxRccurveTo(c, dx1, dy1, dx2, dy2, dx3, 0)
				csctxRccurveTo(c, dx4, 0, dx5, dy5, dx6, -(dy1 + dy2 + dy5))
			case 0x25: // flex1
				if sp < 11 {
					return errors.New("flex1 stack")
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
				if math.Abs(float64(dx)) >= 0.01 || math.Abs(float64(dy)) >= 0.01 {
					dy6 = -dy
				} else {
					dx6 = -dx
				}
				csctxRccurveTo(c, dx1, dy1, dx2, dy2, dx3, dy3)
				csctxRccurveTo(c, dx4, dy4, dx5, dy5, dx6, dy6)

			default:
				return errors.New("unknown two-byte escape")
			}
		default:
			if b0 != 255 && b0 != 28 && b0 != 32 {
				return errors.New("reserved operator")
			}

			// push immediate
			if b0 == 255 {
				f = float32(bufGet32(&b)) / 0x10000
			} else {
				bufSkip(&b, -1)
				f = float32(cffInt(&b))
			}
			if sp >= 48 {
				return errors.New("push stack overflow")
			}
			sp++
			s[sp] = f
			clear_stack = false
		}
		if clear_stack {
			sp = 0
		}
	}
	return errors.New("no endchar")
}
