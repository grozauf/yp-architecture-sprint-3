package producer

import "management/swagger"

type CommandResult struct {
	Action string               `json:"action"`
	Info   *swagger.DeviceInfo  `json:"info"`
	Value  *swagger.DeviceValue `json:"value"`
	Err    *string              `json:"err"`
}
