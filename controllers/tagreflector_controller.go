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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/davecgh/go-spew/spew"
	"github.com/pthomison/errcheck"
	tagreflectorv1alpha1 "github.com/pthomison/tag-watcher/api/v1alpha1"
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

	// Reqest TagReflector Object
	tr, err := r.Get(ctx, req.NamespacedName)
	if client.IgnoreNotFound(err) != nil {
		// requeue in hopes that the error is transient
		return ctrl.Result{}, err
	} else if err != nil {
		// if object has been deleted, ignore
		return ctrl.Result{}, nil
	}

	// Setup Regexes To Test Tags Later On
	// TODO: Don't crash if regex doesn't compile; Report error and skip processing
	match := regexp.MustCompile(tr.Spec.Regex.Match)
	ignore := regexp.MustCompile(tr.Spec.Regex.Ignore)

	// Find All Tags Associated With The Spec Repository
	tags := ListRepository(tr.Spec.SourceRepository)

	// Create The Status Map If Needed
	if tr.Status.MatchedTags == nil {
		tr.Status.MatchedTags = make(map[string]*tagreflectorv1alpha1.MatchedTagStatus)
	}

	// Filter Tags Down Via Spec Regex
	regexFunc := func(tag string) bool {
		return !ignore.MatchString(tag) && match.FindString(tag) == tag
	}
	// Add a MatchedTagStatus object for desired tags
	for _, t := range tags {
		if regexFunc(t) && tr.Status.MatchedTags[t] == nil {
			tr.Status.MatchedTags[t] = &tagreflectorv1alpha1.MatchedTagStatus{
				Tag: t,
			}
		}
	}

	// Write Tags w/o Hash Data
	err = r.StatusUpdate(ctx, tr)
	if err != nil {
		return ctrl.Result{}, err
	}

	for i := range tr.Status.MatchedTags {
		sourceImage := fmt.Sprintf("%v:%v", tr.Spec.SourceRepository, tr.Status.MatchedTags[i].Tag)

		digest, err := GetImageDigest(sourceImage)
		errcheck.Check(err)

		sourceHash := digest

		// var destinationImage string
		// if tr.Spec.ReflectorSuffix == "" {
		// 	destinationImage = fmt.Sprintf("%v:%v", tr.Spec.DestinationRegistry, tr.Status.MatchedTags[i].Tag)
		// } else {
		// 	destinationImage = fmt.Sprintf("%v:%v-%v", tr.Spec.DestinationRegistry, tr.Status.MatchedTags[i].Tag, tr.Spec.ReflectorSuffix)
		// }
		destinationImage := fmt.Sprintf("%v:%v", tr.Spec.DestinationRegistry, tr.Status.MatchedTags[i].Tag)

		destinationHash, _ := GetImageDigest(destinationImage)

		imageDone := func() bool {
			return sourceHash != "" &&
				destinationHash != "" &&
				tr.Status.MatchedTags[i].SourceDigest == sourceHash &&
				tr.Status.MatchedTags[i].DestinationDigest == destinationHash
		}

		if imageDone() {
			fmt.Printf("Source && Destination Digests Matches; Skipping %v\n", sourceImage)
			continue
		}

		// Create a build request && execute it
		// br := containerutils.BuildReqest{
		// 	CTX:              ctx,
		// 	Obj:              tr,
		// 	Tag:              tr.Status.MatchedTags[i].Tag,
		// 	SourceImage:      sourceImage,
		// 	DestinationImage: destinationImage,
		// }
		// br.Build()

		// Exec Plugins
		a := tr.Spec.Action
		switch {
		case a.DockerBuild != nil:

		case a.Copy != nil:
			CopyImage(sourceImage, destinationImage)
		default:
			CopyImage(sourceImage, destinationImage)
		}

		// Don't update the SourceDigest until imageDone() invocation
		tr.Status.MatchedTags[i].SourceDigest = sourceHash

		digest, err = GetImageDigest(destinationImage)
		if err != nil {
			spew.Dump(ListRepository(tr.Spec.DestinationRegistry))
			errcheck.Check(err)
		}

		tr.Status.MatchedTags[i].DestinationDigest = digest

		// spew.Dump(tr.Status.MatchedTags[i])

		// Update the status with the image hash
		err = r.StatusUpdate(ctx, tr)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TagReflectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tagreflectorv1alpha1.TagReflector{}).
		Complete(r)
}

// Utility Wrapper Functions

func (r *TagReflectorReconciler) Get(ctx context.Context, name types.NamespacedName) (*tagreflectorv1alpha1.TagReflector, error) {
	tr := &tagreflectorv1alpha1.TagReflector{}
	err := r.Client.Get(ctx, name, tr)
	return tr, err
}

func (r *TagReflectorReconciler) StatusUpdate(ctx context.Context, tr *tagreflectorv1alpha1.TagReflector) error {
	err := r.Status().Update(ctx, tr)
	return err
}
