// Package common holds shared resource configuration helpers.
package common

import (
	"context"
	"fmt"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/password"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
)

const passwordKey = "password"

// ClusterPasswordSecretRefSetter is implemented by cluster-scoped managed
// resources that expose a *xpv2.SecretKeySelector PasswordSecretRef.
type ClusterPasswordSecretRefSetter interface {
	SetPasswordSecretRef(ref *xpv2.SecretKeySelector)
}

// NamespacedPasswordSecretRefSetter is implemented by namespaced managed
// resources that expose a *xpv2.LocalSecretKeySelector PasswordSecretRef.
type NamespacedPasswordSecretRefSetter interface {
	SetPasswordSecretRef(ref *xpv2.LocalSecretKeySelector)
}

// PasswordGenerator returns a NewInitializerFn:
//   - BYOP: if byopSecretRefPath is already set, no-op.
//   - Auto-generate: if writeConnectionSecretPath is set, generate a password,
//     write it to that secret under "password", and point passwordSecretRef at it.
func PasswordGenerator(byopSecretRefPath, writeConnectionSecretPath string) config.NewInitializerFn {
	return func(cl client.Client) managed.Initializer {
		return managed.InitializerFn(func(ctx context.Context, mg xpresource.Managed) error {
			paved, err := fieldpath.PaveObject(mg)
			if err != nil {
				return fmt.Errorf("cannot pave object: %w", err)
			}
			byop, err := checkBYOP(paved, byopSecretRefPath)
			if err != nil {
				return err
			}
			if byop {
				return nil
			}
			name, ns, err := resolveConnRef(paved, writeConnectionSecretPath, mg.GetNamespace())
			if err != nil || name == "" {
				return err
			}
			return reconcilePassword(ctx, cl, mg, name, ns)
		})
	}
}

func checkBYOP(paved *fieldpath.Paved, path string) (bool, error) {
	sel := map[string]any{}
	if err := paved.GetValueInto(path, &sel); err != nil {
		if fieldpath.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("cannot read %s: %w", path, err)
	}
	name, _ := sel["name"].(string)
	return name != "", nil
}

func resolveConnRef(paved *fieldpath.Paved, path, defaultNS string) (name, ns string, err error) {
	ref := map[string]any{}
	if err := paved.GetValueInto(path, &ref); err != nil {
		if fieldpath.IsNotFound(err) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("cannot read %s: %w", path, err)
	}
	name, _ = ref["name"].(string)
	ns, _ = ref["namespace"].(string)
	if ns == "" {
		ns = defaultNS
	}
	return name, ns, nil
}

func reconcilePassword(ctx context.Context, cl client.Client, mg xpresource.Managed, name, ns string) error {
	s := &corev1.Secret{}
	getErr := cl.Get(ctx, types.NamespacedName{Namespace: ns, Name: name}, s)
	if xpresource.IgnoreNotFound(getErr) != nil {
		return fmt.Errorf("cannot get connection secret: %w", getErr)
	}
	if getErr == nil && len(s.Data[passwordKey]) != 0 {
		return setPasswordSecretRef(ctx, cl, mg, name, ns)
	}
	return generateAndApply(ctx, cl, mg, s, name, ns)
}

func generateAndApply(ctx context.Context, cl client.Client, mg xpresource.Managed, s *corev1.Secret, name, ns string) error {
	pw, err := password.Generate()
	if err != nil {
		return fmt.Errorf("cannot generate password: %w", err)
	}
	s.SetName(name)
	s.SetNamespace(ns)
	s.Type = xpresource.SecretTypeConnection
	if !meta.WasCreated(s) {
		meta.AddOwnerReference(s, meta.AsController(meta.TypedReferenceTo(mg, mg.GetObjectKind().GroupVersionKind())))
	}
	if s.Data == nil {
		s.Data = make(map[string][]byte, 1)
	}
	s.Data[passwordKey] = []byte(pw)
	if err := xpresource.NewAPIPatchingApplicator(cl).Apply(ctx, s); err != nil {
		return fmt.Errorf("cannot apply password secret: %w", err)
	}
	return setPasswordSecretRef(ctx, cl, mg, name, ns)
}

func setPasswordSecretRef(ctx context.Context, cl client.Client, mg xpresource.Managed, name, namespace string) error {
	switch setter := mg.(type) {
	case ClusterPasswordSecretRefSetter:
		setter.SetPasswordSecretRef(&xpv2.SecretKeySelector{
			SecretReference: xpv2.SecretReference{Name: name, Namespace: namespace},
			Key:             passwordKey,
		})
	case NamespacedPasswordSecretRefSetter:
		setter.SetPasswordSecretRef(&xpv2.LocalSecretKeySelector{
			LocalSecretReference: xpv2.LocalSecretReference{Name: name},
			Key:                  passwordKey,
		})
	default:
		return nil
	}
	if err := cl.Update(ctx, mg); err != nil {
		return fmt.Errorf("cannot update managed resource with password secret ref: %w", err)
	}
	return nil
}
