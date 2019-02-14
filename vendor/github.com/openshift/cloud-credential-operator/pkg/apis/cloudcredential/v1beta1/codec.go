/*
Copyright 2018 The OpenShift Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// NewScheme creates a new Scheme
func NewScheme() (*runtime.Scheme, error) {
	return SchemeBuilder.Build()
}

// AWSProviderCodec is a runtime codec for the AWS provider.
// +k8s:deepcopy-gen=false
type AWSProviderCodec struct {
	encoder runtime.Encoder
	decoder runtime.Decoder
}

// NewCodec creates a serializer/deserializer for the provider configuration
func NewCodec() (*AWSProviderCodec, error) {
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}
	codecFactory := serializer.NewCodecFactory(scheme)
	encoder, err := newEncoder(&codecFactory)
	if err != nil {
		return nil, err
	}
	codec := AWSProviderCodec{
		encoder: encoder,
		decoder: codecFactory.UniversalDecoder(SchemeGroupVersion),
	}
	return &codec, nil
}

// EncodeProvider serializes an object to the provider spec.
func (codec *AWSProviderCodec) EncodeProviderSpec(in runtime.Object) (*runtime.RawExtension, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}

// DecodeProviderSpec deserializes an object from the provider config.
func (codec *AWSProviderCodec) DecodeProviderSpec(providerConfig *runtime.RawExtension, out runtime.Object) (*AWSProviderSpec, error) {
	obj, _, err := codec.decoder.Decode(providerConfig.Raw, nil, out)
	if err != nil {
		return nil, fmt.Errorf("decoding failure: %v", err)
	}
	s, ok := obj.(*AWSProviderSpec)
	if !ok {
		return nil, fmt.Errorf("error casting to AWSProviderSpec")
	}
	return s, nil
}

// EncodeProviderStatus serializes the provider status.
func (codec *AWSProviderCodec) EncodeProviderStatus(in runtime.Object) (*runtime.RawExtension, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}

// DecodeProviderStatus deserializes the provider status.
func (codec *AWSProviderCodec) DecodeProviderStatus(providerStatus *runtime.RawExtension, out runtime.Object) (*AWSProviderStatus, error) {
	if providerStatus != nil {
		obj, _, err := codec.decoder.Decode(providerStatus.Raw, nil, out)
		if err != nil {
			return nil, fmt.Errorf("decoding failure: %v", err)
		}
		s, ok := obj.(*AWSProviderStatus)
		if !ok {
			return nil, fmt.Errorf("error casting to AWSProviderStatus")
		}
		return s, nil
	}
	return &AWSProviderStatus{}, nil
}

func newEncoder(codecFactory *serializer.CodecFactory) (runtime.Encoder, error) {
	serializerInfos := codecFactory.SupportedMediaTypes()
	if len(serializerInfos) == 0 {
		return nil, fmt.Errorf("unable to find any serlializers")
	}
	encoder := codecFactory.EncoderForVersion(serializerInfos[0].Serializer, SchemeGroupVersion)
	return encoder, nil
}
