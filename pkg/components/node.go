package components

import (
	"bytes"
	_ "embed"
	"errors"
	"strconv"
	"text/template"

	"github.com/beevik/etree"
	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

//go:embed templates/domain.xml
var NodeXMLTemplate string
var NodeXML = template.Must(template.New("node").Parse(NodeXMLTemplate))

type DiskConfig struct{
	Pool string `json:"pool"`
	Name string `json:"name"`
}

type CdromConfig struct{
	ImagePath string 	`json:"iso_path"`
}

type NetworkDeviceConfig struct{
	Name string 			`json:"name"`
}

type Node struct{
	Name 		string 		`json:"name"`
	Memory 		int 		`json:"memory"`
	CPUS 		int 		`json:"cpus"`

	Disk 		*DiskConfig 			`json:"disk"`
	Cdrom 		*CdromConfig 			`json:"cdrom"`
	Networks	[]*NetworkDeviceConfig 	`json:"networks"`
}

func (ns *Node) GetXPaths() map[string]string{
	return map[string]string{
		"name": "/domain/name",
		"memory": "/domain/memory[@unit='KiB']",
		"cpus": "/domain/vcpu",
		"disk": "/domain/devices/disk[@type='volume']/source",
		"cdrom_iso": "/domain/devices/disk[@device='cdrom']/source",
		"networks": "/domain/interface[@type='network']/source",
	}
}

func (ns *Node) Load(data map[string][]*etree.Element) error{
	requiredKeys := []string{"name", "memory", "cpus", "disk"}
	err := EnsureXMLNode(data, ns.GetXPaths(), requiredKeys)
	if (err != nil){
		return err
	}

	// load name
	ns.Name = data["name"][0].Text()
	
	// load mem
	mem, err := strconv.Atoi(data["memory"][0].Text())
	if (err != nil){
		return errors.New("Invalid memory value")
	}
	ns.Memory = mem/(1024*1024)

	// load cpus
	cpus, err := strconv.Atoi(data["cpus"][0].Text())
	if (err != nil){
		return errors.New("Invalid memory value")
	}
	ns.CPUS = cpus

	// load disk
	diskAttrs, err := UnpackRequiredAttrs([]string{"pool", "volume"}, data["disk"][0])
	if (err != nil){
		return err
	}
	ns.Disk = &DiskConfig{
		Pool: diskAttrs["pool"],
		Name: diskAttrs["volume"],
	}

	// load networks
	ns.Networks = []*NetworkDeviceConfig{}
	for _, net := range data["networks"]{
		data, err := UnpackRequiredAttrs([]string{"network"}, net) 
		if (err != nil){
			return err
		}

		ns.Networks = append(ns.Networks, &NetworkDeviceConfig{ Name: data["name"] })
	}

	// load cdrom iso
	if (len(data["cdrom_iso"]) > 0){
		data, err := UnpackRequiredAttrs([]string{"file"}, data["cdrom_iso"][0]) 
		if (err != nil){
			return err
		}
		ns.Cdrom = &CdromConfig{ ImagePath: data["file"] }
	}

	return nil
}

func NodeExists(l *libvirt.Libvirt, name string) (bool, error){
	nodes, _, err :=  l.ConnectListAllDomains(1, 0)
	if (err != nil){
		return false, err
	}

	for _, n := range nodes{
		if (name == n.Name){
			return true, nil
		}
	}
	
	return false, nil
}

func GetNode(l *libvirt.Libvirt, name string) (*Node, error){
	log := logrus.WithField("node", name)

	exists, err := NodeExists(l, name)
	if (err != nil){
		return nil, err
	}

	if (exists == false){
		log.Debug("Node does not exist")
		return nil, nil
	}

	nobj, err := l.DomainLookupByName(name)
	if (err != nil){
		return nil, err
	}

	// load network data
	nxml, err := l.DomainGetXMLDesc(nobj, 0)
	log.Debug("Parsing node XML:\n"+nxml)
	state, err := LoadXMLData(&Node{}, []byte(nxml))

	if (err != nil){
		log.WithError(err).Debug("Failed loading from XML")
		return nil, err
	}

	return state, nil
}

func SetNode(l *libvirt.Libvirt, cfg *Node) error{
	buf := new(bytes.Buffer)
	err := NodeXML.Execute(buf, cfg)
	if (err != nil){
		logrus.WithError(err).Debug("Failed encoding node state")
		return err
	}

	_, err = l.DomainDefineXML(buf.String())
	if (err != nil){
		logrus.WithError(err).Debug("Failed setting node state")
		return err
	}

	return nil
}

func ListNodesNames(l *libvirt.Libvirt) ([]string, error){
	nodes, _, err :=  l.ConnectListAllDomains(1, 0)
	if (err != nil){
		return []string{}, err
	}

	nodesNames := make([]string, len(nodes))
	for k, n := range nodes{
		nodesNames[k] = n.Name
	}

	return nodesNames, nil
}