package enum

import (
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	goast "go/ast"
	"strings"
)

func GenerateFromFile(inputFile string) ([]byte, error) {
	fileNode, fileSet, err := ast.ParseFile(inputFile)
	if err != nil {
		return nil, err
	}

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
	}

	return nil, nil
}
