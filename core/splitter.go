package core

func SplitByLines(records []Record, linesPerFile int) [][]Record {
	if linesPerFile <= 0 {
		linesPerFile = 1
	}
	var chunks [][]Record
	for i := 0; i < len(records); i += linesPerFile {
		end := i + linesPerFile
		if end > len(records) {
			end = len(records)
		}
		chunks = append(chunks, records[i:end])
	}
	return chunks
}

func SplitByFiles(records []Record, numFiles int) [][]Record {
	if numFiles <= 0 {
		numFiles = 1
	}
	total := len(records)
	base := total / numFiles
	remainder := total % numFiles

	var chunks [][]Record
	start := 0
	for i := 0; i < numFiles; i++ {
		size := base
		if i < remainder {
			size++
		}
		if start >= total {
			chunks = append(chunks, []Record{})
			continue
		}
		end := start + size
		if end > total {
			end = total
		}
		chunks = append(chunks, records[start:end])
		start = end
	}
	return chunks
}
