package common

import (
	"strconv"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// defaultUsername is the ClickHouse service admin user.
const defaultUsername = "default"

var endpointProtocols = []string{"https", "nativesecure", "mysql"}

// EndpointConnectionDetails emits neutral connection facts (every endpoint's
// host/port, private DNS, username) so consumers can build any connection
// shape. Password comes from the default mapping, so it is not re-emitted.
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

// digString returns the string at the key path, or "" if missing/wrong type.
func digString(attr map[string]any, keys ...string) string {
	v := dig(attr, keys...)
	s, _ := v.(string)
	return s
}

// digInt returns the int at the key path (float64 coerced), or 0.
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

// dig walks keys through nested maps, unwrapping single-element list blocks.
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

// asMap coerces v to a map, unwrapping a single-element list. Nil if neither.
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
