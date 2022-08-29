package main

import (
	"fmt"
	"os"
)

//Ingredients and price for one cup espresso
const espressoWater uint = 250
const espressoMilk uint = 0
const espressoCoffeeBeans uint = 16
const espressoPrice uint = 4

//Ingredients and price for one cup latte
const latteWater uint = 350
const latteMilk uint = 75
const latteCoffeeBeans uint = 20
const lattePrice uint = 7

//Ingredients and price for one cup cappuccino
const cappuccinoWater uint = 200
const cappuccinoMilk uint = 100
const cappuccinoCoffeeBeans uint = 12
const cappuccinoPrice uint = 6

const disposableCup uint = 1 //one empty cup

var supplies = [5]uint{400, 540, 120, 9, 550} // Amount {water, milk, beans, cups, money} respectively

func actionOptions() {
	// There we "say" machine what we want. "Actions" on next line
	fmt.Print("Write action (buy, fill, take, remaining, exit):\n")
	var action string
	fmt.Scan(&action)
	switch action {
	case "buy": // by the coffee
		machineMenu()
	case "fill": // add some ingredients into machine
		fillTheMachine()
	case "take": // take all money
		takeMoney()
	case "remaining": // show remaining ingredients
		showSupplies()
	case "exit": // exit
		os.Exit(0)
	default:
		fmt.Print("Can`t run this action!")
		actionOptions()
	}
}

func machineMenu() {
	// User can choose favorite coffee
	fmt.Print("\nWhat do you want to buy? 1 - espresso, 2 - latte, 3 - cappuccino, back - to main menu\n")
	var choice string
	var coffeeChoice [5]uint
	fmt.Scan(&choice)
	switch choice {
	case "1":
		coffeeChoice = [5]uint{espressoWater, espressoMilk, espressoCoffeeBeans, disposableCup, espressoPrice}
	case "2":
		coffeeChoice = [5]uint{latteWater, latteMilk, latteCoffeeBeans, disposableCup, lattePrice}
	case "3":
		coffeeChoice = [5]uint{cappuccinoWater, cappuccinoMilk, cappuccinoCoffeeBeans, disposableCup, cappuccinoPrice}
	case "back":
		actionOptions()
	}
	checkStockBalance(coffeeChoice)
}

func showSupplies() {
	fmt.Printf("\nThe coffee machine has:\n%d ml of water\n%d ml of milk\n"+
		"%d g of coffee beans\n%d disposable cups\n$%d of money\n\n",
		supplies[0], supplies[1], supplies[2], supplies[3], supplies[4])
	actionOptions()
}

func checkStockBalance(coffeeCup [5]uint) {
	// check if there are enough ingredients in the machine
	resources := [4]string{"water", "milk", "coffee beans", "disposable cups"}
	for index, element := range coffeeCup[:4] {
		if element > supplies[index] {
			fmt.Printf("Sorry, not enough %s!\n", resources[index])
			actionOptions()
		}
	}
	fmt.Print("I have enough resources, making you a coffee!\n")
	procedure(coffeeCup)
}

func procedure(coffeeCup [5]uint) {
	// stocks decrease, money increases
	for index := range supplies {
		if index == 4 {
			supplies[index] += coffeeCup[index]
		} else {
			supplies[index] -= coffeeCup[index]
		}
	}
	actionOptions()
}

func fillTheMachine() {
	// add some ingredients
	names := []string{"ml of water", "ml of milk", "grams of coffee beans", "disposable cups of coffee"}
	for index, name := range names {
		fmt.Printf("Write how many %s you want to add:\n", name)
		var amount uint
		fmt.Scan(&amount)
		supplies[index] += amount
	}
	actionOptions()
}

func takeMoney() {
	fmt.Printf("I gave you $%d\n", supplies[4])
	supplies[4] = 0
	actionOptions()
}

func main() {
	actionOptions()
}
