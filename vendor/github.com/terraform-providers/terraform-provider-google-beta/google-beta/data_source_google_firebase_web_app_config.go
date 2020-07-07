package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleFirebaseWebappConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleFirebaseWebappConfigRead,

		Schema: map[string]*schema.Schema{
			"web_app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The id of the Firebase web App.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The project id of the Firebase web App.`,
			},
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API key associated with the web App.`,
			},
			"auth_domain": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The domain Firebase Auth configures for OAuth redirects, in the format:

projectId.firebaseapp.com`,
			},
			"database_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default Firebase Realtime Database URL.`,
			},
			"location_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The ID of the project's default GCP resource location. The location is one of the available GCP resource
locations.

This field is omitted if the default GCP resource location has not been finalized yet. To set your project's
default GCP resource location, call defaultLocation.finalize after you add Firebase services to your project.`,
			},
			"measurement_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The unique Google-assigned identifier of the Google Analytics web stream associated with the Firebase Web App.
Firebase SDKs use this ID to interact with Google Analytics APIs.

This field is only present if the App is linked to a web stream in a Google Analytics App + Web property.
Learn more about this ID and Google Analytics web streams in the Analytics documentation.

To generate a measurementId and link the Web App with a Google Analytics web stream,
call projects.addGoogleAnalytics.`,
			},
			"messaging_sender_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sender ID for use with Firebase Cloud Messaging.`,
			},
			"storage_bucket": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default Cloud Storage for Firebase storage bucket name.`,
			},
		},
	}

}

func dataSourceGoogleFirebaseWebappConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Get("web_app_id").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{FirebaseBasePath}}projects/{{project}}/webApps/{{web_app_id}}/config")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("FirebaseWebApp config %q", d.Id()))
	}

	err = d.Set("api_key", res["apiKey"])
	if err != nil {
		return err
	}
	err = d.Set("auth_domain", res["authDomain"])
	if err != nil {
		return err
	}
	err = d.Set("database_url", res["databaseURL"])
	if err != nil {
		return err
	}
	err = d.Set("location_id", res["locationId"])
	if err != nil {
		return err
	}
	err = d.Set("measurement_id", res["measurementId"])
	if err != nil {
		return err
	}
	err = d.Set("messaging_sender_id", res["messagingSenderId"])
	if err != nil {
		return err
	}
	err = d.Set("storage_bucket", res["storageBucket"])
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}
