package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// databaseContainerSpec - database sidecar
// Add support for StartupProbe
func databaseContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.DatabaseSidecar.DatabaseImage,
		Name:  "database",

		Env: []corev1.EnvVar{
			{
				Name:  "MYSQL_ROOT_PASSWORD",
				Value: cr.Spec.DatabaseSidecar.DatabaseRootPassword,
			},
			{
				Name:  "MYSQL_DATABASE",
				Value: cr.Spec.DatabaseSidecar.DatabaseName,
			},
			{
				Name:  "MYSQL_USER",
				Value: cr.Spec.DatabaseSidecar.DatabaseUser,
			},
			{
				Name:  "MYSQL_PASSWORD",
				Value: cr.Spec.DatabaseSidecar.DatabaseUserPassword,
			},
		},

		Ports: []corev1.ContainerPort{{
			ContainerPort: 3306,
			Name:          "database",
		}},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "data-storage",
				MountPath: cr.Spec.DatabaseSidecar.DatabaseStoreMountPath,
			},
		},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	return container
}
