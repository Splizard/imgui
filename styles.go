package imgui

func StyleColorsDark(dst *ImGuiStyle) {
	var style *ImGuiStyle = dst
	if dst == nil {
		style = GetStyle()
	}
	var colors = &style.Colors

	colors[ImGuiCol_Text] = ImVec4{1.00, 1.00, 1.00, 1.00}
	colors[ImGuiCol_TextDisabled] = ImVec4{0.50, 0.50, 0.50, 1.00}
	colors[ImGuiCol_WindowBg] = ImVec4{0.06, 0.06, 0.06, 0.94}
	colors[ImGuiCol_ChildBg] = ImVec4{0.00, 0.00, 0.00, 0.00}
	colors[ImGuiCol_PopupBg] = ImVec4{0.08, 0.08, 0.08, 0.94}
	colors[ImGuiCol_Border] = ImVec4{0.43, 0.43, 0.50, 0.50}
	colors[ImGuiCol_BorderShadow] = ImVec4{0.00, 0.00, 0.00, 0.00}
	colors[ImGuiCol_FrameBg] = ImVec4{0.16, 0.29, 0.48, 0.54}
	colors[ImGuiCol_FrameBgHovered] = ImVec4{0.26, 0.59, 0.98, 0.40}
	colors[ImGuiCol_FrameBgActive] = ImVec4{0.26, 0.59, 0.98, 0.67}
	colors[ImGuiCol_TitleBg] = ImVec4{0.04, 0.04, 0.04, 1.00}
	colors[ImGuiCol_TitleBgActive] = ImVec4{0.16, 0.29, 0.48, 1.00}
	colors[ImGuiCol_TitleBgCollapsed] = ImVec4{0.00, 0.00, 0.00, 0.51}
	colors[ImGuiCol_MenuBarBg] = ImVec4{0.14, 0.14, 0.14, 1.00}
	colors[ImGuiCol_ScrollbarBg] = ImVec4{0.02, 0.02, 0.02, 0.53}
	colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.31, 0.31, 0.31, 1.00}
	colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.41, 0.41, 0.41, 1.00}
	colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.51, 0.51, 0.51, 1.00}
	colors[ImGuiCol_CheckMark] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_SliderGrab] = ImVec4{0.24, 0.52, 0.88, 1.00}
	colors[ImGuiCol_SliderGrabActive] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_Button] = ImVec4{0.26, 0.59, 0.98, 0.40}
	colors[ImGuiCol_ButtonHovered] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_ButtonActive] = ImVec4{0.06, 0.53, 0.98, 1.00}
	colors[ImGuiCol_Header] = ImVec4{0.26, 0.59, 0.98, 0.31}
	colors[ImGuiCol_HeaderHovered] = ImVec4{0.26, 0.59, 0.98, 0.80}
	colors[ImGuiCol_HeaderActive] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_Separator] = colors[ImGuiCol_Border]
	colors[ImGuiCol_SeparatorHovered] = ImVec4{0.10, 0.40, 0.75, 0.78}
	colors[ImGuiCol_SeparatorActive] = ImVec4{0.10, 0.40, 0.75, 1.00}
	colors[ImGuiCol_ResizeGrip] = ImVec4{0.26, 0.59, 0.98, 0.20}
	colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.26, 0.59, 0.98, 0.67}
	colors[ImGuiCol_ResizeGripActive] = ImVec4{0.26, 0.59, 0.98, 0.95}
	colors[ImGuiCol_Tab] = ImLerpVec4(&colors[ImGuiCol_Header], &colors[ImGuiCol_TitleBgActive], 0.80)
	colors[ImGuiCol_TabHovered] = colors[ImGuiCol_HeaderHovered]
	colors[ImGuiCol_TabActive] = ImLerpVec4(&colors[ImGuiCol_HeaderActive], &colors[ImGuiCol_TitleBgActive], 0.60)
	colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&colors[ImGuiCol_Tab], &colors[ImGuiCol_TitleBg], 0.80)
	colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&colors[ImGuiCol_TabActive], &colors[ImGuiCol_TitleBg], 0.40)
	colors[ImGuiCol_PlotLines] = ImVec4{0.61, 0.61, 0.61, 1.00}
	colors[ImGuiCol_PlotLinesHovered] = ImVec4{1.00, 0.43, 0.35, 1.00}
	colors[ImGuiCol_PlotHistogram] = ImVec4{0.90, 0.70, 0.00, 1.00}
	colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.00, 0.60, 0.00, 1.00}
	colors[ImGuiCol_TableHeaderBg] = ImVec4{0.19, 0.19, 0.20, 1.00}
	colors[ImGuiCol_TableBorderStrong] = ImVec4{0.31, 0.31, 0.35, 1.00} // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableBorderLight] = ImVec4{0.23, 0.23, 0.25, 1.00}  // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableRowBg] = ImVec4{0.00, 0.00, 0.00, 0.00}
	colors[ImGuiCol_TableRowBgAlt] = ImVec4{1.00, 1.00, 1.00, 0.06}
	colors[ImGuiCol_TextSelectedBg] = ImVec4{0.26, 0.59, 0.98, 0.35}
	colors[ImGuiCol_DragDropTarget] = ImVec4{1.00, 1.00, 0.00, 0.90}
	colors[ImGuiCol_NavHighlight] = ImVec4{0.26, 0.59, 0.98, 1.00}
	colors[ImGuiCol_NavWindowingHighlight] = ImVec4{1.00, 1.00, 1.00, 0.70}
	colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.80, 0.80, 0.80, 0.20}
	colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.80, 0.80, 0.80, 0.35}
}

func StyleColorsClassic(dst *ImGuiStyle) {
	var style = dst
	if style == nil {
		style = GetStyle()
	}
	var colors = style.Colors

	colors[ImGuiCol_Text] = ImVec4{0.9, 0.9, 0.9, 1.0}
	colors[ImGuiCol_TextDisabled] = ImVec4{0.6, 0.6, 0.6, 1.0}
	colors[ImGuiCol_WindowBg] = ImVec4{0.0, 0.0, 0.0, 0.8}
	colors[ImGuiCol_ChildBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	colors[ImGuiCol_PopupBg] = ImVec4{0.1, 0.1, 0.1, 0.9}
	colors[ImGuiCol_Border] = ImVec4{0.5, 0.5, 0.5, 0.5}
	colors[ImGuiCol_BorderShadow] = ImVec4{0.0, 0.0, 0.0, 0.0}
	colors[ImGuiCol_FrameBg] = ImVec4{0.4, 0.4, 0.4, 0.3}
	colors[ImGuiCol_FrameBgHovered] = ImVec4{0.4, 0.4, 0.6, 0.4}
	colors[ImGuiCol_FrameBgActive] = ImVec4{0.4, 0.4, 0.6, 0.6}
	colors[ImGuiCol_TitleBg] = ImVec4{0.2, 0.2, 0.5, 0.8}
	colors[ImGuiCol_TitleBgActive] = ImVec4{0.3, 0.3, 0.6, 0.8}
	colors[ImGuiCol_TitleBgCollapsed] = ImVec4{0.4, 0.4, 0.8, 0.2}
	colors[ImGuiCol_MenuBarBg] = ImVec4{0.4, 0.4, 0.5, 0.8}
	colors[ImGuiCol_ScrollbarBg] = ImVec4{0.2, 0.2, 0.3, 0.6}
	colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.4, 0.4, 0.8, 0.3}
	colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.4, 0.4, 0.8, 0.4}
	colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.4, 0.3, 0.8, 0.6}
	colors[ImGuiCol_CheckMark] = ImVec4{0.9, 0.9, 0.9, 0.5}
	colors[ImGuiCol_SliderGrab] = ImVec4{1.0, 1.0, 1.0, 0.3}
	colors[ImGuiCol_SliderGrabActive] = ImVec4{0.4, 0.3, 0.8, 0.6}
	colors[ImGuiCol_Button] = ImVec4{0.3, 0.4, 0.6, 0.6}
	colors[ImGuiCol_ButtonHovered] = ImVec4{0.4, 0.4, 0.7, 0.7}
	colors[ImGuiCol_ButtonActive] = ImVec4{0.4, 0.5, 0.8, 1.0}
	colors[ImGuiCol_Header] = ImVec4{0.4, 0.4, 0.9, 0.4}
	colors[ImGuiCol_HeaderHovered] = ImVec4{0.4, 0.4, 0.9, 0.8}
	colors[ImGuiCol_HeaderActive] = ImVec4{0.5, 0.5, 0.8, 0.8}
	colors[ImGuiCol_Separator] = ImVec4{0.5, 0.5, 0.5, 0.6}
	colors[ImGuiCol_SeparatorHovered] = ImVec4{0.6, 0.6, 0.7, 1.0}
	colors[ImGuiCol_SeparatorActive] = ImVec4{0.7, 0.7, 0.9, 1.0}
	colors[ImGuiCol_ResizeGrip] = ImVec4{1.0, 1.0, 1.0, 0.1}
	colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.7, 0.8, 1.0, 0.6}
	colors[ImGuiCol_ResizeGripActive] = ImVec4{0.7, 0.8, 1.0, 0.9}
	colors[ImGuiCol_Tab] = ImLerpVec4(&colors[ImGuiCol_Header], &colors[ImGuiCol_TitleBgActive], 0.8)
	colors[ImGuiCol_TabHovered] = colors[ImGuiCol_HeaderHovered]
	colors[ImGuiCol_TabActive] = ImLerpVec4(&colors[ImGuiCol_HeaderActive], &colors[ImGuiCol_TitleBgActive], 0.6)
	colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&colors[ImGuiCol_Tab], &colors[ImGuiCol_TitleBg], 0.8)
	colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&colors[ImGuiCol_TabActive], &colors[ImGuiCol_TitleBg], 0.4)
	colors[ImGuiCol_PlotLines] = ImVec4{1.0, 1.0, 1.0, 1.0}
	colors[ImGuiCol_PlotLinesHovered] = ImVec4{0.9, 0.7, 0.0, 1.0}
	colors[ImGuiCol_PlotHistogram] = ImVec4{0.9, 0.7, 0.0, 1.0}
	colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.0, 0.6, 0.0, 1.0}
	colors[ImGuiCol_TableHeaderBg] = ImVec4{0.2, 0.2, 0.3, 1.0}
	colors[ImGuiCol_TableBorderStrong] = ImVec4{0.3, 0.3, 0.4, 1.0} // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableBorderLight] = ImVec4{0.2, 0.2, 0.2, 1.0}  // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableRowBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	colors[ImGuiCol_TableRowBgAlt] = ImVec4{1.0, 1.0, 1.0, 0.0}
	colors[ImGuiCol_TextSelectedBg] = ImVec4{0.0, 0.0, 1.0, 0.3}
	colors[ImGuiCol_DragDropTarget] = ImVec4{1.0, 1.0, 0.0, 0.9}
	colors[ImGuiCol_NavHighlight] = colors[ImGuiCol_HeaderHovered]
	colors[ImGuiCol_NavWindowingHighlight] = ImVec4{1.0, 1.0, 1.0, 0.7}
	colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.8, 0.8, 0.8, 0.2}
	colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.2, 0.2, 0.2, 0.3}
}

// Those light colors are better suited with a thicker font than the default one + FrameBorder
func StyleColorsLight(dst *ImGuiStyle) {
	var style = dst
	if style == nil {
		style = GetStyle()
	}
	var colors = style.Colors

	colors[ImGuiCol_Text] = ImVec4{0.0, 0.0, 0.0, 1.0}
	colors[ImGuiCol_TextDisabled] = ImVec4{0.6, 0.6, 0.6, 1.0}
	colors[ImGuiCol_WindowBg] = ImVec4{0.9, 0.9, 0.9, 1.0}
	colors[ImGuiCol_ChildBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	colors[ImGuiCol_PopupBg] = ImVec4{1.0, 1.0, 1.0, 0.9}
	colors[ImGuiCol_Border] = ImVec4{0.0, 0.0, 0.0, 0.3}
	colors[ImGuiCol_BorderShadow] = ImVec4{0.0, 0.0, 0.0, 0.0}
	colors[ImGuiCol_FrameBg] = ImVec4{1.0, 1.0, 1.0, 1.0}
	colors[ImGuiCol_FrameBgHovered] = ImVec4{0.2, 0.5, 0.9, 0.4}
	colors[ImGuiCol_FrameBgActive] = ImVec4{0.2, 0.5, 0.9, 0.6}
	colors[ImGuiCol_TitleBg] = ImVec4{0.9, 0.9, 0.9, 1.0}
	colors[ImGuiCol_TitleBgActive] = ImVec4{0.8, 0.8, 0.8, 1.0}
	colors[ImGuiCol_TitleBgCollapsed] = ImVec4{1.0, 1.0, 1.0, 0.5}
	colors[ImGuiCol_MenuBarBg] = ImVec4{0.8, 0.8, 0.8, 1.0}
	colors[ImGuiCol_ScrollbarBg] = ImVec4{0.9, 0.9, 0.9, 0.5}
	colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.6, 0.6, 0.6, 0.8}
	colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.4, 0.4, 0.4, 0.8}
	colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.4, 0.4, 0.4, 1.0}
	colors[ImGuiCol_CheckMark] = ImVec4{0.2, 0.5, 0.9, 1.0}
	colors[ImGuiCol_SliderGrab] = ImVec4{0.2, 0.5, 0.9, 0.7}
	colors[ImGuiCol_SliderGrabActive] = ImVec4{0.4, 0.5, 0.8, 0.6}
	colors[ImGuiCol_Button] = ImVec4{0.2, 0.5, 0.9, 0.4}
	colors[ImGuiCol_ButtonHovered] = ImVec4{0.2, 0.5, 0.9, 1.0}
	colors[ImGuiCol_ButtonActive] = ImVec4{0.0, 0.5, 0.9, 1.0}
	colors[ImGuiCol_Header] = ImVec4{0.2, 0.5, 0.9, 0.3}
	colors[ImGuiCol_HeaderHovered] = ImVec4{0.2, 0.5, 0.9, 0.8}
	colors[ImGuiCol_HeaderActive] = ImVec4{0.2, 0.5, 0.9, 1.0}
	colors[ImGuiCol_Separator] = ImVec4{0.3, 0.3, 0.3, 0.6}
	colors[ImGuiCol_SeparatorHovered] = ImVec4{0.1, 0.4, 0.8, 0.7}
	colors[ImGuiCol_SeparatorActive] = ImVec4{0.1, 0.4, 0.8, 1.0}
	colors[ImGuiCol_ResizeGrip] = ImVec4{0.3, 0.3, 0.3, 0.1}
	colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.2, 0.5, 0.9, 0.6}
	colors[ImGuiCol_ResizeGripActive] = ImVec4{0.2, 0.5, 0.9, 0.9}
	colors[ImGuiCol_Tab] = ImLerpVec4(&colors[ImGuiCol_Header], &colors[ImGuiCol_TitleBgActive], 0.9)
	colors[ImGuiCol_TabHovered] = colors[ImGuiCol_HeaderHovered]
	colors[ImGuiCol_TabActive] = ImLerpVec4(&colors[ImGuiCol_HeaderActive], &colors[ImGuiCol_TitleBgActive], 0.6)
	colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&colors[ImGuiCol_Tab], &colors[ImGuiCol_TitleBg], 0.8)
	colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&colors[ImGuiCol_TabActive], &colors[ImGuiCol_TitleBg], 0.4)
	colors[ImGuiCol_PlotLines] = ImVec4{0.3, 0.3, 0.3, 1.0}
	colors[ImGuiCol_PlotLinesHovered] = ImVec4{1.0, 0.4, 0.3, 1.0}
	colors[ImGuiCol_PlotHistogram] = ImVec4{0.9, 0.7, 0.0, 1.0}
	colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.0, 0.4, 0.0, 1.0}
	colors[ImGuiCol_TableHeaderBg] = ImVec4{0.7, 0.8, 0.9, 1.0}
	colors[ImGuiCol_TableBorderStrong] = ImVec4{0.5, 0.5, 0.6, 1.0} // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableBorderLight] = ImVec4{0.6, 0.6, 0.7, 1.0}  // Prefer using Alpha=1.0 here
	colors[ImGuiCol_TableRowBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	colors[ImGuiCol_TableRowBgAlt] = ImVec4{0.3, 0.3, 0.3, 0.0}
	colors[ImGuiCol_TextSelectedBg] = ImVec4{0.2, 0.5, 0.9, 0.3}
	colors[ImGuiCol_DragDropTarget] = ImVec4{0.2, 0.5, 0.9, 0.9}
	colors[ImGuiCol_NavHighlight] = colors[ImGuiCol_HeaderHovered]
	colors[ImGuiCol_NavWindowingHighlight] = ImVec4{0.7, 0.7, 0.7, 0.7}
	colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.2, 0.2, 0.2, 0.2}
	colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.2, 0.2, 0.2, 0.3}
}
