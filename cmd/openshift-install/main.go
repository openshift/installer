package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/stock"
)

var (
	installConfigCommand = kingpin.Command("install-config", "Generate the Install Config asset")

	dirFlag  = kingpin.Flag("dir", "assets directory").Default(".").String()
	logLevel = kingpin.Flag("log-level", "log level (e.g. \"debug\")").Default("info").Enum("debug", "info", "warn", "error", "fatal", "panic")
)

func main() {
	command := kingpin.Parse()

	assetStock := stock.EstablishStock(*dirFlag)

	var targetAsset asset.Asset

	switch command {
	case installConfigCommand.FullCommand():
		targetAsset = assetStock.InstallConfig()
	}

	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		// By definition we should never enter this condition since kingpin should be guarding against incorrect values.
		log.Fatalf("invalid log-level: %v", err)
	}
	log.SetLevel(l)

	assetStore := &asset.StoreImpl{}
	st, err := assetStore.Fetch(targetAsset)
	if err != nil {
		log.Fatalf("failed to generate asset: %v", err)
		os.Exit(1)
	}

	if err := st.PersistToFile(); err != nil {
		log.Fatalf("failed to write target to disk: %v", err)
		os.Exit(1)
	}
}
