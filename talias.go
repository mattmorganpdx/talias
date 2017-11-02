package main

import (
	"bufio"
	"fmt"
	/*"io"
	"io/ioutil"*/
	"os"
	"strconv"
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
	currentTimestamp := 0
	for i := 0; i < len(history); i++  {
		line := history[i]
		if isComment(line) {
			timeStamp, err := strconv.Atoi(line[1:])
			check(err)
			currentTimestamp = timeStamp
		} else {
			var lineCmd CmdInfo
			lineCmd.command = line
			lineCmd.commandNumber = len(cmdInfo) + 1
			lineCmd.timestamp = currentTimestamp
			cmdInfo = append(cmdInfo, lineCmd)
		}
	}

	return cmdInfo
}

func main() {
	histFile := "/home/mmorgan/.bash_history"

	lines, err := readLines(histFile)
	check(err)

	cmdHistory := buildCmdHistory(lines)
	cmdHistoryLength := len(cmdHistory)

	for i := cmdHistoryLength - 10; i < cmdHistoryLength; i++ {
		fmt.Println(cmdHistory[i])
	}
}
