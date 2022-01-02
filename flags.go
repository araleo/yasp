package main

import (
	"errors"
	"flag"
	"os"
	"path"
)

var HELPERS = map[string]string{
	"dir":     "Please enter the absolute path of a directory or a dot to use the cwd.",
	"ignore":  "Please enter the name of the ignore like file or a dot to use the default.",
	"command": "Please enter a command to run: ls, print, todo, env, diag, issues, snitch",
}

// ParseFlags declares and parses the flags accepted by the program
// and returns a map in which the keys are flags' names and values pointers to their content.
func ParseFlags() map[string]*string {
	dirFlag := flag.String("d", "", HELPERS["dir"])
	ignoreFlag := flag.String("i", "", HELPERS["ignore"])
	commandFlag := flag.String("c", "", HELPERS["command"])
	flag.Parse()

	return map[string]*string{
		"dir":     dirFlag,
		"ignore":  ignoreFlag,
		"command": commandFlag,
	}
}

// ValidateCommandFlag validates the -c flag, in which the user must provide
// the desired operation for the program to run.
// ls = list a directory, print = lists prints commands in code, todo = lists todos in code
func ValidateCommandFlag(commandFlag *string) (string, error) {
	if *commandFlag == "" {
		return "", errors.New(HELPERS["command"])
	}

	return *commandFlag, nil
}

// ValidateDirFlag validates the -d flag, in which the user must provide
// the absolute path to a directory or a dot to represent the cwd.
func ValidateDirFlag(dirFlag *string) (string, error) {
	if *dirFlag == "" || *dirFlag == "." {
		return os.Getwd()
	}

	return *dirFlag, nil
}

// ValidateIgnoreFlag validates the -i flag, in which the user must provide
// the absolute path to a ignore like file or a dot to use the default .gitnore.
func ValidateIgnoreFlag(ignoreFlag *string, config *Config) string {
	if *ignoreFlag == "" || *ignoreFlag == "." {
		*ignoreFlag = path.Join(config.Ignore.Path, config.Ignore.File)
	}

	return *ignoreFlag
}
