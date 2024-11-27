package metis

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"sync"
	"testing"
	"time"
)

// TestMultipleGoroutinesConsistency runs the generator in multiple goroutines simulating different machines
func TestMultipleGoroutinesConsistency(t *testing.T) {
	seedBase := time.Now().UnixNano()
	numGoroutines := 10
	iterations := 1000
	min := uint64(0)
	max := uint64(100)

	var wg sync.WaitGroup
	results := make([][]int64, numGoroutines)

	mu := sync.Mutex{}
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			gen := NewCLTGenerator()
			localResult := make([]int64, iterations)
			for j := 0; j < iterations; j++ {
				seed := seedBase + int64(j)
				value, err := gen.GenerateRandomInt(seed, min, max)
				if err != nil {
					t.Errorf("Error generating random int: %v", err)
				}
				localResult[j] = value
			}

			mu.Lock()
			results[index] = localResult
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	for i := 1; i < numGoroutines; i++ {
		if !compareResults(results[0], results[i]) {
			t.Errorf("Inconsistent results between goroutine 0 and goroutine %d", i)
			return
		}
	}
}

// compareResults checks if two slices have the same values
func compareResults(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestGenerateRandomDistribution tests the generator and outputs the frequency distribution to CSV
func TestGenerateRandomDistribution(t *testing.T) {
	seedBase := time.Now().UnixNano()
	testCases := []struct {
		iterations int
		filename   string
		min        uint64
		max        uint64
	}{
		{500, "test_500_iterations.csv", 0, 10},
		{5000, "test_5000_iterations.csv", 0, 100},
		{50000, "test_50000_iterations.csv", 0, 5000},
		{100000, "test_100000_iterations.csv", 0, 10000},
		// Additional test case for exceeding maxSampleSize to verify normal distribution
		{1000000, "test_1000000_max_sample_size_iterations.csv", 0, uint64(maxInt63) / 6},
		{1000000, "test_1000000_double_max_sample_size_iterations.csv", 0, uint64(maxInt63) / 2},
		{1000000, "test_1000000_max_limit_iterations.csv", 0, uint64(maxInt63)},
	}

	gen := NewCLTGenerator()

	for _, tc := range testCases {
		counts := make(map[int64]int)
		for i := 0; i < tc.iterations; i++ {
			seed := seedBase + int64(i)
			value, err := gen.GenerateRandomInt(seed, tc.min, tc.max)
			if err != nil {
				t.Errorf("Error generating random int: %v", err)
			}
			counts[value]++
		}

		// Save results to CSV file
		file, err := os.Create(tc.filename)
		if err != nil {
			t.Errorf("Error creating file %s: %v", tc.filename, err)
			continue
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"Number", "Occurrences"})

		type pair struct {
			num   int64
			count int
		}

		pairs := make([]pair, 0, len(counts))

		for num, count := range counts {
			pairs = append(pairs, pair{num, count})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].num < pairs[j].num
		})

		for _, p := range pairs {
			writer.Write([]string{fmt.Sprintf("%d", p.num), fmt.Sprintf("%d", p.count)})
		}
	}
}

// TestGenerateRandomIntEdgeCases tests edge cases for the GenerateRandomInt function
func TestGenerateRandomIntEdgeCases(t *testing.T) {
	gen := NewCLTGenerator()
	seed := time.Now().UnixNano()

	testCases := []struct {
		min          uint64
		max          uint64
		expectsError bool
	}{
		// Edge cases
		{0, 0, false},                                      // Min equals Max
		{0, uint64(maxInt63), false},                       // Max is max int63
		{0, uint64(maxInt63) + 1, true},                    // Max exceeds int63
		{uint64(maxInt63) + 1, uint64(maxInt63) + 2, true}, // Both Min and Max exceed int63
		{100, 50, true},                                    // Min greater than Max (invalid range)
		{uint64(maxInt63) - 10, uint64(maxInt63), false},   // Near max int63, should work
		{1, uint64(maxInt63) / 2, false},                   // Large range but within int63, should work
		{uint64(maxInt63) / 2, uint64(maxInt63), false},    // Upper half of int63 range
		// Small range with large min and max values
		{uint64(maxInt63) - 5, uint64(maxInt63), false}, // Small range with large values
		{uint64(maxInt63) - 1, uint64(maxInt63), false}, // Min close to max
	}

	for i, tc := range testCases {
		_, err := gen.GenerateRandomInt(seed, tc.min, tc.max)
		if (err != nil) != tc.expectsError {
			t.Errorf("Unexpected error for case %d: min: %d, max: %d, got error: %v", i, tc.min, tc.max, err)
		}
	}
}

// BenchmarkGenerateRandomInt benchmarks the efficiency of the random generation function
func BenchmarkGenerateRandomInt(b *testing.B) {
	gen := NewCLTGenerator()
	seed := time.Now().UnixNano()
	min := uint64(0)
	max := uint64(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.GenerateRandomInt(seed+int64(i), min, max)
		if err != nil {
			b.Errorf("Error generating random int: %v", err)
		}
	}
}
