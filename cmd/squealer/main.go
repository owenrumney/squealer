package main

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/owenrumney/squealer/internal/app/squealer/match"
	"github.com/owenrumney/squealer/internal/app/squealer/scan"
)

var rootcmd = &cobra.Command{
	Use:   "squeal",
	Short: "Search for secrets and squeal about them",
	Long:  `Start commit searching`,
	Run:   startSquealing,
}

func startSquealing(_ *cobra.Command, args []string) {
	var basePath = "./"
	if len(args) > 0 {
		basePath = args[0]
	}
	fmt.Printf("scanning %s\n", basePath)

	mc := match.NewMatcherController()
	_ = mc.Add(`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`)
	_ = mc.Add(`(?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"]`)
	_ = mc.Add(`amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	_ = mc.Add(`(?i)github[_\-\.]?token[\s:,="\]']+?(?-i)[0-9a-zA-Z]{35,40}`)

	_ = mc.Add(`xox[baprs]-([0-9a-zA-Z]{10,48})?`)
	_ = mc.Add(`-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH) )?PRIVATE KEY( BLOCK)?-----`)
	_ = mc.Add(`https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}`)

	scanner, err := scan.New(*mc, basePath)
	if err != nil {
		panic(err)
	}
	err = scanner.Scan()
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := rootcmd.Execute(); err != nil {
		fmt.Printf(err.Error())
	}
}
