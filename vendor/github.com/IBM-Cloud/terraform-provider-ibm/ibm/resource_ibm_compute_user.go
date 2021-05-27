// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

const userCustomerCancelStatus = 1021

func resourceIBMComputeUser() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeUserCreate,
		Read:     resourceIBMComputeUserRead,
		Update:   resourceIBMComputeUserUpdate,
		Delete:   resourceIBMComputeUserDelete,
		Exists:   resourceIBMComputeUserExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "user name",
			},
			"first_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "First name of the user",
			},
			"last_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Last name of the user",
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "email address of the user",
			},
			"company_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "comapany name",
			},
			"address1": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address info of the user",
			},
			"address2": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address info of the user",
			},
			"city": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "City name",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Satate name",
			},
			"country": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Country name",
			},
			"timezone": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "time zone info",
			},
			"user_status": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ACTIVE",
				Description: "user status info",
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				StateFunc: func(v interface{}) string {
					hash := sha256.Sum256([]byte(v.(string)))
					return hex.EncodeToString(hash[:])
				},
				Description: "password for the user",
			},
			"permissions": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "set of persmissions assigned for the user",
			},
			"has_api_key": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "API Key info of the user",
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "API key for the user",
			},
			"ibm_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IBM ID of the  user",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Tags set for the resources",
			},
		},
	}
}

// Create a SoftLayer_User_Customer_CustomerPermission_Permission object from the given string input
func makePermission(p string) datatypes.User_Customer_CustomerPermission_Permission {
	return datatypes.User_Customer_CustomerPermission_Permission{
		KeyName: &p,
	}
}

// Convert a "set" of permission strings to a list of SoftLayer_User_Customer_CustomerPermission_Permissions
func getPermissions(d *schema.ResourceData) []datatypes.User_Customer_CustomerPermission_Permission {
	permissionsSet := d.Get("permissions").(*schema.Set)

	if permissionsSet.Len() == 0 {
		return nil
	}

	permissions := make([]datatypes.User_Customer_CustomerPermission_Permission, 0, permissionsSet.Len())
	for _, elem := range permissionsSet.List() {
		permission := makePermission(elem.(string))

		permissions = append(permissions, permission)
	}
	return permissions
}

func resourceIBMComputeUserCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetUserCustomerService(sess)
	serviceNoRetry := services.GetUserCustomerService(sess.SetRetries(0))

	timezoneID, err := getTimezoneIDByName(sess, d.Get("timezone").(string))
	if err != nil {
		return err
	}

	userStatusID, err := getUserStatusIDByName(sess, d.Get("user_status").(string))
	if err != nil {
		return err
	}

	// Build up our creation options
	opts := datatypes.User_Customer{
		FirstName:    sl.String(d.Get("first_name").(string)),
		LastName:     sl.String(d.Get("last_name").(string)),
		Email:        sl.String(d.Get("email").(string)),
		CompanyName:  sl.String(d.Get("company_name").(string)),
		Address1:     sl.String(d.Get("address1").(string)),
		City:         sl.String(d.Get("city").(string)),
		State:        sl.String(d.Get("state").(string)),
		Country:      sl.String(d.Get("country").(string)),
		TimezoneId:   &timezoneID,
		UserStatusId: &userStatusID,
	}

	if address2, ok := d.GetOk("address2"); ok {
		opts.Address2 = sl.String(address2.(string))
	}

	if username, ok := d.GetOk("username"); ok {
		opts.Username = sl.String(username.(string))
	}

	pass := sl.String(d.Get("password").(string))
	if *pass == "" {
		pass = nil
	}

	res, err := serviceNoRetry.CreateObject(&opts, pass, nil)

	if err != nil {
		return fmt.Errorf("Error creating IBM Cloud User: %s", err)
	}

	d.SetId(strconv.Itoa(*res.Id))
	log.Printf("[INFO] IBM Cloud User: %d", *res.Id)

	permissions := getPermissions(d)

	defaultPortalPermissions := []datatypes.User_Customer_CustomerPermission_Permission{
		{KeyName: sl.String("ACCESS_ALL_GUEST")},
		{KeyName: sl.String("ACCESS_ALL_HARDWARE")},
	}

	log.Printf("Replacing default portal permissions assigned by IBM Cloud with those specified in config")

	// Set the instance ID for the service to act on
	service = service.Id(*res.Id)

	_, err = service.RemoveBulkPortalPermission(defaultPortalPermissions, sl.Bool(true))
	if err != nil {
		return fmt.Errorf("Error removing default portal permissions for IBM Cloud User: %s", err)
	}

	_, err = service.AddBulkPortalPermission(permissions)
	if err != nil {
		return fmt.Errorf("Error setting portal permissions for IBM Cloud User: %s", err)
	}

	create_api_key_flag := d.Get("has_api_key").(bool)
	if create_api_key_flag {
		// We have to create the API key only if the flag is true. If 'false' we do not
		// take the delete action on the API key, as this is the create new user method,
		// and not the edit method.
		_, err = service.AddApiAuthenticationKey()
		if err != nil {
			return fmt.Errorf("Error creating API key: %s", err)
		}
	}

	return resourceIBMComputeUserRead(d, meta)
}

func resourceIBMComputeUserRead(d *schema.ResourceData, meta interface{}) error {
	service := services.GetUserCustomerService(meta.(ClientSession).SoftLayerSession())
	userID, _ := strconv.Atoi(d.Id())

	mask := strings.Join([]string{
		"id",
		"username",
		"email",
		"firstName",
		"lastName",
		"companyName",
		"address1",
		"address2",
		"city",
		"state",
		"country",
		"timezone.shortName",
		"userStatus.keyName",
		"permissions.keyName",
		"apiAuthenticationKeys.authenticationKey",
		"openIdConnectUserName",
	}, ";")

	sluserObj, err := service.Id(userID).Mask(mask).GetObject()
	if err != nil {
		// If the key is somehow already destroyed, mark as
		// successfully gone
		if strings.Contains(err.Error(), "404 Not Found") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving IBM Cloud User: %s", err)
	}

	d.Set("username", sluserObj.Username)
	d.Set("email", sluserObj.Email)
	d.Set("first_name", sluserObj.FirstName)
	d.Set("last_name", sluserObj.LastName)
	d.Set("company_name", sluserObj.CompanyName)
	d.Set("address1", sluserObj.Address1)
	d.Set("address2", sluserObj.Address2)
	d.Set("city", sluserObj.City)
	d.Set("state", sluserObj.State)
	d.Set("country", sluserObj.Country)
	d.Set("timezone", sluserObj.Timezone.ShortName)
	d.Set("user_status", sluserObj.UserStatus.KeyName)

	permissions := make([]string, 0, len(sluserObj.Permissions))
	for _, elem := range sluserObj.Permissions {
		permissions = append(permissions, *elem.KeyName)
	}
	d.Set("permissions", permissions)

	// If present, extract the api key from the SoftLayer response and set the field in the resource
	if len(sluserObj.ApiAuthenticationKeys) > 0 {
		d.Set("api_key", sluserObj.ApiAuthenticationKeys[0].AuthenticationKey) // as its a computed field
		d.Set("has_api_key", true)
	} else {
		d.Set("api_key", "")
		d.Set("has_api_key", false)
	}

	if sluserObj.OpenIdConnectUserName != nil {
		d.Set("ibm_id", sluserObj.OpenIdConnectUserName)
	}

	return nil
}

func resourceIBMComputeUserUpdate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetUserCustomerService(sess)
	serviceNoRetry := services.GetUserCustomerService(sess.SetRetries(0))

	sluid, _ := strconv.Atoi(d.Id())

	mask := strings.Join([]string{
		"id",
		"username",
		"email",
		"firstName",
		"lastName",
		"companyName",
		"address1",
		"address2",
		"city",
		"state",
		"country",
		"timezone.shortName",
		"userStatus.keyName",
		"permissions.keyName",
		"apiAuthenticationKeys.authenticationKey",
		"apiAuthenticationKeys.id",
	}, ";")

	service = service.Id(sluid)
	userObj, err := service.Mask(mask).GetObject()

	// Some fields cannot be updated such as username. Computed fields also cannot be updated
	// by explicitly providing a value. So only update the fields that are editable.
	// Password changes can also not be fully automated, and are not supported
	if d.HasChange("first_name") {
		userObj.FirstName = sl.String(d.Get("first_name").(string))
	}
	if d.HasChange("last_name") {
		userObj.LastName = sl.String(d.Get("last_name").(string))
	}
	if d.HasChange("email") {
		userObj.Email = sl.String(d.Get("email").(string))
	}
	if d.HasChange("company_name") {
		userObj.CompanyName = sl.String(d.Get("company_name").(string))
	}
	if d.HasChange("address1") {
		userObj.Address1 = sl.String(d.Get("address1").(string))
	}
	if d.HasChange("address2") {
		userObj.Address2 = sl.String(d.Get("address2").(string))
	}
	if d.HasChange("city") {
		userObj.City = sl.String(d.Get("city").(string))
	}
	if d.HasChange("state") {
		userObj.State = sl.String(d.Get("state").(string))
	}
	if d.HasChange("country") {
		userObj.Country = sl.String(d.Get("country").(string))
	}
	if d.HasChange("timezone") {
		tzID, err := getTimezoneIDByName(sess, d.Get("timezone").(string))
		if err != nil {
			return err
		}
		userObj.TimezoneId = &tzID
	}
	if d.HasChange("user_status") {
		userStatusID, err := getUserStatusIDByName(sess, d.Get("user_status").(string))
		if err != nil {
			return err
		}
		userObj.UserStatusId = &userStatusID
	}

	_, err = serviceNoRetry.Id(sluid).EditObject(&userObj)
	if err != nil {
		return fmt.Errorf("Error received while editing ibm_compute_user: %s", err)
	}

	if d.HasChange("permissions") {
		old, new := d.GetChange("permissions")

		// 1. Remove old permissions no longer appearing in the new set
		// 2. Add new permissions not already granted

		remove := old.(*schema.Set).Difference(new.(*schema.Set)).List()
		add := new.(*schema.Set).Difference(old.(*schema.Set)).List()

		oldPermissions := make([]datatypes.User_Customer_CustomerPermission_Permission, 0, len(remove))
		newPermissions := make([]datatypes.User_Customer_CustomerPermission_Permission, 0, len(add))

		for _, elem := range remove {
			oldPermissions = append(oldPermissions, makePermission(elem.(string)))
		}

		for _, elem := range add {
			newPermissions = append(newPermissions, makePermission(elem.(string)))
		}

		// 'remove' all old permissions
		_, err = service.RemoveBulkPortalPermission(oldPermissions, sl.Bool(true))
		if err != nil {
			return fmt.Errorf("Error received while removing old permissions from ibm_compute_user: %s", err)
		}

		// 'add' new permission set
		_, err = service.AddBulkPortalPermission(newPermissions)
		if err != nil {
			return fmt.Errorf("Error received while assigning new permissions to ibm_compute_user: %s", err)
		}
	}

	if d.HasChange("has_api_key") {
		// if true, then it means create an api key if none exists. Its a no-op if an api key already exists.
		// else false means, delete the api key if one exists. Its a no-op if no api key exists.
		api_key_flag := d.Get("has_api_key").(bool)

		// Get the current keys.
		keys := userObj.ApiAuthenticationKeys

		// Create a key if flag is true, and a key does not already exist.
		if api_key_flag {
			if len(keys) == 0 { // means key does not exist, so create one.
				key, err := service.AddApiAuthenticationKey()
				if err != nil {
					return fmt.Errorf("Error creating API key while editing ibm_compute_user resource: %s", err)
				}

				d.Set("api_key", key)
			} else {
				d.Set("api_key", keys[0].AuthenticationKey) // as api_key is a computed field
			}
		} else {
			// If false, then delete the key if there was one.
			if len(keys) > 0 {
				success, err := service.RemoveApiAuthenticationKey(keys[0].Id)
				if err != nil {
					return fmt.Errorf("Error deleting API key while editing ibm_compute_user resource: %s", err)
				}

				if !success {
					return fmt.Errorf(
						"The API reported removal of the api key was not successful for %s",
						d.Get("email").(string),
					)
				}
			}
			d.Set("api_key", nil)
		}
	}
	return nil
}

func resourceIBMComputeUserDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetUserCustomerService(sess)

	id, _ := strconv.Atoi(d.Id())

	user := datatypes.User_Customer{
		UserStatusId: sl.Int(userCustomerCancelStatus),
	}

	log.Printf("[INFO] Deleting IBM Cloud user: %d", id)
	_, err := service.Id(id).EditObject(&user)
	if err != nil {
		return fmt.Errorf("Error deleting IBM Cloud user: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMComputeUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetUserCustomerService(meta.(ClientSession).SoftLayerSession())

	id, err := strconv.Atoi(d.Id())

	result, err := service.Id(id).GetObject()

	return result.Id != nil && *result.Id == id && err == nil, nil
}

func getTimezoneIDByName(sess *session.Session, shortName string) (int, error) {
	zones, err := services.GetLocaleTimezoneService(sess).
		Mask("id,shortName").
		GetAllObjects()

	if err != nil {
		return -1, err
	}

	for _, zone := range zones {
		if *zone.ShortName == shortName {
			return *zone.Id, nil
		}
	}

	return -1, fmt.Errorf("Timezone %s could not be found", shortName)

}

func getUserStatusIDByName(sess *session.Session, name string) (int, error) {
	statuses, err := services.GetUserCustomerStatusService(sess).
		Mask("id,keyName").
		GetAllObjects()

	if err != nil {
		return -1, err
	}

	for _, status := range statuses {
		if *status.KeyName == name {
			return *status.Id, nil
		}
	}

	return -1, fmt.Errorf("User status %s could not be found", name)

}
