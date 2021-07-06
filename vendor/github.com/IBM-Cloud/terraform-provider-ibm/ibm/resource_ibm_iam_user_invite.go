// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// MEMBER ...
	MEMBER = "MEMEBER"
	// ACCESS ...
	ACCESS          = "access"
	NOACCESS        = "noacess"
	VIEWONLY        = "viewonly"
	BASICUSER       = "basicuser"
	SUPERUSER       = "superuser"
	MANAGER         = "manager"
	AUDITOR         = "auditor"
	BILLINGMANANGER = "billingmanager"
	DEVELOPER       = "developer"
)

var viewOnly = []string{
	"HARDWARE_VIEW",
	"BANDWIDTH_MANAGE",
	"LICENSE_VIEW",
	"CDN_BANDWIDTH_VIEW",
	"VIRTUAL_GUEST_VIEW",
	"DEDICATED_HOST_VIEW",
}

var noAccess = make([]string, 0)

var basicUser = []string{"HARDWARE_VIEW",
	"USER_MANAGE",
	"BANDWIDTH_MANAGE",
	"DNS_MANAGE",
	"REMOTE_MANAGEMENT",
	"MONITORING_MANAGE",
	"LICENSE_VIEW",
	"IP_ADD",
	"PORT_CONTROL",
	"LOADBALANCER_MANAGE",
	"FIREWALL_MANAGE",
	"SOFTWARE_FIREWALL_MANAGE",
	"ANTI_MALWARE_MANAGE",
	"HOST_ID_MANAGE",
	"VULN_SCAN_MANAGE",
	"NTF_SUBSCRIBER_MANAGE",
	"CDN_BANDWIDTH_VIEW",
	"VIRTUAL_GUEST_VIEW",
	"NETWORK_MESSAGE_DELIVERY_MANAGE",
	"FIREWALL_RULE_MANAGE",
	"DEDICATED_HOST_VIEW",
}

var superUser = []string{"HARDWARE_VIEW",
	"VIEW_CUSTOMER_SOFTWARE_PASSWORD",
	"NETWORK_TUNNEL_MANAGE",
	"CUSTOMER_POST_PROVISION_SCRIPT_MANAGEMENT",
	"VIEW_CPANEL",
	"VIEW_PLESK",
	"VIEW_HELM",
	"VIEW_URCHIN",
	"ADD_SERVICE_STORAGE",
	"USER_MANAGE",
	"SERVER_ADD",
	"SERVER_UPGRADE",
	"SERVER_CANCEL",
	"SERVICE_ADD",
	"SERVICE_UPGRADE",
	"SERVICE_CANCEL",
	"BANDWIDTH_MANAGE",
	"DNS_MANAGE",
	"REMOTE_MANAGEMENT",
	"MONITORING_MANAGE",
	"SERVER_RELOAD",
	"LICENSE_VIEW",
	"IP_ADD",
	"LOCKBOX_MANAGE",
	"NAS_MANAGE",
	"PORT_CONTROL",
	"LOADBALANCER_MANAGE",
	"FIREWALL_MANAGE",
	"SOFTWARE_FIREWALL_MANAGE",
	"ANTI_MALWARE_MANAGE",
	"HOST_ID_MANAGE",
	"VULN_SCAN_MANAGE",
	"NTF_SUBSCRIBER_MANAGE",
	"NETWORK_VLAN_SPANNING",
	"CDN_ACCOUNT_MANAGE",
	"CDN_FILE_MANAGE",
	"CDN_BANDWIDTH_VIEW",
	"NETWORK_ROUTE_MANAGE",
	"VIRTUAL_GUEST_VIEW",
	"INSTANCE_UPGRADE",
	"HOSTNAME_EDIT",
	"NETWORK_MESSAGE_DELIVERY_MANAGE",
	"USER_EVENT_LOG_VIEW",
	"VPN_MANAGE",
	"VIEW_QUANTASTOR",
	"DATACENTER_ACCESS",
	"DATACENTER_ROOM_ACCESS",
	"CUSTOMER_SSH_KEY_MANAGEMENT",
	"FIREWALL_RULE_MANAGE",
	"PUBLIC_IMAGE_MANAGE",
	"SECURITY_CERTIFICATE_VIEW",
	"SECURITY_CERTIFICATE_MANAGE",
	"GATEWAY_MANAGE",
	"SCALE_GROUP_MANAGE",
	"SAML_AUTHENTICATION_MANAGE",
	"MANAGE_SECURITY_GROUPS",
	"PUBLIC_NETWORK_COMPUTE",
	"DEDICATED_HOST_VIEW",
}

var permissionSets = map[string][]string{NOACCESS: noAccess, VIEWONLY: viewOnly,
	BASICUSER: basicUser, SUPERUSER: superUser}

func resourceIBMUserInvite() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMInviteUsers,
		Read:     resourceIBMIAMGetUsers,
		Update:   resourceIBMIAMUpdateUserProfile,
		Delete:   resourceIBMIAMRemoveUser,
		Exists:   resourceIBMIAMGetUserProfileExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{

			"users": {
				Description: "List of ibm id or email of user",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"access_groups": {
				Description: "access group ids to associate the inviting user",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"iam_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"roles": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Role names of the policy definition",
						},

						"resources": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Service name of the policy definition",
									},

									"resource_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of resource instance of the policy definition",
									},

									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region of the policy definition",
									},

									"resource_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource type of the policy definition",
									},

									"resource": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource of the policy definition",
									},

									"resource_group_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the resource group.",
									},

									"attributes": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Set resource attributes in the form of 'name=value,name=value....",
										Elem:        schema.TypeString,
									},
								},
							},
						},
						"account_management": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Give access to all account management services",
						},
					},
				},
			},
			"number_of_invited_users": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of users invited to an account",
			},
			"invited_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"user_id": {
							Description: "ibm id or email of user",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"user_policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"roles": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Role names of the policy definition",
									},

									"resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Service name of the policy definition",
												},

												"resource_instance_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of resource instance of the policy definition",
												},

												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region of the policy definition",
												},

												"resource_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource type of the policy definition",
												},

												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource of the policy definition",
												},

												"resource_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the resource group.",
												},

												"attributes": {
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "Set resource attributes in the form of 'name=value,name=value....",
													Elem:        schema.TypeString,
												},
											},
										},
									},
								},
							},
						},

						"access_groups": {
							Description: "access group ids to associate the inviting user",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": {
										Description: "Name of the access group",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"policies": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},

												"roles": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Role names of the policy definition",
												},

												"resources": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"service": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Service name of the policy definition",
															},

															"resource_instance_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ID of resource instance of the policy definition",
															},

															"region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Region of the policy definition",
															},

															"resource_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource type of the policy definition",
															},

															"resource": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource of the policy definition",
															},

															"resource_group_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ID of the resource group.",
															},

															"attributes": {
																Type:        schema.TypeMap,
																Computed:    true,
																Description: "Set resource attributes in the form of 'name=value,name=value....",
																Elem:        schema.TypeString,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"classic_infra_roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_set": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "permission set for claasic infrastructure",
							ValidateFunc: validateAllowedStringValue([]string{NOACCESS, VIEWONLY, BASICUSER, SUPERUSER}),
						},

						"permissions": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of permissions for claasic infrastructure",
						},
					},
				},
			},
			"cloud_foundry_roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"organization_guid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "GUID of Organization",
						},

						"org_roles": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "roles to be assigned to user in given space",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},

						"spaces": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"space_guid": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "GUID of space",
									},

									"space_roles": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "roles to be assigned to user in given space",
										Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceIBMIAMInviteUsers(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	usersSet := d.Get("users").(*schema.Set)
	usersList := flattenUsersSet(usersSet)
	users := make([]v2.User, 0)
	for _, user := range usersList {
		users = append(users, v2.User{Email: user, AccountRole: MEMBER})
	}
	if len(users) == 0 {
		return fmt.Errorf("Users email not provided")
	}
	var accessGroups = make([]string, 0)
	if data, ok := d.GetOk("access_groups"); ok {
		for _, accessGroup := range data.([]interface{}) {
			accessGroups = append(accessGroups, fmt.Sprintf("%v", accessGroup))
		}
	}

	var accessPolicies []v2.UserPolicy
	if accessPolicyData, ok := d.GetOk("iam_policy"); ok {
		accessPolicies, err = getPolicies(d, meta, accessPolicyData.([]interface{}))
		if err != nil {
			log.Println("IAM Acess policy: ", err.Error())
			return err
		}
	}

	inviteUserPayload := v2.UserInvite{}
	log.Println(inviteUserPayload)
	inviteUserPayload.Users = users
	if len(accessGroups) != 0 {
		inviteUserPayload.AccessGroup = accessGroups
	}
	if len(accessPolicies) != 0 {
		inviteUserPayload.IAMPolicy = accessPolicies
	}

	if infraPermissions := getInfraPermissions(d, meta); len(infraPermissions) != 0 {
		inviteUserPayload.InfrastructureRoles = &v2.InfraPermissions{Permissions: infraPermissions}
	}
	orgRoles, err := getCloudFoundryRoles(d, meta)
	if err != nil {
		return err
	}
	if len(orgRoles) != 0 {
		inviteUserPayload.OrganizationRoles = orgRoles
	}

	accountID, err := getAccountID(d, meta)
	if err != nil {
		return err
	}

	_, InviteUserError := client.InviteUsers(accountID, inviteUserPayload)
	if InviteUserError != nil {
		return InviteUserError
	}
	d.SetId(time.Now().UTC().String())
	return resourceIBMIAMUpdateUserProfile(d, meta)
}

func resourceIBMIAMGetUsers(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	Client := userManagement.UserInvite()
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}
	accountID, err := getAccountID(d, meta)
	if err != nil {
		return err
	}
	res, err := Client.ListUsers(accountID)
	if err != nil {
		return err
	}
	users := make([]string, 0)
	invitedUsers := make([]map[string]interface{}, 0, len(res))

	for _, user := range res {

		if user.AccountID != accountID {
			users = append(users, user.Email)
		}
		/****** For each user *******************
		    1) user_id
		    2) user_level_policies
		    3) List of access groups
		            > Name of access group
		            > acees group level policies
		********************************************/
		//Get User level IAM policies
		policyList, _, err := iamPolicyManagementClient.ListPolicies(&iampolicymanagementv1.ListPoliciesOptions{
			AccountID: core.StringPtr(accountID),
			IamID:     core.StringPtr(user.IamID),
			Type:      core.StringPtr("access"),
		})
		policies := policyList.Policies

		if err != nil {
			return fmt.Errorf("Error retrieving user policies: %s", err)
		}
		userPolicies := make([]map[string]interface{}, 0, len(policies))
		for _, policy := range policies {
			//populate ploicy Roles
			roles := make([]string, len(policy.Roles))
			for i, role := range policy.Roles {
				roles[i] = *role.DisplayName
			}
			//populate policy resources
			resources := flattenPolicyResource(policy.Resources)
			p := map[string]interface{}{
				"id":        policy.ID,
				"roles":     roles,
				"resources": resources,
			}
			userPolicies = append(userPolicies, p)
		}

		// Get AccessGroups associated with user
		retreivedGroups, err := iamuumClient.AccessGroup().List(accountID, user.IamID)
		if err != nil {
			return fmt.Errorf("Error retrieving access groups: %s", err)
		}

		accGroupList := make([]map[string]interface{}, 0, len(retreivedGroups))
		//Get the policies for each access group
		for _, grpData := range retreivedGroups {
			policyList, _, err := iamPolicyManagementClient.ListPolicies(&iampolicymanagementv1.ListPoliciesOptions{
				AccountID:     core.StringPtr(accountID),
				AccessGroupID: core.StringPtr(user.IamID),
			})
			accgrpPolicy := policyList.Policies
			if err != nil {
				return fmt.Errorf("Error retrieving access group policy: %s", err)
			}

			//Fetch access group policies
			grpPolicies := make([]map[string]interface{}, 0, len(accgrpPolicy))
			for _, policy := range accgrpPolicy {
				//populate ploicy Roles
				roles := make([]string, len(policy.Roles))
				for i, role := range policy.Roles {
					roles[i] = *role.DisplayName
				}
				//populate policy resources
				resources := flattenPolicyResource(policy.Resources)
				p := map[string]interface{}{
					"id":        policy.ID,
					"roles":     roles,
					"resources": resources,
				}
				grpPolicies = append(grpPolicies, p)
			}
			//populate name & policies of a access group
			agInfo := map[string]interface{}{
				"name":     grpData.Name,
				"policies": grpPolicies,
			}
			//add agInfo to list of access groups
			accGroupList = append(accGroupList, agInfo)
		}
		userInfo := map[string]interface{}{
			"user_id":       user.Email,
			"user_policies": userPolicies,
			"access_groups": accGroupList,
		}
		invitedUsers = append(invitedUsers, userInfo)
	}
	//set the number of users in an account
	d.Set("number_of_invited_users", len(res)-1)
	d.Set("invited_users", invitedUsers)
	return nil
}

func resourceIBMIAMUpdateUserProfile(d *schema.ResourceData, meta interface{}) error {
	// validate change
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	Client := userManagement.UserInvite()

	if d.HasChange("users") {
		//var removedUsers, addedUsers []string
		accountID, err := getAccountID(d, meta)
		if err != nil {
			return err
		}
		ousrs, nusrs := d.GetChange("users")
		old := ousrs.(*schema.Set)
		new := nusrs.(*schema.Set)

		removed := expandStringList(old.Difference(new).List())
		added := expandStringList(new.Difference(old).List())

		//Update the added users
		if len(added) > 0 {
			users := make([]v2.User, 0)
			for _, user := range added {
				users = append(users, v2.User{Email: user, AccountRole: MEMBER})
			}
			if len(users) == 0 {
				return fmt.Errorf("Users email not provided")
			}

			var accessPolicies []v2.UserPolicy
			if accessPolicyData, ok := d.GetOk("iam_policy"); ok {
				accessPolicies, err = getPolicies(d, meta, accessPolicyData.([]interface{}))
				if err != nil {
					log.Println("IAM Acess policy: ", err.Error())
					return err
				}
			}

			var accessGroups = make([]string, 0)
			if data, ok := d.GetOk("access_groups"); ok {
				for _, accessGroup := range data.([]interface{}) {
					accessGroups = append(accessGroups, fmt.Sprintf("%v", accessGroup))
				}
			}

			infraPermissions := getInfraPermissions(d, meta)
			orgRoles, err := getCloudFoundryRoles(d, meta)
			if err != nil {
				return err
			}

			inviteUserPayload := v2.UserInvite{}

			inviteUserPayload.Users = users
			if len(accessGroups) != 0 {
				inviteUserPayload.AccessGroup = accessGroups
			}
			if len(accessPolicies) != 0 {
				inviteUserPayload.IAMPolicy = accessPolicies
			}
			if len(infraPermissions) != 0 {
				inviteUserPayload.InfrastructureRoles = &v2.InfraPermissions{Permissions: infraPermissions}
			}
			if len(orgRoles) != 0 {
				inviteUserPayload.OrganizationRoles = orgRoles
			}
			_, InviteUserError := Client.InviteUsers(accountID, inviteUserPayload)
			if InviteUserError != nil {
				return InviteUserError
			}
		}

		//Update the removed users
		if len(removed) > 0 {
			for _, user := range removed {
				IAMID, err := getUserIAMID(d, meta, user)
				if err != nil {
					return fmt.Errorf("User's IAM ID not found: %s", err.Error())
				}
				Err := Client.RemoveUsers(accountID, IAMID)
				if Err != nil {
					log.Println("Failed to remove user: ", user)
					return Err
				}
			}
		}

	}
	return resourceIBMIAMGetUsers(d, meta)
}

func resourceIBMIAMRemoveUser(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	Client := userManagement.UserInvite()

	accountID, err := getAccountID(d, meta)
	if err != nil {
		return err
	}

	usersSet := d.Get("users").(*schema.Set)
	usersList := flattenUsersSet(usersSet)
	for _, user := range usersList {
		IAMID, err := getUserIAMID(d, meta, user)

		if err != nil {
			return fmt.Errorf("User's IAM ID not found: %s", err.Error())
		}
		Err := Client.RemoveUsers(accountID, IAMID)
		if Err != nil {
			return Err
		}
	}
	return nil
}

func resourceIBMIAMGetUserProfileExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return false, err
	}
	Client := userManagement.UserInvite()

	accountID, err := getAccountID(d, meta)
	if err != nil {
		return false, err
	}

	usersSet := d.Get("users").(*schema.Set)
	usersList := flattenUsersSet(usersSet)

	res, err := Client.ListUsers(accountID)
	if err != nil {
		return false, err
	}
	var isFound bool
	for _, user := range usersList {

		for _, userInfo := range res {
			if strings.Compare(userInfo.Email, user) == 0 {
				isFound = true
			}
		}
		if !isFound {
			return false, nil
		}
	}
	return true, nil
}

// getAccountID returns accountID
func getAccountID(d *schema.ResourceData, meta interface{}) (string, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return "", err
	}
	return userDetails.userAccount, nil
}

// getUserIAMID ...
func getUserIAMID(d *schema.ResourceData, meta interface{}, user string) (string, error) {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return "", err
	}
	Client := userManagement.UserInvite()

	accountID, err := getAccountID(d, meta)
	if err != nil {
		return "", err
	}

	res, err := Client.ListUsers(accountID)
	if err != nil {
		return "", err
	}

	for _, userInfo := range res {
		if strings.Compare(userInfo.Email, user) == 0 {
			return userInfo.IamID, nil
		}
	}
	return "", nil

}

func getInfraPermissions(d *schema.ResourceData, meta interface{}) []string {
	var infraPermissions = make([]string, 0)
	if data, ok := d.GetOk("classic_infra_roles"); ok {
		for _, resource := range data.([]interface{}) {
			d := resource.(map[string]interface{})
			if permissions, ok := d["permissions"]; ok && permissions != nil {
				for _, value := range permissions.([]interface{}) {
					infraPermissions = append(infraPermissions, fmt.Sprintf("%v", value))
				}
			}
			if permissionSet, ok := d["permission_set"]; ok && permissionSet != nil {
				if permissions, ok := permissionSets[permissionSet.(string)]; ok {
					for _, permission := range permissions {
						infraPermissions = append(infraPermissions, permission)
					}
				}
			}
		}
		return infraPermissions
	}
	return infraPermissions
}

// getPolicies ...
func getPolicies(d *schema.ResourceData, meta interface{}, policies []interface{}) ([]v2.UserPolicy, error) {
	var policyList = make([]v2.UserPolicy, 0)
	for _, policy := range policies {
		p := policy.(map[string]interface{})
		var serviceName string
		resourceAttributes := []iampolicymanagementv1.ResourceAttribute{}
		policyResource := iampolicymanagementv1.PolicyResource{}

		if res, ok := p["resources"]; ok {
			resources := res.([]interface{})
			for _, resource := range resources {
				r, _ := resource.(map[string]interface{})
				serviceName = r["service"].(string)
				if r, ok := r["service"]; ok && r != nil {
					serviceName = r.(string)
					if r.(string) != "" {
						resourceAttr := iampolicymanagementv1.ResourceAttribute{
							Name:     core.StringPtr("serviceName"),
							Value:    core.StringPtr(r.(string)),
							Operator: core.StringPtr("stringEquals"),
						}
						resourceAttributes = append(resourceAttributes, resourceAttr)
					}
				}

				if r, ok := r["resource_instance_id"]; ok {
					if r.(string) != "" {
						resourceAttr := iampolicymanagementv1.ResourceAttribute{
							Name:     core.StringPtr("serviceInstance"),
							Value:    core.StringPtr(r.(string)),
							Operator: core.StringPtr("stringEquals"),
						}
						resourceAttributes = append(resourceAttributes, resourceAttr)
					}
				}

				if r, ok := r["region"]; ok {
					if r.(string) != "" {
						resourceAttr := iampolicymanagementv1.ResourceAttribute{
							Name:     core.StringPtr("region"),
							Value:    core.StringPtr(r.(string)),
							Operator: core.StringPtr("stringEquals"),
						}
						resourceAttributes = append(resourceAttributes, resourceAttr)
					}
				}

				if r, ok := r["resource_type"]; ok {
					if r.(string) != "" {
						resourceAttr := iampolicymanagementv1.ResourceAttribute{
							Name:     core.StringPtr("resourceType"),
							Value:    core.StringPtr(r.(string)),
							Operator: core.StringPtr("stringEquals"),
						}
						resourceAttributes = append(resourceAttributes, resourceAttr)
					}
				}

				if r, ok := r["resource"]; ok {
					if r.(string) != "" {
						resourceAttr := iampolicymanagementv1.ResourceAttribute{
							Name:     core.StringPtr("resource"),
							Value:    core.StringPtr(r.(string)),
							Operator: core.StringPtr("stringEquals"),
						}
						resourceAttributes = append(resourceAttributes, resourceAttr)
					}
				}

				if r, ok := r["resource_group_id"]; ok {
					if r.(string) != "" {
						resourceAttr := iampolicymanagementv1.ResourceAttribute{
							Name:     core.StringPtr("resourceGroupId"),
							Value:    core.StringPtr(r.(string)),
							Operator: core.StringPtr("stringEquals"),
						}
						resourceAttributes = append(resourceAttributes, resourceAttr)
					}
				}

				if r, ok := r["attributes"]; ok {
					for k, v := range r.(map[string]interface{}) {
						resourceAttributes = setResourceAttribute(core.StringPtr(k), v.(*string), resourceAttributes)
					}
				}

			}
		}

		if accountManagement, ok := p["account_management"]; ok && accountManagement.(bool) {
			serviceTypeResourceAttribute := iampolicymanagementv1.ResourceAttribute{
				Name:     core.StringPtr("serviceType"),
				Value:    core.StringPtr("platform_service"),
				Operator: core.StringPtr("stringEquals"),
			}
			resourceAttributes = append(resourceAttributes, serviceTypeResourceAttribute)
		}

		if len(resourceAttributes) == 0 {
			serviceTypeResourceAttribute := iampolicymanagementv1.ResourceAttribute{
				Name:     core.StringPtr("serviceType"),
				Value:    core.StringPtr("service"),
				Operator: core.StringPtr("stringEquals"),
			}
			resourceAttributes = append(resourceAttributes, serviceTypeResourceAttribute)
		}

		accountID, err := getAccountID(d, meta)
		if err != nil {
			return policyList, err
		}

		accountIDResourceAttribute := iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("accountId"),
			Value:    core.StringPtr(accountID),
			Operator: core.StringPtr("stringEquals"),
		}

		resourceAttributes = append(resourceAttributes, accountIDResourceAttribute)

		policyResource.Attributes = resourceAttributes

		iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
		if err != nil {
			return policyList, err
		}

		userDetails, err := meta.(ClientSession).BluemixUserDetails()
		if err != nil {
			return policyList, err
		}

		listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
			AccountID:   &userDetails.userAccount,
			ServiceName: &serviceName,
		}

		roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
		roles := mapRoleListToPolicyRoles(*roleList)

		if err != nil {
			return policyList, err
		}
		var policyRoles = make([]iampolicymanagementv1.PolicyRole, 0)
		if userRoles, ok := p["roles"]; ok {
			policyRoles, err = getRolesFromRoleNames(expandStringList(userRoles.([]interface{})), roles)
			if err != nil {
				return policyList, err
			}
		}

		policyList = append(policyList, v2.UserPolicy{Roles: convertIPMRolesToV1(policyRoles), Resources: convertIPMResourcesToV1(policyResource), Type: ACCESS})
	}
	return policyList, nil
}

func convertIPMRolesToV1(roles []iampolicymanagementv1.PolicyRole) []iampapv1.Role {
	results := make([]iampapv1.Role, len(roles))
	for i, r := range roles {
		results[i] = iampapv1.Role{
			RoleID: *r.RoleID,
		}
	}
	return results
}

func convertIPMResourcesToV1(resource iampolicymanagementv1.PolicyResource) []iampapv1.Resource {
	attributes := make([]iampapv1.Attribute, len(resource.Attributes))
	for i, a := range resource.Attributes {
		attributes[i] = iampapv1.Attribute{
			Name:  *a.Name,
			Value: *a.Value,
		}
	}
	return []iampapv1.Resource{iampapv1.Resource{Attributes: attributes}}
}

// getCloudFoundryRoles ...
func getCloudFoundryRoles(d *schema.ResourceData, meta interface{}) ([]v2.OrgRole, error) {
	cloudFoundryRoles := make([]v2.OrgRole, 0)
	if data, ok := d.GetOk("cloud_foundry_roles"); ok {
		sess, err := meta.(ClientSession).BluemixSession()
		if err != nil {
			return nil, err
		}
		usersSet := d.Get("users").(*schema.Set)
		usersList := flattenUsersSet(usersSet)
		for _, d := range data.([]interface{}) {
			orgRole := v2.OrgRole{}
			role := d.(map[string]interface{})
			orgRole.ID = role["organization_guid"].(string)
			orgRole.Region = sess.Config.Region
			orgRole.Users = usersList
			for _, r := range role["org_roles"].([]interface{}) {
				switch strings.ToLower(r.(string)) {
				case AUDITOR:
					orgRole.Auditors = usersList
				case BILLINGMANANGER:
					orgRole.BillingManagers = usersList
				case MANAGER:
					orgRole.Managers = usersList
				}
			}
			if spaces, ok := role["spaces"]; ok {
				for _, s := range spaces.([]interface{}) {
					spaceInfo := v2.Space{}
					space := s.(map[string]interface{})
					if spaceroles, ok := space["space_roles"]; ok {
						for _, r := range spaceroles.([]interface{}) {
							role := r.(string)
							switch strings.ToLower(role) {
							case AUDITOR:
								spaceInfo.Auditors = usersList
							case DEVELOPER:
								spaceInfo.Developers = usersList
							case MANAGER:
								spaceInfo.Managers = usersList
							}

						}
					}
					if spaceName, ok := space["space_guid"]; ok {
						spaceInfo.ID = spaceName.(string)
					}
					orgRole.Spaces = append(orgRole.Spaces, spaceInfo)
				}
			}
			cloudFoundryRoles = append(cloudFoundryRoles, orgRole)

		}
	}
	return cloudFoundryRoles, nil
}
