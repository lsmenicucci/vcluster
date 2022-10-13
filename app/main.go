package main

import (
	"net"
	"net/netip"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/lsmenicucci/simlab-vcluster/pkg"
	"github.com/sirupsen/logrus"
)

func main(){
	log := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)

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

	cluster := &pkg.ClusterConfig{
		Prefix: "gotest",
		BaseDir: "/home/lsmeni/projetos/simlab-vcluster/test",
		NumWorkers: 2,
		Worker: pkg.ClusterWorkers{ CPUS: 1, Memory: 2, DiskSize: 1 },
		Controller: pkg.ClusterController{ CPUS: 1, Memory: 2, DiskSize: 1, InstallationISO: ""},
		Networks: pkg.ClusterNetworks{
			Internal: pkg.ClusterNetworkConfig{ Address: netip.MustParsePrefix("10.0.2.1/24") },
			External: pkg.ClusterNetworkConfig{ Address: netip.MustParsePrefix("10.0.2.1/24") },
		},
	}

	err = pkg.SetupPool(l, cluster, true)
	if (err != nil){
		log.Error(err)
	}
	err = pkg.SetupVolumes(l, cluster, true)
	if (err != nil){
		log.Error(err)
	}
}