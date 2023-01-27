package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/pthomison/errcheck"
)

const (
	tag = "3"

	srcRegistry  = "docker.io"
	destRegistry = "127.0.0.1:15555"

	srcRepository  = srcRegistry + "/library/python"
	destRepository = destRegistry + "/python"

	srcImage  = srcRepository + ":" + tag
	destImage = destRepository + ":" + tag
)

func main() {
	fmt.Println("Tag Watcher")
	listContainers()

	// spew.Dump(headImage(srcImage))
	// spew.Dump(listRepository(srcRepository))
	// spew.Dump(catalogRegistry(destRegistry))

}

func listContainers() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func copyImage(src string, dest string) {
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

func listRepository(repoStr string) []string {
	repo, err := name.NewRepository(repoStr)
	errcheck.Check(err)

	tags, err := remote.List(repo)
	errcheck.Check(err)

	return tags
}

func catalogRegistry(registryStr string) []string {
	registry, err := name.NewRegistry(registryStr)
	errcheck.Check(err)

	images, err := remote.Catalog(context.Background(), registry)
	errcheck.Check(err)

	return images
}
