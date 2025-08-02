package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/Krish-Parekh/GoGit/internal/object"
)

type TreeEntry struct {
	Mode fs.FileMode
	Path string
	Hash string
}

func writeTree(dir string) (string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == ".git" || filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Printf("error walking directory: %v", err)
		return "", err
	}

	var entries []TreeEntry
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("error reading file %s: %v", file, err)
			continue
		}

		hash, err := object.WriteObject(content, "blob")
		if err != nil {
			log.Printf("error writing blob for file %s: %v", file, err)
			continue
		}

		info, err := os.Stat(file)
		if err != nil {
			log.Printf("error stating file %s: %v", file, err)
			continue
		}

		entries = append(entries, TreeEntry{
			Mode: info.Mode(),
			Path: file,
			Hash: hash,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	var treeContent bytes.Buffer
	for _, entry := range entries {
		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			log.Printf("error decoding hash for %s: %v", entry.Path, err)
			continue
		}
		treeContent.WriteString(fmt.Sprintf("%o %s\x00", entry.Mode.Perm(), entry.Path))
		treeContent.Write(hashBytes)
	}

	hash, err := object.WriteObject(treeContent.Bytes(), "tree")
	if err != nil {
		log.Printf("error writing tree object: %v", err)
		return "", err
	}

	return hash, nil
}

func WriteTreeCommand() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	hash, err := writeTree(cwd)
	if err != nil {
		return err
	}
	fmt.Println(hash)
	return nil
}
