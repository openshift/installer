package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
)

func flattenIdentityApplicationCredentialRolesV3(roles []applicationcredentials.Role) []string {
	var res []string
	for _, role := range roles {
		res = append(res, role.Name)
	}
	return res
}

func expandIdentityApplicationCredentialRolesV3(roles []interface{}) []applicationcredentials.Role {
	var res []applicationcredentials.Role
	for _, role := range roles {
		res = append(res, applicationcredentials.Role{Name: role.(string)})
	}
	return res
}

func flattenIdentityApplicationCredentialAccessRulesV3(rules []applicationcredentials.AccessRule) []map[string]string {
	res := make([]map[string]string, len(rules))
	for i, v := range rules {
		res[i] = map[string]string{
			"id":      v.ID,
			"path":    v.Path,
			"method":  v.Method,
			"service": v.Service,
		}
	}
	return res
}

func expandIdentityApplicationCredentialAccessRulesV3(rules []interface{}) []applicationcredentials.AccessRule {
	var res []applicationcredentials.AccessRule
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
