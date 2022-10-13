package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/netip"
	"os"

	"github.com/lsmenicucci/simlab-vcluster/pkg/xml_template"
)

type ClusterWorkers struct{
	CPUS 			int 	`json:"cpus"`
	Memory			int 	`json:"memory"`
	DiskSize 		int 	`json:"disk_size"`
}

type ClusterController struct{
	CPUS 			int 	`json:"cpus"`
	Memory			int 	`json:"memory"`
	DiskSize 		int 	`json:"disk_size"`
	InstallationISO string 	`json:"installation_iso"`
}

type ClusterNetworkConfig struct{
	Address netip.Prefix `json:"address"`
}

type ClusterNetworks struct{
	Internal ClusterNetworkConfig `json:"internal"`
	External ClusterNetworkConfig `json:"external"`
}

type ClusterConfig struct{
	Prefix 				string 						`json:"prefix"`
	NumWorkers			int 						`json:"workers"`
	BaseDir   			string 						`json:"base_dir"`
	Worker  			ClusterWorkers 				`json:"worker_config"`
	Controller 			ClusterController			`json:"controller_config"`
	Networks 			ClusterNetworks 			`json:"networks"`
}


func (c * ClusterConfig) LoadFromFile(filepath string) error{
	raw, err := os.ReadFile(filepath)

	if (err != nil){
		return err
	}

	err = json.Unmarshal(raw, c)
	if (err != nil){
		return err
	}

	return nil
}

// Names
func (c *ClusterConfig) GetControllerName()string{
	return c.Prefix + "-controller"
}

func (c *ClusterConfig) GetControllerDiskName()string{
	return "controller.qcow2"
}

func (c *ClusterConfig) GetWorkerName(index int)string{
	return fmt.Sprintf("%s-worker-%d", c.Prefix, index)
}

func (c *ClusterConfig) GetWorkerDiskName(index int)string{
	return fmt.Sprintf("worker-%d.qcow2", index)
}

func (c *ClusterConfig) GetStoragePoolName()string{
	return c.Prefix
}

func (c *ClusterConfig) BuildControllerXML() (string, error){
	cfg := xml_template.NodeConfig{
		Name: c.GetControllerName(),
		Memory: c.Controller.Memory,
		CPUS: c.Controller.CPUS,
		Disk: &xml_template.DiskConfig{ Pool: c.Prefix, Name: c.GetControllerDiskName() },
	}

	if (len(c.Controller.InstallationISO) > 0){
		cfg.Cdrom = &xml_template.CdromConfig{ ImagePath: c.Controller.InstallationISO }
	}

	for _, net_suffix := range []string{ "internal", "external" }{
		net_dev := xml_template.NetworkDeviceConfig{ Name: c.Prefix + "-" + net_suffix }
		cfg.Networks = append(cfg.Networks, net_dev)
	}

	buf := new(bytes.Buffer)
	err := xml_template.NodeXML.Execute(buf, cfg)

	return buf.String(), err
}

func (c *ClusterConfig) BuildControllerDiskXML() (string, error){
	cfg := xml_template.VolumeConfig{
		Filename: c.GetControllerDiskName(),
		Capacity: c.Controller.DiskSize,
	}

	buf := new(bytes.Buffer)
	err := xml_template.VolumeXML.Execute(buf, cfg)

	return buf.String(), err
}

func (c *ClusterConfig) BuildWorkerXML(index int) (string, error){
	cfg := xml_template.NodeConfig{
		Name: c.GetWorkerName(index),
		Memory: c.Worker.Memory,
		CPUS: c.Worker.CPUS,
		Disk: &xml_template.DiskConfig{ Pool: c.Prefix, Name: c.GetWorkerDiskName(index) },
	}

	for _, net_suffix := range []string{ "internal" }{
		net_dev := xml_template.NetworkDeviceConfig{ Name: c.Prefix + "-" + net_suffix }
		cfg.Networks = append(cfg.Networks, net_dev)
	}

	buf := new(bytes.Buffer)
	err := xml_template.NodeXML.Execute(buf, cfg)

	return buf.String(), err
}

func (c *ClusterConfig) BuildWorkerDiskXML(index int) (string, error){
	cfg := xml_template.VolumeConfig{
		Filename: c.GetWorkerDiskName(index),
		Capacity: c.Worker.DiskSize,
	}

	buf := new(bytes.Buffer)
	err := xml_template.VolumeXML.Execute(buf, cfg)

	return buf.String(), err
}

func (c *ClusterConfig) BuildInternalNetworkXML() (string, error){
	cfg := xml_template.NetworkConfig{
		Name: c.Prefix + "-internal",
		Internal: true,
		Address: c.Networks.Internal.Address.Addr(),
		Mask: netip.MustParseAddr("255.255.255.0"),
	}

	buf := new(bytes.Buffer)
	err := xml_template.NetworkXML.Execute(buf, cfg)

	return buf.String(), err
}

func (c *ClusterConfig) BuildExternalNetworkXML(hostMacs map[string]string) (string, error){
	cfg := xml_template.NetworkConfig{
		Name: c.Prefix + "-external",
		Internal: false,
		Address: c.Networks.Internal.Address.Addr(),
		Mask: netip.MustParseAddr("255.255.255.0"),
		DHCP: &xml_template.DHCPConfig{},
	}

	cfg.DHCP.Start, cfg.DHCP.End = getAddrRange(c.Networks.External.Address)

	for ip, mac := range hostMacs{
		hostCfg := xml_template.DHCPHost{ MAC: mac, IP: ip }
		cfg.DHCP.Hosts = append(cfg.DHCP.Hosts, hostCfg)
	}

	buf := new(bytes.Buffer)
	err := xml_template.NetworkXML.Execute(buf, cfg)

	return buf.String(), err
}

func (c *ClusterConfig) BuildPoolXML() (string, error){
	cfg := xml_template.PoolConfig{
		Name: c.GetStoragePoolName(),
		Path: c.BaseDir,
	}

	buf := new(bytes.Buffer)
	err := xml_template.PoolXML.Execute(buf, cfg)

	return buf.String(), err
}