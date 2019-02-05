package utils

import (
	"log"
	"strings"
)

// AssertStringContains compares two strings and fails if it isn't contained
func AssertStringContains(mainString string, subString string, errorMessage string) {
	if !strings.Contains(mainString, subString) {
		log.Fatal("Error: " + errorMessage + "\nExpected <<" + mainString + ">> but was <<" + subString + ">>")
	}
}
