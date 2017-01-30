package filter

import (
	"datapipeline/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

var duplicateTests = []struct {
	input    []core.Entry
	expected []core.Entry
}{
	{ // Duplicates must be removed
		input:    []core.Entry{core.Entry{"name": "bar"}, core.Entry{"name": "bar"}},
		expected: []core.Entry{core.Entry{"name": "bar"}},
	},
	{ // First entry with the conflicting field "name" has priority
		input:    []core.Entry{core.Entry{"name": "bar", "testfield": "testvalue"}, core.Entry{"name": "bar"}},
		expected: []core.Entry{core.Entry{"name": "bar", "testfield": "testvalue"}},
	},
	{ // Preserve order of items in pipeline
		input:    []core.Entry{core.Entry{"name": "foo"}, core.Entry{"name": "bar"}, core.Entry{"name": "baz"}},
		expected: []core.Entry{core.Entry{"name": "foo"}, core.Entry{"name": "bar"}, core.Entry{"name": "baz"}},
	},
}

func TestDeduplicate(t *testing.T) {
	for _, tt := range duplicateTests {
		filter := NewDuplicate("name")
		in := make(chan core.Entry)
		out := filter.Process(in)

		go func() {
			for _, entry := range tt.input {
				in <- entry
			}
			close(in)
		}()

		var entries []core.Entry
		for entry := range out {
			entries = append(entries, entry)
		}
		assert.EqualValues(t, tt.expected, entries)
	}
}

var duplicateTestsNonExistingField = []struct {
	input    []core.Entry
	expected []core.Entry
}{
	{
		input:    []core.Entry{core.Entry{"name": "bar"}, core.Entry{"name": "bar"}},
		expected: []core.Entry{core.Entry{"name": "bar"}, core.Entry{"name": "bar"}},
	},
	{
		input:    []core.Entry{core.Entry{"name": "bar", "testfield": "testvalue"}, core.Entry{"name": "bar"}},
		expected: []core.Entry{core.Entry{"name": "bar", "testfield": "testvalue"}, core.Entry{"name": "bar"}},
	},
	{
		input:    []core.Entry{core.Entry{"name": "foo"}, core.Entry{"name": "bar"}, core.Entry{"name": "baz"}},
		expected: []core.Entry{core.Entry{"name": "foo"}, core.Entry{"name": "bar"}, core.Entry{"name": "baz"}},
	},
}

func TestDedupNonExistingField(t *testing.T) {
	for _, tt := range duplicateTestsNonExistingField {
		filter := NewDuplicate("nonexistingfield")
		in := make(chan core.Entry)
		out := filter.Process(in)

		go func() {
			for _, entry := range tt.input {
				in <- entry
			}
			close(in)
		}()

		var entries []core.Entry
		for entry := range out {
			entries = append(entries, entry)
		}
		assert.EqualValues(t, tt.expected, entries)
	}

}
