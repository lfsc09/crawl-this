package cmd

import "github.com/spf13/cobra"

func NewHtmlCmd(opts *CommandOptions) *cobra.Command {
	htmlCmd := &cobra.Command{
		Use:   "html",
		Short: "Crawl a html page.",
		Long: ` Crawl a html page extracting its content and transforming to an output.
		`,
	}

	htmlCmd.Flags().String("link", "", "a start link to crawl")
	htmlCmd.Flags().Int("depth", 2, "depth of the crawl, how many links to follow from the start link")
	htmlCmd.Flags().String("chunk-strategy", "", "the strategy to use for chunking the PDF content (e.g., 'by-sentence', 'by-paragraph')")
	htmlCmd.Flags().Int("chunk-size", 200, "the size of each chunk in characters")
	htmlCmd.Flags().String("output-folder", "out", "folder to save the crawled content, if not specified, it will be saved in the '/out' directory")

	// Configure cobra output streams to use the custom 'Out'
	htmlCmd.SetOut(opts.Out)

	return htmlCmd
}
