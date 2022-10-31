package parser

import (
	"fmt"
	"os"
	"bufio"
	"strings"

	//	"strings"
)

func NewDirParser(filename string) *DirParser {
	return &DirParser{filename}
}

type DirParser struct {
	filename string
}

func (dp *DirParser) readFile() {
	readFile, err := os.Open("test.dir")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	file := bufio.NewScanner(readFile)
	file.Split(bufio.ScanLines)

	stack := DirStack{"", make([]string, 0), false}

	lastTabCount := -1

	if file.Scan() {
		rootDir := file.Text()
		stack.SetRootDir(rootDir)
	}

	for file.Scan() {
		line := file.Text()

		// means it's a comment line
		if strings.HasPrefix(line, "//") {
			continue
		}

		currentTabCount :=  strings.Count(line, "\t")
		if currentTabCount >= lastTabCount { // at the same directory stack
			lastTabCount = currentTabCount
			tabPrefix := strings.Join(make([]string, lastTabCount + 1), "\t")

			if strings.HasPrefix(line, tabPrefix + "> ") { // means it's a directory so we continue the execution by pushing
				stack.Push(line[lastTabCount + 2:])
			} else if strings.HasPrefix(line, tabPrefix + "- ") {
				filename := line[lastTabCount + 2:]
				stack.createFile(filename)
			}

		} else {
			stack.Pop()
		}
	}

	cerr := readFile.Close()
	if cerr != nil {
		panic("Failed to close file!")
	}
}
func (dp *DirParser) Parse() {
	dp.readFile()
}
