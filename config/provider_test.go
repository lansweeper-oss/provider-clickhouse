package config

import (
	"reflect"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// TestExternalTags asserts that resources with a tags map get upjet's tagger,
// and that those with tags as a list do not.
func TestExternalTags(t *testing.T) {
	want := map[string]bool{
		"clickhouse_service":             true,
		"clickhouse_postgres_service":    true,
		"clickhouse_byoc_infrastructure": true,
	}
	tagger := reflect.ValueOf(ujconfig.TagInitializer).Pointer()
	for _, p := range []*ujconfig.Provider{GetProvider(), GetProviderNamespaced()} {
		for name, r := range p.Resources {
			got := false
			for _, fn := range r.InitializerFns {
				got = got || reflect.ValueOf(fn).Pointer() == tagger
			}
			if got != want[name] {
				t.Errorf("%s: tagger registered = %t, want %t", name, got, want[name])
			}
		}
	}
}
