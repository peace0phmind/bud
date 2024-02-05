package enum

//go:generate go run ../../../main.go

// @EnumConfig(sql, ptr, marshal, nocomments)
// @ENUM{pending, inWork, completed, rejected}
type ProjectStatus int
