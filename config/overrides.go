package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/lansweeper-oss/provider-clickhouse/config/common"
)

var gkvOverrideMap = map[string]schema.GroupVersionKind{}

func gvkOverride() config.ResourceOption {
	return func(r *config.Resource) {
		if r.ShortGroup == resourcePrefix {
			r.ShortGroup = "clickhouse"
		}
		if gvk, ok := gkvOverrideMap[r.Name]; ok {
			r.ShortGroup = gvk.Group
			r.Kind = gvk.Kind
			if gvk.Version != "" {
				r.Version = gvk.Version
			}
		}
	}
}

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("clickhouse_service", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"warehouse_id", "backup_configuration", "password"},
		}
		r.InitializerFns = append(r.InitializerFns,
			common.PasswordGenerator(
				"spec.forProvider.passwordSecretRef",
				"spec.writeConnectionSecretToRef",
			))
		if pw, ok := r.TerraformResource.Schema["password"]; ok {
			pw.Description = "Password for the default ClickHouse user.\n" +
				"When passwordSecretRef is set, that secret is used (Bring Your Own Password).\n" +
				"Otherwise a password is auto-generated and written to writeConnectionSecretToRef."
		}
	})
	p.AddResourceConfigurator("clickhouse_clickpipe", func(r *config.Resource) {
		// Upjet only supports primitive types (string, *string, …) as sensitive.
		// The upstream TF provider marks several nested credential blocks as
		// sensitive at the parent level. Clear sensitive on those parent objects
		// so upjet can generate CRDs; leaf string fields keep their own
		// Sensitive flag and become proper SecretRef fields.
		clearSensitiveOnNestedBlocks(r.TerraformResource.Schema)
	})
	p.AddResourceConfigurator("clickhouse_service_private_endpoints_attachment", func(r *config.Resource) {
		r.References = config.References{
			"service_id": {
				TerraformName: "clickhouse_service",
			},
		}
		r.ExternalName.GetExternalNameFn = getExternalNameFromTFParam("service_id")
	})
	p.AddResourceConfigurator("clickhouse_service_transparent_data_encryption_key_association", func(r *config.Resource) {
		r.References = config.References{
			"service_id": {
				TerraformName: "clickhouse_service",
			},
		}
		r.ExternalName.GetExternalNameFn = getExternalNameFromTFParam("service_id")
	})
}

// clearSensitiveOnNestedBlocks recursively walks an SDKv2 schema map and
// clears Sensitive on any entry that contains a nested Resource (i.e. is not a
// leaf). Leaf fields retain their Sensitive flag so upjet generates SecretRef
// fields for them. This is needed because upjet panics when a complex/nested
// type is marked sensitive — only primitive types are supported.
func clearSensitiveOnNestedBlocks(m map[string]*sdkschema.Schema) {
	for _, s := range m {
		if s.Elem == nil {
			continue
		}
		r, ok := s.Elem.(*sdkschema.Resource)
		if !ok {
			continue
		}
		s.Sensitive = false
		clearSensitiveOnNestedBlocks(r.Schema)
	}
}
