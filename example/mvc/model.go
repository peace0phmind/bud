package mvc

import (
	"github.com/peace0phmind/bud/example/mvc/base"
)

type Camera struct {
	base.UUIDBase
	Number  int    `json:"number" gorm:"uniqueIndex:idx_store_number;column:number;not null"` // 编号
	RtspURL string `json:"rtspUrl" gorm:"column:rtsp_url;type:varchar(256)"`                  // rtsp url
	Remark  string `json:"remark" gorm:"column:remark;type:varchar(1024)"`                    // 备注
}

func (*Camera) TableName() string {
	return "camera"
}
