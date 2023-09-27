package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var InputFileName *string
var OutputFileName *string

func main() {

	//retrieving flags -i and -o (input and output file name)
	InputFileName = flag.String("i", "", "Gets the input file name")
	OutputFileName = flag.String("o", "", "Gets the output file name")

	flag.Parse()

	readFile, err := os.Open(*InputFileName)

	if err != nil { //in case it doesn't open
		fmt.Println(err)
		return
	}

	//reading the input of the text file line by line and inserting it into a list.
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	fmt.Println(fileLines)

	err = readFile.Close()
	if err != nil {
		return
	}
}

func WriteOutput() {
	file, errs := os.Create(*OutputFileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, errs = file.WriteString("Hello, world")
	if errs != nil {
		fmt.Println("Failed to write to file:", errs)
		return
	}
	fmt.Println("Wrote to file: ", *OutputFileName)
}
