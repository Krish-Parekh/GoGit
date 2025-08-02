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

/*
Structure of the object:
<object type> <space> <content length> <null byte> <content>
Example:
blob 10\0Hello World

As discussed in reader.go, the objects are stored in the .git/objects folder.

1. The first 2 characters of the hash are used to create the folder name.
2. The remaining 38 characters are used to create the file name.
3. The content is compressed using zlib.
4. The compressed content is written to the file.

*/

func WriteObject(content []byte, objectType string) (string, error) {

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
