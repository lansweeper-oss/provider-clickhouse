package common

import (
	"strconv"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// defaultUsername is the built-in ClickHouse service admin user that the
// generated/BYOP password belongs to.
const defaultUsername = "default"

// endpointProtocols are the per-protocol endpoint blocks exposed in the
// clickhouse_service Terraform state.
var endpointProtocols = []string{"https", "nativesecure", "mysql"}

// EndpointConnectionDetails returns an AdditionalConnectionDetailsFn that emits
// neutral, policy-free connection facts into the connection secret. It does NOT
// pick a protocol or auth strategy — it surfaces every endpoint the service
// exposes so the consumer can assemble whatever connection shape it needs.
//
// Keys emitted (empty values skipped):
//
//	<protocol>_host / <protocol>_port   for https, nativesecure, mysql
//	private_dns_hostname                 (private endpoint, when configured)
//	username                             ("default" service admin user)
//
// The password is already written by the resource's default connection-detail
// mapping, so it is intentionally not re-emitted here.
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

		if len(out) == 0 {
			return nil, nil
		}
		return out, nil
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
