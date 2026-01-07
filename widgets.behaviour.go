package imgui

func SplitterBehavior(bb *ImRect, id ImGuiID, axis ImGuiAxis, size1 *float, size2 *float, min_size1 float, min_size2 float, hover_extend float, hover_visibility_delay float) bool {
	var g = GImGui
	var window = g.CurrentWindow

	var item_flags_backup = g.CurrentItemFlags
	g.CurrentItemFlags |= ImGuiItemFlags_NoNav | ImGuiItemFlags_NoNavDefaultFocus
	var item_add = ItemAdd(bb, id, nil, 0)
	g.CurrentItemFlags = item_flags_backup
	if !item_add {
		return false
	}

	var hovered, held bool
	var bb_interact = *bb
	if axis == ImGuiAxis_Y {
		bb_interact.ExpandVec(ImVec2{0.0, hover_extend})
	} else {
		bb_interact.ExpandVec(ImVec2{hover_extend, 0.0})
	}
	ButtonBehavior(&bb_interact, id, &hovered, &held, ImGuiButtonFlags_FlattenChildren|ImGuiButtonFlags_AllowItemOverlap)
	if hovered {
		g.LastItemData.StatusFlags |= ImGuiItemStatusFlags_HoveredRect // for IsItemHovered(), because bb_interact is larger than bb
	}
	if g.ActiveId != id {
		SetItemAllowOverlap()
	}

	if held || (hovered && g.HoveredIdPreviousFrame == id && g.HoveredIdTimer >= hover_visibility_delay) {
		if axis == ImGuiAxis_Y {
			SetMouseCursor(ImGuiMouseCursor_ResizeNS)
		} else {
			SetMouseCursor(ImGuiMouseCursor_ResizeEW)
		}
	}

	var bb_render = *bb
	if held {
		var mouse_delta_2d = g.IO.MousePos.Sub(g.ActiveIdClickOffset).Sub(bb_interact.Min)
		var mouse_delta = mouse_delta_2d.x
		if axis == ImGuiAxis_Y {
			mouse_delta = mouse_delta_2d.y
		}

		// Minimum pane size
		var size_1_maximum_delta = ImMax(0.0, *size1-min_size1)
		var size_2_maximum_delta = ImMax(0.0, *size2-min_size2)
		if mouse_delta < -size_1_maximum_delta {
			mouse_delta = -size_1_maximum_delta
		}
		if mouse_delta > size_2_maximum_delta {
			mouse_delta = size_2_maximum_delta
		}

		// Apply resize
		if mouse_delta != 0.0 {
			if mouse_delta < 0.0 {
				IM_ASSERT(*size1+mouse_delta >= min_size1)
			}
			if mouse_delta > 0.0 {
				IM_ASSERT(*size2-mouse_delta >= min_size2)
			}
			*size1 += mouse_delta
			*size2 -= mouse_delta
			if axis == ImGuiAxis_X {
				bb_render.Translate(ImVec2{mouse_delta, 0.0})
			} else {
				bb_render.Translate(ImVec2{0.0, mouse_delta})
			}
			MarkItemEdited(id)
		}
	}

	var c = ImGuiCol_Separator
	if held {
		c = ImGuiCol_SeparatorActive
	} else if hovered && g.HoveredIdTimer >= hover_visibility_delay {
		c = ImGuiCol_SeparatorHovered
	}

	// Render
	var col = GetColorU32FromID(c, 1)
	window.DrawList.AddRectFilled(bb_render.Min, bb_render.Max, col, 0.0, 0)

	return held
}

// FIXME: Move more of the code into SliderBehavior()
func sliderBehaviour(bb *ImRect, id ImGuiID, v *float, v_min float, v_max float, format string, flags ImGuiSliderFlags, out_grab_bb *ImRect) bool {
	var g = GImGui
	var style = g.Style

	var axis = ImGuiAxis_X
	if (flags & ImGuiSliderFlags_Vertical) != 0 {
		axis = ImGuiAxis_Y
	}

	var is_logarithmic = (flags & ImGuiSliderFlags_Logarithmic) != 0
	var is_floating_point = true

	var grab_padding float = 2.0
	var slider_sz = (bb.Max.Axis(axis) - bb.Min.Axis(axis)) - grab_padding*2.0
	var grab_sz = style.GrabMinSize
	var v_range = v_min - v_max
	if v_min < v_max {
		v_range = v_max - v_min
	}
	if !is_floating_point && v_range >= 0 { // v_range < 0 may happen on integer overflows
		grab_sz = ImMax((float)(slider_sz/(v_range+1)), style.GrabMinSize) // For integer sliders: if possible have the grab size represent 1 unit
	}
	grab_sz = ImMin(grab_sz, slider_sz)
	var slider_usable_sz = slider_sz - grab_sz
	var slider_usable_pos_min = bb.Min.Axis(axis) + grab_padding + grab_sz*0.5
	var slider_usable_pos_max = bb.Max.Axis(axis) - grab_padding - grab_sz*0.5

	var logarithmic_zero_epsilon float = 0.0 // Only valid when is_logarithmic is true
	var zero_deadzone_halfsize float = 0.0   // Only valid when is_logarithmic is true
	if is_logarithmic {
		// When using logarithmic sliders, we need to clamp to avoid hitting zero, but our choice of clamp value greatly affects slider precision. We attempt to use the specified precision to estimate a good lower bound.
		var decimal_precision int = 1
		if is_floating_point {
			decimal_precision = 3
		}
		logarithmic_zero_epsilon = ImPow(0.1, (float)(decimal_precision))
		zero_deadzone_halfsize = (style.LogSliderDeadzone * 0.5) / ImMax(slider_usable_sz, 1.0)
	}

	// Process interacting with the slider
	var value_changed = false
	if g.ActiveId == id {
		var set_new_value = false
		var clicked_t float = 0.0
		if g.ActiveIdSource == ImGuiInputSource_Mouse {
			if !g.IO.MouseDown[0] {
				ClearActiveID()
			} else {
				var mouse_abs_pos = g.IO.MousePos.Axis(axis)
				clicked_t = 0.0
				if slider_usable_sz > 0.0 {
					clicked_t = ImClamp((mouse_abs_pos-slider_usable_pos_min)/slider_usable_sz, 0.0, 1.0)
				}
				if axis == ImGuiAxis_Y {
					clicked_t = 1.0 - clicked_t
				}
				set_new_value = true
			}
		} else if g.ActiveIdSource == ImGuiInputSource_Nav {
			if g.ActiveIdIsJustActivated {
				g.SliderCurrentAccum = 0.0 // Reset any stored nav delta upon activation
				g.SliderCurrentAccumDirty = false
			}

			var input_delta2 = GetNavInputAmount2d(ImGuiNavDirSourceFlags_Keyboard|ImGuiNavDirSourceFlags_PadDPad, ImGuiInputReadMode_RepeatFast, 0.0, 0.0)
			var input_delta = -input_delta2.y
			if axis == ImGuiAxis_X {
				input_delta = input_delta2.x
			}
			if input_delta != 0.0 {
				var decimal_precision int
				if is_floating_point {
					decimal_precision = 3
				}
				if decimal_precision > 0 {
					input_delta /= 100.0 // Gamepad/keyboard tweak speeds in % of slider bounds
					if IsNavInputDown(ImGuiNavInput_TweakSlow) {
						input_delta /= 10.0
					}
				} else {
					if (v_range >= -100.0 && v_range <= 100.0) || IsNavInputDown(ImGuiNavInput_TweakSlow) {
						// Gamepad/keyboard tweak speeds in integer steps
						if input_delta < 0.0 {
							input_delta = -1.0 / v_range
						} else {
							input_delta = 1.0 / v_range
						}
					} else {
						input_delta /= 100.0
					}
				}
				if IsNavInputDown(ImGuiNavInput_TweakFast) {
					input_delta *= 10.0
				}

				g.SliderCurrentAccum += input_delta
				g.SliderCurrentAccumDirty = true
			}

			var delta = g.SliderCurrentAccum
			if g.NavActivatePressedId == id && !g.ActiveIdIsJustActivated {
				ClearActiveID()
			} else if g.SliderCurrentAccumDirty {
				clicked_t = ScaleRatioFromValueT(*v, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)

				if (clicked_t >= 1.0 && delta > 0.0) || (clicked_t <= 0.0 && delta < 0.0) { // This is to avoid applying the saturation when already past the limits
					set_new_value = false
					g.SliderCurrentAccum = 0.0 // If pushing up against the limits, don't continue to accumulate
				} else {
					set_new_value = true
					var old_clicked_t = clicked_t
					clicked_t = ImSaturate(clicked_t + delta)

					// Calculate what our "new" clicked_t will be, and thus how far we actually moved the slider, and subtract this from the accumulator
					var v_new = ScaleValueFromRatioT(clicked_t, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)
					if (flags & ImGuiSliderFlags_NoRoundToFormat) == 0 {
						v_new = RoundScalarWithFormatT(format, v_new)
					}
					var new_clicked_t = ScaleRatioFromValueT(v_new, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)

					if delta > 0 {
						g.SliderCurrentAccum -= ImMin(new_clicked_t-old_clicked_t, delta)
					} else {
						g.SliderCurrentAccum -= ImMax(new_clicked_t-old_clicked_t, delta)
					}
				}

				g.SliderCurrentAccumDirty = false
			}
		}

		if set_new_value {
			var v_new = ScaleValueFromRatioT(clicked_t, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)

			// Round to user desired precision based on format string
			if (flags & ImGuiSliderFlags_NoRoundToFormat) == 0 {
				v_new = RoundScalarWithFormatT(format, v_new)
			}

			// Apply result
			if *v != v_new {
				*v = v_new
				value_changed = true
			}
		}
	}

	if slider_sz < 1.0 {
		*out_grab_bb = ImRect{bb.Min, bb.Min}
	} else {
		// Output grab position so it can be displayed by the caller
		var grab_t = ScaleRatioFromValueT(*v, v_min, v_max, is_logarithmic, logarithmic_zero_epsilon, zero_deadzone_halfsize)
		if axis == ImGuiAxis_Y {
			grab_t = 1.0 - grab_t
		}
		var grab_pos = ImLerp(slider_usable_pos_min, slider_usable_pos_max, grab_t)
		if axis == ImGuiAxis_X {
			*out_grab_bb = ImRect{ImVec2{grab_pos - grab_sz*0.5, bb.Min.y + grab_padding}, ImVec2{grab_pos + grab_sz*0.5, bb.Max.y - grab_padding}}
		} else {
			*out_grab_bb = ImRect{ImVec2{bb.Min.x + grab_padding, grab_pos - grab_sz*0.5}, ImVec2{bb.Max.x - grab_padding, grab_pos + grab_sz*0.5}}
		}
	}

	return value_changed
}

// For 32-bit and larger types, slider bounds are limited to half the natural type range.
// So e.g. an integer Slider between INT_MAX-10 and INT_MAX will fail, but an integer Slider between INT_MAX/2-10 and INT_MAX/2 will be ok.
// It would be possible to lift that limitation with some work but it doesn't seem to be worth it for sliders.
func SliderBehavior(bb *ImRect, id ImGuiID, data_type ImGuiDataType, p_v any, p_min any, p_max any, format string, flags ImGuiSliderFlags, out_grab_bb *ImRect) bool {
	// Read imgui.cpp "API BREAKING CHANGES" section for 1.78 if you hit this assert.
	IM_ASSERT_USER_ERROR((flags == 1 || (flags&ImGuiSliderFlags_InvalidMask_) == 0), "Invalid ImGuiSliderFlags flag!  Has the 'float power' argument been mistakenly cast to flags? Call function with ImGuiSliderFlags_Logarithmic flags instead.")

	var g = GImGui
	if (g.LastItemData.InFlags&ImGuiItemFlags_ReadOnly != 0) || (flags&ImGuiSliderFlags_ReadOnly != 0) {
		return false
	}

	switch data_type {
	case ImGuiDataType_S32:
		IM_ASSERT(*p_min.(*int) >= IM_S32_MIN/2 && *p_max.(*int) <= IM_S32_MAX/2)
		// Convert int32 to float for slider calculations, then convert back
		v_float := float(*p_v.(*int))
		v_min_float := float(*p_min.(*int))
		v_max_float := float(*p_max.(*int))
		result := sliderBehaviour(bb, id, &v_float, v_min_float, v_max_float, format, flags, out_grab_bb)
		if result {
			*p_v.(*int) = int(v_float)
		}
		return result
	case ImGuiDataType_Float:
		IM_ASSERT(*p_min.(*float) >= -FLT_MAX/2.0 && *p_max.(*float) <= FLT_MAX/2.0)
		return sliderBehaviour(bb, id, p_v.(*float), *p_min.(*float), *p_max.(*float), format, flags, out_grab_bb)
	}

	// FIXME: support other slider types
	// https://github.com/ocornut/imgui/blob/5ee40c8d34bea3009cf462ec963225bd22067e5e/imgui_widgets.cpp#L2959
	IM_ASSERT(false)
	return false
}
