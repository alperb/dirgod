package parser

import (
	. "alperb/dirgod/arguments"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func NewDirParser(args Arguments) *DirParser {
	return &DirParser{args, -1, DirStack{}}
}

type DirParser struct {
	args         Arguments
	lastTabCount int
	stack        DirStack
}

func (dp *DirParser) readFile() {
	readFile, err := os.Open(dp.args.Filename)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	file := bufio.NewScanner(readFile)
	file.Split(bufio.ScanLines)

	dp.stack = DirStack{"", make([]string, 0), false}

	dp.lastTabCount = -1

	if file.Scan() {
		rootDir := file.Text()
		dp.stack.SetRootDir(rootDir)
	}

	for file.Scan() {
		line := file.Text()

		// means it's a comment line
		if strings.HasPrefix(line, "//") {
			continue
		}

		currentTabCount := strings.Count(line, "\t")
		if currentTabCount >= dp.lastTabCount { // at the same directory stack
			dp.lastTabCount = currentTabCount
			dp.performTypeOperation(line)
		} else {
			dp.stack.Pop()
			dp.lastTabCount = -1
		}
	}

	cerr := readFile.Close()
	if cerr != nil {
		panic("Failed to close file!")
	}
}

func (dp *DirParser) performTypeOperation(line string) {
	tabPrefix := strings.Join(make([]string, dp.lastTabCount+1), "\t")

	if strings.HasPrefix(line, tabPrefix+"> ") { // means it's a directory so we continue the execution by pushing
		dp.stack.Push(line[dp.lastTabCount+2:])
		dp.stack.createPath()
	} else if strings.HasPrefix(line, tabPrefix+"- ") {
		filename := line[dp.lastTabCount+2:]
		dp.stack.createFile(filename)
	}
}

func (dp *DirParser) Parse() {
	dp.readFile()
}
