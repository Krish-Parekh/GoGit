package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Krish-Parekh/GoGit/internal/object"
)

/*
ls-tree: list the contents of a tree object

A single tree object contains one or more entries, each of which is the SHA-1 hash of a blob or subtree with its associated mode, type, and filename

Structure of the tree object:
tree <space> <number of entries> <null byte> <entry 1> <entry 2> ... <entry n>

*/

func LsTreeCommand(hash string, nameOnly bool) error {
	objectType, content, err := object.ReadObject(hash)
	if err != nil {
		log.Printf("error reading object %s: %v", hash, err)
		return err
	}

	if objectType != "tree" {
		return fmt.Errorf("object %s is not a tree (is %s)", hash, objectType)
	}

	// Loop through all the entries in the tree object
	for len(content) > 0 {
		space := bytes.IndexByte(content, ' ')
		if space == -1 {
			return fmt.Errorf("invalid tree entry: missing space after mode")
		}
		mode := content[:space]
		content = content[space+1:]

		null := bytes.IndexByte(content, '\x00')
		if null == -1 {
			return fmt.Errorf("invalid tree entry: missing null after filename")
		}
		name := content[:null]
		content = content[null+1:]

		if len(content) < 20 {
			return fmt.Errorf("invalid tree entry: truncated SHA-1 hash")
		}
		sha1Bytes := content[:20]
		content = content[20:]

		if nameOnly {
			fmt.Println(string(name))
		} else {
			entryType, _, err := object.ReadObject(hex.EncodeToString(sha1Bytes))
			if err != nil {
				log.Printf("error reading tree entry object %s: %v", hex.EncodeToString(sha1Bytes), err)
				entryType = "unknown"
			}
			fmt.Printf("%s %s %s\t%s\n", mode, entryType, hex.EncodeToString(sha1Bytes), name)
		}
	}

	return nil
}
