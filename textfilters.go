package imgui

import (
	"strings"
)

// Helper: Growable text buffer for logging/accumulating text
// (this could be called 'ImGuiTextBuilder' / 'ImGuiStringBuilder')
type ImGuiTextBuffer []byte

// Helper: Parse and apply text filters. In format "aaaaa[,bbbb][,ccccc]"
type ImGuiTextFilter struct {
	InputBuf  []byte
	Filters   []string
	CountGrep int
}

func NewImGuiTextFilter(default_filter string) ImGuiTextFilter {
	if default_filter != "" {
		filter := ImGuiTextFilter{InputBuf: []byte(default_filter)}
		filter.Build()
		return filter
	}
	return ImGuiTextFilter{}
}

func (this *ImGuiTextFilter) Draw(label string /*= "Filter (inc,-exc)"*/, width float) bool {
	if width != 0.0 {
		SetNextItemWidth(width)
	}
	var value_changed = InputText(label, &this.InputBuf, 0, nil, nil)
	if value_changed {
		this.Build()
	}
	return value_changed
}

// Helper calling InputText+Build
func (this *ImGuiTextFilter) PassFilter(text string) bool {
	if len(this.Filters) == 0 {
		return true
	}

	for _, f := range this.Filters {
		if len(f) == 0 {
			continue
		}
		if f[0] == '-' {
			if strings.Contains(text, f[1:]) {
				return false
			}
		} else {
			// Grep
			if strings.Contains(text, f) {
				return false
			}
		}
	}

	// Implicit * grep
	if this.CountGrep == 0 {
		return true
	}

	return false
}

func (this *ImGuiTextFilter) Build() {
	this.Filters = strings.Split(string(this.InputBuf), ",")

	this.CountGrep = 0
	for i := range this.Filters {
		this.Filters[i] = strings.TrimSpace(this.Filters[i])
		if len(this.Filters[i]) == 0 {
			continue
		}
		if this.Filters[i][0] != '-' {
			this.CountGrep += 1
		}
	}
}

func (this *ImGuiTextFilter) Clear() {
	this.InputBuf[0] = 0
	this.Build()
}

func (this *ImGuiTextFilter) IsActive() bool { return len(this.Filters) > 0 }
