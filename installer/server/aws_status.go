package server

import (
	"encoding/gob"
	"encoding/json"
	"sync"

	"github.com/coreos/tectonic-installer/installer/server/aws/cloudforms"
)

func init() {
	gob.Register(&TectonicAWSChecker{})
}

// TectonicAWSChecker is a serializable StatusChecker for Tectonic AWS
// clusters.
type TectonicAWSChecker struct {
	AccessKeyID      string
	SecretAccessKey  string
	SessionToken     string
	Region           string
	ControllerDomain string
	TectonicDomain   string
	Cluster          *cloudforms.Cluster
}

// Status checks the state of AWS infrastructure, on-host kubelets, and
// Tectonic components in the cluster.
func (c TectonicAWSChecker) Status() ([]byte, error) {
	type ClusterStatus struct {
		CloudFormation  cloudFormationStatus `json:"cloudFormation"`
		Kubelet         ServiceStatus        `json:"kubelet"`
		TectonicConsole ServiceStatus        `json:"tectonicConsole"`
	}

	status := new(ClusterStatus)
	client := defaultStatusClient

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		// create an AWS client Session
		sess, err := getAWSSession(c.AccessKeyID, c.SecretAccessKey, c.SessionToken, c.Region)
		if err == nil {
			status.CloudFormation = cloudFormationHealth(c.Cluster, sess)
		} else {
			status.CloudFormation = cloudFormationStatus{
				Message: err.Error(),
				Error:   true,
			}
		}
	}()
	go func() {
		defer wg.Done()
		status.Kubelet = KubeletHealth(client, c.ControllerDomain)
	}()
	go func() {
		defer wg.Done()
		status.TectonicConsole = TectonicConsoleHealth(client, c.TectonicDomain)
	}()

	// Wait for health checks to get responses or timeout
	wg.Wait()
	return json.Marshal(status)
}
