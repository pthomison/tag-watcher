package registryutils

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pthomison/errcheck"
)

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

func headImage(image string) *v1.Descriptor {
	ref, err := name.ParseReference(image)
	errcheck.Check(err)

	desc, err := remote.Head(ref)
	errcheck.Check(err)

	return desc
}

func ListRepository(repoStr string) []string {
	repo, err := name.NewRepository(repoStr)
	errcheck.Check(err)

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
