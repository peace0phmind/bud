package bud

import (
	"bytes"
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/bud/enum"
	"golang.org/x/tools/imports"
	"os"
	"path/filepath"
	"strings"
)

func GenerateFile(filename string, outputSuffix string) {
	filename, _ = filepath.Abs(filename)

	fileNode, fileSet, err := ast.ParseFile(filename)
	if err != nil {
		panic(err)
	}

	eg, err := enum.NewGenerator(fileNode, fileSet)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer([]byte{})

	// write package
	buf.WriteString("package " + fileNode.Name.Name)
	buf.WriteString("\n\n")

	// write import
	importList := eg.GetImports()
	buf.WriteString("import (\n")
	for _, imp := range importList {
		buf.WriteString("\t\"" + imp + "\"")
		buf.WriteString("\n")
	}
	buf.WriteString(")\n\n")

	err = eg.WriteConst(buf)
	if err != nil {
		panic(err)
	}

	err = eg.WriteInitFunc(buf)
	if err != nil {
		panic(err)
	}
	buf.WriteString("\n")

	err = eg.WriteBody(buf)
	if err != nil {
		panic(err)
	}

	formatted, err := imports.Process(fileNode.Name.Name, buf.Bytes(), nil)
	if err != nil {
		panic(fmt.Errorf("generate: error formatting code %s\n\n%s", err, buf.String()))
	}

	outFilePath := fmt.Sprintf("%s%s.go", strings.TrimSuffix(filename, filepath.Ext(filename)), outputSuffix)
	if strings.HasSuffix(filename, "_test.go") {
		outFilePath = strings.Replace(outFilePath, "_test"+outputSuffix+".go", outputSuffix+"_test.go", 1)
	}

	mode := int(0o644)
	err = os.WriteFile(outFilePath, formatted, os.FileMode(mode))
	if err != nil {
		panic(fmt.Errorf("failed writing to file %s: %s", outFilePath, err))
	}
}
