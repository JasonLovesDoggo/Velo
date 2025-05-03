package node

type Info struct {
	ID           string            `json:"id"`
	Hostname     string            `json:"hostname"`
	Address      string            `json:"address"`
	Labels       map[string]string `json:"labels"`
	Capacity     Resources         `json:"capacity"`
	Conditions   []string          `json:"conditions"`
	Role         string            `json:"role"`
	Manager      bool              `json:"is_manager"`
	Availability string            `json:"availability"`
}

type Resources struct {
	CPU    int   `json:"cpu_cores"`
	Memory int64 `json:"memory_bytes"`
	Disk   int64 `json:"disk_bytes"`
	GPU    int   `json:"gpu_count"`
}
