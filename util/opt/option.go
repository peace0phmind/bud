package opt

type Option[T any] struct {
	value T
	set   bool
}

func NewOption[T any](t T) *Option[T] {
	return &Option[T]{value: t, set: true}
}

func (o *Option[T]) Get() (t T) {
	if o.set {
		return o.value
	}

	return t
}

func (o *Option[T]) Set(t T) {
	o.set = true
	o.value = t
}

func (o *Option[T]) IsSet() bool {
	return o.set
}
