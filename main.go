package main

import (
	"flag"
	"fmt"
	"github.com/peace0phmind/bud/bud"
	"os"
)

func main() {
	var filename string

	flag.StringVar(&filename, "file", "", "The file to generate bud file.")

	flag.Parse()

	if len(filename) == 0 {
		filename, _ = os.LookupEnv("GOFILE")

		if len(filename) == 0 {
			fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			return
		}
	}

	bud.GenerateFile(filename, "_bud")
}
