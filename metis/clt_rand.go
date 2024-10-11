package metis

import (
	"errors"
	"math/rand"
)

var (
	ErrOverflow     = errors.New("value is too large")
	ErrInvalidRange = errors.New("invalid range")
)

const (
	maxInt63 = int64(1<<62 - 1)
)

// CLTGenerator generates random integers that follow an approximate normal distribution
type CLTGenerator struct{}

// NewCLTGenerator initializes a CLT generator without state dependency
func NewCLTGenerator() *CLTGenerator {
	return &CLTGenerator{}
}

// GenerateRandomInt generates a random integer from min to max (inclusive), based on a provided seed.
// This simulates a normal distribution by summing multiple uniform random variables.
func (g *CLTGenerator) GenerateRandomInt(seed int64, min, max uint64) (int64, error) {
	// Check for invalid range
	if min > max {
		return 0, ErrInvalidRange
	}

	// Check for overflow
	rangeDiff := max - min
	if max > uint64(maxInt63) || min > uint64(maxInt63) || rangeDiff > uint64(maxInt63) {
		return 0, ErrOverflow
	}

	// Number of samples to sum (higher means closer to normal distribution)
	// default to 6 because this is a balanced value between performance and accuracy
	numSamples := int64(6)
	maxSampleSize := maxInt63 / numSamples

	if rangeDiff > uint64(maxSampleSize) {
		// if range is too large change sample size to fit within maxInt63
		// this trades accuracy with safety
		numSamples = maxInt63 / int64(rangeDiff)
		maxSampleSize = maxInt63 / numSamples
	}

	// Create a local new source to avoid disruption from the global random status
	r := rand.New(rand.NewSource(seed))
	sum, minInt, maxInt := int64(0), int64(min), int64(max)
	for i := int64(0); i < numSamples; i++ {
		// rand.Int63n returns a value in the range [0, n), so we need to add 2 here to include max
		randValue := r.Int63n(int64(rangeDiff) + 2)
		sum += randValue
	}

	// Average the sum to bring it back within the desired range
	result := (sum / numSamples) + minInt
	if result < minInt {
		return minInt, nil
	} else if result > maxInt {
		return maxInt, nil
	}
	return result, nil
}
