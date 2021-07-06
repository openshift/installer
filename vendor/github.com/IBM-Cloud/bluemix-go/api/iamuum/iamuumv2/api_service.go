package iamuumv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//IAMUUMServiceAPIv2 is the resource client ...
type IAMUUMServiceAPIv2 interface {
	AccessGroup() AccessGroupRepository
	AccessGroupMember() AccessGroupMemberRepositoryV2
	DynamicRule() DynamicRuleRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//iamService holds the client
type iamuumService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (IAMUUMServiceAPIv2, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.IAMUUMServicev2)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	if config.IAMAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.IAMEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &iamuumService{
		Client: client.New(config, bluemix.IAMUUMServicev2, tokenRefreher),
	}, nil
}

//AccessGroup API
func (a *iamuumService) AccessGroup() AccessGroupRepository {
	return NewAccessGroupRepository(a.Client)
}

// AccessGroupMember API
func (a *iamuumService) AccessGroupMember() AccessGroupMemberRepositoryV2 {
	return NewAccessGroupMemberRepository(a.Client)
}

// Dynamic Rule API
func (a *iamuumService) DynamicRule() DynamicRuleRepository {
	return NewDynamicRuleRepository(a.Client)
}
