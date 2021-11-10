package imgui

import (
	"fmt"
	"runtime"
	"testing"
)

func TestIDs(t *testing.T) {
	fmt.Println(ImHashStr("Hello World", 0, 0))
	fmt.Println(ImHashStr("Beep Boop", 0, 0))
}

func BenchmarkImguiIDs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ImHashStr("Hello World", 0, 0)
	}
}

var pc [1]uintptr

func BenchmarkProgramCounterIDs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.Callers(1, pc[:])
	}
}
