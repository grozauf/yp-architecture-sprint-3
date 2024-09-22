package devices

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
)

var types = []string{"switcher", "gates", "heater"}

type DevicesModule struct {
}

func New() *DevicesModule {
	return &DevicesModule{}
}

func (d *DevicesModule) GetDeviceInfo(ctx context.Context, moduleId, deviceId int) (DeviceInfo, error) {
	return DeviceInfo{
		SerialNumber: uuid.NewString(),
		Status: func() bool {
			return rand.Intn(2) == 0
		}(),
		Type: types[rand.Intn(len(types))],
	}, nil

}

func (d *DevicesModule) SetDeviceTargetValue(ctx context.Context, moduleId, deviceId int, value float32) error {
	return nil
}

func (d *DevicesModule) SetDeviceStatus(ctx context.Context, moduleId, deviceId int, status bool) error {
	return nil
}
