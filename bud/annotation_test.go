package bud

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func toPtr(v Value) *Value {
	return &v
}

func TestParseAnnotation(t *testing.T) {
	//comment := "// test"
	singleComment := "// comment"
	multiComment := "/* comment */"

	tests := []struct {
		name     string
		fileName string
		text     string
		want     *AnnotationGroup
		wantErr  bool
	}{
		{
			name:     "Empty file name and text",
			fileName: "",
			text:     "",
			want:     &AnnotationGroup{},
			wantErr:  false,
		},
		{
			name:     "Valid file name but empty text",
			fileName: "file.go",
			text:     "",
			want:     &AnnotationGroup{},
			wantErr:  false,
		},
		{
			name:     "only one annotation name",
			fileName: "file.go",
			text:     "@tag",
			want: &AnnotationGroup{
				Annotations: []*Annotation{
					{
						Name: Name{Text: "tag"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "two annotation name",
			fileName: "file.go",
			text: `
@tag
@sql
`,
			want: &AnnotationGroup{
				Annotations: []*Annotation{
					{
						Name: Name{Text: "tag"},
					},
					{
						Name: Name{Text: "sql"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "two annotation with params",
			fileName: "file.go",
			text: `
@tag(disable, string = "str\"ing" , int=123, double=456.7, bool = true)
@sql(code int32, name string, message=string)
`,
			want: &AnnotationGroup{
				Annotations: []*Annotation{
					{
						Name: Name{Text: "tag"},
						Params: &Params{List: []*AnnotationParam{
							{
								Key: Key{Text: "disable"},
							},
							{
								Key:   Key{Text: "string"},
								Value: toPtr(any(String{Value: "str\"ing"}).(Value)),
							},
							{
								Key:   Key{Text: "int"},
								Value: toPtr(any(Int{Value: 123}).(Value)),
							},
							{
								Key:   Key{Text: "double"},
								Value: toPtr(any(Float{Value: 456.7}).(Value)),
							},
							{
								Key:   Key{Text: "bool"},
								Value: toPtr(any(Bool{Value: true}).(Value)),
							},
						},
						},
					},
					{
						Name: Name{Text: "sql"},
						Params: &Params{List: []*AnnotationParam{
							{
								Key:   Key{Text: "code"},
								Value: toPtr(any(String{Value: "int32"}).(Value)),
							},
							{
								Key:   Key{Text: "name"},
								Value: toPtr(any(String{Value: "string"}).(Value)),
							},
							{
								Key:   Key{Text: "message"},
								Value: toPtr(any(String{Value: "string"}).(Value)),
							},
						},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "two annotation with params and extends",
			fileName: "file.go",
			text: `
@tag(disable, string = "str\"ing" , int=123, double=456.7, bool = true) {
	Good  
    GoodWithIntValue = 12 
    GoodWithStrValue = "str" 
    GoodWithParams ("string", 123, 456.7, true )  // comment
    GoodWithAll ("string", 123, 456.7, false ) = 89 /* comment */
}
@sql(code int32, name string, message=string)
`,
			want: &AnnotationGroup{
				Annotations: []*Annotation{
					{
						Name: Name{Text: "tag"},
						Params: &Params{List: []*AnnotationParam{
							{
								Key: Key{Text: "disable"},
							},
							{
								Key:   Key{Text: "string"},
								Value: toPtr(any(String{Value: "str\"ing"}).(Value)),
							},
							{
								Key:   Key{Text: "int"},
								Value: toPtr(any(Int{Value: 123}).(Value)),
							},
							{
								Key:   Key{Text: "double"},
								Value: toPtr(any(Float{Value: 456.7}).(Value)),
							},
							{
								Key:   Key{Text: "bool"},
								Value: toPtr(any(Bool{Value: true}).(Value)),
							},
						},
						},
						Extends: &Extends{List: []*AnnotationExtend{
							{
								Name: Name{Text: "Good"},
							},
							{
								Name:  Name{Text: "GoodWithIntValue"},
								Value: toPtr(any(Int{Value: 12}).(Value)),
							},
							{
								Name:  Name{Text: "GoodWithStrValue"},
								Value: toPtr(any(String{Value: "str"}).(Value)),
							},
							{
								Name: Name{Text: "GoodWithParams"},
								Values: []Value{
									any(String{Value: "string"}).(Value),
									any(Int{Value: 123}).(Value),
									any(Float{Value: 456.7}).(Value),
									any(Bool{Value: true}).(Value),
								},
								Comment: &Comment{Text: singleComment},
							},
							{
								Name: Name{Text: "GoodWithAll"},
								Values: []Value{
									any(String{Value: "string"}).(Value),
									any(Int{Value: 123}).(Value),
									any(Float{Value: 456.7}).(Value),
									any(Bool{Value: false}).(Value),
								},
								Value:   toPtr(any(Int{Value: 89}).(Value)),
								Comment: &Comment{Text: multiComment},
							},
						},
						},
					},
					{
						Name: Name{Text: "sql"},
						Params: &Params{List: []*AnnotationParam{
							{
								Key:   Key{Text: "code"},
								Value: toPtr(any(String{Value: "int32"}).(Value)),
							},
							{
								Key:   Key{Text: "name"},
								Value: toPtr(any(String{Value: "string"}).(Value)),
							},
							{
								Key:   Key{Text: "message"},
								Value: toPtr(any(String{Value: "string"}).(Value)),
							},
						},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "two annotation with params and extends on multiple lines and comments",
			fileName: "file.go",
			text: `

// tag comment 1
// tag comment 2
/* tag comment 3
tag comment 4
*/
@tag(
	// upper disable comment
	disable, // disable comment
    // string comment 1
    // string comment 2
	string = "str\"ing" ,
    /* int comment 1
       int comment 2
    */
	int=123
    double=456.7
    bool = true) {
	// comment 1
    // comment 2
	Good  
	/* comment 1
       comment 2
       comment 3 */
    GoodWithIntValue = 12 
    GoodWithStrValue = "str" 
    GoodWithParams ("string", 123, 456.7, true )  // comment
    GoodWithAll ("string", 123, 456.7, false ) = 89 /* comment */
} // tag inline comment
// sql comment 0
/* sql comment 1
 sql comment 2
*/
@sql(code int32, name string, message=string) // sql inline comment
`,
			want: &AnnotationGroup{
				Annotations: []*Annotation{
					{
						Comments: []*Comment{
							{
								Text: "// tag comment 1",
							},
							{
								Text: "// tag comment 2",
							},
							{
								Text: `/* tag comment 3
tag comment 4
*/`,
							},
						},
						Name: Name{Text: "tag"},
						Params: &Params{List: []*AnnotationParam{
							{
								Comments: []*Comment{
									{Text: "// upper disable comment"},
								},
								Key:     Key{Text: "disable"},
								Comment: &Comment{Text: "// disable comment"},
							},
							{
								Comments: []*Comment{
									{Text: "// string comment 1"},
									{Text: "// string comment 2"},
								},
								Key:   Key{Text: "string"},
								Value: toPtr(any(String{Value: "str\"ing"}).(Value)),
							},
							{
								Comments: []*Comment{
									{Text: `/* int comment 1
       int comment 2
    */`},
								},
								Key:   Key{Text: "int"},
								Value: toPtr(any(Int{Value: 123}).(Value)),
							},
							{
								Key:   Key{Text: "double"},
								Value: toPtr(any(Float{Value: 456.7}).(Value)),
							},
							{
								Key:   Key{Text: "bool"},
								Value: toPtr(any(Bool{Value: true}).(Value)),
							},
						},
						},
						Extends: &Extends{List: []*AnnotationExtend{
							{
								Comments: []*Comment{
									{Text: "// comment 1"},
									{Text: "// comment 2"},
								},
								Name: Name{Text: "Good"},
							},
							{
								Comments: []*Comment{
									{Text: `/* comment 1
       comment 2
       comment 3 */`},
								},
								Name:  Name{Text: "GoodWithIntValue"},
								Value: toPtr(any(Int{Value: 12}).(Value)),
							},
							{
								Name:  Name{Text: "GoodWithStrValue"},
								Value: toPtr(any(String{Value: "str"}).(Value)),
							},
							{
								Name: Name{Text: "GoodWithParams"},
								Values: []Value{
									any(String{Value: "string"}).(Value),
									any(Int{Value: 123}).(Value),
									any(Float{Value: 456.7}).(Value),
									any(Bool{Value: true}).(Value),
								},
								Comment: &Comment{Text: singleComment},
							},
							{
								Name: Name{Text: "GoodWithAll"},
								Values: []Value{
									any(String{Value: "string"}).(Value),
									any(Int{Value: 123}).(Value),
									any(Float{Value: 456.7}).(Value),
									any(Bool{Value: false}).(Value),
								},
								Value:   toPtr(any(Int{Value: 89}).(Value)),
								Comment: &Comment{Text: multiComment},
							},
						},
						},
						Comment: &Comment{Text: "// tag inline comment"},
					},
					{
						Comments: []*Comment{
							{
								Text: "// sql comment 0",
							},
							{
								Text: `/* sql comment 1
 sql comment 2
*/`,
							},
						},
						Name: Name{Text: "sql"},
						Params: &Params{List: []*AnnotationParam{
							{
								Key:   Key{Text: "code"},
								Value: toPtr(any(String{Value: "int32"}).(Value)),
							},
							{
								Key:   Key{Text: "name"},
								Value: toPtr(any(String{Value: "string"}).(Value)),
							},
							{
								Key:   Key{Text: "message"},
								Value: toPtr(any(String{Value: "string"}).(Value)),
							},
						},
						},
						Comment: &Comment{
							Text: "// sql inline comment",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAnnotation(tt.fileName, tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAnnotation() error = %v, wantErr %v", err, tt.wantErr)
			}

			opt := cmp.FilterPath(ignorePosFields, cmp.Ignore())
			diff := cmp.Diff(tt.want, got, opt)
			if len(diff) > 0 {
				t.Errorf("ParseAnnotation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func ignorePosFields(path cmp.Path) bool {
	// 遍历路径中的每个步骤
	for _, step := range path {
		if t, ok := step.(cmp.StructField); ok {
			// 如果步骤是结构体字段并且名字为"Pos"，则返回true以忽略
			if t.Name() == "Pos" {
				return true
			}
		}
	}
	// 对于其他字段，不忽略
	return false
}
