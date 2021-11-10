package imgui

import "fmt"

// Widgets: Drag Sliders
// - CTRL+Click on any drag box to turn them into an input box. Manually input values aren't clamped and can go off-bounds.
// - For all the Float2/Float3/Float4/Int2/Int3/Int4 versions of every functions, note that a 'v float[X]' function argument is the same as 'float* v', the array syntax is just a way to document the number of elements that are expected to be accessible. You can pass address of your first element out of a contiguous set, e.g. &myvector.x
// - Adjust format string to decorate the value with a prefix, a suffix, or adapt the editing and display precision e.g. "%.3f" -> 1.234; "%5.2 secs" -> 01.23 secs; "Biscuit: %.0f" -> Biscuit: 1; etc.
// - Format string may also be set to NULL or use the default format ("%f" or "%d").
// - Speed are per-pixel of mouse movement (v_speed=0.2: mouse needs to move by 5 pixels to increase value by 1). For gamepad/keyboard navigation, minimum speed is Max(v_speed, minimum_step_at_given_precision).
// - Use v_min < v_max to clamp edits to given limits. Note that CTRL+Click manual input can override those limits.
// - Use v_max/*= m*/,same with v_min = -FLT_MAX / INT_MIN to a clamping to a minimum.
// - We use the same sets of flags for DragXXX() and SliderXXX() functions as the features are the same and it makes it easier to swap them.
// - Legacy: Pre-1.78 there are DragXXX() function signatures that takes a final `power float=1.0' argument instead of the `ImGuiSliderFlags flags=0' argument.
//   If you get a warning converting a to float ImGuiSliderFlags, read https://github.com/ocornut/imgui/issues/3361
func DragFloat(label string, v *float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flags ImGuiSliderFlags) bool {
	return DragScalar(label, ImGuiDataType_Float, v, v_speed, &v_min, &v_max, format, flags)
}

func DragFloat2(label string, v *[2]float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragFloat3(label string, v *[3]float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragFloat4(label string, v *[4]float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "%.3f"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragFloatRange2(label string, v_current_min *float, v_current_max *float, v_speed float /*= 0*/, v_min float /*= 0*/, v_max float /*= 0*/, format string /*= "*/, format_max string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragInt(label string, v *int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
} // If v_min >= v_max we have no bound
func DragInt2(label string, v [2]int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragInt3(label string, v [3]int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragInt4(label string, v [4]int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "%d"*/, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}
func DragIntRange2(label string, v_current_min *int, v_current_max *int, v_speed float /*= 0*/, v_min int /*= 0*/, v_max int /*= 0*/, format string /*= "*/, format_max string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}

func GetMinimumStepAtDecimalPrecision(decimal_precision int) float {
	var min_steps = [10]float{1.0, 0.1, 0.01, 0.001, 0.0001, 0.00001, 0.000001, 0.0000001, 0.00000001, 0.000000001}
	if decimal_precision < 0 {
		return FLT_MIN
	}
	if decimal_precision < int(len(min_steps)) {
		return min_steps[decimal_precision]
	} else {
		return ImPow(10, float(-decimal_precision))
	}
}

// Note: p_data, p_min and p_max are _pointers_ to a memory address holding the data. For a Drag widget, p_min and p_max are optional.
// Read code of e.g. DragFloat(), DragInt() etc. or examples in 'Demo.Widgets.Data Types' to understand how to use this function directly.
func DragScalar(label string, data_type ImGuiDataType, p_data interface{}, v_speed float /*= 0*/, p_min interface{} /*= L*/, p_max interface{} /*= L*/, format string, flags ImGuiSliderFlags) bool {
	var window = GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var g = GImGui
	var style = g.Style
	var id ImGuiID = window.GetIDs(label)
	var w = CalcItemWidth()

	var label_size ImVec2 = CalcTextSize(label, true, -1)
	var frame_bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(ImVec2{w, label_size.y + style.FramePadding.y*2.0})}

	var padding float
	if label_size.x > 0.0 {
		padding = style.ItemInnerSpacing.x + label_size.x
	}

	var total_bb = ImRect{frame_bb.Min, frame_bb.Max.Add(ImVec2{padding, 0.0})}

	var temp_input_allowed bool = (flags & ImGuiSliderFlags_NoInput) == 0
	ItemSizeRect(&total_bb, style.FramePadding.y)

	var inputable_flags ImGuiItemFlags
	if temp_input_allowed {
		inputable_flags = ImGuiItemFlags_Inputable
	}

	if !ItemAdd(&total_bb, id, &frame_bb, inputable_flags) {
		return false
	}

	// Default format string when passing nil
	if format == "" {
		format = DataTypeGetInfo(data_type).PrintFmt
	}

	// Tabbing or CTRL-clicking on Drag turns it into an InputText
	var hovered = ItemHoverable(&frame_bb, id)
	var temp_input_is_active bool = temp_input_allowed && TempInputIsActive(id)
	if !temp_input_is_active {
		var focus_requested = temp_input_allowed && (g.LastItemData.StatusFlags&ImGuiItemStatusFlags_Focused) != 0
		var clicked = (hovered && g.IO.MouseClicked[0])
		var double_clicked = (hovered && g.IO.MouseDoubleClicked[0])
		if focus_requested || clicked || double_clicked || g.NavActivateId == id || g.NavInputId == id {
			SetActiveID(id, window)
			SetFocusID(id, window)
			FocusWindow(window)
			g.ActiveIdUsingNavDirMask = (1 << ImGuiDir_Left) | (1 << ImGuiDir_Right)
			if temp_input_allowed && (focus_requested || (clicked && g.IO.KeyCtrl) || double_clicked || g.NavInputId == id) {
				temp_input_is_active = true
			}
		}
		// Experimental: simple click (without moving) turns Drag into an InputText
		// FIXME: Currently polling ImGuiConfigFlags_IsTouchScreen, may either poll an hypothetical ImGuiBackendFlags_HasKeyboard and/or an explicit drag settings.
		if g.IO.ConfigDragClickToInputText && temp_input_allowed && !temp_input_is_active {
			if g.ActiveId == id && hovered && g.IO.MouseReleased[0] && !IsMouseDragPastThreshold(0, g.IO.MouseDragThreshold*DRAG_MOUSE_THRESHOLD_FACTOR) {
				g.NavInputId = id
				temp_input_is_active = true
			}
		}
	}

	if temp_input_is_active {
		// Only clamp CTRL+Click input when ImGuiSliderFlags_AlwaysClamp is set
		var is_clamp_input = (flags&ImGuiSliderFlags_AlwaysClamp) != 0 && (p_min == nil || p_max == nil || DataTypeCompare(data_type, p_min, p_max) < 0)

		var x, y interface{}
		if is_clamp_input {
			x, y = p_min, p_max
		}

		return TempInputScalar(&frame_bb, id, label, data_type, p_data, format, x, y)
	}

	var c = ImGuiCol_FrameBg
	switch {
	case g.ActiveId == id:
		c = ImGuiCol_FrameBgActive
	case hovered:
		c = ImGuiCol_FrameBgHovered
	}

	// Draw frame
	var frame_col = GetColorU32FromID(c, 1)
	RenderNavHighlight(&frame_bb, id, 0)
	RenderFrame(frame_bb.Min, frame_bb.Max, frame_col, true, style.FrameRounding)

	// Drag behavior
	var value_changed = DragBehavior(id, data_type, p_data, v_speed, p_min, p_max, format, flags)
	if value_changed {
		MarkItemEdited(id)
	}

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	var value_buf = fmt.Sprint(p_data)
	if g.LogEnabled {
		LogSetNextTextDecoration("{", "}")
	}
	RenderTextClipped(&frame_bb.Min, &frame_bb.Max, value_buf, nil, &ImVec2{0.5, 0.5}, nil)

	if label_size.x > 0.0 {
		RenderText(ImVec2{frame_bb.Max.x + style.ItemInnerSpacing.x, frame_bb.Min.y + style.FramePadding.y}, label, true)
	}

	return value_changed
}

func DragScalarN(label string, data_type ImGuiDataType, p_data interface{}, components int, v_speed float /*= 0*/, p_min interface{} /*= L*/, p_max interface{} /*= L*/, format string, flsgs ImGuiSliderFlags) bool {
	panic("not implemented")
}

func ScaleRatioFromValueT(v, v_min, v_max float, is_logarithmic bool, logarithmic_zero_epsilon, zero_deadzone_halfsize float) float {
	panic("not implemented")
}
func ScaleValueFromRatioT(v, v_min, v_max float, is_logarithmic bool, logarithmic_zero_epsilon, zero_deadzone_halfsize float) float {
	panic("not implemented")
}

func RoundScalarWithFormatT(format string, v float) float {
	panic("not implemented")
}

func DragBehaviorT(v *float, v_speed float, v_min, v_max float, format string, flags ImGuiSliderFlags) bool {
	var g = GImGui
	var axis ImGuiAxis = ImGuiAxis_X
	if (flags & ImGuiSliderFlags_Vertical) != 0 {
		axis = ImGuiAxis_Y
	}
	var is_clamped = (v_min < v_max)
	var is_logarithmic = (flags & ImGuiSliderFlags_Logarithmic) != 0
	var is_floating_point = true

	// Default tweak speed
	if v_speed == 0.0 && is_clamped && (v_max-v_min < FLT_MAX) {
		v_speed = (float)((v_max - v_min) * g.DragSpeedDefaultRatio)
	}
	// Inputs accumulates into g.DragCurrentAccum, which is flushed into the current value as soon as it makes a difference with our precision settings
	var adjust_delta float
	if g.ActiveIdSource == ImGuiInputSource_Mouse && IsMousePosValid(nil) && IsMouseDragPastThreshold(0, g.IO.MouseDragThreshold*DRAG_MOUSE_THRESHOLD_FACTOR) {
		switch axis {
		case ImGuiAxis_X:
			adjust_delta *= g.IO.MouseDelta.x
		case ImGuiAxis_Y:
			adjust_delta *= g.IO.MouseDelta.y
		}
		if g.IO.KeyAlt {
			adjust_delta *= 1.0 / 100.0
		}
		if g.IO.KeyShift {
			adjust_delta *= 10.0
		}
	} else if g.ActiveIdSource == ImGuiInputSource_Nav {
		var decimal_precision int = ImParseFormatPrecision(format, 3)
		amount := GetNavInputAmount2d(ImGuiNavDirSourceFlags_Keyboard|ImGuiNavDirSourceFlags_PadDPad, ImGuiInputReadMode_RepeatFast, 1.0/10.0, 10.0)
		switch axis {
		case ImGuiAxis_X:
			adjust_delta *= amount.x
		case ImGuiAxis_Y:
			adjust_delta *= amount.y
		}
		v_speed = ImMax(v_speed, GetMinimumStepAtDecimalPrecision(decimal_precision))
	}

	adjust_delta *= v_speed

	// For vertical drag we currently assume that Up=higher value (like we do with vertical sliders). This may become a parameter.
	if axis == ImGuiAxis_Y {
		adjust_delta = -adjust_delta
	}

	// For logarithmic use our range is effectively 0..1 so scale the delta into that range
	if is_logarithmic && (v_max-v_min < FLT_MAX) && ((v_max - v_min) > 0.000001) { // Epsilon to avoid /0
		adjust_delta /= (float)(v_max - v_min)
	}

	// Clear current value on activation
	// Avoid altering values and clamping when we are _already_ past the limits and heading in the same direction, so e.g. if range is 0..255, current value is 300 and we are pushing to the right side, keep the 300.
	var is_just_activated bool = g.ActiveIdIsJustActivated
	var is_already_past_limits_and_pushing_outward bool = is_clamped && ((*v >= v_max && adjust_delta > 0.0) || (*v <= v_min && adjust_delta < 0.0))
	if is_just_activated || is_already_past_limits_and_pushing_outward {
		g.DragCurrentAccum = 0.0
		g.DragCurrentAccumDirty = false
	} else if adjust_delta != 0.0 {
		g.DragCurrentAccum += adjust_delta
		g.DragCurrentAccumDirty = true
	}

	if !g.DragCurrentAccumDirty {
		return false
	}

	var v_cur float = *v
	var v_old_ref_for_accum_remainder float

	var logarithmic_zero_epsilon float // Only valid when is_logarithmic is true
	var zero_deadzone_halfsize float   // Drag widgets have no deadzone (as it doesn't make sense)
	if is_logarithmic {
		// When using logarithmic sliders, we need to clamp to avoid hitting zero, but our choice of clamp value greatly affects slider precision. We attempt to use the specified precision to estimate a good lower bound.
		var decimal_precision = ImParseFormatPrecision(format, 3)
		logarithmic_zero_epsilon = ImPow(0.1, (float)(decimal_precision))

		// Convert to parametric space, apply delta, convert back
		var v_old_parametric float = ScaleRatioFromValueT(v_cur, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)
		var v_new_parametric float = v_old_parametric + g.DragCurrentAccum
		v_cur = ScaleValueFromRatioT(v_new_parametric, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)
		v_old_ref_for_accum_remainder = v_old_parametric
	} else {
		v_cur += (float)(g.DragCurrentAccum)
	}

	// Round to user desired precision based on format string
	if flags&ImGuiSliderFlags_NoRoundToFormat == 0 {
		v_cur = RoundScalarWithFormatT(format, v_cur)
	}

	// Preserve remainder after rounding has been applied. This also allow slow tweaking of values.
	g.DragCurrentAccumDirty = false
	if is_logarithmic {
		// Convert to parametric space, apply delta, convert back
		var v_new_parametric float = ScaleRatioFromValueT(v_cur, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)
		g.DragCurrentAccum -= (float)(v_new_parametric - v_old_ref_for_accum_remainder)
	} else {
		g.DragCurrentAccum -= (float)(v_cur - *v)
	}

	// Lose zero sign for float/double
	if v_cur == (float)(-0) {
		v_cur = (float)(0)
	}

	// Clamp values (+ handle overflow/wrap-around for integer types)
	if *v != v_cur && is_clamped {
		if v_cur < v_min || (v_cur > *v && adjust_delta < 0.0 && !is_floating_point) {
			v_cur = v_min
		}
		if v_cur > v_max || (v_cur < *v && adjust_delta > 0.0 && !is_floating_point) {
			v_cur = v_max
		}
	}

	// Apply result
	if *v == v_cur {
		return false
	}

	*v = v_cur
	return true
}

func DragBehavior(id ImGuiID, data_type ImGuiDataType, v interface{}, v_speed float, n interface{}, x interface{}, t string, flags ImGuiSliderFlags) bool {
	// Read imgui.cpp "API BREAKING CHANGES" section for 1.78 if you hit this assert.
	IM_ASSERT_USER_ERROR((flags == 1 || (flags&ImGuiSliderFlags_InvalidMask_) == 0), "Invalid ImGuiSliderFlags flags! Has the 'float power' argument been mistakenly cast to flags? Call function with ImGuiSliderFlags_Logarithmic flags instead.")

	var g = GImGui
	if g.ActiveId == id {
		if g.ActiveIdSource == ImGuiInputSource_Mouse && !g.IO.MouseDown[0] {
			ClearActiveID()
		} else if g.ActiveIdSource == ImGuiInputSource_Nav && g.NavActivatePressedId == id && !g.ActiveIdIsJustActivated {
			ClearActiveID()
		}
	}
	if g.ActiveId != id {
		return false
	}
	if (g.LastItemData.InFlags&ImGuiItemFlags_ReadOnly != 0) || (flags&ImGuiSliderFlags_ReadOnly != 0) {
		return false
	}

	return DragBehaviorT(v.(*float), v_speed, n.(float), x.(float), t, flags)
}
