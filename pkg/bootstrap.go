package pkg

import (
	"github.com/digitalocean/go-libvirt"
	log "github.com/sirupsen/logrus"
)

func SetupAll(l *libvirt.Libvirt, c *ClusterConfig, skip_existing bool) error{
	log.Info("Setting up storage pool")
	err := SetupPool(l, c, skip_existing)
	if (err != nil){
		return err
	}

	log.Info("Setting up volumes")
	err = SetupVolumes(l, c, skip_existing)
	if (err != nil){
		return err
	}
	
	log.Info("Setting up networks")
	err = SetupNetworks(l, c, skip_existing)
	if (err != nil){
		return err
	}

	return nil
}