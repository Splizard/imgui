package imgui

// move content position toward the right, by indent_w, or style.IndentSpacing if indent_w <= 0
func Indent(indent_w float) {
	g := GImGui
	var window = GetCurrentWindow()

	if indent_w == 0 {
		indent_w = g.Style.IndentSpacing
	}

	window.DC.Indent.x += indent_w
	window.DC.CursorPos.x = window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x
}

// move content position back to the left, by indent_w, or style.IndentSpacing if indent_w <= 0
func Unindent(indent_w float) {
	g := GImGui
	var window = GetCurrentWindow()

	if indent_w == 0 {
		indent_w = g.Style.IndentSpacing
	}

	window.DC.Indent.x -= indent_w
	window.DC.CursorPos.x = window.Pos.x + window.DC.Indent.x + window.DC.ColumnsOffset.x
}
