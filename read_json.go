package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type Field struct {
	Name   string
	Length float64
	CType  string
}

type Codogram struct {
	Name    string
	Fields  []Field
	CLength int
}

type Module struct {
	Name      string
	Codograms []Codogram
}

func ParseJsonModule(filename string) (Module, error) {
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

func (m *Module) AddCTypes() error {
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
				fmt.Errorf("unexpected Field.Length")
			}
			m.Codograms[i].Fields[ii].CType = ctype
		}
	}
	return nil
}

func (m *Module) AddCLengths() {
	var len int = 0
	for i, c := range m.Codograms {
		for _, f := range c.Fields {
			len += int(f.Length)
		}
		m.Codograms[i].CLength = (len + 8) / 8
	}

}

func (m *Module) GenerateDotH() (string, error) {
	t, err := template.ParseFiles("h.template")
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.Execute(&b, m)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func main() {

}
