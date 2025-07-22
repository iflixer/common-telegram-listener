package helper

import "strconv"

func StrToInt(s string) (res int) {
	res, _ = strconv.Atoi(s)
	return
}

func StrToInt64(s string) (res int64) {
	res, _ = strconv.ParseInt(s, 10, 0)
	return
}
