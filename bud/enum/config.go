package enum

type Config struct {
	Prefix           string
	NoPrefix         bool   `value:"false"` // 所有生成的枚举不携带类型名称前缀
	StringParse      bool   `value:"true"`
	StringParseName  string `value:"Name"`
	MustParse        bool   `value:"false"`
	Marshal          bool   `value:"false"`
	MarshalName      string `value:"Name"`
	Sql              bool   `value:"false"`
	SqlName          string `value:"Value"`
	Values           bool   `value:"false"` // enum item list
	NoCase           bool   `value:"false"` // case insensitivity
	UseCamelCaseName bool   `value:"true"`
	NoComments       bool   `value:"false"`
	Ptr              bool   `value:"false"`
	ForceUpper       bool   `value:"false"`
	ForceLower       bool   `value:"false"`
}

func (ec *Config) SetForceLower(lower bool) {
	if lower {
		if ec.ForceUpper {
			ec.ForceUpper = false
		}
	}
	ec.ForceLower = lower
}

func (ec *Config) SetForceUpper(upper bool) {
	if upper {
		if ec.ForceLower {
			ec.ForceLower = false
		}
	}
	ec.ForceUpper = upper
}
