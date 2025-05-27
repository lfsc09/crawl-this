package file

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/lfsc09/crawl-this/chunk"
)

const (
	JSONL = "jsonl"
)

type FileStream struct {
	file    *os.File
	fileExt string
	writer  *bufio.Writer
	mu      *sync.Mutex
}

func NewFileStream(outPath string, fileExt string, mu *sync.Mutex) (*FileStream, error) {
	mu.Lock()
	if err := os.Remove(outPath); err != nil {
		if !os.IsNotExist(err) {
			mu.Unlock()
			return nil, err
		}
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		mu.Unlock()
		return nil, err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		mu.Unlock()
		return nil, err
	}
	mu.Unlock()

	return &FileStream{
		file:    outFile,
		fileExt: fileExt,
		writer:  bufio.NewWriter(outFile),
		mu:      mu,
	}, nil
}

// Close closes the file stream
func (fs *FileStream) Close() error {
	if fs.file != nil {
		fs.writer.Flush()
		return fs.file.Close()
	}
	return nil
}

// Write a chunk to the file stream in the specified format
func (fs *FileStream) Write(chunk chunk.Chunk) error {
	if fs.fileExt == JSONL {
		jsonChunk, err := json.Marshal(chunk)
		if err != nil {
			return err
		}

		fs.writer.Write(jsonChunk)
		// JSONL format requires new line after each JSON object
		fs.writer.WriteString("\n")
	}
	return nil
}
