package main

import (
	"log"
	"os"

	"github.com/Krish-Parekh/GoGit/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatalf("usage: gogit <command> [<args>...]")
	}
	command := os.Args[2]
	switch command {
	case "init":
		if err := cmd.InitGitDirectoryCommand("."); err != nil {
			log.Fatalf("init failed: %v", err)
		}
		log.Println("Initialized empty Git repository in ./.git/")
	case "cat-file":
		if len(os.Args) < 5 {
			log.Fatalf("usage: gogit cat-file <flag> <hash>")
		}
		flag := os.Args[3]
		hash := os.Args[4]
		if err := cmd.CatFileCommand(flag, hash); err != nil {
			log.Fatalf("cat-file failed: %v", err)
		}
	case "hash-object":
		if len(os.Args) < 5 {
			log.Fatalf("usage: gogit hash-object <flag> <file>")
		}
		flag := os.Args[3]
		filePath := os.Args[4]
		if err := cmd.HashObjectCommand(flag, filePath); err != nil {
			log.Fatalf("hash-object failed: %v", err)
		}
	case "ls-tree":
		if len(os.Args) < 3 {
			log.Fatalf("usage: gogit ls-tree <hash> [<name-only>]")
		}
		hash := os.Args[3]
		nameOnly := len(os.Args) > 4 && os.Args[4] == "--name-only"
		if err := cmd.LsTreeCommand(hash, nameOnly); err != nil {
			log.Fatalf("ls-tree failed: %v", err)
		}
	case "write-tree":
		if err := cmd.WriteTreeCommand(); err != nil {
			log.Fatalf("write-tree failed: %v", err)
		}
	}
}
