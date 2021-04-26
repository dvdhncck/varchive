package varchive

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/** 

sample output from ffprobe

Stream #0:0[0x1011]: Video: h264 (High) (HDMV / 0x564D4448), yuv420p(top first), 1920x1080 [SAR 1:1 DAR 16:9], 25 fps, 50 tbr, 90k tbn, 50 tbc

*/

type VideoInfo struct {
	Width int64
	Height int64
	Fps float64
	Tbr float64
}

func GetVideoInfoUsingFfProbe(path string) (VideoInfo, error) {

	fullPath, err := filepath.Abs(path)

	if err != nil {
		return VideoInfo{}, err
	}

	args := []string{`-hide_banner`, fullPath}
	
	output := invoke("ffprobe", args)
	lines := strings.Split(output, "\n")

	return ParseVideoInfoFromFfProbe(lines)
}

// :TODO: move to static block?
func getFfProbeOutputParser() string {
	prefix := `Stream\ \#.+Video.+?`
	dimCapture := `(?P<width>\d{2,})x(?P<height>\d{2,})[,\ ]` // integer dimension, e.g. 100x200, must have at least 2 digits per dimension and be followed by whitespace
	
	fpsCapture := `(?P<fps>(?:[0-9]+\.)?[0-9]+)`  // floating point, e.g. xx.yyy or xx or .yyy
	tbrCapture := `(?P<tbr>(?:[0-9]+\.)?[0-9]+)`  // ditto

	return prefix + dimCapture + `.+\ ` + fpsCapture + `\ fps\,\ ` + tbrCapture + `\ tbr.+`
}

func ParseVideoInfoFromFfProbe(lines []string) (VideoInfo, error) {
	r := regexp.MustCompile(getFfProbeOutputParser())

	for _, line := range lines {
		matches := r.FindStringSubmatch(line)

		if matches != nil {
			width, _ := strconv.ParseInt(matches[1], 10, 64)
			height, _ := strconv.ParseInt(matches[2], 10, 64)
			fps, _ := strconv.ParseFloat(matches[3], 64)
			tbr, _ := strconv.ParseFloat(matches[4], 64)
			return VideoInfo{width, height, fps, tbr}, nil
		}
	}

	return VideoInfo{}, errors.New("parse failed")
}