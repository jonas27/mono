package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func visitfunc(path string, d fs.DirEntry, err error) error {
	if strings.HasSuffix(path, ".jpg") {
		log.Println(d)
		jpeg.
	}
	return nil
}

func compress() {
	// items, _ := os.ReadDir(".")
	// for _, item := range items {
	// }
	filepath.WalkDir("static", visitfunc)
	// log.Println(item)
}
