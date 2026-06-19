package common

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/crossplane/upjet/v2/pkg/config"
)

const defaultUsername = "default"

var endpointProtocols = []string{"https", "nativesecure", "mysql"}

// EndpointConnectionDetails emits neutral per-endpoint facts so consumers stay
// free to pick their own protocol/strategy, plus a provider_credentials blob in
// the shape the default ClickHouse app expects.
func EndpointConnectionDetails() config.AdditionalConnectionDetailsFn {
	return func(attr map[string]any) (map[string][]byte, error) {
		out := map[string][]byte{}
		put := func(key, val string) {
			if val != "" {
				out[key] = []byte(val)
			}
		}

		for _, proto := range endpointProtocols {
			put(proto+"_host", digString(attr, "endpoints", proto, "host"))
			if port := digInt(attr, "endpoints", proto, "port"); port != 0 {
				out[proto+"_port"] = []byte(strconv.Itoa(port))
			}
		}
		put("private_dns_hostname", digString(attr, "private_endpoint_config", "private_dns_hostname"))
		put("username", defaultUsername)

		if err := putProviderCredentials(out, attr); err != nil {
			return nil, err
		}

		if len(out) == 0 {
			return nil, nil
		}
		return out, nil
	}
}

// putProviderCredentials adds the assembled credentials blob, preferring the
// private DNS host and the nativesecure endpoint. Skipped when the service is
// not ready (no host or password yet), to avoid writing a half-built blob.
func putProviderCredentials(out map[string][]byte, attr map[string]any) error {
	host := digString(attr, "private_endpoint_config", "private_dns_hostname")
	if host == "" {
		host = digString(attr, "endpoints", "nativesecure", "host")
	}
	password := digString(attr, "password")
	if host == "" || password == "" {
		return nil
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
		return fmt.Errorf("cannot marshal provider_credentials: %w", err)
	}
	out["provider_credentials"] = b
	return nil
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
