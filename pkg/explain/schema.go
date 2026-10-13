package explain

import (
	"github.com/pkg/errors"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func loadSchema(b []byte) (*apiextv1.JSONSchemaProps, error) {
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	apiextv1.AddToScheme(scheme)
	obj, err := runtime.Decode(codecs.UniversalDecoder(apiextv1.SchemeGroupVersion), b)
	if err != nil {
		return nil, err
	}

	crd, ok := obj.(*apiextv1.CustomResourceDefinition)
	if !ok {
		return nil, errors.Errorf("invalid object, should be *apiextv1.CustomResourceDefinition but found %T", obj)
	}
	if len(crd.Spec.Versions) != 1 {
		return nil, errors.New("missing versions in CRD")
	}
	return crd.Spec.Versions[0].Schema.OpenAPIV3Schema, nil
}
