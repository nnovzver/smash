package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

	var hfile = os.Stdout
	var cfile = os.Stdout
	if !stdOut {
		hfile, err = os.Create(strings.TrimSuffix(filename, "json") + "h")
		if err != nil {
			panic(err)
		}
		cfile, err = os.Create(strings.TrimSuffix(filename, "json") + "c")
		if err != nil {
			panic(err)
		}
	}
	err = GenerateCFiles(filename, hfile, cfile)
	if err != nil {
		panic(err)
	}
}
