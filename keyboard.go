package imgui

// Pass in translated ASCII characters for text input.
// - with glfw you can get those from the callback set in glfwSetCharCallback()
// - on Windows you can get those using ToAscii+keyboard state, or via the WM_CHAR message
func (io *ImGuiIO) AddInputCharacter(c rune) {
	if c != 0 {
		if c <= IM_UNICODE_CODEPOINT_MAX {
			io.InputQueueCharacters = append(io.InputQueueCharacters, c)
		} else {
			io.AddInputCharacter(IM_UNICODE_CODEPOINT_INVALID)
		}
	}
}

func (io *ImGuiIO) AddInputCharacters(chars string) {
	for _, c := range chars {
		io.AddInputCharacter(c)
	}
}

// Clear the text input buffer manually
func (io *ImGuiIO) ClearInputCharacters() {
	io.InputQueueCharacters = io.InputQueueCharacters[:0]
}

func GetMergedKeyModFlags() ImGuiKeyModFlags {
	var g = GImGui
	var key_mod_flags ImGuiKeyModFlags = ImGuiKeyModFlags_None
	if g.IO.KeyCtrl {
		key_mod_flags |= ImGuiKeyModFlags_Ctrl
	}
	if g.IO.KeyShift {
		key_mod_flags |= ImGuiKeyModFlags_Shift
	}
	if g.IO.KeyAlt {
		key_mod_flags |= ImGuiKeyModFlags_Alt
	}
	if g.IO.KeySuper {
		key_mod_flags |= ImGuiKeyModFlags_Super
	}
	return key_mod_flags
}

func IsKeyPressed(user_key_index int, repeat bool /*= true*/) bool {
	var g = GImGui
	if user_key_index < 0 {
		return false
	}
	IM_ASSERT(user_key_index >= 0 && user_key_index < int(len(g.IO.KeysDown)))
	var t float = g.IO.KeysDownDuration[user_key_index]
	if t == 0.0 {
		return true
	}
	if repeat && t > g.IO.KeyRepeatDelay {
		return GetKeyPressedAmount(user_key_index, g.IO.KeyRepeatDelay, g.IO.KeyRepeatRate) > 0
	}
	return false
} // was key pressed (went from !Down to Down)? if repeat=true, uses io.KeyRepeatDelay / KeyRepeatRate
