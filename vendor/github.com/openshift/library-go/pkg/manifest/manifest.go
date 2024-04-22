package manifest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	CapabilityAnnotation  = "capability.openshift.io/name"
	DefaultClusterProfile = "self-managed-high-availability"
	featureSetAnnotation  = "release.openshift.io/feature-set"
)

var knownFeatureSets = sets.Set[string]{}

func init() {
	for _, featureSets := range configv1.AllFeatureSets() {
		for featureSet := range featureSets {
			if len(featureSet) == 0 {
				knownFeatureSets.Insert("Default")
				continue
			}
			knownFeatureSets.Insert(string(featureSet))
		}
	}
}

// resourceId uniquely identifies a Kubernetes resource.
// It is used to identify any duplicate resources within
// a given set of manifests.
type resourceId struct {
	// Group identifies a set of API resources exposed together.
	// +optional
	Group string
	// Kind is the name of a particular object schema.
	Kind string
	// Name, sometimes used with the optional Namespace, helps uniquely identify an object.
	Name string
	// Namespace helps uniquely identify an object.
	// +optional
	Namespace string
}

// Manifest stores Kubernetes object in Raw from a file.
// It stores the id and the GroupVersionKind for
// the manifest. Raw and Obj should always be kept in sync
// such that each provides the same data but in different
// formats. To ensure Raw and Obj are always in sync, they
// should not be set directly but rather only be set by
// calling either method ManifestsFromFiles or
// ParseManifests.
type Manifest struct {
	// OriginalFilename is set to the filename this manifest was loaded from.
	// It is not guaranteed to be set or be unique, but will be set when
	// loading from disk to provide a better debug capability.
	OriginalFilename string

	id resourceId

	Raw []byte
	GVK schema.GroupVersionKind

	Obj *unstructured.Unstructured
}

func (r resourceId) equal(id resourceId) bool {
	return reflect.DeepEqual(r, id)
}

func (r resourceId) String() string {
	if len(r.Namespace) == 0 {
		return fmt.Sprintf("Group: %q Kind: %q Name: %q", r.Group, r.Kind, r.Name)
	} else {
		return fmt.Sprintf("Group: %q Kind: %q Namespace: %q Name: %q", r.Group, r.Kind, r.Namespace, r.Name)
	}
}

func (m *Manifest) String() string {
	if m == nil {
		return "nil pointer manifest"
	}

	if m.OriginalFilename != "" {
		return fmt.Sprintf("Filename: %q %s", m.OriginalFilename, m.id.String())
	}

	return m.id.String()
}

func (m Manifest) SameResourceID(manifest Manifest) bool {
	return m.id.equal(manifest.id)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Manifest
// type. It unmarshals bytes of a single kubernetes object to Manifest.
func (m *Manifest) UnmarshalJSON(in []byte) error {
	if m == nil {
		return errors.New("Manifest: UnmarshalJSON on nil pointer")
	}

	// This happens when marshalling
	// <yaml>
	// ---	(this between two `---`)
	// ---
	// <yaml>
	if bytes.Equal(in, []byte("null")) {
		m.Raw = nil
		return nil
	}

	m.Raw = append(m.Raw[0:0], in...)
	udi, _, err := scheme.Codecs.UniversalDecoder().Decode(in, nil, &unstructured.Unstructured{})
	if err != nil {
		return errors.Wrapf(err, "unable to decode manifest")
	}
	ud, ok := udi.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("expected manifest to decode into *unstructured.Unstructured, got %T", ud)
	}
	m.Obj = ud
	return m.populateFromObj()
}

func (m *Manifest) populateFromObj() error {
	m.GVK = m.Obj.GroupVersionKind()
	m.id = resourceId{
		Group:     m.GVK.Group,
		Kind:      m.GVK.Kind,
		Namespace: m.Obj.GetNamespace(),
		Name:      m.Obj.GetName(),
	}
	return validateResourceId(m.id)
}

func getFeatureSets(annotations map[string]string) (sets.Set[string], bool, error) {
	ret := sets.Set[string]{}
	specified := false
	for _, featureSetAnnotation := range []string{featureSetAnnotation} {
		featureSetAnnotationValue, featureSetAnnotationExists := annotations[featureSetAnnotation]
		if featureSetAnnotationExists {
			specified = true
			featureSetAnnotationValues := strings.Split(featureSetAnnotationValue, ",")
			for _, manifestFeatureSet := range featureSetAnnotationValues {
				if !knownFeatureSets.Has(manifestFeatureSet) {
					// never include the manifest if the feature-set annotation is outside of known values
					return nil, specified, fmt.Errorf("unrecognized value %q in %s=%s; known values are: %v", manifestFeatureSet, featureSetAnnotation, featureSetAnnotationValue, strings.Join(sets.List(knownFeatureSets), ","))
				}
			}
			ret.Insert(featureSetAnnotationValues...)
		}
	}

	return ret, specified, nil
}

func checkFeatureSets(requiredFeatureSet string, annotations map[string]string) error {
	requiredAnnotationValue := requiredFeatureSet
	if len(requiredFeatureSet) == 0 {
		requiredAnnotationValue = "Default" // "" in the FeatureSet API is "Default" in the annotation value
	}
	manifestFeatureSets, manifestSpecifiesFeatureSets, err := getFeatureSets(annotations)
	if err != nil {
		return err
	}
	if manifestSpecifiesFeatureSets && !manifestFeatureSets.Has(requiredAnnotationValue) {
		return fmt.Errorf("%q is required, and %s=%s", requiredAnnotationValue, featureSetAnnotation, strings.Join(sets.List(manifestFeatureSets), ","))
	}

	return nil
}

// Include returns an error if the manifest fails an inclusion filter and should be excluded from further
// processing by cluster version operator. Pointer arguments can be set nil to avoid excluding based on that
// filter. For example, setting profile non-nil and capabilities nil will return an error if the manifest's
// profile does not match, but will never return an error about capability issues.
func (m *Manifest) Include(excludeIdentifier *string, requiredFeatureSet *string, profile *string, capabilities *configv1.ClusterVersionCapabilitiesStatus, overrides []configv1.ComponentOverride) error {
	return m.IncludeAllowUnknownCapabilities(excludeIdentifier, requiredFeatureSet, profile, capabilities, overrides, false)
}

// IncludeAllowUnknownCapabilities returns an error if the manifest fails an inclusion filter and should be excluded from
// further processing by cluster version operator. Pointer arguments can be set nil to avoid excluding based on that
// filter. For example, setting profile non-nil and capabilities nil will return an error if the manifest's
// profile does not match, but will never return an error about capability issues. allowUnknownCapabilities only applies
// to capabilities filtering. When set to true a manifest will not be excluded simply because it contains an unknown
// capability. This is necessary to allow updates to an OCP version containing newly defined capabilities.
func (m *Manifest) IncludeAllowUnknownCapabilities(excludeIdentifier *string, requiredFeatureSet *string, profile *string,
	capabilities *configv1.ClusterVersionCapabilitiesStatus, overrides []configv1.ComponentOverride, allowUnknownCapabilities bool) error {

	annotations := m.Obj.GetAnnotations()
	if annotations == nil {
		return fmt.Errorf("no annotations")
	}

	if excludeIdentifier != nil {
		excludeAnnotation := fmt.Sprintf("exclude.release.openshift.io/%s", *excludeIdentifier)
		if v := annotations[excludeAnnotation]; v == "true" {
			return fmt.Errorf("%s=%s", excludeAnnotation, v)
		}
	}

	if requiredFeatureSet != nil {
		err := checkFeatureSets(*requiredFeatureSet, annotations)
		if err != nil {
			return err
		}
	}

	if profile != nil {
		profileAnnotation := fmt.Sprintf("include.release.openshift.io/%s", *profile)
		if val, ok := annotations[profileAnnotation]; ok && val != "true" {
			return fmt.Errorf("unrecognized value %s=%s", profileAnnotation, val)
		} else if !ok {
			return fmt.Errorf("%s unset", profileAnnotation)
		}
	}

	// If there is no capabilities defined in a release then we do not need to check presence of capabilities in the manifest
	if capabilities != nil {
		err := checkResourceEnablement(annotations, capabilities, allowUnknownCapabilities)
		if err != nil {
			return err
		}
	}

	if override := m.getOverrideForManifest(overrides); override != nil && override.Unmanaged {
		return fmt.Errorf("overridden")
	}

	return nil
}

// getOverrideForManifest returns the override when override exists and nil otherwise.
func (m *Manifest) getOverrideForManifest(overrides []configv1.ComponentOverride) *configv1.ComponentOverride {
	for _, override := range overrides {
		namespace := override.Namespace
		if m.id.Namespace == "" {
			namespace = "" // cluster-scoped objects don't have namespace.
		}
		if m.id.equal(resourceId{
			Group:     override.Group,
			Kind:      override.Kind,
			Name:      override.Name,
			Namespace: namespace,
		}) {
			return &override
		}
	}
	return nil
}

// checkResourceEnablement, given resource annotations and defined cluster capabilities, checks if the capability
// annotation exists. If so, each capability name is validated against the known set of capabilities unless
// allowUnknownCapabilities is true. Each valid capability is then checked if it is disabled. If any invalid
// capabilities are found an error is returned listing all invalid capabilities. Otherwise, if any disabled
// capabilities are found an error is returned listing all disabled capabilities.
func checkResourceEnablement(annotations map[string]string, capabilities *configv1.ClusterVersionCapabilitiesStatus,
	allowUnknownCapabilities bool) error {

	caps := getManifestCapabilities(annotations)
	numCaps := len(caps)
	unknownCaps := make([]string, 0, numCaps)
	disabledCaps := make([]string, 0, numCaps)

	for _, c := range caps {

		if !allowUnknownCapabilities {
			var isKnownCap bool
			for _, knownCapability := range capabilities.KnownCapabilities {
				if c == knownCapability {
					isKnownCap = true
				}
			}
			if !isKnownCap {
				unknownCaps = append(unknownCaps, string(c))
				continue
			}
		}

		var isEnabledCap bool
		for _, enabledCapability := range capabilities.EnabledCapabilities {
			if c == enabledCapability {
				isEnabledCap = true
			}

		}
		if !isEnabledCap {
			disabledCaps = append(disabledCaps, string(c))
		}
	}
	if len(unknownCaps) > 0 {
		return fmt.Errorf("unrecognized capability names: %s", strings.Join(unknownCaps, ", "))
	}
	if len(disabledCaps) > 0 {
		return fmt.Errorf("disabled capabilities: %s", strings.Join(disabledCaps, ", "))
	}
	return nil
}

// GetManifestCapabilities returns the manifest's capabilities.
func (m *Manifest) GetManifestCapabilities() []configv1.ClusterVersionCapability {
	annotations := m.Obj.GetAnnotations()
	if annotations == nil {
		return nil
	}
	return getManifestCapabilities(annotations)
}

func getManifestCapabilities(annotations map[string]string) []configv1.ClusterVersionCapability {
	val, ok := annotations[CapabilityAnnotation]

	// check for empty string val to avoid returning length 1 slice of the empty string
	if !ok || val == "" {
		return nil
	}
	caps := strings.Split(val, "+")
	allCaps := make([]configv1.ClusterVersionCapability, len(caps))

	for i, c := range caps {
		allCaps[i] = configv1.ClusterVersionCapability(c)
	}
	return allCaps
}

// ManifestsFromFiles reads files and returns Manifests in the same order.
// 'files' should be list of absolute paths for the manifests on disk. An
// error is returned for each manifest that defines a duplicate resource
// as compared to other manifests defined within the 'files' list.
// Duplicate resources have the same group, kind, name, and namespace.
func ManifestsFromFiles(files []string) ([]Manifest, error) {
	var manifests []Manifest
	ids := make(map[resourceId]bool)
	var errs []error
	for _, file := range files {
		file, err := os.Open(file)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error opening %s", file.Name()))
			continue
		}
		defer file.Close()

		ms, err := ParseManifests(file)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error parsing %s", file.Name()))
			continue
		}
		filename := filepath.Base(file.Name())
		for i, m := range ms {
			ms[i].OriginalFilename = filename
			err = addIfNotDuplicateResource(m, ids)
			if err != nil {
				errs = append(errs, errors.Wrapf(err, "File %s contains", file.Name()))
			}
		}
		manifests = append(manifests, ms...)
	}

	agg := utilerrors.NewAggregate(errs)
	if agg != nil {
		return nil, fmt.Errorf("error loading manifests: %v", agg.Error())
	}

	return manifests, nil
}

// ParseManifests parses a YAML or JSON document that may contain one or more
// kubernetes resources. An error is returned if the input cannot be parsed
// or contains a duplicate resource.
func ParseManifests(r io.Reader) ([]Manifest, error) {
	theseIds := make(map[resourceId]bool)
	d := yaml.NewYAMLOrJSONDecoder(r, 1024)
	var manifests []Manifest
	for {
		m := Manifest{}
		if err := d.Decode(&m); err != nil {
			if err == io.EOF {
				return manifests, nil
			}
			return manifests, errors.Wrapf(err, "error parsing")
		}
		m.Raw = bytes.TrimSpace(m.Raw)
		if len(m.Raw) == 0 || bytes.Equal(m.Raw, []byte("null")) {
			continue
		}
		if err := addIfNotDuplicateResource(m, theseIds); err != nil {
			return manifests, err
		}
		manifests = append(manifests, m)
	}
}

// validateResourceId ensures the id contains the required fields per
// https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/#required-fields.
func validateResourceId(id resourceId) error {
	if id.Kind == "" || id.Name == "" {
		return fmt.Errorf("Resource with fields %s must contain kubernetes required fields kind and name", id)
	}
	return nil
}

func addIfNotDuplicateResource(manifest Manifest, resourceIds map[resourceId]bool) error {
	if _, ok := resourceIds[manifest.id]; !ok {
		resourceIds[manifest.id] = true
		return nil
	}
	return fmt.Errorf("duplicate resource: (%s)", manifest.id)
}
