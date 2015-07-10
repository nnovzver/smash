package main

import (
	"encoding/json"
	"os"
)

type Field struct {
	Name   string
	Length float64
}

type Codogram struct {
	Name   string
	Fields []Field
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

func main() {

}
