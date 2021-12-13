package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
)

func compute_sha256(file string) string {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(file)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func get_all_files_in_directory(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
	println(len(dirs))
	files := get_all_files_in_directory(".")
	has := make(map[string]string)
	dups := make(map[string]string)

	for _, file := range files {
		h := compute_sha256(file)
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
