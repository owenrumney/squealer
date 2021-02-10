package main

import (
	"fmt"
	"github.com/owenrumney/squealer/internal/app/squealer/mertics"
	"github.com/spf13/cobra"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/scan"
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
	configFilePath string
	fromHash       string
	toHash         string
)

func squeal(_ *cobra.Command, args []string) {
	var basePath = "./"
	if len(args) > 0 {
		basePath = args[0]
	}
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		fail(err)
	}

	scanner := getScanner(cfg, basePath)
	err = scanner.Scan()
	if err != nil {
		fail(err)
	}
	metrics := scanner.GetMetrics()
	if !concise {
		fmt.Println(printMetrics(metrics))
	}

	exitCode := int(math.Min(float64(metrics.TransgressionsReported), 1))

	os.Exit(exitCode)
}

func getScanner(cfg *config.Config, basePath string) scan.Scanner {
	scanner, err := scan.NewScanner(scan.ScannerConfig{
		Cfg:      cfg,
		Basepath: basePath,
		Redacted: redacted,
		NoGit:    noGit,
		FromHash: fromHash,
		ToHash:   toHash,
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
	rootcmd.PersistentFlags().BoolVar(&redacted, "redacted", false, "Display the results redacted.")
	rootcmd.PersistentFlags().BoolVar(&concise, "concise", false, "Reduced output.")
	rootcmd.PersistentFlags().BoolVar(&noGit, "no-git", false, "Scan as a directory rather than a git history.")
	rootcmd.PersistentFlags().StringVar(&configFilePath, "config-file", "", "Path to the config file with the rules.")
	rootcmd.PersistentFlags().StringVar(&fromHash, "from-hash", "", "The hash to work back to from the starting hash.")
	rootcmd.PersistentFlags().StringVar(&toHash, "to-hash", "", "The most recent hash to start with.")

	listenForExit()
	if err := rootcmd.Execute(); err != nil {
		fmt.Printf(err.Error())
	}
}

func listenForExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Exiting")
		os.Exit(0)
	}()
}

func fail(err error) {
	fmt.Println(err.Error())
	os.Exit(-1)
}
