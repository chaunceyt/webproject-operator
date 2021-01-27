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

	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_webproject")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Webproject Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWebproject{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("webproject-operator", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Webproject
	err = c.Watch(&source.Kind{Type: &wp.WebProject{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to resource Deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	// Watch for changes to resource ConfigMap
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	// Watch for changes to resource Secret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	// Watch for changes to resource Ingress
	err = c.Watch(&source.Kind{Type: &networkingv1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	// Watch for changes to resource Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	// Watch for changes to resource PersistentVolumeClaim
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	// Watch for changes to resource CronJob
	err = c.Watch(&source.Kind{Type: &v1beta1.CronJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileWebproject implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileWebproject{}

// ReconcileWebproject reconciles a Webproject object
type ReconcileWebproject struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Webproject object and makes changes based on the state read
// and what is in the Webproject.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWebproject) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// reqLogger.Info("Reconciling WebProject")

	// Fetch the Webproject instance
	webproject := &wp.WebProject{}
	err := r.client.Get(context.TODO(), request.NamespacedName, webproject)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	// ===== WEBPROJECT =====

	// ensurePVC - ensure the persitentvolumeclaim for /var/lib/mysql is managed.
	// if the database name is empty do not manage mysql PVC
	if webproject.Spec.DatabaseSidecar.Enabled {
		result, err = r.ensurePVC(request, webproject, r.pvcForMysql(webproject))
		if result != nil {
			return *result, err
		}

	}

	// ensurePVC - ensure the persistentvolumeclaim for web files folder is managed.
	result, err = r.ensurePVC(request, webproject, r.pvcForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureIngress - ensure the ingress object is managed.
	result, err = r.ensureIngress(request, webproject, r.ingressForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureDeployment - ensure the webproject deployment is managed.
	result, err = r.ensureDeployment(request, webproject, r.deploymentForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureSerivce - ensure the k8s service is managed.
	result, err = r.ensureService(request, webproject, r.serviceForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureSerivce - ensure the k8s service is managed.
	result, err = r.ensureService(request, webproject, r.backupServiceForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureEnvConfigMap - manage environmental variable config map for webproject.
	result, err = r.ensureEnvConfigMap(request, webproject, r.envConfigMapForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureCommonConfigMap - manage configmap used to set common environment variables for web and cli containers.
	result, err = r.ensureCommonConfigMap(request, webproject, r.commonConfigMapForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureInitContainerConfigMap - manage initcontainer.sh script configmap
	result, err = r.ensureInitContainerConfigMap(request, webproject, r.initContainerConfigMapForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureWebContainerCronJobConfigMap - manage script configmap
	result, err = r.ensureWebContainerCronJobConfigMap(request, webproject, r.webcontainerCronJobConfigMapForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureDatabaseSidecarCronJobConfigMap - manage script configmap
	result, err = r.ensureDatabaseSidecarCronJobConfigMap(request, webproject, r.databaseSidecarCronJobConfigMapForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureSecret - manage secrets for mysql database, etc
	result, err = r.ensureSecret(request, webproject, r.secretForWebproject(webproject))
	if result != nil {
		return *result, err
	}

	// ensureCronJob - manage backups for database.
	result, err = r.ensureCronJob(request, webproject, r.backupCronJob(webproject))
	if result != nil {
		return *result, err
	}

	// ensureCronJob - manage backups for database.
	result, err = r.ensureWebContainerCronJob(request, webproject, r.webContainerCronJob(webproject))
	if result != nil {
		return *result, err
	}

	// ensureCronJob - manage backups for database.
	result, err = r.ensureSearchContainerCronJob(request, webproject, r.searchContainerCronJob(webproject))
	if result != nil {
		return *result, err
	}

	// ensureCronJob - manage backups for database.
	result, err = r.ensureDatabaseContainerCronJob(request, webproject, r.databaseContainerCronJob(webproject))
	if result != nil {
		return *result, err
	}

	// updateWebProjectStatus
	// TODO: list each object name within webproject.
	WebProject := webproject

	lbs := map[string]string{
		"app.kubernetes.io/name": webproject.Name,
	}

	labelSelector := labels.SelectorFromSet(lbs)
	listOps := &client.ListOptions{Namespace: webproject.Namespace, LabelSelector: labelSelector}

	podList := &corev1.PodList{}
	if err = r.client.List(context.TODO(), podList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	// Count the pods that are pending or running as available
	var available []corev1.Pod
	for _, pod := range podList.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			available = append(available, pod)
		}
	}

	// List of pods
	availableNames := []string{}
	for _, pod := range available {
		availableNames = append(availableNames, pod.ObjectMeta.Name)
	}

	// get configMaps owned by this workload.
	configmapList := &corev1.ConfigMapList{}
	if err = r.client.List(context.TODO(), configmapList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	var configmapAvailable []corev1.ConfigMap
	for _, configmap := range configmapList.Items {
		configmapAvailable = append(configmapAvailable, configmap)
	}

	// List of configMaps
	configmapNames := []string{}
	for _, configmap := range configmapAvailable {
		configmapNames = append(configmapNames, configmap.ObjectMeta.Name)
	}

	// Get secrets owned by this webproject
	secretsList := &corev1.SecretList{}
	if err = r.client.List(context.TODO(), secretsList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	// get secrets owned by this workload.
	var secretsAvailable []corev1.Secret
	for _, secret := range secretsList.Items {
		secretsAvailable = append(secretsAvailable, secret)
	}

	// List of Secrets
	secretsNames := []string{}
	for _, secret := range secretsAvailable {
		secretsNames = append(secretsNames, secret.ObjectMeta.Name)
	}

	// Get deployments owned by this webproject
	deploymentsList := &appsv1.DeploymentList{}
	if err = r.client.List(context.TODO(), deploymentsList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	var deploymentsAvailable []appsv1.Deployment
	for _, deployment := range deploymentsList.Items {
		deploymentsAvailable = append(deploymentsAvailable, deployment)
	}

	// List of Deployments
	deploymentsNames := []string{}
	for _, deployment := range deploymentsAvailable {
		deploymentsNames = append(deploymentsNames, deployment.ObjectMeta.Name)
	}

	// Get deployments owned by this webproject
	cronjobList := &v1beta1.CronJobList{}
	if err = r.client.List(context.TODO(), cronjobList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	var cronjobAvailable []v1beta1.CronJob
	for _, cronjob := range cronjobList.Items {
		cronjobAvailable = append(cronjobAvailable, cronjob)
	}

	// List of Deployments
	cronjobNames := []string{}
	for _, cronjob := range cronjobAvailable {
		cronjobNames = append(cronjobNames, cronjob.ObjectMeta.Name)
	}
	// Update the status if necessary
	status := wp.WebProjectStatus{
		ConfigMapNames:  configmapNames,
		CronJobNames:    cronjobNames,
		DeploymentNames: deploymentsNames,
		PodNames:        availableNames,
		SecretNames:     secretsNames,
	}

	WebProject.Status = status
	err = r.client.Status().Update(context.TODO(), WebProject)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileWebproject) doesIngressExists(wp *wp.WebProject) bool {
	ingress := &networkingv1beta1.Ingress{}
	log.Info("Does Ingress exists")

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      workloadName(wp, "ingress"),
		Namespace: wp.Namespace,
	}, ingress)

	if err != nil {
		log.Error(err, "Ingress not found")
		return false
	}
	return true

}
