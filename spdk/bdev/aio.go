package bdev

type BdevAioInfo struct {
	BdevInfoBasic

	DriverSpecific map[string]BdevAioDriverSpecificInfo `json:"driver_specific"`
}

type BdevAioDriverSpecificInfo struct {
	Filename          string `json:"filename"`
	BlockSizeOverride bool   `json:"block_size_override"`
	Readonly          bool   `json:"readonly"`
}

type BdevAioCreateRequest struct {
	Name      string `json:"name"`
	Filename  string `json:"filename"`
	BlockSize uint64 `json:"block_size"`
}

type BdevAioDeleteRequest struct {
	Name string `json:"name"`
}
