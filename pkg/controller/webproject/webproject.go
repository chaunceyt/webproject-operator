/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webproject

import (
	"context"
	"fmt"

	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
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

func (r *ReconcileWebproject) updateWebProjectStatus(cr *wp.WebProject) error {
	err := r.client.Status().Update(context.TODO(), cr)
	return err
}

// deploymentForWebproject - standard object to manage the webproject's pod.
func (r *ReconcileWebproject) deploymentForWebproject(cr *wp.WebProject) *appsv1.Deployment {
	matchlabels := map[string]string{
		"app.kubernetes.io/name": cr.Name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Spec.ReleaseName,
			Namespace:   cr.Namespace,
			Annotations: cr.Spec.DeploymentAnnotations,
			Labels:      webprojectlabels(cr, "deployment"),
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

	if cr.Spec.ImagePullSecrets.Enabled {
		deployment.Spec.Template.Spec.ImagePullSecrets = append(
			deployment.Spec.Template.Spec.ImagePullSecrets, corev1.LocalObjectReference{
				Name: cr.Spec.ImagePullSecrets.Secretname,
			},
		)
	}

	controllerutil.SetControllerReference(cr, deployment, r.scheme)
	return deployment

}

// backupServiceForWebproject - service responsible for exposing port 3306 for mysql|mariadb container in pod.
// This service is used by the backup cronjob.
// TODO: Determine if we should use network security policy to restrict access to this service.
func (r *ReconcileWebproject) backupServiceForWebproject(cr *wp.WebProject) *corev1.Service {
	matchlabels := map[string]string{
		"app.kubernetes.io/name": cr.Name,
	}

	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "backup-svc"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "service"),
		},

		Spec: corev1.ServiceSpec{
			Selector: matchlabels,

			Ports: []corev1.ServicePort{
				{
					Port: 3306, // externalPort
					TargetPort: intstr.IntOrString{
						Type:   Int,
						IntVal: 3306,
					}, // internalPort
					Protocol: "TCP",
					Name:     "backup-port",
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(cr, service, r.scheme)
	return service

}

func (r *ReconcileWebproject) solrServiceForWebproject(cr *wp.WebProject) *corev1.Service {
	matchlabels := map[string]string{
		"app.kubernetes.io/name": cr.Name,
	}

	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "solr-svc"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "service"),
		},

		Spec: corev1.ServiceSpec{
			Selector: matchlabels,

			Ports: []corev1.ServicePort{
				{
					Port: 8983,
					TargetPort: intstr.IntOrString{
						Type:   Int,
						IntVal: 8983,
					},
					Protocol: "TCP",
					Name:     "solr-port",
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(cr, service, r.scheme)
	return service

}

// serviceForWebproject - service responsible for exposing port 80 of the webcontainer in pod.
func (r *ReconcileWebproject) serviceForWebproject(cr *wp.WebProject) *corev1.Service {
	matchlabels := map[string]string{
		"app.kubernetes.io/name": cr.Name,
	}

	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "svc"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "service"),
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

// configMapForWebproject - configmap that contains the mysql|mariadb user and database variables.
func (r *ReconcileWebproject) envConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "env-config"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"MYSQL_USER":     cr.Spec.DatabaseSidecar.DatabaseUser,
			"MYSQL_DATABASE": cr.Spec.DatabaseSidecar.DatabaseName,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// secretForWebproject - secret that contains the password for mysql|mariadb database user.
func (r *ReconcileWebproject) secretForWebproject(cr *wp.WebProject) *corev1.Secret {

	dep := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "secret"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string][]byte{
			"MYSQL_PASSWORD": []byte(cr.Spec.DatabaseSidecar.DatabaseUserPassword),
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// commonConfigMapForWebproject - configmap that contains the script code for initcontainer.
func (r *ReconcileWebproject) webcontainerCronJobConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "web-container-cron"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"script.sh": cr.Spec.InitContainerScript,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// commonConfigMapForWebproject - configmap that contains the script code for initcontainer.
func (r *ReconcileWebproject) initContainerConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "init-container"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"init-container.sh": cr.Spec.InitContainerScript,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// commonConfigMapForWebproject - configmap that contains common environment variables defined for webproject.
func (r *ReconcileWebproject) commonConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "common-config"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: cr.Spec.CommonConfig,
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// ingressForWebproject - ingress object containing the domains for the webproject.
func (r *ReconcileWebproject) ingressForWebproject(cr *wp.WebProject) *networkingv1beta1.Ingress {

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

	subDomains := webprojectDomainNames(cr)
	ingressSpec := networkingv1beta1.IngressSpec{
		TLS: []networkingv1beta1.IngressTLS{
			{
				Hosts: subDomains,
			},
		},
	}

	ingress := &networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        workloadName(cr, "ingress"),
			Labels:      webprojectlabels(cr, "ingress"),
			Namespace:   cr.Namespace,
			Annotations: cr.Spec.IngressAnnotations,
		},
		Spec: ingressSpec,
	}

	for _, domain := range subDomains {
		ingress.Spec.Rules = append(
			ingress.Spec.Rules, networkingv1beta1.IngressRule{
				Host: domain,
				IngressRuleValue: networkingv1beta1.IngressRuleValue{
					HTTP: &networkingv1beta1.HTTPIngressRuleValue{
						Paths: ingressPaths,
					},
				},
			},
		)
	}

	controllerutil.SetControllerReference(cr, ingress, r.scheme)
	return ingress
}

// pvcForWebproject - persistent volume claim for static files.
// TODO: add support to create VolumeSnapshot from current pvc and use that PVC for webproject
func (r *ReconcileWebproject) pvcForWebproject(cr *wp.WebProject) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "files"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "storage"),
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
// TODO: add support to create VolumeSnapshot from current pvc and use that PVC for webproject
func (r *ReconcileWebproject) pvcForMysql(cr *wp.WebProject) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "data"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "storage"),
		},

		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},

			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.DatabaseSidecar.DatabaseStorageSize),
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.scheme)
	return pvc
}

// webProjectPodSpect - pod for webproject with multiple sidecars.
func webProjectPodSpec(cr *wp.WebProject) corev1.PodSpec {
	webpod := corev1.PodSpec{
		AutomountServiceAccountToken: createBool(false),
		Containers: []corev1.Container{
			webContainerSpec(cr),
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
		},
	}

	// TODO:
	// - append the initcontainer if initcontainer.enabled.
	// - append volume aws secret is enabled.

	if cr.Spec.AWSSecretName != "" {
		webpod.InitContainers = append(webpod.InitContainers, corev1.Container{
			Name:            "webdata",
			Image:           cr.Spec.CLISidecar.Image,
			Command:         []string{"bash", "-c", "/script/init-container.sh"},
			ImagePullPolicy: corev1.PullIfNotPresent,
			SecurityContext: &corev1.SecurityContext{
				RunAsNonRoot:             createBool(false),
				ReadOnlyRootFilesystem:   createBool(false),
				AllowPrivilegeEscalation: createBool(false),
			},
			Env: []corev1.EnvVar{
				{
					Name:  "RELEASE_NAME",
					Value: cr.Spec.ReleaseName,
				},
			},

			VolumeMounts: []corev1.VolumeMount{
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
			},
		})
		webpod.Volumes = append(webpod.Volumes, corev1.Volume{
			Name: "aws-credentials",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: cr.Spec.AWSSecretName,
				},
			},
		})
		webpod.Volumes = append(webpod.Volumes, corev1.Volume{})
	}

	// append database sidecar
	if cr.Spec.DatabaseSidecar.Enabled {
		webpod.Containers = append(webpod.Containers, databaseContainerSpec(cr))
		webpod.Volumes = append(webpod.Volumes, corev1.Volume{
			Name: "data-storage",
			VolumeSource: corev1.VolumeSource{

				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: workloadName(cr, "data"),
				},
			},
		})
	}

	// append cli sidecar
	if cr.Spec.CLISidecar.Enabled {
		webpod.Containers = append(webpod.Containers, cliContainerSpec(cr))
	}

	// append cache sidecar
	// cache engines supported: memcached and redis
	if cr.Spec.CacheSidecar.Enabled {
		if cr.Spec.CacheSidecar.Engine == "memcached" {
			webpod.Containers = append(webpod.Containers, memcachedCacheContainerSpec(cr))

		} else if cr.Spec.CacheSidecar.Engine == "redis" {
			webpod.Containers = append(webpod.Containers, redisCacheContainerSpec(cr))
		}
	}

	// append search sidecar
	// search engines supported: solr and elasticsearch
	if cr.Spec.SearchSidecar.Enabled {
		if cr.Spec.SearchSidecar.Engine == "es" {
			webpod.Containers = append(webpod.Containers, elasticSearchContainerSpec(cr))
		} else if cr.Spec.SearchSidecar.Engine == "solr" {
			webpod.Containers = append(webpod.Containers, solrSearchContainerSpec(cr))
		}

		webpod.Volumes = append(webpod.Volumes, corev1.Volume{
			Name: "search-data",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	}

	return webpod
}

// labels - labels used on all objects.
func webprojectlabels(cr *wp.WebProject, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      cr.Name,
		"app.kubernetes.io/part-of":   cr.Spec.ReleaseName,
		"app.kubernetes.io/component": component,
		"app.kubernetes.io/version":   cr.Spec.ReleaseName,
		"release":                     cr.Spec.ReleaseName,
		"provider":                    "webproject-operator",
	}
}

// webContainerSpec - primary contianer for webproject
// Add support for StartupProbe
func webContainerSpec(cr *wp.WebProject) corev1.Container {
	container := corev1.Container{
		Image: cr.Spec.WebContainer.Image,
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

		Ports: []corev1.ContainerPort{{
			ContainerPort: 80,
			Name:          "web-port",
		}},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: 5,
			PeriodSeconds:       2,
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{
						IntVal: 80,
					},
				},
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
		SecurityContext: &corev1.SecurityContext{
			RunAsNonRoot:             createBool(false),
			AllowPrivilegeEscalation: createBool(false),
		},
	}

	if cr.Spec.AWSSecretName != "" {
		container.EnvFrom = append(container.EnvFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: workloadName(cr, "aws-secret"),
				},
			},
		})
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

// workloadName
func workloadName(cr *wp.WebProject, workloadType string) string {
	return cr.Name + "-" + workloadType
}

func webprojectDomainNames(cr *wp.WebProject) []string {
	subDomains := []string{}
	domains := cr.Spec.IngressHosts

	for _, domain := range domains {
		str := fmt.Sprintf("release-%s-"+cr.Spec.ProjectDomainName, domain)
		subDomains = append(subDomains, str)
	}
	return subDomains
}

// createInt32 - helper function
func createInt32(x int32) *int32 {
	return &x
}

func createBool(x bool) *bool {
	return &x
}

func createInt64(i int64) *int64 {
	return &i
}
