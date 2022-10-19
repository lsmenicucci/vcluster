package pkg

import (
	"net"
	"net/netip"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
	"github.com/lsmenicucci/vcluster/pkg/components"

	"github.com/kr/pretty"
)

func getAddrRange(prefix netip.Prefix) (netip.Addr, netip.Addr){
	start := prefix.Addr()
	end := start

	for {
		if (prefix.Contains(end.Next()) == false){
			return start,end.Prev()
		}
		end = end.Next()
	}
} 

func DialLibvirt() (*libvirt.Libvirt, error){
	log := logrus.New()

	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Errorf("Is libvirt's daemon running?")
		return nil, err
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
		return nil, err
	}

	return l, nil
}

func Test(l *libvirt.Libvirt) {
	cfg := &components.StoragePool{
		Name: "gotest",
		Path: "/home/lsmeni/projetos/simlab-vcluster/test-cluster",
	}

	components.SetStoragePool(l, cfg)
	_, err := components.GetStoragePool(l, "gotest")

	if (err != nil){
		logrus.WithError(err).Info("Error")
	}

	c := &Cluster{}
	c.LoadCluster(l, "gotest")
	logrus.Printf("%# v", pretty.Formatter(c))
}