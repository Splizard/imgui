package imgui

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"unsafe"
)

const FLT_MIN = math.SmallestNonzeroFloat32
const FLT_MAX = math.MaxFloat32
const INT_MAX = math.MaxInt32

type (
	// ImGuiLayoutType Use your programming IDE "Go to definition" facility on the names of the center columns to find the actual flags/enum lists.
	ImGuiLayoutType          int // -> enum ImGuiLayoutType_         // Enum: Horizontal or vertical
	ImGuiItemFlags           int // -> enum ImGuiItemFlags_          // Flags: for PushItemFlag()
	ImGuiItemStatusFlags     int // -> enum ImGuiItemStatusFlags_    // Flags: for DC.LastItemStatusFlags
	ImGuiOldColumnFlags      int // -> enum ImGuiOldColumnFlags_     // Flags: for BeginColumns()
	ImGuiNavHighlightFlags   int // -> enum ImGuiNavHighlightFlags_  // Flags: for RenderNavHighlight()
	ImGuiNavDirSourceFlags   int // -> enum ImGuiNavDirSourceFlags_  // Flags: for GetNavInputAmount2d()
	ImGuiNavMoveFlags        int // -> enum ImGuiNavMoveFlags_       // Flags: for navigation requests
	ImGuiNextItemDataFlags   int // -> enum ImGuiNextItemDataFlags_  // Flags: for SetNextItemXXX() functions
	ImGuiNextWindowDataFlags int // -> enum ImGuiNextWindowDataFlags_// Flags: for SetNextWindowXXX() functions
	ImGuiSeparatorFlags      int // -> enum ImGuiSeparatorFlags_     // Flags: for SeparatorEx()
	ImGuiTextFlags           int // -> enum ImGuiTextFlags_          // Flags: for TextEx()
	ImGuiTooltipFlags        int // -> enum ImGuiTooltipFlags_       // Flags: for BeginTooltipEx()
)

type ImGuiErrorLogCallback func(user_data any, fmt string, args ...any)

// guiContext Current context pointer. Implicitly used by all Dear ImGui functions. Always assumed to be != nil.
//   - ImGui::CreateContext() will automatically set this pointer if it is nil.
//     Change to a different context by calling ImGui::SetCurrentContext().
//   - Important: Dear ImGui functions are not thread-safe because of this pointer.
//     If you want thread-safety to allow N threads to access N different contexts:
//   - Change this variable to use thread local storage so each thread can refer to a different context, in your imconfig.h:
//     struct ImGuiContext;
//     extern thread_local ImGuiContext* MyImGuiTLS;
//     #define guiContext MyImGuiTLS
//     And then define MyImGuiTLS in one of your cpp files. Note that thread_local is a C++11 keyword, earlier C++ uses compiler-specific keyword.
//   - Future development aims to make this context pointer explicit to all calls. Also read https://github.com/ocornut/imgui/issues/586
//   - If you need a finite number of contexts, you may compile and use multiple instances of the ImGui code from a different namespace.
//   - DLL users: read comments above.
var guiContext *ImGuiContext

func IMGUI_DEBUG_LOG(format string, args ...any) {
	fmt.Printf(fmt.Sprintf("[%05d] ", guiContext.FrameCount)+format, args...)
}

func IM_ASSERT_USER_ERROR(x bool, msg string) {
	if !x {
		panic(msg)
	}
}

const IM_PI = 3.14159265358979323846
const IM_NEWLINE = "\n"
const IM_TABSIZE = 4

// IM_F32_TO_INT8_UNBOUND Unsaturated, for display purpose
func IM_F32_TO_INT8_UNBOUND(val float32) int {
	var x float
	if val >= 0 {
		x = 0.5
	} else {
		x = -0.5
	}
	return (int)(val*255 + x)
}

// IM_F32_TO_INT8_SAT Saturated, always output 0..255
func IM_F32_TO_INT8_SAT(val float32) int {
	return (int)(ImSaturate(val)*255.0 + 0.5)
}

func IM_FLOOR(val float32) float { return (float)((int)(val)) }
func IM_ROUND(val float32) float { return (float)((int)(val + 0.5)) }

func ImHashData(ptr unsafe.Pointer, data_size uintptr, seed ImU32) ImGuiID {
	var crc = ^seed
	var data = (*byte)(ptr)
	var crc32_lut = &GCrc32LookupTable
	for i := uintptr(0); i < data_size; i++ {
		crc = (crc >> 8) ^ crc32_lut[(crc&0xFF)^uint(*data)]
		data = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(data)) + 1))
	}
	return ^crc
}

func ImIsPowerOfTwoInt(v int) bool    { return v != 0 && (v&(v-1)) == 0 }
func ImIsPowerOfTwoLong(v int64) bool { return v != 0 && (v&(v-1)) == 0 }

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

func ImCharIsBlankA(c char) bool { return c == ' ' || c == '\t' }
func ImCharIsBlankW(c rune) bool { return c == ' ' || c == '\t' || c == 0x3000 }

// ImTextCharToUtf8_inline Based on stb_to_utf8() from github.com/nothings/stb/
func ImTextCharToUtf8_inline(buf []char, buf_size int, c uint) int {
	if c < 0x80 {
		buf[0] = char(c)
		return 1
	}
	if c < 0x800 {
		if buf_size < 2 {
			return 0
		}
		buf[0] = char(0xc0 + (c >> 6))
		buf[1] = char(0x80 + (c & 0x3f))
		return 2
	}
	if c < 0x10000 {
		if buf_size < 3 {
			return 0
		}
		buf[0] = char(0xe0 + (c >> 12))
		buf[1] = char(0x80 + ((c >> 6) & 0x3f))
		buf[2] = char(0x80 + (c & 0x3f))
		return 3
	}
	if c <= 0x10FFFF {
		if buf_size < 4 {
			return 0
		}
		buf[0] = char(0xf0 + (c >> 18))
		buf[1] = char(0x80 + ((c >> 12) & 0x3f))
		buf[2] = char(0x80 + ((c >> 6) & 0x3f))
		buf[3] = char(0x80 + (c & 0x3f))
		return 4
	}
	// Invalid code point, the max unicode is 0x10FFFF
	return 0
}

// ImTextCharToUtf8 Helpers: UTF-8 <> wchar conversions
func ImTextCharToUtf8(out_buf [5]char, c uint) []char {
	count := ImTextCharToUtf8_inline(out_buf[:], 5, c)
	out_buf[count] = 0
	return out_buf[:]
}

func ImTextStrToUtf8(out_buf []byte, out_buf_size int, in_text []ImWchar, in_text_end []ImWchar) int {
	var out_buf_offset int32 = 0
	var final_buf_offset int32 = 0

	for i := int32(0); i < out_buf_size && i < int32(len(in_text)); i++ {
		c := in_text[i]
		final_buf_offset = i + out_buf_offset
		if c < 0x80 {
			out_buf[final_buf_offset] = char(c)
		} else {
			out_buf_offset += ImTextCharToUtf8_inline(
				out_buf[final_buf_offset:],
				int32(len(out_buf))-final_buf_offset,
				uint32(c),
			)
		}
	}

	return final_buf_offset
}

func ImTextStrFromUtf8(out_buf []ImWchar, out_buf_size int, text string, in_remaining *string) int {
	var count int
	for i, char := range text {
		if count >= int(len(out_buf)) {
			*in_remaining = text[i:]
			return count
		}
		out_buf[i] = char
		count++
	}
	return int(len(text))
}

func ImTextCountCharsFromUtf8(in_text string) int {
	var count int
	for _, val := range in_text {
		var c rune
		ImTextCharFromUtf8(&c, string(val))
		if c == 0 {
			break
		}
		count++
	}
	return count
}

// ImTextCountUtf8BytesFromChar Not optimal but we very rarely use this function.
func ImTextCountUtf8BytesFromChar(in_text, in_text_end []char) int {
	var unused rune = 0
	return ImTextCharFromUtf8(&unused, string(in_text))
}

func ImTextCountUtf8BytesFromStr(in_text, in_text_end []ImWchar) int {
	// return number of bytes to express string in UTF-8
	var bytes_count int = 0
	for _, c := range in_text {
		c_char := char(c)
		if c_char < 0x80 {
			bytes_count++
		} else {
			bytes_count += ImTextCountUtf8BytesFromChar([]char{c_char}, nil)
		}
	}
	return bytes_count
}

type ImFileHandle = *os.File

func ImFileOpen(filename string, mode string) ImFileHandle { panic("not implemented") }
func ImFileClose(file ImFileHandle) bool                   { panic("not implemented") }
func ImFileGetSize(file ImFileHandle) ImU64                { panic("not implemented") }
func ImFileRead(data []byte, size, count ImU64, file ImFileHandle) ImU64 {
	panic("not implemented")
}
func ImFileWrite(data []byte, size, count ImU64, file ImFileHandle) ImU64 {
	panic("not implemented")
}

// ImFileLoadToMemory Helper: Load file content into memory
// Memory allocated with IM_ALLOC(), must be freed by user using IM_FREE() == ImGui::MemFree()
// This can't really be used with "rt" because fseek size won't match read size.
func ImFileLoadToMemory(filename, mode string, out_file_size *size_t, padding_bytes int) []byte {
	b, _ := os.ReadFile(filename)
	return b
}

func ImBitArrayTestBit(arr []ImU32, n int) bool {
	var mask uint32 = 1 << (uint(n) & 31)
	return (arr[n>>5] & mask) != 0
}

func ImBitArrayClearBit(arr []ImU32, n int) {
	var mask uint32 = 1 << (uint(n) & 31)
	arr[n>>5] &= ^mask
}

func ImBitArraySetBit(arr []ImU32, n int) {
	var mask uint32 = 1 << (uint(n) & 31)
	arr[n>>5] |= mask
}

func ImBitArraySetBitRange(arr []ImU32, n, n2 int) {
	n2--
	for n <= n2 {
		var a_mod = n & 31
		var b_mod int
		if n2 > (n | 31) {
			b_mod = 31
		} else {
			b_mod = (n2 & 31) + 1
		}
		var mask ImU32 = (1 << b_mod) - 1&^(1<<a_mod)
		arr[n>>5] |= mask
		n = (n + 32) &^ 31
	}
}

type ImBitVector []ImU32

func (v ImBitVector) SetBitRange(n, n2 int) { // Works on range [n..n2)
	n2--
	for n <= n2 {
		var a_mod = (n & 31)
		var b_mod int
		if n2 > (n | 31) {
			b_mod = 31
		} else {
			b_mod = (n2 & 31) + 1
		}
		var mask = (ImU32)(((ImU64)(1<<b_mod))-1) & ^(ImU32)(((ImU64)(1<<a_mod))-1)
		v[n>>5] |= mask
		n = (n + 32) & ^31
	}
}

func (v *ImBitVector) Create(sz int) {
	*v = make([]ImU32, (uint(sz)+31)>>5)
}

func (v *ImBitVector) Clear() {
	*v = (*v)[:0]
}

func (v *ImBitVector) TestBit(n int) bool {
	return ImBitArrayTestBit(*v, n)
}

func (v *ImBitVector) SetBit(n int) {
	ImBitArraySetBit(*v, n)
}

func (v *ImBitVector) ClearBit(n int) {
	ImBitArrayClearBit(*v, n)
}

type ImSpan struct {
	Data any
}

func (s *ImSpan) Set(data any) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		panic("not implemented")
	}
	s.Data = data
}

func (s *ImSpan) Size() int {
	return int32(reflect.ValueOf(s.Data).Len())
}

func (s *ImSpan) IndexFromPointer(ptr any) int {
	panic("not implemented")
}

type ImDrawListSharedData struct {
	TexUvWhitePixel       ImVec2
	Font                  *ImFont
	FontSize              float
	CurveTessellationTol  float
	CircleSegmentMaxError float
	ClipRectFullscreen    ImVec4
	InitialFlags          ImDrawListFlags
	ArcFastVtx            [IM_DRAWLIST_ARCFAST_TABLE_SIZE]ImVec2
	ArcFastRadiusCutoff   float
	CircleSegmentCounts   [64]ImU8
	TexUvLines            []ImVec4
}

func NewImDrawListSharedData() ImDrawListSharedData {
	var this ImDrawListSharedData
	for i := range this.ArcFastVtx {
		var a = ((float)(i) * 2 * IM_PI) / (float)(len(this.ArcFastVtx))
		this.ArcFastVtx[i] = ImVec2{ImCos(a), ImSin(a)}
	}
	this.ArcFastRadiusCutoff = IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(IM_DRAWLIST_ARCFAST_SAMPLE_MAX, this.CircleSegmentMaxError)
	return this
}

func (d *ImDrawListSharedData) SetCircleTessellationMaxError(max_error float) {
	if d.CircleSegmentMaxError == max_error {
		return
	}
	IM_ASSERT(max_error > 0.0)
	d.CircleSegmentMaxError = max_error
	for i := range d.CircleSegmentCounts {
		var radius = (float)(i)
		if i > 0 {
			d.CircleSegmentCounts[i] = uint8(IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(radius, d.CircleSegmentMaxError))
		}
	}
	d.ArcFastRadiusCutoff = IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(IM_DRAWLIST_ARCFAST_SAMPLE_MAX, d.CircleSegmentMaxError)
}

type ImDrawDataBuilder [2][]*ImDrawList

func (b *ImDrawDataBuilder) Clear() {
	for i := range b {
		b[i] = b[i][:0]
	}
}

func (b *ImDrawDataBuilder) GetDrawListCount() int {
	var count int
	for i := range b {
		count += int(len(b[i]))
	}
	return count
}

func (b *ImDrawDataBuilder) FlattenIntoSingleLayer() {
	var n = len(b[0])
	var size = int(n)
	for i := 1; i < len(*b); i++ {
		size += int(len(b[i]))
	}

	//TODO/FIXME b could be wrong
	if size < int(len(b[0])) {
		b[0] = b[0][0:size]
	} else {
		b[0] = append(b[0], make([]*ImDrawList, 0, size-int(len(b[0])))...)
	}

	for layer_n := 1; layer_n < len(*b); layer_n++ {
		var layer = &b[layer_n]
		if len(*layer) == 0 {
			continue
		}
		copy(b[0][n:], *layer)
		n += len(*layer)
		*layer = (*layer)[:0]
	}
}

type ImGuiDataTypeTempStorage any

// ImGuiDataTypeInfo Type information associated to one ImGuiDataType. Retrieve with DataTypeGetInfo().
type ImGuiDataTypeInfo struct {
	Size     size_t // Size in bytes
	Name     string // Short descriptive name for the type, for debugging
	PrintFmt string // Default printf format for the type
	ScanFmt  string // Default scanf format for the type
}

// ImGuiColorMod Stacked color modifier, backup of modified data so we can restore it
type ImGuiColorMod struct {
	Col         ImGuiCol
	BackupValue ImVec4
}

// ImGuiComboPreviewData Storage data for BeginComboPreview()/EndComboPreview()
type ImGuiComboPreviewData struct {
	PreviewRect                  ImRect
	BackupCursorPos              ImVec2
	BackupCursorMaxPos           ImVec2
	BackupCursorPosPrevLine      ImVec2
	BackupPrevLineTextBaseOffset float
	BackupLayout                 ImGuiLayoutType
}

// ImGuiGroupData Stacked storage data for BeginGroup()/EndGroup()
type ImGuiGroupData struct {
	WindowID                           ImGuiID
	BackupCursorPos                    ImVec2
	BackupCursorMaxPos                 ImVec2
	BackupIndent                       ImVec1
	BackupGroupOffset                  ImVec1
	BackupCurrLineSize                 ImVec2
	BackupCurrLineTextBaseOffset       float
	BackupActiveIdIsAlive              ImGuiID
	BackupActiveIdPreviousFrameIsAlive bool
	BackupHoveredIdIsAlive             bool
	EmitItem                           bool
}

type ImGuiMenuColumns struct {
	TotalWidth     ImU32
	NextTotalWidth ImU32
	Spacing        ImU16
	OffsetIcon     ImU16 // Always zero for now
	OffsetLabel    ImU16 // Offsets are locked in Update()
	OffsetShortcut ImU16
	OffsetMark     ImU16
	Widths         [4]ImU16 // Width of:   Icon, Label, Shortcut, Mark  (accumulators for current frame)
}

func NewImGuiMenuColumns() ImGuiMenuColumns { return ImGuiMenuColumns{} }

// Update Helpers for internal use
func (c *ImGuiMenuColumns) Update(spacing float, window_reappearing bool) {
	if window_reappearing {
		c.Widths = [4]ImU16{}
	}
	c.Spacing = (ImU16)(spacing)
	c.CalcNextTotalWidth(true)
	c.Widths = [4]ImU16{}
	c.TotalWidth = c.NextTotalWidth
	c.NextTotalWidth = 0
}

func (c *ImGuiMenuColumns) DeclColumns(w_icon float, w_label float, w_shortcut float, w_mark float) float {
	c.Widths[0] = uint16(max(int(c.Widths[0]), ((int)(w_icon))))
	c.Widths[1] = uint16(max(int(c.Widths[1]), (int)(w_label)))
	c.Widths[2] = uint16(max(int(c.Widths[2]), (int)(w_shortcut)))
	c.Widths[3] = uint16(max(int(c.Widths[3]), (int)(w_mark)))
	c.CalcNextTotalWidth(false)
	return (float)(max(int(c.TotalWidth), int(c.NextTotalWidth)))
}

func (c *ImGuiMenuColumns) CalcNextTotalWidth(update_offsets bool) {
	var offset ImU16 = 0
	var want_spacing = false
	for i := 0; i < len(c.Widths); i++ {
		var width = c.Widths[i]
		if want_spacing && width > 0 {
			offset += c.Spacing
		}
		want_spacing = want_spacing || (width > 0)
		if update_offsets {
			if i == 1 {
				c.OffsetLabel = offset
			}
			if i == 2 {
				c.OffsetShortcut = offset
			}
			if i == 3 {
				c.OffsetMark = offset
			}
		}
		offset += width
	}
	c.NextTotalWidth = uint(offset)
}

// ImGuiInputTextState Internal state of the currently focused/edited text input box
// For a given item ID, access with ImGui::GetInputTextState()
type ImGuiInputTextState struct {
	ID                   ImGuiID   // widget id owning the text state
	CurLenW, CurLenA     int       // we need to maintain our buffer length in both UTF-8 and wchar format. UTF-8 length is valid even if TextA is not.
	TextW                []ImWchar // edit buffer, we need to persist but can't guarantee the persistence of the user-provided buffer. so we copy into own buffer.
	TextA                []char    // temporary UTF8 buffer for callbacks and other operations. this is not updated in every code-path! size=capacity.
	InitialTextA         []char    // backup of end-user buffer at the time of focus (in UTF-8, unaltered)
	TextAIsValid         bool      // temporary UTF8 buffer is not initially valid before we make the widget active (until then we pull the data from user argument)
	BufCapacityA         int       // end-user buffer capacity
	ScrollX              float     // horizontal scrolling/offset
	Stb                  STB_TexteditState
	CursorAnim           float // timer for cursor blink, reset on every user action so the cursor reappears immediately
	CursorFollow         bool  // set when we want scrolling to follow the current cursor position (not always!)
	SelectedAllMouseLock bool  // after a double-click to select all, we ignore further mouse drags to update selection
	Edited               bool  // edited this frame
	Flags                ImGuiInputTextFlags
	UserCallback         ImGuiInputTextCallback
	UserCallbackData     any
}

func (s *ImGuiInputTextState) ClearText() {
	s.CurLenW = 0
	s.CurLenA = 0
	s.TextW[0] = 0
	s.TextA[0] = 0
	s.CursorClamp()
}

func (s *ImGuiInputTextState) GetUndoAvailCount() int {
	return int(s.Stb.undostate.undo_point)
}

func (s *ImGuiInputTextState) GetRedoAvailCount() int {
	return int(STB_TEXTEDIT_UNDOSTATECOUNT - s.Stb.undostate.redo_point)
}

func (s *ImGuiInputTextState) OnKeyPressed(key int) {
	stb_textedit_key(s, &s.Stb, STB_TEXTEDIT_KEYTYPE(key))
	s.CursorFollow = true
	s.CursorAnimReset()
}

func (s *ImGuiInputTextState) CursorAnimReset() {
	s.CursorAnim = -0.30
}

func (s *ImGuiInputTextState) CursorClamp() {
	s.Stb.cursor = min(s.Stb.cursor, s.CurLenW)
	s.Stb.select_start = min(s.Stb.select_start, s.CurLenW)
	s.Stb.select_end = min(s.Stb.select_end, s.CurLenW)
}

func (s *ImGuiInputTextState) HasSelection() bool {
	return s.Stb.select_start != s.Stb.select_end
}

func (s *ImGuiInputTextState) ClearSelection() {
	s.Stb.select_start = s.Stb.cursor
	s.Stb.select_end = s.Stb.cursor
}

func (s *ImGuiInputTextState) GetCursorPos() int {
	return s.Stb.cursor
}

func (s *ImGuiInputTextState) GetSelectionStart() int {
	return s.Stb.select_start
}

func (s *ImGuiInputTextState) GetSelectionEnd() int {
	return s.Stb.select_end
}

func (s *ImGuiInputTextState) SelectAll() {
	s.Stb.select_start = 0
	s.Stb.cursor = s.CurLenW
	s.Stb.select_end = s.CurLenW
	s.Stb.has_preferred_x = 0
}

type ImGuiPopupData struct {
	PopupId        ImGuiID      // Set on OpenPopup()
	Window         *ImGuiWindow // Resolved on BeginPopup() - may stay unresolved if user never calls OpenPopup()
	SourceWindow   *ImGuiWindow // Set on OpenPopup() copy of NavWindow at the time of opening the popup
	OpenFrameCount int          // Set on OpenPopup()
	OpenParentId   ImGuiID      // Set on OpenPopup(), we need this to differentiate multiple menu sets from each others (e.guiContext. inside menu bar vs loose menu items)
	OpenPopupPos   ImVec2       // Set on OpenPopup(), preferred popup position (typically == OpenMousePos when using mouse)
	OpenMousePos   ImVec2       // Set on OpenPopup(), copy of mouse position at the time of opening popup
}

type ImGuiNextWindowData struct {
	Flags                ImGuiNextWindowDataFlags
	PosCond              ImGuiCond
	SizeCond             ImGuiCond
	CollapsedCond        ImGuiCond
	PosVal               ImVec2
	PosPivotVal          ImVec2
	SizeVal              ImVec2
	ContentSizeVal       ImVec2
	ScrollVal            ImVec2
	CollapsedVal         bool
	SizeConstraintRect   ImRect
	SizeCallback         ImGuiSizeCallback
	SizeCallbackUserData any
	BgAlphaVal           float
	MenuBarOffsetMinVal  ImVec2
}

func (d *ImGuiNextWindowData) ClearFlags() {
	d.Flags = ImGuiNextWindowDataFlags_None
}

type ImGuiNextItemData struct {
	Flags        ImGuiNextItemDataFlags
	Width        float     // Set by SetNextItemWidth()
	FocusScopeId ImGuiID   // Set by SetNextItemMultiSelectData() (!= 0 signify value has been set, so it's an alternate version of HasSelectionData, we don't use Flags for this because they are cleared too early. This is mostly used for debugging)
	OpenCond     ImGuiCond // Set by SetNextItemOpen()
	OpenVal      bool      // Set by SetNextItemOpen()
}

func (d *ImGuiNextItemData) ClearFlags() {
	d.Flags = ImGuiNextItemDataFlags_None
}

type ImGuiLastItemData struct {
	ID          ImGuiID
	InFlags     ImGuiItemFlags       // See ImGuiItemFlags_
	StatusFlags ImGuiItemStatusFlags // See ImGuiItemStatusFlags_
	Rect        ImRect               // Full rectangle
	NavRect     ImRect               // Navigation scoring rectangle (not displayed)
	DisplayRect ImRect               // Display rectangle (only if ImGuiItemStatusFlags_HasDisplayRect is set)
}

// ImGuiWindowStackData Data saved for each window pushed into the stack
type ImGuiWindowStackData struct {
	Window                   *ImGuiWindow
	ParentLastItemDataBackup ImGuiLastItemData
}

type ImGuiShrinkWidthItem struct {
	Index int
	Width float
}

type ImGuiPtrOrIndex struct {
	Ptr   any // Either field can be set, not both. e.guiContext. Dock node tab bars are loose while BeginTabBar() ones are in a pool.
	Index int // Usually index in a main pool.
}

func ImGuiPtr(ptr any) ImGuiPtrOrIndex {
	return ImGuiPtrOrIndex{Ptr: ptr}
}

func ImGuiIndex(index int) ImGuiPtrOrIndex {
	return ImGuiPtrOrIndex{Index: index}
}

type ImGuiNavItemData struct {
	Window       *ImGuiWindow // Init,Move    // Best candidate window (result->ItemWindow->RootWindowForNav == request->Window)
	ID           ImGuiID      // Init,Move    // Best candidate item ID
	FocusScopeId ImGuiID      // Init,Move    // Best candidate focus scope ID
	RectRel      ImRect       // Init,Move    // Best candidate bounding box in window relative space
	DistBox      float        //      Move    // Best candidate box distance to current NavId
	DistCenter   float        //      Move    // Best candidate center distance to current NavId
	DistAxial    float        //      Move    // Best candidate axial distance to current NavId
}

func NewImGuiNavItemData() ImGuiNavItemData {
	return ImGuiNavItemData{
		DistBox:    FLT_MAX,
		DistCenter: FLT_MAX,
		DistAxial:  FLT_MAX,
	}
}

func (d *ImGuiNavItemData) Clear() {
	d.DistBox = FLT_MAX
	d.DistCenter = FLT_MAX
	d.DistAxial = FLT_MAX
}

// ImGuiOldColumnData Storage data for a single column for legacy Columns() api
type ImGuiOldColumnData struct {
	OffsetNorm             float // Column start offset, normalized 0.0 (far left) -> 1.0 (far right)
	OffsetNormBeforeResize float
	Flags                  ImGuiOldColumnFlags // Not exposed
	ClipRect               ImRect
}

type ImGuiOldColumns struct {
	ID                       ImGuiID
	Flags                    ImGuiOldColumnFlags
	IsFirstFrame             bool
	IsBeingResized           bool
	Current                  int
	Count                    int
	OffMinX, OffMaxX         float // Offsets from HostWorkRect.Min.x
	LineMinY, LineMaxY       float
	HostCursorPosY           float  // Backup of CursorPos at the time of BeginColumns()
	HostCursorMaxPosX        float  // Backup of CursorMaxPos at the time of BeginColumns()
	HostInitialClipRect      ImRect // Backup of ClipRect at the time of BeginColumns()
	HostBackupClipRect       ImRect // Backup of ClipRect during PushColumnsBackground()/PopColumnsBackground()
	HostBackupParentWorkRect ImRect // Backup of WorkRect at the time of BeginColumns()
	Columns                  []ImGuiOldColumnData
	Splitter                 ImDrawListSplitter
}

// ImGuiViewportP ImGuiViewport Private/Internals fields (cardinal sin: we are using inheritance!)
// Every instance of ImGuiViewport is in fact a ImGuiViewportP.
type ImGuiViewportP = ImGuiViewport

func NewImGuiViewportP() ImGuiViewportP {
	return ImGuiViewportP{
		DrawListsLastFrame: [2]int{-1, -1},
	}
}

func (p *ImGuiViewportP) CalcWorkRectPos(off_min *ImVec2) ImVec2 {
	return ImVec2{p.Pos.x + off_min.x, p.Pos.y + off_min.y}
}

func (p *ImGuiViewportP) CalcWorkRectSize(off_min *ImVec2, off_max *ImVec2) ImVec2 {
	return ImVec2{max(0.0, p.Size.x-off_min.x+off_max.x), max(0.0, p.Size.y-off_min.y+off_max.y)}
}

func (p *ImGuiViewportP) UpdateWorkRect() {
	p.WorkPos = p.CalcWorkRectPos(&p.WorkOffsetMin)
	p.WorkSize = p.CalcWorkRectSize(&p.WorkOffsetMin, &p.WorkOffsetMax)
}

func (p *ImGuiViewportP) GetMainRect() ImRect {
	return ImRect{ImVec2{p.Pos.x, p.Pos.y}, ImVec2{p.Pos.x + p.Size.x, p.Pos.y + p.Size.y}}
}

func (p *ImGuiViewportP) GetWorkRect() ImRect {
	return ImRect{ImVec2{p.WorkPos.x, p.WorkPos.y}, ImVec2{p.WorkPos.x + p.WorkSize.x, p.WorkPos.y + p.WorkSize.y}}
}

func (p *ImGuiViewportP) GetBuildWorkRect() ImRect {
	var pos = p.CalcWorkRectPos(&p.BuildWorkOffsetMin)
	var size = p.CalcWorkRectSize(&p.BuildWorkOffsetMin, &p.BuildWorkOffsetMax)
	return ImRect{ImVec2{pos.x, pos.y}, ImVec2{pos.x + size.x, pos.y + size.y}}
}

// ImGuiWindowSettings Windows data saved in imgui.ini file
// Because we never destroy or rename ImGuiWindowSettings, we can store the names in a separate buffer easily.
// (this is designed to be stored in a ImChunkStream buffer, with the variable-length Name following our structure)
type ImGuiWindowSettings struct {
	ID        ImGuiID
	Pos       ImVec2ih
	Size      ImVec2ih
	Collapsed bool
	WantApply bool // Set when loaded from .ini data (to enable merging/loading .ini data into an already running context)
	name      string
}

func (s *ImGuiWindowSettings) GetName() string {
	return s.name
}

type ImGuiSettingsHandler struct {
	TypeName   string // Short description stored in .ini file. Disallowed characters: '[' ']'
	TypeHash   ImGuiID
	ClearAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                           // Clear all settings data
	ReadInitFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                           // Read: Called before reading (in registration order)
	ReadOpenFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, name string) any          // Read: Called when entering into a new ini entry e.guiContext. "[Window][Name]"
	ReadLineFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, entry any, line string)   // Read: Called for every line of text within an ini entry
	ApplyAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                           // Read: Called after reading (in registration order)
	WriteAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, out_buf *ImGuiTextBuffer) // Write: Output every entries into 'out_buf'
	UserData   any
}

type ImGuiMetricsConfig struct {
	ShowWindowsRects         bool
	ShowWindowsBeginOrder    bool
	ShowTablesRects          bool
	ShowDrawCmdMesh          bool
	ShowDrawCmdBoundingBoxes bool
	ShowWindowsRectsType     int
	ShowTablesRectsType      int
}

func NewImGuiMetricsConfig() ImGuiMetricsConfig {
	return ImGuiMetricsConfig{
		ShowDrawCmdMesh:          true,
		ShowDrawCmdBoundingBoxes: true,
		ShowWindowsRectsType:     -1,
		ShowTablesRectsType:      -1,
	}
}

type ImGuiStackSizes struct {
	SizeOfIDStack         short
	SizeOfColorStack      short
	SizeOfStyleVarStack   short
	SizeOfFontStack       short
	SizeOfFocusScopeStack short
	SizeOfGroupStack      short
	SizeOfBeginPopupStack short
}

func (s *ImGuiStackSizes) SetToCurrentState() {
	window := guiContext.CurrentWindow
	s.SizeOfIDStack = (short)(len(window.IDStack))
	s.SizeOfColorStack = (short)(len(guiContext.ColorStack))
	s.SizeOfStyleVarStack = (short)(len(guiContext.StyleVarStack))
	s.SizeOfFontStack = (short)(len(guiContext.FontStack))
	s.SizeOfFocusScopeStack = (short)(len(guiContext.FocusScopeStack))
	s.SizeOfGroupStack = (short)(len(guiContext.GroupStack))
	s.SizeOfBeginPopupStack = (short)(len(guiContext.BeginPopupStack))
}

func (s *ImGuiStackSizes) CompareWithCurrentState() {
	window := guiContext.CurrentWindow

	// Window stacks
	// NOT checking: DC.ItemWidth, DC.TextWrapPos (per window) to allow user to conveniently push once and not pop (they are cleared on Begin)
	IM_ASSERT_USER_ERROR(s.SizeOfIDStack == short(len(window.IDStack)), "PushID/PopID or TreeNode/TreePop Mismatch!")

	// Global stacks
	// For color, style and font stacks there is an incentive to use Push/Begin/Pop/.../End patterns, so we relax our checks a little to allow them.
	IM_ASSERT_USER_ERROR(s.SizeOfGroupStack == short(len(guiContext.GroupStack)), "BeginGroup/EndGroup Mismatch!")
	IM_ASSERT_USER_ERROR(s.SizeOfBeginPopupStack == short(len(guiContext.BeginPopupStack)), "BeginPopup/EndPopup or BeginMenu/EndMenu Mismatch!")
	IM_ASSERT_USER_ERROR(s.SizeOfColorStack >= short(len(guiContext.ColorStack)), "PushStyleColor/PopStyleColor Mismatch!")
	IM_ASSERT_USER_ERROR(s.SizeOfStyleVarStack >= short(len(guiContext.StyleVarStack)), "PushStyleVar/PopStyleVar Mismatch!")
	IM_ASSERT_USER_ERROR(s.SizeOfFontStack >= short(len(guiContext.FontStack)), "PushFont/PopFont Mismatch!")
	IM_ASSERT_USER_ERROR(s.SizeOfFocusScopeStack == short(len(guiContext.FocusScopeStack)), "PushFocusScope/PopFocusScope Mismatch!")
}

type ImGuiContextHookCallback func(ctx *ImGuiContext, hook *ImGuiContextHook)

type ImGuiContextHook struct {
	HookId   ImGuiID // A unique ID assigned by AddContextHook()
	Type     ImGuiContextHookType
	Owner    ImGuiID
	Callback ImGuiContextHookCallback
	UserData any
}

// ImGuiWindowTempData Transient per-window data, reset at the beginning of the frame. This used to be called ImGuiDrawContext, hence the DC variable name in ImGuiWindow.
// (That's theory, in practice the delimitation between ImGuiWindow and ImGuiWindowTempData is quite tenuous and could be reconsidered..)
// (This doesn't need a constructor because we zero-clear it as part of ImGuiWindow and all frame-temporary data are setup on Begin)
type ImGuiWindowTempData struct {
	// Layout
	CursorPos              ImVec2  // Current emitting position, in absolute coordinates.
	CursorPosPrevLine      ImVec2  // Previous line emitting position
	CursorStartPos         ImVec2  // Initial position after Begin(), generally == window->Pos
	CursorMaxPos           ImVec2  // Used to implicitly calculate ContentSize at the beginning of next frame, for auto-resize
	IdealMaxPos            ImVec2  // Implicitly calculate the size of our contents, always extending. Saved/restored as ImVec2(FLT_MAX, FLT_MAX) is we don't fit in those constraints.
	CurrLineSize           ImVec2  // Current line size
	PrevLineSize           ImVec2  // Size of previous line
	CurrLineTextBaseOffset float32 // Baseline offset from top of line to top of text after rescaling for text height
	PrevLineTextBaseOffset float32 // Baseline offset from top of line to top of text before rescaling for text height
	Indent                 ImVec1  // Indentation / start position from left of window (increased by TreePush/TreePop, etc.)
	ColumnsOffset          ImVec1  // Offset to the current column (if ColumnsCurrent > 0). FIXME: This and the above should be a stack to allow use cases like Tree->Column->Tree. Need revamp columns API.
	GroupOffset            ImVec1  // Offset to the current group (if any)
	// Keyboard/Gamepad navigation
	NavLayerCurrent          ImGuiNavLayer // Current layer, 0..31 (we currently only use 0..1)
	NavLayersActiveMask      ImU32         // Which layers have been written to (result from previous frame)
	NavLayersActiveMaskNext  ImU32         // Which layers have been written to (accumulator for current frame)
	NavFocusScopeIdCurrent   ImGuiID       // Current focus scope ID while appending
	NavHideHighlightOneFrame bool          // Hide highlight for one frame
	NavHasScroll             bool          // Set when scrolling can be used (ScrollMax > 0.0f)
	// Miscellaneous
	MenuBarAppending          bool             // FIXME: Remove this
	MenuBarOffset             ImVec2           // MenuBarOffset.x is sort of equivalent of a per-layer CursorPos.x, saved/restored as we switch to the menu bar. The only situation when MenuBarOffset.y is > 0 if when (SafeAreaPadding.y > FramePadding.y), often used on TVs.
	MenuColumns               ImGuiMenuColumns // Simplified columns storage for menu items measurement
	TreeDepth                 int              // Current tree depth.
	TreeJumpToParentOnPopMask ImU32            // Store a copy of !guiContext.NavIdIsAlive for TreeDepth 0..31.. Could be turned into a ImU64 if necessary.
	ChildWindows              []*ImGuiWindow
	StateStorage              ImGuiStorage
	CurrentColumns            *ImGuiOldColumns
	CurrentTableIdx           int
	LayoutType                ImGuiLayoutType
	ParentLayoutType          ImGuiLayoutType
	FocusCounterRegular       int
	FocusCounterTabStop       int
	// Local parameters stacks
	// We store the current settings outside of the vectors to increase memory locality (reduce cache misses). The vectors are rarely modified. Also it allows us to not heap allocate for short-lived windows which are not using those settings.
	ItemWidth         float // Current item width (>0.0: width in pixels, <0.0: align xx pixels to the right of window)
	TextWrapPos       float // Current text wrap position
	ItemWidthStack    []float
	TextWrapPosStack  []float
	StackSizesOnBegin ImGuiStackSizes
}

type ImGuiWindow struct {
	Name                string           // Window name, owned by the window.
	ID                  ImGuiID          // == ImHashStr(Name)
	Flags               ImGuiWindowFlags // See enum ImGuiWindowFlags_
	Pos                 ImVec2           // Position (always rounded-up to nearest pixel)
	Size                ImVec2           // Current size (==SizeFull or collapsed title bar size)
	SizeFull            ImVec2           // Size when non collapsed
	ContentSize         ImVec2           // Size of contents/scrollable client area (calculated from the extents reach of the cursor) from previous frame. Does not include window decoration or window padding.
	ContentSizeIdeal    ImVec2
	ContentSizeExplicit ImVec2  // Size of contents/scrollable client area explicitly request by the user via SetNextWindowContentSize().
	WindowPadding       ImVec2  // Window padding at the time of begin.
	WindowRounding      float   // Window rounding at the time of Begin(). May be clamped lower to avoid rendering artifacts with title bar, menu bar etc.
	WindowBorderSize    float   // Window border size at the time of begin.
	MoveId              ImGuiID // == window->GetID("#MOVE")
	ChildId             ImGuiID // ID of corresponding item in parent window (for navigation to return from child window to parent window)
	Scroll              ImVec2  // Current scrolling position
	ScrollMax           ImVec2  // Scrollable maximum position
	ScrollTarget        ImVec2  // target scroll position. stored as cursor position with scrolling canceled out, so the highest point is always 0.0f. (FLT_MAX for no change)

	ScrollTargetCenterRatio        ImVec2 // 0.0f = scroll so that target position is at top, 0.5f = scroll so that target position is centered
	ScrollTargetEdgeSnapDist       ImVec2 // 0.0f = no snapping, >0.0f snapping threshold
	ScrollbarSizes                 ImVec2 // Size taken by each scrollbars on their smaller axis. Pay attention! ScrollbarSizes.x == width of the vertical scrollbar, ScrollbarSizes.y = height of the horizontal scrollbar.
	ScrollbarX, ScrollbarY         bool   // Are scrollbars visible?
	Active                         bool   // Set to true on Begin(), unless Collapsed
	WasActive                      bool
	WriteAccessed                  bool // Set to true when any widget access the current window
	Collapsed                      bool // Set when collapsing window to become only title-bar
	WantCollapseToggle             bool
	SkipItems                      bool    // Set when items can safely be all clipped (e.guiContext. window not visible or collapsed)
	Appearing                      bool    // Set during the frame where the window is appearing (or re-appearing)
	Hidden                         bool    // Do not display (== HiddenFrames*** > 0)
	IsFallbackWindow               bool    // Set on the "Debug##Default" window.
	HasCloseButton                 bool    // Set when the window has a close button (p_open != NULL)
	ResizeBorderHeld               int8    // Current border being held for resize (-1: none, otherwise 0-3)
	BeginCount                     short   // Number of Begin() during the current frame (generally 0 or 1, 1+ if appending via multiple Begin/End pairs)
	BeginOrderWithinParent         short   // Begin() order within immediate parent window, if we are a child window. Otherwise 0.
	BeginOrderWithinContext        short   // Begin() order within entire imgui context. This is mostly used for debugging submission order related issues.
	FocusOrder                     short   // Order within WindowsFocusOrder[], altered when windows are focused.
	PopupId                        ImGuiID // ID in the popup stack when this window is used as a popup/menu (because we use generic Name/ID for recycling)
	AutoFitFramesX, AutoFitFramesY ImS8
	AutoFitChildAxises             ImS8
	AutoFitOnlyGrows               bool
	AutoPosLastDirection           ImGuiDir
	HiddenFramesCanSkipItems       ImS8      // Hide the window for N frames
	HiddenFramesCannotSkipItems    ImS8      // Hide the window for N frames while allowing items to be submitted so we can measure their size
	HiddenFramesForRenderOnly      ImS8      // Hide the window until frame N at Render() time only
	DisableInputsFrames            ImS8      // Disable window interactions for N frames
	SetWindowPosAllowFlags         ImGuiCond //8 :         // store acceptable condition flags for SetNextWindowPos() use.
	SetWindowSizeAllowFlags        ImGuiCond //8 :        // store acceptable condition flags for SetNextWindowSize() use.
	SetWindowCollapsedAllowFlags   ImGuiCond //8:   // store acceptable condition flags for SetNextWindowCollapsed() use.
	SetWindowPosVal                ImVec2    // store window position when using a non-zero Pivot (position set needs to be processed when we know the window size)
	SetWindowPosPivot              ImVec2    // store window pivot for positioning. ImVec2(0, 0) when positioning fromcorner top-left ImVec2(0.5f, 0.5f)centering for ImVec2(1, 1) for bottom right.

	IDStack []ImGuiID           // ID stack. ID are hashes seeded with the value at the top of the stack. (In theory this should be in the TempData structure)
	DC      ImGuiWindowTempData // Temporary per-window data, reset at the beginning of the frame. This used to be called ImGuiDrawContext, hence the "DC" variable name.

	// The best way to understand what those rectangles are is to use the 'Metrics->Tools->Show Windows Rectangles' viewer.
	// The main 'OuterRect', omitted as a field, is window->Rect().
	OuterRectClipped  ImRect   // == Window->Rect() just after setup in Begin(). == window->Rect() for root window.
	InnerRect         ImRect   // Inner rectangle (omit title bar, menu bar, scroll bar)
	InnerClipRect     ImRect   // == InnerRect shrunk by WindowPadding*0.5f on each side, clipped within viewport or parent clip rect.
	WorkRect          ImRect   // Initially covers the whole scrolling region. Reduced by containers e.guiContext columns/tables when active. Shrunk by WindowPadding*1.0f on each side. This is meant to replace ContentRegionRect over time (from 1.71+ onward).
	ParentWorkRect    ImRect   // Backup of WorkRect before entering a container such as columns/tables. Used by e.guiContext. SpanAllColumns functions to easily access. Stacked containers are responsible for maintaining this. // FIXME-WORKRECT: Could be a stack?
	ClipRect          ImRect   // Current clipping/scissoring rectangle, evolve as we are using PushClipRect(), etc. == DrawList->clip_rect_stack.back().
	ContentRegionRect ImRect   // FIXME: This is currently confusing/misleading. It is essentially WorkRect but not handling of scrolling. We currently rely on it as right/bottom aligned sizing operation need some size to rely on.
	HitTestHoleSize   ImVec2ih // Define an optional rectangular hole where mouse will pass-through the window.
	HitTestHoleOffset ImVec2ih

	LastFrameActive  int   // Last frame number the window was Active.
	LastTimeActive   float // Last timestamp the window was Active (using float as we don't need high precision there)
	ItemWidthDefault float
	StateStorage     ImGuiStorage
	ColumnsStorage   []ImGuiOldColumns
	FontWindowScale  float // User scale multiplier per-window, via SetWindowFontScale()
	SettingsOffset   int   // Offset into SettingsWindows[] (offsets are always valid as we only grow the array from the back)

	DrawList                       *ImDrawList // == &DrawListInst (for backward compatibility reason with code using imgui_internal.h we keep this a pointer)
	DrawListInst                   ImDrawList
	ParentWindow                   *ImGuiWindow // If we are a child _or_ popup window, this is pointing to our parent. Otherwise NULL.
	RootWindow                     *ImGuiWindow // Point to ourself or first ancestor that is not a child window == Top-level window.
	RootWindowForTitleBarHighlight *ImGuiWindow // Point to ourself or first ancestor which will display TitleBgActive color when this window is active.
	RootWindowForNav               *ImGuiWindow // Point to ourself or first ancestor which doesn't have the NavFlattened flag.

	NavLastChildNavWindow *ImGuiWindow                 // When going to the menu bar, we remember the child window we came from. (This could probably be made implicit if we kept guiContext.Windows sorted by last focused including child window.)
	NavLastIds            [ImGuiNavLayer_COUNT]ImGuiID // Last known NavId for this window, per layer (0/1)
	NavRectRel            [ImGuiNavLayer_COUNT]ImRect  // Reference rectangle, in window relative space

	MemoryDrawListIdxCapacity int // Backup of last idx/vtx count, so when waking up the window we can preallocate and avoid iterative alloc/copy
	MemoryDrawListVtxCapacity int
	MemoryCompacted           bool // Set when window extraneous data have been garbage collected
}

func NewImGuiWindow(context *ImGuiContext, name string) *ImGuiWindow {
	var id = ImHashStr(name, 0, 0)
	var window = ImGuiWindow{
		Name:                         name,
		ID:                           id,
		IDStack:                      []ImGuiID{id},
		ScrollTarget:                 ImVec2{FLT_MAX, FLT_MAX},
		ScrollTargetCenterRatio:      ImVec2{0.5, 0.5},
		AutoFitFramesX:               -1,
		AutoFitFramesY:               -1,
		AutoPosLastDirection:         ImGuiDir_None,
		SetWindowPosAllowFlags:       ImGuiCond_Always | ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing,
		SetWindowSizeAllowFlags:      ImGuiCond_Always | ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing,
		SetWindowCollapsedAllowFlags: ImGuiCond_Always | ImGuiCond_Once | ImGuiCond_FirstUseEver | ImGuiCond_Appearing,
		SetWindowPosVal:              ImVec2{FLT_MAX, FLT_MAX},
		SetWindowPosPivot:            ImVec2{FLT_MAX, FLT_MAX},
		LastFrameActive:              -1,
		LastTimeActive:               -1.0,
		FontWindowScale:              1.0,
		SettingsOffset:               -1,
	}
	window.MoveId = window.GetIDs("#MOVE")
	window.DrawList = &window.DrawListInst
	window.DrawList._Data = &context.DrawListSharedData
	window.DrawList._OwnerName = name
	return &window
}

func (w *ImGuiWindow) GetIDs(str string) ImGuiID {
	var seed = w.IDStack[len(w.IDStack)-1]
	var id = ImHashStr(str, 0, seed)
	KeepAliveID(id)
	return id
}

func (w *ImGuiWindow) GetIDInterface(ptr any) ImGuiID {
	rvalue := reflect.ValueOf(ptr)

	// .Elem() will panic if it's not an interface or a pointer
	if rvalue.Kind() == reflect.Interface || rvalue.Kind() == reflect.Pointer {
		rvalue = rvalue.Elem()
	}

	// If we can't get the address of the value, make it a reference
	if !rvalue.CanAddr() {
		rvalue = reflect.ValueOf(&ptr).Elem()
	}

	var seed = w.IDStack[len(w.IDStack)-1]
	var id = ImHashData(unsafe.Pointer(rvalue.UnsafeAddr()), rvalue.Type().Size(), seed)
	KeepAliveID(id)
	return id
}

func (w *ImGuiWindow) GetIDInt(n int) ImGuiID {
	var seed = w.IDStack[len(w.IDStack)-1]
	var id = ImHashData(unsafe.Pointer(&n), unsafe.Sizeof(n), seed)
	KeepAliveID(id)
	return id
}

func (w *ImGuiWindow) GetIDNoKeepAlive(str string) ImGuiID {
	var seed = w.IDStack[len(w.IDStack)-1]
	var id = ImHashStr(str, 0, seed)
	return id
}

func (w *ImGuiWindow) GetIDNoKeepAliveInterface(ptr any) ImGuiID {
	rvalue := reflect.ValueOf(ptr).Elem()
	var seed = w.IDStack[len(w.IDStack)-1]
	var id = ImHashData(unsafe.Pointer(rvalue.UnsafeAddr()), rvalue.Type().Size(), seed)
	return id
}

func (w *ImGuiWindow) GetIDNoKeepAliveInt(n int) ImGuiID {
	var seed = w.IDStack[len(w.IDStack)-1]
	var id = ImHashData(unsafe.Pointer(&n), unsafe.Sizeof(n), seed)
	return id
}

// GetIDFromRectangle This is only used in rare/specific situations to manufacture an ID out of nowhere.
func (w *ImGuiWindow) GetIDFromRectangle(r_abs ImRect) ImGuiID {
	var seed = w.IDStack[len(w.IDStack)-1]
	var r_rel = [4]int{(int)(r_abs.Min.x - w.Pos.x), (int)(r_abs.Min.y - w.Pos.y), (int)(r_abs.Max.x - w.Pos.x), (int)(r_abs.Max.y - w.Pos.y)}
	var id = ImHashData(unsafe.Pointer(&r_rel), unsafe.Sizeof(r_rel), seed)
	KeepAliveID(id)
	return id
}

func (w *ImGuiWindow) Rect() ImRect {
	return ImRect{ImVec2{w.Pos.x, w.Pos.y}, ImVec2{w.Pos.x + w.Size.x, w.Pos.y + w.Size.y}}
}

func (w *ImGuiWindow) CalcFontSize() float {
	var scale = guiContext.FontBaseSize * w.FontWindowScale
	if w.ParentWindow != nil {
		scale *= w.ParentWindow.FontWindowScale
	}
	//return 20 //TODO/FIXME
	return scale
}

func (w *ImGuiWindow) TitleBarHeight() float {
	if w.Flags&ImGuiWindowFlags_NoTitleBar != 0 {
		return 0.0
	}
	return w.CalcFontSize() + guiContext.Style.FramePadding.y*2.0
}

func (w *ImGuiWindow) TitleBarRect() ImRect {
	return ImRect{ImVec2{w.Pos.x, w.Pos.y}, ImVec2{w.Pos.x + w.SizeFull.x, w.Pos.y + w.TitleBarHeight()}}
}

func (w *ImGuiWindow) MenuBarHeight() float {
	if w.Flags&ImGuiWindowFlags_MenuBar != 0 {
		return w.DC.MenuBarOffset.y + w.CalcFontSize() + guiContext.Style.FramePadding.y*2.0
	}
	return 0
}

func (w *ImGuiWindow) MenuBarRect() ImRect {
	var y1 = w.Pos.y + w.TitleBarHeight()
	return ImRect{ImVec2{w.Pos.x, y1}, ImVec2{w.Pos.x + w.SizeFull.x, y1 + w.MenuBarHeight()}}
}

var IM_COL32_DISABLE = IM_COL32(0, 0, 0, 1) // Special sentinel code which cannot be used as a regular color.

const IMGUI_TABLE_MAX_COLUMNS = 64               // sizeof(ImU64) * 8. This is solely because we frequently encode columns set in a ImU64.
const IMGUI_TABLE_MAX_DRAW_CHANNELS = (4 + 64*2) // See TableSetupDrawChannels()

// ImGuiTableColumnIdx Our current column maximum is 64 but we may raise that in the future.
type ImGuiTableColumnIdx = ImS8
type ImGuiTableDrawChannelIdx = ImU8

// ImGuiTableColumn [Internal] sizeof() ~ 104
// We use the terminology "Enabled" to refer to a column that is not Hidden by user/api.
// We use the terminology "Clipped" to refer to a column that is out of sight because of scrolling/clipping.
// This is in contrast with some user-facing api such as IsItemVisible() / IsRectVisible() which use "Visible" to mean "not clipped".
type ImGuiTableColumn struct {
	Flags                    ImGuiTableColumnFlags // Flags after some patching (not directly same as provided by user). See ImGuiTableColumnFlags_
	WidthGiven               float                 // Final/actual width visible == (MaxX - MinX), locked in TableUpdateLayout(). May be > WidthRequest to honor minimum width, may be < WidthRequest to honor shrinking columns down in tight space.
	MinX                     float                 // Absolute positions
	MaxX                     float
	WidthRequest             float   // Master width absolute value when !(Flags & _WidthStretch). When Stretch this is derived every frame from StretchWeight in TableUpdateLayout()
	WidthAuto                float   // Automatic width
	StretchWeight            float   // Master width weight when (Flags & _WidthStretch). Often around ~1.0f initially.
	InitStretchWeightOrWidth float   // Value passed to TableSetupColumn(). For Width it is a content width (_without padding_).
	ClipRect                 ImRect  // Clipping rectangle for the column
	UserID                   ImGuiID // Optional, value passed to TableSetupColumn()
	WorkMinX                 float   // Contents region min ~(MinX + CellPaddingX + CellSpacingX1) == cursor start position when entering column
	WorkMaxX                 float   // Contents region max ~(MaxX - CellPaddingX - CellSpacingX2)
	ItemWidth                float   // Current item width for the column, preserved across rows
	ContentMaxXFrozen        float   // Contents maximum position for frozen rows (apart from headers), from which we can infer content width. TableHeader() automatically softclip itself + report ideal desired size, to avoid creating extraneous draw calls
	ContentMaxXUnfrozen      float
	ContentMaxXHeadersUsed   float // Contents maximum position for headers rows (regardless of freezing). TableHeader() automatically softclip itself + report ideal desired size, to avoid creating extraneous draw calls
	ContentMaxXHeadersIdeal  float
	NameOffset               ImS16                    // Offset into parent ColumnsNames[]
	DisplayOrder             ImGuiTableColumnIdx      // Index within Table's IndexToDisplayOrder[] (column may be reordered by users)
	IndexWithinEnabledSet    ImGuiTableColumnIdx      // Index within enabled/visible set (<= IndexToDisplayOrder)
	PrevEnabledColumn        ImGuiTableColumnIdx      // Index of prev enabled/visible column within Columns[], -1 if first enabled/visible column
	NextEnabledColumn        ImGuiTableColumnIdx      // Index of next enabled/visible column within Columns[], -1 if last enabled/visible column
	SortOrder                ImGuiTableColumnIdx      // Index of this column within sort specs, -1 if not sorting on this column, 0 for single-sort, may be >0 on multi-sort
	DrawChannelCurrent       ImGuiTableDrawChannelIdx // Index within DrawSplitter.Channels[]
	DrawChannelFrozen        ImGuiTableDrawChannelIdx // Draw channels for frozen rows (often headers)
	DrawChannelUnfrozen      ImGuiTableDrawChannelIdx // Draw channels for unfrozen rows
	IsEnabled                bool                     // IsUserEnabled && (Flags & ImGuiTableColumnFlags_Disabled) == 0
	IsUserEnabled            bool                     // Is the column not marked Hidden by the user? (unrelated to being off view, e.guiContext. clipped by scrolling).
	IsUserEnabledNextFrame   bool
	IsVisibleX               bool // Is actually in view (e.guiContext. overlapping the host window clipping rectangle, not scrolled).
	IsVisibleY               bool
	IsRequestOutput          bool // Return value for TableSetColumnIndex() / TableNextColumn(): whether we request user to output contents or not.
	IsSkipItems              bool // Do we want item submissions to this column to be completely ignored (no layout will happen).
	IsPreserveWidthAuto      bool
	NavLayerCurrent          ImS8               // ImGuiNavLayer in 1 byte
	AutoFitQueue             ImU8               // Queue of 8 values for the next 8 frames to request auto-fit
	CannotSkipItemsQueue     ImU8               // Queue of 8 values for the next 8 frames to disable Clipped/SkipItem
	SortDirection            ImGuiSortDirection //2 //:                                           // ImGuiSortDirection_Ascending or ImGuiSortDirection_Descending
	SortDirectionsAvailCount ImU8               //2 //:                                           // Number of available sort directions (0 to 3)
	SortDirectionsAvailMask  ImU8               //4 //:                                           // Mask of available sort directions (1-bit each)
	SortDirectionsAvailList  ImU8               // Ordered of available sort directions (2-bits each)
}

func NewImGuiTableColumn() ImGuiTableColumn {
	return ImGuiTableColumn{
		StretchWeight:         -1,
		WidthRequest:          -1,
		NameOffset:            -1,
		DisplayOrder:          -1,
		IndexWithinEnabledSet: -1,
		PrevEnabledColumn:     -1,
		NextEnabledColumn:     -1,
		SortOrder:             -1,
		DrawChannelCurrent:    (ImU8)(255),
		DrawChannelFrozen:     (ImU8)(255),
		DrawChannelUnfrozen:   (ImU8)(255),
	}
}

// ImGuiTableCellData Transient cell data stored per row.
// sizeof() ~ 6
type ImGuiTableCellData struct {
	BgColor ImU32               // Actual color
	Column  ImGuiTableColumnIdx // Column number
}

// ImGuiTable FIXME-TABLE: more transient data could be stored in a per-stacked table structure: DrawSplitter, SortSpecs, incoming RowData
type ImGuiTable struct {
	ID                         ImGuiID
	Flags                      ImGuiTableFlags
	RawData                    any                   // Single allocation to hold Columns[], DisplayOrderToIndex[] and RowCellData[]
	TempData                   *ImGuiTableTempData   // Transient data while table is active. Point within guiContext.CurrentTableStack[]
	Columns                    []ImGuiTableColumn    // ImGuiTableColumn Point within RawData[]
	DisplayOrderToIndex        []ImGuiTableColumnIdx // ImGuiTableColumnIdx Point within RawData[]. Store display order of columns (when not reordered, the values are 0...Count-1)
	RowCellData                []ImGuiTableCellData  // ImGuiTableCellData Point within RawData[]. Store cells background requests for current row.
	EnabledMaskByDisplayOrder  ImU64                 // Column DisplayOrder -> IsEnabled map
	EnabledMaskByIndex         ImU64                 // Column Index -> IsEnabled map (== not hidden by user/api) in a format adequate for iterating column without touching cold data
	VisibleMaskByIndex         ImU64                 // Column Index -> IsVisibleX|IsVisibleY map (== not hidden by user/api && not hidden by scrolling/cliprect)
	RequestOutputMaskByIndex   ImU64                 // Column Index -> IsVisible || AutoFit (== expect user to submit items)
	SettingsLoadedFlags        ImGuiTableFlags       // Which data were loaded from the .ini file (e.guiContext. when order is not altered we won't save order)
	SettingsOffset             int                   // Offset in guiContext.SettingsTables
	LastFrameActive            int
	ColumnsCount               int // Number of columns declared in BeginTable()
	CurrentRow                 int
	CurrentColumn              int
	InstanceCurrent            ImS16 // Count of BeginTable() calls with same ID in the same frame (generally 0). This is a little bit similar to BeginCount for a window, but multiple table with same ID look are multiple tables, they are just synched.
	InstanceInteracted         ImS16 // Mark which instance (generally 0) of the same ID is being interacted with
	RowPosY1                   float
	RowPosY2                   float
	RowMinHeight               float // Height submitted to TableNextRow()
	RowTextBaseline            float
	RowIndentOffsetX           float
	RowFlags                   ImGuiTableRowFlags // Current row flags, see ImGuiTableRowFlags_
	LastRowFlags               ImGuiTableRowFlags
	RowBgColorCounter          int      // Counter for alternating background colors (can be fast-forwarded by e.guiContext clipper), not same as CurrentRow because header rows typically don't increase this.
	RowBgColor                 [2]ImU32 // Background color override for current row.
	BorderColorStrong          ImU32
	BorderColorLight           ImU32
	BorderX1                   float
	BorderX2                   float
	HostIndentX                float
	MinColumnWidth             float
	OuterPaddingX              float
	CellPaddingX               float // Padding from each borders
	CellPaddingY               float
	CellSpacingX1              float // Spacing between non-bordered cells
	CellSpacingX2              float
	LastOuterHeight            float // Outer height from last frame
	LastFirstRowHeight         float // Height of first row from last frame
	InnerWidth                 float // User value passed to BeginTable(), see comments at the top of BeginTable() for details.
	ColumnsGivenWidth          float // Sum of current column width
	ColumnsAutoFitWidth        float // Sum of ideal column width in order nothing to be clipped, used for auto-fitting and content width submission in outer window
	ResizedColumnNextWidth     float
	ResizeLockMinContentsX2    float  // Lock minimum contents width while resizing down in order to not create feedback loops. But we allow growing the table.
	RefScale                   float  // Reference scale to be able to rescale columns on font/dpi changes.
	OuterRect                  ImRect // Note: for non-scrolling table, OuterRect.Max.y is often FLT_MAX until EndTable(), unless a height has been specified in BeginTable().
	InnerRect                  ImRect // InnerRect but without decoration. As with OuterRect, for non-scrolling tables, InnerRect.Max.y is
	WorkRect                   ImRect
	InnerClipRect              ImRect
	BgClipRect                 ImRect              // We use this to cpu-clip cell background color fill
	Bg0ClipRectForDrawCmd      ImRect              // Actual ImDrawCmd clip rect for BG0/1 channel. This tends to be == OuterWindow->ClipRect at BeginTable() because output in BG0/BG1 is cpu-clipped
	Bg2ClipRectForDrawCmd      ImRect              // Actual ImDrawCmd clip rect for BG2 channel. This tends to be a correct, tight-fit, because output to BG2 are done by widgets relying on regular ClipRect.
	HostClipRect               ImRect              // This is used to check if we can eventually merge our columns draw calls into the current draw call of the current window.
	HostBackupInnerClipRect    ImRect              // Backup of InnerWindow->ClipRect during PushTableBackground()/PopTableBackground()
	OuterWindow                *ImGuiWindow        // Parent window for the table
	InnerWindow                *ImGuiWindow        // Window holding the table data (== OuterWindow or a child window)
	ColumnsNames               []string            // Contiguous buffer holding columns names
	DrawSplitter               *ImDrawListSplitter // Shortcut to TempData->DrawSplitter while in table. Isolate draw commands per columns to avoid switching clip rect constantly
	SortSpecsSingle            ImGuiTableColumnSortSpecs
	SortSpecsMulti             []ImGuiTableColumnSortSpecs // FIXME-OPT: Using a small-vector pattern would be good.
	SortSpecs                  ImGuiTableSortSpecs         // Public facing sorts specs, this is what we return in TableGetSortSpecs()
	SortSpecsCount             ImGuiTableColumnIdx
	ColumnsEnabledCount        ImGuiTableColumnIdx      // Number of enabled columns (<= ColumnsCount)
	ColumnsEnabledFixedCount   ImGuiTableColumnIdx      // Number of enabled columns (<= ColumnsCount)
	DeclColumnsCount           ImGuiTableColumnIdx      // Count calls to TableSetupColumn()
	HoveredColumnBody          ImGuiTableColumnIdx      // Index of column whose visible region is being hovered. Important: == ColumnsCount when hovering empty region after the right-most column!
	HoveredColumnBorder        ImGuiTableColumnIdx      // Index of column whose right-border is being hovered (for resizing).
	AutoFitSingleColumn        ImGuiTableColumnIdx      // Index of single column requesting auto-fit.
	ResizedColumn              ImGuiTableColumnIdx      // Index of column being resized. Reset when InstanceCurrent==0.
	LastResizedColumn          ImGuiTableColumnIdx      // Index of column being resized from previous frame.
	HeldHeaderColumn           ImGuiTableColumnIdx      // Index of column header being held.
	ReorderColumn              ImGuiTableColumnIdx      // Index of column being reordered. (not cleared)
	ReorderColumnDir           ImGuiTableColumnIdx      // -1 or +1
	LeftMostEnabledColumn      ImGuiTableColumnIdx      // Index of left-most non-hidden column.
	RightMostEnabledColumn     ImGuiTableColumnIdx      // Index of right-most non-hidden column.
	LeftMostStretchedColumn    ImGuiTableColumnIdx      // Index of left-most stretched column.
	RightMostStretchedColumn   ImGuiTableColumnIdx      // Index of right-most stretched column.
	ContextPopupColumn         ImGuiTableColumnIdx      // Column right-clicked on, of -1 if opening context menu from a neutral/empty spot
	FreezeRowsRequest          ImGuiTableColumnIdx      // Requested frozen rows count
	FreezeRowsCount            ImGuiTableColumnIdx      // Actual frozen row count (== FreezeRowsRequest, or == 0 when no scrolling offset)
	FreezeColumnsRequest       ImGuiTableColumnIdx      // Requested frozen columns count
	FreezeColumnsCount         ImGuiTableColumnIdx      // Actual frozen columns count (== FreezeColumnsRequest, or == 0 when no scrolling offset)
	RowCellDataCurrent         ImGuiTableColumnIdx      // Index of current RowCellData[] entry in current row
	DummyDrawChannel           ImGuiTableDrawChannelIdx // Redirect non-visible columns here.
	Bg2DrawChannelCurrent      ImGuiTableDrawChannelIdx // For Selectable() and other widgets drawing across columns after the freezing line. Index within DrawSplitter.Channels[]
	Bg2DrawChannelUnfrozen     ImGuiTableDrawChannelIdx
	IsLayoutLocked             bool // Set by TableUpdateLayout() which is called when beginning the first row.
	IsInsideRow                bool // Set when inside TableBeginRow()/TableEndRow().
	IsInitializing             bool
	IsSortSpecsDirty           bool
	IsUsingHeaders             bool // Set when the first row had the ImGuiTableRowFlags_Headers flag.
	IsContextPopupOpen         bool // Set when default context menu is open (also see: ContextPopupColumn, InstanceInteracted).
	IsSettingsRequestLoad      bool
	IsSettingsDirty            bool // Set when table settings have changed and needs to be reported into ImGuiTableSetttings data.
	IsDefaultDisplayOrder      bool // Set when display order is unchanged from default (DisplayOrder contains 0...Count-1)
	IsResetAllRequest          bool
	IsResetDisplayOrderRequest bool
	IsUnfrozenRows             bool // Set when we got past the frozen row.
	IsDefaultSizingPolicy      bool // Set if user didn't explicitly set a sizing policy in BeginTable()
	MemoryCompacted            bool
	HostSkipItems              bool // Backup of InnerWindow->SkipItem at the end of BeginTable(), because we will overwrite InnerWindow->SkipItem on a per-column basis
}

func NewImGuiTable() ImGuiTable {
	return ImGuiTable{
		LastFrameActive: -1,
	}
}

// ImGuiTableTempData Transient data that are only needed between BeginTable() and EndTable(), those buffers are shared (1 per level of stacked table).
// - Accessing those requires chasing an extra pointer so for very frequently used data we leave them in the main table structure.
// - We also leave out of this structure data that tend to be particularly useful for debugging/metrics.
type ImGuiTableTempData struct {
	TableIndex                   int     // Index in guiContext.Tables.Buf[] pool
	LastTimeActive               float32 // Last timestamp this structure was used
	UserOuterSize                ImVec2  // outer_size.x passed to BeginTable()
	DrawSplitter                 ImDrawListSplitter
	HostBackupWorkRect           ImRect  // Backup of InnerWindow->WorkRect at the end of BeginTable()
	HostBackupParentWorkRect     ImRect  // Backup of InnerWindow->ParentWorkRect at the end of BeginTable()
	HostBackupPrevLineSize       ImVec2  // Backup of InnerWindow->DC.PrevLineSize at the end of BeginTable()
	HostBackupCurrLineSize       ImVec2  // Backup of InnerWindow->DC.CurrLineSize at the end of BeginTable()
	HostBackupCursorMaxPos       ImVec2  // Backup of InnerWindow->DC.CursorMaxPos at the end of BeginTable()
	HostBackupColumnsOffset      ImVec1  // Backup of OuterWindow->DC.ColumnsOffset at the end of BeginTable()
	HostBackupItemWidth          float32 // Backup of OuterWindow->DC.ItemWidth at the end of BeginTable()
	HostBackupItemWidthStackSize int     // Backup of OuterWindow->DC.ItemWidthStack.Size at the end of BeginTable()
}

func NewImGuiTableTempData() ImGuiTableTempData {
	return ImGuiTableTempData{
		LastTimeActive: -1.0,
	}
}

// ImGuiTableColumnSettings sizeof() ~ 12
type ImGuiTableColumnSettings struct {
	WidthOrWeight float
	UserID        ImGuiID
	Index         ImGuiTableColumnIdx
	DisplayOrder  ImGuiTableColumnIdx
	SortOrder     ImGuiTableColumnIdx
	SortDirection ImU8
	IsEnabled     ImU8 // "Visible" in ini file
	IsStretch     ImU8
}

func NewImGuiTableColumnSettings() ImGuiTableColumnSettings {
	return ImGuiTableColumnSettings{
		Index: -1,
	}
}

// ImGuiTableSettings This is designed to be stored in a single ImChunkStream (1 header followed by N ImGuiTableColumnSettings, etc.)
type ImGuiTableSettings struct {
	ID              ImGuiID         // Set to 0 to invalidate/delete the setting
	SaveFlags       ImGuiTableFlags // Indicate data we want to save using the Resizable/Reorderable/Sortable/Hideable flags (could be using its own flags..)
	RefScale        float           // Reference scale to be able to rescale columns on font/dpi changes.
	ColumnsCount    ImGuiTableColumnIdx
	ColumnsCountMax ImGuiTableColumnIdx // Maximum number of columns this settings instance can store, we can recycle a settings instance with lower number of columns but not higher
	WantApply       bool                // Set when loaded from .ini data (to enable merging/loading .ini data into an already running context)

	Columns []ImGuiTableColumnSettings
}

// ImFontBuilderIO This structure is likely to evolve as we add support for incremental atlas updates
type ImFontBuilderIO struct {
	FontBuilder_Build func(atlas *ImFontAtlas) bool
}
