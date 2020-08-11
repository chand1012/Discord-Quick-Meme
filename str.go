// For all functions that manipulate strings
package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

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
		fmt.Println("Error compiling regexp:", err)
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
		fmt.Println("Error compiling regexp:", err)
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

// I am sure there is a more efficient way to do this
// but this works
func interfaceToStringSlice(inter interface{}) []string {
	var stringSlice []string
	interString := fmt.Sprintf("%v", inter)
	re, _ := regexp.Compile("[^0-9a-zA-Z ]")
	interString = re.ReplaceAllString(interString, "")
	stringSlice = strings.Split(interString, " ")
	return stringSlice
}

func matchRegexList(expressions []string, testStr string) bool {
	for _, item := range expressions {
		compiled, err := regexp.Compile(item)
		if err != nil {
			fmt.Println("Error compiling regexp", item)
			fmt.Println(err.Error())
			fmt.Println("Skipping.")
			continue
		}
		if compiled.MatchString(testStr) {
			return true
		}
	}
	return false
}

// from https://stackoverflow.com/a/10030772/5178731
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
// This is for TESTING ONLY
func randString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
