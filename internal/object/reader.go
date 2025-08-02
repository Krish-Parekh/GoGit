package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadObject(hash string) (string, []byte, error) {
	if len(hash) != 40 {
		return "", nil, fmt.Errorf("invalid hash length: %s", hash)
	}

	folder := hash[:2]
	fileName := hash[2:]
	filePath := filepath.Join(".git", "objects", folder, fileName)

	compressedContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("error reading object file %s: %w", filePath, err)
	}

	zlibReader, err := zlib.NewReader(bytes.NewReader(compressedContent))
	if err != nil {
		return "", nil, fmt.Errorf("error creating zlib reader: %w", err)
	}
	defer zlibReader.Close()

	decompressedContent, err := io.ReadAll(zlibReader)
	if err != nil {
		return "", nil, fmt.Errorf("error decompressing object content: %w", err)
	}

	nullIndex := bytes.IndexByte(decompressedContent, '\x00')
	if nullIndex == -1 {
		return "", nil, fmt.Errorf("invalid object format: could not find null byte separator")
	}

	header := decompressedContent[:nullIndex]
	content := decompressedContent[nullIndex+1:]

	headerParts := strings.SplitN(string(header), " ", 2)
	if len(headerParts) != 2 {
		return "", nil, fmt.Errorf("invalid object header format: %s", header)
	}

	objectType := headerParts[0]
	sizeStr := headerParts[1]

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", nil, fmt.Errorf("invalid object size in header: %s", sizeStr)
	}
	if size != len(content) {
		return "", nil, fmt.Errorf("object size mismatch: header says %d, actual is %d", size, len(content))
	}

	return objectType, content, nil
}
