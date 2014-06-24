package multistream

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestTypeWriter(t *testing.T) {
	b := bytes.NewBuffer(nil)
	w := newTypeWriter("stdout", b)
	w.Write([]byte("hello"))
	expected, _ := json.Marshal(&data{
		Type: "stdout",
		Data: []byte("hello"),
	})
	if got := b.Bytes(); !bytes.Equal(got, expected) {
		t.Fatal("Expected %v, got %v", string(expected), string(got))
	}
}
