package filter

import "datapipeline/core"
import "sort"

type sortFilter struct {
	field string
	order SortOrder
}

// NewSort instantiates a new filter, which gathers all entries of the input
// channel until the input channel is closed. Then, it sorts all entries
// according to the given field and order and writes the entries to the output
// channel.
func NewSort(field string, order SortOrder) core.Pipe {
	return &sortFilter{field: field, order: order}

}

// SortOrder describes either ascending or descending order of items
type SortOrder int

const (
	// OrderAscending describes an ascending sort order
	OrderAscending SortOrder = iota
	// OrderDescending describes a descending sort order
	OrderDescending
)

type entrySorter struct {
	field   string
	order   SortOrder
	entries []core.Entry
}

func (es entrySorter) Len() int {
	return len(es.entries)
}

func (es entrySorter) Less(i, j int) bool {
	if es.order == OrderAscending {
		return es.entries[i][es.field] < es.entries[j][es.field]
	}
	return es.entries[i][es.field] > es.entries[j][es.field]
}

func (es entrySorter) Swap(i, j int) {
	es.entries[j], es.entries[i] = es.entries[i], es.entries[j]
}

func (sf *sortFilter) Process(in <-chan core.Entry) <-chan core.Entry {
	out := make(chan core.Entry)
	go func() {
		var entries []core.Entry
		for entry := range in {
			entries = append(entries, entry)
		}
		sort.Sort(entrySorter{field: sf.field, order: sf.order, entries: entries})

		for _, entry := range entries {
			out <- entry
		}

		close(out)
	}()
	return out
}
