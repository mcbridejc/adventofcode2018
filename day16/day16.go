package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func Execute(instruction string, A int, B int, C int, reg [4]int) [4]int {
	vA := -1
	vB := -1
	if A < 4 {
		vA = reg[A]
	} 
	if B < 4 {
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
	opcode int
	A int
	B int
	C int
}

type Example struct {
	opcode int
	A int
	B int
	C int
	initialReg [4]int
	resultReg [4]int
}

func TryExamples(examples []Example) [][]string {
	options := make([][]string, 0)
	for _, ex := range examples {
		consistentInstructions := make([]string, 0)
		for _, instr := range AllInstructions {
			result := Execute(instr, ex.A, ex.B, ex.C, ex.initialReg)
			if result == ex.resultReg {
				consistentInstructions = append(consistentInstructions, instr)
			}
		}
		options = append(options, consistentInstructions)
	}
	return options
}

func ReadInput(filepath string) ([]Example, []Instruction) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	examples := make([]Example, 0)
	program := make([]Instruction, 0)

	scanner.Scan()
	
	for strings.Contains(scanner.Text(), "Before") {
		var nextEx Example
		fmt.Sscanf(scanner.Text(), "Before: [%d, %d, %d, %d]", &nextEx.initialReg[0], &nextEx.initialReg[1], &nextEx.initialReg[2], &nextEx.initialReg[3])
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d %d %d %d", &nextEx.opcode, &nextEx.A, &nextEx.B, &nextEx.C)
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "After:  [%d, %d, %d, %d]", &nextEx.resultReg[0], &nextEx.resultReg[1], &nextEx.resultReg[2], &nextEx.resultReg[3])
		scanner.Scan() 
		scanner.Scan() // consume blank line
		examples = append(examples, nextEx)
	}

	scanner.Scan()
	scanner.Scan() 

	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		var instr Instruction
		fmt.Sscanf(scanner.Text(), "%d %d %d %d", &instr.opcode, &instr.A, &instr.B, &instr.C)
		program = append(program, instr)
	}
	return examples, program
}

func RunTest(instr string, A int, B int, C int, initial [4]int, exp [4]int) {
	var res [4]int
	res = Execute(instr, A, B, C, initial)
	if res != exp {
		fmt.Println("FAIL!", instr, A, B, C, initial, " - > ", res)
		panic("Fail")
	}
}
func Tests() {
	RunTest("addr", 2, 3, 0, [4]int{3, 2, 1, 0}, [4]int{1, 2, 1, 0})
	RunTest("addi", 0, 12, 0, [4]int{3, 9, 9, 9}, [4]int{15, 9, 9, 9})
	RunTest("banr", 0, 1, 2, [4]int{3, 1, 0, 0}, [4]int{3, 1, 1, 0})
	RunTest("banr", 0, 1, 2, [4]int{3, 6, 0, 0}, [4]int{3, 6, 2, 0})
	RunTest("bori", 3, 7, 0, [4]int{1, 2, 3, 4}, [4]int{7, 2, 3, 4})
}

func SolveOpcodes(examples []Example, consistentOptions [][]string) map[int]string {
	reduced := 1000 // arbitrary large number
	opcodeMap := make(map[int]string)
	instrMap := make(map[string]bool)
	for reduced > 0 {
		reduced = 0
		for i := 0; i < len(examples); i++ {
			remainingOptions := make([]string, 0)
			for _, opt := range consistentOptions[i] {
				if _, present := instrMap[opt]; !present {
					remainingOptions = append(remainingOptions, opt)
				}
			}
			if len(remainingOptions) == 1 {
				opcodeMap[examples[i].opcode] = remainingOptions[0]
				instrMap[remainingOptions[0]] = true
				reduced++
			}
		}
	}
	return opcodeMap
}

func main() {

	Tests()
	examples, program := ReadInput("day16_input.txt")

	fmt.Printf("Read %d examples\n", len(examples))
	fmt.Printf("Read program with %d instructions\n", len(program))

	consistentOptions := TryExamples(examples)
	part1Count := 0
	for _, o := range consistentOptions {
		if len(o) >= 3 {
			part1Count++
		}
	}
	fmt.Printf("There are %d examples with 3 consistent instructions\n", part1Count)

	opcodeMap := SolveOpcodes(examples, consistentOptions)
	if len(opcodeMap) == len(AllInstructions) {
		fmt.Println("Successfully identified all opcodes")
	} else {
		panic("Could not identify all opcodes")
	}

	reg := [4]int{0, 0, 0, 0}
	for _, instr := range program {
		reg = Execute(opcodeMap[instr.opcode], instr.A, instr.B, instr.C, reg)
	}

	fmt.Println("Final register values: ", reg)
	

}