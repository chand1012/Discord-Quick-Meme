package main

import (
	"testing"

	"github.com/joho/godotenv"
)

// Testing
// Nothing here yet

// Benchmarks

// BenchmarkRedditConnection testing reddit speed
func BenchmarkRedditConnection(b *testing.B) {
	godotenv.Load()
	for i := 0; i < b.N; i++ {
		_ = PingReddit()
	}
}
