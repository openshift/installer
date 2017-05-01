package server

import (
	"encoding/gob"
	"encoding/json"
	"sync"
)

func init() {
	gob.Register(&TectonicMetalChecker{})
}

// TectonicMetalChecker is a serializable StatusChecker for Tectonic bare-metal
// clusters.
type TectonicMetalChecker struct {
	Controllers    []Node
	Workers        []Node
	TectonicDomain string
}

// Status checks the state of etcd, on-host kublets, and Tectonic components
// in the cluster.
func (c TectonicMetalChecker) Status() ([]byte, error) {
	type Kubelet struct {
		Controllers []ServiceStatus `json:"controllers"`
		Workers     []ServiceStatus `json:"workers"`
	}
	type ClusterStatus struct {
		Etcd            []ServiceStatus `json:"etcd"`
		Kubelet         Kubelet         `json:"kubelet"`
		TectonicConsole ServiceStatus   `json:"tectonicConsole"`
	}

	status := &ClusterStatus{
		Etcd: make([]ServiceStatus, len(c.Controllers)),
		Kubelet: Kubelet{
			Controllers: make([]ServiceStatus, len(c.Controllers)),
			Workers:     make([]ServiceStatus, len(c.Workers)),
		},
	}
	client := defaultStatusClient

	var wg sync.WaitGroup
	for i, node := range c.Controllers {
		wg.Add(2)
		go func(i int, node Node) {
			defer wg.Done()
			status.Etcd[i] = EtcdHealth(client, node.Name)
		}(i, node)
		go func(i int, node Node) {
			defer wg.Done()
			status.Kubelet.Controllers[i] = KubeletHealth(client, node.Name)
		}(i, node)
	}
	for i, node := range c.Workers {
		wg.Add(1)
		go func(i int, node Node) {
			defer wg.Done()
			status.Kubelet.Workers[i] = KubeletHealth(client, node.Name)
		}(i, node)
	}
	// Tectonic Console must be externally resolvable
	wg.Add(1)
	go func(endpoint string) {
		defer wg.Done()
		status.TectonicConsole = TectonicConsoleHealth(client, endpoint)
	}(c.TectonicDomain)

	// Wait for health checks to get responses or timeout
	wg.Wait()
	return json.Marshal(status)
}
