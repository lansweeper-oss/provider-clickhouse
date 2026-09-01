package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// SetPasswordSecretRef sets the PasswordSecretRef on the namespaced
// Service ForProvider spec.
func (s *Service) SetPasswordSecretRef(ref *xpv2.LocalSecretKeySelector) {
	s.Spec.ForProvider.PasswordSecretRef = ref
}
