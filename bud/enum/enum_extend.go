package enum

import (
	"bytes"
	"fmt"
	"github.com/peace0phmind/bud/structure"
	"reflect"
)

type EnumExtend struct {
	enum    *Enum
	idx     int
	Name    string
	Type    reflect.Kind
	Comment string
}

func (ee *EnumExtend) GetExtendValueMap() string {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(fmt.Sprintf("var _%sMap%s = map[%s]%s{\n", ee.enum.Name, ee.Name, ee.enum.Name, ee.Type.String()))
	if ee.idx == 0 {
		index := 0
		for _, item := range ee.enum.Items {
			nextIndex := index + len(item.GetName())
			buf.WriteString(fmt.Sprintf("	%s: _%s%s[%d:%d],\n", item.GetCodeName(), ee.enum.Name, ee.Name, index, nextIndex))
			index = nextIndex
		}
	} else {
		for _, item := range ee.enum.Items {
			switch ee.Type {
			case reflect.String:
				buf.WriteString(fmt.Sprintf("	%s: \"%s\",\n", item.GetCodeName(), structure.MustConvertTo[string](item.ExtendData[ee.idx])))
			default:
				buf.WriteString(fmt.Sprintf("	%s: %s,\n", item.GetCodeName(), structure.MustConvertToKind(item.ExtendData[ee.idx], ee.Type)))
			}
		}
	}
	buf.WriteString("}\n")

	return buf.String()
}
