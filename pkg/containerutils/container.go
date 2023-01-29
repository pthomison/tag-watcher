package containerutils

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/pthomison/errcheck"
	"github.com/pthomison/tag-watcher/pkg/utils"
)

type Request struct {
	client *client.Client
	ctx    context.Context
}

func NewRequest() *Request {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	errcheck.Check(err)

	return &Request{
		client: cli,
		ctx:    context.Background(),
	}
}

func (r *Request) CreateContainer(c *container.Config) string {
	containerName := fmt.Sprintf("%v-%v", "testing", utils.RandomString(10))

	id, err := r.client.ContainerCreate(r.ctx, c, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, containerName)
	errcheck.Check(err)

	return id.ID
}

func (r *Request) StartContainer(c *container.Config) string {
	containerID := r.CreateContainer(c)

	err := r.client.ContainerStart(r.ctx, containerID, types.ContainerStartOptions{})
	errcheck.Check(err)

	return containerID
}

func (r *Request) DeleteContainer(id string) {
	err := r.client.ContainerRemove(r.ctx, id, types.ContainerRemoveOptions{
		Force: true,
	})
	errcheck.Check(err)
}

func (r *Request) ExecContainer(id string, command []string) {
	exec, err := r.client.ContainerExecCreate(r.ctx, id, types.ExecConfig{
		Cmd: command,
	})
	errcheck.Check(err)

	err = r.client.ContainerExecStart(r.ctx, exec.ID, types.ExecStartCheck{})
	errcheck.Check(err)
}

func (r *Request) CommitContainer(id string, ref string) string {
	resp, err := r.client.ContainerCommit(r.ctx, id, types.ContainerCommitOptions{
		Reference: ref,
	})
	errcheck.Check(err)

	return resp.ID
}

func (r *Request) PullImage(image string) {
	body, err := r.client.ImagePull(r.ctx, image, types.ImagePullOptions{})
	errcheck.Check(err)
	ioutil.ReadAll(body)
	body.Close()
}

func (r *Request) PushImage(image string) {
	body, err := r.client.ImagePush(r.ctx, image, types.ImagePushOptions{
		RegistryAuth: "123",
	})
	errcheck.Check(err)
	ioutil.ReadAll(body)
	body.Close()
}
