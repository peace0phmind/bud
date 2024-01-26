package enum

import (
	"github.com/peace0phmind/bud/bud/ast"
	goast "go/ast"
	"go/token"
)

type SingletonGenerator struct {
	ast.BaseGenerator[Singleton]
}

func newSingletonGenerator(allSingleton []*Singleton) (*SingletonGenerator, error) {
	result := &SingletonGenerator{}

	return result, nil
}

func (sg *SingletonGenerator) GetImports() []string {
	return []string{"fmt"}
}

func Generate(fileNode *goast.File, fileSet *token.FileSet) (ast.Generator, error) {
	return nil, nil
}
