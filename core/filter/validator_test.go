package filter

import (
	"github.com/birdayz/datapipeline/core"
	"testing"

	"github.com/stretchr/testify/require"
)

type countingMockValidator struct {
	entriesProcessed int
}

func (v *countingMockValidator) Validate(entry core.Entry) bool {
	v.entriesProcessed++
	return true
}

type inValidator struct {
}

func (iv *inValidator) Validate(entry core.Entry) bool {
	return false
}

// TestValidatingFilter checks if all validators of the filter are evaluated
func TestAllValidatorsUsed(t *testing.T) {
	val1 := &countingMockValidator{}
	val2 := &countingMockValidator{}
	val3 := &countingMockValidator{}
	filter := NewValidator(val1, val2, val3)
	in := make(chan core.Entry)
	out := filter.Process(in)

	in <- core.Entry{}
	close(in)

	for _ = range out {
	}
	require.EqualValues(t, 1, val1.entriesProcessed)
	require.EqualValues(t, 1, val2.entriesProcessed)
	require.EqualValues(t, 1, val3.entriesProcessed)
}

// Test if a fail in the first validator skips the remaining validators
func TestValidatorFailStop(t *testing.T) {
	val1 := &inValidator{}
	val2 := &countingMockValidator{}
	val3 := &countingMockValidator{}
	filter := NewValidator(val1, val2, val3)
	in := make(chan core.Entry)
	out := filter.Process(in)

	in <- core.Entry{}
	close(in)

	itemsReceived := 0
	for _ = range out {
		itemsReceived++
	}
	require.EqualValues(t, 0, itemsReceived)
	require.EqualValues(t, 0, val2.entriesProcessed)
	require.EqualValues(t, 0, val3.entriesProcessed)

}
