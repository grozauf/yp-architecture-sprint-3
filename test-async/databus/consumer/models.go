package consumer

import "test_async/swagger"

type CommandTelemetryResult struct {
	Action  string                   `json:"action"`
	Values  []swagger.TelemetryValue `json:"values"`
	HasMore bool                     `json:"hasMore"`
	Err     *string                  `json:"err"`
}

