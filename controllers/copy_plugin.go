package controllers

import (
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pthomison/errcheck"
)

func CopyImage(src string, dest string) {
	srcRef := mustParseReference(src)
	destRef := mustParseReference(dest)

	image, err := remote.Image(srcRef)
	errcheck.Check(err)

	mustMoveImage(image, destRef)
}
