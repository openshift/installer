package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/metadata"
	"github.com/openshift/installer/pkg/destroy"
	_ "github.com/openshift/installer/pkg/destroy/libvirt"
	"github.com/openshift/installer/pkg/terraform"
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
		terraformVersion, err := terraform.Version()
		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok && len(exitError.Stderr) > 0 {
				logrus.Error(string(exitError.Stderr))
			}
			logrus.Fatalf("Failed to calculate Terraform version: %v", err)
		}
		fmt.Println(terraformVersion)
		return
	}

	var targetAssets []asset.WritableAsset
	switch command {
	case installConfigCommand.FullCommand():
		targetAssets = []asset.WritableAsset{&installconfig.InstallConfig{}}
	case ignitionConfigsCommand.FullCommand():
		targetAssets = []asset.WritableAsset{
			&bootstrap.Bootstrap{},
			&machine.Master{},
			&machine.Worker{},
		}
	case manifestsCommand.FullCommand():
		targetAssets = []asset.WritableAsset{
			&manifests.Manifests{},
			&manifests.Tectonic{},
		}
	case clusterCommand.FullCommand():
		targetAssets = []asset.WritableAsset{
			&cluster.TerraformVariables{},
			&kubeconfig.Admin{},
			&cluster.Cluster{},
			&metadata.Metadata{},
		}
	}

	switch command {
	case installConfigCommand.FullCommand(),
		ignitionConfigsCommand.FullCommand(),
		manifestsCommand.FullCommand(),
		clusterCommand.FullCommand():
		assetStore := &asset.StoreImpl{}
		for _, a := range targetAssets {
			err := assetStore.Fetch(a)
			if err != nil {
				logrus.Fatalf("Failed to generate %s: %v", a.Name(), err)
			}

			if err := asset.PersistToFile(a, *dirFlag); err != nil {
				logrus.Fatalf("Failed to write asset (%s) to disk: %v", a.Name(), err)
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
