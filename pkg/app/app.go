package app

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/zarldev/ffakes/pkg/app/info"
	"github.com/zarldev/ffakes/pkg/generator"
)

var (
	help, version         bool
	interfaceList, output string
	verbose               bool
	err                   error
)

func Run() {
	interfaces := parseCliFlags()

	if help {
		printHelp()
		return
	}

	if version {
		printVersion()
		return
	}
	if output == "" {
		output = "."
	}
	configureSlog(verbose)

	if len(interfaces) == 0 {
		fmt.Println("you must provide at least one interface to fake")
		return
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if the file is a Go file
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse and generate fakes for each file in the directory
		err = generator.ParseAndGenerate(path, interfaces, output)
		if err != nil {
			fmt.Printf("Failed to generate fakes: %s", err.Error())
			os.Exit(1)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed to walk directory: %s", err.Error())
		os.Exit(1)
	}
}

func configureSlog(v bool) {
	// Configure slog
	// set global logger with custom options
	w := os.Stderr
	level := slog.LevelDebug
	if !v {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen,
		}),
	))
}

func parseCliFlags() []string {
	flag.BoolVar(&help, "help", false,
		"Print help information")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&version, "version", false,
		"Print version information")
	flag.BoolVar(&version, "v", false, "")
	flag.StringVar(&interfaceList, "interfaces", "",
		"Comma separated list of interfaces to fake")
	flag.StringVar(&interfaceList, "i", "", "")
	interfaces := strings.Split(interfaceList, ",")
	flag.StringVar(&output, "output", "",
		"Output directory for the generated files")
	flag.StringVar(&output, "o", "", "")
	flag.BoolVar(&verbose, "vv", false, "")
	flag.BoolVar(&verbose, "vverbose", false,
		"Print verbose output")

	flag.Parse()
	return interfaces
}

func printHelp() {
	printTitle()
	fmt.Println("Usage: ffakes [options] filename")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func printVersion() {
	printTitle()
	fmt.Printf("version: %s\n", info.AppInfo.Version())
	fmt.Printf("build date: %s\n", info.AppInfo.BuildDate())
	fmt.Printf("git sha: %s\n", info.AppInfo.GitSHA())
}

var asciiArt = ` _____ _____ _____ __ ________ _____ 
/   __/   __/  _  |  |  /   __/  ___>
|   __|   __|  _  |  _ <|   __|___  |
\__/  \__/  \__|__|__|__\_____<_____/`

func printTitle() {
	fmt.Println(asciiArt)
}
