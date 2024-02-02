package opt

import (
	"database/sql/driver"
	"testing"
)

type JsonSqlNullTest struct {
}

func (jsn *JsonSqlNullTest) Scan(value interface{}) (err error) {
	return nil
}

func (jsn JsonSqlNullTest) Value() (driver.Value, error) {
	return nil, nil
}

func TestJsonSqlNull(t *testing.T) {
	var jsn *JsonSqlNull[*JsonSqlNullTest]
	jsn = NewJsonSqlNull(&JsonSqlNullTest{})
	jsn.Scan("1")
}
