package varchive

import (
	"flag"
	"fmt"
	"os"
)

type Settings struct {
	paths                []string
	dryRun               bool
	verbose              bool
	outputRoot           string
	logToFile            string
	consoleOutputAllowed bool
	singleThread         bool
	maxParallelTasks     int
	width                string
	height               string
	fps                  string
	quality              int
	fixAudio             bool
	decomb               bool
	reportSizes          bool
}

var settings = Settings{}

func ParseArguments() {

	flag.BoolVar(&settings.verbose, "verbose", false, "be verbose")
	flag.BoolVar(&settings.dryRun, "dryRun", false, "don't affect anything")
	flag.BoolVar(&settings.singleThread, "singleThread", false, "do not parallelise tasks")

	flag.IntVar(&settings.maxParallelTasks, "maxParallelTasks", 4, "maximum number of tasks to have running at any one time.\n")

	flag.BoolVar(&settings.decomb, "decomb", false, "use de-interlacing")
	flag.BoolVar(&settings.fixAudio, "fixAudio", false, "attempt to repair the dodgy audio found on some older files")
	flag.IntVar(&settings.quality, "quality", 20, "encode quality.\nSmaller numbers are better quality, but slower to encode\n")

	flag.BoolVar(&settings.reportSizes, "reportSizes", false, "scan all files and report their video geometry.\nDoes not do any transcoding.")

	flag.StringVar(&settings.outputRoot, "outputRoot", "out",
		"location for output files.\nWill be created if required.\n")

	flag.StringVar(&settings.logToFile, "log", "",
		"location for log files.\n (default is nothing, i.e. log to standard output)")

	flag.StringVar(&settings.fps, "fps", "",
		"frames-per-second for output file.\n (default is 'do not adjust')")

	flag.StringVar(&settings.width, "width", "",
		"pixel width of output file.\n (default is 'do not adjust')")

	flag.StringVar(&settings.height, "height", "",
		"pixel height of output file.\n (default is 'do not adjust')")

	flag.Parse()

	settings.paths = flag.Args()

	if len(settings.paths) == 0 {
		fmt.Println("At least one path is required.\n\nExciting options include:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if settings.singleThread {
		settings.maxParallelTasks = 1
	}

	if settings.maxParallelTasks < 1 {
		fatal("--maxParallelTasks must be 1 or more")
	}

	// special override when we know the ncurses based output is not active
	settings.consoleOutputAllowed = settings.reportSizes
}
