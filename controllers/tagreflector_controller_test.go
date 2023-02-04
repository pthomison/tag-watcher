package controllers

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pthomison/tag-watcher/api/v1alpha1"
)

var _ = Describe("TagReflector controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		TagReflectorName      = "test-tag-reflector"
		TagReflectorNamespace = "default"

		// timeout  = time.Second * 10
		// duration = time.Second * 10
		// interval = time.Millisecond * 250
	)

	Context("Basic Tag Reflector - Copy Definition", func() {
		It("Should do some stuff", func() {
			By("By creating a new TagReflector")
			tagReflector := &v1alpha1.TagReflector{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "tagreflector.operator.pthomison.com/v1alpha1",
					Kind:       "TagReflector",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      TagReflectorName,
					Namespace: TagReflectorNamespace,
				},
				Spec: v1alpha1.TagReflectorSpec{
					Repository:          "docker.io/library/python",
					DestinationRegistry: fmt.Sprintf("%v/python", registryUrl),
					Regex: v1alpha1.TagRegex{
						Match:  "3.9.16-alpine(.*)",
						Ignore: "windows",
					},
				},
			}
			Expect(k8sClient.Create(ctx, tagReflector)).Should(Succeed())

			time.Sleep(5000 * time.Millisecond)

			name := types.NamespacedName{
				Name:      TagReflectorName,
				Namespace: TagReflectorNamespace,
			}
			tagReflector = &v1alpha1.TagReflector{}
			Expect(k8sClient.Get(ctx, name, tagReflector)).Should(Succeed())

			// spew.Dump(tagReflector)

			// Expect(len(tagReflector.Status.MatchedTags)).Should(Not(BeZero()))

		})
	})
})
