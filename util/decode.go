package util

import (
	"fmt"
	"strconv"
	"unicode"
)

// DecodeSize decodes a given size in string format. Eg:- 10MB
// And returns the actual size calculated in Bytes
func DecodeSize(s string) int64 {
	digits := make([]rune, 0)
	chars := make([]rune, 0)

	for _, c := range s {
		if unicode.IsDigit(c) {
			digits = append(digits, c)
			continue
		}
		if unicode.IsLetter(c) {
			chars = append(chars, c)
			continue
		}
	}

	size := fmt.Sprintf("%s", string(digits))
	sizeIdentifier := fmt.Sprintf("%s", string(chars))

	sizeconv, _ := strconv.Atoi(size)

	sizemap := map[string]int64{
		"B":  int64(sizeconv),
		"KB": int64(1024 * sizeconv),
		"MB": int64(1024 * 1024 * sizeconv),
		"GB": int64(1024 * 1024 * 1024 * sizeconv),
	}

	actualsize, exists := sizemap[sizeIdentifier]
	if !exists {
		return -1 // in case of unsupported size identifiers return -1
	}

	return actualsize
}
