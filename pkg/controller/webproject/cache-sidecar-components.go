package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// memcachedCacheContainerSpec - cache sidecar (memcached or redis)
// TODO
// - Make immutable. (ReadOnlyRootFilesystem: true)
func memcachedCacheContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.CacheSidecar.Image,
		Name:  "cache",
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(11211),
			Name:          "memcached",
		}},
		Command: []string{"memcached", "-m", "128", "-o", "modern", "-vv"},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	return container
}

// redisCacheContainerSpec - cache sidecar (memcached or redis)
// TODO
// - Make immutable. (ReadOnlyRootFilesystem: true)
func redisCacheContainerSpec(cr *wp.WebProject) corev1.Container {
	redisCommand := []string{"redis-server", "/opt/redis/redis.conf"}

	password := cr.Spec.CacheSidecar.RedisPassword
	probeCmd := []string{"sh", "-c", "redis-cli -a " + password + " ping"}

	container := corev1.Container{
		Name:         "cache",
		Image:        cr.Spec.CacheSidecar.Image,
		Command:      redisCommand,
		VolumeMounts: getRedisVolumeMounts(cr),
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(6379),
			Name:          "redis",
		}},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: 30,
			TimeoutSeconds:      5,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: probeCmd,
				},
			},
		},
		LivenessProbe: &corev1.Probe{
			InitialDelaySeconds: 30,
			TimeoutSeconds:      5,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: probeCmd,
				},
			},
		},
	}

	return container
}
