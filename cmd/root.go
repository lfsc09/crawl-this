package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd(opts *CommandOptions) *cobra.Command {
	if opts == nil {
		opts = &CommandOptions{
			Out: os.Stdout,
		}
	}

	rootCmd := &cobra.Command{
		Use:   "crawlthis",
		Short: "CrawlThis is a web crawler and scraper",
		Long:  "CrawlThis is a web crawler and scraper that allows you to extract data from websites and pdfs.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Configure cobra output streams to use the custom 'Out'
	rootCmd.SetOut(opts.Out)

	rootCmd.AddCommand(NewPdfCmd(opts))
	rootCmd.AddCommand(NewHtmlCmd(opts))

	// Disable automatic call of `--help` during errors
	rootCmd.SilenceUsage = true

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd(nil)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
