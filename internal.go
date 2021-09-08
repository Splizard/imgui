package imgui

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"unsafe"
)

const FLT_MAX = math.MaxFloat32

/*func NewImGuiContext(atlas *ImFontAtlas) *ImGuiContext {
	if atlas == nil {
		atlas = NewFontAtlas()
	}
	var ctx = new(ImGuiContext)
	ctx.IO.Fonts = atlas

	return &ctx
}*/

type ImGuiContextHook struct{}          // Hook for extensions like ImGuiTestEngine
type ImGuiTabBar struct{}               // Storage for a tab bar
type ImGuiTabItem struct{}              // Storage for a tab item (within a tab bar)
type ImGuiTable struct{}                // Storage for a table
type ImGuiTableColumn struct{}          // Storage for one column of a table
type ImGuiTableTempData struct{}        // Temporary storage for one table (one per table in the stack), shared between tables.
type ImGuiTableSettings struct{}        // Storage for a table .ini settings
type ImGuiTableColumnsSettings struct{} // Storage for a column .ini settings
type ImGuiWindow struct{}               // Storage for one window
type ImGuiWindowTempData struct{}       // Temporary storage for one window (that's the data which in theory we could ditch at the end of the frame, in practice we currently keep it for each window)

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
	TexUvLines            *ImVec4
}

func NewImDrawListSharedData() ImDrawListSharedData { panic("not implemented") }

func (this *ImDrawListSharedData) SetCircleTessellationMaxError(max_error float) {
	this.CircleSegmentMaxError = max_error
}

type ImDrawDataBuilder [2][]*ImDrawList

func (this *ImDrawDataBuilder) Clear() {
	for i := range this {
		this[i] = this[i][0:]
	}
}

func (this *ImDrawDataBuilder) GetDrawListCount() int {
	var count int
	for i := range this {
		count += int(len(this[i]))
	}
	return count
}

func (this *ImDrawDataBuilder) FlattenIntoSingleLayer() { panic("not implemented") }

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

func (this ImGuiNextItemData) ClearFlags() {
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
type ImGuiViewportP struct {
	ImGuiViewport

	DrawListsLastFrame [2]int         // Last frame number the background (0) and foreground (1) draw lists were used
	DrawLists          [2]*ImDrawList // Convenience background (0) and foreground (1) draw lists. We use them to draw software mouser cursor when io.MouseDrawCursor is set and to draw most debug overlays.
	DrawDataBuilder    ImDrawDataBuilder

	WorkOffsetMin      ImVec2 // Work Area: Offset from Pos to top-left corner of Work Area. Generally (0,0) or (0,+main_menu_bar_height). Work Area is Full Area but without menu-bars/status-bars (so WorkArea always fit inside Pos/Size!)
	WorkOffsetMax      ImVec2 // Work Area: Offset from Pos+Size to bottom-right corner of Work Area. Generally (0,0) or (0,-status_bar_height).
	BuildWorkOffsetMin ImVec2 // Work Area: Offset being built during current frame. Generally >= 0.0f.
	BuildWorkOffsetMax ImVec2 // Work Area: Offset being built during current frame. Generally <= 0.0f.
}

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
	ClearAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                                           // Clear all settings data
	ReadInitFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                                           // Read: Called before reading (in registration order)
	ReadOpenFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, name string)                              // Read: Called when entering into a new ini entry e.g. "[Window][Name]"
	ReadLineFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, entry *ImGuiSettingsHandler, line string) // Read: Called for every line of text within an ini entry
	ApplyAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler)                                           // Read: Called after reading (in registration order)
	WriteAllFn func(ctx *ImGuiContext, handler *ImGuiSettingsHandler, out_buf *ImGuiTextBuffer)                 // Write: Output every entries into 'out_buf'
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

func (this *ImGuiStackSizes) SetToCurrentState()       { panic("not implemented") }
func (this *ImGuiStackSizes) CompareWithCurrentState() { panic("not implemented") }
