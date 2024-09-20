package handler

import (
	"context"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"

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
	perPage := pointer.Get(params.PerPage)
	if perPage == 0 {
		perPage = 100
	}
	if perPage > 500 {
		perPage = 500
	}
	page := pointer.Get(params.Page)

	values, hasMore, err := h.repo.GetPaginated(c.Request.Context(), moduleId, deviceId, perPage, page*perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	result := swagger.PaginatedValues{
		HasMore: hasMore,
		Values: func() []swagger.TelemetryValue {
			r := make([]swagger.TelemetryValue, 0, len(values))
			for i := range values {
				r = append(r, swagger.TelemetryValue{
					Value:       values[i].Value,
					OccuranceAt: values[i].OccuranceAt.String(),
				})
			}

			return r
		}(),
	}
	c.JSON(http.StatusOK, result)
}

func (h *handler) TelemetryLatests(c *gin.Context, moduleId int, deviceId int) {
	value, err := h.repo.GetLatest(c.Request.Context(), moduleId, deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(
		http.StatusOK,
		swagger.TelemetryValue{
			Value:       value.Value,
			OccuranceAt: value.OccuranceAt.String(),
		},
	)
}
