package multistream

import (
	"fmt"
	"io"
)

type multiStream struct {
	output  io.Writer
	streams map[string]io.Reader
}

// merges output from streams into `output`, classified them as necessaries
// based on the provided name
func (c *multiStream) Flush() error {
	errDone := make(chan error, len(c.streams))
	for name, reader := range c.streams {
		outR, outW := io.Pipe()
		tw := newTypeWriter(name, outW)

		// add type to the stream output
		go func(reader io.Reader, tw io.Writer) {
			_, err := io.Copy(tw, reader)
			outW.CloseWithError(err)
		}(reader, tw)

		// write typed information back to the writer
		go func(outR io.Reader) {
			_, err := io.Copy(c.output, outR)
			errDone <- err
		}(outR)
	}

	errs := []error{}
	for i := 0; i < len(c.streams); i++ {
		err := <-errDone
		if err != nil && err != io.EOF {
			errs = append(errs)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("Got error %v", errs)
	}
	return nil
}

func MultiStream(output io.Writer, streams map[string]io.Reader) *multiStream {
	return &multiStream{
		output:  output,
		streams: streams,
	}
}
