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
var ProgramCounter int = 96
var hasBreak bool = false

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

		decimalNum := ConvertBinaryToDecimal(number)

		formattedString := fmt.Sprintf("%-38s", opcode+" "+number)
		//38 characters
		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " B   #" + strconv.FormatInt(decimalNum, 10) + "\n"
	}

	if strings.Index(opcode, "10001010000") == 0 { //AND instruction NOE
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt := instructionCode[16:22]
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			shamt+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, R%d", rd, rs, rt)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " AND " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "10001011000") == 0 { //ADD instruction NATHAN
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt := instructionCode[16:22]
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			shamt+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, R%d", rd, rs, rt)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " ADD " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "1001000100") == 0 { //ADDI instruction NATHAN
		opcode = instructionCode[0:10]
		immediate := instructionCode[10:22]
		rnNum := instructionCode[22:27]
		rdNum := instructionCode[27:32]

		rn, _ := strconv.ParseInt(rnNum, 2, 5)
		rd, _ := strconv.ParseInt(rdNum, 2, 5)

		decimalNum := ConvertBinaryToDecimal(immediate)

		formattedString := fmt.Sprintf("%-36s", opcode+" "+immediate+" "+rnNum+" "+rdNum)
		formattedRegisterString := fmt.Sprintf("R%d, R%d", rd, rn)

		return formattedString + "\t" + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + "ADDI " +
			formattedRegisterString + ", #" + strconv.FormatInt(decimalNum, 10) + "\n"

	}

	if strings.Index(opcode, "10101010000") == 0 { //ORR instruction NATHAN
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt := instructionCode[16:22]
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			shamt+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, R%d", rd, rs, rt)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " ORR " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "10110100") == 0 { //CBZ instruction
		opcode = instructionCode[:8]
		immediate := instructionCode[8:27]
		rd := instructionCode[27:]

		immNum, _ := strconv.ParseInt(rd, 2, 32)

		formattedString := fmt.Sprintf("%-38s", opcode+" "+immediate+" "+rd)

		decimalNum := ConvertBinaryToDecimal(immediate)

		formattedRegisterString := fmt.Sprintf("R%d, #%d", immNum, decimalNum)
		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) +
			" CBZ " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "10110101") == 0 { //CBNZ instruction
		opcode = instructionCode[:8]
		immediate := instructionCode[8:27]
		rd := instructionCode[27:]

		immNum, _ := strconv.ParseInt(rd, 2, 32)

		formattedString := fmt.Sprintf("%-38s", opcode+" "+immediate+" "+rd)

		decimalNum := ConvertBinaryToDecimal(immediate)

		formattedRegisterString := fmt.Sprintf("R%d, #%d", immNum, decimalNum)
		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) +
			" CBNZ " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "11001011000") == 0 { //SUB instruction NATHAN
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt := instructionCode[16:22]
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			shamt+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, R%d", rd, rs, rt)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " SUB " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "1101000100") == 0 { //SUBI instruction NATHAN
		opcode = instructionCode[0:10]
		immediate := instructionCode[10:22]
		rnNum := instructionCode[22:27]
		rdNum := instructionCode[27:32]

		rn, _ := strconv.ParseInt(rnNum, 2, 5)
		rd, _ := strconv.ParseInt(rdNum, 2, 5)

		decimalNum := ConvertBinaryToDecimal(immediate)

		formattedString := fmt.Sprintf("%-36s", opcode+" "+immediate+" "+rnNum+" "+rdNum)
		formattedRegisterString := fmt.Sprintf("R%d, R%d", rd, rn)

		return formattedString + "\t" + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + "SUBI " +
			formattedRegisterString + ", #" + strconv.FormatInt(decimalNum, 10) + "\n"
	}

	if strings.Index(opcode, "110100101") == 0 { //MOVZ instruction MAITLAND

	}

	if strings.Index(opcode, "111100101") == 0 { //MOVK instruction MAITLAND

	}

	if strings.Index(opcode, "11010011010") == 0 { //LSR instruction MAITLAND
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt, _ := strconv.ParseInt(instructionCode[16:22], 2, 6)
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			fmt.Sprintf("%06s", strconv.FormatInt(shamt, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, #%d", rd, rs, shamt)
		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + "    LSR    " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "11010011011") == 0 { //LSL instruction MAITLAND
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt, _ := strconv.ParseInt(instructionCode[16:22], 2, 6)
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			fmt.Sprintf("%06s", strconv.FormatInt(shamt, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, #%d", rd, rs, shamt)
		return formattedString + " " + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + "   LSL    " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "11111000000") == 0 { //STUR instruction
		immediate := instructionCode[11:20]
		op2 := instructionCode[20:22]
		rt := instructionCode[22:27]
		rd := instructionCode[27:]

		formattedString := fmt.Sprintf("%-38s", opcode+" "+immediate+" "+op2+" "+rt+" "+rd)

		immNum, _ := strconv.ParseInt(immediate, 2, 64)
		rtNum, _ := strconv.ParseInt(rt, 2, 64)
		rdNum, _ := strconv.ParseInt(rd, 2, 64)

		formattedRegisterString := fmt.Sprintf("R%d, [R%d, #%d]", rdNum, rtNum, immNum)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " STUR " +
			formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "11111000010") == 0 { //LDUR instruction
		immediate := instructionCode[11:20]
		op2 := instructionCode[20:22]
		rt := instructionCode[22:27]
		rd := instructionCode[27:]

		formattedString := fmt.Sprintf("%-38s", opcode+" "+immediate+" "+op2+" "+rt+" "+rd)

		immNum, _ := strconv.ParseInt(immediate, 2, 64)
		rtNum, _ := strconv.ParseInt(rt, 2, 64)
		rdNum, _ := strconv.ParseInt(rd, 2, 64)

		formattedRegisterString := fmt.Sprintf("R%d, [R%d, #%d]", rdNum, rtNum, immNum)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " LDUR " +
			formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "11010011100") == 0 { //ASR instruction
		opcode = instructionCode[0:11]
		rmNum := instructionCode[11:16]
		shamtNum := instructionCode[16:22]
		rnNum := instructionCode[22:27]
		rdNum := instructionCode[27:32]

		shamt, _ := strconv.ParseInt(shamtNum, 2, 5)
		rn, _ := strconv.ParseInt(rnNum, 2, 5)
		rd, _ := strconv.ParseInt(rdNum, 2, 5)

		formattedString := fmt.Sprintf("%-36s", opcode+" "+rmNum+" "+shamtNum+" "+rnNum+" "+rdNum)
		formattedRegisterString := fmt.Sprintf("ASR R%d, R%d, #%d\n", rd, rn, shamt)

		return formattedString + "\t" + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + formattedRegisterString
	}

	if strings.Index(opcode, "00000000000") == 0 { //NOP instruction
		return fmt.Sprintf("%-38s%-4d%s\n", instructionCode, ProgramCounter, "NOP")
	}

	if strings.Index(opcode, "11101010000") == 0 { //EOR instruction NATHAN
		rt, _ := strconv.ParseInt(instructionCode[11:16], 2, 5) //target register
		shamt := instructionCode[16:22]
		rs, _ := strconv.ParseInt(instructionCode[22:27], 2, 5) //source register
		rd, _ := strconv.ParseInt(instructionCode[27:32], 2, 5) //destination register

		formattedString := fmt.Sprintf("%-38s", opcode+" "+fmt.Sprintf("%05s", strconv.FormatInt(rt, 2))+" "+
			shamt+" "+fmt.Sprintf("%05s", strconv.FormatInt(rs, 2))+" "+fmt.Sprintf("%05s", strconv.FormatInt(rd, 2)))

		formattedRegisterString := fmt.Sprintf("R%d, R%d, R%d", rd, rs, rt)

		return formattedString + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + " EOR " + formattedRegisterString + "\n"
	}

	if strings.Index(opcode, "11111110110") == 0 { //BREAK instruction.
		hasBreak = true
		return fmt.Sprintf("%-38s", "11111110 110 11110 11111 11111 100111") + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) +
			" BREAK" + "\n"
	}

	if hasBreak {
		constNum := ConvertBinaryToDecimal(instructionCode)
		instFormat := fmt.Sprintf("%-38s", instructionCode)
		return instFormat + fmt.Sprintf("%-4s", strconv.FormatInt(int64(ProgramCounter), 10)) + fmt.Sprintf("%d\n", constNum)
	}

	return fmt.Sprintf("%-38s%-4s INSTURCTION NOT FOUND\n", instructionCode, strconv.FormatInt(int64(ProgramCounter), 10))
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
