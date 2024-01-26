package enum

import (
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/factory"
	"github.com/peace0phmind/bud/stream"
	"github.com/peace0phmind/bud/structure"
	goast "go/ast"
	"reflect"
)

type EnumConfig struct {
	NoPrefix    bool   `value:"false"` // 所有生成的枚举不携带类型名称前缀
	Marshal     bool   `value:"true"`
	MarshalName string `value:"Name"`
	Sql         bool   `value:"false"`
	SqlName     string `value:"value"`
	Names       bool   `value:"false"` // enum name list
	Values      bool   `value:"true"`  // enum item list
	NoCase      bool   `value:"true"`  // case insensitivity
	MustParse   bool   `value:"false"`
}

type Enum struct {
	Name    string
	Type    reflect.Kind
	Comment string
	Extends []EnumExtend
	Items   []EnumItem
	Config  *EnumConfig
}

type EnumExtend struct {
	Name    string
	Type    reflect.Kind
	Comment string
}

type EnumItem struct {
	Name       string
	Value      any
	Comment    string
	ExtendData []any
}

func (e *Enum) UpdateExtends(a *ast.Annotation) error {
	if a.Params != nil && len(a.Params.List) > 0 {
		for _, p := range a.Params.List {
			if p.Value == nil {
				return errors.New(fmt.Sprintf("Enum %s's extend field %s's type is empty", e.Name, p.Key.Text))
			}
			typeName, err := structure.ConvertTo[string](p.Value.Value())
			if err != nil {
				return errors.New(fmt.Sprintf("Enum %s's extend field %s's type parse error: %v", e.Name, p.Key.Text, err))
			}
			t, err := getKindByName(typeName)
			if err != nil {
				return errors.New(fmt.Sprintf("enum type err: %v", err))
			}

			comment := ast.GetCommentsText(p.Comments)
			if len(comment) == 0 {
				comment = ast.GetCommentText(p.Comment)
			}

			e.Extends = append(e.Extends, EnumExtend{
				Name:    p.Key.Text,
				Type:    t,
				Comment: comment,
			})
		}
	}

	return nil
}

func (e *Enum) UpdateEnumItem(a *ast.Annotation) error {
	if a.Extends != nil && len(a.Extends.List) > 0 {
		for _, ex := range a.Extends.List {
			if len(e.Extends) != len(ex.Values) {
				return errors.New("enum data number not equals with extend type")
			}

			var value any
			if ex.Value != nil {
				value = ex.Value.Value()
			}
			comment := ast.GetCommentsText(ex.Comments)
			if len(comment) == 0 {
				comment = ast.GetCommentText(ex.Comment)
			}

			ei := EnumItem{
				Name:    ex.Name.Text,
				Value:   value,
				Comment: comment,
			}

			ei.ExtendData = stream.Map[ast.Value, any](stream.Of(ex.Values), func(value ast.Value) (any, error) {
				return value.Value(), nil
			}).MustToSlice()

			e.Items = append(e.Items, ei)
		}
		return nil
	}

	return errors.New("Enum must have some items")
}

func (e *Enum) CheckValid() error {

	return nil
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

	t, err := getKindByName(fmt.Sprintf("%s", ts.Type))
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
