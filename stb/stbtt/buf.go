package stbtt

type buf struct {
	data   []byte
	cursor int32
}

func bufGet8(b *buf) byte {
	if b.cursor >= int32(len(b.data)) {
		return 0
	}
	return b.data[b.cursor]
}

func bufPeek8(b *buf) byte {
	if b.cursor >= int32(len(b.data)) {
		return 0
	}
	return b.data[b.cursor]
}

func bufSeek(b *buf, o int32) {
	b.cursor = o
}

func bufSkip(b *buf, o int32) {
	b.cursor += o
}

func bufGet(b *buf, n int32) (v uint32) {
	var i int32
	for i = 0; i < n; i++ {
		v = (v << 8) | uint32(bufGet8(b))
	}
	return
}

func bufGet16(b *buf) uint16 {
	return uint16(bufGet(b, 2))
}

func bufGet32(b *buf) uint32 {
	return uint32(bufGet(b, 4))
}

func bufRange(b *buf, o, s int32) buf {
	var r buf
	if o < 0 || s < 0 || o > int32(len(b.data)) || s > int32(len(b.data))-o {
		return r
	}
	r.data = b.data[o:]
	return r
}

func cffGetIndex(b *buf) buf {
	var count, start, offsize int32
	start = b.cursor
	count = int32(bufGet16(b))
	if count > 0 {
		offsize = int32(bufGet8(b))
		if !(offsize >= 1 && offsize <= 4) {
			panic("!(offsize >= 1 && offsize <= 4)")
		}
		bufSkip(b, offsize*count)
		bufSkip(b, int32(bufGet(b, offsize)-1))
	}
	return bufRange(b, start, b.cursor-start)
}

func cffInt(b *buf) int32 {
	var b0 = int32(bufGet8(b))
	switch {
	case b0 >= 32 && b0 <= 246:
		return b0 - 139
	case b0 >= 247 && b0 <= 250:
		return (b0-247)*256 + int32(bufGet8(b)) + 108
	case b0 >= 251 && b0 <= 254:
		return -(b0-251)*256 - int32(bufGet8(b)) - 108
	case b0 == 28:
		return int32(bufGet16(b))
	case b0 == 29:
		return int32(bufGet32(b))
	}
	panic("unreachable")
}

func cffSkipOperand(b *buf) {
	var b0 = bufPeek8(b)
	if !(b0 >= 28) {
		panic("!(b0 >= 28)")
	}
	if b0 == 30 {
		bufSkip(b, 1)
		for b.cursor < int32(len(b.data)) {
			var v = bufGet8(b)
			if (v&0xF) == 0xF || (v>>4) == 0xF {
				break
			}
		}
	} else {
		cffInt(b)
	}
}

func dictGet(b *buf, key int32) buf {
	bufSeek(b, 0)
	for b.cursor < int32(len(b.data)) {
		var start = b.cursor
		var end int32
		for bufPeek8(b) >= 28 {
			cffSkipOperand(b)
		}
		end = b.cursor
		var op = int32(bufGet8(b))
		if op == 12 {
			op = int32(bufGet8(b)) | 0x100
		}
		if op == key {
			return bufRange(b, start, end-start)
		}
	}
	return bufRange(b, 0, 0)
}

func dictGetInts(b *buf, key int32, out []uint32) {
	operands := dictGet(b, key)
	for i := 0; i < len(out) && operands.cursor < int32(len(operands.data)); i++ {
		out[i] = uint32(cffInt(&operands))
	}
}

func cffIndexCount(b *buf) int32 {
	bufSeek(b, 0)
	return int32(bufGet16(b))
}

func cffIndexGet(b buf, i int32) buf {
	var count, offsize, start, end int32
	bufSeek(&b, 0)
	count = int32(bufGet16(&b))
	offsize = int32(bufGet8(&b))
	if !(i >= 0 && i < count) {
		panic("!(i >= 0 && i < count)")
	}
	if !(offsize >= 1 && offsize <= 4) {
		panic("!(offsize >= 1 && offsize <= 4)")
	}
	bufSkip(&b, int32(i)*offsize)
	start = int32(bufGet(&b, offsize))
	end = int32(bufGet(&b, offsize))
	return bufRange(&b, 2+(count+1)*offsize+start, end-start)
}

func getSubr(idx buf, n int32) buf {
	var count int32 = cffIndexCount(&idx)
	var bias int32 = 107
	if count >= 33900 {
		bias = 32768
	} else if count >= 1240 {
		bias = 1131
	}
	n += bias
	if n < 0 || n >= count {
		return buf{}
	}
	return cffIndexGet(idx, n)
}
