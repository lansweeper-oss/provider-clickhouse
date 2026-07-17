// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	service "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/clickhouse/service"
	cdcinfrastructure "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/clickpipe/cdcinfrastructure"
	clickpipe "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/clickpipe/clickpipe"
	reverseprivateendpoint "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/clickpipes/reverseprivateendpoint"
	reverseprivateendpointcustomprivatedns "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/clickpipes/reverseprivateendpointcustomprivatedns"
	role "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/iam/role"
	settings "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/organization/settings"
	servicepostgres "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/postgres/service"
	providerconfig "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/providerconfig"
	assignment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/role/assignment"
	privateendpointsattachment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/service/privateendpointsattachment"
	scheduledscaling "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/service/scheduledscaling"
	transparentdataencryptionkeyassociation "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/service/transparentdataencryptionkeyassociation"
	upgradewindow "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/service/upgradewindow"
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
		role.Setup,
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
		service.SetupGated,
		cdcinfrastructure.SetupGated,
		clickpipe.SetupGated,
		reverseprivateendpoint.SetupGated,
		reverseprivateendpointcustomprivatedns.SetupGated,
		role.SetupGated,
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
		service.SetupWebhookWithManager,
		cdcinfrastructure.SetupWebhookWithManager,
		clickpipe.SetupWebhookWithManager,
		reverseprivateendpoint.SetupWebhookWithManager,
		reverseprivateendpointcustomprivatedns.SetupWebhookWithManager,
		role.SetupWebhookWithManager,
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
