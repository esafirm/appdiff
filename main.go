package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"

	"github.com/esafirm/appdiff/zipper"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"

	NewFile  = "\033[1;33m[!] %s : New File!\033[0m\n"
	Increase = "\033[1;33m[>] %s : %d => %d\033[0m\n"
	Decrease = "\033[1;34m[<] %s : %d => %d\033[0m\n"
	Same     = "\033[1;36m[=] %s : %d => %d\033[0m\n"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: apkdiff <new_apk> <old_apk>")
		return
	}

	firstApk := os.Args[1]
	secondApk := os.Args[2]

	firstDir, _ := ioutil.TempDir("", "apk")
	secondDir, _ := ioutil.TempDir("", "apk")

	err := unzip(firstApk, firstDir)
	err = unzip(secondApk, secondDir)

	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := ioutil.ReadDir(firstDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	allData := make([]string, len(files))

	fmt.Println("Comparing files…")

	for index, f := range files {
		var secondDirFileName = filepath.Join(secondDir, f.Name())
		var secondFileInfo, err = os.Stat(secondDirFileName)
		var secondSize = secondFileInfo.Size()

		var name = f.Name()
		var firstSize = f.Size()

		allData[index] = fmt.Sprintf("%s, %d,, %s, %d, %d\n", name, firstSize, name, secondSize, firstSize-secondSize)

		if err != nil {
			fmt.Printf(NewFile, name)
			continue
		}

		if firstSize > secondSize {
			fmt.Printf(Increase, name, firstSize, secondSize)
		} else if firstSize < secondSize {
			fmt.Printf(Decrease, name, firstSize, secondSize)
		} else {
			fmt.Printf(Same, name, firstSize, secondSize)
		}
	}

	copyToClipboard(allData)
}

func copyToClipboard(allData []string) {
	var buffer bytes.Buffer

	for _, data := range allData {
		buffer.WriteString(data)
	}

	clipboard.WriteAll(buffer.String())

	fmt.Println("\n\nAll data has been copied to clipboard!")
}

func unzip(path string, destDir string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	_, err = zipper.Unzip(absPath, destDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
