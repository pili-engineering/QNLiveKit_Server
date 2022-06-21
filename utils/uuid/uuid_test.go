package uuid

import (
	"testing"
)

func TestUUID(t *testing.T) {
	Init(1)

	id := Gen()

	ti := GetTimeFromUUID(id)

	if ti.IsZero() {
		t.Fatal("gen uuid failed")
	}

	tt := GetTimeFromUUID("")
	if !tt.IsZero() {
		t.Fatal("unexpected response")
	}
}
