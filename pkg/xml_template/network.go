package xml_template

import (
	_ "embed"
	"net/netip"
	"text/template"
)

type DHCPHost struct{
	MAC string
	IP string
}

type DHCPConfig struct{
	Start netip.Addr
	End netip.Addr
	Hosts []DHCPHost
}

type NetworkConfig struct{
	Name string
	Internal bool
	Address netip.Addr
	Mask netip.Addr
	DHCP *DHCPConfig
}

//go:embed templates/network.xml
var NetworkXMLTemplate string
var NetworkXML = template.Must(template.New("network").Parse(NetworkXMLTemplate))