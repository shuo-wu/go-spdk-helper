package bdev

type BdevRaidInfo struct {
	Name                   string   `json:"name"`
	StripSizeKb            uint32   `json:"strip_size_kb"`
	State                  string   `json:"state"`
	RaidLevel              string   `json:"raid_level"`
	NumBaseBdevs           uint8    `json:"num_base_bdevs"`
	NumBaseBdevsDiscovered uint8    `json:"num_base_bdevs_discovered"`
	BaseBdevsList          []string `json:"base_bdevs_list"`
}

type BdevRaidCreateRequest struct {
	Name        string   `json:"name"`
	RaidLevel   string   `json:"raid_level"`
	BaseBdevs   []string `json:"base_bdevs"`
	StripSizeKb uint32   `json:"strip_size_kb"`
}

type BdevRaidDeleteRequest struct {
	Name string `json:"name"`
}

type BdevRaidGetBdevsRequest struct {
	Category string `json:"category"`
}
