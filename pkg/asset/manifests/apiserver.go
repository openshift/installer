package manifests

import (
	"context"
	"path"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
)

var APIServerCfgFilename = path.Join(manifestDir, "cluster-apiserver-02-config.yml")

type APIServer struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*APIServer)(nil)

func (*APIServer) Name() string {
	return "APIServer Config"
}

func (*APIServer) Dependencies() []asset.Asset {
	return nil
}

func (a *APIServer) Generate(_ context.Context, _ asset.Parents) error {
	config := &configv1.APIServer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "APIServer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: configv1.APIServerSpec{
			TLSAdherence: configv1.TLSAdherencePolicyStrictAllComponents,
			TLSSecurityProfile: &configv1.TLSSecurityProfile{
				Type:   configv1.TLSProfileModernType,
				Modern: &configv1.ModernTLSProfile{},
			},
		},
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from APIServer", a.Name())
	}

	a.FileList = []*asset.File{
		{
			Filename: APIServerCfgFilename,
			Data:     configData,
		},
	}

	return nil
}

func (a *APIServer) Files() []*asset.File {
	return a.FileList
}

func (a *APIServer) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
