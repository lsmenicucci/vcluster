package pkg

import (
	"errors"

	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

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

func SetupNetwork(l *libvirt.Libvirt, name string, xml string, skip_existing bool) error{
	log := logrus.WithField("network", name)

	exists, err :=  NetworkExists(l, name)
	if (err != nil){
		log.WithError(err).Debug("Could check if network exists")
		return err
	}

	if (exists){
		log.Info("Network already exists")
		if (skip_existing){
			log.Debug("Skipping network creation")
			return nil
		}else{
			return errors.New("Network already exists")
		}
	}

	log.Debug("Creating network")
	log.Debug("Network XML:\n"+xml)
	_, err = l.NetworkCreateXML(xml)
	if (err != nil){
		log.WithError(err).Debug("Could not create network")
		return err
	}

	return nil
}

func SetupNetworks(l *libvirt.Libvirt, c *ClusterConfig, skip_existing bool) error{
	log := logrus.New()

	// Create internal network
	inetName := c.GetInternalNetworkName()
	inetCfg, err := c.BuildInternalNetworkXML()
	if (err != nil){
		log.WithError(err).Debug("Failed generating internal network XML")
		return err
	}

	err = SetupNetwork(l, inetName, inetCfg, skip_existing)
	if (err != nil){
		log.WithError(err).Debug("Failed creating internal network")
		return err
	}

	// Create external network
	enetName := c.GetExternalNetworkName()
	enetCfg, err := c.BuildExternalNetworkXML(map[string]string{})
	if (err != nil){
		log.WithError(err).Debug("Failed generating external network XML")
		return err
	}

	err = SetupNetwork(l, enetName, enetCfg, skip_existing)
	if (err != nil){
		log.WithError(err).Debug("Failed creating external network")
		return err
	}

	return nil
}
