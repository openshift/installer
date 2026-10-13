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

	"github.com/openshift/installer/pkg/types/gcp"
	gcpValidation "github.com/openshift/installer/pkg/types/gcp/validation"
)

// Platform collects GCP-specific configuration.
func Platform() (*gcp.Platform, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	project, err := selectProject(ctx)
	if err != nil {
		return nil, err
	}

	region, err := selectRegion(ctx, project)
	if err != nil {
		return nil, err
	}

	return &gcp.Platform{
		ProjectID: project,
		Region:    region,
	}, nil
}

func selectProject(ctx context.Context) (string, error) {
	ssn, err := GetSession(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get session")
	}
	defaultProject := ssn.Credentials.ProjectID

	client := &Client{
		ssn: ssn,
	}

	projects, err := client.GetProjects(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get projects")
	} else if len(projects) == 0 {
		return "", fmt.Errorf("failed to get projects for the given service principal, please check your permissions")
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

func getValidatedRegions(computeRegions []string) map[string]string {
	validatedRegions := make(map[string]string)
	for _, region := range computeRegions {
		// Only add validated regions
		if value, ok := gcpValidation.Regions[region]; ok {
			validatedRegions[region] = value
		}
	}

	return validatedRegions
}

func selectRegion(ctx context.Context, project string) (string, error) {
	ssn, err := GetSession(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get session")
	}

	client := &Client{
		ssn: ssn,
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	computeRegions, err := client.GetRegions(ctx, project)
	if err != nil || len(computeRegions) == 0 {
		return "", errors.Wrap(err, "failed to get regions")
	}

	validRegions := getValidatedRegions(computeRegions)

	defaultRegion := "us-central1"
	defaultRegionName := ""
	longRegions := make([]string, 0, len(validRegions))
	shortRegions := make([]string, 0, len(validRegions))
	for key, value := range validRegions {
		shortRegions = append(shortRegions, key)
		regionDesc := fmt.Sprintf("%s (%s)", key, value)
		longRegions = append(longRegions, regionDesc)

		if defaultRegionName == "" && key == defaultRegion {
			defaultRegionName = regionDesc
		}
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

	if defaultRegionName == "" && len(longRegions) > 0 {
		defaultRegionName = longRegions[0]
	}

	var selectedRegion string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The GCP region to be used for installation.",
				Default: defaultRegionName,
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
