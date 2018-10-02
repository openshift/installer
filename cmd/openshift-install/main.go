package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/stock"
	"github.com/openshift/installer/pkg/destroy"
	_ "github.com/openshift/installer/pkg/destroy/libvirt"
)

var (
	installConfigCommand   = kingpin.Command("install-config", "Generate the Install Config asset")
	ignitionConfigsCommand = kingpin.Command("ignition-configs", "Generate the Ignition Config assets")
	manifestsCommand       = kingpin.Command("manifests", "Generate the Kubernetes manifests")
	clusterCommand         = kingpin.Command("cluster", "Create an OpenShift cluster")
	versionCommand         = kingpin.Command("version", "Print version information and exit")

	destroyCommand = kingpin.Command("destroy-cluster", "Destroy an OpenShift cluster")

	dirFlag  = kingpin.Flag("dir", "assets directory").Default(".").String()
	logLevel = kingpin.Flag("log-level", "log level (e.g. \"debug\")").Default("info").Enum("debug", "info", "warn", "error")

	version = "was not built correctly" // set in hack/build.sh
)

func main() {
	command := kingpin.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
	if level, err := logrus.ParseLevel(*logLevel); err == nil {
		logrus.SetLevel(level)
	} else {
		// By definition we should never enter this condition since kingpin should be guarding against incorrect values.
		logrus.Panicf("Invalid log-level: %v", err)
	}

	if command == versionCommand.FullCommand() {
		fmt.Printf("%s %s\n", os.Args[0], version)
		return
	}

	assetStock := stock.EstablishStock()

	var targetAssets []asset.Asset
	switch command {
	case installConfigCommand.FullCommand():
		targetAssets = []asset.Asset{assetStock.InstallConfig()}
	case ignitionConfigsCommand.FullCommand():
		targetAssets = []asset.Asset{
			assetStock.BootstrapIgnition(),
			assetStock.MasterIgnition(),
			assetStock.WorkerIgnition(),
		}
	case manifestsCommand.FullCommand():
		targetAssets = []asset.Asset{
			assetStock.Manifests(),
			assetStock.Tectonic(),
		}
	case clusterCommand.FullCommand():
		targetAssets = []asset.Asset{
			assetStock.TFVars(),
			assetStock.KubeconfigAdmin(),
			assetStock.Cluster(),
			assetStock.Metadata(),
		}
	}

	switch command {
	case installConfigCommand.FullCommand(),
		ignitionConfigsCommand.FullCommand(),
		manifestsCommand.FullCommand(),
		clusterCommand.FullCommand():
		assetStore := &asset.StoreImpl{}
		for _, asset := range targetAssets {
			st, err := assetStore.Fetch(asset)
			if err != nil {
				logrus.Fatalf("Failed to generate asset: %v", err)
			}

			if err := st.PersistToFile(*dirFlag); err != nil {
				logrus.Fatalf("Failed to write asset (%s) to disk: %v", asset.Name(), err)
			}
		}
	case destroyCommand.FullCommand():
		destroyer, err := destroy.New(logrus.StandardLogger(), *dirFlag)
		if err != nil {
			logrus.Fatalf("Failed while preparing to destroy cluster: %v", err)
		}
		if err := destroyer.Run(); err != nil {
			logrus.Fatalf("Failed to destroy cluster: %v", err)
		}
	}
}
