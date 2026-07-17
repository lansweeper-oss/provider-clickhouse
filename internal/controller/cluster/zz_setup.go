// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	clickpipe "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickhouse/clickpipe"
	role "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickhouse/role"
	service "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickhouse/service"
	cdcinfrastructure "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipe/cdcinfrastructure"
	reverseprivateendpoint "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipes/reverseprivateendpoint"
	reverseprivateendpointcustomprivatedns "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickpipes/reverseprivateendpointcustomprivatedns"
	settings "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/organization/settings"
	servicepostgres "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/postgres/service"
	providerconfig "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/providerconfig"
	assignment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/role/assignment"
	privateendpointsattachment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/privateendpointsattachment"
	scheduledscaling "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/scheduledscaling"
	transparentdataencryptionkeyassociation "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/transparentdataencryptionkeyassociation"
	upgradewindow "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/upgradewindow"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		clickpipe.Setup,
		role.Setup,
		service.Setup,
		cdcinfrastructure.Setup,
		reverseprivateendpoint.Setup,
		reverseprivateendpointcustomprivatedns.Setup,
		settings.Setup,
		servicepostgres.Setup,
		providerconfig.Setup,
		assignment.Setup,
		privateendpointsattachment.Setup,
		scheduledscaling.Setup,
		transparentdataencryptionkeyassociation.Setup,
		upgradewindow.Setup,
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
		clickpipe.SetupGated,
		role.SetupGated,
		service.SetupGated,
		cdcinfrastructure.SetupGated,
		reverseprivateendpoint.SetupGated,
		reverseprivateendpointcustomprivatedns.SetupGated,
		settings.SetupGated,
		servicepostgres.SetupGated,
		providerconfig.SetupGated,
		assignment.SetupGated,
		privateendpointsattachment.SetupGated,
		scheduledscaling.SetupGated,
		transparentdataencryptionkeyassociation.SetupGated,
		upgradewindow.SetupGated,
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
		clickpipe.SetupWebhookWithManager,
		role.SetupWebhookWithManager,
		service.SetupWebhookWithManager,
		cdcinfrastructure.SetupWebhookWithManager,
		reverseprivateendpoint.SetupWebhookWithManager,
		reverseprivateendpointcustomprivatedns.SetupWebhookWithManager,
		settings.SetupWebhookWithManager,
		servicepostgres.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		assignment.SetupWebhookWithManager,
		privateendpointsattachment.SetupWebhookWithManager,
		scheduledscaling.SetupWebhookWithManager,
		transparentdataencryptionkeyassociation.SetupWebhookWithManager,
		upgradewindow.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
