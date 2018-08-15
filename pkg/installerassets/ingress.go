package installerassets

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/assets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ingressConfigRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/cluster-ingress-02-config.yaml",
		RebuildHelper: ingressConfigRebuilder,
	}

	parents, err := asset.GetParents(ctx, getByName, "base-domain", "cluster-name")
	if err != nil {
		return nil, err
	}

	baseDomain := string(parents["base-domain"].Data)
	clusterName := string(parents["cluster-name"].Data)

	config := &configv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.IngressSpec{
			Domain: fmt.Sprintf("apps.%s.%s", clusterName, baseDomain),
		},
	}

	asset.Data, err = yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	Rebuilders["manifests/cluster-ingress-02-config.yaml"] = ingressConfigRebuilder
}
