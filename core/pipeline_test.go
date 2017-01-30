package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mockPipe struct {
	itemsProcessed int
}

func (tp *mockPipe) Process(in <-chan Entry) <-chan Entry {
	out := make(chan Entry)
	go func() {
		for entry := range in {
			tp.itemsProcessed++
			out <- entry
		}
		close(out)
	}()
	return out
}

func TestPipelineNoEntriesLost(t *testing.T) {
	var pipes []Pipe
	pipes = append(pipes, &mockPipe{})
	pipes = append(pipes, &mockPipe{})
	p := NewPipeline(pipes...)

	var entries []Entry
	entries = append(entries, Entry{"testkey": "testvalue"}, Entry{"testkey2": "testvalue2"})

	for _, entry := range entries {
		p.In <- entry
	}
	close(p.In)

	for _ = range p.Out {

	}

	for _, pipe := range pipes {
		require.EqualValues(t, len(entries), pipe.(*mockPipe).itemsProcessed)
	}
}

// Tests if the pipeline works correctly if no pipe is passed at all (hence,
// items added to input can be read from output directly)
func TestPipelineNoPipes(t *testing.T) {
	p := NewPipeline()

	entry := Entry{}
	go func() {
		p.In <- entry
		close(p.In)
	}()

	recvEntry := <-p.Out
	require.Equal(t, entry, recvEntry)
}
