package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	pathsToInclude := os.Args[1:]
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("ERROR: Unable to determine cwd.")
		os.Exit(1)
	}

	fileHashSum := []byte{}

	for i := 0; i < len(pathsToInclude); i++ {
		dir := path.Join(cwd, pathsToInclude[i])
		err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("ERROR: Walking directory.")
				return err
			}

			// Ignore .git folder for it contains large pack files.
			if file.IsDir() && file.Name() == ".git" {
				return filepath.SkipDir
			}

			hashContent := []byte{}
			// Include hash up to this point.
			hashContent = append(hashContent, fileHashSum...)
			// Include file/folder name in hash.
			hashContent = append(hashContent, []byte(file.Name())...)

			if file.IsDir() || file.Mode()&os.ModeSymlink != 0 {
				fmt.Printf("\033[2K\r%s %s", "Processing:", file.Name())
			} else {
				fileContents, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Println("ERROR: Reading file contents:", path)
					return err
				}
				// Include file contents in hash.
				hashContent = append(hashContent, fileContents...)
			}
			sum := md5.Sum(hashContent)
			fileHashSum = sum[:]
			return nil
		})

		if err != nil {
			fmt.Println("ERROR: Unable to hash folders and files.")
			fmt.Println(err)
			os.Exit(1)
		}
	}

	masterHash := md5.Sum(fileHashSum)
	fmt.Printf("\033[2K\r%s %x\n", "Signature:", masterHash)
	fmt.Printf("Task completed in %s\n", time.Since(start))
}
