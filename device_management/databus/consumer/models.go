package consumer

import "management/swagger"

type Device struct {
	Id       int `json:"id"`
	ModuleId int `json:"module_id"`
}

type Command struct {
	Action string               `json:"action"`
	Device Device               `json:"device"`
	Value  *swagger.DeviceValue `json:"value"`
}
