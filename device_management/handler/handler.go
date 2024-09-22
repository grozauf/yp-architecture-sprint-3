package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	conv "management/converter"
	"management/devices"
	"management/swagger"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получить информацию по устройству
	// (GET /manage/{module_id}/{device_id}/info)
	DeviceInfo(c *gin.Context, moduleId int, deviceId int)
	// Установить статус или целевое значение устройства
	// (POST /manage/{module_id}/{device_id}/set)
	DeviceSetValue(c *gin.Context, moduleId int, deviceId int)
}

type DevicesModuleService interface {
	GetDeviceInfo(ctx context.Context, moduleId, deviceId int) (devices.DeviceInfo, error)
	SetDeviceTargetValue(ctx context.Context, moduleId, deviceId int, value float32) error
	SetDeviceStatus(ctx context.Context, moduleId, deviceId int, status bool) error
}

type handler struct {
	moduleSrv DevicesModuleService
}

func New(moduleSrv DevicesModuleService) *handler {
	return &handler{
		moduleSrv: moduleSrv,
	}

}

func (h *handler) DeviceInfo(c *gin.Context, moduleId int, deviceId int) {
	info, err := h.moduleSrv.GetDeviceInfo(c.Request.Context(), moduleId, deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, conv.ToSwaggerDeivceInfo(info))
}

func (h *handler) DeviceSetValue(c *gin.Context, moduleId int, deviceId int) {
	defer c.Request.Body.Close()

	var deviceValue swagger.DeviceValue
	err := c.BindJSON(&deviceValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	switch deviceValue.ValueName {
	case swagger.TargetValue:
		if deviceValue.TargetValue == nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		err = h.moduleSrv.SetDeviceTargetValue(c.Request.Context(), moduleId, deviceId, *deviceValue.TargetValue)
	case swagger.Status:
		if deviceValue.Status == nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		err = h.moduleSrv.SetDeviceStatus(c.Request.Context(), moduleId, deviceId, *deviceValue.Status)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, deviceValue)
}
