package main

import (
	"log"

	"github.com/coreos/tectonic-installer/installer/pkg/workflow"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	dryRunFlag                = kingpin.Flag("dry-run", "Just pretend, but don't do anything").Bool()
	clusterInstallCommand     = kingpin.Command("install", "Create a new Tectonic cluster")
	clusterFullInstallCommand = clusterInstallCommand.Command("full", "Create a new Tectonic cluster").Default()
	clusterAssetsCommand      = clusterInstallCommand.Command("assets", "Generate Tectonic assets.")
	clusterBootstrapCommand   = clusterInstallCommand.Command("bootstrap", "Create a single bootstrap node Tectonic cluster.")
	clusterJoinCommand        = clusterInstallCommand.Command("join", "Create master and worker nodes to join an exisiting Tectonic cluster.")
	clusterDeleteCommand      = kingpin.Command("delete", "Delete an existing Tectonic cluster")
	deleteClusterDir          = clusterDeleteCommand.Arg("dir", "The name of the cluster to delete").String()
	clusterConfigFlag         = clusterInstallCommand.Flag("config", "Cluster specification file").Required().ExistingFile()
)

func main() {
	// TODO: actually do proper error handling
	switch kingpin.Parse() {
	case clusterFullInstallCommand.FullCommand():
		{
			w := workflow.NewInstallWorkflow(*clusterConfigFlag)
			if err := w.Execute(); err != nil {
				log.Fatal(err)
			}
		}
	case clusterAssetsCommand.FullCommand():
		{
			w := workflow.NewAssetsWorkflow(*clusterConfigFlag)
			if err := w.Execute(); err != nil {
				log.Fatal(err)
			}
		}
	case clusterBootstrapCommand.FullCommand():
		{
			w := workflow.NewBootstrapWorkflow(*clusterConfigFlag)
			if err := w.Execute(); err != nil {
				log.Fatal(err)
			}
		}
	case clusterJoinCommand.FullCommand():
		{
			w := workflow.NewJoinWorkflow(*clusterConfigFlag)
			if err := w.Execute(); err != nil {
				log.Fatal(err)
			}
		}
	case clusterDeleteCommand.FullCommand():
		{
			w := workflow.NewDestroyWorkflow(*deleteClusterDir)
			if err := w.Execute(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
