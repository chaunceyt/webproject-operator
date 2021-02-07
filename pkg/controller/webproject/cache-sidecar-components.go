package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// memcachedCacheContainerSpec - cache sidecar (memcached or redis)
// TODO
// - Make immutable. (ReadOnlyRootFilesystem: true)
func memcachedCacheContainerSpec(cr *wp.WebProject) corev1.Container {
	image := "memcached:1.6.9"
	container := corev1.Container{
		Name:  "cache",
		Image: image,
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(11211),
			Name:          "memcached",
		}},
		Command:   []string{"memcached", "-m", "128", "-o", "modern", "-vv"},
		Resources: cr.Spec.CacheSidecar.Resources,
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(true),
			RunAsNonRoot:             createBool(true),
			RunAsUser:                createInt64(1000),
			RunAsGroup:               createInt64(1000),
		},
		Env: []corev1.EnvVar{
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
			{
				Name: "POD_ID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.uid",
					},
				},
			},
			{
				Name: "POD_NAMESPACE",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.namespace",
					},
				},
			},
		},
	}

	return container
}

// redisCacheContainerSpec - cache sidecar (memcached or redis)
// TODO
// - Make immutable. (ReadOnlyRootFilesystem: true) [WIP]
func redisCacheContainerSpec(cr *wp.WebProject) corev1.Container {
	image := "redis:6.0.10-alpine"
	command := []string{"redis-server", "/opt/redis/redis.conf"}

	passwordParam := ""
	password := cr.Spec.CacheSidecar.RedisPassword
	if password != "" {
		passwordParam = " -a " + password
	}
	probeCmd := []string{"sh", "-c", "redis-cli", passwordParam, "ping"}

	container := corev1.Container{
		Name:         "cache",
		Image:        image,
		Command:      command,
		Resources:    cr.Spec.CacheSidecar.Resources,
		VolumeMounts: getRedisVolumeMounts(cr),
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(6379),
			Name:          "redis",
		}},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(true),
			RunAsNonRoot:             createBool(true),
			RunAsUser:                createInt64(1000),
			RunAsGroup:               createInt64(1000),
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
		Env: []corev1.EnvVar{
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
			{
				Name: "POD_ID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.uid",
					},
				},
			},
			{
				Name: "POD_NAMESPACE",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.namespace",
					},
				},
			},
		},
	}

	return container
}
