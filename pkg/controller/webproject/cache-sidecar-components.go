package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// memcachedCacheContainerSpec - cache sidecar (memcached or redis)
func memcachedCacheContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.CacheSidecar.Image,
		Name:  "cache",
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(11211),
			Name:          "memcached",
		}},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	return container
}

// redisCacheContainerSpec - cache sidecar (memcached or redis)
// Add support for StartupProbe
func redisCacheContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.CacheSidecar.Image,
		Name:  "cache",
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(6379),
			Name:          "redis",
		}},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	return container
}
