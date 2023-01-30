package containerutils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/pthomison/tag-watcher/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type BuildReqest struct {
	CTX context.Context
	Obj *v1alpha1.TagReflector
	Tag string
}

func (b *BuildReqest) Build() string {
	_ = log.FromContext(b.CTX)

	baseImage := fmt.Sprintf("%v:%v", b.Obj.Spec.Repository, b.Tag)
	destImage := fmt.Sprintf("%v:%v-%v", b.Obj.Spec.DestinationRegistry, b.Tag, b.Obj.Spec.ReflectorSuffix)

	fmt.Printf("Building %v\n", destImage)

	cli := NewRequest()
	cli.PullImage(baseImage)

	buildContainer := cli.StartContainer(&container.Config{
		Image: baseImage,
		Cmd:   []string{"sleep", "9999"},
	})
	defer cli.DeleteContainer(buildContainer)

	for _, action := range b.Obj.Spec.Actions {
		cli.ExecContainer(buildContainer, action.Command.Args)
	}

	hash := cli.CommitContainer(buildContainer, destImage)
	cli.PushImage(destImage)

	return hash
}
