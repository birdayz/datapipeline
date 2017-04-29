package filter

import (
	"github.com/birdayz/datapipeline/core"
	"testing"

	"github.com/stretchr/testify/require"
)

type pipeTest []struct {
	input        []core.Entry
	expectedAsc  []core.Entry
	expectedDesc []core.Entry
}

var sortNameTests = pipeTest{
	{
		input:        []core.Entry{core.Entry{"name": "xy"}, core.Entry{"name": "f"}, core.Entry{"name": "ab"}},
		expectedAsc:  []core.Entry{core.Entry{"name": "ab"}, core.Entry{"name": "f"}, core.Entry{"name": "xy"}},
		expectedDesc: []core.Entry{core.Entry{"name": "xy"}, core.Entry{"name": "f"}, core.Entry{"name": "ab"}},
	},
	{
		input:        []core.Entry{core.Entry{"name": "1"}, core.Entry{"name": "a"}},
		expectedAsc:  []core.Entry{core.Entry{"name": "1"}, core.Entry{"name": "a"}},
		expectedDesc: []core.Entry{core.Entry{"name": "a"}, core.Entry{"name": "1"}},
	},
	{
		input:        []core.Entry{core.Entry{"name": ""}, core.Entry{"name": "a"}},
		expectedAsc:  []core.Entry{core.Entry{"name": ""}, core.Entry{"name": "a"}},
		expectedDesc: []core.Entry{core.Entry{"name": "a"}, core.Entry{"name": ""}},
	},
}

func TestSortByNameAsc(t *testing.T) {
	for _, test := range sortNameTests {
		filter := NewSort("name", OrderAscending)
		in := make(chan core.Entry)
		out := filter.Process(in)

		for _, entry := range test.input {
			in <- entry
		}
		close(in)

		var entries []core.Entry
		for entry := range out {
			entries = append(entries, entry)
		}
		require.EqualValues(t, test.expectedAsc, entries)
	}
}

func TestSortByNameDesc(t *testing.T) {
	for _, test := range sortNameTests {
		filter := NewSort("name", OrderDescending)
		in := make(chan core.Entry)
		out := filter.Process(in)

		for _, entry := range test.input {
			in <- entry
		}
		close(in)

		var entries []core.Entry
		for entry := range out {
			entries = append(entries, entry)
		}
		require.EqualValues(t, test.expectedDesc, entries)
	}
}
