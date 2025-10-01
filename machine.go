package clerk

type Machine struct {
	APIResource
	Object          string `json:"object"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	InstanceID      string `json:"instance_id"`
	DefaultTokenTTL int64  `json:"default_token_ttl"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type MachineWithScopedMachines struct {
	APIResource
	Machine
	ScopedMachines []*Machine `json:"scoped_machines"`
}

type MachineWithScopedMachinesAndSecretKey struct {
	APIResource
	MachineWithScopedMachines
	SecretKey string `json:"secret_key"`
}

type MachineList struct {
	APIResource
	Machines   []*MachineWithScopedMachines `json:"data"`
	TotalCount int64                        `json:"total_count"`
}

type MachineScope struct {
	APIResource
	Object        string `json:"object"`
	FromMachineID string `json:"from_machine_id"`
	ToMachineID   string `json:"to_machine_id"`
	CreatedAt     int64  `json:"created_at"`
}

type DeletedMachineScope struct {
	APIResource
	Object        string `json:"object"`
	FromMachineID string `json:"from_machine_id"`
	ToMachineID   string `json:"to_machine_id"`
	Deleted       bool   `json:"deleted"`
}

type MachineSecretKey struct {
	APIResource
	Object string `json:"object"`
	Secret string `json:"secret"`
}
