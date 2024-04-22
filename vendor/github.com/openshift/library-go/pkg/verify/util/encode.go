package util

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	coreScheme  = runtime.NewScheme()
	coreCodecs  = serializer.NewCodecFactory(coreScheme)
	coreEncoder runtime.Encoder
)

func init() {
	if err := corev1.AddToScheme(coreScheme); err != nil {
		panic(err)
	}
	coreEncoderCodecFactory := serializer.NewCodecFactory(coreScheme)
	coreEncoder = coreEncoderCodecFactory.LegacyCodec(corev1.SchemeGroupVersion)
}

// ReadConfigMap reads config map object from bytes. nil is returned if
// the object cannot be decoded as a config map.
func ReadConfigMap(objBytes []byte) (*corev1.ConfigMap, error) {
	requiredObj, err := runtime.Decode(coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion), objBytes)
	if err != nil {
		return nil, err
	}
	cm, ok := requiredObj.(*corev1.ConfigMap)
	if ok {
		return cm, nil
	}
	return nil, nil
}

// ConfigMapAsBytes returns given config map as bytes.
func ConfigMapAsBytes(cm *corev1.ConfigMap) ([]byte, error) {
	return runtime.Encode(coreEncoder, cm)
}
