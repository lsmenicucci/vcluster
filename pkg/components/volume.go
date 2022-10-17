package components

import (
	"bytes"
	_ "embed"
	"errors"
	"text/template"

	"github.com/beevik/etree"
	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

//go:embed templates/volume.xml
var VolumeXMLTemplate string
var VolumeXML = template.Must(template.New("volume").Parse(VolumeXMLTemplate))

type Volume struct{
	Filename string
	Capacity int
}

func (ns *Volume) GetXPaths() map[string]string{
	return map[string]string{
		"name": "/network/name",
		"netip": "/network/ip",
		"foward": "/network/foward[@mode='nat']",
	}
}

func (ns *Volume) Load(data map[string][]*etree.Element) error{
	err := EnsureXMLNode(data, ns.GetXPaths(), []string{"name", "netip"})
	if (err != nil){
		return err
	}
	// load name
	ns.Filename = data["name"][0].Text()

	return nil
}

func VolumeExists(l *libvirt.Libvirt, poolName string, name string) (bool, error){
	poolExists, err := StoragePoolExists(l, poolName)
	if (err != nil){
		return false, err
	}

	if (poolExists == false){
		return false, nil
	}

	pool, err := l.StoragePoolLookupByName(poolName)
	if (err != nil){
		return false, err
	}

	vols, _, err :=  l.StoragePoolListAllVolumes(pool, 1, 0)
	if (err != nil){
		return false, err
	}

	for _, n := range vols{
		if (name == n.Name){
			return true, nil
		}
	}
	
	return false, nil
}

func GetVolume(l *libvirt.Libvirt, poolName string, name string) (*Volume, error){
	log := logrus.WithFields(logrus.Fields{ "pool": name, "volume": name })

	exists, err := VolumeExists(l, poolName, name)
	if (err != nil){
		return nil, err
	}

	if (exists == false){
		return nil, nil
	}

	pool, err := l.StoragePoolLookupByName(poolName)
	if (err != nil){
		return nil, err
	}

	nobj, err := l.StorageVolLookupByName(pool, name)
	if (err != nil){
		return nil, err
	}

	// load network data
	nxml, err := l.StorageVolGetXMLDesc(nobj, 0)
	log.Debug("Parsing volume XML:\n"+nxml)
	state, err := LoadXMLData(&Volume{}, []byte(nxml))

	if (err != nil){
		log.WithError(err).Debug("Failed loading from XML")
		return nil, err
	}

	return state, nil
}

func SetVolume(l *libvirt.Libvirt, poolName string, cfg *Volume) error{
	// get pool first
	poolExists, err := StoragePoolExists(l, poolName)
	if (err != nil){
		return err
	}

	if (poolExists == false){
		return errors.New("Storage pool does not exists")
	}

	pool, err := l.StoragePoolLookupByName(poolName)
	if (err != nil){
		return err
	}

	// create volume
	buf := new(bytes.Buffer)
	err = VolumeXML.Execute(buf, cfg)
	if (err != nil){
		logrus.WithError(err).Debug("Failed encoding network state")
		return err
	}

	_, err = l.StorageVolCreateXML(pool, buf.String(), 0)
	if (err != nil){
		logrus.WithError(err).Debug("Failed setting network state")
		return err
	}

	return nil
}

func ListVolumesNames(l *libvirt.Libvirt, poolName string) ([]string, error){
	volsNames := []string{}

	poolExists, err := StoragePoolExists(l, poolName)
	if (err != nil){
		return volsNames, err
	}

	if (poolExists == false){
		return volsNames, nil
	}

	pool, err := l.StoragePoolLookupByName(poolName)
	if (err != nil){
		return volsNames, err
	}

	vols, _, err :=  l.StoragePoolListAllVolumes(pool, 1, 0)
	if (err != nil){
		return volsNames, err
	}

	for _, n := range vols{
		volsNames = append(volsNames, n.Name)
	}

	return volsNames, nil
}