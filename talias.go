package main

import (
	"bufio"
	"fmt"
	/*"io"
	"io/ioutil"*/
	"os"
	"strconv"
)

// Struct to hold th shell command info
type CmdInfo struct {
	command string
	commandNumber int
	timestamp int
}

// The worlds most generic error handler ... but it gets the job done.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// This just returns the whole contents of a file as a string array
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

// Checks if line is a comment. This should be the timestamp of one or more commands
func isComment(line string) bool {
	if string(line[0]) == "#" {
		return true
	}
	return false
}

func buildCmdHistory(history []string) []CmdInfo {
	var cmdInfo []CmdInfo
	// Initialize the timestamp var so we can reset it as we find it in the array
	currentTimestamp := 0
	for i := 0; i < len(history); i++  {
		line := history[i]
		if isComment(line) {
			timeStamp, err := strconv.Atoi(line[1:])
			check(err)
			currentTimestamp = timeStamp
		} else {
			lineCmd := CmdInfo{	line,
								len(cmdInfo) + 1,
								currentTimestamp}
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

	// Print the last 10 commands
	for i := cmdHistoryLength - 10; i < cmdHistoryLength; i++ {
		fmt.Println(cmdHistory[i])
	}
}
