package opsworks

import (
	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// @SDKResource("aws_opsworks_mysql_layer", name="MySQL Layer")
// @Tags(identifierAttribute="arn")
func ResourceMySQLLayer() *schema.Resource {
	layerType := &opsworksLayerType{
		TypeName:         opsworks.LayerTypeDbMaster,
		DefaultLayerName: "MySQL",

		Attributes: map[string]*opsworksLayerTypeAttribute{
			"root_password": {
				AttrName:  opsworks.LayerAttributesKeysMysqlRootPassword,
				Type:      schema.TypeString,
				WriteOnly: true,
			},
			"root_password_on_all_instances": {
				AttrName: opsworks.LayerAttributesKeysMysqlRootPasswordUbiquitous,
				Type:     schema.TypeBool,
				Default:  true,
			},
		},
	}

	return layerType.resourceSchema()
}
