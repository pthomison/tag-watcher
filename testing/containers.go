package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pthomison/errcheck"
)

func listContainers() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	errcheck.Check(err)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	errcheck.Check(err)

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

// f, err := os.CreateTemp("", "test-build-archive-*.tar")
// errcheck.Check(err)
// // defer f.Close()
// fname := f.Name()
// // defer os.Remove(f.Name())
// fmt.Println(fname)

// gw := gzip.NewWriter(f)
// tarball := tar.NewWriter(f)
// defer tarball.Close()
// defer gw.Close()

// hdr := &tar.Header{
// 	Name: "Dockerfile",
// 	Mode: 0777,
// 	Size: int64(len(dockerfile)),
// }

// err = tarball.WriteHeader(hdr)
// errcheck.Check(err)

// _, err = tarball.Write([]byte(dockerfile))
// errcheck.Check(err)

// tarball.Close()
// f.Close()

// f, err = os.Open(fname)
// errcheck.Check(err)

// tb := tar.NewReader(f)

// spew.Dump(tb.Next())

// f.Write([]byte(dockerfile))

// response, err := cli.ImageBuild(context.Background(), tb, types.ImageBuildOptions{})
// errcheck.Check(err)
// defer response.Body.Close()

// b, err := ioutil.ReadAll(response.Body)
// errcheck.Check(err)

// fmt.Println(string(b))
