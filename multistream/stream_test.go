package multistream

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func contain(got []data, item data) bool {
	for _, g := range got {
		if reflect.DeepEqual(g, item) {
			return true
		}
	}
	return false
}

// assert arriving data that can be out of order
func assertData(t *testing.T, output string, expected []data) {
	d := json.NewDecoder(bytes.NewBufferString(output))
	got := []data{}
	for {
		g := data{}
		if err := d.Decode(&g); err != nil {
			break
		}
		got = append(got, g)
	}
	if len(got) != len(expected) {
		t.Fatalf("Expected %v to have %v items, got %v", output,
			len(expected), len(got))
	}
	for _, e := range expected {
		if !contain(got, e) {
			t.Fatalf("Expected to have item %v-%v", e.Type, string(e.Data))
		}
	}
}

func TestMultiStream(t *testing.T) {
	output := bytes.NewBuffer(nil)
	stdout := strings.NewReader("helloworld")
	stderr := strings.NewReader("newerror")

	w := MultiStream(output, map[string]io.Reader{
		"stdout": stdout,
		"stderr": stderr,
	})

	errDone := make(chan error, 1)
	go func() {
		errDone <- w.Flush()
	}()

	after := time.After(10 * time.Millisecond)
	select {
	case <-after:
		t.Fatal("timed out")
	case err := <-errDone:
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		strout := strings.TrimSpace(output.String())
		assertData(t, strout, []data{
			{Type: "stdout", Data: []byte("helloworld")},
			{Type: "stderr", Data: []byte("newerror")},
		})
	}
}
