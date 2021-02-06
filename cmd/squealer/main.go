package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"

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
	concise        = true
	configFilePath string
	fromHash       string
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

	scanner, err := scan.NewGitScanner(scan.NewScannerConfig(basePath, redacted, cfg))
	if err != nil {
		panic(err)
	}
	err = scanner.Scan()
	if err != nil {
		panic(err)
	}
	scanner.Shutdown(concise)
}

func main() {
	rootcmd.PersistentFlags().BoolVar(&redacted, "redacted", true, "Display the results redacted")
	rootcmd.PersistentFlags().BoolVar(&concise, "concise", false, "Reduced output")
	rootcmd.PersistentFlags().StringVar(&configFilePath, "config-file", "", "Path to the config file with the rules")

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
