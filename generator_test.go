package main

import (
	"reflect"
	"strings"
	"testing"
)

func simplifyString(s string) string {
	var r string
	for _, l := range strings.Split(s, "\n") {
		l = strings.Trim(l, " \n\t")
		if len(l) != 0 {
			r += l + "\n"
		}
	}
	return r
}

var simpleModule Module = Module{
	Name: "Simple",
	Codograms: []Codogram{
		{
			Name: "First",
			Fields: []Field{
				{Name: "i", Length: 2, Type: "const", Const: 2},
				{Name: "j", Length: 7, Type: "const", Const: 4},
				{Name: "k", Length: 16, Type: "simple", Enum: "APPLE:1, PEAR:2"},
				{Name: "l", Length: 32, Type: "simple"},
				{Name: "t", Length: 7, Type: "temp"},
				{Name: "b", Length: 128, Type: "blob"},
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
				{Name: "t", Length: 7},
				{Name: "b", Length: 128},
			},
			CLength: 24,
			CMarshal: `
ch[0] |= (c->i<<6)&MASK(7, 6);

ch[0] |= (c->j>>1)&MASK(5, 0);
ch[1] |= (c->j<<7)&MASK(7, 7);

ch[1] |= (c->k>>9)&MASK(6, 0);
ch[2] |= (c->k>>1)&MASK(7, 0);
ch[3] |= (c->k<<7)&MASK(7, 7);

ch[3] |= (c->l>>25)&MASK(6, 0);
ch[4] |= (c->l>>17)&MASK(7, 0);
ch[5] |= (c->l>>9)&MASK(7, 0);
ch[6] |= (c->l>>1)&MASK(7, 0);
ch[7] |= (c->l<<7)&MASK(7, 7);
`,
			CUnmarshal: `
c->i |= (ch[0]>>6)&MASK(1, 0);

c->j |= (ch[0]<<1)&MASK(6, 1);
c->j |= (ch[1]>>7)&MASK(0, 0);

c->k |= (ch[1]<<9)&MASK(15, 9);
c->k |= (ch[2]<<1)&MASK(8, 1);
c->k |= (ch[3]>>7)&MASK(0, 0);

c->l |= (ch[3]<<25)&MASK(31, 25);
c->l |= (ch[4]<<17)&MASK(24, 17);
c->l |= (ch[5]<<9)&MASK(16, 9);
c->l |= (ch[6]<<1)&MASK(8, 1);
c->l |= (ch[7]>>7)&MASK(0, 0);
`,
			CTest: `
i |= (ch[0]>>6)&MASK(1, 0);

j |= (ch[0]<<1)&MASK(6, 1);
j |= (ch[1]>>7)&MASK(0, 0);
`,

			CMacros: `
#define First__BUFSIZE 24
#define First__I 2
#define First__J 4
#define First__APPLE 1
#define First__PEAR 2
#define First__b_OFFSET 8
#define First__b_SIZE 16
`,
		},
	},
}

func TestParseJsonModule(t *testing.T) {
	m, err := parseJsonModule("tests/simple_proto.json")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(m, simpleModule) {
		t.Errorf("m, orig not equal\n m =\n%+v\n orig =\n%+v\n", m, simpleModule)
	}
}

func TestAddCLengths(t *testing.T) {
	m := simpleModule
	m.addCLengths()
	res := m.Codograms[0].CLength
	orig := simpleModuleCode.Codograms[0].CLength
	if res != orig {
		t.Errorf("Unexpected CLength %d != %d\n", res, orig)
	}
}

func TestAddCTypes(t *testing.T) {
	m := simpleModule
	m.addCTypes()
	for i, f := range m.Codograms[0].Fields {
		res := f.CType
		orig := simpleModuleCode.Codograms[0].Fields[i].CType
		if !reflect.DeepEqual(res, orig) {
			t.Errorf("Unexpected CType\n res = %+v\n orig = %+v\n", res, orig)
		}
	}
}

func TestAddCMarshal(t *testing.T) {
	m := simpleModule
	m.addCCode()
	res := simplifyString(m.Codograms[0].CMarshal)
	orig := simplifyString(simpleModuleCode.Codograms[0].CMarshal)
	if res != orig {
		t.Errorf("Unexpected CMarshal\n res = %s\n orig = %s\n", res, orig)
	}
}

func TestAddCUnmarshal(t *testing.T) {
	m := simpleModule
	m.addCCode()
	res := simplifyString(m.Codograms[0].CUnmarshal)
	orig := simplifyString(simpleModuleCode.Codograms[0].CUnmarshal)
	if res != orig {
		t.Errorf("Unexpected CUnmarshal\n res = %s\n orig = %s\n", res, orig)
	}
}

func TestAddCTest(t *testing.T) {
	m := simpleModule
	m.addCCode()
	res := simplifyString(m.Codograms[0].CTest)
	orig := simplifyString(simpleModuleCode.Codograms[0].CTest)
	if res != orig {
		t.Errorf("Unexpected CTest\n res =\n %s\n orig =\n %s\n", res, orig)
	}
}

func TestAddCMacros(t *testing.T) {
	m := simpleModule
	m.addCCode()
	res := simplifyString(m.Codograms[0].CMacros)
	orig := simplifyString(simpleModuleCode.Codograms[0].CMacros)
	if res != orig {
		t.Errorf("Unexpected CMacros\n res =\n %s\n orig =\n %s\n", res, orig)
	}
}

func TestErrorAlignedBlob(t *testing.T) {
	_, err := parseJsonModule("tests/error_aligned_blob.json")
	if err == nil {
		t.Error(err)
	}
}
