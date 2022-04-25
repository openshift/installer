package gcp

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	gcpclient "github.com/openshift/installer/pkg/client/gcp"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/gcp/validation"
)

// Platform collects GCP-specific configuration.
func Platform() (*gcp.Platform, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	project, err := selectProject(ctx)
	if err != nil {
		return nil, err
	}

	region, err := selectRegion(project)
	if err != nil {
		return nil, err
	}

	return &gcp.Platform{
		ProjectID: project,
		Region:    region,
	}, nil
}

func selectProject(ctx context.Context) (string, error) {
	ssn, err := gcpclient.GetSession(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get session")
	}
	defaultProject := ssn.Credentials.ProjectID

	client := &gcpclient.Client{
		SSN: ssn,
	}

	projects, err := client.GetProjects(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get projects")
	}

	var options []string
	ids := make(map[string]string)

	var defaultValue string

	for id, name := range projects {
		option := fmt.Sprintf("%s (%s)", name, id)
		ids[option] = id
		if id == defaultProject {
			defaultValue = option
		}
		options = append(options, option)
	}
	sort.Strings(options)

	var selectedProject string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Project ID",
				Help:    "The project id where the cluster will be provisioned. The default is taken from the provided service account.",
				Default: defaultValue,
				Options: options,
			},
		},
	}, &selectedProject)

	selectedProject = ids[selectedProject]
	return selectedProject, err
}

func selectRegion(project string) (string, error) {
	longRegions := make([]string, 0, len(validation.Regions))
	shortRegions := make([]string, 0, len(validation.Regions))
	for id, location := range validation.Regions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	var regionTransform survey.Transformer = func(ans interface{}) interface{} {
		switch v := ans.(type) {
		case core.OptionAnswer:
			return core.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[0], Index: v.Index}
		case string:
			return strings.SplitN(v, " ", 2)[0]
		}
		return ""
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)

	defaultRegion := "us-central1"
	var selectedRegion string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The GCP region to be used for installation.",
				Default: fmt.Sprintf("%s (%s)", defaultRegion, validation.Regions[defaultRegion]),
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(core.OptionAnswer).Value
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
	}, &selectedRegion)
	if err != nil {
		return "", err
	}

	return selectedRegion, nil
}
