package nutanix

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kdomanski/iso9660"
	"github.com/pkg/errors"
	nutanixClientV3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	diskLabel        = "config-2"
	isoFile          = "bootstrap-ign.iso"
	metadataFilePath = "openstack/latest/meta_data.json"
	userDataFilePath = "openstack/latest/user_data"
	sleepTime        = 10 * time.Second
	timeout          = 5 * time.Minute
)

type metadataCloudInit struct {
	UUID string `json:"uuid"`
}

func bootISOImageName(infraID string) string {
	return fmt.Sprintf("%s-%s", infraID, isoFile)
}

func createBootstrapISO(infraID, userData string) (string, error) {
	id := uuid.New()
	metaObj := &metadataCloudInit{
		UUID: id.String(),
	}
	fullISOFile := bootISOImageName(infraID)
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

func waitForTasks(clientV3 nutanixClientV3.Service, taskUUIDs []string) error {
	for _, t := range taskUUIDs {
		err := waitForTask(clientV3, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func waitForTask(clientV3 nutanixClientV3.Service, taskUUID string) error {
	finished := false
	var err error
	for start := time.Now(); time.Since(start) < timeout; {
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
		return errors.Errorf("timeout while waiting for task UUID: %s", taskUUID)
	}

	return nil
}

func isTaskFinished(clientV3 nutanixClientV3.Service, taskUUID string) (bool, error) {
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
	return false, errors.Errorf("Retrieved unexpected task status: %s", status)

}

func getTaskStatus(clientV3 nutanixClientV3.Service, taskUUID string) (string, error) {
	v, err := clientV3.GetTask(taskUUID)

	if err != nil {
		return "", err
	}

	if *v.Status == "INVALID_UUID" || *v.Status == "FAILED" {
		return *v.Status, errors.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
	}
	return *v.Status, nil
}

// RHCOSImageName is the unique image name for a given cluster
func RHCOSImageName(infraID string) string {
	return fmt.Sprintf("%s-rhcos", infraID)
}
