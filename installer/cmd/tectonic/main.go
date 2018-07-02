package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/coreos/tectonic-installer/installer/pkg/workflow"
)

var (
	clusterInitCommand        = kingpin.Command("init", "Initialize a new Tectonic cluster")
	clusterInitDomainArg      = clusterInitCommand.Arg("domain", "Domain for the new Tectonic cluster").Required().String()
	clusterInitNameFlag       = clusterInitCommand.Flag("name", "Name for the new Tectonic cluster").String()
	clusterInitLicenseFlag    = clusterInitCommand.Flag("license", "License path for the new Tectonic cluster").ExistingFile()
	clusterInitPullSecretFlag = clusterInitCommand.Flag("pullsecret", "PullSecret path for the new Tectonic cluster").ExistingFile()

	clusterImportCommand    = kingpin.Command("import", "Imports a Tectonic cluster from a config")
	clusterImportConfigFlag = clusterImportCommand.Flag("config", "Cluster specification file").ExistingFile()

	clusterInstallCommand          = kingpin.Command("install", "Create a new Tectonic cluster")
	clusterInstallTLSCommand       = clusterInstallCommand.Command("tls", "Generate TLS Certificates.")
	clusterInstallTLSNewCommand    = clusterInstallCommand.Command("newtls", "Generate TLS Certificates, using a new engine (experimental)")
	clusterInstallAssetsCommand    = clusterInstallCommand.Command("assets", "Generate Tectonic assets.")
	clusterInstallBootstrapCommand = clusterInstallCommand.Command("bootstrap", "Create a single bootstrap node Tectonic cluster.")
	clusterInstallFullCommand      = clusterInstallCommand.Command("full", "Create a new Tectonic cluster").Default()
	clusterInstallJoinCommand      = clusterInstallCommand.Command("join", "Create master and worker nodes to join an exisiting Tectonic cluster.")
	clusterInstallDirFlag          = clusterInstallCommand.Flag("dir", "Cluster directory").Default(".").ExistingDir()

	clusterDestroyCommand = kingpin.Command("destroy", "Destroy an existing Tectonic cluster")
	clusterDestroyDirFlag = clusterDestroyCommand.Flag("dir", "Cluster directory").Default(".").ExistingDir()

	convertCommand    = kingpin.Command("convert", "Convert a tfvars.json to a Tectonic config.yaml")
	convertConfigFlag = convertCommand.Flag("config", "tfvars.json file").Required().ExistingFile()

	logLevel = kingpin.Flag("log-level", "log level (e.g. \"debug\")").Default("info").Enum("debug", "info", "warn", "error", "fatal", "panic")
)

func main() {
	var w workflow.Workflow

	switch kingpin.Parse() {
	case clusterInitCommand.FullCommand():
		w = workflow.InitWorkflow(*clusterInitDomainArg, *clusterInitNameFlag, *clusterInitLicenseFlag, *clusterInitPullSecretFlag)
	case clusterImportCommand.FullCommand():
		w = workflow.ImportWorkflow(*clusterImportConfigFlag)
	case clusterInstallFullCommand.FullCommand():
		w = workflow.InstallFullWorkflow(*clusterInstallDirFlag)
	case clusterInstallTLSCommand.FullCommand():
		w = workflow.InstallTLSWorkflow(*clusterInstallDirFlag)
	case clusterInstallTLSNewCommand.FullCommand():
		w = workflow.InstallTLSNewWorkflow(*clusterInstallDirFlag)
	case clusterInstallAssetsCommand.FullCommand():
		w = workflow.InstallAssetsWorkflow(*clusterInstallDirFlag)
	case clusterInstallBootstrapCommand.FullCommand():
		w = workflow.InstallBootstrapWorkflow(*clusterInstallDirFlag)
	case clusterInstallJoinCommand.FullCommand():
		w = workflow.InstallJoinWorkflow(*clusterInstallDirFlag)
	case clusterDestroyCommand.FullCommand():
		w = workflow.DestroyWorkflow(*clusterDestroyDirFlag)
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
