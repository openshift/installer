package baremetal

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
)

type defIgnition struct {
	Name     string
	PoolName string
	Content  string
}

func (ign *defIgnition) createFile() (string, error) {
	tempFile, err := os.CreateTemp("", ign.Name)
	if err != nil {
		return "", fmt.Errorf("Error creating tmp file: %v", err)
	}
	defer tempFile.Close()

	if _, err := tempFile.WriteString(ign.Content); err != nil {
		return "", fmt.Errorf("Cannot write Ignition object to temporary " +
			"ignition file")
	}

	return tempFile.Name(), nil
}

func (ign *defIgnition) CreateAndUpload(client *libvirt.Libvirt) (string, error) {
	pool, err := client.StoragePoolLookupByName(ign.PoolName)
	if err != nil {
		return "", fmt.Errorf("can't find storage pool '%s'", ign.PoolName)
	}

	err = client.StoragePoolRefresh(pool, 0)
	if err != nil {
		return "", fmt.Errorf("failed to refresh pool %v", err)
	}

	volumeDef := newVolume(ign.Name)

	ignFile, err := ign.createFile()
	if err != nil {
		return "", err
	}
	defer func() {
		if err = os.Remove(ignFile); err != nil {
			logrus.Errorf("Error while removing tmp Ignition file: %v", err)
		}
	}()

	img, err := newImage(ignFile)
	if err != nil {
		return "", err
	}

	size, err := img.Size()
	if err != nil {
		return "", err
	}

	volumeDef.Capacity.Unit = "B"
	volumeDef.Capacity.Value = size
	volumeDef.Target.Format.Type = "raw"

	volumeDefXML, err := xml.Marshal(volumeDef)
	if err != nil {
		return "", fmt.Errorf("Error serializing libvirt volume: %s", err)
	}

	volume, err := client.StorageVolCreateXML(pool, string(volumeDefXML), 0)
	if err != nil {
		return "", fmt.Errorf("Error creating libvirt volume for Ignition %s: %s", ign.Name, err)
	}

	err = img.Import(newCopier(client, volume, volumeDef.Capacity.Value), volumeDef)
	if err != nil {
		return "", fmt.Errorf("Error while uploading ignition file %s: %s", img.String(), err)
	}

	return "", nil
}
