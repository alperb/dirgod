package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DirStack struct {
	rootDir       string
	Stack         []string
	isPathCreated bool
}

func (ds *DirStack) createPath() (bool, error) {
	if ds.isPathCreated {
		return true, nil
	}
	joined := strings.Join(ds.Stack, "/")
	err := os.MkdirAll(strings.Trim(joined, " "), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Failed to create path: %s", joined))
	}
	ds.isPathCreated = true
	return true, nil
}

func (ds *DirStack) decideOperation(filename string) int {
	if strings.Count(filename, "<>") == 1 {
		return 3
	} else if strings.Count(filename, "<<") == 1 {
		return 2
	} else if strings.Count(filename, "<") == 1 {
		return 1
	}
	return 0
}

func (ds *DirStack) createFile(filename string) {
	_, perr := ds.createPath()
	if perr != nil {
		fmt.Println(perr)
		panic("Failed to create subpath!")
	}
	parsedFileName := filename

	operation := ds.decideOperation(filename)
	switch operation {
	case 1: // file content is going to be copied
		opIndex := strings.Index(filename, "<")
		parsedFileName = parsedFileName[:opIndex-1]
		copyingFilePath := filename[opIndex+2:]
		defer ds.copyFileContent(copyingFilePath, parsedFileName)
	case 2: // file is going to be moved
		opIndex := strings.Index(filename, "<<")
		parsedFileName = parsedFileName[:opIndex-1]
		movingFilePath := filename[opIndex+3:]
		defer ds.moveFile(movingFilePath, parsedFileName)
		return
	case 3: // file is going to be linked
		opIndex := strings.Index(filename, "<>")
		parsedFileName = parsedFileName[:opIndex-1]
		linkedFilePath := filename[opIndex+3:]
		defer ds.linkFile(linkedFilePath, parsedFileName)
		return
	}

	ds.Push(parsedFileName)
	joined := strings.Join(ds.Stack, "/")

	f, err := os.Create(strings.Trim(joined, " "))

	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Failed to create file: %s", filename))
	}
	defer f.Close()
	defer ds.Pop()
}

func (ds *DirStack) SetRootDir(dir string) {
	ds.rootDir = dir
	ds.Push(dir)
}

func (ds *DirStack) Push(dir string) {
	ds.Stack = append(ds.Stack, dir)
	ds.isPathCreated = false
}

func (ds *DirStack) Pop() {
	if len(ds.Stack) > 1 {
		length := len(ds.Stack)
		ds.Stack = ds.Stack[:length-1]
		ds.isPathCreated = false
	}
}

func (ds *DirStack) moveFile(from string, to string) {
	err := os.Rename(from, filepath.Join(strings.Join(ds.Stack, "/"), "/", to))
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Failed to move file: %s", from))
	}
}

func (ds *DirStack) linkFile(from string, to string) {
	err := os.Symlink(from, filepath.Join(strings.Join(ds.Stack, "/"), "/", to))
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Failed to link file: %s", from))
	}
}

func (ds *DirStack) copyFileContent(from string, to string) {
	// lets make here open the file once to reduce fd use
	f, cerr := os.OpenFile(filepath.Join(strings.Join(ds.Stack, "/"), "/", to), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if cerr != nil {
		fmt.Println(cerr)
		panic("Cannot open created file to write!")
	}
	defer f.Close()

	f2, cerr := os.ReadFile(from)
	_, err := f.WriteString(string(f2))
	if err != nil {
		fmt.Println(err)
		panic("Failed to write to file!")
	}
}
