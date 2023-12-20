package imgui

import (
	"fmt"
	"reflect"
)

// Widgets: Regular Sliders
// - CTRL+Click on any slider to turn them into an input box. Manually input values aren't clamped and can go off-bounds.
// - Adjust format string to decorate the value with a prefix, a suffix, or adapt the editing and display precision e.g. "%.3f" -> 1.234; "%5.2 secs" -> 01.23 secs; "Biscuit: %.0f" -> Biscuit: 1; etc.
// - Format string may also be set to NULL or use the default format ("%f" or "%d").
// - Legacy: Pre-1.78 there are SliderXXX() function signatures that takes a final `power float=1.0' argument instead of the `ImGuiSliderFlags flags=0' argument.
//   If you get a warning converting a to float ImGuiSliderFlags, read https://github.com/ocornut/imgui/issues/3361

func SliderFloat(label string, v *float, v_min float, v_max float, format string /*= "%.3f"*/, flags ImGuiSliderFlags) bool {
	return SliderScalar(label, ImGuiDataType_Float, v, &v_min, &v_max, format, flags)
}

// adjust format to decorate the value with a prefix or a suffix for in-slider labels or unit display.

func SliderFloat2(label string, v *[2]float, v_min float, v_max float, format string /*= "%.3f"*/, flags ImGuiSliderFlags) bool {
	return SliderScalarN(label, ImGuiDataType_Float, v[:], v_min, v_max, format, flags)
}
func SliderFloat3(label string, v *[3]float, v_min float, v_max float, format string /*= "%.3f"*/, flags ImGuiSliderFlags) bool {
	return SliderScalarN(label, ImGuiDataType_Float, v[:], v_min, v_max, format, flags)
}
func SliderFloat4(label string, v *[4]float, v_min float, v_max float, format string /*= "%.3f"*/, flags ImGuiSliderFlags) bool {
	return SliderScalarN(label, ImGuiDataType_Float, v[:], v_min, v_max, format, flags)
}

func SliderAngle(label string, v_rad *float, v_degrees_min float /*= 0*/, v_degrees_max float /*= 0*/, format string /* = "%.0f deg"*/, flags ImGuiSliderFlags) bool {
	if format == "" {
		format = "%.0f deg"
	}
	var v_deg = (*v_rad) * 360.0 / (2 * IM_PI)
	var value_changed = SliderFloat(label, &v_deg, v_degrees_min, v_degrees_max, format, flags)
	*v_rad = v_deg * (2 * IM_PI) / 360.0
	return value_changed
}

func SliderInt(label string, v *int, v_min int, v_max int, format string /*= "%d"*/, flags ImGuiSliderFlags) bool {
	return SliderScalar(label, ImGuiDataType_S32, v, &v_min, &v_max, format, flags)
}
func SliderInt2(label string, v [2]int, v_min int, v_max int, format string /*= "%d"*/, flags ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderInt3(label string, v [3]int, v_min int, v_max int, format string /*= "%d"*/, flags ImGuiSliderFlags) bool {
	panic("not implemented")
}
func SliderInt4(label string, v [4]int, v_min int, v_max int, format string /*= "%d"*/, flags ImGuiSliderFlags) bool {
	panic("not implemented")
}

// Note: p_data, p_min and p_max are _pointers_ to a memory address holding the data. For a slider, they are all required.
// Read code of e.g. SliderFloat(), SliderInt() etc. or examples in 'Demo.Widgets.Data Types' to understand how to use this function directly.
func SliderScalar(label string, data_type ImGuiDataType, p_data any, p_min any, p_max any, format string, flags ImGuiSliderFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	var style = g.Style
	var id = window.GetIDs(label)
	var w = CalcItemWidth()

	var label_size = CalcTextSize(label, true, -1)
	var frame_bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(ImVec2{w, label_size.y + style.FramePadding.y*2.0})}

	var padding float
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var total_bb = ImRect{frame_bb.Min, frame_bb.Max.Add(ImVec2{padding, 0.0})}

	var temp_input_allowed = (flags & ImGuiSliderFlags_NoInput) == 0
	ItemSizeRect(&total_bb, style.FramePadding.y)

	var inputable ImGuiItemFlags
	if temp_input_allowed {
		inputable = ImGuiItemFlags_Inputable
	}

	if !ItemAdd(&total_bb, id, &frame_bb, inputable) {
		return false
	}

	// Default format string when passing nil
	if format == "" {
		format = DataTypeGetInfo(data_type).PrintFmt
	}

	// Tabbing or CTRL-clicking on Slider turns it into an input box
	var hovered = ItemHoverable(&frame_bb, id)
	var temp_input_is_active = temp_input_allowed && TempInputIsActive(id)
	if !temp_input_is_active {
		var focus_requested = temp_input_allowed && (g.LastItemData.StatusFlags&ImGuiItemStatusFlags_Focused) != 0
		var clicked = (hovered && g.IO.MouseClicked[0])
		if focus_requested || clicked || g.NavActivateId == id || g.NavInputId == id {
			SetActiveID(id, window)
			SetFocusID(id, window)
			FocusWindow(window)
			g.ActiveIdUsingNavDirMask |= (1 << ImGuiDir_Left) | (1 << ImGuiDir_Right)
			if temp_input_allowed && (focus_requested || (clicked && g.IO.KeyCtrl) || g.NavInputId == id) {
				temp_input_is_active = true
			}
		}
	}

	if temp_input_is_active {
		// Only clamp CTRL+Click input when ImGuiSliderFlags_AlwaysClamp is set
		var is_clamp_input = (flags & ImGuiSliderFlags_AlwaysClamp) != 0

		var min, max any
		if is_clamp_input {
			min = p_min
			max = p_max
		}

		return TempInputScalar(&frame_bb, id, label, data_type, p_data, format, min, max)
	}

	// Draw frame
	var c = ImGuiCol_FrameBg
	if g.ActiveId == id {
		c = ImGuiCol_FrameBgActive
	} else if hovered {
		c = ImGuiCol_FrameBgHovered
	}

	var frame_col = GetColorU32FromID(c, 1)
	RenderNavHighlight(&frame_bb, id, 0)
	RenderFrame(frame_bb.Min, frame_bb.Max, frame_col, true, g.Style.FrameRounding)

	// Slider behavior
	var grab_bb ImRect
	var value_changed = SliderBehavior(&frame_bb, id, data_type, p_data, p_min, p_max, format, flags, &grab_bb)
	if value_changed {
		MarkItemEdited(id)
	}

	// Render grab
	if grab_bb.Max.x > grab_bb.Min.x {
		var c = ImGuiCol_SliderGrab
		if g.ActiveId == id {
			c = ImGuiCol_SliderGrabActive
		}

		window.DrawList.AddRectFilled(grab_bb.Min, grab_bb.Max, GetColorU32FromID(c, 1), style.GrabRounding, 0)
	}

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	p_data_val := reflect.ValueOf(p_data).Elem()
	var value_buf = fmt.Sprintf(format, p_data_val)
	if g.LogEnabled {
		LogSetNextTextDecoration("{", "}")
	}
	RenderTextClipped(&frame_bb.Min, &frame_bb.Max, value_buf, nil, &ImVec2{0.5, 0.5}, nil)

	if label_size.x > 0.0 {
		RenderText(ImVec2{frame_bb.Max.x + style.ItemInnerSpacing.x, frame_bb.Min.y + style.FramePadding.y}, label, true)
	}

	return value_changed
}

// Add multiple sliders on 1 line for compact edition of multiple components
func SliderScalarN(label string, data_type ImGuiDataType, p_data []float, p_min float, p_max float, format string, flags ImGuiSliderFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	var value_changed = false
	BeginGroup()
	PushString(label)
	PushMultiItemsWidths(1, CalcItemWidth())
	for i := 0; i < len(p_data); i++ {
		PushID(int(i))
		if i > 0 {
			SameLine(0, g.Style.ItemInnerSpacing.x)
		}
		value_changed = SliderScalar("", data_type, &p_data[i], p_min, p_max, format, flags) || value_changed
		PopID()
		PopItemWidth()
	}
	PopID()

	SameLine(0, g.Style.ItemInnerSpacing.x)
	TextEx(label, 0)

	EndGroup()
	return value_changed
}

func VSliderFloat(label string, size ImVec2, v *float, v_min float, v_max float, format string /*= "%.3f"*/, flags ImGuiSliderFlags) bool {
	return VSliderScalar(label, size, ImGuiDataType_Float, v, &v_min, &v_max, format, flags)
}

func VSliderInt(label string, size ImVec2, v *int, v_min int, v_max int, format string /*= "%d"*/, flags ImGuiSliderFlags) bool {
	return VSliderScalar(label, size, ImGuiDataType_S32, v, &v_min, &v_max, format, flags)
}

func VSliderScalar(label string, size ImVec2, data_type ImGuiDataType, p_data any, p_min any, p_max any, format string, flags ImGuiSliderFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	var style = g.Style
	var id = window.GetIDs(label)

	var label_size = CalcTextSize(label, true, -1)
	var frame_bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)}

	var padding float
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var bb = ImRect{frame_bb.Min, frame_bb.Max.Add(ImVec2{padding, 0.0})}

	ItemSizeRect(&bb, style.FramePadding.y)
	if !ItemAdd(&frame_bb, id, nil, 0) {
		return false
	}

	// Default format string when passing nil
	if format == "" {
		format = DataTypeGetInfo(data_type).PrintFmt
	}

	var hovered = ItemHoverable(&frame_bb, id)
	if (hovered && g.IO.MouseClicked[0]) || g.NavActivateId == id || g.NavInputId == id {
		SetActiveID(id, window)
		SetFocusID(id, window)
		FocusWindow(window)
		g.ActiveIdUsingNavDirMask |= (1 << ImGuiDir_Up) | (1 << ImGuiDir_Down)
	}

	// Draw frame
	var c = ImGuiCol_FrameBg
	if g.ActiveId == id {
		c = ImGuiCol_FrameBgActive
	} else if hovered {
		c = ImGuiCol_FrameBgHovered
	}

	var frame_col = GetColorU32FromID(c, 1)
	RenderNavHighlight(&frame_bb, id, 0)
	RenderFrame(frame_bb.Min, frame_bb.Max, frame_col, true, g.Style.FrameRounding)

	// Slider behavior
	var grab_bb ImRect
	var value_changed = SliderBehavior(&frame_bb, id, data_type, p_data, p_min, p_max, format, flags|ImGuiSliderFlags_Vertical, &grab_bb)
	if value_changed {
		MarkItemEdited(id)
	}

	// Render grab
	if grab_bb.Max.y > grab_bb.Min.y {
		var c = ImGuiCol_SliderGrab
		if g.ActiveId == id {
			c = ImGuiCol_SliderGrabActive
		}
		window.DrawList.AddRectFilled(grab_bb.Min, grab_bb.Max, GetColorU32FromID(c, 1), style.GrabRounding, 0)
	}

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	// For the vertical slider we allow centered text to overlap the frame padding
	var value_buf = fmt.Sprintf(format, p_data)
	RenderTextClipped(&ImVec2{frame_bb.Min.x, frame_bb.Min.y + style.FramePadding.y}, &frame_bb.Max, value_buf, nil, &ImVec2{0.5, 0.0}, nil)
	if label_size.x > 0.0 {
		RenderText(ImVec2{frame_bb.Max.x + style.ItemInnerSpacing.x, frame_bb.Min.y + style.FramePadding.y}, label, true)
	}

	return value_changed
}
