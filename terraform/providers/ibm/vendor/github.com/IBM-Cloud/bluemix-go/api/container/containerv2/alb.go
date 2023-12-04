package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type alb struct {
	client *client.Client
}

type AlbConfig struct {
	AlbBuild             string `json:"albBuild"`
	AlbID                string `json:"albID"`
	AlbType              string `json:"albType"`
	AuthBuild            string `json:"authBuild"`
	Cluster              string `json:"cluster"`
	CreatedDate          string `json:"createdDate"`
	DisableDeployment    bool   `json:"disableDeployment"`
	Enable               bool   `json:"enable"`
	LoadBalancerHostname string `json:"loadBalancerHostname"`
	Name                 string `json:"name"`
	NumOfInstances       string `json:"numOfInstances"`
	Resize               bool   `json:"resize"`
	State                string `json:"state"`
	Status               string `json:"status"`
	ZoneAlb              string `json:"zone"`
}

// ALBClusterHealthCheckConfig configuration for ALB in-cluster health check
type ALBClusterHealthCheckConfig struct {
	Cluster string `json:"cluster"`
	Enable  bool   `json:"enable"`
}

type AlbCreateResp struct {
	Cluster string `json:"cluster"`
	Alb     string `json:"alb"`
}

type AlbCreateReq struct {
	Cluster         string `json:"cluster"`
	EnableByDefault bool   `json:"enableByDefault"`
	Type            string `json:"type"`
	ZoneAlb         string `json:"zone"`
	IngressImage    string `json:"ingressImage"`
}

type AlbImageVersions struct {
	DefaultK8sVersion    string   `json:"defaultK8sVersion"`
	SupportedK8sVersions []string `json:"supportedK8sVersions"`
}

// ALBLBConfig ingress ALB load balancer configuration
type ALBLBConfig struct {
	Cluster       string                    `json:"cluster" description:"The ID or name of the cluster" binding:"required"`
	Type          string                    `json:"type,omitempty" description:"type of load balancers to configure ('public' or 'private')"`
	ProxyProtocol *ALBLBProxyProtocolConfig `json:"proxyProtocol,omitempty" description:"PROXY protocol related configurations"`
}

// ALBLBProxyProtocolConfig ingress ALB load balancer Proxy Protocol configuration
type ALBLBProxyProtocolConfig struct {
	Enable        bool     `json:"enable" description:"PROXY protocol state"`
	CIDR          []string `json:"cidr,omitempty" description:"trusted address ranges"`
	HeaderTimeout int      `json:"headerTimeout,omitempty" description:"timeout value for receiving PROXY protocol headers"`
}

// AlbUpdateResp
type AlbUpdateResp struct {
	ClusterID string `json:"clusterID"`
}

// AutoscaleDetails
type AutoscaleDetails struct {
	Config *AutoscaleConfig `json:"config,omitempty" description:"Autoscaling configuration"`
}

// AutoscaleConfig
type AutoscaleConfig struct {
	MinReplicas           int    `json:"minReplicas" description:"Minimum number of replicas"`
	MaxReplicas           int    `json:"maxReplicas" description:"Maximum number of replicas"`
	CustomMetrics         string `json:"customMetrics,omitempty" description:"An array of MetricSpec (autoscaling.k8s.io/v2) encoded as JSON"`
	CPUAverageUtilization int    `json:"cpuAverageUtilization,omitempty" description:"CPU Average Utilization"`
}

type ClusterALB struct {
	ID                      string      `json:"id"`
	Region                  string      `json:"region"`
	DataCenter              string      `json:"dataCenter"`
	IsPaid                  bool        `json:"isPaid"`
	PublicIngressHostname   string      `json:"publicIngressHostname"`
	PublicIngressSecretName string      `json:"publicIngressSecretName"`
	ALBs                    []AlbConfig `json:"alb"`
}

// IgnoredIngressStatusErrors
type IgnoredIngressStatusErrors struct {
	Cluster       string   `json:"cluster" description:"the ID or name of the cluster"`
	IgnoredErrors []string `json:"ignoredErrors" description:"list of error codes that the user wants to ignore"`
}

// IngressStatusState
type IngressStatusState struct {
	Cluster string `json:"cluster" description:"the ID or name of the cluster"`
	Enable  bool   `json:"enable" description:"true or false to enable or disable ingress status"`
}

// IngressStatus struct for the top level ingress status, for a cluster
type IngressStatus struct {
	Cluster                string `json:"cluster"`
	Status                 string `json:"status"`
	NonTranslatedStatus    string `json:"nonTranslatedStatus"`
	Message                string `json:"message"`
	StatusList             []IngressComponentStatus
	GeneralComponentStatus []V2IngressComponentStatus `json:"generalComponentStatus,omitempty"`
	ALBStatus              []V2IngressComponentStatus `json:"albStatus,omitempty"`
	RouterStatus           []V2IngressComponentStatus `json:"routerStatus,omitempty"`
	SubdomainStatus        []V2IngressComponentStatus `json:"subdomainStatus,omitempty"`
	SecretStatus           []V2IngressComponentStatus `json:"secretStatus,omitempty"`
	IgnoredErrors          []string                   `json:"ignoredErrors" description:"list of error codes that the user wants to ignore"`
}

// IngressComponentStatus status of individual ingress component
type IngressComponentStatus struct {
	Component string `json:"component"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

// UpdateALBReq is the body of the v2 Update ALB API endpoint
type UpdateALBReq struct {
	ClusterID string   `json:"cluster" description:"The ID of the cluster on which the update ALB action shall be performed"`
	ALBBuild  string   `json:"albBuild" description:"The version of the build to which the ALB should be updated"`
	ALBList   []string `json:"albList" description:"The list of ALBs that should be updated to the requested albBuild"`
}

// V2IngressComponentStatus status of individual ingress component
type V2IngressComponentStatus struct {
	Component string   `json:"component,omitempty"`
	Status    []string `json:"status,omitempty"`
}

// Clusters interface
type Alb interface {
	AddIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error
	CreateAlb(albCreateReq AlbCreateReq, target ClusterTargetHeader) (AlbCreateResp, error)
	DisableAlb(disableAlbReq AlbConfig, target ClusterTargetHeader) error
	EnableAlb(enableAlbReq AlbConfig, target ClusterTargetHeader) error
	GetALBAutoscaleConfiguration(clusterNameOrID, albID string, target ClusterTargetHeader) (AutoscaleDetails, error)
	GetAlb(albid string, target ClusterTargetHeader) (AlbConfig, error)
	GetAlbClusterHealthCheckConfig(clusterNameOrID string, target ClusterTargetHeader) (ALBClusterHealthCheckConfig, error)
	GetIgnoredIngressStatusErrors(clusterNameOrID string, target ClusterTargetHeader) (IgnoredIngressStatusErrors, error)
	GetIngressLoadBalancerConfig(clusterNameOrID, lbType string, target ClusterTargetHeader) (ALBLBConfig, error)
	GetIngressStatus(clusterNameOrID string, target ClusterTargetHeader) (IngressStatus, error)
	ListAlbImages(target ClusterTargetHeader) (AlbImageVersions, error)
	ListClusterAlbs(clusterNameOrID string, target ClusterTargetHeader) ([]AlbConfig, error)
	RemoveALBAutoscaleConfiguration(clusterNameOrID, albID string, target ClusterTargetHeader) error
	RemoveIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error
	SetALBAutoscaleConfiguration(clusterNameOrID, albID string, autoscaleDetails AutoscaleDetails, target ClusterTargetHeader) error
	SetAlbClusterHealthCheckConfig(albHealthCheckReq ALBClusterHealthCheckConfig, target ClusterTargetHeader) error
	SetIngressStatusState(ingressStatusStateReq IngressStatusState, target ClusterTargetHeader) error
	UpdateAlb(updateAlbReq UpdateALBReq, target ClusterTargetHeader) error
	UpdateIngressLoadBalancerConfig(lbConfig ALBLBConfig, target ClusterTargetHeader) error
}

func newAlbAPI(c *client.Client) Alb {
	return &alb{
		client: c,
	}
}

// CreateAlb create an ALB in a specified zone and cluster
func (r *alb) CreateAlb(albCreateReq AlbCreateReq, target ClusterTargetHeader) (AlbCreateResp, error) {
	var successV AlbCreateResp
	_, err := r.client.Post("/v2/alb/vpc/createAlb", albCreateReq, &successV, target.ToMap())
	return successV, err
}

// DisableAlb disable an ALB in your cluster
func (r *alb) DisableAlb(disableAlbReq AlbConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/vpc/disableAlb", disableAlbReq, nil, target.ToMap())
	return err
}

// EnableAlb enable an ALB in your cluster
func (r *alb) EnableAlb(enableAlbReq AlbConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/vpc/enableAlb", enableAlbReq, nil, target.ToMap())
	return err
}

// GetAlb returns with the details of an ALB
func (r *alb) GetAlb(albID string, target ClusterTargetHeader) (AlbConfig, error) {
	var successV AlbConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/getAlb?albID=%s", albID), &successV, target.ToMap())
	return successV, err
}

// UpdateAlb update one or more ALBs. To update your ALB to a specified image version, automatic updates must be disabled
func (r *alb) UpdateAlb(updateAlbReq UpdateALBReq, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/updateAlb", updateAlbReq, nil, target.ToMap())
	return err
}

// ListClusterALBs returns the list of albs available for cluster
func (r *alb) ListClusterAlbs(clusterNameOrID string, target ClusterTargetHeader) ([]AlbConfig, error) {
	var successV ClusterALB
	rawURL := fmt.Sprintf("v2/alb/getClusterAlbs?cluster=%s", clusterNameOrID)
	_, err := r.client.Get(rawURL, &successV, target.ToMap())
	return successV.ALBs, err
}

// ListAlbImages lists the default and the supported ALB image versions
func (r *alb) ListAlbImages(target ClusterTargetHeader) (AlbImageVersions, error) {
	var successV AlbImageVersions
	_, err := r.client.Get("v2/alb/getAlbImages", &successV, target.ToMap())
	return successV, err
}

// GetIngressStatus returns the ingress status report for the cluster
func (r *alb) GetIngressStatus(clusterNameOrID string, target ClusterTargetHeader) (IngressStatus, error) {
	var successV IngressStatus
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/getStatus?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// GetAlbClusterHealthCheckConfig returns the ALB in-cluster healthcheck config
func (r *alb) GetAlbClusterHealthCheckConfig(clusterNameOrID string, target ClusterTargetHeader) (ALBClusterHealthCheckConfig, error) {
	var successV ALBClusterHealthCheckConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/getIngressClusterHealthcheck?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// SetAlbClusterHealthCheckConfig configure the ALB in-cluster healthcheck
func (r *alb) SetAlbClusterHealthCheckConfig(albHealthCheckReq ALBClusterHealthCheckConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/setIngressClusterHealthcheck", albHealthCheckReq, nil, target.ToMap())
	return err
}

// GetIgnoredIngressStatusErrors lists of error codes that the user wants to ignore
func (r *alb) GetIgnoredIngressStatusErrors(clusterNameOrID string, target ClusterTargetHeader) (IgnoredIngressStatusErrors, error) {
	var successV IgnoredIngressStatusErrors
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/listIgnoredIngressStatusErrors?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// AddIgnoredIngressStatusErrors ignore one or more ingress status error
func (r *alb) AddIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/addIgnoredIngressStatusErrors", ignoredErrorsReq, nil, target.ToMap())
	return err
}

// RemoveIgnoredIngressStatusErrors remove one or more ignored ingress status error
func (r *alb) RemoveIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.DeleteWithBody("/v2/alb/removeIgnoredIngressStatusErrors", ignoredErrorsReq, nil, target.ToMap())
	return err
}

// SetIngressStatusState set the state of the ingress status for a cluster
func (r *alb) SetIngressStatusState(ingressStatusStateReq IngressStatusState, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/setIngressStatusState", ingressStatusStateReq, nil, target.ToMap())
	return err
}

// GetIngressLoadBalancerConfig get the configuration of load balancers for Ingress ALBs
func (r *alb) GetIngressLoadBalancerConfig(clusterNameOrID, lbType string, target ClusterTargetHeader) (ALBLBConfig, error) {
	var successV ALBLBConfig
	_, err := r.client.Get(fmt.Sprintf("/ingress/v2/load-balancer/configuration?cluster=%s&type=%s", clusterNameOrID, lbType), &successV, target.ToMap())
	return successV, err
}

// UpdateIngressLoadBalancerConfig update the configuration of load balancers for Ingress ALBs
func (r *alb) UpdateIngressLoadBalancerConfig(lbConfig ALBLBConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Patch("/ingress/v2/load-balancer/configuration", lbConfig, nil, target.ToMap())
	return err
}

// GetALBAutoscaleConfiguration get the autoscaling configuration for an ALB
func (r *alb) GetALBAutoscaleConfiguration(clusterNameOrID, albID string, target ClusterTargetHeader) (AutoscaleDetails, error) {
	var successV AutoscaleDetails
	_, err := r.client.Get(fmt.Sprintf("/ingress/v2/clusters/%s/albs/%s/autoscale", clusterNameOrID, albID), &successV, target.ToMap())
	return successV, err
}

// SetALBAutoscaleConfiguration set the autoscaling configuration for an ALB
func (r *alb) SetALBAutoscaleConfiguration(clusterNameOrID, albID string, autoscaleDetails AutoscaleDetails, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Put(fmt.Sprintf("/ingress/v2/clusters/%s/albs/%s/autoscale", clusterNameOrID, albID), autoscaleDetails, nil, target.ToMap())
	return err
}

// RemoveALBAutoscaleConfiguration delete the autoscaling configuration for an ALB
func (r *alb) RemoveALBAutoscaleConfiguration(clusterNameOrID, albID string, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Delete(fmt.Sprintf("/ingress/v2/clusters/%s/albs/%s/autoscale", clusterNameOrID, albID), nil, nil, target.ToMap())
	return err
}
