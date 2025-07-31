package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
)

func initGitDir(baseDir string) error {
	// Create .git directory and subdirectories
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		fullPath := baseDir + string(os.PathSeparator) + dir
		// 0755 is the permission for the directory (read, write, execute for owner, read and execute for group and others)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			log.Printf("error creating directory: %v", err)
			return fmt.Errorf("error creating directory: %w", err)
		}
	}
	// Create HEAD file
	headFileContents := []byte("ref: refs/heads/main\n")
	headPath := baseDir + string(os.PathSeparator) + ".git" + string(os.PathSeparator) + "HEAD"
	// 0644 is the permission for the file (read and write for owner, read for group and others)
	if err := os.WriteFile(headPath, headFileContents, 0644); err != nil {
		log.Printf("error writing file: %v", err)
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func catFile(flag string, hash string) error {
	switch flag {
	case "-p":
		// Folder is the first 2 characters of zlib | sha1sum
		folder := hash[:2]
		// File name is the rest of the hash
		fileName := hash[2:]
		// File path is the .git/objects/folder/fileName
		filePath := ".git" + string(os.PathSeparator) + "objects" + string(os.PathSeparator) + folder + string(os.PathSeparator) + fileName

		// Read the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("error reading file: %v", err)
			return fmt.Errorf("error reading file: %w", err)
		}

		// Decompress the file
		reader, err := zlib.NewReader(bytes.NewReader(content))
		if err != nil {
			log.Printf("error creating zlib reader: %v", err)
			return fmt.Errorf("error creating zlib reader: %w", err)
		}
		defer reader.Close()

		// Decompress the file
		decompressed, err := io.ReadAll(reader)
		if err != nil {
			log.Printf("error decompressing file: %v", err)
			return fmt.Errorf("error decompressing file: %w", err)
		}

		// Find the null byte because the header is separated by a null byte
		nullIndex := bytes.IndexByte(decompressed, '\x00')
		if nullIndex == -1 {
			log.Printf("error finding null byte: %v", err)
			return fmt.Errorf("error finding null byte: %w", err)
		}

		// Print the content after the null byte
		fmt.Print(string(decompressed[nullIndex+1:]))
		return nil
	default:
		return fmt.Errorf("unknown flag: %s", flag)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatalf("usage: gogit <command> [<args>...]")
	}

	switch command := os.Args[2]; command {
	case "init":
		if err := initGitDir("."); err != nil {
			log.Fatalf("init failed: %v", err)
		}
		log.Println("initialized git directory")
	case "cat-file":
		flag := os.Args[3]
		hash := os.Args[4]
		if err := catFile(flag, hash); err != nil {
			log.Fatalf("cat-file failed: %v", err)
		}
		log.Println("cat-file successful")

	default:
		fmt.Println("command: ", command)
		log.Fatalf("unknown command: %s", command)
	}
}
