/*
Copyright 2021.

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

package controllers

import (
	"bitbucket.org/klovercloud/klovercloud-operator/constant"
	"bitbucket.org/klovercloud/klovercloud-operator/helper"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"

	servicev1alpha1 "bitbucket.org/klovercloud/klovercloud-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KlovercloudFacadeReconciler reconciles a KlovercloudFacade object
type KlovercloudFacadeReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Context context.Context
	Err     error
}

//+kubebuilder:rbac:groups=service.klovercloud.com,resources=klovercloudfacades,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=service.klovercloud.com,resources=klovercloudfacades/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=service.klovercloud.com,resources=klovercloudfacades/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *KlovercloudFacadeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Err = nil
	r.Context = ctx
	r.Log = ctrl.Log.WithName("reconcile")
	r.Log.Info(fmt.Sprintf("'%s' reconcile", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)

	r.checkAndApplyDeployment()
	r.Log.Info(fmt.Sprintf("'%s' deployment applied", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KlovercloudFacadeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicev1alpha1.KlovercloudFacade{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Pod{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *KlovercloudFacadeReconciler) Init() {
	r.Err = nil
	r.Log = ctrl.Log.WithName("init")
	r.Context = context.Background()

	r.checkForControllerCacheStart()
	r.Log.Info("controller cache has been started", "controller", constant.KlovercloudFacadeController)

	r.checkAndCreateKlovercloudNamespace()
	r.Log.Info("'klovercloud' namespace exists", "controller", constant.KlovercloudFacadeController)

	r.checkAndCreateCustomResource()
	r.Log.Info("'klovercloud' custom resource deployed", "controller", constant.KlovercloudFacadeController)

	r.Err = nil
}

func (r *KlovercloudFacadeReconciler) checkForControllerCacheStart() {
	/*
		Waiting for the cache to be started
	*/
	time.Sleep(2 * time.Second)
	namespace := &corev1.Namespace{}
	for {
		r.Log.Info("checking for controller cache", "controller", constant.KlovercloudFacadeController)
		r.Err = r.Get(r.Context, types.NamespacedName{Name: constant.KlovercloudNamespace}, namespace)
		if r.Err != nil && helper.Common().CheckSubstrings(r.Err.Error(), "cache", "not", "started") == true {
			r.Log.Error(r.Err, "controller cache has not started yet. rechecking after 2 seconds", "controller", constant.KlovercloudFacadeController)
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}
	r.Err = nil
}

func (r *KlovercloudFacadeReconciler) checkAndCreateKlovercloudNamespace() {
	/*
		If klovercloud namespace does not exists then creating new one
	*/
	namespace := &corev1.Namespace{}
	for {
		r.Log.Info(fmt.Sprintf("checking for '%s' namespace", constant.KlovercloudNamespace), "controller", constant.KlovercloudFacadeController)
		r.Err = r.Get(r.Context, types.NamespacedName{Name: constant.KlovercloudNamespace}, namespace)
		if r.Err != nil {
			if helper.Common().CheckSubstrings(r.Err.Error(), "cache", "not", "started") == true {
				r.checkForControllerCacheStart()

			} else if errors.IsNotFound(r.Err) == true {
				r.Log.Info(fmt.Sprintf("creating new '%s' namespace", constant.KlovercloudNamespace), "controller", constant.KlovercloudFacadeController)
				_ = r.Create(r.Context, helper.Common().KlovercloudNamespace())
			} else {
				break
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	r.Err = nil
}

func (r *KlovercloudFacadeReconciler) checkAndCreateCustomResource() {
	cr := &servicev1alpha1.KlovercloudFacade{}
	for {
		r.Log.Info(fmt.Sprintf("checking for '%s' cr", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
		r.Err = r.Get(r.Context, types.NamespacedName{Name: constant.KlovercloudFacadeDeployment, Namespace: constant.KlovercloudNamespace}, cr)
		if r.Err != nil {
			if helper.Common().CheckSubstrings(r.Err.Error(), "cache", "not", "started") == true {
				r.checkForControllerCacheStart()

			} else if errors.IsNotFound(r.Err) == true {
				r.Log.Info(fmt.Sprintf("creating new '%s' cr", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
				_ = r.Create(r.Context, helper.KlovercloudFacade().CustomResourceV1Alpha1())
			} else {
				break
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	r.Err = nil
}

func (r *KlovercloudFacadeReconciler) checkAndApplyDeployment() {
	deployment := &appsv1.Deployment{}
	for {
		r.Log.Info(fmt.Sprintf("checking for '%s' deployment", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
		r.Err = r.Get(r.Context, types.NamespacedName{Name: constant.KlovercloudFacadeDeployment, Namespace: constant.KlovercloudNamespace}, deployment)
		if r.Err != nil {
			if helper.Common().CheckSubstrings(r.Err.Error(), "cache", "not", "started") == true {
				r.checkForControllerCacheStart()

			} else if errors.IsNotFound(r.Err) == true {
				r.Log.Info(fmt.Sprintf("creating new '%s' deployment", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
				err := r.Create(r.Context, helper.KlovercloudFacade().Deployment())
				if err != nil {
					r.Log.Error(err, fmt.Sprintf("unable to create new '%s' deployment. retrying..", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
				}

				time.Sleep(1 * time.Second)
				r.Err = r.Get(r.Context, types.NamespacedName{Name: constant.KlovercloudFacadeDeployment, Namespace: constant.KlovercloudNamespace}, deployment)
				if r.Err == nil {
					break
				}
				time.Sleep(1 * time.Second)

			} else {
				r.Log.Error(r.Err, fmt.Sprintf("something went wrong while fetching '%s' deployment. retrying..", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
				time.Sleep(3 * time.Second)
			}

		} else {
			r.Log.Info(fmt.Sprintf("updating '%s' deployment", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
			r.Err = r.Update(r.Context, helper.KlovercloudFacade().Deployment())
			if r.Err != nil {
				r.Log.Error(r.Err, fmt.Sprintf("unable to update '%s' deployment. retrying..", constant.KlovercloudFacadeDeployment), "controller", constant.KlovercloudFacadeController)
			} else {
				break
			}
		}
	}
	r.Err = nil
}
