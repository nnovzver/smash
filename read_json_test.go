package main

import (
	"reflect"
	"testing"
)

func TestParseJsonModule(t *testing.T) {
	m, err := ParseJsonModule("examples/simple_proto.json")
	if err != nil {
		t.Error(err)
	}
	orig := Module{}
	orig.Name = "Simple"
	orig.Codograms = make([]Codogram, 1)
	orig.Codograms[0].Name = "First"
	orig.Codograms[0].Fields = make([]Field, 4)
	orig.Codograms[0].Fields[0].Name = "i"
	orig.Codograms[0].Fields[0].Length = 2
	orig.Codograms[0].Fields[1].Name = "j"
	orig.Codograms[0].Fields[1].Length = 7
	orig.Codograms[0].Fields[2].Name = "k"
	orig.Codograms[0].Fields[2].Length = 16
	orig.Codograms[0].Fields[3].Name = "l"
	orig.Codograms[0].Fields[3].Length = 32

	if !reflect.DeepEqual(m, orig) {
		t.Error("m, orig not equal")
	}
}
