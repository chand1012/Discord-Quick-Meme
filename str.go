// For all functions that manipulate strings
package main

import (
	"regexp"
	"strings"
)

func stringInSlice(s string, a []string) bool {
	for _, thing := range a {
		if thing == s {
			return true
		}
	}
	return false
}

func textFilter(input string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	outputString := reg.ReplaceAllString(input, "")
	return outputString
}

// ContainsAnySubstring Checks if any of the substrings in the array are in the test string
func ContainsAnySubstring(testString string, strArray []string) bool {
	for _, str := range strArray {
		if strings.Contains(testString, str) {
			return true
		}
	}
	return false
}

func textFilterSlice(input []string) []string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	var returnSlice []string
	for _, thing := range input {
		output := reg.ReplaceAllString(thing, "")
		returnSlice = append(returnSlice, output)
	}
	return returnSlice
}
