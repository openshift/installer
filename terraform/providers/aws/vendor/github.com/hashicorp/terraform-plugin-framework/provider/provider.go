package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Provider is the core interface that all Terraform providers must implement.
//
// Providers can optionally implement these additional concepts:
//
//   - Resources: ProviderWithResources or (deprecated)
//     ProviderWithGetResources.
//   - Data Sources: ProviderWithDataSources or (deprecated)
//     ProviderWithGetDataSources.
//   - Validation: Schema-based via tfsdk.Attribute or entire configuration
//     via ProviderWithConfigValidators or ProviderWithValidateConfig.
//   - Meta Schema: ProviderWithMetaSchema
type Provider interface {
	// GetSchema returns the schema for this provider's configuration. If
	// this provider has no configuration, return an empty schema.Schema.
	GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics)

	// Configure is called at the beginning of the provider lifecycle, when
	// Terraform sends to the provider the values the user specified in the
	// provider configuration block. These are supplied in the
	// ConfigureProviderRequest argument.
	// Values from provider configuration are often used to initialise an
	// API client, which should be stored on the struct implementing the
	// Provider interface.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)

	// DataSources returns a slice of functions to instantiate each DataSource
	// implementation.
	//
	// The data source type name is determined by the DataSource implementing
	// the Metadata method. All data sources must have unique names.
	DataSources(context.Context) []func() datasource.DataSource

	// Resources returns a slice of functions to instantiate each Resource
	// implementation.
	//
	// The resource type name is determined by the Resource implementing
	// the Metadata method. All resources must have unique names.
	Resources(context.Context) []func() resource.Resource
}

// ProviderWithConfigValidators is an interface type that extends Provider to include declarative validations.
//
// Declaring validation using this methodology simplifies implementation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ProviderWithConfigValidators interface {
	Provider

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ConfigValidators(context.Context) []ConfigValidator
}

// ProviderWithMetadata is an interface type that extends Provider to
// return its type name, such as examplecloud, and other
// metadata, such as version.
//
// Implementing this method will populate the
// [datasource.MetadataRequest.ProviderTypeName] and
// [resource.MetadataRequest.ProviderTypeName] fields automatically.
type ProviderWithMetadata interface {
	Provider

	// Metadata should return the metadata for the provider, such as
	// a type name and version data.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)
}

// ProviderWithMetaSchema is a provider with a provider meta schema.
// This functionality is currently experimental and subject to change or break
// without warning; it should only be used by providers that are collaborating
// on its use with the Terraform team.
type ProviderWithMetaSchema interface {
	Provider

	// GetMetaSchema returns the provider meta schema.
	GetMetaSchema(context.Context) (tfsdk.Schema, diag.Diagnostics)
}

// ProviderWithValidateConfig is an interface type that extends Provider to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single provider. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ProviderWithValidateConfig interface {
	Provider

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
