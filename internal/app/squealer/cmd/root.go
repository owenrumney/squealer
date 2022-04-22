package cmd

import (
	"fmt"
	"math"
	"os"

	"github.com/owenrumney/squealer/pkg/squealer"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/owenrumney/squealer/internal/pkg/formatters"
	"github.com/owenrumney/squealer/internal/pkg/metrics"
	"github.com/owenrumney/squealer/pkg/config"
)

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

func Root() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "squealer",
		Short: "Search for secrets and squeal about them",
		Long:  `Telling tales on your secret leaking`,
		RunE:  squeal,
	}
	configureFlags(rootCommand)
	return rootCommand
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func configureFlags(command *cobra.Command) {
	command.PersistentFlags().BoolVar(&redacted, "redacted", redacted, "Display the results redacted.")
	command.PersistentFlags().BoolVar(&concise, "concise", concise, "Reduced output.")
	command.PersistentFlags().BoolVar(&noGit, "no-git", noGit, "Scan as a directory rather than a git history.")
	command.PersistentFlags().BoolVar(&debug, "debug", debug, "Include debug output.")
	command.PersistentFlags().BoolVar(&everything, "everything", everything, "Scan all commits.... everywhere.")
	command.PersistentFlags().StringVar(&configFilePath, "config-file", configFilePath, "Path to the config file with the rules.")
	command.PersistentFlags().StringVar(&fromHash, "from-hash", fromHash, "The hash to work back to from the starting hash.")
	command.PersistentFlags().StringVar(&toHash, "to-hash", toHash, "The most recent hash to start with.")
	command.PersistentFlags().StringVar(&format, "output-format", format, "The format that the output should come in (default, json, sarif.")
	command.PersistentFlags().StringVar(&commitListFile, "commits-file", commitListFile, "Provide a file with the commits to check per line (git rev-list master..HEAD)")

}

func squeal(_ *cobra.Command, args []string) error {
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
		return err
	}

	scanner, err := getScanner(cfg, basePath)
	if err != nil {
		return err
	}
	transgressions, err := scanner.Scan()
	if err != nil {
		return err
	}

	output, err := formatters.GetFormatter(format).PrintTransgressions(transgressions, redacted)
	if err != nil {
		log.WithError(err).Error(err.Error())
	}

	fmt.Println(output)

	metrics := scanner.GetMetrics()
	if !concise {
		_, _ = fmt.Fprint(os.Stderr, printMetrics(metrics))
	}

	exitCode := int(math.Min(float64(metrics.TransgressionsReported), 1))
	os.Exit(exitCode)
	return nil
}

func getScanner(cfg *config.Config, basePath string) (squealer.Scanner, error) {
	scanner, err := squealer.New(
		squealer.OptionWithConfig(*cfg),
		squealer.OptionRedactedSecrets(redacted),
		squealer.OptionNoGitScan(noGit),
		squealer.OptionWithBasePath(basePath),
		squealer.OptionWithFromHash(fromHash),
		squealer.OptionWithToHash(toHash),
		squealer.OptionWithScanEverything(everything),
		squealer.OptionWithCommitListFile(commitListFile),
	)
	return *scanner, err
}

func printMetrics(metrics *metrics.Metrics) string {
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
