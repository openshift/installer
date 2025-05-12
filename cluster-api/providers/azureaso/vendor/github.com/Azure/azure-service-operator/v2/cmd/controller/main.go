/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package main

import (
	"flag"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/Azure/azure-service-operator/v2/cmd/controller/app"
	"github.com/Azure/azure-service-operator/v2/cmd/controller/logging"
	"github.com/Azure/azure-service-operator/v2/internal/version"
)

func main() {
	// Set up to parse command line flags
	exeName := os.Args[0] + " " + version.BuildVersion
	flagSet := flag.NewFlagSet(exeName, flag.ExitOnError)

	// Create a temporary logger for while we get set up
	log := logging.Create(&logging.Config{})

	ctx := ctrl.SetupSignalHandler()

	// Add application and logging flags
	appFlags := app.InitFlags(flagSet)
	logFlags := logging.InitFlags(flagSet)
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Error(err, "failed to parse cmdline flags")
		os.Exit(1)
	}

	// Replace the logger with a configured one
	log = logging.Create(logFlags)
	ctrl.SetLogger(log)
	log.Info("Launching with flags", "flags", appFlags.String())

	mgr := app.SetupControllerManager(ctx, log, appFlags)
	log.Info("starting manager")
	if err = mgr.Start(ctx); err != nil {
		log.Error(err, "failed to start manager")
		os.Exit(1)
	}
}
