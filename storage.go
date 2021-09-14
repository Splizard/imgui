package imgui

import (
	"unsafe"
)

// [Internal]
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

// Helper: Key->Value storage
// Typically you don't have to worry about this since a storage is held within each Window.
// We use it to e.g. store collapse state for a tree (Int 0/1)
// This is optimized for efficient lookup (dichotomy into a contiguous buffer) and rare insertion (typically tied to user interactions aka max once a frame)
// You can use it as custom user storage for temporary values. Declare your own storage if, for example:
// - You want to manipulate the open/close state of a particular sub-tree in your interface (tree node uses Int 0/1 to store their state).
// - You want to store custom debug data easily without adding or editing structures in your code (probably not efficient, but convenient)
// Types are NOT stored, so it is up to you to make sure your Key don't collide with different types.
type ImGuiStorage struct {
	Data     map[ImGuiID]int
	Pointers map[ImGuiID]interface{}
}

func (this *ImGuiStorage) Clear() {
	this.Data = map[ImGuiID]int{}
	this.Pointers = map[ImGuiID]interface{}{}
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
	panic("not implemented")
}

func (this *ImGuiStorage) SetBool(key ImGuiID, val bool) {
	panic("not implemented")
}

func (this *ImGuiStorage) GetFloat(key ImGuiID, default_val float) float {
	panic("not implemented")
}

func (this *ImGuiStorage) SetFloat(key ImGuiID, val float) {
	panic("not implemented")
}

func (this *ImGuiStorage) GetInterface(key ImGuiID) interface{} {
	return this.Pointers[key]
}

func (this *ImGuiStorage) SetInterface(key ImGuiID, val interface{}) {
	if this.Pointers == nil {
		this.Pointers = make(map[ImGuiID]interface{})
	}
	this.Pointers[key] = val
}

// - Get***Ref() functions finds pair, insert on demand if missing, return pointer. Useful if you intend to do Get+Set.
// - References are only valid until a new value is added to the storage. Calling a Set***() function or a Get***Ref() function invalidates the pointer.
// - A typical use case where this is convenient for quick hacking (e.g. add storage during a live Edit&Continue session if you can't modify existing struct)
//      float* pvar ImGui::GetFloatRef(key) = ImGui::SliderFloat("var", pvar, 100.0f) 0, some_var *pvar +=

func (this *ImGuiStorage) GetIntRef(key ImGuiID, default_val int) *int {
	panic("not implemented")
}

func (this *ImGuiStorage) GetBoolRef(key ImGuiID, default_val bool) bool {
	panic("not implemented")
}

func (this *ImGuiStorage) GetFloatRef(key ImGuiID, default_val float) float {
	panic("not implemented")
}

func (this *ImGuiStorage) GetInterfaceRef(key ImGuiID, default_val interface{}) *interface{} {
	panic("not implemented")
}

// Use on your own storage if you know only integer are being stored (open/close all tree nodes)
func (this *ImGuiStorage) SetAllInt(val int) { panic("not implemented") }

// For quicker full rebuild of a storage (instead of an incremental one), you may add all your contents and then sort once.
func (this *ImGuiStorage) BuildSortByKey() {}
