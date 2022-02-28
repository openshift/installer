package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
)

func flattenIdentityApplicationCredentialRolesV3(roles []applicationcredentials.Role) []string {
	res := make([]string, 0, len(roles))
	for _, role := range roles {
		res = append(res, role.Name)
	}
	return res
}

func expandIdentityApplicationCredentialRolesV3(roles []interface{}) []applicationcredentials.Role {
	res := make([]applicationcredentials.Role, 0, len(roles))
	for _, role := range roles {
		res = append(res, applicationcredentials.Role{Name: role.(string)})
	}
	return res
}

func flattenIdentityApplicationCredentialAccessRulesV3(rules []applicationcredentials.AccessRule) []map[string]string {
	res := make([]map[string]string, 0, len(rules))
	for _, role := range rules {
		res = append(res, map[string]string{
			"id":      role.ID,
			"path":    role.Path,
			"method":  role.Method,
			"service": role.Service,
		})
	}
	return res
}

func expandIdentityApplicationCredentialAccessRulesV3(rules []interface{}) []applicationcredentials.AccessRule {
	res := make([]applicationcredentials.AccessRule, 0, len(rules))
	for _, v := range rules {
		rule := v.(map[string]interface{})
		res = append(res,
			applicationcredentials.AccessRule{
				ID:      rule["id"].(string),
				Path:    rule["path"].(string),
				Method:  rule["method"].(string),
				Service: rule["service"].(string),
			},
		)
	}
	return res
}

func applicationCredentialCleanupAccessRulesV3(client *gophercloud.ServiceClient, userID string, id string, rules []applicationcredentials.AccessRule) error {
	for _, rule := range rules {
		log.Printf("[DEBUG] Cleaning up %q access rule from the %q application credential", rule.ID, id)
		err := applicationcredentials.DeleteAccessRule(client, userID, rule.ID).ExtractErr()
		if err != nil {
			switch err.(type) {
			case gophercloud.ErrDefault403:
				// handle "May not delete access rule in use", when the same access rule is
				// used by other application credential
				log.Printf("[DEBUG] Error delete %q access rule from the %q application credential: %s", rule.ID, id, err)
				continue
			case gophercloud.ErrDefault404:
				// access rule was already deleted
				continue
			default:
				return fmt.Errorf("failed to delete %q access rule from the %q application credential: %s", rule.ID, id, err)
			}
		}
	}
	return nil
}
