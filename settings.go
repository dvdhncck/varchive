package varchive

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	flag.IntVar(&settings.maxParallelTasks, "maxParallelTasks", 4, "maximum number of tasks to have running at any one time.\n")

	flag.BoolVar(&settings.decomb, "decomb", false, "use de-interlacing")
	flag.BoolVar(&settings.fixAudio, "fixAudio", false, "attempt to repair the dodgy audio found on some older files")
	flag.IntVar(&settings.quality, "quality", 20, "encode quality.\nSmaller numbers are better quality, but slower to encode\n")

	flag.StringVar(&settings.outputRoot, "outputRoot", "out",
		"location for output files.\nWill be created if required.\n")

	flag.StringVar(&settings.logToFile, "log", "",
		"location for log files.\n (default is nothing, i.e. log to standard output)")

	flag.StringVar(&settings.width, "width", "",
		"pixel width of output files.\n (default is 'do not adjust')")

	flag.StringVar(&settings.height, "height", "",
		"pixel height of output files.\n (default is 'do not adjust')")

	flag.Parse()

	settings.paths = flag.Args()

	if len(settings.paths) == 0 {
		fmt.Println("At least one path is required.\n\nExciting options include:")
		flag.PrintDefaults()
		os.Exit(1)
	}

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
