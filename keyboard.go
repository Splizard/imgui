package imgui

// Inputs Utilities: Keyboard
// - For 'user_key_index int' you can use your own indices/enums according to how your backend/engine stored them in io.KeysDown[].
// - We don't know the meaning of those value. You can use GetKeyIndex() to map a ImGuiKey_ value into the user index.

// == tab stop enable. Allow focusing using TAB/Shift-TAB, enabled by default but you can disable it for certain widgets
func PushAllowKeyboardFocus(allow_keyboard_focus bool) {
	PushItemFlag(ImGuiItemFlags_NoTabStop, !allow_keyboard_focus)
}
func PopAllowKeyboardFocus() {
	PopItemFlag()
}

// map ImGuiKey_* values into user's key index. == io.KeyMap[key]
func GetKeyIndex(imgui_key ImGuiKey) int {
	IM_ASSERT(imgui_key >= 0 && imgui_key < ImGuiKey_COUNT)
	var g = GImGui
	return g.IO.KeyMap[imgui_key]
}

// Note that dear imgui doesn't know the semantic of each entry of io.KeysDown[]!
// Use your own indices/enums according to how your backend/engine stored them into io.KeysDown[]!
// is key being held. == io.KeysDown[user_key_index].
func IsKeyDown(user_key_index int) bool {
	if user_key_index < 0 {
		return false
	}
	var g = GImGui
	IM_ASSERT(user_key_index >= 0 && user_key_index < int(len(g.IO.KeysDown)))
	return g.IO.KeysDown[user_key_index]
}

// was key released (went from Down to !Down)?
func IsKeyReleased(user_key_index int) bool {
	var g = GImGui
	if user_key_index < 0 {
		return false
	}
	IM_ASSERT(user_key_index >= 0 && user_key_index < int(len(g.IO.KeysDown)))
	return g.IO.KeysDownDurationPrev[user_key_index] >= 0.0 && !g.IO.KeysDown[user_key_index]
}

// uses provided repeat rate/delay. return a count, most often 0 or 1 but might be >1 if RepeatRate is small enough that DeltaTime > RepeatRate
func GetKeyPressedAmount(key_index int, repeat_delay float, repeat_rate float) int {
	var g = GImGui
	if key_index < 0 {
		return 0
	}
	IM_ASSERT(key_index >= 0 && key_index < int(len(g.IO.KeysDown)))
	var t = g.IO.KeysDownDuration[key_index]
	return CalcTypematicRepeatAmount(t-g.IO.DeltaTime, t, repeat_delay, repeat_rate)
}

// attention: misleading name! manually override io.WantCaptureKeyboard flag next frame (said flag is entirely left for your application to handle). e.g. force capture keyboard when your widget is being hovered. This is equivalent to setting "io.WantCaptureKeyboard = want_capture_keyboard_value"  {panic("not implemented")} after the next NewFrame() call.
func CaptureKeyboardFromApp(want_capture_keyboard_value bool /*= true*/) {
	if want_capture_keyboard_value {
		GImGui.WantCaptureKeyboardNextFrame = 1
	} else {
		GImGui.WantCaptureKeyboardNextFrame = 0
	}
}

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
