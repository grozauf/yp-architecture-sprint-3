package handler

import (
	"context"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"

	conv "telemetry/converter"
	"telemetry/repository"
	"telemetry/swagger"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получение данных телеметрии
	// (GET /devices/{module_id}/{device_id}/telemetry)
	TelemetryPaginated(c *gin.Context, moduleId int, deviceId int, params swagger.TelemetryPaginatedParams)
	// Получение последнего значения телеметрии
	// (GET /devices/{module_id}/{device_id}/telemetry/latest)
	TelemetryLatests(c *gin.Context, moduleId int, deviceId int)
}

type TelemetryRepository interface {
	GetLatest(
		ctx context.Context,
		moduleId, delviceId int,
	) (repository.TelemetryValue, error)

	GetPaginated(
		ctx context.Context,
		moduleId, deviceId int,
		limit, offset int,
	) ([]repository.TelemetryValue, bool, error)
}

type handler struct {
	repo TelemetryRepository
}

func New(repo TelemetryRepository) *handler {
	return &handler{
		repo: repo,
	}

}

func (h *handler) TelemetryPaginated(c *gin.Context, moduleId int, deviceId int, params swagger.TelemetryPaginatedParams) {
	perPage := conv.NormalizePerPage(pointer.Get(params.PerPage))
	page := pointer.Get(params.Page)

	values, hasMore, err := h.repo.GetPaginated(c.Request.Context(), moduleId, deviceId, perPage, page*perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, conv.RepoToSwaggerPaginatedValues(values, hasMore))
}

func (h *handler) TelemetryLatests(c *gin.Context, moduleId int, deviceId int) {
	value, err := h.repo.GetLatest(c.Request.Context(), moduleId, deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, conv.RepoToSwaggerTelemetryValue(value))
}
