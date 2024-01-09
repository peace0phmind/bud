package base

import (
	"gorm.io/gorm"
)

/*
PageParams

	page: 从1开始表示第一页
		如果page等于0，返回第一页；如果page等于1，还是返回第一页
		如果page小于0,则不进行分页
	size: 表示分页大小
*/
type PageParams struct {
	Page      int `json:"page"`
	Size      int `json:"size"`
	TotalPage int `json:"totalPage"`
	Total     int `json:"total"`
}

func (pp *PageParams) CheckDefault() {
	if pp.Size <= 0 {
		pp.Size = 10
	}

	// 分页从1开始，如果传递空分页对象则初始化为页码1
	if pp.Page == 0 {
		pp.Page = 1
	}

	if pp.Page < 0 {
		pp.Size = pp.Total
		pp.TotalPage = 1
	} else {
		pp.TotalPage = (pp.Total + pp.Size - 1) / pp.Size
	}
}

func (pp *PageParams) QueryCondition(db *gorm.DB) *gorm.DB {
	if db == nil {
		return nil
	}

	if pp.Page < 0 {
		return db
	} else {
		return db.Offset((pp.Page - 1) * pp.Size).Limit(pp.Size)
	}
}
