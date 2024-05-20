package openstack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
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
	refSplit := strings.Split(ref, "/")
	uuid := refSplit[len(refSplit)-1]
	return uuid
}

func expandKeyManagerContainerV1SecretRefs(secretRefs *schema.Set) []containers.SecretRef {
	l := make([]containers.SecretRef, 0, len(secretRefs.List()))

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
	m := make([]map[string]interface{}, 0, len(sr))

	for _, v := range sr {
		m = append(m, map[string]interface{}{
			"name":       v.Name,
			"secret_ref": v.SecretRef,
		})
	}

	return m
}

func flattenKeyManagerContainerV1Consumers(cr []containers.ConsumerRef) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(cr))

	for _, v := range cr {
		m = append(m, map[string]interface{}{
			"name": v.Name,
			"url":  v.URL,
		})
	}

	return m
}
