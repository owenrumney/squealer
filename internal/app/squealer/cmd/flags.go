package cmd

import "github.com/spf13/cobra"

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

func configureFlags(command *cobra.Command) {
	command.PersistentFlags().BoolVar(&redacted, "redacted", redacted, "Display the results redacted.")
	command.PersistentFlags().BoolVar(&concise, "concise", concise, "Reduced output.")
	command.PersistentFlags().BoolVar(&noGit, "no-git", noGit, "Scan as a directory rather than a git history.")
	command.PersistentFlags().BoolVar(&debug, "debug", debug, "Include debug output.")
	command.PersistentFlags().BoolVar(&everything, "everything", everything, "Scan all commits.... everywhere.")
	command.PersistentFlags().StringVar(&configFilePath, "config-file", configFilePath, "Path to the config file with the rules.")
	command.PersistentFlags().StringVar(&fromHash, "from-hash", fromHash, "The hash to work back to from the starting hash.")
	command.PersistentFlags().StringVar(&toHash, "to-hash", toHash, "The most recent hash to start with.")
	command.PersistentFlags().StringVarP(&format, "output-format", "f", format, "The format that the output should come in (default, json, sarif.")
	command.PersistentFlags().StringVar(&commitListFile, "commits-file", commitListFile, "Provide a file with the commits to check per line (git rev-list master..HEAD)")

}
