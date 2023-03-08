// NOTE: Boilerplate only.  Ignore this file.

// Package v1beta1 contains API Schema definitions for the libvirtproviderconfig v1beta1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig
// +k8s:defaulter-gen=TypeMeta
// +groupName=libvirtproviderconfig.openshift.io
package v1beta1

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"sigs.k8s.io/yaml"

	machinev1 "github.com/openshift/api/machine/v1beta1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "libvirtproviderconfig.openshift.io", Version: "v1beta1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// LibvirtProviderConfigCodec contains encoder/decoder to convert this types from/to serialize data
// +k8s:deepcopy-gen=false
type LibvirtProviderConfigCodec struct {
	encoder runtime.Encoder
	decoder runtime.Decoder
}

// NewScheme creates a new Scheme
func NewScheme() (*runtime.Scheme, error) {
	return SchemeBuilder.Build()
}

// NewCodec returns a encode/decoder for this API
func NewCodec() (*LibvirtProviderConfigCodec, error) {
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}
	codecFactory := serializer.NewCodecFactory(scheme)
	encoder, err := newEncoder(&codecFactory)
	if err != nil {
		return nil, err
	}
	codec := LibvirtProviderConfigCodec{
		encoder: encoder,
		decoder: codecFactory.UniversalDecoder(SchemeGroupVersion),
	}
	return &codec, nil
}

// DecodeFromProviderSpec decodes a serialised ProviderConfig into an object
func (codec *LibvirtProviderConfigCodec) DecodeFromProviderSpec(providerConfig machinev1.ProviderSpec, out runtime.Object) error {
	if providerConfig.Value != nil {
		// TODO(jchaloup): revert back to using `Decode` once installer and mao have started using
		// libvirtprovider apis pivoted under openshift.io api group
		// _, _, err := codec.decoder.Decode(providerConfig.Value.Raw, nil, out)
		// if err != nil {
		if err := yaml.Unmarshal(providerConfig.Value.Raw, out); err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

// EncodeToProviderSpec encodes an object into a serialised ProviderConfig
func (codec *LibvirtProviderConfigCodec) EncodeToProviderSpec(in runtime.Object) (*machinev1.ProviderSpec, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &machinev1.ProviderSpec{
		Value: &runtime.RawExtension{Raw: buf.Bytes()},
	}, nil
}

// EncodeProviderStatus encodes an object into serialised data
func (codec *LibvirtProviderConfigCodec) EncodeProviderStatus(in runtime.Object) (*runtime.RawExtension, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}

	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}

// DecodeProviderStatus decodes a serialised providerStatus into an object
func (codec *LibvirtProviderConfigCodec) DecodeProviderStatus(providerStatus *runtime.RawExtension, out runtime.Object) error {
	if providerStatus != nil {
		// TODO(jchaloup): revert back to using `Decode` once installer and mao have started using
		// libvirtprovider apis pivoted under openshift.io api group
		// _, _, err := codec.decoder.Decode(providerStatus.Raw, nil, out)
		// if err != nil {
		if err := yaml.Unmarshal(providerStatus.Raw, out); err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

func newEncoder(codecFactory *serializer.CodecFactory) (runtime.Encoder, error) {
	serializerInfos := codecFactory.SupportedMediaTypes()
	if len(serializerInfos) == 0 {
		return nil, fmt.Errorf("unable to find any serlializers")
	}
	encoder := codecFactory.EncoderForVersion(serializerInfos[0].Serializer, SchemeGroupVersion)
	return encoder, nil
}
