package imgui

func ColorConvertFloat4ToU32(in ImVec4) ImU32 {
	var out ImU32
	out = ((ImU32)(IM_F32_TO_INT8_SAT(in.x))) << IM_COL32_R_SHIFT
	out |= ((ImU32)(IM_F32_TO_INT8_SAT(in.y))) << IM_COL32_G_SHIFT
	out |= ((ImU32)(IM_F32_TO_INT8_SAT(in.z))) << IM_COL32_B_SHIFT
	out |= ((ImU32)(IM_F32_TO_INT8_SAT(in.w))) << IM_COL32_A_SHIFT
	return out
}
