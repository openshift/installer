package main

import (
	"log"
	"os"

	"github.com/coreos/tectonic-installer/installer/pkg/workflow"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	clusterInitCommand    = kingpin.Command("init", "Initialize a new Tectonic cluster")
	clusterInitConfigFlag = clusterInitCommand.Flag("config", "Cluster specification file").Required().ExistingFile()

	clusterInstallCommand          = kingpin.Command("install", "Create a new Tectonic cluster")
	clusterInstallAssetsCommand    = clusterInstallCommand.Command("assets", "Generate Tectonic assets.")
	clusterInstallBootstrapCommand = clusterInstallCommand.Command("bootstrap", "Create a single bootstrap node Tectonic cluster.")
	clusterInstallFullCommand      = clusterInstallCommand.Command("full", "Create a new Tectonic cluster").Default()
	clusterInstallJoinCommand      = clusterInstallCommand.Command("join", "Create master and worker nodes to join an exisiting Tectonic cluster.")
	clusterInstallDirFlag          = clusterInstallCommand.Flag("dir", "Cluster directory").Default(".").ExistingDir()

	clusterDestroyCommand = kingpin.Command("destroy", "Destroy an existing Tectonic cluster")
	clusterDestroyDirFlag = clusterDestroyCommand.Arg("dir", "Cluster directory").Default(".").ExistingDir()
)

func main() {
	var w workflow.Workflow

	switch kingpin.Parse() {
	case clusterInitCommand.FullCommand():
		w = workflow.NewInitWorkflow(*clusterInitConfigFlag)
	case clusterInstallFullCommand.FullCommand():
		w = workflow.NewInstallFullWorkflow(*clusterInstallDirFlag)
	case clusterInstallAssetsCommand.FullCommand():
		w = workflow.NewInstallAssetsWorkflow(*clusterInstallDirFlag)
	case clusterInstallBootstrapCommand.FullCommand():
		w = workflow.NewInstallBootstrapWorkflow(*clusterInstallDirFlag)
	case clusterInstallJoinCommand.FullCommand():
		w = workflow.NewInstallJoinWorkflow(*clusterInstallDirFlag)
	case clusterDestroyCommand.FullCommand():
		w = workflow.NewDestroyWorkflow(*clusterDestroyDirFlag)
	}

	if err := w.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
