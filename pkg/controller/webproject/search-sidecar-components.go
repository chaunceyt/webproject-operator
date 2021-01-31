package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// solrSearchContainerSpec - Solr search sidecar
// TODO:
// - Create service + ingress giving access to solr admin
// - Add support for StartupProbe.
// - Add config option for name of configmap containing solr config.
// - Add support for PVC
func solrSearchContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.SearchSidecar.Image,
		Name:  "search",
		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(8983),
			Name:          "search-port",
		}},
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	if cr.Spec.SearchSidecar.CronJob.Enabled {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      "search-cron-script",
			MountPath: "/opt/script",
		})
	}

	return container
}

// elasticSearchContainerSpec - ElasticSearch sidecar
// TODO:
// - Add support for PVC
// - Add support for StartupProbe.
// - Use for logging solution for webproject. (fluentbit sidecar + kibania)
func elasticSearchContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.SearchSidecar.Image,
		Name:  "search",
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: int32(9200),
				Name:          "reset-port",
			},
			{
				ContainerPort: int32(9300),
				Name:          "intra-node-port",
			},
		},
		Lifecycle: &corev1.Lifecycle{
			PostStart: &corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"sh", "-c", "chown -R 1000:1000 /usr/share/elasticsearch/data"},
				},
			},
		},
		/*StartupProbe: &corev1.Probe{
			InitialDelaySeconds: 5,
			PeriodSeconds:       2,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"sh", "-c", "chown -R 1000:1000 /usr/share/elasticsearch/data"},
				},
			},
		},*/
		Env: []corev1.EnvVar{
			{
				Name:  "discovery.type",
				Value: "single-node",
			},
			{
				Name:  "ES_JAVA_OPTS",
				Value: "-Xms512m -Xmx512m",
			},
		},
		VolumeMounts: getElasticSearchVolumeMounts(cr),
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: createBool(false),
			ReadOnlyRootFilesystem:   createBool(false),
			RunAsNonRoot:             createBool(false),
		},
	}

	if cr.Spec.SearchSidecar.CronJob.Enabled {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      "search-cron-script",
			MountPath: "/opt/script",
		})
	}

	return container
}
