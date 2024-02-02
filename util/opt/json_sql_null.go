package opt

type JsonSqlNull[T SqlNullInf] struct {
	JsonNull[T]
	SqlNull[T]
}

func NewJsonSqlNull[T SqlNullInf](t T) *JsonSqlNull[T] {
	result := &JsonSqlNull[T]{}
	result.JsonV = t
	result.SqlV = t
	return result
}
