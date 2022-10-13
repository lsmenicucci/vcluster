package pkg

import (
	"errors"
	"os"
	"github.com/digitalocean/go-libvirt"
)

func defineNode(l *libvirt.Libvirt, c NodeConfig, overwrite bool) error{
	domains, _, err := l.ConnectListAllDomains(1, 0)
	
	if err != nil {
		return err
	}

	for _, d := range domains{
		if d.Name == c.Name && overwrite == false{
			return errors.New("Domain already exists and overwrite is not intended")
		}
	}

	_ = NodeXML.Execute(os.Stdout, c)

	return nil
}



