package main

import "testing"

//Tests

//TestJSONExtract tests if the subs file can be extracted from
func TestJSONExtract(t *testing.T) {
	subMap := SubExtract("subs.json")

	for k := range subMap {
		_, ok := subMap[k]
		if !ok {
			t.Error("There was an error extracting all of the subs from the map.")
		}
	}
}

// Benchmarks

// BenchmarkJSONExtract test how long it takes to get the subreddits from the file
func BenchmarkJSONExtract(b *testing.B) {
	filename := "subs.json"

	for i := 0; i < b.N; i++ {
		_ = SubExtract(filename)
	}
}
