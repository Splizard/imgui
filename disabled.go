package imgui

// Disabling [BETA API]
// - Disable all user interactions and dim items visuals (applying style.DisabledAlpha over current colors)
// - Those can be nested but it cannot be used to enable an already disabled section (a single BeginDisabled(true) in the stack is enough to keep everything disabled)
// - BeginDisabled(false) essentially does nothing useful but is provided to facilitate use of boolean expressions. If you can a calling BeginDisabled(False)/EndDisabled() best to a it.

// BeginDisabled BeginDisabled()/EndDisabled()
// - Those can be nested but it cannot be used to enable an already disabled section (a single BeginDisabled(true) in the stack is enough to keep everything disabled)
// - Visually this is currently altering alpha, but it is expected that in a future styling system this would work differently.
// - Feedback welcome at https://github.com/ocornut/imgui/issues/211
// - BeginDisabled(false) essentially does nothing useful but is provided to facilitate use of boolean expressions. If you can avoid calling BeginDisabled(False)/EndDisabled() best to avoid it.
// - Optimized shortcuts instead of PushStyleVar() + PushItemFlag()
func BeginDisabled(disabled bool /*= true*/) {
	var was_disabled = (guiContext.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0
	if !was_disabled && disabled {
		guiContext.DisabledAlphaBackup = guiContext.Style.Alpha
		guiContext.Style.Alpha *= guiContext.Style.DisabledAlpha // PushStyleVar(ImGuiStyleVar_Alpha, guiContext.Style.Alpha * guiContext.Style.DisabledAlpha);
	}
	if was_disabled || disabled {
		guiContext.CurrentItemFlags |= ImGuiItemFlags_Disabled
	}
	guiContext.ItemFlagsStack = append(guiContext.ItemFlagsStack, guiContext.CurrentItemFlags)
}

func EndDisabled() {
	var was_disabled = (guiContext.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0
	//PopItemFlag();
	guiContext.ItemFlagsStack = guiContext.ItemFlagsStack[:len(guiContext.ItemFlagsStack)-1]
	guiContext.CurrentItemFlags = guiContext.ItemFlagsStack[len(guiContext.ItemFlagsStack)-1]
	if was_disabled && (guiContext.CurrentItemFlags&ImGuiItemFlags_Disabled) == 0 {
		guiContext.Style.Alpha = guiContext.DisabledAlphaBackup //PopStyleVar();
	}
}
