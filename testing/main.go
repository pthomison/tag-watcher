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

	// ScanRegistry(destRegistry)

	// platform, err := v1.ParsePlatform("linux/arm64")
	// errcheck.Check(err)

	// registryutils.CopyImage(srcImage, destImage)

	// // srcHead := registryutils.HeadImage(srcImage)
	// // spew.Dump(srcHead)

	// dstHead := registryutils.HeadImage(destImage)
	// spew.Dump(dstHead)
	// spew.Dump(dstHead)

	// fmt.Printf("%+v\n", registryutils.ListRepository(srcRepository))
	// fmt.Printf("%+v\n", registryutils.Get("sha256:7efc1ae7e6e9c5263d87845cb00f6ab7f6b27670cae29c9d93fa7910d6ab12c0"))
	registryutils.Hack()

}

func ScanRegistry(registry string) {
	imgs := registryutils.CatalogRegistry(registry)

	for _, img := range imgs {
		repo := fmt.Sprintf("%v/%v", registry, img)
		tags := registryutils.ListRepository(repo)
		for _, tag := range tags {
			fmt.Printf("%v:%v\n", repo, tag)
		}

	}
}
