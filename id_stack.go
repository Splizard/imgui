package imgui

// ID stack/scopes
// Read the FAQ (docs/FAQ.md or http://dearimgui.org/faq) for more details about how ID are handled in dear imgui.
// - Those questions are answered and impacted by understanding of the ID stack system:
//   - "Q: Why is my widget not reacting when I click on it?"
//   - "Q: How can I have widgets with an empty label?"
//   - "Q: How can I have multiple widgets with the same label?"
// - Short version: ID are hashes of the entire ID stack. If you are creating widgets in a loop you most likely
//   want to push a unique identifier (e.g. object pointer, loop index) to uniquely differentiate them.
// - You can also use the "Label##foobar" syntax within widget label to distinguish them from each others.
// - In this header file we use the "label"/"name" terminology to denote a string that will be displayed + used as an ID,
//   whereas "str_id" denote a string that is only used as an ID and not normally displayed.

// PushOverrideID Push given value as-is at the top of the ID stack (whereas PushID combines old and new hashes)
func PushOverrideID(id ImGuiID) {
	g := GImGui
	window := g.CurrentWindow
	window.IDStack = append(window.IDStack, id)
}

// GetIDWithSeed Helper to avoid a common series of PushOverrideID . GetID() . PopID() call
// (note that when using this pattern, TestEngine's "Stack Tool" will tend to not display the intermediate stack level.
//
//	for that to work we would need to do PushOverrideID() . ItemAdd() . PopID() which would alter widget code a little more)
func GetIDWithSeed(str string, seed ImGuiID) ImGuiID {
	var id = ImHashStr(str, uintptr(len(str)), seed)
	KeepAliveID(id)
	return id
}

func PushString(str_id string) {
	g := GImGui
	window := g.CurrentWindow
	id := window.GetIDNoKeepAlive(str_id)
	window.IDStack = append(window.IDStack, id)
}

// PushInterface push pointer into the ID stack (will hash pointer).
func PushInterface(ptr_id any) {
	g := GImGui
	window := g.CurrentWindow
	id := window.GetIDNoKeepAliveInterface(ptr_id)
	window.IDStack = append(window.IDStack, id)
}

// PushID push integer into the ID stack (will hash integer).
func PushID(int_id int) {
	g := GImGui
	window := g.CurrentWindow
	id := window.GetIDNoKeepAliveInt(int_id)
	window.IDStack = append(window.IDStack, id)
}

func PopID() {
	window := GImGui.CurrentWindow
	IM_ASSERT(len(window.IDStack) > 1) // Too many PopID(), or could be popping in a wrong/different window?
	window.IDStack = window.IDStack[:len(window.IDStack)-1]
} // pop from the ID stack.

func GetIDFromString(str_id string) ImGuiID {
	return GImGui.CurrentWindow.GetIDs(str_id)

} // calculate unique ID (hash of whole ID stack + given parameter). e.g. if you want to query into ImGuiStorage yourself

func GetIDs(str_id_begin string) ImGuiID {
	return GImGui.CurrentWindow.GetIDs(str_id_begin)
}

func GetIDFromInterface(ptr_id any) ImGuiID {
	window := GImGui.CurrentWindow
	return window.GetIDInterface(ptr_id)
}
