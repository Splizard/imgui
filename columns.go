package imgui

func PushColumnsBackground() {
	var window = GetCurrentWindowRead()
	var columns = window.DC.CurrentColumns
	if columns.Count == 1 {
		return
	}

	// Optimization: avoid SetCurrentChannel() + PushClipRect()
	columns.HostBackupClipRect = window.ClipRect
	SetWindowClipRectBeforeSetChannel(window, &columns.HostInitialClipRect)
	columns.Splitter.SetCurrentChannel(window.DrawList, 0)
}
