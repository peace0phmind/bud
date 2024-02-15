package opt

import (
	"database/sql"
	"database/sql/driver"
)

type SqlNullInf interface {
	sql.Scanner
	driver.Valuer
}

type SqlNull[T any] struct {
	SqlV  T
	Valid bool
}

func NewSqlNull[T any](t T) *SqlNull[T] {
	return &SqlNull[T]{SqlV: t, Valid: true}
}

func (sn *SqlNull[T]) Scan(value any) error {
	sn.Valid = false

	if value == nil {
		return nil
	}

	if scanner, ok := any(&sn.SqlV).(sql.Scanner); ok {
		err := scanner.Scan(value)
		if err == nil {
			sn.Valid = true
			return nil
		}
		return err
	}

	return nil
}

func (sn *SqlNull[T]) Value() (driver.Value, error) {
	if !sn.Valid {
		return nil, nil
	}

	if valuer, ok := any(sn.SqlV).(driver.Valuer); ok {
		return valuer.Value()
	}

	return nil, nil
}
