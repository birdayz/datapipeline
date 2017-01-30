package core

type nopFilter struct{}

// NewNop constructs a new No-Operation filter.
func NewNop() Pipe {
	return &nopFilter{}
}

// Process implements the core.Pipe interface.
func (n *nopFilter) Process(in <-chan Entry) <-chan Entry {
	out := make(chan Entry)
	go func() {
		for entry := range in {
			out <- entry
		}
		close(out)
	}()
	return out
}
