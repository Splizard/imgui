package imgui

// Clipboard Utilities
// - Also see the LogToClipboard() function to capture GUI into clipboard, or easily output text data to the clipboard.

func GetClipboardText() string {
	var g = GImGui
	if g.IO.GetClipboardTextFn != nil {
		return g.IO.GetClipboardTextFn(g.IO.ClipboardUserData)
	}
	return ""
}

func SetClipboardText(text string) {
	var g = GImGui
	if g.IO.SetClipboardTextFn != nil {
		g.IO.SetClipboardTextFn(g.IO.ClipboardUserData, text)
	}
}

// GetClipboardTextFn_DefaultImpl Local Dear ImGui-only clipboard implementation, if user hasn't defined better clipboard handlers.
func GetClipboardTextFn_DefaultImpl(any) string {
	var g = GImGui
	if len(g.ClipboardHandlerData) == 0 {
		return ""
	}
	return string(g.ClipboardHandlerData)
}

func SetClipboardTextFn_DefaultImpl(_ any, text string) {
	var g = GImGui
	g.ClipboardHandlerData = g.ClipboardHandlerData[:0]
	g.ClipboardHandlerData = []byte(text)
}
