package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type installStateConf struct {
	InstallConfig struct {
		Config struct {
			Metadata struct {
				Name              string      `json:"name"`
				CreationTimestamp interface{} `json:"creationTimestamp"`
			} `json:"metadata"`
			ClusterID  string `json:"clusterID"`
			SSHKey     string `json:"sshKey"`
			BaseDomain string `json:"baseDomain"`
			Networking struct {
				Type            string `json:"type"`
				ServiceCIDR     string `json:"serviceCIDR"`
				ClusterNetworks []struct {
					Cidr             string `json:"cidr"`
					HostSubnetLength int    `json:"hostSubnetLength"`
				} `json:"clusterNetworks"`
			} `json:"networking"`
			Machines []struct {
				Name     string `json:"name"`
				Replicas int    `json:"replicas"`
				Platform struct {
				} `json:"platform"`
			} `json:"machines"`
			Platform struct {
				Libvirt struct {
					URI                    string `json:"URI"`
					DefaultMachinePlatform struct {
						Image string `json:"image"`
					} `json:"defaultMachinePlatform"`
					Network struct {
						If      string `json:"if"`
						IPRange string `json:"ipRange"`
					} `json:"network"`
					MasterIPs interface{} `json:"masterIPs"`
				} `json:"libvirt"`
			} `json:"platform"`
			PullSecret string `json:"pullSecret"`
		} `json:"config"`
		File struct {
			Filename string `json:"Filename"`
			Data     string `json:"Data"`
		} `json:"file"`
	} `json:"*installconfig.InstallConfig"`
	ICBaseDomain struct {
		BaseDomain string `json:"BaseDomain"`
	} `json:"*installconfig.baseDomain"`
	ICClusterID struct {
		ClusterID string `json:"ClusterID"`
	} `json:"*installconfig.clusterID"`
	ICClusterName struct {
		ClusterName string `json:"ClusterName"`
	} `json:"*installconfig.clusterName"`
	ICPlatform struct {
		Libvirt struct {
			URI                    string `json:"URI"`
			DefaultMachinePlatform struct {
				Image string `json:"image"`
			} `json:"defaultMachinePlatform"`
			Network struct {
				If      string `json:"if"`
				IPRange string `json:"ipRange"`
			} `json:"network"`
			MasterIPs interface{} `json:"masterIPs"`
		} `json:"libvirt"`
	} `json:"*installconfig.platform"`
	ICPullSecret struct {
		PullSecret string `json:"PullSecret"`
	} `json:"*installconfig.pullSecret"`
	ICSSHPublicKey struct {
		Key string `json:"Key"`
	} `json:"*installconfig.sshPublicKey"`
}

// DataToValidate as expected base domain and cluster name
type DataToValidate struct {
	BaseDomain    string
	ClusterName   string
	SSH           string
	PullSecret    string
	ConnectionURI string
}

// ValidateInstallStateConfig check that the file content is the same of the env variables
func ValidateInstallStateConfig(filePath string, expectedData DataToValidate) {
	fileName := ".openshift_install_state.json"

	openshiftInstState, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening the file:" + filePath)
		log.Fatalf("Failed while open file: %v", err)
	}

	defer openshiftInstState.Close()
	byteValue, _ := ioutil.ReadAll(openshiftInstState)

	var installState installStateConf
	json.Unmarshal(byteValue, &installState)

	AssertStringContains(installState.InstallConfig.Config.BaseDomain,
		expectedData.BaseDomain,
		"*installconfig.InstallConfig.config.baseDomain not found in "+fileName)
	AssertStringContains(installState.ICBaseDomain.BaseDomain,
		expectedData.BaseDomain,
		"*installconfig.baseDomain.BaseDomain not found in "+fileName)

	AssertStringContains(installState.InstallConfig.Config.Metadata.Name,
		expectedData.ClusterName,
		"*installconfig.config.metadata.name not found in "+fileName)
	AssertStringContains(installState.ICClusterName.ClusterName,
		expectedData.ClusterName,
		"*installconfig.clusterName.ClusterName not found in "+fileName)

	// OPENSHIFT_INSTALL_PLATFORM is implicitly validated by next checks. If its not Libvirt platform, the Name and URI will not be accessible
	AssertStringContains(installState.InstallConfig.Config.Platform.Libvirt.URI,
		expectedData.ConnectionURI,
		"*installconfig.InstallConfig.platform.libvirt.URI not found in "+fileName)
	AssertStringContains(installState.ICPlatform.Libvirt.URI,
		expectedData.ConnectionURI,
		"*installconfig.platform.libvirt.URI not found in "+fileName)

	if strings.Contains(expectedData.SSH, "id_rsa.pub") {
		shhPubFileContent, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub")
		if err != nil {
			fmt.Println("Error: failed reading the ssh file from ~/.ssh/id_rsa.pub")
			fmt.Print(err)
		}
		AssertStringContains(string(installState.InstallConfig.Config.SSHKey),
			string(shhPubFileContent),
			"*installconfig.InstallConfig.config.sshKey not found in "+fileName)
		AssertStringContains(string(installState.ICSSHPublicKey.Key),
			string(shhPubFileContent),
			"*installconfig.sshPublicKey.Key not found in "+fileName)
	} else {
		AssertStringContains(string(installState.InstallConfig.Config.SSHKey),
			"",
			"*installconfig.InstallConfig.config.sshKey not found in "+fileName)
		AssertStringContains(string(installState.ICSSHPublicKey.Key),
			"",
			"*installconfig.sshPublicKey.Key not found in "+fileName)
	}

	AssertStringContains(string(installState.InstallConfig.Config.PullSecret),
		expectedData.PullSecret,
		"*installconfig.InstallConfig.config.pullSecret not found in "+fileName)
	AssertStringContains(string(installState.ICPullSecret.PullSecret),
		expectedData.PullSecret,
		"*installconfig.pullSecret.PullSecret not found in "+fileName)
}
