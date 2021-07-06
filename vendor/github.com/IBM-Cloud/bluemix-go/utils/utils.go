package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/models"

	"github.com/IBM-Cloud/bluemix-go/crn"
)

func GetLocationFromTargetCRN(crnResource string) string {
	if strings.HasPrefix(crnResource, "bluemix-") {
		return crnResource[len("bluemix-"):]
	} else if strings.HasPrefix(crnResource, "staging-") {
		return crnResource[len("staging-"):]
	} else {
		return crnResource
	}
}

func GenerateSpaceCRN(region models.Region, orgID string, spaceID string) crn.CRN {
	spaceCRN := crn.New(CloudName(region), CloudType(region))
	spaceCRN.ServiceName = crn.ServiceBluemix
	spaceCRN.Region = region.Name
	spaceCRN.ScopeType = crn.ScopeOrganization
	spaceCRN.Scope = orgID
	spaceCRN.ResourceType = crn.ResourceTypeCFSpace
	spaceCRN.Resource = spaceID
	return spaceCRN
}

func CloudName(region models.Region) string {
	regionID := region.ID
	if regionID == "" {
		return ""
	}

	splits := strings.Split(regionID, ":")
	if len(splits) != 3 {
		return ""
	}

	customer := splits[0]
	if customer != "ibm" {
		return customer
	}

	deployment := splits[1]
	switch {
	case deployment == "yp":
		return "bluemix"
	case strings.HasPrefix(deployment, "ys"):
		return "staging"
	default:
		return ""
	}
}

func CloudType(region models.Region) string {
	return region.Type
}

func GenerateBoundToCRN(region models.Region, accountID string) crn.CRN {
	var boundTo crn.CRN
	if region.Type == "dedicated" {
		// cname and ctype are hard coded for dedicated
		boundTo = crn.New("bluemix", "public")
	} else {
		boundTo = crn.New(CloudName(region), CloudType(region))
	}

	boundTo.ScopeType = crn.ScopeAccount
	boundTo.Scope = accountID
	return boundTo
}

func GetRolesFromRoleNamesV2(roleNames []string, roles []iampapv2.Role) ([]iampapv2.Role, error) {

	filteredRoles := []iampapv2.Role{}
	for _, roleName := range roleNames {
		role, err := FindRoleByNameV2(roles, roleName)
		if err != nil {
			return []iampapv2.Role{}, err
		}
		filteredRoles = append(filteredRoles, role)
	}
	return filteredRoles, nil
}
func FindRoleByNameV2(supported []iampapv2.Role, name string) (iampapv2.Role, error) {
	for _, role := range supported {
		if role.DisplayName == name {
			return role, nil
		}
	}
	supportedRoles := getSupportedRolesStringV2(supported)
	return iampapv2.Role{}, bmxerror.New(ErrCodeRRoleDoesnotExist,
		fmt.Sprintf("%s was not found. Valid roles are %s", name, supportedRoles))

}

func getSupportedRolesStringV2(supported []iampapv2.Role) string {
	rolesStr := ""
	for index, role := range supported {
		if index != 0 {
			rolesStr += ", "
		}
		rolesStr += role.DisplayName
	}
	return rolesStr
}

func GetRolesFromRoleNames(roleNames []string, roles []models.PolicyRole) ([]models.PolicyRole, error) {

	filteredRoles := []models.PolicyRole{}
	for _, roleName := range roleNames {
		role, err := FindRoleByName(roles, roleName)
		if err != nil {
			return []models.PolicyRole{}, err
		}
		filteredRoles = append(filteredRoles, role)
	}
	return filteredRoles, nil
}

const ErrCodeRRoleDoesnotExist = "RoleDoesnotExist"

func FindRoleByName(supported []models.PolicyRole, name string) (models.PolicyRole, error) {
	for _, role := range supported {
		if role.DisplayName == name {
			return role, nil
		}
	}
	supportedRoles := getSupportedRolesString(supported)
	return models.PolicyRole{}, bmxerror.New(ErrCodeRRoleDoesnotExist,
		fmt.Sprintf("%s was not found. Valid roles are %s", name, supportedRoles))

}

func getSupportedRolesString(supported []models.PolicyRole) string {
	rolesStr := ""
	for index, role := range supported {
		if index != 0 {
			rolesStr += ", "
		}
		rolesStr += role.DisplayName
	}
	return rolesStr
}

func EscapeUrlParm(urlParm string) string {
	if strings.Contains(urlParm, "/") {
		newUrlParm := url.PathEscape(urlParm)
		return newUrlParm
	}
	return urlParm
}
