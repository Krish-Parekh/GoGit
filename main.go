package main

import (
	"fmt"
	"log"
	"os"
)

func initGitDir(baseDir string) error {
	// Create .git directory and subdirectories
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		fullPath := baseDir + string(os.PathSeparator) + dir
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			log.Printf("error creating directory: %v", err)
			return fmt.Errorf("error creating directory: %w", err)
		}
	}
	// Create HEAD file
	headFileContents := []byte("ref: refs/heads/main\n")
	headPath := baseDir + string(os.PathSeparator) + ".git" + string(os.PathSeparator) + "HEAD"
	if err := os.WriteFile(headPath, headFileContents, 0644); err != nil {
		log.Printf("error writing file: %v", err)
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
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
	default:
		fmt.Println("command: ", command)
		log.Fatalf("unknown command: %s", command)
	}
}
