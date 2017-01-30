package core

// A Marshaler is supposed to retrieve entries from the given input channel
// (which may be the output of a pipeline, for example) and marshals the
// contents in an arbitrary (implementation-specific) way.
type Marshaler interface {
	Marshal(in <-chan Entry) error
}
