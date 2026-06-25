// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	service "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/clickhouse/service"
	settings "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/organization/settings"
	providerconfig "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/providerconfig"
	privateendpointsattachment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/service/privateendpointsattachment"
	transparentdataencryptionkeyassociation "github.com/lansweeper-oss/provider-clickhouse/internal/controller/namespaced/service/transparentdataencryptionkeyassociation"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		service.Setup,
		settings.Setup,
		providerconfig.Setup,
		privateendpointsattachment.Setup,
		transparentdataencryptionkeyassociation.Setup,
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
		settings.SetupGated,
		providerconfig.SetupGated,
		privateendpointsattachment.SetupGated,
		transparentdataencryptionkeyassociation.SetupGated,
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
		settings.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		privateendpointsattachment.SetupWebhookWithManager,
		transparentdataencryptionkeyassociation.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
