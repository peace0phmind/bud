package enum

import (
	"embed"
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/stream"
	goast "go/ast"
	"go/token"
	"strings"
	"text/template"
)

//go:embed enum.tmpl
var enumTmpl embed.FS

type EnumGenerator struct {
	ast.BaseGenerator[Enum]
}

func newEnumGenerator(allEnums []*Enum) *EnumGenerator {
	result := &EnumGenerator{}

	tmpl := template.New("enum")
	tmpl = template.Must(tmpl.ParseFS(enumTmpl, "*.tmpl"))

	result.Tmpl = tmpl

	result.DataList = stream.Must(stream.Of(allEnums).Sort(func(x, y *Enum) int { return strings.Compare(x.Name, y.Name) }).ToSlice())

	return result
}

func (eg *EnumGenerator) GetImports() []string {
	return []string{"errors", "fmt"}
}

func NewGenerator(fileNode *goast.File, fileSet *token.FileSet) (ast.Generator, error) {

	allEnums := ast.InspectMapper[goast.TypeSpec, Enum](fileNode, fileSet, func(ts *goast.TypeSpec) *Enum {
		var comment string
		if ts.Doc != nil {
			comment = ts.Doc.Text()
		} else if ts.Comment != nil {
			comment = ts.Comment.Text()
		}

		if len(comment) > 0 && (strings.Contains(comment, "@e") || strings.Contains(comment, "@E")) {
			ag, err1 := ast.ParseAnnotation(ts.Name.Name, comment)
			if err1 != nil {
				fmt.Printf("parse annotation err: %v", err1)
				return nil
			}

			enum, err := annotationGroupToEnum(ag, ts)
			if err != nil {
				panic(fmt.Sprintf("update extends err: %v", err))
			}

			return enum
		}

		return nil
	})

	if len(allEnums) > 0 {
		// create enums
		return newEnumGenerator(allEnums), nil
	}

	return nil, nil
}
