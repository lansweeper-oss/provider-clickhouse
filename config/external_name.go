package config

import (
	"context"
	"regexp"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// sentinelUUID substitutes for any external-name that is not a real ClickHouse
// Cloud resource id (empty, or the k8s metadata.name Crossplane defaults to
// when the external-name annotation is unset). Upstream
// terraform-provider-clickhouse Read handlers forward the id to the ClickHouse
// API without validation; non-UUID values route to list endpoints and return
// arrays that fail to unmarshal into ResponseWithResult[api.Service].
//
// Value must be a well-formed RFC 4122 UUID (valid version nibble 1-5,
// variant nibble 8/9/a/b): the ClickHouse Cloud API rejects malformed UUIDs
// with 400 BAD_REQUEST, which upstream's api.IsNotFound (strings.HasPrefix
// "status: 404") does not catch. This UUID cannot be assigned by the API,
// so the request returns 404 and is handled cleanly by syncServiceState.
//
// Deliberately v1 (only version/variant bits set, all other bits zero), NOT
// v7 (which ClickHouse Cloud itself issues for real resource ids): v7 risks
// colliding with a real service id and routing the sentinel request to an
// existing resource instead of returning 404.
const sentinelUUID = "00000000-0000-1000-8000-000000000000"

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
