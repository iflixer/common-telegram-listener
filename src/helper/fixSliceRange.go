package helper

func FixSliceRange(start, end, length int) (int, int) {
	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}
	if end > length {
		end = length
	}
	return start, end
}
