package pkg

import (
	"github.com/digitalocean/go-libvirt"
//	"github.com/sirupsen/logrus"
)

func NodeExists(l *libvirt.Libvirt, name string) (bool, error){
	domains, _, err :=  l.ConnectListAllDomains(1, 0)
	if (err != nil){
		return false, err
	}

	for _, d := range domains{
		if (name == d.Name){
			return true, nil
		}
	}

	return false, nil
}