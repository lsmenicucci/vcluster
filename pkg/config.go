package pkg

import (
	"encoding/json"
	"net/netip"
	"os"
)

type ClusterWorkerConfig struct{
	CPUS 			int 	`json:cpus`
	Memory			int 	`json:memory`
	DiskSize 		int 	`json:disk_size`
}

type ClusterControllerConfig struct{
	ClusterWorkerConfig
	InstallationISO string 	`json:installation_iso`
}

type ClusterNetworkConfig struct{
	Address netip.Prefix `json:address`
}

type ClusterNetworks struct{
	Internal ClusterNetworkConfig `json:internal`
	External ClusterNetworkConfig `json:external`
}

type ClusterConfig struct{
	Prefix 				string 						`json:prefix`
	NumWorkers			int 						`json:workers`
	WorkerConfig 		ClusterWorkerConfig  		`json:worker_config`
	ControllerConfig 	ClusterControllerConfig 	`json:controller_config`
	Networks 			ClusterNetworks 			`json:networks`
}


func Load(filepath string) (*ClusterConfig, error){
	raw, err := os.ReadFile(filepath)

	if (err != nil){
		return nil, err
	}

	c := &ClusterConfig{}
	err = json.Unmarshal(raw, c)
	if (err != nil){
		return nil, err
	}

	return c, nil
}