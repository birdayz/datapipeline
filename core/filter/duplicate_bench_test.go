package filter

import (
	"datapipeline/core"
	"testing"
)

// Benchmark the performance of the DuplicateFilter with an arbitrary number of
// Entries. A list with ten items and 20% duplicates is used, and concatenated
// with itself to get a larger list. Note that writing/reading from the pipe
// channels is included in the benchtime as well and might impact the results,
// hence not only the sorting operation itself is benchmarked.
func benchmarkRemoveDuplicates(b *testing.B, sizeMultipleOfTen int) {
	var tenEntries = []core.Entry{
		core.Entry{"name": "duplicate"},
		core.Entry{"name": "1"},
		core.Entry{"name": "2"},
		core.Entry{"name": "3"},
		core.Entry{"name": "4"},
		core.Entry{"name": "5"},
		core.Entry{"name": "6"},
		core.Entry{"name": "7"},
		core.Entry{"name": "8"},
		core.Entry{"name": "duplicate"},
	}
	filter := NewDuplicate("name")
	// Generate a larger list of entries
	allEntries := make([]core.Entry, 0, len(tenEntries)*sizeMultipleOfTen)
	for i := 0; i < sizeMultipleOfTen; i++ {
		allEntries = append(allEntries, tenEntries...)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := make(chan core.Entry)
		go func() {
			for _, entry := range allEntries {
				in <- entry
			}
			close(in)

		}()
		out := filter.Process(in)
		for _ = range out {
		}
	}
	b.StopTimer()
}

func BenchmarkRemoveDuplicates_1k(b *testing.B) {
	benchmarkRemoveDuplicates(b, 100)
}

func BenchmarkRemoveDuplicates_10k(b *testing.B) {
	benchmarkRemoveDuplicates(b, 1000)
}
func BenchmarkRemoveDuplicates_100k(b *testing.B) {
	benchmarkRemoveDuplicates(b, 10000)
}

func BenchmarkRemoveDuplicates_1M(b *testing.B) {
	benchmarkRemoveDuplicates(b, 100000)
}
