package secretsmanager

import (
	"fmt"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"strconv"
	"strings"
	"time"
)

func getRegion(originalClient *secretsmanagerv2.SecretsManagerV2, d *schema.ResourceData) string {
	_, ok := d.GetOk("region")
	if ok {
		return d.Get("region").(string)
	} else {
		// extract region from base URL (provider config)
		// base url is like that : "https://<private.>secrets-manager.<region>.<rest of domain>"
		baseUrl := originalClient.Service.GetServiceURL()
		u := strings.Replace(baseUrl, "private.", "", 1)
		return strings.Split(u, ".")[1]
	}
}

// Clone the base secrets manager client and set the API endpoint per the instance
func getEndpointType(originalClient *secretsmanagerv2.SecretsManagerV2, d *schema.ResourceData) string {
	_, ok := d.GetOk("endpoint_type")
	if ok {
		return d.Get("endpoint_type").(string)
	} else {
		baseUrl := originalClient.Service.GetServiceURL()

		if strings.Contains(baseUrl, "private.") {
			return "private"
		} else {
			return "public"
		}
	}
}

// Clone the base secrets manager client and set the API endpoint per the instance
func getClientWithInstanceEndpoint(originalClient *secretsmanagerv2.SecretsManagerV2, instanceId string, region string, endpointType string) *secretsmanagerv2.SecretsManagerV2 {
	// build the api endpoint
	domain := "appdomain.cloud"
	if strings.Contains(os.Getenv("IBMCLOUD_IAM_API_ENDPOINT"), "test") {
		domain = "test.appdomain.cloud"
	}
	var endpoint string
	if endpointType == "private" {
		endpoint = fmt.Sprintf("https://%s.private.%s.secrets-manager.%s", instanceId, region, domain)
	} else {
		endpoint = fmt.Sprintf("https://%s.%s.secrets-manager.%s", instanceId, region, domain)
	}

	// clone the client and set endpoint
	newClient := &secretsmanagerv2.SecretsManagerV2{
		Service: originalClient.Service.Clone(),
	}
	newClient.Service.SetServiceURL(endpoint)
	return newClient
}

// Add the fields needed for building the instance endpoint to the given schema
func AddInstanceFields(resource *schema.Resource) *schema.Resource {
	resource.Schema["instance_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the Secrets Manager instance.",
	}
	resource.Schema["region"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "The region of the Secrets Manager instance.",
	}
	resource.Schema["endpoint_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "public or private.",
	}

	return resource
}

func StringIsIntBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		vs, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		v, err := strconv.Atoi(vs)

		if err != nil {
			errors = append(errors, fmt.Errorf("expected %s to represent an integer", k))
			return warnings, errors
		}

		if v < min || v > max {
			errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return warnings, errors
		}

		return warnings, errors
	}
}

func DateTimeToRFC3339(dt *strfmt.DateTime) (s string) {
	if dt != nil {
		s = time.Time(*dt).Format(time.RFC3339)
	}
	return
}
