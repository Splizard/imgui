package imgui

import (
	"fmt"
	"unsafe"
)

// Helper for ColorPicker4()
func RenderArrowsForVerticalBar(draw_list *ImDrawList, pos, half_sz ImVec2, bar_w, alpha float) {
	var alpha8 = byte(IM_F32_TO_INT8_SAT(alpha))
	RenderArrowPointingAt(draw_list, ImVec2{pos.x + half_sz.x + 1, pos.y}, ImVec2{half_sz.x + 2, half_sz.y + 1}, ImGuiDir_Right, IM_COL32(0, 0, 0, alpha8))
	RenderArrowPointingAt(draw_list, ImVec2{pos.x + half_sz.x, pos.y}, half_sz, ImGuiDir_Right, IM_COL32(255, 255, 255, alpha8))
	RenderArrowPointingAt(draw_list, ImVec2{pos.x + bar_w - half_sz.x - 1, pos.y}, ImVec2{half_sz.x + 2, half_sz.y + 1}, ImGuiDir_Left, IM_COL32(0, 0, 0, alpha8))
	RenderArrowPointingAt(draw_list, ImVec2{pos.x + bar_w - half_sz.x, pos.y}, half_sz, ImGuiDir_Left, IM_COL32(255, 255, 255, alpha8))
}

// Widgets: Color Editor/Picker (tip: the ColorEdit* functions have a little color square that can be left-clicked to open a picker, and right-clicked to open an option menu.)
// - Note that in C++ a 'v float[X]' function argument is the _same_ as 'float* v', the array syntax is just a way to document the number of elements that are expected to be accessible.
// - You can pass the address of a first element float out of a contiguous structure, e.g. &myvector.x
func ColorEdit3(label string, col *[3]float, flags ImGuiColorEditFlags) bool {
	var col4 = [4]float{col[0], col[1], col[2], 1.0}
	result := ColorEdit4(label, &col4, flags|ImGuiColorEditFlags_NoAlpha)
	col[0] = col4[0]
	col[1] = col4[1]
	col[2] = col4[2]
	return result
}

// Edit colors components (each component in 0.0f..1.0f range).
// See enum ImGuiColorEditFlags_ for available options. e.g. Only access 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
// With typical options: Left-click on color square to open color picker. Right-click to open option menu. CTRL-Click over input fields to edit them and TAB to go to next item.
func ColorEdit4(label string, col *[4]float, flags ImGuiColorEditFlags) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	style := g.Style
	var square_sz = GetFrameHeight()
	var w_full = CalcItemWidth()
	var w_button float = 0.0
	if (flags & ImGuiColorEditFlags_NoSmallPreview) == 0 {
		w_button = (square_sz + style.ItemInnerSpacing.x)
	}
	var w_inputs = w_full - w_button
	g.NextItemData.ClearFlags()

	BeginGroup()
	PushString(label)

	// If we're not showing any slider there's no point in doing any HSV conversions
	var flags_untouched = flags
	if (flags & ImGuiColorEditFlags_NoInputs) != 0 {
		flags = (flags & (^ImGuiColorEditFlags_DisplayMask_)) | ImGuiColorEditFlags_DisplayRGB | ImGuiColorEditFlags_NoOptions
	}

	// Context menu: display and modify options (before defaults are applied)
	if flags&ImGuiColorEditFlags_NoOptions != 0 {
		ColorEditOptionsPopup(*col, flags)
	}

	// Read stored options
	if flags&ImGuiColorEditFlags_DisplayMask_ == 0 {
		flags |= (g.ColorEditOptions & ImGuiColorEditFlags_DisplayMask_)
	}
	if flags&ImGuiColorEditFlags_DataTypeMask_ == 0 {
		flags |= (g.ColorEditOptions & ImGuiColorEditFlags_DataTypeMask_)
	}
	if flags&ImGuiColorEditFlags_PickerMask_ == 0 {
		flags |= (g.ColorEditOptions & ImGuiColorEditFlags_PickerMask_)
	}
	if flags&ImGuiColorEditFlags_InputMask_ == 0 {
		flags |= (g.ColorEditOptions & ImGuiColorEditFlags_InputMask_)
	}
	flags |= (g.ColorEditOptions & ^(ImGuiColorEditFlags_DisplayMask_ | ImGuiColorEditFlags_DataTypeMask_ | ImGuiColorEditFlags_PickerMask_ | ImGuiColorEditFlags_InputMask_))

	// FIXME (port): these asserts always fail for some reason
	// IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_DisplayMask_))) // Check that only 1 is selected
	// IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_InputMask_)))   // Check that only 1 is selected

	var alpha = (flags & ImGuiColorEditFlags_NoAlpha) == 0
	var hdr = (flags & ImGuiColorEditFlags_HDR) != 0
	var components int = 3
	if alpha {
		components = 4
	}

	var a float = 1
	if alpha {
		a = col[3]
	}

	// Convert to the formats we need
	var f = [4]float{col[0], col[1], col[2], a}
	if (flags&ImGuiColorEditFlags_InputHSV) != 0 && (flags&ImGuiColorEditFlags_DisplayRGB) != 0 {
		ColorConvertHSVtoRGB(f[0], f[1], f[2], &f[0], &f[1], &f[2])
	} else if (flags&ImGuiColorEditFlags_InputRGB) != 0 && (flags&ImGuiColorEditFlags_DisplayHSV) != 0 {
		// Hue is lost when converting from greyscale rgb (saturation=0). Restore it.
		ColorConvertRGBtoHSV(f[0], f[1], f[2], &f[0], &f[1], &f[2])
		if g.ColorEditLastColor == [3]float32{col[0], col[1], col[2]} {
			if f[1] == 0 {
				f[0] = g.ColorEditLastHue
			}
			if f[2] == 0 {
				f[1] = g.ColorEditLastSat
			}
		}
	}
	var i = [4]int{IM_F32_TO_INT8_UNBOUND(f[0]), IM_F32_TO_INT8_UNBOUND(f[1]), IM_F32_TO_INT8_UNBOUND(f[2]), IM_F32_TO_INT8_UNBOUND(f[3])}

	var value_changed = false
	var value_changed_as_float = false

	var pos = window.DC.CursorPos
	var inputs_offset_x float = 0.0
	if style.ColorButtonPosition == ImGuiDir_Left {
		inputs_offset_x = w_button
	}
	window.DC.CursorPos.x = pos.x + inputs_offset_x

	if (flags&(ImGuiColorEditFlags_DisplayRGB|ImGuiColorEditFlags_DisplayHSV)) != 0 && (flags&ImGuiColorEditFlags_NoInputs) == 0 {
		// RGB/HSV 0..255 Sliders
		var w_item_one = max(1.0, IM_FLOOR((w_inputs-(style.ItemInnerSpacing.x)*float(components-1))/(float)(components)))
		var w_item_last = max(1.0, IM_FLOOR(w_inputs-(w_item_one+style.ItemInnerSpacing.x)*float(components-1)))

		var mformat = "M:000"
		if (flags & ImGuiColorEditFlags_Float) != 0 {
			mformat = "M:0.000"
		}

		var hide_prefix = (w_item_one <= CalcTextSize(mformat, true, -1).x)
		var ids = [4]string{"##X", "##Y", "##Z", "##W"}
		var fmt_table_int = [3][4]string{
			{"%3d", "%3d", "%3d", "%3d"},         // Short display
			{"R:%3d", "G:%3d", "B:%3d", "A:%3d"}, // Long display for RGBA
			{"H:%3d", "S:%3d", "V:%3d", "A:%3d"}, // Long display for HSVA
		}
		var fmt_table_float = [3][4]string{
			{"%0.3f", "%0.3f", "%0.3f", "%0.3f"},         // Short display
			{"R:%0.3f", "G:%0.3f", "B:%0.3f", "A:%0.3f"}, // Long display for RGBA
			{"H:%0.3f", "S:%0.3f", "V:%0.3f", "A:%0.3f"}, // Long display for HSVA
		}
		var fmt_idx int
		if hide_prefix {
			fmt_idx = 0
		} else {
			if (flags & ImGuiColorEditFlags_DisplayHSV) != 0 {
				fmt_idx = 2
			} else {
				fmt_idx = 1
			}
		}

		for n := int(0); n < components; n++ {
			if n > 0 {
				SameLine(0, style.ItemInnerSpacing.x)
			}
			if n+1 < components {
				SetNextItemWidth(w_item_one)
			} else {
				SetNextItemWidth(w_item_last)
			}

			var h1 float = 1
			if hdr {
				h1 = 0
			}
			var h2 int = 255
			if hdr {
				h2 = 0
			}

			// FIXME: When ImGuiColorEditFlags_HDR flag is passed HS values snap in weird ways when SV values go below 0.
			if (flags & ImGuiColorEditFlags_Float) != 0 {
				value_changed = DragFloat(ids[n], &f[n], 1.0/255.0, 0.0, h1, fmt_table_float[fmt_idx][n], 0) || value_changed
				value_changed_as_float = value_changed_as_float || value_changed
			} else {
				value_changed = DragInt(ids[n], &i[n], 1.0, 0, h2, fmt_table_int[fmt_idx][n], 0) || value_changed
			}
			if flags&ImGuiColorEditFlags_NoOptions == 0 {
				OpenPopupOnItemClick("context", 0)
			}
		}
	} else if (flags&ImGuiColorEditFlags_DisplayHex) != 0 && (flags&ImGuiColorEditFlags_NoInputs) == 0 {
		// RGB Hexadecimal Input
		var buf []byte
		if alpha {
			buf = []byte(fmt.Sprintf("#%02X%02X%02X%02X", ImClampInt(i[0], 0, 255), ImClampInt(i[1], 0, 255), ImClampInt(i[2], 0, 255), ImClampInt(i[3], 0, 255)))
		} else {
			buf = []byte(fmt.Sprintf("#%02X%02X%02X", ImClampInt(i[0], 0, 255), ImClampInt(i[1], 0, 255), ImClampInt(i[2], 0, 255)))
		}
		SetNextItemWidth(w_inputs)
		if InputText("##Text", &buf, ImGuiInputTextFlags_CharsHexadecimal|ImGuiInputTextFlags_CharsUppercase, nil, NewImDrawListSharedData()) {
			value_changed = true
			var p = buf
			for p[0] == '#' || ImCharIsBlankA(p[0]) {
				p = p[1:]
			}
			i[0] = 0
			i[1] = 0
			i[2] = 0
			i[3] = 0xFF // alpha default to 255 is not parsed by scanf (e.g. inputting #FFFFFF omitting alpha)
			if alpha {
				fmt.Sscanf(string(p), "%02X%02X%02X%02X", &i[0], &i[1], &i[2], &i[3]) // Treat at unsigned (%X is unsigned)
			} else {
				fmt.Sscanf(string(p), "%02X%02X%02X", &i[0], &i[1], &i[2])
			}
		}
		if flags&ImGuiColorEditFlags_NoOptions == 0 {
			OpenPopupOnItemClick("context", 0)
		}
	}

	var picker_active_window *ImGuiWindow
	if flags&ImGuiColorEditFlags_NoSmallPreview == 0 {
		var button_offset_x float
		if !((flags&ImGuiColorEditFlags_NoInputs) != 0 || (style.ColorButtonPosition == ImGuiDir_Left)) {
			button_offset_x = w_inputs + style.ItemInnerSpacing.x
		}
		window.DC.CursorPos = ImVec2{pos.x + button_offset_x, pos.y}

		var a float = 1
		if alpha {
			a = col[3]
		}

		var col_v4 = ImVec4{col[0], col[1], col[2], a}
		if ColorButton("##ColorButton", col_v4, flags, ImVec2{}) {
			if flags&ImGuiColorEditFlags_NoPicker == 0 {
				// Store current color and open a picker
				g.ColorPickerRef = col_v4
				OpenPopup("picker", 0)
				pos := g.LastItemData.Rect.GetBL().Add(ImVec2{-1, style.ItemSpacing.y})
				SetNextWindowPos(&pos, 0, ImVec2{})
			}
		}
		if flags&ImGuiColorEditFlags_NoOptions == 0 {
			OpenPopupOnItemClick("context", 0)
		}

		if BeginPopup("picker", 0) {
			picker_active_window = g.CurrentWindow
			if len(label) > 0 {
				TextEx(label, 0)
				Spacing()
			}
			var picker_flags_to_forward = ImGuiColorEditFlags_DataTypeMask_ | ImGuiColorEditFlags_PickerMask_ | ImGuiColorEditFlags_InputMask_ | ImGuiColorEditFlags_HDR | ImGuiColorEditFlags_NoAlpha | ImGuiColorEditFlags_AlphaBar
			var picker_flags = (flags_untouched & picker_flags_to_forward) | ImGuiColorEditFlags_DisplayMask_ | ImGuiColorEditFlags_NoLabel | ImGuiColorEditFlags_AlphaPreviewHalf
			SetNextItemWidth(square_sz * 12.0) // Use 256 + bar sizes?
			value_changed = ColorPicker4("##picker", col, picker_flags, []float{g.ColorPickerRef.x, g.ColorPickerRef.y, g.ColorPickerRef.z, g.ColorPickerRef.w}) || value_changed
			EndPopup()
		}
	}

	if len(label) > 0 && flags&ImGuiColorEditFlags_NoLabel == 0 {
		var text_offset_x float
		if flags&ImGuiColorEditFlags_NoInputs != 0 {
			text_offset_x = w_button
		} else {
			text_offset_x = w_full + style.ItemInnerSpacing.x
		}
		window.DC.CursorPos = ImVec2{pos.x + text_offset_x, pos.y + style.FramePadding.y}
		TextEx(label, 0)
	}

	// Convert back
	if value_changed && picker_active_window == nil {
		if !value_changed_as_float {
			for n := 0; n < 4; n++ {
				f[n] = float(i[n]) / 255.0
			}
		}
		if flags&ImGuiColorEditFlags_DisplayHSV != 0 && flags&ImGuiColorEditFlags_InputRGB != 0 {
			g.ColorEditLastHue = f[0]
			g.ColorEditLastSat = f[1]
			ColorConvertHSVtoRGB(f[0], f[1], f[2], &f[0], &f[1], &f[2])
			copy(g.ColorEditLastColor[:], f[:3])
		}
		if (flags&ImGuiColorEditFlags_DisplayRGB) != 0 && (flags&ImGuiColorEditFlags_InputHSV) != 0 {
			ColorConvertRGBtoHSV(f[0], f[1], f[2], &f[0], &f[1], &f[2])
		}

		col[0] = f[0]
		col[1] = f[1]
		col[2] = f[2]
		if alpha {
			col[3] = f[3]
		}
	}

	PopID()
	EndGroup()

	// Drag and Drop Target
	// NB: The flag test is merely an optional micro-optimization, BeginDragDropTarget() does the same test.
	if (g.LastItemData.StatusFlags&ImGuiItemStatusFlags_HoveredRect != 0) && (flags&ImGuiColorEditFlags_NoDragDrop) == 0 && BeginDragDropTarget() {
		var accepted_drag_drop = false
		if payload := AcceptDragDropPayload(IMGUI_PAYLOAD_TYPE_COLOR_3F, 0); payload != nil {
			data := payload.Data.([3]float32)
			copy(col[:], data[:3]) // Preserve alpha if any //-V512
			value_changed = true
			accepted_drag_drop = true
		}
		if payload := AcceptDragDropPayload(IMGUI_PAYLOAD_TYPE_COLOR_4F, 0); payload != nil {
			data := payload.Data.([4]float32)
			copy(col[:], data[:components])
			value_changed = true
			accepted_drag_drop = true
		}

		// Drag-drop payloads are always RGB
		if accepted_drag_drop && flags&ImGuiColorEditFlags_InputHSV != 0 {
			ColorConvertRGBtoHSV(col[0], col[1], col[2], &col[0], &col[1], &col[2])
		}
		EndDragDropTarget()
	}

	// When picker is being actively used, use its active id so IsItemActive() will function on ColorEdit4().
	if picker_active_window != nil && g.ActiveId != 0 && g.ActiveIdWindow == picker_active_window {
		g.LastItemData.ID = g.ActiveId
	}

	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	return value_changed
}

func ColorPicker3(label string, col *[3]float, flags ImGuiColorEditFlags) bool {
	var col4 = [4]float{col[0], col[1], col[2], 1.0}
	if !ColorPicker4(label, &col4, flags|ImGuiColorEditFlags_NoAlpha, nil) {
		return false
	}
	col[0] = col4[0]
	col[1] = col4[1]
	col[2] = col4[2]
	return true
}

// Note: ColorPicker4() only accesses 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
// (In C++ the 'float col[4]' notation for a function argument is equivalent to 'float* col', we only specify a size to facilitate understanding of the code.)
// FIXME: we adjust the big color square height based on item width, which may cause a flickering feedback loop (if automatic height makes a vertical scrollbar appears, affecting automatic width..)
// FIXME: this is trying to be aware of style.Alpha but not fully correct. Also, the color wheel will have overlapping glitches with (style.Alpha < 1.0)
func ColorPicker4(label string, col *[4]float, flags ImGuiColorEditFlags, ref_col []float) bool {
	g := GImGui
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	var draw_list = window.DrawList
	style := g.Style
	io := g.IO

	var width = CalcItemWidth()
	g.NextItemData.ClearFlags()

	PushString(label)
	BeginGroup()

	if (flags & ImGuiColorEditFlags_NoSidePreview) == 0 {
		flags |= ImGuiColorEditFlags_NoSmallPreview
	}

	// Context menu: display and store options.
	if flags&ImGuiColorEditFlags_NoOptions == 0 {
		ColorPickerOptionsPopup(col, flags)
	}

	// Read stored options
	if flags&ImGuiColorEditFlags_PickerMask_ == 0 {
		if (g.ColorEditOptions & ImGuiColorEditFlags_PickerMask_) != 0 {
			flags |= g.ColorEditOptions & ImGuiColorEditFlags_PickerMask_
		} else {
			flags |= ImGuiColorEditFlags_DefaultOptions_ & ImGuiColorEditFlags_PickerMask_
		}
	}
	if flags&ImGuiColorEditFlags_InputMask_ == 0 {
		if (g.ColorEditOptions & ImGuiColorEditFlags_InputMask_) != 0 {
			flags |= g.ColorEditOptions & ImGuiColorEditFlags_InputMask_
		} else {
			flags |= ImGuiColorEditFlags_DefaultOptions_ & ImGuiColorEditFlags_InputMask_
		}
	}
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_PickerMask_))) // Check that only 1 is selected
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_InputMask_)))  // Check that only 1 is selected
	if flags&ImGuiColorEditFlags_NoOptions == 0 {
		flags |= (g.ColorEditOptions & ImGuiColorEditFlags_AlphaBar)
	}

	// Setup
	var components int = 4
	if (flags & ImGuiColorEditFlags_NoAlpha) != 0 {
		components = 3
	}
	var alpha_bar = (flags&ImGuiColorEditFlags_AlphaBar) != 0 && (flags&ImGuiColorEditFlags_NoAlpha) == 0
	var bar float = 1
	if alpha_bar {
		bar = 2
	}
	var picker_pos = window.DC.CursorPos
	var square_sz = GetFrameHeight()
	var bars_width = square_sz                                                              // Arbitrary smallish width of Hue/Alpha picking bars
	var sv_picker_size = max(bars_width*1, width-bar*(bars_width+style.ItemInnerSpacing.x)) // Saturation/Value picking box
	var bar0_pos_x = picker_pos.x + sv_picker_size + style.ItemInnerSpacing.x
	var bar1_pos_x = bar0_pos_x + bars_width + style.ItemInnerSpacing.x
	var bars_triangles_half_sz = IM_FLOOR(bars_width * 0.20)

	var backup_initial_col [4]float
	copy(backup_initial_col[:], col[:components])

	var wheel_thickness = sv_picker_size * 0.08
	var wheel_r_outer = sv_picker_size * 0.50
	var wheel_r_inner = wheel_r_outer - wheel_thickness
	var wheel_center = ImVec2{picker_pos.x + (sv_picker_size+bars_width)*0.5, picker_pos.y + sv_picker_size*0.5}

	// Note: the triangle is displayed rotated with triangle_pa pointing to Hue, but most coordinates stays unrotated for logic.
	var triangle_r = wheel_r_inner - float((int)(sv_picker_size*0.027))
	var triangle_pa = ImVec2{triangle_r, 0.0}                           // Hue point.
	var triangle_pb = ImVec2{triangle_r * -0.5, triangle_r * -0.866025} // Black point.
	var triangle_pc = ImVec2{triangle_r * -0.5, triangle_r * +0.866025} // White point.

	var H, S, V = col[0], col[1], col[2]
	var R, G, B = col[0], col[1], col[2]
	if (flags & ImGuiColorEditFlags_InputRGB) != 0 {
		// Hue is lost when converting from greyscale rgb (saturation=0). Restore it.
		ColorConvertRGBtoHSV(R, G, B, &H, &S, &V)
		if g.ColorEditLastColor == [3]float32{col[0], col[1], col[2]} {
			if S == 0 {
				H = g.ColorEditLastHue
			}
			if V == 0 {
				S = g.ColorEditLastSat
			}
		}
	} else if flags&ImGuiColorEditFlags_InputHSV != 0 {
		ColorConvertHSVtoRGB(H, S, V, &R, &G, &B)
	}

	var value_changed, value_changed_h, value_changed_sv bool

	PushItemFlag(ImGuiItemFlags_NoNav, true)
	if (flags & ImGuiColorEditFlags_PickerHueWheel) != 0 {
		// Hue wheel + SV triangle logic
		InvisibleButton("hsv", ImVec2{sv_picker_size + style.ItemInnerSpacing.x + bars_width, sv_picker_size}, 0)
		if IsItemActive() {
			var initial_off = g.IO.MouseClickedPos[0].Sub(wheel_center)
			var current_off = g.IO.MousePos.Sub(wheel_center)
			var initial_dist2 = ImLengthSqrVec2(initial_off)
			if initial_dist2 >= (wheel_r_inner-1)*(wheel_r_inner-1) && initial_dist2 <= (wheel_r_outer+1)*(wheel_r_outer+1) {
				// Interactive with Hue wheel
				H = ImAtan2(current_off.y, current_off.x) / IM_PI * 0.5
				if H < 0.0 {
					H += 1.0
				}
				value_changed = true
				value_changed_h = true
			}
			var cos_hue_angle = ImCos(-H * 2.0 * IM_PI)
			var sin_hue_angle = ImSin(-H * 2.0 * IM_PI)
			if ImTriangleContainsPoint(&triangle_pa, &triangle_pb, &triangle_pc, ImRotate(&initial_off, cos_hue_angle, sin_hue_angle)) {
				// Interacting with SV triangle
				var current_off_unrotated = *ImRotate(&current_off, cos_hue_angle, sin_hue_angle)
				if !ImTriangleContainsPoint(&triangle_pa, &triangle_pb, &triangle_pc, &current_off_unrotated) {
					current_off_unrotated = ImTriangleClosestPoint(&triangle_pa, &triangle_pb, &triangle_pc, &current_off_unrotated)
				}
				var uu, vv, ww float
				ImTriangleBarycentricCoords(&triangle_pa, &triangle_pb, &triangle_pc, &current_off_unrotated, &uu, &vv, &ww)
				V = ImClamp(1.0-vv, 0.0001, 1.0)
				S = ImClamp(uu/V, 0.0001, 1.0)
				value_changed = true
				value_changed_sv = true
			}
		}
		if flags&ImGuiColorEditFlags_NoOptions != 0 {
			OpenPopupOnItemClick("context", 0)
		}
	} else if flags&ImGuiColorEditFlags_PickerHueBar != 0 {
		// SV rectangle logic
		InvisibleButton("sv", ImVec2{sv_picker_size, sv_picker_size}, 0)
		if IsItemActive() {
			S = ImSaturate((io.MousePos.x - picker_pos.x) / (sv_picker_size - 1))
			V = 1.0 - ImSaturate((io.MousePos.y-picker_pos.y)/(sv_picker_size-1))
			value_changed = true
			value_changed_sv = true
		}
		if flags&ImGuiColorEditFlags_NoOptions != 0 {
			OpenPopupOnItemClick("context", 0)
		}

		// Hue bar logic
		SetCursorScreenPos(ImVec2{bar0_pos_x, picker_pos.y})
		InvisibleButton("hue", ImVec2{bars_width, sv_picker_size}, 0)
		if IsItemActive() {
			H = ImSaturate((io.MousePos.y - picker_pos.y) / (sv_picker_size - 1))
			value_changed = true
			value_changed_h = true
		}
	}

	// Alpha bar logic
	if alpha_bar {
		SetCursorScreenPos(ImVec2{bar1_pos_x, picker_pos.y})
		InvisibleButton("alpha", ImVec2{bars_width, sv_picker_size}, 0)
		if IsItemActive() {
			col[3] = 1.0 - ImSaturate((io.MousePos.y-picker_pos.y)/(sv_picker_size-1))
			value_changed = true
		}
	}
	PopItemFlag() // ImGuiItemFlags_NoNav

	if flags&ImGuiColorEditFlags_NoSidePreview == 0 {
		SameLine(0, style.ItemInnerSpacing.x)
		BeginGroup()
	}

	if flags&ImGuiColorEditFlags_NoLabel == 0 {
		if len(label) > 0 {
			if (flags & ImGuiColorEditFlags_NoSidePreview) != 0 {
				SameLine(0, style.ItemInnerSpacing.x)
			}
			TextEx(label, 0)
		}
	}

	if flags&ImGuiColorEditFlags_NoSidePreview == 0 {
		PushItemFlag(ImGuiItemFlags_NoNavDefaultFocus, true)

		a := col[3]
		if (flags & ImGuiColorEditFlags_NoAlpha) != 0 {
			a = 1.0
		}

		var col_v4 = ImVec4{col[0], col[1], col[2], a}
		if flags&ImGuiColorEditFlags_NoLabel != 0 {
			Text("Current")
		}

		var sub_flags_to_forward = ImGuiColorEditFlags_InputMask_ | ImGuiColorEditFlags_HDR | ImGuiColorEditFlags_AlphaPreview | ImGuiColorEditFlags_AlphaPreviewHalf | ImGuiColorEditFlags_NoTooltip
		ColorButton("##current", col_v4, (flags & sub_flags_to_forward), ImVec2{square_sz * 3, square_sz * 2})
		if ref_col != nil {
			Text("Original")

			a := ref_col[3]
			if (flags & ImGuiColorEditFlags_NoAlpha) != 0 {
				a = 1.0
			}

			var ref_col_v4 = ImVec4{ref_col[0], ref_col[1], ref_col[2], a}
			if ColorButton("##original", ref_col_v4, (flags & sub_flags_to_forward), ImVec2{square_sz * 3, square_sz * 2}) {
				copy(col[:], ref_col[:components])
				value_changed = true
			}
		}
		PopItemFlag()
		EndGroup()
	}

	// Convert back color to RGB
	if value_changed_h || value_changed_sv {
		if (flags & ImGuiColorEditFlags_InputRGB) != 0 {
			h := H
			if H >= 1.0 {
				h = H - 10*1e-6
			}
			s := S
			if S > 0.0 {
				s = S - 10*1e-6
			}
			v := V
			if V > 0.0 {
				v = 1e-6
			}

			ColorConvertHSVtoRGB(h, s, v, &col[0], &col[1], &col[2])
			g.ColorEditLastHue = H
			g.ColorEditLastSat = S
			copy(g.ColorEditLastColor[:], col[:3])
		} else if flags&ImGuiColorEditFlags_InputHSV != 0 {
			col[0] = H
			col[1] = S
			col[2] = V
		}
	}

	// R,G,B and H,S,V slider color editor
	var value_changed_fix_hue_wrap = false
	if flags&ImGuiColorEditFlags_NoInputs == 0 {
		bar := bar0_pos_x
		if alpha_bar {
			bar = bar1_pos_x
		}

		PushItemWidth(bar + bars_width - picker_pos.x)
		var sub_flags_to_forward = ImGuiColorEditFlags_DataTypeMask_ | ImGuiColorEditFlags_InputMask_ | ImGuiColorEditFlags_HDR | ImGuiColorEditFlags_NoAlpha | ImGuiColorEditFlags_NoOptions | ImGuiColorEditFlags_NoSmallPreview | ImGuiColorEditFlags_AlphaPreview | ImGuiColorEditFlags_AlphaPreviewHalf
		var sub_flags = (flags & sub_flags_to_forward) | ImGuiColorEditFlags_NoPicker
		if flags&ImGuiColorEditFlags_DisplayRGB != 0 || (flags&ImGuiColorEditFlags_DisplayMask_) == 0 {
			if ColorEdit4("##rgb", col, sub_flags|ImGuiColorEditFlags_DisplayRGB) {
				// FIXME: Hackily differentiating using the DragInt (ActiveId != 0 && !ActiveIdAllowOverlap) vs. using the InputText or DropTarget.
				// For the later we don't want to run the hue-wrap canceling code. If you are well versed in HSV picker please provide your input! (See #2050)
				value_changed_fix_hue_wrap = (g.ActiveId != 0 && !g.ActiveIdAllowOverlap)
				value_changed = true
			}
		}
		if flags&ImGuiColorEditFlags_DisplayHSV != 0 || (flags&ImGuiColorEditFlags_DisplayMask_) == 0 {
			value_changed = value_changed || (ColorEdit4("##hsv", col, sub_flags|ImGuiColorEditFlags_DisplayHSV))
		}
		if flags&ImGuiColorEditFlags_DisplayHex != 0 || (flags&ImGuiColorEditFlags_DisplayMask_) == 0 {
			value_changed = value_changed || (ColorEdit4("##hex", col, sub_flags|ImGuiColorEditFlags_DisplayHex))
		}
		PopItemWidth()
	}

	// Try to cancel hue wrap (after ColorEdit4 call), if any
	if value_changed_fix_hue_wrap && (flags&ImGuiColorEditFlags_InputRGB) != 0 {
		var new_H, new_S, new_V float
		ColorConvertRGBtoHSV(col[0], col[1], col[2], &new_H, &new_S, &new_V)
		if new_H <= 0 && H > 0 {
			if new_V <= 0 && V != new_V {
				v := new_V
				if new_V <= 0 {
					v = V * 0.5
				}
				ColorConvertHSVtoRGB(H, S, v, &col[0], &col[1], &col[2])
			} else if new_S <= 0 {
				s := new_S
				if new_S <= 0 {
					s = S * 0.5
				}
				ColorConvertHSVtoRGB(H, s, new_V, &col[0], &col[1], &col[2])
			}
		}
	}

	if value_changed {
		if (flags & ImGuiColorEditFlags_InputRGB) != 0 {
			R = col[0]
			G = col[1]
			B = col[2]
			ColorConvertRGBtoHSV(R, G, B, &H, &S, &V)
			if g.ColorEditLastColor == [3]float32{col[0], col[1], col[2]} { // Fix local Hue as display below will use it immediately.
				if S == 0 {
					H = g.ColorEditLastHue
				}
				if V == 0 {
					S = g.ColorEditLastSat
				}
			}
		} else if flags&ImGuiColorEditFlags_InputHSV != 0 {
			H = col[0]
			S = col[1]
			V = col[2]
			ColorConvertHSVtoRGB(H, S, V, &R, &G, &B)
		}
	}

	var style_alpha8 = byte(IM_F32_TO_INT8_SAT(style.Alpha))
	var col_black = IM_COL32(0, 0, 0, style_alpha8)
	var col_white = IM_COL32(255, 255, 255, style_alpha8)
	var col_midgrey = IM_COL32(128, 128, 128, style_alpha8)
	var col_hues = [6 + 1]uint{IM_COL32(255, 0, 0, style_alpha8), IM_COL32(255, 255, 0, style_alpha8), IM_COL32(0, 255, 0, style_alpha8), IM_COL32(0, 255, 255, style_alpha8), IM_COL32(0, 0, 255, style_alpha8), IM_COL32(255, 0, 255, style_alpha8), IM_COL32(255, 0, 0, style_alpha8)}

	var hue_color_f = ImVec4{1, 1, 1, style.Alpha}
	ColorConvertHSVtoRGB(H, 1, 1, &hue_color_f.x, &hue_color_f.y, &hue_color_f.z)
	var hue_color32 = ColorConvertFloat4ToU32(hue_color_f)
	var user_col32_striped_of_alpha = ColorConvertFloat4ToU32(ImVec4{R, G, B, style.Alpha}) // Important: this is still including the main rendering/style alpha!!

	var sv_cursor_pos ImVec2

	if (flags & ImGuiColorEditFlags_PickerHueWheel) != 0 {
		// Render Hue Wheel
		var aeps = 0.5 / wheel_r_outer // Half a pixel arc length in radians (2pi cancels out).
		var segment_per_arc = max(4, (int)(wheel_r_outer/12))
		for n := 0; n < 6; n++ {
			var a0 = float(n)/6.0*2.0*IM_PI - aeps
			var a1 = float(n+1.0)/6.0*2.0*IM_PI + aeps
			var vert_start_idx = int(len(draw_list.VtxBuffer))
			draw_list.PathArcTo(wheel_center, (wheel_r_inner+wheel_r_outer)*0.5, a0, a1, segment_per_arc)
			draw_list.PathStroke(col_white, 0, wheel_thickness)
			var vert_end_idx = int(len(draw_list.VtxBuffer))

			// Paint colors over existing vertices
			var gradient_p0 = ImVec2{wheel_center.x + ImCos(a0)*wheel_r_inner, wheel_center.y + ImSin(a0)*wheel_r_inner}
			var gradient_p1 = ImVec2{wheel_center.x + ImCos(a1)*wheel_r_inner, wheel_center.y + ImSin(a1)*wheel_r_inner}
			ShadeVertsLinearColorGradientKeepAlpha(draw_list, vert_start_idx, vert_end_idx, gradient_p0, gradient_p1, col_hues[n], col_hues[n+1])
		}

		// Render Cursor + preview on Hue Wheel
		var cos_hue_angle = ImCos(H * 2.0 * IM_PI)
		var sin_hue_angle = ImSin(H * 2.0 * IM_PI)
		var hue_cursor_pos = ImVec2{wheel_center.x + cos_hue_angle*(wheel_r_inner+wheel_r_outer)*0.5, wheel_center.y + sin_hue_angle*(wheel_r_inner+wheel_r_outer)*0.5}
		var hue_cursor_rad float
		if value_changed_h {
			hue_cursor_rad = wheel_thickness * 0.65
		} else {
			hue_cursor_rad = wheel_thickness * 0.55
		}
		var hue_cursor_segments = ImClampInt((int)(hue_cursor_rad/1.4), 9, 32)
		draw_list.AddCircleFilled(hue_cursor_pos, hue_cursor_rad, hue_color32, hue_cursor_segments)
		draw_list.AddCircle(hue_cursor_pos, hue_cursor_rad+1, col_midgrey, hue_cursor_segments, 1)
		draw_list.AddCircle(hue_cursor_pos, hue_cursor_rad, col_white, hue_cursor_segments, 1)

		// Render SV triangle (rotated according to hue)
		var tra = wheel_center.Add(*ImRotate(&triangle_pa, cos_hue_angle, sin_hue_angle))
		var trb = wheel_center.Add(*ImRotate(&triangle_pb, cos_hue_angle, sin_hue_angle))
		var trc = wheel_center.Add(*ImRotate(&triangle_pc, cos_hue_angle, sin_hue_angle))
		var uv_white = GetFontTexUvWhitePixel()
		draw_list.PrimReserve(6, 6)
		draw_list.PrimVtx(tra, &uv_white, hue_color32)
		draw_list.PrimVtx(trb, &uv_white, hue_color32)
		draw_list.PrimVtx(trc, &uv_white, col_white)
		draw_list.PrimVtx(tra, &uv_white, 0)
		draw_list.PrimVtx(trb, &uv_white, col_black)
		draw_list.PrimVtx(trc, &uv_white, 0)
		draw_list.AddTriangle(&tra, &trb, trc, col_midgrey, 1.5)
		from := ImLerpVec2(&trc, &tra, ImSaturate(S))
		sv_cursor_pos = ImLerpVec2(&from, &trb, ImSaturate(1-V))
	} else if flags&ImGuiColorEditFlags_PickerHueBar != 0 {
		// Render SV Square
		draw_list.AddRectFilledMultiColor(picker_pos, picker_pos.Add(ImVec2{sv_picker_size, sv_picker_size}), col_white, hue_color32, hue_color32, col_white)
		draw_list.AddRectFilledMultiColor(picker_pos, picker_pos.Add(ImVec2{sv_picker_size, sv_picker_size}), 0, 0, col_black, col_black)
		RenderFrameBorder(picker_pos, picker_pos.Add(ImVec2{sv_picker_size, sv_picker_size}), 0.0)
		sv_cursor_pos.x = ImClamp(IM_ROUND(picker_pos.x+ImSaturate(S)*sv_picker_size), picker_pos.x+2, picker_pos.x+sv_picker_size-2) // Sneakily prevent the circle to stick out too much
		sv_cursor_pos.y = ImClamp(IM_ROUND(picker_pos.y+ImSaturate(1-V)*sv_picker_size), picker_pos.y+2, picker_pos.y+sv_picker_size-2)

		// Render Hue Bar
		for i := 0; i < 6; i++ {
			draw_list.AddRectFilledMultiColor(ImVec2{bar0_pos_x, picker_pos.y + float(i)*(sv_picker_size/6)}, ImVec2{bar0_pos_x + bars_width, picker_pos.y + float(i+1)*(sv_picker_size/6)}, col_hues[i], col_hues[i], col_hues[i+1], col_hues[i+1])
		}
		var bar0_line_y = IM_ROUND(picker_pos.y + H*sv_picker_size)
		RenderFrameBorder(ImVec2{bar0_pos_x, picker_pos.y}, ImVec2{bar0_pos_x + bars_width, picker_pos.y + sv_picker_size}, 0.0)
		RenderArrowsForVerticalBar(draw_list, ImVec2{bar0_pos_x - 1, bar0_line_y}, ImVec2{bars_triangles_half_sz + 1, bars_triangles_half_sz}, bars_width+2.0, style.Alpha)
	}

	// Render cursor/preview circle (clamp S/V within 0..1 range because floating points colors may lead HSV values to be out of range)
	var sv_cursor_rad float = 6.0
	if value_changed_sv {
		sv_cursor_rad = 10.0
	}
	draw_list.AddCircleFilled(sv_cursor_pos, sv_cursor_rad, user_col32_striped_of_alpha, 12)
	draw_list.AddCircle(sv_cursor_pos, sv_cursor_rad+1, col_midgrey, 12, 1)
	draw_list.AddCircle(sv_cursor_pos, sv_cursor_rad, col_white, 12, 1)

	// Render alpha bar
	if alpha_bar {
		var alpha = ImSaturate(col[3])
		var bar1_bb = ImRect{ImVec2{bar1_pos_x, picker_pos.y}, ImVec2{bar1_pos_x + bars_width, picker_pos.y + sv_picker_size}}
		RenderColorRectWithAlphaCheckerboard(draw_list, bar1_bb.Min, bar1_bb.Max, 0, bar1_bb.GetWidth()/2.0, ImVec2{}, 0, 0)
		draw_list.AddRectFilledMultiColor(bar1_bb.Min, bar1_bb.Max, user_col32_striped_of_alpha, user_col32_striped_of_alpha, user_col32_striped_of_alpha&^IM_COL32_A_MASK, user_col32_striped_of_alpha&^IM_COL32_A_MASK)
		var bar1_line_y = IM_ROUND(picker_pos.y + (1.0-alpha)*sv_picker_size)
		RenderFrameBorder(bar1_bb.Min, bar1_bb.Max, 0.0)
		RenderArrowsForVerticalBar(draw_list, ImVec2{bar1_pos_x - 1, bar1_line_y}, ImVec2{bars_triangles_half_sz + 1, bars_triangles_half_sz}, bars_width+2.0, style.Alpha)
	}

	EndGroup()

	if value_changed && backup_initial_col == [4]float{col[0], col[1], col[2], col[3]} {
		value_changed = false
	}
	if value_changed {
		MarkItemEdited(g.LastItemData.ID)
	}

	PopID()

	return value_changed
}

// display a color square/button, hover for details, return true when pressed.
// A little color square. Return true when clicked.
// FIXME: May want to display/ignore the alpha component in the color display? Yet show it in the tooltip.
// 'desc_id' is not called 'label' because we don't display it next to the button, but only in the tooltip.
// Note that 'col' may be encoded in HSV if ImGuiColorEditFlags_InputHSV is set.
func ColorButton(desc_id string, col ImVec4, flags ImGuiColorEditFlags, size ImVec2) bool {
	window := GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	g := GImGui
	var id = window.GetIDs(desc_id)
	var default_size = GetFrameHeight()
	if size.x == 0.0 {
		size.x = default_size
	}
	if size.y == 0.0 {
		size.y = default_size
	}
	var bb = ImRect{window.DC.CursorPos, window.DC.CursorPos.Add(size)}

	var padding float
	if size.y >= default_size {
		padding = g.Style.FramePadding.y
	}

	ItemSizeRect(&bb, padding)
	if !ItemAdd(&bb, id, nil, 0) {
		return false
	}

	var hovered, held bool
	var pressed = ButtonBehavior(&bb, id, &hovered, &held, 0)

	if flags&ImGuiColorEditFlags_NoAlpha != 0 {
		flags &= ^(ImGuiColorEditFlags_AlphaPreview | ImGuiColorEditFlags_AlphaPreviewHalf)
	}

	var col_rgb = col
	if flags&ImGuiColorEditFlags_InputHSV != 0 {
		ColorConvertHSVtoRGB(col_rgb.x, col_rgb.y, col_rgb.z, &col_rgb.x, &col_rgb.y, &col_rgb.z)
	}

	var col_rgb_without_alpha = ImVec4{col_rgb.x, col_rgb.y, col_rgb.z, 1.0}
	var grid_step = min(size.x, size.y) / 2.99
	var rounding = min(g.Style.FrameRounding, grid_step*0.5)
	var bb_inner = bb
	var off float
	if (flags & ImGuiColorEditFlags_NoBorder) == 0 {
		off = -0.75 // The border (using Col_FrameBg) tends to look off when color is near-opaque and rounding is enabled. This offset seemed like a good middle ground to reduce those artifacts.
		bb_inner.Expand(off)
	}
	if (flags&ImGuiColorEditFlags_AlphaPreviewHalf) != 0 && col_rgb.w < 1.0 {
		var mid_x = IM_ROUND((bb_inner.Min.x + bb_inner.Max.x) * 0.5)
		RenderColorRectWithAlphaCheckerboard(window.DrawList, ImVec2{bb_inner.Min.x + grid_step, bb_inner.Min.y}, bb_inner.Max, GetColorU32FromVec(col_rgb), grid_step, ImVec2{-grid_step + off, off}, rounding, ImDrawFlags_RoundCornersRight)
		window.DrawList.AddRectFilled(bb_inner.Min, ImVec2{mid_x, bb_inner.Max.y}, GetColorU32FromVec(col_rgb_without_alpha), rounding, ImDrawFlags_RoundCornersLeft)
	} else {
		// Because GetColorU32() multiplies by the global style Alpha and we don't want to display a checkerboard if the source code had no alpha
		var col_source = col_rgb_without_alpha
		if (flags & ImGuiColorEditFlags_AlphaPreview) != 0 {
			col_source = col_rgb
		}
		if col_source.w < 1.0 {
			RenderColorRectWithAlphaCheckerboard(window.DrawList, bb_inner.Min, bb_inner.Max, GetColorU32FromVec(col_source), grid_step, ImVec2{off, off}, rounding, 0)
		} else {
			window.DrawList.AddRectFilled(bb_inner.Min, bb_inner.Max, GetColorU32FromVec(col_source), rounding, 0)
		}
	}
	RenderNavHighlight(&bb, id, 0)
	if (flags & ImGuiColorEditFlags_NoBorder) == 0 {
		if g.Style.FrameBorderSize > 0.0 {
			RenderFrameBorder(bb.Min, bb.Max, rounding)
		} else {
			window.DrawList.AddRect(bb.Min, bb.Max, GetColorU32FromID(ImGuiCol_FrameBg, 1), rounding, 0, 1) // Color button are often in need of some sort of border
		}
	}

	// Drag and Drop Source
	// NB: The ActiveId test is merely an optional micro-optimization, BeginDragDropSource() does the same test.
	if g.ActiveId == id && (flags&ImGuiColorEditFlags_NoDragDrop) == 0 && BeginDragDropSource(0) {
		if flags&ImGuiColorEditFlags_NoAlpha != 0 {
			SetDragDropPayload(IMGUI_PAYLOAD_TYPE_COLOR_3F, &col_rgb, unsafe.Sizeof(float(0))*3, ImGuiCond_Once)
		} else {
			SetDragDropPayload(IMGUI_PAYLOAD_TYPE_COLOR_4F, &col_rgb, unsafe.Sizeof(float(0))*4, ImGuiCond_Once)
		}
		ColorButton(desc_id, col, flags, ImVec2{})
		SameLine(0, 0)
		TextEx("Color", 0)
		EndDragDropSource()
	}

	// Tooltip
	if flags&ImGuiColorEditFlags_NoTooltip == 0 && hovered {
		ColorTooltip(desc_id, [4]float32{col.x, col.y, col.z, col.w}, flags&(ImGuiColorEditFlags_InputMask_|ImGuiColorEditFlags_NoAlpha|ImGuiColorEditFlags_AlphaPreview|ImGuiColorEditFlags_AlphaPreviewHalf))
	}

	return pressed
}

// initialize current options (generally on application startup) if you want to select a default format, picker type, etc.
// User will be able to change many settings, unless you pass the _NoOptions flag to your calls.
func SetColorEditOptions(flags ImGuiColorEditFlags) {
	g := GImGui
	if (flags & ImGuiColorEditFlags_DisplayMask_) == 0 {
		flags |= ImGuiColorEditFlags_DefaultOptions_ & ImGuiColorEditFlags_DisplayMask_
	}
	if (flags & ImGuiColorEditFlags_DataTypeMask_) == 0 {
		flags |= ImGuiColorEditFlags_DefaultOptions_ & ImGuiColorEditFlags_DataTypeMask_
	}
	if (flags & ImGuiColorEditFlags_PickerMask_) == 0 {
		flags |= ImGuiColorEditFlags_DefaultOptions_ & ImGuiColorEditFlags_PickerMask_
	}
	if (flags & ImGuiColorEditFlags_InputMask_) == 0 {
		flags |= ImGuiColorEditFlags_DefaultOptions_ & ImGuiColorEditFlags_InputMask_
	}
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_DisplayMask_)))  // Check only 1 option is selected
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_DataTypeMask_))) // Check only 1 option is selected
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_PickerMask_)))   // Check only 1 option is selected
	IM_ASSERT(ImIsPowerOfTwoInt(int(flags & ImGuiColorEditFlags_InputMask_)))    // Check only 1 option is selected
	g.ColorEditOptions = flags
}

// Color
// Note: only access 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
func ColorTooltip(text string, col [4]float, flags ImGuiColorEditFlags) {
	g := GImGui

	BeginTooltipEx(0, ImGuiTooltipFlags_OverridePreviousTooltip)
	if len(text) > 0 {
		TextEx(text, 0)
		Separator()
	}

	var a = col[3]
	if flags&ImGuiColorEditFlags_NoAlpha != 0 {
		a = 1
	}

	var sz = ImVec2{g.FontSize*3 + g.Style.FramePadding.y*2, g.FontSize*3 + g.Style.FramePadding.y*2}
	var cf = ImVec4{col[0], col[1], col[2], a}
	var cr, cg, cb, ca = IM_F32_TO_INT8_SAT(col[0]), IM_F32_TO_INT8_SAT(col[1]), IM_F32_TO_INT8_SAT(col[2]), IM_F32_TO_INT8_SAT(col[3])
	if (flags & ImGuiColorEditFlags_NoAlpha) != 0 {
		ca = 255
	}
	ColorButton("##preview", cf, (flags&(ImGuiColorEditFlags_InputMask_|ImGuiColorEditFlags_NoAlpha|ImGuiColorEditFlags_AlphaPreview|ImGuiColorEditFlags_AlphaPreviewHalf))|ImGuiColorEditFlags_NoTooltip, sz)
	SameLine(0, 0)
	if (flags&ImGuiColorEditFlags_InputRGB) != 0 || (flags&ImGuiColorEditFlags_InputMask_) == 0 {
		if (flags & ImGuiColorEditFlags_NoAlpha) != 0 {
			Text("#%02X%02X%02X\nR: %d, G: %d, B: %d\n(%.3f, %.3f, %.3f)", cr, cg, cb, cr, cg, cb, col[0], col[1], col[2])
		} else {
			Text("#%02X%02X%02X%02X\nR:%d, G:%d, B:%d, A:%d\n(%.3f, %.3f, %.3f, %.3f)", cr, cg, cb, ca, cr, cg, cb, ca, col[0], col[1], col[2], col[3])
		}
	} else if flags&ImGuiColorEditFlags_InputHSV != 0 {
		if flags&ImGuiColorEditFlags_NoAlpha != 0 {
			Text("H: %.3f, S: %.3f, V: %.3f", col[0], col[1], col[2])
		} else {
			Text("H: %.3f, S: %.3f, V: %.3f, A: %.3f", col[0], col[1], col[2], col[3])
		}
	}
	EndTooltip()
}

func ColorEditOptionsPopup(col [4]float, flags ImGuiColorEditFlags) {
	var allow_opt_inputs = (flags & ImGuiColorEditFlags_DisplayMask_) == 0
	var allow_opt_datatype = (flags & ImGuiColorEditFlags_DataTypeMask_) == 0
	if (!allow_opt_inputs && !allow_opt_datatype) || !BeginPopup("context", 0) {
		return
	}
	g := GImGui
	var opts = g.ColorEditOptions
	if allow_opt_inputs {
		if RadioButtonBool("RGB", (opts&ImGuiColorEditFlags_DisplayRGB) != 0) {
			opts = (opts & ^ImGuiColorEditFlags_DisplayMask_) | ImGuiColorEditFlags_DisplayRGB
		}
		if RadioButtonBool("HSV", (opts&ImGuiColorEditFlags_DisplayHSV) != 0) {
			opts = (opts & ^ImGuiColorEditFlags_DisplayMask_) | ImGuiColorEditFlags_DisplayHSV
		}
		if RadioButtonBool("Hex", (opts&ImGuiColorEditFlags_DisplayHex) != 0) {
			opts = (opts & ^ImGuiColorEditFlags_DisplayMask_) | ImGuiColorEditFlags_DisplayHex
		}
	}
	if allow_opt_datatype {
		if allow_opt_inputs {
			Separator()
		}
		if RadioButtonBool("0..255", (opts&ImGuiColorEditFlags_Uint8) != 0) {
			opts = (opts & ^ImGuiColorEditFlags_DataTypeMask_) | ImGuiColorEditFlags_Uint8
		}
		if RadioButtonBool("0.00..1.00", (opts&ImGuiColorEditFlags_Float) != 0) {
			opts = (opts & ^ImGuiColorEditFlags_DataTypeMask_) | ImGuiColorEditFlags_Float
		}
	}

	if allow_opt_inputs || allow_opt_datatype {
		Separator()
	}

	if (ButtonEx("Copy as..", &ImVec2{-1, 0}, 0)) {
		OpenPopup("Copy", 0)
	}

	if BeginPopup("Copy", 0) {
		var cr, cg, cb, ca = IM_F32_TO_INT8_SAT(col[0]), IM_F32_TO_INT8_SAT(col[1]), IM_F32_TO_INT8_SAT(col[2]), IM_F32_TO_INT8_SAT(col[3])
		if flags&ImGuiColorEditFlags_NoAlpha != 0 {
			ca = 255
		}

		a := col[3]
		if flags&ImGuiColorEditFlags_NoAlpha != 0 {
			a = 1
		}

		var buf = fmt.Sprintf("%v,%v,%v,%v", col[0], col[1], col[2], a)
		if Selectable(buf, false, 0, ImVec2{}) {
			SetClipboardText(buf)
		}
		buf = fmt.Sprintf("%v,%v,%v,%v", cr, cg, cb, ca)
		if Selectable(buf, false, 0, ImVec2{}) {
			SetClipboardText(buf)
		}
		buf = fmt.Sprintf("%v,%v,%v", cr, cg, cb)
		if Selectable(buf, false, 0, ImVec2{}) {
			SetClipboardText(buf)
		}
		if (flags & ImGuiColorEditFlags_NoAlpha) == 0 {
			buf = fmt.Sprintf("%v,%v,%v,%v", cr, cg, cb, ca)
			if Selectable(buf, false, 0, ImVec2{}) {
				SetClipboardText(buf)
			}
		}
		EndPopup()
	}

	g.ColorEditOptions = opts
	EndPopup()
}

func ColorPickerOptionsPopup(ref_col *[4]float, flags ImGuiColorEditFlags) {
	var allow_opt_picker = flags&ImGuiColorEditFlags_PickerMask_ == 0
	var allow_opt_alpha_bar = flags&ImGuiColorEditFlags_NoAlpha == 0 && flags&ImGuiColorEditFlags_AlphaBar == 0
	if (!allow_opt_picker && !allow_opt_alpha_bar) || !BeginPopup("context", 0) {
		return
	}
	g := GImGui
	if allow_opt_picker {
		var picker_size = ImVec2{g.FontSize * 8, max(g.FontSize*8-(GetFrameHeight()+g.Style.ItemInnerSpacing.x), 1.0)} // FIXME: Picker size copied from main picker function
		PushItemWidth(picker_size.x)
		for picker_type := 0; picker_type < 2; picker_type++ {
			// Draw small/thumbnail version of each picker type (over an invisible button for selection)
			if picker_type > 0 {
				Separator()
			}
			PushInterface(picker_type)
			var picker_flags = ImGuiColorEditFlags_NoInputs | ImGuiColorEditFlags_NoOptions | ImGuiColorEditFlags_NoLabel | ImGuiColorEditFlags_NoSidePreview | (flags & ImGuiColorEditFlags_NoAlpha)
			if picker_type == 0 {
				picker_flags |= ImGuiColorEditFlags_PickerHueBar
			}
			if picker_type == 1 {
				picker_flags |= ImGuiColorEditFlags_PickerHueWheel
			}
			var backup_pos = GetCursorScreenPos()
			if Selectable("##selectable", false, 0, picker_size) { // By default, Selectable() is closing popup
				g.ColorEditOptions = (g.ColorEditOptions & ^ImGuiColorEditFlags_PickerMask_) | (picker_flags & ImGuiColorEditFlags_PickerMask_)
			}
			SetCursorScreenPos(backup_pos)
			var previewing_ref_col [4]float32

			var components int = 4
			if (picker_flags & ImGuiColorEditFlags_NoAlpha) != 0 {
				components = 3
			}

			copy(previewing_ref_col[:], ref_col[:components])
			ColorPicker4("##previewing_picker", &previewing_ref_col, picker_flags, nil)
			PopID()
		}
		PopItemWidth()
	}
	if allow_opt_alpha_bar {
		if allow_opt_picker {
			Separator()
		}
		CheckboxFlagsInt("Alpha Bar", (*int)(&g.ColorEditOptions), int(ImGuiColorEditFlags_AlphaBar))
	}
	EndPopup()
}
