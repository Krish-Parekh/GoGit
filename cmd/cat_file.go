package cmd

import (
	"fmt"
	"log"

	"github.com/Krish-Parekh/GoGit/internal/object"
)

/*
-p: print the contents of the object
-t: print the type of the object
-s: print the size of the object
*/

func CatFileCommand(flag string, hash string) error {
	switch flag {
	case "-p":
		_, content, err := object.ReadObject(hash)
		if err != nil {
			log.Printf("error reading object %s: %v", hash, err)
			return err
		}
		fmt.Print(string(content))
	case "-t":
		objectType, _, err := object.ReadObject(hash)
		if err != nil {
			log.Printf("error reading object %s: %v", hash, err)
			return err
		}
		fmt.Println(objectType)
	case "-s":
		_, content, err := object.ReadObject(hash)
		if err != nil {
			log.Printf("error reading object %s: %v", hash, err)
			return err
		}
		fmt.Println(len(content))
	default:
		return fmt.Errorf("unknown flag for cat-file: %s", flag)
	}
	return nil
}
