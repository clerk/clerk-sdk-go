package clerk

type Machine struct {
	APIResource
	Object     string `json:"object"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	InstanceID string `json:"instance_id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type MachineList struct {
	APIResource
	Machines   []*Machine `json:"data"`
	TotalCount int64      `json:"total_count"`
}
