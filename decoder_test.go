package rootio

import (
	"bytes"
	"testing"
)

func TestDecoder(t *testing.T) {
	data := make([]byte, 32)
	dec := NewDecoder(bytes.NewBuffer(data))

	if dec.Len() != 32 {
		t.Fatalf("expected len=%v. got %v", len(data), dec.Len())
	}
	start := dec.Pos()
	if start != 0 {
		t.Fatalf("expected start=%v. got %v", 0, start)
	}

	var x int16
	err := dec.readBin(&x)
	if err != nil {
		t.Fatalf("error reading int16: %v", err)
	}

	pos := dec.Pos()
	if pos != 2 {
		t.Fatalf("expected pos=%v. got %v", 16, pos)
	}
}