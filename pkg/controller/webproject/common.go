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
	"reflect"

	wpv1 "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// create docker config for workload.

// DockerConfigJson represents ~/.docker/config.json file info.
type DockerConfigJson struct {
	Auths DockerConfig `json:"auths"`
}

// DockerConfig represents the config file used by the docker CLI.
type DockerConfig map[string]DockerConfigEntry

// DockerConfigEntry represents an Auth entry in the dockerconfigjson.
type DockerConfigEntry struct {
	Username string
	Password string
	Email    string
}

func (r *ReconcileWebproject) ensureDeployment(request reconcile.Request, instance *wpv1.WebProject, dep *appsv1.Deployment) (*reconcile.Result, error) {
	found := &appsv1.Deployment{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}

		// Deployment was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	deploy := r.deploymentForWebproject(instance)
	if !reflect.DeepEqual(found.Spec.Template.Spec, deploy.Spec.Template.Spec) {
		err := r.client.Update(context.Background(), deploy)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil

	}
	return nil, nil

}

func (r *ReconcileWebproject) ensureService(request reconcile.Request, instance *wpv1.WebProject, s *corev1.Service) (*reconcile.Result, error) {
	found := &corev1.Service{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileWebproject) ensurePVC(request reconcile.Request, instance *wpv1.WebProject, s *corev1.PersistentVolumeClaim) (*reconcile.Result, error) {
	found := &corev1.PersistentVolumeClaim{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the PVC
		log.Info("Creating a new PVC", "PVC.Namespace", s.Namespace, "PVC.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new PVC", "PVC.Namespace", s.Namespace, "PVC.Name", s.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the pvc not existing
		log.Error(err, "Failed to get PVC")
		return &reconcile.Result{}, err
	}

	return nil, nil

}

func (r *ReconcileWebproject) ensureIngress(request reconcile.Request, instance *wpv1.WebProject, ing *networkingv1beta1.Ingress) (*reconcile.Result, error) {
	found := &networkingv1beta1.Ingress{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      ing.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the ingress
		log.Info("Creating a new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
		err = r.client.Create(context.TODO(), ing)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the ingress not existing
		log.Error(err, "Failed to get Ingress")
		return &reconcile.Result{}, err
	}

	ingress := r.ingressForWebproject(instance)

	if !reflect.DeepEqual(found.Spec.Rules, ingress.Spec.Rules) {
		err := r.client.Update(context.Background(), ingress)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Ingress.Namespace", found.Namespace, "Ingress.Name", found.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil

	}
	return nil, nil
}

func (r *ReconcileWebproject) ensureEnvConfigMap(request reconcile.Request, instance *wpv1.WebProject, cm *corev1.ConfigMap) (*reconcile.Result, error) {
	found := &corev1.ConfigMap{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      cm.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the configmap
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the configmap not existing
		log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileWebproject) ensureInitContainerConfigMap(request reconcile.Request, instance *wpv1.WebProject, cm *corev1.ConfigMap) (*reconcile.Result, error) {
	found := &corev1.ConfigMap{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      cm.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the configmap
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the configmap not existing
		log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}

	configMap := r.initContainerConfigMapForWebproject(instance)
	if !reflect.DeepEqual(found.Data, configMap.Data) {
		log.Info("Updating ConfigMap", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
		err := r.client.Update(context.Background(), configMap)
		if err != nil {
			log.Error(err, "Failed to update ConfigMap", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil
	}

	return nil, nil
}

func (r *ReconcileWebproject) ensureCommonConfigMap(request reconcile.Request, instance *wpv1.WebProject, cm *corev1.ConfigMap) (*reconcile.Result, error) {
	found := &corev1.ConfigMap{}
	ctx := context.Background()

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      cm.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the configmap
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the configmap not existing
		log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}

	configMap := r.commonConfigMapForWebproject(instance)
	if !equality.Semantic.DeepDerivative(found.Data, configMap.Data) {
		log.Info("Updating ConfigMap", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
		err := r.client.Update(ctx, configMap)
		if err != nil {
			log.Error(err, "Failed to update ConfigMap", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
			return &reconcile.Result{}, err
		}
		return nil, nil
	}

	return nil, nil
}

func (r *ReconcileWebproject) ensureSecret(request reconcile.Request, instance *wpv1.WebProject, secret *corev1.Secret) (*reconcile.Result, error) {
	found := &corev1.Secret{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the secret
		log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful - return and requeue
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileWebproject) ensureDockerConfigSecret(request reconcile.Request, instance *wpv1.WebProject, secret *corev1.Secret) (*reconcile.Result, error) {
	found := &corev1.Secret{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the secret
		log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}

	return nil, nil
}
