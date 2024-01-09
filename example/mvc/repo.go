package mvc

import (
	"github.com/peace0phmind/bud/example/mvc/base"
	"github.com/peace0phmind/bud/factory"
)

var _cameraRepo = factory.Singleton[CameraRepo]().MustBuilder()

type CameraRepo struct {
	base.BaseRepo[Camera]
}
