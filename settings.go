package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

type Settings struct {
	paths []string

	dryRun     bool
	verbose    bool
	encodeMode string
	quality    int
	fixAudio   bool
	decomb     bool

	validEncodeModes      []string
	validEncodeModeString string
}

func (rf *Settings) validate() {
	if rf.encodeMode == "" {
		fatal("--encodeMode is required")
	}

	if notIn(rf.validEncodeModes, rf.encodeMode) {
		fatal(fmt.Sprintf("--encodeMode must be one of %s", rf.validEncodeModeString))
	}
}

func ParseArguments() *Settings {
	settings := new(Settings)

	settings.validEncodeModes = []string{"basic", "decomb", "minimise"}
	settings.validEncodeModeString = fmt.Sprintf("'%s'", strings.Join(settings.validEncodeModes, "','"))

	settings.dryRun = *(flag.Bool("dryRun", false, "don't affect anything"))
	settings.verbose = *(flag.Bool("verbose", false, "be verbose"))

	settings.decomb = *(flag.Bool("decomb", false, "use de-interlacing"))
	settings.fixAudio = *(flag.Bool("fixAudio", false, "fix dodgy audio (mystery audio stream on some older files"))
	settings.quality = *(flag.Int("quality", 20, "encode quality. Default 20. Smaller numbers are better quality, but slower to encode"))

	flag.StringVar(&settings.encodeMode, "encodeMode", "",
		fmt.Sprintf("encode mode.\nValid modes: %s", settings.validEncodeModeString))

	flag.Parse()

	settings.paths = flag.Args()

	if settings.verbose {
		log.Printf("%v", settings)
	}

	settings.validate()

	return settings
}

func notIn(options []string, thing string) bool {
	for _, o := range options {
		if o == thing {
			return false
		}
	}
	return true
}
