package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/openshift/installer/pkg/asset"
)

var (
	installConfigCommand = kingpin.Command("install-config", "Generate the Install Config asset")

	logLevel = kingpin.Flag("log-level", "log level (e.g. \"debug\")").Default("info").Enum("debug", "info", "warn", "error", "fatal", "panic")
)

func main() {
	assetStock := asset.EstablishStock()

	var targetAsset asset.Asset

	switch kingpin.Parse() {
	case installConfigCommand.FullCommand():
		targetAsset = assetStock.InstallConfig
	}

	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		// By definition we should never enter this condition since kingpin should be guarding against incorrect values.
		log.Fatalf("invalid log-level: %v", err)
	}
	log.SetLevel(l)

	assetStore := &asset.StoreImpl{}
	assetState, err := assetStore.Fetch(targetAsset)
	if err != nil {
		log.Fatalf("failed to generate asset: %v", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("***** ASSET *****")
	for i, c := range assetState.Contents {
		fmt.Printf("*** Content %v ***\n", i)
		fmt.Println(string(c.Data))
	}
}
