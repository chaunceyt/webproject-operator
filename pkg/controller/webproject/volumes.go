package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// getWebProjectVolumes - volumes for web container.
func getWebProjectVolumes(cr *wp.WebProject) []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "webroot",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "files-storage",
			VolumeSource: corev1.VolumeSource{

				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: workloadName(cr, "files"),
				},
			},
		},
		{
			Name: "log-dir",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	if cr.Spec.DatabaseSidecar.Enabled {
		volumes = append(volumes, corev1.Volume{
			Name: "data-storage",
			VolumeSource: corev1.VolumeSource{

				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: workloadName(cr, "data"),
				},
			},
		})
	}

	if cr.Spec.SearchSidecar.Enabled {
		volumes = append(volumes, corev1.Volume{
			Name: "search-data",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	}

	if cr.Spec.WebContainer.CronJob.Enabled {
		volumes = append(volumes, corev1.Volume{
			Name: "webcontainer-cron-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "webcontainer-cron-script"),
					},
					DefaultMode: createInt32(0744),
				},
			},
		})
	}

	if cr.Spec.SearchSidecar.CronJob.Enabled {
		volumes = append(volumes, corev1.Volume{
			Name: "search-cron-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "search-cron-script"),
					},
					DefaultMode: createInt32(0744),
				},
			},
		})
	}

	if cr.Spec.DatabaseSidecar.CronJob.Enabled {
		volumes = append(volumes, corev1.Volume{
			Name: "database-cron-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "database-cron-script"),
					},
					DefaultMode: createInt32(0744),
				},
			},
		})
	}

	if cr.Spec.CacheSidecar.Engine == "redis" {
		volumes = append(volumes, corev1.Volume{
			Name: "redis-conf",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "redis-conf"),
					},
					DefaultMode: createInt32(0644),
				},
			},
		})

	}

	return volumes
}

// getWebProjectVolumeMounts - volumeMounts for web container.
func getWebProjectVolumeMounts(cr *wp.WebProject) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "webroot",
			MountPath: "/var/www",
		},
		{
			Name:      "files-storage",
			MountPath: cr.Spec.FileStorageMountPath,
		},
	}

	if cr.Spec.WebContainer.CronJob.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "webcontainer-cron-script",
			MountPath: "/opt/script",
		})
	}

	return volumeMounts
}

// getRedisVolumeMounts - volumeMounts for redis container.
func getRedisVolumeMounts(cr *wp.WebProject) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "redis-conf",
			MountPath: "/opt/redis/redis.conf",
			SubPath:   "redis.conf",
		},
		{
			Name:      "log-dir",
			MountPath: "/var/log/webproject",
		},
	}
	return volumeMounts
}

// getSolrVolumeMounts - volumeMounts for solr container.
func getSolrVolumeMounts(cr *wp.WebProject) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}
	return volumeMounts
}

// getElasticSearchVolumeMounts - volumeMount for elasticsearch container.
func getElasticSearchVolumeMounts(cr *wp.WebProject) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "search-data",
			MountPath: "/usr/share/elasticsearch/data",
		},
	}
	return volumeMounts
}

// getDatabaseSidecarVolumeMounts - volumeMounts for database container.
func getDatabaseSidecarVolumeMounts(cr *wp.WebProject) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "data-storage",
			MountPath: cr.Spec.DatabaseSidecar.DatabaseStoreMountPath,
		},
	}

	if cr.Spec.DatabaseSidecar.CronJob.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "database-cron-script",
			MountPath: "/opt/script",
		})
	}

	return volumeMounts
}

// getInitContainerVolumeMounts - volumeMounts for initcontainer
func getInitContainerVolumeMounts(cr *wp.WebProject) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "webroot",
			MountPath: "/data",
		},
		{
			Name:      "files-storage",
			MountPath: "/cmsfiles",
		},
		{
			Name:      "aws-credentials",
			MountPath: "/aws",
		},
		{
			Name:      "init-container",
			MountPath: "/script/init-container.sh",
			SubPath:   "init-container.sh",
		},
	}

	return volumeMounts
}
