package enum

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/peace0phmind/bud/util"
	"reflect"
)

const (
	EnumItemName  = "Name"
	EnumItemValue = "Value"
)

type EnumItem struct {
	enum              *Enum
	idx               int
	Name              string
	Value             any
	DocComment        string
	LineComment       string
	ExtendData        []any
	IsBlankIdentifier bool
}

// GetCodeName return the item name used in code
func (ei *EnumItem) GetCodeName() string {
	if ei.IsBlankIdentifier {
		return BlankIdentifier
	}

	casedName := ei.Name
	if ei.enum.Config.UseCamelCaseName {
		casedName = strcase.ToCamel(ei.Name)
	} else {
		casedName = util.Capitalize(ei.Name)
	}

	if ei.enum.Config.NoPrefix {
		return ei.enum.Config.Prefix + casedName
	} else {
		return ei.enum.Config.Prefix + ei.enum.Name + casedName
	}
}

// GetName return the item real name, default equals with the code name, or an extent named `Name`
func (ei *EnumItem) GetName() string {
	return ei.ExtendData[0].(string)
}

func (ei *EnumItem) GetConstLine() string {
	if ei.Value == nil {
		if ei.idx == 0 {
			return fmt.Sprintf("%s %s = iota", ei.GetCodeName(), ei.enum.Name)
		} else {
			return ei.GetCodeName()
		}
	} else {
		if ei.enum.Type == reflect.String {
			return fmt.Sprintf("%s %s = \"%s\"", ei.GetCodeName(), ei.enum.Name, ei.Value)
		} else {
			return fmt.Sprintf("%s %s = %v", ei.GetCodeName(), ei.enum.Name, ei.Value)
		}
	}
}
