package parse

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), "=")
		name, val := line[0], line[1]
		os.Setenv(name, val)
	}

}
