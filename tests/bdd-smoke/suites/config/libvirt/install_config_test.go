package config_test

import (
	"os"
	"testing"

	"github.com/openshift/installer/tests/bdd-smoke/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for create install-config in libvirt")
}

var _ = Describe("Feature: Check that install-config command generates the install configuration asset for libvirt", func() {

	openshiftInstallerPath := os.Getenv("GOPATH") + "/src/github.com/openshift/installer"
	openshiftInstallBinPath := openshiftInstallerPath + "/bin/openshift-install"

	AfterSuite(func() {
		args := []string{"-rf",
			"scenario1/",
			"scenario3/",
			"scenario4/",
			"scenario5/",
			"tests/bdd-smoke/suites/config/libvirt/scenario6/",
			".openshift_install_state.json",
			".openshift_install.log",
			"install-config.yaml"}
		utils.RunOSCommandWithArgs("rm", args, openshiftInstallerPath)
	})

	// tests for openshift-install create install-config
	Context("Scenario: Generate configuration asset specifying a directory", func() {
		testData := &utils.DataToValidate{
			BaseDomain:    "tt.testing",
			ClusterName:   "test1",
			SSH:           "id_rsa.pub",
			PullSecret:    `{"auths": {"quay.io": {"auth": "aValidPullSecret","email": ""}}}`,
			ConnectionURI: "qemu+tcp://192.168.122.1/system",
		}

		It("Given that the openshift-install is compiled", func() {
			Ω(openshiftInstallBinPath).Should(BeAnExistingFile())
		})
		It("When I execute openshift-install create install-config --dir=scenario1", func() {
			args := []string{"create", "install-config", "--dir=scenario1"}
			commandLineArgs := [][]string{
				{"SSH Public Key", testData.SSH},
				{"Platform", "libvirt"},
				{"Libvirt Connection URI", testData.ConnectionURI},
				{"Base Domain", testData.BaseDomain},
				{"Cluster Name", testData.ClusterName},
				{"Pull Secret", testData.PullSecret},
			}
			cmdOut := utils.RunInstallerWithSurvey(args, openshiftInstallerPath, commandLineArgs)
			Ω(cmdOut).Should(ContainSubstring(".ssh/id_rsa.pub"))
			Ω(cmdOut).Should(ContainSubstring("Platform libvirt"))
			Ω(cmdOut).Should(ContainSubstring("Libvirt Connection URI qemu+tcp://192.168.122.1/system"))
			Ω(cmdOut).Should(ContainSubstring("Base Domain tt.testing"))
			Ω(cmdOut).Should(ContainSubstring("Cluster Name test1"))
		})
		Specify("Then I expect to see the install-config.yaml file in the path ./scenario1/", func() {
			Ω(openshiftInstallerPath + "/scenario1/install-config.yaml").Should(BeAnExistingFile())
		})
		Specify("And I expect install-config.yaml file contains the info specified in the environment variables", func() {
			utils.ValidateInstallConfig(openshiftInstallerPath+"/scenario1/install-config.yaml", *testData)
		})
		Specify("And I expect .openshift_install_state.json file contains the info specified in the environment variables", func() {
			utils.ValidateInstallStateConfig(openshiftInstallerPath+"/scenario1/.openshift_install_state.json", *testData)
		})
	})

	Context("Scenario: Generate configuration asset without specifying a directory with log level in debug", func() {
		testData := &utils.DataToValidate{
			BaseDomain:    "tt.testing",
			ClusterName:   "test1",
			SSH:           "id_rsa.pub",
			PullSecret:    `{"auths": {"quay.io": {"auth": "aValidPullSecret","email": ""}}}`,
			ConnectionURI: "qemu+tcp://192.168.122.1/system",
		}
		It("Given that the openshift-install is compiled", func() {
			Ω(openshiftInstallBinPath).Should(BeAnExistingFile())
		})
		It("When I execute bin/openshift-install install-config --log-level=debug", func() {
			args := []string{"create", "install-config", "--log-level=debug"}
			commandLineArgs := [][]string{
				{"SSH Public Key", testData.SSH},
				{"Platform", "libvirt"},
				{"Libvirt Connection URI", testData.ConnectionURI},
				{"Base Domain", testData.BaseDomain},
				{"Cluster Name", testData.ClusterName},
				{"Pull Secret", testData.PullSecret},
			}
			cmdOut := utils.RunInstallerWithSurvey(args, openshiftInstallerPath, commandLineArgs)
			Ω(cmdOut).Should(ContainSubstring(".ssh/id_rsa.pub"))
			Ω(cmdOut).Should(ContainSubstring("Platform libvirt"))
			Ω(cmdOut).Should(ContainSubstring("Libvirt Connection URI qemu+tcp://192.168.122.1/system"))
			Ω(cmdOut).Should(ContainSubstring("Base Domain tt.testing"))
			Ω(cmdOut).Should(ContainSubstring("Cluster Name test1"))
		})
		Specify("Then I expect to see the install-config.yaml file in the path ./", func() {
			Ω(openshiftInstallerPath + "/install-config.yaml").Should(BeAnExistingFile())
		})
		Specify("And I expect to see the .openshift_install_state.json file in the path ./", func() {
			Ω(openshiftInstallerPath + "/.openshift_install_state.json").Should(BeAnExistingFile())
		})
		Specify("And I expect install-config.yaml file contains the info specified in the environment variables", func() {
			utils.ValidateInstallConfig(openshiftInstallerPath+"/install-config.yaml", *testData)
		})
		Specify("And I expect .openshift_install_state.json file contains the info specified in the environment variables", func() {
			utils.ValidateInstallStateConfig(openshiftInstallerPath+"/.openshift_install_state.json", *testData)
		})
	})

	Context("Scenario: Generate configuration asset specifying a directory with log level in error", func() {
		cmdOut := ""
		testData := &utils.DataToValidate{
			BaseDomain:    "tt.testing",
			ClusterName:   "test1",
			SSH:           "id_rsa.pub",
			PullSecret:    `{"auths": {"quay.io": {"auth": "aValidPullSecret","email": ""}}}`,
			ConnectionURI: "qemu+tcp://192.168.122.1/system",
		}
		It("Given that the openshift-install is compiled", func() {
			Ω(openshiftInstallBinPath).Should(BeAnExistingFile())
		})
		It("When I execute bin/openshift-install install-config --dir=scenario3 --log-level=error", func() {
			args := []string{"create", "install-config", "--dir=scenario3", "--log-level=error"}
			commandLineArgs := [][]string{
				{"SSH Public Key", testData.SSH},
				{"Platform", "libvirt"},
				{"Libvirt Connection URI", testData.ConnectionURI},
				{"Base Domain", testData.BaseDomain},
				{"Cluster Name", testData.ClusterName},
				{"Pull Secret", testData.PullSecret},
			}
			cmdOut := utils.RunInstallerWithSurvey(args, openshiftInstallerPath, commandLineArgs)
			Ω(cmdOut).Should(ContainSubstring(".ssh/id_rsa.pub"))
			Ω(cmdOut).Should(ContainSubstring("Platform libvirt"))
			Ω(cmdOut).Should(ContainSubstring("Libvirt Connection URI qemu+tcp://192.168.122.1/system"))
			Ω(cmdOut).Should(ContainSubstring("Base Domain tt.testing"))
			Ω(cmdOut).Should(ContainSubstring("Cluster Name test1"))
		})
		Specify("Then I expect to see the install-config.yaml file in the path ./scenario3", func() {
			Ω(openshiftInstallerPath + "/scenario3/install-config.yaml").Should(BeAnExistingFile())
		})
		Specify("And I expect install-config.yaml file contains the info specified in the environment variables", func() {
			utils.ValidateInstallConfig(openshiftInstallerPath+"/scenario3/install-config.yaml", *testData)
		})
		Specify("And I expect .openshift_install_state.json file contains the info specified in the environment variables", func() {
			utils.ValidateInstallStateConfig(openshiftInstallerPath+"/scenario3/.openshift_install_state.json", *testData)
		})
		Specify("And I expect to not to get any error in output log from the console", func() {
			Ω(cmdOut).Should(Equal(""))
		})
	})

	Context("Scenario: Generate configuration asset specifying a directory running it from other path", func() {
		cmdOut := ""
		testData := &utils.DataToValidate{
			BaseDomain:    "tt.testing",
			ClusterName:   "test1",
			SSH:           "id_rsa.pub",
			PullSecret:    `{"auths": {"quay.io": {"auth": "aValidPullSecret","email": ""}}}`,
			ConnectionURI: "qemu+tcp://192.168.122.1/system",
		}
		It("Given that the openshift-install is compiled", func() {
			Ω(openshiftInstallBinPath).Should(BeAnExistingFile())
		})
		It("When I execute bin/openshift-install install-config --dir=scenario6 --log-level=error", func() {
			args := []string{"create", "install-config", "--dir=scenario6"}
			commandLineArgs := [][]string{
				{"SSH Public Key", testData.SSH},
				{"Platform", "libvirt"},
				{"Libvirt Connection URI", testData.ConnectionURI},
				{"Base Domain", testData.BaseDomain},
				{"Cluster Name", testData.ClusterName},
				{"Pull Secret", testData.PullSecret},
			}
			cmdOut := utils.RunInstallerWithSurvey(args, ".", commandLineArgs)
			Ω(cmdOut).Should(ContainSubstring(".ssh/id_rsa.pub"))
			Ω(cmdOut).Should(ContainSubstring("Platform libvirt"))
			Ω(cmdOut).Should(ContainSubstring("Libvirt Connection URI qemu+tcp://192.168.122.1/system"))
			Ω(cmdOut).Should(ContainSubstring("Base Domain tt.testing"))
			Ω(cmdOut).Should(ContainSubstring("Cluster Name test1"))
		})
		Specify("Then I expect to see the install-config.yaml file in the path ./scenario6", func() {
			Ω("scenario6/install-config.yaml").Should(BeAnExistingFile())
		})
		Specify("And I expect install-config.yaml file contains the info specified in the environment variables", func() {
			utils.ValidateInstallConfig("scenario6/install-config.yaml", *testData)
		})
		Specify("And I expect .openshift_install_state.json file contains the info specified in the environment variables", func() {
			utils.ValidateInstallStateConfig("scenario6/.openshift_install_state.json", *testData)
		})
		Specify("And I expect to not to get any error in output log from the console", func() {
			Ω(cmdOut).Should(Equal(""))
		})
	})

	Context("Scenario: Use previously generated configuration asset without specifying a directory with log level in debug", func() {
		testData := &utils.DataToValidate{
			BaseDomain:    "tt.testing",
			ClusterName:   "test1",
			SSH:           "id_rsa.pub",
			PullSecret:    `{"auths": {"quay.io": {"auth": "aValidPullSecret","email": ""}}}`,
			ConnectionURI: "qemu+tcp://192.168.122.1/system",
		}
		It("Given that the openshift-install is compiled", func() {
			Ω(openshiftInstallBinPath).Should(BeAnExistingFile())
		})
		It("When I execute bin/openshift-install install-config --log-level=debug", func() {
			args := []string{"create", "install-config", "--log-level=debug"}
			commandLineArgs := [][]string{
				{"SSH Public Key", testData.SSH},
				{"Platform", "libvirt"},
				{"Libvirt Connection URI", testData.ConnectionURI},
				{"Base Domain", testData.BaseDomain},
				{"Cluster Name", testData.ClusterName},
				{"Pull Secret", testData.PullSecret},
			}
			cmdOut := utils.RunInstallerWithSurvey(args, openshiftInstallerPath, commandLineArgs)
			Ω(cmdOut).Should(ContainSubstring(""))
		})
		Specify("Then I expect to see the install-config.yaml file in the path ./", func() {
			Ω(openshiftInstallerPath + "/install-config.yaml").Should(BeAnExistingFile())
		})
		Specify("And I expect to see the .openshift_install_state.json file in the path ./", func() {
			Ω(openshiftInstallerPath + "/.openshift_install_state.json").Should(BeAnExistingFile())
		})
		Specify("And I expect install-config.yaml file contains the info specified in the environment variables", func() {
			utils.ValidateInstallConfig(openshiftInstallerPath+"/install-config.yaml", *testData)
		})
		Specify("And I expect .openshift_install_state.json file contains the info specified in the environment variables", func() {
			utils.ValidateInstallStateConfig(openshiftInstallerPath+"/.openshift_install_state.json", *testData)
		})
	})
})
