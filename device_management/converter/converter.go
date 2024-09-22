package converter

import (
	"github.com/AlekSi/pointer"

	"management/devices"
	"management/swagger"
)

func ToSwaggerDeivceInfo(info devices.DeviceInfo) swagger.DeviceInfo {
	return swagger.DeviceInfo{
		SerialNumber: info.SerialNumber,
		Status:       info.Status,
		Type:         info.Type,
	}
}

func ErrorToStringOrNil(err error) *string {
	if err == nil {
		return nil
	}
	return pointer.ToString(err.Error())
}
