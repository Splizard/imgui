package imgui

import "fmt"

// BeginTooltip Tooltips
// - Tooltip are windows following the mouse. They do not take focus away.
// begin/append a tooltip window. to create full-featured tooltip (with any kind of items).
func BeginTooltip() {
	BeginTooltipEx(ImGuiWindowFlags_None, ImGuiTooltipFlags_None)
}

func EndTooltip() {
	IM_ASSERT(GetCurrentWindowRead().Flags&ImGuiWindowFlags_Tooltip != 0) // Mismatched BeginTooltip()/EndTooltip() calls
	End()
}

// SetTooltip set a text-only tooltip, typically use with ImGui::IsItemHovered(). override any previous call to SetTooltip().
func SetTooltip(format string, args ...interface{}) {
	BeginTooltipEx(0, ImGuiTooltipFlags_OverridePreviousTooltip)
	Text(format, args...)
	EndTooltip()
}

func BeginTooltipEx(extra_flags ImGuiWindowFlags, tooltip_flags ImGuiTooltipFlags) {
	var g = GImGui

	if g.DragDropWithinSource || g.DragDropWithinTarget {
		// The default tooltip position is a little offset to give space to see the context menu (it's also clamped within the current viewport/monitor)
		// In the context of a dragging tooltip we try to reduce that offset and we enforce following the cursor.
		// Whatever we do we want to call SetNextWindowPos() to enforce a tooltip position and disable clipping the tooltip without our display area, like regular tooltip do.
		//ImVec2 tooltip_pos = g.IO.MousePos - g.ActiveIdClickOffset - g.Style.WindowPadding;
		var tooltip_pos ImVec2 = g.IO.MousePos.Add(ImVec2{16 * g.Style.MouseCursorScale, 8 * g.Style.MouseCursorScale})
		SetNextWindowPos(&tooltip_pos, 0, ImVec2{})
		SetNextWindowBgAlpha(g.Style.Colors[ImGuiCol_PopupBg].w * 0.60)
		//PushStyleVar(ImGuiStyleVar_Alpha, g.Style.Alpha * 0.60f); // This would be nice but e.g ColorButton with checkboard has issue with transparent colors :(
		tooltip_flags |= ImGuiTooltipFlags_OverridePreviousTooltip
	}

	var window_name = fmt.Sprintf("##Tooltip_%02d", g.TooltipOverrideCount)
	if tooltip_flags&ImGuiTooltipFlags_OverridePreviousTooltip != 0 {
		if window := FindWindowByName(window_name); window != nil {
			if window.Active {
				// Hide previous tooltip from being displayed. We can't easily "reset" the content of a window so we create a new one.
				window.Hidden = true
				window.HiddenFramesCanSkipItems = 1 // FIXME: This may not be necessary?
				g.TooltipOverrideCount++
				window_name = fmt.Sprintf("##Tooltip_%02d", g.TooltipOverrideCount)
			}
		}
	}
	var flags ImGuiWindowFlags = ImGuiWindowFlags_Tooltip | ImGuiWindowFlags_NoInputs | ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoMove | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoSavedSettings | ImGuiWindowFlags_AlwaysAutoResize
	Begin(window_name, nil, flags|extra_flags)
}
