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
	"sort"
	"strings"
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

type ShellCmdMap map[int]CmdInfo

type TaliasCmdMap map[string]TaliasCmd

func (t TaliasCmdMap) updateAllStatus() TaliasCmdMap {
	taliasCmdMap := make(TaliasCmdMap)
	for k, v := range t {
		v.expire()
		v.Active = isAliasActive(v.Alias)
		taliasCmdMap[k] = v
	}
	return taliasCmdMap
}

func (t *TaliasCmd) expire() {
	if t.ExpirationDate.Before(time.Now()) && isAliasActive(t.Alias) {
		deactivateAlias(t.Alias)
	}
}

func (t *TaliasCmd) extend() {
	//
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
	delAlias			bool
	delAliasName		string
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
		filepath.Join(userHome, ".talias", "talias.conf"),
		false,
		""}

	flag.BoolVar(&appContext.listHistory, "l", false, "list history")
	flag.BoolVar(&appContext.listTaliasData, "L", false, "list aliases")
	flag.StringVar(&appContext.addAliasName, "a", "REQUIRED", "add alias <name>")
	flag.StringVar(&appContext.delAliasName, "d", "REQUIRED", "delete alias <name>")
	flag.Parse()

	if appContext.addAliasName != "REQUIRED" {
		appContext.addAlias = true
	}

	if appContext.delAliasName != "REQUIRED" {
		appContext.delAlias = true
	}

	return appContext
}

// The worlds most generic error handler ... but it gets the job done.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Check if talias in is in path
func checkPath() {
	for _, dir := range strings.Split(os.Getenv("PATH"),":") {
		if dir == ctx.AliasDir {
			return
		}
	}

	fmt.Println("Warning:", ctx.AliasDir, "is not in your path")
}

// This just returns the whole contents of a file as a string array
func readLines(path string) []string {
	file, err := os.Open(path)
	check(err)

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	check(scanner.Err())
	return lines
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

// Load shell history from file
func loadHistoryDataMap() ShellCmdMap {
	shellCmdMap := make(ShellCmdMap)
	historyLines := readLines(ctx.HistFile)
	currentTimeStamp := int64(0)
	for _, line := range historyLines {
		commentCheck := isTimeStamp(line)
		if commentCheck >= 0 {
			currentTimeStamp = commentCheck
		} else {
			mapIndex := len(shellCmdMap) + 1
			shellCmdMap[mapIndex] = CmdInfo{
				line,
				mapIndex,
				currentTimeStamp}
		}
	}
	return shellCmdMap
}

// List last N shell commands from history
func (m ShellCmdMap) listHistory(page int) {
	for i := len(m) + 1 - (ctx.ListHistoryNumber * page);
		i <= len(m) - (ctx.ListHistoryNumber * (page - 1));
		i++ {
			fmt.Println(m[i].commandNumber, " | ", m[i].command)
	}
	// currently not printing time.Unix(m[i].timestamp, 0) but may use later
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
func (taliasData TaliasCmdMap) writeDataFile() {
	taliasJson, err := json.Marshal(taliasData)
	check(err)
	err = ioutil.WriteFile(ctx.DataFile, taliasJson, 0644)
	check(err)
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
	fmt.Print("Enter command number (or 'p' for previous ", ctx.ListHistoryNumber, ") : ")
	text, _ := reader.ReadString('\n')
	if text[:len(text)-1] == "p" {
		return -1
	}
	cmdNum, err := strconv.Atoi(text[:len(text)-1])
	if err != nil {
		return -2
	}
	return cmdNum
}

// Add alias script
func addAliasScript(info CmdInfo, alias string) bool {
	aliasFile := filepath.Join(ctx.AliasDir, alias)
	f, err := os.Create(aliasFile)
	check(err)

	defer f.Close()

	f.WriteString("#!/bin/bash\n")
	f.WriteString("set -e\n")
	f.WriteString(info.command + " \"$@\"" + "\n")

	f.Sync()

	os.Chmod(aliasFile, 0755)

	return true
}

// Adds an alias to the database and creates its script
func addAlias(cmdMap ShellCmdMap, taliasData TaliasCmdMap) {
	cmdNum, page := -1, 0
	for {
		if cmdNum == -1 {
			page = page + 1
			cmdMap.listHistory(page)
			cmdNum = readInput()
		} else if cmdNum > len(cmdMap) || cmdNum < -1 {
			fmt.Println("ERROR: Please select a number below ", len(cmdMap) +1, " or Ctrl-c" )
			cmdMap.listHistory(page)
			cmdNum = readInput()
		} else {
			break
		}
	}
	addAliasScript(cmdMap[cmdNum], ctx.addAliasName)
	newAlias := TaliasCmd{cmdMap[cmdNum].command,
		ctx.addAliasName,
		time.Now(),
		time.Now().Add(time.Hour * ctx.Expiration),
		true}
	taliasData[ctx.addAliasName] = newAlias
	taliasData.writeDataFile()
	fmt.Println(cmdMap[cmdNum].command + " aliased as " + ctx.addAliasName)
}

func deactivateAlias(alias string) {
	err := os.Remove(filepath.Join(ctx.AliasDir, alias))
	check(err)
}

func delAlias(taliasData TaliasCmdMap) {
	if isAliasActive(ctx.delAliasName) {
		deactivateAlias(ctx.delAliasName)
	}
	delete(taliasData, ctx.delAliasName)
	taliasData.writeDataFile()
}

// Make a directory if it doesn't already exist
func mkDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

// List Talias metadata
func (taliasData TaliasCmdMap) listTaliasData() {
	// We want the print out consistent so we need to get the keys and print them in order
	var keys []string
	for k := range taliasData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Registered Commands =======================================")
	for _, k := range keys {
		talias := taliasData[k]
		fmt.Println("alias:", talias.Alias, "\n",
			"command: ", talias.Command, "\n",
			"expired: ", talias.ExpirationDate.Before(time.Now()), "\n",
			"active: ", talias.Active, "\n",
			"==========================================================")
	}
}

func main() {
	// Remember ctx is global
	ctx = initTaliasContext()
	writeConfFile(&ctx)

	checkPath()

	mkDir(ctx.TaliasHome)
	mkDir(ctx.AliasDir)

	taliasData := loadDataFile().updateAllStatus()
	cmdMap := loadHistoryDataMap()

	if ctx.listHistory {
		cmdMap.listHistory(1)
	}

	if ctx.listTaliasData {
		taliasData.listTaliasData()
	}

	if ctx.addAlias {
		addAlias(cmdMap, taliasData)
		os.Exit(0)
	}

	if ctx.delAlias {
		delAlias(taliasData)
		os.Exit(0)
	}

}
