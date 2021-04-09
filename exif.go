package varchive

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetVideoInfo(path string) (int64, int64, error) {

	fullPath, err := filepath.Abs(path)

	if err != nil {
		return 0, 0, err
	}

	args := []string{fullPath}
	output := invoke("exiftool", args)

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Image Size") {
			return parseDimensions(line)
		}
	}
	return 0, 0, errors.New("ImageSize tag not found")
}


func parseDimensions(line string) (int64, int64, error) {
	// line is of the form:
	//     Image Size      : 640x320

	r := regexp.MustCompile(`Image\ Size\s+\:\s+(?P<width>\d+)x(?P<height>\d+)`)
	matches := r.FindStringSubmatch(line)
	names := r.SubexpNames()

	width, height := int64(0), int64(0)

	for i, _ := range matches {
		switch names[i] {
		case ``:
			// ignore the whole group match
		case `width`:
			w, err := strconv.ParseInt(matches[i], 10, 64)
			if err == nil {
				width = w
			} else {
				return 0, 0, err
			}
		case `height`:
			h, err := strconv.ParseInt(matches[i], 10, 64)
			if err == nil {
				height = h
			} else {
				return 0, 0, err
			}
		default:
			return 0, 0, errors.New("unexpected matching group")
		}
	}
	return width, height, nil
}
