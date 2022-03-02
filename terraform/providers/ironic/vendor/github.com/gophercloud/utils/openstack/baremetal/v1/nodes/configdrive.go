package nodes

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// A ConfigDrive struct will be used to create a base64-encoded, gzipped ISO9660 image for use with Ironic.
type ConfigDrive struct {
	UserData       UserDataBuilder        `json:"user_data"`
	MetaData       map[string]interface{} `json:"meta_data"`
	NetworkData    map[string]interface{} `json:"network_data"`
	Version        string                 `json:"-"`
	BuildDirectory string                 `json:"-"`
}

// Interface to let us specify a raw string, or a map for the user data
type UserDataBuilder interface {
	ToUserData() ([]byte, error)
}

type UserDataMap map[string]interface{}
type UserDataString string

// Converts a UserDataMap to JSON-string
func (data UserDataMap) ToUserData() ([]byte, error) {
	return json.MarshalIndent(data, "", "\t")
}

func (data UserDataString) ToUserData() ([]byte, error) {
	return []byte(data), nil
}

type ConfigDriveBuilder interface {
	ToConfigDrive() (string, error)
}

// Writes out a ConfigDrive to a temporary directory, and returns the path
func (configDrive ConfigDrive) ToDirectory() (string, error) {
	// Create a temporary directory for our config drive
	directory, err := ioutil.TempDir(configDrive.BuildDirectory, "gophercloud")
	if err != nil {
		return "", err
	}

	// Build up the paths for OpenStack
	var version string
	if configDrive.Version == "" {
		version = "latest"
	} else {
		version = configDrive.Version
	}

	path := filepath.FromSlash(directory + "/openstack/" + version)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	// Dump out user data
	if configDrive.UserData != nil {
		userDataPath := filepath.FromSlash(path + "/user_data")
		data, err := configDrive.UserData.ToUserData()
		if err != nil {
			return "", err
		}

		if err := ioutil.WriteFile(userDataPath, data, 0644); err != nil {
			return "", err
		}
	}

	// Dump out meta data
	if configDrive.MetaData != nil {
		metaDataPath := filepath.FromSlash(path + "/meta_data.json")
		data, err := json.Marshal(configDrive.MetaData)
		if err != nil {
			return "", err
		}

		if err := ioutil.WriteFile(metaDataPath, data, 0644); err != nil {
			return "", err
		}
	}

	// Dump out network data
	if configDrive.NetworkData != nil {
		networkDataPath := filepath.FromSlash(path + "/network_data.json")
		data, err := json.Marshal(configDrive.NetworkData)
		if err != nil {
			return "", err
		}

		if err := ioutil.WriteFile(networkDataPath, data, 0644); err != nil {
			return "", err
		}
	}

	return directory, nil
}

// Writes out the ConfigDrive struct to a directory structure, and then
// packs it as a base64-encoded gzipped ISO9660 image.
func (configDrive ConfigDrive) ToConfigDrive() (string, error) {
	directory, err := configDrive.ToDirectory()
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(directory)

	// Pack result as gzipped ISO9660 file
	result, err := PackDirectoryAsISO(directory)
	if err != nil {
		return "", err
	}

	// Return as base64-encoded data
	return base64.StdEncoding.EncodeToString(result), nil
}
