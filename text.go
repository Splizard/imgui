package imgui

// push word-wrapping position for Text*() commands. < 0.0: no wrapping; 0.0: wrap to end of window (or column)  {panic("not implemented")} > 0.0: wrap at 'wrap_pos_x' position in window local space
func PushTextWrapPos(wrap_local_pos_x float) {
	var window = GetCurrentWindow()
	window.DC.TextWrapPosStack = append(window.DC.TextWrapPosStack, window.DC.TextWrapPos)
	window.DC.TextWrapPos = wrap_local_pos_x
}

func PopTextWrapPos() {
	var window = GetCurrentWindow()
	window.DC.TextWrapPos = window.DC.TextWrapPosStack[len(window.DC.TextWrapPosStack)-1]
	window.DC.TextWrapPosStack = window.DC.TextWrapPosStack[:len(window.DC.TextWrapPosStack)-1]
}
