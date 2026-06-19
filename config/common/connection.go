package common

import (
	"encoding/json"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"
)

const defaultUsername = "default"

// ProviderCredentialsConnectionDetails emits a single provider_credentials key
// whose JSON value is the shape the ClickHouse app consumes. Host prefers the
// private DNS name, falling back to the nativesecure endpoint. Skipped until
// host and password are known, to avoid writing a half-built blob.
func ProviderCredentialsConnectionDetails() config.AdditionalConnectionDetailsFn {
	return func(attr map[string]any) (map[string][]byte, error) {
		host := digString(attr, "private_endpoint_config", "private_dns_hostname")
		if host == "" {
			host = digString(attr, "endpoints", "nativesecure", "host")
		}
		password := digString(attr, "password")
		if host == "" || password == "" {
			return nil, nil
		}

		creds := map[string]any{
			"host":     host,
			"port":     digInt(attr, "endpoints", "nativesecure", "port"),
			"protocol": "nativesecure",
			"auth_config": map[string]any{
				"strategy": "password",
				"username": defaultUsername,
				"password": password,
			},
		}
		b, err := json.Marshal(creds)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal provider_credentials: %w", err)
		}
		return map[string][]byte{"provider_credentials": b}, nil
	}
}

func digString(attr map[string]any, keys ...string) string {
	s, _ := dig(attr, keys...).(string)
	return s
}

func digInt(attr map[string]any, keys ...string) int {
	switch n := dig(attr, keys...).(type) {
	case float64: // TF JSON numbers decode as float64
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

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

func asMap(v any) map[string]any {
	switch t := v.(type) {
	case map[string]any:
		return t
	case []any: // TF blocks can render as a single-element list
		if len(t) == 0 {
			return nil
		}
		m, _ := t[0].(map[string]any)
		return m
	default:
		return nil
	}
}
