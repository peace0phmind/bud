package mvc

import (
	"github.com/peace0phmind/bud/example/mvc/base"
	"github.com/peace0phmind/bud/factory"
)

var _cameraRest = factory.Singleton[CameraRest]().MustBuilder()

type CameraRest struct {
	base.BaseRest[Camera]
}

func (cr *CameraRest) MustInitOnce() {
	cr.RestName = "camera"
	cr.Repo = _cameraRepo()
}
