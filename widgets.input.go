package imgui

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Create text input in place of another active widget (e.g. used when doing a CTRL+Click on drag/slider widgets)
// FIXME: Facilitate using this in variety of other situations.
func TempInputText(bb *ImRect, id ImGuiID, label string, buf *[]byte, flags ImGuiInputTextFlags) bool {
	// On the first frame, g.TempInputTextId == 0, then on subsequent frames it becomes == id.
	// We clear ActiveID on the first frame to allow the InputText() taking it back.
	var g = GImGui
	var init = (g.TempInputId != id)
	if init {
		ClearActiveID()
	}

	g.CurrentWindow.DC.CursorPos = bb.Min
	size := bb.GetSize()
	var value_changed = InputTextEx(label, "", buf, &size, flags|ImGuiInputTextFlags_MergedItem, nil, nil)
	if init {
		// First frame we started displaying the InputText widget, we expect it to take the active id.
		IM_ASSERT(g.ActiveId == id)
		g.TempInputId = g.ActiveId
	}
	return value_changed
}

// Note that Drag/Slider functions are only forwarding the min/max values clamping values if the ImGuiSliderFlags_AlwaysClamp flag is set!
// This is intended: this way we allow CTRL+Click manual input to set a value out of bounds, for maximum flexibility.
// However this may not be ideal for all uses, as some user code may break on out of bound values.
func TempInputScalar(bb *ImRect, id ImGuiID, label string, data_type ImGuiDataType, p_data interface{}, format string, p_clamp_min interface{}, p_clamp_max interface{}) bool {
	var g = GImGui
	var data_buf = []byte(strings.TrimSpace(fmt.Sprintf(format, p_data)))

	var flags ImGuiInputTextFlags = ImGuiInputTextFlags_AutoSelectAll | ImGuiInputTextFlags_NoMarkEdited
	flags |= ImGuiInputTextFlags_CharsDecimal
	if data_type == ImGuiDataType_Float || data_type == ImGuiDataType_Double {
		flags = ImGuiInputTextFlags_CharsScientific
	}

	var value_changed = false
	if TempInputText(bb, id, label, &data_buf, flags) {
		// Backup old value
		var data_backup = reflect.ValueOf(p_data).Elem().Interface()

		// Apply new value (or operations) then clamp
		DataTypeApplyOpFromText(string(data_buf), string(g.InputTextState.InitialTextA), data_type, p_data, "")
		if p_clamp_min != nil || p_clamp_max != nil {
			if p_clamp_min != nil && p_clamp_max != nil && DataTypeCompare(data_type, p_clamp_min, p_clamp_max) > 0 {
				p_clamp_min, p_clamp_max = p_clamp_max, p_clamp_min
			}
			DataTypeClamp(data_type, p_data, p_clamp_min, p_clamp_max)
		}

		// Only mark as edited if new value is different
		value_changed = reflect.ValueOf(p_data).Elem().Interface() != data_backup
		if value_changed {
			MarkItemEdited(id)
		}
	}
	return value_changed
}

func TempInputIsActive(id ImGuiID) bool {
	var g *ImGuiContext = GImGui
	return (g.ActiveId == id && g.TempInputId == id)
}

func GetInputTextState(id ImGuiID) *ImGuiInputTextState {
	var g *ImGuiContext = GImGui
	if g.InputTextState.ID == id {
		return &g.InputTextState
	}
	return nil
} // Get input text state if active

func InputText(label string, char *[]byte, flags ImGuiInputTextFlags, callback ImGuiInputTextCallback /*= L*/, user_data interface{}) bool {
	IM_ASSERT((flags & ImGuiInputTextFlags_Multiline) == 0) // call InputTextMultiline()
	return InputTextEx(label, "", char, &ImVec2{}, flags, callback, user_data)
}

func InputTextMultiline(label string, buf *[]byte, size ImVec2, flags ImGuiInputTextFlags, callback ImGuiInputTextCallback /*= L*/, user_data interface{}) bool {
	return InputTextEx(label, "", buf, &size, flags|ImGuiInputTextFlags_Multiline, callback, user_data)
}

func InputTextWithHint(label string, hint string, char *[]byte, flags ImGuiInputTextFlags, callback ImGuiInputTextCallback /*= L*/, user_data interface{}) bool {
	IM_ASSERT((flags & ImGuiInputTextFlags_Multiline) == 0) // call InputTextMultiline()
	return InputTextEx(label, hint, char, &ImVec2{}, flags, callback, user_data)
}

func InputFloat(label string, v *float, step, step_fast float, format string, flags ImGuiInputTextFlags) bool {
	flags |= ImGuiInputTextFlags_CharsScientific

	var step_arg *float
	if step > 0.0 {
		step_arg = &step
	}

	var step_fast_arg *float
	if step_fast > 0.0 {
		step_fast_arg = &step_fast
	}

	return InputScalarFloat32(label, v, step_arg, step_fast_arg, format, flags)
}

func InputFloat2(label string, v *[2]float, format string, flags ImGuiInputTextFlags) bool {
	return InputScalarFloat32s(label, (*v)[:], nil, nil, format, flags)
}

func InputFloat3(label string, v *[3]float, format string /*= "%.3f"*/, flags ImGuiInputTextFlags) bool {
	return InputScalarFloat32s(label, (*v)[:], nil, nil, format, flags)
}

func InputFloat4(label string, v *[4]float, format string /*= "%.3f"*/, flags ImGuiInputTextFlags) bool {
	return InputScalarFloat32s(label, (*v)[:], nil, nil, format, flags)
}

func InputInt(label string, v *int, step int /*= 1*/, step_fast int /*= 100*/, flags ImGuiInputTextFlags) bool {
	// Hexadecimal input provided as a convenience but the flag name is awkward. Typically you'd use InputText() to parse your own data, if you want to handle prefixes.
	var format = "%v"
	if (flags & ImGuiInputTextFlags_CharsHexadecimal) != 0 {
		format = "%08X"
	}

	var step_arg *int
	if step > 0.0 {
		step_arg = &step
	}

	var step_fast_arg *int
	if step_fast > 0.0 {
		step_fast_arg = &step_fast
	}

	return InputScalarInt32(label, v, step_arg, step_fast_arg, format, flags)
}

func InputInt2(label string, v *[2]int, flags ImGuiInputTextFlags) bool {
	return InputScalarInt32s(label, v[:], nil, nil, "%d", flags)
}

func InputInt3(label string, v *[3]int, flags ImGuiInputTextFlags) bool {
	return InputScalarInt32s(label, v[:], nil, nil, "%d", flags)
}

func InputInt4(label string, v *[4]int, flags ImGuiInputTextFlags) bool {
	return InputScalarInt32s(label, v[:], nil, nil, "%d", flags)
}

func InputDouble(label string, v *double, step double /*= 0*/, step_fast double /*= 0*/, format string /*= "%.6f"*/, flags ImGuiInputTextFlags) bool {
	flags |= ImGuiInputTextFlags_CharsScientific

	var step_arg *double
	if step > 0.0 {
		step_arg = &step
	}

	var step_fast_arg *double
	if step_fast > 0.0 {
		step_fast_arg = &step_fast
	}

	return InputScalarFloat64(label, v, step_arg, step_fast_arg, format, flags)
}

func InputTextCalcLineCount(text string) int {
	return int(strings.Count(text, "\n"))
}

func InputTextCalcTextSizeW(text []ImWchar, remaining *[]ImWchar, out_offset *ImVec2, stop_on_new_line bool) ImVec2 {
	var g = GImGui
	var font = g.Font
	var line_height float = g.FontSize
	var scale float = line_height / font.FontSize

	var text_size ImVec2
	var line_width float

	var s = text
	for len(s) > 0 {
		var c = s[0]
		s = s[1:]
		if c == '\n' {
			text_size.x = ImMax(text_size.x, line_width)
			text_size.y += line_height
			line_width = 0.0
			if stop_on_new_line {
				break
			}
			continue
		}
		if c == '\r' {
			continue
		}

		var char_width float = font.GetCharAdvance(rune(c)) * scale
		line_width += char_width
	}

	if text_size.x < line_width {
		text_size.x = line_width
	}

	if out_offset != nil {
		*out_offset = ImVec2{line_width, text_size.y + line_height} // offset allow for the possibility of sitting after a trailing \n
	}

	if line_width > 0 || text_size.y == 0.0 { // whereas size.y will ignore the trailing \n
		text_size.y += line_height
	}

	if remaining != nil {
		*remaining = s
	}

	return text_size
}

func resize(buf *[]byte, new_size int) {
	if new_size > int(len(*buf)) {
		var new_buf = make([]byte, new_size)
		copy(new_buf, *buf)
		*buf = new_buf
	} else {
		*buf = (*buf)[:new_size]
	}
}

func resizeRune(buf *[]rune, new_size int) {
	if new_size > int(len(*buf)) {
		var new_buf = make([]rune, new_size)
		copy(new_buf, *buf)
		*buf = new_buf
	} else {
		*buf = (*buf)[:new_size]
	}
}

func runesContains(runes []rune, c rune) bool {
	for _, r := range runes {
		if r == c {
			return true
		}
	}
	return false
}

// Return false to discard a character.
func InputTextFilterCharacter(p_char *rune, flags ImGuiInputTextFlags, callback ImGuiInputTextCallback, user_data interface{}, input_source ImGuiInputSource) bool {
	IM_ASSERT(input_source == ImGuiInputSource_Keyboard || input_source == ImGuiInputSource_Clipboard)
	var c = *p_char

	// Filter non-printable (NB: isprint is unreliable! see #2467)
	var apply_named_filters = true
	if c < 0x20 {
		var pass = false
		pass = pass || (c == '\n' && (flags&ImGuiInputTextFlags_Multiline != 0))
		pass = pass || (c == '\t' && (flags&ImGuiInputTextFlags_AllowTabInput != 0))
		if !pass {
			return false
		}
		apply_named_filters = false // Override named filters below so newline and tabs can still be inserted.
	}

	if input_source != ImGuiInputSource_Clipboard {
		// We ignore Ascii representation of delete (emitted from Backspace on OSX, see #2578, #2817)
		if c == 127 {
			return false
		}

		// Filter private Unicode range. GLFW on OSX seems to send private characters for special keys like arrow keys (FIXME)
		if c >= 0xE000 && c <= 0xF8FF {
			return false
		}
	}

	// Filter Unicode ranges we are not handling in this build
	if c > IM_UNICODE_CODEPOINT_MAX {
		return false
	}

	// Generic named filters
	if apply_named_filters && (flags&(ImGuiInputTextFlags_CharsDecimal|ImGuiInputTextFlags_CharsHexadecimal|ImGuiInputTextFlags_CharsUppercase|ImGuiInputTextFlags_CharsNoBlank|ImGuiInputTextFlags_CharsScientific) != 0) {
		// The libc allows overriding locale, with e.g. 'setlocale(LC_NUMERIC, "de_DE.UTF-8");' which affect the output/input of printf/scanf.
		// The standard mandate that programs starts in the "C" locale where the decimal point is '.'.
		// We don't really intend to provide widespread support for it, but out of empathy for people stuck with using odd API, we support the bare minimum aka overriding the decimal point.
		// Change the default decimal_point with:
		//   ImGui::GetCurrentContext()->PlatformLocaleDecimalPoint = *localeconv()->decimal_point;
		var g = GImGui
		var c_decimal_point rune = (rune)(g.PlatformLocaleDecimalPoint)

		// Allow 0-9 . - + * /
		if flags&ImGuiInputTextFlags_CharsDecimal != 0 {
			if !(c >= '0' && c <= '9') && (c != c_decimal_point) && (c != '-') && (c != '+') && (c != '*') && (c != '/') {
				return false
			}
		}

		// Allow 0-9 . - + * / e E
		if flags&ImGuiInputTextFlags_CharsScientific != 0 {
			if !(c >= '0' && c <= '9') && (c != c_decimal_point) && (c != '-') && (c != '+') && (c != '*') && (c != '/') && (c != 'e') && (c != 'E') {
				return false
			}
		}

		// Allow 0-9 a-F A-F
		if flags&ImGuiInputTextFlags_CharsHexadecimal != 0 {
			if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') && !(c >= 'A' && c <= 'F') {
				return false
			}
		}

		// Turn a-z into A-Z
		if flags&ImGuiInputTextFlags_CharsUppercase != 0 {
			if c >= 'a' && c <= 'z' {
				c += (rune)('A' - 'a')
				*p_char = (c)
			}
		}

		if flags&ImGuiInputTextFlags_CharsNoBlank != 0 {
			if ImCharIsBlankW(c) {
				return false
			}
		}
	}

	// Custom callback filter
	if flags&ImGuiInputTextFlags_CallbackCharFilter != 0 {
		var callback_data ImGuiInputTextCallbackData
		callback_data.EventFlag = ImGuiInputTextFlags_CallbackCharFilter
		callback_data.EventChar = (ImWchar)(c)
		callback_data.Flags = flags
		callback_data.UserData = user_data
		if callback(&callback_data) != 0 {
			return false
		}
		*p_char = callback_data.EventChar
		if callback_data.EventChar == 0 {
			return false
		}
	}

	return true
}

func ImStrbolW(buf_mid_line []ImWchar, buf_begin []ImWchar) []ImWchar { // find beginning-of-line
	// FIXME: this is probably wrong
	/*
		while (buf_mid_line > buf_begin && buf_mid_line[-1] != '\n')
		    buf_mid_line--;
		return buf_mid_line;
	*/

	var i int
	for i = int(len(buf_mid_line) - 1); i > 0 && buf_mid_line[i] != '\n'; i-- {
		// noop
	}

	return []ImWchar{ImWchar(i)}
}

/*

	ABSOLUTE MONSTER OF A FUNCTION AHEAD
    VERY DIFFICULT TO PORT AND LIKELY DOESN'T WORK
	IN THE CURRENT STATE

*/

// Edit a string of text
//   - buf_size account for the zero-terminator, so a buf_size of 6 can hold "Hello" but not "Hello!".
//     This is so we can easily call InputText() on static arrays using ARRAYSIZE() and to match
//     Note that in std::string world, capacity() would omit 1 byte used by the zero-terminator.
//   - When active, hold on a privately held copy of the text (and apply back to 'buf'). So changing 'buf' while the InputText is active has no effect.
//   - If you want to use ImGui::InputText() with std::string, see misc/cpp/imgui_stdlib.h
//
// (FIXME: Rather confusing and messy function, among the worse part of our codebase, expecting to rewrite a V2 at some point.. Partly because we are
//
//	doing UTF8 > U16 > UTF8 conversions on the go to easily interface with stb_textedit. Ideally should stay in UTF-8 all the time. See https://github.com/nothings/stb/issues/188)
func InputTextEx(label string, hint string, buf *[]byte, size_arg *ImVec2, flags ImGuiInputTextFlags, callback ImGuiInputTextCallback, callback_user_data interface{}) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	IM_ASSERT(buf != nil && int(len(*buf)) >= 0)
	// TODO: check these asserts
	IM_ASSERT(!((flags&ImGuiInputTextFlags_CallbackHistory) == 0 && (flags&ImGuiInputTextFlags_Multiline != 0)))        // Can't use both together (they both use up/down keys)
	IM_ASSERT(!((flags&ImGuiInputTextFlags_CallbackCompletion) == 0 && (flags&ImGuiInputTextFlags_AllowTabInput != 0))) // Can't use both together (they both use tab key)

	var g = GImGui
	var io = g.IO
	var style = g.Style

	var RENDER_SELECTION_WHEN_INACTIVE = false
	var is_multiline = (flags & ImGuiInputTextFlags_Multiline) != 0
	var is_readonly = (flags & ImGuiInputTextFlags_ReadOnly) != 0
	var is_password = (flags & ImGuiInputTextFlags_Password) != 0
	var is_undoable = (flags & ImGuiInputTextFlags_NoUndoRedo) == 0
	var is_resizable = (flags & ImGuiInputTextFlags_CallbackResize) != 0
	if is_resizable {
		IM_ASSERT(callback != nil) // Must provide a callback if you set the ImGuiInputTextFlags_CallbackResize flag!
	}

	if is_multiline { // Open group before calling GetID() because groups tracks id created within their scope,
		BeginGroup()
	}

	var id = window.GetIDs(label)
	var label_size = CalcTextSize(label, true, -1)

	var padding float
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	fontsize := label_size.y
	if is_multiline {
		fontsize = g.FontSize * 8.0
	}

	var frame_size = CalcItemSize(*size_arg, CalcItemWidth(), fontsize+style.FramePadding.y*2.0) // Arbitrary default of 8 lines high for multi-line
	var total_size = ImVec2{frame_size.x + padding, frame_size.y}

	var frame_bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(frame_size)}
	var total_bb = ImRect{frame_bb.Min, frame_bb.Min.Add(total_size)}

	var draw_window = window
	var inner_size = frame_size
	var item_status_flags ImGuiItemStatusFlags = 0
	if is_multiline {
		var backup_pos ImVec2 = window.DC.CursorPos
		ItemSizeRect(&total_bb, style.FramePadding.y)
		if !ItemAdd(&total_bb, id, &frame_bb, ImGuiItemFlags_Inputable) {
			EndGroup()
			return false
		}
		item_status_flags = g.LastItemData.StatusFlags
		window.DC.CursorPos = backup_pos

		// We reproduce the contents of BeginChildFrame() in order to provide 'label' so our window internal data are easier to read/debug.
		PushStyleColorVec(ImGuiCol_ChildBg, &style.Colors[ImGuiCol_FrameBg])
		PushStyleFloat(ImGuiStyleVar_ChildRounding, style.FrameRounding)
		PushStyleFloat(ImGuiStyleVar_ChildBorderSize, style.FrameBorderSize)

		size := frame_bb.GetSize()

		var child_visible = BeginChildEx(label, id, &size, true, ImGuiWindowFlags_NoMove)
		PopStyleVar(2)
		PopStyleColor(1)
		if !child_visible {
			EndChild()
			EndGroup()
			return false
		}
		draw_window = g.CurrentWindow                                                   // Child window
		draw_window.DC.NavLayersActiveMaskNext |= (1 << draw_window.DC.NavLayerCurrent) // This is to ensure that EndChild() will display a navigation highlight so we can "enter" into it.
		draw_window.DC.CursorPos = draw_window.DC.CursorPos.Add(style.FramePadding)
		inner_size.x -= draw_window.ScrollbarSizes.x
	} else {
		// Support for internal ImGuiInputTextFlags_MergedItem flag, which could be redesigned as an ItemFlags if needed (with test performed in ItemAdd)
		ItemSizeRect(&total_bb, style.FramePadding.y)
		if (flags & ImGuiInputTextFlags_MergedItem) == 0 {
			if !ItemAdd(&total_bb, id, &frame_bb, ImGuiItemFlags_Inputable) {
				return false
			}
		}
		item_status_flags = g.LastItemData.StatusFlags
	}
	var hovered = ItemHoverable(&frame_bb, id)
	if hovered {
		g.MouseCursor = ImGuiMouseCursor_TextInput
	}

	// We are only allowed to access the state if we are already the active widget.
	var state *ImGuiInputTextState = GetInputTextState(id)

	var focus_requested_by_code = (item_status_flags & ImGuiItemStatusFlags_FocusedByCode) != 0
	var focus_requested_by_tabbing = (item_status_flags & ImGuiItemStatusFlags_FocusedByTabbing) != 0

	var user_clicked = hovered && io.MouseClicked[0]
	var user_nav_input_start = (g.ActiveId != id) && ((g.NavInputId == id) || (g.NavActivateId == id && g.NavInputSource == ImGuiInputSource_Keyboard))
	var user_scroll_finish = is_multiline && state != nil && g.ActiveId == 0 && g.ActiveIdPreviousFrame == GetWindowScrollbarID(draw_window, ImGuiAxis_Y)
	var user_scroll_active = is_multiline && state != nil && g.ActiveId == GetWindowScrollbarID(draw_window, ImGuiAxis_Y)

	var clear_active_id = false
	var select_all = (g.ActiveId != id) && ((flags&ImGuiInputTextFlags_AutoSelectAll) != 0 || user_nav_input_start) && (!is_multiline)

	var scroll_y float = FLT_MAX
	if is_multiline {
		scroll_y = draw_window.Scroll.y
	}

	var init_changed_specs = (state != nil && (state.Stb.single_line != 0) != !is_multiline)
	var init_make_active = (user_clicked || user_scroll_finish || user_nav_input_start || focus_requested_by_code || focus_requested_by_tabbing)
	var init_state = (init_make_active || user_scroll_active)
	if (init_state && g.ActiveId != id) || init_changed_specs {
		// Access state even if we don't own it yet.
		state = &g.InputTextState
		state.CursorAnimReset()

		// Take a copy of the initial buffer value (both in original UTF-8 format and converted to wchar)
		// From the moment we focused we are ignoring the content of 'buf' (unless we are in read-only mode)
		var buf_len = (int)(len(*buf))

		// UTF-8. we use +1 to make sure that .Data is always pointing to at least an empty string.
		resize(&state.InitialTextA, buf_len+1)
		copy(state.InitialTextA, *buf)

		// Start edition
		var buf_end string
		resizeRune(&state.TextW, int(len(*buf)+1))
		resizeRune(&state.TextW, int(len(*buf)+1)) // wchar count <= UTF-8 count. we use +1 to make sure that .Data is always pointing to at least an empty string.
		state.TextA = state.TextA[:0]
		state.TextAIsValid = false // TextA is not valid yet (we will display buf until then)
		state.CurLenW = ImTextStrFromUtf8(state.TextW, int(len(*buf)), string(*buf), &buf_end)
		state.CurLenA = (int)(len(*buf) - len(buf_end)) // We can't get the result from ImStrncpy() above because it is not UTF-8 aware. Here we'll cut off malformed UTF-8.

		// Preserve cursor position and undo/redo stack if we come back to same widget
		// FIXME: For non-readonly widgets we might be able to require that TextAIsValid && TextA == buf ? (untested) and discard undo stack if user buffer has changed.
		var recycle_state = (state.ID == id && !init_changed_specs)
		if recycle_state {
			// Recycle existing cursor/selection/undo stack but clamp position
			// Note a single mouse click will override the cursor/position immediately by calling stb_textedit_click handler.
			state.CursorClamp()
		} else {
			state.ID = id
			state.ScrollX = 0.0
			stb_textedit_initialize_state(&state.Stb, bool2int(!is_multiline))
			if !is_multiline && focus_requested_by_code {
				select_all = true
			}
		}
		if (flags & ImGuiInputTextFlags_AlwaysOverwrite) != 0 {
			state.Stb.insert_mode = 1 // stb field name is indeed incorrect (see #2863)
		}
		if !is_multiline && (focus_requested_by_tabbing || (user_clicked && io.KeyCtrl)) {
			select_all = true
		}
	}

	if g.ActiveId != id && init_make_active {
		IM_ASSERT(state != nil && state.ID == id)
		SetActiveID(id, window)
		SetFocusID(id, window)
		FocusWindow(window)

		// Declare our inputs
		IM_ASSERT(ImGuiNavInput_COUNT < 32)
		g.ActiveIdUsingNavDirMask |= (1 << ImGuiDir_Left) | (1 << ImGuiDir_Right)
		if is_multiline || (flags&ImGuiInputTextFlags_CallbackHistory != 0) {
			g.ActiveIdUsingNavDirMask |= (1 << ImGuiDir_Up) | (1 << ImGuiDir_Down)
		}
		g.ActiveIdUsingNavInputMask |= (1 << ImGuiNavInput_Cancel)
		g.ActiveIdUsingKeyInputMask |= ((ImU64)(1 << ImGuiKey_Home)) | ((ImU64)(1 << ImGuiKey_End))
		if is_multiline {
			g.ActiveIdUsingKeyInputMask |= ((ImU64)(1 << ImGuiKey_PageUp)) | ((ImU64)(1 << ImGuiKey_PageDown))
		}
		if flags&(ImGuiInputTextFlags_CallbackCompletion|ImGuiInputTextFlags_AllowTabInput) != 0 { // Disable keyboard tabbing out as we will use the \t character.
			g.ActiveIdUsingKeyInputMask |= ((ImU64)(1 << ImGuiKey_Tab))
		}
	}

	// We have an edge case if ActiveId was set through another widget (e.g. widget being swapped), clear id immediately (don't wait until the end of the function)
	if g.ActiveId == id && state == nil {
		ClearActiveID()
	}

	// Release focus when we click outside
	if g.ActiveId == id && io.MouseClicked[0] && !init_state && !init_make_active { //-V560
		clear_active_id = true
	}

	// Lock the decision of whether we are going to take the path displaying the cursor or selection
	var render_cursor = (g.ActiveId == id) || (state != nil && user_scroll_active)
	var render_selection = state != nil && state.HasSelection() && (RENDER_SELECTION_WHEN_INACTIVE || render_cursor)
	var value_changed = false
	var enter_pressed = false

	// When read-only we always use the live data passed to the function
	// FIXME-OPT: Because our selection/cursor code currently needs the wide text we need to convert it when active, which is not ideal :(
	if is_readonly && state != nil && (render_cursor || render_selection) {
		var buf_end string
		resizeRune(&state.TextW, int(len(*buf))+1)
		state.CurLenW = ImTextStrFromUtf8(state.TextW, int(len(state.TextW)), string(*buf), &buf_end)
		state.CurLenA = (int)(len(*buf) - len(buf_end))
		state.CursorClamp()
		render_selection = render_selection && state.HasSelection()
	}

	// Select the buffer to render.
	var buf_display_from_state = (render_cursor || render_selection || g.ActiveId == id) && !is_readonly && state != nil && state.TextAIsValid

	b := *buf
	if buf_display_from_state {
		b = state.TextA
	}

	var is_displaying_hint = (hint != "" && b[0] == 0)

	// Password pushes a temporary font with only a fallback glyph
	if is_password && !is_displaying_hint {
		var glyph = g.Font.FindGlyph('*')
		var password_font = &g.InputTextPasswordFont
		password_font.FontSize = g.Font.FontSize
		password_font.Scale = g.Font.Scale
		password_font.Ascent = g.Font.Ascent
		password_font.Descent = g.Font.Descent
		password_font.ContainerAtlas = g.Font.ContainerAtlas
		password_font.FallbackGlyph = glyph
		password_font.FallbackAdvanceX = glyph.AdvanceX
		IM_ASSERT(len(password_font.Glyphs) == 0 && len(password_font.IndexAdvanceX) == 0 && len(password_font.IndexLookup) == 0)
		PushFont(password_font)
	}

	// Process mouse inputs and character inputs
	var backup_current_text_length int = 0
	if g.ActiveId == id {
		IM_ASSERT(state != nil)
		backup_current_text_length = state.CurLenA
		state.Edited = false
		state.BufCapacityA = int(len(*buf))
		state.Flags = flags
		state.UserCallback = callback
		state.UserCallbackData = callback_user_data

		// Although we are active we don't prevent mouse from hovering other elements unless we are interacting right now with the widget.
		// Down the line we should have a cleaner library-wide concept of Selected vs Active.
		g.ActiveIdAllowOverlap = !io.MouseDown[0]
		g.WantTextInputNextFrame = 1

		// Edit in progress
		var mouse_x = (io.MousePos.x - frame_bb.Min.x - style.FramePadding.x) + state.ScrollX

		var mouse_y = g.FontSize * 0.5
		if is_multiline {
			mouse_y = io.MousePos.y - draw_window.DC.CursorPos.y
		}

		var is_osx = io.ConfigMacOSXBehaviors
		if select_all || (hovered && !is_osx && io.MouseDoubleClicked[0]) {
			state.SelectAll()
			state.SelectedAllMouseLock = true
		} else if hovered && is_osx && io.MouseDoubleClicked[0] {
			// Double-click select a word only, OS X style (by simulating keystrokes)
			state.OnKeyPressed(STB_TEXTEDIT_K_WORDLEFT)
			state.OnKeyPressed(STB_TEXTEDIT_K_WORDRIGHT | STB_TEXTEDIT_K_SHIFT)
		} else if io.MouseClicked[0] && !state.SelectedAllMouseLock {
			if hovered {
				stb_textedit_click(state, &state.Stb, mouse_x, mouse_y)
				state.CursorAnimReset()
			}
		} else if io.MouseDown[0] && !state.SelectedAllMouseLock && (io.MouseDelta.x != 0.0 || io.MouseDelta.y != 0.0) {
			stb_textedit_drag(state, &state.Stb, mouse_x, mouse_y)
			state.CursorAnimReset()
			state.CursorFollow = true
		}
		if state.SelectedAllMouseLock && !io.MouseDown[0] {
			state.SelectedAllMouseLock = false
		}

		// It is ill-defined whether the backend needs to send a \t character when pressing the TAB keys.
		// Win32 and GLFW naturally do it but not SDL.
		var ignore_char_inputs = (io.KeyCtrl && !io.KeyAlt) || (is_osx && io.KeySuper)
		if (flags&ImGuiInputTextFlags_AllowTabInput != 0) && IsKeyPressedMap(ImGuiKey_Tab, true) && !ignore_char_inputs && !io.KeyShift && !is_readonly {
			if !runesContains(io.InputQueueCharacters, '\t') {
				var c rune = '\t' // Insert TAB
				if InputTextFilterCharacter(&c, flags, callback, callback_user_data, ImGuiInputSource_Keyboard) {
					state.OnKeyPressed((int)(c)) // TODO: check this
				}
			}
		}

		// Process regular text input (before we check for Return because using some IME will effectively send a Return?)
		// We ignore CTRL inputs, but need to allow ALT+CTRL as some keyboards (e.g. German) use AltGR (which _is_ Alt+Ctrl) to input certain characters.
		if len(io.InputQueueCharacters) > 0 {
			if !ignore_char_inputs && !is_readonly && !user_nav_input_start {
				for n := range io.InputQueueCharacters {
					// Insert character if they pass filtering
					var c = (rune)(io.InputQueueCharacters[n])
					if c == '\t' && io.KeyShift {
						continue
					}
					if InputTextFilterCharacter(&c, flags, callback, callback_user_data, ImGuiInputSource_Keyboard) {
						state.OnKeyPressed((int)(c))
					}
				}
			}

			// Consume characters
			io.InputQueueCharacters = io.InputQueueCharacters[:0]
		}
	}

	// Process other shortcuts/key-presses
	var cancel_edit = false
	if g.ActiveId == id && !g.ActiveIdIsJustActivated && !clear_active_id {
		IM_ASSERT(state != nil)
		IM_ASSERT_USER_ERROR(io.KeyMods == GetMergedKeyModFlags(), "Mismatching io.KeyCtrl/io.KeyShift/io.KeyAlt/io.KeySuper vs io.KeyMods") // We rarely do this check, but if anything let's do it here.

		var row_count_per_page = ImMaxInt((int)((inner_size.y-style.FramePadding.y)/g.FontSize), 1)
		state.Stb.row_count_per_page = row_count_per_page

		var k_mask int
		if io.KeyShift {
			k_mask = STB_TEXTEDIT_K_SHIFT
		}
		var is_osx = io.ConfigMacOSXBehaviors
		var is_osx_shift_shortcut = is_osx && (io.KeyMods == (ImGuiKeyModFlags_Super | ImGuiKeyModFlags_Shift))
		var is_wordmove_key_down = io.KeyCtrl
		var is_startend_key_down = is_osx && io.KeySuper && !io.KeyCtrl && !io.KeyAlt // OS X style: Line/Text Start and End using Cmd+Arrows instead of Home/End

		if is_osx {
			is_wordmove_key_down = io.KeyAlt // OS X style: Text editing cursor movement using Alt instead of Ctrl
		}

		var is_ctrl_key_only = (io.KeyMods == ImGuiKeyModFlags_Ctrl)
		var is_shift_key_only = (io.KeyMods == ImGuiKeyModFlags_Shift)
		var is_shortcut_key = (io.KeyMods == ImGuiKeyModFlags_Ctrl)

		if g.IO.ConfigMacOSXBehaviors {
			is_shortcut_key = (io.KeyMods == ImGuiKeyModFlags_Super)
		}

		var is_cut = ((is_shortcut_key && IsKeyPressedMap(ImGuiKey_X, true)) || (is_shift_key_only && IsKeyPressedMap(ImGuiKey_Delete, true))) && !is_readonly && !is_password && (!is_multiline || state.HasSelection())
		var is_copy = ((is_shortcut_key && IsKeyPressedMap(ImGuiKey_C, true)) || (is_ctrl_key_only && IsKeyPressedMap(ImGuiKey_Insert, true))) && !is_password && (!is_multiline || state.HasSelection())
		var is_paste = ((is_shortcut_key && IsKeyPressedMap(ImGuiKey_V, true)) || (is_shift_key_only && IsKeyPressedMap(ImGuiKey_Insert, true))) && !is_readonly
		var is_undo = ((is_shortcut_key && IsKeyPressedMap(ImGuiKey_Z, true)) && !is_readonly && is_undoable)
		var is_redo = ((is_shortcut_key && IsKeyPressedMap(ImGuiKey_Y, true)) || (is_osx_shift_shortcut && IsKeyPressedMap(ImGuiKey_Z, true))) && !is_readonly && is_undoable

		if IsKeyPressedMap(ImGuiKey_LeftArrow, true) {
			var k int = STB_TEXTEDIT_K_LEFT
			if is_startend_key_down {
				k = STB_TEXTEDIT_K_LINESTART
			} else {
				if is_wordmove_key_down {
					k = STB_TEXTEDIT_K_WORDLEFT
				}
			}
			state.OnKeyPressed(k | k_mask)
		} else if IsKeyPressedMap(ImGuiKey_RightArrow, true) {
			var k int = STB_TEXTEDIT_K_RIGHT
			if is_startend_key_down {
				k = STB_TEXTEDIT_K_LINEEND
			} else {
				if is_wordmove_key_down {
					k = STB_TEXTEDIT_K_WORDRIGHT
				}
			}
			state.OnKeyPressed(k | k_mask)
		} else if IsKeyPressedMap(ImGuiKey_UpArrow, true) && is_multiline {
			if io.KeyCtrl {
				setScrollY(draw_window, ImMax(draw_window.Scroll.y-g.FontSize, 0.0))
			} else {
				var k int = STB_TEXTEDIT_K_UP
				if is_startend_key_down {
					k = STB_TEXTEDIT_K_TEXTSTART
				}
				state.OnKeyPressed(k | k_mask)
			}
		} else if IsKeyPressedMap(ImGuiKey_DownArrow, true) && is_multiline {
			if io.KeyCtrl {
				setScrollY(draw_window, ImMin(draw_window.Scroll.y+g.FontSize, GetScrollMaxY()))
			} else {
				var k int = STB_TEXTEDIT_K_DOWN
				if is_startend_key_down {
					k = STB_TEXTEDIT_K_TEXTEND
				}
				state.OnKeyPressed(k | k_mask)
			}
		} else if IsKeyPressedMap(ImGuiKey_PageUp, true) && is_multiline {
			state.OnKeyPressed(STB_TEXTEDIT_K_PGUP | k_mask)
			scroll_y -= float(row_count_per_page) * g.FontSize
		} else if IsKeyPressedMap(ImGuiKey_PageDown, true) && is_multiline {
			state.OnKeyPressed(STB_TEXTEDIT_K_PGDOWN | k_mask)
			scroll_y += float(row_count_per_page) * g.FontSize
		} else if IsKeyPressedMap(ImGuiKey_Home, true) {
			var k int = STB_TEXTEDIT_K_LINESTART
			if is_startend_key_down {
				k = STB_TEXTEDIT_K_TEXTSTART
			}
			state.OnKeyPressed(k | k_mask)
		} else if IsKeyPressedMap(ImGuiKey_End, true) {
			var k int = STB_TEXTEDIT_K_LINEEND
			if is_startend_key_down {
				k = STB_TEXTEDIT_K_TEXTEND
			}
			state.OnKeyPressed(k | k_mask)
		} else if IsKeyPressedMap(ImGuiKey_Delete, true) && !is_readonly {
			state.OnKeyPressed(STB_TEXTEDIT_K_DELETE | k_mask)
		} else if IsKeyPressedMap(ImGuiKey_Backspace, true) && !is_readonly {
			if !state.HasSelection() {
				if is_wordmove_key_down {
					state.OnKeyPressed(STB_TEXTEDIT_K_WORDLEFT | STB_TEXTEDIT_K_SHIFT)
				} else if is_osx && io.KeySuper && !io.KeyAlt && !io.KeyCtrl {
					state.OnKeyPressed(STB_TEXTEDIT_K_LINESTART | STB_TEXTEDIT_K_SHIFT)
				}
			}
			state.OnKeyPressed(STB_TEXTEDIT_K_BACKSPACE | k_mask)
		} else if IsKeyPressedMap(ImGuiKey_Enter, true) || IsKeyPressedMap(ImGuiKey_KeyPadEnter, true) {
			var ctrl_enter_for_new_line = (flags & ImGuiInputTextFlags_CtrlEnterForNewLine) != 0
			if !is_multiline || (ctrl_enter_for_new_line && !io.KeyCtrl) || (!ctrl_enter_for_new_line && io.KeyCtrl) {
				enter_pressed = true
				clear_active_id = true
			} else if !is_readonly {
				var c rune = '\n' // Insert new line
				if InputTextFilterCharacter(&c, flags, callback, callback_user_data, ImGuiInputSource_Keyboard) {
					state.OnKeyPressed((int)(c))
				}
			}
		} else if IsKeyPressedMap(ImGuiKey_Escape, true) {
			clear_active_id = true
			cancel_edit = true
		} else if is_undo || is_redo {
			var action int = STB_TEXTEDIT_K_REDO
			if is_undo {
				action = STB_TEXTEDIT_K_UNDO
			}
			state.OnKeyPressed(action)
			state.ClearSelection()
		} else if is_shortcut_key && IsKeyPressedMap(ImGuiKey_A, true) {
			state.SelectAll()
			state.CursorFollow = true
		} else if is_cut || is_copy {
			// Cut, Copy
			if io.SetClipboardTextFn != nil {
				var ib int = 0
				if state.HasSelection() {
					ib = ImMinInt(state.Stb.select_start, state.Stb.select_end)
				}
				var ie int
				if state.HasSelection() {
					ie = ImMaxInt(state.Stb.select_start, state.Stb.select_end)
				} else {
					ie = state.CurLenW
				}

				var clipboard_data_len = ImTextCountUtf8BytesFromStr(state.TextW[ib:], state.TextW[ie:]) + 1
				var clipboard_data = make([]byte, clipboard_data_len)
				ImTextStrToUtf8(clipboard_data, clipboard_data_len, state.TextW[ib:], state.TextW[ie:])
				SetClipboardText(string(clipboard_data))
			}
			if is_cut {
				if !state.HasSelection() {
					state.SelectAll()
				}
				state.CursorFollow = true
				stb_textedit_cut(state, &state.Stb)
			}
		} else if is_paste {
			if clipboard := GetClipboardText(); clipboard != "" {
				// Filter pasted buffer
				var clipboard_len = (int)(len(clipboard))
				var clipboard_filtered = make([]ImWchar, clipboard_len+1)
				var clipboard_filtered_len int = 0
				for s := clipboard; len(s) > 0; {
					var c rune
					s = s[ImTextCharFromUtf8(&c, s):]
					if c == 0 {
						break
					}
					if !InputTextFilterCharacter(&c, flags, callback, callback_user_data, ImGuiInputSource_Clipboard) {
						continue
					}
					clipboard_filtered[clipboard_filtered_len] = (ImWchar)(c)
					clipboard_filtered_len++
				}
				clipboard_filtered[clipboard_filtered_len] = 0
				if clipboard_filtered_len > 0 { // If everything was filtered, ignore the pasting operation
					stb_textedit_paste(state, &state.Stb, clipboard_filtered, clipboard_filtered_len)
					state.CursorFollow = true
				}
			}
		}

		// Update render selection flag after events have been handled, so selection highlight can be displayed during the same frame.
		render_selection = render_selection || (state.HasSelection() && (RENDER_SELECTION_WHEN_INACTIVE || render_cursor))
	}

	// Process callbacks and apply result back to user's buffer.
	if g.ActiveId == id {
		IM_ASSERT(state != nil)
		var apply_new_text []byte
		var apply_new_text_length int
		if cancel_edit {
			// Restore initial value. Only return true if restoring to the initial value changes the current buffer contents.
			if !is_readonly && !bytes.Equal(*buf, state.InitialTextA) {
				// Push records into the undo stack so we can CTRL+Z the revert operation itself
				apply_new_text = state.InitialTextA
				apply_new_text_length = int(len(state.InitialTextA) - 1)
				var w_text []ImWchar
				if apply_new_text_length > 0 {
					resizeRune(&w_text, ImTextCountCharsFromUtf8(string(apply_new_text[:apply_new_text_length]))+1)
					remaining := string(apply_new_text[apply_new_text_length:])
					ImTextStrFromUtf8(w_text, int(len(w_text)), string(apply_new_text), &remaining)
				}

				var l int
				if apply_new_text_length > 0 {
					l = int(len(w_text)) - 1
				}

				stb_textedit_replace(state, &state.Stb, w_text, l)
			}
		}

		// When using 'ImGuiInputTextFlags_EnterReturnsTrue' as a special case we reapply the live buffer back to the input buffer before clearing ActiveId, even though strictly speaking it wasn't modified on this frame.
		// If we didn't do that, code like InputInt() with ImGuiInputTextFlags_EnterReturnsTrue would fail.
		// This also allows the user to use InputText() with ImGuiInputTextFlags_EnterReturnsTrue without maintaining any user-side storage (please note that if you use this property along ImGuiInputTextFlags_CallbackResize you can end up with your temporary string object unnecessarily allocating once a frame, either store your string data, either if you don't then don't use ImGuiInputTextFlags_CallbackResize).
		var apply_edit_back_to_user_buffer = !cancel_edit || (enter_pressed && (flags&ImGuiInputTextFlags_EnterReturnsTrue) != 0)
		if apply_edit_back_to_user_buffer {
			// Apply new value immediately - copy modified buffer back
			// Note that as soon as the input box is active, the in-widget value gets priority over any underlying modification of the input buffer
			// FIXME: We actually always render 'buf' when calling DrawList.AddText, making the comment above incorrect.
			// FIXME-OPT: CPU waste to do this every time the widget is active, should mark dirty state from the stb_textedit callbacks.
			if !is_readonly {
				state.TextAIsValid = true
				resize(&state.TextA, int(len(state.TextW)*4+1))
				ImTextStrToUtf8(state.TextA, int(len(state.TextA)), state.TextW, nil)
			}

			// User callback
			if (flags & (ImGuiInputTextFlags_CallbackCompletion | ImGuiInputTextFlags_CallbackHistory | ImGuiInputTextFlags_CallbackEdit | ImGuiInputTextFlags_CallbackAlways)) != 0 {
				IM_ASSERT(callback != nil)

				// The reason we specify the usage semantic (Completion/History) is that Completion needs to disable keyboard TABBING at the moment.
				var event_flag ImGuiInputTextFlags = 0
				var event_key ImGuiKey = ImGuiKey_COUNT
				if (flags&ImGuiInputTextFlags_CallbackCompletion) != 0 && IsKeyPressedMap(ImGuiKey_Tab, true) {
					event_flag = ImGuiInputTextFlags_CallbackCompletion
					event_key = ImGuiKey_Tab
				} else if (flags&ImGuiInputTextFlags_CallbackHistory) != 0 && IsKeyPressedMap(ImGuiKey_UpArrow, true) {
					event_flag = ImGuiInputTextFlags_CallbackHistory
					event_key = ImGuiKey_UpArrow
				} else if (flags&ImGuiInputTextFlags_CallbackHistory) != 0 && IsKeyPressedMap(ImGuiKey_DownArrow, true) {
					event_flag = ImGuiInputTextFlags_CallbackHistory
					event_key = ImGuiKey_DownArrow
				} else if (flags&ImGuiInputTextFlags_CallbackEdit != 0) && state.Edited {
					event_flag = ImGuiInputTextFlags_CallbackEdit
				} else if flags&ImGuiInputTextFlags_CallbackAlways != 0 {
					event_flag = ImGuiInputTextFlags_CallbackAlways
				}

				if event_flag != 0 {
					var callback_data ImGuiInputTextCallbackData
					callback_data.EventFlag = event_flag
					callback_data.Flags = flags
					callback_data.UserData = callback_user_data

					callback_data.EventKey = event_key
					callback_data.Buf = state.TextA
					callback_data.BufTextLen = state.CurLenA
					callback_data.BufSize = state.BufCapacityA
					callback_data.BufDirty = false

					// We have to convert from wchar-positions to UTF-8-positions, which can be pretty slow (an incentive to ditch the ImWchar buffer, see https://github.com/nothings/stb/issues/188)
					var text []rune = state.TextW
					var utf8_cursor_pos = ImTextCountUtf8BytesFromStr(text, text[state.Stb.cursor:])
					callback_data.CursorPos = utf8_cursor_pos
					var utf8_selection_start = ImTextCountUtf8BytesFromStr(text, text[state.Stb.select_start:])
					callback_data.SelectionStart = utf8_selection_start
					var utf8_selection_end = ImTextCountUtf8BytesFromStr(text, text[state.Stb.select_end:])
					callback_data.SelectionEnd = utf8_selection_end

					// Call user code
					callback(&callback_data)

					// Read back what user may have modified
					IM_ASSERT(bytes.Equal(callback_data.Buf, state.TextA)) // Invalid to modify those fields
					IM_ASSERT(callback_data.BufSize == state.BufCapacityA)
					IM_ASSERT(callback_data.Flags == flags)
					var buf_dirty = callback_data.BufDirty
					if callback_data.CursorPos != utf8_cursor_pos || buf_dirty {
						// TODO (port): check if ImTextCountCharsFromUtf8 works correctly, and the line below
						// state->Stb.cursor = ImTextCountCharsFromUtf8(callback_data.Buf, callback_data.Buf + callback_data.CursorPos)
						state.Stb.cursor = ImTextCountCharsFromUtf8(string(callback_data.Buf[:callback_data.CursorPos]))
						state.CursorFollow = true
					}
					if callback_data.SelectionStart != utf8_selection_start || buf_dirty {
						if callback_data.SelectionStart == callback_data.CursorPos {
							state.Stb.select_start = state.Stb.cursor
						} else {
							state.Stb.select_start = ImTextCountCharsFromUtf8(string(callback_data.Buf[:callback_data.SelectionStart]))
						}
					}
					if callback_data.SelectionEnd != utf8_selection_end || buf_dirty {
						if callback_data.SelectionEnd == callback_data.SelectionStart {
							state.Stb.select_end = state.Stb.select_start
						} else {
							state.Stb.select_end = ImTextCountCharsFromUtf8(string(callback_data.Buf[:callback_data.SelectionEnd]))
						}
					}
					if buf_dirty {
						IM_ASSERT(callback_data.BufTextLen == (int)(len(callback_data.Buf))) // You need to maintain BufTextLen if you change the text!
						if callback_data.BufTextLen > backup_current_text_length && is_resizable {
							resizeRune(&state.TextW, int(len(state.TextW))+(callback_data.BufTextLen-backup_current_text_length))
						}
						state.CurLenW = ImTextStrFromUtf8(state.TextW, int(len(state.TextW)), string(callback_data.Buf), nil)
						state.CurLenA = callback_data.BufTextLen // Assume correct length and valid UTF-8 from user, saves us an extra strlen()
						state.CursorAnimReset()
					}
				}
			}

			// Will copy result string if modified
			if !is_readonly && !bytes.Equal(state.TextA, *buf) {
				apply_new_text = state.TextA
				apply_new_text_length = state.CurLenA
			}
		}

		// Copy result to user buffer
		if apply_new_text != nil {
			// We cannot test for 'backup_current_text_length != apply_new_text_length' here because we have no guarantee that the size
			// of our owned buffer matches the size of the string object held by the user, and by design we allow InputText() to be used
			// without any storage on user's side.
			IM_ASSERT(apply_new_text_length >= 0)
			var buf_size int
			if is_resizable {
				var callback_data ImGuiInputTextCallbackData
				callback_data.EventFlag = ImGuiInputTextFlags_CallbackResize
				callback_data.Flags = flags
				callback_data.Buf = *buf
				callback_data.BufTextLen = apply_new_text_length
				callback_data.BufSize = ImMaxInt(int(len(*buf)), apply_new_text_length+1)
				callback_data.UserData = callback_user_data
				callback(&callback_data)
				// *buf = callback_data.Buf[:callback_data.BufSize]
				*buf = callback_data.Buf
				buf_size = callback_data.BufSize
				apply_new_text_length = ImMinInt(callback_data.BufTextLen, buf_size-1)
				IM_ASSERT(apply_new_text_length <= buf_size)
			}
			//IMGUI_DEBUG_LOG("InputText(\"%s\"): apply_new_text length %d\n", label, apply_new_text_length);

			// If the underlying buffer resize was denied or not carried to the next frame, apply_new_text_length+1 may be >= buf_size.
			copy(*buf, apply_new_text[:ImMinInt(apply_new_text_length+1, buf_size)])
			value_changed = true
		}

		// Clear temporary user storage
		state.Flags = ImGuiInputTextFlags_None
		state.UserCallback = nil
		state.UserCallbackData = nil
	}

	// Release active ID at the end of the function (so e.g. pressing Return still does a final application of the value)
	if clear_active_id && g.ActiveId == id {
		ClearActiveID()
	}

	// Render frame
	if !is_multiline {
		RenderNavHighlight(&frame_bb, id, 0)
		RenderFrame(frame_bb.Min, frame_bb.Max, GetColorU32FromID(ImGuiCol_FrameBg, 1), true, style.FrameRounding)
	}

	var clip_rect = ImVec4{frame_bb.Min.x, frame_bb.Min.y, frame_bb.Min.x + inner_size.x, frame_bb.Min.y + inner_size.y} // Not using frame_bb.Max because we have adjusted size
	var draw_pos ImVec2
	if is_multiline {
		draw_pos = draw_window.DC.CursorPos
	} else {
		draw_pos = frame_bb.Min.Add(style.FramePadding)
	}
	var text_size ImVec2

	// Set upper limit of single-line InputTextEx() at 2 million characters strings. The current pathological worst case is a long line
	// without any carriage return, which would makes ImFont::RenderText() reserve too many vertices and probably crash. Avoid it altogether.
	// Note that we only use this limit on single-line InputText(), so a pathologically large line on a InputTextMultiline() would still crash.
	var buf_display_max_length = 2 * 1024 * 1024
	var buf_display = *buf //-V595
	if buf_display_from_state {
		buf_display = state.TextA
	}
	var buf_display_end []byte // We have specialized paths below for setting the length
	if is_displaying_hint {
		buf_display = []byte(hint)
		buf_display_end = buf_display[len(hint):]
	}

	// Render text. We currently only render selection when the widget is active or while scrolling.
	// FIXME: We could remove the '&& render_cursor' to keep rendering selection when inactive.
	if render_cursor || render_selection {
		IM_ASSERT(state != nil)
		if !is_displaying_hint {
			buf_display_end = buf_display[state.CurLenA:]
		}

		// Render text (with cursor and selection)
		// This is going to be messy. We need to:
		// - Display the text (this alone can be more easily clipped)
		// - Handle scrolling, highlight selection, display cursor (those all requires some form of 1d.2d cursor position calculation)
		// - Measure text height (for scrollbar)
		// We are attempting to do most of that in **one main pass** to minimize the computation cost (non-negligible for large amount of text) + 2nd pass for selection rendering (we could merge them by an extra refactoring effort)
		// FIXME: This should occur on buf_display but we'd need to maintain cursor/select_start/select_end for UTF-8.
		var text_begin = state.TextW
		var cursor_offset, select_start_offset ImVec2

		{
			// Find lines numbers straddling 'cursor' (slot 0) and 'select_start' (slot 1) positions.
			var searches_input_ptr [2][]ImWchar
			var searches_result_line_no = [2]int{-1000, -1000}
			var searches_remaining = 0
			if render_cursor {
				searches_input_ptr[0] = text_begin[state.Stb.cursor:]
				searches_result_line_no[0] = -1
				searches_remaining++
			}
			if render_selection {
				searches_input_ptr[1] = text_begin[ImMinInt(state.Stb.select_start, state.Stb.select_end):]
				searches_result_line_no[1] = -1
				searches_remaining++
			}

			// Iterate all lines to find our line numbers
			// In multi-line mode, we never exit the loop until all lines are counted, so add one extra to the searches_remaining counter.
			if is_multiline {
				searches_remaining += 1
			}

			var line_count int = 0
			// for (const ImWchar* s = text_begin; (s = (const ImWchar*)wcschr((const wchar_t*)s, (wchar_t)'\n')) != nil; s++)  // FIXME-OPT: Could use this when wchar_t are 16-bit
			for s := text_begin; len(s) > 0; s = s[1:] {
				if s[0] == '\n' {
					line_count++
					if searches_result_line_no[0] == -1 && len(s) < len(searches_input_ptr[0]) {
						searches_result_line_no[0] = line_count
						searches_remaining--
						if searches_remaining <= 0 {
							break
						}
					}
					if searches_result_line_no[1] == -1 && len(s) >= len(searches_input_ptr[1]) {
						searches_result_line_no[1] = line_count
						searches_remaining--
						if searches_remaining <= 0 {
							break
						}
					}
				}
			}
			line_count++
			if searches_result_line_no[0] == -1 {
				searches_result_line_no[0] = line_count
			}
			if searches_result_line_no[1] == -1 {
				searches_result_line_no[1] = line_count
			}

			// Calculate 2d position by finding the beginning of the line and measuring distance
			cursor_offset.x = InputTextCalcTextSizeW(ImStrbolW(searches_input_ptr[0], text_begin), &searches_input_ptr[0], nil, false).x
			cursor_offset.y = float(searches_result_line_no[0]) * g.FontSize
			if searches_result_line_no[1] >= 0 {
				select_start_offset.x = InputTextCalcTextSizeW(ImStrbolW(searches_input_ptr[1], text_begin), &searches_input_ptr[1], nil, false).x
				select_start_offset.y = float(searches_result_line_no[1]) * g.FontSize
			}

			// Store text height (note that we haven't calculated text width at all, see GitHub issues #383, #1224)
			if is_multiline {
				text_size = ImVec2{inner_size.x, float(line_count) * g.FontSize}
			}
		}

		// Scroll
		if render_cursor && state.CursorFollow {
			// Horizontal scroll in chunks of quarter width
			if (flags & ImGuiInputTextFlags_NoHorizontalScroll) == 0 {
				var scroll_increment_x = inner_size.x * 0.25
				var visible_width = inner_size.x - style.FramePadding.x
				if cursor_offset.x < state.ScrollX {
					state.ScrollX = IM_FLOOR(ImMax(0.0, cursor_offset.x-scroll_increment_x))
				} else if cursor_offset.x-visible_width >= state.ScrollX {
					state.ScrollX = IM_FLOOR(cursor_offset.x - visible_width + scroll_increment_x)
				}
			} else {
				state.ScrollX = 0.0
			}

			// Vertical scroll
			if is_multiline {
				// Test if cursor is vertically visible
				if cursor_offset.y-g.FontSize < scroll_y {
					scroll_y = ImMax(0.0, cursor_offset.y-g.FontSize)
				} else if cursor_offset.y-inner_size.y >= scroll_y {
					scroll_y = cursor_offset.y - inner_size.y + style.FramePadding.y*2.0
				}
				var scroll_max_y = ImMax((text_size.y+style.FramePadding.y*2.0)-inner_size.y, 0.0)
				scroll_y = ImClamp(scroll_y, 0.0, scroll_max_y)
				draw_pos.y += (draw_window.Scroll.y - scroll_y) // Manipulate cursor pos immediately avoid a frame of lag
				draw_window.Scroll.y = scroll_y
			}

			state.CursorFollow = false
		}

		// Draw selection
		var draw_scroll = ImVec2{state.ScrollX, 0.0}
		if render_selection {
			var text_selected_begin = text_begin[ImMinInt(state.Stb.select_start, state.Stb.select_end):]
			var text_selected_end = text_begin[ImMaxInt(state.Stb.select_start, state.Stb.select_end):]

			var t float = 0.6
			if render_cursor {
				t = 1
			}

			var bg_color = GetColorU32FromID(ImGuiCol_TextSelectedBg, t) // FIXME: current code flow mandate that render_cursor is always true here, we are leaving the transparent one for tests.
			var bg_offy_up float = -1.0                                  // FIXME: those offsets should be part of the style? they don't play so well with multi-line selection.
			var bg_offy_dn float = 2.0

			if is_multiline {
				bg_offy_up = 0
				bg_offy_dn = 0.0
			}

			var rect_pos ImVec2 = draw_pos.Add(select_start_offset).Sub(draw_scroll)
			var slice = text_selected_begin[:len(text_selected_begin)-len(text_selected_end)]
			for i, p := range slice {
				if rect_pos.y > clip_rect.w+g.FontSize {
					break
				}
				if rect_pos.y < clip_rect.y {
					//p = (const ImWchar*)wmemchr((const wchar_t*)p, '\n', text_selected_end - p);  // FIXME-OPT: Could use this when wchar_t are 16-bit
					//p = p ? p + 1 : text_selected_end;
					for i < len(text_selected_begin)-len(text_selected_end) {
						if p == '\n' {
							break
						}
					}
				} else {
					var pp = slice[i:]
					var rect_size ImVec2 = InputTextCalcTextSizeW(slice[i:], &pp, nil, true)
					if rect_size.x <= 0.0 {
						rect_size.x = IM_FLOOR(g.Font.GetCharAdvance((ImWchar)(' ')) * 0.50) // So we can see selected empty lines
					}
					var rect = ImRect{rect_pos.Add(ImVec2{0.0, bg_offy_up - g.FontSize}), rect_pos.Add(ImVec2{rect_size.x, bg_offy_dn})}
					rect.ClipWith(ImRectFromVec4(&clip_rect))
					if rect.Overlaps(ImRectFromVec4(&clip_rect)) {
						draw_window.DrawList.AddRectFilled(rect.Min, rect.Max, bg_color, 0, 0)
					}
				}
				rect_pos.x = draw_pos.x - draw_scroll.x
				rect_pos.y += g.FontSize
			}
		}

		// We test for 'buf_display_max_length' as a way to avoid some pathological cases (e.g. single-line 1 MB string) which would make ImDrawList crash.
		if is_multiline || len(buf_display)-len(buf_display_end) < buf_display_max_length {
			var c = ImGuiCol_Text
			if is_displaying_hint {
				c = ImGuiCol_TextDisabled
			}
			var col = GetColorU32FromID(c, 1)
			var clip = &clip_rect
			if is_multiline {
				clip = nil
			}
			draw_window.DrawList.AddTextV(g.Font, g.FontSize, draw_pos.Sub(draw_scroll), col, string(buf_display[:len(buf_display)-len(buf_display_end)]), 0.0, clip)
		}

		// Draw blinking cursor
		if render_cursor {
			state.CursorAnim += io.DeltaTime
			var cursor_is_visible = (!g.IO.ConfigInputTextCursorBlink) || (state.CursorAnim <= 0.0) || ImFmod(state.CursorAnim, 1.20) <= 0.80
			var f = draw_pos.Add(cursor_offset).Sub(draw_scroll)
			var cursor_screen_pos ImVec2 = *ImFloorVec(&f)
			var cursor_screen_rect = ImRect{ImVec2{cursor_screen_pos.x, cursor_screen_pos.y - g.FontSize + 0.5}, ImVec2{cursor_screen_pos.x + 1.0, cursor_screen_pos.y - 1.5}}
			if cursor_is_visible && cursor_screen_rect.Overlaps(ImRectFromVec4(&clip_rect)) {
				bl := cursor_screen_rect.GetBL()
				draw_window.DrawList.AddLine(&cursor_screen_rect.Min, &bl, GetColorU32FromID(ImGuiCol_Text, 1), 1)
			}

			// Notify OS of text input position for advanced IME (-1 x offset so that Windows IME can cover our cursor. Bit of an extra nicety.)
			if !is_readonly {
				g.PlatformImePos = ImVec2{cursor_screen_pos.x - 1.0, cursor_screen_pos.y - g.FontSize}
			}
		}
	} else {
		// Render text only (no selection, no cursor)
		if is_multiline {
			buf_display_end = buf_display[:]
			text_size = ImVec2{inner_size.x, float(bytes.Count(buf_display, []byte{'\n'})) * g.FontSize} // We don't need width
		} else if !is_displaying_hint && g.ActiveId == id {
			buf_display_end = buf_display[state.CurLenA:]
		} else if !is_displaying_hint {
			buf_display_end = buf_display[len(buf_display):]
		}

		if is_multiline || (len(buf_display)-len(buf_display_end)) < buf_display_max_length {
			var c = ImGuiCol_Text
			if is_displaying_hint {
				c = ImGuiCol_TextDisabled
			}
			var col ImU32 = GetColorU32FromID(c, 1)
			var clip = &clip_rect
			if is_multiline {
				clip = nil
			}
			draw_window.DrawList.AddTextV(g.Font, g.FontSize, draw_pos, col, string(buf_display[:len(buf_display)-len(buf_display_end)]), 0.0, clip)
		}
	}

	if is_password && !is_displaying_hint {
		PopFont()
	}

	if is_multiline {
		Dummy(ImVec2{text_size.x, text_size.y + style.FramePadding.y})
		EndChild()
		EndGroup()
	}

	// Log as text
	if g.LogEnabled && (!is_password || is_displaying_hint) {
		LogSetNextTextDecoration("{", "}")
		LogRenderedText(&draw_pos, string(buf_display[:len(buf_display)-len(buf_display_end)]))
	}

	if label_size.x > 0 {
		RenderText(ImVec2{frame_bb.Max.x + style.ItemInnerSpacing.x, frame_bb.Min.y + style.FramePadding.y}, label, true)
	}

	if value_changed && (flags&ImGuiInputTextFlags_NoMarkEdited) == 0 {
		MarkItemEdited(id)
	}

	if (flags & ImGuiInputTextFlags_EnterReturnsTrue) != 0 {
		return enter_pressed
	} else {
		return value_changed
	}
}
