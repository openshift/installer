// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisEdgeFunctionsActions                         = "cis_edge_functions_actions"
	cisEdgeFunctionsActionEtag                      = "etag"
	cisEdgeFunctionsActionHandlers                  = "handlers"
	cisEdgeFunctionsActionRoutes                    = "routes"
	cisEdgeFunctionsActionTriggerID                 = "trigger_id"
	cisEdgeFunctionsActionRoutePattern              = "pattern_url"
	cisEdgeFunctionsActionRouteActionName           = "action_name"
	cisEdgeFunctionsActionRouteRequestLimitFailOpen = "request_limit_fail_open"
	cisEdgeFunctionsActionCreatedOn                 = "created_on"
	cisEdgeFunctionsActionModifiedOn                = "modified_on"
)

func dataSourceIBMCISEdgeFunctionsActions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISEdgeFunctionsActionsRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDataDiff,
			},
			cisEdgeFunctionsActions: {
				Type:        schema.TypeList,
				Description: "List of edge functions actions",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisEdgeFunctionsActionEtag: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function action etag",
						},
						cisEdgeFunctionsActionHandlers: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Edge function action handlers",
						},
						cisEdgeFunctionsActionCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function action script created on",
						},
						cisEdgeFunctionsActionModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function action script modified on",
						},
						cisEdgeFunctionsActionRoutes: {
							Type:        schema.TypeList,
							Description: "List of edge function action routes",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisEdgeFunctionsActionTriggerID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Edge function action script identifier",
									},
									cisEdgeFunctionsActionRouteActionName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Edge function action script name",
									},
									cisEdgeFunctionsActionRoutePattern: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Edge function action pattern",
									},
									cisEdgeFunctionsActionRouteRequestLimitFailOpen: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Edge function action script request limit fail open",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISEdgeFunctionsActionsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewListEdgeFunctionsActionsOptions()
	result, _, err := cisClient.ListEdgeFunctionsActions(opt)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	scriptInfo := make([]map[string]interface{}, 0)
	for _, script := range result.Result {
		routes := make([]map[string]interface{}, 0)
		for _, route := range script.Routes {
			r := map[string]interface{}{
				cisEdgeFunctionsActionTriggerID:                 *route.ID,
				cisEdgeFunctionsActionRoutePattern:              *route.Pattern,
				cisEdgeFunctionsActionRouteActionName:           *route.Script,
				cisEdgeFunctionsActionRouteRequestLimitFailOpen: *route.RequestLimitFailOpen,
			}
			routes = append(routes, r)
		}
		handlers := make([]string, 0)
		for _, h := range script.Handlers {
			handlers = append(handlers, h)
		}
		l := map[string]interface{}{
			cisEdgeFunctionsActionEtag:       *script.Etag,
			cisEdgeFunctionsActionHandlers:   handlers,
			cisEdgeFunctionsActionCreatedOn:  (*script.CreatedOn).String(),
			cisEdgeFunctionsActionModifiedOn: (*script.ModifiedOn).String(),
			cisEdgeFunctionsActionRoutes:     routes,
		}
		scriptInfo = append(scriptInfo, l)
	}
	d.SetId(dataSourceIBMCISEdgeFunctionsActionsID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisEdgeFunctionsActions, scriptInfo)
	return nil
}

func dataSourceIBMCISEdgeFunctionsActionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
