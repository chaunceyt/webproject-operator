package webproject

import (
	wpv1 "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	Int intstr.Type = iota
)

func (r *ReconcileWebproject) deploymentForWebproject(cr *wpv1.WebProject) *appsv1.Deployment {
	matchlabels := map[string]string{
		"app.kubernetes.io/name": cr.Name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ReleaseName,
			Namespace: cr.Namespace,
			Labels:    labels(cr, "deployment"),
		},

		Spec: appsv1.DeploymentSpec{
			// Only wait 5 min (instead of 10min) for failed deployment.
			ProgressDeadlineSeconds: createInt32(300),
			Selector: &metav1.LabelSelector{
				MatchLabels: matchlabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchlabels,
				},
				Spec: webProjectPodSpec(cr),
			},
		},
	}

	controllerutil.SetControllerReference(cr, deployment, r.scheme)
	return deployment

}

func (r *ReconcileWebproject) serviceForWebproject(cr *wpv1.WebProject) *corev1.Service {
	matchlabels := map[string]string{
		"app.kubernetes.io/name": cr.Name,
	}

	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "svc"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "service"),
		},

		Spec: corev1.ServiceSpec{
			Selector: matchlabels,

			Ports: []corev1.ServicePort{
				{
					Port: 80, // externalPort
					TargetPort: intstr.IntOrString{
						Type:   Int,
						IntVal: 80,
					}, // internalPort
					Protocol: "TCP",
					Name:     "port",
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(cr, service, r.scheme)
	return service

}

// configMapForWebproject returns a webproject configmap object
func (r *ReconcileWebproject) envConfigMapForWebproject(cr *wpv1.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "env-config"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "config"),
		},
		Data: map[string]string{
			"MYSQL_USER":     cr.Spec.DatabaseUser,
			"MYSQL_DATABASE": cr.Spec.DatabaseName,
		},
	}
	// Set Operator instance as the owner and controller
	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// configMapForWebproject returns a webproject configmap object
func (r *ReconcileWebproject) secretForWebproject(cr *wpv1.WebProject) *corev1.Secret {

	dep := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "secret"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "config"),
		},
		Data: map[string][]byte{
			"MYSQL_PASSWORD": []byte(cr.Spec.DatabaseUserPassword),
		},
	}
	// Set Operator instance as the owner and controller
	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// awsSecretForWebproject returns a webproject configmap object
func (r *ReconcileWebproject) awsSecretForWebproject(cr *wpv1.WebProject) *corev1.Secret {

	dep := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "aws-secret"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "config"),
		},
		Data: map[string][]byte{
			"AWS_ACCESS_KEY_ID":     []byte("changeme"),
			"AWS_SECRET_ACCESS_KEY": []byte("changeme"),
			"AWS_DEFAULT_REGION":    []byte("changeme"),
			"AWS_BUCKET":            []byte("changeme"),
		},
	}
	// Set Operator instance as the owner and controller
	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// commonConfigMapForWebproject returns a webproject configmap object
func (r *ReconcileWebproject) commonConfigMapForWebproject(cr *wpv1.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "common-config"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "config"),
		},
		Data: map[string]string{
			"BUILD_ID":               "changeme",
			"DOCROOT":                "changeme",
			"PROJECT_ENV":            "changeme",
			"CI":                     "true",
			"PHP_MAX_EXECUTION_TIME": "changeme",
			"PHP_MEMORY_LIMIT":       "changeme",
		},
	}
	// Set Operator instance as the owner and controller
	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// serviceForWebproject returns a webproject Ingress object
func (r *ReconcileWebproject) ingressForWebproject(cr *wpv1.WebProject) *networkingv1beta1.Ingress {

	ingressPaths := []networkingv1beta1.HTTPIngressPath{
		networkingv1beta1.HTTPIngressPath{
			Path: "/",
			Backend: networkingv1beta1.IngressBackend{
				ServiceName: workloadName(cr, "svc"),
				ServicePort: intstr.IntOrString{
					Type:   Int,
					IntVal: 80,
				},
			},
		},
	}
	domainname := webprojectDomainName(cr)
	domain := []string{domainname}
	ingressSpec := networkingv1beta1.IngressSpec{
		TLS: []networkingv1beta1.IngressTLS{
			{
				Hosts: domain,
			},
		},
		Rules: []networkingv1beta1.IngressRule{
			{
				Host: domainname,
				IngressRuleValue: networkingv1beta1.IngressRuleValue{
					HTTP: &networkingv1beta1.HTTPIngressRuleValue{
						Paths: ingressPaths,
					},
				},
			},
		},
	}
	// Allow webapp to handle ssl redirects - nginx.ingress.kubernetes.io/ssl-redirect: "false"
	// Add auth - nginx.ingress.kubernetes.io/auth-url: https://auth.domain.com/prod/auth
	// Add signin - nginx.ingress.kubernetes.io/auth-signin: https://auth.domain.com/prod/signin
	// Add support for rewriting of target - "nginx.ingress.kubernetes.io/rewrite-target":    "/$2",
	// Add nginx.ingress.kubernetes.io/auth-tls-verify-client: "off" and nginx.ingress.kubernetes.io/backend-protocol: HTTPS
	// if the project is using gatsby custom certs.
	ingress := &networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "ingress"),
			Labels:    labels(cr, "ingress"),
			Namespace: cr.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                   "nginx",
				"nginx.ingress.kubernetes.io/proxy-body-size":   "0",
				"nginx.ingress.kubernetes.io/proxy-buffer-size": "16k",
			},
		},
		Spec: ingressSpec,
	}

	// Set Operator instance as the owner and controller
	controllerutil.SetControllerReference(cr, ingress, r.scheme)
	return ingress
}

// pvcForWebproject - persistent volume claim for static files.
func (r *ReconcileWebproject) pvcForWebproject(cr *wpv1.WebProject) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "files"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "storage"),
		},

		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},

			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.FileStorageSize),
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.scheme)
	return pvc

}

// pvcForMysql - persistent volume claim for mysql|mariadb data path /var/lib/mysql
func (r *ReconcileWebproject) pvcForMysql(cr *wpv1.WebProject) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "data"),
			Namespace: cr.Namespace,
			Labels:    labels(cr, "storage"),
		},

		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},

			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.DatabaseStorageSize),
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.scheme)
	return pvc

}

// webProjectPodSpect - pod for webproject with multiple sidecars.
func webProjectPodSpec(cr *wpv1.WebProject) corev1.PodSpec {
	return corev1.PodSpec{
		AutomountServiceAccountToken: createBool(false),
		Containers: []corev1.Container{
			webContainerSpec(cr),
			cliContainerSpec(cr),
			databaseContainerSpec(cr),
			cacheContainerSpec(cr),
		},

		Volumes: []corev1.Volume{

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
				Name: "data-storage",
				VolumeSource: corev1.VolumeSource{

					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: workloadName(cr, "data"),
					},
				},
			},
		},
	}
}

// labels - labels used on all objects.
func labels(cr *wpv1.WebProject, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      cr.Name,
		"app.kubernetes.io/part-of":   cr.Spec.ReleaseName,
		"app.kubernetes.io/component": component,
		"app.kubernetes.io/version":   cr.Spec.ReleaseName,
		"release":                     cr.Spec.ReleaseName,
	}
}

// webContainerSpec - primary contianer for webproject
func webContainerSpec(cr *wpv1.WebProject) corev1.Container {
	return corev1.Container{
		Image: cr.Spec.WebImage,
		Name:  "web",
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "env-config"),
					},
				},
			},
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "common-config"),
					},
				},
			},
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "secret"),
					},
				},
			},
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: workloadName(cr, "aws-secret"),
					},
				},
			},
		},
		Env: []corev1.EnvVar{
			{
				Name: "DB_HOST",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
			{
				Name: "CACHE_HOST",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
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

		Ports: []corev1.ContainerPort{{
			ContainerPort: 80,
			Name:          "web-port",
		}},

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
	}
}

// cliContainerSpec - cli sidecar
func cliContainerSpec(cr *wpv1.WebProject) corev1.Container {
	return corev1.Container{
		Image: "outrigger/cli:2-php7.3",
		Name:  "cli",

		Env: []corev1.EnvVar{
			{
				Name: "DB_HOST",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
			{
				Name: "CACHE_HOST",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
			{
				Name:  "MYSQL_DATABASE",
				Value: cr.Spec.DatabaseName,
			},
			{
				Name:  "MYSQL_USER",
				Value: cr.Spec.DatabaseUser,
			},
			{
				Name:  "MYSQL_PASSWORD",
				Value: cr.Spec.DatabaseUserPassword,
			},
		},

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
	}
}

// databaseContainerSpec - database sidecar
func databaseContainerSpec(cr *wpv1.WebProject) corev1.Container {
	return corev1.Container{
		Image: cr.Spec.DatabaseImage,
		Name:  "database",

		Env: []corev1.EnvVar{
			{
				Name:  "MYSQL_ROOT_PASSWORD",
				Value: cr.Spec.DatabaseRootPassword,
			},
			{
				Name:  "MYSQL_DATABASE",
				Value: cr.Spec.DatabaseName,
			},
			{
				Name:  "MYSQL_USER",
				Value: cr.Spec.DatabaseUser,
			},
			{
				Name:  "MYSQL_PASSWORD",
				Value: cr.Spec.DatabaseUserPassword,
			},
		},

		Ports: []corev1.ContainerPort{{
			ContainerPort: 3306,
			Name:          "database",
		}},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "data-storage",
				MountPath: cr.Spec.DatabaseStoreMountPath,
			},
		},
	}
}

// cacheContainerSpec - cache sidecar (memcached or redis)
func cacheContainerSpec(cr *wpv1.WebProject) corev1.Container {
	return corev1.Container{
		Image: cr.Spec.CacheImage,
		Name:  "cache",

		Ports: []corev1.ContainerPort{{
			ContainerPort: int32(cr.Spec.CachePort),
			Name:          "cache-port",
		}},
	}
}

// workloadName
func workloadName(cr *wpv1.WebProject, workloadType string) string {
	return cr.Name + "-" + workloadType
}

func webprojectDomainName(cr *wpv1.WebProject) string {
	return cr.Spec.ReleaseName + "." + cr.Spec.ProjectDomainName
}

// createInt32 - helper function
func createInt32(x int32) *int32 {
	return &x
}

func createBool(x bool) *bool {
	return &x
}
