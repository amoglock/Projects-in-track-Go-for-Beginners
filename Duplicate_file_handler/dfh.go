package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	Size int64
	Path []string
}

var files = make(map[int64][]string) //files with same size [size]:paths
var allFiles []File

func main() {
	enterTheFolder()
	traverse()
	sortingOptions()
}

//enterTheFolder checks whether the directory name is passed to the command line
func enterTheFolder() {
	if len(os.Args) != 2 {
		fmt.Println("Directory is not specified")
		os.Exit(0)
	}
}

//traverse over all files in directory
func traverse() {
	var extension string
	fmt.Println("Enter file format:")
	fmt.Scanln(&extension)

	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() { //if is not directory
			if extension == filepath.Ext(path) || extension == "" { //find what need extension || all files
				fileInfo, err := os.Stat(path)
				if err != nil {
					log.Fatal(err)
				}
				files[fileInfo.Size()] = append(files[fileInfo.Size()], filepath.Join(filepath.Dir(path), fileInfo.Name()))
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for i, j := range files { //turn filtered files map into struct
		allFiles = append(allFiles, File{Size: i, Path: j})
	}
}

func sortingOptions() {
	fmt.Print("Size sorting options:\n1. Descending\n2. Ascending\n\n")
	action()
}

func action() {
	var option int
	fmt.Println("Enter a sorting option:")
	_, err := fmt.Scan(&option)
	if err != nil {
		log.Fatal(err)
	}
	switch option {
	case 1:
		sortingDuplicates(option)
	case 2:
		sortingDuplicates(option)
	default:
		fmt.Println("Wrong option")
		action()
	}
}

//sortingDuplicates sorts files by size depending on the option
//finds files with the same size and put them in a separate list. Then we work with him
func sortingDuplicates(option int) {
	var sameSizesFiles []File
	if option == 1 {
		sort.Slice(allFiles, func(i, j int) bool {
			return allFiles[i].Size > allFiles[j].Size
		})
	} else if option == 2 {
		sort.Slice(allFiles, func(i, j int) bool {
			return allFiles[i].Size < allFiles[j].Size
		})
	}
	for _, path := range allFiles {
		if len(path.Path) > 1 {
			sameSizesFiles = append(sameSizesFiles, path)
			fmt.Printf("%d bytes\n", path.Size)
			for _, i := range path.Path {
				fmt.Println(i)
			}
			fmt.Println()
		}
	}
	checkDuplicatesOptions(&sameSizesFiles)
}

func checkDuplicatesOptions(listOfFiles *[]File) {
	var option string
	fmt.Println("Check for duplicates?")
	fmt.Scan(&option)
	switch option {
	case "yes":
		fmt.Println()
		calculateHash(listOfFiles)
	case "no":
		os.Exit(0)
	default:
		checkDuplicatesOptions(listOfFiles)
	}
}

//Find the hash of each file. Sort the files with the same hash and put them in forDelete map
//Assign each file a number in order.
func calculateHash(listOfFiles *[]File) {
	count := 1
	var forDelete = make(map[int]string)
	for _, path := range *listOfFiles {
		fmt.Printf("%d bytes\n", path.Size)
		var hashFiles = make(map[string][]string)
		for _, filePath := range path.Path {
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}
			md5Hash := md5.New()
			if _, err := io.Copy(md5Hash, file); err != nil {
				log.Fatal(err)
			}
			file.Close()
			hash := fmt.Sprintf("%x\n", md5Hash.Sum(nil))
			hashFiles[hash] = append(hashFiles[hash], filePath)
		}
		for k, v := range hashFiles {
			if len(v) > 1 {
				fmt.Printf("Hash: %s", k)
				for _, i := range v {
					fmt.Printf("%d. %s\n", count, i)
					forDelete[count] = i
					count++
				}
				fmt.Println()
			}
		}
	}
	deleteAction(forDelete)
}

func deleteAction(forDelete map[int]string) {
	var action string
	fmt.Println("Delete files?")
	fmt.Scanln(&action)
	switch action {
	case "yes":
		numbers := inputNumbers(len(forDelete))

		deleting(forDelete, numbers...)
	case "no":
		os.Exit(0)
	default:
		fmt.Println("Wrong option")
		deleteAction(forDelete)
	}
}

//Input numbers files to delete
func inputNumbers(count int) []int {
	fmt.Println("Enter file numbers to delete:")
	scanner := bufio.NewScanner(os.Stdin)
	var numbers []int
	for scanner.Scan() {
		line := scanner.Text()
		lineArray := strings.Split(line, " ")
		for _, number := range lineArray {
			val, err := strconv.Atoi(number)
			if err != nil {
				fmt.Println("Wrong format")
				inputNumbers(count)
			}
			numbers = append(numbers, val)
		}
		break
	}
	return numbers
}

//Delete files by number. Calculate freed up space
func deleting(forDelete map[int]string, numbers ...int) {
	allBytes := 0
	for _, number := range numbers {
		fileInfo, err := os.Stat(forDelete[number])
		if err != nil {
			continue
		}
		allBytes += int(fileInfo.Size())
		err = os.Remove(forDelete[number])
		if err != nil {
			continue
		}
	}
	fmt.Printf("Total freed up space: %d bytes\n", allBytes)
}
