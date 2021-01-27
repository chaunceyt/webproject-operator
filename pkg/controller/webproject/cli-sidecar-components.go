package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// cliContainerSpec - cli sidecar
// Add support for StartupProbe
func cliContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.CLISidecar.Image,
		Name:  "cli",
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "webroot",
				MountPath: "/var/www",
			},
			{
				Name:      "files-storage",
				MountPath: cr.Spec.FileStorageMountPath,
			},
		},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			RunAsNonRoot:             createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
		},
	}

	if cr.Spec.DatabaseSidecar.Enabled {
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "DB_HOST",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		})
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "MYSQL_DATABASE",
			Value: cr.Spec.DatabaseSidecar.DatabaseName,
		})

		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "MYSQL_USER",
			Value: cr.Spec.DatabaseSidecar.DatabaseUser,
		})
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "MYSQL_PASSWORD",
			Value: cr.Spec.DatabaseSidecar.DatabaseUserPassword,
		})
	}

	// append cache env var.
	if cr.Spec.CacheSidecar.Enabled {
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "CACHE_HOST",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		})
	}

	// append search env var.
	if cr.Spec.SearchSidecar.Enabled {
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "SEARCH_HOST",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		})
	}

	return container
}
