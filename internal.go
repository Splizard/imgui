package imgui

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

/*func NewImGuiContext(atlas *ImFontAtlas) *ImGuiContext {
	if atlas == nil {
		atlas = NewFontAtlas()
	}
	var ctx = new(ImGuiContext)
	ctx.IO.Fonts = atlas

	return &ctx
}*/

type ImDrawDataBuilder struct{}         // Helper to build a ImDrawData instance
type ImGuiColorMod struct{}             // Stacked color modifier, backup of modified data so we can restore it
type ImGuiContextHook struct{}          // Hook for extensions like ImGuiTestEngine
type ImGuiDataTypeInfo struct{}         // Type information associated to a ImGuiDataType enum
type ImGuiGroupData struct{}            // Stacked storage data for BeginGroup()/EndGroup()
type ImGuiInputTextState struct{}       // Internal state of the currently focused/edited text input box
type ImGuiLastItemData struct{}         // Status storage for last submitted items
type ImGuiMenuColumns struct{}          // Simple column measurement, currently used for MenuItem() only
type ImGuiNavItemData struct{}          // Result of a gamepad/keyboard directional navigation move query result
type ImGuiMetricsConfig struct{}        // Storage for ShowMetricsWindow() and DebugNodeXXX() functions
type ImGuiNextWindowData struct{}       // Storage for SetNextWindow** functions
type ImGuiNextItemData struct{}         // Storage for SetNextItem** functions
type ImGuiOldColumnData struct{}        // Storage data for a single column for legacy Columns() api
type ImGuiOldColumns struct{}           // Storage data for a columns set for legacy Columns() api
type ImGuiPopupData struct{}            // Storage for current popup stack
type ImGuiSettingsHandler struct{}      // Storage for one type registered in the .ini file
type ImGuiStackSizes struct{}           // Storage of stack sizes for debugging/asserting
type ImGuiStyleMod struct{}             // Stacked style modifier, backup of modified data so we can restore it
type ImGuiTabBar struct{}               // Storage for a tab bar
type ImGuiTabItem struct{}              // Storage for a tab item (within a tab bar)
type ImGuiTable struct{}                // Storage for a table
type ImGuiTableColumn struct{}          // Storage for one column of a table
type ImGuiTableTempData struct{}        // Temporary storage for one table (one per table in the stack), shared between tables.
type ImGuiTableSettings struct{}        // Storage for a table .ini settings
type ImGuiTableColumnsSettings struct{} // Storage for a column .ini settings
type ImGuiWindow struct{}               // Storage for one window
type ImGuiWindowTempData struct{}       // Temporary storage for one window (that's the data which in theory we could ditch at the end of the frame, in practice we currently keep it for each window)
type ImGuiWindowSettings struct{}       // Storage for a window .ini settings (we keep one of those even if the actual window wasn't instanced during this session)

// Use your programming IDE "Go to definition" facility on the names of the center columns to find the actual flags/enum lists.
type ImGuiLayoutType int          // -> enum ImGuiLayoutType_         // Enum: Horizontal or vertical
type ImGuiItemFlags int           // -> enum ImGuiItemFlags_          // Flags: for PushItemFlag()
type ImGuiItemStatusFlags int     // -> enum ImGuiItemStatusFlags_    // Flags: for DC.LastItemStatusFlags
type ImGuiOldColumnFlags int      // -> enum ImGuiOldColumnFlags_     // Flags: for BeginColumns()
type ImGuiNavHighlightFlags int   // -> enum ImGuiNavHighlightFlags_  // Flags: for RenderNavHighlight()
type ImGuiNavDirSourceFlags int   // -> enum ImGuiNavDirSourceFlags_  // Flags: for GetNavInputAmount2d()
type ImGuiNavMoveFlags int        // -> enum ImGuiNavMoveFlags_       // Flags: for navigation requests
type ImGuiNextItemDataFlags int   // -> enum ImGuiNextItemDataFlags_  // Flags: for SetNextItemXXX() functions
type ImGuiNextWindowDataFlags int // -> enum ImGuiNextWindowDataFlags_// Flags: for SetNextWindowXXX() functions
type ImGuiSeparatorFlags int      // -> enum ImGuiSeparatorFlags_     // Flags: for SeparatorEx()
type ImGuiTextFlags int           // -> enum ImGuiTextFlags_          // Flags: for TextEx()
type ImGuiTooltipFlags int        // -> enum ImGuiTooltipFlags_       // Flags: for BeginTooltipEx()

type ImGuiErrorLogCallback func(user_data interface{}, fmt string)

var GImGui *ImGuiContext

func IMGUI_DEBUG_LOG(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%05d] ", GImGui.FrameCount)+format, args...)
}

func IM_ASSERT_USER_ERROR(x bool, msg string) {
	if !x {
		panic(msg)
	}
}

const IM_PI = 3.14159265358979323846
const IM_NEWLINE = "\n"
const IM_TABSIZE = 4

// Unsaturated, for display purpose
func IM_F32_TO_INT8_UNBOUND(val float32) int {
	var x float
	if val >= 0 {
		x = 0.5
	} else {
		x = -0.5
	}
	return (int)(val*255 + x)
}

// Saturated, always output 0..255
func IM_F32_TO_INT8_SAT(val float32) int {
	return (int)(ImSaturate(val)*255.0 + 0.5)
}

func IM_FLOOR(val float32) float {
	return (float)((int)(val))
}

func IM_ROUND(val float32) float {
	return (float)((int)(val + 0.5))
}

func ImHashData(data unsafe.Pointer, data_size uintptr, seed ImU32) ImGuiID { panic("not implemented") }
func ImHashStr(data string, data_size uintptr, seed ImU32) ImGuiID          { panic("not implemented") }

func ImAlphaBlendColors(col_a, col_b ImU32) ImU32 { panic("not implemented") }

func ImIsPowerOfTwoInt(v int) bool {
	return v != 0 && (v&(v-1)) == 0
}

func ImIsPowerOfTwoLong(v int64) bool {
	return v != 0 && (v&(v-1)) == 0
}

func ImUpperPowerOfTwo(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func ImStricmp(str1, str2 string) int                                { panic("not implemented") }
func ImStrncpy(dst []byte, src string, count uintptr)                { panic("not implemented") }
func ImStrdup(str string) []byte                                     { panic("not implemented") }
func ImStrdupcpy(dst []byte, p_dst_size *uintptr, str string) []byte { panic("not implemented") }

func ImStrchrRange(str_begin, str_end string, c char) string { panic("not implemented") }
func ImStrlenW(str string) int                               { panic("not implemented") }
func ImStreolRange(str, str_end string) string               { panic("not implemented") } // End end-of-line
func ImStrbolW(buf_mid_line, buf_begin []ImWchar) []ImWchar  { panic("not implemented") } // Find beginning-of-line
func ImStristr(haystack string, haystack_end string, needle string, needle_end string) string {
	panic("not implemented")
} //
func ImStrTrimBlanks(str []byte) {
	panic("not implemented")
}
func ImStrSkipBlank(str string) string { panic("not implemented") }

func ImFormatString(buf []byte, buf_size uintptr, fmt string, args ...interface{}) int {
	panic("not implemented")
}
func ImFormatStringV(buf []byte, buf_size uintptr, fmt string, args []interface{}) int {
	panic("not implemented")
}
func ImParseFormatFindStart(format string) string { panic("not implemented") }
func ImParseFormatFindEnd(format string) string   { panic("not implemented") }
func ImParseFormatTrimDecorations(format string, buf []byte, buf_size size_t) string {
	panic("not implemented")
}
func ImParseFormatPrecision(format string, default_value int) int { panic("not implemented") }
func ImCharIsBlankA(c char) bool                                  { return c == ' ' || c == '\t' }
func ImCharIsBlankW(c uint) bool                                  { return c == ' ' || c == '\t' || c == 0x3000 }

// Helpers: UTF-8 <> wchar conversions
func ImTextCharToUtf8(out_buf [5]char, c uint) string { panic("not implemented") } // return out_buf
func ImTextStrToUtf8(out_buf []byte, out_buf_size int, in_text []ImWchar, in_text_end []ImWchar) int {
	panic("not implemented")
}                                                                         // return out_// return output UTF-8 bytes count
func ImTextCharFromUtf8(out_char []uint, in_text, in_text_end string) int { panic("not implemented") } // read one character. return input UTF-8 bytes count
func ImTextStrFromUtf8(out_buf []ImWchar, out_buf_size int, in_text, in_text_end string, in_remaining *string) int {
	panic("not implemented")
}                                                                    // return input UTF-8 bytes count
func ImTextCountCharsFromUtf8(in_text, in_text_end string) int       { panic("not implemented") } // return number of UTF-8 code-points (NOT bytes count)
func ImTextCountUtf8BytesFromChar(in_text, in_text_end string) int   { panic("not implemented") } // return number of bytes to express one char in UTF-8
func ImTextCountUtf8BytesFromStr(in_text, in_text_end []ImWchar) int { panic("not implemented") } // return number of bytes to express string in UTF-8

type ImFileHandle = os.File

func ImFileOpen(filename string, mode string) ImFileHandle { panic("not implemented") }
func ImFileClose(file ImFileHandle) bool                   { panic("not implemented") }
func ImFileGetSize(file ImFileHandle) ImU64                { panic("not implemented") }
func ImFileRead(data []byte, size, count ImU64, file ImFileHandle) ImU64 {
	panic("not implemented")
}
func ImFileWrite(data []byte, size, count ImU64, file ImFileHandle) ImU64 {
	panic("not implemented")
}

func ImFileLoadToMemory(filename, mode string, out_file_size *size_t, padding_bytes int) []byte {
	panic("not implemented")
}

func ImBitArrayTestBit(arr []ImU32, n int) bool {
	var mask uint32 = 1 << (n & 31)
	return (arr[n>>5] & mask) != 0
}

func ImBitArrayClearBit(arr []ImU32, n int) {
	var mask uint32 = 1 << (n & 31)
	arr[n>>5] &= ^mask
}

func ImBitArraySetBit(arr []ImU32, n int) {
	var mask uint32 = 1 << (n & 31)
	arr[n>>5] |= mask
}

func ImBitArraySetBitRange(arr []ImU32, n, n2 int) {
	n2--
	for n <= n2 {
		var a_mod int = n & 31
		var b_mod int
		if n2 > (n | 31) {
			b_mod = 31
		} else {
			b_mod = (n2 & 31) + 1
		}
		var mask ImU32 = (1 << b_mod) - 1 & ^(1<<a_mod)
		arr[n>>5] |= mask
		n = (n + 32) & ^31
	}
}

type ImBitVector []ImU32

func (this *ImBitVector) Create(sz int) {
	*this = make([]ImU32, (sz+31)>>5)
}

func (this *ImBitVector) Clear() {
	*this = (*this)[0:]
}

func (this *ImBitVector) TestBit(n int) bool {
	return ImBitArrayTestBit(*this, n)
}

func (this *ImBitVector) SetBit(n int) {
	ImBitArraySetBit(*this, n)
}

func (this *ImBitVector) ClearBit(n int) {
	ImBitArrayClearBit(*this, n)
}

type ImSpan struct {
	Data interface{}
}

func (this ImSpan) Set(data interface{}) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		panic("not implemented")
	}
	this.Data = data
}

func (this ImSpan) Size() int {
	return int32(reflect.ValueOf(this.Data).Len())
}

func (this ImSpan) IndexFromPointer(ptr interface{}) int {
	panic("not implemented")
}
