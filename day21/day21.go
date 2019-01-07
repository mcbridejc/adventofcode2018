package main

import (
	"bufio"
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

func RunProgram(program []Instruction, initialReg [6]int, ipReg int, trace bool, maxCycles int, breakpoints []int) [6]int {
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
		for _, bp := range breakpoints {
			if loadAddr == bp {
				fmt.Println("Breaking @ ", loadAddr)
				return reg
			}
		}
		cycle++
	}
	fmt.Println("Final register values: ", reg)
	return reg
}

func main() {
	// Simulate program
	// The first value output is the answer to part 1, as this is the way 
	// to terminate the program as quickly as possible
	SimulateProgram(16128384)

	//// Execute the program using the emulated machine from day19, breaking at 
	//// the first r0 == r4 check (addr 28). This is not really necessary, but 
	//// is useful for validating that the golang translation is correct
	// ipReg, program := ReadInput("day21_input.txt")
	// reg := [6]int{0, 0, 0, 0, 0, 0}
	// reg = RunProgram(program, reg, ipReg, true, 100000, []int{28})

	// For part 2, we must assume the function is periodic (which is must be, as
	// its output is limited to 24 bits) and find the last value in the first
	// cycle. Empirically, I found that the first output value (16128384) is 
	// not included in the repeating pattern, so we must find where the repeating 
	// pattern starts. 
	Next := NewSequenceGenerator()
	history := make([]int, 0)
	found := false
	for !found{
		x := Next()
		for _, y := range history {
			if x == y {
				found = true
				fmt.Printf("Found repeat @ %d. last value: %d\n", len(history), history[len(history)-1])
				break
			}
		}
		history = append(history, x)
	}
}
