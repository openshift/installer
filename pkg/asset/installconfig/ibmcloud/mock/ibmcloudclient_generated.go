// Code generated by MockGen. DO NOT EDIT.
// Source: ./client.go
//
// Generated by this command:
//
//	mockgen -source=./client.go -destination=./mock/ibmcloudclient_generated.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	s3 "github.com/IBM/ibm-cos-sdk-go/service/s3"
	dnsrecordsv1 "github.com/IBM/networking-go-sdk/dnsrecordsv1"
	iamidentityv1 "github.com/IBM/platform-services-go-sdk/iamidentityv1"
	resourcecontrollerv2 "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	resourcemanagerv2 "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	vpcv1 "github.com/IBM/vpc-go-sdk/vpcv1"
	responses "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/responses"
	types "github.com/openshift/installer/pkg/types"
	gomock "go.uber.org/mock/gomock"
)

// MockAPI is a mock of API interface.
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
	isgomock struct{}
}

// MockAPIMockRecorder is the mock recorder for MockAPI.
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance.
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// AttachFloatingIP mocks base method.
func (m *MockAPI) AttachFloatingIP(ctx context.Context, instanceName, instanceID, region, resourceGroupName string) (*vpcv1.FloatingIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachFloatingIP", ctx, instanceName, instanceID, region, resourceGroupName)
	ret0, _ := ret[0].(*vpcv1.FloatingIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AttachFloatingIP indicates an expected call of AttachFloatingIP.
func (mr *MockAPIMockRecorder) AttachFloatingIP(ctx, instanceName, instanceID, region, resourceGroupName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachFloatingIP", reflect.TypeOf((*MockAPI)(nil).AttachFloatingIP), ctx, instanceName, instanceID, region, resourceGroupName)
}

// CreateCISDNSRecord mocks base method.
func (m *MockAPI) CreateCISDNSRecord(ctx context.Context, cisInstanceCRN, zoneID, recordName, cname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCISDNSRecord", ctx, cisInstanceCRN, zoneID, recordName, cname)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCISDNSRecord indicates an expected call of CreateCISDNSRecord.
func (mr *MockAPIMockRecorder) CreateCISDNSRecord(ctx, cisInstanceCRN, zoneID, recordName, cname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCISDNSRecord", reflect.TypeOf((*MockAPI)(nil).CreateCISDNSRecord), ctx, cisInstanceCRN, zoneID, recordName, cname)
}

// CreateCOSBucket mocks base method.
func (m *MockAPI) CreateCOSBucket(ctx context.Context, cosInstanceID, bucketName, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCOSBucket", ctx, cosInstanceID, bucketName, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCOSBucket indicates an expected call of CreateCOSBucket.
func (mr *MockAPIMockRecorder) CreateCOSBucket(ctx, cosInstanceID, bucketName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCOSBucket", reflect.TypeOf((*MockAPI)(nil).CreateCOSBucket), ctx, cosInstanceID, bucketName, region)
}

// CreateCOSInstance mocks base method.
func (m *MockAPI) CreateCOSInstance(ctx context.Context, cosName, resourceGroupID string) (*resourcecontrollerv2.ResourceInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCOSInstance", ctx, cosName, resourceGroupID)
	ret0, _ := ret[0].(*resourcecontrollerv2.ResourceInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCOSInstance indicates an expected call of CreateCOSInstance.
func (mr *MockAPIMockRecorder) CreateCOSInstance(ctx, cosName, resourceGroupID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCOSInstance", reflect.TypeOf((*MockAPI)(nil).CreateCOSInstance), ctx, cosName, resourceGroupID)
}

// CreateCOSObject mocks base method.
func (m *MockAPI) CreateCOSObject(ctx context.Context, sourceData []byte, fileName, cosInstanceID, bucketName, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCOSObject", ctx, sourceData, fileName, cosInstanceID, bucketName, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCOSObject indicates an expected call of CreateCOSObject.
func (mr *MockAPIMockRecorder) CreateCOSObject(ctx, sourceData, fileName, cosInstanceID, bucketName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCOSObject", reflect.TypeOf((*MockAPI)(nil).CreateCOSObject), ctx, sourceData, fileName, cosInstanceID, bucketName, region)
}

// CreateDNSServicesDNSRecord mocks base method.
func (m *MockAPI) CreateDNSServicesDNSRecord(ctx context.Context, dnsInstanceID, zoneID, recordName, cname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDNSServicesDNSRecord", ctx, dnsInstanceID, zoneID, recordName, cname)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDNSServicesDNSRecord indicates an expected call of CreateDNSServicesDNSRecord.
func (mr *MockAPIMockRecorder) CreateDNSServicesDNSRecord(ctx, dnsInstanceID, zoneID, recordName, cname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDNSServicesDNSRecord", reflect.TypeOf((*MockAPI)(nil).CreateDNSServicesDNSRecord), ctx, dnsInstanceID, zoneID, recordName, cname)
}

// CreateDNSServicesPermittedNetwork mocks base method.
func (m *MockAPI) CreateDNSServicesPermittedNetwork(ctx context.Context, dnsInstanceID, dnsZoneID, vpcCRN string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDNSServicesPermittedNetwork", ctx, dnsInstanceID, dnsZoneID, vpcCRN)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDNSServicesPermittedNetwork indicates an expected call of CreateDNSServicesPermittedNetwork.
func (mr *MockAPIMockRecorder) CreateDNSServicesPermittedNetwork(ctx, dnsInstanceID, dnsZoneID, vpcCRN any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDNSServicesPermittedNetwork", reflect.TypeOf((*MockAPI)(nil).CreateDNSServicesPermittedNetwork), ctx, dnsInstanceID, dnsZoneID, vpcCRN)
}

// CreateIAMAuthorizationPolicy mocks base method.
func (m *MockAPI) CreateIAMAuthorizationPolicy(tx context.Context, sourceServiceName, sourceServiceResourceType, targetServiceName, targetServiceInstanceID string, roles []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIAMAuthorizationPolicy", tx, sourceServiceName, sourceServiceResourceType, targetServiceName, targetServiceInstanceID, roles)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIAMAuthorizationPolicy indicates an expected call of CreateIAMAuthorizationPolicy.
func (mr *MockAPIMockRecorder) CreateIAMAuthorizationPolicy(tx, sourceServiceName, sourceServiceResourceType, targetServiceName, targetServiceInstanceID, roles any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIAMAuthorizationPolicy", reflect.TypeOf((*MockAPI)(nil).CreateIAMAuthorizationPolicy), tx, sourceServiceName, sourceServiceResourceType, targetServiceName, targetServiceInstanceID, roles)
}

// CreateResourceGroup mocks base method.
func (m *MockAPI) CreateResourceGroup(ctx context.Context, rgName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateResourceGroup", ctx, rgName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateResourceGroup indicates an expected call of CreateResourceGroup.
func (mr *MockAPIMockRecorder) CreateResourceGroup(ctx, rgName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateResourceGroup", reflect.TypeOf((*MockAPI)(nil).CreateResourceGroup), ctx, rgName)
}

// DeleteCOSBucket mocks base method.
func (m *MockAPI) DeleteCOSBucket(ctx context.Context, cosInstanceID, bucketName, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCOSBucket", ctx, cosInstanceID, bucketName, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCOSBucket indicates an expected call of DeleteCOSBucket.
func (mr *MockAPIMockRecorder) DeleteCOSBucket(ctx, cosInstanceID, bucketName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCOSBucket", reflect.TypeOf((*MockAPI)(nil).DeleteCOSBucket), ctx, cosInstanceID, bucketName, region)
}

// DeleteCOSInstance mocks base method.
func (m *MockAPI) DeleteCOSInstance(ctx context.Context, cosInstanceID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCOSInstance", ctx, cosInstanceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCOSInstance indicates an expected call of DeleteCOSInstance.
func (mr *MockAPIMockRecorder) DeleteCOSInstance(ctx, cosInstanceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCOSInstance", reflect.TypeOf((*MockAPI)(nil).DeleteCOSInstance), ctx, cosInstanceID)
}

// DeleteCOSObject mocks base method.
func (m *MockAPI) DeleteCOSObject(ctx context.Context, cosInstanceID, bucketName, objectKey, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCOSObject", ctx, cosInstanceID, bucketName, objectKey, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCOSObject indicates an expected call of DeleteCOSObject.
func (mr *MockAPIMockRecorder) DeleteCOSObject(ctx, cosInstanceID, bucketName, objectKey, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCOSObject", reflect.TypeOf((*MockAPI)(nil).DeleteCOSObject), ctx, cosInstanceID, bucketName, objectKey, region)
}

// DeleteFloatingIP mocks base method.
func (m *MockAPI) DeleteFloatingIP(ctx context.Context, floatingIPID, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFloatingIP", ctx, floatingIPID, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFloatingIP indicates an expected call of DeleteFloatingIP.
func (mr *MockAPIMockRecorder) DeleteFloatingIP(ctx, floatingIPID, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFloatingIP", reflect.TypeOf((*MockAPI)(nil).DeleteFloatingIP), ctx, floatingIPID, region)
}

// DeleteSecurityGroup mocks base method.
func (m *MockAPI) DeleteSecurityGroup(ctx context.Context, securityGroupID, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecurityGroup", ctx, securityGroupID, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecurityGroup indicates an expected call of DeleteSecurityGroup.
func (mr *MockAPIMockRecorder) DeleteSecurityGroup(ctx, securityGroupID, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecurityGroup", reflect.TypeOf((*MockAPI)(nil).DeleteSecurityGroup), ctx, securityGroupID, region)
}

// DeleteSecurityGroupTargetBinding mocks base method.
func (m *MockAPI) DeleteSecurityGroupTargetBinding(ctx context.Context, securityGroupID, targetID, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecurityGroupTargetBinding", ctx, securityGroupID, targetID, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecurityGroupTargetBinding indicates an expected call of DeleteSecurityGroupTargetBinding.
func (mr *MockAPIMockRecorder) DeleteSecurityGroupTargetBinding(ctx, securityGroupID, targetID, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecurityGroupTargetBinding", reflect.TypeOf((*MockAPI)(nil).DeleteSecurityGroupTargetBinding), ctx, securityGroupID, targetID, region)
}

// GetAPIKey mocks base method.
func (m *MockAPI) GetAPIKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAPIKey indicates an expected call of GetAPIKey.
func (mr *MockAPIMockRecorder) GetAPIKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKey", reflect.TypeOf((*MockAPI)(nil).GetAPIKey))
}

// GetAuthenticatorAPIKeyDetails mocks base method.
func (m *MockAPI) GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthenticatorAPIKeyDetails", ctx)
	ret0, _ := ret[0].(*iamidentityv1.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthenticatorAPIKeyDetails indicates an expected call of GetAuthenticatorAPIKeyDetails.
func (mr *MockAPIMockRecorder) GetAuthenticatorAPIKeyDetails(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthenticatorAPIKeyDetails", reflect.TypeOf((*MockAPI)(nil).GetAuthenticatorAPIKeyDetails), ctx)
}

// GetCISInstance mocks base method.
func (m *MockAPI) GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCISInstance", ctx, crnstr)
	ret0, _ := ret[0].(*resourcecontrollerv2.ResourceInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCISInstance indicates an expected call of GetCISInstance.
func (mr *MockAPIMockRecorder) GetCISInstance(ctx, crnstr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCISInstance", reflect.TypeOf((*MockAPI)(nil).GetCISInstance), ctx, crnstr)
}

// GetCOSBucketByName mocks base method.
func (m *MockAPI) GetCOSBucketByName(ctx context.Context, cosInstanceID, bucketName, region string) (*s3.Bucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCOSBucketByName", ctx, cosInstanceID, bucketName, region)
	ret0, _ := ret[0].(*s3.Bucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCOSBucketByName indicates an expected call of GetCOSBucketByName.
func (mr *MockAPIMockRecorder) GetCOSBucketByName(ctx, cosInstanceID, bucketName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCOSBucketByName", reflect.TypeOf((*MockAPI)(nil).GetCOSBucketByName), ctx, cosInstanceID, bucketName, region)
}

// GetCOSInstanceByName mocks base method.
func (m *MockAPI) GetCOSInstanceByName(ctx context.Context, cosName string) (*resourcecontrollerv2.ResourceInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCOSInstanceByName", ctx, cosName)
	ret0, _ := ret[0].(*resourcecontrollerv2.ResourceInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCOSInstanceByName indicates an expected call of GetCOSInstanceByName.
func (mr *MockAPIMockRecorder) GetCOSInstanceByName(ctx, cosName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCOSInstanceByName", reflect.TypeOf((*MockAPI)(nil).GetCOSInstanceByName), ctx, cosName)
}

// GetDNSInstance mocks base method.
func (m *MockAPI) GetDNSInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDNSInstance", ctx, crnstr)
	ret0, _ := ret[0].(*resourcecontrollerv2.ResourceInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDNSInstance indicates an expected call of GetDNSInstance.
func (mr *MockAPIMockRecorder) GetDNSInstance(ctx, crnstr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDNSInstance", reflect.TypeOf((*MockAPI)(nil).GetDNSInstance), ctx, crnstr)
}

// GetDNSInstancePermittedNetworks mocks base method.
func (m *MockAPI) GetDNSInstancePermittedNetworks(ctx context.Context, dnsID, dnsZone string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDNSInstancePermittedNetworks", ctx, dnsID, dnsZone)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDNSInstancePermittedNetworks indicates an expected call of GetDNSInstancePermittedNetworks.
func (mr *MockAPIMockRecorder) GetDNSInstancePermittedNetworks(ctx, dnsID, dnsZone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDNSInstancePermittedNetworks", reflect.TypeOf((*MockAPI)(nil).GetDNSInstancePermittedNetworks), ctx, dnsID, dnsZone)
}

// GetDNSRecordsByName mocks base method.
func (m *MockAPI) GetDNSRecordsByName(ctx context.Context, crnstr, zoneID, recordName string) ([]dnsrecordsv1.DnsrecordDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDNSRecordsByName", ctx, crnstr, zoneID, recordName)
	ret0, _ := ret[0].([]dnsrecordsv1.DnsrecordDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDNSRecordsByName indicates an expected call of GetDNSRecordsByName.
func (mr *MockAPIMockRecorder) GetDNSRecordsByName(ctx, crnstr, zoneID, recordName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDNSRecordsByName", reflect.TypeOf((*MockAPI)(nil).GetDNSRecordsByName), ctx, crnstr, zoneID, recordName)
}

// GetDNSZoneIDByName mocks base method.
func (m *MockAPI) GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDNSZoneIDByName", ctx, name, publish)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDNSZoneIDByName indicates an expected call of GetDNSZoneIDByName.
func (mr *MockAPIMockRecorder) GetDNSZoneIDByName(ctx, name, publish any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDNSZoneIDByName", reflect.TypeOf((*MockAPI)(nil).GetDNSZoneIDByName), ctx, name, publish)
}

// GetDNSZones mocks base method.
func (m *MockAPI) GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]responses.DNSZoneResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDNSZones", ctx, publish)
	ret0, _ := ret[0].([]responses.DNSZoneResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDNSZones indicates an expected call of GetDNSZones.
func (mr *MockAPIMockRecorder) GetDNSZones(ctx, publish any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDNSZones", reflect.TypeOf((*MockAPI)(nil).GetDNSZones), ctx, publish)
}

// GetDedicatedHostByName mocks base method.
func (m *MockAPI) GetDedicatedHostByName(ctx context.Context, name, region string) (*vpcv1.DedicatedHost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDedicatedHostByName", ctx, name, region)
	ret0, _ := ret[0].(*vpcv1.DedicatedHost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDedicatedHostByName indicates an expected call of GetDedicatedHostByName.
func (mr *MockAPIMockRecorder) GetDedicatedHostByName(ctx, name, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDedicatedHostByName", reflect.TypeOf((*MockAPI)(nil).GetDedicatedHostByName), ctx, name, region)
}

// GetDedicatedHostProfiles mocks base method.
func (m *MockAPI) GetDedicatedHostProfiles(ctx context.Context, region string) ([]vpcv1.DedicatedHostProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDedicatedHostProfiles", ctx, region)
	ret0, _ := ret[0].([]vpcv1.DedicatedHostProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDedicatedHostProfiles indicates an expected call of GetDedicatedHostProfiles.
func (mr *MockAPIMockRecorder) GetDedicatedHostProfiles(ctx, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDedicatedHostProfiles", reflect.TypeOf((*MockAPI)(nil).GetDedicatedHostProfiles), ctx, region)
}

// GetEncryptionKey mocks base method.
func (m *MockAPI) GetEncryptionKey(ctx context.Context, keyCRN string) (*responses.EncryptionKeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEncryptionKey", ctx, keyCRN)
	ret0, _ := ret[0].(*responses.EncryptionKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEncryptionKey indicates an expected call of GetEncryptionKey.
func (mr *MockAPIMockRecorder) GetEncryptionKey(ctx, keyCRN any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEncryptionKey", reflect.TypeOf((*MockAPI)(nil).GetEncryptionKey), ctx, keyCRN)
}

// GetFloatingIPByName mocks base method.
func (m *MockAPI) GetFloatingIPByName(ctx context.Context, floatingIPName, region string) (*vpcv1.FloatingIP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloatingIPByName", ctx, floatingIPName, region)
	ret0, _ := ret[0].(*vpcv1.FloatingIP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFloatingIPByName indicates an expected call of GetFloatingIPByName.
func (mr *MockAPIMockRecorder) GetFloatingIPByName(ctx, floatingIPName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloatingIPByName", reflect.TypeOf((*MockAPI)(nil).GetFloatingIPByName), ctx, floatingIPName, region)
}

// GetIBMCloudRegions mocks base method.
func (m *MockAPI) GetIBMCloudRegions(ctx context.Context) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIBMCloudRegions", ctx)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIBMCloudRegions indicates an expected call of GetIBMCloudRegions.
func (mr *MockAPIMockRecorder) GetIBMCloudRegions(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIBMCloudRegions", reflect.TypeOf((*MockAPI)(nil).GetIBMCloudRegions), ctx)
}

// GetLoadBalancer mocks base method.
func (m *MockAPI) GetLoadBalancer(ctx context.Context, loadBalancerID string) (*vpcv1.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoadBalancer", ctx, loadBalancerID)
	ret0, _ := ret[0].(*vpcv1.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoadBalancer indicates an expected call of GetLoadBalancer.
func (mr *MockAPIMockRecorder) GetLoadBalancer(ctx, loadBalancerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoadBalancer", reflect.TypeOf((*MockAPI)(nil).GetLoadBalancer), ctx, loadBalancerID)
}

// GetResourceGroup mocks base method.
func (m *MockAPI) GetResourceGroup(ctx context.Context, nameOrID string) (*resourcemanagerv2.ResourceGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceGroup", ctx, nameOrID)
	ret0, _ := ret[0].(*resourcemanagerv2.ResourceGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourceGroup indicates an expected call of GetResourceGroup.
func (mr *MockAPIMockRecorder) GetResourceGroup(ctx, nameOrID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceGroup", reflect.TypeOf((*MockAPI)(nil).GetResourceGroup), ctx, nameOrID)
}

// GetResourceGroups mocks base method.
func (m *MockAPI) GetResourceGroups(ctx context.Context) ([]resourcemanagerv2.ResourceGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceGroups", ctx)
	ret0, _ := ret[0].([]resourcemanagerv2.ResourceGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourceGroups indicates an expected call of GetResourceGroups.
func (mr *MockAPIMockRecorder) GetResourceGroups(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceGroups", reflect.TypeOf((*MockAPI)(nil).GetResourceGroups), ctx)
}

// GetSSHKeyByPublicKey mocks base method.
func (m *MockAPI) GetSSHKeyByPublicKey(ctx context.Context, publicKey string) (*vpcv1.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSSHKeyByPublicKey", ctx, publicKey)
	ret0, _ := ret[0].(*vpcv1.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSSHKeyByPublicKey indicates an expected call of GetSSHKeyByPublicKey.
func (mr *MockAPIMockRecorder) GetSSHKeyByPublicKey(ctx, publicKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSSHKeyByPublicKey", reflect.TypeOf((*MockAPI)(nil).GetSSHKeyByPublicKey), ctx, publicKey)
}

// GetSecurityGroupByName mocks base method.
func (m *MockAPI) GetSecurityGroupByName(ctx context.Context, sgName, vpcID, region string) (*vpcv1.SecurityGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecurityGroupByName", ctx, sgName, vpcID, region)
	ret0, _ := ret[0].(*vpcv1.SecurityGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecurityGroupByName indicates an expected call of GetSecurityGroupByName.
func (mr *MockAPIMockRecorder) GetSecurityGroupByName(ctx, sgName, vpcID, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecurityGroupByName", reflect.TypeOf((*MockAPI)(nil).GetSecurityGroupByName), ctx, sgName, vpcID, region)
}

// GetSubnet mocks base method.
func (m *MockAPI) GetSubnet(ctx context.Context, subnetID string) (*vpcv1.Subnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnet", ctx, subnetID)
	ret0, _ := ret[0].(*vpcv1.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnet indicates an expected call of GetSubnet.
func (mr *MockAPIMockRecorder) GetSubnet(ctx, subnetID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnet", reflect.TypeOf((*MockAPI)(nil).GetSubnet), ctx, subnetID)
}

// GetSubnetByName mocks base method.
func (m *MockAPI) GetSubnetByName(ctx context.Context, subnetName, region string) (*vpcv1.Subnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnetByName", ctx, subnetName, region)
	ret0, _ := ret[0].(*vpcv1.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnetByName indicates an expected call of GetSubnetByName.
func (mr *MockAPIMockRecorder) GetSubnetByName(ctx, subnetName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnetByName", reflect.TypeOf((*MockAPI)(nil).GetSubnetByName), ctx, subnetName, region)
}

// GetVPC mocks base method.
func (m *MockAPI) GetVPC(ctx context.Context, vpcID string) (*vpcv1.VPC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVPC", ctx, vpcID)
	ret0, _ := ret[0].(*vpcv1.VPC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVPC indicates an expected call of GetVPC.
func (mr *MockAPIMockRecorder) GetVPC(ctx, vpcID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVPC", reflect.TypeOf((*MockAPI)(nil).GetVPC), ctx, vpcID)
}

// GetVPCByName mocks base method.
func (m *MockAPI) GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVPCByName", ctx, vpcName)
	ret0, _ := ret[0].(*vpcv1.VPC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVPCByName indicates an expected call of GetVPCByName.
func (mr *MockAPIMockRecorder) GetVPCByName(ctx, vpcName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVPCByName", reflect.TypeOf((*MockAPI)(nil).GetVPCByName), ctx, vpcName)
}

// GetVPCZonesForRegion mocks base method.
func (m *MockAPI) GetVPCZonesForRegion(ctx context.Context, region string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVPCZonesForRegion", ctx, region)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVPCZonesForRegion indicates an expected call of GetVPCZonesForRegion.
func (mr *MockAPIMockRecorder) GetVPCZonesForRegion(ctx, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVPCZonesForRegion", reflect.TypeOf((*MockAPI)(nil).GetVPCZonesForRegion), ctx, region)
}

// GetVPCs mocks base method.
func (m *MockAPI) GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVPCs", ctx, region)
	ret0, _ := ret[0].([]vpcv1.VPC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVPCs indicates an expected call of GetVPCs.
func (mr *MockAPIMockRecorder) GetVPCs(ctx, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVPCs", reflect.TypeOf((*MockAPI)(nil).GetVPCs), ctx, region)
}

// GetVSI mocks base method.
func (m *MockAPI) GetVSI(ctx context.Context, instanceID, region string) (*vpcv1.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVSI", ctx, instanceID, region)
	ret0, _ := ret[0].(*vpcv1.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVSI indicates an expected call of GetVSI.
func (mr *MockAPIMockRecorder) GetVSI(ctx, instanceID, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVSI", reflect.TypeOf((*MockAPI)(nil).GetVSI), ctx, instanceID, region)
}

// GetVSIProfiles mocks base method.
func (m *MockAPI) GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVSIProfiles", ctx)
	ret0, _ := ret[0].([]vpcv1.InstanceProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVSIProfiles indicates an expected call of GetVSIProfiles.
func (mr *MockAPIMockRecorder) GetVSIProfiles(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVSIProfiles", reflect.TypeOf((*MockAPI)(nil).GetVSIProfiles), ctx)
}

// ListCOSBuckets mocks base method.
func (m *MockAPI) ListCOSBuckets(ctx context.Context, cosInstanceID, region string) (*s3.ListBucketsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCOSBuckets", ctx, cosInstanceID, region)
	ret0, _ := ret[0].(*s3.ListBucketsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCOSBuckets indicates an expected call of ListCOSBuckets.
func (mr *MockAPIMockRecorder) ListCOSBuckets(ctx, cosInstanceID, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCOSBuckets", reflect.TypeOf((*MockAPI)(nil).ListCOSBuckets), ctx, cosInstanceID, region)
}

// ListCOSObjects mocks base method.
func (m *MockAPI) ListCOSObjects(ctx context.Context, cosInstanceID, bucketName, region string) (*s3.ListObjectsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCOSObjects", ctx, cosInstanceID, bucketName, region)
	ret0, _ := ret[0].(*s3.ListObjectsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCOSObjects indicates an expected call of ListCOSObjects.
func (mr *MockAPIMockRecorder) ListCOSObjects(ctx, cosInstanceID, bucketName, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCOSObjects", reflect.TypeOf((*MockAPI)(nil).ListCOSObjects), ctx, cosInstanceID, bucketName, region)
}

// SetVPCServiceURLForRegion mocks base method.
func (m *MockAPI) SetVPCServiceURLForRegion(ctx context.Context, region string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetVPCServiceURLForRegion", ctx, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetVPCServiceURLForRegion indicates an expected call of SetVPCServiceURLForRegion.
func (mr *MockAPIMockRecorder) SetVPCServiceURLForRegion(ctx, region any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVPCServiceURLForRegion", reflect.TypeOf((*MockAPI)(nil).SetVPCServiceURLForRegion), ctx, region)
}
