package marshaler

import (
	"bytes"
	"datapipeline/core"
	"io"

	"gopkg.in/yaml.v2"
)

// jsonMarshaler marshals all entries to the YAML format.
type yamlMarshaler struct {
	writer io.Writer
}

// NewYAML constructs a new yamlMarshaler.
func NewYAML(writer io.Writer) core.Marshaler {
	return &yamlMarshaler{writer: writer}
}

// Marshal implements the core.Marshaler interface.
func (y *yamlMarshaler) Marshal(in <-chan core.Entry) error {
	var rows []map[string]string
	for entry := range in {
		rows = append(rows, entry)
	}

	// (jb) This is suboptimal: we have to marshal into []byte and
	// afterwards write it, because the yaml API does not expose a method to
	// directly marshal into a writer.
	data, err := yaml.Marshal(rows)
	if err != nil {
		return err
	}
	_, err = io.Copy(y.writer, bytes.NewReader(data))
	return err
}
