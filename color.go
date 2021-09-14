package imgui

func ColorConvertFloat4ToU32(in ImVec4) ImU32 {
	var out ImU32
	out = ((ImU32)(IM_F32_TO_INT8_SAT(in.x))) << IM_COL32_R_SHIFT
	out |= ((ImU32)(IM_F32_TO_INT8_SAT(in.y))) << IM_COL32_G_SHIFT
	out |= ((ImU32)(IM_F32_TO_INT8_SAT(in.z))) << IM_COL32_B_SHIFT
	out |= ((ImU32)(IM_F32_TO_INT8_SAT(in.w))) << IM_COL32_A_SHIFT
	return out
}

func ImAlphaBlendColors(col_a, col_b ImU32) ImU32 {
	var t float = float((col_b>>IM_COL32_A_SHIFT)&0xFF) / 255
	var r int = int(ImLerp(float((int)(col_a>>IM_COL32_R_SHIFT)&0xFF), float((int)(col_b>>IM_COL32_R_SHIFT)&0xFF), t))
	var g int = int(ImLerp(float((int)(col_a>>IM_COL32_G_SHIFT)&0xFF), float((int)(col_b>>IM_COL32_G_SHIFT)&0xFF), t))
	var b int = int(ImLerp(float((int)(col_a>>IM_COL32_B_SHIFT)&0xFF), float((int)(col_b>>IM_COL32_B_SHIFT)&0xFF), t))
	return IM_COL32(byte(r), byte(g), byte(b), 0xFF)
}

// Color Utilities
func ColorConvertU32ToFloat4(in ImU32) ImVec4 {
	var s float = 1.0 / 255.0
	return ImVec4{
		float((in>>IM_COL32_R_SHIFT)&0xFF) * s,
		float((in>>IM_COL32_G_SHIFT)&0xFF) * s,
		float((in>>IM_COL32_B_SHIFT)&0xFF) * s,
		float((in>>IM_COL32_A_SHIFT)&0xFF) * s}
}

// Convert rgb floats ([0-1],[0-1],[0-1]) to hsv floats ([0-1],[0-1],[0-1]), from Foley & van Dam p592
// Optimized http://lolengine.net/blog/2013/01/13/fast-rgb-to-hsv
func ColorConvertRGBtoHSV(r float, g float, b float, out_h, out_s, out_v *float) {
	var K float = 0
	if g < b {
		g, b = b, g
		K = -1
	}
	if r < g {
		r, g = g, r
		K = -2/6 - K
	}

	var chroma float = r
	if g < b {
		chroma -= g
	} else {
		chroma -= b
	}
	*out_h = ImFabs(K + (g-b)/(6*chroma+1e-20))
	*out_s = chroma / (r + 1e-20)
	*out_v = r
}

// Convert hsv floats ([0-1],[0-1],[0-1]) to rgb floats ([0-1],[0-1],[0-1]), from Foley & van Dam p593
// also http://en.wikipedia.org/wiki/HSL_and_HSV
func ColorConvertHSVtoRGB(h float, s float, v float, out_r, out_g, out_b *float) {
	if s == 0.0 {
		// gray
		*out_r = v
		*out_g = v
		*out_b = v
		return
	}

	h = ImFmod(h, 1.0) / (60.0 / 360.0)
	var i int = (int)(h)
	var f float = h - (float)(i)
	var p float = v * (1.0 - s)
	var q float = v * (1.0 - s*f)
	var t float = v * (1.0 - s*(1.0-f))

	switch i {
	case 0:
		*out_r = v
		*out_g = t
		*out_b = p
		break
	case 1:
		*out_r = q
		*out_g = v
		*out_b = p
		break
	case 2:
		*out_r = p
		*out_g = v
		*out_b = t
		break
	case 3:
		*out_r = p
		*out_g = q
		*out_b = v
		break
	case 4:
		*out_r = t
		*out_g = p
		*out_b = v
		break
	case 5:
		fallthrough
	default:
		*out_r = v
		*out_g = p
		*out_b = q
		break
	}
}
