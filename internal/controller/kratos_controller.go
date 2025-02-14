/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"google.golang.org/api/container/v1"
	operatorv1 "operator.kratos.io/kratos/api/v1"
	"operator.kratos.io/kratos/internal/cloud"
)

// KratosReconciler reconciles a Kratos object
type KratosReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	GkeService *container.Service
}

// +kubebuilder:rbac:groups=operator.opertor.kratos.io,resources=kratos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.opertor.kratos.io,resources=kratos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.opertor.kratos.io,resources=kratos/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kratos object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *KratosReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var kratos operatorv1.Kratos
	if err := r.Get(ctx, req.NamespacedName, &kratos); err != nil {
		log.Error(err, "unable to fetch Kratos")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !kratos.DeletionTimestamp.IsZero() {
		if containsString(kratos.Finalizers, "kratos.finalizer") {
			if err := cloud.DeleteCluster(ctx, r.GkeService, &kratos); err != nil {
				log.Error(err, "failed to delete GKE cluster")
				return ctrl.Result{}, err
			}
			kratos.Finalizers = removeString(kratos.Finalizers, "kratos.finalizer")
			if err := r.Update(ctx, &kratos); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !containsString(kratos.Finalizers, "kratos.finalizer") {
		kratos.Finalizers = append(kratos.Finalizers, "kratos.finalizer")
		if err := r.Update(ctx, &kratos); err != nil {
			return ctrl.Result{}, err
		}
	}

	exists, err := cloud.ClusterExists(ctx, r.GkeService, &kratos)
	if err != nil {
		log.Error(err, "failed to check cluster existence")
		return ctrl.Result{}, err
	}

	if !exists {
		if err := cloud.CreateCluster(ctx, r.GkeService, &kratos); err != nil {
			log.Error(err, "failed to create GKE cluster")
			return ctrl.Result{}, err
		}
	} else {
		if err := cloud.UpdateCluster(ctx, r.GkeService, &kratos); err != nil {
			log.Error(err, "failed to update GKE cluster")
			return ctrl.Result{}, err
		}
	}

	kratos.Status.Phase = operatorv1.PhaseRunning
	if err := r.Status().Update(ctx, &kratos); err != nil {
		log.Error(err, "failed to update Kratos status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KratosReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.Kratos{}).
		Named("kratos").
		Complete(r)
}

// Helper functions to handle finalizers
func containsString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func removeString(slice []string, str string) []string {
	result := []string{}
	for _, v := range slice {
		if v != str {
			result = append(result, v)
		}
	}
	return result
}
