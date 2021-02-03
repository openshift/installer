// Package equinixmetal collects equinixmetal-specific configuration.
package equinixmetal

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/equinixmetal"
	"github.com/openshift/installer/pkg/validate"
)

const (
	DefaultFacility = "EWR1" // Parsippany, NJ, US
)

// Platform collects equinixmetal-specific configuration.
func Platform() (*equinixmetal.Platform, error) {
	facilityCode, err := selectFacility()
	if err != nil {
		return nil, err
	}

	projectID, err := selectProject()
	if err != nil {
		return nil, err
	}

	return &equinixmetal.Platform{
		Facility:  facilityCode,
		ProjectID: projectID,
	}, nil
}

func selectProject() (string, error) {
	var projectID string
	/*
		//TODO(displague) offer a mapping of project names to project ids, project
		// names may be duplicated between organizations so names may not be unique.

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		conn, err := getConnection(Config{APIKey: apiKey, APIURL: apiURL})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create Equinix Metal connection")
		}


		client := &Client{Conn: conn}
		projects, err := client.ListProjects(ctx)

		if err != nil {
			return nil, errors.Wrap(err, "failed to list Equinix Metal projects")
		}

		projectNames := []string{}
		for _, p := range projects {
			projectNames = append(projectNames, p.Name)
		}
	*/

	err := survey.Ask([]*survey.Question{{
		Prompt: &survey.Input{
			Message: "Equinix Metal Project ID",
			Help:    "The Equinix Metal project id to use for installation",
		},
	}}, &projectID)

	if err != nil {
		return "", err
	}

	return projectID, nil
}

func selectFacility() (string, error) {
	var facilityID string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Equinix Metal Facility Code",
				Help:    "The Equinix Metal Facility code (this is the short name, e.g. 'ewr1')",
				Default: DefaultFacility,
			},
		},
	}, &facilityID)
	if err != nil {
		return "", err
	}
	return facilityID, nil
}

func askForConfig() (*Config, error) {
	var apiURL, apiKey string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Equinix Metal API URL",
				Help:    "The base URL for accessing the Equinix Metal API",
				Default: "https://api.equinix.com/metal/v1",
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, apiURL)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Equinix Metal API Key",
				Help:    "The User or Project Equinix Metal API Key to access the Equinix Metal API",
			},
		},
	}, &apiKey)
	if err != nil {
		return nil, err
	}

	return &Config{
		APIKey: apiKey,
		APIURL: apiURL,
	}, nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}
