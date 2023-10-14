package main

import (
	"path/filepath"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getLocalControlPlaneClient() (client.Client, error) {
	envtestCfg, err := clientcmd.BuildConfigFromFlags("", filepath.Join(command.RootOpts.Dir, "auth", "envtest.kubeconfig"))
	if err != nil {
		return nil, err
	}
	return client.New(envtestCfg, client.Options{
		Scheme: scheme.Scheme,
	})
}
