package main

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	level := 0
	var lines []bool
	err = fillTree(out, path, printFiles, &level, files, lines)
	if err != nil {
		return err
	}
	return nil
}

func fillTree(out io.Writer, path string, printFiles bool, level *int, files []fs.FileInfo, lines []bool) error {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() <= files[j].Name()
	})

	for index, f := range files {
		isLastDirectory := func() bool {
			if printFiles {
				return false
			}
			for _, f := range files[index+1:] {
				if f.IsDir() {
					return false
				}
			}
			return true
		}
		isLastIndex := (index == len(files)-1) || isLastDirectory()

		err := printBranch(out, printFiles, level, f, isLastIndex, lines)
		if err != nil {
			return err
		}

		if f.IsDir() {
			// we don't add a line to the files inside the directory, that is last in outer directory
			if len(lines) <= *level {
				lines = append(lines, !isLastIndex)
			} else {
				lines[*level] = !isLastIndex
			}
			err := printDirectory(out, printFiles, level, lines, f.Name(), path+"/"+f.Name())
			if err != nil {
				return err
			}
		} else if printFiles {
			err := printFile(out, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printBranch(out io.Writer, printFiles bool, level *int, file fs.FileInfo, isLastIndex bool, lines []bool) error {
	fileIsDir := file.IsDir()
	if fileIsDir || printFiles {
		elementSign := "├───"
		lastSign := "└───"
		sign := elementSign
		if isLastIndex {
			sign = lastSign
		}
		err := addShift(out, level, lines)
		if err != nil {
			return err
		}
		_, err = out.Write([]byte(sign))
		if err != nil {
			return err
		}
	}
	return nil
}

func addShift(out io.Writer, level *int, lines []bool) error {
	for i := 0; i < *level; i++ {
		shift := "│\t"
		if !lines[i] {
			shift = "\t"
		}
		_, err := out.Write([]byte(shift))
		if err != nil {
			return err
		}
	}
	return nil
}

func printDirectory(out io.Writer, printFiles bool, level *int, lines []bool, fileName string, newPath string) error {
	_, err := out.Write([]byte(fileName + "\n"))
	if err != nil {
		return err
	}
	newFiles, err := ioutil.ReadDir(newPath)
	if err != nil {
		return err
	}
	*level += 1
	err = fillTree(out, newPath, printFiles, level, newFiles, lines)
	*level -= 1
	lines = lines[:*level]
	if err != nil {
		return err
	}
	return nil
}

func printFile(out io.Writer, file fs.FileInfo) error {
	size := strconv.FormatInt(file.Size(), 10) + "b"
	if file.Size() == 0 {
		size = "empty"
	}
	_, err := out.Write([]byte(file.Name() + " (" + size + ")\n"))
	if err != nil {
		return err
	}
	return nil
}
