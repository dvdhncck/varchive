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
	singleThread bool

	outputRoot string

	maxParallelTasks int

	encodeMode string
	quality    int
	fixAudio   bool
	decomb     bool
}

func ParseArguments() *Settings {

	validEncodeModes := []string{"basic", "decomb", "minimise"}
	validEncodeModeString := fmt.Sprintf("'%s'", strings.Join(validEncodeModes, "','"))

	settings := new(Settings)

	flag.BoolVar(&settings.verbose, "verbose", false, "be verbose")
	flag.BoolVar(&settings.dryRun, "dryRun", false, "don't affect anything")
	flag.BoolVar(&settings.singleThread, "singleThread", false, "do not parallelise tasks")

	flag.IntVar(&settings.maxParallelTasks, "maxParallelTasks", 4, "maximum number of tasks to have running at any one time.\n  Default 4.")
	
	flag.BoolVar(&settings.decomb, "decomb", false, "use de-interlacing")
	flag.BoolVar(&settings.fixAudio, "fixAudio", false, "fix dodgy audio (mystery audio stream on some older files")
	flag.IntVar(&settings.quality, "quality", 20, "encode quality. Default 20.\n  Smaller numbers are better quality, but slower to encode")

	flag.StringVar(&settings.encodeMode, "encodeMode", "",
		fmt.Sprintf("encode mode.\nValid modes: %s", validEncodeModeString))

	flag.StringVar(&settings.outputRoot, "outputRoot", "varchive",
		fmt.Sprint("location for output files.\n  Default is './varchive'"))

	flag.Parse()

	settings.paths = flag.Args()

	if settings.verbose {
		log.Printf("%v", settings)
	}

	if settings.singleThread {
		settings.maxParallelTasks = 1
	}

	if settings.maxParallelTasks < 1 {
		fatal("--maxParallelTasks must be 1 or more")
	}

	if settings.dryRun {
		log.Print("DRY RUN")
	}
	
	if settings.encodeMode == "" {
		fatal("--encodeMode is required")
	}

	if notIn(validEncodeModes, settings.encodeMode) {
		fatal(fmt.Sprintf("--encodeMode must be one of %s", validEncodeModeString))
	}
	
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
