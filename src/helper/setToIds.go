package helper

import "strings"

func SetToIds(set string) (res []int) {
	set = strings.ReplaceAll(set, "#", "")
	s := strings.Split(set, ",")
	for _, r := range s {
		res = append(res, StrToInt(r))
	}
	return
}
