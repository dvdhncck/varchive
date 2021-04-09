package varchive

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

type Settings struct {
	paths []string

	dryRun           bool
	verbose          bool
	singleThread     bool
	outputRoot       string
	logToFile        string
	maxParallelTasks int
	encodeMode       string
	width            string
	height           string
	quality          int
	fixAudio         bool
	decomb           bool
}

var settings = Settings{}

func IsVerbose() bool {
	return settings.verbose
}

func ParseArguments() {

	validEncodeModes := []string{"basic", "decomb", "minimise"}
	validEncodeModeString := fmt.Sprintf("'%s'", strings.Join(validEncodeModes, "','"))

	flag.BoolVar(&settings.verbose, "verbose", false, "be verbose")
	flag.BoolVar(&settings.dryRun, "dryRun", false, "don't affect anything")
	flag.BoolVar(&settings.singleThread, "singleThread", false, "do not parallelise tasks")

	flag.IntVar(&settings.maxParallelTasks, "maxParallelTasks", 4, "maximum number of tasks to have running at any one time.\n  Default 4.")

	flag.BoolVar(&settings.decomb, "decomb", false, "use de-interlacing")
	flag.BoolVar(&settings.fixAudio, "fixAudio", false, "fix dodgy audio (mystery audio stream on some older files")
	flag.IntVar(&settings.quality, "quality", 20, "encode quality. Default 20.\n  Smaller numbers are better quality, but slower to encode")

	flag.StringVar(&settings.encodeMode, "encodeMode", "",
		fmt.Sprintf("encode mode.\nValid modes: %s", validEncodeModeString))

	flag.StringVar(&settings.outputRoot, "outputRoot", "out",
		"location for output files.\n  Default is './out'")

	flag.StringVar(&settings.logToFile, "log", "",
		"location for lof files.\n  Default is nothing, i.e. log to standard output")

	flag.StringVar(&settings.width, "width", "",
		"pixel width of output files.\n  Default is 'do not adjust'")

	flag.StringVar(&settings.height, "height", "",
		"pixel height of output files.\n  Default is 'do not adjust'")

	geometry := flag.String("geometry", "", "geometry of output video.\n  Default is 'do not adjust'")

	flag.Parse()

	settings.paths = flag.Args()

	if settings.verbose {
		log.Printf("Settings: %v", settings)
	}

	if settings.singleThread {
		settings.maxParallelTasks = 1
	}

	if settings.maxParallelTasks < 1 {
		fatal("--maxParallelTasks must be 1 or more")
	}

	if settings.encodeMode == "" {
		fatal("--encodeMode is required")
	}

	if notIn(validEncodeModes, settings.encodeMode) {
		fatal(fmt.Sprintf("--encodeMode must be one of %s", validEncodeModeString))
	}

	if *geometry != "" {
		fatal("too much geoms")
	}

	if settings.dryRun {
		log.Print("Dry run mode enabled")
	}
}

func notIn(options []string, thing string) bool {
	for _, o := range options {
		if o == thing {
			return false
		}
	}
	return true
}
