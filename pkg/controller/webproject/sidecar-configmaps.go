package webproject

import (
	"bytes"
	"fmt"
	"html/template"

	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// databaseSidecarCronJobConfigMapForWebproject - configmap that contains the script code.
func (r *ReconcileWebproject) databaseSidecarCronJobConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "database-cron-script"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"script.sh": cr.Spec.DatabaseSidecar.CronJob.Script,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// redisConfigMapForWebproject - configmap that contains configuration for redis caching.
func (r *ReconcileWebproject) redisConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {
	// create template for redis.
	redisConfigTemplate := `loglevel notice
logfile /dev/stdout
timeout 300
tcp-keepalive 0
databases 1
maxmemory 2mb
maxmemory-policy allkeys-lru
`

	tmpl, err := template.New("redis").Parse(redisConfigTemplate)
	if err != nil {
		panic(err)
	}

	var tplOutput bytes.Buffer
	if err := tmpl.Execute(&tplOutput, cr); err != nil {
		panic(err)
	}

	redisConfigFileContent := tplOutput.String()

	if cr.Spec.CacheSidecar.RedisPassword != "" {
		redisConfigFileContent = fmt.Sprintf("%s\nrequirepass %s\n", redisConfigFileContent, cr.Spec.CacheSidecar.RedisPassword)
	}

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "redis-conf"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"redis.conf": redisConfigFileContent,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// searchSidecarCronJobConfigMapForWebproject - configmap that contains the script code.
func (r *ReconcileWebproject) searchSidecarCronJobConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "search-cron-script"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"script.sh": cr.Spec.SearchSidecar.CronJob.Script,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

// webContainerCronJobConfigMapForWebproject - configmap that contains the script code.
func (r *ReconcileWebproject) webContainerCronJobConfigMapForWebproject(cr *wp.WebProject) *corev1.ConfigMap {

	dep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "webcontainer-cron-script"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "config"),
		},
		Data: map[string]string{
			"script.sh": cr.Spec.WebContainer.CronJob.Script,
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}
