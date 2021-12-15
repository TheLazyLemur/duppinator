/*
Copyright © 2021 Dan Rousseau <danrousseau@protonmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	hashDb map[string]string = nil
	dir    string
)

// dedupCmd represents the dedup command
var dedupCmd = &cobra.Command{
	Use:   "dedup",
	Short: "Recurse through directories and find duplicate files, replacing them with symlinks",
	Long: `Recurse through directories and find duplicate files, replacing them with symlinks of the original file.
	This command IS destructive and will replace files with symlinks... Use with caution.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dedup called")
		println("Starting")
		recurse_through_directories(dir)
	},
}

func init() {
	rootCmd.AddCommand(dedupCmd)
	dedupCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to start searching from.")
	hashDb = make(map[string]string)
}

func recurse_through_directories(directory string) {
	println("In directory" + directory)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	dirs := make([]string, 0)

	for _, f := range files {

		if should_skip_file(f.Name(), directory) {
			continue
		}

		if f.IsDir() {
			fp := filepath.Join(directory, f.Name())
			dirs = append(dirs, fp)
			continue
		}

		fp := filepath.Join(directory, f.Name())
		hash := compute_sha256(fp)
		if hash == "" {
			continue
		}

		if _, ok := hashDb[hash]; ok {
			println("Dup found")
			println("Original:" + hashDb[hash])
			fp := filepath.Join(directory, f.Name())
			sym_link(fp, hashDb[hash])
		} else {
			fp := filepath.Join(directory, f.Name())
			hashDb[hash] = fp
		}
	}

	for _, d := range dirs {
		recurse_through_directories(d)
	}
}

func should_skip_file(file string, directory string) bool {
	if file[0:1] == "." {
		return true
	}

	fp := filepath.Join(directory, file)
	fi, err := os.Lstat(fp)
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}
	return false
}

func sym_link(from string, to string) {
	os.Remove(from)
	err := os.Symlink(to, from)
	if err != nil {
		log.Fatal("Failed symlink", err)
	}
}

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
