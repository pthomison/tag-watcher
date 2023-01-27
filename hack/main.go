package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(headImage(srcImage))
	spew.Dump(listRepository(srcRepository))
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
