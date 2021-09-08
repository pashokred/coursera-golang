package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	var files []os.FileInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir failed: %w", err)
	}
	level := 0
	var lines []bool
	err = fillTree(out, path, printFiles, &level, files, lines)
	if err != nil {
		return fmt.Errorf("fillTree failed: %w", err)
	}
	return nil
}

func fillTree(out io.Writer, path string, printFiles bool, level *int, files []os.FileInfo, lines []bool) error {
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
				return fmt.Errorf("printDirectory on level %v and with directory name %s failed: %w", *level, f.Name(), err)
			}
		} else if printFiles {
			err := printFile(out, f)
			if err != nil {
				return fmt.Errorf("printFIle on level %v and with name %s failed: %w", *level, f.Name(), err)
			}
		}
	}
	return nil
}

func printBranch(out io.Writer, printFiles bool, level *int, file os.FileInfo, isLastIndex bool, lines []bool) error {
	fileIsDir := file.IsDir()
	if fileIsDir || printFiles {
		elementSign, lastSign := "├───", "└───"
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
			return fmt.Errorf("out.Write on level %v and with file name %s failed: %w", *level, file.Name(), err)
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
			return fmt.Errorf("addShift out.Write on level %v failed: %w", *level, err)
		}
	}
	return nil
}

func printDirectory(out io.Writer, printFiles bool, level *int, lines []bool, fileName string, newPath string) error {
	_, err := out.Write([]byte(fileName + "\n"))
	if err != nil {
		return fmt.Errorf("printDirectory out.Write on level %v and with file name %s failed: %w", *level, fileName, err)
	}
	newFiles, err := ioutil.ReadDir(newPath)
	if err != nil {
		return fmt.Errorf("printDirectory ioutil.ReadDir on level %v and with file name %s failed: %w", *level, fileName, err)
	}
	*level += 1
	err = fillTree(out, newPath, printFiles, level, newFiles, lines)
	if err != nil {
		return err
	}
	*level -= 1
	lines = lines[:*level]

	return nil
}

func printFile(out io.Writer, file os.FileInfo) error {
	size := fmt.Sprintf("%vb", file.Size())
	if file.Size() == 0 {
		size = "empty"
	}
	_, err := out.Write([]byte(file.Name() + " (" + size + ")\n"))
	if err != nil {
		return fmt.Errorf("printDirectory out.Write with file name %s failed: %w", file.Name(), err)
	}
	return nil
}
