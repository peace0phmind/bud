package opt

import (
	"database/sql"
	"database/sql/driver"
)

type SqlNullInf interface {
	sql.Scanner
	driver.Valuer
}

type SqlNull[T SqlNullInf] struct {
	SqlV  T
	Valid bool
}

func NewSqlNull[T SqlNullInf](t T) *SqlNull[T] {
	return &SqlNull[T]{SqlV: t}
}

func (sn *SqlNull[T]) Scan(value any) error {
	sn.Valid = false

	if value == nil {
		return nil
	}

	err := sn.SqlV.Scan(value)
	if err == nil {
		sn.Valid = true
		return nil
	}

	return err
}

func (sn *SqlNull[T]) Value() (driver.Value, error) {
	if !sn.Valid {
		return nil, nil
	}

	return sn.SqlV.Value()
}
