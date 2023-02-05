package controllers

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
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
					SourceRepository:    "docker.io/library/python",
					DestinationRegistry: fmt.Sprintf("%v/python", destRegistryUrl),
					Regex: v1alpha1.TagRegex{
						Match:  "3.9.16-alpine(.*)",
						Ignore: "windows",
					},
				},
			}
			Expect(k8sClient.Create(ctx, tagReflector)).Should(Succeed())

			name := types.NamespacedName{
				Name:      TagReflectorName,
				Namespace: TagReflectorNamespace,
			}

			waitFor := func(timeout time.Duration, test func(*v1alpha1.TagReflector) bool) {
				timeoutChan := make(chan bool, 1)
				go func() {
					time.Sleep(timeout)
					timeoutChan <- true
				}()
				for {
					select {
					case <-timeoutChan:
						return
					default:
						tagReflector = &v1alpha1.TagReflector{}
						Expect(k8sClient.Get(ctx, name, tagReflector)).Should(Succeed())

						if test(tagReflector) {
							return
						}
					}

				}
			}

			waitFor(10000*time.Millisecond, func(tr *v1alpha1.TagReflector) bool {

				if len(tr.Status.MatchedTags) != 3 {
					return false
				}

				for _, mt := range tr.Status.MatchedTags {
					if mt.DestinationDigest == "" || mt.SourceDigest == "" {
						return false
					}
				}

				return true
			})

			spew.Dump(tagReflector)

			Expect(len(tagReflector.Status.MatchedTags)).Should(Not(BeZero()))

		})
	})
})
