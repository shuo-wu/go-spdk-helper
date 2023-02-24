package bdev

type LvstoreInfo struct {
	Uuid              string `json:"uuid"`
	Name              string `json:"name"`
	BaseBdev          string `json:"base_bdev"`
	TotalDataClusters int    `json:"total_data_clusters"`
	FreeClusters      int    `json:"free_clusters"`
	BlockSize         int    `json:"block_size"`
	ClusterSize       int    `json:"cluster_size"`
}

type BdevLvolInfo struct {
	BdevInfoBasic

	DriverSpecific map[string]BdevLvolDriverSpecificInfo `json:"driver_specific"`
}

type BdevLvolDriverSpecificInfo struct {
	LvolStoreUuid string `json:"lvol_store_uuid"`
	BaseBdev      string `json:"base_bdev"`
	ThinProvision bool   `json:"thin_provision"`
	Snapshot      bool   `json:"snapshot"`
	Clone         bool   `json:"clone"`
}

type BdevLvolCreateLvstoreRequest struct {
	BdevName string `json:"bdev_name"`
	LvsName  string `json:"lvs_name"`
	// ClusterSz                 uint32 `json:"cluster_sz,omitempty"`
	// ClearMethod               string `json:"clear_method,omitempty"`
	// NumMdPagesPerClusterRatio uint32 `json:"num_md_pages_per_cluster_ratio,omitempty"`
}

type BdevLvolDeleteLvstoreRequest struct {
	Uuid    string `json:"uuid,omitempty"`
	LvsName string `json:"lvs_name,omitempty"`
}

type BdevLvolGetLvstoreRequest struct {
	Uuid    string `json:"uuid,omitempty"`
	LvsName string `json:"lvs_name,omitempty"`
}

type BdevLvolCreateRequest struct {
	LvsName  string `json:"lvs_name"`
	LvolName string `json:"lvol_name"`
	Size     uint64 `json:"size"`

	//ClearMethod   string `json:"clear_method"`
	//ThinProvision bool `json:"thin_provision"`
}

type BdevLvolDeleteRequest struct {
	Name string `json:"name"`
}

type BdevLvolSnapshotRequest struct {
	LvolName     string `json:"lvol_name"`
	SnapshotName string `json:"snapshot_name"`
}

type BdevLvolCloneRequest struct {
	SnapshotName string `json:"snapshot_name"`
	CloneName    string `json:"clone_name"`
}

type BdevLvolDecoupleParentRequest struct {
	Name string `json:"name"`
}

type BdevLvolResizeRequest struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
}
