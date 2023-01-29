package main

import (
	"fmt"

	"github.com/pthomison/tag-watcher/pkg/registryutils"
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
	fmt.Println("Tag Watcher - Testing")
	// r := containerutils.NewRequest()

	registry := destRegistry

	imgs := registryutils.CatalogRegistry(registry)

	for _, img := range imgs {
		repo := fmt.Sprintf("%v/%v", destRegistry, img)
		tags := registryutils.ListRepository(repo)
		for _, tag := range tags {
			fmt.Printf("%v:%v\n", repo, tag)
		}

	}

}

// type Request struct {
// 	client *client.Client
// 	ctx    context.Context
// }

// func NewRequest() *Request {
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	errcheck.Check(err)

// 	return &Request{
// 		client: cli,
// 		ctx:    context.Background(),
// 	}
// }

// func (r *Request) StartContainer() string {
// 	id, err := r.client.ContainerCreate(r.ctx, &container.Config{
// 		Image: "python:3",
// 		Cmd: []string{
// 			"python3",
// 			"-m",
// 			"http.server",
// 			"8080",
// 		},
// 	}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, "testing")
// 	errcheck.Check(err)

// 	err = r.client.ContainerStart(r.ctx, id.ID, types.ContainerStartOptions{})
// 	errcheck.Check(err)

// 	return id.ID
// }

// func (r *Request) DeleteContainer(id string) {
// 	err := r.client.ContainerRemove(r.ctx, id, types.ContainerRemoveOptions{
// 		Force: true,
// 	})
// 	errcheck.Check(err)
// }

// func (r *Request) BuildImage() string {
// 	createResp, err := r.client.ContainerCreate(r.ctx, &container.Config{
// 		Image: "python:3",
// 		Cmd: []string{
// 			"python3",
// 			"-m",
// 			"http.server",
// 			"8080",
// 		},
// 	}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, "testing")
// 	errcheck.Check(err)

// 	id := createResp.ID

// 	err = r.client.ContainerStart(r.ctx, id, types.ContainerStartOptions{})
// 	errcheck.Check(err)

// 	defer r.DeleteContainer(id)

// 	exec, err := r.client.ContainerExecCreate(r.ctx, id, types.ExecConfig{
// 		Cmd: []string{
// 			"touch",
// 			"/iwashere",
// 		},
// 	})
// 	errcheck.Check(err)

// 	err = r.client.ContainerExecStart(r.ctx, exec.ID, types.ExecStartCheck{})
// 	errcheck.Check(err)

// 	resp, err := r.client.ContainerCommit(r.ctx, id, types.ContainerCommitOptions{
// 		Reference: "test-image:latest",
// 	})
// 	errcheck.Check(err)

// 	return resp.ID
// }
