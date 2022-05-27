package controller_utils

func MemoryToByte(memString string) int {
	if memString == `0` || memString == `` {
		return 0
	}
	memByte := 0
	for _, c := range memString {
		if c >= '0' && c <= '9' {
			memByte = memByte*10 + int(c-'0')
		} else if c == 'K' || c == 'k' {
			return memByte * 1024
		} else if c == 'M' || c == 'm' {
			return memByte * 1024 * 1024
		} else if c == 'G' || c == 'g' {
			return memByte * 1024 * 1024 * 1024
		}
	}
	return 0
}
