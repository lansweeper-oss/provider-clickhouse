package config

import (
	"context"
	"regexp"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// sentinelUUID is substituted for any external-name that is not a real
// ClickHouse Cloud resource id (empty, or the k8s metadata.name that
// Crossplane defaults to when no external-name annotation is set). Upstream
// terraform-provider-clickhouse Read handlers forward the id to the
// ClickHouse API without validation; non-UUID or empty values hit wrong
// endpoints and return arrays that fail to unmarshal into
// ResponseWithResult[api.Service]. The sentinel is a syntactically valid
// UUID with an invalid v4 version nibble, so it cannot collide with a real
// id and routes to GET /services/<sentinel> → 404 → api.IsNotFound.
const sentinelUUID = "ffffffff-ffff-ffff-ffff-ffffffffffff"

// uuidRe matches a canonical UUID. Used to detect a real ClickHouse Cloud
// resource id versus a pre-create placeholder (k8s name or empty).
var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"clickhouse_organization_settings":                               config.IdentifierFromProvider,
	"clickhouse_service":                                             withSentinelWhenNotUUID(config.IdentifierFromProvider),
	"clickhouse_service_private_endpoints_attachment":                withSentinelWhenNotUUID(config.ParameterAsIdentifier("service_id")),
	"clickhouse_service_transparent_data_encryption_key_association": withSentinelWhenNotUUID(config.ParameterAsIdentifier("service_id")),
}

// withSentinelWhenNotUUID wraps an ExternalName so GetIDFn returns sentinelUUID
// whenever the parent would return a value that is not a canonical UUID.
// Covers empty (IdentifierFromProvider pre-create) and k8s metadata.name
// fallback (Crossplane default when external-name annotation is unset).
func withSentinelWhenNotUUID(base config.ExternalName) config.ExternalName {
	parentGetID := base.GetIDFn
	base.GetIDFn = func(ctx context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
		id, err := parentGetID(ctx, externalName, parameters, setup)
		if err != nil {
			return "", err
		}
		if !uuidRe.MatchString(id) {
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
