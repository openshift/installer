/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/blang/semver/v4"
	ignV3Types "github.com/coreos/ignition/v2/config/v3_4/types"
	"github.com/go-logr/logr"
	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	cosSession "github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/cos"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	ignV2Types "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/ignition"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/options"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

const cosURLDomain = "cloud-object-storage.appdomain.cloud"

// PowerVSMachineScopeParams defines the input parameters used to create a new PowerVSMachineScope.
type PowerVSMachineScopeParams struct {
	Logger            logr.Logger
	Client            client.Client
	Cluster           *clusterv1.Cluster
	Machine           *clusterv1.Machine
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1.IBMPowerVSMachine
	IBMPowerVSImage   *infrav1.IBMPowerVSImage
	ServiceEndpoint   []endpoints.ServiceEndpoint
	DHCPIPCacheStore  cache.Store
}

// PowerVSMachineScope defines a scope defined around a Power VS Machine.
type PowerVSMachineScope struct {
	Client client.Client

	IBMPowerVSClient  powervs.PowerVS
	IBMVPCClient      vpc.Vpc
	ResourceClient    resourcecontroller.ResourceController
	Cluster           *clusterv1.Cluster
	Machine           *clusterv1.Machine
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1.IBMPowerVSMachine
	IBMPowerVSImage   *infrav1.IBMPowerVSImage
	ServiceEndpoint   []endpoints.ServiceEndpoint
	DHCPIPCacheStore  cache.Store
}

// NewPowerVSMachineScope creates a new PowerVSMachineScope from the supplied parameters.
func NewPowerVSMachineScope(params PowerVSMachineScopeParams) (scope *PowerVSMachineScope, err error) { //nolint:gocyclo
	scope = &PowerVSMachineScope{}

	if params.Client == nil {
		err = errors.New("client is required when creating a MachineScope")
		return nil, err
	}
	scope.Client = params.Client

	if params.Machine == nil {
		err = errors.New("machine is required when creating a MachineScope")
		return nil, err
	}
	scope.Machine = params.Machine

	if params.Cluster == nil {
		err = errors.New("cluster is required when creating a MachineScope")
		return nil, err
	}
	scope.Cluster = params.Cluster

	if params.IBMPowerVSMachine == nil {
		err = errors.New("ibmPowerVSMachine machine is required when creating a MachineScope")
		return nil, err
	}
	scope.IBMPowerVSMachine = params.IBMPowerVSMachine
	scope.IBMPowerVSCluster = params.IBMPowerVSCluster
	scope.IBMPowerVSImage = params.IBMPowerVSImage
	scope.ServiceEndpoint = params.ServiceEndpoint

	if params.Logger == (logr.Logger{}) {
		params.Logger = klog.Background()
	}
	if params.Logger.V(DEBUGLEVEL).Enabled() {
		core.SetLoggingLevel(core.LevelDebug)
	}

	// Create Resource Controller client.
	var serviceOption resourcecontroller.ServiceOptions
	// Fetch the resource controller endpoint.
	rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), params.ServiceEndpoint)
	if rcEndpoint != "" {
		serviceOption.URL = rcEndpoint
		params.Logger.V(3).Info("Overriding the default resource controller endpoint", "resourceControllerEndpoint", rcEndpoint)
	}

	rc, err := resourcecontroller.NewService(serviceOption)
	if err != nil {
		return nil, err
	}
	scope.ResourceClient = rc

	var serviceInstanceID, serviceInstanceName string
	if params.IBMPowerVSMachine.Spec.ServiceInstanceID != "" {
		serviceInstanceID = params.IBMPowerVSMachine.Spec.ServiceInstanceID
	} else if params.IBMPowerVSMachine.Spec.ServiceInstance != nil && params.IBMPowerVSMachine.Spec.ServiceInstance.ID != nil {
		serviceInstanceID = *params.IBMPowerVSMachine.Spec.ServiceInstance.ID
	} else {
		serviceInstanceName = fmt.Sprintf("%s-%s", params.IBMPowerVSCluster.GetName(), "serviceInstance")
		if params.IBMPowerVSCluster.Spec.ServiceInstance != nil && params.IBMPowerVSCluster.Spec.ServiceInstance.Name != nil {
			serviceInstanceName = *params.IBMPowerVSCluster.Spec.ServiceInstance.Name
		}
	}
	serviceInstance, err := rc.GetServiceInstance(serviceInstanceID, serviceInstanceName, params.IBMPowerVSCluster.Spec.Zone)
	if err != nil {
		params.Logger.Error(err, "failed to get PowerVS service instance details", "serviceInstanceName", serviceInstanceName, "serviceInstanceID", serviceInstanceID)
		return nil, err
	}
	if serviceInstance == nil {
		return nil, fmt.Errorf("PowerVS service instance %s is not yet created", serviceInstanceName)
	}
	if *serviceInstance.State != string(infrav1.ServiceInstanceStateActive) {
		return nil, fmt.Errorf("PowerVS service instance name: %s id: %s is not in active state", serviceInstanceName, serviceInstanceID)
	}
	serviceInstanceID = *serviceInstance.GUID

	region := endpoints.ConstructRegionFromZone(*serviceInstance.RegionID)
	scope.SetRegion(region)
	scope.SetZone(*serviceInstance.RegionID)

	serviceOptions := powervs.ServiceOptions{
		IBMPIOptions: &ibmpisession.IBMPIOptions{
			Debug: params.Logger.V(DEBUGLEVEL).Enabled(),
			Zone:  *serviceInstance.RegionID,
		},
		CloudInstanceID: serviceInstanceID,
	}

	// Fetch the service endpoint.
	if svcEndpoint := endpoints.FetchPVSEndpoint(region, params.ServiceEndpoint); svcEndpoint != "" {
		serviceOptions.IBMPIOptions.URL = svcEndpoint
		params.Logger.V(3).Info("Overriding the default PowerVS service endpoint", "serviceEndpoint", svcEndpoint)
	}

	c, err := powervs.NewService(serviceOptions)
	if err != nil {
		err = fmt.Errorf("failed to create PowerVS service")
		return nil, err
	}
	c.WithClients(serviceOptions)

	scope.IBMPowerVSClient = c
	scope.DHCPIPCacheStore = params.DHCPIPCacheStore

	var vpcRegion string
	if params.IBMPowerVSCluster.Spec.VPC == nil || params.IBMPowerVSCluster.Spec.VPC.Region == nil {
		vpcRegion, err = regionUtil.VPCRegionForPowerVSRegion(scope.GetRegion())
		if err != nil {
			return nil, fmt.Errorf("failed to create VPC client, error getting VPC region %v", err)
		}
	} else {
		vpcRegion = *params.IBMPowerVSCluster.Spec.VPC.Region
	}
	svcEndpoint := endpoints.FetchVPCEndpoint(vpcRegion, params.ServiceEndpoint)
	vpcClient, err := vpc.NewService(svcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create IBM VPC client: %w", err)
	}
	scope.IBMVPCClient = vpcClient
	return scope, nil
}

func (m *PowerVSMachineScope) ensureInstanceUnique(instanceName string) (*models.PVMInstanceReference, error) {
	instances, err := m.IBMPowerVSClient.GetAllInstance()
	if err != nil {
		return nil, err
	}
	for _, ins := range instances.PvmInstances {
		if *ins.ServerName == instanceName {
			return ins, nil
		}
	}
	return nil, nil
}

// CreateMachine creates a PowerVS machine.
func (m *PowerVSMachineScope) CreateMachine(ctx context.Context) (*models.PVMInstanceReference, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)

	machineSpec := m.IBMPowerVSMachine.Spec

	instanceReply, err := m.ensureInstanceUnique(m.IBMPowerVSMachine.Name)
	if err != nil {
		return nil, err
	} else if instanceReply != nil {
		log.Info("PowerVS instance already exists")
		return instanceReply, nil
	}

	// Check if create request has been already triggered.
	// If InstanceReadyCondition is Unknown then return and wait for it to get updated.
	for _, con := range m.IBMPowerVSMachine.Status.Conditions {
		if con.Type == infrav1.InstanceReadyCondition && con.Status == corev1.ConditionUnknown {
			return nil, nil
		}
		if con.Type == infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition && con.Status == corev1.ConditionUnknown {
			return nil, nil
		}
	}

	// TODO(karthik-k-n): Fix this
	userData, userDataErr := m.resolveUserData(ctx)
	if userDataErr != nil {
		return nil, fmt.Errorf("failed to resolve userdata %w", userDataErr)
	}

	memory := float64(machineSpec.MemoryGiB)

	var processors float64
	switch machineSpec.Processors.Type {
	case intstr.Int:
		processors = float64(machineSpec.Processors.IntVal)
	case intstr.String:
		processors, err = strconv.ParseFloat(machineSpec.Processors.StrVal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Processors(%s) to float64", machineSpec.Processors.StrVal)
		}
	}

	var imageID *string
	if m.IBMPowerVSImage != nil {
		imageID = &m.IBMPowerVSImage.Status.ImageID
	} else {
		imageID, err = getImageID(machineSpec.Image, m)
		if err != nil {
			record.Warnf(m.IBMPowerVSMachine, "FailedRetriveImage", "Failed image retrival - %v", err)
			return nil, fmt.Errorf("error getting image ID: %v", err)
		}
		log.V(3).Info("Retrieved image id", "imageID", *imageID)
	}
	network := machineSpec.Network
	if network.ID == nil && network.Name == nil && network.RegEx == nil {
		// if the network is nil, Fetch from cluster.
		if m.IBMPowerVSCluster.Status.Network != nil && m.IBMPowerVSCluster.Status.Network.ID != nil {
			network.ID = m.IBMPowerVSCluster.Status.Network.ID
		}
	}

	networkID, err := getNetworkID(network, m)
	if err != nil {
		record.Warnf(m.IBMPowerVSMachine, "FailedRetrieveNetwork", "Failed network retrieval - %v", err)
		return nil, fmt.Errorf("error getting network ID: %v", err)
	}
	log.V(3).Info("Retrieved network id", "networkID", *networkID)

	procType := strings.ToLower(string(machineSpec.ProcessorType))

	params := &p_cloud_p_vm_instances.PcloudPvminstancesPostParams{
		Body: &models.PVMInstanceCreate{
			ImageID: imageID,
			Networks: []*models.PVMInstanceAddNetwork{
				{
					NetworkID: networkID,
				},
			},
			ServerName: &m.IBMPowerVSMachine.Name,
			Memory:     &memory,
			Processors: &processors,
			ProcType:   &procType,
			SysType:    machineSpec.SystemType,
			UserData:   userData,
		},
	}
	if machineSpec.SSHKey != "" {
		params.Body.KeyPairName = machineSpec.SSHKey
	}
	log.V(3).Info("Creating PowerVS instance", "params", params)
	_, err = m.IBMPowerVSClient.CreateInstance(params.Body)
	if err != nil {
		record.Warnf(m.IBMPowerVSMachine, "FailedCreateInstance", "Failed instance creation - %v", err)
		return nil, err
	}
	record.Eventf(m.IBMPowerVSMachine, "SuccessfulCreateInstance", "Created Instance %q", m.IBMPowerVSMachine.Name)
	return nil, nil
}

func (m *PowerVSMachineScope) resolveUserData(ctx context.Context) (string, error) {
	userData, err := m.GetRawBootstrapData()
	if err != nil {
		return "", err
	}
	if m.UseIgnition() {
		data, err := m.ignitionUserData(ctx, userData)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}
	return base64.StdEncoding.EncodeToString(userData), err
}

func getIgnitionVersion(scope *PowerVSMachineScope) string {
	if scope.IBMPowerVSCluster.Spec.Ignition == nil {
		scope.IBMPowerVSCluster.Spec.Ignition = &infrav1.Ignition{}
	}
	if scope.IBMPowerVSCluster.Spec.Ignition.Version == "" {
		scope.IBMPowerVSCluster.Spec.Ignition.Version = infrav1.DefaultIgnitionVersion
	}
	return scope.IBMPowerVSCluster.Spec.Ignition.Version
}

func (m *PowerVSMachineScope) bootstrapDataKey() string {
	// Use machine name as object key.
	return path.Join(m.Role(), m.Name())
}

// Role returns the machine role from the labels.
func (m *PowerVSMachineScope) Role() string {
	if util.IsControlPlaneMachine(m.Machine) {
		return "control-plane"
	}
	return "node"
}

// Name returns the IBMPowerVSMachine name.
func (m *PowerVSMachineScope) Name() string {
	return m.IBMPowerVSMachine.Name
}

func (m *PowerVSMachineScope) createIgnitionData(ctx context.Context, data []byte) (string, error) {
	log := ctrl.LoggerFrom(ctx)
	if len(data) == 0 {
		return "", fmt.Errorf("user data is empty")
	}

	cosClient, err := m.createCOSClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create COS client %w", err)
	}
	key := m.bootstrapDataKey()
	log.V(3).Info("Bootstrap data key name", "key", key)

	bucket := m.bucketName()
	region := m.bucketRegion()
	if region == "" {
		return "", fmt.Errorf("failed to determine COS bucket region, both bucket region and VPC region not set")
	}

	if _, err := cosClient.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return "", fmt.Errorf("failed to push object to COS bucket %w", err)
	}

	objHost := fmt.Sprintf("%s.s3.%s.%s", bucket, region, cosURLDomain)

	cosServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.COS), m.ServiceEndpoint)
	if cosServiceEndpoint != "" {
		log.V(3).Info("Overriding the default COS endpoint in ignition URL", "cosEndpoint", cosServiceEndpoint)
		cosURL, _ := url.Parse(cosServiceEndpoint)
		if cosURL.Scheme != "" {
			objHost = fmt.Sprintf("%s.%s", bucket, cosURL.Host)
		} else {
			objHost = fmt.Sprintf("%s.%s", bucket, cosServiceEndpoint)
		}
	}
	objectURL := &url.URL{
		Scheme: "https",
		Host:   objHost,
		Path:   key,
	}
	log.V(3).Info("Generated Ignition URL", "objectURL", objectURL.String())

	return objectURL.String(), nil
}

func (m *PowerVSMachineScope) ignitionUserData(ctx context.Context, userData []byte) ([]byte, error) {
	objectURL, err := m.createIgnitionData(ctx, userData)
	if err != nil {
		return nil, fmt.Errorf("failed to create user data object %w", err)
	}

	auth, err := authenticator.GetIAMAuthenticator()
	if err != nil {
		return nil, err
	}

	iamtoken, err := auth.GetToken()
	if err != nil {
		return nil, err
	}
	if iamtoken == "" {
		return nil, fmt.Errorf("IAM token is empty")
	}
	token := "Bearer " + iamtoken

	ignVersion := getIgnitionVersion(m)
	semver, err := semver.ParseTolerant(ignVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ignition version %q: %w", ignVersion, err)
	}

	switch semver.Major {
	case 2:
		ignData := &ignV2Types.Config{
			Ignition: ignV2Types.Ignition{
				Version: semver.String(),
				Config: ignV2Types.IgnitionConfig{
					Replace: &ignV2Types.ConfigReference{
						Source: objectURL,
						HTTPHeaders: ignV2Types.HTTPHeaders{
							{
								Name:  "Authorization",
								Value: token,
							},
						},
					},
				},
			},
		}
		return json.Marshal(ignData)
	case 3:
		ignData := &ignV3Types.Config{
			Ignition: ignV3Types.Ignition{
				Version: semver.String(),
				Config: ignV3Types.IgnitionConfig{
					Replace: ignV3Types.Resource{
						Source: aws.String(objectURL),
						HTTPHeaders: ignV3Types.HTTPHeaders{
							{
								Name:  "Authorization",
								Value: aws.String(token),
							},
						},
					},
				},
			},
		}
		return json.Marshal(ignData)
	default:
		return nil, fmt.Errorf("unsupported ignition version %q", ignVersion)
	}
}

// UseIgnition returns true if Ignition is set in IBMPowerVSCluster.
func (m *PowerVSMachineScope) UseIgnition() bool {
	return m.IBMPowerVSCluster.Spec.Ignition != nil
}

// DeleteMachine deletes the power vs machine associated with machine instance id and service instance id.
func (m *PowerVSMachineScope) DeleteMachine() error {
	if err := m.IBMPowerVSClient.DeleteInstance(m.IBMPowerVSMachine.Status.InstanceID); err != nil {
		record.Warnf(m.IBMPowerVSMachine, "FailedDeleteInstance", "Failed instance deletion - %v", err)
		return err
	}
	record.Eventf(m.IBMPowerVSMachine, "SuccessfulDeleteInstance", "Deleted Instance %q", m.IBMPowerVSMachine.Name)
	return nil
}

// DeleteMachineIgnition deletes the ignition associated with machine.
func (m *PowerVSMachineScope) DeleteMachineIgnition(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	_, err := m.GetRawBootstrapData()
	if err != nil {
		return err
	}
	if !m.UseIgnition() {
		log.V(3).Info("Machine is not using user data of type ignition")
		return nil
	}
	cosClient, err := m.createCOSClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create COS client %w", err)
	}

	bucket := m.bucketName()
	objs, _ := cosClient.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	for _, j := range objs.Contents {
		if strings.Contains(*j.Key, m.Name()) {
			if _, err := cosClient.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    j.Key,
			}); err != nil {
				record.Warnf(m.IBMPowerVSMachine, "FailedDeleteMachineIgnition", "Failed machine ignition deletion - %v", err)
				return fmt.Errorf("failed to delete COS object %w", err)
			}
		}
	}
	record.Eventf(m.IBMPowerVSMachine, "SuccessfulDeleteMachineIgnition", "Deleted machine ignition %q", m.IBMPowerVSMachine.Name)
	return nil
}

// createCOSClient creates a new cosClient from the supplied parameters.
func (m *PowerVSMachineScope) createCOSClient(ctx context.Context) (cos.Cos, error) {
	log := ctrl.LoggerFrom(ctx)
	var cosInstanceName string
	if m.IBMPowerVSCluster.Spec.CosInstance == nil || m.IBMPowerVSCluster.Spec.CosInstance.Name == "" {
		cosInstanceName = fmt.Sprintf("%s-%s", m.IBMPowerVSCluster.GetName(), "cosinstance")
	} else {
		cosInstanceName = m.IBMPowerVSCluster.Spec.CosInstance.Name
	}

	serviceInstance, err := m.ResourceClient.GetInstanceByName(cosInstanceName, resourcecontroller.CosResourceID, resourcecontroller.CosResourcePlanID)
	if err != nil {
		log.Error(err, "failed to get COS service instance", "name", cosInstanceName)
		return nil, err
	}
	if serviceInstance == nil {
		log.V(3).Info("COS service instance is nil")
		return nil, errors.New("COS service instance is nil")
	}
	if *serviceInstance.State != string(infrav1.ServiceInstanceStateActive) {
		return nil, fmt.Errorf("COS service instance is not in active state, current state: %s", *serviceInstance.State)
	}

	props, err := authenticator.GetProperties()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service properties: %w", err)
	}
	apiKey := props["APIKEY"]
	if apiKey == "" {
		fmt.Printf("IBM Cloud API key is not provided, set %s environmental variable", "IBMCLOUD_API_KEY")
	}

	region := m.bucketRegion()
	if region == "" {
		return nil, fmt.Errorf("failed to determine COS bucket region, both bucket region and VPC region not set")
	}

	serviceEndpoint := fmt.Sprintf("s3.%s.%s", region, cosURLDomain)
	// Fetch the COS service endpoint.
	cosServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.COS), m.ServiceEndpoint)
	if cosServiceEndpoint != "" {
		log.V(3).Info("Overriding the default COS endpoint", "cosEndpoint", cosServiceEndpoint)
		serviceEndpoint = cosServiceEndpoint
	}

	cosOptions := cos.ServiceOptions{
		Options: &cosSession.Options{
			Config: aws.Config{
				Endpoint: &serviceEndpoint,
				Region:   &region,
			},
		},
	}

	cosClient, err := cos.NewService(cosOptions, apiKey, *serviceInstance.GUID)
	if err != nil {
		return nil, fmt.Errorf("failed to create COS client: %w", err)
	}
	return cosClient, nil
}

// GetRawBootstrapData returns the bootstrap data if present.
func (m *PowerVSMachineScope) GetRawBootstrapData() ([]byte, error) {
	if m.Machine == nil || m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, errors.New("failed to retrieve bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Machine.Namespace, Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.Client.Get(context.TODO(), key, secret); err != nil {
		return nil, fmt.Errorf("failed to retrieve bootstrap data secret: %v", err)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, errors.New("failed to retrieve bootstrap data: secret value key is missing")
	}

	return value, nil
}

func getImageID(image *infrav1.IBMPowerVSResourceReference, m *PowerVSMachineScope) (*string, error) {
	if image.ID != nil {
		return image.ID, nil
	} else if image.Name != nil {
		images, err := m.GetImages()
		if err != nil {
			return nil, err
		}
		for _, img := range images.Images {
			if *image.Name == *img.Name {
				return img.ImageID, nil
			}
		}
	} else {
		return nil, fmt.Errorf("both image ID and image Name can't be nil")
	}
	return nil, fmt.Errorf("failed to find an image ID")
}

// GetImages will get list of images for the powervs service instance.
func (m *PowerVSMachineScope) GetImages() (*models.Images, error) {
	return m.IBMPowerVSClient.GetAllImage()
}

func getNetworkID(network infrav1.IBMPowerVSResourceReference, m *PowerVSMachineScope) (*string, error) {
	if network.ID != nil {
		return network.ID, nil
	} else if network.Name != nil {
		networks, err := m.GetNetworks()
		if err != nil {
			return nil, err
		}
		for _, nw := range networks.Networks {
			if *network.Name == *nw.Name {
				return nw.NetworkID, nil
			}
		}
		return nil, fmt.Errorf("failed to find a network ID with name %s", *network.Name)
	} else if network.RegEx != nil {
		networks, err := m.GetNetworks()
		if err != nil {
			return nil, err
		}
		re, err := regexp.Compile(*network.RegEx)
		if err != nil {
			return nil, err
		}
		// In case of multiple network names matches the provided regular expression the first matched network will be selected.
		for _, nw := range networks.Networks {
			if match := re.Match([]byte(*nw.Name)); match {
				return nw.NetworkID, nil
			}
		}
		return nil, fmt.Errorf("failed to find a network ID with RegEx %s", *network.RegEx)
	}
	return nil, fmt.Errorf("ID, Name and RegEx can't be nil")
}

// GetNetworks will get list of networks for the powervs service instance.
func (m *PowerVSMachineScope) GetNetworks() (*models.Networks, error) {
	return m.IBMPowerVSClient.GetAllNetwork()
}

// SetReady will set the status as ready for the machine.
func (m *PowerVSMachineScope) SetReady() {
	m.IBMPowerVSMachine.Status.Ready = true
}

// SetNotReady will set status as not ready for the machine.
func (m *PowerVSMachineScope) SetNotReady() {
	m.IBMPowerVSMachine.Status.Ready = false
}

// SetFailureReason will set status FailureReason for the machine.
func (m *PowerVSMachineScope) SetFailureReason(reason string) {
	//nolint:staticcheck
	m.IBMPowerVSMachine.Status.FailureReason = &reason
}

// SetFailureMessage will set status FailureMessage for the machine.
func (m *PowerVSMachineScope) SetFailureMessage(message string) {
	//nolint:staticcheck
	m.IBMPowerVSMachine.Status.FailureMessage = &message
}

// IsReady will return the status for the machine.
func (m *PowerVSMachineScope) IsReady() bool {
	return m.IBMPowerVSMachine.Status.Ready
}

// SetInstanceID will set the instance id for the machine.
func (m *PowerVSMachineScope) SetInstanceID(id *string) {
	if id != nil {
		m.IBMPowerVSMachine.Status.InstanceID = *id
	}
}

// GetInstanceID will get the instance id for the machine.
func (m *PowerVSMachineScope) GetInstanceID() string {
	return m.IBMPowerVSMachine.Status.InstanceID
}

// SetHealth will set the health status for the machine.
func (m *PowerVSMachineScope) SetHealth(health *models.PVMInstanceHealth) {
	if health != nil {
		m.IBMPowerVSMachine.Status.Health = health.Status
	}
}

// SetAddresses will set the addresses for the machine.
func (m *PowerVSMachineScope) SetAddresses(ctx context.Context, instance *models.PVMInstance) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)
	var addresses []corev1.NodeAddress
	// Setting the name of the vm to the InternalDNS and Hostname as the vm uses that as hostname.
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: *instance.ServerName,
	})
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeHostName,
		Address: *instance.ServerName,
	})
	for _, network := range instance.Networks {
		if strings.TrimSpace(network.IPAddress) != "" {
			addresses = append(addresses, corev1.NodeAddress{
				Type:    corev1.NodeInternalIP,
				Address: strings.TrimSpace(network.IPAddress),
			})
		}
		if strings.TrimSpace(network.ExternalIP) != "" {
			addresses = append(addresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: strings.TrimSpace(network.ExternalIP),
			})
		}
	}
	m.IBMPowerVSMachine.Status.Addresses = addresses
	if len(addresses) > 2 {
		// If the address length is more than 2 means either NodeInternalIP or NodeExternalIP is updated so return
		return
	}
	// In this case there is no IP found under instance.Networks, So try to fetch the IP from cache or DHCP server

	// Look for DHCP IP from the cache
	obj, exists, err := m.DHCPIPCacheStore.GetByKey(*instance.ServerName)
	if err != nil {
		log.Error(err, "failed to fetch the DHCP IP address from cache store")
	} else if exists {
		log.V(3).Info("Found IP for machine from DHCP cache", "IP", obj.(powervs.VMip).IP)
		addresses = append(addresses, corev1.NodeAddress{
			Type:    corev1.NodeInternalIP,
			Address: obj.(powervs.VMip).IP,
		})
		m.IBMPowerVSMachine.Status.Addresses = addresses
		return
	}
	// Fetch the VM network ID
	network := m.IBMPowerVSMachine.Spec.Network
	if network.ID == nil && network.Name == nil && network.RegEx == nil {
		// if the network is nil, Fetch from cluster.
		if m.IBMPowerVSCluster.Status.Network != nil && m.IBMPowerVSCluster.Status.Network.ID != nil {
			network.ID = m.IBMPowerVSCluster.Status.Network.ID
		}
	}
	networkID, err := getNetworkID(network, m)
	if err != nil {
		log.Error(err, "failed to fetch network id from network resource")
		return
	}
	log.V(3).Info("Retrieved network id", "networkID", *networkID)
	// Fetch the details of the network attached to the VM
	var pvmNetwork *models.PVMInstanceNetwork
	for _, network := range instance.Networks {
		if network.NetworkID == *networkID {
			pvmNetwork = network
			log.V(3).Info("Found network attached to machine", "networkID", network.NetworkID)
		}
	}
	if pvmNetwork == nil {
		log.V(3).Info("Failed to get network attached to machine", "networkID", *networkID)
		return
	}
	// Get all the DHCP servers
	dhcpServer, err := m.IBMPowerVSClient.GetAllDHCPServers()
	if err != nil {
		log.Error(err, "failed to get DHCP server")
		return
	}
	// Get the Details of DHCP server associated with the network
	var dhcpServerDetails *models.DHCPServerDetail
	for _, server := range dhcpServer {
		if server.Network == nil || server.Network.ID == nil {
			log.V(3).Info("Skipping the DHCP server as its network details is nil", "dhcpServerID", *server.ID)
			continue
		}
		if *server.Network.ID == *networkID {
			log.V(3).Info("Found DHCP server for network", "dhcpServerID", *server.ID, "networkID", *networkID)
			dhcpServerDetails, err = m.IBMPowerVSClient.GetDHCPServer(*server.ID)
			if err != nil {
				log.Error(err, "failed to get DHCP server details", "dhcpServerID", *server.ID)
				return
			}
			break
		}
	}
	if dhcpServerDetails == nil {
		errStr := fmt.Errorf("DHCP server details is nil")
		log.Error(errStr, "DHCP server associated with network is nil", "networkID", *networkID)
		return
	}

	// Fetch the VM IP using VM's mac from DHCP server lease
	var internalIP *string
	for _, lease := range dhcpServerDetails.Leases {
		if *lease.InstanceMacAddress == pvmNetwork.MacAddress {
			log.V(3).Info("Found internal IP for machine from DHCP lease", "IP", *lease.InstanceIP)
			internalIP = lease.InstanceIP
			break
		}
	}
	if internalIP == nil {
		errStr := errors.New("internal IP is nil")
		log.Error(errStr, "failed to get internal IP, DHCP lease not found for machine with MAC in DHCP network",
			"mac", pvmNetwork.MacAddress, "dhcpServerID", *dhcpServerDetails.ID)
		return
	}
	log.V(3).Info("Found internal IP for VM from DHCP lease", "IP", *internalIP)
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalIP,
		Address: *internalIP,
	})
	// Update the cache with the ip and VM name
	err = m.DHCPIPCacheStore.Add(powervs.VMip{
		Name: *instance.ServerName,
		IP:   *internalIP,
	})
	if err != nil {
		log.Error(err, "failed to update the DHCP cache store with the IP", "IP", *internalIP)
	}
	m.IBMPowerVSMachine.Status.Addresses = addresses
}

// SetInstanceState will set the state for the machine.
func (m *PowerVSMachineScope) SetInstanceState(status *string) {
	m.IBMPowerVSMachine.Status.InstanceState = infrav1.PowerVSInstanceState(*status)
}

// GetInstanceState will get the state for the machine.
func (m *PowerVSMachineScope) GetInstanceState() infrav1.PowerVSInstanceState {
	return m.IBMPowerVSMachine.Status.InstanceState
}

// SetRegion will set the region for the machine.
func (m *PowerVSMachineScope) SetRegion(region string) {
	m.IBMPowerVSMachine.Status.Region = &region
}

// GetRegion will get the region for the machine.
func (m *PowerVSMachineScope) GetRegion() string {
	if m.IBMPowerVSMachine.Status.Region == nil {
		return ""
	}
	return *m.IBMPowerVSMachine.Status.Region
}

// SetZone will set the zone for the machine.
func (m *PowerVSMachineScope) SetZone(zone string) {
	m.IBMPowerVSMachine.Status.Zone = &zone
}

// GetZone will get the zone for the machine.
func (m *PowerVSMachineScope) GetZone() string {
	if m.IBMPowerVSMachine.Status.Zone == nil {
		return ""
	}
	return *m.IBMPowerVSMachine.Status.Zone
}

// GetServiceInstanceID returns the service instance id.
func (m *PowerVSMachineScope) GetServiceInstanceID() (string, error) {
	if m.IBMPowerVSCluster.Status.ServiceInstance != nil && m.IBMPowerVSCluster.Status.ServiceInstance.ID != nil {
		return *m.IBMPowerVSCluster.Status.ServiceInstance.ID, nil
	}
	if m.IBMPowerVSCluster.Spec.ServiceInstanceID != "" {
		return m.IBMPowerVSCluster.Spec.ServiceInstanceID, nil
	}
	if m.IBMPowerVSCluster.Spec.ServiceInstance != nil && m.IBMPowerVSCluster.Spec.ServiceInstance.ID != nil {
		return *m.IBMPowerVSCluster.Spec.ServiceInstance.ID, nil
	}
	// If we are not able to find service instance id, derive it from name if defined.
	if m.IBMPowerVSCluster.Spec.ServiceInstance != nil && m.IBMPowerVSCluster.Spec.ServiceInstance.Name == nil {
		return "", fmt.Errorf("failed to find service instance id as both name and id are not set")
	}
	serviceInstance, err := m.ResourceClient.GetServiceInstance("", *m.IBMPowerVSCluster.Spec.ServiceInstance.Name, ptr.To(m.GetZone()))
	if err != nil {
		return "", err
	}
	// It's safe to directly dereference GUID as its already done in NewPowerVSMachineScope
	return *serviceInstance.GUID, nil
}

// SetProviderID will set the provider id for the machine.
func (m *PowerVSMachineScope) SetProviderID(instanceID string) error {
	if options.ProviderIDFormatType(options.ProviderIDFormat) != options.ProviderIDFormatV2 {
		return fmt.Errorf("invalid value for ProviderIDFormat")
	}

	serviceInstanceID, err := m.GetServiceInstanceID()
	if err != nil {
		return err
	}
	m.IBMPowerVSMachine.Spec.ProviderID = ptr.To(fmt.Sprintf("ibmpowervs://%s/%s/%s/%s", m.GetRegion(), m.GetZone(), serviceInstanceID, instanceID))
	return nil
}

// GetMachineInternalIP returns the machine's internal IP.
func (m *PowerVSMachineScope) GetMachineInternalIP() string {
	for _, address := range m.IBMPowerVSMachine.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return address.Address
		}
	}
	return ""
}

// CreateVPCLoadBalancerPoolMember creates a member in load balancer pool.
func (m *PowerVSMachineScope) CreateVPCLoadBalancerPoolMember(ctx context.Context) (*vpcv1.LoadBalancerPoolMember, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)
	loadBalancers := make([]infrav1.VPCLoadBalancerSpec, 0)
	if len(m.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		loadBalancer := infrav1.VPCLoadBalancerSpec{
			Name:   fmt.Sprintf("%s-loadbalancer", m.IBMPowerVSCluster.Name),
			Public: ptr.To(true),
		}
		loadBalancers = append(loadBalancers, loadBalancer)
	}
	for index, loadBalancer := range m.IBMPowerVSCluster.Spec.LoadBalancers {
		if loadBalancer.Name == "" {
			loadBalancer.Name = fmt.Sprintf("%s-loadbalancer-%d", m.IBMPowerVSCluster.Name, index)
		}
		loadBalancers = append(loadBalancers, loadBalancer)
	}

	for _, lb := range loadBalancers {
		var lbID *string
		if m.IBMPowerVSCluster.Status.LoadBalancers == nil {
			return nil, fmt.Errorf("failed to find VPC load balancer ID")
		}
		if val, ok := m.IBMPowerVSCluster.Status.LoadBalancers[lb.Name]; ok {
			lbID = val.ID
		} else {
			return nil, fmt.Errorf("failed to find VPC load balancer ID")
		}
		loadBalancer, _, err := m.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
			ID: lbID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find VPC load balancer details: %w", err)
		}
		if *loadBalancer.ProvisioningStatus != string(infrav1.VPCLoadBalancerStateActive) {
			return nil, fmt.Errorf("VPC load balancer is not in active state, current state %s", *loadBalancer.ProvisioningStatus)
		}
		if len(loadBalancer.Pools) == 0 {
			return nil, fmt.Errorf("no pools exist for the VPC load balancer %s", lb.Name)
		}

		internalIP := m.GetMachineInternalIP()

		// lbAdditionalListeners is a mapping of additionalListener's port-protocol to the additionalListener as defined in the specification
		// It will be used later to get the default pool associated with the listener
		lbAdditionalListeners := map[string]infrav1.AdditionalListenerSpec{}
		for _, additionalListener := range lb.AdditionalListeners {
			if additionalListener.Protocol == nil {
				additionalListener.Protocol = &infrav1.VPCLoadBalancerListenerProtocolTCP
			}
			lbAdditionalListeners[fmt.Sprintf("%d-%s", additionalListener.Port, *additionalListener.Protocol)] = additionalListener
		}

		// loadBalancerListeners is a mapping of the loadBalancer listener's defaultPoolName to the additionalListener
		// as the default pool name might be empty in spec and should be fetched from the cloud's listener
		loadBalancerListeners := map[string]infrav1.AdditionalListenerSpec{}
		for _, listener := range loadBalancer.Listeners {
			listenerOptions := &vpcv1.GetLoadBalancerListenerOptions{}
			listenerOptions.SetLoadBalancerID(*loadBalancer.ID)
			listenerOptions.SetID(*listener.ID)
			loadBalancerListener, _, err := m.IBMVPCClient.GetLoadBalancerListener(listenerOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list %s load balancer listener: %w", *listener.ID, err)
			}
			if additionalListener, ok := lbAdditionalListeners[fmt.Sprintf("%d-%s", *loadBalancerListener.Port, *loadBalancerListener.Protocol)]; ok {
				if loadBalancerListener.DefaultPool != nil {
					loadBalancerListeners[*loadBalancerListener.DefaultPool.Name] = additionalListener
				}
				// loadBalancerListeners map is created only with the listeners provided in the spec,
				// and targetPort is populated only if there is an entry in the map.
				// Inorder for the default pool 6443 to be added to all control plane machines, creating an entry in the map for the same.
			} else if loadBalancerListener.Port != nil && *loadBalancerListener.Port == int64(6443) {
				protocol := infrav1.VPCLoadBalancerListenerProtocol(*loadBalancerListener.Protocol)
				listener := infrav1.AdditionalListenerSpec{
					Port:     *loadBalancerListener.Port,
					Protocol: &protocol,
				}
				if loadBalancerListener.DefaultPool != nil {
					loadBalancerListeners[*loadBalancerListener.DefaultPool.Name] = listener
				} else {
					log.V(3).Error(fmt.Errorf("unable to get the default pool details"), "default pool is nil", "port", loadBalancerListener.Port)
				}
			}
		}
		// Update each LoadBalancer pool
		// For each pool, get the additionalListener associated with the pool from the loadBalancerListeners map.
		for _, pool := range loadBalancer.Pools {
			log.V(3).Info("Updating LoadBalancer pool member", "pool", *pool.Name, "loadBalancerName", *loadBalancer.Name, "IP", internalIP)
			listOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{}
			listOptions.SetLoadBalancerID(*loadBalancer.ID)
			listOptions.SetPoolID(*pool.ID)
			listLoadBalancerPoolMembers, _, err := m.IBMVPCClient.ListLoadBalancerPoolMembers(listOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list %s VPC load balancer pool: %w", *pool.Name, err)
			}
			var targetPort int64
			var alreadyRegistered bool

			if loadBalancerListener, ok := loadBalancerListeners[*pool.Name]; ok {
				targetPort = loadBalancerListener.Port
				log.V(3).Info("Checking if machine label matches with the label selector in listener", "machineLabel", m.IBMPowerVSMachine.Labels, "labelSelector", loadBalancerListener.Selector)
				selector, err := metav1.LabelSelectorAsSelector(&loadBalancerListener.Selector)
				if err != nil {
					log.V(5).Error(err, "Skipping listener addition, failed to get label selector from spec selector")
					continue
				}

				if selector.Empty() && !util.IsControlPlaneMachine(m.Machine) {
					log.V(3).Info("Skipping listener addition as the selector is empty and not a control plane machine")
					continue
				}
				// Skip adding the listener if the selector does not match
				if !selector.Empty() && !selector.Matches(labels.Set(m.IBMPowerVSMachine.Labels)) {
					log.V(3).Info("Skip adding listener, machine label doesn't match with the listener label selector", "pool", *pool.Name, "IP", internalIP)
					continue
				}
			}

			for _, member := range listLoadBalancerPoolMembers.Members {
				if target, ok := member.Target.(*vpcv1.LoadBalancerPoolMemberTarget); ok {
					if *target.Address == internalIP {
						alreadyRegistered = true
						log.V(3).Info("Target IP already configured for pool", "IP", internalIP, "poolName", *pool.Name)
					}
				}
			}

			if alreadyRegistered {
				log.V(3).Info("PoolMember already exist", "poolName", *pool.Name, "IP", internalIP, "targetPort", targetPort)
				continue
			}

			// make sure that LoadBalancer is in active state
			loadBalancer, _, err := m.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: loadBalancer.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to fetch VPC load balancer details with ID: %s error: %v", *lbID, err)
			}
			if *loadBalancer.ProvisioningStatus != string(infrav1.VPCLoadBalancerStateActive) {
				log.V(3).Info("Unable to update pool for VPC load balancer as it is not in active state", "loadBalancerName", *loadBalancer.Name, "loadBalancerState", *loadBalancer.ProvisioningStatus)
				return nil, fmt.Errorf("VPC load balancer %s not in active state to update pool member", *loadBalancer.Name)
			}

			options := &vpcv1.CreateLoadBalancerPoolMemberOptions{}
			options.SetPort(targetPort)
			options.SetLoadBalancerID(*loadBalancer.ID)
			options.SetPoolID(*pool.ID)
			options.SetTarget(&vpcv1.LoadBalancerPoolMemberTargetPrototype{
				Address: &internalIP,
			})
			log.V(3).Info("Creating VPC load balancer pool member", "options", options)
			loadBalancerPoolMember, _, err := m.IBMVPCClient.CreateLoadBalancerPoolMember(options)
			if err != nil {
				return nil, fmt.Errorf("failed to create VPC load balancer %s pool member %w", *loadBalancer.Name, err)
			}
			log.Info("Created VPC load balancer pool member", "loadBalancerID", *loadBalancerPoolMember.ID)
			return loadBalancerPoolMember, nil
		}
	}
	return nil, nil
}

// APIServerPort returns the APIServerPort.
func (m *PowerVSMachineScope) APIServerPort() int32 {
	if m.Cluster.Spec.ClusterNetwork.APIServerPort > 0 {
		return m.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return infrav1.DefaultAPIServerPort
}

// TODO: reuse getServiceName function instead.
func (m *PowerVSMachineScope) bucketName() string {
	if m.IBMPowerVSCluster.Spec.CosInstance != nil && m.IBMPowerVSCluster.Spec.CosInstance.BucketName != "" {
		return m.IBMPowerVSCluster.Spec.CosInstance.BucketName
	}
	return fmt.Sprintf("%s-%s", m.IBMPowerVSCluster.GetName(), "cosbucket")
}

// TODO: duplicate function, optimize it.
func (m *PowerVSMachineScope) bucketRegion() string {
	if m.IBMPowerVSCluster.Spec.CosInstance != nil && m.IBMPowerVSCluster.Spec.CosInstance.BucketRegion != "" {
		return m.IBMPowerVSCluster.Spec.CosInstance.BucketRegion
	}
	// if the bucket region is not set, use vpc region
	vpcDetails := m.IBMPowerVSCluster.Spec.VPC
	if vpcDetails != nil && vpcDetails.Region != nil {
		return *vpcDetails.Region
	}
	return ""
}
