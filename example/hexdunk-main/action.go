package main

//editor user actions that touch/need the activefile.

import (
	"fmt"
	"io"
	"os"

	"github.com/snhmibby/filebuf"
)

func actionGoto() {
	IntDialog(DialogGoto)
}

func actionGotoAddr(addr int64) {
	tab := ActiveTab()
	if tab == nil {
		panic("Goto: tab or file is nil (shouldn't happen)")
	}
	tab.setCursor(addr)
}

func actionMove(move int64) {
	tab := ActiveTab()
	if tab == nil {
		panic("Move: tab or file is nil (shouldn't happen)")
	}
	tab.setCursor(tab.view.cursor + move)
	tab.view.SetSelection(0, 0)
}

func actionRedo() {
	file := ActiveFile()
	if file == nil {
		panic("Redo: file is nil (shouldn't happen)")
	} else {
		file.Redo()
	}
}

func actionUndo() {
	file := ActiveFile()
	if file == nil {
		panic("Undo: file is nil (shouldn't happen)")
	} else {
		file.Undo()
	}
}

func actionInsert(b byte) {
	tab := ActiveTab()
	file := ActiveFile()
	if tab == nil || file == nil {
		panic("Insert: tab or file is nil (shouldn't happen)")
	}
	off := tab.view.cursor
	file.buf.Insert1(off, b)
	tab.view.SetSelection(0, 0)
	tab.setCursor(off + 1)
	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() (int64, int64) {
			file.Cut(off, 1)
			return off, 0
		},
		redo: func() (int64, int64) {
			file.buf.Insert1(off, b)
			return off, 0
		},
	})
}

func actionOverWrite(b byte) {
	tab := ActiveTab()
	file := ActiveFile()
	if tab == nil || file == nil {
		panic("Overwrite: tab or file is nil (shouldn't happen)")
	}
	off := tab.view.cursor
	file.buf.Seek(off, io.SeekStart)
	var overwritten_ = make([]byte, 1)
	_, err := file.buf.Read(overwritten_)
	if err != nil {
		ErrorDialog("Overwrite", fmt.Sprintf("Couldn't read from buffer: %v.", err))
		return
	}
	overwritten := overwritten_[0]
	file.buf.Remove(off, 1)
	file.buf.Insert1(off, b)
	tab.setCursor(off + 1)
	tab.view.SetSelection(0, 0)

	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() (int64, int64) {
			file.buf.Remove(off, 1)
			file.buf.Insert1(off, overwritten)
			return off, 0
		},
		redo: func() (int64, int64) {
			file.buf.Remove(off, 1)
			file.buf.Insert1(off, b)
			return off + 1, 0
		},
	})
}

func actionCut() {
	tab := ActiveTab()
	file := ActiveFile()
	if tab == nil || file == nil {
		panic("Cut: tab or file is nil (shouldn't happen)")
	}
	off, size := tab.view.Selection()
	cut, err := file.Cut(off, size) //XXX BUG this possibly creates hidden filetree copies (need to find them on saving!)
	if err != nil {
		ErrorDialog(fmt.Sprintf("Cut(%d, %d)", off, size), fmt.Sprint(err))
		return
	}

	HD.ClipBoard = cut
	tab.setCursor(off)
	tab.view.SetSelection(0, 0)
	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() (int64, int64) {
			file.buf.Paste(off, cut)
			return off, cut.Size()
		},
		redo: func() (int64, int64) {
			file.buf.Cut(off, size)
			return off, 0
		},
	})
}

func actionCopy() {
	file := ActiveFile()
	tab := ActiveTab()
	if tab == nil || file == nil {
		panic("Copy: tab or file is nil (shouldn't happen)")
	}
	off, size := tab.view.Selection()
	cpy, err := file.Copy(off, size)
	if err != nil {
		ErrorDialog(fmt.Sprintf("Copy(%d, %d)", off, size), fmt.Sprint(err))
	}
	HD.ClipBoard = cpy
	tab.setCursor(off)
	tab.view.SetSelection(off, 0)
}

//paste in front cursor
func actionPaste() {
	if HD.ClipBoard == nil {
		return //nothing to Paste
	}

	file := ActiveFile()
	tab := ActiveTab()
	if file == nil || tab == nil {
		panic("Paste: tab or file is nil (shouldn't happen)")
	}
	off := tab.view.cursor
	buf := HD.ClipBoard //XXX BUG this creates hidden copies of a file-based tree
	file.Paste(off, buf)
	tab.view.SetSelection(off, buf.Size())

	file.emptyRedo()
	file.addUndo(Undo{
		undo: func() (int64, int64) {
			file.buf.Cut(off, buf.Size())
			return off, 0
		},
		redo: func() (int64, int64) {
			file.buf.Paste(off, buf)
			return off, buf.Size()
		},
	})
}

func actionNewFile() {
	tmpPath, err := os.CreateTemp("", "NewFile*")
	if err != nil {
		ErrorDialog("NewFile", fmt.Sprintf("Cannot create tmp file: %v", err))
	}
	_, err = OpenHexFile(tmpPath.Name())
	if err != nil {
		ErrorDialog("NewFile", fmt.Sprintf("Couldn't open %s: %v", tmpPath.Name(), err))
	}
}

//callbacks for dialogs are set in the draw() layout function
func actionOpen(p string) {
	_, err := OpenHexFile(p)
	if err != nil {
		title := fmt.Sprintf("Opening File <%s>", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}
}

func actionOpenFile() {
	FileDialog(DialogOpen)
}

func actionWriteFile(p string) {
	hf := ActiveFile()
	f, err := os.CreateTemp("", "")
	if err != nil {
		title := fmt.Sprintf("Opening File <%s> for saving.", p)
		msg := fmt.Sprint(err)
		ErrorDialog(title, msg)
	}
	/*
	 *hf.buf.Seek(0, io.SeekStart)
	 *n, err := io.Copy(f, hf.buf)
	 *if err != nil || n != hf.buf.Size() {
	 *    os.Remove(f.Name())
	 *    title := fmt.Sprintf("Writing File <%s>", p)
	 *    msg := fmt.Sprintf("Written %d bytes (expected %d)\nError: %v", n, hf.buf.Size(), err)
	 *    ErrorDialog(title, msg)
	 *}
	 */
	//use the iter interface
	hf.buf.Iter(func(slice []byte) bool {
		var n int
		n, err = f.Write(slice)
		return n != len(slice) || err != nil
	})
	if err != nil {
		os.Remove(f.Name())
		title := fmt.Sprintf("Writing File <%s>", p)
		msg := fmt.Sprintf("Error: %v", err)
		ErrorDialog(title, msg)
	}
	err = os.Rename(f.Name(), p)
	if err != nil {
		os.Remove(f.Name())
		title := fmt.Sprintf("Naming File <%s>", p)
		msg := fmt.Sprintf("Couldn't rename tmp file <%s> to <%s>", f.Name(), p)
		ErrorDialog(title, msg)
	}

	//TODO XXX open a new buffer on the whole file again here?
	//this would refresh the working tree buffer
	hf.buf, err = filebuf.OpenFile(p)
	if err != nil {
		panic("TODO: handle OpenFile error after save")
	}
}

func actionSaveFile() {
	if ActiveFile() != nil {
		//TODO:
		//if not is tmp/new file {
		//   save to real file name
		//else use save-as dialog
		FileDialog(DialogSaveAs)
	}
}

func actionSaveAs() {
	if ActiveFile() != nil {
		FileDialog(DialogSaveAs)
	}
}

func actionCloseTab() {
	if HD.ActiveTab >= 0 {
		CloseTab(HD.ActiveTab)
	}
}

func actionQuit() {
	//TODO: "do you want to save unsaved changes..." dialog
	os.Exit(0)
}
