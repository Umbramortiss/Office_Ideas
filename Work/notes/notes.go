package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
)

var (
	Green = color.New(color.FgGreen).Add(color.Bold)
	let   string
)

const (
    letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
    numBytes = "0123456789"
)


func main() {

	toolControl()

	

	/* notes := phonetics_words()
	fmt.Println("Notes:")
	for _, note := range notes {
		fmt.Println(note)
	}

}

func notes_list() []string {

	notesDir := "Work_Notes"
	nList, err := filepath.Glob(filepath.Join(notesDir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	var notes []string
	for _, file := range nList {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, string(content))
	}
	return notes
}

func listTitle() []string {

	notesDir := "Work_Notes"
	nList, err := filepath.Glob(filepath.Join(notesDir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	var titles []string
	for _, file := range nList {
		titles = append(titles, file)
	}
	return titles
}

func chooseNote(notes []string, titles []string) {
	fmt.Println("Available Notes: ")
	for i, title := range titles {
		fmt.Printf("%d. %s\n", i+1, extractFileName(title))
	}
	selectedNoteIndex := userInput("Enter the number of the note: ")
	num, err := strconv.Atoi(selectedNoteIndex)
	if err != nil || num < 1 || num > len(titles) {
		fmt.Println("Invalid input. Please enter a valid note number.")
		return
	}

	selectedNote := notes[num-1]
	fmt.Println("You chose note number: ", num)
	fmt.Println("Note Content: ")
	fmt.Println(selectedNote)

	copyChoice := userInput("Copy this note to clipboard? (y/n): ")
	if strings.ToLower(copyChoice) == "y" {
		err := clipboard.WriteAll(selectedNote)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Note copied to clipboard!")
	} else {
		fmt.Println("Note not copied to clipboard.")
	}

}

func userInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func extractFileName(notes string) string {
	baseName := filepath.Base(notes)

	baseNameWithoutExtension := strings.TrimSuffix(baseName, ".txt")

	parts := strings.Split(baseNameWithoutExtension, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return baseNameWithoutExtension
}

func generatePassword(length int, includeNumber bool, includeSpecial bool) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var password []byte
	var charSource string

	if includeNumber {
		charSource += "0123456789"
	}
	if includeSpecial {
		charSource += "!@#$%^&*()_+=-"
	}
	charSource += charset

	for i := 0; i < length; i++ {
		randNum := rand.Intn(len(charSource))
		password = append(password, charSource[randNum])
	}
	return string(password)
}

func toolControl(){

	Green.Println("Choose your Work option: ")
	fmt.Scanln(&option)

	switch option {
	case 1:
		notes := notes_list()
		chooseNote(notes)
	
	case 2:
		password := generatePassword(12, true, true)
    	fmt.Println("Random password:", password)
	
	default:
		Green.Println("Invalid option")
	}

}
