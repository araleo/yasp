package main

import "strings"

// buildRegexes returns a map with the regex string for each command.
func buildRegexes(config *Config) map[string]string {
	return map[string]string{
		"print":  BuildRegexString(config.Print.Commands),
		"todos":  BuildRegexString(config.Todos.Commands),
		"issues": BuildRegexString(config.Issues.Commands),
	}
}

// BuildRegexString builds a string formatted and ready to compiled with the regexp package.
func BuildRegexString(allCommands string) string {
	commands := strings.Split(allCommands, ",")
	expr := ``

	for _, command := range commands {
		command = strings.ReplaceAll(command, "(", "\\(")
		command = strings.ReplaceAll(command, ")", "")
		expr += command + "|"
	}

	return strings.TrimSuffix(expr, "|")
}
