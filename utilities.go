package varchive

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"log"
	"os/exec"
	"strings"
)

func fatal(message string) {

	fmt.Println(message)
	os.Exit(1)
}

func invoke(command string, args []string) string {

	if settings.verbose {
		log.Printf("Invoking: %s %s", command, strings.Join(args, ` `))
	}

	if ! settings.dryRun {
		output, err := exec.Command(command, args...).Output()

		if err == nil {
			if settings.verbose {
				log.Printf("Return ok, stdout: %v", string(output))
			}
			return string(output)
		} else {
			log.Fatal(fmt.Sprintf("Failed: %s %s\nErr: %v\nStdout: %v",
				command, strings.Join(args, ` `), err, string(output)))
		}
	}
	
	return ""
}

func createOutputRootIfRequired() {
	if _, err := os.Stat(settings.outputRoot); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(settings.outputRoot, 0755)
		}
	}
}

func failIfConcatenationFileAlreadyExists(path string) {
	if _, err := os.Stat(path); err == nil {
		fatal(fmt.Sprintf("%s exists, will not overwrite", path))
	}
}

func lastBitOfPath(path string) string {
	return filepath.Base(path)
}

// do we actually need this? if we use os.Command, daft filenames should not be a problem
func sanitisePath(path string) string {
	// re := regexp.MustCompile(`\s+`)
	// return re.ReplaceAllString(path, ``)
	return path
}

func makeTemporaryFile(extension string) string {
	file, err := ioutil.TempFile("", "varchive.*"+extension)
	if err != nil {
		fatal(err.Error())
	}
	defer os.Remove(file.Name())
	return file.Name()
}

func removeTemporaryFile(path string) {
	// ignore any errors (which will probably be "file not found")
	os.Remove(path)
}

func getFileExtension(path string) string {
	return filepath.Ext(path)
}

func niceSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}

func niceTime(seconds float64) string {
	if math.IsInf(seconds, +1) || seconds < 0 {
		return "---:--:--"
	}

	const spm = 60
	const sph = 60 * 60
	h, m, s := 0, 0, int64(seconds)
	for s > sph {
		h++
		s -= sph
	}
	for s > spm {
		m++
		s -= spm
	}
	return fmt.Sprintf("%03d:%02d:%02d", h, m, s)
}

func getSuitableDisplay() Display {
	if settings.dryRun {
		return NewNoOpDisplay()
	} else {
		return NewDisplay()
	}
}