package crd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/yaml"
)

var (
	testEnv   *envtest.Environment
	dynClient dynamic.Interface
	envtestErr error

	installConfigGVR = schema.GroupVersionResource{
		Group:    "install.openshift.io",
		Version:  "v1",
		Resource: "installconfigs",
	}
)

func TestMain(m *testing.M) {
	crd := loadCRDRaw()

	// The CRD has known schema issues from other platforms (Azure identity
	// defaults, Nutanix map-type, PowerVC uniqueItems, AWS/ibmcloud CEL cost)
	// that prevent the API server from accepting it. Keep only the platforms
	// under test so envtest can install the CRD for behavioral testing.
	patchCRDForEnvtest(crd)

	testEnv = &envtest.Environment{
		CRDs:                  []*apiextensionsv1.CustomResourceDefinition{crd},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	if err != nil {
		envtestErr = fmt.Errorf("failed to start envtest: %w", err)
		fmt.Fprintf(os.Stderr, "warning: %v (behavioral tests will be skipped)\n", envtestErr)
		code := m.Run()
		os.Exit(code)
	}

	dynClient, err = dynamic.NewForConfig(cfg)
	if err != nil {
		testEnv.Stop()
		envtestErr = fmt.Errorf("failed to create dynamic client: %w", err)
		fmt.Fprintf(os.Stderr, "warning: %v (behavioral tests will be skipped)\n", envtestErr)
		code := m.Run()
		os.Exit(code)
	}

	code := m.Run()
	testEnv.Stop()
	os.Exit(code)
}

func loadCRDRaw() *apiextensionsv1.CustomResourceDefinition {
	data, err := os.ReadFile(filepath.Join("..", "..", "data", "data", "install.openshift.io_installconfigs.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load CRD: %v\n", err)
		os.Exit(1)
	}
	var crd apiextensionsv1.CustomResourceDefinition
	if err := yaml.Unmarshal(data, &crd); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal CRD: %v\n", err)
		os.Exit(1)
	}
	return &crd
}

// patchCRDForEnvtest strips the CRD down to only the top-level required
// fields and the platform sections under test. The full CRD exceeds the
// API server's CEL cost budget due to its size and unconstrained strings
// across many platforms. This only affects behavioral tests — schema tests
// load the unmodified CRD YAML directly.
func patchCRDForEnvtest(crd *apiextensionsv1.CustomResourceDefinition) {
	keepPlatforms := map[string]bool{"vsphere": true}

	for i := range crd.Spec.Versions {
		schema := crd.Spec.Versions[i].Schema
		if schema == nil || schema.OpenAPIV3Schema == nil {
			continue
		}
		root := schema.OpenAPIV3Schema

		keepTopLevel := map[string]bool{
			"baseDomain": true,
			"metadata":   true,
			"platform":   true,
			"pullSecret": true,
		}
		for name := range root.Properties {
			if !keepTopLevel[name] {
				delete(root.Properties, name)
			}
		}

		stripPlatforms(root, "platform", keepPlatforms)

		root.Required = filterStrings(root.Required, keepTopLevel)

		delete(root.Properties, "controlPlane")
		delete(root.Properties, "compute")
		delete(root.Properties, "arbiter")

		boundStringLengths(root)
		boundArraySizes(root)
	}
}

func stripPlatforms(root *apiextensionsv1.JSONSchemaProps, key string, keep map[string]bool) {
	plat, ok := root.Properties[key]
	if !ok {
		return
	}
	for name := range plat.Properties {
		if !keep[name] {
			delete(plat.Properties, name)
		}
	}
	root.Properties[key] = plat
}

func filterStrings(list []string, keep map[string]bool) []string {
	var result []string
	for _, s := range list {
		if keep[s] {
			result = append(result, s)
		}
	}
	return result
}

// boundArraySizes walks the schema tree and adds maxItems=64
// to all arrays that lack it. Without maxItems, the CEL cost
// estimator multiplies per-item rule costs by an unbounded factor.
func boundArraySizes(schema *apiextensionsv1.JSONSchemaProps) {
	mi := int64(64)
	if schema.Type == "array" && schema.MaxItems == nil {
		schema.MaxItems = &mi
	}
	for name, prop := range schema.Properties {
		boundArraySizes(&prop)
		schema.Properties[name] = prop
	}
	if schema.Items != nil && schema.Items.Schema != nil {
		boundArraySizes(schema.Items.Schema)
	}
}

// boundStringLengths walks the schema tree and adds maxLength=256
// to all string-typed fields and array items that lack it. Without
// maxLength, the CEL cost estimator assumes unbounded string length
// and rejects the CRD. This is test infrastructure only.
func boundStringLengths(schema *apiextensionsv1.JSONSchemaProps) {
	ml := int64(256)
	if schema.Type == "string" && schema.MaxLength == nil {
		schema.MaxLength = &ml
	}
	for name, prop := range schema.Properties {
		boundStringLengths(&prop)
		schema.Properties[name] = prop
	}
	if schema.Items != nil && schema.Items.Schema != nil {
		boundStringLengths(schema.Items.Schema)
	}
	for i := range schema.AllOf {
		boundStringLengths(&schema.AllOf[i])
	}
	for i := range schema.AnyOf {
		boundStringLengths(&schema.AnyOf[i])
	}
	for i := range schema.OneOf {
		boundStringLengths(&schema.OneOf[i])
	}
}

// skipIfNoEnvtest skips the test if envtest is not available.
func skipIfNoEnvtest(t *testing.T) {
	t.Helper()
	if envtestErr != nil {
		t.Skipf("envtest not available: %v", envtestErr)
	}
	if dynClient == nil {
		t.Skip("envtest not initialized")
	}
}

// loadCRD reads and unmarshals the InstallConfig CRD YAML.
func loadCRD(t *testing.T) *apiextensionsv1.CustomResourceDefinition {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "data", "data", "install.openshift.io_installconfigs.yaml"))
	require.NoError(t, err, "failed to load CRD")

	var crd apiextensionsv1.CustomResourceDefinition
	require.NoError(t, yaml.Unmarshal(data, &crd), "failed to unmarshal CRD")
	return &crd
}

// rootSchema returns the top-level OpenAPI v3 schema from the CRD.
func rootSchema(t *testing.T, crd *apiextensionsv1.CustomResourceDefinition) *apiextensionsv1.JSONSchemaProps {
	t.Helper()
	require.NotEmpty(t, crd.Spec.Versions, "CRD has no versions")
	require.NotNil(t, crd.Spec.Versions[0].Schema, "CRD version has no schema")
	require.NotNil(t, crd.Spec.Versions[0].Schema.OpenAPIV3Schema, "CRD version has no OpenAPI v3 schema")
	return crd.Spec.Versions[0].Schema.OpenAPIV3Schema
}

// schemaAtPath walks the CRD schema tree by property name.
// Use "[]" to descend into an array's item schema.
func schemaAtPath(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, path ...string) *apiextensionsv1.JSONSchemaProps {
	t.Helper()
	current := schema
	for _, segment := range path {
		if segment == "[]" {
			require.NotNil(t, current.Items, "expected array items at path segment '[]'")
			require.NotNil(t, current.Items.Schema, "expected single schema for array items")
			current = current.Items.Schema
			continue
		}
		prop, ok := current.Properties[segment]
		require.True(t, ok, "missing property %q in schema (available: %v)", segment, propertyNames(current))
		current = &prop
	}
	return current
}

func propertyNames(schema *apiextensionsv1.JSONSchemaProps) []string {
	names := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		names = append(names, k)
	}
	return names
}

// requireCELRule asserts the schema has an XValidation rule matching the given string.
func requireCELRule(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, rule string) {
	t.Helper()
	for _, v := range schema.XValidations {
		if v.Rule == rule {
			return
		}
	}
	rules := make([]string, 0, len(schema.XValidations))
	for _, v := range schema.XValidations {
		rules = append(rules, v.Rule)
	}
	t.Errorf("CEL rule %q not found; existing rules: %v", rule, rules)
}

// requireEnum asserts the schema's Enum constraint contains the given values.
// It checks both the direct Enum field and any AllOf entries that define enums.
func requireEnum(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, values ...string) {
	t.Helper()
	allEnums := collectEnums(schema)
	require.NotEmpty(t, allEnums, "expected enum constraint")
	for _, v := range values {
		assert.Contains(t, allEnums, v, "enum missing value %q", v)
	}
}

func collectEnums(schema *apiextensionsv1.JSONSchemaProps) []string {
	var result []string
	for _, e := range schema.Enum {
		var s string
		if err := yaml.Unmarshal(e.Raw, &s); err == nil {
			result = append(result, s)
		}
	}
	for _, a := range schema.AllOf {
		result = append(result, collectEnums(&a)...)
	}
	return result
}

// requireFormat asserts the schema's Format field. Pass "" to assert no format is set.
func requireFormat(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, format string) {
	t.Helper()
	assert.Equal(t, format, schema.Format)
}

// requirePattern asserts the schema's Pattern constraint.
func requirePattern(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, pattern string) {
	t.Helper()
	assert.Equal(t, pattern, schema.Pattern, "unexpected pattern")
}

// requireStringLength asserts MinLength/MaxLength. Pass nil to skip either bound.
func requireStringLength(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, min, max *int64) {
	t.Helper()
	if min != nil {
		require.NotNil(t, schema.MinLength, "expected minLength to be set")
		assert.Equal(t, *min, *schema.MinLength, "unexpected minLength")
	}
	if max != nil {
		require.NotNil(t, schema.MaxLength, "expected maxLength to be set")
		assert.Equal(t, *max, *schema.MaxLength, "unexpected maxLength")
	}
}

// requireArrayItems asserts MinItems/MaxItems. Pass nil to skip either bound.
func requireArrayItems(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, min, max *int64) {
	t.Helper()
	if min != nil {
		require.NotNil(t, schema.MinItems, "expected minItems to be set")
		assert.Equal(t, *min, *schema.MinItems, "unexpected minItems")
	}
	if max != nil {
		require.NotNil(t, schema.MaxItems, "expected maxItems to be set")
		assert.Equal(t, *max, *schema.MaxItems, "unexpected maxItems")
	}
}

// requireNumericRange asserts Minimum/Maximum. Pass nil to skip either bound.
func requireNumericRange(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, min, max *float64) {
	t.Helper()
	if min != nil {
		require.NotNil(t, schema.Minimum, "expected minimum to be set")
		assert.Equal(t, *min, *schema.Minimum, "unexpected minimum")
	}
	if max != nil {
		require.NotNil(t, schema.Maximum, "expected maximum to be set")
		assert.Equal(t, *max, *schema.Maximum, "unexpected maximum")
	}
}

// requireRequired asserts the schema's Required list includes the given field names.
func requireRequired(t *testing.T, schema *apiextensionsv1.JSONSchemaProps, fields ...string) {
	t.Helper()
	for _, f := range fields {
		assert.Contains(t, schema.Required, f, "field %q not in required list", f)
	}
}

// int64Ptr returns a pointer to the given int64 value.
func int64Ptr(v int64) *int64 { return &v }

// float64Ptr returns a pointer to the given float64 value.
func float64Ptr(v float64) *float64 { return &v }

// createCR creates an InstallConfig CR and returns the error from the API server.
// On success, registers a cleanup to delete the resource.
func createCR(t *testing.T, obj map[string]any) error {
	t.Helper()
	u := &unstructured.Unstructured{Object: obj}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "install.openshift.io",
		Version: "v1",
		Kind:    "InstallConfig",
	})

	created, err := dynClient.Resource(installConfigGVR).Namespace("default").Create(
		context.Background(), u, metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	t.Cleanup(func() {
		dynClient.Resource(installConfigGVR).Namespace("default").Delete(
			context.Background(), created.GetName(), metav1.DeleteOptions{},
		)
	})
	return nil
}
