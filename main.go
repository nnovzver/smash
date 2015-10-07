package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var stdOut bool
var hOnly bool
var cOnly bool
var outputDir string

func init() {
	flag.BoolVar(&stdOut, "s", false, "print generated code to stdout")
	flag.BoolVar(&hOnly, "h", false, "generate .h header only")
	flag.BoolVar(&cOnly, "c", false, "generate .c source only")
	flag.StringVar(&outputDir, "o", "", "output directory. Default the same directory where .json")
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

	if strings.HasSuffix(filename, ".json") != true {
		fmt.Fprintf(os.Stderr, "File extension should be .json: %s\n", filename)
	}

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
