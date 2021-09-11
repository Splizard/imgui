package imgui

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
