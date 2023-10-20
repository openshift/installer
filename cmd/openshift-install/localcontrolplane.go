package main

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/cmd/openshift-install/command"
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
