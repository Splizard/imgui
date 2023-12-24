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
	var was_disabled = (g.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0
	if !was_disabled && disabled {
		g.DisabledAlphaBackup = g.Style.Alpha
		g.Style.Alpha *= g.Style.DisabledAlpha // PushStyleVar(ImGuiStyleVar_Alpha, g.Style.Alpha * g.Style.DisabledAlpha);
	}
	if was_disabled || disabled {
		g.CurrentItemFlags |= ImGuiItemFlags_Disabled
	}
	g.ItemFlagsStack = append(g.ItemFlagsStack, g.CurrentItemFlags)
}

func EndDisabled() {
	var was_disabled = (g.CurrentItemFlags & ImGuiItemFlags_Disabled) != 0
	//PopItemFlag();
	g.ItemFlagsStack = g.ItemFlagsStack[:len(g.ItemFlagsStack)-1]
	g.CurrentItemFlags = g.ItemFlagsStack[len(g.ItemFlagsStack)-1]
	if was_disabled && (g.CurrentItemFlags&ImGuiItemFlags_Disabled) == 0 {
		g.Style.Alpha = g.DisabledAlphaBackup //PopStyleVar();
	}
}
