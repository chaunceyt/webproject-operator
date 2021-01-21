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
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Webproject

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &networkingv1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wp.WebProject{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
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

	// ensureSecret - manage secrets for mysql database, etc
	result, err = r.ensureSecret(request, webproject, r.secretForWebproject(webproject))
	if result != nil {
		return *result, err
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
