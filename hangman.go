/*
 * Hangman
 * Copyright (C) 2021  Lena Voytek
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	// Show help if requested
	if sendHelp(os.Args) {
		return
	}

	// Get the gallows and dictionary for this game
	stages := ExtractGallows(getGallowsFromArgs(os.Args))
	dict := ExtractDictionary(getDictionaryFromArgs(os.Args))

	if stages == nil || dict == nil {
		return
	}

	currentWord := ChooseWord(dict)

	var guessedLetters []string
	var nextGuess string

	isLowercaseWord := regexp.MustCompile(`^[a-z]+$`).MatchString

	// Play the game
	for !PrintStage(currentWord, guessedLetters, stages) {
		fmt.Print("Guess> ")
		fmt.Scanln(&nextGuess)

		// Add letters in nextGuess to guessed letters
		if isLowercaseWord(nextGuess) {
			for _, c := range nextGuess {
				letterAlreadyGuessed := false

				for _, guess := range guessedLetters {
					if string(c) == guess {
						letterAlreadyGuessed = true
						break
					}
				}

				if !letterAlreadyGuessed {
					guessedLetters = append(guessedLetters, string(c))
				}
			}
		}
	}
}

// PrintStage displays the current state of the game in the terminal. This
// includes the current body parts, the gallows, the word outline, and all
// guessed letters. The function returns true if the game has finished, and
// false otherwise
func PrintStage(currentWord string, guessedLetters []string, gallows []string) bool {
	// Get all the bad letter guesses
	var badGuesses []string
	for _, c := range guessedLetters {
		if !strings.Contains(currentWord, c) {
			badGuesses = append(badGuesses, c)
		}
	}

	// Split current gallows string into individual lines
	stageLines := strings.Split(gallows[min(len(badGuesses), len(gallows)-1)], "\n")

	// Get the longest line of the stage to offset the bad guess box from
	badGuessOffset := 0

	for _, stageLine := range stageLines {
		if len(stageLine) > badGuessOffset {
			badGuessOffset = len(stageLine)
		}
	}

	badGuessOffset += 3

	for i := 0; i < 10; i++ {
		fmt.Println()
	}

	if len(stageLines) >= 6 {
		startBadGuessesLine := (len(stageLines) - 6) / 2

		for i, stageLine := range stageLines {
			if i >= startBadGuessesLine && i < startBadGuessesLine+6 {
				fmt.Print(stageLine)
				fmt.Print(strings.Repeat(" ", (badGuessOffset - len(stageLine))))

				if i == startBadGuessesLine {
					fmt.Println(" ___________________")
				} else if i == startBadGuessesLine+5 {
					fmt.Println("|___________________|")
				} else {
					fmt.Print("| ")

					startIndex := (i - startBadGuessesLine - 1) * 9

					for i := startIndex; i < startIndex+9; i++ {
						if i < len(badGuesses) {
							fmt.Print(badGuesses[i] + " ")
						} else {
							fmt.Print("  ")
						}
					}

					fmt.Println("|")
				}

			} else {
				fmt.Println(stageLine)
			}
		}

	} else {
		//TODO Make guess printing work for smaller gallows
	}

	// Print guessed letters and the word outline
	fmt.Println()

	wordComplete := true

	for _, c := range currentWord {
		letterGuessed := false

		for _, guess := range guessedLetters {
			if string(c) == guess {
				letterGuessed = true
				break
			}
		}

		if letterGuessed {
			fmt.Print(string(c))
		} else {
			wordComplete = false
			fmt.Print(" ")
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", len(currentWord)))

	// Print win or loss screen if game over
	if len(badGuesses) >= len(gallows)-1 {
		PrintLoseScreen(currentWord)
		return true
	}

	if wordComplete {
		PrintWinScreen()
		return true
	}

	return false
}

func PrintWinScreen() {
	fmt.Println("__  __               _       ___")
	fmt.Println("\\ \\/ /___  __  __   | |     / (_)___")
	fmt.Println(" \\  / __ \\/ / / /   | | /| / / / __ \\")
	fmt.Println(" / / /_/ / /_/ /    | |/ |/ / / / / /")
	fmt.Println("/_/\\____/\\__,_/     |__/|__/_/_/ /_/")
}

func PrintLoseScreen(word string) {
	fmt.Println("__  __               __")
	fmt.Println("\\ \\/ /___  __  __   / /   ____  ________")
	fmt.Println(" \\  / __ \\/ / / /  / /   / __ \\/ ___/ _ \\")
	fmt.Println(" / / /_/ / /_/ /  / /___/ /_/ (__  )  __/")
	fmt.Println("/_/\\____/\\__,_/  /_____/\\____/____/\\___/")
	fmt.Println()
	fmt.Println("The word was: " + word)
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

// min gets the minimum value between two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getGallowsFromArgs(args []string) string {
	checkNextItem := false
	gallowsFile := "default.gallows"

	for _, arg := range args {
		if checkNextItem {
			gallowsFile = arg
		} else if arg == "--gallows" || arg == "-g" {
			checkNextItem = true
		}
	}

	return gallowsFile
}

func getDictionaryFromArgs(args []string) string {
	checkNextItem := false
	dictFile := "words.txt"

	for _, arg := range args {
		if checkNextItem {
			dictFile = arg
		} else if arg == "--dictionary" || arg == "-d" {
			checkNextItem = true
		}
	}

	return dictFile
}

// sendHelp displays the help section if requested in the arguments.
// It returns true if the help screen was shown
func sendHelp(args []string) bool {
	helpNeeded := false

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			helpNeeded = true
			break
		}
	}

	if helpNeeded {
		fmt.Println("Options:")
		fmt.Println("\t-d, --dictionary [filename]\tProvide a custom dictionary file, which is a set of lowercase words split up by new lines")
		fmt.Println("\t-g, --gallows [filename]\tProvide a custom gallows and body progression design, which is a set of stages separated by two new lines")
		fmt.Println("\t-h, --help\t\t\tPrint this screen")
	}

	return helpNeeded
}
