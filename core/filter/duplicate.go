package filter

import "datapipeline/core"

type duplicateFilter struct {
	dedupField string
}

// NewDuplicate instantiates a new filter which removes duplicates in the
// pipeline. Duplicates are identified by the given dedupField. Only one entry
// with the same value for the given field is allowed. In case of duplicates,
// the first entry received from the input channel has priority.
func NewDuplicate(dedupField string) core.Pipe {
	return &duplicateFilter{dedupField: dedupField}
}

func (df *duplicateFilter) Process(in <-chan core.Entry) <-chan core.Entry {
	out := make(chan core.Entry)

	go func() {
		fieldValueSeen := make(map[string]bool)
		for entry := range in {
			if dedupFieldValue, ok := entry[df.dedupField]; ok {
				if _, seen := fieldValueSeen[dedupFieldValue]; seen {
					// drop this entry, it is a duplicate
					continue
				}
				fieldValueSeen[entry[df.dedupField]] = true
			}
			out <- entry
		}
		close(out)
	}()
	return out
}
