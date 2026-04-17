package config

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// sentinelUUID is substituted for an empty external-name so the upstream
// terraform-provider-clickhouse Read handlers do not issue a GET on the
// collection endpoint (trailing slash on empty id), which returns a JSON array
// and fails to unmarshal into ResponseWithResult[api.Service]. A non-v4 UUID
// is used so it cannot collide with a real service id assigned by ClickHouse
// Cloud. GET /services/<sentinel> returns 404, handled cleanly by
// api.IsNotFound in the upstream provider.
const sentinelUUID = "ffffffff-ffff-ffff-ffff-ffffffffffff"

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"clickhouse_organization_settings":                               config.IdentifierFromProvider,
	"clickhouse_service":                                             withSentinelOnEmpty(config.IdentifierFromProvider),
	"clickhouse_service_private_endpoints_attachment":                withSentinelOnEmpty(config.ParameterAsIdentifier("service_id")),
	"clickhouse_service_transparent_data_encryption_key_association": withSentinelOnEmpty(config.ParameterAsIdentifier("service_id")),
}

// withSentinelOnEmpty wraps an ExternalName so GetIDFn returns sentinelUUID
// whenever the parent would return an empty id. Required for resources whose
// upstream Read handler calls the ClickHouse API without guarding against an
// empty identifier.
func withSentinelOnEmpty(base config.ExternalName) config.ExternalName {
	parentGetID := base.GetIDFn
	base.GetIDFn = func(ctx context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
		id, err := parentGetID(ctx, externalName, parameters, setup)
		if err != nil {
			return "", err
		}
		if id == "" {
			return sentinelUUID, nil
		}
		return id, nil
	}
	return base
}

// ExternalNameConfigured returns the list of possible external name
// configurations for this provider.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		l[i] = name
		i++
	}
	return l
}

// ExternalNameConfigurations applies all external name configurations for each
// group resource separately.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}
