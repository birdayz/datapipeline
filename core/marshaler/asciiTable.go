package marshaler

import (
	"github.com/birdayz/datapipeline/core"
	"io"

	"github.com/olekukonko/tablewriter"
)

// asciiTableMarshaler marshals all entries into a fancy ascii table.
type asciiTableMarshaler struct {
	fields []string
	writer io.Writer
}

// NewASCIITable constructs a new asciiTableMarshaler.
func NewASCIITable(fields []string, writer io.Writer) core.Marshaler {
	return &asciiTableMarshaler{fields: fields, writer: writer}
}

// Marshal implements the core.Marshal interface.
func (s *asciiTableMarshaler) Marshal(in <-chan core.Entry) error {
	table := tablewriter.NewWriter(s.writer)
	table.SetHeader(s.fields)

	for entry := range in {
		row := make([]string, len(s.fields))
		for i, field := range s.fields {
			row[i] = entry[field]
		}
		table.Append(row)
	}
	table.Render()
	return nil
}
