package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

func hashObject(flag string, filePath string) error {
	switch flag {
	case "-w":
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("error reading file: %v", err)
			return fmt.Errorf("error reading file: %w", err)
		}
		// use sha1sum to hash the content
		hash := sha1.New()
		hash.Write(content)
		hashValue := hash.Sum(nil)
		hashString := hex.EncodeToString(hashValue)

		// write the hash to the .git/objects/folder/fileName
		folder := hashString[:2]
		fileName := hashString[2:]

		folderPath := ".git" + string(os.PathSeparator) + "objects" + string(os.PathSeparator) + folder
		filePath := folderPath + string(os.PathSeparator) + fileName

		// create the folder if it doesn't exist
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			log.Printf("error creating directory: %v", err)
			return fmt.Errorf("error creating directory: %w", err)
		}

		// create the file if it doesn't exist
		if _, err := os.Create(filePath); err != nil {
			log.Printf("error creating file: %v", err)
			return fmt.Errorf("error creating file: %w", err)
		}

		// compress the content
		compressed := bytes.NewBuffer(nil)
		writer, err := zlib.NewWriterLevel(compressed, zlib.NoCompression)
		if err != nil {
			log.Printf("error creating zlib writer: %v", err)
			return fmt.Errorf("error creating zlib writer: %w", err)
		}
		writer.Write(content)
		writer.Close()

		// write the compressed content to the file
		if err := os.WriteFile(filePath, compressed.Bytes(), 0644); err != nil {
			log.Printf("error writing file: %v", err)
			return fmt.Errorf("error writing file: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unknown flag: %s", flag)
	}
}

func lsTree(hash string) error {
	folder := hash[:2]
	fileName := hash[2:]

	filePath := ".git" + string(os.PathSeparator) + "objects" + string(os.PathSeparator) + folder + string(os.PathSeparator) + fileName

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading file: %v", err)
		return fmt.Errorf("error reading file: %w", err)
	}

	zlibReader, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		log.Printf("error creating zlib reader: %v", err)
		return fmt.Errorf("error creating zlib reader: %w", err)
	}
	defer zlibReader.Close()

	reader := bufio.NewReader(zlibReader)
	objectType, _ := reader.ReadString(' ')
	objectType = strings.TrimSpace(objectType)

	if objectType != "tree" {
		log.Printf("error: not a tree object: %v", err)
		return fmt.Errorf("error: not a tree object: %w", err)
	}

	lengthStr, err := reader.ReadString(0)
	lengthStr = lengthStr[:len(lengthStr)-1]
	if err != nil {
		log.Printf("error reading length string: %v", err)
		return fmt.Errorf("error reading length string: %w", err)
	}
	objSize, _ := strconv.Atoi(lengthStr)

	if objSize == 0 {
		log.Printf("error: object file %s is empty", filePath)
		return fmt.Errorf("error: object file %s is empty", filePath)
	}

	fileHash := make([]byte, 20)
	for {
		fileMode, err := reader.ReadString(' ')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error reading file mode: %v", err)
			return fmt.Errorf("error reading file mode: %w", err)
		}
		fileMode = "000000" + fileMode[:len(fileMode)-1]
		fileMode = fileMode[len(fileMode)-6:]

		name, err := reader.ReadString('\000')
		if err != nil {
			log.Printf("error reading file name: %v", err)
			return fmt.Errorf("error reading file name: %w", err)
		}
		name = name[:len(name)-1]

		_, err = reader.Read(fileHash)
		if err != nil {
			log.Printf("error reading hash: %v", err)
			return fmt.Errorf("error reading hash: %w", err)
		}

		hash := hex.EncodeToString(fileHash)
		folder := hash[:2]
		fileName := hash[2:]
		filePath := ".git" + string(os.PathSeparator) + "objects" + string(os.PathSeparator) + folder + string(os.PathSeparator) + fileName
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("error reading file: %v", err)
			return fmt.Errorf("error reading file: %w", err)
		}

		zlibReader, err := zlib.NewReader(bytes.NewReader(content))
		if err != nil {
			log.Printf("error creating zlib reader: %v", err)
			return fmt.Errorf("error creating zlib reader: %w", err)
		}

		reader := bufio.NewReader(zlibReader)
		objectType, _ := reader.ReadString(' ')
		objectType = strings.TrimSpace(objectType)
		lengthStr, _ := reader.ReadString(0)
		lengthStr = lengthStr[:len(lengthStr)-1]
		objSize, _ := strconv.Atoi(lengthStr)

		fmt.Printf("%s %s %x\t%d\t%s\n", fileMode, objectType, fileHash, objSize, name)
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
		flag := os.Args[3]
		hash := os.Args[4]
		if err := catFile(flag, hash); err != nil {
			log.Fatalf("cat-file failed: %v", err)
		}
		log.Println("cat-file successful")
	case "hash-object":
		flag := os.Args[3]
		filePath := os.Args[4]
		if err := hashObject(flag, filePath); err != nil {
			log.Fatalf("hash-object failed: %v", err)
		}
		log.Println("hash-object successful")
	case "ls-tree":
		hash := os.Args[3]
		if err := lsTree(hash); err != nil {
			log.Fatalf("ls-tree failed: %v", err)
		}
		log.Println("ls-tree successful")

	default:
		fmt.Println("command: ", command)
		log.Fatalf("unknown command: %s", command)
	}
}
