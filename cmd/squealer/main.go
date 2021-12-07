package main

import (
	"fmt"
	"math"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/owenrumney/squealer/internal/app/squealer/formatters"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"github.com/owenrumney/squealer/internal/app/squealer/scan"
	"github.com/owenrumney/squealer/pkg/config"
)

var rootcmd = &cobra.Command{
	Use:   "squealer",
	Short: "Search for secrets and squeal about them",
	Long:  `Telling tales on your secret leaking`,
	Run:   squeal,
}

var (
	redacted       = false
	concise        = false
	noGit          = false
	debug          = false
	everything     = false
	configFilePath string
	fromHash       string
	toHash         string
	commitListFile string
	format         string
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func squeal(_ *cobra.Command, args []string) {
	if concise {
		log.SetLevel(log.FatalLevel)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	var basePath = "./"
	if len(args) > 0 {
		basePath = args[0]
	}
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		fail(err)
	}

	scanner := getScanner(cfg, basePath)
	transgressions, err := scanner.Scan()
	if err != nil {
		fail(err)
	}

	output, err := formatters.GetFormatter(format).PrintTransgressions(transgressions, redacted)
	if err != nil {
		log.WithError(err).Error(err.Error())
	}

	fmt.Printf(output)

	metrics := scanner.GetMetrics()
	if !concise {
		_, _ = fmt.Fprint(os.Stderr, printMetrics(metrics))
	}

	exitCode := int(math.Min(float64(metrics.TransgressionsReported), 1))

	log.Infof("Exit code: %d", exitCode)
	os.Exit(exitCode)
}

func getScanner(cfg *config.Config, basePath string) scan.Scanner {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:            cfg,
		Basepath:       basePath,
		Redacted:       redacted,
		NoGit:          noGit,
		FromHash:       fromHash,
		ToHash:         toHash,
		Everything:     everything,
		CommitListFile: commitListFile,
	})
	if err != nil {
		fail(err)
	}
	return scanner
}

func printMetrics(metrics *mertics.Metrics) string {
	duration, _ := metrics.Duration()
	return fmt.Sprintf(`
Processing:
  duration:     %4.2fs
  commits:      %d
  commit files: %d

transgressionMap:
  identified:   %d
  ignored:      %d
  reported:     %d

`,
		duration,
		metrics.CommitsProcessed,
		metrics.FilesProcessed,
		metrics.TransgressionsFound,
		metrics.TransgressionsIgnored,
		metrics.TransgressionsReported)
}

func main() {
	rootcmd.PersistentFlags().BoolVar(&redacted, "redacted", redacted, "Display the results redacted.")
	rootcmd.PersistentFlags().BoolVar(&concise, "concise", concise, "Reduced output.")
	rootcmd.PersistentFlags().BoolVar(&noGit, "no-git", noGit, "Scan as a directory rather than a git history.")
	rootcmd.PersistentFlags().BoolVar(&debug, "debug", debug, "Include debug output.")
	rootcmd.PersistentFlags().BoolVar(&everything, "everything", everything, "Scan all commits.... everywhere.")
	rootcmd.PersistentFlags().StringVar(&configFilePath, "config-file", configFilePath, "Path to the config file with the rules.")
	rootcmd.PersistentFlags().StringVar(&fromHash, "from-hash", fromHash, "The hash to work back to from the starting hash.")
	rootcmd.PersistentFlags().StringVar(&toHash, "to-hash", toHash, "The most recent hash to start with.")
	rootcmd.PersistentFlags().StringVar(&format, "output-format", format, "The format that the output should come in (default, json, sarif.")
	rootcmd.PersistentFlags().StringVar(&commitListFile, "commits-file", commitListFile, "Provide a file with the commits to check per line (git rev-list master..HEAD)")

	if err := rootcmd.Execute(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	log.WithError(err).Error(err.Error())
	os.Exit(-1)
}
