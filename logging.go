package imgui

import (
	"fmt"
	"os"
)

// LogBegin Logging/Capture
// . BeginCapture() when we design v2 api, for now stay under the radar by using the old name.
func LogBegin(ltype ImGuiLogType, auto_open_depth int) {
	var g = GImGui
	var window = g.CurrentWindow
	IM_ASSERT(!g.LogEnabled)
	IM_ASSERT(g.LogFile == nil)
	IM_ASSERT(g.LogBuffer.Len() == 0)
	g.LogEnabled = true
	g.LogType = ltype
	g.LogNextPrefix = ""
	g.LogNextSuffix = ""
	g.LogDepthRef = window.DC.TreeDepth
	g.LogDepthToExpand = g.LogDepthToExpandDefault
	if auto_open_depth >= 0 {
		g.LogDepthToExpand = auto_open_depth
	}
	g.LogLinePosY = FLT_MAX
	g.LogLineFirstItem = true
}

// LogToBuffer Start logging/capturing to internal buffer
func LogToBuffer(auto_open_depth int /*= -1*/) {
	var g = GImGui
	if g.LogEnabled {
		return
	}
	LogBegin(ImGuiLogType_Buffer, auto_open_depth)
}

// LogRenderedText Internal version that takes a position to decide on newline placement and pad items according to their depth.
// We split text into individual lines to add current tree level padding
// FIXME: This code is a little complicated perhaps, considering simplifying the whole system.
func LogRenderedText(ref_pos *ImVec2, text string) {
	LogText(text)
}

// LogSetNextTextDecoration Important: doesn't copy underlying data, use carefully (prefix/suffix must be in scope at the time of the next LogRenderedText)
func LogSetNextTextDecoration(prefix string, suffix string) {
	var g = GImGui
	g.LogNextPrefix = prefix
	g.LogNextSuffix = suffix
}

// Logging/Capture
// - All text output from the interface can be captured into tty/file/clipboard. By default, tree nodes are automatically opened during logging.

// LogToTTY start logging to tty (stdout)
func LogToTTY(auto_open_depth int /*= -1*/) {
	var g = GImGui
	if g.LogEnabled {
		return
	}
	LogBegin(ImGuiLogType_TTY, auto_open_depth)
	g.LogFile = os.Stdout
}

// LogToFile Start logging/capturing text output to given file
func LogToFile(auto_open_depth int /*= 1*/, filename string) {
	var g = GImGui
	if g.LogEnabled {
		return
	}

	// FIXME: We could probably open the file in text mode "at", however note that clipboard/buffer logging will still
	// be subject to outputting OS-incompatible carriage return if within strings the user doesn't use IM_NEWLINE.
	// By opening the file in binary mode "ab" we have consistent output everywhere.
	if filename == "" {
		filename = g.IO.LogFilename
	}
	if filename == "" || filename[0] == 0 {
		return
	}
	var f = ImFileOpen(filename, "ab")
	if f == nil {
		IM_ASSERT(false)
		return
	}

	LogBegin(ImGuiLogType_File, auto_open_depth)
	g.LogFile = f
}

// LogToClipboard start logging to OS clipboard
func LogToClipboard(auto_open_depth int /*= -1*/) {
	var g = GImGui
	if g.LogEnabled {
		return
	}
	LogBegin(ImGuiLogType_Clipboard, auto_open_depth)
}

func LogFinish() {
	var g = GImGui
	if !g.LogEnabled {
		return
	}

	LogText(IM_NEWLINE)
	switch g.LogType {
	case ImGuiLogType_TTY:
		//g.LogFile
	case ImGuiLogType_File:
		ImFileClose(g.LogFile)
	case ImGuiLogType_Buffer:
	case ImGuiLogType_Clipboard:
		if g.LogBuffer.Len() > 0 {
			SetClipboardText(g.LogBuffer.String())
		}
	case ImGuiLogType_None:
		IM_ASSERT(false)
	}

	g.LogEnabled = false
	g.LogType = ImGuiLogType_None
	g.LogFile = nil
	g.LogBuffer.Reset()
} // stop logging (close file, etc.)

// LogButtons helper to display buttons for logging to tty/file/clipboard
func LogButtons() {
	var g = GImGui

	PushString("LogButtons")
	var log_to_tty = Button("Log To TTY")
	SameLine(0, 0)
	var log_to_file = Button("Log To File")
	SameLine(0, 0)
	var log_to_clipboard = Button("Log To Clipboard")
	SameLine(0, 0)
	PushAllowKeyboardFocus(false)
	SetNextItemWidth(80.0)
	SliderInt("Default Depth", &g.LogDepthToExpandDefault, 0, 9, "", 0)
	PopAllowKeyboardFocus()
	PopID()

	// Start logging at the end of the function so that the buttons don't appear in the log
	if log_to_tty {
		LogToTTY(-1)
	}
	if log_to_file {
		LogToFile(-1, "")
	}
	if log_to_clipboard {
		LogToClipboard(-1)
	}
}

// LogText pass text data straight to log (without being displayed)
func LogText(format string, args ...any) {
	var g = GImGui
	if !g.LogEnabled {
		return
	}

	if g.LogFile != nil {
		fmt.Fprintf(g.LogFile, format, args...)
	} else {
		fmt.Fprintf(&g.LogBuffer, format, args...)
	}
}
