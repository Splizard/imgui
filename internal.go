package imgui

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"unsafe"
)

const FLT_MAX = math.MaxFloat32
const INT_MAX = math.MaxInt32

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

// Current context pointer. Implicitly used by all Dear ImGui functions. Always assumed to be != nil.
// - ImGui::CreateContext() will automatically set this pointer if it is nil.
//   Change to a different context by calling ImGui::SetCurrentContext().
// - Important: Dear ImGui functions are not thread-safe because of this pointer.
//   If you want thread-safety to allow N threads to access N different contexts:
//   - Change this variable to use thread local storage so each thread can refer to a different context, in your imconfig.h:
//         struct ImGuiContext;
//         extern thread_local ImGuiContext* MyImGuiTLS;
//         #define GImGui MyImGuiTLS
//     And then define MyImGuiTLS in one of your cpp files. Note that thread_local is a C++11 keyword, earlier C++ uses compiler-specific keyword.
//   - Future development aims to make this context pointer explicit to all calls. Also read https://github.com/ocornut/imgui/issues/586
//   - If you need a finite number of contexts, you may compile and use multiple instances of the ImGui code from a different namespace.
// - DLL users: read comments above.
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

func ImHashData(ptr unsafe.Pointer, data_size uintptr, seed ImU32) ImGuiID {
	var crc ImU32 = ^seed
	var data *byte = (*byte)(ptr)
	var crc32_lut = &GCrc32LookupTable
	for i := uintptr(0); i < data_size; i++ {
		crc = (crc >> 8) ^ crc32_lut[(crc&0xFF)^uint(*data)]
		data = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(data)) + 1))
	}
	return ^crc
}

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
} // return out_// return output UTF-8 bytes count

func ImTextStrFromUtf8(out_buf []ImWchar, out_buf_size int, text string, in_remaining *string) int {
	var count int
	for i, char := range text {
		if count >= int(len(out_buf)) {
			*in_remaining = text[i:]
			return int(count)
		}
		out_buf[i] = char
		count++
	}
	return int(len(text))
}

// return number of UTF-8 code-points (NOT bytes count)
func ImTextCountCharsFromUtf8(in_text string) int {
	var count int
	for range in_text {
		count++
	}
	return count
}

func ImTextCountUtf8BytesFromChar(in_text, in_text_end string) int   { panic("not implemented") } // return number of bytes to express one char in UTF-8
func ImTextCountUtf8BytesFromStr(in_text, in_text_end []ImWchar) int { panic("not implemented") } // return number of bytes to express string in UTF-8

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

// Helper: Load file content into memory
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
		var a_mod int = n & 31
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

func (this *ImBitVector) Create(sz int) {
	*this = make([]ImU32, (uint(sz)+31)>>5)
}

func (this *ImBitVector) Clear() {
	*this = (*this)[:0]
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

func (this *ImSpan) Set(data interface{}) {
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
		var a float = ((float)(i) * 2 * IM_PI) / (float)(len(this.ArcFastVtx))
		this.ArcFastVtx[i] = ImVec2{ImCos(a), ImSin(a)}
	}
	this.ArcFastRadiusCutoff = IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(IM_DRAWLIST_ARCFAST_SAMPLE_MAX, this.CircleSegmentMaxError)
	return this
}

func (this *ImDrawListSharedData) SetCircleTessellationMaxError(max_error float) {
	if this.CircleSegmentMaxError == max_error {
		return
	}
	IM_ASSERT(max_error > 0.0)
	this.CircleSegmentMaxError = max_error
	for i := range this.CircleSegmentCounts {
		var radius float = (float)(i)
		if i > 0 {
			this.CircleSegmentCounts[i] = uint8(IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(radius, this.CircleSegmentMaxError))
		}
	}
	this.ArcFastRadiusCutoff = IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC_R(IM_DRAWLIST_ARCFAST_SAMPLE_MAX, this.CircleSegmentMaxError)
}

type ImDrawDataBuilder [2][]*ImDrawList

func (this *ImDrawDataBuilder) Clear() {
	for i := range this {
		this[i] = this[i][:0]
	}
}

func (this *ImDrawDataBuilder) GetDrawListCount() int {
	var count int
	for i := range this {
		count += int(len(this[i]))
	}
	return count
}

func (Layers *ImDrawDataBuilder) FlattenIntoSingleLayer() {
	var n = len(Layers[0])
	var size int = int(n)
	for i := 1; i < len(*Layers); i++ {
		size += int(len(Layers[i]))
	}

	//TODO/FIXME this could be wrong
	if size < int(len(Layers[0])) {
		Layers[0] = Layers[0][0:size]
	} else {
		Layers[0] = append(Layers[0], make([]*ImDrawList, 0, size-int(len(Layers[0])))...)
	}

	for layer_n := 1; layer_n < len(*Layers); layer_n++ {
		var layer *[]*ImDrawList = &Layers[layer_n]
		if len(*layer) == 0 {
			continue
		}
		copy(Layers[0][n:], *layer)
		n += len(*layer)
		*layer = (*layer)[:0]
	}
}

type ImGuiDataTypeTempStorage [8]byte

// Type information associated to one ImGuiDataType. Retrieve with DataTypeGetInfo().
type ImGuiDataTypeInfo struct {
	Size     size_t // Size in bytes
	Name     string // Short descriptive name for the type, for debugging
	PrintFmt string // Default printf format for the type
	ScanFmt  string // Default scanf format for the type
}

// Stacked color modifier, backup of modified data so we can restore it
type ImGuiColorMod struct {
	Col         ImGuiCol
	BackupValue ImVec4
}

type ImGuiStyleMod struct {
	VarIdx      ImGuiStyleVar
	BackupValue [2]int
}

func NewImGuiStyleModInt(idx ImGuiStyleVar, v int) ImGuiStyleMod {
	return ImGuiStyleMod{VarIdx: idx, BackupValue: [2]int{v, 0}}
}

func NewImGuiStyleModFloat(idx ImGuiStyleVar, v float32) ImGuiStyleMod {
	return ImGuiStyleMod{VarIdx: idx, BackupValue: [2]int{*(*int)(unsafe.Pointer(&v)), 0}}
}

func NewImGuiStyleModVec(idx ImGuiStyleVar, v ImVec2) ImGuiStyleMod {
	return ImGuiStyleMod{VarIdx: idx, BackupValue: [2]int{*(*int)(unsafe.Pointer(&v.x)), *(*int)(unsafe.Pointer(&v.y))}}
}

// Storage data for BeginComboPreview()/EndComboPreview()
type ImGuiComboPreviewData struct {
	PreviewRect                  ImRect
	BackupCursorPos              ImVec2
	BackupCursorMaxPos           ImVec2
	BackupCursorPosPrevLine      ImVec2
	BackupPrevLineTextBaseOffset float
	BackupLayout                 ImGuiLayoutType
}

// Stacked storage data for BeginGroup()/EndGroup()
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

func NewImGuiMenuColumns() ImGuiMenuColumns {
	return ImGuiMenuColumns{}
}

func (this *ImGuiMenuColumns) Update(spacing float, window_reappearing bool) {
	if window_reappearing {
		this.Spacing = ImU16(spacing)
		this.NextTotalWidth = this.TotalWidth
		this.OffsetLabel = 0
		this.OffsetShortcut = 0
		this.OffsetMark = 0
		this.OffsetIcon = 0
	}
}

func (this *ImGuiMenuColumns) DeclColumns(w_icon float, w_label float, w_shortcut float, w_mark float) float {
	this.Widths[0] = ImU16(w_icon)
	this.Widths[1] = ImU16(w_label)
	this.Widths[2] = ImU16(w_shortcut)
	this.Widths[3] = ImU16(w_mark)
	this.TotalWidth = ImU32(this.Widths[0]) + ImU32(this.Widths[1]) + ImU32(this.Widths[2]) + ImU32(this.Widths[3])
	return float(this.TotalWidth)
}

func (this *ImGuiMenuColumns) CalcNextTotalWidth(update_offsets bool) {
	this.NextTotalWidth = ImU32(this.Widths[0]) + ImU32(this.Widths[1]) + ImU32(this.Widths[2]) + ImU32(this.Widths[3])
	if update_offsets {
		this.OffsetLabel = this.Widths[0]
		this.OffsetShortcut = this.Widths[0] + this.Widths[1]
		this.OffsetMark = this.OffsetShortcut + this.Widths[2]
		this.OffsetIcon = this.OffsetMark + this.Widths[3]
	}
}

// Internal state of the currently focused/edited text input box
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
	UserCallbackData     interface{}
}

func (this *ImGuiInputTextState) ClearText() {
	this.CurLenW = 0
	this.CurLenA = 0
	this.TextW[0] = 0
	this.TextA[0] = 0
	this.CursorClamp()
}

func (this *ImGuiInputTextState) GetUndoAvailCount() int {
	return int(this.Stb.undostate.undo_point)
}

func (this *ImGuiInputTextState) GetRedoAvailCount() int {
	return int(STB_TEXTEDIT_UNDOSTATECOUNT - this.Stb.undostate.redo_point)
}

func (this *ImGuiInputTextState) OnKeyPressed(key int) {
	panic("not implemented")
}

func (this *ImGuiInputTextState) CursorAnimReset() {
	this.CursorAnim = -0.30
}

func (this *ImGuiInputTextState) CursorClamp() {
	this.Stb.cursor = ImMinInt(this.Stb.cursor, this.CurLenW)
	this.Stb.select_start = ImMinInt(this.Stb.select_start, this.CurLenW)
	this.Stb.select_end = ImMinInt(this.Stb.select_end, this.CurLenW)
}

func (this *ImGuiInputTextState) HasSelection() bool {
	return this.Stb.select_start != this.Stb.select_end
}

func (this *ImGuiInputTextState) ClearSelection() {
	this.Stb.select_start = this.Stb.cursor
	this.Stb.select_end = this.Stb.cursor
}

func (this *ImGuiInputTextState) GetCursorPos() int {
	return this.Stb.cursor
}

func (this *ImGuiInputTextState) GetSelectionStart() int {
	return this.Stb.select_start
}

func (this *ImGuiInputTextState) GetSelectionEnd() int {
	return this.Stb.select_end
}

func (this *ImGuiInputTextState) SelectAll() {
	this.Stb.select_start = 0
	this.Stb.cursor = this.CurLenW
	this.Stb.select_end = this.CurLenW
	this.Stb.has_preferred_x = 0
}

type ImGuiPopupData struct {
	PopupId        ImGuiID      // Set on OpenPopup()
	Window         *ImGuiWindow // Resolved on BeginPopup() - may stay unresolved if user never calls OpenPopup()
	SourceWindow   *ImGuiWindow // Set on OpenPopup() copy of NavWindow at the time of opening the popup
	OpenFrameCount int          // Set on OpenPopup()
	OpenParentId   ImGuiID      // Set on OpenPopup(), we need this to differentiate multiple menu sets from each others (e.g. inside menu bar vs loose menu items)
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
	SizeCallbackUserData interface{}
	BgAlphaVal           float
	MenuBarOffsetMinVal  ImVec2
}

func (this *ImGuiNextWindowData) ClearFlags() {
	this.Flags = ImGuiNextWindowDataFlags_None
}

type ImGuiNextItemData struct {
	Flags        ImGuiNextItemDataFlags
	Width        float     // Set by SetNextItemWidth()
	FocusScopeId ImGuiID   // Set by SetNextItemMultiSelectData() (!= 0 signify value has been set, so it's an alternate version of HasSelectionData, we don't use Flags for this because they are cleared too early. This is mostly used for debugging)
	OpenCond     ImGuiCond // Set by SetNextItemOpen()
	OpenVal      bool      // Set by SetNextItemOpen()
}

func (this *ImGuiNextItemData) ClearFlags() {
	this.Flags = ImGuiNextItemDataFlags_None
}

type ImGuiLastItemData struct {
	ID          ImGuiID
	InFlags     ImGuiItemFlags       // See ImGuiItemFlags_
	StatusFlags ImGuiItemStatusFlags // See ImGuiItemStatusFlags_
	Rect        ImRect               // Full rectangle
	NavRect     ImRect               // Navigation scoring rectangle (not displayed)
	DisplayRect ImRect               // Display rectangle (only if ImGuiItemStatusFlags_HasDisplayRect is set)
}

// Data saved for each window pushed into the stack
type ImGuiWindowStackData struct {
	Window                   *ImGuiWindow
	ParentLastItemDataBackup ImGuiLastItemData
}

type ImGuiShrinkWidthItem struct {
	Index int
	Width float
}

type ImGuiPtrOrIndex struct {
	Ptr   interface{} // Either field can be set, not both. e.g. Dock node tab bars are loose while BeginTabBar() ones are in a pool.
	Index int         // Usually index in a main pool.
}

func ImGuiPtr(ptr interface{}) ImGuiPtrOrIndex {
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

func (this *ImGuiNavItemData) Clear() {
	this.DistBox = FLT_MAX
	this.DistCenter = FLT_MAX
	this.DistAxial = FLT_MAX
}

// Storage data for a single column for legacy Columns() api
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

// ImGuiViewport Private/Internals fields (cardinal sin: we are using inheritance!)
// Every instance of ImGuiViewport is in fact a ImGuiViewportP.
type ImGuiViewportP = ImGuiViewport

func NewImGuiViewportP() ImGuiViewportP {
	return ImGuiViewportP{
		DrawListsLastFrame: [2]int{-1, -1},
	}
}

func (this *ImGuiViewportP) CalcWorkRectPos(off_min *ImVec2) ImVec2 {
	return ImVec2{this.Pos.x + off_min.x, this.Pos.y + off_min.y}
}

func (this *ImGuiViewportP) CalcWorkRectSize(off_min *ImVec2, off_max *ImVec2) ImVec2 {
	return ImVec2{ImMax(0.0, this.Size.x-off_min.x+off_max.x), ImMax(0.0, this.Size.y-off_min.y+off_max.y)}
}

func (this *ImGuiViewportP) UpdateWorkRect() {
	this.WorkPos = this.CalcWorkRectPos(&this.WorkOffsetMin)
	this.WorkSize = this.CalcWorkRectSize(&this.WorkOffsetMin, &this.WorkOffsetMax)
}

func (this *ImGuiViewportP) GetMainRect() ImRect {
	return ImRect{ImVec2{this.Pos.x, this.Pos.y}, ImVec2{this.Pos.x + this.Size.x, this.Pos.y + this.Size.y}}
}

func (this *ImGuiViewportP) GetWorkRect() ImRect {
	return ImRect{ImVec2{this.WorkPos.x, this.WorkPos.y}, ImVec2{this.WorkPos.x + this.WorkSize.x, this.WorkPos.y + this.WorkSize.y}}
}

func (this *ImGuiViewportP) GetBuildWorkRect() ImRect {
	var pos ImVec2 = this.CalcWorkRectPos(&this.BuildWorkOffsetMin)
	var size ImVec2 = this.CalcWorkRectSize(&this.BuildWorkOffsetMin, &this.BuildWorkOffsetMax)
	return ImRect{ImVec2{pos.x, pos.y}, ImVec2{pos.x + size.x, pos.y + size.y}}
}

// Windows data saved in imgui.ini file
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

func (this *ImGuiWindowSettings) GetName() string {
	return this.name
}

type ImGuiSettingsHandler struct {
	TypeName   string // Short description stored in .ini file. Disallowed characters: '[' ']'
	TypeHash   ImGuiID
	ClearAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                                 // Clear all settings data
	ReadInitFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                                 // Read: Called before reading (in registration order)
	ReadOpenFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, name string) interface{}        // Read: Called when entering into a new ini entry e.g. "[Window][Name]"
	ReadLineFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, entry interface{}, line string) // Read: Called for every line of text within an ini entry
	ApplyAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                                 // Read: Called after reading (in registration order)
	WriteAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, out_buf *ImGuiTextBuffer)       // Write: Output every entries into 'out_buf'
	UserData   interface{}
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

func (this *ImGuiStackSizes) SetToCurrentState() {
	var g = GImGui
	var window = g.CurrentWindow
	this.SizeOfIDStack = (short)(len(window.IDStack))
	this.SizeOfColorStack = (short)(len(g.ColorStack))
	this.SizeOfStyleVarStack = (short)(len(g.StyleVarStack))
	this.SizeOfFontStack = (short)(len(g.FontStack))
	this.SizeOfFocusScopeStack = (short)(len(g.FocusScopeStack))
	this.SizeOfGroupStack = (short)(len(g.GroupStack))
	this.SizeOfBeginPopupStack = (short)(len(g.BeginPopupStack))
}

func (this *ImGuiStackSizes) CompareWithCurrentState() {
	var g = GImGui
	var window = g.CurrentWindow

	// Window stacks
	// NOT checking: DC.ItemWidth, DC.TextWrapPos (per window) to allow user to conveniently push once and not pop (they are cleared on Begin)
	IM_ASSERT_USER_ERROR(this.SizeOfIDStack == short(len(window.IDStack)), "PushID/PopID or TreeNode/TreePop Mismatch!")

	// Global stacks
	// For color, style and font stacks there is an incentive to use Push/Begin/Pop/.../End patterns, so we relax our checks a little to allow them.
	IM_ASSERT_USER_ERROR(this.SizeOfGroupStack == short(len(g.GroupStack)), "BeginGroup/EndGroup Mismatch!")
	IM_ASSERT_USER_ERROR(this.SizeOfBeginPopupStack == short(len(g.BeginPopupStack)), "BeginPopup/EndPopup or BeginMenu/EndMenu Mismatch!")
	IM_ASSERT_USER_ERROR(this.SizeOfColorStack >= short(len(g.ColorStack)), "PushStyleColor/PopStyleColor Mismatch!")
	IM_ASSERT_USER_ERROR(this.SizeOfStyleVarStack >= short(len(g.StyleVarStack)), "PushStyleVar/PopStyleVar Mismatch!")
	IM_ASSERT_USER_ERROR(this.SizeOfFontStack >= short(len(g.FontStack)), "PushFont/PopFont Mismatch!")
	IM_ASSERT_USER_ERROR(this.SizeOfFocusScopeStack == short(len(g.FocusScopeStack)), "PushFocusScope/PopFocusScope Mismatch!")
}

type ImGuiContextHookCallback func(ctx *ImGuiContext, hook *ImGuiContextHook)

type ImGuiContextHook struct {
	HookId   ImGuiID // A unique ID assigned by AddContextHook()
	Type     ImGuiContextHookType
	Owner    ImGuiID
	Callback ImGuiContextHookCallback
	UserData interface{}
}

type ImGuiContext struct {
	Initialized                        bool
	IO                                 ImGuiIO
	Style                              ImGuiStyle
	Font                               *ImFont
	FontSize                           float
	FontBaseSize                       float
	DrawListSharedData                 ImDrawListSharedData
	Time                               double
	FrameCount                         int
	FrameCountEnded                    int
	FrameCountRendered                 int
	WithinFrameScope                   bool // Set by NewFrame(), cleared by EndFrame()
	WithinFrameScopeWithImplicitWindow bool // Set by NewFrame(), cleared by EndFrame() when the implicit debug window has been pushed
	WithinEndChild                     bool // Set within EndChild()
	GcCompactAll                       bool

	// Windows state
	Windows                        []*ImGuiWindow // Windows, sorted in display order, back to front
	WindowsFocusOrder              []*ImGuiWindow // Root windows, sorted in focus order, back to front.
	WindowsTempSortBuffer          []*ImGuiWindow // Temporary buffer used in EndFrame() to reorder windows so parents are kept before their child
	CurrentWindowStack             []ImGuiWindowStackData
	WindowsById                    ImGuiStorage // Map window's ImGuiID to *ImGuiWindow
	WindowsActiveCount             int          // Number of unique windows submitted by frame
	WindowsHoverPadding            ImVec2       // Padding around resizable windows for which hovering on counts as hovering the window == ImMax(style.TouchExtraPadding, WINDOWS_HOVER_PADDING)
	CurrentWindow                  *ImGuiWindow // Window being drawn into
	HoveredWindow                  *ImGuiWindow // Window the mouse is hovering. Will typically catch mouse inputs.
	HoveredWindowUnderMovingWindow *ImGuiWindow // Hovered window ignoring MovingWindow. Only set if MovingWindow is set.
	MovingWindow                   *ImGuiWindow // Track the window we clicked on (in order to preserve focus). The actual window that is moved is generally MovingWindow->RootWindow.
	WheelingWindow                 *ImGuiWindow // Track the window we started mouse-wheeling on. Until a timer elapse or mouse has moved, generally keep scrolling the same window even if during the course of scrolling the mouse ends up hovering a child window.
	WheelingWindowRefMousePos      ImVec2
	WheelingWindowTimer            float

	// Item/widgets state and tracking information
	HoveredId                                ImGuiID // Hovered widget, filled during the frame
	HoveredIdPreviousFrame                   ImGuiID
	HoveredIdAllowOverlap                    bool
	HoveredIdUsingMouseWheel                 bool // Hovered widget will use mouse wheel. Blocks scrolling the underlying window.
	HoveredIdPreviousFrameUsingMouseWheel    bool
	HoveredIdDisabled                        bool    // At least one widget passed the rect test, but has been discarded by disabled flag or popup inhibit. May be true even if HoveredId == 0.
	HoveredIdTimer                           float   // Measure contiguous hovering time
	HoveredIdNotActiveTimer                  float   // Measure contiguous hovering time where the item has not been active
	ActiveId                                 ImGuiID // Active widget
	ActiveIdIsAlive                          ImGuiID // Active widget has been seen this frame (we can't use a bool as the ActiveId may change within the frame)
	ActiveIdTimer                            float
	ActiveIdIsJustActivated                  bool // Set at the time of activation for one frame
	ActiveIdAllowOverlap                     bool // Active widget allows another widget to steal active id (generally for overlapping widgets, but not always)
	ActiveIdNoClearOnFocusLoss               bool // Disable losing active id if the active id window gets unfocused.
	ActiveIdHasBeenPressedBefore             bool // Track whether the active id led to a press (this is to allow changing between PressOnClick and PressOnRelease without pressing twice). Used by range_select branch.
	ActiveIdHasBeenEditedBefore              bool // Was the value associated to the widget Edited over the course of the Active state.
	ActiveIdHasBeenEditedThisFrame           bool
	ActiveIdUsingMouseWheel                  bool   // Active widget will want to read mouse wheel. Blocks scrolling the underlying window.
	ActiveIdUsingNavDirMask                  ImU32  // Active widget will want to read those nav move requests (e.g. can activate a button and move away from it)
	ActiveIdUsingNavInputMask                ImU32  // Active widget will want to read those nav inputs.
	ActiveIdUsingKeyInputMask                ImU64  // Active widget will want to read those key inputs. When we grow the ImGuiKey enum we'll need to either to order the enum to make useful keys come first, either redesign this into e.g. a small array.
	ActiveIdClickOffset                      ImVec2 // Clicked offset from upper-left corner, if applicable (currently only set by ButtonBehavior)
	ActiveIdWindow                           *ImGuiWindow
	ActiveIdSource                           ImGuiInputSource // Activating with mouse or nav (gamepad/keyboard)
	ActiveIdMouseButton                      ImGuiMouseButton
	ActiveIdPreviousFrame                    ImGuiID
	ActiveIdPreviousFrameIsAlive             bool
	ActiveIdPreviousFrameHasBeenEditedBefore bool
	ActiveIdPreviousFrameWindow              *ImGuiWindow
	LastActiveId                             ImGuiID // Store the last non-zero ActiveId, useful for animation.
	LastActiveIdTimer                        float   // Store the last non-zero ActiveId timer since the beginning of activation, useful for animation.

	// Next window/item data
	CurrentItemFlags ImGuiItemFlags      // == g.ItemFlagsStack.back()
	NextItemData     ImGuiNextItemData   // Storage for SetNextItem** functions
	LastItemData     ImGuiLastItemData   // Storage for last submitted item (setup by ItemAdd)
	NextWindowData   ImGuiNextWindowData // Storage for SetNextWindow** functions

	// Shared stacks
	ColorStack      []ImGuiColorMod  // Stack for PushStyleColor()/PopStyleColor() - inherited by Begin()
	StyleVarStack   []ImGuiStyleMod  // Stack for PushStyleVar()/PopStyleVar() - inherited by Begin()
	FontStack       []*ImFont        // Stack for PushFont()/PopFont() - inherited by Begin()
	FocusScopeStack []ImGuiID        // Stack for PushFocusScope()/PopFocusScope() - not inherited by Begin(), unless child window
	ItemFlagsStack  []ImGuiItemFlags // Stack for PushItemFlag()/PopItemFlag() - inherited by Begin()
	GroupStack      []ImGuiGroupData // Stack for BeginGroup()/EndGroup() - not inherited by Begin()
	OpenPopupStack  []ImGuiPopupData // Which popups are open (persistent)
	BeginPopupStack []ImGuiPopupData // Which level of BeginPopup() we are in (reset every frame)

	// Viewports
	Viewports []*ImGuiViewportP // Active viewports (Size==1 in 'master' branch). Each viewports hold their copy of ImDrawData.

	// Gamepad/keyboard Navigation
	NavWindow                  *ImGuiWindow // Focused window for navigation. Could be called 'FocusWindow'
	NavId                      ImGuiID      // Focused item for navigation
	NavFocusScopeId            ImGuiID      // Identify a selection scope (selection code often wants to "clear other items" when landing on an item of the selection set)
	NavActivateId              ImGuiID      // ~~ (g.ActiveId == 0) && IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0, also set when calling ActivateItem()
	NavActivateDownId          ImGuiID      // ~~ IsNavInputDown(ImGuiNavInput_Activate) ? NavId : 0
	NavActivatePressedId       ImGuiID      // ~~ IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0
	NavInputId                 ImGuiID      // ~~ IsNavInputPressed(ImGuiNavInput_Input) ? NavId : 0
	NavJustTabbedId            ImGuiID      // Just tabbed to this id.
	NavJustMovedToId           ImGuiID      // Just navigated to this id (result of a successfully MoveRequest).
	NavJustMovedToFocusScopeId ImGuiID      // Just navigated to this focus scope id (result of a successfully MoveRequest).
	NavJustMovedToKeyMods      ImGuiKeyModFlags
	NavNextActivateId          ImGuiID          // Set by ActivateItem(), queued until next frame.
	NavInputSource             ImGuiInputSource // Keyboard or Gamepad mode? THIS WILL ONLY BE None or NavGamepad or NavKeyboard.
	NavLayer                   ImGuiNavLayer    // Layer we are navigating on. For now the system is hard-coded for 0=main contents and 1=menu/title bar, may expose layers later.
	NavIdTabCounter            int              // == NavWindow->DC.FocusIdxTabCounter at time of NavId processing
	NavIdIsAlive               bool             // Nav widget has been seen this frame ~~ NavRectRel is valid
	NavMousePosDirty           bool             // When set we will update mouse position if (io.ConfigFlags & ImGuiConfigFlags_NavEnableSetMousePos) if set (NB: this not enabled by default)
	NavDisableHighlight        bool             // When user starts using mouse, we hide gamepad/keyboard highlight (NB: but they are still available, which is why NavDisableHighlight isn't always != NavDisableMouseHover)
	NavDisableMouseHover       bool             // When user starts using gamepad/keyboard, we hide mouse hovering highlight until mouse is touched again.

	// Navigation: Init & Move Requests
	NavAnyRequest             bool // ~~ NavMoveRequest || NavInitRequest this is to perform early out in ItemAdd()
	NavInitRequest            bool // Init request for appearing window to select first item
	NavInitRequestFromMove    bool
	NavInitResultId           ImGuiID // Init request result (first item of the window, or one for which SetItemDefaultFocus() was called)
	NavInitResultRectRel      ImRect  // Init request result rectangle (relative to parent window)
	NavMoveSubmitted          bool    // Move request submitted, will process result on next NewFrame()
	NavMoveScoringItems       bool    // Move request submitted, still scoring incoming items
	NavMoveForwardToNextFrame bool
	NavMoveFlags              ImGuiNavMoveFlags
	NavMoveKeyMods            ImGuiKeyModFlags
	NavMoveDir                ImGuiDir // Direction of the move request (left/right/up/down)
	NavMoveDirForDebug        ImGuiDir
	NavMoveClipDir            ImGuiDir         // FIXME-NAV: Describe the purpose of this better. Might want to rename?
	NavScoringRect            ImRect           // Rectangle used for scoring, in screen space. Based of window->NavRectRel[], modified for directional navigation scoring.
	NavScoringDebugCount      int              // Metrics for debugging
	NavMoveResultLocal        ImGuiNavItemData // Best move request candidate within NavWindow
	NavMoveResultLocalVisible ImGuiNavItemData // Best move request candidate within NavWindow that are mostly visible (when using ImGuiNavMoveFlags_AlsoScoreVisibleSet flag)
	NavMoveResultOther        ImGuiNavItemData // Best move request candidate within NavWindow's flattened hierarchy (when using ImGuiWindowFlags_NavFlattened flag)

	// Navigation: Windowing (CTRL+TAB for list, or Menu button + keys or directional pads to move/resize)
	NavWindowingTarget         *ImGuiWindow // Target window when doing CTRL+Tab (or Pad Menu + FocusPrev/Next), this window is temporarily displayed top-most!
	NavWindowingTargetAnim     *ImGuiWindow // Record of last valid NavWindowingTarget until DimBgRatio and NavWindowingHighlightAlpha becomes 0.0f, so the fade-out can stay on it.
	NavWindowingListWindow     *ImGuiWindow // Internal window actually listing the CTRL+Tab contents
	NavWindowingTimer          float
	NavWindowingHighlightAlpha float
	NavWindowingToggleLayer    bool

	// Legacy Focus/Tabbing system (older than Nav, active even if Nav is disabled, misnamed. FIXME-NAV: This needs a redesign!)
	TabFocusRequestCurrWindow         *ImGuiWindow //
	TabFocusRequestNextWindow         *ImGuiWindow //
	TabFocusRequestCurrCounterRegular int          // Any item being requested for focus, stored as an index (we on layout to be stable between the frame pressing TAB and the next frame, semi-ouch)
	TabFocusRequestCurrCounterTabStop int          // Tab item being requested for focus, stored as an index
	TabFocusRequestNextCounterRegular int          // Stored for next frame
	TabFocusRequestNextCounterTabStop int          // "
	TabFocusPressed                   bool         // Set in NewFrame() when user pressed Tab

	// Render
	DimBgRatio  float // 0.0..1.0 animation when fading in a dimming background (for modal window and CTRL+TAB list)
	MouseCursor ImGuiMouseCursor

	// Drag and Drop
	DragDropActive                  bool
	DragDropWithinSource            bool // Set when within a BeginDragDropXXX/EndDragDropXXX block for a drag source.
	DragDropWithinTarget            bool // Set when within a BeginDragDropXXX/EndDragDropXXX block for a drag target.
	DragDropSourceFlags             ImGuiDragDropFlags
	DragDropSourceFrameCount        int
	DragDropMouseButton             ImGuiMouseButton
	DragDropPayload                 ImGuiPayload
	DragDropTargetRect              ImRect // Store rectangle of current target candidate (we favor small targets when overlapping)
	DragDropTargetId                ImGuiID
	DragDropAcceptFlags             ImGuiDragDropFlags
	DragDropAcceptIdCurrRectSurface float    // Target item surface (we resolve overlapping targets by prioritizing the smaller surface)
	DragDropAcceptIdCurr            ImGuiID  // Target item id (set at the time of accepting the payload)
	DragDropAcceptIdPrev            ImGuiID  // Target item id from previous frame (we need to store this to allow for overlapping drag and drop targets)
	DragDropAcceptFrameCount        int      // Last time a target expressed a desire to accept the source
	DragDropHoldJustPressedId       ImGuiID  // Set when holding a payload just made ButtonBehavior() return a press.
	DragDropPayloadBufHeap          []byte   // We don't expose the ImVector<> directly, ImGuiPayload only holds pointer+size
	DragDropPayloadBufLocal         [16]byte // Local buffer for small payloads

	// Table
	CurrentTable                *ImGuiTable
	CurrentTableStackIdx        int
	Tables                      map[ImGuiID]*ImGuiTable
	TablesTempDataStack         []ImGuiTableTempData
	TablesLastTimeActive        []float // Last used timestamp of each tables (SOA, for efficient GC)
	DrawChannelsTempMergeBuffer []ImDrawChannel

	// Tab bars
	CurrentTabBar      *ImGuiTabBar
	TabBars            map[ImGuiID]ImGuiTabBar
	CurrentTabBarStack []ImGuiPtrOrIndex
	ShrinkWidthBuffer  []ImGuiShrinkWidthItem

	// Widget state
	MouseLastValidPos               ImVec2
	InputTextState                  ImGuiInputTextState
	InputTextPasswordFont           ImFont
	TempInputId                     ImGuiID             // Temporary text input when CTRL+clicking on a slider, etc.
	ColorEditOptions                ImGuiColorEditFlags // Store user options for color edit widgets
	ColorEditLastHue                float               // Backup of last Hue associated to LastColor[3], so we can restore Hue in lossy RGB<>HSV round trips
	ColorEditLastSat                float               // Backup of last Saturation associated to LastColor[3], so we can restore Saturation in lossy RGB<>HSV round trips
	ColorEditLastColor              [3]float
	ColorPickerRef                  ImVec4 // Initial/reference color at the time of opening the color picker.
	ComboPreviewData                ImGuiComboPreviewData
	SliderCurrentAccum              float // Accumulated slider delta when using navigation controls.
	SliderCurrentAccumDirty         bool  // Has the accumulated slider delta changed since last time we tried to apply it?
	DragCurrentAccumDirty           bool
	DragCurrentAccum                float // Accumulator for dragging modification. Always high-precision, not rounded by end-user precision settings
	DragSpeedDefaultRatio           float // If speed == 0.0f, uses (max-min) * DragSpeedDefaultRatio
	DisabledAlphaBackup             float // Backup for style.Alpha for BeginDisabled()
	ScrollbarClickDeltaToGrabCenter float // Distance between mouse and center of grab box, normalized in parent space. Use storage?
	TooltipOverrideCount            int
	TooltipSlowDelay                float     // Time before slow tooltips appears (FIXME: This is temporary until we merge in tooltip timer+priority work)
	ClipboardHandlerData            []char    // If no custom clipboard handler is defined
	MenusIdSubmittedThisFrame       []ImGuiID // A list of menu IDs that were rendered at least once

	// Platform support
	PlatformImePos             ImVec2 // Cursor position request & last passed to the OS Input Method Editor
	PlatformImeLastPos         ImVec2
	PlatformLocaleDecimalPoint char // '.' or *localeconv()->decimal_point

	// Settings
	SettingsLoaded     bool
	SettingsDirtyTimer float                  // Save .ini Settings to memory when time reaches zero
	SettingsIniData    ImGuiTextBuffer        // In memory .ini settings
	SettingsHandlers   []ImGuiSettingsHandler // List of .ini settings handlers
	SettingsWindows    []ImGuiWindowSettings  // ImGuiWindow .ini settings entries
	SettingsTables     []ImGuiTableSettings   // ImGuiTable .ini settings entries
	Hooks              []ImGuiContextHook     // Hooks for extensions (e.g. test engine)
	HookIdNext         ImGuiID                // Next available HookId

	// Capture/Logging
	LogEnabled              bool            // Currently capturing
	LogType                 ImGuiLogType    // Capture target
	LogFile                 ImFileHandle    // If != NULL log to stdout/ file
	LogBuffer               ImGuiTextBuffer // Accumulation buffer when log to clipboard. This is pointer so our GImGui static constructor doesn't call heap allocators.
	LogNextPrefix           string
	LogNextSuffix           string
	LogLinePosY             float
	LogLineFirstItem        bool
	LogDepthRef             int
	LogDepthToExpand        int
	LogDepthToExpandDefault int // Default/stored value for LogDepthMaxExpand if not specified in the LogXXX function call.

	// Debug Tools
	DebugItemPickerActive  bool    // Item picker is active (started with DebugStartItemPicker())
	DebugItemPickerBreakId ImGuiID // Will call IM_DEBUG_BREAK() when encountering this id
	DebugMetricsConfig     ImGuiMetricsConfig

	// Misc
	FramerateSecPerFrame         [120]float // Calculate estimate of framerate for user over the last 2 seconds.
	FramerateSecPerFrameIdx      int
	FramerateSecPerFrameCount    int
	FramerateSecPerFrameAccum    float
	WantCaptureMouseNextFrame    int // Explicit capture via CaptureKeyboardFromApp()/CaptureMouseFromApp() sets those flags
	WantCaptureKeyboardNextFrame int
	WantTextInputNextFrame       int
	TempBuffer                   [1024 * 31]byte // Temporary text buffer
}

func NewImGuiContext(atlas *ImFontAtlas) ImGuiContext {
	if atlas == nil {
		ptr := NewImFontAtlas()
		atlas = &ptr
	}
	var io = NewImGuiIO()
	io.Fonts = atlas
	return ImGuiContext{
		IO:                                io,
		DrawListSharedData:                NewImDrawListSharedData(),
		Style:                             NewImGuiStyle(),
		FrameCountEnded:                   -1,
		FrameCountRendered:                -1,
		ActiveIdClickOffset:               ImVec2{-1, -1},
		ActiveIdMouseButton:               -1,
		NavIdTabCounter:                   INT_MAX,
		TabFocusRequestCurrCounterRegular: INT_MAX,
		TabFocusRequestNextCounterRegular: INT_MAX,
		TabFocusRequestCurrCounterTabStop: INT_MAX,
		TabFocusRequestNextCounterTabStop: INT_MAX,
		MouseCursor:                       ImGuiMouseCursor_Arrow,
		DragDropSourceFrameCount:          -1,
		DragDropMouseButton:               -1,
		DragDropAcceptFrameCount:          -1,
		CurrentTableStackIdx:              -1,
		ColorEditLastColor:                [3]float{FLT_MAX, FLT_MAX, FLT_MAX},
		DragSpeedDefaultRatio:             1 / 100.0,
		TooltipSlowDelay:                  0.5,
		PlatformImePos:                    ImVec2{FLT_MAX, FLT_MAX},
		PlatformImeLastPos:                ImVec2{FLT_MAX, FLT_MAX},
		PlatformLocaleDecimalPoint:        '.',
		LogLinePosY:                       FLT_MAX,
		LogDepthToExpand:                  2,
		LogDepthToExpandDefault:           2,
		WantCaptureMouseNextFrame:         -1,
		WantCaptureKeyboardNextFrame:      -1,
		WantTextInputNextFrame:            -1,
	}
}

// Transient per-window data, reset at the beginning of the frame. This used to be called ImGuiDrawContext, hence the DC variable name in ImGuiWindow.
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
	TreeJumpToParentOnPopMask ImU32            // Store a copy of !g.NavIdIsAlive for TreeDepth 0..31.. Could be turned into a ImU64 if necessary.
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
	SkipItems                      bool    // Set when items can safely be all clipped (e.g. window not visible or collapsed)
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
	WorkRect          ImRect   // Initially covers the whole scrolling region. Reduced by containers e.g columns/tables when active. Shrunk by WindowPadding*1.0f on each side. This is meant to replace ContentRegionRect over time (from 1.71+ onward).
	ParentWorkRect    ImRect   // Backup of WorkRect before entering a container such as columns/tables. Used by e.g. SpanAllColumns functions to easily access. Stacked containers are responsible for maintaining this. // FIXME-WORKRECT: Could be a stack?
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

	NavLastChildNavWindow *ImGuiWindow                 // When going to the menu bar, we remember the child window we came from. (This could probably be made implicit if we kept g.Windows sorted by last focused including child window.)
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

func (this *ImGuiWindow) GetIDs(str string) ImGuiID {
	var seed ImGuiID = this.IDStack[len(this.IDStack)-1]
	var id ImGuiID = ImHashStr(str, 0, seed)
	KeepAliveID(id)
	return id
}
func (this *ImGuiWindow) GetIDInterface(ptr interface{}) ImGuiID { panic("not implemented") }

func (this *ImGuiWindow) GetIDInt(n int) ImGuiID {
	var seed ImGuiID = this.IDStack[len(this.IDStack)-1]
	var id ImGuiID = ImHashData(unsafe.Pointer(&n), unsafe.Sizeof(n), seed)
	KeepAliveID(id)
	return id
}

func (this *ImGuiWindow) GetIDNoKeepAlive(str string) ImGuiID {
	var seed ImGuiID = this.IDStack[len(this.IDStack)-1]
	var id ImGuiID = ImHashStr(str, 0, seed)
	return id
}

func (this *ImGuiWindow) GetIDNoKeepAliveInterface(ptr interface{}) ImGuiID { panic("not implemented") }
func (this *ImGuiWindow) GetIDNoKeepAliveInt(n int) ImGuiID                 { panic("not implemented") }
func (this *ImGuiWindow) GetIDFromRectangle(r_abs ImRect) ImGuiID           { panic("not implemented") }

func (this *ImGuiWindow) Rect() ImRect {
	return ImRect{ImVec2{this.Pos.x, this.Pos.y}, ImVec2{this.Pos.x + this.Size.x, this.Pos.y + this.Size.y}}
}

func (this *ImGuiWindow) CalcFontSize() float {
	var g *ImGuiContext = GImGui
	var scale float = g.FontBaseSize * this.FontWindowScale
	if this.ParentWindow != nil {
		scale *= this.ParentWindow.FontWindowScale
	}
	//return 20 //TODO/FIXME
	return scale
}

func (this *ImGuiWindow) TitleBarHeight() float {
	var g *ImGuiContext = GImGui
	if this.Flags&ImGuiWindowFlags_NoTitleBar != 0 {
		return 0.0
	}
	return this.CalcFontSize() + g.Style.FramePadding.y*2.0
}

func (this *ImGuiWindow) TitleBarRect() ImRect {
	return ImRect{ImVec2{this.Pos.x, this.Pos.y}, ImVec2{this.Pos.x + this.SizeFull.x, this.Pos.y + this.TitleBarHeight()}}
}

func (this *ImGuiWindow) MenuBarHeight() float {
	var g *ImGuiContext = GImGui
	if this.Flags&ImGuiWindowFlags_MenuBar != 0 {
		return this.DC.MenuBarOffset.y + this.CalcFontSize() + g.Style.FramePadding.y*2.0
	}
	return 0
}

func (this *ImGuiWindow) MenuBarRect() ImRect {
	var y1 float = this.Pos.y + this.TitleBarHeight()
	return ImRect{ImVec2{this.Pos.x, y1}, ImVec2{this.Pos.x + this.SizeFull.x, y1 + this.MenuBarHeight()}}
}

type ImGuiTabItem struct {
	ID                ImGuiID
	Flags             ImGuiTabItemFlags
	LastFrameVisible  int
	LastFrameSelected int   // This allows us to infer an ordered list of the last activated tabs with little maintenance
	Offset            float // Position relative to beginning of tab
	Width             float // Width currently displayed
	ContentWidth      float // Width of label, stored during BeginTabItem() call
	NameOffset        int   // When Window==NULL, offset to name within parent ImGuiTabBar::TabsNames
	BeginOrder        int   // BeginTabItem() order, used to re-order tabs after toggling ImGuiTabBarFlags_Reorderable
	IndexDuringLayout int   // Index only used during TabBarLayout()
	WantClose         bool  // Marked as closed by SetTabItemClosed()
}

func NewImGuiTabItem() ImGuiTabItem {
	return ImGuiTabItem{
		LastFrameVisible:  -1,
		LastFrameSelected: -1,
		NameOffset:        -1,
		BeginOrder:        -1,
		IndexDuringLayout: -1,
	}
}

type ImGuiTabBar struct {
	Tabs                            []ImGuiTabItem
	Flags                           ImGuiTabBarFlags
	ID                              ImGuiID         // Zero for tab-bars used by docking
	SelectedTabId                   ImGuiID         // Selected tab/window
	NextSelectedTabId               ImGuiID         // Next selected tab/window. Will also trigger a scrolling animation
	VisibleTabId                    ImGuiID         // Can occasionally be != SelectedTabId (e.g. when previewing contents for CTRL+TAB preview)
	CurrFrameVisible                int             //
	PrevFrameVisible                int             //
	BarRect                         ImRect          //
	CurrTabsContentsHeight          float           //
	PrevTabsContentsHeight          float           // Record the height of contents submitted below the tab bar
	WidthAllTabs                    float           // Actual width of all tabs (locked during layout)
	WidthAllTabsIdeal               float           // Ideal width if all tabs were visible and not clipped
	ScrollingAnim                   float           //
	ScrollingTarget                 float           //
	ScrollingTargetDistToVisibility float           //
	ScrollingSpeed                  float           //
	ScrollingRectMinX               float           //
	ScrollingRectMaxX               float           //
	ReorderRequestTabId             ImGuiID         //
	ReorderRequestOffset            int             //
	BeginCount                      int             //
	WantLayout                      bool            //
	VisibleTabWasSubmitted          bool            // Set to true when a new tab item or button has been added to the tab bar during last frame
	TabsAddedNew                    bool            // Set to true when a new tab item or button has been added to the tab bar during last frame
	TabsActiveCount                 int             // Number of tabs submitted this frame.
	LastTabItemIdx                  int             // Index of last BeginTabItem() tab for use by EndTabItem()
	ItemSpacingY                    float           //
	FramePadding                    ImVec2          // style.FramePadding locked at the time of BeginTabBar()
	BackupCursorPos                 ImVec2          //
	TabsNames                       ImGuiTextBuffer // For non-docking tab bar we re-append names in a contiguous buffer.
}

func NewImGuiTabBar() ImGuiTabBar {
	panic("not implemented")
	return ImGuiTabBar{}
}

func (this ImGuiTabBar) GetTabOrder(tab *ImGuiTabItem) int {
	for i := range this.Tabs {
		if tab == &this.Tabs[i] {
			return int(i)
		}
	}
	return -1
}

func (this ImGuiTabBar) GetTabName(tab *ImGuiTabItem) string {
	IM_ASSERT(tab.NameOffset != -1 && tab.NameOffset < int(len(this.TabsNames)))
	return string(this.TabsNames[tab.NameOffset:]) //TODO/FIXME zero termination
}

var IM_COL32_DISABLE = IM_COL32(0, 0, 0, 1) // Special sentinel code which cannot be used as a regular color.

const IMGUI_TABLE_MAX_COLUMNS = 64               // sizeof(ImU64) * 8. This is solely because we frequently encode columns set in a ImU64.
const IMGUI_TABLE_MAX_DRAW_CHANNELS = (4 + 64*2) // See TableSetupDrawChannels()

// Our current column maximum is 64 but we may raise that in the future.
type ImGuiTableColumnIdx = ImS8
type ImGuiTableDrawChannelIdx = ImU8

// [Internal] sizeof() ~ 104
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
	IsUserEnabled            bool                     // Is the column not marked Hidden by the user? (unrelated to being off view, e.g. clipped by scrolling).
	IsUserEnabledNextFrame   bool
	IsVisibleX               bool // Is actually in view (e.g. overlapping the host window clipping rectangle, not scrolled).
	IsVisibleY               bool
	IsRequestOutput          bool // Return value for TableSetColumnIndex() / TableNextColumn(): whether we request user to output contents or not.
	IsSkipItems              bool // Do we want item submissions to this column to be completely ignored (no layout will happen).
	IsPreserveWidthAuto      bool
	NavLayerCurrent          ImS8 // ImGuiNavLayer in 1 byte
	AutoFitQueue             ImU8 // Queue of 8 values for the next 8 frames to request auto-fit
	CannotSkipItemsQueue     ImU8 // Queue of 8 values for the next 8 frames to disable Clipped/SkipItem
	SortDirection            ImU8 //2 //:                                           // ImGuiSortDirection_Ascending or ImGuiSortDirection_Descending
	SortDirectionsAvailCount ImU8 //2 //:                                           // Number of available sort directions (0 to 3)
	SortDirectionsAvailMask  ImU8 //4 //:                                           // Mask of available sort directions (1-bit each)
	SortDirectionsAvailList  ImU8 // Ordered of available sort directions (2-bits each)
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

// Transient cell data stored per row.
// sizeof() ~ 6
type ImGuiTableCellData struct {
	BgColor ImU32               // Actual color
	Column  ImGuiTableColumnIdx // Column number
}

// FIXME-TABLE: more transient data could be stored in a per-stacked table structure: DrawSplitter, SortSpecs, incoming RowData
type ImGuiTable struct {
	ID                         ImGuiID
	Flags                      ImGuiTableFlags
	RawData                    interface{}         // Single allocation to hold Columns[], DisplayOrderToIndex[] and RowCellData[]
	TempData                   *ImGuiTableTempData // Transient data while table is active. Point within g.CurrentTableStack[]
	Columns                    ImSpan              // ImGuiTableColumn Point within RawData[]
	DisplayOrderToIndex        ImSpan              // ImGuiTableColumnIdx Point within RawData[]. Store display order of columns (when not reordered, the values are 0...Count-1)
	RowCellData                ImSpan              // ImGuiTableCellData Point within RawData[]. Store cells background requests for current row.
	EnabledMaskByDisplayOrder  ImU64               // Column DisplayOrder -> IsEnabled map
	EnabledMaskByIndex         ImU64               // Column Index -> IsEnabled map (== not hidden by user/api) in a format adequate for iterating column without touching cold data
	VisibleMaskByIndex         ImU64               // Column Index -> IsVisibleX|IsVisibleY map (== not hidden by user/api && not hidden by scrolling/cliprect)
	RequestOutputMaskByIndex   ImU64               // Column Index -> IsVisible || AutoFit (== expect user to submit items)
	SettingsLoadedFlags        ImGuiTableFlags     // Which data were loaded from the .ini file (e.g. when order is not altered we won't save order)
	SettingsOffset             int                 // Offset in g.SettingsTables
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
	RowBgColorCounter          int      // Counter for alternating background colors (can be fast-forwarded by e.g clipper), not same as CurrentRow because header rows typically don't increase this.
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
	ColumnsNames               ImGuiTextBuffer     // Contiguous buffer holding columns names
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

// Transient data that are only needed between BeginTable() and EndTable(), those buffers are shared (1 per level of stacked table).
// - Accessing those requires chasing an extra pointer so for very frequently used data we leave them in the main table structure.
// - We also leave out of this structure data that tend to be particularly useful for debugging/metrics.
type ImGuiTableTempData struct {
	TableIndex                   int     // Index in g.Tables.Buf[] pool
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

// sizeof() ~ 12
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

// This is designed to be stored in a single ImChunkStream (1 header followed by N ImGuiTableColumnSettings, etc.)
type ImGuiTableSettings struct {
	ID              ImGuiID         // Set to 0 to invalidate/delete the setting
	SaveFlags       ImGuiTableFlags // Indicate data we want to save using the Resizable/Reorderable/Sortable/Hideable flags (could be using its own flags..)
	RefScale        float           // Reference scale to be able to rescale columns on font/dpi changes.
	ColumnsCount    ImGuiTableColumnIdx
	ColumnsCountMax ImGuiTableColumnIdx // Maximum number of columns this settings instance can store, we can recycle a settings instance with lower number of columns but not higher
	WantApply       bool                // Set when loaded from .ini data (to enable merging/loading .ini data into an already running context)
}

// This structure is likely to evolve as we add support for incremental atlas updates
type ImFontBuilderIO struct {
	FontBuilder_Build func(atlas *ImFontAtlas) bool
}
