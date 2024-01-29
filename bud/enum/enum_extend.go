package enum

import (
	"bytes"
	"fmt"
	"github.com/peace0phmind/bud/structure"
	"reflect"
	"strings"
)

type EnumExtend struct {
	enum                *Enum
	idx                 int
	enum2ExtendRendered bool
	extend2EnumRendered bool
	Name                string
	Type                reflect.Kind
	Comment             string
}

func (ee *EnumExtend) Enum() *Enum {
	return ee.enum
}

func (ee *EnumExtend) Enum2ExtendVarName() string {
	return fmt.Sprintf("_%sMap%s", ee.enum.Name, ee.Name)
}

func (ee *EnumExtend) Enum2ExtendMap() string {
	if ee.enum2ExtendRendered == true {
		return ""
	}

	if ee.Name == EnumItemName && ee.enum.Type == reflect.String {
		return ""
	}

	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(fmt.Sprintf("var %s = map[%s]%s{\n", ee.Enum2ExtendVarName(), ee.enum.Name, ee.Type.String()))
	if ee.idx == 0 {
		index := 0
		for _, item := range ee.enum.GetItems() {
			nextIndex := index + len(item.GetName())
			buf.WriteString(fmt.Sprintf("	%s: _%s%s[%d:%d],\n", item.GetCodeName(), ee.enum.Name, ee.Name, index, nextIndex))
			index = nextIndex
		}
	} else {
		for _, item := range ee.enum.GetItems() {
			switch ee.Type {
			case reflect.String:
				buf.WriteString(fmt.Sprintf("	%s: \"%s\",\n", item.GetCodeName(), structure.MustConvertTo[string](item.ExtendData[ee.idx])))
			default:
				buf.WriteString(fmt.Sprintf("	%s: %s,\n", item.GetCodeName(), structure.MustConvertToKind(item.ExtendData[ee.idx], ee.Type)))
			}
		}
	}
	buf.WriteString("}\n")

	ee.enum2ExtendRendered = true

	return buf.String()
}

func (ee *EnumExtend) Extend2EnumVarName() string {
	return fmt.Sprintf("_%s%sMap", ee.enum.Name, ee.Name)
}

func (ee *EnumExtend) Extend2EnumMap() string {
	if ee.extend2EnumRendered == true {
		return ""
	}

	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(fmt.Sprintf("var %s = map[%s]%s{\n", ee.Extend2EnumVarName(), ee.Type.String(), ee.enum.Name))

	index := 0
	for _, item := range ee.enum.GetItems() {
		var itemValue = ""
		nextIndex := index + len(item.GetName())

		if ee.Type == reflect.String {
			if ee.Name == EnumItemName && ee.enum.Type != reflect.String {
				itemValue = item.GetName()
				buf.WriteString(fmt.Sprintf("	_%sName[%d:%d]: %s,\n", ee.enum.Name, index, nextIndex, item.GetCodeName()))
			} else {
				itemValue = structure.MustConvertTo[string](item.ExtendData[ee.idx])
				buf.WriteString(fmt.Sprintf("	\"%s\": %s,\n", itemValue, item.GetCodeName()))
			}
		} else {
			buf.WriteString(fmt.Sprintf("	%d: %s,\n", structure.MustConvertToKind(item.ExtendData[ee.idx], ee.Type), item.GetCodeName()))
		}
		if ee.enum.Config.NoCase && ee.Type == reflect.String && (itemValue != strings.ToLower(itemValue)) {
			if ee.Name == EnumItemName && ee.enum.Type != reflect.String {
				buf.WriteString(fmt.Sprintf("	strings.ToLower(_%sName[%d:%d]): %s,\n", ee.enum.Name, index, nextIndex, item.GetCodeName()))
			} else {
				buf.WriteString(fmt.Sprintf("	\"%s\": %s,\n", strings.ToLower(itemValue), item.GetCodeName()))
			}
		}

		index = nextIndex
	}
	buf.WriteString("}\n")

	ee.extend2EnumRendered = true

	return buf.String()
}
