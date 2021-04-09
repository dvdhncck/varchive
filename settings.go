package varchive

import (
	"flag"
	"log"
)

type Settings struct {
	paths            []string
	dryRun           bool
	verbose          bool
	outputRoot       string
	logToFile        string
	singleThread     bool
	maxParallelTasks int
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

	flag.BoolVar(&settings.verbose, "verbose", false, "be verbose")
	flag.BoolVar(&settings.dryRun, "dryRun", false, "don't affect anything")
	flag.BoolVar(&settings.singleThread, "singleThread", false, "do not parallelise tasks")

	flag.IntVar(&settings.maxParallelTasks, "maxParallelTasks", 4, "maximum number of tasks to have running at any one time.\n  Default 4.")

	flag.BoolVar(&settings.decomb, "decomb", false, "use de-interlacing")
	flag.BoolVar(&settings.fixAudio, "fixAudio", false, "fix dodgy audio (mystery audio stream on some older files")
	flag.IntVar(&settings.quality, "quality", 20, "encode quality. Default 20.\n  Smaller numbers are better quality, but slower to encode")

	flag.StringVar(&settings.outputRoot, "outputRoot", "out",
		"location for output files.\n  Default is './out'")

	flag.StringVar(&settings.logToFile, "log", "",
		"location for lof files.\n  Default is nothing, i.e. log to standard output")

	flag.StringVar(&settings.width, "width", "",
		"pixel width of output files.\n  Default is 'do not adjust'")

	flag.StringVar(&settings.height, "height", "",
		"pixel height of output files.\n  Default is 'do not adjust'")

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

	if settings.dryRun {
		log.Print("Dry run mode enabled")
	}
}
