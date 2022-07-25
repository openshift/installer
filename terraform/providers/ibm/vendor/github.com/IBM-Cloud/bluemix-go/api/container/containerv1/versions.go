package containerv1

import (
	"github.com/IBM-Cloud/bluemix-go/client"
)

//KubeVersion ...
type KubeVersion struct {
	Major   int
	Minor   int
	Patch   int
	Default bool
}

type V1Version map[string][]KubeVersion

//KubeVersions interface
type KubeVersions interface {
	List(target ClusterTargetHeader) ([]KubeVersion, error)
	ListV1(target ClusterTargetHeader) (V1Version, error)
}

type version struct {
	client *client.Client
}

func newKubeVersionAPI(c *client.Client) KubeVersions {
	return &version{
		client: c,
	}
}

//List ...
func (v *version) List(target ClusterTargetHeader) ([]KubeVersion, error) {
	versions := []KubeVersion{}
	_, err := v.client.Get("/v1/kube-versions", &versions, target.ToMap())
	if err != nil {
		return nil, err
	}
	return versions, err
}

func (v *version) ListV1(target ClusterTargetHeader) (V1Version, error) {
	v1ver := V1Version{}
	_, err := v.client.Get("/v1/versions", &v1ver, target.ToMap())
	if err != nil {
		return nil, err
	}

	return v1ver, err
}
