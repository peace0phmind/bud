package factory

import (
	"github.com/peace0phmind/bud/factory"
	"github.com/stretchr/testify/assert"
	"testing"
)

type animal struct {
	Name string
}

func (a *animal) MustInitOnce() {
	a.Name = "animal"
}

type cat struct {
	animal // 嵌入 animal 结构体
	Name   string
}

func (c *cat) MustInitOnce() {
	c.animal.MustInitOnce()

	c.Name = "cat"
}

type dog struct {
	animal // 嵌入 animal 结构体
}

func (d *dog) MustInitOnce() {
	d.animal.MustInitOnce()

	d.Name = "dog"
}

func (c *cat) Meow() string {
	return "meow"
}

func TestSingletonBuilder(t *testing.T) {
	var Cat = factory.Singleton[cat]().MustBuilder()

	assert.Equal(t, "meow", Cat().Meow())
	assert.Equal(t, "cat", Cat().Name)
	assert.Equal(t, "animal", Cat().animal.Name)

	var Dog = factory.Singleton[dog]().MustBuilder()

	assert.Equal(t, "dog", Dog().Name)
	assert.Equal(t, "dog", Dog().animal.Name)
}
