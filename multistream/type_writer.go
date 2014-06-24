package multistream

import (
	"encoding/json"
	"io"
)

type data struct {
	Type string
	Data []byte
}

type typeWriter struct {
	dataType string
	output   io.Writer
}

func (w *typeWriter) Write(p []byte) (int, error) {
	b, err := json.Marshal(&data{Type: w.dataType, Data: p})
	if err != nil {
		return 0, err
	}
	if _, err := w.output.Write(b); err != nil {
		return 0, err
	}
	return len(p), nil
}

func newTypeWriter(dataType string, output io.Writer) *typeWriter {
	return &typeWriter{
		dataType: dataType,
		output:   output,
	}
}
