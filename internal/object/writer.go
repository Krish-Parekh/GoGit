package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func WriteObject(content []byte) (string, error) {
	objectType := "blob"

	header := fmt.Appendf(nil, "%s %d\x00", objectType, len(content))

	fullContent := append(header, content...)

	hash := sha1.Sum(fullContent)
	hashString := hex.EncodeToString(hash[:])

	folder := hashString[:2]
	fileName := hashString[2:]
	folderPath := filepath.Join(".git", "objects", folder)
	filePath := filepath.Join(folderPath, fileName)

	if _, err := os.Stat(filePath); err == nil {
		return hashString, nil
	}

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", fmt.Errorf("error creating object directory %s: %w", folderPath, err)
	}

	var compressedContent bytes.Buffer
	writer := zlib.NewWriter(&compressedContent)
	if _, err := writer.Write(fullContent); err != nil {
		return "", fmt.Errorf("error compressing object content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing zlib writer: %w", err)
	}

	if err := os.WriteFile(filePath, compressedContent.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("error writing object file %s: %w", filePath, err)
	}

	return hashString, nil
}
