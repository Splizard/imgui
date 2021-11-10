package imgui

// Shade functions (write over already created vertices)

// Generic linear color gradient, write to RGB fields, leave A untouched.
func ShadeVertsLinearColorGradientKeepAlpha(draw_list *ImDrawList, vert_start_idx int, vert_end_idx int, gradient_p0 ImVec2, gradient_p1 ImVec2, col0 ImU32, col1 ImU32) {
	var gradient_extent ImVec2 = gradient_p1.Sub(gradient_p0)
	var gradient_inv_length2 float = 1.0 / ImLengthSqrVec2(gradient_extent)
	var col0_r = (int)(col0>>IM_COL32_R_SHIFT) & 0xFF
	var col0_g = (int)(col0>>IM_COL32_G_SHIFT) & 0xFF
	var col0_b = (int)(col0>>IM_COL32_B_SHIFT) & 0xFF
	var col_delta_r = ((int)(col1>>IM_COL32_R_SHIFT) & 0xFF) - col0_r
	var col_delta_g = ((int)(col1>>IM_COL32_G_SHIFT) & 0xFF) - col0_g
	var col_delta_b = ((int)(col1>>IM_COL32_B_SHIFT) & 0xFF) - col0_b
	for vert_idx := vert_start_idx; vert_idx < vert_end_idx; vert_idx++ {
		vert := draw_list.VtxBuffer[vert_idx]
		var diff = vert.pos.Sub(gradient_p0)
		var d = ImDot(&diff, &gradient_extent)
		var t = ImClamp(d*gradient_inv_length2, 0.0, 1.0)
		var r = (int)(float(col0_r) + float(col_delta_r)*t)
		var g = (int)(float(col0_g) + float(col_delta_g)*t)
		var b = (int)(float(col0_b) + float(col_delta_b)*t)
		vert.col = (uint(r) << IM_COL32_R_SHIFT) | (uint(g) << IM_COL32_G_SHIFT) | (uint(b) << IM_COL32_B_SHIFT) | (vert.col & IM_COL32_A_MASK)
	}
}

// Distribute UV over (a, b) rectangle
func ShadeVertsLinearUV(t *ImDrawList, vert_start_idx int, vert_end_idx int, a *ImVec2, b *ImVec2, uv_a *ImVec2, uv_b *ImVec2, clamp bool) {
	var size = b.Sub(*a)
	var uv_size = uv_b.Sub(*uv_a)

	var scale ImVec2
	if size.x != 0.0 {
		scale.x = uv_size.x / size.x
	}
	if size.y != 0.0 {
		scale.y = uv_size.y / size.y
	}

	if clamp {
		var min = ImMinVec2(uv_a, uv_b)
		var max = ImMinVec2(uv_a, uv_b)
		for vertex_idx := vert_start_idx; vertex_idx < vert_end_idx; vertex_idx++ {
			vertex := &t.VtxBuffer[vertex_idx]
			d := ImVec2{vertex.pos.x, vertex.pos.y}.Sub(*a)
			v := uv_a.Add(*ImMul(&d, &scale))
			vertex.uv = ImClampVec2(&v, &min, max)
		}
	} else {
		for vertex_idx := vert_start_idx; vertex_idx < vert_end_idx; vertex_idx++ {
			vertex := &t.VtxBuffer[vertex_idx]
			v := ImVec2{vertex.pos.x, vertex.pos.y}.Sub(*a)
			vertex.uv = uv_a.Add(*ImMul(&v, &scale))
		}
	}
}
