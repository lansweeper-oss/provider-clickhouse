package v1alpha1

import (
	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
)

// SetPasswordSecretRef sets the PasswordSecretRef on the namespaced
// Service ForProvider spec.
func (s *Service) SetPasswordSecretRef(ref *v1.LocalSecretKeySelector) {
	s.Spec.ForProvider.PasswordSecretRef = ref
}
