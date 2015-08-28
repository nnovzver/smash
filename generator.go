package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ConstId  = "const"
	SimpleId = "simple"
	BlobId   = "blob"
	TempId   = "temp"
)

func GetConstId() string {
	return ConstId
}

func GetBlobId() string {
	return BlobId
}

func GetTempId() string {
	return TempId
}

func BytesInBits(bits float64) float64 {
	return bits / 8
}

type Field struct {
	Name       string
	Length     float64
	Type       string
	Const      float64
	Enum       string
	CType      string
	BlobOffset int
	BlobSize   int
}

type Codogram struct {
	Name       string
	Fields     []Field
	CLength    int
	CMarshal   string
	CUnmarshal string
	CTest      string
	CMacros    string
}

type Module struct {
	Name      string
	FileName  string
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

	err = m.checkModule()
	if err != nil {
		return Module{}, err
	}

	return m, nil
}

func (m *Module) checkModule() error {
	// check blob boundaries
	for _, c := range m.Codograms {
		bitsBefore := 0
		for _, f := range c.Fields {
			if f.Type == BlobId && (bitsBefore%8 != 0 || int(f.Length)%8 != 0) {
				return fmt.Errorf("Blob %s in codogram %s not aligned!", f.Name, c.Name)
			}
			bitsBefore += int(f.Length)
		}
	}

	return nil
}

func (m *Module) addCTypes() error {
	for i, c := range m.Codograms {
		for ii, f := range c.Fields {
			var ctype string
			if f.Type == SimpleId || f.Type == ConstId {
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
		m.Codograms[i].CLength = (len + 7) / 8
		len = 0
	}

}

func (m *Module) addFileName(filename string) {
	basename := strings.TrimSuffix(filepath.Base(filename), ".json")
	m.FileName = basename
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
		var macrosCode string

		// add codogram size macros
		codogramSize := 0
		for _, f := range c.Fields {
			codogramSize += int(f.Length)
		}
		macrosCode += fmt.Sprintf("#define %s__SIZE %d\n",
			c.Name, codogramSize/8)

		for fidx, f := range c.Fields {
			// generate Macros
			if f.Type == ConstId {
				macrosCode += fmt.Sprintf("#define %s__%s %d\n",
					c.Name, strings.ToUpper(f.Name), int(f.Const))
			}
			if f.Type == BlobId {
				macrosCode += fmt.Sprintf("#define %s__%s_OFFSET %d\n",
					c.Name, f.Name, byteIndex)
				macrosCode += fmt.Sprintf("#define %s__%s_SIZE %d\n",
					c.Name, f.Name, int(f.Length)/8)

				m.Codograms[cidx].Fields[fidx].BlobOffset = byteIndex
				m.Codograms[cidx].Fields[fidx].BlobSize = int(f.Length) / 8
				continue
			}
			if f.Enum != "" {
				for _, s := range strings.Split(f.Enum, ",") {
					val := strings.Split(s, ":")
					if val != nil {
						macrosCode += fmt.Sprintf("#define %s__%s %s\n",
							c.Name,
							strings.Trim(val[0], " \t\n"),
							strings.Trim(val[1], " \t\n"))
					}
				}
			}

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
				// generate Marshal if field not temporary
				if f.Type != TempId {
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
		}

		m.Codograms[cidx].CMarshal = marshalCode
		m.Codograms[cidx].CUnmarshal = unmarshalCode
		m.Codograms[cidx].CTest = testCode
		m.Codograms[cidx].CMacros = macrosCode
	}
}

func GenerateCFiles(jfilename string) (string, error) {
	// parse json module
	m, err := parseJsonModule(jfilename)
	if err != nil {
		return "", err
	}

	err = m.addCTypes()
	if err != nil {
		return "", err
	}
	m.addCLengths()
	m.addCCode()
	m.addFileName(jfilename)

	// create buffers
	hbuf := new(bytes.Buffer)
	cbuf := new(bytes.Buffer)

	// generate code
	ht := template.New("h_template")
	ht = ht.Funcs(template.FuncMap{"getBlobId": GetBlobId})
	ht = ht.Funcs(template.FuncMap{"getTempId": GetTempId})
	ht = ht.Funcs(template.FuncMap{"bytesInBits": BytesInBits})
	ht, err = ht.Parse(h_template)
	if err != nil {
		return "", err
	}
	err = ht.Execute(hbuf, m)
	if err != nil {
		return "", err
	}

	ct := template.New("c_template")
	ct = ct.Funcs(template.FuncMap{"getConstId": GetConstId})
	ct = ct.Funcs(template.FuncMap{"getBlobId": GetBlobId})
	ct, err = ct.Parse(c_template)
	if err != nil {
		return "", err
	}
	err = ct.Execute(cbuf, m)
	if err != nil {
		return "", err
	}

	// string for stdout
	var str string

	// create files
	if hOnly == cOnly || hOnly {
		str += string(hbuf.String())
		hfile, err := os.Create(filepath.Dir(jfilename) + "/" + m.FileName + ".gen.h")
		if err != nil {
			return "", err
		}

		_, err = hbuf.WriteTo(hfile)
		if err != nil {
			return "", err
		}
	}

	if hOnly == cOnly || cOnly {
		str += string(cbuf.String())
		cfile, err := os.Create(filepath.Dir(jfilename) + "/" + m.FileName + ".gen.c")
		if err != nil {
			return "", err
		}

		_, err = cbuf.WriteTo(cfile)
		if err != nil {
			return "", err
		}
	}

	return str, nil
}
