package webproject

import (
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
