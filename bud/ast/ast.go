package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func ParseFromFile(inputFile string) ([]byte, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, inputFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("generate: error parsing input file '%s': %s", inputFile, err)
	}

	InspectAnnotation(f, fileSet)

	return nil, nil
}

func InspectAnnotation(fileNode *ast.File, fileSet *token.FileSet) map[string]*ast.TypeSpec {
	annotations := make(map[string]*ast.TypeSpec)
	// Inspect the AST and find all structs.
	ast.Inspect(fileNode, func(n ast.Node) bool {
		//fmt.Printf("Node: %#v\n", n)
		switch x := n.(type) {
		case *ast.GenDecl:
		case *ast.FuncDecl:
		case *ast.Ident:
			if x.Obj != nil {
				if x.Obj.Kind == ast.Typ {
					if ts, ok := x.Obj.Decl.(*ast.TypeSpec); ok {
						if ts.Doc == nil || ts.Comment == nil {

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

func FindCommentGroup(fileNode *ast.File, fileSet *token.FileSet, pos token.Pos) (doc, comment *ast.CommentGroup) {
	indentPos := fileSet.Position(pos)

	for _, commentGroup := range fileNode.Comments {
		commentGroupPos := fileSet.Position(commentGroup.End())

		if commentGroupPos.Line == indentPos.Line && indentPos.Offset < commentGroupPos.Offset {
			comment = commentGroup
		}

		if commentGroupPos.Line+1 == indentPos.Line && commentGroupPos.Offset < indentPos.Offset {
			doc = commentGroup
		}
	}

	return
}
