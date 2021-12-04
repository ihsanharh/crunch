package utils

import "fmt"

func Fmt(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}

func Is(value bool, a, n string) string {
	if value {
		return a
	} else {
		return n
	}
}