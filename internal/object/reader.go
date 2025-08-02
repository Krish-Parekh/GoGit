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

/*
The hash comes from the sha1 of the content. If you take the zlib compressed data and pipe it through sha1sum, you get the filename.

Hash is done using SHA1SUM algorithm, for which the output size is 20 bytes (160 bits/40 hex characters).
The first 2 characters of the hash are used to create the folder name.
The remaining 38 characters are used to create the file name.

The objects are stored in the .git/objects folder.
for example, if hash is f713b3c87b42cd63f791a27aff9743ea990f89fb, then the object is stored in the .git/objects/f7/13b3c87b42cd63f791a27aff9743ea990f89fb folder.


Git makes use of zlib to compress the objects before storing them (https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)

For us to read the object, we need to decompress the data and then read the header and the content.
Structure of the object:
<object type> <space> <content length> <null byte> <content>
Example:
blob 10\0Hello World
*/

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
