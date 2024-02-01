package enum

import (
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/bud/ast"
	"github.com/peace0phmind/bud/stream"
	"github.com/peace0phmind/bud/structure"
	"github.com/peace0phmind/bud/util"
	goast "go/ast"
	"reflect"
)

const (
	BlankIdentifier = "_"
)

type Enum struct {
	Name    string
	Type    reflect.Kind
	Comment string
	Attrs   []*Attribute
	Items   []*Item
	Config  *Config
}

func (e *Enum) UpdateAttributes(a *ast.Annotation) error {
	if a.Params != nil && len(a.Params.List) > 0 {
		for idx, p := range a.Params.List {
			if p.Value == nil {
				return errors.New(fmt.Sprintf("Enum %s's attribute %s's type is empty", e.Name, p.Key.Text))
			}
			typeName, err := structure.ConvertTo[string](p.Value.Value())
			if err != nil {
				return errors.New(fmt.Sprintf("Enum %s's attribute %s's type parse error: %v", e.Name, p.Key.Text, err))
			}
			t, err := getEnumAttributeKindByName(typeName)
			if err != nil {
				return errors.New(fmt.Sprintf("enum type err: %v", err))
			}

			comment := ast.GetCommentsText(p.Comments)
			if len(comment) == 0 {
				comment = ast.GetCommentText(p.Comment)
			}

			e.Attrs = append(e.Attrs, &Attribute{
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

func (e *Enum) UpdateItems(a *ast.Annotation) error {
	if a.Extends != nil && len(a.Extends.List) > 0 {
		for idx, ex := range a.Extends.List {
			if len(e.Attrs) != len(ex.Values) {
				return errors.New("enum data number not equals with enum attribute type")
			}

			var value any
			if ex.Value != nil {
				value = ex.Value.Value()
			}

			ei := &Item{
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
				ei.AttributeData = stream.Must(stream.Map[ast.Value, any](stream.Of(ex.Values), func(value ast.Value) (any, error) {
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
	// check e.Attribute exist name equals "Name" and type is string
	for _, ex := range e.Attrs {
		if ex.Name == ItemName && ex.Type != reflect.String {
			return errors.New("enum attribute 'Name' must have type string")
		}
	}

	// check e.Attribute name is unique
	attributeNames := make(map[string]bool)
	for _, ex := range e.Attrs {
		if attributeNames[ex.Name] {
			return fmt.Errorf("enum attribute names must be unique, %s", ex.Name)
		}
		attributeNames[ex.Name] = true
	}

	// check e.Items name is unique
	itemNames := make(map[string]bool)
	for _, item := range e.GetItems() {
		if itemNames[item.Name] {
			return fmt.Errorf("enum item names must be unique, %s", item.Name)
		}
		itemNames[item.Name] = true
	}

	// if e.Attribute is empty or e.Attribute haven't a ItemName item, then use item's name to create it
	if fee := e.FindAttributeByName(ItemName); fee == nil {
		for _, ee := range e.Attrs {
			ee.idx += 1
		}

		ee := &Attribute{
			enum:    e,
			idx:     0,
			Name:    ItemName,
			Type:    reflect.String,
			Comment: "",
		}
		e.Attrs = append([]*Attribute{ee}, e.Attrs...)

		for _, ei := range e.GetItems() {
			ei.AttributeData = append([]any{ei.Name}, ei.AttributeData...)
		}
	} else {
		for _, ei := range e.GetItems() {
			if isBlankIdentifier(ei.AttributeData[fee.idx]) {
				ei.AttributeData[fee.idx] = ei.Name
			}
		}
	}

	// check config names and type
	spee := e.FindAttributeByName(e.Config.StringParseName)
	if spee == nil {
		return errors.New("enum config string parse name must exist in enum attributes")
	} else {
		if !stream.Must(stream.Of(enumTypes).Contains(spee.Type, func(x, y reflect.Kind) (bool, error) { return x == y, nil })) {
			return errors.New("StringParseName's type muse be number or string")
		}
	}

	mee := e.FindAttributeByName(e.Config.MarshalName)
	if mee == nil {
		return errors.New("enum config marshal name must exist in enum attributes")
	} else {
		if !stream.Must(stream.Of(enumTypes).Contains(mee.Type, func(x, y reflect.Kind) (bool, error) { return x == y, nil })) {
			return errors.New("MarshalName's type muse be number or string")
		}
	}

	//see := e.FindAttributeByName(e.Config.SqlName)
	//if see == nil {
	//	return errors.New("enum config sql name must exist in enum attributes")
	//} else {
	//	if !stream.Must(stream.Of(enumTypes).Contains(see.Type, func(x, y reflect.Kind) (bool, error) { return x == y, nil })) {
	//		return errors.New("SqlName's type muse be number or string")
	//	}
	//}

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
		if stream.Must(stream.Of(e.GetItems()).AnyMatch(func(item *Item) (bool, error) { return item.Value != nil, nil })) {
			value := 0
			for _, item := range e.GetItems() {
				if item.Value == nil {
					item.Value = value
					value += 1
				} else {
					item.Value = structure.MustConvertTo[int](item.Value)
					value = item.Value.(int) + 1
				}
				item.Value = structure.MustConvertToKind(item.Value, e.Type)
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

func (e *Enum) FindAttributeByName(name string) *Attribute {
	if len(e.Attrs) > 0 {
		for _, ee := range e.Attrs {
			if ee.Name == name {
				return ee
			}
		}
	}

	return nil
}

func (e *Enum) GetItems() []*Item {
	return stream.Must(stream.Of(e.Items).Filter(func(item *Item) (bool, error) {
		return !item.IsBlankIdentifier, nil
	}).ToSlice())
}

func annotationGroupToEnumConfig(ag *ast.AnnotationGroup, globalConfig *Config) (*Config, error) {
	enumConfAnnotation := ag.FindAnnotationByName("EnumConfig")

	if enumConfAnnotation != nil {
		ec, err := ast.AnnotationParamsTo[Config](structure.Clone(globalConfig), enumConfAnnotation)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("pase annotation err: %v", err))
		}
		return ec, nil
	}

	if globalConfig != nil {
		return structure.Clone(globalConfig), nil
	} else {
		return nil, nil
	}
}

func annotationGroupToEnum(ag *ast.AnnotationGroup, ts *goast.TypeSpec, globalConfig *Config) (*Enum, error) {
	enumAnnotation := ag.FindAnnotationByName("enum")
	if enumAnnotation == nil {
		return nil, nil
	}

	ec, err := annotationGroupToEnumConfig(ag, globalConfig)
	if err != nil {
		return nil, err
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

	err = enum.UpdateAttributes(enumAnnotation)
	if err != nil {
		return nil, err
	}

	err = enum.UpdateItems(enumAnnotation)
	if err != nil {
		return nil, err
	}

	err = enum.CheckValid()
	if err != nil {
		return nil, err
	}

	return enum, nil
}
