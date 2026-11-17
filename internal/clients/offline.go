package clients

import (
	"context"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// offlineProvider never configures a ClickHouse Cloud API client. Resources
// get no provider data, so they skip the plan-time checks that call the API
// (e.g. the organization eligibility check on services) and plans are computed
// without credentials and without network access.
//
// shortcut: clickhouse_organization_settings' ModifyPlan reads core_dumps_enabled
// from the API without a nil-client check, so planning its creation with that
// field unset panics (the diff server recovers and fails that one request).
// Drop this note once upstream guards the call as it does for services.
type offlineProvider struct {
	fwprovider.Provider
}

func (offlineProvider) Configure(context.Context, fwprovider.ConfigureRequest, *fwprovider.ConfigureResponse) {
}

// OfflineTerraformSetupBuilder returns the terraform.SetupFn used by the diff
// server. It ignores the ProviderConfig entirely: no credentials are read and
// no connection to ClickHouse Cloud is opened.
func OfflineTerraformSetupBuilder(frameworkProvider fwprovider.Provider) terraform.SetupFn {
	return func(context.Context, client.Client, resource.Managed) (terraform.Setup, error) {
		return terraform.Setup{
			FrameworkProvider: offlineProvider{Provider: frameworkProvider},
			Configuration:     map[string]any{},
		}, nil
	}
}
