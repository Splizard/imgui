package imgui

import "fmt"

// Apply to existing windows (if any)
func WindowSettingsHandler_ApplyAll(ctx *ImGuiContext, _ *ImGuiSettingsHandler) {
	var g = ctx
	for i := range g.SettingsWindows {
		settings := &g.SettingsWindows[i]
		if settings.WantApply {
			if window := FindWindowByID(settings.ID); window != nil {
				ApplyWindowSettings(window, settings)
			}
			settings.WantApply = false
		}
	}
}

func ApplyWindowSettings(window *ImGuiWindow, settings *ImGuiWindowSettings) {
	window.Pos = *ImFloorVec(&ImVec2{float(settings.Pos.x), float(settings.Pos.y)})
	if settings.Size.x > 0 && settings.Size.y > 0 {
		window.Size = *ImFloorVec(&ImVec2{float(settings.Size.x), float(settings.Size.y)})
		window.SizeFull = *ImFloorVec(&ImVec2{float(settings.Size.x), float(settings.Size.y)})
	}
	window.Collapsed = settings.Collapsed
}

func WindowSettingsHandler_ClearAll(ctx *ImGuiContext, _ *ImGuiSettingsHandler) {
	var g = ctx
	for i := range g.Windows {
		g.Windows[i].SettingsOffset = -1
	}
	g.SettingsWindows = g.SettingsWindows[0:]
}

func WindowSettingsHandler_ReadOpen(_ *ImGuiContext, _ *ImGuiSettingsHandler, name string) interface{} {
	var settings = FindOrCreateWindowSettings(name)
	var id = settings.ID
	*settings = ImGuiWindowSettings{} // Clear existing if recycling previous entry
	settings.ID = id
	settings.WantApply = true
	return settings
}

func WindowSettingsHandler_ReadLine(_ *ImGuiContext, _ *ImGuiSettingsHandler, entry interface{}, line string) {
	var settings = entry.(*ImGuiWindowSettings)
	var x, y int
	var i int
	if n, _ := fmt.Sscanf(line, "Pos=%i,%i", &x, &y); n == 2 {
		settings.Pos = ImVec2ih{(short)(x), (short)(y)}
	} else if n, _ := fmt.Sscanf(line, "Size=%i,%i", &x, &y); n == 2 {
		settings.Size = ImVec2ih{(short)(x), (short)(y)}
	} else if n, _ := fmt.Sscanf(line, "Collapsed=%d", &i); n == 1 {
		settings.Collapsed = (i != 0)
	}
}

func WindowSettingsHandler_WriteAll(ctx *ImGuiContext, handler *ImGuiSettingsHandler, buf *ImGuiTextBuffer) {
	// Gather data from windows that were active during this session
	// (if a window wasn't opened in this session we preserve its settings)
	var g = ctx
	for i := range g.Windows {
		var window = g.Windows[i]
		if window.Flags&ImGuiWindowFlags_NoSavedSettings != 0 {
			continue
		}

		var settings *ImGuiWindowSettings
		if window.SettingsOffset != -1 {
			settings = &g.SettingsWindows[window.SettingsOffset]
		} else {
			settings = FindOrCreateWindowSettings(window.Name)
		}

		if settings == nil {
			settings = CreateNewWindowSettings(window.Name)

			window.SettingsOffset = -1
			for i := range g.SettingsWindows {
				if settings == &g.SettingsWindows[i] {
					window.SettingsOffset = int(i)
					break
				}
			}
		}
		IM_ASSERT(settings.ID == window.ID)
		settings.Pos = ImVec2ih{short(window.Pos.x), short(window.Pos.y)}
		settings.Size = ImVec2ih{short(window.SizeFull.x), short(window.SizeFull.y)}

		settings.Collapsed = window.Collapsed
	}

	// Write to text buffer
	for i := range g.SettingsWindows {
		settings := &g.SettingsWindows[i]
		var settings_name = settings.GetName()
		*buf = append(*buf, []byte(fmt.Sprintf("[%s][%s]\n", handler.TypeName, settings_name))...)
		*buf = append(*buf, []byte(fmt.Sprintf("Pos=%d,%d\n", settings.Pos.x, settings.Pos.y))...)
		*buf = append(*buf, []byte(fmt.Sprintf("Size=%d,%d\n", settings.Size.x, settings.Size.y))...)
		*buf = append(*buf, []byte(fmt.Sprintf("Collapsed=%d\n", settings.Collapsed))...)
		*buf = append(*buf, []byte("\n")...)
	}
}
