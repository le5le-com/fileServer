package utils

import (
	"strconv"
	"unsafe"
)

func Int(v interface{}) int {
	switch reply := v.(type) {
	case int:
		return reply
	case int8:
		return int(reply)
	case int16:
		return int(reply)
	case int32:
		return int(reply)
	case int64:
		return int(reply)
	case []byte:
		n, _ := strconv.ParseInt(string(reply), 10, 0)
		return int(n)
	case nil:
		return 0
	case float64:
		return int(reply)
	}

	return 0
}

func String(v interface{}) string {
	switch reply := v.(type) {
	case string:
		return reply
	case int:
		return strconv.Itoa(reply)
	case int8:
		return strconv.Itoa(int(reply))
	case int16:
		return strconv.Itoa(int(reply))
	case int32:
		return strconv.Itoa(int(reply))
	case int64:
		return strconv.Itoa(int(reply))
	case []byte:
		return *(*string)(unsafe.Pointer(&reply))
	case nil:
		return ""
	case bool:
		if reply {
			return "true"
		}
		return "false"
	}

	return ""
}
