package main

import "testing"

//Testing

// TestMRISAConnection makes sure that a connection to MRISA can be made
func TestMRISAConnection(t *testing.T) {
	testImg := "https://upload.wikimedia.org/wikipedia/commons/thumb/a/af/Tux.png/220px-Tux.png"

	valOne, valTwo := imageSearch(testImg)

	if valOne == nil && valTwo == nil {
		t.Error("Was expecting two string slices, got nil for both.")
	}

}

//Benchmarks

// BenchmarkMRISAConnection tests how long a connection takes to the MRISA server
func BenchmarkMRISAConnection(b *testing.B) {
	testImg := "https://raw.githubusercontent.com/chand1012/Discord-Quick-Meme/master/canary-icon.png"

	for i := 0; i < b.N; i++ {
		_, _ = imageSearch(testImg)
	}
}
