package imgui

import "fmt"

// Widgets: Text
// raw text without formatting. Roughly equivalent to Text("%s", text) but: A) doesn't require null terminated string if 'text_end' is specified, B) it's faster, no memory copy is done, no buffer size limits, recommended for long chunks of text.
func TextUnformatted(text string) {
	TextEx(text, ImGuiTextFlags_NoWidthForLargeClippedText)
}

// shortcut for PushStyleColor(ImGuiCol_Text, style.Colors[ImGuiCol_TextDisabled]); Text(fmt, ...); PopStyleColor()  {panic("not implemented")}
func TextDisabled(format string, args ...any) {
	g := GImGui
	PushStyleColorVec(ImGuiCol_Text, &g.Style.Colors[ImGuiCol_TextDisabled])
	if format[0] == '%' && format[1] == 's' && format[2] == 0 {
		TextEx(fmt.Sprintf(format, args...), ImGuiTextFlags_NoWidthForLargeClippedText) // Skip formatting
	} else {
		Text(format, args...)
	}
	PopStyleColor(1)
}

// shortcut for PushTextWrapPos(0.0); Text(fmt, ...); PopTextWrapPos()  {panic("not implemented")}. Note that this won't work on an auto-resizing window if there's no other widgets to extend the window width, yoy may need to set a size using SetNextWindowSize().
func TextWrapped(format string, args ...any) {
	g := GImGui
	var need_backup = (g.CurrentWindow.DC.TextWrapPos < 0.0) // Keep existing wrap position if one is already set
	if need_backup {
		PushTextWrapPos(0.0)
	}
	if format[0] == '%' && format[1] == 's' && format[2] == 0 {
		TextEx(fmt.Sprintf(format, args...), ImGuiTextFlags_NoWidthForLargeClippedText) // Skip formatting
	} else {
		Text(format, args...)
	}
	if need_backup {
		PopTextWrapPos()
	}
}

// display text+label aligned the same way as value+label widgets
func LabelText(label string, format string, args ...any) {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	g := GImGui
	style := g.Style
	var w = CalcItemWidth()

	var value = fmt.Sprintf(format, args...)
	var value_size = CalcTextSize(value, false, -1)
	var label_size = CalcTextSize(label, true, -1)

	var pos = window.DC.CursorPos
	var value_bb = ImRect{pos, pos.Add(ImVec2{w, value_size.y + style.FramePadding.y*2})}

	var padding float
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var total_bb = ImRect{pos, pos.Add(ImVec2{w + padding, max(value_size.y, label_size.y) + style.FramePadding.y*2})}
	ItemSizeRect(&total_bb, style.FramePadding.y)
	if !ItemAdd(&total_bb, 0, nil, 0) {
		return
	}

	// Render
	min := value_bb.Min.Add(style.FramePadding)
	RenderTextClipped(&min, &value_bb.Max, value, &value_size, &ImVec2{}, nil)
	if label_size.x > 0.0 {
		RenderText(ImVec2{value_bb.Max.x + style.ItemInnerSpacing.x, value_bb.Min.y + style.FramePadding.y}, label, true)
	}
}

// Text with a little bullet aligned to the typical tree node.
// shortcut for Bullet()+Text()
func BulletText(format string, args ...any) {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	g := GImGui
	style := g.Style

	var text = fmt.Sprintf(format, args...)
	var label_size = CalcTextSize(text, false, -1)

	var padding float
	if label_size.x > 0.0 { // Empty text doesn't add padding
		padding = (label_size.x + style.FramePadding.x*2)
	}

	var total_size = ImVec2{g.FontSize + padding, label_size.y}
	var pos = window.DC.CursorPos
	pos.y += window.DC.CurrLineTextBaseOffset
	ItemSizeVec(&total_size, 0.0)
	var bb = ImRect{pos, pos.Add(total_size)}
	if !ItemAdd(&bb, 0, nil, 0) {
		return
	}

	// Render
	var text_col = GetColorU32FromID(ImGuiCol_Text, 1)
	RenderBullet(window.DrawList, bb.Min.Add(ImVec2{style.FramePadding.x + g.FontSize*0.5, g.FontSize * 0.5}), text_col)
	RenderText(bb.Min.Add(ImVec2{g.FontSize + style.FramePadding.x*2, 0.0}), text, false)
}

// draw a small circle + keep the cursor on the same line. advance cursor x position by GetTreeNodeToLabelSpacing(), same distance that TreeNode() uses
func Bullet() {
	window := GetCurrentWindow()
	if window.SkipItems {
		return
	}

	g := GImGui
	style := g.Style
	var line_height = max(min(window.DC.CurrLineSize.y, g.FontSize+g.Style.FramePadding.y*2), g.FontSize)
	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(ImVec2{g.FontSize, line_height})}
	ItemSizeRect(&bb, 0)
	if !ItemAdd(&bb, 0, nil, 0) {
		SameLine(0, style.FramePadding.x*2)
		return
	}

	// Render and stay on same line
	var text_col = GetColorU32FromID(ImGuiCol_Text, 1)
	RenderBullet(window.DrawList, bb.Min.Add(ImVec2{style.FramePadding.x + g.FontSize*0.5, line_height * 0.5}), text_col)
	SameLine(0, style.FramePadding.x*2.0)
}
