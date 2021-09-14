package imgui

// Debugging enums
const (
	WRT_OuterRect = iota
	WRT_InnerRect
	WRT_OuterRectClipped
	WRT_InnerClipRect
	WRT_WorkRect
	WRT_Content
	WRT_ContentIdeal
	WRT_ContentRegionRect
	WRT_Count
)

var wrt_rects_names = []string{
	"OuterRect", "OuterRectClipped", "InnerRect", "InnerClipRect", "WorkRect", "Content", "ContentIdeal", "ContentRegionRect",
}

const (
	TRT_OuterRect = iota
	TRT_InnerRect
	TRT_WorkRect
	TRT_HostClipRect
	TRT_InnerClipRect
	TRT_BackgroundClipRect
	TRT_ColumnsRect
	TRT_ColumnsWorkRect
	TRT_ColumnsClipRect
	TRT_ColumnsContentHeadersUsed
	TRT_ColumnsContentHeadersIdeal
	TRT_ColumnsContentFrozen
	TRT_ColumnsContentUnfrozen
	TRT_Count
)

var trt_rects_names = []string{
	"OuterRect", "InnerRect", "WorkRect", "HostClipRect", "InnerClipRect", "BackgroundClipRect", "ColumnsRect", "ColumnsWorkRect", "ColumnsClipRect", "ColumnsContentHeadersUsed", "ColumnsContentHeadersIdeal", "ColumnsContentFrozen", "ColumnsContentUnfrozen",
}

// [DEBUG] Item picker tool - start with DebugStartItemPicker() - useful to visually select an item and break into its call-stack.
func UpdateDebugToolItemPicker() {
	var g = GImGui
	g.DebugItemPickerBreakId = 0
	if g.DebugItemPickerActive {
		var hovered_id ImGuiID = g.HoveredIdPreviousFrame
		SetMouseCursor(ImGuiMouseCursor_Hand)
		if IsKeyPressedMap(ImGuiKey_Escape, true) {
			g.DebugItemPickerActive = false
		}
		if IsMouseClicked(0, true) && hovered_id != 0 {
			g.DebugItemPickerBreakId = hovered_id
			g.DebugItemPickerActive = false
		}
		SetNextWindowBgAlpha(0.60)
		BeginTooltip()
		Text("HoveredId: 0x%08X", hovered_id)
		Text("Press ESC to abort picking.")
		var c *ImVec4
		if hovered_id != 0 {
			c = GetStyleColorVec4(ImGuiCol_Text)
		} else {
			c = GetStyleColorVec4(ImGuiCol_TextDisabled)
		}
		TextColored(c, "Click to break in debugger!")
		EndTooltip()
	}
}

// create Metrics/Debugger window. display Dear ImGui internals: windows, draw commands, various internal state, etc.
func ShowMetricsWindow(p_open *bool) {
	if !Begin("Dear ImGui Metrics/Debugger", p_open, 0) {
		End()
		return
	}

	var g = GImGui
	var io = g.IO
	var cfg = g.DebugMetricsConfig

	Text("Dear ImGui %s", GetVersion())
	Text("Application average %.3f ms/frame (%.1f FPS)", 1000.0/io.Framerate, io.Framerate)
	Text("%d vertices, %d indices (%d triangles)", io.MetricsRenderVertices, io.MetricsRenderIndices, io.MetricsRenderIndices/3)
	Text("%d active windows (%d visible)", io.MetricsActiveWindows, io.MetricsRenderWindows)
	Text("%d active allocations", io.MetricsActiveAllocations)

	Separator()

	if cfg.ShowWindowsRectsType < 0 {
		cfg.ShowWindowsRectsType = WRT_WorkRect
	}
	if cfg.ShowTablesRectsType < 0 {
		cfg.ShowTablesRectsType = TRT_WorkRect
	}

	/*GetTableRect := func(table *ImGuiTable, rect_type, n int) ImRect {
		switch rect_type {
		case TRT_OuterRect:
			return table.OuterRect
		case TRT_InnerRect:
			return table.InnerRect
		case TRT_WorkRect:
			return table.WorkRect
		case TRT_HostClipRect:
			return table.HostClipRect
		case TRT_InnerClipRect:
			return table.InnerClipRect
		case TRT_BackgroundClipRect:
			return table.BgClipRect
		case TRT_ColumnsRect:
			var c = &table.Columns[n]
			return ImRect(c.MinX, table.InnerClipRect.Min.y, c.MaxX, table.InnerClipRect.Min.y+table.LastOuterHeight)
		case TRT_ColumnsWorkRect:
			var c = &table.Columns[n]
			return ImRect(c.WorkMinX, table.WorkRect.Min.y, c.WorkMaxX, table.WorkRect.Max.y)
		case TRT_ColumnsClipRect:
			var c = &table.Columns[n]
			return c.ClipRect
		case TRT_ColumnsContentHeadersUsed:
			var c = &table.Columns[n] // Note: y1/y2 not always accurate
			return ImRect(c.WorkMinX, table.InnerClipRect.Min.y, c.ContentMaxXHeadersUsed, table.InnerClipRect.Min.y+table.LastFirstRowHeight)
		case TRT_ColumnsContentHeadersIdeal:
			var c = &table.Columns[n]
			return ImRect(c.WorkMinX, table.InnerClipRect.Min.y, c.ContentMaxXHeadersIdeal, table.InnerClipRect.Min.y+table.LastFirstRowHeight)
		case TRT_ColumnsContentFrozen:
			var c = &table.Columns[n]
			return ImRect(c.WorkMinX, table.InnerClipRect.Min.y, c.ContentMaxXFrozen, table.InnerClipRect.Min.y+table.LastFirstRowHeight)
		case TRT_ColumnsContentUnfrozen:
			var c = &table.Columns[n]
			return ImRect(c.WorkMinX, table.InnerClipRect.Min.y+table.LastFirstRowHeight, c.ContentMaxXUnfrozen, table.InnerClipRect.Max.y)
		}
		IM_ASSERT(false)
		return ImRect{}
	}*/

	/*GetWindowRect := func(window *ImGuiWindow, rect_type int) ImRect {
		switch rect_type {
		case WRT_OuterRect:
			return window.Rect()
		case WRT_OuterRectClipped:
			return window.OuterRectClipped
		case WRT_InnerRect:
			return window.InnerRect
		case WRT_InnerClipRect:
			return window.InnerClipRect
		case WRT_WorkRect:
			return window.WorkRect
		case WRT_Content:
			var min ImVec2 = window.InnerRect.Min.Sub(window.Scroll).Add(window.WindowPadding)
			return ImRect{min, min.Add(window.ContentSize)}
		case WRT_ContentIdeal:
			var min ImVec2 = window.InnerRect.Min.Sub(window.Scroll).Add(window.WindowPadding)
			return ImRect{min, min.Add(window.ContentSizeIdeal)}
		case WRT_ContentRegionRect:
			return window.ContentRegionRect
		}
		IM_ASSERT(false)
		return ImRect{}
	}*/

	//Tools
	if TreeNode("Tools") {
		if Button("Item Picker") {
			DebugStartItemPicker()
		}

		TreePop()
	}

	End()
}
