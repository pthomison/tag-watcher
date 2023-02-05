package controllers

import (
	"crypto/tls"
	"math/rand"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pthomison/errcheck"
)

func ListRepository(repoStr string) []string {
	repo, err := name.NewRepository(repoStr)
	errcheck.Check(err)

	tags, err := remote.List(repo)
	errcheck.Check(err)

	return tags
}

func GetImageDigest(imageReg string) (string, error) {
	head, err := headImage(imageReg)
	if err != nil {
		return "", err
	}
	digest := head.Digest
	return digest.String(), nil
}

func headImage(image string) (*v1.Descriptor, error) {
	ref, err := name.ParseReference(image)
	if err != nil {
		return nil, err
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	desc, err := remote.Head(ref)
	if err != nil {
		return nil, err
	}

	return desc, nil
}

func mustMoveImage(img v1.Image, dest name.Reference) {
	spew.Dump(dest)
	err := remote.Write(dest, img)
	errcheck.Check(err)
}

func mustParseReference(s string) name.Reference {
	ref, err := name.ParseReference(s)
	spew.Dump(err)
	errcheck.Check(err)
	return ref
}

// func CatalogRegistry(registryStr string) []string {
// 	registry, err := name.NewRegistry(registryStr)
// 	errcheck.Check(err)

// 	images, err := remote.Catalog(context.Background(), registry)
// 	errcheck.Check(err)

// 	return images
// }

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
