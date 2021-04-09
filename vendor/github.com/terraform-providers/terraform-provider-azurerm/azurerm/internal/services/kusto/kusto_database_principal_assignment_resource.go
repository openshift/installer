package kusto

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKustoDatabasePrincipalAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKustoDatabasePrincipalAssignmentCreate,
		Read:   resourceKustoDatabasePrincipalAssignmentRead,
		Delete: resourceKustoDatabasePrincipalAssignmentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoClusterName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoDatabaseName,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoDatabasePrincipalAssignmentName,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"principal_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"principal_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"principal_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.PrincipalTypeApp),
					string(kusto.PrincipalTypeGroup),
					string(kusto.PrincipalTypeUser),
				}, false),
			},

			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.Admin),
					string(kusto.Ingestor),
					string(kusto.Monitor),
					string(kusto.User),
					string(kusto.UnrestrictedViewers),
					string(kusto.Viewer),
				}, false),
			},
		},
	}
}

func resourceKustoDatabasePrincipalAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Database Principal Assignment creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)
	name := d.Get("name").(string)

	existing, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_kusto_database_principal_assignment", *existing.ID)
	}

	tenantID := d.Get("tenant_id").(string)
	principalID := d.Get("principal_id").(string)
	principalType := d.Get("principal_type").(string)
	role := d.Get("role").(string)

	principalAssignment := kusto.DatabasePrincipalAssignment{
		DatabasePrincipalProperties: &kusto.DatabasePrincipalProperties{
			TenantID:      utils.String(tenantID),
			PrincipalID:   utils.String(principalID),
			PrincipalType: kusto.PrincipalType(principalType),
			Role:          kusto.DatabasePrincipalRole(role),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, databaseName, name, principalAssignment)
	if err != nil {
		return fmt.Errorf("failed creating Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed waiting for completion of Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
	if err != nil {
		return fmt.Errorf("failed to retrieve Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read ID for Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q)", name, resourceGroup, clusterName, databaseName)
	}

	d.SetId(*resp.ID)

	return resourceKustoDatabasePrincipalAssignmentRead(d, meta)
}

func resourceKustoDatabasePrincipalAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabasePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.PrincipalAssignmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", id.PrincipalAssignmentName, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)
	d.Set("name", id.PrincipalAssignmentName)

	tenantID := ""
	if resp.TenantID != nil {
		tenantID = *resp.TenantID
	}

	tenantName := ""
	if resp.TenantName != nil {
		tenantName = *resp.TenantName
	}

	principalID := ""
	if resp.PrincipalID != nil {
		principalID = *resp.PrincipalID
	}

	principalName := ""
	if resp.PrincipalName != nil {
		principalName = *resp.PrincipalName
	}

	principalType := string(resp.PrincipalType)
	role := string(resp.Role)

	d.Set("tenant_id", tenantID)
	d.Set("tenant_name", tenantName)
	d.Set("principal_id", principalID)
	d.Set("principal_name", principalName)
	d.Set("principal_type", principalType)
	d.Set("role", role)

	return nil
}

func resourceKustoDatabasePrincipalAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasePrincipalAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabasePrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.PrincipalAssignmentName)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", id.PrincipalAssignmentName, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Database Principal Assignment %q (Resource Group %q, Cluster %q, Database %q): %+v", id.PrincipalAssignmentName, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	return nil
}

func validateAzureRMKustoDatabasePrincipalAssignmentName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 260 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}
