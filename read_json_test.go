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

var simpleModuleCode Module = Module{
	Name: "Simple",
	Codograms: []Codogram{
		{
			Name: "First",
			Fields: []Field{
				{Name: "i", Length: 2, CType: "uint8_t"},
				{Name: "j", Length: 7, CType: "uint8_t"},
				{Name: "k", Length: 16, CType: "uint16_t"},
				{Name: "l", Length: 32, CType: "uint32_t"},
			},
			CLength: 8,
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

func TestAddCLengths(t *testing.T) {
	m := simpleModule
	m.AddCLengths()
	res := m.Codograms[0].CLength
	orig := simpleModuleCode.Codograms[0].CLength
	if res != orig {
		t.Errorf("Unexpected CLength %d != %d\n", res, orig)
	}
}

func TestAddCTypes(t *testing.T) {
	m := simpleModule
	m.AddCTypes()
	res := m.Codograms[0].Fields
	orig := simpleModuleCode.Codograms[0].Fields
	if !reflect.DeepEqual(res, orig) {
		t.Errorf("Unexpected CType\n res = %+v\n orig = %+v\n", res, orig)
	}
}

func TestGenerateDotH(t *testing.T) {
	m := simpleModule
	err := m.AddCTypes()
	if err != nil {
		t.Error(err)
	}
	g, err := m.GenerateDotH()
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
