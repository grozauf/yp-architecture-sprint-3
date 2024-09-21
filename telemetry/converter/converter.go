package converter

import (
	"github.com/AlekSi/pointer"

	"telemetry/repository"
	"telemetry/swagger"
)

func RepoToSwaggerPaginatedValues(values []repository.TelemetryValue, hasMore bool) swagger.PaginatedValues {
	return swagger.PaginatedValues{
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
}

func RepoToSwaggerTelemetryValue(value repository.TelemetryValue) swagger.TelemetryValue {
	return swagger.TelemetryValue{
		Value:       value.Value,
		OccuranceAt: value.OccuranceAt.String(),
	}
}

func NormalizePerPage(perPage int) int {
	if perPage == 0 {
		perPage = 100
	}
	if perPage > 500 {
		perPage = 500
	}
	return perPage
}

func ErrorToStringOrNil(err error) *string {
	if err == nil {
		return nil
	}
	return pointer.ToString(err.Error())
}
