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

	"github.com/docker/docker/api/types/container"
	"github.com/pthomison/errcheck"
	"github.com/pthomison/tag-watcher/api/v1alpha1"
	tagreflectorv1alpha1 "github.com/pthomison/tag-watcher/api/v1alpha1"
	"github.com/pthomison/tag-watcher/pkg/containerutils"
	"github.com/pthomison/tag-watcher/pkg/registryutils"
)

// TagReflectorReconciler reconciles a TagReflector object
type TagReflectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tagreflector.operator.pthomison.com,resources=tagreflectors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tagreflector.operator.pthomison.com,resources=tagreflectors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tagreflector.operator.pthomison.com,resources=tagreflectors/finalizers,verbs=update

func (r *TagReflectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	tr := &v1alpha1.TagReflector{}
	err := r.Client.Get(ctx, req.NamespacedName, tr)
	errcheck.Check(err)

	re := regexp.MustCompile(tr.Spec.Regex)
	tags := registryutils.ListRepository(tr.Spec.Repository)

	for _, t := range tags {
		if re.FindString(t) == t {
			fmt.Println(t)
			br := BuildReqest{
				ctx: ctx,
				obj: tr,
				// repo:    tr.Spec.Registry,
				tag: t,
				// command: tr.Spec.Commands,
			}
			br.Build()
		}
	}

	return ctrl.Result{}, nil
}

type BuildReqest struct {
	ctx context.Context
	obj *v1alpha1.TagReflector

	// repo    string
	tag string
	// command []string
}

func (b *BuildReqest) Build() error {
	_ = log.FromContext(b.ctx)

	dockerReq := containerutils.NewRequest()

	buildContainer := dockerReq.StartContainer(&container.Config{
		Image: fmt.Sprintf("%v:%v", b.obj.Spec.Repository, b.tag),
		Cmd: []string{
			"python3",
			"-m",
			"http.server",
			"8080",
		},
	})

	defer dockerReq.DeleteContainer(buildContainer)

	dockerReq.ExecContainer(buildContainer, b.obj.Spec.Commands)

	dest := fmt.Sprintf("%v:%v-%v", b.obj.Spec.DestinationRegistry, b.tag, b.obj.Spec.ReflectorSuffix)

	dockerReq.CommitContainer(buildContainer, dest)

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TagReflectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tagreflectorv1alpha1.TagReflector{}).
		Complete(r)
}
