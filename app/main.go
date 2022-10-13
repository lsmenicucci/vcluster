package main 

import (
	"log"
	"net"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/lsmenicucci/simlab-vcluster/pkg"
)

func main(){
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to dial libvirt: %v", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	domains, _, err := l.ConnectListAllDomains(1, 0)
	if err != nil {
		log.Fatalf("failed to retrieve domains: %v", err)
	}

	for _, dom := range domains {
		log.Println(dom.Name)
	}

	ncfg := pkg.NodeConfig{
		Name:"gotest-head",
		Memory: 2,
		Cpus: 4,
		Disk: &pkg.DiskConfig{ Name: "gotest-hed", Pool: "gotest" },
		Cdrom: nil,
		Networks: []pkg.NetworkDeviceConfig{},
	}

	pkg.DefineNode(l, ncfg, true)
}