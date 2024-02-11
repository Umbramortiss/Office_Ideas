package main


import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/fatih/color"
)

func main() {
	for {
		words := phonetics_words()
		relatedWords(words)
		firstLetter(words)
		if let == "Done" {
			break
		}
	}

}

var (
	Green = color.New(color.FgGreen).Add(color.Bold)
	let   string
	//go:embed words.txt
	content string
)

func letter() {

	fmt.Print("Enter your desired letter: ")
	fmt.Scanln(&let)
	if let == "`,\",<, >,&" {
		log.Fatal("invalid charcters")
	}

}

func phonetics_words() []string {

	wordList, err := os.Open(content)

	if err != nil {
		log.Fatal(err)
	}

	defer wordList.Close()

	scanner := bufio.NewScanner(wordList)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return words
}

func relatedWords(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)

		}

	}

	return result
}

func firstLetter(s []string) {
	letter()
	var letters []string
	for _, str := range s {
		if string(str[0]) == let {
			letters = append(letters, str)
		}
	}
	if len(letters) > 0 {
		fmt.Println(let + " is for " + letters[rand.Intn(len(letters))])
		Green.Println("Type Done when finished ")
	} else {
		if let != "Done" {
			fmt.Println("No words found starting with letter " + let)
		}

	}

}

/*func randomWord(words []string) string {
	rand.Seed(time.Now().UnixNano())
	return words[rand.Intn(len(words))]
} */
