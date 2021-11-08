package imgui

import "fmt"

// Widgets: Text
// raw text without formatting. Roughly equivalent to Text("%s", text) but: A) doesn't require null terminated string if 'text_end' is specified, B) it's faster, no memory copy is done, no buffer size limits, recommended for long chunks of text.
func TextUnformatted(text string) {
	TextEx(text, ImGuiTextFlags_NoWidthForLargeClippedText)
}

// shortcut for PushStyleColor(ImGuiCol_Text, style.Colors[ImGuiCol_TextDisabled]); Text(fmt, ...); PopStyleColor()  {panic("not implemented")}
func TextDisabled(format string, args ...interface{}) {
	var g = GImGui
	PushStyleColorVec(ImGuiCol_Text, &g.Style.Colors[ImGuiCol_TextDisabled])
	if format[0] == '%' && format[1] == 's' && format[2] == 0 {
		TextEx(fmt.Sprintf(format, args...), ImGuiTextFlags_NoWidthForLargeClippedText) // Skip formatting
	} else {
		Text(format, args...)
	}
	PopStyleColor(1)
}

func TextWrapped(fmt string, args ...interface{})             { panic("not implemented") } // shortcut for PushTextWrapPos(0.0); Text(fmt, ...); PopTextWrapPos()  {panic("not implemented")}. Note that this won't work on an auto-resizing window if there's no other widgets to extend the window width, yoy may need to set a size using SetNextWindowSize().
func LabelText(label string, fmt string, args ...interface{}) { panic("not implemented") } // display text+label aligned the same way as value+label widgets
func BulletText(fmt string, args ...interface{})              { panic("not implemented") } // shortcut for Bullet()+Text()
