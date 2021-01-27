package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// solrSearchContainerSpec - Solr search sidecar
// TODO: create service + ingress giving access to solr admin
// Add configmap to bring in config.
// Add support for StartupProbe.
// Add support for PVC
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

	return container
}

// elasticSearchContainerSpec - ElasticSearch sidecar
// Add support for StartupProbe.
// TODO: Add support for PVC
// Use for logging solution for webproject. fluentbit sidecar + kibania
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
		StartupProbe: &corev1.Probe{
			InitialDelaySeconds: 5,
			PeriodSeconds:       2,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"sh", "-c", "chown -R 1000:1000 /usr/share/elasticsearch/data"},
				},
			},
		},
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
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "search-data",
				MountPath: "/usr/share/elasticsearch/data",
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
