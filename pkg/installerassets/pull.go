package installerassets

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/validate"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getPullSecret(ctx context.Context) (data []byte, err error) {
	value := os.Getenv("OPENSHIFT_INSTALL_PULL_SECRET")
	if value != "" {
		err := validate.JSON([]byte(value))
		if err != nil {
			return nil, err
		}
		return []byte(value), nil
	}

	path := os.Getenv("OPENSHIFT_INSTALL_PULL_SECRET_PATH")
	if path != "" {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = validate.JSON(data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "Pull Secret",
			Help:    "The container registry pull secret for this cluster.",
		},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			return validate.JSON([]byte(ans.(string)))
		}),
	}

	var response string
	err = survey.Ask([]*survey.Question{question}, &response)
	if err != nil {
		return nil, errors.Wrap(err, "ask")
	}

	return []byte(response), nil
}

func pullSecretRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/pull.json",
		RebuildHelper: pullSecretRebuilder,
	}

	parents, err := asset.GetParents(ctx, getByName, "pull-secret")
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coreos-pull-secret",
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: parents["pull-secret"].Data,
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}

	asset.Data, err = json.Marshal(secret)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	Defaults["pull-secret"] = getPullSecret
	Rebuilders["manifests/pull.json"] = pullSecretRebuilder
}
