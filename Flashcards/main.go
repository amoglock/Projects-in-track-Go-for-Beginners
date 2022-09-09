package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Card struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
	Mistakes   int    `json:"mistakes"`
}

var cardsMap = make(map[string]Card)
var logList []string
var tag []string

func main() {
	importFile := flag.String("import_from", "", "Enter file from import.")
	exportFile := flag.String("export_to", "", "Enter file to export.")
	flag.Parse()

	if *exportFile != "" {
		tag = append(tag, *exportFile)
	}
	if *importFile != "" {
		importCards(*importFile, "tag")
	}

	menu()
}
func menu() {

	appendToLogList("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):\n")
	action := input()
	switch action {
	case "add":
		add()
	case "remove":
		deleteCard()
	case "import":
		appendToLogList("File name")
		fileName := input()
		importCards(fileName, "")
	case "export":
		appendToLogList("File name:")
		fileName := input()
		exportCards(fileName, "")
	case "ask":
		checkCardsAmount()
	case "log":
		saveLog()
	case "hardest card":
		hardestCard()
	case "reset stats":
		resetStats(&cardsMap)
	case "exit":
		if len(tag) == 1 {
			exportCards(tag[0], "exit")
		}
		fmt.Print("Bye bye!")
		os.Exit(0)
	default:
		menu()
	}

}

// All lines input (terms or definitions)
func input() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	logList = append(logList, line)
	return line
}

func add() {
	defer menu()
	var term, definition string
	appendToLogList("The card:\n")
	term = input()
	checkTerm(&term)
	appendToLogList("The definition:\n")
	definition = input()
	checkDefinition(&definition)
	cardsMap[term] = Card{term, definition, 0}
	appendToLogList(fmt.Sprintf("The pair (\"%s\":\"%s\") has been added.\n\n", term, definition))
}

func deleteCard() {
	defer menu()
	appendToLogList("Which card?")
	cardForDelete := input()
	if _, ok := cardsMap[cardForDelete]; ok {
		delete(cardsMap, cardForDelete)
		appendToLogList("The card has been removed.\n\n")
		return
	}
	appendToLogList(fmt.Sprintf("Can't remove \"%s\": there is no such card.\n\n", cardForDelete))
}

func importCards(fileName string, command string) {
	file, err := os.Open(fileName)
	if err != nil {
		appendToLogList("File not found\n\n")
		menu()
	}
	defer file.Close()
	a := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		var line Card
		err := json.Unmarshal([]byte(data), &line)
		if err != nil {
			log.Fatal(err)
		}
		cardsMap[line.Term] = line
		a++
	}
	appendToLogList(fmt.Sprintf("%d cards have been loaded.\n\n", a))
	if command == "tag" {
		return
	}
	menu()
}

func exportCards(fileName string, command string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, val := range cardsMap {
		objJson, err := json.Marshal(val)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(file, string(objJson))
	}
	appendToLogList(fmt.Sprintf("%d cards have been saved.\n\n", len(cardsMap)))
	if command == "exit" {
		fmt.Print("Bye bye!")
		os.Exit(0)
	}
	menu()
}

func checkTerm(term *string) {
	if _, ok := cardsMap[*term]; ok {
		appendToLogList(fmt.Sprintf("The term \"%s\" already exists. Try again:\n", *term))
		*term = input()
		checkTerm(term)
	}
}

func checkDefinition(definition *string) {
	for _, val := range cardsMap {
		if *definition == val.Definition {
			appendToLogList(fmt.Sprintf("The definition \"%s\" already exists. Try again:\n", val.Definition))
			*definition = input()
			checkDefinition(definition)
		}
	}
}

func game(times int) {
	i := 0
	for i < times {
		for _, val := range cardsMap {
			if i == times {
				menu()
			}
			appendToLogList(fmt.Sprintf("Print the definition of \"%s\":\n", val.Term))
			answer := input()
			switch answer {
			case val.Definition:
				appendToLogList("Correct!\n")
			default:
				add := checkAnotherTerm(answer)
				appendToLogList(fmt.Sprintf("Wrong. The right answer is \"%s\"%s\n", val.Definition, add))
				cardsMap[val.Term] = Card{val.Term, val.Definition, val.Mistakes + 1}
			}
			i++
		}
	}
}

func checkAnotherTerm(a string) string {
	for _, val := range cardsMap {
		if a == val.Definition {
			return fmt.Sprintf(", but your definition is correct for \"%s\".", val.Term)
		}
	}
	return "."
}

func checkCardsAmount() {
	var times int
	appendToLogList("How many times to ask?\n")
	fmt.Scan(&times)
	if len(cardsMap) == 0 {
		appendToLogList("I don't have that many cards. Am is empty!\n")
	} else {
		game(times)
	}
	menu()
}

func saveLog() {
	appendToLogList("File name\n")
	fileName := input()
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, val := range logList {
		_, err := fmt.Fprintln(file, val) // append the additional line
		if err != nil {
			log.Fatal(err)
		}
	}
	appendToLogList("The log has been saved.\n")
	menu()
}

func appendToLogList(line string) {
	logList = append(logList, strings.TrimSpace(line))
	fmt.Print(line)
}

func hardestCard() {
	var mistakes []string
	max := 1
	for _, val := range cardsMap {
		if val.Mistakes == max {
			mistakes = append(mistakes, val.Term)
		} else if val.Mistakes > max {
			max = val.Mistakes
			mistakes = nil
			mistakes = append(mistakes, val.Term)
		}
	}
	if len(mistakes) == 0 {
		appendToLogList("There are no cards with errors.\n\n")
	} else if len(mistakes) == 1 {
		appendToLogList(fmt.Sprintf("The hardest card is \"%s\". You have %d errors answering it.\n\n", strings.Join(mistakes, "\", \""), max))
	} else {
		appendToLogList(fmt.Sprintf("The hardest cards are \"%s\". You have %d errors answering it.\n\n", strings.Join(mistakes, "\", \""), max))
	}
	menu()
}

func resetStats(cards *map[string]Card) {
	for _, val := range *cards {
		cardsMap[val.Term] = Card{val.Term, val.Definition, 0}
	}
	appendToLogList("Card statistics have been reset.\n")
	menu()
}
