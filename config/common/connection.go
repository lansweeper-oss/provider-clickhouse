package common

import (
	"strconv"

	"github.com/crossplane/upjet/v2/pkg/config"
)

const defaultUsername = "default"

var endpointProtocols = []string{"https", "nativesecure", "mysql"}

// EndpointConnectionDetails emits neutral facts instead of a fixed shape so
// consumers stay free to pick their own protocol and auth strategy.
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
		// password is omitted: the default connection-detail mapping already
		// emits it, and re-emitting collides and errors.

		if len(out) == 0 {
			return nil, nil
		}
		return out, nil
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
