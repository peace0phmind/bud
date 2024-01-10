package mvc

import (
	"github.com/peace0phmind/bud/example/mvc/base"
	"github.com/peace0phmind/bud/factory"
)

var _cameraRest = factory.Singleton[CameraRest]().Getter()

type CameraRest struct {
	base.BaseRest[Camera]
}

func (cr *CameraRest) Init() {
	cr.RestName = "camera"
	cr.Repo = _cameraRepo()
}
