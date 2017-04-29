package unmarshaler

import (
	"github.com/birdayz/datapipeline/core"
	"encoding/csv"
	"io"
)

// csvUnmarshaler reads from a CSV file and produces a core.Entry for each row
// in the CSV file.
type csvUnmarshaler struct {
	reader io.Reader
}

// NewCSV reads from a reader which outputs CSV contents. CSV file rows are
// transformed to the core.Entry datatype, with each column being a key/value
// pair.
func NewCSV(reader io.Reader) core.Unmarshaler {
	return &csvUnmarshaler{reader: reader}
}

// Unmarshal implements the core.Unmarshaler interface.
func (c *csvUnmarshaler) Unmarshal(out chan<- core.Entry) error {
	r := csv.NewReader(c.reader)
	defer close(out)

	// Read header first
	headers, err := r.Read()
	if err != nil && err != io.EOF {
		return err
	}
	numFields := len(headers)

	for record, err := r.Read(); err != io.EOF; record, err = r.Read() {
		if err != nil {
			return err
		}
		row := make(core.Entry, numFields)
		for idx, key := range headers {
			row[key] = record[idx]
		}
		out <- row
	}
	return nil
}
