package controllers

import (
	"crypto/tls"
	"net/http"

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

// Currently Unused
// func CopyImage(src string, dest string) {
// 	ref, err := name.ParseReference(src)
// 	errcheck.Check(err)

// 	remoteRef, err := name.ParseReference(dest)
// 	errcheck.Check(err)

// 	image, err := remote.Image(ref)
// 	errcheck.Check(err)

// 	err = remote.Write(remoteRef, image)
// 	errcheck.Check(err)
// }

// func CatalogRegistry(registryStr string) []string {
// 	registry, err := name.NewRegistry(registryStr)
// 	errcheck.Check(err)

// 	images, err := remote.Catalog(context.Background(), registry)
// 	errcheck.Check(err)

// 	return images
// }
