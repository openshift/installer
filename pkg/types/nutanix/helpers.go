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

	machinev1 "github.com/openshift/api/machine/v1"
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
func RHCOSImageName(p *Platform, infraID string) string {
	imgName := p.PreloadedOSImageName
	if imgName == "" {
		imgName = fmt.Sprintf("%s-rhcos", infraID)
	}

	return imgName
}

// CategoryKey returns the cluster specific category key name.
func CategoryKey(infraID string) string {
	categoryKey := fmt.Sprintf("%s%s", categoryKeyPrefix, infraID)
	return categoryKey
}

// GetGPUList returns a list of VMGpus for the given list of GPU identifiers in the Prism Element (uuid).
func GetGPUList(ctx context.Context, client *nutanixclientv3.Client, gpus []machinev1.NutanixGPU, peUUID string) ([]*nutanixclientv3.VMGpu, error) {
	vmGPUs := make([]*nutanixclientv3.VMGpu, 0)

	if len(gpus) == 0 {
		return vmGPUs, nil
	}

	peGPUs, err := GetGPUsForPE(ctx, client, peUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve GPUs of the Prism Element cluster (uuid: %v): %w", peUUID, err)
	}
	if len(peGPUs) == 0 {
		return nil, fmt.Errorf("no available GPUs found in Prism Element cluster (uuid: %s)", peUUID)
	}

	for _, gpu := range gpus {
		foundGPU, err := GetGPUFromList(ctx, client, gpu, peGPUs)
		if err != nil {
			return nil, err
		}
		vmGPUs = append(vmGPUs, foundGPU)
	}

	return vmGPUs, nil
}

// GetGPUFromList returns the VMGpu matching the input reqirements from the provided list of GPU devices.
func GetGPUFromList(ctx context.Context, client *nutanixclientv3.Client, gpu machinev1.NutanixGPU, gpuDevices []*nutanixclientv3.GPU) (*nutanixclientv3.VMGpu, error) {
	for _, gd := range gpuDevices {
		if gd.Status != "UNUSED" {
			continue
		}

		if (gpu.Type == machinev1.NutanixGPUIdentifierDeviceID && gd.DeviceID != nil && *gpu.DeviceID == int32(*gd.DeviceID)) ||
			(gpu.Type == machinev1.NutanixGPUIdentifierName && *gpu.Name == gd.Name) {
			return &nutanixclientv3.VMGpu{
				DeviceID: gd.DeviceID,
				Mode:     &gd.Mode,
				Vendor:   &gd.Vendor,
			}, nil
		}
	}

	return nil, fmt.Errorf("no available GPU found that matches required GPU inputs")
}

// GetGPUsForPE returns all the GPU devices for the given Prism Element (uuid).
func GetGPUsForPE(ctx context.Context, client *nutanixclientv3.Client, peUUID string) ([]*nutanixclientv3.GPU, error) {
	gpus := make([]*nutanixclientv3.GPU, 0)
	hosts, err := client.V3.ListAllHost(ctx)
	if err != nil {
		return gpus, fmt.Errorf("failed to get hosts from Prism Central: %w", err)
	}

	for _, host := range hosts.Entities {
		if host == nil ||
			host.Status == nil ||
			host.Status.ClusterReference == nil ||
			host.Status.ClusterReference.UUID != peUUID ||
			host.Status.Resources == nil ||
			len(host.Status.Resources.GPUList) == 0 {
			continue
		}

		for _, peGpu := range host.Status.Resources.GPUList {
			if peGpu != nil {
				gpus = append(gpus, peGpu)
			}
		}
	}

	return gpus, nil
}

// FindImageUUIDByName retrieves the image resource uuid by the given image name from PC.
func FindImageUUIDByName(ctx context.Context, ntnxclient *nutanixclientv3.Client, imageName string) (*string, error) {
	res, err := ntnxclient.V3.ListImage(ctx, &nutanixclientv3.DSMetadata{
		Filter: utils.StringPtr(fmt.Sprintf("name==%s", imageName)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find image by name %q in PC/PE. err: %w", imageName, err)
	}

	if len(res.Entities) == 0 {
		return nil, fmt.Errorf("no image with name %q is found in PC/PE", imageName)
	}

	if len(res.Entities) > 1 {
		return nil, fmt.Errorf("found more than one (%v) images with name %q in PC/PE", len(res.Entities), imageName)
	}

	return res.Entities[0].Metadata.UUID, nil
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
