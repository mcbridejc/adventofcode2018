package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var AllInstructions = [...]string {
	"addr",
	"addi",
	"mulr",
	"muli",
	"banr",
	"bani",
	"borr",
	"bori",
	"setr",
	"seti",
	"gtir",
	"gtri",
	"gtrr",
	"eqir",
	"eqri",
	"eqrr"}

func Execute(instr Instruction, reg [6]int) [6]int {
	vA := -1
	vB := -1
	A := instr.A
	B := instr.B
	C := instr.C
	instruction := instr.opcode
	if A < 	6 {
		vA = reg[A]
	} 
	if B < 6 {
		vB = reg[B]
	}
	switch instruction {
	case "addr":
		reg[C] = vA + vB
	case "addi":
		reg[C] = vA + B
	case "mulr":
		reg[C] = vA * vB
	case "muli":
		reg[C] = vA * B
	case "banr":
		reg[C] = vA & vB
	case "bani":
		reg[C] = vA & B
	case "borr":
		reg[C] = vA | vB
	case "bori":
		reg[C] = vA | B
	case "setr":
		reg[C] = vA
	case "seti":
		reg[C] = A
	case "gtir":
		reg[C] = 0
		if A > vB { reg[C] = 1 } 
	case "gtri":
		reg[C] = 0
		if vA > B { reg[C] = 1 }
	case "gtrr":
		reg[C] = 0
		if vA > vB { reg[C] = 1 }
	case "eqir":
		reg[C] = 0
		if A == vB { reg[C] = 1 } 
	case "eqri":
		reg[C] = 0
		if vA == B { reg[C] = 1 }
	case "eqrr":
		reg[C] = 0
		if vA == vB { reg[C] = 1 }
	default:
		panic("Unknown instruction")
	}
	return reg
}

type Instruction struct {
	opcode string
	A int
	B int
	C int
}

func ReadInput(filepath string) (int, []Instruction) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	program := make([]Instruction, 0)
	ipReg := -1
	scanner.Scan()
	
	fmt.Sscanf(scanner.Text(), "#ip %d", &ipReg)
	for scanner.Scan() {
		var instr Instruction
		fmt.Sscanf(scanner.Text(), "%s %d %d %d", &instr.opcode, &instr.A, &instr.B, &instr.C)
		program = append(program, instr)
	}
	return ipReg, program
}

func RunProgram(program []Instruction, initialReg [6]int, ipReg int, trace bool, maxCycles int) [6]int {
	reg := initialReg
	loadAddr := reg[ipReg]
	cycle := 0
	for cycle < maxCycles{
		if loadAddr >= len(program) {
			fmt.Println("Halting")
			break
		}
		instr := program[loadAddr]
		if trace {
			fmt.Printf("%03d: PC=%02d %s %d %d %d [%v]\n", cycle, reg[ipReg], instr.opcode, instr.A, instr.B, instr.C, reg)
		}
		reg = Execute(instr, reg)
		reg[ipReg] += 1
		loadAddr = reg[ipReg]
		cycle++
	}
	fmt.Println("Final register values: ", reg)
	return reg
}

func main() {
	trace := flag.Bool("trace", false, "Enable instruction trace output")
	flag.Parse()
	ipReg, program := ReadInput("day19_input.txt")

	fmt.Printf("Assigned reg %d to IP\n", ipReg)
	fmt.Printf("Read program with %d instructions\n", len(program))

	fmt.Println("Part 1\n------")
	reg := [6]int{0, 0, 0, 0, 0, 0}
	reg = RunProgram(program, reg, ipReg, *trace, 100000000)

	fmt.Println("Part 2\n------")
	reg = [6]int{1, 0, 0, 0, 0, 0}
	reg = RunProgram(program, reg, ipReg, *trace, 1000)
	fmt.Println("r5 argument is ", reg[5])

	// See annotated_program.txt and pseudocode.txt
	// Reverse engineering the assembly program shows that the main 
	// loop (labeled by me as 'compute') is summing all the possible factors of
	// the value in r5, by brute force. This isn't feasible. However, 
	// there's a much more efficient way to find theh factors: prime fatorization. 
	// Instead of implement this, I used an online calculator: 
	// https://www.calculatorsoup.com/calculators/math/factors.php?input=10551345&action=solve
}