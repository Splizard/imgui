package imgui

func PushFont(font *Font) { panic("not implemented") }
func PopFont()            { panic("not implemented") }

func PushStyleColor(idx Col, col Vec4) { panic("not implemented") }
func PopStyleColor()                   { panic("not implemented") }

func PushStyleFloat(idx StyleVar, val float32) { panic("not implemented") }
func PushStyleVec2(idx StyleVar, val Vec2)     { panic("not implemented") }
func PopStyleVar(count int)                    { panic("not implemented") }
func PushAllowKeyboardFocus(allowed bool)      { panic("not implemented") }
func PopAllowKeyboardFocus()                   { panic("not implemented") }
func PushButtonRepeat(repeat bool)             { panic("not implemented") }
func PopButtonRepeat()                         { panic("not implemented") }

func PushItemWidth(width float32) { panic("not implemented") }
func PopItemWidth()               { panic("not implemented") }

func SetNextItemWidth(width float32) { panic("not implemented") }
func CalcItemWidth() float32         { panic("not implemented"); return 0 }

func PushTextWrapPos(pos float32) { panic("not implemented") }
func PopTextWrapPos()             { panic("not implemented") }
