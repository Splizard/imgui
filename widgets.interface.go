package imgui

import (
	"fmt"
	"reflect"
	"unsafe"
)

var GDataTypeInfo = []ImGuiDataTypeInfo{
	{unsafe.Sizeof(int8(0)), "S8", "%v", "%v"}, // ImGuiDataType_S8
	{unsafe.Sizeof(uint8(0)), "U8", "%v", "%v"},
	{unsafe.Sizeof(int16(0)), "S16", "%v", "%v"}, // ImGuiDataType_S16
	{unsafe.Sizeof(uint16(0)), "U16", "%v", "%v"},
	{unsafe.Sizeof(int(0)), "S32", "%v", "%v"}, // ImGuiDataType_S32
	{unsafe.Sizeof(uint(0)), "U32", "%v", "%v"},
	{unsafe.Sizeof(ImS64(0)), "S64", "%v", "%v"}, // ImGuiDataType_S64
	{unsafe.Sizeof(ImU64(0)), "U64", "%v", "%v"},
	{unsafe.Sizeof(float(0)), "float", "%.3f", "%f"},  // ImGuiDataType_Float (float are promoted to double in va_arg)
	{unsafe.Sizeof(double(0)), "double", "%f", "%lf"}, // ImGuiDataType_Double
}

// PatchFormatStringFloatToInt FIXME-LEGACY: Prior to 1.61 our DragInt() function internally used floats and because of this the compile-time default value for format was "%.0f".
// Even though we changed the compile-time default, we expect users to have carried %f around, which would break the display of DragInt() calls.
// To honor backward compatibility we are rewriting the format string, unless IMGUI_DISABLE_OBSOLETE_FUNCTIONS is enabled. What could possibly go wrong?!
func PatchFormatStringFloatToInt(format string) string {
	if format[0] == '%' && format[1] == '.' && format[2] == '0' && format[3] == 'f' && format[4] == 0 { // Fast legacy path for "%.0f" which is expected to be the most common case.
		return "%v"
	}
	return format
}

func DataTypeFormatString(data_type ImGuiDataType, p_data any, format string) string {
	return fmt.Sprintf(format, p_data)
}

func DataTypeApplyOp(data_type ImGuiDataType, op int, output any, arg_1 any, arg_2 any) {
	//FIXME (porting) overflow handling was removed, need to add it back?
	IM_ASSERT(op == '+' || op == '-')
	switch data_type {
	case ImGuiDataType_S8:
		if op == '+' {
			*(output.(*ImS8)) = *(arg_1.(*ImS8)) + *(arg_2.(*ImS8))
		}
		if op == '-' {
			*(output.(*ImS8)) = *(arg_1.(*ImS8)) - *(arg_2.(*ImS8))
		}
		return
	case ImGuiDataType_U8:
		if op == '+' {
			*(output.(*ImU8)) = *(arg_1.(*ImU8)) + *(arg_2.(*ImU8))
		}
		if op == '-' {
			*(output.(*ImU8)) = *(arg_1.(*ImU8)) - *(arg_2.(*ImU8))
		}
		return
	case ImGuiDataType_S16:
		if op == '+' {
			*(output.(*ImS16)) = *(arg_1.(*ImS16)) + *(arg_2.(*ImS16))
		}
		if op == '-' {
			*(output.(*ImS16)) = *(arg_1.(*ImS16)) - *(arg_2.(*ImS16))
		}
		return
	case ImGuiDataType_U16:
		if op == '+' {
			*(output.(*ImU16)) = *(arg_1.(*ImU16)) + *(arg_2.(*ImU16))
		}
		if op == '-' {
			*(output.(*ImU16)) = *(arg_1.(*ImU16)) - *(arg_2.(*ImU16))
		}
		return
	case ImGuiDataType_S32:
		if op == '+' {
			*(output.(*ImS32)) = *(arg_1.(*ImS32)) + *(arg_2.(*ImS32))
		}
		if op == '-' {
			*(output.(*ImS32)) = *(arg_1.(*ImS32)) - *(arg_2.(*ImS32))
		}
		return
	case ImGuiDataType_U32:
		if op == '+' {
			*(output.(*ImU32)) = *(arg_1.(*ImU32)) + *(arg_2.(*ImU32))
		}
		if op == '-' {
			*(output.(*ImU32)) = *(arg_1.(*ImU32)) - *(arg_2.(*ImU32))
		}
		return
	case ImGuiDataType_S64:
		if op == '+' {
			*(output.(*ImS64)) = *(arg_1.(*ImS64)) + *(arg_2.(*ImS64))
		}
		if op == '-' {
			*(output.(*ImS64)) = *(arg_1.(*ImS64)) - *(arg_2.(*ImS64))
		}
		return
	case ImGuiDataType_U64:
		if op == '+' {
			*(output.(*ImU64)) = *(arg_1.(*ImU64)) + *(arg_2.(*ImU64))
		}
		if op == '-' {
			*(output.(*ImU64)) = *(arg_1.(*ImU64)) - *(arg_2.(*ImU64))
		}
		return
	case ImGuiDataType_Float:
		if op == '+' {
			*(output.(*float)) = *(arg_1.(*float)) + *(arg_2.(*float))
		}
		if op == '-' {
			*(output.(*float)) = *(arg_1.(*float)) - *(arg_2.(*float))
		}
		return
	case ImGuiDataType_Double:
		if op == '+' {
			*(output.(*double)) = *(arg_1.(*double)) + *(arg_2.(*double))
		}
		if op == '-' {
			*(output.(*double)) = *(arg_1.(*double)) - *(arg_2.(*double))
		}
		return
	case ImGuiDataType_COUNT:
		break
	}
	IM_ASSERT(false)
}

// DataTypeApplyOpFromText User can input math operators (e.g. +100) to edit a numerical values.
// NB: This is _not_ a full expression evaluator. We should probably add one and replace this dumb mess..
func DataTypeApplyOpFromText(buf string, initial_value_buf string, data_type ImGuiDataType, p_data any, format string) bool {
	for buf[0] == ' ' || buf[0] == '\t' {
		buf = buf[1:]
	}

	// We don't support '-' op because it would conflict with inputing negative value.
	// Instead you can use +-100 to subtract from an existing value
	var op = buf[0]
	if op == '+' || op == '*' || op == '/' {
		buf = buf[1:]
		for buf[0] == ' ' || buf[0] == '\t' {
			buf = buf[1:]
		}
	} else {
		op = 0
	}
	if len(buf) == 0 {
		return false
	}

	// Copy the value in an opaque buffer so we can compare at the end of the function if it changed at all.
	var type_info = DataTypeGetInfo(data_type)

	if format == "" {
		format = type_info.ScanFmt
	}

	// FIXME-LEGACY: The aim is to remove those operators and write a proper expression evaluator at some point..
	var arg1i int = 0
	if data_type == ImGuiDataType_S32 {
		var v = (p_data).(*int)
		var data_backup = *v
		var arg0i = *v
		var arg1f = 0.0
		if n, _ := fmt.Sscanf(initial_value_buf, format, &arg0i); op != 0 && n < 1 {
			return false
		}
		// Store operand in a float so we can use fractional value for multipliers (*1.1), but constant always parsed as integer so we can fit big integers (e.g. 2000000003) past float precision
		if op == '+' {
			if n, _ := fmt.Sscanf(buf, "%d", &arg1i); n != 0 {
				*v = arg0i + arg1i
			}
			// Add (use "+-" to subtract)
		} else if op == '*' {
			if n, _ := fmt.Sscanf(buf, "%f", &arg1f); n != 0 {
				*v = arg0i * int(arg1f)
			}
			// Multiply
		} else if op == '/' {
			if n, _ := fmt.Sscanf(buf, "%f", &arg1f); n != 0 && arg1f != 0.0 {
				*v = arg0i / int(arg1f)
			}
			// Divide
		} else {
			// Assign constant
			if n, _ := fmt.Sscanf(buf, format, &arg1i); n == 1 {
				*v = arg1i
			}
		}

		return data_backup != *v
	} else if data_type == ImGuiDataType_Float {
		// For floats we have to ignore format with precision (e.g. "%.2f") because sscanf doesn't take them in
		format = "%f"
		var v = p_data.(*float)
		var data_backup = *v
		var arg0f, arg1f float = *v, 0.0
		if n, _ := fmt.Sscanf(initial_value_buf, format, &arg0f); n < 1 && op != 0 {
			return false
		}
		if n, _ := fmt.Sscanf(buf, format, &arg1f); n < 1 {
			return false
		}
		if op == '+' {
			*v = arg0f + arg1f
			// Add (use "+-" to subtract)
		} else if op == '*' {
			*v = arg0f * arg1f
			// Multiply
		} else if op == '/' {
			if arg1f != 0.0 {
				*v = arg0f / arg1f
			}
			// Divide
		} else {
			*v = arg1f // Assign constant
		}

		return data_backup != *v
	} else if data_type == ImGuiDataType_Double {
		format = "%lf" // scanf differentiate float/double unlike printf which forces everything to double because of ellipsis
		var v = p_data.(*double)
		var data_backup = *v
		var arg0f, arg1f = *v, 0.0
		if n, _ := fmt.Sscanf(initial_value_buf, format, &arg0f); n < 1 && op != 0 {
			return false
		}
		if n, _ := fmt.Sscanf(buf, format, &arg1f); n < 1 {
			return false
		}
		if op == '+' {
			*v = arg0f + arg1f
			// Add (use "+-" to subtract)
		} else if op == '*' {
			*v = arg0f * arg1f
			// Multiply
		} else if op == '/' {
			if arg1f != 0.0 {
				*v = arg0f / arg1f
			}
			// Divide
		} else {
			*v = arg1f // Assign constant
		}
		return data_backup != *v
	} else if data_type == ImGuiDataType_U32 || data_type == ImGuiDataType_S64 || data_type == ImGuiDataType_U64 {
		var data_backup = reflect.ValueOf(p_data).Elem().Interface()

		// All other types assign constant
		// We don't bother handling support for legacy operators since they are a little too crappy. Instead we will later implement a proper expression evaluator in the future.
		if n, _ := fmt.Sscanf(buf, format, p_data); n < 1 {
			return false
		}

		return data_backup != reflect.ValueOf(p_data).Elem().Interface()
	} else {
		var data_backup = reflect.ValueOf(p_data).Elem().Interface()

		// Small types need a 32-bit buffer to receive the result from scanf()
		var v32 int
		if n, _ := fmt.Sscanf(buf, format, &v32); n < 1 {
			return false
		}
		if data_type == ImGuiDataType_S8 {
			*p_data.(*int8) = int8(ImClampInt(v32, -128, 127))
		} else if data_type == ImGuiDataType_U8 {
			*p_data.(*uint8) = uint8(ImClampInt(v32, 0, 255))
		} else if data_type == ImGuiDataType_S16 {
			*p_data.(*int16) = int16(ImClampInt(v32, -32768, 32767))
		} else if data_type == ImGuiDataType_U16 {
			*p_data.(*uint16) = uint16(ImClampInt(v32, 0, 65535))
		} else {
			IM_ASSERT(false)
		}

		return data_backup != reflect.ValueOf(p_data).Elem().Interface()
	}
}

func DataTypeCompare(data_type ImGuiDataType, arg_1 any, arg_2 any) int {
	switch data_type {
	case ImGuiDataType_S8:
		return int(arg_1.(int8) - arg_2.(int8))
	case ImGuiDataType_U8:
		return int(arg_1.(uint8) - arg_2.(uint8))
	case ImGuiDataType_S16:
		return int(arg_1.(int16) - arg_2.(int16))
	case ImGuiDataType_U16:
		return int(arg_1.(uint16) - arg_2.(uint16))
	case ImGuiDataType_S32:
		return arg_1.(int32) - arg_2.(int32)
	case ImGuiDataType_U32:
		return int(arg_1.(uint32) - arg_2.(uint32))
	case ImGuiDataType_S64:
		return int(arg_1.(int64) - arg_2.(int64))
	case ImGuiDataType_U64:
		return int(arg_1.(uint64) - arg_2.(uint64))
	case ImGuiDataType_Float:
		return int(arg_1.(float32) - arg_2.(float32))
	case ImGuiDataType_Double:
		return int(arg_1.(float64) - arg_2.(float64))
	}
	IM_ASSERT(false)
	return 0
}

func DataTypeClamp(data_type ImGuiDataType, a any, n any, x any) bool {
	switch data_type {
	case ImGuiDataType_S8:
		*a.(*int8) = int8(ImClampInt(int(a.(int8)), int(n.(int8)), int(x.(int8))))
	case ImGuiDataType_U8:
		*a.(*uint8) = uint8(ImClampInt(int(a.(uint8)), int(n.(uint8)), int(x.(uint8))))
	case ImGuiDataType_S16:
		*a.(*int16) = int16(ImClampInt(int(a.(int16)), int(n.(int16)), int(x.(int16))))
	case ImGuiDataType_U16:
		*a.(*uint16) = uint16(ImClampInt(int(a.(uint16)), int(n.(uint16)), int(x.(uint16))))
	case ImGuiDataType_S32:
		*a.(*int32) = ImClampInt(a.(int32), n.(int32), x.(int32))
	case ImGuiDataType_U32:
		*a.(*uint32) = uint32(ImClampInt(int(a.(uint32)), int(n.(uint32)), int(x.(uint32))))
	case ImGuiDataType_S64:
		*a.(*int64) = ImClampInt64(a.(int64), n.(int64), x.(int64))
	case ImGuiDataType_U64:
		*a.(*uint64) = ImClampUint64(a.(uint64), n.(uint64), x.(uint64))
	case ImGuiDataType_Float:
		*a.(*float32) = ImClamp(a.(float32), n.(float32), x.(float32))
	case ImGuiDataType_Double:
		*a.(*float64) = ImClamp64(a.(float64), n.(float64), x.(float64))
	}
	IM_ASSERT(false)
	return false
}

func DataTypeGetInfo(data_type ImGuiDataType) *ImGuiDataTypeInfo {
	IM_ASSERT(data_type >= 0 && data_type < ImGuiDataType_COUNT)
	return &GDataTypeInfo[data_type]
}
