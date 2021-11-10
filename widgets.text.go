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

// Text with a little bullet aligned to the typical tree node.
// shortcut for Bullet()+Text()
func BulletText(format string, args ...interface{}) {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return
	}

	var g = GImGui
	var style = g.Style

	var text = fmt.Sprintf(format, args...)
	var label_size = CalcTextSize(text, false, -1)

	var padding float
	if label_size.x > 0.0 { // Empty text doesn't add padding
		padding = (label_size.x + style.FramePadding.x*2)
	}

	var total_size = ImVec2{g.FontSize + padding, label_size.y}
	var pos ImVec2 = window.DC.CursorPos
	pos.y += window.DC.CurrLineTextBaseOffset
	ItemSizeVec(&total_size, 0.0)
	var bb = ImRect{pos, pos.Add(total_size)}
	if !ItemAdd(&bb, 0, nil, 0) {
		return
	}

	// Render
	var text_col ImU32 = GetColorU32FromID(ImGuiCol_Text, 1)
	RenderBullet(window.DrawList, bb.Min.Add(ImVec2{style.FramePadding.x + g.FontSize*0.5, g.FontSize * 0.5}), text_col)
	RenderText(bb.Min.Add(ImVec2{g.FontSize + style.FramePadding.x*2, 0.0}), text, false)
}
