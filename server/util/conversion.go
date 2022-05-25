package util

import "strconv"

// bool relation
func RequireBool(value string) bool {
	result, _ := strconv.ParseBool(value)
	return result
}

// int relation
func Int64(value string) int64 {
	i64, _ := strconv.ParseInt(value, 10, 64)
	return i64
}

func Uint64(value string) (uint64, error) {
	ui64, err := strconv.ParseUint(value, 10, 64)
	return ui64, err
}

func Uint32(value string) uint32 {
	ui64, _ := strconv.ParseUint(value, 10, 32)
	return uint32(ui64)
}

func Uint16(value string) uint16 {
	ui64, _ := strconv.ParseUint(value, 10, 16)
	return uint16(ui64)
}
