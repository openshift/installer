package tshelpers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	imageapi "github.com/openshift/api/image/v1"
	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	"github.com/openshift/library-go/pkg/image/dockerv1client"
)

// FakeOCPRegistry creates a very minimal Docker registry for publishing
// a single fake OCP release image in fixed repo, plus a bunch of
// additional images required by the Agent-based installer.
// The registry is configured to provide just the minimal amount of data
// required by the tests.
type FakeOCPRegistry struct {
	mux    *http.ServeMux
	server *httptest.Server

	blobs     map[string][]byte
	manifests map[string][]byte
	tags      map[string]string

	releaseDigest string
	tmpDir        string
}

// NewFakeOCPRegistry creates a new instance of the fake registry.
func NewFakeOCPRegistry(tmpDir string) *FakeOCPRegistry {
	return &FakeOCPRegistry{
		blobs:     make(map[string][]byte),
		manifests: make(map[string][]byte),
		tags:      make(map[string]string),
		tmpDir:    tmpDir,
	}
}

// Start configures the handlers, brings up the local server for the
// registry and pre-load the required data for publishing an OCP
// release image.
func (fr *FakeOCPRegistry) Start() error {
	fr.mux = http.NewServeMux()

	// Ping handler
	fr.mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Docker-Distribution-Api-Version", "registry/2.0")
		if err := json.NewEncoder(w).Encode(make(map[string]interface{})); err != nil {
			log.Println(err)
		}
	})

	// This handler is invoked when retrieving the image manifest
	fr.mux.HandleFunc("/v2/ocp/release/manifests/{digest}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		digest := r.PathValue("digest")
		manifest, found := fr.manifests[digest]
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if _, err := w.Write(manifest); err != nil {
			log.Println(err)
		}
	})

	// Generic blobs handler used to serve both the image config and data
	fr.mux.HandleFunc("/v2/ocp/release/blobs/{digest}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		digest := r.PathValue("digest")
		blob, found := fr.blobs[digest]
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if _, err := w.Write(blob); err != nil {
			log.Println(err)
		}
	})

	// Catch all
	fr.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	err := fr.newTLSServer(fr.mux.ServeHTTP)
	if err != nil {
		return err
	}
	fr.server.StartTLS()

	err = fr.setupReleasePayload()
	if err != nil {
		return err
	}

	return nil
}

func (fr *FakeOCPRegistry) pullSpec(digest string) string {
	return fmt.Sprintf("%s/ocp/release@%s", fr.server.URL[len("https://"):], digest)
}

// ReleasePullspec provides an handy method to get the release pull spec.
func (fr *FakeOCPRegistry) ReleasePullspec() string {
	return fr.pullSpec(fr.releaseDigest)
}

func (fr *FakeOCPRegistry) addTarFile(tw *tar.Writer, name string, data []byte) {
	header := &tar.Header{
		Name: name,
		Mode: 0600,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(header); err != nil {
		log.Println(err)
	}
	if _, err := tw.Write(data); err != nil {
		log.Println(err)
	}
}

// Creates a small ISO but good enough to be processed
// by ABI.
func (fr *FakeOCPRegistry) makeMinimalISO() ([]byte, error) {
	tempDir, err := os.MkdirTemp(fr.tmpDir, "nodejoiner-it")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	files := map[string][]byte{
		"iso/images/ignition.img":       []byte("ignitionimg"),
		"iso/images/pxeboot/initrd.img": []byte("initrdimg"),
		"iso/images/pxeboot/rootfs.img": []byte("rootfsimg"),
		"iso/images/pxeboot/vmlinuz":    []byte("vmlinuz"),
		"iso/images/efiboot.img":        []byte("efibootimg"),
		"iso/boot.catalog":              []byte("bootcatalog"),
		// The following files are required to allow the correct embedding of
		// the kernel args (ie, fips=1).
		"iso/EFI/redhat/grub.cfg": []byte("\n############### COREOS_KARG_EMBED_AREA"),
		"iso/coreos/kargs.json":   []byte(`{"files": [{"path": "EFI/redhat/grub.cfg"}]}`),
	}
	for file, content := range files {
		dir := filepath.Dir(file)
		fullDir := filepath.Join(tempDir, dir)

		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return nil, err
		}
		fullPath := filepath.Join(tempDir, file)
		f, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if _, err = f.Write(content); err != nil {
			return nil, err
		}
	}

	baseIso := filepath.Join(tempDir, "baseiso-nj.iso")
	if err := isoeditor.Create(baseIso, filepath.Join(tempDir, "iso"), "nj"); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(baseIso)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fr *FakeOCPRegistry) setupReleasePayload() error {
	// agent-installer-utils image
	if _, err := fr.PushImage("agent-installer-utils", func(tw *tar.Writer) error {
		// fake agent-tui files
		fr.addTarFile(tw, "usr/bin/agent-tui", []byte("foo-data"))
		fr.addTarFile(tw, "usr/lib64/libnmstate.so.2", []byte("foo-data"))
		return nil
	}); err != nil {
		return err
	}

	// machine-os-images image
	if _, err := fr.PushImage("machine-os-images", func(tw *tar.Writer) error {
		// fake base ISO
		isoData, err := fr.makeMinimalISO()
		if err != nil {
			return err
		}
		fr.addTarFile(tw, "coreos/coreos-x86_64.iso", isoData)
		return nil
	}); err != nil {
		return err
	}

	// release image
	releaseDigest, err := fr.PushImage("release-99.0.0", func(tw *tar.Writer) error {
		// images-references file
		imageReferences := imageapi.ImageStream{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ImageStream",
				APIVersion: "image.openshift.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"io.openshift.build.versions": "99.0.0",
				},
			},
			Spec: imageapi.ImageStreamSpec{},
		}
		for tag, digest := range fr.tags {
			imageReferences.Spec.Tags = append(imageReferences.Spec.Tags, imageapi.TagReference{
				Name: tag,
				From: &corev1.ObjectReference{
					Name: fr.pullSpec(digest),
					Kind: "DockerImage",
				},
			})
		}
		data, err := json.Marshal(&imageReferences)
		if err != nil {
			return err
		}
		fr.addTarFile(tw, "release-manifests/image-references", data)

		// release-metadata file
		type CincinnatiMetadata struct {
			Kind string `json:"kind"`

			Version  string   `json:"version"`
			Previous []string `json:"previous"`
			Next     []string `json:"next,omitempty"`

			Metadata map[string]interface{} `json:"metadata,omitempty"`
		}
		releaseMetadata := CincinnatiMetadata{
			Kind:    "cincinnati-metadata-v0",
			Version: "99.0.0",
		}
		data, err = json.Marshal(releaseMetadata)
		if err != nil {
			return err
		}
		fr.addTarFile(tw, "release-manifests/release-metadata", data)

		return nil
	})
	if err != nil {
		return err
	}
	fr.releaseDigest = releaseDigest

	return nil
}

func (fr *FakeOCPRegistry) newTLSServer(handler http.HandlerFunc) error {
	fr.server = httptest.NewUnstartedServer(handler)
	cert, err := fr.generateSelfSignedCert()
	if err != nil {
		return fmt.Errorf("error configuring server cert: %w", err)
	}
	fr.server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
	}
	return nil
}

// Close shutdowns the fake registry server.
func (fr *FakeOCPRegistry) Close() {
	fr.server.Close()
}

func (fr *FakeOCPRegistry) generateSelfSignedCert() (tls.Certificate, error) {
	// Generate the private key
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}
	// Generate the serial number
	sn, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return tls.Certificate{}, err
	}
	// Create the certificate template
	template := x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"Day2 AddNodes Tester & Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	return tls.X509KeyPair(certPEM, keyPEM)
}

// PushImage adds an image to the registry, storing the content provided into a single layer.
func (fr *FakeOCPRegistry) PushImage(tag string, blobFn func(tw *tar.Writer) error) (string, error) {
	// Create the image config. Just a few fields are required for oc commands.
	config := dockerv1client.DockerImageConfig{
		ID:           uuid.New().String(),
		Architecture: "amd64",
		OS:           "linux",
		Created:      time.Now(),
	}
	configData, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	configDigest := fr.sha(configData)
	fr.blobs[configDigest] = configData

	// Create the image blob data, as a gzipped tar content.
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	err = blobFn(tw)
	if err != nil {
		return "", err
	}
	tw.Close()
	gw.Close()
	blobData := buf.Bytes()
	blobDigest := fr.sha(blobData)
	fr.blobs[blobDigest] = blobData

	// Create the image manifest.
	manifest := dockerv1client.DockerImageManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Config: dockerv1client.Descriptor{
			MediaType: "application/vnd.docker.container.image.v1+json",
			Digest:    configDigest,
		},
		Layers: []dockerv1client.Descriptor{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Digest:    blobDigest,
			},
		},
		Name: "ocp/release",
		Tag:  tag,
	}
	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return "", err
	}
	manifestDigest := fr.sha(manifestData)
	fr.manifests[manifestDigest] = manifestData

	fr.tags[tag] = manifestDigest

	return manifestDigest, nil
}

func (fr *FakeOCPRegistry) sha(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:]))
}
