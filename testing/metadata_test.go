package main

import (
	"davidhancock.com/varchive"
	"fmt"
	"testing"
)

func Test_parseDimensionsUsingFfProbe(t *testing.T) {

	info, err := varchive.GetVideoInfoUsingFfProbe(`test-data/one/sample 2.mpg`) // path relative to this .go file

	if err == nil {
		fmt.Printf("%d x %d @ %f (%f)", info.Width, info.Height, info.Fps, info.Tbr)
	} else {
		t.Fatal(err)
	}

}

func Test_ffProbeOutputCapture(t *testing.T) {

	tests := []string{
		// this one is confusing as '0x564D4448' looks suspiciously like a dimension
		`Stream #0:0[0x1011]: Video: h264 (High) (HDMV / 0x564D4448), yuv420p(top first), 1920x1080 [SAR 1:1 DAR 16:9], 125 fps, 50 tbr, 90k tbn, 50 tbc`,
		// this one has decimals in the frame rate
		`Stream #0:0[0x1e0]: Video: mpeg1video, yuv420p(tv), 640x320 [SAR 1:1 DAR 2:1], 104857 kb/s, 25.52 fps, 24.99 tbr,...`,
		// this one doesnt have the aspect ratio part '[SAR 1:1 DAR 2:1]'
		`Stream #0:0[0x1e0]: Video: mpeg1video, yuv420p(tv), 200x17 [foo], 104857 kb/s, 1.23 fps, 0.99 tbr,...`,
	}
	
	expected := []varchive.VideoInfo{
		{Width: 1920, Height: 1080, Fps: 125, Tbr: 50},
		{Width: 640, Height: 320, Fps: 25.52, Tbr: 24.99},
		{Width: 200, Height: 17, Fps: 1.23, Tbr: .99},
	}

	for index, test := range tests {
		info, err := varchive.ParseVideoInfoFromFfProbe([]string{test})
		if err != nil {
			t.Fatal(fmt.Sprintf(`#%d: unexpected error: %s`, index, err.Error()))
		}
		if info.Width != expected[index].Width {
			t.Fatal(fmt.Sprintf(`#%d: width wrong, wanted %d, got %d`, index, expected[index].Width, info.Width))
		}
		if info.Height != expected[index].Height {
			t.Fatal(fmt.Sprintf(`#%d: height wrong, wanted %d, got %d`, index, expected[index].Height, info.Height))
		}
		if info.Fps != expected[index].Fps {
			t.Fatal(fmt.Sprintf(`#%d: fps wrong, wanted %.2f, got %2.f`, index, expected[index].Fps, info.Fps))
		}
		if info.Tbr != expected[index].Tbr {
			t.Fatal(fmt.Sprintf(`#%d: tbr wrong, wanted %.2f, got %2.f`, index, expected[index].Tbr, info.Tbr))
		}
	}
}

