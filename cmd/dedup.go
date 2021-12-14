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

	"github.com/spf13/cobra"
)

// dedupCmd represents the dedup command
var dedupCmd = &cobra.Command{
	Use:   "dedup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dedup called")
		files := get_all_files_in_directory("/home/dan/Workspace/")
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
	},
}

func init() {
	rootCmd.AddCommand(dedupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dedupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dedupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

func get_all_files_in_directory(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("failed to read files in directory", err)
	}

	var file_names []string
	for _, f := range files {
		if f.IsDir() == false {
			file_names = append(file_names, dir+f.Name())
		}
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
