package imgui

// Clipboard Utilities
// - Also see the LogToClipboard() function to capture GUI into clipboard, or easily output text data to the clipboard.

func GetClipboardText() string {
	if guiContext.IO.GetClipboardTextFn != nil {
		return guiContext.IO.GetClipboardTextFn(guiContext.IO.ClipboardUserData)
	}
	return ""
}

func SetClipboardText(text string) {
	if guiContext.IO.SetClipboardTextFn != nil {
		guiContext.IO.SetClipboardTextFn(guiContext.IO.ClipboardUserData, text)
	}
}

// GetClipboardTextFn_DefaultImpl Local Dear ImGui-only clipboard implementation, if user hasn't defined better clipboard handlers.
func GetClipboardTextFn_DefaultImpl(any) string {
	if len(guiContext.ClipboardHandlerData) == 0 {
		return ""
	}
	return string(guiContext.ClipboardHandlerData)
}

func SetClipboardTextFn_DefaultImpl(_ any, text string) {
	guiContext.ClipboardHandlerData = guiContext.ClipboardHandlerData[:0]
	guiContext.ClipboardHandlerData = []byte(text)
}
