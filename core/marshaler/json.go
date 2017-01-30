package marshaler

import (
	"datapipeline/core"
	"encoding/json"
	"io"
)

// jsonMarshaler marshals all entries into one json array.
type jsonMarshaler struct {
	writer io.Writer
}

// NewJSON constructs a new jsonMarshaler.
func NewJSON(writer io.Writer) core.Marshaler {
	return &jsonMarshaler{writer: writer}
}

// Marshal implements the core.Marshaler interface.
func (m *jsonMarshaler) Marshal(in <-chan core.Entry) error {
	var rows []map[string]string
	for entry := range in {
		rows = append(rows, entry)
	}

	encoder := json.NewEncoder(m.writer)
	return encoder.Encode(rows) // TODO
}
