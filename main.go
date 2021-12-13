package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
)

func compute_sha256(file string) string {

	fileToOpen, err := os.Open(file)
	fileInfo, err := fileToOpen.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		return ""
	}

	hasher := sha256.New()
	s, err := ioutil.ReadFile(file)
	hasher.Write(s)
	if err != nil {
		log.Fatal("Failed to compute hash", err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func get_all_files_in_directory(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("failed to read file", err)
	}

	var file_names []string
	for _, f := range files {
		file_names = append(file_names, f.Name())
	}

	return file_names
}

func sym_link(from string, to string) {
	os.Remove(to)
	err := os.Symlink(from, to)
	if err != nil {
		log.Fatal("Failed symlink", err)
	}
}

func get_directories_recursively(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	return dirs
}

func main() {
	dirs := get_directories_recursively("/home/dan/Workspace")
	files := get_all_files_in_directory(".")
	has := make(map[string]string)
	dups := make(map[string]string)

	for _, file := range files {
		h := compute_sha256(file)
		if h == "" {
			continue
		}
		if _, ok := has[h]; ok {
			dups[h] = file
			continue
		}
		has[h] = file
	}

	for _, file := range dups {
		sym_link(file, has[compute_sha256(file)])
	}
}
