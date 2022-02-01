package validator

import "strings"

func checkStringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
