package xml_template

import (
	_ "embed"
	"text/template"
)

type DiskConfig struct{
	Pool string
	Name string
}

type CdromConfig struct{
	ImagePath string
}

type NetworkDeviceConfig struct{
	Name string
}

type NodeConfig struct{
	Name string
	Memory int
	CPUS int 

	Disk *DiskConfig
	Cdrom *CdromConfig 
	Networks []NetworkDeviceConfig
}

func (nc *NodeConfig) AddNetwork(net *NetworkConfig){
	dev := NetworkDeviceConfig{ Name: net.Name }
	nc.Networks = append(nc.Networks, dev)
}

//go:embed templates/domain.xml
var NodeXMLTemplate string
var NodeXML = template.Must(template.New("node").Parse(NodeXMLTemplate))