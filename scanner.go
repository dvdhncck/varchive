package varchive

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileWithSize struct {
	path string
	size int64
}

type FilesWithSize []*FileWithSize

func ScanPaths() map[string]FilesWithSize {

	if settings.verbose {
		Log("Scanning paths")
	}

	pathsAndFiles := make(map[string]FilesWithSize)

	for _, path := range settings.paths {

		if settings.verbose { 
			Log("Scanning %s....", path)
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			fatal(err.Error())
		}

		filesForPath := FilesWithSize{}

		if fileInfo.IsDir() {
			if settings.verbose { 
				Log("%s is a directory", path)
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
			Log("%v is a file and will be ignored", path)
		}

		pathsAndFiles[path] = filesForPath
	}

	return pathsAndFiles
}

