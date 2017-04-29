package marshaler

import "github.com/birdayz/datapipeline/core"

// NopMarshaler drops all entries read from the input channel
type nopMarshaler struct{}

// NewNop constructs a new nopMarshaler.
func NewNop() core.Marshaler {
	return &nopMarshaler{}
}

// Marshal implements the core.Marshaler interface
func (m *nopMarshaler) Marshal(in <-chan core.Entry) error {
	for _ = range in {
	}
	return nil
}
