package containerutils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/pthomison/errcheck"
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

func (r *Request) StartContainer(c *container.Config) string {
	id, err := r.client.ContainerCreate(r.ctx, c, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, "testing")
	errcheck.Check(err)

	err = r.client.ContainerStart(r.ctx, id.ID, types.ContainerStartOptions{})
	errcheck.Check(err)

	return id.ID
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

func (r *Request) BuildImage() string {
	createResp, err := r.client.ContainerCreate(r.ctx, &container.Config{
		Image: "python:3",
		Cmd: []string{
			"python3",
			"-m",
			"http.server",
			"8080",
		},
	}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, "testing")
	errcheck.Check(err)

	id := createResp.ID

	err = r.client.ContainerStart(r.ctx, id, types.ContainerStartOptions{})
	errcheck.Check(err)

	defer r.DeleteContainer(id)

	exec, err := r.client.ContainerExecCreate(r.ctx, id, types.ExecConfig{
		Cmd: []string{
			"touch",
			"/iwashere",
		},
	})
	errcheck.Check(err)

	err = r.client.ContainerExecStart(r.ctx, exec.ID, types.ExecStartCheck{})
	errcheck.Check(err)

	resp, err := r.client.ContainerCommit(r.ctx, id, types.ContainerCommitOptions{
		Reference: "test-image:latest",
	})
	errcheck.Check(err)

	return resp.ID
}
