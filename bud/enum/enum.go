package enum

import (
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/factory"
	"github.com/peace0phmind/bud/stream"
	"github.com/peace0phmind/bud/structure"
	"github.com/peace0phmind/bud/util"
	goast "go/ast"
	"reflect"
)

const (
	BlankIdentifier = "_"
)

type EnumConfig struct {
	NoPrefix         bool `value:"false"` // 所有生成的枚举不携带类型名称前缀
	Prefix           string
	StringParse      bool   `value:"true"`
	StringParseName  string `value:"Name"`
	MustParse        bool   `value:"false"`
	Marshal          bool   `value:"false"`
	MarshalName      string `value:"Name"`
	Sql              bool   `value:"false"`
	SqlName          string `value:"Value"`
	Names            bool   `value:"false"` // enum name list
	Values           bool   `value:"false"` // enum item list
	NoCase           bool   `value:"false"` // case insensitivity
	UseCamelCaseName bool   `value:"true"`
	NoComments       bool   `value:"false"`
	Ptr              bool   `value:"false"`
	ForceLower       bool   `value:"false"`
}

type Enum struct {
	Name    string
	Type    reflect.Kind
	Comment string
	Extends []*EnumExtend
	Items   []*EnumItem
	Config  *EnumConfig
}

func (e *Enum) UpdateExtends(a *ast.Annotation) error {
	if a.Params != nil && len(a.Params.List) > 0 {
		for idx, p := range a.Params.List {
			if p.Value == nil {
				return errors.New(fmt.Sprintf("Enum %s's extend field %s's type is empty", e.Name, p.Key.Text))
			}
			typeName, err := structure.ConvertTo[string](p.Value.Value())
			if err != nil {
				return errors.New(fmt.Sprintf("Enum %s's extend field %s's type parse error: %v", e.Name, p.Key.Text, err))
			}
			t, err := getEnumExtendKindByName(typeName)
			if err != nil {
				return errors.New(fmt.Sprintf("enum type err: %v", err))
			}

			comment := ast.GetCommentsText(p.Comments)
			if len(comment) == 0 {
				comment = ast.GetCommentText(p.Comment)
			}

			e.Extends = append(e.Extends, &EnumExtend{
				enum:    e,
				idx:     idx,
				Name:    util.Capitalize(p.Key.Text),
				Type:    t,
				Comment: comment,
			})
		}
	}

	return nil
}

func (e *Enum) UpdateEnumItem(a *ast.Annotation) error {
	if a.Extends != nil && len(a.Extends.List) > 0 {
		for idx, ex := range a.Extends.List {
			if len(e.Extends) != len(ex.Values) {
				return errors.New("enum data number not equals with extend type")
			}

			var value any
			if ex.Value != nil {
				value = ex.Value.Value()
			}

			ei := &EnumItem{
				enum:        e,
				idx:         idx,
				Name:        ex.Name.Text,
				Value:       value,
				DocComment:  ast.GetCommentsText(ex.Comments),
				LineComment: ast.GetCommentText(ex.Comment),
			}

			if ei.Name == BlankIdentifier {
				ei.IsBlankIdentifier = true
			}

			if !ei.IsBlankIdentifier {
				ei.ExtendData = stream.Must(stream.Map[ast.Value, any](stream.Of(ex.Values), func(value ast.Value) (any, error) {
					return value.Value(), nil
				}).ToSlice())
			}

			e.Items = append(e.Items, ei)
		}
		return nil
	}

	return errors.New("Enum must have some items")
}

func (e *Enum) CheckValid() error {
	// check e.Extend exist name equals "Name" and type is string
	for _, ex := range e.Extends {
		if ex.Name == EnumItemName && ex.Type != reflect.String {
			return errors.New("enum extend field 'Name' must have type string")
		}
	}

	// check e.Extend name is unique
	extendNames := make(map[string]bool)
	for _, ex := range e.Extends {
		if extendNames[ex.Name] {
			return fmt.Errorf("enum extend field names must be unique, %s", ex.Name)
		}
		extendNames[ex.Name] = true
	}

	// check e.Items name is unique
	itemNames := make(map[string]bool)
	for _, item := range e.GetItems() {
		if itemNames[item.Name] {
			return fmt.Errorf("enum item names must be unique, %s", item.Name)
		}
		itemNames[item.Name] = true
	}

	// if e.Extend is empty or e.Extend haven't a EnumItemName item, then use item's name to create it
	if fee := e.FindExtendByName(EnumItemName); fee == nil {
		for _, ee := range e.Extends {
			ee.idx += 1
		}

		ee := &EnumExtend{
			enum:    e,
			idx:     0,
			Name:    EnumItemName,
			Type:    reflect.String,
			Comment: "",
		}
		e.Extends = append([]*EnumExtend{ee}, e.Extends...)

		for _, ei := range e.GetItems() {
			ei.ExtendData = append([]any{ei.Name}, ei.ExtendData...)
		}
	} else {
		for _, ei := range e.GetItems() {
			if isBlankIdentifier(ei.ExtendData[fee.idx]) {
				ei.ExtendData[fee.idx] = ei.Name
			}
		}
	}

	// check and set item value
	if e.Type == reflect.String {
		for _, ei := range e.GetItems() {
			if ei.Value == nil {
				ei.Value = ei.GetName()
			} else {
				ei.Value = structure.MustConvertTo[string](ei.Value)
			}
		}
	} else {
		if stream.Must(stream.Of(e.GetItems()).AnyMatch(func(item *EnumItem) (bool, error) { return item.Value != nil, nil })) {
			value := 0
			for _, item := range e.GetItems() {
				if item.Value == nil {
					item.Value = value
					value += 1
				} else {
					item.Value = structure.MustConvertTo[int](item.Value)
					value = item.Value.(int) + 1
				}
			}
		}
	}

	return nil
}

func isBlankIdentifier(value any) bool {
	if bi, ok := value.(string); ok {
		return bi == BlankIdentifier
	}
	return false
}

func (e *Enum) FindExtendByName(name string) *EnumExtend {
	if len(e.Extends) > 0 {
		for _, ee := range e.Extends {
			if ee.Name == name {
				return ee
			}
		}
	}

	return nil
}

func (e *Enum) GetItems() []*EnumItem {
	return stream.Must(stream.Of(e.Items).Filter(func(item *EnumItem) (bool, error) {
		return !item.IsBlankIdentifier, nil
	}).ToSlice())
}

func annotationGroupToEnum(ag *ast.AnnotationGroup, ts *goast.TypeSpec) (*Enum, error) {
	enumAnnotation := ag.FindAnnotationByName("enum")
	enumConfAnnotation := ag.FindAnnotationByName("EnumConfig")

	if enumAnnotation == nil {
		return nil, nil
	}

	var ec *EnumConfig = nil
	var err error
	if enumConfAnnotation != nil {
		ec, err = ast.AnnotationParamsTo[EnumConfig](enumConfAnnotation)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("pase annotation err: %v", err))
		}
	} else {
		ec = factory.New[EnumConfig]()
	}

	enum := &Enum{
		Config: ec,
	}

	t, err := getEnumKindByName(fmt.Sprintf("%s", ts.Type))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("enum type err: %v", err))
	}

	enum.Name = ts.Name.Name
	enum.Type = t

	enum.Comment = ast.GetCommentsText(enumAnnotation.Comments)
	if len(enum.Comment) == 0 {
		enum.Comment = ast.GetCommentText(enumAnnotation.Comment)
	}

	err = enum.UpdateExtends(enumAnnotation)
	if err != nil {
		return nil, err
	}

	err = enum.UpdateEnumItem(enumAnnotation)
	if err != nil {
		return nil, err
	}

	err = enum.CheckValid()
	if err != nil {
		return nil, err
	}

	return enum, nil
}
