package v1alpha1

import (
	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
)

// SetPasswordSecretRef sets the PasswordSecretRef on the cluster-scoped
// Service ForProvider spec.
func (s *Service) SetPasswordSecretRef(ref *v1.SecretKeySelector) {
	s.Spec.ForProvider.PasswordSecretRef = ref
}
