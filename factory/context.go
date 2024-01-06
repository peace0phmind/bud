package factory

var Context = Singleton[context]().MustBuilder()

type context struct {
}

func (c *context) Get(instanceName string) any {
	return nil
}
