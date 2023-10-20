/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package app

import (
	"flag"
	"fmt"

	"k8s.io/klog/v2"

	"github.com/Azure/azure-service-operator/v2/internal/version"
)

type Flags struct {
	MetricsAddr          string
	HealthAddr           string
	WebhookPort          int
	WebhookCertDir       string
	EnableLeaderElection bool
	EnableCRDManagement  bool
	CRDPatterns          string // This is a ; delimited string containing a collection of patterns
	PreUpgradeCheck      bool
}

func (f Flags) String() string {
	return fmt.Sprintf(
		"MetricsAddr: %s, HealthAddr: %s, Webhook=Port: %d, WebhookCertDir: %s, EnableLeaderElection: %t, EnableCRDManagement: %t, CRDPatterns: %s, PreUpgradeCheck: %t",
		f.MetricsAddr,
		f.HealthAddr,
		f.WebhookPort,
		f.WebhookCertDir,
		f.EnableLeaderElection,
		f.EnableCRDManagement,
		f.CRDPatterns,
		f.PreUpgradeCheck)
}

func ParseFlags(args []string) (Flags, error) {
	exeName := args[0] + " " + version.BuildVersion
	flagSet := flag.NewFlagSet(exeName, flag.ExitOnError)
	klog.InitFlags(flagSet)

	var metricsAddr string
	var healthAddr string
	var webhookPort int
	var webhookCertDir string
	var enableLeaderElection bool
	var enableCRDmanagement bool
	var crdPatterns string
	var preUpgradeCheck bool

	// default here for 'MetricsAddr' is set to "0", which sets metrics to be disabled if 'metrics-addr' flag is omitted.
	flagSet.StringVar(&metricsAddr, "metrics-addr", "0", "The address the metric endpoint binds to.")
	flagSet.StringVar(&healthAddr, "health-addr", "", "The address the healthz endpoint binds to.")
	flagSet.IntVar(&webhookPort, "webhook-port", 9443, "The port the webhook endpoint binds to.")
	flagSet.StringVar(&webhookCertDir, "webhook-cert-dir", "", "The directory the webhook server's certs are stored.")
	flagSet.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controllers manager. Enabling this will ensure there is only one active controllers manager.")
	flagSet.BoolVar(&enableCRDmanagement, "enable-crd-management", true,
		"Allows the operator to manage all custom resource definitions, including upgrades.")
	flagSet.StringVar(&crdPatterns, "crd-pattern", "", "Install these CRDs. CRDs already in the cluster will also always be upgraded.")
	flagSet.BoolVar(&preUpgradeCheck, "pre-upgrade-check", false,
		"Enable pre upgrade check to check if existing crds contain helm 'keep' policy.")

	flagSet.Parse(args[1:]) //nolint:errcheck

	return Flags{
		MetricsAddr:          metricsAddr,
		HealthAddr:           healthAddr,
		WebhookPort:          webhookPort,
		WebhookCertDir:       webhookCertDir,
		EnableLeaderElection: enableLeaderElection,
		EnableCRDManagement:  enableCRDmanagement,
		CRDPatterns:          crdPatterns,
		PreUpgradeCheck:      preUpgradeCheck,
	}, nil
}
