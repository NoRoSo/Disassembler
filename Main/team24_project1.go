package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var InputFileName *string
var OutputFileName *string
var ProgramCounter int = 0

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

	WriteOutput(fileLines)

	err = readFile.Close()
	if err != nil {
		return
	}
}

func WriteOutput(fileLine []string) {
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

	for i := 0; i < len(fileLine); i++ {
		tempString := CreateString(fileLine[i])
		_, errs = file.WriteString(tempString)
		if errs != nil {
			fmt.Println("Failed to write to file:", errs)
			return
		}
		ProgramCounter += 4
	}

	fmt.Println("Wrote to file: ", *OutputFileName)
}

func CreateString(instructionCode string) string {
	opcode := instructionCode[0:11] //the opcode of an instruction, will change depending on instruction

	if strings.Index(opcode, "000101") == 0 { //B instruction
		opcode = instructionCode[0:6] //an example of opcode getting changed
		number := instructionCode[6:] //the bounds of where the number begins

		var decimalNum int64
		if number[0] == '0' {
			decimalNum, _ = strconv.ParseInt(number, 2, 28) // the _ is "error"
		} else {
			decimalNum, _ = strconv.ParseInt("100000000000000000000000000", 2, 28)
			temp, _ := strconv.ParseInt(number, 2, 28)
			decimalNum -= temp
			decimalNum *= -1
			fmt.Println(decimalNum)
		}
		formattedString := fmt.Sprintf("%-38s", opcode+" "+number)
		//38 characters
		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + "B   #" + strconv.FormatInt(decimalNum, 10) + "\n"
	}

	if strings.Index(opcode, "10001010000") == 0 { //AND instruction

	}

	if strings.Index(opcode, "10001011000") == 0 { //ADD instruction

	}

	if strings.Index(opcode, "1001000100") == 0 { //ADDI instruction

	}

	if strings.Index(opcode, "10101010000") == 0 { //ORR instruction

	}

	if strings.Index(opcode, "10110100") == 0 { //CBZ instruction

	}

	if strings.Index(opcode, "10110101") == 0 { //CBNZ instruction

	}

	if strings.Index(opcode, "11001011000") == 0 { //SUB instruction

	}

	if strings.Index(opcode, "1101000100") == 0 { //SUBI instruction

	}

	if strings.Index(opcode, "110100101") == 0 { //MOVZ instruction

	}

	if strings.Index(opcode, "111100101") == 0 { //MOVK instruction

	}

	if strings.Index(opcode, "11010011010") == 0 { //LSR instruction

	}

	if strings.Index(opcode, "11010011011") == 0 { //LSL instruction

	}

	if strings.Index(opcode, "11111000000") == 0 { //STUR instruction

	}

	if strings.Index(opcode, "11111000010") == 0 { //LDUR instruction

	}

	if strings.Index(opcode, "11010011100") == 0 { //ASR instruction

	}

	if strings.Index(opcode, "00000000") == 0 { //NOP instruction

	}

	if strings.Index(opcode, "11101010000") == 0 { //EOR instruction

	}

	return "\n"
}
