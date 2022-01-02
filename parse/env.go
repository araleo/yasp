package parse

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type EnvConfig struct {
	Path string
	File string
	Vars string
}

func LoadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), "=")
		if len(line) == 2 {
			name, val := line[0], line[1]
			os.Setenv(name, val)
		}
	}
}

// CheckEnv receives the data for an env config file and returns true if it finds a file like it in the project or false otherwise.
func CheckEnv(env EnvConfig) bool {
	expectedVars := strings.Split(env.Vars, ",")

	filepath := path.Join(env.Path, env.File)
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Can't find the env file at %s\n", filepath)
		return false
	}

	lines := strings.Split(string(file), "\n")
	foundVars := make([]string, 0)
	for _, line := range lines {
		keyVal := strings.Split(line, "=")
		if len(keyVal) != 2 {
			fmt.Printf("Invalid value for variable %s in file %s\n", line, filepath)
			return false
		}
		k := keyVal[0]
		foundVars = append(foundVars, k)
	}

	foundStr := strings.Join(foundVars, ",")
	for _, expected := range expectedVars {
		if !strings.Contains(foundStr, expected) {
			fmt.Printf("Can't find the variable %s in the %s file\n", expected, filepath)
			return false
		}
	}

	return true
}
