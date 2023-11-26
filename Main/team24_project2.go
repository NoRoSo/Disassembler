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
var ProgramCounter = 96
var hasBreak = false
var cycleNum = 1

var instructionMap = map[int64]Instruction{}
var registers [1024]int64
var maxRegisters = 32

type Instruction struct {
	instructionName string //The instruction name
	storeRegister   int64  //The register where the output gets stored
	targetRegister1 int64  //The first target register (if you read it, it will be the middle register)
	targetRegister3 int64  //The first target register (if you read it, it will be the middle register)
	immediateValue  int64  //If needed on certain instructions, use this to store immediate values (e.g: B, ADDI, SUBI, etc.)
	programCounter  int64  //The PC of this instruction.
}

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

	//initializing the registers to zeros
	for i := 0; i < 1024; i++ {
		registers[i] = 0
	}

	WriteOutput(fileLines)

	err = readFile.Close()
	if err != nil {
		return
	}
}

func WriteOutput(fileLine []string) {
	file, errs := os.Create(*OutputFileName + "_dis.txt")
	file2, errs := os.Create(*OutputFileName + "_sim.txt")

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

	currentInstruction := instructionMap[int64(96)]
	ProgramCounter = 96
	//for currentInstruction.instructionName != "BREAK" {
	for ProgramCounter <= 100 {
		tempString2 := createSimString(currentInstruction)
		_, errs = file2.WriteString(tempString2)
		cycleNum++
		currentInstruction = instructionMap[int64(ProgramCounter)]
	}

	err := file.Close()
	if err != nil {
		return
	}

	err = file2.Close()
	if err != nil {
		return
	}

}

func createSimString(instruction Instruction) string {
	simOutput := "=====================\n"
	simOutput += fmt.Sprintf("cycle:%d\t%d\t", cycleNum, instruction.programCounter)

	switch instruction.instructionName {
	case "ADD":
		//ADD code goes here COMPLETE
		simOutput += fmt.Sprintf("ADD\tR%d,\tR%d,\tR%d\n", instruction.storeRegister, instruction.targetRegister1, instruction.targetRegister3)
		registers[instruction.storeRegister] = registers[instruction.targetRegister1] + registers[instruction.targetRegister3]
		ProgramCounter += 4
	case "ADDI":
		//ADDI goes here COMPLETE
		simOutput += fmt.Sprintf("ADDI\tR%d,\tR%d,\t#%d\n", instruction.storeRegister, instruction.targetRegister1, instruction.immediateValue)
		registers[instruction.storeRegister] = registers[instruction.targetRegister1] + instruction.immediateValue
		ProgramCounter += 4
	case "SUB":
		//SUB goes here
	case "SUBI":
		//SUBI goes here
	case "AND":
		//AND goes here
	case "ORR":
		//ORR goes here
	case "EOR":
		//EOR goes here
	case "B":
		//B goes here

	case "CBZ":
		//CBZ goes here
	case "CBNZ":
		//CBNZ goes here
	case "MOVZ":
		//MOVZ goes here
	case "MOVK":
		//MOVK goes here
	case "LSR":
		//LSR goes here
	case "LSL":
		//LSL goes here
	case "STUR":
		//STUR goes here
	case "LDUR":
		//LDUR goes here
	case "ASR":
		//ASR goes here
	case "NOP":
		//NOP goes here COMPLETE
		simOutput += "NOP\n"
	case "BREAK":
		//BREAK goes here COMPLETE
		simOutput += "BREAK\n"
	}

	//required, to output registers and data after completing the actual processing of instructions
	simOutput += outputRegisters()
	simOutput += outputData()
	return simOutput
}

func CreateString(instructionCode string) string {
	opcode := instructionCode[0:11] //the opcode of an instruction, it will change depending on instruction

	if hasBreak {
		constNum := ConvertBinaryToDecimal(instructionCode)
		instFormat := fmt.Sprintf("%s\t", instructionCode)
		return instFormat + fmt.Sprintf("%s", strconv.FormatInt(int64(ProgramCounter), 10)) + "\t" + fmt.Sprintf("%d\n", constNum)
	}

	if strings.Index(opcode, "000101") == 0 { //B instruction
		opcode = instructionCode[0:6] //an example of opcode getting changed
		number := instructionCode[6:] //the bounds of where the number begins

		decimalNum := ConvertBinaryToDecimal(number)

		formattedString := fmt.Sprintf("%s\t", opcode+" "+number)
		instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), instructionName: "B", immediateValue: decimalNum}
		//38 characters
		return formattedString + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) + "B\t#" + strconv.FormatInt(decimalNum, 10) + "\n"
	}

	if strings.Index(opcode, "10001010000") == 0 { //AND instruction NOE
		return R_Instruction(instructionCode, "AND")
	}

	if strings.Index(opcode, "10001011000") == 0 { //ADD instruction NATHAN

		return R_Instruction(instructionCode, "ADD")
	}

	if strings.Index(opcode, "1001000100") == 0 { //ADDI instruction NATHAN
		return I_Instruction(instructionCode, "ADDI")
	}

	if strings.Index(opcode, "10101010000") == 0 { //ORR instruction NATHAN
		return R_Instruction(instructionCode, "ORR")
	}

	if strings.Index(opcode, "10110100") == 0 { //CBZ instruction
		return BranchInstruction(instructionCode, "CBZ")
	}

	if strings.Index(opcode, "10110101") == 0 { //CBNZ instruction
		return BranchInstruction(instructionCode, "CBNZ")
	}

	if strings.Index(opcode, "11001011000") == 0 { //SUB instruction NATHAN
		return R_Instruction(instructionCode, "SUB")
	}

	if strings.Index(opcode, "1101000100") == 0 { //SUBI instruction NATHAN
		return I_Instruction(instructionCode, "SUBI")
	}

	if strings.Index(opcode, "110100101") == 0 { //MOVZ instruction MAITLAND
		return IM_Instruction(instructionCode, "MOVZ")
	}

	if strings.Index(opcode, "111100101") == 0 { //MOVK instruction MAITLAND
		return IM_Instruction(instructionCode, "MOVK")
	}

	if strings.Index(opcode, "11010011010") == 0 { //LSR instruction MAITLAND
		return R_Instruction_Modified(instructionCode, "LSR")
	}

	if strings.Index(opcode, "11010011011") == 0 { //LSL instruction MAITLAND
		return R_Instruction_Modified(instructionCode, "LSL")
	}

	if strings.Index(opcode, "11111000000") == 0 { //STUR instruction NOE
		return D_Instruction(instructionCode, "STUR")
	}

	if strings.Index(opcode, "11111000010") == 0 { //LDUR instruction NOE
		return D_Instruction(instructionCode, "LDUR")
	}

	if strings.Index(opcode, "11010011100") == 0 { //ASR instruction
		return R_Instruction_Modified(instructionCode, "ASR")
	}

	if strings.Index(opcode, "00000000000") == 0 { //NOP instruction
		instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), instructionName: "NOP"}
		return fmt.Sprintf("0000000000000000000000000000000\t%d\t%s\n", ProgramCounter, "NOP")
	}

	if strings.Index(opcode, "11101010000") == 0 { //EOR instruction NATHAN
		return R_Instruction(instructionCode, "EOR")
	}

	if strings.Index(opcode, "11111110110") == 0 { //BREAK instruction.
		hasBreak = true
		instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), instructionName: "BREAK"}
		return fmt.Sprintf("%s\t", "1 11111 10110 11110 11111 11111 100111") + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) +
			"BREAK" + "\n"
	}

	return "Unknown Instruction\n"
}

func R_Instruction(instructionCode string, instructionName string) string { // Code for ADD, SUB, AND, ORR
	opcode := instructionCode[0:11]
	rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
	shamt := instructionCode[16:22]
	rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
	rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

	formattedString := fmt.Sprintf("%s\t", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
		shamt+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

	formattedRegisterString := fmt.Sprintf("R%d, R%d, R%d", rd, rs, rt)

	instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), storeRegister: rd, targetRegister1: rs, targetRegister3: rt, instructionName: instructionName}
	return formattedString + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) + instructionName + "\t" + formattedRegisterString + "\n"
}

func R_Instruction_Modified(instructionCode string, instructionName string) string { // Code for LSR, LSL, ASR (needed since these use an immediate instead)
	opcode := instructionCode[0:11]
	rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
	shamt, _ := strconv.ParseInt(instructionCode[16:22], 2, 6)
	rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
	rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

	formattedString := fmt.Sprintf("%s\t", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
		fmt.Sprintf("%06s", strconv.FormatInt(shamt, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

	formattedRegisterString := fmt.Sprintf("R%d, R%d, #%d", rd, rs, shamt)

	instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), storeRegister: rd, targetRegister1: rs, immediateValue: shamt, instructionName: instructionName}
	return formattedString + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) + instructionName + "\t" + formattedRegisterString + "\n"
}

func D_Instruction(instructionCode string, instructionName string) string { //Code for LDUR and STUR
	opcode := instructionCode[0:11]
	immediate := instructionCode[11:20]
	op2 := instructionCode[20:22]
	rt := instructionCode[22:27]
	rd := instructionCode[27:]

	formattedString := fmt.Sprintf("%s\t", opcode+" "+immediate+" "+op2+" "+rt+" "+rd)

	immNum, _ := strconv.ParseInt(immediate, 2, 64)
	rtNum, _ := strconv.ParseInt(rt, 2, 64)
	rdNum, _ := strconv.ParseInt(rd, 2, 64)

	formattedRegisterString := fmt.Sprintf("R%d, [R%d, #%d]", rdNum, rtNum, immNum)

	instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), storeRegister: rdNum, targetRegister1: rtNum, immediateValue: immNum, instructionName: instructionName}
	return formattedString + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) + instructionName + "\t" +
		formattedRegisterString + "\n"
}

func I_Instruction(instructionCode string, instructionName string) string { // Code for ADDI, SUBI
	opcode := instructionCode[0:10]
	immediate := instructionCode[10:22]
	rnNum := instructionCode[22:27]
	rdNum := instructionCode[27:32]

	rn, _ := strconv.ParseInt(rnNum, 2, 5)
	rd, _ := strconv.ParseInt(rdNum, 2, 5)

	decimalNum := ConvertBinaryToDecimal("0" + immediate)

	formattedString := fmt.Sprintf("%s\t", opcode+" "+immediate+" "+rnNum+" "+rdNum)
	formattedRegisterString := fmt.Sprintf("R%d, R%d", rd, rn)

	instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), storeRegister: rd, targetRegister1: rn, immediateValue: decimalNum, instructionName: instructionName}
	return formattedString + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) + instructionName + "\t" +
		formattedRegisterString + ", #" + strconv.FormatInt(decimalNum, 10) + "\n"
}

func IM_Instruction(instructionCode string, instructionName string) string { //Code for MOVK and MOVZ
	opCode := instructionCode[0:9]
	shamtStr := instructionCode[9:11]
	ImediateStr := instructionCode[11:27]
	rdStr := instructionCode[27:32]

	ImediateNum, _ := strconv.ParseInt(ImediateStr, 2, 64)
	rdNum, _ := strconv.ParseInt(rdStr, 2, 64)

	formattedString := fmt.Sprintf("%s %s %s %s", opCode, shamtStr, ImediateStr, rdStr)

	formattedString = fmt.Sprintf("%s\t", formattedString)

	var shamtValue int64 = 100
	shamt, _ := strconv.ParseInt(shamtStr, 2, 2)

	switch shamt {
	case 0:
		shamtValue = 0
	case 1:
		shamtValue = 16
	case 2:
		shamtValue = 32
	case 3:
		shamtValue = 48
	default:
		fmt.Println("Invalid shamt value:", shamt)
		return ""
	}

	instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), storeRegister: rdNum, targetRegister1: ImediateNum, immediateValue: shamtValue, instructionName: instructionName}

	return fmt.Sprintf("%s%d\t%s\tR%s, %d, LSL %d\n", formattedString, ProgramCounter, instructionName, strconv.FormatInt(rdNum, 10), ImediateNum, shamtValue)
}

func BranchInstruction(instructionCode string, instructionName string) string { //Code for CBZ and CBNZ
	opcode := instructionCode[:8]
	immediate := instructionCode[8:27]
	rd := instructionCode[27:]

	immNum, _ := strconv.ParseInt(rd, 2, 32)
	formattedString := fmt.Sprintf("%s\t", opcode+" "+immediate+" "+rd)
	decimalNum := ConvertBinaryToDecimal(immediate)
	formattedRegisterString := fmt.Sprintf("R%d, #%d", immNum, decimalNum)

	instructionMap[int64(ProgramCounter)] = Instruction{programCounter: int64(ProgramCounter), targetRegister1: immNum, immediateValue: decimalNum, instructionName: instructionName}
	return formattedString + fmt.Sprintf("%s\t", strconv.FormatInt(int64(ProgramCounter), 10)) +
		instructionName + "\t" + formattedRegisterString + "\n"
}

// this takes a binary string, 2's complement or not (depending on the first bit) and converts it to decimal.
func ConvertBinaryToDecimal(binaryString string) int64 {
	if binaryString[0] == '0' {
		num, _ := strconv.ParseInt(binaryString, 2, 64)
		return num
	} else {
		num, _ := strconv.ParseInt(binaryString, 2, 64)
		mask := int64(1 << uint(len(binaryString)))
		num -= mask
		return num
	}
}

func resizeRegisterList(registerNum int) {
	for registerNum > maxRegisters {
		maxRegisters += 8
	}
}

func outputRegisters() string {
	outputString := "registers:\n"

	starterNum := 8
	for ; starterNum <= maxRegisters; starterNum += 8 {
		outputString += fmt.Sprintf("r%02d:", starterNum-8)
		for i := starterNum - 8; i < starterNum; i++ {
			outputString += fmt.Sprintf("\t%d", registers[i])
		}
		outputString += "\n"
	}

	return outputString
}

func outputData() string {
	outputString := "data:\n"

	return outputString
}
