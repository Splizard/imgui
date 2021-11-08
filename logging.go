package imgui

import (
	"fmt"
	"os"
)

// Logging/Capture
// . BeginCapture() when we design v2 api, for now stay under the radar by using the old name.
func LogBegin(ltype ImGuiLogType, auto_open_depth int) {
	var g = GImGui
	var window = g.CurrentWindow
	IM_ASSERT(g.LogEnabled == false)
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

// Start logging/capturing to internal buffer
func LogToBuffer(auto_open_depth int /*= -1*/) {
	var g = GImGui
	if g.LogEnabled {
		return
	}
	LogBegin(ImGuiLogType_Buffer, auto_open_depth)
}

// Internal version that takes a position to decide on newline placement and pad items according to their depth.
// We split text into individual lines to add current tree level padding
// FIXME: This code is a little complicated perhaps, considering simplifying the whole system.
func LogRenderedText(ref_pos *ImVec2, text string) {
	var g = GImGui
	var window = g.CurrentWindow

	var prefix = g.LogNextPrefix
	var suffix = g.LogNextSuffix
	g.LogNextPrefix = ""
	g.LogNextSuffix = ""

	text_end := FindRenderedTextEnd(text)

	var log_new_line = ref_pos != nil && (ref_pos.y > g.LogLinePosY+g.Style.FramePadding.y+1)
	if ref_pos != nil {
		g.LogLinePosY = ref_pos.y
	}
	if log_new_line {
		LogText(IM_NEWLINE)
		g.LogLineFirstItem = true
	}

	if prefix != "" {
		LogRenderedText(ref_pos, prefix) // Calculate end ourself to ensure "##" are included here.
	}

	// Re-adjust padding if we have popped out of our starting depth
	if g.LogDepthRef > window.DC.TreeDepth {
		g.LogDepthRef = window.DC.TreeDepth
	}
	var tree_depth = (window.DC.TreeDepth - g.LogDepthRef)

	var text_remaining = text
	for {
		// Split the string. Each new line (after a '\n') is followed by indentation corresponding to the current depth of our log entry.
		// We don't add a trailing \n yet to allow a subsequent item on the same line to be captured.
		var line_start = text_remaining
		var line_end = ImStreolRange(line_start, text_end)
		var is_last_line bool = (line_end == text_end)
		if line_start != line_end || !is_last_line {
			var line_length int = (int)(len(line_end) - len(line_start))
			var indentation int = 1
			if g.LogLineFirstItem {
				indentation = tree_depth * 4
			}
			LogText("%*s%.*s", indentation, "", line_length, line_start)
			g.LogLineFirstItem = false
			if line_end[0] == '\n' {
				LogText(IM_NEWLINE)
				g.LogLineFirstItem = true
			}
		}
		if is_last_line {
			break
		}
		text_remaining = line_end[1:]
	}

	if suffix != "" {
		LogRenderedText(ref_pos, suffix)
	}
}

// Important: doesn't copy underlying data, use carefully (prefix/suffix must be in scope at the time of the next LogRenderedText)
func LogSetNextTextDecoration(prefix string, suffix string) {
	var g = GImGui
	g.LogNextPrefix = prefix
	g.LogNextSuffix = suffix
}

// Logging/Capture
// - All text output from the interface can be captured into tty/file/clipboard. By default, tree nodes are automatically opened during logging.

// start logging to tty (stdout)
func LogToTTY(auto_open_depth int /*= -1*/) {
	var g = GImGui
	if g.LogEnabled {
		return
	}
	LogBegin(ImGuiLogType_TTY, auto_open_depth)
	g.LogFile = os.Stdout
}

// Start logging/capturing text output to given file
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
	var f ImFileHandle = ImFileOpen(filename, "ab")
	if f == nil {
		IM_ASSERT(false)
		return
	}

	LogBegin(ImGuiLogType_File, auto_open_depth)
	g.LogFile = f
}

// start logging to OS clipboard
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
		break
	case ImGuiLogType_File:
		ImFileClose(g.LogFile)
		break
	case ImGuiLogType_Buffer:
		break
	case ImGuiLogType_Clipboard:
		if g.LogBuffer.Len() > 0 {
			SetClipboardText(g.LogBuffer.String())
		}
		break
	case ImGuiLogType_None:
		IM_ASSERT(false)
		break
	}

	g.LogEnabled = false
	g.LogType = ImGuiLogType_None
	g.LogFile = nil
	g.LogBuffer.Reset()
} // stop logging (close file, etc.)

// helper to display buttons for logging to tty/file/clipboard
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

// pass text data straight to log (without being displayed)
func LogText(format string, args ...interface{}) {
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
