package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func buildPrefix(lastItemsInFolder []bool) string {
	var prefix = ""
	for i := 1; i < len(lastItemsInFolder) - 1; i++ {
		if lastItemsInFolder[i] {
			prefix += "\t"
		} else {
			prefix += "│\t"
		}
	}
	if lastItemsInFolder[len(lastItemsInFolder) - 1] {
		prefix += "└───"
	} else {
		prefix += "├───"
	}
	return prefix
}

type ByName []os.DirEntry
func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func processNode(nestLevel int, out io.Writer, path string, printFiles bool, prefixArr []bool) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if file != nil {
		var fileStat, error = file.Stat()
		if error != nil {
			return error
		}

		var prefix = buildPrefix(prefixArr)
		if nestLevel == 0 {
			prefix = ""
		}
		if fileStat.IsDir() {
			var dirs, error = file.ReadDir(0)
			if error != nil {
				return error
			}

			if nestLevel != 0 {
				fmt.Fprintln(out, prefix+filepath.Base(file.Name()))
			}

			sort.Sort(ByName(dirs))

			if !printFiles {
				filteredDirs := []os.DirEntry{}
				for _, dir := range dirs {
					if dir.IsDir() {
						filteredDirs = append(filteredDirs, dir)
					}
				}
				dirs = filteredDirs
			}
			for i := 0; i < len(dirs); i++ {
				var newPrefixArr []bool
				if i == len(dirs) - 1 {
					newPrefixArr = append(prefixArr, true)
				} else {
					newPrefixArr = append(prefixArr, false)
				}
				processNode(nestLevel + 1, out, path + "/" + dirs[i].Name(), printFiles,  newPrefixArr)
			}
		} else {
			if printFiles {
				size := fileStat.Size()
				sizeStr := "empty"
				if size != 0 {
					sizeStr = strconv.FormatInt(size, 10) + "b"
				}
				if nestLevel != 0 {
					fmt.Fprintln(out, prefix+filepath.Base(file.Name())+" ("+sizeStr+")")
				}
			}
		}
		defer file.Close()
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var arr []bool
	arr = append(arr, true )
	return processNode(0, out, path, printFiles, arr)
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
