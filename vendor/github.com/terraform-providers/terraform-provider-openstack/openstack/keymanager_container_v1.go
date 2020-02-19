package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func keyManagerContainerV1WaitForContainerDeletion(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := containers.Delete(kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func keyManagerContainerV1Type(v string) containers.ContainerType {
	var ctype containers.ContainerType

	switch v {
	case "generic":
		ctype = containers.GenericContainer
	case "rsa":
		ctype = containers.RSAContainer
	case "certificate":
		ctype = containers.CertificateContainer
	}

	return ctype
}

func keyManagerContainerV1WaitForContainerCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		container, err := containers.Get(kmClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}

		if container.Status == "ERROR" {
			return "", container.Status, fmt.Errorf("Error creating container")
		}

		return container, container.Status, nil
	}
}

func keyManagerContainerV1GetUUIDfromContainerRef(ref string) string {
	// container ref has form https://{barbican_host}/v1/containers/{container_uuid}
	// so we are only interested in the last part
	ref_split := strings.Split(ref, "/")
	uuid := ref_split[len(ref_split)-1]
	return uuid
}

func expandKeyManagerContainerV1SecretRefs(secretRefs *schema.Set) []containers.SecretRef {
	var l []containers.SecretRef

	for _, v := range secretRefs.List() {
		if v, ok := v.(map[string]interface{}); ok {
			var s containers.SecretRef

			if v, ok := v["secret_ref"]; ok {
				s.SecretRef = v.(string)
			}
			if v, ok := v["name"]; ok {
				s.Name = v.(string)
			}

			l = append(l, s)
		}
	}

	return l
}

func flattenKeyManagerContainerV1SecretRefs(sr []containers.SecretRef) []map[string]interface{} {
	var m []map[string]interface{}

	for _, v := range sr {
		m = append(m, map[string]interface{}{
			"name":       v.Name,
			"secret_ref": v.SecretRef,
		})
	}

	return m
}

func flattenKeyManagerContainerV1Consumers(cr []containers.ConsumerRef) []map[string]interface{} {
	var m []map[string]interface{}

	for _, v := range cr {
		m = append(m, map[string]interface{}{
			"name": v.Name,
			"url":  v.URL,
		})
	}

	return m
}
