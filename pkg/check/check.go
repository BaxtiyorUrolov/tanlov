package check

import (

	"unicode"
)

func PhoneNumber(phone string) bool {
	for _, r := range phone {
		if r == '+' {
			continue
		} else if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}


