package main

import (
	"davidhancock.com/varchive"
	"fmt"
	"regexp"
	"testing"
)

func Test_parseDimensionsUsingffProbe(t *testing.T) {

	info, err := varchive.GetVideoInfoUsingFfProbe(`test-data/one/sample 2.mpg`) // path relative to this .go file

	if err == nil {
		fmt.Printf("%d x %d @ %f (%f)", info.Width, info.Height, info.Fps, info.Tbr)
	} else {
		t.Fatal(err)
	}

}

func Test_regexCapturing(t *testing.T) {

	r := regexp.MustCompile(varchive.GetFfProbeOutputParser())
	
	test := `Stream #0:0[0x1e0]: Video: mpeg1video, yuv420p(tv), 640x320 [SAR 1:1 DAR 2:1], 104857 kb/s, 25.52 fps, 24.99 tbr,...`
	matches := r.FindStringSubmatch(test)
	if matches[1] != `640` {
		t.Fatal(`width wrong`)
	}
	if matches[2] != `320` {
		t.Fatal(`width wrong`)
	}
	if matches[3] != `25.52` {
		t.Fatal(`width wrong`)
	}
	if matches[4] != `24.99` {
		t.Fatal(`width wrong`)
	}

}
