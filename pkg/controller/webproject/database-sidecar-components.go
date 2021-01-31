package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// databaseContainerSpec - database sidecar
// pvcForMysql is responsible for the PVC data-storage
// TODO:
// - Add support for StartupProbe
// - Implement backup solution.
/*
Backup solution

Storage options
- pvc
- AWS s3
- GCP bucket

pvc

- 2gb PVC default
- nginx-alpine pod with pvc to backup-volume mounted under docroot.
- basic auth to download backup
- use cronjob to run mysqldump of database
- use kubectl to copy database dump file to pod with backup-volume pvc
*/
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

		VolumeMounts: getDatabaseSidecarVolumeMounts(cr),

		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	return container
}
