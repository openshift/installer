package softwaredefinedstorage

import (
	"log"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// No Operation SDS Struct Defined
type noopSds struct{}

func NewSdsNoop() Sds {
	return &noopSds{}
}

func (noop noopSds) PreWorkerReplace(worker v2.Worker) error {
	log.Println("In NoopSds PreWorkerReplace")
	return nil
}

func (noop noopSds) PostWorkerReplace(worker v2.Worker) error {
	log.Println("In NoopSds PostWorkerReplace")
	return nil
}
