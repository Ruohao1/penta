package ports

type Port struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
	Service  string `json:"service"`
}
