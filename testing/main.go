package main

import (
	"fmt"
)

func main() {
	fmt.Println("Tag Watcher - Testing")
	// registryutils.Hack()

}

// func ScanRegistry(registry string) {
// 	imgs := registryutils.CatalogRegistry(registry)

// 	for _, img := range imgs {
// 		repo := fmt.Sprintf("%v/%v", registry, img)
// 		tags := registryutils.ListRepository(repo)
// 		for _, tag := range tags {
// 			fmt.Printf("%v:%v\n", repo, tag)
// 		}

// 	}
// }
