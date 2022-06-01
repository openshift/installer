package agent

import (
	"github.com/sirupsen/logrus"

	ainstaller_ctrl "github.com/openshift/assisted-installer/src/assisted_installer_controller/"
	ainstaller_common "github.com/openshift/assisted-installer/src/common"
	ainstaller_invclient "github.com/openshift/assisted-installer/src/inventory_client"
	ainstaller_utils "github.com/openshift/assisted-installer/src/utils"
	aservice_models "github.com/openshift/assisted-service/models"
)

type WaitForBootstrap struct {
	log *logrus.Logger
	inv inventory_client.InventoryClient
	kc  k8s_client.K8SClient
}

// WaitFor wait for the installation complete triggered by the agent installer.
func WaitFor() error {
	logrus.Info("WaitFor command")

	return nil
}

func WaitForBootstrap(b *WaitForBootstrap) error {
	logrus.Info("WaitForBootstrap command")

	// Assisted Service Models
	// - cluster : https://github.com/openshift/assisted-service/blob/master/models/cluster.go
	// - host: https://github.com/openshift/assisted-service/blob/master/models/host.go
	// - host_stage: https://github.com/openshift/assisted-service/blob/master/models/host.go

	// Assisted Installer functions
	// - invclient.GetHosts
	// - common.GetHostsInStatus
	// - utils.GenerateRequestContext

	// Openshift Installer
	// - waitForBootstrapComplete cmd/openshift-install/create.go#341

	// psuedo --

	// while bootstrap != not complete

	// 	check for assisted service
	// 	if assisted service is up log host status from assisted service ?

	// after bootstrap complete return true

	// Snippets from assisted installer ---

	// ctxReq := ainstaller_utils.GenerateRequestContext()
	// log := ainstaller_utils.RequestIDLogger(ctxReq, c.log)

	// nodesMap, err := ainstaller_invclient.GetHosts(ctxReq, log, []string{aservice_models.HostStatusDisabled})
	// hostsInProgressMap := ainstaller_common.GetHostsInStatus(nodesMap, []string{aservice_models.HostStatusInstalled}, false)
	// errNodesMap := ainstaller_common.GetHostsInStatus(nodesMap, []string{aservice_models.HostStatusError}, true)

	// ainstaller_ctrl.logHostStatus(log, nodesMap)

	return nil
}
