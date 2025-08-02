package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Krish-Parekh/GoGit/internal/object"
)

func HashObjectCommand(flag string, filePath string) error {
	if flag != "-w" {
		return fmt.Errorf("unknown flag for hash-object: %s (only -w is supported)", flag)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading file %s: %v", filePath, err)
		return err
	}

	hash, err := object.WriteObject(content)
	if err != nil {
		log.Printf("error writing object for file %s: %v", filePath, err)
		return err
	}

	fmt.Println(hash)

	return nil
}
