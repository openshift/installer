package smoke

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	// These consts are extension name's we want to ensure are present in
	// BigQuery results.
	accountIDExtension              = "accountID"
	certificatesStrategyExtension   = "certificatesStrategy"
	installerPlatformExtension      = "installerPlatform"
	tectonicUpdaterEnabledExtension = "tectonicUpdaterEnabled"
	// extensionsNameKey is the key for extension names in BigQuery results.
	extensionsNameKey = "extensions_name"
	// extensionsValueKey is the key for extension values in BigQuery results.
	extensionsValueKey = "extensions_value"
	// bigQuerySpecEnv is the environment variable containing the spec of the BigQuery table to test for cluster metrics.
	bigQuerySpecEnv = "SMOKE_BIGQUERY_SPEC"
)

func testQA(t *testing.T) {
	t.Run("StatsEmitterLogs", testGetStatsEmitterLogs)
	t.Run("BigQueryData", testGetBigQueryData)
}

func getStatsEmitterLogs(t *testing.T) error {
	c, _ := newClient(t)
	expected := "report successfully sent"
	namespace := "tectonic-system"
	podPrefix := "tectonic-stats-emitter"
	logs, err := validatePodLogging(c, namespace, podPrefix)
	if err != nil {
		return fmt.Errorf("failed to gather logs for %s/%s, %v", namespace, podPrefix, err)
	}
	if !bytes.Contains(logs, []byte(expected)) {
		return fmt.Errorf("expected logs to contain %q", expected)
	}
	return nil
}

func testGetStatsEmitterLogs(t *testing.T) {
	max := 1 * time.Minute
	err := retry(getStatsEmitterLogs, t, 5*time.Second, max)
	if err != nil {
		t.Fatalf("Failed to verify stats-emitter success in logs in %v.", max)
	}
	t.Log("Successfully verified stats-emitter success in logs.")
}

// getBigQueryData finds the Tectonic cluster ID from the Tectonic configmap
// in Kubernetes and uses it, along with the provided BigQuery spec, to query
// BigQuery for the metrics for the given cluster.
func getBigQueryData(t *testing.T) error {
	// Parse BigQuery spec.
	bigQuerySpec := os.Getenv(bigQuerySpecEnv)
	project, dataset, table, err := parseBigQuerySpec(bigQuerySpec)
	if err != nil {
		return fmt.Errorf("failed to parse BigQuery spec: %v", err)
	}
	// Get Tectonic cluster configuration.
	cm, err := getTectonicClusterConfig(t)
	if err != nil {
		return fmt.Errorf("failed to get Tectonic cluster configuration: %v", err)
	}
	cid, ok := cm.Data["clusterID"]
	if !ok {
		return errors.New("failed to find cluster ID in ConfigMap")
	}
	// Initialize BigQuery client.
	ctx := context.Background()
	// This assumes that:
	//  a) a GCE ServiceAccount has been created for this app
	//  b) the ServiceAccount is an owner of the dataset for this app
	//  c) the credentials for the ServiceAccount are in a file
	//  d) env GOOGLE_APPLICATION_CREDENTIALS=<path to credentials file>
	bq, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to create BigQuery client: %v", err)
	}
	// Get cluster stats extensions from BigQuery.
	q := bq.Query(`SELECT
 extensions.name,
 extensions.value,
FROM
  FLATTEN([` + fmt.Sprintf("%s:%s.%s", project, dataset, table) + `], extensions)
WHERE
 clusterID = '` + cid + `'
GROUP BY
 extensions.name,
 extensions.value`)
	expected := make(map[string]string)
	found := make(map[string]string)
	// extensions is an array of the tested stats extensions.
	var extensions = []string{accountIDExtension, certificatesStrategyExtension, installerPlatformExtension, tectonicUpdaterEnabledExtension}
	for _, name := range extensions {
		// Some extensions are not in the ConfigMap and so do not have
		// expected values. Instead, we just expect them to be present
		// in BigQuery and do not care about their values.
		if value, ok := cm.Data[name]; ok {
			expected[name] = value
		}
	}
	it, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read query results: %v", err)
	}
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to get next row: %v", err)
		}
		n, _ := row[extensionsNameKey]
		name, ok := n.(string)
		if !ok {
			return fmt.Errorf("expected extension name to be a string")
		}
		v, _ := row[extensionsValueKey]
		value, ok := v.(string)
		if !ok {
			return fmt.Errorf("expected extension value to be a string")
		}
		found[name] = value
	}
	// Ensure stats extensions are in BigQuery.
	var wrong []string
	for _, name := range extensions {
		expectedValue, ok := expected[name]
		// If the extension does not have an expected value,
		// then just check if it is present in BigQuery at all.
		if !ok {
			if _, ok := found[name]; !ok {
				wrong = append(wrong, fmt.Sprintf("did not find extension %q", name))
			}
			continue
		}
		if foundValue, _ := found[name]; foundValue != expectedValue {
			wrong = append(wrong, fmt.Sprintf("expected extension %q to be %q, got %q", name, expectedValue, foundValue))
		}
	}
	if len(wrong) != 0 {
		return fmt.Errorf("failed to find extensions in BigQuery results: %s", strings.Join(wrong, "; "))
	}
	return nil
}

func testGetBigQueryData(t *testing.T) {
	max := 1 * time.Minute
	err := retry(getBigQueryData, t, 10*time.Second, max)
	if err != nil {
		t.Fatalf("Failed to verify stats-emitter data in BigQuery in %v.", max)
	}
	t.Log("Successfully verified stats-emitter data in BigQuery.")
}

// bqre is a regular expression for parse BigQuery specs.
var bqre = regexp.MustCompile(`^bigquery://([^.]+)\.([^.]+)\.([^.]+)$`)

// parseBigQuerySpec parses a spec formatted as `bigquery://project.dataset.table`.
// The 3 string returns are project, dataset, and table respectively.
// This will return an error if it does not believe the argument is a BigQuery spec,
// or if it believes the argument is a biquery spec but it can't parse it properly.
func parseBigQuerySpec(spec string) (string, string, string, error) {
	if !strings.HasPrefix(spec, "bigquery://") {
		return "", "", "", errors.New("BigQuery spec must begin with \"bigquery://\"")
	}
	subs := bqre.FindStringSubmatch(spec)
	if len(subs) != 4 {
		return "", "", "", fmt.Errorf("invalid BigQuery spec: %q", spec)
	}
	return subs[1], subs[2], subs[3], nil
}

// getTectonicClusterConfig gets the cluster's configuration from the tectonic-config ConfigMap.
func getTectonicClusterConfig(t *testing.T) (*v1.ConfigMap, error) {
	configmapName := "tectonic-config"
	c, _ := newClient(t)
	cm, err := c.Core().ConfigMaps(tectonicSystemNamespace).Get(configmapName, meta_v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to find ConfigMap %q: %v", configmapName, err)
	}
	return cm, nil
}
