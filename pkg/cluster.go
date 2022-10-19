package pkg

import (
	"strings"

	"github.com/digitalocean/go-libvirt"
	"github.com/lsmenicucci/vcluster/pkg/components"
)

type Cluster struct{
	StoragePool *components.StoragePool 	`json:"storage_pool"`
	Volumes 	[]*components.Volume 		`json:"volumes"`
	Networks 	[]*components.Network 		`json:"networks"`
	Nodes 		[]*components.Node 			`json:"nodes"`
}

func (c *Cluster) SetCluster(l *libvirt.Libvirt) error{	
	err := components.SetStoragePool(l, c.StoragePool)
	if (err != nil){
		return err
	}
	
	for _, v := range c.Volumes{
		err := components.SetVolume(l, c.StoragePool.Name, v)
		if (err != nil){
			return err
		}
	}

	for _, n := range c.Networks{
		err := components.SetNetwork(l, n)
		if (err != nil){
			return err
		}
	}

	for _, n := range c.Nodes{
		err := components.SetNode(l, n)
		if (err != nil){
			return err
		}
	}

	return nil
}

func (c *Cluster) LoadCluster(l *libvirt.Libvirt, prefix string) error{
	// load nodes
	nodeNames, err := components.ListNodesNames(l)
	if (err != nil){
		return err
	}

	c.Nodes = []*components.Node{}
	for _, n := range nodeNames{
		if (strings.HasPrefix(n, prefix)){
			node, err := components.GetNode(l, n)
			if (err != nil){
				return err
			}
			c.Nodes = append(c.Nodes, node)
		}
	}

	// load networks
	netNames, err := components.ListNetworksNames(l)
	if (err != nil){
		return err
	}

	c.Networks = []*components.Network{}
	for _, n := range netNames{
		if (strings.HasPrefix(n, prefix)){
			net, err := components.GetNetwork(l, n)
			if (err != nil){
				return err
			}
			c.Networks = append(c.Networks, net)
		}
	}

	// load storage pool
	pool, err := components.GetStoragePool(l, prefix)
	if (err != nil){
		return err
	}
	c.StoragePool = pool

	// load volumes
	if (c.StoragePool != nil){
		volNames, err := components.ListVolumesNames(l, c.StoragePool.Name)
		if (err != nil){
			return err
		}

		c.Volumes = make([]*components.Volume, len(volNames))
		for k, n := range volNames{
			vol, err := components.GetVolume(l, c.StoragePool.Name, n)
			if (err != nil){
				return err
			}
			c.Volumes[k] = vol
		}
	}

	return nil
}