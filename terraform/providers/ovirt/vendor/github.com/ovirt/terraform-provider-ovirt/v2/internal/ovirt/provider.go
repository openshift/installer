package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

var providerSchema = map[string]*schema.Schema{
	"username": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Username and realm for oVirt authentication. Required when mock = false. Example: `admin@internal`",
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Password for oVirt authentication. Required when mock = false.",
	},
	"url": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "URL for the oVirt engine API. Required when mock = false. Example: `https://example.com/ovirt-engine/api/`",
	},
	"extra_headers": {
		Type:        schema.TypeMap,
		Optional:    true,
		Elem:        schema.TypeString,
		Description: "Additional HTTP headers to set on each API call.",
	},
	"tls_insecure": {
		Type:             schema.TypeBool,
		Optional:         true,
		ValidateDiagFunc: validateTLSInsecure,
		Description:      "Disable certificate verification when connecting the Engine. This is not recommended. Setting this option is incompatible with other `tls_` options.",
	},
	"tls_system": {
		Type:             schema.TypeBool,
		Optional:         true,
		ValidateDiagFunc: validateTLSSystem,
		Description:      "Use the system certificate pool to verify the Engine certificate. This does not work on Windows. Can be used in parallel with other `tls_` options, one tls_ option is required when mock = false.",
	},
	"tls_ca_bundle": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Validate the Engine certificate against the provided CA certificates. The certificate chain passed should be in PEM format. Can be used in parallel with other `tls_` options, one `tls_` option is required when mock = false.",
	},
	"tls_ca_files": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Validate the Engine certificate against the CA certificates provided in the files in this parameter. The files should contain certificates in PEM format. Can be used in parallel with other tls_ options, one tls_ option is required when mock = false.",
		// Validating TypeList fields is not yet supported in Terraform.
		//ValidateDiagFunc: validateFilesExist,
	},
	"tls_ca_dirs": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Validate the engine certificate against the CA certificates provided in the specified directories. The directory should contain only files with certificates in PEM format. Can be used in parallel with other tls_ options, one tls_ option is required when mock = false.",
		// Validating TypeList fields is not yet supported in Terraform.
		//ValidateDiagFunc: validateDirsExist,
	},
	"mock": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "When set to true, the Terraform provider runs against an internal simulation. This should only be used for testing when an oVirt engine is not available as the mock backend does not persist state across runs. When set to false, one of the tls_ options is required.",
	},
}

// New returns a new Terraform provider schema for oVirt.
func New() func() *schema.Provider {
	return newProvider(newTerraformLogger()).getProvider
}

func newProvider(logger ovirtclientlog.Logger) providerInterface {
	helper, err := ovirtclient.NewMockTestHelper(
		logger,
	)
	if err != nil {
		panic(err)
	}
	return &provider{
		testHelper: helper,
	}
}

type providerInterface interface {
	getTestHelper() ovirtclient.TestHelper
	getProvider() *schema.Provider
	getProviderFactories() map[string]func() (*schema.Provider, error)
}

type provider struct {
	testHelper ovirtclient.TestHelper
	client     ovirtclient.Client
}

func (p *provider) getTestHelper() ovirtclient.TestHelper {
	return p.testHelper
}

func (p *provider) getProvider() *schema.Provider {
	return &schema.Provider{
		Schema:               providerSchema,
		ConfigureContextFunc: p.configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"ovirt_affinity_group":           p.affinityGroupResource(),
			"ovirt_vm_affinity_group":        p.vmAffinityGroupResource(),
			"ovirt_vm":                       p.vmResource(),
			"ovirt_vm_graphics_consoles":     p.vmGraphicsConsolesResource(),
			"ovirt_vm_start":                 p.vmStartResource(),
			"ovirt_vm_tag":                   p.vmTagResource(),
			"ovirt_vm_optimize_cpu_settings": p.vmOptimizeCPUSettingsResource(),
			"ovirt_disk":                     p.diskResource(),
			"ovirt_disk_resize":              p.diskResizeResource(),
			"ovirt_vm_disks_resize":          p.vmDisksResizeResource(),
			"ovirt_disk_from_image":          p.diskFromImageResource(),
			"ovirt_disk_attachment":          p.diskAttachmentResource(),
			"ovirt_disk_attachments":         p.diskAttachmentsResource(),
			"ovirt_nic":                      p.nicResource(),
			"ovirt_tag":                      p.tagResource(),
			"ovirt_template":                 p.templateResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ovirt_blank_template":            p.blankTemplateDataSource(),
			"ovirt_disk_attachments":          p.diskAttachmentsDataSource(),
			"ovirt_template_disk_attachments": p.templateDiskAttachmentsDataSource(),
			"ovirt_cluster_hosts":             p.clusterHostsDataSource(),
			"ovirt_templates":                 p.templatesDataSource(),
			"ovirt_affinity_group":            p.affinityGroupDataSource(),
			"ovirt_wait_for_ip":               p.waitForIPDataSource(),
		},
	}
}

func (p *provider) getProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"ovirt": func() (*schema.Provider, error) { //nolint:unparam
			return p.getProvider(), nil
		},
	}
}

func (p *provider) configureProvider(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if mock, ok := data.GetOk("mock"); ok && mock == true {
		p.client = p.testHelper.GetClient()
		return p, diags
	}

	url, diags := extractString(data, "url", diags)
	username, diags := extractString(data, "username", diags)
	password, diags := extractString(data, "password", diags)

	tls := ovirtclient.TLS()
	if insecure, ok := data.GetOk("tls_insecure"); ok && insecure == true {
		tls.Insecure()
	}
	if system, ok := data.GetOk("tls_system"); ok && system == true {
		tls.CACertsFromSystem()
	}

	caFiles, diags := getStringSliceFromResource("tls_ca_files", data, diags)
	for _, caFile := range caFiles {
		tls.CACertsFromFile(caFile)
	}

	caDirs, diags := getStringSliceFromResource("tls_ca_dirs", data, diags)
	for _, caDir := range caDirs {
		tls.CACertsFromDir(caDir)
	}

	if caBundle, ok := data.GetOk("tls_ca_bundle"); ok {
		caCerts, ok := caBundle.(string)
		if !ok {
			diags = append(
				diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "The tls_ca_bundle option is not a string",
					Detail:   "The tls_ca_bundle option must be a string containing the CA certificates in PEM format",
				},
			)
		} else {
			tls.CACertsFromMemory([]byte(caCerts))
		}
	}

	if len(diags) != 0 {
		return nil, diags
	}

	client, err := ovirtclient.New(
		url,
		username,
		password,
		tls,
		&terraformLogger{
			ctx: ctx,
		},
		nil,
	)
	if err != nil {
		diags = append(
			diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Failed to create oVirt client",
				Detail:        err.Error(),
				AttributePath: nil,
			},
		)
		return nil, diags
	}
	p.client = client
	return p, diags
}

func getStringSliceFromResource(fieldName string, data *schema.ResourceData, diags diag.Diagnostics) ([]string, diag.Diagnostics) {
	value, ok := data.GetOk(fieldName)
	if !ok {
		return nil, diags
	}

	values, ok := value.([]interface{})
	if !ok {
		diags = append(
			diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("The %s option must be a list", fieldName),
				Detail:   fmt.Sprintf("The %s option must be a list, but got %v", fieldName, value),
			},
		)
		return nil, diags
	}

	result := []string{}
	for _, val := range values {
		valueStr, ok := val.(string)
		if !ok {
			diags = append(
				diags, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       fmt.Sprintf("The %s option must be a list of strings", fieldName),
					Detail:        fmt.Sprintf("The %s option must be a list of strings. Value %v is not a string", fieldName, val),
					AttributePath: nil,
				},
			)
		} else {
			result = append(result, valueStr)
		}
	}
	return result, diags
}
