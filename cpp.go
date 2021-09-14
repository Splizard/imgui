package imgui

type double = float64
type int = int32
type uint = uint32
type float = float32
type size_t = uintptr
type char = byte

func isfalse(x int) bool {
	return x == 0
}

func istrue(x int) bool {
	return x != 0
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
