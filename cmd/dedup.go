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

var (
	hashDb map[string]string = nil
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
		println("Starting")
		if hashDb == nil {
			hashDb := make(map[string]string)
			hashDb["helo"] = "helo"
		}
		recurse_through_directories("/home/dan/Pictures/")
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

func recurse_through_directories(directory string) {
	//
	println("In directory" + directory)
	if hashDb == nil {
		hashDb = make(map[string]string)
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() == false {
			hash := compute_sha256(directory + f.Name())
			if hash != "" {
				if _, ok := hashDb[hash]; ok {
					println("Dup found")
					println("Original:" + hashDb[hash])
					sym_link(directory+f.Name(), hashDb[hash])
					continue
				}
				hashDb[hash] = directory + f.Name()
			}
		} else {
			recurse_through_directories(directory + f.Name() + "/")
		}
	}
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
