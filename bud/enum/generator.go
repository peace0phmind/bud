package enum

import (
	"embed"
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/factory"
	"github.com/peace0phmind/bud/stream"
	"github.com/peace0phmind/bud/util"
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

	funcs := template.FuncMap{}
	funcs["IA"] = util.IndefiniteArticle
	tmpl.Funcs(funcs)

	result.Tmpl = template.Must(tmpl.ParseFS(enumTmpl, "*.tmpl"))

	result.DataList = stream.Must(stream.Of(allEnums).Sort(func(x, y *Enum) int { return strings.Compare(x.Name, y.Name) }).ToSlice())

	return result
}

func (eg *EnumGenerator) GetImports() []string {
	return []string{"errors", "fmt"}
}

func NewGenerator(fileNode *goast.File, fileSet *token.FileSet) (ast.Generator, error) {
	// get global enum config
	var ec *Config = nil
	for _, cg := range fileNode.Comments {
		if strings.HasPrefix(cg.List[len(cg.List)-1].Text, "//go:generate") {
			ag, err := ast.ParseAnnotation(fileNode.Name.Name, cg.Text())
			if err != nil {
				return nil, fmt.Errorf("parse annotation err: %v", err)
			}

			ec1, err := annotationGroupToEnumConfig(ag, nil)
			if err != nil {
				return nil, fmt.Errorf("parse annotation to enum config err: %v", err)
			}
			ec = ec1
			break
		}
	}

	if ec == nil {
		ec = factory.New[Config]()
	}

	// get all enums
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

			enum, err := annotationGroupToEnum(ag, ts, ec)
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
