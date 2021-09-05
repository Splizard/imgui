package imgui

func Separator() { panic("not implemented") }

func SameLine(offset, spacing float32) { panic("not implemented") }

func NewLine() { panic("not implemented") }

func Spacing() { panic("not implemented") }

func Dummy(size Vec2) { panic("not implemented") }

func Indent(indent float32)   { panic("not implemented") }
func Unindent(indent float32) { panic("not implemented") }

func BeginGroup() { panic("not implemented") }
func EndGroup()   { panic("not implemented") }

func GetCursorPos() Vec2      { panic("not implemented"); return Vec2{} }
func GetCursorPosX() float32  { panic("not implemented"); return 0 }
func GetCursorPosY() float32  { panic("not implemented"); return 0 }
func SetCursorPos(p Vec2)     { panic("not implemented") }
func SetCursorPosX(x float32) { panic("not implemented") }
func SetCursorPosY(y float32) { panic("not implemented") }

func GetCursorStartPos() Vec2  { panic("not implemented"); return Vec2{} }
func GetCursorScreenPos() Vec2 { panic("not implemented"); return Vec2{} }

func SetCursorScreenPos(pos Vec2)           { panic("not implemented") }
func AlignTextToFramePadding()              { panic("not implemented") }
func GetTextLineHeight() float32            { panic("not implemented"); return 0 }
func GetTextLineHeightWithSpacing() float32 { panic("not implemented"); return 0 }

func GetFrameHeight() float32            { panic("not implemented"); return 0 }
func GetFrameHeightWithSpacing() float32 { panic("not implemented"); return 0 }
