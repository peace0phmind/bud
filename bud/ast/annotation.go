package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/peace0phmind/bud/factory"
	"github.com/peace0phmind/bud/structure"
	"reflect"
	"strings"
	"text/scanner"
)

type Key struct {
	Pos  lexer.Position
	Text string `@Ident "="?`
}

type Name struct {
	Pos  lexer.Position
	Text string `@Ident`
}

type Comment struct {
	Pos  lexer.Position
	Text string `@Comment`
}

type AnnotationGroup struct {
	Annotations []*Annotation `@@*`
}

type ClosedParenthesis struct {
	Pos               lexer.Position
	ClosedParenthesis string `")"`
}

type Params struct {
	List              []*AnnotationParam `"(" @@*`
	ClosedParenthesis ClosedParenthesis  `@@`
}

type ClosedBracket struct {
	Pos           lexer.Position
	ClosedBracket string `"}"`
}

type Extends struct {
	List          []*AnnotationExtend `"{" @@*`
	ClosedBracket ClosedBracket       `@@`
}

type Annotation struct {
	BeforeUseless *string    `(~(Comment | "@"))*`
	Comments      []*Comment `@@*`
	Name          Name       `"@" @@`
	Params        *Params    `@@?`
	Extends       *Extends   `@@?`
	Comment       *Comment   `@@?`
	AfterUseless  *string    `(~(Comment | "@"))*`
}

type AnnotationParam struct {
	Pos      lexer.Position
	Comments []*Comment `@@*`
	Key      Key        `@@`
	Value    *Value     `@@? ","?`
	Comment  *Comment   `@@?`
}

type AnnotationExtend struct {
	Pos      lexer.Position
	Comments []*Comment `@@*`
	Name     Name       `@@`
	Values   []Value    `("(" @@* ")")?`
	Value    *Value     `("=" @@)?`
	Comment  *Comment   `@@?`
}

type Value interface{ value() any }

type Float struct {
	Value float64 `@Float ","? `
}

func (f Float) value() any { return f.Value }

type Int struct {
	Value int `@Int ","? `
}

func (i Int) value() any {
	return i.Value
}

type String struct {
	Value string `@(String | Ident) ","? `
}

func (s String) value() any {
	return s.Value
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type Bool struct {
	Value Boolean `@("true" | "false") ","? `
}

func (b Bool) value() any {
	return bool(b.Value)
}

//type Unknown struct {
//	Value string `@Ident ","? `
//}
//
//func (u Unknown) value() {}

var annotationParser = participle.MustBuild[AnnotationGroup](
	participle.Lexer(lexer.NewTextScannerLexer(func(s *scanner.Scanner) {
		s.Mode &^= scanner.SkipComments
	})),
	participle.Union[Value](Bool{}, Float{}, Int{}, String{}),
	participle.Unquote("String"),
)

func fixComments(annotationGroup *AnnotationGroup, err error) (*AnnotationGroup, error) {
	if err != nil {
		return annotationGroup, err
	}

	for ai, annotation := range annotationGroup.Annotations {
		if annotation.Params != nil {
			for pi, param := range annotation.Params.List {
				if param.Comment != nil &&
					param.Comment.Pos.Line != param.Key.Pos.Line &&
					pi+1 < len(annotation.Params.List) {
					annotation.Params.List[pi+1].Comments = append([]*Comment{param.Comment}, annotation.Params.List[pi+1].Comments...)
					param.Comment = nil
				}
			}
		}

		if annotation.Extends != nil {
			for ei, extend := range annotation.Extends.List {
				if extend.Comment != nil &&
					extend.Comment.Pos.Line != extend.Name.Pos.Line &&
					ei+1 < len(annotation.Extends.List) {
					annotation.Extends.List[ei+1].Comments = append([]*Comment{extend.Comment}, annotation.Extends.List[ei+1].Comments...)
					extend.Comment = nil
				}
			}
		}

		if annotation.Comment != nil &&
			!(annotation.Comment.Pos.Line == annotation.Name.Pos.Line ||
				(annotation.Params != nil && annotation.Params.ClosedParenthesis.Pos.Line == annotation.Comment.Pos.Line) ||
				(annotation.Extends != nil && annotation.Extends.ClosedBracket.Pos.Line == annotation.Comment.Pos.Line)) &&
			ai+1 < len(annotationGroup.Annotations) {
			annotationGroup.Annotations[ai+1].Comments = append([]*Comment{annotation.Comment}, annotationGroup.Annotations[ai+1].Comments...)
			annotation.Comment = nil
		}
	}

	return annotationGroup, err
}

var defaultBoolValue = any(Bool{Value: true}).(Value)

func AnnotationParamsTo[T any](a *Annotation) (t *T, err error) {

	t = factory.New[T]()

	if a.Params != nil {
		err = structure.WalkField(t, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type) error {
			switch fieldValue.Kind() {
			case reflect.Ptr, reflect.Struct:
				return nil
			}

			var ap *AnnotationParam = nil

			for _, p := range a.Params.List {
				if strings.EqualFold(structField.Name, p.Key.Text) {
					ap = p
					break
				}
			}

			if ap == nil {
				return nil
			}

			if fieldValue.Kind() == reflect.Bool && ap.Value == nil {
				ap.Value = &defaultBoolValue
			}

			if ap.Value != nil {
				return structure.MapToValue((*ap.Value).value(), fieldValue)
			}

			return nil
		})
	}

	return
}

func ParseAnnotation(fileName string, text string) (*AnnotationGroup, error) {
	return fixComments(annotationParser.ParseString(fileName, text))
}
