package registryutils

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pthomison/errcheck"
)

func Hack() {

	remoteImage := "index.docker.io/library/python:3"
	destImage := "127.0.0.1:15555/python:3"
	// remoteManifest := "index.docker.io/library/python@sha256:7efc1ae7e6e9c5263d87845cb00f6ab7f6b27670cae29c9d93fa7910d6ab12c0"

	CopyImage(remoteImage, destImage)

	spew.Dump(GetImageDigest(remoteImage))
	spew.Dump(GetImageDigest(destImage))

}

func GetManifest(imageRef string) *v1.Manifest {
	image := GetImage(imageRef)

	manifest, err := image.Manifest()
	errcheck.Check(err)

	return manifest
}

func GetImageDigest(imageReg string) string {
	digest, err := GetImage(imageReg).Digest()
	errcheck.Check(err)
	return digest.String()
}

func GetImage(imageRef string) v1.Image {
	descriptor := Get(imageRef)

	image, err := descriptor.Image()
	errcheck.Check(err)

	return image
}

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

func HeadImage(image string) *v1.Descriptor {
	ref, err := name.ParseReference(image)
	errcheck.Check(err)

	desc, err := remote.Head(ref)
	errcheck.Check(err)

	return desc
}

func Get(s string) *remote.Descriptor {
	ref, err := name.ParseReference(s)
	errcheck.Check(err)

	d, err := remote.Get(ref)
	errcheck.Check(err)

	return d
}

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
