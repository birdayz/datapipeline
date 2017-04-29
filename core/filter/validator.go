package filter

import "github.com/birdayz/datapipeline/core"

// A Validator defines the logic how an Entry is validated. It returns true if
// the Entry is valid.
type Validator interface {
	Validate(entry core.Entry) bool
}

// ValidatingFilter implements the core.Pipe interface, hence it is a stage in a
// pipeline. ValidatingFilter defines how Validation is orchestrated. This is a
// naive implementation: Entries are retrieved from the input channel. For each
// Entry retrieved, the given validators are called synchronously until a
// validation error occurs or all validators are passed. This procedure is
// performed for each Entry in the input channel, until the channel is closed.
// The policy in case of invalid Entries is, that these Entries are simply
// dropped (not passed to the output channel).
type ValidatingFilter struct {
	validators []Validator
}

// NewValidator constructs a new Pipe, which uses the given validators to decide
// if incoming entries should be dropped.
func NewValidator(validators ...Validator) core.Pipe {
	return &ValidatingFilter{validators: validators}
}

// Process items from in chan, validate, pass to out chan if valid.
func (vf *ValidatingFilter) Process(in <-chan core.Entry) <-chan core.Entry {
	out := make(chan core.Entry)
	go func() {
	entries:
		for entry := range in {
			for _, validator := range vf.validators {
				if valid := validator.Validate(entry); !valid {
					continue entries
				}
			}
			out <- entry
		}
		close(out)
	}()
	return out
}
