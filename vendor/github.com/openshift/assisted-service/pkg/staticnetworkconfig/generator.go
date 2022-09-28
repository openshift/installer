package staticnetworkconfig

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/executer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/json"
)

type Config struct {
	MaxConcurrentGenerations int64 `envconfig:"MAX_CONCURRENT_NMSTATECTL_GENERATIONS" default:"30"`
}

type StaticNetworkConfigData struct {
	FilePath     string
	FileContents string
}

//go:generate mockgen -source=generator.go -package=staticnetworkconfig -destination=mock_generator.go
type StaticNetworkConfig interface {
	GenerateStaticNetworkConfigData(ctx context.Context, hostsYAMLS string) ([]StaticNetworkConfigData, error)
	FormatStaticNetworkConfigForDB(staticNetworkConfig []*models.HostStaticNetworkConfig) (string, error)
	ValidateStaticConfigParams(ctx context.Context, staticNetworkConfig []*models.HostStaticNetworkConfig) error
}

type StaticNetworkConfigGenerator struct {
	Config
	log logrus.FieldLogger
	sem *semaphore.Weighted
}

func New(log logrus.FieldLogger, cfg Config) StaticNetworkConfig {
	return &StaticNetworkConfigGenerator{
		Config: cfg,
		log:    log,
		sem:    semaphore.NewWeighted(cfg.MaxConcurrentGenerations)}
}

func (s *StaticNetworkConfigGenerator) GenerateStaticNetworkConfigData(ctx context.Context, staticNetworkConfigStr string) ([]StaticNetworkConfigData, error) {

	staticNetworkConfig, err := s.decodeStaticNetworkConfig(staticNetworkConfigStr)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to decode static network config")
		return nil, err
	}
	s.log.Infof("Start configuring static network for %d hosts", len(staticNetworkConfig))
	filesList := []StaticNetworkConfigData{}
	for i, hostConfig := range staticNetworkConfig {
		hostFileList, err := s.generateHostStaticNetworkConfigData(ctx, hostConfig, fmt.Sprintf("host%d", i))
		if err != nil {
			err = errors.Wrapf(err, "failed to create static config for host %d", i)
			s.log.Error(err)
			return nil, err
		}
		filesList = append(filesList, hostFileList...)
	}
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) generateHostStaticNetworkConfigData(ctx context.Context, hostConfig *models.HostStaticNetworkConfig, hostDir string) ([]StaticNetworkConfigData, error) {
	hostYAML := hostConfig.NetworkYaml
	macInterfaceMapping := s.formatMacInterfaceMap(hostConfig.MacInterfaceMap)
	result, err := s.executeNMStatectl(ctx, hostYAML)
	if err != nil {
		return nil, err
	}
	filesList, err := s.createNMConnectionFiles(result, hostDir)
	if err != nil {
		s.log.WithError(err).Errorf("failed to create NM connection files")
		return nil, err
	}
	mapConfigData := StaticNetworkConfigData{
		FilePath:     filepath.Join(hostDir, "mac_interface.ini"),
		FileContents: macInterfaceMapping,
	}
	filesList = append(filesList, mapConfigData)
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) executeNMStatectl(ctx context.Context, hostYAML string) (string, error) {
	err := s.sem.Acquire(ctx, 1)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to lock semaphore for nmstatectl execution")
		return "", err
	}
	defer s.sem.Release(1)

	executer := &executer.CommonExecuter{}
	f, err := executer.TempFile("", "host-config")
	if err != nil {
		s.log.WithError(err).Errorf("Failed to create temp file")
		return "", err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	_, err = f.WriteString(hostYAML)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to write host config to temp file")
		return "", err
	}
	if err = f.Sync(); err != nil {
		s.log.WithError(err).Warn("Failed to sync file")
	}
	if err = f.Close(); err != nil {
		s.log.WithError(err).Warn("Failed to close file")
	}
	stdout, stderr, retCode := executer.ExecuteWithContext(ctx, "nmstatectl", "gc", f.Name())
	if retCode != 0 {
		s.log.Errorf("<nmstatectl gc> failed, errorCode %d, stderr %s, input yaml <%s>", retCode, stderr, hostYAML)
		errMsg := strings.Split(stderr, "Error:")
		return "", fmt.Errorf("failed to execute 'nmstatectl gc', error: %s", strings.TrimSpace(errMsg[len(errMsg)-1]))
	}
	return stdout, nil
}

// create NMConnectionFiles formats the nmstate output into a list of file data
// Nothing is written to the local filesystem
func (s *StaticNetworkConfigGenerator) createNMConnectionFiles(nmstateOutput, hostDir string) ([]StaticNetworkConfigData, error) {
	var hostNMConnections map[string]interface{}
	err := yaml.Unmarshal([]byte(nmstateOutput), &hostNMConnections)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to unmarshal nmstate output")
		return nil, err
	}
	if _, found := hostNMConnections["NetworkManager"]; !found {
		return nil, errors.Errorf("nmstate generated an empty NetworkManager config file content")
	}
	filesList := []StaticNetworkConfigData{}
	connectionsList := hostNMConnections["NetworkManager"].([]interface{})
	if len(connectionsList) == 0 {
		return nil, errors.Errorf("nmstate generated an empty NetworkManager config file content")
	}
	for _, connection := range connectionsList {
		connectionElems := connection.([]interface{})
		fileName := connectionElems[0].(string)
		fileContents, err := s.formatNMConnection(connectionElems[1].(string))
		if err != nil {
			return nil, err
		}
		s.log.Infof("Adding NMConnection file <%s>", fileName)
		newFile := StaticNetworkConfigData{
			FilePath:     filepath.Join(hostDir, fileName),
			FileContents: fileContents,
		}
		filesList = append(filesList, newFile)
	}
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) formatNMConnection(nmConnection string) (string, error) {
	ini.PrettyFormat = false
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, []byte(nmConnection))
	if err != nil {
		s.log.WithError(err).Errorf("Failed to load the ini format string %s", nmConnection)
		return "", err
	}
	connectionSection := cfg.Section("connection")
	_, err = connectionSection.NewKey("autoconnect", "true")
	if err != nil {
		s.log.WithError(err).Errorf("Failed to add autoconnect key to section connection")
		return "", err
	}
	_, err = connectionSection.NewKey("autoconnect-priority", "1")
	if err != nil {
		s.log.WithError(err).Errorf("Failed to add autoconnect-priority key to section connection")
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = cfg.WriteTo(buf)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to output nmconnection ini file to buffer")
		return "", err
	}
	return buf.String(), nil
}

func (s *StaticNetworkConfigGenerator) ValidateStaticConfigParams(ctx context.Context, staticNetworkConfig []*models.HostStaticNetworkConfig) error {
	var err *multierror.Error
	for i, hostConfig := range staticNetworkConfig {
		err = multierror.Append(err, s.validateMacInterfaceName(i, hostConfig.MacInterfaceMap))
		if validateErr := s.validateNMStateYaml(ctx, hostConfig.NetworkYaml); validateErr != nil {
			err = multierror.Append(err, fmt.Errorf("failed to validate network yaml for host %d, %s", i, validateErr))
		}
	}
	return err.ErrorOrNil()
}

func (s *StaticNetworkConfigGenerator) validateMacInterfaceName(hostIdx int, macInterfaceMap models.MacInterfaceMap) error {
	interfaceCheck := make(map[string]struct{}, len(macInterfaceMap))
	macCheck := make(map[string]struct{}, len(macInterfaceMap))
	for _, macInterface := range macInterfaceMap {
		interfaceCheck[macInterface.LogicalNicName] = struct{}{}
		macCheck[macInterface.MacAddress] = struct{}{}
	}
	if len(interfaceCheck) < len(macInterfaceMap) || len(macCheck) < len(macInterfaceMap) {
		return fmt.Errorf("MACs and Interfaces for host %d must be unique", hostIdx)
	}
	return nil
}

func (s *StaticNetworkConfigGenerator) validateNMStateYaml(ctx context.Context, networkYaml string) error {
	result, err := s.executeNMStatectl(ctx, networkYaml)
	if err != nil {
		return err
	}

	// Check that the file content can be created
	// This doesn't write anything to the local filesystem
	_, err = s.createNMConnectionFiles(result, "temphostdir")
	return err
}

func compareMapInterfaces(intf1, intf2 *models.MacInterfaceMapItems0) bool {
	if intf1.LogicalNicName != intf2.LogicalNicName {
		return intf1.LogicalNicName < intf2.LogicalNicName
	}
	return intf1.MacAddress < intf2.MacAddress
}

func compareMacInterfaceMaps(map1, map2 models.MacInterfaceMap) bool {
	if len(map1) != len(map2) {
		return len(map1) < len(map2)
	}
	for i := range map1 {
		less := compareMapInterfaces(map1[i], map2[i])
		greater := compareMapInterfaces(map2[i], map1[i])
		if less || greater {
			return less
		}
	}
	return false
}

func sortStaticNetworkConfig(staticNetworkConfig []*models.HostStaticNetworkConfig) {
	for i := range staticNetworkConfig {
		item := staticNetworkConfig[i]
		sort.SliceStable(item.MacInterfaceMap, func(i, j int) bool {
			return compareMapInterfaces(item.MacInterfaceMap[i], item.MacInterfaceMap[j])
		})
	}
	sort.SliceStable(staticNetworkConfig, func(i, j int) bool {
		hostConfig1 := staticNetworkConfig[i]
		hostConfig2 := staticNetworkConfig[j]
		if hostConfig1.NetworkYaml != hostConfig2.NetworkYaml {
			return hostConfig1.NetworkYaml < hostConfig2.NetworkYaml
		}
		return compareMacInterfaceMaps(hostConfig1.MacInterfaceMap, hostConfig2.MacInterfaceMap)
	})
}

func (s *StaticNetworkConfigGenerator) FormatStaticNetworkConfigForDB(staticNetworkConfig []*models.HostStaticNetworkConfig) (string, error) {
	if len(staticNetworkConfig) == 0 {
		return "", nil
	}
	sortStaticNetworkConfig(staticNetworkConfig)
	b, err := json.Marshal(&staticNetworkConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to JSON Marshal static network config")
	}
	return string(b), nil
}

func (s *StaticNetworkConfigGenerator) decodeStaticNetworkConfig(staticNetworkConfigStr string) (staticNetworkConfig []*models.HostStaticNetworkConfig, err error) {
	if staticNetworkConfigStr == "" {
		return
	}
	err = json.Unmarshal([]byte(staticNetworkConfigStr), &staticNetworkConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON Unmarshal static network config %s", staticNetworkConfigStr)
	}
	return
}

func (s *StaticNetworkConfigGenerator) formatMacInterfaceMap(macInterfaceMap models.MacInterfaceMap) string {
	lines := make([]string, len(macInterfaceMap))
	for i, entry := range macInterfaceMap {
		lines[i] = fmt.Sprintf("%s=%s", entry.MacAddress, entry.LogicalNicName)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}
