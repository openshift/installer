package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/openshift/installer/installer/pkg/workflow"
)

var (
	clusterInitCommand    = kingpin.Command("init", "Initialize a new Tectonic cluster")
	clusterInitConfigFlag = clusterInitCommand.Flag("config", "Cluster specification file").Required().ExistingFile()

	clusterInstallCommand = kingpin.Command("install", "Create a new Tectonic cluster")
	clusterInstallDirFlag = clusterInstallCommand.Flag("dir", "Cluster directory").Default(".").ExistingDir()

	clusterDestroyCommand   = kingpin.Command("destroy", "Destroy an existing Tectonic cluster")
	clusterDestroyDirFlag   = clusterDestroyCommand.Flag("dir", "Cluster directory").Default(".").ExistingDir()
	clusterDestroyContOnErr = clusterDestroyCommand.Flag("continue-on-error", "Log errors, but attempt to continue cleaning up the cluster.  THIS MAY LEAK RESOURCES, because you may not have enough state left after a partial removal to be able to perform a second destroy.").Default("false").Bool()

	convertCommand    = kingpin.Command("convert", "Convert a tfvars.json to a Tectonic config.yaml")
	convertConfigFlag = convertCommand.Flag("config", "tfvars.json file").Required().ExistingFile()

	logLevel = kingpin.Flag("log-level", "log level (e.g. \"debug\")").Default("info").Enum("debug", "info", "warn", "error", "fatal", "panic")
)

func main() {
	var w workflow.Workflow

	switch kingpin.Parse() {
	case clusterInitCommand.FullCommand():
		w = workflow.InitWorkflow(*clusterInitConfigFlag)
	case clusterInstallCommand.FullCommand():
		w = workflow.InstallWorkflow(*clusterInstallDirFlag)
	case clusterDestroyCommand.FullCommand():
		w = workflow.DestroyWorkflow(*clusterDestroyDirFlag, *clusterDestroyContOnErr)
	case convertCommand.FullCommand():
		w = workflow.ConvertWorkflow(*convertConfigFlag)
	}

	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		// By definition we should never enter this condition since kingpin should be guarding against incorrect values.
		log.Fatalf("invalid log-level: %v", err)
	}
	log.SetLevel(l)

	if err := w.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
