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

	fmt.Println("--- starting ---")

	tr := &v1alpha1.TagReflector{}
	err := r.Client.Get(ctx, req.NamespacedName, tr)
	errcheck.Check(err)

	match := regexp.MustCompile(tr.Spec.Regex.Match)
	ignore := regexp.MustCompile(tr.Spec.Regex.Ignore)
	tags := registryutils.ListRepository(tr.Spec.Repository)

	if tr.Status.MatchedTags == nil {
		tr.Status.MatchedTags = make(map[string]*tagreflectorv1alpha1.MatchedTagStatus)
	}

	for _, t := range tags {
		if !ignore.MatchString(t) && match.FindString(t) == t {
			if tr.Status.MatchedTags[t] == nil {
				tr.Status.MatchedTags[t] = &v1alpha1.MatchedTagStatus{
					Tag: t,
				}
			}
		}
	}

	err = r.Status().Update(ctx, tr)
	if err != nil {
		return ctrl.Result{}, err
	}

	for i := range tr.Status.MatchedTags {
		if tr.Status.MatchedTags[i].Hash != "" {
			continue
		}

		br := BuildReqest{
			ctx: ctx,
			obj: tr,
			tag: tr.Status.MatchedTags[i].Tag,
		}
		tr.Status.MatchedTags[i].Hash = br.Build()

		err = r.Status().Update(ctx, tr)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	fmt.Println("--- ending ---")

	return ctrl.Result{}, nil
}

type BuildReqest struct {
	ctx context.Context
	obj *v1alpha1.TagReflector
	tag string
}

func (b *BuildReqest) Build() string {
	_ = log.FromContext(b.ctx)

	baseImage := fmt.Sprintf("%v:%v", b.obj.Spec.Repository, b.tag)
	destImage := fmt.Sprintf("%v:%v-%v", b.obj.Spec.DestinationRegistry, b.tag, b.obj.Spec.ReflectorSuffix)

	fmt.Printf("Building %v\n", destImage)

	cli := containerutils.NewRequest()
	cli.PullImage(baseImage)

	buildContainer := cli.StartContainer(&container.Config{
		Image: baseImage,
		Cmd:   []string{"sleep", "9999"},
	})
	defer cli.DeleteContainer(buildContainer)

	for _, action := range b.obj.Spec.Actions {
		cli.ExecContainer(buildContainer, action.Command.Args)
	}

	hash := cli.CommitContainer(buildContainer, destImage)
	cli.PushImage(destImage)

	return hash
}

// SetupWithManager sets up the controller with the Manager.
func (r *TagReflectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tagreflectorv1alpha1.TagReflector{}).
		Complete(r)
}
