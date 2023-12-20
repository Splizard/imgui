package imgui

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// LoadIniSettingsFromDisk Settings/.Ini Utilities
// - The disk functions are automatically called if io.IniFilename != NULL (default is "imgui.ini").
// - Set io.IniFilename to NULL to load/save manually. Read io.WantSaveIniSettings description about handling .ini saving manually.
// - Important: default value "imgui.ini" is relative to current working dir! Most apps will want to lock this to an absolute path (e.g. same path as executables).
func LoadIniSettingsFromDisk(ini_filename string) {
	var file_data_size uintptr = 0
	var file_data []byte = ImFileLoadToMemory(ini_filename, "rb", &file_data_size, 0)
	if file_data == nil {
		return
	}
	LoadIniSettingsFromMemory(file_data, (size_t)(file_data_size))
} // call after CreateContext() and before the first call to NewFrame(). NewFrame() automatically calls LoadIniSettingsFromDisk(io.IniFilename).

func LoadIniSettingsFromMemory(buf []byte, ini_size uintptr) {
	var g = GImGui
	IM_ASSERT(g.Initialized)
	//IM_ASSERT(!g.WithinFrameScope && "Cannot be called between NewFrame() and EndFrame()");
	//IM_ASSERT(g.SettingsLoaded == false && g.FrameCount == 0);

	// Call pre-read handlers
	// Some types will clear their data (e.g. dock information) some types will allow merge/override (window)
	for handler_n := range g.SettingsHandlers {
		if g.SettingsHandlers[handler_n].ReadInitFn != nil {
			g.SettingsHandlers[handler_n].ReadInitFn(g, &g.SettingsHandlers[handler_n])
		}
	}

	var reader = bufio.NewReader(bytes.NewReader(buf))
	var entry_handler *ImGuiSettingsHandler
	var entry_data interface{} = nil

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = line[:len(line)-1]

		if line == "" {
			continue
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			splits := strings.SplitN(line[1:len(line)-1], "][", 2)
			if len(splits) == 2 {
				var settings_type = splits[0]
				var settings_id = splits[1]

				entry_handler = FindSettingsHandler(settings_type)
				if entry_handler != nil {
					entry_data = entry_handler.ReadOpenFn(g, entry_handler, settings_id)
				}
			}
		} else if entry_handler != nil {
			entry_handler.ReadLineFn(g, entry_handler, entry_data, line)
		}
	}

	g.SettingsLoaded = true

	// Call post-read handlers
	for handler_n := range g.SettingsHandlers {
		if g.SettingsHandlers[handler_n].ApplyAllFn != nil {
			g.SettingsHandlers[handler_n].ApplyAllFn(g, &g.SettingsHandlers[handler_n])
		}
	}
} // call after CreateContext() and before the first call to NewFrame() to provide .ini data from your own data source.

func SaveIniSettingsToDisk(ini_filename string) {
	var g = GImGui
	g.SettingsDirtyTimer = 0.0
	if ini_filename == "" {
		return
	}

	var ini_data_size size_t = 0
	var ini_data []byte = SaveIniSettingsToMemory(&ini_data_size)

	os.WriteFile(ini_filename, ini_data, 0666)
} // this is automatically called (if io.IniFilename is not empty) a few seconds after any modification that should be reflected in the .ini file (and also by DestroyContext).

func SaveIniSettingsToMemory(out_size *uintptr) []byte {
	var g = GImGui
	g.SettingsDirtyTimer = 0.0
	g.SettingsIniData = g.SettingsIniData[:0]
	for handler_n := range g.SettingsHandlers {
		var handler *ImGuiSettingsHandler = &g.SettingsHandlers[handler_n]
		handler.WriteAllFn(g, handler, &g.SettingsIniData)
	}
	if out_size != nil {
		*out_size = (size_t)(len(g.SettingsIniData))
	}
	return g.SettingsIniData
} // return a zero-terminated string with the .ini data which you can save by your own mean. call when io.WantSaveIniSettings is set, then save data by your own mean and clear io.WantSaveIniSettings.

// MarkIniSettingsDirty Settings
func MarkIniSettingsDirty() {
	var g = GImGui
	if g.SettingsDirtyTimer <= 0.0 {
		g.SettingsDirtyTimer = g.IO.IniSavingRate
	}
}

func MarkIniSettingsDirtyWindow(window *ImGuiWindow) {
	var g = GImGui
	if window.Flags&ImGuiWindowFlags_NoSavedSettings == 0 {
		if g.SettingsDirtyTimer <= 0.0 {
			g.SettingsDirtyTimer = g.IO.IniSavingRate
		}
	}
}

func ClearIniSettings() {
	var g = GImGui
	g.SettingsIniData = g.SettingsIniData[:0]
	for handler_n := range g.SettingsHandlers {
		if g.SettingsHandlers[handler_n].ClearAllFn != nil {
			g.SettingsHandlers[handler_n].ClearAllFn(g, &g.SettingsHandlers[handler_n])
		}
	}
}

func CreateNewWindowSettings(name string) *ImGuiWindowSettings {
	var g = GImGui

	if index := strings.Index(name, "###"); index != -1 {
		name = name[index:]
	}

	var settings ImGuiWindowSettings

	settings.ID = ImHashStr(name, 0, 0)
	settings.name = name

	g.SettingsWindows = append(g.SettingsWindows, settings)

	return &g.SettingsWindows[len(g.SettingsWindows)-1]
}

func FindWindowSettings(id ImGuiID) *ImGuiWindowSettings {
	var g = GImGui
	for i := range g.SettingsWindows {
		settings := &g.SettingsWindows[i]
		if settings.ID == id {
			return settings
		}
	}
	return nil
}

func FindOrCreateWindowSettings(name string) *ImGuiWindowSettings {
	if settings := FindWindowSettings(ImHashStr(name, 0, 0)); settings != nil {
		return settings
	}
	return CreateNewWindowSettings(name)
}

func FindSettingsHandler(name string) *ImGuiSettingsHandler {
	var g = GImGui
	var type_hash ImGuiID = ImHashStr(name, 0, 0)
	for handler_n := range g.SettingsHandlers {
		if g.SettingsHandlers[handler_n].TypeHash == type_hash {
			return &g.SettingsHandlers[handler_n]
		}
	}
	return nil
}

// UpdateSettings Called by NewFrame()
func UpdateSettings() {
	// Load settings on first frame (if not explicitly loaded manually before)
	var g *ImGuiContext = GImGui
	if !g.SettingsLoaded {
		IM_ASSERT(len(g.SettingsWindows) == 0)
		if g.IO.IniFilename != "" {
			LoadIniSettingsFromDisk(g.IO.IniFilename)
		}
		g.SettingsLoaded = true
	}

	// Save settings (with a delay after the last modification, so we don't spam disk too much)
	if g.SettingsDirtyTimer > 0.0 {
		g.SettingsDirtyTimer -= g.IO.DeltaTime
		if g.SettingsDirtyTimer <= 0.0 {
			if g.IO.IniFilename != "" {
				SaveIniSettingsToDisk(g.IO.IniFilename)
			} else {
				g.IO.WantSaveIniSettings = true // Let user know they can call SaveIniSettingsToMemory(). user will need to clear io.WantSaveIniSettings themselves.
			}
			g.SettingsDirtyTimer = 0
		}
	}
}

// WindowSettingsHandler_ApplyAll Apply to existing windows (if any)
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
	g.SettingsWindows = g.SettingsWindows[:0]
}

func WindowSettingsHandler_ReadOpen(_ *ImGuiContext, _ *ImGuiSettingsHandler, name string) interface{} {
	var settings = FindOrCreateWindowSettings(name)
	var id = settings.ID
	*settings = ImGuiWindowSettings{} // Clear existing if recycling previous entry
	settings.ID = id
	settings.name = name
	settings.WantApply = true
	return settings
}

func WindowSettingsHandler_ReadLine(_ *ImGuiContext, _ *ImGuiSettingsHandler, entry interface{}, line string) {
	var settings = entry.(*ImGuiWindowSettings)
	var x, y int
	var i bool

	if n, _ := fmt.Sscanf(line, "Pos=%v,%v", &x, &y); n == 2 {
		settings.Pos = ImVec2ih{(short)(x), (short)(y)}
	} else if n, _ := fmt.Sscanf(line, "Size=%v,%v", &x, &y); n == 2 {
		settings.Size = ImVec2ih{(short)(x), (short)(y)}
	} else if n, _ := fmt.Sscanf(line, "Collapsed=%v", &i); n == 1 {
		settings.Collapsed = i
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
		settings.Size = ImVec2ih{short(window.Size.x), short(window.Size.y)}

		settings.Collapsed = window.Collapsed
	}

	// Write to text buffer
	for i := range g.SettingsWindows {
		settings := &g.SettingsWindows[i]
		var settings_name = settings.GetName()
		*buf = append(*buf, []byte(fmt.Sprintf("[%s][%s]\n", handler.TypeName, settings_name))...)
		*buf = append(*buf, []byte(fmt.Sprintf("Pos=%d,%d\n", settings.Pos.x, settings.Pos.y))...)
		*buf = append(*buf, []byte(fmt.Sprintf("Size=%d,%d\n", settings.Size.x, settings.Size.y))...)
		*buf = append(*buf, []byte(fmt.Sprintf("Collapsed=%v\n", settings.Collapsed))...)
		*buf = append(*buf, []byte("\n")...)
	}
}
