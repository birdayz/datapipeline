package core

// Pipeline exposes an input channel and an output channel.
type Pipeline struct {
	In  chan<- Entry
	Out <-chan Entry
}

// An Entry represents one item that may be passed into a pipe(line). Each Entry
// contains an arbitrary number of key/value pairs.
type Entry map[string]string

// Pipe reads entries from a given input channel. The pipe returns an output
// channel. Input entries should be processed in a goroutine. Entries may be
// omitted/filtered, hence not all entries retrieved from the input channel must
// be passed to the output channel. Implementations MUST close the output
// channel once the input channel is closed (e.g. after a range loop over the
// input channel). Process must be called before entries are written into the
// input channel, otherwise adding items to the input channel blocks.
type Pipe interface {
	Process(in <-chan Entry) (out <-chan Entry)
}

// NewPipeline constructs a Pipeline by creating an input channel, and
// connecting the input channel to the first pipe. The output channel of the
// pipeline is the output channel of the last pipe. All pipes are connected in
// the same order as passed to the NewPipeline(). The input channel MUST be
// closed by the caller once all entries are sent. Pipes must NOT be nil.
func NewPipeline(pipes ...Pipe) *Pipeline {
	pipelineIn := make(chan Entry)

	var nextPipeIn <-chan Entry = pipelineIn
	for _, pipe := range pipes {
		nextPipeIn = pipe.Process(nextPipeIn)
	}

	p := &Pipeline{In: pipelineIn, Out: nextPipeIn}
	return p
}
