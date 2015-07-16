package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s json_module_filename\n", os.Args[0])
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		printUsage()
		return
	}
	filename := flag.Args()[0]

	_, err := os.Stat(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No such file or directory: %s\n", filename)
		return
	}

	hfile, err := os.Create(strings.TrimSuffix(filename, "json") + "h")
	if err != nil {
		panic(err)
	}
	cfile, err := os.Create(strings.TrimSuffix(filename, "json") + "c")
	if err != nil {
		panic(err)
	}
	err = GenerateCFiles(filename, hfile, cfile)
	if err != nil {
		panic(err)
	}
}
