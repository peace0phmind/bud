package mvc

import (
	"testing"
)

func TestCameraRest(t *testing.T) {

	c := _cameraRest()
	print(c.Self)
	print(c.RestName)
}
