package common

import (
	"encoding/json"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"
)

const (
	// providerCredentialsKey is the single connection-secret key holding the
	// JSON-encoded credentials blob consumed downstream.
	providerCredentialsKey = "provider_credentials"
	// nativesecureProtocol is the ClickHouse native protocol over TLS.
	nativesecureProtocol = "nativesecure"
	// defaultUsername is the built-in ClickHouse service user.
	defaultUsername = "default"
)

// ProviderCredentialsConnectionDetails returns an AdditionalConnectionDetailsFn
// that emits a single "provider_credentials" key into the connection secret.
//
// The value is JSON of the form:
//
//	{
//	  "host": "<private_dns_hostname || endpoints.nativesecure.host>",
//	  "port": <endpoints.nativesecure.port>,
//	  "protocol": "nativesecure",
//	  "auth_config": {"strategy": "password", "username": "default", "password": "<plaintext>"}
//	}
//
// When the service is not yet ready (no host or no password in the Terraform
// state) it emits nothing rather than an incomplete/erroring secret.
func ProviderCredentialsConnectionDetails() config.AdditionalConnectionDetailsFn {
	return func(attr map[string]any) (map[string][]byte, error) {
		host := digString(attr, "private_endpoint_config", "private_dns_hostname")
		if host == "" {
			host = digString(attr, "endpoints", "nativesecure", "host")
		}
		port := digInt(attr, "endpoints", "nativesecure", "port")
		passwd := digString(attr, "password")

		if host == "" || passwd == "" {
			return nil, nil
		}

		creds := map[string]any{
			"host":     host,
			"port":     port,
			"protocol": nativesecureProtocol,
			"auth_config": map[string]any{
				"strategy": "password",
				"username": defaultUsername,
				"password": passwd,
			},
		}
		b, err := json.Marshal(creds)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal %s: %w", providerCredentialsKey, err)
		}
		return map[string][]byte{providerCredentialsKey: b}, nil
	}
}

// digString walks nested Terraform-state attributes by key and returns the
// string at the end of the path, or "" if any segment is missing or not the
// expected type. Single-element list blocks are transparently unwrapped.
func digString(attr map[string]any, keys ...string) string {
	v := dig(attr, keys...)
	s, _ := v.(string)
	return s
}

// digInt is like digString but returns an int, coercing float64 (the type
// JSON numbers decode to) into int. Returns 0 when missing or not numeric.
func digInt(attr map[string]any, keys ...string) int {
	switch n := dig(attr, keys...).(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

// dig walks keys through nested map[string]any values. When it encounters a
// []any (a Terraform block rendered as a single-element list), it descends into
// the first element. Returns nil if the path cannot be resolved.
func dig(attr map[string]any, keys ...string) any {
	var cur any = attr
	for _, k := range keys {
		m := asMap(cur)
		if m == nil {
			return nil
		}
		cur = m[k]
	}
	return cur
}

// asMap normalises a value into a map[string]any, unwrapping a single-element
// []any block if necessary. Returns nil when not map-shaped.
func asMap(v any) map[string]any {
	switch t := v.(type) {
	case map[string]any:
		return t
	case []any:
		if len(t) == 0 {
			return nil
		}
		m, _ := t[0].(map[string]any)
		return m
	default:
		return nil
	}
}
