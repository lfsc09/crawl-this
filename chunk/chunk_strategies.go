package chunk

func ChunkTextSegment(text string, chunkSize int, addMeta any) []Chunk {
	var chunks []Chunk
	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize {
		end := min(i+chunkSize, len(runes))
		chunkText := string(runes[i:end])
		var chunkMeta map[string]any

		switch addMeta.(type) {
		case FileMeta:
			chunkMeta["source"] = chunkMeta["source"].(string)
			chunkMeta["page"] = chunkMeta["page"].(int)
		}
		chunkMeta["length"] = len(chunkText)

		chunks = append(chunks, Chunk{
			ChunkIndex: len(chunks) + 1,
			Text:       chunkText,
			Metadata:   chunkMeta,
		})
	}

	return chunks
}
