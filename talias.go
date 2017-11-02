package main

import (
	"bufio"
	"fmt"
	/*"io"
	"io/ioutil"*/
	"os"
)

type CmdInfo struct {
	command string
	commandNumber int
	timestamp int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func isComment(line string) bool {
	if string(line[0]) == "#" {
		return true
	}
	return false
}

func buildCmdHistory(history []string) []CmdInfo {
	var cmdInfo []CmdInfo
	//for i, line := range history

	return cmdInfo
}

func main() {
	histFile := "/home/mmorgan/.bash_history"

	lines, err := readLines(histFile)
	check(err)

	histLength := len(lines)

	for i := (histLength - 10); i < histLength ; i += 1 {
		if isComment(lines[i]) {
			fmt.Println(i, lines[i])
		}
	}
}
