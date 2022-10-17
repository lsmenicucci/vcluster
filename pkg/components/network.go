package components

import (
	"bytes"
	_ "embed"
	"errors"
	"net/netip"
	"text/template"

	"github.com/beevik/etree"
	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

//go:embed templates/network.xml
var NetworkXMLTemplate string
var NetworkXML = template.Must(template.New("network").Parse(NetworkXMLTemplate))

type DHCPHost struct{
	MAC string `json:"mac"`
	IP 	string `json:"ip"`
}

type DHCPConfig struct{
	Start 	netip.Addr 	`json:"start"`
	End 	netip.Addr 	`json:"end"`
	Hosts 	[]DHCPHost 	`json:"hosts"`
}

type Network struct{
	Name 		string 		`json:"name"`
	Internal 	bool 		`json:"internal"`
	Address 	netip.Addr 	`json:"address"`
	Mask 		netip.Addr 	`json:"mask"`
	DHCP 		*DHCPConfig `json:"dhcp"`
}

func (ns *Network) GetXPaths() map[string]string{
	return map[string]string{
		"name": "/network/name",
		"netip": "/network/ip",
		"foward": "/network/foward[@mode='nat']",
	}
}

func (ns *Network) Load(data map[string][]*etree.Element) error{
	err := EnsureXMLNode(data, ns.GetXPaths(), []string{"name", "netip"})
	if (err != nil){
		return err
	}
	// load name
	ns.Name = data["name"][0].Text()

	// load net addr
	addr, err := netip.ParseAddr(data["netip"][0].SelectAttrValue("address", ""))
	if (err != nil){
		return errors.New("Failed parsing address")
	}
	ns.Address = addr

	// load net mask
	mask, err := netip.ParseAddr(data["netip"][0].SelectAttrValue("netmask", ""))
	if (err != nil){
		return errors.New("Failed parsing netmask")
	}
	ns.Mask = mask

	// load internal flag
	if (len(data["foward"]) == 0){
		ns.Internal = true
	}else{
		ns.Internal = false
	}

	return nil
}

func NetworkExists(l *libvirt.Libvirt, name string) (bool, error){
	networks, _, err :=  l.ConnectListAllNetworks(1, 0)
	if (err != nil){
		return false, err
	}

	for _, n := range networks{
		if (name == n.Name){
			return true, nil
		}
	}
	
	return false, nil
}

func GetNetwork(l *libvirt.Libvirt, name string) (*Network, error){
	log := logrus.WithField("network", name)
	exists, err := NetworkExists(l, name)
	if (err != nil){
		return nil, err
	}

	if (exists == false){
		return nil, nil
	}

	nobj, err := l.NetworkLookupByName(name)
	if (err != nil){
		return nil, err
	}

	// load network data
	nxml, err := l.NetworkGetXMLDesc(nobj, 0)
	log.Debug("Parsing network XML:\n"+nxml)
	state, err := LoadXMLData(&Network{}, []byte(nxml))

	if (err != nil){
		log.WithError(err).Debug("Failed loading from XML")
		return nil, err
	}

	return state, nil
}

func SetNetwork(l *libvirt.Libvirt, cfg *Network) error{
	buf := new(bytes.Buffer)
	err := NetworkXML.Execute(buf, cfg)
	if (err != nil){
		logrus.WithError(err).Debug("Failed encoding network state")
		return err
	}

	_, err = l.NetworkDefineXML(buf.String())
	if (err != nil){
		logrus.WithError(err).Debug("Failed setting network state")
		return err
	}

	return nil
}


func ListNetworksNames(l *libvirt.Libvirt) ([]string, error){
	networks, _, err :=  l.ConnectListAllNetworks(1, 0)
	if (err != nil){
		return []string{}, err
	}

	netNames := make([]string, len(networks))
	for k, n := range networks{
		netNames[k] = n.Name
	}

	return netNames, nil
}