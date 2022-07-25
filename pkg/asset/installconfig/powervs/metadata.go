package powervs

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	gohttp "net/http"
	"sync"
)

//go:generate mockgen -source=./metadata.go -destination=./mock/powervsmetadata_generated.go -package=mock

// MetadataAPI represents functions that eventually call out to the API
type MetadataAPI interface {
	AccountID(ctx context.Context) (string, error)
	APIKey(ctx context.Context) (string, error)
	CISInstanceCRN(ctx context.Context) (string, error)
}

// Metadata holds additional metadata for InstallConfig resources that
// do not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	BaseDomain string

	accountID      string
	apiKey         string
	cisInstanceCRN string
	client         *Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(baseDomain string) *Metadata {
	return &Metadata{BaseDomain: baseDomain}
}

// AccountID returns the IBM Cloud account ID associated with the authentication
// credentials.
func (m *Metadata) AccountID(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.accountID == "" {
		apiKeyDetails, err := m.client.GetAuthenticatorAPIKeyDetails(ctx)
		if err != nil {
			return "", err
		}

		m.accountID = *apiKeyDetails.AccountID
	}

	return m.accountID, nil
}

// APIKey returns the IBM Cloud account API Key associated with the authentication
// credentials.
func (m *Metadata) APIKey(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.apiKey == "" {
		m.apiKey = m.client.GetAPIKey()
	}

	return m.apiKey, nil
}

// getCISInstanceCRN gets the CRN name for the specified base domain.
func getCISInstanceCRN(APIKey string, BaseDomain string) (string, error) {
	var CISInstanceCRN string = ""
	var bxSession *bxsession.Session
	var err error
	var tokenProviderEndpoint string = "https://iam.cloud.ibm.com"
	var tokenRefresher *authentication.IAMAuthRepository
	var authenticator *core.IamAuthenticator
	var controllerSvc *resourcecontrollerv2.ResourceControllerV2
	var listInstanceOptions *resourcecontrollerv2.ListResourceInstancesOptions
	var listResourceInstancesResponse *resourcecontrollerv2.ResourceInstancesList
	var instance resourcecontrollerv2.ResourceInstance
	var zonesService *zonesv1.ZonesV1
	var listZonesOptions *zonesv1.ListZonesOptions
	var listZonesResponse *zonesv1.ListZonesResp

	bxSession, err = bxsession.New(&bluemix.Config{
		BluemixAPIKey:         APIKey,
		TokenProviderEndpoint: &tokenProviderEndpoint,
		Debug:                 false,
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: bxsession.New: %v", err)
	}
	tokenRefresher, err = authentication.NewIAMAuthRepository(bxSession.Config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: authentication.NewIAMAuthRepository: %v", err)
	}
	err = tokenRefresher.AuthenticateAPIKey(bxSession.Config.BluemixAPIKey)
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: tokenRefresher.AuthenticateAPIKey: %v", err)
	}
	authenticator = &core.IamAuthenticator{
		ApiKey: APIKey,
	}
	err = authenticator.Validate()
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: authenticator.Validate: %v", err)
	}
	// Instantiate the service with an API key based IAM authenticator
	controllerSvc, err = resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
		ServiceName:   "cloud-object-storage",
		URL:           "https://resource-controller.cloud.ibm.com",
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: creating ControllerV2 Service: %v", err)
	}
	listInstanceOptions = controllerSvc.NewListResourceInstancesOptions()
	listInstanceOptions.SetResourceID(cisServiceID)
	listResourceInstancesResponse, _, err = controllerSvc.ListResourceInstances(listInstanceOptions)
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: ListResourceInstances: %v", err)
	}
	for _, instance = range listResourceInstancesResponse.Resources {
		authenticator = &core.IamAuthenticator{
			ApiKey: APIKey,
		}

		err = authenticator.Validate()
		if err != nil {
		}

		zonesService, err = zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: authenticator,
			Crn:           instance.CRN,
		})
		if err != nil {
			return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: NewZonesV1: %v", err)
		}
		listZonesOptions = zonesService.NewListZonesOptions()
		listZonesResponse, _, err = zonesService.ListZones(listZonesOptions)
		if listZonesResponse == nil {
			return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: ListZones: %v", err)
		}
		for _, zone := range listZonesResponse.Result {
			if *zone.Status == "active" {
				if *zone.Name == BaseDomain {
					CISInstanceCRN = *instance.CRN
				}
			}
		}
	}

	return CISInstanceCRN, nil
}

// CISInstanceCRN returns the Cloud Internet Services instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) CISInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.apiKey == "" {
		m.apiKey = m.client.GetAPIKey()
	}

	if m.cisInstanceCRN == "" {
		var cisInstanceCRN string = ""
		var err error

		cisInstanceCRN, err = getCISInstanceCRN(m.apiKey, m.BaseDomain)
		if err != nil {
			return "", err
		}

		m.cisInstanceCRN = cisInstanceCRN
	}

	return m.cisInstanceCRN, nil
}

// SetCISInstanceCRN sets Cloud Internet Services instance CRN to a string value.
func (m *Metadata) SetCISInstanceCRN(crn string) {
	m.cisInstanceCRN = crn
}
