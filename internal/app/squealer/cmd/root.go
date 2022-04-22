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

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

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

	scanMetrics := scanner.GetMetrics()
	if !concise {
		_, _ = fmt.Fprint(os.Stderr, printMetrics(scanMetrics))
	}

	exitCode := int(math.Min(float64(scanMetrics.TransgressionsReported), 1))
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
