package main

import (
	"log"
	"flag"
	"fmt"
	"strings"
)

type RunFlags struct {
	paths      []string

	dryRun     bool
	encodeMode string
	quality    int
	fixAudio   bool
	decomb     bool

	validEncodeModes      []string
	validEncodeModeString string
}

func (rf *RunFlags) validate() {
	if rf.encodeMode == "" {
		fatal("--encodeMode is required")
	}

	if notIn(rf.validEncodeModes, rf.encodeMode) {
		fatal(fmt.Sprintf("--encodeMode must be one of %s", rf.validEncodeModeString))
	}
}

func ParseArguments() *RunFlags {
	runFlags := new(RunFlags)

	runFlags.validEncodeModes = []string{"basic", "decomb", "minimise"}
	runFlags.validEncodeModeString = fmt.Sprintf("'%s'", strings.Join(rf.validEncodeModes, "','"))
	
	runFlags.dryRun = *(flag.Bool("dryRun", false, "don't affect anything"))
	runFlags.decomb = *(flag.Bool("decomb", false, "use de-interlacing"))
	runFlags.fixAudio = *(flag.Bool("fixAudio", false, "fix dodgy audio (mystery audio stream on some older files"))
	runFlags.quality = *(flag.Int("quality", 20, "encode quality. Default 20. Smaller numbers are better quality, but slower to encode"))

	flag.StringVar(&runFlags.encodeMode, "encodeMode", "",
		fmt.Sprintf("encode mode.\nValid modes: %s", runFlags.validEncodeModeString))

	flag.Parse()

	runFlags.paths = flag.Args()

	log.Printf("%v", runFlags)

	runFlags.validate()

	return runFlags
}

func notIn(options []string, thing string) bool {
	for _, o := range options {
		if o == thing {
			return false
		}
	}
	return true
}
