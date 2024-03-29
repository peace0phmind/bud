package enum

import (
	"bytes"
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
				isValue: false,
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
			} else {
				if len(e.Attrs) != len(ex.Values) {
					return errors.New("enum data number not equals with enum attribute type")
				}
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

func (e *Enum) CheckValid() (err error) {
	if err = e.checkAttributeNameUnique(); err != nil {
		return err
	}

	if err = e.checkItemNameUnique(); err != nil {
		return err
	}

	if err = e.checkAndUpdateNameAttribute(); err != nil {
		return err
	}

	if err = e.checkAndUpdateValueAttribute(); err != nil {
		return err
	}

	if err = e.Config.CheckValid(); err != nil {
		return err
	}

	return nil
}

// checkItemNameUnique checks if the item names in the Enum are unique.
// It returns an error if a duplicate name is found.
func (e *Enum) checkItemNameUnique() error {
	itemNames := make(map[string]bool)
	for _, item := range e.GetItems() {
		if itemNames[item.Name] {
			return fmt.Errorf("enum item names must be unique, %s", item.Name)
		}
		itemNames[item.Name] = true
	}
	return nil
}

// checkAttributeNameUnique checks if the attribute names in the Enum are unique.
// It returns an error if a duplicate name is found.
func (e *Enum) checkAttributeNameUnique() error {
	// check e.Attribute name is unique
	attributeNames := make(map[string]bool)
	for _, ex := range e.Attrs {
		if attributeNames[ex.Name] {
			return fmt.Errorf("enum attribute names must be unique, %s", ex.Name)
		}
		attributeNames[ex.Name] = true
	}
	return nil
}

// checkAndUpdateNameAttribute checks if the attribute 'Name' exists in the enum.
// If it doesn't exist, it adds the attribute 'Name' to the enum and updates the item data.
// If the attribute 'Name' exists, it checks if its type is string and updates the item data accordingly.
// It returns an error if the attribute 'Name' exists and its type is not string.
func (e *Enum) checkAndUpdateNameAttribute() error {
	if nameAttr := e.FindAttributeByName(ItemName); nameAttr == nil {
		for _, ee := range e.Attrs {
			ee.idx += 1
		}

		nameAttr = &Attribute{
			enum:    e,
			idx:     0,
			isValue: false,
			Name:    ItemName,
			Type:    reflect.String,
			Comment: "",
		}
		e.Attrs = append([]*Attribute{nameAttr}, e.Attrs...)

		for _, ei := range e.GetItems() {
			ei.AttributeData = append([]any{ei.Name}, ei.AttributeData...)
		}
	} else {
		if nameAttr.Type != reflect.String {
			return errors.New("enum attribute 'Name' must have type string")
		}

		for _, ei := range e.GetItems() {
			if isBlankIdentifier(ei.AttributeData[nameAttr.idx]) {
				ei.AttributeData[nameAttr.idx] = ei.Name
			}
		}
	}

	return nil
}

func (e *Enum) checkAndUpdateValueAttribute() error {
	if valueAttr := e.FindAttributeByName(ItemValue); valueAttr != nil {
		return errors.New("\"Value\" is a reserved attribute in enum and cannot appear in named parameters. However, it can be directly specified after \"=\".")
	}

	valueAttr := &Attribute{
		enum:    e,
		idx:     -1,
		isValue: true,
		Name:    ItemValue,
		Type:    e.Type,
		Comment: "",
	}
	e.Attrs = append(e.Attrs, valueAttr)

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
			nextValue := 0
			for _, item := range e.Items {
				if item.Value == nil {
					item.Value = nextValue
					nextValue += 1
				} else {
					item.Value = structure.MustConvertTo[int](item.Value)
					nextValue = item.Value.(int) + 1
				}
				item.Value = structure.MustConvertToKind(item.Value, e.Type)
			}
		}
	}

	return nil
}

func (e *Enum) Names() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("var _%sNames = []string{\n", e.Name))
	index := 0
	for _, item := range e.GetItems() {
		if e.Type == reflect.String {
			buf.WriteString(fmt.Sprintf("\t\"%s\",\n", item.GetName()))
		} else {
			nextIndex := index + len(item.GetName())
			buf.WriteString(fmt.Sprintf("\t_%sName[%d:%d],\n", e.Name, index, nextIndex))
			index = nextIndex
		}
	}
	buf.WriteString("}")
	return buf.String()
}

func (e *Enum) EmptyEnumValue() string {
	if e.Type == reflect.String {
		return "\"\""
	} else {
		return fmt.Sprintf("%s(0)", e.Name)
	}
}

func isBlankIdentifier(value any) bool {
	if bi, ok := value.(string); ok {
		return bi == BlankIdentifier
	}
	return false
}

func (e *Enum) FindAttributeByName(name string) *Attribute {
	capName := util.Capitalize(name)
	if len(e.Attrs) > 0 {
		for _, ee := range e.Attrs {
			if ee.Name == capName {
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
	ec.enum = enum

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
