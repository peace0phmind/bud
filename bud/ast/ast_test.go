package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// @Enum eabc
type Eabc int

// @Enum abc
type Abc struct {
}

func (abc *Abc) Init(
	// @Enum a
	a int,
	b int, // @Enum b
) {
}

func TestAAA(t *testing.T) {
	const src = `package main

//ExampleFunction is an example
func ExampleFunction(
    //Param1 is the first param
    Param1 int,
	// Param2 is for aaa
    // @Enum Param2 is the second param
    Param2 string, // @123
) {
    //function body
}`
	// 创建一个新的token.FileSet
	fset := token.NewFileSet()

	// 解析源码以获得AST
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	//file, err := ast.ParseFile(fset, "./annotation.go", nil, ast.ParseComments)
	if err != nil {
		panic(err)
	}

	// 遍历AST中的所有声明
	for _, decl := range file.Decls {
		// 确保声明是函数
		if fn, ok := decl.(*ast.FuncDecl); ok {
			// 遍历函数的参数
			for _, field := range fn.Type.Params.List {
				for _, name := range field.Names {
					// 找到Param2参数
					if name.Name == "Param2" {
						// 获取Param2参数的位置
						doc, comment := FindCommentGroup(file, fset, field.Pos())
						if doc != nil {
							fmt.Println(strings.TrimSpace(doc.Text()))
						}
						if comment != nil {
							fmt.Println(strings.TrimSpace(comment.Text()))
						}
					}
				}
			}
		}
	}
}
