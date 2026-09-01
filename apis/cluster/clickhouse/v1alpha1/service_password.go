package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// SetPasswordSecretRef sets the PasswordSecretRef on the cluster-scoped
// Service ForProvider spec.
func (s *Service) SetPasswordSecretRef(ref *xpv2.SecretKeySelector) {
	s.Spec.ForProvider.PasswordSecretRef = ref
}
