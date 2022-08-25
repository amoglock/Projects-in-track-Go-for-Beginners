package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func inputCommand() []string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter a command and data:")
	scanner.Scan()
	command := strings.Split(scanner.Text(), " ")
	return command
}

func createNote(notes *[]string, command []string) {
	*notes = append(*notes, strings.Join(command, " "))
	fmt.Println("[OK] The note was successfully created")
}

func clearNotes() []string {
	fmt.Println("[OK] All notes were successfully deleted")
	return make([]string, 0)
}

func printNotes(notes *[]string) {
	for index, element := range *notes {
		if len(element) > 0 {
			fmt.Printf("[Info] %d: %s\n", index+1, element)
		}
	}
}

func checkPosition(command string, maxNotes int) int {
	val, err := strconv.Atoi(command)
	if err != nil {
		err = fmt.Errorf("[Error] Invalid position: %s\n", command)
		fmt.Print(err)
	}
	if val > maxNotes {
		err = fmt.Errorf("[Error] Position %d is out of the boundary [1, %d]\n", val, maxNotes)
		fmt.Print(err)
		val = 0
	}

	return val
}

func checkCommand(command []string, notes []string, maxNotes int) {
	if len(command) == 1 {
		fmt.Println("[Error] Missing position argument")
	} else if len(command) == 2 {
		fmt.Println("[Error] Missing note argument")
	} else {
		position := checkPosition(command[1], maxNotes)
		if position > 0 {
			checkNote(position, notes, command)
		}
	}
}

func checkNote(position int, notes, command []string) {
	if position > len(notes) {
		fmt.Print("[Error] There is nothing to update\n")
	} else {
		updateNote(command[2:], notes, position)
	}
}

func updateNote(command, notes []string, position int) {
	s := strings.Join(command, " ")
	notes[position-1] = s
	fmt.Printf("[OK] The note at position %d was successfully updated\n", position)
}

func deleteNote(command []string, notes []string, maxNotes int) []string {
	if len(command) == 1 {
		fmt.Println("[Error] Missing position argument")
	} else {
		position := checkPosition(command[1], maxNotes)
		if position > len(notes) {
			fmt.Print("[Error] There is nothing to delete\n")
		} else {
			if position > 0 {
				copy(notes[position-1:], notes[position:])
				notes = notes[:len(notes)-1]
				fmt.Printf("[OK] The note at position %d was successfully deleted\n", position)
			}
		}
	}
	return notes
}

func usersInput(notes []string, maxNotes int) {
	for {
		command := inputCommand()
		switch command[0] {
		case "exit":
			fmt.Println("[Info] Bye!")
			os.Exit(0)
		case "create":
			if len(command[1:]) < 1 {
				err := fmt.Errorf("[Error] Missing note argument")
				fmt.Println(err)
			} else if len(notes) < maxNotes {
				createNote(&notes, command[1:])
			} else {
				err := fmt.Errorf("[Error] Notepad is full")
				fmt.Println(err)
			}
		case "clear":
			notes = clearNotes()
		case "list":
			if len(notes) > 0 {
				printNotes(&notes)
			} else {
				fmt.Println("[Info] Notepad is empty")
			}
		case "update":
			checkCommand(command, notes, maxNotes)
		case "delete":
			notes = deleteNote(command, notes, maxNotes)
		default:
			err := fmt.Errorf("[Error] Unknown command")
			fmt.Println(err)
		}
	}
}

func main() {
	fmt.Println("Enter the maximum number of notes:")
	var maxNotes int
	fmt.Scan(&maxNotes)
	notes := make([]string, 0)
	usersInput(notes, maxNotes)
}
