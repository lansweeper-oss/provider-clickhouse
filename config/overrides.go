package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	})
	p.AddResourceConfigurator("clickhouse_service_private_endpoints_attachment", func(r *config.Resource) {
		r.References = config.References{
			"service_id": {
				TerraformName: "clickhouse_service",
			},
		}
	})
	p.AddResourceConfigurator("clickhouse_service_transparent_data_encryption_key_association", func(r *config.Resource) {
		r.References = config.References{
			"service_id": {
				TerraformName: "clickhouse_service",
			},
		}
	})
}
