package filesystem

import "regexp"

func isValidChar(input string) bool {
	validCharPattern := "^[a-zA-Z0-9 ]*$"
	matched, _ := regexp.MatchString(validCharPattern, input)
	return matched
}
