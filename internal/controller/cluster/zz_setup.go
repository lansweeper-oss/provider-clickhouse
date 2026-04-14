// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	service "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/clickhouse/service"
	settings "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/organization/settings"
	providerconfig "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/providerconfig"
	privateendpointsattachment "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/privateendpointsattachment"
	transparentdataencryptionkeyassociation "github.com/lansweeper-oss/provider-clickhouse/internal/controller/cluster/service/transparentdataencryptionkeyassociation"
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
