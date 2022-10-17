package components

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/beevik/etree"
	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

//go:embed templates/pool.xml
var PoolXMLTemplate string
var PoolXML = template.Must(template.New("pool").Parse(PoolXMLTemplate))

type StoragePool struct{
	Path string
	Name string
}

func (ns *StoragePool) GetXPaths() map[string]string{
	return map[string]string{
		"name": "/network/name",
		"netip": "/network/ip",
		"foward": "/network/foward[@mode='nat']",
	}
}

func (ns *StoragePool) Load(data map[string][]*etree.Element) error{
	err := EnsureXMLNode(data, ns.GetXPaths(), []string{"name", "netip"})
	if (err != nil){
		return err
	}
	// load name
	ns.Name = data["name"][0].Text()


	return nil
}

func StoragePoolExists(l *libvirt.Libvirt, name string) (bool, error){
	pools, _, err :=  l.ConnectListAllStoragePools(1, 0)
	if (err != nil){
		return false, err
	}

	for _, n := range pools{
		if (name == n.Name){
			return true, nil
		}
	}
	
	return false, nil
}

func GetStoragePool(l *libvirt.Libvirt, name string) (*StoragePool, error){
	log := logrus.WithField("pool", name)
	exists, err := StoragePoolExists(l, name)
	if (err != nil){
		return nil, err
	}

	if (exists == false){
		return nil, nil
	}

	nobj, err := l.StoragePoolLookupByName(name)
	if (err != nil){
		return nil, err
	}

	if (&nobj != nil){
		nobj = nobj
	}

	// load network data
	nxml, err := l.StoragePoolGetXMLDesc(nobj, 0)
	log.Debug("Parsing pool XML:\n"+nxml)
	state, err := LoadXMLData(&StoragePool{}, []byte(nxml))

	if (err != nil){
		log.WithError(err).Debug("Failed loading from XML")
		return nil, err
	}

	return state, nil
}

func SetStoragePool(l *libvirt.Libvirt, cfg *StoragePool) error{
	buf := new(bytes.Buffer)
	err := PoolXML.Execute(buf, cfg)
	if (err != nil){
		logrus.WithError(err).Debug("Failed encoding network state")
		return err
	}

	_, err = l.StoragePoolDefineXML(buf.String(), 0)
	if (err != nil){
		logrus.WithError(err).Debug("Failed setting network state")
		return err
	}

	return nil
}
