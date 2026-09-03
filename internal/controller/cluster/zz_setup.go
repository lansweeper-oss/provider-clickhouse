// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	service "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickhouse/service"
	cdcinfrastructure "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipe/cdcinfrastructure"
	clickpipe "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipe/clickpipe"
	reverseprivateendpoint "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipes/reverseprivateendpoint"
	reverseprivateendpointcustomprivatedns "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipes/reverseprivateendpointcustomprivatedns"
	alert "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/alert"
	connection "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/connection"
	dashboard "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/dashboard"
	role "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/role"
	savedsearch "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/savedsearch"
	source "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/source"
	team "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/team"
	teammember "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/teammember"
	webhook "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickstack/webhook"
	roleiam "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/iam/role"
	settings "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/organization/settings"
	servicepostgres "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/postgres/service"
	providerconfig "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/providerconfig"
	assignment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/role/assignment"
	privateendpointsattachment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/privateendpointsattachment"
	scheduledscaling "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/scheduledscaling"
	transparentdataencryptionkeyassociation "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/transparentdataencryptionkeyassociation"
	upgradewindow "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/upgradewindow"
	attachment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/udf/attachment"
	udf "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/udf/udf"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		service.Setup,
		cdcinfrastructure.Setup,
		clickpipe.Setup,
		reverseprivateendpoint.Setup,
		reverseprivateendpointcustomprivatedns.Setup,
		alert.Setup,
		connection.Setup,
		dashboard.Setup,
		role.Setup,
		savedsearch.Setup,
		source.Setup,
		team.Setup,
		teammember.Setup,
		webhook.Setup,
		roleiam.Setup,
		settings.Setup,
		servicepostgres.Setup,
		providerconfig.Setup,
		assignment.Setup,
		privateendpointsattachment.Setup,
		scheduledscaling.Setup,
		transparentdataencryptionkeyassociation.Setup,
		upgradewindow.Setup,
		attachment.Setup,
		udf.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		service.SetupGated,
		cdcinfrastructure.SetupGated,
		clickpipe.SetupGated,
		reverseprivateendpoint.SetupGated,
		reverseprivateendpointcustomprivatedns.SetupGated,
		alert.SetupGated,
		connection.SetupGated,
		dashboard.SetupGated,
		role.SetupGated,
		savedsearch.SetupGated,
		source.SetupGated,
		team.SetupGated,
		teammember.SetupGated,
		webhook.SetupGated,
		roleiam.SetupGated,
		settings.SetupGated,
		servicepostgres.SetupGated,
		providerconfig.SetupGated,
		assignment.SetupGated,
		privateendpointsattachment.SetupGated,
		scheduledscaling.SetupGated,
		transparentdataencryptionkeyassociation.SetupGated,
		upgradewindow.SetupGated,
		attachment.SetupGated,
		udf.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		service.SetupWebhookWithManager,
		cdcinfrastructure.SetupWebhookWithManager,
		clickpipe.SetupWebhookWithManager,
		reverseprivateendpoint.SetupWebhookWithManager,
		reverseprivateendpointcustomprivatedns.SetupWebhookWithManager,
		alert.SetupWebhookWithManager,
		connection.SetupWebhookWithManager,
		dashboard.SetupWebhookWithManager,
		role.SetupWebhookWithManager,
		savedsearch.SetupWebhookWithManager,
		source.SetupWebhookWithManager,
		team.SetupWebhookWithManager,
		teammember.SetupWebhookWithManager,
		webhook.SetupWebhookWithManager,
		roleiam.SetupWebhookWithManager,
		settings.SetupWebhookWithManager,
		servicepostgres.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		assignment.SetupWebhookWithManager,
		privateendpointsattachment.SetupWebhookWithManager,
		scheduledscaling.SetupWebhookWithManager,
		transparentdataencryptionkeyassociation.SetupWebhookWithManager,
		upgradewindow.SetupWebhookWithManager,
		attachment.SetupWebhookWithManager,
		udf.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
