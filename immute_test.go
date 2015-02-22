package immute

import (
	"fmt"
	"testing"
)

func TestIntSequence(t *testing.T) {
	seq := CreateMap(map[interface{}]interface{}{1: "sic", 3: "luc"})

	if seq == nil {
		t.Fatalf("Sequence does not match Sequencer interface", seq)
	}

	if seq.Length() != 2 {
		t.Fatalf("Sequence has wrong length it should be 2", seq, seq.Length())
	}

	filter = seq.Filter(func(i interface{}, k interface{}) {
		fmt.Println("k:", i, k)
		return i.(string) == "sic"
	}, func(c int, f interface{}) {})
}
