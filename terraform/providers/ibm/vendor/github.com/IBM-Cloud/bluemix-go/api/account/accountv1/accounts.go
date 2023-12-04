package accountv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type AccountUser struct {
	UserId      string `json:"userId"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	State       string `json:"state"`
	IbmUniqueId string `json:"ibmUniqueId"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	CreatedOn   string `json:"createdOn"`
	VerifiedOn  string `json:"verifiedOn"`
	Id          string `json:"id"`
	UaaGuid     string `json:"uaaGuid"`
	AccountId   string `json:"accountId"`
	Role        string `json:"role"`
	InvitedOn   string `json:"invitedOn"`
	Photo       string `json:"photo"`
}

// Accounts ...
type Accounts interface {
	GetAccountUsers(accountGuid string) ([]models.AccountUser, error)
	InviteAccountUser(accountGuid string, userEmail string) (AccountInviteResponse, error)
	DeleteAccountUser(accountGuid string, userGuid string) error
	FindAccountUserByUserId(accountGuid string, userId string) (*models.AccountUser, error)
}

type account struct {
	client *client.Client
}

type Metadata struct {
	Guid       string   `json:"guid"`
	Url        string   `json:"url"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	VerifiedAt string   `json:"verified_at"`
	Identity   Identity `json:"identity"`
}

type AccountUserEntity struct {
	AccountId   string `json:"account_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	State       string `json:"state"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
	Role        string `json:"role"`
	Photo       string `json:"photo"`
}

type AccountUserMetadata Metadata

type Identity struct {
	Id         string `json:"id"`
	UserName   string `json:"username"`
	Realmid    string `json:"realmid"`
	Identifier string `json:"identifier"`
}

// Account Invites ...
type AccountInviteResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	State string `json:"state"`
}

type AccountUserQueryResponse struct {
	Metadata     Metadata
	AccountUsers []AccountUserResource `json:"resources"`
}

type AccountResource struct {
	Metadata Metadata
	Entity   AccountEntity
}

type AccountEntity struct {
	Name          string                       `json:"name"`
	Type          string                       `json:"type"`
	State         string                       `json:"state"`
	OwnerIamId    string                       `json:"owner_iam_id"`
	CountryCode   string                       `json:"country_code"`
	CurrencyCode  string                       `json:"currency_code"`
	Organizations []models.AccountOrganization `json:"organizations_region"`
}

func (resource AccountResource) ToModel() models.V2Account {
	return models.V2Account{
		Guid:          resource.Metadata.Guid,
		Name:          resource.Entity.Name,
		Type:          resource.Entity.Type,
		State:         resource.Entity.State,
		OwnerIamId:    resource.Entity.OwnerIamId,
		CountryCode:   resource.Entity.CountryCode,
		Organizations: resource.Entity.Organizations,
	}
}

// AccountUserResource is the original user information returned by V2 endpoint (for listing)
type AccountUserResource struct {
	ID             string `json:"id"`
	IAMID          string `json:"iam_id"`
	Realm          string `json:"realm"`
	UserID         string `json:"user_id"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	State          string `json:"state"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phonenumber"`
	AltPhoneNumber string `json:"altphonenumber"`
	Photo          string `json:"photo"`
	InvitedOn      string `json:"invitedOn"`
	AddedOn        string `json:"added_on"`
	AccountID      string `json:"account_id"`

	Linkages []struct {
		Origin string `json:"origin"`
		ID     string `json:"id"`
	} `json:"linkages"`
}

func (resource AccountUserResource) ToModel() models.AccountUser {
	var uaaGUID string
	for _, linkage := range resource.Linkages {
		if linkage.Origin == "UAA" {
			uaaGUID = linkage.ID
			break
		}
	}
	user := models.AccountUser{
		Id:          resource.ID,
		UserId:      resource.UserID,
		FirstName:   resource.FirstName,
		LastName:    resource.LastName,
		State:       resource.State,
		Email:       resource.Email,
		Phonenumber: resource.PhoneNumber,
		AccountId:   resource.AccountID,
		Photo:       resource.Photo,
		IbmUniqueId: resource.IAMID,
		UaaGuid:     uaaGUID,
		AddedOn:     resource.AddedOn,
		InvitedOn:   resource.InvitedOn,
	}

	return user
}

func newAccountAPI(c *client.Client) Accounts {
	return &account{
		client: c,
	}
}

// GetAccountUser ...
func (a *account) GetAccountUsers(accountGuid string) ([]models.AccountUser, error) {
	var users []models.AccountUser

	resp, err := a.client.GetPaginated(fmt.Sprintf("/v2/accounts/%s/users", accountGuid),
		accountv2.NewAccountPaginatedResources(AccountUserResource{}),
		func(resource interface{}) bool {
			if accountUser, ok := resource.(AccountUserResource); ok {
				users = append(users, accountUser.ToModel())
				return true
			}
			return false
		})

	if resp.StatusCode == 404 {
		return []models.AccountUser{}, bmxerror.New(ErrCodeNoAccountExists,
			fmt.Sprintf("No Account exists with account id:%q", accountGuid))
	}

	return users, err
}

// Deprecated: User Invite is deprecated from accounts use UserInvite from userManagement
func (a *account) InviteAccountUser(accountGuid string, userEmail string) (AccountInviteResponse, error) {
	type userEntity struct {
		Email       string `json:"email"`
		AccountRole string `json:"account_role"`
	}

	payload := struct {
		Users []userEntity `json:"users"`
	}{
		Users: []userEntity{
			{
				Email:       userEmail,
				AccountRole: "MEMBER",
			},
		},
	}

	resp := AccountInviteResponse{}

	_, err := a.client.Post(fmt.Sprintf("/v2/accounts/%s/users", accountGuid), payload, &resp)
	return resp, err
}

func (a *account) DeleteAccountUser(accountGuid string, userGuid string) error {
	_, err := a.client.Delete(fmt.Sprintf("/v2/accounts/%s/users/%s", accountGuid, userGuid))

	return err
}

func (a *account) FindAccountUserByUserId(accountGuid string, userId string) (*models.AccountUser, error) {
	queryResp := AccountUserQueryResponse{}

	req := rest.GetRequest(*a.client.Config.Endpoint+fmt.Sprintf("/v2/accounts/%s/users", accountGuid)).
		Query("user_id", userId)

	response, err := a.client.SendRequest(req,
		&queryResp)

	if err != nil {
		switch response.StatusCode {
		case 404:
			return nil, nil
		default:
			return nil, err
		}
	} else if len(queryResp.AccountUsers) == 0 {
		return nil, nil
	} else {
		accountUser := queryResp.AccountUsers[0].ToModel()
		return &accountUser, nil
	}
}
