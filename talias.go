package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"flag"
	"time"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
)

// Global app context variable
var ctx TaliasContext

// Struct to hold the shell command info
type CmdInfo struct {
	command       string
	commandNumber int
	timestamp     int64
}

// Struct to hold talias metadata
type TaliasCmd struct {
	Command            string
	Alias              string
	InitializationDate time.Time
	ExpirationDate     time.Time
	Active			   bool
}

type TaliasCmdMap map[string]TaliasCmd

func (t TaliasCmdMap) updateAllStatus() TaliasCmdMap {
	taliasCmdMap := make(TaliasCmdMap)
	for k, v := range t {
		v.Active = isAliasActive(v.Alias)
		taliasCmdMap[k] = v
	}
	return taliasCmdMap
}

// Struct to hold app context
type TaliasContext struct {
	listHistory         bool
	ListHistoryNumber   int
	addAlias            bool
	addAliasName        string
	listAliases         bool
	purgeExpiredAliases bool
	HistFile            string
	TaliasHome          string
	AliasDir            string
	DataFile            string
	listTaliasData      bool
	Expiration          time.Duration
	configFile			string
}

// Initialize app context
func initTaliasContext() TaliasContext {
	userHome := os.Getenv("HOME")
	appContext := TaliasContext{
		true,
		10,
		false,
		"",
		false,
		false,
		filepath.Join(userHome, ".bash_history"),
		filepath.Join(userHome, ".talias"),
		filepath.Join(userHome, ".talias", "bin"),
		filepath.Join(userHome, ".talias", "talias.db"),
		false,
		72,
		filepath.Join(userHome, ".talias", "talias.conf")}

	flag.BoolVar(&appContext.listHistory, "l", false, "list history")
	flag.BoolVar(&appContext.listTaliasData, "L", false, "list aliases")
	flag.StringVar(&appContext.addAliasName, "a", "REQUIRED", "add alias <name>")
	flag.Parse()

	if appContext.addAliasName != "REQUIRED" {
		appContext.addAlias = true
	}

	return appContext
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

// Check if an alias in the db currently has a script in place
func isAliasActive(alias string) bool {
	fullPath := filepath.Join(ctx.AliasDir, alias)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Build the array of all the available shell history
func buildCmdHistory(history []string) []CmdInfo {
	var cmdInfo []CmdInfo
	// Initialize the timestamp var so we can reset it as we find it in the array
	var currentTimestamp int64
	currentTimestamp = 0
	for i := 0; i < len(history); i++ {
		line := history[i]
		commentCheck := isTimeStamp(line)
		if commentCheck >= 0 {
			currentTimestamp = commentCheck
		} else {
			lineCmd := CmdInfo{line,
				len(cmdInfo) + 1,
				currentTimestamp}
			cmdInfo = append(cmdInfo, lineCmd)
		}
	}

	return cmdInfo
}

// Build a map so that commands can be referenced by id number
func buildCmdHistoryMap(cmdInfo []CmdInfo) map[int]CmdInfo {
	cmdMap := make(map[int]CmdInfo)
	for i, cmd := range cmdInfo {
		cmdMap[i+1] = cmd
	}
	return cmdMap
}

// Load Json Metadata
func loadDataFile() TaliasCmdMap {
	taliasCmdMap := make(TaliasCmdMap)
	raw, err := ioutil.ReadFile(ctx.DataFile)
	if ! os.IsNotExist(err) {
		check(err)
	}

	err = json.Unmarshal(raw, &taliasCmdMap)
	return taliasCmdMap
}

// Write Json Metadata
func writeDataFile(taliasData *TaliasCmdMap) bool {
	taliasJson, err := json.Marshal(taliasData)
	check(err)
	err = ioutil.WriteFile(ctx.DataFile, taliasJson, 0644)
	check(err)
	return true
}

func writeConfFile(taliasConf *TaliasContext) bool {
	taliasJson, err := json.Marshal(taliasConf)
	check(err)
	err = ioutil.WriteFile(ctx.configFile, taliasJson, 0644)
	check(err)
	return true
}

// Read user input of command number to create alias
func readInput() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter command number: ")
	text, _ := reader.ReadString('\n')
	cmdNum, err := strconv.Atoi(text[:len(text)-1])
	check(err)
	return cmdNum
}

// Add alias script
func addAlias(info CmdInfo, alias string) bool {
	aliasFile := filepath.Join(ctx.AliasDir, alias)
	f, err := os.Create(aliasFile)
	check(err)

	defer f.Close()

	f.WriteString("#!/bin/bash\n")
	f.WriteString("set -e\n")
	f.WriteString(info.command + " $*" + "\n")

	f.Sync()

	os.Chmod(aliasFile, 0755)

	return true
}

// List last N shell commands from history
func listHistory(cmdHistoryLength int, cmdHistory []CmdInfo, cmdCount int) {
	for i := cmdHistoryLength - cmdCount; i < cmdHistoryLength; i++ {
		fmt.Println(cmdHistory[i].commandNumber, time.Unix(cmdHistory[i].timestamp, 0), cmdHistory[i].command)
	}
}

// Make a directory if it doesn't already exist
func mkDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

// List Talias metadata
func (taliasData TaliasCmdMap) listTaliasData() {
	fmt.Println("Registered Commands =======================================")
	for _, talias := range taliasData {
		fmt.Println("alias:", talias.Alias, "\n",
			"command: ", talias.Command, "\n",
			"expired: ", talias.InitializationDate.After(talias.ExpirationDate), "\n",
			"active: ", talias.Active, "\n",
			"==========================================================")
	}
}

func main() {
	// Remember ctx is global
	ctx = initTaliasContext()

	taliasData := loadDataFile().updateAllStatus()

	mkDir(ctx.TaliasHome)
	mkDir(ctx.AliasDir)

	lines, err := readLines(ctx.HistFile)
	check(err)

	cmdHistory := buildCmdHistory(lines)
	cmdHistoryLength := len(cmdHistory)
	cmdHistoryMap := buildCmdHistoryMap(cmdHistory)

	// Print the last N commands
	if ctx.listHistory {
		listHistory(cmdHistoryLength, cmdHistory, ctx.ListHistoryNumber)
	}

	if ctx.listTaliasData {
		taliasData.listTaliasData()
	}

	if ctx.addAlias {
		listHistory(cmdHistoryLength, cmdHistory, ctx.ListHistoryNumber)
		cmdNum := readInput()
		addAlias(cmdHistoryMap[cmdNum], ctx.addAliasName)
		newAlias := TaliasCmd{cmdHistoryMap[cmdNum].command,
			ctx.addAliasName,
			time.Now(),
			time.Now().Add(time.Hour * ctx.Expiration),
			true}
		taliasData[ctx.addAliasName] = newAlias
		writeDataFile(&taliasData)
		fmt.Println(cmdHistoryMap[cmdNum].command + " aliased as " + ctx.addAliasName)
	}

	writeDataFile(&taliasData)
	writeConfFile(&ctx)
}
