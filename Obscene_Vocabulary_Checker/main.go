package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func openFile() map[string]struct{} {
	tabooWords := make(map[string]struct{})
	var fileName string
	fmt.Scan(&fileName)
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range strings.Split(string(data), "\n") {
		tabooWords[strings.ToLower(name)] = struct{}{}
	}
	return tabooWords
}

// Input and return a line to check
func checkInputLine() string {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	line = strings.Trim(line, ".")
	return line
}

// Outputs a new line with the replaced words if they were in the forbidden list
func printLine(tabooWords map[string]struct{}) {
	line := checkInputLine()
	if line != "exit" {
		splitLine := strings.Split(line, " ")
		for _, value := range splitLine {
			line = strings.Replace(line, value, checkWord(value, tabooWords), 1)
		}
		fmt.Println(line)
		printLine(tabooWords)
	}
	fmt.Print("Bye!\n")
}

// Check word for ban
func checkWord(value string, tabooWords map[string]struct{}) string {
	if _, ok := tabooWords[strings.ToLower(value)]; ok {
		// Replace each letter with a "*" symbol
		maskedValue := strings.Repeat("*", len(value))
		return maskedValue
	}
	return value
}

func main() {
	tabooWords := openFile()
	printLine(tabooWords)
}
