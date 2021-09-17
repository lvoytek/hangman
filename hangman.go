package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	stages := ExtractGallows("default.gallows")

	fmt.Println(stages[5])

	currentWord := ChooseWord(ExtractDictionary("words.txt"))
	fmt.Println(currentWord)
}

// ExtractGallows returns a set of stages that the gallows and body can be in
// from the contents of a gallows text file. The file is formatted such that
// each stage is in sequential order separated by two empty lines containing
// only \n characters. The function takes in a file path input. If the file
// does not exist then an error is logged and nil is returned
func ExtractGallows(gallowsFile string) []string {
	fileBytes, err := ioutil.ReadFile(gallowsFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	fileString := string(fileBytes)
	fileString = strings.Trim(fileString, "\n")

	var gallowsStages []string

	for strings.Count(fileString, "\n\n\n") > 0 {
		gallowsStages = append(gallowsStages, fileString[0:strings.Index(fileString, "\n\n\n")])
		fileString = fileString[strings.Index(fileString, "\n\n\n")+3:]
	}

	gallowsStages = append(gallowsStages, fileString)

	return gallowsStages
}

// ChooseWord picks a random word from a string array and returns it.
// If there are no strings in the list, return an empty string
func ChooseWord(wordList []string) string {
	if len(wordList) == 0 {
		return ""
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	timeBasedRand := rand.New(s1)

	return wordList[timeBasedRand.Intn(len(wordList))]
}

// ExtractDictionary reads from a dictionary file and return a string
// array containing strings with only lowercase letters log an error
// and returns nil if the file is invalid
// A dictionary in this case is a text file containing words separated by
// '\n' new lines.
func ExtractDictionary(dictionaryFile string) []string {
	file, err := os.Open(dictionaryFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	dictScanner := bufio.NewScanner(file)
	dictScanner.Split(bufio.ScanLines)

	isLowercaseWord := regexp.MustCompile(`^[a-z]+$`).MatchString

	var outputWords []string
	var currentWord string

	for dictScanner.Scan() {
		currentWord = dictScanner.Text()

		if isLowercaseWord(currentWord) {
			outputWords = append(outputWords, currentWord)
		}
	}

	return outputWords
}
