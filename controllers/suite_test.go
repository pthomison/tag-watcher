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
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/pthomison/errcheck"
	tagreflectorv1alpha1 "github.com/pthomison/tag-watcher/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg             *rest.Config
	k8sClient       client.Client // You'll be using this client in your tests.
	testEnv         *envtest.Environment
	ctx             context.Context
	cancel          context.CancelFunc
	srcRegistry     *httptest.Server
	srcRegistryUrl  string
	destRegistry    *httptest.Server
	destRegistryUrl string
)

// Bootstrap Ginko
// https://onsi.github.io/ginkgo/
func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = tagreflectorv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&TagReflectorReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()

		extract := func(s *httptest.Server) string {
			u, err := url.Parse(s.URL)
			Expect(err).ToNot(HaveOccurred())

			r := fmt.Sprintf("%v%v", u.Host, u.Path)

			Expect(r).NotTo(BeNil())

			return r
		}

		srcRegistry = httptest.NewServer(registry.New())
		destRegistry = httptest.NewServer(registry.New())
		srcRegistryUrl = extract(srcRegistry)
		destRegistryUrl = extract(destRegistry)

		// var err error
		// registryServer, err = registry.TLS("registry.pthomison.com")
		// Expect(err).ToNot(HaveOccurred())

		//
		// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	}()

	// wait for registies to be available
	for {
		if srcRegistryUrl != "" && destRegistryUrl != "" {
			break
		}
	}

	// generate image
	randomImage := func() v1.Image {
		rnd, err := random.Image(1024, 1)
		errcheck.Check(err)
		return rnd
	}

	i := randomImage()
	imageName := fmt.Sprintf("%v/random:testing", srcRegistryUrl)
	fmt.Println(imageName)

	mustMoveImage(i, mustParseReference(imageName))
	fmt.Println(imageName)

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	cancel()
	srcRegistry.Close()
	destRegistry.Close()
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
