package producer

const (
	CommandTelemetryLatest    = "latest"
	CommandTelemetryPaginated = "paginated"
)

type Device struct {
	Id       int `json:"id"`
	ModuleId int `json:"module_id"`
}

type PaginationParams struct {
	Page    *int `json:"page"`
	PerPage *int `json:"per_page"`
}

type CommandTelemetryIn struct {
	Action          string            `json:"action"`
	Device          Device            `json:"device"`
	PaginatedParams *PaginationParams `json:"paginated_params"`
}

type ManagementCommand struct {
	Action string `json:"action"`
	Device Device `json:"device"`
}
