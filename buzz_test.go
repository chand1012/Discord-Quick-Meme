package main

import "testing"

func TestBuzzWords(t *testing.T) {
	test := getABuzzWord()

	if test == "" {
		t.Error("Could not get a buzzword.")
	}
}
