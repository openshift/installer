package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type installConf struct {
	APIVersion string `yaml:"apiVersion"`
	BaseDomain string `yaml:"baseDomain"`
	Machines   []struct {
		Name     string `yaml:"name"`
		Platform struct {
		} `yaml:"platform"`
		Replicas int `yaml:"replicas"`
	} `yaml:"machines"`
	Metadata struct {
		CreationTimestamp interface{} `yaml:"creationTimestamp"`
		Name              string      `yaml:"name"`
	} `yaml:"metadata"`
	Networking struct {
		ClusterNetworks []struct {
			Cidr             string `yaml:"cidr"`
			HostSubnetLength int    `yaml:"hostSubnetLength"`
		} `yaml:"clusterNetworks"`
		MachineCIDR string `yaml:"machineCIDR"`
		ServiceCIDR string `yaml:"serviceCIDR"`
		Type        string `yaml:"type"`
	} `yaml:"networking"`
	Platform struct {
		Libvirt struct {
			URI     string `yaml:"URI"`
			Network struct {
				If string `yaml:"if"`
			} `yaml:"network"`
		} `yaml:"libvirt"`
	} `yaml:"platform"`
	PullSecret string `yaml:"pullSecret"`
	SSHKey     string `yaml:"sshKey"`
}

func (c *installConf) readInstallConf(filePath string) *installConf {

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

//ValidateInstallConfig check that the file content is the same of the env variables
func ValidateInstallConfig(filePath string, expectedData DataToValidate) {
	fileName := "install-config.yaml"
	var installConfiguration installConf
	installConfiguration.readInstallConf(filePath)

	AssertStringContains(installConfiguration.BaseDomain,
		expectedData.BaseDomain,
		"baseDomain not found in "+fileName)

	// OPENSHIFT_INSTALL_PLATFORM is implicitly validated by next checks. If its not Libvirt platform, the Name and URI will not be accessible
	AssertStringContains(installConfiguration.Platform.Libvirt.URI,
		expectedData.ConnectionURI,
		"platform.libvirt.URI not found in "+fileName)

	if strings.Contains(expectedData.SSH, "id_rsa.pub") {
		shhPubFileContent, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub")
		if err != nil {
			fmt.Println("Error: failed reading the ssh file from ~/.ssh/id_rsa.pub")
			fmt.Print(err)
		}

		AssertStringContains(string(installConfiguration.SSHKey),
			string(shhPubFileContent),
			"sshKey not found in "+fileName)
	} else {
		AssertStringContains(string(installConfiguration.SSHKey),
			"",
			"*installconfig.InstallConfig.config.sshKey not found in "+fileName)
	}

	AssertStringContains(string(installConfiguration.PullSecret),
		expectedData.PullSecret,
		"pullSecret not found in "+fileName)
}
