package verify

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"

	"github.com/openshift/library-go/pkg/manifest"
	"github.com/openshift/library-go/pkg/verify/store"
	"github.com/openshift/library-go/pkg/verify/store/configmap"
	"github.com/openshift/library-go/pkg/verify/store/parallel"
	"github.com/openshift/library-go/pkg/verify/store/sigstore"
	"github.com/openshift/library-go/pkg/verify/util"
)

const (
	// ReleaseAnnotationConfigMapVerifier is an annotation set on a config map in the
	// release payload to indicate that this config map controls signing for the payload.
	// Only the first config map within the payload should be used, regardless of whether
	// it has data. See NewFromConfigMapData for more.
	ReleaseAnnotationConfigMapVerifier = "release.openshift.io/verification-config-map"

	// verifierPublicKeyPrefix is the unique portion of the key used within a config map
	// identifying data field containing one or more GPG public keys in ASCII form that
	// must have signed the release image by digest.
	verifierPublicKeyPrefix = "verifier-public-key-"

	// storePrefix is the unique portion of the key used within a config map identifying
	// data field containing a URL (scheme http://, or https://) location that contains
	// signatures.
	storePrefix = "store-"
)

// GetSignaturesAsConfigmap returns the given signatures in a config map. Uses
// util.DigestToKeyPrefix to replace colon with dash when saving digest to config map.
func GetSignaturesAsConfigmap(digest string, signatures [][]byte) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: configmap.NamespaceLabelConfigMap,
			Labels: map[string]string{
				configmap.ReleaseLabelConfigMap: "",
			},
		},
		BinaryData: make(map[string][]byte),
	}
	prefix, err := util.DigestToKeyPrefix(digest, "-")
	if err != nil {
		return nil, err
	}
	cm.Name = prefix
	for i, v := range signatures {
		cm.BinaryData[fmt.Sprintf("%s-%d", prefix, i+1)] = v
	}
	return cm, nil
}

// NewFromManifests fetches the first config map in the manifest list with the correct annotation.
// It returns an error if the data is not valid, or no verifier if a config map wth the required
// annotation is not found. See the verify package for more details on the algorithm for verification.
// If the annotation is set, a verifier or error is always returned.
func NewFromManifests(manifests []manifest.Manifest, clientBuilder sigstore.HTTPClient) (Interface, error) {
	for _, manifest := range manifests {
		configMap, err := util.ReadConfigMap(manifest.Raw)

		// configMap will be nil if this is not a config map
		if err != nil || configMap == nil {
			continue
		}
		if _, ok := configMap.Annotations[ReleaseAnnotationConfigMapVerifier]; !ok {
			continue
		}
		src := fmt.Sprintf("the config map %s/%s", configMap.Namespace, configMap.Name)
		data, _, err := unstructured.NestedStringMap(manifest.Obj.Object, "data")
		if err != nil {
			return nil, errors.Wrapf(err, "%s is not valid: %v", src, err)
		}
		verifier, err := newFromConfigMapData(src, data, clientBuilder)
		if err != nil {
			return nil, err
		}
		return verifier, nil
	}
	return nil, nil
}

// newFromConfigMapData expects to receive the data field of the first config map in the release
// image payload with the annotation "release.openshift.io/verification-config-map". Only the
// first payload item in lexographic order will be considered - all others are ignored. The
// verifier returned by this method
//
// The presence of one or more config maps instructs the CVO to verify updates before they are
// downloaded.
//
// The keys within the config map in the data field define how verification is performed:
//
// verifier-public-key-*: One or more GPG public keys in ASCII form that must have signed the
//
//	release image by digest.
//
// store-*: A URL (scheme file://, http://, or https://) location that contains signatures. These
//
//	signatures are in the atomic container signature format. The URL will have the digest
//	of the image appended to it as "<STORE>/<ALGO>=<DIGEST>/signature-<NUMBER>" as described
//	in the container image signing format. The docker-image-manifest section of the
//	signature must match the release image digest. Signatures are searched starting at
//	NUMBER 1 and incrementing if the signature exists but is not valid. The signature is a
//	GPG signed and encrypted JSON message. The file store is provided for testing only at
//	the current time, although future versions of the CVO might allow host mounting of
//	signatures.
//
// See https://github.com/containers/image/blob/ab49b0a48428c623a8f03b41b9083d48966b34a9/docs/signature-protocols.md
// for a description of the signature store
//
// The returned verifier will require that any new release image will only be considered verified
// if each provided public key has signed the release image digest. The signature may be in any
// store and the lookup order is internally defined.
func newFromConfigMapData(src string, data map[string]string, clientBuilder sigstore.HTTPClient) (Interface, error) {
	verifiers := make(map[string]openpgp.EntityList)
	var stores []store.Store
	for k, v := range data {
		switch {
		case strings.HasPrefix(k, verifierPublicKeyPrefix):
			keyring, err := loadArmoredOrUnarmoredGPGKeyRing([]byte(v))
			if err != nil {
				return nil, errors.Wrapf(err, "%s has an invalid key %q that must be a GPG public key: %v", src, k, err)
			}
			verifiers[k] = keyring
		case strings.HasPrefix(k, storePrefix):
			v = strings.TrimSpace(v)
			u, err := url.Parse(v)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "file") {
				return nil, fmt.Errorf("%s has an invalid key %q: must be a valid URL with scheme file://, http://, or https://", src, k)
			}
			if u.Scheme == "file" {
				stores = append(stores, &fileStore{
					directory: u.Path,
				})
			} else {
				stores = append(stores, &sigstore.Store{
					URI:        u,
					HTTPClient: clientBuilder,
				})
			}
		default:
			klog.Warningf("An unexpected key was found in %s and will be ignored (expected store-* or verifier-public-key-*): %s", src, k)
		}
	}
	if len(stores) == 0 {
		return nil, fmt.Errorf("%s did not provide any signature stores to read from and cannot be used", src)
	}
	if len(verifiers) == 0 {
		return nil, fmt.Errorf("%s did not provide any GPG public keys to verify signatures from and cannot be used", src)
	}

	return NewReleaseVerifier(verifiers, &parallel.Store{Stores: stores}), nil
}

func loadArmoredOrUnarmoredGPGKeyRing(data []byte) (openpgp.EntityList, error) {
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(data))
	if err == nil {
		return keyring, nil
	}
	return openpgp.ReadKeyRing(bytes.NewReader(data))
}
