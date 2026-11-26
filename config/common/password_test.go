package common

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/lansweeper-oss/provider-clickhouse/apis/cluster/clickhouse/v1alpha1"
)

func TestPasswordGeneratorManagementPolicies(t *testing.T) {
	cases := map[string]struct {
		policies  xpv2.ManagementPolicies
		wantWrite bool
	}{
		"ObserveOnly":              {policies: xpv2.ManagementPolicies{xpv2.ManagementActionObserve}},
		"ObserveLateInitialize":    {policies: xpv2.ManagementPolicies{xpv2.ManagementActionObserve, xpv2.ManagementActionLateInitialize}},
		"ObserveDelete":            {policies: xpv2.ManagementPolicies{xpv2.ManagementActionObserve, xpv2.ManagementActionDelete}},
		"All":                      {policies: xpv2.ManagementPolicies{xpv2.ManagementActionAll}, wantWrite: true},
		"ObserveCreate":            {policies: xpv2.ManagementPolicies{xpv2.ManagementActionObserve, xpv2.ManagementActionCreate}, wantWrite: true},
		"ObserveUpdate":            {policies: xpv2.ManagementPolicies{xpv2.ManagementActionObserve, xpv2.ManagementActionUpdate}, wantWrite: true},
		"EmptyMeansFullManagement": {policies: nil, wantWrite: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			svc := &v1alpha1.Service{}
			svc.SetName("svc")
			svc.SetManagementPolicies(tc.policies)
			svc.SetWriteConnectionSecretToReference(&xpv2.SecretReference{Name: "conn", Namespace: "ns"})

			calls := 0
			cl := &test.MockClient{
				// An existing secret with a password skips generation and only sets the ref.
				MockGet: func(_ context.Context, _ client.ObjectKey, obj client.Object) error {
					calls++
					obj.(*corev1.Secret).Data = map[string][]byte{passwordKey: []byte("pw")}
					return nil
				},
				MockUpdate: func(_ context.Context, _ client.Object, _ ...client.UpdateOption) error {
					calls++
					return nil
				},
			}

			init := PasswordGenerator("spec.forProvider.passwordSecretRef", "spec.writeConnectionSecretToRef")(cl)
			require.NoError(t, init.Initialize(context.Background(), svc))

			if !tc.wantWrite {
				require.Zero(t, calls, "client must not be called")
				require.Nil(t, svc.Spec.ForProvider.PasswordSecretRef)
				return
			}
			require.Equal(t, 2, calls, "expected Get and Update")
			require.Equal(t, &xpv2.SecretKeySelector{Name: "conn", Namespace: "ns", Key: passwordKey},
				svc.Spec.ForProvider.PasswordSecretRef)
		})
	}
}
