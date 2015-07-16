package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

const (
	ConstId = 1
)

func GetConstId() float64 {
	return float64(ConstId)
}

type Field struct {
	Name   string
	Length float64
	Type   float64
	Const  float64
	CType  string
}

type Codogram struct {
	Name       string
	Fields     []Field
	CLength    int
	CMarshal   string
	CUnmarshal string
	CTest      string
}

type Module struct {
	Name      string
	Codograms []Codogram
}

func parseJsonModule(filename string) (Module, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Module{}, err
	}

	dec := json.NewDecoder(f)
	var m Module
	err = dec.Decode(&m)
	if err != nil {
		return Module{}, err
	}

	return m, nil
}

func (m *Module) addCTypes() error {
	for i, c := range m.Codograms {
		for ii, f := range c.Fields {
			var ctype string
			switch {
			case 0 < f.Length && f.Length <= 8:
				ctype = "uint8_t"
			case f.Length <= 16:
				ctype = "uint16_t"
			case f.Length <= 32:
				ctype = "uint32_t"
			default:
				return fmt.Errorf("unexpected Field.Length")
			}
			m.Codograms[i].Fields[ii].CType = ctype
		}
	}
	return nil
}

func (m *Module) addCLengths() {
	var len int = 0
	for i, c := range m.Codograms {
		for _, f := range c.Fields {
			len += int(f.Length)
		}
		m.Codograms[i].CLength = (len + 8) / 8
	}

}

func (m *Module) addCCode() {
	whereToShift := func(freeBitsInByte, highBitInField int) int {
		return freeBitsInByte - highBitInField - 1
	}

	for cidx, c := range m.Codograms {
		freeBitsInByte := 8
		byteIndex := 0
		var marshalCode string
		var unmarshalCode string
		var testCode string
		for _, f := range c.Fields {
			bitsToMarshal := int(f.Length)
			bytesForField := ((bitsToMarshal - freeBitsInByte) + 15) / 8

			for i := 0; i < bytesForField; i++ {
				fieldShift := whereToShift(freeBitsInByte, bitsToMarshal-1)
				var maskEnd int
				var shiftSymbol string
				if fieldShift < 0 {
					// field is not fully packed
					fieldShift = -fieldShift
					shiftSymbol = ">>"
					maskEnd = 0
				} else {
					// field is fully packed
					shiftSymbol = "<<"
					maskEnd = freeBitsInByte - bitsToMarshal
				}
				// generate Marshal
				marshalCode += fmt.Sprintf("  ch[%d] |= (c->%s%s%d)&MASK(%d, %d);\n",
					byteIndex, f.Name, shiftSymbol, fieldShift,
					freeBitsInByte-1, maskEnd)
				// generate Unmarshal
				var unmarshalCodeTemp string
				if shiftSymbol == ">>" {
					unmarshalCodeTemp = fmt.Sprintf("  c->%s |= (ch[%d]%s%d)&MASK(%d, %d);\n",
						f.Name, byteIndex, "<<", fieldShift,
						freeBitsInByte-1+fieldShift, maskEnd+fieldShift)
				} else {
					unmarshalCodeTemp = fmt.Sprintf("  c->%s |= (ch[%d]%s%d)&MASK(%d, %d);\n",
						f.Name, byteIndex, ">>", fieldShift,
						freeBitsInByte-1-fieldShift, maskEnd-fieldShift)
				}
				unmarshalCode += unmarshalCodeTemp
				// generate Test
				if f.Type == ConstId {
					testCode += strings.Replace(unmarshalCodeTemp, "c->", "", 1)
				}

				if maskEnd == 0 {
					// field is not fully packed
					byteIndex++
					bitsToMarshal -= freeBitsInByte
					freeBitsInByte = 8
				} else {
					// field is fully packed
					freeBitsInByte -= bitsToMarshal
				}
			}
			marshalCode += "\n"
			unmarshalCode += "\n"
			testCode += "\n"
		}
		m.Codograms[cidx].CMarshal = marshalCode
		m.Codograms[cidx].CUnmarshal = unmarshalCode
		m.Codograms[cidx].CTest = testCode
	}
}

func GenerateCFiles(jfilename string, wr io.Writer) error {
	m, err := parseJsonModule(jfilename)
	if err != nil {
		return err
	}
	err = m.addCTypes()
	if err != nil {
		return err
	}
	m.addCLengths()
	m.addCCode()

	ht := template.New("h.template")
	ht, err = ht.ParseFiles("h.template")
	if err != nil {
		return err
	}
	err = ht.Execute(wr, m)
	if err != nil {
		return err
	}

	ct := template.New("c.template")
	ct = ct.Funcs(template.FuncMap{"getConstId": GetConstId})
	ct, err = ct.ParseFiles("c.template")
	if err != nil {
		return err
	}
	err = ct.Execute(wr, m)
	if err != nil {
		return err
	}

	return nil
}
