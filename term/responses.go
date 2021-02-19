package term

import (
	"fmt"
	"strconv"
	"strings"
)

func UserDeclined(response string) bool {
	return strings.ToLower(response)[0] == 'n'
}

func UserAccepted(response string) bool {
	return strings.ToLower(response)[0] != 'n'
}

func ValidateNumericChoice(min int, max int) func(string) string {
	return func(userInput string) string {
		index, err := strconv.Atoi(userInput)
		if err != nil || index < min || index > max {
			return fmt.Sprintf("Sorry, I didn't understand that. Give me a number between %v and %v", min, max)
		}
		return ""
	}
}
