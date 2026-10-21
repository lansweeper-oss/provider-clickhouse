package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	apisCluster "github.com/lansweeper-oss/provider-clickhouse/apis/cluster"
	apisNamespaced "github.com/lansweeper-oss/provider-clickhouse/apis/namespaced"
)

func TestProviderGroups(t *testing.T) {
	s := runtime.NewScheme()
	require.NoError(t, clientgoscheme.AddToScheme(s))
	require.NoError(t, apisCluster.AddToScheme(s))
	require.NoError(t, apisNamespaced.AddToScheme(s))

	got := providerGroups(s, "clickhouse.crossplane.io", "clickhouse.m.crossplane.io")

	assert.Contains(t, got, "clickhouse.crossplane.io")
	assert.Contains(t, got, "clickhouse.m.crossplane.io")
	assert.Contains(t, got, "iam.clickhouse.crossplane.io")
	assert.Contains(t, got, "iam.clickhouse.m.crossplane.io")
	assert.NotContains(t, got, "")
	assert.NotContains(t, got, "apps")
}
