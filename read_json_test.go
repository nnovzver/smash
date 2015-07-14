package main

import (
	"reflect"
	"strings"
	"testing"
)

var simpleModule Module = Module{
	Name: "Simple",
	Codograms: []Codogram{
		{
			Name: "First",
			Fields: []Field{
				{Name: "i", Length: 2},
				{Name: "j", Length: 7},
				{Name: "k", Length: 16},
				{Name: "l", Length: 32},
			},
		},
	},
}

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

func TestGenerateDotH(t *testing.T) {
	m := simpleModule
	err := m.AddCTypes()
	if err != nil {
		t.Error(err)
	}
	g, err := GenerateDotH(&m)
	if err != nil {
		t.Error(err)
	}
	orig := `
typedef struct First {
  uint8_t i; // 2 bits
  uint8_t j; // 7 bits
  uint16_t k; // 16 bits
  uint32_t l; // 32 bits
} First;
`
	replacer := strings.NewReplacer(" ", "",
		"\t", "",
		"\n", "")
	g = replacer.Replace(g)
	orig = replacer.Replace(orig)
	if g != orig {
		t.Errorf("g != orig:\n%s\n%s", g, orig)
	}
}
