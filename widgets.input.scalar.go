package imgui

import "fmt"

// Note: p_data, p_step, p_step_fast are _pointers_ to a memory address holding the data. For an Input widget, p_step and p_step_fast are optional.
// Read code of e.g. InputFloat(), InputInt() etc. or examples in 'Demo.Widgets.Data Types' to understand how to use this function directly.
func InputScalarInt64(label string, p_data, p_step, p_step_fast *int64, format string, flags ImGuiInputTextFlags) bool {
	data_type := ImGuiDataType_S64

	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style

	if format == "" {
		format = "%v"
	}

	var buf = []byte(fmt.Sprintf(format, *p_data))

	var value_changed = false
	if (flags & (ImGuiInputTextFlags_CharsHexadecimal | ImGuiInputTextFlags_CharsScientific)) == 0 {
		flags |= ImGuiInputTextFlags_CharsDecimal
	}
	flags |= ImGuiInputTextFlags_AutoSelectAll
	flags |= ImGuiInputTextFlags_NoMarkEdited // We call MarkItemEdited() ourselves by comparing the actual data rather than the string.

	if p_step != nil {
		var button_size = GetFrameHeight()

		BeginGroup() // The only purpose of the group here is to allow the caller to query item data e.g. IsItemActive()
		PushString(label)
		SetNextItemWidth(ImMax(1.0, CalcItemWidth()-(button_size+style.ItemInnerSpacing.x)*2))
		if InputText("", &buf, flags, nil, nil) { // PushId(label) + "" gives us the expected ID from outside point of view
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}

		// Step buttons
		var backup_frame_padding = style.FramePadding
		style.FramePadding.x = style.FramePadding.y
		var button_flags = ImGuiButtonFlags_Repeat | ImGuiButtonFlags_DontClosePopups
		if flags&ImGuiInputTextFlags_ReadOnly != 0 {
			BeginDisabled(true)
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("-", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '-', p_data, p_data, step)
			value_changed = true
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("+", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '+', p_data, p_data, step)
			value_changed = true
		}
		if (flags & ImGuiInputTextFlags_ReadOnly) != 0 {
			EndDisabled()
		}

		SameLine(0, style.ItemInnerSpacing.x)
		TextEx(label, 0)
		style.FramePadding = backup_frame_padding

		PopID()
		EndGroup()
	} else {
		if InputText(label, &buf, flags, nil, nil) {
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}
	}
	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}

func InputScalarInt64s(label string, p_data []int64, p_step, p_step_fast *int64, format string, flags ImGuiInputTextFlags) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	components := int(len(p_data))

	var g = GImGui
	var value_changed = false
	BeginGroup()
	PushString(label)
	PushMultiItemsWidths(components, CalcItemWidth())
	for i := int(0); i < components; i++ {
		p_data := &p_data[i]
		PushID(i)
		if i > 0 {
			SameLine(0, g.Style.ItemInnerSpacing.x)
		}
		value_changed = InputScalarInt64("", p_data, p_step, p_step_fast, format, flags) || value_changed
		PopID()
		PopItemWidth()
	}
	PopID()

	SameLine(0.0, g.Style.ItemInnerSpacing.x)
	TextEx(label, 0)

	EndGroup()
	return value_changed
}

// Note: p_data, p_step, p_step_fast are _pointers_ to a memory address holding the data. For an Input widget, p_step and p_step_fast are optional.
// Read code of e.g. InputFloat(), InputInt() etc. or examples in 'Demo.Widgets.Data Types' to understand how to use this function directly.
func InputScalarInt32(label string, p_data, p_step, p_step_fast *int32, format string, flags ImGuiInputTextFlags) bool {
	data_type := ImGuiDataType_S32

	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style

	if format == "" {
		format = "%v"
	}

	var buf = []byte(fmt.Sprintf(format, *p_data))

	var value_changed = false
	if (flags & (ImGuiInputTextFlags_CharsHexadecimal | ImGuiInputTextFlags_CharsScientific)) == 0 {
		flags |= ImGuiInputTextFlags_CharsDecimal
	}
	flags |= ImGuiInputTextFlags_AutoSelectAll
	flags |= ImGuiInputTextFlags_NoMarkEdited // We call MarkItemEdited() ourselves by comparing the actual data rather than the string.

	if p_step != nil {
		var button_size = GetFrameHeight()

		BeginGroup() // The only purpose of the group here is to allow the caller to query item data e.g. IsItemActive()
		PushString(label)
		SetNextItemWidth(ImMax(1.0, CalcItemWidth()-(button_size+style.ItemInnerSpacing.x)*2))
		if InputText("", &buf, flags, nil, nil) { // PushId(label) + "" gives us the expected ID from outside point of view
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}

		// Step buttons
		var backup_frame_padding = style.FramePadding
		style.FramePadding.x = style.FramePadding.y
		var button_flags = ImGuiButtonFlags_Repeat | ImGuiButtonFlags_DontClosePopups
		if flags&ImGuiInputTextFlags_ReadOnly != 0 {
			BeginDisabled(true)
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("-", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '-', p_data, p_data, step)
			value_changed = true
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("+", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '+', p_data, p_data, step)
			value_changed = true
		}
		if (flags & ImGuiInputTextFlags_ReadOnly) != 0 {
			EndDisabled()
		}

		SameLine(0, style.ItemInnerSpacing.x)
		TextEx(label, 0)
		style.FramePadding = backup_frame_padding

		PopID()
		EndGroup()
	} else {
		if InputText(label, &buf, flags, nil, nil) {
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}
	}
	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}

func InputScalarInt32s(label string, p_data []int32, p_step, p_step_fast *int32, format string, flags ImGuiInputTextFlags) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	components := int(len(p_data))

	var g = GImGui
	var value_changed = false
	BeginGroup()
	PushString(label)
	PushMultiItemsWidths(components, CalcItemWidth())
	for i := int(0); i < components; i++ {
		p_data := &p_data[i]
		PushID(i)
		if i > 0 {
			SameLine(0, g.Style.ItemInnerSpacing.x)
		}
		value_changed = InputScalarInt32("", p_data, p_step, p_step_fast, format, flags) || value_changed
		PopID()
		PopItemWidth()
	}
	PopID()

	SameLine(0.0, g.Style.ItemInnerSpacing.x)
	TextEx(label, 0)

	EndGroup()
	return value_changed
}

// Note: p_data, p_step, p_step_fast are _pointers_ to a memory address holding the data. For an Input widget, p_step and p_step_fast are optional.
// Read code of e.g. InputFloat(), InputInt() etc. or examples in 'Demo.Widgets.Data Types' to understand how to use this function directly.
func InputScalarFloat64(label string, p_data, p_step, p_step_fast *float64, format string, flags ImGuiInputTextFlags) bool {
	data_type := ImGuiDataType_Double

	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style

	if format == "" {
		format = "%v"
	}

	var buf = []byte(fmt.Sprintf(format, *p_data))

	var value_changed = false
	if (flags & (ImGuiInputTextFlags_CharsHexadecimal | ImGuiInputTextFlags_CharsScientific)) == 0 {
		flags |= ImGuiInputTextFlags_CharsDecimal
	}
	flags |= ImGuiInputTextFlags_AutoSelectAll
	flags |= ImGuiInputTextFlags_NoMarkEdited // We call MarkItemEdited() ourselves by comparing the actual data rather than the string.

	if p_step != nil {
		var button_size = GetFrameHeight()

		BeginGroup() // The only purpose of the group here is to allow the caller to query item data e.g. IsItemActive()
		PushString(label)
		SetNextItemWidth(ImMax(1.0, CalcItemWidth()-(button_size+style.ItemInnerSpacing.x)*2))
		if InputText("", &buf, flags, nil, nil) { // PushId(label) + "" gives us the expected ID from outside point of view
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}

		// Step buttons
		var backup_frame_padding = style.FramePadding
		style.FramePadding.x = style.FramePadding.y
		var button_flags = ImGuiButtonFlags_Repeat | ImGuiButtonFlags_DontClosePopups
		if flags&ImGuiInputTextFlags_ReadOnly != 0 {
			BeginDisabled(true)
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("-", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '-', p_data, p_data, step)
			value_changed = true
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("+", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '+', p_data, p_data, step)
			value_changed = true
		}
		if (flags & ImGuiInputTextFlags_ReadOnly) != 0 {
			EndDisabled()
		}

		SameLine(0, style.ItemInnerSpacing.x)
		TextEx(label, 0)
		style.FramePadding = backup_frame_padding

		PopID()
		EndGroup()
	} else {
		if InputText(label, &buf, flags, nil, nil) {
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}
	}
	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}

func InputScalarFloat64s(label string, p_data []float64, p_step, p_step_fast *float64, format string, flags ImGuiInputTextFlags) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	components := int(len(p_data))

	var g = GImGui
	var value_changed = false
	BeginGroup()
	PushString(label)
	PushMultiItemsWidths(components, CalcItemWidth())
	for i := int(0); i < components; i++ {
		p_data := &p_data[i]
		PushID(i)
		if i > 0 {
			SameLine(0, g.Style.ItemInnerSpacing.x)
		}
		value_changed = InputScalarFloat64("", p_data, p_step, p_step_fast, format, flags) || value_changed
		PopID()
		PopItemWidth()
	}
	PopID()

	SameLine(0.0, g.Style.ItemInnerSpacing.x)
	TextEx(label, 0)

	EndGroup()
	return value_changed
}

// Note: p_data, p_step, p_step_fast are _pointers_ to a memory address holding the data. For an Input widget, p_step and p_step_fast are optional.
// Read code of e.g. InputFloat(), InputInt() etc. or examples in 'Demo.Widgets.Data Types' to understand how to use this function directly.
func InputScalarFloat32(label string, p_data, p_step, p_step_fast *float32, format string, flags ImGuiInputTextFlags) bool {
	data_type := ImGuiDataType_Double

	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style

	if format == "" {
		format = "%v"
	}

	var buf = []byte(fmt.Sprintf(format, *p_data))

	var value_changed = false
	if (flags & (ImGuiInputTextFlags_CharsHexadecimal | ImGuiInputTextFlags_CharsScientific)) == 0 {
		flags |= ImGuiInputTextFlags_CharsDecimal
	}
	flags |= ImGuiInputTextFlags_AutoSelectAll
	flags |= ImGuiInputTextFlags_NoMarkEdited // We call MarkItemEdited() ourselves by comparing the actual data rather than the string.

	if p_step != nil {
		var button_size = GetFrameHeight()

		BeginGroup() // The only purpose of the group here is to allow the caller to query item data e.g. IsItemActive()
		PushString(label)
		SetNextItemWidth(ImMax(1.0, CalcItemWidth()-(button_size+style.ItemInnerSpacing.x)*2))
		if InputText("", &buf, flags, nil, nil) { // PushId(label) + "" gives us the expected ID from outside point of view
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}

		// Step buttons
		var backup_frame_padding = style.FramePadding
		style.FramePadding.x = style.FramePadding.y
		var button_flags = ImGuiButtonFlags_Repeat | ImGuiButtonFlags_DontClosePopups
		if flags&ImGuiInputTextFlags_ReadOnly != 0 {
			BeginDisabled(true)
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("-", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '-', p_data, p_data, step)
			value_changed = true
		}
		SameLine(0, style.ItemInnerSpacing.x)
		if ButtonEx("+", &ImVec2{button_size, button_size}, button_flags) {

			step := p_step
			if g.IO.KeyCtrl && p_step_fast != nil {
				step = p_step_fast
			}

			DataTypeApplyOp(data_type, '+', p_data, p_data, step)
			value_changed = true
		}
		if (flags & ImGuiInputTextFlags_ReadOnly) != 0 {
			EndDisabled()
		}

		SameLine(0, style.ItemInnerSpacing.x)
		TextEx(label, 0)
		style.FramePadding = backup_frame_padding

		PopID()
		EndGroup()
	} else {
		if InputText(label, &buf, flags, nil, nil) {
			value_changed = DataTypeApplyOpFromText(string(buf), string(g.InputTextState.InitialTextA), data_type, p_data, format)
		}
	}
	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}

func InputScalarFloat32s(label string, p_data []float32, p_step, p_step_fast *float32, format string, flags ImGuiInputTextFlags) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	components := int(len(p_data))

	var g = GImGui
	var value_changed = false
	BeginGroup()
	PushString(label)
	PushMultiItemsWidths(components, CalcItemWidth())
	for i := int(0); i < components; i++ {
		p_data := &p_data[i]
		PushID(i)
		if i > 0 {
			SameLine(0, g.Style.ItemInnerSpacing.x)
		}
		value_changed = InputScalarFloat32("", p_data, p_step, p_step_fast, format, flags) || value_changed
		PopID()
		PopItemWidth()
	}
	PopID()

	SameLine(0.0, g.Style.ItemInnerSpacing.x)
	TextEx(label, 0)

	EndGroup()
	return value_changed
}
