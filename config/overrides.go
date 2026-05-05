package config

import (
	"encoding/json"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gkvOverrideMap = map[string]schema.GroupVersionKind{}

// serviceCredentialsConnectionDetails emits a "credentials" key into the
// clickhouse_service connection secret containing the JSON blob expected by
// provider-clickhousedbops ProviderConfig (host, port, protocol, auth_config).
// Sourced from the TF attributes after apply: password from attr.password,
// host/port from attr.endpoints.nativesecure. Lets downstream ProviderConfigs
// consume the Service connection secret directly without a composed Object.
// See https://github.com/crossplane/upjet/blob/main/docs/configuring-a-resource.md#additional-sensitive-fields-and-custom-connection-details
func serviceCredentialsConnectionDetails(attr map[string]any) (map[string][]byte, error) {
	password, _ := attr["password"].(string)
	if password == "" {
		return nil, nil
	}
	endpoints, _ := attr["endpoints"].(map[string]any)
	nativeSecure, _ := endpoints["nativesecure"].(map[string]any)
	host, _ := nativeSecure["host"].(string)
	if host == "" {
		return nil, nil
	}
	port, _ := nativeSecure["port"].(float64)
	if port == 0 {
		port = 9440
	}
	creds := map[string]any{
		"host":     host,
		"port":     int(port),
		"protocol": "nativesecure",
		"auth_config": map[string]any{
			"strategy": "password",
			"username": "default",
			"password": password,
		},
	}
	b, err := json.Marshal(creds)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal credentials")
	}
	return map[string][]byte{"credentials": b}, nil
}

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
		r.Sensitive.AdditionalConnectionDetailsFn = serviceCredentialsConnectionDetails
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
