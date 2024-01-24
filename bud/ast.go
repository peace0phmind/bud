package bud

import (
	"fmt"
	"github.com/peace0phmind/bud/bud/enum"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func GenerateFromFile(inputFile string) ([]byte, error) {
	f, err := parser.ParseFile(token.NewFileSet(), inputFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("generate: error parsing input file '%s': %s", inputFile, err)
	}

	enums := InspectAnnotation(f, AnnotationEnum)

	return enum.Generate(enums)
}

func InspectAnnotation(node ast.Node, annotation AnnotationType) map[string]*ast.TypeSpec {
	annotations := make(map[string]*ast.TypeSpec)
	// Inspect the AST and find all structs.
	ast.Inspect(node, func(n ast.Node) bool {
		//fmt.Printf("Node: %#v\n", n)
		switch x := n.(type) {
		case *ast.GenDecl:
			copyGenDeclCommentsToSpecs(x)
		case *ast.FuncDecl:
			copyFuncDeclCommentsToSpecs(x)
		case *ast.Ident:
			if x.Obj != nil {
				if x.Obj.Kind == ast.Typ {
					if ts, ok := x.Obj.Decl.(*ast.TypeSpec); ok {
						if ts.Doc != nil {
							for _, comment := range ts.Doc.List {
								if strings.Contains(comment.Text, string(annotation)) {
									annotations[x.Name] = ts
									break
								}
							}
						}
					}
				}
			}
		case *ast.Field:
			if len(x.Names) > 0 && x.Names[0].Name == "b" {
				if x.Doc != nil {
					comment := x.Doc.Text()
					log.Printf("Found comment for Param2: %s", comment)
				}
			}
		}

		//Return true to continue through the tree
		return true
	})

	return annotations
}

// copied from github.com/abice/go-enum
// copyDocsToSpecs will take the GenDecl level documents and copy them
// to the children Type and Value specs.  I think this is actually working
// around a bug in the AST, but it works for now.
func copyGenDeclCommentsToSpecs(x *ast.GenDecl) {
	// Copy the doc spec to the type or value spec
	// cause they missed this... whoops
	if x.Doc != nil {
		for _, spec := range x.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				if s.Doc == nil {
					s.Doc = x.Doc
				}
			case *ast.ValueSpec:
				if s.Doc == nil {
					s.Doc = x.Doc
				}
			}
		}
	}
}

func copyFuncDeclCommentsToSpecs(x *ast.FuncDecl) {
	// Copy the doc spec to the type or value spec
	// cause they missed this... whoops
	if x.Doc != nil {
		fmt.Sprintf("%+v", x)
	}
}
