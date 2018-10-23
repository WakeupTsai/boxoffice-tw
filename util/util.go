package util

import (
	"strconv"
	"strings"
)

func StringToInt(str string) int {
	result, _ := strconv.Atoi(strings.Replace(str, ",", "", -1))
	return result
}
