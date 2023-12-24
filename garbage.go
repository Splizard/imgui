package imgui

// GcCompactTransientMiscBuffers Garbage collection
func GcCompactTransientMiscBuffers() {
	g.ItemFlagsStack = nil
	g.GroupStack = nil
	TableGcCompactSettings()
}

// GcCompactTransientWindowBuffers Free up/compact internal window buffers, we can use this when a window becomes unused.
// Not freed:
// - ImGuiWindow, ImGuiWindowSettings, Name, StateStorage, ColumnsStorage (may hold useful data)
// This should have no noticeable visual effect. When the window reappear however, expect new allocation/buffer growth/copy cost.
func GcCompactTransientWindowBuffers(window *ImGuiWindow) {
	window.MemoryCompacted = true
	window.MemoryDrawListIdxCapacity = int(cap(window.DrawList.IdxBuffer))
	window.MemoryDrawListVtxCapacity = int(cap(window.DrawList.VtxBuffer))
	window.IDStack = nil
	window.DrawList._ClearFreeMemory()
	window.DC.ChildWindows = nil
	window.DC.ItemWidthStack = nil
	window.DC.TextWrapPosStack = nil
}

func GcAwakeTransientWindowBuffers(window *ImGuiWindow) {
	// We stored capacity of the ImDrawList buffer to reduce growth-caused allocation/copy when awakening.
	// The other buffers tends to amortize much faster.
	window.MemoryCompacted = false

	//reserve
	window.DrawList.IdxBuffer = append(window.DrawList.IdxBuffer, make([]ImDrawIdx, window.MemoryDrawListIdxCapacity-int(len(window.DrawList.IdxBuffer)))...)
	window.DrawList.VtxBuffer = append(window.DrawList.VtxBuffer, make([]ImDrawVert, window.MemoryDrawListVtxCapacity-int(len(window.DrawList.VtxBuffer)))...)

	window.MemoryDrawListIdxCapacity = 0
	window.MemoryDrawListVtxCapacity = 0
}
