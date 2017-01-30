/*
Package core implements a pipeline which consists of an input channel, an
output channel and a number of pipes in between. Entries written into the input
channel are passed through each pipe until the output channel is reached. Each
pipe may or may not pass entries to the next pipe.
    ____________________________
 ->(______|______|______|______()->
 In  Pipe   Pipe   Pipe   Pipe   Out
*/
package core
