package unmarshaler

import (
	"bytes"
	"datapipeline/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	expected := []core.Entry{core.Entry{"col1": "val1", "col2": "val2", "col3": "val3"}, core.Entry{"col1": "val21", "col2": "val22", "col3": "val23"}}
	testInput := []byte("col1,col2,col3\nval1,val2,val3\nval21,val22,val23")

	u := NewCSV(bytes.NewReader(testInput))

	out := make(chan core.Entry)
	go func() {
		err := u.Unmarshal(out)
		require.NoError(t, err)
	}()
	var entries []core.Entry
	for entry := range out {
		entries = append(entries, entry)
	}
	require.Equal(t, expected, entries)
}

func TestUnmarshalGarbage(t *testing.T) {
	f := bytes.NewReader([]byte{0xFF, 0x00, 0xFE})
	u := NewCSV(f)

	out := make(chan core.Entry)
	go func() {
		err := u.Unmarshal(out)
		require.NoError(t, err)
	}()
	var entries []core.Entry
	for entry := range out {
		entries = append(entries, entry)
	}
	require.EqualValues(t, 0, len(entries))
}
