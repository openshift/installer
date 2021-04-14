// Package equinixmetal collects equinixmetal-specific configuration.
package equinixmetal

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/equinixmetal"
	"github.com/openshift/installer/pkg/validate"
)

const (
	DefaultFacility = "da11" // Dallas, TX, US
	DefaultMetro    = "SV"   // Silicon Valley, US
)

// Platform collects equinixmetal-specific configuration.
func Platform() (*equinixmetal.Platform, error) {
	facilityCode, err := selectFacility()
	if err != nil {
		return nil, err
	}

	metroCode, err := selectMetro()
	if err != nil {
		return nil, err
	}

	projectID, err := selectProject()
	if err != nil {
		return nil, err
	}

	return &equinixmetal.Platform{
		Facility:  facilityCode,
		Metro:     metroCode,
		ProjectID: projectID,
	}, nil
}

func selectProject() (string, error) {
	var projectID string

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
		projectNames = append(projectNames, p.ID+" "+p.Name)
	}

	err := survey.Ask([]*survey.Question{{
		Prompt: &survey.Input{
			Message: "Equinix Metal Project ID",
			Help:    "The Equinix Metal project id to use for installation",
			Options: projectName,
		},
		Validate: survey.ComposeValidators(survey.Required),
	}}, &projectID)

	if err != nil {
		return "", err
	}

	parts := strings.Split(projectID, " ")
	return parts[0], nil
}

func selectFacility() (string, error) {
	var facilityID string

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	conn, err := getConnection(Config{APIKey: apiKey, APIURL: apiURL})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Equinix Metal connection")
	}

	client := &Client{Conn: conn}
	facilities, err := client.ListFacilities(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to list Equinix Metal facilities")
	}

	facilitiesNames := []string{}
	for _, f := range facilities {
		facilityNames = append(facilityNames, f.Code+" "+f.Name)
	}

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Equinix Metal Facility Code",
				Help:    "The Equinix Metal Facility code (this is the short name, e.g. 'da11')",
				Default: DefaultFacility,
				Options: facilitiesNames,
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &facilityID)

	if err != nil {
		return "", err
	}
	return facilityID, nil
}

func selectMetro() (string, error) {
	var metroID string

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	conn, err := getConnection(Config{APIKey: apiKey, APIURL: apiURL})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Equinix Metal connection")
	}

	client := &Client{Conn: conn}
	metros, err := client.ListMetros(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to list Equinix Metal metros")
	}

	metrosNames := []string{}
	for _, m := range metros {
		metroNames = append(metroNames, m.Code+" "+m.Name)
	}

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Equinix Metal Metro Code",
				Help:    "The Equinix Metal Metro code (this is the short name, e.g. 'SV')",
				Default: DefaultMetro,
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &metroID)

	if err != nil {
		return "", err
	}
	return metroID, nil
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
			Validate: survey.ComposeValidators(survey.Required),
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
