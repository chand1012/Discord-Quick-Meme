// For all functions that manipulate strings
package main

import (
	"log"
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
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		log.Println("Error compiling regexp:", err)
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
	if err != nil {
		log.Println("Error compiling regexp:", err)
		return nil
	}
	var returnSlice []string
	for _, thing := range input {
		output := reg.ReplaceAllString(thing, "")
		returnSlice = append(returnSlice, output)
	}
	return returnSlice
}

func matchRegexList(expressions []string, testStr string) bool {
	for _, item := range expressions {
		compiled, err := regexp.Compile(item)
		if err != nil {
			log.Println("Error compiling regexp", item)
			log.Println(err.Error())
			log.Println("Skipping.")
			continue
		}
		if compiled.MatchString(testStr) {
			return true
		}
	}
	return false
}
