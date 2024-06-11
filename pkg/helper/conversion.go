package helper

import "strconv"

func ParseStringToUint(input string) uint {
	result, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0
	}

	return uint(result)
}
