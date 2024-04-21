package filesystem

import "regexp"

func isValidInput(input string) bool {
	if len(input) > 20 {
		return false
	}
	validCharPattern := "^[a-zA-Z0-9 ]*$"
	matched, _ := regexp.MatchString(validCharPattern, input)
	return matched
}
