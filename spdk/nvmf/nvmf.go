package nvmf

type NvmfCreateTransportRequest struct {
	Trtype string `json:"trtype"`
}

type NvmfGetTransportRequest struct {
	TgtName string `json:"tgt_name,omitempty"`
	Trtype  string `json:"trtype,omitempty"`
}

type NvmfTransport struct {
	Trtype              string `json:"trtype"`
	MaxQueueDepth       uint32 `json:"max_queue_depth"`
	MaxIoQpairsPerCtrlr uint32 `json:"max_io_qpairs_per_ctrlr"`
	InCapsuleDataSize   uint32 `json:"in_capsule_data_size"`
	MaxIoSize           uint32 `json:"max_io_size"`
	IoUnitSize          uint32 `json:"io_unit_size"`
	MaxAqDepth          uint32 `json:"max_aq_depth"`
	NumSharedBuffers    uint32 `json:"num_shared_buffers"`
	BufCacheSize        uint32 `json:"buf_cache_size"`
	SockPriority        uint32 `json:"sock_priority"`
	AbortTimeoutSec     uint32 `json:"abort_timeout_sec"`
	DifInsertOrStrip    bool   `json:"dif_insert_or_strip"`
	Zcopy               bool   `json:"zcopy"`
	C2HSuccess          bool   `json:"c2h_success"`
}

type NvmfCreateSubsystemRequest struct {
	Nqn string `json:"nqn"`
}

type NvmfDeleteSubsystemRequest struct {
	Nqn     string `json:"nqn"`
	TgtName string `json:"tgt_name,omitempty"`
}

type NvmfGetSubsystemsRequest struct {
	TgtName string `json:"tgt_name,omitempty"`
}

type NvmfSubsystem struct {
	Nqn             string                       `json:"nqn"`
	Subtype         string                       `json:"subtype"`
	ListenAddresses []NvmfSubsystemListenAddress `json:"listen_addresses"`
	AllowAnyHost    bool                         `json:"allow_any_host"`
	Hosts           []NvmfSubsystemHost          `json:"hosts"`
	SerialNumber    string                       `json:"serial_number,omitempty"`
	ModelNumber     string                       `json:"model_number,omitempty"`
	MaxNamespaces   uint32                       `json:"max_namespaces,omitempty"`
	MinCntlid       uint32                       `json:"min_cntlid,omitempty"`
	MaxCntlid       uint32                       `json:"max_cntlid,omitempty"`
	Namespaces      []NvmfSubsystemNamespace     `json:"namespaces"`
}

type NvmfSubsystemListenAddress struct {
	Trtype  string `json:"trtype"`
	Adrfam  string `json:"adrfam"`
	Traddr  string `json:"traddr"`
	Trsvcid string `json:"trsvcid"`
}

type NvmfSubsystemNamespace struct {
	Nsid     uint32 `json:"nsid,omitempty"`
	BdevName string `json:"bdev_name"`
	Nguid    string `json:"nguid,omitempty"`
	Eui64    string `json:"eui64,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
	Anagrpid string `json:"anagrpid,omitempty"`
	PtplFile string `json:"ptpl_file,omitempty"`
}

type NvmfSubsystemHost struct {
	Nqn string `json:"nqn"`
}

type NvmfSubsystemAddNsRequest struct {
	Nqn       string                 `json:"nqn"`
	Namespace NvmfSubsystemNamespace `json:"namespace"`
	TgtName   string                 `json:"tgt_name,omitempty"`
}

type NvmfSubsystemRemoveNsRequest struct {
	Nqn     string `json:"nqn"`
	Nsid    uint32 `json:"nsid"`
	TgtName string `json:"tgt_name,omitempty"`
}

type NvmfSubsystemAddListenerRequest struct {
	Nqn           string        `json:"nqn"`
	ListenAddress ListenAddress `json:"listen_address"`

	TgtName string `json:"tgt_name,omitempty"`
}

type ListenAddress struct {
	Trtype  string `json:"trtype"`
	Adrfam  string `json:"adrfam"`
	Traddr  string `json:"traddr"`
	Trsvcid string `json:"trsvcid"`
}

type NvmfSubsystemRemoveListenerRequest struct {
	Nqn           string        `json:"nqn"`
	ListenAddress ListenAddress `json:"listen_address"`

	TgtName string `json:"tgt_name,omitempty"`
}
