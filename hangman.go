package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	stages := ExtractGallows("default.gallows")

	fmt.Println(stages[5])
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
