package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type GuardRecord struct {
	sleepTotals [60]int
	shifts      int
}

func (r *GuardRecord) SleepProbabilityAtMinute(minute int) float32 {
	return float32(r.sleepTotals[minute]) / float32(r.shifts)
}
func (r *GuardRecord) SleepProbabilities() (p [60]float32) {
	for i := 0; i < 60; i += 1 {
		p[i] = r.SleepProbabilityAtMinute(i)
	}
	return
}
func (r *GuardRecord) TotalSleep() (totalSleep int) {
	for _, t := range r.sleepTotals {
		totalSleep += t
	}
	return
}
func (r *GuardRecord) SleepPerNight() float32 {

	return float32(r.TotalSleep()) / float32(r.shifts)
}
func (r *GuardRecord) IncrementShifts() {
	r.shifts += 1
}
func (r *GuardRecord) AddSleep(start_minute int, wake_minute int) {
	for i := start_minute; i < wake_minute; i += 1 {
		r.sleepTotals[i] += 1
	}
}
func (r *GuardRecord) Shifts() int {
	return r.shifts
}
func (r *GuardRecord) SleepTotals() []int {
	return r.sleepTotals[:]
}

func NewGuardRecord() *GuardRecord {
	var r GuardRecord
	r.shifts = 0
	return &r
}

func CheckLine(expr string, s string) (guard_id int, err error) {
	re := regexp.MustCompile(expr)
	match := re.FindStringSubmatch(s)
	if match == nil {
		err = errors.New("No guard found")
		return
	}
	guard_id, err = strconv.Atoi(match[1])
	return
}

func ReadGuardRecords(filepath string) map[int]*GuardRecord {
	f, err := os.Open(filepath)
	check(err)

	scanner := bufio.NewScanner(f)
	var entries []string

	// Read all lines and sort them alphabetically
	// (which has the ultimate effect of sorting them chronologically)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}
	sort.Strings(entries)

	guards := make(map[int]*GuardRecord)

	var last_guard_id int
	var last_sleep_minute int

	for _, e := range entries {
		if guard_id, err := CheckLine("Guard #(\\d*) begins shift", e); err == nil {
			if _, present := guards[guard_id]; !present {
				guards[guard_id] = NewGuardRecord()
			}
			guards[guard_id].IncrementShifts()
			last_guard_id = guard_id
			continue
		}

		if time, err := CheckLine("\\[\\d\\d\\d\\d-\\d\\d-\\d\\d \\d\\d:(\\d\\d)\\] falls asleep", e); err == nil {
			last_sleep_minute = time
			continue
		}

		if time, err := CheckLine("\\[\\d\\d\\d\\d-\\d\\d-\\d\\d \\d\\d:(\\d\\d)\\] wakes up", e); err == nil {
			guards[last_guard_id].AddSleep(last_sleep_minute, time)
		}
	}
	return guards
}

func main() {
	guards := ReadGuardRecords("day4_input.txt")

	for k, v := range guards {
		fmt.Printf("Guard #%d: %d shifts\n", k, v.Shifts())
		for i := 0; i < 60; i += 1 {
			fmt.Printf("%d ", v.SleepTotals()[i])
		}
		fmt.Printf("\n")
		for i := 0; i < 60; i += 1 {
			fmt.Printf("% 2d ", int(v.SleepProbabilityAtMinute(i)*100))
		}
		fmt.Printf("\n")
		fmt.Println(v.SleepPerNight())
	}

	var max_sleep int
	var sleepiest_guard_id int
	for k, v := range guards {
		s := v.TotalSleep() //v.SleepPerNight()
		if s > max_sleep {
			max_sleep = s
			sleepiest_guard_id = k
		}
	}

	var max_probability float32
	var sleepiest_minute int
	for i, p := range guards[sleepiest_guard_id].SleepProbabilities() {
		if p > max_probability {
			max_probability = p
			sleepiest_minute = i
		}
	}

	fmt.Printf("Step 1: Sleepiest: Guard %d at minute %d (%d)\n", sleepiest_guard_id, sleepiest_minute, sleepiest_guard_id*sleepiest_minute)

	max_probability = 0.0
	for guard_id, guard := range guards {
		for minute := 0; minute < 60; minute += 1 {
			p := float32(guard.SleepTotals()[minute])
			if p > max_probability {
				max_probability = p
				sleepiest_guard_id = guard_id
				sleepiest_minute = minute
			}
		}
	}

	fmt.Printf("Step 2: Sleepiest: Guard %d at minute %d (%d)\n", sleepiest_guard_id, sleepiest_minute, sleepiest_guard_id*sleepiest_minute)
}
