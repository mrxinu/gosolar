package gosolar

type Subnet struct {
	Address        string `json:"Address"`
	CIDR           string `json:"CIDR"`
	Comments       string `json:"Comments"`
	AddressMask    string `json:"AddressMask"`
	DisplayName    string `json:"DisplayName"`
	FriendlyName   string `json:"FriendlyName"`
	TotalCount     int    `json:"totalCount"`
	UsedCount      int    `json:"UsedCount"`
	AvailableCount int    `json:"AvailableCount"`
	ReservedCount  int    `json:"ReservedCount"`
	//TransientCount string `json"Transient"`
}
