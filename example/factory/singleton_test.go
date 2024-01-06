package factory

import (
	"github.com/peace0phmind/bud/factory"
	"github.com/stretchr/testify/assert"
	"testing"
)

type cat struct {
}

func (c *cat) Meow() string {
	return "meow"
}

func TestSingletonBuilder(t *testing.T) {
	var Cat = factory.Singleton[cat]().MustBuilder()
	assert.Equal(t, "meow", Cat().Meow())
}
