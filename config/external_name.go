package config

import (
	"context"
	"regexp"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/pkg/errors"
)

const (
	clickhouseService = "clickhouse_service"
	clickstackGroup   = "clickstack"
	serviceIDParam    = "service_id"
	teamParam         = "team"
	clickstackTeam    = "clickhouse_clickstack_team"
	clickstackSource  = "clickhouse_clickstack_source"
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
var uuidRe = regexp.MustCompile(`^[[:xdigit:]]{8}(-[[:xdigit:]]{4}){3}-[[:xdigit:]]{12}$`)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// ClickHouse Cloud resources
	"clickhouse_clickpipe_cdc_infrastructure":                           withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_clickpipe":                                              config.IdentifierFromProvider,
	"clickhouse_clickpipes_reverse_private_endpoint_custom_private_dns": config.IdentifierFromProvider,
	"clickhouse_clickpipes_reverse_private_endpoint":                    config.IdentifierFromProvider,
	"clickhouse_organization_settings":                                  config.IdentifierFromProvider,
	"clickhouse_role_assignment":                                        withSentinelWhenNotUUID(config.ParameterAsIdentifier("role_id")),
	"clickhouse_role":                                                   withSentinelWhenNotUUID(config.IdentifierFromProvider),
	"clickhouse_service_private_endpoints_attachment":                   withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_service_scheduled_scaling":                              withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_service_transparent_data_encryption_key_association":    withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_service_upgrade_window":                                 withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_udf_attachment":                                         withSentinelWhenNotUUID(config.ParameterAsIdentifier(serviceIDParam)),
	"clickhouse_udf":                                                    identifierFromParameter("function_name"),
	clickhouseService:                                                   withSentinelWhenNotUUID(config.IdentifierFromProvider),
	// Postgres resources
	"clickhouse_postgres_service": withSentinelWhenNotUUID(config.IdentifierFromProvider),
	// ClickStack resources
	"clickhouse_clickstack_alert":        config.IdentifierFromProvider,
	"clickhouse_clickstack_connection":   config.IdentifierFromProvider,
	"clickhouse_clickstack_dashboard":    config.IdentifierFromProvider,
	"clickhouse_clickstack_role":         config.IdentifierFromProvider,
	"clickhouse_clickstack_saved_search": config.IdentifierFromProvider,
	clickstackSource:                     config.IdentifierFromProvider,
	clickstackTeam:                       config.IdentifierFromProvider,
	"clickhouse_clickstack_team_member":  withOptionalPrefix(identifierFromParameterOrAnnotation("email"), "team"),
	"clickhouse_clickstack_webhook":      config.IdentifierFromProvider,
}

// identifierFromParameter returns an ExternalName config that keeps "param" in spec.forProvider
// and mirrors it to the external-name annotation after observe.
// Based on IdentifierFromProvider, but reads "param" from the tfstate.
func identifierFromParameter(param string) config.ExternalName {
	e := config.IdentifierFromProvider
	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		v, ok := tfstate[param]
		if !ok {
			return "", errors.Errorf("cannot find %s in tfstate", param)
		}
		s, ok := v.(string)
		if !ok {
			return "", errors.Errorf("%s in tfstate is not a string", param)
		}
		return s, nil
	}
	return e
}

// withOptionalPrefix wraps an ExternalName so GetIDFn prepends "prefix/" when the given parameter
// is non-empty. Covers resources whose TF import accepts both "value" and "prefix/value" forms.
func withOptionalPrefix(base config.ExternalName, prefixParam string) config.ExternalName {
	parentGetID := base.GetIDFn
	base.GetIDFn = func(ctx context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
		id, err := parentGetID(ctx, externalName, parameters, setup)
		if err != nil {
			return "", err
		}
		if prefix, _ := parameters[prefixParam].(string); prefix != "" {
			return prefix + "/" + id, nil
		}
		return id, nil
	}
	return base
}

// identifierFromParameterOrAnnotation builds an ExternalName whose GetIDFn
// prefers the named parameter from spec.forProvider, falling back to the
// external-name annotation when the parameter is empty (import scenario).
func identifierFromParameterOrAnnotation(param string) config.ExternalName {
	e := identifierFromParameter(param)
	e.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, _ map[string]any) (string, error) {
		if v, _ := parameters[param].(string); v != "" {
			return v, nil
		}
		if externalName != "" {
			return externalName, nil
		}
		return "", errors.Errorf("cannot determine ID: neither %s parameter nor external-name annotation is set", param)
	}
	return e
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
