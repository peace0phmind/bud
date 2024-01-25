package enum

//go:generate go-enum --marshal --values --nocomments --nocase

// AnnotationType
// ENUM
//
//	enum="@Enum"
//
// )

type Annotation string

const AnnotationEnum Annotation = "@enum"

type EnumConfig struct {
	NoPrefix    bool   `value:"false"` // 所有生成的枚举不携带类型名称前缀
	Marshal     bool   `value:"true"`
	MarshalName string `value:"Name"`
	Sql         bool   `value:"false"`
	SqlName     string `value:"Value"`
	Names       bool   `value:"false"` // enum name list
	Values      bool   `value:"true"`  // enum item list
	NoCase      bool   `value:"true"`  // case insensitivity
	MustParse   bool   `value:"false"`
}
