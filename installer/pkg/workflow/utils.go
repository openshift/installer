package workflow

import (
	"errors"
	"log"
	"time"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func waitForNcg(m *metadata) error {
	kubeconfigPath := m.statePath + kubeConfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	retries := 180
	wait := 10
	for retries > 0 {
		// client will error until api sever is up
		ds, _ := client.DaemonSets("kube-system").Get("ncg")
		log.Printf("Waiting for NCG to be running, this might take a while...")
		if ds.Status.NumberReady >= 1 {
			return nil
		}
		time.Sleep(time.Second * time.Duration(wait))
		retries--
	}
	return errors.New("NCG is not running")
}

func destroyCname(m *metadata) error {
	return runTfCommand(m.statePath, "destroy", "-force", "-state=bootstrap.tfstate", "-target=aws_route53_record.tectonic_ncg", tectonic.FindTemplatesForStep("bootstrap"))
}

func importAutoScalingGroup(m *metadata) error {
	bp := m.statePath
	var err error
	err = runTfCommand(bp, "import", "-state=joining.tfstate", "-config="+tectonic.FindTemplatesForStep("joining"), "aws_autoscaling_group.masters", m.clusterName+"-masters")
	if err != nil {
		return err
	}
	err = runTfCommand(bp, "import", "-state=joining.tfstate", "-config="+tectonic.FindTemplatesForStep("joining"), "aws_autoscaling_group.workers", m.clusterName+"-workers")
	if err != nil {
		return err
	}
	return nil

}
