package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/araleo/yasp/git"
	"github.com/araleo/yasp/parse"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Ignore struct {
		Path string
		File string
	}
	Print struct {
		Commands string
	}
	Todos struct {
		Commands string
	}
	Issues struct {
		Commands string
	}
	Env map[string]parse.EnvConfig
}

const PREFIX = "|_ "

// TODO go through the tree only once while searching for todos and prints and not twice as it is now.
func main() {
	parse.LoadDotEnv()
	git.LoadConfig()

	config := Config{}
	loadConfig(&config)

	regexMap := buildRegexes(&config)

	rootDir, command, ignoredNames := parseFlags(&config)

	if command == "ls" {
		depth := 0
		printDir(rootDir, &depth, ignoredNames)
	}

	if command == "diag" {
		lines := 0

		fmt.Println("\nChecking \033[35menv\033[0m files and variables...")
		checkEnvs(&config)

		fmt.Println("\nI found these \033[35mprint\033[0m statements in the code:")
		walkDir(rootDir, ignoredNames, regexMap["print"], false, &lines)

		fmt.Println("\nI found these \033[35mtodo\033[0m statements in the code:")
		walkDir(rootDir, ignoredNames, regexMap["todos"], false, nil)

		fmt.Printf("\nFor this diagnosis I went through %d lines of code.\n", lines)
	}

	if command == "env" {
		fmt.Println("\nChecking \033[35menv\033[0m files and variables...")
		checkEnvs(&config)
	}

	if command == "print" {
		fmt.Println("\nI found these \033[35mprint\033[0m statements in the code:")
		walkDir(rootDir, ignoredNames, regexMap["print"], false, nil)
	}

	if command == "todo" {
		fmt.Println("\nI found these \033[35mtodo\033[0m statements in the code:")
		walkDir(rootDir, ignoredNames, regexMap["todos"], false, nil)
	}

	if command == "snitch" {
		fmt.Println("\nReporting these unreported issues:")
		walkDir(rootDir, ignoredNames, regexMap["issues"], true, nil)
	}

	if command == "issues" {
		fmt.Println("\nCurrent GitLab issues:")
		git.ListIssues()
	}
}

// parseFlags parses the flags provided by the user and returns their respective values
func parseFlags(config *Config) (string, string, []string) {
	flagsMap := ParseFlags()

	rootDir, err := ValidateDirFlag(flagsMap["dir"])
	checkErr(err)

	command, err := ValidateCommandFlag(flagsMap["command"])
	checkErr(err)

	ignoreFile := ValidateIgnoreFlag(flagsMap["ignore"], config)
	ignoredNames := loadIgnore(ignoreFile)

	return rootDir, command, ignoredNames
}

// checkEnvs checks all envs specified in the yasp config file for their respective config files.
func checkEnvs(config *Config) {
	for envName, envConfigs := range config.Env {
		ok := parse.CheckEnv(envConfigs)
		if ok {
			fmt.Println(envName + " env files seem ok.")
		}
	}
}

// reportIssue receives a string in the todo <title> format submits it as an issue on Gitlab and returns the created issue iid.
func reportIssue(lineText string) int {
	issueTitle := strings.Split(lineText, "! ")[1]
	formatedTitle := strings.ReplaceAll(issueTitle, " ", "%20")
	createdId := git.CreateIssue(formatedTitle)
	return createdId
}

// searchForPattern receives a filepath as a string and returns all ocourrences of pattern in the content of the file
func searchForPattern(filepath string, pattern *regexp.Regexp, snitch bool, lines *int) {
	file, err := os.Open(filepath)
	checkErr(err)
	defer file.Close()

	fileDir, fileName := path.Split(filepath)
	dirName := path.Base(fileDir)

	line := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := strings.TrimSpace(scanner.Text())
		matched := pattern.FindString(lineText)

		if matched != "" {
			coloredLine := strings.ReplaceAll(lineText, matched, "\033[35m"+matched+"\033[0m")
			fmt.Println(path.Join(dirName, fileName), line, coloredLine)
			if snitch {
				reportIssue(lineText)
			}
		}
		line++
	}

	if lines != nil {
		*lines += line
	}
}

// walkDir recursively walks through the received directory checking for occurrences of pattern in every file found.
func walkDir(dirPath string, ignored []string, regexString string, snitch bool, lines *int) {
	files := getFiles(dirPath, ignored)

	pattern, err := regexp.Compile(regexString)
	checkErr(err)

	for _, file := range files {
		filepath := path.Join(dirPath, file.Name())

		if file.IsDir() {
			walkDir(filepath, ignored, regexString, snitch, lines)
		} else {
			searchForPattern(filepath, pattern, snitch, lines)
		}
	}
}

// printDir recursively prints all the files and directories of a directory formatted according to their depth.
func printDir(dirPath string, depth *int, ignored []string) {
	files := getFiles(dirPath, ignored)

	if *depth == 0 {
		fmt.Println(path.Base(dirPath))
	}

	for _, file := range files {
		*depth++
		fullPrefix := strings.Repeat("  ", *depth-1) + PREFIX
		filepath := path.Join(dirPath, file.Name())
		fmt.Println(fullPrefix + file.Name())
		if file.IsDir() {
			printDir(filepath, depth, ignored)
		}
		*depth--
	}
}

// getFiles returns all files from the dir param, except for the ones listed in the ignored param.
func getFiles(dir string, ignored []string) []fs.DirEntry {
	var validFiles []fs.DirEntry
	files, err := os.ReadDir(dir)
	checkErr(err)

	for _, file := range files {
		if !checkIgnore(file.Name(), ignored) {
			validFiles = append(validFiles, file)
		}
	}

	return validFiles
}

// checkIgnore returns true if the fileName is in the ignored slice and false otherwise.
func checkIgnore(fileName string, ignored []string) bool {
	if len(ignored) == 0 {
		return false
	}

	for _, ignoredName := range ignored {
		if fileName == ignoredName {
			return true
		}
	}
	return false
}

// loadIgnore reads the ignore like file and returns a slice where each element is a line in the read file.
func loadIgnore(ignoreFileName string) []string {
	if ignoreFileName == "" {
		return make([]string, 0)
	}

	file, err := os.ReadFile(ignoreFileName)
	checkErr(err)

	fields := strings.Fields(string(file))
	fields = append(fields, ".yaspignore", "yasp.yml")

	return fields
}

// loadConfig reads the yml config file and loads the data into the Config struct
func loadConfig(config *Config) {
	file, err := os.ReadFile("yasp.yml")
	checkErr(err)

	err = yaml.Unmarshal(file, config)
	checkErr(err)
}

// TODO! Refactor error logging.
// checkErr logs all errors as fatals.
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
