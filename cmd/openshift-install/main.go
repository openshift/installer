package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/stock"
	"github.com/openshift/installer/pkg/destroy"
)

var (
	installConfigCommand   = kingpin.Command("install-config", "Generate the Install Config asset")
	ignitionConfigsCommand = kingpin.Command("ignition-configs", "Generate the Ignition Config assets")
	manifestsCommand       = kingpin.Command("manifests", "Generate the Kubernetes manifests")
	clusterCommand         = kingpin.Command("cluster", "Create an OpenShift cluster")

	destroyCommand = kingpin.Command("destroy-cluster", "Destroy an OpenShift cluster")

	dirFlag  = kingpin.Flag("dir", "assets directory").Default(".").String()
	logLevel = kingpin.Flag("log-level", "log level (e.g. \"debug\")").Default("warn").Enum("debug", "info", "warn", "error", "fatal", "panic")
)

func main() {
	command := kingpin.Parse()
	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		// By definition we should never enter this condition since kingpin should be guarding against incorrect values.
		log.Fatalf("invalid log-level: %v", err)
	}
	log.SetLevel(l)

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
				log.Fatalf("failed to generate asset: %v", err)
				os.Exit(1)
			}

			if err := st.PersistToFile(*dirFlag); err != nil {
				log.Fatalf("failed to write target to disk: %v", err)
				os.Exit(1)
			}
		}
	case destroyCommand.FullCommand():
		destroyer, err := destroy.NewDestroyer(l, *dirFlag)
		if err != nil {
			log.Fatalf("failed to create destroyer: %v", err)
			os.Exit(1)
		}
		if err := destroyer.Run(); err != nil {
			log.Fatalf("destroy failed: %v", err)
			os.Exit(1)
		}

	}

}
