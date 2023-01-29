/*
Copyright 2023.

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
	"context"
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/pthomison/errcheck"
	"github.com/pthomison/tag-watcher/api/v1alpha1"
	tagreflectorv1alpha1 "github.com/pthomison/tag-watcher/api/v1alpha1"
	"github.com/pthomison/tag-watcher/pkg/registry"
)

// TagReflectorReconciler reconciles a TagReflector object
type TagReflectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tagreflector.operator.pthomison.com,resources=tagreflectors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tagreflector.operator.pthomison.com,resources=tagreflectors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tagreflector.operator.pthomison.com,resources=tagreflectors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TagReflector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TagReflectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	tr := &v1alpha1.TagReflector{}
	err := r.Client.Get(ctx, req.NamespacedName, tr)
	errcheck.Check(err)

	re := regexp.MustCompile(tr.Spec.Regex)
	tags := registry.ListRepository(tr.Spec.Registry)

	// s := strings.Join(tags, " ")

	for _, t := range tags {
		if re.FindString(t) == t {
			fmt.Println(t)
		}
	}

	// fmt.Println(re.String())
	// fmt.Println(s)

	// matchedtags := re.FindAll([]byte(s), -1)

	// // fmt.Println(matchedtags)

	// for _, v := range matchedtags {
	// 	fmt.Println(string(v))
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TagReflectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tagreflectorv1alpha1.TagReflector{}).
		Complete(r)
}
