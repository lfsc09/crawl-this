package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	ledongthuc_pdf "github.com/ledongthuc/pdf"
	"github.com/lfsc09/crawl-this/chunk"
	"github.com/lfsc09/crawl-this/file"
	"github.com/lfsc09/crawl-this/log"
	"github.com/spf13/cobra"
	rscio_pdf "rsc.io/pdf"
)

func NewPdfCmd(opts *CommandOptions) *cobra.Command {
	pdfCmd := &cobra.Command{
		Use:   "pdf",
		Short: "Crawl a PDF document.",
		Long: `Crawl a PDF document extracting its content and transforming it to an output format (JSONL by default).
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, _ := cmd.Flags().GetString("file")
			// chunkStrategy, _ := cmd.Flags().GetString("chunk-strategy")
			chunkSize, _ := cmd.Flags().GetInt("chunk-size")
			outputFolder, _ := cmd.Flags().GetString("output-folder")
			verbose, _ := cmd.Flags().GetBool("verbose")

			logger := log.NewLog(opts.Out, verbose)

			if filePath == "" {
				return fmt.Errorf("provide a file path using --file")
			}

			logger.Log(fmt.Sprintf("Crawling PDF file: %s\n", filePath))
			logger.Log(fmt.Sprintf("Chunking content with size %d characters\n\n", chunkSize))

			// Validate and create the output filepath and file
			var mu sync.Mutex
			filename := filepath.Base(filePath)
			outPath := outputFolder + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + file.JSONL
			fileStream, err := file.NewFileStream(outPath, file.JSONL, &mu)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			defer fileStream.Close()

			// Extract content from the PDF file
			pages, err := extracContentFromPdf_LEDONGTHUC(filePath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			logger.Log(fmt.Sprintf("Extracted %d pages from the PDF document\n", len(pages)))

			// Chunk the content and write to the file stream
			for pageNum, pageContent := range pages {
				logger.Log(fmt.Sprintf("Page [%d] content '%v'\n", pageNum+1, pageContent))
				chunks := chunk.ChunkTextSegment(pageContent, chunkSize, chunk.FileMeta{
					Source: filePath,
					Page:   pageNum + 1,
				})
				logger.Log(fmt.Sprintf("Page [%d] has %d chunks\n", pageNum+1, len(chunks)))
				for _, chunkContent := range chunks {
					fileStream.Write(chunkContent)
				}
			}

			return nil
		},
	}

	pdfCmd.Flags().String("file", "", "the file path of the PDF document to crawl")
	pdfCmd.Flags().String("chunk-strategy", "", "the strategy to use for chunking the PDF content (e.g., 'by-sentence', 'by-paragraph')")
	pdfCmd.Flags().Int("chunk-size", 200, "the size of each chunk in characters")
	pdfCmd.Flags().String("output-folder", "out", "folder to save the crawled content, if not specified, it will be saved in the '/out' directory")
	pdfCmd.Flags().Bool("verbose", false, "enable verbose output")

	// Configure cobra output streams to use the custom 'Out'
	pdfCmd.SetOut(opts.Out)

	return pdfCmd
}

// This function reads the PDF file and return its content as raw text
func extractContentFromPdf_RSCIO(filePath string) ([]string, error) {
	openedFile, err := rscio_pdf.Open(filePath)
	if err != nil {
		return nil, err
	}

	var pages []string
	for i := range openedFile.NumPage() {
		page := openedFile.Page(i + 1)
		if page.V.IsNull() {
			continue
		}
		content := page.V.RawString()
		pages = append(pages, content)
	}

	return pages, nil
}

func extracContentFromPdf_LEDONGTHUC(filePath string) ([]string, error) {
	openedFile, r, err := ledongthuc_pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	var pages []string
	totalPages := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		page := r.Page(pageIndex)
		if page.V.IsNull() || page.V.Key("Contents").Kind() == ledongthuc_pdf.Null {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		pages = append(pages, content)
	}

	return pages, nil
}
