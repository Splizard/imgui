package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"unicode"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
	B "github.com/snhmibby/filebuf"
)

const hexwidgetEditPopup = "hexwidget###editPopup"

func numHexDigits(addr int64) int {
	hexdigits := 1
	for addr > 0 {
		hexdigits++
		addr /= 16
	}
	return hexdigits
}

func addrLabel(addr int64, nDigits int) string {
	return fmt.Sprintf("%0*X:", nDigits, addr)
}

type HexViewWidget struct {
	state *ViewState

	id     string
	buffer *B.Buffer

	width           float32
	height          float32
	charWidth       float32
	charHeight      float32
	addressBarWidth float32
}

func HexView(id string, b *B.Buffer, st *ViewState) *HexViewWidget {
	h := &HexViewWidget{id: id, buffer: b}
	h.state = st
	return h
}

func bytesPerLine(width, charwidth float32) int {
	//to display 1 byte takes 4 characters: 2 for hexdump, 1 trailing space and 1 print
	maxChars := int(width / (4 * charwidth))

	//round to multiple of 4
	maxChars -= maxChars % 4
	if maxChars == 0 {
		maxChars = 4
	}
	return maxChars
}

//update transient state variables and helpers
func (h *HexViewWidget) update() {
	h.width, h.height = G.GetAvailableRegion()
	sz := I.CalcTextSize("F", true, 0)
	h.charWidth, h.charHeight = sz.X, sz.Y

	size := h.buffer.Size()
	nDigits := numHexDigits(size)
	h.addressBarWidth, _ = G.CalcTextSize(addrLabel(size, nDigits))

	h.state.bytesPerLine = int64(bytesPerLine(h.width-h.addressBarWidth, h.charWidth))
	h.state.linesPerScreen = int64(h.height / h.charHeight)

	h.state.topAddr = int64(I.ScrollY()/h.charHeight) * h.state.bytesPerLine

	if h.state.shouldScroll {
		h.ScrollTo(h.state.scrollToAddr)
		h.state.shouldScroll = false
	}
}

func (h *HexViewWidget) onScreen(addr int64) bool {
	top := h.state.topAddr
	fin := top + h.state.bytesPerLine*h.state.linesPerScreen
	return addr >= top && addr < fin
}

//to make the cut&copy keys work with a selection or the item under cursor
func (h *HexViewWidget) selectMinimal1(f func()) func() {
	return func() {
		if h.state.cursor < h.buffer.Size() && h.state.selectionSize == 0 {
			h.state.selectionStart = h.state.cursor
			h.state.selectionSize++
		}
		f()
	}
}

//I would like to have this in main.go?
func (h *HexViewWidget) handleKeys() {
	keymap := map[G.Key]func(){
		//movement
		G.KeyDown:  func() { actionMove(-h.state.bytesPerLine) },
		G.KeyK:     func() { actionMove(-h.state.bytesPerLine) },
		G.KeyUp:    func() { actionMove(+h.state.bytesPerLine) },
		G.KeyJ:     func() { actionMove(+h.state.bytesPerLine) },
		G.KeyLeft:  func() { actionMove(-1) },
		G.KeyH:     func() { actionMove(-1) },
		G.KeyRight: func() { actionMove(+1) },
		G.KeyL:     func() { actionMove(+1) },
		G.KeyG:     actionGoto,

		//edit
		G.KeyX: h.selectMinimal1(actionCut),
		G.KeyY: h.selectMinimal1(actionCopy),
		G.KeyP: actionPaste,
		G.KeyI: func() { h.state.SetSelection(0, 0); h.state.editmode = InsertMode },
		G.KeyO: func() { h.state.SetSelection(0, 0); h.state.editmode = OverwriteMode },
		G.KeyU: actionUndo,
		G.KeyR: actionRedo,
	}
	//other modes are handled by the edit-input-widget in the hex dump
	if h.state.editmode == NormalMode && G.IsWindowFocused(G.FocusedFlagsNone) {
		for k, f := range keymap {
			if G.IsKeyPressed(k) {
				f()
			}
		}
	}
}

func printByte(b byte) string {
	if unicode.IsGraphic(rune(b)) {
		return string(b)
	} else {
		return "."
	}
}

func (h *HexViewWidget) updateSelection(addr int64) {
	s := h.state
	if h.buffer.Size() <= 0 {
		s.selectionSize = 0
		s.selectionStart = 0
	}

	if s.dragging {
		//update mouse drag
		if addr < s.dragstart {
			s.selectionStart = addr
			s.selectionSize = s.dragstart - addr + 1
		} else {
			s.selectionStart = s.dragstart
			s.selectionSize = addr - s.dragstart + 1
		}
	} else {
		//shift-click
		off, size := s.Selection()
		if size == 0 {
			//new selection from cursor to click addr
			if addr < s.cursor {
				s.selectionStart = addr
				s.selectionSize = s.cursor - addr + 1
			} else {
				s.selectionStart = s.cursor
				s.selectionSize = addr - s.cursor
			}
		} else {
			//update either the begin or the end of an existing selection with the click addr
			if addr-off < off+size-addr { //addr is closer to start than end of selection?
				//shift bottom of selection to include addr
				s.selectionStart = addr
				s.selectionSize += off - addr
			} else {
				//shift end of selection to include up to addr
				extra := addr - (off + size) + 1
				s.selectionSize += extra
			}
		}
	}

	//We allow EOF to be a selectable/editable field and such, chomp it here
	if s.inSelection(h.buffer.Size()) {
		s.selectionSize = h.buffer.Size() - s.selectionStart
	}
}

func (h *HexViewWidget) printBG(addr int64, cursorw, selectw int) {
	canvas := G.GetCanvas()
	pos := G.GetCursorScreenPos()

	if addr == h.state.cursor {
		cursorBG := color.RGBA{R: 255, G: 100, B: 000, A: 255}
		if h.state.editmode != NormalMode {
			cursorBG.B = 255
		}
		rect := image.Pt(cursorw*int(h.charWidth), int(h.charHeight))
		canvas.AddRectFilled(pos, pos.Add(rect), cursorBG, 0, 0)
	}

	if h.state.inSelection(addr) {
		selectionBG := color.RGBA{R: 50, G: 30, B: 150, A: 100}
		rect := image.Pt(selectw*int(h.charWidth), int(h.charHeight))
		canvas.AddRectFilled(pos, pos.Add(rect), selectionBG, 0, 0)
	}
}

//little hack for vec2
func vec2Abs(v I.Vec2) float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

func mouseMoved() bool {
	delta := vec2Abs(G.Context.IO().GetMouseDelta())
	//fmt.Println("mousedelta", delta)
	return delta > 0
}

//make an input-cell
func (h *HexViewWidget) BuildInput(addr int64) {
	h.printBG(addr, 2, 3)
	InputHex("inputhex", h.cancelInput, h.advanceInput).Build()
}

//a cell is a piece of text that corresponds to a file-offset.
//it can be clicked and dragged
func (h *HexViewWidget) BuildCell(addr int64, txt string) {
	if len(txt) == 1 {
		h.printBG(addr, 1, 1)
	} else {
		h.printBG(addr, 2, 3)
	}
	I.Text(txt)

	//Mouse handling, clicks and drags
	//XXX all of this should be moved to the parent widget for efficiency
	//but is it a bit tedious to calculate what the mouse-position represents
	if !G.IsItemHovered() {
		return
	}
	if h.state.dragging {
		h.updateSelection(addr)
	}
	if G.IsMouseDown(G.MouseButtonLeft) {
		h.state.editmode = NormalMode
		if !h.state.dragging {
			if G.IsKeyDown(G.KeyLeftShift) || G.IsKeyDown(G.KeyRightShift) {
				h.updateSelection(addr)
			} else if mouseMoved() {
				h.state.dragstart = addr
				h.state.selectionStart = addr
				h.state.selectionSize = 0
				h.state.dragging = true
			} else {
				h.state.selectionSize = 0
			}
		}
		h.state.cursor = addr // should be updated after call to updateSelection
		if h.state.dragging && addr == h.buffer.Size() {
			h.state.cursor = addr - 1
		}
	}
	if G.IsMouseReleased(G.MouseButtonLeft) {
		h.state.dragging = false
	}
}

//to be called from Build() function, prints hexdump of byte and handles keyclicks and such
func (h *HexViewWidget) BuildHexCell(addr int64, b byte) {
	var hex string

	//put 1 'empty` box at EOF to be able to put the cursor there (for appending to a file)
	endAddr := h.buffer.Size()
	if h.onScreen(h.state.cursor) && h.state.editmode == InsertMode {
		endAddr++
	}
	if addr == endAddr {
		hex = "   "
	} else {
		hex = fmt.Sprintf("%02X ", b)
	}
	h.BuildCell(addr, hex)
}

//to be called from Build() function, prints readable interpretation of byte and handles keyclicks and such
func (h *HexViewWidget) BuildStrCell(addr int64, b byte) {
	str := printByte(b)
	h.BuildCell(addr, str)
}

func (h *HexViewWidget) Build() {
	//use a child widget with NoMove flags, so that dragging events gets passed to the
	//widget, instead of dragging the window
	G.Child().Border(false).Flags(G.WindowFlagsNoMove).Layout(
		G.Custom(h.printWidget),
		G.ContextMenu().Layout(menuEdit()),
	).Build()
}

//callback for edit-widget
func (h *HexViewWidget) cancelInput() {
	h.state.editmode = NormalMode
}

//callback for edit-widget
func (h *HexViewWidget) advanceInput(b byte) {
	if h.state.cursor >= h.buffer.Size() {
		h.state.editmode = InsertMode
	}
	switch h.state.editmode {
	case InsertMode:
		actionInsert(b)
	case OverwriteMode:
		actionOverWrite(b)
	}
}

func (h *HexViewWidget) printWidget() {
	I.PushStyleVarVec2(I.StyleVarFramePadding, I.Vec2{})
	I.PushStyleVarVec2(I.StyleVarItemSpacing, I.Vec2{})
	I.PushStyleVarVec2(I.StyleVarCellPadding, I.Vec2{})
	defer I.PopStyleVarV(3)

	h.update()
	h.handleKeys() //XXX this should be somewhere else??

	flags := I.TableFlags_BordersOuter | I.TableFlags_SizingFixedFit
	if I.BeginTable("HexDumpTable", 3, flags, I.Vec2{}, 0) {
		defer I.EndTable()
		I.TableSetupColumn("Offset", 0, h.addressBarWidth, 0)
		I.TableSetupColumn("HexDump", 0, 3*h.charWidth*float32(h.state.bytesPerLine), 0)
		I.TableSetupColumn("Readable", 0, float32(h.state.bytesPerLine)*h.charWidth, 0)

		//XXX this is simple solution, it was a bit difficult to put the logic for 3 different edit modes
		//into 1 dump function and to take care of only printing 1 edit-input per window,
		//so i have 3 big functions that basically do the same...

		//XXX it is necessary to take care only 1 edit-input gets print per hex dump.
		//because input is not cleared (consumed), then the next edit would get the input also :))

		//for the 1st implementation split into 3 different functions.
		//should be unified in the future!!
		switch h.state.editmode {
		case NormalMode:
			h.printNormalDump()
		case OverwriteMode:
			h.printOverWriteDump()
		case InsertMode:
			h.printInsertDump()
		}
	}
}

func (h *HexViewWidget) printInsertDump() {
	//print the hex dump using a listclipper
	numLines := (h.buffer.Size() + h.state.bytesPerLine - 1) / h.state.bytesPerLine
	lineBuffer := make([]byte, int(h.state.bytesPerLine)) //buffer to read the bytes for 1 line
	maxAddr := numHexDigits(h.buffer.Size())              //saved for printing address

	var seenCursor = false //the input-handling means
	var clip I.ListClipper
	//dumb hack: do numlines + 10 because on big files, the last few lines get chopped off
	//due to floating point errors in scrolling calculations
	clip.BeginV(int(numLines+10), h.charHeight)
	defer clip.End()
	for clip.Step() {
		for lnum := clip.DisplayStart; lnum < clip.DisplayEnd; lnum++ {
			offs := int64(lnum) * h.state.bytesPerLine
			if offs > h.buffer.Size() {
				break
			}

			//read data for this line
			if seenCursor {
				h.buffer.Seek(offs-1, io.SeekStart)
			} else {
				h.buffer.Seek(offs, io.SeekStart)
			}
			n, e := h.buffer.Read(lineBuffer)
			if e != nil && e != io.EOF {
				panic(e) //XXX not very elegant
			}

			//address
			I.TableNextColumn()
			I.Text(addrLabel(offs, maxAddr))

			//hex dump
			var cursorOffset = 0 //if the cursor is in this line, linuBuffer reads should offset
			I.TableNextColumn()
			for i := 0; i < n; i++ {
				if i != 0 {
					I.SameLine()
				}
				addr := offs + int64(i)
				if addr == h.state.cursor && !seenCursor {
					seenCursor = true
					cursorOffset = 1
					h.BuildInput(addr)
				} else {
					h.BuildHexCell(addr-int64(cursorOffset), lineBuffer[i-cursorOffset])
				}
			}

			//allow to select EOF
			if n != int(h.state.bytesPerLine) {
				if n != 0 {
					I.SameLine()
				}
				addr := offs + int64(n)
				if addr == h.state.cursor && !seenCursor {
					seenCursor = true
					cursorOffset = 1
					h.BuildInput(addr)
				} else {
					if cursorOffset == 0 {
						//just EOF
						h.BuildHexCell(addr, 0)
					} else {
						h.BuildHexCell(addr, lineBuffer[n-1])
					}
				}
			}

			//readable string
			I.TableNextColumn()
			cursorOffset = 0
			for i := 0; i < n; i++ {
				if i != 0 {
					I.SameLine()
				}
				addr := offs + int64(i)
				if addr == h.state.cursor {
					cursorOffset = 1
					h.BuildStrCell(addr, 0)
				} else {
					h.BuildStrCell(addr, lineBuffer[i-cursorOffset])
				}
			}
		}
	}
}

func (h *HexViewWidget) printOverWriteDump() {
	//print the hex dump using a listclipper
	numLines := (h.buffer.Size() + h.state.bytesPerLine - 1) / h.state.bytesPerLine
	lineBuffer := make([]byte, int(h.state.bytesPerLine)) //buffer to read the bytes for 1 line
	maxAddr := numHexDigits(h.buffer.Size())              //saved for printing address

	var seenCursor = false //the input-handling means
	var clip I.ListClipper
	//dumb hack: do numlines + 10 because on big files, the last few lines get chopped off
	//due to floating point errors in scrolling calculations
	clip.BeginV(int(numLines+10), h.charHeight)
	defer clip.End()
	for clip.Step() {
		for lnum := clip.DisplayStart; lnum < clip.DisplayEnd; lnum++ {
			offs := int64(lnum) * h.state.bytesPerLine
			if offs > h.buffer.Size() {
				break
			}

			//read data for this line
			h.buffer.Seek(offs, io.SeekStart)
			n, e := h.buffer.Read(lineBuffer)
			if e != nil && e != io.EOF {
				panic(e) //XXX not very elegant
			}

			//address
			I.TableNextColumn()
			I.Text(addrLabel(offs, maxAddr))

			//hex dump
			I.TableNextColumn()
			for i := 0; i < n; i++ {
				if i != 0 {
					I.SameLine()
				}
				addr := offs + int64(i)
				if addr == h.state.cursor && !seenCursor {
					seenCursor = true
					h.BuildInput(addr)
				} else {
					h.BuildHexCell(addr, lineBuffer[i])
				}
			}
			//allow to select EOF
			if n != int(h.state.bytesPerLine) {
				if n != 0 {
					I.SameLine()
				}
				addr := offs + int64(n)
				if addr == h.state.cursor && !seenCursor {
					seenCursor = true
					h.BuildInput(addr)
				} else {
					h.BuildHexCell(addr, 0)
				}
			}

			//readable string
			I.TableNextColumn()
			for i := 0; i < n; i++ {
				if i != 0 {
					I.SameLine()
				}
				h.BuildStrCell(offs+int64(i), lineBuffer[i])
			}
		}
	}
}

func (h *HexViewWidget) printNormalDump() {
	//print the hex dump using a listclipper
	numLines := (h.buffer.Size() + h.state.bytesPerLine - 1) / h.state.bytesPerLine
	lineBuffer := make([]byte, int(h.state.bytesPerLine)) //buffer to read the bytes for 1 line
	maxAddr := numHexDigits(h.buffer.Size())              //saved for printing address

	var clip I.ListClipper
	//dumb hack: do numlines + 10 because on big files, the last few lines get chopped off
	//due to floating point errors in scrolling calculations
	clip.BeginV(int(numLines+10), h.charHeight)
	defer clip.End()
	for clip.Step() {
		for lnum := clip.DisplayStart; lnum < clip.DisplayEnd; lnum++ {
			offs := int64(lnum) * h.state.bytesPerLine
			if offs > h.buffer.Size() {
				break
			}

			//read data for this line
			h.buffer.Seek(offs, io.SeekStart)
			n, e := h.buffer.Read(lineBuffer)
			if e != nil && e != io.EOF {
				panic(e) //XXX not very elegant
			}

			//address
			I.TableNextColumn()
			I.Text(addrLabel(offs, maxAddr))

			//hex dump
			I.TableNextColumn()
			for i := 0; i < n; i++ {
				if i != 0 {
					I.SameLine()
				}
				addr := offs + int64(i)
				h.BuildHexCell(addr, lineBuffer[i])
			}

			//allow to select EOF
			if n != int(h.state.bytesPerLine) {
				if n != 0 {
					I.SameLine()
				}
				addr := offs + int64(n)
				h.BuildHexCell(addr, 0)
			}

			//readable string
			I.TableNextColumn()
			for i := 0; i < n; i++ {
				if i != 0 {
					I.SameLine()
				}
				h.BuildStrCell(offs+int64(i), lineBuffer[i])
			}
		}
	}
}

func (h *HexViewWidget) ScrollTo(addr int64) {
	bpl := h.state.bytesPerLine
	top := h.state.topAddr
	switch {
	case h.state.cursor < top:
		//scroll up, addr should be in the first line
		//make first line the one that contains addr
		top = (h.state.cursor / bpl) * bpl

	case h.state.cursor > top+bpl*h.state.linesPerScreen:
		//scroll down, addr should be in the last line
		a := h.state.cursor - h.state.linesPerScreen*bpl
		top = ((a + bpl - 1) / bpl) * bpl
	default:
		//addr is already on screen
	}
	/* clamp address */
	if top < 0 {
		top = 0
	}
	if top > h.buffer.Size() {
		top = h.buffer.Size()
	}
	I.SetScrollY(float32(float64(top/bpl) * float64(h.charHeight)))
	h.state.topAddr = top
}
