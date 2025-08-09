/*
Copyright 2022 Nutanix

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package context

import (
	"context"
	"fmt"
	"sync"

	"github.com/nutanix-cloud-native/prism-go-client/utils"
	prismclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	ctrl "sigs.k8s.io/controller-runtime"
	ctlclient "sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
)

var (
	RemoteClientCache = map[ctlclient.ObjectKey]ctlclient.Client{}
	cacheLock         = &sync.Mutex{}
)

// ClusterContext is a context used with a NutanixCluster reconciler
type ClusterContext struct {
	Context       context.Context
	NutanixClient *prismclientv3.Client

	Cluster        *capiv1.Cluster
	NutanixCluster *infrav1.NutanixCluster
}

// MachineContext is a context used with a NutanixMachine reconciler
type MachineContext struct {
	Context       context.Context
	NutanixClient *prismclientv3.Client

	Cluster        *capiv1.Cluster
	Machine        *capiv1.Machine
	NutanixCluster *infrav1.NutanixCluster
	NutanixMachine *infrav1.NutanixMachine

	// The VM ip address
	IP string
}

// IsControlPlaneMachine returns true if the provided resource is
// a member of the control plane.
func IsControlPlaneMachine(nma *infrav1.NutanixMachine) bool {
	if nma == nil {
		return false
	}
	_, ok := nma.GetLabels()[capiv1.MachineControlPlaneNameLabel]
	return ok
}

// GetNutanixMachinesInCluster gets a cluster's NutanixMachine resources.
func (clctx *ClusterContext) GetNutanixMachinesInCluster(client ctlclient.Client) ([]*infrav1.NutanixMachine, error) {
	clusterName := clctx.NutanixCluster.Name
	clusterNamespace := clctx.NutanixCluster.Namespace
	labels := map[string]string{capiv1.ClusterNameLabel: clusterName}
	machineList := &infrav1.NutanixMachineList{}

	err := client.List(clctx.Context, machineList,
		ctlclient.InNamespace(clusterNamespace), ctlclient.MatchingLabels(labels))
	if err != nil {
		return nil, err
	}

	ntxMachines := make([]*infrav1.NutanixMachine, len(machineList.Items))
	for i := range machineList.Items {
		ntxMachines[i] = &machineList.Items[i]
	}

	return ntxMachines, nil
}

func (clctx *ClusterContext) SetFailureStatus(failureReason string, failureMessage error) {
	log := ctrl.LoggerFrom(clctx.Context)
	log.Error(failureMessage, fmt.Sprintf("cluster failed: %s", failureReason))
	clctx.NutanixCluster.Status.FailureMessage = utils.StringPtr(fmt.Sprintf("%v", failureMessage))
	clctx.NutanixCluster.Status.FailureReason = &failureReason
}

func (clctx *MachineContext) SetFailureStatus(failureReason string, failureMessage error) {
	log := ctrl.LoggerFrom(clctx.Context)
	log.Error(failureMessage, fmt.Sprintf("machine failed: %s", failureReason))
	clctx.NutanixMachine.Status.FailureMessage = utils.StringPtr(fmt.Sprintf("%v", failureMessage))
	clctx.NutanixMachine.Status.FailureReason = &failureReason
}

func GetRemoteClient(ctx context.Context, client ctlclient.Client, clusterKey ctlclient.ObjectKey) (ctlclient.Client, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	remoteClient, ok := RemoteClientCache[clusterKey]
	if ok {
		return remoteClient, nil
	}
	remoteClient, err := remote.NewClusterClient(ctx, "remote-cluster-cache", client, clusterKey)
	if err != nil {
		return nil, err
	}
	RemoteClientCache[clusterKey] = remoteClient

	return remoteClient, nil
}

func RemoveRemoteClient(clusterKey ctlclient.ObjectKey) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	delete(RemoteClientCache, clusterKey)
}
