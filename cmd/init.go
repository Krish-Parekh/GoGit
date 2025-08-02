package cmd

import (
	"fmt"
	"log"
	"os"
)

func InitGitDirectoryCommand(baseDirectory string) error {
	// Create .git directory and subdirectories
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		fullPath := baseDirectory + string(os.PathSeparator) + dir
		// 0755 is the permission for the directory (read, write, execute for owner, read and execute for group and others)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			log.Printf("error creating directory: %v", err)
			return fmt.Errorf("error creating directory: %w", err)
		}
	}
	// Create HEAD file
	headFileContents := []byte("ref: refs/heads/main\n")
	headPath := baseDirectory + string(os.PathSeparator) + ".git" + string(os.PathSeparator) + "HEAD"
	// 0644 is the permission for the file (read and write for owner, read for group and others)
	if err := os.WriteFile(headPath, headFileContents, 0644); err != nil {
		log.Printf("error writing file: %v", err)
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
