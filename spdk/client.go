package spdk

import (
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/longhorn/go-spdk-helper/jsonrpc"
	"github.com/longhorn/go-spdk-helper/spdk/bdev"
	"github.com/longhorn/go-spdk-helper/spdk/nvmf"
	"github.com/longhorn/go-spdk-helper/types"
)

const (
	DefaultTimeoutInSecond = 30
)

type Client struct {
	conn net.Conn

	jsonCli *jsonrpc.Client
}

func NewClient() (*Client, error) {
	conn, err := net.Dial(types.DefaultJsonServerNetwork, types.DefaultUnixDomainSocketPath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening socket for spdk client")
	}

	return &Client{
		conn:    conn,
		jsonCli: jsonrpc.NewClient(conn),
	}, nil
}

// BdevGetBdevs get information about block devices (bdevs).
//
//	"name": Optional. If this is not specified, the function will list all block devices.
//
//	"timeout":  0 by default, meaning the method returns immediately whether the bdev exists or not.
func (c *Client) BdevGetBdevs(name string, timeout uint64) ([]bdev.BdevInfo, error) {
	req := bdev.BdevGetBdevsRequest{
		Name:    name,
		Timeout: timeout,
	}

	outputJson, err := c.jsonCli.SendCommand("bdev_get_bdevs", req)
	if err != nil {
		return nil, err
	}
	bdevInfoList, ok := outputJson.([]bdev.BdevInfo)
	if !ok {
		return nil, fmt.Errorf("invalid output of BdevGetBdevs: %v", outputJson)
	}

	return bdevInfoList, nil
}

// BdevAioCreate constructs Linux AIO bdev.
func (c *Client) BdevAioCreate(filename, name string, blockSize uint64) (string, error) {
	req := bdev.BdevAioCreateRequest{
		Name:      name,
		Filename:  filename,
		BlockSize: blockSize,
	}

	result, err := c.jsonCli.SendCommand("bdev_aio_create", req)
	if err != nil {
		return "", err
	}
	bdevName, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid result of BdevAioCreate: %v", result)
	}

	return bdevName, nil
}

// BdevAioDelete deletes Linux AIO bdev.
func (c *Client) BdevAioDelete(name string) (bool, error) {
	req := bdev.BdevAioDeleteRequest{
		Name: name,
	}

	result, err := c.jsonCli.SendCommand("bdev_aio_delete", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid result of BdevAioDelete: %v", result)
	}

	return deleted, nil
}

// BdevAioGet will list all aio bdevs if a name is not specified.
//
//	"timeout": 0 by default, meaning the method returns immediately whether the aio exists or not.
func (c *Client) BdevAioGet(name string, timeout uint64) ([]bdev.BdevAioInfo, error) {
	req := bdev.BdevGetBdevsRequest{
		Name:    name,
		Timeout: timeout,
	}

	result, err := c.jsonCli.SendCommand("bdev_get_bdevs", req)
	if err != nil {
		return nil, err
	}
	bdevInfoList, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid output of BdevGetBdevs: %v", result)
	}

	res := []bdev.BdevAioInfo{}
	for _, b := range bdevInfoList {
		a, ok := b.(bdev.BdevAioInfo)
		if !ok {
			continue
		}
		res = append(res, a)
	}

	return res, nil
}

// BdevLvolCreateLvstore constructs a logical volume store.
func (c *Client) BdevLvolCreateLvstore(bdevName, lvsName string) (string, error) {
	req := bdev.BdevLvolCreateLvstoreRequest{
		BdevName: bdevName,
		LvsName:  lvsName,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_create_lvstore", req)
	if err != nil {
		return "", err
	}
	uuid, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid output of BdevLvolCreateLvstore: %v", result)
	}

	return uuid, nil
}

// BdevLvolDeleteLvstore destroys a logical volume store. It receives either lvs_name or UUID.
func (c *Client) BdevLvolDeleteLvstore(lvsName, uuid string) (bool, error) {
	req := bdev.BdevLvolDeleteLvstoreRequest{
		LvsName: lvsName,
		Uuid:    uuid,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_delete_lvstore", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of BdevLvolDeleteLvstore: %v", result)
	}

	return deleted, nil
}

// BdevLvolGetLvstore receives either lvs_name or UUID.
func (c *Client) BdevLvolGetLvstore(lvsName, uuid string) ([]bdev.LvstoreInfo, error) {
	req := bdev.BdevLvolGetLvstoreRequest{
		LvsName: lvsName,
		Uuid:    uuid,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_get_lvstores", req)
	if err != nil {
		return nil, err
	}
	lvstoreInfoList, ok := result.([]bdev.LvstoreInfo)
	if !ok {
		return nil, fmt.Errorf("invalid output of BdevLvolGetLvstore: %v", result)
	}

	return lvstoreInfoList, nil
}

// BdevLvolCreate constructs a logical volume store.
func (c *Client) BdevLvolCreate(lvstoreName, lvolName string, size uint64) (string, error) {
	req := bdev.BdevLvolCreateRequest{
		LvsName:  lvstoreName,
		LvolName: lvolName,
		Size:     size,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_create", req)
	if err != nil {
		return "", err
	}
	uuid, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid output of BdevLvolCreate: %v", result)
	}

	return uuid, nil
}

// BdevLvolDelete destroys a logical volume.
//
//	"name": UUID or alias of the logical volume. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>.
func (c *Client) BdevLvolDelete(name string) (bool, error) {
	req := bdev.BdevLvolDeleteRequest{
		Name: name,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_delete", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of BdevLvolDelete: %v", result)
	}

	return deleted, nil
}

// BdevLvolGet gets information about bdev lvols.
//
//	"name": UUID or alias of the logical volume. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>.
//		 	This is optional. If this is not specified, the function will list all block devices.
//
//	"timeout": 0 by default, meaning the method returns immediately whether the lvol exists or not.
func (c *Client) BdevLvolGet(name string, timeout uint64) ([]bdev.BdevLvolInfo, error) {
	req := bdev.BdevGetBdevsRequest{
		Name:    name,
		Timeout: timeout,
	}

	result, err := c.jsonCli.SendCommand("bdev_get_bdevs", req)
	if err != nil {
		return nil, err
	}
	bdevInfoList, ok := result.([]bdev.BdevLvolInfo)
	if !ok {
		return nil, fmt.Errorf("invalid output of BdevLvolGet: %v", result)
	}

	res := []bdev.BdevLvolInfo{}
	for _, b := range bdevInfoList {
		if b.ProductName != bdev.BdevProductNameLvol {
			continue
		}
		if _, ok := b.DriverSpecific[bdev.BdevDriverNameLvol]; !ok {
			continue
		}
		res = append(res, b)
	}

	return res, nil
}

// BdevLvolSnapshot capture a snapshot of the current state of a logical volume as a new bdev lvol.
//
//	"name": UUID or alias of the logical volume to create a snapshot from. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>.
func (c *Client) BdevLvolSnapshot(name, snapshotName string) (string, error) {
	req := bdev.BdevLvolSnapshotRequest{
		LvolName:     name,
		SnapshotName: snapshotName,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_snapshot", req)
	if err != nil {
		return "", err
	}
	uuid, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid output of BdevLvolSnapshot: %v", result)
	}

	return uuid, nil
}

// BdevLvolClone creates a logical volume based on a snapshot.
//
//	"name": UUID or alias of the logical volume/snapshot to clone. The alias of a lvol is <LVSTORE NAME>/<SNAPSHOT or LVOL NAME>.
func (c *Client) BdevLvolClone(name, cloneName string) (string, error) {
	req := bdev.BdevLvolCloneRequest{
		SnapshotName: name,
		CloneName:    cloneName,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_clone", req)
	if err != nil {
		return "", err
	}
	uuid, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid output of BdevLvolClone: %v", result)
	}

	return uuid, nil
}

// BdevLvolDecoupleParent decouples the parent of a logical volume.
// For unallocated clusters which is allocated in the parent, they are allocated and copied from the parent,
// but for unallocated clusters which is thin provisioned in the parent, they are kept thin provisioned. Then all dependencies on the parent are removed.
//
//	"name": UUID or alias of the logical volume to decouple the parent of it. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>.
func (c *Client) BdevLvolDecoupleParent(name string) (string, error) {
	req := bdev.BdevLvolDecoupleParentRequest{
		Name: name,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_decouple_parent", req)
	if err != nil {
		return "", err
	}
	uuid, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid output of BdevLvolDecoupleParent: %v", result)
	}

	return uuid, nil
}

// BdevLvolResize resizes a logical volume.
//
//	"name": UUID or alias of the logical volume to resize.
//
//	"size": Desired size of the logical volume in bytes.
func (c *Client) BdevLvolResize(name string, size uint64) (bool, error) {
	req := bdev.BdevLvolResizeRequest{
		Name: name,
		Size: size,
	}

	result, err := c.jsonCli.SendCommand("bdev_lvol_resize", req)
	if err != nil {
		return false, err
	}
	resized, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of BdevLvolResize: %v", result)
	}

	return resized, nil
}

// BdevRaidCreate constructs a new RAID bdev.
//
//	"name": a RAID bdev name rather than an alias or a UUID.
//
//	"raidLevel": RAID level. It can be "0"/"raid0", "1"/"raid1", "5f"/"raid5f", or "concat".

func (c *Client) BdevRaidCreate(name, raidLevel string, stripSizeKb uint32, baseBdevs []string) (bool, error) {
	req := bdev.BdevRaidCreateRequest{
		Name:        name,
		RaidLevel:   raidLevel,
		StripSizeKb: stripSizeKb,
		BaseBdevs:   baseBdevs,
	}

	result, err := c.jsonCli.SendCommand("bdev_raid_create", req)
	if err != nil {
		return false, err
	}
	created, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of BdevRaidCreate: %v", result)
	}

	return created, nil
}

// BdevRaidDelete destroys a logical volume.
func (c *Client) BdevRaidDelete(name string) (bool, error) {
	req := bdev.BdevRaidDeleteRequest{
		Name: name,
	}

	result, err := c.jsonCli.SendCommand("bdev_raid_delete", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of BdevRaidDelete: %v", result)
	}

	return deleted, nil
}

// BdevRaidGetBdevs is used to list all the raid bdev details based on the input category requested.
//
//	"category": This should be one of 'all', 'online', 'configuring' or 'offline'.
//	  'all' means all the raid bdevs whether they are online or configuring or offline.
//	  'online' is the raid bdev which is registered with bdev layer.
//	  'configuring' is the raid bdev which does not have full configuration discovered yet.
//	  'offline' is the raid bdev which is not registered with bdev as of now and it has encountered any error or user has requested to offline the raid bdev.
func (c *Client) BdevRaidGetBdevs(category string) ([]bdev.BdevRaidInfo, error) {
	req := bdev.BdevRaidGetBdevsRequest{
		Category: category,
	}

	result, err := c.jsonCli.SendCommand("bdev_raid_get_bdevs", req)
	if err != nil {
		return nil, err
	}
	bdevRaidInfoList, ok := result.([]bdev.BdevRaidInfo)
	if !ok {
		return nil, fmt.Errorf("invalid output of BdevRaidGetBdevs: %v", result)
	}

	return bdevRaidInfoList, nil
}

// BdevNvmeAttachController constructs NVMe bdev.
//
//	"name": Name of the NVMe controller, prefix for each bdev name.
//
//	"trtype": NVMe-oF target trtype: "tcp", "rdma" or "pcie"
//
//	"traddr": NVMe-oF target address: ip or BDF
//
//	"subnqn": NVMe-oF target subnqn. It can be the nvmf subsystem nqn.
//
//	"trsvcid": NVMe-oF target trsvcid: port number
//
//	"adrfam": NVMe-oF target adrfam: ipv4, ipv6, ib, fc, intra_host
func (c *Client) BdevNvmeAttachController(name, subnqn, trtype, adrfam, traddr, trsvcid string) (string, error) {
	req := bdev.BdevNvmeAttachControllerRequest{
		Name:    name,
		Subnqn:  subnqn,
		Trtype:  trtype,
		Adrfam:  adrfam,
		Traddr:  traddr,
		Trsvcid: trsvcid,
	}

	result, err := c.jsonCli.SendCommand("bdev_nvme_attach_controller", req)
	if err != nil {
		return "", err
	}
	bdevName, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("invalid output of BdevNvmeAttachController: %v", result)
	}

	return bdevName, nil
}

// BdevNvmeDetachController detach NVMe controller and delete any associated bdevs.
//
//	"name": Name of the NVMe controller.
func (c *Client) BdevNvmeDetachController(name string) (bool, error) {
	req := bdev.BdevNvmeDetachControllerRequest{
		Name: name,
	}

	result, err := c.jsonCli.SendCommand("bdev_nvme_detach_controller", req)
	if err != nil {
		return false, err
	}
	detached, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of BdevNvmeDetachController: %v", result)
	}

	return detached, nil
}

// BdevNvmeGetControllers gets information about bdev NVMe controllers.
//
//	"name": Name of the NVMe controller. Optional. If this is not specified, the function will list all NVMe controllers.
func (c *Client) BdevNvmeGetControllers(name string) ([]bdev.BdevNvmeControllerInfo, error) {
	req := bdev.BdevNvmeGetControllersRequest{
		Name: name,
	}

	result, err := c.jsonCli.SendCommand("bdev_nvme_get_controllers", req)
	if err != nil {
		return nil, err
	}
	controllerInfoList, ok := result.([]bdev.BdevNvmeControllerInfo)
	if !ok {
		return nil, fmt.Errorf("invalid output of BdevNvmeGetControllers: %v", result)
	}

	return controllerInfoList, nil
}

// NvmfCreateTransport initializes an NVMe-oF transport with the given options.
//
//	"trtype": Transport type, "tcp" or "rdma"
func (c *Client) NvmfCreateTransport(trtype string) (bool, error) {
	req := nvmf.NvmfCreateTransportRequest{
		Trtype: trtype,
	}

	result, err := c.jsonCli.SendCommand("nvmf_create_transport", req)
	if err != nil {
		return false, err
	}
	created, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of NvmfCreateTransport: %v", result)
	}

	return created, nil
}

// NvmfGetTransport lists all transports if no parameters specified.
//
//	"trtype": Transport type, "tcp" or "rdma"
//
//	"tgtName": Parent NVMe-oF target name.
func (c *Client) NvmfGetTransport(trtype, tgtName string) ([]nvmf.NvmfTransport, error) {
	req := nvmf.NvmfGetTransportRequest{
		Trtype:  trtype,
		TgtName: tgtName,
	}

	result, err := c.jsonCli.SendCommand("nvmf_get_transport", req)
	if err != nil {
		return nil, err
	}
	transportList, ok := result.([]nvmf.NvmfTransport)
	if !ok {
		return nil, fmt.Errorf("invalid output of NvmfGetTransport: %v", result)
	}

	return transportList, nil
}

// NvmfCreateSubsystem constructs an NVMe over Fabrics target subsystem..
//
//	"nqn": Subsystem NQN.
func (c *Client) NvmfCreateSubsystem(nqn string) (bool, error) {
	req := nvmf.NvmfCreateSubsystemRequest{
		Nqn: nqn,
	}

	result, err := c.jsonCli.SendCommand("nvmf_create_subsystem", req)
	if err != nil {
		return false, err
	}
	created, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of NvmfCreateSubsystem: %v", result)
	}

	return created, nil
}

// NvmfDeleteSubsystem constructs an NVMe over Fabrics target subsystem..
//
//	"nqn": Subsystem NQN.
//
//	"tgtName": Parent NVMe-oF target name. Optional.
func (c *Client) NvmfDeleteSubsystem(nqn, targetName string) (bool, error) {
	req := nvmf.NvmfDeleteSubsystemRequest{
		Nqn:     nqn,
		TgtName: targetName,
	}

	result, err := c.jsonCli.SendCommand("nvmf_delete_subsystem", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of NvmfDeleteSubsystem: %v", result)
	}

	return deleted, nil
}

// NvmfGetSubsystems lists all subsystem for the specified NVMe-oF target.
//
//	"tgtName": Parent NVMe-oF target name.
func (c *Client) NvmfGetSubsystems(tgtName string) ([]nvmf.NvmfSubsystem, error) {
	req := nvmf.NvmfGetSubsystemsRequest{
		TgtName: tgtName,
	}

	result, err := c.jsonCli.SendCommand("nvmf_get_subsystems", req)
	if err != nil {
		return nil, err
	}
	subsystemList, ok := result.([]nvmf.NvmfSubsystem)
	if !ok {
		return nil, fmt.Errorf("invalid output of NvmfGetSubsystems: %v", result)
	}

	return subsystemList, nil
}

// NvmfSubsystemAddNs constructs an NVMe over Fabrics target subsystem..
//
//	"nqn": Subsystem NQN.
//
//	"bdevName": Name of bdev to expose as a namespace.
func (c *Client) NvmfSubsystemAddNs(nqn, bdevName string) (uint32, error) {
	req := nvmf.NvmfSubsystemAddNsRequest{
		Nqn:       nqn,
		Namespace: nvmf.NvmfSubsystemNamespace{BdevName: bdevName},
	}

	result, err := c.jsonCli.SendCommand("nvmf_subsystem_add_ns", req)
	if err != nil {
		return 0, err
	}
	nsid, ok := result.(uint32)
	if !ok {
		return 0, fmt.Errorf("invalid output of NvmfSubsystemAddNs: %v", result)
	}

	return nsid, nil
}

// NvmfSubsystemRemoveNs constructs an NVMe over Fabrics target subsystem..
//
//	"nqn": Subsystem NQN.
//
//	"nsid": Namespace ID.
func (c *Client) NvmfSubsystemRemoveNs(nqn string, nsid uint32) (bool, error) {
	req := nvmf.NvmfSubsystemRemoveNsRequest{
		Nqn:  nqn,
		Nsid: nsid,
	}

	result, err := c.jsonCli.SendCommand("nvmf_subsystem_remove_ns", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of NvmfSubsystemRemoveNs: %v", result)
	}

	return deleted, nil
}

// NvmfSubsystemsGetNss lists all namespaces for the specified NVMe-oF target subsystem if bdev name or NSID is not specified.
//
//	"nqn": Subsystem NQN.
//
//	"bdevName": Name of bdev to expose as a namespace. Optional. It's better not to specify this and "nsid" simultaneously.
//
//	"nsid": Namespace ID. Optional. It's better not to specify this and "bdevName" simultaneously.
func (c *Client) NvmfSubsystemsGetNss(nqn, bdevName string, nsid uint32) ([]nvmf.NvmfSubsystemNamespace, error) {
	req := nvmf.NvmfGetSubsystemsRequest{}

	result, err := c.jsonCli.SendCommand("nvmf_get_subsystems", req)
	if err != nil {
		return nil, err
	}
	subsystemList, ok := result.([]nvmf.NvmfSubsystem)
	if !ok {
		return nil, fmt.Errorf("invalid output of NvmfGetSubsystems: %v", result)
	}

	nsList := []nvmf.NvmfSubsystemNamespace{}
	for _, subsystem := range subsystemList {
		if subsystem.Nqn != nqn {
			continue
		}
		for _, ns := range subsystem.Namespaces {
			if nsid > 0 && ns.Nsid != nsid {
				continue
			}
			if bdevName != "" && ns.BdevName != bdevName {
				continue
			}
			nsList = append(nsList, ns)
		}
	}

	return nsList, nil
}

// NvmfSubsystemAddListener adds a new listen address to an NVMe-oF subsystem.
//
//		"nqn": Subsystem NQN.
//
//		"trtype": NVMe-oF target trtype: "tcp", "rdma" or "pcie"
//
//	 "adrfam": Address family ("IPv4", "IPv6", "IB", or "FC")
//
//		"traddr": NVMe-oF target address: ip or BDF
//
//		"trsvcid": NVMe-oF target trsvcid: port number
func (c *Client) NvmfSubsystemAddListener(nqn, trtype, traddr, trsvcid, adrfam string) (bool, error) {
	req := nvmf.NvmfSubsystemAddListenerRequest{
		Nqn: nqn,
		ListenAddress: nvmf.ListenAddress{
			Trtype:  trtype,
			Adrfam:  adrfam,
			Traddr:  traddr,
			Trsvcid: trsvcid,
		},
	}

	result, err := c.jsonCli.SendCommand("nvmf_subsystem_add_listener", req)
	if err != nil {
		return false, err
	}
	created, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of NvmfSubsystemAddListener: %v", result)
	}

	return created, nil
}

// NvmfSubsystemRemoveListener removes a listen address from an NVMe-oF subsystem.
//
//		"nqn": Subsystem NQN.
//
//		"trtype": NVMe-oF target trtype: "tcp", "rdma" or "pcie"
//
//	 "adrfam": Address family ("IPv4", "IPv6", "IB", or "FC")
//
//		"traddr": NVMe-oF target address: ip or BDF
//
//		"trsvcid": NVMe-oF target trsvcid: port number
func (c *Client) NvmfSubsystemRemoveListener(nqn, trtype, traddr, trsvcid, adrfam string) (bool, error) {
	req := nvmf.NvmfSubsystemRemoveListenerRequest{
		Nqn: nqn,
		ListenAddress: nvmf.ListenAddress{
			Trtype:  trtype,
			Adrfam:  adrfam,
			Traddr:  traddr,
			Trsvcid: trsvcid,
		},
	}

	result, err := c.jsonCli.SendCommand("nvmf_subsystem_remove_listener", req)
	if err != nil {
		return false, err
	}
	deleted, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("invalid output of NvmfSubsystemRemoveListener: %v", result)
	}

	return deleted, nil
}
