package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/ClickHouse/terraform-provider-clickhouse/pkg/provider"
	"github.com/ClickHouse/terraform-provider-clickhouse/pkg/resource"
)

const (
	resourcePrefix = "clickhouse"
	modulePath     = "github.com/lansweeper-oss/provider-clickhouse"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithIncludeList(nil),
		ujconfig.WithTerraformPluginFrameworkIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithRootGroup(resourcePrefix+".crossplane.io"),
		ujconfig.WithTerraformPluginFrameworkProvider(provider.NewBuilder(resource.GetResourceFactories())()),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
			gvkOverride(),
		))

	Configure(pc)

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithIncludeList(nil),
		ujconfig.WithTerraformPluginFrameworkIncludeList(ExternalNameConfigured()),
		ujconfig.WithShortName(resourcePrefix),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithRootGroup(resourcePrefix+".m.crossplane.io"),
		ujconfig.WithTerraformPluginFrameworkProvider(provider.NewBuilder(resource.GetResourceFactories())()),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
			gvkOverride(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	Configure(pc)

	pc.ConfigureResources()
	return pc
}
