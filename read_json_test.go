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

	if !reflect.DeepEqual(m, simpleModule) {
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
