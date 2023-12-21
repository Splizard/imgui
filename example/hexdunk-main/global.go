package main

//global variables and definitions

import (
	"fmt"
	"io/fs"

	B "github.com/snhmibby/filebuf"
)

const (
	ProgramName = "HexDunk"

	//dialog ids
	DialogOpen   = "Open"         //fileDialog, callback: actionOpen
	DialogSaveAs = "Save As"      //fileDialog, callback: actionWriteFile
	DialogGoto   = "Goto Address" //intDialog,  callback: actionGotoAddr
)

//should not be an opaque closure, but a struct with an action-enum.
//that way, we can adjust copy-buffers when a file gets saved (i.e. referenced
//portions of the to-be-saved file should be removed through all buffers throughout
//the program on a save, so that they don't 'change' when the file gets written)
type Undo struct {
	undo, redo func() (int64, int64) //return affected region
}

//an opened file
type HexFile struct {
	name       string
	buf        *B.Buffer
	dirty      bool
	stats      fs.FileInfo
	undo, redo []Undo
}

//each tab is a view on an opened file
type HexTab struct {
	name string
	view *ViewState
}

type ViewState struct {
	//generic state
	cursor         int64 //address (byte offset in file)
	topAddr        int64 //address on top of the screen
	bytesPerLine   int64 //number of 'dunked' bytes per line
	linesPerScreen int64 //number of lines per screen
	editmode       editMode

	//current selection
	selectionStart, selectionSize int64

	//selection mouse dragging
	dragging  bool
	dragstart int64

	//update scroll-position (communicate to imgui scrollbar during widget build)
	shouldScroll bool
	scrollToAddr int64
}

type editMode int

const (
	NormalMode editMode = iota
	InsertMode
	OverwriteMode
)

//global variables
type Globals struct {
	// All opened files, index by file-path
	Files map[string]*HexFile

	// All tabs (every tab is a view on an opened file)
	Tabs []HexTab

	//Index of active tab in display
	ActiveTab int

	//Current copy/paste buffer
	//TODO: could be something nice, a circular buffer, named buffers, etc
	ClipBoard *B.Buffer
}

var HD Globals = Globals{
	Tabs:      make([]HexTab, 0),
	ActiveTab: -1,
	Files:     make(map[string]*HexFile),
}

func ActiveTab() *HexTab {
	if HD.ActiveTab >= 0 {
		//consistency check
		if ActiveFile() == nil {
			panic("impossible")
		}
		return &HD.Tabs[HD.ActiveTab]
	}
	return nil
}

func ActiveFile() *HexFile {
	if HD.ActiveTab >= 0 {
		hf, ok := HD.Files[HD.Tabs[HD.ActiveTab].name]
		if !ok {
			panic("tab opened on closed file")
		}
		return hf
	}
	return nil
}

/* Tab methods */

//set the tab cursor to addr (make sure it is in file-range) and scroll to it
func (tab *HexTab) setCursor(addr int64) {
	hf, ok := HD.Files[HD.Tabs[HD.ActiveTab].name]
	if !ok {
		panic("tab opened on closed file")
	}
	hf.ClampAddr(&addr)
	tab.view.cursor = addr
	tab.view.ScrollTo(addr)
}

/* ViewState methods (should/could also be hextab* methods */

func (view *ViewState) SetSelection(begin, size int64) {
	view.selectionStart, view.selectionSize = begin, size
}

func (view *ViewState) Selection() (begin, size int64) {
	return view.selectionStart, view.selectionSize
}

func (st *ViewState) inSelection(addr int64) bool {
	off, size := st.Selection()
	return addr >= off && addr < off+size
}

func (st *ViewState) ScrollTo(addr int64) {
	st.shouldScroll = true
	st.scrollToAddr = addr
}

/* some utility functions */

//mkErr will create a properly formatted error message
func mkErr(msg string, e error) error {
	err := fmt.Errorf("%s: %v", msg, e)
	return err
}
