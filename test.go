package main

import "testing"

// Testing
// Nothing here yet

// Benchmarks

// BenchmarkRedditConnection testing reddit speed
func BenchmarkRedditConnection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PingReddit()
	}
}
