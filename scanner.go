package varchive

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileWithSize struct {
	path string
	size int64
}

type FilesWithSize []*FileWithSize

func ScanPaths() map[string]FilesWithSize {

	pathsAndFiles := make(map[string]FilesWithSize)

	for _, path := range settings.paths {

		if settings.verbose { 
			log.Printf("Scanning %s....\n", path)
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		filesForPath := FilesWithSize{}

		if fileInfo.IsDir() {
			if settings.verbose { 
				log.Printf("%s is a directory\n", path)
			}
			err := filepath.Walk(path, func(walkedPath string, fileInfo os.FileInfo, err error) error {
				if err == nil {
					if walkedPath != path { // the path itself is included in the results of Walk(path,...)
						if fileInfo.IsDir() {
							fatal(fmt.Sprintf("Recursive directories are not handled (%v)", walkedPath))
						} else {
							filesForPath = append(filesForPath, &FileWithSize{walkedPath, fileInfo.Size()})
						}
					} else {
						return err
					}
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		} else {
			fatal(fmt.Sprintf("Only paths can be specified, %v is a file", path))
		}

		pathsAndFiles[path] = filesForPath
	}

	return pathsAndFiles
}

