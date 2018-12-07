package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type Task struct {
	name string
	complete bool
	inprogress bool
	deps []*Task
}

func NewTask(name string) (*Task) { 
	var task Task
	task.name = name
	task.deps = make([]*Task, 0)
	return &task
}

type TaskCollection struct {
	tasks []*Task
}

func NewTaskCollection() (*TaskCollection) {
	var tc TaskCollection
	tc.tasks = make([]*Task, 0)
	return &tc
}

func FindOrCreateTask(tc *TaskCollection, name string) *Task {
	for _, t := range tc.tasks { 
		if t.name == name {
			return t
		}
	}

	newTask := NewTask(name)
	tc.tasks = append(tc.tasks, newTask)

	return newTask
}

func ReadTasks(filepath string) *TaskCollection {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	tc := NewTaskCollection()
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile("Step (\\w*) must be finished before step (\\w*) can begin.")
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			panic("Regex fail")
		}

		dependentName := matches[1]
		dependeeName := matches[2]
		dependee := FindOrCreateTask(tc, dependeeName)
		dependent := FindOrCreateTask(tc, dependentName)

		dependee.deps = append(dependee.deps, dependent)
	}
	return tc
}

func GetReadyTasks(tc *TaskCollection) []*Task {
	readyTasks := make([]*Task, 0)
	for _, t := range tc.tasks {
		if t.complete || t.inprogress {
			continue
		}
		ready := true
		for _, dep := range t.deps {
			if !dep.complete {
				ready = false
			}
		}
		if ready {
			readyTasks = append(readyTasks, t)
		}
	}
	sort.Slice(readyTasks, func (i, j int) bool { return readyTasks[i].name < readyTasks[j].name })
	return readyTasks
}

func Part1(tc *TaskCollection) {
	taskLog := ""
	for {
		readyTasks := GetReadyTasks(tc)
		if len(readyTasks) == 0 {
			break
		}
		fmt.Printf("Ready tasks: ")
		for _, t := range readyTasks {
			fmt.Printf("%s ", t.name)
		}
		fmt.Printf("\n")
		taskLog += readyTasks[0].name
		readyTasks[0].complete = true
	}

	fmt.Println("Final task order: ", taskLog)
}

const NumWorkers = 5
const Debug = false
func Part2(tc *TaskCollection) {
	time := 0
	// An array of work-time remaining values for all our elves
	workerLoad := make([]int, NumWorkers)
	workerActiveTask := make([]*Task, NumWorkers)

	for {
		// Queue up tasks as long as there are tasks ready and free workers
		for {
			readyTasks := GetReadyTasks(tc)
			if len(readyTasks) == 0 {
				break
			}
			readyWorkerIdx := -1
			for i, load := range workerLoad {
				if load == 0 {
					readyWorkerIdx = i
					break
				}
			}
			if readyWorkerIdx == -1 {
				break
			}
			taskTime := 60 + (int(readyTasks[0].name[0]) - int('A') + 1)
			if Debug {
				fmt.Printf("Starting task %s for %d\n", readyTasks[0].name, taskTime)
			}
			readyTasks[0].inprogress = true
			workerActiveTask[readyWorkerIdx] = readyTasks[0]
			workerLoad[readyWorkerIdx] = taskTime
		}
		
		if Debug {
			fmt.Printf("%d: ", time)
			for _, t := range workerActiveTask {
				if t == nil {
					fmt.Printf(". ")
				} else {
					fmt.Printf("%s ", t.name)
				}
			}
			fmt.Printf("\n")
		}

		// Decrease work remaining on active tasks, mark them complete as appropriate
		totalRemaining := 0
		for i, _ := range workerLoad {
			totalRemaining += workerLoad[i]
			if workerLoad[i] > 1 {
				workerLoad[i] -= 1
			} else if workerLoad[i] == 1 {
				workerLoad[i] = 0
				workerActiveTask[i].complete = true
				workerActiveTask[i] = nil
			}
		}

		if totalRemaining == 0 {
			fmt.Printf("Complete at t=%d\n", time)
			break
		}

		time += 1
	}

}

func main() {
	tasks := ReadTasks("day7_input.txt")

	fmt.Printf("Read %d tasks\n", len(tasks.tasks))

	fmt.Println("Part 1")
	fmt.Println("------")
	Part1(tasks)

	for _, t := range tasks.tasks {
		t.complete = false
	}

	fmt.Println("Part 2")
	fmt.Println("------")
	Part2(tasks)
}