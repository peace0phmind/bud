package bud

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

func TestName(t *testing.T) {
	GenerateFromFile("./annotation.go")
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
	//file, err := parser.ParseFile(fset, "./annotation.go", nil, parser.ParseComments)
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
						paramPos := fset.Position(field.Pos())

						// 遍历文件中的所有注释
						for _, commentGroup := range file.Comments {
							commentGroupPos := fset.Position(commentGroup.End())
							fmt.Printf("commentGroupPos: %#+v \n", commentGroupPos)
							for _, comment := range commentGroup.List {
								// 获取注释的位置
								commentPos := fset.Position(comment.Slash)
								fmt.Printf("commentPos: %#+v \n", commentPos)
								fmt.Printf("paramPos: %#+v \n", paramPos)

								// 如果注释在Param2参数之前，且在同一行或紧接在前一行
								if (commentPos.Line == paramPos.Line && paramPos.Offset < commentPos.Offset) ||
									(commentPos.Offset < paramPos.Offset && commentPos.Line+1 == paramPos.Line) {
									fmt.Println(strings.TrimSpace(comment.Text))
								}
							}
						}
					}
				}
			}
		}
	}
}
