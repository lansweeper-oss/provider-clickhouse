package config

import (
	"context"
	"regexp"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/pkg/errors"
)

const (
	clickhouseService = "clickhouse_service"
	serviceIDParam    = "service_id"
)

// getExternalNameFromServiceID returns a GetExternalNameFn that reads
// "service_id" from tfstate. Use for TF resources whose schema does not expose
// an "id" attribute, where upjet's default IDAsExternalName fails post-apply
// with "cannot find id in tfstate".
func getExternalNameFromServiceID() func(map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		v, ok := tfstate[serviceIDParam]
		if !ok {
			return "", errors.Errorf("cannot find %s in tfstate", serviceIDParam)
		}
		s, ok := v.(string)
		if !ok {
			return "", errors.Errorf("%s in tfstate is not a string", serviceIDParam)
		}
		return s, nil
	}
}

// sentinelUUID stands in for external-names that are not real ClickHouse Cloud
// ids (empty, or the k8s metadata.name default). It must be a well-formed UUID
// the API will never assign, so the Read returns a clean 404 instead of
// unmarshal errors from list endpoints or 400s that upstream's IsNotFound
// misses. v1 chosen over v7 to avoid collision with API-issued ids.
const sentinelUUID = "00000000-0000-1000-8000-000000000000"

// uuidRe matches a canonical UUID. Used to detect a real ClickHouse Cloud
// resource id versus a pre-create placeholder (k8s name or empty).
var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"clickhouse_clickpipe_cdc_infrastructure":                           withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_clickpipe":                                              config.IdentifierFromProvider,
	"clickhouse_clickpipes_reverse_private_endpoint_custom_private_dns": config.IdentifierFromProvider,
	"clickhouse_clickpipes_reverse_private_endpoint":                    config.IdentifierFromProvider,
	"clickhouse_organization_settings":                                  config.IdentifierFromProvider,
	"clickhouse_postgres_service":                                       withSentinelWhenNotUUID(config.IdentifierFromProvider),
	"clickhouse_role_assignment":                                        withSentinelWhenNotUUID(config.ParameterAsIdentifier("role_id")),
	"clickhouse_role":                                                   withSentinelWhenNotUUID(config.IdentifierFromProvider),
	"clickhouse_service_private_endpoints_attachment":                   withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_service_scheduled_scaling":                              withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_service_transparent_data_encryption_key_association":    withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_service_upgrade_window":                                 withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	clickhouseService:                                                   withSentinelWhenNotUUID(config.IdentifierFromProvider),
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
