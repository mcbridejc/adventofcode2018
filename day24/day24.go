package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type AttackType int

type Group struct {
	id int
	isInfection bool
	units int
	initiative int
	hitpoints int
	attackType string
	attackDamage int
	immunities []string
	weaknesses []string

	selectedTarget *Group
}

func (g *Group) HasWeakness(dtype string) bool {
	for _, t := range g.weaknesses {
		if t == dtype {
			return true
		}
	}
	return false
}

func (g *Group) HasImmunity(dtype string) bool {
	for _, t := range g.immunities {
		if t == dtype {
			return true
		}
	}
	return false
}

func (g *Group) EffectivePower() int {
	return g.units * g.attackDamage
}

// Compute how much damage this group would do to another
func (g *Group) Damage(other *Group) int {
	if other.HasWeakness(g.attackType) {
		return g.EffectivePower() * 2
	} else if other.HasImmunity(g.attackType) {
		return 0
	} else {
		return g.EffectivePower()
	}
}

func ParseLine(line string) (Group, bool) {
	re1 := regexp.MustCompile("(\\d+) units each with (\\d+) hit points (\\(.*\\) )?with an attack that does (\\d+) (\\w+) damage at initiative (\\d+)")
	re2 := regexp.MustCompile("(\\w+) to (.*)")
	match := re1.FindStringSubmatch(line)
	if match == nil {
		return Group{}, false
	}
	group := Group{}
	group.units, _ = strconv.Atoi(match[1])
	group.hitpoints, _ = strconv.Atoi(match[2])
	group.attackDamage, _ = strconv.Atoi(match[4])
	group.attackType = match[5]
	group.initiative, _ = strconv.Atoi(match[6])
	if len(match[3]) > 0 {
		// Remove leading and trailing parens
		s := match[3][1:len(match[3])-2]

		// Split it on any semicolons
		commands := strings.Split(s, ";")
		for _, cmdString := range commands {
			match = re2.FindStringSubmatch(cmdString)
			if match == nil {
				panic("Couldn't match specialty string")
			}
			if match[1] == "weak" {
				group.weaknesses = strings.Split(match[2], ", ")
			} else if match[1] == "immune" {
				group.immunities = strings.Split(match[2], ", ")
			} else {
				panic("Unrecognized specialty")
			}
		}
	}
	return group, true
}

func ReadInput(filepath string) ([]*Group, []*Group) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	infection := make([]*Group, 0)
	immuneSystem := make([]*Group, 0)

	scanner := bufio.NewScanner(f)

	scanner.Scan()

	id := 1
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		}
		group, ok := ParseLine(scanner.Text())
		if !ok {
			panic("Couldn't parse line")
		}
		group.id = id
		id++
		immuneSystem = append(immuneSystem, &group)
	}

	id = 1
	scanner.Scan() // skip a line
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		}
		group, ok := ParseLine(scanner.Text())
		if !ok {
			panic("Couldn't parse line")
		}
		group.isInfection = true
		group.id = id
		id++
		infection = append(infection, &group)
	}

	return immuneSystem, infection
}

func LessBySelectOrder(a *Group, b *Group) bool {
	if a.EffectivePower() == b.EffectivePower() {
		return a.initiative > b.initiative
	} 
	return a.EffectivePower() > b.EffectivePower()
}

func LessByAttackOrder(a *Group, b *Group) bool {
	return a.initiative > b.initiative
}

func SortGroups(groups []*Group, SortFn func(*Group, *Group) bool) []*Group {
	newList := make([]*Group, len(groups))
	copy(newList, groups)
	sort.Slice(newList, func (i, j int) bool {return SortFn(newList[i], newList[j])})
	return newList
}


func DoAttackRound(immune, infection []*Group) ([]*Group, []*Group, int) {
	sortedAll := append(immune, infection...)
	sortedAll = SortGroups(sortedAll, LessBySelectOrder)
	targetable := append(immune, infection...)

	totalUnitsKilled := 0
	chosen := make(map[*Group]bool)
	// Select targets
	for _, chooser := range sortedAll {
		targetIdx := -1
		maxDamage := -1
		for i, candidate := range targetable {
			if candidate.isInfection == chooser.isInfection {
				continue; // They are on the same team
			}
			if _, found := chosen[candidate]; found {
				continue; // This group is already targetted
			}
			appliedDamage := chooser.Damage(candidate)
			if appliedDamage == 0 || appliedDamage < maxDamage {
				continue; // We would do no damage, or less damage than the target we've already chosen
			}
			if appliedDamage > maxDamage {
				targetIdx = i
				maxDamage = appliedDamage
			} else if appliedDamage == maxDamage {
				// We found another candidate with equal damage, sort by tie breaking rules
				if candidate.EffectivePower() > targetable[targetIdx].EffectivePower() {
					targetIdx = i
				} else if candidate.EffectivePower() == targetable[targetIdx].EffectivePower() {
					if candidate.initiative > targetable[targetIdx].initiative {
						targetIdx = i
					}
				}
			} 
		}
		// If a selection was made, save it on the chooser object and remove it
		// from the targetable list because units can only be targeted once in a round
		if targetIdx >= 0 {
			// if chooser.isInfection {
			// 	fmt.Printf("Infection ")
			// } else {
			// 	fmt.Printf("Immune ")
			// }
			// fmt.Printf("group %d chooses target group %d\n", chooser.id, targetable[targetIdx].id)
			chooser.selectedTarget = targetable[targetIdx]
			chosen[targetable[targetIdx]] = true
		}
	}

	// Deal damage
	sortedAll = SortGroups(sortedAll, LessByAttackOrder)
	for _, attacker := range sortedAll {
		if attacker.units == 0 {
			continue
		}
		if attacker.selectedTarget != nil {
			damage := attacker.Damage(attacker.selectedTarget)
			unitsKilled := damage / attacker.selectedTarget.hitpoints
			// if attacker.isInfection {
			// 	fmt.Printf("Infection group %d deals %d to defending group %d, killing %d units\n", 
			// 		attacker.id, damage, attacker.selectedTarget.id, unitsKilled)
			// } else {
			// 	fmt.Printf("Immune group %d deals %d to defending group %d, killing %d units\n",
			// 		attacker.id, damage, attacker.selectedTarget.id, unitsKilled)
			// }
			if unitsKilled >= attacker.selectedTarget.units {
				totalUnitsKilled += attacker.selectedTarget.units
				attacker.selectedTarget.units = 0
				
			} else {
				totalUnitsKilled += unitsKilled
				attacker.selectedTarget.units -= unitsKilled
			}
			// Clear the target
			attacker.selectedTarget = nil
		}
	}

	newImmune := make([]*Group, 0)
	newInfection := make([]*Group, 0)
	// Copy any groups that aren't dead
	for _, g := range immune {
		if g.units > 0 {
			newImmune = append(newImmune, g)
		}
	}
	for _, g := range infection {
		if g.units > 0 {
			newInfection = append(newInfection, g)
		}
	}
	return newImmune, newInfection, totalUnitsKilled
}

func CopyGroups(toCopy []*Group) []*Group {
	ret := make([]*Group, len(toCopy))
	for i, g := range toCopy {
		ret[i] = new(Group)
		*ret[i] = *g
	}
	return ret
}

func main() {
	initImmuneGroups, initInfectionGroups := ReadInput("day24_input.txt")

	immuneGroups := CopyGroups(initImmuneGroups)
	infectionGroups := CopyGroups(initInfectionGroups)

	fmt.Printf("%d immune groups: \n", len(immuneGroups))
	for _, grp := range immuneGroups {
		fmt.Println(grp)
	}

	fmt.Printf("%d infection groups: \n", len(infectionGroups))
	for _, grp := range infectionGroups {
		fmt.Println(grp)
	}

	round := 0
	unitsKilled := 0
	for {
		immuneGroups, infectionGroups, unitsKilled = DoAttackRound(immuneGroups, infectionGroups)
		round++
		fmt.Printf("After round %d: \n  Immune: %d groups\n  Infection: %d groups\n\n", round, len(immuneGroups), len(infectionGroups))
		
		if len(immuneGroups) == 0 || len(infectionGroups) == 0 || unitsKilled == 0 {
			break
		}
	}

	if len(immuneGroups) == 0 {
		units := 0
		for _, g := range infectionGroups {
			units += g.units
		}
		fmt.Printf("Infection wins with %d units remaining\n", units)
	} else if len(infectionGroups) == 0 {
		units := 0
		for _, g := range immuneGroups {
			units += g.units
		}
		fmt.Printf("Immune wins with %d units remaining\n", units)
	} else {
		fmt.Printf("Draw\n")
	}

	fmt.Println("Part 2...finding boost")
	boost := 0
	for {
		// Reset to initial groups
		immuneGroups = CopyGroups(initImmuneGroups)
		infectionGroups = CopyGroups(initInfectionGroups)
		boost++
		for _, g := range immuneGroups {
			g.attackDamage += boost
		}

		fmt.Println("Trying boost ", boost)
		// Run to completion
		unitsKilled := 0
		for {
			immuneGroups, infectionGroups, unitsKilled = DoAttackRound(immuneGroups, infectionGroups)
			
			if len(immuneGroups) == 0 || len(infectionGroups) == 0  || unitsKilled == 0 {
				break
			}
		}
		if len(infectionGroups) == 0 {
			// Immune won!
			units := 0
			for _, g := range immuneGroups {
				units += g.units
			}
			fmt.Printf("Immune wins with %d units remaining\n", units)
			break
		}
	}

	fmt.Println("Required boost is ", boost)
}