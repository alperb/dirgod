package parser

import (
	. "alperb/dirgod/arguments"
	. "alperb/dirgod/logger"
	"bufio"
	"os"
	"strings"
)

func NewDirParser(args Arguments) *DirParser {
	return &DirParser{args, -1, DirStack{}, nil}
}

type DirParser struct {
	args         Arguments
	lastTabCount int
	stack        DirStack
	logger       *Logger
}

func (dp *DirParser) readFile() {
	dp.logger.Log("Reading file " + dp.args.Filename + "...")
	readFile, err := os.Open(dp.args.Filename)

	if err != nil {
		dp.logger.Log("Error while reading file: " + err.Error())
		panic(err)
	}
	dp.logger.Log("File read successfully")

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
			dp.logger.Debug("Skipping comment line: " + line)
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

		dp.logger.Log("Creating directory: " + strings.Join(dp.stack.Stack, "/"))

		dp.stack.createPath()

		dp.logger.Log("Directory created: " + strings.Join(dp.stack.Stack, "/"))
	} else if strings.HasPrefix(line, tabPrefix+"- ") {
		filename := line[dp.lastTabCount+2:]
		dp.logger.Log("Creating file: " + strings.Join(dp.stack.Stack, "/") + "/" + filename)
		dp.stack.createFile(filename)
		dp.logger.Log("File created: " + strings.Join(dp.stack.Stack, "/") + "/" + filename)
	}
}

func (dp *DirParser) Parse() {
	dp.logger = NewLogger(dp.args)
	dp.readFile()
}
