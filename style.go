package imgui

type Style struct{}

func StyleColorsDark(dst *Style)    { panic("not implemented") }
func StyleColorsLight(dst *Style)   { panic("not implemented") }
func StyleColorsClassic(dst *Style) { panic("not implemented") }

func GetFont() *Font {
	panic("not implemented")
	return nil
}

func GetFontSize() float32 {
	panic("not implemented")
	return 0
}

func GetFontTexUvWhitePixel() Vec2 {
	panic("not implemented")
	return Vec2{}
}

func GetColorU32(col uint32) uint32 { panic("not implemented"); return 0 }

func GetStyleColorVec4(idx Col) Vec4 {
	panic("not implemented")
	return Vec4{}
}
