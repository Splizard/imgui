package imgui

import (
	"unsafe"
)

/*
	The storage functions have been modified to use an ordinary Go map,
	this is probably inefficient, but convinient for now - Quentin.
*/

// ImGuiStoragePair [Internal]
type ImGuiStoragePair struct {
	key ImGuiID
	val int
}

func (this *ImGuiStoragePair) SetInt(key ImGuiID, val int) {
	this.key = key
	this.val = val
}

func (this *ImGuiStoragePair) SetFloat(key ImGuiID, val float) {
	this.key = key
	*(*float)(unsafe.Pointer(&this.val)) = val
}

// ImGuiStorage Helper: Key->Value storage
// Typically you don't have to worry about this since a storage is held within each Window.
// We use it to e.guiContext. store collapse state for a tree (Int 0/1)
// This is optimized for efficient lookup (dichotomy into a contiguous buffer) and rare insertion (typically tied to user interactions aka max once a frame)
// You can use it as custom user storage for temporary values. Declare your own storage if, for example:
// - You want to manipulate the open/close state of a particular sub-tree in your interface (tree node uses Int 0/1 to store their state).
// - You want to store custom debug data easily without adding or editing structures in your code (probably not efficient, but convenient)
// Types are NOT stored, so it is up to you to make sure your Key don't collide with different types.
type ImGuiStorage struct {
	Data     map[ImGuiID]int
	Pointers map[ImGuiID]any
}

func (this *ImGuiStorage) Clear() {
	this.Data = map[ImGuiID]int{}
	this.Pointers = map[ImGuiID]any{}
}

func (this *ImGuiStorage) GetInt(key ImGuiID, default_val int) int {
	val, ok := this.Data[key]
	if !ok {
		return default_val
	}
	return val
}

func (this *ImGuiStorage) SetInt(key ImGuiID, val int) {
	if this.Data == nil {
		this.Data = map[ImGuiID]int{}
	}
	this.Data[key] = val
}

func (this *ImGuiStorage) GetBool(key ImGuiID, default_val bool) bool {
	var def int
	if default_val {
		def = 1
	}
	return this.GetInt(key, def) != 0
}

func (this *ImGuiStorage) SetBool(key ImGuiID, val bool) {
	if val {
		this.SetInt(key, 1)
	} else {
		this.SetInt(key, 0)
	}
}

func (this *ImGuiStorage) GetFloat(key ImGuiID, default_val float) float {
	val := this.GetInt(key, *(*int)(unsafe.Pointer(&default_val)))
	return *(*float)(unsafe.Pointer(&val))
}

func (this *ImGuiStorage) SetFloat(key ImGuiID, val float) {
	this.SetInt(key, *(*int)(unsafe.Pointer(&val)))
}

func (this *ImGuiStorage) GetInterface(key ImGuiID) any {
	return this.Pointers[key]
}

func (this *ImGuiStorage) SetInterface(key ImGuiID, val any) {
	if this.Pointers == nil {
		this.Pointers = make(map[ImGuiID]any)
	}
	this.Pointers[key] = val
}

// SetAllInt Use on your own storage if you know only integer are being stored (open/close all tree nodes)
func (this *ImGuiStorage) SetAllInt(val int) {
	for key := range this.Data {
		this.Data[key] = val
	}
}
