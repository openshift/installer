package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/openpgp"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/openshift/library-go/pkg/verify/store"
	"github.com/openshift/library-go/pkg/verify/store/serial"
	"github.com/openshift/library-go/pkg/verify/util"
)

// Interface performs verification of the provided content. The default implementation
// in this package uses the container signature format defined at https://github.com/containers/image
// to authenticate that a given release image digest has been signed by a trusted party.
type Interface interface {
	// Verify should return nil if the provided release digest has sufficient signatures to be considered
	// valid. It should return an error in all other cases.
	Verify(ctx context.Context, releaseDigest string) error

	// Signatures returns a copy of any cached signatures that have been validated
	// so far. It may return no signatures.
	Signatures() map[string][][]byte

	// Verifiers returns a copy of the verifiers in this payload.
	Verifiers() map[string]openpgp.EntityList

	// AddStore adds additional stores for signature verification.
	AddStore(additionalStore store.Store)
}

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}

type rejectVerifier struct{}

func (rejectVerifier) Verify(ctx context.Context, releaseDigest string) error {
	return fmt.Errorf("verification is not possible")
}

func (rejectVerifier) Signatures() map[string][][]byte {
	return nil
}

func (rejectVerifier) Verifiers() map[string]openpgp.EntityList {
	return nil
}

func (rejectVerifier) AddStore(additionalStore store.Store) {
}

// Reject fails always fails verification.
var Reject Interface = rejectVerifier{}

// maxSignatureSearch prevents unbounded recursion on malicious signature stores (if
// an attacker was able to take ownership of the store to perform DoS on clusters).
const maxSignatureSearch = 10

// validReleaseDigest is a verification rule to filter clearly invalid digests.
var validReleaseDigest = regexp.MustCompile(`^[a-z0-9]+(?:[.+_-][a-z0-9]+)*:[a-zA-Z0-9=_-]+$`)

// releaseVerifier implements a signature intersection operation on a provided release
// digest - all verifiers must have at least one valid signature attesting the release
// digest. If any failure occurs the caller should assume the content is unverified.
type releaseVerifier struct {
	verifiers map[string]openpgp.EntityList

	// store is the store from which release signatures are retrieved.
	store store.Store

	lock           sync.Mutex
	signatureCache map[string][][]byte
}

// NewReleaseVerifier creates a release verifier for the provided inputs.
func NewReleaseVerifier(verifiers map[string]openpgp.EntityList, store store.Store) Interface {
	return &releaseVerifier{
		verifiers: verifiers,
		store:     store,

		signatureCache: make(map[string][][]byte),
	}
}

// Verifiers returns a copy of the verifiers in this payload.
func (v *releaseVerifier) Verifiers() map[string]openpgp.EntityList {
	out := make(map[string]openpgp.EntityList, len(v.verifiers))
	for k, v := range v.verifiers {
		out[k] = v
	}
	return out
}

// String summarizes the verifier for human consumption
func (v *releaseVerifier) String() string {
	var keys []string
	for name := range v.verifiers {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.Grow(256)
	fmt.Fprintf(&builder, "All release image digests must have GPG signatures from")
	if len(keys) == 0 {
		fmt.Fprint(&builder, " <ERROR: no verifiers>")
	}
	for _, name := range keys {
		verifier := v.verifiers[name]
		fmt.Fprintf(&builder, " %s (", name)
		for i, entity := range verifier {
			if i != 0 {
				fmt.Fprint(&builder, ", ")
			}
			if entity.PrimaryKey != nil {
				fmt.Fprintf(&builder, strings.ToUpper(fmt.Sprintf("%x", entity.PrimaryKey.Fingerprint)))
				fmt.Fprint(&builder, ": ")
			}
			count := 0
			for identityName := range entity.Identities {
				if count != 0 {
					fmt.Fprint(&builder, ", ")
				}
				fmt.Fprintf(&builder, "%s", identityName)
				count++
			}
		}
		fmt.Fprint(&builder, ")")
	}
	fmt.Fprintf(&builder, " - will check for signatures in containers/image format at")
	if v.store == nil {
		fmt.Fprint(&builder, " <ERROR: no store>")
	} else {
		fmt.Fprintf(&builder, " %s", v.store)
	}
	return builder.String()
}

// Verify ensures that at least one valid signature exists for an image with digest
// matching release digest in any of the provided locations for all verifiers, or returns
// an error.
func (v *releaseVerifier) Verify(ctx context.Context, releaseDigest string) error {
	if len(v.verifiers) == 0 || v.store == nil {
		return fmt.Errorf("the release verifier is incorrectly configured, unable to verify digests")
	}
	if len(releaseDigest) == 0 {
		return fmt.Errorf("release images that are not accessed via digest cannot be verified")
	}
	if !validReleaseDigest.MatchString(releaseDigest) {
		return fmt.Errorf("the provided release image digest has an invalid format: %q", releaseDigest)
	}

	if v.hasVerified(releaseDigest) {
		return nil
	}

	remaining := make(map[string]openpgp.EntityList, len(v.verifiers))
	for k, v := range v.verifiers {
		remaining[k] = v
	}

	var signedWith [][]byte
	var errs []error
	err := v.store.Signatures(ctx, "", releaseDigest, func(ctx context.Context, signature []byte, errIn error) (done bool, err error) {
		if errIn != nil {
			klog.V(4).Infof("error retrieving signature for %s: %v", releaseDigest, errIn)
			errs = append(errs, fmt.Errorf("%s: %w", time.Now().Format(time.RFC3339), errIn))
			return false, nil
		}
		for k, keyring := range remaining {
			content, _, err := verifySignatureWithKeyring(bytes.NewReader(signature), keyring)
			if err != nil {
				klog.V(4).Infof("keyring %q could not verify signature for %s: %v", k, releaseDigest, err)
				errs = append(errs, fmt.Errorf("%s: %w", time.Now().Format(time.RFC3339), err))
				continue
			}
			if err := verifyAtomicContainerSignature(content, releaseDigest); err != nil {
				klog.V(4).Infof("signature for %s is not valid: %v", releaseDigest, err)
				errs = append(errs, fmt.Errorf("%s: %w", time.Now().Format(time.RFC3339), err))
				continue
			}
			delete(remaining, k)
			signedWith = append(signedWith, signature)
		}
		return len(remaining) == 0, nil
	})
	if err != nil {
		klog.V(4).Infof("Failed to retrieve signatures for %s: %v", releaseDigest, err)
		errs = append(errs, fmt.Errorf("%s: %w", time.Now().Format(time.RFC3339), err))
	}

	if len(remaining) > 0 {
		remainingKeyRings := make([]string, 0, len(remaining))
		for k := range remaining {
			remainingKeyRings = append(remainingKeyRings, k)
		}
		err := &wrapError{
			msg: fmt.Sprintf("unable to verify %s against keyrings: %s", releaseDigest, strings.Join(remainingKeyRings, ", ")),
			err: errors.NewAggregate(errs),
		}
		klog.V(4).Info(err.Error())
		return err
	}

	v.cacheVerification(releaseDigest, signedWith)

	return nil
}

// Signatures returns a copy of any cached signatures that have been validated
// so far. It may return no signatures.
func (v *releaseVerifier) Signatures() map[string][][]byte {
	copied := make(map[string][][]byte)
	v.lock.Lock()
	defer v.lock.Unlock()
	for k, v := range v.signatureCache {
		copied[k] = v
	}
	return copied
}

// AddStore adds additional stores for signature verification.
func (v *releaseVerifier) AddStore(additionalStore store.Store) {
	v.store = &serial.Store{Stores: []store.Store{additionalStore, v.store}}
}

// hasVerified returns true if the digest has already been verified.
func (v *releaseVerifier) hasVerified(releaseDigest string) bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	_, ok := v.signatureCache[releaseDigest]
	return ok
}

const maxSignatureCacheSize = 64

// cacheVerification caches the result of signature check for a digest for later retrieval.
func (v *releaseVerifier) cacheVerification(releaseDigest string, signedWith [][]byte) {
	v.lock.Lock()
	defer v.lock.Unlock()

	if len(signedWith) == 0 || len(releaseDigest) == 0 || v.signatureCache == nil {
		return
	}
	// remove the new entry
	delete(v.signatureCache, releaseDigest)
	// ensure the cache doesn't grow beyond our cap
	for k := range v.signatureCache {
		if len(v.signatureCache) < maxSignatureCacheSize {
			break
		}
		delete(v.signatureCache, k)
	}
	v.signatureCache[releaseDigest] = signedWith
}

type fileStore struct {
	directory string
}

// Signatures reads signatures as "signature-1", "signature-2", etc. out of a digest-based subdirectory.
func (s *fileStore) Signatures(ctx context.Context, name string, digest string, fn store.Callback) error {
	digestPathSegment, err := util.DigestToKeyPrefix(digest, "=")
	if err != nil {
		return err
	}

	base := filepath.Join(s.directory, digestPathSegment, "signature-")
	for i := 1; i < maxSignatureSearch; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		path := base + strconv.Itoa(i)
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			klog.V(4).Infof("unable to load signature: %v", err)
			done, err := fn(ctx, nil, err)
			if done || err != nil {
				return err
			}
			continue
		}
		done, err := fn(ctx, data, nil)
		if done || err != nil {
			return err
		}
	}
	return nil
}

func (s *fileStore) String() string {
	return fmt.Sprintf("file://%s", s.directory)
}

// verifySignatureWithKeyring performs a containers/image verification of the provided signature
// message, checking for the integrity and authenticity of the provided message in r. It will return
// the identity of the signer if successful along with the message contents.
func verifySignatureWithKeyring(r io.Reader, keyring openpgp.EntityList) ([]byte, string, error) {
	md, err := openpgp.ReadMessage(r, keyring, nil, nil)
	if err != nil {
		return nil, "", fmt.Errorf("could not read the message: %v", err)
	}
	if !md.IsSigned {
		return nil, "", fmt.Errorf("not signed")
	}
	content, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, "", err
	}
	if md.SignatureError != nil {
		return nil, "", fmt.Errorf("signature error: %v", md.SignatureError)
	}
	if md.SignedBy == nil {
		return nil, "", fmt.Errorf("invalid signature")
	}
	if md.Signature != nil {
		if md.Signature.SigLifetimeSecs != nil {
			expiry := md.Signature.CreationTime.Add(time.Duration(*md.Signature.SigLifetimeSecs) * time.Second)
			if time.Now().After(expiry) {
				return nil, "", fmt.Errorf("signature expired on %s", expiry)
			}
		}
	} else if md.SignatureV3 == nil {
		return nil, "", fmt.Errorf("unexpected openpgp.MessageDetails: neither Signature nor SignatureV3 is set")
	}

	// follow conventions in containers/image
	return content, strings.ToUpper(fmt.Sprintf("%x", md.SignedBy.PublicKey.Fingerprint)), nil
}

// An atomic container signature has the following schema:
//
//	{
//		"critical": {
//				"type": "atomic container signature",
//				"image": {
//						"docker-manifest-digest": "sha256:817a12c32a39bbe394944ba49de563e085f1d3c5266eb8e9723256bc4448680e"
//				},
//				"identity": {
//						"docker-reference": "docker.io/library/busybox:latest"
//				}
//		},
//		"optional": {
//				"creator": "some software package v1.0.1-35",
//				"timestamp": 1483228800,
//		}
//	}
type signature struct {
	Critical criticalSignature `json:"critical"`
	Optional optionalSignature `json:"optional"`
}

type criticalSignature struct {
	Type     string           `json:"type"`
	Image    criticalImage    `json:"image"`
	Identity criticalIdentity `json:"identity"`
}

type criticalImage struct {
	DockerManifestDigest string `json:"docker-manifest-digest"`
}

type criticalIdentity struct {
	DockerReference string `json:"docker-reference"`
}

type optionalSignature struct {
	Creator   string `json:"creator"`
	Timestamp int64  `json:"timestamp"`
}

// verifyAtomicContainerSignature verifiers that the provided data authenticates the
// specified release digest. If error is returned the provided data does NOT authenticate
// the release digest and the signature must be ignored.
func verifyAtomicContainerSignature(data []byte, releaseDigest string) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	var sig signature
	if err := d.Decode(&sig); err != nil {
		return fmt.Errorf("the signature is not valid JSON: %v", err)
	}
	if sig.Critical.Type != "atomic container signature" {
		return fmt.Errorf("signature is not the correct type")
	}
	if len(sig.Critical.Identity.DockerReference) == 0 {
		return fmt.Errorf("signature must have an identity")
	}
	if sig.Critical.Image.DockerManifestDigest != releaseDigest {
		return fmt.Errorf("signature digest does not match")
	}
	return nil
}
