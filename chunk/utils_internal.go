package chunk

type Chunk struct {
	ChunkIndex int    `json:"chunk_index"`
	Text       string `json:"text"`
	// Metadata map[string]any `json:"metadata,omitempty"`
	Metadata map[string]any `json:"metadata"`
}

type FileMeta struct {
	Source string
	Page   int
}
