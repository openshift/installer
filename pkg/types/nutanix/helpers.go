package nutanix

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/google/uuid"
	"github.com/kdomanski/iso9660"
	"github.com/nutanix-cloud-native/prism-go-client/utils"
	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"
	"k8s.io/utils/ptr"
)

const (
	diskLabel        = "config-2"
	isoFile          = "bootstrap-ign.iso"
	metadataFilePath = "openstack/latest/meta_data.json"
	userDataFilePath = "openstack/latest/user_data"
	sleepTime        = 10 * time.Second
	// DefaultPrismAPICallTimeout is 10 minutes.
	DefaultPrismAPICallTimeout = int(10)

	// Category Key format: "kubernetes-io-cluster-<cluster-id>".
	categoryKeyPrefix = "kubernetes-io-cluster-"
	// CategoryValueOwned is the category value representing owned by the cluster.
	CategoryValueOwned = "owned"
	// CategoryValueShared is the category value representing shared by the cluster.
	CategoryValueShared = "shared"
)

var (
	// prismAPICallTimeoutDuration is the timeout for the prism-api calls.
	// It can be changed by calling SetPrismAPICallTimeoutDuration().
	prismAPICallTimeoutDuration = time.Duration(DefaultPrismAPICallTimeout) * time.Minute
)

type metadataCloudInit struct {
	UUID string `json:"uuid"`
}

// SetPrismAPICallTimeoutDuration sets and returns the timeout duration value for prism-api calls.
func SetPrismAPICallTimeoutDuration(p *Platform) time.Duration {
	if p.PrismAPICallTimeout != nil {
		prismAPICallTimeoutDuration = time.Duration(*p.PrismAPICallTimeout) * time.Minute
	}
	return prismAPICallTimeoutDuration
}

// BootISOImageName is the image name for Bootstrap node for a given infraID.
func BootISOImageName(infraID string) string {
	return fmt.Sprintf("%s-%s", infraID, isoFile)
}

// BootISOImagePath is the image path for Bootstrap node for a given infraID and path.
func BootISOImagePath(path, infraID string) string {
	imgName := BootISOImageName(infraID)
	application := "openshift-installer"
	subdir := "image_cache"
	fullISOFile := filepath.Join(path, application, subdir, imgName)
	return fullISOFile
}

// CreateBootstrapISO creates a ISO for the bootstrap node.
func CreateBootstrapISO(infraID, userData string) (string, error) {
	id := uuid.New()
	metaObj := &metadataCloudInit{
		UUID: id.String(),
	}
	metadata, err := json.Marshal(metaObj)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to marshal metadata struct to json"))
	}

	writer, err := iso9660.NewWriter()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to create writer: %s", err))
	}

	defer writer.Cleanup()

	userDataReader := strings.NewReader(userData)
	err = writer.AddFile(userDataReader, userDataFilePath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to add file: %s", err))
	}

	metadataReader := strings.NewReader(string(metadata))
	err = writer.AddFile(metadataReader, metadataFilePath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to add file: %s", err))
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to fetch user cache dir")
	}

	fullISOFile := BootISOImagePath(cacheDir, infraID)
	fullISOFileDir, err := filepath.Abs(filepath.Dir(fullISOFile))
	if err != nil {
		return "", errors.Wrap(err, "unable to extract parent directory from bootstrap iso filepath")
	}

	_, err = os.Stat(fullISOFileDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(fullISOFileDir, 0755)
			if err != nil {
				return "", errors.Wrap(err, fmt.Sprintf("failed to create %s", fullISOFileDir))
			}
		} else {
			return "", errors.Wrap(err, fmt.Sprintf("cannot access %s", fullISOFileDir))
		}
	}

	outputFile, err := os.OpenFile(fullISOFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to create file: %s", err))
	}

	err = writer.WriteTo(outputFile, diskLabel)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to write ISO image: %s", err))
	}

	err = outputFile.Close()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to close output file: %s", err))
	}

	return fullISOFile, nil
}

// WaitForTasks is a wrapper for WaitForTask.
func WaitForTasks(clientV3 nutanixclientv3.Service, taskUUIDs []string) error {
	for _, t := range taskUUIDs {
		err := WaitForTask(clientV3, t)
		if err != nil {
			return err
		}
	}
	return nil
}

// WaitForTask waits until a queued task has been finished or timeout has been reached.
func WaitForTask(clientV3 nutanixclientv3.Service, taskUUID string) error {
	finished := false
	var err error
	start := time.Now()

	for time.Since(start) < prismAPICallTimeoutDuration {
		finished, err = isTaskFinished(clientV3, taskUUID)
		if err != nil {
			return err
		}
		if finished {
			break
		}
		time.Sleep(sleepTime)
	}
	if !finished {
		return errors.Errorf("timeout while waiting for task UUID: %s, used_time: %s", taskUUID, time.Since(start))
	}

	return nil
}

func isTaskFinished(clientV3 nutanixclientv3.Service, taskUUID string) (bool, error) {
	isFinished := map[string]bool{
		"QUEUED":    false,
		"RUNNING":   false,
		"SUCCEEDED": true,
	}
	status, err := getTaskStatus(clientV3, taskUUID)
	if err != nil {
		return false, err
	}
	if val, ok := isFinished[status]; ok {
		return val, nil
	}
	return false, errors.Errorf("retrieved unexpected task status: %s", status)
}

func getTaskStatus(clientV3 nutanixclientv3.Service, taskUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	v, err := clientV3.GetTask(ctx, taskUUID)

	if err != nil {
		return "", err
	}

	if *v.Status == "INVALID_UUID" || *v.Status == "FAILED" {
		return *v.Status, errors.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
	}
	return *v.Status, nil
}

// RHCOSImageName is the unique image name for a given cluster.
func RHCOSImageName(infraID string) string {
	return fmt.Sprintf("%s-rhcos", infraID)
}

// CategoryKey returns the cluster specific category key name.
func CategoryKey(infraID string) string {
	categoryKey := fmt.Sprintf("%s%s", categoryKeyPrefix, infraID)
	return categoryKey
}

// InsertHostnameIgnition inserts the file "/etc/hostname" with the given hostname to the provided Ignition config data.
func InsertHostnameIgnition(ignData []byte, hostname string) ([]byte, error) {
	ignConfig := &igntypes.Config{}
	if err := json.Unmarshal(ignData, &ignConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Ignition config: %w", err)
	}

	hostnameFile := igntypes.File{
		Node: igntypes.Node{
			Path:      "/etc/hostname",
			Overwrite: ptr.To(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: ptr.To(420),
			Contents: igntypes.Resource{
				Source: ptr.To(dataurl.EncodeBytes([]byte(hostname))),
			},
		},
	}

	if ignConfig.Storage.Files == nil {
		ignConfig.Storage.Files = make([]igntypes.File, 0)
	}
	ignConfig.Storage.Files = append(ignConfig.Storage.Files, hostnameFile)

	ign, err := json.Marshal(ignConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ignition data: %w", err)
	}

	return ign, nil
}
