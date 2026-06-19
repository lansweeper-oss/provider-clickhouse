package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigString(t *testing.T) {
	cases := map[string]struct {
		attr map[string]any
		keys []string
		want string
	}{
		"nil map":         {attr: nil, keys: []string{"a"}, want: ""},
		"missing key":     {attr: map[string]any{}, keys: []string{"a"}, want: ""},
		"wrong type":      {attr: map[string]any{"a": 1}, keys: []string{"a"}, want: ""},
		"top level":       {attr: map[string]any{"a": "v"}, keys: []string{"a"}, want: "v"},
		"nested map":      {attr: map[string]any{"a": map[string]any{"b": "v"}}, keys: []string{"a", "b"}, want: "v"},
		"missing nested":  {attr: map[string]any{"a": map[string]any{}}, keys: []string{"a", "b"}, want: ""},
		"list-wrapped":    {attr: map[string]any{"a": []any{map[string]any{"b": "v"}}}, keys: []string{"a", "b"}, want: "v"},
		"empty list":      {attr: map[string]any{"a": []any{}}, keys: []string{"a", "b"}, want: ""},
		"intermediate nil": {attr: map[string]any{"a": nil}, keys: []string{"a", "b"}, want: ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, digString(tc.attr, tc.keys...))
		})
	}
}

func TestDigInt(t *testing.T) {
	cases := map[string]struct {
		attr map[string]any
		keys []string
		want int
	}{
		"nil map":       {attr: nil, keys: []string{"p"}, want: 0},
		"missing":       {attr: map[string]any{}, keys: []string{"p"}, want: 0},
		"float64":       {attr: map[string]any{"p": float64(9440)}, keys: []string{"p"}, want: 9440},
		"int":           {attr: map[string]any{"p": 9440}, keys: []string{"p"}, want: 9440},
		"wrong type":    {attr: map[string]any{"p": "9440"}, keys: []string{"p"}, want: 0},
		"nested float":  {attr: map[string]any{"e": map[string]any{"p": float64(8443)}}, keys: []string{"e", "p"}, want: 8443},
		"list-wrapped":  {attr: map[string]any{"e": []any{map[string]any{"p": float64(9440)}}}, keys: []string{"e", "p"}, want: 9440},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, digInt(tc.attr, tc.keys...))
		})
	}
}

type providerCredentials struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	AuthConfig struct {
		Strategy string `json:"strategy"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth_config"`
}

func TestProviderCredentialsConnectionDetails(t *testing.T) {
	fn := ProviderCredentialsConnectionDetails()

	t.Run("happy path uses nativesecure endpoint", func(t *testing.T) {
		attr := map[string]any{
			"password": "s3cret",
			"endpoints": map[string]any{
				"nativesecure": map[string]any{
					"host": "abc.clickhouse.cloud",
					"port": float64(9440),
				},
			},
		}
		out, err := fn(attr)
		require.NoError(t, err)
		require.Contains(t, out, providerCredentialsKey)

		var got providerCredentials
		require.NoError(t, json.Unmarshal(out[providerCredentialsKey], &got))
		assert.Equal(t, "abc.clickhouse.cloud", got.Host)
		assert.Equal(t, 9440, got.Port)
		assert.Equal(t, "nativesecure", got.Protocol)
		assert.Equal(t, "password", got.AuthConfig.Strategy)
		assert.Equal(t, "default", got.AuthConfig.Username)
		assert.Equal(t, "s3cret", got.AuthConfig.Password)
	})

	t.Run("private dns hostname preferred over endpoint host", func(t *testing.T) {
		attr := map[string]any{
			"password": "s3cret",
			"private_endpoint_config": map[string]any{
				"private_dns_hostname": "private.dns.internal",
			},
			"endpoints": map[string]any{
				"nativesecure": map[string]any{
					"host": "public.clickhouse.cloud",
					"port": float64(9440),
				},
			},
		}
		out, err := fn(attr)
		require.NoError(t, err)
		var got providerCredentials
		require.NoError(t, json.Unmarshal(out[providerCredentialsKey], &got))
		assert.Equal(t, "private.dns.internal", got.Host)
		assert.Equal(t, 9440, got.Port)
	})

	t.Run("falls back to endpoint host when private dns empty", func(t *testing.T) {
		attr := map[string]any{
			"password": "s3cret",
			"private_endpoint_config": map[string]any{
				"private_dns_hostname": "",
			},
			"endpoints": map[string]any{
				"nativesecure": map[string]any{
					"host": "public.clickhouse.cloud",
					"port": float64(9440),
				},
			},
		}
		out, err := fn(attr)
		require.NoError(t, err)
		var got providerCredentials
		require.NoError(t, json.Unmarshal(out[providerCredentialsKey], &got))
		assert.Equal(t, "public.clickhouse.cloud", got.Host)
	})

	t.Run("missing password returns nil", func(t *testing.T) {
		attr := map[string]any{
			"endpoints": map[string]any{
				"nativesecure": map[string]any{"host": "h", "port": float64(9440)},
			},
		}
		out, err := fn(attr)
		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("missing host returns nil", func(t *testing.T) {
		attr := map[string]any{"password": "s3cret"}
		out, err := fn(attr)
		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("host present but nativesecure port absent yields port 0", func(t *testing.T) {
		attr := map[string]any{
			"password": "s3cret",
			"private_endpoint_config": map[string]any{
				"private_dns_hostname": "private.dns.internal",
			},
		}
		out, err := fn(attr)
		require.NoError(t, err)
		var got providerCredentials
		require.NoError(t, json.Unmarshal(out[providerCredentialsKey], &got))
		assert.Equal(t, "private.dns.internal", got.Host)
		assert.Equal(t, 0, got.Port)
	})
}
