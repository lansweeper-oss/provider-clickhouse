package clients

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// offlineProvider never configures a ClickHouse Cloud API client. Resources
// get no provider data, so they skip the plan-time checks that call the API
// (e.g. the organization eligibility check on services) and plans are computed
// without credentials and without network access.
type offlineProvider struct {
	fwprovider.Provider
	log logging.Logger
}

func (offlineProvider) Configure(context.Context, fwprovider.ConfigureRequest, *fwprovider.ConfigureResponse) {
}

func (p offlineProvider) Resources(ctx context.Context) []func() fwresource.Resource {
	factories := p.Provider.Resources(ctx)
	offline := make([]func() fwresource.Resource, len(factories))
	for i, f := range factories {
		offline[i] = func() fwresource.Resource { return offlineResource{Resource: f(), log: p.log} }
	}
	return offline
}

// offlineResource plans a resource without an API client. Some upstream
// ModifyPlan implementations call the client without a nil check (e.g.
// clickhouse_organization_settings reading core_dumps_enabled on create), so a
// panic is recovered and the plan reverts to the framework's, as if the
// resource had no ModifyPlan: values only the API knows stay unknown.
//
// It only exposes Configure and ModifyPlan, the optional interfaces the
// framework's PlanResourceChange uses; the diff server calls nothing else.
type offlineResource struct {
	fwresource.Resource
	log logging.Logger
}

func (r offlineResource) Configure(ctx context.Context, req fwresource.ConfigureRequest, resp *fwresource.ConfigureResponse) {
	if c, ok := r.Resource.(fwresource.ResourceWithConfigure); ok {
		c.Configure(ctx, req, resp)
	}
}

func (r offlineResource) ModifyPlan(ctx context.Context, req fwresource.ModifyPlanRequest, resp *fwresource.ModifyPlanResponse) {
	m, ok := r.Resource.(fwresource.ResourceWithModifyPlan)
	if !ok {
		return
	}
	plan, replace, diags := resp.Plan, resp.RequiresReplace, resp.Diagnostics
	defer func() {
		if p := recover(); p != nil {
			resp.Plan, resp.RequiresReplace, resp.Diagnostics = plan, replace, diags
			r.log.Info("Skipped a ModifyPlan that needs the ClickHouse Cloud API", "resource", fmt.Sprintf("%T", r.Resource), "panic", fmt.Sprint(p))
		}
	}()
	m.ModifyPlan(ctx, req, resp)
}

// OfflineTerraformSetupBuilder returns the terraform.SetupFn used by the diff
// server. It ignores the ProviderConfig entirely: no credentials are read and
// no connection to ClickHouse Cloud is opened.
func OfflineTerraformSetupBuilder(frameworkProvider fwprovider.Provider, log logging.Logger) terraform.SetupFn {
	return func(context.Context, client.Client, resource.Managed) (terraform.Setup, error) {
		return terraform.Setup{
			FrameworkProvider: offlineProvider{Provider: frameworkProvider, log: log},
			Configuration:     map[string]any{},
		}, nil
	}
}
