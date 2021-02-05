package main

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/owenrumney/squealer/internal/app/squealer/config"
	"github.com/owenrumney/squealer/internal/app/squealer/scan"
)

var rootcmd = &cobra.Command{
	Use:   "squeal",
	Short: "Search for secrets and squeal about them",
	Long:  `Start commit searching`,
	Run:   startSquealing,
}

var (
	redacted       = false
	configFilePath string
)

func startSquealing(_ *cobra.Command, args []string) {
	var basePath = "./"
	if len(args) > 0 {
		basePath = args[0]
	}
	fmt.Printf("scanning %s\n", basePath)

	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	scanner, err := scan.NewGitScanner(basePath, cfg)
	if err != nil {
		panic(err)
	}
	err = scanner.Scan()
	if err != nil {
		panic(err)
	}

	scanner.ShowMetrics()
}

func main() {
	rootcmd.PersistentFlags().BoolVar(&redacted, "redacted", false, "Display the results redacted")
	rootcmd.PersistentFlags().StringVar(&configFilePath, "config-file", "", "Path to the config file with the rules")

	if err := rootcmd.Execute(); err != nil {
		fmt.Printf(err.Error())
	}
}
