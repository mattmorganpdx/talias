package main

import (
	"bufio"
	"fmt"
	/*"io"
	"io/ioutil"*/
	"os"
	"strconv"
	"flag"
	"time"
)

var TALIAS_DIR = "/home/mmorgan/.talias/"

// Struct to hold th shell command info
type CmdInfo struct {
	command string
	commandNumber int
	timestamp int64
}

type TaliasCmd struct {
	id int
	command string
	alias string
	initializationDate int
	expirationDate int
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
func isTimeStamp(line string) int64 {
	if len(line) != 0 {
		if string(line[0]) == "#" {
			timeStamp, err := strconv.ParseInt(line[1:], 10, 64)
			if err == nil {
				return timeStamp
			}
		}
	}
	return -1
}

func buildCmdHistory(history []string) []CmdInfo {
	var cmdInfo []CmdInfo
	// Initialize the timestamp var so we can reset it as we find it in the array
	var currentTimestamp int64
	currentTimestamp = 0
	for i := 0; i < len(history); i++  {
		line := history[i]
		commentCheck := isTimeStamp(line)
		if commentCheck >= 0  {
			currentTimestamp = commentCheck
		} else {
			lineCmd := CmdInfo{ line,
								len(cmdInfo) + 1,
								currentTimestamp}
			cmdInfo = append(cmdInfo, lineCmd)
		}
	}

	return cmdInfo
}

func buildCmdHistoryMap(cmdInfo []CmdInfo) map[int]CmdInfo {
	cmdMap := make(map[int]CmdInfo)
	for i, cmd := range cmdInfo {
		cmdMap[i + 1] = cmd
	}
	return cmdMap
}

func loadDataFile(filename string) []TaliasCmd {
	var taliasCmd []TaliasCmd
	placeHolderCmd := TaliasCmd{0,
								"ls -ltr",
								"lsltr",
								0,
								1}
	taliasCmd = append(taliasCmd, placeHolderCmd)
	return taliasCmd
}

func writeDataFile(filename string, command []TaliasCmd) bool {
	return true
}

func readInput() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter command number: ")
	text, _ := reader.ReadString('\n')
	cmdNum, err := strconv.Atoi(text[:len(text) - 1])
	check(err)
	return cmdNum
}

func addAlias(info CmdInfo, alias string) bool {
	aliasFile := TALIAS_DIR + "/" + alias
	f, err := os.Create(aliasFile)
	check(err)

	defer f.Close()

	f.WriteString("#!/bin/bash\n")
	f.WriteString("set -e\n")
	f.WriteString(info.command + "\n")

	f.Sync()

	os.Chmod(aliasFile, 0755)

	return true
}

func main() {
	histFile := "/home/mmorgan/.bash_history"

	numbPtr := flag.Bool("l", false, "list history")
	flag.Parse()

	lines, err := readLines(histFile)
	check(err)

	cmdHistory := buildCmdHistory(lines)
	cmdHistoryLength := len(cmdHistory)
	cmdHistoryMap := buildCmdHistoryMap(cmdHistory)

	// Print the last 10 commands
	if *numbPtr {
		for i := cmdHistoryLength - 10; i < cmdHistoryLength; i++ {
			fmt.Println(cmdHistory[i].commandNumber, time.Unix(cmdHistory[i].timestamp, 0), cmdHistory[i].command)
		}
	}

	taliasData := loadDataFile("/tmp/.talias")

	for _, talias := range taliasData {
		println(talias.command)
	}

	cmdNum := readInput()
	fmt.Println(cmdHistoryMap[cmdNum].command)
	addAlias(cmdHistoryMap[cmdNum], "tea")
}
