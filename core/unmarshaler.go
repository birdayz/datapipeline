package core

// An Unmarshaler reads data from an arbitrary source and transforms it to
// Entries. Each entry is passed to the given channel. This channel may be the
// input channel of a pipeline. Only "real" errors should be returned, so e.g.
// io.EOF should not be returned because it does not indicate that something
// went wrong.
type Unmarshaler interface {
	Unmarshal(out chan<- Entry) error
}
