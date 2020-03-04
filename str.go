// For all functions that manipulate strings
package main

import (
	"regexp"
	"strconv"
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
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	errCheck("Error compiling regexp", err, false)
	if err != nil {
		return "" // return empty string because more errors would occur otherwise
	}
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
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	errCheck("Error compiling regexp", err, false)
	if err != nil {
		return nil
	}
	var returnSlice []string
	for _, thing := range input {
		output := reg.ReplaceAllString(thing, "")
		returnSlice = append(returnSlice, output)
	}
	return returnSlice
}

// this if for "how long ago"
// ex: "5 hours ago"
// ex: "4 days ago"
func timeStrToSeconds(stamp string) int64 {
	var stampLength string
	var finalTime int64
	lengths := []string{"second", "minute", "hour", "day", "year"}
	for _, l := range lengths {
		if strings.Contains(stamp, l) {
			stampLength = l
			break
		}
	}
	spaceIndex := strings.Index(stamp, " ")
	timeUnknown, _ := strconv.ParseInt(stamp[:spaceIndex], 10, 64)
	switch stampLength {
	case "second":
		finalTime = timeUnknown

	case "minute":
		finalTime = timeUnknown * 60

	case "hour":
		finalTime = timeUnknown * 3600

	case "day":
		finalTime = timeUnknown * 3600 * 24

	case "year":
		finalTime = timeUnknown * 3600 * 8760
	}
	return finalTime
}
