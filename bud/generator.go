package bud

import (
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/bud/enum"
	"go/parser"
	"go/token"
)

func GenerateFromFile(inputFile string) ([]byte, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, inputFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("generate: error parsing input file '%s': %s", inputFile, err)
	}

	enums := ast.InspectAnnotation(f, fileSet)

	return enum.Generate(enums)
}
