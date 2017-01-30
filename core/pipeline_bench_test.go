package core

import "testing"

func BenchmarkInOut(b *testing.B) {
	entry := Entry{"test": "bla", "lol": "x"}
	pipeline := NewPipeline()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			pipeline.In <- entry
		}
		close(pipeline.In)
	}()

	for _ = range pipeline.Out {
	}
}
func BenchmarkInOut_2(b *testing.B) {
	entry := Entry{"test": "bla", "lol": "x"}
	pipeline := NewPipeline(NewNop())
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			pipeline.In <- entry
		}
		close(pipeline.In)
	}()

	for _ = range pipeline.Out {
	}
}

func BenchmarkInOut_4(b *testing.B) {
	entry := Entry{"test": "bla", "lol": "x"}
	pipeline := NewPipeline(NewNop(), NewNop(), NewNop())
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			pipeline.In <- entry
		}
		close(pipeline.In)
	}()

	for _ = range pipeline.Out {
	}
}

// Benchmark the instantiation cost of a single Entry
func BenchmarkCreateEntry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Entry{}
	}
}
