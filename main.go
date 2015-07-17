package main

import (
	"flag"
	"fmt"
	"os"
)

var stdOut bool

func init() {
	flag.BoolVar(&stdOut, "s", false, "print generated code to stdout")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] json_module_filename\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}
	filename := flag.Args()[0]

	_, err := os.Stat(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No such file or directory: %s\n", filename)
		flag.Usage()
		return
	}

	code, err := GenerateCFiles(filename)
	if err != nil {
		panic(err)
	}
	if stdOut {
		fmt.Println(code)
	}
}
