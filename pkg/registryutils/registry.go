package registryutils

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pthomison/errcheck"
)

// func Hack() {

// 	remoteImage := "index.docker.io/library/python:3"
// 	destImage := "127.0.0.1:15555/python:3"
// 	// remoteManifest := "index.docker.io/library/python@sha256:7efc1ae7e6e9c5263d87845cb00f6ab7f6b27670cae29c9d93fa7910d6ab12c0"

// 	CopyImage(remoteImage, destImage)

// 	spew.Dump(GetImageDigest(remoteImage))
// 	spew.Dump(GetImageDigest(destImage))

// }

// func GetManifest(imageRef string) *v1.Manifest {
// 	image := GetImage(imageRef)
// 	if image == nil {
// 		return nil
// 	}

// 	manifest, err := image.Manifest()
// 	errcheck.Check(err)

// 	return manifest
// }

func GetImageDigest(imageReg string) (string, error) {
	// image := GetImage(imageReg)
	// if image == nil {
	// 	return ""
	// }
	// digest, err := image.Digest()
	// errcheck.Check(err)
	// return digest.String()

	head, err := HeadImage(imageReg)
	if err != nil {
		return "", err
	}
	digest := head.Digest
	return digest.String(), nil

}

// func GetImage(imageRef string) v1.Image {
// 	descriptor := Get(imageRef)
// 	if descriptor == nil {
// 		return nil
// 	}

// 	image, err := descriptor.Image()
// 	errcheck.Check(err)

// 	return image
// }

func CopyImage(src string, dest string) {
	ref, err := name.ParseReference(src)
	errcheck.Check(err)

	remoteRef, err := name.ParseReference(dest)
	errcheck.Check(err)

	image, err := remote.Image(ref)
	errcheck.Check(err)

	err = remote.Write(remoteRef, image)
	errcheck.Check(err)
}

func HeadImage(image string) (*v1.Descriptor, error) {
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

// func Get(s string) *remote.Descriptor {
// 	ref, err := name.ParseReference(s)
// 	errcheck.Check(err)

// 	d, err := remote.Get(ref)
// 	if err != nil {
// 		return nil
// 	}

// 	return d
// }

func ListRepository(repoStr string) []string {
	repo, err := name.NewRepository(repoStr)
	errcheck.Check(err)

	// TODO: Figure out a clean platform parsing strategy/why the below doesn't work
	// platform, err := v1.ParsePlatform("linux/arm64")
	// errcheck.Check(err)
	// tags, err := remote.List(repo, remote.WithPlatform(*platform))

	tags, err := remote.List(repo)
	errcheck.Check(err)

	return tags
}

func CatalogRegistry(registryStr string) []string {
	registry, err := name.NewRegistry(registryStr)
	errcheck.Check(err)

	images, err := remote.Catalog(context.Background(), registry)
	errcheck.Check(err)

	return images
}
