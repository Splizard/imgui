package imgui

type double = float64
type int = int32
type uint = uint32
type float = float32
type size_t = uintptr
type char = byte

/*
func isfalse(x int) bool {
	return x == 0
}
*/

func istrue(x int) bool {
	return x != 0
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

const (
	IM_S8_MIN  = -128
	IM_S8_MAX  = 127
	IM_U8_MIN  = 0
	IM_U8_MAX  = 255
	IM_S16_MIN = -32768
	IM_S16_MAX = 32767
	IM_U16_MIN = 0
	IM_U16_MAX = 65535
	IM_S32_MIN = -2147483648
	IM_S32_MAX = 2147483647
	IM_U32_MIN = 0
	IM_U32_MAX = 4294967295
	IM_S64_MIN = -9223372036854775808
	IM_S64_MAX = 9223372036854775807
	IM_U64_MIN = 0
	IM_U64_MAX = 18446744073709551615
)
