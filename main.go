package main

import (
	"flag"
	"fmt"
	"github.com/peace0phmind/bud/bud"
	"os"
)

func main() {
	var filename string
	var fileSuffix string

	flag.StringVar(&filename, "file", "", "The file to generate bud file.")
	flag.StringVar(&fileSuffix, "file-suffix", "_bud", "Changes the default filename suffix of _bud to something else.")

	flag.Parse()

	if len(filename) == 0 {
		filename, _ = os.LookupEnv("GOFILE")

		if len(filename) == 0 {
			fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			return
		}
	}

	bud.GenerateFile(filename, fileSuffix)
}
