package pkg

import (
	"errors"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

func volExists(name string, volList []libvirt.StorageVol) bool{
	for _, v := range volList{
		if (v.Name == name){
			return true
		}
	}

	return false
}

func SetupPool(l *libvirt.Libvirt, c *ClusterConfig, skip_existing bool) error{
	poolName := c.GetStoragePoolName()
	log := logrus.WithField("pool", poolName)

	// check pre existing pools
	log.Debug("Listing available storage pools")
	pools, _ ,err := l.ConnectListAllStoragePools(1, 0)
	if (err != nil){
		return errors.New(fmt.Sprintf("Error fetching storage pools: %e", err))
	}

	for _, p := range pools{
		if (p.Name == poolName){
			log.Info("Pool aready exists")
			if (skip_existing){
				return nil
			}else{
				return errors.New("Pool aready exists")
			}
		}
	}
	
	// write pool
	poolCfg, err := c.BuildPoolXML()
	if (err != nil){
		return err
	} 

	// flag = 1 enables XML validation
	log.Debug("Defining storage pool")
	_, err = l.StoragePoolCreateXML(poolCfg, 1)
	if (err != nil){
		return errors.New(fmt.Sprintf("Error while defining storage pool: %e", err))
	}

	return nil
}

func SetupVolume(l *libvirt.Libvirt, pool libvirt.StoragePool, name string, xml string, skip_existing bool) error{
	log := logrus.WithFields(logrus.Fields{ "pool": pool.Name, "volume": name})

	volumes, _, err := l.StoragePoolListAllVolumes(pool, 1, 0)
	if (err != nil){
		log.WithError(err).Debug("Failed listing volumes")
		return err
	}

	if (volExists(name, volumes) == true){
		log.Info("Volume already exists")
		if (skip_existing){
			return nil
		}else{
			return errors.New("Volume already exists in pool")
		}
	}

	log.Debug("Creating volume from XML")
	_, err = l.StorageVolCreateXML(pool, xml, 0)
	if (err != nil){
		log.Logger.WithError(err).Debug("Failed creating volume from XML")
		return err
	}

	return nil
}

func SetupVolumes(l *libvirt.Libvirt, c *ClusterConfig, skip_existing bool) error{
	poolName := c.GetStoragePoolName()
	log := logrus.WithField("pool", poolName)

	pool, err := l.StoragePoolLookupByName(poolName)
	if (err != nil){
		log.Debug(err)
		return err
	}

	// define controller volume
	log.Debug("Generating controller disk XML")
	diskXml, err := c.BuildControllerDiskXML()
	if (err != nil){
		return err
	}

	err = SetupVolume(l,pool, c.GetControllerDiskName(), diskXml, skip_existing)
	if (err != nil){
		log.WithError(err).Debug("Failed creating controller volume")
		return errors.New("Error creating controller volume: " + err.Error())
	}

	// define worker volumes
	for k := 0; k < c.NumWorkers; k++{
		log.Debug("Generating worker disk XML")
		diskXml, err = c.BuildWorkerDiskXML(k)
		if (err != nil){
			return err
		}

		err = SetupVolume(l,pool, c.GetWorkerDiskName(k), diskXml, skip_existing)
		if (err != nil){
			log.WithError(err).Debug("Failed creating worker volume")
			return errors.New("Error creating worker volume: " + err.Error())
		}
	}

	return nil
}
