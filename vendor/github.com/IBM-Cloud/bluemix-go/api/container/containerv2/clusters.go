package containerv2

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

//ClusterCreateRequest ...
type ClusterCreateRequest struct {
	DisablePublicServiceEndpoint bool             `json:"disablePublicServiceEndpoint"`
	KubeVersion                  string           `json:"kubeVersion" description:"kubeversion of cluster"`
	Billing                      string           `json:"billing,omitempty"`
	PodSubnet                    string           `json:"podSubnet"`
	Provider                     string           `json:"provider"`
	ServiceSubnet                string           `json:"serviceSubnet"`
	Name                         string           `json:"name" binding:"required" description:"The cluster's name"`
	DefaultWorkerPoolEntitlement string           `json:"defaultWorkerPoolEntitlement"`
	CosInstanceCRN               string           `json:"cosInstanceCRN"`
	WorkerPools                  WorkerPoolConfig `json:"workerPool"`
}

type WorkerPoolConfig struct {
	DiskEncryption bool              `json:"diskEncryption,omitempty"`
	Entitlement    string            `json:"entitlement"`
	Flavor         string            `json:"flavor"`
	Isolation      string            `json:"isolation,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	Name           string            `json:"name" binding:"required" description:"The workerpool's name"`
	VpcID          string            `json:"vpcID"`
	WorkerCount    int               `json:"workerCount"`
	Zones          []Zone            `json:"zones"`
}

// type Label struct {
// 	AdditionalProp1 string `json:"additionalProp1,omitempty"`
// 	AdditionalProp2 string `json:"additionalProp2,omitempty"`
// 	AdditionalProp3 string `json:"additionalProp3,omitempty"`
// }

type Zone struct {
	ID       string `json:"id,omitempty" description:"The id"`
	SubnetID string `json:"subnetID,omitempty"`
}

//ClusterInfo ...
type ClusterInfo struct {
	CreatedDate       string        `json:"createdDate"`
	DataCenter        string        `json:"dataCenter"`
	ID                string        `json:"id"`
	Location          string        `json:"location"`
	Entitlement       string        `json:"entitlement"`
	MasterKubeVersion string        `json:"masterKubeVersion"`
	Name              string        `json:"name"`
	Region            string        `json:"region"`
	ResourceGroupID   string        `json:"resourceGroup"`
	State             string        `json:"state"`
	IsPaid            bool          `json:"isPaid"`
	Addons            []Addon       `json:"addons"`
	OwnerEmail        string        `json:"ownerEmail"`
	Type              string        `json:"type"`
	TargetVersion     string        `json:"targetVersion"`
	ServiceSubnet     string        `json:"serviceSubnet"`
	ResourceGroupName string        `json:"resourceGroupName"`
	Provider          string        `json:"provider"`
	PodSubnet         string        `json:"podSubnet"`
	MultiAzCapable    bool          `json:"multiAzCapable"`
	APIUser           string        `json:"apiUser"`
	ServerURL         string        `json:"serverURL"`
	MasterURL         string        `json:"masterURL"`
	DisableAutoUpdate bool          `json:"disableAutoUpdate"`
	WorkerZones       []string      `json:"workerZones"`
	Vpcs              []string      `json:"vpcs"`
	CRN               string        `json:"crn"`
	VersionEOS        string        `json:"versionEOS"`
	ServiceEndpoints  Endpoints     `json:"serviceEndpoints"`
	Lifecycle         LifeCycleInfo `json:"lifecycle"`
	WorkerCount       int           `json:"workerCount"`
	Ingress           IngresInfo    `json:"ingress"`
	Features          Feat          `json:"features"`
}
type Feat struct {
	KeyProtectEnabled bool `json:"keyProtectEnabled"`
	PullSecretApplied bool `json:"pullSecretApplied"`
}
type IngresInfo struct {
	HostName   string `json:"hostname"`
	SecretName string `json:"secretName"`
}
type LifeCycleInfo struct {
	ModifiedDate             string `json:"modifiedDate"`
	MasterStatus             string `json:"masterStatus"`
	MasterStatusModifiedDate string `json:"masterStatusModifiedDate"`
	MasterHealth             string `json:"masterHealth"`
	MasterState              string `json:"masterState"`
}

//ClusterTargetHeader ...
type ClusterTargetHeader struct {
	AccountID     string
	ResourceGroup string
	Provider      string // supported providers e.g vpc-classic , vpc-gen2, satellite
}
type Endpoints struct {
	PrivateServiceEndpointEnabled bool   `json:"privateServiceEndpointEnabled"`
	PrivateServiceEndpointURL     string `json:"privateServiceEndpointURL"`
	PublicServiceEndpointEnabled  bool   `json:"publicServiceEndpointEnabled"`
	PublicServiceEndpointURL      string `json:"publicServiceEndpointURL"`
}

type Addon struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

//ClusterCreateResponse ...
type ClusterCreateResponse struct {
	ID string `json:"clusterID"`
}

//Clusters interface
type Clusters interface {
	Create(params ClusterCreateRequest, target ClusterTargetHeader) (ClusterCreateResponse, error)
	List(target ClusterTargetHeader) ([]ClusterInfo, error)
	Delete(name string, target ClusterTargetHeader, deleteDependencies ...bool) error
	GetCluster(name string, target ClusterTargetHeader) (*ClusterInfo, error)
	GetClusterConfigDetail(name, homeDir string, admin bool, target ClusterTargetHeader) (containerv1.ClusterKeyInfo, error)
	StoreConfigDetail(name, baseDir string, admin bool, createCalicoConfig bool, target ClusterTargetHeader) (string, containerv1.ClusterKeyInfo, error)

	//TODO Add other opertaions
}
type clusters struct {
	client     *client.Client
	pathPrefix string
}

const (
	accountIDHeader     = "X-Auth-Resource-Account"
	resourceGroupHeader = "X-Auth-Resource-Group"
)

//ToMap ...
func (c ClusterTargetHeader) ToMap() map[string]string {
	m := make(map[string]string, 3)
	m[accountIDHeader] = c.AccountID
	m[resourceGroupHeader] = c.ResourceGroup
	return m
}

func newClusterAPI(c *client.Client) Clusters {
	return &clusters{
		client: c,
		//pathPrefix: "/v2/vpc/",
	}
}

//List ...
func (r *clusters) List(target ClusterTargetHeader) ([]ClusterInfo, error) {
	clusters := []ClusterInfo{}
	var err error
	if target.Provider != "satellite" {
		getClustersPath := "/v2/vpc/getClusters"
		if len(target.Provider) > 0 {
			getClustersPath = fmt.Sprintf(getClustersPath+"?provider=%s", url.QueryEscape(target.Provider))
		}
		_, err := r.client.Get(getClustersPath, &clusters, target.ToMap())
		if err != nil {
			return nil, err
		}
	}
	if len(target.Provider) == 0 || target.Provider == "satellite" {
		// get satellite clusters
		satelliteClusters := []ClusterInfo{}
		_, err = r.client.Get("/v2/satellite/getClusters", &satelliteClusters, target.ToMap())
		if err != nil && target.Provider == "satellite" {
			// return error only when provider is satellite. Else ignore error and return VPC clusters
			trace.Logger.Println("Unable to get the satellite clusters ", err)
			return nil, err
		}
		clusters = append(clusters, satelliteClusters...)
	}
	return clusters, nil
}

//Create ...
func (r *clusters) Create(params ClusterCreateRequest, target ClusterTargetHeader) (ClusterCreateResponse, error) {
	var cluster ClusterCreateResponse
	_, err := r.client.Post("/v2/vpc/createCluster", params, &cluster, target.ToMap())
	return cluster, err
}

//Delete ...
func (r *clusters) Delete(name string, target ClusterTargetHeader, deleteDependencies ...bool) error {
	var rawURL string
	if len(deleteDependencies) != 0 {
		rawURL = fmt.Sprintf("/v1/clusters/%s?deleteResources=%t", name, deleteDependencies[0])
	} else {
		rawURL = fmt.Sprintf("/v1/clusters/%s", name)
	}
	_, err := r.client.Delete(rawURL, target.ToMap())
	return err
}

//GetClusterByIDorName
func (r *clusters) GetCluster(name string, target ClusterTargetHeader) (*ClusterInfo, error) {
	ClusterInfo := &ClusterInfo{}
	rawURL := fmt.Sprintf("/v2/vpc/getCluster?cluster=%s", name)
	_, err := r.client.Get(rawURL, &ClusterInfo, target.ToMap())
	if err != nil {
		return nil, err
	}
	return ClusterInfo, err
}
func (r *ClusterInfo) IsStagingSatelliteCluster() bool {
	return strings.Index(r.ServerURL, "stg") > 0 && r.Provider == "satellite"
}

//FindWithOutShowResourcesCompatible ...
func (r *clusters) FindWithOutShowResourcesCompatible(name string, target ClusterTargetHeader) (ClusterInfo, error) {
	rawURL := fmt.Sprintf("/v2/getCluster?v1-compatible&cluster=%s", name)
	cluster := ClusterInfo{}
	_, err := r.client.Get(rawURL, &cluster, target.ToMap())
	if err != nil {
		return cluster, err
	}
	// Handle VPC cluster.  ServerURL is blank for v2/vpc clusters
	if cluster.ServerURL == "" {
		cluster.ServerURL = cluster.MasterURL
	}
	return cluster, err
}

//GetClusterConfigDetail ...
func (r *clusters) GetClusterConfigDetail(name, dir string, admin bool, target ClusterTargetHeader) (containerv1.ClusterKeyInfo, error) {
	clusterkey := containerv1.ClusterKeyInfo{}
	// Block to add token for openshift clusters (This can be temporary until iks team handles openshift clusters)
	clusterInfo, err := r.FindWithOutShowResourcesCompatible(name, target)
	if err != nil {
		// Assuming an error means that this is a vpc cluster, and we're returning existing kubeconfig
		// When we add support for vpcs on openshift clusters, we may want revisit this
		return clusterkey, err
	}

	if !helpers.FileExists(dir) {
		return clusterkey, fmt.Errorf("Path: %q, to download the config doesn't exist", dir)
	}
	postBody := map[string]interface{}{
		"cluster": name,
		"format":  "zip",
	}
	rawURL := fmt.Sprintf("/v2/applyRBACAndGetKubeconfig")
	if admin {
		postBody["admin"] = true
	}
	if clusterInfo.Provider == "satellite" {
		postBody["endpointType"] = "link"
		postBody["admin"] = true
	}
	resultDir := containerv1.ComputeClusterConfigDir(dir, name, admin)
	const kubeConfigName = "config.yml"
	err = os.MkdirAll(resultDir, 0755)
	if err != nil {
		return clusterkey, fmt.Errorf("Error creating directory to download the cluster config")
	}
	downloadPath := filepath.Join(resultDir, "config.zip")
	trace.Logger.Println("Will download the kubeconfig at", downloadPath)

	var out *os.File
	if out, err = os.Create(downloadPath); err != nil {
		return clusterkey, err
	}
	defer out.Close()
	defer helpers.RemoveFile(downloadPath)
	_, err = r.client.Post(rawURL, postBody, out, target.ToMap())
	if err != nil {
		return clusterkey, err
	}
	trace.Logger.Println("Downloaded the kubeconfig at", downloadPath)
	if err = helpers.Unzip(downloadPath, resultDir); err != nil {
		return clusterkey, err
	}
	defer helpers.RemoveFilesWithPattern(resultDir, "[^(.yml)|(.pem)]$")
	var kubeyml string
	files, _ := ioutil.ReadDir(resultDir)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".zip") {
			fileContent, _ := ioutil.ReadFile(resultDir + "/" + f.Name())
			if f.Name() == "admin-key.pem" {
				clusterkey.AdminKey = string(fileContent)
			}
			if f.Name() == "admin.pem" {
				clusterkey.Admin = string(fileContent)
			}
			if strings.HasPrefix(f.Name(), "ca") && strings.HasSuffix(f.Name(), ".pem") {
				clusterkey.ClusterCACertificate = string(fileContent)
			}
			old := filepath.Join(resultDir, f.Name())
			new := filepath.Join(resultDir, f.Name())
			if strings.HasSuffix(f.Name(), ".yaml") {
				new = filepath.Join(path.Clean(resultDir), "/", path.Clean(kubeConfigName))
				kubeyml = new
			}
			err := os.Rename(old, new)
			if err != nil {
				return clusterkey, fmt.Errorf("Couldn't rename: %q", err)
			}
		}
	}
	if resultDir == "" {
		return clusterkey, errors.New("Unable to locate kube config in zip archive")
	}

	kubefile, _ := ioutil.ReadFile(kubeyml)
	var yamlConfig containerv1.ConfigFile
	err = yaml.Unmarshal(kubefile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}
	if len(yamlConfig.Clusters) != 0 {
		clusterkey.Host = yamlConfig.Clusters[0].Cluster.Server
	}
	if len(yamlConfig.Users) != 0 {
		clusterkey.Token = yamlConfig.Users[0].User.AuthProvider.Config.IDToken
	}

	// Block to add token for openshift clusters (This can be temporary until iks team handles openshift clusters)
	clusterInfo, err = r.FindWithOutShowResourcesCompatible(name, target)
	if err != nil {
		// Assuming an error means that this is a vpc cluster, and we're returning existing kubeconfig
		// When we add support for vpcs on openshift clusters, we may want revisit this
		clusterkey.FilePath, _ = filepath.Abs(kubeyml)
		return clusterkey, err
	}
	if clusterInfo.Type == "openshift" && clusterInfo.Provider != "satellite" {
		trace.Logger.Println("Debug: type is openshift trying login to get token")
		var yamlConfig []byte
		if yamlConfig, err = ioutil.ReadFile(kubeyml); err != nil {
			return clusterkey, err
		}
		yamlConfig, err = r.FetchOCTokenForKubeConfig(yamlConfig, &clusterInfo, clusterInfo.IsStagingSatelliteCluster())
		if err != nil {
			return clusterkey, err
		}
		err = ioutil.WriteFile(kubeyml, yamlConfig, 0644) // 0644 is irrelevant here, since file already exists.
		if err != nil {
			return clusterkey, err
		}
		openshiftyml, _ := ioutil.ReadFile(kubeyml)
		var openshiftyaml containerv1.ConfigFileOpenshift
		err = yaml.Unmarshal(openshiftyml, &openshiftyaml)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
		}
		openshiftusers := openshiftyaml.Users
		for _, usr := range openshiftusers {
			if strings.HasPrefix(usr.Name, "IAM") {
				clusterkey.Token = usr.User.Token
			}
		}
		if len(openshiftyaml.Clusters) != 0 {
			clusterkey.Host = openshiftyaml.Clusters[0].Cluster.Server
		}
		clusterkey.ClusterCACertificate = ""

	}
	clusterkey.FilePath, _ = filepath.Abs(kubeyml)
	return clusterkey, err
}

//StoreConfigDetail ...
func (r *clusters) StoreConfigDetail(name, dir string, admin, createCalicoConfig bool, target ClusterTargetHeader) (string, containerv1.ClusterKeyInfo, error) {
	clusterkey := containerv1.ClusterKeyInfo{}
	clusterInfo, err := r.FindWithOutShowResourcesCompatible(name, target)
	if err != nil {
		return "", clusterkey, err
	}
	postBody := map[string]interface{}{
		"cluster": name,
		"format":  "zip",
	}

	var calicoConfig string
	if !helpers.FileExists(dir) {
		return "", clusterkey, fmt.Errorf("Path: %q, to download the config doesn't exist", dir)
	}
	rawURL := fmt.Sprintf("/v2/applyRBACAndGetKubeconfig")
	if admin {
		postBody["admin"] = true
	}
	if clusterInfo.Provider == "satellite" {
		postBody["endpointType"] = "link"
		postBody["admin"] = true
	}
	if createCalicoConfig {
		postBody["network"] = true
	}
	resultDir := containerv1.ComputeClusterConfigDir(dir, name, admin)
	err = os.MkdirAll(resultDir, 0755)
	if err != nil {
		return "", clusterkey, fmt.Errorf("Error creating directory to download the cluster config")
	}
	downloadPath := filepath.Join(resultDir, "config.zip")
	trace.Logger.Println("Will download the kubeconfig at", downloadPath)

	var out *os.File
	if out, err = os.Create(downloadPath); err != nil {
		return "", clusterkey, err
	}
	defer out.Close()
	defer helpers.RemoveFile(downloadPath)
	_, err = r.client.Post(rawURL, postBody, out, target.ToMap())
	if err != nil {
		return "", clusterkey, err
	}
	trace.Logger.Println("Downloaded the kubeconfig at", downloadPath)
	if err = helpers.Unzip(downloadPath, resultDir); err != nil {
		return "", clusterkey, err
	}
	trace.Logger.Println("Downloaded the kubec", resultDir)

	unzipConfigPath := resultDir
	trace.Logger.Println("Located unzipped directory: ", unzipConfigPath)
	files, _ := ioutil.ReadDir(unzipConfigPath)
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".zip") {
			fileContent, _ := ioutil.ReadFile(unzipConfigPath + "/" + f.Name())
			if f.Name() == "admin-key.pem" {
				clusterkey.AdminKey = string(fileContent)
			}
			if f.Name() == "admin.pem" {
				clusterkey.Admin = string(fileContent)
			}
			if strings.HasPrefix(f.Name(), "ca") && strings.HasSuffix(f.Name(), ".pem") {
				clusterkey.ClusterCACertificate = string(fileContent)
			}
			old := filepath.Join(unzipConfigPath, f.Name())
			new := filepath.Join(unzipConfigPath, f.Name())
			err := os.Rename(old, new)
			if err != nil {
				return "", clusterkey, fmt.Errorf("Couldn't rename: %q", err)
			}
		}
	}
	baseDirFiles, err := ioutil.ReadDir(resultDir)
	if err != nil {
		return "", clusterkey, err
	}

	if createCalicoConfig {
		// Proccess calico golang template file if it exists
		calicoConfig, err = containerv1.GenerateCalicoConfig(resultDir)
		if err != nil {
			return "", clusterkey, err
		}
	}
	var kubeconfigFileName string
	for _, baseDirFile := range baseDirFiles {
		if strings.Contains(baseDirFile.Name(), ".yaml") {
			kubeconfigFileName = fmt.Sprintf("%s/%s", resultDir, baseDirFile.Name())
			break
		}
	}
	if kubeconfigFileName == "" {
		return "", clusterkey, errors.New("Unable to locate kube config in zip archive")
	}
	kubefile, _ := ioutil.ReadFile(kubeconfigFileName)
	var yamlConfig containerv1.ConfigFile
	err = yaml.Unmarshal(kubefile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}
	if len(yamlConfig.Clusters) != 0 {
		clusterkey.Host = yamlConfig.Clusters[0].Cluster.Server
	}
	if len(yamlConfig.Users) != 0 {
		clusterkey.Token = yamlConfig.Users[0].User.AuthProvider.Config.IDToken
	}

	// Block to add token for openshift clusters (This can be temporary until iks team handles openshift clusters)
	clusterInfo, err = r.FindWithOutShowResourcesCompatible(name, target)
	if err != nil {
		// Assuming an error means that this is a vpc cluster, and we're returning existing kubeconfig
		// When we add support for vpcs on openshift clusters, we may want revisit this
		clusterkey.FilePath = kubeconfigFileName
		return calicoConfig, clusterkey, nil
	}

	if clusterInfo.Type == "openshift" && clusterInfo.Provider != "satellite" {
		trace.Logger.Println("Cluster Type is openshift trying login to get token")
		var yamlConfig []byte
		if yamlConfig, err = ioutil.ReadFile(kubeconfigFileName); err != nil {
			return "", clusterkey, err
		}
		yamlConfig, err = r.FetchOCTokenForKubeConfig(yamlConfig, &clusterInfo, clusterInfo.IsStagingSatelliteCluster())
		if err != nil {
			return "", clusterkey, err
		}
		err = ioutil.WriteFile(kubeconfigFileName, yamlConfig, 0644) // check about permissions and truncate
		if err != nil {
			return "", clusterkey, err
		}
		openshiftyml, _ := ioutil.ReadFile(kubeconfigFileName)
		var openshiftyaml containerv1.ConfigFileOpenshift
		err = yaml.Unmarshal(openshiftyml, &openshiftyaml)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
		}
		openshiftusers := openshiftyaml.Users
		for _, usr := range openshiftusers {
			if strings.HasPrefix(usr.Name, "IAM") {
				clusterkey.Token = usr.User.Token
			}
		}
		if len(openshiftyaml.Clusters) != 0 {
			clusterkey.Host = openshiftyaml.Clusters[0].Cluster.Server
		}
		clusterkey.ClusterCACertificate = ""

	}
	clusterkey.FilePath = kubeconfigFileName
	return calicoConfig, clusterkey, nil
}
