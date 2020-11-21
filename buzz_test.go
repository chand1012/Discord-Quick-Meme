package main

import "testing"

func TestBuzzWords(t *testing.T) {
	test := getABuzzWord()

	if len(test) == 0 {
		t.Error("Could not get a buzzword.")
	}
}
