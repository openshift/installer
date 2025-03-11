package catalogmanagement

import (
	"context"
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func SIToSS(i []interface{}) []string {
	var ss []string
	for _, iface := range i {
		ss = append(ss, iface.(string))
	}
	return ss
}

// FetchOfferingWithAllVersions fetches an offering and its additional versions
func FetchOfferingWithAllVersions(ctx context.Context, client *catalogmanagementv1.CatalogManagementV1, getOfferingOptions *catalogmanagementv1.GetOfferingOptions) (*catalogmanagementv1.Offering, *core.DetailedResponse, error) {
	offering, response, err := client.GetOfferingWithContext(ctx, getOfferingOptions)
	if err != nil {
		return offering, response, err
	}

	// Fetch additional versions for the offering
	if err := FetchAdditionalVersions(ctx, client, offering); err != nil {
		return offering, response, err
	}

	return offering, response, nil
}

// FetchAdditionalVersions fetches all versions for each kind in the offering
func FetchAdditionalVersions(ctx context.Context, client *catalogmanagementv1.CatalogManagementV1, offering *catalogmanagementv1.Offering) error {
	// Define the recursive function to fetch versions for a kind
	var fetchVersionsForKind func(kind *catalogmanagementv1.Kind, start *string) error

	fetchVersionsForKind = func(kind *catalogmanagementv1.Kind, start *string) error {
		fmt.Println("START STRING: ", start)
		// Fetch versions for the kind
		getVersionsOptions := &catalogmanagementv1.GetVersionsOptions{
			CatalogIdentifier: offering.CatalogID,
			OfferingID:        offering.ID,
			KindID:            kind.ID,
			Start:             start,
		}

		result, _, err := client.GetVersionsWithContext(ctx, getVersionsOptions)
		if err != nil {
			return err
		}

		if result.Versions != nil {
			// Append fetched versions to the kind
			kind.Versions = append(kind.Versions, result.Versions...)
		}

		// Check if there are more versions to fetch
		if result.Next != nil && result.Next.Start != nil && *result.Next.Start != "" {
			// If there are more versions to fetch, call recursively
			return fetchVersionsForKind(kind, result.Next.Start)
		}

		return nil
	}

	// Iterate over kinds in the offering and fetch additional versions
	for i := range offering.Kinds {
		kind := &offering.Kinds[i] // Use a pointer to modify the original kind
		if kind.AllVersions != nil && kind.AllVersions.Next != nil && kind.AllVersions.Next.Start != nil && *kind.AllVersions.Next.Start != "" {
			// Load additional versions for this kind
			if err := fetchVersionsForKind(kind, kind.AllVersions.Next.Start); err != nil {
				return err
			}
		}
	}

	return nil
}
