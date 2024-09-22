package repository

import (
	"context"
	"math/rand"
	"time"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) GetLatest(ctx context.Context, moduleId, delviceId int) (TelemetryValue, error) {
	return TelemetryValue{
		Value:       rand.Float32(),
		OccuranceAt: time.Now(),
	}, nil
}

func (r *Repository) GetPaginated(ctx context.Context, moduleId, deviceId int, limit, offset int) ([]TelemetryValue, bool, error) {
	result := make([]TelemetryValue, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, TelemetryValue{
			Value:       rand.Float32(),
			OccuranceAt: time.Now(),
		})
	}
	return result, true, nil
}
