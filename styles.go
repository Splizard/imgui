package imgui

func StyleColorsDark(style *ImGuiStyle) {
	if style == nil {
		style = GetStyle()
	}
	style.Colors[ImGuiCol_Text] = ImVec4{1.00, 1.00, 1.00, 1.00}
	style.Colors[ImGuiCol_TextDisabled] = ImVec4{0.50, 0.50, 0.50, 1.00}
	style.Colors[ImGuiCol_WindowBg] = ImVec4{0.06, 0.06, 0.06, 0.94}
	style.Colors[ImGuiCol_ChildBg] = ImVec4{0.00, 0.00, 0.00, 0.00}
	style.Colors[ImGuiCol_PopupBg] = ImVec4{0.08, 0.08, 0.08, 0.94}
	style.Colors[ImGuiCol_Border] = ImVec4{0.43, 0.43, 0.50, 0.50}
	style.Colors[ImGuiCol_BorderShadow] = ImVec4{0.00, 0.00, 0.00, 0.00}
	style.Colors[ImGuiCol_FrameBg] = ImVec4{0.16, 0.29, 0.48, 0.54}
	style.Colors[ImGuiCol_FrameBgHovered] = ImVec4{0.26, 0.59, 0.98, 0.40}
	style.Colors[ImGuiCol_FrameBgActive] = ImVec4{0.26, 0.59, 0.98, 0.67}
	style.Colors[ImGuiCol_TitleBg] = ImVec4{0.04, 0.04, 0.04, 1.00}
	style.Colors[ImGuiCol_TitleBgActive] = ImVec4{0.16, 0.29, 0.48, 1.00}
	style.Colors[ImGuiCol_TitleBgCollapsed] = ImVec4{0.00, 0.00, 0.00, 0.51}
	style.Colors[ImGuiCol_MenuBarBg] = ImVec4{0.14, 0.14, 0.14, 1.00}
	style.Colors[ImGuiCol_ScrollbarBg] = ImVec4{0.02, 0.02, 0.02, 0.53}
	style.Colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.31, 0.31, 0.31, 1.00}
	style.Colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.41, 0.41, 0.41, 1.00}
	style.Colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.51, 0.51, 0.51, 1.00}
	style.Colors[ImGuiCol_CheckMark] = ImVec4{0.26, 0.59, 0.98, 1.00}
	style.Colors[ImGuiCol_SliderGrab] = ImVec4{0.24, 0.52, 0.88, 1.00}
	style.Colors[ImGuiCol_SliderGrabActive] = ImVec4{0.26, 0.59, 0.98, 1.00}
	style.Colors[ImGuiCol_Button] = ImVec4{0.26, 0.59, 0.98, 0.40}
	style.Colors[ImGuiCol_ButtonHovered] = ImVec4{0.26, 0.59, 0.98, 1.00}
	style.Colors[ImGuiCol_ButtonActive] = ImVec4{0.06, 0.53, 0.98, 1.00}
	style.Colors[ImGuiCol_Header] = ImVec4{0.26, 0.59, 0.98, 0.31}
	style.Colors[ImGuiCol_HeaderHovered] = ImVec4{0.26, 0.59, 0.98, 0.80}
	style.Colors[ImGuiCol_HeaderActive] = ImVec4{0.26, 0.59, 0.98, 1.00}
	style.Colors[ImGuiCol_Separator] = style.Colors[ImGuiCol_Border]
	style.Colors[ImGuiCol_SeparatorHovered] = ImVec4{0.10, 0.40, 0.75, 0.78}
	style.Colors[ImGuiCol_SeparatorActive] = ImVec4{0.10, 0.40, 0.75, 1.00}
	style.Colors[ImGuiCol_ResizeGrip] = ImVec4{0.26, 0.59, 0.98, 0.20}
	style.Colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.26, 0.59, 0.98, 0.67}
	style.Colors[ImGuiCol_ResizeGripActive] = ImVec4{0.26, 0.59, 0.98, 0.95}
	style.Colors[ImGuiCol_Tab] = ImLerpVec4(&style.Colors[ImGuiCol_Header], &style.Colors[ImGuiCol_TitleBgActive], 0.80)
	style.Colors[ImGuiCol_TabHovered] = style.Colors[ImGuiCol_HeaderHovered]
	style.Colors[ImGuiCol_TabActive] = ImLerpVec4(&style.Colors[ImGuiCol_HeaderActive], &style.Colors[ImGuiCol_TitleBgActive], 0.60)
	style.Colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&style.Colors[ImGuiCol_Tab], &style.Colors[ImGuiCol_TitleBg], 0.80)
	style.Colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&style.Colors[ImGuiCol_TabActive], &style.Colors[ImGuiCol_TitleBg], 0.40)
	style.Colors[ImGuiCol_PlotLines] = ImVec4{0.61, 0.61, 0.61, 1.00}
	style.Colors[ImGuiCol_PlotLinesHovered] = ImVec4{1.00, 0.43, 0.35, 1.00}
	style.Colors[ImGuiCol_PlotHistogram] = ImVec4{0.90, 0.70, 0.00, 1.00}
	style.Colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.00, 0.60, 0.00, 1.00}
	style.Colors[ImGuiCol_TableHeaderBg] = ImVec4{0.19, 0.19, 0.20, 1.00}
	style.Colors[ImGuiCol_TableBorderStrong] = ImVec4{0.31, 0.31, 0.35, 1.00} // Prefer using Alpha=1.0 here
	style.Colors[ImGuiCol_TableBorderLight] = ImVec4{0.23, 0.23, 0.25, 1.00}  // Prefer using Alpha=1.0 here
	style.Colors[ImGuiCol_TableRowBg] = ImVec4{0.00, 0.00, 0.00, 0.00}
	style.Colors[ImGuiCol_TableRowBgAlt] = ImVec4{1.00, 1.00, 1.00, 0.06}
	style.Colors[ImGuiCol_TextSelectedBg] = ImVec4{0.26, 0.59, 0.98, 0.35}
	style.Colors[ImGuiCol_DragDropTarget] = ImVec4{1.00, 1.00, 0.00, 0.90}
	style.Colors[ImGuiCol_NavHighlight] = ImVec4{0.26, 0.59, 0.98, 1.00}
	style.Colors[ImGuiCol_NavWindowingHighlight] = ImVec4{1.00, 1.00, 1.00, 0.70}
	style.Colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.80, 0.80, 0.80, 0.20}
	style.Colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.80, 0.80, 0.80, 0.35}
}

func StyleColorsClassic(style *ImGuiStyle) {
	if style == nil {
		style = GetStyle()
	}

	style.Colors[ImGuiCol_Text] = ImVec4{0.9, 0.9, 0.9, 1.0}
	style.Colors[ImGuiCol_TextDisabled] = ImVec4{0.6, 0.6, 0.6, 1.0}
	style.Colors[ImGuiCol_WindowBg] = ImVec4{0.0, 0.0, 0.0, 0.8}
	style.Colors[ImGuiCol_ChildBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	style.Colors[ImGuiCol_PopupBg] = ImVec4{0.1, 0.1, 0.1, 0.9}
	style.Colors[ImGuiCol_Border] = ImVec4{0.5, 0.5, 0.5, 0.5}
	style.Colors[ImGuiCol_BorderShadow] = ImVec4{0.0, 0.0, 0.0, 0.0}
	style.Colors[ImGuiCol_FrameBg] = ImVec4{0.4, 0.4, 0.4, 0.3}
	style.Colors[ImGuiCol_FrameBgHovered] = ImVec4{0.4, 0.4, 0.6, 0.4}
	style.Colors[ImGuiCol_FrameBgActive] = ImVec4{0.4, 0.4, 0.6, 0.6}
	style.Colors[ImGuiCol_TitleBg] = ImVec4{0.2, 0.2, 0.5, 0.8}
	style.Colors[ImGuiCol_TitleBgActive] = ImVec4{0.3, 0.3, 0.6, 0.8}
	style.Colors[ImGuiCol_TitleBgCollapsed] = ImVec4{0.4, 0.4, 0.8, 0.2}
	style.Colors[ImGuiCol_MenuBarBg] = ImVec4{0.4, 0.4, 0.5, 0.8}
	style.Colors[ImGuiCol_ScrollbarBg] = ImVec4{0.2, 0.2, 0.3, 0.6}
	style.Colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.4, 0.4, 0.8, 0.3}
	style.Colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.4, 0.4, 0.8, 0.4}
	style.Colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.4, 0.3, 0.8, 0.6}
	style.Colors[ImGuiCol_CheckMark] = ImVec4{0.9, 0.9, 0.9, 0.5}
	style.Colors[ImGuiCol_SliderGrab] = ImVec4{1.0, 1.0, 1.0, 0.3}
	style.Colors[ImGuiCol_SliderGrabActive] = ImVec4{0.4, 0.3, 0.8, 0.6}
	style.Colors[ImGuiCol_Button] = ImVec4{0.3, 0.4, 0.6, 0.6}
	style.Colors[ImGuiCol_ButtonHovered] = ImVec4{0.4, 0.4, 0.7, 0.7}
	style.Colors[ImGuiCol_ButtonActive] = ImVec4{0.4, 0.5, 0.8, 1.0}
	style.Colors[ImGuiCol_Header] = ImVec4{0.4, 0.4, 0.9, 0.4}
	style.Colors[ImGuiCol_HeaderHovered] = ImVec4{0.4, 0.4, 0.9, 0.8}
	style.Colors[ImGuiCol_HeaderActive] = ImVec4{0.5, 0.5, 0.8, 0.8}
	style.Colors[ImGuiCol_Separator] = ImVec4{0.5, 0.5, 0.5, 0.6}
	style.Colors[ImGuiCol_SeparatorHovered] = ImVec4{0.6, 0.6, 0.7, 1.0}
	style.Colors[ImGuiCol_SeparatorActive] = ImVec4{0.7, 0.7, 0.9, 1.0}
	style.Colors[ImGuiCol_ResizeGrip] = ImVec4{1.0, 1.0, 1.0, 0.1}
	style.Colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.7, 0.8, 1.0, 0.6}
	style.Colors[ImGuiCol_ResizeGripActive] = ImVec4{0.7, 0.8, 1.0, 0.9}
	style.Colors[ImGuiCol_Tab] = ImLerpVec4(&style.Colors[ImGuiCol_Header], &style.Colors[ImGuiCol_TitleBgActive], 0.8)
	style.Colors[ImGuiCol_TabHovered] = style.Colors[ImGuiCol_HeaderHovered]
	style.Colors[ImGuiCol_TabActive] = ImLerpVec4(&style.Colors[ImGuiCol_HeaderActive], &style.Colors[ImGuiCol_TitleBgActive], 0.6)
	style.Colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&style.Colors[ImGuiCol_Tab], &style.Colors[ImGuiCol_TitleBg], 0.8)
	style.Colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&style.Colors[ImGuiCol_TabActive], &style.Colors[ImGuiCol_TitleBg], 0.4)
	style.Colors[ImGuiCol_PlotLines] = ImVec4{1.0, 1.0, 1.0, 1.0}
	style.Colors[ImGuiCol_PlotLinesHovered] = ImVec4{0.9, 0.7, 0.0, 1.0}
	style.Colors[ImGuiCol_PlotHistogram] = ImVec4{0.9, 0.7, 0.0, 1.0}
	style.Colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.0, 0.6, 0.0, 1.0}
	style.Colors[ImGuiCol_TableHeaderBg] = ImVec4{0.2, 0.2, 0.3, 1.0}
	style.Colors[ImGuiCol_TableBorderStrong] = ImVec4{0.3, 0.3, 0.4, 1.0} // Prefer using Alpha=1.0 here
	style.Colors[ImGuiCol_TableBorderLight] = ImVec4{0.2, 0.2, 0.2, 1.0}  // Prefer using Alpha=1.0 here
	style.Colors[ImGuiCol_TableRowBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	style.Colors[ImGuiCol_TableRowBgAlt] = ImVec4{1.0, 1.0, 1.0, 0.0}
	style.Colors[ImGuiCol_TextSelectedBg] = ImVec4{0.0, 0.0, 1.0, 0.3}
	style.Colors[ImGuiCol_DragDropTarget] = ImVec4{1.0, 1.0, 0.0, 0.9}
	style.Colors[ImGuiCol_NavHighlight] = style.Colors[ImGuiCol_HeaderHovered]
	style.Colors[ImGuiCol_NavWindowingHighlight] = ImVec4{1.0, 1.0, 1.0, 0.7}
	style.Colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.8, 0.8, 0.8, 0.2}
	style.Colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.2, 0.2, 0.2, 0.3}
}

// StyleColorsLight Those light colors are better suited with a thicker font than the default one + FrameBorder
func StyleColorsLight(style *ImGuiStyle) {
	if style == nil {
		style = GetStyle()
	}

	style.Colors[ImGuiCol_Text] = ImVec4{0.0, 0.0, 0.0, 1.0}
	style.Colors[ImGuiCol_TextDisabled] = ImVec4{0.6, 0.6, 0.6, 1.0}
	style.Colors[ImGuiCol_WindowBg] = ImVec4{0.9, 0.9, 0.9, 1.0}
	style.Colors[ImGuiCol_ChildBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	style.Colors[ImGuiCol_PopupBg] = ImVec4{1.0, 1.0, 1.0, 0.9}
	style.Colors[ImGuiCol_Border] = ImVec4{0.0, 0.0, 0.0, 0.3}
	style.Colors[ImGuiCol_BorderShadow] = ImVec4{0.0, 0.0, 0.0, 0.0}
	style.Colors[ImGuiCol_FrameBg] = ImVec4{1.0, 1.0, 1.0, 1.0}
	style.Colors[ImGuiCol_FrameBgHovered] = ImVec4{0.2, 0.5, 0.9, 0.4}
	style.Colors[ImGuiCol_FrameBgActive] = ImVec4{0.2, 0.5, 0.9, 0.6}
	style.Colors[ImGuiCol_TitleBg] = ImVec4{0.9, 0.9, 0.9, 1.0}
	style.Colors[ImGuiCol_TitleBgActive] = ImVec4{0.8, 0.8, 0.8, 1.0}
	style.Colors[ImGuiCol_TitleBgCollapsed] = ImVec4{1.0, 1.0, 1.0, 0.5}
	style.Colors[ImGuiCol_MenuBarBg] = ImVec4{0.8, 0.8, 0.8, 1.0}
	style.Colors[ImGuiCol_ScrollbarBg] = ImVec4{0.9, 0.9, 0.9, 0.5}
	style.Colors[ImGuiCol_ScrollbarGrab] = ImVec4{0.6, 0.6, 0.6, 0.8}
	style.Colors[ImGuiCol_ScrollbarGrabHovered] = ImVec4{0.4, 0.4, 0.4, 0.8}
	style.Colors[ImGuiCol_ScrollbarGrabActive] = ImVec4{0.4, 0.4, 0.4, 1.0}
	style.Colors[ImGuiCol_CheckMark] = ImVec4{0.2, 0.5, 0.9, 1.0}
	style.Colors[ImGuiCol_SliderGrab] = ImVec4{0.2, 0.5, 0.9, 0.7}
	style.Colors[ImGuiCol_SliderGrabActive] = ImVec4{0.4, 0.5, 0.8, 0.6}
	style.Colors[ImGuiCol_Button] = ImVec4{0.2, 0.5, 0.9, 0.4}
	style.Colors[ImGuiCol_ButtonHovered] = ImVec4{0.2, 0.5, 0.9, 1.0}
	style.Colors[ImGuiCol_ButtonActive] = ImVec4{0.0, 0.5, 0.9, 1.0}
	style.Colors[ImGuiCol_Header] = ImVec4{0.2, 0.5, 0.9, 0.3}
	style.Colors[ImGuiCol_HeaderHovered] = ImVec4{0.2, 0.5, 0.9, 0.8}
	style.Colors[ImGuiCol_HeaderActive] = ImVec4{0.2, 0.5, 0.9, 1.0}
	style.Colors[ImGuiCol_Separator] = ImVec4{0.3, 0.3, 0.3, 0.6}
	style.Colors[ImGuiCol_SeparatorHovered] = ImVec4{0.1, 0.4, 0.8, 0.7}
	style.Colors[ImGuiCol_SeparatorActive] = ImVec4{0.1, 0.4, 0.8, 1.0}
	style.Colors[ImGuiCol_ResizeGrip] = ImVec4{0.3, 0.3, 0.3, 0.1}
	style.Colors[ImGuiCol_ResizeGripHovered] = ImVec4{0.2, 0.5, 0.9, 0.6}
	style.Colors[ImGuiCol_ResizeGripActive] = ImVec4{0.2, 0.5, 0.9, 0.9}
	style.Colors[ImGuiCol_Tab] = ImLerpVec4(&style.Colors[ImGuiCol_Header], &style.Colors[ImGuiCol_TitleBgActive], 0.9)
	style.Colors[ImGuiCol_TabHovered] = style.Colors[ImGuiCol_HeaderHovered]
	style.Colors[ImGuiCol_TabActive] = ImLerpVec4(&style.Colors[ImGuiCol_HeaderActive], &style.Colors[ImGuiCol_TitleBgActive], 0.6)
	style.Colors[ImGuiCol_TabUnfocused] = ImLerpVec4(&style.Colors[ImGuiCol_Tab], &style.Colors[ImGuiCol_TitleBg], 0.8)
	style.Colors[ImGuiCol_TabUnfocusedActive] = ImLerpVec4(&style.Colors[ImGuiCol_TabActive], &style.Colors[ImGuiCol_TitleBg], 0.4)
	style.Colors[ImGuiCol_PlotLines] = ImVec4{0.3, 0.3, 0.3, 1.0}
	style.Colors[ImGuiCol_PlotLinesHovered] = ImVec4{1.0, 0.4, 0.3, 1.0}
	style.Colors[ImGuiCol_PlotHistogram] = ImVec4{0.9, 0.7, 0.0, 1.0}
	style.Colors[ImGuiCol_PlotHistogramHovered] = ImVec4{1.0, 0.4, 0.0, 1.0}
	style.Colors[ImGuiCol_TableHeaderBg] = ImVec4{0.7, 0.8, 0.9, 1.0}
	style.Colors[ImGuiCol_TableBorderStrong] = ImVec4{0.5, 0.5, 0.6, 1.0} // Prefer using Alpha=1.0 here
	style.Colors[ImGuiCol_TableBorderLight] = ImVec4{0.6, 0.6, 0.7, 1.0}  // Prefer using Alpha=1.0 here
	style.Colors[ImGuiCol_TableRowBg] = ImVec4{0.0, 0.0, 0.0, 0.0}
	style.Colors[ImGuiCol_TableRowBgAlt] = ImVec4{0.3, 0.3, 0.3, 0.0}
	style.Colors[ImGuiCol_TextSelectedBg] = ImVec4{0.2, 0.5, 0.9, 0.3}
	style.Colors[ImGuiCol_DragDropTarget] = ImVec4{0.2, 0.5, 0.9, 0.9}
	style.Colors[ImGuiCol_NavHighlight] = style.Colors[ImGuiCol_HeaderHovered]
	style.Colors[ImGuiCol_NavWindowingHighlight] = ImVec4{0.7, 0.7, 0.7, 0.7}
	style.Colors[ImGuiCol_NavWindowingDimBg] = ImVec4{0.2, 0.2, 0.2, 0.2}
	style.Colors[ImGuiCol_ModalWindowDimBg] = ImVec4{0.2, 0.2, 0.2, 0.3}
}
